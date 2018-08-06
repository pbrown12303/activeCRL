// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
	"testing"
)

func TestBaseElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	var be baseElement
	// Test id
	// Can't use the locking version of GetId because we are about to change the id
	if be.getIdNoLock() != "" {
		t.Error("baseElement identifier not nil before setting")
	}

	// This will change the id
	be.initializeBaseElement()
	if be.GetId(hl) == "" {
		t.Error("baseElement identifier nil after setting")
	}

	// Test version
	if be.GetVersion(hl) != 0 {
		t.Error("baseElement version not 0 before increment")
	}
	be.internalIncrementVersion()
	if be.GetVersion(hl) != 1 {
		t.Error("baseElement version not 1 after increment")
	}

}

func TestGetUri(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	e1 := uOfD.NewElement(hl)
	testUri := "testUri"
	SetUri(e1, testUri, hl)
	if GetUri(e1, hl) != testUri {
		t.Error("Uri not returned")
	}
}
