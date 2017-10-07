package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/satori/go.uuid"
)

type baseElement struct {
	sync.Mutex
	id      uuid.UUID
	version int
	uOfD    *UniverseOfDiscourse
}

func (bePtr *baseElement) cloneAttributes(source baseElement) {
	bePtr.id = source.id
	bePtr.version = source.version
	bePtr.uOfD = source.uOfD
}

func (bePtr *baseElement) GetId(hl *HeldLocks) uuid.UUID {
	if hl != nil {
		hl.LockBaseElement(bePtr)
	}
	return bePtr.id
}

func (bePtr *baseElement) getIdNoLock() uuid.UUID {
	return bePtr.id
}

func (bePtr *baseElement) GetUniverseOfDiscourse(hl *HeldLocks) *UniverseOfDiscourse {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	return bePtr.uOfD
}

func (bePtr *baseElement) GetVersion(hl *HeldLocks) int {
	if hl != nil {
		hl.LockBaseElement(bePtr)
	}
	return bePtr.version
}

// initializeBaseElement() initializes the uuid. Note that initialization does
// not increment the version counter nor does it notify other objects that a change
// has occurred.
func (bePtr *baseElement) initializeBaseElement() {
	bePtr.id = uuid.NewV4()
}

// internalIncrementVersion() increments the version counter. Note that it does not
// notify other objects that a change has occurred.
func (bePtr *baseElement) internalIncrementVersion() {
	bePtr.version++
}

func (bePtr *baseElement) isEquivalent(be *baseElement, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	hl.LockBaseElement(be)
	if bePtr.id != be.id {
		log.Printf("Equivalence failed: ids do not match \n")
		return false
	}
	if bePtr.version != be.version {
		log.Printf("Equivalence failed: versions do not match \n")
		return false
	}
	return true
}

func (elPtr *baseElement) marshalBaseElementFields(buffer *bytes.Buffer) error {
	buffer.WriteString(fmt.Sprintf("\"Id\":\"%s\",", elPtr.id.String()))
	buffer.WriteString(fmt.Sprintf("\"Version\":\"%d\",", elPtr.version))
	return nil
}

func (bePtr *baseElement) printBaseElement(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	log.Printf("%s  id: %s \n", prefix, bePtr.id.String())
	log.Printf("%s  version %d \n", prefix, bePtr.version)
}

func (be *baseElement) recoverBaseElementFields(unmarshaledData *map[string]json.RawMessage) error {
	// Id
	var recoveredId string
	err := json.Unmarshal((*unmarshaledData)["Id"], &recoveredId)
	if err != nil {
		log.Printf("Recovery of BaseElement.id as string failed\n")
		return err
	}
	be.id, err = uuid.FromString(recoveredId)
	if err != nil {
		log.Printf("Conversion of string to uuid failed\n")
		return err
	}
	// Version
	var recoveredVersion string
	err = json.Unmarshal((*unmarshaledData)["Version"], &recoveredVersion)
	if err != nil {
		log.Printf("Recovery of BaseElement.version failed\n")
		return err
	}
	be.version, err = strconv.Atoi(recoveredVersion)
	if err != nil {
		log.Printf("Conversion of BaseElement.version failed\n")
		return err
	}
	return nil
}

// setUniverseOfDiscourse() sets the uOfD of which this object is a member. Strictly
// speaking, this is not an attribute of the baseElement, but rather a context in which
// the baseElement is operating in which the baseElement may be able to locate other objects
// by id.
func (bePtr *baseElement) setUniverseOfDiscourse(uOfD *UniverseOfDiscourse, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	bePtr.uOfD = uOfD
}

func (bePtr *baseElement) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock Base Element %p\n", bePtr)
	}
	bePtr.Lock()
}

func (bePtr *baseElement) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock Base Element %p\n", bePtr)
	}
	bePtr.Unlock()
}

type BaseElement interface {
	GetId(*HeldLocks) uuid.UUID
	getIdNoLock() uuid.UUID
	GetUniverseOfDiscourse(*HeldLocks) *UniverseOfDiscourse
	GetVersion(*HeldLocks) int
	internalIncrementVersion()
	setUniverseOfDiscourse(*UniverseOfDiscourse, *HeldLocks)
	TraceableLock()
	TraceableUnlock()
}

func GetName(be BaseElement, hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be)
	switch be.(type) {
	case Value:
		val := be.(Value)
		return val.getName(hl)
	case Element:
		el := be.(Element)
		nl := el.GetNameLiteral(hl)
		if nl != nil {
			return nl.GetLiteralValue(hl)
		}
	}
	return ""
}

func GetOwningElement(be BaseElement, hl *HeldLocks) Element {
	if be == nil {
		return nil
	}
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be)
	switch be.(type) {
	case Value:
		return be.(Value).getOwningElement(hl)
	case Element:
		oep := be.(Element).GetOwningElementPointer(hl)
		if oep != nil {
			return oep.GetElement(hl)
		}
	}
	return nil
}

func GetUri(be BaseElement, hl *HeldLocks) string {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be)
	switch be.(type) {
	case Value:
		val := be.(Value)
		return val.getUri(hl)
	case Element:
		el := be.(Element)
		nl := el.GetUriLiteral(hl)
		if nl != nil {
			return nl.GetLiteralValue(hl)
		}

	}
	return ""
}

func SetOwningElement(be BaseElement, parent Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be)
	switch be.(type) {
	case Element:
		elPtr := be.(Element)
		oldParent := GetOwningElement(elPtr, hl)
		if oldParent == nil && parent == nil {
			return // Nothing to do
		} else if oldParent != nil && parent != nil && oldParent.GetId(hl) != parent.GetId(hl) {
			return // Nothing to do
		}
		oep := elPtr.GetOwningElementPointer(hl)
		if oep == nil {
			oep = elPtr.GetUniverseOfDiscourse(hl).NewOwningElementPointer(hl)
			//			log.Printf("In case Element of SetOwningElement, created OwningElementPointer and about to call SetOwningElement")
			SetOwningElement(oep, elPtr, hl)
		}
		oep.SetElement(parent, hl)
	case Value:
		//		log.Printf("In case Value of SetOwningElement")
		//		Print(be, "", hl)
		//		Print(parent, "", hl)
		val := be.(Value)
		oldOwningElement := val.getOwningElement(hl)
		if oldOwningElement == nil && parent == nil {
			return // Nothing to do
		} else if oldOwningElement != nil && parent != nil && oldOwningElement.GetId(hl) == parent.GetId(hl) {
			return // Nothing to do
		}
		if val.getOwningElement(hl) != nil {
			removeOwnedBaseElement(val.getOwningElement(hl), val, hl)
		}
		preChange(val, hl)
		val.setOwningElement(parent, hl)
		notification := NewChangeNotification(val, MODIFY, nil)
		postChange(val, notification, hl)

		if val.getOwningElement(hl) != nil {
			//			log.Printf("In case Value of SetOwningElement about to call addOwnedBaseElement")
			addOwnedBaseElement(val.getOwningElement(hl), val, hl)
		}

	}
}

func SetUri(be BaseElement, uri string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be)
	switch be.(type) {
	case Element:
		el := be.(Element)
		nl := el.GetUriLiteral(hl)
		if nl != nil {
			hl.LockBaseElement(nl)
		} else {
			nlp := el.GetUriLiteralPointer(hl)
			if nlp == nil {
				nlp = be.GetUniverseOfDiscourse(hl).NewUriLiteralPointer(hl)
				SetOwningElement(nlp, el, hl)
			}
			nl = be.GetUniverseOfDiscourse(hl).NewLiteral(hl)
			SetOwningElement(nl, el, hl)
			nlp.SetLiteral(nl, hl)
		}
		nl.SetLiteralValue(uri, hl)
	case Value:
		preChange(be, hl)
		be.(Value).setUri(uri, hl)
		notification := NewChangeNotification(be, MODIFY, nil)
		postChange(be, notification, hl)

	}
}
