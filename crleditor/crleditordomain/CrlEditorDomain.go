package crleditordomain

import (
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructures"
	"github.com/pbrown12303/activeCRL/crldiagram"
)

// EditorDomainURI is the URI for the domain of the editor
var EditorDomainURI = "http://activeCrl.com/crlEditor/EditorDomain"

// EditorURI is the URI for accessing the CrlEditor
var EditorURI = "http://activeCrl.com/crlEditor/Editor"

// EditorSettingsURI is the URI for accessing the CrlEditor settings
var EditorSettingsURI = EditorURI + "/Settings"

// EditorOpenDiagramsURI is the URI for accessing the CrlEditor list of open diagrams
var EditorOpenDiagramsURI = EditorSettingsURI + "/OpenDiagrams"

// editorDropReferenceAsLinkURI is the URI for the setting parameter DropReferenceAsLink
var editorDropReferenceAsLinkURI = EditorSettingsURI + "/DropReferenceAsLink"

// editorDropRefinmentAsLinkURI is the URI for the setting parameter DropRefinementAsLink
var editorDropRefinementAsLinkURI = EditorSettingsURI + "/DropRefinementAsLink"

// TreeViewsURI identifies the TreeViews concept
var TreeViewsURI = EditorDomainURI + "/TreeViews"

// TreeNodeManagerURI identifies the ManageNodes concept
var TreeNodeManagerURI = TreeViewsURI + "/TreeNodeManager"

// TreeNodeManagerUofDReferenceURI identifies the ManageNodesUofDReference
var TreeNodeManagerUofDReferenceURI = TreeNodeManagerURI + "/UofDReference"

// ViewNodeURI identifies the ViewNode concept
var ViewNodeURI = TreeViewsURI + "/ViewNode"

// ViewNodeElementReferenceURI identifies the ViewNodeElementReference concept
var ViewNodeElementReferenceURI = ViewNodeURI + "/ElementReference"

// DiagramViewMonitorURI identifies the diagram view monitor
var DiagramViewMonitorURI = EditorDomainURI + "/DiagramViewMonitor"

// BuildEditorDomain programmatically constructs the EditorDomain
func BuildEditorDomain(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	// EditorDomain
	// Assumes that the DiagramConceptSpace has already been added to the uOfD
	domain, _ := uOfD.NewElement(hl, EditorDomainURI)
	domain.SetLabel("EditorDomain", hl)

	settings, _ := uOfD.NewElement(hl, EditorSettingsURI)
	settings.SetLabel("Settings", hl)
	settings.SetOwningConcept(domain, hl)

	diagram := uOfD.GetElementWithURI(crldiagram.CrlDiagramURI)
	editorOpenDiagramsList, err := crldatastructures.NewList(uOfD, diagram, hl, EditorOpenDiagramsURI)
	if err != nil {
		return nil, err
	}
	//	editorOpenDiagramsList.SetURI(EditorOpenDiagramsURI, hl)
	editorOpenDiagramsList.SetLabel("EditorOpenDiagrams", hl)
	editorOpenDiagramsList.SetOwningConcept(settings, hl)

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
