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

func TestBaseElementPointerFunctionsIds(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var BaseElementPointerFunctionsUri string = CoreFunctionsPrefix + "BaseElementPointerFunctions"
	validateElementId(t, uOfD, hl, BaseElementPointerFunctionsUri)
	//
	//var BaseElementPointerCreateUri string = CoreFunctionsPrefix + "BaseElementPointer/Create"
	validateElementId(t, uOfD, hl, BaseElementPointerCreateUri)
	//var BaseElementPointerCreateCreatedBaseElementPointerRefUri = CoreFunctionsPrefix + "BaseElementPointer/Create/CreatedBaseElementPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementPointerCreateCreatedBaseElementPointerRefUri)
	//
	//var BaseElementPointerGetBaseElementUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement"
	validateElementId(t, uOfD, hl, BaseElementPointerGetBaseElementUri)
	//var BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement/SourceBaseElementPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri)
	//var BaseElementPointerGetBaseElementIndicatedBaseElementRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement/IndicatedBaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementPointerGetBaseElementIndicatedBaseElementRefUri)
	//
	//var BaseElementPointerGetBaseElementIdUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId"
	validateElementId(t, uOfD, hl, BaseElementPointerGetBaseElementIdUri)
	//var BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId/SourceBaseElementPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri)
	//var BaseElementPointerGetBaseElementIdCreatedLiteralUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, BaseElementPointerGetBaseElementIdCreatedLiteralUri)
	//
	//var BaseElementPointerGetBaseElementVersionUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion"
	validateElementId(t, uOfD, hl, BaseElementPointerGetBaseElementVersionUri)
	//var BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion/SourceBaseElementPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri)
	//var BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri)
	//
	//var BaseElementPointerSetBaseElementUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement"
	validateElementId(t, uOfD, hl, BaseElementPointerSetBaseElementUri)
	//var BaseElementPointerSetBaseElementBaseElementRefUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement/BaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementPointerSetBaseElementBaseElementRefUri)
	//var BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement/ModifiedBaseElementPointerRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri)
}

func TestCreateBaseElementPointerFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createBaseElementPointer := uOfD.GetElementWithUri(BaseElementPointerCreateUri)
	if createBaseElementPointer == nil {
		t.Error("CreateBaseElementPointer not found")
	}
	createdBaseElementReference := uOfD.GetBaseElementReferenceWithUri(BaseElementPointerCreateCreatedBaseElementPointerRefUri)
	if createdBaseElementReference == nil {
		t.Error("CreatedBaseElementReference not found")
	}

	// Now create the instance of the function
	createBaseElementPointerInstance := uOfD.NewElement(hl)
	createBaseElementPointerInstanceIdentifier := createBaseElementPointerInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createBaseElementPointer, hl)

	refinementInstance.SetRefinedElement(createBaseElementPointerInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundReference := core.GetChildBaseElementReferenceWithAncestorUri(createBaseElementPointerInstance, BaseElementPointerCreateCreatedBaseElementPointerRefUri, hl)
	foundReferenceIdentifier := uuid.Nil
	var createdBaseElementPointer core.BaseElementPointer
	createdBaseElementPointerIdentifier := uuid.Nil
	if foundReference == nil {
		t.Error("Reference not created")
	} else {
		foundReferenceIdentifier = foundReference.GetId(hl)
		foundBaseElement := foundReference.GetReferencedBaseElement(hl)
		if foundBaseElement == nil {
			t.Error("Element not created")
		} else {
			switch foundBaseElement.(type) {
			case core.BaseElementPointer:
				createdBaseElementPointer = foundBaseElement.(core.BaseElementPointer)
				createdBaseElementPointerIdentifier = createdBaseElementPointer.GetId(hl)
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdBaseElementPointer == nil {
		t.Error("createdBaseElementPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdBaseElementPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createBaseElementPointerInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundReferenceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdBaseElementPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createBaseElementPointerInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundReferenceIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildBaseElementReferenceWithAncestorUri(redoneInstance, BaseElementPointerCreateCreatedBaseElementPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdBaseElementPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.BaseElementReference).GetReferencedBaseElement(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestGetBaseElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getBaseElement := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementUri)
	if getBaseElement == nil {
		t.Errorf("GetBaseElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getBaseElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getBaseElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetBaseElement()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIndicatedBaseElementRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceBaseElement := uOfD.NewElement(hl)
	sourceBaseElementPointer := uOfD.NewBaseElementPointer(hl)
	sourceBaseElementPointer.SetBaseElement(sourceBaseElement, hl)
	sourceReference.SetReferencedBaseElement(sourceBaseElementPointer, hl)

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
		if targetBaseElement != sourceBaseElement {
			t.Errorf("Target element value incorrect")
		}
	}
}

func TestGetBaseElementId(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getBaseElementId := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementIdUri)
	if getBaseElementVersion == nil {
		t.Errorf("GetBaseElementId function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getBaseElementId, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getBaseElementId, hl) != true {
		t.Errorf("Replicate is not refinement of GetBaseElementId()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIdCreatedLiteralUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceBaseElement := uOfD.NewElement(hl)
	sourceName := "SourceName"
	core.SetName(sourceBaseElement, sourceName, hl)
	sourceBaseElementPointer := uOfD.NewBaseElementPointer(hl)
	sourceBaseElementPointer.SetBaseElement(sourceBaseElement, hl)
	sourceReference.SetReferencedBaseElement(sourceBaseElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := targetReference.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != sourceBaseElementPointer.GetBaseElementId(hl).String() {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestGetBaseElementVersion(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getBaseElementVersion := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementVersionUri)
	if getBaseElementVersion == nil {
		t.Errorf("GetBaseElementVersion function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getBaseElementVersion, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getBaseElementVersion, hl) != true {
		t.Errorf("Replicate is not refinement of GetBaseElementVersion()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceBaseElement := uOfD.NewElement(hl)
	sourceName := "SourceName"
	core.SetName(sourceBaseElement, sourceName, hl)
	sourceBaseElementPointer := uOfD.NewBaseElementPointer(hl)
	sourceBaseElementPointer.SetBaseElement(sourceBaseElement, hl)
	sourceReference.SetReferencedBaseElement(sourceBaseElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := targetReference.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != strconv.Itoa(sourceBaseElementPointer.GetBaseElementVersion(hl)) {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestSetBaseElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setBaseElement := uOfD.GetElementWithUri(BaseElementPointerSetBaseElementUri)
	if setBaseElement == nil {
		t.Errorf("SetBaseElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setBaseElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setBaseElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetBaseElement()")
	}
	baseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerSetBaseElementBaseElementRefUri, hl)
	if baseElementReference == nil {
		t.Errorf("BaseElementReference child not found")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceBaseElement := uOfD.NewElement(hl)
	targetBaseElementPointer := uOfD.NewBaseElementPointer(hl)
	targetReference.SetReferencedBaseElement(targetBaseElementPointer, hl)
	baseElementReference.SetReferencedBaseElement(sourceBaseElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if targetBaseElementPointer.GetBaseElement(hl) != sourceBaseElement {
		t.Errorf("BaseElementPointer value not set")
		core.Print(baseElementReference, "BaseElementReference: ", hl)
		core.Print(targetReference, "TargetReference: ", hl)
	}
}
