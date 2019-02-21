package editor

import (
	"reflect"
	"strconv"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagram"
	//	"log"
)

func addDiagramViewFunctionsToUofD(uOfD core.UniverseOfDiscourse) {
	uOfD.AddFunction(crldiagram.CrlDiagramURI, updateDiagramView)
	uOfD.AddFunction(crldiagram.CrlDiagramElementURI, updateDiagramElementView)
}

func getLinkAdditionalParameters(link core.Element, hl *core.HeldLocks) map[string]string {
	modelElement := crldiagram.GetReferencedModelElement(link, hl)
	modelElementType := ""
	if modelElement != nil {
		modelElementType = reflect.TypeOf(modelElement).String()
	}
	additionalParameters := map[string]string{
		"DisplayLabel": crldiagram.GetDisplayLabel(link, hl),
		"Icon":         GetIconPath(crldiagram.GetReferencedModelElement(link, hl), hl),
		"OwnerID":      link.GetOwningConceptID(hl),
		"Abstractions": crldiagram.GetAbstractionDisplayLabel(link, hl),
		"LinkType":     modelElementType,
		"LinkSourceID": crldiagram.GetLinkSource(link, hl).GetConceptID(hl),
		"LinkTargetID": crldiagram.GetLinkTarget(link, hl).GetConceptID(hl)}
	return additionalParameters
}

func getNodeAdditionalParameters(node core.Element, hl *core.HeldLocks) map[string]string {
	additionalParameters := map[string]string{
		"DisplayLabel":        crldiagram.GetDisplayLabel(node, hl),
		"DisplayLabelYOffset": strconv.FormatFloat(crldiagram.GetDisplayLabelYOffset(node, hl), 'f', -1, 64),
		"NodeHeight":          strconv.FormatFloat(crldiagram.GetNodeHeight(node, hl), 'f', -1, 64),
		"NodeWidth":           strconv.FormatFloat(crldiagram.GetNodeWidth(node, hl), 'f', -1, 64),
		"NodeX":               strconv.FormatFloat(crldiagram.GetNodeX(node, hl), 'f', -1, 64),
		"NodeY":               strconv.FormatFloat(crldiagram.GetNodeY(node, hl), 'f', -1, 64),
		"Icon":                GetIconPath(crldiagram.GetReferencedModelElement(node, hl), hl),
		"OwnerID":             node.GetOwningConceptID(hl),
		"Abstractions":        crldiagram.GetAbstractionDisplayLabel(node, hl)}
	return additionalParameters
}

func updateDiagramElementView(diagramElement core.Element, changeNotification *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.ReadLockElement(diagramElement)
	if diagramElement.IsRefinementOfURI(crldiagram.CrlDiagramNodeURI, hl) {
		switch changeNotification.GetNatureOfChange() {
		case core.ChildChanged:
			additionalParameters := getNodeAdditionalParameters(diagramElement, hl)
			CrlEditorSingleton.SendNotification("UpdateDiagramNode", diagramElement.GetConceptID(hl), diagramElement, additionalParameters)
		}
		return
	} else if diagramElement.IsRefinementOfURI(crldiagram.CrlDiagramLinkURI, hl) {
		switch changeNotification.GetNatureOfChange() {
		case core.ChildChanged:
			additionalParameters := getLinkAdditionalParameters(diagramElement, hl)
			CrlEditorSingleton.SendNotification("UpdateDiagramLink", diagramElement.GetConceptID(hl), diagramElement, additionalParameters)
		}
		return
	}
}

func updateDiagramView(el core.Element, changeNotifications *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	go updateDiagramViewInternal(el, changeNotifications, uOfD)
}

func updateDiagramViewInternal(el core.Element, changeNotifications *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	// changeNotification.Print("updateDiagramView ", hl)
	// Check whether the name has changed
	// Check to see whether a diagram node or edge has been added or removed from the diagram
}
