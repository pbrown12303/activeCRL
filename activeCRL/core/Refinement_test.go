package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewRefinement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	el1 := uOfD.NewRefinement()
	if el1.GetId() == uuid.Nil {
		t.Error("Refinement identifier not properly initialized")
	}
	if el1.GetVersion() != 0 {
		t.Error("Refinement version not properly initialized")
	}
	if el1.getOwnedBaseElements() == nil {
		t.Error("Refinement ownedBaseElements not properly initialized")
	}
}

func TestRefinementOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewRefinement()
	child.SetOwningElement(parent)
	if child.GetOwningElement() != parent {
		t.Error("Child's owner not set properly")
	}
	if child.getOwningElementPointer() == nil {
		t.Error("Child's owningElementPointer not initialized properly")
	}
	if child.getOwningElementPointer().GetOwningElement().GetId() != child.GetId() {
		t.Error("Child's owningElementPointer.getOwningElement() != child")
	}
	if child.getOwningElementPointer().GetElement() != parent {
		t.Error("Child's owningElementPointer.getElement() != parent")
	}
}

func TestSetAbstractElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewRefinement()
	child.SetOwningElement(parent)
	abstractElement := uOfD.NewElement()
	if child.GetAbstractElement() != nil {
		t.Error("Refinement's abstract element not initialized to nil")
	}
	child.SetAbstractElement(abstractElement)
	if child.GetAbstractElement() == nil {
		t.Error("Refinement's abstract element is nil after assignment")
		Print(parent, "   ")
	}
	if child.GetAbstractElement() != nil && child.GetAbstractElement().GetId() != abstractElement.GetId() {
		t.Error("Refinement's abstract element not set properly")
	}
	child.SetAbstractElement(nil)
	if child.GetAbstractElement() != nil {
		t.Error("Refinement's abstract element not nild properly")
	}
}

func TestSetRefinedElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewRefinement()
	child.SetOwningElement(parent)
	refinedElement := uOfD.NewElement()
	if child.GetRefinedElement() != nil {
		t.Error("Refinement's refined element not initialized to nil")
	}
	child.SetRefinedElement(refinedElement)
	if child.GetRefinedElement() == nil {
		t.Error("Refinement's refined element is nil after assignment")
		Print(parent, "   ")
	}
	if child.GetRefinedElement() != nil && child.GetRefinedElement().GetId() != refinedElement.GetId() {
		t.Error("Refinement's refined element not set properly")
	}
	child.SetRefinedElement(nil)
	if child.GetRefinedElement() != nil {
		t.Error("Refinement's refined element not nild properly")
	}
}

func TestRefinementMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewRefinement()
	child.SetOwningElement(parent)
	abstractElement := uOfD.NewElement()
	child.SetAbstractElement(abstractElement)
	refinedElement := uOfD.NewElement()
	child.SetRefinedElement(refinedElement)

	result, err := json.MarshalIndent(parent, "", "   ")
	if err != nil {
		t.Error(err)
	}

	//	fmt.Printf("Encoded Parent \n%s \n", result)

	uOfD2 := NewUniverseOfDiscourse()
	recoveredParent := uOfD2.RecoverElement(result)
	if recoveredParent != nil {
		//		Print(recoveredParent, "")
	}
	if !Equivalent(parent, recoveredParent) {
		t.Error("Recovered parent not equivalent to original parent")
	}
}

func TestRefinementClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewRefinement()
	child.SetOwningElement(parent)
	abstractElement := uOfD.NewElement()
	child.SetAbstractElement(abstractElement)
	refinedElement := uOfD.NewElement()
	child.SetRefinedElement(refinedElement)
	clone := child.(*refinement).clone()
	if !Equivalent(child, clone) {
		Print(child, "   ")
		Print(clone, "   ")
		t.Error("Cloned Refinement not equivalent to original")
	}

}
