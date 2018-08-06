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

func TestNewElementPointerReference(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewElementPointerReference(hl)
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

func TestNewElementPointerReferenceUriId(t *testing.T) {
	var uri string = "http://TestURI/"
	var expectedId string = uuid.NewV5(uuid.NamespaceURL, uri).String()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el1 := uOfD.NewElementPointerReference(hl, uri)
	if expectedId != el1.GetId(hl) {
		t.Errorf("Incorrect UUID")
	}
}

func TestElementPointerReferenceOwnership(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementPointerReference(hl)
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

func TestSetReferencedElementPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementPointerReference(hl)
	SetOwningElement(child, parent, hl)
	if child.GetReferencedElementPointer(hl) != nil {
		t.Error("ElementPointerReference's element pointer not initialized to nil")
	}
	elementPointer := uOfD.NewReferencedElementPointer(hl)
	child.SetReferencedElementPointer(elementPointer, hl)
	if child.GetReferencedElementPointer(hl) == nil {
		t.Error("ElementPointerReference's  element pointer is nil after assignment")
		Print(elementPointer, "   ", hl)
	}
	if child.GetReferencedElementPointer(hl) != nil && child.GetReferencedElementPointer(hl).GetId(hl) != elementPointer.GetId(hl) {
		t.Error("ElementPointerReference's  element pointer not set properly")
	}
	child.SetReferencedElementPointer(nil, hl)
	if child.GetReferencedElementPointer(hl) != nil {
		t.Error("ElementPointerReference's  element pointer not nild properly")
	}
}

func TestElementPointerReferenceMarshal(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementPointerReference(hl)
	SetOwningElement(child, parent, hl)
	elementPointer := uOfD.NewReferencedElementPointer(hl)
	child.SetReferencedElementPointer(elementPointer, hl)

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

func TestElementPointerReferenceClone(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	parent := uOfD.NewElement(hl)
	child := uOfD.NewElementPointerReference(hl)
	SetOwningElement(child, parent, hl)
	elementPointer := uOfD.NewReferencedElementPointer(hl)
	child.SetReferencedElementPointer(elementPointer, hl)

	clone := child.(*elementPointerReference).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("Element clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}

}
