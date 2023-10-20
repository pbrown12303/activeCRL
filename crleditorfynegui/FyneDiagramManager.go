package crleditorfynegui

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
	"strconv"

	"fyne.io/x/fyne/widget/diagramwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crlmapsdomain"
	"github.com/pbrown12303/activeCRL/images"

	mapset "github.com/deckarep/golang-set"
)

const (
	displayLabel = "DisplayLabel"
)

// ToolbarSelection is the enumeration of toolbar items
type ToolbarSelection int

// The enumerated values ToolbarSelection
const (
	CursorSelected ToolbarSelection = iota
	ElementSelected
	LiteralSelected
	ReferenceSelected
	ReferenceLinkSelected
	RefinementSelected
	RefinementLinkSelected
	OwnerPointerSelected
	ReferencedElementPointerSelected
	AbstractElementPointerSelected
	RefinedElementPointerSelected
	OneToOneMapSelected
	CreateRefinementOfConceptSelected
)

// ToString retuns text identifying the ToolbarSelection
func (selection ToolbarSelection) ToString() string {
	switch selection {
	case CursorSelected:
		return "Cursor Selected"
	case ElementSelected:
		return "Element Selected"
	case LiteralSelected:
		return "Literal Selected"
	case ReferenceSelected:
		return "Reference Selected"
	case ReferenceLinkSelected:
		return "Reference Link Selected"
	case RefinementSelected:
		return "Refinement Selected"
	case RefinementLinkSelected:
		return "Refinement Link Selected"
	case OwnerPointerSelected:
		return "Owner Poiner Selected"
	case ReferencedElementPointerSelected:
		return "Referenced Element Pointer Selected"
	case AbstractElementPointerSelected:
		return "Abstract Element Pointer Selected"
	case RefinedElementPointerSelected:
		return "RefinedElementPointer Selected"
	case OneToOneMapSelected:
		return "OneToOne Map Selected"
	case CreateRefinementOfConceptSelected:
		return "Clone Selection As Refinement Selected"
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
	fyneGUI                                *CrlEditorFyneGUI
	diagramArea                            *fyne.Container
	diagramTabs                            map[string]*diagramTab
	toolbar                                *fyne.Container
	toolButtons                            map[ToolbarSelection]*widget.Button
	tabArea                                *container.DocTabs
	diagramObserver                        *diagramObserver
	diagramElementObserver                 *diagramElementObserver
	currentToolbarSelection                ToolbarSelection
	connectionTransactionTransientConcepts mapset.Set
}

// NewFyneDiagramManager creates a diagram manager and associates it with the FyneGUI
func NewFyneDiagramManager(fyneGUI *CrlEditorFyneGUI) *FyneDiagramManager {
	var dm FyneDiagramManager
	dm.createToolbar()
	dm.diagramTabs = make(map[string]*diagramTab)
	dm.connectionTransactionTransientConcepts = mapset.NewSet()
	dm.tabArea = container.NewDocTabs()
	dm.tabArea.OnClosed = diagramClosed
	dm.diagramArea = container.NewBorder(nil, nil, dm.toolbar, nil, dm.tabArea)
	dm.diagramObserver = newDiagramObserver(&dm)
	dm.diagramElementObserver = newDiagramElementObserver(&dm)
	dm.fyneGUI = fyneGUI
	dm.currentToolbarSelection = CursorSelected
	dm.toolButtons[CursorSelected].Importance = widget.HighImportance
	dm.toolButtons[CursorSelected].Refresh()
	return &dm
}

func (dm *FyneDiagramManager) addElementToDiagram(element core.Concept, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramElement {
	if element.IsRefinementOfURI(crldiagramdomain.CrlDiagramNodeURI, trans) {
		return dm.addNodeToDiagram(element, trans, diagramWidget)
	} else if element.IsRefinementOfURI(crldiagramdomain.CrlDiagramLinkURI, trans) {
		return dm.addLinkToDiagram(element, trans, diagramWidget)
	}
	return nil
}

func (dm *FyneDiagramManager) addLinkToDiagram(link core.Concept, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) *FyneCrlDiagramLink {
	crlDiagramSource := crldiagramdomain.GetLinkSource(link, trans)
	if crlDiagramSource == nil {
		// Register for changes so that when sufficient information is present we can add it to the diagram
		link.Register(dm.diagramElementObserver)
		return nil
	}
	fyneSource := diagramWidget.GetDiagramElement(crlDiagramSource.GetConceptID(trans))
	if fyneSource == nil {
		// the source is not in the diagram
		return nil
	}
	fyneSourcePad := fyneSource.GetDefaultConnectionPad()
	crlDiagramTarget := crldiagramdomain.GetLinkTarget(link, trans)
	if crlDiagramTarget == nil {
		// Register for changes so that when sufficient information is present we can add it to the diagram
		link.Register(dm.diagramElementObserver)
		return nil
	}
	targetConceptID := crlDiagramTarget.GetConceptID(trans)
	fyneTarget := diagramWidget.GetDiagramElement(targetConceptID)
	if fyneTarget == nil {
		// the target is not in the diagram
		return nil
	}
	fyneTargetPad := fyneTarget.GetDefaultConnectionPad()
	diagramLink := NewFyneCrlDiagramLink(diagramWidget, link, trans)
	diagramLink.SetSourcePad(fyneSourcePad)
	diagramLink.SetTargetPad(fyneTargetPad)
	link.Register(dm.diagramElementObserver)
	diagramWidget.Refresh()
	return diagramLink
}

func (dm *FyneDiagramManager) addNodeToDiagram(node core.Concept, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramNode {
	nodeID := node.GetConceptID(trans)
	diagramNode := diagramWidget.GetDiagramNode(nodeID)
	if diagramNode == nil {

		diagramNode = NewFyneCrlDiagramNode(node, trans, diagramWidget)
		node.Register(dm.diagramElementObserver)
	}
	return diagramNode
}

func (dm *FyneDiagramManager) cancelLinkTransaction() {
	selectedDiagram := dm.GetSelectedDiagram()
	connectionTransaction := selectedDiagram.ConnectionTransaction
	if connectionTransaction != nil {
		selectedDiagram.RemoveElement(connectionTransaction.Link.GetDiagramElementID())
		selectedDiagram.ConnectionTransaction = nil
	}
	trans, isNew := dm.fyneGUI.editor.GetTransaction()
	if isNew {
		defer dm.fyneGUI.editor.EndTransaction()
	}
	trans.GetUniverseOfDiscourse().DeleteElements(dm.connectionTransactionTransientConcepts, trans)
	dm.connectionTransactionTransientConcepts.Clear()
	dm.fyneGUI.windowContent.Refresh()
}

// closeAllDiagrams closes all of the currently displayed diagrams. It is not an undoable operation
func (dm *FyneDiagramManager) closeAllDiagrams() {
	diagramIDs := []string{}
	for _, diagramTab := range dm.diagramTabs {
		diagramIDs = append(diagramIDs, diagramTab.diagramID)
	}
	for _, diagramID := range diagramIDs {
		dm.closeDiagramNoUndo(diagramID)
	}
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
	dm.fyneGUI.editor.DiagramDisplayRemoved(diagramID)
	dm.fyneGUI.editor.DiagramSelected("")
}

func (dm *FyneDiagramManager) closeDiagramNoUndo(diagramID string) {
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

func (dm *FyneDiagramManager) completeLinkTransaction() {
	selectedDiagram := dm.GetSelectedDiagram()
	connectionTransaction := selectedDiagram.ConnectionTransaction
	if connectionTransaction != nil {
		selectedDiagram.ConnectionTransaction = nil
	}
	dm.connectionTransactionTransientConcepts.Clear()
	dm.setToolbarSelection(CursorSelected)
}

func (dm *FyneDiagramManager) createToolbar() {
	dm.toolbar = container.NewVBox()
	dm.toolButtons = make(map[ToolbarSelection]*widget.Button)
	// Cursor
	button := widget.NewButtonWithIcon("", images.ResourceCursorIconPng, func() {
		dm.setToolbarSelection(CursorSelected)
	})
	dm.toolButtons[CursorSelected] = button
	dm.toolbar.Add(button)
	// Element
	button = widget.NewButtonWithIcon("", images.ResourceElementIconPng, func() {
		dm.setToolbarSelection(ElementSelected)
	})
	dm.toolButtons[ElementSelected] = button
	dm.toolbar.Add(button)
	// Literal
	button = widget.NewButtonWithIcon("", images.ResourceLiteralIconPng, func() {
		dm.setToolbarSelection(LiteralSelected)
	})
	dm.toolButtons[LiteralSelected] = button
	dm.toolbar.Add(button)
	// Reference
	button = widget.NewButtonWithIcon("", images.ResourceReferenceIconPng, func() {
		dm.setToolbarSelection(ReferenceSelected)
	})
	dm.toolButtons[ReferenceSelected] = button
	dm.toolbar.Add(button)
	// ReferenceLink
	button = widget.NewButtonWithIcon("", images.ResourceReferenceLinkIconPng, func() {
		dm.setToolbarSelection(ReferenceLinkSelected)
		dm.startCreateLinkTransaction()
	})
	dm.toolButtons[ReferenceLinkSelected] = button
	dm.toolbar.Add(button)
	// Refinement
	button = widget.NewButtonWithIcon("", images.ResourceRefinementIconPng, func() {
		dm.setToolbarSelection(RefinementSelected)
	})
	dm.toolButtons[RefinementSelected] = button
	dm.toolbar.Add(button)
	// RefinementLink
	button = widget.NewButtonWithIcon("", images.ResourceRefinementLinkIconPng, func() {
		dm.setToolbarSelection(RefinementLinkSelected)
		dm.startCreateLinkTransaction()
	})
	dm.toolButtons[RefinementLinkSelected] = button
	dm.toolbar.Add(button)
	// OwnerPointer
	button = widget.NewButtonWithIcon("", images.ResourceOwnerPointerIconPng, func() {
		dm.setToolbarSelection(OwnerPointerSelected)
		dm.startCreateLinkTransaction()
	})
	dm.toolButtons[OwnerPointerSelected] = button
	dm.toolbar.Add(button)
	// ReferencedElementPointer
	button = widget.NewButtonWithIcon("", images.ResourceElementPointerIconPng, func() {
		dm.setToolbarSelection(ReferencedElementPointerSelected)
		dm.startCreateLinkTransaction()
	})
	dm.toolButtons[ReferencedElementPointerSelected] = button
	dm.toolbar.Add(button)
	// AbstractPointer
	button = widget.NewButtonWithIcon("", images.ResourceAbstractPointerIconPng, func() {
		dm.setToolbarSelection(AbstractElementPointerSelected)
		dm.startCreateLinkTransaction()
	})
	dm.toolButtons[AbstractElementPointerSelected] = button
	dm.toolbar.Add(button)
	// RefinedPointer
	button = widget.NewButtonWithIcon("", images.ResourceRefinedPointerIconPng, func() {
		dm.setToolbarSelection(RefinedElementPointerSelected)
		dm.startCreateLinkTransaction()
	})
	dm.toolButtons[RefinedElementPointerSelected] = button
	dm.toolbar.Add(button)
	// Separator
	separator := widget.NewSeparator()
	dm.toolbar.Add(separator)
	// OneToOne Map
	button = widget.NewButtonWithIcon("", images.ResourceOneToOneIconPng, func() {
		dm.setToolbarSelection(OneToOneMapSelected)
	})
	dm.toolButtons[OneToOneMapSelected] = button
	dm.toolbar.Add(button)
	// Clone Selection As Refinement
	button = widget.NewButtonWithIcon("", images.ResourceRefinedCloneIconPng, func() {
		dm.setToolbarSelection(CreateRefinementOfConceptSelected)
	})
	dm.toolButtons[CreateRefinementOfConceptSelected] = button
	dm.toolbar.Add(button)
}

func (dm *FyneDiagramManager) deleteConceptView(elementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	diagramElement := uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("FyneDiagramManager.deleteDiagramElementView diagramElement not found for elementID " + elementID)
	}
	return uOfD.DeleteElement(diagramElement, trans)
}

func (dm *FyneDiagramManager) diagramElementSelectionChanged(diagramElementID string) {
	editor := dm.fyneGUI.editor
	trans, new := editor.GetTransaction()
	if new {
		defer editor.EndTransaction()
	}
	fyneDiagramElement := dm.GetSelectedDiagram().GetDiagramElement(diagramElementID)
	var crlElement core.Concept
	switch typedElement := fyneDiagramElement.(type) {
	case *FyneCrlDiagramNode:
		crlElement = typedElement.modelElement
	case *FyneCrlDiagramLink:
		crlElement = typedElement.modelElement
	}
	dm.fyneGUI.editor.SelectElement(crlElement, trans)
}

func (dm *FyneDiagramManager) diagramMouseMoved(event *desktop.MouseEvent) {
	if dm.fyneGUI.dragDropTransaction != nil {
		dm.fyneGUI.dragDropTransaction.diagramID = dm.GetSelectedDiagram().ID
		dm.fyneGUI.dragDropTransaction.currentDiagramMousePosition = event.Position
	}
}

func (dm *FyneDiagramManager) diagramMouseOut() {
	if dm.fyneGUI.dragDropTransaction != nil {
		dm.fyneGUI.dragDropTransaction.diagramID = ""
		dm.fyneGUI.dragDropTransaction.currentDiagramMousePosition = fyne.NewPos(-1, -1)
	}
}

func (dm *FyneDiagramManager) diagramTapped(fyneDiagram *diagramwidget.DiagramWidget, event *fyne.PointEvent) {
	trans, new := dm.fyneGUI.editor.GetTransaction()
	if new {
		defer dm.fyneGUI.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	crlDiagram := uOfD.GetElement(fyneDiagram.ID)
	var el core.Concept
	switch dm.currentToolbarSelection {
	case CursorSelected:
		fyneDiagram.ClearSelection()
	case ElementSelected:
		uOfD.MarkUndoPoint()
		el, _ = uOfD.NewElement(trans)
		el.SetLabel(dm.fyneGUI.editor.GetDefaultElementLabel(), trans)
	case LiteralSelected:
		uOfD.MarkUndoPoint()
		el, _ = uOfD.NewLiteral(trans)
		el.SetLabel(dm.fyneGUI.editor.GetDefaultLiteralLabel(), trans)
	case ReferenceSelected:
		uOfD.MarkUndoPoint()
		el, _ = uOfD.NewReference(trans)
		el.SetLabel(dm.fyneGUI.editor.GetDefaultReferenceLabel(), trans)
	case RefinementSelected:
		uOfD.MarkUndoPoint()
		el, _ = uOfD.NewRefinement(trans)
		el.SetLabel(dm.fyneGUI.editor.GetDefaultRefinementLabel(), trans)
	case AbstractElementPointerSelected, OwnerPointerSelected, ReferencedElementPointerSelected, ReferenceLinkSelected, RefinedElementPointerSelected, RefinementLinkSelected:
		uOfD.MarkUndoPoint()
	case OneToOneMapSelected:
		sourceMap := uOfD.GetElementWithURI(crlmapsdomain.CrlOneToOneMapURI)
		el, _ = uOfD.CreateReplicateAsRefinement(sourceMap, trans)
		el.SetOwningConcept(crlDiagram.GetOwningConcept(trans), trans)
	case CreateRefinementOfConceptSelected:
		selection := FyneGUISingleton.editor.GetCurrentSelection()
		if selection != nil {
			el, _ = uOfD.CreateRefinementOfConcept(selection, selection.GetLabel(trans), trans)
			el.SetOwningConcept(crlDiagram.GetOwningConcept(trans), trans)
		}
	}

	if el != nil {
		elID := el.GetConceptID(trans)
		el.SetOwningConceptID(crlDiagram.GetOwningConceptID(trans), trans)
		dm.fyneGUI.editor.SelectElement(el, trans)

		// Now the view
		x := event.Position.X
		y := event.Position.Y
		newNode, err := crldiagramdomain.NewDiagramNode(uOfD, trans)
		if err != nil {
			log.Print(err)
			return
		}
		newNode.Register(dm.diagramElementObserver)
		crldiagramdomain.SetNodeX(newNode, float64(x), trans)
		crldiagramdomain.SetNodeY(newNode, float64(y), trans)
		newNode.SetLabel(el.GetLabel(trans), trans)
		crldiagramdomain.SetReferencedModelConcept(newNode, el, trans)
		crldiagramdomain.SetDisplayLabel(newNode, el.GetLabel(trans), trans)
		newNode.SetOwningConcept(crlDiagram, trans)
		dm.selectElementInDiagram(elID, fyneDiagram, trans)
		dm.ElementSelected(elID, trans)
	} else {
		dm.ElementSelected("", trans)
	}
	dm.setToolbarSelection(CursorSelected)
}

func (dm *FyneDiagramManager) displayDiagram(diagram core.Concept, trans *core.Transaction) error {
	diagramID := diagram.GetConceptID(trans)
	tabItem := dm.diagramTabs[diagramID]
	if tabItem == nil {
		diagramWidget := diagramwidget.NewDiagramWidget(diagramID)
		diagramWidget.OnTappedCallback = dm.diagramTapped
		diagramWidget.MouseMovedCallback = dm.diagramMouseMoved
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
		diagramWidget.LinkSegmentMouseDownSecondaryCallback = dm.linkMouseDown
		dm.fyneGUI.editor.DiagramDisplayed(diagramID)
	}
	dm.tabArea.Select(tabItem.tab)
	return nil
}

// ElementSelected selects the element in each diagram (if present)
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

// GetDrawingArea returns the diagram area container
func (dm *FyneDiagramManager) GetDrawingArea() *fyne.Container {
	return dm.diagramArea
}

// GetSelectedDiagram returns the currently selected DiagramWidget
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
	dm.diagramObserver = newDiagramObserver(dm)
	dm.diagramElementObserver = newDiagramElementObserver(dm)
	dm.tabArea.SetItems([]*container.TabItem{})
	dm.tabArea.Select(nil)
	dm.currentToolbarSelection = CursorSelected
	dm.toolButtons[CursorSelected].Importance = widget.HighImportance
	dm.toolButtons[CursorSelected].Refresh()
	dm.connectionTransactionTransientConcepts.Clear()
	dm.closeAllDiagrams()
}

func (dm *FyneDiagramManager) refreshGUI(trans *core.Transaction) {
	diagramIDs := []string{}
	for _, diagramTab := range dm.diagramTabs {
		diagramIDs = append(diagramIDs, diagramTab.diagramID)
	}
	for _, diagramID := range diagramIDs {
		if !dm.fyneGUI.editor.IsDiagramDisplayed(diagramID, trans) {
			dm.closeDiagram(diagramID)
		}
	}
	editor := dm.fyneGUI.editor
	for _, diagramID := range editor.GetSettings().OpenDiagrams {
		editor.GetDiagramManager().DisplayDiagram(diagramID, trans)
		diagram := dm.getDiagramWidget(diagramID)
		for _, diagramElement := range diagram.GetDiagramElements() {
			diagram.RemoveElement(diagramElement.GetDiagramElementID())
		}
		crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagramID)
		dm.populateDiagram(crlDiagram, trans)
		dm.selectElementInDiagram(editor.GetSettings().Selection, diagram, trans)
		diagram.Refresh()
	}
	dm.SelectDiagram(editor.GetSettings().CurrentDiagram)
}

// isConnectionAllowed is the callback function for determining acceptable link connections
func (dm *FyneDiagramManager) isConnectionAllowed(fyneLink diagramwidget.DiagramLink, linkEnd diagramwidget.LinkEnd, pad diagramwidget.ConnectionPad) bool {
	trans, new := dm.fyneGUI.editor.GetTransaction()
	if new {
		defer dm.fyneGUI.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	crlLink := uOfD.GetElement(fyneLink.GetDiagramElementID())
	crlPadOwner := uOfD.GetElement(pad.GetPadOwner().GetDiagramElementID())
	if crlLink.IsRefinementOfURI(crldiagramdomain.CrlDiagramReferenceLinkURI, trans) {
		return true
	} else if crlLink.IsRefinementOfURI(crldiagramdomain.CrlDiagramAbstractPointerURI, trans) {
		padOwnerModelElement := crldiagramdomain.GetReferencedModelConcept(crlPadOwner, trans)
		if padOwnerModelElement == nil {
			return false
		}
		switch linkEnd {
		case diagramwidget.SOURCE:
			switch padOwnerModelElement.GetConceptType() {
			case core.Refinement:
				return true
			}
			return false
		case diagramwidget.TARGET:
			if crlPadOwner.IsRefinementOfURI(crldiagramdomain.CrlDiagramPointerURI, trans) {
				return false
			}
			return true
		}
	} else if crlLink.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementPointerURI, trans) {
		padOwnerModelElement := crldiagramdomain.GetReferencedModelConcept(crlPadOwner, trans)
		if padOwnerModelElement == nil {
			return false
		}
		switch linkEnd {
		case diagramwidget.SOURCE:
			switch padOwnerModelElement.GetConceptType() {
			case core.Reference:
				return true
			}
			return false
		case diagramwidget.TARGET:
			return true
		}
	} else if crlLink.IsRefinementOfURI(crldiagramdomain.CrlDiagramOwnerPointerURI, trans) {
		switch linkEnd {
		case diagramwidget.SOURCE:
			return true
		case diagramwidget.TARGET:
			if crlPadOwner.IsRefinementOfURI(crldiagramdomain.CrlDiagramPointerURI, trans) {
				return false
			}
			if crlPadOwner != crldiagramdomain.GetLinkSource(crlLink, trans) {
				// an element cannot own itself
				return true
			}
		}
	} else if crlLink.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinedPointerURI, trans) {
		padOwnerModelElement := crldiagramdomain.GetReferencedModelConcept(crlPadOwner, trans)
		if padOwnerModelElement == nil {
			return false
		}
		switch linkEnd {
		case diagramwidget.SOURCE:
			switch padOwnerModelElement.GetConceptType() {
			case core.Refinement:
				return true
			}
			return false
		case diagramwidget.TARGET:
			if crlPadOwner.IsRefinementOfURI(crldiagramdomain.CrlDiagramPointerURI, trans) {
				return false
			}
			return true
		}
	} else if crlLink.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinementLinkURI, trans) {
		return !crlPadOwner.IsRefinementOfURI(crldiagramdomain.CrlDiagramPointerURI, trans)
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
			case ReferenceLinkSelected:
				linkModelElement := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				sourceModelElement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				linkModelElement.SetOwningConcept(sourceModelElement, trans)
				link.Show()
			case RefinementLinkSelected:
				linkModelElement := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				sourceModelElement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				switch linkModelElement.GetConceptType() {
				case core.Refinement:
					linkModelElement.SetOwningConcept(sourceModelElement, trans)
					linkModelElement.SetRefinedConcept(sourceModelElement, trans)
					link.Show()
				}
			case AbstractElementPointerSelected:
				currentModelRefinement := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				newModelRefinement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				if currentModelRefinement != newModelRefinement {
					if currentModelRefinement != nil {
						currentModelRefinement.SetAbstractConcept(nil, trans)
					}
					crlLinkTarget := crldiagramdomain.GetLinkTarget(crlLink, trans)
					targetModelElement := crldiagramdomain.GetReferencedModelConcept(crlLinkTarget, trans)
					switch newModelRefinement.GetConceptType() {
					case core.Refinement:
						newModelRefinement.SetAbstractConcept(targetModelElement, trans)
						crldiagramdomain.SetReferencedModelConcept(crlLink, newModelRefinement, trans)
						typedLink.modelElement = newModelRefinement
					}
				}
			case OwnerPointerSelected:
				currentLinkModelConcept := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				newLinkModelConcept := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				if currentLinkModelConcept != newLinkModelConcept {
					if currentLinkModelConcept != nil {
						currentLinkModelConcept.SetOwningConcept(nil, trans)
					}
					crlLinkTarget := crldiagramdomain.GetLinkTarget(crlLink, trans)
					targetModelElement := crldiagramdomain.GetReferencedModelConcept(crlLinkTarget, trans)
					newLinkModelConcept.SetOwningConcept(targetModelElement, trans)
					crldiagramdomain.SetReferencedModelConcept(crlLink, newLinkModelConcept, trans)
					typedLink.modelElement = newLinkModelConcept
				}
			case ReferencedElementPointerSelected:
				currentModelReference := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				newModelReference := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				if currentModelReference != newModelReference {
					crldiagramdomain.SetReferencedModelConcept(crlLink, newModelReference, trans)
					attributeName := core.NoAttribute
					if currentModelReference != nil {
						switch currentModelReference.GetConceptType() {
						case core.Reference:
							attributeName = currentModelReference.GetReferencedAttributeName(trans)
							currentModelReference.SetReferencedConcept(nil, core.NoAttribute, trans)
						}
					}
					crlLinkTarget := crldiagramdomain.GetLinkTarget(crlLink, trans)
					targetModelElement := crldiagramdomain.GetReferencedModelConcept(crlLinkTarget, trans)
					switch newModelReference.GetConceptType() {
					case core.Reference:
						newModelReference.SetReferencedConcept(targetModelElement, attributeName, trans)
						crldiagramdomain.SetReferencedModelConcept(crlLink, newModelReference, trans)
						typedLink.modelElement = newModelReference
					}
				}
			case RefinedElementPointerSelected:
				currentModelRefinement := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				newModelRefinement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				if currentModelRefinement != newModelRefinement {
					if currentModelRefinement != nil {
						switch currentModelRefinement.GetConceptType() {
						case core.Refinement:
							currentModelRefinement.SetRefinedConcept(nil, trans)
						}
					}
					crlLinkTarget := crldiagramdomain.GetLinkTarget(crlLink, trans)
					targetModelElement := crldiagramdomain.GetReferencedModelConcept(crlLinkTarget, trans)
					switch newModelRefinement.GetConceptType() {
					case core.Refinement:
						newModelRefinement.SetRefinedConcept(targetModelElement, trans)
						crldiagramdomain.SetReferencedModelConcept(crlLink, newModelRefinement, trans)
						typedLink.modelElement = newModelRefinement
					}
				}
			}
		case "target":
			crldiagramdomain.SetLinkTarget(crlLink, crlNewPadOwner, trans)
			switch typedLink.linkType {
			case ReferenceLinkSelected:
				linkModelElement := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				targetModelElement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				newPadOwner := newPad.GetPadOwner()
				attributeName := getAttributeNameBasedOnTargetType(newPadOwner)
				switch linkModelElement.GetConceptType() {
				case core.Reference:
					linkModelElement.SetReferencedConcept(targetModelElement, attributeName, trans)
				}
			case RefinementLinkSelected:
				linkModelElement := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				targetModelElement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				switch linkModelElement.GetConceptType() {
				case core.Refinement:
					linkModelElement.SetAbstractConcept(targetModelElement, trans)
				}
			case AbstractElementPointerSelected:
				crlModelRefinement := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				currentAbstractElement := crlModelRefinement.GetAbstractConcept(trans)
				newAbstractElement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				if currentAbstractElement != newAbstractElement {
					switch crlModelRefinement.GetConceptType() {
					case core.Refinement:
						crlModelRefinement.SetAbstractConcept(newAbstractElement, trans)
					}
				}
			case OwnerPointerSelected:
				crlLinkParent := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				targetModelElement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				if crlLinkParent != nil && crlLinkParent.GetOwningConcept(trans) != targetModelElement {
					crlLinkParent.SetOwningConcept(targetModelElement, trans)
				}
			case ReferencedElementPointerSelected:
				crlModelReference := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				currentReferencedElement := crlModelReference.GetReferencedConcept(trans)
				newReferencedElement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				if currentReferencedElement != newReferencedElement {
					attributeName := getAttributeNameBasedOnTargetType(newPad.GetPadOwner())
					switch crlModelReference.GetConceptType() {
					case core.Reference:
						crlModelReference.SetReferencedConcept(newReferencedElement, attributeName, trans)
					}
				}
			case RefinedElementPointerSelected:
				crlModelRefinement := crldiagramdomain.GetReferencedModelConcept(crlLink, trans)
				currentRefinedElement := crlModelRefinement.GetRefinedConcept(trans)
				newRefinedElement := crldiagramdomain.GetReferencedModelConcept(crlNewPadOwner, trans)
				if currentRefinedElement != newRefinedElement {
					switch crlModelRefinement.GetConceptType() {
					case core.Refinement:
						crlModelRefinement.SetRefinedConcept(newRefinedElement, trans)

					}
				}
			}
			dm.completeLinkTransaction()
		}
	}
	return nil
}

func getAttributeNameBasedOnTargetType(newPadOwner diagramwidget.DiagramElement) core.AttributeName {
	var attributeName core.AttributeName = core.NoAttribute
	if newPadOwner == nil {
		return attributeName
	}
	typedPadOwner := newPadOwner.GetDiagram().GetDiagramLink(newPadOwner.GetDiagramElementID())
	switch castPadOwner := typedPadOwner.(type) {
	case *FyneCrlDiagramLink:
		switch castPadOwner.linkType {
		case OwnerPointerSelected:
			attributeName = core.OwningConceptID
		case ReferencedElementPointerSelected:
			attributeName = core.ReferencedConceptID
		case AbstractElementPointerSelected:
			attributeName = core.AbstractConceptID
		case RefinedElementPointerSelected:
			attributeName = core.RefinedConceptID
		}
	}
	return attributeName
}

func (dm *FyneDiagramManager) linkMouseDown(link diagramwidget.DiagramLink, event *desktop.MouseEvent) {
	switch typedLink := link.(type) {
	case *FyneCrlDiagramLink:
		ShowSecondaryPopup(typedLink, event)
	}
}

func (dm *FyneDiagramManager) nullifyReferencedConcept(fcde FyneCrlDiagramElement) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	ref := fcde.GetModelElement()
	if ref == nil {
		return errors.New("FyneDiagramManager.nullifyReferencedConcept called with nil model element")
	}
	switch ref.GetConceptType() {
	case core.Reference:
		err := ref.SetReferencedConceptID("", core.NoAttribute, trans)
		if err != nil {
			return errors.Wrap(err, "FyneDiagramManager.nullifyReferencedConcept failed")
		}
	}
	return nil
}

// populateDiagram adds all elements to the diagram
func (dm *FyneDiagramManager) populateDiagram(diagram core.Concept, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	diagramWidget := dm.getDiagramWidget(diagram.GetConceptID(trans))
	nodes := diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, trans)
	for _, node := range nodes {
		dm.addNodeToDiagram(node, trans, diagramWidget)
	}
	links := diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramLinkURI, trans)
	// Since links may have other links as source or target, the source or target may not have been added to the
	// diagram yet. DeferredLinkIDs keeps track of those that were not added due to the absence of the source or target
	desiredLinkIDs := mapset.NewSet()
	for _, link := range links {
		desiredLinkIDs.Add(link.GetConceptID(trans))
	}
	workingLinkIDs := desiredLinkIDs.Clone()
	for workingLinkIDs.Cardinality() > 0 {
		deferredLinkIDs := mapset.NewSet()
		workingIterator := workingLinkIDs.Iterator()
		for entry := range workingIterator.C {
			linkID := entry.(string)
			diagramLink := diagramWidget.GetDiagramLink(linkID)
			if diagramLink == nil {
				link := uOfD.GetElement(linkID)
				addedLink := dm.addLinkToDiagram(link, trans, diagramWidget)
				if addedLink == nil {
					deferredLinkIDs.Add(linkID)
				}
			}
		}
		workingLinkIDs = deferredLinkIDs.Clone()
	}
	diagramWidget.Refresh()
	return nil
}

// SelectDiagram selects the tab whose diagram has the indicated ID
func (dm *FyneDiagramManager) SelectDiagram(diagramID string) {
	tabItem := dm.diagramTabs[diagramID]
	if tabItem != nil {
		dm.tabArea.Select(tabItem.tab)
	}
	dm.fyneGUI.editor.DiagramSelected(diagramID)
}

func (dm *FyneDiagramManager) selectElementInDiagram(elementID string, diagram *diagramwidget.DiagramWidget, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	diagram.ClearSelectionNoCallback()
	if elementID == "" {
		return nil
	}
	foundDiagramElementID := ""
	for _, fyneDiagramElement := range diagram.GetDiagramElements() {
		id := fyneDiagramElement.GetDiagramElementID()
		crlDiagramElement := uOfD.GetElement(id)
		if crlDiagramElement != nil {
			crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
			if crlModelElement != nil {
				if crlModelElement.GetConceptID(trans) == elementID {
					foundDiagramElementID = id
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

func (dm *FyneDiagramManager) showOwnedConcepts(elementID string, recursive bool, skipRefinements bool) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	return dm.showOwnedConceptsImpl(uOfD, elementID, recursive, skipRefinements, trans)
}

func (dm *FyneDiagramManager) showOwnedConceptsImpl(uOfD *core.UniverseOfDiscourse, elementID string, recursive bool, skipRefinements bool, trans *core.Transaction) error {
	diagramElement := uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showOwnedConcepts diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(trans)
	if diagram == nil {
		return errors.New("diagramManager.showOwnedConcepts diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelConcept(diagramElement, trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showOwnedConcepts modelConcept not found for elementID " + elementID)
	}
	it := modelConcept.GetOwnedConceptIDs(trans).Iterator()
	var xOffset float64 = 0
	xPreferencesOffset := FyneGUISingleton.editor.GetUserPreferences().HorizontalLayoutSpacing
	yPreferencesOffset := FyneGUISingleton.editor.GetUserPreferences().VerticalLayoutSpacing
	for id := range it.C {
		child := uOfD.GetElement(id.(string))
		if child == nil {
			return errors.New("Child Concept is nil for id " + id.(string))
		}
		if skipRefinements && child.GetConceptType() == core.Refinement {
			continue
		}
		diagramChildConcept := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, child, trans)
		if diagramChildConcept == nil {
			diagramChildConcept, _ = crldiagramdomain.NewDiagramNode(uOfD, trans)
			crldiagramdomain.SetReferencedModelConcept(diagramChildConcept, child, trans)
			crldiagramdomain.SetDisplayLabel(diagramChildConcept, child.GetLabel(trans), trans)
			diagramElementX := crldiagramdomain.GetNodeX(diagramElement, trans)
			diagramElementY := crldiagramdomain.GetNodeY(diagramElement, trans)
			diagramElementHeight := crldiagramdomain.GetNodeHeight(diagramElement, trans)
			crldiagramdomain.SetNodeX(diagramChildConcept, diagramElementX+xOffset, trans)
			crldiagramdomain.SetNodeY(diagramChildConcept, diagramElementY+diagramElementHeight+yPreferencesOffset, trans)
			diagramChildConcept.SetOwningConcept(diagram, trans)
			xOffset = xOffset + xPreferencesOffset + crldiagramdomain.GetNodeWidth(diagramChildConcept, trans)
		}
		ownerPointer := crldiagramdomain.GetOwnerPointer(diagram, diagramElement, trans)
		if ownerPointer == nil {
			ownerPointer, _ = crldiagramdomain.NewDiagramOwnerPointer(uOfD, trans)
			crldiagramdomain.SetReferencedModelConcept(ownerPointer, child, trans)
			crldiagramdomain.SetLinkSource(ownerPointer, diagramChildConcept, trans)
			crldiagramdomain.SetLinkTarget(ownerPointer, diagramElement, trans)
			ownerPointer.SetOwningConcept(diagram, trans)
		}
		if recursive {
			dm.showOwnedConceptsImpl(uOfD, diagramChildConcept.GetConceptID(trans), recursive, skipRefinements, trans)
		}
	}
	return nil
}

func (dm *FyneDiagramManager) showOwner(elementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	diagramElement := uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showOwner diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(trans)
	if diagram == nil {
		return errors.New("diagramManager.showOwner diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelConcept(diagramElement, trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showOwner modelConcept not found for elementID " + elementID)
	}
	modelConceptOwner := modelConcept.GetOwningConcept(trans)
	if modelConceptOwner == nil {
		return errors.New("Owner is nil")
	}
	diagramConceptOwner := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelConceptOwner, trans)
	if diagramConceptOwner == nil {
		diagramConceptOwner, _ = crldiagramdomain.NewDiagramNode(uOfD, trans)
		crldiagramdomain.SetReferencedModelConcept(diagramConceptOwner, modelConceptOwner, trans)
		crldiagramdomain.SetDisplayLabel(diagramConceptOwner, modelConceptOwner.GetLabel(trans), trans)
		diagramElementX := crldiagramdomain.GetNodeX(diagramElement, trans)
		diagramElementY := crldiagramdomain.GetNodeY(diagramElement, trans)
		crldiagramdomain.SetNodeX(diagramConceptOwner, diagramElementX, trans)
		crldiagramdomain.SetNodeY(diagramConceptOwner, diagramElementY-100, trans)
		diagramConceptOwner.SetOwningConcept(diagram, trans)
	}
	ownerPointer := crldiagramdomain.GetOwnerPointer(diagram, diagramElement, trans)
	if ownerPointer == nil {
		ownerPointer, _ = crldiagramdomain.NewDiagramOwnerPointer(uOfD, trans)
		crldiagramdomain.SetReferencedModelConcept(ownerPointer, modelConcept, trans)
		crldiagramdomain.SetLinkSource(ownerPointer, diagramElement, trans)
		crldiagramdomain.SetLinkTarget(ownerPointer, diagramConceptOwner, trans)
		ownerPointer.SetOwningConcept(diagram, trans)
	}
	return nil
}

func (dm *FyneDiagramManager) setToolbarSelection(sel ToolbarSelection) {
	if sel != dm.currentToolbarSelection {
		if dm.GetSelectedDiagram().ConnectionTransaction != nil {
			dm.cancelLinkTransaction()
		}
		dm.currentToolbarSelection = sel
		dm.toolButtons[sel].Importance = widget.HighImportance
		dm.toolButtons[sel].Refresh()
		for i := CursorSelected; i <= CreateRefinementOfConceptSelected; i++ {
			if i != sel {
				dm.toolButtons[i].Importance = widget.LowImportance
				dm.toolButtons[i].Refresh()
			}
		}
	}
}

func (dm *FyneDiagramManager) showAbstractConcept(elementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	diagramElement := uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showAbstractConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(trans)
	if diagram == nil {
		return errors.New("diagramManager.showAbstractConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelConcept(diagramElement, trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showAbstractConcept modelConcept not found for elementID " + elementID)
	}
	var modelRefinement core.Concept
	switch modelConcept.GetConceptType() {
	case core.Refinement:
		modelRefinement = modelConcept
	default:
		return errors.New("diagramManager.showAbstractConcept modelConcept is not a Refinement")
	}
	modelAbstractConcept := modelRefinement.GetAbstractConcept(trans)
	if modelAbstractConcept == nil {
		return errors.New("Abstract Concept is nil")
	}
	diagramAbstractConcept := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelAbstractConcept, trans)
	if diagramAbstractConcept == nil {
		diagramAbstractConcept, _ = crldiagramdomain.NewDiagramNode(uOfD, trans)
		crldiagramdomain.SetReferencedModelConcept(diagramAbstractConcept, modelAbstractConcept, trans)
		crldiagramdomain.SetDisplayLabel(diagramAbstractConcept, modelAbstractConcept.GetLabel(trans), trans)
		diagramElementX := crldiagramdomain.GetNodeX(diagramElement, trans)
		diagramElementY := crldiagramdomain.GetNodeY(diagramElement, trans)
		crldiagramdomain.SetNodeX(diagramAbstractConcept, diagramElementX, trans)
		crldiagramdomain.SetNodeY(diagramAbstractConcept, diagramElementY-100, trans)
		diagramAbstractConcept.SetOwningConcept(diagram, trans)
	}
	elementPointer := crldiagramdomain.GetElementPointer(diagram, diagramElement, trans)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramAbstractPointer(uOfD, trans)
		crldiagramdomain.SetReferencedModelConcept(elementPointer, modelConcept, trans)
		crldiagramdomain.SetLinkSource(elementPointer, diagramElement, trans)
		crldiagramdomain.SetLinkTarget(elementPointer, diagramAbstractConcept, trans)
		elementPointer.SetOwningConcept(diagram, trans)
	}
	return nil
}

func (dm *FyneDiagramManager) showReferencedConcept(elementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	return dm.showReferencedConceptImpl(uOfD, elementID, trans)
}

func (dm *FyneDiagramManager) showReferencedConceptImpl(uOfD *core.UniverseOfDiscourse, elementID string, trans *core.Transaction) error {
	diagramElement := uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showReferencedConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(trans)
	if diagram == nil {
		return errors.New("diagramManager.showReferencedConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelConcept(diagramElement, trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showReferencedConcept modelConcept not found for elementID " + elementID)
	}
	var modelReference core.Concept
	switch modelConcept.GetConceptType() {
	case core.Reference:
		modelReference = modelConcept
	default:
		return nil
	}
	modelReferencedConcept := modelReference.GetReferencedConcept(trans)
	if modelReferencedConcept == nil {
		return nil
	}
	var diagramReferencedConcept core.Concept
	switch modelReference.GetReferencedAttributeName(trans) {
	case core.NoAttribute, core.LiteralValue:
		diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelReferencedConcept, trans)
		if diagramReferencedConcept == nil {
			diagramReferencedConcept, _ = crldiagramdomain.NewDiagramNode(uOfD, trans)
			crldiagramdomain.SetReferencedModelConcept(diagramReferencedConcept, modelReferencedConcept, trans)
			crldiagramdomain.SetDisplayLabel(diagramReferencedConcept, modelReferencedConcept.GetLabel(trans), trans)
			diagramElementX := crldiagramdomain.GetNodeX(diagramElement, trans)
			diagramElementY := crldiagramdomain.GetNodeY(diagramElement, trans)
			diagramElementWidth := crldiagramdomain.GetNodeWidth(diagramElement, trans)
			xOffset := FyneGUISingleton.editor.GetUserPreferences().HorizontalLayoutSpacing
			crldiagramdomain.SetNodeX(diagramReferencedConcept, diagramElementX+diagramElementWidth+xOffset, trans)
			crldiagramdomain.SetNodeY(diagramReferencedConcept, diagramElementY, trans)
			diagramReferencedConcept.SetOwningConcept(diagram, trans)
		}
	case core.OwningConceptID:
		diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConceptOwnerPointer(diagram, modelReferencedConcept, trans)
		if diagramReferencedConcept == nil {
			return errors.New("No representation of the owner pointer currently exists in this diagram")
		}
	case core.ReferencedConceptID:
		switch modelReferencedConcept.GetConceptType() {
		case core.Reference:
			diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConceptElementPointer(diagram, modelReferencedConcept, trans)
			if diagramReferencedConcept == nil {
				return errors.New("No representation of the referenced concept pointer currently exists in this diagram")
			}
		}
	case core.AbstractConceptID:
		switch modelReferencedConcept.GetConceptType() {
		case core.Refinement:
			diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConceptAbstractPointer(diagram, modelReferencedConcept, trans)
		}
		if diagramReferencedConcept == nil {
			return errors.New("No representation of the abstract concept pointer currently exists in this diagram")
		}
	case core.RefinedConceptID:
		switch modelReferencedConcept.GetConceptType() {
		case core.Refinement:
			diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConceptRefinedPointer(diagram, modelReferencedConcept, trans)
			if diagramReferencedConcept == nil {
				return errors.New("No representation of the refined concept pointer currently exists in this diagram")
			}
		}
	}
	elementPointer := crldiagramdomain.GetElementPointer(diagram, diagramElement, trans)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramElementPointer(uOfD, trans)
		crldiagramdomain.SetReferencedModelConcept(elementPointer, modelConcept, trans)
		crldiagramdomain.SetLinkSource(elementPointer, diagramElement, trans)
		crldiagramdomain.SetLinkTarget(elementPointer, diagramReferencedConcept, trans)
		elementPointer.SetOwningConcept(diagram, trans)
	}
	return nil
}

func (dm *FyneDiagramManager) showReferencedConceptsRecursively(elementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	return dm.showReferencedConceptsRecursivelyImpl(uOfD, elementID, trans)
}

func (dm *FyneDiagramManager) showReferencedConceptsRecursivelyImpl(uOfD *core.UniverseOfDiscourse, elementID string, trans *core.Transaction) error {
	diagramElement := uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showReferencedConceptsRecursivelyImpl diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(trans)
	if diagram == nil {
		return errors.New("diagramManager.showOwnedConcepts diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelConcept(diagramElement, trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showReferencedConceptsRecursivelyImpl modelConcept not found for elementID " + elementID)
	}
	if modelConcept.GetConceptType() == core.Reference {
		err := dm.showReferencedConceptImpl(uOfD, elementID, trans)
		if err != nil {
			return errors.Wrap(err, "FyneDiagramManager.showReferencedConceptsRecursivelyImpl failed")
		}
	}
	it := modelConcept.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		child := uOfD.GetElement(id.(string))
		if child == nil {
			return errors.New("Child Concept is nil for id " + id.(string))
		}
		diagramChildConcept := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, child, trans)
		if diagramChildConcept != nil {
			err := dm.showReferencedConceptsRecursivelyImpl(uOfD, diagramChildConcept.GetConceptID(trans), trans)
			if err != nil {
				return errors.Wrap(err, "FyneDiagramManager.showReferencedConceptsRecursivelyImpl failed")
			}
		}
	}
	return nil
}

func (dm *FyneDiagramManager) showRefinedConcept(elementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	diagramElement := uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showRefinedConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(trans)
	if diagram == nil {
		return errors.New("diagramManager.showRefinedConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelConcept(diagramElement, trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showRefinedConcept modelConcept not found for elementID " + elementID)
	}
	var modelRefinement core.Concept
	switch modelConcept.GetConceptType() {
	case core.Refinement:
		modelRefinement = modelConcept
	default:
		return errors.New("diagramManager.showRefinedConcept modelConcept is not a Refinement")
	}
	modelRefinedConcept := modelRefinement.GetRefinedConcept(trans)
	if modelRefinedConcept == nil {
		return errors.New("Refined Concept is nil")
	}
	diagramRefinedConcept := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelRefinedConcept, trans)
	if diagramRefinedConcept == nil {
		diagramRefinedConcept, _ = crldiagramdomain.NewDiagramNode(uOfD, trans)
		crldiagramdomain.SetReferencedModelConcept(diagramRefinedConcept, modelRefinedConcept, trans)
		crldiagramdomain.SetDisplayLabel(diagramRefinedConcept, modelRefinedConcept.GetLabel(trans), trans)
		diagramElementX := crldiagramdomain.GetNodeX(diagramElement, trans)
		diagramElementY := crldiagramdomain.GetNodeY(diagramElement, trans)
		crldiagramdomain.SetNodeX(diagramRefinedConcept, diagramElementX, trans)
		crldiagramdomain.SetNodeY(diagramRefinedConcept, diagramElementY-100, trans)
		diagramRefinedConcept.SetOwningConcept(diagram, trans)
	}
	elementPointer := crldiagramdomain.GetElementPointer(diagram, diagramElement, trans)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramRefinedPointer(uOfD, trans)
		crldiagramdomain.SetReferencedModelConcept(elementPointer, modelConcept, trans)
		crldiagramdomain.SetLinkSource(elementPointer, diagramElement, trans)
		crldiagramdomain.SetLinkTarget(elementPointer, diagramRefinedConcept, trans)
		elementPointer.SetOwningConcept(diagram, trans)
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
	case RefinementLinkSelected, ReferenceLinkSelected, AbstractElementPointerSelected, OwnerPointerSelected, ReferencedElementPointerSelected, RefinedElementPointerSelected:
		trans, new := dm.fyneGUI.editor.GetTransaction()
		if new {
			defer dm.fyneGUI.editor.EndTransaction()
		}
		uOfD := trans.GetUniverseOfDiscourse()
		var crlLink core.Concept
		var fyneLink diagramwidget.DiagramLink
		switch dm.currentToolbarSelection {
		case ReferenceLinkSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramReferenceLink(uOfD, trans)
			crlModelReference, _ := uOfD.NewReference(trans)
			dm.connectionTransactionTransientConcepts.Add(crlModelReference.GetConceptID(trans))
			crlModelReference.SetLabel(FyneGUISingleton.editor.GetDefaultReferenceLabel(), trans)
			crldiagramdomain.SetReferencedModelConcept(crlLink, crlModelReference, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
			fyneLink.Hide()
		case RefinementLinkSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramRefinementLink(uOfD, trans)
			crlModelRefinement, _ := uOfD.NewRefinement(trans)
			dm.connectionTransactionTransientConcepts.Add(crlModelRefinement.GetConceptID(trans))
			crlModelRefinement.SetLabel(FyneGUISingleton.editor.GetDefaultRefinementLabel(), trans)
			crldiagramdomain.SetReferencedModelConcept(crlLink, crlModelRefinement, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
			fyneLink.Hide()
		case AbstractElementPointerSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramAbstractPointer(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case OwnerPointerSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramOwnerPointer(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case ReferencedElementPointerSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramElementPointer(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case RefinedElementPointerSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramRefinedPointer(uOfD, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		}
		crlDiagram := uOfD.GetElement(currentDiagram.ID)
		crlLink.SetOwningConcept(crlDiagram, trans)
		dm.connectionTransactionTransientConcepts.Add(crlLink.GetConceptID(trans))
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
					if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramNodeBGColorURI, trans) {
						bgColor := crldiagramdomain.GetBGColor(crlDiagramElement, trans)
						log.Printf("Background Color: %s", bgColor)
						goColor := getGoColor(bgColor)
						fyneDiagramElement.SetBackgroundColor(goColor)
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
				case core.IndicatedConceptChanged:
					reportingConcept := trans.GetUniverseOfDiscourse().GetElement(ownedConceptChangedNotification.GetReportingElementID())
					if reportingConcept != nil && reportingConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementModelReferenceURI, trans) {
						indicatedNotification := ownedConceptChangedNotification.GetUnderlyingChange()
						switch indicatedNotification.GetNatureOfChange() {
						case core.ReferencedConceptChanged:
							if indicatedNotification.GetAfterConceptState().ReferencedConceptID == "" {
								trans.GetUniverseOfDiscourse().DeleteElement(crlDiagramElement, trans)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func getCrlColor(goColor color.Color) string {
	switch typedColor := goColor.(type) {
	case color.NRGBA:
		r, g, b, a := typedColor.RGBA()
		goColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	switch typedColor := goColor.(type) {
	case color.RGBA:
		red := fmt.Sprintf("%02x", typedColor.R)
		green := fmt.Sprintf("%02x", typedColor.G)
		blue := fmt.Sprintf("%02x", typedColor.B)
		a := fmt.Sprintf("%02x", typedColor.A)
		crlColor := "x" + red + green + blue + a
		return crlColor
	}
	return ""
}

func getGoColor(lineColor string) color.Color {
	if lineColor == "" {
		return color.Transparent
	}
	redString := lineColor[1:3]
	red, _ := strconv.ParseUint(redString, 16, 8)
	greenString := lineColor[3:5]
	green, _ := strconv.ParseUint(greenString, 16, 8)
	blueString := lineColor[5:7]
	blue, _ := strconv.ParseUint(blueString, 16, 8)
	aString := "ff"
	if len(lineColor) == 9 {
		aString = lineColor[7:9]
	}
	a, _ := strconv.ParseUint(aString, 16, 8)
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
