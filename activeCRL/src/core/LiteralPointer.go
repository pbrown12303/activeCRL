package core

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (lpPtr *literalPointer) GetLiteral() Literal {
	lpPtr.Lock()
	defer lpPtr.Unlock()
	if lpPtr.literal == nil && lpPtr.getLiteralIdentifier() != uuid.Nil && lpPtr.uOfD != nil {
		lpPtr.literal = lpPtr.uOfD.getLiteral(lpPtr.getLiteralIdentifier().String())
	}
	return lpPtr.literal
}

func NewNameLiteralPointer(uOfD *UniverseOfDiscourse) LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.setLiteralPointerRole(NAME)
	uOfD.AddBaseElement(&lp)
	return &lp
}

func NewDefinitionLiteralPointer(uOfD *UniverseOfDiscourse) LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.setLiteralPointerRole(DEFINITION)
	uOfD.AddBaseElement(&lp)
	return &lp
}

func NewUriLiteralPointer(uOfD *UniverseOfDiscourse) LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.setLiteralPointerRole(URI)
	uOfD.AddBaseElement(&lp)
	return &lp
}

func NewValueLiteralPointer(uOfD *UniverseOfDiscourse) LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.setLiteralPointerRole(VALUE)
	uOfD.AddBaseElement(&lp)
	return &lp
}

func (lpPtr *literalPointer) GetLiteralIdentifier() uuid.UUID {
	lpPtr.Lock()
	defer lpPtr.Unlock()
	return lpPtr.getLiteralIdentifier()
}

func (lpPtr *literalPointer) getLiteralIdentifier() uuid.UUID {
	return lpPtr.literalId
}

func (lpPtr *literalPointer) getLiteralPointerRole() LiteralPointerRole {
	return lpPtr.literalPointerRole
}

func (lpPtr *literalPointer) GetLiteralVersion() int {
	lpPtr.Lock()
	defer lpPtr.Unlock()
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
		fmt.Printf("Equivalence failed: indicated literal ids do not match \n")
		return false
	}
	if bePtr.literalVersion != be.literalVersion {
		fmt.Printf("Equivalence failed: indicated literal versions do not match \n")
		return false
	}
	if bePtr.literalPointerRole != be.literalPointerRole {
		fmt.Printf("Equivalence failed: literal pointer roles do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer)
}

func (elPtr *literalPointer) MarshalJSON() ([]byte, error) {
	elPtr.Lock()
	defer elPtr.Unlock()
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
	fmt.Printf("%sIndicated LiteralId: %s \n", prefix, lpPtr.literalId.String())
	fmt.Printf("%sIndicated LiteralVersion: %d \n", prefix, lpPtr.literalVersion)
	fmt.Printf("%sLiteralPointerRole: %d \n", prefix, lpPtr.literalPointerRole)
}

func (lp *literalPointer) recoverLiteralPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	err := lp.pointer.recoverPointerFields(unmarshaledData)
	if err != nil {
		fmt.Printf("LiteralPointer's Recovery of PointerFields failed\n")
		return err
	}
	// Literal ID
	var recoveredLiteralId string
	err = json.Unmarshal((*unmarshaledData)["LiteralId"], &recoveredLiteralId)
	if err != nil {
		fmt.Printf("LiteralPointer's Recovery of LiteralId failed\n")
		return err
	}
	lp.literalId, err = uuid.FromString(recoveredLiteralId)
	if err != nil {
		fmt.Printf("LiteralPointer's conversion of LiteralId failed\n")
		return err
	}
	// Version
	var recoveredLiteralVersion string
	err = json.Unmarshal((*unmarshaledData)["LiteralVersion"], &recoveredLiteralVersion)
	if err != nil {
		fmt.Printf("LiteralPointer's Recovery of LiteralVersion failed\n")
		return err
	}
	lp.literalVersion, err = strconv.Atoi(recoveredLiteralVersion)
	if err != nil {
		fmt.Printf("Conversion of LiteralPointer.literalVersion failed\n")
		return err
	}
	// Literal pointer role
	var recoveredLiteralPointerRole string
	err = json.Unmarshal((*unmarshaledData)["LiteralPointerRole"], &recoveredLiteralPointerRole)
	if err != nil {
		fmt.Printf("LiteralPointer's Recovery of LiteralPointerRole failed\n")
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
	lpPtr.Lock()
	defer lpPtr.Unlock()
	if literal != nil {
		literal.Lock()
		defer literal.Unlock()
	}
	lpPtr.setLiteral(literal)
}

func (lpPtr *literalPointer) setLiteral(literal Literal) {
	if literal != lpPtr.literal {
		lpPtr.literal = literal
		if literal != nil {
			lpPtr.literalId = literal.getId()
			lpPtr.literalVersion = literal.getVersion()
		} else {
			lpPtr.literalId = uuid.Nil
			lpPtr.literalVersion = 0
		}
	}
}

func (lpPtr *literalPointer) setLiteralPointerRole(role LiteralPointerRole) {
	lpPtr.literalPointerRole = role
}

func (lpPtr *literalPointer) SetOwningElement(element Element) {
	lpPtr.Lock()
	defer lpPtr.Unlock()
	if element != nil {
		element.Lock()
		defer element.Unlock()
	}
	lpPtr.setOwningElement(element)
}

func (lpPtr *literalPointer) setOwningElement(element Element) {
	if element != lpPtr.getOwningElement() {
		if lpPtr.getOwningElement() != nil {
			lpPtr.getOwningElement().removeOwnedBaseElement(lpPtr)
		}
		lpPtr.owningElement = element
		if lpPtr.getOwningElement() != nil {
			lpPtr.getOwningElement().addOwnedBaseElement(lpPtr)
		}
	}
}

type LiteralPointer interface {
	Pointer
	GetLiteral() Literal
	GetLiteralIdentifier() uuid.UUID
	getLiteralPointerRole() LiteralPointerRole
	GetLiteralVersion() int
	setLiteral(Literal)
	SetLiteral(Literal)
}
