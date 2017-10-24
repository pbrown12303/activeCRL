// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/satori/go.uuid"
	//	"log"
	"strconv"
	"sync"
	"testing"
	//	"time"
)

func TestElementPointerPointerFunctionsIds(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var ElementPointerPointerFunctionsUri string = CoreFunctionsPrefix + "ElementPointerPointerFunctions"
	validateElementId(t, uOfD, hl, ElementPointerPointerFunctionsUri)
	//
	//var ElementPointerPointerCreateElementPointerPointerUri string = CoreFunctionsPrefix + "ElementPointerPointer/CreateAbstractElementPointerPointer"
	validateElementId(t, uOfD, hl, ElementPointerPointerCreateElementPointerPointerUri)
	//var ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri = CoreFunctionsPrefix + "ElementPointerPointer/CreateAbstractElementPointerPointer/CreatedElementPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri)
	//
	//var ElementPointerPointerGetElementPointerUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointer"
	validateElementId(t, uOfD, hl, ElementPointerPointerGetElementPointerUri)
	//var ElementPointerPointerGetElementPointerSourceElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointer/SourceElementPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, ElementPointerPointerGetElementPointerSourceElementPointerPointerRefUri)
	//var ElementPointerPointerGetElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointer/IndicatedElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerPointerGetElementPointerIndicatedElementPointerRefUri)
	//
	//var ElementPointerPointerGetElementPointerIdUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerId"
	validateElementId(t, uOfD, hl, ElementPointerPointerGetElementPointerIdUri)
	//var ElementPointerPointerGetElementPointerIdSourceElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerId/SourceElementPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, ElementPointerPointerGetElementPointerIdSourceElementPointerPointerRefUri)
	//var ElementPointerPointerGetElementPointerIdCreatedLiteralUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerId/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementPointerPointerGetElementPointerIdCreatedLiteralUri)
	//
	//var ElementPointerPointerGetElementPointerVersionUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerVersion"
	validateElementId(t, uOfD, hl, ElementPointerPointerGetElementPointerVersionUri)
	//var ElementPointerPointerGetElementPointerVersionSourceElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerVersion/SourceElementPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, ElementPointerPointerGetElementPointerVersionSourceElementPointerPointerRefUri)
	//var ElementPointerPointerGetElementPointerVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerVersion/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementPointerPointerGetElementPointerVersionCreatedLiteralRefUri)
	//
	//var ElementPointerPointerSetElementPointerUri string = CoreFunctionsPrefix + "ElementPointerPointer/SetElementPointer"
	validateElementId(t, uOfD, hl, ElementPointerPointerSetElementPointerUri)
	//var ElementPointerPointerSetElementPointerElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/SetElementPointer/ElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerPointerSetElementPointerElementPointerRefUri)
	//var ElementPointerPointerSetElementPointerModifiedElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/SetElementPointer/ModifiedElementPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, ElementPointerPointerSetElementPointerModifiedElementPointerPointerRefUri)
}

func TestCreateElementPointerPointerFunction(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createElementPointerPointer := uOfD.GetElementWithUri(ElementPointerPointerCreateElementPointerPointerUri)
	if createElementPointerPointer == nil {
		t.Error("CreateElementPointerPointer not found")
	}
	createdElementPointerPointerRef := uOfD.GetBaseElementReferenceWithUri(ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri)
	if createdElementPointerPointerRef == nil {
		t.Error("CreatedElementPointerPointerRef not found")
	}

	// Now create the instance of the function
	createElementPointePointerInstance := uOfD.NewElement(hl)
	createElementPointePointerInstanceIdentifier := createElementPointePointerInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createElementPointerPointer, hl)
	refinementInstance.SetRefinedElement(createElementPointePointerInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundReference := core.GetChildBaseElementReferenceWithAncestorUri(createElementPointePointerInstance, ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri, hl)
	foundReferenceIdentifier := uuid.Nil
	var createdElementPointerPointer core.ElementPointerPointer
	createdElementPointerPointerIdentifier := uuid.Nil
	if foundReference == nil {
		t.Error("Reference not created")
	} else {
		foundReferenceIdentifier = foundReference.GetId(hl)
		foundBaseElement := foundReference.GetReferencedBaseElement(hl)
		if foundBaseElement == nil {
			t.Error("ElementPointerPointer not created")
		} else {
			switch foundBaseElement.(type) {
			case core.ElementPointerPointer:
				createdElementPointerPointer = foundBaseElement.(core.ElementPointerPointer)
				createdElementPointerPointerIdentifier = createdElementPointerPointer.GetId(hl)
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdElementPointerPointer == nil {
		t.Error("createdElementPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdElementPointerPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createElementPointePointerInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundReferenceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdElementPointerPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createElementPointePointerInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundReferenceIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildBaseElementReferenceWithAncestorUri(redoneInstance, ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdElementPointerPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.BaseElementReference).GetReferencedBaseElement(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestGetElementPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getElementPointer := uOfD.GetElementWithUri(ElementPointerPointerGetElementPointerUri)
	if getElementPointer == nil {
		t.Errorf("GetElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getElementPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetElementPointer()")
	}
	sourceElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerSourceElementPointerPointerRefUri, hl)
	if sourceElementPointerPointerRef == nil {
		t.Errorf("SourceElementPointerPointerRef child not found")
	}
	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		t.Errorf("IndicatedElementPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceElementPointer := uOfD.NewReferencedElementPointer(hl)
	sourceElementPointerPointer := uOfD.NewElementPointerPointer(hl)
	sourceElementPointerPointer.SetElementPointer(sourceElementPointer, hl)
	sourceElementPointerPointerRef.SetReferencedBaseElement(sourceElementPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)

	indicatedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	if indicatedElementPointer != sourceElementPointer {
		t.Errorf("Target element value incorrect")
	}
}

func TestGetElementPointerId(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getElementPointerId := uOfD.GetElementWithUri(ElementPointerPointerGetElementPointerIdUri)
	if getElementPointerId == nil {
		t.Errorf("GetElementPointerId function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementPointerId, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getElementPointerId, hl) != true {
		t.Errorf("Replicate is not refinement of GetElementPointerId()")
	}
	sourceElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerIdSourceElementPointerPointerRefUri, hl)
	if sourceElementPointerPointerRef == nil {
		t.Errorf("SourceReference child not found")
	}
	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerIdCreatedLiteralUri, hl)
	if createdLiteralRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElementPointer := uOfD.NewReferencedElementPointer(hl)
	sourceElementPointerPointer := uOfD.NewElementPointerPointer(hl)
	sourceElementPointerPointer.SetElementPointer(sourceElementPointer, hl)
	sourceElementPointerPointerRef.SetReferencedBaseElement(sourceElementPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != sourceElementPointerPointer.GetElementPointerId(hl).String() {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestGetElementPointerVersion(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getElementPointerVersion := uOfD.GetElementWithUri(ElementPointerPointerGetElementPointerVersionUri)
	if getElementPointerVersion == nil {
		t.Errorf("GetElementPointerVersion function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementPointerVersion, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getElementPointerVersion, hl) != true {
		t.Errorf("Replicate is not refinement of GetElementPointerVersion()")
	}
	sourceElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerVersionSourceElementPointerPointerRefUri, hl)
	if sourceElementPointerPointerRef == nil {
		t.Errorf("SourceElementPointerPointerRef child not found")
	}
	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerVersionCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		t.Errorf("CreatedLiteralRef child not found")
	}

	// Now test target reference update functionality
	sourceElementPointer := uOfD.NewReferencedElementPointer(hl)
	// Force the version to change
	core.SetUri(sourceElementPointer, "Test URI", hl)
	sourceElementPointerPointer := uOfD.NewElementPointerPointer(hl)
	sourceElementPointerPointer.SetElementPointer(sourceElementPointer, hl)
	sourceElementPointerPointerRef.SetReferencedBaseElement(sourceElementPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != strconv.Itoa(sourceElementPointerPointer.GetElementPointerVersion(hl)) {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestSetElementPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setElementPointer := uOfD.GetElementWithUri(ElementPointerPointerSetElementPointerUri)
	if setElementPointer == nil {
		t.Errorf("SetElementPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, setElementPointer, hl) != true {
		t.Errorf("Replicate is not refinement of SetElementPointer()")
	}
	elementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerPointerSetElementPointerElementPointerRefUri, hl)
	if elementPointerRef == nil {
		t.Errorf("ElementReference child not found")
	}
	modifiedElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerPointerSetElementPointerModifiedElementPointerPointerRefUri, hl)
	if modifiedElementPointerPointerRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElementPointer := uOfD.NewReferencedElementPointer(hl)
	elementPointerRef.SetReferencedElementPointer(sourceElementPointer, hl)

	targetElementPointerPointer := uOfD.NewElementPointerPointer(hl)
	modifiedElementPointerPointerRef.SetReferencedBaseElement(targetElementPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if targetElementPointerPointer.GetElementPointer(hl) != sourceElementPointer {
		t.Errorf("ElementPointerPointer value not set")
		core.Print(elementPointerRef, "ElementPointerRef: ", hl)
		core.Print(modifiedElementPointerPointerRef, "ModifiedElementPointerPointerRef: ", hl)
	}
}
