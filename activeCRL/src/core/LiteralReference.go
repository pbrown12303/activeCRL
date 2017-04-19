package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type literalReference struct {
	reference
}

func NewLiteralReference(uOfD *UniverseOfDiscourse) LiteralReference {
	var el literalReference
	el.initializeLiteralReference()
	uOfD.addBaseElement(&el)
	return &el
}

func (lrPtr *literalReference) GetReferencedLiteral() Literal {
	rep := lrPtr.getReferencedLiteralPointer()
	if rep != nil {
		return rep.GetLiteral()
	}
	return nil
}

func (lrPtr *literalReference) getReferencedLiteralPointer() LiteralPointer {
	for _, be := range lrPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *literalPointer:
			if be.(*literalPointer).getLiteralPointerRole() == VALUE {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

func (elPtr *literalReference) initializeLiteralReference() {
	elPtr.initializeReference()
}

func (bePtr *literalReference) isEquivalent(be *literalReference) bool {
	var elementPtr *element = &bePtr.element
	return elementPtr.isEquivalent(&be.element)
}

func (elPtr *literalReference) MarshalJSON() ([]byte, error) {
	//	fmt.Printf("MarshalJSON called on Literal Reference \n")
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalLiteralReferenceFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *literalReference) marshalLiteralReferenceFields(buffer *bytes.Buffer) error {
	return elPtr.element.marshalElementFields(buffer)
}

func (elPtr *literalReference) printLiteralReference(prefix string) {
	elPtr.printElement(prefix)
}

func (el *literalReference) recoverLiteralReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.element.recoverElementFields(unmarshaledData)
}

func (elPtr *literalReference) setOwningElement(owningElement Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(elPtr.uOfD)
		oep.setOwningElement(elPtr)
	}
	oep.SetElement(owningElement)
}

func (lrPtr *literalReference) SetReferencedLiteral(el Literal) {
	ep := lrPtr.getReferencedLiteralPointer()
	if ep == nil {
		ep = NewValueLiteralPointer(lrPtr.uOfD)
		ep.setOwningElement(lrPtr)
	}
	ep.SetLiteral(el)
}

type LiteralReference interface {
	Reference
	GetReferencedLiteral() Literal
	SetReferencedLiteral(Literal)
}
