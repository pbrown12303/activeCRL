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

type ToolbarSelection int

const (
	CURSOR ToolbarSelection = iota
	ELEMENT
	LITERAL
	REFERENCE
	REFERENCE_LINK
	REFINEMENT
	REFINEMENT_LINK
	OWNER_POINTER
	REFERENCED_ELEMENT_POINTER
	ABSTRACT_ELEMENT_POINTER
	REFINED_ELEMENT_POINTER
)

func (selection ToolbarSelection) ToString() string {
	switch selection {
	case CURSOR:
		return "Cursor"
	case ELEMENT:
		return "Element"
	case LITERAL:
		return "Literal"
	case REFERENCE:
		return "Reference"
	case REFERENCE_LINK:
		return "Reference Link"
	case REFINEMENT:
		return "Refinement"
	case REFINEMENT_LINK:
		return "Refinement Link"
	case OWNER_POINTER:
		return "Owner Poiner"
	case REFERENCED_ELEMENT_POINTER:
		return "Referenced Element Pointer"
	case ABSTRACT_ELEMENT_POINTER:
		return "Abstract Element Pointer"
	case REFINED_ELEMENT_POINTER:
		return "RefinedElementPointer"
	}
	return ""
}

type diagramTab struct {
	diagramID string
	tab       *container.TabItem
	diagram   *diagramwidget.DiagramWidget
}

// FyneDiagramManager manages the relationship between the fyne DiagramWidgets and the
// underlying CRL model. It is a component of the  FyneGUI
type FyneDiagramManager struct {
	fyneGUI                         *CrlEditorFyneGUI
	diagramArea                     *fyne.Container
	diagramTabs                     map[string]*diagramTab
	toolbar                         *fyne.Container
	toolButtons                     map[ToolbarSelection]*widget.Button
	tabArea                         *container.DocTabs
	diagramObserver                 *diagramObserver
	diagramElementObserver          *diagramElementObserver
	currentToolbarSelection         ToolbarSelection
	inProgressCreateLinkTransaction *core.Transaction
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
	dm.currentToolbarSelection = CURSOR
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
	diagramLink := NewFyneCrlDiagramLink(diagramWidget, link, trans)
	diagramLink.SetSourcePad(fyneSourcePad)
	diagramLink.SetTargetPad(fyneTargetPad)
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
		if diagram != nil {
			diagram.Deregister(dm.diagramObserver)
		}
	}
}

func (dm *FyneDiagramManager) createToolbar() {
	dm.toolbar = container.NewVBox()
	dm.toolButtons = make(map[ToolbarSelection]*widget.Button)
	// Cursor
	button := widget.NewButtonWithIcon("", images.ResourceCursorIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = CURSOR
	}
	dm.toolButtons[CURSOR] = button
	dm.toolbar.Add(button)
	// Element
	button = widget.NewButtonWithIcon("", images.ResourceElementIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = ELEMENT
	}
	dm.toolButtons[ELEMENT] = button
	dm.toolbar.Add(button)
	// Literal
	button = widget.NewButtonWithIcon("", images.ResourceLiteralIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = LITERAL
	}
	dm.toolButtons[LITERAL] = button
	dm.toolbar.Add(button)
	// Reference
	button = widget.NewButtonWithIcon("", images.ResourceReferenceIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = REFERENCE
	}
	dm.toolButtons[REFERENCE] = button
	dm.toolbar.Add(button)
	// ReferenceLink
	button = widget.NewButtonWithIcon("", images.ResourceReferenceLinkIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = REFERENCE_LINK
		dm.startCreateLinkTransaction()
	}
	dm.toolButtons[REFERENCE_LINK] = button
	dm.toolbar.Add(button)
	// Refinement
	button = widget.NewButtonWithIcon("", images.ResourceRefinementIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = REFINEMENT
	}
	dm.toolButtons[REFINEMENT] = button
	dm.toolbar.Add(button)
	// RefinementLink
	button = widget.NewButtonWithIcon("", images.ResourceRefinementLinkIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = REFINEMENT_LINK
		dm.startCreateLinkTransaction()
	}
	dm.toolButtons[REFINEMENT_LINK] = button
	dm.toolbar.Add(button)
	// OwnerPointer
	button = widget.NewButtonWithIcon("", images.ResourceOwnerPointerIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = OWNER_POINTER
		dm.startCreateLinkTransaction()
	}
	dm.toolButtons[OWNER_POINTER] = button
	dm.toolbar.Add(button)
	// ReferencedElementPointer
	button = widget.NewButtonWithIcon("", images.ResourceElementPointerIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = REFERENCED_ELEMENT_POINTER
		dm.startCreateLinkTransaction()
	}
	dm.toolButtons[REFERENCED_ELEMENT_POINTER] = button
	dm.toolbar.Add(button)
	// AbstractPointer
	button = widget.NewButtonWithIcon("", images.ResourceAbstractPointerIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = ABSTRACT_ELEMENT_POINTER
		dm.startCreateLinkTransaction()
	}
	dm.toolButtons[ABSTRACT_ELEMENT_POINTER] = button
	dm.toolbar.Add(button)
	// RefinedPointer
	button = widget.NewButtonWithIcon("", images.ResourceRefinedPointerIconPng, nil)
	button.OnTapped = func() {
		dm.currentToolbarSelection = REFINED_ELEMENT_POINTER
		dm.startCreateLinkTransaction()
	}
	dm.toolButtons[REFINED_ELEMENT_POINTER] = button
	dm.toolbar.Add(button)
}

func (dm *FyneDiagramManager) diagramElementSelectionChanged(diagramElementID string) {
	editor := dm.fyneGUI.editor
	trans, new := editor.GetTransaction()
	if new {
		defer editor.EndTransaction()
	}
	dm.fyneGUI.editor.SelectElementUsingIDString(diagramElementID, trans)
}

func (dm *FyneDiagramManager) diagramTapped(fyneDiagram *diagramwidget.DiagramWidget, event *fyne.PointEvent) {
	trans, new := dm.fyneGUI.editor.GetTransaction()
	if new {
		defer dm.fyneGUI.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	crlDiagram := uOfD.GetElement(fyneDiagram.ID)
	var el core.Element
	switch dm.currentToolbarSelection {
	case CURSOR:
		fyneDiagram.ClearSelection()
	case ELEMENT:
		el, _ = uOfD.NewElement(trans)
		el.SetLabel(dm.fyneGUI.editor.GetDefaultElementLabel(), trans)
	case LITERAL:
		el, _ = uOfD.NewLiteral(trans)
		el.SetLabel(dm.fyneGUI.editor.GetDefaultLiteralLabel(), trans)
	case REFERENCE:
		el, _ = uOfD.NewReference(trans)
		el.SetLabel(dm.fyneGUI.editor.GetDefaultReferenceLabel(), trans)
	case REFINEMENT:
		el, _ = uOfD.NewRefinement(trans)
		el.SetLabel(dm.fyneGUI.editor.GetDefaultRefinementLabel(), trans)
	}

	if el != nil {
		el.SetOwningConceptID(crlDiagram.GetOwningConceptID(trans), trans)
		dm.fyneGUI.editor.SelectElement(el, trans)

		// Now the view
		x := event.Position.X
		y := event.Position.Y
		var newNode core.Element
		newNode, _ = crldiagramdomain.NewDiagramNode(uOfD, trans)
		newNode.Register(dm.diagramElementObserver)
		crldiagramdomain.SetNodeX(newNode, float64(x), trans)
		crldiagramdomain.SetNodeY(newNode, float64(y), trans)
		newNode.SetLabel(el.GetLabel(trans), trans)
		crldiagramdomain.SetReferencedModelElement(newNode, el, trans)
		crldiagramdomain.SetDisplayLabel(newNode, el.GetLabel(trans), trans)
		newNode.SetOwningConcept(crlDiagram, trans)
	}

	dm.ElementSelected("", trans)
	dm.currentToolbarSelection = CURSOR

}

func (dm *FyneDiagramManager) displayDiagram(diagram core.Element, trans *core.Transaction) error {
	diagramID := diagram.GetConceptID(trans)
	tabItem := dm.diagramTabs[diagramID]
	if tabItem == nil {
		diagramWidget := diagramwidget.NewDiagramWidget(diagramID)
		diagramWidget.OnTappedCallback = dm.diagramTapped
		scrollingContainer := container.NewScroll(diagramWidget)
		tabItem = &diagramTab{
			diagramID: diagramID,
			tab:       container.NewTabItem(diagram.GetLabel(trans), scrollingContainer),
			diagram:   diagramWidget,
		}
		dm.diagramTabs[diagramID] = tabItem
		dm.tabArea.Append(tabItem.tab)
		diagram.Register(dm.diagramObserver)
		dm.populateDiagram(diagram, trans)
		diagramWidget.LinkConnectionChangedCallback = func(link diagramwidget.DiagramLink, end string, oldPad diagramwidget.ConnectionPad, newPad diagramwidget.ConnectionPad) {
			dm.linkConnectionChanged(link, end, oldPad, newPad)
		}
		diagramWidget.PrimaryDiagramElementSelectionChangedCallback = func(id string) {
			dm.diagramElementSelectionChanged(id)
		}
		diagramWidget.IsConnectionAllowedCallback = func(link diagramwidget.DiagramLink, linkEnd diagramwidget.LinkEnd, pad diagramwidget.ConnectionPad) bool {
			return dm.isConnectionAllowed(link, linkEnd, pad)
		}
	}
	dm.tabArea.Select(tabItem.tab)
	return nil
}

func (dm *FyneDiagramManager) SelectDiagram(diagramID string) {
	tabItem := dm.diagramTabs[diagramID]
	dm.tabArea.Select(tabItem.tab)
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

func (dm *FyneDiagramManager) GetSelectedDiagram() *diagramwidget.DiagramWidget {
	selectedTabItem := dm.tabArea.Selected()
	for _, diagramTab := range dm.diagramTabs {
		if diagramTab.tab == selectedTabItem {
			return diagramTab.diagram
		}
	}
	return nil
}

func (dm *FyneDiagramManager) initialize() {
	diagramIDs := []string{}
	for _, diagramTab := range dm.diagramTabs {
		diagramIDs = append(diagramIDs, diagramTab.diagramID)
	}
	for _, diagramID := range diagramIDs {
		dm.closeDiagram(diagramID)
	}
}

// isConnectionAllowed is the callback function for determining acceptable link connections
func (dm *FyneDiagramManager) isConnectionAllowed(fyneLink diagramwidget.DiagramLink, linkEnd diagramwidget.LinkEnd, pad diagramwidget.ConnectionPad) bool {
	trans, new := dm.fyneGUI.editor.GetTransaction()
	if new {
		defer dm.fyneGUI.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	link := uOfD.GetElement(fyneLink.GetDiagramElementID())
	padOwner := uOfD.GetElement(pad.GetPadOwner().GetDiagramElementID())
	if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramReferenceLinkURI, trans) {
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramAbstractPointerURI, trans) {
		padOwnerModelElement := crldiagramdomain.GetReferencedModelElement(padOwner, trans)
		if padOwnerModelElement == nil {
			return false
		}
		switch linkEnd {
		case diagramwidget.SOURCE:
			switch padOwnerModelElement.(type) {
			case core.Refinement:
				return true
			}
			return false
		case diagramwidget.TARGET:
			return true
		}
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementPointerURI, trans) {
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramOwnerPointerURI, trans) {
		switch linkEnd {
		case diagramwidget.SOURCE:
			return true
		case diagramwidget.TARGET:
			if padOwner != crldiagramdomain.GetLinkSource(link, trans) {
				// an element cannot own itself
				return true
			}
		}
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinedPointerURI, trans) {
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinementLinkURI, trans) {
	}
	return false
}

// linkConnectionChanged is the callback for changes in link connections
func (dm *FyneDiagramManager) linkConnectionChanged(link diagramwidget.DiagramLink, end string, oldPad diagramwidget.ConnectionPad, newPad diagramwidget.ConnectionPad) error {
	switch typedLink := link.(type) {
	case *FyneCrlDiagramLink:
		trans, new := dm.fyneGUI.editor.GetTransaction()
		if new {
			defer dm.fyneGUI.editor.EndTransaction()
		}
		uOfD := trans.GetUniverseOfDiscourse()
		crlLink := uOfD.GetElement(link.GetDiagramElementID())
		if crlLink == nil {
			return errors.New("in FyneDiagramManager.linkConnectionChanged CrlLink not found")
		}
		crlNewPadOwner := uOfD.GetElement(newPad.GetPadOwner().GetDiagramElementID())
		if crlNewPadOwner == nil {
			return errors.New("in FyneDiagramManager.linkConnectionChanged CrlLink not found")
		}
		switch end {
		case "source":
			crldiagramdomain.SetLinkSource(crlLink, crlNewPadOwner, trans)
			switch typedLink.linkType {
			case REFERENCE_LINK:
			case REFINEMENT_LINK:
			case ABSTRACT_ELEMENT_POINTER:
			case OWNER_POINTER:
				currentLinkParent := crldiagramdomain.GetReferencedModelElement(crlLink, trans)
				newLinkParent := crldiagramdomain.GetReferencedModelElement(crlNewPadOwner, trans)
				if currentLinkParent != newLinkParent {
					if currentLinkParent != nil {
						currentLinkParent.SetOwningConcept(nil, trans)
					}
					crlLinkTarget := crldiagramdomain.GetLinkTarget(crlLink, trans)
					targetModelElement := crldiagramdomain.GetReferencedModelElement(crlLinkTarget, trans)
					newLinkParent.SetOwningConcept(targetModelElement, trans)
					crldiagramdomain.SetReferencedModelElement(crlLink, newLinkParent, trans)
				}
			case REFERENCED_ELEMENT_POINTER:
			case REFINED_ELEMENT_POINTER:
			}
		case "target":
			crldiagramdomain.SetLinkTarget(crlLink, crlNewPadOwner, trans)
			switch typedLink.linkType {
			case REFERENCE_LINK:
			case REFINEMENT_LINK:
			case ABSTRACT_ELEMENT_POINTER:
			case OWNER_POINTER:
				linkParent := crldiagramdomain.GetReferencedModelElement(crlLink, trans)
				targetModelElement := crldiagramdomain.GetReferencedModelElement(crlNewPadOwner, trans)
				if linkParent != nil && linkParent.GetOwningConcept(trans) != targetModelElement {
					linkParent.SetOwningConcept(targetModelElement, trans)
				}
			case REFERENCED_ELEMENT_POINTER:
			case REFINED_ELEMENT_POINTER:
			}
		}
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
	diagram.ClearSelectionNoCallback()
	foundDiagramElementID := ""
	for key := range diagram.GetDiagramElements() {
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

func (dm *FyneDiagramManager) startCreateLinkTransaction() {
	currentDiagram := dm.GetSelectedDiagram()
	if currentDiagram == nil {
		// nothing to do
		return
	}
	// Only start if the current toolbar selection is for a link or pointer
	switch dm.currentToolbarSelection {
	case REFINEMENT_LINK, REFERENCE_LINK, ABSTRACT_ELEMENT_POINTER, OWNER_POINTER, REFERENCED_ELEMENT_POINTER, REFINED_ELEMENT_POINTER:
		trans, new := dm.fyneGUI.editor.GetTransaction()
		if new {
			dm.inProgressCreateLinkTransaction = trans
		}
		uOfD := trans.GetUniverseOfDiscourse()
		var crlLink core.Element
		var fyneLink diagramwidget.DiagramLink
		switch dm.currentToolbarSelection {
		case REFINEMENT_LINK:
			crlLink, _ = crldiagramdomain.NewDiagramRefinementLink(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case REFERENCE_LINK:
			crlLink, _ = crldiagramdomain.NewDiagramReferenceLink(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case ABSTRACT_ELEMENT_POINTER:
			crlLink, _ = crldiagramdomain.NewDiagramAbstractPointer(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case OWNER_POINTER:
			crlLink, _ = crldiagramdomain.NewDiagramOwnerPointer(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case REFERENCED_ELEMENT_POINTER:
			crlLink, _ = crldiagramdomain.NewDiagramElementPointer(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case REFINED_ELEMENT_POINTER:
			crlLink, _ = crldiagramdomain.NewDiagramRefinedPointer(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		}
		crlDiagram := uOfD.GetElement(currentDiagram.ID)
		crlLink.SetOwningConcept(crlDiagram, trans)
		currentDiagram.StartNewLinkConnectionTransaction(fyneLink)
	}
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
	if crlDiagramElement != nil {
		fyneDiagramElement := diagramWidget.GetDiagramElement(elementID)
		if fyneDiagramElement == nil || reflect.ValueOf(fyneDiagramElement).IsNil() {
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
