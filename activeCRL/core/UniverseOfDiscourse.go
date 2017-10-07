package core

import (
	"errors"
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

type UniverseOfDiscourse struct {
	sync.RWMutex
	baseElementMap            *StringBaseElementMap
	baseElementListenerMap    *StringBaseElementPointerListMap
	elementListenerMap        *StringElementPointerListMap
	elementPointerListenerMap *StringElementPointerPointerListMap
	idUriMap                  *StringStringMap
	literalListenerMap        *StringLiteralPointerListMap
	literalPointerListenerMap *StringLiteralPointerPointerListMap
	undoMgr                   *undoManager
	uriBaseElementMap         *StringBaseElementMap
}

func NewUniverseOfDiscourse() *UniverseOfDiscourse {
	var uOfD UniverseOfDiscourse
	uOfD.baseElementMap = NewStringBaseElementMap()
	uOfD.baseElementListenerMap = NewStringBaseElementPointerListMap()
	uOfD.elementListenerMap = NewStringElementPointerListMap()
	uOfD.elementPointerListenerMap = NewStringElementPointerPointerListMap()
	uOfD.idUriMap = NewStringStringMap()
	uOfD.literalListenerMap = NewStringLiteralPointerListMap()
	uOfD.literalPointerListenerMap = NewStringLiteralPointerPointerListMap()
	uOfD.undoMgr = NewUndoManager()
	uOfD.uriBaseElementMap = NewStringBaseElementMap()
	uOfD.GetCoreConceptSpace()
	return &uOfD
}

func (uOfDPtr *UniverseOfDiscourse) AddBaseElement(be BaseElement, hl *HeldLocks) error {
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
	uOfDPtr.baseElementMap.SetEntry(be.GetId(hl).String(), be)
	uri := GetUri(be, hl)
	if uri != "" {
		uOfDPtr.uriBaseElementMap.SetEntry(uri, be)
		uOfDPtr.idUriMap.SetEntry(be.GetId(hl).String(), uri)
	}
	uOfDPtr.undoMgr.markNewBaseElement(be, hl)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) addBaseElementForUndo(be BaseElement, hl *HeldLocks) error {
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
	uOfDPtr.baseElementMap.SetEntry(be.GetId(hl).String(), be)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) addBaseElementListener(baseElement BaseElement, baseElementPointer BaseElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(baseElement)
	if baseElement != nil {
		elementId := baseElement.GetId(hl).String()
		uOfDPtr.baseElementListenerMap.AddEntry(elementId, baseElementPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) addElementListener(element Element, elementPointer ElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	if element != nil {
		elementId := element.GetId(hl).String()
		uOfDPtr.elementListenerMap.AddEntry(elementId, elementPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) addElementPointerListener(elementPointer ElementPointer, elementPointerPointer ElementPointerPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(elementPointer)
	if elementPointer != nil {
		elementId := elementPointer.GetId(hl).String()
		uOfDPtr.elementPointerListenerMap.AddEntry(elementId, elementPointerPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) addLiteralListener(literal Literal, literalPointer LiteralPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(literal)
	if literal != nil {
		literalId := literal.GetId(hl).String()
		uOfDPtr.literalListenerMap.AddEntry(literalId, literalPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) addLiteralPointerListener(literalPointer LiteralPointer, literalPointerPointer LiteralPointerPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(literalPointer)
	if literalPointer != nil {
		literalId := literalPointer.GetId(hl).String()
		uOfDPtr.literalPointerListenerMap.AddEntry(literalId, literalPointerPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) DeleteBaseElement(be BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be)
	if be != nil {
		//		log.Printf("About to delete BaseElement with id: %s\n", be.GetId(hl).String())
		SetOwningElement(be, nil, hl)
		uOfDPtr.removeBaseElement(be, hl)
		beId := be.GetId(hl).String()
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

func (uOfDPtr *UniverseOfDiscourse) GetBaseElement(id string) BaseElement {
	return uOfDPtr.baseElementMap.GetEntry(id)
}

func (uOfDPtr *UniverseOfDiscourse) GetBaseElementReferenceWithUri(uri string) BaseElementReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case BaseElementReference:
		return be.(BaseElementReference)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetBaseElementWithUri(uri string) BaseElement {
	return uOfDPtr.uriBaseElementMap.GetEntry(uri)
}

func (uOfDPtr *UniverseOfDiscourse) GetCoreConceptSpace() Element {
	coreConceptSpace := uOfDPtr.GetElementWithUri(CoreConceptSpaceUri)
	if coreConceptSpace == nil {
		coreConceptSpace = uOfDPtr.RecoverElement([]byte(serializedCore))
	}
	return coreConceptSpace
}

func (uOfDPtr *UniverseOfDiscourse) GetElement(id string) Element {
	// No locking required
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case Element:
		return be.(Element)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetElementPointer(id string) ElementPointer {
	// No locking required
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case ElementPointer:
		return be.(ElementPointer)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetElementWithUri(uri string) Element {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case Element:
		return be.(Element)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetElementPointerReferenceWithUri(uri string) ElementPointerReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case ElementPointerReference:
		return be.(ElementPointerReference)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetElementReferenceWithUri(uri string) ElementReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case ElementReference:
		return be.(ElementReference)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetLiteral(id string) Literal {
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case Literal:
		return be.(Literal)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetLiteralWithUri(uri string) Literal {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case Literal:
		return be.(Literal)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetLiteralPointer(id string) LiteralPointer {
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case LiteralPointer:
		return be.(LiteralPointer)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetLiteralReferenceWithUri(uri string) LiteralReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case LiteralReference:
		return be.(LiteralReference)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) GetLiteralPointerReferenceWithUri(uri string) LiteralPointerReference {
	be := uOfDPtr.GetBaseElementWithUri(uri)
	switch be.(type) {
	case LiteralPointerReference:
		return be.(LiteralPointerReference)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getRefinement(id string) Refinement {
	be := uOfDPtr.baseElementMap.GetEntry(id)
	switch be.(type) {
	case Refinement:
		return be.(Refinement)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) MarkUndoPoint() {
	uOfDPtr.undoMgr.MarkUndoPoint()
}

// NewElement() creates an initialized Element. No locking is required since the existence of
// the element is unknown outside this routine
func (uOfD *UniverseOfDiscourse) NewElement(hl *HeldLocks) Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el element
	el.initializeElement()
	uOfD.AddBaseElement(&el, hl)
	return &el
}

// NewAbstractElementPointer() creates and intitializes an elementPointer to play the role of an AbstractElementPointer
func (uOfD *UniverseOfDiscourse) NewAbstractElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = ABSTRACT_ELEMENT
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

// NewBaseElementPointer() creates and intitializes an elementPointer to play the role of an AbstractElementPointer
func (uOfD *UniverseOfDiscourse) NewBaseElementPointer(hl *HeldLocks) BaseElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep baseElementPointer
	ep.initializeBaseElementPointer()
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

func (uOfD *UniverseOfDiscourse) NewBaseElementReference(hl *HeldLocks) BaseElementReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el baseElementReference
	el.initializeBaseElementReference()
	uOfD.AddBaseElement(&el, hl)
	return &el
}

// NewRefinedElementPointer() creates and intitializes an elementPointer to play the role of an RefinedElementPointer
func (uOfD *UniverseOfDiscourse) NewRefinedElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = REFINED_ELEMENT
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

// NewOwningElementPointer() creates and intitializes an elementPointer to play the role of an OwningElementPointer
func (uOfD *UniverseOfDiscourse) NewOwningElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = OWNING_ELEMENT
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

// NewReferencedElementPointer() creates and intitializes an elementPointer to play the role of an ReferencedElementPointer
func (uOfD *UniverseOfDiscourse) NewReferencedElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = REFERENCED_ELEMENT
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

func (uOfD *UniverseOfDiscourse) NewElementPointerPointer(hl *HeldLocks) ElementPointerPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep elementPointerPointer
	ep.initializeElementPointerPointer()
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

func (uOfD *UniverseOfDiscourse) NewElementPointerReference(hl *HeldLocks) ElementPointerReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el elementPointerReference
	el.initializeElementPointerReference()
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfD *UniverseOfDiscourse) NewElementReference(hl *HeldLocks) ElementReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el elementReference
	el.initializeElementReference()
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfD *UniverseOfDiscourse) NewLiteral(hl *HeldLocks) Literal {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lit literal
	lit.initializeLiteral()
	uOfD.AddBaseElement(&lit, hl)
	return &lit
}

func (uOfD *UniverseOfDiscourse) NewNameLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = NAME
	uOfD.AddBaseElement(&lp, hl)
	return &lp
}

func (uOfD *UniverseOfDiscourse) NewDefinitionLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = DEFINITION
	uOfD.AddBaseElement(&lp, hl)
	return &lp
}

func (uOfD *UniverseOfDiscourse) NewUriLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = URI
	uOfD.AddBaseElement(&lp, hl)
	return &lp
}

func (uOfD *UniverseOfDiscourse) NewValueLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = VALUE
	uOfD.AddBaseElement(&lp, hl)
	return &lp
}

func (uOfD *UniverseOfDiscourse) NewLiteralPointerPointer(hl *HeldLocks) LiteralPointerPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var ep literalPointerPointer
	ep.initializeLiteralPointerPointer()
	uOfD.AddBaseElement(&ep, hl)
	return &ep
}

func (uOfD *UniverseOfDiscourse) NewLiteralPointerReference(hl *HeldLocks) LiteralPointerReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el literalPointerReference
	el.initializeLiteralPointerReference()
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfD *UniverseOfDiscourse) NewLiteralReference(hl *HeldLocks) LiteralReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el literalReference
	el.initializeLiteralReference()
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfD *UniverseOfDiscourse) NewRefinement(hl *HeldLocks) Refinement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	var el refinement
	el.initializeRefinement()
	uOfD.AddBaseElement(&el, hl)
	return &el
}

func (uOfDPtr *UniverseOfDiscourse) notifyElementListeners(notification *ChangeNotification, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.notifyElementListeners() called with nil HeldLocks")
	}
	switch notification.changedObject.(type) {
	case Element:
		id := notification.changedObject.GetId(hl).String()
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

func (uOfD *UniverseOfDiscourse) RecoverElement(data []byte) Element {
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

func (uOfDPtr *UniverseOfDiscourse) Redo(hl *HeldLocks) {
	if hl == nil {
		hl := NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	uOfDPtr.undoMgr.redo(uOfDPtr, hl)
}

func (uOfDPtr *UniverseOfDiscourse) removeBaseElement(be BaseElement, hl *HeldLocks) error {
	if hl == nil {
		return errors.New("UniverseOfDiscourse.removeBaseElement called with nil HeldLocks")
	}
	if be == nil {
		return errors.New("UniverseOfDiscource removeBaseElement failed because base element was nil")
	}
	hl.LockBaseElement(be)
	uOfDPtr.baseElementMap.DeleteEntry(be.GetId(hl).String())
	url := GetUri(be, hl)
	if url != "" {
		uOfDPtr.uriBaseElementMap.DeleteEntry(url)
		uOfDPtr.idUriMap.DeleteEntry(be.GetId(hl).String())
	}
	uOfDPtr.undoMgr.markRemovedBaseElement(be, hl)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) removeBaseElementForUndo(be BaseElement, hl *HeldLocks) {
	if be != nil {
		hl.LockBaseElement(be)
		if uOfDPtr.undoMgr.debugUndo == true {
			log.Printf("Removing base element for undo, id: %s\n", be.GetId(hl).String())
			Print(be, "RemovedBaseElement: ", hl)
		}
		uOfDPtr.baseElementMap.DeleteEntry(be.GetId(hl).String())
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeBaseElementListener(baseElement BaseElement, baseElementPointer BaseElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if baseElement != nil {
		elementId := baseElement.GetId(hl).String()
		uOfDPtr.baseElementListenerMap.RemoveEntry(elementId, baseElementPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeElementListener(element Element, elementPointer ElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if element != nil {
		elementId := element.GetId(hl).String()
		uOfDPtr.elementListenerMap.RemoveEntry(elementId, elementPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeElementPointerListener(elementPointer ElementPointer, elementPointerPointer ElementPointerPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if elementPointer != nil {
		elementId := elementPointer.GetId(hl).String()
		uOfDPtr.elementPointerListenerMap.RemoveEntry(elementId, elementPointerPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeLiteralListener(literal Literal, literalPointer LiteralPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if literal != nil {
		literalId := literal.GetId(hl).String()
		uOfDPtr.literalListenerMap.RemoveEntry(literalId, literalPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeLiteralPointerListener(literalPointer LiteralPointer, literalPointerPointer LiteralPointerPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	if literalPointer != nil {
		elementId := literalPointer.GetId(hl).String()
		uOfDPtr.literalPointerListenerMap.RemoveEntry(elementId, literalPointerPointer)
	}
}

func (uOfDPtr *UniverseOfDiscourse) restoreUriIndexRecursively(be BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	uri := GetUri(be, hl)
	if uri != "" {
		uOfDPtr.uriBaseElementMap.SetEntry(uri, be)
		uOfDPtr.idUriMap.SetEntry(be.GetId(hl).String(), uri)
	}

	switch be.(type) {
	case Element:
		for _, child := range be.(Element).GetOwnedBaseElements(hl) {
			uOfDPtr.restoreUriIndexRecursively(child, hl)
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) SetDebugUndo(newSetting bool) {
	uOfDPtr.undoMgr.setDebugUndo(newSetting)
}

func (uOfDPtr *UniverseOfDiscourse) SetRecordingUndo(newSetting bool) {
	uOfDPtr.undoMgr.setRecordingUndo(newSetting)
}

func (uOfDPtr *UniverseOfDiscourse) SetUniverseOfDiscourseRecursively(be BaseElement, hl *HeldLocks) {
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

func (uOfDPtr *UniverseOfDiscourse) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock Universe of Discourse %p\n", uOfDPtr)
	}
	uOfDPtr.Lock()
}

func (uOfDPtr *UniverseOfDiscourse) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock Universe of Discourse %p\n", uOfDPtr)
	}
	uOfDPtr.Unlock()
}

func (uOfDPtr *UniverseOfDiscourse) Undo(hl *HeldLocks) {
	if hl == nil {
		hl := NewHeldLocks(nil)
		hl.ReleaseLocks()
	}
	uOfDPtr.undoMgr.undo(uOfDPtr, hl)
}
