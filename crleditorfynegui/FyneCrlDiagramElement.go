package crleditorfynegui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/diagramwidget"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
)

// FyneCrlDiagramElement serves as a mapping between the Fyne DiagramElement and the core DiagramElement
type FyneCrlDiagramElement interface {
	GetDiagramElement() *crldiagramdomain.CrlDiagramElement
	GetDiagramElementID() string
	GetModelElement() *core.Concept
	GetModelElementID() string
	GetFyneProperties() diagramwidget.DiagramElementProperties
	SetFyneProperties(diagramwidget.DiagramElementProperties)
	Refresh()
}

var _ fyne.Widget = (*FyneCrlDiagramNode)(nil)
var _ diagramwidget.DiagramElement = (*FyneCrlDiagramNode)(nil)
var _ diagramwidget.DiagramNode = (*FyneCrlDiagramNode)(nil)
var _ fyne.Tappable = (*FyneCrlDiagramNode)(nil)

// FyneCrlDiagramNode is an extension to diagramwidget.DiagramNode that serves as a binding
// between the diagramwidget nodes and the crldiagramdomain diagram noddes
type FyneCrlDiagramNode struct {
	diagramwidget.BaseDiagramNode
	crlDiagramNode  *crldiagramdomain.CrlDiagramNode
	modelElement    *core.Concept
	entryWidget     *widget.Entry
	abstractionText *canvas.Text
	labelBinding    binding.String
	// abstractionTextBinding binding.String
}

// NewFyneCrlDiagramNode creates a fyne node that corresponds to the supplied crldiagram node
func NewFyneCrlDiagramNode(crlNode *crldiagramdomain.CrlDiagramNode, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramNode {
	newNode := &FyneCrlDiagramNode{}
	nodeID := crlNode.ToCore().GetConceptID(trans)
	newNode.crlDiagramNode = crlNode
	newNode.modelElement = crlNode.ToCrlDiagramElement().GetReferencedModelConcept(trans)
	nodeIcon := widget.NewIcon(getIconResource(newNode.modelElement, trans))

	abstractionString := crlNode.ToCrlDiagramElement().GetAbstractionDisplayLabel(trans)
	newNode.abstractionText = canvas.NewText(abstractionString, color.Black)
	newNode.abstractionText.TextStyle = fyne.TextStyle{Bold: false, Italic: true, Monospace: false, Symbol: false, TabWidth: 4}

	hBox := container.NewHBox(nodeIcon, newNode.abstractionText)
	nodeLabel := crlNode.ToCrlDiagramElement().GetDisplayLabel(trans)
	newNode.labelBinding = binding.NewString()
	newNode.labelBinding.Set(nodeLabel)
	newNode.labelBinding.AddListener(binding.NewDataListener(func() { newNode.labelChanged() }))
	newNode.entryWidget = widget.NewEntryWithData(newNode.labelBinding)
	newNode.entryWidget.Wrapping = fyne.TextWrapOff
	newNode.entryWidget.Scroll = container.ScrollNone
	newNode.entryWidget.Validator = nil
	newNode.entryWidget.Refresh()

	newNode.MovedCallback = func() {
		newNode.nodeMoved()
	}

	nodeContainer := container.NewVBox(hBox, newNode.entryWidget)

	diagramwidget.InitializeBaseDiagramNode(newNode, diagramWidget, nodeContainer, nodeID)
	// Size isn't available until after initialization
	newNode.abstractionText.TextSize = newNode.GetFyneProperties().CaptionTextSize
	x := crlNode.GetNodeX(trans)
	y := crlNode.GetNodeY(trans)
	fynePosition := fyne.NewPos(float32(x), float32(y))
	newNode.Move(fynePosition)
	fgColor := crlNode.ToCrlDiagramElement().GetLineColor(trans)
	bgColor := crlNode.ToCrlDiagramElement().GetBGColor(trans)
	newNode.SetForegroundColor(getGoColor(fgColor))
	newNode.SetBackgroundColor(getGoColor(bgColor))
	newNode.Refresh()
	return newNode
}

// GetDiagramElement returns the crl diagram element associated with the link
func (fcdn *FyneCrlDiagramNode) GetDiagramElement() *crldiagramdomain.CrlDiagramElement {
	return fcdn.crlDiagramNode.ToCrlDiagramElement()
}

// GetDiagramElementID returns the ID of the crl diagram eleent associagted with the link
func (fcdn *FyneCrlDiagramNode) GetDiagramElementID() string {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	return fcdn.crlDiagramNode.ToCore().GetConceptID(trans)
}

// GetModelElement returns the crl model element represented by the link
func (fcdn *FyneCrlDiagramNode) GetModelElement() *core.Concept {
	return fcdn.modelElement
}

// GetModelElementID returns the ID of the crl model eleent represented by the link
func (fcdn *FyneCrlDiagramNode) GetModelElementID() string {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	return fcdn.modelElement.GetConceptID(trans)
}

// GetFyneProperties returns the fyne DiagramElementProperties of the diagram link
func (fcdn *FyneCrlDiagramNode) GetFyneProperties() diagramwidget.DiagramElementProperties {
	return fcdn.GetProperties()
}

func (fcdn *FyneCrlDiagramNode) labelChanged() {
	newValue, _ := fcdn.labelBinding.Get()
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	if fcdn.modelElement != nil {
		fcdn.modelElement.SetLabel(newValue, trans)
	}
}

// MouseDown shows the secondary popup for right mouse
func (fcdn *FyneCrlDiagramNode) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonSecondary {
		ShowSecondaryPopup(fcdn, event)
	}
}

func setCrlDiagramElementProperties(diagramElementID string, properties diagramwidget.DiagramElementProperties) {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.MarkUndoPoint()
	crlDiagramElement := crldiagramdomain.GetCrlDiagramElement(diagramElementID, trans)
	if crlDiagramElement == nil {
		return
	}
	crlFGColor := getCrlColor(properties.ForegroundColor)
	crlBGColor := getCrlColor(properties.BackgroundColor)
	crlDiagramElement.SetLineColor(crlFGColor, trans)
	crlDiagramElement.SetBGColor(crlBGColor, trans)
}

// ShowSecondaryPopup actually displays the secondary popup - it is used for both Nodes and Links
func ShowSecondaryPopup(fcde FyneCrlDiagramElement, event *desktop.MouseEvent) {
	items := []*fyne.MenuItem{}
	showModelConceptItem := fyne.NewMenuItem("Show Concept in Navigator", func() {
		FyneGUISingleton.treeManager.ShowElementInTree(fcde.GetModelElement())
	})
	items = append(items, showModelConceptItem)

	showDiagramElementItem := fyne.NewMenuItem("Show Diagram Element in Navigator", func() {
		FyneGUISingleton.treeManager.ShowElementInTree(fcde.GetDiagramElement().ToCore())
	})
	items = append(items, showDiagramElementItem)

	showOwnerItem := fyne.NewMenuItem("Show Owner", func() {
		FyneGUISingleton.diagramManager.showOwner(fcde.GetDiagramElementID())
	})
	items = append(items, showOwnerItem)

	showOwnedConceptsItem := fyne.NewMenuItem("Show Owned Conecpts", func() {
		FyneGUISingleton.diagramManager.showOwnedConcepts(fcde.GetDiagramElementID(), false, false)
	})
	items = append(items, showOwnedConceptsItem)

	showOwnedConceptsSkipRefimementsItem := fyne.NewMenuItem("Show Owned Conecpts Skip Refinements", func() {
		FyneGUISingleton.diagramManager.showOwnedConcepts(fcde.GetDiagramElementID(), false, true)
	})
	items = append(items, showOwnedConceptsSkipRefimementsItem)

	showOwnedConceptsRecursivelyItem := fyne.NewMenuItem("Show Owned Conecpts Recursively", func() {
		FyneGUISingleton.diagramManager.showOwnedConcepts(fcde.GetDiagramElementID(), true, false)
	})
	items = append(items, showOwnedConceptsRecursivelyItem)

	showOwnedConceptsRecursivelySkipRefimementsItem := fyne.NewMenuItem("Show Owned Conecpts Recursively Skip Refinements", func() {
		FyneGUISingleton.diagramManager.showOwnedConcepts(fcde.GetDiagramElementID(), true, true)
	})
	items = append(items, showOwnedConceptsRecursivelySkipRefimementsItem)

	showReferencedConcepsRecursivelyItem := fyne.NewMenuItem("Show Referenced Concepts Recursively", func() {
		FyneGUISingleton.diagramManager.showReferencedConceptsRecursively(fcde.GetDiagramElementID())
	})
	items = append(items, showReferencedConcepsRecursivelyItem)

	switch fcde.GetModelElement().GetConceptType() {
	case core.Reference:
		showReferencedConceptItem := fyne.NewMenuItem("Show Referenced Concept", func() {
			FyneGUISingleton.diagramManager.showReferencedConcept(fcde.GetDiagramElementID())
		})
		nullifyReferencedConceptItem := fyne.NewMenuItem("Nullify Referenced Concept", func() {
			FyneGUISingleton.diagramManager.nullifyReferencedConcept(fcde)
		})
		items = append(items, showReferencedConceptItem, nullifyReferencedConceptItem)
	case core.Refinement:
		showAbstractConceptItem := fyne.NewMenuItem("Show Abstract Concept", func() {
			FyneGUISingleton.diagramManager.showAbstractConcept(fcde.GetDiagramElementID())
		})
		showRefinedConceptItem := fyne.NewMenuItem("Show Refined Concept", func() {
			FyneGUISingleton.diagramManager.showRefinedConcept(fcde.GetDiagramElementID())
		})
		items = append(items, showAbstractConceptItem, showRefinedConceptItem)
	}
	deleteConceptViewItem := fyne.NewMenuItem("Delete Concept View", func() {
		FyneGUISingleton.diagramManager.deleteConceptView(fcde.GetDiagramElementID())
	})

	editFormatItem := fyne.NewMenuItem("Edit Format", func() {
		ShowFyneFormatDialog(fcde.GetFyneProperties(), func(properties diagramwidget.DiagramElementProperties) {
			fcde.SetFyneProperties(properties)
			setCrlDiagramElementProperties(fcde.GetDiagramElementID(), properties)
			fcde.Refresh()

		})
	})
	copyFormatItem := fyne.NewMenuItem("Copy Format", func() {
		if FyneGUISingleton.propertiesClipboard == nil {
			FyneGUISingleton.propertiesClipboard = &diagramwidget.DiagramElementProperties{}
		}
		*(FyneGUISingleton.propertiesClipboard) = fcde.GetFyneProperties()
	})
	pasteFormatItem := fyne.NewMenuItem("Paste Format", func() {
		if FyneGUISingleton.propertiesClipboard != nil {
			fcde.SetFyneProperties(*(FyneGUISingleton.propertiesClipboard))
			setCrlDiagramElementProperties(fcde.GetDiagramElementID(), *(FyneGUISingleton.propertiesClipboard))
			fcde.Refresh()
		}
	})
	items = append(items, deleteConceptViewItem, editFormatItem, copyFormatItem, pasteFormatItem)
	bringForwardItem := fyne.NewMenuItem("Bring Forward", func() {
		FyneGUISingleton.diagramManager.GetSelectedDiagram().BringForward(fcde.GetDiagramElementID())
	})
	bringToFrontItem := fyne.NewMenuItem("Bring To Front", func() {
		FyneGUISingleton.diagramManager.GetSelectedDiagram().BringToFront(fcde.GetDiagramElementID())
	})
	sendBackwardItem := fyne.NewMenuItem("Send Backward", func() {
		FyneGUISingleton.diagramManager.GetSelectedDiagram().SendBackward(fcde.GetDiagramElementID())
	})
	sendToBackItem := fyne.NewMenuItem("Send To Back", func() {
		FyneGUISingleton.diagramManager.GetSelectedDiagram().SendToBack(fcde.GetDiagramElementID())
	})
	items = append(items, bringForwardItem, bringToFrontItem, sendBackwardItem, sendToBackItem)
	menu := fyne.NewMenu("Diagram Element Popup", items...)
	popup := widget.NewPopUpMenu(menu, FyneGUISingleton.window.Canvas())
	popup.Move(event.AbsolutePosition)
	popup.Show()
}

// MouseUp responds to mouse up events
func (fcdn *FyneCrlDiagramNode) MouseUp(event *desktop.MouseEvent) {
}

func (fcdn *FyneCrlDiagramNode) nodeMoved() {
	newPosition := fcdn.Position()
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	currentX := fcdn.crlDiagramNode.GetNodeX(trans)
	currentY := fcdn.crlDiagramNode.GetNodeY(trans)
	if newPosition.X != float32(currentX) {
		fcdn.crlDiagramNode.SetNodeX(float64(newPosition.X), trans)
	}
	if newPosition.Y != float32(currentY) {
		fcdn.crlDiagramNode.SetNodeY(float64(newPosition.Y), trans)
	}
}

// SetFyneProperties sets the fyne DiagramElementProperties of the diagram link
func (fcdn *FyneCrlDiagramNode) SetFyneProperties(properties diagramwidget.DiagramElementProperties) {
	fcdn.SetProperties(properties)
}

var _ diagramwidget.DiagramLink = (*FyneCrlDiagramLink)(nil)

// FyneCrlDiagramLink is an extension to the diagramwidget.DiagramLink that serves as a binding between
// the fyne link and the crldiagramdomain link
type FyneCrlDiagramLink struct {
	diagramwidget.BaseDiagramLink
	diagramLink       *crldiagramdomain.CrlDiagramLink
	modelElement      *core.Concept
	labelAnchoredText *diagramwidget.AnchoredText
	linkType          ToolbarSelection
}

// NewFyneCrlDiagramLink creates a fyne link that corresponds to the supplied crldiagramdomain link
func NewFyneCrlDiagramLink(diagramWidget *diagramwidget.DiagramWidget, link *crldiagramdomain.CrlDiagramLink, trans *core.Transaction) *FyneCrlDiagramLink {
	diagramLink := &FyneCrlDiagramLink{}
	diagramLink.diagramLink = link
	diagramLink.modelElement = link.ToCrlDiagramElement().GetReferencedModelConcept(trans)
	diagramwidget.InitializeBaseDiagramLink(diagramLink, diagramWidget, link.ToCore().GetConceptID(trans))
	// Display labels are not appropriate for pointers
	if !link.ToCore().IsRefinementOfURI(crldiagramdomain.CrlDiagramPointerURI, trans) {
		linkLabel := link.ToCrlDiagramElement().GetDisplayLabel(trans)
		diagramLink.labelAnchoredText = diagramLink.AddMidpointAnchoredText(displayLabel, linkLabel)
		displayedTextBinding := diagramLink.labelAnchoredText.GetDisplayedTextBinding()
		displayedTextBinding.Set(linkLabel)
		displayedTextBinding.AddListener(binding.NewDataListener(func() { diagramLink.labelChanged() }))
	}
	if link.ToCore().IsRefinementOfURI(crldiagramdomain.CrlDiagramReferenceLinkURI, trans) {
		diagramLink.AddTargetDecoration(createReferenceArrowhead())
		diagramLink.AddSourceDecoration(createDiamond())
		diagramLink.linkType = ReferenceLinkSelected
	} else if link.ToCore().IsRefinementOfURI(crldiagramdomain.CrlDiagramAbstractPointerURI, trans) {
		diagramLink.AddSourceDecoration(createRefinementTriangle())
		diagramLink.linkType = AbstractElementPointerSelected
	} else if link.ToCore().IsRefinementOfURI(crldiagramdomain.CrlDiagramElementPointerURI, trans) {
		diagramLink.AddTargetDecoration(createReferenceArrowhead())
		diagramLink.linkType = ReferencedElementPointerSelected
	} else if link.ToCore().IsRefinementOfURI(crldiagramdomain.CrlDiagramOwnerPointerURI, trans) {
		diagramLink.AddTargetDecoration(createDiamond())
		diagramLink.linkType = OwnerPointerSelected
	} else if link.ToCore().IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinedPointerURI, trans) {
		diagramLink.AddSourceDecoration(createMirrorRefinementTriangle())
		diagramLink.linkType = RefinedElementPointerSelected
	} else if link.ToCore().IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinementLinkURI, trans) {
		diagramLink.AddMidpointDecoration(createRefinementTriangle())
		diagramLink.linkType = RefinementLinkSelected
	}
	// Some remedial work here for crlLinks that were initially saved without a fgColor, with the assumption
	// that links never have a transparent color
	black := color.RGBA{0, 0, 0, 255}
	grey := color.RGBA{153, 153, 153, 255}
	fgColor := link.ToCrlDiagramElement().GetLineColor(trans)
	if fgColor == "" {
		if link.ToCore().IsRefinementOfURI(crldiagramdomain.CrlDiagramPointerURI, trans) {
			fgColor = getCrlColor(grey)
		} else {
			fgColor = getCrlColor(black)
		}
		link.ToCrlDiagramElement().SetLineColor(fgColor, trans)
	}
	bgColor := link.ToCrlDiagramElement().GetBGColor(trans)
	diagramLink.SetForegroundColor(getGoColor(fgColor))
	diagramLink.SetBackgroundColor(getGoColor(bgColor))
	diagramLink.Refresh()
	return diagramLink
}

// GetDiagramElement returns the crl diagram element associated with the link
func (fcdl *FyneCrlDiagramLink) GetDiagramElement() *crldiagramdomain.CrlDiagramElement {
	return fcdl.diagramLink.ToCrlDiagramElement()
}

// GetDiagramElementID returns the ID of the crl diagram eleent associagted with the link
func (fcdl *FyneCrlDiagramLink) GetDiagramElementID() string {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	return fcdl.diagramLink.ToCore().GetConceptID(trans)
}

// GetModelElement returns the crl model element represented by the link
func (fcdl *FyneCrlDiagramLink) GetModelElement() *core.Concept {
	return fcdl.modelElement
}

// GetModelElementID returns the ID of the crl model eleent represented by the link
func (fcdl *FyneCrlDiagramLink) GetModelElementID() string {
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	return fcdl.modelElement.GetConceptID(trans)
}

// GetFyneProperties returns the fyne DiagramElementProperties of the diagram link
func (fcdl *FyneCrlDiagramLink) GetFyneProperties() diagramwidget.DiagramElementProperties {
	return fcdl.GetProperties()
}

func (fcdl *FyneCrlDiagramLink) labelChanged() {
	newValue, _ := fcdl.labelAnchoredText.GetDisplayedTextBinding().Get()
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	if fcdl.modelElement != nil {
		if fcdl.modelElement.GetLabel(trans) != newValue {
			fcdl.modelElement.SetLabel(newValue, trans)
		}
	}
}

func (fcdl *FyneCrlDiagramLink) referencePositionChanged() {
	newValue, _ := fcdl.labelAnchoredText.GetDisplayedTextBinding().Get()
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	if fcdl.modelElement != nil {
		if fcdl.modelElement.GetLabel(trans) != newValue {
			fcdl.modelElement.SetLabel(newValue, trans)
		}
	}
}

// SetLabel sets the label for the fyne DiagramLink
func (fcdl *FyneCrlDiagramLink) SetLabel(label string) {
	fcdl.labelAnchoredText.GetDisplayedTextBinding().Set(label)
}

// SetFyneProperties sets the fyne DiagramElementProperties of the diagram link
func (fcdl *FyneCrlDiagramLink) SetFyneProperties(properties diagramwidget.DiagramElementProperties) {
	fcdl.SetProperties(properties)
}
