package core

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"sync"

	"github.com/satori/go.uuid"
)

// universeOfDiscourse represents the scope of relevant concepts
type universeOfDiscourse struct {
	element
	computeFunctions      functions
	executedCalls         chan *pendingFunctionCall
	idURIMap              *StringStringMap
	undoManager           *undoManager
	unresolvedPointersMap *stringCachedPointersMap
	uriElementMap         *StringElementMap
	uuidElementMap        *StringElementMap
	waitGroup             sync.WaitGroup
}

// NewUniverseOfDiscourse creates and initializes a new UniverseOfDiscourse
func NewUniverseOfDiscourse() UniverseOfDiscourse {
	var uOfD universeOfDiscourse
	uOfD.computeFunctions = make(map[string][]crlExecutionFunction)
	uOfD.idURIMap = NewStringStringMap()
	uOfD.undoManager = newUndoManager()
	uOfD.unresolvedPointersMap = newStringCachedPointersMap()
	uOfD.uriElementMap = NewStringElementMap()
	uOfD.uuidElementMap = NewStringElementMap()
	uOfDID, _ := uOfD.generateConceptID(UniverseOfDiscourseURI)
	uOfD.initializeElement(uOfDID)
	hl := uOfD.NewHeldLocks()
	uOfD.SetUniverseOfDiscourse(&uOfD, hl)
	uOfD.AddFunction(coreHousekeepingURI, coreHousekeeping)
	buildCoreConceptSpace(&uOfD, hl)
	uOfD.SetUniverseOfDiscourse(&uOfD, hl)
	hl.ReleaseLocksAndWait()
	return &uOfD
}

func (uOfDPtr *universeOfDiscourse) addElement(el Element, hl *HeldLocks) error {
	if el == nil {
		return errors.New("UniverseOfDiscource addElement() failed because element was nil")
	}
	hl.WriteLockElement(el)
	uuid := el.GetConceptID(hl)
	if uuid == "" {
		return errors.New("UniverseOfDiscource addBaseElement() failed because UUID was nil")
	}
	uOfDPtr.uuidElementMap.SetEntry(el.getConceptIDNoLock(), el)
	uri := el.GetURI(hl)
	if uri != "" {
		uOfDPtr.uriElementMap.SetEntry(uri, el)
		uOfDPtr.idURIMap.SetEntry(el.GetConceptID(hl), uri)
	}
	uOfDPtr.undoManager.markNewElement(el, hl)
	// Make sure it's in the owner's ownedConcepts set
	owner := el.GetOwningConcept(hl)
	if owner != nil {
		owner.addOwnedConcept(uuid, hl)
	}
	// Fix up any unresolved pointers
	uOfDPtr.unresolvedPointersMap.resolveCachedPointers(el, hl)
	// TODO: FIX THIS
	// notification := NewChangeNotification(el, ADD, "AddBaseElement", nil)
	// uOfDPtr.uOfDChanged(notification, hl)
	return nil
}

func (uOfDPtr *universeOfDiscourse) addElementForUndo(el Element, hl *HeldLocks) error {
	if el == nil {
		return errors.New("UniverseOfDiscource addElementForUndo() failed because element was nil")
	}
	if el != nil {
		hl.ReadLockElement(el)
	}
	if uOfDPtr.undoManager.debugUndo == true {
		log.Printf("Adding element for undo, id: %s\n", el.GetConceptID(hl))
		Print(el, "Added Element: ", hl)
	}
	uOfDPtr.uuidElementMap.SetEntry(el.GetConceptID(hl), el)
	return nil
}

func (uOfDPtr *universeOfDiscourse) AddFunction(uri string, function crlExecutionFunction) {
	uOfDPtr.computeFunctions[string(uri)] = append(uOfDPtr.computeFunctions[string(uri)], function)
}

func (uOfDPtr *universeOfDiscourse) addUnresolvedPointer(pointer *cachedPointer) {
	uOfDPtr.unresolvedPointersMap.addCachedPointer(pointer)
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

// ClearUniverseOfDiscourse() removes the element from the uOfD of which this element is a member. Strictly
// speaking, this is not an attribute of the elment, but rather a context in which
// the element is operating in which the element may be able to locate other objects
// by id.
func (uOfDPtr *universeOfDiscourse) ClearUniverseOfDiscourse(el Element, hl *HeldLocks) error {
	if el.GetIsCore() {
		return errors.New("ClearUniverseOfDiscourse called on a CRL core concept")
	}
	hl.WriteLockElement(el)
	if el.GetUniverseOfDiscourse(hl) != uOfDPtr {
		return errors.New("ClearUniverseOfDiscourse called on an Element in a different UofD")
	}
	if el.IsReadOnly(hl) {
		return errors.New("SetUniverseOfDiscourse called on read-only Element")
	}
	uOfDPtr.preChange(el, hl)
	uOfDPtr.removeElement(el, hl)
	el.setUniverseOfDiscourse(nil, hl)
	return nil
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
	element.FindAbstractions(&abstractions, hl)
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

func (uOfDPtr *universeOfDiscourse) generateConceptID(uri ...string) (string, error) {
	var conceptID string
	if len(uri) == 0 {
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

func (uOfDPtr *universeOfDiscourse) getURIElementMap() *StringElementMap {
	return uOfDPtr.uriElementMap
}

func (uOfDPtr *universeOfDiscourse) getWaitGroup() *sync.WaitGroup {
	return &uOfDPtr.waitGroup
}

func (uOfDPtr *universeOfDiscourse) IsRecordingUndo() bool {
	return uOfDPtr.undoManager.recordingUndo
}

// NewConceptChangeNotification creates a ChangeNotification that records state of the concept prior to the
// change. Note that this MUST be called prior to making any changes to the concept.
func (uOfDPtr *universeOfDiscourse) NewConceptChangeNotification(changingConcept Element, hl *HeldLocks) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = changingConcept
	notification.natureOfChange = ConceptChanged
	// Since this function is invoked by the *element methods for Literals, References, and Refinements, we play a
	// game to get the full datatype cloned
	notification.priorConceptState = clone(uOfDPtr.GetElement(changingConcept.getConceptIDNoLock()), hl)
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
	var el element
	el.initializeElement(conceptID)
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
	hl.waitGroup = uOfDPtr.getWaitGroup()
	hl.functionCallManager = newFunctionCallManager(hl.uOfD)
	return &hl
}

func (uOfDPtr *universeOfDiscourse) NewLiteral(hl *HeldLocks, uri ...string) (Literal, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	var lit literal
	lit.initializeLiteral(conceptID)
	hl.WriteLockElement(&lit)
	uOfDPtr.SetUniverseOfDiscourse(&lit, hl)
	return &lit, nil
}

func (uOfDPtr *universeOfDiscourse) NewReference(hl *HeldLocks, uri ...string) (Reference, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	var ref reference
	ref.initializeReference(conceptID)
	hl.WriteLockElement(&ref)
	uOfDPtr.SetUniverseOfDiscourse(&ref, hl)
	return &ref, nil
}

func (uOfDPtr *universeOfDiscourse) NewRefinement(hl *HeldLocks, uri ...string) (Refinement, error) {
	conceptID, err := uOfDPtr.generateConceptID(uri...)
	if err != nil {
		return nil, err
	}
	var ref refinement
	ref.initializeRefinement(conceptID)
	hl.WriteLockElement(&ref)
	uOfDPtr.SetUniverseOfDiscourse(&ref, hl)
	return &ref, nil
}

func (uOfDPtr *universeOfDiscourse) NewUniverseOfDiscourseChangeNotification(underlyingChange *ChangeNotification) *ChangeNotification {
	var notification ChangeNotification
	notification.reportingElement = uOfDPtr
	notification.natureOfChange = UofDChanged
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
	functionIdentifiers := uOfDPtr.findFunctions(el, notification, hl)
	for _, functionIdentifier := range functionIdentifiers {
		if TraceChange == true {
			log.Printf("queueFunctionExecutions calling function, URI: %s", functionIdentifier)
			log.Printf("Function target: %T %s", el, el.getConceptIDNoLock())
			notification.Print("Notification: ", hl)
		}
		hl.functionCallManager.addFunctionCall(functionIdentifier, el, notification)
	}
}

func (uOfDPtr *universeOfDiscourse) removeElement(el Element, hl *HeldLocks) error {
	if el == nil {
		return errors.New("UniverseOfDiscource removeElement failed elcause Element was nil")
	}
	hl.WriteLockElement(el)
	uuid := el.GetConceptID(hl)
	uOfDPtr.uuidElementMap.DeleteEntry(uuid)
	uri := el.GetURI(hl)
	if uri != "" {
		uOfDPtr.uriElementMap.DeleteEntry(uri)
		uOfDPtr.idURIMap.DeleteEntry(el.GetConceptID(hl))
	}
	uOfDPtr.undoManager.markRemovedElement(el, hl)
	// Remove element from owner's child list
	owner := el.GetOwningConcept(hl)
	if owner != nil {
		owner.GetOwningConcept(hl).removeOwnedConcept(uuid, hl)
	}
	// Remove Element from all cached pointers
	for _, listener := range *el.(*element).listeners.CopyMap() {
		switch listener.(type) {
		case *reference:
			ref := listener.(*reference)
			cachedPointer := ref.referencedConcept
			if cachedPointer.getIndicatedConcept() == el {
				cachedPointer.setIndicatedConcept(nil)
				uOfDPtr.unresolvedPointersMap.addCachedPointer(cachedPointer)
			}
		case *refinement:
			ref := listener.(*refinement)
			abstractCachedPointer := ref.abstractConcept
			if abstractCachedPointer.getIndicatedConcept() == el {
				abstractCachedPointer.setIndicatedConcept(nil)
				uOfDPtr.unresolvedPointersMap.addCachedPointer(abstractCachedPointer)
			}
			refinedCachedPointer := ref.refinedConcept
			if refinedCachedPointer.getIndicatedConcept() == el {
				refinedCachedPointer.setIndicatedConcept(nil)
				uOfDPtr.unresolvedPointersMap.addCachedPointer(refinedCachedPointer)
			}
		}
	}
	// TODO: Fix notification for removeElement
	// notification := NewChangeNotification(el, ADD, "removeBaseElement", nil)
	// uOfDPtr.uOfDChanged(notification, hl)
	return nil
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
	// TODO: Fix these
	// uOfD.SetUniverseOfDiscourseRecursively(recoveredElement, hl)
	// restoreValueOwningElementFieldsRecursively(recoveredElement.(Element), hl)
	// uOfD.restoreUriIndexRecursively(recoveredElement, hl)
	uOfDPtr.addElement(recoveredElement, hl)
	return recoveredElement, nil
}

// SetUniverseOfDiscourse() sets the uOfD of which this element is a member. Strictly
// speaking, this is not an attribute of the elment, but rather a context in which
// the element is operating in which the element may be able to locate other objects
// by id.
func (uOfDPtr *universeOfDiscourse) SetUniverseOfDiscourse(el Element, hl *HeldLocks) error {
	hl.WriteLockElement(el)
	currentUofD := el.GetUniverseOfDiscourse(hl)
	if currentUofD != uOfDPtr {
		if el.GetIsCore() {
			return errors.New("SetUniverseOfDiscourse called on a CRL core concept")
		}
		if currentUofD != nil {
			return errors.New("SetUniverseOfDiscourse called on an Element in another uOfD")
		}
		if el.IsReadOnly(hl) {
			return errors.New("SetUniverseOfDiscourse called on read-only Element")
		}
		uOfDPtr.preChange(el, hl)
		el.setUniverseOfDiscourse(uOfDPtr, hl)
		uOfDPtr.addElement(el, hl)
	}
	return nil
}

func (uOfDPtr *universeOfDiscourse) unmarshalPolymorphicElement(data []byte, result *Element, hl *HeldLocks) error {
	// TODO: Fix unmarshall
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
		recoveredElement.initializeElement("")
		*result = &recoveredElement
		err = recoveredElement.recoverElementFields(&unmarshaledData, hl)
		if err != nil {
			return err
		}
	case "*core.reference":
		//		fmt.Printf("Switch choice *core.elementReference \n")
		var recoveredReference reference
		recoveredReference.uOfD = uOfDPtr
		recoveredReference.initializeReference("")
		*result = &recoveredReference
		err = recoveredReference.recoverReferenceFields(&unmarshaledData, hl)
		if err != nil {
			return err
		}
	case "*core.literal":
		//		fmt.Printf("Switch choice *core.literal \n")
		var recoveredLiteral literal
		recoveredLiteral.uOfD = uOfDPtr
		recoveredLiteral.initializeLiteral("")
		*result = &recoveredLiteral
		err = recoveredLiteral.recoverLiteralFields(&unmarshaledData, hl)
		if err != nil {
			return err
		}
	case "*core.refinement":
		var recoveredRefinement refinement
		recoveredRefinement.uOfD = uOfDPtr
		recoveredRefinement.initializeRefinement("")
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
	addElement(Element, *HeldLocks) error
	AddFunction(string, crlExecutionFunction)
	addUnresolvedPointer(*cachedPointer)
	changeURIForElement(Element, string, string) error
	ClearUniverseOfDiscourse(Element, *HeldLocks) error
	findFunctions(Element, *ChangeNotification, *HeldLocks) []string
	getComputeFunctions() *functions
	GetElement(string) Element
	GetElementWithURI(string) Element
	getExecutedCalls() chan *pendingFunctionCall
	GetFunctions(string) []crlExecutionFunction
	GetLiteral(string) Literal
	GetLiteralWithURI(string) Literal
	GetReference(string) Reference
	GetReferenceWithURI(string) Reference
	GetRefinement(string) Refinement
	GetRefinementWithURI(string) Refinement
	getURIElementMap() *StringElementMap
	getWaitGroup() *sync.WaitGroup
	IsRecordingUndo() bool
	// MarkUndoPoint()
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
	removeElement(Element, *HeldLocks) error
	SetUniverseOfDiscourse(Element, *HeldLocks) error
	// RecoverElement([]byte) Element
	// Redo(*HeldLocks)
	// SetDebugUndo(bool)
	// SetRecordingUndo(bool)
	// SetUniverseOfDiscourseRecursively(BaseElement, *HeldLocks)
	// Undo(*HeldLocks)
	// uOfDChanged(*ChangeNotification, *HeldLocks)
	// updateUriIndices(BaseElement, *HeldLocks)
}
