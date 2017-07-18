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

func (lPtr *literal) GetLiteralValue() string {
	lPtr.TraceableLock()
	defer lPtr.TraceableUnlock()
	return lPtr.getLiteralValue()
}

func (lPtr *literal) getLiteralValue() string {
	return lPtr.literalValue
}

func (lPtr *literal) GetNameNoLock() string {
	lPtr.TraceableLock()
	defer lPtr.TraceableUnlock()
	return lPtr.getLiteralValue()
}

func (lPtr *literal) initializeLiteral() {
	lPtr.initializeValue()
}

func (lPtr *literal) isEquivalent(lit *literal) bool {
	if lPtr.literalValue != lit.literalValue {
		log.Printf("Literal values not equivalent - v1: %s v2: %s \n", lPtr.literalValue, lit.literalValue)
		return false
	}
	var valuePtr *value = &lPtr.value
	return valuePtr.isEquivalent(&lit.value)
}

func (lPtr *literal) MarshalJSON() ([]byte, error) {
	lPtr.TraceableLock()
	defer lPtr.TraceableUnlock()
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

func (lPtr *literal) printLiteral(prefix string) {
	lPtr.printValue(prefix)
	log.Printf("%sliteralValue: %s \n", prefix, lPtr.literalValue)
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

func (lPtr *literal) SetLiteralValue(newValue string) {
	lPtr.TraceableLock()
	defer lPtr.TraceableUnlock()
	lPtr.setLiteralValue(newValue)
}

func (lPtr *literal) setLiteralValue(newValue string) {
	if lPtr.literalValue != newValue {
		preChange(lPtr)
		lPtr.literalValue = newValue
		notification := NewChangeNotification(lPtr, MODIFY, nil)
		postChange(lPtr, notification)
	}
}

func (lPtr *literal) SetOwningElement(el Element) {
	lPtr.TraceableLock()
	defer lPtr.TraceableUnlock()
	lPtr.SetOwningElementNoLock(el)
}

func (lPtr *literal) SetOwningElementNoLock(el Element) {
	if lPtr.getOwningElement() != el {
		if lPtr.owningElement != nil {
			lPtr.owningElement.removeOwnedBaseElement(lPtr)
		}

		preChange(lPtr)
		lPtr.owningElement = el
		notification := NewChangeNotification(lPtr, MODIFY, nil)
		postChange(lPtr, notification)

		if lPtr.owningElement != nil {
			lPtr.owningElement.addOwnedBaseElement(lPtr)
		}
	}
}

// internalSetOwningElement() is an internal function used only in unmarshal
func (lPtr *literal) internalSetOwningElement(el Element) {
	lPtr.owningElement = el
	if lPtr.owningElement != nil {
		lPtr.owningElement.internalAddOwnedBaseElement(lPtr)
	}
}

func (lPtr *literal) SetUri(uri string) {
	lPtr.TraceableLock()
	defer lPtr.TraceableUnlock()
	lPtr.SetUriNoLock(uri)
}

func (lPtr *literal) SetUriNoLock(uri string) {
	preChange(lPtr)
	lPtr.uri = uri
	notification := NewChangeNotification(lPtr, MODIFY, nil)
	postChange(lPtr, notification)
}

type Literal interface {
	Value
	GetLiteralValue() string
	getLiteralValue() string
	setLiteralValue(string)
	SetLiteralValue(string)
}
