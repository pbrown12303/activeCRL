package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type literal struct {
	element
	LiteralValue string
}

func (lPtr *literal) clone(hl *HeldLocks) Literal {
	hl.ReadLockElement(lPtr)
	var clonedLiteral literal
	clonedLiteral.initializeLiteral("", "")
	clonedLiteral.cloneAttributes(lPtr, hl)
	return &clonedLiteral
}

func (lPtr *literal) cloneAttributes(source *literal, hl *HeldLocks) {
	lPtr.element.cloneAttributes(&source.element, hl)
	lPtr.LiteralValue = source.LiteralValue
	lPtr.element.cloneAttributes(&source.element, hl)
}

func (lPtr *literal) GetLiteralValue(hl *HeldLocks) string {
	hl.ReadLockElement(lPtr)
	return lPtr.LiteralValue
}

func (lPtr *literal) initializeLiteral(conceptID string, uri string) {
	lPtr.initializeElement(conceptID, uri)
}

func (lPtr *literal) isEquivalent(hl1 *HeldLocks, ref *literal, hl2 *HeldLocks, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) > 0 {
		print = printExceptions[0]
	}
	hl1.ReadLockElement(lPtr)
	hl2.ReadLockElement(ref)
	if ref.LiteralValue != lPtr.LiteralValue {
		if print {
			log.Printf("In literal.isEquivalent, LiteralValues do not match")
		}
		return false
	}
	return lPtr.element.isEquivalent(hl1, &ref.element, hl2, print)
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
	buffer.WriteString(fmt.Sprintf("\"LiteralValue\":\"%s\",", lPtr.LiteralValue))
	lPtr.element.marshalElementFields(buffer)
	return nil
}

// recoverLiteralFields() is used when de-serializing an element. The activities in restoring the
// literal are not considered changes so the version counter is not incremented and the monitors of this
// element are not notified of chaanges.
func (lPtr *literal) recoverLiteralFields(unmarshaledData *map[string]json.RawMessage, hl *HeldLocks) error {
	err := lPtr.recoverElementFields(unmarshaledData, hl)
	if err != nil {
		return err
	}
	// LiteralValue
	var recoveredLiteralValue string
	err = json.Unmarshal((*unmarshaledData)["LiteralValue"], &recoveredLiteralValue)
	if err != nil {
		log.Printf("Recovery of Element.LiteralValue as string failed\n")
		return err
	}
	lPtr.LiteralValue = recoveredLiteralValue
	return nil
}

func (lPtr *literal) SetLiteralValue(value string, hl *HeldLocks) error {
	hl.WriteLockElement(lPtr)
	editableError := lPtr.editableError(hl)
	if editableError != nil {
		return editableError
	}
	if lPtr.LiteralValue != value {
		lPtr.uOfD.preChange(lPtr, hl)
		notification := lPtr.uOfD.NewConceptChangeNotification(lPtr, hl)
		lPtr.incrementVersion(hl)
		lPtr.LiteralValue = value
		lPtr.uOfD.queueFunctionExecutions(lPtr, notification, hl)
	}
	return nil
}

// Literal is a concept that is, literally, a literal
type Literal interface {
	Element
	GetLiteralValue(*HeldLocks) string
	SetLiteralValue(string, *HeldLocks) error
}
