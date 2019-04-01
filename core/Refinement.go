package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type refinement struct {
	element
	AbstractConceptID      string
	abstractConcept        *cachedPointer
	AbstractConceptVersion int
	RefinedConceptID       string
	refinedConcept         *cachedPointer
	RefinedConceptVersion  int
}

func (rPtr *refinement) clone(hl *HeldLocks) Refinement {
	hl.ReadLockElement(rPtr)
	var ref refinement
	ref.initializeRefinement("", "")
	ref.cloneAttributes(rPtr, hl)
	return &ref
}

func (rPtr *refinement) cloneAttributes(source *refinement, hl *HeldLocks) {
	rPtr.element.cloneAttributes(&source.element, hl)
	rPtr.AbstractConceptID = source.AbstractConceptID
	//	rPtr.abstractConcept = newCachedPointer(rPtr.getConceptIDNoLock(), false)
	rPtr.abstractConcept.setIndicatedConceptID(source.AbstractConceptID)
	rPtr.abstractConcept.setIndicatedConcept(source.abstractConcept.getIndicatedConcept())
	rPtr.abstractConcept.parentConceptID = source.abstractConcept.parentConceptID
	rPtr.AbstractConceptVersion = source.AbstractConceptVersion
	rPtr.RefinedConceptID = source.RefinedConceptID
	//	rPtr.refinedConcept = newCachedPointer(rPtr.getConceptIDNoLock(), false)
	rPtr.refinedConcept.setIndicatedConceptID(source.RefinedConceptID)
	rPtr.refinedConcept.setIndicatedConcept(source.refinedConcept.getIndicatedConcept())
	rPtr.refinedConcept.parentConceptID = source.refinedConcept.parentConceptID
	rPtr.RefinedConceptVersion = source.RefinedConceptVersion
}

func (rPtr *refinement) GetAbstractConcept(hl *HeldLocks) Element {
	hl.ReadLockElement(rPtr)
	cachedConcept := rPtr.abstractConcept.getIndicatedConcept()
	if cachedConcept == nil && rPtr.AbstractConceptID != "" {
		cachedConcept = rPtr.uOfD.GetElement(rPtr.AbstractConceptID)
		if cachedConcept != nil {
			rPtr.abstractConcept.setIndicatedConcept(cachedConcept)
		}
	}
	return cachedConcept
}

func (rPtr *refinement) getAbstractConceptNoLock() Element {
	return rPtr.abstractConcept.getIndicatedConcept()
}

func (rPtr *refinement) GetAbstractConceptID(hl *HeldLocks) string {
	hl.ReadLockElement(rPtr)
	return rPtr.AbstractConceptID
}

func (rPtr *refinement) getAbstractConceptIDNoLock() string {
	return rPtr.AbstractConceptID
}

func (rPtr *refinement) GetAbstractConceptVersion(hl *HeldLocks) int {
	hl.ReadLockElement(rPtr)
	return rPtr.AbstractConceptVersion
}

func (rPtr *refinement) GetRefinedConcept(hl *HeldLocks) Element {
	hl.ReadLockElement(rPtr)
	cachedConcept := rPtr.refinedConcept.getIndicatedConcept()
	if cachedConcept == nil && rPtr.RefinedConceptID != "" {
		cachedConcept = rPtr.uOfD.GetElement(rPtr.RefinedConceptID)
		if cachedConcept != nil {
			rPtr.refinedConcept.setIndicatedConcept(cachedConcept)
		}
	}
	return cachedConcept
}

func (rPtr *refinement) getRefinedConceptNoLock() Element {
	return rPtr.refinedConcept.getIndicatedConcept()
}

func (rPtr *refinement) GetRefinedConceptID(hl *HeldLocks) string {
	hl.ReadLockElement(rPtr)
	return rPtr.RefinedConceptID
}

func (rPtr *refinement) getRefinedConceptIDNoLock() string {
	return rPtr.RefinedConceptID
}

func (rPtr *refinement) GetRefinedConceptVersion(hl *HeldLocks) int {
	hl.ReadLockElement(rPtr)
	return rPtr.RefinedConceptVersion
}

func (rPtr *refinement) initializeRefinement(conceptID string, uri string) {
	rPtr.initializeElement(conceptID, uri)
	rPtr.abstractConcept = newCachedPointer(rPtr.getConceptIDNoLock(), false)
	rPtr.refinedConcept = newCachedPointer(rPtr.getConceptIDNoLock(), false)
}

func (rPtr *refinement) isEquivalent(hl1 *HeldLocks, ref *refinement, hl2 *HeldLocks) bool {
	hl1.ReadLockElement(rPtr)
	hl2.ReadLockElement(ref)
	if rPtr.AbstractConceptID != ref.AbstractConceptID {
		return false
	}
	if rPtr.abstractConcept.isEquivalent(ref.abstractConcept) != true {
		return false
	}
	if rPtr.AbstractConceptVersion != ref.AbstractConceptVersion {
		return false
	}
	if rPtr.RefinedConceptID != ref.RefinedConceptID {
		return false
	}
	if rPtr.refinedConcept.isEquivalent(ref.refinedConcept) != true {
		return false
	}
	if rPtr.RefinedConceptVersion != ref.RefinedConceptVersion {
		return false
	}
	return rPtr.element.isEquivalent(hl1, &ref.element, hl2)
}

// MarshalJSON produces a byte string JSON representation of the Element
func (rPtr *refinement) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(rPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := rPtr.marshalRefinementFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (rPtr *refinement) marshalRefinementFields(buffer *bytes.Buffer) error {
	buffer.WriteString(fmt.Sprintf("\"AbstractConceptID\":\"%s\",", rPtr.AbstractConceptID))
	buffer.WriteString(fmt.Sprintf("\"AbstractConceptVersion\":\"%d\",", rPtr.AbstractConceptVersion))
	buffer.WriteString(fmt.Sprintf("\"RefinedConceptID\":\"%s\",", rPtr.RefinedConceptID))
	buffer.WriteString(fmt.Sprintf("\"RefinedConceptVersion\":\"%d\",", rPtr.RefinedConceptVersion))
	rPtr.element.marshalElementFields(buffer)
	return nil
}

// recoverRefinementFields() is used when de-serializing an element. The activities in restoring the
// element are not considered changes so the version counter is not incremented and the monitors of this
// element are not notified of chaanges.
func (rPtr *refinement) recoverRefinementFields(unmarshaledData *map[string]json.RawMessage, hl *HeldLocks) error {
	err := rPtr.recoverElementFields(unmarshaledData, hl)
	if err != nil {
		return err
	}
	// AbstractConceptID
	var recoveredAbstractConceptID string
	err = json.Unmarshal((*unmarshaledData)["AbstractConceptID"], &recoveredAbstractConceptID)
	if err != nil {
		log.Printf("Recovery of Refinement.AbstractConceptID as string failed\n")
		return err
	}
	rPtr.AbstractConceptID = recoveredAbstractConceptID
	rPtr.abstractConcept.setIndicatedConceptID(recoveredAbstractConceptID)
	rPtr.abstractConcept.parentConceptID = rPtr.getConceptIDNoLock()
	foundAbstractConcept := rPtr.uOfD.GetElement(recoveredAbstractConceptID)
	if foundAbstractConcept == nil {
		rPtr.uOfD.addUnresolvedPointer(rPtr.abstractConcept)
	} else {
		rPtr.abstractConcept.setIndicatedConcept(foundAbstractConcept)
	}
	// AbstractConceptVersion
	var recoveredAbstractConceptVersion string
	err = json.Unmarshal((*unmarshaledData)["AbstractConceptVersion"], &recoveredAbstractConceptVersion)
	if err != nil {
		log.Printf("Recovery of Refinement.AbstractConceptVersion failed\n")
		return err
	}
	rPtr.AbstractConceptVersion, err = strconv.Atoi(recoveredAbstractConceptVersion)
	if err != nil {
		log.Printf("Conversion of Refinement.AbstractConceptVersion to integer failed\n")
		return err
	}
	// RefinedConceptID
	var recoveredRefinedConceptID string
	err = json.Unmarshal((*unmarshaledData)["RefinedConceptID"], &recoveredRefinedConceptID)
	if err != nil {
		log.Printf("Recovery of Refinement.RefinedConceptID as string failed\n")
		return err
	}
	rPtr.RefinedConceptID = recoveredRefinedConceptID
	rPtr.refinedConcept.setIndicatedConceptID(recoveredRefinedConceptID)
	rPtr.refinedConcept.parentConceptID = rPtr.getConceptIDNoLock()
	foundRefinedConcept := rPtr.uOfD.GetElement(recoveredRefinedConceptID)
	if foundRefinedConcept == nil {
		rPtr.uOfD.addUnresolvedPointer(rPtr.refinedConcept)
	} else {
		rPtr.refinedConcept.setIndicatedConcept(foundRefinedConcept)
	}
	// RefinedConceptVersion
	var recoveredRefinedConceptVersion string
	err = json.Unmarshal((*unmarshaledData)["RefinedConceptVersion"], &recoveredRefinedConceptVersion)
	if err != nil {
		log.Printf("Recovery of Refinement.RefinedConceptVersion failed\n")
		return err
	}
	rPtr.RefinedConceptVersion, err = strconv.Atoi(recoveredRefinedConceptVersion)
	if err != nil {
		log.Printf("Conversion of Refinement.RefinedConceptVersion to integer failed\n")
		return err
	}
	return nil
}

// SetAbstractConcept sets the abstract concept using the ID of the supplied Element
func (rPtr *refinement) SetAbstractConcept(el Element, hl *HeldLocks) error {
	hl.WriteLockElement(rPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return rPtr.SetAbstractConceptID(id, hl)
}

func (rPtr *refinement) SetAbstractConceptID(acID string, hl *HeldLocks) error {
	hl.WriteLockElement(rPtr)
	editableError := rPtr.editableError(hl)
	if editableError != nil {
		return editableError
	}
	if rPtr.AbstractConceptID != acID {
		rPtr.uOfD.preChange(rPtr, hl)
		if rPtr.AbstractConceptID != "" {
			rPtr.uOfD.GetElement(rPtr.AbstractConceptID).removeListener(rPtr.ConceptID, hl)
		}
		var newAbstractConcept Element
		if acID != "" {
			newAbstractConcept = rPtr.uOfD.GetElement(acID)
			if newAbstractConcept != nil {
				newAbstractConcept.addListener(rPtr.ConceptID, hl)
			}
		}
		notification := rPtr.uOfD.NewConceptChangeNotification(rPtr, hl)
		rPtr.AbstractConceptID = acID
		rPtr.abstractConcept.setIndicatedConcept(newAbstractConcept)
		rPtr.abstractConcept.setIndicatedConceptID(acID)
		if newAbstractConcept != nil {
			rPtr.AbstractConceptVersion = newAbstractConcept.GetVersion(hl)
		} else {
			rPtr.AbstractConceptVersion = 0
		}
		rPtr.uOfD.queueFunctionExecutions(rPtr, notification, hl)
	}
	return nil
}

func (rPtr *refinement) SetRefinedConcept(el Element, hl *HeldLocks) error {
	hl.WriteLockElement(rPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return rPtr.SetRefinedConceptID(id, hl)
}

func (rPtr *refinement) SetRefinedConceptID(rcID string, hl *HeldLocks) error {
	hl.WriteLockElement(rPtr)
	editableError := rPtr.editableError(hl)
	if editableError != nil {
		return editableError
	}
	if rPtr.RefinedConceptID != rcID {
		rPtr.uOfD.preChange(rPtr, hl)
		if rPtr.RefinedConceptID != "" {
			rPtr.uOfD.GetElement(rPtr.RefinedConceptID).removeListener(rPtr.ConceptID, hl)
		}
		var newRefinedConcept Element
		if rcID != "" {
			newRefinedConcept = rPtr.uOfD.GetElement(rcID)
			if newRefinedConcept != nil {
				newRefinedConcept.addListener(rPtr.ConceptID, hl)
			}
		}
		notification := rPtr.uOfD.NewConceptChangeNotification(rPtr, hl)
		rPtr.RefinedConceptID = rcID
		rPtr.refinedConcept.setIndicatedConcept(newRefinedConcept)
		rPtr.refinedConcept.setIndicatedConceptID(rcID)
		if newRefinedConcept != nil {
			rPtr.RefinedConceptVersion = newRefinedConcept.GetVersion(hl)
		} else {
			rPtr.RefinedConceptVersion = 0
		}
		rPtr.uOfD.queueFunctionExecutions(rPtr, notification, hl)
	}
	return nil
}

// Refinement is the reification of a refinement association between an abstract Element and a refined Element
type Refinement interface {
	Element
	GetAbstractConceptID(*HeldLocks) string
	getAbstractConceptIDNoLock() string
	GetAbstractConcept(*HeldLocks) Element
	getAbstractConceptNoLock() Element
	GetAbstractConceptVersion(*HeldLocks) int
	GetRefinedConceptID(*HeldLocks) string
	getRefinedConceptIDNoLock() string
	GetRefinedConcept(*HeldLocks) Element
	getRefinedConceptNoLock() Element
	GetRefinedConceptVersion(*HeldLocks) int
	SetAbstractConcept(Element, *HeldLocks) error
	SetAbstractConceptID(string, *HeldLocks) error
	SetRefinedConcept(Element, *HeldLocks) error
	SetRefinedConceptID(string, *HeldLocks) error
}
