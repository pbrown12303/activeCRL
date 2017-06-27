package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewLiteralReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	el1 := uOfD.NewLiteralReference()
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

func TestLiteralReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteralReference()
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

func TestSetReferencedLiteral(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteralReference()
	child.SetOwningElement(parent)
	literal := uOfD.NewLiteral()
	if child.GetReferencedLiteral() != nil {
		t.Error("LiteralReference's referenced literal not initialized to nil")
	}
	child.SetReferencedLiteral(literal)
	if child.GetReferencedLiteral() == nil {
		t.Error("LiteralReference's referenced literal is nil after assignment")
		Print(parent, "   ")
	}
	if child.GetReferencedLiteral() != nil && child.GetReferencedLiteral().GetId() != literal.GetId() {
		t.Error("LiteralReference's referenced literal not set properly")
	}
	child.SetReferencedLiteral(nil)
	if child.GetReferencedLiteral() != nil {
		t.Error("LiteralReference's referenced literal not nild properly")
	}
}

func TestLiteralReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteralReference()
	child.SetOwningElement(parent)
	literal := uOfD.NewLiteral()
	child.SetReferencedLiteral(literal)
	//	fmt.Printf("Parent before encoding \n")
	//	Print(parent, "   ")

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

func TestLiteralReferenceClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewLiteralReference()
	child.SetOwningElement(parent)
	literal := uOfD.NewLiteral()
	child.SetReferencedLiteral(literal)
	clone := child.(*literalReference).clone()
	if !Equivalent(child, clone) {
		Print(child, "   ")
		Print(clone, "   ")
		t.Error("Cloned LiteralReference not equivalent to original")
	}
}
