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
	it := uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Iterator()
	defer it.Stop()
	for id := range it.C {
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
	if el == nil {
		return errors.New("treeManger.addNode called with nil element")
	}
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(crldiagram.IsDiagram(el, hl))}
	conceptState, err2 := core.NewConceptState(el)
	if err2 != nil {
		return errors.Wrap(err2, "treeManager.addNode failed")
	}
	notificationResponse, err3 := CrlEditorSingleton.GetClientNotificationManager().SendNotification("AddTreeNode", el.GetConceptID(hl), conceptState, additionalParameters)
	if err3 != nil {
		return errors.Wrap(err3, "TreeManager.addNode failed")
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
	it := uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Iterator()
	defer it.Stop()
	for id := range it.C {
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
	if el == nil {
		return errors.New("treeManager.changeNode called with nil Element")
	}
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(crldiagram.IsDiagram(el, hl))}
	conceptState, err2 := core.NewConceptState(el)
	if err2 != nil {
		return errors.Wrap(err2, "treeManager.addNode failed")
	}
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("ChangeTreeNode", conceptState.ConceptID, conceptState, additionalParameters)
	if err != nil {
		return errors.Wrap(err, "TreeManager.changeNode failed")
	}
	if notificationResponse != nil && notificationResponse.Result != 0 {
		return errors.New("In TreeManager.changeNode, notification response was: " + notificationResponse.ErrorMessage)
	}
	return nil
}

// removeNode removes the tree node
func (tmPtr *treeManager) removeNode(elID string, hl *core.HeldLocks) error {
	if elID == "" {
		return errors.New("treeManager.removeNode called with no ConceptID")
	}
	notificationResponse, err := CrlEditorSingleton.GetClientNotificationManager().SendNotification("DeleteTreeNode", elID, nil, nil)
	if err != nil {
		return errors.Wrap(err, "TreeManager.removeNode failed")
	}
	if notificationResponse.Result != 0 {
		return errors.New("In TreeManager.removeNode, notification response was: " + notificationResponse.ErrorMessage)
	}
	return nil
}

func (tmPtr *treeManager) getChangeNotificationBelowUofD(changeNotification *core.ChangeNotification) *core.ChangeNotification {
	if changeNotification.GetChangedConceptID() == "" { // only happens when uOfD is the reporting element
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
