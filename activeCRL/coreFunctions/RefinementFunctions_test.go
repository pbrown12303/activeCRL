package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"testing"
	"time"
)

func TestCreateRefinementFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createRefinement := uOfD.GetElementWithUri(RefinementCreateUri)
	if createRefinement == nil {
		t.Error("CreateRefinement not found")
	}
	createdRefinementRef := uOfD.GetElementReferenceWithUri(RefinementCreateCreatedRefinementRefUri)
	if createdRefinementRef == nil {
		t.Error("CreatedRefinementRef not found")
		core.Print(createRefinement, "CreateRefinement: ", hl)
	}

	createRefinementInstance := uOfD.NewElement(hl)
	createRefinementInstanceIdentifier := createRefinementInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createRefinement, hl)
	refinementInstance.SetRefinedElement(createRefinementInstance, hl)
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	foundRefinementRef := core.GetChildElementReferenceWithAncestorUri(createRefinementInstance, RefinementCreateCreatedRefinementRefUri, hl)
	foundRefinementRefIdentifier := ""
	var createdRefinement core.Refinement
	createdRefinementIdentifier := ""
	if foundRefinementRef == nil {
		t.Error("RefinementRef not created")
	} else {
		foundRefinementRefIdentifier = foundRefinementRef.GetId(hl).String()
		foundRefinement := foundRefinementRef.GetReferencedElement(hl)
		if foundRefinement == nil {
			t.Error("Refinement not created")
		} else {
			switch foundRefinement.(type) {
			case core.Refinement:
				createdRefinement = foundRefinement.(core.Refinement)
				createdRefinementIdentifier = createdRefinement.GetId(hl).String()
			default:
				t.Error("Created object of wrong type")
			}
		}
	}
	if createdRefinement == nil {
		t.Error("createdRefinement is nil")
	}
	newlyCreatedRefinement := uOfD.GetBaseElement(createdRefinementIdentifier)
	if newlyCreatedRefinement == nil {
		t.Error("Created object not in UofD")
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createRefinementInstanceIdentifier) != nil {
		t.Error("Refinement creation not undone")
	}
	if uOfD.GetElement(foundRefinementRefIdentifier) != nil {
		t.Error("Refinement creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createRefinementInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Refinement creation not redone")
	}
	redoneRefinementRef := uOfD.GetElement(foundRefinementRefIdentifier)
	if redoneRefinementRef == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, RefinementCreateCreatedRefinementRefUri, hl) != redoneRefinementRef {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedRefinement := uOfD.GetBaseElement(createdRefinementIdentifier)
		if redoneCreatedRefinement == nil {
			t.Error("Created refinement not redone")
		} else {
			if redoneRefinementRef.(core.ElementReference).GetReferencedElement(hl) != redoneCreatedRefinement {
				t.Error("Reference pointer to created refinement not restored")
			}
		}
	}
}

func TestRefinementGetAbstractElementPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getAbstractElementPointer := uOfD.GetElementWithUri(RefinementGetAbstractElementPointerUri)
	if getAbstractElementPointer == nil {
		t.Errorf("GetAbstractElementPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getAbstractElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(getAbstractElementPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetAbstractElementPointer()")
	}
	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementPointerSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		t.Errorf("sourceRefinementRef child not found")
	}
	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, RefinementGetAbstractElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		t.Errorf("indicatedElementPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceRefinement := uOfD.NewRefinement(hl)
	dummyElement := uOfD.NewElement(hl)
	sourceRefinement.SetAbstractElement(dummyElement, hl)
	sourceRefinementRef.SetReferencedElement(sourceRefinement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetReferencedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	if targetReferencedElementPointer == nil {
		t.Errorf("Target ReferencedElementPointer not found")
		core.Print(sourceRefinementRef, "SourceReference: ", hl)
		core.Print(indicatedElementPointerRef, "TargetReference: ", hl)
	} else {
		if targetReferencedElementPointer != sourceRefinement.GetAbstractElementPointer(hl) {
			t.Errorf("Target ElementPointer value incorrect")
		}
	}
}

func TestGetAbstractElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getAbstractElement := uOfD.GetElementWithUri(RefinementGetAbstractElementUri)
	if getAbstractElement == nil {
		t.Errorf("GetAbstractElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getAbstractElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(getAbstractElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetAbstractElement()")
	}
	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		t.Errorf("sourceRefinementRef child not found")
	}
	indicatedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementIndicatedElementRefUri, hl)
	if indicatedElementRef == nil {
		t.Errorf("indicatedElementRef child not found")
	}

	// Now test target reference update functionality
	sourceRefinement := uOfD.NewRefinement(hl)
	dummyElement := uOfD.NewElement(hl)
	sourceRefinement.SetAbstractElement(dummyElement, hl)
	sourceRefinementRef.SetReferencedElement(sourceRefinement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetElement := indicatedElementRef.GetReferencedElement(hl)
	if targetElement == nil {
		t.Errorf("Target element pointer not found")
		core.Print(sourceRefinementRef, "SourceReference: ", hl)
		core.Print(indicatedElementRef, "TargetReference: ", hl)
	} else {
		if targetElement != sourceRefinement.GetAbstractElement(hl) {
			t.Errorf("Target element pointer value incorrect")
		}
	}
}

func TestSetAbstractElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	setAbstractElement := uOfD.GetElementWithUri(RefinementSetAbstractElementUri)
	if setAbstractElement == nil {
		t.Errorf("SetAbstractElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setAbstractElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(setAbstractElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetAbstractElement()")
	}
	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetAbstractElementSourceElementRefUri, hl)
	if sourceElementRef == nil {
		t.Errorf("SourceElementRef child not found")
	}
	modifiedRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetAbstractElementModifiedRefinementRefUri, hl)
	if modifiedRefinementRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElement := uOfD.NewElement(hl)
	targetRefinement := uOfD.NewRefinement(hl)
	modifiedRefinementRef.SetReferencedElement(targetRefinement, hl)
	sourceElementRef.SetReferencedElement(sourceElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	if targetRefinement.GetAbstractElement(hl) != sourceElement {
		t.Errorf("Refinement value not set")
		core.Print(sourceElementRef, "ElementRef: ", hl)
		core.Print(modifiedRefinementRef, "TargetReference: ", hl)
	}
}

func TestRefinementGetRefinedElementPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getRefinedElementPointer := uOfD.GetElementWithUri(RefinementGetRefinedElementPointerUri)
	if getRefinedElementPointer == nil {
		t.Errorf("GetRefinedElementPointer function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getRefinedElementPointer, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(getRefinedElementPointer, hl) != true {
		t.Errorf("Replicate is not refinement of GetRefinedElementPointer()")
	}
	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementPointerSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		t.Errorf("sourceRefinementRef child not found")
	}
	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, RefinementGetRefinedElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		t.Errorf("indicatedElementPointerRef child not found")
	}

	// Now test target reference update functionality
	sourceRefinement := uOfD.NewRefinement(hl)
	dummyElement := uOfD.NewElement(hl)
	sourceRefinement.SetRefinedElement(dummyElement, hl)
	sourceRefinementRef.SetReferencedElement(sourceRefinement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetReferencedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	if targetReferencedElementPointer == nil {
		t.Errorf("Target ReferencedElementPointer not found")
		core.Print(sourceRefinementRef, "SourceReference: ", hl)
		core.Print(indicatedElementPointerRef, "TargetReference: ", hl)
	} else {
		if targetReferencedElementPointer != sourceRefinement.GetRefinedElementPointer(hl) {
			t.Errorf("Target ElementPointer value incorrect")
		}
	}
}

func TestGetRefinedElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getRefinedElement := uOfD.GetElementWithUri(RefinementGetRefinedElementUri)
	if getRefinedElement == nil {
		t.Errorf("GetRefinedElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getRefinedElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(getRefinedElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetRefinedElement()")
	}
	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		t.Errorf("sourceRefinementRef child not found")
	}
	indicatedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementIndicatedElementRefUri, hl)
	if indicatedElementRef == nil {
		t.Errorf("indicatedElementRef child not found")
	}

	// Now test target reference update functionality
	sourceRefinement := uOfD.NewRefinement(hl)
	dummyElement := uOfD.NewElement(hl)
	sourceRefinement.SetRefinedElement(dummyElement, hl)
	sourceRefinementRef.SetReferencedElement(sourceRefinement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetElement := indicatedElementRef.GetReferencedElement(hl)
	if targetElement == nil {
		t.Errorf("Target element pointer not found")
		core.Print(sourceRefinementRef, "SourceReference: ", hl)
		core.Print(indicatedElementRef, "TargetReference: ", hl)
	} else {
		if targetElement != sourceRefinement.GetRefinedElement(hl) {
			t.Errorf("Target element pointer value incorrect")
		}
	}
}

func TestSetRefinedElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	setRefinedElement := uOfD.GetElementWithUri(RefinementSetRefinedElementUri)
	if setRefinedElement == nil {
		t.Errorf("SetRefinedElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setRefinedElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(setRefinedElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetRefinedElement()")
	}
	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetRefinedElementSourceElementRefUri, hl)
	if sourceElementRef == nil {
		t.Errorf("SourceElementRef child not found")
	}
	modifiedRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetRefinedElementModifiedRefinementRefUri, hl)
	if modifiedRefinementRef == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	sourceElement := uOfD.NewElement(hl)
	targetRefinement := uOfD.NewRefinement(hl)
	modifiedRefinementRef.SetReferencedElement(targetRefinement, hl)
	sourceElementRef.SetReferencedElement(sourceElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	if targetRefinement.GetRefinedElement(hl) != sourceElement {
		t.Errorf("Refinement value not set")
		core.Print(sourceElementRef, "ElementRef: ", hl)
		core.Print(modifiedRefinementRef, "TargetReference: ", hl)
	}
}
