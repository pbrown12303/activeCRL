// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"sync"
	"testing"
)

func TestNewElementPointerPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	owner := uOfD.NewElementPointerReference(hl)
	epp := uOfD.NewElementPointerPointer(hl)
	SetOwningElement(epp, owner, hl)
	if GetOwningElement(epp, hl) != owner {
		t.Error("Element pointer pointer's owner not set properly")
	}
	var found bool = false
	for _, be := range owner.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == epp.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("Element Pointer Pointer not found in parent's OwnedBaseElements \n")
	}
	if owner.GetElementPointerPointer(hl) != epp {
		t.Error("Owner.getElementPointerPointer() did not return Referenced Element Pointer")
	}
}

func TestNewElementPointerPointerUriId(t *testing.T) {
	var uri string = "http://TestURI/"
	var expectedId string = uuid.NewV5(uuid.NamespaceURL, uri).String()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	epp := uOfD.NewElementPointerPointer(hl, uri)
	if expectedId != epp.GetId(hl) {
		t.Errorf("Incorrect UUID")
	}
}

func TestSetElementPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	owner := uOfD.NewElementPointerReference(hl)
	epp := uOfD.NewElementPointerPointer(hl)
	SetOwningElement(epp, owner, hl)
	target := uOfD.NewReferencedElementPointer(hl)
	SetOwningElement(target, owner, hl)
	epp.SetElementPointer(target, hl)
	if epp.GetElementPointer(hl) != target {
		t.Error("ElementPointer not set to target properly")
	}
	epp.SetElementPointer(nil, hl)
	if epp.GetElementPointer(hl) != nil {
		t.Error("ElementPointer not set to nil properly")
	}
}

func TestElementPointerPointerMarshal(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	owner := uOfD.NewElementPointerReference(hl)
	epp := uOfD.NewElementPointerPointer(hl)
	SetOwningElement(epp, owner, hl)
	target := uOfD.NewReferencedElementPointer(hl)
	SetOwningElement(target, owner, hl)
	epp.SetElementPointer(target, hl)

	result, err := json.MarshalIndent(owner, "", "   ")
	if err != nil {
		t.Error(err)
	}

	uOfD2 := NewUniverseOfDiscourse(hl)
	recoveredOwner := uOfD2.RecoverElement(result)
	if !Equivalent(owner, recoveredOwner, hl) {
		t.Error("Recovered owner not equivalent to original owner")
	}
}

func TestElementPointerPointerClone(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	owner := uOfD.NewElementPointerReference(hl)
	epp := uOfD.NewElementPointerPointer(hl)
	SetOwningElement(epp, owner, hl)
	target := uOfD.NewReferencedElementPointer(hl)
	SetOwningElement(target, owner, hl)
	epp.SetElementPointer(target, hl)

	clone := epp.(*elementPointerPointer).clone()
	if !Equivalent(epp, clone, hl) {
		t.Error("ElementPointerPointer clone failed")
		Print(epp, "   ", hl)
		Print(clone, "   ", hl)
	}

}
