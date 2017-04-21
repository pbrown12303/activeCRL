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
	uOfD.AddBaseElement(&el)
	return &el
}

func (lrPtr *literalReference) GetReferencedLiteral() Literal {
	lrPtr.Lock()
	defer lrPtr.Unlock()
	return lrPtr.getReferencedLiteral()
}

func (lrPtr *literalReference) getReferencedLiteral() Literal {
	rep := lrPtr.getReferencedLiteralPointer()
	if rep != nil {
		return rep.GetLiteral()
	}
	return nil
}

func (lrPtr *literalReference) getReferencedLiteralPointer() LiteralPointer {
	for _, be := range lrPtr.getOwnedBaseElements() {
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
	elPtr.Lock()
	defer elPtr.Unlock()
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

func (elPtr *literalReference) SetOwningElement(parent Element) {
	elPtr.Lock()
	defer elPtr.Unlock()
	oldParent := elPtr.getOwningElement()
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
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		oep.Lock()
		defer oep.Unlock()
	}
	elPtr.setOwningElement(parent)
}

func (elPtr *literalReference) setOwningElement(parent Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(elPtr.uOfD)
		oep.setOwningElement(elPtr)
	}
	oep.setElement(parent)
}

func (lrPtr *literalReference) SetReferencedLiteral(el Literal) {
	lrPtr.Lock()
	defer lrPtr.Unlock()
	ep := lrPtr.getReferencedLiteralPointer()
	if ep != nil {
		ep.Lock()
		defer ep.Unlock()
	}
	if el != nil {
		el.Lock()
		defer el.Unlock()
	}
	lrPtr.setReferencedLiteral(el)
}

func (lrPtr *literalReference) setReferencedLiteral(el Literal) {
	ep := lrPtr.getReferencedLiteralPointer()
	if ep == nil {
		ep = NewValueLiteralPointer(lrPtr.uOfD)
		ep.setOwningElement(lrPtr)
	}
	ep.setLiteral(el)
}

type LiteralReference interface {
	Reference
	GetReferencedLiteral() Literal
	SetReferencedLiteral(Literal)
}
