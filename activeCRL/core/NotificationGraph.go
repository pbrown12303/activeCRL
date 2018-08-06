package core

import (
	"github.com/awalterschulze/gographviz"
	"reflect"
	"strconv"
)

type notificationGraph struct {
	graph            *gographviz.Graph
	callSequence     int
	callAnnotation   map[string]string
	nodeBaseElements map[string]BaseElement
}

func NewNotificationGraph(notification *ChangeNotification, hl *HeldLocks) *notificationGraph {
	var nGraph notificationGraph
	nGraph.callAnnotation = make(map[string]string)
	nGraph.nodeBaseElements = make(map[string]BaseElement)
	nGraph.graph = gographviz.NewGraph()
	nGraph.graph.SetDir(true)
	nGraph.graph.SetStrict(true)
	nGraph.addNotification(notification, hl)
	for _, node := range nGraph.graph.Nodes.Nodes {
		node.Attrs["label"] = nGraph.makeLabel(nGraph.nodeBaseElements[node.Name], hl)
	}
	return &nGraph
}

func (ngPtr *notificationGraph) addNotification(notification *ChangeNotification, hl *HeldLocks) {
	changedObject := notification.changedObject
	changedObjectId := makeGraphId(changedObject, hl)
	ngPtr.nodeBaseElements[changedObjectId] = changedObject
	ngPtr.makeNode(changedObject, hl)

	ngPtr.callAnnotation[changedObjectId] = ngPtr.callAnnotation[changedObjectId] + "<TR><TD>" + strconv.Itoa(ngPtr.callSequence) + ":" + notification.origin + "</TD></TR>"
	ngPtr.callSequence--

	ngPtr.graphParentsRecursively(changedObject, hl)

	switch changedObject.(type) {
	case ElementPointer:
		indicatedElement := changedObject.(ElementPointer).GetElement(hl)
		if indicatedElement != nil {
			indicatedElementId := makeGraphId(indicatedElement, hl)
			ngPtr.nodeBaseElements[indicatedElementId] = indicatedElement
			ngPtr.makeNode(indicatedElement, hl)
			ngPtr.makeIndicatedElementEdge(changedObjectId, indicatedElementId, hl)
			ngPtr.graphParentsRecursively(indicatedElement, hl)
		}
	}

	if notification.underlyingChange != nil {
		ngPtr.addNotification(notification.underlyingChange, hl)
	}
}

func (ngPtr *notificationGraph) getGraph() *gographviz.Graph {
	return ngPtr.graph
}

func (ngPtr *notificationGraph) graphParentsRecursively(child BaseElement, hl *HeldLocks) {
	parent := GetOwningElement(child, hl)
	if parent != nil {
		childObjectId := makeGraphId(child, hl)
		parentGraphId := makeGraphId(parent, hl)
		ngPtr.nodeBaseElements[parentGraphId] = parent
		ngPtr.makeNode(parent, hl)
		ngPtr.makeOwnerEdge(parentGraphId, childObjectId, hl)
		ngPtr.graphParentsRecursively(parent, hl)
	}

}

func (ngPtr *notificationGraph) makeNode(be BaseElement, hl *HeldLocks) {
	id := makeGraphId(be, hl)
	if ngPtr.graph.IsNode(id) != true {
		nodeAttrs := make(map[string]string)
		nodeAttrs["shape"] = "none"
		ngPtr.graph.AddNode("", id, nodeAttrs)
	}
}

func makeGraphId(be BaseElement, hl *HeldLocks) string {
	var graphId string = `"` + be.GetId(hl) + `"`
	return graphId
}

func (ngPtr *notificationGraph) makeIndicatedElementEdge(parentId string, childId string, hl *HeldLocks) {
	ownerEdgeAttrs := make(map[string]string)
	ngPtr.graph.AddEdge(parentId, childId, true, ownerEdgeAttrs)
}

func (ngPtr *notificationGraph) makeLabel(be BaseElement, hl *HeldLocks) string {
	typeString := reflect.TypeOf(be).String()
	graphId := makeGraphId(be, hl)
	augmentation := ""
	switch be.(type) {
	case ElementPointer:
		augmentation = ":" + be.(ElementPointer).GetElementPointerRole(hl).RoleToString()
	}
	label := "<<TABLE><TR><TD>" + typeString + augmentation + "</TD></TR><TR><TD>" + GetLabel(be, hl) + "</TD></TR><TR><TD>" + graphId + "</TD></TR>" + ngPtr.callAnnotation[graphId] + "</TABLE>>"
	return label
}

func (ngPtr *notificationGraph) makeOwnerEdge(parentId string, childId string, hl *HeldLocks) {
	ownerEdgeAttrs := make(map[string]string)
	ownerEdgeAttrs["arrowhead"] = "none"
	ownerEdgeAttrs["arrowtail"] = "diamond"
	ownerEdgeAttrs["dir"] = "both"
	ngPtr.graph.AddEdge(parentId, childId, true, ownerEdgeAttrs)
}
