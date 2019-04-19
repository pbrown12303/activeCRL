package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// AttributeName indicates the attribute being referenced (if any):
type AttributeName int

// NoAttribute indicates that no attribute is being referenced
// OwningConceptID     indicates that the OwningConceptID is being referenced
// ReferencedConceptID indicates that the ReferencedConceptID is being referenced
// AbstractConceptID   indicates that the AbstractConceptID is being referenced
// RefinedConceptID    indicates that the RefinedConceptID is being referenced
// LiteralValue       indicates that the LiteralValue is being referenced
const (
	NoAttribute         = AttributeName(0)
	OwningConceptID     = AttributeName(1)
	ReferencedConceptID = AttributeName(2)
	AbstractConceptID   = AttributeName(3)
	RefinedConceptID    = AttributeName(4)
	LiteralValue        = AttributeName(5)
)

func (an AttributeName) String() string {
	switch an {
	case NoAttribute:
		return "NoAttribute"
	case OwningConceptID:
		return "OwningConceptID"
	case ReferencedConceptID:
		return "ReferencedConceptID"
	case AbstractConceptID:
		return "AbstractConceptID"
	case RefinedConceptID:
		return "RefinedConceptID"
	case LiteralValue:
		return "LiteralValue"
	}
	return "Undefined"
}

// FindAttributeName takes a string version of the name and returns the corresponding AttributeName enumeration value
func FindAttributeName(stringName string) (AttributeName, error) {
	switch stringName {
	case "NoAttribute":
		return NoAttribute, nil
	case "OwningConceptID":
		return OwningConceptID, nil
	case "ReferencedConceptID":
		return ReferencedConceptID, nil
	case "AbstractConceptID":
		return AbstractConceptID, nil
	case "RefinedConceptID":
		return RefinedConceptID, nil
	case "LiteralValue":
		return LiteralValue, nil
	}
	return NoAttribute, errors.New("NewAttribute value not found for stringName: " + stringName)
}

type reference struct {
	element
	ReferencedConceptID string
	// referencedConcept is a cache for convenience
	referencedConcept        *cachedPointer
	ReferencedAttributeName  AttributeName
	ReferencedConceptVersion int
}

func (rPtr *reference) clone(hl *HeldLocks) Reference {
	hl.ReadLockElement(rPtr)
	var ref reference
	ref.initializeReference("", "")
	ref.cloneAttributes(rPtr, hl)
	return &ref
}

func (rPtr *reference) cloneAttributes(source *reference, hl *HeldLocks) {
	rPtr.element.cloneAttributes(&source.element, hl)
	rPtr.ReferencedConceptID = source.ReferencedConceptID
	rPtr.referencedConcept.setIndicatedConceptID(source.referencedConcept.getIndicatedConceptID())
	rPtr.referencedConcept.setIndicatedConcept(source.referencedConcept.getIndicatedConcept())
	rPtr.referencedConcept.parentConceptID = source.referencedConcept.parentConceptID
	rPtr.ReferencedConceptVersion = source.ReferencedConceptVersion
}

// GetReferencedConcept returns the element representing  the concept being referenced
// Note that this is a cached value
func (rPtr *reference) GetReferencedConcept(hl *HeldLocks) Element {
	hl.ReadLockElement(rPtr)
	return rPtr.getReferencedConceptNoLock()
}

func (rPtr *reference) getReferencedConceptNoLock() Element {
	cachedConcept := rPtr.referencedConcept.getIndicatedConcept()
	if cachedConcept == nil && rPtr.ReferencedConceptID != "" {
		cachedConcept = rPtr.uOfD.GetElement(rPtr.ReferencedConceptID)
		if cachedConcept != nil {
			rPtr.referencedConcept.setIndicatedConcept(cachedConcept)
		}
	}
	return cachedConcept
}

// GetReferencedConceptID returns the identifier of the concept being referenced
func (rPtr *reference) GetReferencedConceptID(hl *HeldLocks) string {
	hl.ReadLockElement(rPtr)
	return rPtr.ReferencedConceptID
}

// GetReferencedConceptAttribute returns an indicator of which attribute is being referenced (if any)
func (rPtr *reference) GetReferencedAttributeName(hl *HeldLocks) AttributeName {
	hl.ReadLockElement(rPtr)
	return rPtr.ReferencedAttributeName
}

// GetReferencedConceptVersion returns the last known version of the referenced concept
func (rPtr *reference) GetReferencedConceptVersion(hl *HeldLocks) int {
	hl.ReadLockElement(rPtr)
	return rPtr.ReferencedConceptVersion
}

func (rPtr *reference) initializeReference(conceptID string, uri string) {
	rPtr.initializeElement(conceptID, uri)
	rPtr.referencedConcept = newCachedPointer(rPtr.getConceptIDNoLock(), false)
}

func (rPtr *reference) isEquivalent(hl1 *HeldLocks, el *reference, hl2 *HeldLocks) bool {
	hl1.ReadLockElement(rPtr)
	hl2.ReadLockElement(el)
	if rPtr.ReferencedConceptID != el.ReferencedConceptID {
		return false
	}
	if rPtr.ReferencedAttributeName != el.ReferencedAttributeName {
		return false
	}
	if rPtr.referencedConcept.isEquivalent(el.referencedConcept) != true {
		return false
	}
	if rPtr.ReferencedConceptVersion != el.ReferencedConceptVersion {
		return false
	}
	return rPtr.element.isEquivalent(hl1, &el.element, hl2)
}

// MarshalJSON produces a byte string JSON representation of the Element
func (rPtr *reference) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(rPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := rPtr.marshalReferenceFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (rPtr *reference) marshalReferenceFields(buffer *bytes.Buffer) error {
	buffer.WriteString(fmt.Sprintf("\"ReferencedConceptID\":\"%s\",", rPtr.ReferencedConceptID))
	buffer.WriteString(fmt.Sprintf("\"ReferencedAttributeName\":\"%s\",", rPtr.ReferencedAttributeName.String()))
	buffer.WriteString(fmt.Sprintf("\"ReferencedConceptVersion\":\"%d\",", rPtr.ReferencedConceptVersion))
	rPtr.element.marshalElementFields(buffer)
	return nil
}

// recoverReferenceFields() is used when de-serializing an element. The activities in restoring the
// element are not considered changes so the version counter is not incremented and the monitors of this
// element are not notified of chaanges.
func (rPtr *reference) recoverReferenceFields(unmarshaledData *map[string]json.RawMessage, hl *HeldLocks) error {
	err := rPtr.recoverElementFields(unmarshaledData, hl)
	if err != nil {
		return err
	}
	// ReferencedConceptID
	var recoveredReferencedConceptID string
	err = json.Unmarshal((*unmarshaledData)["ReferencedConceptID"], &recoveredReferencedConceptID)
	if err != nil {
		log.Printf("Recovery of Reference.ReferencedConceptID as string failed\n")
		return err
	}
	rPtr.ReferencedConceptID = recoveredReferencedConceptID
	rPtr.referencedConcept.setIndicatedConceptID(recoveredReferencedConceptID)
	rPtr.referencedConcept.parentConceptID = rPtr.getConceptIDNoLock()
	foundReferencedConcept := rPtr.uOfD.GetElement(recoveredReferencedConceptID)
	if foundReferencedConcept == nil {
		rPtr.uOfD.addUnresolvedPointer(rPtr.referencedConcept)
	} else {
		rPtr.referencedConcept.setIndicatedConcept(foundReferencedConcept)
	}
	// ReferencedAttributeName
	var recoveredReferencedConceptAttributeName string
	err = json.Unmarshal((*unmarshaledData)["ReferencedAttributeName"], &recoveredReferencedConceptAttributeName)
	if err != nil {
		log.Printf("Recovery of Reference.ReferencedAttributeName as string failed\n")
		return err
	}
	var attributeName AttributeName
	attributeName, err = FindAttributeName(recoveredReferencedConceptAttributeName)
	if err != nil {
		log.Printf("Conversion of Reference.ReferencedAttributeName to AttributeName failed\n")
		return err
	}
	rPtr.ReferencedAttributeName = attributeName
	// ReferencedConceptVersion
	var recoveredReferencedConceptVersion string
	err = json.Unmarshal((*unmarshaledData)["ReferencedConceptVersion"], &recoveredReferencedConceptVersion)
	if err != nil {
		log.Printf("Recovery of Reference.ReferencedConceptVersion failed\n")
		return err
	}
	rPtr.ReferencedConceptVersion, err = strconv.Atoi(recoveredReferencedConceptVersion)
	if err != nil {
		log.Printf("Conversion of Reference.ReferencedConceptVersion to integer failed\n")
		return err
	}
	return nil
}

// SetReferencedConcept sets the referenced concept by calling SetReferencedConceptID using the ID of the
// supplied Element
func (rPtr *reference) SetReferencedConcept(el Element, hl *HeldLocks) error {
	hl.WriteLockElement(rPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return rPtr.SetReferencedConceptID(id, hl)
}

// SetReferencedConceptID sets the referenced concept using the supplied ID.
func (rPtr *reference) SetReferencedConceptID(rcID string, hl *HeldLocks) error {
	hl.WriteLockElement(rPtr)
	editableError := rPtr.editableError(hl)
	if editableError != nil {
		return editableError
	}
	if rPtr.ReferencedConceptID != rcID {
		rPtr.uOfD.preChange(rPtr, hl)
		if rPtr.ReferencedConceptID != "" {
			oldReferencedConcept := rPtr.uOfD.GetElement(rPtr.ReferencedConceptID)
			if oldReferencedConcept != nil {
				oldReferencedConcept.removeListener(rPtr.ConceptID, hl)
			}
		}
		var newReferencedConcept Element
		if rcID != "" {
			newReferencedConcept = rPtr.uOfD.GetElement(rcID)
			if newReferencedConcept != nil {
				newReferencedConcept.addListener(rPtr.ConceptID, hl)
			}
		}
		notification := rPtr.uOfD.NewConceptChangeNotification(rPtr, hl)
		rPtr.ReferencedConceptID = rcID
		rPtr.referencedConcept.setIndicatedConcept(newReferencedConcept)
		rPtr.referencedConcept.setIndicatedConceptID(rcID)
		if newReferencedConcept == nil {
			rPtr.ReferencedConceptVersion = 0
		} else {
			rPtr.ReferencedConceptVersion = newReferencedConcept.GetVersion(hl)
		}
		rPtr.uOfD.queueFunctionExecutions(rPtr, notification, hl)
	}
	return nil
}

// SetReferencedConceptAttribute sets the value indicating whether a specific attribute of the referenced concept is being
// referenced
func (rPtr *reference) SetReferencedAttributeName(attributeName AttributeName, hl *HeldLocks) error {
	hl.WriteLockElement(rPtr)
	editableError := rPtr.editableError(hl)
	if editableError != nil {
		return editableError
	}
	if rPtr.ReferencedAttributeName != attributeName {
		rPtr.uOfD.preChange(rPtr, hl)
		notification := rPtr.uOfD.NewConceptChangeNotification(rPtr, hl)
		rPtr.ReferencedAttributeName = attributeName
		rPtr.uOfD.queueFunctionExecutions(rPtr, notification, hl)
	}
	return nil
}

// Reference represents a concept that is a pointer to another concept
type Reference interface {
	Element
	GetReferencedConcept(*HeldLocks) Element
	GetReferencedConceptID(*HeldLocks) string
	GetReferencedAttributeName(*HeldLocks) AttributeName
	GetReferencedConceptVersion(*HeldLocks) int
	getReferencedConceptNoLock() Element
	SetReferencedConcept(Element, *HeldLocks) error
	SetReferencedAttributeName(attributeName AttributeName, hl *HeldLocks) error
	SetReferencedConceptID(string, *HeldLocks) error
}
