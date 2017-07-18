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

func (lprPtr *literalPointerReference) GetLiteralPointer() LiteralPointer {
	lprPtr.TraceableLock()
	defer lprPtr.TraceableUnlock()
	return lprPtr.getLiteralPointer()
}

func (lprPtr *literalPointerReference) getLiteralPointer() LiteralPointer {
	rep := lprPtr.getLiteralPointerPointer()
	if rep != nil {
		return rep.getLiteralPointer()
	}
	return nil
}

func (lprPtr *literalPointerReference) getLiteralPointerPointer() LiteralPointerPointer {
	for _, be := range lprPtr.getOwnedBaseElements() {
		switch be.(type) {
		case *literalPointerPointer:
			return be.(LiteralPointerPointer)
		}
	}
	return nil
}

func (elPtr *literalPointerReference) initializeLiteralPointerReference() {
	elPtr.initializeReference()
}

func (bePtr *literalPointerReference) isEquivalent(be *literalPointerReference) bool {
	var referencePtr *reference = &bePtr.reference
	return referencePtr.isEquivalent(&be.reference)
}

func (elPtr *literalPointerReference) MarshalJSON() ([]byte, error) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
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

func (elPtr *literalPointerReference) printLiteralPointerReference(prefix string) {
	elPtr.printReference(prefix)
}

func (el *literalPointerReference) recoverLiteralPointerReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.reference.recoverReferenceFields(unmarshaledData)
}

func (elPtr *literalPointerReference) SetOwningElement(parent Element) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	oldParent := elPtr.getOwningElement()
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
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		oep.TraceableLock()
		defer oep.TraceableUnlock()
	}
	elPtr.SetOwningElementNoLock(parent)
}

func (elPtr *literalPointerReference) SetOwningElementNoLock(parent Element) {
	if elPtr.getOwningElement() != parent {
		oep := elPtr.getOwningElementPointer()
		if oep == nil {
			oep = elPtr.uOfD.NewOwningElementPointer()
			oep.SetOwningElementNoLock(elPtr)
		}
		oep.setElement(parent)
	}
}

func (lprPtr *literalPointerReference) SetLiteralPointer(lp LiteralPointer) {
	lprPtr.TraceableLock()
	defer lprPtr.TraceableUnlock()
	ep := lprPtr.getLiteralPointerPointer()
	if ep != nil {
		ep.TraceableLock()
		defer ep.TraceableUnlock()
	}
	lprPtr.setLiteralPointer(lp)
}

func (lprPtr *literalPointerReference) setLiteralPointer(lp LiteralPointer) {
	if lprPtr.getLiteralPointer() != lp {
		ep := lprPtr.getLiteralPointerPointer()
		if ep == nil {
			ep = lprPtr.uOfD.NewLiteralPointerPointer()
			ep.SetOwningElementNoLock(lprPtr)
		}
		ep.setLiteralPointer(lp)
	}
}

type LiteralPointerReference interface {
	Reference
	GetLiteralPointer() LiteralPointer
	getLiteralPointerPointer() LiteralPointerPointer
	SetLiteralPointer(LiteralPointer)
}
