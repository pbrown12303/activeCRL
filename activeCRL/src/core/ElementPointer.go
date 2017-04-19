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

func (epPtr *elementPointer) GetElement() Element {
	if epPtr.element == nil && epPtr.GetElementIdentifier() != uuid.Nil && epPtr.uOfD != nil {
		epPtr.element = epPtr.uOfD.getElement(epPtr.GetElementIdentifier().String())
	}
	return epPtr.element
}

func (epPtr *elementPointer) GetName() string {
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

func NewAbstractElementPointer(uOfD *UniverseOfDiscourse) ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.setElementPointerRole(ABSTRACT_ELEMENT)
	uOfD.addBaseElement(&ep)
	return &ep
}

func NewRefinedElementPointer(uOfD *UniverseOfDiscourse) ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.setElementPointerRole(REFINED_ELEMENT)
	uOfD.addBaseElement(&ep)
	return &ep
}

func NewOwningElementPointer(uOfD *UniverseOfDiscourse) ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.setElementPointerRole(OWNING_ELEMENT)
	uOfD.addBaseElement(&ep)
	return &ep
}

func NewReferencedElementPointer(uOfD *UniverseOfDiscourse) ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.setElementPointerRole(REFERENCED_ELEMENT)
	uOfD.addBaseElement(&ep)
	return &ep
}

func (epPtr *elementPointer) GetElementIdentifier() uuid.UUID {
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

func (epPtr *elementPointer) SetElement(element Element) {
	if element != epPtr.element {
		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.GetOwningElement() != nil {
				epPtr.element.removeOwnedBaseElement(epPtr.GetOwningElement())
			}
		}
		epPtr.element = element
		if element != nil {
			epPtr.elementId = element.GetId()
			epPtr.elementVersion = element.GetVersion()
		} else {
			epPtr.elementId = uuid.Nil
			epPtr.elementVersion = 0
		}
		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.GetOwningElement() != nil {
				epPtr.element.addOwnedBaseElement(epPtr.GetOwningElement())
			}
		}
	}
}

func (epPtr *elementPointer) setElementPointerRole(role ElementPointerRole) {
	epPtr.elementPointerRole = role
}

func (epPtr *elementPointer) setOwningElement(element Element) {
	if element != epPtr.GetOwningElement() {
		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.GetOwningElement() != nil {
				epPtr.element.removeOwnedBaseElement(epPtr.GetOwningElement())
			}
		}
		if epPtr.GetOwningElement() != nil {
			epPtr.GetOwningElement().removeOwnedBaseElement(epPtr)
		}
		epPtr.owningElement = element
		if epPtr.GetOwningElement() != nil {
			epPtr.GetOwningElement().addOwnedBaseElement(epPtr)
		}
		if epPtr.getElementPointerRole() == OWNING_ELEMENT {
			if epPtr.element != nil && epPtr.GetOwningElement() != nil {
				epPtr.element.addOwnedBaseElement(epPtr.GetOwningElement())
			}
		}
	}
}

type ElementPointer interface {
	Pointer
	GetElement() Element
	GetElementIdentifier() uuid.UUID
	GetElementVersion() int
	SetElement(Element)
}
