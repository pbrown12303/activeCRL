package editor

import (
	//	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"github.com/satori/go.uuid"
	"log"
	//	"strconv"
	"sync"
)

var TreeViewsUri string = EditorUri + "/TreeViews"
var ManageNodesUri string = TreeViewsUri + "/ManageNodes"
var ManageNodesUofDReferenceUri string = ManageNodesUri + "/UofDReference"
var ViewNodeUri string = TreeViewsUri + "/ViewNode"
var ViewNodeBaseElementReferenceUri string = ViewNodeUri + "/BaseElementReference"

// treeViewManageNodes() is the callback function that manaages the tree view when base elements in the Universe of Discourse change.
// The changes being sought are the addition, removal, and re-parenting of base elements and the changes in their names.
func treeViewManageNodes(instance core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	if CrlEditorSingleton == nil {
		log.Printf("CrlEditorSingleton is nil")
	}
	treeManager := CrlEditorSingleton.GetTreeManager()
	if treeManager == nil {
		log.Printf("TreeManager is nil")
	}
	log.Printf("treeViewManageNodes called, notifications length: %d", len(changeNotifications))
	for _, changeNotification := range changeNotifications {
		underlyingChangeNotification := treeManager.getChangeNotificationBelowUofD(changeNotification)
		if underlyingChangeNotification != nil {
			// this is the notification we are interested in
			// Find the changed base element
			changedBaseElement := underlyingChangeNotification.GetChangedBaseElement()
			changedBaseElementId := changedBaseElement.GetId(hl).String()

			// Tracing
			//			underlyingChangeNotification.Print("Change Notification "+strconv.Itoa(i)+": ", hl)
			//			js.Global.Set("changedBaseElementId", changedBaseElementId)

			// Now see if the node view exists
			changedBaseElementNodeViewId := changedBaseElementId + treeNodeSuffix
			changedBaseElementNodeView := jquery.NewJQuery(treeManager.treeId).Call("jstree", "get_node", changedBaseElementNodeViewId)

			// Tracing
			//			js.Global.Set("treeManagerJquery", jquery.NewJQuery(treeManager.treeId))
			//			js.Global.Set("changedBaseElementNodeView", changedBaseElementNodeView)

			if changedBaseElementNodeView.Length == 0 {
				// Node does not exist. Create it
				// First, determine whether this is a root element or a child
				var parentTreeNodeId string
				parentTreeNodeId = "#"
				parent := core.GetOwningElement(changedBaseElement, hl)
				if parent != nil {
					parentTreeNodeId = parent.GetId(hl).String() + treeNodeSuffix
				}
				treeManager.AddNode(changedBaseElement, parentTreeNodeId, hl)
			} else {
				// Node exists - update it
				// See if parent has changed
				currentTreeParentId := changedBaseElementNodeView.Attr("parent")
				currentParent := core.GetOwningElement(changedBaseElement, hl)
				currentParentId := "#" // the jstree version of a nil parent
				if currentParent != nil {
					currentParentId = currentParent.GetId(hl).String() + treeNodeSuffix
				}
				if currentTreeParentId != currentParentId {
					jquery.NewJQuery(treeManager.treeId).Call("jstree", "cut", changedBaseElementId)
					jquery.NewJQuery(treeManager.treeId).Call("jstree", "paste", currentParentId, "last")
				}

				// See if the name has changed
				changedBaseElementLabel := core.GetLabel(changedBaseElement, hl)
				if changedBaseElementNodeView.Attr("text") != changedBaseElementLabel {
					jquery.NewJQuery(treeManager.treeId).Call("jstree", "rename_node", changedBaseElementNodeViewId, changedBaseElementLabel)
				}
			}
		}
	}
}

func BuildTreeViews(conceptSpace core.Element, hl *core.HeldLocks) {
	uOfD := conceptSpace.GetUniverseOfDiscourse(hl)

	// TreeViews
	treeViews := uOfD.NewElement(hl, TreeViewsUri)
	core.SetLabel(treeViews, "TreeViews", hl)
	core.SetUri(treeViews, TreeViewsUri, hl)
	core.SetOwningElement(treeViews, conceptSpace, hl)

	// ManageNodes
	manageNodes := uOfD.NewElement(hl, ManageNodesUri)
	core.SetLabel(manageNodes, "ManageNodes", hl)
	core.SetUri(manageNodes, ManageNodesUri, hl)
	core.SetOwningElement(manageNodes, treeViews, hl)
	// ManageNodes UofD Reference
	uOfDReference := uOfD.NewElementReference(hl, ManageNodesUofDReferenceUri)
	core.SetLabel(uOfDReference, "UofDReference", hl)
	core.SetUri(uOfDReference, ManageNodesUofDReferenceUri, hl)
	core.SetOwningElement(uOfDReference, manageNodes, hl)

	// ViewNode
	viewNode := uOfD.NewElement(hl, ViewNodeUri)
	core.SetLabel(viewNode, "ViewNode", hl)
	core.SetUri(viewNode, ViewNodeUri, hl)
	core.SetOwningElement(viewNode, treeViews, hl)
	// ViewNode BaseElementReference
	baseElementReference := uOfD.NewBaseElementReference(hl, ViewNodeBaseElementReferenceUri)
	core.SetLabel(baseElementReference, "BaseElementReference", hl)
	core.SetUri(baseElementReference, ViewNodeBaseElementReferenceUri, hl)
	core.SetOwningElement(baseElementReference, viewNode, hl)
}

func registerTreeViewFunctions() {
	core.GetCore().AddFunction(ManageNodesUri, treeViewManageNodes)
}
