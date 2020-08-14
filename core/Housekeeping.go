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
	elCurrentState, err := NewConceptState(el)
	if err != nil {
		return errors.Wrap(err, "coreHousekeeping failed")
	}
	switch notification.GetNatureOfChange() {
	case ConceptChanged:
		// Notify Universe of Discourse
		uOfDChangedNotification := uOfD.newUniverseOfDiscourseChangeNotification(notification)
		err := uOfD.queueFunctionExecutions(uOfD, uOfDChangedNotification, hl)
		if err != nil {
			return errors.Wrap(err, "Housekeeping.go coreHousekeeping failed")
		}
		// Send ChildChanged to owner
		owner := el.GetOwningConcept(hl)
		if owner != nil {
			childChangedNotification := uOfD.NewForwardingChangeNotification(el, notification.GetBeforeState(), notification.GetAfterState(), ChildChanged, notification, hl)
			err = uOfD.queueFunctionExecutions(owner, childChangedNotification, hl)
			if err != nil {
				return errors.Wrap(err, "Housekeeping.go coreHousekeeping failed")
			}
		}
		// If owner has changed, send ChildChanged to old owner as well
		oldOwner := uOfD.GetElement(notification.GetBeforeState().OwningConceptID)
		if oldOwner != nil && oldOwner != owner {
			childChangedNotification := uOfD.NewForwardingChangeNotification(el, notification.GetBeforeState(), notification.GetAfterState(), ChildChanged, notification, hl)
			err = uOfD.queueFunctionExecutions(oldOwner, childChangedNotification, hl)
			if err != nil {
				return errors.Wrap(err, "Housekeeping.go coreHousekeeping failed")
			}
		}
		// Send IndicatedConceptChanged to listeners
		el.notifyListeners(notification, hl)
	case AbstractionChanged:
		// Increment version
		el.incrementVersion(hl)
		// Send ChildAbstractionChanged to owner
		owner := el.GetOwningConcept(hl)
		if owner != nil {
			childAbstractionChangedNotification := uOfD.NewForwardingChangeNotification(el, elCurrentState, elCurrentState, ChildAbstractionChanged, notification, hl)
			err := uOfD.queueFunctionExecutions(owner, childAbstractionChangedNotification, hl)
			if err != nil {
				return errors.Wrap(err, "Housekeeping.go coreHousekeeping failed")
			}
		}
		// Send IndicatedConceptChanged or AbstractionChanged to listeners
		el.notifyListeners(notification, hl)
	case ChildAbstractionChanged:
		// Increment version
		el.incrementVersion(hl)
		// Send ChildAbstractionChanged to owner
		owner := el.GetOwningConcept(hl)
		if owner != nil {
			childAbstractionChangedNotification := uOfD.NewForwardingChangeNotification(el, elCurrentState, elCurrentState, ChildAbstractionChanged, notification, hl)
			err := uOfD.queueFunctionExecutions(owner, childAbstractionChangedNotification, hl)
			if err != nil {
				return errors.Wrap(err, "Housekeeping.go coreHousekeeping failed")
			}
		}
		// Send IndicatedConceptChanged or AbstractionChanged to listeners
		el.notifyListeners(notification, hl)
	case ChildChanged:
		// Increment version
		el.incrementVersion(hl)
		// Send ChildChanged to owner
		owner := el.GetOwningConcept(hl)
		if owner != nil {
			childChangedNotification := uOfD.NewForwardingChangeNotification(el, elCurrentState, elCurrentState, ChildChanged, notification, hl)
			err := uOfD.queueFunctionExecutions(owner, childChangedNotification, hl)
			if err != nil {
				return errors.Wrap(err, "Housekeeping.go coreHousekeeping failed")
			}
		}
		// Send IndicatedConceptChanged to listeners
		el.notifyListeners(notification, hl)
	case IndicatedConceptChanged:
		owner := el.GetOwningConcept(hl)
		if owner != nil && !notification.isReferenced(owner) {
			indicatedConceptChangedNotification := uOfD.NewForwardingChangeNotification(el, notification.GetBeforeState(), notification.GetAfterState(), IndicatedConceptChanged, notification, hl)
			err := uOfD.queueFunctionExecutions(owner, indicatedConceptChangedNotification, hl)
			if err != nil {
				return errors.Wrap(err, "Housekeeping.go coreHousekeeping failed")
			}
		}
	case UofDConceptChanged, UofDConceptAdded, UofDConceptRemoved:
		// Send IndicatedConceptChanged to listeners
		err := el.notifyListeners(notification, hl)
		if err != nil {
			return errors.Wrap(err, "Housekeeping.go coreHousekeeping failed")
		}
	}
	return nil
}
