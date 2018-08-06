// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	//	"fmt"
	"github.com/satori/go.uuid"
	"sync"
	"testing"
)

func TestNewRefinement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewRefinement(hl)
	if el1.GetId(hl) == "" {
		t.Error("Refinement identifier not properly initialized")
	}
	if el1.GetVersion(hl) != 0 {
		t.Error("Refinement version not properly initialized")
	}
	if el1.GetOwnedBaseElements(hl) != nil {
		t.Error("Refinement ownedBaseElements not properly initialized")
	}
}

func TestNewRefinementUriId(t *testing.T) {
	var uri string = "http://TestURI/"
	var expectedId string = uuid.NewV5(uuid.NamespaceURL, uri).String()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewRefinement(hl, uri)
	if expectedId != el1.GetId(hl) {
		t.Errorf("Incorrect UUID")
	}
}

func TestRefinementOwnership(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	if GetOwningElement(child, hl) != parent {
		t.Error("Child's owner not set properly")
	}
	if child.GetOwningElementPointer(hl) == nil {
		t.Error("Child's owningElementPointer not initialized properly")
	}
	if GetOwningElement(child.GetOwningElementPointer(hl), hl).GetId(hl) != child.GetId(hl) {
		t.Error("Child's owningElementPointer.getOwningElement() != child")
	}
	if child.GetOwningElementPointer(hl).GetElement(hl) != parent {
		t.Error("Child's owningElementPointer.getElement() != parent")
	}
}

func TestSetAbstractElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	abstractElement := uOfD.NewElement(hl)
	if child.GetAbstractElement(hl) != nil {
		t.Error("Refinement's abstract element not initialized to nil")
	}
	child.SetAbstractElement(abstractElement, hl)
	if child.GetAbstractElement(hl) == nil {
		t.Error("Refinement's abstract element is nil after assignment")
		Print(parent, "   ", hl)
	}
	if child.GetAbstractElement(hl) != nil && child.GetAbstractElement(hl).GetId(hl) != abstractElement.GetId(hl) {
		t.Error("Refinement's abstract element not set properly")
	}
	child.SetAbstractElement(nil, hl)
	if child.GetAbstractElement(hl) != nil {
		t.Error("Refinement's abstract element not nild properly")
	}
}

func TestSetRefinedElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	refinedElement := uOfD.NewElement(hl)
	if child.GetRefinedElement(hl) != nil {
		t.Error("Refinement's refined element not initialized to nil")
	}
	child.SetRefinedElement(refinedElement, hl)
	if child.GetRefinedElement(hl) == nil {
		t.Error("Refinement's refined element is nil after assignment")
		Print(parent, "   ", hl)
	}
	if child.GetRefinedElement(hl) != nil && child.GetRefinedElement(hl).GetId(hl) != refinedElement.GetId(hl) {
		t.Error("Refinement's refined element not set properly")
	}
	child.SetRefinedElement(nil, hl)
	if child.GetRefinedElement(hl) != nil {
		t.Error("Refinement's refined element not nild properly")
	}
}

func TestRefinementMarshal(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	abstractElement := uOfD.NewElement(hl)
	child.SetAbstractElement(abstractElement, hl)
	refinedElement := uOfD.NewElement(hl)
	child.SetRefinedElement(refinedElement, hl)

	result, err := json.MarshalIndent(parent, "", "   ")
	if err != nil {
		t.Error(err)
	}

	//	fmt.Printf("Encoded Parent \n%s \n", result)

	uOfD2 := NewUniverseOfDiscourse(hl)
	recoveredParent := uOfD2.RecoverElement(result)
	if recoveredParent != nil {
		//		Print(recoveredParent, "")
	}
	if !Equivalent(parent, recoveredParent, hl) {
		t.Error("Recovered parent not equivalent to original parent")
	}
}

func TestRefinementClone(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewRefinement(hl)
	SetOwningElement(child, parent, hl)
	abstractElement := uOfD.NewElement(hl)
	child.SetAbstractElement(abstractElement, hl)
	refinedElement := uOfD.NewElement(hl)
	child.SetRefinedElement(refinedElement, hl)
	clone := child.(*refinement).clone()
	if !Equivalent(child, clone, hl) {
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
		t.Error("Cloned Refinement not equivalent to original")
	}

}
