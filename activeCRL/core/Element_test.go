// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	//	"log"
	"sync"
	"testing"
)

func TestCloneElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el := uOfD.NewElement(hl)
	SetLabel(el, "E1", hl)
	SetUri(el, "E1.testDomain.com", hl)
	SetDefinition(el, "The definition of E1", hl)
	clone := el.(*element).clone()
	if !Equivalent(el, clone, hl) {
		t.Error("Element clone failed")
		Print(el, "   ", hl)
		Print(clone, "   ", hl)
	}
}

func TestElementMarshal(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElement(hl)
	SetOwningElement(child, parent, hl)
	testLabel := "TestLabel"
	testUri := "TestUri"
	SetLabel(parent, testLabel, hl)
	SetUri(parent, testUri, hl)

	result, err := json.MarshalIndent(parent, "", "   ")
	if err != nil {
		t.Error(err)
	}

	uOfD2 := NewUniverseOfDiscourse(hl)
	recoveredParent := uOfD2.RecoverElement(result)
	if !Equivalent(parent, recoveredParent, hl) {
		t.Error("Recovered parent not equivalent to original parent")
	}
	if GetLabel(recoveredParent, hl) != testLabel {
		t.Error("Recovered test name incorrect")
	}
	if GetUri(recoveredParent, hl) != testUri {
		t.Error("Recovered test uri incorrect")
	}
	if recoveredParent.GetLabelLiteral(hl).getOwningElement(hl) != recoveredParent {
		t.Error("Recovered LabelLiteral owning element not restored properly")
	}
	if recoveredParent.GetUriLiteral(hl).getOwningElement(hl) != recoveredParent {
		t.Error("Recovered UriLiteral owning element not restored properly")
	}
}

func TestElementOwnership(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElement(hl)
	SetOwningElement(child, parent, hl)
	if GetOwningElement(child, hl) != parent {
		t.Error("Child's owner not set properly")
	}
	if child.GetOwningElementPointer(hl) == nil {
		t.Error("Child's owningElementPointer not initialized properly")
	}
	if GetOwningElement(child.GetOwningElementPointer(hl), hl) != child {
		t.Error("Child's owningElementPointer.getOwningElement() != child")
	}
	if child.GetOwningElementPointer(hl).GetElement(hl) != parent {
		t.Error("Child's owningElementPointer.getElement() != parent")
	}
	var found bool = false
	for _, be := range parent.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == child.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("Parent does not contain child in getOwnedBaseElements()")
	}
	found = false
	for _, be := range parent.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == child.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("Parent does not contain child in GetOwnedBaseElements()")
	}
}

func TestGetAbstractElementsRecursively(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	refinedElement := uOfD.NewElement(hl)
	abstractElements := uOfD.GetAbstractElementsRecursively(refinedElement, hl)
	if len(abstractElements) != 0 {
		t.Error("AbstractElements length not 0\n")
	}
	abstractElement1 := uOfD.NewElement(hl)
	refinement1 := uOfD.NewRefinement(hl)
	refinement1.SetAbstractElement(abstractElement1, hl)
	refinement1.SetRefinedElement(refinedElement, hl)
	abstractElements = uOfD.GetAbstractElementsRecursively(refinedElement, hl)
	if len(abstractElements) != 1 {
		t.Error("AbstractElements length != 1")
	}
	if abstractElements[0] != abstractElement1 {
		t.Error("AbstractElements first element not abstractElement1")
	}
	abstractElement2 := uOfD.NewElement(hl)
	refinement2 := uOfD.NewRefinement(hl)
	refinement2.SetAbstractElement(abstractElement2, hl)
	refinement2.SetRefinedElement(refinedElement, hl)
	abstractElements = uOfD.GetAbstractElementsRecursively(refinedElement, hl)
	if len(abstractElements) != 2 {
		t.Error("Abstractions length != 2")
	}
	if abstractElements[1] != abstractElement2 {
		t.Error("AbstractElements second element not abstractElement2")
	}
	abstractElement3 := uOfD.NewElement(hl)
	refinement3 := uOfD.NewRefinement(hl)
	refinement3.SetAbstractElement(abstractElement3, hl)
	refinement3.SetRefinedElement(abstractElement1, hl)
	abstractElements = uOfD.GetAbstractElementsRecursively(refinedElement, hl)
	if len(abstractElements) != 3 {
		t.Error("Abstractions length != 3")
	}
	if abstractElements[2] != abstractElement3 {
		t.Error("AbstractElements third element not abstractElement3")
	}
	abstractElement4 := uOfD.NewElement(hl)
	refinement4 := uOfD.NewRefinement(hl)
	refinement4.SetAbstractElement(abstractElement4, hl)
	refinement4.SetRefinedElement(abstractElement2, hl)
	abstractElements = uOfD.GetAbstractElementsRecursively(refinedElement, hl)
	if len(abstractElements) != 4 {
		t.Error("Abstractions length != 4")
	}
	if abstractElements[3] != abstractElement4 {
		t.Error("AbstractElements fourth element not abstractElement4")
	}
}

func TestGetImmediateAbstractElements(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	refinedElement := uOfD.NewElement(hl)
	abstractElements := uOfD.getImmediateAbstractElements(refinedElement, hl)
	if len(abstractElements) != 0 {
		t.Error("AbstractElements length not 0\n")
	}
	abstractElement1 := uOfD.NewElement(hl)
	refinement1 := uOfD.NewRefinement(hl)
	refinement1.SetAbstractElement(abstractElement1, hl)
	refinement1.SetRefinedElement(refinedElement, hl)
	abstractElements = uOfD.getImmediateAbstractElements(refinedElement, hl)
	if len(abstractElements) != 1 {
		t.Error("AbstractElements length != 1")
	}
	if abstractElements[0] != abstractElement1 {
		t.Error("AbstractElements first element not abstractElement1")
	}
	abstractElement2 := uOfD.NewElement(hl)
	refinement2 := uOfD.NewRefinement(hl)
	refinement2.SetAbstractElement(abstractElement2, hl)
	refinement2.SetRefinedElement(refinedElement, hl)
	abstractElements = uOfD.getImmediateAbstractElements(refinedElement, hl)
	if len(abstractElements) != 2 {
		t.Error("Abstractions length != 2")
	}
	if abstractElements[1] != abstractElement2 {
		t.Error("AbstractElements second element not abstractElement2")
	}
}

func TestIsRefinementOf(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	abstraction1 := uOfD.NewElement(hl)
	abstraction2 := uOfD.NewElementReference(hl)
	refinedElement := uOfD.NewElementReference(hl)
	refinement1 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement1, refinedElement, hl)
	refinement1.SetAbstractElement(abstraction1, hl)
	refinement1.SetRefinedElement(refinedElement, hl)
	refinement2 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement2, refinedElement, hl)
	refinement2.SetAbstractElement(abstraction2, hl)
	refinement2.SetRefinedElement(refinedElement, hl)

	if uOfD.IsRefinementOf(refinedElement, abstraction1, hl) == false {
		t.Errorf("Refinement of Element failed")
	}
	if uOfD.IsRefinementOf(refinedElement, abstraction2, hl) == false {
		t.Errorf("Refinement of ElementReference failed")
	}
}

func TestNewElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewElement(hl)
	if el1.GetId(hl) == "" {
		t.Error("Element identifier not properly initialized")
	}
	if el1.GetVersion(hl) != 0 {
		t.Error("Element version not properly initialized")
	}
	if el1.GetOwnedBaseElements(hl) != nil {
		t.Error("Element ownedBaseElements not properly initialized")
	}
}

func TestNewElementUriId(t *testing.T) {
	var uri string = "http://TestURI/"
	var expectedId string = uuid.NewV5(uuid.NamespaceURL, uri).String()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewElement(hl, uri)
	if expectedId != el1.GetId(hl) {
		t.Errorf("Incorrect UUID")
	}

}

func TestSetDefinition(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	var testLabel string = "Test Label"
	SetDefinition(parent, testLabel, hl)
	if parent.GetDefinition(hl) != testLabel {
		t.Error("GetLabel() value not equal to assigned name")
	}
	if parent.GetDefinitionLiteral(hl) == nil {
		t.Error("getLabelLiteral() is nil after name assigned")
	}
	if parent.GetDefinitionLiteralPointer(hl) == nil {
		t.Error("getLabelLiteralPointer() is nil after name assigned")

	}
}

func TestSetLabel(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	var testLabel string = "Test Label"
	SetLabel(parent, testLabel, hl)
	if GetLabel(parent, hl) != testLabel {
		t.Error("GetLabel() value not equal to assigned name")
	}
	if parent.GetLabelLiteral(hl) == nil {
		t.Error("getLabelLiteral() is nil after name assigned")
	}
	if parent.GetLabelLiteralPointer(hl) == nil {
		t.Error("getLabelLiteralPointer() is nil after name assigned")

	}
}

func TestSetUri(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	var testLabel string = "Test Label"
	SetUri(parent, testLabel, hl)
	if GetUri(parent, hl) != testLabel {
		t.Error("GetLabel() value not equal to assigned name")
	}
	if parent.GetUriLiteral(hl) == nil {
		t.Error("getLabelLiteral() is nil after name assigned")
	}
	if parent.GetUriLiteralPointer(hl) == nil {
		t.Error("getLabelLiteralPointer() is nil after name assigned")

	}
}

func TestVersionWithParentChange(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	oldParent := uOfD.NewElement(hl)
	newParent := uOfD.NewElement(hl)
	elementX := uOfD.NewElement(hl)
	SetOwningElement(elementX, oldParent, hl)
	oldParentInitialVersion := oldParent.GetVersion(hl)
	newParentInitialVersion := newParent.GetVersion(hl)
	elementXInitialVersion := elementX.GetVersion(hl)
	elementXOwnerPointer := elementX.GetOwningElementPointer(hl)
	elementXOwnerPointerInitialVersion := elementXOwnerPointer.GetVersion(hl)
	elementXOwnerPointer.SetElement(newParent, hl)
	if GetOwningElement(elementX, hl) != newParent {
		t.Error("elementX Owner not changed properly")
	}
	if !(elementXOwnerPointer.GetVersion(hl) > elementXOwnerPointerInitialVersion) {
		t.Error("Owning element pointer version not incremented")
	}
	if !(elementX.GetVersion(hl) > elementXInitialVersion) {
		t.Error("elementX version not incremented")
	}
	if !(oldParent.GetVersion(hl) > oldParentInitialVersion) {
		t.Error("old parent version not incremented")
	}
	if !(newParent.GetVersion(hl) > newParentInitialVersion) {
		t.Error("new parent version not incremented")
	}

}

func TestVersionWithParentChangeAndCommonGrandparent(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	grandparent := uOfD.NewElement(hl)
	grandparentPreviousVersion := grandparent.GetVersion(hl)
	oldParent := uOfD.NewElement(hl)
	oldParentPreviousVersion := oldParent.GetVersion(hl)
	SetOwningElement(oldParent, grandparent, hl)
	if !(grandparent.GetVersion(hl) > grandparentPreviousVersion) {
		t.Error("Grandparent version not incremented when old parent added as child")
	}
	grandparentPreviousVersion = grandparent.GetVersion(hl)
	if !(oldParent.GetVersion(hl) > oldParentPreviousVersion) {
		t.Error("Old parent version not incremented when added as child to grandparent")
	}
	oldParentPreviousVersion = oldParent.GetVersion(hl)

	newParent := uOfD.NewElement(hl)
	newParentPreviousVersion := newParent.GetVersion(hl)
	SetOwningElement(newParent, grandparent, hl)
	if !(grandparent.GetVersion(hl) > grandparentPreviousVersion) {
		t.Error("Grandparent version not incremented when new parent added as child")
	}
	grandparentPreviousVersion = grandparent.GetVersion(hl)
	if !(newParent.GetVersion(hl) > newParentPreviousVersion) {
		t.Error("New parent version not incremented when added as child to grandparent")
	}
	newParentPreviousVersion = newParent.GetVersion(hl)

	elementX := uOfD.NewElement(hl)
	elementXPreviousVersion := elementX.GetVersion(hl)
	SetOwningElement(elementX, oldParent, hl)
	if !(oldParent.GetVersion(hl) > oldParentPreviousVersion) {
		t.Error("Old parent version not incremented when elementX added as child")
	}
	oldParentPreviousVersion = oldParent.GetVersion(hl)
	if !(elementX.GetVersion(hl) > elementXPreviousVersion) {
		t.Error("elementX version not incremented when added as a child to oldParent")
	}
	elementXPreviousVersion = elementX.GetVersion(hl)

	elementXOwnerPointer := elementX.GetOwningElementPointer(hl)
	elementXOwnerPointerInitialVersion := elementXOwnerPointer.GetVersion(hl)
	elementXOwnerPointer.SetElement(newParent, hl)
	if GetOwningElement(elementX, hl) != newParent {
		t.Error("elementX Owner not changed properly")
	}
	if !(elementXOwnerPointer.GetVersion(hl) > elementXOwnerPointerInitialVersion) {
		t.Error("Owning element pointer version not incremented")
	}
	if !(elementX.GetVersion(hl) > elementXPreviousVersion) {
		t.Error("elementX version not incremented when parent changed")
	}
	if !(oldParent.GetVersion(hl) > oldParentPreviousVersion) {
		t.Error("old parent version not incremented when elementX removed as child")
	}
	if !(newParent.GetVersion(hl) > newParentPreviousVersion) {
		t.Error("new parent version not incremented when elementX added as child")
	}
	if !(grandparent.GetVersion(hl) > grandparentPreviousVersion) {
		t.Error("Grandparent version not incremented when elementX parent changed to new parent")
	}
}
