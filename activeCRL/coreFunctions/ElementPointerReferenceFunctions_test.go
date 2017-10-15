// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/satori/go.uuid"
	//	"log"
	"sync"
	"testing"
	//	"time"
)

func TestElementPointerReferenceFunctionsIds(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var ElementPointerReferenceFunctionsUri string = CoreFunctionsPrefix + "ElementPointerReferenceFunctions"
	validateElementId(t, uOfD, hl, ElementPointerReferenceFunctionsUri)
	//
	//var ElementPointerReferenceCreateUri string = CoreFunctionsPrefix + "ElementPointerReference/Create"
	validateElementId(t, uOfD, hl, ElementPointerReferenceCreateUri)
	//var ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri = CoreFunctionsPrefix + "ElementPointerReference/Create/CreatedElementPointerReferenceRef"
	validateElementReferenceId(t, uOfD, hl, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri)
	//
	//var ElementPointerReferenceGetReferencedElementPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer"
	validateElementId(t, uOfD, hl, ElementPointerReferenceGetReferencedElementPointerUri)
	//var ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer/SourceElementPointerReferenceRef"
	validateElementReferenceId(t, uOfD, hl, ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri)
	//var ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer/IndicatedElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri)
	//
	//var ElementPointerReferenceGetElementPointerPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer"
	validateElementId(t, uOfD, hl, ElementPointerReferenceGetElementPointerPointerUri)
	//var ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer/SourceElementPointerReferenceRef"
	validateElementReferenceId(t, uOfD, hl, ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri)
	//var ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer/IndicatedElementPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri)
	//
	//var ElementPointerReferenceSetReferencedElementPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer"
	validateElementId(t, uOfD, hl, ElementPointerReferenceSetReferencedElementPointerUri)
	//var ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer/SourceElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri)
	//var ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer/ModifiedElementPointerReferenceRef"
	validateElementReferenceId(t, uOfD, hl, ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri)
}

func TestCreateElementPointerReferenceFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createElementPointerReference := uOfD.GetElementWithUri(ElementPointerReferenceCreateUri)
	if createElementPointerReference == nil {
		t.Error("CreateElementPointerReference not found")
	}
	createdElementPointerReferenceRef := uOfD.GetElementReferenceWithUri(ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri)
	if createdElementPointerReferenceRef == nil {
		t.Error("CreatedElementPointerReferenceRef not found")
		core.Print(createElementPointerReference, "CreateElementPointerReference: ", hl)
	}

	createElementPointerReferenceInstance := uOfD.NewElement(hl)
	createElementPointerReferenceInstanceIdentifier := createElementPointerReferenceInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createElementPointerReference, hl)

	refinementInstance.SetRefinedElement(createElementPointerReferenceInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundElementPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(createElementPointerReferenceInstance, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri, hl)
	foundElementPointerReferenceRefIdentifier := uuid.Nil
	var createdElementPointerReference core.ElementPointerReference
	createdElementPointerReferenceIdentifier := uuid.Nil
	if foundElementPointerReferenceRef == nil {
		t.Error("ElementPointerReferenceRef not created")
	} else {
		foundElementPointerReferenceRefIdentifier = foundElementPointerReferenceRef.GetId(hl)
		foundElementPointerReference := foundElementPointerReferenceRef.GetReferencedElement(hl)
		if foundElementPointerReference == nil {
			t.Error("ElementPointerReference not created")
		} else {
			switch foundElementPointerReference.(type) {
			case core.ElementPointerReference:
				createdElementPointerReference = foundElementPointerReference.(core.ElementPointerReference)
				createdElementPointerReferenceIdentifier = createdElementPointerReference.GetId(hl)
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdElementPointerReference == nil {
		t.Error("createdElementPointerReference is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdElementPointerReferenceIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created object not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createElementPointerReferenceInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundElementPointerReferenceRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createElementPointerReferenceInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReferenceReference := uOfD.GetElement(foundElementPointerReferenceRefIdentifier)
	if redoneReferenceReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri, hl) != redoneReferenceReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdElementPointerReferenceIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReferenceReference.(core.ElementReference).GetReferencedElement(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestGetElementPointerPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getElementPointerPointer := uOfD.GetElementWithUri(ElementPointerReferenceGetElementPointerPointerUri)
	if getElementPointerPointer == nil {
		t.Errorf("GetElementPointerPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getElementPointerPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetElementPointerPointer()")
	}
	sourceElementPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri, hl)
	if sourceElementPointerReferenceRef == nil {
		t.Errorf("sourceElementPointerReferenceRef child not found")
	}
	indicatedElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri, hl)
	if indicatedElementPointerPointerRef == nil {
		t.Errorf("indicatedElementPointerPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceElementPointerReference := uOfD.NewElementPointerReference(hl)
	dummyElementPointer := uOfD.NewReferencedElementPointer(hl)
	sourceElementPointerReference.SetReferencedElementPointer(dummyElementPointer, hl)
	sourceElementPointerReferenceRef.SetReferencedElement(sourceElementPointerReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetElementPointerPointer := indicatedElementPointerPointerRef.GetReferencedBaseElement(hl)
	if targetElementPointerPointer == nil {
		t.Errorf("Target ElementPointerPointer not found")
		core.Print(sourceElementPointerReferenceRef, "SourceReference: ", hl)
		core.Print(indicatedElementPointerPointerRef, "TargetReference: ", hl)
	} else {
		if targetElementPointerPointer != sourceElementPointerReference.GetElementPointerPointer(hl) {
			t.Errorf("Target ElementPointerPointer value incorrect")
		}
	}
}

func TestGetReferencedElementPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getReferencedElementPointer := uOfD.GetElementWithUri(ElementPointerReferenceGetReferencedElementPointerUri)
	if getReferencedElementPointer == nil {
		t.Errorf("GetReferencedElementPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getReferencedElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getReferencedElementPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetReferencedElementPointer()")
	}
	sourceElementPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri, hl)
	if sourceElementPointerReferenceRef == nil {
		t.Errorf("sourceElementPointerReferenceRef child not found")
	}
	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		t.Errorf("indicatedElementPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceElementPointerReference := uOfD.NewElementPointerReference(hl)
	dummyElementPointer := uOfD.NewReferencedElementPointer(hl)
	sourceElementPointerReference.SetReferencedElementPointer(dummyElementPointer, hl)
	sourceElementPointerReferenceRef.SetReferencedElement(sourceElementPointerReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	if targetElementPointer == nil {
		t.Errorf("Target element pointer not found")
		core.Print(sourceElementPointerReferenceRef, "SourceReference: ", hl)
		core.Print(indicatedElementPointerRef, "TargetReference: ", hl)
	} else {
		if targetElementPointer != sourceElementPointerReference.GetReferencedElementPointer(hl) {
			t.Errorf("Target element pointer value incorrect")
		}
	}
}

func TestSetReferencedElementPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setReferencedElementPointer := uOfD.GetElementWithUri(ElementPointerReferenceSetReferencedElementPointerUri)
	if setReferencedElementPointer == nil {
		t.Errorf("SetReferencedElementPointer function representation not found")
	} else {
		//		 Create the instance
		replicate := core.CreateReplicateAsRefinement(setReferencedElementPointer, hl)

		// Locks must be released to allow function to execute
		hl.ReleaseLocks()
		wg.Wait()

		// Now check the replication
		if uOfD.IsRefinementOf(replicate, setReferencedElementPointer, hl) != true {
			t.Errorf("Replicate is not refinement of SetReferencedElementPointer()")
		}
		sourceElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri, hl)
		if sourceElementPointerRef == nil {
			t.Errorf("SourceElementPointerRef child not found")
		}
		modifiedElementPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri, hl)
		if modifiedElementPointerReferenceRef == nil {
			t.Errorf("TargetReference child not found")
		}

		// Now test target reference update functionality
		sourceElementPointer := uOfD.NewReferencedElementPointer(hl)
		targetElementPointerReference := uOfD.NewElementPointerReference(hl)
		modifiedElementPointerReferenceRef.SetReferencedElement(targetElementPointerReference, hl)
		sourceElementPointerRef.SetReferencedElementPointer(sourceElementPointer, hl)

		// Locks must be released to allow function to execute
		hl.ReleaseLocks()
		wg.Wait()

		hl.LockBaseElement(replicate)
		if targetElementPointerReference.GetReferencedElementPointer(hl) != sourceElementPointer {
			t.Errorf("ElementPointerReference value not set")
			core.Print(sourceElementPointerRef, "ElementPointerRef: ", hl)
			core.Print(modifiedElementPointerReferenceRef, "TargetReference: ", hl)
		}
	}
}