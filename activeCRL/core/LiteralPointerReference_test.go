package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewLiteralPointerReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	el1 := uOfD.NewLiteralPointerReference()
	if el1.GetId() == uuid.Nil {
		t.Error("Element identifier not properly initialized")
	}
	if el1.GetVersion() != 0 {
		t.Error("Element version not properly initialized")
	}
	if el1.getOwnedBaseElements() == nil {
		t.Error("Element ownedBaseElements not properly initialized")
	}
}

func TestLiteralPointerReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteralPointerReference()
	child.SetOwningElement(parent)
	if child.GetOwningElement() != parent {
		t.Error("Child's owner not set properly")
	}
	if child.getOwningElementPointer() == nil {
		t.Error("Child's owningElementPointer not initialized properly")
	}
	if child.getOwningElementPointer().GetOwningElement().GetId() != child.GetId() {
		t.Error("Child's owningElementPointer.getOwningElement() != child")
	}
	if child.getOwningElementPointer().GetElement() != parent {
		t.Error("Child's owningElementPointer.getElement() != parent")
	}
}

func TestSetReferencedLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteralPointerReference()
	child.SetOwningElement(parent)
	if child.GetLiteralPointer() != nil {
		t.Error("LiteralPointerReference's element pointer not initialized to nil")
	}
	literalPointer := uOfD.NewValueLiteralPointer()
	child.SetLiteralPointer(literalPointer)
	if child.GetLiteralPointer() == nil {
		t.Error("LiteralPointerReference's  element pointer is nil after assignment")
		Print(literalPointer, "   ")
	}
	if child.GetLiteralPointer() != nil && child.GetLiteralPointer().GetId() != literalPointer.GetId() {
		t.Error("LiteralPointerReference's  element pointer not set properly")
	}
	child.SetLiteralPointer(nil)
	if child.GetLiteralPointer() != nil {
		t.Error("LiteralPointerReference's  element pointer not nild properly")
	}
}

func TestLiteralPointerReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteralPointerReference()
	child.SetOwningElement(parent)
	literalPointer := uOfD.NewValueLiteralPointer()
	child.SetLiteralPointer(literalPointer)

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
	if !Equivalent(parent, recoveredParent) {
		t.Error("Recovered parent not equivalent to original parent")
	}
}

func TestLiteralPointerReferenceClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteralPointerReference()
	child.SetOwningElement(parent)
	literalPointer := uOfD.NewValueLiteralPointer()
	child.SetLiteralPointer(literalPointer)
	clone := child.(*literalPointerReference).clone()
	if !Equivalent(child, clone) {
		t.Error("LiteralPointerReference clone failed")
		Print(child, "   ")
		Print(clone, "   ")
	}
}
