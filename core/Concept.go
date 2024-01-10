package core

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

	"github.com/pkg/errors"

	mapset "github.com/deckarep/golang-set"
)

// ConceptType indicates the type of the core concept being represented
type ConceptType int

// Element is the representation of a simple concept
// Literal is the representatio of a literal string
// Reference is the representation of a reference to another concept
// Refinement is the representation of an abstraction-refinement relationship
const (
	Element ConceptType = iota
	Literal
	Reference
	Refinement
)

// ConceptTypeToString returns the string representation of the ConceptType
func ConceptTypeToString(conceptType ConceptType) string {
	switch conceptType {
	case Element:
		return "Element"
	case Literal:
		return "Literal"
	case Reference:
		return "Reference"
	case Refinement:
		return "Refinement"
	}
	return ""
}

// StringToConceptType returns the ConceptType corresponding to the string
func StringToConceptType(s string) (ConceptType, error) {
	switch s {
	case "Element":
		return Element, nil
	case "Literal":
		return Literal, nil
	case "Reference":
		return Reference, nil
	case "Refinement":
		return Refinement, nil
	default:
		return 0, errors.New("Invalid concept name: " + s)
	}
}

// Concept is the core representation of a concept
type Concept struct {
	sync.RWMutex
	Label     string
	ConceptID string
	ConceptType
	Definition              string
	AbstractConceptID       string
	LiteralValue            string
	OwningConceptID         string
	ReferencedConceptID     string
	ReferencedAttributeName AttributeName
	RefinedConceptID        string
	Version                 *versionCounter
	uOfD                    *UniverseOfDiscourse
	URI                     string
	IsCore                  bool
	ReadOnly                bool
	observers               mapset.Set
}

// addOwnedConcept adds the indicated Element as a child (owned) concept.
// This is purely an internal housekeeping method. Note that
// no checking of whether the Element is read-only is performed here. This check
// is performed by the child
func (cPtr *Concept) addOwnedConcept(ownedConceptID string, trans *Transaction) {
	trans.ReadLockElement(cPtr)
	if !cPtr.uOfD.ownedIDsMap.ContainsMappedValue(cPtr.ConceptID, ownedConceptID) {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ addOwnedConcept")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		cPtr.Version.incrementVersion()
		cPtr.uOfD.ownedIDsMap.addMappedValue(cPtr.GetConceptID(trans), ownedConceptID)
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
}

// addRecoveredOwnedConcept adds the indicated Element as a child (owned) concept without incrementing
// the version.
// This is purely an internal housekeeping method. Note that
// no checking of whether the Element is read-only is performed here. This check
// is performed by the child
func (cPtr *Concept) addRecoveredOwnedConcept(ownedConceptID string, trans *Transaction) {
	trans.ReadLockElement(cPtr)
	if !cPtr.uOfD.ownedIDsMap.ContainsMappedValue(cPtr.ConceptID, ownedConceptID) {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ addRecoveredOwnedConcept")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		cPtr.uOfD.ownedIDsMap.addMappedValue(cPtr.ConceptID, ownedConceptID)
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
}

// addListener adds the indicated Element as a listening concept.
// This is an internal housekeeping method.
func (cPtr *Concept) addListener(listeningConceptID string, trans *Transaction) {
	trans.ReadLockElement(cPtr)
	if !cPtr.uOfD.listenersMap.ContainsMappedValue(cPtr.ConceptID, listeningConceptID) {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ addListener")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		cPtr.uOfD.listenersMap.addMappedValue(cPtr.ConceptID, listeningConceptID)
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
}

// clone is an internal function that makes a copy of the given element - including its
// identifier. This is done only to support undo/redo: the clone should NEVER be added to the
// universe of discourse
func (cPtr *Concept) clone(trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	// The newly made clone never gets locked
	var cl Concept
	cl.initializeConcept(cPtr.ConceptType, "", "")
	cl.cloneAttributes(cPtr, trans)
	return &cl
}

// cloneAttributes is a supporting function for clone
func (cPtr *Concept) cloneAttributes(source *Concept, trans *Transaction) {
	cPtr.AbstractConceptID = source.AbstractConceptID
	cPtr.ConceptID = source.ConceptID
	cPtr.Definition = source.Definition
	cPtr.Label = source.Label
	cPtr.LiteralValue = source.LiteralValue
	cPtr.IsCore = source.IsCore
	cPtr.OwningConceptID = source.OwningConceptID
	cPtr.ReadOnly = source.ReadOnly
	cPtr.ReferencedConceptID = source.ReferencedConceptID
	cPtr.RefinedConceptID = source.RefinedConceptID
	cPtr.uOfD = source.uOfD
	cPtr.URI = source.URI
	cPtr.Version.counter = source.Version.counter
}

// GetAbstractConcept returns the abstract concept, which will be nil unless this is a Refinement
func (cPtr *Concept) GetAbstractConcept(trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	return cPtr.uOfD.GetElement(cPtr.AbstractConceptID)
}

func (cPtr *Concept) getAbstractConceptNoLock() *Concept {
	return cPtr.uOfD.GetElement(cPtr.AbstractConceptID)
}

// GetAbstractConceptID returns the ID of the abstract concept, which will be the empty string unless this is a Refinement
func (cPtr *Concept) GetAbstractConceptID(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.AbstractConceptID
}

func (cPtr *Concept) getAbstractConceptIDNoLock() string {
	return cPtr.AbstractConceptID
}

// GetConceptID returns the conceptID
func (cPtr *Concept) GetConceptID(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.ConceptID
}

// getConceptIDNoLock returns the conceptID without locking the Element.
// it is intended to support infrastructure functions only. Since the
// conceptID is never edited, this ought to be a safe operation. Even in
// cloning, in which the ConceptID is explicitly set, no other thread is
// even aware of the existence of the clone at the time the ID is set, so this
// ought to be safe.
func (cPtr *Concept) getConceptIDNoLock() string {
	return cPtr.ConceptID
}

// GetConceptType returns the type of the concept
func (cPtr *Concept) GetConceptType() ConceptType {
	return cPtr.ConceptType
}

// GetDefinition returns the definition if one exists
func (cPtr *Concept) GetDefinition(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.Definition
}

// GetFirstOwnedConceptRefinedFrom returns the first child that has the indicated abstraction as
// one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (cPtr *Concept) GetFirstOwnedConceptRefinedFrom(abstraction *Concept, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		if element.IsRefinementOf(abstraction, trans) {
			it.Stop()
			return element
		}
	}
	return nil
}

// GetFirstOwnedConceptRefinedFromURI returns the first child that has the abstraction indicated
// by the URI as one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (cPtr *Concept) GetFirstOwnedConceptRefinedFromURI(abstractionURI string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	abstraction := cPtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		return cPtr.GetFirstOwnedConceptRefinedFrom(abstraction, trans)
	}
	return nil
}

// GetFirstOwnedLiteralRefinementOf returns the first child literal that has the indicated
// abstraction as one of its abstractions.
func (cPtr *Concept) GetFirstOwnedLiteralRefinementOf(abstraction *Concept, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		switch element.GetConceptType() {
		case Literal:
			if element.IsRefinementOf(abstraction, trans) {
				it.Stop()
				return element
			}
		}
	}
	return nil
}

// GetFirstOwnedLiteralRefinementOfURI returns the first child literal that has the abstraction indicated
// by the URI as one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (cPtr *Concept) GetFirstOwnedLiteralRefinementOfURI(abstractionURI string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	abstraction := cPtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		return cPtr.GetFirstOwnedLiteralRefinementOf(abstraction, trans)
	}
	return nil
}

// GetFirstOwnedReferenceRefinedFrom returns the first child reference that has the indicated
// abstraction as one of its abstractions.
func (cPtr *Concept) GetFirstOwnedReferenceRefinedFrom(abstraction *Concept, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	ownedIDs := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID)
	it := ownedIDs.Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		switch element.GetConceptType() {
		case Reference:
			if element.IsRefinementOf(abstraction, trans) {
				it.Stop()
				return element
			}
		}
	}
	return nil
}

// GetFirstOwnedReferenceRefinedFromURI returns the first child reference that has the abstraction indicated
// by the URI as one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (cPtr *Concept) GetFirstOwnedReferenceRefinedFromURI(abstractionURI string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	uOfD := cPtr.uOfD
	if uOfD == nil {
		return nil
	}
	abstraction := uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		return cPtr.GetFirstOwnedReferenceRefinedFrom(abstraction, trans)
	}
	return nil
}

// GetFirstOwnedRefinementRefinedFrom returns the first child refinement that has the indicated
// abstraction as one of its abstractions.
func (cPtr *Concept) GetFirstOwnedRefinementRefinedFrom(abstraction *Concept, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		switch element.GetConceptType() {
		case Refinement:
			if element.IsRefinementOf(abstraction, trans) {
				it.Stop()
				return element
			}
		}
	}
	return nil
}

// GetFirstOwnedRefinementRefinedFromURI returns the first child refinement that has the abstraction indicated
// by the URI as one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (cPtr *Concept) GetFirstOwnedRefinementRefinedFromURI(abstractionURI string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	abstraction := cPtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		return cPtr.GetFirstOwnedRefinementRefinedFrom(abstraction, trans)
	}
	return nil
}

// GetFirstOwnedConceptWithURI returns the first owned concept with the indicated URI
func (cPtr *Concept) GetFirstOwnedConceptWithURI(uri string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		if element.GetURI(trans) == uri {
			it.Stop()
			return element
		}
	}
	return nil
}

// GetFirstOwnedLiteralRefinedFrom returns the first owned literal that is refined from the supplied abstract concept
func (cPtr *Concept) GetFirstOwnedLiteralRefinedFrom(abstraction *Concept, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		switch element.GetConceptType() {
		case Literal:
			if element.IsRefinementOf(abstraction, trans) {
				it.Stop()
				return element
			}
		}
	}
	return nil
}

// GetFirstOwnedLiteralRefinedFromURI returns the first owned literal that is refined from the given URI
func (cPtr *Concept) GetFirstOwnedLiteralRefinedFromURI(uri string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	abstraction := cPtr.uOfD.GetElementWithURI(uri)
	if abstraction != nil {
		return cPtr.GetFirstOwnedLiteralRefinedFrom(abstraction, trans)
	}
	return nil
}

// GetFirstOwnedLiteralWithURI returns the first owned literal that has the indicated URI
func (cPtr *Concept) GetFirstOwnedLiteralWithURI(uri string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		switch element.GetConceptType() {
		case Literal:
			if element.GetURI(trans) == uri {
				it.Stop()
				return element
			}
		}
	}
	return nil
}

// GetFirstOwnedReferenceWithURI returns the first owned reference with the indicated URI
func (cPtr *Concept) GetFirstOwnedReferenceWithURI(uri string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		if element.GetURI(trans) == uri {
			it.Stop()
			return element
		}
	}
	return nil
}

// GetFirstOwnedRefinementWithURI returns the first owned refinement with the indicated URI
func (cPtr *Concept) GetFirstOwnedRefinementWithURI(uri string, trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		if element.GetURI(trans) == uri {
			it.Stop()
			return element
		}
	}
	return nil
}

// Deregister removes the registration of an Observer
func (cPtr *Concept) Deregister(observer Observer) error {
	cPtr.observers.Remove(observer)
	return nil
}

// FindAbstractions adds all found abstractions to supplied map
func (cPtr *Concept) FindAbstractions(abstractions map[string]*Concept, trans *Transaction) {
	it := cPtr.uOfD.listenersMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		listener := cPtr.uOfD.GetElement(id.(string))
		abstraction := listener.GetAbstractConcept(trans)
		if abstraction != nil && abstraction.getConceptIDNoLock() != cPtr.getConceptIDNoLock() {
			abstractions[abstraction.GetConceptID(trans)] = abstraction
			abstraction.FindAbstractions(abstractions, trans)
		}
	}
}

// FindImmediateAbstractions adds all immediate abstractions to supplied map
func (cPtr *Concept) FindImmediateAbstractions(abstractions map[string]*Concept, trans *Transaction) {
	// There are no abstractions without the uOfD context
	if cPtr.uOfD == nil {
		return
	}
	it := cPtr.uOfD.listenersMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		listener := cPtr.uOfD.GetElement(id.(string))
		switch listener.GetConceptType() {
		case Refinement:
			abstraction := listener.GetAbstractConcept(trans)
			if abstraction != nil && abstraction.getConceptIDNoLock() != cPtr.getConceptIDNoLock() {
				abstractions[abstraction.GetConceptID(trans)] = abstraction
			}
		}
	}
}

// GetIsCore returns true if the element is one of the core elements of CRL. The purpose of this
// function is to prevent SetReadOnly(true) on concepts that are built-in to CRL. Locking is
// not necessary as this value is set when the object is created and never expected to change
func (cPtr *Concept) GetIsCore(trans *Transaction) bool {
	trans.ReadLockElement(cPtr)
	return cPtr.IsCore
}

// GetLabel returns the label if one exists
func (cPtr *Concept) GetLabel(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.Label
}

func (cPtr *Concept) getLabelNoLock() string {
	return cPtr.Label
}

// GetLiteralValue returns the literal value
func (cPtr *Concept) GetLiteralValue(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.LiteralValue
}

// GetOwningConceptID returns the ID of the concept that owns this one (if any)
func (cPtr *Concept) GetOwningConceptID(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.OwningConceptID
}

// GetOwnedConceptIDs returns the set of IDs owned by this concept. Note that if this Element is not
// presently in a uOfD it returns the empty set
func (cPtr *Concept) GetOwnedConceptIDs(trans *Transaction) mapset.Set {
	if cPtr.uOfD == nil {
		return mapset.NewSet()
	}
	return cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID)
}

// GetOwnedConcepts returns the element's owned concepts if
func (cPtr *Concept) GetOwnedConcepts(trans *Transaction) map[string]*Concept {
	ownedConcepts := make(map[string]*Concept)
	if cPtr.uOfD == nil {
		return ownedConcepts
	}
	it := cPtr.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		if element != nil {
			ownedConcepts[id.(string)] = element
		}
	}
	return ownedConcepts
}

// GetOwnedConceptsRefinedFrom returns the owned concepts with the indicated abstraction as
// one of their abstractions.
func (cPtr *Concept) GetOwnedConceptsRefinedFrom(abstraction *Concept, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		if element.IsRefinementOf(abstraction, trans) {
			matches[element.GetConceptID(trans)] = element
		}
	}
	return matches
}

// GetOwnedConceptsRefinedFromURI returns the owned concepts that have the abstraction indicated
// by the URI as one of their abstractions.
func (cPtr *Concept) GetOwnedConceptsRefinedFromURI(abstractionURI string, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	abstraction := cPtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
		for id := range it.C {
			element := cPtr.uOfD.GetElement(id.(string))
			if element.IsRefinementOf(abstraction, trans) {
				matches[element.GetConceptID(trans)] = element
			}
		}
	}
	return matches
}

// GetOwnedDescendantsRefinedFrom returns the owned concepts with the indicated abstraction as
// one of their abstractions.
func (cPtr *Concept) GetOwnedDescendantsRefinedFrom(abstraction *Concept, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	if abstraction != nil {
		// it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
		descendantIDs := mapset.NewSet()
		cPtr.uOfD.GetConceptsOwnedConceptIDsRecursively(cPtr.ConceptID, descendantIDs, trans)
		it := descendantIDs.Iterator()
		for id := range it.C {
			element := cPtr.uOfD.GetElement(id.(string))
			if element.IsRefinementOf(abstraction, trans) {
				matches[element.GetConceptID(trans)] = element
			}
		}
	}
	return matches
}

// GetOwnedDescendantsRefinedFromURI returns the descendant concepts that have the indicated abstraction
// by the URI as one of their abstractions.
func (cPtr *Concept) GetOwnedDescendantsRefinedFromURI(abstractionURI string, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	abstraction := cPtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		// it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
		descendantIDs := mapset.NewSet()
		cPtr.uOfD.GetConceptsOwnedConceptIDsRecursively(cPtr.ConceptID, descendantIDs, trans)
		it := descendantIDs.Iterator()
		for id := range it.C {
			element := cPtr.uOfD.GetElement(id.(string))
			if element.IsRefinementOf(abstraction, trans) {
				matches[element.GetConceptID(trans)] = element
			}
		}
	}
	return matches
}

// GetOwnedLiteralsRefinedFrom returns the owned literals that have the indicated
// abstraction as one of their abstractions.
func (cPtr *Concept) GetOwnedLiteralsRefinedFrom(abstraction *Concept, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		switch element.GetConceptType() {
		case Literal:
			if element.IsRefinementOf(abstraction, trans) {
				matches[element.GetConceptID(trans)] = element
			}
		}
	}
	return matches
}

// GetOwnedLiteralsRefinedFromURI returns the child literals that have the abstraction indicated
// by the URI as one of their abstractions.
func (cPtr *Concept) GetOwnedLiteralsRefinedFromURI(abstractionURI string, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	abstraction := cPtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
		for id := range it.C {
			element := cPtr.uOfD.GetElement(id.(string))
			switch element.GetConceptType() {
			case Literal:
				if element.IsRefinementOf(abstraction, trans) {
					matches[element.GetConceptID(trans)] = element
				}
			}
		}
	}
	return matches
}

// GetOwnedReferencesRefinedFrom returns the owned references that have the indicated
// abstraction as one of their abstractions.
func (cPtr *Concept) GetOwnedReferencesRefinedFrom(abstraction *Concept, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		switch element.GetConceptType() {
		case Reference:
			if element.IsRefinementOf(abstraction, trans) {
				matches[element.GetConceptID(trans)] = element
			}
		}
	}
	return matches
}

// GetOwnedReferencesRefinedFromURI returns the owned references that have the abstraction indicated
// by the URI as one of their abstractions.
func (cPtr *Concept) GetOwnedReferencesRefinedFromURI(abstractionURI string, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	abstraction := cPtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
		defer it.Stop()
		for id := range it.C {
			element := cPtr.uOfD.GetElement(id.(string))
			switch element.GetConceptType() {
			case Reference:
				if element.IsRefinementOf(abstraction, trans) {
					matches[element.GetConceptID(trans)] = element
				}
			}
		}
	}
	return matches
}

// GetOwnedRefinementsRefinedFrom returns the owned refinements that have the indicated
// abstraction as one of their abstractions.
func (cPtr *Concept) GetOwnedRefinementsRefinedFrom(abstraction *Concept, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		element := cPtr.uOfD.GetElement(id.(string))
		switch element.GetConceptType() {
		case Refinement:
			if element.IsRefinementOf(abstraction, trans) {
				matches[element.GetConceptID(trans)] = element
			}
		}
	}
	return matches
}

// GetOwnedRefinementsRefinedFromURI returns the owned refinements that have the abstraction indicated
// by the URI as one of its abstractions.
func (cPtr *Concept) GetOwnedRefinementsRefinedFromURI(abstractionURI string, trans *Transaction) map[string]*Concept {
	trans.ReadLockElement(cPtr)
	matches := map[string]*Concept{}
	abstraction := cPtr.uOfD.GetElementWithURI(abstractionURI)
	if abstraction != nil {
		it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
		defer it.Stop()
		for id := range it.C {
			element := cPtr.uOfD.GetElement(id.(string))
			switch element.GetConceptType() {
			case Refinement:
				if element.IsRefinementOf(abstraction, trans) {
					matches[element.GetConceptID(trans)] = element
				}
			}
		}
	}
	return matches
}

// GetOwningConcept returns the Element representing the concept that owns this one (if any)
func (cPtr *Concept) GetOwningConcept(trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	if cPtr.uOfD != nil {
		return cPtr.uOfD.GetElement(cPtr.OwningConceptID)
	}
	return nil
}

// getOwningConceptNoLock returns the Element representing the concept that owns this one (if any)
func (cPtr *Concept) getOwningConceptNoLock() *Concept {
	if cPtr.uOfD != nil {
		return cPtr.uOfD.GetElement(cPtr.OwningConceptID)
	}
	return nil
}

// getOwningConceptIDNoLock returns the Element representing the concept that owns this one (if any)
func (cPtr *Concept) getOwningConceptIDNoLock() string {
	return cPtr.OwningConceptID
}

// GetReferencedConcept returns the element representing  the concept being referenced
// Note that this is a cached value
func (cPtr *Concept) GetReferencedConcept(trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	return cPtr.getReferencedConceptNoLock()
}

func (cPtr *Concept) getReferencedConceptNoLock() *Concept {
	return cPtr.uOfD.GetElement(cPtr.ReferencedConceptID)
}

// GetReferencedConceptID returns the identifier of the concept being referenced
func (cPtr *Concept) GetReferencedConceptID(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.ReferencedConceptID
}

// GetReferencedAttributeName returns an indicator of which attribute is being referenced (if any)
func (cPtr *Concept) GetReferencedAttributeName(trans *Transaction) AttributeName {
	trans.ReadLockElement(cPtr)
	return cPtr.ReferencedAttributeName
}

// GetReferencedAttributeValue returns the string value of the referenced attribute (if any)
func (cPtr *Concept) GetReferencedAttributeValue(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	if cPtr.ReferencedConceptID != "" {
		referencedConcept := cPtr.GetReferencedConcept(trans)
		if referencedConcept != nil {
			if cPtr.ReferencedAttributeName == OwningConceptID {
				return referencedConcept.GetOwningConceptID(trans)
			}
			switch referencedConcept.GetConceptType() {
			case Reference:
				if cPtr.ReferencedAttributeName == ReferencedConceptID {
					return referencedConcept.GetReferencedConceptID(trans)
				}
			case Refinement:
				if cPtr.ReferencedAttributeName == AbstractConceptID {
					return referencedConcept.GetAbstractConceptID(trans)
				}
				if cPtr.ReferencedAttributeName == RefinedConceptID {
					return referencedConcept.GetRefinedConceptID(trans)
				}
			case Literal:
				if cPtr.ReferencedAttributeName == LiteralValue {
					return referencedConcept.GetLiteralValue(trans)
				}
			}
		}
	}
	return ""
}

// GetRefinedConcept returns the refined concept which will be nil unless this is a Refinement
func (cPtr *Concept) GetRefinedConcept(trans *Transaction) *Concept {
	trans.ReadLockElement(cPtr)
	return cPtr.uOfD.GetElement(cPtr.RefinedConceptID)
}

func (cPtr *Concept) getRefinedConceptNoLock() *Concept {
	return cPtr.uOfD.GetElement(cPtr.RefinedConceptID)
}

// GetRefinedConceptID returns the is of the refined concept which will be nil unless this is a Refinement
func (cPtr *Concept) GetRefinedConceptID(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.RefinedConceptID
}

func (cPtr *Concept) getRefinedConceptIDNoLock() string {
	return cPtr.RefinedConceptID
}

// GetUniverseOfDiscourse returns the UniverseOfDiscourse in which the element instance resides
func (cPtr *Concept) GetUniverseOfDiscourse(trans *Transaction) *UniverseOfDiscourse {
	trans.ReadLockElement(cPtr)
	return cPtr.uOfD
}

// getUniverseOfDiscourseNoLock returns the UniverseOfDiscourse in which the element instance resides
func (cPtr *Concept) getUniverseOfDiscourseNoLock() *UniverseOfDiscourse {
	return cPtr.uOfD
}

// GetURI returns the URI string associated with the element if there is one
func (cPtr *Concept) GetURI(trans *Transaction) string {
	trans.ReadLockElement(cPtr)
	return cPtr.URI
}

// getURINoLock returns the URI string associated with the element if there is one
func (cPtr *Concept) getURINoLock() string {
	return cPtr.URI
}

// GetVersion returns the version of the element
func (cPtr *Concept) GetVersion(trans *Transaction) int {
	trans.ReadLockElement(cPtr)
	return cPtr.Version.getVersion()
}

// IsRefinementOf returns true if the given abstraction is contained in the abstractions set
// of this element. No locking is required since the StringIntMap does its own locking
func (cPtr *Concept) IsRefinementOf(abstraction *Concept, trans *Transaction) bool {
	trans.ReadLockElement(cPtr)
	// Check to see whether the abstraction is one of the core classes
	abstractionURI := abstraction.GetURI(trans)
	switch abstractionURI {
	case ElementURI:
		return true
	case LiteralURI:
		switch cPtr.GetConceptType() {
		case Literal:
			return true
		}
	case ReferenceURI:
		switch cPtr.GetConceptType() {
		case Reference:
			return true
		}
	case RefinementURI:
		switch cPtr.GetConceptType() {
		case Refinement:
			return true
		}
	}
	it := trans.uOfD.listenersMap.GetMappedValues(cPtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		listener := cPtr.uOfD.GetElement(id.(string))
		switch listener.GetConceptType() {
		case Refinement:
			foundAbstraction := listener.GetAbstractConcept(trans)
			if foundAbstraction == nil {
				continue
			}
			if foundAbstraction.getConceptIDNoLock() == cPtr.ConceptID {
				continue
			}
			if foundAbstraction == abstraction {
				return true
			}
			if foundAbstraction != nil {
				foundRecursively := foundAbstraction.IsRefinementOf(abstraction, trans)
				if foundRecursively {
					return true
				}
			}
		}
	}
	return false
}

// IsRefinementOfURI returns true if this is a refinement of the indicated URI
func (cPtr *Concept) IsRefinementOfURI(uri string, trans *Transaction) bool {
	trans.ReadLockElement(cPtr)
	if cPtr.uOfD == nil {
		return false
	}
	abstraction := cPtr.uOfD.GetElementWithURI(uri)
	if abstraction == nil {
		return false
	}
	return cPtr.IsRefinementOf(abstraction, trans)
}

// initializeConcept creates the identifier (using the uri if supplied) and
// creates the abstractions, ownedConcepts, and referrencingConcpsts maps.
// Note that initialization is not considered a change, so the version counter is not incremented
// nor are monitors of this element notified of changes.
func (cPtr *Concept) initializeConcept(conceptType ConceptType, identifier string, uri string) {
	cPtr.ConceptID = identifier
	cPtr.ConceptType = conceptType
	cPtr.Version = newVersionCounter()
	cPtr.URI = uri
	cPtr.observers = mapset.NewSet()
}

// IsReadOnly returns a boolean indicating whether the concept can be modified.
func (cPtr *Concept) IsReadOnly(trans *Transaction) bool {
	trans.ReadLockElement(cPtr)
	return cPtr.ReadOnly
}

// isEditable checks to see if the element cannot be edited because it
// is either a core element or has been marked readOnly.
func (cPtr *Concept) isEditable(trans *Transaction) bool {
	if cPtr.GetIsCore(trans) || cPtr.IsReadOnly(trans) {
		return false
	}
	return true
}

// isEquivalent only checks the element attributes. It ignores the uOfD.
func (cPtr *Concept) isEquivalent(hl1 *Transaction, el *Concept, hl2 *Transaction, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) > 0 {
		print = printExceptions[0]
	}
	hl1.ReadLockElement(cPtr)
	hl2.ReadLockElement(el)
	if cPtr.AbstractConceptID != el.AbstractConceptID {
		if print {
			log.Printf("In refinement.isEquivalent, AbstractConecptIDs do not match")
		}
		return false
	}
	if cPtr.ConceptID != el.ConceptID {
		if print {
			log.Printf("In element.isEquivalent, ConceptIDs do not match")
		}
		return false
	}
	if cPtr.Definition != el.Definition {
		if print {
			log.Printf("In element.isEquivalent, Definitions do not match")
		}
		return false
	}
	if cPtr.IsCore != el.IsCore {
		if print {
			log.Printf("In element.isEquivalent, IsCore do not match")
		}
		return false
	}
	if cPtr.Label != el.Label {
		if print {
			log.Printf("In element.isEquivalent, Labels do not match")
		}
		return false
	}
	if cPtr.LiteralValue != el.LiteralValue {
		if print {
			log.Printf("In literal.isEquivalent, LiteralValues do not match")
		}
		return false
	}
	if cPtr.OwningConceptID != el.OwningConceptID {
		if print {
			log.Printf("In element.isEquivalent, OwningConceptIDs do not match")
		}
		return false
	}
	if cPtr.ReadOnly != el.ReadOnly {
		if print {
			log.Printf("In element.isEquivalent, ReadOnly does not match")
		}
		return false
	}
	if cPtr.ReferencedConceptID != el.ReferencedConceptID {
		if print {
			log.Printf("In reference.IsEquivalent, ReferencedConceptIDs do not match")
		}
		return false
	}
	if cPtr.ReferencedAttributeName != el.ReferencedAttributeName {
		if print {
			log.Printf("In reference.IsEquivalent, ReferencedAttributeNames do not match")
		}
		return false
	}
	if cPtr.RefinedConceptID != el.RefinedConceptID {
		if print {
			log.Printf("In refinement.isEquivalent, RefinedConecptIDs do not match")
		}
		return false
	}
	if cPtr.Version.getVersion() != el.Version.getVersion() {
		if print {
			log.Printf("In element.isEquivalent, Versions do not match")
		}
		return false
	}
	if cPtr.URI != el.URI {
		if print {
			log.Printf("In element.isEquivalent, URIs do not match")
		}
		return false
	}
	return true
}

// IsLiteral() returns true if the concept is a literal
func (cPtr *Concept) IsLiteral() bool {
	return cPtr.ConceptType == Literal
}

// IsOwnedConcept returns true if the supplied element is an owned concept. Note that
// there is an interval of time during editing in which the child's owner will be set but the child
// has not yet been added to the element's OwnedConcepts list. Similarly, there is an interval of time
// during editing during which the child's owner has been changed but the original owner's OwnedConcept
// list has not yet been updated.
func (cPtr *Concept) IsOwnedConcept(el *Concept, trans *Transaction) bool {
	trans.ReadLockElement(cPtr)
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	defer it.Stop()
	for id := range it.C {
		child := cPtr.uOfD.GetElement(id.(string))
		if el.GetConceptID(trans) == child.GetConceptID(trans) {
			return true
		}
	}
	return false
}

// IsReference() returns true if the concept is a reference
func (cPtr *Concept) IsReference() bool {
	return cPtr.ConceptType == Reference
}

// IsRefinement() returns true if the concept is a refinement
func (cPtr *Concept) IsRefinement() bool {
	return cPtr.ConceptType == Refinement
}

// MarshalJSON produces a byte string JSON representation of the Element
func (cPtr *Concept) MarshalJSON() ([]byte, error) {
	type AliasConcept Concept
	return json.Marshal(&struct {
		ConceptType string
		*AliasConcept
		IsCore                  string
		ReadOnly                string
		ReferencedAttributeName string
		Version                 string
	}{
		ConceptType:             ConceptTypeToString(cPtr.ConceptType),
		IsCore:                  strconv.FormatBool(cPtr.IsCore),
		ReadOnly:                strconv.FormatBool(cPtr.ReadOnly),
		ReferencedAttributeName: cPtr.ReferencedAttributeName.String(),
		Version:                 strconv.Itoa(cPtr.Version.getVersion()),
		AliasConcept:            (*AliasConcept)(cPtr),
	})
}

// notifyObservers passes the notification to all registered Observers
func (cPtr *Concept) notifyObservers(notification *ChangeNotification, trans *Transaction) error {
	// it := cPtr.observers.Iterator()
	// defer it.Stop()
	for _, observer := range cPtr.observers.ToSlice() {
		err := observer.(Observer).Update(notification, trans)
		if err != nil {
			return errors.Wrap(err, "element.notifyObservers failed")
		}
	}
	return nil
}

func (cPtr *Concept) notifyPointerOwners(notification *ChangeNotification, trans *Transaction) error {
	trans.ReadLockElement(cPtr)
	uOfD := cPtr.uOfD
	if uOfD != nil {
		it := uOfD.listenersMap.GetMappedValues(cPtr.ConceptID).Iterator()
		for id := range it.C {
			listener := uOfD.GetElement(id.(string))
			indicatedConceptChangeNotification, err := uOfD.NewForwardingChangeNotification(listener, IndicatedConceptChanged, notification, trans)
			if err != nil {
				return errors.Wrap(err, "element.notifyPointerOwners failed")
			}
			err = uOfD.callAssociatedFunctions(listener, indicatedConceptChangeNotification, trans)
			if err != nil {
				it.Stop()
				return errors.Wrap(err, "element.notifyPointerOwners failed")
			}
			err = listener.notifyOwner(indicatedConceptChangeNotification, trans)
			if err != nil {
				it.Stop()
				return errors.Wrap(err, "element.notifyPointerOwners failed")
			}
			err = listener.notifyObservers(indicatedConceptChangeNotification, trans)
			if err != nil {
				it.Stop()
				return errors.Wrap(err, "element.notifyPointerOwners failed")
			}
		}
	}
	return nil
}

// notifyOwner informs the owner that the concept has changed state
func (cPtr *Concept) notifyOwner(notification *ChangeNotification, trans *Transaction) error {
	trans.ReadLockElement(cPtr)
	switch notification.natureOfChange {
	case OwningConceptChanged:
		oldOwnerID := notification.beforeConceptState.OwningConceptID
		newOwnerID := notification.afterConceptState.OwningConceptID
		if oldOwnerID != "" {
			oldOwner := cPtr.uOfD.GetElement(oldOwnerID)
			if oldOwner != nil {
				ownedConceptChangeNotification, err := cPtr.uOfD.NewForwardingChangeNotification(oldOwner, OwnedConceptChanged, notification, trans)
				if err != nil {
					return errors.Wrap(err, "element.notifyOwner failed")
				}
				err = cPtr.uOfD.callAssociatedFunctions(oldOwner, ownedConceptChangeNotification, trans)
				if err != nil {
					return errors.Wrap(err, "element.notifyOwner failed")
				}
				err = oldOwner.notifyObservers(ownedConceptChangeNotification, trans)
				if err != nil {
					return errors.Wrap(err, "element.notifyOwner failed")
				}
			}
		}
		if newOwnerID != "" {
			newOwner := cPtr.uOfD.GetElement(newOwnerID)
			if newOwner != nil {
				ownedConceptChangeNotification, err := cPtr.uOfD.NewForwardingChangeNotification(newOwner, OwnedConceptChanged, notification, trans)
				if err != nil {
					return errors.Wrap(err, "element.notifyOwner failed")
				}
				err = cPtr.uOfD.callAssociatedFunctions(newOwner, ownedConceptChangeNotification, trans)
				if err != nil {
					return errors.Wrap(err, "element.notifyOwner failed")
				}
				err = newOwner.notifyObservers(ownedConceptChangeNotification, trans)
				if err != nil {
					return errors.Wrap(err, "element.notifyOwner failed")
				}
			}
		}
	case ConceptChanged, ReferencedConceptChanged, AbstractConceptChanged, RefinedConceptChanged, IndicatedConceptChanged:
		owner := cPtr.GetOwningConcept(trans)
		if owner != nil {
			ownedConceptChangeNotification, err := cPtr.uOfD.NewForwardingChangeNotification(owner, OwnedConceptChanged, notification, trans)
			if err != nil {
				return errors.Wrap(err, "element.notifyOwner failed")
			}
			err = cPtr.uOfD.callAssociatedFunctions(owner, ownedConceptChangeNotification, trans)
			if err != nil {
				return errors.Wrap(err, "element.notifyOwner failed")
			}
			err = owner.notifyObservers(ownedConceptChangeNotification, trans)
			if err != nil {
				return errors.Wrap(err, "element.notifyOwner failed")
			}
		}
	}
	return nil
}

// propagateChange() distributes the change notification to relevant parties
func (cPtr *Concept) propagateChange(notification *ChangeNotification, trans *Transaction) error {
	var err error = nil
	switch notification.natureOfChange {
	case ConceptChanged, OwningConceptChanged, OwnedConceptChanged, ReferencedConceptChanged, AbstractConceptChanged, RefinedConceptChanged:
		err = notification.uOfD.callAssociatedFunctions(cPtr, notification, trans)
		if err != nil {
			return errors.Wrap(err, "element.propagateChange failed")
		}
		err = cPtr.notifyPointerOwners(notification, trans)
		if err != nil {
			return errors.Wrap(err, "element.propagateChange failed")
		}
		if cPtr.uOfD != nil {
			err = cPtr.notifyOwner(notification, trans)
		}
		if err != nil {
			return errors.Wrap(err, "element.propagateChange failed")
		}
		err = cPtr.notifyObservers(notification, trans)
		if err != nil {
			return errors.Wrap(err, "element.propagateChange failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.NotifyUofDObservers(notification, trans)
		}
		if err != nil {
			return errors.Wrap(err, "element.propagateChange failed")
		}
	case ConceptAdded, ConceptRemoved:
		if cPtr.uOfD != nil {
			cPtr.uOfD.NotifyUofDObservers(notification, trans)
		}
		if err != nil {
			return errors.Wrap(err, "element.propagateChange failed")
		}
	}
	return nil
}

// tickle sends the notification to the targetElement. Its sole purpose is to trigger any functions
// that may be associated with that Element.
func (cPtr *Concept) tickle(targetElement *Concept, notification *ChangeNotification, trans *Transaction) error {
	var err error = nil
	switch notification.natureOfChange {
	case Tickle:
		err = notification.uOfD.callAssociatedFunctions(targetElement, notification, trans)
		if err != nil {
			return errors.Wrap(err, "element.tickle failed")
		}
		err = cPtr.notifyObservers(notification, trans)
		if err != nil {
			return errors.Wrap(err, "element.trigger failed")
		}
	}
	return nil
}

// removeListener removes the indicated Element as a listening concept.
func (cPtr *Concept) removeListener(listeningConceptID string, trans *Transaction) {
	trans.ReadLockElement(cPtr)
	if cPtr.uOfD.undoManager.debugUndo {
		log.Print("+++")
		log.Print("+++ removeListener")
		log.Print("+++")
	}
	cPtr.uOfD.preChange(cPtr, trans)
	cPtr.uOfD.listenersMap.removeMappedValue(cPtr.ConceptID, listeningConceptID)
	if cPtr.uOfD != nil {
		cPtr.uOfD.postChange(cPtr, trans)
	}
}

// Register adds the registration of an Observer
func (cPtr *Concept) Register(observer Observer) error {
	cPtr.observers.Add(observer)
	return nil
}

// removeOwnedConcept removes the indicated Element as a child (owned) concept.
func (cPtr *Concept) removeOwnedConcept(ownedConceptID string, trans *Transaction) error {
	trans.ReadLockElement(cPtr)
	if cPtr.IsReadOnly(trans) {
		return errors.New("Element.removedOwnedConcept called on read-only Element")
	}
	if cPtr.uOfD.undoManager.debugUndo {
		log.Print("+++")
		log.Print("+++ removeOwnedConcept")
		log.Print("+++")
	}
	cPtr.uOfD.preChange(cPtr, trans)
	cPtr.Version.incrementVersion()
	cPtr.uOfD.ownedIDsMap.removeMappedValue(cPtr.ConceptID, ownedConceptID)
	if cPtr.uOfD != nil {
		cPtr.uOfD.postChange(cPtr, trans)
	}
	return nil
}

// SetAbstractConcept sets the abstract concept using the ID of the supplied Element
func (cPtr *Concept) SetAbstractConcept(el *Concept, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("refinement.SetAbstractConcept failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return cPtr.SetAbstractConceptID(id, trans)
}

// SetAbstractConceptID sets the abstract concept ID
func (cPtr *Concept) SetAbstractConceptID(acID string, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("refinement.SetAbstractConceptID failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.isEditable(trans) {
		return errors.New("refinement.SetAbstractConceptID failed because the refinement is not editable")
	}
	if cPtr.AbstractConceptID != acID {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetAbstractConcepetID")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
		}
		cPtr.Version.incrementVersion()
		var oldAbstractConcept *Concept
		if cPtr.AbstractConceptID != "" {
			oldAbstractConcept = cPtr.uOfD.GetElement(cPtr.AbstractConceptID)
			if oldAbstractConcept != nil {
				oldAbstractConcept.removeListener(cPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
				}
			} else {
				// This case can arise if the abstract concept is not currently loaded
				cPtr.uOfD.listenersMap.removeMappedValue(cPtr.AbstractConceptID, cPtr.ConceptID)
			}
		}
		var newAbstractConcept *Concept
		if acID != "" {
			newAbstractConcept = cPtr.uOfD.GetElement(acID)
			if newAbstractConcept != nil {
				newAbstractConcept.addListener(cPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
				}
			}
		}
		cPtr.AbstractConceptID = acID
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "refinement.SetAbstractConceptID failed")
		}
		err = cPtr.uOfD.SendPointerChangeNotification(cPtr, AbstractConceptChanged, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "refinement.SetAbstractConceptID failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// SetDefinition sets the definition of the Element
func (cPtr *Concept) SetDefinition(def string, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetDefinition failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.isEditable(trans) {
		return errors.New("element.SetDefinition failed because the element is not editable")
	}
	if cPtr.Definition != def {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetDefinition")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "element.SetDefinition failed")
		}
		cPtr.Version.incrementVersion()
		cPtr.Definition = def
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetDefinition failed")
		}
		err = cPtr.uOfD.SendConceptChangeNotification(cPtr, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "element.SetDefinition failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// SetIsCore sets the flag indicating that the element is a Core concept and cannot be edited. Once set, this flag cannot be cleared.
func (cPtr *Concept) SetIsCore(trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetIsCore failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.IsCore {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetIsCore")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "element.SetIsCore failed")
		}
		cPtr.Version.incrementVersion()
		cPtr.IsCore = true
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetIsCore failed")
		}
		err = cPtr.uOfD.SendConceptChangeNotification(cPtr, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "element.SetIsCore failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// SetIsCoreRecursively recursively sets the flag indicating that the element is a Core concept and cannot be edited. Once set, this flag cannot be cleared.
func (cPtr *Concept) SetIsCoreRecursively(trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetIsCoreRecursively failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	err := cPtr.SetIsCore(trans)
	if err != nil {
		return errors.Wrap(err, "Element.SetIsCoreRecursively failed")
	}
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		el := cPtr.uOfD.GetElement(id.(string))
		err = el.SetIsCoreRecursively(trans)
		if err != nil {
			it.Stop()
			return errors.Wrap(err, "Element.SetIsCoreRecursively failed")
		}
	}
	return nil
}

// SetLabel sets the label of the Element
func (cPtr *Concept) SetLabel(label string, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetLabel failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.isEditable(trans) {
		return errors.New("element.SetLabel failed because the element is not editable")
	}
	if cPtr.Label != label {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetLabel")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "element.SetLabel failed")
		}
		cPtr.Version.incrementVersion()
		cPtr.Label = label
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetLabel failed")
		}
		err = cPtr.uOfD.SendConceptChangeNotification(cPtr, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "element.SetLabel failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// SetLiteralValue sets the literal value
func (cPtr *Concept) SetLiteralValue(value string, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("literal.SetLiteralValue failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.isEditable(trans) {
		return errors.New("literal.SetLiteralValue failed because the literal is not editable")
	}
	if cPtr.LiteralValue != value {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetLiteralValue")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "literal.SetLiteralValue failed")
		}
		cPtr.Version.incrementVersion()
		cPtr.LiteralValue = value
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "literal.SetLiteralValue failed")
		}
		err = cPtr.uOfD.SendConceptChangeNotification(cPtr, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "literal.SetLiteralValue failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// SetOwningConcept takes the ID of the supplied concept and call SetOwningConceptID. It first checks to
// determine whether the new owner is editable and will throw an error if it is not
func (cPtr *Concept) SetOwningConcept(el *Concept, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetOwningConcept failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	id := ""
	if el != nil {
		if !el.isEditable(trans) {
			return errors.New("element.SetOwningConcept called with an owner that is not editable")
		}
		id = el.getConceptIDNoLock()
	}
	err := cPtr.SetOwningConceptID(id, trans)
	if err != nil {
		errors.Wrap(err, "element.SetOwningConcept failed")
	}
	return nil
}

// SetOwningConceptID sets the ID of the owning concept for the element
// Design Note: the argument is the identifier rather than the Element to ensure
// the correct type of the owning concept is recorded.
func (cPtr *Concept) SetOwningConceptID(ocID string, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetOwningConceptID failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.isEditable(trans) {
		return errors.New("element.SetOwningConceptID failed because the element is not editable")
	}
	if ocID == cPtr.ConceptID {
		return errors.New("element.SetOwningConceptID called with itself as owner")
	}
	newOwner := cPtr.uOfD.GetElement(ocID)
	if newOwner != nil && !newOwner.isEditable(trans) {
		return errors.New("element.SetOwningConceptID called with new owner not editable")
	}
	oldOwner := cPtr.GetOwningConcept(trans)
	if oldOwner != nil && !oldOwner.isEditable(trans) {
		return errors.New("element.SetOwningConceptID called with old owner not editable")
	}
	// Do nothing if there is no change
	if cPtr.OwningConceptID != ocID {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetOwningConceptID")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "element.SetOwningConceptID failed")
		}
		if oldOwner != nil {
			oldOwner.removeOwnedConcept(cPtr.ConceptID, trans)
			if err != nil {
				return errors.Wrap(err, "element.SetOwningConceptID failed")
			}
		}
		cPtr.Version.incrementVersion()
		if newOwner != nil {
			newOwner.addOwnedConcept(cPtr.ConceptID, trans)
			if err != nil {
				return errors.Wrap(err, "element.SetOwningConceptID failed")
			}
		}
		cPtr.OwningConceptID = ocID
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetOwningConceptID failed")
		}
		err = cPtr.uOfD.SendPointerChangeNotification(cPtr, OwningConceptChanged, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "element.SetOwningConceptID failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// SetReadOnly provides a mechanism for preventing modifications to concepts. It will throw an error
// if the concept is one of the CRL core concepts, as these can never be made writable. It will also
// throw an error if its owner is read only and this call tries to set read only false.
func (cPtr *Concept) SetReadOnly(value bool, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetReadOnly failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if cPtr.GetIsCore(trans) {
		return errors.New("element.SetReadOnly failed because element is a core element")
	}
	if cPtr.GetOwningConcept(trans) != nil {
		if cPtr.GetOwningConcept(trans).IsReadOnly(trans) && !value {
			return errors.New("element.SetReadOnly failed because the owner is read only")
		}
	}
	if cPtr.ReadOnly != value {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetReadOnly")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "element.SetReadOnly failed")
		}
		cPtr.Version.incrementVersion()
		cPtr.ReadOnly = value
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetDeSetReadOnlyfinition failed")
		}
		err = cPtr.uOfD.SendConceptChangeNotification(cPtr, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "element.SetDeSetReadOnlyfinition failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// SetReadOnlyRecursively makes this concept and all its descendants read-only
func (cPtr *Concept) SetReadOnlyRecursively(value bool, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetReadOnlyRecursively failed because the element uOfD is nil")
	}
	err := cPtr.SetReadOnly(value, trans)
	if err != nil {
		return errors.Wrap(err, "Element.SetReadOnlyRecursively failed")
	}
	it := cPtr.uOfD.ownedIDsMap.GetMappedValues(cPtr.ConceptID).Iterator()
	for id := range it.C {
		el := cPtr.uOfD.GetElement(id.(string))
		err = el.SetReadOnlyRecursively(value, trans)
		if err != nil {
			it.Stop()
			return errors.Wrap(err, "Element.SetReadOnlyRecursively failed")
		}
	}
	return nil
}

// SetReferencedConcept sets the referenced concept by calling SetReferencedConceptID using the ID of the
// supplied Element
func (cPtr *Concept) SetReferencedConcept(el *Concept, attributeName AttributeName, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("reference.SetReferencedConcept failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return cPtr.SetReferencedConceptID(id, attributeName, trans)
}

// SetReferencedConceptID sets the referenced concept using the supplied ID.
func (cPtr *Concept) SetReferencedConceptID(rcID string, attributeName AttributeName, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("reference.SetReferencedConceptID failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.isEditable(trans) {
		return errors.New("reference.SetReferencedConceptID failed because the reference is not editable")
	}
	var newReferencedConcept *Concept
	var oldReferencedConcept *Concept
	if cPtr.ReferencedConceptID != rcID || cPtr.ReferencedAttributeName != attributeName {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetReferencedConceptID")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		if rcID != "" {
			newReferencedConcept = cPtr.uOfD.GetElement(rcID)
			switch attributeName {
			case ReferencedConceptID:
				switch newReferencedConcept.GetConceptType() {
				case Reference:
				default:
					return errors.New("In reference.SetReferencedConceptID, the ReferencedAttributeName was ReferencedConceptID, but the referenced concept is not a Reference")
				}
			case AbstractConceptID, RefinedConceptID:
				switch newReferencedConcept.GetConceptType() {
				case Refinement:
				default:
					return errors.New("In reference.SetReferencedConceptID, the ReferencedAttributeName was AbstractConceptID or RefinedConceptID, but the referenced concept is not a Refinement")
				}
			}
			if newReferencedConcept != nil {
				newReferencedConcept.addListener(cPtr.ConceptID, trans)
			}
		}
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "reference.SetReferencedConceptID failed")
		}
		cPtr.Version.incrementVersion()
		if cPtr.ReferencedConceptID != "" {
			oldReferencedConcept = cPtr.uOfD.GetElement(cPtr.ReferencedConceptID)
			if oldReferencedConcept != nil {
				oldReferencedConcept.removeListener(cPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "reference.SetReferencedConceptID failed")
				}
			}
		}
		if rcID != "" {
			if newReferencedConcept != nil {
				newReferencedConcept.addListener(cPtr.ConceptID, trans)
			}
		}
		cPtr.ReferencedConceptID = rcID
		cPtr.ReferencedAttributeName = attributeName
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "reference.SetReferencedConceptID failed")
		}
		err = cPtr.uOfD.SendPointerChangeNotification(cPtr, ReferencedConceptChanged, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "reference.SetReferencedConceptID failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// SetRefinedConcept sets the refined concept
func (cPtr *Concept) SetRefinedConcept(el *Concept, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("refinement.SetRefinedConcept failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	id := ""
	if el != nil {
		id = el.getConceptIDNoLock()
	}
	return cPtr.SetRefinedConceptID(id, trans)
}

// SetRefinedConceptID sets the refined concept ID
func (cPtr *Concept) SetRefinedConceptID(rcID string, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("refinement.SetRefinedConceptID failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.isEditable(trans) {
		return errors.New("refinement.SetReferencedConceptID failed because the refinement is not editable")
	}
	if cPtr.RefinedConceptID != rcID {
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetRefinedConceptID")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
		}
		cPtr.Version.incrementVersion()
		var oldRefinedConcept *Concept
		if cPtr.RefinedConceptID != "" {
			oldRefinedConcept = cPtr.uOfD.GetElement(cPtr.RefinedConceptID)
			if oldRefinedConcept != nil {
				oldRefinedConcept.removeListener(cPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
				}
			} else {
				// This case can arise if the abstract concept is not currently loaded
				cPtr.uOfD.listenersMap.removeMappedValue(cPtr.RefinedConceptID, cPtr.ConceptID)
			}
		}
		var newRefinedConcept *Concept
		if rcID != "" {
			newRefinedConcept = cPtr.uOfD.GetElement(rcID)
			if newRefinedConcept != nil {
				newRefinedConcept.addListener(cPtr.ConceptID, trans)
				if err != nil {
					return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
				}
			}
		}
		cPtr.RefinedConceptID = rcID
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "refinement.SetRefinedConceptID failed")
		}
		err = cPtr.uOfD.SendPointerChangeNotification(cPtr, RefinedConceptChanged, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "refinement.SetRefinedConceptID failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// setUniverseOfDiscourse is intended to be called only by the UniverseOfDiscourse
func (cPtr *Concept) setUniverseOfDiscourse(uOfD *UniverseOfDiscourse, trans *Transaction) {
	trans.WriteLockElement(cPtr)
	cPtr.uOfD = uOfD
}

// SetURI sets the URI of the Element
func (cPtr *Concept) SetURI(uri string, trans *Transaction) error {
	if cPtr.uOfD == nil {
		return errors.New("element.SetURI failed because the element uOfD is nil")
	}
	trans.WriteLockElement(cPtr)
	if !cPtr.isEditable(trans) {
		return errors.New("element.SetURI failed because the elementis not editable")
	}
	if cPtr.URI != uri {
		foundElement := cPtr.uOfD.GetElementWithURI(uri)
		if foundElement != nil && foundElement.GetConceptID(trans) != cPtr.ConceptID {
			return errors.New("Element already exists with URI " + uri)
		}
		if cPtr.uOfD.undoManager.debugUndo {
			log.Print("+++")
			log.Print("+++ SetURI")
			log.Print("+++")
		}
		cPtr.uOfD.preChange(cPtr, trans)
		beforeState, err := NewConceptState(cPtr)
		if err != nil {
			return errors.Wrap(err, "element.SetURI failed")
		}
		cPtr.uOfD.changeURIForElement(cPtr, cPtr.URI, uri)
		cPtr.Version.incrementVersion()
		cPtr.URI = uri
		afterState, err2 := NewConceptState(cPtr)
		if err2 != nil {
			return errors.Wrap(err2, "element.SetURI failed")
		}
		err = cPtr.uOfD.SendConceptChangeNotification(cPtr, beforeState, afterState, trans)
		if err != nil {
			return errors.Wrap(err, "element.SetURI failed")
		}
		if cPtr.uOfD != nil {
			cPtr.uOfD.postChange(cPtr, trans)
		}
	}
	return nil
}

// TraceableReadLock locks the concept in a traceable way
func (cPtr *Concept) TraceableReadLock(trans *Transaction) {
	if TraceLocks {
		log.Printf("HL %p about to read lock Element %p %s\n", trans, cPtr, cPtr.Label)
	}
	cPtr.RLock()
}

// TraceableWriteLock locks the concept in a traceable way
func (cPtr *Concept) TraceableWriteLock(trans *Transaction) {
	if TraceLocks {
		log.Printf("HL %p about to write lock Element %p %s\n", trans, cPtr, cPtr.Label)
	}
	cPtr.Lock()
}

// TraceableReadUnlock unlocks the concept in a traceable way
func (cPtr *Concept) TraceableReadUnlock(trans *Transaction) {
	if TraceLocks {
		log.Printf("HL %p about to read unlock Element %p %s\n", trans, cPtr, cPtr.Label)
	}
	cPtr.RUnlock()
}

// TraceableWriteUnlock unlocks the concept in a traceable way
func (cPtr *Concept) TraceableWriteUnlock(trans *Transaction) {
	if TraceLocks {
		log.Printf("HL %p about to write unlock Element %p %s\n", trans, cPtr, cPtr.Label)
	}
	cPtr.Unlock()
}

// UnmarshalJSON unmarshals the JSON into the current concept
func (cPtr *Concept) UnmarshalJSON(data []byte) error {
	type AliasConcept Concept
	aux := &struct {
		ConceptType             string
		IsCore                  string
		ReadOnly                string
		ReferencedAttributeName string
		Version                 string
		*AliasConcept
	}{
		AliasConcept: (*AliasConcept)(cPtr),
	}
	var err error
	if err = json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "concept.UnmarshalJSON failed")
	}
	cPtr.ConceptType, err = StringToConceptType(aux.ConceptType)
	if err = json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "concept.UnmarshalJSON failed")
	}
	cPtr.IsCore, err = strconv.ParseBool(aux.IsCore)
	if err = json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "concept.UnmarshalJSON failed")
	}
	cPtr.ReadOnly, err = strconv.ParseBool(aux.ReadOnly)
	if err = json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "concept.UnmarshalJSON failed")
	}
	cPtr.ReferencedAttributeName, err = FindAttributeName(aux.ReferencedAttributeName)
	if err = json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "concept.UnmarshalJSON failed")
	}
	cPtr.Version.counter, err = strconv.Atoi(aux.Version)
	if err = json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "concept.UnmarshalJSON failed")
	}
	return nil
}

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
