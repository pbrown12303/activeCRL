package fyneeditor

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"
)

// CrlFyneEditor is the Crl Editor built with Fyne
type CrlFyneEditor struct {
	app    fyne.App
	window fyne.Window
}

// NewCrlFyneEditor returns an initialized CrlFyneEditor
func NewCrlFyneEditor() *CrlFyneEditor {
	var editor CrlFyneEditor
	editor.app = app.New()
	editor.window = editor.app.NewWindow("Crl Editor")
	editor.window.SetMainMenu(buildCrlFyneEditorMenu(editor.window))
	editor.window.SetMaster()

	top := widget.NewLabel("top bar")
	left := canvas.NewText("left", color.White)
	middle := canvas.NewText("content", color.White)
	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(top, nil, left, nil),
		top, left, middle)
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

// GetWindow returns the main window of the CrlFyneEditor
func (editor *CrlFyneEditor) GetWindow() fyne.Window {
	return editor.window
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
