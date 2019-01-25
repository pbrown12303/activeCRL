package editor

import (
	"log"
	"strconv"

	"github.com/pbrown12303/activeCRL/core"
)

const treeNodeSuffix = "TreeNode"

// TreeManager manages the client's tree display of the uOfD
type TreeManager struct {
	manageNodesFunction core.Element
	treeID              string
	rootElements        map[string]core.Element
	uOfD                core.UniverseOfDiscourse
}

// NewTreeManager creates an instance of the TreeManager
func NewTreeManager(uOfD core.UniverseOfDiscourse, treeID string, hl *core.HeldLocks) *TreeManager {
	var treeManager TreeManager
	treeManager.uOfD = uOfD
	treeManager.treeID = treeID
	treeManager.rootElements = make(map[string]core.Element)

	// Set up the tree view
	var err error
	treeManager.manageNodesFunction, err = uOfD.CreateReplicateAsRefinementFromURI(ManageTreeNodesURI, hl)
	if err != nil {
		log.Print(err)
	}
	uOfDReference := treeManager.manageNodesFunction.GetFirstOwnedReferenceRefinedFromURI(ManageNodesUofDReferenceURI, hl)
	uOfDReference.SetReferencedConcept(uOfD, hl)
	treeManager.manageNodesFunction.SetIsCoreRecursively(hl)

	return &treeManager
}

// AddChildren adds the OwnedConcepts of the supplied Element to the client's tree
func (tmPtr *TreeManager) AddChildren(el core.Element, hl *core.HeldLocks) {
	switch el.(type) {
	case core.Element:
		for _, child := range el.GetOwnedConcepts(hl) {
			tmPtr.AddNode(child, hl)
			tmPtr.AddChildren(child, hl)
		}
	}
}

// AddNode adds a node to the tree
func (tmPtr *TreeManager) AddNode(el core.Element, hl *core.HeldLocks) {
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(IsDiagram(el, hl))}
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("AddTreeNode", el.GetConceptID(hl), el, additionalParameters)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if notificationResponse.Result != 0 {
		log.Print(notificationResponse.ErrorMessage)
	}
}

// ChangeNode updates the tree node
func (tmPtr *TreeManager) ChangeNode(el core.Element, hl *core.HeldLocks) {
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("ChangeTreeNode", el.GetConceptID(hl), el, nil)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if notificationResponse.Result != 0 {
		log.Print(notificationResponse.ErrorMessage)
	}
}

// RemoveNode removes the tree node
func (tmPtr *TreeManager) RemoveNode(el core.Element, hl *core.HeldLocks) {
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("DeleteTreeNode", el.GetConceptID(hl), el, nil)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if notificationResponse.Result != 0 {
		log.Print(notificationResponse.ErrorMessage)
	}
}

func (tmPtr *TreeManager) getChangeNotificationBelowUofD(changeNotification *core.ChangeNotification) *core.ChangeNotification {
	if changeNotification.GetReportingElement() == tmPtr.uOfD {
		return changeNotification.GetUnderlyingChange()
	} else if changeNotification.GetUnderlyingChange() != nil {
		return tmPtr.getChangeNotificationBelowUofD(changeNotification.GetUnderlyingChange())
	}
	return nil
}

// InitializeTree initializes the display of the tree in the client
func (tmPtr *TreeManager) InitializeTree(hl *core.HeldLocks) {
	for _, el := range tmPtr.uOfD.GetElements() {
		if el.GetOwningConcept(hl) == nil {
			tmPtr.AddNode(el, hl)
			tmPtr.AddChildren(el, hl)
		}
	}
}

// IsDiagram returns true if the supplied element is a crldiagram
func IsDiagram(el core.Element, hl *core.HeldLocks) bool {
	switch el.(type) {
	case core.Element:
		return el.IsRefinementOf(CrlEditorSingleton.GetDiagramManager().abstractDiagram, hl)
	}
	return false
}

// func getIDWithoutSuffix(stringWithSuffix string, suffix string) string {
// 	if len(stringWithSuffix) > len(suffix) && stringWithSuffix[len(stringWithSuffix)-len(suffix):] == suffix {
// 		return stringWithSuffix[:len(stringWithSuffix)-len(suffix)]
// 	}
// 	return ""
// }

// func getIDWithoutPrefix(stringWithPrefix string, prefix string) string {
// 	if len(stringWithPrefix) > len(prefix) && stringWithPrefix[:len(prefix)] == prefix {
// 		return stringWithPrefix[len(prefix):]
// 	}
// 	return ""
// }
