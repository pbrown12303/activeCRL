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

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructures"
	"github.com/pbrown12303/activeCRL/crldatatypes"
	"github.com/pbrown12303/activeCRL/crldiagram"
	"github.com/pbrown12303/activeCRL/crleditor/crleditordomain"

	"github.com/sqweek/dialog"
)

// CrlEditorSingleton is the singleton instance of the CrlEditor
var CrlEditorSingleton *CrlEditor

// CrlLogClientNotifications enables logging of client notifications when set to true
var CrlLogClientNotifications = false

// CrlLogClientRequests enables the logging of client requests when set to true
var CrlLogClientRequests = false

// CrlEditorUserPreferences are the user preferences for the editor
type CrlEditorUserPreferences struct {
	WorkspacePath               string
	DropDiagramReferenceAsLink  bool
	DropDiagramRefinementAsLink bool
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
	clientNotificationManager     *ClientNotificationManager
	userFolder                    string
	userPreferences               *CrlEditorUserPreferences
	settings                      core.Element
	currentSelection              core.Element
	cutBuffer                     map[string]core.Element
	diagramManager                *diagramManager
	initialized                   bool
	treeDragSelection             core.Element
	treeManager                   *treeManager
	uOfD                          *core.UniverseOfDiscourse
	workspaceFiles                map[string]*workspaceFile
	workingConceptSpace           core.Element
	defaultConceptSpaceLabelCount int
	defaultElementLabelCount      int
	defaultLiteralLabelCount      int
	defaultReferenceLabelCount    int
	defaultRefinementLabelCount   int
}

// InitializeCrlEditorSingleton initializes the CrlEditor singleton instance. It should be called once
// when the editor web page is created
func InitializeCrlEditorSingleton(userFolderArg string) {
	var editor CrlEditor
	editor.initialized = false
	CrlEditorSingleton = &editor
	CrlEditorSingleton.initializeEditor(userFolderArg)
	editor.initialized = true
	log.Printf("Editor initialized")
}

// initializeEditor must be called before any other editor operation. It should not be used to reinitialize the editor
func (edPtr *CrlEditor) initializeEditor(userFolderArg string) {
	edPtr.uOfD = core.NewUniverseOfDiscourse()
	hl := edPtr.uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	edPtr.initializeUofD(hl)

	if userFolderArg == "" {
		edPtr.userFolder = home
	} else {
		edPtr.userFolder = userFolderArg
	}
	var userPreferences CrlEditorUserPreferences
	edPtr.userPreferences = &userPreferences
	edPtr.cutBuffer = make(map[string]core.Element)
	edPtr.workspaceFiles = make(map[string]*workspaceFile)

	domain := edPtr.uOfD.GetElementWithURI(crleditordomain.EditorDomainURI)
	BuildTreeViewManager(domain, hl)
	BuildDiagramViewMonitor(domain, hl)

	edPtr.createTreeManager(edPtr, "#uOfD", hl)
	edPtr.diagramManager = newDiagramManager(edPtr)
	edPtr.clientNotificationManager = newClientNotificationManager()

}

// reinitializeEditor is used to re-initialize the editor
func (edPtr *CrlEditor) reinitializeEditor() {
	edPtr.uOfD = core.NewUniverseOfDiscourse()
	hl := edPtr.uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	edPtr.initializeUofD(hl)

	edPtr.cutBuffer = make(map[string]core.Element)
	edPtr.workspaceFiles = make(map[string]*workspaceFile)

	domain := edPtr.uOfD.GetElementWithURI(crleditordomain.EditorDomainURI)
	BuildTreeViewManager(domain, hl)
	BuildDiagramViewMonitor(domain, hl)

	edPtr.createSettings(hl)
	edPtr.createTreeManager(edPtr, "#uOfD", hl)
	edPtr.diagramManager = newDiagramManager(edPtr)
	edPtr.ResetDefaultLabelCounts()
	// edPtr.clientNotificationManager = newClientNotificationManager()
	edPtr.userPreferences.WorkspacePath = ""
	edPtr.uOfD.SetRecordingUndo(true)
}

func (edPtr *CrlEditor) initializeUofD(hl *core.HeldLocks) error {
	crldatatypes.BuildCrlDataTypesConceptSpace(edPtr.uOfD, hl)
	crldatastructures.BuildCrlDataStructuresConceptSpace(edPtr.uOfD, hl)
	hl.ReleaseLocksAndWait()
	crldiagram.BuildCrlDiagramConceptSpace(edPtr.uOfD, hl)
	hl.ReleaseLocksAndWait()
	_, err := AddEditorConceptSpaceToUofD(edPtr.uOfD, hl)
	if err != nil {
		return err
	}
	hl.ReleaseLocksAndWait()
	// Create editor working concept space
	edPtr.workingConceptSpace, _ = edPtr.uOfD.NewElement(hl)
	edPtr.workingConceptSpace.SetLabel("EditorWorkingCS", hl)
	registerDiagramViewMonitorFunctions(edPtr.uOfD)
	return nil
}

// AddEditorConceptSpaceToUofD adds the concepts representing the various editor views to the universe of discourse
func AddEditorConceptSpaceToUofD(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	conceptSpace := uOfD.GetElementWithURI(crleditordomain.EditorURI)
	if conceptSpace == nil {
		var err error
		conceptSpace, err = crleditordomain.BuildEditorDomain(uOfD, hl)
		if err != nil {
			return nil, err
		}
	}
	return conceptSpace, nil
}

// createTreeManager creates an instance of the TreeManager
func (edPtr *CrlEditor) createTreeManager(editor *CrlEditor, treeID string, hl *core.HeldLocks) error {
	var tm treeManager
	tm.editor = editor
	tm.treeID = treeID
	edPtr.treeManager = &tm
	var err error
	tm.treeNodeManager, err = edPtr.uOfD.CreateReplicateAsRefinementFromURI(TreeNodeManagerURI, hl)
	if err != nil {
		return err
	}
	treeNodeManagerUOfDReference := tm.treeNodeManager.GetFirstOwnedReferenceRefinedFromURI(TreeNodeManagerUofDReferenceURI, hl)
	treeNodeManagerUOfDReference.SetReferencedConcept(edPtr.uOfD, hl)
	tm.treeNodeManager.SetOwningConcept(edPtr.workingConceptSpace, hl)
	tm.treeNodeManager.SetIsCoreRecursively(hl)
	hl.ReleaseLocksAndWait()
	registerTreeViewFunctions(edPtr.uOfD)
	return nil
}

// createSettings creates the concept space for settings and adds it to the workspace
func (edPtr *CrlEditor) createSettings(hl *core.HeldLocks) error {

	newSettings, err := edPtr.uOfD.CreateReplicateAsRefinementFromURI(crleditordomain.EditorSettingsURI, hl)
	if err != nil {
		return err
	}
	edPtr.settings = newSettings
	openDiagrams := edPtr.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
	diagram := edPtr.uOfD.GetElementWithURI(crldiagram.CrlDiagramURI)
	crldatastructures.SetListType(openDiagrams, diagram, hl)
	return nil
}

// ClearWorkspace closes the current workspace without saving the root elements
func (edPtr *CrlEditor) ClearWorkspace(hl *core.HeldLocks) error {
	var err error
	rootElements := edPtr.uOfD.GetRootElements(hl)
	for id, wf := range edPtr.workspaceFiles {
		if rootElements[id] == nil {
			edPtr.deleteFile(wf, hl)
			delete(edPtr.workspaceFiles, id)
		}
	}
	edPtr.reinitializeEditor()
	edPtr.initializeClientState(hl)
	return err
}

// CloseWorkspace closes the current workspace, saving the root elements
func (edPtr *CrlEditor) CloseWorkspace(hl *core.HeldLocks) error {
	var err error
	if edPtr.userPreferences.WorkspacePath != "" {
		err = edPtr.SaveWorkspace(hl)
		if err != nil {
			return err
		}
		for _, wsf := range edPtr.workspaceFiles {
			err = wsf.File.Close()
			if err != nil {
				return err
			}
		}
	}
	edPtr.reinitializeEditor()
	edPtr.initializeClientState(hl)
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
	qualifiedFilename := edPtr.userPreferences.WorkspacePath + "/" + wf.Info.Name()
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

func (edPtr *CrlEditor) getDefaultConceptSpaceLabel() string {
	edPtr.defaultConceptSpaceLabelCount++
	countString := strconv.Itoa(edPtr.defaultConceptSpaceLabelCount)
	return "ConceptSpace" + countString
}

func (edPtr *CrlEditor) getDefaultElementLabel() string {
	edPtr.defaultElementLabelCount++
	countString := strconv.Itoa(edPtr.defaultElementLabelCount)
	return "Element" + countString
}

func (edPtr *CrlEditor) getDefaultLiteralLabel() string {
	edPtr.defaultLiteralLabelCount++
	countString := strconv.Itoa(edPtr.defaultLiteralLabelCount)
	return "Literal" + countString
}

func (edPtr *CrlEditor) getDefaultReferenceLabel() string {
	edPtr.defaultReferenceLabelCount++
	countString := strconv.Itoa(edPtr.defaultReferenceLabelCount)
	return "Reference" + countString
}

func (edPtr *CrlEditor) getDefaultRefinementLabel() string {
	edPtr.defaultRefinementLabelCount++
	countString := strconv.Itoa(edPtr.defaultRefinementLabelCount)
	return "Refinement" + countString
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

// GetDropDiagramReferenceAsLink returns true if dropped references are shown as links
func (edPtr *CrlEditor) GetDropDiagramReferenceAsLink(hl *core.HeldLocks) bool {
	return edPtr.userPreferences.DropDiagramReferenceAsLink
}

// GetDropDiagramRefinementAsLink returns true if dropped refinements are shown as links
func (edPtr *CrlEditor) GetDropDiagramRefinementAsLink(hl *core.HeldLocks) bool {
	return edPtr.userPreferences.DropDiagramRefinementAsLink
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
func (edPtr *CrlEditor) GetUofD() *core.UniverseOfDiscourse {
	return edPtr.uOfD
}

// GetUserPreferences returns the current user's preferences
func (edPtr *CrlEditor) GetUserPreferences() *CrlEditorUserPreferences {
	return edPtr.userPreferences
}

// getUserPreferencesPath returns the path to the user preferences
func (edPtr *CrlEditor) getUserPreferencesPath() string {
	return edPtr.userFolder + "/.crleditoruserpreferences"
}

// GetWorkspacePath return the path to the current workspace
func (edPtr *CrlEditor) GetWorkspacePath() string {
	return edPtr.userPreferences.WorkspacePath
}

// InitializeClient sets the client state after a browser refresh.
func (edPtr *CrlEditor) InitializeClient() error {
	<-webSocketReady
	uOfD := edPtr.GetUofD()
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	if edPtr.IsInitialized() == false {
		time.Sleep(100 * time.Millisecond)
	}
	return edPtr.initializeClientState(hl)
}

// initializeClientState sets the client state at any desired time
func (edPtr *CrlEditor) initializeClientState(hl *core.HeldLocks) error {
	err := edPtr.getTreeManager().initializeTree(hl)
	if err != nil {
		return err
	}
	edPtr.SendUserPreferences(hl)
	edPtr.SendDebugSettings()
	edPtr.SendWorkspacePath()
	edPtr.SendClearDiagrams()
	openDiagrams := edPtr.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
	if openDiagrams == nil {
		return errors.New("In CrlEditor.initializeClientState, openDiagrams is nil")
	}
	openDiagramReference, err2 := crldatastructures.GetFirstMemberReference(openDiagrams, hl)
	if err2 != nil {
		return err2
	}
	for openDiagramReference != nil {
		diagram := openDiagramReference.GetReferencedConcept(hl)
		if diagram == nil {
			return errors.New("Failed to load diagram with ID: " + openDiagramReference.GetReferencedConceptID(hl))
		}
		err2 = edPtr.diagramManager.displayDiagram(diagram, hl)
		if err2 != nil {
			return err2
		}
		openDiagramReference, _ = crldatastructures.GetNextMemberReference(openDiagramReference, hl)
	}
	hl.ReleaseLocksAndWait()
	edPtr.SendClientInitializationComplete()
	return nil
}

// IsInitialized returns true if the editor's initialization is complete
func (edPtr *CrlEditor) IsInitialized() bool {
	return edPtr.initialized
}

// LoadUserPreferences loads the user preferences saved in the user's home directory
func (edPtr *CrlEditor) LoadUserPreferences(workspaceArg string) error {
	path := edPtr.getUserPreferencesPath()
	_, err := os.Stat(path)
	if err != nil {
		// it is OK to not find the file
		edPtr.userPreferences.WorkspacePath = workspaceArg
		return nil
	}
	fileSettings, err2 := ioutil.ReadFile(path)
	if err2 != nil {
		return err
	}
	err = json.Unmarshal(fileSettings, edPtr.userPreferences)
	if err != nil {
		return err
	}
	return nil
}

// newFile creates a file with the name being the ConceptID of the supplied Element and returns the workspaceFile struct
// HACK: For troubleshooting, prepend the ConceptLabel to the filename so that the individual files may be better identified without
// opening them. This, however, will result in duplicate files should the label of the Concept be changed.
// TODO: remove this hack
func (edPtr *CrlEditor) newFile(el core.Element, hl *core.HeldLocks) (*workspaceFile, error) {
	if edPtr.userPreferences.WorkspacePath == "" {
		return nil, errors.New("CrlEditor.NewFile called with no settings.WorkspacePath defined")
	}
	// HACK: here's the hack
	filename := edPtr.userPreferences.WorkspacePath + "/" + el.GetLabel(hl) + "--" + el.GetConceptID(hl) + ".acrl"
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
	file, err := os.OpenFile(edPtr.userPreferences.WorkspacePath+"/"+fileInfo.Name(), mode, fileInfo.Mode())
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

// LoadWorkspace loads the workspace currently designated by the userPreferences.WorkspacePath. If the path is empty, it is a no-op.
func (edPtr *CrlEditor) LoadWorkspace(hl *core.HeldLocks) error {
	files, err := ioutil.ReadDir(edPtr.userPreferences.WorkspacePath)
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
			if workspaceFile.ConceptSpace.IsRefinementOfURI(crleditordomain.EditorSettingsURI, hl) {
				if edPtr.settings != nil {
					edPtr.uOfD.DeleteElement(edPtr.settings, hl)
					hl.ReleaseLocksAndWait()
				}
				edPtr.settings = workspaceFile.ConceptSpace
			}
		}
	}
	if edPtr.settings == nil {
		err = edPtr.createSettings(hl)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadWorkspaceOnly loads the workspace currently designated by the userPreferences.WorkspacePath. If the path is empty, it is a no-op.

// openWorkspace sets the path to the folder to be used as a workspace. It is the implementation of a request from the client.
func (edPtr *CrlEditor) openWorkspace(hl *core.HeldLocks) error {
	if edPtr.userPreferences.WorkspacePath != "" {
		return errors.New("Cannot open another workspace in the same editor - close existing workspace first")
	}
	path, err2 := edPtr.SelectWorkspace()
	if err2 != nil {
		return err2
	}
	return edPtr.openWorkspaceImpl(path, hl)
}

func (edPtr *CrlEditor) openWorkspaceImpl(path string, hl *core.HeldLocks) error {
	edPtr.userPreferences.WorkspacePath = path
	err := edPtr.LoadWorkspace(hl)
	if err != nil {
		return err
	}
	return edPtr.initializeClientState(hl)
}

// OpenWorkspaceProgrammatically is intended for use in automated testing scenarios
func (edPtr *CrlEditor) OpenWorkspaceProgrammatically(path string, hl *core.HeldLocks) error {
	defer hl.ReleaseLocksAndWait()
	if path == "" {
		return errors.New("OpenWorkspaceProgrammatically called with empty path")
	}
	return edPtr.openWorkspaceImpl(path, hl)
}

// Redo performs an undo on the uOfD and refreshes the interface
func (edPtr *CrlEditor) Redo(hl *core.HeldLocks) error {
	edPtr.uOfD.Redo(hl)
	edPtr.initializeClientState(hl)
	return nil
}

// ResetDefaultLabelCounts re-initializes the default counters for all new model elements
func (edPtr *CrlEditor) ResetDefaultLabelCounts() {
	edPtr.defaultConceptSpaceLabelCount = 0
	edPtr.defaultElementLabelCount = 0
	edPtr.defaultLiteralLabelCount = 0
	edPtr.defaultReferenceLabelCount = 0
	edPtr.defaultRefinementLabelCount = 0
	edPtr.diagramManager.ResetDefaultLabelCounts()
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

// saveUserPreferences saves the current user preferences to the user's home directory
func (edPtr *CrlEditor) saveUserPreferences() error {
	f, err := os.OpenFile(edPtr.getUserPreferencesPath(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	serializedUserPreferences, err2 := json.Marshal(edPtr.userPreferences)
	if err2 != nil {
		return err2
	}
	_, err = f.Write(serializedUserPreferences)
	if err != nil {
		return err
	}
	err = f.Truncate(int64(len(serializedUserPreferences)))
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
			if workspaceFile != nil {
				err = edPtr.saveFile(workspaceFile, hl)
				if err != nil {
					return err
				}
			} else {
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

// SelectWorkspace opens a dialog for the user to select a workspace
func (edPtr *CrlEditor) SelectWorkspace() (string, error) {
	return dialog.Directory().Title("Select a directory for your workspace").Browse()
}

// SendClearDiagrams tells the client to close all displayed diagrams
func (edPtr *CrlEditor) SendClearDiagrams() {
	edPtr.SendNotification("ClearDiagrams", "", nil, nil)
}

// SendClientInitializationComplete tells the client that all initialization activities have been performed
func (edPtr *CrlEditor) SendClientInitializationComplete() {
	edPtr.SendNotification("InitializationComplete", "", nil, nil)
}

// SendDebugSettings sends the trace settings to the client so that they can be edited
func (edPtr *CrlEditor) SendDebugSettings() {
	params := make(map[string]string)
	params["EnableNotificationTracing"] = strconv.FormatBool(edPtr.GetTraceChange())
	params["OmitHousekeepingCalls"] = strconv.FormatBool(edPtr.GetOmitHousekeepingCalls())
	params["OmitManageTreeNodesCalls"] = strconv.FormatBool(edPtr.GetOmitManageTreeNodesCalls())
	edPtr.SendNotification("DebugSettings", "", nil, params)
}

// SendUserPreferences sends the editor settings to the client so that they can be edited
func (edPtr *CrlEditor) SendUserPreferences(hl *core.HeldLocks) {
	params := make(map[string]string)
	params["DropReferenceAsLink"] = strconv.FormatBool(edPtr.GetDropDiagramReferenceAsLink(hl))
	params["DropRefinementAsLink"] = strconv.FormatBool(edPtr.GetDropDiagramRefinementAsLink(hl))
	edPtr.SendNotification("UserPreferences", "", nil, params)
}

// SendNotification calls the ClientNotificationManager method of the same name and returns the result.
func (edPtr *CrlEditor) SendNotification(notificationDescription string, id string, el core.Element, additionalParameters map[string]string) (*NotificationResponse, error) {
	return edPtr.GetClientNotificationManager().SendNotification(notificationDescription, id, el, additionalParameters)
}

// SendWorkspacePath sends the new workspace path to the client
func (edPtr *CrlEditor) SendWorkspacePath() {
	params := make(map[string]string)
	params["WorkspacePath"] = edPtr.userPreferences.WorkspacePath
	edPtr.SendNotification("WorkspacePath", "", nil, params)
}

// SetAdHocTrace sets the value of the core.AdHocTrace variable used in troubleshooting
func (edPtr *CrlEditor) SetAdHocTrace(status bool) {
	core.AdHocTrace = status
}

// SetDropDiagramReferenceAsLink returns true if dropped references are shown as links
func (edPtr *CrlEditor) SetDropDiagramReferenceAsLink(value bool, hl *core.HeldLocks) {
	edPtr.userPreferences.DropDiagramReferenceAsLink = value
}

// SetDropDiagramRefinementAsLink returns true if dropped refinements are shown as links
func (edPtr *CrlEditor) SetDropDiagramRefinementAsLink(value bool, hl *core.HeldLocks) {
	edPtr.userPreferences.DropDiagramRefinementAsLink = value
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

// SetWorkspacePath sets the user's preference WorkspacePath value.
func (edPtr *CrlEditor) SetWorkspacePath(path string) error {
	edPtr.userPreferences.WorkspacePath = path
	return edPtr.saveUserPreferences()
}

// Undo performs an undo on the uOfD and refreshes the interface
func (edPtr *CrlEditor) Undo(hl *core.HeldLocks) error {
	edPtr.uOfD.Undo(hl)
	return edPtr.initializeClientState(hl)
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

// UpdateUserPreferences updates the values of the editor settings and sends a notification of the change to the client.
func (edPtr *CrlEditor) UpdateUserPreferences(request *Request, hl *core.HeldLocks) {
	value, _ := strconv.ParseBool(request.AdditionalParameters["DropReferenceAsLink"])
	edPtr.SetDropDiagramReferenceAsLink(value, hl)
	value, _ = strconv.ParseBool(request.AdditionalParameters["DropRefinementAsLink"])
	edPtr.SetDropDiagramRefinementAsLink(value, hl)
	edPtr.saveUserPreferences()
	edPtr.SendUserPreferences(hl)
}
