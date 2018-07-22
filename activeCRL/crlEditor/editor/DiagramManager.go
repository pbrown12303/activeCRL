package editor

import (
	//	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/crlDiagram"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
)

const diagramContainerSuffix = "DiagramContainer"
const diagramSuffix = "Diagram"

var defaultLabelCount int
var diagramTabCount int
var diagramContainerCount int
var jointGraphCount int
var jointBaseElementNodeCount int

type DiagramManager struct {
	abstractDiagram                  core.Element
	crlDiagramUUIDToCrlDiagram       map[uuid.UUID]core.Element
	crlDiagramIdToDiagramContainer   map[string]*js.Object
	crlDiagramIdToJointGraph         map[string]*js.Object
	crlDiagramNodeUUIDToJointElement map[uuid.UUID]*js.Object
	diagramContainerIdToCrlDiagram   map[string]core.Element
	jointElementIdToCrlDiagramNode   map[string]core.Element
	jointGraphIdToCrlDiagram         map[string]core.Element
	jointPapers                      map[string]*js.Object
}

func NewDiagramManager() *DiagramManager {
	var diagramManager DiagramManager
	diagramManager.abstractDiagram = CrlEditorSingleton.uOfD.GetElementWithUri(crlDiagram.CrlDiagramUri)
	diagramManager.crlDiagramUUIDToCrlDiagram = make(map[uuid.UUID]core.Element)
	diagramManager.crlDiagramIdToDiagramContainer = make(map[string]*js.Object)
	diagramManager.crlDiagramIdToJointGraph = make(map[string]*js.Object)
	diagramManager.crlDiagramNodeUUIDToJointElement = make(map[uuid.UUID]*js.Object)
	diagramManager.diagramContainerIdToCrlDiagram = make(map[string]core.Element)
	diagramManager.jointElementIdToCrlDiagramNode = make(map[string]core.Element)
	diagramManager.jointGraphIdToCrlDiagram = make(map[string]core.Element)
	diagramManager.jointPapers = make(map[string]*js.Object)
	defineNodeViews()
	return &diagramManager
}

func (dm *DiagramManager) addViewFunctionsToUofD() {
	uOfD := CrlEditorSingleton.uOfD
	hl := CrlEditorSingleton.hl
	addDiagramViewFunctionsToUofD(uOfD, hl)
	//	addDiagramNodeViewFunctionsToUofD(uOfD, hl)
	//	addDiagramLinkViewFunctionsToUofD(uOfD, hl)
}

func createDiagramTabPrefix() string {
	diagramTabCount++
	countString := strconv.Itoa(diagramTabCount)
	return "DiagramTab" + countString
}

func createDiagramContainerPrefix() string {
	diagramContainerCount++
	countString := strconv.Itoa(diagramContainerCount)
	return "DiagramView" + countString
}

func createJointBaseElementNodePrefix() string {
	jointBaseElementNodeCount++
	countString := strconv.Itoa(jointGraphCount)
	return "JointBaseElementNode" + countString
}

func createJointGraphPrefix() string {
	jointGraphCount++
	countString := strconv.Itoa(jointGraphCount)
	return "DiagramGraph" + countString
}

func (dmPtr *DiagramManager) DisplayDiagram(diagram core.Element, hl *core.HeldLocks) {
	diagramId := diagram.GetId(hl)
	diagramIdString := diagramId.String()
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

func getDefaultDiagramLabel() string {
	defaultLabelCount++
	countString := strconv.Itoa(defaultLabelCount)
	return "Diagram" + countString
}

func (dmPtr *DiagramManager) NewDiagram() core.Element {
	// Insert name prompt here
	name := getDefaultDiagramLabel()
	hl := CrlEditorSingleton.hl
	defer hl.ReleaseLocks()
	uOfD := CrlEditorSingleton.uOfD
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

func onDragover(event *js.Object, data *js.Object) {
	event.Call("preventDefault")
}

func onDiagramManagerDrop(event *js.Object) {
	hl := CrlEditorSingleton.hl
	event.Call("preventDefault")
	js.Global.Set("dropEvent", event)
	log.Printf("On Drop called")
	httpDiagramContainerId := event.Get("target").Get("parentElement").Get("parentElement").Get("id").String()
	be := CrlEditorSingleton.GetTreeDragSelection()
	js.Global.Set("diagramDroppedBaseElement", be.GetId(hl).String())
	js.Global.Get("console").Call("log", "In onDiagramManagerDrop")

	addNodeView(httpDiagramContainerId, be, event.Get("layerX").Float(), event.Get("layerY").Float(), hl)

	CrlEditorSingleton.SelectBaseElement(be)
	CrlEditorSingleton.SetTreeDragSelection(nil)
}

func onDiagramManagerCellPointerDown(cellView *js.Object, event *js.Object, x *js.Object, y *js.Object) {
	js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown")
	hl := CrlEditorSingleton.hl
	crlJointId := cellView.Get("model").Get("crlJointId").String()
	js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown crlJointId = "+crlJointId)
	diagramManager := CrlEditorSingleton.diagramManager
	diagramNode := diagramManager.jointElementIdToCrlDiagramNode[crlJointId]
	if diagramNode == nil {
		js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown diagramNode is nil")
	} else {
		js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown diagramNode id = "+diagramNode.GetId(hl).String())
	}

	be, err := crlDiagram.GetReferencedBaseElement(diagramNode, hl)
	if err == nil {
		CrlEditorSingleton.SelectBaseElement(be)
	}
}

func onMakeDiagramVisible(e jquery.Event) {
	httpDiagramContainerId := e.Get("target").Call("getAttribute", "viewId").String()
	//	js.Global.Get("console").Call("log", "In : onMakeDiagramVisible with: "+httpDiagramContainerId)
	//	js.Global.Set("clickEvent", e)
	//	js.Global.Set("clickEventTarget", e.Get("target"))
	//	js.Global.Set("clickEventViewId", e.Get("target").Call("getAttribute", "viewId"))
	makeDiagramVisible(httpDiagramContainerId)
}

func makeDiagramVisible(httpDiagramContainerId string) {
	x := js.Global.Get("document").Call("getElementsByClassName", "crlDiagramContainer")
	lengthString := strconv.Itoa(x.Length())
	js.Global.Get("console").Call("log", "List length: "+lengthString)
	for i := 0; i < x.Length(); i++ {
		js.Global.Get("console").Call("log", "Container id: ", x.Index(i).Get("id").String())
		if x.Index(i).Get("id").String() == httpDiagramContainerId {
			x.Index(i).Get("style").Set("display", "block")
			js.Global.Get("console").Call("log", "Showing: "+httpDiagramContainerId)
		} else {
			x.Index(i).Get("style").Set("display", "none")
			js.Global.Get("console").Call("log", "Hiding: "+httpDiagramContainerId)
		}
	}

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
