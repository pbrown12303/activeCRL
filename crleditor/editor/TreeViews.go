package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/core"

	//	"github.com/satori/go.uuid"
	"log"
)

// TreeViewsURI identifies the TreeViews concept
var TreeViewsURI = EditorURI + "/TreeViews"

// ManageNodesURI identifies the ManageNodes concept
var ManageNodesURI = TreeViewsURI + "/ManageNodes"

// ManageNodesUofDReferenceURI identifies the ManageNodesUofDReference
var ManageNodesUofDReferenceURI = ManageNodesURI + "/UofDReference"

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

	if CrlEditorSingleton == nil {
		log.Printf("CrlEditorSingleton is nil")
	}
	treeManager := CrlEditorSingleton.GetTreeManager()
	if treeManager == nil {
		log.Printf("TreeManager is nil")
	}
	underlyingChangeNotification := treeManager.getChangeNotificationBelowUofD(changeNotification)
	if underlyingChangeNotification != nil {
		// this is the notification we are interested in
		// Find the changed base element
		changedElement := underlyingChangeNotification.GetReportingElement()
		changedElementID := changedElement.GetConceptID(hl)

		// Tracing
		if core.AdHocTrace == true {
			underlyingChangeNotification.Print("Change Notification: ", hl)
			js.Global.Set("changedElementID", changedElementID)
		}

		// Now see if the node view exists
		changedElementNodeViewID := changedElementID + treeNodeSuffix
		changedElementNodeView := jquery.NewJQuery(treeManager.treeID).Call("jstree", "get_node", changedElementNodeViewID)

		// Tracing
		if core.AdHocTrace == true {
			js.Global.Set("treeManagerJquery", jquery.NewJQuery(treeManager.treeID))
			js.Global.Set("changedElementNodeView", changedElementNodeView)
		}

		if changedElementNodeView.Length == 0 {
			// Node does not exist. Create it

			// Tracing
			if core.AdHocTrace == true {
				log.Printf("----- Node does not exist")
			}

			// First, determine whether this is a root element or a child
			var parentTreeNodeID string
			parentTreeNodeID = "#"
			parent := changedElement.GetOwningConcept(hl)
			if parent != nil {
				parentTreeNodeID = parent.GetConceptID(hl) + treeNodeSuffix
			}
			treeManager.AddNode(changedElement, parentTreeNodeID, hl)
		} else {
			// Node exists - update it

			// Tracing
			if core.AdHocTrace == true {
				log.Printf("----- Node exists")
			}

			// See if parent has changed
			currentTreeParentID := changedElementNodeView.Attr("parent")
			currentParent := changedElement.GetOwningConcept(hl)
			currentParentID := "#" // the jstree version of a nil parent
			if currentParent != nil {
				currentParentID = currentParent.GetConceptID(hl) + treeNodeSuffix
			}
			if currentTreeParentID != currentParentID {
				jquery.NewJQuery(treeManager.treeID).Call("jstree", "cut", changedElementID)
				jquery.NewJQuery(treeManager.treeID).Call("jstree", "paste", currentParentID, "last")
			}

			// See if the name has changed
			changedBaseElementLabel := changedElement.GetLabel(hl)
			if changedElementNodeView.Attr("text") != changedBaseElementLabel {
				jquery.NewJQuery(treeManager.treeID).Call("jstree", "rename_node", changedElementNodeViewID, changedBaseElementLabel)
			}
		}
	}
}

// BuildTreeViews builds the concepts related to TreeViews and adds them to the uOfD
func BuildTreeViews(conceptSpace core.Element, hl *core.HeldLocks) {
	uOfD := conceptSpace.GetUniverseOfDiscourse(hl)

	// TreeViews
	treeViews, _ := uOfD.NewElement(hl, TreeViewsURI)
	treeViews.SetLabel("TreeViews", hl)
	treeViews.SetURI(TreeViewsURI, hl)
	treeViews.SetOwningConcept(conceptSpace, hl)

	// ManageNodes
	manageNodes, _ := uOfD.NewElement(hl, ManageNodesURI)
	manageNodes.SetLabel("ManageNodes", hl)
	manageNodes.SetURI(ManageNodesURI, hl)
	manageNodes.SetOwningConcept(treeViews, hl)
	// ManageNodes UofD Reference
	uOfDReference, _ := uOfD.NewReference(hl, ManageNodesUofDReferenceURI)
	uOfDReference.SetLabel("UofDReference", hl)
	uOfDReference.SetURI(ManageNodesUofDReferenceURI, hl)
	uOfDReference.SetOwningConcept(manageNodes, hl)

	// ViewNode
	viewNode, _ := uOfD.NewElement(hl, ViewNodeURI)
	viewNode.SetLabel("ViewNode", hl)
	viewNode.SetURI(ViewNodeURI, hl)
	viewNode.SetOwningConcept(treeViews, hl)
	// ViewNode BaseElementReference
	reference, _ := uOfD.NewReference(hl, ViewNodeElementReferenceURI)
	reference.SetLabel("BaseElementReference", hl)
	reference.SetURI(ViewNodeElementReferenceURI, hl)
	reference.SetOwningConcept(viewNode, hl)
}

func registerTreeViewFunctions() {
	//	core.GetCore().AddFunction(ManageNodesURI, treeViewManageNodes)
}
