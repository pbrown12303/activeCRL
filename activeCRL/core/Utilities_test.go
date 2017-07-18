package core

import (
	"testing"
)

func TestEquivalence(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElement()
	child.SetOwningElement(parent)
	if Equivalent(parent, parent) != true {
		t.Errorf("Equivalence test failed")
	}
}

func TestGetChildWithAncestorUri(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElement()
	child.SetOwningElement(parent)
	ancestor := uOfD.NewElement()
	refinement := uOfD.NewRefinement()
	refinement.SetAbstractElement(ancestor)
	refinement.SetRefinedElement(child)
	uri := "testingUri"
	foundChild := GetChildElementWithAncestorUri(parent, uri)
	if foundChild != nil {
		t.Errorf("Child found when it should not have been")
	}
	ancestor.SetUri(uri)
	if len(child.GetAbstractElementsRecursivelyNoLock()) != 1 {
		t.Errorf("Ancestor set size != 1")
	}
	foundAncestor := child.GetAbstractElementsRecursivelyNoLock()[0]
	if foundAncestor.GetUri() != uri {
		t.Errorf("Ancestor uri not set")
	}
	foundChild = GetChildElementWithAncestorUri(parent, uri)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}

func TestGetChildElementWithURI(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElement()
	child.SetOwningElement(parent)
	uri := "testingUri"
	foundChild := GetChildElementWithUri(parent, uri)
	if foundChild != nil {
		t.Errorf("Child found when it should not have been")
	}
	child.SetUri(uri)
	foundChild = GetChildElementWithUri(parent, uri)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}

func TestGetChildElementReferenceWithURI(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElementReference()
	child.SetOwningElement(parent)
	uri := "testingUri"
	foundChild := GetChildElementReferenceWithUri(parent, uri)
	if foundChild != nil {
		t.Errorf("Child found when it should not have been")
	}
	child.SetUri(uri)
	foundChild = GetChildElementReferenceWithUri(parent, uri)
	if foundChild == nil {
		t.Errorf("Child not found")
	}
}
