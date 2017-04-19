package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"reflect"
)

type element struct {
	baseElement
	ownedBaseElements map[string]BaseElement
}

func NewElement(uOfD *UniverseOfDiscourse) Element {
	var el element
	el.initializeElement()
	uOfD.addBaseElement(&el)
	return &el
}

func (elPtr *element) addOwnedBaseElement(be BaseElement) {
	if be != nil && be.GetId() != uuid.Nil {
		elPtr.ownedBaseElements[be.GetId().String()] = be
	}
}

func (elPtr *element) GetDefinition() string {
	nl := elPtr.getDefinitionLiteral()
	if nl != nil {
		return nl.GetLiteralValue()
	}
	return ""
}

func (elPtr *element) getDefinitionLiteral() Literal {
	nlp := elPtr.getDefinitionLiteralPointer()
	if nlp != nil {
		return nlp.GetLiteral()
	}
	return nil
}

func (elPtr *element) getDefinitionLiteralPointer() LiteralPointer {
	for _, be := range elPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *literalPointer:
			if be.(*literalPointer).getLiteralPointerRole() == DEFINITION {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

func (elPtr *element) GetName() string {
	nl := elPtr.getNameLiteral()
	if nl != nil {
		return nl.GetLiteralValue()
	}
	return ""
}

func (elPtr *element) getNameLiteral() Literal {
	nlp := elPtr.getNameLiteralPointer()
	if nlp != nil {
		return nlp.GetLiteral()
	}
	return nil
}

func (elPtr *element) getNameLiteralPointer() LiteralPointer {
	for _, be := range elPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *literalPointer:
			if be.(*literalPointer).getLiteralPointerRole() == NAME {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

func (elPtr *element) GetOwnedBaseElements() map[string]BaseElement {
	return elPtr.ownedBaseElements
}

func (elPtr *element) GetOwningElement() Element {
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		return oep.GetElement()
	}
	return nil
}

func (elPtr *element) getOwningElementPointer() ElementPointer {
	for _, be := range elPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *elementPointer:
			if be.(*elementPointer).getElementPointerRole() == OWNING_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

func (elPtr *element) GetUri() string {
	nl := elPtr.getUriLiteral()
	if nl != nil {
		return nl.GetLiteralValue()
	}
	return ""
}

func (elPtr *element) getUriLiteral() Literal {
	nlp := elPtr.getUriLiteralPointer()
	if nlp != nil {
		return nlp.GetLiteral()
	}
	return nil
}

func (elPtr *element) getUriLiteralPointer() LiteralPointer {
	for _, be := range elPtr.GetOwnedBaseElements() {
		switch be.(type) {
		case *literalPointer:
			if be.(*literalPointer).getLiteralPointerRole() == URI {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

func (elPtr *element) initializeElement() {
	elPtr.initializeBaseElement()
	elPtr.ownedBaseElements = make(map[string]BaseElement)
}

func (bePtr *element) isEquivalent(be *element) bool {
	if len(bePtr.ownedBaseElements) != len(be.ownedBaseElements) {
		fmt.Printf("Equivalence failed: Owned Base Elements lenght does not match \n")
		return false
	}
	for key, value := range bePtr.ownedBaseElements {
		beValue := be.ownedBaseElements[key]
		if beValue == nil {
			fmt.Printf("Equivalence failed: no value found for Owned Base Element key %s \n", key)
			return false
		}
		if !Equivalent(value, beValue) {
			fmt.Printf("Equivalence failed: values do not match for Owned Base Element key %s \n", key)
			fmt.Printf("First element's value: \n")
			Print(value, "   ")
			fmt.Printf("Second element's value: \n")
			Print(beValue, "   ")
			return false
		}
	}
	var baseElementPtr *baseElement = &bePtr.baseElement
	return baseElementPtr.isEquivalent(&be.baseElement)
}

func (elPtr *element) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalElementFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *element) marshalElementFields(buffer *bytes.Buffer) error {
	elPtr.baseElement.marshalBaseElementFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"OwnedBaseElements\":{"))
	count := len(elPtr.ownedBaseElements)
	for key, value := range elPtr.ownedBaseElements {
		count--
		buffer.WriteString(fmt.Sprintf("\"%s\":", key))
		encodedObject, err := json.Marshal(value)
		if err != nil {
			return err
		}
		buffer.Write(encodedObject)
		if count > 0 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(fmt.Sprintf("}"))
	return nil
}

var printCount int = 0

func (elPtr *element) printElement(prefix string) {
	if printCount < 100 {
		printCount++
		elPtr.printBaseElement(prefix)
		fmt.Printf("%sOwned Base Elements: count %d \n", prefix, len(elPtr.GetOwnedBaseElements()))
		extendedPrefix := prefix + "   "
		for _, be := range elPtr.GetOwnedBaseElements() {
			Print(be, extendedPrefix)
		}
	}
}

func (el *element) recoverElementFields(unmarshaledData *map[string]json.RawMessage) error {
	err := el.baseElement.recoverBaseElementFields(unmarshaledData)
	if err != nil {
		return err
	}
	var obeMap map[string]json.RawMessage
	err = json.Unmarshal((*unmarshaledData)["OwnedBaseElements"], &obeMap)
	if err != nil {
		fmt.Printf("Recovery of Element.OwnedBaseElements failed\n")
		return err
	}
	for _, rawBe := range obeMap {
		var recoveredBaseElement BaseElement
		err = unmarshalPolymorphicBaseElement(rawBe, &recoveredBaseElement)
		if err != nil {
			fmt.Printf("Polymorphic Recovery of one Element.OwnedBaseElements failed\n")
			return err
		}
		el.addOwnedBaseElement(recoveredBaseElement)
	}
	return nil
}

func (elPtr *element) removeOwnedBaseElement(be BaseElement) {
	if be != nil && be.GetId() != uuid.Nil {
		delete(elPtr.ownedBaseElements, be.GetId().String())
	}
}

func (elPtr *element) SetDefinition(name string) {
	nl := elPtr.getDefinitionLiteral()
	if nl == nil {
		nlp := elPtr.getDefinitionLiteralPointer()
		if nlp == nil {
			nlp = NewDefinitionLiteralPointer(elPtr.getUniverseOfDiscourse())
			nlp.setOwningElement(elPtr)
		}
		nl = NewLiteral(elPtr.getUniverseOfDiscourse())
		nl.setOwningElement(elPtr)
		nlp.SetLiteral(nl)
	}
	nl.SetLiteralValue(name)
}

func (elPtr *element) SetName(name string) {
	nl := elPtr.getNameLiteral()
	if nl == nil {
		nlp := elPtr.getNameLiteralPointer()
		if nlp == nil {
			nlp = NewNameLiteralPointer(elPtr.getUniverseOfDiscourse())
			nlp.setOwningElement(elPtr)
		}
		nl = NewLiteral(elPtr.getUniverseOfDiscourse())
		nl.setOwningElement(elPtr)
		nlp.SetLiteral(nl)
	}
	nl.SetLiteralValue(name)
}

func (elPtr *element) setOwningElement(owningElement Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = NewOwningElementPointer(elPtr.uOfD)
		oep.setOwningElement(elPtr)
	}
	oep.SetElement(owningElement)
}

func (elPtr *element) SetUri(name string) {
	nl := elPtr.getUriLiteral()
	if nl == nil {
		nlp := elPtr.getUriLiteralPointer()
		if nlp == nil {
			nlp = NewUriLiteralPointer(elPtr.getUniverseOfDiscourse())
			nlp.setOwningElement(elPtr)
		}
		nl = NewLiteral(elPtr.getUniverseOfDiscourse())
		nl.setOwningElement(elPtr)
		nlp.SetLiteral(nl)
	}
	nl.SetLiteralValue(name)
}

type Element interface {
	BaseElement
	addOwnedBaseElement(BaseElement)
	GetDefinition() string
	getDefinitionLiteral() Literal
	getDefinitionLiteralPointer() LiteralPointer
	getNameLiteral() Literal
	getNameLiteralPointer() LiteralPointer
	GetOwnedBaseElements() map[string]BaseElement
	getOwningElementPointer() ElementPointer
	GetUri() string
	getUriLiteral() Literal
	getUriLiteralPointer() LiteralPointer
	removeOwnedBaseElement(BaseElement)
	SetDefinition(string)
	SetName(string)
	SetUri(string)
}
