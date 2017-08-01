package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewElementPointerReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	el1 := uOfD.NewElementPointerReference(hl)
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

func TestElementPointerReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementPointerReference(hl)
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

func TestSetReferencedElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementPointerReference(hl)
	SetOwningElement(child, parent, hl)
	if child.GetElementPointer(hl) != nil {
		t.Error("ElementPointerReference's element pointer not initialized to nil")
	}
	elementPointer := uOfD.NewReferencedElementPointer(hl)
	child.SetElementPointer(elementPointer, hl)
	if child.GetElementPointer(hl) == nil {
		t.Error("ElementPointerReference's  element pointer is nil after assignment")
		Print(elementPointer, "   ", hl)
	}
	if child.GetElementPointer(hl) != nil && child.GetElementPointer(hl).GetId(hl) != elementPointer.GetId(hl) {
		t.Error("ElementPointerReference's  element pointer not set properly")
	}
	child.SetElementPointer(nil, hl)
	if child.GetElementPointer(hl) != nil {
		t.Error("ElementPointerReference's  element pointer not nild properly")
	}
}

func TestElementPointerReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementPointerReference(hl)
	SetOwningElement(child, parent, hl)
	elementPointer := uOfD.NewReferencedElementPointer(hl)
	child.SetElementPointer(elementPointer, hl)

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

func TestElementPointerReferenceClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementPointerReference(hl)
	SetOwningElement(child, parent, hl)
	elementPointer := uOfD.NewReferencedElementPointer(hl)
	child.SetElementPointer(elementPointer, hl)

	clone := child.(*elementPointerReference).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("Element clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}

}
