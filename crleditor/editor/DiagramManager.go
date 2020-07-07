package editor

import (
	"github.com/pkg/errors"
	"strconv"

	//	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructures"
	"github.com/pbrown12303/activeCRL/crldiagram"
	"github.com/pbrown12303/activeCRL/crleditor/crleditordomain"
)

const diagramContainerSuffix = "DiagramContainer"
const diagramSuffix = "Diagram"

// diagramManager manages the diagram portion of the UI and all interactions with it
type diagramManager struct {
	crlEditor                *CrlEditor
	defaultDiagramLabelCount int
}

func newDiagramManager(crlEditor *CrlEditor) *diagramManager {
	var dm diagramManager
	dm.crlEditor = crlEditor
	uOfD := CrlEditorSingleton.GetUofD()
	addDiagramViewFunctionsToUofD(uOfD)
	return &dm
}

func (dmPtr *diagramManager) abstractPointerChanged(linkID string, sourceID string, targetID string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.crlEditor.uOfD
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.abstractPointerChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagram.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.elementPoiabstractPointerChangednterChanged called with model source not found")
	}
	var modelRefinement core.Refinement
	switch modelSource.(type) {
	case core.Refinement:
		modelRefinement = modelSource.(core.Refinement)
		break
	default:
		return "", errors.New("diagramManager.abstractPointerChanged called with source not being a Refinement")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.abstractPointerChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagram.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.abstractPointerChanged called with model target not found")
	}
	var err error
	var diagramPointer core.Element
	if linkID == "" {
		// this is a new link
		diagramPointer, err = crldiagram.NewDiagramAbstractPointer(uOfD, hl)
		if err != nil {
			return "", err
		}
		diagramPointer.SetOwningConceptID(diagramSource.GetOwningConceptID(hl), hl)
		crldiagram.SetReferencedModelElement(diagramPointer, modelSource, hl)
		crldiagram.SetLinkSource(diagramPointer, diagramSource, hl)
		crldiagram.SetLinkTarget(diagramPointer, diagramTarget, hl)
		modelRefinement.SetAbstractConcept(modelTarget, hl)
		dmPtr.crlEditor.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramPointer = uOfD.GetElement(linkID)
		if diagramPointer == nil {
			return "", errors.New("diagramManager.abstractPointerChanged called with diagramPointer not found in diagram")
		}
		if diagramSource != crldiagram.GetLinkSource(diagramPointer, hl) {
			crldiagram.SetLinkSource(diagramPointer, diagramSource, hl)
		}
		if diagramTarget != crldiagram.GetLinkTarget(diagramPointer, hl) {
			crldiagram.SetLinkTarget(diagramPointer, diagramTarget, hl)
		}
		if modelSource != crldiagram.GetReferencedModelElement(diagramPointer, hl) {
			crldiagram.SetReferencedModelElement(diagramPointer, modelSource, hl)
		}
		if modelTarget != modelRefinement.GetAbstractConcept(hl) {
			modelRefinement.SetAbstractConcept(modelTarget, hl)
		}
	}

	return diagramPointer.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) addConceptView(request *Request, hl *core.HeldLocks) (core.Element, error) {
	uOfD := dmPtr.crlEditor.GetUofD()
	diagramID := request.AdditionalParameters["DiagramID"]
	diagram := uOfD.GetElement(diagramID)
	el := uOfD.GetElement(dmPtr.crlEditor.GetTreeDragSelectionID(hl))
	if el == nil {
		return nil, errors.New("Indicated model element not found in addNodeView, ID: " + request.RequestConceptID)
	}
	var x, y float64
	var err error
	x, err = strconv.ParseFloat(request.AdditionalParameters["NodeX"], 64)
	if err != nil {
		return nil, err
	}
	y, err = strconv.ParseFloat(request.AdditionalParameters["NodeY"], 64)
	if err != nil {
		return nil, err
	}

	createAsLink := false
	switch el.(type) {
	case core.Reference:
		createAsLink = CrlEditorSingleton.GetDropDiagramReferenceAsLink(hl)
	case core.Refinement:
		createAsLink = CrlEditorSingleton.GetDropDiagramRefinementAsLink(hl)
	}

	var newElement core.Element
	if createAsLink {
		var modelSourceConcept core.Element
		var modelTargetConcept core.Element
		switch el.(type) {
		case core.Reference:
			newElement, err = crldiagram.NewDiagramReferenceLink(uOfD, hl)
			if err != nil {
				return nil, err
			}
			reference := el.(core.Reference)
			modelSourceConcept = reference.GetOwningConcept(hl)
			modelTargetConcept = reference.GetReferencedConcept(hl)
		case core.Refinement:
			newElement, err = crldiagram.NewDiagramRefinementLink(uOfD, hl)
			if err != nil {
				return nil, err
			}
			refinement := el.(core.Refinement)
			modelSourceConcept = refinement.GetRefinedConcept(hl)
			modelTargetConcept = refinement.GetAbstractConcept(hl)
		}
		if modelSourceConcept == nil {
			return nil, errors.New("In addConceptView for link, modelSourceConcept is nil")
		}
		if modelTargetConcept == nil {
			return nil, errors.New("In addConceptView for link, modelTargetConcept is nil")
		}
		diagramSourceElement := crldiagram.GetFirstElementRepresentingConcept(diagram, modelSourceConcept, hl)
		if diagramSourceElement == nil {
			return nil, errors.New("In addConceptView for reference link, diagramSourceElement is nil")
		}
		diagramTargetElement := crldiagram.GetFirstElementRepresentingConcept(diagram, modelTargetConcept, hl)
		if diagramTargetElement == nil {
			return nil, errors.New("In addConceptView for reference link, diagramTargetElement is nil")
		}
		crldiagram.SetLinkSource(newElement, diagramSourceElement, hl)
		crldiagram.SetLinkTarget(newElement, diagramTargetElement, hl)
	} else {
		newElement, err = crldiagram.NewDiagramNode(uOfD, hl)
		if err != nil {
			return nil, err
		}
		crldiagram.SetNodeX(newElement, x, hl)
		crldiagram.SetNodeY(newElement, y, hl)
	}

	newElement.SetLabel(el.GetLabel(hl), hl)
	crldiagram.SetReferencedModelElement(newElement, el, hl)
	crldiagram.SetDisplayLabel(newElement, el.GetLabel(hl), hl)

	newElement.SetOwningConceptID(diagramID, hl)
	hl.ReleaseLocksAndWait()

	return newElement, nil
}

// addDiagramToDisplayedList adds the diagram to the list of displayed diagrams
func (dmPtr *diagramManager) addDiagramToDisplayedList(diagram core.Element, hl *core.HeldLocks) error {
	if !dmPtr.isDiagramDisplayed(diagram, hl) {
		openDiagrams := dmPtr.crlEditor.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
		_, err := crldatastructures.AppendListMember(openDiagrams, diagram, hl)
		if err != nil {
			return errors.Wrap(err, "diagramManager.addDiagramToDisplayedList failed")
		}
	}
	return nil
}

func (dmPtr *diagramManager) closeDiagramView(diagramID string, hl *core.HeldLocks) error {
	diagram := dmPtr.crlEditor.uOfD.GetElement(diagramID)
	if diagram == nil {
		return errors.New("In diagram manager.closeDiagramView, the suppied diagram was not found in the uOfD")
	}
	if !diagram.IsRefinementOfURI(crldiagram.CrlDiagramURI, hl) {
		return errors.New("In diagramManager.closeDiagramView, the supplied diagram is not a refinement of CrlDiagramURI")
	}
	// If the diagram is in the list of displayed diagrams, remove it
	if dmPtr.isDiagramDisplayed(diagram, hl) {
		dmPtr.removeDiagramFromDisplayedList(diagram, hl)
	}
	_, err := SendNotification("CloseDiagramView", diagramID, nil, map[string]string{})
	return err
}

func (dmPtr *diagramManager) deleteDiagramElementView(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.crlEditor.uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.deleteDiagramElementView diagramElement not found for elementID " + elementID)
	}
	dEls := mapset.NewSet(diagramElement.GetConceptID(hl))
	return dmPtr.crlEditor.uOfD.DeleteElements(dEls, hl)
}

// diagramChanged handles the update of the diagram tab when the diagram name changes
func (dmPtr *diagramManager) diagramLabelChanged(diagramID string, diagram core.Element, hl *core.HeldLocks) {
	SendNotification("DiagramLabelChanged", diagramID, diagram, map[string]string{})
}

// diagramClick handles the creation of a new Element and adding a node representation of it to the diagram
func (dmPtr *diagramManager) diagramClick(request *Request, hl *core.HeldLocks) error {
	uOfD := dmPtr.crlEditor.uOfD
	diagramID := request.AdditionalParameters["DiagramID"]
	diagram := uOfD.GetElement(diagramID)
	if diagram == nil {
		return errors.New("diagramManager.diagramClick called with DiagramID that does not exist")
	}
	var el core.Element
	var err error
	nodeType := request.AdditionalParameters["NodeType"]
	if nodeType == "" {
		return errors.New("diagramManager.diagramClick called with no NodeType")
	}
	switch nodeType {
	case "Element":
		el, err = uOfD.NewElement(hl)
		el.SetLabel(dmPtr.crlEditor.getDefaultElementLabel(), hl)
	case "Literal":
		el, err = uOfD.NewLiteral(hl)
		el.SetLabel(dmPtr.crlEditor.getDefaultLiteralLabel(), hl)
	case "Reference":
		el, err = uOfD.NewReference(hl)
		el.SetLabel(dmPtr.crlEditor.getDefaultReferenceLabel(), hl)
	case "Refinement":
		el, err = uOfD.NewRefinement(hl)
		el.SetLabel(dmPtr.crlEditor.getDefaultRefinementLabel(), hl)
	case "Diagram":
		el, err = crldiagram.NewDiagram(uOfD, hl)
		el.SetLabel(dmPtr.getDefaultDiagramLabel(), hl)
	}
	if err != nil {
		return err
	}
	el.SetOwningConceptID(diagram.GetOwningConceptID(hl), hl)
	dmPtr.crlEditor.SelectElement(el, hl)

	// Now the view
	var x, y float64
	x, err = strconv.ParseFloat(request.AdditionalParameters["NodeX"], 64)
	if err != nil {
		return err
	}
	y, err = strconv.ParseFloat(request.AdditionalParameters["NodeY"], 64)
	if err != nil {
		return err
	}
	var newNode core.Element
	newNode, err = crldiagram.NewDiagramNode(uOfD, hl)
	if err != nil {
		return err
	}
	crldiagram.SetNodeX(newNode, x, hl)
	crldiagram.SetNodeY(newNode, y, hl)
	newNode.SetLabel(el.GetLabel(hl), hl)
	crldiagram.SetReferencedModelElement(newNode, el, hl)
	crldiagram.SetDisplayLabel(newNode, el.GetLabel(hl), hl)

	newNode.SetOwningConceptID(diagramID, hl)
	hl.ReleaseLocksAndWait()
	dmPtr.crlEditor.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})

	return nil
}

// diagramDrop evaluates the request resulting from a drop in the diagram
func (dmPtr *diagramManager) diagramDrop(request *Request, hl *core.HeldLocks) error {
	_, err := dmPtr.addConceptView(request, hl)
	if err != nil {
		return err
	}
	dmPtr.crlEditor.SetTreeDragSelection("")
	return nil
}

// DiagramViewHasBeenClosed notifies the server that the client has closed the diagram view
func (dmPtr *diagramManager) DiagramViewHasBeenClosed(diagramID string, hl *core.HeldLocks) error {
	diagram := dmPtr.crlEditor.uOfD.GetElement(diagramID)
	if !diagram.IsRefinementOfURI(crldiagram.CrlDiagramURI, hl) {
		return errors.New("In diagramManager.closeDiagramView, the supplied diagram is not a refinement of CrlDiagramURI")
	}
	// If the diagram is in the list of displayed diagrams, remove it
	if dmPtr.isDiagramDisplayed(diagram, hl) {
		dmPtr.removeDiagramFromDisplayedList(diagram, hl)
	}
	return nil
}

// displayDiagram tells the client to display the indicated diagram.
func (dmPtr *diagramManager) displayDiagram(diagram core.Element, hl *core.HeldLocks) error {
	if !diagram.IsRefinementOfURI(crldiagram.CrlDiagramURI, hl) {
		return errors.New("In diagramManager.displayDiagram, the supplied diagram is not a refinement of CrlDiagramURI")
	}
	// Make sure the diagram is in the list of displayed diagrams
	if !dmPtr.isDiagramDisplayed(diagram, hl) {
		err3 := dmPtr.addDiagramToDisplayedList(diagram, hl)
		if err3 != nil {
			return errors.Wrap(err3, "diagramManager.displayDiagram failed")
		}
	}
	// make sure there is a monitor on the diagram so we know when it has been deleted
	err2 := dmPtr.verifyMonitorPresent(diagram, hl)
	if err2 != nil {
		return err2
	}
	// Tell the client to display the diagram
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("DisplayDiagram", diagram.GetConceptID(hl), diagram, nil)
	if err != nil {
		return err
	}
	if notificationResponse == nil {
		return errors.New("In diagramManager.displayDiagram the notification response was nil")
	}
	if notificationResponse.Result != 0 {
		return err
	}
	return dmPtr.refreshDiagram(diagram, hl)
}

func (dmPtr *diagramManager) getDefaultDiagramLabel() string {
	dmPtr.defaultDiagramLabelCount++
	countString := strconv.Itoa(dmPtr.defaultDiagramLabelCount)
	return "Diagram" + countString
}

func (dmPtr *diagramManager) elementPointerChanged(linkID string, sourceID string, targetID string, targetAttributeName string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.crlEditor.uOfD
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.elementPointerChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagram.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.elementPointerChanged called with model source not found")
	}
	var modelReference core.Reference
	switch modelSource.(type) {
	case core.Reference:
		modelReference = modelSource.(core.Reference)
		break
	default:
		return "", errors.New("diagramManager.elementPointerChanged called with source not being a Reference")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.elementPointerChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagram.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.elementPointerChanged called with model target not found")
	}
	var err error
	var diagramPointer core.Element
	attributeName := core.NoAttribute
	switch targetAttributeName {
	case "OwningConceptID":
		attributeName = core.OwningConceptID
	case "ReferencedConceptID":
		attributeName = core.ReferencedConceptID
	case "AbstractConceptID":
		attributeName = core.AbstractConceptID
	case "RefinedConceptID":
		attributeName = core.RefinedConceptID
	}
	modelReference.SetReferencedAttributeName(attributeName, hl)
	if linkID == "" {
		// this is a new link
		diagramPointer, err = crldiagram.NewDiagramElementPointer(uOfD, hl)
		if err != nil {
			return "", err
		}
		diagramPointer.SetOwningConceptID(diagramSource.GetOwningConceptID(hl), hl)
		crldiagram.SetReferencedModelElement(diagramPointer, modelSource, hl)
		crldiagram.SetLinkSource(diagramPointer, diagramSource, hl)
		crldiagram.SetLinkTarget(diagramPointer, diagramTarget, hl)
		modelReference.SetReferencedConcept(modelTarget, hl)
		dmPtr.crlEditor.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramPointer = uOfD.GetElement(linkID)
		if diagramPointer == nil {
			return "", errors.New("diagramManager.elementPointerChanged called with diagramPointer not found in diagram")
		}
		if diagramSource != crldiagram.GetLinkSource(diagramPointer, hl) {
			crldiagram.SetLinkSource(diagramPointer, diagramSource, hl)
		}
		if diagramTarget != crldiagram.GetLinkTarget(diagramPointer, hl) {
			crldiagram.SetLinkTarget(diagramPointer, diagramTarget, hl)
		}
		if modelSource != crldiagram.GetReferencedModelElement(diagramPointer, hl) {
			crldiagram.SetReferencedModelElement(diagramPointer, modelSource, hl)
		}
		if modelTarget != modelReference.GetReferencedConcept(hl) {
			modelReference.SetReferencedConcept(modelTarget, hl)
		}
	}

	return diagramPointer.GetConceptID(hl), nil
}

// isDiagramDisplayed returns true if the diagram is in the list of displayed diagrams
func (dmPtr *diagramManager) isDiagramDisplayed(diagram core.Element, hl *core.HeldLocks) bool {
	openDiagrams := dmPtr.crlEditor.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
	return crldatastructures.IsListMember(openDiagrams, diagram, hl)
}

// newDiagram creates a new crldiagram
func (dmPtr *diagramManager) newDiagram(hl *core.HeldLocks) core.Element {
	// Insert name prompt here
	name := dmPtr.getDefaultDiagramLabel()
	uOfD := CrlEditorSingleton.GetUofD()
	diagram, err := crldiagram.NewDiagram(uOfD, hl)
	if err != nil {
		log.Print(err)
	}
	diagram.SetLabel(name, hl)
	hl.ReleaseLocksAndWait()
	return diagram
}

func (dmPtr *diagramManager) ownerPointerChanged(linkID string, sourceID string, targetID string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.crlEditor.uOfD
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.ownerPointerChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagram.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.ownerPointerChanged called with model source not found")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.ownerPointerChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagram.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.ownerPointerChanged called with model target not found")
	}
	var err error
	var diagramPointer core.Element
	if linkID == "" {
		// this is a new link
		diagramPointer, err = crldiagram.NewDiagramOwnerPointer(uOfD, hl)
		if err != nil {
			return "", err
		}
		diagramPointer.SetOwningConceptID(diagramSource.GetOwningConceptID(hl), hl)
		crldiagram.SetReferencedModelElement(diagramPointer, modelSource, hl)
		crldiagram.SetLinkSource(diagramPointer, diagramSource, hl)
		crldiagram.SetLinkTarget(diagramPointer, diagramTarget, hl)
		modelSource.SetOwningConcept(modelTarget, hl)
		dmPtr.crlEditor.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramPointer = uOfD.GetElement(linkID)
		if diagramPointer == nil {
			return "", errors.New("diagramManager.ownerPointerChanged called with diagramPointer not found in diagram")
		}
		if diagramSource != crldiagram.GetLinkSource(diagramPointer, hl) {
			crldiagram.SetLinkSource(diagramPointer, diagramSource, hl)
		}
		if diagramTarget != crldiagram.GetLinkTarget(diagramPointer, hl) {
			crldiagram.SetLinkTarget(diagramPointer, diagramTarget, hl)
		}
		if modelSource != crldiagram.GetReferencedModelElement(diagramPointer, hl) {
			crldiagram.SetReferencedModelElement(diagramPointer, modelSource, hl)
		}
		if modelTarget != modelSource.GetOwningConcept(hl) {
			modelSource.SetOwningConcept(modelTarget, hl)
		}
	}

	return diagramPointer.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) refinedPointerChanged(linkID string, sourceID string, targetID string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.crlEditor.uOfD
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.refinedPointerChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagram.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.elementPoirefinedPointerChangednterChanged called with model source not found")
	}
	var modelRefinement core.Refinement
	switch modelSource.(type) {
	case core.Refinement:
		modelRefinement = modelSource.(core.Refinement)
		break
	default:
		return "", errors.New("diagramManager.refinedPointerChanged called with source not being a Refinement")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.refinedPointerChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagram.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.refinedPointerChanged called with model target not found")
	}
	var err error
	var diagramPointer core.Element
	if linkID == "" {
		// this is a new link
		diagramPointer, err = crldiagram.NewDiagramRefinedPointer(uOfD, hl)
		if err != nil {
			return "", err
		}
		diagramPointer.SetOwningConceptID(diagramSource.GetOwningConceptID(hl), hl)
		crldiagram.SetReferencedModelElement(diagramPointer, modelSource, hl)
		crldiagram.SetLinkSource(diagramPointer, diagramSource, hl)
		crldiagram.SetLinkTarget(diagramPointer, diagramTarget, hl)
		modelRefinement.SetRefinedConcept(modelTarget, hl)
		dmPtr.crlEditor.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramPointer = uOfD.GetElement(linkID)
		if diagramPointer == nil {
			return "", errors.New("diagramManager.refinedPointerChanged called with diagramPointer not found in diagram")
		}
		if diagramSource != crldiagram.GetLinkSource(diagramPointer, hl) {
			crldiagram.SetLinkSource(diagramPointer, diagramSource, hl)
		}
		if diagramTarget != crldiagram.GetLinkTarget(diagramPointer, hl) {
			crldiagram.SetLinkTarget(diagramPointer, diagramTarget, hl)
		}
		if modelSource != crldiagram.GetReferencedModelElement(diagramPointer, hl) {
			crldiagram.SetReferencedModelElement(diagramPointer, modelSource, hl)
		}
		if modelTarget != modelRefinement.GetRefinedConcept(hl) {
			modelRefinement.SetRefinedConcept(modelTarget, hl)
		}
	}

	return diagramPointer.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) ReferenceLinkChanged(linkID string, sourceID string, targetID string, targetAttributeName string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.crlEditor.uOfD
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagram.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with model source not found")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagram.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with model target not found")
	}
	diagram := diagramSource.GetOwningConcept(hl)
	var err error
	var diagramLink core.Element
	var modelElement core.Element
	attributeName := core.NoAttribute
	switch targetAttributeName {
	case "OwningConceptID":
		attributeName = core.OwningConceptID
	case "ReferencedConceptID":
		attributeName = core.ReferencedConceptID
	case "AbstractConceptID":
		attributeName = core.AbstractConceptID
	case "RefinedConceptID":
		attributeName = core.RefinedConceptID
	}
	if linkID == "" {
		// this is a new reference
		newReference, _ := uOfD.NewReference(hl)
		newReference.SetReferencedConcept(modelTarget, hl)
		newReference.SetOwningConcept(modelSource, hl)
		newReference.SetReferencedAttributeName(attributeName, hl)
		diagramLink, err = crldiagram.NewDiagramReferenceLink(uOfD, hl)
		if err != nil {
			return "", err
		}
		diagramLink.SetOwningConcept(diagram, hl)
		crldiagram.SetReferencedModelElement(diagramLink, newReference, hl)
		crldiagram.SetLinkSource(diagramLink, diagramSource, hl)
		crldiagram.SetLinkTarget(diagramLink, diagramTarget, hl)
		modelElement = newReference
		dmPtr.crlEditor.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramLink = uOfD.GetElement(linkID)
		modelElement = crldiagram.GetReferencedModelElement(diagramLink, hl)
		if modelElement != nil {
			switch modelElement.(type) {
			case core.Reference:
				reference := modelElement.(core.Reference)
				if diagramLink == nil {
					return "", errors.New("diagramManager.refinementLinkChanged called with diagramPointer not found in diagram")
				}
				if diagramSource != crldiagram.GetLinkSource(diagramLink, hl) {
					reference.SetOwningConcept(modelSource, hl)
				}
				if diagramTarget != crldiagram.GetLinkTarget(diagramLink, hl) {
					reference.SetReferencedConcept(modelTarget, hl)
				}
				reference.SetReferencedAttributeName(attributeName, hl)
			}
		}
	}
	dmPtr.crlEditor.SelectElement(modelElement, hl)

	return diagramLink.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) RefinementLinkChanged(linkID string, sourceID string, targetID string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.crlEditor.uOfD
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagram.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with model source not found")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagram.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with model target not found")
	}
	diagram := diagramSource.GetOwningConcept(hl)
	diagramOwner := diagram.GetOwningConcept(hl)
	var err error
	var diagramLink core.Element
	var modelElement core.Element
	if linkID == "" {
		// this is a new refinement
		newRefinement, _ := uOfD.NewRefinement(hl)
		newRefinement.SetRefinedConcept(modelSource, hl)
		newRefinement.SetAbstractConcept(modelTarget, hl)
		newRefinement.SetOwningConcept(diagramOwner, hl)
		diagramLink, err = crldiagram.NewDiagramRefinementLink(uOfD, hl)
		if err != nil {
			return "", err
		}
		diagramLink.SetOwningConcept(diagram, hl)
		crldiagram.SetReferencedModelElement(diagramLink, newRefinement, hl)
		crldiagram.SetLinkSource(diagramLink, diagramSource, hl)
		crldiagram.SetLinkTarget(diagramLink, diagramTarget, hl)
		modelElement = newRefinement
		dmPtr.crlEditor.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramLink = uOfD.GetElement(linkID)
		modelElement = crldiagram.GetReferencedModelElement(diagramLink, hl)
		if modelElement != nil {
			switch modelElement.(type) {
			case core.Refinement:
				refinement := modelElement.(core.Refinement)
				if diagramLink == nil {
					return "", errors.New("diagramManager.refinementLinkChanged called with diagramPointer not found in diagram")
				}
				if diagramSource != crldiagram.GetLinkSource(diagramLink, hl) {
					refinement.SetRefinedConcept(modelSource, hl)
				}
				if diagramTarget != crldiagram.GetLinkTarget(diagramLink, hl) {
					refinement.SetAbstractConcept(modelTarget, hl)
				}
			}
		}
	}
	dmPtr.crlEditor.SelectElement(modelElement, hl)

	return diagramLink.GetConceptID(hl), nil
}

// refreshDiagramUsingURI finds the diagram and resends all diagram elements to the browser
func (dmPtr *diagramManager) refreshDiagramUsingURI(diagramID string, hl *core.HeldLocks) error {
	diagram := dmPtr.crlEditor.uOfD.GetElement(diagramID)
	if diagram == nil {
		return errors.New("In diagramManager.refreshDiagram, diagram not found for ID: " + diagramID)
	}
	return dmPtr.refreshDiagram(diagram, hl)
}

// refreshDiagram resends all diagram elements to the browser
func (dmPtr *diagramManager) refreshDiagram(diagram core.Element, hl *core.HeldLocks) error {
	nodes := diagram.GetOwnedConceptsRefinedFromURI(crldiagram.CrlDiagramNodeURI, hl)
	for _, node := range nodes {
		additionalParameters := getNodeAdditionalParameters(node, hl)
		notificationResponse, err := CrlEditorSingleton.SendNotification("AddDiagramNode", node.GetConceptID(hl), node, additionalParameters)
		if err != nil {
			return err
		}
		if notificationResponse.Result != 0 {
			return errors.New(notificationResponse.ErrorMessage)
		}
	}
	links := diagram.GetOwnedConceptsRefinedFromURI(crldiagram.CrlDiagramLinkURI, hl)
	for _, link := range links {
		additionalParameters := getLinkAdditionalParameters(link, hl)
		notificationResponse, err := CrlEditorSingleton.SendNotification("AddDiagramLink", link.GetConceptID(hl), link, additionalParameters)
		if err != nil {
			return err
		}
		if notificationResponse.Result != 0 {
			return errors.New(notificationResponse.ErrorMessage)
		}
	}
	return nil
}

// removeDiagramFromDisplayedList adds the diagram to the list of displayed diagrams
func (dmPtr *diagramManager) removeDiagramFromDisplayedList(diagram core.Element, hl *core.HeldLocks) {
	if dmPtr.isDiagramDisplayed(diagram, hl) {
		openDiagrams := dmPtr.crlEditor.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
		crldatastructures.RemoveListMember(openDiagrams, diagram, hl)
	}
}

// ResetDefaultLabelCounts
func (dmPtr *diagramManager) ResetDefaultLabelCounts() {
	dmPtr.defaultDiagramLabelCount = 0
}

// setDiagramNodePosition sets the position of the diagram node
func (dmPtr *diagramManager) setDiagramNodePosition(nodeID string, x float64, y float64, hl *core.HeldLocks) {
	node := CrlEditorSingleton.GetUofD().GetElement(nodeID)
	if node == nil {
		// This can happen when the concept space containing the diagram is deleted???
		log.Print("In setDiagramNodePosition node not found for nodeID: " + nodeID)
		return
	}
	crldiagram.SetNodeX(node, x, hl)
	crldiagram.SetNodeY(node, y, hl)
}

func (dmPtr *diagramManager) showAbstractConcept(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.crlEditor.uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showAbstractConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showAbstractConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagram.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showAbstractConcept modelConcept not found for elementID " + elementID)
	}
	var modelRefinement core.Refinement
	switch modelConcept.(type) {
	case core.Refinement:
		modelRefinement = modelConcept.(core.Refinement)
		break
	default:
		return errors.New("diagramManager.showAbstractConcept modelConcept is not a Refinement")
	}
	modelAbstractConcept := modelRefinement.GetAbstractConcept(hl)
	if modelAbstractConcept == nil {
		return errors.New("Abstract Concept is nil")
	}
	diagramAbstractConcept := crldiagram.GetFirstElementRepresentingConcept(diagram, modelAbstractConcept, hl)
	if diagramAbstractConcept == nil {
		diagramAbstractConcept, _ = crldiagram.NewDiagramNode(dmPtr.crlEditor.uOfD, hl)
		crldiagram.SetReferencedModelElement(diagramAbstractConcept, modelAbstractConcept, hl)
		diagramAbstractConcept.SetOwningConcept(diagram, hl)
		diagramElementX := crldiagram.GetNodeX(diagramElement, hl)
		diagramElementY := crldiagram.GetNodeY(diagramElement, hl)
		crldiagram.SetNodeX(diagramAbstractConcept, diagramElementX, hl)
		crldiagram.SetNodeY(diagramAbstractConcept, diagramElementY-100, hl)
	}
	elementPointer := crldiagram.GetElementPointer(diagram, diagramElement, hl)
	if elementPointer == nil {
		elementPointer, _ = crldiagram.NewDiagramAbstractPointer(dmPtr.crlEditor.uOfD, hl)
		elementPointer.SetOwningConcept(diagram, hl)
		crldiagram.SetReferencedModelElement(elementPointer, modelConcept, hl)
		crldiagram.SetLinkSource(elementPointer, diagramElement, hl)
		crldiagram.SetLinkTarget(elementPointer, diagramAbstractConcept, hl)
	}
	return nil
}

func (dmPtr *diagramManager) showOwner(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.crlEditor.uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showOwner diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showOwner diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagram.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showOwner modelConcept not found for elementID " + elementID)
	}
	modelConceptOwner := modelConcept.GetOwningConcept(hl)
	if modelConceptOwner == nil {
		return errors.New("Owner is nil")
	}
	diagramConceptOwner := crldiagram.GetFirstElementRepresentingConcept(diagram, modelConceptOwner, hl)
	if diagramConceptOwner == nil {
		diagramConceptOwner, _ = crldiagram.NewDiagramNode(dmPtr.crlEditor.uOfD, hl)
		crldiagram.SetReferencedModelElement(diagramConceptOwner, modelConceptOwner, hl)
		diagramConceptOwner.SetOwningConcept(diagram, hl)
		diagramElementX := crldiagram.GetNodeX(diagramElement, hl)
		diagramElementY := crldiagram.GetNodeY(diagramElement, hl)
		crldiagram.SetNodeX(diagramConceptOwner, diagramElementX, hl)
		crldiagram.SetNodeY(diagramConceptOwner, diagramElementY-100, hl)
	}
	ownerPointer := crldiagram.GetOwnerPointer(diagram, diagramElement, hl)
	if ownerPointer == nil {
		ownerPointer, _ = crldiagram.NewDiagramOwnerPointer(dmPtr.crlEditor.uOfD, hl)
		ownerPointer.SetOwningConcept(diagram, hl)
		crldiagram.SetReferencedModelElement(ownerPointer, modelConcept, hl)
		crldiagram.SetLinkSource(ownerPointer, diagramElement, hl)
		crldiagram.SetLinkTarget(ownerPointer, diagramConceptOwner, hl)
	}
	return nil
}

func (dmPtr *diagramManager) showReferencedConcept(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.crlEditor.uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showReferencedConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showReferencedConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagram.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showReferencedConcept modelConcept not found for elementID " + elementID)
	}
	var modelReference core.Reference
	switch modelConcept.(type) {
	case core.Reference:
		modelReference = modelConcept.(core.Reference)
		break
	default:
		return errors.New("diagramManager.showReferencedConcept modelConcept is not a Reference")
	}
	modelReferencedConcept := modelReference.GetReferencedConcept(hl)
	if modelReferencedConcept == nil {
		return errors.New("Referenced Concept is nil")
	}
	diagramReferencedConcept := crldiagram.GetFirstElementRepresentingConcept(diagram, modelReferencedConcept, hl)
	if diagramReferencedConcept == nil {
		diagramReferencedConcept, _ = crldiagram.NewDiagramNode(dmPtr.crlEditor.uOfD, hl)
		crldiagram.SetReferencedModelElement(diagramReferencedConcept, modelReferencedConcept, hl)
		diagramReferencedConcept.SetOwningConcept(diagram, hl)
		diagramElementX := crldiagram.GetNodeX(diagramElement, hl)
		diagramElementY := crldiagram.GetNodeY(diagramElement, hl)
		crldiagram.SetNodeX(diagramReferencedConcept, diagramElementX, hl)
		crldiagram.SetNodeY(diagramReferencedConcept, diagramElementY-100, hl)
	}
	elementPointer := crldiagram.GetElementPointer(diagram, diagramElement, hl)
	if elementPointer == nil {
		elementPointer, _ = crldiagram.NewDiagramElementPointer(dmPtr.crlEditor.uOfD, hl)
		elementPointer.SetOwningConcept(diagram, hl)
		crldiagram.SetReferencedModelElement(elementPointer, modelConcept, hl)
		crldiagram.SetLinkSource(elementPointer, diagramElement, hl)
		crldiagram.SetLinkTarget(elementPointer, diagramReferencedConcept, hl)
	}
	return nil
}

func (dmPtr *diagramManager) showRefinedConcept(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.crlEditor.uOfD.GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showRefinedConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showRefinedConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagram.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showRefinedConcept modelConcept not found for elementID " + elementID)
	}
	var modelRefinement core.Refinement
	switch modelConcept.(type) {
	case core.Refinement:
		modelRefinement = modelConcept.(core.Refinement)
		break
	default:
		return errors.New("diagramManager.showRefinedConcept modelConcept is not a Refinement")
	}
	modelRefinedConcept := modelRefinement.GetRefinedConcept(hl)
	if modelRefinedConcept == nil {
		return errors.New("Refined Concept is nil")
	}
	diagramRefinedConcept := crldiagram.GetFirstElementRepresentingConcept(diagram, modelRefinedConcept, hl)
	if diagramRefinedConcept == nil {
		diagramRefinedConcept, _ = crldiagram.NewDiagramNode(dmPtr.crlEditor.uOfD, hl)
		crldiagram.SetReferencedModelElement(diagramRefinedConcept, modelRefinedConcept, hl)
		diagramRefinedConcept.SetOwningConcept(diagram, hl)
		diagramElementX := crldiagram.GetNodeX(diagramElement, hl)
		diagramElementY := crldiagram.GetNodeY(diagramElement, hl)
		crldiagram.SetNodeX(diagramRefinedConcept, diagramElementX, hl)
		crldiagram.SetNodeY(diagramRefinedConcept, diagramElementY-100, hl)
	}
	elementPointer := crldiagram.GetElementPointer(diagram, diagramElement, hl)
	if elementPointer == nil {
		elementPointer, _ = crldiagram.NewDiagramRefinedPointer(dmPtr.crlEditor.uOfD, hl)
		elementPointer.SetOwningConcept(diagram, hl)
		crldiagram.SetReferencedModelElement(elementPointer, modelConcept, hl)
		crldiagram.SetLinkSource(elementPointer, diagramElement, hl)
		crldiagram.SetLinkTarget(elementPointer, diagramRefinedConcept, hl)
	}
	return nil
}

func (dmPtr *diagramManager) verifyMonitorPresent(diagram core.Element, hl *core.HeldLocks) error {
	workingConceptSpace := dmPtr.crlEditor.workingConceptSpace
	for _, monitor := range workingConceptSpace.GetOwnedReferencesRefinedFromURI(crleditordomain.DiagramViewMonitorURI, hl) {
		if monitor.GetReferencedConcept(hl) == diagram {
			return nil
		}
	}
	newMonitor, err := dmPtr.crlEditor.uOfD.CreateReplicateAsRefinementFromURI(crleditordomain.DiagramViewMonitorURI, hl)
	if err != nil {
		return err
	}
	newMonitor.SetOwningConcept(workingConceptSpace, hl)
	newMonitor.(core.Reference).SetReferencedConcept(diagram, hl)
	return nil
}
