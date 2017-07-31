package core

import (
	"bytes"
	"encoding/json"
)

type reference struct {
	element
}

func (elPtr *reference) cloneAttributes(source reference) {
	elPtr.element.cloneAttributes(source.element)
}

func (elPtr *reference) initializeReference() {
	elPtr.initializeElement()
}
func (bePtr *reference) isEquivalent(be *reference, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	var elementPtr *element = &bePtr.element
	return elementPtr.isEquivalent(&be.element, hl)
}

func (elPtr *reference) marshalReferenceFields(buffer *bytes.Buffer) error {
	return elPtr.element.marshalElementFields(buffer)
}

func (elPtr *reference) printReference(prefix string, hl *HeldLocks) {
	elPtr.printElement(prefix, hl)
}

func (el *reference) recoverReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.element.recoverElementFields(unmarshaledData)
}

type Reference interface {
	Element
}
