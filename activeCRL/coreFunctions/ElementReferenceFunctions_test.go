// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"sync"
	"testing"
	//	"time"
)

func TestElementReferenceFunctionsIds(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var ElementReferenceFunctionsUri string = CoreFunctionsPrefix + "ElementReferenceFunctions"
	validateElementId(t, uOfD, hl, ElementReferenceFunctionsUri)
	//
	//var ElementReferenceCreateUri string = CoreFunctionsPrefix + "ElementReference/Create"
	validateElementId(t, uOfD, hl, ElementReferenceCreateUri)
	//var ElementReferenceCreateCreatedElementReferenceRefUri = CoreFunctionsPrefix + "ElementReference/Create/CreatedElementReferenceRef"
	validateElementReferenceId(t, uOfD, hl, ElementReferenceCreateCreatedElementReferenceRefUri)
	//
	//var ElementReferenceGetReferencedElementUri string = CoreFunctionsPrefix + "ElementReference/GetReferencedElement"
	validateElementId(t, uOfD, hl, ElementReferenceGetReferencedElementUri)
	//var ElementReferenceGetReferencedElementSourceElementReferenceRefUri = CoreFunctionsPrefix + "ElementReference/GetReferencedElement/SourceElementReferenceRef"
	validateElementReferenceId(t, uOfD, hl, ElementReferenceGetReferencedElementSourceElementReferenceRefUri)
	//var ElementReferenceGetReferencedElementIndicatedElementRefUri string = CoreFunctionsPrefix + "ElementReference/GetReferencedElement/IndicatedElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementReferenceGetReferencedElementIndicatedElementRefUri)
	//
	//var ElementReferenceGetElementPointerUri string = CoreFunctionsPrefix + "ElementReference/GetElementPointer"
	validateElementId(t, uOfD, hl, ElementReferenceGetElementPointerUri)
	//var ElementReferenceGetElementPointerSourceElementReferenceRefUri string = CoreFunctionsPrefix + "ElementReference/GetElementPointer/SourceElementReferenceRef"
	validateElementReferenceId(t, uOfD, hl, ElementReferenceGetElementPointerSourceElementReferenceRefUri)
	//var ElementReferenceGetElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "ElementReference/GetElementPointer/IndicatedElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementReferenceGetElementPointerIndicatedElementPointerRefUri)
	//
	//var ElementReferenceSetReferencedElementUri string = CoreFunctionsPrefix + "ElementReference/SetReferencedElement"
	validateElementId(t, uOfD, hl, ElementReferenceSetReferencedElementUri)
	//var ElementReferenceSetReferencedElementSourceElementRefUri string = CoreFunctionsPrefix + "ElementReference/SetReferencedElement/SourceElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementReferenceSetReferencedElementSourceElementRefUri)
	//var ElementReferenceSetReferencedElementModifiedElementReferenceRefUri string = CoreFunctionsPrefix + "ElementReference/SetReferencedElement/ModifiedElementReferenceRef"
	validateElementReferenceId(t, uOfD, hl, ElementReferenceSetReferencedElementModifiedElementReferenceRefUri)
}

func TestCreateElementReferenceFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createElementReference := uOfD.GetElementWithUri(ElementReferenceCreateUri)
	if createElementReference == nil {
		t.Error("CreateElementReference not found")
	}
	createdElementReferenceRef := uOfD.GetElementReferenceWithUri(ElementReferenceCreateCreatedElementReferenceRefUri)
	if createdElementReferenceRef == nil {
		t.Error("CreatedElementReferenceRef not found")
		core.Print(createElementReference, "CreateElementReference: ", hl)
	}

	createElementReferenceInstance := uOfD.NewElement(hl)
	createElementReferenceInstanceIdentifier := createElementReferenceInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createElementReference, hl)

	refinementInstance.SetRefinedElement(createElementReferenceInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundElementReferenceRef := core.GetChildElementReferenceWithAncestorUri(createElementReferenceInstance, ElementReferenceCreateCreatedElementReferenceRefUri, hl)
	foundElementReferenceRefIdentifier := ""
	var createdElementReference core.ElementReference
	createdElementReferenceIdentifier := ""
	if foundElementReferenceRef == nil {
		t.Error("ElementReferenceRef not created")
	} else {
		foundElementReferenceRefIdentifier = foundElementReferenceRef.GetId(hl).String()
		foundElementReference := foundElementReferenceRef.GetReferencedElement(hl)
		if foundElementReference == nil {
			t.Error("ElementReference not created")
		} else {
			switch foundElementReference.(type) {
			case core.ElementReference:
				createdElementReference = foundElementReference.(core.ElementReference)
				createdElementReferenceIdentifier = createdElementReference.GetId(hl).String()
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdElementReference == nil {
		t.Error("createdElementReference is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdElementReferenceIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created object not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createElementReferenceInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundElementReferenceRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createElementReferenceInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReferenceReference := uOfD.GetElement(foundElementReferenceRefIdentifier)
	if redoneReferenceReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, ElementReferenceCreateCreatedElementReferenceRefUri, hl) != redoneReferenceReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdElementReferenceIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReferenceReference.(core.ElementReference).GetReferencedElement(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestElementReferenceGetElementPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getElementPointer := uOfD.GetElementWithUri(ElementReferenceGetElementPointerUri)
	if getElementPointer == nil {
		t.Errorf("GetElementPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getElementPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetElementPointer()")
	}
	sourceElementReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceGetElementPointerSourceElementReferenceRefUri, hl)
	if sourceElementReferenceRef == nil {
		t.Errorf("sourceElementReferenceRef child not found")
	}
	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementReferenceGetElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		t.Errorf("indicatedElementPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceElementReference := uOfD.NewElementReference(hl)
	dummyElement := uOfD.NewElement(hl)
	sourceElementReference.SetReferencedElement(dummyElement, hl)
	sourceElementReferenceRef.SetReferencedElement(sourceElementReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetReferencedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	if targetReferencedElementPointer == nil {
		t.Errorf("Target ReferencedElementPointer not found")
		core.Print(sourceElementReferenceRef, "SourceReference: ", hl)
		core.Print(indicatedElementPointerRef, "TargetReference: ", hl)
	} else {
		if targetReferencedElementPointer != sourceElementReference.GetElementPointer(hl) {
			t.Errorf("Target ElementPointer value incorrect")
		}
	}
}

func TestGetReferencedElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getReferencedElement := uOfD.GetElementWithUri(ElementReferenceGetReferencedElementUri)
	if getReferencedElement == nil {
		t.Errorf("GetReferencedElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getReferencedElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getReferencedElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetReferencedElement()")
	}
	sourceElementReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceGetReferencedElementSourceElementReferenceRefUri, hl)
	if sourceElementReferenceRef == nil {
		t.Errorf("sourceElementReferenceRef child not found")
	}
	indicatedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceGetReferencedElementIndicatedElementRefUri, hl)
	if indicatedElementRef == nil {
		t.Errorf("indicatedElementRef child not found")
	}

	// Now test target reference update functionality
	sourceElementReference := uOfD.NewElementReference(hl)
	dummyElement := uOfD.NewElement(hl)
	sourceElementReference.SetReferencedElement(dummyElement, hl)
	sourceElementReferenceRef.SetReferencedElement(sourceElementReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetElement := indicatedElementRef.GetReferencedElement(hl)
	if targetElement == nil {
		t.Errorf("Target element pointer not found")
		core.Print(sourceElementReferenceRef, "SourceReference: ", hl)
		core.Print(indicatedElementRef, "TargetReference: ", hl)
	} else {
		if targetElement != sourceElementReference.GetReferencedElement(hl) {
			t.Errorf("Target element pointer value incorrect")
		}
	}
}

func TestSetReferencedElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setReferencedElement := uOfD.GetElementWithUri(ElementReferenceSetReferencedElementUri)
	if setReferencedElement == nil {
		t.Errorf("SetReferencedElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setReferencedElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setReferencedElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetReferencedElement()")
	}
	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceSetReferencedElementSourceElementRefUri, hl)
	if sourceElementRef == nil {
		t.Errorf("SourceElementRef child not found")
	}
	modifiedElementReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceSetReferencedElementModifiedElementReferenceRefUri, hl)
	if modifiedElementReferenceRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElement := uOfD.NewElement(hl)
	targetElementReference := uOfD.NewElementReference(hl)
	modifiedElementReferenceRef.SetReferencedElement(targetElementReference, hl)
	sourceElementRef.SetReferencedElement(sourceElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if targetElementReference.GetReferencedElement(hl) != sourceElement {
		t.Errorf("ElementReference value not set")
		core.Print(sourceElementRef, "ElementRef: ", hl)
		core.Print(modifiedElementReferenceRef, "TargetReference: ", hl)
	}
}
