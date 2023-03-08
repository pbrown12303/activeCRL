package crleditorfynegui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/images"
)

type FyneDiagramManager struct {
	diagramArea *fyne.Container
	diagramTabs map[string]*container.TabItem
	toolbar     *fyne.Container
	toolButtons map[string]*widget.Button
	tabArea     *container.DocTabs
	cursorTool  *widget.Button
	elementTool *widget.Button
	LiteralTool *widget.Button
}

func NewFyneDiagramManager() *FyneDiagramManager {
	var dm FyneDiagramManager
	dm.createToolbar()
	dm.diagramTabs = make(map[string]*container.TabItem)
	dm.tabArea = container.NewDocTabs()
	dm.diagramArea = container.NewBorder(nil, nil, dm.toolbar, nil, dm.tabArea)

	return &dm
}

func (dm *FyneDiagramManager) createToolbar() {
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

func (dm *FyneDiagramManager) displayDiagram(diagram core.Element, trans *core.Transaction) error {
	diagramID := diagram.GetConceptID(trans)
	tabItem := dm.diagramTabs[diagramID]
	if tabItem == nil {
		tabItem = container.NewTabItem(diagram.GetLabel(trans), container.NewWithoutLayout())
		dm.diagramTabs[diagramID] = tabItem
		dm.tabArea.Append(tabItem)
	}
	return nil
}

func (dm *FyneDiagramManager) GetDrawingArea() *fyne.Container {
	return dm.diagramArea
}
