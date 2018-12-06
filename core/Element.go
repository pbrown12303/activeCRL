package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"sync"
)

// element is the root representation of a concept
type element struct {
	sync.RWMutex
	ConceptID       string
	Definition      string
	Label           string
	listeners       *StringElementMap
	IsCore          bool
	ownedConcepts   *StringElementMap
	owningConcept   *cachedPointer
	OwningConceptID string
	readOnly        bool
	Version         *versionCounter
	uOfD            UniverseOfDiscourse
	URI             string
}

// addOwnedConcept adds the indicated Element as a child (owned) concept.
// This is purely an internal housekeeping method. Note that
// no checking of whether the Element is read-only is performed here. This check
// is performed by the child
func (ePtr *element) addOwnedConcept(ownedConceptID string, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	if ePtr.ownedConcepts.GetEntry(ownedConceptID) == nil {
		ePtr.uOfD.preChange(ePtr, hl)
		ePtr.incrementVersion(hl)
		ownedConcept := ePtr.GetUniverseOfDiscourse(hl).GetElement(ownedConceptID)
		if ownedConcept != nil {
			ePtr.ownedConcepts.SetEntry(ownedConceptID, ownedConcept)
		}
	}
}

// addRecoveredOwnedConcept adds the indicated Element as a child (owned) concept without incrementing
// the version.
// This is purely an internal housekeeping method. Note that
// no checking of whether the Element is read-only is performed here. This check
// is performed by the child
func (ePtr *element) addRecoveredOwnedConcept(ownedConceptID string, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	if ePtr.ownedConcepts.GetEntry(ownedConceptID) == nil {
		ePtr.uOfD.preChange(ePtr, hl)
		ownedConcept := ePtr.GetUniverseOfDiscourse(hl).GetElement(ownedConceptID)
		if ownedConcept != nil {
			ePtr.ownedConcepts.SetEntry(ownedConceptID, ownedConcept)
		}
	}
}

// addListener adds the indicated Element as a listening concept.
// This is an internal housekeeping method.
func (ePtr *element) addListener(listeningConceptID string, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	if ePtr.listeners.GetEntry(listeningConceptID) == nil {
		ePtr.uOfD.preChange(ePtr, hl)
		listeningConcept := ePtr.GetUniverseOfDiscourse(hl).GetElement(listeningConceptID)
		if listeningConcept != nil {
			ePtr.listeners.SetEntry(listeningConceptID, listeningConcept)
		}
	}
}

// clone is an internal function that makes a copy of the given element - including its
// identifier. This is done only to support undo/redo: the clone should NEVER be added to the
// universe of discourse
func (ePtr *element) clone(hl *HeldLocks) *element {
	hl.ReadLockElement(ePtr)
	// The newly made clone never gets locked
	var cl element
	cl.initializeElement("")
	cl.cloneAttributes(ePtr, hl)
	return &cl
}

// cloneAttributes is a supporting function for clone
func (ePtr *element) cloneAttributes(source *element, hl *HeldLocks) {
	ePtr.ConceptID = source.ConceptID
	//	ePtr.ownedConcepts = NewStringElementMap()
	for k, v := range *source.ownedConcepts.CopyMap() {
		ePtr.ownedConcepts.SetEntry(k, v)
	}
	ePtr.OwningConceptID = source.OwningConceptID
	//	ePtr.owningConcept = newCachedPointer(ePtr.getConceptIDNoLock(), true)
	ePtr.owningConcept.setIndicatedConceptID(source.owningConcept.getIndicatedConceptID())
	ePtr.owningConcept.setIndicatedConcept(source.owningConcept.getIndicatedConcept())
	ePtr.owningConcept.parentConceptID = source.owningConcept.parentConceptID
	ePtr.listeners = NewStringElementMap()
	for k, v := range *source.listeners.CopyMap() {
		ePtr.listeners.SetEntry(k, v)
	}
	//	ePtr.Version = newVersionCounter()
	ePtr.Version.counter = source.Version.counter
	ePtr.readOnly = source.readOnly
	ePtr.uOfD = source.uOfD
}

// editableError checks to see if the element cannot be edited because it
// is either a core element or has been marked readOnly.
func (ePtr *element) editableError() error {
	if ePtr.GetIsCore() {
		return errors.New("Element.SetOwningConceptID called on core Element")
	}
	if ePtr.readOnly {
		return errors.New("Element.SetOwningConcept called on read-only Element")
	}
	return nil
}

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

// GetFirstChildWithAbstraction returns the first child that has the indicated abstraction as
// one of its abstractions. Note that there is no ordering of children so in the event that
// there is more than one child with the given abstraction the result is nondeterministic.
func (ePtr *element) GetFirstChildWithAbstraction(abstraction Element, hl *HeldLocks) Element {
	hl.ReadLockElement(ePtr)
	for _, element := range *ePtr.ownedConcepts.CopyMap() {
		if element.HasAbstraction(abstraction, hl) {
			return element
		}
	}
	return nil
}

// FindAbstractions adds all found abstractions to supplied map
func (ePtr *element) FindAbstractions(abstractions *map[string]Element, hl *HeldLocks) {
	for _, listener := range *ePtr.listeners.CopyMap() {
		switch listener.(type) {
		case *refinement:
			abstraction := listener.(*refinement).GetAbstractConcept(hl)
			if abstraction != nil && abstraction.getConceptIDNoLock() != ePtr.getConceptIDNoLock() {
				(*abstractions)[abstraction.GetConceptID(hl)] = abstraction
				abstraction.FindAbstractions(abstractions, hl)
			}
		}
	}
}

// GetGetLabel returns the label if one exists
func (ePtr *element) GetLabel(hl *HeldLocks) string {
	hl.ReadLockElement(ePtr)
	return ePtr.Label
}

func (ePtr *element) GetOwnedConcepts(hl *HeldLocks) *map[string]Element {
	hl.ReadLockElement(ePtr)
	ownedConcepts := make(map[string]Element)
	for key, value := range *ePtr.ownedConcepts.CopyMap() {
		ownedConcepts[key] = value
	}
	return &ownedConcepts
}

// GetOwningConceptID returns the ID of the concept that owns this one (if any)
func (ePtr *element) GetOwningConceptID(hl *HeldLocks) string {
	hl.ReadLockElement(ePtr)
	return ePtr.OwningConceptID
}

// GetOwningConcept returns the Element representing the concept that owns this one (if any)
func (ePtr *element) GetOwningConcept(hl *HeldLocks) Element {
	hl.ReadLockElement(ePtr)
	return ePtr.owningConcept.getIndicatedConcept()
}

// GetUniverseOfDiscourse returns the UniverseOfDiscourse in which the element instance resides
func (ePtr *element) GetUniverseOfDiscourse(hl *HeldLocks) UniverseOfDiscourse {
	hl.ReadLockElement(ePtr)
	return ePtr.uOfD
}

// getUniverseOfDiscourseNoLock returns the UniverseOfDiscourse in which the element instance resides
func (ePtr *element) getUniverseOfDiscourseNoLock() UniverseOfDiscourse {
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

// HasAbstraction returns true if the given abstraction is contained in the abstractions set
// of this element. No locking is required since the StringIntMap does its own locking
func (ePtr *element) HasAbstraction(abstraction Element, hl *HeldLocks) bool {
	hl.ReadLockElement(ePtr)
	for _, listener := range *ePtr.listeners.CopyMap() {
		switch listener.(type) {
		case *refinement:
			foundAbstraction := listener.(*refinement).GetAbstractConcept(hl)
			if foundAbstraction == abstraction {
				return true
			}
			if foundAbstraction != nil {
				foundRecursively := foundAbstraction.HasAbstraction(abstraction, hl)
				if foundRecursively {
					return true
				}
			}
		}
	}
	return false
}

func (ePtr *element) incrementVersion(hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	ePtr.uOfD.preChange(ePtr, hl)
	ePtr.Version.incrementVersion()
	if ePtr.OwningConceptID != "" {
		ePtr.uOfD.GetElement(ePtr.OwningConceptID).incrementVersion(hl)
	}
}

// initializeElement creates the identifier (using the uri if supplied) and
// creates the abstractions, ownedConcepts, and referrencingConcpsts maps.
// Note that initialization is not considered a change, so the version counter is not incremented
// nor are monitors of this element notified of changes.
func (ePtr *element) initializeElement(identifier string) {
	ePtr.ConceptID = identifier
	ePtr.ownedConcepts = NewStringElementMap()
	ePtr.owningConcept = newCachedPointer(ePtr.getConceptIDNoLock(), true)
	ePtr.listeners = NewStringElementMap()
	ePtr.Version = newVersionCounter()
}

// GetIsCore returns true if the element is one of the core elements of CRL. The purpose of this
// function is to prevent SetReadOnly(true) on concepts that are built-in to CRL. Locking is
// not necessary as this value is set when the object is created and never expected to change
func (ePtr *element) GetIsCore() bool {
	return ePtr.IsCore
}

// IsReadOnly returns a boolean indicating whether the concept can be modified.
func (ePtr *element) IsReadOnly(hl *HeldLocks) bool {
	hl.ReadLockElement(ePtr)
	return ePtr.readOnly
}

// isEquivalent only checks the element attributes. It ignores the uOfD.
func (ePtr *element) isEquivalent(hl1 *HeldLocks, el *element, hl2 *HeldLocks) bool {
	hl1.ReadLockElement(ePtr)
	hl2.ReadLockElement(el)
	if ePtr.ConceptID != el.ConceptID {
		return false
	}
	if ePtr.Definition != el.Definition {
		return false
	}
	if ePtr.IsCore != el.IsCore {
		return false
	}
	if ePtr.Label != el.Label {
		return false
	}
	if ePtr.listeners.IsEquivalent(el.listeners) != true {
		return false
	}
	if ePtr.ownedConcepts.IsEquivalent(el.ownedConcepts) != true {
		return false
	}
	if ePtr.OwningConceptID != el.OwningConceptID {
		return false
	}
	if ePtr.owningConcept.isEquivalent(el.owningConcept) == false {
		return false
	}
	if ePtr.readOnly != el.readOnly {
		return false
	}
	if ePtr.Version.getVersion() != el.Version.getVersion() {
		return false
	}
	if ePtr.URI != el.URI {
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
	for _, child := range *ePtr.ownedConcepts.CopyMap() {
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
	buffer.WriteString(fmt.Sprintf("\"IsCore\":\"%t\"", ePtr.IsCore))
	return nil
}

func (ePtr *element) notifyListeners(underlyingNotification *ChangeNotification, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	indicatedConceptChanged := ePtr.uOfD.NewForwardingChangeNotification(ePtr, IndicatedConceptChanged, underlyingNotification)
	abstractionChanged := ePtr.uOfD.NewForwardingChangeNotification(ePtr, AbstractionChanged, underlyingNotification)
	for _, listener := range *ePtr.listeners.CopyMap() {
		switch listener.(type) {
		case *refinement:
			if listener.(*refinement).GetAbstractConcept(hl) == ePtr {
				refinedConcept := listener.(*refinement).GetRefinedConcept(hl)
				if refinedConcept != nil {
					ePtr.uOfD.queueFunctionExecutions(refinedConcept, abstractionChanged, hl)
				}
			} else {
				ePtr.uOfD.queueFunctionExecutions(listener, indicatedConceptChanged, hl)
			}
		default:
			ePtr.uOfD.queueFunctionExecutions(listener, indicatedConceptChanged, hl)
		}
	}
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
	ePtr.owningConcept.parentConceptID = recoveredConceptID
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
	ePtr.owningConcept.setIndicatedConceptID(recoveredOwningConceptID)
	owningConcept := ePtr.uOfD.GetElement(recoveredOwningConceptID)
	if owningConcept == nil {
		ePtr.uOfD.addUnresolvedPointer(ePtr.owningConcept)
	} else {
		ePtr.owningConcept.setIndicatedConcept(owningConcept)
		owningConcept.addOwnedConcept(ePtr.getConceptIDNoLock(), hl)
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
	return nil
}

// removeListener removes the indicated Element as a listening concept.
func (ePtr *element) removeListener(listeningConceptID string, hl *HeldLocks) {
	hl.ReadLockElement(ePtr)
	ePtr.uOfD.preChange(ePtr, hl)
	ePtr.listeners.DeleteEntry(listeningConceptID)
}

// removeOwnedConcept removes the indicated Element as a child (owned) concept.
func (ePtr *element) removeOwnedConcept(ownedConceptID string, hl *HeldLocks) error {
	hl.ReadLockElement(ePtr)
	ePtr.uOfD.preChange(ePtr, hl)
	if ePtr.IsReadOnly(hl) {
		return errors.New("Element.removedOwnedConcept called on read-only Element")
	}
	ePtr.ownedConcepts.DeleteEntry(ownedConceptID)
	return nil
}

// SetDefinition sets the definition of the Element
func (ePtr *element) SetDefinition(def string, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	editableError := ePtr.editableError()
	if editableError != nil {
		return editableError
	}
	if ePtr.Definition != def {
		ePtr.uOfD.preChange(ePtr, hl)
		notification := ePtr.uOfD.NewConceptChangeNotification(ePtr, hl)
		ePtr.uOfD.queueFunctionExecutions(ePtr, notification, hl)
		ePtr.Definition = def
	}
	return nil
}

func (ePtr *element) setIsCore(value bool, hl *HeldLocks) {
	hl.WriteLockElement(ePtr)
	ePtr.IsCore = value
}

// SetLabel sets the label of the Element
func (ePtr *element) SetLabel(label string, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	editableError := ePtr.editableError()
	if editableError != nil {
		return editableError
	}
	if ePtr.Label != label {
		ePtr.uOfD.preChange(ePtr, hl)
		notification := ePtr.uOfD.NewConceptChangeNotification(ePtr, hl)
		ePtr.uOfD.queueFunctionExecutions(ePtr, notification, hl)
		ePtr.Label = label
	}
	return nil
}

// SetOwningConceptID sets the ID of the owning concept for the element
// Design Note: the argument is the identifier rather than the Element to ensure
// the correct type of the owning concept is recorded. For example, if a method
// of *element calls this method with itself as the argument, the actual type
// recorded would be *element even if the actual caller is a *literal, *reference, or
// *refinement
func (ePtr *element) SetOwningConceptID(ocID string, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	editableError := ePtr.editableError()
	if editableError != nil {
		return editableError
	}
	// Do nothing if there is no change
	if ePtr.OwningConceptID != ocID {
		ePtr.uOfD.preChange(ePtr, hl)
		notification := ePtr.uOfD.NewConceptChangeNotification(ePtr, hl)
		ePtr.uOfD.queueFunctionExecutions(ePtr, notification, hl)
		oldOwner := ePtr.GetOwningConcept(hl)
		if oldOwner != nil {
			oldOwner.removeOwnedConcept(ePtr.ConceptID, hl)
		}
		var newOwner Element
		ePtr.OwningConceptID = ocID
		ePtr.owningConcept.setIndicatedConceptID(ocID)
		if ocID != "" {
			newOwner = ePtr.uOfD.GetElement(ocID)
			ePtr.owningConcept.setIndicatedConcept(newOwner)
			newOwner.addOwnedConcept(ePtr.ConceptID, hl)
		} else {
			ePtr.owningConcept.setIndicatedConcept(nil)
		}
	}
	return nil
}

// SetReadOnly provides a mechanism for preventing modifications to concepts. It will throw an error
// if the concept is one of the CRL core concepts, as these can never be made writable. It will also throw
// an error if there is an owner and it is read only
func (ePtr *element) SetReadOnly(value bool, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	editableError := ePtr.editableError()
	if editableError != nil {
		return editableError
	}
	if ePtr.GetOwningConcept(hl) != nil {
		ownerEditableError := ePtr.GetOwningConcept(hl).editableError()
		if ownerEditableError != nil {
			return ownerEditableError
		}
	}
	if ePtr.readOnly != value {
		ePtr.uOfD.preChange(ePtr, hl)
		notification := ePtr.uOfD.NewConceptChangeNotification(ePtr, hl)
		ePtr.uOfD.queueFunctionExecutions(ePtr, notification, hl)
		ePtr.readOnly = value
	}
	return nil
}

// setUniverseOfDiscourse is intended to be called only by the UniverseOfDiscourse
func (ePtr *element) setUniverseOfDiscourse(uOfD UniverseOfDiscourse, hl *HeldLocks) {
	hl.WriteLockElement(ePtr)
	ePtr.uOfD = uOfD
}

// SetURI sets the URI of the Element
func (ePtr *element) SetURI(uri string, hl *HeldLocks) error {
	hl.WriteLockElement(ePtr)
	editableError := ePtr.editableError()
	if editableError != nil {
		return editableError
	}
	if ePtr.URI != uri {
		ePtr.uOfD.preChange(ePtr, hl)
		ePtr.uOfD.changeURIForElement(ePtr, ePtr.URI, uri)
		notification := ePtr.uOfD.NewConceptChangeNotification(ePtr, hl)
		ePtr.uOfD.queueFunctionExecutions(ePtr, notification, hl)
		ePtr.URI = uri
	}
	return nil
}

func (ePtr *element) TraceableReadLock(hl *HeldLocks) {
	if TraceLocks {
		log.Printf("HL %p about to read lock Element %p\n", hl, ePtr)
	}
	ePtr.RLock()
}

func (ePtr *element) TraceableWriteLock(hl *HeldLocks) {
	if TraceLocks {
		log.Printf("HL %p about to write lock Element %p\n", hl, ePtr)
	}
	ePtr.Lock()
}

func (ePtr *element) TraceableReadUnlock(hl *HeldLocks) {
	if TraceLocks {
		log.Printf("HL %p about to read unlock Element %p\n", hl, ePtr)
	}
	ePtr.RUnlock()
}

func (ePtr *element) TraceableWriteUnlock(hl *HeldLocks) {
	if TraceLocks {
		log.Printf("HL %p about to write unlock Element %p\n", hl, ePtr)
	}
	ePtr.Unlock()
}

// Element is the representation of a concept
type Element interface {
	addListener(string, *HeldLocks)
	addOwnedConcept(string, *HeldLocks)
	addRecoveredOwnedConcept(string, *HeldLocks)
	editableError() error
	FindAbstractions(*map[string]Element, *HeldLocks)
	GetConceptID(*HeldLocks) string
	getConceptIDNoLock() string
	GetDefinition(*HeldLocks) string
	GetFirstChildWithAbstraction(Element, *HeldLocks) Element
	GetIsCore() bool
	GetLabel(*HeldLocks) string
	GetOwnedConcepts(*HeldLocks) *map[string]Element
	GetOwningConceptID(*HeldLocks) string
	GetOwningConcept(*HeldLocks) Element
	GetUniverseOfDiscourse(*HeldLocks) UniverseOfDiscourse
	getUniverseOfDiscourseNoLock() UniverseOfDiscourse
	GetURI(*HeldLocks) string
	GetVersion(*HeldLocks) int
	HasAbstraction(Element, *HeldLocks) bool
	incrementVersion(*HeldLocks)
	IsOwnedConcept(Element, *HeldLocks) bool
	IsReadOnly(*HeldLocks) bool
	MarshalJSON() ([]byte, error)
	notifyListeners(*ChangeNotification, *HeldLocks)
	removeListener(string, *HeldLocks)
	removeOwnedConcept(string, *HeldLocks) error
	SetDefinition(string, *HeldLocks) error
	setIsCore(bool, *HeldLocks)
	SetLabel(string, *HeldLocks) error
	SetOwningConceptID(string, *HeldLocks) error
	SetReadOnly(bool, *HeldLocks) error
	setUniverseOfDiscourse(UniverseOfDiscourse, *HeldLocks)
	SetURI(string, *HeldLocks) error
	TraceableReadLock(*HeldLocks)
	TraceableWriteLock(*HeldLocks)
	TraceableReadUnlock(*HeldLocks)
	TraceableWriteUnlock(*HeldLocks)
}
