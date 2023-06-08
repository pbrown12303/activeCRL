package crleditorfynegui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crleditor"
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
	// The following attributes are kept for testing purposes
	// File Menu Items
	newDomainItem           *fyne.MenuItem
	selectConceptWithIDItem *fyne.MenuItem
	saveWorkspaceItem       *fyne.MenuItem
	closeWorkspaceItem      *fyne.MenuItem
	clearWorkspaceItem      *fyne.MenuItem
	openWorkspaceItem       *fyne.MenuItem
	userPreferencesItem     *fyne.MenuItem
	// Edit Menu Items
	undoItem *fyne.MenuItem
	redoItem *fyne.MenuItem
	// Debug Menu Items
	debugSettingsItem     *fyne.MenuItem
	displayCallGraphsItem *fyne.MenuItem
	// Help Menu Items
	helpItem *fyne.MenuItem
	// Main Menu Items
	mainMenu  *fyne.MainMenu
	fileMenu  *fyne.Menu
	editMenu  *fyne.Menu
	debugMenu *fyne.Menu
	helpMenu  *fyne.Menu
}

// NewFyneGUI returns an initialized FyneGUI
func NewFyneGUI(crlEditor *crleditor.Editor) *CrlEditorFyneGUI {
	fyneGUI := &CrlEditorFyneGUI{}
	fyneGUI.app = app.New()
	initializeFyneGUI(fyneGUI, crlEditor)
	return fyneGUI
}

func initializeFyneGUI(fyneGUI *CrlEditorFyneGUI, crlEditor *crleditor.Editor) {
	FyneGUISingleton = fyneGUI
	fyneGUI.editor = crlEditor
	InitBindings()
	fyneGUI.app.Settings().SetTheme(&fyneGuiTheme{})
	fyneGUI.treeManager = NewFyneTreeManager(fyneGUI)
	fyneGUI.propertyManager = NewFynePropertyManager()
	fyneGUI.diagramManager = NewFyneDiagramManager(fyneGUI)
	fyneGUI.window = fyneGUI.app.NewWindow("Crl Editor    Workspace: " + crlEditor.GetWorkspacePath())
	fyneGUI.buildCrlFyneEditorMenus()
	fyneGUI.window.SetMainMenu(fyneGUI.mainMenu)
	fyneGUI.window.SetMaster()

	leftSide := container.NewVSplit(fyneGUI.treeManager.tree, fyneGUI.propertyManager.properties)
	drawingArea := fyneGUI.diagramManager.GetDrawingArea()

	content := container.NewHSplit(leftSide, drawingArea)

	fyneGUI.window.SetContent(content)
}

// buildCrlFyneEditorMenu builds the main menu for the Crl Fyne Editor
func (fyneGUI *CrlEditorFyneGUI) buildCrlFyneEditorMenus() {
	// File Menu Items
	fyneGUI.newDomainItem = fyne.NewMenuItem("New Domain", func() {
		trans, isNew := fyneGUI.editor.GetTransaction()
		if isNew {
			defer fyneGUI.editor.EndTransaction()
		}
		uOfD := fyneGUI.editor.GetUofD()
		uOfD.MarkUndoPoint()
		cs, _ := uOfD.NewElement(trans)
		cs.SetLabel(fyneGUI.editor.GetDefaultDomainLabel(), trans)
		fyneGUI.editor.SelectElement(cs, trans)
	})
	fyneGUI.selectConceptWithIDItem = fyne.NewMenuItem("Select Concept With ID", nil)
	fyneGUI.saveWorkspaceItem = fyne.NewMenuItem("Save Workspace", func() {
		trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
		if isNew {
			defer fyneGUI.editor.EndTransaction()
		}
		crleditor.CrlEditorSingleton.SaveWorkspace(trans)
	})
	fyneGUI.closeWorkspaceItem = fyne.NewMenuItem("Close Workspace", func() {
		trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
		if isNew {
			defer fyneGUI.editor.EndTransaction()
		}
		crleditor.CrlEditorSingleton.CloseWorkspace(trans)
	})
	fyneGUI.clearWorkspaceItem = fyne.NewMenuItem("Clear Workspace", func() {
		trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
		if isNew {
			defer fyneGUI.editor.EndTransaction()
		}
		crleditor.CrlEditorSingleton.ClearWorkspace(trans)
	})
	fyneGUI.openWorkspaceItem = fyne.NewMenuItem("Open Workspace", func() {
		err := crleditor.CrlEditorSingleton.OpenWorkspace()
		if err != nil {
			errorMsg := widget.NewLabel(err.Error())
			popup := widget.NewPopUp(errorMsg, fyneGUI.window.Canvas())
			popup.Show()
		}
	})
	fyneGUI.userPreferencesItem = fyne.NewMenuItem("UserPreferences", func() { fmt.Println("User Preferences") })

	// Edit Menu Items
	fyneGUI.undoItem = fyne.NewMenuItem("Undo", nil)
	fyneGUI.redoItem = fyne.NewMenuItem("Redo", nil)

	// Debug Menu Items
	fyneGUI.debugSettingsItem = fyne.NewMenuItem("Debug Settings", nil)
	fyneGUI.displayCallGraphsItem = fyne.NewMenuItem("Display Call Graphs", nil)

	// Help Menu Items
	fyneGUI.helpItem = fyne.NewMenuItem("Help", func() { fmt.Println("Help Menu") })

	// Main Menu
	fyneGUI.fileMenu = fyne.NewMenu("File", fyneGUI.newDomainItem, fyne.NewMenuItemSeparator(), fyneGUI.saveWorkspaceItem, fyneGUI.closeWorkspaceItem, fyneGUI.clearWorkspaceItem, fyneGUI.openWorkspaceItem, fyne.NewMenuItemSeparator(), fyneGUI.userPreferencesItem)
	fyneGUI.editMenu = fyne.NewMenu("Edit", fyneGUI.selectConceptWithIDItem, fyneGUI.undoItem, fyneGUI.redoItem)
	fyneGUI.debugMenu = fyne.NewMenu("Debug", fyneGUI.debugSettingsItem, fyneGUI.displayCallGraphsItem)
	fyneGUI.helpMenu = fyne.NewMenu("Help", fyneGUI.helpItem)

	fyneGUI.mainMenu = fyne.NewMainMenu(fyneGUI.fileMenu, fyneGUI.editMenu, fyneGUI.debugMenu, fyneGUI.helpMenu)
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
	gui.GetWindow().SetTitle("Crl Editor         Workspace: " + gui.editor.GetWorkspacePath())
	gui.treeManager.initialize()
	gui.diagramManager.initialize()
	for _, openDiagramID := range gui.editor.GetSettings().OpenDiagrams {
		diagram := gui.editor.GetUofD().GetElement(openDiagramID)
		if diagram == nil {
			log.Printf("In FyneGui.initializeClientState: Failed to load diagram with ID: %s", openDiagramID)
		} else {
			err := gui.diagramManager.displayDiagram(diagram, hl)
			if err != nil {
				return errors.Wrap(err, "In FyneGUI.initializeClientState diagram "+diagram.GetLabel(hl)+" did not display")
			}
		}
	}
	return nil
}

// func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
// 	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
// 		focused.TypedShortcut(s)
// 	}
// }
