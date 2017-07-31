package core

import (
	"log"
)

var TraceChange bool

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
	if TraceChange == true {
		log.Printf("PostChange called")
		PrintNotification(notification, hl)
	}
	if notification.GetDepth() < 10 {
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
			for _, function := range GetCore().findFunctions(be.(Element), hl) {
				go function(be.(Element), notification)
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
	if TraceChange == true {
		log.Printf("At end of PostChange, changed object:")
		Print(notification.changedObject, "", hl)
	}
}

func propagateChange(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	if notification.GetDepth() < 10 {
		parent := GetOwningElement(be, hl)
		//		newNotification := NewChangeNotification(be, MODIFY, notification)
		switch be.(type) {
		case Element:
			for _, function := range GetCore().findFunctions(be.(Element), hl) {
				go function(be.(Element), notification)
			}
		}
		if parent != nil {
			propagateChange(parent, notification, hl)
		}
		switch be.(type) {
		case Element:
			be.GetUniverseOfDiscourse(hl).notifyElementListeners(notification, hl)
		}
	}
}
