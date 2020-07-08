package editor

import (
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crleditor/crleditordomain"
	"github.com/pkg/errors"

	//	"github.com/satori/go.uuid"
	"log"
)

// treeViewManageNodes() is the callback function that manaages the tree view when base elements in the Universe of Discourse change.
// The changes being sought are the addition, removal, and re-parenting of base elements and the changes in their names.
func treeViewManageNodes(instance core.Element, changeNotification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
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
			return errors.New("treeViewManageNodes called with IndicatedConceptChanged but no underlying chanage")
		}
		switch underlyingChange.GetNatureOfChange() {
		case core.IndicatedConceptChanged:
			secondUnderlyingChange := underlyingChange.GetUnderlyingChange()
			if secondUnderlyingChange == nil {
				return errors.New("treeViewManageNodes called with IndicatedConceptChanged but no underlying chanage")
			}
			switch secondUnderlyingChange.GetNatureOfChange() {
			case core.UofDConceptAdded:
				changedElement := secondUnderlyingChange.GetPriorState()
				treeManager.addNode(changedElement, hl)
			case core.UofDConceptChanged:
				thirdUnderlyingChange := secondUnderlyingChange.GetUnderlyingChange()
				if thirdUnderlyingChange == nil {
					return errors.New("treeViewManageNodes called with UofDConceptChanged but no thirdUnderlyingChange chanage")
				}
				changedElement := thirdUnderlyingChange.GetReportingElement()
				treeManager.changeNode(changedElement, hl)
			case core.UofDConceptRemoved:
				changedElement := secondUnderlyingChange.GetPriorState()
				treeManager.removeNode(changedElement, hl)
			}
		}
	}
	return nil
}

func registerTreeViewFunctions(uOfD *core.UniverseOfDiscourse) {
	uOfD.AddFunction(crleditordomain.TreeNodeManagerURI, treeViewManageNodes)
}
