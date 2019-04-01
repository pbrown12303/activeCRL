package editor

import (
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
	var linkType string
	var represents string
	if link.IsRefinementOfURI(crldiagram.CrlDiagramRefinementLinkURI, hl) {
		linkType = "RefinementLink"
		represents = "Refinement"
	} else if link.IsRefinementOfURI(crldiagram.CrlDiagramReferenceLinkURI, hl) {
		linkType = "ReferenceLink"
		represents = "Reference"
	} else if link.IsRefinementOfURI(crldiagram.CrlDiagramOwnerPointerURI, hl) {
		linkType = "OwnerPointer"
		represents = "OwnerPointer"
	} else if link.IsRefinementOfURI(crldiagram.CrlDiagramElementPointerURI, hl) {
		linkType = "ElementPointer"
		represents = "ElementPointer"
	} else if link.IsRefinementOfURI(crldiagram.CrlDiagramAbstractPointerURI, hl) {
		linkType = "AbstractPointer"
		represents = "AbstractPointer"
	} else if link.IsRefinementOfURI(crldiagram.CrlDiagramRefinedPointerURI, hl) {
		linkType = "RefinedPointer"
		represents = "RefinedPointer"
	}
	linkSource := crldiagram.GetLinkSource(link, hl)
	linkSourceID := ""
	if linkSource != nil {
		linkSourceID = linkSource.GetConceptID(hl)
	}
	linkTarget := crldiagram.GetLinkTarget(link, hl)
	linkTargetID := ""
	if linkTarget != nil {
		linkTargetID = linkTarget.GetConceptID(hl)
	}
	additionalParameters := map[string]string{
		"DisplayLabel": crldiagram.GetDisplayLabel(link, hl),
		"Icon":         GetIconPath(crldiagram.GetReferencedModelElement(link, hl), hl),
		"OwnerID":      link.GetOwningConceptID(hl),
		"Abstractions": crldiagram.GetAbstractionDisplayLabel(link, hl),
		"LinkType":     linkType,
		"LinkSourceID": linkSourceID,
		"LinkTargetID": linkTargetID,
		"Represents":   represents}
	return additionalParameters
}

func getNodeAdditionalParameters(node core.Element, hl *core.HeldLocks) map[string]string {
	var represents string
	referencedModelElement := crldiagram.GetReferencedModelElement(node, hl)
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
		"DisplayLabel":        crldiagram.GetDisplayLabel(node, hl),
		"DisplayLabelYOffset": strconv.FormatFloat(crldiagram.GetDisplayLabelYOffset(node, hl), 'f', -1, 64),
		"NodeHeight":          strconv.FormatFloat(crldiagram.GetNodeHeight(node, hl), 'f', -1, 64),
		"NodeWidth":           strconv.FormatFloat(crldiagram.GetNodeWidth(node, hl), 'f', -1, 64),
		"NodeX":               strconv.FormatFloat(crldiagram.GetNodeX(node, hl), 'f', -1, 64),
		"NodeY":               strconv.FormatFloat(crldiagram.GetNodeY(node, hl), 'f', -1, 64),
		"Icon":                GetIconPath(referencedModelElement, hl),
		"OwnerID":             node.GetOwningConceptID(hl),
		"Abstractions":        crldiagram.GetAbstractionDisplayLabel(node, hl),
		"Represents":          represents}
	return additionalParameters
}

// updateDiagramElementView attaches to the diagram element and updates the client display of the diagram based on changes to the diagramElement
func updateDiagramElementView(diagramElement core.Element, changeNotification *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.ReadLockElement(diagramElement)
	if diagramElement.GetUniverseOfDiscourse(hl) != uOfD {
		// The diagram element has been removed from the universe of discourse
		priorState := changeNotification.GetPriorState()
		if priorState != nil {
			additionalParameters := map[string]string{"OwnerID": priorState.GetOwningConceptID(hl)}
			SendNotification("DeleteDiagramElement", diagramElement.GetConceptID(hl), priorState, additionalParameters)
		}
		return
	}
	if crldiagram.IsDiagramNode(diagramElement, hl) {
		switch changeNotification.GetNatureOfChange() {
		case core.ChildChanged:
			additionalParameters := getNodeAdditionalParameters(diagramElement, hl)
			SendNotification("UpdateDiagramNode", diagramElement.GetConceptID(hl), diagramElement, additionalParameters)
		}
		return
	} else if crldiagram.IsDiagramLink(diagramElement, hl) {
		switch changeNotification.GetNatureOfChange() {
		case core.ChildChanged:
			additionalParameters := getLinkAdditionalParameters(diagramElement, hl)
			SendNotification("UpdateDiagramLink", diagramElement.GetConceptID(hl), diagramElement, additionalParameters)
		}
		return
	}
}

func updateDiagramView(diagram core.Element, changeNotification *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.ReadLockElement(diagram)
	switch changeNotification.GetNatureOfChange() {
	case core.ChildChanged:
		reportingElement := changeNotification.GetReportingElement()
		if reportingElement.IsRefinementOfURI(crldiagram.CrlDiagramElementURI, hl) {
			// The reporting element is one of the diagram elements
			underlyingChange := changeNotification.GetUnderlyingChange()
			if underlyingChange != nil {
				switch underlyingChange.GetNatureOfChange() {
				case core.ConceptChanged:
					// The diagram element itself has changed. Need to check to see if it is an ownership change that requires
					// removal from the diagram
					oldChildOwner := underlyingChange.GetPriorState().GetOwningConcept(hl)
					currentChildOwner := underlyingChange.GetReportingElement().GetOwningConcept(hl)
					if oldChildOwner == diagram && currentChildOwner != diagram {
						priorState := underlyingChange.GetPriorState()
						additionalParameters := map[string]string{"OwnerID": diagram.GetConceptID(hl)}
						SendNotification("DeleteDiagramElement", priorState.GetConceptID(hl), priorState, additionalParameters)
					} else if oldChildOwner != diagram && currentChildOwner == diagram {
						newElement := underlyingChange.GetReportingElement()
						if crldiagram.IsDiagramNode(newElement, hl) {
							additionalParameters := getNodeAdditionalParameters(newElement, hl)
							SendNotification("AddDiagramNode", newElement.GetConceptID(hl), newElement, additionalParameters)
						} else if crldiagram.IsDiagramLink(newElement, hl) {
							additionalParameters := getLinkAdditionalParameters(newElement, hl)
							SendNotification("AddDiagramLink", newElement.GetConceptID(hl), newElement, additionalParameters)
						}
					}
				case core.ChildChanged:
					// a Child of the diagram element changed. Check to see whether it is a model reference and, if so, whether the model elemment still exists
					if crldiagram.IsModelReference(underlyingChange.GetReportingElement(), hl) {
						secondUnderlyingChange := underlyingChange.GetUnderlyingChange()
						switch secondUnderlyingChange.GetNatureOfChange() {
						case core.ConceptChanged:
							// check to see whether the modelElement is nil
							modelReference := secondUnderlyingChange.GetReportingElement().(core.Reference)
							if modelReference.GetReferencedConcept(hl) == nil {
								additionalParameters := map[string]string{"OwnerID": diagram.GetConceptID(hl)}
								SendNotification("DeleteDiagramElement", reportingElement.GetConceptID(hl), reportingElement, additionalParameters)
							}
						}
					}
				}
			}
		}
	case core.ConceptChanged:
		// The diagram itself has changed. Check to see whether is has been deleted
		if diagram.GetOwningConcept(hl) == nil {
			SendNotification("DeleteDiagram", diagram.GetConceptID(hl), diagram, nil)
		}
	}
}
