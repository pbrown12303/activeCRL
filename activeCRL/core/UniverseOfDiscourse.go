// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	//	"github.com/satori/go.uuid"
	"log"
)

type universeOfDiscourse struct {
	element
	baseElementMap            *UUIDBaseElementMap
	baseElementListenerMap    *UUIDBaseElementPointerListMap
	elementListenerMap        *UUIDElementPointerListMap
	elementPointerListenerMap *UUIDElementPointerPointerListMap
	idUriMap                  *UUIDStringMap
	literalListenerMap        *UUIDLiteralPointerListMap
	literalPointerListenerMap *UUIDLiteralPointerPointerListMap
	undoMgr                   *undoManager
	uriBaseElementMap         *StringBaseElementMap
}

func NewUniverseOfDiscourse(hl *HeldLocks) UniverseOfDiscourse {
	var uOfD universeOfDiscourse
	uOfD.baseElementMap = NewUUIDBaseElementMap()
	uOfD.baseElementListenerMap = NewUUIDBaseElementPointerListMap()
	uOfD.elementListenerMap = NewUUIDElementPointerListMap()
	uOfD.elementPointerListenerMap = NewUUIDElementPointerPointerListMap()
	uOfD.idUriMap = NewUUIDStringMap()
	uOfD.literalListenerMap = NewUUIDLiteralPointerListMap()
	uOfD.literalPointerListenerMap = NewUUIDLiteralPointerPointerListMap()
	uOfD.undoMgr = NewUndoManager()
	uOfD.uriBaseElementMap = NewStringBaseElementMap()
	uOfD.initializeElement(UniverseOfDiscourseUri)
	buildCoreConceptSpace(&uOfD, hl)
	uOfD.AddBaseElement(&uOfD, hl)
	return &uOfD
}

func (uOfDPtr *universeOfDiscourse) AddBaseElement(be BaseElement, hl *HeldLocks) error {
	if be == nil {
		return errors.New("UniverseOfDiscource addBaseElement() failed because base element was nil")
	}
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if be.GetId(hl) == "" {
		return errors.New("UniverseOfDiscource addBaseElement() failed because UUID was nil")
	}
	hl.LockBaseElement(be)
	be.setUniverseOfDiscourse(uOfDPtr, hl)
	uOfDPtr.baseElementMap.SetEntry(be.GetId(hl), be)
	uri := GetUri(be, hl)
	if uri != "" {
		uOfDPtr.uriBaseElementMap.SetEntry(uri, be)
		uOfDPtr.idUriMap.SetEntry(be.GetId(hl), uri)
	}
	uOfDPtr.undoMgr.markNewBaseElement(be, hl)
	notification := NewChangeNotification(be, ADD, "AddBaseElement", nil)
	uOfDPtr.uOfDChanged(notification, hl)
	return nil
}

func (uOfDPtr *universeOfDiscourse) addBaseElementForUndo(be BaseElement, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.addBaseElementForUndo() called with nil HeldLocks")
	}
	if be == nil {
		return errors.New("UniverseOfDiscource addBaseElementForUndo() failed because base element was nil")
	}
	if be != nil {
		hl.LockBaseElement(be)
	}
	if uOfDPtr.undoMgr.debugUndo == true {
		log.Printf("Adding base element for undo, id: %s\n", be.GetId(hl))
		Print(be, "AddedBaseElement: ", hl)
	}
	uOfDPtr.baseElementMap.SetEntry(be.GetId(hl), be)
	return nil
}

func (uOfDPtr *universeOfDiscourse) addBaseElementListener(baseElement BaseElement, baseElementPointer BaseElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(baseElement)
	if baseElement != nil {
		elementId := baseElement.GetId(hl)
		uOfDPtr.baseElementListenerMap.AddEntry(elementId, baseElementPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) addElementListener(element Element, elementPointer ElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	if element != nil {
		elementId := element.GetId(hl)
		uOfDPtr.elementListenerMap.AddEntry(elementId, elementPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) addElementPointerListener(elementPointer ElementPointer, elementPointerPointer ElementPointerPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elementPointer)
	if elementPointer != nil {
		elementId := elementPointer.GetId(hl)
		uOfDPtr.elementPointerListenerMap.AddEntry(elementId, elementPointerPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) addLiteralListener(literal Literal, literalPointer LiteralPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(literal)
	if literal != nil {
		literalId := literal.GetId(hl)
		uOfDPtr.literalListenerMap.AddEntry(literalId, literalPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) addLiteralPointerListener(literalPointer LiteralPointer, literalPointerPointer LiteralPointerPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(literalPointer)
	if literalPointer != nil {
		literalId := literalPointer.GetId(hl)
		uOfDPtr.literalPointerListenerMap.AddEntry(literalId, literalPointerPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) DeleteBaseElement(be BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be)
	if be != nil {
		//		log.Printf("About to delete BaseElement with id: %s\n", be.GetId(hl).String())
		SetOwningElement(be, nil, hl)
		uOfDPtr.removeBaseElement(be, hl)
		beId := be.GetId(hl)
		bepl := uOfDPtr.baseElementListenerMap.GetEntry(beId)
		if bepl != nil {
			for _, bep := range *bepl {
				bep.SetBaseElement(nil, hl)
			}
		}
		switch be.(type) {
		case Element:
			epl := uOfDPtr.elementListenerMap.GetEntry(beId)
			if epl != nil {
				for _, elementPointer := range *epl {
					elementPointer.SetElement(nil, hl)
				}
				for _, child := range be.(Element).GetOwnedBaseElements(hl) {
					uOfDPtr.DeleteBaseElement(child, hl)
				}
			}
		case ElementPointer:
		case Literal:
		case LiteralPointer:
		}
	}
}

// GetAbstractElementsRecursively() returns all of the elements abstractions
func (uOfDPtr *universeOfDiscourse) GetAbstractElementsRecursively(el Element, hl *HeldLocks) []Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(el)
	abstractElements := uOfDPtr.getImmediateAbstractElements(el, hl)
	var ancestors []Element
	for _, element := range abstractElements {
		for _, ancestor := range uOfDPtr.GetAbstractElementsRecursively(element, hl) {
			ancestors = append(ancestors, ancestor)
		}
	}
	for _, ancestor := range ancestors {
		abstractElements = append(abstractElements, ancestor)
	}
	return abstractElements
}

func (uOfDPtr *universeOfDiscourse) GetBaseElement(id string) BaseElement {
	return uOfDPtr.baseElementMap.GetEntry(id)
}

func (uOfDPtr *universeOfDiscourse) GetBaseElements() []BaseElement {
	uOfDPtr.baseElementMap.TraceableLock()
	defer uOfDPtr.baseElementMap.TraceableUnlock()
	var baseElements []BaseElement
	for _, be := range uOfDPtr.baseElementMap.baseElementMap {
		baseElements = append(baseElements, be)
	}
	return baseElements
}

func (uOfDPtr *universeOfDiscourse) GetBaseElementReferenceWithUri(uri string) BaseElementReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case BaseElementReference:
		return be.(BaseElementReference)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetBaseElementWithUri(uri string) BaseElement {
	return uOfDPtr.uriBaseElementMap.GetEntry(uri)
}

func (uOfDPtr *universeOfDiscourse) GetCoreConceptSpace() Element {
	return uOfDPtr.GetElementWithUri(CoreConceptSpaceUri)
}

func (uOfDPtr *universeOfDiscourse) GetElement(id string) Element {
	// No locking required
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case Element:
		return be.(Element)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetElementPointer(id string) ElementPointer {
	// No locking required
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case ElementPointer:
		return be.(ElementPointer)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetElementWithUri(uri string) Element {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case Element:
		return be.(Element)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetElementPointerReferenceWithUri(uri string) ElementPointerReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case ElementPointerReference:
		return be.(ElementPointerReference)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetElementReferenceWithUri(uri string) ElementReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case ElementReference:
		return be.(ElementReference)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) getImmediateAbstractElements(el Element, hl *HeldLocks) []Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(el)
	var abstractElements []Element
	abstractions := uOfDPtr.getImmediateAbstractions(el, hl)
	if abstractions != nil {
		for _, abstraction := range abstractions {
			if abstraction.GetAbstractElement(hl) != nil {
				abstractElements = append(abstractElements, abstraction.GetAbstractElement(hl))
			}
		}
	}
	return abstractElements
}

func (uOfDPtr *universeOfDiscourse) getImmediateAbstractions(el Element, hl *HeldLocks) []Refinement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(el)
	var abstractions []Refinement
	ePtrs := uOfDPtr.elementListenerMap.GetEntry(el.GetId(hl))
	if ePtrs != nil {
		for _, ePtr := range *ePtrs {
			if ePtr.GetElementPointerRole(hl) == REFINED_ELEMENT {
				abstractions = append(abstractions, GetOwningElement(ePtr, hl).(Refinement))
			}
		}
	}
	return abstractions
}

func (uOfDPtr *universeOfDiscourse) getImmediateRefinements(el Element, hl *HeldLocks) []Refinement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(el)
	var refinements []Refinement
	ePtrs := uOfDPtr.elementListenerMap.GetEntry(el.GetId(hl))
	if ePtrs != nil {
		for _, ePtr := range *ePtrs {
			if ePtr.GetElementPointerRole(hl) == ABSTRACT_ELEMENT {
				refinements = append(refinements, GetOwningElement(ePtr, hl).(Refinement))
			}
		}
	}
	return refinements
}

func (uOfDPtr *universeOfDiscourse) GetLiteral(id string) Literal {
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case Literal:
		return be.(Literal)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetLiteralWithUri(uri string) Literal {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case Literal:
		return be.(Literal)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetLiteralPointer(id string) LiteralPointer {
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case LiteralPointer:
		return be.(LiteralPointer)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetLiteralReferenceWithUri(uri string) LiteralReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case LiteralReference:
		return be.(LiteralReference)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetLiteralPointerReferenceWithUri(uri string) LiteralPointerReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case LiteralPointerReference:
		return be.(LiteralPointerReference)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) getRefinement(id string) Refinement {
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case Refinement:
		return be.(Refinement)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) IsRefinementOf(refinedElement Element, abstractElement Element, hl *HeldLocks) bool {
	for _, candidateElement := range uOfDPtr.GetAbstractElementsRecursively(refinedElement, hl) {
		if candidateElement == abstractElement {
			return true
		}
	}
	return false
}

func (uOfDPtr *universeOfDiscourse) IsRecordingUndo() bool {
	return uOfDPtr.undoMgr.recordingUndo
}

func (uOfDPtr *universeOfDiscourse) MarkUndoPoint() {
	uOfDPtr.undoMgr.MarkUndoPoint()
}

// NewElement() creates an initialized Element. No locking is required since the existence of
// the element is unknown outside this routine
func (uOfD *universeOfDiscourse) NewElement(hl *HeldLocks, uri ...string) Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el element
	el.initializeElement(uri...)
	uOfD.AddBaseElement(&el, hl)
	return &el
}

// NewAbstractElementPointer() creates and intitializes an elementPointer to play the role of an AbstractElementPointer
func (uOfD *universeOfDiscourse) NewAbstractElementPointer(hl *HeldLocks, uri ...string) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointer
	ep.initializeElementPointer(uri...)
	ep.elementPointerRole = ABSTRACT_ELEMENT
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

// NewBaseElementPointer() creates and intitializes an elementPointer to play the role of an AbstractElementPointer
func (uOfD *universeOfDiscourse) NewBaseElementPointer(hl *HeldLocks, uri ...string) BaseElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep baseElementPointer
	ep.initializeBaseElementPointer(uri...)
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

func (uOfD *universeOfDiscourse) NewBaseElementReference(hl *HeldLocks, uri ...string) BaseElementReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el baseElementReference
	el.initializeBaseElementReference(uri...)
	uOfD.AddBaseElement(&el, hl)
	return &el
}

// NewRefinedElementPointer() creates and intitializes an elementPointer to play the role of an RefinedElementPointer
func (uOfD *universeOfDiscourse) NewRefinedElementPointer(hl *HeldLocks, uri ...string) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointer
	ep.initializeElementPointer(uri...)
	ep.elementPointerRole = REFINED_ELEMENT
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

// NewOwningElementPointer() creates and intitializes an elementPointer to play the role of an OwningElementPointer
func (uOfD *universeOfDiscourse) NewOwningElementPointer(hl *HeldLocks, uri ...string) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointer
	ep.initializeElementPointer(uri...)
	ep.elementPointerRole = OWNING_ELEMENT
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

// NewReferencedElementPointer() creates and intitializes an elementPointer to play the role of an ReferencedElementPointer
func (uOfD *universeOfDiscourse) NewReferencedElementPointer(hl *HeldLocks, uri ...string) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointer
	ep.initializeElementPointer(uri...)
	ep.elementPointerRole = REFERENCED_ELEMENT
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

func (uOfD *universeOfDiscourse) NewElementPointerPointer(hl *HeldLocks, uri ...string) ElementPointerPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointerPointer
	ep.initializeElementPointerPointer(uri...)
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

func (uOfD *universeOfDiscourse) NewElementPointerReference(hl *HeldLocks, uri ...string) ElementPointerReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el elementPointerReference
	el.initializeElementPointerReference(uri...)
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfD *universeOfDiscourse) NewElementReference(hl *HeldLocks, uri ...string) ElementReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el elementReference
	el.initializeElementReference(uri...)
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfD *universeOfDiscourse) NewLiteral(hl *HeldLocks, uri ...string) Literal {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lit literal
	lit.initializeLiteral(uri...)
	uOfD.AddBaseElement(&lit, hl)
	return &lit
}

func (uOfD *universeOfDiscourse) NewLabelLiteralPointer(hl *HeldLocks, uri ...string) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lp literalPointer
	lp.initializeLiteralPointer(uri...)
	lp.literalPointerRole = NAME
	uOfD.AddBaseElement(&lp, hl)
	return &lp
}

func (uOfD *universeOfDiscourse) NewDefinitionLiteralPointer(hl *HeldLocks, uri ...string) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lp literalPointer
	lp.initializeLiteralPointer(uri...)
	lp.literalPointerRole = DEFINITION
	uOfD.AddBaseElement(&lp, hl)
	return &lp
}

func (uOfD *universeOfDiscourse) NewUriLiteralPointer(hl *HeldLocks, uri ...string) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lp literalPointer
	lp.initializeLiteralPointer(uri...)
	lp.literalPointerRole = URI
	uOfD.AddBaseElement(&lp, hl)
	return &lp
}

func (uOfD *universeOfDiscourse) NewValueLiteralPointer(hl *HeldLocks, uri ...string) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lp literalPointer
	lp.initializeLiteralPointer(uri...)
	lp.literalPointerRole = VALUE
	uOfD.AddBaseElement(&lp, hl)
	return &lp
}

func (uOfD *universeOfDiscourse) NewLiteralPointerPointer(hl *HeldLocks, uri ...string) LiteralPointerPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep literalPointerPointer
	ep.initializeLiteralPointerPointer(uri...)
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

func (uOfD *universeOfDiscourse) NewLiteralPointerReference(hl *HeldLocks, uri ...string) LiteralPointerReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el literalPointerReference
	el.initializeLiteralPointerReference(uri...)
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfD *universeOfDiscourse) NewLiteralReference(hl *HeldLocks, uri ...string) LiteralReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el literalReference
	el.initializeLiteralReference(uri...)
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfD *universeOfDiscourse) NewRefinement(hl *HeldLocks, uri ...string) Refinement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el refinement
	el.initializeRefinement(uri...)
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfDPtr *universeOfDiscourse) notifyBaseElementListeners(notification *ChangeNotification, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.notifyElementListeners() called with nil HeldLocks")
	}
	id := notification.changedObject.GetId(hl)
	bepl := uOfDPtr.baseElementListenerMap.GetEntry(id)
	if bepl != nil {
		for _, baseElementPointer := range *bepl {
			// Must suppress circular notifications
			if notification.isReferenced(baseElementPointer) == false {
				newNotification := NewChangeNotification(baseElementPointer, MODIFY, "notifyBaseElementListeners", notification)
				indicatedBaseElementChanged(baseElementPointer, newNotification, hl)
			}
		}
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) notifyElementListeners(notification *ChangeNotification, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.notifyElementListeners() called with nil HeldLocks")
	}
	switch notification.changedObject.(type) {
	case Element:
		id := notification.changedObject.GetId(hl)
		epl := uOfDPtr.elementListenerMap.GetEntry(id)
		if epl != nil {
			for _, elementPointer := range *epl {
				// Determine whether the pointer is an OWNING_ELEMENT pointer. If it is, see if the change is a
				// modification. We do not want to propagate modifications to owners down to their owned elements.
				isOwningElementPointer := elementPointer.GetElementPointerRole(hl) == OWNING_ELEMENT
				isModification := notification.natureOfChange == MODIFY
				if !(isOwningElementPointer && isModification) {
					newNotification := NewChangeNotification(elementPointer, MODIFY, "notifyElementListeners", notification)
					indicatedBaseElementChanged(elementPointer, newNotification, hl)
				}
			}
		}
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) notifyElementPointerListeners(notification *ChangeNotification, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.notifyElementPointerListeners() called with nil HeldLocks")
	}
	switch notification.changedObject.(type) {
	case ElementPointer:
		id := notification.changedObject.GetId(hl)
		epl := uOfDPtr.elementPointerListenerMap.GetEntry(id)
		if epl != nil {
			for _, elementPointerPointer := range *epl {
				// Must suppress circular notifications
				if notification.isReferenced(elementPointerPointer) == false {
					newNotification := NewChangeNotification(elementPointerPointer, MODIFY, "notifyElementPointerListeners", notification)
					indicatedBaseElementChanged(elementPointerPointer, newNotification, hl)
				}
			}
		}
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) notifyListeners(be BaseElement, notification *ChangeNotification, hl *HeldLocks) {
	// If the element that changed (the base element) is already a source of change in the notification's prior history, we have a circular
	// reference. No notification should be performed
	priorChange := notification.underlyingChange
	if priorChange != nil && priorChange.isReferenced(notification.changedObject) {
		return
	}

	uOfDPtr.notifyBaseElementListeners(notification, hl)
	switch be.(type) {
	case Element:
		uOfDPtr.notifyElementListeners(notification, hl)
	case ElementPointer:
		uOfDPtr.notifyElementPointerListeners(notification, hl)
	case Literal:
		uOfDPtr.notifyLiteralListeners(notification, hl)
	case LiteralPointer:
		uOfDPtr.notifyLiteralPointerListeners(notification, hl)
	case UniverseOfDiscourse:
		uOfDPtr.notifyElementListeners(notification, hl)
	}
}

func (uOfDPtr *universeOfDiscourse) notifyLiteralListeners(notification *ChangeNotification, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.notifyLiteralListeners() called with nil HeldLocks")
	}
	switch notification.changedObject.(type) {
	case Literal:
		id := notification.changedObject.GetId(hl)
		lpl := uOfDPtr.literalListenerMap.GetEntry(id)
		if lpl != nil {
			for _, literalPointer := range *lpl {
				// Must suppress circular notifications
				if notification.isReferenced(literalPointer) == false {
					newNotification := NewChangeNotification(literalPointer, MODIFY, "notifyLiteralListeners", notification)
					indicatedBaseElementChanged(literalPointer, newNotification, hl)
				}
			}
		}
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) notifyLiteralPointerListeners(notification *ChangeNotification, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.notifyLiteralPointerListeners() called with nil HeldLocks")
	}
	switch notification.changedObject.(type) {
	case LiteralPointer:
		id := notification.changedObject.GetId(hl)
		epl := uOfDPtr.literalPointerListenerMap.GetEntry(id)
		if epl != nil {
			for _, literalPointerPointer := range *epl {
				// Must suppress circular notifications
				if notification.isReferenced(literalPointerPointer) == false {
					newNotification := NewChangeNotification(literalPointerPointer, MODIFY, "notifyLiteralPointerListeners", notification)
					indicatedBaseElementChanged(literalPointerPointer, newNotification, hl)
				}
			}
		}
	}
	return nil
}

func (uOfD *universeOfDiscourse) RecoverElement(data []byte) Element {
	if len(data) == 0 {
		return nil
	}
	var recoveredElement BaseElement
	err := unmarshalPolymorphicBaseElement(data, &recoveredElement)
	if err != nil {
		log.Printf("Error recovering Element: %s \n", err)
		return nil
	}
	hl := NewHeldLocks(nil)
	defer hl.ReleaseLocks()
	uOfD.SetUniverseOfDiscourseRecursively(recoveredElement, hl)
	restoreValueOwningElementFieldsRecursively(recoveredElement.(Element), hl)
	uOfD.restoreUriIndexRecursively(recoveredElement, hl)
	return recoveredElement.(Element)
}

func (uOfDPtr *universeOfDiscourse) Redo(hl *HeldLocks) {
	if hl == nil {
		hl := NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	uOfDPtr.undoMgr.redo(uOfDPtr, hl)
}

func (uOfDPtr *universeOfDiscourse) removeBaseElement(be BaseElement, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.removeBaseElement called with nil HeldLocks")
	}
	if be == nil {
		return errors.New("UniverseOfDiscource removeBaseElement failed because base element was nil")
	}
	hl.LockBaseElement(be)
	uOfDPtr.baseElementMap.DeleteEntry(be.GetId(hl))
	url := GetUri(be, hl)
	if url != "" {
		uOfDPtr.uriBaseElementMap.DeleteEntry(url)
		uOfDPtr.idUriMap.DeleteEntry(be.GetId(hl))
	}
	uOfDPtr.undoMgr.markRemovedBaseElement(be, hl)
	notification := NewChangeNotification(be, ADD, "removeBaseElement", nil)
	uOfDPtr.uOfDChanged(notification, hl)
	return nil
}

func (uOfDPtr *universeOfDiscourse) removeBaseElementForUndo(be BaseElement, hl *HeldLocks) {
	if be != nil {
		hl.LockBaseElement(be)
		if uOfDPtr.undoMgr.debugUndo == true {
			log.Printf("Removing base element for undo, id: %s\n", be.GetId(hl))
			Print(be, "RemovedBaseElement: ", hl)
		}
		uOfDPtr.baseElementMap.DeleteEntry(be.GetId(hl))
	}
}

func (uOfDPtr *universeOfDiscourse) removeBaseElementListener(baseElement BaseElement, baseElementPointer BaseElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if baseElement != nil {
		elementId := baseElement.GetId(hl)
		uOfDPtr.baseElementListenerMap.RemoveEntry(elementId, baseElementPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) removeElementListener(element Element, elementPointer ElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if element != nil {
		elementId := element.GetId(hl)
		uOfDPtr.elementListenerMap.RemoveEntry(elementId, elementPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) removeElementPointerListener(elementPointer ElementPointer, elementPointerPointer ElementPointerPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if elementPointer != nil {
		elementId := elementPointer.GetId(hl)
		uOfDPtr.elementPointerListenerMap.RemoveEntry(elementId, elementPointerPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) removeLiteralListener(literal Literal, literalPointer LiteralPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if literal != nil {
		literalId := literal.GetId(hl)
		uOfDPtr.literalListenerMap.RemoveEntry(literalId, literalPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) removeLiteralPointerListener(literalPointer LiteralPointer, literalPointerPointer LiteralPointerPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if literalPointer != nil {
		elementId := literalPointer.GetId(hl)
		uOfDPtr.literalPointerListenerMap.RemoveEntry(elementId, literalPointerPointer)
	}
}

func (uOfDPtr *universeOfDiscourse) restoreUriIndexRecursively(be BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	uri := GetUri(be, hl)
	if uri != "" {
		uOfDPtr.uriBaseElementMap.SetEntry(uri, be)
		uOfDPtr.idUriMap.SetEntry(be.GetId(hl), uri)
	}

	switch be.(type) {
	case Element:
		for _, child := range be.(Element).GetOwnedBaseElements(hl) {
			uOfDPtr.restoreUriIndexRecursively(child, hl)
		}
	}
}

func (uOfDPtr *universeOfDiscourse) SetDebugUndo(newSetting bool) {
	uOfDPtr.undoMgr.setDebugUndo(newSetting)
}

func (uOfDPtr *universeOfDiscourse) SetRecordingUndo(newSetting bool) {
	uOfDPtr.undoMgr.setRecordingUndo(newSetting)
}

func (uOfDPtr *universeOfDiscourse) SetUniverseOfDiscourseRecursively(be BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	uOfDPtr.AddBaseElement(be, hl)
	switch be.(type) {
	case Element:
		for _, child := range be.(Element).GetOwnedBaseElements(hl) {
			uOfDPtr.SetUniverseOfDiscourseRecursively(child, hl)
		}
	}
}

//func (uOfDPtr *universeOfDiscourse) TraceableLock() {
//	if TraceLocks {
//		log.Printf("About to lock Universe of Discourse %p\n", uOfDPtr)
//	}
//	uOfDPtr.Lock()
//}
//
//func (uOfDPtr *universeOfDiscourse) TraceableUnlock() {
//	if TraceLocks {
//		log.Printf("About to unlock Universe of Discourse %p\n", uOfDPtr)
//	}
//	uOfDPtr.Unlock()
//}
//
func (uOfDPtr *universeOfDiscourse) Undo(hl *HeldLocks) {
	if hl == nil {
		hl := NewHeldLocks(nil)
		hl.ReleaseLocks()
	}
	uOfDPtr.undoMgr.undo(uOfDPtr, hl)
}

func (uOfD *universeOfDiscourse) updateUriIndices(be BaseElement, hl *HeldLocks) {
	id := be.GetId(hl)
	oldUri := uOfD.idUriMap.GetEntry(id)
	newUri := GetUri(be, hl)
	if oldUri != newUri {
		if oldUri != "" {
			uOfD.uriBaseElementMap.DeleteEntry(oldUri)
		}
		if newUri == "" {
			uOfD.idUriMap.DeleteEntry(id)
		} else {
			uOfD.idUriMap.SetEntry(id, newUri)
			uOfD.uriBaseElementMap.SetEntry(newUri, be)
		}
	}
}

func (uOfDPtr *universeOfDiscourse) uOfDChanged(notification *ChangeNotification, hl *HeldLocks) {
	if TraceChange == true {
		log.Printf("uOfDChanged called, uOfD ID: %s \n", uOfDPtr.GetId(hl))
		notification.Print("uOfDChanged Incoming Notification: ", hl)
	}
	newNotification := NewChangeNotification(uOfDPtr, MODIFY, "uOfDChanged", notification)
	indicatedBaseElementChanged(uOfDPtr, newNotification, hl)
}

type UniverseOfDiscourse interface {
	Element
	AddBaseElement(BaseElement, *HeldLocks) error
	addBaseElementListener(BaseElement, BaseElementPointer, *HeldLocks)
	DeleteBaseElement(BaseElement, *HeldLocks)
	GetAbstractElementsRecursively(Element, *HeldLocks) []Element
	GetBaseElement(string) BaseElement
	GetBaseElements() []BaseElement
	GetBaseElementReferenceWithUri(string) BaseElementReference
	GetBaseElementWithUri(string) BaseElement
	GetCoreConceptSpace() Element
	GetElement(string) Element
	GetElementPointer(string) ElementPointer
	GetElementWithUri(string) Element
	GetElementPointerReferenceWithUri(string) ElementPointerReference
	GetElementReferenceWithUri(string) ElementReference
	getImmediateAbstractElements(Element, *HeldLocks) []Element
	getImmediateAbstractions(Element, *HeldLocks) []Refinement
	getImmediateRefinements(Element, *HeldLocks) []Refinement
	GetLiteral(string) Literal
	GetLiteralWithUri(string) Literal
	GetLiteralPointer(string) LiteralPointer
	GetLiteralReferenceWithUri(string) LiteralReference
	GetLiteralPointerReferenceWithUri(string) LiteralPointerReference
	IsRefinementOf(Element, Element, *HeldLocks) bool
	IsRecordingUndo() bool
	MarkUndoPoint()
	NewElement(*HeldLocks, ...string) Element
	NewAbstractElementPointer(*HeldLocks, ...string) ElementPointer
	NewBaseElementPointer(*HeldLocks, ...string) BaseElementPointer
	NewBaseElementReference(*HeldLocks, ...string) BaseElementReference
	NewRefinedElementPointer(*HeldLocks, ...string) ElementPointer
	NewOwningElementPointer(*HeldLocks, ...string) ElementPointer
	NewReferencedElementPointer(*HeldLocks, ...string) ElementPointer
	NewElementPointerPointer(*HeldLocks, ...string) ElementPointerPointer
	NewElementPointerReference(*HeldLocks, ...string) ElementPointerReference
	NewElementReference(*HeldLocks, ...string) ElementReference
	NewLiteral(*HeldLocks, ...string) Literal
	NewLabelLiteralPointer(*HeldLocks, ...string) LiteralPointer
	NewDefinitionLiteralPointer(*HeldLocks, ...string) LiteralPointer
	NewUriLiteralPointer(*HeldLocks, ...string) LiteralPointer
	NewValueLiteralPointer(*HeldLocks, ...string) LiteralPointer
	NewLiteralPointerPointer(*HeldLocks, ...string) LiteralPointerPointer
	NewLiteralPointerReference(*HeldLocks, ...string) LiteralPointerReference
	NewLiteralReference(*HeldLocks, ...string) LiteralReference
	NewRefinement(*HeldLocks, ...string) Refinement
	notifyBaseElementListeners(*ChangeNotification, *HeldLocks) error
	notifyElementListeners(*ChangeNotification, *HeldLocks) error
	notifyElementPointerListeners(*ChangeNotification, *HeldLocks) error
	notifyListeners(BaseElement, *ChangeNotification, *HeldLocks)
	notifyLiteralListeners(*ChangeNotification, *HeldLocks) error
	notifyLiteralPointerListeners(*ChangeNotification, *HeldLocks) error
	RecoverElement([]byte) Element
	Redo(*HeldLocks)
	removeBaseElementListener(BaseElement, BaseElementPointer, *HeldLocks)
	SetDebugUndo(bool)
	SetRecordingUndo(bool)
	SetUniverseOfDiscourseRecursively(BaseElement, *HeldLocks)
	Undo(*HeldLocks)
	uOfDChanged(*ChangeNotification, *HeldLocks)
	updateUriIndices(BaseElement, *HeldLocks)
}
