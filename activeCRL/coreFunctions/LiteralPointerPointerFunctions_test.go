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

func TestLiteralPointerPointerFunctionsIds(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var LiteralPointerPointerFunctionsUri string = CoreFunctionsPrefix + "LiteralPointerPointerFunctions"
	validateElementId(t, uOfD, hl, LiteralPointerPointerFunctionsUri)
	//
	//var LiteralPointerPointerCreateLiteralPointerPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/CreateAbstractLiteralPointerPointer"
	validateElementId(t, uOfD, hl, LiteralPointerPointerCreateLiteralPointerPointerUri)
	//var LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri = CoreFunctionsPrefix + "LiteralPointerPointer/CreateAbstractLiteralPointerPointer/CreatedLiteralPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri)
	//
	//var LiteralPointerPointerGetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer"
	validateElementId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerUri)
	//var LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer/SourceLiteralPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri)
	//var LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer/IndicatedLiteralPointerRef"
	validateLiteralPointerReferenceId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri)
	//
	//var LiteralPointerPointerGetLiteralPointerIdUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId"
	validateElementId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerIdUri)
	//var LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId/SourceLiteralPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri)
	//var LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri)
	//
	//var LiteralPointerPointerGetLiteralPointerVersionUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion"
	validateElementId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerVersionUri)
	//var LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion/SourceLiteralPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri)
	//var LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri)
	//
	//var LiteralPointerPointerSetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer"
	validateElementId(t, uOfD, hl, LiteralPointerPointerSetLiteralPointerUri)
	//var LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer/LiteralPointerRef"
	validateLiteralPointerReferenceId(t, uOfD, hl, LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri)
	//var LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer/ModifiedLiteralPointerPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri)
}

func TestCreateLiteralPointerPointerFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createLiteralPointerPointer := uOfD.GetElementWithUri(LiteralPointerPointerCreateLiteralPointerPointerUri)
	if createLiteralPointerPointer == nil {
		t.Error("CreateLiteralPointerPointer not found")
	}
	createdLiteralPointerPointerRef := uOfD.GetBaseElementReferenceWithUri(LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri)
	if createdLiteralPointerPointerRef == nil {
		t.Error("CreatedLiteralPointerPointerRef not found")
	}

	// Now create the instance of the function
	createElementPointePointerInstance := uOfD.NewElement(hl)
	createElementPointePointerInstanceIdentifier := createElementPointePointerInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createLiteralPointerPointer, hl)
	refinementInstance.SetRefinedElement(createElementPointePointerInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundReference := core.GetChildBaseElementReferenceWithAncestorUri(createElementPointePointerInstance, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri, hl)
	foundReferenceIdentifier := uuid.Nil
	var createdLiteralPointerPointer core.LiteralPointerPointer
	createdLiteralPointerPointerIdentifier := uuid.Nil
	if foundReference == nil {
		t.Error("Reference not created")
	} else {
		foundReferenceIdentifier = foundReference.GetId(hl)
		foundBaseElement := foundReference.GetReferencedBaseElement(hl)
		if foundBaseElement == nil {
			t.Error("LiteralPointerPointer not created")
		} else {
			switch foundBaseElement.(type) {
			case core.LiteralPointerPointer:
				createdLiteralPointerPointer = foundBaseElement.(core.LiteralPointerPointer)
				createdLiteralPointerPointerIdentifier = createdLiteralPointerPointer.GetId(hl)
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdLiteralPointerPointer == nil {
		t.Error("createdLiteralPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdLiteralPointerPointerIdentifier)
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
	if uOfD.GetElement(createdLiteralPointerPointerIdentifier) != nil {
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
		if core.GetChildBaseElementReferenceWithAncestorUri(redoneInstance, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdLiteralPointerPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.BaseElementReference).GetReferencedBaseElement(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestGetLiteralPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getLiteralPointer := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerUri)
	if getLiteralPointer == nil {
		t.Errorf("GetElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getLiteralPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetLiteralPointer()")
	}
	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		t.Errorf("SourceLiteralPointerPointerRef child not found")
	}
	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		t.Errorf("IndicatedLiteralPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceLiteralPointer := uOfD.NewNameLiteralPointer(hl)
	sourceLiteralPointerPointer := uOfD.NewLiteralPointerPointer(hl)
	sourceLiteralPointerPointer.SetLiteralPointer(sourceLiteralPointer, hl)
	sourceLiteralPointerPointerRef.SetReferencedBaseElement(sourceLiteralPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)

	indicatedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	if indicatedLiteralPointer != sourceLiteralPointer {
		t.Errorf("Target element value incorrect")
	}
}

func TestGetLiteralPointerId(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getLiteralPointerId := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerIdUri)
	if getLiteralPointerId == nil {
		t.Errorf("GetLiteralPointerId function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getLiteralPointerId, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getLiteralPointerId, hl) != true {
		t.Errorf("Replicate is not refinement of GetLiteralPointerId()")
	}
	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		t.Errorf("SourceReference child not found")
	}
	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri, hl)
	if createdLiteralRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceLiteralPointer := uOfD.NewNameLiteralPointer(hl)
	sourceLiteralPointerPointer := uOfD.NewLiteralPointerPointer(hl)
	sourceLiteralPointerPointer.SetLiteralPointer(sourceLiteralPointer, hl)
	sourceLiteralPointerPointerRef.SetReferencedBaseElement(sourceLiteralPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != sourceLiteralPointerPointer.GetLiteralPointerId(hl).String() {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestGetLiteralPointerVersion(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getLiteralPointerVersion := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerVersionUri)
	if getLiteralPointerVersion == nil {
		t.Errorf("GetLiteralPointerVersion function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getLiteralPointerVersion, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getLiteralPointerVersion, hl) != true {
		t.Errorf("Replicate is not refinement of GetLiteralPointerVersion()")
	}
	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		t.Errorf("SourceLiteralPointerPointerRef child not found")
	}
	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		t.Errorf("CreatedLiteralRef child not found")
	}

	// Now test target reference update functionality
	sourceLiteralPointer := uOfD.NewNameLiteralPointer(hl)
	// Force the version to change
	core.SetUri(sourceLiteralPointer, "Test URI", hl)
	sourceLiteralPointerPointer := uOfD.NewLiteralPointerPointer(hl)
	sourceLiteralPointerPointer.SetLiteralPointer(sourceLiteralPointer, hl)
	sourceLiteralPointerPointerRef.SetReferencedBaseElement(sourceLiteralPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != strconv.Itoa(sourceLiteralPointerPointer.GetLiteralPointerVersion(hl)) {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestSetLiteralPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setLiteralPointer := uOfD.GetElementWithUri(LiteralPointerPointerSetLiteralPointerUri)
	if setLiteralPointer == nil {
		t.Errorf("SetLiteralPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setLiteralPointer, hl) != true {
		t.Errorf("Replicate is not refinement of SetLiteralPointer()")
	}
	literalPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri, hl)
	if literalPointerRef == nil {
		t.Errorf("ElementReference child not found")
	}
	modifiedLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri, hl)
	if modifiedLiteralPointerPointerRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceLiteralPointer := uOfD.NewNameLiteralPointer(hl)
	literalPointerRef.SetReferencedLiteralPointer(sourceLiteralPointer, hl)

	targetLiteralPointerPointer := uOfD.NewLiteralPointerPointer(hl)
	modifiedLiteralPointerPointerRef.SetReferencedBaseElement(targetLiteralPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if targetLiteralPointerPointer.GetLiteralPointer(hl) != sourceLiteralPointer {
		t.Errorf("LiteralPointerPointer value not set")
		core.Print(literalPointerRef, "LiteralPointerRef: ", hl)
		core.Print(modifiedLiteralPointerPointerRef, "ModifiedLiteralPointerPointerRef: ", hl)
	}
}
