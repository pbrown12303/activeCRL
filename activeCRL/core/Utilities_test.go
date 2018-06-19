// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
	"testing"
)

func TestEquivalence(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElement(hl)
	SetOwningElement(child, parent, hl)
	if Equivalent(parent, parent, hl) != true {
		t.Errorf("Equivalence test failed")
	}
}

func TestGetChildElementWithAncestorUri(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElement(hl)
	SetOwningElement(child, parent, hl)
	ancestor := uOfD.NewElement(hl)
	refinement := uOfD.NewRefinement(hl)
	refinement.SetAbstractElement(ancestor, hl)
	refinement.SetRefinedElement(child, hl)
	uri := "testingUri"
	foundChild := GetChildElementWithAncestorUri(parent, uri, hl)
	if foundChild != nil {
		t.Errorf("Child found when it should not have been")
	}
	SetUri(ancestor, uri, hl)
	if len(uOfD.GetAbstractElementsRecursively(child, hl)) != 1 {
		t.Errorf("Ancestor set size != 1")
	}
	foundAncestor := uOfD.GetAbstractElementsRecursively(child, hl)[0]
	if GetUri(foundAncestor, hl) != uri {
		t.Errorf("Ancestor uri not set")
	}
	foundChild = GetChildElementWithAncestorUri(parent, uri, hl)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}

func TestGetChildElementReferenceWithAncestorUri(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementReference(hl)
	SetOwningElement(child, parent, hl)
	ancestor := uOfD.NewElementReference(hl)
	refinement := uOfD.NewRefinement(hl)
	refinement.SetAbstractElement(ancestor, hl)
	refinement.SetRefinedElement(child, hl)
	uri := "testingUri"
	foundChild := GetChildElementReferenceWithAncestorUri(parent, uri, hl)
	if foundChild != nil {
		t.Errorf("Child found when it should not have been")
	}
	SetUri(ancestor, uri, hl)
	if len(uOfD.GetAbstractElementsRecursively(child, hl)) != 1 {
		t.Errorf("Ancestor set size != 1")
	}
	foundAncestor := uOfD.GetAbstractElementsRecursively(child, hl)[0]
	if GetUri(foundAncestor, hl) != uri {
		t.Errorf("Ancestor uri not set")
	}
	foundChild = GetChildElementReferenceWithAncestorUri(parent, uri, hl)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}

func TestGetChildElementWithURI(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElement(hl)
	SetOwningElement(child, parent, hl)
	uri := "testingUri"
	foundChild := GetChildElementWithUri(parent, uri, hl)
	if foundChild != nil {
		t.Errorf("Child found when it should not have been")
	}
	SetUri(child, uri, hl)
	foundChild = GetChildElementWithUri(parent, uri, hl)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}

func TestGetChildElementReferenceWithURI(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementReference(hl)
	SetOwningElement(child, parent, hl)
	uri := "testingUri"
	foundChild := GetChildElementReferenceWithUri(parent, uri, hl)
	if foundChild != nil {
		t.Errorf("Child found when it should not have been")
	}
	SetUri(child, uri, hl)
	foundChild = GetChildElementReferenceWithUri(parent, uri, hl)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}

func TestReplicateAsRefinement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)

	// Create original
	original := uOfD.NewElement(hl)
	SetLabel(original, "Root", hl)
	oChild1 := uOfD.NewElement(hl)
	SetOwningElement(oChild1, original, hl)
	oChild1Label := "Element"
	SetLabel(oChild1, oChild1Label, hl)
	oChild2 := uOfD.NewElementPointerReference(hl)
	SetOwningElement(oChild2, original, hl)
	oChild2Label := "ElementPointerReference"
	SetLabel(oChild2, oChild2Label, hl)
	oChild3 := uOfD.NewElementReference(hl)
	SetOwningElement(oChild3, original, hl)
	oChild3Label := "ElementReference"
	SetLabel(oChild3, oChild3Label, hl)
	oChild4 := uOfD.NewLiteralPointerReference(hl)
	SetOwningElement(oChild4, original, hl)
	oChild4Label := "LiteralPointerReference"
	SetLabel(oChild4, oChild4Label, hl)
	oChild5 := uOfD.NewLiteralReference(hl)
	SetOwningElement(oChild5, original, hl)
	oChild5Label := "LiteralReference"
	SetLabel(oChild5, oChild5Label, hl)

	replicate := uOfD.NewElement(hl)
	ReplicateAsRefinement(original, replicate, hl)

	if uOfD.IsRefinementOf(replicate, original, hl) == false {
		t.Errorf("Replicate not refinement of Original")
	}

	var foundChild1Replicate bool = false
	var foundChild2Replicate bool = false
	var foundChild3Replicate bool = false
	var foundChild4Replicate bool = false
	var foundChild5Replicate bool = false

	for _, replicateChild := range replicate.GetOwnedElements(hl) {
		if uOfD.IsRefinementOf(replicateChild, oChild1, hl) {
			foundChild1Replicate = true
		}
		if uOfD.IsRefinementOf(replicateChild, oChild2, hl) {
			foundChild2Replicate = true
		}
		if uOfD.IsRefinementOf(replicateChild, oChild3, hl) {
			foundChild3Replicate = true
		}
		if uOfD.IsRefinementOf(replicateChild, oChild4, hl) {
			foundChild4Replicate = true
		}
		if uOfD.IsRefinementOf(replicateChild, oChild5, hl) {
			foundChild5Replicate = true
		}
	}

	if foundChild1Replicate == false {
		t.Errorf("Child1 Replicate not found")
	}
	if foundChild2Replicate == false {
		t.Errorf("Child2 Replicate not found")
	}
	if foundChild3Replicate == false {
		t.Errorf("Child3 Replicate not found")
	}
	if foundChild4Replicate == false {
		t.Errorf("Child4 Replicate not found")
	}
	if foundChild5Replicate == false {
		t.Errorf("Child5 Replicate not found")
	}

	childCount := len(replicate.GetOwnedBaseElements(hl))

	// Now test to make sure children are not duplicated
	ReplicateAsRefinement(original, replicate, hl)
	if len(replicate.GetOwnedBaseElements(hl)) != childCount {
		t.Errorf("ReplicateAsRefinement is not idempotent")
	}
}
