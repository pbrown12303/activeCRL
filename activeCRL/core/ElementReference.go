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

func (erPtr *elementReference) GetReferencedElement(hl *HeldLocks) Element {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(erPtr)
	rep := erPtr.GetElementPointer(hl)
	if rep != nil {
		return rep.GetElement(hl)
	}
	return nil
}

func (erPtr *elementReference) GetElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(erPtr)
	for _, be := range erPtr.ownedBaseElements {
		switch be.(type) {
		case ElementPointer:
			if be.(ElementPointer).GetElementPointerRole(hl) == REFERENCED_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

func (elPtr *elementReference) initializeElementReference() {
	elPtr.initializeReference()
}

func (bePtr *elementReference) isEquivalent(be *elementReference, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	var referencePtr *reference = &bePtr.reference
	return referencePtr.isEquivalent(&be.reference, hl)
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

func (elPtr *elementReference) printElementReference(prefix string, hl *HeldLocks) {
	elPtr.printReference(prefix, hl)
}

func (el *elementReference) recoverElementReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.reference.recoverReferenceFields(unmarshaledData)
}

func (erPtr *elementReference) SetReferencedElement(el Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(erPtr)
	if erPtr.GetReferencedElement(hl) != el {
		ep := erPtr.GetElementPointer(hl)
		if ep == nil {
			ep = erPtr.uOfD.NewReferencedElementPointer(hl)
			SetOwningElement(ep, erPtr, hl)
		}
		ep.SetElement(el, hl)
	}
}

type ElementReference interface {
	Reference
	GetReferencedElement(*HeldLocks) Element
	GetElementPointer(*HeldLocks) ElementPointer
	SetReferencedElement(Element, *HeldLocks)
}
