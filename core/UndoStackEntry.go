// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

type UndoChangeType int

const (
	Marker UndoChangeType = iota
	Creation
	Deletion
	Change
)

type undoRedoStackEntry struct {
	changeType     UndoChangeType
	priorState     Element
	changedElement Element
}

func NewUndoRedoStackEntry(changeType UndoChangeType, priorState Element, changedElement Element) *undoRedoStackEntry {
	var entry undoRedoStackEntry
	entry.changeType = changeType
	entry.priorState = priorState
	entry.changedElement = changedElement
	return &entry
}
