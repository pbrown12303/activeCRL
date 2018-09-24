package editor

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
)

const treeNodeSuffix = "TreeNode"

// The Tree Manager itself
type TreeManager struct {
	manageNodesFunction core.Element
	treeId              string
	rootElements        map[string]core.BaseElement
	uOfD                core.UniverseOfDiscourse
}

func NewTreeManager(uOfD core.UniverseOfDiscourse, treeId string, hl *core.HeldLocks) *TreeManager {
	var treeManager TreeManager
	treeManager.uOfD = uOfD
	treeManager.treeId = treeId
	treeManager.rootElements = make(map[string]core.BaseElement)

	uOfDQuery := jquery.NewJQuery("#uOfD")
	uOfDQuery.Call("jstree", js.M{"core": js.M{"check_callback": true,
		"multiple": false}})

	// Set up the selection callback
	uOfDQuery.On("select_node.jstree", func(e jquery.Event, selection *js.Object) {
		js.Global.Set("treeSelection", selection)
		treeNodeId := selection.Get("node").Get("id").String()
		baseElementId := getIdWithoutSuffix(treeNodeId, "TreeNode")
		log.Printf("Selected node id: %s", treeNodeId)
		log.Printf("Selected base element id: %s", baseElementId)
		CrlEditorSingleton.SelectBaseElementUsingIDString(baseElementId)
	})

	// Set up the tree drag start callback
	uOfDQuery.On("dragstart", func(e jquery.Event, data *js.Object) {
		onTreeDragStart(e, data)
	})

	// Set up the tree view
	var err error
	treeManager.manageNodesFunction, err = core.CreateReplicateAsRefinementFromUri(uOfD, ManageNodesUri, hl)
	if err != nil {
		log.Print(err)
	}
	uOfDReference := core.GetChildElementReferenceWithAncestorUri(treeManager.manageNodesFunction, ManageNodesUofDReferenceUri, hl)
	uOfDReference.SetReferencedElement(uOfD, hl)

	return &treeManager
}

func (tmPtr *TreeManager) AddChildren(be core.BaseElement, hl *core.HeldLocks) {
	switch be.(type) {
	case core.Element:
		parentId := be.GetId(hl)
		jTreeParentId := parentId + "TreeNode"
		for _, child := range be.(core.Element).GetOwnedBaseElements(hl) {
			tmPtr.AddNode(child, jTreeParentId, hl)
			tmPtr.AddChildren(child, hl)
		}
	}
}

func (tmPtr *TreeManager) AddNode(be core.BaseElement, parentId string, hl *core.HeldLocks) {
	id := be.GetId(hl) + treeNodeSuffix
	name := core.GetLabel(be, hl)
	var icon string
	switch be.(type) {
	case core.ElementPointer:
		icon = "/icons/ElementPointerIcon.svg"
	case core.ElementPointerPointer:
		icon = "/icons/ElementPointerPointerIcon.svg"
	case core.ElementPointerReference:
		icon = "/icons/ElementPointerReferenceIcon.svg"
	case core.ElementReference:
		icon = "/icons/ElementReferenceIcon.svg"
	case core.Literal:
		icon = "/icons/LiteralIcon.svg"
	case core.LiteralPointer:
		icon = "/icons/LiteralPointerIcon.svg"
	case core.LiteralPointerPointer:
		icon = "/icons/LiteralPointerPointerIcon.svg"
	case core.LiteralPointerReference:
		icon = "/icons/LiteralPointerReferenceIcon.svg"
	case core.LiteralReference:
		icon = "/icons/LiteralReferenceIcon.svg"
	case core.Refinement:
		icon = "/icons/RefinementIcon.svg"
	case core.Element:
		if IsDiagram(be, hl) {
			icon = "/icons/DiagramIcon.svg"
		} else {
			icon = "/icons/ElementIcon.svg"
		}
	}
	jquery.NewJQuery(tmPtr.treeId).Call("jstree", "create_node", parentId, js.M{"id": id,
		"text": name,
		"icon": icon}, "last")
}

func (tmPtr *TreeManager) getChangeNotificationBelowUofD(changeNotification *core.ChangeNotification) *core.ChangeNotification {
	if changeNotification.GetChangedBaseElement() == tmPtr.uOfD {
		return changeNotification.GetUnderlyingChangeNotification()
	} else if changeNotification.GetUnderlyingChangeNotification() != nil {
		return tmPtr.getChangeNotificationBelowUofD(changeNotification.GetUnderlyingChangeNotification())
	}
	return nil
}

func (tmPtr *TreeManager) InitializeTree(hl *core.HeldLocks) {
	for _, be := range tmPtr.uOfD.GetBaseElements() {
		if core.GetOwningElement(be, hl) == nil {
			tmPtr.AddNode(be, "#", hl)
			tmPtr.AddChildren(be, hl)
		}
	}
}

func IsDiagram(be core.BaseElement, hl *core.HeldLocks) bool {
	switch be.(type) {
	case core.Element:
		return CrlEditorSingleton.GetUofD().IsRefinementOf(be.(core.Element), CrlEditorSingleton.GetDiagramManager().abstractDiagram, hl)
	}
	return false
}

func getIdWithoutSuffix(stringWithSuffix string, suffix string) string {
	if len(stringWithSuffix) > len(suffix) && stringWithSuffix[len(stringWithSuffix)-len(suffix):] == suffix {
		return stringWithSuffix[:len(stringWithSuffix)-len(suffix)]
	} else {
		return ""
	}
}

func onTreeDragStart(e jquery.Event, data *js.Object) {
	parentId := e.Get("target").Get("parentElement").Get("id").String()
	log.Printf("On Tree Drag Start called, ParentId = %s", parentId)
	selectedBaseElementId := getIdWithoutSuffix(parentId, treeNodeSuffix)
	log.Printf("selectedBaseElementId = %s", selectedBaseElementId)
	be := CrlEditorSingleton.GetUofD().GetBaseElement(selectedBaseElementId)
	CrlEditorSingleton.SetTreeDragSelection(be)
}
