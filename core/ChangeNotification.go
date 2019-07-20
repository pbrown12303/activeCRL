// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"log"
	"strconv"
)

// NatureOfChange indicates the type of base element change:
type NatureOfChange int

// AbstractionChanged indicates that an abstraction of the Element has been modified
// ChildAbstractionChanged indicates that an abstracton of the indicated child has been modified
// ChildChanged indicates that a child of the Element has been changed
// ConceptChanged indicates that a field of the concept has been modified
// ConceptDeleted indicates that the concept has been deleted
// IndicatedConceptChanged indicates that an Element indicated by a pointer has changed
// UofDConceptChanged indicates that a concept in the UofD has changed
const (
	AbstractionChanged      = NatureOfChange(1)
	ChildAbstractionChanged = NatureOfChange(2)
	ChildChanged            = NatureOfChange(3)
	ConceptChanged          = NatureOfChange(4)
	IndicatedConceptChanged = NatureOfChange(5)
	UofDConceptAdded        = NatureOfChange(6)
	UofDConceptChanged      = NatureOfChange(7)
	UofDConceptRemoved      = NatureOfChange(8)
)

func (noc NatureOfChange) String() string {
	switch noc {
	case AbstractionChanged:
		return "AbstractionChanged"
	case ChildAbstractionChanged:
		return "ChildAbstractionChanged"
	case ChildChanged:
		return "ChildChanged"
	case ConceptChanged:
		return "ConceptChanged"
	case IndicatedConceptChanged:
		return "IndicatedConceptChanged"
	case UofDConceptAdded:
		return "UofDConceptAdded"
	case UofDConceptChanged:
		return "UofDConceptChanged"
	case UofDConceptRemoved:
		return "UofDConceptRemoved"
	}
	return "Undefined"
}

// ChangeNotification records the metadata regarding a change to a Element. It provides
// the nature of the change, the old and new values, and the reporting Element.
// It also provides the underlying change that triggered this one (if any)
type ChangeNotification struct {
	priorState       Element
	natureOfChange   NatureOfChange
	reportingElement Element
	underlyingChange *ChangeNotification
	uOfD             *UniverseOfDiscourse
}

// GetDepth returns the depth of the nested notifications within the current notification
func (cnPtr *ChangeNotification) GetDepth() int {
	return cnPtr.getDepth(0)
}

func (cnPtr *ChangeNotification) getDepth(currentDepth int) int {
	newDepth := currentDepth + 1
	if cnPtr.underlyingChange != nil {
		return cnPtr.underlyingChange.getDepth(newDepth)
	}
	return newDepth
}

func (cnPtr *ChangeNotification) isReferenced(el Element) bool {
	if cnPtr.reportingElement == el {
		return true
	} else if cnPtr.underlyingChange != nil {
		return cnPtr.underlyingChange.isReferenced(el)
	}
	return false
}

// GetNatureOfChange returns the NatureOFChange
func (cnPtr *ChangeNotification) GetNatureOfChange() NatureOfChange {
	return cnPtr.natureOfChange
}

// GetPriorState returns the state, which is a clone of the Element prior to the change
// Note that while this is an Element, it is NOT a member of the UniverseOfDiscourse
func (cnPtr *ChangeNotification) GetPriorState() Element {
	return cnPtr.priorState
}

// GetReportingElement returns the Element reporting the change
func (cnPtr *ChangeNotification) GetReportingElement() Element {
	return cnPtr.reportingElement
}

// GetReportingElementID returns the ID of the Element reporting the change
func (cnPtr *ChangeNotification) GetReportingElementID() string {
	return cnPtr.reportingElement.getConceptIDNoLock()
}

// GetUnderlyingChange returns the change notification that triggered the change being
// reported in this ChangeNotification
func (cnPtr *ChangeNotification) GetUnderlyingChange() *ChangeNotification {
	return cnPtr.underlyingChange
}

// Print prints the change notification for diagnostic purposes to the log
func (cnPtr *ChangeNotification) Print(prefix string, hl *HeldLocks) {
	if EnableNotificationPrint == true {
		startCount := 0
		cnPtr.printRecursively(prefix, hl, startCount)
	}
}

// printRecursively prints the change notification for diagnostic purposes to the log. The startCount
// indicates the depth of nesting of the print so that the printout can el indented appropriately.
func (cnPtr *ChangeNotification) printRecursively(prefix string, hl *HeldLocks, startCount int) {
	notificationType := "+++ " + cnPtr.natureOfChange.String()
	log.Printf("%s%s: \n", prefix, "### Notification Level: "+strconv.Itoa(startCount)+" Type: "+notificationType)
	if cnPtr.reportingElement == nil {
		log.Printf(prefix + "Reporting Element is nil")
	} else {
		log.Printf(prefix+"  Reporting Element: %T", cnPtr.GetReportingElement())
		log.Printf(prefix+"  Reporting Element ID: %s", cnPtr.reportingElement.getConceptIDNoLock())
		log.Printf(prefix+"  Reporting Element Version: %d", cnPtr.GetReportingElement().GetVersion(hl))
		log.Printf(prefix+"  Reporting Element Label: %s", cnPtr.GetReportingElement().GetLabel(hl))
		priorState := cnPtr.GetPriorState()
		if priorState != nil {
			jsonString, _ := priorState.MarshalJSON()
			log.Printf(prefix+"  PriorState: %s", jsonString)
		} else {
			log.Printf(prefix + "  PriorState is nil")
		}
	}
	if cnPtr.underlyingChange != nil {
		cnPtr.underlyingChange.printRecursively(prefix+"      ", hl, startCount-1)
	}
	log.Printf(prefix + "End of notification")
}
