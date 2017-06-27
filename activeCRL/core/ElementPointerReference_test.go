package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewElementPointerReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	el1 := uOfD.NewElementPointerReference()
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

func TestElementPointerReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElementPointerReference()
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

func TestSetReferencedElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElementPointerReference()
	child.SetOwningElement(parent)
	if child.GetElementPointer() != nil {
		t.Error("ElementPointerReference's element pointer not initialized to nil")
	}
	elementPointer := uOfD.NewReferencedElementPointer()
	child.SetElementPointer(elementPointer)
	if child.GetElementPointer() == nil {
		t.Error("ElementPointerReference's  element pointer is nil after assignment")
		Print(elementPointer, "   ")
	}
	if child.GetElementPointer() != nil && child.GetElementPointer().GetId() != elementPointer.GetId() {
		t.Error("ElementPointerReference's  element pointer not set properly")
	}
	child.SetElementPointer(nil)
	if child.GetElementPointer() != nil {
		t.Error("ElementPointerReference's  element pointer not nild properly")
	}
}

func TestElementPointerReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElementPointerReference()
	child.SetOwningElement(parent)
	elementPointer := uOfD.NewReferencedElementPointer()
	child.SetElementPointer(elementPointer)

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

func TestElementPointerReferenceClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElementPointerReference()
	child.SetOwningElement(parent)
	elementPointer := uOfD.NewReferencedElementPointer()
	child.SetElementPointer(elementPointer)

	clone := child.(*elementPointerReference).clone()
	if !Equivalent(child, clone) {
		t.Error("Element clone failed")
		Print(child, "   ")
		Print(clone, "   ")
	}

}
