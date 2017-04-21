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
	uOfD.AddBaseElement(&el)
	return &el
}

func (rPtr *refinement) GetAbstractElement() Element {
	rPtr.Lock()
	defer rPtr.Unlock()
	return rPtr.getAbstractElement()
}

func (rPtr *refinement) getAbstractElement() Element {
	rep := rPtr.getAbstractElementPointer()
	if rep != nil {
		return rep.GetElement()
	}
	return nil
}

func (rPtr *refinement) getAbstractElementPointer() ElementPointer {
	for _, be := range rPtr.getOwnedBaseElements() {
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
	rPtr.Lock()
	defer rPtr.Unlock()
	return rPtr.getRefinedElement()
}

func (rPtr *refinement) getRefinedElement() Element {
	rep := rPtr.getRefinedElementPointer()
	if rep != nil {
		return rep.GetElement()
	}
	return nil
}

func (rPtr *refinement) getRefinedElementPointer() ElementPointer {
	for _, be := range rPtr.getOwnedBaseElements() {
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
	elPtr.Lock()
	defer elPtr.Unlock()
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
	rPtr.Lock()
	defer rPtr.Unlock()
	ep := rPtr.getAbstractElementPointer()
	if ep != nil {
		ep.Lock()
		defer ep.Unlock()
	}
	if el != nil {
		el.Lock()
		defer el.Unlock()
	}
	rPtr.setAbstractElement(el)
}

func (rPtr *refinement) setAbstractElement(el Element) {
	ep := rPtr.getAbstractElementPointer()
	if ep == nil {
		ep = NewAbstractElementPointer(rPtr.uOfD)
		ep.setOwningElement(rPtr)
	}
	ep.setElement(el)
}

func (elPtr *refinement) SetOwningElement(parent Element) {
	elPtr.Lock()
	defer elPtr.Unlock()
	oldParent := elPtr.getOwningElement()
	if oldParent == nil && parent == nil {
		return // Nothing to do
	} else if oldParent != nil && parent != nil && oldParent.getId() != parent.getId() {
		return // Nothing to do
	}
	if oldParent != nil {
		oldParent.Lock()
		defer oldParent.Unlock()
	}
	if parent != nil {
		parent.Lock()
		defer parent.Unlock()
	}
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		oep.Lock()
		defer oep.Unlock()
	}
	elPtr.setOwningElement(parent)
}

func (elPtr *refinement) setOwningElement(parent Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(elPtr.uOfD)
		oep.setOwningElement(elPtr)
	}
	oep.setElement(parent)
}

func (rPtr *refinement) SetRefinedElement(el Element) {
	rPtr.Lock()
	defer rPtr.Unlock()
	ep := rPtr.getRefinedElementPointer()
	if ep != nil {
		ep.Lock()
		defer ep.Unlock()
	}
	if el != nil {
		el.Lock()
		defer el.Unlock()
	}
	rPtr.setRefinedElement(el)

}

func (rPtr *refinement) setRefinedElement(el Element) {
	ep := rPtr.getRefinedElementPointer()
	if ep == nil {
		ep = NewRefinedElementPointer(rPtr.uOfD)
		ep.setOwningElement(rPtr)
	}
	ep.setElement(el)
}

type Refinement interface {
	Element
	GetAbstractElement() Element
	GetRefinedElement() Element
	SetAbstractElement(Element)
	SetRefinedElement(Element)
}
