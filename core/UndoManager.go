// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"log"
	"runtime/debug"
	"sync"
)

type undoManager struct {
	sync.Mutex
	debugUndo     bool
	recordingUndo bool
	redoStack     undoStack
	undoStack     undoStack
	uOfD          *UniverseOfDiscourse
}

// NewUndoManager creates and initializes the manager for the undo/redo functionality
func newUndoManager(uOfD *UniverseOfDiscourse) *undoManager {
	var undoMgr undoManager
	undoMgr.debugUndo = false
	undoMgr.recordingUndo = false
	undoMgr.uOfD = uOfD
	return &undoMgr
}

// markChangedElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markChangedElement(changedElement Element, hl *Transaction) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	hl.ReadLockElement(changedElement)
	if undoMgr.debugUndo {
		debug.PrintStack()
	}
	priorState := clone(changedElement, hl)
	priorOwnedElements := undoMgr.uOfD.ownedIDsMap.GetMappedValues(changedElement.GetConceptID(hl)).Clone()
	priorListeners := undoMgr.uOfD.listenersMap.GetMappedValues(changedElement.GetConceptID(hl)).Clone()
	if undoMgr.recordingUndo {
		undoMgr.undoStack.Push(newUndoRedoStackEntry(Change, priorState, priorOwnedElements, priorListeners, changedElement))
	}
	return nil
}

// markNewElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markNewElement(el Element, hl *Transaction) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if hl == nil {
		return errors.New("UndoManager.markNewElement called with nil HeldLocks")
	}
	hl.ReadLockElement(el)
	if undoMgr.debugUndo {
		debug.PrintStack()
	}
	if undoMgr.recordingUndo {
		clone := clone(el, hl)
		priorOwnedElements := undoMgr.uOfD.ownedIDsMap.GetMappedValues(el.GetConceptID(hl)).Clone()
		priorListeners := undoMgr.uOfD.listenersMap.GetMappedValues(el.GetConceptID(hl)).Clone()
		stackEntry := newUndoRedoStackEntry(Creation, clone, priorOwnedElements, priorListeners, el)
		if undoMgr.debugUndo {
			PrintStackEntry(stackEntry, hl)
		}
		undoMgr.undoStack.Push(stackEntry)
	}
	return nil
}

// markRemoveElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markRemovedElement(el Element, hl *Transaction) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if hl == nil {
		return errors.New("UndoManager.markRemovedElement called with nil HeldLocks")
	}
	hl.ReadLockElement(el)
	if undoMgr.debugUndo {
		debug.PrintStack()
	}
	if undoMgr.recordingUndo {
		clone := clone(el, hl)
		priorOwnedElements := undoMgr.uOfD.ownedIDsMap.GetMappedValues(el.GetConceptID(hl)).Clone()
		priorListeners := undoMgr.uOfD.listenersMap.GetMappedValues(el.GetConceptID(hl)).Clone()
		undoMgr.undoStack.Push(newUndoRedoStackEntry(Deletion, clone, priorOwnedElements, priorListeners, el))
	}
	return nil
}

// MarkUndoPoint() If undo is enabled, puts a marker on the undo stack.
func (undoMgr *undoManager) MarkUndoPoint() {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if undoMgr.recordingUndo {
		undoMgr.undoStack.Push(newUndoRedoStackEntry(Marker, nil, nil, nil, nil))
	}
}

// PrintUndoStack prints the undo stack. It is intended only for debugging.
func PrintUndoStack(s undoStack, stackName string, uOfD *UniverseOfDiscourse) {
	hl := uOfD.NewTransaction()
	defer hl.ReleaseLocksAndWait()
	log.Printf("%s:", stackName)
	for _, entry := range s {
		var changeType string
		switch entry.changeType {
		case Creation:
			{
				changeType = "Creation"
			}
		case Deletion:
			{
				changeType = "Deletion"
			}
		case Change:
			{
				changeType = "Change"
			}
		case Marker:
			{
				changeType = "Marker"
			}
		}
		log.Printf("   Change type: %s", changeType)
		log.Printf("   Prior state:")
		Print(entry.priorState, "      ", hl)
		log.Printf("   Changed element:")
		Print(entry.changedElement, "      ", hl)
	}
}

// PrintStackEntry prints the entry on the stack. It is intended for use only in debugging
func PrintStackEntry(entry *undoRedoStackEntry, hl *Transaction) {
	var changeType string
	switch entry.changeType {
	case Creation:
		{
			changeType = "Creation"
		}
	case Deletion:
		{
			changeType = "Deletion"
		}
	case Change:
		{
			changeType = "Change"
		}
	case Marker:
		{
			changeType = "Marker"
		}
	}
	log.Printf("   Change type: %s", changeType)
	log.Printf("   Prior state:")
	Print(entry.priorState, "      ", hl)
	log.Printf("   Changed element:")
	Print(entry.changedElement, "      ", hl)
}

func (undoMgr *undoManager) redo(hl *Transaction) {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	uOfD := undoMgr.uOfD
	for len(undoMgr.redoStack) > 0 {
		currentEntry := undoMgr.redoStack.Pop()
		var currentID string
		if currentEntry.changedElement != nil {
			currentID = currentEntry.changedElement.GetConceptID(hl)
		}
		if currentEntry.changeType == Marker {
			undoMgr.undoStack.Push(currentEntry)
			return
		} else if currentEntry.changeType == Creation {
			// Update owner's owned elements. Since we are redoing a creation, if the changed element has an owner
			// we must add the current element to the owner's ownedElements
			if currentEntry.changedElement.GetOwningConceptID(hl) != "" {
				uOfD.ownedIDsMap.AddMappedValue(currentEntry.changedElement.GetOwningConceptID(hl), currentID)
			}
			// Update listeners. If this is a reference or refinement pointing to another element, add this element to the other element's listener's set
			switch currentEntry.changedElement.(type) {
			case *reference:
				referencedElementID := currentEntry.changedElement.(*reference).ReferencedConceptID
				if referencedElementID != "" {
					uOfD.listenersMap.AddMappedValue(referencedElementID, currentID)
				}
			case *refinement:
				abstractID := currentEntry.changedElement.(*refinement).AbstractConceptID
				if abstractID != "" {
					uOfD.listenersMap.AddMappedValue(abstractID, currentID)
				}
				refinedID := currentEntry.changedElement.(*refinement).RefinedConceptID
				if refinedID != "" {
					uOfD.listenersMap.AddMappedValue(refinedID, currentID)
				}
			}
			// Update the uriUUIDMap
			uri := currentEntry.changedElement.GetURI(hl)
			if uri != "" {
				uOfD.uriUUIDMap.SetEntry(uri, currentID)
			}
			undoMgr.undoStack.Push(currentEntry)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			// this was a new element
			uOfD.addElementForUndo(currentEntry.changedElement, hl)
			uOfD.ownedIDsMap.SetMappedValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
		} else if currentEntry.changeType == Deletion {
			// Update owner's owned elements. Since we are redoing a deletion, if the changed element has an owner
			// we must remove the current element from the owner's ownedElements
			if currentEntry.priorState.GetOwningConceptID(hl) != "" {
				uOfD.ownedIDsMap.RemoveMappedValue(currentEntry.priorState.GetOwningConceptID(hl), currentEntry.priorState.GetConceptID(hl))
			}
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			switch currentEntry.priorState.(type) {
			case *reference:
				referencedElementID := currentEntry.priorState.(*reference).ReferencedConceptID
				if referencedElementID != "" {
					uOfD.listenersMap.RemoveMappedValue(referencedElementID, currentEntry.priorState.GetConceptID(hl))
				}
			case *refinement:
				abstractID := currentEntry.priorState.(*refinement).AbstractConceptID
				if abstractID != "" {
					uOfD.listenersMap.RemoveMappedValue(abstractID, currentEntry.priorState.GetConceptID(hl))
				}
				refinedID := currentEntry.priorState.(*refinement).RefinedConceptID
				if refinedID != "" {
					uOfD.listenersMap.RemoveMappedValue(refinedID, currentEntry.priorState.GetConceptID(hl))
				}
			}
			// Update the uriUUIDMap
			uri := currentEntry.priorState.GetURI(hl)
			if uri != "" {
				uOfD.uriUUIDMap.DeleteEntry(uri)
			}
			undoMgr.undoStack.Push(currentEntry)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			// this was an deleted element
			uOfD.removeElementForUndo(currentEntry.changedElement, hl)
			uOfD.ownedIDsMap.SetMappedValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
		} else if currentEntry.changeType == Change {
			// If the owner changes, update owner's owned elements.
			currentOwnerID := currentEntry.changedElement.GetOwningConceptID(hl)
			priorOwnerID := currentEntry.priorState.GetOwningConceptID(hl)
			if currentOwnerID != priorOwnerID {
				if currentOwnerID != "" {
					uOfD.ownedIDsMap.AddMappedValue(currentOwnerID, currentID)
				}
				if priorOwnerID != "" {
					uOfD.ownedIDsMap.RemoveMappedValue(priorOwnerID, currentID)
				}
			}
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			switch currentEntry.changedElement.(type) {
			case *reference:
				currentReferencedElementID := currentEntry.changedElement.(*reference).ReferencedConceptID
				priorReferencedElementID := currentEntry.priorState.(*reference).ReferencedConceptID
				if currentReferencedElementID != priorReferencedElementID {
					if currentReferencedElementID != "" {
						uOfD.listenersMap.RemoveMappedValue(currentReferencedElementID, currentID)
					}
					if priorReferencedElementID != "" {
						uOfD.listenersMap.AddMappedValue(priorReferencedElementID, currentID)
					}
				}
			case *refinement:
				currentAbstractID := currentEntry.changedElement.(*refinement).AbstractConceptID
				priorAbstractID := currentEntry.priorState.(*refinement).AbstractConceptID
				if currentAbstractID != priorAbstractID {
					if currentAbstractID != "" {
						uOfD.listenersMap.RemoveMappedValue(currentAbstractID, currentID)
					}
					if priorAbstractID != "" {
						uOfD.listenersMap.AddMappedValue(priorAbstractID, currentID)
					}
				}
				currentRefinedID := currentEntry.changedElement.(*refinement).RefinedConceptID
				priorRefinedID := currentEntry.priorState.(*refinement).RefinedConceptID
				if currentRefinedID != priorRefinedID {
					if currentRefinedID != "" {
						uOfD.listenersMap.RemoveMappedValue(currentRefinedID, currentID)
					}
					if priorRefinedID != "" {
						uOfD.listenersMap.AddMappedValue(priorRefinedID, currentID)
					}
				}
			}
			// Update the uriUUIDMap
			currentURI := currentEntry.changedElement.GetURI(hl)
			priorURI := currentEntry.priorState.GetURI(hl)
			if currentURI != priorURI {
				if currentURI != "" {
					uOfD.uriUUIDMap.DeleteEntry(currentURI)
				}
				if priorURI != "" {
					uOfD.uriUUIDMap.SetEntry(priorURI, currentID)
				}
			}
			clone := clone(currentEntry.changedElement, hl)
			priorOwnedElements := undoMgr.uOfD.ownedIDsMap.GetMappedValues(currentID).Clone()
			priorListeners := undoMgr.uOfD.listenersMap.GetMappedValues(currentID).Clone()
			undoEntry := newUndoRedoStackEntry(Change, clone, priorOwnedElements, priorListeners, currentEntry.changedElement)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			uOfD.ownedIDsMap.SetMappedValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
			undoMgr.undoStack.Push(undoEntry)
		}
	}
}

// restoreState is used as part of the undo process. It changes the currentState object
// to have the priorState.
func (undoMgr *undoManager) restoreState(priorState Element, currentState Element, hl *Transaction) {
	if undoMgr.debugUndo {
		log.Printf("Restoring State")
		log.Printf("   Current state:")
		Print(currentState, "      ", hl)
		log.Printf("   Prior state")
		Print(priorState, "      ", hl)
	}
	switch currentState.(type) {
	case *element:
		currentState.(*element).cloneAttributes(priorState.(*element), hl)
	case *reference:
		currentState.(*reference).cloneAttributes(priorState.(*reference), hl)
	case *literal:
		currentState.(*literal).cloneAttributes(priorState.(*literal), hl)
	case *refinement:
		currentState.(*refinement).cloneAttributes(priorState.(*refinement), hl)
	default:
		log.Printf("restoreState called with unhandled type %T\n", currentState)
	}
}

// func (undoMgr *undoManager) setDebugUndo(newSetting bool) {
// 	undoMgr.TraceableLock()
// 	defer undoMgr.TraceableUnlock()
// 	undoMgr.debugUndo = newSetting
// }

func (undoMgr *undoManager) setRecordingUndo(newSetting bool) {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	undoMgr.recordingUndo = newSetting
}

func (undoMgr *undoManager) TraceableLock() {
	// if TraceLocks {
	// 	log.Printf("About to lock Undo Manager %p\n", undoMgr)
	// }
	undoMgr.Lock()
}

func (undoMgr *undoManager) TraceableUnlock() {
	// if TraceLocks {
	// 	log.Printf("About to unlock Undo Manager %p\n", undoMgr)
	// }
	undoMgr.Unlock()
}

func (undoMgr *undoManager) undo(hl *Transaction) {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	uOfD := undoMgr.uOfD
	firstEntry := true
	for len(undoMgr.undoStack) > 0 {
		currentEntry := undoMgr.undoStack.Pop()
		var currentID string
		if currentEntry.changedElement != nil {
			currentID = currentEntry.changedElement.GetConceptID(hl)
		}
		if currentEntry.changeType == Marker {
			if firstEntry {
				undoMgr.redoStack.Push(currentEntry)
			} else {
				// Put it back on the undo stack
				undoMgr.undoStack.Push(currentEntry)
				return
			}
		} else if currentEntry.changeType == Creation {
			// Update owner's owned elements. Since we are undoing a creation, if the changed element has an owner
			// we must remove the current element from the owner's ownedElements
			if currentEntry.changedElement.GetOwningConceptID(hl) != "" {
				uOfD.ownedIDsMap.RemoveMappedValue(currentEntry.changedElement.GetOwningConceptID(hl), currentID)
			}
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			switch currentEntry.changedElement.(type) {
			case *reference:
				referencedElementID := currentEntry.changedElement.(*reference).ReferencedConceptID
				if referencedElementID != "" {
					uOfD.listenersMap.RemoveMappedValue(referencedElementID, currentID)
				}
			case *refinement:
				abstractID := currentEntry.changedElement.(*refinement).AbstractConceptID
				if abstractID != "" {
					uOfD.listenersMap.RemoveMappedValue(abstractID, currentID)
				}
				refinedID := currentEntry.changedElement.(*refinement).RefinedConceptID
				if refinedID != "" {
					uOfD.listenersMap.RemoveMappedValue(refinedID, currentID)
				}
			}
			// Update the uriUUIDMap
			uri := currentEntry.changedElement.GetURI(hl)
			if uri != "" {
				uOfD.uriUUIDMap.DeleteEntry(uri)
			}
			undoMgr.redoStack.Push(currentEntry)
			uOfD.removeElementForUndo(currentEntry.changedElement, hl)
			uOfD.ownedIDsMap.SetMappedValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
		} else if currentEntry.changeType == Deletion {
			// Update owner's owned elements. Since we are undoing a deletion, if the changed element has an owner
			// we must add the current element from the owner's ownedElements
			if currentEntry.priorState.GetOwningConceptID(hl) != "" {
				uOfD.ownedIDsMap.AddMappedValue(currentEntry.priorState.GetOwningConceptID(hl), currentEntry.priorState.GetConceptID(hl))
			}
			// Update listeners. If this is a reference or refinement pointing to another element, add this element to the other element's listener's set
			switch currentEntry.priorState.(type) {
			case *reference:
				referencedElementID := currentEntry.priorState.(*reference).ReferencedConceptID
				if referencedElementID != "" {
					uOfD.listenersMap.AddMappedValue(referencedElementID, currentEntry.priorState.GetConceptID(hl))
				}
			case *refinement:
				abstractID := currentEntry.priorState.(*refinement).AbstractConceptID
				if abstractID != "" {
					uOfD.listenersMap.AddMappedValue(abstractID, currentEntry.priorState.GetConceptID(hl))
				}
				refinedID := currentEntry.priorState.(*refinement).RefinedConceptID
				if refinedID != "" {
					uOfD.listenersMap.AddMappedValue(refinedID, currentEntry.priorState.GetConceptID(hl))
				}
			}
			// Update the uriUUIDMap
			uri := currentEntry.priorState.GetURI(hl)
			if uri != "" {
				uOfD.uriUUIDMap.SetEntry(uri, currentID)
			}
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			undoMgr.redoStack.Push(currentEntry)
			uOfD.addElementForUndo(currentEntry.changedElement, hl)
			uOfD.ownedIDsMap.SetMappedValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
		} else if currentEntry.changeType == Change {
			// If the owner changes, update owner's owned elements.
			currentOwnerID := currentEntry.changedElement.GetOwningConceptID(hl)
			priorOwnerID := currentEntry.priorState.GetOwningConceptID(hl)
			if currentOwnerID != priorOwnerID {
				if currentOwnerID != "" {
					uOfD.ownedIDsMap.RemoveMappedValue(currentOwnerID, currentID)
				}
				if priorOwnerID != "" {
					uOfD.ownedIDsMap.AddMappedValue(priorOwnerID, currentID)
				}
			}
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			switch currentEntry.changedElement.(type) {
			case *reference:
				currentReferencedElementID := currentEntry.changedElement.(*reference).ReferencedConceptID
				priorReferencedElementID := currentEntry.priorState.(*reference).ReferencedConceptID
				if currentReferencedElementID != priorReferencedElementID {
					if currentReferencedElementID != "" {
						uOfD.listenersMap.RemoveMappedValue(currentReferencedElementID, currentID)
					}
					if priorReferencedElementID != "" {
						uOfD.listenersMap.AddMappedValue(priorReferencedElementID, currentID)
					}
				}
			case *refinement:
				currentAbstractID := currentEntry.changedElement.(*refinement).AbstractConceptID
				priorAbstractID := currentEntry.priorState.(*refinement).AbstractConceptID
				if currentAbstractID != priorAbstractID {
					if currentAbstractID != "" {
						uOfD.listenersMap.RemoveMappedValue(currentAbstractID, currentID)
					}
					if priorAbstractID != "" {
						uOfD.listenersMap.AddMappedValue(priorAbstractID, currentID)
					}
				}
				currentRefinedID := currentEntry.changedElement.(*refinement).RefinedConceptID
				priorRefinedID := currentEntry.priorState.(*refinement).RefinedConceptID
				if currentRefinedID != priorRefinedID {
					if currentRefinedID != "" {
						uOfD.listenersMap.RemoveMappedValue(currentRefinedID, currentID)
					}
					if priorRefinedID != "" {
						uOfD.listenersMap.AddMappedValue(priorRefinedID, currentID)
					}
				}
			}
			// Update the uriUUIDMap
			currentURI := currentEntry.changedElement.GetURI(hl)
			priorURI := currentEntry.priorState.GetURI(hl)
			if currentURI != priorURI {
				if currentURI != "" {
					uOfD.uriUUIDMap.DeleteEntry(currentURI)
				}
				if priorURI != "" {
					uOfD.uriUUIDMap.SetEntry(priorURI, currentID)
				}
			}
			clone := clone(currentEntry.changedElement, hl)
			priorOwnedElements := uOfD.ownedIDsMap.GetMappedValues(currentID).Clone()
			priorListeners := uOfD.listenersMap.GetMappedValues(currentID).Clone()
			redoEntry := newUndoRedoStackEntry(Change, clone, priorOwnedElements, priorListeners, currentEntry.changedElement)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			uOfD.ownedIDsMap.SetMappedValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
			undoMgr.redoStack.Push(redoEntry)
		}
		firstEntry = false
	}
}
