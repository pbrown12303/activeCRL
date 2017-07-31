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

func (eprPtr *elementPointerReference) GetElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eprPtr)
	rep := eprPtr.getElementPointerPointer(hl)
	if rep != nil {
		return rep.GetElementPointer(hl)
	}
	return nil
}

func (eprPtr *elementPointerReference) getElementPointerPointer(hl *HeldLocks) ElementPointerPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eprPtr)
	for _, be := range eprPtr.ownedBaseElements {
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

func (bePtr *elementPointerReference) isEquivalent(be *elementPointerReference, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	var referencePtr *reference = &bePtr.reference
	return referencePtr.isEquivalent(&be.reference, hl)
}

func (elPtr *elementPointerReference) MarshalJSON() ([]byte, error) {
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

func (elPtr *elementPointerReference) printElementPointerReference(prefix string, hl *HeldLocks) {
	elPtr.printReference(prefix, hl)
}

func (el *elementPointerReference) recoverElementPointerReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.reference.recoverReferenceFields(unmarshaledData)
}

//func (erPtr *elementPointerReference) SetOwningElement(parent Element, hl *HeldLocks) {
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
//	if oldParent != parent {
//		oep := erPtr.getOwningElementPointer(hl)
//		if oep == nil {
//			oep = erPtr.uOfD.NewOwningElementPointer(hl)
//			oep.SetOwningElement(erPtr, hl)
//		}
//		oep.SetElement(parent, hl)
//	}
//}

func (eprPtr *elementPointerReference) SetElementPointer(el ElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eprPtr)
	if eprPtr.GetElementPointer(hl) != el {
		ep := eprPtr.getElementPointerPointer(hl)
		if ep == nil {
			ep = eprPtr.uOfD.NewElementPointerPointer(hl)
			SetOwningElement(ep, eprPtr, hl)
		}
		ep.SetElementPointer(el, hl)
	}
}

type ElementPointerReference interface {
	Reference
	GetElementPointer(*HeldLocks) ElementPointer
	getElementPointerPointer(*HeldLocks) ElementPointerPointer
	SetElementPointer(ElementPointer, *HeldLocks)
}
