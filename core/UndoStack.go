// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

type undoStack []*undoRedoStackEntry

func (s undoStack) Empty() bool {
	return len(s) == 0
}

func (s undoStack) Peek() *undoRedoStackEntry {
	return s[len(s)-1]
}

func (s *undoStack) Push(entry *undoRedoStackEntry) {
	(*s) = append((*s), entry)
}

func (s *undoStack) Pop() *undoRedoStackEntry {
	entry := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return entry
}
