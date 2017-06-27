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

func (eprPtr *elementPointerReference) clone() *elementPointerReference {
	var clone elementPointerReference
	clone.ownedBaseElements = make(map[string]BaseElement)
	clone.cloneAttributes(*eprPtr)
	return &clone
}

func (eprPtr *elementPointerReference) cloneAttributes(source elementPointerReference) {
	eprPtr.reference.cloneAttributes(source.reference)
}

func (eprPtr *elementPointerReference) GetElementPointer() ElementPointer {
	eprPtr.traceableLock()
	defer eprPtr.traceableUnlock()
	rep := eprPtr.getElementPointerPointer()
	if rep != nil {
		rep.traceableLock()
		defer rep.traceableUnlock()
	}
	return eprPtr.getElementPointer()
}

func (eprPtr *elementPointerReference) getElementPointer() ElementPointer {
	rep := eprPtr.getElementPointerPointer()
	if rep != nil {
		return rep.getElementPointer()
	}
	return nil
}

func (eprPtr *elementPointerReference) getElementPointerPointer() ElementPointerPointer {
	for _, be := range eprPtr.getOwnedBaseElements() {
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
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
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

func (erPtr *elementPointerReference) SetOwningElement(parent Element) {
	erPtr.traceableLock()
	defer erPtr.traceableUnlock()
	oldParent := erPtr.getOwningElement()
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
	oep := erPtr.getOwningElementPointer()
	if oep != nil {
		oep.traceableLock()
		defer oep.traceableUnlock()
	}
	erPtr.setOwningElement(parent)
}

func (erPtr *elementPointerReference) setOwningElement(owningElement Element) {
	if erPtr.getOwningElement() != owningElement {
		oep := erPtr.getOwningElementPointer()
		if oep == nil {
			oep = erPtr.uOfD.NewOwningElementPointer()
			oep.setOwningElement(erPtr)
		}
		oep.setElement(owningElement)
	}
}

func (eprPtr *elementPointerReference) SetElementPointer(el ElementPointer) {
	eprPtr.traceableLock()
	defer eprPtr.traceableUnlock()
	ep := eprPtr.getElementPointerPointer()
	if ep != nil {
		ep.traceableLock()
		defer ep.traceableUnlock()
	}
	eprPtr.setElementPointer(el)
}

func (eprPtr *elementPointerReference) setElementPointer(el ElementPointer) {
	if eprPtr.getElementPointer() != el {
		ep := eprPtr.getElementPointerPointer()
		if ep == nil {
			ep = eprPtr.uOfD.NewElementPointerPointer()
			ep.setOwningElement(eprPtr)
		}
		ep.setElementPointer(el)
	}
}

type ElementPointerReference interface {
	Reference
	GetElementPointer() ElementPointer
	getElementPointerPointer() ElementPointerPointer
	setElementPointer(ElementPointer)
	SetElementPointer(ElementPointer)
}
