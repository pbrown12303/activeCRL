package core

import (
	"encoding/json"
	//	"fmt"
	"testing"
)

func TestNewNameLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewNameLiteralPointer(hl)
	SetOwningElement(child, parent, hl)
	if GetOwningElement(child, hl) != parent {
		t.Error("Child's owner not set properly")
	}
	var found bool = false
	for _, be := range parent.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == child.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.GetLiteralPointerRole(hl) != NAME {
		t.Error("LiteralPointer role not NAME \n")
	}
}

func TestDefinitionNameLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewDefinitionLiteralPointer(hl)
	SetOwningElement(child, parent, hl)
	if GetOwningElement(child, hl) != parent {
		t.Error("Child's owner not set properly")
	}
	var found bool = false
	for _, be := range parent.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == child.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.GetLiteralPointerRole(hl) != DEFINITION {
		t.Error("LiteralPointer role not DEFINITION \n")
	}
}

func TestNewUriLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewUriLiteralPointer(hl)
	SetOwningElement(child, parent, hl)
	if GetOwningElement(child, hl) != parent {
		t.Error("Child's owner not set properly")
	}
	var found bool = false
	for _, be := range parent.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == child.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.GetLiteralPointerRole(hl) != URI {
		t.Error("LiteralPointer role not URI \n")
	}
}

func TestNewValueLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewValueLiteralPointer(hl)
	SetOwningElement(child, parent, hl)
	if GetOwningElement(child, hl) != parent {
		t.Error("Child's owner not set properly")
	}
	var found bool = false
	for _, be := range parent.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == child.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.GetLiteralPointerRole(hl) != VALUE {
		t.Error("LiteralPointer role not VALUE \n")
	}
}

func TestSetLiteral(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewNameLiteralPointer(hl)
	SetOwningElement(child, parent, hl)
	literal := uOfD.NewLiteral(hl)
	SetOwningElement(literal, parent, hl)
	if child.GetLiteral(hl) != nil {
		t.Error("LiteralPointer's Literal not initially nil \n")
	}
	child.SetLiteral(literal, hl)
	if child.GetLiteral(hl) != literal {
		t.Error("LiteralPointer's Literal not properly set after assignment \n")
	}
}

func TestLiteralPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewNameLiteralPointer(hl)
	SetOwningElement(child, parent, hl)
	literal := uOfD.NewLiteral(hl)
	SetOwningElement(literal, parent, hl)
	child.SetLiteral(literal, hl)

	result, err := json.MarshalIndent(parent, "", "   ")
	if err != nil {
		t.Error(err)
	}

	//	fmt.Printf("Encoded Parent \n%s \n", result)

	uOfD2 := NewUniverseOfDiscourse()
	recoveredParent := uOfD2.RecoverElement(result)
	if !Equivalent(parent, recoveredParent, hl) {
		t.Error("Recovered parent not equivalent to original parent")
	}
}

func TestLiteralPointerClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewNameLiteralPointer(hl)
	SetOwningElement(child, parent, hl)
	literal := uOfD.NewLiteral(hl)
	SetOwningElement(literal, parent, hl)
	child.SetLiteral(literal, hl)
	clone := child.(*literalPointer).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("LiteralPointer clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}
}
