package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewElementReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	el1 := uOfD.NewElementReference(hl)
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

func TestElementReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementReference(hl)
	SetOwningElement(child, parent, hl)
	if child.GetOwningElement(hl) != parent {
		t.Error("Child's owner not set properly")
	}
	if child.GetOwningElementPointer(hl) == nil {
		t.Error("Child's owningElementPointer not initialized properly")
	}
	if GetOwningElement(child.GetOwningElementPointer(hl), hl).GetId(hl) != child.GetId(hl) {
		t.Error("Child's owningElementPointer.getOwningElement() != child")
	}
	if child.GetOwningElementPointer(hl).GetElement(hl) != parent {
		t.Error("Child's owningElementPointer.getElement() != parent")
	}
}

func TestSetReferencedElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementReference(hl)
	SetOwningElement(child, parent, hl)
	if child.GetReferencedElement(hl) != nil {
		t.Error("ElementReference's referenced element not initialized to nil")
	}
	child.SetReferencedElement(parent, hl)
	if child.GetReferencedElement(hl) == nil {
		t.Error("ElementReference's referenced element is nil after assignment")
		Print(parent, "   ", hl)
	}
	if child.GetReferencedElement(hl) != nil && child.GetReferencedElement(hl).GetId(hl) != parent.GetId(hl) {
		t.Error("ElementReference's referenced element not set properly")
	}
	child.SetReferencedElement(nil, hl)
	if child.GetReferencedElement(hl) != nil {
		t.Error("ElementReference's referenced element not nild properly")
	}
}

func TestElementReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementReference(hl)
	SetOwningElement(child, parent, hl)
	child.SetReferencedElement(parent, hl)

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

func TestElementReferenceClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementReference(hl)
	SetOwningElement(child, parent, hl)
	child.SetReferencedElement(parent, hl)
	clone := child.(*elementReference).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("ElementReference clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}
}
