package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewElementReference(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	el1 := NewElementReference(uOfD)
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

func TestElementReferenceOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewElementReference(uOfD)
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

func TestSetReferencedElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewElementReference(uOfD)
	child.SetOwningElement(parent)
	if child.GetReferencedElement() != nil {
		t.Error("ElementReference's referenced element not initialized to nil")
	}
	child.SetReferencedElement(parent)
	if child.GetReferencedElement() == nil {
		t.Error("ElementReference's referenced element is nil after assignment")
		Print(parent, "   ")
	}
	if child.GetReferencedElement() != nil && child.GetReferencedElement().GetId() != parent.GetId() {
		t.Error("ElementReference's referenced element not set properly")
	}
	child.SetReferencedElement(nil)
	if child.GetReferencedElement() != nil {
		t.Error("ElementReference's referenced element not nild properly")
	}
}

func TestElementReferenceMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewElementReference(uOfD)
	child.SetOwningElement(parent)
	child.SetReferencedElement(parent)

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
