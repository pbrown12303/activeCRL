package core

import (
	//	"log"
	"sync"
	"testing"
)

var functionCalled bool
var wg sync.WaitGroup

func trialFunction(element Element, changeNotification *ChangeNotification) {
	defer wg.Done()
	functionCalled = true
}

func TesFunctionExecution(t *testing.T) {
	hl := NewHeldLocks()
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
	wg.Add(1)
	refinement.SetRefinedElement(child, hl)
	wg.Wait()
	TraceChange = false

	if functionCalled == false {
		t.Errorf("TrialFunction not called after abstraction created")
	}

	// Now test to see if SetName() also triggers the function
	// The SetName() call is going to result in six change notification function calls
	wg.Add(6)
	functionCalled = false
	SetName(child, "Child", hl)
	wg.Wait()

	if functionCalled == false {
		t.Errorf("TrialFunction not called after child.SetName()")
	}
}
