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
************************** IDsSortedByLabel ******************************
 */

// IDsSortedByLabel implements the sort.Interface for []string based on the string
// being the ID of a core.Element sorted by the Label of the Element
type IDsSortedByLabel []string

func (a IDsSortedByLabel) Len() int      { return len(a) }
func (a IDsSortedByLabel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a IDsSortedByLabel) Less(i, j int) bool {
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
	treeNodes    map[string]*fyneTreeNode
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
	ftm.treeNodes = make(map[string]*fyneTreeNode)
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
func (ftm *FyneTreeManager) ShowElementInTree(element core.Concept) {
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
	sort.Sort(IDsSortedByLabel(ids))
	return ids
}

// IsBranch always returns true - all elements are potentially branches
func IsBranch(uid string) bool {
	// All elements in the uOfD are potentially branches
	return true
}

// CreateNode creates a tree node
func CreateNode(branch bool) fyne.CanvasObject {
	return newFyneTreeNode()
}

// UpdateNode updates the data in the indicated node
func UpdateNode(uid string, branch bool, node fyne.CanvasObject) {
	tn := node.(*fyneTreeNode)
	tn.id = uid
	icon := getIconResourceByID(uid)
	if icon != nil {
		tn.icon.SetResource(icon)
	}
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
func getIconResource(el core.Concept, trans *core.Transaction) *fyne.StaticResource {
	isDiagram := crldiagramdomain.IsDiagram(el, trans)
	if el != nil {
		switch el.GetConceptType() {
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
	}
	return nil
}

/*
***************** TREE NODE ****************************
 */

var _ fyne.Draggable = (*fyneTreeNode)(nil)
var _ fyne.Tappable = (*fyneTreeNode)(nil)
var _ fyne.SecondaryTappable = (*fyneTreeNode)(nil)

type fyneTreeNode struct {
	widget.BaseWidget
	id    string
	icon  *widget.Icon
	label *widget.Label
	box   *fyne.Container
}

func newFyneTreeNode() *fyneTreeNode {
	tn := &fyneTreeNode{}
	tn.BaseWidget.ExtendBaseWidget(tn)
	tn.icon = widget.NewIcon(images.ResourceElementIconPng)
	// tn.icon.Resize(fyne.NewSize(20, 20))
	tn.label = widget.NewLabel("short")
	tn.box = container.NewHBox(tn.icon, tn.label)
	return tn
}

func (tn *fyneTreeNode) CreateRenderer() fyne.WidgetRenderer {
	return newFyneTreeNodeRenderer(tn)
}

func (tn *fyneTreeNode) Cursor() desktop.StandardCursor {
	return FyneGUISingleton.activeCursor
}

func (tn *fyneTreeNode) DragEnd() {
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

func (tn *fyneTreeNode) Dragged(event *fyne.DragEvent) {
	if FyneGUISingleton.dragDropTransaction == nil {
		FyneGUISingleton.dragDropTransaction = &dragDropTransaction{id: tn.id}
	}
	FyneGUISingleton.activeCursor = desktop.CrosshairCursor
	FyneGUISingleton.windowContent.Refresh()
}

func (tn *fyneTreeNode) Tapped(event *fyne.PointEvent) {
	FyneGUISingleton.treeManager.tree.Select(tn.id)
}

func (tn *fyneTreeNode) TappedSecondary(event *fyne.PointEvent) {
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

type fyneTreeNodeRenderer struct {
	tn *fyneTreeNode
}

func newFyneTreeNodeRenderer(tn *fyneTreeNode) *fyneTreeNodeRenderer {
	tnr := &fyneTreeNodeRenderer{}
	tnr.tn = tn
	return tnr
}

func (tnr *fyneTreeNodeRenderer) Destroy() {
}

func (tnr *fyneTreeNodeRenderer) Layout(size fyne.Size) {
	tnr.tn.box.Resize(tnr.tn.box.MinSize())
}

func (tnr *fyneTreeNodeRenderer) MinSize() fyne.Size {
	return tnr.tn.box.MinSize()
}

func (tnr *fyneTreeNodeRenderer) Objects() []fyne.CanvasObject {
	obj := []fyne.CanvasObject{}
	obj = append(obj, tnr.tn.box)
	return obj
}

func (tnr *fyneTreeNodeRenderer) Refresh() {
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
