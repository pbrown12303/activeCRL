package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"strconv"
)

type baseElement struct {
	id      uuid.UUID
	version int
	uOfD    *UniverseOfDiscourse
}

func (bePtr *baseElement) GetId() uuid.UUID {
	return bePtr.id
}

func (bePtr *baseElement) getUniverseOfDiscourse() *UniverseOfDiscourse {
	return bePtr.uOfD
}

func (bePtr *baseElement) GetVersion() int {
	return bePtr.version
}

func (bePtr *baseElement) initializeBaseElement() {
	bePtr.id = uuid.NewV4()
}

func (bePtr *baseElement) incrementVersion() {
	bePtr.version++
}

func (bePtr *baseElement) isEquivalent(be *baseElement) bool {
	if bePtr.id != be.id {
		fmt.Printf("Equivalence failed: ids do not match \n")
		return false
	}
	if bePtr.version != be.version {
		fmt.Printf("Equivalence failed: versions do not match \n")
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
	fmt.Printf("%sid: %s \n", prefix, bePtr.id.String())
	fmt.Printf("%sversion %d \n", prefix, bePtr.version)
}

func (be *baseElement) recoverBaseElementFields(unmarshaledData *map[string]json.RawMessage) error {
	// Id
	var recoveredId string
	err := json.Unmarshal((*unmarshaledData)["Id"], &recoveredId)
	if err != nil {
		fmt.Printf("Recovery of BaseElement.id as string failed\n")
		return err
	}
	be.id, err = uuid.FromString(recoveredId)
	if err != nil {
		fmt.Printf("Conversion of string to uuid failed\n")
		return err
	}
	// Version
	var recoveredVersion string
	err = json.Unmarshal((*unmarshaledData)["Version"], &recoveredVersion)
	if err != nil {
		fmt.Printf("Recovery of BaseElement.version failed\n")
		return err
	}
	be.version, err = strconv.Atoi(recoveredVersion)
	if err != nil {
		fmt.Printf("Conversion of BaseElement.version failed\n")
		return err
	}
	return nil
}

func (bePtr *baseElement) setUniverseOfDiscourse(uOfD *UniverseOfDiscourse) {
	bePtr.uOfD = uOfD
}

type BaseElement interface {
	GetId() uuid.UUID
	GetName() string
	GetOwningElement() Element
	getUniverseOfDiscourse() *UniverseOfDiscourse
	GetVersion() int
	setOwningElement(Element)
	setUniverseOfDiscourse(*UniverseOfDiscourse)
}
