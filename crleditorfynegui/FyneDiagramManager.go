package crleditorfynegui

import (
	"image/color"
	"reflect"

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

// FyneDiagramManager manages the relationship between the fyne DiagramWidgets and the
// underlying CRL model. It is a component of the  FyneGUI
type FyneDiagramManager struct {
	fyneGUI                *FyneGUI
	diagramArea            *fyne.Container
	diagramTabs            map[string]*container.TabItem
	toolbar                *fyne.Container
	toolButtons            map[string]*widget.Button
	tabArea                *container.DocTabs
	cursorTool             *widget.Button
	elementTool            *widget.Button
	LiteralTool            *widget.Button
	diagramObserver        *diagramObserver
	diagramElementObserver *diagramElementObserver
}

// NewFyneDiagramManager creates a diagram manager and associates it with the FyneGUI
func NewFyneDiagramManager(fyneGUI *FyneGUI) *FyneDiagramManager {
	var dm FyneDiagramManager
	dm.createToolbar()
	dm.diagramTabs = make(map[string]*container.TabItem)
	dm.tabArea = container.NewDocTabs()
	dm.tabArea.OnClosed = diagramClosed
	dm.diagramArea = container.NewBorder(nil, nil, dm.toolbar, nil, dm.tabArea)
	dm.diagramObserver = newDiagramObserver(&dm)
	dm.diagramElementObserver = newDiagramElementObserver(&dm)
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
		scrollingContainer := container.NewScroll(diagramwidget.NewDiagramWidget(diagramID))
		tabItem = container.NewTabItem(diagram.GetLabel(trans), scrollingContainer)
		dm.diagramTabs[diagramID] = tabItem
		dm.tabArea.Append(tabItem)
		diagram.Register(dm.diagramObserver)
		dm.populateDiagram(diagram, trans)
	}
	return nil
}

func (dm *FyneDiagramManager) getDiagramWidget(diagramID string) *diagramwidget.DiagramWidget {
	tabItem := dm.diagramTabs[diagramID]
	diagramWidget := tabItem.Content.(*container.Scroll).Content.(*diagramwidget.DiagramWidget)
	return diagramWidget
}

func (dm *FyneDiagramManager) GetDrawingArea() *fyne.Container {
	return dm.diagramArea
}

// populateDiagram adds all elements to the diagram
func (dm *FyneDiagramManager) populateDiagram(diagram core.Element, trans *core.Transaction) error {
	diagramWidget := dm.getDiagramWidget(diagram.GetConceptID(trans))
	nodes := diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, trans)
	for _, node := range nodes {
		// Get the node ID
		// Get the icon
		// Get the abstraction string
		// Build the node content
		// Display the text
		// Now create the node itself
		dm.addNodeToDiagram(node, trans, diagramWidget)
	}
	links := diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramLinkURI, trans)
	for _, link := range links {
		linkID := link.GetConceptID(trans)
		diagramLink := diagramWidget.Links[linkID]
		if diagramLink == nil {
			dm.addLinkToDiagram(link, trans, diagramWidget)
		}
	}
	diagramWidget.Refresh()
	return nil
}

func (dm *FyneDiagramManager) addElementToDiagram(element core.Element, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) {
	if element.IsRefinementOfURI(crldiagramdomain.CrlDiagramNodeURI, trans) {
		dm.addNodeToDiagram(element, trans, diagramWidget)
	} else if element.IsRefinementOfURI(crldiagramdomain.CrlDiagramLinkURI, trans) {
		dm.addLinkToDiagram(element, trans, diagramWidget)
	}
}

func (dm *FyneDiagramManager) addLinkToDiagram(link core.Element, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) {
	crlDiagramSource := crldiagramdomain.GetLinkSource(link, trans)
	if crlDiagramSource == nil {
		// Register for changes so that when sufficient information is present we can add it to the diagram
		link.Register(dm.diagramElementObserver)
		return
	}
	fyneSource := diagramWidget.GetDiagramElement(crlDiagramSource.GetConceptID(trans))
	fyneSourcePad := fyneSource.GetDefaultConnectionPad()
	crlDiagramTarget := crldiagramdomain.GetLinkTarget(link, trans)
	if crlDiagramTarget == nil {
		// Register for changes so that when sufficient information is present we can add it to the diagram
		link.Register(dm.diagramElementObserver)
		return
	}
	fyneTarget := diagramWidget.GetDiagramElement(crlDiagramTarget.GetConceptID(trans))
	fyneTargetPad := fyneTarget.GetDefaultConnectionPad()
	diagramLink := diagramwidget.NewDiagramLink(diagramWidget, fyneSourcePad, fyneTargetPad, link.GetConceptID(trans))
	diagramLink.AddMidpointAnchoredText("displayLabel", crldiagramdomain.GetDisplayLabel(link, trans))
	grey := color.RGBA{153, 153, 153, 255}
	if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramReferenceLinkURI, trans) {
		diagramLink.AddTargetDecoration(createReferenceArrowhead())
		diagramLink.AddSourceDecoration(createDiamond())
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramAbstractPointerURI, trans) {
		diagramLink.AddSourceDecoration(createRefinementTriangle())
		diagramLink.LinkColor = grey
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementPointerURI, trans) {
		diagramLink.AddTargetDecoration(createReferenceArrowhead())
		diagramLink.LinkColor = grey
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramOwnerPointerURI, trans) {
		diagramLink.AddTargetDecoration(createDiamond())
		diagramLink.LinkColor = grey
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinedPointerURI, trans) {
		diagramLink.AddSourceDecoration(createMirrorRefinementTriangle())
		diagramLink.LinkColor = grey
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinementLinkURI, trans) {
		diagramLink.AddMidpointDecoration(createRefinementTriangle())
	}
	link.Register(dm.diagramElementObserver)
}

func (dm *FyneDiagramManager) addNodeToDiagram(node core.Element, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) {
	nodeID := node.GetConceptID(trans)
	diagramNode := diagramWidget.Nodes[nodeID]
	if diagramNode == nil {

		modelElement := crldiagramdomain.GetReferencedModelElement(node, trans)
		nodeIcon := widget.NewIcon(getIconResource(modelElement, trans))

		abstractionString := crldiagramdomain.GetAbstractionDisplayLabel(node, trans)
		abstractionText := canvas.NewText(abstractionString, color.Black)
		abstractionText.TextSize = diagramWidget.DiagramTheme.Size(theme.SizeNameCaptionText)
		abstractionText.TextStyle = fyne.TextStyle{Bold: false, Italic: true, Monospace: false, Symbol: false, TabWidth: 4}

		hBox := container.NewHBox(nodeIcon, abstractionText)
		nodeLabel := crldiagramdomain.GetDisplayLabel(node, trans)
		entryWidget := widget.NewEntry()
		entryWidget.SetText(nodeLabel)
		entryWidget.Wrapping = fyne.TextWrapOff
		entryWidget.Refresh()
		nodeContainer := container.NewVBox(hBox, entryWidget)

		diagramNode = diagramwidget.NewDiagramNode(diagramWidget, nodeContainer, nodeID)
		x := crldiagramdomain.GetNodeX(node, trans)
		y := crldiagramdomain.GetNodeY(node, trans)
		fynePosition := fyne.NewPos(float32(x), float32(y))
		diagramNode.Move(fynePosition)
		diagramNode.Refresh()
		node.Register(dm.diagramElementObserver)
	}
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

// Update is the callback for changes to the core diagram
func (do *diagramObserver) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
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
	case core.OwnedConceptChanged:
		underlyingChange := notification.GetUnderlyingChange()
		beforeState := underlyingChange.GetBeforeConceptState()
		afterState := underlyingChange.GetAfterConceptState()
		diagramWidget := do.diagramManager.getDiagramWidget(notification.GetReportingElementID())
		switch underlyingChange.GetNatureOfChange() {
		case core.OwningConceptChanged:
			if beforeState.OwningConceptID == "" && afterState.OwningConceptID != "" {
				// the element has been added
				uOfD := trans.GetUniverseOfDiscourse()
				element := uOfD.GetElement(afterState.ConceptID)
				do.diagramManager.addElementToDiagram(element, trans, diagramWidget)
			} else {
				// the element has been removed
			}

		}
	}
	return nil
}

type diagramElementObserver struct {
	diagramManager *FyneDiagramManager
}

func newDiagramElementObserver(dm *FyneDiagramManager) *diagramElementObserver {
	deo := diagramElementObserver{diagramManager: dm}
	return &deo
}

// Update is the callback for changes to the core diagram element
func (deo *diagramElementObserver) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
	diagramID := notification.GetReportingElementState().OwningConceptID
	tabItem := deo.diagramManager.diagramTabs[diagramID]
	if tabItem == nil {
		// the diagram is not being displayed
		return nil
	}
	diagramWidget := deo.diagramManager.getDiagramWidget(notification.GetReportingElementState().OwningConceptID)
	elementID := notification.GetReportingElementID()
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(elementID)
	fyneDiagramElement := diagramWidget.GetDiagramElement(elementID)
	if reflect.ValueOf(fyneDiagramElement).IsNil() {
		deo.diagramManager.addElementToDiagram(crlDiagramElement, trans, diagramWidget)
	}
	return nil
}

var referenceArrowHeadPoints []fyne.Position = []fyne.Position{
	{X: 0, Y: 0},
	{X: 8, Y: 5},
	{X: 8, Y: -5},
}

func createReferenceArrowhead() *diagramwidget.Polygon {
	polygon := diagramwidget.NewPolygon(referenceArrowHeadPoints)
	polygon.SetSolid(true)
	polygon.SetClosed(true)
	return polygon
}

var diamondPoints []fyne.Position = []fyne.Position{
	{X: 0, Y: 0},
	{X: 8, Y: 4},
	{X: 16, Y: 0},
	{X: 8, Y: -4},
}

func createDiamond() *diagramwidget.Polygon {
	polygon := diagramwidget.NewPolygon(diamondPoints)
	polygon.SetSolid(true)
	polygon.SetClosed(true)
	return polygon
}

var refinementTrianglePoints []fyne.Position = []fyne.Position{
	{X: 0, Y: 8},
	{X: 16, Y: 0},
	{X: 0, Y: -8},
}

func createRefinementTriangle() *diagramwidget.Polygon {
	polygon := diagramwidget.NewPolygon(refinementTrianglePoints)
	polygon.SetSolid(false)
	polygon.SetClosed(true)
	return polygon
}

var mirrorRefinementTrianglePoints []fyne.Position = []fyne.Position{
	{X: 0, Y: 0},
	{X: 16, Y: 8},
	{X: 16, Y: -8},
}

func createMirrorRefinementTriangle() *diagramwidget.Polygon {
	polygon := diagramwidget.NewPolygon(mirrorRefinementTrianglePoints)
	polygon.SetSolid(false)
	polygon.SetClosed(true)
	return polygon
}
