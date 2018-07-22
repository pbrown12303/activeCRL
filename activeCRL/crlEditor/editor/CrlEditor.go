package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/crlDiagram"
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

var EditorUri string = "http://activeCrl.com/crlEditor/Editor"

var CrlEditorSingleton *CrlEditor

type CrlEditor struct {
	currentSelection  core.BaseElement
	treeDragSelection core.BaseElement
	diagramManager    *DiagramManager
	hl                *core.HeldLocks
	propertiesManager *PropertiesManager
	treeManager       *TreeManager
	uOfD              core.UniverseOfDiscourse
	initialized       bool
}

func InitializeCrlEditorSingleton() {
	var wg sync.WaitGroup
	var editor CrlEditor
	editor.initialized = false
	editor.hl = core.NewHeldLocks(&wg)
	defer editor.hl.ReleaseLocks()

	// Set the value of the singleton
	CrlEditorSingleton = &editor

	// Set up the Universe of Discourse
	editor.uOfD = core.NewUniverseOfDiscourse(editor.hl)
	crlDiagram.AddCoreDiagramToUofD(editor.uOfD, editor.hl)
	AddEditorViewsToUofD(editor.uOfD, editor.hl)

	// Set up the diagram manager
	editor.diagramManager = NewDiagramManager()

	// Set up the tree manager
	editor.treeManager = NewTreeManager(editor.uOfD, "#uOfD", editor.hl)
	editor.treeManager.InitializeTree(editor.hl)

	// Set up the properties manager
	editor.propertiesManager = NewPropertiesManager()

	// Add the event listeners
	editorQuery := jquery.NewJQuery("body")
	editorQuery.On("ondrop", func(e jquery.Event, data *js.Object) {
		onEditorDrop(e, data)
	})

	editor.hl.ReleaseLocksAndWait()
	editor.initialized = true
	log.Printf("Editor initialized")
}

func AddEditorViewsToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	conceptSpace := uOfD.GetElementWithUri(EditorUri)
	if conceptSpace == nil {
		conceptSpace = BuildEditorConceptSpace(uOfD, hl)
		if conceptSpace == nil {
			log.Printf("Build of Editor Concept Space failed")
		}
	}
	return conceptSpace
}

func (edPtr *CrlEditor) GetCurrentSelection() core.BaseElement {
	return edPtr.currentSelection
}

func (edPtr *CrlEditor) GetDiagramManager() *DiagramManager {
	return edPtr.diagramManager
}

func (edPtr *CrlEditor) GetTreeDragSelection() core.BaseElement {
	return edPtr.treeDragSelection
}

func (edPtr *CrlEditor) GetTreeManager() *TreeManager {
	return edPtr.treeManager
}

func (edPtr *CrlEditor) IsInitialized() bool {
	return edPtr.initialized
}

func (edPtr *CrlEditor) SelectBaseElement(be core.BaseElement) core.BaseElement {
	edPtr.currentSelection = be
	edPtr.propertiesManager.BaseElementSelected(edPtr.currentSelection, edPtr.hl)
	return edPtr.currentSelection
}

func (edPtr *CrlEditor) SelectBaseElementUsingIdString(id string) core.BaseElement {
	uuid, _ := uuid.FromString(id)
	edPtr.currentSelection = edPtr.uOfD.GetBaseElement(uuid)
	edPtr.propertiesManager.BaseElementSelected(edPtr.currentSelection, edPtr.hl)
	return edPtr.currentSelection
}

func (edPtr *CrlEditor) SetSelectionDefinition(definition string) {
	switch edPtr.currentSelection.(type) {
	case core.Element:
		core.SetDefinition(edPtr.currentSelection.(core.Element), definition, edPtr.hl)
		edPtr.hl.ReleaseLocksAndWait()
	}
}

func (edPtr *CrlEditor) SetSelectionLabel(name string) {
	switch edPtr.currentSelection.(type) {
	case core.Element:
		core.SetLabel(edPtr.currentSelection.(core.Element), name, edPtr.hl)
		edPtr.hl.ReleaseLocksAndWait()
	}
}

func (edPtr *CrlEditor) SetSelectionUri(uri string) {
	core.SetUri(edPtr.currentSelection, uri, edPtr.hl)
	edPtr.hl.ReleaseLocksAndWait()
}

func (edPtr *CrlEditor) SetTreeDragSelection(be core.BaseElement) {
	edPtr.treeDragSelection = be
}

func BuildEditorConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// EditorConceptSpace
	conceptSpace := uOfD.NewElement(hl, EditorUri)
	core.SetLabel(conceptSpace, "EditorConceptSpace", hl)
	core.SetUri(conceptSpace, EditorUri, hl)

	BuildTreeViews(conceptSpace, hl)

	return conceptSpace
}

func onEditorDrop(e jquery.Event, data *js.Object) {
	CrlEditorSingleton.SetTreeDragSelection(nil)
}

func init() {
	registerTreeViewFunctions()
}
