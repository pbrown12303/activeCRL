// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestUniverseOfDiscourseCreationTime(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	CountHeldLocksInvocations = true
	HeldLocksInvocations = 0
	NewLockInvocations = 0
	startTime := time.Now()
	uOfD := NewUniverseOfDiscourse(hl).(*universeOfDiscourse)
	if uOfD == nil {
		t.Error("Universe of Discourse not created")
	}
	endTime := time.Now()
	CountHeldLocksInvocations = false
	lockInvocations := HeldLocksInvocations
	HeldLocksInvocations = 0
	newLocks := NewLockInvocations
	NewLockInvocations = 0
	duration := endTime.Sub(startTime)
	log.Printf("Create Universe of Discourse duration: %s \n", duration.String())
	log.Printf("Create Universe of Discourse LockBaseElement invocations: %d \n", lockInvocations)
	log.Printf("Create Universe of Discourse LockBaseElement new locks: %d \n", newLocks)
}

func TestHeldLocksOverhead(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl).(*universeOfDiscourse)

	// Element
	element := uOfD.NewElement(hl)

	i := 100000
	startTime := time.Now()
	for i > 0 {
		hl.LockBaseElement(element)
		i--
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	log.Printf("RunTest duration: %s \n", duration.String())

}

func TestGetBaseElementWithUri(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl).(*universeOfDiscourse)

	// Element
	element := uOfD.NewElement(hl)
	SetLabel(element, "Element", hl)
	recoveredElement := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Element")
	if recoveredElement != nil {
		t.Error("Wrong element returned for find Element by URI")
		Print(recoveredElement, "", hl)
	}
	SetUri(element, "http://activeCrl.com/test/Element", hl)
	recoveredElement = uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Element")
	if recoveredElement == nil {
		t.Error("Did not find Element by URI")
	}

	// ElementPointer
	elementPointer := uOfD.NewReferencedElementPointer(hl)
	SetUri(elementPointer, "http://activeCrl.com/test/ElementPointer", hl)
	recoveredElementPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/ElementPointer")
	if recoveredElementPointer == nil {
		t.Error("Did not find ElementPointer by URI")
	}

	// ElementPointerPointer
	elementPointerPointer := uOfD.NewElementPointerPointer(hl)
	SetUri(elementPointerPointer, "http://activeCrl.com/test/ElementPointerPointer", hl)
	recoveredElementPointerPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/ElementPointerPointer")
	if recoveredElementPointerPointer == nil {
		t.Error("Did not find ElementPointerPointer by URI")
	}

	// ElementPointerReference
	elementPointerReference := uOfD.NewElementPointerReference(hl)
	SetLabel(elementPointerReference, "ElementReference", hl)
	SetUri(elementPointerReference, "http://activeCrl.com/test/ElementPointerReference", hl)
	recoveredElementPointerReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/ElementPointerReference")
	if recoveredElementPointerReference == nil {
		t.Error("Did not find ElementPointerReference by URI")
	}

	// ElementReference
	elementReference := uOfD.NewElementReference(hl)
	SetLabel(elementReference, "ElementReference", hl)
	SetUri(elementReference, "http://activeCrl.com/test/ElementReference", hl)
	recoveredElementReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/ElementReference")
	if recoveredElementReference == nil {
		t.Error("Did not find ElementReference by URI")
	}

	// Literal
	literal := uOfD.NewLiteral(hl)
	SetUri(literal, "http://activeCrl.com/test/Literal", hl)
	recoveredLiteral := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Literal")
	if recoveredLiteral == nil {
		t.Error("Did not find Literal by URI")
	}

	// LiteralPointer
	literalPointer := uOfD.NewValueLiteralPointer(hl)
	SetUri(literalPointer, "http://activeCrl.com/test/LiteralPointer", hl)
	recoveredLiteralPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/LiteralPointer")
	if recoveredLiteralPointer == nil {
		t.Error("Did not find LiteralPointer by URI")
	}

	// LiteralPointerPointer
	literalPointerPointer := uOfD.NewLiteralPointerPointer(hl)
	SetUri(literalPointerPointer, "http://activeCrl.com/test/LiteralPointerPointer", hl)
	recoveredLiteralPointerPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/LiteralPointerPointer")
	if recoveredLiteralPointerPointer == nil {
		t.Error("Did not find LiteralPointerPointer by URI")
	}

	// LiteralPointerReference
	literalPointerReference := uOfD.NewLiteralPointerReference(hl)
	SetLabel(literalPointerReference, "LiteralReference", hl)
	SetUri(literalPointerReference, "http://activeCrl.com/test/LiteralPointerReference", hl)
	recoveredLiteralPointerReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/LiteralPointerReference")
	if recoveredLiteralPointerReference == nil {
		t.Error("Did not find LiteralPointerReference by URI")
	}

	// LiteralReference
	literalReference := uOfD.NewLiteralReference(hl)
	SetLabel(literalReference, "LiteralReference", hl)
	SetUri(literalReference, "http://activeCrl.com/test/LiteralReference", hl)
	recoveredLiteralReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/LiteralReference")
	if recoveredLiteralReference == nil {
		t.Error("Did not find LiteralReference by URI")
	}

	// Refinement
	refinement := uOfD.NewRefinement(hl)
	SetLabel(refinement, "Refinement", hl)
	SetUri(refinement, "http://activeCrl.com/test/Refinement", hl)
	recoveredRefinement := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Refinement")
	if recoveredRefinement == nil {
		t.Error("Did not find Refinement by URI")
	}

	// Child of element
	child := uOfD.NewElement(hl)
	SetLabel(child, "Child", hl)
	SetOwningElement(child, element, hl)
	SetUri(child, "http://activeCrl.com/test/Element/Child", hl)
	recoveredChild := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Element/Child")
	if recoveredChild == nil {
		t.Error("Did not find Child by URI")
	}

}

func TestAddElementListener(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl).(*universeOfDiscourse)
	e1 := uOfD.NewElement(hl)
	ep := uOfD.NewReferencedElementPointer(hl)
	ep.SetElement(e1, hl)
	elm := uOfD.elementListenerMap.GetEntry(e1.GetId(hl))
	if elm == nil {
		t.Error("ElementListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("ElementListenerMap entry length != 1")
		} else {
			if (*elm)[0] != ep {
				t.Error("ElementListenerMap entry does not contain ElementPointer")
			}
		}
	}
	ep.SetElement(nil, hl)
	elm = uOfD.elementListenerMap.GetEntry(e1.GetId(hl))
	if elm == nil {
		t.Error("ElementListenerMap entry is nil after SetElement(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("ElementListenerMap entry length != 0")
		}
	}

}

func TestAddElementPointerListener(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl).(*universeOfDiscourse)
	ep := uOfD.NewReferencedElementPointer(hl)
	epp := uOfD.NewElementPointerPointer(hl)
	epp.SetElementPointer(ep, hl)
	elm := uOfD.elementPointerListenerMap.GetEntry(ep.GetId(hl))
	if elm == nil {
		t.Error("ElementPointerListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("ElementPointerListenerMap entry length != 1")
		} else {
			if (*elm)[0] != epp {
				t.Error("ElementPointerListenerMap entry does not contain ElementPointer")
			}
		}
	}
	epp.SetElementPointer(nil, hl)
	elm = uOfD.elementPointerListenerMap.GetEntry(ep.GetId(hl))
	if elm == nil {
		t.Error("ElementListenerMap entry is nil after SetElement(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("ElementListenerMap entry length != 0")
		}
	}

}

func TestAddLiteralListener(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl).(*universeOfDiscourse)
	e1 := uOfD.NewLiteral(hl)
	lp := uOfD.NewLabelLiteralPointer(hl)
	lp.SetLiteral(e1, hl)
	elm := uOfD.literalListenerMap.GetEntry(e1.GetId(hl))
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("LiteralListenerMap entry length != 1")
		} else {
			if (*elm)[0] != lp {
				t.Error("LiteralListenerMap entry does not contain LiteralPointer")
			}
		}
	}
	lp.SetLiteral(nil, hl)
	elm = uOfD.literalListenerMap.GetEntry(e1.GetId(hl))
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil after SetLiteral(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("LiteralListenerMap entry length != 0")
		}
	}

}

func TestAddLiteralPointerListener(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl).(*universeOfDiscourse)
	lp := uOfD.NewLabelLiteralPointer(hl)
	lpp := uOfD.NewLiteralPointerPointer(hl)
	lpp.SetLiteralPointer(lp, hl)
	elm := uOfD.literalPointerListenerMap.GetEntry(lp.GetId(hl))
	if elm == nil {
		t.Error("LiteralPointerListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("LiteralPointerListenerMap entry length != 1")
		} else {
			if (*elm)[0] != lpp {
				t.Error("LiteralPointerListenerMap entry does not contain LiteralPointer")
			}
		}
	}
	lpp.SetLiteralPointer(nil, hl)
	elm = uOfD.literalPointerListenerMap.GetEntry(lp.GetId(hl))
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil after SetLiteral(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("LiteralListenerMap entry length != 0")
		}
	}
}

// Test uOfD Notifications

type TestLock struct {
	sync.Mutex
}

var testLock TestLock
var uOfDNotificationReceived bool = false
var uOfDNotificationTestUrl string = "http://activeCRL.com/test/uOfDNotificationTestUrl"
var testElement BaseElement
var testElementFound bool = false

func uOfDNotificationReceiver(el Element, notifications []*ChangeNotification, wg *sync.WaitGroup) {
	testLock.Lock()
	defer testLock.Unlock()
	uOfDNotificationReceived = true
	for _, notification := range notifications {
		if notification.isReferenced(testElement) {
			testElementFound = true
		}
	}
}

func TestUOfDNotifications(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse(hl)

	GetCore().AddFunction(uOfDNotificationTestUrl, uOfDNotificationReceiver)
	functionRepresentation := uOfD.NewElement(hl)
	SetUri(functionRepresentation, uOfDNotificationTestUrl, hl)

	monitor := uOfD.NewElement(hl)
	reference := uOfD.NewElementReference(hl)
	SetOwningElement(reference, monitor, hl)
	refinement := uOfD.NewRefinement(hl)
	SetOwningElement(refinement, monitor, hl)
	refinement.SetAbstractElement(functionRepresentation, hl)
	refinement.SetRefinedElement(monitor, hl)
	reference.SetReferencedElement(uOfD, hl)
	hl.ReleaseLocksAndWait()

	// Check the static setup
	// Check to see that the reference's referencedElementPointer is registered to listen to the uOfD
	var foundPointer bool = false
	ep := reference.GetElementPointer(hl)
	elementListenerMap := uOfD.(*universeOfDiscourse).elementListenerMap
	for _, fep := range *(*elementListenerMap).elementPointerListMap[uOfD.getIdNoLock()] {
		if fep == ep {
			foundPointer = true
		}
	}
	if foundPointer == false {
		t.Errorf("Pointer to uOfD not in elementListenerMap")
	}

	// Now start the test
	uOfDNotificationReceived = false
	testElement = uOfD.NewElement(hl)
	hl.ReleaseLocksAndWait()

	if uOfDNotificationReceived == false {
		t.Errorf("Notification not received")
	}
	if testElementFound == false {
		t.Errorf("New Element not found in notifications, id: %s", testElement.GetId(hl).String())

	}
}
