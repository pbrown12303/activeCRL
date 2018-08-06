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

func TestNewLiteralReference(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewLiteralReference(hl)
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

func TestNewLiteralReferenceUriId(t *testing.T) {
	var uri string = "http://TestURI/"
	var expectedId string = uuid.NewV5(uuid.NamespaceURL, uri).String()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewLiteralReference(hl, uri)
	if expectedId != el1.GetId(hl) {
		t.Errorf("Incorrect UUID")
	}
}

func TestLiteralReferenceOwnership(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteralReference(hl)
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

func TestSetReferencedLiteral(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteralReference(hl)
	SetOwningElement(child, parent, hl)
	literal := uOfD.NewLiteral(hl)
	if child.GetReferencedLiteral(hl) != nil {
		t.Error("LiteralReference's referenced literal not initialized to nil")
	}
	child.SetReferencedLiteral(literal, hl)
	if child.GetReferencedLiteral(hl) == nil {
		t.Error("LiteralReference's referenced literal is nil after assignment")
		Print(parent, "   ", hl)
	}
	if child.GetReferencedLiteral(hl) != nil && child.GetReferencedLiteral(hl).GetId(hl) != literal.GetId(hl) {
		t.Error("LiteralReference's referenced literal not set properly")
	}
	child.SetReferencedLiteral(nil, hl)
	if child.GetReferencedLiteral(hl) != nil {
		t.Error("LiteralReference's referenced literal not nild properly")
	}
}

func TestLiteralReferenceMarshal(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteralReference(hl)
	SetOwningElement(child, parent, hl)
	literal := uOfD.NewLiteral(hl)
	child.SetReferencedLiteral(literal, hl)
	//	fmt.Printf("Parent before encoding \n")
	//	Print(parent, "   ")

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

func TestLiteralReferenceClone(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewLiteralReference(hl)
	SetOwningElement(child, parent, hl)
	literal := uOfD.NewLiteral(hl)
	child.SetReferencedLiteral(literal, hl)
	clone := child.(*literalReference).clone()
	if !Equivalent(child, clone, hl) {
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
		t.Error("Cloned LiteralReference not equivalent to original")
	}
}
