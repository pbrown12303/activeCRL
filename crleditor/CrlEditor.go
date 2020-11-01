package crleditor

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructuresdomain"
	"github.com/pbrown12303/activeCRL/crldatatypesdomain"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditordomain"
	"github.com/pbrown12303/activeCRL/crlmapsdomain"

	"github.com/sqweek/dialog"
)

// UserPreferences are the user preferences for the editor
type UserPreferences struct {
	WorkspacePath               string
	DropDiagramReferenceAsLink  bool
	DropDiagramRefinementAsLink bool
}

// Editor manages one or more CrlEditors
type Editor struct {
	currentSelection            core.Element
	cutBuffer                   map[string]core.Element
	defaultDomainLabelCount     int
	defaultElementLabelCount    int
	defaultLiteralLabelCount    int
	defaultReferenceLabelCount  int
	defaultRefinementLabelCount int
	defaultDiagramLabelCount    int
	editorGUIs                  []EditorGUI
	exitRequested               bool
	home                        string
	settings                    core.Element
	uOfDManager                 *core.UofDManager
	userPreferences             *UserPreferences
	userFolder                  string
	workspaceManager            *CrlWorkspaceManager
}

// NewEditor returns an initialized Editor
func NewEditor(userFolderArg string) *Editor {
	editor := &Editor{}
	var err error
	editor.home, err = os.UserHomeDir()
	if err != nil {
		log.Fatalf("User home directory not found")
	}
	editor.userPreferences = &UserPreferences{}
	if userFolderArg == "" {
		editor.userFolder = editor.home
	} else {
		editor.userFolder = userFolderArg
	}
	editor.uOfDManager = &core.UofDManager{}
	editor.workspaceManager = NewCrlWorkspaceManager(editor)
	return editor
}

// AddDiagramToDisplayedList adds the diagramID to the list of displayed diagrams
func (editor *Editor) AddDiagramToDisplayedList(diagramID string, hl *core.HeldLocks) error {
	if !editor.IsDiagramDisplayed(diagramID, hl) {
		openDiagrams := editor.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
		_, err := crldatastructuresdomain.AppendStringListMember(openDiagrams, diagramID, hl)
		if err != nil {
			return errors.Wrap(err, "diagramManager.addDiagramToDisplayedList failed")
		}
	}
	return nil
}

// AddEditorGUI adds an editor to the list of editorGUIs being managed by the
func (editor *Editor) AddEditorGUI(editorGUI EditorGUI) error {
	editor.editorGUIs = append(editor.editorGUIs, editorGUI)
	// err := editorGUI.RegisterUofDInitializationFunctions(editor.uOfDManager)
	// if err != nil {
	// 	return errors.Wrap(err, "Editor.AddEditorGUI failed")
	// }
	// err := editorGUI.RegisterUofDPostInitializationFunctions(editor.uOfDManager)
	// if err != nil {
	// 	return errors.Wrap(err, "Editor.AddEditorGUI failed")
	// }
	return nil
}

// ClearWorkspace clears all files in the current workspace that correspond to uOfD root elements
// and then reinitializes all editorGUIs.
func (editor *Editor) ClearWorkspace(hl *core.HeldLocks) error {
	workspacePath := editor.userPreferences.WorkspacePath
	err := editor.workspaceManager.ClearWorkspace(workspacePath, hl)
	if err != nil {
		return errors.Wrap(err, "crleditor.Editor.ClearWorkspace failed")
	}

	err = editor.Initialize(workspacePath, false)
	if err != nil {
		return errors.Wrap(err, "crleditor.Editor.ClearWorkspace failed")
	}
	return nil
}

// CloseDiagramView removes the diagram from the list of displayed diagrams and informs all GUIs
func (editor *Editor) CloseDiagramView(diagramID string, hl *core.HeldLocks) error {
	// If the diagram is in the list of displayed diagrams, remove it
	if editor.IsDiagramDisplayed(diagramID, hl) {
		editor.RemoveDiagramFromDisplayedList(diagramID, hl)
	}
	for _, gui := range editor.editorGUIs {
		err := gui.CloseDiagramView(diagramID, hl)
		if err != nil {
			return errors.Wrap(err, "Editor.CloseDiagramView failed")
		}
	}
	return nil
}

// CloseWorkspace closes the current workspace, saving the root elements
func (editor *Editor) CloseWorkspace(hl *core.HeldLocks) error {
	var err error
	if editor.userPreferences.WorkspacePath != "" {
		err = editor.workspaceManager.CloseWorkspace(hl)
		if err != nil {
			return errors.Wrap(err, "CrlEditor.CloseWorkspace failed")
		}
	}
	// The hl here is from the old UofD. Initialize will create a new one, so we first release the locks on the old one
	hl.ReleaseLocksAndWait()
	editor.SetWorkspacePath("")
	err = editor.Initialize("", false)
	if err != nil {
		return errors.Wrap(err, "crleditor.Editor.CloseWorkspace failed")
	}
	return nil
}

// createSettings creates the concept space for settings and adds it to the workspace
func (editor *Editor) createSettings(hl *core.HeldLocks) error {
	newSettings, err := editor.GetUofD().CreateReplicateAsRefinementFromURI(crleditordomain.EditorSettingsURI, hl)
	if err != nil {
		return errors.Wrap(err, "Editor.createSettings failed")
	}
	editor.settings = newSettings
	openDiagrams := editor.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
	diagram := editor.GetUofD().GetElementWithURI(crldiagramdomain.CrlDiagramURI)
	crldatastructuresdomain.SetListType(openDiagrams, diagram, hl)
	return nil
}

// DeleteElement removes the element from the UniverseOfDiscourse
func (editor *Editor) DeleteElement(elID string, hl *core.HeldLocks) error {
	defer hl.ReleaseLocksAndWait()
	el := editor.GetUofD().GetElement(elID)
	if el != nil {
		// TODO: Populate cut buffer with full set of deleted elements
		// editor.cutBuffer = make(map[string]core.Element)
		// editor.cutBuffer[elID] = el
		err := editor.GetUofD().DeleteElement(el, hl)
		if err != nil {
			return errors.Wrap(err, "Editor.DeleteElement failed")
		}
		editor.SelectElement(nil, hl)
	}
	for _, gui := range editor.editorGUIs {
		err := gui.ElementDeleted(elID, hl)
		if err != nil {
			errors.Wrap(err, "Editor.DeleteElement failed")
		}
	}
	return nil
}

// FileLoaded is used to inform the CrlEditor that a file has been loaded
func (editor *Editor) FileLoaded(el core.Element, hl *core.HeldLocks) {
	for _, editorGUI := range editor.editorGUIs {
		editorGUI.FileLoaded(el, hl)
	}
}

// GetCurrentSelection returns the Element that is the current selection in the editor
func (editor *Editor) GetCurrentSelection() core.Element {
	return editor.currentSelection
}

// GetCurrentSelectionID returns the ConceptID of the currently selected Element
func (editor *Editor) GetCurrentSelectionID(hl *core.HeldLocks) string {
	if editor.currentSelection == nil {
		return ""
	}
	return editor.currentSelection.GetConceptID(hl)
}

// GetDefaultDomainLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultDomainLabel() string {
	editor.defaultDomainLabelCount++
	countString := strconv.Itoa(editor.defaultDomainLabelCount)
	return "Domain" + countString
}

// GetDefaultDiagramLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultDiagramLabel() string {
	editor.defaultDiagramLabelCount++
	countString := strconv.Itoa(editor.defaultDiagramLabelCount)
	return "Diagram" + countString
}

// GetDefaultElementLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultElementLabel() string {
	editor.defaultElementLabelCount++
	countString := strconv.Itoa(editor.defaultElementLabelCount)
	return "Element" + countString
}

// GetDefaultLiteralLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultLiteralLabel() string {
	editor.defaultLiteralLabelCount++
	countString := strconv.Itoa(editor.defaultLiteralLabelCount)
	return "Literal" + countString
}

// GetDefaultReferenceLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultReferenceLabel() string {
	editor.defaultReferenceLabelCount++
	countString := strconv.Itoa(editor.defaultReferenceLabelCount)
	return "Reference" + countString
}

// GetDefaultRefinementLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultRefinementLabel() string {
	editor.defaultRefinementLabelCount++
	countString := strconv.Itoa(editor.defaultRefinementLabelCount)
	return "Refinement" + countString
}

// GetDropDiagramReferenceAsLink returns true if dropped references are shown as links
func (editor *Editor) GetDropDiagramReferenceAsLink(hl *core.HeldLocks) bool {
	return editor.userPreferences.DropDiagramReferenceAsLink
}

// GetDropDiagramRefinementAsLink returns true if dropped refinements are shown as links
func (editor *Editor) GetDropDiagramRefinementAsLink(hl *core.HeldLocks) bool {
	return editor.userPreferences.DropDiagramRefinementAsLink
}

// GetExitRequested returns true if exit has been requested
func (editor *Editor) GetExitRequested() bool {
	return editor.exitRequested
}

// GetNotificationsLimit returns the current value of the core NotificationsLimit used in troubleshooting
func (editor *Editor) GetNotificationsLimit() int {
	return core.GetNotificationsLimit()
}

// getNoSaveDomains returns a map of the editor domains that should not be saved
func (editor *Editor) getNoSaveDomains(hl *core.HeldLocks) map[string]core.Element {
	noSaveDomains := make(map[string]core.Element)
	for _, editor := range editor.editorGUIs {
		editor.GetNoSaveDomains(noSaveDomains, hl)
	}
	return noSaveDomains
}

// GetSettings returns the editor settings
func (editor *Editor) GetSettings() core.Element {
	return editor.settings
}

// GetUofD returns the current UniverseOfDiscourse
func (editor *Editor) GetUofD() *core.UniverseOfDiscourse {
	return editor.uOfDManager.UofD
}

// GetUserPreferences returns the current user's preferences
func (editor *Editor) GetUserPreferences() *UserPreferences {
	return editor.userPreferences
}

// getUserPreferencesPath returns the path to the user preferences
func (editor *Editor) getUserPreferencesPath() string {
	return editor.userFolder + "/.crleditoruserpreferences"
}

// GetWorkspacePath return the path to the current workspace
func (editor *Editor) GetWorkspacePath() string {
	return editor.userPreferences.WorkspacePath
}

// Initialize initializes the uOfD, workspace manager, and all registered editorGUIs
func (editor *Editor) Initialize(workspacePath string, promptWorkspaceSelection bool) error {
	editor.settings = nil
	editor.uOfDManager.Initialize()
	if editor.workspaceManager == nil {
		editor.workspaceManager = NewCrlWorkspaceManager(editor)
	}
	editor.workspaceManager.Initialize()
	editor.workspaceManager.LoadUserPreferences(workspacePath)
	if editor.userPreferences.WorkspacePath != workspacePath {
		editor.SetWorkspacePath(workspacePath)
	}
	if editor.userPreferences.WorkspacePath == "" && promptWorkspaceSelection {
		workspacePath, err := editor.SelectWorkspace()
		if err != nil {
			return errors.Wrap(err, "Editor.Initialize failed")
		}
		err = editor.SetWorkspacePath(workspacePath)
		if err != nil {
			return errors.Wrap(err, "Editor.Initialize failed")
		}
	}
	editor.cutBuffer = make(map[string]core.Element)
	hl := editor.uOfDManager.UofD.NewHeldLocks()
	editor.resetDefaultLabelCounts()

	crldatatypesdomain.BuildCrlDataTypesDomain(editor.GetUofD(), hl)
	crldatastructuresdomain.BuildCrlDataStructuresDomain(editor.GetUofD(), hl)
	crldiagramdomain.BuildCrlDiagramDomain(editor.GetUofD(), hl)
	crleditordomain.BuildEditorDomain(editor.GetUofD(), hl)
	err := crlmapsdomain.BuildCrlMapsDomain(editor.GetUofD(), hl)
	if err != nil {
		return errors.Wrap(err, "Editor.Initialize failed")
	}
	hl.ReleaseLocksAndWait()

	for _, editor := range editor.editorGUIs {
		err = editor.Initialize(hl)
		if err != nil {
			return errors.Wrap(err, "Editor.Initialize failed")
		}
	}
	hl.ReleaseLocksAndWait()

	if editor.userPreferences.WorkspacePath != "" {
		err = editor.workspaceManager.LoadWorkspace(hl)
	}
	if err != nil {
		return errors.Wrap(err, "Editor.Initialize failed")
	}
	if editor.settings == nil {
		err = editor.createSettings(hl)
		if err != nil {
			return errors.Wrap(err, "Editor.Initialize failed")
		}
	}
	hl.ReleaseLocksAndWait()

	for _, editor := range editor.editorGUIs {
		err = editor.InitializeGUI(hl)
		if err != nil {
			return errors.Wrap(err, "Editor.Initialize failed")
		}
	}
	hl.ReleaseLocksAndWait()

	editor.uOfDManager.UofD.SetRecordingUndo(true)
	return nil
}

// InitializeGUI tells all GUIs to initialize their state
func (editor *Editor) InitializeGUI(hl *core.HeldLocks) error {
	for _, gui := range editor.editorGUIs {
		err := gui.InitializeGUI(hl)
		if err != nil {
			return errors.Wrap(err, "Editor.InitializeGUI failed")
		}
	}
	return nil
}

// IsDiagramDisplayed returns true if the diagram is in the list of displayed diagrams
func (editor *Editor) IsDiagramDisplayed(diagramID string, hl *core.HeldLocks) bool {
	openDiagrams := editor.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
	return crldatastructuresdomain.IsStringListMember(openDiagrams, diagramID, hl)
}

// LoadWorkspace tells the editor to load the workspace
func (editor *Editor) LoadWorkspace(hl *core.HeldLocks) error {
	err := editor.workspaceManager.LoadWorkspace(hl)
	if err != nil {
		return errors.Wrap(err, "Editor.LoadWorkspace failed")
	}
	if editor.settings == nil {
		err = editor.createSettings(hl)
		if err != nil {
			return errors.Wrap(err, "Editor.LoadWorkspace failed")
		}
	}
	return nil
}

// OpenWorkspace sets the path to the folder to be used as a workspace. It is the implementation of a request from the client.
func (editor *Editor) OpenWorkspace() error {
	if editor.userPreferences.WorkspacePath != "" {
		return errors.New("Cannot open another workspace in the same editor - close existing workspace first")
	}
	path, err2 := editor.SelectWorkspace()
	if err2 != nil {
		return err2
	}
	return editor.Initialize(path, false)
}

// func (editor *Editor) openWorkspaceImpl(path string, hl *core.HeldLocks) error {
// 	err := editor.Initialize(path, false)
// 	if err != nil {
// 		return errors.Wrap(err, "Editor.openWorkspaceImpl failed")
// 	}
// 	return nil
// }

// // OpenWorkspaceProgrammatically is intended for use in automated testing scenarios
// func (editor *Editor) OpenWorkspaceProgrammatically(path string, hl *core.HeldLocks) error {
// 	defer hl.ReleaseLocksAndWait()
// 	if path == "" {
// 		return errors.New("OpenWorkspaceProgrammatically called with empty path")
// 	}
// 	return editor.openWorkspaceImpl(path, hl)
// }

// Redo performs an undo on the editor.editor.GetUofD() and refreshes the interface
func (editor *Editor) Redo(hl *core.HeldLocks) error {
	editor.GetUofD().Redo(hl)
	err := editor.InitializeGUI(hl)
	if err != nil {
		return errors.Wrap(err, "Editor.Redo failed")
	}
	return nil
}

// RemoveDiagramFromDisplayedList removes the diagramID from the list of displayed diagrams
func (editor *Editor) RemoveDiagramFromDisplayedList(diagramID string, hl *core.HeldLocks) {
	if editor.IsDiagramDisplayed(diagramID, hl) {
		openDiagrams := editor.settings.GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
		crldatastructuresdomain.RemoveStringListMember(openDiagrams, diagramID, hl)
	}
}

// ResetDefaultLabelCounts re-initializes the default counters for all new model elements
func (editor *Editor) resetDefaultLabelCounts() {
	editor.defaultDomainLabelCount = 0
	editor.defaultElementLabelCount = 0
	editor.defaultLiteralLabelCount = 0
	editor.defaultReferenceLabelCount = 0
	editor.defaultRefinementLabelCount = 0
	editor.defaultDiagramLabelCount = 0
}

// SaveUserPreferences saves the current user preferences to the user's home directory
func (editor *Editor) SaveUserPreferences() error {
	f, err := os.OpenFile(editor.getUserPreferencesPath(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	serializedUserPreferences, err2 := json.Marshal(editor.userPreferences)
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

// SaveWorkspace saves the workspace
func (editor *Editor) SaveWorkspace(hl *core.HeldLocks) error {
	err := editor.workspaceManager.SaveWorkspace(hl)
	if err != nil {
		return errors.Wrap(err, "Editor.SaveWorkspace failed")
	}
	return nil
}

// SelectElement selects the indicated Element in the tree, displays the Element in the Properties window, and selects it in the
// current diagram (if present).
func (editor *Editor) SelectElement(el core.Element, hl *core.HeldLocks) error {
	editor.currentSelection = el
	for _, gui := range editor.editorGUIs {
		err := gui.ElementSelected(el, hl)
		if err != nil {
			return errors.Wrap(err, "Editor.SelectElement failed")
		}
	}
	return nil
}

// SelectElementUsingIDString selects the Element whose ConceptID matches the supplied string
func (editor *Editor) SelectElementUsingIDString(id string, hl *core.HeldLocks) error {
	foundElement := editor.GetUofD().GetElement(id)
	if foundElement == nil && id != "" {
		return errors.New("In BrowserGUI.SelectElementUsingIDString, element was not found")
	}
	return editor.SelectElement(foundElement, hl)
}

// SelectWorkspace opens a dialog for the user to select a workspace
func (editor *Editor) SelectWorkspace() (string, error) {
	return dialog.Directory().Title("Select a directory for your workspace").Browse()
}

func (editor *Editor) setSettings(settings core.Element, hl *core.HeldLocks) error {
	if settings == nil {
		return errors.New("Editor.setSettings called with nil settings")
	}
	if settings.IsRefinementOfURI(crleditordomain.EditorSettingsURI, hl) == false {
		return errors.New("Editor.setSettings called with nil settings")
	}
	if editor.settings != nil {
		err := editor.GetUofD().DeleteElement(editor.settings, hl)
		if err != nil {
			return errors.Wrap(err, "Editor.setSettings failed")
		}
		hl.ReleaseLocksAndWait()
	}
	editor.settings = settings
	return nil
}

// SetDropDiagramReferenceAsLink returns true if dropped references are shown as links
func (editor *Editor) SetDropDiagramReferenceAsLink(value bool, hl *core.HeldLocks) {
	editor.userPreferences.DropDiagramReferenceAsLink = value
}

// SetDropDiagramRefinementAsLink returns true if dropped refinements are shown as links
func (editor *Editor) SetDropDiagramRefinementAsLink(value bool, hl *core.HeldLocks) {
	editor.userPreferences.DropDiagramRefinementAsLink = value
}

// SetExitRequested informs the Editor that exit has been requested. Intended to be used by the GUI
func (editor *Editor) SetExitRequested() {
	editor.exitRequested = true
}

// SetSelectionDefinition is a convenience method for setting the Definition of the currently selected Element
func (editor *Editor) SetSelectionDefinition(definition string, hl *core.HeldLocks) {
	editor.currentSelection.SetDefinition(definition, hl)
}

// SetSelectionLabel is a convenience method for setting the Label of the currently selected Element
func (editor *Editor) SetSelectionLabel(name string, hl *core.HeldLocks) {
	editor.currentSelection.SetLabel(name, hl)
}

// SetSelectionURI is a convenience method for setting the URI of the curretly selected Element
func (editor *Editor) SetSelectionURI(uri string, hl *core.HeldLocks) {
	editor.currentSelection.SetURI(uri, hl)
}

// SetWorkspacePath sets the user's preference WorkspacePath value.
func (editor *Editor) SetWorkspacePath(path string) error {
	editor.userPreferences.WorkspacePath = path
	return editor.SaveUserPreferences()
}

// Undo performs an undo on the editor.GetUofD() and refreshes the interface
func (editor *Editor) Undo(hl *core.HeldLocks) error {
	editor.GetUofD().Undo(hl)
	for _, gui := range editor.editorGUIs {
		err := gui.InitializeGUI(hl)
		if err != nil {
			return errors.Wrap(err, "Editor.Undo failed")
		}
	}
	return nil
}

// EditorGUI is the interface for all CrlEditors, independent of implementation technology
type EditorGUI interface {
	CloseDiagramView(diagramID string, hl *core.HeldLocks) error
	ElementDeleted(elID string, hl *core.HeldLocks) error
	ElementSelected(el core.Element, hl *core.HeldLocks) error
	FileLoaded(el core.Element, hl *core.HeldLocks)
	GetNoSaveDomains(noSaveDomains map[string]core.Element, hl *core.HeldLocks)
	Initialize(hl *core.HeldLocks) error
	InitializeGUI(hl *core.HeldLocks) error
}
