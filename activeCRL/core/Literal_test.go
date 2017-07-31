package core

import (
	"encoding/json"
	//	"fmt"
	"testing"
)

func TestNewLiteral(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteral(hl)
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
		t.Error("Literal not found in parent's OwnedBaseElements \n")
	}
}

func TestLiteralMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteral(hl)
	SetOwningElement(child, parent, hl)
	var testString string = "Test String"
	child.SetLiteralValue(testString, hl)

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

func TestLiteralClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteral(hl)
	SetOwningElement(child, parent, hl)
	var testString string = "Test String"
	child.SetLiteralValue(testString, hl)
	clone := child.(*literal).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("Literal clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}
}
