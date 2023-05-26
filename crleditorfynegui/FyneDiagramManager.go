package crleditorfynegui

import (
	"errors"
	"image/color"
	"log"
	"reflect"
	"strconv"

	"fyne.io/x/fyne/widget/diagramwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/images"
)

const (
	displayLabel = "DisplayLabel"
)

type diagramTab struct {
	diagramID string
	tab       *container.TabItem
	diagram   *diagramwidget.DiagramWidget
}

// FyneDiagramManager manages the relationship between the fyne DiagramWidgets and the
// underlying CRL model. It is a component of the  FyneGUI
type FyneDiagramManager struct {
	fyneGUI                *CrlEditorFyneGUI
	diagramArea            *fyne.Container
	diagramTabs            map[string]*diagramTab
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
func NewFyneDiagramManager(fyneGUI *CrlEditorFyneGUI) *FyneDiagramManager {
	var dm FyneDiagramManager
	dm.createToolbar()
	dm.diagramTabs = make(map[string]*diagramTab)
	dm.tabArea = container.NewDocTabs()
	dm.tabArea.OnClosed = diagramClosed
	dm.diagramArea = container.NewBorder(nil, nil, dm.toolbar, nil, dm.tabArea)
	dm.diagramObserver = newDiagramObserver(&dm)
	dm.diagramElementObserver = newDiagramElementObserver(&dm)
	dm.fyneGUI = fyneGUI
	return &dm
}

func (dm *FyneDiagramManager) addElementToDiagram(element core.Element, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramElement {
	if element.IsRefinementOfURI(crldiagramdomain.CrlDiagramNodeURI, trans) {
		return dm.addNodeToDiagram(element, trans, diagramWidget)
	} else if element.IsRefinementOfURI(crldiagramdomain.CrlDiagramLinkURI, trans) {
		return dm.addLinkToDiagram(element, trans, diagramWidget)
	}
	return nil
}

func (dm *FyneDiagramManager) addLinkToDiagram(link core.Element, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) *FyneCrlDiagramLink {
	crlDiagramSource := crldiagramdomain.GetLinkSource(link, trans)
	if crlDiagramSource == nil {
		// Register for changes so that when sufficient information is present we can add it to the diagram
		link.Register(dm.diagramElementObserver)
		return nil
	}
	fyneSource := diagramWidget.GetDiagramElement(crlDiagramSource.GetConceptID(trans))
	fyneSourcePad := fyneSource.GetDefaultConnectionPad()
	crlDiagramTarget := crldiagramdomain.GetLinkTarget(link, trans)
	if crlDiagramTarget == nil {
		// Register for changes so that when sufficient information is present we can add it to the diagram
		link.Register(dm.diagramElementObserver)
		return nil
	}
	fyneTarget := diagramWidget.GetDiagramElement(crlDiagramTarget.GetConceptID(trans))
	fyneTargetPad := fyneTarget.GetDefaultConnectionPad()
	diagramLink := NewFyneCrlDiagramLink(diagramWidget, fyneSourcePad, fyneTargetPad, link, trans)
	link.Register(dm.diagramElementObserver)
	return diagramLink
}

func (dm *FyneDiagramManager) addNodeToDiagram(node core.Element, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramNode {
	nodeID := node.GetConceptID(trans)
	diagramNode := diagramWidget.Nodes[nodeID]
	if diagramNode == nil {

		diagramNode = NewFyneCrlDiagramNode(node, trans, diagramWidget)
		node.Register(dm.diagramElementObserver)
	}
	return diagramNode
}

func (dm *FyneDiagramManager) closeDiagram(diagramID string) {
	tabItem := dm.diagramTabs[diagramID]
	if tabItem != nil {
		dm.tabArea.Remove(tabItem.tab)
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

func (dm *FyneDiagramManager) diagramElementSelectionChanged(diagramElementID string) {
	editor := dm.fyneGUI.editor
	trans, new := editor.GetTransaction()
	if new {
		defer trans.ReleaseLocks()
	}
	dm.fyneGUI.editor.SelectElementUsingIDString(diagramElementID, trans)
}

func (dm *FyneDiagramManager) displayDiagram(diagram core.Element, trans *core.Transaction) error {
	diagramID := diagram.GetConceptID(trans)
	tabItem := dm.diagramTabs[diagramID]
	if tabItem == nil {
		diagramWidget := diagramwidget.NewDiagramWidget(diagramID)
		scrollingContainer := container.NewScroll(diagramWidget)
		newTabItem := &diagramTab{
			diagramID: diagramID,
			tab:       container.NewTabItem(diagram.GetLabel(trans), scrollingContainer),
			diagram:   diagramWidget,
		}
		dm.diagramTabs[diagramID] = newTabItem
		dm.tabArea.Append(newTabItem.tab)
		diagram.Register(dm.diagramObserver)
		dm.populateDiagram(diagram, trans)
		diagramWidget.LinkConnectionChangedCallback = func(link diagramwidget.DiagramLink, end string, oldPad diagramwidget.ConnectionPad, newPad diagramwidget.ConnectionPad) {
			dm.linkConnectionChanged(link, end, oldPad, newPad)
		}
		diagramWidget.PrimaryDiagramElementSelectionChangedCallback = func(id string) {
			dm.diagramElementSelectionChanged(id)
		}
	}
	return nil
}

func (dm *FyneDiagramManager) ElementSelected(id string, trans *core.Transaction) {
	for _, tabItem := range dm.diagramTabs {
		dm.selectElementInDiagram(id, tabItem.diagram, trans)
	}
}

func (dm *FyneDiagramManager) getDiagramWidget(diagramID string) *diagramwidget.DiagramWidget {
	tabItem := dm.diagramTabs[diagramID]
	diagramWidget := tabItem.diagram
	return diagramWidget
}

func (dm *FyneDiagramManager) GetDrawingArea() *fyne.Container {
	return dm.diagramArea
}

// linkConnectionChanged is the callback for changes in link connections
func (dm *FyneDiagramManager) linkConnectionChanged(link diagramwidget.DiagramLink, end string, oldPad diagramwidget.ConnectionPad, newPad diagramwidget.ConnectionPad) error {
	trans, new := dm.fyneGUI.editor.GetTransaction()
	if new {
		defer trans.ReleaseLocks()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	crlLink := uOfD.GetElement(link.GetDiagramElementID())
	if crlLink == nil {
		return errors.New("In FyneDiagramManager.linkConnectionChanged CrlLink not found")
	}
	crlNewPadOwner := uOfD.GetElement(newPad.GetPadOwner().GetDiagramElementID())
	if crlNewPadOwner == nil {
		return errors.New("In FyneDiagramManager.linkConnectionChanged CrlLink not found")
	}
	switch end {
	case "source":
		crldiagramdomain.SetLinkSource(crlLink, crlNewPadOwner, trans)
	case "target":
		crldiagramdomain.SetLinkTarget(crlLink, crlNewPadOwner, trans)
	}
	return nil
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

func (dm *FyneDiagramManager) selectElementInDiagram(elementID string, diagram *diagramwidget.DiagramWidget, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	foundDiagramElementID := ""
	for key, _ := range diagram.GetDiagramElements() {
		crlDiagramElement := uOfD.GetElement(key)
		if crlDiagramElement != nil {
			crlModelElement := crldiagramdomain.GetReferencedModelElement(crlDiagramElement, trans)
			if crlModelElement != nil {
				if crlModelElement.GetConceptID(trans) == elementID {
					foundDiagramElementID = key
					break
				}
			}
		}
	}
	if foundDiagramElementID != "" {
		diagram.SelectDiagramElementNoCallback(foundDiagramElementID)
	}
	return nil
}

func diagramClosed(tabItem *container.TabItem) {
	for k, v := range FyneGUISingleton.diagramManager.diagramTabs {
		if v.tab == tabItem {
			trans, isNew := FyneGUISingleton.editor.GetTransaction()
			if isNew {
				defer FyneGUISingleton.editor.EndTransaction()
			}
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
			tabItem.tab.Text = afterStateLabel
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
				diagramWidget.RemoveElement(afterState.ConceptID)
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
	if notification.GetNatureOfChange() == core.ConceptRemoved {
		return nil
	}
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
		fyneDiagramElement = deo.diagramManager.addElementToDiagram(crlDiagramElement, trans, diagramWidget)
	}
	switch typedElement := fyneDiagramElement.(type) {
	case *FyneCrlDiagramNode:
		switch notification.GetNatureOfChange() {
		case core.OwnedConceptChanged:
			ownedConceptChangedNotification := notification.GetUnderlyingChange()
			switch ownedConceptChangedNotification.GetNatureOfChange() {
			case core.ConceptChanged:
				changedConcept := trans.GetUniverseOfDiscourse().GetElement(ownedConceptChangedNotification.GetChangedConceptID())
				if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementDisplayLabelURI, trans) {
					beforeState := ownedConceptChangedNotification.GetBeforeConceptState()
					afterState := ownedConceptChangedNotification.GetAfterConceptState()
					if afterState.LiteralValue != beforeState.LiteralValue {
						typedElement.labelBinding.Set(afterState.LiteralValue)
						typedElement.entryWidget.Refresh()
						fyneDiagramElement.Refresh()
					}
					return nil
				}
				if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramNodeXURI, trans) {
					x := float32(crldiagramdomain.GetNodeX(crlDiagramElement, trans))
					fynePosition := fyneDiagramElement.Position()
					if x != fynePosition.X {
						fyneDiagramElement.Move(fyne.NewPos(x, fynePosition.Y))

					}
					return nil
				}
				if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramNodeYURI, trans) {
					y := float32(crldiagramdomain.GetNodeY(crlDiagramElement, trans))
					fynePosition := fyneDiagramElement.Position()
					if y != fynePosition.Y {
						fyneDiagramElement.Move(fyne.NewPos(fynePosition.X, y))

					}
					return nil
				}
				if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementLineColorURI, trans) {
					lineColor := crldiagramdomain.GetLineColor(crlDiagramElement, trans)
					log.Printf("Line Color: %s", lineColor)
					goColor := getGoColor(lineColor)
					fyneDiagramElement.SetForegroundColor(goColor)

				}
			}
		}
	case *FyneCrlDiagramLink:
		switch notification.GetNatureOfChange() {
		case core.OwnedConceptChanged:
			ownedConceptChangedNotification := notification.GetUnderlyingChange()
			switch ownedConceptChangedNotification.GetNatureOfChange() {
			case core.ConceptChanged:
				changedConcept := trans.GetUniverseOfDiscourse().GetElement(ownedConceptChangedNotification.GetChangedConceptID())
				if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementDisplayLabelURI, trans) {
					beforeState := ownedConceptChangedNotification.GetBeforeConceptState()
					afterState := ownedConceptChangedNotification.GetAfterConceptState()
					if afterState.LiteralValue != beforeState.LiteralValue {
						typedElement.SetLabel(afterState.LiteralValue)
						fyneDiagramElement.Refresh()
					}
					return nil
				}
				if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementLineColorURI, trans) {
					lineColor := crldiagramdomain.GetLineColor(crlDiagramElement, trans)
					log.Printf("Line Color: %s", lineColor)
					goColor := getGoColor(lineColor)
					fyneDiagramElement.SetForegroundColor(goColor)
					return nil
				}
			}
		}

	}
	return nil
}

func getGoColor(lineColor string) color.RGBA {
	redString := lineColor[1:3]
	red, _ := strconv.ParseUint(redString, 16, 8)
	greenString := lineColor[3:5]
	green, _ := strconv.ParseUint(greenString, 16, 8)
	blueString := lineColor[5:7]
	blue, _ := strconv.ParseUint(blueString, 16, 8)
	a, _ := strconv.ParseUint("ff", 16, 8)
	goColor := color.RGBA{
		uint8(red),
		uint8(green),
		uint8(blue),
		uint8(a),
	}
	return goColor
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
