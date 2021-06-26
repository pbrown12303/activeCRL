package browsergui

import (
	"log"

	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
)

// propertyManager manages the client's display of the selectedElement's properties
type propertyManager struct {
	browserGUI *BrowserGUI
}

// initialize sets up the uOfD monitoring
func (pmPtr *propertyManager) initialize(hl *core.Transaction) error {
	err := pmPtr.browserGUI.GetUofD().Register(pmPtr)
	if err != nil {
		return errors.Wrap(err, "propertyManager.initialize failed")
	}
	return nil
}

// Update  is the callback function that manaages the properties view for the selected element when elements in the Universe of Discourse change.
func (pmPtr *propertyManager) Update(notification *core.ChangeNotification, hl *core.Transaction) error {
	uOfD := hl.GetUniverseOfDiscourse()

	// Tracing
	if core.AdHocTrace {
		log.Printf("In propertyManager.Update()")
	}

	changedElementID := notification.GetChangedConceptID()
	changedElement := uOfD.GetElement(changedElementID)
	if changedElement == nil {
		return errors.New("propertyManager.changeNode called with nil Element")
	}
	additionalParameters := map[string]string{}
	conceptState, err2 := core.NewConceptState(changedElement)
	if err2 != nil {
		return errors.Wrap(err2, "propertyManager.Update failed")
	}
	if pmPtr.browserGUI.editor.GetCurrentSelection() == changedElement && notification.GetNatureOfChange() != core.ConceptRemoved {
		notificationResponse, err := pmPtr.browserGUI.GetClientNotificationManager().SendNotification("UpdateProperties", conceptState.ConceptID, conceptState, additionalParameters)
		if err != nil {
			return errors.Wrap(err, "propertyManager.Update failed")
		}
		if notificationResponse != nil && notificationResponse.Result != 0 {
			return errors.New("In propertyManager.Update, notification response was: " + notificationResponse.ErrorMessage)
		}
	}

	return nil
}
