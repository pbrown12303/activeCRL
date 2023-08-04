// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"log"
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
func (undoMgr *undoManager) markChangedElement(changedElement Element, trans *Transaction) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	trans.ReadLockElement(changedElement)
	priorState := clone(changedElement, trans)
	priorOwnedElements := undoMgr.uOfD.ownedIDsMap.GetMappedValues(changedElement.GetConceptID(trans)).Clone()
	priorListeners := undoMgr.uOfD.listenersMap.GetMappedValues(changedElement.GetConceptID(trans)).Clone()
	priorUofD := ""
	if changedElement.getUniverseOfDiscourseNoLock() != nil {
		priorUofD = changedElement.getUniverseOfDiscourseNoLock().id
	}
	if undoMgr.recordingUndo {
		stackEntry := newUndoRedoStackEntry(Change, priorState, priorOwnedElements, priorListeners, priorUofD, changedElement)
		if undoMgr.debugUndo {
			PrintStackEntry(stackEntry, trans)
		}
		undoMgr.undoStack.Push(stackEntry)
	}
	return nil
}

// markNewElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markNewElement(el Element, trans *Transaction) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if trans == nil {
		return errors.New("UndoManager.markNewElement called with nil HeldLocks")
	}
	trans.ReadLockElement(el)
	if undoMgr.recordingUndo {
		clone := clone(el, trans)
		priorOwnedElements := undoMgr.uOfD.ownedIDsMap.GetMappedValues(el.GetConceptID(trans)).Clone()
		priorListeners := undoMgr.uOfD.listenersMap.GetMappedValues(el.GetConceptID(trans)).Clone()
		priorUofD := ""
		if el.getUniverseOfDiscourseNoLock() != nil {
			priorUofD = el.getUniverseOfDiscourseNoLock().id
		}
		stackEntry := newUndoRedoStackEntry(Creation, clone, priorOwnedElements, priorListeners, priorUofD, el)
		if undoMgr.debugUndo {
			PrintStackEntry(stackEntry, trans)
		}
		undoMgr.undoStack.Push(stackEntry)
	}
	return nil
}

// markRemoveElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markRemovedElement(el Element, trans *Transaction) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if trans == nil {
		return errors.New("UndoManager.markRemovedElement called with nil HeldLocks")
	}
	trans.ReadLockElement(el)
	if undoMgr.recordingUndo {
		clone := clone(el, trans)
		priorOwnedElements := undoMgr.uOfD.ownedIDsMap.GetMappedValues(el.GetConceptID(trans)).Clone()
		priorListeners := undoMgr.uOfD.listenersMap.GetMappedValues(el.GetConceptID(trans)).Clone()
		priorUofD := ""
		if el.getUniverseOfDiscourseNoLock() != nil {
			priorUofD = el.getUniverseOfDiscourseNoLock().id
		}
		stackEntry := newUndoRedoStackEntry(Deletion, clone, priorOwnedElements, priorListeners, priorUofD, el)
		if undoMgr.debugUndo {
			PrintStackEntry(stackEntry, trans)
		}
		undoMgr.undoStack.Push(stackEntry)
	}
	return nil
}

// MarkUndoPoint() If undo is enabled, puts a marker on the undo stack.
func (undoMgr *undoManager) MarkUndoPoint() {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if undoMgr.recordingUndo {
		if undoMgr.debugUndo {
			log.Print("***** MARK UNDO POINT *****")
		}
		undoMgr.undoStack.Push(newUndoRedoStackEntry(Marker, nil, nil, nil, "", nil))
	}
}

// PrintUndoStack prints the undo stack. It is intended only for debugging.
func PrintUndoStack(s undoStack, stackName string, uOfD *UniverseOfDiscourse) {
	trans := uOfD.NewTransaction()
	defer trans.ReleaseLocks()
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
		log.Printf("Change type: %s", changeType)
		log.Printf("   Prior state:")
		Print(entry.priorState, "      ", trans)
		log.Printf("   Changed element:")
		log.Printf("   Prior UofD: %s", entry.priorUofD)
		Print(entry.changedElement, "      ", trans)
	}
}

// PrintStackEntry prints the entry on the stack. It is intended for use only in debugging
func PrintStackEntry(entry *undoRedoStackEntry, trans *Transaction) {
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
	log.Printf("Change type: %s", changeType)
	log.Printf("   Prior state:")
	Print(entry.priorState, "      ", trans)
	log.Printf("   Changed element:")
	Print(entry.changedElement, "      ", trans)
	log.Printf("   Prior UofD: %s", entry.priorUofD)
	log.Printf("   Undo/Redo Stack Entry priorOwnedElements: %v", trans.uOfD.ownedIDsMap.GetMappedValues(entry.changedElement.getConceptIDNoLock()))

}

func (undoMgr *undoManager) redo(trans *Transaction) {
	if undoMgr.debugUndo {
		log.Print("***** BEGIN REDO ****")
	}
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	uOfD := undoMgr.uOfD
	for len(undoMgr.redoStack) > 0 {
		currentEntry := undoMgr.redoStack.Pop()
		var currentID string
		if currentEntry.changedElement != nil {
			currentID = currentEntry.changedElement.GetConceptID(trans)
		}
		if currentEntry.changeType == Marker {
			undoMgr.undoStack.Push(currentEntry)
			return
		} else if currentEntry.changeType == Creation {
			// Update listeners. If this is a reference or refinement pointing to another element, add this element to the other element's listener's set
			switch currentEntry.changedElement.(type) {
			case *reference:
				referencedElementID := currentEntry.changedElement.(*reference).ReferencedConceptID
				if referencedElementID != "" {
					uOfD.listenersMap.addMappedValue(referencedElementID, currentID)
				}
			case *refinement:
				abstractID := currentEntry.changedElement.(*refinement).AbstractConceptID
				if abstractID != "" {
					uOfD.listenersMap.addMappedValue(abstractID, currentID)
				}
				refinedID := currentEntry.changedElement.(*refinement).RefinedConceptID
				if refinedID != "" {
					uOfD.listenersMap.addMappedValue(refinedID, currentID)
				}
			}
			// Update the uriUUIDMap
			uri := currentEntry.changedElement.GetURI(trans)
			if uri != "" {
				uOfD.uriUUIDMap.SetEntry(uri, currentID)
			}
			undoMgr.undoStack.Push(currentEntry)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, trans)
			// this was a new element
			uOfD.addElementForUndo(currentEntry.changedElement, trans)
			uOfD.setOwnedIDsMapValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
			uOfD.setUUIDElementMapEntry(currentEntry.changedElement.getConceptIDNoLock(), currentEntry.changedElement)
		} else if currentEntry.changeType == Deletion {
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			switch currentEntry.priorState.(type) {
			case *reference:
				referencedElementID := currentEntry.priorState.(*reference).ReferencedConceptID
				if referencedElementID != "" {
					uOfD.listenersMap.removeMappedValue(referencedElementID, currentEntry.priorState.GetConceptID(trans))
				}
			case *refinement:
				abstractID := currentEntry.priorState.(*refinement).AbstractConceptID
				if abstractID != "" {
					uOfD.listenersMap.removeMappedValue(abstractID, currentEntry.priorState.GetConceptID(trans))
				}
				refinedID := currentEntry.priorState.(*refinement).RefinedConceptID
				if refinedID != "" {
					uOfD.listenersMap.removeMappedValue(refinedID, currentEntry.priorState.GetConceptID(trans))
				}
			}
			// Update the uriUUIDMap
			uri := currentEntry.priorState.GetURI(trans)
			if uri != "" {
				uOfD.uriUUIDMap.DeleteEntry(uri)
			}
			undoMgr.undoStack.Push(currentEntry)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, trans)
			// this was an deleted element
			uOfD.removeElementForUndo(currentEntry.changedElement, trans)
			uOfD.setOwnedIDsMapValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
			uOfD.deleteUUIDElementMapEntry(currentEntry.changedElement.getConceptIDNoLock())
		} else if currentEntry.changeType == Change {
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			switch currentEntry.changedElement.(type) {
			case *reference:
				currentReferencedElementID := currentEntry.changedElement.(*reference).ReferencedConceptID
				priorReferencedElementID := currentEntry.priorState.(*reference).ReferencedConceptID
				if currentReferencedElementID != priorReferencedElementID {
					if currentReferencedElementID != "" {
						uOfD.listenersMap.removeMappedValue(currentReferencedElementID, currentID)
					}
					if priorReferencedElementID != "" {
						uOfD.listenersMap.addMappedValue(priorReferencedElementID, currentID)
					}
				}
			case *refinement:
				currentAbstractID := currentEntry.changedElement.(*refinement).AbstractConceptID
				priorAbstractID := currentEntry.priorState.(*refinement).AbstractConceptID
				if currentAbstractID != priorAbstractID {
					if currentAbstractID != "" {
						uOfD.listenersMap.removeMappedValue(currentAbstractID, currentID)
					}
					if priorAbstractID != "" {
						uOfD.listenersMap.addMappedValue(priorAbstractID, currentID)
					}
				}
				currentRefinedID := currentEntry.changedElement.(*refinement).RefinedConceptID
				priorRefinedID := currentEntry.priorState.(*refinement).RefinedConceptID
				if currentRefinedID != priorRefinedID {
					if currentRefinedID != "" {
						uOfD.listenersMap.removeMappedValue(currentRefinedID, currentID)
					}
					if priorRefinedID != "" {
						uOfD.listenersMap.addMappedValue(priorRefinedID, currentID)
					}
				}
			}
			// Update the uriUUIDMap
			currentURI := currentEntry.changedElement.GetURI(trans)
			priorURI := currentEntry.priorState.GetURI(trans)
			if currentURI != priorURI {
				if currentURI != "" {
					uOfD.uriUUIDMap.DeleteEntry(currentURI)
				}
				if priorURI != "" {
					uOfD.uriUUIDMap.SetEntry(priorURI, currentID)
				}
			}
			clone := clone(currentEntry.changedElement, trans)
			priorOwnedElements := undoMgr.uOfD.ownedIDsMap.GetMappedValues(currentID).Clone()
			priorListeners := undoMgr.uOfD.listenersMap.GetMappedValues(currentID).Clone()
			undoEntry := newUndoRedoStackEntry(Change, clone, priorOwnedElements, priorListeners, currentEntry.priorUofD, currentEntry.changedElement)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, trans)
			uOfD.setOwnedIDsMapValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
			if currentEntry.priorUofD != uOfD.id {
				uOfD.deleteUUIDElementMapEntry(currentID)
			} else if currentEntry.priorUofD == uOfD.id {
				uOfD.setUUIDElementMapEntry(currentID, currentEntry.changedElement)
			}
			undoMgr.undoStack.Push(undoEntry)
		}
	}
}

// restoreState is used as part of the undo process. It changes the currentState object
// to have the priorState.
func (undoMgr *undoManager) restoreState(priorState Element, currentState Element, trans *Transaction) {
	if undoMgr.debugUndo {
		log.Printf("Restoring State")
		log.Printf("   Current state:")
		Print(currentState, "      ", trans)
		log.Printf("   Prior state")
		Print(priorState, "      ", trans)
	}
	switch currentState.(type) {
	case *element:
		currentState.(*element).cloneAttributes(priorState.(*element), trans)
	case *reference:
		currentState.(*reference).cloneAttributes(priorState.(*reference), trans)
	case *literal:
		currentState.(*literal).cloneAttributes(priorState.(*literal), trans)
	case *refinement:
		currentState.(*refinement).cloneAttributes(priorState.(*refinement), trans)
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

func (undoMgr *undoManager) undo(trans *Transaction) {
	if undoMgr.debugUndo {
		log.Print("***** BEGIN UNDO ****")
	}
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	uOfD := undoMgr.uOfD
	firstEntry := true
	for len(undoMgr.undoStack) > 0 {
		currentEntry := undoMgr.undoStack.Pop()
		var currentID string
		if currentEntry.changedElement != nil {
			currentID = currentEntry.changedElement.GetConceptID(trans)
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
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			switch currentEntry.changedElement.(type) {
			case *reference:
				referencedElementID := currentEntry.changedElement.(*reference).ReferencedConceptID
				if referencedElementID != "" {
					uOfD.listenersMap.removeMappedValue(referencedElementID, currentID)
				}
			case *refinement:
				abstractID := currentEntry.changedElement.(*refinement).AbstractConceptID
				if abstractID != "" {
					uOfD.listenersMap.removeMappedValue(abstractID, currentID)
				}
				refinedID := currentEntry.changedElement.(*refinement).RefinedConceptID
				if refinedID != "" {
					uOfD.listenersMap.removeMappedValue(refinedID, currentID)
				}
			}
			// Update the uriUUIDMap
			uri := currentEntry.changedElement.GetURI(trans)
			if uri != "" {
				uOfD.uriUUIDMap.DeleteEntry(uri)
			}
			undoMgr.redoStack.Push(currentEntry)
			uOfD.removeElementForUndo(currentEntry.changedElement, trans)
			uOfD.setOwnedIDsMapValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
			uOfD.deleteUUIDElementMapEntry(currentEntry.changedElement.getConceptIDNoLock())
		} else if currentEntry.changeType == Deletion {
			// Update listeners. If this is a reference or refinement pointing to another element, add this element to the other element's listener's set
			switch currentEntry.priorState.(type) {
			case *reference:
				referencedElementID := currentEntry.priorState.(*reference).ReferencedConceptID
				if referencedElementID != "" {
					uOfD.listenersMap.addMappedValue(referencedElementID, currentEntry.priorState.GetConceptID(trans))
				}
			case *refinement:
				abstractID := currentEntry.priorState.(*refinement).AbstractConceptID
				if abstractID != "" {
					uOfD.listenersMap.addMappedValue(abstractID, currentEntry.priorState.GetConceptID(trans))
				}
				refinedID := currentEntry.priorState.(*refinement).RefinedConceptID
				if refinedID != "" {
					uOfD.listenersMap.addMappedValue(refinedID, currentEntry.priorState.GetConceptID(trans))
				}
			}
			// Update the uriUUIDMap
			uri := currentEntry.priorState.GetURI(trans)
			if uri != "" {
				uOfD.uriUUIDMap.SetEntry(uri, currentID)
			}
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, trans)
			undoMgr.redoStack.Push(currentEntry)
			uOfD.addElementForUndo(currentEntry.changedElement, trans)
			uOfD.setOwnedIDsMapValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
			uOfD.setUUIDElementMapEntry(currentEntry.changedElement.getConceptIDNoLock(), currentEntry.changedElement)
		} else if currentEntry.changeType == Change {
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			switch currentEntry.changedElement.(type) {
			case *reference:
				currentReferencedElementID := currentEntry.changedElement.(*reference).ReferencedConceptID
				priorReferencedElementID := currentEntry.priorState.(*reference).ReferencedConceptID
				if currentReferencedElementID != priorReferencedElementID {
					if currentReferencedElementID != "" {
						uOfD.listenersMap.removeMappedValue(currentReferencedElementID, currentID)
					}
					if priorReferencedElementID != "" {
						uOfD.listenersMap.addMappedValue(priorReferencedElementID, currentID)
					}
				}
			case *refinement:
				currentAbstractID := currentEntry.changedElement.(*refinement).AbstractConceptID
				priorAbstractID := currentEntry.priorState.(*refinement).AbstractConceptID
				if currentAbstractID != priorAbstractID {
					if currentAbstractID != "" {
						uOfD.listenersMap.removeMappedValue(currentAbstractID, currentID)
					}
					if priorAbstractID != "" {
						uOfD.listenersMap.addMappedValue(priorAbstractID, currentID)
					}
				}
				currentRefinedID := currentEntry.changedElement.(*refinement).RefinedConceptID
				priorRefinedID := currentEntry.priorState.(*refinement).RefinedConceptID
				if currentRefinedID != priorRefinedID {
					if currentRefinedID != "" {
						uOfD.listenersMap.removeMappedValue(currentRefinedID, currentID)
					}
					if priorRefinedID != "" {
						uOfD.listenersMap.addMappedValue(priorRefinedID, currentID)
					}
				}
			}
			// Update the uriUUIDMap
			currentURI := currentEntry.changedElement.GetURI(trans)
			priorURI := currentEntry.priorState.GetURI(trans)
			if currentURI != priorURI {
				if currentURI != "" {
					uOfD.uriUUIDMap.DeleteEntry(currentURI)
				}
				if priorURI != "" {
					uOfD.uriUUIDMap.SetEntry(priorURI, currentID)
				}
			}
			clone := clone(currentEntry.changedElement, trans)
			priorOwnedElements := uOfD.ownedIDsMap.GetMappedValues(currentID).Clone()
			priorListeners := uOfD.listenersMap.GetMappedValues(currentID).Clone()
			redoEntry := newUndoRedoStackEntry(Change, clone, priorOwnedElements, priorListeners, currentEntry.priorUofD, currentEntry.changedElement)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, trans)
			uOfD.setOwnedIDsMapValues(currentID, currentEntry.priorOwnedElements)
			uOfD.listenersMap.SetMappedValues(currentID, currentEntry.priorListeners)
			if currentEntry.priorUofD != uOfD.id {
				uOfD.deleteUUIDElementMapEntry(currentID)
			} else if currentEntry.priorUofD == uOfD.id {
				uOfD.setUUIDElementMapEntry(currentID, currentEntry.changedElement)
			}
			undoMgr.redoStack.Push(redoEntry)
		}
		firstEntry = false
	}
}
