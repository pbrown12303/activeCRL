package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type elementPointerReference struct {
	reference
}

func NewElementPointerReference(uOfD *UniverseOfDiscourse) ElementPointerReference {
	var el elementPointerReference
	el.initializeElementPointerReference()
	uOfD.addBaseElement(&el)
	return &el
}

func (eprPtr *elementPointerReference) GetElementPointer() ElementPointer {
	rep := eprPtr.getElementPointerPointer()
	if rep != nil {
		return rep.GetElementPointer()
	}
	return nil
}

func (eprPtr *elementPointerReference) getElementPointerPointer() ElementPointerPointer {
	for _, be := range eprPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *elementPointerPointer:
			return be.(ElementPointerPointer)
		}
	}
	return nil
}

func (elPtr *elementPointerReference) initializeElementPointerReference() {
	elPtr.initializeReference()
}

func (bePtr *elementPointerReference) isEquivalent(be *elementPointerReference) bool {
	var referencePtr *reference = &bePtr.reference
	return referencePtr.isEquivalent(&be.reference)
}

func (elPtr *elementPointerReference) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalElementPointerReferenceFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *elementPointerReference) marshalElementPointerReferenceFields(buffer *bytes.Buffer) error {
	return elPtr.reference.marshalReferenceFields(buffer)
}

func (elPtr *elementPointerReference) printElementPointerReference(prefix string) {
	elPtr.printReference(prefix)
}

func (el *elementPointerReference) recoverElementPointerReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.reference.recoverReferenceFields(unmarshaledData)
}

func (elPtr *elementPointerReference) setOwningElement(owningElement Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(elPtr.uOfD)
		oep.setOwningElement(elPtr)
	}
	oep.SetElement(owningElement)
}

func (eprPtr *elementPointerReference) SetElementPointer(el ElementPointer) {
	ep := eprPtr.getElementPointerPointer()
	if ep == nil {
		ep = NewElementPointerPointer(eprPtr.uOfD)
		ep.setOwningElement(eprPtr)
	}
	ep.SetElementPointer(el)
}

type ElementPointerReference interface {
	Reference
	GetElementPointer() ElementPointer
	SetElementPointer(ElementPointer)
}
