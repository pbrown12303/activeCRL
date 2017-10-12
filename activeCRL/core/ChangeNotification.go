// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
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
	underlyingChange *ChangeNotification
}

func NewChangeNotification(baseElement BaseElement, natureOfChange NatureOfChange, underlyingChange *ChangeNotification) *ChangeNotification {
	var notification ChangeNotification
	notification.changedObject = baseElement
	notification.natureOfChange = natureOfChange
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

func (notification *ChangeNotification) Print(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	notificationType := ""
	switch notification.natureOfChange {
	case ADD:
		notificationType = "Add"
	case MODIFY:
		notificationType = "Modify"
	case REMOVE:
		notificationType = "Remove"
	}
	log.Printf("%s%s: \n", prefix, notificationType)
	if notification.changedObject == nil {
		log.Printf(prefix + "Changed object is nil")
	} else {
		log.Printf(prefix + "Changed object is not nil")
		Print(notification.changedObject, prefix+"   ", hl)
	}
	if notification.underlyingChange != nil {
		notification.underlyingChange.Print(prefix+"      ", hl)
	}
	log.Printf(prefix + "End of notification")
}

// abstractionChanged() is used by refinements to inform their refinedElements when they have changed. It does no locking.
func abstractionChanged(element Element, notification *ChangeNotification, hl *HeldLocks) {
	preChange(element, hl)
	postChange(element, notification, hl)
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
//      Notification consists of calling propagateChange() on the pointer
func postChange(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	if notificationsLimit > 0 {
		if notificationsCount > notificationsLimit {
			return
		}
		notificationsCount++
	}
	uOfD := be.GetUniverseOfDiscourse(hl)
	// Increment the version
	be.internalIncrementVersion()
	// Update uri indices
	id := be.GetId(hl)
	oldUri := uOfD.(*universeOfDiscourse).idUriMap.GetEntry(id)
	newUri := GetUri(be, hl)
	if oldUri != newUri {
		if oldUri != "" {
			uOfD.(*universeOfDiscourse).uriBaseElementMap.DeleteEntry(oldUri)
		}
		if newUri == "" {
			uOfD.(*universeOfDiscourse).idUriMap.DeleteEntry(id)
		} else {
			uOfD.(*universeOfDiscourse).idUriMap.SetEntry(id, newUri)
			uOfD.(*universeOfDiscourse).uriBaseElementMap.SetEntry(newUri, be)
		}
	}
	// Initiate function execution
	switch be.(type) {
	case Element:
		for _, labeledFunction := range GetCore().FindFunctions(be.(Element), notification, hl) {
			if TraceChange == true {
				log.Printf("PostChange calling function, URI: %s", labeledFunction.label)
				Print(be, labeledFunction.label+" Target: ", hl)
				notification.Print("Notification: ", hl)
			}
			hl.functionCallManager.AddFunctionCall(labeledFunction, be.(Element), notification)
		}
	}
	// Notify parents of change
	parent := GetOwningElement(be, hl)
	if parent != nil {
		childChanged(parent, notification, hl)
	}
	if uOfD != nil {
		propagateChange(uOfD, notification, hl)
	}
	// Notify listeners
	be.GetUniverseOfDiscourse(hl).notifyBaseElementListeners(notification, hl)
	switch be.(type) {
	case Element:
		uOfD.(*universeOfDiscourse).notifyElementListeners(notification, hl)
	case ElementPointer:
		be.GetUniverseOfDiscourse(hl).notifyElementPointerListeners(notification, hl)
	case Literal:
		be.GetUniverseOfDiscourse(hl).notifyLiteralListeners(notification, hl)
	case LiteralPointer:
		be.GetUniverseOfDiscourse(hl).notifyLiteralPointerListeners(notification, hl)
	}
}

// propagageChange(BaseElement) spreads the knowledge that the base element has changed. It does
// several things:
//   1) if the object is an element it checks to see whether the element has functions associated
//      with it and queues up the calls to the functions.
//   2) if the object is any kind of Pointer it updates the pointer's record of its indicated object version
//   3) It calls propagateChange() on the parent element
//      It also notifies the universe of discourse of the change so that listeners to it become
//      aware of the change
//   4) If the object is of a type that can have a listener (i.e. a pointer to it), the listeners are
//      notified
func propagateChange(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	if notificationsLimit > 0 {
		if notificationsCount > notificationsLimit {
			return
		}
		notificationsCount++
	}
	parent := GetOwningElement(be, hl)
	switch be.(type) {
	case Element:
		for _, labeledFunction := range GetCore().FindFunctions(be.(Element), notification, hl) {
			if TraceChange == true {
				log.Printf("PropagateChange calling function, URI: %s", labeledFunction.label)
				Print(be, labeledFunction.label+" Target: ", hl)
				notification.Print("Notification: ", hl)
			}
			hl.functionCallManager.AddFunctionCall(labeledFunction, be.(Element), notification)
		}
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
	// Notify parent
	if parent != nil {
		propagateChange(parent, notification, hl)
	}
	uOfD := be.GetUniverseOfDiscourse(hl)
	if uOfD != nil {
		// Notify Universe of Discourse
		propagateChange(uOfD, notification, hl)
		// Notify listeners
		uOfD.notifyBaseElementListeners(notification, hl)
		switch be.(type) {
		case Element:
			uOfD.notifyElementListeners(notification, hl)
		case ElementPointer:
			uOfD.notifyElementPointerListeners(notification, hl)
		case Literal:
			uOfD.notifyLiteralListeners(notification, hl)
		case LiteralPointer:
			uOfD.notifyLiteralPointerListeners(notification, hl)
		}

	}
}
