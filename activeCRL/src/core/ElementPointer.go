package core

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (epPtr *elementPointer) GetElement() Element {
	epPtr.traceableLock()
	defer epPtr.traceableUnlock()
	return epPtr.getElement()
}

// getElement() assumes that all relevant locking is being managed elsewhere
func (epPtr *elementPointer) getElement() Element {
	if epPtr.element == nil && epPtr.getElementIdentifier() != uuid.Nil && epPtr.uOfD != nil {
		epPtr.element = epPtr.uOfD.getElement(epPtr.getElementIdentifier().String())
	}
	return epPtr.element
}

func (epPtr *elementPointer) GetName() string {
	epPtr.traceableLock()
	defer epPtr.traceableUnlock()
	switch epPtr.getElementPointerRole() {
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

// NewAbstractElementPointer() creates and intitializes an elementPointer to play the role of an AbstractElementPointer
func NewAbstractElementPointer(uOfD *UniverseOfDiscourse) ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = ABSTRACT_ELEMENT
	uOfD.AddBaseElement(&ep)
	return &ep
}

// NewRefinedElementPointer() creates and intitializes an elementPointer to play the role of an RefinedElementPointer
func NewRefinedElementPointer(uOfD *UniverseOfDiscourse) ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = REFINED_ELEMENT
	uOfD.AddBaseElement(&ep)
	return &ep
}

// NewOwningElementPointer() creates and intitializes an elementPointer to play the role of an OwningElementPointer
func NewOwningElementPointer(uOfD *UniverseOfDiscourse) ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = OWNING_ELEMENT
	uOfD.AddBaseElement(&ep)
	return &ep
}

// NewReferencedElementPointer() creates and intitializes an elementPointer to play the role of an ReferencedElementPointer
func NewReferencedElementPointer(uOfD *UniverseOfDiscourse) ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = REFERENCED_ELEMENT
	uOfD.AddBaseElement(&ep)
	return &ep
}

// GetElementIdentifier() locks the element pointer and returns the element identifier, releasing the lock in the process
func (epPtr *elementPointer) GetElementIdentifier() uuid.UUID {
	epPtr.traceableLock()
	defer epPtr.traceableUnlock()
	return epPtr.getElementIdentifier()
}

// getElementIdentifier() returns the element identifier without locking
func (epPtr *elementPointer) getElementIdentifier() uuid.UUID {
	return epPtr.elementId
}

func (epPtr *elementPointer) getElementPointerRole() ElementPointerRole {
	return epPtr.elementPointerRole
}

func (epPtr *elementPointer) GetElementVersion() int {
	return epPtr.elementVersion
}

func (epPtr *elementPointer) initializeElementPointer() {
	epPtr.initializePointer()
}

func (bePtr *elementPointer) isEquivalent(be *elementPointer) bool {
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
	return pointerPtr.isEquivalent(&be.pointer)
}

func (elPtr *elementPointer) MarshalJSON() ([]byte, error) {
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
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

func (epPtr *elementPointer) printElementPointer(prefix string) {
	epPtr.printPointer(prefix)
	fmt.Printf("%sIndicated ElementID: %s \n", prefix, epPtr.elementId.String())
	fmt.Printf("%sIndicated ElementVersion: %d \n", prefix, epPtr.elementVersion)
	fmt.Printf("%sElementPointerRole: %d \n", prefix, epPtr.elementPointerRole)
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
func (epPtr *elementPointer) SetElement(element Element) {
	epPtr.traceableLock()
	defer epPtr.traceableUnlock()
	oldElement := epPtr.getElement()
	if oldElement == nil && element == nil {
		return // Nothing to do
	} else if oldElement != nil && element != nil && oldElement.getId() == element.getId() {
		return // Nothing to do
	}
	if element != nil {
		element.traceableLock() // We need to lock the element to make sure it's version doesn't change during this operation
		defer element.traceableUnlock()
	}
	if epPtr.getElementPointerRole() == OWNING_ELEMENT {
		// We have some additional unwiring and wiring to do in this case
		if oldElement != nil {
			oldElement.traceableLock()
			defer oldElement.traceableUnlock()
		}
		if epPtr.getOwningElement() != nil {
			epPtr.getOwningElement().traceableLock()
			defer epPtr.getOwningElement().traceableUnlock()
		}
	}
	epPtr.setElement(element)
}

// setElement() is intended for internal use within the core. It assumes that all relevant
// objects (parent, child, the element pointer itself) have already been locked. All of the
// operations it invokes are also non-locking
func (epPtr *elementPointer) setElement(element Element) {
	if element != epPtr.element {
		// If this is an owningElementPointer, some bookkeeping of the oldOwner
		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.getOwningElement() != nil {
				epPtr.element.removeOwnedBaseElement(epPtr.getOwningElement())
			}
		}

		// Now the actual change of the pointer
		preChange(epPtr)
		epPtr.element = element
		if element != nil {
			epPtr.elementId = element.getId()
			epPtr.elementVersion = element.getVersion()
		} else {
			epPtr.elementId = uuid.Nil
			epPtr.elementVersion = 0
		}
		postChange(epPtr)

		// If this is an owningElementPointer, some bookkeeping of the newOwner
		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.getOwningElement() != nil {
				epPtr.element.addOwnedBaseElement(epPtr.getOwningElement())
			}
		}
	}
}

// SetOwningElement() actually manages relationships between a number of objects,
// particularly when the pointer is the OWNING_ELEMENT pointer for its owner.
// Because of the complex wiring between the objects, we have to lock all relevant
// objects here and then use non-locking worker methods
func (epPtr *elementPointer) SetOwningElement(newOwningElement Element) {
	epPtr.traceableLock()
	defer epPtr.traceableUnlock()
	oldOwningElement := epPtr.getOwningElement()
	if oldOwningElement == nil && newOwningElement == nil {
		return // Nothing to do
	} else if oldOwningElement != nil && newOwningElement != nil && oldOwningElement.getId() == newOwningElement.getId() {
		return // Nothing to do
	}
	if oldOwningElement != nil {
		oldOwningElement.traceableLock()
		defer oldOwningElement.traceableUnlock()
	}
	if newOwningElement != nil {
		newOwningElement.traceableLock()
		defer newOwningElement.traceableUnlock()
	}
	if epPtr.getElementPointerRole() == OWNING_ELEMENT {
		// In this case the element being pointed to will also be impacted
		if epPtr.getElement() != nil {
			epPtr.getElement().traceableLock()
			defer epPtr.getElement().traceableUnlock()
		}
	}
	epPtr.setOwningElement(newOwningElement)
}

// setOwningElement() is a non-locking function that sets the ownership of the element pointer.
// It adjusts the ownedBaseElement set of both the old and new owner. In addition, if it is an
// owningElementPointer, it adjusts the ownedBaseElement set of the owner's owner.
func (epPtr *elementPointer) setOwningElement(element Element) {
	if element != epPtr.getOwningElement() {
		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.getOwningElement() != nil {
				epPtr.element.removeOwnedBaseElement(epPtr.getOwningElement())
			}
		}

		if epPtr.getOwningElement() != nil {
			epPtr.getOwningElement().removeOwnedBaseElement(epPtr)
		}

		preChange(epPtr)
		epPtr.owningElement = element
		postChange(epPtr)

		if epPtr.getOwningElement() != nil {
			epPtr.getOwningElement().addOwnedBaseElement(epPtr)
		}

		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.getOwningElement() != nil {
				epPtr.element.addOwnedBaseElement(epPtr.getOwningElement())
			}
		}
	}
}

// internalSetOwningElement() is an internal function used only when unmarshaling.
func (epPtr *elementPointer) internalSetOwningElement(element Element) {
	if element != epPtr.getOwningElement() {
		epPtr.owningElement = element
		if epPtr.getOwningElement() != nil {
			epPtr.getOwningElement().internalAddOwnedBaseElement(epPtr)
		}

		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.getOwningElement() != nil {
				epPtr.element.internalAddOwnedBaseElement(epPtr.getOwningElement())
			}
		}
	}
}

type ElementPointer interface {
	Pointer
	getElement() Element
	GetElement() Element
	GetElementIdentifier() uuid.UUID
	getElementPointerRole() ElementPointerRole
	GetElementVersion() int
	setElement(Element)
	SetElement(Element)
}
