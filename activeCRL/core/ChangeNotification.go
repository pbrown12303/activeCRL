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

// abstractionChanged() is used by refinements to inform their refinedElements when they have changed. It does no locking.
func abstractionChanged(element Element, notification *ChangeNotification, hl *HeldLocks) {
	preChange(element, hl)
	postChange(element, notification, hl)
}

func preChange(be BaseElement, hl *HeldLocks) {
	if be != nil && be.GetUniverseOfDiscourse(hl).undoMgr.recordingUndo == true {
		be.GetUniverseOfDiscourse(hl).undoMgr.markChangedBaseElement(be, hl)
	}
}

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
	id := be.GetId(hl).String()
	oldUri := uOfD.idUriMap.GetEntry(id)
	newUri := GetUri(be, hl)
	if oldUri != newUri {
		if oldUri != "" {
			uOfD.uriBaseElementMap.DeleteEntry(oldUri)
		}
		if newUri == "" {
			uOfD.idUriMap.DeleteEntry(id)
		} else {
			uOfD.idUriMap.SetEntry(id, newUri)
			uOfD.uriBaseElementMap.SetEntry(newUri, be)
		}
	}
	// Initiate function execution
	switch be.(type) {
	case Element:
		for _, labeledFunction := range GetCore().FindFunctions(be.(Element), notification, hl) {
			if TraceChange == true {
				log.Printf("PostChange calling function, URI: %s", labeledFunction.label)
				//				time.Sleep(10000000 * time.Nanosecond)
				Print(be, labeledFunction.label+" Target: ", hl)
				//				time.Sleep(10000000 * time.Nanosecond)
				//				for _, abstraction := range be.(Element).getImmediateAbstractions(hl) {
				//					Print(abstraction, labeledFunction.label+" Abstraction: ", hl)
				//					//					time.Sleep(10000000 * time.Nanosecond)
				//				}
				PrintNotification(notification, "Notification: ", hl)
				//				time.Sleep(1000000 * time.Nanosecond)
			}
			go labeledFunction.function(be.(Element), notification)
		}
	}
	// Notify parents of change
	parent := GetOwningElement(be, hl)
	if parent != nil {
		childChanged(parent, notification, hl)
	}
	// Notify listeners
	switch be.(type) {
	case Element:
		uOfD.notifyElementListeners(notification, hl)
	}
}

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
				//				time.Sleep(10000000 * time.Nanosecond)
				Print(be, labeledFunction.label+" Target: ", hl)
				//				time.Sleep(100000000 * time.Nanosecond)
				//				for _, abstraction := range be.(Element).getImmediateAbstractions(hl) {
				//					Print(abstraction, labeledFunction.label+" Abstraction: ", hl)
				//					//					time.Sleep(10000000 * time.Nanosecond)
				//				}
				PrintNotification(notification, "Notification: ", hl)
				//				time.Sleep(10000000 * time.Nanosecond)
			}
			go labeledFunction.function(be.(Element), notification)
		}
	case ElementPointer:
		ep := be.(ElementPointer)
		target := notification.changedObject
		ep.setElementVersion(target.GetVersion(hl), hl)
	}
	if parent != nil {
		propagateChange(parent, notification, hl)
	}
	switch be.(type) {
	case Element:
		be.GetUniverseOfDiscourse(hl).notifyElementListeners(notification, hl)
	}
}
