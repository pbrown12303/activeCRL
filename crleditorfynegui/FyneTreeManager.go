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

// TODO Remove this after fyne transaction approach is determined
func GetTransaction() *core.Transaction {
	if crleditor.CrlEditorSingleton != nil && crleditor.CrlEditorSingleton.GetInProgressTransaction() != nil {
		return crleditor.CrlEditorSingleton.GetInProgressTransaction()
	}
	return crleditor.CrlEditorSingleton.GetUofD().NewTransaction()
}

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
	tree *widget.Tree
}

// NewFyneTreeManager returns an initialized FyneTreeManager
func NewFyneTreeManager() *FyneTreeManager {
	var treeManager FyneTreeManager
	treeManager.tree = widget.NewTree(GetChildUIDs, IsBranch, CreateNode, UpdateNode)
	treeManager.tree.ExtendBaseWidget(treeManager.tree)
	treeManager.tree.Show()
	return &treeManager
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
		labelItem, _ := structBinding.GetItem("Label")
		label.Bind(labelItem.(binding.String))
	}
	contents[0].Show()
}

// getIconResourceByID returns the icon image resource to be used in representing the given Element in the tree
func getIconResourceByID(id string) *fyne.StaticResource {
	el := crleditor.CrlEditorSingleton.GetUofD().GetElement(id)
	trans := GetTransaction()
	defer trans.ReleaseLocks()
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