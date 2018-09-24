package editor

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/crlDiagram"
)

// EditorURI is the URI for accessing the CrlEditor
var EditorURI = "http://activeCrl.com/crlEditor/Editor"

// CrlEditorSingleton is the singleton instance of the CrlEditor
var CrlEditorSingleton CrlEditor

var debugSettingsDialog jquery.JQuery
var displayGraphDialog jquery.JQuery

type crlEditor struct {
	currentSelection  core.BaseElement
	treeDragSelection core.BaseElement
	diagramManager    *DiagramManager
	propertiesManager *PropertiesManager
	treeManager       *TreeManager
	hl                *core.HeldLocks
	uOfD              core.UniverseOfDiscourse
	initialized       bool
}

// InitializeCrlEditorSingleton initializes the CrlEditor singleton instance. It should be called once
// when the editor web page is created
func InitializeCrlEditorSingleton() {
	var wg sync.WaitGroup
	var editor crlEditor
	editor.initialized = false
	editor.hl = core.NewHeldLocks(&wg)
	defer editor.hl.ReleaseLocks()

	// Set the value of the singleton
	CrlEditorSingleton = &editor

	// Set up the Universe of Discourse
	editor.uOfD = core.NewUniverseOfDiscourse(editor.hl)
	crlDiagram.AddCoreDiagramToUofD(editor.uOfD, editor.hl)
	AddEditorViewsToUofD(editor.uOfD, editor.hl)

	// Set up the diagram manager
	editor.diagramManager = newDiagramManager()

	// Set up the tree manager
	editor.treeManager = NewTreeManager(editor.uOfD, "#uOfD", editor.hl)
	editor.treeManager.InitializeTree(editor.hl)

	// Set up the properties manager
	editor.propertiesManager = NewPropertiesManager()

	// Add the event listeners
	editorQuery := jquery.NewJQuery("body")
	editorQuery.On("ondrop", func(e jquery.Event, data *js.Object) {
		onEditorDrop(e, data)
	})

	editor.hl.ReleaseLocksAndWait()
	editor.initialized = true
	log.Printf("Editor initialized")
}

// AddEditorViewsToUofD adds the concepts representing the various editor views to the universe of discurse
func AddEditorViewsToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	conceptSpace := uOfD.GetElementWithUri(EditorURI)
	if conceptSpace == nil {
		conceptSpace = BuildEditorConceptSpace(uOfD, hl)
		if conceptSpace == nil {
			log.Printf("Build of Editor Concept Space failed")
		}
	}
	return conceptSpace
}

// DebugSettingsOK is the callback function for the Debug Settings dialog OK button.
// It updaates the debug settings with the values from the dialog
func DebugSettingsOK(e jquery.Event, data *js.Object) {
	//	log.Printf("DebugSettingsOK called")
	//	js.Global.Set("maxTracingDepth", debugSettingsDialog.Find("#maxTracingDepth"))
	maxTracingDepth, err1 := strconv.Atoi(debugSettingsDialog.Find("#maxTracingDepth").Val())
	if err1 != nil {
		log.Printf(err1.Error())
		return
	}
	//	js.Global.Set("enableTracing", debugSettingsDialog.Find("#enableTracing"))
	enableTracing, err2 := strconv.ParseBool(debugSettingsDialog.Find("#enableTracing").Val())
	if err2 != nil {
		log.Printf(err2.Error())
		return
	}
	//	log.Printf("Debug Settings depth: %d enabled: %t \n", maxTracingDepth, enableTracing)
	//	js.Global.Set("debugSettingsDialog", debugSettingsDialog)
	core.TraceChange = enableTracing
	if enableTracing == true {
		core.SetNotificationsLimit(maxTracingDepth)
	} else {
		core.SetNotificationsLimit(0)
	}
	debugSettingsDialog.Call("dialog", "close")
}

// DebugSettings creates and displays the Debug Settings dialog so that the debug settings can be updated from the UI
func (edPtr *crlEditor) DebugSettings() {
	//	log.Printf("DebugSettings called")
	if jquery.IsEmptyObject(debugSettingsDialog) {
		debugSettingsDialog = jquery.NewJQuery("<div class='dialog' title='Notification Trace Settings'></div>").Call("dialog", js.M{
			"resizable": false,
			"height":    200,
			"modal":     true,
			"buttons":   js.M{"OK": DebugSettingsOK}})
		debugSettingsDialog.SetHtml("<label for='maxTracingDepth'>Max Depth</label>" +
			"<input type='number' id='maxTracingDepth' placeholder='1'> <br>" +
			"<label for='enableTracing'>Enable Notification Tracing</label>" +
			"<input type='checkbox' id='enableTracing' value='1'> <br>")
	}
	//	js.Global.Set("newDialog", debugSettingsDialog)
	debugSettingsDialog.Call("dialog", "open")
	//	jquery.NewJQuery("#notificationTraceSettingsDialog").Call("dialog", "open")
}

// DisplayGraph opens a new tab and displays the selected graph
func DisplayGraph(e jquery.Event, data *js.Object) {
	selectedGraphIndexString := displayGraphDialog.Find("#selectedGraph").Val()
	selectedGraphIndex, err := strconv.Atoi(selectedGraphIndexString)
	if err != nil {
		log.Printf(err.Error())
	}
	displayGraphDialog.Call("dialog", "close")
	graphs := core.GetFunctionCallGraphs()
	log.Printf("Number of function call graphs: %d\n", len(graphs))
	log.Printf("Graphs: %#v\n", graphs)
	if selectedGraphIndex > 0 && selectedGraphIndex <= len(graphs) {
		newTab := js.Global.Get("window").Call("open", "", "_blank")
		log.Printf("Selected Graph: %#v\n", graphs[selectedGraphIndex-1])
		graphString := graphs[selectedGraphIndex-1].GetGraph().String()
		log.Printf("Graph String: %s\n", graphString)
		graphStringEscapedQuotes := strings.Replace(graphString, "\"", "\\\"", -1)
		graphStringEscapedQuotesNoNL := strings.Replace(graphStringEscapedQuotes, "\n", "", -1)
		htmlString := "<head>" +
			"  <meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\">\n" +
			"  <meta charset=\"utf-8\">\n" +
			"  <title>Function Call Graph Display</title>\n" +
			"  <link rel=\"stylesheet\" href=\"/js/jquery-ui-1.12.1.custom/jquery-ui.css\">\n" +
			"  <script src=\"/js/jquery-ui-1.12.1.custom/external/jquery/jquery.js\"></script>\n" +
			"  <script src=\"/js/viz.js\"></script>\n" +
			"  <script src=\"/js/full.render.js\"></script>\n" +
			"</head>\n" +
			"<body>\n" +
			"  <script>\n" +
			"    var graphString =\"" + graphStringEscapedQuotesNoNL + "\"\n" +
			"    var viz = new Viz();\n" +
			"    viz.renderSVGElement(graphString)\n" +
			"    .then(function(element) {\n" +
			"      document.body.appendChild(element);\n" +
			"    })\n" +
			"    .catch(error => {\n" +
			"      // Create a new Viz instance (@see Caveats page for more info) \n" +
			"      viz = new Viz();\n" +
			"      // Possibly display the error\n" +
			"      console.error(error);\n" +
			"    });\n" +
			"  </script>\n" +
			"</body>\n"
		js.Global.Set("htmlString", htmlString)
		newTab.Get("document").Call("write", htmlString)
	}
}

// DisplayGraphDialog opens a dialog in which a trace can be selected for display
func (edPtr *crlEditor) DisplayGraphDialog() {
	if jquery.IsEmptyObject(displayGraphDialog) {
		displayGraphDialog = jquery.NewJQuery("<div class='dialog' title='Display Function Call Graph'></div>").Call("dialog", js.M{
			"resizable": false,
			"height":    300,
			"width":     400,
			"modal":     true,
			"buttons":   js.M{"DisplaySelectedTrace": DisplayGraph}})
	}
	displayGraphDialog.SetHtml("<p id=\"numberOfAvailableFunctionCallGraphs\">Number of available function calls: " + strconv.Itoa(len(core.GetFunctionCallGraphs())) + " </output><br>" +
		"<label for=\"selectedGraph\">Graph To Display</label>\n" +
		"<input type=\"number\" id=\"selectedGraph\" placeholder=\"0\" min=\"0\" max=\"CrlEditor.GetNumberOfFunctionCalls()\">")
	js.Global.Set("displayTraceDialog", displayGraphDialog)
	js.Global.Set("numberOfAvailableFunctionCallGraphs", displayGraphDialog.Find("#numberOfAvailableFunctionCallGraphs"))
	//	displayTraceDialog.Find("#numberOfAvailableFunctionCallGraphs").SetText("Number of available traces: " + strconv.Itoa(len(core.GetNotificationGraphs())))
	displayGraphDialog.Call("dialog", "open")
}

func (edPtr *crlEditor) GetAdHocTrace() bool {
	return core.AdHocTrace
}

func (edPtr *crlEditor) GetCurrentSelection() core.BaseElement {
	return edPtr.currentSelection
}

func (edPtr *crlEditor) GetCurrentSelectionID() string {
	return edPtr.currentSelection.GetId(edPtr.hl)
}

func (edPtr *crlEditor) GetDiagramManager() *DiagramManager {
	return edPtr.diagramManager
}

func (edPtr *crlEditor) getHeldLocks() *core.HeldLocks {
	return edPtr.hl
}

func (edPtr *crlEditor) GetNumberOfFunctionCalls() int {
	return len(core.GetFunctionCallGraphs())
}

func (edPtr *crlEditor) GetPropertiesManager() *PropertiesManager {
	return edPtr.propertiesManager
}

func (edPtr *crlEditor) GetTraceChange() bool {
	return core.TraceChange
}

func (edPtr *crlEditor) GetTraceChangeLimit() int {
	return core.GetNotificationsLimit()
}

func (edPtr *crlEditor) GetTreeDragSelection() core.BaseElement {
	return edPtr.treeDragSelection
}

func (edPtr *crlEditor) GetTreeDragSelectionID() string {
	return edPtr.treeDragSelection.GetId(edPtr.hl)
}

func (edPtr *crlEditor) GetTreeManager() *TreeManager {
	return edPtr.treeManager
}

func (edPtr *crlEditor) GetUofD() core.UniverseOfDiscourse {
	return edPtr.uOfD
}

func (edPtr *crlEditor) IsInitialized() bool {
	return edPtr.initialized
}

func (edPtr *crlEditor) RunTest() string {
	log.Printf("In RunTest \n")
	currentSelection := edPtr.GetCurrentSelection()
	if currentSelection == nil {
		log.Printf("In RunTest, current selection is nil \n")
		return ""
	}
	i := 100000
	startTime := time.Now()
	for i > 0 {
		edPtr.hl.LockBaseElement(currentSelection)
		i--
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("RunTest duration: %s \n", duration.String())
	return duration.String()
}

func (edPtr *crlEditor) SelectBaseElement(be core.BaseElement) core.BaseElement {
	edPtr.currentSelection = be
	edPtr.propertiesManager.BaseElementSelected(edPtr.currentSelection, edPtr.hl)
	return edPtr.currentSelection
}

func (edPtr *crlEditor) SelectBaseElementUsingIDString(id string) core.BaseElement {
	edPtr.currentSelection = edPtr.uOfD.GetBaseElement(id)
	edPtr.propertiesManager.BaseElementSelected(edPtr.currentSelection, edPtr.hl)
	return edPtr.currentSelection
}

func (edPtr *crlEditor) SetAdHocTrace(status bool) {
	core.AdHocTrace = status
}

func (edPtr *crlEditor) SetSelectionDefinition(definition string) {
	switch edPtr.currentSelection.(type) {
	case core.Element:
		core.SetDefinition(edPtr.currentSelection.(core.Element), definition, edPtr.hl)
		edPtr.hl.ReleaseLocksAndWait()
	}
}

func (edPtr *crlEditor) SetSelectionLabel(name string) {
	switch edPtr.currentSelection.(type) {
	case core.Element:
		core.SetLabel(edPtr.currentSelection.(core.Element), name, edPtr.hl)
		edPtr.hl.ReleaseLocksAndWait()
	}
}

func (edPtr *crlEditor) SetSelectionURI(uri string) {
	core.SetUri(edPtr.currentSelection, uri, edPtr.hl)
	edPtr.hl.ReleaseLocksAndWait()
}

func (edPtr *crlEditor) SetTraceChange(newValue bool) {
	core.TraceChange = newValue
}

func (edPtr *crlEditor) SetTraceChangeLimit(limit int) {
	core.SetNotificationsLimit(limit)
}

func (edPtr *crlEditor) SetTreeDragSelection(be core.BaseElement) {
	edPtr.treeDragSelection = be
}

// BuildEditorConceptSpace programmatically constructs the EditorConceptSpace
func BuildEditorConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// EditorConceptSpace
	conceptSpace := uOfD.NewElement(hl, EditorURI)
	core.SetLabel(conceptSpace, "EditorConceptSpace", hl)
	core.SetUri(conceptSpace, EditorURI, hl)

	BuildTreeViews(conceptSpace, hl)

	return conceptSpace
}

func onEditorDrop(e jquery.Event, data *js.Object) {
	CrlEditorSingleton.SetTreeDragSelection(nil)
}

func init() {
	registerTreeViewFunctions()
}

// CrlEditor type is the central component of the CrlEditor. It manages the subordinate managers (Property, Tree, Diagram)
// and the singleton instances of the Universe of Discourse and HeldLocks shared by all editing operations.
type CrlEditor interface {
	DebugSettings()
	DisplayGraphDialog()
	GetAdHocTrace() bool
	GetCurrentSelection() core.BaseElement
	GetCurrentSelectionID() string
	GetDiagramManager() *DiagramManager
	getHeldLocks() *core.HeldLocks
	GetPropertiesManager() *PropertiesManager
	GetTraceChange() bool
	GetTraceChangeLimit() int
	GetTreeDragSelection() core.BaseElement
	GetTreeDragSelectionID() string
	GetTreeManager() *TreeManager
	GetUofD() core.UniverseOfDiscourse
	IsInitialized() bool
	RunTest() string
	SelectBaseElement(be core.BaseElement) core.BaseElement
	SelectBaseElementUsingIDString(id string) core.BaseElement
	SetAdHocTrace(bool)
	SetSelectionDefinition(definition string)
	SetSelectionLabel(name string)
	SetSelectionURI(uri string)
	SetTraceChange(bool)
	SetTraceChangeLimit(int)
	SetTreeDragSelection(be core.BaseElement)
}
