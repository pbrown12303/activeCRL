// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"reflect"
)

type element struct {
	baseElement
	ownedBaseElements map[string]BaseElement
}

// addOwnedBaseElement() adds the indicated base element as a child (owned)
// base element of this object. Calling this method is considered a change to the element
// and will result in monitors being notified of changes.
func addOwnedBaseElement(elPtr Element, be BaseElement, hl *HeldLocks) {
	preChange(elPtr, hl)
	elPtr.internalAddOwnedBaseElement(be, hl)
	notification := NewChangeNotification(elPtr, ADD, nil)
	postChange(elPtr, notification, hl)
}

// childChanged() is used by ownedBaseElements to inform their parents when they have changed. It does no locking.
func childChanged(el Element, notification *ChangeNotification, hl *HeldLocks) {
	if TraceChange == true {
		log.Printf("childChanged called on Element %s \n", el.GetId(hl).String())
		notification.Print("ChildChanged Incoming Notification: ", hl)
	}
	preChange(el, hl)
	newNotification := NewChangeNotification(el, MODIFY, notification)
	switch el.(type) {
	case Refinement:
		refinedElement := el.(Refinement).GetRefinedElement(hl)
		refinedElementPointer := el.(Refinement).GetRefinedElementPointer(hl)
		if refinedElement != nil {
			cn := notification.getReferencingChangeNotification(refinedElementPointer)
			if cn != nil && cn.underlyingChange == nil {
				abstractionChanged(refinedElement, newNotification, hl)
			}
		}
	}
	postChange(el, newNotification, hl)
}

func (elPtr *element) clone() *element {
	var cl element
	cl.ownedBaseElements = make(map[string]BaseElement)
	cl.cloneAttributes(*elPtr)
	return &cl
}

func (elPtr *element) cloneAttributes(source element) {
	elPtr.baseElement.cloneAttributes(source.baseElement)
	for key, _ := range elPtr.ownedBaseElements {
		delete(elPtr.ownedBaseElements, key)
	}
	for key, value := range source.ownedBaseElements {
		elPtr.ownedBaseElements[key] = value
	}
}

// GetAbstractElementsRecursivelyNoLock() returns all of the elements abstractions
func (elPtr *element) GetAbstractElementsRecursively(hl *HeldLocks) []Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	abstractElements := elPtr.getImmediateAbstractElements(hl)
	var ancestors []Element
	for _, element := range abstractElements {
		for _, ancestor := range element.GetAbstractElementsRecursively(hl) {
			ancestors = append(ancestors, ancestor)
		}
	}
	for _, ancestor := range ancestors {
		abstractElements = append(abstractElements, ancestor)
	}
	return abstractElements
}

func (elPtr *element) GetDefinition(hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	nlp := elPtr.GetDefinitionLiteralPointer(hl)
	if nlp != nil {
		nl := nlp.GetLiteral(hl)
		if nl != nil {
			return nl.GetLiteralValue(hl)
		}
	}
	return ""
}

func (elPtr *element) GetDefinitionLiteral(hl *HeldLocks) Literal {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	nlp := elPtr.GetDefinitionLiteralPointer(hl)
	if nlp != nil {
		return nlp.GetLiteral(hl)
	}
	return nil
}

func (elPtr *element) GetDefinitionLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	for _, be := range elPtr.ownedBaseElements {
		switch be.(type) {
		case LiteralPointer:
			if be.(LiteralPointer).GetLiteralPointerRole(hl) == DEFINITION {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

func (elPtr *element) getImmediateAbstractElements(hl *HeldLocks) []Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	var abstractElements []Element
	abstractions := elPtr.getImmediateAbstractions(hl)
	if abstractions != nil {
		for _, abstraction := range abstractions {
			if abstraction.GetAbstractElement(hl) != nil {
				abstractElements = append(abstractElements, abstraction.GetAbstractElement(hl))
			}
		}
	}
	return abstractElements
}

func (elPtr *element) getImmediateAbstractions(hl *HeldLocks) []Refinement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	var abstractions []Refinement
	ePtrs := elPtr.uOfD.elementListenerMap.GetEntry(elPtr.GetId(hl).String())
	if ePtrs != nil {
		for _, ePtr := range *ePtrs {
			if ePtr.GetElementPointerRole(hl) == REFINED_ELEMENT {
				abstractions = append(abstractions, GetOwningElement(ePtr, hl).(Refinement))
			}
		}
	}
	return abstractions
}

func (elPtr *element) getImmediateRefinements(hl *HeldLocks) []Refinement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	var refinements []Refinement
	ePtrs := elPtr.uOfD.elementListenerMap.GetEntry(elPtr.GetId(hl).String())
	if ePtrs != nil {
		for _, ePtr := range *ePtrs {
			if ePtr.GetElementPointerRole(hl) == ABSTRACT_ELEMENT {
				refinements = append(refinements, GetOwningElement(ePtr, hl).(Refinement))
			}
		}
	}
	return refinements
}

func (elPtr *element) GetNameLiteral(hl *HeldLocks) Literal {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	nlp := elPtr.GetNameLiteralPointer(hl)
	if nlp != nil {
		return nlp.GetLiteral(hl)
	}
	return nil
}

func (elPtr *element) GetNameLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	for _, be := range elPtr.GetOwnedBaseElements(hl) {
		switch be.(type) {
		case LiteralPointer:
			if be.(LiteralPointer).GetLiteralPointerRole(hl) == NAME {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

func (elPtr *element) GetOwnedBaseElements(hl *HeldLocks) []BaseElement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	var obe []BaseElement
	for _, be := range elPtr.ownedBaseElements {
		obe = append(obe, be)
	}
	return obe
}

func (elPtr *element) GetOwnedElements(hl *HeldLocks) []Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	var obe []Element
	for _, be := range elPtr.ownedBaseElements {
		switch be.(type) {
		case Element:
			obe = append(obe, be.(Element))
		}
	}
	return obe
}

//func (elPtr *element) GetOwningElement(hl *HeldLocks) Element {
//	if hl == nil {
//		hl = NewHeldLocks(nil)
//		defer hl.ReleaseLocks()
//	}
//	hl.LockBaseElement(elPtr)
//	oep := elPtr.GetOwningElementPointer(hl)
//	if oep != nil {
//		return oep.GetElement(hl)
//	}
//	return nil
//}

func (elPtr *element) GetOwningElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	for _, be := range elPtr.ownedBaseElements {
		switch be.(type) {
		case *elementPointer:
			if be.(ElementPointer).GetElementPointerRole(hl) == OWNING_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

func (elPtr *element) GetUri(hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	ul := elPtr.GetUriLiteral(hl)
	if ul != nil {
		return ul.GetLiteralValue(hl)
	}
	return ""
}

func (elPtr *element) GetUriLiteral(hl *HeldLocks) Literal {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	nlp := elPtr.GetUriLiteralPointer(hl)
	if nlp != nil {
		return nlp.GetLiteral(hl)
	}
	return nil
}

func (elPtr *element) GetUriLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	for _, be := range elPtr.ownedBaseElements {
		switch be.(type) {
		case LiteralPointer:
			if be.(LiteralPointer).GetLiteralPointerRole(hl) == URI {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

// initializeElement() creates the ownedBaseElements map and calls initializeBaseElement().
// Note that initialization is not considered a change, so the version counter is not incremented
// nor are monitors of this element notified of changes.
func (elPtr *element) initializeElement() {
	elPtr.initializeBaseElement()
	elPtr.ownedBaseElements = make(map[string]BaseElement)
}

// internalAddOwnedBaseElement() adds the indicated base element as a child (owned)
// base element of this object. Calling this method is not considered a change to the element
// and will not result in monitors being notified of changes.
func (elPtr *element) internalAddOwnedBaseElement(be BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	if be != nil && be.GetId(hl) != uuid.Nil {
		elPtr.ownedBaseElements[be.GetId(hl).String()] = be
	}
}

// internalRemoveOwnedBaseElement() removes the indicated baseElement from the ownedBaseElements
// map. Note that this is not considered a change and that the version counter will not be incremented and
// the monitors of this element will not be notified of the change.
func (elPtr *element) internalRemoveOwnedBaseElement(be BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	if be != nil && be.GetId(hl) != uuid.Nil {
		delete(elPtr.ownedBaseElements, be.GetId(hl).String())
	}
}

// isEquivalent is a non-locking function that compares this element against another to see
// if the other element and its substructure are equivalent
func (bePtr *element) isEquivalent(be *element, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	if len(bePtr.ownedBaseElements) != len(be.ownedBaseElements) {
		//		log.Printf("Equivalence failed: Owned Base Elements lenght does not match \n")
		return false
	}
	for key, value := range bePtr.ownedBaseElements {
		beValue := be.ownedBaseElements[key]
		if beValue == nil {
			//			log.Printf("Equivalence failed: no value found for Owned Base Element key %s \n", key)
			return false
		}
		if !Equivalent(value, beValue, hl) {
			//			log.Printf("Equivalence failed: values do not match for Owned Base Element key %s \n", key)
			//			log.Printf("First element's value: \n")
			//			Print(value, "   ")
			//			log.Printf("Second element's value: \n")
			//			Print(beValue, "   ")
			return false
		}
	}
	var baseElementPtr *baseElement = &bePtr.baseElement
	return baseElementPtr.isEquivalent(&be.baseElement, hl)
}

func (ePtr *element) IsOwnedBaseElement(be BaseElement, hl *HeldLocks) bool {
	for key, _ := range ePtr.ownedBaseElements {
		if key == be.GetId(hl).String() {
			return true
		}
	}
	return false
}

func (elPtr *element) IsRefinementOf(el Element, hl *HeldLocks) bool {
	for _, abstractElement := range elPtr.GetAbstractElementsRecursively(hl) {
		if el == abstractElement {
			return true
		}
	}
	return false
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

func (elPtr *element) printElement(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elPtr)
	// We use the prefix lenth to curtail infinite recursion - circular ownership
	if len(prefix) > 300 {
		log.Printf("Prefix length exceeds 300")
		return
	}
	elPtr.printBaseElement(prefix, hl)
	log.Printf("%s  Owned Base Elements: count %d \n", prefix, len(elPtr.ownedBaseElements))
	extendedPrefix := prefix + "   "
	for _, be := range elPtr.ownedBaseElements {
		Print(be, extendedPrefix, hl)
	}
}

// recoverElementFields() is used when de-serializing an element. The activities in restoring the
// element are not considered changes so the version counter is not incremented and the monitors of this
// element are not notified of chaanges.
func (el *element) recoverElementFields(unmarshaledData *map[string]json.RawMessage) error {
	err := el.baseElement.recoverBaseElementFields(unmarshaledData)
	if err != nil {
		return err
	}
	var obeMap map[string]json.RawMessage
	err = json.Unmarshal((*unmarshaledData)["OwnedBaseElements"], &obeMap)
	if err != nil {
		log.Printf("Recovery of Element.OwnedBaseElements failed\n")
		return err
	}
	for _, rawBe := range obeMap {
		var recoveredBaseElement BaseElement
		err = unmarshalPolymorphicBaseElement(rawBe, &recoveredBaseElement)
		if err != nil {
			log.Printf("Polymorphic Recovery of one Element.OwnedBaseElements failed\n")
			return err
		}
		el.internalAddOwnedBaseElement(recoveredBaseElement, nil)
	}
	return nil
}

// removeOwnedBaseElement() removes the indicated baseElement from the ownedBaseElements
// map. Note that this is considered a change and that the version counter will be incremented and
// the monitors of this element will be notified of the change.
func removeOwnedBaseElement(elPtr Element, be BaseElement, hl *HeldLocks) {
	preChange(elPtr, hl)
	elPtr.internalRemoveOwnedBaseElement(be, hl)
	notification := NewChangeNotification(elPtr, REMOVE, nil)
	postChange(elPtr, notification, hl)
}

type Element interface {
	BaseElement
	GetAbstractElementsRecursively(*HeldLocks) []Element
	GetDefinition(*HeldLocks) string
	GetDefinitionLiteral(*HeldLocks) Literal
	GetDefinitionLiteralPointer(*HeldLocks) LiteralPointer
	getImmediateAbstractElements(*HeldLocks) []Element
	getImmediateAbstractions(*HeldLocks) []Refinement
	getImmediateRefinements(*HeldLocks) []Refinement
	GetNameLiteral(*HeldLocks) Literal
	GetNameLiteralPointer(*HeldLocks) LiteralPointer
	GetOwnedBaseElements(*HeldLocks) []BaseElement
	GetOwnedElements(*HeldLocks) []Element
	//	GetOwningElement(*HeldLocks) Element
	GetOwningElementPointer(*HeldLocks) ElementPointer
	GetUriLiteral(*HeldLocks) Literal
	GetUriLiteralPointer(*HeldLocks) LiteralPointer
	internalAddOwnedBaseElement(BaseElement, *HeldLocks)
	internalRemoveOwnedBaseElement(BaseElement, *HeldLocks)
	IsOwnedBaseElement(BaseElement, *HeldLocks) bool
	IsRefinementOf(Element, *HeldLocks) bool
	MarshalJSON() ([]byte, error)
}

func SetDefinition(el Element, definition string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(el)
	nl := el.GetDefinitionLiteral(hl)
	if nl == nil {
		nlp := el.GetDefinitionLiteralPointer(hl)
		if nlp == nil {
			nlp = el.GetUniverseOfDiscourse(hl).NewDefinitionLiteralPointer(hl)
			SetOwningElement(nlp, el, hl)
		}
		nl = el.GetUniverseOfDiscourse(hl).NewLiteral(hl)
		SetOwningElement(nl, el, hl)
		nlp.SetLiteral(nl, hl)
	}
	nl.SetLiteralValue(definition, hl)
}

func SetName(el Element, name string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(el)
	nl := el.GetNameLiteral(hl)
	if nl == nil {
		nlp := el.GetNameLiteralPointer(hl)
		if nlp == nil {
			nlp = el.GetUniverseOfDiscourse(hl).NewNameLiteralPointer(hl)
			SetOwningElement(nlp, el, hl)
		}
		nl = el.GetUniverseOfDiscourse(hl).NewLiteral(hl)
		SetOwningElement(nl, el, hl)
		nlp.SetLiteral(nl, hl)
	}
	nl.SetLiteralValue(name, hl)
}
