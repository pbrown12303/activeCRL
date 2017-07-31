package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewRefinement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	el1 := uOfD.NewRefinement(hl)
	if el1.GetId(hl) == uuid.Nil {
		t.Error("Refinement identifier not properly initialized")
	}
	if el1.GetVersion(hl) != 0 {
		t.Error("Refinement version not properly initialized")
	}
	if el1.GetOwnedBaseElements(hl) != nil {
		t.Error("Refinement ownedBaseElements not properly initialized")
	}
}

func TestRefinementOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	if child.GetOwningElement(hl) != parent {
		t.Error("Child's owner not set properly")
	}
	if child.getOwningElementPointer(hl) == nil {
		t.Error("Child's owningElementPointer not initialized properly")
	}
	if GetOwningElement(child.getOwningElementPointer(hl), hl).GetId(hl) != child.GetId(hl) {
		t.Error("Child's owningElementPointer.getOwningElement() != child")
	}
	if child.getOwningElementPointer(hl).GetElement(hl) != parent {
		t.Error("Child's owningElementPointer.getElement() != parent")
	}
}

func TestSetAbstractElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	abstractElement := uOfD.NewElement(hl)
	if child.GetAbstractElement(hl) != nil {
		t.Error("Refinement's abstract element not initialized to nil")
	}
	child.SetAbstractElement(abstractElement, hl)
	if child.GetAbstractElement(hl) == nil {
		t.Error("Refinement's abstract element is nil after assignment")
		Print(parent, "   ", hl)
	}
	if child.GetAbstractElement(hl) != nil && child.GetAbstractElement(hl).GetId(hl) != abstractElement.GetId(hl) {
		t.Error("Refinement's abstract element not set properly")
	}
	child.SetAbstractElement(nil, hl)
	if child.GetAbstractElement(hl) != nil {
		t.Error("Refinement's abstract element not nild properly")
	}
}

func TestSetRefinedElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	refinedElement := uOfD.NewElement(hl)
	if child.GetRefinedElement(hl) != nil {
		t.Error("Refinement's refined element not initialized to nil")
	}
	child.SetRefinedElement(refinedElement, hl)
	if child.GetRefinedElement(hl) == nil {
		t.Error("Refinement's refined element is nil after assignment")
		Print(parent, "   ", hl)
	}
	if child.GetRefinedElement(hl) != nil && child.GetRefinedElement(hl).GetId(hl) != refinedElement.GetId(hl) {
		t.Error("Refinement's refined element not set properly")
	}
	child.SetRefinedElement(nil, hl)
	if child.GetRefinedElement(hl) != nil {
		t.Error("Refinement's refined element not nild properly")
	}
}

func TestRefinementMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	abstractElement := uOfD.NewElement(hl)
	child.SetAbstractElement(abstractElement, hl)
	refinedElement := uOfD.NewElement(hl)
	child.SetRefinedElement(refinedElement, hl)

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
	if !Equivalent(parent, recoveredParent, hl) {
		t.Error("Recovered parent not equivalent to original parent")
	}
}

func TestRefinementClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	abstractElement := uOfD.NewElement(hl)
	child.SetAbstractElement(abstractElement, hl)
	refinedElement := uOfD.NewElement(hl)
	child.SetRefinedElement(refinedElement, hl)
	clone := child.(*refinement).clone()
	if !Equivalent(child, clone, hl) {
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
		t.Error("Cloned Refinement not equivalent to original")
	}

}
