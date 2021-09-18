package fynegui

import (
	"fyne.io/fyne/widget"
)

// FyneTreeManager is the manager of the fyne tree in the CrlFyneEditor
type FyneTreeManager struct {
	tree *widget.Tree
}

// NewFyneTreeManager returns an initialized FyneTreeManager
func NewFyneTreeManager() *FyneTreeManager {
	var treeManager FyneTreeManager
	// TODO: update this for the latest version of fyne
	//	treeManager.tree = widget.NewTree()
	return &treeManager
}
