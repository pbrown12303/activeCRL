package core

import (
	"bytes"
	"encoding/json"
)

type reference struct {
	element
}

func (elPtr *reference) initializeReference() {
	elPtr.initializeElement()
}
func (bePtr *reference) isEquivalent(be *reference) bool {
	var elementPtr *element = &bePtr.element
	return elementPtr.isEquivalent(&be.element)
}

func (elPtr *reference) marshalReferenceFields(buffer *bytes.Buffer) error {
	return elPtr.element.marshalElementFields(buffer)
}

func (elPtr *reference) printReference(prefix string) {
	elPtr.printElement(prefix)
}

func (el *reference) recoverReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.element.recoverElementFields(unmarshaledData)
}

type Reference interface {
	Element
}
