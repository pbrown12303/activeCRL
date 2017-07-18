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
func (elPtr *element) addOwnedBaseElement(be BaseElement) {
	preChange(elPtr)
	elPtr.internalAddOwnedBaseElement(be)
	notification := NewChangeNotification(elPtr, ADD, nil)
	postChange(elPtr, notification)
}

// childChanged() is used by ownedBaseElements to inform their parents when they have changed. It does no locking.
func (elPtr *element) childChanged(notification *ChangeNotification) {
	preChange(elPtr)
	newNotification := NewChangeNotification(elPtr, MODIFY, notification)
	postChange(elPtr, newNotification)
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

// GetAbstractElementsRecursivelyNoLock() is a non-locking element that returns all of the elements abstractions
func (elPtr *element) GetAbstractElementsRecursivelyNoLock() []Element {
	abstractElements := elPtr.getImmediateAbstractElements()
	var ancestors []Element
	for _, element := range abstractElements {
		for _, ancestor := range element.GetAbstractElementsRecursivelyNoLock() {
			ancestors = append(ancestors, ancestor)
		}
	}
	for _, ancestor := range ancestors {
		abstractElements = append(abstractElements, ancestor)
	}
	return abstractElements
}

// GetDefinition() is the public method to retrieve the definition. It locks the element and, if present, the definitionLiteralPointer
// and the literal. It then calls the non-locking getDefinition()
func (elPtr *element) GetDefinition() string {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	nlp := elPtr.getDefinitionLiteralPointer()
	if nlp != nil {
		nlp.TraceableLock()
		defer nlp.TraceableUnlock()
	}
	nl := elPtr.getDefinitionLiteral()
	if nl != nil {
		nl.TraceableLock()
		defer nl.TraceableUnlock()
	}
	return elPtr.getDefinition()
}

// getDefinition() is an internal method that actually gets the name. If there is a definitionLiteralPointer that
// points to a literal, it returns the value of the literal. Otherwise it returns the empty string. This method does
// no locking.
func (elPtr *element) getDefinition() string {
	nlp := elPtr.getDefinitionLiteralPointer()
	if nlp != nil {
		nl := nlp.getLiteral()
		if nl != nil {
			return nl.getLiteralValue()
		}
	}
	return ""
}

func (elPtr *element) getDefinitionLiteral() Literal {
	nlp := elPtr.getDefinitionLiteralPointer()
	if nlp != nil {
		return nlp.getLiteral()
	}
	return nil
}

func (elPtr *element) getDefinitionLiteralPointer() LiteralPointer {
	for _, be := range elPtr.getOwnedBaseElements() {
		switch be.(type) {
		case *literalPointer:
			if be.(*literalPointer).getLiteralPointerRole() == DEFINITION {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

func (elPtr *element) getImmediateAbstractElements() []Element {
	var abstractElements []Element
	abstractions := elPtr.getImmediateAbstractions()
	if abstractions != nil {
		for _, abstraction := range abstractions {
			if abstraction.getAbstractElement() != nil {
				abstractElements = append(abstractElements, abstraction.getAbstractElement())
			}
		}
	}
	return abstractElements
}

func (elPtr *element) getImmediateAbstractions() []Refinement {
	var abstractions []Refinement
	ePtrs := elPtr.uOfD.elementListenerMap[elPtr.getId().String()]
	if ePtrs != nil {
		for _, ePtr := range *ePtrs {
			if ePtr.getElementPointerRole() == REFINED_ELEMENT {
				abstractions = append(abstractions, ePtr.getOwningElement().(Refinement))
			}
		}
	}
	return abstractions
}

func (elPtr *element) getImmediateRefinements() []Refinement {
	var refinements []Refinement
	ePtrs := elPtr.uOfD.elementListenerMap[elPtr.getId().String()]
	if ePtrs != nil {
		for _, ePtr := range *ePtrs {
			if ePtr.getElementPointerRole() == ABSTRACT_ELEMENT {
				refinements = append(refinements, ePtr.getOwningElement().(Refinement))
			}
		}
	}
	return refinements
}

// GetName() locks the element and, if they are not nil, the nameLiteralPointer and name literal. It then
// returns the result of calling the non-locking GetNameNoLock()
func (elPtr *element) GetName() string {
	elPtr.TraceableRLock()
	defer elPtr.TraceableRUnlock()
	nlp := elPtr.getNameLiteralPointer()
	if nlp != nil {
		nlp.TraceableLock()
		defer nlp.TraceableUnlock()
	}
	nl := elPtr.getNameLiteral()
	if nl != nil {
		nl.TraceableLock()
		defer nl.TraceableUnlock()
	}
	return elPtr.GetNameNoLock()
}

// GetNameNoLock() is a non-locking function that returns the name string.
func (elPtr *element) GetNameNoLock() string {
	nl := elPtr.getNameLiteral()
	if nl != nil {
		return nl.getLiteralValue()
	}
	return ""
}

// getNameLiteral() is a non-locking function that returns the name literal or nil if there is none.
func (elPtr *element) getNameLiteral() Literal {
	nlp := elPtr.getNameLiteralPointer()
	if nlp != nil {
		return nlp.getLiteral()
	}
	return nil
}

// getNameLiteralPointer() is a non-locking function that walks the ownedBaseElements set and returns
// the first member that is a literalPointer with the role set to NAME
func (elPtr *element) getNameLiteralPointer() LiteralPointer {
	for _, be := range elPtr.getOwnedBaseElements() {
		switch be.(type) {
		case *literalPointer:
			if be.(*literalPointer).getLiteralPointerRole() == NAME {
				return be.(LiteralPointer)
			}
		}
	}
	return nil
}

// getOwnedBaseElements() is a non-locking function that returns the element's ownedBaseElements map.
func (elPtr *element) getOwnedBaseElements() map[string]BaseElement {
	return elPtr.ownedBaseElements
}

// GetOwnedBaseElements() is a non-locking function that returns an array containing the ownedBaseElements
func (elPtr *element) GetOwnedBaseElements() []BaseElement {
	var obe []BaseElement
	for _, be := range elPtr.ownedBaseElements {
		obe = append(obe, be)
	}
	return obe
}

// GetOwnedElementsNoLock() is a non-locking function that returns an array containing the
// Elements belonging to ownedBaseElements
func (elPtr *element) GetOwnedElementsNoLock() []Element {
	var obe []Element
	for _, be := range elPtr.ownedBaseElements {
		switch be.(type) {
		case Element:
			obe = append(obe, be.(Element))
		}
	}
	return obe
}

// GetOwningElement is a locking function that locks the element and, if present, the owningElementPointer. It then
// returns the value of the non-locking getOwningElement()
func (elPtr *element) GetOwningElement() Element {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		oep.TraceableLock()
		defer oep.TraceableUnlock()
	}
	return elPtr.getOwningElement()
}

// getOwningElement() is a non-locking function that uses the owningElementPointer to locate the owningElement and return it.
func (elPtr *element) getOwningElement() Element {
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		return oep.getElement()
	}
	return nil
}

// getOwningElementPointer() is a non-locking function that walks the ownedBaseElements and returns the first
// elementPointer whose role is set to OWNING_ELEMENT
func (elPtr *element) getOwningElementPointer() ElementPointer {
	for _, be := range elPtr.getOwnedBaseElements() {
		switch be.(type) {
		case *elementPointer:
			if be.(*elementPointer).getElementPointerRole() == OWNING_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

// GetUniverseOfDiscourse() is a locking function that returns a pointer to the UniverseOfDiscourse to which this element belongs
func (elPtr *element) GetUniverseOfDiscourse() *UniverseOfDiscourse {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	return elPtr.GetUniverseOfDiscourseNoLock()
}

// GetUniverseOfDiscourseNoLock() is a non-locking function that returns a pointer to the UniverseOfDiscourse to which this element belongs
func (elPtr *element) GetUniverseOfDiscourseNoLock() *UniverseOfDiscourse {
	return elPtr.uOfD
}

// GetUri() is a locking function that locks the element and, if present, the uriLiteralPointer and uriLiteral. It then
// returns the result of the non-locking GetUriNoLock()
func (elPtr *element) GetUri() string {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	ulp := elPtr.getUriLiteralPointer()
	if ulp != nil {
		ulp.TraceableLock()
		defer ulp.TraceableUnlock()
	}
	ul := elPtr.getUriLiteral()
	if ul != nil {
		ul.TraceableLock()
		ul.TraceableUnlock()
	}
	return elPtr.GetUriNoLock()
}

// GetUriNoLock() is a non-locking function that uses the uriLiteralPointer to locate the uriLiteral and return its string value.
func (elPtr *element) GetUriNoLock() string {
	ul := elPtr.getUriLiteral()
	if ul != nil {
		return ul.GetLiteralValue()
	}
	return ""
}

// getUriLiteral() is a non-locking function that uses the uriLiteralPointer to locate and return the uriLiteral.
func (elPtr *element) getUriLiteral() Literal {
	nlp := elPtr.getUriLiteralPointer()
	if nlp != nil {
		return nlp.getLiteral()
	}
	return nil
}

// getUriLiteralPointer() is a non-locking function that walks the element's ownedBaseElements and returns the first
// literalPointer whose role is URI
func (elPtr *element) getUriLiteralPointer() LiteralPointer {
	for _, be := range elPtr.getOwnedBaseElements() {
		switch be.(type) {
		case *literalPointer:
			if be.(*literalPointer).getLiteralPointerRole() == URI {
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
func (elPtr *element) internalAddOwnedBaseElement(be BaseElement) {
	if be != nil && be.getId() != uuid.Nil {
		elPtr.ownedBaseElements[be.getId().String()] = be
	}
}

// internalRemoveOwnedBaseElement() removes the indicated baseElement from the ownedBaseElements
// map. Note that this is not considered a change and that the version counter will not be incremented and
// the monitors of this element will not be notified of the change.
func (elPtr *element) internalRemoveOwnedBaseElement(be BaseElement) {
	if be != nil && be.getId() != uuid.Nil {
		delete(elPtr.ownedBaseElements, be.getId().String())
	}
}

// isEquivalent is a non-locking function that compares this element against another to see
// if the other element and its substructure are equivalent
func (bePtr *element) isEquivalent(be *element) bool {
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
		if !Equivalent(value, beValue) {
			//			log.Printf("Equivalence failed: values do not match for Owned Base Element key %s \n", key)
			//			log.Printf("First element's value: \n")
			//			Print(value, "   ")
			//			log.Printf("Second element's value: \n")
			//			Print(beValue, "   ")
			return false
		}
	}
	var baseElementPtr *baseElement = &bePtr.baseElement
	return baseElementPtr.isEquivalent(&be.baseElement)
}

func (elPtr *element) MarshalJSON() ([]byte, error) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
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
		log.Printf("%sOwned Base Elements: count %d \n", prefix, len(elPtr.getOwnedBaseElements()))
		extendedPrefix := prefix + "   "
		for _, be := range elPtr.getOwnedBaseElements() {
			Print(be, extendedPrefix)
		}
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
		el.internalAddOwnedBaseElement(recoveredBaseElement)
	}
	return nil
}

// removeOwnedBaseElement() removes the indicated baseElement from the ownedBaseElements
// map. Note that this is considered a change and that the version counter will be incremented and
// the monitors of this element will be notified of the change.
func (elPtr *element) removeOwnedBaseElement(be BaseElement) {
	preChange(elPtr)
	elPtr.internalRemoveOwnedBaseElement(be)
	notification := NewChangeNotification(elPtr, REMOVE, nil)
	postChange(elPtr, notification)
}

// SetDefinition() updates the literal containing the definition. If needed both the
// literal and the literalPointer pointing to it are created. This method locks the element,
// and, indirectly, increments the version and notifies monitors of the change.
//
func (elPtr *element) SetDefinition(definition string) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	nl := elPtr.getDefinitionLiteral()
	if nl != nil {
		nl.TraceableLock()
		defer nl.TraceableUnlock()
	}
	nlp := elPtr.getDefinitionLiteralPointer()
	if nlp != nil {
		nlp.TraceableLock()
		defer nlp.TraceableUnlock()
	}
	elPtr.setDefinition(definition)
}

// setDefinition() updates the literal containing the definition. If needed, both the
// literal and the literalPointer pointing to it are created. This method does not lock the element.
// It does not directly increment the version and notify monitors of the change. It is making changes to
// subordinate objects, i.e. the definition literal and definition literal pointer. These objects will, in turn
// notify the element (their parent) of the change.
func (elPtr *element) setDefinition(definition string) {
	nl := elPtr.getDefinitionLiteral()
	if nl == nil {
		nlp := elPtr.getDefinitionLiteralPointer()
		if nlp == nil {
			nlp = elPtr.getUniverseOfDiscourse().NewDefinitionLiteralPointer()
			nlp.SetOwningElementNoLock(elPtr)
		}
		nl = elPtr.getUniverseOfDiscourse().NewLiteral()
		nl.SetOwningElementNoLock(elPtr)
		nlp.setLiteral(nl)
	}
	nl.setLiteralValue(definition)
}

// SetName() is a locking function that locks the element and, if present, the nameLiteralPointer and nameLiteral. It
// then calls the non-locking SetNameNoLock() to actually set the name value.
func (elPtr *element) SetName(name string) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	nl := elPtr.getNameLiteral()
	if nl != nil {
		nl.TraceableLock()
		defer nl.TraceableUnlock()
	}
	nlp := elPtr.getNameLiteralPointer()
	if nlp != nil {
		nlp.TraceableLock()
		defer nlp.TraceableUnlock()
	}
	elPtr.SetNameNoLock(name)
}

// SetNameNoLock() is a non-locking function that checks for the existence of the nameLiteralPointer and nameLiteral, creating
// them if necessary. It then sets the value of the nameLiteral to the indicated string.
func (elPtr *element) SetNameNoLock(name string) {
	nl := elPtr.getNameLiteral()
	if nl == nil {
		nlp := elPtr.getNameLiteralPointer()
		if nlp == nil {
			nlp = elPtr.getUniverseOfDiscourse().NewNameLiteralPointer()
			nlp.SetOwningElementNoLock(elPtr)
		}
		nl = elPtr.getUniverseOfDiscourse().NewLiteral()
		nl.SetOwningElementNoLock(elPtr)
		nlp.setLiteral(nl)
	}
	nl.setLiteralValue(name)
}

// SetOwningElement() manages the owning element poiner belonging to this element.
// There are potentially four objects involved: the parent, the old parent (if
// there is one), the child (this element), and the owningElementPointer (oep).
// Because of the complexity of the wiring, all involved objects are locked here and
// the worker methods do not do any locking.
func (elPtr *element) SetOwningElement(parent Element) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	oldParent := elPtr.getOwningElement()
	if oldParent == nil && parent == nil {
		return // Nothing to do
	} else if oldParent != nil && parent != nil && oldParent.getId() != parent.getId() {
		return // Nothing to do
	}
	if oldParent != nil {
		oldParent.TraceableLock()
		defer oldParent.TraceableUnlock()
	}
	if parent != nil {
		parent.TraceableLock()
		defer parent.TraceableUnlock()
	}
	oep := elPtr.getOwningElementPointer()
	if oep != nil {
		oep.TraceableLock()
		defer oep.TraceableUnlock()
	}
	elPtr.SetOwningElementNoLock(parent)
}

// SetOwningElementNoLock() is a non-locking function that checks for the existence of the owningElementPointer, creating
// it if necessary. It then sets the owningElementPointer to point to the indicated element. Note that a side-effect
// of this action increments the version numbers of the owningElementPointer, this element, old and new owningElements, and all their
// owners, recureively.
func (elPtr *element) SetOwningElementNoLock(parent Element) {
	oep := elPtr.getOwningElementPointer()
	if oep == nil {
		oep = elPtr.uOfD.NewOwningElementPointer()
		oep.SetOwningElementNoLock(elPtr)
	}
	oep.setElement(parent)
}

func (elPtr *element) SetUri(uri string) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	nl := elPtr.getUriLiteral()
	if nl != nil {
		nl.TraceableLock()
		defer nl.TraceableUnlock()
	}
	nlp := elPtr.getUriLiteralPointer()
	if nlp != nil {
		nlp.TraceableLock()
		defer nlp.TraceableUnlock()
	}
	elPtr.SetUriNoLock(uri)
}

func (elPtr *element) SetUriNoLock(uri string) {
	nl := elPtr.getUriLiteral()
	if nl == nil {
		nlp := elPtr.getUriLiteralPointer()
		if nlp == nil {
			nlp = elPtr.getUniverseOfDiscourse().NewUriLiteralPointer()
			nlp.SetOwningElementNoLock(elPtr)
		}
		nl = elPtr.getUniverseOfDiscourse().NewLiteral()
		nl.SetOwningElementNoLock(elPtr)
		nlp.setLiteral(nl)
	}
	nl.setLiteralValue(uri)
}

type Element interface {
	BaseElement
	addOwnedBaseElement(BaseElement)
	childChanged(*ChangeNotification)
	GetAbstractElementsRecursivelyNoLock() []Element
	GetDefinition() string
	getDefinitionLiteral() Literal
	getDefinitionLiteralPointer() LiteralPointer
	getImmediateAbstractElements() []Element
	getImmediateAbstractions() []Refinement
	getImmediateRefinements() []Refinement
	GetName() string
	getNameLiteral() Literal
	getNameLiteralPointer() LiteralPointer
	getOwnedBaseElements() map[string]BaseElement
	GetOwnedBaseElements() []BaseElement
	GetOwnedElementsNoLock() []Element
	getOwningElementPointer() ElementPointer
	GetUniverseOfDiscourse() *UniverseOfDiscourse
	GetUniverseOfDiscourseNoLock() *UniverseOfDiscourse
	getUriLiteral() Literal
	getUriLiteralPointer() LiteralPointer
	internalAddOwnedBaseElement(BaseElement)
	internalRemoveOwnedBaseElement(BaseElement)
	MarshalJSON() ([]byte, error)
	removeOwnedBaseElement(BaseElement)
	SetDefinition(string)
	SetName(string)
	SetNameNoLock(string)
}
