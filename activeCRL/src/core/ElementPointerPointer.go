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

func (eppPtr *elementPointerPointer) GetElementPointer() ElementPointer {
	eppPtr.Lock()
	defer eppPtr.Unlock()
	if eppPtr.elementPointer == nil && eppPtr.getElementPointerIdentifier() != uuid.Nil && eppPtr.uOfD != nil {
		eppPtr.elementPointer = eppPtr.uOfD.getElementPointer(eppPtr.getElementPointerIdentifier().String())
	}
	return eppPtr.elementPointer
}

func (eppPtr *elementPointerPointer) GetName() string {
	return "elementPointerPointer"
}

func (eppPtr *elementPointerPointer) GetElementPointerIdentifier() uuid.UUID {
	eppPtr.Lock()
	defer eppPtr.Unlock()
	return eppPtr.getElementPointerIdentifier()
}

func (eppPtr *elementPointerPointer) getElementPointerIdentifier() uuid.UUID {
	return eppPtr.elementPointerId
}

func (eppPtr *elementPointerPointer) GetElementPointerVersion() int {
	eppPtr.Lock()
	defer eppPtr.Unlock()
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
	elPtr.Lock()
	defer elPtr.Unlock()
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
	eppPtr.Lock()
	defer eppPtr.Unlock()
	if elementPointer != nil {
		elementPointer.Lock()
		defer elementPointer.Unlock()
	}
	eppPtr.setElementPointer(elementPointer)
}

func (eppPtr *elementPointerPointer) setElementPointer(elementPointer ElementPointer) {
	if elementPointer != eppPtr.elementPointer {
		eppPtr.elementPointer = elementPointer
		if elementPointer != nil {
			eppPtr.elementPointerId = elementPointer.getId()
			eppPtr.elementPointerVersion = elementPointer.getVersion()
		} else {
			eppPtr.elementPointerId = uuid.Nil
			eppPtr.elementPointerVersion = 0
		}
	}
}

func (eppPtr *elementPointerPointer) SetOwningElement(element Element) {
	eppPtr.Lock()
	defer eppPtr.Unlock()
	eppPtr.setOwningElement(element)
}

func (eppPtr *elementPointerPointer) setOwningElement(element Element) {
	if element != eppPtr.GetOwningElement() {
		if eppPtr.GetOwningElement() != nil {
			eppPtr.GetOwningElement().removeOwnedBaseElement(eppPtr)
		}
		eppPtr.owningElement = element
		if eppPtr.GetOwningElement() != nil {
			eppPtr.GetOwningElement().addOwnedBaseElement(eppPtr)
		}
	}
}

type ElementPointerPointer interface {
	Pointer
	GetElementPointer() ElementPointer
	GetElementPointerIdentifier() uuid.UUID
	GetElementPointerVersion() int
	setElementPointer(ElementPointer)
	SetElementPointer(ElementPointer)
}
