package browsergui

import (
	"github.com/pkg/errors"
	"strconv"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	// "github.com/pbrown12303/activeCRL/crleditorbrowserguidomain"
	//	"log"
)

// func addDiagramViewFunctionsToUofD(uOfD *core.UniverseOfDiscourse) {
// 	uOfD.AddFunction(crldiagramdomain.CrlDiagramURI, updateDiagramView)
// 	uOfD.AddFunction(crldiagramdomain.CrlDiagramElementURI, updateDiagramElementView)
// }

// diagramElementManager manages the GUI display of diagram elements
type diagramElementManager struct {
	diagramManager  *diagramManager
	diagramElements map[string]core.Element
}

func newDiagramElementManager(diagramManager *diagramManager) *diagramElementManager {
	dem := &diagramElementManager{}
	dem.diagramManager = diagramManager
	dem.diagramElements = map[string]core.Element{}
	return dem
}

func getLinkAdditionalParameters(link core.Element, hl *core.HeldLocks) map[string]string {
	var linkType string
	var represents string
	if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinementLinkURI, hl) {
		linkType = "RefinementLink"
		represents = "Refinement"
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramReferenceLinkURI, hl) {
		linkType = "ReferenceLink"
		represents = "Reference"
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramOwnerPointerURI, hl) {
		linkType = "OwnerPointer"
		represents = "OwnerPointer"
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementPointerURI, hl) {
		linkType = "ElementPointer"
		represents = "ElementPointer"
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramAbstractPointerURI, hl) {
		linkType = "AbstractPointer"
		represents = "AbstractPointer"
	} else if link.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinedPointerURI, hl) {
		linkType = "RefinedPointer"
		represents = "RefinedPointer"
	}
	linkSource := crldiagramdomain.GetLinkSource(link, hl)
	linkSourceID := ""
	if linkSource != nil {
		linkSourceID = linkSource.GetConceptID(hl)
	}
	linkTarget := crldiagramdomain.GetLinkTarget(link, hl)
	linkTargetID := ""
	if linkTarget != nil {
		linkTargetID = linkTarget.GetConceptID(hl)
	}
	additionalParameters := map[string]string{
		"DisplayLabel": crldiagramdomain.GetDisplayLabel(link, hl),
		"Icon":         GetIconPath(crldiagramdomain.GetReferencedModelElement(link, hl), hl),
		"OwnerID":      link.GetOwningConceptID(hl),
		"Abstractions": crldiagramdomain.GetAbstractionDisplayLabel(link, hl),
		"LinkType":     linkType,
		"LinkSourceID": linkSourceID,
		"LinkTargetID": linkTargetID,
		"Represents":   represents}
	return additionalParameters
}

func getNodeAdditionalParameters(node core.Element, hl *core.HeldLocks) map[string]string {
	var represents string
	referencedModelElement := crldiagramdomain.GetReferencedModelElement(node, hl)
	if referencedModelElement != nil {
		switch referencedModelElement.(type) {
		case core.Literal:
			represents = "Literal"
		case core.Reference:
			represents = "Reference"
		case core.Refinement:
			represents = "Refinement"
		case core.Element:
			represents = "Element"
		}
	}
	additionalParameters := map[string]string{
		"DisplayLabel":        crldiagramdomain.GetDisplayLabel(node, hl),
		"DisplayLabelYOffset": strconv.FormatFloat(crldiagramdomain.GetDisplayLabelYOffset(node, hl), 'f', -1, 64),
		"NodeHeight":          strconv.FormatFloat(crldiagramdomain.GetNodeHeight(node, hl), 'f', -1, 64),
		"NodeWidth":           strconv.FormatFloat(crldiagramdomain.GetNodeWidth(node, hl), 'f', -1, 64),
		"NodeX":               strconv.FormatFloat(crldiagramdomain.GetNodeX(node, hl), 'f', -1, 64),
		"NodeY":               strconv.FormatFloat(crldiagramdomain.GetNodeY(node, hl), 'f', -1, 64),
		"Icon":                GetIconPath(referencedModelElement, hl),
		"OwnerID":             node.GetOwningConceptID(hl),
		"Abstractions":        crldiagramdomain.GetAbstractionDisplayLabel(node, hl),
		"Represents":          represents,
		"LineColor":           crldiagramdomain.GetLineColor(node, hl),
		"BGColor":             crldiagramdomain.GetBGColor(node, hl)}
	return additionalParameters
}

// Update updates the client display of the diagram based on changes to the diagramElement
func (demPtr *diagramElementManager) Update(changeNotification *core.ChangeNotification, hl *core.HeldLocks) error {
	uOfD := hl.GetUniverseOfDiscourse()
	// if the reportingElementState is nil, this notification comes from the uOfD. We can ignore these (I think)
	if changeNotification.GetReportingElementState() == nil {
		return nil
	}
	diagramElement := uOfD.GetElement(changeNotification.GetReportingElementID())
	if diagramElement == nil {
		return errors.New("In DiagramElementManager.Update, diagramElement not found in uOfD")
	}
	hl.ReadLockElement(diagramElement)
	if diagramElement.GetUniverseOfDiscourse(hl) != uOfD {
		// The diagram element has been removed from the universe of discourse
		priorState := changeNotification.GetBeforeConceptState()
		if priorState != nil {
			additionalParameters := map[string]string{"OwnerID": priorState.OwningConceptID}
			SendNotification("DeleteDiagramElement", diagramElement.GetConceptID(hl), priorState, additionalParameters)
		}
		diagramElement.Deregister(demPtr)
		delete(demPtr.diagramElements, diagramElement.GetConceptID(hl))
		return nil
	}
	if crldiagramdomain.IsDiagramNode(diagramElement, hl) {
		switch changeNotification.GetNatureOfChange() {
		case core.ForwardedChange:
			additionalParameters := getNodeAdditionalParameters(diagramElement, hl)
			conceptState, err := core.NewConceptState(diagramElement)
			if err != nil {
				return errors.Wrap(err, "DiagramView.go updateDiagrmElementView failed")
			}
			SendNotification("UpdateDiagramNode", diagramElement.GetConceptID(hl), conceptState, additionalParameters)
		case core.ConceptChanged:
			currentConcept := changeNotification.GetAfterConceptState()
			priorConcept := changeNotification.GetBeforeConceptState()
			if currentConcept != nil && priorConcept != nil && currentConcept.Label != priorConcept.Label {
				currentLabel := currentConcept.Label
				diagramElement.SetLabel(currentLabel, hl)
				crldiagramdomain.SetDisplayLabel(diagramElement, currentLabel, hl)
				additionalParameters := getNodeAdditionalParameters(diagramElement, hl)
				conceptState, err := core.NewConceptState(diagramElement)
				if err != nil {
					return errors.Wrap(err, "DiagramView.go updateDiagrmElementView failed")
				}
				SendNotification("UpdateDiagramNode", diagramElement.GetConceptID(hl), conceptState, additionalParameters)
			}
		}
		return nil
	} else if crldiagramdomain.IsDiagramLink(diagramElement, hl) {
		switch changeNotification.GetNatureOfChange() {
		case core.ForwardedChange:
			additionalParameters := getLinkAdditionalParameters(diagramElement, hl)
			conceptState, err := core.NewConceptState(diagramElement)
			if err != nil {
				return errors.Wrap(err, "DiagramView.go updateDiagrmElementView failed")
			}
			SendNotification("UpdateDiagramLink", diagramElement.GetConceptID(hl), conceptState, additionalParameters)
		case core.ConceptChanged:
			currentConcept := changeNotification.GetAfterConceptState()
			priorConcept := changeNotification.GetBeforeConceptState()
			if currentConcept != nil && priorConcept != nil && currentConcept.Label != priorConcept.Label {
				currentLabel := currentConcept.Label
				diagramElement.SetLabel(currentLabel, hl)
				crldiagramdomain.SetDisplayLabel(diagramElement, currentLabel, hl)
				additionalParameters := getLinkAdditionalParameters(diagramElement, hl)
				conceptState, err := core.NewConceptState(diagramElement)
				if err != nil {
					return errors.Wrap(err, "DiagramView.go updateDiagrmElementView failed")
				}
				SendNotification("UpdateDiagramLink", diagramElement.GetConceptID(hl), conceptState, additionalParameters)
			}
		}
		return nil
	}
	return nil
}

// // diagramViewMonitor is the callback function that manages the diagram view in the gui.
// func diagramViewMonitor(instance core.Element, changeNotification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
// 	// The instance here is the reference that is monitoring the diagram
// 	hl := uOfD.NewHeldLocks()
// 	defer hl.ReleaseLocks()

// 	switch changeNotification.GetNatureOfChange() {
// 	case core.ConceptChanged:
// 		// When the diagram is deleted, the reference to it in this object becomes nil resulting in a ConceptChanged
// 		switch instance.(type) {
// 		case core.Reference:
// 			if changeNotification.GetAfterState().ReferencedConceptID == "" && changeNotification.GetBeforeState().ReferencedConceptID != "" {
// 				err := BrowserGUISingleton.editor.CloseDiagramView(changeNotification.GetBeforeState().ReferencedConceptID, hl)
// 				if err != nil {
// 					errors.Wrap(err, "DiagramView.go diagramViewMonitor failed")
// 				}
// 				uOfD.DeleteElement(instance, hl)
// 			}
// 		}
// 	case core.IndicatedConceptChanged:
// 		underlyingChange := changeNotification.GetUnderlyingChange()
// 		switch underlyingChange.GetNatureOfChange() {
// 		case core.ConceptChanged:
// 			diagram := uOfD.GetElement(underlyingChange.GetAfterState().ConceptID)
// 			if diagram != nil && underlyingChange.GetBeforeState().Label != underlyingChange.GetAfterState().Label {
// 				diagramID := underlyingChange.GetAfterState().ConceptID
// 				BrowserGUISingleton.getDiagramManager().diagramLabelChanged(diagramID, diagram, hl)
// 			}
// 		}
// 	}
// 	return nil
// }

// func registerDiagramViewMonitorFunctions(uOfD *core.UniverseOfDiscourse) {
// 	uOfD.AddFunction(crleditorbrowserguidomain.DiagramViewMonitorURI, diagramViewMonitor)
// }
