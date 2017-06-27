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

func (rPtr *refinement) clone() *refinement {
	var clone refinement
	clone.ownedBaseElements = make(map[string]BaseElement)
	clone.cloneAttributes(*rPtr)
	return &clone
}

func (rPtr *refinement) cloneAttributes(source refinement) {
	rPtr.element.cloneAttributes(source.element)
}

func (rPtr *refinement) GetAbstractElement() Element {
	rPtr.traceableLock()
	defer rPtr.traceableUnlock()
	return rPtr.getAbstractElement()
}

func (rPtr *refinement) getAbstractElement() Element {
	rep := rPtr.getAbstractElementPointer()
	if rep != nil {
		return rep.getElement()
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
	rPtr.traceableLock()
	defer rPtr.traceableUnlock()
	return rPtr.getRefinedElement()
}

func (rPtr *refinement) getRefinedElement() Element {
	rep := rPtr.getRefinedElementPointer()
	if rep != nil {
		return rep.getElement()
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
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
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
	rPtr.traceableLock()
	defer rPtr.traceableUnlock()
	ep := rPtr.getAbstractElementPointer()
	if ep != nil {
		ep.traceableLock()
		defer ep.traceableUnlock()
	}
	if el != nil {
		el.traceableLock()
		defer el.traceableUnlock()
	}
	rPtr.setAbstractElement(el)
}

func (rPtr *refinement) setAbstractElement(el Element) {
	if rPtr.getAbstractElement() != el {
		ep := rPtr.getAbstractElementPointer()
		if ep == nil {
			ep = rPtr.uOfD.NewAbstractElementPointer()
			ep.setOwningElement(rPtr)
		}
		ep.setElement(el)
	}
}

func (elPtr *refinement) SetOwningElement(parent Element) {
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
	oldParent := elPtr.getOwningElement()
	if oldParent == nil && parent == nil {
		return // Nothing to do
	} else if oldParent != nil && parent != nil && oldParent.getId() != parent.getId() {
		return // Nothing to do
	}
	if oldParent != nil {
		oldParent.traceableLock()
		defer oldParent.traceableUnlock()
	}
	if parent != nil {
		parent.traceableLock()
		defer parent.traceableUnlock()
	}
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		oep.traceableLock()
		defer oep.traceableUnlock()
	}
	elPtr.setOwningElement(parent)
}

func (elPtr *refinement) setOwningElement(parent Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = elPtr.uOfD.NewOwningElementPointer()
		oep.setOwningElement(elPtr)
	}
	oep.setElement(parent)
}

func (rPtr *refinement) SetRefinedElement(el Element) {
	rPtr.traceableLock()
	defer rPtr.traceableUnlock()
	ep := rPtr.getRefinedElementPointer()
	if ep != nil {
		ep.traceableLock()
		defer ep.traceableUnlock()
	}
	if el != nil {
		el.traceableLock()
		defer el.traceableUnlock()
	}
	rPtr.setRefinedElement(el)

}

func (rPtr *refinement) setRefinedElement(el Element) {
	if rPtr.getRefinedElement() != el {
		ep := rPtr.getRefinedElementPointer()
		if ep == nil {
			ep = rPtr.uOfD.NewRefinedElementPointer()
			ep.setOwningElement(rPtr)
		}
		ep.setElement(el)
	}
}

type Refinement interface {
	Element
	GetAbstractElement() Element
	getAbstractElementPointer() ElementPointer
	GetRefinedElement() Element
	getRefinedElementPointer() ElementPointer
	SetAbstractElement(Element)
	SetRefinedElement(Element)
}
