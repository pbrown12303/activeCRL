package core

import (
	"errors"
	"log"
	"runtime/debug"
	"sync"

	"github.com/satori/go.uuid"
)

type elementPointerList *[]ElementPointer
type elementPointerPointerList *[]ElementPointerPointer
type literalPointerList *[]LiteralPointer
type literalPointerPointerList *[]LiteralPointerPointer

type UniverseOfDiscourse struct {
	sync.RWMutex
	baseElementMap            map[string]BaseElement
	uriBaseElementMap         map[string]BaseElement
	elementListenerMap        map[string]elementPointerList
	elementPointerListenerMap map[string]elementPointerPointerList
	literalListenerMap        map[string]literalPointerList
	literalPointerListenerMap map[string]literalPointerPointerList
	recordingUndo             bool
	undoStack                 undoStack
	redoStack                 undoStack
	debugUndo                 bool
}

func NewUniverseOfDiscourse() *UniverseOfDiscourse {
	var uOfD UniverseOfDiscourse
	uOfD.baseElementMap = make(map[string]BaseElement)
	uOfD.uriBaseElementMap = make(map[string]BaseElement)
	uOfD.elementListenerMap = make(map[string]elementPointerList)
	uOfD.elementPointerListenerMap = make(map[string]elementPointerPointerList)
	uOfD.literalListenerMap = make(map[string]literalPointerList)
	uOfD.literalPointerListenerMap = make(map[string]literalPointerPointerList)
	uOfD.recordingUndo = false
	uOfD.debugUndo = false
	return &uOfD
}

func (uOfDPtr *UniverseOfDiscourse) AddBaseElement(be BaseElement) error {
	//	log.Printf("Locking UofD\n")
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	if be != nil {
		be.traceableLock()
		defer be.traceableUnlock()
	}
	return uOfDPtr.addBaseElement(be)
}

func (uOfDPtr *UniverseOfDiscourse) addBaseElement(be BaseElement) error {
	if be == nil {
		return errors.New("UniverseOfDiscource addBaseElement failed because base element was nil")
	}
	//	log.Printf("Locking %T: %s \n", be, be.getId().String())
	//	log.Printf("BaseElement: %+v \n", be)
	//	log.Printf("Got the lock for %T: %s \n", be, be.getId().String())
	if be.getId() == uuid.Nil {
		return errors.New("UniverseOfDiscource addBaseElement failed because UUID was nil")
	}
	oldUOfD := be.getUniverseOfDiscourse()
	if oldUOfD != nil {
		if oldUOfD == uOfDPtr {
			return nil
		} else {
			log.Printf("Locking old UofD\n")
			oldUOfD.traceableLock()
			defer oldUOfD.traceableUnlock()
			oldUOfD.removeBaseElement(be)
		}
	}
	//	log.Printf("Adding be to UofD map")
	uOfDPtr.baseElementMap[be.getId().String()] = be
	//	log.Printf("Setting be's uOfD")
	be.setUniverseOfDiscourse(uOfDPtr)
	uOfDPtr.markNewBaseElement(be)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) addBaseElementForUndo(be BaseElement) {
	//	log.Printf("Locking UofD\n")
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	if be != nil {
		be.traceableLock()
		defer be.traceableUnlock()
	}
	uOfDPtr.baseElementMap[be.getId().String()] = be
}

func (uOfDPtr *UniverseOfDiscourse) addElementListener(element Element, elementPointer ElementPointer) {
	if element != nil {
		elementId := element.getId().String()
		currentList := uOfDPtr.elementListenerMap[elementId]
		if currentList != nil && len(*currentList) > 0 {
			for i := range *currentList {
				if (*currentList)[i] == elementPointer {
					// element is already in list
					return
				}
			}
		}
		if currentList == nil {
			var newList [1]ElementPointer
			newList[0] = elementPointer
			newSlice := newList[:]
			uOfDPtr.elementListenerMap[elementId] = &newSlice
		} else {
			updatedList := append(*currentList, elementPointer)
			uOfDPtr.elementListenerMap[elementId] = &updatedList
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) addElementPointerListener(elementPointer ElementPointer, elementPointerPointer ElementPointerPointer) {
	if elementPointer != nil {
		elementId := elementPointer.getId().String()
		currentList := uOfDPtr.elementPointerListenerMap[elementId]
		if currentList != nil && len(*currentList) > 0 {
			for i := range *currentList {
				if (*currentList)[i] == elementPointerPointer {
					// elementPointer is already in list
					return
				}
			}
		}
		if currentList == nil {
			var newList [1]ElementPointerPointer
			newList[0] = elementPointerPointer
			newSlice := newList[:]
			uOfDPtr.elementPointerListenerMap[elementId] = &newSlice
		} else {
			updatedList := append(*currentList, elementPointerPointer)
			uOfDPtr.elementPointerListenerMap[elementId] = &updatedList
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) addLiteralListener(literal Literal, literalPointer LiteralPointer) {
	if literal != nil {
		literalId := literal.getId().String()
		currentList := uOfDPtr.literalListenerMap[literalId]
		if currentList != nil && len(*currentList) > 0 {
			for i := range *currentList {
				if (*currentList)[i] == literalPointer {
					// literal is already in list
					return
				}
			}
		}
		if currentList == nil {
			var newList [1]LiteralPointer
			newList[0] = literalPointer
			newSlice := newList[:]
			uOfDPtr.literalListenerMap[literalId] = &newSlice
		} else {
			updatedList := append(*currentList, literalPointer)
			uOfDPtr.literalListenerMap[literalId] = &updatedList
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) addLiteralPointerListener(literalPointer LiteralPointer, literalPointerPointer LiteralPointerPointer) {
	if literalPointer != nil {
		literalId := literalPointer.getId().String()
		currentList := uOfDPtr.literalPointerListenerMap[literalId]
		if currentList != nil && len(*currentList) > 0 {
			for i := range *currentList {
				if (*currentList)[i] == literalPointerPointer {
					// literalPointer is already in list
					return
				}
			}
		}
		if currentList == nil {
			var newList [1]LiteralPointerPointer
			newList[0] = literalPointerPointer
			newSlice := newList[:]
			uOfDPtr.literalPointerListenerMap[literalId] = &newSlice
		} else {
			updatedList := append(*currentList, literalPointerPointer)
			uOfDPtr.literalPointerListenerMap[literalId] = &updatedList
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) getBaseElement(id string) BaseElement {
	return uOfDPtr.baseElementMap[id]
}

func (uOfDPtr *UniverseOfDiscourse) GetBaseElementWithUri(uri string) BaseElement {
	uOfDPtr.traceableRLock()
	defer uOfDPtr.traceableRUnlock()
	return uOfDPtr.getBaseElementWithUri(uri)
}

func (uOfDPtr *UniverseOfDiscourse) getBaseElementWithUri(uri string) BaseElement {
	//	return uOfDPtr.uriBaseElementMap[uri]
	// For now we just brute force it:
	for _, be := range uOfDPtr.baseElementMap {
		if be.getUri() == uri {
			return be
		}
	}
	return nil
}

func (uOfD *UniverseOfDiscourse) GetCoreConceptSpace() Element {
	return uOfD.RecoverElement([]byte(serializedCore))
}

func (uOfDPtr *UniverseOfDiscourse) GetElement(id string) Element {
	uOfDPtr.traceableRLock()
	defer uOfDPtr.traceableRUnlock()
	return uOfDPtr.getElement(id)
}

func (uOfDPtr *UniverseOfDiscourse) getElement(id string) Element {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *element:
		return be.(Element)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getElementPointer(id string) ElementPointer {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *elementPointer:
		return be.(ElementPointer)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getLiteral(id string) Literal {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *literal:
		return be.(Literal)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getLiteralPointer(id string) LiteralPointer {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *literalPointer:
		return be.(LiteralPointer)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getRefinement(id string) Refinement {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *refinement:
		return be.(Refinement)
	}
	return nil
}

// markChangedBaseElement() If undo is enabled, updates the undo stack.
// This function locks the UniverseOfDiscourse
func (uOfDPtr *UniverseOfDiscourse) markChangedBaseElement(changedElement BaseElement) {
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	if uOfDPtr.debugUndo == true {
		debug.PrintStack()
	}
	clone := clone(changedElement)
	if uOfDPtr.recordingUndo {
		uOfDPtr.undoStack.Push(NewUndoRedoStackEntry(Change, clone, changedElement))
	}
}

// markNewBaseElement() If undo is enabled, updates the undo stack.
// This function does not lock the UniverseOfDiscourse - the caller is expected to manage the locking
func (uOfDPtr *UniverseOfDiscourse) markNewBaseElement(be BaseElement) {
	if uOfDPtr.debugUndo == true {
		debug.PrintStack()
	}
	clone := clone(be)
	if uOfDPtr.recordingUndo {
		uOfDPtr.undoStack.Push(NewUndoRedoStackEntry(Creation, clone, be))
	}
}

// markNewBaseElement() If undo is enabled, updates the undo stack.
// This function does not lock the UniverseOfDiscourse - the caller is expected to manage the locking
func (uOfDPtr *UniverseOfDiscourse) markRemovedBaseElement(be BaseElement) {
	// Caller is expected to manage locking
	if uOfDPtr.debugUndo == true {
		debug.PrintStack()
	}
	clone := clone(be)
	if uOfDPtr.recordingUndo {
		uOfDPtr.undoStack.Push(NewUndoRedoStackEntry(Deletion, clone, be))
	}
}

// markUndoPoint() If undo is enabled, puts a marker on the undo stack.
// This function locks the UniverseOfDiscourse
func (uOfDPtr *UniverseOfDiscourse) markUndoPoint() {
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	if uOfDPtr.recordingUndo {
		uOfDPtr.undoStack.Push(NewUndoRedoStackEntry(Marker, nil, nil))
	}
}

// NewElement() creates an initialized Element. No locking is required since the existence of
// the element is unknown outside this routine
func (uOfD *UniverseOfDiscourse) NewElement() Element {
	var el element
	el.initializeElement()
	uOfD.AddBaseElement(&el)
	return &el
}

// NewAbstractElementPointer() creates and intitializes an elementPointer to play the role of an AbstractElementPointer
func (uOfD *UniverseOfDiscourse) NewAbstractElementPointer() ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = ABSTRACT_ELEMENT
	uOfD.AddBaseElement(&ep)
	return &ep
}

// NewRefinedElementPointer() creates and intitializes an elementPointer to play the role of an RefinedElementPointer
func (uOfD *UniverseOfDiscourse) NewRefinedElementPointer() ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = REFINED_ELEMENT
	uOfD.AddBaseElement(&ep)
	return &ep
}

// NewOwningElementPointer() creates and intitializes an elementPointer to play the role of an OwningElementPointer
func (uOfD *UniverseOfDiscourse) NewOwningElementPointer() ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = OWNING_ELEMENT
	uOfD.AddBaseElement(&ep)
	return &ep
}

// NewReferencedElementPointer() creates and intitializes an elementPointer to play the role of an ReferencedElementPointer
func (uOfD *UniverseOfDiscourse) NewReferencedElementPointer() ElementPointer {
	var ep elementPointer
	ep.initializeElementPointer()
	ep.elementPointerRole = REFERENCED_ELEMENT
	uOfD.AddBaseElement(&ep)
	return &ep
}

func (uOfD *UniverseOfDiscourse) NewElementPointerPointer() ElementPointerPointer {
	var ep elementPointerPointer
	ep.initializeElementPointerPointer()
	uOfD.AddBaseElement(&ep)
	return &ep
}

func (uOfD *UniverseOfDiscourse) NewElementPointerReference() ElementPointerReference {
	var el elementPointerReference
	el.initializeElementPointerReference()
	uOfD.AddBaseElement(&el)
	return &el
}

func (uOfD *UniverseOfDiscourse) NewElementReference() ElementReference {
	var el elementReference
	el.initializeElementReference()
	uOfD.AddBaseElement(&el)
	return &el
}

func (uOfD *UniverseOfDiscourse) NewLiteral() Literal {
	var lit literal
	lit.initializeLiteral()
	uOfD.AddBaseElement(&lit)
	return &lit
}

func (uOfD *UniverseOfDiscourse) NewNameLiteralPointer() LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = NAME
	uOfD.AddBaseElement(&lp)
	return &lp
}

func (uOfD *UniverseOfDiscourse) NewDefinitionLiteralPointer() LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = DEFINITION
	uOfD.AddBaseElement(&lp)
	return &lp
}

func (uOfD *UniverseOfDiscourse) NewUriLiteralPointer() LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = URI
	uOfD.AddBaseElement(&lp)
	return &lp
}

func (uOfD *UniverseOfDiscourse) NewValueLiteralPointer() LiteralPointer {
	var lp literalPointer
	lp.initializeLiteralPointer()
	lp.literalPointerRole = VALUE
	uOfD.AddBaseElement(&lp)
	return &lp
}

func (uOfD *UniverseOfDiscourse) NewLiteralPointerPointer() LiteralPointerPointer {
	var ep literalPointerPointer
	ep.initializeLiteralPointerPointer()
	uOfD.AddBaseElement(&ep)
	return &ep
}

func (uOfD *UniverseOfDiscourse) NewLiteralPointerReference() LiteralPointerReference {
	var el literalPointerReference
	el.initializeLiteralPointerReference()
	uOfD.AddBaseElement(&el)
	return &el
}

func (uOfD *UniverseOfDiscourse) NewLiteralReference() LiteralReference {
	var el literalReference
	el.initializeLiteralReference()
	uOfD.AddBaseElement(&el)
	return &el
}

func (uOfD *UniverseOfDiscourse) NewRefinement() Refinement {
	var el refinement
	el.initializeRefinement()
	uOfD.AddBaseElement(&el)
	return &el
}

func (uOfDPtr *UniverseOfDiscourse) notifyElementListeners(notification *ChangeNotification) {
	switch notification.changedObject.(type) {
	case Element:
		id := notification.changedObject.getId().String()
		epl := uOfDPtr.elementListenerMap[id]
		if epl != nil {
			for _, elementPointer := range *epl {
				// Must suppress circular notifications
				if notification.isReferenced(elementPointer) == false {
					propagateChange(elementPointer, notification)
				}
			}
		}
	}
}

func (uOfD *UniverseOfDiscourse) RecoverElement(data []byte) Element {
	if len(data) == 0 {
		return nil
	}
	var recoveredElement BaseElement
	err := unmarshalPolymorphicBaseElement(data, &recoveredElement)
	//	fmt.Printf("Recovered Element: \n")
	//	Print(recoveredElement, "   ")
	if err != nil {
		log.Printf("Error recovering Element: %s \n", err)
		return nil
	}
	uOfD.SetUniverseOfDiscourseRecursively(recoveredElement)
	restoreValueOwningElementFieldsRecursively(recoveredElement.(Element))
	return recoveredElement.(Element)
}

func (uOfDPtr *UniverseOfDiscourse) redo() {
	for len(uOfDPtr.redoStack) > 0 {
		currentEntry := uOfDPtr.redoStack.Pop()
		if currentEntry.changeType == Marker {
			uOfDPtr.undoStack.Push(currentEntry)
			return
		} else if currentEntry.changeType == Creation {
			uOfDPtr.undoStack.Push(currentEntry)
			uOfDPtr.restoreState(currentEntry.priorState, currentEntry.changedElement)
			// this was a new element
			uOfDPtr.addBaseElementForUndo(currentEntry.changedElement)
		} else if currentEntry.changeType == Deletion {
			uOfDPtr.undoStack.Push(currentEntry)
			uOfDPtr.restoreState(currentEntry.priorState, currentEntry.changedElement)
			// this was an deleted element
			uOfDPtr.removeBaseElementForUndo(currentEntry.changedElement)
		} else {
			clone := clone(currentEntry.changedElement)
			undoEntry := NewUndoRedoStackEntry(Change, clone, currentEntry.changedElement)
			uOfDPtr.restoreState(currentEntry.priorState, currentEntry.changedElement)
			uOfDPtr.undoStack.Push(undoEntry)
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) RemoveBaseElement(be BaseElement) error {
	//	log.Printf("Locking UofD\n")
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	if be != nil {
		be.traceableLock()
		defer be.traceableUnlock()
	}
	return uOfDPtr.removeBaseElement(be)
}

func (uOfDPtr *UniverseOfDiscourse) removeBaseElement(be BaseElement) error {
	if be == nil {
		return errors.New("UniverseOfDiscource removeBaseElement failed because base element was nil")
	}
	delete(uOfDPtr.baseElementMap, be.getId().String())
	uOfDPtr.markRemovedBaseElement(be)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) removeBaseElementForUndo(be BaseElement) {
	//	log.Printf("Locking UofD\n")
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	if be != nil {
		be.traceableLock()
		defer be.traceableUnlock()
	}
	delete(uOfDPtr.baseElementMap, be.getId().String())
}

func (uOfDPtr *UniverseOfDiscourse) removeElementListener(element Element, elementPointer ElementPointer) {
	if element != nil {
		elementId := element.getId().String()
		currentList := uOfDPtr.elementListenerMap[elementId]
		if currentList != nil && len(*currentList) > 0 {
			for i := range *currentList {
				if (*currentList)[i] == elementPointer {
					copy((*currentList)[i:], (*currentList)[i+1:])
					updatedList := (*currentList)[:len(*currentList)-1]
					uOfDPtr.elementListenerMap[elementId] = &updatedList
					return
				}
			}
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeElementPointerListener(elementPointer ElementPointer, elementPointerPointer ElementPointerPointer) {
	if elementPointer != nil {
		elementId := elementPointer.getId().String()
		currentList := uOfDPtr.elementPointerListenerMap[elementId]
		if currentList != nil && len(*currentList) > 0 {
			for i := range *currentList {
				if (*currentList)[i] == elementPointerPointer {
					copy((*currentList)[i:], (*currentList)[i+1:])
					updatedList := (*currentList)[:len(*currentList)-1]
					uOfDPtr.elementPointerListenerMap[elementId] = &updatedList
					return
				}
			}
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeLiteralListener(literal Literal, literalPointer LiteralPointer) {
	if literal != nil {
		literalId := literal.getId().String()
		currentList := uOfDPtr.literalListenerMap[literalId]
		if currentList != nil && len(*currentList) > 0 {
			for i := range *currentList {
				if (*currentList)[i] == literalPointer {
					copy((*currentList)[i:], (*currentList)[i+1:])
					updatedList := (*currentList)[:len(*currentList)-1]
					uOfDPtr.literalListenerMap[literalId] = &updatedList
					return
				}
			}
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeLiteralPointerListener(literalPointer LiteralPointer, literalPointerPointer LiteralPointerPointer) {
	if literalPointer != nil {
		elementId := literalPointer.getId().String()
		currentList := uOfDPtr.literalPointerListenerMap[elementId]
		if currentList != nil && len(*currentList) > 0 {
			for i := range *currentList {
				if (*currentList)[i] == literalPointerPointer {
					copy((*currentList)[i:], (*currentList)[i+1:])
					updatedList := (*currentList)[:len(*currentList)-1]
					uOfDPtr.literalPointerListenerMap[elementId] = &updatedList
					return
				}
			}
		}
	}
}

// restoreState is used as part of the undo process. It changes the currentState object
// to have the priorState.
func (uOfDPtr *UniverseOfDiscourse) restoreState(priorState BaseElement, currentState BaseElement) {
	if uOfDPtr.debugUndo == true {
		log.Printf("Restoring State")
		log.Printf("   Current state:")
		Print(currentState, "      ")
		log.Printf("   Prior state")
		Print(priorState, "      ")
	}
	switch currentState.(type) {
	case *element:
		currentState.(*element).cloneAttributes(*priorState.(*element))
	case *elementPointer:
		currentState.(*elementPointer).cloneAttributes(*priorState.(*elementPointer))
	case *elementPointerPointer:
		currentState.(*elementPointerPointer).cloneAttributes(*priorState.(*elementPointerPointer))
	case *elementPointerReference:
		currentState.(*elementPointerReference).cloneAttributes(*priorState.(*elementPointerReference))
	case *elementReference:
		currentState.(*elementReference).cloneAttributes(*priorState.(*elementReference))
	case *literal:
		currentState.(*literal).cloneAttributes(*priorState.(*literal))
	case *literalPointer:
		currentState.(*literalPointer).cloneAttributes(*priorState.(*literalPointer))
	case *literalPointerPointer:
		currentState.(*literalPointerPointer).cloneAttributes(*priorState.(*literalPointerPointer))
	case *literalPointerReference:
		currentState.(*literalPointerReference).cloneAttributes(*priorState.(*literalPointerReference))
	case *literalReference:
		currentState.(*literalReference).cloneAttributes(*priorState.(*literalReference))
	case *refinement:
		currentState.(*refinement).cloneAttributes(*priorState.(*refinement))
	default:
		log.Printf("restoreState called with unhandled type %T\n", currentState)
	}
}

func (uOfDPtr *UniverseOfDiscourse) SetRecordingUndo(newSetting bool) {
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	uOfDPtr.setRecordingUndo(newSetting)
}

func (uOfDPtr *UniverseOfDiscourse) setRecordingUndo(newSetting bool) {
	uOfDPtr.recordingUndo = newSetting
}

func (uOfDPtr *UniverseOfDiscourse) SetUniverseOfDiscourseRecursively(be BaseElement) {
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	uOfDPtr.setUniverseOfDiscourseRecursively(be)
}

func (uOfDPtr *UniverseOfDiscourse) setUniverseOfDiscourseRecursively(be BaseElement) {
	uOfDPtr.addBaseElement(be)
	switch be.(type) {
	case *element:
		for _, child := range be.(*element).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *elementPointerReference:
		for _, child := range be.(*elementPointerReference).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *elementReference:
		for _, child := range be.(*elementReference).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *literalPointerReference:
		for _, child := range be.(*literalPointerReference).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *literalReference:
		for _, child := range be.(*literalReference).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *refinement:
		for _, child := range be.(*refinement).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *elementPointer, *elementPointerPointer, *literal, *literalPointer, *literalPointerPointer:
	// Do nothing
	default:
		log.Printf("UniverseOfDiscourse.setUniverseOfDiscourseRecursively is missing case for %T\n", be)
	}
}

func (uOfDPtr *UniverseOfDiscourse) traceableLock() {
	if TraceLocks {
		log.Printf("About to lock Universe of Discourse %p\n", uOfDPtr)
	}
	uOfDPtr.Lock()
}

func (uOfDPtr *UniverseOfDiscourse) traceableRLock() {
	if TraceLocks {
		log.Printf("About to lock Universe of Discourse %p\n", uOfDPtr)
	}
	uOfDPtr.RLock()
}

func (uOfDPtr *UniverseOfDiscourse) traceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock Universe of Discourse %p\n", uOfDPtr)
	}
	uOfDPtr.Unlock()
}

func (uOfDPtr *UniverseOfDiscourse) traceableRUnlock() {
	if TraceLocks {
		log.Printf("About to unlock Universe of Discourse %p\n", uOfDPtr)
	}
	uOfDPtr.RUnlock()
}

func (uOfDPtr *UniverseOfDiscourse) undo() {
	firstEntry := true
	for len(uOfDPtr.undoStack) > 0 {
		currentEntry := uOfDPtr.undoStack.Pop()
		if currentEntry.changeType == Marker {
			if firstEntry {
				uOfDPtr.redoStack.Push(currentEntry)
			} else {
				// Put it back on the undo stack
				uOfDPtr.undoStack.Push(currentEntry)
				return
			}
		} else if currentEntry.changeType == Creation {
			uOfDPtr.redoStack.Push(currentEntry)
			uOfDPtr.removeBaseElementForUndo(currentEntry.changedElement)
		} else if currentEntry.changeType == Deletion {
			uOfDPtr.restoreState(currentEntry.priorState, currentEntry.changedElement)
			uOfDPtr.redoStack.Push(currentEntry)
			uOfDPtr.addBaseElementForUndo(currentEntry.changedElement)
		} else if currentEntry.changeType == Change {
			clone := clone(currentEntry.changedElement)
			redoEntry := NewUndoRedoStackEntry(Change, clone, currentEntry.changedElement)
			uOfDPtr.restoreState(currentEntry.priorState, currentEntry.changedElement)
			uOfDPtr.redoStack.Push(redoEntry)
		}
		firstEntry = false
	}
}
