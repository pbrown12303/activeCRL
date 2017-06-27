package core

import (
	"encoding/json"
	//	"fmt"
	"testing"
)

func TestNewLiteral(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteral()
	child.SetOwningElement(parent)
	if child.GetOwningElement() != parent {
		t.Error("Child's owner not set properly")
	}
	var found bool = false
	for key, _ := range parent.getOwnedBaseElements() {
		if key == child.GetId().String() {
			found = true
		}
	}
	if found == false {
		t.Error("Literal not found in parent's OwnedBaseElements \n")
	}
}

func TestLiteralMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteral()
	child.SetOwningElement(parent)
	var testString string = "Test String"
	child.SetLiteralValue(testString)

	result, err := json.MarshalIndent(parent, "", "   ")
	if err != nil {
		t.Error(err)
	}

	//	fmt.Printf("Encoded Parent \n%s \n", result)

	uOfD2 := NewUniverseOfDiscourse()
	recoveredParent := RecoverElement(result, uOfD2)
	if recoveredParent != nil {
		//		Print(recoveredParent, "")
	}
	if !Equivalent(parent, recoveredParent) {
		t.Error("Recovered parent not equivalent to original parent")
	}
}

func TestLiteralClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteral()
	child.SetOwningElement(parent)
	var testString string = "Test String"
	child.SetLiteralValue(testString)
	clone := child.(*literal).clone()
	if !Equivalent(child, clone) {
		t.Error("Literal clone failed")
		Print(child, "   ")
		Print(clone, "   ")
	}
}
