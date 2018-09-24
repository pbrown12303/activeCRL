package core

import (
	"log"
	"strconv"
)

// FunctionCallGraph is a graphical representation of the function calls made when HeldLocks.ReleaseLocks is called
type FunctionCallGraph struct {
	baseGraph
	functionName     crlExecutionFunctionArrayIdentifier
	executingElement Element
}

// NewFunctionCallGraph creates a function call graph for the indicated function, executing element, and notifications
func NewFunctionCallGraph(functionID crlExecutionFunctionArrayIdentifier, executingElement Element, notifications []*ChangeNotification) *FunctionCallGraph {
	var fcGraph FunctionCallGraph
	graphName := "FunctionCallGraph"
	fcGraph.initializeBaseGraph(graphName)
	fcGraph.parentGraphNodePrefix[graphName] = ""
	fcGraph.functionName = functionID
	fcGraph.executingElement = executingElement
	executingElementNodeID := fcGraph.makeNode(executingElement, graphName)
	fcGraph.nodeBaseElements[executingElementNodeID] = executingElement
	subgraphCount := 0
	for _, notification := range notifications {
		subgraphName := "cluster_" + strconv.Itoa(subgraphCount)
		subgraphPrefix := "NG_" + strconv.Itoa(subgraphCount) + "_"
		fcGraph.parentGraphNodePrefix[subgraphName] = subgraphPrefix
		err := fcGraph.graph.AddSubGraph(graphName, subgraphName, map[string]string{"label": subgraphPrefix})
		if err != nil {
			log.Printf("Error in FunctionCallGraph.NewFunctionCallGraph adding a subgraph")
			log.Printf(err.Error())
		}
		notification.GetDepth()
		fcGraph.addNotification(notification, subgraphName)
		err = fcGraph.graph.AddEdge(executingElementNodeID, fcGraph.rootNodeIDs[subgraphName], true, map[string]string{})
		if err != nil {
			log.Printf("Error in FunctionCallGraph.NewFunctionCallGraph adding an edge from executing element to the subgraph root node")
			log.Printf(err.Error())
		}
		subgraphCount++
	}
	for _, node := range fcGraph.graph.Nodes.Nodes {
		node.Attrs["label"] = fcGraph.makeLabel(node.Name, fcGraph.nodeToGraphName[node.Name])
	}
	return &fcGraph
}
