// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

func (pPtr *pointer) initializePointer(uri ...string) {
	pPtr.initializeValue(uri...)
}

func (pPtr *pointer) isEquivalent(be *pointer, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(pPtr)
	hl.LockBaseElement(be)
	var valuePtr *value = &pPtr.value
	return valuePtr.isEquivalent(&be.value, hl)
}

func (elPtr *pointer) marshalPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.value.marshalValueFields(buffer)
	return err
}

func (pPtr *pointer) printPointer(prefix string, hl *HeldLocks) {
	pPtr.printValue(prefix, hl)
}

func (el *pointer) recoverPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.value.recoverValueFields(unmarshaledData)
}

type Pointer interface {
	Value
}
