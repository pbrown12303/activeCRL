package browsergui

import (
	"github.com/pkg/errors"
	"strconv"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditorbrowserguidomain"
	//	"log"
)

func addDiagramViewFunctionsToUofD(uOfD *core.UniverseOfDiscourse) {
	uOfD.AddFunction(crldiagramdomain.CrlDiagramURI, updateDiagramView)
	uOfD.AddFunction(crldiagramdomain.CrlDiagramElementURI, updateDiagramElementView)
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

// updateDiagramElementView attaches to the diagram element and updates the client display of the diagram based on changes to the diagramElement
func updateDiagramElementView(diagramElement core.Element, changeNotification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.ReadLockElement(diagramElement)
	if diagramElement.GetUniverseOfDiscourse(hl) != uOfD {
		// The diagram element has been removed from the universe of discourse
		priorState := changeNotification.GetBeforeState()
		if priorState != nil {
			additionalParameters := map[string]string{"OwnerID": priorState.OwningConceptID}
			SendNotification("DeleteDiagramElement", diagramElement.GetConceptID(hl), priorState, additionalParameters)
		}
		return nil
	}
	if crldiagramdomain.IsDiagramNode(diagramElement, hl) {
		switch changeNotification.GetNatureOfChange() {
		case core.ChildChanged:
			additionalParameters := getNodeAdditionalParameters(diagramElement, hl)
			conceptState, err := core.NewConceptState(diagramElement)
			if err != nil {
				return errors.Wrap(err, "DiagramView.go updateDiagrmElementView failed")
			}
			SendNotification("UpdateDiagramNode", diagramElement.GetConceptID(hl), conceptState, additionalParameters)
		case core.IndicatedConceptChanged:
			underlyingNotification := changeNotification.GetUnderlyingChange()
			if underlyingNotification != nil {
				switch underlyingNotification.GetNatureOfChange() {
				case core.IndicatedConceptChanged:
					secondUnderlyingNotification := underlyingNotification.GetUnderlyingChange()
					if secondUnderlyingNotification != nil {
						switch secondUnderlyingNotification.GetNatureOfChange() {
						case core.ConceptChanged:
							currentConcept := secondUnderlyingNotification.GetAfterState()
							priorConcept := secondUnderlyingNotification.GetBeforeState()
							if currentConcept != nil && priorConcept != nil && currentConcept.Label != priorConcept.Label {
								currentLabel := currentConcept.Label
								diagramElement.SetLabel(currentLabel, hl)
								crldiagramdomain.SetDisplayLabel(diagramElement, currentLabel, hl)
							}
						}
					}
				}
			}
		}
		return nil
	} else if crldiagramdomain.IsDiagramLink(diagramElement, hl) {
		switch changeNotification.GetNatureOfChange() {
		case core.ChildChanged:
			additionalParameters := getLinkAdditionalParameters(diagramElement, hl)
			conceptState, err := core.NewConceptState(diagramElement)
			if err != nil {
				return errors.Wrap(err, "DiagramView.go updateDiagrmElementView failed")
			}
			SendNotification("UpdateDiagramLink", diagramElement.GetConceptID(hl), conceptState, additionalParameters)
		}
		return nil
	}
	return nil
}

// updateDiagramView attaches to the diagram view and handles additions and removals of diagram elements from the diagram view
// Note that it cannot delete the diagram view in the GUI because this function will never get called: once the diagram has been
// deleted, queuing of functions related to it is suppressed. That's what the DiagramViewMonitor is for.
func updateDiagramView(diagram core.Element, changeNotification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.ReadLockElement(diagram)
	switch changeNotification.GetNatureOfChange() {
	case core.ChildChanged:
		// In the case of deletion, the reportingElement may be nil at this point, i.e. already removed from UofD
		reportingElementState := changeNotification.GetAfterState()
		underlyingChange := changeNotification.GetUnderlyingChange()
		if underlyingChange != nil {
			switch underlyingChange.GetNatureOfChange() {
			case core.ConceptChanged:
				// The diagram element itself has changed. Need to check to see if it is an ownership change that requires
				// removal from the diagram
				oldChildOwner := uOfD.GetElement(underlyingChange.GetBeforeState().OwningConceptID)
				currentChildOwner := uOfD.GetElement(underlyingChange.GetAfterState().OwningConceptID)
				if oldChildOwner == diagram && currentChildOwner != diagram {
					beforeState := underlyingChange.GetBeforeState()
					additionalParameters := map[string]string{"OwnerID": diagram.GetConceptID(hl)}
					SendNotification("DeleteDiagramElement", beforeState.ConceptID, beforeState, additionalParameters)
				} else if oldChildOwner != diagram && currentChildOwner == diagram {
					newElement := uOfD.GetElement(underlyingChange.GetAfterState().ConceptID)
					if newElement == nil {
						return errors.New("In DiagramView.go updateDiagramView, newElement was nil")
					}
					conceptState, err := core.NewConceptState(newElement)
					if err != nil {
						return errors.Wrap(err, "DiagramView.go updateDiagramView failed")
					}
					if crldiagramdomain.IsDiagramNode(newElement, hl) {
						additionalParameters := getNodeAdditionalParameters(newElement, hl)
						SendNotification("AddDiagramNode", newElement.GetConceptID(hl), conceptState, additionalParameters)
					} else if crldiagramdomain.IsDiagramLink(newElement, hl) {
						additionalParameters := getLinkAdditionalParameters(newElement, hl)
						SendNotification("AddDiagramLink", newElement.GetConceptID(hl), conceptState, additionalParameters)
					}
				}
			case core.ChildChanged:
				// a Child of the diagram element changed. Check to see whether it is a model reference and, if so, whether the model elemment still exists
				if crldiagramdomain.IsModelReference(uOfD.GetElement(underlyingChange.GetChangedConceptID()), hl) {
					secondUnderlyingChange := underlyingChange.GetUnderlyingChange()
					switch secondUnderlyingChange.GetNatureOfChange() {
					case core.ConceptChanged:
						// check to see whether the modelElement is nil
						modelReference := uOfD.GetElement(secondUnderlyingChange.GetChangedConceptID()).(core.Reference)
						if modelReference.GetReferencedConcept(hl) == nil {
							additionalParameters := map[string]string{"OwnerID": diagram.GetConceptID(hl)}
							SendNotification("DeleteDiagramElement", reportingElementState.ConceptID, reportingElementState, additionalParameters)
						}
					}
				}
			}
		}
	}
	return nil
}

// diagramViewMonitor is the callback function that manages the diagram view in the gui.
func diagramViewMonitor(instance core.Element, changeNotification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
	// The instance here is the reference that is monitoring the diagram
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocks()

	switch changeNotification.GetNatureOfChange() {
	case core.ConceptChanged:
		// When the diagram is deleted, the reference to it in this object becomes nil resulting in a ConceptChanged
		switch instance.(type) {
		case core.Reference:
			if changeNotification.GetAfterState().ReferencedConceptID == "" && changeNotification.GetBeforeState().ReferencedConceptID != "" {
				err := BrowserGUISingleton.editor.CloseDiagramView(changeNotification.GetBeforeState().ReferencedConceptID, hl)
				if err != nil {
					errors.Wrap(err, "DiagramView.go diagramViewMonitor failed")
				}
				uOfD.DeleteElement(instance, hl)
			}
		}
	case core.IndicatedConceptChanged:
		underlyingChange := changeNotification.GetUnderlyingChange()
		switch underlyingChange.GetNatureOfChange() {
		case core.ConceptChanged:
			diagram := uOfD.GetElement(underlyingChange.GetAfterState().ConceptID)
			if diagram != nil && underlyingChange.GetBeforeState().Label != underlyingChange.GetAfterState().Label {
				diagramID := underlyingChange.GetAfterState().ConceptID
				BrowserGUISingleton.getDiagramManager().diagramLabelChanged(diagramID, diagram, hl)
			}
		}
	}
	return nil
}

func registerDiagramViewMonitorFunctions(uOfD *core.UniverseOfDiscourse) {
	uOfD.AddFunction(crleditorbrowserguidomain.DiagramViewMonitorURI, diagramViewMonitor)
}
