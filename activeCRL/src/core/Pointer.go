package core

import (
	"bytes"
	"encoding/json"
)

type pointer struct {
	value
}

func (pPtr *pointer) cloneAttributes(source pointer) {
	pPtr.value.cloneAttributes(source.value)
}

func (pPtr *pointer) initializePointer() {
	pPtr.initializeValue()
}

func (pPtr *pointer) isEquivalent(be *pointer) bool {
	var valuePtr *value = &pPtr.value
	return valuePtr.isEquivalent(&be.value)
}

func (elPtr *pointer) marshalPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.value.marshalValueFields(buffer)
	return err
}

func (pPtr *pointer) printPointer(prefix string) {
	pPtr.printValue(prefix)
}

func (el *pointer) recoverPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.value.recoverValueFields(unmarshaledData)
}

type Pointer interface {
	Value
}
