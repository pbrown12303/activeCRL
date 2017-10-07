// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	//	"fmt"
	"sync"
	"testing"
)

func TestNewLiteralPointerPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	parent := uOfD.NewLiteralPointerReference(hl)
	child := uOfD.NewLiteralPointerPointer(hl)
	SetOwningElement(child, parent, hl)
	if GetOwningElement(child, hl) != parent {
		t.Error("Child's owner not set properly")
	}
	var found bool = false
	for _, be := range parent.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == child.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("LiteralPointerPointer not found in parent's OwnedBaseElements \n")
	}
	if parent.GetLiteralPointerPointer(hl) != child {
		t.Error("LiteralPointerReference.GetLiteralPointer() did not return child")
	}
}

func TestSetLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	parent := uOfD.NewLiteralPointerReference(hl)
	child := uOfD.NewLiteralPointerPointer(hl)
	SetOwningElement(child, parent, hl)
	literalPointer := uOfD.NewValueLiteralPointer(hl)
	SetOwningElement(literalPointer, parent, hl)
	if child.GetLiteralPointer(hl) != nil {
		t.Error("LiteralPointer's Literal not initially nil \n")
	}
	child.SetLiteralPointer(literalPointer, hl)
	if child.GetLiteralPointer(hl) != literalPointer {
		t.Error("LiteralPointerPointer's LiteralPointer not properly set after assignment \n")
	}
}

func TestLiteralPointerPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	parent := uOfD.NewLiteralPointerReference(hl)
	child := uOfD.NewLiteralPointerPointer(hl)
	SetOwningElement(child, parent, hl)
	literalPointer := uOfD.NewValueLiteralPointer(hl)
	SetOwningElement(literalPointer, parent, hl)
	child.SetLiteralPointer(literalPointer, hl)

	result, err := json.MarshalIndent(parent, "", "   ")
	if err != nil {
		t.Error(err)
	}

	//	fmt.Printf("Encoded Parent \n%s \n", result)

	uOfD2 := NewUniverseOfDiscourse()
	recoveredParent := uOfD2.RecoverElement(result)
	if !Equivalent(parent, recoveredParent, hl) {
		t.Error("Recovered parent not equivalent to original parent")
	}
}

func TestLiteralPointerPointerClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	parent := uOfD.NewLiteralPointerReference(hl)
	child := uOfD.NewLiteralPointerPointer(hl)
	SetOwningElement(child, parent, hl)
	literalPointer := uOfD.NewValueLiteralPointer(hl)
	SetOwningElement(literalPointer, parent, hl)
	child.SetLiteralPointer(literalPointer, hl)
	clone := child.(*literalPointerPointer).clone()
	if !Equivalent(child, clone, hl) {
		t.Error("LiteralPointerPointer clone failed")
		Print(child, "   ", hl)
		Print(clone, "   ", hl)
	}

}
