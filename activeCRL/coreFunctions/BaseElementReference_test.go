package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
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

	//	// Now create the instance of the function
	log.Printf("Before refinement")
	core.Print(createBaseElementReference, "CreateBaseElementReference: ", hl)
	createElement := uOfD.GetElementWithUri(ElementCreateUri)
	core.Print(createElement, "CreateElement: ", hl)
	//	core.PrintUriIndexJustIdentifiers(uOfD, hl)
	createBaseElementReferenceInstance := uOfD.NewElement(hl)
	//	//	createBaseElementReferenceInstanceIdentifier := createBaseElementReferenceInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createBaseElementReference, hl)
	//
	refinementInstance.SetRefinedElement(createBaseElementReferenceInstance, hl)
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	log.Printf("After refinement")
	createElement = uOfD.GetElementWithUri(ElementCreateUri)
	core.Print(createElement, "CreateElement: ", hl)
	//	core.PrintUriIndexJustIdentifiers(uOfD, hl)

	//	foundBaseElementReference := core.GetChildElementReferenceWithAncestorUri(createBaseElementReferenceInstance, BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri, hl)
	//	//	foundBaseElementReferenceIdentifier := ""
	//	var createdBaseElementReference core.BaseElementReference
	//	createdBaseElementReferenceIdentifier := ""
	//	if foundBaseElementReference == nil {
	//		t.Error("BaseElementReference not created")
	//	}
	//	else {
	//		//		foundBaseElementReferenceIdentifier = foundBaseElementReference.GetId(hl).String()
	//		foundBaseElementReference := foundBaseElementReference.GetReferencedElement(hl)
	//		if foundBaseElementReference == nil {
	//			t.Error("Element not created")
	//		} else {
	//			switch foundBaseElementReference.(type) {
	//			case core.BaseElementReference:
	//				createdBaseElementReference = foundBaseElementReference.(core.BaseElementReference)
	//				createdBaseElementReferenceIdentifier = createdBaseElementReference.GetId(hl).String()
	//			default:
	//				t.Error("Created object of wrong type")
	//			}
	//		}
	//	}
	//	if createdBaseElementReference == nil {
	//		t.Error("createdBaseElementReference is nil")
	//	}
	//	newlyCreatedElement := uOfD.GetBaseElement(createdBaseElementReferenceIdentifier)
	//	if newlyCreatedElement == nil {
	//		t.Error("Created object not in UofD")
	//	}
	//
	//	// Now undo
	//	uOfD.Undo(hl)
	//	if uOfD.GetElement(createBaseElementReferenceInstanceIdentifier) != nil {
	//		t.Error("Element creation not undone")
	//	}
	//	if uOfD.GetElement(foundBaseElementReferenceIdentifier) != nil {
	//		t.Error("Element creation not undone")
	//	}
	//	if uOfD.GetElement(createdBaseElementReferenceIdentifier) != nil {
	//		t.Error("Element creation not undone")
	//	}
	//
	//	// Now Redo
	//	uOfD.Redo(hl)
	//	redoneInstance := uOfD.GetElement(createBaseElementReferenceInstanceIdentifier)
	//	if redoneInstance == nil {
	//		t.Error("Element creation not redone")
	//	}
	//	redoneReference := uOfD.GetElement(foundBaseElementReferenceIdentifier)
	//	if redoneReference == nil {
	//		t.Error("Reference creation not redone")
	//	} else {
	//		if core.GetChildBaseElementReferenceWithAncestorUri(redoneInstance, BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri, hl) != redoneReference {
	//			t.Error("Reference not restored as child of function instance")
	//		}
	//		redoneCreatedElement := uOfD.GetBaseElement(createdBaseElementReferenceIdentifier)
	//		if redoneCreatedElement == nil {
	//			t.Error("Created element not redone")
	//		} else {
	//			if redoneReference.(core.BaseElementReference).GetReferencedBaseElement(hl) != redoneCreatedElement {
	//				t.Error("Reference pointer to created element not restored")
	//			}
	//		}
	//	}
}
