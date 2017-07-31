package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"log"
	"reflect"
)

type refinement struct {
	element
}

func (rPtr *refinement) clone() *refinement {
	var clone refinement
	clone.ownedBaseElements = make(map[string]BaseElement)
	clone.cloneAttributes(*rPtr)
	return &clone
}

func (rPtr *refinement) cloneAttributes(source refinement) {
	rPtr.element.cloneAttributes(source.element)
}

func (rPtr *refinement) GetAbstractElement(hl *HeldLocks) Element {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(rPtr)
	rep := rPtr.getAbstractElementPointer(hl)
	if rep != nil {
		return rep.GetElement(hl)
	}
	return nil
}

func (rPtr *refinement) getAbstractElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(rPtr)
	for _, be := range rPtr.ownedBaseElements {
		switch be.(type) {
		case ElementPointer:
			if be.(ElementPointer).GetElementPointerRole(hl) == ABSTRACT_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

func (rPtr *refinement) GetRefinedElement(hl *HeldLocks) Element {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(rPtr)
	rep := rPtr.getRefinedElementPointer(hl)
	if rep != nil {
		return rep.GetElement(hl)
	}
	return nil
}

func (rPtr *refinement) getRefinedElementPointer(hl *HeldLocks) ElementPointer {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(rPtr)
	for _, be := range rPtr.ownedBaseElements {
		switch be.(type) {
		case ElementPointer:
			if be.(ElementPointer).GetElementPointerRole(hl) == REFINED_ELEMENT {
				return be.(ElementPointer)
			}
		}
	}
	return nil
}

func (rPtr *refinement) initializeRefinement() {
	rPtr.initializeElement()
}

func (bePtr *refinement) isEquivalent(be *refinement, hl *HeldLocks) bool {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(bePtr)
	var elementPtr *element = &bePtr.element
	return elementPtr.isEquivalent(&be.element, hl)
}

func (elPtr *refinement) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	typeName := reflect.TypeOf(elPtr).String()
	buffer.WriteString(fmt.Sprintf("\"Type\":\"%s\",", typeName))
	err := elPtr.marshalRefinementFields(buffer)
	buffer.WriteString("}")
	return buffer.Bytes(), err
}

func (elPtr *refinement) marshalRefinementFields(buffer *bytes.Buffer) error {
	return elPtr.element.marshalElementFields(buffer)
}

func (elPtr *refinement) printRefinement(prefix string, hl *HeldLocks) {
	elPtr.printElement(prefix, hl)
}

func (el *refinement) recoverRefinementFields(unmarshaledData *map[string]json.RawMessage) error {
	return el.element.recoverElementFields(unmarshaledData)
}

func (rPtr *refinement) SetAbstractElement(el Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(rPtr)
	if rPtr.GetAbstractElement(hl) != el {
		ep := rPtr.getAbstractElementPointer(hl)
		if ep == nil {
			ep = rPtr.uOfD.NewAbstractElementPointer(hl)
			SetOwningElement(ep, rPtr, hl)
		}
		ep.SetElement(el, hl)
	}
}

func (rPtr *refinement) SetRefinedElement(el Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks()
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(rPtr)
	if rPtr.GetRefinedElement(hl) != el {
		ep := rPtr.getRefinedElementPointer(hl)
		if ep == nil {
			ep = rPtr.uOfD.NewRefinedElementPointer(hl)
			SetOwningElement(ep, rPtr, hl)
		}
		ep.SetElement(el, hl)
	}
}

type Refinement interface {
	Element
	GetAbstractElement(*HeldLocks) Element
	getAbstractElementPointer(*HeldLocks) ElementPointer
	GetRefinedElement(*HeldLocks) Element
	getRefinedElementPointer(*HeldLocks) ElementPointer
	SetAbstractElement(Element, *HeldLocks)
	SetRefinedElement(Element, *HeldLocks)
}
