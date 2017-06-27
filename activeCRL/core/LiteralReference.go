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

func (lrPtr *literalReference) clone() *literalReference {
	var clone literalReference
	clone.ownedBaseElements = make(map[string]BaseElement)
	clone.cloneAttributes(*lrPtr)
	return &clone
}

func (lrPtr *literalReference) cloneAttributes(source literalReference) {
	lrPtr.reference.cloneAttributes(source.reference)
}

func (lrPtr *literalReference) GetReferencedLiteral() Literal {
	lrPtr.traceableLock()
	defer lrPtr.traceableUnlock()
	return lrPtr.getReferencedLiteral()
}

func (lrPtr *literalReference) getReferencedLiteral() Literal {
	rep := lrPtr.getReferencedLiteralPointer()
	if rep != nil {
		return rep.getLiteral()
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
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
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
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
	oldParent := elPtr.getOwningElement()
	if oldParent == nil && parent == nil {
		return // Nothing to do
	} else if oldParent != nil && parent != nil && oldParent.getId() != parent.getId() {
		return // Nothing to do
	}
	if oldParent != nil {
		oldParent.traceableLock()
		defer oldParent.traceableUnlock()
	}
	if parent != nil {
		parent.traceableLock()
		defer parent.traceableUnlock()
	}
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		oep.traceableLock()
		defer oep.traceableUnlock()
	}
	elPtr.setOwningElement(parent)
}

func (elPtr *literalReference) setOwningElement(parent Element) {
	if elPtr.getOwningElement() != parent {
		oep := elPtr.getOwningElementPointer()
		if oep == nil {
			oep = elPtr.uOfD.NewOwningElementPointer()
			oep.setOwningElement(elPtr)
		}
		oep.setElement(parent)
	}
}

func (lrPtr *literalReference) SetReferencedLiteral(el Literal) {
	lrPtr.traceableLock()
	defer lrPtr.traceableUnlock()
	ep := lrPtr.getReferencedLiteralPointer()
	if ep != nil {
		ep.traceableLock()
		defer ep.traceableUnlock()
	}
	if el != nil {
		el.traceableLock()
		defer el.traceableUnlock()
	}
	lrPtr.setReferencedLiteral(el)
}

func (lrPtr *literalReference) setReferencedLiteral(el Literal) {
	if lrPtr.getReferencedLiteral() != el {
		ep := lrPtr.getReferencedLiteralPointer()
		if ep == nil {
			ep = lrPtr.uOfD.NewValueLiteralPointer()
			ep.setOwningElement(lrPtr)
		}
		ep.setLiteral(el)
	}
}

type LiteralReference interface {
	Reference
	GetReferencedLiteral() Literal
	SetReferencedLiteral(Literal)
}
