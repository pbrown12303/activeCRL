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
	sync.RWMutex
	id      uuid.UUID
	version int
	uOfD    *UniverseOfDiscourse
}

func (bePtr *baseElement) cloneAttributes(source baseElement) {
	bePtr.id = source.id
	bePtr.version = source.version
	bePtr.uOfD = source.uOfD
}

// GetId locks the element, reads the id, and returns, releasing the lock
func (bePtr *baseElement) GetId() uuid.UUID {
	bePtr.TraceableLock()
	defer bePtr.TraceableUnlock()
	return bePtr.getId()
}

// getId returns the id without locking
func (bePtr *baseElement) getId() uuid.UUID {
	return bePtr.id
}

func (bePtr *baseElement) getUniverseOfDiscourse() *UniverseOfDiscourse {
	return bePtr.uOfD
}

// GetVersion() Locks the element and returns the version, releasing the lock
func (bePtr *baseElement) GetVersion() int {
	bePtr.TraceableLock()
	defer bePtr.TraceableUnlock()
	return bePtr.getVersion()
}

// getVersion() returns the version withoug locking
func (bePtr *baseElement) getVersion() int {
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

func (bePtr *baseElement) isEquivalent(be *baseElement) bool {
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

func (bePtr *baseElement) printBaseElement(prefix string) {
	log.Printf("%sid: %s \n", prefix, bePtr.id.String())
	log.Printf("%sversion %d \n", prefix, bePtr.version)
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
func (bePtr *baseElement) setUniverseOfDiscourse(uOfD *UniverseOfDiscourse) {
	bePtr.uOfD = uOfD
}

func (bePtr *baseElement) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock Base Element %p\n", bePtr)
	}
	bePtr.Lock()
}

func (bePtr *baseElement) TraceableRLock() {
	if TraceLocks {
		log.Printf("About to lock Base Element %p\n", bePtr)
	}
	bePtr.RLock()
}

func (bePtr *baseElement) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock Base Element %p\n", bePtr)
	}
	bePtr.Unlock()
}

func (bePtr *baseElement) TraceableRUnlock() {
	if TraceLocks {
		log.Printf("About to unlock Base Element %p\n", bePtr)
	}
	bePtr.RUnlock()
}

type BaseElement interface {
	getId() uuid.UUID
	GetId() uuid.UUID
	GetNameNoLock() string
	getOwningElement() Element
	GetOwningElement() Element
	getUniverseOfDiscourse() *UniverseOfDiscourse
	GetUri() string
	GetUriNoLock() string
	getVersion() int
	GetVersion() int
	internalIncrementVersion()
	SetOwningElement(Element)
	SetOwningElementNoLock(Element)
	setUniverseOfDiscourse(*UniverseOfDiscourse)
	SetUri(string)
	SetUriNoLock(string)
	TraceableLock()
	TraceableUnlock()
}
