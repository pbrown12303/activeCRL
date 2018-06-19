// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

func (vPtr *value) getOwningElement(hl *HeldLocks) Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(vPtr)
	return vPtr.owningElement
}

func (vPtr *value) getUri(hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(vPtr)
	return vPtr.uri
}

func (vPtr *value) initializeValue(uri ...string) {
	vPtr.initializeBaseElement(uri...)
}

func (vPtr *value) isEquivalent(be *value, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(vPtr)
	hl.LockBaseElement(be)
	if (vPtr.owningElement == nil && be.owningElement != nil) || (vPtr.owningElement != nil && be.owningElement == nil) {
		log.Printf("Equivalence failed: Value's Owning Elements do not match - one is nil and the other is not \n")
		log.Printf("First value: %#v \n", vPtr)
		log.Printf("First value's owner: \n")
		Print(vPtr.owningElement, "   ", hl)
		log.Printf("Second value: %#v \n", be)
		log.Printf("Second value's owner: \n")
		Print(be.owningElement, "   ", hl)
		return false
	}
	if vPtr.owningElement != nil && be.owningElement != nil && vPtr.owningElement.GetId(hl) != be.owningElement.GetId(hl) {
		log.Printf("Equivalence failed: Value's Owning Elements do not match - they have different identifiers\n")
		log.Printf("First value's owner: \n")
		Print(vPtr.owningElement, "   ", hl)
		log.Printf("Second value's owner: \n")
		Print(be.owningElement, "   ", hl)
		return false
	}
	var baseElementPtr *baseElement = &vPtr.baseElement
	return baseElementPtr.isEquivalent(&be.baseElement, hl)
}

func (vPtr *value) marshalValueFields(buffer *bytes.Buffer) error {
	vPtr.baseElement.marshalBaseElementFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"Uri\":\"%s\",", vPtr.uri))
	return nil
}

func (vPtr *value) printValue(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(vPtr)
	vPtr.printBaseElement(prefix, hl)
	log.Printf("%s  uri: %s \n", prefix, vPtr.getUri(hl))
	if vPtr.getOwningElement(hl) == nil {
		log.Printf("%s  owningElmentIdentifier: %s \n", prefix, "nil")
	} else {
		log.Printf("%s  owningElmentIdentifier: %s \n", prefix, vPtr.owningElement.GetId(hl).String())
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

func (el *value) setOwningElement(oe Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(el)
	el.owningElement = oe
}

type Value interface {
	BaseElement
	getLabel(*HeldLocks) string
	getOwningElement(*HeldLocks) Element
	getUri(*HeldLocks) string
	setOwningElement(Element, *HeldLocks)
	setUri(string, *HeldLocks)
}
