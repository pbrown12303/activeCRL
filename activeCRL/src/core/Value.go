package core

import (
	"bytes"
	"encoding/json"
	"log"
)

type value struct {
	baseElement
	owningElement Element
}

func (vPtr *value) cloneAttributes(source value) {
	vPtr.baseElement.cloneAttributes(source.baseElement)
	vPtr.owningElement = source.owningElement
}

func (vPtr *value) GetOwningElement() Element {
	vPtr.traceableLock()
	defer vPtr.traceableUnlock()
	return vPtr.getOwningElement()
}

func (vPtr *value) getOwningElement() Element {
	return vPtr.owningElement
}

func (vPtr *value) initializeValue() {
	vPtr.initializeBaseElement()
}

func (vPtr *value) isEquivalent(be *value) bool {
	if (vPtr.owningElement == nil && be.owningElement != nil) || (vPtr.owningElement != nil && be.owningElement == nil) {
		log.Printf("Equivalence failed: Value's Owning Elements do not match - one is nil and the other is not \n")
		log.Printf("First value: %#v \n", vPtr)
		log.Printf("First value's owner: \n")
		Print(vPtr.owningElement, "   ")
		log.Printf("Second value: %#v \n", be)
		log.Printf("Second value's owner: \n")
		Print(be.owningElement, "   ")
		return false
	}
	if vPtr.owningElement != nil && be.owningElement != nil && vPtr.owningElement.getId() != be.owningElement.getId() {
		log.Printf("Equivalence failed: Value's Owning Elements do not match - they have different identifiers\n")
		log.Printf("First value's owner: \n")
		Print(vPtr.owningElement, "   ")
		log.Printf("Second value's owner: \n")
		Print(be.owningElement, "   ")
		return false
	}
	var baseElementPtr *baseElement = &vPtr.baseElement
	return baseElementPtr.isEquivalent(&be.baseElement)
}

func (vPtr *value) marshalValueFields(buffer *bytes.Buffer) error {
	vPtr.baseElement.marshalBaseElementFields(buffer)
	return nil
}

func (vPtr *value) printValue(prefix string) {
	vPtr.printBaseElement(prefix)
	if vPtr.getOwningElement() == nil {
		log.Printf("%sowningElmentIdentifier: %s \n", prefix, "nil")
	} else {
		log.Printf("%sowningElmentIdentifier: %s \n", prefix, vPtr.owningElement.getId().String())
	}
}

func (el *value) recoverValueFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.baseElement.recoverBaseElementFields(unmarshaledData)
}

type Value interface {
	BaseElement
}
