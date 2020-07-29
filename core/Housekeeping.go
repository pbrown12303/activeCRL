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
		uOfD.queueFunctionExecutions(uOfD, uOfDChangedNotification, hl)
		// Send ChildChanged to owner
		owner := el.GetOwningConcept(hl)
		if owner != nil {
			childChangedNotification := uOfD.NewForwardingChangeNotification(el, notification.GetBeforeState(), notification.GetAfterState(), ChildChanged, notification, hl)
			uOfD.queueFunctionExecutions(owner, childChangedNotification, hl)
		}
		// If owner has changed, send ChildChanged to old owner as well
		oldOwner := uOfD.GetElement(notification.GetBeforeState().OwningConceptID)
		if oldOwner != nil && oldOwner != owner {
			childChangedNotification := uOfD.NewForwardingChangeNotification(el, notification.GetBeforeState(), notification.GetAfterState(), ChildChanged, notification, hl)
			uOfD.queueFunctionExecutions(oldOwner, childChangedNotification, hl)
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
			uOfD.queueFunctionExecutions(owner, childAbstractionChangedNotification, hl)
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
			uOfD.queueFunctionExecutions(owner, childAbstractionChangedNotification, hl)
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
			uOfD.queueFunctionExecutions(owner, childChangedNotification, hl)
		}
		// Send IndicatedConceptChanged to listeners
		el.notifyListeners(notification, hl)
	case IndicatedConceptChanged:
		owner := el.GetOwningConcept(hl)
		if owner != nil && !notification.isReferenced(owner) {
			indicatedConceptChangedNotification := uOfD.NewForwardingChangeNotification(el, notification.GetBeforeState(), notification.GetAfterState(), IndicatedConceptChanged, notification, hl)
			uOfD.queueFunctionExecutions(owner, indicatedConceptChangedNotification, hl)
		}
	case UofDConceptChanged, UofDConceptAdded, UofDConceptRemoved:
		// Send IndicatedConceptChanged to listeners
		el.notifyListeners(notification, hl)
	}
	return nil
}
