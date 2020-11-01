package core

import (
	"github.com/pkg/errors"
)

var coreHousekeepingURI = CorePrefix + "coreHousekeeping"

// coreHousekeeping does the housekeeping for the core concepts
func coreHousekeeping(el Element, notification *ChangeNotification, uOfD *UniverseOfDiscourse) error {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.ReadLockElement(el)
	// Notify listeners
	err := el.notifyListeners(notification, hl)
	if err != nil {
		return errors.Wrap(err, "coreHousekeeping failed")
	}
	// Notify owner if needed
	if el.GetOwningConcept(hl) != nil && el.GetForwardNotificationsToOwner(hl) == true && notification.GetNatureOfChange() != OwningConceptChanged {
		forwardingNotification, err := uOfD.NewForwardingChangeNotification(el, ForwardedChange, notification, hl)
		if err != nil {
			return errors.Wrap(err, "coreHousekeeping failed")
		}
		err = uOfD.queueFunctionExecutions(el.GetOwningConcept(hl), forwardingNotification, hl)
		if err != nil {
			return errors.Wrap(err, "element.SetDefinition failed")
		}
	}
	return nil
}
