package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewBaseElementReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	el1 := uOfD.NewBaseElementReference()
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

func TestBaseElementReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewBaseElementReference()
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

func TestReferenceSetBaseElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewBaseElementReference()
	child.SetOwningElement(parent)
	if child.GetBaseElement() != nil {
		t.Error("BaseElementReference's base element not initialized to nil")
	}
	child.SetBaseElement(parent)
	if child.GetBaseElement() == nil {
		t.Error("BaseElementReference's base element is nil after assignment")
		Print(parent, "   ")
	}
	if child.GetBaseElement() != nil && child.GetBaseElement().GetId() != parent.GetId() {
		t.Error("BaseElementReference's base element not set properly")
	}
	child.SetBaseElement(nil)
	if child.GetBaseElement() != nil {
		t.Error("BaseElementReference's base element not nild properly")
	}
}

func TestBaseElementReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewBaseElementReference()
	child.SetOwningElement(parent)
	child.SetBaseElement(parent)

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

func TestBaseElementReferenceClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewBaseElementReference()
	child.SetOwningElement(parent)
	child.SetBaseElement(parent)
	clone := child.(*baseElementReference).clone()
	if !Equivalent(child, clone) {
		t.Error("ElementReference clone failed")
		Print(child, "   ")
		Print(clone, "   ")
	}
}
