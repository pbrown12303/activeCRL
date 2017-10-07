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

func TestCreateLiteralPointerReferenceFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createLiteralPointerReference := uOfD.GetElementWithUri(LiteralPointerReferenceCreateUri)
	if createLiteralPointerReference == nil {
		t.Error("CreateLiteralPointerReference not found")
	}
	createdLiteralPointerReferenceRef := uOfD.GetElementReferenceWithUri(LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri)
	if createdLiteralPointerReferenceRef == nil {
		t.Error("CreatedLiteralPointerReferenceRef not found")
		core.Print(createLiteralPointerReference, "CreateLiteralPointerReference: ", hl)
	}

	createLiteralPointerReferenceInstance := uOfD.NewElement(hl)
	createLiteralPointerReferenceInstanceIdentifier := createLiteralPointerReferenceInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createLiteralPointerReference, hl)

	refinementInstance.SetRefinedElement(createLiteralPointerReferenceInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundLiteralPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(createLiteralPointerReferenceInstance, LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri, hl)
	foundLiteralPointerReferenceRefIdentifier := ""
	var createdLiteralPointerReference core.LiteralPointerReference
	createdLiteralPointerReferenceIdentifier := ""
	if foundLiteralPointerReferenceRef == nil {
		t.Error("LiteralPointerReferenceRef not created")
	} else {
		foundLiteralPointerReferenceRefIdentifier = foundLiteralPointerReferenceRef.GetId(hl).String()
		foundLiteralPointerReference := foundLiteralPointerReferenceRef.GetReferencedElement(hl)
		if foundLiteralPointerReference == nil {
			t.Error("LiteralPointerReference not created")
		} else {
			switch foundLiteralPointerReference.(type) {
			case core.LiteralPointerReference:
				createdLiteralPointerReference = foundLiteralPointerReference.(core.LiteralPointerReference)
				createdLiteralPointerReferenceIdentifier = createdLiteralPointerReference.GetId(hl).String()
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdLiteralPointerReference == nil {
		t.Error("createdLiteralPointerReference is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdLiteralPointerReferenceIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created object not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createLiteralPointerReferenceInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundLiteralPointerReferenceRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createLiteralPointerReferenceInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReferenceReference := uOfD.GetElement(foundLiteralPointerReferenceRefIdentifier)
	if redoneReferenceReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri, hl) != redoneReferenceReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdLiteralPointerReferenceIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReferenceReference.(core.ElementReference).GetReferencedElement(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestGetLiteralPointerPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getLiteralPointerPointer := uOfD.GetElementWithUri(LiteralPointerReferenceGetLiteralPointerPointerUri)
	if getLiteralPointerPointer == nil {
		t.Errorf("GetLiteralPointerPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getLiteralPointerPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getLiteralPointerPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetLiteralPointerPointer()")
	}
	sourceLiteralPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralPointerReferenceGetLiteralPointerPointerSourceLiteralPointerReferenceRefUri, hl)
	if sourceLiteralPointerReferenceRef == nil {
		t.Errorf("sourceLiteralPointerReferenceRef child not found")
	}
	indicatedLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerReferenceGetLiteralPointerPointerIndicatedLiteralPointerPointerRefUri, hl)
	if indicatedLiteralPointerPointerRef == nil {
		t.Errorf("indicatedLiteralPointerPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceLiteralPointerReference := uOfD.NewLiteralPointerReference(hl)
	dummyLiteralPointer := uOfD.NewNameLiteralPointer(hl)
	sourceLiteralPointerReference.SetReferencedLiteralPointer(dummyLiteralPointer, hl)
	sourceLiteralPointerReferenceRef.SetReferencedElement(sourceLiteralPointerReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteralPointerPointer := indicatedLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	if targetLiteralPointerPointer == nil {
		t.Errorf("Target LiteralPointerPointer not found")
		core.Print(sourceLiteralPointerReferenceRef, "SourceReference: ", hl)
		core.Print(indicatedLiteralPointerPointerRef, "TargetReference: ", hl)
	} else {
		if targetLiteralPointerPointer != sourceLiteralPointerReference.GetLiteralPointerPointer(hl) {
			t.Errorf("Target LiteralPointerPointer value incorrect")
		}
	}
}

func TestGetReferencedLiteralPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getReferencedLiteralPointer := uOfD.GetElementWithUri(LiteralPointerReferenceGetReferencedLiteralPointerUri)
	if getReferencedLiteralPointer == nil {
		t.Errorf("GetReferencedLiteralPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getReferencedLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getReferencedLiteralPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetReferencedLiteralPointer()")
	}
	sourceLiteralPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralPointerReferenceGetReferencedLiteralPointerSourceLiteralPointerReferenceRefUri, hl)
	if sourceLiteralPointerReferenceRef == nil {
		t.Errorf("sourceLiteralPointerReferenceRef child not found")
	}
	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerReferenceGetReferencedLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		t.Errorf("indicatedLiteralPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceLiteralPointerReference := uOfD.NewLiteralPointerReference(hl)
	dummyLiteralPointer := uOfD.NewNameLiteralPointer(hl)
	sourceLiteralPointerReference.SetReferencedLiteralPointer(dummyLiteralPointer, hl)
	sourceLiteralPointerReferenceRef.SetReferencedElement(sourceLiteralPointerReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	if targetLiteralPointer == nil {
		t.Errorf("Target element pointer not found")
		core.Print(sourceLiteralPointerReferenceRef, "SourceReference: ", hl)
		core.Print(indicatedLiteralPointerRef, "TargetReference: ", hl)
	} else {
		if targetLiteralPointer != sourceLiteralPointerReference.GetReferencedLiteralPointer(hl) {
			t.Errorf("Target element pointer value incorrect")
		}
	}
}

func TestSetReferencedLiteralPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	setReferencedLiteralPointer := uOfD.GetElementWithUri(LiteralPointerReferenceSetReferencedLiteralPointerUri)
	if setReferencedLiteralPointer == nil {
		t.Errorf("SetReferencedLiteralPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setReferencedLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setReferencedLiteralPointer, hl) != true {
		t.Errorf("Replicate is not refinement of SetReferencedLiteralPointer()")
	}
	sourceLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerReferenceSetReferencedLiteralPointerSourceLiteralPointerRefUri, hl)
	if sourceLiteralPointerRef == nil {
		t.Errorf("SourceLiteralPointerRef child not found")
	}
	modifiedLiteralPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralPointerReferenceSetReferencedLiteralPointerModifiedLiteralPointerReferenceRefUri, hl)
	if modifiedLiteralPointerReferenceRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceLiteralPointer := uOfD.NewNameLiteralPointer(hl)
	targetLiteralPointerReference := uOfD.NewLiteralPointerReference(hl)
	modifiedLiteralPointerReferenceRef.SetReferencedElement(targetLiteralPointerReference, hl)
	sourceLiteralPointerRef.SetReferencedLiteralPointer(sourceLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if targetLiteralPointerReference.GetReferencedLiteralPointer(hl) != sourceLiteralPointer {
		t.Errorf("LiteralPointerReference value not set")
		core.Print(sourceLiteralPointerRef, "LiteralPointerRef: ", hl)
		core.Print(modifiedLiteralPointerReferenceRef, "TargetReference: ", hl)
	}
}
