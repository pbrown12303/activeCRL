package fynegui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crlfynebindings"
)

func getUofD() *core.UniverseOfDiscourse {
	return crleditor.CrlEditorSingleton.GetUofD()
}

// FyneGUI is the Crl Editor built with Fyne
type FyneGUI struct {
	app             fyne.App
	drawingManager  *FyneDrawingManager
	propertyManager *FynePropertyManager
	treeManager     *FyneTreeManager
	window          fyne.Window
}

// NewFyneGUI returns an initialized FyneGUI
func NewFyneGUI(crleditor *crleditor.Editor) *FyneGUI {
	var editor FyneGUI
	editor.app = app.New()
	crlfynebindings.InitBindings()
	editor.app.Settings().SetTheme(&fyneGuiTheme{})
	editor.treeManager = NewFyneTreeManager()
	editor.propertyManager = NewFynePropertyManager()
	editor.drawingManager = NewFyneDrawingManager()
	editor.window = editor.app.NewWindow("Crl Editor")
	editor.window.SetMainMenu(buildCrlFyneEditorMenu(editor.window))
	editor.window.SetMaster()

	leftSide := container.NewVSplit(editor.treeManager.tree, editor.propertyManager.properties)
	drawingArea := editor.drawingManager.GetDrawingArea()

	content := container.NewHSplit(leftSide, drawingArea)

	editor.window.SetContent(content)
	return &editor
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

// GetWindow returns the main window of the FyneGUI
func (gui *FyneGUI) GetWindow() fyne.Window {
	return gui.window
}

// CloseDiagramView
func (gui *FyneGUI) CloseDiagramView(diagramID string, hl *core.Transaction) error {
	return nil
}

// ElementDeleted
func (gui *FyneGUI) ElementDeleted(elID string, hl *core.Transaction) error {
	return nil
}

// ElementSelected
func (gui *FyneGUI) ElementSelected(el core.Element, hl *core.Transaction) error {
	uid := ""
	if el != nil {
		uid = el.GetConceptID(hl)
	}
	gui.propertyManager.displayProperties(uid)
	return nil
}

// FileLoaded
func (gui *FyneGUI) FileLoaded(el core.Element, hl *core.Transaction) {
	// noop
}

// GetNoSaveDomains
func (gui *FyneGUI) GetNoSaveDomains(noSaveDomains map[string]core.Element, hl *core.Transaction) {
	// noop
}

// Initialize
func (gui *FyneGUI) Initialize(hl *core.Transaction) error {
	return nil
}

// InitializeGUI
func (gui *FyneGUI) InitializeGUI(hl *core.Transaction) error {
	return nil
}

// RegisterUofDInitializationFunctions
func (gui *FyneGUI) RegisterUofDInitializationFunctions(uOfDManager *core.UofDManager) error {
	return nil
}

// RegisterUofDPostInitializationFunctions
func (gui *FyneGUI) RegisterUofDPostInitializationFunctions(uOfDManager *core.UofDManager) error {
	return nil
}

// func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
// 	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
// 		focused.TypedShortcut(s)
// 	}
// }
