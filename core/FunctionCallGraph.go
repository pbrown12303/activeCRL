package core

import (
	"log"
	"reflect"
)

// FunctionCallGraph is a graphical representation of the function calls made when HeldLocks.ReleaseLocks is called
type FunctionCallGraph struct {
	baseGraph
	functionName     string
	executingElement Element
}

// NewFunctionCallGraph creates a function call graph for the indicated function, executing element, and notifications
func NewFunctionCallGraph(functionID string, executingElement Element, notification *ChangeNotification, hl *HeldLocks) *FunctionCallGraph {
	var fcGraph FunctionCallGraph
	graphName := "FunctionCallGraph"
	fcGraph.initializeBaseGraph(graphName)
	fcGraph.parentGraphNodePrefix[graphName] = ""
	fcGraph.functionName = functionID
	fcGraph.executingElement = executingElement
	executingElementID := executingElement.getConceptIDNoLock()
	executingElementType := reflect.TypeOf(executingElement).String()
	executingElementLabel := executingElement.getLabelNoLock()
	executingElementNodeID := fcGraph.makeNode(executingElementID, GetConceptTypeString(executingElement), executingElement.getLabelNoLock(), graphName, true, functionID)
	fcGraph.nodeElementLabels[executingElementNodeID] = executingElementLabel
	fcGraph.nodeElementLabels[executingElementNodeID] = executingElementType
	fcGraph.graphParentsRecursively(executingElement, "")
	fcGraph.addNotification(notification, graphName)
	err := fcGraph.graph.AddEdge(executingElementNodeID, fcGraph.rootNodeIDs[graphName], true, map[string]string{})
	if err != nil {
		log.Printf("Error in FunctionCallGraph.NewFunctionCallGraph adding an edge from executing element to the subgraph root node")
		log.Printf(err.Error())
	}
	for _, node := range fcGraph.graph.Nodes.Nodes {
		if node.Name == executingElementNodeID {
			node.Attrs["label"] = fcGraph.makeLabel(node.Name, fcGraph.nodeToGraphName[node.Name], functionID)
		} else {
			node.Attrs["label"] = fcGraph.makeLabel(node.Name, fcGraph.nodeToGraphName[node.Name], "")
		}
	}
	return &fcGraph
}
