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

func TestCreateBaseElementReferenceFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createBaseElementReference := uOfD.GetElementWithUri(BaseElementReferenceCreateUri)
	if createBaseElementReference == nil {
		t.Error("CreateBaseElementReference not found")
	}
	createdBaseElementReferenceReference := uOfD.GetElementReferenceWithUri(BaseElementReferenceCreateCreatedBaseElementReferenceRefUri)
	if createdBaseElementReferenceReference == nil {
		t.Error("CreatedBaseElementReferenceReference not found")
		core.Print(createBaseElementReference, "CreateBaseElementReference: ", hl)
	}

	createBaseElementReferenceInstance := uOfD.NewElement(hl)
	createBaseElementReferenceInstanceIdentifier := createBaseElementReferenceInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createBaseElementReference, hl)

	refinementInstance.SetRefinedElement(createBaseElementReferenceInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundBaseElementReferenceReference := core.GetChildElementReferenceWithAncestorUri(createBaseElementReferenceInstance, BaseElementReferenceCreateCreatedBaseElementReferenceRefUri, hl)
	foundBaseElementReferenceReferenceIdentifier := ""
	var createdBaseElementReference core.BaseElementReference
	createdBaseElementReferenceIdentifier := ""
	if foundBaseElementReferenceReference == nil {
		t.Error("BaseElementReference not created")
	} else {
		foundBaseElementReferenceReferenceIdentifier = foundBaseElementReferenceReference.GetId(hl).String()
		foundBaseElementReference := foundBaseElementReferenceReference.GetReferencedElement(hl)
		if foundBaseElementReference == nil {
			t.Error("Element not created")
		} else {
			switch foundBaseElementReference.(type) {
			case core.BaseElementReference:
				createdBaseElementReference = foundBaseElementReference.(core.BaseElementReference)
				createdBaseElementReferenceIdentifier = createdBaseElementReference.GetId(hl).String()
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdBaseElementReference == nil {
		t.Error("createdBaseElementReference is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdBaseElementReferenceIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created object not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createBaseElementReferenceInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundBaseElementReferenceReferenceIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createBaseElementReferenceInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReferenceReference := uOfD.GetElement(foundBaseElementReferenceReferenceIdentifier)
	if redoneReferenceReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, BaseElementReferenceCreateCreatedBaseElementReferenceRefUri, hl) != redoneReferenceReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdBaseElementReferenceIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReferenceReference.(core.ElementReference).GetReferencedElement(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestGetBaseElementPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getBaseElementPointer := uOfD.GetElementWithUri(BaseElementReferenceGetBaseElementPointerUri)
	if getBaseElementPointer == nil {
		t.Errorf("GetBaseElementPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getBaseElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getBaseElementPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetBaseElementPointer()")
	}
	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementReferenceGetBaseElementPointerSourceBaseElementReferenceRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementReferenceGetBaseElementPointerIndicatedBaseElementPointerRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceBaseElementReference := uOfD.NewBaseElementReference(hl)
	dummyElement := uOfD.NewElement(hl)
	sourceBaseElementReference.SetReferencedBaseElement(dummyElement, hl)
	sourceReference.SetReferencedElement(sourceBaseElementReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetBaseElement := targetReference.GetReferencedBaseElement(hl)
	if targetBaseElement == nil {
		t.Errorf("Target element not found")
		core.Print(sourceReference, "SourceReference: ", hl)
		core.Print(targetReference, "TargetReference: ", hl)
	} else {
		if targetBaseElement != sourceBaseElementReference.GetBaseElementPointer(hl) {
			t.Errorf("Target element value incorrect")
		}
	}
}

func TestGetReferencedBaseElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getReferencedBaseElement := uOfD.GetElementWithUri(BaseElementReferenceGetReferencedBaseElementUri)
	if getReferencedBaseElement == nil {
		t.Errorf("GetReferencedBaseElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getReferencedBaseElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getReferencedBaseElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetReferencedBaseElement()")
	}
	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementReferenceGetReferencedBaseElementSourceBaseElementReferenceRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementReferenceGetReferencedBaseElementIndicatedBaseElementRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceBaseElementReference := uOfD.NewBaseElementReference(hl)
	dummyElement := uOfD.NewElement(hl)
	sourceBaseElementReference.SetReferencedBaseElement(dummyElement, hl)
	sourceReference.SetReferencedElement(sourceBaseElementReference, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetBaseElement := targetReference.GetReferencedBaseElement(hl)
	if targetBaseElement == nil {
		t.Errorf("Target element not found")
		core.Print(sourceReference, "SourceReference: ", hl)
		core.Print(targetReference, "TargetReference: ", hl)
	} else {
		if targetBaseElement != sourceBaseElementReference.GetReferencedBaseElement(hl) {
			t.Errorf("Target element value incorrect")
		}
	}
}

func TestSetReferencedBaseElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	setReferencedBaseElement := uOfD.GetElementWithUri(BaseElementReferenceSetReferencedBaseElementUri)
	if setReferencedBaseElement == nil {
		t.Errorf("SetReferencedBaseElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setReferencedBaseElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setReferencedBaseElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetReferencedBaseElement()")
	}
	baseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementReferenceSetReferencedBaseElementSourceBaseElementRefUri, hl)
	if baseElementReference == nil {
		t.Errorf("BaseElementReference child not found")
	}
	targetReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementReferenceSetReferencedBaseElementModifiedBaseElementReferenceRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceBaseElement := uOfD.NewElement(hl)
	targetBaseElementReference := uOfD.NewBaseElementReference(hl)
	targetReference.SetReferencedElement(targetBaseElementReference, hl)
	baseElementReference.SetReferencedBaseElement(sourceBaseElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if targetBaseElementReference.GetReferencedBaseElement(hl) != sourceBaseElement {
		t.Errorf("BaseElementReference value not set")
		core.Print(baseElementReference, "BaseElementReference: ", hl)
		core.Print(targetReference, "TargetReference: ", hl)
	}
}
