package core

import (
	"log"
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
	executingElementNodeID := fcGraph.makeNode(executingElement, graphName, true)
	fcGraph.nodeElements[executingElementNodeID] = executingElement
	fcGraph.graphParentsRecursively(executingElement, "")
	fcGraph.addNotification(notification, graphName)
	err := fcGraph.graph.AddEdge(executingElementNodeID, fcGraph.rootNodeIDs[graphName], true, map[string]string{})
	if err != nil {
		log.Printf("Error in FunctionCallGraph.NewFunctionCallGraph adding an edge from executing element to the subgraph root node")
		log.Printf(err.Error())
	}
	for _, node := range fcGraph.graph.Nodes.Nodes {
		node.Attrs["label"] = fcGraph.makeLabel(node.Name, fcGraph.nodeToGraphName[node.Name])
	}
	return &fcGraph
}
