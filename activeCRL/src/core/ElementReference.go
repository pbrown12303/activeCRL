package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type elementReference struct {
	reference
}

func NewElementReference(uOfD *UniverseOfDiscourse) ElementReference {
	var el elementReference
	el.initializeElementReference()
	uOfD.addBaseElement(&el)
	return &el
}

func (erPtr *elementReference) GetReferencedElement() Element {
	rep := erPtr.getReferencedElementPointer()
	if rep != nil {
		return rep.GetElement()
	}
	return nil
}

func (erPtr *elementReference) getReferencedElementPointer() ElementPointer {
	for _, be := range erPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *elementPointer:
			if be.(*elementPointer).getElementPointerRole() == REFERENCED_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

func (elPtr *elementReference) initializeElementReference() {
	elPtr.initializeReference()
}

func (bePtr *elementReference) isEquivalent(be *elementReference) bool {
	var referencePtr *reference = &bePtr.reference
	return referencePtr.isEquivalent(&be.reference)
}

func (elPtr *elementReference) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalElementReferenceFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *elementReference) marshalElementReferenceFields(buffer *bytes.Buffer) error {
	return elPtr.reference.marshalReferenceFields(buffer)
}

func (elPtr *elementReference) printElementReference(prefix string) {
	elPtr.printReference(prefix)
}

func (el *elementReference) recoverElementReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.reference.recoverReferenceFields(unmarshaledData)
}

func (elPtr *elementReference) setOwningElement(owningElement Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(elPtr.uOfD)
		oep.setOwningElement(elPtr)
	}
	oep.SetElement(owningElement)
}

func (erPtr *elementReference) SetReferencedElement(el Element) {
	ep := erPtr.getReferencedElementPointer()
	if ep == nil {
		ep = NewReferencedElementPointer(erPtr.uOfD)
		ep.setOwningElement(erPtr)
	}
	ep.SetElement(el)
}

type ElementReference interface {
	Reference
	GetReferencedElement() Element
	SetReferencedElement(Element)
}
