package editor

import (
	"strconv"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagram"
	//	"log"
)

func addDiagramViewFunctionsToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {
	uOfD.AddFunction(crldiagram.CrlDiagramURI, updateDiagramView)
	uOfD.AddFunction(crldiagram.CrlDiagramNodeURI, updateDiagramNodeView)
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

func updateDiagramNodeView(node core.Element, changeNotification *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.ReadLockElement(node)
	switch changeNotification.GetNatureOfChange() {
	case core.ChildChanged:
		additionalParameters := getNodeAdditionalParameters(node, hl)
		CrlEditorSingleton.GetClientNotificationManager().SendNotification("UpdateDiagramNode", node.GetConceptID(hl), node, additionalParameters)
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

func init() {
	//	core.GetCore().AddFunction(crlDiagram.CrlDiagramUri, updateDiagramView)
}
