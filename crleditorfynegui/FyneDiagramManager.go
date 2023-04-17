package crleditorfynegui

import (
	"image/color"

	"fyne.io/x/fyne/widget/diagramwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/images"
)

type FyneDiagramManager struct {
	fyneGUI         *FyneGUI
	diagramArea     *fyne.Container
	diagramTabs     map[string]*container.TabItem
	toolbar         *fyne.Container
	toolButtons     map[string]*widget.Button
	tabArea         *container.DocTabs
	cursorTool      *widget.Button
	elementTool     *widget.Button
	LiteralTool     *widget.Button
	diagramObserver *diagramObserver
}

func NewFyneDiagramManager(fyneGUI *FyneGUI) *FyneDiagramManager {
	var dm FyneDiagramManager
	dm.createToolbar()
	dm.diagramTabs = make(map[string]*container.TabItem)
	dm.tabArea = container.NewDocTabs()
	dm.tabArea.OnClosed = diagramClosed
	dm.diagramArea = container.NewBorder(nil, nil, dm.toolbar, nil, dm.tabArea)
	dm.diagramObserver = newDiagramObserver(&dm)
	return &dm
}

func (dm *FyneDiagramManager) closeDiagram(diagramID string) {
	tabItem := dm.diagramTabs[diagramID]
	if tabItem != nil {
		dm.tabArea.Remove(tabItem)
		delete(dm.diagramTabs, diagramID)
		diagram := dm.fyneGUI.editor.GetUofD().GetElement(diagramID)
		diagram.Deregister(dm.diagramObserver)
	}
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
		tabItem = container.NewTabItem(diagram.GetLabel(trans), diagramwidget.NewDiagramWidget(diagramID))
		dm.diagramTabs[diagramID] = tabItem
		dm.tabArea.Append(tabItem)
		diagram.Register(dm.diagramObserver)
		dm.refreshDiagram(diagram, trans)
	}
	return nil
}

func (dm *FyneDiagramManager) GetDrawingArea() *fyne.Container {
	return dm.diagramArea
}

// refreshDiagram resends all diagram elements to the browser
func (dm *FyneDiagramManager) refreshDiagram(diagram core.Element, trans *core.Transaction) error {
	tabItem := dm.diagramTabs[diagram.GetConceptID(trans)]
	diagramWidget := tabItem.Content.(*diagramwidget.DiagramWidget)
	nodes := diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, trans)
	for _, node := range nodes {
		// Get the node ID
		nodeID := node.GetConceptID(trans)
		// Get the icon
		modelElement := crldiagramdomain.GetReferencedModelElement(node, trans)
		nodeIcon := widget.NewIcon(getIconResource(modelElement, trans))
		// Get the abstraction string
		abstractionString := crldiagramdomain.GetAbstractionDisplayLabel(node, trans)
		abstractionText := canvas.NewText(abstractionString, color.Black)
		abstractionText.TextSize = diagramWidget.DiagramTheme.Size(theme.SizeNameCaptionText)
		abstractionText.TextStyle = fyne.TextStyle{Bold: false, Italic: true, Monospace: false, Symbol: false, TabWidth: 4}
		// Build the node content
		hBox := container.NewHBox(nodeIcon, abstractionText)
		nodeLabel := crldiagramdomain.GetDisplayLabel(node, trans)
		entryWidget := widget.NewEntry()
		entryWidget.SetText(nodeLabel)
		entryWidget.Wrapping = fyne.TextWrapOff
		entryWidget.Refresh() // Display the text
		nodeContainer := container.NewVBox(hBox, entryWidget)
		// Now create the node itself
		diagramNode := diagramwidget.NewDiagramNode(diagramWidget, nodeContainer, nodeID)
		x := crldiagramdomain.GetNodeX(node, trans)
		y := crldiagramdomain.GetNodeY(node, trans)
		fynePosition := fyne.NewPos(float32(x), float32(y))
		diagramNode.Move(fynePosition)
		diagramNode.Refresh()
	}
	// links := diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramLinkURI, hl)
	// for _, link := range links {
	// additionalParameters := getLinkAdditionalParameters(link, hl)
	// conceptState, err2 := core.NewConceptState(link)
	// if err2 != nil {
	// 	return errors.Wrap(err2, "diagramManager.refreshDiagram failed")
	// }
	// link.Register(dmPtr.elementManager)
	// notificationResponse, err := BrowserGUISingleton.SendNotification("AddDiagramLink", link.GetConceptID(hl), conceptState, additionalParameters)
	// if err != nil {
	// 	return errors.Wrap(err, "diagramManager.refreshDiagram failed")
	// }
	// if notificationResponse.Result != 0 {
	// 	return errors.New(notificationResponse.ErrorMessage)
	// }
	// }
	diagramWidget.Refresh()
	return nil
}

func diagramClosed(tabItem *container.TabItem) {
	for k, v := range FyneGUISingleton.diagramManager.diagramTabs {
		if v == tabItem {
			trans := GetTransaction()
			defer trans.ReleaseLocks()
			delete(FyneGUISingleton.diagramManager.diagramTabs, k)
			crleditor.CrlEditorSingleton.CloseDiagramView(k, trans)
			return
		}
	}
}

// diagramObserver monitors the core diagram for changes relevant to the displayed diagram
type diagramObserver struct {
	diagramManager *FyneDiagramManager
}

func newDiagramObserver(dm *FyneDiagramManager) *diagramObserver {
	do := diagramObserver{diagramManager: dm}
	return &do
}

// updateDiagram is the callback for changes to the core diagram
func (do *diagramObserver) Update(notification *core.ChangeNotification, heldLocks *core.Transaction) error {
	tabItem := do.diagramManager.diagramTabs[notification.GetReportingElementID()]
	if tabItem == nil {
		// the diagram is not being displayed
		return nil
	}
	switch notification.GetNatureOfChange() {
	case core.ConceptChanged:
		beforeStateLabel := notification.GetBeforeConceptState().Label
		afterStateLabel := notification.GetAfterConceptState().Label
		if beforeStateLabel != afterStateLabel {
			tabItem.Text = afterStateLabel
			do.diagramManager.tabArea.Refresh()
		}
	}
	return nil
}
