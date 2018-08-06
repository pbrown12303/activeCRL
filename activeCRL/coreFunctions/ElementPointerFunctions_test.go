// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestElementPointerFunctionsIds(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var ElementPointerFunctionsUri string = CoreFunctionsPrefix + "ElementPointerFunctions"
	validateElementId(t, uOfD, hl, ElementPointerFunctionsUri)
	//
	//var ElementPointerCreateAbstractElementPointerUri string = CoreFunctionsPrefix + "ElementPointer/CreateAbstractElementPointer"
	validateElementId(t, uOfD, hl, ElementPointerCreateAbstractElementPointerUri)
	//var ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri = CoreFunctionsPrefix + "ElementPointer/CreateAbstractElementPointer/CreatedElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri)
	//
	//var ElementPointerCreateRefinedElementPointerUri string = CoreFunctionsPrefix + "ElementPointer/CreateRefinedElementPointer"
	validateElementId(t, uOfD, hl, ElementPointerCreateRefinedElementPointerUri)
	//var ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri = CoreFunctionsPrefix + "ElementPointer/CreateRefinedElementPointer/CreatedElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri)
	//
	//var ElementPointerCreateOwningElementPointerUri string = CoreFunctionsPrefix + "ElementPointer/CreateOwningElementPointer"
	validateElementId(t, uOfD, hl, ElementPointerCreateOwningElementPointerUri)
	//var ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri = CoreFunctionsPrefix + "ElementPointer/CreateOwningElementPointer/CreatedElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri)
	//
	//var ElementPointerCreateReferencedElementPointerUri string = CoreFunctionsPrefix + "ElementPointer/CreateReferencedElementPointer"
	validateElementId(t, uOfD, hl, ElementPointerCreateReferencedElementPointerUri)
	//var ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri = CoreFunctionsPrefix + "ElementPointer/CreateReferencedElementPointer/CreatedElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri)
	//
	//var ElementPointerGetElementUri string = CoreFunctionsPrefix + "ElementPointer/GetElement"
	validateElementId(t, uOfD, hl, ElementPointerGetElementUri)
	//var ElementPointerGetElementSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElement/SourceElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerGetElementSourceElementPointerRefUri)
	//var ElementPointerGetElementIndicatedElementRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElement/IndicatedElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementPointerGetElementIndicatedElementRefUri)
	//
	//var ElementPointerGetElementIdUri string = CoreFunctionsPrefix + "ElementPointer/GetElementId"
	validateElementId(t, uOfD, hl, ElementPointerGetElementIdUri)
	//var ElementPointerGetElementIdSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementId/SourceElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerGetElementIdSourceElementPointerRefUri)
	//var ElementPointerGetElementIdCreatedLiteralUri string = CoreFunctionsPrefix + "ElementPointer/GetElementId/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementPointerGetElementIdCreatedLiteralUri)
	//
	//var ElementPointerGetElementPointerRoleUri string = CoreFunctionsPrefix + "ElementPointer/GetElementPointerRole"
	validateElementId(t, uOfD, hl, ElementPointerGetElementPointerRoleUri)
	//var ElementPointerGetElementPointerRoleSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementPointerRole/SourceElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerGetElementPointerRoleSourceElementPointerRefUri)
	//var ElementPointerGetElementPointerRoleCreatedLiteralRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementPointerRole/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementPointerGetElementPointerRoleCreatedLiteralRefUri)
	//
	//var ElementPointerGetElementVersionUri string = CoreFunctionsPrefix + "ElementPointer/GetElementVersion"
	validateElementId(t, uOfD, hl, ElementPointerGetElementVersionUri)
	//var ElementPointerGetElementVersionSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementVersion/SourceElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerGetElementVersionSourceElementPointerRefUri)
	//var ElementPointerGetElementVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementVersion/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementPointerGetElementVersionCreatedLiteralRefUri)
	//
	//var ElementPointerSetElementUri string = CoreFunctionsPrefix + "ElementPointer/SetElement"
	validateElementId(t, uOfD, hl, ElementPointerSetElementUri)
	//var ElementPointerSetElementElementRefUri string = CoreFunctionsPrefix + "ElementPointer/SetElement/ElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementPointerSetElementElementRefUri)
	//var ElementPointerSetElementModifiedElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/SetElement/ModifiedElementPointerRef"
	validateElementPointerReferenceId(t, uOfD, hl, ElementPointerSetElementModifiedElementPointerRefUri)
}

func TestCreateAbstractElementPointerFunction(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createAbstractElementPointer := uOfD.GetElementWithUri(ElementPointerCreateAbstractElementPointerUri)
	if createAbstractElementPointer == nil {
		t.Error("CreateAbstractElementPointer not found")
	}
	createdElementPointerRef := uOfD.GetElementPointerReferenceWithUri(ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri)
	if createdElementPointerRef == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createElementPointerFunctionInstance := uOfD.NewElement(hl)
	createElementPointerFunctionInstanceIdentifier := createElementPointerFunctionInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createAbstractElementPointer, hl)
	refinementInstance.SetRefinedElement(createElementPointerFunctionInstance, hl)
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	foundElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(createElementPointerFunctionInstance, ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri, hl)
	foundElementPointerRefIdentifier := ""
	var foundElementPointer core.ElementPointer
	foundElementPointerIdentifier := ""
	if foundElementPointerRef == nil {
		t.Error("Reference not created")
	} else {
		foundElementPointerRefIdentifier = foundElementPointerRef.GetId(hl)
		foundElementPointer = foundElementPointerRef.GetReferencedElementPointer(hl)
		if foundElementPointer == nil {
			t.Error("ElementPointer not created")
		} else {
			foundElementPointerIdentifier = foundElementPointer.GetId(hl)
		}
	}
	if foundElementPointer == nil {
		t.Error("foundElementPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(foundElementPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createElementPointerFunctionInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundElementPointerRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundElementPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createElementPointerFunctionInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundElementPointerRefIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementPointerReferenceWithAncestorUri(redoneInstance, ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(foundElementPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.ElementPointerReference).GetReferencedElementPointer(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestCreateRefinedElementPointerFunction(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createRefinedElementPointer := uOfD.GetElementWithUri(ElementPointerCreateRefinedElementPointerUri)
	if createRefinedElementPointer == nil {
		t.Error("CreateRefinedElementPointer not found")
	}
	createdElementPointerRef := uOfD.GetElementPointerReferenceWithUri(ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri)
	if createdElementPointerRef == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createElementPointerFunctionInstance := uOfD.NewElement(hl)
	createElementPointerFunctionInstanceIdentifier := createElementPointerFunctionInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createRefinedElementPointer, hl)
	refinementInstance.SetRefinedElement(createElementPointerFunctionInstance, hl)
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	foundElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(createElementPointerFunctionInstance, ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri, hl)
	foundElementPointerRefIdentifier := ""
	var createdElementPointer core.ElementPointer
	createdElementPointerIdentifier := ""
	if foundElementPointerRef == nil {
		t.Error("Reference not created")
	} else {
		foundElementPointerRefIdentifier = foundElementPointerRef.GetId(hl)
		foundElementPointer := foundElementPointerRef.GetReferencedElementPointer(hl)
		if foundElementPointer == nil {
			t.Error("ElementPointer not created")
		} else {
			switch foundElementPointer.(type) {
			case core.ElementPointer:
				createdElementPointer = foundElementPointer.(core.ElementPointer)
				createdElementPointerIdentifier = createdElementPointer.GetId(hl)
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdElementPointer == nil {
		t.Error("createdElementPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdElementPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createElementPointerFunctionInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundElementPointerRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdElementPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createElementPointerFunctionInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundElementPointerRefIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementPointerReferenceWithAncestorUri(redoneInstance, ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdElementPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.ElementPointerReference).GetReferencedElementPointer(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestCreateOwningElementPointerFunction(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createOwningElementPointer := uOfD.GetElementWithUri(ElementPointerCreateOwningElementPointerUri)
	if createOwningElementPointer == nil {
		t.Error("CreateOwningElementPointer not found")
	}
	createdElementPointerRef := uOfD.GetElementPointerReferenceWithUri(ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri)
	if createdElementPointerRef == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createElementPointerFunctionInstance := uOfD.NewElement(hl)
	createElementPointerFunctionInstanceIdentifier := createElementPointerFunctionInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createOwningElementPointer, hl)
	refinementInstance.SetRefinedElement(createElementPointerFunctionInstance, hl)
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	foundElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(createElementPointerFunctionInstance, ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri, hl)
	foundElementPointerRefIdentifier := ""
	var createdElementPointer core.ElementPointer
	createdElementPointerIdentifier := ""
	if foundElementPointerRef == nil {
		t.Error("Reference not created")
	} else {
		foundElementPointerRefIdentifier = foundElementPointerRef.GetId(hl)
		foundElementPointer := foundElementPointerRef.GetReferencedElementPointer(hl)
		if foundElementPointer == nil {
			t.Error("ElementPointer not created")
		} else {
			switch foundElementPointer.(type) {
			case core.ElementPointer:
				createdElementPointer = foundElementPointer.(core.ElementPointer)
				createdElementPointerIdentifier = createdElementPointer.GetId(hl)
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdElementPointer == nil {
		t.Error("createdElementPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdElementPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createElementPointerFunctionInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundElementPointerRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdElementPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createElementPointerFunctionInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundElementPointerRefIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementPointerReferenceWithAncestorUri(redoneInstance, ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdElementPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.ElementPointerReference).GetReferencedElementPointer(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestCreateReferencedElementPointerFunction(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createReferencedElementPointer := uOfD.GetElementWithUri(ElementPointerCreateReferencedElementPointerUri)
	if createReferencedElementPointer == nil {
		t.Error("CreateReferencedElementPointer not found")
	}
	createdElementPointerRef := uOfD.GetElementPointerReferenceWithUri(ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri)
	if createdElementPointerRef == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createElementPointerFunctionInstance := uOfD.NewElement(hl)
	createElementPointerFunctionInstanceIdentifier := createElementPointerFunctionInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createReferencedElementPointer, hl)
	refinementInstance.SetRefinedElement(createElementPointerFunctionInstance, hl)
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	foundElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(createElementPointerFunctionInstance, ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri, hl)
	foundElementPointerRefIdentifier := ""
	var createdElementPointer core.ElementPointer
	createdElementPointerIdentifier := ""
	if foundElementPointerRef == nil {
		t.Error("Reference not created")
	} else {
		foundElementPointerRefIdentifier = foundElementPointerRef.GetId(hl)
		foundElementPointer := foundElementPointerRef.GetReferencedElementPointer(hl)
		if foundElementPointer == nil {
			t.Error("ElementPointer not created")
		} else {
			switch foundElementPointer.(type) {
			case core.ElementPointer:
				createdElementPointer = foundElementPointer.(core.ElementPointer)
				createdElementPointerIdentifier = createdElementPointer.GetId(hl)
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdElementPointer == nil {
		t.Error("createdElementPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdElementPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createElementPointerFunctionInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundElementPointerRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdElementPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createElementPointerFunctionInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundElementPointerRefIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementPointerReferenceWithAncestorUri(redoneInstance, ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdElementPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.ElementPointerReference).GetReferencedElementPointer(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestGetElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getElement := uOfD.GetElementWithUri(ElementPointerGetElementUri)
	if getElement == nil {
		t.Errorf("GetElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetElement()")
	}
	sourceReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerGetElementSourceElementPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerGetElementIndicatedElementRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElement := uOfD.NewElement(hl)
	sourceElementPointer := uOfD.NewReferencedElementPointer(hl)
	sourceElementPointer.SetElement(sourceElement, hl)
	sourceReference.SetReferencedElementPointer(sourceElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetElement := targetReference.GetReferencedElement(hl)
	if targetElement == nil {
		t.Errorf("Target element not found")
		core.Print(sourceReference, "SourceReference: ", hl)
		core.Print(targetReference, "TargetReference: ", hl)
	} else {
		if targetElement != sourceElement {
			t.Errorf("Target element value incorrect")
		}
	}
}

func TestGetElementId(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getElementId := uOfD.GetElementWithUri(ElementPointerGetElementIdUri)
	if getElementId == nil {
		t.Errorf("GetElementId function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementId, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getElementId, hl) != true {
		t.Errorf("Replicate is not refinement of GetElementId()")
	}
	sourceReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerGetElementIdSourceElementPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerGetElementIdCreatedLiteralUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElement := uOfD.NewElement(hl)
	sourceLabel := "SourceLabel"
	core.SetLabel(sourceElement, sourceLabel, hl)
	sourceElementPointer := uOfD.NewReferencedElementPointer(hl)
	sourceElementPointer.SetElement(sourceElement, hl)
	sourceReference.SetReferencedElementPointer(sourceElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetLiteral := targetReference.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != sourceElementPointer.GetElementId(hl) {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestGetElementVersion(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getElementVersion := uOfD.GetElementWithUri(ElementPointerGetElementVersionUri)
	if getElementVersion == nil {
		t.Errorf("GetElementVersion function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementVersion, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getElementVersion, hl) != true {
		t.Errorf("Replicate is not refinement of GetElementVersion()")
	}
	sourceReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerGetElementVersionSourceElementPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerGetElementVersionCreatedLiteralRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElement := uOfD.NewElement(hl)
	sourceLabel := "SourceLabel"
	core.SetLabel(sourceElement, sourceLabel, hl)
	sourceElementPointer := uOfD.NewReferencedElementPointer(hl)
	sourceElementPointer.SetElement(sourceElement, hl)
	sourceReference.SetReferencedElementPointer(sourceElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetLiteral := targetReference.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != strconv.Itoa(sourceElementPointer.GetElementVersion(hl)) {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestSetElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setElement := uOfD.GetElementWithUri(ElementPointerSetElementUri)
	if setElement == nil {
		t.Errorf("SetElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, setElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetBaseElement()")
	}
	elementReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerSetElementElementRefUri, hl)
	if elementReference == nil {
		t.Errorf("ElementReference child not found")
	}
	targetReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerSetElementModifiedElementPointerRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElement := uOfD.NewElement(hl)
	targetElementPointer := uOfD.NewReferencedElementPointer(hl)
	targetReference.SetReferencedElementPointer(targetElementPointer, hl)
	elementReference.SetReferencedElement(sourceElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	if targetElementPointer.GetElement(hl) != sourceElement {
		t.Errorf("ElementPointer value not set")
		core.Print(elementReference, "ElementReference: ", hl)
		core.Print(targetReference, "TargetReference: ", hl)
	}
}
