package fynegui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pbrown12303/activeCRL/images"
)

type FyneDrawingManager struct {
	drawingArea *fyne.Container
	toolbar     *fyne.Container
	toolButtons map[string]*widget.Button
	tabArea     *container.DocTabs
	cursorTool  *widget.Button
	elementTool *widget.Button
	LiteralTool *widget.Button
}

func NewFyneDrawingManager() *FyneDrawingManager {
	var dm FyneDrawingManager
	dm.createToolbar()
	dm.tabArea = container.NewDocTabs()
	dm.drawingArea = container.NewBorder(nil, nil, dm.toolbar, nil, dm.tabArea)

	return &dm
}

func (dm *FyneDrawingManager) GetDrawingArea() *fyne.Container {
	return dm.drawingArea
}

func (dm *FyneDrawingManager) createToolbar() {
	dm.toolbar = container.NewVBox()
	dm.toolButtons = make(map[string]*widget.Button)
	// Cursor
	button := widget.NewButtonWithIcon("", images.ResourceCursorIconPng, nil)
	dm.toolButtons["Cursor"] = button
	dm.toolbar.Add(button)
	// Element
	button = widget.NewButtonWithIcon("", images.ResourceElementIconPng, nil)
	dm.toolButtons["Element"] = button
	dm.toolbar.Add(button)
	// Literal
	button = widget.NewButtonWithIcon("", images.ResourceLiteralIconPng, nil)
	dm.toolButtons["Literal"] = button
	dm.toolbar.Add(button)
	// Reference
	button = widget.NewButtonWithIcon("", images.ResourceReferenceIconPng, nil)
	dm.toolButtons["Reference"] = button
	dm.toolbar.Add(button)
	// ReferenceLink
	button = widget.NewButtonWithIcon("", images.ResourceReferenceLinkIconPng, nil)
	dm.toolButtons["ReferenceLink"] = button
	dm.toolbar.Add(button)
	// Refinement
	button = widget.NewButtonWithIcon("", images.ResourceRefinementIconPng, nil)
	dm.toolButtons["Refinement"] = button
	dm.toolbar.Add(button)
	// RefinementLink
	button = widget.NewButtonWithIcon("", images.ResourceRefinementLinkIconPng, nil)
	dm.toolButtons["RefinementLink"] = button
	dm.toolbar.Add(button)
	// OwnerPointer
	button = widget.NewButtonWithIcon("", images.ResourceOwnerPointerIconPng, nil)
	dm.toolButtons["OwnerPointer"] = button
	dm.toolbar.Add(button)
	// ElementPointer
	button = widget.NewButtonWithIcon("", images.ResourceElementPointerIconPng, nil)
	dm.toolButtons["ElementPointer"] = button
	dm.toolbar.Add(button)
	// AbstractPointer
	button = widget.NewButtonWithIcon("", images.ResourceAbstractPointerIconPng, nil)
	dm.toolButtons["AbstractPointer"] = button
	dm.toolbar.Add(button)
	// RefinedPointer
	button = widget.NewButtonWithIcon("", images.ResourceRefinedPointerIconPng, nil)
	dm.toolButtons["RefinedPointer"] = button
	dm.toolbar.Add(button)
}
