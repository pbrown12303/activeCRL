// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"os"
	"strconv"
	//	"time"
)

var TraceChange bool
var notificationsLimit int
var notificationsCount int

type NatureOfChange int

const (
	ADD NatureOfChange = iota
	MODIFY
	REMOVE
)

type ChangeNotification struct {
	changedObject    BaseElement
	natureOfChange   NatureOfChange
	origin           string
	underlyingChange *ChangeNotification
}

func NewChangeNotification(baseElement BaseElement, natureOfChange NatureOfChange, origin string, underlyingChange *ChangeNotification) *ChangeNotification {
	var notification ChangeNotification
	notification.changedObject = baseElement
	notification.natureOfChange = natureOfChange
	notification.origin = origin
	notification.underlyingChange = underlyingChange
	return &notification
}

// LimitNotifications() is provided as a debugging aid. It limits the number of change notifications allowed.
// A value of 0 is unlimited.
func LimitNotifications(limit int) {
	notificationsLimit = limit
	notificationsCount = 0
}

func (cnPtr *ChangeNotification) GetDepth() int {
	return cnPtr.getDepth(0)
}

func (cnPtr *ChangeNotification) getDepth(currentDepth int) int {
	newDepth := currentDepth + 1
	if cnPtr.underlyingChange != nil {
		return cnPtr.underlyingChange.getDepth(newDepth)
	}
	return newDepth
}

func (cnPtr *ChangeNotification) isReferenced(be BaseElement) bool {
	if cnPtr.changedObject == be {
		return true
	} else if cnPtr.underlyingChange != nil {
		return cnPtr.underlyingChange.isReferenced(be)
	}
	return false
}

func (cnPtr *ChangeNotification) GetChangedBaseElement() BaseElement {
	return cnPtr.changedObject
}

func (cnPtr *ChangeNotification) getReferencingChangeNotification(be BaseElement) *ChangeNotification {
	if cnPtr.changedObject == be {
		return cnPtr
	} else {
		if cnPtr.underlyingChange != nil {
			return cnPtr.underlyingChange.getReferencingChangeNotification(be)
		}
	}
	return nil
}

func (cnPtr *ChangeNotification) GetUnderlyingChangeNotification() *ChangeNotification {
	return cnPtr.underlyingChange
}

func (notification *ChangeNotification) Print(prefix string, hl *HeldLocks) {
	startCount := 0
	notification.PrintRecursively(prefix, hl, startCount)
}

func (notification *ChangeNotification) PrintRecursively(prefix string, hl *HeldLocks, startCount int) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	notificationType := ""
	switch notification.natureOfChange {
	case ADD:
		notificationType = "+++ Add"
	case MODIFY:
		notificationType = "~~~ Modify"
	case REMOVE:
		notificationType = "--- Remove"
	}
	log.Printf("%s%s: \n", prefix, "### Notification Level: "+strconv.Itoa(startCount)+" Type: "+notificationType)
	log.Printf("%s Origin: %s \n", prefix, notification.origin)
	if notification.changedObject == nil {
		log.Printf(prefix + "Changed object is nil")
	} else {
		log.Printf(prefix + "Changed object is not nil")
		log.Printf(prefix+"  Type: %T", notification.changedObject)
		log.Printf(prefix+"  Id: %s", notification.changedObject.GetId(hl))
		log.Printf(prefix+"  Version: %d", notification.changedObject.GetVersion(hl))
		//		Print(notification.changedObject, prefix+"   ", hl)
	}
	if notification.underlyingChange != nil {
		notification.underlyingChange.PrintRecursively(prefix+"      ", hl, startCount-1)
	}
	log.Printf(prefix + "End of notification")
}

// abstractionChanged() is used by refinements to inform their refinedElements when they have changed. It does no locking.
func abstractionChanged(element Element, notification *ChangeNotification, hl *HeldLocks) {
	preChange(element, hl)
	postChange(element, notification, hl)
}

// childChanged() is used by ownedBaseElements to inform their parents when they have changed. It does no locking.
func childChanged(el Element, notification *ChangeNotification, hl *HeldLocks) {
	if TraceChange == true {
		log.Printf("childChanged called on %T identifier: %s \n", el, el.GetId(hl))
		notification.Print("ChildChanged Incoming Notification: ", hl)
	}
	// First check for circular notifications. We do not want to propagate these
	if notification.isReferenced(el) {
		return
	}
	preChange(el, hl)
	newNotification := NewChangeNotification(el, MODIFY, "childChanged", notification)
	switch el.(type) {
	case Refinement:
		refinedElement := el.(Refinement).GetRefinedElement(hl)
		refinedElementPointer := el.(Refinement).GetRefinedElementPointer(hl)
		if refinedElement != nil {
			cn := notification.getReferencingChangeNotification(refinedElementPointer)
			if cn != nil && cn.underlyingChange == nil {
				abstractionChanged(refinedElement, newNotification, hl)
			}
		}
	}
	postChange(el, newNotification, hl)
}

func notifyListeners(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	underlyingChange := notification.underlyingChange
	if underlyingChange != nil {
		if underlyingChange.isReferenced(notification.changedObject) {
			return
		}
	}
	uOfD := be.GetUniverseOfDiscourse(hl)
	if uOfD != nil {
		uOfD.notifyListeners(be, notification, hl)
	}
}

func notifyParent(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	parent := GetOwningElement(be, hl)
	if parent != nil {
		childChanged(parent, notification, hl)
	}
}

func notifyUniverseOfDiscourse(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	uOfD := be.GetUniverseOfDiscourse(hl)
	if uOfD != nil {
		if uOfD != be {
			uOfD.uOfDChanged(notification, hl)
		}
	}

}

func queueFunctionExecutions(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	switch be.(type) {
	case Element:
		functionIdentifiers := GetCore().FindFunctions(be.(Element), notification, hl)
		for _, functionIdentifier := range functionIdentifiers {
			if TraceChange == true {
				log.Printf("queueFunctionExecutions calling function, URI: %s", functionIdentifier)
				Print(be, string(functionIdentifier)+"Function Target: ", hl)
				notification.Print("Notification: ", hl)
			}
			hl.functionCallManager.AddFunctionCall(functionIdentifier, be.(Element), notification)
		}
	}

}

func preChange(be BaseElement, hl *HeldLocks) {
	if be != nil && be.GetUniverseOfDiscourse(hl).IsRecordingUndo() == true {
		be.GetUniverseOfDiscourse(hl).(*universeOfDiscourse).undoMgr.markChangedBaseElement(be, hl)
	}
}

// This function is called after a change is made to a base element. It does several things:
//   1) It increments the base element version
//   2) If the uri is changed, it updates the indexes associated with uri's
//   3) If there are associated functions, it queues up the call to those functions
//   4) It notifies the parent (owningElement) that this ownedElement has changed
//      It also notifies the universe of discourse of the change so that listeners to it become
//      aware of the change
//   5) If there are any pointers to the base element it notifies them of the change.
//      Notification consists of calling indicatedBaseElementChanged() on the pointer
func postChange(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	if TraceChange == true {
		log.Printf("In post change, %T identifier: %s \n", be, be.getIdNoLock())
		notification.Print("PostChange Incoming Notification: ", hl)
	}
	if notificationsLimit > 0 {
		if notificationsCount > notificationsLimit {
			return
		}
		notificationsCount++
	}

	// Increment the version
	be.internalIncrementVersion()
	// Update uri indices
	updateUriIndices(be, hl)
	// Initiate function execution
	queueFunctionExecutions(be, notification, hl)
	// Notify parents of change
	notifyParent(be, notification, hl)
	// Notify Universe of Discourse
	notifyUniverseOfDiscourse(be, notification, hl)
	// Notify listeners
	notifyListeners(be, notification, hl)
}

// indicatedBaseElementChanged(BaseElement) spreads the knowledge that the base element has changed. It does
// several things:
//   1) if the receiving object is an element it checks to see whether the element has functions associated
//      with it and queues up the calls to the functions.
//   2) if the object is any kind of Pointer it updates the pointer's record of its indicated object version
//   3) It calls indicatedBaseElementChanged() on the parent element
//      It also notifies the universe of discourse of the change so that listeners to it become
//      aware of the change
//   4) If the object is of a type that can have a listener (i.e. a pointer to it), the listeners are
//      notified
func indicatedBaseElementChanged(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	if TraceChange == true {
		log.Printf("In indicatedBaseElementChanged, be identifier: %s \n", be.getIdNoLock())
		notification.Print("indicatedBaseElementChanged Incoming Notification: ", hl)
		var filename string = "NotificationGraph" + strconv.Itoa(notificationsCount)
		file, err := os.Create(filename)
		if err != nil {
			log.Printf("Error: %s", err)
		}
		nGraph := NewNotificationGraph(notification, hl)
		file.WriteString(nGraph.getGraph().String())
	}
	if notificationsLimit > 0 {
		if notificationsCount > notificationsLimit {
			return
		}
		notificationsCount++
	}
	if AdHocTrace == true {
		if notificationsCount > 900000 {
			TraceChange = true
		}
	}

	// Suppress circular notifications
	underlyingChange := notification.underlyingChange
	if underlyingChange != nil {
		if underlyingChange.isReferenced(notification.changedObject) {
			return
		}
	}

	// Initiate function executions
	queueFunctionExecutions(be, notification, hl)
	// For pointers, update the version number of the referenced object. This does not count as a change
	// to the pointer itself
	updatePointerVersions(be, notification, hl)
	// Notify parent
	notifyParent(be, notification, hl)
	// Notify Universe of Discourse
	notifyUniverseOfDiscourse(be, notification, hl)
	// Notify listeners
	notifyListeners(be, notification, hl)
}

func updatePointerVersions(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	switch be.(type) {
	case ElementPointer:
		ep := be.(ElementPointer)
		target := notification.changedObject
		ep.setElementVersion(target.GetVersion(hl), hl)
	case ElementPointerPointer:
		ep := be.(ElementPointerPointer)
		target := notification.changedObject
		ep.setElementPointerVersion(target.GetVersion(hl), hl)
	case LiteralPointer:
		ep := be.(LiteralPointer)
		target := notification.changedObject
		ep.setLiteralVersion(target.GetVersion(hl), hl)
	case LiteralPointerPointer:
		ep := be.(LiteralPointerPointer)
		target := notification.changedObject
		ep.setLiteralPointerVersion(target.GetVersion(hl), hl)
	}
}

func updateUriIndices(be BaseElement, hl *HeldLocks) {
	uOfD := be.GetUniverseOfDiscourse(hl)
	if uOfD != nil {
		uOfD.updateUriIndices(be, hl)
	}
}
