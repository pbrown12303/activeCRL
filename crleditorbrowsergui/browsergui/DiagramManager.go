package browsergui

import (
	"strconv"

	"github.com/pkg/errors"

	//	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	// "github.com/pbrown12303/activeCRL/crleditorbrowserguidomain"
)

// diagramManager manages the diagram portion of the UI and all interactions with it
type diagramManager struct {
	browserGUI     *BrowserGUI
	diagrams       map[string]core.Element
	elementManager *diagramElementManager
}

func newDiagramManager(browserGUI *BrowserGUI) *diagramManager {
	dm := &diagramManager{}
	dm.browserGUI = browserGUI
	dm.diagrams = map[string]core.Element{}
	dm.elementManager = newDiagramElementManager(dm)
	return dm
}

func (dmPtr *diagramManager) abstractPointerChanged(linkID string, sourceID string, targetID string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.browserGUI.GetUofD()
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.abstractPointerChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagramdomain.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.elementPoiabstractPointerChangednterChanged called with model source not found")
	}
	var modelRefinement core.Refinement
	switch typedModelSource := modelSource.(type) {
	case core.Refinement:
		modelRefinement = typedModelSource
	default:
		return "", errors.New("diagramManager.abstractPointerChanged called with source not being a Refinement")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.abstractPointerChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagramdomain.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.abstractPointerChanged called with model target not found")
	}
	var err error
	var diagramPointer core.Element
	if linkID == "" {
		// this is a new link
		diagramPointer, err = crldiagramdomain.NewDiagramAbstractPointer(uOfD, hl)
		if err != nil {
			return "", err
		}
		crldiagramdomain.SetReferencedModelElement(diagramPointer, modelSource, hl)
		crldiagramdomain.SetLinkSource(diagramPointer, diagramSource, hl)
		crldiagramdomain.SetLinkTarget(diagramPointer, diagramTarget, hl)
		modelRefinement.SetAbstractConcept(modelTarget, hl)
		diagramPointer.SetOwningConceptID(diagramSource.GetOwningConceptID(hl), hl)
		dmPtr.browserGUI.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramPointer = uOfD.GetElement(linkID)
		if diagramPointer == nil {
			return "", errors.New("diagramManager.abstractPointerChanged called with diagramPointer not found in diagram")
		}
		if diagramSource != crldiagramdomain.GetLinkSource(diagramPointer, hl) {
			crldiagramdomain.SetLinkSource(diagramPointer, diagramSource, hl)
		}
		if diagramTarget != crldiagramdomain.GetLinkTarget(diagramPointer, hl) {
			crldiagramdomain.SetLinkTarget(diagramPointer, diagramTarget, hl)
		}
		if modelSource != crldiagramdomain.GetReferencedModelElement(diagramPointer, hl) {
			crldiagramdomain.SetReferencedModelElement(diagramPointer, modelSource, hl)
		}
		if modelTarget != modelRefinement.GetAbstractConcept(hl) {
			modelRefinement.SetAbstractConcept(modelTarget, hl)
		}
	}

	return diagramPointer.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) addConceptView(request *Request, hl *core.HeldLocks) (core.Element, error) {
	uOfD := dmPtr.browserGUI.GetUofD()
	diagramID := request.AdditionalParameters["DiagramID"]
	diagram := uOfD.GetElement(diagramID)
	el := uOfD.GetElement(dmPtr.browserGUI.GetTreeDragSelectionID(hl))
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
	return dmPtr.addConceptViewImpl(uOfD, diagram, el, x, y, hl)
}

// addConceptViewImpl creates the concept view and adds it to the diagram
func (dmPtr *diagramManager) addConceptViewImpl(uOfD *core.UniverseOfDiscourse, diagram core.Element, el core.Element, x float64, y float64, hl *core.HeldLocks) (core.Element, error) {

	createAsLink := false
	switch el.(type) {
	case core.Reference:
		createAsLink = dmPtr.browserGUI.editor.GetDropDiagramReferenceAsLink(hl)
	case core.Refinement:
		createAsLink = dmPtr.browserGUI.editor.GetDropDiagramRefinementAsLink(hl)
	}

	var newElement core.Element
	var err error
	if createAsLink {
		var modelSourceConcept core.Element
		var modelTargetConcept core.Element
		switch elTyped := el.(type) {
		case core.Reference:
			newElement, err = crldiagramdomain.NewDiagramReferenceLink(uOfD, hl)
			if err != nil {
				return nil, err
			}
			reference := elTyped
			modelSourceConcept = reference.GetOwningConcept(hl)
			modelTargetConcept = reference.GetReferencedConcept(hl)
		case core.Refinement:
			newElement, err = crldiagramdomain.NewDiagramRefinementLink(uOfD, hl)
			if err != nil {
				return nil, err
			}
			refinement := elTyped
			modelSourceConcept = refinement.GetRefinedConcept(hl)
			modelTargetConcept = refinement.GetAbstractConcept(hl)
		}
		if modelSourceConcept == nil {
			return nil, errors.New("In addConceptView for link, modelSourceConcept is nil")
		}
		if modelTargetConcept == nil {
			return nil, errors.New("In addConceptView for link, modelTargetConcept is nil")
		}
		diagramSourceElement := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelSourceConcept, hl)
		if diagramSourceElement == nil {
			return nil, errors.New("In addConceptView for reference link, diagramSourceElement is nil")
		}
		diagramTargetElement := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelTargetConcept, hl)
		if diagramTargetElement == nil {
			return nil, errors.New("In addConceptView for reference link, diagramTargetElement is nil")
		}
		crldiagramdomain.SetLinkSource(newElement, diagramSourceElement, hl)
		crldiagramdomain.SetLinkTarget(newElement, diagramTargetElement, hl)
	} else {
		newElement, err = crldiagramdomain.NewDiagramNode(uOfD, hl)
		if err != nil {
			return nil, err
		}
		crldiagramdomain.SetNodeX(newElement, x, hl)
		crldiagramdomain.SetNodeY(newElement, y, hl)
		crldiagramdomain.SetLineColor(newElement, "#000000", hl)
	}

	err = newElement.SetLabel(el.GetLabel(hl), hl)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.addConceptView failed")
	}
	crldiagramdomain.SetReferencedModelElement(newElement, el, hl)
	crldiagramdomain.SetDisplayLabel(newElement, el.GetLabel(hl), hl)

	err = newElement.SetOwningConceptID(diagram.GetConceptID(hl), hl)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.addConceptView failed")
	}
	err = newElement.Register(dmPtr.elementManager)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.addConceptView failed")
	}
	hl.ReleaseLocksAndWait()

	return newElement, nil
}

// addCopyWithRefinement creates a copy with refinement of the selected item and places it in the diagram.
func (dmPtr *diagramManager) addCopyWithRefinement(request *Request, hl *core.HeldLocks) (core.Element, error) {
	uOfD := dmPtr.browserGUI.GetUofD()
	diagramID := request.AdditionalParameters["DiagramID"]
	diagram := uOfD.GetElement(diagramID)
	if diagram == nil {
		return nil, errors.New("Diagram not found in diagramManager.addNodeView, ID: " + request.RequestConceptID)
	}
	el := uOfD.GetElement(dmPtr.browserGUI.GetTreeDragSelectionID(hl))
	if el == nil {
		return nil, errors.New("Indicated model element not found in diagramManager.addNodeView, ID: " + request.RequestConceptID)
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
	copy, err := uOfD.CreateReplicateAsRefinement(el, hl)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.addCopyWithRefinement failed")
	}
	copy.SetOwningConcept(diagram.GetOwningConcept(hl), hl)
	hl.ReleaseLocksAndWait()
	_, err2 := dmPtr.addConceptViewImpl(uOfD, diagram, copy, x, y, hl)
	if err2 != nil {
		return nil, errors.Wrap(err2, "diagramManager.addCopyWithRefinement failed")
	}
	return copy, nil
}

func (dmPtr *diagramManager) addDiagram(ownerID string, hl *core.HeldLocks) (core.Element, error) {
	diagram, err := dmPtr.newDiagram(hl)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.addDiagram failed")
	}
	err = diagram.SetOwningConceptID(ownerID, hl)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.addDiagram failed")
	}
	err = dmPtr.browserGUI.editor.SelectElement(diagram, hl)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.addDiagram failed")
	}
	hl.ReleaseLocksAndWait()
	err = dmPtr.displayDiagram(diagram, hl)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.addDiagram failed")
	}
	hl.ReleaseLocksAndWait()
	return diagram, nil
}

func (dmPtr *diagramManager) deleteDiagramElementView(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.browserGUI.GetUofD().GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.deleteDiagramElementView diagramElement not found for elementID " + elementID)
	}
	dEls := mapset.NewSet(diagramElement.GetConceptID(hl))
	return dmPtr.browserGUI.GetUofD().DeleteElements(dEls, hl)
}

// diagramClick handles the creation of a new Element and adding a node representation of it to the diagram
func (dmPtr *diagramManager) diagramClick(request *Request, hl *core.HeldLocks) error {
	uOfD := dmPtr.browserGUI.GetUofD()
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
		el.SetLabel(dmPtr.browserGUI.editor.GetDefaultElementLabel(), hl)
	case "Literal":
		el, err = uOfD.NewLiteral(hl)
		el.SetLabel(dmPtr.browserGUI.editor.GetDefaultLiteralLabel(), hl)
	case "Reference":
		el, err = uOfD.NewReference(hl)
		el.SetLabel(dmPtr.browserGUI.editor.GetDefaultReferenceLabel(), hl)
	case "Refinement":
		el, err = uOfD.NewRefinement(hl)
		el.SetLabel(dmPtr.browserGUI.editor.GetDefaultRefinementLabel(), hl)
	case "Diagram":
		el, err = crldiagramdomain.NewDiagram(uOfD, hl)
		el.SetLabel(dmPtr.browserGUI.editor.GetDefaultDiagramLabel(), hl)
	}
	if err != nil {
		return errors.Wrap(err, "diagramManager.diagramClick failed")
	}
	el.SetOwningConceptID(diagram.GetOwningConceptID(hl), hl)
	dmPtr.browserGUI.editor.SelectElement(el, hl)

	// Now the view
	var x, y float64
	x, err = strconv.ParseFloat(request.AdditionalParameters["NodeX"], 64)
	if err != nil {
		return errors.Wrap(err, "diagramManager.diagramClick failed")
	}
	y, err = strconv.ParseFloat(request.AdditionalParameters["NodeY"], 64)
	if err != nil {
		return errors.Wrap(err, "diagramManager.diagramClick failed")
	}
	var newNode core.Element
	newNode, err = crldiagramdomain.NewDiagramNode(uOfD, hl)
	if err != nil {
		return errors.Wrap(err, "diagramManager.diagramClick failed")
	}
	crldiagramdomain.SetNodeX(newNode, x, hl)
	crldiagramdomain.SetNodeY(newNode, y, hl)
	newNode.SetLabel(el.GetLabel(hl), hl)
	crldiagramdomain.SetReferencedModelElement(newNode, el, hl)
	crldiagramdomain.SetDisplayLabel(newNode, el.GetLabel(hl), hl)

	newNode.SetOwningConceptID(diagramID, hl)
	err = newNode.Register(dmPtr.elementManager)
	if err != nil {
		return errors.Wrap(err, "diagramManager.diagramClick failed")
	}
	dmPtr.browserGUI.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	hl.ReleaseLocksAndWait()

	return nil
}

// diagramDrop evaluates the request resulting from a drop in the diagram
func (dmPtr *diagramManager) diagramDrop(request *Request, hl *core.HeldLocks) error {
	if request.AdditionalParameters["Shift"] == "false" {
		_, err := dmPtr.addConceptView(request, hl)
		if err != nil {
			return errors.Wrap(err, "diagramManager.diagramDrop failed")
		}
		dmPtr.browserGUI.SetTreeDragSelection("")
	} else {
		_, err := dmPtr.addCopyWithRefinement(request, hl)
		if err != nil {
			return errors.Wrap(err, "diagramManager.diagramDrop failed")
		}
	}
	dmPtr.browserGUI.SetTreeDragSelection("")
	return nil
}

// DiagramViewHasBeenClosed notifies the server that the client has closed the diagram view
func (dmPtr *diagramManager) DiagramViewHasBeenClosed(diagramID string, hl *core.HeldLocks) error {
	if dmPtr.browserGUI.editor.IsDiagramDisplayed(diagramID, hl) {
		dmPtr.browserGUI.editor.RemoveDiagramFromDisplayedList(diagramID, hl)
	}
	return nil
}

// displayDiagram tells the client to display the indicated diagram.
func (dmPtr *diagramManager) displayDiagram(diagram core.Element, hl *core.HeldLocks) error {
	diagramID := diagram.GetConceptID(hl)
	if !diagram.IsRefinementOfURI(crldiagramdomain.CrlDiagramURI, hl) {
		return errors.New("In diagramManager.displayDiagram, the supplied diagram is not a refinement of CrlDiagramURI")
	}
	// Make sure the diagram is in the list of displayed diagrams
	if !dmPtr.browserGUI.editor.IsDiagramDisplayed(diagramID, hl) {
		err3 := dmPtr.browserGUI.editor.AddDiagramToDisplayedList(diagramID, hl)
		if err3 != nil {
			return errors.Wrap(err3, "diagramManager.displayDiagram failed")
		}
	}
	// make sure there is a monitor on the diagram so we know when it has been deleted
	err := diagram.Register(dmPtr)
	if err != nil {
		return errors.Wrap(err, "diagramManager.displayDiagram failed")
	}
	// Tell the client to display the diagram
	conceptState, err2 := core.NewConceptState(diagram)
	if err2 != nil {
		return errors.Wrap(err2, "diagramManager.displayDiagram failed")
	}
	notificationResponse, err := BrowserGUISingleton.GetClientNotificationManager().SendNotification("DisplayDiagram", diagram.GetConceptID(hl), conceptState, nil)
	if err != nil {
		return errors.Wrap(err, "diagramManager.displayDiagram failed")
	}
	if notificationResponse == nil {
		return errors.New("In diagramManager.displayDiagram the notification response was nil")
	}
	if notificationResponse.Result != 0 {
		return errors.New("In diagramManager.displayDiagram, notificationResponse was not 0")
	}
	return dmPtr.refreshDiagram(diagram, hl)
}

func (dmPtr *diagramManager) formatChanged(diagramElement core.Element, lineColor string, bgColor string, hl *core.HeldLocks) error {
	crldiagramdomain.SetLineColor(diagramElement, lineColor, hl)
	crldiagramdomain.SetBGColor(diagramElement, bgColor, hl)
	return nil
}

func (dmPtr *diagramManager) elementPointerChanged(linkID string, sourceID string, targetID string, targetAttributeName string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.browserGUI.GetUofD()
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.elementPointerChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagramdomain.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.elementPointerChanged called with model source not found")
	}
	var modelReference core.Reference
	switch typedModelSource := modelSource.(type) {
	case core.Reference:
		modelReference = typedModelSource
	default:
		return "", errors.New("diagramManager.elementPointerChanged called with source not being a Reference")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.elementPointerChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagramdomain.GetReferencedModelElement(diagramTarget, hl)
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
	if linkID == "" {
		// this is a new link
		diagramPointer, err = crldiagramdomain.NewDiagramElementPointer(uOfD, hl)
		if err != nil {
			return "", err
		}
		crldiagramdomain.SetReferencedModelElement(diagramPointer, modelSource, hl)
		crldiagramdomain.SetLinkSource(diagramPointer, diagramSource, hl)
		crldiagramdomain.SetLinkTarget(diagramPointer, diagramTarget, hl)
		if attributeName == core.NoAttribute {
			modelReference.SetReferencedConcept(modelTarget, attributeName, hl)
		} else {
			// for references to pointers, it is the pointer owner that is the referenced concept
			modelReference.SetReferencedConcept(modelSource, attributeName, hl)
		}
		diagramPointer.SetOwningConceptID(diagramSource.GetOwningConceptID(hl), hl)
		dmPtr.browserGUI.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramPointer = uOfD.GetElement(linkID)
		if diagramPointer == nil {
			return "", errors.New("diagramManager.elementPointerChanged called with diagramPointer not found in diagram")
		}
		if diagramSource != crldiagramdomain.GetLinkSource(diagramPointer, hl) {
			crldiagramdomain.SetLinkSource(diagramPointer, diagramSource, hl)
		}
		if diagramTarget != crldiagramdomain.GetLinkTarget(diagramPointer, hl) {
			crldiagramdomain.SetLinkTarget(diagramPointer, diagramTarget, hl)
		}
		if modelSource != crldiagramdomain.GetReferencedModelElement(diagramPointer, hl) {
			crldiagramdomain.SetReferencedModelElement(diagramPointer, modelSource, hl)
		}
		if attributeName == core.NoAttribute {
			if modelTarget != modelReference.GetReferencedConcept(hl) {
				modelReference.SetReferencedConcept(modelTarget, attributeName, hl)
			}
		} else {
			// for references to pointers, it is the pointer owner that is the referenced concept
			if modelSource != modelReference.GetReferencedConcept(hl) {
				modelReference.SetReferencedConcept(modelSource, attributeName, hl)
			}

		}
	}

	return diagramPointer.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) initialize() error {
	dmPtr.diagrams = map[string]core.Element{}
	return nil
}

// newDiagram creates a new crldiagram
func (dmPtr *diagramManager) newDiagram(hl *core.HeldLocks) (core.Element, error) {
	// Insert name prompt here
	name := dmPtr.browserGUI.editor.GetDefaultDiagramLabel()
	uOfD := BrowserGUISingleton.GetUofD()
	diagram, err := crldiagramdomain.NewDiagram(uOfD, hl)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.newDiagram failed")
	}
	diagram.SetLabel(name, hl)
	hl.ReleaseLocksAndWait()
	dmPtr.diagrams[diagram.GetConceptID(hl)] = diagram
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.newDiagram failed")
	}
	err = diagram.Register(dmPtr)
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.newDiagram failed")
	}
	return diagram, nil
}

func (dmPtr *diagramManager) ownerPointerChanged(linkID string, sourceID string, targetID string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.browserGUI.GetUofD()
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.ownerPointerChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagramdomain.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.ownerPointerChanged called with model source not found")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.ownerPointerChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagramdomain.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.ownerPointerChanged called with model target not found")
	}
	var err error
	var diagramPointer core.Element
	if linkID == "" {
		// this is a new link
		diagramPointer, err = crldiagramdomain.NewDiagramOwnerPointer(uOfD, hl)
		if err != nil {
			return "", err
		}
		crldiagramdomain.SetReferencedModelElement(diagramPointer, modelSource, hl)
		crldiagramdomain.SetLinkSource(diagramPointer, diagramSource, hl)
		crldiagramdomain.SetLinkTarget(diagramPointer, diagramTarget, hl)
		modelSource.SetOwningConcept(modelTarget, hl)
		diagramPointer.SetOwningConceptID(diagramSource.GetOwningConceptID(hl), hl)
		dmPtr.browserGUI.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramPointer = uOfD.GetElement(linkID)
		if diagramPointer == nil {
			return "", errors.New("diagramManager.ownerPointerChanged called with diagramPointer not found in diagram")
		}
		if diagramSource != crldiagramdomain.GetLinkSource(diagramPointer, hl) {
			crldiagramdomain.SetLinkSource(diagramPointer, diagramSource, hl)
		}
		if diagramTarget != crldiagramdomain.GetLinkTarget(diagramPointer, hl) {
			crldiagramdomain.SetLinkTarget(diagramPointer, diagramTarget, hl)
		}
		if modelSource != crldiagramdomain.GetReferencedModelElement(diagramPointer, hl) {
			crldiagramdomain.SetReferencedModelElement(diagramPointer, modelSource, hl)
		}
		if modelTarget != modelSource.GetOwningConcept(hl) {
			modelSource.SetOwningConcept(modelTarget, hl)
		}
	}

	return diagramPointer.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) refinedPointerChanged(linkID string, sourceID string, targetID string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.browserGUI.GetUofD()
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.refinedPointerChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagramdomain.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.elementPoirefinedPointerChangednterChanged called with model source not found")
	}
	var modelRefinement core.Refinement
	switch typedModelSource := modelSource.(type) {
	case core.Refinement:
		modelRefinement = typedModelSource.(core.Refinement)
	default:
		return "", errors.New("diagramManager.refinedPointerChanged called with source not being a Refinement")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.refinedPointerChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagramdomain.GetReferencedModelElement(diagramTarget, hl)
	if modelTarget == nil {
		return "", errors.New("diagramManager.refinedPointerChanged called with model target not found")
	}
	var err error
	var diagramPointer core.Element
	if linkID == "" {
		// this is a new link
		diagramPointer, err = crldiagramdomain.NewDiagramRefinedPointer(uOfD, hl)
		if err != nil {
			return "", err
		}
		crldiagramdomain.SetReferencedModelElement(diagramPointer, modelSource, hl)
		crldiagramdomain.SetLinkSource(diagramPointer, diagramSource, hl)
		crldiagramdomain.SetLinkTarget(diagramPointer, diagramTarget, hl)
		modelRefinement.SetRefinedConcept(modelTarget, hl)
		diagramPointer.SetOwningConceptID(diagramSource.GetOwningConceptID(hl), hl)
		dmPtr.browserGUI.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramPointer = uOfD.GetElement(linkID)
		if diagramPointer == nil {
			return "", errors.New("diagramManager.refinedPointerChanged called with diagramPointer not found in diagram")
		}
		if diagramSource != crldiagramdomain.GetLinkSource(diagramPointer, hl) {
			crldiagramdomain.SetLinkSource(diagramPointer, diagramSource, hl)
		}
		if diagramTarget != crldiagramdomain.GetLinkTarget(diagramPointer, hl) {
			crldiagramdomain.SetLinkTarget(diagramPointer, diagramTarget, hl)
		}
		if modelSource != crldiagramdomain.GetReferencedModelElement(diagramPointer, hl) {
			crldiagramdomain.SetReferencedModelElement(diagramPointer, modelSource, hl)
		}
		if modelTarget != modelRefinement.GetRefinedConcept(hl) {
			modelRefinement.SetRefinedConcept(modelTarget, hl)
		}
	}

	return diagramPointer.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) ReferenceLinkChanged(linkID string, sourceID string, targetID string, targetAttributeName string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.browserGUI.GetUofD()
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagramdomain.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with model source not found")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagramdomain.GetReferencedModelElement(diagramTarget, hl)
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
		newReference.SetReferencedConcept(modelTarget, attributeName, hl)
		newReference.SetOwningConcept(modelSource, hl)
		diagramLink, err = crldiagramdomain.NewDiagramReferenceLink(uOfD, hl)
		if err != nil {
			return "", err
		}
		crldiagramdomain.SetReferencedModelElement(diagramLink, newReference, hl)
		crldiagramdomain.SetLinkSource(diagramLink, diagramSource, hl)
		crldiagramdomain.SetLinkTarget(diagramLink, diagramTarget, hl)
		modelElement = newReference
		diagramLink.SetOwningConcept(diagram, hl)
		err = diagramLink.Register(dmPtr.elementManager)
		if err != nil {
			return "", errors.Wrap(err, "diagramManager.ReferenceLinkChanged failed")
		}
		dmPtr.browserGUI.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramLink = uOfD.GetElement(linkID)
		modelElement = crldiagramdomain.GetReferencedModelElement(diagramLink, hl)
		if modelElement != nil {
			switch typedModelElement := modelElement.(type) {
			case core.Reference:
				reference := typedModelElement
				if diagramLink == nil {
					return "", errors.New("diagramManager.refinementLinkChanged called with diagramPointer not found in diagram")
				}
				if reference.GetOwningConcept(hl) != modelSource {
					reference.SetOwningConcept(modelSource, hl)
				}
				if reference.GetReferencedConcept(hl) != modelTarget {
					reference.SetReferencedConcept(modelTarget, attributeName, hl)
				}
			}
		}
	}
	dmPtr.browserGUI.editor.SelectElement(modelElement, hl)

	return diagramLink.GetConceptID(hl), nil
}

func (dmPtr *diagramManager) RefinementLinkChanged(linkID string, sourceID string, targetID string, hl *core.HeldLocks) (string, error) {
	uOfD := dmPtr.browserGUI.GetUofD()
	diagramSource := uOfD.GetElement(sourceID)
	if diagramSource == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with sourceID not found in diagram")
	}
	modelSource := crldiagramdomain.GetReferencedModelElement(diagramSource, hl)
	if modelSource == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with model source not found")
	}
	diagramTarget := uOfD.GetElement(targetID)
	if diagramTarget == nil {
		return "", errors.New("diagramManager.refinementLinkChanged called with targetID not found in diagram")
	}
	modelTarget := crldiagramdomain.GetReferencedModelElement(diagramTarget, hl)
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
		diagramLink, err = crldiagramdomain.NewDiagramRefinementLink(uOfD, hl)
		if err != nil {
			return "", err
		}
		crldiagramdomain.SetReferencedModelElement(diagramLink, newRefinement, hl)
		crldiagramdomain.SetLinkSource(diagramLink, diagramSource, hl)
		crldiagramdomain.SetLinkTarget(diagramLink, diagramTarget, hl)
		diagramLink.SetOwningConcept(diagram, hl)
		modelElement = newRefinement
		err = modelElement.Register(dmPtr.elementManager)
		if err != nil {
			return "", errors.Wrap(err, "diagramManager.ReferenceLinkChanged failed")
		}
		dmPtr.browserGUI.SendNotification("ClearToolbarSelection", "", nil, map[string]string{})
	} else {
		diagramLink = uOfD.GetElement(linkID)
		modelElement = crldiagramdomain.GetReferencedModelElement(diagramLink, hl)
		if modelElement != nil {
			switch typedModelElement := modelElement.(type) {
			case core.Refinement:
				refinement := typedModelElement
				if diagramLink == nil {
					return "", errors.New("diagramManager.refinementLinkChanged called with diagramPointer not found in diagram")
				}
				if diagramSource != crldiagramdomain.GetLinkSource(diagramLink, hl) {
					refinement.SetRefinedConcept(modelSource, hl)
				}
				if diagramTarget != crldiagramdomain.GetLinkTarget(diagramLink, hl) {
					refinement.SetAbstractConcept(modelTarget, hl)
				}
			}
		}
	}
	dmPtr.browserGUI.editor.SelectElement(modelElement, hl)

	return diagramLink.GetConceptID(hl), nil
}

// refreshDiagramUsingURI finds the diagram and resends all diagram elements to the browser
func (dmPtr *diagramManager) refreshDiagramUsingURI(diagramID string, hl *core.HeldLocks) error {
	diagram := dmPtr.browserGUI.GetUofD().GetElement(diagramID)
	if diagram == nil {
		return errors.New("In diagramManager.refreshDiagram, diagram not found for ID: " + diagramID)
	}
	return dmPtr.refreshDiagram(diagram, hl)
}

// refreshDiagram resends all diagram elements to the browser
func (dmPtr *diagramManager) refreshDiagram(diagram core.Element, hl *core.HeldLocks) error {
	nodes := diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, hl)
	for _, node := range nodes {
		additionalParameters := getNodeAdditionalParameters(node, hl)
		conceptState, err2 := core.NewConceptState(node)
		if err2 != nil {
			return errors.Wrap(err2, "diagramManager.refreshDiagram failed")
		}
		node.Register(dmPtr.elementManager)
		notificationResponse, err := BrowserGUISingleton.SendNotification("AddDiagramNode", node.GetConceptID(hl), conceptState, additionalParameters)
		if err != nil {
			return errors.Wrap(err, "diagramManager.refreshDiagram failed")
		}
		if notificationResponse.Result != 0 {
			return errors.New(notificationResponse.ErrorMessage)
		}
	}
	links := diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramLinkURI, hl)
	for _, link := range links {
		additionalParameters := getLinkAdditionalParameters(link, hl)
		conceptState, err2 := core.NewConceptState(link)
		if err2 != nil {
			return errors.Wrap(err2, "diagramManager.refreshDiagram failed")
		}
		link.Register(dmPtr.elementManager)
		notificationResponse, err := BrowserGUISingleton.SendNotification("AddDiagramLink", link.GetConceptID(hl), conceptState, additionalParameters)
		if err != nil {
			return errors.Wrap(err, "diagramManager.refreshDiagram failed")
		}
		if notificationResponse.Result != 0 {
			return errors.New(notificationResponse.ErrorMessage)
		}
	}
	return nil
}

// setDiagramNodePosition sets the position of the diagram node
func (dmPtr *diagramManager) setDiagramNodePosition(nodeID string, x float64, y float64, hl *core.HeldLocks) {
	node := BrowserGUISingleton.GetUofD().GetElement(nodeID)
	if node == nil {
		// This can happen when the concept space containing the diagram is deleted???
		log.Print("In setDiagramNodePosition node not found for nodeID: " + nodeID)
		return
	}
	crldiagramdomain.SetNodeX(node, x, hl)
	crldiagramdomain.SetNodeY(node, y, hl)
}

func (dmPtr *diagramManager) showAbstractConcept(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.browserGUI.GetUofD().GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showAbstractConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showAbstractConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showAbstractConcept modelConcept not found for elementID " + elementID)
	}
	var modelRefinement core.Refinement
	switch typedModelConcept := modelConcept.(type) {
	case core.Refinement:
		modelRefinement = typedModelConcept
	default:
		return errors.New("diagramManager.showAbstractConcept modelConcept is not a Refinement")
	}
	modelAbstractConcept := modelRefinement.GetAbstractConcept(hl)
	if modelAbstractConcept == nil {
		return errors.New("Abstract Concept is nil")
	}
	diagramAbstractConcept := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelAbstractConcept, hl)
	if diagramAbstractConcept == nil {
		diagramAbstractConcept, _ = crldiagramdomain.NewDiagramNode(dmPtr.browserGUI.GetUofD(), hl)
		crldiagramdomain.SetReferencedModelElement(diagramAbstractConcept, modelAbstractConcept, hl)
		crldiagramdomain.SetDisplayLabel(diagramAbstractConcept, modelAbstractConcept.GetLabel(hl), hl)
		diagramElementX := crldiagramdomain.GetNodeX(diagramElement, hl)
		diagramElementY := crldiagramdomain.GetNodeY(diagramElement, hl)
		crldiagramdomain.SetNodeX(diagramAbstractConcept, diagramElementX, hl)
		crldiagramdomain.SetNodeY(diagramAbstractConcept, diagramElementY-100, hl)
		diagramAbstractConcept.SetOwningConcept(diagram, hl)
	}
	elementPointer := crldiagramdomain.GetElementPointer(diagram, diagramElement, hl)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramAbstractPointer(dmPtr.browserGUI.GetUofD(), hl)
		crldiagramdomain.SetReferencedModelElement(elementPointer, modelConcept, hl)
		crldiagramdomain.SetLinkSource(elementPointer, diagramElement, hl)
		crldiagramdomain.SetLinkTarget(elementPointer, diagramAbstractConcept, hl)
		elementPointer.SetOwningConcept(diagram, hl)
	}
	return nil
}

func (dmPtr *diagramManager) showOwnedConcepts(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.browserGUI.GetUofD().GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showOwnedConcepts diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showOwnedConcepts diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showOwnedConcepts modelConcept not found for elementID " + elementID)
	}
	it := modelConcept.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	var offset float64
	for id := range it.C {
		child := dmPtr.browserGUI.GetUofD().GetElement(id.(string))
		if child == nil {
			return errors.New("Child Concept is nil for id " + id.(string))
		}
		diagramChildConcept := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, child, hl)
		if diagramChildConcept == nil {
			diagramChildConcept, _ = crldiagramdomain.NewDiagramNode(dmPtr.browserGUI.GetUofD(), hl)
			crldiagramdomain.SetReferencedModelElement(diagramChildConcept, child, hl)
			crldiagramdomain.SetDisplayLabel(diagramChildConcept, child.GetLabel(hl), hl)
			diagramElementX := crldiagramdomain.GetNodeX(diagramElement, hl)
			diagramElementY := crldiagramdomain.GetNodeY(diagramElement, hl)
			crldiagramdomain.SetNodeX(diagramChildConcept, diagramElementX+offset, hl)
			crldiagramdomain.SetNodeY(diagramChildConcept, diagramElementY+50, hl)
			diagramChildConcept.SetOwningConcept(diagram, hl)
		}
		ownerPointer := crldiagramdomain.GetOwnerPointer(diagram, diagramElement, hl)
		if ownerPointer == nil {
			ownerPointer, _ = crldiagramdomain.NewDiagramOwnerPointer(dmPtr.browserGUI.GetUofD(), hl)
			crldiagramdomain.SetReferencedModelElement(ownerPointer, modelConcept, hl)
			crldiagramdomain.SetLinkSource(ownerPointer, diagramChildConcept, hl)
			crldiagramdomain.SetLinkTarget(ownerPointer, diagramElement, hl)
			ownerPointer.SetOwningConcept(diagram, hl)
		}
		offset = offset + 50
	}
	return nil
}

func (dmPtr *diagramManager) showOwner(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.browserGUI.GetUofD().GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showOwner diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showOwner diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showOwner modelConcept not found for elementID " + elementID)
	}
	modelConceptOwner := modelConcept.GetOwningConcept(hl)
	if modelConceptOwner == nil {
		return errors.New("Owner is nil")
	}
	diagramConceptOwner := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelConceptOwner, hl)
	if diagramConceptOwner == nil {
		diagramConceptOwner, _ = crldiagramdomain.NewDiagramNode(dmPtr.browserGUI.GetUofD(), hl)
		crldiagramdomain.SetReferencedModelElement(diagramConceptOwner, modelConceptOwner, hl)
		crldiagramdomain.SetDisplayLabel(diagramConceptOwner, modelConceptOwner.GetLabel(hl), hl)
		diagramElementX := crldiagramdomain.GetNodeX(diagramElement, hl)
		diagramElementY := crldiagramdomain.GetNodeY(diagramElement, hl)
		crldiagramdomain.SetNodeX(diagramConceptOwner, diagramElementX, hl)
		crldiagramdomain.SetNodeY(diagramConceptOwner, diagramElementY-100, hl)
		diagramConceptOwner.SetOwningConcept(diagram, hl)
	}
	ownerPointer := crldiagramdomain.GetOwnerPointer(diagram, diagramElement, hl)
	if ownerPointer == nil {
		ownerPointer, _ = crldiagramdomain.NewDiagramOwnerPointer(dmPtr.browserGUI.GetUofD(), hl)
		crldiagramdomain.SetReferencedModelElement(ownerPointer, modelConcept, hl)
		crldiagramdomain.SetLinkSource(ownerPointer, diagramElement, hl)
		crldiagramdomain.SetLinkTarget(ownerPointer, diagramConceptOwner, hl)
		ownerPointer.SetOwningConcept(diagram, hl)
	}
	return nil
}

func (dmPtr *diagramManager) showReferencedConcept(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.browserGUI.GetUofD().GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showReferencedConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showReferencedConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showReferencedConcept modelConcept not found for elementID " + elementID)
	}
	var modelReference core.Reference
	switch typedModelConcept := modelConcept.(type) {
	case core.Reference:
		modelReference = typedModelConcept
	default:
		return errors.New("diagramManager.showReferencedConcept modelConcept is not a Reference")
	}
	modelReferencedConcept := modelReference.GetReferencedConcept(hl)
	if modelReferencedConcept == nil {
		return errors.New("Referenced Concept is nil")
	}
	var diagramReferencedConcept core.Element
	switch modelReference.GetReferencedAttributeName(hl) {
	case core.NoAttribute, core.LiteralValue:
		diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelReferencedConcept, hl)
		if diagramReferencedConcept == nil {
			diagramReferencedConcept, _ = crldiagramdomain.NewDiagramNode(dmPtr.browserGUI.GetUofD(), hl)
			crldiagramdomain.SetReferencedModelElement(diagramReferencedConcept, modelReferencedConcept, hl)
			crldiagramdomain.SetDisplayLabel(diagramReferencedConcept, modelReferencedConcept.GetLabel(hl), hl)
			diagramElementX := crldiagramdomain.GetNodeX(diagramElement, hl)
			diagramElementY := crldiagramdomain.GetNodeY(diagramElement, hl)
			crldiagramdomain.SetNodeX(diagramReferencedConcept, diagramElementX, hl)
			crldiagramdomain.SetNodeY(diagramReferencedConcept, diagramElementY-100, hl)
			diagramReferencedConcept.SetOwningConcept(diagram, hl)
		}
	case core.OwningConceptID:
		diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConceptOwnerPointer(diagram, modelReferencedConcept, hl)
		if diagramReferencedConcept == nil {
			return errors.New("No representation of the owner pointer currently exists in this diagram")
		}
	case core.ReferencedConceptID:
		switch typedModelReferencedConcept := modelReferencedConcept.(type) {
		case core.Reference:
			diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConceptElementPointer(diagram, typedModelReferencedConcept, hl)
			if diagramReferencedConcept == nil {
				return errors.New("No representation of the referenced concept pointer currently exists in this diagram")
			}
		}
	case core.AbstractConceptID:
		switch typedModelReferencedConcept := modelReferencedConcept.(type) {
		case core.Refinement:
			diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConceptAbstractPointer(diagram, typedModelReferencedConcept, hl)
		}
		if diagramReferencedConcept == nil {
			return errors.New("No representation of the abstract concept pointer currently exists in this diagram")
		}
	case core.RefinedConceptID:
		switch typedModelReferencedConcept := modelReferencedConcept.(type) {
		case core.Refinement:
			diagramReferencedConcept = crldiagramdomain.GetFirstElementRepresentingConceptRefinedPointer(diagram, typedModelReferencedConcept, hl)
			if diagramReferencedConcept == nil {
				return errors.New("No representation of the refined concept pointer currently exists in this diagram")
			}
		}
	}
	elementPointer := crldiagramdomain.GetElementPointer(diagram, diagramElement, hl)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramElementPointer(dmPtr.browserGUI.GetUofD(), hl)
		crldiagramdomain.SetReferencedModelElement(elementPointer, modelConcept, hl)
		crldiagramdomain.SetLinkSource(elementPointer, diagramElement, hl)
		crldiagramdomain.SetLinkTarget(elementPointer, diagramReferencedConcept, hl)
		elementPointer.SetOwningConcept(diagram, hl)
	}
	return nil
}

func (dmPtr *diagramManager) showRefinedConcept(elementID string, hl *core.HeldLocks) error {
	diagramElement := dmPtr.browserGUI.GetUofD().GetElement(elementID)
	if diagramElement == nil {
		return errors.New("diagramManager.showRefinedConcept diagramElement not found for elementID " + elementID)
	}
	diagram := diagramElement.GetOwningConcept(hl)
	if diagram == nil {
		return errors.New("diagramManager.showRefinedConcept diagram not found for elementID " + elementID)
	}
	modelConcept := crldiagramdomain.GetReferencedModelElement(diagramElement, hl)
	if modelConcept == nil {
		return errors.New("diagramManager.showRefinedConcept modelConcept not found for elementID " + elementID)
	}
	var modelRefinement core.Refinement
	switch typedModelConcept := modelConcept.(type) {
	case core.Refinement:
		modelRefinement = typedModelConcept
	default:
		return errors.New("diagramManager.showRefinedConcept modelConcept is not a Refinement")
	}
	modelRefinedConcept := modelRefinement.GetRefinedConcept(hl)
	if modelRefinedConcept == nil {
		return errors.New("Refined Concept is nil")
	}
	diagramRefinedConcept := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelRefinedConcept, hl)
	if diagramRefinedConcept == nil {
		diagramRefinedConcept, _ = crldiagramdomain.NewDiagramNode(dmPtr.browserGUI.GetUofD(), hl)
		crldiagramdomain.SetReferencedModelElement(diagramRefinedConcept, modelRefinedConcept, hl)
		crldiagramdomain.SetDisplayLabel(diagramRefinedConcept, modelRefinedConcept.GetLabel(hl), hl)
		diagramElementX := crldiagramdomain.GetNodeX(diagramElement, hl)
		diagramElementY := crldiagramdomain.GetNodeY(diagramElement, hl)
		crldiagramdomain.SetNodeX(diagramRefinedConcept, diagramElementX, hl)
		crldiagramdomain.SetNodeY(diagramRefinedConcept, diagramElementY-100, hl)
		diagramRefinedConcept.SetOwningConcept(diagram, hl)
	}
	elementPointer := crldiagramdomain.GetElementPointer(diagram, diagramElement, hl)
	if elementPointer == nil {
		elementPointer, _ = crldiagramdomain.NewDiagramRefinedPointer(dmPtr.browserGUI.GetUofD(), hl)
		crldiagramdomain.SetReferencedModelElement(elementPointer, modelConcept, hl)
		crldiagramdomain.SetLinkSource(elementPointer, diagramElement, hl)
		crldiagramdomain.SetLinkTarget(elementPointer, diagramRefinedConcept, hl)
		elementPointer.SetOwningConcept(diagram, hl)
	}
	return nil
}

// Update handles additions and removals of diagram elements from the diagram view
// Note that it cannot delete the diagram view in the GUI because this function will never get called: once the diagram has been
// deleted, queuing of functions related to it is suppressed. That's what the DiagramViewMonitor is for.
func (dmPtr *diagramManager) Update(notification *core.ChangeNotification, hl *core.HeldLocks) error {
	uOfD := hl.GetUniverseOfDiscourse()
	switch notification.GetNatureOfChange() {
	case core.ConceptRemoved:
		if notification.GetBeforeConceptState() != nil {
			err := dmPtr.browserGUI.editor.CloseDiagramView(notification.GetBeforeConceptState().ConceptID, hl)
			if err != nil {
				return errors.Wrap(err, "diagramManager.Update failed")
			}
		} else {
			return errors.New("diagramManager.Update called with ConceptRemoved but beforeConceptState being nil")
		}
	case core.OwningConceptChanged:
		if notification.GetAfterReferencedState() == nil ||
			(notification.GetBeforeReferencedState() != nil && notification.GetBeforeReferencedState().ConceptID != notification.GetAfterReferencedState().ConceptID) {
			// If the diagram was the owner but is no longer the owner, then remove the diagram element view
			beforeState := notification.GetBeforeConceptState()
			ownerID := ""
			if notification.GetBeforeReferencedState() != nil {
				ownerID = notification.GetBeforeReferencedState().ConceptID
			}
			additionalParameters := map[string]string{"OwnerID": ownerID}
			SendNotification("DeleteDiagramElement", beforeState.ConceptID, beforeState, additionalParameters)
		} else if notification.GetBeforeReferencedState() == nil ||
			(notification.GetAfterReferencedState() != nil && notification.GetBeforeReferencedState().ConceptID != notification.GetAfterReferencedState().ConceptID) {
			// we have to add the diagram element view
			afterState := notification.GetAfterConceptState()
			newElement := uOfD.GetElement(afterState.ConceptID)
			newElement.Register(dmPtr.elementManager)
			if crldiagramdomain.IsDiagramNode(newElement, hl) {
				additionalParameters := getNodeAdditionalParameters(newElement, hl)
				_, err := SendNotification("AddDiagramNode", newElement.GetConceptID(hl), afterState, additionalParameters)
				return err
			} else if crldiagramdomain.IsDiagramLink(newElement, hl) {
				additionalParameters := getLinkAdditionalParameters(newElement, hl)
				_, err := SendNotification("AddDiagramLink", newElement.GetConceptID(hl), afterState, additionalParameters)
				return err
			}
		}
	case core.ConceptChanged:
		beforeState := notification.GetBeforeConceptState()
		afterState := notification.GetAfterConceptState()
		if beforeState != nil && afterState != nil {
			if beforeState.Label != afterState.Label {
				_, err := SendNotification("DiagramLabelChanged", afterState.ConceptID, afterState, map[string]string{})
				return err
			}
		}
	}
	return nil
}
