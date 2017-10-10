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

	"github.com/satori/go.uuid"
)

type elementPointerPointer struct {
	pointer
	elementPointer        ElementPointer
	elementPointerId      uuid.UUID
	elementPointerVersion int
}

func (eppPtr *elementPointerPointer) clone() *elementPointerPointer {
	var clone elementPointerPointer
	clone.cloneAttributes(*eppPtr)
	return &clone
}

func (eppPtr *elementPointerPointer) cloneAttributes(source elementPointerPointer) {
	eppPtr.pointer.cloneAttributes(source.pointer)
	eppPtr.elementPointer = source.elementPointer
	eppPtr.elementPointerId = source.elementPointerId
	eppPtr.elementPointerVersion = source.elementPointerVersion
}

func (eppPtr *elementPointerPointer) GetElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eppPtr)
	if eppPtr.elementPointer == nil && eppPtr.GetElementPointerId(hl) != uuid.Nil && eppPtr.uOfD != nil {
		eppPtr.elementPointer = eppPtr.uOfD.GetElementPointer(eppPtr.GetElementPointerId(hl).String())
	}
	return eppPtr.elementPointer
}

func (eppPtr *elementPointerPointer) getName(hl *HeldLocks) string {
	// No locking required - it's a constant
	return "elementPointerPointer"
}

func (eppPtr *elementPointerPointer) GetElementPointerId(hl *HeldLocks) uuid.UUID {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eppPtr)
	return eppPtr.elementPointerId
}

func (eppPtr *elementPointerPointer) GetElementPointerVersion(hl *HeldLocks) int {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eppPtr)
	return eppPtr.elementPointerVersion
}

func (eppPtr *elementPointerPointer) initializeElementPointerPointer(uri ...string) {
	eppPtr.initializePointer(uri...)
}

func (bePtr *elementPointerPointer) isEquivalent(be *elementPointerPointer, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	if bePtr.elementPointerId != be.elementPointerId {
		fmt.Printf("Equivalence failed: indicated elementPointerPointer ids do not match \n")
		return false
	}
	if bePtr.elementPointerVersion != be.elementPointerVersion {
		fmt.Printf("Equivalence failed: indicated elementPointerPointer versions do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer, hl)
}

func (elPtr *elementPointerPointer) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.maarshalElementPointerPointerFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *elementPointerPointer) maarshalElementPointerPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.pointer.marshalPointerFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"ElementPointerId\":\"%s\",", elPtr.elementPointerId.String()))
	buffer.WriteString(fmt.Sprintf("\"ElementPointerVersion\":\"%d\"", elPtr.elementPointerVersion))
	return err
}

func (eppPtr *elementPointerPointer) printElementPointerPointer(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eppPtr)
	eppPtr.printPointer(prefix, hl)
	log.Printf("%s  Indicated ElementPointerID: %s \n", prefix, eppPtr.elementPointerId.String())
	log.Printf("%s  Indicated ElementPointerVersion: %d \n", prefix, eppPtr.elementPointerVersion)
}

func (ep *elementPointerPointer) recoverElementPointerPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	err := ep.pointer.recoverPointerFields(unmarshaledData)
	if err != nil {
		fmt.Printf("ElementPointerPointer's Recovery of PointerFields failed\n")
		return err
	}
	// ElementPointer ID
	var recoveredElementPointerId string
	err = json.Unmarshal((*unmarshaledData)["ElementPointerId"], &recoveredElementPointerId)
	if err != nil {
		fmt.Printf("ElementPointerPointer's Recovery of ElementPointerId failed\n")
		return err
	}
	ep.elementPointerId, err = uuid.FromString(recoveredElementPointerId)
	if err != nil {
		fmt.Printf("ElementPointerPointer's conversion of ElementPointerId failed\n")
		return err
	}
	// Version
	var recoveredElementPointerVersion string
	err = json.Unmarshal((*unmarshaledData)["ElementPointerVersion"], &recoveredElementPointerVersion)
	if err != nil {
		fmt.Printf("ElementPointerPointer's Recovery of ElementPointerVersion failed\n")
		return err
	}
	ep.elementPointerVersion, err = strconv.Atoi(recoveredElementPointerVersion)
	if err != nil {
		fmt.Printf("Conversion of ElementPointerPointer.elementPointerVersion failed\n")
		return err
	}
	return nil
}

func (eppPtr *elementPointerPointer) SetElementPointer(elementPointer ElementPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eppPtr)
	if elementPointer != eppPtr.elementPointer {
		preChange(eppPtr, hl)
		if eppPtr.elementPointer != nil {
			eppPtr.uOfD.removeElementPointerListener(eppPtr.elementPointer, eppPtr, hl)
		}
		eppPtr.elementPointer = elementPointer
		if elementPointer != nil {
			eppPtr.elementPointerId = elementPointer.GetId(hl)
			eppPtr.elementPointerVersion = elementPointer.GetVersion(hl)
			eppPtr.uOfD.addElementPointerListener(elementPointer, eppPtr, hl)
		} else {
			eppPtr.elementPointerId = uuid.Nil
			eppPtr.elementPointerVersion = 0
		}
		notification := NewChangeNotification(eppPtr, MODIFY, nil)
		postChange(eppPtr, notification, hl)
	}
}

func (eppPtr *elementPointerPointer) SetOwningElement(element Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eppPtr)
	oldOwningElement := eppPtr.getOwningElement(hl)
	if element != oldOwningElement {
		if oldOwningElement != nil {
			removeOwnedBaseElement(oldOwningElement, eppPtr, hl)
		}

		preChange(eppPtr, hl)
		eppPtr.owningElement = element
		notification := NewChangeNotification(eppPtr, MODIFY, nil)
		postChange(eppPtr, notification, hl)

		if element != nil {
			addOwnedBaseElement(element, eppPtr, hl)
		}
	}
}

// internalSetOwningElement() is an internal function used only in unmarshal
func (eppPtr *elementPointerPointer) internalSetOwningElement(element Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(eppPtr)
	if element != eppPtr.getOwningElement(hl) {
		eppPtr.owningElement = element
		if eppPtr.getOwningElement(hl) != nil {
			eppPtr.getOwningElement(hl).internalAddOwnedBaseElement(eppPtr, hl)
		}
	}
}

func (epPtr *elementPointerPointer) setUri(uri string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(epPtr)
	preChange(epPtr, hl)
	epPtr.uri = uri
	notification := NewChangeNotification(epPtr, MODIFY, nil)
	postChange(epPtr, notification, hl)
}

type ElementPointerPointer interface {
	Pointer
	GetElementPointer(*HeldLocks) ElementPointer
	GetElementPointerId(*HeldLocks) uuid.UUID
	GetElementPointerVersion(*HeldLocks) int
	SetElementPointer(ElementPointer, *HeldLocks)
}
