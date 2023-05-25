package crleditorfynegui

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	fyneGUI *FyneGUI
	tree    *widget.Tree
}

// NewFyneTreeManager returns an initialized FyneTreeManager
func NewFyneTreeManager(fyneGUI *FyneGUI) *FyneTreeManager {
	var treeManager FyneTreeManager
	treeManager.fyneGUI = fyneGUI
	treeManager.tree = widget.NewTree(GetChildUIDs, IsBranch, CreateNode, UpdateNode)
	treeManager.tree.ExtendBaseWidget(treeManager.tree)
	treeManager.tree.OnSelected = func(uid string) { treeManager.onNodeSelected(uid) }
	treeManager.tree.Show()
	return &treeManager
}

func (ftm *FyneTreeManager) ElementSelected(uid string) {
	ftm.tree.ScrollTo(uid)
	ftm.tree.Select(uid)
	trans, new := ftm.fyneGUI.editor.GetTransaction()
	if new {
		defer trans.ReleaseLocks()
	}
	ftm.openParentsRecursively(uid, trans)
}

func (ftm *FyneTreeManager) onNodeSelected(id string) {
	trans, new := ftm.fyneGUI.editor.GetTransaction()
	if new {
		defer trans.ReleaseLocks()
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
	icon := widget.NewIcon(images.ResourceElementIconPng)
	label := widget.NewLabel("short")
	box := container.NewHBox(icon, label)
	return box
}

func UpdateNode(uid string, branch bool, node fyne.CanvasObject) {
	contents := node.(*fyne.Container).Objects
	contents[0].(*widget.Icon).SetResource(getIconResourceByID(uid))
	label := contents[1].(*widget.Label)
	if uid == "" {
		label.SetText("uOfD")
	} else {
		conceptBinding := GetConceptStateBinding(uid)
		structBinding := *conceptBinding.GetBoundData()
		if structBinding != nil {
			labelItem, _ := structBinding.GetItem("Label")
			label.Bind(labelItem.(binding.String))
		}
	}
	contents[0].Show()
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
