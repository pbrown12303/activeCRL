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

func NewLiteralPointerReference(uOfD *UniverseOfDiscourse) LiteralPointerReference {
	var el literalPointerReference
	el.initializeLiteralPointerReference()
	uOfD.addBaseElement(&el)
	return &el
}

func (lprPtr *literalPointerReference) GetLiteralPointer() LiteralPointer {
	rep := lprPtr.getLiteralPointerPointer()
	if rep != nil {
		return rep.GetLiteralPointer()
	}
	return nil
}

func (lprPtr *literalPointerReference) getLiteralPointerPointer() LiteralPointerPointer {
	for _, be := range lprPtr.GetOwnedBaseElements() {
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

func (elPtr *literalPointerReference) setOwningElement(owningElement Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(elPtr.uOfD)
		oep.setOwningElement(elPtr)
	}
	oep.SetElement(owningElement)
}

func (lprPtr *literalPointerReference) SetLiteralPointer(el LiteralPointer) {
	ep := lprPtr.getLiteralPointerPointer()
	if ep == nil {
		ep = NewLiteralPointerPointer(lprPtr.uOfD)
		ep.setOwningElement(lprPtr)
	}
	ep.SetLiteralPointer(el)
}

type LiteralPointerReference interface {
	Reference
	GetLiteralPointer() LiteralPointer
	SetLiteralPointer(LiteralPointer)
}
