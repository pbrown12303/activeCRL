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
	//	"time"
)

func TestCreateNameLiteralPointerFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createNameLiteralPointer := uOfD.GetElementWithUri(LiteralPointerCreateNameLiteralPointerUri)
	if createNameLiteralPointer == nil {
		t.Error("CreateNameLiteralPointer not found")
	}
	createdLiteralPointerRef := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateNameLiteralPointerCreatedLiteralPointerRefUri)
	if createdLiteralPointerRef == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createLiteralPointerInstance := uOfD.NewElement(hl)
	createLiteralPointerInstanceIdentifier := createLiteralPointerInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createNameLiteralPointer, hl)
	refinementInstance.SetRefinedElement(createLiteralPointerInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(createLiteralPointerInstance, LiteralPointerCreateNameLiteralPointerCreatedLiteralPointerRefUri, hl)
	foundLiteralPointerRefIdentifier := ""
	var createdLiteralPointer core.LiteralPointer
	createdLiteralPointerIdentifier := ""
	if foundLiteralPointerRef == nil {
		t.Error("Reference not created")
	} else {
		foundLiteralPointerRefIdentifier = foundLiteralPointerRef.GetId(hl).String()
		createdLiteralPointer = foundLiteralPointerRef.GetReferencedLiteralPointer(hl)
		if createdLiteralPointer == nil {
			t.Error("LiteralPointer not created")
		} else {
			createdLiteralPointerIdentifier = createdLiteralPointer.GetId(hl).String()
		}
	}
	if createdLiteralPointer == nil {
		t.Error("createdLiteralPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdLiteralPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetLiteral(createLiteralPointerInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetLiteral(foundLiteralPointerRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetLiteral(createdLiteralPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createLiteralPointerInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundLiteralPointerRefIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildLiteralPointerReferenceWithAncestorUri(redoneInstance, LiteralPointerCreateNameLiteralPointerCreatedLiteralPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdLiteralPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.LiteralPointerReference).GetReferencedLiteralPointer(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestCreateDefinitionLiteralPointerFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createDefinitionLiteralPointer := uOfD.GetElementWithUri(LiteralPointerCreateDefinitionLiteralPointerUri)
	if createDefinitionLiteralPointer == nil {
		t.Error("CreateDefinitionLiteralPointer not found")
	}
	createdLiteralPointerRef := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri)
	if createdLiteralPointerRef == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createLiteralPointerInstance := uOfD.NewElement(hl)
	createLiteralPointerInstanceIdentifier := createLiteralPointerInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createDefinitionLiteralPointer, hl)
	refinementInstance.SetRefinedElement(createLiteralPointerInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(createLiteralPointerInstance, LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri, hl)
	foundLiteralPointerRefIdentifier := ""
	var createdLiteralPointer core.LiteralPointer
	createdLiteralPointerIdentifier := ""
	if foundLiteralPointerRef == nil {
		t.Error("Reference not created")
	} else {
		foundLiteralPointerRefIdentifier = foundLiteralPointerRef.GetId(hl).String()
		createdLiteralPointer = foundLiteralPointerRef.GetReferencedLiteralPointer(hl)
		if createdLiteralPointer == nil {
			t.Error("LiteralPointer not created")
		} else {
			createdLiteralPointerIdentifier = createdLiteralPointer.GetId(hl).String()
		}
	}
	if createdLiteralPointer == nil {
		t.Error("createdLiteralPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdLiteralPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createLiteralPointerInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundLiteralPointerRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdLiteralPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createLiteralPointerInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundLiteralPointerRefIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildLiteralPointerReferenceWithAncestorUri(redoneInstance, LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdLiteralPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.LiteralPointerReference).GetReferencedLiteralPointer(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestCreateUriLiteralPointerFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createUriLiteralPointer := uOfD.GetElementWithUri(LiteralPointerCreateUriLiteralPointerUri)
	if createUriLiteralPointer == nil {
		t.Error("CreateUriLiteralPointer not found")
	}
	createdLiteralPointerRef := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri)
	if createdLiteralPointerRef == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createLiteralPointerInstance := uOfD.NewElement(hl)
	createLiteralPointerInstanceIdentifier := createLiteralPointerInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createUriLiteralPointer, hl)
	refinementInstance.SetRefinedElement(createLiteralPointerInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(createLiteralPointerInstance, LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri, hl)
	foundLiteralPointerRefIdentifier := ""
	var createdLiteralPointer core.LiteralPointer
	createdLiteralPointerIdentifier := ""
	if foundLiteralPointerRef == nil {
		t.Error("Reference not created")
	} else {
		foundLiteralPointerRefIdentifier = foundLiteralPointerRef.GetId(hl).String()
		createdLiteralPointer = foundLiteralPointerRef.GetReferencedLiteralPointer(hl)
		if createdLiteralPointer == nil {
			t.Error("LiteralPointer not created")
		} else {
			createdLiteralPointerIdentifier = createdLiteralPointer.GetId(hl).String()
		}
	}
	if createdLiteralPointer == nil {
		t.Error("createdLiteralPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdLiteralPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createLiteralPointerInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundLiteralPointerRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdLiteralPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createLiteralPointerInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundLiteralPointerRefIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildLiteralPointerReferenceWithAncestorUri(redoneInstance, LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdLiteralPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.LiteralPointerReference).GetReferencedLiteralPointer(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestCreateValueLiteralPointerFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createValueLiteralPointer := uOfD.GetElementWithUri(LiteralPointerCreateValueLiteralPointerUri)
	if createValueLiteralPointer == nil {
		t.Error("CreateValueLiteralPointer not found")
	}
	createdLiteralPointerRef := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri)
	if createdLiteralPointerRef == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createLiteralPointerInstance := uOfD.NewElement(hl)
	createLiteralPointerInstanceIdentifier := createLiteralPointerInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createValueLiteralPointer, hl)
	refinementInstance.SetRefinedElement(createLiteralPointerInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(createLiteralPointerInstance, LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri, hl)
	foundLiteralPointerRefIdentifier := ""
	var createdLiteralPointer core.LiteralPointer
	createdLiteralPointerIdentifier := ""
	if foundLiteralPointerRef == nil {
		t.Error("Reference not created")
	} else {
		foundLiteralPointerRefIdentifier = foundLiteralPointerRef.GetId(hl).String()
		createdLiteralPointer = foundLiteralPointerRef.GetReferencedLiteralPointer(hl)
		if createdLiteralPointer == nil {
			t.Error("LiteralPointer not created")
		} else {
			createdLiteralPointerIdentifier = createdLiteralPointer.GetId(hl).String()
		}
	}
	if createdLiteralPointer == nil {
		t.Error("createdLiteralPointer is nil")
	}
	newlyCreatedElement := uOfD.GetBaseElement(createdLiteralPointerIdentifier)
	if newlyCreatedElement == nil {
		t.Error("Created element not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createLiteralPointerInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundLiteralPointerRefIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdLiteralPointerIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createLiteralPointerInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundLiteralPointerRefIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildLiteralPointerReferenceWithAncestorUri(redoneInstance, LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetBaseElement(createdLiteralPointerIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.LiteralPointerReference).GetReferencedLiteralPointer(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}

func TestGetLiteral(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getElement := uOfD.GetElementWithUri(LiteralPointerGetLiteralUri)
	if getElement == nil {
		t.Errorf("GetLiteral function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetLiteral()")
	}
	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	indicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	sourceLiteralPointer := uOfD.NewValueLiteralPointer(hl)
	sourceLiteralPointer.SetLiteral(sourceLiteral, hl)
	sourceReference.SetReferencedLiteralPointer(sourceLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	indicatedLiteral := indicatedLiteralRef.GetReferencedLiteral(hl)
	if indicatedLiteral == nil {
		t.Errorf("Target element not found")
		core.Print(sourceReference, "SourceReference: ", hl)
		core.Print(indicatedLiteralRef, "TargetReference: ", hl)
	} else {
		if indicatedLiteral != sourceLiteral {
			t.Errorf("Target element value incorrect")
		}
	}
}

func TestGetLiteralId(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getElementId := uOfD.GetElementWithUri(LiteralPointerGetLiteralIdUri)
	if getElementId == nil {
		t.Errorf("GetLiteralId function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementId, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getElementId, hl) != true {
		t.Errorf("Replicate is not refinement of GetLiteralId()")
	}
	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIdSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIdCreatedLiteralUri, hl)
	if createdLiteralRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	sourceLiteralValue := "SourceName"
	sourceLiteral.SetLiteralValue(sourceLiteralValue, hl)
	sourceLiteralPointer := uOfD.NewValueLiteralPointer(hl)
	sourceLiteralPointer.SetLiteral(sourceLiteral, hl)
	sourceReference.SetReferencedLiteralPointer(sourceLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != sourceLiteralPointer.GetLiteralId(hl).String() {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestGetLiteralVersion(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getElementVersion := uOfD.GetElementWithUri(LiteralPointerGetLiteralVersionUri)
	if getElementVersion == nil {
		t.Errorf("GetLiteralVersion function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getElementVersion, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(getElementVersion, hl) != true {
		t.Errorf("Replicate is not refinement of GetLiteralVersion()")
	}
	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralVersionCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	sourceLiteralValue := "SourceName"
	sourceLiteral.SetLiteralValue(sourceLiteralValue, hl)
	sourceLiteralPointer := uOfD.NewValueLiteralPointer(hl)
	sourceLiteralPointer.SetLiteral(sourceLiteral, hl)
	sourceReference.SetReferencedLiteralPointer(sourceLiteralPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	targetLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != strconv.Itoa(sourceLiteralPointer.GetLiteralVersion(hl)) {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestSetLiteral(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	setElement := uOfD.GetElementWithUri(LiteralPointerSetLiteralUri)
	if setElement == nil {
		t.Errorf("SetLiteral function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetBaseElement()")
	}
	literalRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerSetLiteralLiteralRefUri, hl)
	if literalRef == nil {
		t.Errorf("ElementReference child not found")
	}
	modifiedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerSetLiteralModifiedLiteralPointerRefUri, hl)
	if modifiedLiteralPointerRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	modifiedLiteralPointer := uOfD.NewValueLiteralPointer(hl)
	modifiedLiteralPointerRef.SetReferencedLiteralPointer(modifiedLiteralPointer, hl)
	literalRef.SetReferencedLiteral(sourceLiteral, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if modifiedLiteralPointer.GetLiteral(hl) != sourceLiteral {
		t.Errorf("LiteralPointer value not set")
		core.Print(literalRef, "LiteralRef: ", hl)
		core.Print(modifiedLiteralPointerRef, "ModifiedLiteralPointerRef: ", hl)
	}
}
