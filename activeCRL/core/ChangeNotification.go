package core

import ()

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
	newNotification := NewChangeNotification(element, MODIFY, notification)
	postChange(element, newNotification)
}

func preChange(be BaseElement) {
	if be != nil && be.getUniverseOfDiscourse().recordingUndo == true {
		be.getUniverseOfDiscourse().markChangedBaseElement(be)
	}
}

func postChange(be BaseElement, notification *ChangeNotification) {
	be.internalIncrementVersion()
	parent := be.getOwningElement()
	newNotification := NewChangeNotification(be, MODIFY, notification)
	switch be.(type) {
	case Element:
		for _, function := range GetCore().findFunctions(be.(Element)) {
			go function(be.(Element), newNotification)
		}
	}
	if parent != nil {
		parent.childChanged(newNotification)
	}
	switch be.(type) {
	case Element:
		be.getUniverseOfDiscourse().notifyElementListeners(newNotification)
	}
}

func propagateChange(be BaseElement, notification *ChangeNotification) {
	parent := be.getOwningElement()
	newNotification := NewChangeNotification(be, MODIFY, notification)
	switch be.(type) {
	case Element:
		for _, function := range GetCore().findFunctions(be.(Element)) {
			go function(be.(Element), newNotification)
		}
	}
	if parent != nil {
		propagateChange(parent, newNotification)
	}
	switch be.(type) {
	case Element:
		be.getUniverseOfDiscourse().notifyElementListeners(newNotification)
	}
}
