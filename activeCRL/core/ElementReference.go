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

func (erPtr *elementReference) clone() *elementReference {
	var clone elementReference
	clone.ownedBaseElements = make(map[string]BaseElement)
	clone.cloneAttributes(*erPtr)
	return &clone
}

func (erPtr *elementReference) cloneAttributes(source elementReference) {
	erPtr.reference.cloneAttributes(source.reference)
}

func (erPtr *elementReference) GetReferencedElement() Element {
	erPtr.TraceableLock()
	defer erPtr.TraceableUnlock()
	return erPtr.getReferencedElement()
}

func (erPtr *elementReference) getReferencedElement() Element {
	rep := erPtr.getReferencedElementPointer()
	if rep != nil {
		return rep.getElement()
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
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
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
	erPtr.TraceableLock()
	defer erPtr.TraceableUnlock()
	oldParent := erPtr.getOwningElement()
	if oldParent == nil && parent == nil {
		return // Nothing to do
	} else if oldParent != nil && parent != nil && oldParent.getId() != parent.getId() {
		return // Nothing to do
	}
	if oldParent != nil {
		oldParent.TraceableLock()
		defer oldParent.TraceableUnlock()
	}
	if parent != nil {
		parent.TraceableLock()
		defer parent.TraceableUnlock()
	}
	oep := erPtr.getOwningElementPointer()
	if oep != nil {
		oep.TraceableLock()
		defer oep.TraceableUnlock()
	}
	erPtr.SetOwningElementNoLock(parent)
}

func (erPtr *elementReference) SetOwningElementNoLock(owningElement Element) {
	if erPtr.getOwningElement() != owningElement {
		oep := erPtr.getOwningElementPointer()
		if oep == nil {
			oep = erPtr.uOfD.NewOwningElementPointer()
			oep.SetOwningElementNoLock(erPtr)
		}
		oep.setElement(owningElement)
	}
}

func (erPtr *elementReference) SetReferencedElement(el Element) {
	erPtr.TraceableLock()
	defer erPtr.TraceableUnlock()
	ep := erPtr.getReferencedElementPointer()
	if ep != nil {
		ep.TraceableLock()
		defer ep.TraceableUnlock()
	}
	erPtr.setReferencedElement(el)
}

func (erPtr *elementReference) setReferencedElement(el Element) {
	if erPtr.getReferencedElement() != el {
		ep := erPtr.getReferencedElementPointer()
		if ep == nil {
			ep = erPtr.uOfD.NewReferencedElementPointer()
			ep.SetOwningElementNoLock(erPtr)
		}
		ep.setElement(el)
	}
}

type ElementReference interface {
	Reference
	GetReferencedElement() Element
	getReferencedElementPointer() ElementPointer
	SetReferencedElement(Element)
}
