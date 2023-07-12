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

var _ fyne.Widget = (*FyneCrlDiagramNode)(nil)
var _ diagramwidget.DiagramElement = (*FyneCrlDiagramNode)(nil)
var _ diagramwidget.DiagramNode = (*FyneCrlDiagramNode)(nil)
var _ fyne.Tappable = (*FyneCrlDiagramNode)(nil)

// FyneCrlDiagramNode is an extension to diagramwidget.DiagramNode that serves as a binding
// between the diagramwidget nodes and the crldiagramdomain diagram noddes
type FyneCrlDiagramNode struct {
	diagramwidget.BaseDiagramNode
	diagramElement  core.Element
	modelElement    core.Element
	entryWidget     *widget.Entry
	abstractionText *canvas.Text
	labelBinding    binding.String
	// abstractionTextBinding binding.String
}

// NewFyneCrlDiagramNode creates a fyne node that corresponds to the supplied crldiagram node
func NewFyneCrlDiagramNode(node core.Element, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramNode {
	newNode := &FyneCrlDiagramNode{}
	nodeID := node.GetConceptID(trans)
	newNode.diagramElement = node
	newNode.modelElement = crldiagramdomain.GetReferencedModelConcept(node, trans)
	nodeIcon := widget.NewIcon(getIconResource(newNode.modelElement, trans))

	abstractionString := crldiagramdomain.GetAbstractionDisplayLabel(node, trans)
	newNode.abstractionText = canvas.NewText(abstractionString, color.Black)
	newNode.abstractionText.TextStyle = fyne.TextStyle{Bold: false, Italic: true, Monospace: false, Symbol: false, TabWidth: 4}

	hBox := container.NewHBox(nodeIcon, newNode.abstractionText)
	nodeLabel := crldiagramdomain.GetDisplayLabel(node, trans)
	newNode.labelBinding = binding.NewString()
	newNode.labelBinding.Set(nodeLabel)
	newNode.labelBinding.AddListener(binding.NewDataListener(func() { newNode.labelChanged() }))
	newNode.entryWidget = widget.NewEntryWithData(newNode.labelBinding)
	newNode.entryWidget.Wrapping = fyne.TextWrapOff
	newNode.entryWidget.Validator = nil
	newNode.entryWidget.Refresh()

	newNode.MovedCallback = func() {
		newNode.nodeMoved()
	}

	nodeContainer := container.NewVBox(hBox, newNode.entryWidget)

	diagramwidget.InitializeBaseDiagramNode(newNode, diagramWidget, nodeContainer, nodeID)
	// Size isn't available until after initialization
	newNode.abstractionText.TextSize = newNode.GetProperties().CaptionTextSize
	x := crldiagramdomain.GetNodeX(node, trans)
	y := crldiagramdomain.GetNodeY(node, trans)
	fynePosition := fyne.NewPos(float32(x), float32(y))
	newNode.Move(fynePosition)
	newNode.Refresh()
	return newNode
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

func (fcdn *FyneCrlDiagramNode) MouseDown(event *desktop.MouseEvent) {
	if event.Button == desktop.MouseButtonSecondary {
		items := []*fyne.MenuItem{}
		showModelConceptItem := fyne.NewMenuItem("Show Concept in Navigator", func() {
			FyneGUISingleton.treeManager.ShowElementInTree(fcdn.modelElement)
		})
		items = append(items, showModelConceptItem)
		showDiagramElementItem := fyne.NewMenuItem("Show Diagram Element in Navigator", func() {
			FyneGUISingleton.treeManager.ShowElementInTree(fcdn.diagramElement)
		})
		items = append(items, showDiagramElementItem)
		showOwnerItem := fyne.NewMenuItem("Show Owner", func() {
			FyneGUISingleton.diagramManager.showOwner(fcdn.GetDiagramElementID())
		})
		items = append(items, showOwnerItem)
		showOwnedConceptsItem := fyne.NewMenuItem("Show Owned Conecpts", func() {
			FyneGUISingleton.diagramManager.showOwnedConcepts(fcdn.GetDiagramElementID())
		})
		items = append(items, showOwnedConceptsItem)
		switch fcdn.modelElement.(type) {
		case core.Reference:
			showReferencedConceptItem := fyne.NewMenuItem("Show Referenced Concept", func() {
				FyneGUISingleton.diagramManager.showReferencedConcept(fcdn.GetDiagramElementID())
			})
			nullifyReferencedConceptItem := fyne.NewMenuItem("Nullify Referenced Concept", func() {
				FyneGUISingleton.diagramManager.nullifyReferencedConcept(fcdn)
			})
			items = append(items, showReferencedConceptItem, nullifyReferencedConceptItem)
		case core.Refinement:
			showAbstractConceptItem := fyne.NewMenuItem("Show Abstract Concept", func() {
				FyneGUISingleton.diagramManager.showAbstractConcept(fcdn.GetDiagramElementID())
			})
			showRefinedConceptItem := fyne.NewMenuItem("Show Refined Concept", func() {
				FyneGUISingleton.diagramManager.showRefinedConcept(fcdn.GetDiagramElementID())
			})
			items = append(items, showAbstractConceptItem, showRefinedConceptItem)
		}
		deleteDiagramElementViewItem := fyne.NewMenuItem("Delete Diagram Element View", func() {
			FyneGUISingleton.diagramManager.deleteDiagramElementView(fcdn.GetDiagramElementID())
		})
		// <a class="show" onclick="crlBringToFront()">Bring To Front</a>
		editFormatItem := fyne.NewMenuItem("Edit Format", func() {
			ShowFyneFormatDialog(fcdn.GetProperties(), func(properties diagramwidget.DiagramElementProperties) {
				fcdn.SetProperties(properties)
				fcdn.Refresh()
				// fcdn.GetDiagram().ForceRepaint()
			})
		})
		// <a class="show" onclick="crlEditFormat()">Edit Format</a>
		// <a class="show" onclick="crlCopyFormat()">Copy Format</a>
		// <a class="show" onclick="crlPasteFormat()">Paste Format</a>
		items = append(items, deleteDiagramElementViewItem, editFormatItem)
		menu := fyne.NewMenu("Diagram Element Popup", items...)
		popup := widget.NewPopUpMenu(menu, FyneGUISingleton.window.Canvas())
		popup.Move(event.AbsolutePosition)
		popup.Show()
	}
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
	currentX := crldiagramdomain.GetNodeX(fcdn.diagramElement, trans)
	currentY := crldiagramdomain.GetNodeY(fcdn.diagramElement, trans)
	if newPosition.X != float32(currentX) {
		crldiagramdomain.SetNodeX(fcdn.diagramElement, float64(newPosition.X), trans)
	}
	if newPosition.Y != float32(currentY) {
		crldiagramdomain.SetNodeY(fcdn.diagramElement, float64(newPosition.Y), trans)
	}
}

var _ diagramwidget.DiagramLink = (*FyneCrlDiagramLink)(nil)

// FyneCrlDiagramLink is an extension to the diagramwidget.DiagramLink that serves as a binding between
// the fyne link and the crldiagramdomain link
type FyneCrlDiagramLink struct {
	diagramwidget.BaseDiagramLink
	diagramElement    core.Element
	modelElement      core.Element
	labelAnchoredText *diagramwidget.AnchoredText
	linkType          ToolbarSelection
}

// NewFyneCrlDiagramLink creates a fyne link that corresponds to the supplied crldiagramdomain link
func NewFyneCrlDiagramLink(diagramWidget *diagramwidget.DiagramWidget, link core.Element, trans *core.Transaction) *FyneCrlDiagramLink {
	diagramLink := &FyneCrlDiagramLink{}
	diagramLink.diagramElement = link
	diagramLink.modelElement = crldiagramdomain.GetReferencedModelConcept(link, trans)
	diagramwidget.InitializeBaseDiagramLink(diagramLink, diagramWidget, link.GetConceptID(trans))
	// Display labels are not appropriate for pointers
	if !link.IsRefinementOfURI(crldiagramdomain.CrlDiagramPointerURI, trans) {
		linkLabel := crldiagramdomain.GetDisplayLabel(link, trans)
		diagramLink.labelAnchoredText = diagramLink.AddMidpointAnchoredText(displayLabel, linkLabel)
		displayedTextBinding := diagramLink.labelAnchoredText.GetDisplayedTextBinding()
		displayedTextBinding.Set(linkLabel)
		displayedTextBinding.AddListener(binding.NewDataListener(func() { diagramLink.labelChanged() }))
	}
	grey := color.RGBA{153, 153, 153, 255}
	if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramReferenceLinkURI, trans) {
		diagramLink.AddTargetDecoration(createReferenceArrowhead())
		diagramLink.AddSourceDecoration(createDiamond())
		diagramLink.linkType = REFERENCE_LINK
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramAbstractPointerURI, trans) {
		diagramLink.AddSourceDecoration(createRefinementTriangle())
		diagramLink.SetForegroundColor(grey)
		diagramLink.linkType = ABSTRACT_ELEMENT_POINTER
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementPointerURI, trans) {
		diagramLink.AddTargetDecoration(createReferenceArrowhead())
		diagramLink.SetForegroundColor(grey)
		diagramLink.linkType = REFERENCED_ELEMENT_POINTER
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramOwnerPointerURI, trans) {
		diagramLink.AddTargetDecoration(createDiamond())
		diagramLink.SetForegroundColor(grey)
		diagramLink.linkType = OWNER_POINTER
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinedPointerURI, trans) {
		diagramLink.AddSourceDecoration(createMirrorRefinementTriangle())
		diagramLink.SetForegroundColor(grey)
		diagramLink.linkType = REFINED_ELEMENT_POINTER
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinementLinkURI, trans) {
		diagramLink.AddMidpointDecoration(createRefinementTriangle())
		diagramLink.linkType = REFINEMENT_LINK
	}
	diagramLink.Refresh()
	return diagramLink
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

func (fcdl *FyneCrlDiagramLink) SetLabel(label string) {
	fcdl.labelAnchoredText.GetDisplayedTextBinding().Set(label)
}
