package crleditor

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructuresdomain"
	"github.com/pbrown12303/activeCRL/crldatatypesdomain"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crlmapsdomain"

	"github.com/sqweek/dialog"
)

// UserPreferences are the user preferences for the editor
type UserPreferences struct {
	WorkspacePath               string
	DropDiagramReferenceAsLink  bool
	DropDiagramRefinementAsLink bool
	HorizontalLayoutSpacing     float64
	VerticalLayoutSpacing       float64
}

// Settings reflect the current status of the editing session
type Settings struct {
	DefaultDomainLabelCount     int
	DefaultElementLabelCount    int
	DefaultLiteralLabelCount    int
	DefaultReferenceLabelCount  int
	DefaultRefinementLabelCount int
	DefaultDiagramLabelCount    int
	Selection                   string
	OpenDiagrams                []string
	CurrentDiagram              string
}

// CrlEditorSingleton is the unique single instance of the Editor in an editng session
var CrlEditorSingleton *Editor

// Editor manages one or more CrlEditors
type Editor struct {
	currentSelection           *core.Concept
	cutBuffer                  map[string]*core.Concept
	editorGUIs                 []EditorGUI
	exitRequested              bool
	home                       string
	settings                   *Settings
	uOfDManager                *core.UofDManager
	diagramManager             *DiagramManager
	userPreferences            *UserPreferences
	userFolder                 string
	workspaceManager           *CrlWorkspaceManager
	inProgressTransaction      *core.Transaction
	transientSelection         *core.Concept
	transientDisplayedDiagrams *core.Concept
	transientCurrentDiagram    *core.Concept
	undoRedoInProgress         bool
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
	editor.settings = &Settings{}
	editor.settings.OpenDiagrams = []string{}
	if userFolderArg == "" {
		editor.userFolder = editor.home
	} else {
		editor.userFolder = userFolderArg
	}
	editor.uOfDManager = &core.UofDManager{}
	editor.workspaceManager = NewCrlWorkspaceManager(editor)
	editor.diagramManager = NewDiagramManager(editor)
	return editor
}

// addDiagramToDisplayedList adds the diagramID to the list of displayed diagrams
func (editor *Editor) addDiagramToDisplayedList(diagramID string, trans *core.Transaction) error {
	if !editor.IsDiagramDisplayed(diagramID, trans) {
		editor.UpdateOpenDiagrams(append(editor.settings.OpenDiagrams, diagramID), trans)
	}
	return nil
}

// AddEditorGUI adds an editor to the list of editorGUIs being managed by the editor
// and initializes the GUI
func (editor *Editor) AddEditorGUI(editorGUI EditorGUI) error {
	editor.editorGUIs = append(editor.editorGUIs, editorGUI)
	trans, isNew := editor.GetTransaction()
	if isNew {
		defer editor.EndTransaction()
	}
	editorGUI.Initialize(trans)
	editorGUI.RefreshGUI(trans)
	return nil
}

// ClearWorkspace clears all files in the current workspace that correspond to uOfD root elements
// and then reinitializes all editorGUIs.
func (editor *Editor) ClearWorkspace(trans *core.Transaction) error {
	workspacePath := editor.userPreferences.WorkspacePath
	err := editor.workspaceManager.ClearWorkspace(workspacePath, trans)
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
func (editor *Editor) CloseDiagramView(diagramID string, trans *core.Transaction) error {
	// If the diagram is in the list of displayed diagrams, remove it
	if editor.IsDiagramDisplayed(diagramID, trans) {
		editor.RemoveDiagramFromDisplayedList(diagramID, trans)
	}
	for _, gui := range editor.editorGUIs {
		err := gui.CloseDiagramView(diagramID, trans)
		if err != nil {
			return errors.Wrap(err, "Editor.CloseDiagramView failed")
		}
	}
	return nil
}

// CloseWorkspace closes the current workspace, saving the root elements
func (editor *Editor) CloseWorkspace(trans *core.Transaction) error {
	var err error
	if editor.userPreferences.WorkspacePath != "" {
		err = editor.workspaceManager.CloseWorkspace(trans)
		if err != nil {
			return errors.Wrap(err, "CrlEditor.CloseWorkspace failed")
		}
	}
	// The trans here is from the old UofD. Initialize will create a new one, so we first release the locks on the old one
	editor.EndTransaction()
	editor.SetWorkspacePath("")
	err = editor.Initialize("", false)
	if err != nil {
		return errors.Wrap(err, "crleditor.Editor.CloseWorkspace failed")
	}
	return nil
}

// DeleteElement removes the element from the UniverseOfDiscourse
func (editor *Editor) DeleteElement(elID string, trans *core.Transaction) error {
	el := editor.GetUofD().GetElement(elID)
	if el != nil {
		// TODO: Populate cut buffer with full set of deleted elements
		// editor.cutBuffer = make(map[string]core.Element)
		// editor.cutBuffer[elID] = el
		if el.IsRefinementOfURI(crldiagramdomain.CrlDiagramURI, trans) {
			if editor.IsDiagramDisplayed(elID, trans) {
				editor.CloseDiagramView(elID, trans)
			}
		}
		err := editor.GetUofD().DeleteElement(el, trans)
		if err != nil {
			return errors.Wrap(err, "Editor.DeleteElement failed")
		}
		editor.SelectElement(nil, trans)
	}
	for _, gui := range editor.editorGUIs {
		err := gui.ConceptDeleted(elID, trans)
		if err != nil {
			errors.Wrap(err, "Editor.DeleteElement failed")
		}
	}
	return nil
}

// DiagramDisplayed ensures that the indicated diagram is on the list of currently displayed diagrams
func (editor *Editor) DiagramDisplayed(id string) {
	var found bool
	for _, foundID := range editor.settings.OpenDiagrams {
		if foundID == id {
			found = true
		}
	}
	if !found {
		trans, new := editor.GetTransaction()
		if new {
			defer editor.EndTransaction()
		}
		editor.UpdateOpenDiagrams(append(editor.settings.OpenDiagrams, id), trans)
	}
}

// DiagramDisplayRemoved ensures that the indicated diagram is not on the list of currently displayed diagrams
func (editor *Editor) DiagramDisplayRemoved(id string) {
	for i, foundID := range editor.settings.OpenDiagrams {
		if foundID == id {
			trans, new := editor.GetTransaction()
			if new {
				defer editor.EndTransaction()
			}
			editor.UpdateOpenDiagrams(append(editor.settings.OpenDiagrams[:i], editor.settings.OpenDiagrams[i+1:]...), trans)
			return
		}
	}
}

// DiagramSelected ensures that the setting is updated to indicate the currently displayed diagram
func (editor *Editor) DiagramSelected(id string) {
	trans, new := editor.GetTransaction()
	if new {
		defer editor.EndTransaction()
	}
	editor.UpdateCurrentDiagram(id, trans)
}

// EndTransaction releases the transaction locks and clears the in-progress transaction
func (editor *Editor) EndTransaction() {
	if editor.inProgressTransaction != nil {
		editor.inProgressTransaction.ReleaseLocks()
		editor.inProgressTransaction = nil
	}
}

// FileLoaded is used to inform the CrlEditor that a file has been loaded
func (editor *Editor) FileLoaded(el *core.Concept, trans *core.Transaction) {
	for _, editorGUI := range editor.editorGUIs {
		editorGUI.FileLoaded(el, trans)
	}
}

// GetCurrentSelection returns the Element that is the current selection in the editor
func (editor *Editor) GetCurrentSelection() *core.Concept {
	return editor.currentSelection
}

// GetCurrentSelectionID returns the ConceptID of the currently selected Element
func (editor *Editor) GetCurrentSelectionID(trans *core.Transaction) string {
	if editor.currentSelection == nil || reflect.ValueOf(editor.currentSelection).IsNil() {
		return ""
	}
	return editor.currentSelection.GetConceptID(trans)
}

// GetDefaultDomainLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultDomainLabel() string {
	editor.settings.DefaultDomainLabelCount++
	countString := strconv.Itoa(editor.settings.DefaultDomainLabelCount)
	return "Domain" + countString
}

// GetDefaultDiagramLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultDiagramLabel() string {
	editor.settings.DefaultDiagramLabelCount++
	countString := strconv.Itoa(editor.settings.DefaultDiagramLabelCount)
	return "Diagram" + countString
}

// GetDefaultElementLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultElementLabel() string {
	editor.settings.DefaultElementLabelCount++
	countString := strconv.Itoa(editor.settings.DefaultElementLabelCount)
	return "Element" + countString
}

// GetDefaultLiteralLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultLiteralLabel() string {
	editor.settings.DefaultLiteralLabelCount++
	countString := strconv.Itoa(editor.settings.DefaultLiteralLabelCount)
	return "Literal" + countString
}

// GetDefaultReferenceLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultReferenceLabel() string {
	editor.settings.DefaultReferenceLabelCount++
	countString := strconv.Itoa(editor.settings.DefaultReferenceLabelCount)
	return "Reference" + countString
}

// GetDefaultRefinementLabel increments the default label count and returns a label containing the new count
func (editor *Editor) GetDefaultRefinementLabel() string {
	editor.settings.DefaultRefinementLabelCount++
	countString := strconv.Itoa(editor.settings.DefaultRefinementLabelCount)
	return "Refinement" + countString
}

// GetDiagramManager returns the diagram manager
func (editor *Editor) GetDiagramManager() *DiagramManager {
	return editor.diagramManager
}

// GetDropDiagramReferenceAsLink returns true if dropped references are shown as links
func (editor *Editor) GetDropDiagramReferenceAsLink(trans *core.Transaction) bool {
	return editor.userPreferences.DropDiagramReferenceAsLink
}

// GetDropDiagramRefinementAsLink returns true if dropped refinements are shown as links
func (editor *Editor) GetDropDiagramRefinementAsLink(trans *core.Transaction) bool {
	return editor.userPreferences.DropDiagramRefinementAsLink
}

// GetExitRequested returns true if exit has been requested
func (editor *Editor) GetExitRequested() bool {
	return editor.exitRequested
}

// getNoSaveDomains returns a map of the editor domains that should not be saved
func (editor *Editor) getNoSaveDomains(trans *core.Transaction) map[string]*core.Concept {
	noSaveDomains := make(map[string]*core.Concept)
	for _, editor := range editor.editorGUIs {
		editor.GetNoSaveDomains(noSaveDomains, trans)
	}
	return noSaveDomains
}

// GetSettings returns the editor settings
func (editor *Editor) GetSettings() *Settings {
	return editor.settings
}

// GetTransaction returns the in-progress transaction, if there is one, and otherwise creates a new transaction
func (editor *Editor) GetTransaction() (*core.Transaction, bool) {
	if editor.inProgressTransaction != nil {
		return editor.inProgressTransaction, false
	}
	editor.inProgressTransaction = editor.GetUofD().NewTransaction()
	return editor.inProgressTransaction, true
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

// getSettingsPath returns the path to the user preferences
func (editor *Editor) getSettingsPath() string {
	return editor.GetWorkspacePath() + "/.settings"
}

// GetWorkspacePath return the path to the current workspace
func (editor *Editor) GetWorkspacePath() string {
	return editor.userPreferences.WorkspacePath
}

// Initialize initializes the uOfD, workspace manager, and all registered editorGUIs
func (editor *Editor) Initialize(workspacePath string, promptWorkspaceSelection bool) error {
	editor.settings = &Settings{}
	editor.settings.OpenDiagrams = []string{}
	editor.inProgressTransaction = nil
	editor.uOfDManager.Initialize()
	trans, isNew := editor.GetTransaction()
	if isNew {
		defer editor.EndTransaction()
	}
	crldatatypesdomain.BuildCrlDataTypesDomain(editor.GetUofD(), trans)
	crldatastructuresdomain.BuildCrlDataStructuresDomain(editor.GetUofD(), trans)
	crldiagramdomain.BuildCrlDiagramDomain(editor.GetUofD(), trans)
	err := crlmapsdomain.BuildCrlMapsDomain(editor.GetUofD(), trans)
	if err != nil {
		return errors.Wrap(err, "Editor.Initialize failed")
	}
	uOfD := trans.GetUniverseOfDiscourse()
	editor.transientCurrentDiagram, _ = uOfD.NewOwnedLiteral(core.Transient, "TransientCurrentDiagram", trans)
	editor.transientDisplayedDiagrams, _ = uOfD.NewOwnedLiteral(core.Transient, "TransientDisplayedDiagrams", trans)
	editor.transientSelection, _ = uOfD.NewOwnedLiteral(core.Transient, "TransientSelection", trans)
	if editor.workspaceManager == nil {
		editor.workspaceManager = NewCrlWorkspaceManager(editor)
	}
	editor.workspaceManager.Initialize()
	editor.workspaceManager.LoadUserPreferences(workspacePath)
	if workspacePath != "" && editor.userPreferences.WorkspacePath != workspacePath {
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
	editor.cutBuffer = make(map[string]*core.Concept)

	for _, editorGUI := range editor.editorGUIs {
		err = editorGUI.Initialize(trans)
		if err != nil {
			return errors.Wrap(err, "Editor.Initialize failed")
		}
	}

	if editor.userPreferences.WorkspacePath != "" {
		err = editor.workspaceManager.LoadWorkspace(trans)
	}
	if err != nil {
		return errors.Wrap(err, "Editor.Initialize failed")
	}

	err = editor.RefreshGUI(trans)
	if err != nil {
		return errors.Wrap(err, "Editor.Initialize failed")
	}

	editor.uOfDManager.UofD.SetRecordingUndo(true)
	return nil
}

// RefreshGUI tells all GUIs to initialize their state
func (editor *Editor) RefreshGUI(trans *core.Transaction) error {
	for _, gui := range editor.editorGUIs {
		err := gui.RefreshGUI(trans)
		if err != nil {
			return errors.Wrap(err, "Editor.RefreshGUI failed")
		}
	}
	return nil
}

// IsDiagramDisplayed returns true if the diagram is in the list of displayed diagrams
func (editor *Editor) IsDiagramDisplayed(diagramID string, trans *core.Transaction) bool {
	for _, openDiagramID := range editor.settings.OpenDiagrams {
		if openDiagramID == diagramID {
			return true
		}
	}
	return false
}

// LoadWorkspace tells the editor to load the workspace
func (editor *Editor) LoadWorkspace(trans *core.Transaction) error {
	err := editor.workspaceManager.LoadWorkspace(trans)
	if err != nil {
		return errors.Wrap(err, "Editor.LoadWorkspace failed")
	}
	return nil
}

// OpenWorkspace sets the path to the folder to be used as a workspace. It is the implementation of a request from the client.
func (editor *Editor) OpenWorkspace(trans *core.Transaction) error {
	if editor.userPreferences.WorkspacePath != "" {
		return errors.New("Cannot open another workspace in the same editor - close existing workspace first")
	}
	path, err := editor.SelectWorkspace()
	if err != nil {
		return err
	}
	err = editor.workspaceManager.LoadSettings(trans)
	if err != nil {
		return err
	}
	return editor.Initialize(path, false)
}

// Redo performs an undo on the editor.editor.GetUofD() and refreshes the interface
func (editor *Editor) Redo(trans *core.Transaction) error {
	editor.undoRedoInProgress = true
	editor.GetUofD().Redo(trans)
	editor.SelectElementUsingIDString(editor.transientSelection.GetLiteralValue(trans), trans)
	editor.settings.Selection = editor.transientSelection.GetLiteralValue(trans)
	var recoveredOpenDiagrams []string
	recoverdJSONOpenDiagrams := editor.transientDisplayedDiagrams.GetLiteralValue(trans)
	json.Unmarshal([]byte(recoverdJSONOpenDiagrams), &recoveredOpenDiagrams)
	editor.settings.OpenDiagrams = recoveredOpenDiagrams
	editor.settings.CurrentDiagram = editor.transientCurrentDiagram.GetLiteralValue(trans)
	err := editor.RefreshGUI(trans)
	if err != nil {
		return errors.Wrap(err, "Editor.Redo failed")
	}
	editor.undoRedoInProgress = false
	return nil
}

// RemoveDiagramFromDisplayedList removes the diagramID from the list of displayed diagrams
func (editor *Editor) RemoveDiagramFromDisplayedList(diagramID string, trans *core.Transaction) {
	if editor.IsDiagramDisplayed(diagramID, trans) {
		for i, openDiagramID := range editor.settings.OpenDiagrams {
			if openDiagramID == diagramID {
				newList := make([]string, 0)
				newList = append(newList, editor.settings.OpenDiagrams[:i]...)
				editor.UpdateOpenDiagrams(append(newList, editor.settings.OpenDiagrams[i+1:]...), trans)
				return
			}
		}
	}
}

// SaveSettings saves the settings to the workspace
func (editor *Editor) SaveSettings() error {
	f, err := os.OpenFile(editor.getSettingsPath(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	serializedSettings, err2 := json.Marshal(editor.settings)
	if err2 != nil {
		return err2
	}
	_, err = f.Write(serializedSettings)
	if err != nil {
		return err
	}
	err = f.Truncate(int64(len(serializedSettings)))
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
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
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

// SaveWorkspace saves the workspace
func (editor *Editor) SaveWorkspace(trans *core.Transaction) error {
	err := editor.workspaceManager.SaveWorkspace(trans)
	if err != nil {
		return errors.Wrap(err, "Editor.SaveWorkspace failed")
	}
	err = editor.SaveSettings()
	if err != nil {
		return errors.Wrap(err, "Editor.SaveWorkspace failed to save settings")
	}
	return nil
}

// SelectElement selects the indicated Element in the tree, displays the Element in the Properties window, and selects it in the
// current diagram (if present).
func (editor *Editor) SelectElement(el *core.Concept, trans *core.Transaction) error {
	if editor.currentSelection != el {
		editor.currentSelection = el
		for _, gui := range editor.editorGUIs {
			err := gui.ConceptSelected(el, trans)
			if err != nil {
				return errors.Wrap(err, "Editor.SelectElement failed")
			}
		}
		elID := ""
		if el != nil {
			elID = el.GetConceptID(trans)
		}
		if !editor.undoRedoInProgress {
			editor.transientSelection.SetLiteralValue(elID, trans)
		}
		editor.settings.Selection = elID
	}
	return nil
}

// SelectElementUsingIDString selects the Element whose ConceptID matches the supplied string
func (editor *Editor) SelectElementUsingIDString(id string, trans *core.Transaction) error {
	foundElement := editor.GetUofD().GetElement(id)
	if foundElement == nil && id != "" {
		return errors.New("In Editor.SelectElementUsingIDString, element was not found")
	}
	return editor.SelectElement(foundElement, trans)
}

// SelectWorkspace opens a dialog for the user to select a workspace
func (editor *Editor) SelectWorkspace() (string, error) {
	return dialog.Directory().Title("Select a directory for your workspace").Browse()
}

// SetDropDiagramReferenceAsLink returns true if dropped references are shown as links
func (editor *Editor) SetDropDiagramReferenceAsLink(value bool, trans *core.Transaction) {
	editor.userPreferences.DropDiagramReferenceAsLink = value
}

// SetDropDiagramRefinementAsLink returns true if dropped refinements are shown as links
func (editor *Editor) SetDropDiagramRefinementAsLink(value bool, trans *core.Transaction) {
	editor.userPreferences.DropDiagramRefinementAsLink = value
}

// SetExitRequested informs the Editor that exit has been requested. Intended to be used by the GUI
func (editor *Editor) SetExitRequested() {
	editor.exitRequested = true
}

// SetSelectionDefinition is a convenience method for setting the Definition of the currently selected Element
func (editor *Editor) SetSelectionDefinition(definition string, trans *core.Transaction) {
	editor.currentSelection.SetDefinition(definition, trans)
}

// SetSelectionLabel is a convenience method for setting the Label of the currently selected Element
func (editor *Editor) SetSelectionLabel(name string, trans *core.Transaction) {
	editor.currentSelection.SetLabel(name, trans)
}

// SetSelectionURI is a convenience method for setting the URI of the curretly selected Element
func (editor *Editor) SetSelectionURI(uri string, trans *core.Transaction) {
	editor.currentSelection.SetURI(uri, trans)
}

// SetWorkspacePath sets the user's preference WorkspacePath value.
func (editor *Editor) SetWorkspacePath(path string) error {
	editor.userPreferences.WorkspacePath = path
	return editor.SaveUserPreferences()
}

// Undo performs an undo on the editor.GetUofD() and refreshes the interface
func (editor *Editor) Undo(trans *core.Transaction) error {
	editor.undoRedoInProgress = true
	editor.GetUofD().Undo(trans)
	editor.SelectElementUsingIDString(editor.transientSelection.GetLiteralValue(trans), trans)
	editor.settings.Selection = editor.transientSelection.GetLiteralValue(trans)
	var recoveredOpenDiagrams []string
	recoverdJSONOpenDiagrams := editor.transientDisplayedDiagrams.GetLiteralValue(trans)
	json.Unmarshal([]byte(recoverdJSONOpenDiagrams), &recoveredOpenDiagrams)
	editor.settings.OpenDiagrams = recoveredOpenDiagrams
	editor.settings.CurrentDiagram = editor.transientCurrentDiagram.GetLiteralValue(trans)
	for _, gui := range editor.editorGUIs {
		err := gui.RefreshGUI(trans)
		if err != nil {
			return errors.Wrap(err, "Editor.Undo failed")
		}
	}
	editor.undoRedoInProgress = false
	return nil
}

// UpdateCurrentDiagram updates the setting value for current diagram in a way that allows the change to be undone
func (editor *Editor) UpdateCurrentDiagram(diagramID string, trans *core.Transaction) {
	editor.settings.CurrentDiagram = diagramID
	if !editor.undoRedoInProgress {
		editor.transientCurrentDiagram.SetLiteralValue(diagramID, trans)
	}
}

// UpdateOpenDiagrams updates the setting value for open diagrams in a way that allows the change to be undone
func (editor *Editor) UpdateOpenDiagrams(diagramIDs []string, trans *core.Transaction) {
	editor.settings.OpenDiagrams = diagramIDs
	if !editor.undoRedoInProgress {
		jsonOpenDiagrams, _ := json.Marshal(editor.settings.OpenDiagrams)
		stringifiedJSONOpenDiagrams := string(jsonOpenDiagrams)
		editor.transientDisplayedDiagrams.SetLiteralValue(stringifiedJSONOpenDiagrams, trans)
	}
}

// EditorGUI is the interface for all CrlEditors, independent of implementation technology
type EditorGUI interface {
	CloseDiagramView(diagramID string, trans *core.Transaction) error
	DisplayDiagram(diagram *crldiagramdomain.CrlDiagram, trans *core.Transaction) error
	ConceptDeleted(elID string, trans *core.Transaction) error
	ConceptSelected(el *core.Concept, trans *core.Transaction) error
	FileLoaded(el *core.Concept, trans *core.Transaction)
	GetNoSaveDomains(noSaveDomains map[string]*core.Concept, trans *core.Transaction)
	Initialize(trans *core.Transaction) error
	RefreshGUI(trans *core.Transaction) error
}
