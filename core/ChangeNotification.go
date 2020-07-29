// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"reflect"
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

// ConceptState is a flattened representation of all concept types. It is used to capture the current state of a concept
type ConceptState struct {
	// Element fields
	ConceptID       string
	ConceptType     string
	OwningConceptID string
	Label           string
	Definition      string
	URI             string
	Version         string
	IsCore          string
	ReadOnly        string
	// Literal fields
	LiteralValue string
	// Reference fields
	ReferencedConceptID      string
	ReferencedAttributeName  string
	ReferencedConceptVersion string
	// Refinement Fields
	AbstractConceptID      string
	AbstractConceptVersion string
	RefinedConceptID       string
	RefinedConceptVersion  string
}

// NewConceptState copies the state of an Element into a ConceptState struct
func NewConceptState(el Element) (*ConceptState, error) {
	if el == nil {
		return nil, errors.New("NewConceptState called with nil element")
	}
	mJSON, err := el.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "NewConceptState failed")
	}
	var newConceptState ConceptState
	err = json.Unmarshal([]byte(mJSON), &newConceptState)
	if err != nil {
		return nil, errors.Wrap(err, "NewConceptState failed")
	}
	newConceptState.ConceptType = reflect.TypeOf(el).String()
	return &newConceptState, nil
}

// ChangeNotification records the metadata regarding a change to a Element. It provides
// the nature of the change, the old and new values, and the reporting Element.
// It also provides the underlying change that triggered this one (if any)
type ChangeNotification struct {
	reportingElementID    string
	reportingElementLabel string
	reportingElementType  string
	natureOfChange        NatureOfChange
	beforeState           *ConceptState
	afterState            *ConceptState
	underlyingChange      *ChangeNotification
	uOfD                  *UniverseOfDiscourse
}

// GetAfterState returns the state of the Element after the change
func (cnPtr *ChangeNotification) GetAfterState() *ConceptState {
	return cnPtr.afterState
}

// GetBeforeState returns the state of the Element before the change
// Note that while this is an Element, it is NOT a member of the UniverseOfDiscourse
func (cnPtr *ChangeNotification) GetBeforeState() *ConceptState {
	return cnPtr.beforeState
}

// GetChangedConceptID returns the ID of the Element impacted by the change
func (cnPtr *ChangeNotification) GetChangedConceptID() string {
	if cnPtr.afterState != nil {
		return cnPtr.afterState.ConceptID
	} else if cnPtr.beforeState != nil {
		return cnPtr.beforeState.ConceptID
	}
	return ""
}

// GetChangedConceptLabel returns the label of the Element impacted by the change
func (cnPtr *ChangeNotification) GetChangedConceptLabel() string {
	if cnPtr.afterState != nil {
		return cnPtr.afterState.Label
	} else if cnPtr.beforeState != nil {
		return cnPtr.beforeState.Label
	}
	return ""
}

// GetChangedConceptType returns the typeString of the Element impacted by the change
func (cnPtr *ChangeNotification) GetChangedConceptType() string {
	if cnPtr.afterState != nil {
		return cnPtr.afterState.ConceptType
	} else if cnPtr.beforeState != nil {
		return cnPtr.beforeState.ConceptType
	}
	return ""
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

// GetNatureOfChange returns the NatureOFChange
func (cnPtr *ChangeNotification) GetNatureOfChange() NatureOfChange {
	return cnPtr.natureOfChange
}

// GetReportingElementID returns the ID of the element sending the notification
func (cnPtr *ChangeNotification) GetReportingElementID() string {
	return cnPtr.reportingElementID
}

// GetReportingElementLabel returns the Label of the element sending the notification
func (cnPtr *ChangeNotification) GetReportingElementLabel() string {
	return cnPtr.reportingElementLabel
}

// GetReportingElementType returns the Type of the element sending the notification
func (cnPtr *ChangeNotification) GetReportingElementType() string {
	return cnPtr.reportingElementType
}

// GetUnderlyingChange returns the change notification that triggered the change being
// reported in this ChangeNotification
func (cnPtr *ChangeNotification) GetUnderlyingChange() *ChangeNotification {
	return cnPtr.underlyingChange
}

func (cnPtr *ChangeNotification) isReferenced(el Element) bool {
	if cnPtr.GetChangedConceptID() == el.getConceptIDNoLock() {
		return true
	} else if cnPtr.underlyingChange != nil {
		return cnPtr.underlyingChange.isReferenced(el)
	}
	return false
}

// Print prints the change notification for diagnostic purposes to the log
func (cnPtr *ChangeNotification) Print(prefix string, hl *HeldLocks) {
	if EnableNotificationPrint == true {
		startCount := 0
		cnPtr.printRecursively(prefix, hl, startCount)
	}
}

// printRecursively prints the change notification for diagnostic purposes to the log. The startCount
// indicates the depth of nesting of the print so that the printout can be indented appropriately.
func (cnPtr *ChangeNotification) printRecursively(prefix string, hl *HeldLocks, startCount int) {
	notificationType := "+++ " + cnPtr.natureOfChange.String()
	log.Printf("%s%s: \n", prefix, "### Notification Level: "+strconv.Itoa(startCount)+" Type: "+notificationType)
	if cnPtr.afterState != nil {
		log.Printf(prefix+"  AfterState: %+v", cnPtr.afterState)
	}
	if cnPtr.beforeState != nil {
		log.Printf(prefix+"  BeforeState: %s", cnPtr.beforeState)
	}
	if cnPtr.underlyingChange != nil {
		cnPtr.underlyingChange.printRecursively(prefix+"      ", hl, startCount-1)
	}
	log.Printf(prefix + "End of notification")
}
