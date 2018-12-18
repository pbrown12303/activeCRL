package editor

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/core"
)

const treeNodeSuffix = "TreeNode"

// TreeManager manages the client's tree display of the uOfD
type TreeManager struct {
	manageNodesFunction core.Element
	treeID              string
	rootElements        map[string]core.Element
	uOfD                core.UniverseOfDiscourse
}

// NewTreeManager creates an instance of the TreeManager
func NewTreeManager(uOfD core.UniverseOfDiscourse, treeID string, hl *core.HeldLocks) *TreeManager {
	var treeManager TreeManager
	treeManager.uOfD = uOfD
	treeManager.treeID = treeID
	treeManager.rootElements = make(map[string]core.Element)

	uOfDQuery := jquery.NewJQuery("#uOfD")
	uOfDQuery.Call("jstree", js.M{"core": js.M{"check_callback": true,
		"multiple": false}})

	// Set up the selection callback
	uOfDQuery.On("select_node.jstree", func(e jquery.Event, selection *js.Object) {
		js.Global.Set("treeSelection", selection)
		treeNodeID := selection.Get("node").Get("id").String()
		elementID := getIDWithoutSuffix(treeNodeID, "TreeNode")
		log.Printf("Selected node id: %s", treeNodeID)
		log.Printf("Selected base element id: %s", elementID)
		CrlEditorSingleton.SelectElementUsingIDString(elementID)
	})

	// Set up the tree drag start callback
	uOfDQuery.On("dragstart", func(e jquery.Event, data *js.Object) {
		onTreeDragStart(e, data)
	})

	// Set up the tree view
	var err error
	treeManager.manageNodesFunction, err = uOfD.CreateReplicateAsRefinementFromURI(ManageNodesURI, hl)
	if err != nil {
		log.Print(err)
	}
	uOfDReference := treeManager.manageNodesFunction.GetFirstChildReferenceWithAbstractionURI(ManageNodesUofDReferenceURI, hl)
	uOfDReference.SetReferencedConcept(uOfD, hl)

	return &treeManager
}

// AddChildren adds the OwnedConcepts of the supplied Element to the client's tree
func (tmPtr *TreeManager) AddChildren(el core.Element, hl *core.HeldLocks) {
	switch el.(type) {
	case core.Element:
		parentID := el.GetConceptID(hl)
		jTreeParentID := parentID + "TreeNode"
		for _, child := range *el.GetOwnedConcepts(hl) {
			tmPtr.AddNode(child, jTreeParentID, hl)
			tmPtr.AddChildren(child, hl)
		}
	}
}

// AddNode adds a node to the tree
func (tmPtr *TreeManager) AddNode(el core.Element, parentID string, hl *core.HeldLocks) {
	id := el.GetConceptID(hl) + treeNodeSuffix
	name := el.GetLabel(hl)
	var icon string
	switch el.(type) {
	case core.Reference:
		icon = "/icons/ElementReferenceIcon.svg"
	case core.Literal:
		icon = "/icons/LiteralIcon.svg"
	case core.Refinement:
		icon = "/icons/RefinementIcon.svg"
	case core.Element:
		if IsDiagram(el, hl) {
			icon = "/icons/DiagramIcon.svg"
		} else {
			icon = "/icons/ElementIcon.svg"
		}
	}
	jquery.NewJQuery(tmPtr.treeID).Call("jstree", "create_node", parentID, js.M{"id": id,
		"text": name,
		"icon": icon}, "last")
}

func (tmPtr *TreeManager) getChangeNotificationBelowUofD(changeNotification *core.ChangeNotification) *core.ChangeNotification {
	if changeNotification.GetReportingElement() == tmPtr.uOfD {
		return changeNotification.GetUnderlyingChange()
	} else if changeNotification.GetUnderlyingChange() != nil {
		return tmPtr.getChangeNotificationBelowUofD(changeNotification.GetUnderlyingChange())
	}
	return nil
}

// InitializeTree initializes the display of the tree in the client
func (tmPtr *TreeManager) InitializeTree(hl *core.HeldLocks) {
	for _, el := range *tmPtr.uOfD.GetElements() {
		if el.GetOwningConcept(hl) == nil {
			tmPtr.AddNode(el, "#", hl)
			tmPtr.AddChildren(el, hl)
		}
	}
}

// IsDiagram returns true if the supplied element is a crldiagram
func IsDiagram(el core.Element, hl *core.HeldLocks) bool {
	switch el.(type) {
	case core.Element:
		return el.HasAbstraction(CrlEditorSingleton.GetDiagramManager().abstractDiagram, hl)
	}
	return false
}

func getIDWithoutSuffix(stringWithSuffix string, suffix string) string {
	if len(stringWithSuffix) > len(suffix) && stringWithSuffix[len(stringWithSuffix)-len(suffix):] == suffix {
		return stringWithSuffix[:len(stringWithSuffix)-len(suffix)]
	}
	return ""
}

func onTreeDragStart(e jquery.Event, data *js.Object) {
	parentID := e.Get("target").Get("parentElement").Get("id").String()
	log.Printf("On Tree Drag Start called, ParentId = %s", parentID)
	selectedBaseElementID := getIDWithoutSuffix(parentID, treeNodeSuffix)
	log.Printf("selectedBaseElementID = %s", selectedBaseElementID)
	el := CrlEditorSingleton.GetUofD().GetElement(selectedBaseElementID)
	CrlEditorSingleton.SetTreeDragSelection(el)
}
