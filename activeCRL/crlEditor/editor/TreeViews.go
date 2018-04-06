package editor

import (
	//	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/satori/go.uuid"
	"log"
	//	"strconv"
	"sync"
)

var TreeViewsUri string = EditorUri + "/TreeViews"
var ManageNodesUri string = TreeViewsUri + "/ManageNodes"
var ManageNodesUofDReferenceUri string = ManageNodesUri + "/UofDReference"
var ViewNodeUri string = TreeViewsUri + "/ViewNode"
var ViewNodeBaseElementReferenceUri string = ViewNodeUri + "/BaseElementReference"

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
	var changedBaseElements map[uuid.UUID]core.BaseElement = make(map[uuid.UUID]core.BaseElement)
	log.Printf("treeViewManageNodes called, notifications length: %d", len(changeNotifications))
	for _, changeNotification := range changeNotifications {
		underlyingChangeNotification := treeManager.getChangeNotificationBelowUofD(changeNotification)
		if underlyingChangeNotification != nil {
			// this is the notification we are interested in
			changedBaseElement := underlyingChangeNotification.GetChangedBaseElement()
			//			underlyingChangeNotification.Print("Change Notification "+strconv.Itoa(i)+": ", hl)
			changedBaseElements[changedBaseElement.GetId(hl)] = changedBaseElement
		}
	}
	for _, changedBaseElement := range changedBaseElements {
		changedBaseElementId := changedBaseElement.GetId(hl).String()
		currentTreeNode := jquery.NewJQuery(treeManager.treeId).Call("jstree", "get_node", changedBaseElement.GetId(hl).String())
		if currentTreeNode.Length == 0 {
			// Node does not exist. Create it
			parentId := "#"
			parent := core.GetOwningElement(changedBaseElement, hl)
			if parent != nil {
				parentId = parent.GetId(hl).String()
			}
			treeManager.AddNode(changedBaseElement, parentId, hl)
		} else {
			// Node exists - update it
			// See if parent has changed
			currentTreeParentId := currentTreeNode.Attr("parent")
			currentParent := core.GetOwningElement(changedBaseElement, hl)
			currentParentId := "#" // the jstree version of a nil parent
			if currentParent != nil {
				currentParentId = currentParent.GetId(hl).String()
			}
			if currentTreeParentId != currentParentId {
				jquery.NewJQuery(treeManager.treeId).Call("jstree", "cut", changedBaseElementId)
				jquery.NewJQuery(treeManager.treeId).Call("jstree", "paste", currentParentId, "last")
			}

			// See if the name has changed
			changedBaseElementName := core.GetName(changedBaseElement, hl)
			if currentTreeNode.Attr("text") != changedBaseElementName {
				jquery.NewJQuery(treeManager.treeId).Call("jstree", "rename_node", changedBaseElementId, changedBaseElementName)
			}
		}
	}
}

func treeViewViewNode(instance core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {

}

func BuildTreeViews(conceptSpace core.Element, hl *core.HeldLocks) {
	uOfD := conceptSpace.GetUniverseOfDiscourse(hl)

	// TreeViews
	treeViews := uOfD.NewElement(hl, TreeViewsUri)
	core.SetName(treeViews, "TreeViews", hl)
	core.SetUri(treeViews, TreeViewsUri, hl)
	core.SetOwningElement(treeViews, conceptSpace, hl)

	// ManageNodes
	manageNodes := uOfD.NewElement(hl, ManageNodesUri)
	core.SetName(manageNodes, "ManageNodes", hl)
	core.SetUri(manageNodes, ManageNodesUri, hl)
	core.SetOwningElement(manageNodes, treeViews, hl)
	// ManageNodes UofD Reference
	uOfDReference := uOfD.NewElementReference(hl, ManageNodesUofDReferenceUri)
	core.SetName(uOfDReference, "UofDReference", hl)
	core.SetUri(uOfDReference, ManageNodesUofDReferenceUri, hl)
	core.SetOwningElement(uOfDReference, manageNodes, hl)

	// ViewNode
	viewNode := uOfD.NewElement(hl, ViewNodeUri)
	core.SetName(viewNode, "ViewNode", hl)
	core.SetUri(viewNode, ViewNodeUri, hl)
	core.SetOwningElement(viewNode, treeViews, hl)
	// ViewNode BaseElementReference
	baseElementReference := uOfD.NewBaseElementReference(hl, ViewNodeBaseElementReferenceUri)
	core.SetName(baseElementReference, "BaseElementReference", hl)
	core.SetUri(baseElementReference, ViewNodeBaseElementReferenceUri, hl)
	core.SetOwningElement(baseElementReference, viewNode, hl)
}

func registerTreeViewFunctions() {
	core.GetCore().AddFunction(ManageNodesUri, treeViewManageNodes)
	core.GetCore().AddFunction(ViewNodeUri, treeViewViewNode)

}
