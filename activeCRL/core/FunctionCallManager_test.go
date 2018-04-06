// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
	"testing"
)

var function1CallCount int

func trialFunction1(element Element, changeNotifications []*ChangeNotification, wg *sync.WaitGroup) {
	function1CallCount++
}

var function2CallCount int

func trialFunction2(element Element, changeNotifications []*ChangeNotification, wg *sync.WaitGroup) {
	function2CallCount++
}

func TestFunctionCallManager(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)
	el := uOfD.NewElement(hl)
	fcm := hl.functionCallManager

	// Add the first functioncall
	var functionIdentifier1 crlExecutionFunctionIdentifier = "Function1"
	GetCore().AddFunction(string(functionIdentifier1), trialFunction1)
	cn1a := NewChangeNotification(el, ADD, "TestFunctionCallManager", nil)
	fcm.AddFunctionCall(functionIdentifier1, el, cn1a)
	if fcm.functionTargetMap[functionIdentifier1] == nil {
		t.Errorf("Labeled function not found in function target map\n")
		fcm.Print("Pending Function Calls: ")
	} else {
		enm := fcm.functionTargetMap[functionIdentifier1]
		if len(enm[el]) != 1 {
			t.Errorf("ElementNotificationsMap length != 1")
		}
		if enm[el][0] != cn1a {
			t.Errorf("ElementNotificationsMap[0] != cn1a")
		}
	}

	// Now call the first function with a different change notification
	cn1b := NewChangeNotification(el, MODIFY, "TestFunctionCallManager", nil)
	fcm.AddFunctionCall(functionIdentifier1, el, cn1b)
	if fcm.functionTargetMap[functionIdentifier1] == nil {
		t.Errorf("Labeled function not found in function target map\n")
		fcm.Print("Pending Function Calls: ")
	} else {
		enm := fcm.functionTargetMap[functionIdentifier1]
		if len(enm[el]) != 2 {
			t.Errorf("ElementNotificationsMap length != 2")
			fcm.Print("")
		}
		if enm[el][1] != cn1b {
			t.Errorf("ElementNotificationsMap[1] != cn1b")
		}
	}

	// Add a call to the second function
	var functionIdentifier2 crlExecutionFunctionIdentifier = "Function2"
	GetCore().AddFunction(string(functionIdentifier2), trialFunction2)
	cn2a := NewChangeNotification(el, ADD, "TestFunctionCallManager", nil)
	fcm.AddFunctionCall(functionIdentifier2, el, cn2a)
	if fcm.functionTargetMap[functionIdentifier2] == nil {
		t.Errorf("Labeled function 2 not found in function target map\n")
		fcm.Print("Pending Function Calls: ")
	} else {
		enm := fcm.functionTargetMap[functionIdentifier2]
		if len(enm[el]) != 1 {
			t.Errorf("ElementNotificationsMap length != 1")
		}
		if enm[el][0] != cn2a {
			t.Errorf("ElementNotificationsMap[0] != cn2a")
		}
	}

	// Now call the second function again with a different change notification
	cn2b := NewChangeNotification(el, MODIFY, "TestFunctionCallManager", nil)
	fcm.AddFunctionCall(functionIdentifier2, el, cn2b)
	if fcm.functionTargetMap[functionIdentifier2] == nil {
		t.Errorf("Labeled function not found in function target map\n")
		fcm.Print("Pending Function Calls: ")
	} else {
		enm := fcm.functionTargetMap[functionIdentifier2]
		if len(enm[el]) != 2 {
			t.Errorf("ElementNotificationsMap length != 2")
			fcm.Print("")
		}
		if enm[el][1] != cn2b {
			t.Errorf("ElementNotificationsMap[1] != cn2b")
		}
	}
	hl.ReleaseLocks()
	hl.waitGroup.Wait()

	if function1CallCount != 1 {
		t.Errorf("Function 1 call count incorrect: %d", function1CallCount)
	}
	if function2CallCount != 1 {
		t.Errorf("Function 2 call count incorrect: %d", function2CallCount)
	}

}
