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
				//				time.Sleep(10000000 * time.Nanosecond)
				Print(be, labeledFunction.label+" Target: ", hl)
				//				time.Sleep(10000000 * time.Nanosecond)
				//				for _, abstraction := range be.(Element).getImmediateAbstractions(hl) {
				//					Print(abstraction, labeledFunction.label+" Abstraction: ", hl)
				//					//					time.Sleep(10000000 * time.Nanosecond)
				//				}
				notification.Print("Notification: ", hl)
				//				time.Sleep(1000000 * time.Nanosecond)
			}
			hl.functionCallManager.AddFunctionCall(labeledFunction, be.(Element), notification)
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
		uOfD.(*universeOfDiscourse).notifyElementListeners(notification, hl)
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
				notification.Print("Notification: ", hl)
				//				time.Sleep(10000000 * time.Nanosecond)
			}
			hl.functionCallManager.AddFunctionCall(labeledFunction, be.(Element), notification)
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
		be.GetUniverseOfDiscourse(hl).(*universeOfDiscourse).notifyElementListeners(notification, hl)
	}
}
