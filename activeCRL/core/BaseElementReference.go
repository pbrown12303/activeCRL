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

func (erPtr *baseElementReference) GetBaseElement(hl *HeldLocks) BaseElement {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(erPtr)
	rep := erPtr.getBaseElementPointer(hl)
	if rep != nil {
		return rep.GetBaseElement(hl)
	}
	return nil
}

func (erPtr *baseElementReference) getBaseElementPointer(hl *HeldLocks) BaseElementPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(erPtr)
	for _, be := range erPtr.ownedBaseElements {
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

func (bePtr *baseElementReference) isEquivalent(be *baseElementReference, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	var referencePtr *reference = &bePtr.reference
	return referencePtr.isEquivalent(&be.reference, hl)
}

func (elPtr *baseElementReference) MarshalJSON() ([]byte, error) {
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

func (elPtr *baseElementReference) printBaseElementReference(prefix string, hl *HeldLocks) {
	elPtr.printReference(prefix, hl)
}

func (el *baseElementReference) recoverBaseElementReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.reference.recoverReferenceFields(unmarshaledData)
}

//func (erPtr *baseElementReference) SetOwningElement(parent Element, hl *HeldLocks) {
//	if hl == nil {
//		hl = NewHeldLocks()
//		defer hl.ReleaseLocks()
//	}
//	hl.LockBaseElement(erPtr)
//	oldParent := erPtr.GetOwningElement(hl)
//	if oldParent == nil && parent == nil {
//		return // Nothing to do
//	} else if oldParent != nil && parent != nil && oldParent.GetId(hl) != parent.GetId(hl) {
//		return // Nothing to do
//	}
//	if oldParent != nil {
//		hl.LockBaseElement(oldParent)
//	}
//	if parent != nil {
//		hl.LockBaseElement(parent)
//	}
//	oep := erPtr.getOwningElementPointer(hl)
//	if oep != nil {
//		hl.LockBaseElement(oep)
//	} else {
//		oep = erPtr.uOfD.NewOwningElementPointer(hl)
//		oep.SetOwningElement(erPtr, hl)
//	}
//	oep.SetElement(parent, hl)
//}

func (erPtr *baseElementReference) SetBaseElement(el BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(erPtr)
	if erPtr.GetBaseElement(hl) != el {
		ep := erPtr.getBaseElementPointer(hl)
		if ep == nil {
			ep = erPtr.uOfD.NewBaseElementPointer(hl)
			SetOwningElement(ep, erPtr, hl)
		}
		ep.SetBaseElement(el, hl)
	}
}

type BaseElementReference interface {
	Reference
	GetBaseElement(*HeldLocks) BaseElement
	getBaseElementPointer(*HeldLocks) BaseElementPointer
	SetBaseElement(BaseElement, *HeldLocks)
}
