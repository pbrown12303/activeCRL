package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type literalPointerReference struct {
	reference
}

func (lprPtr *literalPointerReference) clone() *literalPointerReference {
	var clone literalPointerReference
	clone.ownedBaseElements = make(map[string]BaseElement)
	clone.cloneAttributes(*lprPtr)
	return &clone
}

func (lprPtr *literalPointerReference) cloneAttributes(source literalPointerReference) {
	lprPtr.reference.cloneAttributes(source.reference)
}

func (lprPtr *literalPointerReference) GetLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lprPtr)
	rep := lprPtr.GetLiteralPointerPointer(hl)
	if rep != nil {
		return rep.GetLiteralPointer(hl)
	}
	return nil
}

func (lprPtr *literalPointerReference) GetLiteralPointerPointer(hl *HeldLocks) LiteralPointerPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lprPtr)
	for _, be := range lprPtr.ownedBaseElements {
		switch be.(type) {
		case LiteralPointerPointer:
			return be.(LiteralPointerPointer)
		}
	}
	return nil
}

func (elPtr *literalPointerReference) initializeLiteralPointerReference() {
	elPtr.initializeReference()
}

func (bePtr *literalPointerReference) isEquivalent(be *literalPointerReference, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	var referencePtr *reference = &bePtr.reference
	return referencePtr.isEquivalent(&be.reference, hl)
}

func (elPtr *literalPointerReference) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalLiteralPointerReferenceFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *literalPointerReference) marshalLiteralPointerReferenceFields(buffer *bytes.Buffer) error {
	return elPtr.reference.marshalReferenceFields(buffer)
}

func (elPtr *literalPointerReference) printLiteralPointerReference(prefix string, hl *HeldLocks) {
	elPtr.printReference(prefix, hl)
}

func (el *literalPointerReference) recoverLiteralPointerReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.reference.recoverReferenceFields(unmarshaledData)
}

//func (elPtr *literalPointerReference) SetOwningElement(parent Element, hl *HeldLocks) {
//	if hl == nil {
//		hl = NewHeldLocks()
//		defer hl.ReleaseLocks()
//	}
//	hl.LockBaseElement(elPtr)
//	oldParent := elPtr.GetOwningElement(hl)
//	if oldParent == nil && parent == nil {
//		return // Nothing to do
//	} else if oldParent != nil && parent != nil && oldParent.GetId(hl) == parent.GetId(hl) {
//		return // Nothing to do
//	}
//	oep := elPtr.getOwningElementPointer(hl)
//	if oep == nil {
//		oep = elPtr.uOfD.NewOwningElementPointer(hl)
//		oep.SetOwningElement(elPtr, hl)
//	}
//	oep.SetElement(parent, hl)
//}

func (lprPtr *literalPointerReference) SetLiteralPointer(lp LiteralPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lprPtr)
	if lprPtr.GetLiteralPointer(hl) != lp {
		ep := lprPtr.GetLiteralPointerPointer(hl)
		if ep == nil {
			ep = lprPtr.uOfD.NewLiteralPointerPointer(hl)
			SetOwningElement(ep, lprPtr, hl)
		}
		ep.SetLiteralPointer(lp, hl)
	}
}

type LiteralPointerReference interface {
	Reference
	GetLiteralPointer(*HeldLocks) LiteralPointer
	GetLiteralPointerPointer(*HeldLocks) LiteralPointerPointer
	SetLiteralPointer(LiteralPointer, *HeldLocks)
}
