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

type literalPointer struct {
	pointer
	literal            Literal
	literalId          uuid.UUID
	literalVersion     int
	literalPointerRole LiteralPointerRole
}

func NewNameLiteralPointer(uOfD *UniverseOfDiscourse) LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = NAME
	uOfD.AddBaseElement(&lp)
	return &lp
}

func NewDefinitionLiteralPointer(uOfD *UniverseOfDiscourse) LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = DEFINITION
	uOfD.AddBaseElement(&lp)
	return &lp
}

func NewUriLiteralPointer(uOfD *UniverseOfDiscourse) LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = URI
	uOfD.AddBaseElement(&lp)
	return &lp
}

func NewValueLiteralPointer(uOfD *UniverseOfDiscourse) LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = VALUE
	uOfD.AddBaseElement(&lp)
	return &lp
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
func (lpPtr *literalPointer) GetLiteral() Literal {
	lpPtr.traceableLock()
	defer lpPtr.traceableUnlock()
	if lpPtr.literal == nil && lpPtr.getLiteralIdentifier() != uuid.Nil && lpPtr.uOfD != nil {
		lpPtr.literal = lpPtr.uOfD.getLiteral(lpPtr.getLiteralIdentifier().String())
	}
	if lpPtr.literal != nil {
		lpPtr.literal.traceableLock()
		defer lpPtr.literal.traceableUnlock()
	}
	return lpPtr.getLiteral()
}

// getLiteral() is a non-locking internal function that simply returns the literal pointed to by this
// literal pointer. If the literal value is nil, it checks to see whether there is an identifier for it present
// and, if so, attemtps to find it using the uOfD.
func (lpPtr *literalPointer) getLiteral() Literal {
	if lpPtr.literal == nil && lpPtr.getLiteralIdentifier() != uuid.Nil && lpPtr.uOfD != nil {
		lpPtr.literal = lpPtr.uOfD.getLiteral(lpPtr.getLiteralIdentifier().String())
	}
	return lpPtr.literal
}

func (lpPtr *literalPointer) GetLiteralIdentifier() uuid.UUID {
	lpPtr.traceableLock()
	defer lpPtr.traceableUnlock()
	return lpPtr.getLiteralIdentifier()
}

func (lpPtr *literalPointer) getLiteralIdentifier() uuid.UUID {
	return lpPtr.literalId
}

func (lpPtr *literalPointer) getLiteralPointerRole() LiteralPointerRole {
	return lpPtr.literalPointerRole
}

func (lpPtr *literalPointer) GetLiteralVersion() int {
	lpPtr.traceableLock()
	defer lpPtr.traceableUnlock()
	return lpPtr.getLiteralVersion()
}

func (lpPtr *literalPointer) getLiteralVersion() int {
	return lpPtr.literalVersion
}

func (lpPtr *literalPointer) GetName() string {
	switch lpPtr.getLiteralPointerRole() {
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

func (lpPtr *literalPointer) initializeLiteralPointer() {
	lpPtr.initializePointer()
}

func (bePtr *literalPointer) isEquivalent(be *literalPointer) bool {
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
	return pointerPtr.isEquivalent(&be.pointer)
}

func (elPtr *literalPointer) MarshalJSON() ([]byte, error) {
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
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

func (lpPtr *literalPointer) printLiteralPointer(prefix string) {
	lpPtr.printPointer(prefix)
	log.Printf("%sIndicated LiteralId: %s \n", prefix, lpPtr.literalId.String())
	log.Printf("%sIndicated LiteralVersion: %d \n", prefix, lpPtr.literalVersion)
	log.Printf("%sLiteralPointerRole: %d \n", prefix, lpPtr.literalPointerRole)
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

func (lpPtr *literalPointer) SetLiteral(literal Literal) {
	lpPtr.traceableLock()
	defer lpPtr.traceableUnlock()
	if literal != nil {
		literal.traceableLock()
		defer literal.traceableUnlock()
	}
	lpPtr.setLiteral(literal)
}

func (lpPtr *literalPointer) setLiteral(literal Literal) {
	if literal != lpPtr.literal {
		preChange(lpPtr)
		lpPtr.literal = literal
		if literal != nil {
			lpPtr.literalId = literal.getId()
			lpPtr.literalVersion = literal.getVersion()
		} else {
			lpPtr.literalId = uuid.Nil
			lpPtr.literalVersion = 0
		}
		postChange(lpPtr)
	}
}

func (lpPtr *literalPointer) SetOwningElement(element Element) {
	lpPtr.traceableLock()
	defer lpPtr.traceableUnlock()
	if element != nil {
		element.traceableLock()
		defer element.traceableUnlock()
	}
	lpPtr.setOwningElement(element)
}

func (lpPtr *literalPointer) setOwningElement(element Element) {
	if element != lpPtr.getOwningElement() {
		if lpPtr.getOwningElement() != nil {
			lpPtr.getOwningElement().removeOwnedBaseElement(lpPtr)
		}

		preChange(lpPtr)
		lpPtr.owningElement = element
		postChange(lpPtr)

		if lpPtr.getOwningElement() != nil {
			lpPtr.getOwningElement().addOwnedBaseElement(lpPtr)
		}
	}
}

// internalSetOwningElement() is an internal function used only in unmarshal
func (lpPtr *literalPointer) internalSetOwningElement(element Element) {
	if element != lpPtr.getOwningElement() {
		lpPtr.owningElement = element
		if lpPtr.getOwningElement() != nil {
			lpPtr.getOwningElement().internalAddOwnedBaseElement(lpPtr)
		}
	}
}

func (lpPtr *literalPointer) SetUri(uri string) {
	lpPtr.traceableLock()
	defer lpPtr.traceableUnlock()
	lpPtr.setUri(uri)
}

func (lpPtr *literalPointer) setUri(uri string) {
	preChange(lpPtr)
	lpPtr.uri = uri
	postChange(lpPtr)
}

type LiteralPointer interface {
	Pointer
	GetLiteral() Literal
	getLiteral() Literal
	GetLiteralIdentifier() uuid.UUID
	getLiteralPointerRole() LiteralPointerRole
	GetLiteralVersion() int
	setLiteral(Literal)
	SetLiteral(Literal)
}
