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

// ConceptChanged indicates that an attribute of the concept that is NOT an element reference has changed
// OwningConceptChanged indicates that the ownership of the concept has changed
// ReferencedConceptChanged indicates that a different Element is being referenced
// AbstractConceptChanged indicates that a different Element is now the abstract concept
// RefinedConceptChanged indicates that a different Element is now the refined concept
const (
	ConceptAdded             = NatureOfChange(1)
	ConceptChanged           = NatureOfChange(2)
	ConceptRemoved           = NatureOfChange(3)
	OwningConceptChanged     = NatureOfChange(4)
	ReferencedConceptChanged = NatureOfChange(5)
	AbstractConceptChanged   = NatureOfChange(6)
	RefinedConceptChanged    = NatureOfChange(7)
	ForwardedChange          = NatureOfChange(8)
)

func (noc NatureOfChange) String() string {
	switch noc {
	case ConceptAdded:
		return "ConceptAdded"
	case ConceptChanged:
		return "ConceptChanged"
	case ConceptRemoved:
		return "ConceptRemoved"
	case OwningConceptChanged:
		return "OwningConceptChanged"
	case ReferencedConceptChanged:
		return "ReferencedConceptChanged"
	case AbstractConceptChanged:
		return "AbstractConceptChanged"
	case RefinedConceptChanged:
		return "RefinedConceptChanged"
	case ForwardedChange:
		return "ForwardedChange"
	}
	return "Undefined"
}

// ConceptState is a flattened representation of all concept types. It is used to capture the current state of a concept
type ConceptState struct {
	// Element fields
	ConceptID                   string
	ConceptType                 string
	OwningConceptID             string
	Label                       string
	Definition                  string
	URI                         string
	Version                     string
	IsCore                      string
	ReadOnly                    string
	ForwardNotificationsToOwner string
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

// ChangeNotification records the data and metadata regarding a change to a Element. It provides
// the nature of the change, the old and new values, and the reporting Element.
// It also provides the underlying change that triggered this one (if any)
type ChangeNotification struct {
	reportingElementState *ConceptState
	natureOfChange        NatureOfChange
	beforeConceptState    *ConceptState
	afterConceptState     *ConceptState
	beforeReferencedState *ConceptState
	afterReferencedState  *ConceptState
	underlyingChange      *ChangeNotification
	uOfD                  *UniverseOfDiscourse
}

// GetAfterConceptState returns the state of the Element after the change
func (cnPtr *ChangeNotification) GetAfterConceptState() *ConceptState {
	return cnPtr.afterConceptState
}

// GetAfterReferencedState returns the state of the referenced Element after the change
func (cnPtr *ChangeNotification) GetAfterReferencedState() *ConceptState {
	return cnPtr.afterReferencedState
}

// GetBeforeConceptState returns the state of the Element before the change
func (cnPtr *ChangeNotification) GetBeforeConceptState() *ConceptState {
	return cnPtr.beforeConceptState
}

// GetBeforeReferencedState returns the state of the referenced Element before the change
func (cnPtr *ChangeNotification) GetBeforeReferencedState() *ConceptState {
	return cnPtr.beforeReferencedState
}

// GetChangedConceptID returns the ID of the Element impacted by the change
func (cnPtr *ChangeNotification) GetChangedConceptID() string {
	if cnPtr.afterConceptState != nil {
		return cnPtr.afterConceptState.ConceptID
	} else if cnPtr.beforeConceptState != nil {
		return cnPtr.beforeConceptState.ConceptID
	}
	return ""
}

// GetChangedConceptLabel returns the label of the Element impacted by the change
func (cnPtr *ChangeNotification) GetChangedConceptLabel() string {
	if cnPtr.afterConceptState != nil {
		return cnPtr.afterConceptState.Label
	} else if cnPtr.beforeConceptState != nil {
		return cnPtr.beforeConceptState.Label
	}
	return ""
}

// GetChangedConceptType returns the typeString of the Element impacted by the change
func (cnPtr *ChangeNotification) GetChangedConceptType() string {
	if cnPtr.afterConceptState != nil {
		return cnPtr.afterConceptState.ConceptType
	} else if cnPtr.beforeConceptState != nil {
		return cnPtr.beforeConceptState.ConceptType
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
	return cnPtr.reportingElementState.ConceptID
}

// GetReportingElementLabel returns the Label of the element sending the notification
func (cnPtr *ChangeNotification) GetReportingElementLabel() string {
	return cnPtr.reportingElementState.Label
}

// GetReportingElementState returns the State of the element sending the notification
// If this is nil, the report is coming from the uOfD
func (cnPtr *ChangeNotification) GetReportingElementState() *ConceptState {
	return cnPtr.reportingElementState
}

// GetReportingElementType returns the Type of the element sending the notification
func (cnPtr *ChangeNotification) GetReportingElementType() string {
	return cnPtr.reportingElementState.ConceptType
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
	if cnPtr.afterConceptState != nil {
		log.Printf(prefix+"  AfterState: %+v", cnPtr.afterConceptState)
	}
	if cnPtr.beforeConceptState != nil {
		log.Printf(prefix+"  BeforeState: %s", cnPtr.beforeConceptState)
	}
	if cnPtr.underlyingChange != nil {
		cnPtr.underlyingChange.printRecursively(prefix+"      ", hl, startCount-1)
	}
	log.Printf(prefix + "End of notification")
}
