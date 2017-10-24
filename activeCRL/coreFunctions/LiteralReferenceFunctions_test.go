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

func TestLiteralReferenceFunctionsIds(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var LiteralReferenceFunctionsUri string = CoreFunctionsPrefix + "LiteralReferenceFunctions"
	validateElementId(t, uOfD, hl, LiteralReferenceFunctionsUri)
	//
	//var LiteralReferenceCreateUri string = CoreFunctionsPrefix + "LiteralReference/Create"
	validateElementId(t, uOfD, hl, LiteralReferenceCreateUri)
	//var LiteralReferenceCreateCreatedLiteralReferenceRefUri = CoreFunctionsPrefix + "LiteralReference/Create/CreatedLiteralReferenceRef"
	validateElementReferenceId(t, uOfD, hl, LiteralReferenceCreateCreatedLiteralReferenceRefUri)
	//
	//var LiteralReferenceGetReferencedLiteralUri string = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral"
	validateElementId(t, uOfD, hl, LiteralReferenceGetReferencedLiteralUri)
	//var LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral/SourceLiteralReferenceRef"
	validateElementReferenceId(t, uOfD, hl, LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri)
	//var LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral/IndicatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri)
	//
	//var LiteralReferenceGetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer"
	validateElementId(t, uOfD, hl, LiteralReferenceGetLiteralPointerUri)
	//var LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer/SourceLiteralReferenceRef"
	validateElementReferenceId(t, uOfD, hl, LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri)
	//var LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer/IndicatedLiteralPointerRef"
	validateLiteralPointerReferenceId(t, uOfD, hl, LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri)
	//
	//var LiteralReferenceSetReferencedLiteralUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral"
	validateElementId(t, uOfD, hl, LiteralReferenceSetReferencedLiteralUri)
	//var LiteralReferenceSetReferencedLiteralSourceLiteralRefUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral/SourceLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralReferenceSetReferencedLiteralSourceLiteralRefUri)
	//var LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral/ModifiedLiteralReferenceRef"
	validateElementReferenceId(t, uOfD, hl, LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri)
}

func TestCreateLiteralReferenceFunction(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createLiteralReference := uOfD.GetElementWithUri(LiteralReferenceCreateUri)
	if createLiteralReference == nil {
		t.Error("CreateLiteralReference not found")
	}
	createdLiteralReferenceRef := uOfD.GetElementReferenceWithUri(LiteralReferenceCreateCreatedLiteralReferenceRefUri)
	if createdLiteralReferenceRef == nil {
		t.Error("CreatedLiteralReferenceRef not found")
		core.Print(createLiteralReference, "CreateLiteralReference: ", hl)
	}

	createLiteralReferenceInstance := uOfD.NewElement(hl)
	createLiteralReferenceInstanceIdentifier := createLiteralReferenceInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createLiteralReference, hl)

	refinementInstance.SetRefinedElement(createLiteralReferenceInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(createLiteralReferenceInstance, LiteralReferenceCreateCreatedLiteralReferenceRefUri, hl)
	foundLiteralReferenceRefIdentifier := uuid.Nil
	var createdLiteralReference core.LiteralReference
	createdLiteralReferenceIdentifier := uuid.Nil
	if foundLiteralReferenceRef == nil {
		t.Error("LiteralReferenceRef not created")
	} else {
		foundLiteralReferenceRefIdentifier = foundLiteralReferenceRef.GetId(hl)
		foundLiteralReference := foundLiteralReferenceRef.GetReferencedElement(hl)
		if foundLiteralReference == nil {
			t.Error("LiteralReference not created")
		} else {
			switch foundLiteralReference.(type) {
			case core.LiteralReference:
				createdLiteralReference = foundLiteralReference.(core.LiteralReference)
				createdLiteralReferenceIdentifier = createdLiteralReference.GetId(hl)
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdLiteralReference == nil {
		t.Error("createdLiteralReference is nil")
	}
	newlyCreatedLiteral := uOfD.GetBaseElement(createdLiteralReferenceIdentifier)
	if newlyCreatedLiteral == nil {
		t.Error("Created object not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createLiteralReferenceInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundLiteralReferenceRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createLiteralReferenceInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneLiteralReferenceRef := uOfD.GetElement(foundLiteralReferenceRefIdentifier)
	if redoneLiteralReferenceRef == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, LiteralReferenceCreateCreatedLiteralReferenceRefUri, hl) != redoneLiteralReferenceRef {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedLiteralReference := uOfD.GetBaseElement(createdLiteralReferenceIdentifier)
		if redoneCreatedLiteralReference == nil {
			t.Error("Created literal not redone")
		} else {
			if redoneLiteralReferenceRef.(core.ElementReference).GetReferencedElement(hl) != redoneCreatedLiteralReference {
				t.Error("Reference pointer to created literal not restored")
			}
		}
	}
}

func TestLiteralReferenceGetLiteralPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getLiteralPointer := uOfD.GetElementWithUri(LiteralReferenceGetLiteralPointerUri)
	if getLiteralPointer == nil {
		t.Errorf("GetLiteralPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getLiteralPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetLiteralPointer()")
	}
	sourceLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri, hl)
	if sourceLiteralReferenceRef == nil {
		t.Errorf("sourceLiteralReferenceRef child not found")
	}
	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		t.Errorf("indicatedLiteralPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceLiteralReference := uOfD.NewLiteralReference(hl)
	dummyLiteral := uOfD.NewLiteral(hl)
	sourceLiteralReference.SetReferencedLiteral(dummyLiteral, hl)
	sourceLiteralReferenceRef.SetReferencedElement(sourceLiteralReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetReferencedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	if targetReferencedLiteralPointer == nil {
		t.Errorf("Target ReferencedElementPointer not found")
		core.Print(sourceLiteralReferenceRef, "SourceReference: ", hl)
		core.Print(indicatedLiteralPointerRef, "TargetReference: ", hl)
	} else {
		if targetReferencedLiteralPointer != sourceLiteralReference.GetLiteralPointer(hl) {
			t.Errorf("Target ElementPointer value incorrect")
		}
	}
}

func TestGetReferencedLiteral(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getReferencedLiteral := uOfD.GetElementWithUri(LiteralReferenceGetReferencedLiteralUri)
	if getReferencedLiteral == nil {
		t.Errorf("GetReferencedLiteral function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getReferencedLiteral, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getReferencedLiteral, hl) != true {
		t.Errorf("Replicate is not refinement of GetReferencedLiteral()")
	}
	sourceLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri, hl)
	if sourceLiteralReferenceRef == nil {
		t.Errorf("sourceLiteralReferenceRef child not found")
	}
	indicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		t.Errorf("indicatedLiteralRef child not found")
	}

	// Now test target reference update functionality
	sourceLiteralReference := uOfD.NewLiteralReference(hl)
	dummyLiteral := uOfD.NewLiteral(hl)
	sourceLiteralReference.SetReferencedLiteral(dummyLiteral, hl)
	sourceLiteralReferenceRef.SetReferencedElement(sourceLiteralReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := indicatedLiteralRef.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target element pointer not found")
		core.Print(sourceLiteralReferenceRef, "SourceReference: ", hl)
		core.Print(indicatedLiteralRef, "TargetReference: ", hl)
	} else {
		if targetLiteral != sourceLiteralReference.GetReferencedLiteral(hl) {
			t.Errorf("Target element pointer value incorrect")
		}
	}
}

func TestSetReferencedLiteral(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setReferencedLiteral := uOfD.GetElementWithUri(LiteralReferenceSetReferencedLiteralUri)
	if setReferencedLiteral == nil {
		t.Errorf("SetReferencedLiteral function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setReferencedLiteral, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, setReferencedLiteral, hl) != true {
		t.Errorf("Replicate is not refinement of SetReferencedLiteral()")
	}
	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralReferenceSetReferencedLiteralSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		t.Errorf("SourceLiteralRef child not found")
	}
	modifiedLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri, hl)
	if modifiedLiteralReferenceRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	targetLiteralReference := uOfD.NewLiteralReference(hl)
	modifiedLiteralReferenceRef.SetReferencedElement(targetLiteralReference, hl)
	sourceLiteralRef.SetReferencedLiteral(sourceLiteral, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if targetLiteralReference.GetReferencedLiteral(hl) != sourceLiteral {
		t.Errorf("LiteralReference value not set")
		core.Print(sourceLiteralRef, "ElementRef: ", hl)
		core.Print(modifiedLiteralReferenceRef, "TargetReference: ", hl)
	}
}
