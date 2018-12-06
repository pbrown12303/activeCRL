package core

import (
	"log"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

type baseGraph struct {
	graph          *gographviz.Graph
	callAnnotation map[string]string
	nodeElements   map[string]Element
	// rootNodeIDs maps a parent graph to the ID of its root node. Each graph or subgraph has a root node.
	rootNodeIDs map[string]string
	// parentGraphNodePrefix maps each parentGraph to the prefix used for its nodes
	parentGraphNodePrefix map[string]string
	// parentGraphCallSequence keeps track of the call sequence for each graph
	parentGraphCallSequence map[string]int
	// nodeToGraphName keeps track of the graph (or subgraph) name in whieh the node resides
	nodeToGraphName map[string]string
}

func (bgPtr *baseGraph) initializeBaseGraph(graphName string) {
	bgPtr.graph = gographviz.NewGraph()
	bgPtr.graph.SetDir(true)
	bgPtr.graph.SetStrict(true)
	bgPtr.graph.SetName(graphName)
	bgPtr.callAnnotation = make(map[string]string)
	bgPtr.nodeElements = make(map[string]Element)
	bgPtr.rootNodeIDs = make(map[string]string)
	bgPtr.parentGraphNodePrefix = make(map[string]string)
	bgPtr.parentGraphCallSequence = make(map[string]int)
	bgPtr.nodeToGraphName = make(map[string]string)
}

// addNotification adds a notification to a graph (its parent graph). The changed object is added to the graph if not already present.
// If this is the root notification, the ID of the changed object becomes the rootNodeID for the parentGraph.
// If the changed object does not exist as a node, a node is created and an annotation is added to indicate the type of notification
// and the position in the notification hierarcy. If the changed object already exists then a new annotation is just added
func (bgPtr *baseGraph) addNotification(notification *ChangeNotification, parentGraph string) {
	changedObject := notification.GetReportingElement()
	changedObjectID := bgPtr.makeNode(changedObject, parentGraph)
	// By definition, the root notification's changed object is the root node
	if bgPtr.rootNodeIDs[parentGraph] == "" {
		bgPtr.rootNodeIDs[parentGraph] = changedObjectID
	}
	bgPtr.nodeElements[changedObjectID] = changedObject

	bgPtr.callAnnotation[changedObjectID] = bgPtr.callAnnotation[changedObjectID] + "<TR><TD>" + strconv.Itoa(bgPtr.parentGraphCallSequence[parentGraph]) + ":" + notification.reportingElement.getConceptIDNoLock() + "</TD></TR>"
	bgPtr.parentGraphCallSequence[parentGraph]--

	bgPtr.graphParentsRecursively(changedObject, parentGraph)

	switch changedObject.(type) {
	case Reference:
		indicatedElement := changedObject.(Reference).getReferencedConceptNoLock()
		if indicatedElement != nil {
			indicatedElementID := makeGraphID(indicatedElement, bgPtr.parentGraphNodePrefix[parentGraph])
			bgPtr.nodeElements[indicatedElementID] = indicatedElement
			bgPtr.makeNode(indicatedElement, parentGraph)
			bgPtr.makeIndicatedElementEdge(changedObjectID, indicatedElementID)
			bgPtr.graphParentsRecursively(indicatedElement, parentGraph)
		}
	}

	if notification.underlyingChange != nil {
		bgPtr.addNotification(notification.underlyingChange, parentGraph)
	}
}

func (bgPtr *baseGraph) getRootNodeID(parentGraph string) string {
	return bgPtr.rootNodeIDs[parentGraph]
}

func (bgPtr *baseGraph) graphParentsRecursively(child Element, parentGraph string) {
	// TODO
	// parent := getOwningConceptNoLock(child)
	// if parent != nil {
	// 	childObjectID := makeGraphID(child, bgPtr.parentGraphNodePrefix[parentGraph])
	// 	parentGraphID := makeGraphID(parent, bgPtr.parentGraphNodePrefix[parentGraph])
	// 	bgPtr.nodeElements[parentGraphID] = parent
	// 	bgPtr.makeNode(parent, parentGraph)
	// 	bgPtr.makeOwnerEdge(parentGraphID, childObjectID)
	// 	bgPtr.graphParentsRecursively(parent, parentGraph)
	// }

}

// GetGraph returns the grqaphviz.Graph
func (bgPtr *baseGraph) GetGraph() *gographviz.Graph {
	return bgPtr.graph
}

func makeGraphID(be Element, prefix string) string {
	var graphID = "\"" + prefix + be.getConceptIDNoLock() + "\""
	return graphID
}

func (bgPtr *baseGraph) makeIndicatedElementEdge(parentID string, childID string) {
	ownerEdgeAttrs := make(map[string]string)
	err := bgPtr.graph.AddEdge(parentID, childID, true, ownerEdgeAttrs)
	if err != nil {
		log.Printf("Error in BaseGraph.makeIndicatedElementEdge")
		log.Printf(err.Error())
	}
}

func (bgPtr *baseGraph) makeNode(be Element, parentGraph string) string {
	id := makeGraphID(be, bgPtr.parentGraphNodePrefix[parentGraph])
	// TODO
	// if bgPtr.graph.IsNode(id) != true {
	// 	nodeAttrs := make(map[string]string)
	// 	nodeAttrs["shape"] = "none"
	// 	typeString := reflect.TypeOf(be).String()
	// 	nodeAttrs["label"] = "<<TABLE><TR><TD>" + typeString + "</TD></TR><TR><TD>" + getLabelNoLock(be) + "</TD></TR><TR><TD>" + id + "</TD></TR></TABLE>>"
	// 	err := bgPtr.graph.AddNode(parentGraph, id, nodeAttrs)
	// 	if err != nil {
	// 		log.Printf("Error in BaseGraph.makeNode")
	// 		log.Printf(err.Error())
	// 	}
	// 	bgPtr.nodeToGraphName[id] = parentGraph
	// }
	return id
}

func (bgPtr *baseGraph) makeLabel(graphID string, parentGraph string) string {
	// TODO
	// be := bgPtr.nodeElements[graphID]
	// if be == nil {
	// 	log.Printf("In BaseGraph.makeLabel with nil BaseElement for graphID %s\n", graphID)
	// 	return ""
	// }
	// typeString := reflect.TypeOf(be).String()
	// augmentation := ""
	// switch be.(type) {
	// case ElementPointer:
	// 	augmentation = ":" + be.(ElementPointer).getElementPointerRoleNoLock().RoleToString()
	// }
	// label := "<<TABLE><TR><TD>" + typeString + augmentation + "</TD></TR><TR><TD>" + getLabelNoLock(be) + "</TD></TR><TR><TD>" + graphID + "</TD></TR>" + bgPtr.callAnnotation[graphID] + "</TABLE>>"
	// return label
	return ""
}

func (bgPtr *baseGraph) makeOwnerEdge(parentID string, childID string) {
	ownerEdgeAttrs := make(map[string]string)
	ownerEdgeAttrs["arrowhead"] = "none"
	ownerEdgeAttrs["arrowtail"] = "diamond"
	ownerEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(parentID, childID, true, ownerEdgeAttrs)
	if err != nil {
		log.Printf("Error in BaseGraph.makeOwnerEdge")
		log.Printf(err.Error())
	}
}
