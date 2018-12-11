package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type reference struct {
	element
	ReferencedConceptID string
	// referencedConcept is a cache for convenience
	referencedConcept        *cachedPointer
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
		rPtr.uOfD.queueFunctionExecutions(rPtr, notification, hl)
		rPtr.ReferencedConceptID = rcID
		rPtr.referencedConcept.setIndicatedConcept(newReferencedConcept)
		rPtr.referencedConcept.setIndicatedConceptID(rcID)
		if newReferencedConcept == nil {
			rPtr.ReferencedConceptVersion = 0
		} else {
			rPtr.ReferencedConceptVersion = newReferencedConcept.GetVersion(hl)
		}
	}
	return nil
}

// Reference represents a concept that is a pointer to another concept
type Reference interface {
	Element
	GetReferencedConcept(*HeldLocks) Element
	GetReferencedConceptID(*HeldLocks) string
	GetReferencedConceptVersion(*HeldLocks) int
	getReferencedConceptNoLock() Element
	SetReferencedConcept(Element, *HeldLocks) error
	SetReferencedConceptID(string, *HeldLocks) error
}
