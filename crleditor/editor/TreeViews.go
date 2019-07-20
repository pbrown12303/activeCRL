package editor

import (
	"github.com/pbrown12303/activeCRL/core"

	//	"github.com/satori/go.uuid"
	"log"
)

// TreeViewsURI identifies the TreeViews concept
var TreeViewsURI = editorURI + "/TreeViews"

// ManageTreeNodesURI identifies the ManageNodes concept
var ManageTreeNodesURI = TreeViewsURI + "/ManageTreeNodes"

// ManageNodesUofDReferenceURI identifies the ManageNodesUofDReference
var ManageNodesUofDReferenceURI = ManageTreeNodesURI + "/UofDReference"

// ViewNodeURI identifies the ViewNode concept
var ViewNodeURI = TreeViewsURI + "/ViewNode"

// ViewNodeElementReferenceURI identifies the ViewNodeElementReference concept
var ViewNodeElementReferenceURI = ViewNodeURI + "/ElementReference"

// treeViewManageNodes() is the callback function that manaages the tree view when base elements in the Universe of Discourse change.
// The changes being sought are the addition, removal, and re-parenting of base elements and the changes in their names.
func treeViewManageNodes(instance core.Element, changeNotification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocks()

	// Tracing
	if core.AdHocTrace == true {
		log.Printf("In treeViewManageNodes()")
	}

	treeManager := CrlEditorSingleton.getTreeManager()

	switch changeNotification.GetNatureOfChange() {
	case core.IndicatedConceptChanged:
		underlyingChange := changeNotification.GetUnderlyingChange()
		if underlyingChange == nil {
			log.Printf("treeViewManageNodes called with IndicatedConceptChanged but no underlying chanage")
			return
		}
		switch underlyingChange.GetNatureOfChange() {
		case core.IndicatedConceptChanged:
			secondUnderlyingChange := underlyingChange.GetUnderlyingChange()
			if secondUnderlyingChange == nil {
				log.Printf("treeViewManageNodes called with IndicatedConceptChanged but no underlying chanage")
				return
			}
			switch secondUnderlyingChange.GetNatureOfChange() {
			case core.UofDConceptAdded:
				changedElement := secondUnderlyingChange.GetPriorState()
				treeManager.addNode(changedElement, hl)
			case core.UofDConceptChanged:
				thirdUnderlyingChange := secondUnderlyingChange.GetUnderlyingChange()
				if thirdUnderlyingChange == nil {
					log.Printf("treeViewManageNodes called with UofDConceptChanged but no thirdUnderlyingChange chanage")
					return
				}
				changedElement := thirdUnderlyingChange.GetReportingElement()
				treeManager.changeNode(changedElement, hl)
			case core.UofDConceptRemoved:
				changedElement := secondUnderlyingChange.GetPriorState()
				treeManager.removeNode(changedElement, hl)
			}
		}
	}

}

// BuildTreeViews builds the concepts related to TreeViews and adds them to the uOfD
func BuildTreeViews(conceptSpace core.Element, hl *core.HeldLocks) {
	uOfD := conceptSpace.GetUniverseOfDiscourse(hl)

	// TreeViews
	treeViews, _ := uOfD.NewElement(hl, TreeViewsURI)
	treeViews.SetLabel("TreeViews", hl)
	treeViews.SetURI(TreeViewsURI, hl)
	treeViews.SetOwningConcept(conceptSpace, hl)
	treeViews.SetIsCore(hl)

	// ManageNodes
	manageNodes, _ := uOfD.NewElement(hl, ManageTreeNodesURI)
	manageNodes.SetLabel("ManageTreeNodes", hl)
	manageNodes.SetURI(ManageTreeNodesURI, hl)
	manageNodes.SetOwningConcept(treeViews, hl)
	manageNodes.SetIsCore(hl)
	// ManageNodes UofD Reference
	uOfDReference, _ := uOfD.NewReference(hl, ManageNodesUofDReferenceURI)
	uOfDReference.SetLabel("UofDReference", hl)
	uOfDReference.SetURI(ManageNodesUofDReferenceURI, hl)
	uOfDReference.SetOwningConcept(manageNodes, hl)
	uOfDReference.SetIsCore(hl)

	// ViewNode
	viewNode, _ := uOfD.NewElement(hl, ViewNodeURI)
	viewNode.SetLabel("ViewNode", hl)
	viewNode.SetURI(ViewNodeURI, hl)
	viewNode.SetOwningConcept(treeViews, hl)
	viewNode.SetIsCore(hl)
	// ViewNode BaseElementReference
	reference, _ := uOfD.NewReference(hl, ViewNodeElementReferenceURI)
	reference.SetLabel("ElementReference", hl)
	reference.SetURI(ViewNodeElementReferenceURI, hl)
	reference.SetOwningConcept(viewNode, hl)
	reference.SetIsCore(hl)
}

func registerTreeViewFunctions(uOfD *core.UniverseOfDiscourse) {
	uOfD.AddFunction(ManageTreeNodesURI, treeViewManageNodes)
}
