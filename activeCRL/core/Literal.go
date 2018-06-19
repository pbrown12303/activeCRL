// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type literal struct {
	value
	literalValue string
}

func (lPtr *literal) clone() *literal {
	var clone literal
	clone.cloneAttributes(*lPtr)
	return &clone
}

func (lPtr *literal) cloneAttributes(source literal) {
	lPtr.value.cloneAttributes(source.value)
	lPtr.literalValue = source.literalValue
}

func (lPtr *literal) GetLiteralValue(hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lPtr)
	return lPtr.literalValue
}

func (lPtr *literal) getLabel(hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lPtr)
	return lPtr.GetLiteralValue(hl)
}

func (lPtr *literal) initializeLiteral(uri ...string) {
	lPtr.initializeValue(uri...)
}

func (lPtr *literal) isEquivalent(lit *literal, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lPtr)
	if lPtr.literalValue != lit.literalValue {
		log.Printf("Literal values not equivalent - v1: %s v2: %s \n", lPtr.literalValue, lit.literalValue)
		return false
	}
	var valuePtr *value = &lPtr.value
	return valuePtr.isEquivalent(&lit.value, hl)
}

func (lPtr *literal) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(lPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := lPtr.marshalLiteralFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (lPtr *literal) marshalLiteralFields(buffer *bytes.Buffer) error {
	lPtr.value.marshalValueFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"LiteralValue\":\"%s\"", lPtr.literalValue))
	return nil
}

func (lPtr *literal) printLiteral(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lPtr)
	lPtr.printValue(prefix, hl)
	log.Printf("%s  literalValue: %s \n", prefix, lPtr.literalValue)
	owningElementType := ""
	if lPtr.owningElement != nil {
		owningElementType = reflect.TypeOf(lPtr.owningElement).String()
	}
	log.Printf("%s  owningElementType: %s \n", prefix, owningElementType)
}

func (lPtr *literal) recoverLiteralFields(unmarshaledData *map[string]json.RawMessage) error {
	err := lPtr.recoverValueFields(unmarshaledData)
	if err != nil {
		fmt.Printf("Literal's Recovery of ValueFields failed\n")
		return err
	}
	// Element ID
	var recoveredLiteralValue string
	err = json.Unmarshal((*unmarshaledData)["LiteralValue"], &recoveredLiteralValue)
	if err != nil {
		fmt.Printf("ElementPointer's Recovery of ElementId failed\n")
		return err
	}
	lPtr.literalValue = recoveredLiteralValue
	return nil
}

func (lPtr *literal) SetLiteralValue(newValue string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lPtr)
	if lPtr.literalValue != newValue {
		preChange(lPtr, hl)
		lPtr.literalValue = newValue
		notification := NewChangeNotification(lPtr, MODIFY, "SetLiteralValue", nil)
		postChange(lPtr, notification, hl)
	}
}

func (lPtr *literal) SetOwningElement(el Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lPtr)
	if lPtr.owningElement != el {
		if lPtr.owningElement != nil {
			removeOwnedBaseElement(lPtr.owningElement, lPtr, hl)
		}

		preChange(lPtr, hl)
		lPtr.owningElement = el
		notification := NewChangeNotification(lPtr, MODIFY, "SetOwningElement", nil)
		postChange(lPtr, notification, hl)

		if lPtr.owningElement != nil {
			addOwnedBaseElement(lPtr.owningElement, lPtr, hl)
		}
	}
}

// internalSetOwningElement() is an internal function used only in unmarshal
func (lPtr *literal) internalSetOwningElement(el Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lPtr)
	lPtr.owningElement = el
	if lPtr.owningElement != nil {
		lPtr.owningElement.internalAddOwnedBaseElement(lPtr, hl)
	}
}

func (lPtr *literal) setUri(uri string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lPtr)
	if uri != lPtr.uri {
		preChange(lPtr, hl)
		lPtr.uri = uri
		notification := NewChangeNotification(lPtr, MODIFY, "setUri", nil)
		postChange(lPtr, notification, hl)
	}
}

type Literal interface {
	Value
	GetLiteralValue(*HeldLocks) string
	SetLiteralValue(string, *HeldLocks)
}
