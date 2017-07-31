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
