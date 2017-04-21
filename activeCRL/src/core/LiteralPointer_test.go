package core

import (
	"encoding/json"
	//	"fmt"
	"testing"
)

func TestNewNameLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewNameLiteralPointer(uOfD)
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
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.getLiteralPointerRole() != NAME {
		t.Error("LiteralPointer role not NAME \n")
	}
}

func TestDefinitionNameLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewDefinitionLiteralPointer(uOfD)
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
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.getLiteralPointerRole() != DEFINITION {
		t.Error("LiteralPointer role not DEFINITION \n")
	}
}

func TestNewUriLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewUriLiteralPointer(uOfD)
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
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.getLiteralPointerRole() != URI {
		t.Error("LiteralPointer role not URI \n")
	}
}

func TestNewValueLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewValueLiteralPointer(uOfD)
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
		t.Error("LiteralPointer not found in parent's OwnedBaseElements \n")
	}
	if child.getLiteralPointerRole() != VALUE {
		t.Error("LiteralPointer role not VALUE \n")
	}
}

func TestSetLiteral(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewNameLiteralPointer(uOfD)
	child.SetOwningElement(parent)
	literal := NewLiteral(uOfD)
	literal.SetOwningElement(parent)
	if child.GetLiteral() != nil {
		t.Error("LiteralPointer's Literal not initially nil \n")
	}
	child.SetLiteral(literal)
	if child.GetLiteral() != literal {
		t.Error("LiteralPointer's Literal not properly set after assignment \n")
	}
}

func TestLiteralPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewNameLiteralPointer(uOfD)
	child.SetOwningElement(parent)
	literal := NewLiteral(uOfD)
	literal.SetOwningElement(parent)
	child.SetLiteral(literal)

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
