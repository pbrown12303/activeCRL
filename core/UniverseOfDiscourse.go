package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	uuid "github.com/satori/go.uuid"
)

// universeOfDiscourse represents the scope of relevant concepts
type universeOfDiscourse struct {
	element
	computeFunctions      functions
	executedCalls         chan *pendingFunctionCall
	undoManager           *undoManager
	unresolvedPointersMap *stringCachedPointersMap
	uriElementMap         *StringElementMap
	uuidElementMap        *StringElementMap
}

// NewUniverseOfDiscourse creates and initializes a new UniverseOfDiscourse
func NewUniverseOfDiscourse() UniverseOfDiscourse {
	var uOfD universeOfDiscourse
	uOfD.computeFunctions = make(map[string][]crlExecutionFunction)
	uOfD.undoManager = newUndoManager()
	uOfD.unresolvedPointersMap = newStringCachedPointersMap()
	uOfD.uriElementMap = NewStringElementMap()
	uOfD.uuidElementMap = NewStringElementMap()
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

func (uOfDPtr *universeOfDiscourse) addElement(el Element, inRecovery bool, hl *HeldLocks) error {
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
		uOfDPtr.uriElementMap.SetEntry(uri, el)
	}
	owner := uOfDPtr.GetElement(el.GetOwningConceptID(hl))
	if owner != nil {
		if inRecovery {
			owner.addRecoveredOwnedConcept(uuid, hl)
		} else {
			owner.addOwnedConcept(uuid, hl)
		}
	}
	// Add element to all listener's lists
	switch el.(type) {
	case *reference:
		ref := el.(*reference)
		referencedConcept := uOfDPtr.GetElement(ref.GetReferencedConceptID(hl))
		if referencedConcept != nil {
			ref.referencedConcept.setIndicatedConcept(referencedConcept)
			referencedConcept.addListener(uuid, hl)
		}
	case *refinement:
		ref := el.(*refinement)
		abstractConcept := uOfDPtr.GetElement(ref.GetAbstractConceptID(hl))
		if abstractConcept != nil {
			ref.abstractConcept.setIndicatedConcept(abstractConcept)
			abstractConcept.addListener(uuid, hl)
		}
		refinedConcept := uOfDPtr.GetElement(ref.GetRefinedConceptID(hl))
		if refinedConcept != nil {
			ref.refinedConcept.setIndicatedConcept(refinedConcept)
			refinedConcept.addListener(uuid, hl)
		}
	}
	uOfDPtr.cacheUnresolvedPointers(el, hl)
	uOfDPtr.unresolvedPointersMap.resolveCachedPointers(el, hl)
	return nil
}

func (uOfDPtr *universeOfDiscourse) addElementForUndo(el Element, hl *HeldLocks) error {
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
		uOfDPtr.uriElementMap.SetEntry(uri, el)
	}
	uOfDPtr.cacheUnresolvedPointers(el, hl)
	uOfDPtr.unresolvedPointersMap.resolveCachedPointers(el, hl)
	return nil
}

func (uOfDPtr *universeOfDiscourse) AddFunction(uri string, function crlExecutionFunction) {
	uOfDPtr.computeFunctions[string(uri)] = append(uOfDPtr.computeFunctions[string(uri)], function)
}

func (uOfDPtr *universeOfDiscourse) addUnresolvedPointer(pointer *cachedPointer) {
	uOfDPtr.unresolvedPointersMap.addCachedPointer(pointer)
}

func (uOfDPtr *universeOfDiscourse) cacheUnresolvedPointers(el Element, hl *HeldLocks) {
	if el.GetOwningConcept(hl) == nil && el.GetOwningConceptID(hl) != "" {
		uOfDPtr.addUnresolvedPointer(el.(*element).owningConcept)
	}
	switch el.(type) {
	case *reference:
		ref := el.(*reference)
		if ref.GetReferencedConcept(hl) == nil && ref.GetReferencedConceptID(hl) != "" {
			uOfDPtr.addUnresolvedPointer(ref.referencedConcept)
		}
	case *refinement:
		ref := el.(*refinement)
		if ref.GetAbstractConcept(hl) == nil && ref.GetAbstractConceptID(hl) != "" {
			uOfDPtr.addUnresolvedPointer(ref.abstractConcept)
		}
		if ref.GetRefinedConcept(hl) == nil && ref.GetRefinedConceptID(hl) != "" {
			uOfDPtr.addUnresolvedPointer(ref.refinedConcept)
		}
	}
}

func (uOfDPtr *universeOfDiscourse) changeURIForElement(el Element, oldURI string, newURI string) error {
	if oldURI != "" && uOfDPtr.uriElementMap.GetEntry(oldURI) == el {
		uOfDPtr.uriElementMap.DeleteEntry(oldURI)
	}
	if newURI != "" {
		if uOfDPtr.uriElementMap.GetEntry(newURI) != nil {
			return errors.New("Attempted to assign a URI that is already in use")
		}
		uOfDPtr.uriElementMap.SetEntry(newURI, el)
	}
	return nil
}

// CreateReplicateAsRefinement replicates the indicated Element and all of its descendent Elements
// except that descendant Refinements are not replicated.
// For each replicated Element, a Refinement is created with the abstractElement being the original and the refinedElement
// being the replica. The root replicated element is returned.
func (uOfDPtr *universeOfDiscourse) CreateReplicateAsRefinement(original Element, hl *HeldLocks, newURI ...string) Element {
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
func (uOfDPtr *universeOfDiscourse) CreateReplicateAsRefinementFromURI(originalURI string, hl *HeldLocks, newURI ...string) (Element, error) {
	original := uOfDPtr.GetElementWithURI(originalURI)
	if original == nil {
		return nil, fmt.Errorf("In CreateReplicateAsRefinementFromURI Element with uri %s not found", originalURI)
	}
	return uOfDPtr.CreateReplicateAsRefinement(original, hl, newURI...), nil
}

func (uOfDPtr *universeOfDiscourse) findFunctions(element Element, notification *ChangeNotification, hl *HeldLocks) []string {
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

func (uOfDPtr *universeOfDiscourse) deleteElement(el Element, deletedElements map[string]Element, hl *HeldLocks) error {
	if el == nil {
		return errors.New("UniverseOfDiscource removeElement failed elcause Element was nil")
	}
	hl.WriteLockElement(el)
	uOfDPtr.undoManager.markRemovedElement(el, hl)
	uuid := el.GetConceptID(hl)
	uri := el.GetURI(hl)
	if uri != "" {
		uOfDPtr.uriElementMap.DeleteEntry(uri)
	}
	// Remove element from owner's child list
	owner := el.GetOwningConcept(hl)
	if owner != nil && deletedElements[owner.GetConceptID(hl)] == nil {
		el.SetOwningConceptID("", hl)
	}
	// Remove element from all listener's lists
	switch el.(type) {
	case *reference:
		ref := el.(*reference)
		referencedConcept := ref.GetReferencedConcept(hl)
		if referencedConcept != nil {
			referencedConcept.removeListener(uuid, hl)
		}
	case *refinement:
		ref := el.(*refinement)
		abstractConcept := ref.GetAbstractConcept(hl)
		if abstractConcept != nil {
			abstractConcept.removeListener(uuid, hl)
		}
		refinedConcept := ref.GetRefinedConcept(hl)
		if refinedConcept != nil {
			refinedConcept.removeListener(uuid, hl)
		}
	}
	for _, listener := range el.getListeners(hl) {
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
	// Remove Element from all cached pointers
	uOfDPtr.uncacheUnresolvedPointers(el, hl)
	uOfDPtr.unresolveIncomingPointers(el, hl)
	uOfDPtr.uuidElementMap.DeleteEntry(uuid)
	el.setUniverseOfDiscourse(nil, hl)
	return nil
}

// DeleteElement() removes a single element and its descentants from the uOfD. Pointers to the elements from other elements are set to nil.
func (uOfDPtr *universeOfDiscourse) DeleteElement(element Element, hl *HeldLocks) error {
	elements := map[string]Element{element.GetConceptID(hl): element}
	element.GetOwnedConceptsRecursively(elements, hl)
	return uOfDPtr.DeleteElements(elements, hl)
}

// DeleteElements() removes the elements from the uOfD. Pointers to the elements from elements not being deleted are set to nil.
func (uOfDPtr *universeOfDiscourse) DeleteElements(elements map[string]Element, hl *HeldLocks) error {
	for _, el := range elements {
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
	for _, el := range elements {
		hl.WriteLockElement(el)
		uOfDPtr.preChange(el, hl)
		uOfDPtr.queueFunctionExecutions(uOfDPtr, uOfDPtr.newConceptRemovedNotification(el, hl), hl)
		uOfDPtr.deleteElement(el, elements, hl)
		el.setUniverseOfDiscourse(nil, hl)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) generateConceptID(uri ...string) (string, error) {
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
func (uOfDPtr *universeOfDiscourse) getComputeFunctions() *functions {
	return &uOfDPtr.computeFunctions
}

// GetElement returns the Element with the conceptID
func (uOfDPtr *universeOfDiscourse) GetElement(conceptID string) Element {
	return uOfDPtr.uuidElementMap.GetEntry(conceptID)
}

func (uOfDPtr *universeOfDiscourse) GetElements() map[string]Element {
	return uOfDPtr.uuidElementMap.CopyMap()
}

// GetElementWithURI returns the Element with the given URI
func (uOfDPtr *universeOfDiscourse) GetElementWithURI(uri string) Element {
	return uOfDPtr.uriElementMap.GetEntry(uri)
}

func (uOfDPtr *universeOfDiscourse) getExecutedCalls() chan *pendingFunctionCall {
	return uOfDPtr.executedCalls
}

func (uOfDPtr *universeOfDiscourse) GetFunctions(uri string) []crlExecutionFunction {
	return uOfDPtr.computeFunctions[string(uri)]
}

// GetIDForURI returns a V5 UUID derived from the given URI. If the given URI
// is not valid it returns the empty string.
func (uOfDPtr *universeOfDiscourse) GetIDForURI(uri string) string {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return ""
	}
	return uuid.NewV5(uuid.NamespaceURL, uri).String()
}

// GetLiteral returns the literal with the indicated ID (if found)
func (uOfDPtr *universeOfDiscourse) GetLiteral(conceptID string) Literal {
	el := uOfDPtr.GetElement(conceptID)
	switch el.(type) {
	case *literal:
		return el.(*literal)
	}
	return nil
}

// GetLiteralWithURI returns the literal with the indicated URI (if found)
func (uOfDPtr *universeOfDiscourse) GetLiteralWithURI(uri string) Literal {
	el := uOfDPtr.GetElementWithURI(uri)
	switch el.(type) {
	case *literal:
		return el.(*literal)
	}
	return nil
}

// GetReference returns the reference with the indicated ID (if found)
func (uOfDPtr *universeOfDiscourse) GetReference(conceptID string) Reference {
	el := uOfDPtr.GetElement(conceptID)
	switch el.(type) {
	case *reference:
		return el.(*reference)
	}
	return nil
}

// GetReferenceWithURI returns the reference with the indicated URI (if found)
func (uOfDPtr *universeOfDiscourse) GetReferenceWithURI(uri string) Reference {
	el := uOfDPtr.GetElementWithURI(uri)
	switch el.(type) {
	case *reference:
		return el.(*reference)
	}
	return nil
}

// GetRefinement returns the refinement with the indicated ID (if found)
func (uOfDPtr *universeOfDiscourse) GetRefinement(conceptID string) Refinement {
	el := uOfDPtr.GetElement(conceptID)
	switch el.(type) {
	case *refinement:
		return el.(*refinement)
	}
	return nil
}

// GetRefinementWithURI returns the refinement with the indicated URI (if found)
func (uOfDPtr *universeOfDiscourse) GetRefinementWithURI(uri string) Refinement {
	el := uOfDPtr.GetElementWithURI(uri)
	switch el.(type) {
	case *refinement:
		return el.(*refinement)
	}
	return nil
}

// GetRootElements returns all elements that do not have owners
func (uOfDPtr *universeOfDiscourse) GetRootElements(hl *HeldLocks) map[string]Element {
	allElements := uOfDPtr.GetElements()
	rootElements := make(map[string]Element)
	for id, el := range allElements {
		if el.GetOwningConceptID(hl) == "" {
			rootElements[id] = el
		}
	}
	return rootElements
}

func (uOfDPtr *universeOfDiscourse) getURIElementMap() *StringElementMap {
	return uOfDPtr.uriElementMap
}

// IsRecordingUndo reveals whether undo recording is on
func (uOfDPtr *universeOfDiscourse) IsRecordingUndo() bool {
	return uOfDPtr.undoManager.recordingUndo
}

// MarkUndoPoint marks a point on the undo stack. The next undo operation will undo everything back to this point.
func (uOfDPtr *universeOfDiscourse) MarkUndoPoint() {
	uOfDPtr.undoManager.MarkUndoPoint()
}

// MarshalConceptSpace creates a JSON representation of an element and all of its descendants
func (uOfDPtr *universeOfDiscourse) MarshalConceptSpace(el Element, hl *HeldLocks) ([]byte, error) {
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

func (uOfDPtr *universeOfDiscourse) marshalConceptRecursively(el Element, hl *HeldLocks) ([]byte, error) {
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
	for _, child := range el.GetOwnedConcepts(hl) {
		marshaledChild, err := uOfDPtr.marshalConceptRecursively(child, hl)
		if err != nil {
			return result, err
		}
		result = append(result, marshaledChild...)
	}
	return result, nil
}

// newConceptAddedNotification creates a UniverseOfDiscourseAdded notification
func (uOfDPtr *universeOfDiscourse) newConceptAddedNotification(concept Element, hl *HeldLocks) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = uOfDPtr
	notification.natureOfChange = UofDConceptAdded
	notification.priorState = clone(concept, hl)
	notification.uOfD = uOfDPtr
	return &notification
}

// NewConceptChangeNotification creates a ChangeNotification that records state of the concept prior to the
// change. Note that this MUST be called prior to making any changes to the concept.
func (uOfDPtr *universeOfDiscourse) NewConceptChangeNotification(changingConcept Element, hl *HeldLocks) *ChangeNotification {
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
func (uOfDPtr *universeOfDiscourse) newConceptRemovedNotification(concept Element, hl *HeldLocks) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = uOfDPtr
	notification.natureOfChange = UofDConceptRemoved
	notification.priorState = clone(concept, hl)
	notification.uOfD = uOfDPtr
	return &notification
}

// NewForwardingNotification creates a ChangeNotification that records the reason for the change to the element,
// including the nature of the change, an indication of which component originated the change, and whether there
// was a preceeding notification that triggered this change.
func (uOfDPtr *universeOfDiscourse) NewForwardingChangeNotification(reportingElement Element, natureOfChange NatureOfChange, underlyingChange *ChangeNotification) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = reportingElement
	notification.natureOfChange = natureOfChange
	notification.underlyingChange = underlyingChange
	notification.uOfD = uOfDPtr
	return &notification
}

func (uOfDPtr *universeOfDiscourse) NewElement(hl *HeldLocks, uri ...string) (Element, error) {
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
func (uOfDPtr *universeOfDiscourse) NewHeldLocks() *HeldLocks {
	var hl HeldLocks
	hl.readLocks = make(map[string]Element)
	hl.writeLocks = make(map[string]Element)
	hl.uOfD = uOfDPtr
	hl.functionCallManager = newFunctionCallManager(hl.uOfD)
	return &hl
}

func (uOfDPtr *universeOfDiscourse) NewLiteral(hl *HeldLocks, uri ...string) (Literal, error) {
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

func (uOfDPtr *universeOfDiscourse) NewReference(hl *HeldLocks, uri ...string) (Reference, error) {
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

func (uOfDPtr *universeOfDiscourse) NewRefinement(hl *HeldLocks, uri ...string) (Refinement, error) {
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

func (uOfDPtr *universeOfDiscourse) NewUniverseOfDiscourseChangeNotification(underlyingChange *ChangeNotification) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = uOfDPtr
	notification.natureOfChange = UofDConceptChanged
	notification.underlyingChange = underlyingChange
	notification.uOfD = uOfDPtr
	return &notification
}

func (uOfDPtr *universeOfDiscourse) preChange(el Element, hl *HeldLocks) {
	if el != nil && uOfDPtr.IsRecordingUndo() == true {
		uOfDPtr.undoManager.markChangedElement(el, hl)
	}
}

func (uOfDPtr *universeOfDiscourse) queueFunctionExecutions(el Element, notification *ChangeNotification, hl *HeldLocks) {
	if el == nil {
		log.Printf("universeOfDiscourse.queueFunctionExecution called with a nil Element")
		return
	}
	if el.GetUniverseOfDiscourse(hl) == nil {
		// Functions do not get executed on elements that are no longer in a Universe of Discourse
		return
	}
	if notification.GetNatureOfChange() == 0 {
		log.Printf("universeOfDiscourse.queueFunctionExecution called without of NatureOfChange")
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
func (uOfDPtr *universeOfDiscourse) Redo(hl *HeldLocks) {
	uOfDPtr.undoManager.redo(uOfDPtr, hl)
}

func (uOfDPtr *universeOfDiscourse) removeElementForUndo(el Element, hl *HeldLocks) {
	if el != nil {
		hl.ReadLockElement(el)
		if uOfDPtr.undoManager.debugUndo == true {
			log.Printf("Removing element for undo, id: %s\n", el.GetConceptID(hl))
			Print(el, "Removed Element: ", hl)
		}
		uOfDPtr.uuidElementMap.DeleteEntry(el.GetConceptID(hl))
	}
}

// RecoverConceptSpace reconstructs a concept space from its JSON representation
func (uOfDPtr *universeOfDiscourse) RecoverConceptSpace(data []byte, hl *HeldLocks) (Element, error) {
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
func (uOfDPtr *universeOfDiscourse) RecoverElement(data []byte, hl *HeldLocks) (Element, error) {
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

func (uOfDPtr *universeOfDiscourse) removeUnresolvedPointer(pointer *cachedPointer) {
	uOfDPtr.unresolvedPointersMap.removeCachedPointer(pointer)
}

// replicateAsRefinement replicates the structure of the original in the replicate, ignoring
// Refinements The name from each original element is copied into the name of the
// corresponding replicate element. This function is idempotent: if applied to an existing structure,
// Elements of that structure that have existing Refinement relationships with original Elements
// will not be re-created.
func (uOfDPtr *universeOfDiscourse) replicateAsRefinement(original Element, replicate Element, hl *HeldLocks) {
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

	for _, originalChild := range original.GetOwnedConcepts(hl) {
		switch originalChild.(type) {
		case Refinement:
			continue
		}
		var replicateChild Element
		// For each original child, determine whether there is already a replicate child that
		// has the original child as one of its abstractions. This is replicateChild
		for _, currentChild := range replicate.GetOwnedConcepts(hl) {
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
func (uOfDPtr *universeOfDiscourse) SetRecordingUndo(newSetting bool) {
	uOfDPtr.undoManager.setRecordingUndo(newSetting)
}

// SetUniverseOfDiscourse() sets the uOfD of which this element is a member. Strictly
// speaking, this is not an attribute of the elment, but rather a context in which
// the element is operating in which the element may be able to locate other objects
// by id.
func (uOfDPtr *universeOfDiscourse) SetUniverseOfDiscourse(el Element, hl *HeldLocks) error {
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
func (uOfDPtr *universeOfDiscourse) Undo(hl *HeldLocks) {
	uOfDPtr.undoManager.undo(uOfDPtr, hl)
}

// uncacheUnresolvedPointers is a support function used when removing an Element from a uOfD.
// It removes all of the element's unresolved pointers from the uOfD cache
func (uOfDPtr *universeOfDiscourse) uncacheUnresolvedPointers(el Element, hl *HeldLocks) {
	if el.GetOwningConcept(hl) == nil && el.GetOwningConceptID(hl) != "" {
		uOfDPtr.removeUnresolvedPointer(el.getOwningConceptPointer())
	}
	switch el.(type) {
	case *reference:
		ref := el.(*reference)
		if ref.GetReferencedConcept(hl) == nil && ref.GetReferencedConceptID(hl) != "" {
			uOfDPtr.removeUnresolvedPointer(ref.referencedConcept)
		}
	case *refinement:
		ref := el.(*refinement)
		if ref.GetAbstractConcept(hl) == nil && ref.GetAbstractConceptID(hl) != "" {
			uOfDPtr.removeUnresolvedPointer(ref.abstractConcept)
		}
		if ref.GetRefinedConcept(hl) == nil && ref.GetRefinedConceptID(hl) != "" {
			uOfDPtr.removeUnresolvedPointer(ref.refinedConcept)
		}
	}
}

func (uOfDPtr *universeOfDiscourse) unmarshalPolymorphicElement(data []byte, result *Element, hl *HeldLocks) error {
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

// unresolveIncomingPointers is a support function used when removing an Element fron the uOfD.
// It finds all incoming cachedPointers (owningConcept, referencedConcept, abstractConcept, refinedConcept),
// sets the cachedPointer's indicatedConcept value to nil, and adds the cachedPointer to the uOfD's unresolved
// pointers
func (uOfDPtr *universeOfDiscourse) unresolveIncomingPointers(el Element, hl *HeldLocks) {
	for _, ownedConcept := range el.GetOwnedConcepts(hl) {
		// verify that it's pointing to this Element
		cp := ownedConcept.getOwningConceptPointer()
		if cp.getIndicatedConcept() == el {
			cp.setIndicatedConcept(nil)
			uOfDPtr.addUnresolvedPointer(cp)
		}
	}
	for _, listener := range el.getListeners(hl) {
		switch listener.(type) {
		case *reference:
			ref := listener.(*reference)
			cp := ref.referencedConcept
			if cp.getIndicatedConcept() == el {
				cp.setIndicatedConcept(nil)
				uOfDPtr.addUnresolvedPointer(cp)
			}
		case *refinement:
			ref := listener.(*refinement)
			acCp := ref.abstractConcept
			if acCp.getIndicatedConcept() == el {
				acCp.setIndicatedConcept(nil)
				uOfDPtr.addUnresolvedPointer(acCp)
			}
			rcCP := ref.refinedConcept
			if rcCP.getIndicatedConcept() == el {
				rcCP.setIndicatedConcept(nil)
				uOfDPtr.addUnresolvedPointer(rcCP)
			}
		}
	}
}

func (uOfDPtr *universeOfDiscourse) uriValidForConceptID(uri ...string) error {
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

// UniverseOfDiscourse defines the set of concepts (Elements) currently in scope
type UniverseOfDiscourse interface {
	Element
	addElement(el Element, inRecovery bool, hl *HeldLocks) error
	AddFunction(string, crlExecutionFunction)
	addUnresolvedPointer(*cachedPointer)
	changeURIForElement(Element, string, string) error
	CreateReplicateAsRefinement(Element, *HeldLocks, ...string) Element
	CreateReplicateAsRefinementFromURI(string, *HeldLocks, ...string) (Element, error)
	deleteElement(Element, map[string]Element, *HeldLocks) error
	DeleteElement(Element, *HeldLocks) error
	DeleteElements(map[string]Element, *HeldLocks) error
	findFunctions(Element, *ChangeNotification, *HeldLocks) []string
	generateConceptID(...string) (string, error)
	getComputeFunctions() *functions
	GetElement(string) Element
	GetElements() map[string]Element
	GetElementWithURI(string) Element
	getExecutedCalls() chan *pendingFunctionCall
	GetFunctions(string) []crlExecutionFunction
	GetLiteral(string) Literal
	GetLiteralWithURI(string) Literal
	GetReference(string) Reference
	GetReferenceWithURI(string) Reference
	GetRefinement(string) Refinement
	GetRefinementWithURI(string) Refinement
	GetRootElements(*HeldLocks) map[string]Element
	getURIElementMap() *StringElementMap
	IsRecordingUndo() bool
	MarkUndoPoint()
	MarshalConceptSpace(Element, *HeldLocks) ([]byte, error)
	NewConceptChangeNotification(Element, *HeldLocks) *ChangeNotification
	NewForwardingChangeNotification(Element, NatureOfChange, *ChangeNotification) *ChangeNotification
	NewElement(*HeldLocks, ...string) (Element, error)
	NewHeldLocks() *HeldLocks
	NewReference(*HeldLocks, ...string) (Reference, error)
	NewLiteral(*HeldLocks, ...string) (Literal, error)
	NewRefinement(*HeldLocks, ...string) (Refinement, error)
	NewUniverseOfDiscourseChangeNotification(*ChangeNotification) *ChangeNotification
	preChange(Element, *HeldLocks)
	queueFunctionExecutions(Element, *ChangeNotification, *HeldLocks)
	RecoverElement([]byte, *HeldLocks) (Element, error)
	RecoverConceptSpace([]byte, *HeldLocks) (Element, error)
	SetUniverseOfDiscourse(Element, *HeldLocks) error
	// RecoverElement([]byte) Element
	Redo(*HeldLocks)
	// SetDebugUndo(bool)
	SetRecordingUndo(bool)
	// SetUniverseOfDiscourseRecursively(BaseElement, *HeldLocks)
	Undo(*HeldLocks)
	// uOfDChanged(*ChangeNotification, *HeldLocks)
	// updateUriIndices(BaseElement, *HeldLocks)
}
