package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type baseElementReference struct {
	reference
}

func (erPtr *baseElementReference) clone() *baseElementReference {
	var clone baseElementReference
	clone.ownedBaseElements = make(map[string]BaseElement)
	clone.cloneAttributes(*erPtr)
	return &clone
}

func (erPtr *baseElementReference) cloneAttributes(source baseElementReference) {
	erPtr.reference.cloneAttributes(source.reference)
}

func (erPtr *baseElementReference) GetBaseElement() BaseElement {
	erPtr.TraceableLock()
	defer erPtr.TraceableUnlock()
	return erPtr.getBaseElement()
}

func (erPtr *baseElementReference) getBaseElement() BaseElement {
	rep := erPtr.getBaseElementPointer()
	if rep != nil {
		return rep.getBaseElement()
	}
	return nil
}

func (erPtr *baseElementReference) getBaseElementPointer() BaseElementPointer {
	for _, be := range erPtr.getOwnedBaseElements() {
		switch be.(type) {
		case *baseElementPointer:
			return be.(BaseElementPointer)
		}
	}
	return nil
}

func (elPtr *baseElementReference) initializeBaseElementReference() {
	elPtr.initializeReference()
}

func (bePtr *baseElementReference) isEquivalent(be *baseElementReference) bool {
	var referencePtr *reference = &bePtr.reference
	return referencePtr.isEquivalent(&be.reference)
}

func (elPtr *baseElementReference) MarshalJSON() ([]byte, error) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalBaseElementReferenceFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *baseElementReference) marshalBaseElementReferenceFields(buffer *bytes.Buffer) error {
	return elPtr.reference.marshalReferenceFields(buffer)
}

func (elPtr *baseElementReference) printBaseElementReference(prefix string) {
	elPtr.printReference(prefix)
}

func (el *baseElementReference) recoverBaseElementReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.reference.recoverReferenceFields(unmarshaledData)
}

func (erPtr *baseElementReference) SetOwningElement(parent Element) {
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

func (erPtr *baseElementReference) SetOwningElementNoLock(owningElement Element) {
	if erPtr.getOwningElement() != owningElement {
		oep := erPtr.getOwningElementPointer()
		if oep == nil {
			oep = erPtr.uOfD.NewOwningElementPointer()
			oep.SetOwningElementNoLock(erPtr)
		}
		oep.setElement(owningElement)
	}
}

func (erPtr *baseElementReference) SetBaseElement(el BaseElement) {
	erPtr.TraceableLock()
	defer erPtr.TraceableUnlock()
	ep := erPtr.getBaseElementPointer()
	if ep != nil {
		ep.TraceableLock()
		defer ep.TraceableUnlock()
	}
	erPtr.setBaseElement(el)
}

func (erPtr *baseElementReference) setBaseElement(el BaseElement) {
	if erPtr.getBaseElement() != el {
		ep := erPtr.getBaseElementPointer()
		if ep == nil {
			ep = erPtr.uOfD.NewBaseElementPointer()
			ep.SetOwningElementNoLock(erPtr)
		}
		ep.setBaseElement(el)
	}
}

type BaseElementReference interface {
	Reference
	GetBaseElement() BaseElement
	getBaseElementPointer() BaseElementPointer
	SetBaseElement(BaseElement)
}
