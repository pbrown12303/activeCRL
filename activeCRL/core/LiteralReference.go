// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

func (lrPtr *literalReference) GetReferencedLiteral(hl *HeldLocks) Literal {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lrPtr)
	rep := lrPtr.GetLiteralPointer(hl)
	if rep != nil {
		return rep.GetLiteral(hl)
	}
	return nil
}

func (lrPtr *literalReference) GetLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lrPtr)
	for _, be := range lrPtr.ownedBaseElements {
		switch be.(type) {
		case LiteralPointer:
			if be.(LiteralPointer).GetLiteralPointerRole(hl) == VALUE {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

func (elPtr *literalReference) initializeLiteralReference() {
	elPtr.initializeReference()
}

func (bePtr *literalReference) isEquivalent(be *literalReference, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	var elementPtr *element = &bePtr.element
	return elementPtr.isEquivalent(&be.element, hl)
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

func (elPtr *literalReference) printLiteralReference(prefix string, hl *HeldLocks) {
	elPtr.printElement(prefix, hl)
}

func (el *literalReference) recoverLiteralReferenceFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.element.recoverElementFields(unmarshaledData)
}

func (lrPtr *literalReference) SetReferencedLiteral(el Literal, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lrPtr)
	if lrPtr.GetReferencedLiteral(hl) != el {
		ep := lrPtr.GetLiteralPointer(hl)
		if ep == nil {
			ep = lrPtr.uOfD.NewValueLiteralPointer(hl)
			SetOwningElement(ep, lrPtr, hl)
		}
		ep.SetLiteral(el, hl)
	}
}

type LiteralReference interface {
	Reference
	GetReferencedLiteral(*HeldLocks) Literal
	GetLiteralPointer(*HeldLocks) LiteralPointer
	SetReferencedLiteral(Literal, *HeldLocks)
}
