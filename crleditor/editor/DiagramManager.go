package editor

import (
	"errors"
	"strconv"

	"github.com/pbrown12303/activeCRL/crldiagram"
	//	"fmt"
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/pbrown12303/activeCRL/core"
)

const diagramContainerSuffix = "DiagramContainer"
const diagramSuffix = "Diagram"

// diagramManager manages the diagram portion of the UI and all interactions with it
type diagramManager struct {
	crlEditor *CrlEditor
}

func newDiagramManager(crlEditor *CrlEditor) *diagramManager {
	var dm diagramManager
	dm.crlEditor = crlEditor
	uOfD := CrlEditorSingleton.GetUofD()
	addDiagramViewFunctionsToUofD(uOfD)
	return &dm
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
		createAsLink = CrlEditorSingleton.GetEditorSettings().DropReferenceAsLink
	case core.Refinement:
		createAsLink = CrlEditorSingleton.GetEditorSettings().DropRefinementAsLink
	}

	var newElement core.Element
	if createAsLink {
		newElement, err = uOfD.CreateReplicateAsRefinementFromURI(crldiagram.CrlDiagramLinkURI, hl)
		if err != nil {
			return nil, err
		}
		var modelSourceConcept core.Element
		var modelTargetConcept core.Element
		switch el.(type) {
		case core.Reference:
			reference := el.(core.Reference)
			modelSourceConcept = reference.GetOwningConcept(hl)
			modelTargetConcept = reference.GetReferencedConcept(hl)
		case core.Refinement:
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
		newElement, err = uOfD.CreateReplicateAsRefinementFromURI(crldiagram.CrlDiagramNodeURI, hl)
		if err != nil {
			js.Global.Get("console").Call("log", "Failed to create CrlDiagramNode"+err.Error())
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

// diagramDrop evaluates the request resulting from a drop in the diagram
func (dmPtr *diagramManager) diagramDrop(request *Request, hl *core.HeldLocks) error {
	_, err := dmPtr.addConceptView(request, hl)
	if err != nil {
		return err
	}
	dmPtr.crlEditor.SetTreeDragSelection("")
	return nil
}

// displayDiagram tells the client to display the indicated diagram.
func (dmPtr *diagramManager) displayDiagram(diagram core.Element, hl *core.HeldLocks) {
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("DisplayDiagram", diagram.GetConceptID(hl), diagram, nil)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if notificationResponse.Result != 0 {
		log.Print(notificationResponse.ErrorMessage)
		return
	}
	nodes := diagram.GetOwnedConceptsRefinedFromURI(crldiagram.CrlDiagramNodeURI, hl)
	for _, node := range nodes {
		additionalParameters := getNodeAdditionalParameters(node, hl)
		notificationResponse, err = CrlEditorSingleton.SendNotification("AddDiagramNode", node.GetConceptID(hl), node, additionalParameters)
		if err != nil {
			log.Printf(err.Error())
			break
		}
		if notificationResponse.Result != 0 {
			log.Print(notificationResponse.ErrorMessage)
			break
		}
	}
	links := diagram.GetOwnedConceptsRefinedFromURI(crldiagram.CrlDiagramLinkURI, hl)
	for _, link := range links {
		additionalParameters := getLinkAdditionalParameters(link, hl)
		notificationResponse, err = CrlEditorSingleton.SendNotification("AddDiagramLink", link.GetConceptID(hl), link, additionalParameters)
		if err != nil {
			log.Printf(err.Error())
			break
		}
		if notificationResponse.Result != 0 {
			log.Print(notificationResponse.ErrorMessage)
			break
		}
	}

}

// newDiagram creates a new crldiagram
func (dmPtr *diagramManager) newDiagram(hl *core.HeldLocks) core.Element {
	// Insert name prompt here
	name := getDefaultDiagramLabel()
	uOfD := CrlEditorSingleton.GetUofD()
	diagram, err := uOfD.CreateReplicateAsRefinementFromURI(crldiagram.CrlDiagramURI, hl)
	if err != nil {
		log.Print(err)
	}
	diagram.SetLabel(name, hl)
	hl.ReleaseLocksAndWait()
	return diagram
}

// setDiagramNodePosition sets the position of the diagram node
func (dmPtr *diagramManager) setDiagramNodePosition(nodeID string, x float64, y float64, hl *core.HeldLocks) {
	node := CrlEditorSingleton.GetUofD().GetElement(nodeID)
	if node == nil {
		log.Print("In SetDiagramNodePosition node not found for nodeID: " + nodeID)
		return
	}
	crldiagram.SetNodeX(node, x, hl)
	crldiagram.SetNodeY(node, y, hl)
}
