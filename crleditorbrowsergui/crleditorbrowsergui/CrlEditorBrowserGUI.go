package crleditorbrowsergui

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
	// "github.com/pbrown12303/activeCRL/crleditorbrowserguidomain"
)

// BrowserGUISingleton is the singleton instance of the BrowserGUI
var BrowserGUISingleton *CrlEditorBrowserGUI

// CrlLogClientNotifications enables logging of client notifications when set to true
var CrlLogClientNotifications = false

// CrlLogClientRequests enables the logging of client requests when set to true
var CrlLogClientRequests = false

// CrlEditorBrowserGUI is the browser gui for the CrlEditor. It manages the subordinate managers (Property, Tree, Diagram)
type CrlEditorBrowserGUI struct {
	editor                    *crleditor.Editor
	clientNotificationManager *ClientNotificationManager
	diagramManager            *diagramManager
	initialized               bool
	serverRunning             bool
	startBrowser              bool
	treeDragSelection         core.Element
	treeManager               *treeManager
	propertyManager           *propertyManager
	workingDomain             core.Element
}

// InitializeBrowserGUISingleton initializes the BrowserGUI singleton instance. It should be called once
// when the editor web page is created
func InitializeBrowserGUISingleton(editor *crleditor.Editor, startBrowser bool) {
	browserGUI := &CrlEditorBrowserGUI{}
	browserGUI.editor = editor
	browserGUI.initialized = false
	browserGUI.startBrowser = startBrowser
	BrowserGUISingleton = browserGUI
}

// CloseDiagramView closes the gui display of the diagram
func (bgPtr *CrlEditorBrowserGUI) CloseDiagramView(diagramID string, hl *core.Transaction) error {
	_, err := SendNotification("CloseDiagramView", diagramID, nil, map[string]string{})
	if err != nil {
		return errors.Wrap(err, "BrowserGUI.CloseDiagramView failed")
	}
	return nil
}

// createPropertyManager creates an instance of the propertyManager
func (bgPtr *CrlEditorBrowserGUI) createPropertyManager() error {
	pm := &propertyManager{}
	pm.browserGUI = bgPtr
	bgPtr.propertyManager = pm
	return nil
}

// createTreeManager creates an instance of the TreeManager
func (bgPtr *CrlEditorBrowserGUI) createTreeManager(treeID string) error {
	tm := &treeManager{}
	tm.browserGUI = bgPtr
	tm.treeID = treeID
	bgPtr.treeManager = tm
	return nil
}

// DisplayCallGraph opens a new tab and displays the selected graph
func (bgPtr *CrlEditorBrowserGUI) DisplayCallGraph(indexString string, hl *core.Transaction) error {
	index, err := strconv.ParseInt(indexString, 10, 32)
	if err != nil {
		return err
	}
	if index == -1 {
		// Display them all
		for _, functionCallGraph := range core.GetFunctionCallGraphs() {
			err := bgPtr.displayCallGraph(functionCallGraph, hl)
			if err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
		}
	}

	numberOfGraphs := len(core.GetFunctionCallGraphs())
	if index < 0 || index > int64(numberOfGraphs-1) {
		return errors.New("In BrowserGUI.DisplayCallGraph, index is out of bounds")
	}

	functionCallGraph := core.GetFunctionCallGraphs()[index]
	if functionCallGraph == nil {
		return errors.New("In BrowserGUI.DisplayCallGraph, function call graph is nil for index " + indexString)
	}
	return bgPtr.displayCallGraph(functionCallGraph, hl)

}

func (bgPtr *CrlEditorBrowserGUI) displayCallGraph(functionCallGraph *core.FunctionCallGraph, hl *core.Transaction) error {
	graph := functionCallGraph.GetGraph()
	if graph == nil {
		return errors.New("In BrowserGUI.displayCallGraph, graph is nil")
	}
	graphString := graph.String()
	if strings.Contains(graphString, "error") {
		return errors.New("In BrowserGUI.displayCallGraph the graph string contained an error: " + graphString)
	}
	_, err := SendNotification("DisplayGraph", "", nil, map[string]string{"GraphString": graphString})
	return err
}

// ElementDeleted is used to inform the gui that the element has been deleted
func (bgPtr *CrlEditorBrowserGUI) ElementDeleted(elID string, hl *core.Transaction) error {
	return nil
}

// ElementSelected selects the indicated Element in the tree, displays the Element in the Properties window, and selects it in the
// current diagram (if present).
func (bgPtr *CrlEditorBrowserGUI) ElementSelected(el core.Element, hl *core.Transaction) error {
	elID := ""
	var conceptState *core.ConceptState
	var err error
	if el != nil {
		elID = el.GetConceptID(hl)
		conceptState, err = core.NewConceptState(el)
		if err != nil {
			return errors.Wrap(err, "BrowserGUI.ElementSelected failed")
		}
	}
	_, err = SendNotification("ElementSelected", elID, conceptState, nil)
	if err != nil {
		return errors.Wrap(err, "In BrowserGUI.SelectElement, SendNotification failed")
	}
	return nil
}

// DisplayDiagram tells the diagramManager to display the diagram
func (bgPtr *CrlEditorBrowserGUI) DisplayDiagram(diagram core.Element, trans *core.Transaction) error {
	return bgPtr.diagramManager.displayDiagram(diagram, trans)
}

// FileLoaded informs the BrowserGUI that a file has been loaded
func (bgPtr *CrlEditorBrowserGUI) FileLoaded(el core.Element, hl *core.Transaction) {
	bgPtr.treeManager.addNodeRecursively(el, hl)
}

// GetAdHocTrace returns the value of the AdHocTrace variable used in troubleshooting
func (bgPtr *CrlEditorBrowserGUI) GetAdHocTrace() bool {
	return core.AdHocTrace
}

// GetAvailableGraphCount returns the number of available call graphs
func (bgPtr *CrlEditorBrowserGUI) GetAvailableGraphCount() int {
	return len(core.GetFunctionCallGraphs())
}

// GetClientNotificationManager returns the ClientNotificationManager used to send notifications to the client
func (bgPtr *CrlEditorBrowserGUI) GetClientNotificationManager() *ClientNotificationManager {
	return bgPtr.clientNotificationManager
}

// getDiagramManager returns the DiagramManager
func (bgPtr *CrlEditorBrowserGUI) getDiagramManager() *diagramManager {
	return bgPtr.diagramManager
}

// GetNoSaveDomains reports gui-specific domains that should not be saved
func (bgPtr *CrlEditorBrowserGUI) GetNoSaveDomains(noSaveDomains map[string]core.Element, hl *core.Transaction) {
	if bgPtr.workingDomain != nil {
		noSaveDomains[bgPtr.workingDomain.GetConceptID(hl)] = bgPtr.workingDomain
	}
}

// GetNumberOfFunctionCalls returns the number of function calls in the graph
func (bgPtr *CrlEditorBrowserGUI) GetNumberOfFunctionCalls() int {
	return len(core.GetFunctionCallGraphs())
}

// GetOmitDiagramRelatedCalls returns the value of core.OmitDiagramRelatedCalls used in troubleshooting
func (bgPtr *CrlEditorBrowserGUI) GetOmitDiagramRelatedCalls() bool {
	return core.OmitDiagramRelatedCalls
}

// GetOmitHousekeepingCalls returns the value of core.OmitHousekeepingCalls used in troubleshooting
func (bgPtr *CrlEditorBrowserGUI) GetOmitHousekeepingCalls() bool {
	return core.OmitHousekeepingCalls
}

// GetOmitManageTreeNodesCalls returns the value of core.OmitManageTreeNodesCalls used in troubleshooting
func (bgPtr *CrlEditorBrowserGUI) GetOmitManageTreeNodesCalls() bool {
	return core.OmitManageTreeNodesCalls
}

// GetTraceChange returns the value of the core.TraceChange variable used in troubleshooting
func (bgPtr *CrlEditorBrowserGUI) GetTraceChange() bool {
	return core.TraceChange
}

// GetIconPath returns the path to the icon to be used in representing the given Element
func GetIconPath(el core.Element, hl *core.Transaction) string {
	isDiagram := crldiagramdomain.IsDiagram(el, hl)
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

// GetTreeDragSelection returns the Element currently being dragged from the tree
func (bgPtr *CrlEditorBrowserGUI) GetTreeDragSelection() core.Element {
	return bgPtr.treeDragSelection
}

// GetTreeDragSelectionID returns the ConceptID of the Element being dragged from the tree
func (bgPtr *CrlEditorBrowserGUI) GetTreeDragSelectionID(hl *core.Transaction) string {
	return bgPtr.treeDragSelection.GetConceptID(hl)
}

// getTreeManager returns the TreeManager
func (bgPtr *CrlEditorBrowserGUI) getTreeManager() *treeManager {
	return bgPtr.treeManager
}

// GetUofD returns the UniverseOfDiscourse being edited by this editor
func (bgPtr *CrlEditorBrowserGUI) GetUofD() *core.UniverseOfDiscourse {
	return bgPtr.editor.GetUofD()
}

// Initialize must be called before any editor operation.
func (bgPtr *CrlEditorBrowserGUI) Initialize(hl *core.Transaction) error {
	if bgPtr.treeManager == nil {
		bgPtr.createTreeManager("#uOfD()")
	}
	err := bgPtr.treeManager.initialize(hl)
	if err != nil {
		errors.Wrap(err, "BrowserGUI.Initialize failed")
	}
	if bgPtr.diagramManager == nil {
		bgPtr.diagramManager = newDiagramManager(bgPtr)
	}
	err = bgPtr.diagramManager.initialize()
	if err != nil {
		errors.Wrap(err, "BrowserGUI.Initialize failed")
	}
	if bgPtr.clientNotificationManager == nil {
		bgPtr.clientNotificationManager = newClientNotificationManager()
	}
	if bgPtr.propertyManager == nil {
		bgPtr.createPropertyManager()
	}
	err = bgPtr.propertyManager.initialize(hl)
	if err != nil {
		errors.Wrap(err, "BrowserGUI.Initialize failed")
	}
	return nil
}

// InitializeGUI sets the client state after a browser refresh.
func (bgPtr *CrlEditorBrowserGUI) InitializeGUI(hl *core.Transaction) error {
	if !bgPtr.serverRunning {
		go bgPtr.StartServer()
		for !bgPtr.IsInitialized() {
			time.Sleep(100 * time.Millisecond)
		}
		bgPtr.serverRunning = true
	}
	err := bgPtr.initializeClientState(hl)
	if err != nil {
		return errors.Wrap(err, "Error in BrowserGUI.InitializeGUI")
	}
	return nil
}

// initializeClientState sets the client state at any desired time
func (bgPtr *CrlEditorBrowserGUI) initializeClientState(hl *core.Transaction) error {
	err := bgPtr.getTreeManager().initializeTree(hl)
	if err != nil {
		return errors.Wrap(err, "BrowserGUI.initializeClientState failed")
	}
	bgPtr.SendUserPreferences(hl)
	bgPtr.SendDebugSettings()
	bgPtr.SendWorkspacePath()
	bgPtr.SendClearDiagrams()
	_, err = SendNotification("ElementSelected", "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "BrowserGUI.initializeClientState failed")
	}
	for _, openDiagramID := range bgPtr.editor.GetSettings().OpenDiagrams {
		diagram := bgPtr.editor.GetUofD().GetElement(openDiagramID)
		if diagram == nil {
			log.Printf("In BrowserGui.initializeClientState: uOfD does not contain diagram with ID: %s", openDiagramID)
		} else {
			err = bgPtr.diagramManager.displayDiagram(diagram, hl)
			if err != nil {
				return errors.Wrap(err, "In BrowserGUI.initializeClientState diagram "+diagram.GetLabel(hl)+" did not display")
			}
		}
	}
	bgPtr.SendClientInitializationComplete()
	return nil
}

// IsInitialized returns true if the editor's initialization is complete
func (bgPtr *CrlEditorBrowserGUI) IsInitialized() bool {
	return bgPtr.initialized
}

func (bgPtr *CrlEditorBrowserGUI) nullifyReferencedConcept(refID string, hl *core.Transaction) error {
	ref := bgPtr.editor.GetUofD().GetReference(refID)
	if ref == nil {
		return errors.New("BrowserGUI.nullifyReferencedConcept called with refID not found in bgPtr.editor.GetUofD()")
	}
	err := ref.SetReferencedConceptID("", core.NoAttribute, hl)
	if err != nil {
		return errors.Wrap(err, "BrowserGUI.nullifyReferencedConcept failed")
	}
	return nil
}

// SendClearDiagrams tells the client to close all displayed diagrams
func (bgPtr *CrlEditorBrowserGUI) SendClearDiagrams() {
	bgPtr.SendNotification("ClearDiagrams", "", nil, nil)
}

// SendClientInitializationComplete tells the client that all initialization activities have been performed
func (bgPtr *CrlEditorBrowserGUI) SendClientInitializationComplete() {
	bgPtr.SendNotification("InitializationComplete", "", nil, nil)
}

// SendDebugSettings sends the trace settings to the client so that they can be edited
func (bgPtr *CrlEditorBrowserGUI) SendDebugSettings() {
	params := make(map[string]string)
	params["EnableNotificationTracing"] = strconv.FormatBool(bgPtr.GetTraceChange())
	params["OmitHousekeepingCalls"] = strconv.FormatBool(bgPtr.GetOmitHousekeepingCalls())
	params["OmitManageTreeNodesCalls"] = strconv.FormatBool(bgPtr.GetOmitManageTreeNodesCalls())
	params["OmitDiagramRelatedCalls"] = strconv.FormatBool(bgPtr.GetOmitDiagramRelatedCalls())
	bgPtr.SendNotification("DebugSettings", "", nil, params)
}

// SendUserPreferences sends the editor settings to the client so that they can be edited
func (bgPtr *CrlEditorBrowserGUI) SendUserPreferences(hl *core.Transaction) {
	params := make(map[string]string)
	params["DropReferenceAsLink"] = strconv.FormatBool(bgPtr.editor.GetDropDiagramReferenceAsLink(hl))
	params["DropRefinementAsLink"] = strconv.FormatBool(bgPtr.editor.GetDropDiagramRefinementAsLink(hl))
	bgPtr.SendNotification("UserPreferences", "", nil, params)
}

// SendNotification calls the ClientNotificationManager method of the same name and returns the result.
func (bgPtr *CrlEditorBrowserGUI) SendNotification(notificationDescription string, id string, elState *core.ConceptState, additionalParameters map[string]string) (*NotificationResponse, error) {
	return bgPtr.GetClientNotificationManager().SendNotification(notificationDescription, id, elState, additionalParameters)
}

// SendWorkspacePath sends the new workspace path to the client
func (bgPtr *CrlEditorBrowserGUI) SendWorkspacePath() {
	params := make(map[string]string)
	params["WorkspacePath"] = bgPtr.editor.GetUserPreferences().WorkspacePath
	bgPtr.SendNotification("WorkspacePath", "", nil, params)
}

// SetAdHocTrace sets the value of the core.AdHocTrace variable used in troubleshooting
func (bgPtr *CrlEditorBrowserGUI) SetAdHocTrace(status bool) {
	core.AdHocTrace = status
}

// SetInitialized tells the BrowserGUI that sockets have been initialized
func (bgPtr *CrlEditorBrowserGUI) SetInitialized() {
	bgPtr.initialized = true
}

// SetTraceChange sets the value of the core.TraceChange variable used in troubleshooting
func (bgPtr *CrlEditorBrowserGUI) SetTraceChange(newValue bool, omitHousekeepingCalls bool, omitManageTreeNodesCalls bool, omitDiagramRelatedCalls bool) {
	core.TraceChange = newValue
	core.OmitHousekeepingCalls = omitHousekeepingCalls
	core.OmitManageTreeNodesCalls = omitManageTreeNodesCalls
	core.OmitDiagramRelatedCalls = omitDiagramRelatedCalls
	core.ClearFunctionCallGraphs()
}

// SetTreeDragSelection identifies the Element as the one being dragged from the tree
func (bgPtr *CrlEditorBrowserGUI) SetTreeDragSelection(elID string) {
	bgPtr.treeDragSelection = bgPtr.GetUofD().GetElement(elID)
}

// ShowConceptInTree shows the concept in the tree
func (bgPtr *CrlEditorBrowserGUI) ShowConceptInTree(concept core.Element, hl *core.Transaction) error {
	if concept == nil {
		return errors.New("BrowserGUI.ShowConceptInTree called with nil concept")
	}
	conceptState, err := core.NewConceptState(concept)
	if err != nil {
		return errors.Wrap(err, "BrowserGUI.ShowConceptInTree failed")
	}
	_, err2 := bgPtr.SendNotification("ShowTreeNode", concept.GetConceptID(hl), conceptState, nil)
	if err2 != nil {
		return errors.Wrap(err, "BrowserGUI.ShowConceptInTree failed")
	}
	return nil
}

// UpdateDebugSettings updates the debug-related settings and sends a notification to the client
func (bgPtr *CrlEditorBrowserGUI) UpdateDebugSettings(request *Request) {
	traceChange, err := strconv.ParseBool(request.AdditionalParameters["EnableNotificationTracing"])
	if err != nil {
		log.Print(err.Error())
		return
	}
	omitHousekeeingCalls, err := strconv.ParseBool(request.AdditionalParameters["OmitHousekeepingCalls"])
	if err != nil {
		log.Print(err.Error())
		return
	}
	omitManageTreeNodesCalls, err := strconv.ParseBool(request.AdditionalParameters["OmitManageTreeNodesCalls"])
	if err != nil {
		log.Print(err.Error())
		return
	}
	omitDiagramRelatedCalls, err := strconv.ParseBool(request.AdditionalParameters["OmitDiagramRelatedCalls"])
	if err != nil {
		log.Print(err.Error())
		return
	}
	bgPtr.SetTraceChange(traceChange, omitHousekeeingCalls, omitManageTreeNodesCalls, omitDiagramRelatedCalls)
	bgPtr.SendDebugSettings()
}

// UpdateUserPreferences updates the values of the editor settings and sends a notification of the change to the client.
func (bgPtr *CrlEditorBrowserGUI) UpdateUserPreferences(request *Request, hl *core.Transaction) {
	value, _ := strconv.ParseBool(request.AdditionalParameters["DropReferenceAsLink"])
	bgPtr.editor.SetDropDiagramReferenceAsLink(value, hl)
	value, _ = strconv.ParseBool(request.AdditionalParameters["DropRefinementAsLink"])
	bgPtr.editor.SetDropDiagramRefinementAsLink(value, hl)
	bgPtr.editor.SaveUserPreferences()
	bgPtr.SendUserPreferences(hl)
}
