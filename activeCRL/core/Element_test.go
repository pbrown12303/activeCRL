package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNewElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	el1 := uOfD.NewElement()
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

func TestElementOwnership(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElement()
	child.SetOwningElement(parent)
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
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElement()
	child.SetOwningElement(parent)

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

func TestSetName(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
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
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
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
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
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

func TestVersionWithParentChange(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	oldParent := uOfD.NewElement()
	newParent := uOfD.NewElement()
	elementX := uOfD.NewElement()
	elementX.SetOwningElement(oldParent)
	oldParentInitialVersion := oldParent.GetVersion()
	newParentInitialVersion := newParent.GetVersion()
	elementXInitialVersion := elementX.GetVersion()
	elementXOwnerPointer := elementX.getOwningElementPointer()
	elementXOwnerPointerInitialVersion := elementXOwnerPointer.GetVersion()
	elementXOwnerPointer.SetElement(newParent)
	if elementX.GetOwningElement() != newParent {
		t.Error("elementX Owner not changed properly")
	}
	if !(elementXOwnerPointer.GetVersion() > elementXOwnerPointerInitialVersion) {
		t.Error("Owning element pointer version not incremented")
	}
	if !(elementX.GetVersion() > elementXInitialVersion) {
		t.Error("elementX version not incremented")
	}
	if !(oldParent.GetVersion() > oldParentInitialVersion) {
		t.Error("old parent version not incremented")
	}
	if !(newParent.GetVersion() > newParentInitialVersion) {
		t.Error("new parent version not incremented")
	}

}

func TestVersionWithParentChangeAndCommonGrandparent(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	grandparent := uOfD.NewElement()
	grandparentPreviousVersion := grandparent.GetVersion()
	oldParent := uOfD.NewElement()
	oldParentPreviousVersion := oldParent.GetVersion()
	oldParent.SetOwningElement(grandparent)
	if !(grandparent.GetVersion() > grandparentPreviousVersion) {
		t.Error("Grandparent version not incremented when old parent added as child")
	}
	grandparentPreviousVersion = grandparent.GetVersion()
	if !(oldParent.GetVersion() > oldParentPreviousVersion) {
		t.Error("Old parent version not incremented when added as child to grandparent")
	}
	oldParentPreviousVersion = oldParent.GetVersion()

	newParent := uOfD.NewElement()
	newParentPreviousVersion := newParent.GetVersion()
	newParent.SetOwningElement(grandparent)
	if !(grandparent.GetVersion() > grandparentPreviousVersion) {
		t.Error("Grandparent version not incremented when new parent added as child")
	}
	grandparentPreviousVersion = grandparent.GetVersion()
	if !(newParent.GetVersion() > newParentPreviousVersion) {
		t.Error("New parent version not incremented when added as child to grandparent")
	}
	newParentPreviousVersion = newParent.GetVersion()

	elementX := uOfD.NewElement()
	elementXPreviousVersion := elementX.GetVersion()
	elementX.SetOwningElement(oldParent)
	if !(oldParent.GetVersion() > oldParentPreviousVersion) {
		t.Error("Old parent version not incremented when elementX added as child")
	}
	oldParentPreviousVersion = oldParent.GetVersion()
	if !(elementX.GetVersion() > elementXPreviousVersion) {
		t.Error("elementX version not incremented when added as a child to oldParent")
	}
	elementXPreviousVersion = elementX.GetVersion()

	elementXOwnerPointer := elementX.getOwningElementPointer()
	elementXOwnerPointerInitialVersion := elementXOwnerPointer.GetVersion()
	elementXOwnerPointer.SetElement(newParent)
	if elementX.GetOwningElement() != newParent {
		t.Error("elementX Owner not changed properly")
	}
	if !(elementXOwnerPointer.GetVersion() > elementXOwnerPointerInitialVersion) {
		t.Error("Owning element pointer version not incremented")
	}
	if !(elementX.GetVersion() > elementXPreviousVersion) {
		t.Error("elementX version not incremented when parent changed")
	}
	if !(oldParent.GetVersion() > oldParentPreviousVersion) {
		t.Error("old parent version not incremented when elementX removed as child")
	}
	if !(newParent.GetVersion() > newParentPreviousVersion) {
		t.Error("new parent version not incremented when elementX added as child")
	}
	if !(grandparent.GetVersion() > grandparentPreviousVersion) {
		t.Error("Grandparent version not incremented when elementX parent changed to new parent")
	}
}

func TestCloneElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	el := uOfD.NewElement()
	el.SetName("E1")
	el.SetUri("E1.testDomain.com")
	el.SetDefinition("The definition of E1")
	clone := el.(*element).clone()
	if !Equivalent(el, clone) {
		t.Error("Element clone failed")
		Print(el, "   ")
		Print(clone, "   ")
	}
}
