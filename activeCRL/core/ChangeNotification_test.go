// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	//	"log"
	"sync"
	"testing"
)

var functionCalled bool
var wg sync.WaitGroup

func trialFunction(element Element, changeNotifications []*ChangeNotification, wg *sync.WaitGroup) {
	defer wg.Done()
	functionCalled = true
}

func TesFunctionExecution(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	functionCalled = false
	uOfD := NewUniverseOfDiscourse()
	uri := "FunctionAncestor"
	GetCore().AddFunction(uri, trialFunction)
	functionAncestor := uOfD.NewElement(hl)
	SetUri(functionAncestor, uri, hl)
	child := uOfD.NewElement(hl)
	refinement := uOfD.NewRefinement(hl)
	refinement.SetAbstractElement(functionAncestor, hl)

	// SetRefinedElement should trigger the function
	//	TraceChange = true
	//	wg.Add(1)
	refinement.SetRefinedElement(child, hl)
	wg.Wait()
	//	TraceChange = false

	if functionCalled == false {
		t.Errorf("TrialFunction not called after abstraction created")
	}

	// Now test to see if SetName() also triggers the function
	// The SetName() call is going to result in six change notification function calls
	//	wg.Add(6)
	functionCalled = false
	SetName(child, "Child", hl)
	wg.Wait()

	if functionCalled == false {
		t.Errorf("TrialFunction not called after child.SetName()")
	}
}
