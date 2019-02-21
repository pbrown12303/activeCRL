package editor

import (
	"github.com/pbrown12303/activeCRL/core"

	//	"github.com/satori/go.uuid"
	"log"
)

// TreeViewsURI identifies the TreeViews concept
var TreeViewsURI = editorURI + "/TreeViews"

// ManageTreeNodesURI identifies the ManageNodes concept
var ManageTreeNodesURI = TreeViewsURI + "/ManageTreeNodes"

// ManageNodesUofDReferenceURI identifies the ManageNodesUofDReference
var ManageNodesUofDReferenceURI = ManageTreeNodesURI + "/UofDReference"

// ViewNodeURI identifies the ViewNode concept
var ViewNodeURI = TreeViewsURI + "/ViewNode"

// ViewNodeElementReferenceURI identifies the ViewNodeElementReference concept
var ViewNodeElementReferenceURI = ViewNodeURI + "/ElementReference"

// treeViewManageNodes() is the callback function that manaages the tree view when base elements in the Universe of Discourse change.
// The changes being sought are the addition, removal, and re-parenting of base elements and the changes in their names.
func treeViewManageNodes(instance core.Element, changeNotification *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocks()

	// Tracing
	if core.AdHocTrace == true {
		log.Printf("In treeViewManageNodes()")
	}

	treeManager := CrlEditorSingleton.getTreeManager()

	switch changeNotification.GetNatureOfChange() {
	case core.IndicatedConceptChanged:
		underlyingChange := changeNotification.GetUnderlyingChange()
		if underlyingChange == nil {
			log.Printf("treeViewManageNodes called with IndicatedConceptChanged but no underlying chanage")
			return
		}
		switch underlyingChange.GetNatureOfChange() {
		case core.IndicatedConceptChanged:
			secondUnderlyingChange := underlyingChange.GetUnderlyingChange()
			if secondUnderlyingChange == nil {
				log.Printf("treeViewManageNodes called with IndicatedConceptChanged but no underlying chanage")
				return
			}
			switch secondUnderlyingChange.GetNatureOfChange() {
			case core.UofDConceptAdded:
				changedElement := secondUnderlyingChange.GetConceptState()
				treeManager.addNode(changedElement, hl)
			case core.UofDConceptChanged:
				thirdUnderlyingChange := secondUnderlyingChange.GetUnderlyingChange()
				if thirdUnderlyingChange == nil {
					log.Printf("treeViewManageNodes called with UofDConceptChanged but no thirdUnderlyingChange chanage")
					return
				}
				changedElement := thirdUnderlyingChange.GetConceptState()
				treeManager.changeNode(changedElement, hl)
			case core.UofDConceptRemoved:
				changedElement := secondUnderlyingChange.GetConceptState()
				treeManager.removeNode(changedElement, hl)
			}
		}
	}

	// // this is the notification we are interested in
	// // Find the changed base element
	// changedElement := changeNotification.GetReportingElement()
	// changedElementID := changeNotification.GetReportingElementID()

	// // Now see if the node view exists
	// changedElementNodeViewID := changedElementID + treeNodeSuffix
	// changedElementNodeView := jquery.NewJQuery(treeManager.treeID).Call("jstree", "get_node", changedElementNodeViewID)

	// if changedElementNodeView.Length == 0 {
	// 	// Node does not exist. Create it

	// 	// Tracing
	// 	if core.AdHocTrace == true {
	// 		log.Printf("----- Node does not exist")
	// 	}

	// 	// First, determine whether this is a root element or a child
	// 	var parentTreeNodeID string
	// 	parentTreeNodeID = "#"
	// 	parent := changedElement.GetOwningConcept(hl)
	// 	if parent != nil {
	// 		parentTreeNodeID = parent.GetConceptID(hl) + treeNodeSuffix
	// 	}
	// 	treeManager.AddNode(changedElement, parentTreeNodeID, hl)
	// } else {
	// 	// Node exists - update it

	// 	// Tracing
	// 	if core.AdHocTrace == true {
	// 		log.Printf("----- Node exists")
	// 	}

	// 	// See if parent has changed
	// 	currentTreeParentID := changedElementNodeView.Attr("parent")
	// 	currentParent := changedElement.GetOwningConcept(hl)
	// 	currentParentID := "#" // the jstree version of a nil parent
	// 	if currentParent != nil {
	// 		currentParentID = currentParent.GetConceptID(hl) + treeNodeSuffix
	// 	}
	// 	if currentTreeParentID != currentParentID {
	// 		jquery.NewJQuery(treeManager.treeID).Call("jstree", "cut", changedElementID)
	// 		jquery.NewJQuery(treeManager.treeID).Call("jstree", "paste", currentParentID, "last")
	// 	}

	// 	// See if the name has changed
	// 	changedBaseElementLabel := changedElement.GetLabel(hl)
	// 	if changedElementNodeView.Attr("text") != changedBaseElementLabel {
	// 		jquery.NewJQuery(treeManager.treeID).Call("jstree", "rename_node", changedElementNodeViewID, changedBaseElementLabel)
	// 	}
	// }
}

// BuildTreeViews builds the concepts related to TreeViews and adds them to the uOfD
func BuildTreeViews(conceptSpace core.Element, hl *core.HeldLocks) {
	uOfD := conceptSpace.GetUniverseOfDiscourse(hl)

	// TreeViews
	treeViews, _ := uOfD.NewElement(hl, TreeViewsURI)
	treeViews.SetLabel("TreeViews", hl)
	treeViews.SetURI(TreeViewsURI, hl)
	treeViews.SetOwningConcept(conceptSpace, hl)
	treeViews.SetIsCore(hl)

	// ManageNodes
	manageNodes, _ := uOfD.NewElement(hl, ManageTreeNodesURI)
	manageNodes.SetLabel("ManageTreeNodes", hl)
	manageNodes.SetURI(ManageTreeNodesURI, hl)
	manageNodes.SetOwningConcept(treeViews, hl)
	manageNodes.SetIsCore(hl)
	// ManageNodes UofD Reference
	uOfDReference, _ := uOfD.NewReference(hl, ManageNodesUofDReferenceURI)
	uOfDReference.SetLabel("UofDReference", hl)
	uOfDReference.SetURI(ManageNodesUofDReferenceURI, hl)
	uOfDReference.SetOwningConcept(manageNodes, hl)
	uOfDReference.SetIsCore(hl)

	// ViewNode
	viewNode, _ := uOfD.NewElement(hl, ViewNodeURI)
	viewNode.SetLabel("ViewNode", hl)
	viewNode.SetURI(ViewNodeURI, hl)
	viewNode.SetOwningConcept(treeViews, hl)
	viewNode.SetIsCore(hl)
	// ViewNode BaseElementReference
	reference, _ := uOfD.NewReference(hl, ViewNodeElementReferenceURI)
	reference.SetLabel("ElementReference", hl)
	reference.SetURI(ViewNodeElementReferenceURI, hl)
	reference.SetOwningConcept(viewNode, hl)
	reference.SetIsCore(hl)
}

func registerTreeViewFunctions(uOfD core.UniverseOfDiscourse) {
	uOfD.AddFunction(ManageTreeNodesURI, treeViewManageNodes)
}
