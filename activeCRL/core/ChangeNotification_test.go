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

func TestFunctionExecution(t *testing.T) {
	functionCalled = false
	uOfD := NewUniverseOfDiscourse()
	uri := "FunctionAncestor"
	GetCore().computeFunctions[uri] = trialFunction
	functionAncestor := uOfD.NewElement()
	functionAncestor.SetUri(uri)
	child := uOfD.NewElement()
	refinement := uOfD.NewRefinement()
	refinement.SetAbstractElement(functionAncestor)

	// SetRefinedElement should trigger the function
	wg.Add(1)
	refinement.SetRefinedElement(child)
	wg.Wait()

	if functionCalled == false {
		t.Errorf("TrialFunction not called after abstraction created")
	}

	// Now test to see if SetName() also triggers the function
	// The SetName() call is going to result in six change notification function calls
	wg.Add(6)
	functionCalled = false
	child.SetName("Child")
	wg.Wait()

	if functionCalled == false {
		t.Errorf("TrialFunction not called after child.SetName()")
	}
}
