package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"sort"
	"strconv"

	"github.com/pkg/errors"

	mapset "github.com/deckarep/golang-set"
	uuid "github.com/satori/go.uuid"
)

// UniverseOfDiscourse represents the scope of relevant concepts
type UniverseOfDiscourse struct {
	id                  string
	computeFunctions    functions
	executedCalls       chan *functionCallRecord
	undoManager         *undoManager
	uriUUIDMap          *StringStringMap
	uuidElementMap      *StringElementMap
	inProgressDeletions *StringElementMap
	ownedIDsMap         *OneToNStringMap
	listenersMap        *OneToNStringMap
	abstractionsMap     *OneToNStringMap
	observers           mapset.Set
}

// NewUniverseOfDiscourse creates and initializes a new UniverseOfDiscourse
func NewUniverseOfDiscourse() *UniverseOfDiscourse {
	var uOfD UniverseOfDiscourse
	newUUID, err := uuid.NewV4()
	if err != nil {
		log.Printf("failure in UniverseOfDiscourse.NewUniverseOfDiscourse %s", err)
	}
	uOfD.id = newUUID.String()
	uOfD.observers = mapset.NewSet()
	uOfD.computeFunctions = make(map[string][]crlExecutionFunction)
	uOfD.undoManager = newUndoManager(&uOfD)
	uOfD.uriUUIDMap = NewStringStringMap()
	uOfD.uuidElementMap = NewStringElementMap()
	uOfD.inProgressDeletions = NewStringElementMap()
	uOfD.ownedIDsMap = NewOneToNStringMap()
	uOfD.listenersMap = NewOneToNStringMap()
	uOfD.abstractionsMap = NewOneToNStringMap()
	trans := uOfD.NewTransaction()
	buildCoreDomain(&uOfD, trans)
	trans.ReleaseLocks()
	return &uOfD
}

// AddAbstractionToConcept creates a refinement with the concept being the refined concept. The new refinement is
// a child of the supplied concept.
func (uOfDPtr *UniverseOfDiscourse) AddAbstractionToConcept(con Concept, abstraction Concept, trans *Transaction) error {
	if con == nil {
		return errors.New("uOfD.AddAbstractionToConcept called with nil concept")
	}
	if abstraction == nil {
		return errors.New("uOfD.AddAbstractionToConcept called with nil abstraction")
	}
	refinement, err := uOfDPtr.NewRefinement(trans)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.AddAbstractionToConcept failed: ")
	}
	refinement.SetOwningConcept(con, trans)
	refinement.SetAbstractConcept(abstraction, trans)
	refinement.SetRefinedConcept(con, trans)
	refinement.SetLabel("Refines "+abstraction.GetLabel(trans), trans)
	return nil
}

// AddAbstractionURIToConcept creates a refinement with the concept being the refined concept. The new refinement is
// a child of the supplied concept.
func (uOfDPtr *UniverseOfDiscourse) AddAbstractionURIToConcept(con Concept, abstractionURI string, trans *Transaction) error {
	abstraction := uOfDPtr.GetElementWithURI(abstractionURI)
	if abstraction == nil {
		return errors.New("oOfD.AddAbstractionURIToConcept abstraction not found for abstractionURI: " + abstractionURI)
	}
	return uOfDPtr.AddAbstractionToConcept(con, abstraction, trans)
}

func (uOfDPtr *UniverseOfDiscourse) addElement(el Concept, inRecovery bool, trans *Transaction) error {
	if el == nil {
		return errors.New("UniverseOfDiscource addElement() failed because element was nil")
	}
	trans.WriteLockElement(el)
	uuid := el.GetConceptID(trans)
	if uuid == "" {
		return errors.New("UniverseOfDiscource addElement() failed because UUID was nil")
	}
	uOfDPtr.undoManager.markNewElement(el, trans)

	uOfDPtr.addElementForUndo(el, trans)

	uOfDPtr.postChange(el, trans)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) addElementForUndo(el Concept, trans *Transaction) error {
	if el == nil {
		return errors.New("UniverseOfDiscource addElementForUndo() failed because element was nil")
	}
	trans.WriteLockElement(el)
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("Adding element for undo, id: %s\n", el.GetConceptID(trans))
		Print(el, "Added Element: ", trans)
	}
	uOfDPtr.setUUIDElementMapEntry(el.GetConceptID(trans), el)
	uuid := el.GetConceptID(trans)
	uri := el.GetURI(trans)
	if uri != "" {
		uOfDPtr.uriUUIDMap.SetEntry(uri, el.GetConceptID(trans))
	}
	ownerID := el.GetOwningConceptID(trans)
	if ownerID != "" {
		uOfDPtr.addMappedValueToOwnedIDsMap(ownerID, uuid)
	}
	// Add element to all listener's lists
	switch el.GetConceptType() {
	case Reference:
		referencedConceptID := el.GetReferencedConceptID(trans)
		if referencedConceptID != "" {
			uOfDPtr.addMappedValueToListenersMap(referencedConceptID, uuid)
		}
	case Refinement:
		abstractConceptID := el.GetAbstractConceptID(trans)
		if abstractConceptID != "" {
			uOfDPtr.addMappedValueToListenersMap(abstractConceptID, uuid)
		}
		refinedConceptID := el.GetRefinedConceptID(trans)
		if refinedConceptID != "" {
			uOfDPtr.addMappedValueToListenersMap(refinedConceptID, uuid)
		}
	}

	return nil
}

// AddFunction registers a function with the indicated uri
func (uOfDPtr *UniverseOfDiscourse) AddFunction(uri string, function crlExecutionFunction) {
	uOfDPtr.computeFunctions[string(uri)] = append(uOfDPtr.computeFunctions[string(uri)], function)
}

func (uOfDPtr *UniverseOfDiscourse) changeURIForElement(el Concept, oldURI string, newURI string) error {
	if oldURI != "" && uOfDPtr.uriUUIDMap.GetEntry(oldURI) == el.getConceptIDNoLock() {
		uOfDPtr.uriUUIDMap.DeleteEntry(oldURI)
	}
	if newURI != "" {
		if uOfDPtr.uriUUIDMap.GetEntry(newURI) != "" {
			return errors.New("Attempted to assign a URI that is already in use")
		}
		uOfDPtr.uriUUIDMap.SetEntry(newURI, el.getConceptIDNoLock())
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) addMappedValueToListenersMap(listenerID string, value string) {
	uOfDPtr.listenersMap.addMappedValue(listenerID, value)
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("      addMappedValueToListenersMap listenerID: %s value: %s resultingOwnedIDs: %v", listenerID, value, uOfDPtr.listenersMap.GetMappedValues(listenerID))
	}
}

func (uOfDPtr *UniverseOfDiscourse) addMappedValueToOwnedIDsMap(ownerID string, ownedConceptID string) {
	uOfDPtr.ownedIDsMap.addMappedValue(ownerID, ownedConceptID)
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("      addMappedValueToOwnedIDsMap ownerID: %s newOwnedConceptID: %s resultingOwnedIDs: %v", ownerID, ownedConceptID, uOfDPtr.ownedIDsMap.GetMappedValues(ownerID))
	}
}

// Clone makes an exact copy of the UniverseOfDiscourse and all its contents except for the undo/redo stack. All Elements are new objects,
// but all the identifiers are retained from the original uOfD.
func (uOfDPtr *UniverseOfDiscourse) Clone(trans *Transaction) *UniverseOfDiscourse {
	newUofD := NewUniverseOfDiscourse()

	// uOfD.computeFunctions = make(map[string][]crlExecutionFunction)
	var uri string
	var functionArray []crlExecutionFunction
	for uri, functionArray = range uOfDPtr.computeFunctions {
		// Housekeeping functions are already present in a new uOfD
		if uri != "http://activeCrl.com/core/coreHousekeeping" {
			newUofD.computeFunctions[uri] = append(newUofD.computeFunctions[uri], functionArray...)
		}
	}

	for uri, uuid := range uOfDPtr.uriUUIDMap.CopyMap() {
		newUofD.uriUUIDMap.SetEntry(uri, uuid)
	}

	for id, el := range uOfDPtr.uuidElementMap.CopyMap() {
		newElement := clone(el, trans)
		newUofD.setUUIDElementMapEntry(id, newElement)
	}
	for key, strings := range uOfDPtr.ownedIDsMap.CopyMap() {
		newUofD.ownedIDsMap.SetMappedValues(key, strings)
	}

	// uOfD.listenersMap = NewOneToNStringMap()
	for key, strings := range uOfDPtr.listenersMap.CopyMap() {
		newUofD.listenersMap.SetMappedValues(key, strings)
	}

	// uOfD.abstractionsMap = NewOneToNStringMap()
	for key, strings := range uOfDPtr.abstractionsMap.CopyMap() {
		newUofD.abstractionsMap.SetMappedValues(key, strings)
	}

	return newUofD
}

// CreateOwnedRefinementOfConcept creates a refinement of the concept identified by the URI, setting its owner and label
func (uOfDPtr *UniverseOfDiscourse) CreateOwnedRefinementOfConcept(original Concept, owner Concept, label string, trans *Transaction, newURI ...string) (Concept, error) {
	newConcept, err := uOfDPtr.CreateRefinementOfConcept(original, label, trans, newURI...)
	if err != nil {
		return nil, errors.Wrap(err, "CreateOwnedRefinementOfConceptURI failed")
	}
	newConcept.SetOwningConcept(owner, trans)
	return newConcept, nil
}

// CreateOwnedRefinementOfConceptURI creates a refinement of the concept identified by the URI, setting its owner and label
func (uOfDPtr *UniverseOfDiscourse) CreateOwnedRefinementOfConceptURI(originalURI string, owner Concept, label string, trans *Transaction, newURI ...string) (Concept, error) {
	original := uOfDPtr.GetElementWithURI(originalURI)
	if original == nil {
		return nil, errors.New("In CreateOwnedRefinementOfConceptURI, no concept found with uri: " + originalURI)
	}
	return uOfDPtr.CreateOwnedRefinementOfConcept(original, owner, label, trans, newURI...)
}

// CreateRefinementOfConcept creates a new concept of the same type as the original and makes the new concept a refinement of
// the original
func (uOfDPtr *UniverseOfDiscourse) CreateRefinementOfConcept(original Concept, label string, trans *Transaction, newURI ...string) (Concept, error) {
	refinedConcept, err := uOfDPtr.NewConcept(original.GetConceptType(), trans, newURI...)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.CreateRefinementOfConcept uOfD.NewConcept failed")
	}
	switch original.GetConceptType() {
	case Reference:
		err = refinedConcept.SetReferencedConceptID("", original.GetReferencedAttributeName(trans), trans)
		if err != nil {
			return nil, errors.Wrap(err, "UniverseOfDiscourse.CreateRefinementOfConcept refinedConcept.SetReferencedConceptID failed")
		}
	}
	err = refinedConcept.SetLabel(label, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.CreateRefinementOfConcept refinedConcept.SetLabel failed")
	}
	refinementURI := ""
	if len(newURI) == 1 && newURI[0] != "" {
		refinementURI = newURI[0] + original.GetConceptID(trans) + "/Refinement"
	}
	refinement, err := uOfDPtr.NewRefinement(trans, refinementURI)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.CreateRefinementOfConcept failed: ")
	}
	refinement.SetOwningConcept(refinedConcept, trans)
	refinement.SetAbstractConcept(original, trans)
	refinement.SetRefinedConcept(refinedConcept, trans)
	refinement.SetLabel("Refines "+original.GetLabel(trans), trans)
	return refinedConcept, nil
}

// CreateRefinementOfConceptURI creates a new concept of the same type as the indicated URI and makes the new concept a refinement of it
func (uOfDPtr *UniverseOfDiscourse) CreateRefinementOfConceptURI(uri string, label string, trans *Transaction, newURI ...string) (Concept, error) {
	abstractConcept := uOfDPtr.GetElementWithURI(uri)
	if abstractConcept == nil {
		return nil, errors.New("In CreateRefinementOfConceptURI, no concept found with URI: " + uri)
	}
	return uOfDPtr.CreateRefinementOfConcept(abstractConcept, label, trans, newURI...)
}

// CreateReplicateAsRefinement replicates the indicated Element and all of its descendent Elements
// except that descendant Refinements are not replicated.
// For each replicated Element, a Refinement is created with the abstractElement being the original and the refinedElement
// being the replica. The root replicated element is returned.
//
// Deprecated: Create specific New<type>() functions instead that construct the desired concept with expected children
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateAsRefinement(original Concept, trans *Transaction, newURI ...string) (Concept, error) {
	replicate, err := uOfDPtr.NewConcept(original.GetConceptType(), trans, newURI...)
	if err != nil {
		return nil, err
	}
	err = uOfDPtr.replicateAsRefinement(original, replicate, trans, newURI...)
	if err != nil {
		return nil, err
	}
	return replicate, nil
}

// CreateReplicateAsRefinementFromURI replicates the Element indicated by the URI
//
// Deprecated: Create specific New<type>() functions instead that construct the desired concept with expected children
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateAsRefinementFromURI(originalURI string, trans *Transaction, newURI ...string) (Concept, error) {
	original := uOfDPtr.GetElementWithURI(originalURI)
	if original == nil {
		return nil, fmt.Errorf("in CreateReplicateAsRefinementFromURI Element with uri %s not found", originalURI)
	}
	return uOfDPtr.CreateReplicateAsRefinement(original, trans, newURI...)
}

// CreateReplicateLiteralAsRefinement replicates the supplied Literal and makes all elements of the replicate
// refinements of the original elements
//
// Deprecated: Create specific New<type>() functions instead that construct the desired concept with expected children
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateLiteralAsRefinement(original Concept, trans *Transaction, newURI ...string) (Concept, error) {
	replicate, err := uOfDPtr.NewLiteral(trans, newURI...)
	if err != nil {
		return nil, err
	}
	err = uOfDPtr.replicateAsRefinement(original, replicate, trans, newURI...)
	if err != nil {
		return nil, err
	}
	return replicate, nil
}

// CreateReplicateLiteralAsRefinementFromURI replicates the Literal indicated by the URI
//
// Deprecated: Create specific New<type>() functions instead that construct the desired concept with expected children
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateLiteralAsRefinementFromURI(originalURI string, trans *Transaction, newURI ...string) (Concept, error) {
	original := uOfDPtr.GetLiteralWithURI(originalURI)
	if original == nil {
		return nil, fmt.Errorf("in CreateReplicateLiteralAsRefinementFromURI Element with uri %s not found", originalURI)
	}
	return uOfDPtr.CreateReplicateLiteralAsRefinement(original, trans, newURI...)
}

// CreateReplicateReferenceAsRefinement replicates the supplied reference and makes all elements of the replicate
// refinements of the original elements
//
// Deprecated: Create specific New<type>() functions instead that construct the desired concept with expected children
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateReferenceAsRefinement(original Concept, trans *Transaction, newURI ...string) (Concept, error) {
	if original.GetConceptType() != Reference {
		return nil, errors.New("CreateReplicateReferenceAsRefinement called with non-Reference argument")
	}
	replicate, err := uOfDPtr.NewReference(trans, newURI...)
	if err != nil {
		return nil, err
	}
	err = uOfDPtr.replicateAsRefinement(original, replicate, trans, newURI...)
	if err != nil {
		return nil, err
	}
	return replicate, nil
}

// CreateReplicateReferenceAsRefinementFromURI replicates the Reference indicated by the URI
//
// Deprecated: Create specific New<type>() functions instead that construct the desired concept with expected children
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateReferenceAsRefinementFromURI(originalURI string, trans *Transaction, newURI ...string) (Concept, error) {
	original := uOfDPtr.GetReferenceWithURI(originalURI)
	if original == nil {
		return nil, fmt.Errorf("in CreateReplicateAsRefinementFromURI Element with uri %s not found", originalURI)
	}
	return uOfDPtr.CreateReplicateReferenceAsRefinement(original, trans, newURI...)
}

// CreateReplicateRefinementAsRefinement replicates the supplied refinement and makes all elements of the replicate
// refinements of the original elements
//
// Deprecated: Create specific New<type>() functions instead that construct the desired concept with expected children
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateRefinementAsRefinement(original Concept, trans *Transaction, newURI ...string) (Concept, error) {
	replicate, err := uOfDPtr.NewRefinement(trans, newURI...)
	if err != nil {
		return nil, err
	}
	err = uOfDPtr.replicateAsRefinement(original, replicate, trans, newURI...)
	if err != nil {
		return nil, err
	}
	return replicate, nil
}

// CreateReplicateRefinementAsRefinementFromURI replicates the Refinement indicated by the URI
//
// Deprecated: Create specific New<type>() functions instead that construct the desired concept with expected children
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateRefinementAsRefinementFromURI(originalURI string, trans *Transaction, newURI ...string) (Concept, error) {
	original := uOfDPtr.GetRefinementWithURI(originalURI)
	if original == nil {
		return nil, fmt.Errorf("in CreateReplicateRefinementAsRefinementFromURI Element with uri %s not found", originalURI)
	}
	return uOfDPtr.CreateReplicateRefinementAsRefinement(original, trans, newURI...)
}

func (uOfDPtr *UniverseOfDiscourse) findFunctions(element Concept, notification *ChangeNotification, trans *Transaction) []string {
	var functionIdentifiers []string
	if element == nil {
		return functionIdentifiers
	}
	// Now find functions associated with self and abstractions
	selfAndAbstractions := make(map[string]Concept)
	selfAndAbstractions[element.getConceptIDNoLock()] = element
	element.FindAbstractions(selfAndAbstractions, trans)
	for _, candidate := range selfAndAbstractions {
		uri := candidate.GetURI(trans)
		if uri != "" {
			functions := uOfDPtr.computeFunctions[uri]
			if functions != nil {
				functionIdentifiers = append(functionIdentifiers, uri)
			}
		}
	}
	return functionIdentifiers
}

func (uOfDPtr *UniverseOfDiscourse) deleteElement(el Concept, trans *Transaction) error {
	if el == nil {
		return errors.New("UniverseOfDiscource removeElement failed elcause Element was nil")
	}
	trans.WriteLockElement(el)
	beforeState, err := NewConceptState(el)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.deleteElement failed")
	}
	uOfDPtr.undoManager.markRemovedElement(el, trans)

	uuid := el.GetConceptID(trans)
	uri := el.GetURI(trans)
	if uri != "" {
		uOfDPtr.uriUUIDMap.DeleteEntry(uri)
	}
	// Remove element from owner's child list
	ownerID := el.GetOwningConceptID(trans)
	if ownerID != "" {
		el.SetOwningConceptID("", trans)
	}
	it := uOfDPtr.listenersMap.GetMappedValues(uuid).Iterator()
	for id := range it.C {
		listener := uOfDPtr.GetElement(id.(string))
		if listener != nil {
			switch listener.GetConceptType() {
			case Reference:
				listener.SetReferencedConcept(nil, NoAttribute, trans)
			case Refinement:
				if listener.GetAbstractConcept(trans) == el {
					listener.SetAbstractConcept(nil, trans)
				} else if listener.GetRefinedConcept(trans) == el {
					listener.SetRefinedConcept(nil, trans)
				}
			}
		}
	}
	// Spread the news
	conceptRemovedNotification := uOfDPtr.newUofDConceptRemovedNotification(beforeState, trans)
	err = uOfDPtr.NotifyUofDObservers(conceptRemovedNotification, trans)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.deleteElement failed")
	}
	err = el.notifyObservers(conceptRemovedNotification, trans)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.deleteElement failed")
	}
	// Remove element from all listener's lists
	switch el.GetConceptType() {
	case Reference:
		referencedConceptID := el.GetReferencedConceptID(trans)
		if referencedConceptID != "" {
			uOfDPtr.removeMappedValueFromListenersMap(referencedConceptID, uuid)
		}
	case Refinement:
		abstractConceptID := el.GetAbstractConceptID(trans)
		if abstractConceptID != "" {
			uOfDPtr.removeMappedValueFromListenersMap(abstractConceptID, uuid)
		}
		refinedConceptID := el.GetRefinedConceptID(trans)
		if refinedConceptID != "" {
			uOfDPtr.removeMappedValueFromListenersMap(refinedConceptID, uuid)
		}
	}
	uOfDPtr.deleteKeyFromListenersMap(uuid)
	uOfDPtr.abstractionsMap.DeleteKey(uuid)
	uOfDPtr.deleteKeyFromOwnedIDsMap(uuid)
	uOfDPtr.deleteUUIDElementMapEntry(uuid)
	// Finally, remove from the universe of discourse
	el.setUniverseOfDiscourse(nil, trans)
	return nil
}

// DeleteElement removes a single element and its descentants from the uOfD. Pointers to the elements from other elements are set to nil.
func (uOfDPtr *UniverseOfDiscourse) DeleteElement(element Concept, trans *Transaction) error {
	id := element.GetConceptID(trans)
	elements := mapset.NewSet(id)
	return uOfDPtr.DeleteElements(elements, trans)
}

// DeleteElements removes the elements whose IDs are in the set from the uOfD. Pointers to the elements from elements not being deleted are set to nil.
func (uOfDPtr *UniverseOfDiscourse) DeleteElements(elements mapset.Set, trans *Transaction) error {
	it := elements.Iterator()
	for id := range it.C {
		el := uOfDPtr.GetElement(id.(string))
		if el == nil {
			break
		}
		if el.GetIsCore(trans) {
			it.Stop()
			return errors.New("UniverseOfDiscourse.DeleteElements called on a CRL core concept")
		}
		if el.GetUniverseOfDiscourse(trans) != uOfDPtr {
			it.Stop()
			return errors.New("UniverseOfDiscourse.DeleteElements called on an Element in a different UofD")
		}
		if el.IsReadOnly(trans) {
			it.Stop()
			return errors.New("UniverseOfDiscourse.DeleteElements called on read-only Element")
		}
		uOfDPtr.inProgressDeletions.SetEntry(id.(string), el)
		descendants := mapset.NewSet()
		uOfDPtr.GetConceptsOwnedConceptIDsRecursively(id.(string), descendants, trans)
		it3 := descendants.Iterator()
		for descendantID := range it3.C {
			descendant := uOfDPtr.GetElement(descendantID.(string))
			if descendant != nil {
				uOfDPtr.inProgressDeletions.SetEntry(descendantID.(string), descendant)
			}
		}
	}
	for _, el := range uOfDPtr.inProgressDeletions.CopyMap() {
		if el != nil {
			trans.WriteLockElement(el)
			uOfDPtr.deleteElement(el, trans)
		}
	}
	uOfDPtr.inProgressDeletions.Clear()
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) deleteUUIDElementMapEntry(conceptID string) {
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("About to delete UUIDElementMapEntry ID: %s", conceptID)
	}
	uOfDPtr.uuidElementMap.DeleteEntry(conceptID)

}

func (uOfDPtr *UniverseOfDiscourse) deleteKeyFromListenersMap(key string) {
	uOfDPtr.listenersMap.DeleteKey(key)
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("      deleteKeyFromListenersMap key: %s", key)
	}
}

func (uOfDPtr *UniverseOfDiscourse) deleteKeyFromOwnedIDsMap(key string) {
	uOfDPtr.ownedIDsMap.DeleteKey(key)
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("      deleteKeyFromOwnedIDsMap key: %s", key)
	}
}

// Deregister removes the registration of an Observer
func (uOfDPtr *UniverseOfDiscourse) Deregister(observer Observer) error {
	uOfDPtr.observers.Remove(observer)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) generateConceptID(uri ...string) (string, error) {
	var conceptID string
	if len(uri) == 0 || (len(uri) == 1 && uri[0] == "") {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return "", errors.Wrap(err, "failure in UniverseOfDiscourse.generateConceptID")
		}
		conceptID = newUUID.String()
	} else {
		if len(uri) == 1 {
			_, err := url.ParseRequestURI(uri[0])
			if err != nil {
				return "", errors.New("Invalid URI provided for initializing Element")
			}
			conceptID = uuid.NewV5(uuid.NamespaceURL, uri[0]).String()
		} else {
			return "", errors.New("Invalid URI provided for initializing Element")
		}
	}
	return conceptID, nil
}

// getComputeFunctions returns a pointer to the compute functions. It is intended for the exclusive use of the
// FunctionCallManager
// func (uOfDPtr *UniverseOfDiscourse) getComputeFunctions() *functions {
// 	return &uOfDPtr.computeFunctions
// }

// GetElement returns the Element with the conceptID
func (uOfDPtr *UniverseOfDiscourse) GetElement(conceptID string) Concept {
	return uOfDPtr.uuidElementMap.GetEntry(conceptID)
}

// GetElementLabel returns the label of the Element with the conceptID
func (uOfDPtr *UniverseOfDiscourse) GetElementLabel(conceptID string) string {
	element := uOfDPtr.uuidElementMap.GetEntry(conceptID)
	if element != nil {
		return element.getLabelNoLock()
	}
	return "<deleted>"
}

// GetElements returns the Elements in the uOfD mapped by their ConceptIDs
func (uOfDPtr *UniverseOfDiscourse) GetElements() map[string]Concept {
	return uOfDPtr.uuidElementMap.CopyMap()
}

// GetElementWithURI returns the Element with the given URI
func (uOfDPtr *UniverseOfDiscourse) GetElementWithURI(uri string) Concept {
	return uOfDPtr.GetElement(uOfDPtr.uriUUIDMap.GetEntry(uri))
}

func (uOfDPtr *UniverseOfDiscourse) getExecutedCalls() chan *functionCallRecord {
	return uOfDPtr.executedCalls
}

// getFunctions returns the array of functions associatee with the given URI
func (uOfDPtr *UniverseOfDiscourse) getFunctions(uri string) []crlExecutionFunction {
	return uOfDPtr.computeFunctions[string(uri)]
}

// GetIDForURI returns a V5 UUID derived from the given URI. If the given URI
// is not valid it returns the empty string.
func (uOfDPtr *UniverseOfDiscourse) GetIDForURI(uri string) string {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return ""
	}
	return uuid.NewV5(uuid.NamespaceURL, uri).String()
}

// GetListenerIDs returns the set of listener IDs for the indicated ID
func (uOfDPtr *UniverseOfDiscourse) GetListenerIDs(id string) mapset.Set {
	return uOfDPtr.listenersMap.GetMappedValues(id)
}

// GetLiteral returns the literal with the indicated ID (if found)
func (uOfDPtr *UniverseOfDiscourse) GetLiteral(conceptID string) Concept {
	el := uOfDPtr.GetElement(conceptID)
	if el == nil {
		return nil
	}
	switch el.GetConceptType() {
	case Literal:
		return el
	}
	return nil
}

// GetLiteralWithURI returns the literal with the indicated URI (if found)
func (uOfDPtr *UniverseOfDiscourse) GetLiteralWithURI(uri string) Concept {
	el := uOfDPtr.GetElementWithURI(uri)
	if el == nil {
		return nil
	}
	switch el.GetConceptType() {
	case Literal:
		return el
	}
	return nil
}

// GetConceptsOwnedConceptIDs returns the set of owned concepts for the indicated ID
func (uOfDPtr *UniverseOfDiscourse) GetConceptsOwnedConceptIDs(id string) mapset.Set {
	return uOfDPtr.ownedIDsMap.GetMappedValues(id)
}

// GetConceptsOwnedConceptIDsRecursively returns the IDs of owned concepts
func (uOfDPtr *UniverseOfDiscourse) GetConceptsOwnedConceptIDsRecursively(rootID string, descendants mapset.Set, trans *Transaction) {
	it := uOfDPtr.ownedIDsMap.GetMappedValues(rootID).Iterator()
	for id := range it.C {
		descendants.Add(id.(string))
		uOfDPtr.GetConceptsOwnedConceptIDsRecursively(id.(string), descendants, trans)
	}
}

// GetReference returns the reference with the indicated ID (if found)
func (uOfDPtr *UniverseOfDiscourse) GetReference(conceptID string) Concept {
	el := uOfDPtr.GetElement(conceptID)
	if el == nil {
		return nil
	}
	switch el.GetConceptType() {
	case Reference:
		return el
	}
	return nil
}

// GetReferenceWithURI returns the reference with the indicated URI (if found)
func (uOfDPtr *UniverseOfDiscourse) GetReferenceWithURI(uri string) Concept {
	el := uOfDPtr.GetElementWithURI(uri)
	if el == nil {
		return nil
	}
	switch el.GetConceptType() {
	case Reference:
		return el
	}
	return nil
}

// GetRefinement returns the refinement with the indicated ID (if found)
func (uOfDPtr *UniverseOfDiscourse) GetRefinement(conceptID string) Concept {
	el := uOfDPtr.GetElement(conceptID)
	if el == nil {
		return nil
	}
	switch el.GetConceptType() {
	case Refinement:
		return el
	}
	return nil
}

// GetRefinementWithURI returns the refinement with the indicated URI (if found)
func (uOfDPtr *UniverseOfDiscourse) GetRefinementWithURI(uri string) Concept {
	el := uOfDPtr.GetElementWithURI(uri)
	if el == nil {
		return nil
	}
	switch el.GetConceptType() {
	case Refinement:
		return el
	}
	return nil
}

// GetRootElements returns all elements that do not have owners
func (uOfDPtr *UniverseOfDiscourse) GetRootElements(trans *Transaction) map[string]Concept {
	allElements := uOfDPtr.GetElements()
	rootElements := make(map[string]Concept)
	for id, el := range allElements {
		if el.GetOwningConceptID(trans) == "" && el.GetURI(trans) != TransientURI {
			rootElements[id] = el
		}
	}
	return rootElements
}

// GetRootElementIDs returns all elements that do not have owners
func (uOfDPtr *UniverseOfDiscourse) GetRootElementIDs() []string {
	allElements := uOfDPtr.GetElements()
	var ids []string
	for id, el := range allElements {
		if el.getOwningConceptIDNoLock() == "" && el.getURINoLock() != TransientURI {
			ids = append(ids, id)
		}
	}
	return ids
}

func (uOfDPtr *UniverseOfDiscourse) getURIUUIDMap() *StringStringMap {
	return uOfDPtr.uriUUIDMap
}

// IsEquivalent returns true if all of the root elements in the uOfD are recursively equivalent
func (uOfDPtr *UniverseOfDiscourse) IsEquivalent(hl1 *Transaction, uOfD2 *UniverseOfDiscourse, hl2 *Transaction, printExceptions ...bool) bool {
	var printEquivalenceExceptions bool
	if len(printExceptions) > 0 {
		printEquivalenceExceptions = printExceptions[0]
	}
	// Functions
	// uOfD.computeFunctions = make(map[string][]crlExecutionFunction)
	if len(uOfDPtr.computeFunctions) != len(uOfD2.computeFunctions) {
		if printEquivalenceExceptions {
			log.Printf("Length of compute functions map not equivalent")
		}
		return false
	}
	var uri string
	var functionArray []crlExecutionFunction
	for uri, functionArray = range uOfDPtr.computeFunctions {
		if len(functionArray) != len(uOfD2.computeFunctions[uri]) {
			if printEquivalenceExceptions {
				log.Printf("Length of compute functions array not equivalent for uri: %s", uri)
			}
			return false
		}
		var crlFunction crlExecutionFunction
		var crlFunction2 crlExecutionFunction
		var i int
		for i, crlFunction = range functionArray {
			crlFunction2 = uOfD2.computeFunctions[uri][i]
			if reflect.ValueOf(crlFunction).Pointer() != reflect.ValueOf(crlFunction2).Pointer() {
				if printEquivalenceExceptions {
					log.Printf("The %dth compute function is not equivalent for uri: %s", i, uri)
				}
				return false
			}
		}
	}

	// uOfD.uriUUIDMap = NewStringStringMap()
	if !uOfDPtr.uriUUIDMap.IsEquivalent(uOfD2.uriUUIDMap, printEquivalenceExceptions) {
		if printEquivalenceExceptions {
			log.Printf("uriUUDIMap not equivalent")
		}
		return false
	}

	// uOfD.uuidElementMap = NewStringElementMap()
	if !uOfDPtr.uuidElementMap.IsEquivalent(uOfD2.uuidElementMap) {
		if printEquivalenceExceptions {
			log.Printf("uriUUDIMap keys not equivalent")
		}
		return false
	}

	// uOfD.ownedIDsMap
	if !uOfDPtr.ownedIDsMap.IsEquivalent(uOfD2.ownedIDsMap) {
		if printEquivalenceExceptions {
			log.Printf("ownedIDsMap not equivalent")
		}
		return false
	}

	// uOfD.listenersMap
	if !uOfDPtr.listenersMap.IsEquivalent(uOfD2.listenersMap) {
		if printEquivalenceExceptions {
			log.Printf("listenersMap not equivalent")
			map1Keys := uOfDPtr.listenersMap.GetSortedKeys()
			map2Keys := uOfD2.listenersMap.GetSortedKeys()
			for _, map1Key := range map1Keys {
				map1Listeners := uOfDPtr.listenersMap.GetMappedValues(map1Key)
				map2Listeners := uOfD2.listenersMap.GetMappedValues(map1Key)
				map1ListenerKeys := []string{}
				map1ListenerIterator := map1Listeners.Iterator()
				for listenerKey := range map1ListenerIterator.C {
					map1ListenerKeys = append(map1ListenerKeys, listenerKey.(string))
				}
				sort.Slice(map1ListenerKeys, func(i, j int) bool {
					return map1ListenerKeys[i] < map1ListenerKeys[j]
				})
				map2ListenerKeys := []string{}
				map2ListenerIterator := map2Listeners.Iterator()
				for listenerKey := range map2ListenerIterator.C {
					map2ListenerKeys = append(map2ListenerKeys, listenerKey.(string))
				}
				sort.Slice(map2ListenerKeys, func(i, j int) bool {
					return map2ListenerKeys[i] < map2ListenerKeys[j]
				})
				if map1Listeners.Cardinality() != map2Listeners.Cardinality() || !reflect.DeepEqual(map1ListenerKeys, map2ListenerKeys) {
					log.Printf("  Map1 %d keys; Map 2 %d keys", len(map1Keys), len(map2Keys))
					log.Printf("    Map1 Key: %s has %d listeners and Map2 has %d listeners", map1Key, map1Listeners.Cardinality(), map2Listeners.Cardinality())
					log.Printf("      Map1 Key %s listeners: %v", map1Key, map1ListenerKeys)
					log.Printf("      Map2 Key %s listeners: %v", map1Key, map2ListenerKeys)
				}
			}
			for _, map2Key := range map2Keys {
				if uOfDPtr.listenersMap.GetMappedValues(map2Key) == nil {
					log.Printf("    Map2 Key: %s not found in Map1", map2Key)
					map2Listeners := uOfD2.listenersMap.GetMappedValues(map2Key)
					map2ListenerKeys := []string{}
					map2ListenerIterator := map2Listeners.Iterator()
					for listenerKey := range map2ListenerIterator.C {
						map2ListenerKeys = append(map2ListenerKeys, listenerKey.(string))
					}
					log.Printf("      Map2 Key %s listeners: %v", map2Key, map2ListenerKeys)
				}
			}
		}
		return false
	}

	rootElements1 := uOfDPtr.GetRootElements(hl1)
	rootElements2 := uOfD2.GetRootElements(hl2)
	for id1, el1 := range rootElements1 {
		el2 := rootElements2[id1]
		if el2 == nil || !RecursivelyEquivalent(el1, hl1, el2, hl2, printEquivalenceExceptions) {
			return false
		}
	}

	return true
}

// IsRecordingUndo reveals whether undo recording is on
func (uOfDPtr *UniverseOfDiscourse) IsRecordingUndo() bool {
	// TODO Remove this debugging code
	if uOfDPtr == nil {
		log.Fatal("In UniverseOfDiscourse.IsRecordingUndo() will nil uOfDPtr")
	}
	return uOfDPtr.undoManager.recordingUndo
}

// MarkUndoPoint marks a point on the undo stack. The next undo operation will undo everything back to this point.
func (uOfDPtr *UniverseOfDiscourse) MarkUndoPoint() {
	uOfDPtr.undoManager.MarkUndoPoint()
}

// MarshalDomain creates a JSON representation of an element and all of its descendants
func (uOfDPtr *UniverseOfDiscourse) MarshalDomain(el Concept, trans *Transaction) ([]byte, error) {
	var result []byte
	result = append(result, []byte("[")...)
	marshaledConcept, err := uOfDPtr.marshalConceptRecursively(el, trans)
	if err != nil {
		return result, err
	}
	// The last byte of marshaledConcept is going to be a comma we don't want
	result = append(result, marshaledConcept[0:len(marshaledConcept)-1]...)
	result = append(result, []byte("]")...)
	return result, nil
}

func (uOfDPtr *UniverseOfDiscourse) marshalConceptRecursively(el Concept, trans *Transaction) ([]byte, error) {
	var result []byte
	if el == nil {
		return result, errors.New("UniverseOfDiscourse.marshalConceptRecursively called with nil concept")
	}
	marshaledElement, err := el.MarshalJSON()
	if err != nil {
		return result, err
	}
	result = append(result, marshaledElement...)
	result = append(result, []byte(",")...)
	elID := el.GetConceptID(trans)
	it := uOfDPtr.GetConceptsOwnedConceptIDs(elID).Iterator()
	for id := range it.C {
		child := uOfDPtr.GetElement(id.(string))
		marshaledChild, err := uOfDPtr.marshalConceptRecursively(child, trans)
		if err != nil {
			it.Stop()
			return result, err
		}
		result = append(result, marshaledChild...)
	}
	return result, nil
}

// newUofDConceptAddedNotification creates a UofDConceptAdded notification
func (uOfDPtr *UniverseOfDiscourse) newUofDConceptAddedNotification(afterState *ConceptState, trans *Transaction) *ChangeNotification {
	var notification ChangeNotification
	notification.afterConceptState = afterState
	notification.natureOfChange = ConceptAdded
	notification.uOfD = uOfDPtr
	return &notification
}

// SendConceptChangeNotification creates a ConceptChangeNotification
func (uOfDPtr *UniverseOfDiscourse) SendConceptChangeNotification(reportingElement Concept, beforeState *ConceptState, afterState *ConceptState, trans *Transaction) error {
	notification := &ChangeNotification{}
	reportingConceptState, err := NewConceptState(reportingElement)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.NewConceptChangeNotification failed")
	}
	notification.reportingElementState = reportingConceptState
	notification.beforeConceptState = beforeState
	notification.afterConceptState = afterState
	notification.natureOfChange = ConceptChanged
	notification.uOfD = uOfDPtr
	err = reportingElement.propagateChange(notification, trans)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.SendConceptChangeNotification failed")
	}
	return nil
}

// SendPointerChangeNotification creates a PointerChangeNotification and sends it to the relevant parties
func (uOfDPtr *UniverseOfDiscourse) SendPointerChangeNotification(reportingElement Concept, natureOfChange NatureOfChange, beforeConceptState *ConceptState, afterConceptState *ConceptState, trans *Transaction) error {
	notification := &ChangeNotification{}
	reportingConceptState, err := NewConceptState(reportingElement)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.NewPointerChangeNotification failed")
	}
	notification.reportingElementState = reportingConceptState
	notification.beforeConceptState = beforeConceptState
	notification.afterConceptState = afterConceptState
	notification.natureOfChange = natureOfChange
	notification.uOfD = uOfDPtr
	reportingElement.propagateChange(notification, trans)
	return nil
}

// SendTickleNotification creates a Tickle notification and sends it to the indicated target element. Its purpose is to
// trigger the execution of any functions associated with the target element
func (uOfDPtr *UniverseOfDiscourse) SendTickleNotification(reportingElement Concept, targetElement Concept, trans *Transaction) error {
	notification := &ChangeNotification{}
	reportingConceptState, err := NewConceptState(reportingElement)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.SendTickleNotification failed")
	}
	notification.reportingElementState = reportingConceptState
	notification.beforeConceptState = nil
	notification.afterConceptState = nil
	notification.natureOfChange = Tickle
	notification.uOfD = uOfDPtr
	reportingElement.tickle(targetElement, notification, trans)
	return nil
}

// newUofDConceptRemovedNotification creates a UniverseOfDiscourseRemoved notification
func (uOfDPtr *UniverseOfDiscourse) newUofDConceptRemovedNotification(beforeState *ConceptState, trans *Transaction) *ChangeNotification {
	var notification ChangeNotification
	notification.natureOfChange = ConceptRemoved
	notification.beforeConceptState = beforeState
	notification.uOfD = uOfDPtr
	return &notification
}

// NewForwardingChangeNotification creates a ChangeNotification that records the reason for the change to the element,
// including the nature of the change, an indication of which component originated the change, and whether there
// was a preceeding notification that triggered this change.
func (uOfDPtr *UniverseOfDiscourse) NewForwardingChangeNotification(reportingElement Concept, natureOfChange NatureOfChange, underlyingChange *ChangeNotification, trans *Transaction) (*ChangeNotification, error) {
	notification := &ChangeNotification{}
	reportingElementState, err := NewConceptState(reportingElement)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewForwardingChangeNotification failed")
	}
	notification.reportingElementState = reportingElementState
	notification.natureOfChange = natureOfChange
	notification.underlyingChange = underlyingChange
	notification.uOfD = uOfDPtr
	return notification, nil
}

// NewConcept creates and initializes a new Concept
func (uOfDPtr *UniverseOfDiscourse) NewConcept(conceptType ConceptType, trans *Transaction, uri ...string) (Concept, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var el concept
	el.initializeConcept(conceptType, conceptID, actualURI)
	trans.WriteLockElement(&el)
	uOfDPtr.SetUniverseOfDiscourse(&el, trans)
	if actualURI != "" {
		el.SetURI(actualURI, trans)
	}
	return &el, nil
}

// NewElement creates and initializes a new Element
func (uOfDPtr *UniverseOfDiscourse) NewElement(trans *Transaction, uri ...string) (Concept, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var el concept
	el.initializeConcept(Element, conceptID, actualURI)
	trans.WriteLockElement(&el)
	uOfDPtr.SetUniverseOfDiscourse(&el, trans)
	if actualURI != "" {
		el.SetURI(actualURI, trans)
	}
	return &el, nil
}

// NewTransaction creates and initializes a HeldLocks structure utilizing the supplied WaitGroup
func (uOfDPtr *UniverseOfDiscourse) NewTransaction() *Transaction {
	var trans Transaction
	trans.readLocks = make(map[string]Concept)
	trans.writeLocks = make(map[string]Concept)
	trans.inProgressCalls = make(map[string]bool)
	trans.uOfD = uOfDPtr
	return &trans
}

// NewLiteral creates and initializes a new Literal
func (uOfDPtr *UniverseOfDiscourse) NewLiteral(trans *Transaction, uri ...string) (Concept, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var lit concept
	lit.initializeConcept(Literal, conceptID, actualURI)
	trans.WriteLockElement(&lit)
	uOfDPtr.SetUniverseOfDiscourse(&lit, trans)
	if actualURI != "" {
		lit.SetURI(actualURI, trans)
	}
	return &lit, nil
}

// NewOwnedElement creates an element (with optional URI) and sets its owner and label
func (uOfDPtr *UniverseOfDiscourse) NewOwnedElement(owner Concept, label string, trans *Transaction, uri ...string) (Concept, error) {
	el, err := uOfDPtr.NewElement(trans, uri...)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedElement failed")
	}
	err = el.SetLabel(label, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedElement failed")
	}
	err = el.SetOwningConcept(owner, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedElement failed")
	}
	return el, nil
}

// NewOwnedLiteral creates a literal (with optional URI) and sets its owner and label
func (uOfDPtr *UniverseOfDiscourse) NewOwnedLiteral(owner Concept, label string, trans *Transaction, uri ...string) (Concept, error) {
	lit, err := uOfDPtr.NewLiteral(trans, uri...)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedLiteral failed")
	}
	err = lit.SetLabel(label, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedLiteral failed")
	}
	err = lit.SetOwningConcept(owner, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedLiteral failed")
	}
	return lit, nil
}

// NewOwnedReference creates a reference (with optional URI) and sets its owner and label
func (uOfDPtr *UniverseOfDiscourse) NewOwnedReference(owner Concept, label string, trans *Transaction, uri ...string) (Concept, error) {
	ref, err := uOfDPtr.NewReference(trans, uri...)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedReference failed")
	}
	err = ref.SetLabel(label, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedReference failed")
	}
	err = ref.SetOwningConcept(owner, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedReference failed")
	}
	return ref, nil
}

// NewOwnedRefinement creates a refinement (with optional URI) and sets its owner and label
func (uOfDPtr *UniverseOfDiscourse) NewOwnedRefinement(owner Concept, label string, abstractConcept Concept, refinedConcept Concept, trans *Transaction, uri ...string) (Concept, error) {
	ref, err := uOfDPtr.NewRefinement(trans, uri...)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	err = ref.SetLabel(label, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	err = ref.SetOwningConcept(owner, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	err = ref.SetAbstractConcept(abstractConcept, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	err = ref.SetRefinedConcept(refinedConcept, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	return ref, nil
}

// NewCompleteRefinement creates a refinement (with optional URI) and sets its abstract and refined references, sets the label, and
// makes the refined element the owner
func (uOfDPtr *UniverseOfDiscourse) NewCompleteRefinement(abstractElement Concept, refinedElement Concept, label string, trans *Transaction, uri ...string) (Concept, error) {
	ref, err := uOfDPtr.NewRefinement(trans, uri...)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	err = ref.SetLabel(label, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	err = ref.SetAbstractConcept(abstractElement, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	err = ref.SetRefinedConcept(refinedElement, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	err = ref.SetOwningConcept(refinedElement, trans)
	if err != nil {
		return nil, errors.Wrap(err, "UniverseOfDiscourse.NewOwnedRefinement failed")
	}
	return ref, nil
}

// NewReference creates and initializes a new Reference
func (uOfDPtr *UniverseOfDiscourse) NewReference(trans *Transaction, uri ...string) (Concept, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var ref concept
	ref.initializeConcept(Reference, conceptID, actualURI)
	trans.WriteLockElement(&ref)
	uOfDPtr.SetUniverseOfDiscourse(&ref, trans)
	if actualURI != "" {
		ref.SetURI(actualURI, trans)
	}
	return &ref, nil
}

// NewRefinement creates and initializes a new Refinement
func (uOfDPtr *UniverseOfDiscourse) NewRefinement(trans *Transaction, uri ...string) (Concept, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var ref concept
	ref.initializeConcept(Refinement, conceptID, actualURI)
	trans.WriteLockElement(&ref)
	uOfDPtr.SetUniverseOfDiscourse(&ref, trans)
	if actualURI != "" {
		ref.SetURI(actualURI, trans)
	}
	return &ref, nil
}

// NotifyUofDObservers passes the notification to all registered Observers
func (uOfDPtr *UniverseOfDiscourse) NotifyUofDObservers(notification *ChangeNotification, trans *Transaction) error {
	it := uOfDPtr.observers.Iterator()
	for observer := range it.C {
		err := observer.(Observer).Update(notification, trans)
		if err != nil {
			it.Stop()
			return errors.Wrap(err, "element.NotifyUofDObservers failed")
		}
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) callAssociatedFunctions(el Concept, notification *ChangeNotification, trans *Transaction) error {
	if el == nil {
		return errors.New("UniverseOfDiscourse.queueFunctionExecution called with a nil Element")
	}
	if el.GetUniverseOfDiscourse(trans) == nil {
		// Functions do not get executed on elements that are no longer in a Universe of Discourse
		return nil
	}
	if notification.GetNatureOfChange() == 0 {
		return errors.New("UniverseOfDiscourse.callAssociatedFunctions called without of NatureOfChange")
	}
	functionIdentifiers := uOfDPtr.findFunctions(el, notification, trans)
	for _, functionIdentifier := range functionIdentifiers {
		if TraceLocks || TraceChange {
			omitTrace := (OmitManageTreeNodesCalls && functionIdentifier == "http://activeCrl.com/crlEditor/Editor/TreeViews/ManageTreeNodes") ||
				(OmitDiagramRelatedCalls && isDiagramRelatedFunction(functionIdentifier))
			if !omitTrace {
				log.Printf("Calling function with URI: %s notification: %s target: %p", functionIdentifier, notification.GetNatureOfChange().String(), el)
				notification.Print("      Notification: ", trans)
				log.Printf("  Function target: %T %s %s %p", el, el.getConceptIDNoLock(), el.GetLabel(trans), el)
			}
		}
		err := trans.callFunctions(functionIdentifier, el, notification)
		if err != nil {
			return errors.Wrap(err, "UniverseOfDiscourse.callAssociatedFunctions failed")
		}
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) postChange(el Concept, trans *Transaction) {
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("   After change:")
		Print(el, "      ", trans)
		log.Printf("   AfterOwnedElements: %v", uOfDPtr.ownedIDsMap.GetMappedValues(el.getConceptIDNoLock()))
		log.Printf("   After Listeners: %v", uOfDPtr.listenersMap.GetMappedValues(el.getConceptIDNoLock()))
	}
}

func (uOfDPtr *UniverseOfDiscourse) preChange(el Concept, trans *Transaction) {
	if el != nil && uOfDPtr.IsRecordingUndo() {
		uOfDPtr.undoManager.markChangedElement(el, trans)
	}
}

// RecoverDomain reconstructs a concept space from its JSON representation
func (uOfDPtr *UniverseOfDiscourse) RecoverDomain(data []byte, trans *Transaction) (Concept, error) {
	var unmarshaledData []json.RawMessage
	var conceptSpace Concept
	err := json.Unmarshal(data, &unmarshaledData)
	if err != nil {
		return nil, err
	}
	for _, data := range unmarshaledData {
		var el Concept
		el, err = uOfDPtr.RecoverElement(data, trans)
		if err != nil {
			return nil, err
		}
		if el.GetOwningConceptID(trans) == "" {
			if conceptSpace == nil {
				conceptSpace = el
			} else {
				log.Printf("In UniverseOfDiscourse.RecoverDomain more than one element does not have an owner: %s %s", el.GetLabel(trans), el.GetConceptID(trans))
			}
		}
	}
	return conceptSpace, nil
}

// RecoverElement reconstructs an Element (or subclass) from its JSON representation
func (uOfDPtr *UniverseOfDiscourse) RecoverElement(data []byte, trans *Transaction) (Concept, error) {
	if len(data) == 0 {
		err := errors.New("RecoverElement called with no data")
		return nil, err
	}
	var recoveredElement *concept = &concept{}
	recoveredElement.Version = newVersionCounter()
	recoveredElement.observers = mapset.NewSet()

	err := recoveredElement.UnmarshalJSON(data)
	if err != nil {
		log.Printf("Error recovering Element: %s \n", err)
		return nil, err
	}
	recoveredElement.uOfD = uOfDPtr
	uOfDPtr.addElement(recoveredElement, true, trans)
	return recoveredElement, nil
}

// Redo redoes the last undo, if any
func (uOfDPtr *UniverseOfDiscourse) Redo(trans *Transaction) {
	uOfDPtr.undoManager.redo(trans)
}

func (uOfDPtr *UniverseOfDiscourse) removeElementForUndo(el Concept, trans *Transaction) {
	if el != nil {
		trans.ReadLockElement(el)
		elID := el.GetConceptID(trans)
		if uOfDPtr.undoManager.debugUndo {
			log.Printf("Removing element for undo, id: %s\n", elID)
			Print(el, "Removed Element: ", trans)
		}
		uri := el.GetURI(trans)
		if uri != "" {
			uOfDPtr.uriUUIDMap.DeleteEntry(uri)
		}
		uOfDPtr.deleteUUIDElementMapEntry(elID)
		// Remove element from all listener's lists
		switch el.GetConceptType() {
		case Reference:
			referencedConceptID := el.GetReferencedConceptID(trans)
			if referencedConceptID != "" {
				uOfDPtr.removeMappedValueFromListenersMap(referencedConceptID, elID)
			}
		case Refinement:
			abstractConceptID := el.GetAbstractConceptID(trans)
			if abstractConceptID != "" {
				uOfDPtr.removeMappedValueFromListenersMap(abstractConceptID, elID)
			}
			refinedConceptID := el.GetRefinedConceptID(trans)
			if refinedConceptID != "" {
				uOfDPtr.removeMappedValueFromListenersMap(refinedConceptID, elID)
			}
		}
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeMappedValueFromListenersMap(listenerID string, value string) {
	uOfDPtr.listenersMap.removeMappedValue(listenerID, value)
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("      removeMappedValueFromListenersMap listenerID: %s value: %s resultingOwnedIDs: %v", listenerID, value, uOfDPtr.listenersMap.GetMappedValues(listenerID))
	}
}

func (uOfDPtr *UniverseOfDiscourse) removeMappedValueFromOwnedIDsMap(ownerID string, ownedConceptID string) {
	uOfDPtr.ownedIDsMap.removeMappedValue(ownerID, ownedConceptID)
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("      removeMappedValueFromOwnedIDsMap ownerID: %s removedConceptID: %s resultingOwnedIDs: %v", ownerID, ownedConceptID, uOfDPtr.ownedIDsMap.GetMappedValues(ownerID))
	}
}

// replicateAsRefinement replicates the structure of the original in the replicate, ignoring
// Refinements The name from each original element is copied into the name of the
// corresponding replicate element. Most attributes
// are not replicated, specifically any pointers, ReadOnly, Definition, IsCore, Version, and observers.
// This function is idempotent: if applied to an existing structure,
// Elements of that structure that have existing Refinement relationships with original Elements
// will not be re-created.
func (uOfDPtr *UniverseOfDiscourse) replicateAsRefinement(original Concept, replicate Concept, trans *Transaction, uri ...string) error {
	trans.ReadLockElement(original)
	trans.WriteLockElement(replicate)

	// Set the attributes - but no IDs
	err := replicate.SetLabel("Instance of "+original.GetLabel(trans), trans)
	if err != nil {
		return errors.Wrap(err, "UniverseOfDiscourse.replicateAsRefinement replicate.SetLabel failed")
	}
	switch original.GetConceptType() {
	case Reference:
		switch replicate.GetConceptType() {
		case Reference:
			replicate.SetReferencedConcept(nil, original.GetReferencedAttributeName(trans), trans)
		}
	}

	// Determine whether there is already a refinement in place; if not, create it
	if !replicate.IsRefinementOf(original, trans) {
		refinementURI := ""
		if len(uri) == 1 && uri[0] != "" {
			refinementURI = uri[0] + original.GetConceptID(trans) + "/Refinement"
		}
		refinement, err := uOfDPtr.NewRefinement(trans, refinementURI)
		if err != nil {
			return errors.Wrap(err, "UniverseOfDiscourse.replicateAsRefinement failed: ")
		}
		refinement.SetOwningConcept(replicate, trans)
		refinement.SetAbstractConcept(original, trans)
		refinement.SetRefinedConcept(replicate, trans)
		refinement.SetLabel("Refines "+original.GetLabel(trans), trans)
	}

	// Now determine which children need to be replicated
	originalID := original.GetConceptID(trans)
	replicateID := replicate.GetConceptID(trans)
	it := uOfDPtr.GetConceptsOwnedConceptIDs(originalID).Iterator()
	newChildCount := 0
	for id := range it.C {
		newChildURI := ""
		originalChild := uOfDPtr.GetElement(id.(string))
		switch originalChild.GetConceptType() {
		case Refinement:
			continue
		}
		var replicateChild Concept
		// For each original child, determine whether there is already a replicate child that
		// has the original child as one of its abstractions. This is replicateChild
		it2 := uOfDPtr.GetConceptsOwnedConceptIDs(replicateID).Iterator()
		for id := range it2.C {
			currentChild := uOfDPtr.GetElement(id.(string))
			switch currentChild.GetConceptType() {
			case Refinement:
				continue
			}
			currentChildAbstractions := make(map[string]Concept)
			currentChild.FindAbstractions(currentChildAbstractions, trans)
			for _, currentChildAbstraction := range currentChildAbstractions {
				if currentChildAbstraction == originalChild {
					replicateChild = currentChild
				}
			}
		}
		// If the replicate child is nil at this point, there is no existing replicate child that corresponds
		// to the original child - create one.
		if replicateChild == nil {
			newChildCount++
			if uri != nil && uri[0] != "" {
				newChildURI = uri[0] + ".child" + strconv.Itoa(newChildCount)
			} else {
				newChildURI = ""
			}
			var replicateError error
			replicateChild, replicateError = uOfDPtr.NewConcept(originalChild.GetConceptType(), trans, newChildURI)
			if replicateError != nil {
				return errors.Wrap(replicateError, "UniverseOfDiscourse.replicateAsRefinement failed: ")
			}
			if replicateChild != nil {
				replicateChild.SetOwningConcept(replicate, trans)
				err := uOfDPtr.replicateAsRefinement(originalChild, replicateChild, trans, newChildURI)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Register adds the registration of an Observer
func (uOfDPtr *UniverseOfDiscourse) Register(observer Observer) error {
	uOfDPtr.observers.Add(observer)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) setMappedValuesForListenersMap(listenerID string, value mapset.Set) {
	uOfDPtr.listenersMap.SetMappedValues(listenerID, value)
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("      setMappedValuesForListenersMap listenerID: %s value: %v", listenerID, value)
	}
}

func (uOfDPtr *UniverseOfDiscourse) setOwnedIDsMapValues(currentOwnerID string, values mapset.Set) {
	if uOfDPtr.undoManager.debugUndo {
		log.Printf("About to set ownedIDsMap values ownerID: %v values: %v", currentOwnerID, values)
	}
	uOfDPtr.ownedIDsMap.SetMappedValues(currentOwnerID, values)
}

// SetRecordingUndo turns undo/redo recording on and off
func (uOfDPtr *UniverseOfDiscourse) SetRecordingUndo(newSetting bool) {
	uOfDPtr.undoManager.setRecordingUndo(newSetting)
}

// SetUniverseOfDiscourse sets the uOfD of which this element is a member. Strictly
// speaking, this is not an attribute of the elment, but rather a context in which
// the element is operating in which the element may be able to locate other objects
// by id.
func (uOfDPtr *UniverseOfDiscourse) SetUniverseOfDiscourse(el Concept, trans *Transaction) error {
	trans.WriteLockElement(el)
	currentUofD := el.GetUniverseOfDiscourse(trans)
	if currentUofD != uOfDPtr {
		if el.GetIsCore(trans) {
			return errors.New("SetUniverseOfDiscourse called on a CRL core concept")
		}
		if currentUofD != nil {
			return errors.New("SetUniverseOfDiscourse called on an Element in another uOfD")
		}
		if el.IsReadOnly(trans) {
			return errors.New("SetUniverseOfDiscourse called on read-only Element")
		}
		el.setUniverseOfDiscourse(uOfDPtr, trans)
		uOfDPtr.addElement(el, false, trans)
		elementState, err := NewConceptState(el)
		if err != nil {
			return errors.Wrap(err, "UniverseOfDiscourse.SetUniverseOfDiscourse failed")
		}
		conceptAddedNotification := uOfDPtr.newUofDConceptAddedNotification(elementState, trans)
		uOfDPtr.NotifyUofDObservers(conceptAddedNotification, trans)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) setUUIDElementMapEntry(conceptID string, entry Concept) {
	if uOfDPtr.undoManager.debugUndo {
		serializedElement, _ := entry.MarshalJSON()
		log.Printf("About to add UUIDElementMapEntry ID: %s, entry: %v", conceptID, string(serializedElement))
	}
	uOfDPtr.uuidElementMap.SetEntry(conceptID, entry)

}

// StartDebugUndo starts the logging of undo/redo activity
func (uOfDPtr *UniverseOfDiscourse) StartDebugUndo() {
	uOfDPtr.undoManager.debugUndo = true
}

// StopDebugUndo ends the logging of undo/redo activity
func (uOfDPtr *UniverseOfDiscourse) StopDebugUndo() {
	uOfDPtr.undoManager.debugUndo = false
}

// Undo undoes all the changes up to the last UndoMarker or the beginning of Undo, whichever comes first.
func (uOfDPtr *UniverseOfDiscourse) Undo(trans *Transaction) {
	uOfDPtr.undoManager.undo(trans)
}

func (uOfDPtr *UniverseOfDiscourse) uriValidForConceptID(uri ...string) error {
	if len(uri) == 0 {
		return nil
	}
	if len(uri) == 1 {
		_, err := url.ParseRequestURI(uri[0])
		if err != nil {
			return errors.New("Invalid uri provided for Element conceptID: " + uri[0])
		}
		id := uuid.NewV5(uuid.NamespaceURL, uri[0]).String()
		if uOfDPtr.GetElement(id) != nil {
			return errors.New("A conceptID already exists for URI: " + uri[0])
		}
	}
	if len(uri) > 1 {
		return errors.New("Too many values provided for URI")
	}
	return nil
}
