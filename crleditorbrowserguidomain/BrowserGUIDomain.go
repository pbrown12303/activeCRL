package crleditorbrowserguidomain

import (
	"github.com/pbrown12303/activeCRL/core"
)

// BrowserGUIDomainURI is the URI for accessing the BrowserGUI
var BrowserGUIDomainURI = "http://activeCrl.com/crleditorbrowserguidomain/CrlBrowserGUIDomain"

// TreeViewsURI identifies the TreeViews concept
var TreeViewsURI = BrowserGUIDomainURI + "/TreeViews"

// TreeNodeManagerURI identifies the ManageNodes concept
var TreeNodeManagerURI = TreeViewsURI + "/TreeNodeManager"

// TreeNodeManagerUofDReferenceURI identifies the ManageNodesUofDReference
var TreeNodeManagerUofDReferenceURI = TreeNodeManagerURI + "/UofDReference"

// ViewNodeURI identifies the ViewNode concept
var ViewNodeURI = TreeViewsURI + "/ViewNode"

// ViewNodeElementReferenceURI identifies the ViewNodeElementReference concept
var ViewNodeElementReferenceURI = ViewNodeURI + "/ElementReference"

// DiagramViewMonitorURI identifies the diagram view monitor
var DiagramViewMonitorURI = BrowserGUIDomainURI + "/DiagramViewMonitor"

// BuildBrowserGUIDomain programmatically constructs the BrowserGUIDomain
func BuildBrowserGUIDomain(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	// EditorDomain
	// Assumes that the DiagramDomain has already been added to the uOfD
	domain, _ := uOfD.NewElement(hl, BrowserGUIDomainURI)
	domain.SetLabel("CrlEditorBrowserGuiDomain", hl)

	// TreeViews
	treeViews, _ := uOfD.NewElement(hl, TreeViewsURI)
	treeViews.SetLabel("TreeViews", hl)
	treeViews.SetURI(TreeViewsURI, hl)
	treeViews.SetOwningConcept(domain, hl)

	// ManageNodes
	manageNodes, _ := uOfD.NewElement(hl, TreeNodeManagerURI)
	manageNodes.SetLabel("TreeNodeManager", hl)
	manageNodes.SetURI(TreeNodeManagerURI, hl)
	manageNodes.SetOwningConcept(treeViews, hl)
	// ManageNodes UofD Reference
	uOfDReference, _ := uOfD.NewReference(hl, TreeNodeManagerUofDReferenceURI)
	uOfDReference.SetLabel("UofDReference", hl)
	uOfDReference.SetURI(TreeNodeManagerUofDReferenceURI, hl)
	uOfDReference.SetOwningConcept(manageNodes, hl)

	// ViewNode
	viewNode, _ := uOfD.NewElement(hl, ViewNodeURI)
	viewNode.SetLabel("ViewNode", hl)
	viewNode.SetURI(ViewNodeURI, hl)
	viewNode.SetOwningConcept(treeViews, hl)
	// ViewNode BaseElementReference
	reference, _ := uOfD.NewReference(hl, ViewNodeElementReferenceURI)
	reference.SetLabel("ElementReference", hl)
	reference.SetURI(ViewNodeElementReferenceURI, hl)
	reference.SetOwningConcept(viewNode, hl)

	// DiagramViewMonitor
	diagramViewMonitor, _ := uOfD.NewReference(hl, DiagramViewMonitorURI)
	diagramViewMonitor.SetLabel("DiagramViewMonitor", hl)
	diagramViewMonitor.SetOwningConcept(domain, hl)

	domain.SetIsCoreRecursively(hl)

	return domain, nil
}
