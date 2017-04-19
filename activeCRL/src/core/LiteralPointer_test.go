package core

import (
	"encoding/json"
	//	"fmt"
	"testing"
)

func TestNewNameLiteralPointer(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewNameLiteralPointer(&uOfD)
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
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.getLiteralPointerRole() != NAME {
		t.Error("LiteralPointer role not NAME \n")
	}
}

func TestDefinitionNameLiteralPointer(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewDefinitionLiteralPointer(&uOfD)
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
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.getLiteralPointerRole() != DEFINITION {
		t.Error("LiteralPointer role not DEFINITION \n")
	}
}

func TestNewUriLiteralPointer(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewUriLiteralPointer(&uOfD)
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
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.getLiteralPointerRole() != URI {
		t.Error("LiteralPointer role not URI \n")
	}
}

func TestNewValueLiteralPointer(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewValueLiteralPointer(&uOfD)
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
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.getLiteralPointerRole() != VALUE {
		t.Error("LiteralPointer role not VALUE \n")
	}
}

func TestSetLiteral(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewNameLiteralPointer(&uOfD)
	child.setOwningElement(parent)
	literal := NewLiteral(&uOfD)
	literal.setOwningElement(parent)
	if child.GetLiteral() != nil {
		t.Error("LiteralPointer's Literal not initially nil \n")
	}
	child.SetLiteral(literal)
	if child.GetLiteral() != literal {
		t.Error("LiteralPointer's Literal not properly set after assignment \n")
	}
}

func TestLiteralPointerMarshal(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewNameLiteralPointer(&uOfD)
	child.setOwningElement(parent)
	literal := NewLiteral(&uOfD)
	literal.setOwningElement(parent)
	child.SetLiteral(literal)

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
