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
}

// NewUndoManager creates and initializes the manager for the undo/redo functionality
func newUndoManager() *undoManager {
	var undoMgr undoManager
	undoMgr.debugUndo = false
	undoMgr.recordingUndo = false
	return &undoMgr
}

// markChangedElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markChangedElement(changedElement Element, hl *HeldLocks) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	hl.ReadLockElement(changedElement)
	if undoMgr.debugUndo == true {
		debug.PrintStack()
	}
	clone := clone(changedElement, hl)
	if undoMgr.recordingUndo {
		undoMgr.undoStack.Push(NewUndoRedoStackEntry(Change, clone, changedElement))
	}
	return nil
}

// markNewElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markNewElement(be Element, hl *HeldLocks) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if hl == nil {
		return errors.New("UndoManager.markNewElement called with nil HeldLocks")
	}
	hl.ReadLockElement(be)
	if undoMgr.debugUndo == true {
		debug.PrintStack()
	}
	if undoMgr.recordingUndo {
		clone := clone(be, hl)
		stackEntry := NewUndoRedoStackEntry(Creation, clone, be)
		if undoMgr.debugUndo == true {
			PrintStackEntry(stackEntry, hl)
		}
		undoMgr.undoStack.Push(stackEntry)
	}
	return nil
}

// markRemoveElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markRemovedElement(be Element, hl *HeldLocks) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if hl == nil {
		return errors.New("UndoManager.markRemovedElement called with nil HeldLocks")
	}
	hl.ReadLockElement(be)
	if undoMgr.debugUndo == true {
		debug.PrintStack()
	}
	if undoMgr.recordingUndo {
		clone := clone(be, hl)
		undoMgr.undoStack.Push(NewUndoRedoStackEntry(Deletion, clone, be))
	}
	return nil
}

// MarkUndoPoint() If undo is enabled, puts a marker on the undo stack.
func (undoMgr *undoManager) MarkUndoPoint() {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if undoMgr.recordingUndo {
		undoMgr.undoStack.Push(NewUndoRedoStackEntry(Marker, nil, nil))
	}
}

func PrintUndoStack(s undoStack, stackName string, uOfD UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocks()
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

func PrintStackEntry(entry *undoRedoStackEntry, hl *HeldLocks) {
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

func (undoMgr *undoManager) redo(uOfD UniverseOfDiscourse, hl *HeldLocks) {
	// TODO:
	// undoMgr.TraceableLock()
	// defer undoMgr.TraceableUnlock()
	// if hl == nil {
	// 	hl = NewHeldLocks(nil)
	// 	defer hl.ReleaseLocks()
	// }
	// for len(undoMgr.redoStack) > 0 {
	// 	currentEntry := undoMgr.redoStack.Pop()
	// 	if currentEntry.changeType == Marker {
	// 		undoMgr.undoStack.Push(currentEntry)
	// 		return
	// 	} else if currentEntry.changeType == Creation {
	// 		undoMgr.undoStack.Push(currentEntry)
	// 		undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
	// 		// this was a new element
	// 		uOfD.(*universeOfDiscourse).addElementForUndo(currentEntry.changedElement, hl)
	// 	} else if currentEntry.changeType == Deletion {
	// 		undoMgr.undoStack.Push(currentEntry)
	// 		undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
	// 		// this was an deleted element
	// 		uOfD.(*universeOfDiscourse).removeElementForUndo(currentEntry.changedElement, hl)
	// 	} else {
	// 		clone := clone(currentEntry.changedElement)
	// 		undoEntry := NewUndoRedoStackEntry(Change, clone, currentEntry.changedElement)
	// 		undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
	// 		undoMgr.undoStack.Push(undoEntry)
	// 	}
	// }
}

// restoreState is used as part of the undo process. It changes the currentState object
// to have the priorState.
func (undoMgr *undoManager) restoreState(priorState Element, currentState Element, hl *HeldLocks) {
	if undoMgr.debugUndo == true {
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

func (undoMgr *undoManager) setDebugUndo(newSetting bool) {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	undoMgr.debugUndo = newSetting
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

func (undoMgr *undoManager) undo(uOfD UniverseOfDiscourse, hl *HeldLocks) {
	// TODO:
	// undoMgr.TraceableLock()
	// defer undoMgr.TraceableUnlock()
	// firstEntry := true
	// for len(undoMgr.undoStack) > 0 {
	// 	currentEntry := undoMgr.undoStack.Pop()
	// 	if currentEntry.changeType == Marker {
	// 		if firstEntry {
	// 			undoMgr.redoStack.Push(currentEntry)
	// 		} else {
	// 			// Put it back on the undo stack
	// 			undoMgr.undoStack.Push(currentEntry)
	// 			return
	// 		}
	// 	} else if currentEntry.changeType == Creation {
	// 		undoMgr.redoStack.Push(currentEntry)
	// 		uOfD.(*universeOfDiscourse).removeElementForUndo(currentEntry.changedElement, hl)
	// 	} else if currentEntry.changeType == Deletion {
	// 		undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
	// 		undoMgr.redoStack.Push(currentEntry)
	// 		uOfD.(*universeOfDiscourse).addElementForUndo(currentEntry.changedElement, hl)
	// 	} else if currentEntry.changeType == Change {
	// 		clone := clone(currentEntry.changedElement)
	// 		redoEntry := NewUndoRedoStackEntry(Change, clone, currentEntry.changedElement)
	// 		undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
	// 		undoMgr.redoStack.Push(redoEntry)
	// 	}
	// 	firstEntry = false
	// }
}
