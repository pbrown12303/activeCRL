package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/satori/go.uuid"
)

type baseElementPointer struct {
	pointer
	baseElement        BaseElement
	baseElementId      uuid.UUID
	baseElementVersion int
}

func (bepPtr *baseElementPointer) clone() *baseElementPointer {
	var bep baseElementPointer
	bep.cloneAttributes(*bepPtr)
	return &bep
}

func (bepPtr *baseElementPointer) cloneAttributes(source baseElementPointer) {
	bepPtr.pointer.cloneAttributes(source.pointer)
	bepPtr.baseElement = source.baseElement
	bepPtr.baseElementId = source.baseElementId
	bepPtr.baseElementVersion = source.baseElementVersion
}

func (bepPtr *baseElementPointer) baseElementChanged(notification *ChangeNotification) {
	// Circular references need to be detected and curtailed, hence the isReferenced() call
	if bepPtr.GetOwningElement() != nil && notification.isReferenced(bepPtr) == false {
		newNotification := NewChangeNotification(bepPtr, MODIFY, notification)
		bepPtr.GetOwningElement().childChanged(newNotification)
	}

}

func (bepPtr *baseElementPointer) GetBaseElement() BaseElement {
	bepPtr.TraceableLock()
	defer bepPtr.TraceableUnlock()
	return bepPtr.getBaseElement()
}

// getElement() assumes that all relevant locking is being managed elsewhere
func (bepPtr *baseElementPointer) getBaseElement() BaseElement {
	if bepPtr.baseElement == nil && bepPtr.getBaseElementIdentifier() != uuid.Nil && bepPtr.uOfD != nil {
		bepPtr.baseElement = bepPtr.uOfD.getBaseElement(bepPtr.getBaseElementIdentifier().String())
	}
	return bepPtr.baseElement
}

// GetNameNoLock() is a non-locking function that returns the name
func (bepPtr *baseElementPointer) GetNameNoLock() string {
	return "baseElementPointer"
}

// GetBaseElementIdentifier() locks the vase element pointer and returns the base element identifier, releasing the lock in the process
func (bepPtr *baseElementPointer) GetBaseElementIdentifier() uuid.UUID {
	bepPtr.TraceableLock()
	defer bepPtr.TraceableUnlock()
	return bepPtr.getBaseElementIdentifier()
}

// getBaseElementIdentifier() returns the base element identifier without locking
func (bepPtr *baseElementPointer) getBaseElementIdentifier() uuid.UUID {
	return bepPtr.baseElementId
}

func (bepPtr *baseElementPointer) GetBaseElementVersion() int {
	return bepPtr.baseElementVersion
}

func (bepPtr *baseElementPointer) initializeBaseElementPointer() {
	bepPtr.initializePointer()
}

func (bePtr *baseElementPointer) isEquivalent(be *baseElementPointer) bool {
	if bePtr.baseElementId != be.baseElementId {
		fmt.Printf("Equivalence failed: indicated base element ids do not match \n")
		return false
	}
	if bePtr.baseElementVersion != be.baseElementVersion {
		fmt.Printf("Equivalence failed: indicated base element versions do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer)
}

func (elPtr *baseElementPointer) MarshalJSON() ([]byte, error) {
	elPtr.TraceableLock()
	defer elPtr.TraceableUnlock()
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalBaseElementPointerFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *baseElementPointer) marshalBaseElementPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.pointer.marshalPointerFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"BaseElementId\":\"%s\",", elPtr.baseElementId.String()))
	buffer.WriteString(fmt.Sprintf("\"BaseElementVersion\":\"%d\"", elPtr.baseElementVersion))
	return err
}

func (bepPtr *baseElementPointer) printBaseElementPointer(prefix string) {
	bepPtr.printPointer(prefix)
	log.Printf("%sIndicated BaseElementID: %s \n", prefix, bepPtr.baseElementId.String())
	log.Printf("%sIndicated BaseElementVersion: %d \n", prefix, bepPtr.baseElementVersion)
}

func (ep *baseElementPointer) recoverBaseElementPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	err := ep.pointer.recoverPointerFields(unmarshaledData)
	if err != nil {
		fmt.Printf("BaseElementPointer's Recovery of PointerFields failed\n")
		return err
	}
	// Element ID
	var recoveredElementId string
	err = json.Unmarshal((*unmarshaledData)["BaseElementId"], &recoveredElementId)
	if err != nil {
		fmt.Printf("BaseElementPointer's Recovery of BaseElementId failed\n")
		return err
	}
	ep.baseElementId, err = uuid.FromString(recoveredElementId)
	if err != nil {
		fmt.Printf("BaseElementPointer's conversion of BaseElementId failed\n")
		return err
	}
	// Version
	var recoveredElementVersion string
	err = json.Unmarshal((*unmarshaledData)["BaseElementVersion"], &recoveredElementVersion)
	if err != nil {
		fmt.Printf("BaseElementPointer's Recovery of BaseElementVersion failed\n")
		return err
	}
	ep.baseElementVersion, err = strconv.Atoi(recoveredElementVersion)
	if err != nil {
		fmt.Printf("Conversion of BaseElementPointer.elementVersion failed\n")
		return err
	}
	return nil
}

// SetBaseElement() establishes the element to which this pointer points. If this pointer
// happens to be an OWNING_ELEMENT pointer, there is a side-effect in which this pointer's
// owner is removed as a child from the old target element and added as a child to the new
// target element. Locking must take this into account.
func (bepPtr *baseElementPointer) SetBaseElement(baseElement BaseElement) {
	bepPtr.TraceableLock()
	defer bepPtr.TraceableUnlock()
	oldBaseElement := bepPtr.getBaseElement()
	if oldBaseElement == nil && baseElement == nil {
		return // Nothing to do
	} else if oldBaseElement != nil && baseElement != nil && oldBaseElement.getId() == baseElement.getId() {
		return // Nothing to do
	}
	if baseElement != nil {
		baseElement.TraceableLock() // We need to lock the element to make sure it's version doesn't change during this operation
		defer baseElement.TraceableUnlock()
	}
	bepPtr.setBaseElement(baseElement)
}

// setBaseElement() is intended for internal use within the core. It assumes that all relevant
// objects (parent, child, the element pointer itself) have already been locked. All of the
// operations it invokes are also non-locking
func (bepPtr *baseElementPointer) setBaseElement(baseElement BaseElement) {
	if baseElement != bepPtr.baseElement {
		oldPtr := bepPtr.baseElement
		preChange(bepPtr)
		if oldPtr != nil {
			bepPtr.uOfD.removeBaseElementListener(oldPtr, bepPtr)
		}
		bepPtr.baseElement = baseElement
		if baseElement != nil {
			bepPtr.baseElementId = baseElement.getId()
			bepPtr.baseElementVersion = baseElement.getVersion()
			bepPtr.uOfD.addBaseElementListener(baseElement, bepPtr)
		} else {
			bepPtr.baseElementId = uuid.Nil
			bepPtr.baseElementVersion = 0
		}
		notification := NewChangeNotification(bepPtr, MODIFY, nil)
		postChange(bepPtr, notification)
	}
}

// SetOwningElement() actually manages relationships between a number of objects,
// particularly when the pointer is the OWNING_ELEMENT pointer for its owner.
// Because of the complex wiring between the objects, we have to lock all relevant
// objects here and then use non-locking worker methods
func (bepPtr *baseElementPointer) SetOwningElement(newOwningElement Element) {
	bepPtr.TraceableLock()
	defer bepPtr.TraceableUnlock()
	oldOwningElement := bepPtr.getOwningElement()
	if oldOwningElement == nil && newOwningElement == nil {
		return // Nothing to do
	} else if oldOwningElement != nil && newOwningElement != nil && oldOwningElement.getId() == newOwningElement.getId() {
		return // Nothing to do
	}
	if oldOwningElement != nil {
		oldOwningElement.TraceableLock()
		defer oldOwningElement.TraceableUnlock()
	}
	if newOwningElement != nil {
		newOwningElement.TraceableLock()
		defer newOwningElement.TraceableUnlock()
	}
	bepPtr.SetOwningElementNoLock(newOwningElement)
}

// SetOwningElementNoLock() is a non-locking function that sets the ownership of the element pointer.
// It adjusts the ownedBaseElement set of both the old and new owner. In addition, if it is an
// owningElementPointer, it adjusts the ownedBaseElement set of the owner's owner.
func (bepPtr *baseElementPointer) SetOwningElementNoLock(element Element) {
	if element != bepPtr.getOwningElement() {

		if bepPtr.getOwningElement() != nil {
			bepPtr.getOwningElement().removeOwnedBaseElement(bepPtr)
		}

		preChange(bepPtr)
		bepPtr.owningElement = element
		notification := NewChangeNotification(bepPtr, MODIFY, nil)
		postChange(bepPtr, notification)

		if bepPtr.getOwningElement() != nil {
			bepPtr.getOwningElement().addOwnedBaseElement(bepPtr)
		}

	}
}

// internalSetOwningElement() is an internal function used only when unmarshaling.
func (bepPtr *baseElementPointer) internalSetOwningElement(element Element) {
	if element != bepPtr.getOwningElement() {
		bepPtr.owningElement = element
		if bepPtr.getOwningElement() != nil {
			bepPtr.getOwningElement().internalAddOwnedBaseElement(bepPtr)
		}
	}
}

func (bepPtr *baseElementPointer) SetUri(uri string) {
	bepPtr.TraceableLock()
	defer bepPtr.TraceableUnlock()
	bepPtr.SetUriNoLock(uri)
}

func (bepPtr *baseElementPointer) SetUriNoLock(uri string) {
	preChange(bepPtr)
	bepPtr.uri = uri
	notification := NewChangeNotification(bepPtr, MODIFY, nil)
	postChange(bepPtr, notification)
}

type BaseElementPointer interface {
	Pointer
	baseElementChanged(*ChangeNotification)
	getBaseElement() BaseElement
	GetBaseElement() BaseElement
	GetBaseElementIdentifier() uuid.UUID
	GetBaseElementVersion() int
	setBaseElement(BaseElement)
	SetBaseElement(BaseElement)
}
