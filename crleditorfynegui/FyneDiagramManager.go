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

func (dm *FyneDiagramManager) addElementToDiagram(element *crldiagramdomain.CrlDiagramElement, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramElement {
	if crldiagramdomain.IsDiagramNode(element.ToCore(), trans) {
		return dm.addNodeToDiagram((*crldiagramdomain.CrlDiagramNode)(element), trans, diagramWidget)
	} else if crldiagramdomain.IsDiagramLink(element.ToCore(), trans) {
		return dm.addLinkToDiagram((*crldiagramdomain.CrlDiagramLink)(element), trans, diagramWidget)
	}
	return nil
}

func (dm *FyneDiagramManager) addLinkToDiagram(link *crldiagramdomain.CrlDiagramLink, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) *FyneCrlDiagramLink {
	crlDiagramSource := link.GetLinkSource(trans)
	if crlDiagramSource == nil {
		// Register for changes so that when sufficient information is present we can add it to the diagram
		link.ToCore().Register(dm.diagramElementObserver)
		return nil
	}
	fyneSource := diagramWidget.GetDiagramElement(crlDiagramSource.ToCore().GetConceptID(trans))
	if fyneSource == nil {
		// the source is not in the diagram
		return nil
	}
	fyneSourcePad := fyneSource.GetDefaultConnectionPad()
	crlDiagramTarget := link.GetLinkTarget(trans)
	if crlDiagramTarget == nil {
		// Register for changes so that when sufficient information is present we can add it to the diagram
		link.ToCore().Register(dm.diagramElementObserver)
		return nil
	}
	targetConceptID := crlDiagramTarget.ToCore().GetConceptID(trans)
	fyneTarget := diagramWidget.GetDiagramElement(targetConceptID)
	if fyneTarget == nil {
		// the target is not in the diagram
		return nil
	}
	fyneTargetPad := fyneTarget.GetDefaultConnectionPad()
	diagramLink := NewFyneCrlDiagramLink(diagramWidget, link, trans)
	diagramLink.SetSourcePad(fyneSourcePad)
	diagramLink.SetTargetPad(fyneTargetPad)
	link.ToCore().Register(dm.diagramElementObserver)
	diagramWidget.Refresh()
	return diagramLink
}

func (dm *FyneDiagramManager) addNodeToDiagram(node *crldiagramdomain.CrlDiagramNode, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramNode {
	nodeID := node.ToCore().GetConceptID(trans)
	diagramNode := diagramWidget.GetDiagramNode(nodeID)
	if diagramNode == nil {

		diagramNode = NewFyneCrlDiagramNode(node, trans, diagramWidget)
		node.ToCore().Register(dm.diagramElementObserver)
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
	FyneGUISingleton.treeManager.tree.Refresh()
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
	var crlElement *core.Concept
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
	crlDiagram := crldiagramdomain.GetCrlDiagram(fyneDiagram.ID, trans)
	var el *core.Concept
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
		el, _ = crlmapsdomain.NewOneToOneMap(uOfD, trans)
		el.SetOwningConcept(crlDiagram.ToCore().GetOwningConcept(trans), trans)
	case CreateRefinementOfConceptSelected:
		selection := FyneGUISingleton.editor.GetCurrentSelection()
		if selection != nil {
			el, _ = uOfD.CreateRefinementOfConcept(selection, selection.GetLabel(trans), trans)
			el.SetOwningConcept(crlDiagram.ToCore().GetOwningConcept(trans), trans)
		}
	}

	if el != nil {
		elID := el.GetConceptID(trans)
		el.SetOwningConceptID(crlDiagram.ToCore().GetOwningConceptID(trans), trans)
		dm.fyneGUI.editor.SelectElement(el, trans)

		// Now the view
		x := event.Position.X
		y := event.Position.Y
		newNode, err := crldiagramdomain.NewDiagramNode(trans)
		if err != nil {
			log.Print(err)
			return
		}
		newNode.ToCore().Register(dm.diagramElementObserver)
		newNode.ToCrlDiagramElement().SetLineColor("x000000ff", trans)
		newNode.SetNodeX(float64(x), trans)
		newNode.SetNodeY(float64(y), trans)
		newNode.ToCore().SetLabel(el.GetLabel(trans), trans)
		newNode.ToCrlDiagramElement().SetReferencedModelConcept(el, trans)
		newNode.ToCrlDiagramElement().SetDisplayLabel(el.GetLabel(trans), trans)
		newNode.ToCrlDiagramElement().SetDiagram(crlDiagram, trans)
		dm.selectElementInDiagram(elID, fyneDiagram, trans)
		dm.ElementSelected(elID, trans)
	} else {
		dm.ElementSelected("", trans)
	}
	dm.setToolbarSelection(CursorSelected)
}

func (dm *FyneDiagramManager) displayDiagram(diagram *crldiagramdomain.CrlDiagram, trans *core.Transaction) error {
	diagramID := diagram.ToCore().GetConceptID(trans)
	tabItem := dm.diagramTabs[diagramID]
	if tabItem == nil {
		diagramWidget := diagramwidget.NewDiagramWidget(diagramID)
		diagramWidget.OnTappedCallback = dm.diagramTapped
		diagramWidget.MouseMovedCallback = dm.diagramMouseMoved
		scrollingContainer := container.NewScroll(diagramWidget)
		tabItem = &diagramTab{
			diagramID: diagramID,
			tab:       container.NewTabItem(diagram.ToCore().GetLabel(trans), scrollingContainer),
			diagram:   diagramWidget,
		}
		dm.diagramTabs[diagramID] = tabItem
		dm.tabArea.Append(tabItem.tab)
		diagram.ToCore().Register(dm.diagramObserver)
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
		crlDiagram := crldiagramdomain.GetCrlDiagram(diagramID, trans)
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
	crlLink := crldiagramdomain.GetCrlDiagramLink(fyneLink.GetDiagramElementID(), trans)
	crlPadOwner := crldiagramdomain.GetCrlDiagramElement(pad.GetPadOwner().GetDiagramElementID(), trans)
	if crlLink.IsReferenceLink(trans) {
		return true
	} else if crlLink.IsAbstractPointer(trans) {
		padOwnerModelElement := crlPadOwner.GetReferencedModelConcept(trans)
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
			if crlPadOwner.IsDiagramPointer(trans) {
				return false
			}
			return true
		}
	} else if crlLink.IsElementPointer(trans) {
		padOwnerModelElement := crlPadOwner.GetReferencedModelConcept(trans)
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
	} else if crlLink.IsOwnerPointer(trans) {
		switch linkEnd {
		case diagramwidget.SOURCE:
			return true
		case diagramwidget.TARGET:
			if crlPadOwner.IsDiagramPointer(trans) {
				return false
			}
			if crlPadOwner != crlLink.GetLinkSource(trans) {
				// an element cannot own itself
				return true
			}
		}
	} else if crlLink.IsRefinedPointer(trans) {
		padOwnerModelElement := crlPadOwner.GetReferencedModelConcept(trans)
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
			if crlPadOwner.IsDiagramPointer(trans) {
				return false
			}
			return true
		}
	} else if crlLink.IsRefinementLink(trans) {
		return !crlPadOwner.IsDiagramPointer(trans)
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
		crlLink := crldiagramdomain.GetCrlDiagramLink(link.GetDiagramElementID(), trans)
		if crlLink == nil {
			return errors.New("in FyneDiagramManager.linkConnectionChanged CrlLink not found")
		}
		crlNewPadOwner := crldiagramdomain.GetCrlDiagramElement(newPad.GetPadOwner().GetDiagramElementID(), trans)
		if crlNewPadOwner == nil {
			return errors.New("in FyneDiagramManager.linkConnectionChanged CrlLink not found")
		}
		switch end {
		case "source":
			crlLink.SetLinkSource(crlNewPadOwner, trans)
			switch typedLink.linkType {
			case ReferenceLinkSelected:
				linkModelElement := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				sourceModelElement := crlNewPadOwner.GetReferencedModelConcept(trans)
				linkModelElement.SetOwningConcept(sourceModelElement, trans)
				link.Show()
			case RefinementLinkSelected:
				linkModelElement := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				sourceModelElement := crlNewPadOwner.GetReferencedModelConcept(trans)
				switch linkModelElement.GetConceptType() {
				case core.Refinement:
					linkModelElement.SetOwningConcept(sourceModelElement, trans)
					linkModelElement.SetRefinedConcept(sourceModelElement, trans)
					link.Show()
				}
			case AbstractElementPointerSelected:
				currentModelRefinement := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				newModelRefinement := crlNewPadOwner.GetReferencedModelConcept(trans)
				if currentModelRefinement != newModelRefinement {
					if currentModelRefinement != nil {
						currentModelRefinement.SetAbstractConcept(nil, trans)
					}
					crlLinkTarget := crlLink.GetLinkTarget(trans)
					targetModelElement := crlLinkTarget.GetReferencedModelConcept(trans)
					switch newModelRefinement.GetConceptType() {
					case core.Refinement:
						newModelRefinement.SetAbstractConcept(targetModelElement, trans)
						crlLink.ToCrlDiagramElement().SetReferencedModelConcept(newModelRefinement, trans)
						typedLink.modelElement = newModelRefinement
					}
				}
			case OwnerPointerSelected:
				currentLinkModelConcept := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				newLinkModelConcept := crlNewPadOwner.GetReferencedModelConcept(trans)
				if currentLinkModelConcept != newLinkModelConcept {
					if currentLinkModelConcept != nil {
						currentLinkModelConcept.SetOwningConcept(nil, trans)
					}
					crlLinkTarget := crlLink.GetLinkTarget(trans)
					targetModelElement := crlLinkTarget.GetReferencedModelConcept(trans)
					newLinkModelConcept.SetOwningConcept(targetModelElement, trans)
					crlLink.ToCrlDiagramElement().SetReferencedModelConcept(newLinkModelConcept, trans)
					typedLink.modelElement = newLinkModelConcept
				}
			case ReferencedElementPointerSelected:
				currentModelReference := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				newModelReference := crlNewPadOwner.GetReferencedModelConcept(trans)
				if currentModelReference != newModelReference {
					crlLink.ToCrlDiagramElement().SetReferencedModelConcept(newModelReference, trans)
					attributeName := core.NoAttribute
					if currentModelReference != nil {
						switch currentModelReference.GetConceptType() {
						case core.Reference:
							attributeName = currentModelReference.GetReferencedAttributeName(trans)
							currentModelReference.SetReferencedConcept(nil, core.NoAttribute, trans)
						}
					}
					crlLinkTarget := crlLink.GetLinkTarget(trans)
					targetModelElement := crlLinkTarget.GetReferencedModelConcept(trans)
					switch newModelReference.GetConceptType() {
					case core.Reference:
						newModelReference.SetReferencedConcept(targetModelElement, attributeName, trans)
						crlLink.ToCrlDiagramElement().SetReferencedModelConcept(newModelReference, trans)
						typedLink.modelElement = newModelReference
					}
				}
			case RefinedElementPointerSelected:
				currentModelRefinement := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				newModelRefinement := crlNewPadOwner.GetReferencedModelConcept(trans)
				if currentModelRefinement != newModelRefinement {
					if currentModelRefinement != nil {
						switch currentModelRefinement.GetConceptType() {
						case core.Refinement:
							currentModelRefinement.SetRefinedConcept(nil, trans)
						}
					}
					crlLinkTarget := crlLink.GetLinkTarget(trans)
					targetModelElement := crlLinkTarget.GetReferencedModelConcept(trans)
					switch newModelRefinement.GetConceptType() {
					case core.Refinement:
						newModelRefinement.SetRefinedConcept(targetModelElement, trans)
						crlLink.ToCrlDiagramElement().SetReferencedModelConcept(newModelRefinement, trans)
						typedLink.modelElement = newModelRefinement
					}
				}
			}
		case "target":
			crlLink.SetLinkTarget(crlNewPadOwner, trans)
			switch typedLink.linkType {
			case ReferenceLinkSelected:
				linkModelElement := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				targetModelElement := crlNewPadOwner.GetReferencedModelConcept(trans)
				newPadOwner := newPad.GetPadOwner()
				attributeName := getAttributeNameBasedOnTargetType(newPadOwner)
				switch linkModelElement.GetConceptType() {
				case core.Reference:
					linkModelElement.SetReferencedConcept(targetModelElement, attributeName, trans)
				}
			case RefinementLinkSelected:
				linkModelElement := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				targetModelElement := crlNewPadOwner.GetReferencedModelConcept(trans)
				switch linkModelElement.GetConceptType() {
				case core.Refinement:
					linkModelElement.SetAbstractConcept(targetModelElement, trans)
				}
			case AbstractElementPointerSelected:
				crlModelRefinement := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				currentAbstractElement := crlModelRefinement.GetAbstractConcept(trans)
				newAbstractElement := crlNewPadOwner.GetReferencedModelConcept(trans)
				if currentAbstractElement != newAbstractElement {
					switch crlModelRefinement.GetConceptType() {
					case core.Refinement:
						crlModelRefinement.SetAbstractConcept(newAbstractElement, trans)
					}
				}
			case OwnerPointerSelected:
				crlLinkParent := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				targetModelElement := crlNewPadOwner.GetReferencedModelConcept(trans)
				if crlLinkParent != nil && crlLinkParent.GetOwningConcept(trans) != targetModelElement {
					crlLinkParent.SetOwningConcept(targetModelElement, trans)
				}
			case ReferencedElementPointerSelected:
				crlModelReference := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				currentReferencedElement := crlModelReference.GetReferencedConcept(trans)
				newReferencedElement := crlNewPadOwner.GetReferencedModelConcept(trans)
				if currentReferencedElement != newReferencedElement {
					attributeName := getAttributeNameBasedOnTargetType(newPad.GetPadOwner())
					switch crlModelReference.GetConceptType() {
					case core.Reference:
						crlModelReference.SetReferencedConcept(newReferencedElement, attributeName, trans)
					}
				}
			case RefinedElementPointerSelected:
				crlModelRefinement := crlLink.ToCrlDiagramElement().GetReferencedModelConcept(trans)
				currentRefinedElement := crlModelRefinement.GetRefinedConcept(trans)
				newRefinedElement := crlNewPadOwner.GetReferencedModelConcept(trans)
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
func (dm *FyneDiagramManager) populateDiagram(diagram *crldiagramdomain.CrlDiagram, trans *core.Transaction) error {
	diagramWidget := dm.getDiagramWidget(diagram.ToCore().GetConceptID(trans))
	nodes := diagram.ToCore().GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, trans)
	for _, node := range nodes {
		dm.addNodeToDiagram((*crldiagramdomain.CrlDiagramNode)(node), trans, diagramWidget)
	}
	links := diagram.ToCore().GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramLinkURI, trans)
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
				link := crldiagramdomain.GetCrlDiagramLink(linkID, trans)
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
	diagram.ClearSelectionNoCallback()
	if elementID == "" {
		return nil
	}
	foundDiagramElementID := ""
	for _, fyneDiagramElement := range diagram.GetDiagramElements() {
		id := fyneDiagramElement.GetDiagramElementID()
		crlDiagramElement := crldiagramdomain.GetCrlDiagramElement(id, trans)
		if crlDiagramElement != nil {
			crlModelElement := crlDiagramElement.GetReferencedModelConcept(trans)
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
	return dm.showOwnedConceptsImpl(elementID, recursive, skipRefinements, trans)
}

func (dm *FyneDiagramManager) showOwnedConceptsImpl(elementID string, recursive bool, skipRefinements bool, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	ownerDiagramElement := crldiagramdomain.GetCrlDiagramElement(elementID, trans)
	if ownerDiagramElement == nil {
		return errors.New("diagramManager.showOwnedConcepts diagramElement not found for elementID " + elementID)
	}
	diagram := ownerDiagramElement.GetDiagram(trans)
	if diagram == nil {
		return errors.New("diagramManager.showOwnedConcepts diagram not found for elementID " + elementID)
	}
	modelConcept := ownerDiagramElement.GetReferencedModelConcept(trans)
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
		diagramChildConcept := diagram.GetFirstElementRepresentingConcept(child, trans)
		if diagramChildConcept == nil {
			newChildNode, _ := crldiagramdomain.NewDiagramNode(trans)
			diagramChildConcept = newChildNode.ToCrlDiagramElement()
			diagramChildConcept.SetReferencedModelConcept(child, trans)
			diagramChildConcept.SetDisplayLabel(child.GetLabel(trans), trans)
			// TODO Address the case in which the diagram element is a link!
			if ownerDiagramElement.IsNode(trans) {
				ownerNode := ownerDiagramElement.ToNode(trans)
				diagramElementX := ownerNode.GetNodeX(trans)
				diagramElementY := ownerNode.GetNodeY(trans)
				diagramElementHeight := ownerNode.GetNodeHeight(trans)
				newChildNode.SetNodeX(diagramElementX+xOffset, trans)
				newChildNode.SetNodeY(diagramElementY+diagramElementHeight+yPreferencesOffset, trans)
				diagramChildConcept.SetDiagram(diagram, trans)
				xOffset = xOffset + xPreferencesOffset + newChildNode.GetNodeWidth(trans)
			}
		}
		ownerPointer := diagram.GetOwnerPointer(ownerDiagramElement, trans)
		if ownerPointer == nil {
			ownerPointer, _ = crldiagramdomain.NewDiagramOwnerPointer(trans)
			ownerPointer.ToCrlDiagramElement().SetReferencedModelConcept(child, trans)
			ownerPointer.SetLinkSource(diagramChildConcept, trans)
			ownerPointer.SetLinkTarget(ownerDiagramElement, trans)
			ownerPointer.ToCrlDiagramElement().SetDiagram(diagram, trans)
		}
		if recursive {
			dm.showOwnedConceptsImpl(diagramChildConcept.ToCore().GetConceptID(trans), recursive, skipRefinements, trans)
		}
	}
	return nil
}

func (dm *FyneDiagramManager) showOwner(diagramElementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	diagramElement := crldiagramdomain.GetCrlDiagramElement(diagramElementID, trans)
	if diagramElement == nil {
		return errors.New("diagramManager.showOwner diagramElement not found for elementID " + diagramElementID)
	}
	diagram := diagramElement.GetDiagram(trans)
	if diagram == nil {
		return errors.New("diagramManager.showOwner diagram not found for elementID " + diagramElementID)
	}
	modelConcept := diagramElement.GetReferencedModelConcept(trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showOwner modelConcept not found for elementID " + diagramElementID)
	}
	modelConceptOwner := modelConcept.GetOwningConcept(trans)
	if modelConceptOwner == nil {
		return errors.New("Owner is nil")
	}
	diagramConceptOwner := diagram.GetFirstElementRepresentingConcept(modelConceptOwner, trans)
	if diagramConceptOwner == nil {
		diagramConceptOwnerNode, _ := crldiagramdomain.NewDiagramNode(trans)
		diagramConceptOwner := diagramConceptOwnerNode.ToCrlDiagramElement()
		diagramConceptOwner.SetReferencedModelConcept(modelConceptOwner, trans)
		diagramConceptOwner.SetDisplayLabel(modelConceptOwner.GetLabel(trans), trans)
		// ToDo address case in which diagramElement is a link
		if diagramElement.IsNode(trans) {
			diagramElementNode := diagramElement.ToNode(trans)
			diagramElementX := diagramElementNode.GetNodeX(trans)
			diagramElementY := diagramElementNode.GetNodeY(trans)
			diagramConceptOwnerNode.SetNodeX(diagramElementX, trans)
			diagramConceptOwnerNode.SetNodeY(diagramElementY-100, trans)
			diagramConceptOwner.SetDiagram(diagram, trans)
		}
	}
	ownerPointer := diagram.GetOwnerPointer(diagramElement, trans)
	if ownerPointer == nil {
		ownerPointer, _ = crldiagramdomain.NewDiagramOwnerPointer(trans)
		ownerPointer.ToCrlDiagramElement().SetReferencedModelConcept(modelConcept, trans)
		ownerPointer.SetLinkSource(diagramElement, trans)
		ownerPointer.SetLinkTarget(diagramConceptOwner, trans)
		ownerPointer.ToCrlDiagramElement().SetDiagram(diagram, trans)
	}
	return nil
}

func (dm *FyneDiagramManager) setToolbarSelection(sel ToolbarSelection) {
	if dm.GetSelectedDiagram() == nil {
		return
	}
	if sel != dm.currentToolbarSelection {
		if dm.GetSelectedDiagram().ConnectionTransaction != nil {
			dm.cancelLinkTransaction()
		}
		dm.currentToolbarSelection = sel
		dm.toolButtons[sel].Importance = widget.HighImportance
		dm.toolButtons[sel].Refresh()
		for i := CursorSelected; i <= CreateRefinementOfConceptSelected; i++ {
			if i != sel {
				dm.toolButtons[i].Importance = widget.MediumImportance
				dm.toolButtons[i].Refresh()
			}
		}
	}
}

func (dm *FyneDiagramManager) showAbstractConcept(diagramElementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	diagramElement := crldiagramdomain.GetCrlDiagramElement(diagramElementID, trans)
	if diagramElement == nil {
		return errors.New("diagramManager.showAbstractConcept diagramElement not found for elementID " + diagramElementID)
	}
	diagram := diagramElement.GetDiagram(trans)
	if diagram == nil {
		return errors.New("diagramManager.showAbstractConcept diagram not found for elementID " + diagramElementID)
	}
	modelConcept := diagramElement.GetReferencedModelConcept(trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showAbstractConcept modelConcept not found for elementID " + diagramElementID)
	}
	var modelRefinement *core.Concept
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
	diagramAbstractConcept := diagram.GetFirstElementRepresentingConcept(modelAbstractConcept, trans)
	if diagramAbstractConcept == nil {
		diagramAbstractConceptNode, _ := crldiagramdomain.NewDiagramNode(trans)
		diagramAbstractConcept = diagramAbstractConceptNode.ToCrlDiagramElement()
		diagramAbstractConcept.SetReferencedModelConcept(modelAbstractConcept, trans)
		diagramAbstractConcept.SetDisplayLabel(modelAbstractConcept.GetLabel(trans), trans)
		// TODO address case in which diagram element is a link
		if diagramElement.IsNode(trans) {
			diagramElementNode := diagramElement.ToNode(trans)
			diagramElementX := diagramElementNode.GetNodeX(trans)
			diagramElementY := diagramElementNode.GetNodeY(trans)
			diagramAbstractConceptNode.SetNodeX(diagramElementX, trans)
			diagramAbstractConceptNode.SetNodeY(diagramElementY-100, trans)
		}
		diagramAbstractConcept.SetDiagram(diagram, trans)

	}
	elementPointer := diagram.GetElementPointer(diagramElement, trans)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramAbstractPointer(trans)
		elementPointer.ToCrlDiagramElement().SetReferencedModelConcept(modelConcept, trans)
		elementPointer.SetLinkSource(diagramElement, trans)
		elementPointer.SetLinkTarget(diagramAbstractConcept, trans)
		elementPointer.ToCrlDiagramElement().SetDiagram(diagram, trans)
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
	return dm.showReferencedConceptImpl(elementID, trans)
}

func (dm *FyneDiagramManager) showReferencedConceptImpl(diagramElementID string, trans *core.Transaction) error {
	diagramElement := crldiagramdomain.GetCrlDiagramElement(diagramElementID, trans)
	if diagramElement == nil {
		return errors.New("diagramManager.showReferencedConcept diagramElement not found for elementID " + diagramElementID)
	}
	diagram := diagramElement.GetDiagram(trans)
	if diagram == nil {
		return errors.New("diagramManager.showReferencedConcept diagram not found for elementID " + diagramElementID)
	}
	modelConcept := diagramElement.GetReferencedModelConcept(trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showReferencedConcept modelConcept not found for elementID " + diagramElementID)
	}
	var modelReference *core.Concept
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
	var diagramReferencedConcept *crldiagramdomain.CrlDiagramElement
	switch modelReference.GetReferencedAttributeName(trans) {
	case core.NoAttribute, core.LiteralValue:
		diagramReferencedConcept = diagram.GetFirstElementRepresentingConcept(modelReferencedConcept, trans)
		if diagramReferencedConcept == nil {
			diagramReferencedConceptNode, _ := crldiagramdomain.NewDiagramNode(trans)
			diagramReferencedConcept = diagramReferencedConceptNode.ToCrlDiagramElement()
			diagramReferencedConcept.SetReferencedModelConcept(modelReferencedConcept, trans)
			diagramReferencedConcept.SetDisplayLabel(modelReferencedConcept.GetLabel(trans), trans)
			// TODO Address case in which diagramElement is a link
			if diagramElement.IsNode(trans) {
				diagramElementNode := diagramElement.ToNode(trans)
				diagramElementX := diagramElementNode.GetNodeX(trans)
				diagramElementY := diagramElementNode.GetNodeY(trans)
				diagramElementWidth := diagramElementNode.GetNodeWidth(trans)
				xOffset := FyneGUISingleton.editor.GetUserPreferences().HorizontalLayoutSpacing
				diagramReferencedConceptNode.SetNodeX(diagramElementX+diagramElementWidth+xOffset, trans)
				diagramReferencedConceptNode.SetNodeY(diagramElementY, trans)
			}
			diagramReferencedConcept.SetDiagram(diagram, trans)
		}
	case core.OwningConceptID:
		diagramReferencedConcept = diagram.GetFirstElementRepresentingConceptOwnerPointer(modelReferencedConcept, trans).ToCrlDiagramElement()
		if diagramReferencedConcept == nil {
			return errors.New("No representation of the owner pointer currently exists in this diagram")
		}
	case core.ReferencedConceptID:
		switch modelReferencedConcept.GetConceptType() {
		case core.Reference:
			diagramReferencedConcept = diagram.GetFirstElementRepresentingConceptElementPointer(modelReferencedConcept, trans).ToCrlDiagramElement()
			if diagramReferencedConcept == nil {
				return errors.New("No representation of the referenced concept pointer currently exists in this diagram")
			}
		}
	case core.AbstractConceptID:
		switch modelReferencedConcept.GetConceptType() {
		case core.Refinement:
			diagramReferencedConcept = diagram.GetFirstElementRepresentingConceptAbstractPointer(modelReferencedConcept, trans).ToCrlDiagramElement()
		}
		if diagramReferencedConcept == nil {
			return errors.New("No representation of the abstract concept pointer currently exists in this diagram")
		}
	case core.RefinedConceptID:
		switch modelReferencedConcept.GetConceptType() {
		case core.Refinement:
			diagramReferencedConcept = diagram.GetFirstElementRepresentingConceptRefinedPointer(modelReferencedConcept, trans).ToCrlDiagramElement()
			if diagramReferencedConcept == nil {
				return errors.New("No representation of the refined concept pointer currently exists in this diagram")
			}
		}
	}
	elementPointer := diagram.GetElementPointer(diagramElement, trans)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramElementPointer(trans)
		elementPointer.ToCrlDiagramElement().SetReferencedModelConcept(modelConcept, trans)
		elementPointer.SetLinkSource(diagramElement, trans)
		elementPointer.SetLinkTarget(diagramReferencedConcept, trans)
		elementPointer.ToCrlDiagramElement().SetDiagram(diagram, trans)
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
	return dm.showReferencedConceptsRecursivelyImpl(elementID, trans)
}

func (dm *FyneDiagramManager) showReferencedConceptsRecursivelyImpl(diagramElementID string, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	diagramElement := crldiagramdomain.GetCrlDiagramElement(diagramElementID, trans)
	if diagramElement == nil {
		return errors.New("diagramManager.showReferencedConceptsRecursivelyImpl diagramElement not found for elementID " + diagramElementID)
	}
	diagram := diagramElement.GetDiagram(trans)
	if diagram == nil {
		return errors.New("diagramManager.showOwnedConcepts diagram not found for elementID " + diagramElementID)
	}
	modelConcept := diagramElement.GetReferencedModelConcept(trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showReferencedConceptsRecursivelyImpl modelConcept not found for elementID " + diagramElementID)
	}
	if modelConcept.GetConceptType() == core.Reference {
		err := dm.showReferencedConceptImpl(diagramElementID, trans)
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
		diagramChildConcept := diagram.GetFirstElementRepresentingConcept(child, trans)
		if diagramChildConcept != nil {
			err := dm.showReferencedConceptsRecursivelyImpl(diagramChildConcept.ToCore().GetConceptID(trans), trans)
			if err != nil {
				return errors.Wrap(err, "FyneDiagramManager.showReferencedConceptsRecursivelyImpl failed")
			}
		}
	}
	return nil
}

func (dm *FyneDiagramManager) showRefinedConcept(diagramElementID string) error {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	diagramElement := crldiagramdomain.GetCrlDiagramElement(diagramElementID, trans)
	if diagramElement == nil {
		return errors.New("diagramManager.showRefinedConcept diagramElement not found for elementID " + diagramElementID)
	}
	diagram := diagramElement.GetDiagram(trans)
	if diagram == nil {
		return errors.New("diagramManager.showRefinedConcept diagram not found for elementID " + diagramElementID)
	}
	modelConcept := diagramElement.GetReferencedModelConcept(trans)
	if modelConcept == nil {
		return errors.New("diagramManager.showRefinedConcept modelConcept not found for elementID " + diagramElementID)
	}
	var modelRefinement *core.Concept
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
	diagramRefinedConcept := diagram.GetFirstElementRepresentingConcept(modelRefinedConcept, trans)
	if diagramRefinedConcept == nil {
		diagramRefeinedConceptNode, _ := crldiagramdomain.NewDiagramNode(trans)
		diagramRefinedConcept = diagramRefeinedConceptNode.ToCrlDiagramElement()
		diagramRefinedConcept.SetReferencedModelConcept(modelRefinedConcept, trans)
		diagramRefinedConcept.SetDisplayLabel(modelRefinedConcept.GetLabel(trans), trans)
		// TODO address the case in which the diagramElement is a link
		if diagramElement.IsNode(trans) {
			diagramElementNode := diagramElement.ToNode(trans)
			diagramElementX := diagramElementNode.GetNodeX(trans)
			diagramElementY := diagramElementNode.GetNodeY(trans)
			diagramRefeinedConceptNode.SetNodeX(diagramElementX, trans)
			diagramRefeinedConceptNode.SetNodeY(diagramElementY-100, trans)
		}
		diagramRefinedConcept.SetDiagram(diagram, trans)
	}
	elementPointer := diagram.GetElementPointer(diagramElement, trans)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramRefinedPointer(trans)
		elementPointer.ToCrlDiagramElement().SetReferencedModelConcept(modelConcept, trans)
		elementPointer.SetLinkSource(diagramElement, trans)
		elementPointer.SetLinkTarget(diagramRefinedConcept, trans)
		elementPointer.ToCrlDiagramElement().SetDiagram(diagram, trans)
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
		var crlLink *crldiagramdomain.CrlDiagramLink
		var fyneLink diagramwidget.DiagramLink
		switch dm.currentToolbarSelection {
		case ReferenceLinkSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramReferenceLink(trans)
			crlModelReference, _ := uOfD.NewReference(trans)
			dm.connectionTransactionTransientConcepts.Add(crlModelReference.GetConceptID(trans))
			crlModelReference.SetLabel(FyneGUISingleton.editor.GetDefaultReferenceLabel(), trans)
			crlLink.ToCrlDiagramElement().SetReferencedModelConcept(crlModelReference, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
			fyneLink.Hide()
		case RefinementLinkSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramRefinementLink(trans)
			crlModelRefinement, _ := uOfD.NewRefinement(trans)
			dm.connectionTransactionTransientConcepts.Add(crlModelRefinement.GetConceptID(trans))
			crlModelRefinement.SetLabel(FyneGUISingleton.editor.GetDefaultRefinementLabel(), trans)
			crlLink.ToCrlDiagramElement().SetReferencedModelConcept(crlModelRefinement, trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
			fyneLink.Hide()
		case AbstractElementPointerSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramAbstractPointer(trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case OwnerPointerSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramOwnerPointer(trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case ReferencedElementPointerSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramElementPointer(trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		case RefinedElementPointerSelected:
			uOfD.MarkUndoPoint()
			crlLink, _ = crldiagramdomain.NewDiagramRefinedPointer(trans)
			fyneLink = NewFyneCrlDiagramLink(currentDiagram, crlLink, trans)
		}
		crlDiagram := crldiagramdomain.GetCrlDiagram(currentDiagram.ID, trans)
		crlLink.ToCrlDiagramElement().SetDiagram(crlDiagram, trans)
		dm.connectionTransactionTransientConcepts.Add(crlLink.ToCore().GetConceptID(trans))
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
				element := crldiagramdomain.GetCrlDiagramElement(afterState.ConceptID, trans)
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
	crlDiagramElement := crldiagramdomain.GetCrlDiagramElement(elementID, trans)
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
						x := float32(crlDiagramElement.ToNode(trans).GetNodeX(trans))
						fynePosition := fyneDiagramElement.Position()
						if x != fynePosition.X {
							fyneDiagramElement.Move(fyne.NewPos(x, fynePosition.Y))

						}
						return nil
					}
					if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramNodeYURI, trans) {
						y := float32(crlDiagramElement.ToNode(trans).GetNodeY(trans))
						fynePosition := fyneDiagramElement.Position()
						if y != fynePosition.Y {
							fyneDiagramElement.Move(fyne.NewPos(fynePosition.X, y))

						}
						return nil
					}
					if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementLineColorURI, trans) {
						lineColor := crlDiagramElement.GetLineColor(trans)
						log.Printf("Line Color: %s", lineColor)
						goColor := getGoColor(lineColor)
						fyneDiagramElement.SetForegroundColor(goColor)
					}
					if changedConcept.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementBGColorURI, trans) {
						bgColor := crlDiagramElement.GetBGColor(trans)
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
						lineColor := crlDiagramElement.GetLineColor(trans)
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
								trans.GetUniverseOfDiscourse().DeleteElement(crlDiagramElement.ToCore(), trans)
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
