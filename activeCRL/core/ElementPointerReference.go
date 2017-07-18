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
	eprPtr.TraceableLock()
	defer eprPtr.TraceableUnlock()
	rep := eprPtr.getElementPointerPointer()
	if rep != nil {
		rep.TraceableLock()
		defer rep.TraceableUnlock()
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
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
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

func (erPtr *elementPointerReference) SetOwningElementNoLock(owningElement Element) {
	if erPtr.getOwningElement() != owningElement {
		oep := erPtr.getOwningElementPointer()
		if oep == nil {
			oep = erPtr.uOfD.NewOwningElementPointer()
			oep.SetOwningElementNoLock(erPtr)
		}
		oep.setElement(owningElement)
	}
}

func (eprPtr *elementPointerReference) SetElementPointer(el ElementPointer) {
	eprPtr.TraceableLock()
	defer eprPtr.TraceableUnlock()
	ep := eprPtr.getElementPointerPointer()
	if ep != nil {
		ep.TraceableLock()
		defer ep.TraceableUnlock()
	}
	eprPtr.setElementPointer(el)
}

func (eprPtr *elementPointerReference) setElementPointer(el ElementPointer) {
	if eprPtr.getElementPointer() != el {
		ep := eprPtr.getElementPointerPointer()
		if ep == nil {
			ep = eprPtr.uOfD.NewElementPointerPointer()
			ep.SetOwningElementNoLock(eprPtr)
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
