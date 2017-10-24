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

func TestNewBaseElementReference(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewBaseElementReference(hl)
	if el1.GetId(hl) == uuid.Nil {
		t.Error("Element identifier not properly initialized")
	}
	if el1.GetVersion(hl) != 0 {
		t.Error("Element version not properly initialized")
	}
	if el1.GetOwnedBaseElements(hl) != nil {
		t.Error("Element ownedBaseElements not properly initialized")
	}
}

func TestNewBaseElementReferenceUriId(t *testing.T) {
	var uri string = "http://TestURI/"
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewBaseElementReference(hl, uri)
	var expectedId uuid.UUID = uuid.NewV5(uuid.NamespaceURL, uri)
	if expectedId != el1.GetId(hl) {
		t.Errorf("Incorrect UUID")
	}

}

func TestBaseElementReferenceOwnership(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewBaseElementReference(hl)
	SetOwningElement(child, parent, hl)
	//	Print(parent, "", hl)
	//	Print(child, "", hl)
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

func TestReferenceSetBaseElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewBaseElementReference(hl)
	SetOwningElement(child, parent, hl)
	if child.GetReferencedBaseElement(hl) != nil {
		t.Error("BaseElementReference's base element not initialized to nil")
	}
	child.SetReferencedBaseElement(parent, hl)
	if child.GetReferencedBaseElement(hl) == nil {
		t.Error("BaseElementReference's base element is nil after assignment")
		Print(parent, "   ", hl)
	}
	if child.GetReferencedBaseElement(hl) != nil && child.GetReferencedBaseElement(hl).GetId(hl) != parent.GetId(hl) {
		t.Error("BaseElementReference's base element not set properly")
	}
	child.SetReferencedBaseElement(nil, hl)
	if child.GetReferencedBaseElement(hl) != nil {
		t.Error("BaseElementReference's base element not nild properly")
	}
}

func TestBaseElementReferenceMarshal(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewBaseElementReference(hl)
	SetOwningElement(child, parent, hl)
	child.SetReferencedBaseElement(parent, hl)

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

func TestBaseElementReferenceClone(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewBaseElementReference(hl)
	SetOwningElement(child, parent, hl)
	child.SetReferencedBaseElement(parent, hl)
	clone := child.(*baseElementReference).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("ElementReference clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}
}
