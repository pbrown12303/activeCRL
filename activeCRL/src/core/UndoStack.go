package core

import ()

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
