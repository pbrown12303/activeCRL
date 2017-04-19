package core

import (
	"errors"
	"log"
)

type UniverseOfDiscourse struct {
	baseElementMap map[string]BaseElement
}

func (uOfDPtr *UniverseOfDiscourse) addBaseElement(be BaseElement) error {
	if be == nil {
		return errors.New("UniverseOfDiscource addBaseElement failed because base element was nil")
	}
	if uOfDPtr.baseElementMap[be.GetId().String()] != nil {
		return errors.New("UniverseOfDiscource addBaseElement failed because UUID was already in map")
	}
	be.setUniverseOfDiscourse(uOfDPtr)
	return nil
}

func (uOfDPtr *UniverseOfDiscourse) getBaseElement(id string) BaseElement {
	return uOfDPtr.baseElementMap[id]
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
