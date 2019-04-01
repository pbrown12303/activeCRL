package core

import (
	"log"
	"reflect"
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
func (bgPtr *baseGraph) addNotification(notification *ChangeNotification, parentGraph string) string {
	changedObject := notification.GetReportingElement()
	changedObjectID := bgPtr.makeNode(changedObject, parentGraph, false)
	// By definition, the root notification's changed object is the root node
	if bgPtr.rootNodeIDs[parentGraph] == "" {
		bgPtr.rootNodeIDs[parentGraph] = changedObjectID
	}
	bgPtr.nodeElements[changedObjectID] = changedObject

	bgPtr.callAnnotation[changedObjectID] = bgPtr.callAnnotation[changedObjectID] + "<TR><TD>" + strconv.Itoa(bgPtr.parentGraphCallSequence[parentGraph]) + ":" + notification.reportingElement.getConceptIDNoLock() + notification.natureOfChange.String() + "</TD></TR>"
	bgPtr.parentGraphCallSequence[parentGraph]--

	bgPtr.graphParentsRecursively(changedObject, parentGraph)

	switch changedObject.(type) {
	case Reference:
		indicatedElement := changedObject.(Reference).getReferencedConceptNoLock()
		if indicatedElement != nil {
			indicatedElementID := makeGraphID(indicatedElement, bgPtr.parentGraphNodePrefix[parentGraph])
			bgPtr.nodeElements[indicatedElementID] = indicatedElement
			bgPtr.makeNode(indicatedElement, parentGraph, false)
			bgPtr.makeIndicatedElementEdge(changedObjectID, indicatedElementID)
			bgPtr.graphParentsRecursively(indicatedElement, parentGraph)
		}
	case Refinement:
		abstractConcept := changedObject.(Refinement).getAbstractConceptNoLock()
		if abstractConcept != nil {
			abstractConceptID := makeGraphID(abstractConcept, bgPtr.parentGraphNodePrefix[parentGraph])
			bgPtr.nodeElements[abstractConceptID] = abstractConcept
			bgPtr.makeNode(abstractConcept, parentGraph, false)
			bgPtr.makeAbstractConceptEdge(changedObjectID, abstractConceptID)
			bgPtr.graphParentsRecursively(abstractConcept, parentGraph)
		}
		refinedConcept := changedObject.(Refinement).getAbstractConceptNoLock()
		if refinedConcept != nil {
			refinedConceptID := makeGraphID(refinedConcept, bgPtr.parentGraphNodePrefix[parentGraph])
			bgPtr.nodeElements[refinedConceptID] = refinedConcept
			bgPtr.makeNode(refinedConcept, parentGraph, false)
			bgPtr.makeRefinedConceptEdge(changedObjectID, refinedConceptID)
			bgPtr.graphParentsRecursively(refinedConcept, parentGraph)
		}
	}

	if notification.underlyingChange != nil {
		underlyingNotificationID := bgPtr.addNotification(notification.underlyingChange, parentGraph)
		bgPtr.makeNotificationEdge(underlyingNotificationID, changedObjectID)
	}
	return changedObjectID
}

func (bgPtr *baseGraph) getRootNodeID(parentGraph string) string {
	return bgPtr.rootNodeIDs[parentGraph]
}

func (bgPtr *baseGraph) graphParentsRecursively(child Element, parentGraph string) {
	parent := child.getOwningConceptPointer().getIndicatedConcept()
	if parent != nil {
		childObjectID := makeGraphID(child, bgPtr.parentGraphNodePrefix[parentGraph])
		parentGraphID := makeGraphID(parent, bgPtr.parentGraphNodePrefix[parentGraph])
		bgPtr.nodeElements[parentGraphID] = parent
		bgPtr.makeNode(parent, parentGraph, false)
		bgPtr.makeOwnerEdge(parentGraphID, childObjectID)
		bgPtr.graphParentsRecursively(parent, parentGraph)
	}
}

// GetGraph returns the grqaphviz.Graph
func (bgPtr *baseGraph) GetGraph() *gographviz.Graph {
	return bgPtr.graph
}

func makeGraphID(be Element, prefix string) string {
	var graphID = prefix + "\"" + be.getConceptIDNoLock() + "\""
	return graphID
}

func (bgPtr *baseGraph) makeAbstractConceptEdge(sourceID string, targetID string) {
	abstractEdgeAttrs := make(map[string]string)
	abstractEdgeAttrs["arrowhead"] = "invempty"
	abstractEdgeAttrs["arrowtail"] = "none"
	abstractEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(sourceID, targetID, true, abstractEdgeAttrs)
	if err != nil {
		log.Printf("Error in BaseGraph.makeAbstractConceptEdge")
		log.Printf(err.Error())
	}
}
func (bgPtr *baseGraph) makeRefinedConceptEdge(sourceID string, targetID string) {
	refinedEdgeAttrs := make(map[string]string)
	refinedEdgeAttrs["arrowhead"] = "none"
	refinedEdgeAttrs["arrowtail"] = "inv"
	refinedEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(sourceID, targetID, true, refinedEdgeAttrs)
	if err != nil {
		log.Printf("Error in BaseGraph.makeRefinedConceptEdge")
		log.Printf(err.Error())
	}
}
func (bgPtr *baseGraph) makeIndicatedElementEdge(sourceID string, targetID string) {
	referenceEdgeAttrs := make(map[string]string)
	referenceEdgeAttrs["arrowhead"] = "normal"
	referenceEdgeAttrs["arrowtail"] = "none"
	referenceEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(sourceID, targetID, true, referenceEdgeAttrs)
	if err != nil {
		log.Printf("Error in BaseGraph.makeIndicatedElementEdge")
		log.Printf(err.Error())
	}
}

func (bgPtr *baseGraph) makeNotificationEdge(sourceID string, targetID string) {
	notificationEdgeAttrs := make(map[string]string)
	notificationEdgeAttrs["arrowhead"] = "open"
	notificationEdgeAttrs["arrowtail"] = "none"
	notificationEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(sourceID, targetID, true, notificationEdgeAttrs)
	if err != nil {
		log.Printf("Error in BaseGraph.makeNotificationEdge")
		log.Printf(err.Error())
	}
}

func (bgPtr *baseGraph) makeNode(el Element, parentGraph string, root bool) string {
	id := makeGraphID(el, bgPtr.parentGraphNodePrefix[parentGraph])
	if bgPtr.graph.IsNode(id) != true {
		nodeAttrs := make(map[string]string)
		if root {
			nodeAttrs["shape"] = "octagon"
		} else {
			nodeAttrs["shape"] = "none"
		}
		typeString := reflect.TypeOf(el).String()
		nodeAttrs["label"] = "<<TABLE><TR><TD>" + typeString + "</TD></TR><TR><TD>" + el.getLabelNoLock() + "</TD></TR><TR><TD>" + id + "</TD></TR></TABLE>>"
		err := bgPtr.graph.AddNode(parentGraph, id, nodeAttrs)
		if err != nil {
			log.Printf("Error in BaseGraph.makeNode")
			log.Printf(err.Error())
		}
		bgPtr.nodeToGraphName[id] = parentGraph
	}
	return id
}

func (bgPtr *baseGraph) makeLabel(graphID string, parentGraph string) string {
	el := bgPtr.nodeElements[graphID]
	if el == nil {
		log.Printf("In BaseGraph.makeLabel with nil Element for graphID %s\n", graphID)
		return ""
	}
	typeString := reflect.TypeOf(el).String()
	label := "<<TABLE><TR><TD>" + typeString + "</TD></TR><TR><TD>" + el.getLabelNoLock() + "</TD></TR><TR><TD>" + graphID + "</TD></TR>" + bgPtr.callAnnotation[graphID] + "</TABLE>>"
	return label
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
