package editor

import (
	"log"
	"strconv"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagram"
)

const treeNodeSuffix = "TreeNode"

// treeManager manages the client's tree display of the uOfD
type treeManager struct {
	treeNodeManager core.Element
	treeID          string
	editor          *CrlEditor
}

// newTreeManager creates an instance of the TreeManager
func newTreeManager(editor *CrlEditor, treeID string) *treeManager {
	var tm treeManager
	tm.editor = editor
	tm.treeID = treeID
	return &tm
}

// addChildren adds the OwnedConcepts of the supplied Element to the client's tree
func (tmPtr *treeManager) addChildren(el core.Element, hl *core.HeldLocks) {
	uOfD := tmPtr.editor.uOfD
	for id := range uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Iterator().C {
		child := uOfD.GetElement(id.(string))
		tmPtr.addNode(child, hl)
		tmPtr.addChildren(child, hl)
	}
}

// addNode adds a node to the tree
func (tmPtr *treeManager) addNode(el core.Element, hl *core.HeldLocks) {
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(crldiagram.IsDiagram(el, hl))}
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("AddTreeNode", el.GetConceptID(hl), el, additionalParameters)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if notificationResponse.Result != 0 {
		log.Print(notificationResponse.ErrorMessage)
	}
}

// addNodeRecursively adds the node and all of its descendants to the tree
func (tmPtr *treeManager) addNodeRecursively(el core.Element, hl *core.HeldLocks) {
	tmPtr.addNode(el, hl)
	uOfD := tmPtr.editor.uOfD
	for id := range uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Iterator().C {
		child := uOfD.GetElement(id.(string))
		tmPtr.addNodeRecursively(child, hl)
	}
}

// changeNode updates the tree node
func (tmPtr *treeManager) changeNode(el core.Element, hl *core.HeldLocks) {
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(crldiagram.IsDiagram(el, hl))}
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("ChangeTreeNode", el.GetConceptID(hl), el, additionalParameters)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if notificationResponse.Result != 0 {
		log.Print(notificationResponse.ErrorMessage)
	}
}

func (tmPtr *treeManager) configureUofD(hl *core.HeldLocks) (core.Element, error) {
	// Set up the tree view
	var err error
	tmPtr.treeNodeManager, err = tmPtr.editor.uOfD.CreateReplicateAsRefinementFromURI(TreeNodeManagerURI, hl)
	if err != nil {
		return nil, err
	}
	uOfDReference := tmPtr.treeNodeManager.GetFirstOwnedReferenceRefinedFromURI(TreeNodeManagerUofDReferenceURI, hl)
	uOfDReference.SetReferencedConcept(tmPtr.editor.uOfD, hl)
	return tmPtr.treeNodeManager, nil
}

// removeNode removes the tree node
func (tmPtr *treeManager) removeNode(el core.Element, hl *core.HeldLocks) {
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("DeleteTreeNode", el.GetConceptID(hl), el, nil)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if notificationResponse.Result != 0 {
		log.Print(notificationResponse.ErrorMessage)
	}
}

func (tmPtr *treeManager) getChangeNotificationBelowUofD(changeNotification *core.ChangeNotification) *core.ChangeNotification {
	if changeNotification.GetReportingElement() == tmPtr.editor.uOfD {
		return changeNotification.GetUnderlyingChange()
	} else if changeNotification.GetUnderlyingChange() != nil {
		return tmPtr.getChangeNotificationBelowUofD(changeNotification.GetUnderlyingChange())
	}
	return nil
}

// initializeTree initializes the display of the tree in the client
func (tmPtr *treeManager) initializeTree(hl *core.HeldLocks) {
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("ClearTree", "", nil, map[string]string{})
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if notificationResponse.Result != 0 {
		log.Print(notificationResponse.ErrorMessage)
	}
	for _, el := range tmPtr.editor.uOfD.GetElements() {
		if el.GetOwningConcept(hl) == nil {
			tmPtr.addNode(el, hl)
			tmPtr.addChildren(el, hl)
		}
	}
}
