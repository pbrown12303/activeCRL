package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/satori/go.uuid"
)

type literalPointerPointer struct {
	pointer
	literalPointer        LiteralPointer
	literalPointerId      uuid.UUID
	literalPointerVersion int
}

func (pllPtr *literalPointerPointer) clone() *literalPointerPointer {
	var clone literalPointerPointer
	clone.cloneAttributes(*pllPtr)
	return &clone
}

func (pllPtr *literalPointerPointer) cloneAttributes(source literalPointerPointer) {
	pllPtr.pointer.cloneAttributes(source.pointer)
	pllPtr.literalPointer = source.literalPointer
	pllPtr.literalPointerId = source.literalPointerId
	pllPtr.literalPointerVersion = source.literalPointerVersion
}

func (pllPtr *literalPointerPointer) GetLiteralPointer() LiteralPointer {
	pllPtr.traceableLock()
	defer pllPtr.traceableUnlock()
	return pllPtr.getLiteralPointer()
}

func (pllPtr *literalPointerPointer) getLiteralPointer() LiteralPointer {
	if pllPtr.literalPointer == nil && pllPtr.getLiteralPointerIdentifier() != uuid.Nil && pllPtr.uOfD != nil {
		pllPtr.literalPointer = pllPtr.uOfD.getLiteralPointer(pllPtr.getLiteralPointerIdentifier().String())
	}
	return pllPtr.literalPointer
}

func (pllPtr *literalPointerPointer) GetName() string {
	return "literalPointerPointer"
}

func (pllPtr *literalPointerPointer) GetLiteralPointerIdentifier() uuid.UUID {
	pllPtr.traceableLock()
	defer pllPtr.traceableUnlock()
	return pllPtr.getLiteralPointerIdentifier()
}

func (pllPtr *literalPointerPointer) getLiteralPointerIdentifier() uuid.UUID {
	return pllPtr.literalPointerId
}

func (pllPtr *literalPointerPointer) GetLiteralPointerVersion() int {
	pllPtr.traceableLock()
	defer pllPtr.traceableUnlock()
	return pllPtr.getLiteralPointerVersion()
}

func (pllPtr *literalPointerPointer) getLiteralPointerVersion() int {
	return pllPtr.literalPointerVersion
}

func (pllPtr *literalPointerPointer) initializeLiteralPointerPointer() {
	pllPtr.initializePointer()
}

func (bePtr *literalPointerPointer) isEquivalent(be *literalPointerPointer) bool {
	if bePtr.literalPointerId != be.literalPointerId {
		fmt.Printf("Equivalence failed: indicated literalPointerPointer ids do not match \n")
		return false
	}
	if bePtr.literalPointerVersion != be.literalPointerVersion {
		fmt.Printf("Equivalence failed: indicated literalPointerPointer versions do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer)
}

func (elPtr *literalPointerPointer) MarshalJSON() ([]byte, error) {
	elPtr.traceableLock()
	defer elPtr.traceableUnlock()
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.maarshalLiteralPointerPointerFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *literalPointerPointer) maarshalLiteralPointerPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.pointer.marshalPointerFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"LiteralPointerId\":\"%s\",", elPtr.literalPointerId.String()))
	buffer.WriteString(fmt.Sprintf("\"LiteralPointerVersion\":\"%d\"", elPtr.literalPointerVersion))
	return err
}

func (pllPtr *literalPointerPointer) printLiteralPointerPointer(prefix string) {
	pllPtr.printPointer(prefix)
	fmt.Printf("%sIndicated LiteralPointerID: %s \n", prefix, pllPtr.literalPointerId.String())
	fmt.Printf("%sIndicated LiteralPointerVersion: %d \n", prefix, pllPtr.literalPointerVersion)
}

func (ep *literalPointerPointer) recoverLiteralPointerPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	err := ep.pointer.recoverPointerFields(unmarshaledData)
	if err != nil {
		fmt.Printf("LiteralPointerPointer's Recovery of PointerFields failed\n")
		return err
	}
	// LiteralPointer ID
	var recoveredLiteralPointerId string
	err = json.Unmarshal((*unmarshaledData)["LiteralPointerId"], &recoveredLiteralPointerId)
	if err != nil {
		fmt.Printf("LiteralPointerPointer's Recovery of LiteralPointerId failed\n")
		return err
	}
	ep.literalPointerId, err = uuid.FromString(recoveredLiteralPointerId)
	if err != nil {
		fmt.Printf("LiteralPointerPointer's conversion of LiteralPointerId failed\n")
		return err
	}
	// Version
	var recoveredLiteralPointerVersion string
	err = json.Unmarshal((*unmarshaledData)["LiteralPointerVersion"], &recoveredLiteralPointerVersion)
	if err != nil {
		fmt.Printf("LiteralPointerPointer's Recovery of LiteralPointerVersion failed\n")
		return err
	}
	ep.literalPointerVersion, err = strconv.Atoi(recoveredLiteralPointerVersion)
	if err != nil {
		fmt.Printf("Conversion of LiteralPointerPointer.literalPointerVersion failed\n")
		return err
	}
	return nil
}

func (pllPtr *literalPointerPointer) SetLiteralPointer(literalPointer LiteralPointer) {
	pllPtr.traceableLock()
	defer pllPtr.traceableUnlock()
	if literalPointer != nil {
		literalPointer.traceableLock()
		defer literalPointer.traceableUnlock()
	}
	pllPtr.setLiteralPointer(literalPointer)
}

func (pllPtr *literalPointerPointer) setLiteralPointer(literalPointer LiteralPointer) {
	if literalPointer != pllPtr.literalPointer {
		preChange(pllPtr)
		pllPtr.literalPointer = literalPointer
		if literalPointer != nil {
			pllPtr.literalPointerId = literalPointer.getId()
			pllPtr.literalPointerVersion = literalPointer.getVersion()
		} else {
			pllPtr.literalPointerId = uuid.Nil
			pllPtr.literalPointerVersion = 0
		}
		postChange(pllPtr)
	}
}

func (pllPtr *literalPointerPointer) SetOwningElement(element Element) {
	pllPtr.traceableLock()
	defer pllPtr.traceableUnlock()
	pllPtr.setOwningElement(element)
}

func (pllPtr *literalPointerPointer) setOwningElement(element Element) {
	if element != pllPtr.getOwningElement() {
		if pllPtr.getOwningElement() != nil {
			pllPtr.getOwningElement().removeOwnedBaseElement(pllPtr)
		}

		preChange(pllPtr)
		pllPtr.owningElement = element
		postChange(pllPtr)

		if pllPtr.getOwningElement() != nil {
			pllPtr.getOwningElement().addOwnedBaseElement(pllPtr)
		}
	}
}

// internalSetOwningElement() is an internal function used only in unmarshal
func (pllPtr *literalPointerPointer) internalSetOwningElement(element Element) {
	if element != pllPtr.GetOwningElement() {
		pllPtr.owningElement = element
		if pllPtr.GetOwningElement() != nil {
			pllPtr.GetOwningElement().internalAddOwnedBaseElement(pllPtr)
		}
	}
}

func (lpPtr *literalPointerPointer) SetUri(uri string) {
	lpPtr.traceableLock()
	defer lpPtr.traceableUnlock()
	lpPtr.setUri(uri)
}

func (lpPtr *literalPointerPointer) setUri(uri string) {
	preChange(lpPtr)
	lpPtr.uri = uri
	postChange(lpPtr)
}

type LiteralPointerPointer interface {
	Pointer
	GetLiteralPointer() LiteralPointer
	getLiteralPointer() LiteralPointer
	GetLiteralPointerIdentifier() uuid.UUID
	GetLiteralPointerVersion() int
	setLiteralPointer(LiteralPointer)
	SetLiteralPointer(LiteralPointer)
}
