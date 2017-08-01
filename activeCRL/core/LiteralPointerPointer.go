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

type literalPointerPointer struct {
	pointer
	literalPointer        LiteralPointer
	literalPointerId      uuid.UUID
	literalPointerVersion int
}

func (pllPtr *literalPointerPointer) clone() *literalPointerPointer {
	var clone literalPointerPointer
	clone.cloneAttributes(*pllPtr)
	return &clone
}

func (pllPtr *literalPointerPointer) cloneAttributes(source literalPointerPointer) {
	pllPtr.pointer.cloneAttributes(source.pointer)
	pllPtr.literalPointer = source.literalPointer
	pllPtr.literalPointerId = source.literalPointerId
	pllPtr.literalPointerVersion = source.literalPointerVersion
}

func (pllPtr *literalPointerPointer) GetLiteralPointer(hl *HeldLocks) LiteralPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(pllPtr)
	if pllPtr.literalPointer == nil && pllPtr.GetLiteralPointerIdentifier(hl) != uuid.Nil && pllPtr.uOfD != nil {
		pllPtr.literalPointer = pllPtr.uOfD.GetLiteralPointer(pllPtr.GetLiteralPointerIdentifier(hl).String())
	}
	return pllPtr.literalPointer
}

func (pllPtr *literalPointerPointer) getName(hl *HeldLocks) string {
	// No locking required - it's a constant
	return "literalPointerPointer"
}

func (pllPtr *literalPointerPointer) GetLiteralPointerIdentifier(hl *HeldLocks) uuid.UUID {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(pllPtr)
	return pllPtr.literalPointerId
}

func (pllPtr *literalPointerPointer) GetLiteralPointerVersion(hl *HeldLocks) int {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(pllPtr)
	return pllPtr.literalPointerVersion
}

func (pllPtr *literalPointerPointer) initializeLiteralPointerPointer() {
	pllPtr.initializePointer()
}

func (bePtr *literalPointerPointer) isEquivalent(be *literalPointerPointer, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	if bePtr.literalPointerId != be.literalPointerId {
		fmt.Printf("Equivalence failed: indicated literalPointerPointer ids do not match \n")
		return false
	}
	if bePtr.literalPointerVersion != be.literalPointerVersion {
		fmt.Printf("Equivalence failed: indicated literalPointerPointer versions do not match \n")
		return false
	}
	var pointerPtr *pointer = &bePtr.pointer
	return pointerPtr.isEquivalent(&be.pointer, hl)
}

func (elPtr *literalPointerPointer) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.maarshalLiteralPointerPointerFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *literalPointerPointer) maarshalLiteralPointerPointerFields(buffer *bytes.Buffer) error {
	err := elPtr.pointer.marshalPointerFields(buffer)
	buffer.WriteString(fmt.Sprintf("\"LiteralPointerId\":\"%s\",", elPtr.literalPointerId.String()))
	buffer.WriteString(fmt.Sprintf("\"LiteralPointerVersion\":\"%d\"", elPtr.literalPointerVersion))
	return err
}

func (pllPtr *literalPointerPointer) printLiteralPointerPointer(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(pllPtr)
	pllPtr.printPointer(prefix, hl)
	log.Printf("%s  Indicated LiteralPointerID: %s \n", prefix, pllPtr.literalPointerId.String())
	log.Printf("%s  Indicated LiteralPointerVersion: %d \n", prefix, pllPtr.literalPointerVersion)
}

func (ep *literalPointerPointer) recoverLiteralPointerPointerFields(unmarshaledData *map[string]json.RawMessage) error {
	err := ep.pointer.recoverPointerFields(unmarshaledData)
	if err != nil {
		fmt.Printf("LiteralPointerPointer's Recovery of PointerFields failed\n")
		return err
	}
	// LiteralPointer ID
	var recoveredLiteralPointerId string
	err = json.Unmarshal((*unmarshaledData)["LiteralPointerId"], &recoveredLiteralPointerId)
	if err != nil {
		fmt.Printf("LiteralPointerPointer's Recovery of LiteralPointerId failed\n")
		return err
	}
	ep.literalPointerId, err = uuid.FromString(recoveredLiteralPointerId)
	if err != nil {
		fmt.Printf("LiteralPointerPointer's conversion of LiteralPointerId failed\n")
		return err
	}
	// Version
	var recoveredLiteralPointerVersion string
	err = json.Unmarshal((*unmarshaledData)["LiteralPointerVersion"], &recoveredLiteralPointerVersion)
	if err != nil {
		fmt.Printf("LiteralPointerPointer's Recovery of LiteralPointerVersion failed\n")
		return err
	}
	ep.literalPointerVersion, err = strconv.Atoi(recoveredLiteralPointerVersion)
	if err != nil {
		fmt.Printf("Conversion of LiteralPointerPointer.literalPointerVersion failed\n")
		return err
	}
	return nil
}

func (pllPtr *literalPointerPointer) SetLiteralPointer(literalPointer LiteralPointer, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(pllPtr)
	if literalPointer != pllPtr.literalPointer {
		preChange(pllPtr, hl)
		if pllPtr.literalPointer != nil {
			pllPtr.uOfD.removeLiteralPointerListener(pllPtr.literalPointer, pllPtr, hl)
		}
		pllPtr.literalPointer = literalPointer
		if literalPointer != nil {
			pllPtr.literalPointerId = literalPointer.GetId(hl)
			pllPtr.literalPointerVersion = literalPointer.GetVersion(hl)
			pllPtr.uOfD.addLiteralPointerListener(literalPointer, pllPtr, hl)
		} else {
			pllPtr.literalPointerId = uuid.Nil
			pllPtr.literalPointerVersion = 0
		}
		notification := NewChangeNotification(pllPtr, MODIFY, nil)
		postChange(pllPtr, notification, hl)
	}
}

func (pllPtr *literalPointerPointer) SetOwningElement(element Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(pllPtr)
	oldOwner := pllPtr.getOwningElement(hl)
	if element != oldOwner {
		if oldOwner != nil {
			removeOwnedBaseElement(oldOwner, pllPtr, hl)
		}

		preChange(pllPtr, hl)
		pllPtr.owningElement = element
		notification := NewChangeNotification(pllPtr, MODIFY, nil)
		postChange(pllPtr, notification, hl)

		if element != nil {
			addOwnedBaseElement(element, pllPtr, hl)
		}
	}
}

// internalSetOwningElement() is an internal function used only in unmarshal
func (pllPtr *literalPointerPointer) internalSetOwningElement(element Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(pllPtr)
	if element != pllPtr.getOwningElement(hl) {
		pllPtr.owningElement = element
		if pllPtr.getOwningElement(hl) != nil {
			pllPtr.getOwningElement(hl).internalAddOwnedBaseElement(pllPtr, hl)
		}
	}
}

func (lpPtr *literalPointerPointer) setUri(uri string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(lpPtr)
	preChange(lpPtr, hl)
	lpPtr.uri = uri
	notification := NewChangeNotification(lpPtr, MODIFY, nil)
	postChange(lpPtr, notification, hl)
}

type LiteralPointerPointer interface {
	Pointer
	GetLiteralPointer(*HeldLocks) LiteralPointer
	GetLiteralPointerIdentifier(*HeldLocks) uuid.UUID
	GetLiteralPointerVersion(*HeldLocks) int
	SetLiteralPointer(LiteralPointer, *HeldLocks)
}
