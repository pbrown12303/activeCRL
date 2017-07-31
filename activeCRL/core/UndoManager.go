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

func NewUndoManager() *undoManager {
	var undoMgr undoManager
	undoMgr.debugUndo = false
	undoMgr.recordingUndo = false
	return &undoMgr
}

// markChangedBaseElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markChangedBaseElement(changedElement BaseElement, hl *HeldLocks) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if hl == nil {
		return errors.New("UndoManager.markChangedBaseElement called with nil HeldLocks")
	}
	hl.LockBaseElement(changedElement)
	if undoMgr.debugUndo == true {
		debug.PrintStack()
	}
	clone := clone(changedElement)
	if undoMgr.recordingUndo {
		undoMgr.undoStack.Push(NewUndoRedoStackEntry(Change, clone, changedElement))
	}
	return nil
}

// markNewBaseElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markNewBaseElement(be BaseElement, hl *HeldLocks) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if hl == nil {
		return errors.New("UndoManager.markNewBaseElement called with nil HeldLocks")
	}
	hl.LockBaseElement(be)
	if undoMgr.debugUndo == true {
		debug.PrintStack()
	}
	clone := clone(be)
	if undoMgr.recordingUndo {
		undoMgr.undoStack.Push(NewUndoRedoStackEntry(Creation, clone, be))
	}
	return nil
}

// markRemoveBaseElement() If undo is enabled, updates the undo stack.
func (undoMgr *undoManager) markRemovedBaseElement(be BaseElement, hl *HeldLocks) error {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if hl == nil {
		return errors.New("UndoManager.markRemovedBaseElement called with nil HeldLocks")
	}
	hl.LockBaseElement(be)
	if undoMgr.debugUndo == true {
		debug.PrintStack()
	}
	clone := clone(be)
	if undoMgr.recordingUndo {
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

func (undoMgr *undoManager) redo(uOfD *UniverseOfDiscourse, hl *HeldLocks) {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	for len(undoMgr.redoStack) > 0 {
		currentEntry := undoMgr.redoStack.Pop()
		if currentEntry.changeType == Marker {
			undoMgr.undoStack.Push(currentEntry)
			return
		} else if currentEntry.changeType == Creation {
			undoMgr.undoStack.Push(currentEntry)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			// this was a new element
			uOfD.addBaseElementForUndo(currentEntry.changedElement, hl)
		} else if currentEntry.changeType == Deletion {
			undoMgr.undoStack.Push(currentEntry)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			// this was an deleted element
			uOfD.removeBaseElementForUndo(currentEntry.changedElement, hl)
		} else {
			clone := clone(currentEntry.changedElement)
			undoEntry := NewUndoRedoStackEntry(Change, clone, currentEntry.changedElement)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			undoMgr.undoStack.Push(undoEntry)
		}
	}
}

// restoreState is used as part of the undo process. It changes the currentState object
// to have the priorState.
func (undoMgr *undoManager) restoreState(priorState BaseElement, currentState BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	if undoMgr.debugUndo == true {
		log.Printf("Restoring State")
		log.Printf("   Current state:")
		Print(currentState, "      ", hl)
		log.Printf("   Prior state")
		Print(priorState, "      ", hl)
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

func (undoMgr *undoManager) setRecordingUndo(newSetting bool) {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	undoMgr.recordingUndo = newSetting
}

func (undoMgr *undoManager) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock Universe of Discourse %p\n", undoMgr)
	}
	undoMgr.Lock()
}

func (undoMgr *undoManager) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock Universe of Discourse %p\n", undoMgr)
	}
	undoMgr.Unlock()
}

func (undoMgr *undoManager) undo(uOfD *UniverseOfDiscourse, hl *HeldLocks) {
	undoMgr.TraceableLock()
	defer undoMgr.TraceableUnlock()
	firstEntry := true
	for len(undoMgr.undoStack) > 0 {
		currentEntry := undoMgr.undoStack.Pop()
		if currentEntry.changeType == Marker {
			if firstEntry {
				undoMgr.redoStack.Push(currentEntry)
			} else {
				// Put it back on the undo stack
				undoMgr.undoStack.Push(currentEntry)
				return
			}
		} else if currentEntry.changeType == Creation {
			undoMgr.redoStack.Push(currentEntry)
			uOfD.removeBaseElementForUndo(currentEntry.changedElement, hl)
		} else if currentEntry.changeType == Deletion {
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			undoMgr.redoStack.Push(currentEntry)
			uOfD.addBaseElementForUndo(currentEntry.changedElement, hl)
		} else if currentEntry.changeType == Change {
			clone := clone(currentEntry.changedElement)
			redoEntry := NewUndoRedoStackEntry(Change, clone, currentEntry.changedElement)
			undoMgr.restoreState(currentEntry.priorState, currentEntry.changedElement, hl)
			undoMgr.redoStack.Push(redoEntry)
		}
		firstEntry = false
	}
}
