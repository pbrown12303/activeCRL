package crleditorfynegui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/diagramwidget"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pkg/errors"
)

var FyneGUISingleton *CrlEditorFyneGUI

// CrlEditorFyneGUI is the Crl Editor built with Fyne
type CrlEditorFyneGUI struct {
	app                 fyne.App
	editor              *crleditor.Editor
	diagramManager      *FyneDiagramManager
	propertyManager     *FynePropertyManager
	treeManager         *FyneTreeManager
	window              fyne.Window
	currentSelectionID  string
	propertiesClipboard *diagramwidget.DiagramElementProperties
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
	debugSettingsItem *fyne.MenuItem
	// Help Menu Items
	helpItem *fyne.MenuItem
	// Main Menu Items
	mainMenu            *fyne.MainMenu
	fileMenu            *fyne.Menu
	editMenu            *fyne.Menu
	debugMenu           *fyne.Menu
	helpMenu            *fyne.Menu
	dragDropTransaction *dragDropTransaction
}

// NewFyneGUI returns an initialized FyneGUI
func NewFyneGUI(crlEditor *crleditor.Editor) *CrlEditorFyneGUI {
	gui := &CrlEditorFyneGUI{}
	gui.app = app.New()
	initializeFyneGUI(gui, crlEditor)
	return gui
}

func initializeFyneGUI(gui *CrlEditorFyneGUI, crlEditor *crleditor.Editor) {
	FyneGUISingleton = gui
	gui.editor = crlEditor
	InitBindings()
	gui.app.Settings().SetTheme(&fyneGuiTheme{})
	gui.treeManager = NewFyneTreeManager(gui)
	gui.propertyManager = NewFynePropertyManager()
	gui.diagramManager = NewFyneDiagramManager(gui)
	gui.window = gui.app.NewWindow("Crl Editor    Workspace: " + crlEditor.GetWorkspacePath())
	gui.buildCrlFyneEditorMenus()
	gui.window.SetMainMenu(gui.mainMenu)
	gui.window.SetMaster()

	leftSide := container.NewVSplit(gui.treeManager.tree, gui.propertyManager.properties)
	drawingArea := gui.diagramManager.GetDrawingArea()

	content := container.NewHSplit(leftSide, drawingArea)

	gui.window.SetContent(content)
}

func (gui *CrlEditorFyneGUI) addDiagram(parentID string) core.Element {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.markUndoPoint()
	uOfD := trans.GetUniverseOfDiscourse()
	newElement, _ := uOfD.CreateReplicateAsRefinementFromURI(crldiagramdomain.CrlDiagramURI, trans)
	newElement.SetLabel(gui.editor.GetDefaultDiagramLabel(), trans)
	newElement.SetOwningConceptID(parentID, trans)
	gui.editor.SelectElement(newElement, trans)
	gui.DisplayDiagram(newElement, trans)
	return newElement
}

func (gui *CrlEditorFyneGUI) addElement(parentID string, label string) core.Element {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.markUndoPoint()
	uOfD := trans.GetUniverseOfDiscourse()
	newElement, _ := uOfD.NewElement(trans)
	if label == "" {
		label = gui.editor.GetDefaultElementLabel()
	}
	newElement.SetLabel(label, trans)
	newElement.SetOwningConceptID(parentID, trans)
	gui.editor.SelectElement(newElement, trans)
	return newElement
}

func (gui *CrlEditorFyneGUI) addLiteral(parentID string, label string) core.Literal {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.markUndoPoint()
	uOfD := trans.GetUniverseOfDiscourse()
	newLiteral, _ := uOfD.NewLiteral(trans)
	if label == "" {
		label = gui.editor.GetDefaultLiteralLabel()
	}
	newLiteral.SetLabel(label, trans)
	newLiteral.SetOwningConceptID(parentID, trans)
	gui.editor.SelectElement(newLiteral, trans)
	return newLiteral
}

func (gui *CrlEditorFyneGUI) addReference(parentID string, label string) core.Reference {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.markUndoPoint()
	uOfD := trans.GetUniverseOfDiscourse()
	newReference, _ := uOfD.NewReference(trans)
	if label == "" {
		label = gui.editor.GetDefaultReferenceLabel()
	}
	newReference.SetLabel(label, trans)
	newReference.SetOwningConceptID(parentID, trans)
	gui.editor.SelectElement(newReference, trans)
	return newReference
}

func (gui *CrlEditorFyneGUI) addRefinement(parentID string, label string) core.Refinement {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.markUndoPoint()
	uOfD := trans.GetUniverseOfDiscourse()
	newRefinement, _ := uOfD.NewRefinement(trans)
	if label == "" {
		label = gui.editor.GetDefaultRefinementLabel()
	}
	newRefinement.SetLabel(label, trans)
	newRefinement.SetOwningConceptID(parentID, trans)
	gui.editor.SelectElement(newRefinement, trans)
	return newRefinement
}

// buildCrlFyneEditorMenu builds the main menu for the Crl Fyne Editor
func (gui *CrlEditorFyneGUI) buildCrlFyneEditorMenus() {
	// File Menu Items
	gui.newDomainItem = fyne.NewMenuItem("New Domain", func() {
		gui.addElement("", gui.editor.GetDefaultDomainLabel())
	})
	gui.selectConceptWithIDItem = fyne.NewMenuItem("Select Concept With ID", func() {
		entryItem := newPastableEntry()
		formItem := widget.NewFormItem("Enter ID of desired element", entryItem)
		dialog.ShowForm("Select Concept", "Select", "Cancel", []*widget.FormItem{formItem}, func(b bool) {
			trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
			if isNew {
				defer gui.editor.EndTransaction()
			}
			gui.editor.SelectElementUsingIDString(entryItem.Text, trans)
		}, gui.window)
	})
	gui.saveWorkspaceItem = fyne.NewMenuItem("Save Workspace", func() {
		trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
		if isNew {
			defer gui.editor.EndTransaction()
		}
		crleditor.CrlEditorSingleton.SaveWorkspace(trans)
	})
	gui.closeWorkspaceItem = fyne.NewMenuItem("Close Workspace", func() {
		trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
		if isNew {
			defer gui.editor.EndTransaction()
		}
		crleditor.CrlEditorSingleton.CloseWorkspace(trans)
	})
	gui.clearWorkspaceItem = fyne.NewMenuItem("Clear Workspace", func() {
		trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
		if isNew {
			defer gui.editor.EndTransaction()
		}
		crleditor.CrlEditorSingleton.ClearWorkspace(trans)
	})
	gui.openWorkspaceItem = fyne.NewMenuItem("Open Workspace", func() {
		trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
		if isNew {
			defer gui.editor.EndTransaction()
		}
		err := crleditor.CrlEditorSingleton.OpenWorkspace(trans)
		if err != nil {
			errorMsg := widget.NewLabel(err.Error())
			popup := widget.NewPopUp(errorMsg, gui.window.Canvas())
			popup.Show()
		}
	})
	gui.userPreferencesItem = fyne.NewMenuItem("UserPreferences", func() {
		preferences := crleditor.UserPreferences{}
		preferences = *gui.editor.GetUserPreferences()
		referenceRadioGroup := widget.NewRadioGroup([]string{"Node", "Link"}, func(s string) {
			if s == "Link" {
				preferences.DropDiagramReferenceAsLink = true
			} else {
				preferences.DropDiagramReferenceAsLink = false
			}
		})
		referenceRadioGroup.Horizontal = true
		if preferences.DropDiagramReferenceAsLink {
			referenceRadioGroup.Selected = "Link"
		} else {
			referenceRadioGroup.Selected = "Node"
		}
		referenceChoice := container.NewHBox(widget.NewLabel("Drop Reference as: "), referenceRadioGroup)
		refinementRadioGroup := widget.NewRadioGroup([]string{"Node", "Link"}, func(s string) {
			if s == "Link" {
				preferences.DropDiagramRefinementAsLink = true
			} else {
				preferences.DropDiagramRefinementAsLink = false
			}
		})
		refinementRadioGroup.Horizontal = true
		if preferences.DropDiagramRefinementAsLink {
			refinementRadioGroup.Selected = "Link"
		} else {
			refinementRadioGroup.Selected = "Node"
		}
		refinementChoice := container.NewHBox(widget.NewLabel("Drop Refinement as: "), refinementRadioGroup)
		vBox := container.NewVBox(referenceChoice, refinementChoice)
		dialog.ShowCustomConfirm("User Preferences", "Save", "Cancel", vBox, func(b bool) {
			if b {
				*gui.editor.GetUserPreferences() = preferences
				gui.editor.SaveUserPreferences()
			}
		}, gui.window)
		fmt.Println("User Preferences")
	})

	// Edit Menu Items
	gui.undoItem = fyne.NewMenuItem("Undo", func() {
		FyneGUISingleton.undo()
	})
	gui.redoItem = fyne.NewMenuItem("Redo", func() {
		FyneGUISingleton.redo()
	})

	// Debug Menu Items
	gui.debugSettingsItem = fyne.NewMenuItem("Debug Settings", func() {
		traceChange := core.TraceChange
		omitManageTreeNodeCalls := core.OmitManageTreeNodesCalls
		omitDiagramRelatedCalls := core.OmitDiagramRelatedCalls
		enableTraceChangeItem := widget.NewCheck("Enable Trace Change", func(value bool) {
			traceChange = value
		})
		enableTraceChangeItem.Checked = traceChange
		omitManageTreeNodeCallsItem := widget.NewCheck("Omit ManageTreeNode calls", func(value bool) {
			omitManageTreeNodeCalls = value
		})
		omitManageTreeNodeCallsItem.Checked = omitManageTreeNodeCalls
		omitDiagramRelatedCallsItem := widget.NewCheck("Omit Diagram-Related calls", func(value bool) {
			omitDiagramRelatedCalls = value
		})
		omitDiagramRelatedCallsItem.Checked = omitDiagramRelatedCalls
		vBox := container.NewVBox(enableTraceChangeItem, omitManageTreeNodeCallsItem, omitDiagramRelatedCallsItem)
		dialog.ShowCustomConfirm("Debug Settings", "Save", "Cancel", vBox, func(b bool) {
			if b {
				core.TraceChange = traceChange
				core.OmitManageTreeNodesCalls = omitManageTreeNodeCalls
				core.OmitDiagramRelatedCalls = omitDiagramRelatedCalls
			}
		}, gui.window)
	})

	// Help Menu Items
	gui.helpItem = fyne.NewMenuItem("Help", func() { fmt.Println("Help Menu") })

	// Main Menu
	gui.fileMenu = fyne.NewMenu("File", gui.newDomainItem, fyne.NewMenuItemSeparator(), gui.saveWorkspaceItem, gui.closeWorkspaceItem, gui.clearWorkspaceItem, gui.openWorkspaceItem, fyne.NewMenuItemSeparator(), gui.userPreferencesItem)
	gui.editMenu = fyne.NewMenu("Edit", gui.selectConceptWithIDItem, gui.undoItem, gui.redoItem)
	gui.debugMenu = fyne.NewMenu("Debug", gui.debugSettingsItem)
	gui.helpMenu = fyne.NewMenu("Help", gui.helpItem)

	gui.mainMenu = fyne.NewMainMenu(gui.fileMenu, gui.editMenu, gui.debugMenu, gui.helpMenu)
}

// CloseDiagramView
func (gui *CrlEditorFyneGUI) CloseDiagramView(diagramID string, trans *core.Transaction) error {
	gui.diagramManager.closeDiagram(diagramID)
	return nil
}

func (gui *CrlEditorFyneGUI) deleteElement(elementID string) {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.editor.DeleteElement(elementID, trans)
	gui.editor.SelectElement(nil, trans)
}

func (gui *CrlEditorFyneGUI) displayDiagram(diagramID string) {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.editor.GetDiagramManager().DisplayDiagram(diagramID, trans)
}

// ElementDeleted - no additional action required
func (gui *CrlEditorFyneGUI) ElementDeleted(elID string, trans *core.Transaction) error {
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
		gui.currentSelectionID = uid
	}
	return nil
}

// DisplayDiagram
func (gui *CrlEditorFyneGUI) DisplayDiagram(diagram core.Element, trans *core.Transaction) error {
	gui.diagramManager.displayDiagram(diagram, trans)
	return nil
}

// FileLoaded - no action required
func (gui *CrlEditorFyneGUI) FileLoaded(el core.Element, trans *core.Transaction) {
}

// GetNoSaveDomains - there aren't any for the CRLEditorFyneGUI
func (gui *CrlEditorFyneGUI) GetNoSaveDomains(noSaveDomains map[string]core.Element, trans *core.Transaction) {
}

// GetWindow returns the main window of the FyneGUI
func (gui *CrlEditorFyneGUI) GetWindow() fyne.Window {
	return gui.window
}

// Initialize
func (gui *CrlEditorFyneGUI) Initialize(trans *core.Transaction) error {
	return nil
}

// InitializeGUI
func (gui *CrlEditorFyneGUI) InitializeGUI(trans *core.Transaction) error {
	gui.GetWindow().SetTitle("Crl Editor         Workspace: " + gui.editor.GetWorkspacePath())
	gui.treeManager.initialize()
	gui.diagramManager.initialize()
	for _, openDiagramID := range gui.editor.GetSettings().OpenDiagrams {
		diagram := gui.editor.GetUofD().GetElement(openDiagramID)
		if diagram == nil {
			log.Printf("In FyneGui.initializeClientState: Failed to load diagram with ID: %s", openDiagramID)
		} else {
			err := gui.diagramManager.displayDiagram(diagram, trans)
			if err != nil {
				return errors.Wrap(err, "In FyneGUI.initializeClientState diagram "+diagram.GetLabel(trans)+" did not display")
			}
		}
	}
	gui.diagramManager.SelectDiagram(gui.editor.GetSettings().CurrentDiagram)
	selectedElement := gui.editor.GetUofD().GetElement(gui.editor.GetSettings().Selection)
	gui.ElementSelected(selectedElement, trans)
	return nil
}

func (gui *CrlEditorFyneGUI) markUndoPoint() {
	uOfD := gui.editor.GetUofD()
	uOfD.MarkUndoPoint()
}

func (gui *CrlEditorFyneGUI) redo() {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.editor.Redo(trans)
}

func (gui *CrlEditorFyneGUI) undo() {
	trans, isNew := gui.editor.GetTransaction()
	if isNew {
		defer gui.editor.EndTransaction()
	}
	gui.editor.Undo(trans)
}

type dragDropTransaction struct {
	id                          string
	diagramID                   string
	currentDiagramMousePosition fyne.Position
}

type pastableEntry struct {
	widget.Entry
}

func newPastableEntry() *pastableEntry {
	pe := &pastableEntry{}
	pe.ExtendBaseWidget(pe)
	return pe
}

func (pe *pastableEntry) TypedShortcut(shortcut fyne.Shortcut) {
	log.Print(shortcut.ShortcutName())
	switch typedShortcut := shortcut.(type) {
	case *fyne.ShortcutPaste:
		typedShortcut.Clipboard = FyneGUISingleton.window.Clipboard()
		pe.SetText(typedShortcut.Clipboard.Content())
		// pe.Refresh()
	}
}

func (pe *pastableEntry) MinSize() fyne.Size {
	if pe.Text == "" {
		return fyne.NewSize(30, 12)
	}
	return pe.Entry.MinSize()
}
