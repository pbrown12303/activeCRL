package core

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type value struct {
	baseElement
	owningElement Element
}

func (vPtr *value) GetOwningElement() Element {
	return vPtr.owningElement
}

func (vPtr *value) initializeValue() {
	vPtr.initializeBaseElement()
}

func (vPtr *value) isEquivalent(be *value) bool {
	if (vPtr.owningElement == nil && be.owningElement != nil) || (vPtr.owningElement != nil && be.owningElement == nil) {
		fmt.Printf("Equivalence failed: Value's Owning Elements do not match - one is nil and the other is not \n")
		fmt.Printf("First value: %#v \n", vPtr)
		fmt.Printf("First value's owner: \n")
		Print(vPtr.owningElement, "   ")
		fmt.Printf("Second value: %#v \n", be)
		fmt.Printf("Second value's owner: \n")
		Print(be.owningElement, "   ")
		return false
	}
	if vPtr.owningElement != nil && be.owningElement != nil && vPtr.owningElement.GetId() != be.owningElement.GetId() {
		fmt.Printf("Equivalence failed: Value's Owning Elements do not match - they have different identifiers\n")
		fmt.Printf("First value's owner: \n")
		Print(vPtr.owningElement, "   ")
		fmt.Printf("Second value's owner: \n")
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
}

func (el *value) recoverValueFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.baseElement.recoverBaseElementFields(unmarshaledData)
}

type Value interface {
	BaseElement
}
