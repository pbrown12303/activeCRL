package core

import (
	"encoding/json"
	//	"fmt"
	"testing"
)

func TestNewLiteral(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewLiteral(&uOfD)
	child.setOwningElement(parent)
	if child.GetOwningElement() != parent {
		t.Error("Child's owner not set properly")
	}
	var found bool = false
	for key, _ := range parent.GetOwnedBaseElements() {
		if key == child.GetId().String() {
			found = true
		}
	}
	if found == false {
		t.Error("Literal not found in parent's OwnedBaseElements \n")
	}
}

func TestLiteralMarshal(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewLiteral(&uOfD)
	child.setOwningElement(parent)
	var testString string = "Test String"
	child.SetLiteralValue(testString)

	result, err := json.MarshalIndent(parent, "", "   ")
	if err != nil {
		t.Error(err)
	}

	//	fmt.Printf("Encoded Parent \n%s \n", result)

	var uOfD2 UniverseOfDiscourse
	recoveredParent := RecoverElement(result, &uOfD2)
	if recoveredParent != nil {
		//		Print(recoveredParent, "")
	}
	if !Equivalent(parent, recoveredParent) {
		t.Error("Recovered parent not equivalent to original parent")
	}

}
