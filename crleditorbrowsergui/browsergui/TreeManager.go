package browsergui

import (
	"log"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
)

// treeManager manages the client's tree display of the uOfD
type treeManager struct {
	treeID     string
	browserGUI *BrowserGUI
}

// addChildren adds the OwnedConcepts of the supplied Element to the client's tree
func (tmPtr *treeManager) addChildren(el core.Element, hl *core.Transaction) error {
	uOfD := tmPtr.browserGUI.GetUofD()
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
func (tmPtr *treeManager) addNode(el core.Element, hl *core.Transaction) error {
	if el == nil {
		return errors.New("treeManger.addNode called with nil element")
	}
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(crldiagramdomain.IsDiagram(el, hl))}
	conceptState, err2 := core.NewConceptState(el)
	if err2 != nil {
		return errors.Wrap(err2, "treeManager.addNode failed")
	}
	notificationResponse, err3 := BrowserGUISingleton.GetClientNotificationManager().SendNotification("AddTreeNode", el.GetConceptID(hl), conceptState, additionalParameters)
	if err3 != nil {
		return errors.Wrap(err3, "TreeManager.addNode failed")
	}
	if notificationResponse != nil && notificationResponse.Result != 0 {
		return errors.New("In TreeManager.addNode, got " + notificationResponse.ErrorMessage)
	}
	return nil
}

// addNodeRecursively adds the node and all of its descendants to the tree
func (tmPtr *treeManager) addNodeRecursively(el core.Element, hl *core.Transaction) error {
	err := tmPtr.addNode(el, hl)
	if err != nil {
		return errors.Wrap(err, "TreeManager.addNodeRecursively failed")
	}
	uOfD := tmPtr.browserGUI.GetUofD()
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
func (tmPtr *treeManager) changeNode(el core.Element, hl *core.Transaction) error {
	if el == nil {
		return errors.New("treeManager.changeNode called with nil Element")
	}
	icon := GetIconPath(el, hl)
	additionalParameters := map[string]string{
		"icon":      icon,
		"isDiagram": strconv.FormatBool(crldiagramdomain.IsDiagram(el, hl))}
	conceptState, err2 := core.NewConceptState(el)
	if err2 != nil {
		return errors.Wrap(err2, "treeManager.addNode failed")
	}
	notificationResponse, err := BrowserGUISingleton.GetClientNotificationManager().SendNotification("ChangeTreeNode", conceptState.ConceptID, conceptState, additionalParameters)
	if err != nil {
		return errors.Wrap(err, "TreeManager.changeNode failed")
	}
	if notificationResponse != nil && notificationResponse.Result != 0 {
		return errors.New("In TreeManager.changeNode, notification response was: " + notificationResponse.ErrorMessage)
	}

	return nil
}

// removeNode removes the tree node
func (tmPtr *treeManager) removeNode(elID string, hl *core.Transaction) error {
	if elID == "" {
		return errors.New("treeManager.removeNode called with no ConceptID")
	}
	notificationResponse, err := BrowserGUISingleton.GetClientNotificationManager().SendNotification("DeleteTreeNode", elID, nil, nil)
	if err != nil {
		return errors.Wrap(err, "TreeManager.removeNode failed")
	}
	if notificationResponse.Result != 0 {
		return errors.New("In TreeManager.removeNode, notification response was: " + notificationResponse.ErrorMessage)
	}
	return nil
}

// func (tmPtr *treeManager) getChangeNotificationBelowUofD(changeNotification *core.ChangeNotification) *core.ChangeNotification {
// 	if changeNotification.GetChangedConceptID() == "" { // only happens when uOfD is the reporting element
// 		return changeNotification.GetUnderlyingChange()
// 	} else if changeNotification.GetUnderlyingChange() != nil {
// 		return tmPtr.getChangeNotificationBelowUofD(changeNotification.GetUnderlyingChange())
// 	}
// 	return nil
// }

// initialize sets up the uOfD monitoring
func (tmPtr *treeManager) initialize(hl *core.Transaction) error {
	err := tmPtr.browserGUI.GetUofD().Register(tmPtr)
	if err != nil {
		return errors.Wrap(err, "treeManager.initialize failed")
	}
	return nil
}

// initializeTree initializes the display of the tree in the client
func (tmPtr *treeManager) initializeTree(hl *core.Transaction) error {
	notificationResponse, err := BrowserGUISingleton.GetClientNotificationManager().SendNotification("ClearTree", "", nil, map[string]string{})
	if err != nil {
		return errors.Wrap(err, "TreeManager.initializeTree failed")
	}
	if notificationResponse == nil {
		return errors.New("treeManager.initializeTree called, but no notificationResponse was received")
	}
	if notificationResponse.Result != 0 {
		return errors.New("In TreeManager.initializeTree, notification response was: " + notificationResponse.ErrorMessage)
	}
	for _, el := range tmPtr.browserGUI.GetUofD().GetElements() {
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

// Update  is the callback function that manaages the tree view when elements in the Universe of Discourse change.
// The changes being sought are the addition, removal, and re-parenting of base elements and the changes in their names.
func (tmPtr *treeManager) Update(notification *core.ChangeNotification, hl *core.Transaction) error {
	uOfD := hl.GetUniverseOfDiscourse()

	// Tracing
	if core.AdHocTrace {
		log.Printf("In treeManager.Update()")
	}

	changedElementID := notification.GetChangedConceptID()
	changedElement := uOfD.GetElement(changedElementID)
	switch notification.GetNatureOfChange() {
	case core.ConceptAdded:
		tmPtr.addNode(changedElement, hl)
	case core.ConceptChanged, core.OwningConceptChanged:
		tmPtr.changeNode(changedElement, hl)
	case core.ConceptRemoved:
		tmPtr.removeNode(changedElementID, hl)
	}

	return nil
}
