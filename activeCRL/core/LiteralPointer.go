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
	"strconv"

	"github.com/satori/go.uuid"
)

type LiteralPointerRole int

const (
	NAME LiteralPointerRole = 1 + iota
	DEFINITION
	URI
	VALUE
)

func (lpr LiteralPointerRole) RoleToString() string {
	role := ""
	switch lpr {
	case NAME:
		role = "NAME"
	case DEFINITION:
		role = "DEFINITION"
	case URI:
		role = "URI"
	case VALUE:
		role = "VALUE"
	}
	return role
}

type literalPointer struct {
	pointer
	literal            Literal
	literalId          uuid.UUID
	literalVersion     int
	literalPointerRole LiteralPointerRole
}

func (lpPtr *literalPointer) clone() *literalPointer {
	var clone literalPointer
	clone.cloneAttributes(*lpPtr)
	return &clone
}

func (lpPtr *literalPointer) cloneAttributes(source literalPointer) {
	lpPtr.pointer.cloneAttributes(source.pointer)
	lpPtr.literal = source.literal
	lpPtr.literalId = source.literalId
	lpPtr.literalVersion = source.literalVersion
	lpPtr.literalPointerRole = source.literalPointerRole
}

// GetLiteral() locks the literal pointer and the literal to which it points. If the literal
// valaue is nil, it checks to see whether there is an identifier for it present and, if so,
// attempts to find it using the uOfD. It then returns the result of calling the non-locking getLiteral()
func (lpPtr *literalPointer) GetLiteral(hl *HeldLocks) Literal {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	if lpPtr.literal == nil && lpPtr.GetLiteralId(hl) != uuid.Nil && lpPtr.uOfD != nil {
		lpPtr.literal = lpPtr.uOfD.GetLiteral(lpPtr.GetLiteralId(hl))
	}
	return lpPtr.literal
}

func (lpPtr *literalPointer) GetLiteralId(hl *HeldLocks) uuid.UUID {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	return lpPtr.literalId
}

func (lpPtr *literalPointer) GetLiteralPointerRole(hl *HeldLocks) LiteralPointerRole {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	return lpPtr.literalPointerRole
}

func (lpPtr *literalPointer) GetLiteralVersion(hl *HeldLocks) int {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	return lpPtr.literalVersion
}

func (lpPtr *literalPointer) getName(hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	switch lpPtr.GetLiteralPointerRole(hl) {
	case NAME:
		return "name"
	case DEFINITION:
		return "definition"
	case URI:
		return "uri"
	case VALUE:
		return "value"
	}
	return ""
}

func (lpPtr *literalPointer) initializeLiteralPointer(uri ...string) {
	lpPtr.initializePointer(uri...)
}

func (bePtr *literalPointer) isEquivalent(be *literalPointer, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	if bePtr.literalId != be.literalId {
		log.Printf("Equivalence failed: indicated literal ids do not match \n")
		return false
	}
	if bePtr.literalVersion != be.literalVersion {
		log.Printf("Equivalence failed: indicated literal versions do not match \n")
		return false
	}
	if bePtr.literalPointerRole != be.literalPointerRole {
		log.Printf("Equivalence failed: literal pointer roles do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer, hl)
}

func (elPtr *literalPointer) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalLiteralPointerFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *literalPointer) marshalLiteralPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.pointer.marshalPointerFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"LiteralId\":\"%s\",", elPtr.literalId.String()))
	buffer.WriteString(fmt.Sprintf("\"LiteralVersion\":\"%d\",", elPtr.literalVersion))
	switch elPtr.literalPointerRole {
	case VALUE:
		buffer.WriteString(fmt.Sprintf("\"LiteralPointerRole\":\"%s\"", "VALUE"))
	case URI:
		buffer.WriteString(fmt.Sprintf("\"LiteralPointerRole\":\"%s\"", "URI"))
	case NAME:
		buffer.WriteString(fmt.Sprintf("\"LiteralPointerRole\":\"%s\"", "NAME"))
	case DEFINITION:
		buffer.WriteString(fmt.Sprintf("\"LiteralPointerRole\":\"%s\"", "DEFINITION"))
	}
	return err
}

func (lpPtr *literalPointer) printLiteralPointer(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	lpPtr.printPointer(prefix, hl)
	log.Printf("%s  Indicated LiteralId: %s \n", prefix, lpPtr.literalId.String())
	log.Printf("%s  Indicated LiteralVersion: %d \n", prefix, lpPtr.literalVersion)
	role := ""
	switch lpPtr.literalPointerRole {
	case VALUE:
		role = "VALUE"
	case URI:
		role = "URI"
	case NAME:
		role = "NAME"
	case DEFINITION:
		role = "DEFINITION"
	}
	log.Printf("%s  LiteralPointerRole: %s \n", prefix, role)
}

func (lp *literalPointer) recoverLiteralPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	err := lp.pointer.recoverPointerFields(unmarshaledData)
	if err != nil {
		log.Printf("LiteralPointer's Recovery of PointerFields failed\n")
		return err
	}
	// Literal ID
	var recoveredLiteralId string
	err = json.Unmarshal((*unmarshaledData)["LiteralId"], &recoveredLiteralId)
	if err != nil {
		log.Printf("LiteralPointer's Recovery of LiteralId failed\n")
		return err
	}
	lp.literalId, err = uuid.FromString(recoveredLiteralId)
	if err != nil {
		log.Printf("LiteralPointer's conversion of LiteralId failed\n")
		return err
	}
	// Version
	var recoveredLiteralVersion string
	err = json.Unmarshal((*unmarshaledData)["LiteralVersion"], &recoveredLiteralVersion)
	if err != nil {
		log.Printf("LiteralPointer's Recovery of LiteralVersion failed\n")
		return err
	}
	lp.literalVersion, err = strconv.Atoi(recoveredLiteralVersion)
	if err != nil {
		log.Printf("Conversion of LiteralPointer.literalVersion failed\n")
		return err
	}
	// Literal pointer role
	var recoveredLiteralPointerRole string
	err = json.Unmarshal((*unmarshaledData)["LiteralPointerRole"], &recoveredLiteralPointerRole)
	if err != nil {
		log.Printf("LiteralPointer's Recovery of LiteralPointerRole failed\n")
		return err
	}
	switch recoveredLiteralPointerRole {
	case "VALUE":
		lp.literalPointerRole = VALUE
	case "URI":
		lp.literalPointerRole = URI
	case "NAME":
		lp.literalPointerRole = NAME
	case "DEFINITION":
		lp.literalPointerRole = DEFINITION
	}
	return nil
}

func (lpPtr *literalPointer) SetLiteral(literal Literal, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	if literal != lpPtr.literal {
		preChange(lpPtr, hl)
		if lpPtr.literal != nil {
			lpPtr.uOfD.removeLiteralListener(lpPtr.literal, lpPtr, hl)
		}
		lpPtr.literal = literal
		if literal != nil {
			lpPtr.literalId = literal.GetId(hl)
			lpPtr.literalVersion = literal.GetVersion(hl)
			lpPtr.uOfD.addLiteralListener(literal, lpPtr, hl)
		} else {
			lpPtr.literalId = uuid.Nil
			lpPtr.literalVersion = 0
		}
		notification := NewChangeNotification(lpPtr, MODIFY, nil)
		postChange(lpPtr, notification, hl)
	}
}

func (lpPtr *literalPointer) SetOwningElement(element Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	if element != lpPtr.owningElement {
		if lpPtr.owningElement != nil {
			removeOwnedBaseElement(lpPtr.owningElement, lpPtr, hl)
		}

		preChange(lpPtr, hl)
		lpPtr.owningElement = element
		notification := NewChangeNotification(lpPtr, MODIFY, nil)
		postChange(lpPtr, notification, hl)

		if lpPtr.owningElement != nil {
			addOwnedBaseElement(lpPtr.owningElement, lpPtr, hl)
		}
	}
}

// internalSetOwningElement() is an internal function used only in unmarshal
func (lpPtr *literalPointer) internalSetOwningElement(element Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	if element != lpPtr.getOwningElement(hl) {
		lpPtr.owningElement = element
		if lpPtr.getOwningElement(hl) != nil {
			lpPtr.getOwningElement(hl).internalAddOwnedBaseElement(lpPtr, hl)
		}
	}
}

func (lpPtr *literalPointer) setUri(uri string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if uri != lpPtr.uri {
		preChange(lpPtr, hl)
		lpPtr.uri = uri
		notification := NewChangeNotification(lpPtr, MODIFY, nil)
		postChange(lpPtr, notification, hl)
	}
}

type LiteralPointer interface {
	Pointer
	GetLiteral(*HeldLocks) Literal
	GetLiteralId(*HeldLocks) uuid.UUID
	GetLiteralPointerRole(*HeldLocks) LiteralPointerRole
	GetLiteralVersion(*HeldLocks) int
	SetLiteral(Literal, *HeldLocks)
}
