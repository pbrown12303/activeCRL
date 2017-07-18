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
func abstractionChanged(element Element, notification *ChangeNotification) {
	preChange(element)
	postChange(element, notification)
}

func preChange(be BaseElement) {
	if be != nil && be.getUniverseOfDiscourse().recordingUndo == true {
		be.getUniverseOfDiscourse().markChangedBaseElement(be)
	}
}

func postChange(be BaseElement, notification *ChangeNotification) {
	if TraceChange == true {
		log.Printf("PostChange called")
		PrintNotification(notification)
	}
	if notification.GetDepth() < 10 {
		be.internalIncrementVersion()
		parent := be.getOwningElement()
		//		newNotification := NewChangeNotification(be, MODIFY, notification)
		switch be.(type) {
		case Element:
			for _, function := range GetCore().findFunctions(be.(Element)) {
				go function(be.(Element), notification)
			}
		}
		if parent != nil {
			parent.childChanged(notification)
		}
		switch be.(type) {
		case Element:
			be.getUniverseOfDiscourse().notifyElementListeners(notification)
		}
	}
}

func propagateChange(be BaseElement, notification *ChangeNotification) {
	if notification.GetDepth() < 10 {
		parent := be.getOwningElement()
		//		newNotification := NewChangeNotification(be, MODIFY, notification)
		switch be.(type) {
		case Element:
			for _, function := range GetCore().findFunctions(be.(Element)) {
				go function(be.(Element), notification)
			}
		}
		if parent != nil {
			propagateChange(parent, notification)
		}
		switch be.(type) {
		case Element:
			be.getUniverseOfDiscourse().notifyElementListeners(notification)
		}
	}
}
