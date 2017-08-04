package core

import (
	"testing"
)

func TestEquivalence(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElement(hl)
	SetOwningElement(child, parent, hl)
	if Equivalent(parent, parent, hl) != true {
		t.Errorf("Equivalence test failed")
	}
}

func TestGetChildElementWithAncestorUri(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
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
	if len(child.GetAbstractElementsRecursively(hl)) != 1 {
		t.Errorf("Ancestor set size != 1")
	}
	foundAncestor := child.GetAbstractElementsRecursively(hl)[0]
	if GetUri(foundAncestor, hl) != uri {
		t.Errorf("Ancestor uri not set")
	}
	foundChild = GetChildElementWithAncestorUri(parent, uri, hl)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}

func TestGetChildElementReferenceWithAncestorUri(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
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
	if len(child.GetAbstractElementsRecursively(hl)) != 1 {
		t.Errorf("Ancestor set size != 1")
	}
	foundAncestor := child.GetAbstractElementsRecursively(hl)[0]
	if GetUri(foundAncestor, hl) != uri {
		t.Errorf("Ancestor uri not set")
	}
	foundChild = GetChildElementReferenceWithAncestorUri(parent, uri, hl)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}

func TestGetChildElementWithURI(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
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
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
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
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()

	// Create original
	original := uOfD.NewElement(hl)
	SetName(original, "Root", hl)
	oChild1 := uOfD.NewElement(hl)
	SetOwningElement(oChild1, original, hl)
	oChild1Name := "Element"
	SetName(oChild1, oChild1Name, hl)
	oChild2 := uOfD.NewElementPointerReference(hl)
	SetOwningElement(oChild2, original, hl)
	oChild2Name := "ElementPointerReference"
	SetName(oChild2, oChild2Name, hl)
	oChild3 := uOfD.NewElementReference(hl)
	SetOwningElement(oChild3, original, hl)
	oChild3Name := "ElementReference"
	SetName(oChild3, oChild3Name, hl)
	oChild4 := uOfD.NewLiteralPointerReference(hl)
	SetOwningElement(oChild4, original, hl)
	oChild4Name := "LiteralPointerReference"
	SetName(oChild4, oChild4Name, hl)
	oChild5 := uOfD.NewLiteralReference(hl)
	SetOwningElement(oChild5, original, hl)
	oChild5Name := "LiteralReference"
	SetName(oChild5, oChild5Name, hl)

	replicate := uOfD.NewElement(hl)
	ReplicateAsRefinement(original, replicate, hl)

	if replicate.IsRefinementOf(original, hl) == false {
		t.Errorf("Replicate not refinement of Original")
	}

	var foundChild1Replicate bool = false
	var foundChild2Replicate bool = false
	var foundChild3Replicate bool = false
	var foundChild4Replicate bool = false
	var foundChild5Replicate bool = false

	for _, replicateChild := range replicate.GetOwnedElements(hl) {
		if replicateChild.IsRefinementOf(oChild1, hl) {
			foundChild1Replicate = true
		}
		if replicateChild.IsRefinementOf(oChild2, hl) {
			foundChild2Replicate = true
		}
		if replicateChild.IsRefinementOf(oChild3, hl) {
			foundChild3Replicate = true
		}
		if replicateChild.IsRefinementOf(oChild4, hl) {
			foundChild4Replicate = true
		}
		if replicateChild.IsRefinementOf(oChild5, hl) {
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
