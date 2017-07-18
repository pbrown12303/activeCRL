package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type value struct {
	baseElement
	owningElement Element
	uri           string
}

func (vPtr *value) cloneAttributes(source value) {
	vPtr.baseElement.cloneAttributes(source.baseElement)
	vPtr.owningElement = source.owningElement
}

func (vPtr *value) GetOwningElement() Element {
	vPtr.TraceableLock()
	defer vPtr.TraceableUnlock()
	return vPtr.getOwningElement()
}

func (vPtr *value) getOwningElement() Element {
	return vPtr.owningElement
}

func (vPtr *value) GetUri() string {
	vPtr.TraceableLock()
	defer vPtr.TraceableUnlock()
	return vPtr.GetUriNoLock()
}

func (vPtr *value) GetUriNoLock() string {
	return vPtr.uri
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
	buffer.WriteString(fmt.Sprintf("\"Uri\":\"%s\",", vPtr.uri))
	return nil
}

func (vPtr *value) printValue(prefix string) {
	vPtr.printBaseElement(prefix)
	log.Printf("%suri: %s \n", prefix, vPtr.GetUriNoLock())
	if vPtr.getOwningElement() == nil {
		log.Printf("%sowningElmentIdentifier: %s \n", prefix, "nil")
	} else {
		log.Printf("%sowningElmentIdentifier: %s \n", prefix, vPtr.owningElement.getId().String())
	}
}

func (el *value) recoverValueFields(unmarshaledData *map[string]json.RawMessage) error {
	// Uri
	var recoveredUri string
	err := json.Unmarshal((*unmarshaledData)["Uri"], &recoveredUri)
	if err != nil {
		log.Printf("Recovery of Value.uri as string failed\n")
		return err
	}
	el.uri = recoveredUri
	return el.baseElement.recoverBaseElementFields(unmarshaledData)
}

type Value interface {
	BaseElement
}
