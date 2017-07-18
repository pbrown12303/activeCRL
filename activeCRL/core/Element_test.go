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
	var found bool = false
	for _, be := range parent.getOwnedBaseElements() {
		if be.GetId() == child.GetId() {
			found = true
		}
	}
	if found == false {
		t.Error("Parent does not contain child in getOwnedBaseElements()")
	}
	found = false
	for _, be := range parent.GetOwnedBaseElements() {
		if be.GetId() == child.GetId() {
			found = true
		}
	}
	if found == false {
		t.Error("Parent does not contain child in GetOwnedBaseElements()")
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
	recoveredParent := uOfD2.RecoverElement(result)
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

func TestGetImmediateAbstractions(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	refinedElement := uOfD.NewElement()
	abstractions := refinedElement.getImmediateAbstractions()
	if len(abstractions) != 0 {
		t.Error("Abstractions length not 0\n")
	}
	abstractElement1 := uOfD.NewElement()
	refinement1 := uOfD.NewRefinement()
	refinement1.SetAbstractElement(abstractElement1)
	refinement1.SetRefinedElement(refinedElement)
	abstractions = refinedElement.getImmediateAbstractions()
	if len(abstractions) != 1 {
		t.Error("Abstractions length != 1")
	}
	if abstractions[0] != refinement1 {
		t.Error("Abstractions first element not refinement1")
	}
	abstractElement2 := uOfD.NewElement()
	refinement2 := uOfD.NewRefinement()
	refinement2.SetAbstractElement(abstractElement2)
	refinement2.SetRefinedElement(refinedElement)
	abstractions = refinedElement.getImmediateAbstractions()
	if len(abstractions) != 2 {
		t.Error("Abstractions length != 2")
	}
	if abstractions[1] != refinement2 {
		t.Error("Abstractions second element not refinement2")
	}
}

func TestGetImmediateAbstractElements(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	refinedElement := uOfD.NewElement()
	abstractElements := refinedElement.getImmediateAbstractElements()
	if len(abstractElements) != 0 {
		t.Error("AbstractElements length not 0\n")
	}
	abstractElement1 := uOfD.NewElement()
	refinement1 := uOfD.NewRefinement()
	refinement1.SetAbstractElement(abstractElement1)
	refinement1.SetRefinedElement(refinedElement)
	abstractElements = refinedElement.getImmediateAbstractElements()
	if len(abstractElements) != 1 {
		t.Error("AbstractElements length != 1")
	}
	if abstractElements[0] != abstractElement1 {
		t.Error("AbstractElements first element not abstractElement1")
	}
	abstractElement2 := uOfD.NewElement()
	refinement2 := uOfD.NewRefinement()
	refinement2.SetAbstractElement(abstractElement2)
	refinement2.SetRefinedElement(refinedElement)
	abstractElements = refinedElement.getImmediateAbstractElements()
	if len(abstractElements) != 2 {
		t.Error("Abstractions length != 2")
	}
	if abstractElements[1] != abstractElement2 {
		t.Error("AbstractElements second element not abstractElement2")
	}
}

func TestGetAbstractElementsRecursively(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	refinedElement := uOfD.NewElement()
	abstractElements := refinedElement.GetAbstractElementsRecursivelyNoLock()
	if len(abstractElements) != 0 {
		t.Error("AbstractElements length not 0\n")
	}
	abstractElement1 := uOfD.NewElement()
	refinement1 := uOfD.NewRefinement()
	refinement1.SetAbstractElement(abstractElement1)
	refinement1.SetRefinedElement(refinedElement)
	abstractElements = refinedElement.GetAbstractElementsRecursivelyNoLock()
	if len(abstractElements) != 1 {
		t.Error("AbstractElements length != 1")
	}
	if abstractElements[0] != abstractElement1 {
		t.Error("AbstractElements first element not abstractElement1")
	}
	abstractElement2 := uOfD.NewElement()
	refinement2 := uOfD.NewRefinement()
	refinement2.SetAbstractElement(abstractElement2)
	refinement2.SetRefinedElement(refinedElement)
	abstractElements = refinedElement.GetAbstractElementsRecursivelyNoLock()
	if len(abstractElements) != 2 {
		t.Error("Abstractions length != 2")
	}
	if abstractElements[1] != abstractElement2 {
		t.Error("AbstractElements second element not abstractElement2")
	}
	abstractElement3 := uOfD.NewElement()
	refinement3 := uOfD.NewRefinement()
	refinement3.SetAbstractElement(abstractElement3)
	refinement3.SetRefinedElement(abstractElement1)
	abstractElements = refinedElement.GetAbstractElementsRecursivelyNoLock()
	if len(abstractElements) != 3 {
		t.Error("Abstractions length != 3")
	}
	if abstractElements[2] != abstractElement3 {
		t.Error("AbstractElements third element not abstractElement3")
	}
	abstractElement4 := uOfD.NewElement()
	refinement4 := uOfD.NewRefinement()
	refinement4.SetAbstractElement(abstractElement4)
	refinement4.SetRefinedElement(abstractElement2)
	abstractElements = refinedElement.GetAbstractElementsRecursivelyNoLock()
	if len(abstractElements) != 4 {
		t.Error("Abstractions length != 4")
	}
	if abstractElements[3] != abstractElement4 {
		t.Error("AbstractElements fourth element not abstractElement4")
	}
}

func TestGetImmediateRefinements(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	abstractElement := uOfD.NewElement()
	refinedElement1 := uOfD.NewElement()
	refinedElement2 := uOfD.NewElement()
	refinements := abstractElement.getImmediateRefinements()
	if len(refinements) != 0 {
		t.Error("Refinements length not 0\n")
	}
	refinement1 := uOfD.NewRefinement()
	refinement1.SetAbstractElement(abstractElement)
	refinement1.SetRefinedElement(refinedElement1)
	refinements = abstractElement.getImmediateRefinements()
	if len(refinements) != 1 {
		t.Error("Refinements length != 1")
	}
	if refinements[0] != refinement1 {
		t.Error("Refinements first element not refinement1")
	}
	refinement2 := uOfD.NewRefinement()
	refinement2.SetAbstractElement(abstractElement)
	refinement2.SetRefinedElement(refinedElement2)
	refinements = abstractElement.getImmediateRefinements()
	if len(refinements) != 2 {
		t.Error("Refinements length != 2")
	}
	if refinements[1] != refinement2 {
		t.Error("Refinements second element not refinement2")
	}
}
