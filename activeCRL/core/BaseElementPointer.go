// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	//	"github.com/satori/go.uuid"
)

type baseElementPointer struct {
	pointer
	baseEl             BaseElement
	baseElementId      string
	baseElementVersion int
}

func (bepPtr *baseElementPointer) clone() *baseElementPointer {
	var bep baseElementPointer
	bep.cloneAttributes(*bepPtr)
	return &bep
}

func (bepPtr *baseElementPointer) cloneAttributes(source baseElementPointer) {
	bepPtr.pointer.cloneAttributes(source.pointer)
	bepPtr.baseEl = source.baseEl
	bepPtr.baseElementId = source.baseElementId
	bepPtr.baseElementVersion = source.baseElementVersion
}

func (bepPtr *baseElementPointer) baseElementChanged(notification *ChangeNotification, hl *HeldLocks) {
	// Circular references need to be detected and curtailed, hence the isReferenced() call
	if bepPtr.getOwningElement(hl) != nil && notification.isReferenced(bepPtr) == false {
		newNotification := NewChangeNotification(bepPtr, MODIFY, "baseElementChanged", notification)
		childChanged(bepPtr.getOwningElement(hl), newNotification, hl)
	}

}

func (bepPtr *baseElementPointer) GetBaseElement(hl *HeldLocks) BaseElement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(&bepPtr.pointer.baseElement)
	if bepPtr.baseEl == nil && bepPtr.GetBaseElementId(hl) != "" && bepPtr.uOfD != nil {
		bepPtr.baseEl = bepPtr.uOfD.GetBaseElement(bepPtr.GetBaseElementId(hl))
	}
	return bepPtr.baseEl
}

func (bepPtr *baseElementPointer) getLabel(hl *HeldLocks) string {
	// No need to lock - this is a constant
	return bepPtr.getLabelNoLock()
}

func (bepPtr *baseElementPointer) getLabelNoLock() string {
	return "baseElementPointer"
}

// GetBaseElementIdentifier() locks the vase element pointer and returns the base element identifier, releasing the lock in the process
func (bepPtr *baseElementPointer) GetBaseElementId(hl *HeldLocks) string {
	if hl != nil {
		hl.LockBaseElement(&bepPtr.pointer.baseElement)
	}
	return bepPtr.baseElementId
}

func (bepPtr *baseElementPointer) GetBaseElementVersion(hl *HeldLocks) int {
	if hl != nil {
		hl.LockBaseElement(bepPtr)
	}
	return bepPtr.baseElementVersion
}

func (bepPtr *baseElementPointer) initializeBaseElementPointer(uri ...string) {
	bepPtr.initializePointer(uri...)
}

func (bePtr *baseElementPointer) isEquivalent(be *baseElementPointer, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	if bePtr.baseElementId != be.baseElementId {
		fmt.Printf("Equivalence failed: indicated base element ids do not match \n")
		return false
	}
	if bePtr.baseElementVersion != be.baseElementVersion {
		fmt.Printf("Equivalence failed: indicated base element versions do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer, hl)
}

func (elPtr *baseElementPointer) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalBaseElementPointerFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *baseElementPointer) marshalBaseElementPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.pointer.marshalPointerFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"BaseElementId\":\"%s\",", elPtr.baseElementId))
	buffer.WriteString(fmt.Sprintf("\"BaseElementVersion\":\"%d\"", elPtr.baseElementVersion))
	return err
}

func (bepPtr *baseElementPointer) printBaseElementPointer(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bepPtr)
	bepPtr.printPointer(prefix, hl)
	log.Printf("%s  Indicated BaseElementID: %s \n", prefix, bepPtr.baseElementId)
	log.Printf("%s  Indicated BaseElementVersion: %d \n", prefix, bepPtr.baseElementVersion)
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
	ep.baseElementId = recoveredElementId
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
// target element.
func (bepPtr *baseElementPointer) SetBaseElement(newBaseElement BaseElement, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(&bepPtr.baseElement)
	oldBaseElement := bepPtr.baseEl
	if oldBaseElement == nil && newBaseElement == nil {
		return // Nothing to do
	} else if oldBaseElement != nil && newBaseElement != nil && oldBaseElement.GetId(hl) == newBaseElement.GetId(hl) {
		return // Nothing to do
	}
	if newBaseElement != nil {
		hl.LockBaseElement(newBaseElement)
	}
	preChange(bepPtr, hl)
	if oldBaseElement != nil {
		bepPtr.uOfD.removeBaseElementListener(oldBaseElement, bepPtr, hl)
	}
	bepPtr.baseEl = newBaseElement
	if newBaseElement != nil {
		bepPtr.baseElementId = newBaseElement.GetId(hl)
		bepPtr.baseElementVersion = newBaseElement.GetVersion(hl)
		bepPtr.uOfD.addBaseElementListener(newBaseElement, bepPtr, hl)
	} else {
		bepPtr.baseElementId = ""
		bepPtr.baseElementVersion = 0
	}
	notification := NewChangeNotification(bepPtr, MODIFY, "SetBaseElement", nil)
	postChange(bepPtr, notification, hl)
}

func (bepPtr *baseElementPointer) SetOwningElement(newOwningElement Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bepPtr)
	oldOwningElement := bepPtr.getOwningElement(hl)
	if oldOwningElement == nil && newOwningElement == nil {
		return // Nothing to do
	} else if oldOwningElement != nil && newOwningElement != nil && oldOwningElement.GetId(hl) == newOwningElement.GetId(hl) {
		return // Nothing to do
	}
	if bepPtr.getOwningElement(hl) != nil {
		removeOwnedBaseElement(bepPtr.getOwningElement(hl), bepPtr, hl)
	}

	preChange(bepPtr, hl)
	bepPtr.owningElement = newOwningElement
	notification := NewChangeNotification(bepPtr, MODIFY, "SetOwningElement", nil)
	postChange(bepPtr, notification, hl)

	if bepPtr.getOwningElement(hl) != nil {
		addOwnedBaseElement(bepPtr.getOwningElement(hl), bepPtr, hl)
	}
}

// internalSetOwningElement() is an internal function used only when unmarshaling.
func (bepPtr *baseElementPointer) internalSetOwningElement(element Element, hl *HeldLocks) {
	if element != bepPtr.getOwningElement(hl) {
		bepPtr.owningElement = element
		if bepPtr.getOwningElement(hl) != nil {
			bepPtr.getOwningElement(hl).internalAddOwnedBaseElement(bepPtr, hl)
		}
	}
}

func (bepPtr *baseElementPointer) setUri(uri string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bepPtr)
	preChange(bepPtr, hl)
	bepPtr.uri = uri
	notification := NewChangeNotification(bepPtr, MODIFY, "setUri", nil)
	postChange(bepPtr, notification, hl)
}

type BaseElementPointer interface {
	Pointer
	baseElementChanged(*ChangeNotification, *HeldLocks)
	GetBaseElement(*HeldLocks) BaseElement
	GetBaseElementId(*HeldLocks) string
	GetBaseElementVersion(*HeldLocks) int
	SetBaseElement(BaseElement, *HeldLocks)
}
