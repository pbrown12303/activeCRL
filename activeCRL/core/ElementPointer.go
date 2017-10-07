package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/satori/go.uuid"
)

type ElementPointerRole int

const (
	ABSTRACT_ELEMENT ElementPointerRole = 1 + iota
	REFINED_ELEMENT
	OWNING_ELEMENT
	REFERENCED_ELEMENT
)

func (epr ElementPointerRole) RoleToString() string {
	role := ""
	switch epr {
	case ABSTRACT_ELEMENT:
		role = "ABSTRACT_ELEMENT"
	case REFINED_ELEMENT:
		role = "REFINED_ELEMENT"
	case OWNING_ELEMENT:
		role = "OWNING_ELEMENT"
	case REFERENCED_ELEMENT:
		role = "REFERENCED_ELEMENT"
	}
	return role
}

type elementPointer struct {
	pointer
	element            Element
	elementId          uuid.UUID
	elementVersion     int
	elementPointerRole ElementPointerRole
}

func (epPtr *elementPointer) clone() *elementPointer {
	var ep elementPointer
	ep.cloneAttributes(*epPtr)
	return &ep
}

func (epPtr *elementPointer) cloneAttributes(source elementPointer) {
	epPtr.pointer.cloneAttributes(source.pointer)
	epPtr.element = source.element
	epPtr.elementId = source.elementId
	epPtr.elementVersion = source.elementVersion
	epPtr.elementPointerRole = source.elementPointerRole

}

func (epPtr *elementPointer) elementChanged(notification *ChangeNotification, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	// Circular references need to be detected and curtailed, hence the isReferenced() call
	if epPtr.getOwningElement(hl) != nil && notification.isReferenced(epPtr) == false {
		newNotification := NewChangeNotification(epPtr, MODIFY, notification)
		childChanged(epPtr.getOwningElement(hl), newNotification, hl)
	}

}

func (epPtr *elementPointer) GetElement(hl *HeldLocks) Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	if epPtr.element == nil && epPtr.GetElementId(hl) != uuid.Nil && epPtr.uOfD != nil {
		epPtr.element = epPtr.uOfD.GetElement(epPtr.GetElementId(hl).String())
	}
	return epPtr.element
}

func (epPtr *elementPointer) getName(hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	switch epPtr.GetElementPointerRole(hl) {
	case ABSTRACT_ELEMENT:
		return "abstractElement"
	case REFINED_ELEMENT:
		return "refinedElement"
	case OWNING_ELEMENT:
		return "owningElement"
	case REFERENCED_ELEMENT:
		return "referencedElement"
	}
	return ""
}

// GetElementIdentifier() locks the element pointer and returns the element identifier, releasing the lock in the process
func (epPtr *elementPointer) GetElementId(hl *HeldLocks) uuid.UUID {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	return epPtr.elementId
}

func (epPtr *elementPointer) GetElementPointerRole(hl *HeldLocks) ElementPointerRole {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	return epPtr.elementPointerRole
}

func (epPtr *elementPointer) GetElementVersion(hl *HeldLocks) int {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	return epPtr.elementVersion
}

func (epPtr *elementPointer) initializeElementPointer() {
	epPtr.initializePointer()
}

func (bePtr *elementPointer) isEquivalent(be *elementPointer, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	if bePtr.elementId != be.elementId {
		fmt.Printf("Equivalence failed: indicated element ids do not match \n")
		return false
	}
	if bePtr.elementVersion != be.elementVersion {
		fmt.Printf("Equivalence failed: indicated element versions do not match \n")
		return false
	}
	if bePtr.elementPointerRole != be.elementPointerRole {
		fmt.Printf("Equivalence failed: element pointer roles do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer, hl)
}

func (elPtr *elementPointer) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalElementPointerFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *elementPointer) marshalElementPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.pointer.marshalPointerFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"ElementId\":\"%s\",", elPtr.elementId.String()))
	buffer.WriteString(fmt.Sprintf("\"ElementVersion\":\"%d\",", elPtr.elementVersion))
	switch elPtr.elementPointerRole {
	case ABSTRACT_ELEMENT:
		buffer.WriteString(fmt.Sprintf("\"ElementPointerRole\":\"%s\"", "ABSTRACT_ELEMENT"))
	case REFINED_ELEMENT:
		buffer.WriteString(fmt.Sprintf("\"ElementPointerRole\":\"%s\"", "REFINED_ELEMENT"))
	case OWNING_ELEMENT:
		buffer.WriteString(fmt.Sprintf("\"ElementPointerRole\":\"%s\"", "OWNING_ELEMENT"))
	case REFERENCED_ELEMENT:
		buffer.WriteString(fmt.Sprintf("\"ElementPointerRole\":\"%s\"", "REFERENCED_ELEMENT"))
	}
	return err
}

func (epPtr *elementPointer) printElementPointer(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	epPtr.printPointer(prefix, hl)
	log.Printf("%s  Indicated ElementID: %s \n", prefix, epPtr.elementId.String())
	log.Printf("%s  Indicated ElementVersion: %d \n", prefix, epPtr.elementVersion)
	role := ""
	switch epPtr.elementPointerRole {
	case ABSTRACT_ELEMENT:
		role = "ABSTRACT_ELEMENT"
	case REFINED_ELEMENT:
		role = "REFINED_ELEMENT"
	case OWNING_ELEMENT:
		role = "OWNING_ELEMENT"
	case REFERENCED_ELEMENT:
		role = "REFERENCED_ELEMENT"
	}
	log.Printf("%s  ElementPointerRole: %s \n", prefix, role)
}

func (ep *elementPointer) recoverElementPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	err := ep.pointer.recoverPointerFields(unmarshaledData)
	if err != nil {
		fmt.Printf("ElementPointer's Recovery of PointerFields failed\n")
		return err
	}
	// Element ID
	var recoveredElementId string
	err = json.Unmarshal((*unmarshaledData)["ElementId"], &recoveredElementId)
	if err != nil {
		fmt.Printf("ElementPointer's Recovery of ElementId failed\n")
		return err
	}
	ep.elementId, err = uuid.FromString(recoveredElementId)
	if err != nil {
		fmt.Printf("ElementPointer's conversion of ElementId failed\n")
		return err
	}
	// Version
	var recoveredElementVersion string
	err = json.Unmarshal((*unmarshaledData)["ElementVersion"], &recoveredElementVersion)
	if err != nil {
		fmt.Printf("ElementPointer's Recovery of ElementVersion failed\n")
		return err
	}
	ep.elementVersion, err = strconv.Atoi(recoveredElementVersion)
	if err != nil {
		fmt.Printf("Conversion of ElementPointer.elementVersion failed\n")
		return err
	}
	// Element pointer role
	var recoveredElementPointerRole string
	err = json.Unmarshal((*unmarshaledData)["ElementPointerRole"], &recoveredElementPointerRole)
	if err != nil {
		fmt.Printf("ElementPointer's Recovery of ElementPointerRole failed\n")
		return err
	}
	switch recoveredElementPointerRole {
	case "ABSTRACT_ELEMENT":
		ep.elementPointerRole = ABSTRACT_ELEMENT
	case "REFINED_ELEMENT":
		ep.elementPointerRole = REFINED_ELEMENT
	case "OWNING_ELEMENT":
		ep.elementPointerRole = OWNING_ELEMENT
	case "REFERENCED_ELEMENT":
		ep.elementPointerRole = REFERENCED_ELEMENT
	}
	return nil
}

// SetElement() establishes the element to which this pointer points. If this pointer
// happens to be an OWNING_ELEMENT pointer, there is a side-effect in which this pointer's
// owner is removed as a child from the old target element and added as a child to the new
// target element. Locking must take this into account.
func (epPtr *elementPointer) SetElement(element Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	oldElement := epPtr.GetElement(hl)
	if oldElement == nil && element == nil {
		return // Nothing to do
	} else if oldElement != nil && element != nil && oldElement.GetId(hl) == element.GetId(hl) {
		return // Nothing to do
	}
	if element != nil {
		hl.LockBaseElement(element)
	}
	// If this is an owningElementPointer, some bookkeeping of the oldOwner
	if epPtr.GetElementPointerRole(hl) == OWNING_ELEMENT {
		if oldElement != nil && epPtr.getOwningElement(hl) != nil {
			removeOwnedBaseElement(epPtr.element, epPtr.getOwningElement(hl), hl)
		}
	}

	// Now the actual change of the pointer
	preChange(epPtr, hl)
	if oldElement != nil {
		epPtr.uOfD.removeElementListener(oldElement, epPtr, hl)
	}
	epPtr.element = element
	if element != nil {
		epPtr.elementId = element.GetId(hl)
		epPtr.elementVersion = element.GetVersion(hl)
		epPtr.uOfD.addElementListener(element, epPtr, hl)
	} else {
		epPtr.elementId = uuid.Nil
		epPtr.elementVersion = 0
	}
	notification := NewChangeNotification(epPtr, MODIFY, nil)
	postChange(epPtr, notification, hl)

	// If this is an owningElementPointer, some bookkeeping of the newOwner
	if epPtr.GetElementPointerRole(hl) == OWNING_ELEMENT {
		if epPtr.element != nil && epPtr.getOwningElement(hl) != nil {
			addOwnedBaseElement(epPtr.element, epPtr.getOwningElement(hl), hl)
		}
	}
}

// setElementVersion() is an internal function used as part of change propagation. Id does
// not trigger any notifications
func (epPtr *elementPointer) setElementVersion(newVersion int, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	epPtr.elementVersion = newVersion
}

// SetOwningElement() actually manages relationships between a number of objects,
// particularly when the pointer is the OWNING_ELEMENT pointer for its owner.
// Because of the complex wiring between the objects, we have to lock all relevant
// objects here and then use non-locking worker methods
func (epPtr *elementPointer) SetOwningElement(newOwningElement Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	oldOwningElement := epPtr.getOwningElement(hl)
	if oldOwningElement == nil && newOwningElement == nil {
		return // Nothing to do
	} else if oldOwningElement != nil && newOwningElement != nil && oldOwningElement.GetId(hl) == newOwningElement.GetId(hl) {
		return // Nothing to do
	}
	if epPtr.GetElementPointerRole(hl) == OWNING_ELEMENT {
		if epPtr.element != nil && oldOwningElement != nil {
			removeOwnedBaseElement(epPtr.element, oldOwningElement, hl)
		}
	}

	if oldOwningElement != nil {
		removeOwnedBaseElement(oldOwningElement, epPtr, hl)
	}

	preChange(epPtr, hl)
	epPtr.owningElement = newOwningElement
	notification := NewChangeNotification(epPtr, MODIFY, nil)
	postChange(epPtr, notification, hl)

	if newOwningElement != nil {
		addOwnedBaseElement(newOwningElement, epPtr, hl)
	}

	if epPtr.GetElementPointerRole(hl) == OWNING_ELEMENT {
		if epPtr.element != nil && newOwningElement != nil {
			addOwnedBaseElement(epPtr.element, newOwningElement, hl)
		}
	}
}

// internalSetOwningElement() is an internal function used only when unmarshaling.
func (epPtr *elementPointer) internalSetOwningElement(element Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	if element != epPtr.getOwningElement(hl) {
		epPtr.owningElement = element
		if epPtr.getOwningElement(hl) != nil {
			epPtr.getOwningElement(hl).internalAddOwnedBaseElement(epPtr, hl)
		}

		if epPtr.GetElementPointerRole(hl) == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.getOwningElement(hl) != nil {
				epPtr.element.internalAddOwnedBaseElement(epPtr.getOwningElement(hl), hl)
			}
		}
	}
}

func (epPtr *elementPointer) setUri(uri string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	preChange(epPtr, hl)
	epPtr.uri = uri
	notification := NewChangeNotification(epPtr, MODIFY, nil)
	postChange(epPtr, notification, hl)
}

type ElementPointer interface {
	Pointer
	elementChanged(*ChangeNotification, *HeldLocks)
	GetElement(*HeldLocks) Element
	GetElementId(*HeldLocks) uuid.UUID
	GetElementPointerRole(*HeldLocks) ElementPointerRole
	GetElementVersion(*HeldLocks) int
	SetElement(Element, *HeldLocks)
	setElementVersion(int, *HeldLocks)
}
