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

func (pllPtr *literalPointerPointer) GetLiteralPointer() LiteralPointer {
	if pllPtr.literalPointer == nil && pllPtr.GetLiteralPointerIdentifier() != uuid.Nil && pllPtr.uOfD != nil {
		pllPtr.literalPointer = pllPtr.uOfD.getLiteralPointer(pllPtr.GetLiteralPointerIdentifier().String())
	}
	return pllPtr.literalPointer
}

func (pllPtr *literalPointerPointer) GetName() string {
	return "literalPointerPointer"
}

func NewLiteralPointerPointer(uOfD *UniverseOfDiscourse) LiteralPointerPointer {
	var ep literalPointerPointer
	ep.initializeLiteralPointerPointer()
	uOfD.addBaseElement(&ep)
	return &ep
}

func (pllPtr *literalPointerPointer) GetLiteralPointerIdentifier() uuid.UUID {
	return pllPtr.literalPointerId
}

func (pllPtr *literalPointerPointer) GetLiteralPointerVersion() int {
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
	if literalPointer != pllPtr.literalPointer {
		pllPtr.literalPointer = literalPointer
		if literalPointer != nil {
			pllPtr.literalPointerId = literalPointer.GetId()
			pllPtr.literalPointerVersion = literalPointer.GetVersion()
		} else {
			pllPtr.literalPointerId = uuid.Nil
			pllPtr.literalPointerVersion = 0
		}
	}
}

func (pllPtr *literalPointerPointer) setOwningElement(element Element) {
	if element != pllPtr.GetOwningElement() {
		if pllPtr.GetOwningElement() != nil {
			pllPtr.GetOwningElement().removeOwnedBaseElement(pllPtr)
		}
		pllPtr.owningElement = element
		if pllPtr.GetOwningElement() != nil {
			pllPtr.GetOwningElement().addOwnedBaseElement(pllPtr)
		}
	}
}

type LiteralPointerPointer interface {
	Pointer
	GetLiteralPointer() LiteralPointer
	GetLiteralPointerIdentifier() uuid.UUID
	GetLiteralPointerVersion() int
	SetLiteralPointer(LiteralPointer)
}
