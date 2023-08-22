package crleditorfynegui

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/images"
)

/*
************************** ByLabel ******************************
 */

// ByLabel implements the sort.Interface for []string based on the string
// being the ID of a core.Element sorted by the Label of the Element
type ByLabel []string

func (a ByLabel) Len() int      { return len(a) }
func (a ByLabel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLabel) Less(i, j int) bool {
	uOfD := FyneGUISingleton.editor.GetUofD()
	iLabel := uOfD.GetElementLabel(a[i])
	jLabel := uOfD.GetElementLabel(a[j])
	return iLabel+a[i] < jLabel+a[j]
}

/*
********************* FyneTreeManager ******************************
 */

// FyneTreeManager is the manager of the fyne tree in the CrlFyneEditor
type FyneTreeManager struct {
	fyneGUI      *CrlEditorFyneGUI
	tree         *widget.Tree
	uofdObserver *uOfDObserver
	treeNodes    map[string]*treeNode
}

// NewFyneTreeManager returns an initialized FyneTreeManager
func NewFyneTreeManager(fyneGUI *CrlEditorFyneGUI) *FyneTreeManager {
	ftm := &FyneTreeManager{}
	ftm.fyneGUI = fyneGUI
	ftm.tree = widget.NewTree(GetChildUIDs, IsBranch, CreateNode, UpdateNode)
	ftm.tree.ExtendBaseWidget(ftm.tree)
	ftm.tree.OnSelected = func(uid string) { ftm.onNodeSelected(uid) }
	ftm.tree.Show()
	ftm.tree.Select("")
	ftm.initialize()
	return ftm
}

// ElementSelected causes the tree manager to select the indicated tree entry
func (ftm *FyneTreeManager) ElementSelected(uid string) {
	if uid == "" {
		ftm.tree.UnselectAll()
		return
	}
	trans, new := ftm.fyneGUI.editor.GetTransaction()
	if new {
		defer ftm.fyneGUI.editor.EndTransaction()
	}
	ftm.openParentsRecursively(uid, trans)
	ftm.tree.ScrollTo(uid)
	ftm.tree.Select(uid)
	ftm.tree.Refresh()
}

func (ftm *FyneTreeManager) initialize() {
	if ftm.uofdObserver == nil {
		ftm.uofdObserver = newUofDObserver(ftm)
	} else {
		ftm.uofdObserver.initialize()
	}
	ftm.treeNodes = make(map[string]*treeNode)
}

func (ftm *FyneTreeManager) onNodeSelected(id string) {
	trans, new := ftm.fyneGUI.editor.GetTransaction()
	if new {
		defer ftm.fyneGUI.editor.EndTransaction()
	}
	ftm.fyneGUI.editor.SelectElementUsingIDString(id, trans)
}

func (ftm *FyneTreeManager) openParentsRecursively(childUID string, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	crlElement := uOfD.GetElement(childUID)
	if crlElement != nil {
		parent := crlElement.GetOwningConcept(trans)
		if parent != nil {
			parentID := parent.GetConceptID(trans)
			ftm.tree.OpenBranch(parentID)
			ftm.openParentsRecursively(parentID, trans)
		}
	}
}

// ShowElementInTree ensures that the indicated element is shown in the tree
func (ftm *FyneTreeManager) ShowElementInTree(element core.Element) {
	if element != nil {
		trans, new := ftm.fyneGUI.editor.GetTransaction()
		if new {
			defer ftm.fyneGUI.editor.EndTransaction()
		}
		uid := element.GetConceptID(trans)
		ftm.tree.ScrollTo(uid)
		ftm.tree.Select(uid)
		ftm.openParentsRecursively(uid, trans)
	}
}

/*****************************
Tree-defining functions
*/

// GetChildUIDs returns an array of the child UIDs
func GetChildUIDs(parentUID string) []string {
	var ids []string
	if parentUID == "" {
		uOfD := FyneGUISingleton.editor.GetUofD()
		if uOfD != nil {
			ids = uOfD.GetRootElementIDs()
		}
	} else {
		iterator := FyneGUISingleton.editor.GetUofD().GetConceptsOwnedConceptIDs(parentUID).Iterator()
		for member := range iterator.C {
			ids = append(ids, member.(string))
		}
	}
	sort.Sort(ByLabel(ids))
	return ids
}

// IsBranch always returns true - all elements are potentially branches
func IsBranch(uid string) bool {
	// All elements in the uOfD are potentially branches
	return true
}

// CreateNode creates a tree node
func CreateNode(branch bool) fyne.CanvasObject {
	return newTreeNode()
}

// UpdateNode updates the data in the indicated node
func UpdateNode(uid string, branch bool, node fyne.CanvasObject) {
	tn := node.(*treeNode)
	tn.id = uid
	tn.icon.SetResource(getIconResourceByID(uid))
	if uid == "" {
		tn.label.SetText("uOfD")
	} else {
		conceptBinding := FyneGUISingleton.GetConceptStateBinding(uid)
		structBinding := *conceptBinding.GetBoundData()
		if structBinding != nil {
			labelItem, _ := structBinding.GetItem("Label")
			tn.label.Bind(labelItem.(binding.String))
		}
	}
	FyneGUISingleton.treeManager.treeNodes[uid] = tn
	tn.Show()
}

// getIconResourceByID returns the icon image resource to be used in representing the given Element in the tree
func getIconResourceByID(id string) *fyne.StaticResource {
	el := crleditor.CrlEditorSingleton.GetUofD().GetElement(id)
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	return getIconResource(el, trans)
}

// getIconResource returns the icon image resource to be used in representing the given Element in the tree
func getIconResource(el core.Element, trans *core.Transaction) *fyne.StaticResource {
	isDiagram := crldiagramdomain.IsDiagram(el, trans)
	switch el.(type) {
	case core.Reference:
		return images.ResourceReferenceIconPng
	case core.Literal:
		return images.ResourceLiteralIconPng
	case core.Refinement:
		return images.ResourceRefinementIconPng
	case core.Element:
		if isDiagram {
			return images.ResourceDiagramIconPng
		}
		return images.ResourceElementIconPng
	}
	return nil
}

/*
***************** TREE NODE ****************************
 */

var _ fyne.Draggable = (*treeNode)(nil)
var _ fyne.Tappable = (*treeNode)(nil)
var _ fyne.SecondaryTappable = (*treeNode)(nil)

type treeNode struct {
	widget.BaseWidget
	id    string
	icon  *widget.Icon
	label *widget.Label
	box   *fyne.Container
}

func newTreeNode() *treeNode {
	tn := &treeNode{}
	tn.BaseWidget.ExtendBaseWidget(tn)
	tn.icon = widget.NewIcon(images.ResourceElementIconPng)
	tn.label = widget.NewLabel("short")
	tn.box = container.NewHBox(tn.icon, tn.label)
	return tn
}

func (tn *treeNode) CreateRenderer() fyne.WidgetRenderer {
	return newTreeNodeRenderer(tn)
}

func (tn *treeNode) Cursor() desktop.StandardCursor {
	return FyneGUISingleton.activeCursor
}

func (tn *treeNode) DragEnd() {
	if FyneGUISingleton.dragDropTransaction != nil {
		ddt := FyneGUISingleton.dragDropTransaction
		if ddt.diagramID != "" && FyneGUISingleton.dragDropTransaction.currentDiagramMousePosition != fyne.NewPos(-1, -1) {
			trans, isNew := FyneGUISingleton.editor.GetTransaction()
			if isNew {
				defer FyneGUISingleton.editor.EndTransaction()
			}
			trans.GetUniverseOfDiscourse().MarkUndoPoint()
			view, _ := FyneGUISingleton.editor.GetDiagramManager().AddConceptView(ddt.diagramID, ddt.id, float64(ddt.currentDiagramMousePosition.X), float64(ddt.currentDiagramMousePosition.Y), trans)
			fyneDiagram := FyneGUISingleton.diagramManager.GetSelectedDiagram()
			fyneDiagram.SelectDiagramElementNoCallback(view.GetConceptID(trans))
			fyneDiagram.Refresh()
		}
		FyneGUISingleton.dragDropTransaction = nil
	}
}

func (tn *treeNode) Dragged(event *fyne.DragEvent) {
	if FyneGUISingleton.dragDropTransaction == nil {
		FyneGUISingleton.dragDropTransaction = &dragDropTransaction{id: tn.id}
	}
	FyneGUISingleton.activeCursor = desktop.CrosshairCursor
	FyneGUISingleton.windowContent.Refresh()
}

func (tn *treeNode) Tapped(event *fyne.PointEvent) {
	FyneGUISingleton.treeManager.tree.Select(tn.id)
}

func (tn *treeNode) TappedSecondary(event *fyne.PointEvent) {
	addElement := fyne.NewMenuItem("Add Child Element", func() {
		FyneGUISingleton.addElement(tn.id, "")
	})
	addDiagram := fyne.NewMenuItem("Add Child Diagram", func() {
		FyneGUISingleton.addDiagram(tn.id)
	})
	addLiteral := fyne.NewMenuItem("Add Child Literal", func() {
		FyneGUISingleton.addLiteral(tn.id, "")
	})
	addReference := fyne.NewMenuItem("Add Child Reference", func() {
		FyneGUISingleton.addReference(tn.id, "")
	})
	addRefinement := fyne.NewMenuItem("Add Child Refinement", func() {
		FyneGUISingleton.addRefinement(tn.id, "")
	})
	childMenu := fyne.NewMenu("Add Child", addDiagram, addElement, addLiteral, addReference, addRefinement)
	childMenuItem := fyne.NewMenuItem("Add Child", func() {
		popup := widget.NewPopUpMenu(childMenu, FyneGUISingleton.window.Canvas())
		popup.ShowAtPosition(event.AbsolutePosition)
	})
	deleteElementItem := fyne.NewMenuItem("Delete", func() {
		FyneGUISingleton.deleteElement(tn.id)
	})
	topMenuItems := []*fyne.MenuItem{}
	topMenuItems = append(topMenuItems, childMenuItem)
	trans, isNew := FyneGUISingleton.editor.GetTransaction()
	if isNew {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	nodeElement := trans.GetUniverseOfDiscourse().GetElement(tn.id)
	if crldiagramdomain.IsDiagram(nodeElement, trans) {
		showDiagramItem := fyne.NewMenuItem("Show Diagram", func() {
			FyneGUISingleton.displayDiagram(tn.id)
		})
		topMenuItems = append(topMenuItems, showDiagramItem)
	}
	topMenuItems = append(topMenuItems, deleteElementItem)
	topMenu := fyne.NewMenu("Top Menu", topMenuItems...)
	popup := widget.NewPopUpMenu(topMenu, FyneGUISingleton.window.Canvas())
	popup.ShowAtPosition(event.AbsolutePosition)
}

type treeNodeRenderer struct {
	tn *treeNode
}

func newTreeNodeRenderer(tn *treeNode) *treeNodeRenderer {
	tnr := &treeNodeRenderer{}
	tnr.tn = tn
	return tnr
}

func (tnr *treeNodeRenderer) Destroy() {
}

func (tnr *treeNodeRenderer) Layout(size fyne.Size) {
}

func (tnr *treeNodeRenderer) MinSize() fyne.Size {
	return tnr.tn.box.MinSize()
}

func (tnr *treeNodeRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{}
	obj = append(obj, tnr.tn.box)
	return obj
}

func (tnr *treeNodeRenderer) Refresh() {
}

/*
********************** uOfDObserver **************************
 */

type uOfDObserver struct {
	ftm  *FyneTreeManager
	uOfD *core.UniverseOfDiscourse
}

func newUofDObserver(ftm *FyneTreeManager) *uOfDObserver {
	uo := &uOfDObserver{}
	uo.ftm = ftm
	uo.initialize()
	return uo
}

func (uo *uOfDObserver) initialize() {
	if uo.uOfD != nil {
		uo.uOfD.Deregister(uo)
	}
	uo.uOfD = uo.ftm.fyneGUI.editor.GetUofD()
	uo.uOfD.Register(uo)
}

// Update is the callback for changes to the core diagram
func (uo *uOfDObserver) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
	switch notification.GetNatureOfChange() {
	case core.ConceptRemoved:
		uo.ftm.tree.Refresh()
	}
	return nil
}
