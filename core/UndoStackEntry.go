// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import mapset "github.com/deckarep/golang-set"

// UndoChangeType identifies the type of undo change
type UndoChangeType int

const (
	// Marker marks the point on the stack at which an undo or redo operation will stop
	Marker UndoChangeType = iota
	// Creation marks the creation of a new Element
	Creation
	// Deletion marks the deletion of an Element
	Deletion
	// Change marks a change to an element
	Change
)

type undoRedoStackEntry struct {
	changeType         UndoChangeType
	priorState         Element
	priorOwnedElements mapset.Set
	priorListeners     mapset.Set
	priorUofD          string
	changedElement     Element
}

func newUndoRedoStackEntry(changeType UndoChangeType, priorState Element, priorOwnedElements mapset.Set, priorListeners mapset.Set, priorUofD string, changedElement Element) *undoRedoStackEntry {
	var entry undoRedoStackEntry
	entry.changeType = changeType
	entry.priorState = priorState
	entry.priorOwnedElements = priorOwnedElements
	entry.priorListeners = priorListeners
	entry.priorUofD = priorUofD
	entry.changedElement = changedElement
	return &entry
}
