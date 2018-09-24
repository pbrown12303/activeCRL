package editor

import (
	//	"fmt"
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/crlDiagram"
)

const diagramContainerSuffix = "DiagramContainer"
const diagramSuffix = "Diagram"

// DiagramManager manages the diagram portion of the UI and all interactions with it
type DiagramManager struct {
	abstractDiagram                  core.Element
	crlDiagramUUIDToCrlDiagram       map[string]core.Element
	crlDiagramIDToDiagramContainer   map[string]*js.Object
	crlDiagramIDToJointGraph         map[string]*js.Object
	crlDiagramNodeUUIDToJointElement map[string]*js.Object
	diagramContainerIDToCrlDiagram   map[string]core.Element
	jointElementIDToCrlDiagramNode   map[string]core.Element
	jointGraphIDToCrlDiagram         map[string]core.Element
	jointPapers                      map[string]*js.Object
}

func newDiagramManager() *DiagramManager {
	var diagramManager DiagramManager
	diagramManager.abstractDiagram = CrlEditorSingleton.GetUofD().GetElementWithUri(crlDiagram.CrlDiagramUri)
	diagramManager.crlDiagramUUIDToCrlDiagram = make(map[string]core.Element)
	diagramManager.crlDiagramIDToDiagramContainer = make(map[string]*js.Object)
	diagramManager.crlDiagramIDToJointGraph = make(map[string]*js.Object)
	diagramManager.crlDiagramNodeUUIDToJointElement = make(map[string]*js.Object)
	diagramManager.diagramContainerIDToCrlDiagram = make(map[string]core.Element)
	diagramManager.jointElementIDToCrlDiagramNode = make(map[string]core.Element)
	diagramManager.jointGraphIDToCrlDiagram = make(map[string]core.Element)
	diagramManager.jointPapers = make(map[string]*js.Object)
	defineNodeViews()
	return &diagramManager
}

func (dmPtr *DiagramManager) addViewFunctionsToUofD() {
	uOfD := CrlEditorSingleton.GetUofD()
	hl := CrlEditorSingleton.getHeldLocks()
	addDiagramViewFunctionsToUofD(uOfD, hl)
	//	addDiagramNodeViewFunctionsToUofD(uOfD, hl)
	//	addDiagramLinkViewFunctionsToUofD(uOfD, hl)
}

// GetCrlDiagramForContainerID returns the CrlDiagram being displayed in the container
func (dmPtr *DiagramManager) GetCrlDiagramForContainerID(containerId string) core.Element {
	return dmPtr.diagramContainerIDToCrlDiagram[containerId]
}

// GetCrlDiagramIDForContainerID returns the identifier of the CrlDiagram being displayed in the container
func (dmPtr *DiagramManager) GetCrlDiagramIDForContainerID(containerID string) string {
	diagram := dmPtr.diagramContainerIDToCrlDiagram[containerID]
	id := ""
	if diagram != nil {
		id = diagram.GetId(CrlEditorSingleton.getHeldLocks())
	}
	return id
}

// DisplayDiagram display diagram creates a diagram tab, if one does not exist, and then displays the diagram.
// If the tab already exists, it simply shows it.
func (dmPtr *DiagramManager) DisplayDiagram(diagram core.Element, hl *core.HeldLocks) {
	diagramID := diagram.GetId(hl)
	diagramIDString := diagramID
	diagramLabel := core.GetLabel(diagram, hl)
	httpDiagramContainer := dmPtr.crlDiagramIDToDiagramContainer[diagramIDString]
	// Construct the container if it is not already present
	if httpDiagramContainer == nil {

		// Tracing
		//		if core.AdHocTrace == true {
		//			js.Global.Get("console").Call("log", "Creating httpDiagramContainer")
		//		}

		httpDiagramContainerID := createDiagramContainerPrefix() + diagramIDString
		topContent := js.Global.Get("top-content")
		httpDiagramContainer = js.Global.Get("document").Call("createElement", "DIV")
		httpDiagramContainer.Set("id", httpDiagramContainerID)
		httpDiagramContainer.Call("setAttribute", "class", "crlDiagramContainer")
		// It is not clear why, but the ondrop callback does not get called unless the ondragover callback is used,
		// even though the callback just calls preventDefault on the dragover event
		httpDiagramContainer.Set("ondragover", onDragover)
		httpDiagramContainer.Set("ondrop", onDiagramManagerDrop)
		httpDiagramContainer.Get("style").Set("display", "none")
		dmPtr.crlDiagramIDToDiagramContainer[diagramIDString] = httpDiagramContainer
		dmPtr.diagramContainerIDToCrlDiagram[httpDiagramContainerID] = diagram
		topContent.Call("appendChild", httpDiagramContainer)
		// Create the new tab
		tabs := js.Global.Get("tabs")
		newTab := js.Global.Get("document").Call("createElement", "button")
		newTab.Set("innerHTML", diagramLabel)
		newTab.Set("className", "w3-bar-item w3-button")
		newTabID := createDiagramTabPrefix() + diagramIDString
		newTab.Set("id", newTabID)
		newTab.Call("setAttribute", "viewId", httpDiagramContainerID)
		newTab.Call("addEventListener", "click", func(e jquery.Event) {
			onMakeDiagramVisible(e)
		})
		tabs.Call("appendChild", newTab, -1)

		jointGraph := dmPtr.crlDiagramIDToJointGraph[httpDiagramContainerID]
		if jointGraph == nil {
			jointGraphID := createJointGraphPrefix() + diagramIDString
			jointGraph = js.Global.Get("joint").Get("dia").Get("Graph").New()
			jointGraph.Set("id", jointGraphID)
			dmPtr.crlDiagramIDToJointGraph[httpDiagramContainerID] = jointGraph
			dmPtr.jointGraphIDToCrlDiagram[jointGraphID] = diagram
		}

		jointPaper := dmPtr.jointPapers[httpDiagramContainerID]
		if jointPaper == nil {
			diagramPaperDiv := js.Global.Get("document").Call("createElement", "DIV")
			httpDiagramContainer.Call("appendChild", diagramPaperDiv)
			jointPaper = js.Global.Get("joint").Get("dia").Get("Paper").New(js.M{
				"el":       []*js.Object{diagramPaperDiv},
				"width":    1000,
				"height":   1000,
				"model":    jointGraph,
				"gridSize": 1})
			dmPtr.jointPapers[httpDiagramContainerID] = jointPaper
			jointPaper.Call("on", "cell:pointerdown", onDiagramManagerCellPointerDown)
		}
	}
	makeDiagramVisible(httpDiagramContainer.Get("id").String())
}

func (dmPtr *DiagramManager) NewDiagram() core.Element {
	// Insert name prompt here
	name := getDefaultDiagramLabel()
	hl := CrlEditorSingleton.getHeldLocks()
	defer hl.ReleaseLocks()
	uOfD := CrlEditorSingleton.GetUofD()
	diagram, err := core.CreateReplicateAsRefinementFromUri(uOfD, crlDiagram.CrlDiagramUri, hl)
	if err != nil {
		log.Print(err)
	}
	core.SetLabel(diagram, name, hl)
	dmPtr.crlDiagramUUIDToCrlDiagram[diagram.GetId(hl)] = diagram

	// Tracing
	//	if core.AdHocTrace == true {
	//		log.Printf("Created diagram with name: %s", name)
	//	}

	dmPtr.DisplayDiagram(diagram, hl)
	return diagram
}

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
