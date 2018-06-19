package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/coreDiagram"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
)

const diagramContainerSuffix = "DiagramContainer"
const diagramSuffix = "Diagram"

var defaultLabelCount int
var diagramTabCount int
var diagramViewCount int
var diagramGraphCount int

type paperProperties struct {
	*js.Object
	el       []*js.Object `js:"el"`
	width    float64      `js:"width"`
	height   float64      `js:"height"`
	model    *js.Object   `js:"model"`
	gridSize float64      `js:"gridSize"`
}

type positionProperties struct {
	*js.Object
	x float64 `js:"x"`
	y float64 `js:"y"`
}

type sizeProperties struct {
	*js.Object
	width  float64 `js:"width"`
	height float64 `js:"height"`
}

type shapeProperties struct {
	*js.Object
	fill string `js:"fill"`
}

type textProperties struct {
	*js.Object
	text string `js:"text"`
	fill string `js:"fill"`
}

type attrProperties struct {
	*js.Object
	rect *shapeProperties `js:"rect"`
	text *textProperties  `js:"text"`
}

type rectProperties struct {
	*js.Object
	position *positionProperties `js:"position"`
	size     *sizeProperties     `js:"size"`
	attrs    *attrProperties     `js:"attrs"`
}

type baseElementProperties struct {
	*js.Object
	position *positionProperties `js:"position"`
	size     *sizeProperties     `js:"size"`
	attrs    *attrProperties     `js:"attrs"`
	name     string              `js:"name"`
}

type nameProperty struct {
	*js.Object
	name string `js:"name"`
}

type DiagramManager struct {
	abstractDiagram           core.Element
	diagrams                  map[uuid.UUID]core.Element
	diagramGraphs             map[string]*js.Object
	diagramFromDiagramGraphId map[string]core.Element
	diagramPapers             map[string]*js.Object
}

func NewDiagramManager() *DiagramManager {
	var diagramManager DiagramManager
	diagramManager.abstractDiagram = CrlEditorSingleton.uOfD.GetElementWithUri(coreDiagram.CrlDiagramUri)
	diagramManager.diagrams = make(map[uuid.UUID]core.Element)
	diagramManager.diagramGraphs = make(map[string]*js.Object)
	diagramManager.diagramFromDiagramGraphId = make(map[string]core.Element)
	diagramManager.diagramPapers = make(map[string]*js.Object)
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

func createDiagramViewPrefix() string {
	diagramViewCount++
	countString := strconv.Itoa(diagramViewCount)
	return "DiagramView" + countString
}

func createDiagramGraphPrefix() string {
	diagramGraphCount++
	countString := strconv.Itoa(diagramGraphCount)
	return "DiagramGraph" + countString
}

func (dmPtr *DiagramManager) DisplayDiagram(diagram core.Element, hl *core.HeldLocks) {
	diagramId := diagram.GetId(hl)
	diagramIdString := diagramId.String()
	diagramLabel := core.GetLabel(diagram, hl)
	diagramViewId := createDiagramViewPrefix() + diagramIdString

	// See if diagram is already in GUI

	topContent := js.Global.Get("top-content")
	crlDiagramContainer := js.Global.Get("document").Call("createElement", "DIV")
	crlDiagramContainer.Set("id", diagramViewId)
	crlDiagramContainer.Call("setAttribute", "class", "crlDiagramContainer")
	// It is not clear why, but the ondrop callback does not get called unless the ondragover callback is used,
	// even though the callback just calls preventDefault on the dragover event
	crlDiagramContainer.Set("ondragover", onDragover)
	crlDiagramContainer.Set("ondrop", onDiagramManagerDrop)
	crlDiagramContainer.Get("style").Set("display", "none")
	topContent.Call("appendChild", crlDiagramContainer)

	tabs := js.Global.Get("tabs")
	newTab := js.Global.Get("document").Call("createElement", "button")
	newTab.Set("innerHTML", diagramLabel)
	newTab.Set("className", "w3-bar-item w3-button")
	newTabId := createDiagramTabPrefix() + diagramIdString
	newTab.Set("id", newTabId)
	newTab.Call("setAttribute", "viewId", diagramViewId)
	//	newTab.Set("onclick", "openDiagramContainer('"+diagramViewId+"')")
	newTab.Call("addEventListener", "click", func(e jquery.Event) {
		onMakeDiagramVisible(e)
	})
	tabs.Call("appendChild", newTab, -1)

	diagramGraph := dmPtr.diagramGraphs[diagramViewId]
	if diagramGraph == nil {
		diagramGraphId := createDiagramGraphPrefix() + diagramIdString
		diagramGraph = js.Global.Get("joint").Get("dia").Get("Graph").New()
		diagramGraph.Set("id", diagramGraphId)
		dmPtr.diagramGraphs[diagramViewId] = diagramGraph
		dmPtr.diagramFromDiagramGraphId[diagramGraphId] = diagram
	}

	diagramPaper := dmPtr.diagramPapers[diagramViewId]
	if diagramPaper == nil {
		diagramPaperDiv := js.Global.Get("document").Call("createElement", "DIV")
		crlDiagramContainer.Call("appendChild", diagramPaperDiv)
		pProps := &paperProperties{Object: js.Global.Get("Object").New()}
		pProps.el = []*js.Object{diagramPaperDiv}
		pProps.width = 600
		pProps.height = 600
		pProps.model = diagramGraph
		pProps.gridSize = 1
		js.Global.Set("pProps", pProps)
		diagramPaper = js.Global.Get("joint").Get("dia").Get("Paper").New(pProps)
		dmPtr.diagramPapers[diagramViewId] = diagramPaper
		diagramPaper.Call("on", "cell:pointerdown", onDiagramManagerCellPointerDown)
	}

	js.Global.Set("diagramGraph", diagramGraph)
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
	diagram, err := core.CreateReplicateAsRefinementFromUri(uOfD, coreDiagram.CrlDiagramUri, hl)
	if err != nil {
		log.Print(err)
	}
	core.SetLabel(diagram, name, hl)
	dmPtr.diagrams[diagram.GetId(hl)] = diagram
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
	diagramManager := CrlEditorSingleton.GetDiagramManager()

	diagramViewId := event.Get("target").Get("parentElement").Get("parentElement").Get("id").String()
	js.Global.Set("dropTarget", event.Get("target"))
	js.Global.Set("dropTargetParent", event.Get("target").Get("parentElement"))
	js.Global.Set("dropTargetParentId", event.Get("target").Get("parentElement").Get("id"))
	js.Global.Get("console").Call("log", "DiagramViewId: "+diagramViewId)
	graph := diagramManager.diagramGraphs[diagramViewId]
	//	diagram := diagramManager.diagramFromDiagramGraphId[graph.Get("id").String()]

	diagramBaseElementProps := &baseElementProperties{Object: js.Global.Get("Object").New()}

	// size
	sizeProp := &sizeProperties{Object: js.Global.Get("Object").New()}
	sizeProp.width = 100.0
	sizeProp.height = 30.0
	diagramBaseElementProps.size = sizeProp

	// position
	positionProp := &positionProperties{Object: js.Global.Get("Object").New()}
	positionProp.x = event.Get("layerX").Float()
	positionProp.y = event.Get("layerY").Float()
	diagramBaseElementProps.position = positionProp

	diagramBaseElement := js.Global.Get("joint").Get("shapes").Get("crl").Get("BaseElement").New(diagramBaseElementProps)

	// name
	be := CrlEditorSingleton.GetTreeDragSelection()
	name := core.GetLabel(be, hl)
	//	nameProps := &nameProperty{Object: js.Global.Get("Object").New()}
	//	nameProps.name = name
	diagramBaseElement.Get("attributes").Set("name", name)
	diagramBaseElement.Call("updateRectangles")

	js.Global.Set("graph", graph)
	js.Global.Set("diagramBaseElement", diagramBaseElement)

	graph.Call("addCell", diagramBaseElement)
	CrlEditorSingleton.SelectBaseElement(be)
	CrlEditorSingleton.SetTreeDragSelection(nil)
}

func onDiagramManagerCellPointerDown(cellView *js.Object, event *js.Object, x *js.Object, y *js.Object) {
	baseElementIdString := cellView.Get("model").Get("id").String()
	log.Printf("Pointerdown on Cell %s", baseElementIdString)
	js.Global.Set("cellView", cellView)
	CrlEditorSingleton.SelectBaseElementUsingIdString(baseElementIdString)
}

func onMakeDiagramVisible(e jquery.Event) {
	diagramViewId := e.Get("target").Call("getAttribute", "viewId").String()
	js.Global.Get("console").Call("log", "In : onMakeDiagramVisible with: "+diagramViewId)
	js.Global.Set("clickEvent", e)
	js.Global.Set("clickEventTarget", e.Get("target"))
	js.Global.Set("clickEventViewId", e.Get("target").Call("getAttribute", "viewId"))
	x := js.Global.Get("document").Call("getElementsByClassName", "crlDiagramContainer")
	lengthString := strconv.Itoa(x.Length())
	js.Global.Get("console").Call("log", "List length: "+lengthString)
	for i := 0; i < x.Length(); i++ {
		js.Global.Get("console").Call("log", "Container id: ", x.Index(i).Get("id").String())
		if x.Index(i).Get("id").String() == diagramViewId {
			x.Index(i).Get("style").Set("display", "block")
			js.Global.Get("console").Call("log", "Showing: "+diagramViewId)
		} else {
			x.Index(i).Get("style").Set("display", "none")
			js.Global.Get("console").Call("log", "Hiding: "+diagramViewId)
		}
	}
}

func (dmPtr *DiagramManager) SetSize() {
	// A test rectangle
	posProp := &positionProperties{Object: js.Global.Get("Object").New()}
	posProp.x = 100
	posProp.y = 30

	sizeProp := &sizeProperties{Object: js.Global.Get("Object").New()}
	sizeProp.width = 100.0
	sizeProp.height = 30.0

	shapeProp := &shapeProperties{Object: js.Global.Get("Object").New()}
	shapeProp.fill = "blue"

	textProp := &textProperties{Object: js.Global.Get("Object").New()}
	textProp.text = "my box"
	textProp.fill = "white"

	attrProp := &attrProperties{Object: js.Global.Get("Object").New()}
	attrProp.rect = shapeProp
	attrProp.text = textProp

	rectProp := &rectProperties{Object: js.Global.Get("Object").New()}
	rectProp.position = posProp
	rectProp.size = sizeProp
	rectProp.attrs = attrProp

	rect := js.Global.Get("joint").Get("shapes").Get("basic").Get("Rect").New(rectProp)

	js.Global.Set("diagramRect", rect)
	js.Global.Set("rectProp", rectProp)
	js.Global.Set("sizeProp", sizeProp)

	js.Global.Get("diagramGraph").Call("addCell", rect)

}
