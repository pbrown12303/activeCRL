package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"strconv"
	"testing"
	"time"
)

func TestCreateBaseElementPointerFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createBaseElementPointer := uOfD.GetElementWithUri(BaseElementPointerCreateUri)
	if createBaseElementPointer == nil {
		t.Error("CreateBaseElementPointer not found")
	}
	createdBaseElementReference := uOfD.GetBaseElementReferenceWithUri(BaseElementPointerCreateCreatedBaseElementPointerReferenceUri)
	if createdBaseElementReference == nil {
		t.Error("CreatedBaseElementReference not found")
	}

	// Now create the instance of the function
	createBaseElementPointerInstance := uOfD.NewElement(hl)
	createBaseElementPointerInstanceIdentifier := createBaseElementPointerInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createBaseElementPointer, hl)

	refinementInstance.SetRefinedElement(createBaseElementPointerInstance, hl)
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	foundReference := core.GetChildBaseElementReferenceWithAncestorUri(createBaseElementPointerInstance, BaseElementPointerCreateCreatedBaseElementPointerReferenceUri, hl)
	foundReferenceIdentifier := ""
	var createdBaseElementPointer core.BaseElementPointer
	createdBaseElementPointerIdentifier := ""
	if foundReference == nil {
		t.Error("Reference not created")
	} else {
		foundReferenceIdentifier = foundReference.GetId(hl).String()
		foundBaseElement := foundReference.GetReferencedBaseElement(hl)
		if foundBaseElement == nil {
			t.Error("Element not created")
		} else {
			switch foundBaseElement.(type) {
			case core.BaseElementPointer:
				createdBaseElementPointer = foundBaseElement.(core.BaseElementPointer)
				createdBaseElementPointerIdentifier = createdBaseElementPointer.GetId(hl).String()
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
		if core.GetChildBaseElementReferenceWithAncestorUri(redoneInstance, BaseElementPointerCreateCreatedBaseElementPointerReferenceUri, hl) != redoneReference {
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
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getBaseElement := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementUri)
	if getBaseElement == nil {
		t.Errorf("GetBaseElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getBaseElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(getBaseElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetBaseElement()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementSourceReferenceUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementTargetBaseElementReferenceUri, hl)
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
	time.Sleep(10000000 * time.Nanosecond)

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
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getBaseElementId := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementIdUri)
	if getBaseElementVersion == nil {
		t.Errorf("GetBaseElementId function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getBaseElementId, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(getBaseElementId, hl) != true {
		t.Errorf("Replicate is not refinement of GetBaseElementId()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIdSourceReferenceUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIdTargetLiteralReferenceUri, hl)
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
	time.Sleep(10000000 * time.Nanosecond)

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
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	getBaseElementVersion := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementVersionUri)
	if getBaseElementVersion == nil {
		t.Errorf("GetBaseElementVersion function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getBaseElementVersion, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(getBaseElementVersion, hl) != true {
		t.Errorf("Replicate is not refinement of GetBaseElementVersion()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementVersionSourceReferenceUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementVersionTargetLiteralReferenceUri, hl)
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
	time.Sleep(10000000 * time.Nanosecond)

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
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	setBaseElement := uOfD.GetElementWithUri(BaseElementPointerSetBaseElementUri)
	if setBaseElement == nil {
		t.Errorf("SetBaseElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setBaseElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if replicate.IsRefinementOf(setBaseElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetBaseElement()")
	}
	baseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerSetBaseElementBaseElementReferenceUri, hl)
	if baseElementReference == nil {
		t.Errorf("BaseElementReference child not found")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerSetBaseElementTargetBaseElementPointerReferenceUri, hl)
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
	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	if targetBaseElementPointer.GetBaseElement(hl) != sourceBaseElement {
		t.Errorf("BaseElementPointer value not set")
		core.Print(baseElementReference, "BaseElementReference: ", hl)
		core.Print(targetReference, "TargetReference: ", hl)
	}
}
