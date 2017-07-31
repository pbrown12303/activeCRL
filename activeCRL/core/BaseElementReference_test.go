package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewBaseElementReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	el1 := uOfD.NewBaseElementReference(hl)
	if el1.GetId(hl) == uuid.Nil {
		t.Error("Element identifier not properly initialized")
	}
	if el1.GetVersion(hl) != 0 {
		t.Error("Element version not properly initialized")
	}
	if el1.GetOwnedBaseElements(hl) != nil {
		t.Error("Element ownedBaseElements not properly initialized")
	}
}

func TestBaseElementReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewBaseElementReference(hl)
	SetOwningElement(child, parent, hl)
	//	Print(parent, "", hl)
	//	Print(child, "", hl)
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

func TestReferenceSetBaseElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewBaseElementReference(hl)
	SetOwningElement(child, parent, hl)
	if child.GetBaseElement(hl) != nil {
		t.Error("BaseElementReference's base element not initialized to nil")
	}
	child.SetBaseElement(parent, hl)
	if child.GetBaseElement(hl) == nil {
		t.Error("BaseElementReference's base element is nil after assignment")
		Print(parent, "   ", hl)
	}
	if child.GetBaseElement(hl) != nil && child.GetBaseElement(hl).GetId(hl) != parent.GetId(hl) {
		t.Error("BaseElementReference's base element not set properly")
	}
	child.SetBaseElement(nil, hl)
	if child.GetBaseElement(hl) != nil {
		t.Error("BaseElementReference's base element not nild properly")
	}
}

func TestBaseElementReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewBaseElementReference(hl)
	SetOwningElement(child, parent, hl)
	child.SetBaseElement(parent, hl)

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

func TestBaseElementReferenceClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewBaseElementReference(hl)
	SetOwningElement(child, parent, hl)
	child.SetBaseElement(parent, hl)
	clone := child.(*baseElementReference).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("ElementReference clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}
}
