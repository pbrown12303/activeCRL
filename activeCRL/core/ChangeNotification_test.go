package core

import (
	//	"log"
	"sync"
	"testing"
)

var functionCalled bool
var wg sync.WaitGroup

func trialFunction(element Element, changeNotification *ChangeNotification) {
	//	PrintMutex.Lock()
	//	defer PrintMutex.Unlock()
	//	log.Printf("Entering trialFunction\n")
	//	Print(element, "+++")
	//	PrintNotification(changeNotification)
	defer wg.Done()
	functionCalled = true
}

func TestFunctionExecution(t *testing.T) {
	functionCalled = false
	uOfD := NewUniverseOfDiscourse()
	uri := "FunctionAncestor"
	GetCore().AddFunction(uri, trialFunction)
	functionAncestor := uOfD.NewElement()
	functionAncestor.SetUri(uri)
	child := uOfD.NewElement()
	refinement := uOfD.NewRefinement()
	refinement.SetAbstractElement(functionAncestor)

	// SetRefinedElement should trigger the function
	//	TraceChange = true
	wg.Add(1)
	refinement.SetRefinedElement(child)
	wg.Wait()
	TraceChange = false

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
