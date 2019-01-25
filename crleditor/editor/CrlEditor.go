package editor

import (
	"log"
	"time"

	"github.com/pbrown12303/activeCRL/crldiagram"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/core"
)

// EditorURI is the URI for accessing the CrlEditor
var EditorURI = "http://activeCrl.com/crlEditor/Editor"

// CrlEditorSingleton is the singleton instance of the CrlEditor
var CrlEditorSingleton *CrlEditor

var debugSettingsDialog jquery.JQuery
var displayGraphDialog jquery.JQuery

// CrlEditor type is the central component of the CrlEditor. It manages the subordinate managers (Property, Tree, Diagram)
// and the singleton instances of the Universe of Discourse and HeldLocks shared by all editing operations.
type CrlEditor struct {
	clientNotificationManager *ClientNotificationManager
	currentSelection          core.Element
	cutBuffer                 map[string]core.Element
	diagramManager            *DiagramManager
	initialized               bool
	propertiesManager         *PropertiesManager
	treeDragSelection         core.Element
	treeManager               *TreeManager
	uOfD                      core.UniverseOfDiscourse
}

// InitializeCrlEditorSingleton initializes the CrlEditor singleton instance. It should be called once
// when the editor web page is created
func InitializeCrlEditorSingleton() {
	var editor CrlEditor
	editor.initialized = false
	editor.uOfD = core.NewUniverseOfDiscourse()
	hl := editor.uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()

	// Set the value of the singleton
	CrlEditorSingleton = &editor

	// Set up the Universe of Discourse
	crldiagram.BuildCrlDiagramConceptSpace(editor.uOfD, hl)
	hl.ReleaseLocksAndWait()
	AddEditorViewsToUofD(editor.uOfD, hl)
	hl.ReleaseLocksAndWait()

	editor.cutBuffer = make(map[string]core.Element)

	// Set up the diagram manager
	editor.diagramManager = newDiagramManager(&editor)

	// Set up the tree manager
	editor.treeManager = NewTreeManager(editor.uOfD, "#uOfD", hl)

	// Set up the properties manager
	editor.propertiesManager = NewPropertiesManager(&editor)

	// Create the ClientNotificationManager
	editor.clientNotificationManager = newClientNotificationManager()

	// TODO: Move this onDrop function to the client
	// Add the event listeners
	// editorQuery := jquery.NewJQuery("body")
	// editorQuery.On("ondrop", func(e jquery.Event, data *js.Object) {
	// 	onEditorDrop(e, data)
	// })
	hl.ReleaseLocksAndWait()
	registerTreeViewFunctions(editor.uOfD)
	editor.initialized = true
	log.Printf("Editor initialized")
}

// InitializeClient sets the client state after a browser refresh.
func InitializeClient() {
	<-webSocketReady
	log.Printf("InitializeClient executing")
	uOfD := CrlEditorSingleton.GetUofD()
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	if CrlEditorSingleton.IsInitialized() == false {
		time.Sleep(1000 * time.Millisecond)
	}
	CrlEditorSingleton.GetTreeManager().InitializeTree(hl)
}

// AddEditorViewsToUofD adds the concepts representing the various editor views to the universe of discurse
func AddEditorViewsToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	conceptSpace := uOfD.GetElementWithURI(EditorURI)
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
	// TODO: Reimplement DebugSettingsOK
	// //	log.Printf("DebugSettingsOK called")
	// //	js.Global.Set("maxTracingDepth", debugSettingsDialog.Find("#maxTracingDepth"))
	// maxTracingDepth, err1 := strconv.Atoi(debugSettingsDialog.Find("#maxTracingDepth").Val())
	// if err1 != nil {
	// 	log.Printf(err1.Error())
	// 	return
	// }
	// //	js.Global.Set("enableTracing", debugSettingsDialog.Find("#enableTracing"))
	// enableTracing, err2 := strconv.ParseBool(debugSettingsDialog.Find("#enableTracing").Val())
	// if err2 != nil {
	// 	log.Printf(err2.Error())
	// 	return
	// }
	// //	log.Printf("Debug Settings depth: %d enabled: %t \n", maxTracingDepth, enableTracing)
	// //	js.Global.Set("debugSettingsDialog", debugSettingsDialog)
	// core.TraceChange = enableTracing
	// if enableTracing == true {
	// 	core.SetNotificationsLimit(maxTracingDepth)
	// } else {
	// 	core.SetNotificationsLimit(0)
	// }
	// debugSettingsDialog.Call("dialog", "close")
}

// DebugSettings creates and displays the Debug Settings dialog so that the debug settings can be updated from the UI
func (edPtr *CrlEditor) DebugSettings() {
	// TODO: Reimplement DebugSettings
	// //	log.Printf("DebugSettings called")
	// if jquery.IsEmptyObject(debugSettingsDialog) {
	// 	debugSettingsDialog = jquery.NewJQuery("<div class='dialog' title='Notification Trace Settings'></div>").Call("dialog", js.M{
	// 		"resizable": false,
	// 		"height":    200,
	// 		"modal":     true,
	// 		"buttons":   js.M{"OK": DebugSettingsOK}})
	// 	debugSettingsDialog.SetHtml("<label for='maxTracingDepth'>Max Depth</label>" +
	// 		"<input type='number' id='maxTracingDepth' placeholder='1'> <br>" +
	// 		"<label for='enableTracing'>Enable Notification Tracing</label>" +
	// 		"<input type='checkbox' id='enableTracing' value='1'> <br>")
	// }
	// //	js.Global.Set("newDialog", debugSettingsDialog)
	// debugSettingsDialog.Call("dialog", "open")
	// //	jquery.NewJQuery("#notificationTraceSettingsDialog").Call("dialog", "open")
}

// Delete removes the element from the UniverseOfDiscourse
func (edPtr *CrlEditor) Delete(elID string) error {
	uOfD := CrlEditorSingleton.GetUofD()
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	el := uOfD.GetElement(elID)
	if el != nil {
		edPtr.cutBuffer = make(map[string]core.Element)
		edPtr.cutBuffer[elID] = el
		uOfD.MarkUndoPoint()
		return uOfD.RemoveElement(el, hl)
	}
	return nil
}

// DisplayGraph opens a new tab and displays the selected graph
func DisplayGraph(e jquery.Event, data *js.Object) {
	// TODO Reimplement DisplayGraph
	// selectedGraphIndexString := displayGraphDialog.Find("#selectedGraph").Val()
	// selectedGraphIndex, err := strconv.Atoi(selectedGraphIndexString)
	// if err != nil {
	// 	log.Printf(err.Error())
	// }
	// displayGraphDialog.Call("dialog", "close")
	// graphs := core.GetFunctionCallGraphs()
	// log.Printf("Number of function call graphs: %d\n", len(graphs))
	// log.Printf("Graphs: %#v\n", graphs)
	// if selectedGraphIndex > 0 && selectedGraphIndex <= len(graphs) {
	// 	newTab := js.Global.Get("window").Call("open", "", "_blank")
	// 	log.Printf("Selected Graph: %#v\n", graphs[selectedGraphIndex-1])
	// 	graphString := graphs[selectedGraphIndex-1].GetGraph().String()
	// 	log.Printf("Graph String: %s\n", graphString)
	// 	graphStringEscapedQuotes := strings.Replace(graphString, "\"", "\\\"", -1)
	// 	graphStringEscapedQuotesNoNL := strings.Replace(graphStringEscapedQuotes, "\n", "", -1)
	// 	htmlString := "<head>" +
	// 		"  <meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\">\n" +
	// 		"  <meta charset=\"utf-8\">\n" +
	// 		"  <title>Function Call Graph Display</title>\n" +
	// 		"  <link rel=\"stylesheet\" href=\"/js/jquery-ui-1.12.1.custom/jquery-ui.css\">\n" +
	// 		"  <script src=\"/js/jquery-ui-1.12.1.custom/external/jquery/jquery.js\"></script>\n" +
	// 		"  <script src=\"/js/viz.js\"></script>\n" +
	// 		"  <script src=\"/js/full.render.js\"></script>\n" +
	// 		"</head>\n" +
	// 		"<body>\n" +
	// 		"  <script>\n" +
	// 		"    var graphString =\"" + graphStringEscapedQuotesNoNL + "\"\n" +
	// 		"    var viz = new Viz();\n" +
	// 		"    viz.renderSVGElement(graphString)\n" +
	// 		"    .then(function(element) {\n" +
	// 		"      document.body.appendChild(element);\n" +
	// 		"    })\n" +
	// 		"    .catch(error => {\n" +
	// 		"      // Create a new Viz instance (@see Caveats page for more info) \n" +
	// 		"      viz = new Viz();\n" +
	// 		"      // Possibly display the error\n" +
	// 		"      console.error(error);\n" +
	// 		"    });\n" +
	// 		"  </script>\n" +
	// 		"</body>\n"
	// 	js.Global.Set("htmlString", htmlString)
	// 	newTab.Get("document").Call("write", htmlString)
	// }
}

// DisplayGraphDialog opens a dialog in which a trace can be selected for display
func (edPtr *CrlEditor) DisplayGraphDialog() {
	// TODO: Reimplement  DisplayGraphDialog
	// if jquery.IsEmptyObject(displayGraphDialog) {
	// 	displayGraphDialog = jquery.NewJQuery("<div class='dialog' title='Display Function Call Graph'></div>").Call("dialog", js.M{
	// 		"resizable": false,
	// 		"height":    300,
	// 		"width":     400,
	// 		"modal":     true,
	// 		"buttons":   js.M{"DisplaySelectedTrace": DisplayGraph}})
	// }
	// displayGraphDialog.SetHtml("<p id=\"numberOfAvailableFunctionCallGraphs\">Number of available function calls: " + strconv.Itoa(len(core.GetFunctionCallGraphs())) + " </output><br>" +
	// 	"<label for=\"selectedGraph\">Graph To Display</label>\n" +
	// 	"<input type=\"number\" id=\"selectedGraph\" placeholder=\"0\" min=\"0\" max=\"CrlEditor.GetNumberOfFunctionCalls()\">")
	// js.Global.Set("displayTraceDialog", displayGraphDialog)
	// js.Global.Set("numberOfAvailableFunctionCallGraphs", displayGraphDialog.Find("#numberOfAvailableFunctionCallGraphs"))
	// //	displayTraceDialog.Find("#numberOfAvailableFunctionCallGraphs").SetText("Number of available traces: " + strconv.Itoa(len(core.GetNotificationGraphs())))
	// displayGraphDialog.Call("dialog", "open")
}

// GetAdHocTrace returns the value of the AdHocTrace variable used in troubleshooting
func (edPtr *CrlEditor) GetAdHocTrace() bool {
	return core.AdHocTrace
}

// GetClientNotificationManager returns the ClientNotificationManager used to send notifications to the client
func (edPtr *CrlEditor) GetClientNotificationManager() *ClientNotificationManager {
	return edPtr.clientNotificationManager
}

// GetCurrentSelection returns the Element that is the current selection in the editor
func (edPtr *CrlEditor) GetCurrentSelection() core.Element {
	return edPtr.currentSelection
}

// GetCurrentSelectionID returns the ConceptID of the currently selected Element
func (edPtr *CrlEditor) GetCurrentSelectionID(hl *core.HeldLocks) string {
	return edPtr.currentSelection.GetConceptID(hl)
}

// GetDiagramManager returns the DiagramManager
func (edPtr *CrlEditor) GetDiagramManager() *DiagramManager {
	return edPtr.diagramManager
}

// GetNumberOfFunctionCalls returns the number of function calls in the graph
func (edPtr *CrlEditor) GetNumberOfFunctionCalls() int {
	return len(core.GetFunctionCallGraphs())
}

// GetPropertiesManager returns the PropertiesManager
func (edPtr *CrlEditor) GetPropertiesManager() *PropertiesManager {
	return edPtr.propertiesManager
}

// GetTraceChange returns the value of the core.TraceChange variable used in troubleshooting
func (edPtr *CrlEditor) GetTraceChange() bool {
	return core.TraceChange
}

// GetIconPath returns the path to the icon to be used in representing the given Element
func GetIconPath(el core.Element, hl *core.HeldLocks) string {
	isDiagram := IsDiagram(el, hl)
	switch el.(type) {
	case core.Reference:
		return "/icons/ElementReferenceIcon.svg"
	case core.Literal:
		return "/icons/LiteralIcon.svg"
	case core.Refinement:
		return "/icons/RefinementIcon.svg"
	case core.Element:
		if isDiagram {
			return "/icons/DiagramIcon.svg"
		}
		return "/icons/ElementIcon.svg"
	}
	return ""
}

// GetNotificationsLimit returns the current value of the core NotificationsLimit used in troubleshooting
func (edPtr *CrlEditor) GetNotificationsLimit() int {
	return core.GetNotificationsLimit()
}

// GetTreeDragSelection returns the Element currently being dragged from the tree
func (edPtr *CrlEditor) GetTreeDragSelection() core.Element {
	return edPtr.treeDragSelection
}

// GetTreeDragSelectionID returns the ConceptID of the Element being dragged from the tree
func (edPtr *CrlEditor) GetTreeDragSelectionID(hl *core.HeldLocks) string {
	return edPtr.treeDragSelection.GetConceptID(hl)
}

// GetTreeManager returns the TreeManager
func (edPtr *CrlEditor) GetTreeManager() *TreeManager {
	return edPtr.treeManager
}

// GetUofD returns the UniverseOfDiscourse being edited by this editor
func (edPtr *CrlEditor) GetUofD() core.UniverseOfDiscourse {
	return edPtr.uOfD
}

// IsInitialized returns true if the editor's initialization is complete
func (edPtr *CrlEditor) IsInitialized() bool {
	return edPtr.initialized
}

// SelectElement selects the indicated Element in the tree, displays the Element in the Properties window, and selects it in the
// current diagram (if present).
func (edPtr *CrlEditor) SelectElement(be core.Element, hl *core.HeldLocks) core.Element {
	edPtr.currentSelection = be
	edPtr.propertiesManager.ElementSelected(edPtr.currentSelection, hl)
	return edPtr.currentSelection
}

// SelectElementUsingIDString selects the Element whose ConceptID matches the supplied string
func (edPtr *CrlEditor) SelectElementUsingIDString(id string, hl *core.HeldLocks) core.Element {
	edPtr.currentSelection = edPtr.uOfD.GetElement(id)
	edPtr.propertiesManager.ElementSelected(edPtr.currentSelection, hl)
	return edPtr.currentSelection
}

// SetAdHocTrace sets the value of the core.AdHocTrace variable used in troubleshooting
func (edPtr *CrlEditor) SetAdHocTrace(status bool) {
	core.AdHocTrace = status
}

// SetSelectionDefinition is a convenience method for setting the Definition of the currently selected Element
func (edPtr *CrlEditor) SetSelectionDefinition(definition string, hl *core.HeldLocks) {
	switch edPtr.currentSelection.(type) {
	case core.Element:
		edPtr.currentSelection.SetDefinition(definition, hl)
	}
}

// SetSelectionLabel is a convenience method for setting the Label of the currently selected Element
func (edPtr *CrlEditor) SetSelectionLabel(name string, hl *core.HeldLocks) {
	switch edPtr.currentSelection.(type) {
	case core.Element:
		edPtr.currentSelection.SetLabel(name, hl)
	}
}

// SetSelectionURI is a convenience method for setting the URI of the curretly selected Element
func (edPtr *CrlEditor) SetSelectionURI(uri string, hl *core.HeldLocks) {
	edPtr.currentSelection.SetURI(uri, hl)
}

// SetTraceChange sets the value of the core.TraceChange variable used in troubleshooting
func (edPtr *CrlEditor) SetTraceChange(newValue bool) {
	core.TraceChange = newValue
}

// SetTraceChangeLimit sets the value of the TraceChangeLimit used in troubleshooting
func (edPtr *CrlEditor) SetTraceChangeLimit(limit int) {
	core.SetNotificationsLimit(limit)
}

// SetTreeDragSelection identifies the Element as the one being dragged from the tree
func (edPtr *CrlEditor) SetTreeDragSelection(elID string) {
	edPtr.treeDragSelection = edPtr.GetUofD().GetElement(elID)
}

// BuildEditorConceptSpace programmatically constructs the EditorConceptSpace
func BuildEditorConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// EditorConceptSpace
	conceptSpace, _ := uOfD.NewElement(hl, EditorURI)
	conceptSpace.SetLabel("EditorConceptSpace", hl)
	conceptSpace.SetURI(EditorURI, hl)
	conceptSpace.SetIsCore(hl)

	BuildTreeViews(conceptSpace, hl)

	return conceptSpace
}
