// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

type universeOfDiscourse struct {
	element
	sync.RWMutex
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

func NewUniverseOfDiscourse() UniverseOfDiscourse {
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
	hl := NewHeldLocks(nil)
	defer hl.ReleaseLocks()
	buildCoreConceptSpace(&uOfD, hl)
	uOfD.initializeElement(UniverseOfDiscourseUri)
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
	if be.GetId(hl) == uuid.Nil {
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
		log.Printf("Adding base element for undo, id: %s\n", be.GetId(hl).String())
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

func (uOfDPtr *universeOfDiscourse) GetBaseElement(id uuid.UUID) BaseElement {
	return uOfDPtr.baseElementMap.GetEntry(id)
}

func (uOfDPtr *universeOfDiscourse) GetBaseElements() []BaseElement {
	var baseElements []BaseElement
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

func (uOfDPtr *universeOfDiscourse) GetElement(id uuid.UUID) Element {
	// No locking required
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case Element:
		return be.(Element)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) GetElementPointer(id uuid.UUID) ElementPointer {
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

func (uOfDPtr *universeOfDiscourse) GetLiteral(id uuid.UUID) Literal {
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

func (uOfDPtr *universeOfDiscourse) GetLiteralPointer(id uuid.UUID) LiteralPointer {
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

func (uOfDPtr *universeOfDiscourse) getRefinement(id uuid.UUID) Refinement {
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case Refinement:
		return be.(Refinement)
	}
	return nil
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

func (uOfD *universeOfDiscourse) NewNameLiteralPointer(hl *HeldLocks, uri ...string) LiteralPointer {
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
				// Must suppress circular notifications
				if notification.isReferenced(elementPointer) == false {
					newNotification := NewChangeNotification(elementPointer, MODIFY, notification)
					propagateChange(elementPointer, newNotification, hl)
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
	return nil
}

func (uOfDPtr *universeOfDiscourse) removeBaseElementForUndo(be BaseElement, hl *HeldLocks) {
	if be != nil {
		hl.LockBaseElement(be)
		if uOfDPtr.undoMgr.debugUndo == true {
			log.Printf("Removing base element for undo, id: %s\n", be.GetId(hl).String())
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

type UniverseOfDiscourse interface {
	Element
	AddBaseElement(BaseElement, *HeldLocks) error
	DeleteBaseElement(BaseElement, *HeldLocks)
	GetBaseElement(uuid.UUID) BaseElement
	GetBaseElements() []BaseElement
	GetBaseElementReferenceWithUri(string) BaseElementReference
	GetBaseElementWithUri(string) BaseElement
	GetCoreConceptSpace() Element
	GetElement(uuid.UUID) Element
	GetElementPointer(uuid.UUID) ElementPointer
	GetElementWithUri(string) Element
	GetElementPointerReferenceWithUri(string) ElementPointerReference
	GetElementReferenceWithUri(string) ElementReference
	GetLiteral(uuid.UUID) Literal
	GetLiteralWithUri(string) Literal
	GetLiteralPointer(uuid.UUID) LiteralPointer
	GetLiteralReferenceWithUri(string) LiteralReference
	GetLiteralPointerReferenceWithUri(string) LiteralPointerReference
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
	NewNameLiteralPointer(*HeldLocks, ...string) LiteralPointer
	NewDefinitionLiteralPointer(*HeldLocks, ...string) LiteralPointer
	NewUriLiteralPointer(*HeldLocks, ...string) LiteralPointer
	NewValueLiteralPointer(*HeldLocks, ...string) LiteralPointer
	NewLiteralPointerPointer(*HeldLocks, ...string) LiteralPointerPointer
	NewLiteralPointerReference(*HeldLocks, ...string) LiteralPointerReference
	NewLiteralReference(*HeldLocks, ...string) LiteralReference
	NewRefinement(*HeldLocks, ...string) Refinement
	RecoverElement([]byte) Element
	Redo(*HeldLocks)
	SetDebugUndo(bool)
	SetRecordingUndo(bool)
	SetUniverseOfDiscourseRecursively(BaseElement, *HeldLocks)
	//	TraceableLock()
	//	TraceableUnlock()
	Undo(*HeldLocks)
}
