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

// ByLabel implements the sort.Interface for []string based on the string
// being the ID of a core.Element sorted by the Label of the Element
type ByLabel []string

func (a ByLabel) Len() int      { return len(a) }
func (a ByLabel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLabel) Less(i, j int) bool {
	uOfD := FyneGUISingleton.editor.GetUofD()
	return uOfD.GetElementLabel(a[i]) < uOfD.GetElementLabel(a[j])
}

// FyneTreeManager is the manager of the fyne tree in the CrlFyneEditor
type FyneTreeManager struct {
	fyneGUI      *CrlEditorFyneGUI
	tree         *widget.Tree
	uofdObserver uOfDObserver
}

// NewFyneTreeManager returns an initialized FyneTreeManager
func NewFyneTreeManager(fyneGUI *CrlEditorFyneGUI) *FyneTreeManager {
	treeManager := &FyneTreeManager{}
	treeManager.fyneGUI = fyneGUI
	treeManager.tree = widget.NewTree(GetChildUIDs, IsBranch, CreateNode, UpdateNode)
	treeManager.tree.ExtendBaseWidget(treeManager.tree)
	treeManager.tree.OnSelected = func(uid string) { treeManager.onNodeSelected(uid) }
	treeManager.tree.Show()
	treeManager.uofdObserver = *newUofDObserver(treeManager)
	return treeManager
}

func (ftm *FyneTreeManager) ElementSelected(uid string) {
	ftm.tree.ScrollTo(uid)
	ftm.tree.Select(uid)
	trans, new := ftm.fyneGUI.editor.GetTransaction()
	if new {
		defer ftm.fyneGUI.editor.EndTransaction()
	}
	ftm.openParentsRecursively(uid, trans)
}

func (ftm *FyneTreeManager) initialize() {
	ftm.tree.Refresh()
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

func GetChildUIDs(parentUid string) []string {
	var ids []string
	if parentUid == "" {
		uOfD := FyneGUISingleton.editor.GetUofD()
		if uOfD != nil {
			ids = uOfD.GetRootElementIDs()
		}
	} else {
		iterator := FyneGUISingleton.editor.GetUofD().GetConceptsOwnedConceptIDs(parentUid).Iterator()
		for member := range iterator.C {
			ids = append(ids, member.(string))
		}
	}
	sort.Sort(ByLabel(ids))
	return ids
}

func IsBranch(uid string) bool {
	// All elements in the uOfD are potentially branches
	return true
}

func CreateNode(branch bool) fyne.CanvasObject {
	return newTreeNode()
}

func UpdateNode(uid string, branch bool, node fyne.CanvasObject) {
	tn := node.(*treeNode)
	tn.id = uid
	tn.icon.SetResource(getIconResourceByID(uid))
	if uid == "" {
		tn.label.SetText("uOfD")
	} else {
		conceptBinding := GetConceptStateBinding(uid)
		structBinding := *conceptBinding.GetBoundData()
		if structBinding != nil {
			labelItem, _ := structBinding.GetItem("Label")
			tn.label.Bind(labelItem.(binding.String))
		}
	}
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

var _ desktop.Mouseable = (*treeNode)(nil)

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

func (tn *treeNode) MouseDown(event *desktop.MouseEvent) {
	switch event.Button {
	case desktop.LeftMouseButton:
		FyneGUISingleton.treeManager.tree.Select(tn.id)
	case desktop.RightMouseButton:
		addElement := fyne.NewMenuItem("Add Child Element", func() {
			trans, isNew := FyneGUISingleton.editor.GetTransaction()
			if isNew {
				defer FyneGUISingleton.editor.EndTransaction()
			}
			uOfD := trans.GetUniverseOfDiscourse()
			newElement, _ := uOfD.NewElement(trans)
			newElement.SetLabel(FyneGUISingleton.editor.GetDefaultElementLabel(), trans)
			newElement.SetOwningConceptID(tn.id, trans)
			FyneGUISingleton.editor.SelectElement(newElement, trans)
		})
		addDiagram := fyne.NewMenuItem("Add Child Diagram", func() {
			trans, isNew := FyneGUISingleton.editor.GetTransaction()
			if isNew {
				defer FyneGUISingleton.editor.EndTransaction()
			}
			uOfD := trans.GetUniverseOfDiscourse()
			newElement, _ := uOfD.CreateReplicateAsRefinementFromURI(crldiagramdomain.CrlDiagramURI, trans)
			newElement.SetLabel(FyneGUISingleton.editor.GetDefaultDiagramLabel(), trans)
			newElement.SetOwningConceptID(tn.id, trans)
			FyneGUISingleton.editor.SelectElement(newElement, trans)
		})
		addLiteral := fyne.NewMenuItem("Add Child Literal", func() {
			trans, isNew := FyneGUISingleton.editor.GetTransaction()
			if isNew {
				defer FyneGUISingleton.editor.EndTransaction()
			}
			uOfD := trans.GetUniverseOfDiscourse()
			newLiteral, _ := uOfD.NewLiteral(trans)
			newLiteral.SetLabel(FyneGUISingleton.editor.GetDefaultLiteralLabel(), trans)
			newLiteral.SetOwningConceptID(tn.id, trans)
			FyneGUISingleton.editor.SelectElement(newLiteral, trans)
		})
		addReference := fyne.NewMenuItem("Add Child Reference", func() {
			trans, isNew := FyneGUISingleton.editor.GetTransaction()
			if isNew {
				defer FyneGUISingleton.editor.EndTransaction()
			}
			uOfD := trans.GetUniverseOfDiscourse()
			newReference, _ := uOfD.NewReference(trans)
			newReference.SetLabel(FyneGUISingleton.editor.GetDefaultReferenceLabel(), trans)
			newReference.SetOwningConceptID(tn.id, trans)
			FyneGUISingleton.editor.SelectElement(newReference, trans)
		})
		addRefinement := fyne.NewMenuItem("Add Child Refinement", func() {
			trans, isNew := FyneGUISingleton.editor.GetTransaction()
			if isNew {
				defer FyneGUISingleton.editor.EndTransaction()
			}
			uOfD := trans.GetUniverseOfDiscourse()
			newRefinement, _ := uOfD.NewRefinement(trans)
			newRefinement.SetLabel(FyneGUISingleton.editor.GetDefaultRefinementLabel(), trans)
			newRefinement.SetOwningConceptID(tn.id, trans)
			FyneGUISingleton.editor.SelectElement(newRefinement, trans)
		})
		childMenu := fyne.NewMenu("Add Child", addDiagram, addElement, addLiteral, addReference, addRefinement)
		childMenuItem := fyne.NewMenuItem("Add Child", func() {
			popup := widget.NewPopUpMenu(childMenu, FyneGUISingleton.window.Canvas())
			popup.ShowAtPosition(event.AbsolutePosition)
		})
		deleteElementItem := fyne.NewMenuItem("Delete", func() {
			trans, isNew := FyneGUISingleton.editor.GetTransaction()
			if isNew {
				defer FyneGUISingleton.editor.EndTransaction()
			}
			FyneGUISingleton.editor.DeleteElement(tn.id, trans)
			FyneGUISingleton.editor.SelectElement(nil, trans)
		})
		topMenuItems := []*fyne.MenuItem{}
		topMenuItems = append(topMenuItems, childMenuItem)
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
		nodeElement := trans.GetUniverseOfDiscourse().GetElement(tn.id)
		if crldiagramdomain.IsDiagram(nodeElement, trans) && !FyneGUISingleton.editor.IsDiagramDisplayed(tn.id, trans) {
			showDiagramItem := fyne.NewMenuItem("Show Diagram", func() {
				FyneGUISingleton.editor.GetDiagramManager().DisplayDiagram(tn.id, trans)
			})
			topMenuItems = append(topMenuItems, showDiagramItem)
		}
		topMenuItems = append(topMenuItems, deleteElementItem)
		topMenu := fyne.NewMenu("Top Menu", topMenuItems...)
		popup := widget.NewPopUpMenu(topMenu, FyneGUISingleton.window.Canvas())
		popup.ShowAtPosition(event.AbsolutePosition)
	}
}

func (tn *treeNode) MouseUp(event *desktop.MouseEvent) {

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

type uOfDObserver struct {
	ftm *FyneTreeManager
}

func newUofDObserver(ftm *FyneTreeManager) *uOfDObserver {
	uo := &uOfDObserver{}
	uo.ftm = ftm
	ftm.fyneGUI.editor.GetUofD().Register(uo)
	return uo
}

// Update is the callback for changes to the core diagram
func (uo *uOfDObserver) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
	switch notification.GetNatureOfChange() {
	case core.ConceptRemoved:
		uo.ftm.tree.Refresh()
	}
	return nil
}
