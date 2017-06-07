package core

import (
	"errors"
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

type UniverseOfDiscourse struct {
	sync.Mutex
	baseElementMap map[string]BaseElement
}

func NewUniverseOfDiscourse() *UniverseOfDiscourse {
	var uOfD UniverseOfDiscourse
	uOfD.baseElementMap = make(map[string]BaseElement)
	return &uOfD
}

func (uOfDPtr *UniverseOfDiscourse) AddBaseElement(be BaseElement) error {
	//	log.Printf("Locking UofD\n")
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	return uOfDPtr.addBaseElement(be)
}

func (uOfDPtr *UniverseOfDiscourse) addBaseElement(be BaseElement) error {
	if be == nil {
		return errors.New("UniverseOfDiscource addBaseElement failed because base element was nil")
	}
	//	log.Printf("Locking %T: %s \n", be, be.getId().String())
	//	log.Printf("BaseElement: %+v \n", be)
	be.traceableLock()
	defer be.traceableUnlock()
	//	log.Printf("Got the lock for %T: %s \n", be, be.getId().String())
	if be.getId() == uuid.Nil {
		return errors.New("UniverseOfDiscource addBaseElement failed because UUID was nil")
	}
	oldUOfD := be.getUniverseOfDiscourse()
	if oldUOfD != nil {
		if oldUOfD == uOfDPtr {
			return nil
		} else {
			log.Printf("Locking old UofD\n")
			oldUOfD.traceableLock()
			defer oldUOfD.traceableUnlock()
			oldUOfD.removeBaseElement(be)
		}
	}
	//	log.Printf("Adding be to UofD map")
	uOfDPtr.baseElementMap[be.getId().String()] = be
	//	log.Printf("Setting be's uOfD")
	be.setUniverseOfDiscourse(uOfDPtr)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getBaseElement(id string) BaseElement {
	return uOfDPtr.baseElementMap[id]
}

func (uOfDPtr *UniverseOfDiscourse) GetElement(id string) Element {
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	return uOfDPtr.getElement(id)
}

func (uOfDPtr *UniverseOfDiscourse) getElement(id string) Element {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *element:
		return be.(Element)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getElementPointer(id string) ElementPointer {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *elementPointer:
		return be.(ElementPointer)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getLiteral(id string) Literal {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *literal:
		return be.(Literal)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getLiteralPointer(id string) LiteralPointer {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *literalPointer:
		return be.(LiteralPointer)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getRefinement(id string) Refinement {
	be := uOfDPtr.baseElementMap[id]
	switch be.(type) {
	case *refinement:
		return be.(Refinement)
	}
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) removeBaseElement(be BaseElement) {
	delete(uOfDPtr.baseElementMap, be.getId().String())
}

func (uOfDPtr *UniverseOfDiscourse) SetUniverseOfDiscourseRecursively(be BaseElement) {
	uOfDPtr.traceableLock()
	defer uOfDPtr.traceableUnlock()
	uOfDPtr.setUniverseOfDiscourseRecursively(be)
}

func (uOfDPtr *UniverseOfDiscourse) setUniverseOfDiscourseRecursively(be BaseElement) {
	uOfDPtr.addBaseElement(be)
	switch be.(type) {
	case *element:
		for _, child := range be.(*element).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *elementPointerReference:
		for _, child := range be.(*elementPointerReference).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *elementReference:
		for _, child := range be.(*elementReference).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *literalPointerReference:
		for _, child := range be.(*literalPointerReference).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *literalReference:
		for _, child := range be.(*literalReference).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *refinement:
		for _, child := range be.(*refinement).ownedBaseElements {
			uOfDPtr.setUniverseOfDiscourseRecursively(child)
		}
	case *elementPointer, *elementPointerPointer, *literal, *literalPointer, *literalPointerPointer:
	// Do nothing
	default:
		log.Printf("UniverseOfDiscourse.setUniverseOfDiscourseRecursively is missing case for %T\n", be)
	}
}

func (uOfDPtr *UniverseOfDiscourse) traceableLock() {
	if traceLocks {
		log.Printf("About to lock Universe of Discourse %p\n", uOfDPtr)
	}
	uOfDPtr.Lock()
}

func (uOfDPtr *UniverseOfDiscourse) traceableUnlock() {
	if traceLocks {
		log.Printf("About to unlock Universe of Discourse %p\n", uOfDPtr)
	}
	uOfDPtr.Unlock()
}
