package editor

import (
	//	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/crlDiagram"
	"log"
)

const diagramContainerSuffix = "DiagramContainer"
const diagramSuffix = "Diagram"

type DiagramManager struct {
	abstractDiagram                  core.Element
	crlDiagramUUIDToCrlDiagram       map[string]core.Element
	crlDiagramIdToDiagramContainer   map[string]*js.Object
	crlDiagramIdToJointGraph         map[string]*js.Object
	crlDiagramNodeUUIDToJointElement map[string]*js.Object
	diagramContainerIdToCrlDiagram   map[string]core.Element
	jointElementIdToCrlDiagramNode   map[string]core.Element
	jointGraphIdToCrlDiagram         map[string]core.Element
	jointPapers                      map[string]*js.Object
}

func NewDiagramManager() *DiagramManager {
	var diagramManager DiagramManager
	diagramManager.abstractDiagram = CrlEditorSingleton.GetUofD().GetElementWithUri(crlDiagram.CrlDiagramUri)
	diagramManager.crlDiagramUUIDToCrlDiagram = make(map[string]core.Element)
	diagramManager.crlDiagramIdToDiagramContainer = make(map[string]*js.Object)
	diagramManager.crlDiagramIdToJointGraph = make(map[string]*js.Object)
	diagramManager.crlDiagramNodeUUIDToJointElement = make(map[string]*js.Object)
	diagramManager.diagramContainerIdToCrlDiagram = make(map[string]core.Element)
	diagramManager.jointElementIdToCrlDiagramNode = make(map[string]core.Element)
	diagramManager.jointGraphIdToCrlDiagram = make(map[string]core.Element)
	diagramManager.jointPapers = make(map[string]*js.Object)
	defineNodeViews()
	return &diagramManager
}

func (dm *DiagramManager) addViewFunctionsToUofD() {
	uOfD := CrlEditorSingleton.GetUofD()
	hl := CrlEditorSingleton.getHeldLocks()
	addDiagramViewFunctionsToUofD(uOfD, hl)
	//	addDiagramNodeViewFunctionsToUofD(uOfD, hl)
	//	addDiagramLinkViewFunctionsToUofD(uOfD, hl)
}

func (dmPtr *DiagramManager) GetCrlDiagramForContainerId(containerId string) core.Element {
	return dmPtr.diagramContainerIdToCrlDiagram[containerId]
}

func (dmPtr *DiagramManager) GetCrlDiagramIdForContainerId(containerId string) string {
	diagram := dmPtr.diagramContainerIdToCrlDiagram[containerId]
	id := ""
	if diagram != nil {
		id = diagram.GetId(CrlEditorSingleton.getHeldLocks())
	}
	return id
}

func (dmPtr *DiagramManager) DisplayDiagram(diagram core.Element, hl *core.HeldLocks) {
	diagramId := diagram.GetId(hl)
	diagramIdString := diagramId
	diagramLabel := core.GetLabel(diagram, hl)
	httpDiagramContainer := dmPtr.crlDiagramIdToDiagramContainer[diagramIdString]
	// Construct the container if it is not already present
	if httpDiagramContainer == nil {
		js.Global.Get("console").Call("log", "Creating httpDiagramContainer")
		httpDiagramContainerId := createDiagramContainerPrefix() + diagramIdString
		topContent := js.Global.Get("top-content")
		httpDiagramContainer = js.Global.Get("document").Call("createElement", "DIV")
		httpDiagramContainer.Set("id", httpDiagramContainerId)
		httpDiagramContainer.Call("setAttribute", "class", "crlDiagramContainer")
		// It is not clear why, but the ondrop callback does not get called unless the ondragover callback is used,
		// even though the callback just calls preventDefault on the dragover event
		httpDiagramContainer.Set("ondragover", onDragover)
		httpDiagramContainer.Set("ondrop", onDiagramManagerDrop)
		httpDiagramContainer.Get("style").Set("display", "none")
		dmPtr.crlDiagramIdToDiagramContainer[diagramIdString] = httpDiagramContainer
		dmPtr.diagramContainerIdToCrlDiagram[httpDiagramContainerId] = diagram
		topContent.Call("appendChild", httpDiagramContainer)
		// Create the new tab
		tabs := js.Global.Get("tabs")
		newTab := js.Global.Get("document").Call("createElement", "button")
		newTab.Set("innerHTML", diagramLabel)
		newTab.Set("className", "w3-bar-item w3-button")
		newTabId := createDiagramTabPrefix() + diagramIdString
		newTab.Set("id", newTabId)
		newTab.Call("setAttribute", "viewId", httpDiagramContainerId)
		newTab.Call("addEventListener", "click", func(e jquery.Event) {
			onMakeDiagramVisible(e)
		})
		tabs.Call("appendChild", newTab, -1)

		jointGraph := dmPtr.crlDiagramIdToJointGraph[httpDiagramContainerId]
		if jointGraph == nil {
			jointGraphId := createJointGraphPrefix() + diagramIdString
			jointGraph = js.Global.Get("joint").Get("dia").Get("Graph").New()
			jointGraph.Set("id", jointGraphId)
			dmPtr.crlDiagramIdToJointGraph[httpDiagramContainerId] = jointGraph
			dmPtr.jointGraphIdToCrlDiagram[jointGraphId] = diagram
		}

		jointPaper := dmPtr.jointPapers[httpDiagramContainerId]
		if jointPaper == nil {
			diagramPaperDiv := js.Global.Get("document").Call("createElement", "DIV")
			httpDiagramContainer.Call("appendChild", diagramPaperDiv)
			jointPaper = js.Global.Get("joint").Get("dia").Get("Paper").New(js.M{
				"el":       []*js.Object{diagramPaperDiv},
				"width":    1000,
				"height":   1000,
				"model":    jointGraph,
				"gridSize": 1})
			dmPtr.jointPapers[httpDiagramContainerId] = jointPaper
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
	log.Printf("Created diagram with name: %s", name)
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
