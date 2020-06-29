package editor

import (
	"strconv"

	"github.com/pkg/errors"

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

// addChildren adds the OwnedConcepts of the supplied Element to the client's tree
func (tmPtr *treeManager) addChildren(el core.Element, hl *core.HeldLocks) error {
	uOfD := tmPtr.editor.uOfD
	for id := range uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Iterator().C {
		child := uOfD.GetElement(id.(string))
		if child == nil {
			return errors.New("In TreeManager.addChildren, no child found for id: " + id.(string))
		}
		err := tmPtr.addNode(child, hl)
		if err != nil {
			return errors.Wrap(err, "TreeManager.addChildren failed")
		}
		err = tmPtr.addChildren(child, hl)
		if err != nil {
			return errors.Wrap(err, "TreeManager.addChildren failed")
		}
	}
	return nil
}

// addNode adds a node to the tree
func (tmPtr *treeManager) addNode(el core.Element, hl *core.HeldLocks) error {
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(crldiagram.IsDiagram(el, hl))}
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("AddTreeNode", el.GetConceptID(hl), el, additionalParameters)
	if err != nil {
		return errors.Wrap(err, "TreeManager.addNode failed")
	}
	if notificationResponse != nil && notificationResponse.Result != 0 {
		return errors.New("In TreeManager.addNode, got " + notificationResponse.ErrorMessage)
	}
	return nil
}

// addNodeRecursively adds the node and all of its descendants to the tree
func (tmPtr *treeManager) addNodeRecursively(el core.Element, hl *core.HeldLocks) error {
	err := tmPtr.addNode(el, hl)
	if err != nil {
		return errors.Wrap(err, "TreeManager.addNodeRecursively failed")
	}
	uOfD := tmPtr.editor.uOfD
	for id := range uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Iterator().C {
		child := uOfD.GetElement(id.(string))
		if child == nil {
			return errors.New("In TreeManager.addNodeRecursively, child not found for id: " + id.(string))
		}
		err = tmPtr.addNodeRecursively(child, hl)
		if err != nil {
			return errors.Wrap(err, "TreeManager.addNodeRecursively failed")
		}
	}
	return nil
}

// changeNode updates the tree node
func (tmPtr *treeManager) changeNode(el core.Element, hl *core.HeldLocks) error {
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(crldiagram.IsDiagram(el, hl))}
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("ChangeTreeNode", el.GetConceptID(hl), el, additionalParameters)
	if err != nil {
		return errors.Wrap(err, "TreeManager.changeNode failed")
	}
	if notificationResponse != nil && notificationResponse.Result != 0 {
		return errors.New("In TreeManager.changeNode, notification response was: " + notificationResponse.ErrorMessage)
	}
	return nil
}

// removeNode removes the tree node
func (tmPtr *treeManager) removeNode(el core.Element, hl *core.HeldLocks) error {
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("DeleteTreeNode", el.GetConceptID(hl), el, nil)
	if err != nil {
		return errors.Wrap(err, "TreeManager.removeNode failed")
	}
	if notificationResponse.Result != 0 {
		return errors.New("In TreeManager.removeNode, notification response was: " + notificationResponse.ErrorMessage)
	}
	return nil
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
func (tmPtr *treeManager) initializeTree(hl *core.HeldLocks) error {
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("ClearTree", "", nil, map[string]string{})
	if err != nil {
		return errors.Wrap(err, "TreeManager.initializeTree failed")
	}
	if notificationResponse == nil {
		return errors.New("treeManager.initializeTree called, but no notificationResponse was received")
	}
	if notificationResponse.Result != 0 {
		return errors.New("In TreeManager.initializeTree, notification response was: " + notificationResponse.ErrorMessage)
	}
	for _, el := range tmPtr.editor.uOfD.GetElements() {
		if el.GetOwningConcept(hl) == nil {
			err = tmPtr.addNode(el, hl)
			if err != nil {
				return errors.Wrap(err, "TreeManager.initialzeTree failed")
			}
			err = tmPtr.addChildren(el, hl)
			if err != nil {
				return errors.Wrap(err, "TreeManager.initializeTree failed")
			}
		}
	}
	return nil
}
