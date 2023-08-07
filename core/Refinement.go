package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/pkg/errors"
)

type refinement struct {
	element
	AbstractConceptID string
	RefinedConceptID  string
}

func (rPtr *refinement) clone(trans *Transaction) Refinement {
	trans.ReadLockElement(rPtr)
	var ref refinement
	ref.initializeRefinement("", "")
	ref.cloneAttributes(rPtr, trans)
	return &ref
}

func (rPtr *refinement) cloneAttributes(source *refinement, trans *Transaction) {
	rPtr.element.cloneAttributes(&source.element, trans)
	rPtr.AbstractConceptID = source.AbstractConceptID
	rPtr.RefinedConceptID = source.RefinedConceptID
}

func (rPtr *refinement) GetAbstractConcept(trans *Transaction) Element {
	trans.ReadLockElement(rPtr)
	return rPtr.uOfD.GetElement(rPtr.AbstractConceptID)
}

func (rPtr *refinement) getAbstractConceptNoLock() Element {
	return rPtr.uOfD.GetElement(rPtr.AbstractConceptID)
}

func (rPtr *refinement) GetAbstractConceptID(trans *Transaction) string {
	trans.ReadLockElement(rPtr)
	return rPtr.AbstractConceptID
}

func (rPtr *refinement) getAbstractConceptIDNoLock() string {
	return rPtr.AbstractConceptID
}

func (rPtr *refinement) GetRefinedConcept(trans *Transaction) Element {
	trans.ReadLockElement(rPtr)
	return rPtr.uOfD.GetElement(rPtr.RefinedConceptID)
}

func (rPtr *refinement) getRefinedConceptNoLock() Element {
	return rPtr.uOfD.GetElement(rPtr.RefinedConceptID)
}

func (rPtr *refinement) GetRefinedConceptID(trans *Transaction) string {
	trans.ReadLockElement(rPtr)
	return rPtr.RefinedConceptID
}

func (rPtr *refinement) getRefinedConceptIDNoLock() string {
	return rPtr.RefinedConceptID
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
	if rPtr.RefinedConceptID != ref.RefinedConceptID {
		if print {
			log.Printf("In refinement.isEquivalent, RefinedConecptIDs do not match")
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
	buffer.WriteString(fmt.Sprintf("\"RefinedConceptID\":\"%s\",", rPtr.RefinedConceptID))
	rPtr.element.marshalElementFields(buffer)
	return nil
}

// recoverRefinementFields() is used when de-serializing an element. The activities in restoring the
// element are not considered changes so the version counter is not incremented and the monitors of this
// element are not notified of chaanges.
func (rPtr *refinement) recoverRefinementFields(unmarshaledData *map[string]json.RawMessage, trans *Transaction) error {
	err := rPtr.recoverElementFields(unmarshaledData, trans)
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
	// RefinedConceptID
	var recoveredRefinedConceptID string
	err = json.Unmarshal((*unmarshaledData)["RefinedConceptID"], &recoveredRefinedConceptID)
	if err != nil {
		log.Printf("Recovery of Refinement.RefinedConceptID as string failed\n")
		return err
	}
	rPtr.RefinedConceptID = recoveredRefinedConceptID
	return nil
}

// SetAbstractConcept sets the abstract concept using the ID of the supplied Element
func (rPtr *refinement) SetAbstractConcept(el Element, trans *Transaction) error {
	if rPtr.uOfD == nil {
		return errors.New("refinement.SetAbstractConcept failed because the element uOfD is nil")
	}
	trans.WriteLockElement(rPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return rPtr.SetAbstractConceptID(id, trans)
}

func (rPtr *refinement) SetAbstractConceptID(acID string, trans *Transaction) error {
	if rPtr.uOfD == nil {
		return errors.New("refinement.SetAbstractConceptID failed because the element uOfD is nil")
	}
	trans.WriteLockElement(rPtr)
	if !rPtr.isEditable(trans) {
		return errors.New("refinement.SetAbstractConceptID failed because the refinement is not editable")
	}
	if rPtr.AbstractConceptID != acID {
		rPtr.uOfD.preChange(rPtr, trans)
		beforeState, err := NewConceptState(rPtr)
		if err != nil {
			return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
		}
		rPtr.incrementVersion(trans)
		var oldAbstractConcept Element
		if rPtr.AbstractConceptID != "" {
			oldAbstractConcept = rPtr.uOfD.GetElement(rPtr.AbstractConceptID)
			if oldAbstractConcept != nil {
				oldAbstractConcept.removeListener(rPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
				}
			} else {
				// This case can arise if the abstract concept is not currently loaded
				rPtr.uOfD.listenersMap.removeMappedValue(rPtr.AbstractConceptID, rPtr.ConceptID)
			}
		}
		var newAbstractConcept Element
		if acID != "" {
			newAbstractConcept = rPtr.uOfD.GetElement(acID)
			if newAbstractConcept != nil {
				newAbstractConcept.addListener(rPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
				}
			}
		}
		rPtr.AbstractConceptID = acID
		afterState, err2 := NewConceptState(rPtr)
		if err2 != nil {
			return errors.Wrap(err2, "refinement.SetAbstractConceptID failed")
		}
		err = rPtr.uOfD.SendPointerChangeNotification(rPtr, AbstractConceptChanged, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
		}
	}
	return nil
}

func (rPtr *refinement) SetRefinedConcept(el Element, trans *Transaction) error {
	if rPtr.uOfD == nil {
		return errors.New("refinement.SetRefinedConcept failed because the element uOfD is nil")
	}
	trans.WriteLockElement(rPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return rPtr.SetRefinedConceptID(id, trans)
}

func (rPtr *refinement) SetRefinedConceptID(rcID string, trans *Transaction) error {
	if rPtr.uOfD == nil {
		return errors.New("refinement.SetRefinedConceptID failed because the element uOfD is nil")
	}
	trans.WriteLockElement(rPtr)
	if !rPtr.isEditable(trans) {
		return errors.New("refinement.SetReferencedConceptID failed because the refinement is not editable")
	}
	if rPtr.RefinedConceptID != rcID {
		rPtr.uOfD.preChange(rPtr, trans)
		beforeState, err := NewConceptState(rPtr)
		if err != nil {
			return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
		}
		rPtr.incrementVersion(trans)
		var oldRefinedConcept Element
		if rPtr.RefinedConceptID != "" {
			oldRefinedConcept = rPtr.uOfD.GetElement(rPtr.RefinedConceptID)
			if oldRefinedConcept != nil {
				oldRefinedConcept.removeListener(rPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
				}
			}
		}
		var newRefinedConcept Element
		if rcID != "" {
			newRefinedConcept = rPtr.uOfD.GetElement(rcID)
			if newRefinedConcept != nil {
				newRefinedConcept.addListener(rPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
				}
			}
		}
		rPtr.RefinedConceptID = rcID
		afterState, err2 := NewConceptState(rPtr)
		if err2 != nil {
			return errors.Wrap(err2, "refinement.SetRefinedConceptID failed")
		}
		err = rPtr.uOfD.SendPointerChangeNotification(rPtr, RefinedConceptChanged, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
		}
	}
	return nil
}

// Refinement is the reification of a refinement association between an abstract Element and a refined Element.
type Refinement interface {
	Element
	GetAbstractConceptID(*Transaction) string
	getAbstractConceptIDNoLock() string
	GetAbstractConcept(*Transaction) Element
	getAbstractConceptNoLock() Element
	GetRefinedConceptID(*Transaction) string
	getRefinedConceptIDNoLock() string
	GetRefinedConcept(*Transaction) Element
	getRefinedConceptNoLock() Element
	SetAbstractConcept(Element, *Transaction) error
	SetAbstractConceptID(string, *Transaction) error
	SetRefinedConcept(Element, *Transaction) error
	SetRefinedConceptID(string, *Transaction) error
}
