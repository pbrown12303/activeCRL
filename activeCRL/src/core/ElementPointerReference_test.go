package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewElementPointerReference(t *testing.T) {
	var uOfD UniverseOfDiscourse
	el1 := NewElementPointerReference(&uOfD)
	if el1.GetId() == uuid.Nil {
		t.Error("Element identifier not properly initialized")
	}
	if el1.GetVersion() != 0 {
		t.Error("Element version not properly initialized")
	}
	if el1.GetOwnedBaseElements() == nil {
		t.Error("Element ownedBaseElements not properly initialized")
	}
}

func TestElementPointerReferenceOwnership(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewElementPointerReference(&uOfD)
	child.setOwningElement(parent)
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
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewElementPointerReference(&uOfD)
	child.setOwningElement(parent)
	if child.GetElementPointer() != nil {
		t.Error("ElementPointerReference's element pointer not initialized to nil")
	}
	elementPointer := NewReferencedElementPointer(&uOfD)
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
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewElementPointerReference(&uOfD)
	child.setOwningElement(parent)
	elementPointer := NewReferencedElementPointer(&uOfD)
	child.SetElementPointer(elementPointer)

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
