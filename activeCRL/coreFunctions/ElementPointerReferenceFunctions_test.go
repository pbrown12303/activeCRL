package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"sync"
	"testing"
	//	"time"
)

func TestCreateElementPointerReferenceFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

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
	createElementPointerReferenceInstanceIdentifier := createElementPointerReferenceInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createElementPointerReference, hl)

	refinementInstance.SetRefinedElement(createElementPointerReferenceInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	foundElementPointerReferenceRef := core.GetChildElementReferenceWithAncestorUri(createElementPointerReferenceInstance, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri, hl)
	foundElementPointerReferenceRefIdentifier := ""
	var createdElementPointerReference core.ElementPointerReference
	createdElementPointerReferenceIdentifier := ""
	if foundElementPointerReferenceRef == nil {
		t.Error("ElementPointerReferenceRef not created")
	} else {
		foundElementPointerReferenceRefIdentifier = foundElementPointerReferenceRef.GetId(hl).String()
		foundElementPointerReference := foundElementPointerReferenceRef.GetReferencedElement(hl)
		if foundElementPointerReference == nil {
			t.Error("ElementPointerReference not created")
		} else {
			switch foundElementPointerReference.(type) {
			case core.ElementPointerReference:
				createdElementPointerReference = foundElementPointerReference.(core.ElementPointerReference)
				createdElementPointerReferenceIdentifier = createdElementPointerReference.GetId(hl).String()
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
	GetCoreFunctionsConceptSpace(uOfD)

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
	if replicate.IsRefinementOf(getElementPointerPointer, hl) != true {
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
	GetCoreFunctionsConceptSpace(uOfD)

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
	if replicate.IsRefinementOf(getReferencedElementPointer, hl) != true {
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
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	setReferencedElementPointer := uOfD.GetElementWithUri(ElementPointerReferenceSetReferencedElementPointerUri)
	if setReferencedElementPointer == nil {
		t.Errorf("SetReferencedElementPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setReferencedElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setReferencedElementPointer, hl) != true {
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
