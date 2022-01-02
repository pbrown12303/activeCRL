package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/pkg/errors"
)

type literal struct {
	element
	LiteralValue string
}

func (lPtr *literal) clone(hl *Transaction) Literal {
	hl.ReadLockElement(lPtr)
	var clonedLiteral literal
	clonedLiteral.initializeLiteral("", "")
	clonedLiteral.cloneAttributes(lPtr, hl)
	return &clonedLiteral
}

func (lPtr *literal) cloneAttributes(source *literal, hl *Transaction) {
	lPtr.element.cloneAttributes(&source.element, hl)
	lPtr.LiteralValue = source.LiteralValue
	lPtr.element.cloneAttributes(&source.element, hl)
}

func (lPtr *literal) GetLiteralValue(hl *Transaction) string {
	hl.ReadLockElement(lPtr)
	return lPtr.LiteralValue
}

func (lPtr *literal) initializeLiteral(conceptID string, uri string) {
	lPtr.initializeElement(conceptID, uri)
}

func (lPtr *literal) isEquivalent(hl1 *Transaction, ref *literal, hl2 *Transaction, printExceptions ...bool) bool {
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
func (lPtr *literal) recoverLiteralFields(unmarshaledData *map[string]json.RawMessage, hl *Transaction) error {
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

func (lPtr *literal) SetLiteralValue(value string, hl *Transaction) error {
	if lPtr.uOfD == nil {
		return errors.New("literal.SetLiteralValue failed because the element uOfD is nil")
	}
	hl.WriteLockElement(lPtr)
	if !lPtr.isEditable(hl) {
		return errors.New("literal.SetLiteralValue failed because the literal is not editable")
	}
	if lPtr.LiteralValue != value {
		lPtr.uOfD.preChange(lPtr, hl)
		beforeState, err := NewConceptState(lPtr)
		if err != nil {
			return errors.Wrap(err, "literal.SetLiteralValue failed")
		}
		lPtr.incrementVersion(hl)
		lPtr.LiteralValue = value
		afterState, err2 := NewConceptState(lPtr)
		if err2 != nil {
			return errors.Wrap(err2, "literal.SetLiteralValue failed")
		}
		err = lPtr.uOfD.SendConceptChangeNotification(lPtr, beforeState, afterState, hl)
		if err != nil {
			return errors.Wrap(err, "literal.SetLiteralValue failed")
		}
	}
	return nil
}

// Literal is a concept that is, literally, a literal
type Literal interface {
	Element
	GetLiteralValue(*Transaction) string
	SetLiteralValue(string, *Transaction) error
}
