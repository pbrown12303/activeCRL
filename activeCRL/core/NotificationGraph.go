package core

// NotificationGraph is a graphviz graphical representation of the structure of a ChangeNotification
type NotificationGraph struct {
	baseGraph
}

// NewNotificationGraph creates a graphviz graph of the information in the ChangeNotification
// The prefix parameter provides the opportunity for the caller to insert a string prefix into the
// node identifier.
func NewNotificationGraph(notification *ChangeNotification) *NotificationGraph {
	var nGraph NotificationGraph
	graphName := "NotificationGraph"
	nGraph.initializeBaseGraph(graphName)
	nGraph.addNotification(notification, graphName)
	for _, node := range nGraph.graph.Nodes.Nodes {
		node.Attrs["label"] = nGraph.makeLabel(node.Name, nGraph.nodeToGraphName[node.Name])
	}
	return &nGraph
}
