package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type refinement struct {
	element
}

func NewRefinement(uOfD *UniverseOfDiscourse) Refinement {
	var el refinement
	el.initializeRefinement()
	uOfD.addBaseElement(&el)
	return &el
}

func (rPtr *refinement) GetAbstractElement() Element {
	rep := rPtr.getAbstractElementPointer()
	if rep != nil {
		return rep.GetElement()
	}
	return nil
}

func (rPtr *refinement) getAbstractElementPointer() ElementPointer {
	for _, be := range rPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *elementPointer:
			if be.(*elementPointer).getElementPointerRole() == ABSTRACT_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

func (rPtr *refinement) GetRefinedElement() Element {
	rep := rPtr.getRefinedElementPointer()
	if rep != nil {
		return rep.GetElement()
	}
	return nil
}

func (rPtr *refinement) getRefinedElementPointer() ElementPointer {
	for _, be := range rPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *elementPointer:
			if be.(*elementPointer).getElementPointerRole() == REFINED_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

func (rPtr *refinement) initializeRefinement() {
	rPtr.initializeElement()
}

func (bePtr *refinement) isEquivalent(be *refinement) bool {
	var elementPtr *element = &bePtr.element
	return elementPtr.isEquivalent(&be.element)
}

func (elPtr *refinement) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalRefinementFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *refinement) marshalRefinementFields(buffer *bytes.Buffer) error {
	return elPtr.element.marshalElementFields(buffer)
}

func (elPtr *refinement) printRefinement(prefix string) {
	elPtr.printElement(prefix)
}

func (el *refinement) recoverRefinementFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.element.recoverElementFields(unmarshaledData)
}

func (rPtr *refinement) SetAbstractElement(el Element) {
	ep := rPtr.getAbstractElementPointer()
	if ep == nil {
		ep = NewAbstractElementPointer(rPtr.uOfD)
		ep.setOwningElement(rPtr)
	}
	ep.SetElement(el)
}

func (elPtr *refinement) setOwningElement(owningElement Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(elPtr.uOfD)
		oep.setOwningElement(elPtr)
	}
	oep.SetElement(owningElement)
}

func (rPtr *refinement) SetRefinedElement(el Element) {
	ep := rPtr.getRefinedElementPointer()
	if ep == nil {
		ep = NewRefinedElementPointer(rPtr.uOfD)
		ep.setOwningElement(rPtr)
	}
	ep.SetElement(el)

}

type Refinement interface {
	Element
	GetAbstractElement() Element
	GetRefinedElement() Element
	SetAbstractElement(Element)
	SetRefinedElement(Element)
}
