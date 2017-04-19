package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewElement(t *testing.T) {
	var uOfD UniverseOfDiscourse
	el1 := NewElement(&uOfD)
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

func TestElementOwnership(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewElement(&uOfD)
	child.setOwningElement(parent)
	if child.GetOwningElement() != parent {
		t.Error("Child's owner not set properly")
	}
	if child.getOwningElementPointer() == nil {
		t.Error("Child's owningElementPointer not initialized properly")
	}
	if child.getOwningElementPointer().GetOwningElement() != child {
		t.Error("Child's owningElementPointer.getOwningElement() != child")
	}
	if child.getOwningElementPointer().GetElement() != parent {
		t.Error("Child's owningElementPointer.getElement() != parent")
	}
}

func TestElementMarshal(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewElement(&uOfD)
	child.setOwningElement(parent)

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

func TestSetName(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	var testName string = "Test Name"
	parent.SetName(testName)
	if parent.GetName() != testName {
		t.Error("GetName() value not equal to assigned name")
	}
	if parent.getNameLiteral() == nil {
		t.Error("getNameLiteral() is nil after name assigned")
	}
	if parent.getNameLiteralPointer() == nil {
		t.Error("getNameLiteralPointer() is nil after name assigned")

	}
}

func TestSetDefinition(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	var testName string = "Test Name"
	parent.SetDefinition(testName)
	if parent.GetDefinition() != testName {
		t.Error("GetName() value not equal to assigned name")
	}
	if parent.getDefinitionLiteral() == nil {
		t.Error("getNameLiteral() is nil after name assigned")
	}
	if parent.getDefinitionLiteralPointer() == nil {
		t.Error("getNameLiteralPointer() is nil after name assigned")

	}
}

func TestSetUri(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	var testName string = "Test Name"
	parent.SetUri(testName)
	if parent.GetUri() != testName {
		t.Error("GetName() value not equal to assigned name")
	}
	if parent.getUriLiteral() == nil {
		t.Error("getNameLiteral() is nil after name assigned")
	}
	if parent.getUriLiteralPointer() == nil {
		t.Error("getNameLiteralPointer() is nil after name assigned")

	}
}
