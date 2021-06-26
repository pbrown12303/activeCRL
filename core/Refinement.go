package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

type refinement struct {
	element
	AbstractConceptID      string
	AbstractConceptVersion int
	RefinedConceptID       string
	RefinedConceptVersion  int
}

func (rPtr *refinement) clone(hl *Transaction) Refinement {
	hl.ReadLockElement(rPtr)
	var ref refinement
	ref.initializeRefinement("", "")
	ref.cloneAttributes(rPtr, hl)
	return &ref
}

func (rPtr *refinement) cloneAttributes(source *refinement, hl *Transaction) {
	rPtr.element.cloneAttributes(&source.element, hl)
	rPtr.AbstractConceptID = source.AbstractConceptID
	//	rPtr.abstractConcept = newCachedPointer(rPtr.getConceptIDNoLock(), false)
	rPtr.AbstractConceptVersion = source.AbstractConceptVersion
	rPtr.RefinedConceptID = source.RefinedConceptID
	//	rPtr.refinedConcept = newCachedPointer(rPtr.getConceptIDNoLock(), false)
	rPtr.RefinedConceptVersion = source.RefinedConceptVersion
}

func (rPtr *refinement) GetAbstractConcept(hl *Transaction) Element {
	hl.ReadLockElement(rPtr)
	return rPtr.uOfD.GetElement(rPtr.AbstractConceptID)
}

func (rPtr *refinement) getAbstractConceptNoLock() Element {
	return rPtr.uOfD.GetElement(rPtr.AbstractConceptID)
}

func (rPtr *refinement) GetAbstractConceptID(hl *Transaction) string {
	hl.ReadLockElement(rPtr)
	return rPtr.AbstractConceptID
}

func (rPtr *refinement) getAbstractConceptIDNoLock() string {
	return rPtr.AbstractConceptID
}

func (rPtr *refinement) GetAbstractConceptVersion(hl *Transaction) int {
	hl.ReadLockElement(rPtr)
	return rPtr.AbstractConceptVersion
}

func (rPtr *refinement) GetRefinedConcept(hl *Transaction) Element {
	hl.ReadLockElement(rPtr)
	return rPtr.uOfD.GetElement(rPtr.RefinedConceptID)
}

func (rPtr *refinement) getRefinedConceptNoLock() Element {
	return rPtr.uOfD.GetElement(rPtr.RefinedConceptID)
}

func (rPtr *refinement) GetRefinedConceptID(hl *Transaction) string {
	hl.ReadLockElement(rPtr)
	return rPtr.RefinedConceptID
}

func (rPtr *refinement) getRefinedConceptIDNoLock() string {
	return rPtr.RefinedConceptID
}

func (rPtr *refinement) GetRefinedConceptVersion(hl *Transaction) int {
	hl.ReadLockElement(rPtr)
	return rPtr.RefinedConceptVersion
}

func (rPtr *refinement) initializeRefinement(conceptID string, uri string) {
	rPtr.initializeElement(conceptID, uri)
}

func (rPtr *refinement) isEquivalent(hl1 *Transaction, ref *refinement, hl2 *Transaction, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) > 0 {
		print = printExceptions[0]
	}
	hl1.ReadLockElement(rPtr)
	hl2.ReadLockElement(ref)
	if rPtr.AbstractConceptID != ref.AbstractConceptID {
		if print {
			log.Printf("In refinement.isEquivalent, AbstractConecptIDs do not match")
		}
		return false
	}
	if rPtr.AbstractConceptVersion != ref.AbstractConceptVersion {
		if print {
			log.Printf("In refinement.isEquivalent, AbstractConecptVersionss do not match")
		}
		return false
	}
	if rPtr.RefinedConceptID != ref.RefinedConceptID {
		if print {
			log.Printf("In refinement.isEquivalent, RefinedConecptIDs do not match")
		}
		return false
	}
	if rPtr.RefinedConceptVersion != ref.RefinedConceptVersion {
		if print {
			log.Printf("In refinement.isEquivalent, RefinedConecptVersions do not match")
		}
		return false
	}
	return rPtr.element.isEquivalent(hl1, &ref.element, hl2, print)
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
func (rPtr *refinement) recoverRefinementFields(unmarshaledData *map[string]json.RawMessage, hl *Transaction) error {
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
func (rPtr *refinement) SetAbstractConcept(el Element, hl *Transaction) error {
	hl.WriteLockElement(rPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return rPtr.SetAbstractConceptID(id, hl)
}

func (rPtr *refinement) SetAbstractConceptID(acID string, hl *Transaction) error {
	hl.WriteLockElement(rPtr)
	if !rPtr.isEditable(hl) {
		return errors.New("refinement.SetAbstractConceptID failed because the refinement is not editable")
	}
	if rPtr.AbstractConceptID != acID {
		rPtr.uOfD.preChange(rPtr, hl)
		beforeState, err := NewConceptState(rPtr)
		if err != nil {
			return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
		}
		rPtr.incrementVersion(hl)
		var beforeReferencedState *ConceptState
		var oldAbstractConcept Element
		if rPtr.AbstractConceptID != "" {
			oldAbstractConcept = rPtr.uOfD.GetElement(rPtr.AbstractConceptID)
			if oldAbstractConcept != nil {
				oldAbstractConcept.removeListener(rPtr.ConceptID, hl)
				beforeReferencedState, err = NewConceptState(oldAbstractConcept)
				if err != nil {
					return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
				}
			} else {
				// This case can arise if the abstract concept is not currently loaded
				rPtr.uOfD.listenersMap.RemoveMappedValue(rPtr.AbstractConceptID, rPtr.ConceptID)
			}
		}
		var newAbstractConcept Element
		var afterReferencedState *ConceptState
		if acID != "" {
			newAbstractConcept = rPtr.uOfD.GetElement(acID)
			if newAbstractConcept != nil {
				newAbstractConcept.addListener(rPtr.ConceptID, hl)
				afterReferencedState, err = NewConceptState(newAbstractConcept)
				if err != nil {
					return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
				}
			}
		}
		rPtr.AbstractConceptID = acID
		if newAbstractConcept != nil {
			rPtr.AbstractConceptVersion = newAbstractConcept.GetVersion(hl)
		} else {
			rPtr.AbstractConceptVersion = 0
		}
		afterState, err2 := NewConceptState(rPtr)
		if err2 != nil {
			return errors.Wrap(err2, "refinement.SetAbstractConceptID failed")
		}
		err = rPtr.uOfD.SendPointerChangeNotification(rPtr, AbstractConceptChanged, beforeState, afterState, beforeReferencedState, afterReferencedState, hl)
		if err != nil {
			return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
		}
	}
	return nil
}

func (rPtr *refinement) SetRefinedConcept(el Element, hl *Transaction) error {
	hl.WriteLockElement(rPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return rPtr.SetRefinedConceptID(id, hl)
}

func (rPtr *refinement) SetRefinedConceptID(rcID string, hl *Transaction) error {
	hl.WriteLockElement(rPtr)
	if !rPtr.isEditable(hl) {
		return errors.New("refinement.SetReferencedConceptID failed because the refinement is not editable")
	}
	if rPtr.RefinedConceptID != rcID {
		rPtr.uOfD.preChange(rPtr, hl)
		beforeState, err := NewConceptState(rPtr)
		if err != nil {
			return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
		}
		rPtr.incrementVersion(hl)
		var beforeReferencedState *ConceptState
		var oldRefinedConcept Element
		if rPtr.RefinedConceptID != "" {
			oldRefinedConcept = rPtr.uOfD.GetElement(rPtr.RefinedConceptID)
			if oldRefinedConcept != nil {
				oldRefinedConcept.removeListener(rPtr.ConceptID, hl)
				beforeReferencedState, err = NewConceptState(oldRefinedConcept)
				if err != nil {
					return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
				}
			}
		}
		var newRefinedConcept Element
		var afterReferencedState *ConceptState
		if rcID != "" {
			newRefinedConcept = rPtr.uOfD.GetElement(rcID)
			if newRefinedConcept != nil {
				newRefinedConcept.addListener(rPtr.ConceptID, hl)
				afterReferencedState, err = NewConceptState(newRefinedConcept)
				if err != nil {
					return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
				}
			}
		}
		rPtr.RefinedConceptID = rcID
		if newRefinedConcept != nil {
			rPtr.RefinedConceptVersion = newRefinedConcept.GetVersion(hl)
		} else {
			rPtr.RefinedConceptVersion = 0
		}
		afterState, err2 := NewConceptState(rPtr)
		if err2 != nil {
			return errors.Wrap(err2, "refinement.SetRefinedConceptID failed")
		}
		err = rPtr.uOfD.SendPointerChangeNotification(rPtr, RefinedConceptChanged, beforeState, afterState, beforeReferencedState, afterReferencedState, hl)
		if err != nil {
			return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
		}
	}
	return nil
}

// Refinement is the reification of a refinement association between an abstract Element and a refined Element
type Refinement interface {
	Element
	GetAbstractConceptID(*Transaction) string
	getAbstractConceptIDNoLock() string
	GetAbstractConcept(*Transaction) Element
	getAbstractConceptNoLock() Element
	GetAbstractConceptVersion(*Transaction) int
	GetRefinedConceptID(*Transaction) string
	getRefinedConceptIDNoLock() string
	GetRefinedConcept(*Transaction) Element
	getRefinedConceptNoLock() Element
	GetRefinedConceptVersion(*Transaction) int
	SetAbstractConcept(Element, *Transaction) error
	SetAbstractConceptID(string, *Transaction) error
	SetRefinedConcept(Element, *Transaction) error
	SetRefinedConceptID(string, *Transaction) error
}
