package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/satori/go.uuid"
	"log"
)

const treeNodeSuffix = "TreeNode"

// Data structures used to set up the jstree
type jstree struct {
	*js.Object
	core    *jstreeCore `js:"core"`
	plugins []string    `js:"plugins"`
}

type jstreeCore struct {
	*js.Object
	check_callback bool `js:"check_callback"`
	multiple       bool `js:"multiple"`
}

type jstreeNode struct {
	*js.Object
	id   string `js:"id"`
	name string `js:"text"`
	icon string `js:"icon"`
}

// The Tree Manager itself
type TreeManager struct {
	manageNodesFunction core.Element
	treeId              string
	rootElements        map[uuid.UUID]core.BaseElement
	uOfD                core.UniverseOfDiscourse
}

func NewTreeManager(uOfD core.UniverseOfDiscourse, treeId string, hl *core.HeldLocks) *TreeManager {
	var treeManager TreeManager
	treeManager.uOfD = uOfD
	treeManager.treeId = treeId
	treeManager.rootElements = make(map[uuid.UUID]core.BaseElement)

	// Set up the jstree
	coreData := &jstreeCore{Object: js.Global.Get("Object").New()}
	coreData.check_callback = true
	coreData.multiple = false
	jstreeData := &jstree{Object: js.Global.Get("Object").New()}
	jstreeData.core = coreData
	//	jstreeData.plugins = []string{"dnd"}
	uOfDQuery := jquery.NewJQuery("#uOfD")
	uOfDQuery.Call("jstree", jstreeData)

	// Set up the selection callback
	uOfDQuery.On("select_node.jstree", func(e jquery.Event, selection *js.Object) {
		js.Global.Set("treeSelection", selection)
		treeNodeId := selection.Get("node").Get("id").String()
		baseElementId := getIdWithoutSuffix(treeNodeId, "TreeNode")
		log.Printf("Selected node id: %s", treeNodeId)
		log.Printf("Selected base element id: %s", baseElementId)
		CrlEditorSingleton.SelectBaseElementUsingIdString(baseElementId)
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
		parentId := be.GetId(hl).String()
		jTreeParentId := parentId + "TreeNode"
		for _, child := range be.(core.Element).GetOwnedBaseElements(hl) {
			tmPtr.AddNode(child, jTreeParentId, hl)
			tmPtr.AddChildren(child, hl)
		}
	}
}

func (tmPtr *TreeManager) AddNode(be core.BaseElement, parentId string, hl *core.HeldLocks) {
	nodeData := &jstreeNode{Object: js.Global.Get("Object").New()}
	nodeData.id = be.GetId(hl).String() + treeNodeSuffix
	nodeData.name = core.GetName(be, hl)
	switch be.(type) {
	case core.ElementPointer:
		nodeData.icon = "/icons/ElementPointerIcon.svg"
	case core.ElementPointerPointer:
		nodeData.icon = "/icons/ElementPointerPointerIcon.svg"
	case core.ElementPointerReference:
		nodeData.icon = "/icons/ElementPointerReferenceIcon.svg"
	case core.ElementReference:
		nodeData.icon = "/icons/ElementReferenceIcon.svg"
	case core.Literal:
		nodeData.icon = "/icons/LiteralIcon.svg"
	case core.LiteralPointer:
		nodeData.icon = "/icons/LiteralPointerIcon.svg"
	case core.LiteralPointerPointer:
		nodeData.icon = "/icons/LiteralPointerPointerIcon.svg"
	case core.LiteralPointerReference:
		nodeData.icon = "/icons/LiteralPointerReferenceIcon.svg"
	case core.LiteralReference:
		nodeData.icon = "/icons/LiteralReferenceIcon.svg"
	case core.Refinement:
		nodeData.icon = "/icons/RefinementIcon.svg"
	case core.Element:
		if IsDiagram(be, hl) {
			nodeData.icon = "/icons/DiagramIcon.svg"
		} else {
			nodeData.icon = "/icons/ElementIcon.svg"
		}
	}
	jquery.NewJQuery(tmPtr.treeId).Call("jstree", "create_node", parentId, nodeData, "last")
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
		return CrlEditorSingleton.uOfD.IsRefinementOf(be.(core.Element), CrlEditorSingleton.diagramManager.abstractDiagram, hl)
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
	selectedBaseElementUUID, err := uuid.FromString(selectedBaseElementId)
	if err == nil {
		be := CrlEditorSingleton.uOfD.GetBaseElement(selectedBaseElementUUID)
		CrlEditorSingleton.SetTreeDragSelection(be)
	}
}
