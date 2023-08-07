package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/pkg/errors"
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
	Label               = AttributeName(6)
	Definition          = AttributeName(7)
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
	case Label:
		return "Label"
	case Definition:
		return "Definition"
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
	ReferencedConceptID     string
	ReferencedAttributeName AttributeName
}

func (rPtr *reference) clone(trans *Transaction) Reference {
	trans.ReadLockElement(rPtr)
	var ref reference
	ref.initializeReference("", "")
	ref.cloneAttributes(rPtr, trans)
	return &ref
}

func (rPtr *reference) cloneAttributes(source *reference, trans *Transaction) {
	rPtr.element.cloneAttributes(&source.element, trans)
	rPtr.ReferencedConceptID = source.ReferencedConceptID
}

// GetReferencedConcept returns the element representing  the concept being referenced
// Note that this is a cached value
func (rPtr *reference) GetReferencedConcept(trans *Transaction) Element {
	trans.ReadLockElement(rPtr)
	return rPtr.getReferencedConceptNoLock()
}

func (rPtr *reference) getReferencedConceptNoLock() Element {
	return rPtr.uOfD.GetElement(rPtr.ReferencedConceptID)
}

// GetReferencedConceptID returns the identifier of the concept being referenced
func (rPtr *reference) GetReferencedConceptID(trans *Transaction) string {
	trans.ReadLockElement(rPtr)
	return rPtr.ReferencedConceptID
}

// GetReferencedAttributeName returns an indicator of which attribute is being referenced (if any)
func (rPtr *reference) GetReferencedAttributeName(trans *Transaction) AttributeName {
	trans.ReadLockElement(rPtr)
	return rPtr.ReferencedAttributeName
}

// GetReferencedAttributeValue returns the string value of the referenced attribute (if any)
func (rPtr *reference) GetReferencedAttributeValue(trans *Transaction) string {
	trans.ReadLockElement(rPtr)
	if rPtr.ReferencedConceptID != "" {
		referencedConcept := rPtr.GetReferencedConcept(trans)
		if referencedConcept != nil {
			if rPtr.ReferencedAttributeName == OwningConceptID {
				return referencedConcept.GetOwningConceptID(trans)
			}
			switch typedReferencedConcept := referencedConcept.(type) {
			case Reference:
				if rPtr.ReferencedAttributeName == ReferencedConceptID {
					return typedReferencedConcept.GetReferencedConceptID(trans)
				}
			case Refinement:
				if rPtr.ReferencedAttributeName == AbstractConceptID {
					return typedReferencedConcept.GetAbstractConceptID(trans)
				}
				if rPtr.ReferencedAttributeName == RefinedConceptID {
					return typedReferencedConcept.GetRefinedConceptID(trans)
				}
			case Literal:
				if rPtr.ReferencedAttributeName == LiteralValue {
					return typedReferencedConcept.GetLiteralValue(trans)
				}
			}
		}
	}
	return ""
}

func (rPtr *reference) initializeReference(conceptID string, uri string) {
	rPtr.initializeElement(conceptID, uri)
}

func (rPtr *reference) isEquivalent(hl1 *Transaction, el *reference, hl2 *Transaction, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) > 0 {
		print = printExceptions[0]
	}
	hl1.ReadLockElement(rPtr)
	hl2.ReadLockElement(el)
	if rPtr.ReferencedConceptID != el.ReferencedConceptID {
		if print {
			log.Printf("In reference.IsEquivalent, ReferencedConceptIDs do not match")
		}
		return false
	}
	if rPtr.ReferencedAttributeName != el.ReferencedAttributeName {
		if print {
			log.Printf("In reference.IsEquivalent, ReferencedAttributeNames do not match")
		}
		return false
	}
	return rPtr.element.isEquivalent(hl1, &el.element, hl2, print)
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
	rPtr.element.marshalElementFields(buffer)
	return nil
}

// recoverReferenceFields() is used when de-serializing an element. The activities in restoring the
// element are not considered changes so the version counter is not incremented and the monitors of this
// element are not notified of chaanges.
func (rPtr *reference) recoverReferenceFields(unmarshaledData *map[string]json.RawMessage, trans *Transaction) error {
	err := rPtr.recoverElementFields(unmarshaledData, trans)
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
	return nil
}

// SetReferencedConcept sets the referenced concept by calling SetReferencedConceptID using the ID of the
// supplied Element
func (rPtr *reference) SetReferencedConcept(el Element, attributeName AttributeName, trans *Transaction) error {
	if rPtr.uOfD == nil {
		return errors.New("reference.SetReferencedConcept failed because the element uOfD is nil")
	}
	trans.WriteLockElement(rPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return rPtr.SetReferencedConceptID(id, attributeName, trans)
}

// SetReferencedConceptID sets the referenced concept using the supplied ID.
func (rPtr *reference) SetReferencedConceptID(rcID string, attributeName AttributeName, trans *Transaction) error {
	if rPtr.uOfD == nil {
		return errors.New("reference.SetReferencedConceptID failed because the element uOfD is nil")
	}
	trans.WriteLockElement(rPtr)
	if !rPtr.isEditable(trans) {
		return errors.New("reference.SetReferencedConceptID failed because the reference is not editable")
	}
	var newReferencedConcept Element
	var oldReferencedConcept Element
	if rPtr.ReferencedConceptID != rcID || rPtr.ReferencedAttributeName != attributeName {
		if rcID != "" {
			newReferencedConcept = rPtr.uOfD.GetElement(rcID)
			switch attributeName {
			case ReferencedConceptID:
				switch newReferencedConcept.(type) {
				case Reference:
				default:
					return errors.New("In reference.SetReferencedConceptID, the ReferencedAttributeName was ReferencedConceptID, but the referenced concept is not a Reference")
				}
			case AbstractConceptID, RefinedConceptID:
				switch newReferencedConcept.(type) {
				case Refinement:
				default:
					return errors.New("In reference.SetReferencedConceptID, the ReferencedAttributeName was AbstractConceptID or RefinedConceptID, but the referenced concept is not a Refinement")
				}
			}
			if newReferencedConcept != nil {
				newReferencedConcept.addListener(rPtr.ConceptID, trans)
			}
		}
		beforeState, err := NewConceptState(rPtr)
		if err != nil {
			return errors.Wrap(err, "reference.SetReferencedConceptID failed")
		}
		rPtr.uOfD.preChange(rPtr, trans)
		rPtr.incrementVersion(trans)
		if rPtr.ReferencedConceptID != "" {
			oldReferencedConcept = rPtr.uOfD.GetElement(rPtr.ReferencedConceptID)
			if oldReferencedConcept != nil {
				oldReferencedConcept.removeListener(rPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "reference.SetReferencedConceptID failed")
				}
			}
		}
		if rcID != "" {
			if newReferencedConcept != nil {
				newReferencedConcept.addListener(rPtr.ConceptID, trans)
			}
		}
		rPtr.ReferencedConceptID = rcID
		rPtr.ReferencedAttributeName = attributeName
		afterState, err2 := NewConceptState(rPtr)
		if err2 != nil {
			return errors.Wrap(err2, "reference.SetReferencedConceptID failed")
		}
		err = rPtr.uOfD.SendPointerChangeNotification(rPtr, ReferencedConceptChanged, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "reference.SetReferencedConceptID failed")
		}
	}
	return nil
}

// Reference represents a concept that is a pointer to another concept.
type Reference interface {
	Element
	GetReferencedConcept(*Transaction) Element
	GetReferencedConceptID(*Transaction) string
	GetReferencedAttributeName(*Transaction) AttributeName
	GetReferencedAttributeValue(*Transaction) string
	getReferencedConceptNoLock() Element
	SetReferencedConcept(Element, AttributeName, *Transaction) error
	SetReferencedConceptID(string, AttributeName, *Transaction) error
}
