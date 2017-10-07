// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"sync"
	"testing"
)

func TestBaseElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	owner := uOfD.NewBaseElementReference(hl)
	rep := uOfD.NewBaseElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	if GetOwningElement(rep, hl) != owner {
		t.Error("Base element pointer's owner not set properly")
	}
	var found bool = false
	for _, be := range owner.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == rep.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("Base Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if owner.GetBaseElementPointer(hl) != rep {
		t.Error("Owner.getBaseElementPointer() did not return Base Element Pointer")
	}
}

func TestSetBaseElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	owner := uOfD.NewBaseElementReference(hl)
	rep := uOfD.NewBaseElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	target := uOfD.NewElement(hl)
	SetOwningElement(target, owner, hl)
	rep.SetBaseElement(target, hl)
	if rep.GetBaseElement(hl) != target {
		t.Error("BaseElementPointer SetBaseElement(target) did not work")
	}
	rep.SetBaseElement(nil, hl)
	if rep.GetBaseElement(hl) != nil {
		t.Error("BaseElementPointer SetBaseElement(nil) did not work")
	}
}

func TestBaseElementPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	owner := uOfD.NewBaseElementReference(hl)
	rep := uOfD.NewBaseElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	target := uOfD.NewElement(hl)
	SetOwningElement(target, owner, hl)
	rep.SetBaseElement(target, hl)

	result, err := json.MarshalIndent(owner, "", "   ")
	if err != nil {
		t.Error(err)
	}

	uOfD2 := NewUniverseOfDiscourse()
	recoveredOwner := uOfD2.RecoverElement(result)
	if !Equivalent(owner, recoveredOwner, hl) {
		t.Error("Recovered owner not equivalent to original owner")
	}
}

func TestBaseElementPointerClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	owner := uOfD.NewBaseElementReference(hl)
	rep := uOfD.NewBaseElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	target := uOfD.NewElement(hl)
	SetOwningElement(target, owner, hl)
	rep.SetBaseElement(target, hl)

	clone := rep.(*baseElementPointer).clone()
	if !Equivalent(rep, clone, hl) {
		t.Error("ElementPointer clone failed")
		Print(rep, "   ", hl)
		Print(clone, "   ", hl)
	}

}
