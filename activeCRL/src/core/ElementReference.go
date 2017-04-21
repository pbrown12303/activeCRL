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
	uOfD.AddBaseElement(&el)
	return &el
}

func (erPtr *elementReference) GetReferencedElement() Element {
	erPtr.Lock()
	defer erPtr.Unlock()
	return erPtr.getReferencedElement()
}

func (erPtr *elementReference) getReferencedElement() Element {
	rep := erPtr.getReferencedElementPointer()
	if rep != nil {
		return rep.GetElement()
	}
	return nil
}

func (erPtr *elementReference) getReferencedElementPointer() ElementPointer {
	for _, be := range erPtr.getOwnedBaseElements() {
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
	elPtr.Lock()
	defer elPtr.Unlock()
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

func (erPtr *elementReference) SetOwningElement(parent Element) {
	erPtr.Lock()
	defer erPtr.Unlock()
	oldParent := erPtr.getOwningElement()
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
	oep := erPtr.getOwningElementPointer()
	if oep != nil {
		oep.Lock()
		defer oep.Unlock()
	}
	erPtr.setOwningElement(parent)
}

func (erPtr *elementReference) setOwningElement(owningElement Element) {
	oep := erPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(erPtr.uOfD)
		oep.setOwningElement(erPtr)
	}
	oep.setElement(owningElement)
}

func (erPtr *elementReference) SetReferencedElement(el Element) {
	erPtr.Lock()
	defer erPtr.Unlock()
	ep := erPtr.getReferencedElementPointer()
	if ep != nil {
		ep.Lock()
		defer ep.Unlock()
	}
	erPtr.setReferencedElement(el)
}

func (erPtr *elementReference) setReferencedElement(el Element) {
	ep := erPtr.getReferencedElementPointer()
	if ep == nil {
		ep = NewReferencedElementPointer(erPtr.uOfD)
		ep.setOwningElement(erPtr)
	}
	ep.setElement(el)
}

type ElementReference interface {
	Reference
	GetReferencedElement() Element
	SetReferencedElement(Element)
}
