package core

import (
	"github.com/pkg/errors"
)

var coreHousekeepingURI = CorePrefix + "coreHousekeeping"

// coreHousekeeping does the housekeeping for the core concepts
func coreHousekeeping(el Element, notification *ChangeNotification, trans *Transaction) error {
	uOfD := trans.uOfD
	trans.ReadLockElement(el)
	// Supress circular notifications
	underlyingNotification := notification.GetUnderlyingChange()
	if underlyingNotification != nil && HasReportedPreviously(el.GetConceptID(trans), underlyingNotification) {
		return nil
	}
	// Notify listeners
	err := el.notifyListeners(notification, trans)
	if err != nil {
		return errors.Wrap(err, "coreHousekeeping failed")
	}
	// Notify owner if needed
	switch el.(type) {
	case Reference:
		if el.GetOwningConcept(trans) != nil && !(notification.GetNatureOfChange() == OwningConceptChanged && notification.GetReportingElementID() != el.GetConceptID(trans)) {
			forwardingNotification, err := uOfD.NewForwardingChangeNotification(el, ForwardedChange, notification, trans)
			if err != nil {
				return errors.Wrap(err, "coreHousekeeping failed")
			}
			err = uOfD.queueFunctionExecutions(el.GetOwningConcept(trans), forwardingNotification, trans)
			if err != nil {
				return errors.Wrap(err, "element.SetDefinition failed")
			}
		}
	}
	return nil
}

// HasReportedPreviously checks to see whether the element was a reporting element in the notification or one of its nested notifications
func HasReportedPreviously(elID string, notification *ChangeNotification) bool {
	if notification.GetReportingElementID() == elID {
		return true
	}
	nestedNotification := notification.GetUnderlyingChange()
	if nestedNotification != nil {
		return HasReportedPreviously(elID, nestedNotification)
	}
	return false
}
