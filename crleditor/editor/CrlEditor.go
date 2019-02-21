package editor

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pbrown12303/activeCRL/crldiagram"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/core"
)

// editorURI is the URI for accessing the CrlEditor
var editorURI = "http://activeCrl.com/crlEditor/Editor"

// CrlEditorSingleton is the singleton instance of the CrlEditor
var CrlEditorSingleton *CrlEditor

// CrlEditorSettings are the configurable behaviors of the editor
type CrlEditorSettings struct {
	DropReferenceAsLink  bool
	DropRefinementAsLink bool
}

type workspaceFile struct {
	File          *os.File
	LoadedVersion int
	Info          os.FileInfo
	ConceptSpace  core.Element
}

// CrlEditor is the central component of the CrlEditor. It manages the subordinate managers (Property, Tree, Diagram)
// and the singleton instances of the Universe of Discourse and HeldLocks shared by all editing operations.
type CrlEditor struct {
	clientNotificationManager *ClientNotificationManager
	crlEditorSettings         *CrlEditorSettings
	currentSelection          core.Element
	cutBuffer                 map[string]core.Element
	diagramManager            *diagramManager
	initialized               bool
	propertiesManager         *PropertiesManager
	treeDragSelection         core.Element
	treeManager               *treeManager
	uOfD                      core.UniverseOfDiscourse
	workspaceFiles            map[string]*workspaceFile
	workspacePath             string
}

// InitializeCrlEditorSingleton initializes the CrlEditor singleton instance. It should be called once
// when the editor web page is created
func InitializeCrlEditorSingleton() {
	var editor CrlEditor
	editor.initialized = false
	var settings CrlEditorSettings
	editor.crlEditorSettings = &settings
	editor.cutBuffer = make(map[string]core.Element)
	editor.workspaceFiles = make(map[string]*workspaceFile)

	editor.uOfD = core.NewUniverseOfDiscourse()
	hl := editor.uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()

	// Set the value of the singleton
	CrlEditorSingleton = &editor

	editor.treeManager = newTreeManager(&editor, "#uOfD")
	editor.diagramManager = newDiagramManager(&editor)
	editor.propertiesManager = NewPropertiesManager(&editor)
	editor.clientNotificationManager = newClientNotificationManager()

	editor.initializeUofD(hl)

	editor.initialized = true
	log.Printf("Editor initialized")
}

func (edPtr *CrlEditor) initializeUofD(hl *core.HeldLocks) {
	crldiagram.BuildCrlDiagramConceptSpace(edPtr.uOfD, hl)
	hl.ReleaseLocksAndWait()
	AddEditorViewsToUofD(edPtr.uOfD, hl)
	hl.ReleaseLocksAndWait()
	edPtr.treeManager.configureUofD(hl)
	hl.ReleaseLocksAndWait()
	registerTreeViewFunctions(edPtr.uOfD)
}

// InitializeClient sets the client state after a browser refresh.
func InitializeClient() {
	<-webSocketReady
	log.Printf("InitializeClient executing")
	uOfD := CrlEditorSingleton.GetUofD()
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	if CrlEditorSingleton.IsInitialized() == false {
		time.Sleep(100 * time.Millisecond)
	}
	CrlEditorSingleton.getTreeManager().initializeTree(hl)
	CrlEditorSingleton.SendEditorSettings()
	CrlEditorSingleton.SendDebugSettings()
}

// AddEditorViewsToUofD adds the concepts representing the various editor views to the universe of discurse
func AddEditorViewsToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	conceptSpace := uOfD.GetElementWithURI(editorURI)
	if conceptSpace == nil {
		conceptSpace = BuildEditorConceptSpace(uOfD, hl)
		if conceptSpace == nil {
			log.Printf("Build of Editor Concept Space failed")
		}
	}
	return conceptSpace
}

// CloseWorkspace closes the current workspace, saving the root elements
func (edPtr *CrlEditor) CloseWorkspace(hl *core.HeldLocks) error {
	err := edPtr.SaveWorkspace(hl)
	if err != nil {
		return err
	}
	hl.ReleaseLocksAndWait()
	edPtr.uOfD = core.NewUniverseOfDiscourse()
	hl2 := edPtr.uOfD.NewHeldLocks()
	defer hl2.ReleaseLocksAndWait()
	edPtr.initializeUofD(hl2)
	edPtr.treeManager.configureUofD(hl2)
	return nil
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

// deleteFile deletes the file from the os
func (edPtr *CrlEditor) deleteFile(wf *workspaceFile, hl *core.HeldLocks) error {
	qualifiedFilename := edPtr.workspacePath + "/" + wf.Info.Name()
	return os.Remove(qualifiedFilename)
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

// getDiagramManager returns the DiagramManager
func (edPtr *CrlEditor) getDiagramManager() *diagramManager {
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
	isDiagram := crldiagram.IsDiagram(el, hl)
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

// GetEditorSettings returns the settings that impact editor behavior
func (edPtr *CrlEditor) GetEditorSettings() *CrlEditorSettings {
	return edPtr.crlEditorSettings
}

// GetTreeDragSelection returns the Element currently being dragged from the tree
func (edPtr *CrlEditor) GetTreeDragSelection() core.Element {
	return edPtr.treeDragSelection
}

// GetTreeDragSelectionID returns the ConceptID of the Element being dragged from the tree
func (edPtr *CrlEditor) GetTreeDragSelectionID(hl *core.HeldLocks) string {
	return edPtr.treeDragSelection.GetConceptID(hl)
}

// getTreeManager returns the TreeManager
func (edPtr *CrlEditor) getTreeManager() *treeManager {
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

// newFile creates a file with the name being the ConceptID of the supplied Element and returns the workspaceFile struct
func (edPtr *CrlEditor) newFile(el core.Element, hl *core.HeldLocks) (*workspaceFile, error) {
	if edPtr.workspacePath == "" {
		return nil, errors.New("CrlEditor.NewFile called with no workspacePath defined")
	}
	filename := edPtr.workspacePath + "/" + el.GetConceptID(hl) + ".acrl"
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	fileInfo, err2 := os.Stat(filename)
	if err2 != nil {
		return nil, err2
	}
	var wf workspaceFile
	wf.ConceptSpace = el
	wf.File = file
	wf.LoadedVersion = el.GetVersion(hl)
	wf.Info = fileInfo
	return &wf, nil
}

// openFile opens the file and returns a workspaceFile struct
func (edPtr *CrlEditor) openFile(fileInfo os.FileInfo, hl *core.HeldLocks) (*workspaceFile, error) {
	writable := (fileInfo.Mode().Perm() & 0200) > 0
	mode := os.O_RDONLY
	if writable {
		mode = os.O_RDWR
	}
	file, err := os.OpenFile(edPtr.workspacePath+"/"+fileInfo.Name(), mode, fileInfo.Mode())
	if err != nil {
		return nil, err
	}
	fileContent := make([]byte, fileInfo.Size())
	_, err = file.Read(fileContent)
	if err != nil {
		return nil, err
	}
	element, err2 := edPtr.uOfD.RecoverConceptSpace(fileContent, hl)
	if err2 != nil {
		return nil, err2
	}
	if !writable {
		element.SetReadOnlyRecursively(true, hl)
	}
	edPtr.treeManager.addNodeRecursively(element, hl)
	var wf workspaceFile
	wf.ConceptSpace = element
	wf.Info = fileInfo
	wf.LoadedVersion = element.GetVersion(hl)
	wf.File = file
	return &wf, nil
}

// openWorkspace sets the path to the folder to be used as a workspace
func (edPtr *CrlEditor) openWorkspace(path string, hl *core.HeldLocks) error {
	if path != edPtr.workspacePath && edPtr.workspacePath != "" {
		return errors.New("Cannot open another workspace in the same editor. A new editor must be started.")
	}
	edPtr.workspacePath = path
	edPtr.SendWorkspacePath()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".acrl") {
			workspaceFile, err := edPtr.openFile(f, hl)
			if err != nil {
				return err
			}
			edPtr.workspaceFiles[workspaceFile.ConceptSpace.GetConceptID(hl)] = workspaceFile
		}
	}
	return nil
}

// saveFile saves the file and updates the fileInfo
func (edPtr *CrlEditor) saveFile(wf *workspaceFile, hl *core.HeldLocks) error {
	hl.ReadLockElement(wf.ConceptSpace)
	if wf.File == nil {
		return errors.New("CrlEditor.SaveFile called with nil file")
	}
	byteArray, err := edPtr.uOfD.MarshalConceptSpace(wf.ConceptSpace, hl)
	if err != nil {
		return err
	}
	var length int
	length, err = wf.File.WriteAt(byteArray, 0)
	if err != nil {
		return err
	}
	err = wf.File.Truncate(int64(length))
	if err != nil {
		return err
	}
	return wf.File.Sync()
}

// SaveWorkspace saves all top-level concepts whose versions are different than the last retrieved version.
func (edPtr *CrlEditor) SaveWorkspace(hl *core.HeldLocks) error {
	rootElements := edPtr.uOfD.GetRootElements(hl)
	var err error
	for id, el := range rootElements {
		if el.GetIsCore(hl) == false {
			workspaceFile := edPtr.workspaceFiles[id]
			if workspaceFile != nil && workspaceFile.LoadedVersion < el.GetVersion(hl) {
				err = edPtr.saveFile(workspaceFile, hl)
				if err != nil {
					return err
				}
				break
			}
			if workspaceFile == nil {
				workspaceFile, err = edPtr.newFile(el, hl)
				if err != nil {
					return err
				}
				edPtr.workspaceFiles[id] = workspaceFile
				err = edPtr.saveFile(workspaceFile, hl)
				if err != nil {
					return err
				}
			}
		}
	}
	for id, wf := range edPtr.workspaceFiles {
		if rootElements[id] == nil {
			edPtr.deleteFile(wf, hl)
			delete(edPtr.workspaceFiles, id)
		}
	}
	return nil
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

// SendDebugSettings sends the trace settings to the client so that they can be edited
func (edPtr *CrlEditor) SendDebugSettings() {
	params := make(map[string]string)
	params["EnableNotificationTracing"] = strconv.FormatBool(edPtr.GetTraceChange())
	edPtr.SendNotification("DebugSettings", "", nil, params)
}

// SendEditorSettings sends the editor settings to the client so that they can be edited
func (edPtr *CrlEditor) SendEditorSettings() {
	params := make(map[string]string)
	params["DropReferenceAsLink"] = strconv.FormatBool(edPtr.crlEditorSettings.DropReferenceAsLink)
	params["DropRefinementAsLink"] = strconv.FormatBool(edPtr.crlEditorSettings.DropRefinementAsLink)
	edPtr.SendNotification("EditorSettings", "", nil, params)
}

// SendNotification calls the ClientNotificationManager method of the same name and returns the result.
func (edPtr *CrlEditor) SendNotification(notificationDescription string, id string, el core.Element, additionalParameters map[string]string) (*NotificationResponse, error) {
	return edPtr.GetClientNotificationManager().SendNotification(notificationDescription, id, el, additionalParameters)
}

// SendWorkspacePath sends the new workspace path to the client
func (edPtr *CrlEditor) SendWorkspacePath() {
	params := make(map[string]string)
	params["WorkspacePath"] = edPtr.workspacePath
	edPtr.SendNotification("WorkspacePath", "", nil, params)
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

// UpdateDebugSettings updates the debug-related settings and sends a notification to the client
func (edPtr *CrlEditor) UpdateDebugSettings(request *Request) {
	traceChange, err := strconv.ParseBool(request.AdditionalParameters["EnableNotificationTracing"])
	if err != nil {
		log.Printf(err.Error())
		return
	}
	edPtr.SetTraceChange(traceChange)
	edPtr.SendDebugSettings()
}

// UpdateEditorSettings updates the values of the editor settings and sends a notification of the change to the client.
func (edPtr *CrlEditor) UpdateEditorSettings(request *Request) {
	edPtr.crlEditorSettings.DropReferenceAsLink, _ = strconv.ParseBool(request.AdditionalParameters["DropReferenceAsLink"])
	edPtr.crlEditorSettings.DropRefinementAsLink, _ = strconv.ParseBool(request.AdditionalParameters["DropRefinementAsLink"])
	edPtr.SendEditorSettings()
}

// BuildEditorConceptSpace programmatically constructs the EditorConceptSpace
func BuildEditorConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// EditorConceptSpace
	conceptSpace, _ := uOfD.NewElement(hl, editorURI)
	conceptSpace.SetLabel("EditorConceptSpace", hl)
	conceptSpace.SetURI(editorURI, hl)
	conceptSpace.SetIsCore(hl)

	BuildTreeViews(conceptSpace, hl)

	return conceptSpace
}
