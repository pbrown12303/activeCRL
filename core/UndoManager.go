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
func (undoMgr *undoManager) markChangedElement(changedElement *Concept, trans *Transaction) error {
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
func (undoMgr *undoManager) markNewElement(el *Concept, trans *Transaction) error {
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
func (undoMgr *undoManager) markRemovedElement(el *Concept, trans *Transaction) error {
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
	log.Printf("   Prior Listeners: %v", entry.priorListeners.ToSlice())
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
		}
		if currentEntry.changeType == Creation {
			// Update listeners. If this is a reference or refinement pointing to another element, add this element to the other element's listener's set
			currentOwnerID := currentEntry.changedElement.OwningConceptID
			if currentOwnerID != "" {
				uOfD.addMappedValueToOwnedIDsMap(currentOwnerID, currentID)
			}
			switch currentEntry.changedElement.GetConceptType() {
			case Reference:
				referencedElementID := currentEntry.changedElement.ReferencedConceptID
				if referencedElementID != "" {
					uOfD.addMappedValueToListenersMap(referencedElementID, currentID)
				}
			case Refinement:
				abstractID := currentEntry.changedElement.AbstractConceptID
				if abstractID != "" {
					uOfD.addMappedValueToListenersMap(abstractID, currentID)
				}
				refinedID := currentEntry.changedElement.RefinedConceptID
				if refinedID != "" {
					uOfD.addMappedValueToListenersMap(refinedID, currentID)
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
			uOfD.setMappedValuesForListenersMap(currentID, currentEntry.priorListeners)
		} else if currentEntry.changeType == Deletion {
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			priorOwnerID := currentEntry.priorState.OwningConceptID
			if priorOwnerID != "" {
				uOfD.removeMappedValueFromOwnedIDsMap(priorOwnerID, currentEntry.priorState.ConceptID)
			}
			switch currentEntry.priorState.GetConceptType() {
			case Reference:
				referencedElementID := currentEntry.priorState.ReferencedConceptID
				if referencedElementID != "" {
					uOfD.removeMappedValueFromListenersMap(referencedElementID, currentEntry.priorState.GetConceptID(trans))
				}
			case Refinement:
				abstractID := currentEntry.priorState.AbstractConceptID
				if abstractID != "" {
					uOfD.removeMappedValueFromListenersMap(abstractID, currentEntry.priorState.GetConceptID(trans))
				}
				refinedID := currentEntry.priorState.RefinedConceptID
				if refinedID != "" {
					uOfD.removeMappedValueFromListenersMap(refinedID, currentEntry.priorState.GetConceptID(trans))
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
			uOfD.setMappedValuesForListenersMap(currentID, currentEntry.priorListeners)
		} else if currentEntry.changeType == Change {
			// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
			currentOwnerID := currentEntry.changedElement.OwningConceptID
			if currentOwnerID != "" {
				uOfD.removeMappedValueFromOwnedIDsMap(currentOwnerID, currentID)
			}
			priorOwnerID := currentEntry.priorState.OwningConceptID
			if priorOwnerID != "" {
				uOfD.addMappedValueToOwnedIDsMap(priorOwnerID, currentEntry.priorState.ConceptID)
			}
			switch currentEntry.changedElement.GetConceptType() {
			case Reference:
				currentReferencedElementID := currentEntry.changedElement.ReferencedConceptID
				priorReferencedElementID := currentEntry.priorState.ReferencedConceptID
				if currentReferencedElementID != priorReferencedElementID {
					if currentReferencedElementID != "" {
						uOfD.removeMappedValueFromListenersMap(currentReferencedElementID, currentID)
					}
					if priorReferencedElementID != "" {
						uOfD.addMappedValueToListenersMap(priorReferencedElementID, currentID)
					}
				}
			case Refinement:
				currentAbstractID := currentEntry.changedElement.AbstractConceptID
				priorAbstractID := currentEntry.priorState.AbstractConceptID
				if currentAbstractID != priorAbstractID {
					if currentAbstractID != "" {
						uOfD.removeMappedValueFromListenersMap(currentAbstractID, currentID)
					}
					if priorAbstractID != "" {
						uOfD.addMappedValueToListenersMap(priorAbstractID, currentID)
					}
				}
				currentRefinedID := currentEntry.changedElement.RefinedConceptID
				priorRefinedID := currentEntry.priorState.RefinedConceptID
				if currentRefinedID != priorRefinedID {
					if currentRefinedID != "" {
						uOfD.removeMappedValueFromListenersMap(currentRefinedID, currentID)
					}
					if priorRefinedID != "" {
						uOfD.addMappedValueToListenersMap(priorRefinedID, currentID)
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
			uOfD.setMappedValuesForListenersMap(currentID, currentEntry.priorListeners)
			if currentEntry.priorUofD != uOfD.id {
				uOfD.deleteUUIDElementMapEntry(currentID)
			}
			undoMgr.undoStack.Push(undoEntry)
		}
	}
}

// restoreState is used as part of the undo process. It changes the currentState object
// to have the priorState.
func (undoMgr *undoManager) restoreState(priorState *Concept, currentState *Concept, trans *Transaction) {
	if undoMgr.debugUndo {
		log.Printf("Restoring State")
		log.Printf("   Current state:")
		Print(currentState, "      ", trans)
		log.Printf("   Prior state")
		Print(priorState, "      ", trans)
	}
	currentState.cloneAttributes(priorState, trans)
}

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
		} else {
			if currentEntry.changeType == Creation {
				// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
				currentOwnerID := currentEntry.changedElement.OwningConceptID
				if currentOwnerID != "" {
					uOfD.removeMappedValueFromOwnedIDsMap(currentOwnerID, currentID)
				}
				switch currentEntry.changedElement.GetConceptType() {
				case Reference:
					referencedElementID := currentEntry.changedElement.ReferencedConceptID
					if referencedElementID != "" {
						uOfD.removeMappedValueFromListenersMap(referencedElementID, currentID)
					}
				case Refinement:
					abstractID := currentEntry.changedElement.AbstractConceptID
					if abstractID != "" {
						uOfD.removeMappedValueFromListenersMap(abstractID, currentID)
					}
					refinedID := currentEntry.changedElement.RefinedConceptID
					if refinedID != "" {
						uOfD.removeMappedValueFromListenersMap(refinedID, currentID)
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
				uOfD.setMappedValuesForListenersMap(currentID, currentEntry.priorListeners)
			} else if currentEntry.changeType == Deletion {
				// Update listeners. If this is a reference or refinement pointing to another element, add this element to the other element's listener's set
				priorOwnerID := currentEntry.priorState.OwningConceptID
				if priorOwnerID != "" {
					uOfD.addMappedValueToOwnedIDsMap(priorOwnerID, currentEntry.priorState.ConceptID)
				}
				switch currentEntry.priorState.GetConceptType() {
				case Reference:
					referencedElementID := currentEntry.priorState.ReferencedConceptID
					if referencedElementID != "" {
						uOfD.addMappedValueToListenersMap(referencedElementID, currentEntry.priorState.GetConceptID(trans))
					}
				case Refinement:
					abstractID := currentEntry.priorState.AbstractConceptID
					if abstractID != "" {
						uOfD.addMappedValueToListenersMap(abstractID, currentEntry.priorState.GetConceptID(trans))
					}
					refinedID := currentEntry.priorState.RefinedConceptID
					if refinedID != "" {
						uOfD.addMappedValueToListenersMap(refinedID, currentEntry.priorState.GetConceptID(trans))
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
				uOfD.setMappedValuesForListenersMap(currentID, currentEntry.priorListeners)
			} else if currentEntry.changeType == Change {
				// Update listeners. If this is a reference or refinement pointing to another element, remove this element from the other element's listener's set
				currentOwnerID := currentEntry.changedElement.OwningConceptID
				if currentOwnerID != "" {
					uOfD.removeMappedValueFromOwnedIDsMap(currentOwnerID, currentID)
				}
				priorOwnerID := currentEntry.priorState.OwningConceptID
				if priorOwnerID != "" {
					uOfD.addMappedValueToOwnedIDsMap(priorOwnerID, currentEntry.priorState.ConceptID)
				}
				switch currentEntry.changedElement.GetConceptType() {
				case Reference:
					currentReferencedElementID := currentEntry.changedElement.ReferencedConceptID
					priorReferencedElementID := currentEntry.priorState.ReferencedConceptID
					if currentReferencedElementID != priorReferencedElementID {
						if currentReferencedElementID != "" {
							uOfD.removeMappedValueFromListenersMap(currentReferencedElementID, currentID)
						}
						if priorReferencedElementID != "" {
							uOfD.addMappedValueToListenersMap(priorReferencedElementID, currentID)
						}
					}
				case Refinement:
					currentAbstractID := currentEntry.changedElement.AbstractConceptID
					priorAbstractID := currentEntry.priorState.AbstractConceptID
					if currentAbstractID != priorAbstractID {
						if currentAbstractID != "" {
							uOfD.removeMappedValueFromListenersMap(currentAbstractID, currentID)
						}
						if priorAbstractID != "" {
							uOfD.addMappedValueToListenersMap(priorAbstractID, currentID)
						}
					}
					currentRefinedID := currentEntry.changedElement.RefinedConceptID
					priorRefinedID := currentEntry.priorState.RefinedConceptID
					if currentRefinedID != priorRefinedID {
						if currentRefinedID != "" {
							uOfD.removeMappedValueFromListenersMap(currentRefinedID, currentID)
						}
						if priorRefinedID != "" {
							uOfD.addMappedValueToListenersMap(priorRefinedID, currentID)
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
				uOfD.setMappedValuesForListenersMap(currentID, currentEntry.priorListeners)
				if currentEntry.priorUofD != uOfD.id {
					uOfD.deleteUUIDElementMapEntry(currentID)
				} else if currentEntry.priorUofD == uOfD.id {
					uOfD.setUUIDElementMapEntry(currentID, currentEntry.changedElement)
				}
				undoMgr.redoStack.Push(redoEntry)
			}
		}
		firstEntry = false
	}
}
