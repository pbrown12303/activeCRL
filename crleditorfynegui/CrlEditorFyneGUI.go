package crleditorfynegui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructuresdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crleditordomain"
	"github.com/pkg/errors"
)

var FyneGUISingleton *CrlEditorFyneGUI

// CrlEditorFyneGUI is the Crl Editor built with Fyne
type CrlEditorFyneGUI struct {
	app                fyne.App
	editor             *crleditor.Editor
	diagramManager     *FyneDiagramManager
	propertyManager    *FynePropertyManager
	treeManager        *FyneTreeManager
	window             fyne.Window
	currentSelectionID string
}

// NewFyneGUI returns an initialized FyneGUI
func NewFyneGUI(crlEditor *crleditor.Editor) *CrlEditorFyneGUI {
	var fyneGUI CrlEditorFyneGUI
	FyneGUISingleton = &fyneGUI
	fyneGUI.editor = crlEditor
	fyneGUI.app = app.New()
	InitBindings()
	fyneGUI.app.Settings().SetTheme(&fyneGuiTheme{})
	fyneGUI.treeManager = NewFyneTreeManager(&fyneGUI)
	fyneGUI.propertyManager = NewFynePropertyManager()
	fyneGUI.diagramManager = NewFyneDiagramManager(&fyneGUI)
	fyneGUI.window = fyneGUI.app.NewWindow("Crl Editor")
	fyneGUI.window.SetMainMenu(buildCrlFyneEditorMenu(fyneGUI.window))
	fyneGUI.window.SetMaster()

	leftSide := container.NewVSplit(fyneGUI.treeManager.tree, fyneGUI.propertyManager.properties)
	drawingArea := fyneGUI.diagramManager.GetDrawingArea()

	content := container.NewHSplit(leftSide, drawingArea)

	fyneGUI.window.SetContent(content)
	return &fyneGUI
}

// buildCrlFyneEditorMenu builds the main menu for the Crl Fyne Editor
func buildCrlFyneEditorMenu(window fyne.Window) *fyne.MainMenu {
	// File Menu Items
	newDomainItem := fyne.NewMenuItem("New Domain", nil)
	selectConceptWithIDItem := fyne.NewMenuItem("Select Concept With ID", nil)
	saveWorkspaceItem := fyne.NewMenuItem("Save Workspace", nil)
	closeWorkspaceItem := fyne.NewMenuItem("Close Workspace", nil)
	clearWorkspaceItem := fyne.NewMenuItem("Clear Workspace", nil)
	openWorkspaceItem := fyne.NewMenuItem("Open Workspace", nil)
	userPreferencesItem := fyne.NewMenuItem("UserPreferences", func() { fmt.Println("User Preferences") })

	// Edit Menu Items
	undoItem := fyne.NewMenuItem("Undo", nil)
	redoItem := fyne.NewMenuItem("Redo", nil)

	// Debug Menu Items
	debugSettingsItem := fyne.NewMenuItem("Debug Settings", nil)
	displayCallGraphsItem := fyne.NewMenuItem("Display Call Graphs", nil)

	// Help Menu Items
	helpItem := fyne.NewMenuItem("Help", func() { fmt.Println("Help Menu") })

	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		fyne.NewMenu("File", newDomainItem, fyne.NewMenuItemSeparator(), saveWorkspaceItem, closeWorkspaceItem, clearWorkspaceItem, openWorkspaceItem, fyne.NewMenuItemSeparator(), userPreferencesItem),
		fyne.NewMenu("Edit", selectConceptWithIDItem, undoItem, redoItem),
		fyne.NewMenu("Debug", debugSettingsItem, displayCallGraphsItem),
		fyne.NewMenu("Help", helpItem),
	)
	return mainMenu
}

// CloseDiagramView
func (gui *CrlEditorFyneGUI) CloseDiagramView(diagramID string, hl *core.Transaction) error {
	gui.diagramManager.closeDiagram(diagramID)
	return nil
}

// ElementDeleted
func (gui *CrlEditorFyneGUI) ElementDeleted(elID string, hl *core.Transaction) error {
	// TODO Implement this
	return nil
}

// ElementSelected
func (gui *CrlEditorFyneGUI) ElementSelected(el core.Element, trans *core.Transaction) error {
	uid := ""
	if el != nil {
		uid = el.GetConceptID(trans)
	}
	if gui.currentSelectionID != uid {
		gui.propertyManager.displayProperties(uid)
		gui.treeManager.ElementSelected(uid)
		gui.diagramManager.ElementSelected(uid, trans)
	}
	return nil
}

// DisplayDiagram
func (gui *CrlEditorFyneGUI) DisplayDiagram(diagram core.Element, trans *core.Transaction) error {
	gui.diagramManager.displayDiagram(diagram, trans)
	return nil
}

// FileLoaded
func (gui *CrlEditorFyneGUI) FileLoaded(el core.Element, hl *core.Transaction) {
	// TODO Implement this
	// noop
}

// GetNoSaveDomains
func (gui *CrlEditorFyneGUI) GetNoSaveDomains(noSaveDomains map[string]core.Element, hl *core.Transaction) {
	// TODO Implement this
	// noop
}

// func (gui *CrlEditorFyneGUI) getUofD() *core.UniverseOfDiscourse {
// 	return gui.editor.GetUofD()
// }

// GetWindow returns the main window of the FyneGUI
func (gui *CrlEditorFyneGUI) GetWindow() fyne.Window {
	return gui.window
}

// Initialize
func (gui *CrlEditorFyneGUI) Initialize(hl *core.Transaction) error {
	return nil
}

// InitializeGUI
func (gui *CrlEditorFyneGUI) InitializeGUI(hl *core.Transaction) error {
	openDiagrams := gui.editor.GetSettings().GetFirstOwnedConceptRefinedFromURI(crleditordomain.EditorOpenDiagramsURI, hl)
	if openDiagrams == nil {
		return errors.New("In FyneGUI.initializeClientState, openDiagrams is nil")
	}
	openDiagramLiteral, err2 := crldatastructuresdomain.GetFirstMemberLiteral(openDiagrams, hl)
	if err2 != nil {
		return errors.Wrap(err2, "In FyneGUI.initializeClientState getting first member literal failed")
	}
	for openDiagramLiteral != nil {
		diagram := gui.editor.GetUofD().GetElement(openDiagramLiteral.GetLiteralValue(hl))
		if diagram == nil {
			log.Printf("In FyneGui.initializeClientState: Failed to load diagram with ID: %s", openDiagramLiteral.GetLiteralValue(hl))
		} else {
			err2 = gui.diagramManager.displayDiagram(diagram, hl)
			if err2 != nil {
				return errors.Wrap(err2, "In FyneGUI.initializeClientState diagram "+diagram.GetLabel(hl)+" did not display")
			}
		}
		openDiagramLiteral, _ = crldatastructuresdomain.GetNextMemberLiteral(openDiagramLiteral, hl)
	}
	return nil
}

// RegisterUofDInitializationFunctions
func (gui *CrlEditorFyneGUI) RegisterUofDInitializationFunctions(uOfDManager *core.UofDManager) error {
	return nil
}

// RegisterUofDPostInitializationFunctions
func (gui *CrlEditorFyneGUI) RegisterUofDPostInitializationFunctions(uOfDManager *core.UofDManager) error {
	return nil
}

// func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
// 	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
// 		focused.TypedShortcut(s)
// 	}
// }
