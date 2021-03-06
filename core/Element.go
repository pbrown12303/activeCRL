package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
	"sync"

	mapset "github.com/deckarep/golang-set"
)

// element is the root representation of a concept
type element struct {
	sync.RWMutex
	ConceptID       string
	Definition      string
	Label           string
	IsCore          bool
	OwningConceptID string
	ReadOnly        bool
	Version         *versionCounter
	uOfD            *UniverseOfDiscourse
	URI             string
	observers       mapset.Set
}

// addOwnedConcept adds the indicated Element as a child (owned) concept.
// This is purely an internal housekeeping method. Note that
// no checking of whether the Element is read-only is performed here. This check
// is performed by the child
func (ePtr *element) addOwnedConcept(ownedConceptID string, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	if ePtr.uOfD.ownedIDsMap.ContainsMappedValue(ePtr.ConceptID, ownedConceptID) == false {
		ePtr.uOfD.preChange(ePtr, hl)
		ePtr.incrementVersion(hl)
		ePtr.uOfD.ownedIDsMap.AddMappedValue(ePtr.GetConceptID(hl), ownedConceptID)
	}
}

// addRecoveredOwnedConcept adds the indicated Element as a child (owned) concept without incrementing
// the version.
// This is purely an internal housekeeping method. Note that
// no checking of whether the Element is read-only is performed here. This check
// is performed by the child
func (ePtr *element) addRecoveredOwnedConcept(ownedConceptID string, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	if ePtr.uOfD.ownedIDsMap.ContainsMappedValue(ePtr.ConceptID, ownedConceptID) == false {
		ePtr.uOfD.preChange(ePtr, hl)
		ePtr.uOfD.ownedIDsMap.AddMappedValue(ePtr.ConceptID, ownedConceptID)
	}
}

// addListener adds the indicated Element as a listening concept.
// This is an internal housekeeping method.
func (ePtr *element) addListener(listeningConceptID string, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	if ePtr.uOfD.listenersMap.ContainsMappedValue(ePtr.ConceptID, listeningConceptID) == false {
		ePtr.uOfD.preChange(ePtr, hl)
		ePtr.uOfD.listenersMap.AddMappedValue(ePtr.ConceptID, listeningConceptID)
	}
}

// clone is an internal function that makes a copy of the given element - including its
// identifier. This is done only to support undo/redo: the clone should NEVER be added to the
// universe of discourse
func (ePtr *element) clone(hl *HeldLocks) *element {
	hl.ReadLockElement(ePtr)
	// The newly made clone never gets locked
	var cl element
	cl.initializeElement("", "")
	cl.cloneAttributes(ePtr, hl)
	return &cl
}

// cloneAttributes is a supporting function for clone
func (ePtr *element) cloneAttributes(source *element, hl *HeldLocks) {
	ePtr.ConceptID = source.ConceptID
	ePtr.Definition = source.Definition
	ePtr.Label = source.Label
	ePtr.IsCore = source.IsCore
	ePtr.OwningConceptID = source.OwningConceptID
	ePtr.ReadOnly = source.ReadOnly
	ePtr.Version.counter = source.Version.counter
	ePtr.uOfD = source.uOfD
	ePtr.URI = source.URI
}

// // editableError checks to see if the element cannot be edited because it
// // is either a core element or has been marked readOnly.
// func (ePtr *element) editableError(hl *HeldLocks) error {
// 	if ePtr.GetIsCore(hl) {
// 		return errors.New("Element.SetOwningConceptID called on core Element")
// 	}
// 	if ePtr.ReadOnly {
// 		return errors.New("Element.SetOwningConcept called on read-only Element")
// 	}
// 	return nil
// }

// GetConceptID returns the conceptID
func (ePtr *element) GetConceptID(hl *HeldLocks) string {
	hl.ReadLockElement(ePtr)
	return ePtr.ConceptID
}

// getConceptIDNoLock returns the conceptID without locking the Element.
// it is intended to support infrastructure functions only. Since the
// conceptID is never edited, this ought to be a safe operation. Even in
// cloning, in which the ConceptID is explicitly set, no other thread is
// even aware of the existence of the clone at the time the ID is set, so this
// ought to be safe.
func (ePtr *element) getConceptIDNoLock() string {
	return ePtr.ConceptID
}

// GetDefinition returns the definition if one exists
func (ePtr *element) GetDefinition(hl *HeldLocks) string {
	hl.ReadLockElement(ePtr)
	return ePtr.Definition
}

// GetFirstOwnedConceptRefinedFrom returns the first child that has the indicated abstraction as
// one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (ePtr *element) GetFirstOwnedConceptRefinedFrom(abstraction Element, hl *HeldLocks) Element {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		if element.IsRefinementOf(abstraction, hl) {
			return element
		}
	}
	return nil
}

// GetFirstOwnedConceptRefinedFromURI returns the first child that has the abstraction indicated
// by the URI as one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (ePtr *element) GetFirstOwnedConceptRefinedFromURI(abstractionURI string, hl *HeldLocks) Element {
	hl.ReadLockElement(ePtr)
	abstraction := ePtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		return ePtr.GetFirstOwnedConceptRefinedFrom(abstraction, hl)
	}
	return nil
}

// GetFirstOwnedLiteralRefinementOf returns the first child literal that has the indicated
// abstraction as one of its abstractions.
func (ePtr *element) GetFirstOwnedLiteralRefinementOf(abstraction Element, hl *HeldLocks) Literal {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case Literal:
			if element.IsRefinementOf(abstraction, hl) {
				return element.(Literal)
			}
		}
	}
	return nil
}

// GetFirstOwnedLiteralRefinementOfURI returns the first child literal that has the abstraction indicated
// by the URI as one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (ePtr *element) GetFirstOwnedLiteralRefinementOfURI(abstractionURI string, hl *HeldLocks) Literal {
	hl.ReadLockElement(ePtr)
	abstraction := ePtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		return ePtr.GetFirstOwnedLiteralRefinementOf(abstraction, hl)
	}
	return nil
}

// GetFirstOwnedReferenceRefinedFrom returns the first child reference that has the indicated
// abstraction as one of its abstractions.
func (ePtr *element) GetFirstOwnedReferenceRefinedFrom(abstraction Element, hl *HeldLocks) Reference {
	hl.ReadLockElement(ePtr)
	ownedIDs := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID)
	it := ownedIDs.Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case Reference:
			if element.(Reference).IsRefinementOf(abstraction, hl) {
				return element.(Reference)
			}
		}
	}
	return nil
}

// GetFirstOwnedReferenceRefinedFromURI returns the first child reference that has the abstraction indicated
// by the URI as one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (ePtr *element) GetFirstOwnedReferenceRefinedFromURI(abstractionURI string, hl *HeldLocks) Reference {
	hl.ReadLockElement(ePtr)
	uOfD := ePtr.uOfD
	if uOfD == nil {
		return nil
	}
	abstraction := uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		return ePtr.GetFirstOwnedReferenceRefinedFrom(abstraction, hl)
	}
	return nil
}

// GetFirstOwnedRefinementRefinedFrom returns the first child refinement that has the indicated
// abstraction as one of its abstractions.
func (ePtr *element) GetFirstOwnedRefinementRefinedFrom(abstraction Element, hl *HeldLocks) Refinement {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case Refinement:
			if element.IsRefinementOf(abstraction, hl) {
				return element.(Refinement)
			}
		}
	}
	return nil
}

// GetFirstOwnedRefinementRefinedFromURI returns the first child refinement that has the abstraction indicated
// by the URI as one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (ePtr *element) GetFirstOwnedRefinementRefinedFromURI(abstractionURI string, hl *HeldLocks) Refinement {
	hl.ReadLockElement(ePtr)
	abstraction := ePtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		return ePtr.GetFirstOwnedRefinementRefinedFrom(abstraction, hl)
	}
	return nil
}

// GetFirstOwnedConceptWithURI
func (ePtr *element) GetFirstOwnedConceptWithURI(uri string, hl *HeldLocks) Element {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		if element.GetURI(hl) == uri {
			return element
		}
	}
	return nil
}

func (ePtr *element) GetFirstOwnedLiteralRefinedFrom(abstraction Element, hl *HeldLocks) Literal {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case Literal:
			if element.IsRefinementOf(abstraction, hl) {
				return element.(Literal)
			}
		}
	}
	return nil
}

func (ePtr *element) GetFirstOwnedLiteralRefinedFromURI(uri string, hl *HeldLocks) Literal {
	hl.ReadLockElement(ePtr)
	abstraction := ePtr.uOfD.GetElementWithURI(uri)
	if abstraction != nil {
		return ePtr.GetFirstOwnedLiteralRefinedFrom(abstraction, hl)
	}
	return nil
}

func (ePtr *element) GetFirstOwnedLiteralWithURI(uri string, hl *HeldLocks) Literal {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case *literal:
			if element.GetURI(hl) == uri {
				return element.(*literal)
			}
		}
	}
	return nil
}

func (ePtr *element) GetFirstOwnedReferenceWithURI(uri string, hl *HeldLocks) Reference {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case *reference:
			if element.GetURI(hl) == uri {
				return element.(*reference)
			}
		}
	}
	return nil
}

func (ePtr *element) GetFirstOwnedRefinementWithURI(uri string, hl *HeldLocks) Refinement {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case *refinement:
			if element.GetURI(hl) == uri {
				return element.(*refinement)
			}
		}
	}
	return nil
}

// Deregister removes the registration of an Observer
func (ePtr *element) Deregister(observer Observer) error {
	ePtr.observers.Remove(observer)
	return nil
}

// FindAbstractions adds all found abstractions to supplied map
func (ePtr *element) FindAbstractions(abstractions map[string]Element, hl *HeldLocks) {
	it := ePtr.uOfD.listenersMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		listener := ePtr.uOfD.GetElement(id.(string))
		switch listener.(type) {
		case *refinement:
			abstraction := listener.(*refinement).GetAbstractConcept(hl)
			if abstraction != nil && abstraction.getConceptIDNoLock() != ePtr.getConceptIDNoLock() {
				abstractions[abstraction.GetConceptID(hl)] = abstraction
				abstraction.FindAbstractions(abstractions, hl)
			}
		}
	}
}

// FindImmediateAbstractions adds all immediate abstractions to supplied map
func (ePtr *element) FindImmediateAbstractions(abstractions map[string]Element, hl *HeldLocks) {
	// There are no abstractions without the uOfD context
	if ePtr.uOfD == nil {
		return
	}
	it := ePtr.uOfD.listenersMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		listener := ePtr.uOfD.GetElement(id.(string))
		switch listener.(type) {
		case *refinement:
			abstraction := listener.(*refinement).GetAbstractConcept(hl)
			if abstraction != nil && abstraction.getConceptIDNoLock() != ePtr.getConceptIDNoLock() {
				abstractions[abstraction.GetConceptID(hl)] = abstraction
			}
		}
	}
}

// GetIsCore returns true if the element is one of the core elements of CRL. The purpose of this
// function is to prevent SetReadOnly(true) on concepts that are built-in to CRL. Locking is
// not necessary as this value is set when the object is created and never expected to change
func (ePtr *element) GetIsCore(hl *HeldLocks) bool {
	hl.ReadLockElement(ePtr)
	return ePtr.IsCore
}

// GetGetLabel returns the label if one exists
func (ePtr *element) GetLabel(hl *HeldLocks) string {
	hl.ReadLockElement(ePtr)
	return ePtr.Label
}

func (ePtr *element) getLabelNoLock() string {
	return ePtr.Label
}

// GetOwningConceptID returns the ID of the concept that owns this one (if any)
func (ePtr *element) GetOwningConceptID(hl *HeldLocks) string {
	hl.ReadLockElement(ePtr)
	return ePtr.OwningConceptID
}

// GetOwnedConceptIDs returns the set of IDs owned by this concept. Note that if this Element is not
// presently in a uOfD it returns the empty set
func (ePtr *element) GetOwnedConceptIDs(hl *HeldLocks) mapset.Set {
	if ePtr.uOfD == nil {
		return mapset.NewSet()
	}
	return ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID)
}

// GetOwnedConcepts returns the element's owned concepts if
func (ePtr *element) GetOwnedConcepts(hl *HeldLocks) map[string]Element {
	ownedConcepts := make(map[string]Element)
	if ePtr.uOfD == nil {
		return ownedConcepts
	}
	it := ePtr.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		if element != nil {
			ownedConcepts[id.(string)] = element
		}
	}
	return ownedConcepts
}

// GetOwnedConceptsRefinedFrom returns the owned concepts with the indicated abstraction as
// one of their abstractions.
func (ePtr *element) GetOwnedConceptsRefinedFrom(abstraction Element, hl *HeldLocks) map[string]Element {
	hl.ReadLockElement(ePtr)
	matches := map[string]Element{}
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		if element.IsRefinementOf(abstraction, hl) {
			matches[element.GetConceptID(hl)] = element
		}
	}
	return matches
}

// GetOwnedConceptsRefinedFromURI returns the owned concepts that have the abstraction indicated
// by the URI as one of their abstractions.
func (ePtr *element) GetOwnedConceptsRefinedFromURI(abstractionURI string, hl *HeldLocks) map[string]Element {
	hl.ReadLockElement(ePtr)
	matches := map[string]Element{}
	abstraction := ePtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
		defer it.Stop()
		for id := range it.C {
			element := ePtr.uOfD.GetElement(id.(string))
			if element.IsRefinementOf(abstraction, hl) {
				matches[element.GetConceptID(hl)] = element
			}
		}
	}
	return matches
}

// GetOwnedDescendantsRefinedFrom returns the owned concepts with the indicated abstraction as
// one of their abstractions.
func (ePtr *element) GetOwnedDescendantsRefinedFrom(abstraction Element, hl *HeldLocks) map[string]Element {
	hl.ReadLockElement(ePtr)
	matches := map[string]Element{}
	if abstraction != nil {
		// it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
		descendantIDs := mapset.NewSet()
		ePtr.uOfD.GetConceptsOwnedConceptIDsRecursively(ePtr.ConceptID, descendantIDs, hl)
		it := descendantIDs.Iterator()
		defer it.Stop()
		for id := range it.C {
			element := ePtr.uOfD.GetElement(id.(string))
			if element.IsRefinementOf(abstraction, hl) {
				matches[element.GetConceptID(hl)] = element
			}
		}
	}
	return matches
}

// GetOwnedDescendantsRefinedFromURI returns the descendant concepts that have the indicated abstraction
// by the URI as one of their abstractions.
func (ePtr *element) GetOwnedDescendantsRefinedFromURI(abstractionURI string, hl *HeldLocks) map[string]Element {
	hl.ReadLockElement(ePtr)
	matches := map[string]Element{}
	abstraction := ePtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		// it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
		descendantIDs := mapset.NewSet()
		ePtr.uOfD.GetConceptsOwnedConceptIDsRecursively(ePtr.ConceptID, descendantIDs, hl)
		it := descendantIDs.Iterator()
		defer it.Stop()
		for id := range it.C {
			element := ePtr.uOfD.GetElement(id.(string))
			if element.IsRefinementOf(abstraction, hl) {
				matches[element.GetConceptID(hl)] = element
			}
		}
	}
	return matches
}

// GetOwnedLiteralsRefinedFrom returns the owned literals that have the indicated
// abstraction as one of their abstractions.
func (ePtr *element) GetOwnedLiteralsRefinedFrom(abstraction Element, hl *HeldLocks) map[string]Literal {
	hl.ReadLockElement(ePtr)
	matches := map[string]Literal{}
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case Literal:
			if element.IsRefinementOf(abstraction, hl) {
				matches[element.GetConceptID(hl)] = element.(Literal)
			}
		}
	}
	return matches
}

// GetOwnedLiteralsRefinedFromURI returns the child literals that have the abstraction indicated
// by the URI as one of their abstractions.
func (ePtr *element) GetOwnedLiteralsRefinedFromURI(abstractionURI string, hl *HeldLocks) map[string]Literal {
	hl.ReadLockElement(ePtr)
	matches := map[string]Literal{}
	abstraction := ePtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
		defer it.Stop()
		for id := range it.C {
			element := ePtr.uOfD.GetElement(id.(string))
			switch element.(type) {
			case Literal:
				if element.IsRefinementOf(abstraction, hl) {
					matches[element.GetConceptID(hl)] = element.(Literal)
				}
			}
		}
	}
	return matches
}

// GetOwnedReferencesRefinedFrom returns the owned references that have the indicated
// abstraction as one of their abstractions.
func (ePtr *element) GetOwnedReferencesRefinedFrom(abstraction Element, hl *HeldLocks) map[string]Reference {
	hl.ReadLockElement(ePtr)
	matches := map[string]Reference{}
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case Reference:
			if element.IsRefinementOf(abstraction, hl) {
				matches[element.GetConceptID(hl)] = element.(Reference)
			}
		}
	}
	return matches
}

// GetOwnedReferencesRefinedFromURI returns the owned references that have the abstraction indicated
// by the URI as one of their abstractions.
func (ePtr *element) GetOwnedReferencesRefinedFromURI(abstractionURI string, hl *HeldLocks) map[string]Reference {
	hl.ReadLockElement(ePtr)
	matches := map[string]Reference{}
	abstraction := ePtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
		defer it.Stop()
		for id := range it.C {
			element := ePtr.uOfD.GetElement(id.(string))
			switch element.(type) {
			case Reference:
				if element.IsRefinementOf(abstraction, hl) {
					matches[element.GetConceptID(hl)] = element.(Reference)
				}
			}
		}
	}
	return matches
}

// GetOwnedRefinementsRefinedFrom returns the owned refinements that have the indicated
// abstraction as one of their abstractions.
func (ePtr *element) GetOwnedRefinementsRefinedFrom(abstraction Element, hl *HeldLocks) map[string]Refinement {
	hl.ReadLockElement(ePtr)
	matches := map[string]Refinement{}
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := ePtr.uOfD.GetElement(id.(string))
		switch element.(type) {
		case Refinement:
			if element.IsRefinementOf(abstraction, hl) {
				matches[element.GetConceptID(hl)] = element.(Refinement)
			}
		}
	}
	return matches
}

// GetOwnedRefinementsRefinedFromURI returns the owned refinements that have the abstraction indicated
// by the URI as one of its abstractions.
func (ePtr *element) GetOwnedRefinementsRefinedFromURI(abstractionURI string, hl *HeldLocks) map[string]Refinement {
	hl.ReadLockElement(ePtr)
	matches := map[string]Refinement{}
	abstraction := ePtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
		defer it.Stop()
		for id := range it.C {
			element := ePtr.uOfD.GetElement(id.(string))
			switch element.(type) {
			case Refinement:
				if element.IsRefinementOf(abstraction, hl) {
					matches[element.GetConceptID(hl)] = element.(Refinement)
				}
			}
		}
	}
	return matches
}

// GetOwningConcept returns the Element representing the concept that owns this one (if any)
func (ePtr *element) GetOwningConcept(hl *HeldLocks) Element {
	hl.ReadLockElement(ePtr)
	if ePtr.uOfD != nil {
		return ePtr.uOfD.GetElement(ePtr.OwningConceptID)
	}
	return nil
}

// getOwningConceptNoLock returns the Element representing the concept that owns this one (if any)
func (ePtr *element) getOwningConceptNoLock() Element {
	if ePtr.uOfD != nil {
		return ePtr.uOfD.GetElement(ePtr.OwningConceptID)
	}
	return nil
}

// GetUniverseOfDiscourse returns the UniverseOfDiscourse in which the element instance resides
func (ePtr *element) GetUniverseOfDiscourse(hl *HeldLocks) *UniverseOfDiscourse {
	hl.ReadLockElement(ePtr)
	return ePtr.uOfD
}

// getUniverseOfDiscourseNoLock returns the UniverseOfDiscourse in which the element instance resides
func (ePtr *element) getUniverseOfDiscourseNoLock() *UniverseOfDiscourse {
	return ePtr.uOfD
}

// GetURI returns the URI string associated with the element if there is one
func (ePtr *element) GetURI(hl *HeldLocks) string {
	hl.ReadLockElement(ePtr)
	return ePtr.URI
}

// GetVersion returns the version of the element
func (ePtr *element) GetVersion(hl *HeldLocks) int {
	hl.ReadLockElement(ePtr)
	return ePtr.Version.getVersion()
}

// IsRefinementOf returns true if the given abstraction is contained in the abstractions set
// of this element. No locking is required since the StringIntMap does its own locking
func (ePtr *element) IsRefinementOf(abstraction Element, hl *HeldLocks) bool {
	hl.ReadLockElement(ePtr)
	// Get the actual element so that we can get the correct type
	fullElement := ePtr.uOfD.GetElement(ePtr.ConceptID)
	// Check to see whether the abstraction is one of the core classes
	abstractionURI := abstraction.GetURI(hl)
	switch abstractionURI {
	case ElementURI:
		return true
	case LiteralURI:
		switch fullElement.(type) {
		case Literal:
			return true
		}
	case ReferenceURI:
		switch fullElement.(type) {
		case Reference:
			return true
		}
	case RefinementURI:
		switch fullElement.(type) {
		case Refinement:
			return true
		}
	}
	it := ePtr.uOfD.listenersMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		listener := ePtr.uOfD.GetElement(id.(string))
		switch listener.(type) {
		case Refinement:
			foundAbstraction := listener.(Refinement).GetAbstractConcept(hl)
			if foundAbstraction == nil {
				continue
			}
			if foundAbstraction.getConceptIDNoLock() == ePtr.ConceptID {
				continue
			}
			if foundAbstraction == abstraction {
				return true
			}
			if foundAbstraction != nil {
				foundRecursively := foundAbstraction.IsRefinementOf(abstraction, hl)
				if foundRecursively {
					return true
				}
			}
		}
	}
	return false
}

func (ePtr *element) IsRefinementOfURI(uri string, hl *HeldLocks) bool {
	hl.ReadLockElement(ePtr)
	if ePtr.uOfD == nil {
		return false
	}
	abstraction := ePtr.uOfD.GetElementWithURI(uri)
	if abstraction == nil {
		return false
	}
	return ePtr.IsRefinementOf(abstraction, hl)
}

func (ePtr *element) incrementVersion(hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	if ePtr.uOfD != nil {
		// UofD may be nil during the deletion of this element
		ePtr.uOfD.preChange(ePtr, hl)
		ePtr.Version.incrementVersion()
		if ePtr.OwningConceptID != "" {
			owningConcept := ePtr.uOfD.GetElement(ePtr.OwningConceptID)
			// the owning concept may also be in the process of deletion
			if owningConcept != nil {
				owningConcept.incrementVersion(hl)
			}
		}
	}
}

// initializeElement creates the identifier (using the uri if supplied) and
// creates the abstractions, ownedConcepts, and referrencingConcpsts maps.
// Note that initialization is not considered a change, so the version counter is not incremented
// nor are monitors of this element notified of changes.
func (ePtr *element) initializeElement(identifier string, uri string) {
	ePtr.ConceptID = identifier
	ePtr.Version = newVersionCounter()
	ePtr.URI = uri
	ePtr.observers = mapset.NewSet()
}

// IsReadOnly returns a boolean indicating whether the concept can be modified.
func (ePtr *element) IsReadOnly(hl *HeldLocks) bool {
	hl.ReadLockElement(ePtr)
	return ePtr.ReadOnly
}

// isEditable checks to see if the element cannot be edited because it
// is either a core element or has been marked readOnly.
func (ePtr *element) isEditable(hl *HeldLocks) bool {
	if ePtr.GetIsCore(hl) || ePtr.IsReadOnly(hl) {
		return false
	}
	return true
}

// isEquivalent only checks the element attributes. It ignores the uOfD.
func (ePtr *element) isEquivalent(hl1 *HeldLocks, el *element, hl2 *HeldLocks, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) > 0 {
		print = printExceptions[0]
	}
	hl1.ReadLockElement(ePtr)
	hl2.ReadLockElement(el)
	if ePtr.ConceptID != el.ConceptID {
		if print {
			log.Printf("In element.isEquivalent, ConceptIDs do not match")
		}
		return false
	}
	if ePtr.Definition != el.Definition {
		if print {
			log.Printf("In element.isEquivalent, Definitions do not match")
		}
		return false
	}
	if ePtr.IsCore != el.IsCore {
		if print {
			log.Printf("In element.isEquivalent, IsCore do not match")
		}
		return false
	}
	if ePtr.Label != el.Label {
		if print {
			log.Printf("In element.isEquivalent, Labels do not match")
		}
		return false
	}
	if ePtr.OwningConceptID != el.OwningConceptID {
		if print {
			log.Printf("In element.isEquivalent, OwningConceptIDs do not match")
		}
		return false
	}
	if ePtr.ReadOnly != el.ReadOnly {
		if print {
			log.Printf("In element.isEquivalent, ReadOnly does not match")
		}
		return false
	}
	if ePtr.Version.getVersion() != el.Version.getVersion() {
		if print {
			log.Printf("In element.isEquivalent, Versions do not match")
		}
		return false
	}
	if ePtr.URI != el.URI {
		if print {
			log.Printf("In element.isEquivalent, URIs do not match")
		}
		return false
	}
	return true
}

// IsOwnedConcept returns true if the supplied element is an owned concept. Note that
// there is an interval of time during editing in which the child's owner will be set but the child
// has not yet been added to the element's OwnedConcepts list. Similarly, there is an interval of time
// during editing during which the child's owner has been changed but the original owner's OwnedConcept
// list has not yet been updated.
func (ePtr *element) IsOwnedConcept(el Element, hl *HeldLocks) bool {
	hl.ReadLockElement(ePtr)
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		child := ePtr.uOfD.GetElement(id.(string))
		if el.GetConceptID(hl) == child.GetConceptID(hl) {
			return true
		}
	}
	return false
}

// MarshalJSON produces a byte string JSON representation of the Element
func (ePtr *element) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(ePtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := ePtr.marshalElementFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (ePtr *element) marshalElementFields(buffer *bytes.Buffer) error {
	buffer.WriteString(fmt.Sprintf("\"ConceptID\":\"%s\",", ePtr.ConceptID))
	buffer.WriteString(fmt.Sprintf("\"OwningConceptID\":\"%s\",", ePtr.OwningConceptID))
	buffer.WriteString(fmt.Sprintf("\"Label\":\"%s\",", ePtr.Label))
	buffer.WriteString(fmt.Sprintf("\"Definition\":\"%s\",", ePtr.Definition))
	buffer.WriteString(fmt.Sprintf("\"URI\":\"%s\",", ePtr.URI))
	buffer.WriteString(fmt.Sprintf("\"Version\":\"%d\",", ePtr.Version.getVersion()))
	buffer.WriteString(fmt.Sprintf("\"IsCore\":\"%t\",", ePtr.IsCore))
	buffer.WriteString(fmt.Sprintf("\"ReadOnly\":\"%t\"", ePtr.ReadOnly))
	return nil
}

// NotifyAll passes the notification to all registered Observers
func (ePtr *element) NotifyAll(notification *ChangeNotification, hl *HeldLocks) error {
	it := ePtr.observers.Iterator()
	defer it.Stop()
	for observer := range it.C {
		err := observer.(Observer).Update(notification, hl)
		if err != nil {
			return errors.Wrap(err, "element.NotifyAll failed")
		}
	}
	return nil
}

func (ePtr *element) notifyListeners(notification *ChangeNotification, hl *HeldLocks) error {
	hl.ReadLockElement(ePtr)
	if ePtr.uOfD != nil {
		it := ePtr.uOfD.listenersMap.GetMappedValues(ePtr.ConceptID).Iterator()
		defer it.Stop()
		for id := range it.C {
			listener := ePtr.uOfD.GetElement(id.(string))
			if listener.GetConceptID(hl) == ePtr.OwningConceptID && notification.GetReportingElementType() == "core.Reference" {
				forwardingChangeNotification, err := ePtr.uOfD.NewForwardingChangeNotification(ePtr, ForwardedChange, notification, hl)
				if err != nil {
					return errors.Wrap(err, "element.notifyListeners failed")
				}
				err = ePtr.uOfD.queueFunctionExecutions(listener, forwardingChangeNotification, hl)
				if err != nil {
					return errors.Wrap(err, "element.notifyListeners failed")
				}
				continue
			}
			switch listener.(type) {
			case Reference:
				if !(notification.GetNatureOfChange() == ReferencedConceptChanged && notification.GetReportingElementID() == listener.(Reference).GetConceptID(hl)) {
					err := ePtr.uOfD.queueFunctionExecutions(listener, notification, hl)
					if err != nil {
						return errors.Wrap(err, "element.notifyListeners failed")
					}
				}
			case Refinement:
				if !((notification.GetNatureOfChange() == AbstractConceptChanged || notification.GetNatureOfChange() == RefinedConceptChanged) && notification.GetReportingElementID() == listener.(Refinement).GetConceptID(hl)) {
					err := ePtr.uOfD.queueFunctionExecutions(listener, notification, hl)
					if err != nil {
						return errors.Wrap(err, "element.notifyListeners failed")
					}
				}
			}
		}
	}
	return nil
}

// recoverElementFields() is used when de-serializing an element. The activities in restoring the
// element are not considered changes so the version counter is not incremented and the monitors of this
// element are not notified of chaanges.
func (ePtr *element) recoverElementFields(unmarshaledData *map[string]json.RawMessage, hl *HeldLocks) error {
	// ConceptID
	var recoveredConceptID string
	err := json.Unmarshal((*unmarshaledData)["ConceptID"], &recoveredConceptID)
	if err != nil {
		log.Printf("Recovery of Element.ConceptID as string failed\n")
		return err
	}
	ePtr.ConceptID = recoveredConceptID
	// Definition
	var recoveredDefinition string
	err = json.Unmarshal((*unmarshaledData)["Definition"], &recoveredDefinition)
	if err != nil {
		log.Printf("Recovery of Element.Definition as string failed\n")
		return err
	}
	ePtr.Definition = recoveredDefinition
	// Label
	var recoveredLabel string
	err = json.Unmarshal((*unmarshaledData)["Label"], &recoveredLabel)
	if err != nil {
		log.Printf("Recovery of Element.Label as string failed\n")
		return err
	}
	ePtr.Label = recoveredLabel
	// IsCore
	var recoveredIsCore string
	err = json.Unmarshal((*unmarshaledData)["IsCore"], &recoveredIsCore)
	if err != nil {
		log.Printf("Recovery of Element.IsCore as string failed\n")
		return err
	}
	ePtr.IsCore, err = strconv.ParseBool(recoveredIsCore)
	if err != nil {
		log.Printf("Conversion of IsCOre from string to bool failed")
		return err
	}
	// OwningConceptID
	var recoveredOwningConceptID string
	err = json.Unmarshal((*unmarshaledData)["OwningConceptID"], &recoveredOwningConceptID)
	if err != nil {
		log.Printf("Recovery of Element.OwningConceptID as string failed\n")
		return err
	}
	ePtr.OwningConceptID = recoveredOwningConceptID
	// ReadOnly
	var recoveredReadOnly string
	err = json.Unmarshal((*unmarshaledData)["ReadOnly"], &recoveredReadOnly)
	if err != nil {
		log.Printf("Recovery of Element.ReadOnly as string failed\n")
		return err
	}
	ePtr.ReadOnly, err = strconv.ParseBool(recoveredReadOnly)
	if err != nil {
		log.Printf("Conversion of ReadOnly from string to bool failed")
		return err
	}
	// Version
	var recoveredVersion string
	err = json.Unmarshal((*unmarshaledData)["Version"], &recoveredVersion)
	if err != nil {
		log.Printf("Recovery of BaseElement.version failed\n")
		return err
	}
	ePtr.Version.counter, err = strconv.Atoi(recoveredVersion)
	if err != nil {
		log.Printf("Conversion of Element.version to integer failed\n")
		return err
	}
	// URI
	var recoveredURI string
	err = json.Unmarshal((*unmarshaledData)["URI"], &recoveredURI)
	if err != nil {
		log.Printf("Recovery of Element.URI as string failed\n")
		return err
	}
	ePtr.URI = recoveredURI
	return nil
}

// removeListener removes the indicated Element as a listening concept.
func (ePtr *element) removeListener(listeningConceptID string, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	ePtr.uOfD.preChange(ePtr, hl)
	ePtr.uOfD.listenersMap.RemoveMappedValue(ePtr.ConceptID, listeningConceptID)
}

// Register adds the registration of an Observer
func (ePtr *element) Register(observer Observer) error {
	ePtr.observers.Add(observer)
	return nil
}

// removeOwnedConcept removes the indicated Element as a child (owned) concept.
func (ePtr *element) removeOwnedConcept(ownedConceptID string, hl *HeldLocks) error {
	hl.ReadLockElement(ePtr)
	if ePtr.IsReadOnly(hl) {
		return errors.New("Element.removedOwnedConcept called on read-only Element")
	}
	ePtr.uOfD.preChange(ePtr, hl)
	ePtr.incrementVersion(hl)
	ePtr.uOfD.ownedIDsMap.RemoveMappedValue(ePtr.ConceptID, ownedConceptID)
	return nil
}

// SetDefinition sets the definition of the Element
func (ePtr *element) SetDefinition(def string, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	if ePtr.isEditable(hl) == false {
		return errors.New("element.SetDefinition failed because the element is not editable")
	}
	if ePtr.Definition != def {
		ePtr.uOfD.preChange(ePtr, hl)
		beforeState, err := NewConceptState(ePtr)
		if err != nil {
			return errors.Wrap(err, "element.SetDefinition failed")
		}
		ePtr.incrementVersion(hl)
		ePtr.Definition = def
		afterState, err2 := NewConceptState(ePtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetDefinition failed")
		}
		err = ePtr.uOfD.SendConceptChangeNotification(ePtr, beforeState, afterState, hl)
		if err != nil {
			return errors.Wrap(err, "element.SetDefinition failed")
		}
	}
	return nil
}

// SetIsCore sets the flag indicating that the element is a Core concept and cannot be edited. Once set, this flag cannot be cleared.
func (ePtr *element) SetIsCore(hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	if ePtr.IsCore != true {
		ePtr.uOfD.preChange(ePtr, hl)
		beforeState, err := NewConceptState(ePtr)
		if err != nil {
			return errors.Wrap(err, "element.SetIsCore failed")
		}
		ePtr.incrementVersion(hl)
		ePtr.IsCore = true
		afterState, err2 := NewConceptState(ePtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetIsCore failed")
		}
		err = ePtr.uOfD.SendConceptChangeNotification(ePtr, beforeState, afterState, hl)
		if err != nil {
			return errors.Wrap(err, "element.SetIsCore failed")
		}
	}
	return nil
}

// SetIsCoreRecursively recursively sets the flag indicating that the element is a Core concept and cannot be edited. Once set, this flag cannot be cleared.
func (ePtr *element) SetIsCoreRecursively(hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	err := ePtr.SetIsCore(hl)
	if err != nil {
		return errors.Wrap(err, "Element.SetIsCoreRecursively failed")
	}
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		el := ePtr.uOfD.GetElement(id.(string))
		err = el.SetIsCoreRecursively(hl)
		if err != nil {
			return errors.Wrap(err, "Element.SetIsCoreRecursively failed")
		}
	}
	return nil
}

// SetLabel sets the label of the Element
func (ePtr *element) SetLabel(label string, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	if ePtr.isEditable(hl) == false {
		return errors.New("element.SetLabel failed because the element is not editable")
	}
	if ePtr.Label != label {
		ePtr.uOfD.preChange(ePtr, hl)
		beforeState, err := NewConceptState(ePtr)
		if err != nil {
			return errors.Wrap(err, "element.SetLabel failed")
		}
		ePtr.incrementVersion(hl)
		ePtr.Label = label
		afterState, err2 := NewConceptState(ePtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetLabel failed")
		}
		err = ePtr.uOfD.SendConceptChangeNotification(ePtr, beforeState, afterState, hl)
		if err != nil {
			return errors.Wrap(err, "element.SetLabel failed")
		}
	}
	return nil
}

// SetOwningConcept takes the ID of the supplied concept and call SetOwningConceptID. It first checks to
// determine whether the new owner is editable and will throw an error if it is not
func (ePtr *element) SetOwningConcept(el Element, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	id := ""
	if el != nil {
		if el.isEditable(hl) == false {
			return errors.New("element.SetOwningConcept called with an owner that is not editable")
		}
		id = el.getConceptIDNoLock()
	}
	err := ePtr.SetOwningConceptID(id, hl)
	if err != nil {
		errors.Wrap(err, "element.SetOwningConcept failed")
	}
	return nil
}

// SetOwningConceptID sets the ID of the owning concept for the element
// Design Note: the argument is the identifier rather than the Element to ensure
// the correct type of the owning concept is recorded.
func (ePtr *element) SetOwningConceptID(ocID string, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	if ePtr.isEditable(hl) == false {
		return errors.New("element.SetOwningConceptID failed because the element is not editable")
	}
	if ocID == ePtr.ConceptID {
		return errors.New("element.SetOwningConceptID called with itself as owner")
	}
	newOwner := ePtr.uOfD.GetElement(ocID)
	if newOwner != nil && newOwner.isEditable(hl) == false {
		return errors.New("element.SetOwningConceptID called with new owner not editable")
	}
	oldOwner := ePtr.GetOwningConcept(hl)
	if oldOwner != nil && oldOwner.isEditable(hl) == false {
		return errors.New("element.SetOwningConceptID called with old owner not editable")
	}
	// Do nothing if there is no change
	if ePtr.OwningConceptID != ocID {
		ePtr.uOfD.preChange(ePtr, hl)
		beforeState, err := NewConceptState(ePtr)
		if err != nil {
			return errors.Wrap(err, "element.SetOwningConceptID failed")
		}
		var ownerBeforeState *ConceptState
		if oldOwner != nil {
			oldOwner.removeOwnedConcept(ePtr.ConceptID, hl)
			ownerBeforeState, err = NewConceptState(oldOwner)
			if err != nil {
				return errors.Wrap(err, "element.SetOwningConceptID failed")
			}
		}
		ePtr.incrementVersion(hl)
		var ownerAfterState *ConceptState
		if newOwner != nil {
			newOwner.addOwnedConcept(ePtr.ConceptID, hl)
			ownerAfterState, err = NewConceptState(newOwner)
			if err != nil {
				return errors.Wrap(err, "element.SetOwningConceptID failed")
			}
		}
		ePtr.OwningConceptID = ocID
		afterState, err2 := NewConceptState(ePtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetOwningConceptID failed")
		}
		err = ePtr.uOfD.SendPointerChangeNotification(ePtr, OwningConceptChanged, beforeState, afterState, ownerBeforeState, ownerAfterState, hl)
		if err != nil {
			return errors.Wrap(err, "element.SetOwningConceptID failed")
		}
	}
	return nil
}

// SetReadOnly provides a mechanism for preventing modifications to concepts. It will throw an error
// if the concept is one of the CRL core concepts, as these can never be made writable. It will also
// throw an error if its owner is read only and this call tries to set read only false.
func (ePtr *element) SetReadOnly(value bool, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	if ePtr.GetIsCore(hl) == true {
		return errors.New("element.SetReadOnly failed because element is a core element")
	}
	if ePtr.GetOwningConcept(hl) != nil {
		if ePtr.GetOwningConcept(hl).IsReadOnly(hl) == true && value == false {
			return errors.New("element.SetReadOnly failed because the owner is read only")
		}
	}
	if ePtr.ReadOnly != value {
		ePtr.uOfD.preChange(ePtr, hl)
		beforeState, err := NewConceptState(ePtr)
		if err != nil {
			return errors.Wrap(err, "element.SetReadOnly failed")
		}
		ePtr.incrementVersion(hl)
		ePtr.ReadOnly = value
		afterState, err2 := NewConceptState(ePtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetDeSetReadOnlyfinition failed")
		}
		err = ePtr.uOfD.SendConceptChangeNotification(ePtr, beforeState, afterState, hl)
		if err != nil {
			return errors.Wrap(err, "element.SetDeSetReadOnlyfinition failed")
		}
	}
	return nil
}

func (ePtr *element) SetReadOnlyRecursively(value bool, hl *HeldLocks) error {
	err := ePtr.SetReadOnly(value, hl)
	if err != nil {
		return errors.Wrap(err, "Element.SetReadOnlyRecursively failed")
	}
	it := ePtr.uOfD.ownedIDsMap.GetMappedValues(ePtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		el := ePtr.uOfD.GetElement(id.(string))
		err = el.SetReadOnlyRecursively(value, hl)
		if err != nil {
			return errors.Wrap(err, "Element.SetReadOnlyRecursively failed")
		}
	}
	return nil
}

// setUniverseOfDiscourse is intended to be called only by the UniverseOfDiscourse
func (ePtr *element) setUniverseOfDiscourse(uOfD *UniverseOfDiscourse, hl *HeldLocks) {
	hl.WriteLockElement(ePtr)
	ePtr.uOfD = uOfD
}

// SetURI sets the URI of the Element
func (ePtr *element) SetURI(uri string, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	if ePtr.isEditable(hl) == false {
		return errors.New("element.SetURI failed because the elementis not editable")
	}
	if ePtr.URI != uri {
		foundElement := ePtr.uOfD.GetElementWithURI(uri)
		if foundElement != nil && foundElement.GetConceptID(hl) != ePtr.ConceptID {
			return errors.New("Element already exists with URI " + uri)
		}
		ePtr.uOfD.preChange(ePtr, hl)
		beforeState, err := NewConceptState(ePtr)
		if err != nil {
			return errors.Wrap(err, "element.SetURI failed")
		}
		ePtr.uOfD.changeURIForElement(ePtr, ePtr.URI, uri)
		ePtr.incrementVersion(hl)
		ePtr.URI = uri
		afterState, err2 := NewConceptState(ePtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetURI failed")
		}
		err = ePtr.uOfD.SendConceptChangeNotification(ePtr, beforeState, afterState, hl)
		if err != nil {
			return errors.Wrap(err, "element.SetURI failed")
		}
	}
	return nil
}

func (ePtr *element) TraceableReadLock(hl *HeldLocks) {
	if TraceLocks {
		log.Printf("HL %p about to read lock Element %p %s\n", hl, ePtr, ePtr.Label)
	}
	ePtr.RLock()
}

func (ePtr *element) TraceableWriteLock(hl *HeldLocks) {
	if TraceLocks {
		log.Printf("HL %p about to write lock Element %p %s\n", hl, ePtr, ePtr.Label)
	}
	ePtr.Lock()
}

func (ePtr *element) TraceableReadUnlock(hl *HeldLocks) {
	if TraceLocks {
		log.Printf("HL %p about to read unlock Element %p %s\n", hl, ePtr, ePtr.Label)
	}
	ePtr.RUnlock()
}

func (ePtr *element) TraceableWriteUnlock(hl *HeldLocks) {
	if TraceLocks {
		log.Printf("HL %p about to write unlock Element %p %s\n", hl, ePtr, ePtr.Label)
	}
	ePtr.Unlock()
}

// Element is the representation of a concept
type Element interface {
	Subject
	addListener(string, *HeldLocks)
	addOwnedConcept(string, *HeldLocks)
	addRecoveredOwnedConcept(string, *HeldLocks)
	// editableError(*HeldLocks) error
	FindAbstractions(map[string]Element, *HeldLocks)
	FindImmediateAbstractions(map[string]Element, *HeldLocks)
	GetConceptID(*HeldLocks) string
	getConceptIDNoLock() string
	GetDefinition(*HeldLocks) string
	GetFirstOwnedConceptRefinedFrom(Element, *HeldLocks) Element
	GetFirstOwnedConceptRefinedFromURI(string, *HeldLocks) Element
	GetFirstOwnedLiteralRefinementOf(Element, *HeldLocks) Literal
	GetFirstOwnedLiteralRefinementOfURI(string, *HeldLocks) Literal
	GetFirstOwnedReferenceRefinedFrom(Element, *HeldLocks) Reference
	GetFirstOwnedReferenceRefinedFromURI(string, *HeldLocks) Reference
	GetFirstOwnedRefinementRefinedFrom(Element, *HeldLocks) Refinement
	GetFirstOwnedRefinementRefinedFromURI(string, *HeldLocks) Refinement
	GetFirstOwnedConceptWithURI(string, *HeldLocks) Element
	GetFirstOwnedLiteralRefinedFrom(Element, *HeldLocks) Literal
	GetFirstOwnedLiteralRefinedFromURI(string, *HeldLocks) Literal
	GetFirstOwnedLiteralWithURI(string, *HeldLocks) Literal
	GetFirstOwnedReferenceWithURI(string, *HeldLocks) Reference
	GetFirstOwnedRefinementWithURI(string, *HeldLocks) Refinement
	GetIsCore(*HeldLocks) bool
	GetLabel(*HeldLocks) string
	getLabelNoLock() string
	// getListeners(*HeldLocks) (mapset.Set, error)
	// GetOwnedConcepts(*HeldLocks) mapset.Set
	// GetOwnedConceptsRecursively(mapset.Set, *HeldLocks)
	GetOwnedConcepts(hl *HeldLocks) map[string]Element
	GetOwnedConceptIDs(hl *HeldLocks) mapset.Set
	GetOwnedConceptsRefinedFrom(Element, *HeldLocks) map[string]Element
	GetOwnedConceptsRefinedFromURI(string, *HeldLocks) map[string]Element
	GetOwnedDescendantsRefinedFrom(Element, *HeldLocks) map[string]Element
	GetOwnedDescendantsRefinedFromURI(string, *HeldLocks) map[string]Element
	GetOwnedLiteralsRefinedFrom(Element, *HeldLocks) map[string]Literal
	GetOwnedLiteralsRefinedFromURI(string, *HeldLocks) map[string]Literal
	GetOwnedReferencesRefinedFrom(Element, *HeldLocks) map[string]Reference
	GetOwnedReferencesRefinedFromURI(string, *HeldLocks) map[string]Reference
	GetOwnedRefinementsRefinedFrom(Element, *HeldLocks) map[string]Refinement
	GetOwnedRefinementsRefinedFromURI(string, *HeldLocks) map[string]Refinement
	GetOwningConceptID(*HeldLocks) string
	GetOwningConcept(*HeldLocks) Element
	getOwningConceptNoLock() Element
	GetUniverseOfDiscourse(*HeldLocks) *UniverseOfDiscourse
	getUniverseOfDiscourseNoLock() *UniverseOfDiscourse
	GetURI(*HeldLocks) string
	GetVersion(*HeldLocks) int
	isEditable(*HeldLocks) bool
	IsRefinementOf(Element, *HeldLocks) bool
	IsRefinementOfURI(string, *HeldLocks) bool
	incrementVersion(*HeldLocks)
	IsOwnedConcept(Element, *HeldLocks) bool
	IsReadOnly(*HeldLocks) bool
	MarshalJSON() ([]byte, error)
	notifyListeners(*ChangeNotification, *HeldLocks) error
	removeListener(string, *HeldLocks)
	removeOwnedConcept(string, *HeldLocks) error
	SetDefinition(string, *HeldLocks) error
	SetIsCore(*HeldLocks) error
	SetIsCoreRecursively(*HeldLocks) error
	SetLabel(string, *HeldLocks) error
	SetOwningConcept(Element, *HeldLocks) error
	SetOwningConceptID(string, *HeldLocks) error
	SetReadOnly(bool, *HeldLocks) error
	SetReadOnlyRecursively(bool, *HeldLocks) error
	setUniverseOfDiscourse(*UniverseOfDiscourse, *HeldLocks)
	SetURI(string, *HeldLocks) error
	TraceableReadLock(*HeldLocks)
	TraceableWriteLock(*HeldLocks)
	TraceableReadUnlock(*HeldLocks)
	TraceableWriteUnlock(*HeldLocks)
}
