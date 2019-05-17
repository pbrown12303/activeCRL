package editor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pbrown12303/activeCRL/crldiagram"

	"github.com/pbrown12303/activeCRL/core"
)

// editorURI is the URI for accessing the CrlEditor
var editorURI = "http://activeCrl.com/crlEditor/Editor"

// CrlEditorSingleton is the singleton instance of the CrlEditor
var CrlEditorSingleton *CrlEditor

// CrlLogClientNotifications enables logging of client notifications when set to true
var CrlLogClientNotifications = false

// CrlLogClientRequests enables the logging of client requests when set to true
var CrlLogClientRequests = false

// CrlEditorSettings are the configurable behaviors of the editor
type CrlEditorSettings struct {
	DropReferenceAsLink  bool
	DropRefinementAsLink bool
	WorkspacePath        string
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
	settings                  *CrlEditorSettings
	currentSelection          core.Element
	cutBuffer                 map[string]core.Element
	diagramManager            *diagramManager
	initialized               bool
	treeDragSelection         core.Element
	treeManager               *treeManager
	uOfD                      core.UniverseOfDiscourse
	workspaceFiles            map[string]*workspaceFile
	workingConceptSpace       core.Element
}

// InitializeCrlEditorSingleton initializes the CrlEditor singleton instance. It should be called once
// when the editor web page is created
func InitializeCrlEditorSingleton() {
	var editor CrlEditor
	editor.initialized = false
	var settings CrlEditorSettings
	editor.settings = &settings
	editor.cutBuffer = make(map[string]core.Element)
	editor.workspaceFiles = make(map[string]*workspaceFile)

	editor.uOfD = core.NewUniverseOfDiscourse()
	hl := editor.uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()

	// Set the value of the singleton
	CrlEditorSingleton = &editor

	editor.treeManager = newTreeManager(&editor, "#uOfD")
	editor.diagramManager = newDiagramManager(&editor)
	editor.clientNotificationManager = newClientNotificationManager()

	editor.initializeUofD(hl)

	editor.initialized = true
	log.Printf("Editor initialized")
}

func (edPtr *CrlEditor) initializeUofD(hl *core.HeldLocks) error {
	crldiagram.BuildCrlDiagramConceptSpace(edPtr.uOfD, hl)
	hl.ReleaseLocksAndWait()
	AddEditorConceptSpaceToUofD(edPtr.uOfD, hl)
	hl.ReleaseLocksAndWait()
	// Create editor working concept space
	edPtr.workingConceptSpace, _ = edPtr.uOfD.NewElement(hl)
	edPtr.workingConceptSpace.SetLabel("EditorWorkingCS", hl)
	treeNodeManager, err := edPtr.treeManager.configureUofD(hl)
	if err != nil {
		return err
	}
	treeNodeManager.SetOwningConcept(edPtr.workingConceptSpace, hl)
	treeNodeManager.SetIsCoreRecursively(hl)
	hl.ReleaseLocksAndWait()
	registerTreeViewFunctions(edPtr.uOfD)
	registerDiagramViewMonitorFunctions(edPtr.uOfD)
	return nil
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

// AddEditorConceptSpaceToUofD adds the concepts representing the various editor views to the universe of discurse
func AddEditorConceptSpaceToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
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
	edPtr.settings.WorkspacePath = ""
	edPtr.workspaceFiles = make(map[string]*workspaceFile)
	hl.ReleaseLocksAndWait()
	edPtr.uOfD = core.NewUniverseOfDiscourse()
	hl2 := edPtr.uOfD.NewHeldLocks()
	defer hl2.ReleaseLocksAndWait()
	err = edPtr.initializeUofD(hl2)
	if err != nil {
		return err
	}
	_, err = SendNotification("Refresh", "", nil, nil)
	return err
}

// Delete removes the element from the UniverseOfDiscourse
func (edPtr *CrlEditor) Delete(elID string) error {
	uOfD := CrlEditorSingleton.GetUofD()
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	el := uOfD.GetElement(elID)
	if el != nil {
		// TODO: Populate cut buffer with full set of deleted elements
		// edPtr.cutBuffer = make(map[string]core.Element)
		// edPtr.cutBuffer[elID] = el
		uOfD.MarkUndoPoint()
		err := uOfD.DeleteElement(el, hl)
		if err != nil {
			return err
		}
		edPtr.SelectElement(nil, hl)
	}
	return nil
}

// deleteFile deletes the file from the os
func (edPtr *CrlEditor) deleteFile(wf *workspaceFile, hl *core.HeldLocks) error {
	qualifiedFilename := edPtr.settings.WorkspacePath + "/" + wf.Info.Name()
	return os.Remove(qualifiedFilename)
}

// DisplayCallGraph opens a new tab and displays the selected graph
func (edPtr *CrlEditor) DisplayCallGraph(indexString string, hl *core.HeldLocks) error {
	index, err := strconv.ParseInt(indexString, 10, 32)
	if err != nil {
		return err
	}
	if index == -1 {
		// Display them all
		for _, functionCallGraph := range core.GetFunctionCallGraphs() {
			err := edPtr.displayCallGraph(functionCallGraph, hl)
			if err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
		}
	}

	numberOfGraphs := len(core.GetFunctionCallGraphs())
	if index < 0 || index > int64(numberOfGraphs-1) {
		return errors.New("In CrlEditor.DisplayCallGraph, index is out of bounds")
	}

	functionCallGraph := core.GetFunctionCallGraphs()[index]
	if functionCallGraph == nil {
		return errors.New("In CrlEditor.DisplayCallGraph, function call graph is nil for index " + indexString)
	}
	return edPtr.displayCallGraph(functionCallGraph, hl)

}

func (edPtr *CrlEditor) displayCallGraph(functionCallGraph *core.FunctionCallGraph, hl *core.HeldLocks) error {
	graph := functionCallGraph.GetGraph()
	if graph == nil {
		return errors.New("In CrlEditor.displayCallGraph, graph is nil")
	}
	graphString := graph.String()
	if strings.Contains(graphString, "error") {
		return errors.New("In CrlEditor.displayCallGraph the graph string contained an error: " + graphString)
	}
	_, err := SendNotification("DisplayGraph", "", nil, map[string]string{"GraphString": graphString})
	return err
}

// GetAdHocTrace returns the value of the AdHocTrace variable used in troubleshooting
func (edPtr *CrlEditor) GetAdHocTrace() bool {
	return core.AdHocTrace
}

// GetClientNotificationManager returns the ClientNotificationManager used to send notifications to the client
func (edPtr *CrlEditor) GetClientNotificationManager() *ClientNotificationManager {
	return edPtr.clientNotificationManager
}

// GetAvailableGraphCount returns the number of available call graphs
func (edPtr *CrlEditor) GetAvailableGraphCount() int {
	return len(core.GetFunctionCallGraphs())
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

// GetOmitHousekeepingCalls returns the value of core.OmitHousekeepingCalls used in troubleshooting
func (edPtr *CrlEditor) GetOmitHousekeepingCalls() bool {
	return core.OmitHousekeepingCalls
}

// GetOmitManageTreeNodesCalls returns the value of core.OmitManageTreeNodesCalls used in troubleshooting
func (edPtr *CrlEditor) GetOmitManageTreeNodesCalls() bool {
	return core.OmitManageTreeNodesCalls
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
		return "/icons/ReferenceIcon.svg"
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
	return edPtr.settings
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

// loadSettings loads the settings saved in the user's home directory
func (edPtr *CrlEditor) loadSettings() error {
	path := home + "/.crleditorsettings"
	_, err := os.Stat(path)
	if err != nil {
		// it is OK to not find the file
		return nil
	}
	fileSettings, err2 := ioutil.ReadFile(path)
	if err2 != nil {
		return err
	}
	err = json.Unmarshal(fileSettings, edPtr.settings)
	if err != nil {
		return err
	}
	return nil
}

// newFile creates a file with the name being the ConceptID of the supplied Element and returns the workspaceFile struct
func (edPtr *CrlEditor) newFile(el core.Element, hl *core.HeldLocks) (*workspaceFile, error) {
	if edPtr.settings.WorkspacePath == "" {
		return nil, errors.New("CrlEditor.NewFile called with no settings.WorkspacePath defined")
	}
	filename := edPtr.settings.WorkspacePath + "/" + el.GetConceptID(hl) + ".acrl"
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
	file, err := os.OpenFile(edPtr.settings.WorkspacePath+"/"+fileInfo.Name(), mode, fileInfo.Mode())
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
	if path != edPtr.settings.WorkspacePath && edPtr.settings.WorkspacePath != "" {
		return errors.New("Cannot open another workspace in the same editor - a new editor must be started")
	}
	edPtr.settings.WorkspacePath = path
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

// saveSettings saves the current settings to the user's home directory
func (edPtr *CrlEditor) saveSettings() error {
	path := home + "/.crleditorsettings"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	serializedSettings, err2 := json.Marshal(edPtr.settings)
	if err2 != nil {
		return err2
	}
	_, err = f.Write(serializedSettings)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

// SaveWorkspace saves all top-level concepts whose versions are different than the last retrieved version.
func (edPtr *CrlEditor) SaveWorkspace(hl *core.HeldLocks) error {
	rootElements := edPtr.uOfD.GetRootElements(hl)
	var err error
	for id, el := range rootElements {
		if el.GetIsCore(hl) == false && edPtr.workingConceptSpace != el {
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
func (edPtr *CrlEditor) SelectElement(el core.Element, hl *core.HeldLocks) core.Element {
	edPtr.currentSelection = el
	selectedID := ""
	if el != nil {
		selectedID = el.GetConceptID(hl)
	}
	_, err := SendNotification("ElementSelected", selectedID, el, nil)
	if err != nil {
		log.Printf(err.Error())
	}
	return edPtr.currentSelection
}

// SelectElementUsingIDString selects the Element whose ConceptID matches the supplied string
func (edPtr *CrlEditor) SelectElementUsingIDString(id string, hl *core.HeldLocks) core.Element {
	foundElement := edPtr.uOfD.GetElement(id)
	return edPtr.SelectElement(foundElement, hl)
}

// SendDebugSettings sends the trace settings to the client so that they can be edited
func (edPtr *CrlEditor) SendDebugSettings() {
	params := make(map[string]string)
	params["EnableNotificationTracing"] = strconv.FormatBool(edPtr.GetTraceChange())
	params["OmitHousekeepingCalls"] = strconv.FormatBool(edPtr.GetOmitHousekeepingCalls())
	params["OmitManageTreeNodesCalls"] = strconv.FormatBool(edPtr.GetOmitManageTreeNodesCalls())
	edPtr.SendNotification("DebugSettings", "", nil, params)
}

// SendEditorSettings sends the editor settings to the client so that they can be edited
func (edPtr *CrlEditor) SendEditorSettings() {
	params := make(map[string]string)
	params["DropReferenceAsLink"] = strconv.FormatBool(edPtr.settings.DropReferenceAsLink)
	params["DropRefinementAsLink"] = strconv.FormatBool(edPtr.settings.DropRefinementAsLink)
	edPtr.SendNotification("EditorSettings", "", nil, params)
}

// SendNotification calls the ClientNotificationManager method of the same name and returns the result.
func (edPtr *CrlEditor) SendNotification(notificationDescription string, id string, el core.Element, additionalParameters map[string]string) (*NotificationResponse, error) {
	return edPtr.GetClientNotificationManager().SendNotification(notificationDescription, id, el, additionalParameters)
}

// SendWorkspacePath sends the new workspace path to the client
func (edPtr *CrlEditor) SendWorkspacePath() {
	params := make(map[string]string)
	params["WorkspacePath"] = edPtr.settings.WorkspacePath
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
func (edPtr *CrlEditor) SetTraceChange(newValue bool, omitHousekeepingCalls bool, omitManageTreeNodesCalls bool) {
	core.TraceChange = newValue
	core.OmitHousekeepingCalls = omitHousekeepingCalls
	core.OmitManageTreeNodesCalls = omitManageTreeNodesCalls
	core.ClearFunctionCallGraphs()
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
	omitHousekeeingCalls, err := strconv.ParseBool(request.AdditionalParameters["OmitHousekeepingCalls"])
	if err != nil {
		log.Printf(err.Error())
		return
	}
	omitManageTreeNodesCalls, err := strconv.ParseBool(request.AdditionalParameters["OmitManageTreeNodesCalls"])
	if err != nil {
		log.Printf(err.Error())
		return
	}
	edPtr.SetTraceChange(traceChange, omitHousekeeingCalls, omitManageTreeNodesCalls)
	edPtr.SendDebugSettings()
}

// UpdateEditorSettings updates the values of the editor settings and sends a notification of the change to the client.
func (edPtr *CrlEditor) UpdateEditorSettings(request *Request) {
	edPtr.settings.DropReferenceAsLink, _ = strconv.ParseBool(request.AdditionalParameters["DropReferenceAsLink"])
	edPtr.settings.DropRefinementAsLink, _ = strconv.ParseBool(request.AdditionalParameters["DropRefinementAsLink"])
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
	BuildDiagramViewMonitor(conceptSpace, hl)

	return conceptSpace
}
