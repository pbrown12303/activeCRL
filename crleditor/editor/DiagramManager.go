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

// DiagramManager manages the diagram portion of the UI and all interactions with it
type DiagramManager struct {
	abstractDiagram core.Element
	crlEditor       *CrlEditor
}

func newDiagramManager(crlEditor *CrlEditor) *DiagramManager {
	var diagramManager DiagramManager
	diagramManager.crlEditor = crlEditor
	diagramManager.abstractDiagram = CrlEditorSingleton.GetUofD().GetElementWithURI(crldiagram.CrlDiagramURI)
	uOfD := CrlEditorSingleton.GetUofD()
	uOfD.AddFunction(crldiagram.CrlDiagramNodeURI, updateDiagramNodeView)
	return &diagramManager
}

func (dmPtr *DiagramManager) addNodeView(request *Request, hl *core.HeldLocks) (core.Element, error) {
	uOfD := dmPtr.crlEditor.GetUofD()
	diagramID := request.AdditionalParameters["DiagramID"]
	el := uOfD.GetElement(dmPtr.crlEditor.GetTreeDragSelectionID(hl))
	if el == nil {
		return nil, errors.New("Indicated model element not found in addNodeView, ID: " + request.RequestConceptID)
	}
	x, err := strconv.ParseFloat(request.AdditionalParameters["NodeX"], 64)
	if err != nil {
		return nil, err
	}
	y, err2 := strconv.ParseFloat(request.AdditionalParameters["NodeY"], 64)
	if err2 != nil {
		return nil, err
	}

	newNode, err := uOfD.CreateReplicateAsRefinementFromURI(crldiagram.CrlDiagramNodeURI, hl)
	if err != nil {
		js.Global.Get("console").Call("log", "Failed to create CrlDiagramNode"+err.Error())
		return nil, err
	}
	newNode.SetLabel(el.GetLabel(hl), hl)
	crldiagram.SetReferencedModelElement(newNode, el, hl)
	crldiagram.SetDisplayLabel(newNode, el.GetLabel(hl), hl)
	crldiagram.SetNodeX(newNode, x, hl)
	crldiagram.SetNodeY(newNode, y, hl)

	newNode.SetOwningConceptID(diagramID, hl)
	hl.ReleaseLocksAndWait()

	return newNode, nil
}

func (dmPtr *DiagramManager) addViewFunctionsToUofD(hl *core.HeldLocks) {
	uOfD := CrlEditorSingleton.GetUofD()
	addDiagramViewFunctionsToUofD(uOfD, hl)
	//	addDiagramNodeViewFunctionsToUofD(uOfD, hl)
	//	addDiagramLinkViewFunctionsToUofD(uOfD, hl)
}

// DiagramDrop evaluates the request resulting from a drop in the diagram
func (dmPtr *DiagramManager) DiagramDrop(request *Request, hl *core.HeldLocks) error {
	_, err := dmPtr.addNodeView(request, hl)
	if err != nil {
		return err
	}
	dmPtr.crlEditor.SetTreeDragSelection("")
	return nil
}

// DisplayDiagram tells the client to display the indicated diagram.
func (dmPtr *DiagramManager) DisplayDiagram(diagram core.Element, hl *core.HeldLocks) {
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
}

// NewDiagram creates a new crldiagram
func (dmPtr *DiagramManager) NewDiagram(hl *core.HeldLocks) core.Element {
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

// SetSize sets the size of the diagram rectangle
func (dmPtr *DiagramManager) SetSize() {
	// A test rectangle
	rect := js.Global.Get("joint").Get("shapes").Get("basic").Get("Rect").New(js.M{
		"position": js.M{
			"x": 100,
			"y": 30},
		"size": js.M{
			"width":  100,
			"height": 30},
		"attrs": js.M{
			"rect": js.M{
				"fill": "blue"},
			"text": js.M{
				"text": "my box",
				"fill": "white"}}})

	js.Global.Set("diagramRect", rect)
	js.Global.Get("jointGraph").Call("addCell", rect)

}
