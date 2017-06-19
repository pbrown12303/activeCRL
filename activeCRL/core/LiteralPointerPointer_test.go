package core

import (
	"encoding/json"
	//	"fmt"
	"testing"
)

func TestNewLiteralPointerPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewLiteralPointerReference(uOfD)
	child := NewLiteralPointerPointer(uOfD)
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
		t.Error("LiteralPointerPointer not found in parent's OwnedBaseElements \n")
	}
	if parent.getLiteralPointerPointer() != child {
		t.Error("LiteralPointerReference.GetLiteralPointer() did not return child")
	}
}

func TestSetLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewLiteralPointerReference(uOfD)
	child := NewLiteralPointerPointer(uOfD)
	child.SetOwningElement(parent)
	literalPointer := NewValueLiteralPointer(uOfD)
	literalPointer.SetOwningElement(parent)
	if child.GetLiteralPointer() != nil {
		t.Error("LiteralPointer's Literal not initially nil \n")
	}
	child.SetLiteralPointer(literalPointer)
	if child.GetLiteralPointer() != literalPointer {
		t.Error("LiteralPointerPointer's LiteralPointer not properly set after assignment \n")
	}
}

func TestLiteralPointerPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewLiteralPointerReference(uOfD)
	child := NewLiteralPointerPointer(uOfD)
	child.SetOwningElement(parent)
	literalPointer := NewValueLiteralPointer(uOfD)
	literalPointer.SetOwningElement(parent)
	child.SetLiteralPointer(literalPointer)

	result, err := json.MarshalIndent(parent, "", "   ")
	if err != nil {
		t.Error(err)
	}

	//	fmt.Printf("Encoded Parent \n%s \n", result)

	uOfD2 := NewUniverseOfDiscourse()
	recoveredParent := RecoverElement(result, uOfD2)
	if !Equivalent(parent, recoveredParent) {
		t.Error("Recovered parent not equivalent to original parent")
	}
}

func TestLiteralPointerPointerClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewLiteralPointerReference(uOfD)
	child := NewLiteralPointerPointer(uOfD)
	child.SetOwningElement(parent)
	literalPointer := NewValueLiteralPointer(uOfD)
	literalPointer.SetOwningElement(parent)
	child.SetLiteralPointer(literalPointer)
	clone := child.(*literalPointerPointer).clone()
	if !Equivalent(child, clone) {
		t.Error("LiteralPointerPointer clone failed")
		Print(child, "   ")
		Print(clone, "   ")
	}

}
