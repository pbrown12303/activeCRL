package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	mapset "github.com/deckarep/golang-set"
	uuid "github.com/satori/go.uuid"
)

// UniverseOfDiscourse represents the scope of relevant concepts
type UniverseOfDiscourse struct {
	element
	computeFunctions functions
	executedCalls    chan *pendingFunctionCall
	undoManager      *undoManager
	uriUUIDMap       *StringStringMap
	uuidElementMap   *StringElementMap
	ownedIDsMap      *OneToNStringMap
	listenersMap     *OneToNStringMap
	abstractionsMap  *OneToNStringMap
}

// NewUniverseOfDiscourse creates and initializes a new UniverseOfDiscourse
func NewUniverseOfDiscourse() *UniverseOfDiscourse {
	var uOfD UniverseOfDiscourse
	uOfD.computeFunctions = make(map[string][]crlExecutionFunction)
	uOfD.undoManager = newUndoManager(&uOfD)
	uOfD.uriUUIDMap = NewStringStringMap()
	uOfD.uuidElementMap = NewStringElementMap()
	uOfD.ownedIDsMap = NewOneToNStringMap()
	uOfD.listenersMap = NewOneToNStringMap()
	uOfD.abstractionsMap = NewOneToNStringMap()
	uOfDID, _ := uOfD.generateConceptID(UniverseOfDiscourseURI)
	uOfD.initializeElement(uOfDID, UniverseOfDiscourseURI)
	uOfD.Label = "UniverseOfDiscourse"
	uOfD.uOfD = &uOfD
	hl := uOfD.NewHeldLocks()
	uOfD.IsCore = true
	uOfD.addElement(&uOfD, false, hl)
	uOfD.AddFunction(coreHousekeepingURI, coreHousekeeping)
	hl.ReleaseLocksAndWait()
	buildCoreConceptSpace(&uOfD, hl)
	hl.ReleaseLocksAndWait()
	return &uOfD
}

func (uOfDPtr *UniverseOfDiscourse) addElement(el Element, inRecovery bool, hl *HeldLocks) error {
	if el == nil {
		return errors.New("UniverseOfDiscource addElement() failed because element was nil")
	}
	hl.WriteLockElement(el)
	uOfDPtr.undoManager.markNewElement(el, hl)
	uuid := el.GetConceptID(hl)
	if uuid == "" {
		return errors.New("UniverseOfDiscource addElement() failed because UUID was nil")
	}
	uOfDPtr.uuidElementMap.SetEntry(el.getConceptIDNoLock(), el)
	uri := el.GetURI(hl)
	if uri != "" {
		uOfDPtr.uriUUIDMap.SetEntry(uri, uuid)
	}
	ownerID := el.GetOwningConceptID(hl)
	if ownerID != "" {
		uOfDPtr.ownedIDsMap.AddMappedValue(ownerID, uuid)
	}
	// Add element to all listener's lists
	switch el.(type) {
	case *reference:
		ref := el.(*reference)
		referencedConceptID := ref.GetReferencedConceptID(hl)
		if referencedConceptID != "" {
			uOfDPtr.listenersMap.AddMappedValue(referencedConceptID, uuid)
		}
	case *refinement:
		ref := el.(*refinement)
		abstractConceptID := ref.GetAbstractConceptID(hl)
		if abstractConceptID != "" {
			uOfDPtr.listenersMap.AddMappedValue(abstractConceptID, uuid)
		}
		refinedConceptID := ref.GetRefinedConceptID(hl)
		if refinedConceptID != "" {
			uOfDPtr.listenersMap.AddMappedValue(refinedConceptID, uuid)
		}
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) addElementForUndo(el Element, hl *HeldLocks) error {
	if el == nil {
		return errors.New("UniverseOfDiscource addElementForUndo() failed because element was nil")
	}
	hl.WriteLockElement(el)
	if uOfDPtr.undoManager.debugUndo == true {
		log.Printf("Adding element for undo, id: %s\n", el.GetConceptID(hl))
		Print(el, "Added Element: ", hl)
	}
	uOfDPtr.uuidElementMap.SetEntry(el.GetConceptID(hl), el)
	uri := el.GetURI(hl)
	if uri != "" {
		uOfDPtr.uriUUIDMap.SetEntry(uri, el.GetConceptID(hl))
	}
	return nil
}

// AddFunction registers a function with the indicated uri
func (uOfDPtr *UniverseOfDiscourse) AddFunction(uri string, function crlExecutionFunction) {
	uOfDPtr.computeFunctions[string(uri)] = append(uOfDPtr.computeFunctions[string(uri)], function)
}

func (uOfDPtr *UniverseOfDiscourse) changeURIForElement(el Element, oldURI string, newURI string) error {
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

// CreateReplicateAsRefinement replicates the indicated Element and all of its descendent Elements
// except that descendant Refinements are not replicated.
// For each replicated Element, a Refinement is created with the abstractElement being the original and the refinedElement
// being the replica. The root replicated element is returned.
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateAsRefinement(original Element, hl *HeldLocks, newURI ...string) Element {
	uri := ""
	if len(newURI) > 0 {
		uri = newURI[0]
	}
	var replicate Element
	switch original.(type) {
	case Literal:
		replicate, _ = uOfDPtr.NewLiteral(hl, uri)
	case Reference:
		replicate, _ = uOfDPtr.NewReference(hl, uri)
	case Refinement:
		replicate, _ = uOfDPtr.NewRefinement(hl, uri)
	case Element:
		replicate, _ = uOfDPtr.NewElement(hl, uri)
	}
	uOfDPtr.replicateAsRefinement(original, replicate, hl)
	return replicate
}

// CreateReplicateAsRefinementFromURI replicates the Element indicated by the URI
func (uOfDPtr *UniverseOfDiscourse) CreateReplicateAsRefinementFromURI(originalURI string, hl *HeldLocks, newURI ...string) (Element, error) {
	original := uOfDPtr.GetElementWithURI(originalURI)
	if original == nil {
		return nil, fmt.Errorf("In CreateReplicateAsRefinementFromURI Element with uri %s not found", originalURI)
	}
	return uOfDPtr.CreateReplicateAsRefinement(original, hl, newURI...), nil
}

func (uOfDPtr *UniverseOfDiscourse) findFunctions(element Element, notification *ChangeNotification, hl *HeldLocks) []string {
	var functionIdentifiers []string
	if element == nil {
		return functionIdentifiers
	}
	// Always add coreHouskeeping
	functionIdentifiers = append(functionIdentifiers, coreHousekeepingURI)
	// Now find functions associated with abstractions
	abstractions := make(map[string]Element)
	element.FindAbstractions(abstractions, hl)
	for _, abstraction := range abstractions {
		uri := abstraction.GetURI(hl)
		if uri != "" {
			functions := uOfDPtr.computeFunctions[uri]
			if functions != nil {
				functionIdentifiers = append(functionIdentifiers, uri)
			}
		}
	}
	return functionIdentifiers
}

func (uOfDPtr *UniverseOfDiscourse) deleteElement(el Element, deletedElements mapset.Set, hl *HeldLocks) error {
	if el == nil {
		return errors.New("UniverseOfDiscource removeElement failed elcause Element was nil")
	}
	hl.WriteLockElement(el)
	uOfDPtr.undoManager.markRemovedElement(el, hl)
	uuid := el.GetConceptID(hl)
	uri := el.GetURI(hl)
	if uri != "" {
		uOfDPtr.uriUUIDMap.DeleteEntry(uri)
	}
	// Remove element from owner's child list
	ownerID := el.GetOwningConceptID(hl)
	if ownerID != "" {
		el.SetOwningConceptID("", hl)
	}
	for id := range uOfDPtr.listenersMap.GetMappedValues(uuid).Iterator().C {
		listener := uOfDPtr.GetElement(id.(string))
		switch listener.(type) {
		case Reference:
			listener.(Reference).SetReferencedConcept(nil, hl)
		case Refinement:
			refinement := listener.(Refinement)
			if refinement.GetAbstractConcept(hl) == el {
				refinement.SetAbstractConcept(nil, hl)
			} else if refinement.GetRefinedConcept(hl) == el {
				refinement.SetRefinedConcept(nil, hl)
			}
		}
	}
	// Remove element from all listener's lists
	switch el.(type) {
	case *reference:
		ref := el.(*reference)
		referencedConceptID := ref.GetReferencedConceptID(hl)
		if referencedConceptID != "" {
			uOfDPtr.listenersMap.RemoveMappedValue(referencedConceptID, uuid)
		}
	case *refinement:
		ref := el.(*refinement)
		abstractConceptID := ref.GetAbstractConceptID(hl)
		if abstractConceptID != "" {
			uOfDPtr.listenersMap.RemoveMappedValue(abstractConceptID, uuid)
		}
		refinedConceptID := ref.GetRefinedConceptID(hl)
		if refinedConceptID != "" {
			uOfDPtr.listenersMap.RemoveMappedValue(refinedConceptID, uuid)
		}
	}
	uOfDPtr.listenersMap.DeleteKey(uuid)
	uOfDPtr.abstractionsMap.DeleteKey(uuid)
	uOfDPtr.ownedIDsMap.DeleteKey(uuid)
	uOfDPtr.uuidElementMap.DeleteEntry(uuid)
	el.setUniverseOfDiscourse(nil, hl)
	return nil
}

// DeleteElement removes a single element and its descentants from the uOfD. Pointers to the elements from other elements are set to nil.
func (uOfDPtr *UniverseOfDiscourse) DeleteElement(element Element, hl *HeldLocks) error {
	id := element.GetConceptID(hl)
	elements := mapset.NewSet(id)
	uOfDPtr.GetOwnedConceptIDsRecursively(id, elements, hl)
	return uOfDPtr.DeleteElements(elements, hl)
}

// DeleteElements removes the elements from the uOfD. Pointers to the elements from elements not being deleted are set to nil.
func (uOfDPtr *UniverseOfDiscourse) DeleteElements(elements mapset.Set, hl *HeldLocks) error {
	for id := range elements.Iterator().C {
		el := uOfDPtr.GetElement(id.(string))
		if el.GetIsCore(hl) {
			return errors.New("ClearUniverseOfDiscourse called on a CRL core concept")
		}
		if el.GetUniverseOfDiscourse(hl).getConceptIDNoLock() != uOfDPtr.getConceptIDNoLock() {
			return errors.New("ClearUniverseOfDiscourse called on an Element in a different UofD")
		}
		if el.IsReadOnly(hl) {
			return errors.New("SetUniverseOfDiscourse called on read-only Element")
		}
	}
	for id := range elements.Iterator().C {
		el := uOfDPtr.GetElement(id.(string))
		hl.WriteLockElement(el)
		uOfDPtr.preChange(el, hl)
		uOfDPtr.queueFunctionExecutions(uOfDPtr, uOfDPtr.newConceptRemovedNotification(el, hl), hl)
		uOfDPtr.deleteElement(el, elements, hl)
		el.setUniverseOfDiscourse(nil, hl)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) generateConceptID(uri ...string) (string, error) {
	var conceptID string
	if len(uri) == 0 || (len(uri) == 1 && uri[0] == "") {
		conceptID = uuid.NewV4().String()
	} else {
		if len(uri) == 1 {
			_, err := url.ParseRequestURI(uri[0])
			if err != nil {
				return "", errors.New("Invalid URI provided for initializing lement")
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
func (uOfDPtr *UniverseOfDiscourse) getComputeFunctions() *functions {
	return &uOfDPtr.computeFunctions
}

// GetElement returns the Element with the conceptID
func (uOfDPtr *UniverseOfDiscourse) GetElement(conceptID string) Element {
	return uOfDPtr.uuidElementMap.GetEntry(conceptID)
}

// GetElements returns the Elements in the uOfD mapped by their ConceptIDs
func (uOfDPtr *UniverseOfDiscourse) GetElements() map[string]Element {
	return uOfDPtr.uuidElementMap.CopyMap()
}

// GetElementWithURI returns the Element with the given URI
func (uOfDPtr *UniverseOfDiscourse) GetElementWithURI(uri string) Element {
	return uOfDPtr.GetElement(uOfDPtr.uriUUIDMap.GetEntry(uri))
}

func (uOfDPtr *UniverseOfDiscourse) getExecutedCalls() chan *pendingFunctionCall {
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

// getListenerIDs returns the set of listener IDs for the indicated ID
func (uOfDPtr *UniverseOfDiscourse) getListenerIDs(id string) mapset.Set {
	return uOfDPtr.listenersMap.GetMappedValues(id)
}

// GetLiteral returns the literal with the indicated ID (if found)
func (uOfDPtr *UniverseOfDiscourse) GetLiteral(conceptID string) Literal {
	el := uOfDPtr.GetElement(conceptID)
	switch el.(type) {
	case *literal:
		return el.(*literal)
	}
	return nil
}

// GetLiteralWithURI returns the literal with the indicated URI (if found)
func (uOfDPtr *UniverseOfDiscourse) GetLiteralWithURI(uri string) Literal {
	el := uOfDPtr.GetElementWithURI(uri)
	switch el.(type) {
	case *literal:
		return el.(*literal)
	}
	return nil
}

// GetOwnedConceptIDs returns the set of owned concepts for the indicated ID
func (uOfDPtr *UniverseOfDiscourse) GetOwnedConceptIDs(id string) mapset.Set {
	return uOfDPtr.ownedIDsMap.GetMappedValues(id)
}

// GetOwnedConceptIDsRecursively returns the IDs of owned concepts
func (uOfDPtr *UniverseOfDiscourse) GetOwnedConceptIDsRecursively(rootID string, descendants mapset.Set, hl *HeldLocks) {
	for id := range uOfDPtr.ownedIDsMap.GetMappedValues(rootID).Iterator().C {
		descendants.Add(id.(string))
		uOfDPtr.GetOwnedConceptIDsRecursively(id.(string), descendants, hl)
	}
}

// GetReference returns the reference with the indicated ID (if found)
func (uOfDPtr *UniverseOfDiscourse) GetReference(conceptID string) Reference {
	el := uOfDPtr.GetElement(conceptID)
	switch el.(type) {
	case *reference:
		return el.(*reference)
	}
	return nil
}

// GetReferenceWithURI returns the reference with the indicated URI (if found)
func (uOfDPtr *UniverseOfDiscourse) GetReferenceWithURI(uri string) Reference {
	el := uOfDPtr.GetElementWithURI(uri)
	switch el.(type) {
	case *reference:
		return el.(*reference)
	}
	return nil
}

// GetRefinement returns the refinement with the indicated ID (if found)
func (uOfDPtr *UniverseOfDiscourse) GetRefinement(conceptID string) Refinement {
	el := uOfDPtr.GetElement(conceptID)
	switch el.(type) {
	case *refinement:
		return el.(*refinement)
	}
	return nil
}

// GetRefinementWithURI returns the refinement with the indicated URI (if found)
func (uOfDPtr *UniverseOfDiscourse) GetRefinementWithURI(uri string) Refinement {
	el := uOfDPtr.GetElementWithURI(uri)
	switch el.(type) {
	case *refinement:
		return el.(*refinement)
	}
	return nil
}

// GetRootElements returns all elements that do not have owners
func (uOfDPtr *UniverseOfDiscourse) GetRootElements(hl *HeldLocks) map[string]Element {
	allElements := uOfDPtr.GetElements()
	rootElements := make(map[string]Element)
	for id, el := range allElements {
		if el.GetOwningConceptID(hl) == "" {
			rootElements[id] = el
		}
	}
	return rootElements
}

func (uOfDPtr *UniverseOfDiscourse) getURIUUIDMap() *StringStringMap {
	return uOfDPtr.uriUUIDMap
}

// IsRecordingUndo reveals whether undo recording is on
func (uOfDPtr *UniverseOfDiscourse) IsRecordingUndo() bool {
	return uOfDPtr.undoManager.recordingUndo
}

// MarkUndoPoint marks a point on the undo stack. The next undo operation will undo everything back to this point.
func (uOfDPtr *UniverseOfDiscourse) MarkUndoPoint() {
	uOfDPtr.undoManager.MarkUndoPoint()
}

// MarshalConceptSpace creates a JSON representation of an element and all of its descendants
func (uOfDPtr *UniverseOfDiscourse) MarshalConceptSpace(el Element, hl *HeldLocks) ([]byte, error) {
	var result []byte
	result = append(result, []byte("[")...)
	marshaledConcept, err := uOfDPtr.marshalConceptRecursively(el, hl)
	if err != nil {
		return result, err
	}
	// The last byte of marshaledConcept is going to be a comma we don't want
	result = append(result, marshaledConcept[0:len(marshaledConcept)-1]...)
	result = append(result, []byte("]")...)
	return result, nil
}

func (uOfDPtr *UniverseOfDiscourse) marshalConceptRecursively(el Element, hl *HeldLocks) ([]byte, error) {
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
	elID := el.GetConceptID(hl)
	for id := range uOfDPtr.GetOwnedConceptIDs(elID).Iterator().C {
		child := uOfDPtr.GetElement(id.(string))
		marshaledChild, err := uOfDPtr.marshalConceptRecursively(child, hl)
		if err != nil {
			return result, err
		}
		result = append(result, marshaledChild...)
	}
	return result, nil
}

// newConceptAddedNotification creates a UniverseOfDiscourseAdded notification
func (uOfDPtr *UniverseOfDiscourse) newConceptAddedNotification(concept Element, hl *HeldLocks) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = uOfDPtr
	notification.natureOfChange = UofDConceptAdded
	notification.priorState = clone(concept, hl)
	notification.uOfD = uOfDPtr
	return &notification
}

// NewConceptChangeNotification creates a ChangeNotification that records state of the concept prior to the
// change. Note that this MUST be called prior to making any changes to the concept.
func (uOfDPtr *UniverseOfDiscourse) NewConceptChangeNotification(changingConcept Element, hl *HeldLocks) *ChangeNotification {
	// Since this function is invoked by the *element methods for Literals, References, and Refinements, we play a
	// game to get the full datatype
	correctedChangingConcept := uOfDPtr.GetElement(changingConcept.getConceptIDNoLock())
	var notification ChangeNotification
	notification.reportingElement = correctedChangingConcept
	notification.natureOfChange = ConceptChanged
	// Since this function is invoked by the *element methods for Literals, References, and Refinements, we play a
	// game to get the full datatype cloned
	notification.priorState = clone(correctedChangingConcept, hl)
	notification.uOfD = uOfDPtr
	return &notification
}

// newConceptRemovedNotification creates a UniverseOfDiscourseRemoved notification
func (uOfDPtr *UniverseOfDiscourse) newConceptRemovedNotification(concept Element, hl *HeldLocks) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = uOfDPtr
	notification.natureOfChange = UofDConceptRemoved
	notification.priorState = clone(concept, hl)
	notification.uOfD = uOfDPtr
	return &notification
}

// NewForwardingChangeNotification creates a ChangeNotification that records the reason for the change to the element,
// including the nature of the change, an indication of which component originated the change, and whether there
// was a preceeding notification that triggered this change.
func (uOfDPtr *UniverseOfDiscourse) NewForwardingChangeNotification(reportingElement Element, natureOfChange NatureOfChange, underlyingChange *ChangeNotification) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = reportingElement
	notification.natureOfChange = natureOfChange
	notification.underlyingChange = underlyingChange
	notification.uOfD = uOfDPtr
	return &notification
}

// NewElement creates and initializes a new Element
func (uOfDPtr *UniverseOfDiscourse) NewElement(hl *HeldLocks, uri ...string) (Element, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var el element
	el.initializeElement(conceptID, actualURI)
	hl.WriteLockElement(&el)
	uOfDPtr.SetUniverseOfDiscourse(&el, hl)
	return &el, nil
}

// NewHeldLocks creates and initializes a HeldLocks structure utilizing the supplied WaitGroup
func (uOfDPtr *UniverseOfDiscourse) NewHeldLocks() *HeldLocks {
	var hl HeldLocks
	hl.readLocks = make(map[string]Element)
	hl.writeLocks = make(map[string]Element)
	hl.uOfD = uOfDPtr
	hl.functionCallManager = newFunctionCallManager(hl.uOfD)
	return &hl
}

// NewLiteral creates and initializes a new Literal
func (uOfDPtr *UniverseOfDiscourse) NewLiteral(hl *HeldLocks, uri ...string) (Literal, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var lit literal
	lit.initializeLiteral(conceptID, actualURI)
	hl.WriteLockElement(&lit)
	uOfDPtr.SetUniverseOfDiscourse(&lit, hl)
	return &lit, nil
}

// NewReference creates and initializes a new Reference
func (uOfDPtr *UniverseOfDiscourse) NewReference(hl *HeldLocks, uri ...string) (Reference, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var ref reference
	ref.initializeReference(conceptID, actualURI)
	hl.WriteLockElement(&ref)
	uOfDPtr.SetUniverseOfDiscourse(&ref, hl)
	return &ref, nil
}

// NewRefinement creates and initializes a new Refinement
func (uOfDPtr *UniverseOfDiscourse) NewRefinement(hl *HeldLocks, uri ...string) (Refinement, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	actualURI := ""
	if len(uri) == 1 {
		actualURI = uri[0]
	}
	var ref refinement
	ref.initializeRefinement(conceptID, actualURI)
	hl.WriteLockElement(&ref)
	uOfDPtr.SetUniverseOfDiscourse(&ref, hl)
	return &ref, nil
}

// newUniverseOfDiscourseChangeNotification creates a new ChangeNotification for a UofD change
func (uOfDPtr *UniverseOfDiscourse) newUniverseOfDiscourseChangeNotification(underlyingChange *ChangeNotification) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = uOfDPtr
	notification.natureOfChange = UofDConceptChanged
	notification.underlyingChange = underlyingChange
	notification.uOfD = uOfDPtr
	return &notification
}

func (uOfDPtr *UniverseOfDiscourse) preChange(el Element, hl *HeldLocks) {
	if el != nil && uOfDPtr.IsRecordingUndo() == true {
		uOfDPtr.undoManager.markChangedElement(el, hl)
	}
}

func (uOfDPtr *UniverseOfDiscourse) queueFunctionExecutions(el Element, notification *ChangeNotification, hl *HeldLocks) {
	if el == nil {
		log.Printf("UniverseOfDiscourse.queueFunctionExecution called with a nil Element")
		return
	}
	if el.GetUniverseOfDiscourse(hl) == nil {
		// Functions do not get executed on elements that are no longer in a Universe of Discourse
		return
	}
	if notification.GetNatureOfChange() == 0 {
		log.Printf("UniverseOfDiscourse.queueFunctionExecution called without of NatureOfChange")
		return
	}
	functionIdentifiers := uOfDPtr.findFunctions(el, notification, hl)
	for _, functionIdentifier := range functionIdentifiers {
		if TraceLocks == true || TraceChange == true {
			omitTrace := (OmitHousekeepingCalls && functionIdentifier == "http://activeCrl.com/core/coreHousekeeping") ||
				(OmitManageTreeNodesCalls && functionIdentifier == "http://activeCrl.com/crlEditor/Editor/TreeViews/ManageTreeNodes")
			if omitTrace == false {
				log.Printf("queueFunctionExecutions adding function, URI: %s notification: %s target: %p", functionIdentifier, notification.GetNatureOfChange().String(), el)
				log.Printf("Function target: %T %s %s %p", el, el.getConceptIDNoLock(), el.GetLabel(hl), el)
				notification.Print("Notification: ", hl)
			}
		}
		hl.functionCallManager.addFunctionCall(functionIdentifier, el, notification)
	}
}

// Redo redoes the last undo, if any
func (uOfDPtr *UniverseOfDiscourse) Redo(hl *HeldLocks) {
	uOfDPtr.undoManager.redo(hl)
}

func (uOfDPtr *UniverseOfDiscourse) removeElementForUndo(el Element, hl *HeldLocks) {
	if el != nil {
		hl.ReadLockElement(el)
		elID := el.GetConceptID(hl)
		if uOfDPtr.undoManager.debugUndo == true {
			log.Printf("Removing element for undo, id: %s\n", elID)
			Print(el, "Removed Element: ", hl)
		}
		uOfDPtr.uuidElementMap.DeleteEntry(elID)
		// TODO fixup for cached ownedConceptIDs, listenerIDs
	}
}

// RecoverConceptSpace reconstructs a concept space from its JSON representation
func (uOfDPtr *UniverseOfDiscourse) RecoverConceptSpace(data []byte, hl *HeldLocks) (Element, error) {
	var unmarshaledData []json.RawMessage
	var conceptSpace Element
	err := json.Unmarshal(data, &unmarshaledData)
	if err != nil {
		return nil, err
	}
	for _, data := range unmarshaledData {
		var el Element
		el, err = uOfDPtr.RecoverElement(data, hl)
		if err != nil {
			return nil, err
		}
		if el.GetOwningConceptID(hl) == "" {
			if conceptSpace == nil {
				conceptSpace = el
			} else {
				return nil, errors.New("In UniverseOfDiscourse.RecoverConceptSpace more than one element does not have an owner")
			}
		}
	}
	return conceptSpace, nil
}

// RecoverElement reconstructs an Element (or subclass) from its JSON representation
func (uOfDPtr *UniverseOfDiscourse) RecoverElement(data []byte, hl *HeldLocks) (Element, error) {
	if len(data) == 0 {
		err := errors.New("RecoverElement called with no data")
		return nil, err
	}
	var recoveredElement Element
	err := uOfDPtr.unmarshalPolymorphicElement(data, &recoveredElement, hl)
	if err != nil {
		log.Printf("Error recovering Element: %s \n", err)
		return nil, err
	}
	uOfDPtr.addElement(recoveredElement, true, hl)
	return recoveredElement, nil
}

// replicateAsRefinement replicates the structure of the original in the replicate, ignoring
// Refinements The name from each original element is copied into the name of the
// corresponding replicate element. This function is idempotent: if applied to an existing structure,
// Elements of that structure that have existing Refinement relationships with original Elements
// will not be re-created.
func (uOfDPtr *UniverseOfDiscourse) replicateAsRefinement(original Element, replicate Element, hl *HeldLocks) {
	hl.ReadLockElement(original)
	hl.WriteLockElement(replicate)

	replicate.SetLabel(original.GetLabel(hl), hl)
	if replicate.IsRefinementOf(original, hl) == false {
		refinement, _ := uOfDPtr.NewRefinement(hl)
		refinement.SetOwningConcept(replicate, hl)
		refinement.SetAbstractConcept(original, hl)
		refinement.SetRefinedConcept(replicate, hl)
		refinement.SetLabel("Refines "+original.GetLabel(hl), hl)
	}
	originalID := original.GetConceptID(hl)
	replicateID := replicate.GetConceptID(hl)
	for id := range uOfDPtr.GetOwnedConceptIDs(originalID).Iterator().C {
		originalChild := uOfDPtr.GetElement(id.(string))
		switch originalChild.(type) {
		case Refinement:
			continue
		}
		var replicateChild Element
		// For each original child, determine whether there is already a replicate child that
		// has the original child as one of its abstractions. This is replicateChild
		for id := range uOfDPtr.GetOwnedConceptIDs(replicateID).Iterator().C {
			currentChild := uOfDPtr.GetElement(id.(string))
			currentChildAbstractions := make(map[string]Element)
			currentChild.FindAbstractions(currentChildAbstractions, hl)
			for _, currentChildAbstraction := range currentChildAbstractions {
				if currentChildAbstraction == originalChild {
					replicateChild = currentChild
				}
			}
		}
		// If the replicate child is nil at this point, there is no existing replicate child that corresponds
		// to the original child - create one.
		if replicateChild == nil {
			switch originalChild.(type) {
			case Reference:
				replicateChild, _ = uOfDPtr.NewReference(hl)
			case Literal:
				replicateChild, _ = uOfDPtr.NewLiteral(hl)
			case Element:
				replicateChild, _ = uOfDPtr.NewElement(hl)
			}
			replicateChild.SetOwningConcept(replicate, hl)
			refinement, _ := uOfDPtr.NewRefinement(hl)
			refinement.SetOwningConcept(replicateChild, hl)
			refinement.SetAbstractConcept(originalChild, hl)
			refinement.SetRefinedConcept(replicateChild, hl)
			refinement.SetLabel("Refines "+originalChild.GetLabel(hl), hl)
			replicateChild.SetLabel(originalChild.GetLabel(hl), hl)
		}
		switch originalChild.(type) {
		case Element:
			uOfDPtr.replicateAsRefinement(originalChild.(Element), replicateChild.(Element), hl)
		}
	}
}

// SetRecordingUndo turns undo/redo recording on and off
func (uOfDPtr *UniverseOfDiscourse) SetRecordingUndo(newSetting bool) {
	uOfDPtr.undoManager.setRecordingUndo(newSetting)
}

// SetUniverseOfDiscourse sets the uOfD of which this element is a member. Strictly
// speaking, this is not an attribute of the elment, but rather a context in which
// the element is operating in which the element may be able to locate other objects
// by id.
func (uOfDPtr *UniverseOfDiscourse) SetUniverseOfDiscourse(el Element, hl *HeldLocks) error {
	hl.WriteLockElement(el)
	currentUofD := el.GetUniverseOfDiscourse(hl)
	if currentUofD != uOfDPtr {
		if el.GetIsCore(hl) {
			return errors.New("SetUniverseOfDiscourse called on a CRL core concept")
		}
		if currentUofD != nil {
			return errors.New("SetUniverseOfDiscourse called on an Element in another uOfD")
		}
		if el.IsReadOnly(hl) {
			return errors.New("SetUniverseOfDiscourse called on read-only Element")
		}
		uOfDPtr.preChange(el, hl)
		uOfDPtr.queueFunctionExecutions(uOfDPtr, uOfDPtr.newConceptAddedNotification(el, hl), hl)
		el.setUniverseOfDiscourse(uOfDPtr, hl)
		uOfDPtr.addElement(el, false, hl)
	}
	return nil
}

// Undo undoes all the changes up to the last UndoMarker or the beginning of Undo, whichever comes first.
func (uOfDPtr *UniverseOfDiscourse) Undo(hl *HeldLocks) {
	uOfDPtr.undoManager.undo(hl)
}

func (uOfDPtr *UniverseOfDiscourse) unmarshalPolymorphicElement(data []byte, result *Element, hl *HeldLocks) error {
	var unmarshaledData map[string]json.RawMessage
	err := json.Unmarshal(data, &unmarshaledData)
	if err != nil {
		return err
	}
	var elementType string
	err = json.Unmarshal(unmarshaledData["Type"], &elementType)
	if err != nil {
		return err
	}
	switch elementType {
	case "*core.element":
		//		fmt.Printf("Switch choice *core.element \n")
		var recoveredElement element
		recoveredElement.uOfD = uOfDPtr
		recoveredElement.initializeElement("", "")
		*result = &recoveredElement
		err = recoveredElement.recoverElementFields(&unmarshaledData, hl)
		if err != nil {
			return err
		}
	case "*core.reference":
		//		fmt.Printf("Switch choice *core.elementReference \n")
		var recoveredReference reference
		recoveredReference.uOfD = uOfDPtr
		recoveredReference.initializeReference("", "")
		*result = &recoveredReference
		err = recoveredReference.recoverReferenceFields(&unmarshaledData, hl)
		if err != nil {
			return err
		}
	case "*core.literal":
		//		fmt.Printf("Switch choice *core.literal \n")
		var recoveredLiteral literal
		recoveredLiteral.uOfD = uOfDPtr
		recoveredLiteral.initializeLiteral("", "")
		*result = &recoveredLiteral
		err = recoveredLiteral.recoverLiteralFields(&unmarshaledData, hl)
		if err != nil {
			return err
		}
	case "*core.refinement":
		var recoveredRefinement refinement
		recoveredRefinement.uOfD = uOfDPtr
		recoveredRefinement.initializeRefinement("", "")
		*result = &recoveredRefinement
		err = recoveredRefinement.recoverRefinementFields(&unmarshaledData, hl)
		if err != nil {
			return err
		}
	default:
		log.Printf("No case for %s in unmarshalPolymorphicBaseElement \n", elementType)
	}
	return nil
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
