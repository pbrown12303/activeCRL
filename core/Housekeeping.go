package core

var coreHousekeepingURI = CorePrefix + "coreHousekeeping"

// coreHousekeeping does the housekeeping for the core concepts
func coreHousekeeping(el Element, notification *ChangeNotification, uOfD UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.WriteLockElement(el)
	switch notification.GetNatureOfChange() {
	case ConceptChanged:
		// Notify Universe of Discourse
		uOfDChangedNotification := uOfD.NewUniverseOfDiscourseChangeNotification(notification)
		uOfD.queueFunctionExecutions(uOfD, uOfDChangedNotification, hl)
		// Send ChildChanged to owner
		owner := el.GetOwningConcept(hl)
		if owner != nil {
			childChangedNotification := uOfD.NewForwardingChangeNotification(el, ChildChanged, notification)
			uOfD.queueFunctionExecutions(owner, childChangedNotification, hl)
		}
		// Send IndicatedConceptChanged to listeners
		el.notifyListeners(notification, hl)
	case AbstractionChanged:
		// Increment version
		el.incrementVersion(hl)
		// Send ChildAbstractionChanged to owner
		owner := el.GetOwningConcept(hl)
		if owner != nil {
			childAbstractionChangedNotification := uOfD.NewForwardingChangeNotification(el, ChildAbstractionChanged, notification)
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
			childAbstractionChangedNotification := uOfD.NewForwardingChangeNotification(el, ChildAbstractionChanged, notification)
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
			childChangedNotification := uOfD.NewForwardingChangeNotification(el, ChildChanged, notification)
			uOfD.queueFunctionExecutions(owner, childChangedNotification, hl)
		}
		// Send IndicatedConceptChanged to listeners
		el.notifyListeners(notification, hl)
	case IndicatedConceptChanged:
		owner := el.GetOwningConcept(hl)
		if owner != nil {
			indicatedConceptChangedNotification := uOfD.NewForwardingChangeNotification(el, IndicatedConceptChanged, notification)
			uOfD.queueFunctionExecutions(owner, indicatedConceptChangedNotification, hl)
		}
	case UofDChanged:
		// Send IndicatedConceptChanged to listeners
		el.notifyListeners(notification, hl)
	}
}
