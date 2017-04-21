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
	uOfD.AddBaseElement(&el)
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
	elPtr.Lock()
	defer elPtr.Unlock()
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

func (erPtr *elementPointerReference) setOwningElement(owningElement Element) {
	oep := erPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(erPtr.uOfD)
		oep.setOwningElement(erPtr)
	}
	oep.setElement(owningElement)
}

func (eprPtr *elementPointerReference) SetElementPointer(el ElementPointer) {
	eprPtr.Lock()
	defer eprPtr.Unlock()
	ep := eprPtr.getElementPointerPointer()
	if ep != nil {
		ep.Lock()
		defer ep.Unlock()
	}
	eprPtr.setElementPointer(el)
}

func (eprPtr *elementPointerReference) setElementPointer(el ElementPointer) {
	ep := eprPtr.getElementPointerPointer()
	if ep == nil {
		ep = NewElementPointerPointer(eprPtr.uOfD)
		ep.setOwningElement(eprPtr)
	}
	ep.setElementPointer(el)
}

type ElementPointerReference interface {
	Reference
	GetElementPointer() ElementPointer
	setElementPointer(ElementPointer)
	SetElementPointer(ElementPointer)
}
