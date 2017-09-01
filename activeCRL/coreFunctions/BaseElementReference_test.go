package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"testing"
	"time"
)

func TestCreateBaseElementReferenceFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createBaseElementReference := uOfD.GetElementWithUri(BaseElementReferenceCreateUri)
	if createBaseElementReference == nil {
		t.Error("CreateBaseElementReference not found")
	}
	createdBaseElementReferenceReference := uOfD.GetElementReferenceWithUri(BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri)
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
	time.Sleep(10000000 * time.Nanosecond)

	foundBaseElementReferenceReference := core.GetChildElementReferenceWithAncestorUri(createBaseElementReferenceInstance, BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri, hl)
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
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri, hl) != redoneReferenceReference {
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
