package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/satori/go.uuid"
)

type elementPointerPointer struct {
	pointer
	elementPointer        ElementPointer
	elementPointerId      uuid.UUID
	elementPointerVersion int
}

func NewElementPointerPointer(uOfD *UniverseOfDiscourse) ElementPointerPointer {
	var ep elementPointerPointer
	ep.initializeElementPointerPointer()
	uOfD.AddBaseElement(&ep)
	return &ep
}

func (eppPtr *elementPointerPointer) clone() *elementPointerPointer {
	var clone elementPointerPointer
	clone.cloneAttributes(*eppPtr)
	return &clone
}

func (eppPtr *elementPointerPointer) cloneAttributes(source elementPointerPointer) {
	eppPtr.pointer.cloneAttributes(source.pointer)
	eppPtr.elementPointer = source.elementPointer
	eppPtr.elementPointerId = source.elementPointerId
	eppPtr.elementPointerVersion = source.elementPointerVersion
}

func (eppPtr *elementPointerPointer) GetElementPointer() ElementPointer {
	eppPtr.traceableLock()
	defer eppPtr.traceableUnlock()
	if eppPtr.elementPointer == nil && eppPtr.getElementPointerIdentifier() != uuid.Nil && eppPtr.uOfD != nil {
		eppPtr.elementPointer = eppPtr.uOfD.getElementPointer(eppPtr.getElementPointerIdentifier().String())
	}
	return eppPtr.elementPointer
}

func (eppPtr *elementPointerPointer) getElementPointer() ElementPointer {
	if eppPtr.elementPointer == nil && eppPtr.getElementPointerIdentifier() != uuid.Nil && eppPtr.uOfD != nil {
		eppPtr.elementPointer = eppPtr.uOfD.getElementPointer(eppPtr.getElementPointerIdentifier().String())
	}
	return eppPtr.elementPointer
}

func (eppPtr *elementPointerPointer) GetName() string {
	return "elementPointerPointer"
}

func (eppPtr *elementPointerPointer) GetElementPointerIdentifier() uuid.UUID {
	eppPtr.traceableLock()
	defer eppPtr.traceableUnlock()
	return eppPtr.getElementPointerIdentifier()
}

func (eppPtr *elementPointerPointer) getElementPointerIdentifier() uuid.UUID {
	return eppPtr.elementPointerId
}

func (eppPtr *elementPointerPointer) GetElementPointerVersion() int {
	eppPtr.traceableLock()
	defer eppPtr.traceableUnlock()
	return eppPtr.getElementPointerVersion()
}

func (eppPtr *elementPointerPointer) getElementPointerVersion() int {
	return eppPtr.elementPointerVersion
}

func (eppPtr *elementPointerPointer) initializeElementPointerPointer() {
	eppPtr.initializePointer()
}

func (bePtr *elementPointerPointer) isEquivalent(be *elementPointerPointer) bool {
	if bePtr.elementPointerId != be.elementPointerId {
		fmt.Printf("Equivalence failed: indicated elementPointerPointer ids do not match \n")
		return false
	}
	if bePtr.elementPointerVersion != be.elementPointerVersion {
		fmt.Printf("Equivalence failed: indicated elementPointerPointer versions do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer)
}

func (elPtr *elementPointerPointer) MarshalJSON() ([]byte, error) {
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.maarshalElementPointerPointerFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *elementPointerPointer) maarshalElementPointerPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.pointer.marshalPointerFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"ElementPointerId\":\"%s\",", elPtr.elementPointerId.String()))
	buffer.WriteString(fmt.Sprintf("\"ElementPointerVersion\":\"%d\"", elPtr.elementPointerVersion))
	return err
}

func (eppPtr *elementPointerPointer) printElementPointerPointer(prefix string) {
	eppPtr.printPointer(prefix)
	fmt.Printf("%sIndicated ElementPointerID: %s \n", prefix, eppPtr.elementPointerId.String())
	fmt.Printf("%sIndicated ElementPointerVersion: %d \n", prefix, eppPtr.elementPointerVersion)
}

func (ep *elementPointerPointer) recoverElementPointerPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	err := ep.pointer.recoverPointerFields(unmarshaledData)
	if err != nil {
		fmt.Printf("ElementPointerPointer's Recovery of PointerFields failed\n")
		return err
	}
	// ElementPointer ID
	var recoveredElementPointerId string
	err = json.Unmarshal((*unmarshaledData)["ElementPointerId"], &recoveredElementPointerId)
	if err != nil {
		fmt.Printf("ElementPointerPointer's Recovery of ElementPointerId failed\n")
		return err
	}
	ep.elementPointerId, err = uuid.FromString(recoveredElementPointerId)
	if err != nil {
		fmt.Printf("ElementPointerPointer's conversion of ElementPointerId failed\n")
		return err
	}
	// Version
	var recoveredElementPointerVersion string
	err = json.Unmarshal((*unmarshaledData)["ElementPointerVersion"], &recoveredElementPointerVersion)
	if err != nil {
		fmt.Printf("ElementPointerPointer's Recovery of ElementPointerVersion failed\n")
		return err
	}
	ep.elementPointerVersion, err = strconv.Atoi(recoveredElementPointerVersion)
	if err != nil {
		fmt.Printf("Conversion of ElementPointerPointer.elementPointerVersion failed\n")
		return err
	}
	return nil
}

func (eppPtr *elementPointerPointer) SetElementPointer(elementPointer ElementPointer) {
	eppPtr.traceableLock()
	defer eppPtr.traceableUnlock()
	if elementPointer != nil {
		elementPointer.traceableLock()
		defer elementPointer.traceableUnlock()
	}
	eppPtr.setElementPointer(elementPointer)
}

func (eppPtr *elementPointerPointer) setElementPointer(elementPointer ElementPointer) {
	if elementPointer != eppPtr.elementPointer {
		preChange(eppPtr)
		eppPtr.elementPointer = elementPointer
		if elementPointer != nil {
			eppPtr.elementPointerId = elementPointer.getId()
			eppPtr.elementPointerVersion = elementPointer.getVersion()
		} else {
			eppPtr.elementPointerId = uuid.Nil
			eppPtr.elementPointerVersion = 0
		}
		postChange(eppPtr)
	}
}

func (eppPtr *elementPointerPointer) SetOwningElement(element Element) {
	eppPtr.traceableLock()
	defer eppPtr.traceableUnlock()
	currentOwner := eppPtr.getOwningElement()
	if currentOwner != element {
		if eppPtr.getOwningElement() != nil {
			currentOwner.traceableLock()
			defer currentOwner.traceableUnlock()
		}
		if element != nil {
			element.traceableLock()
			defer element.traceableUnlock()
		}
		eppPtr.setOwningElement(element)
	}
}

func (eppPtr *elementPointerPointer) setOwningElement(element Element) {
	if element != eppPtr.getOwningElement() {
		if eppPtr.getOwningElement() != nil {
			eppPtr.getOwningElement().removeOwnedBaseElement(eppPtr)
		}

		preChange(eppPtr)
		eppPtr.owningElement = element
		postChange(eppPtr)

		if eppPtr.getOwningElement() != nil {
			eppPtr.getOwningElement().addOwnedBaseElement(eppPtr)
		}
	}
}

// internalSetOwningElement() is an internal function used only in unmarshal
func (eppPtr *elementPointerPointer) internalSetOwningElement(element Element) {
	if element != eppPtr.GetOwningElement() {
		eppPtr.owningElement = element
		if eppPtr.GetOwningElement() != nil {
			eppPtr.GetOwningElement().internalAddOwnedBaseElement(eppPtr)
		}
	}
}

func (epPtr *elementPointerPointer) SetUri(uri string) {
	epPtr.traceableLock()
	defer epPtr.traceableUnlock()
	epPtr.setUri(uri)
}

func (epPtr *elementPointerPointer) setUri(uri string) {
	preChange(epPtr)
	epPtr.uri = uri
	postChange(epPtr)
}

type ElementPointerPointer interface {
	Pointer
	GetElementPointer() ElementPointer
	getElementPointer() ElementPointer
	GetElementPointerIdentifier() uuid.UUID
	GetElementPointerVersion() int
	setElementPointer(ElementPointer)
	SetElementPointer(ElementPointer)
}
