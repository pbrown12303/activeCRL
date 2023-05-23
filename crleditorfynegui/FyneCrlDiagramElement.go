package crleditorfynegui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/diagramwidget"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
)

var _ fyne.Widget = (*FyneCrlDiagramNode)(nil)
var _ diagramwidget.DiagramElement = (*FyneCrlDiagramNode)(nil)
var _ fyne.Tappable = (*FyneCrlDiagramNode)(nil)

type FyneCrlDiagramNode struct {
	diagramwidget.BaseDiagramNode
	diagramElement  core.Element
	modelElement    core.Element
	entryWidget     *widget.Entry
	abstractionText *canvas.Text
	labelBinding    binding.String
	// abstractionTextBinding binding.String
}

func NewFyneCrlDiagramNode(node core.Element, trans *core.Transaction, diagramWidget *diagramwidget.DiagramWidget) diagramwidget.DiagramNode {
	newNode := &FyneCrlDiagramNode{}
	nodeID := node.GetConceptID(trans)
	newNode.diagramElement = node
	newNode.modelElement = crldiagramdomain.GetReferencedModelElement(node, trans)
	nodeIcon := widget.NewIcon(getIconResource(newNode.modelElement, trans))

	abstractionString := crldiagramdomain.GetAbstractionDisplayLabel(node, trans)
	newNode.abstractionText = canvas.NewText(abstractionString, color.Black)
	newNode.abstractionText.TextSize = diagramWidget.DiagramTheme.Size(theme.SizeNameCaptionText)
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
	fcdn.modelElement.SetLabel(newValue, trans)
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

type FyneCrlDiagramLink struct {
	diagramwidget.BaseDiagramLink
	diagramElement    core.Element
	modelElement      core.Element
	labelAnchoredText *diagramwidget.AnchoredText
}

func NewFyneCrlDiagramLink(diagramWidget *diagramwidget.DiagramWidget, fyneSourcePad diagramwidget.ConnectionPad, fyneTargetPad diagramwidget.ConnectionPad, link core.Element, trans *core.Transaction) *FyneCrlDiagramLink {
	diagramLink := &FyneCrlDiagramLink{}
	diagramLink.diagramElement = link
	diagramLink.modelElement = crldiagramdomain.GetReferencedModelElement(link, trans)
	diagramwidget.InitializeBaseDiagramLink(diagramLink, diagramWidget, fyneSourcePad, fyneTargetPad, link.GetConceptID(trans))
	linkLabel := crldiagramdomain.GetDisplayLabel(link, trans)
	diagramLink.labelAnchoredText = diagramLink.AddMidpointAnchoredText(displayLabel, "test")
	displayedTextBinding := diagramLink.labelAnchoredText.GetDisplayedTextBinding()
	displayedTextBinding.Set(linkLabel)
	displayedTextBinding.AddListener(binding.NewDataListener(func() { diagramLink.labelChanged() }))
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
	diagramLink.Refresh()
	return diagramLink
}

func (fcdl *FyneCrlDiagramLink) labelChanged() {
	newValue, _ := fcdl.labelAnchoredText.GetDisplayedTextBinding().Get()
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	if fcdl.modelElement.GetLabel(trans) != newValue {
		fcdl.modelElement.SetLabel(newValue, trans)
	}
}

func (fcdl *FyneCrlDiagramLink) SetLabel(label string) {
	fcdl.labelAnchoredText.GetDisplayedTextBinding().Set(label)
}
