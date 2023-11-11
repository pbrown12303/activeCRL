package core

import (
	"log"
	"reflect"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

type baseGraph struct {
	graph             *gographviz.Graph
	callAnnotation    map[string]string
	nodeElementLabels map[string]string
	nodeElementTypes  map[string]string
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
	bgPtr.nodeElementLabels = make(map[string]string)
	bgPtr.nodeElementTypes = make(map[string]string)
	bgPtr.rootNodeIDs = make(map[string]string)
	bgPtr.parentGraphNodePrefix = make(map[string]string)
	bgPtr.parentGraphCallSequence = make(map[string]int)
	bgPtr.nodeToGraphName = make(map[string]string)
}

// addNotification adds a notification to a graph (its parent graph). The reporting element is added to the graph if not already present.
// If this is the root notification, the ID of the reporting element becomes the rootNodeID for the parentGraph.
// If the reporting element does not exist as a node, a node is created and an annotation is added to indicate the type of notification
// and the position in the notification hierarcy. If the reporting element already exists then a new annotation is just added
func (bgPtr *baseGraph) addNotification(notification *ChangeNotification, parentGraph string) string {
	reportingElement := notification.uOfD.GetElement(notification.GetReportingElementID()) // this will return nil after an element deletion
	reportingElementNodeID := bgPtr.makeNode(notification.GetReportingElementID(), notification.GetReportingElementType(), notification.GetReportingElementLabel(), parentGraph, false, "")
	// By definition, the root notification's changed object is the root node
	if bgPtr.rootNodeIDs[parentGraph] == "" {
		bgPtr.rootNodeIDs[parentGraph] = reportingElementNodeID
	}
	bgPtr.nodeElementLabels[reportingElementNodeID] = notification.GetReportingElementLabel()
	bgPtr.nodeElementTypes[reportingElementNodeID] = notification.GetReportingElementType()

	bgPtr.callAnnotation[reportingElementNodeID] = bgPtr.callAnnotation[reportingElementNodeID] + "<TR><TD>" + strconv.Itoa(bgPtr.parentGraphCallSequence[parentGraph]) + ":" + notification.GetChangedConceptID() + notification.natureOfChange.String() + "</TD></TR>"
	bgPtr.parentGraphCallSequence[parentGraph]--

	bgPtr.graphParentsRecursively(reportingElement, parentGraph)

	switch reportingElement.GetConceptType() {
	case Reference:
		indicatedElement := reportingElement.getReferencedConceptNoLock()
		if indicatedElement != nil {
			indicatedElementID := makeGraphID(indicatedElement.getConceptIDNoLock(), bgPtr.parentGraphNodePrefix[parentGraph])
			bgPtr.nodeElementLabels[indicatedElementID] = indicatedElement.getLabelNoLock()
			bgPtr.nodeElementTypes[indicatedElementID] = reflect.TypeOf(indicatedElement).String()
			bgPtr.makeNode(indicatedElement.getConceptIDNoLock(), GetConceptTypeString(indicatedElement), notification.GetChangedConceptLabel(), parentGraph, false, "")
			bgPtr.makeIndicatedElementEdge(reportingElementNodeID, indicatedElementID)
			bgPtr.graphParentsRecursively(indicatedElement, parentGraph)
		}
	case Refinement:
		abstractConcept := reportingElement.getAbstractConceptNoLock()
		if abstractConcept != nil {
			abstractConceptID := makeGraphID(abstractConcept.getConceptIDNoLock(), bgPtr.parentGraphNodePrefix[parentGraph])
			bgPtr.nodeElementLabels[abstractConceptID] = abstractConcept.getLabelNoLock()
			bgPtr.nodeElementTypes[abstractConceptID] = reflect.TypeOf(abstractConcept).String()
			bgPtr.makeNode(abstractConcept.getConceptIDNoLock(), GetConceptTypeString(abstractConcept), notification.GetChangedConceptLabel(), parentGraph, false, "")
			bgPtr.makeAbstractConceptEdge(reportingElementNodeID, abstractConceptID)
			bgPtr.graphParentsRecursively(abstractConcept, parentGraph)
		}
		refinedConcept := reportingElement.getAbstractConceptNoLock()
		if refinedConcept != nil {
			refinedConceptID := makeGraphID(refinedConcept.getConceptIDNoLock(), bgPtr.parentGraphNodePrefix[parentGraph])
			bgPtr.nodeElementLabels[refinedConceptID] = refinedConcept.getLabelNoLock()
			bgPtr.nodeElementTypes[refinedConceptID] = reflect.TypeOf(refinedConcept).String()
			bgPtr.makeNode(refinedConcept.getConceptIDNoLock(), GetConceptTypeString(refinedConcept), notification.GetChangedConceptLabel(), parentGraph, false, "")
			bgPtr.makeRefinedConceptEdge(reportingElementNodeID, refinedConceptID)
			bgPtr.graphParentsRecursively(refinedConcept, parentGraph)
		}
	}

	if notification.underlyingChange != nil {
		underlyingNotificationID := bgPtr.addNotification(notification.underlyingChange, parentGraph)
		bgPtr.makeNotificationEdge(underlyingNotificationID, reportingElementNodeID)
	}
	return reportingElementNodeID
}

// func (bgPtr *baseGraph) getRootNodeID(parentGraph string) string {
// 	return bgPtr.rootNodeIDs[parentGraph]
// }

func (bgPtr *baseGraph) graphParentsRecursively(child *Concept, parentGraph string) {
	if child == nil {
		return
	}
	parent := child.getOwningConceptNoLock()
	if parent != nil {
		childObjectID := makeGraphID(child.getConceptIDNoLock(), bgPtr.parentGraphNodePrefix[parentGraph])
		parentGraphID := makeGraphID(parent.getConceptIDNoLock(), bgPtr.parentGraphNodePrefix[parentGraph])
		bgPtr.nodeElementLabels[parentGraphID] = parent.getLabelNoLock()
		bgPtr.nodeElementTypes[parentGraphID] = reflect.TypeOf(parent).String()
		bgPtr.makeNode(parent.getConceptIDNoLock(), GetConceptTypeString(parent), parent.getLabelNoLock(), parentGraph, false, "")
		bgPtr.makeOwnerEdge(parentGraphID, childObjectID)
		bgPtr.graphParentsRecursively(parent, parentGraph)
	}
}

// GetGraph returns the grqaphviz.Graph
func (bgPtr *baseGraph) GetGraph() *gographviz.Graph {
	return bgPtr.graph
}

func makeGraphID(conceptID string, prefix string) string {
	var graphID = prefix + "\"" + conceptID + "\""
	return graphID
}

func (bgPtr *baseGraph) makeAbstractConceptEdge(sourceID string, targetID string) {
	abstractEdgeAttrs := make(map[string]string)
	abstractEdgeAttrs["arrowhead"] = "invempty"
	abstractEdgeAttrs["arrowtail"] = "none"
	abstractEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(sourceID, targetID, true, abstractEdgeAttrs)
	if err != nil {
		log.Print("Error in BaseGraph.makeAbstractConceptEdge")
		log.Print(err.Error())
	}
}
func (bgPtr *baseGraph) makeRefinedConceptEdge(sourceID string, targetID string) {
	refinedEdgeAttrs := make(map[string]string)
	refinedEdgeAttrs["arrowhead"] = "none"
	refinedEdgeAttrs["arrowtail"] = "inv"
	refinedEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(sourceID, targetID, true, refinedEdgeAttrs)
	if err != nil {
		log.Print("Error in BaseGraph.makeRefinedConceptEdge")
		log.Print(err.Error())
	}
}
func (bgPtr *baseGraph) makeIndicatedElementEdge(sourceID string, targetID string) {
	referenceEdgeAttrs := make(map[string]string)
	referenceEdgeAttrs["arrowhead"] = "normal"
	referenceEdgeAttrs["arrowtail"] = "none"
	referenceEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(sourceID, targetID, true, referenceEdgeAttrs)
	if err != nil {
		log.Print("Error in BaseGraph.makeIndicatedElementEdge")
		log.Print(err.Error())
	}
}

func (bgPtr *baseGraph) makeNotificationEdge(sourceID string, targetID string) {
	notificationEdgeAttrs := make(map[string]string)
	notificationEdgeAttrs["arrowhead"] = "open"
	notificationEdgeAttrs["arrowtail"] = "none"
	notificationEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(sourceID, targetID, true, notificationEdgeAttrs)
	if err != nil {
		log.Print("Error in BaseGraph.makeNotificationEdge")
		log.Print(err.Error())
	}
}

func (bgPtr *baseGraph) makeNode(conceptID string, typeString string, label string, parentGraph string, root bool, functionName string) string {
	if conceptID == "" || typeString == "" {
		log.Print("baseGraph.makeNode called will empty strings")
	}
	id := makeGraphID(conceptID, bgPtr.parentGraphNodePrefix[parentGraph])
	if !bgPtr.graph.IsNode(id) {
		nodeAttrs := make(map[string]string)
		if root {
			nodeAttrs["shape"] = "none"
			nodeAttrs["fillcolor"] = "yellow"
		} else {
			nodeAttrs["shape"] = "none"
		}
		// typeString := reflect.TypeOf(el).String()
		if root {
			// nodeAttrs["label"] = "<<TABLE HEIGHT='0' WIDTH='0'><TR><TD>" + functionName + "</TD></TR><TR><TD>" + typeString + "</TD></TR><TR><TD>" + label + "</TD></TR><TR><TD>" + id + "</TD></TR></TABLE>>"
			nodeAttrs["label"] = "<<TABLE><TR><TD> functionName </TD></TR><TR><TD>" + typeString + "</TD></TR><TR><TD>" + id + "</TD></TR></TABLE>>"
		} else {
			nodeAttrs["label"] = "<<TABLE><TR><TD>" + typeString + "</TD></TR><TR><TD>" + label + "</TD></TR><TR><TD>" + id + "</TD></TR></TABLE>>"
		}
		err := bgPtr.graph.AddNode(parentGraph, id, nodeAttrs)
		if err != nil {
			log.Print("Error in BaseGraph.makeNode")
			log.Print(err.Error())
		}
		bgPtr.nodeToGraphName[id] = parentGraph
	}
	return id
}

func (bgPtr *baseGraph) makeLabel(graphID string, parentGraph string, functionID string) string {
	elLabel := bgPtr.nodeElementLabels[graphID]
	elType := bgPtr.nodeElementTypes[graphID]
	var label string
	if functionID == "" {
		label = "<<TABLE><TR><TD>" + elType + "</TD></TR><TR><TD>" + elLabel + "</TD></TR><TR><TD>" + graphID + "</TD></TR>" + bgPtr.callAnnotation[graphID] + "</TABLE>>"
	} else {
		label = "<<TABLE><TR><TD BGCOLOR='yellow'>" + functionID + "</TD></TR><TR><TD>" + elType + "</TD></TR><TR><TD>" + elLabel + "</TD></TR><TR><TD>" + graphID + "</TD></TR>" + bgPtr.callAnnotation[graphID] + "</TABLE>>"
	}
	return label
}

func (bgPtr *baseGraph) makeOwnerEdge(parentID string, childID string) {
	ownerEdgeAttrs := make(map[string]string)
	ownerEdgeAttrs["arrowhead"] = "none"
	ownerEdgeAttrs["arrowtail"] = "diamond"
	ownerEdgeAttrs["dir"] = "both"
	err := bgPtr.graph.AddEdge(parentID, childID, true, ownerEdgeAttrs)
	if err != nil {
		log.Print("Error in BaseGraph.makeOwnerEdge")
		log.Print(err.Error())
	}
}
