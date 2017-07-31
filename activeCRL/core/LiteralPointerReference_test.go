package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewLiteralPointerReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	el1 := uOfD.NewLiteralPointerReference(hl)
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

func TestLiteralPointerReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteralPointerReference(hl)
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

func TestSetReferencedLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteralPointerReference(hl)
	SetOwningElement(child, parent, hl)
	if child.GetLiteralPointer(hl) != nil {
		t.Error("LiteralPointerReference's element pointer not initialized to nil")
	}
	literalPointer := uOfD.NewValueLiteralPointer(hl)
	child.SetLiteralPointer(literalPointer, hl)
	if child.GetLiteralPointer(hl) == nil {
		t.Error("LiteralPointerReference's  element pointer is nil after assignment")
		Print(literalPointer, "   ", hl)
	}
	if child.GetLiteralPointer(hl) != nil && child.GetLiteralPointer(hl).GetId(hl) != literalPointer.GetId(hl) {
		t.Error("LiteralPointerReference's  element pointer not set properly")
	}
	child.SetLiteralPointer(nil, hl)
	if child.GetLiteralPointer(hl) != nil {
		t.Error("LiteralPointerReference's  element pointer not nild properly")
	}
}

func TestLiteralPointerReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteralPointerReference(hl)
	SetOwningElement(child, parent, hl)
	literalPointer := uOfD.NewValueLiteralPointer(hl)
	child.SetLiteralPointer(literalPointer, hl)

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

func TestLiteralPointerReferenceClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteralPointerReference(hl)
	SetOwningElement(child, parent, hl)
	literalPointer := uOfD.NewValueLiteralPointer(hl)
	child.SetLiteralPointer(literalPointer, hl)
	clone := child.(*literalPointerReference).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("LiteralPointerReference clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}
}
