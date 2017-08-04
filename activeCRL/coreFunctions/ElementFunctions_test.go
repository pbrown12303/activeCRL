package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"testing"
	"time"
)

func TestCreateElementFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	createElement := uOfD.GetElementWithUri(CreateElememtUri)
	if createElement == nil {
		t.Error("CreateElement not found")
	}
	createdElementReference := uOfD.GetElementReferenceWithUri(CreatedElementReferenceUri)
	if createdElementReference == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createElementInstance := uOfD.NewElement(hl)
	createElementInstanceIdentifier := createElementInstance.GetId(hl).String()
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createElement, hl)

	refinementInstance.SetRefinedElement(createElementInstance, hl)
	hl.ReleaseLocks()
	time.Sleep(10000000 * time.Nanosecond)

	// Check the results
	//	log.Printf("Original instance:")
	//	core.Print(createElementInstance, "...", hl)

	foundReference := core.GetChildElementReferenceWithAncestorUri(createElementInstance, CreatedElementReferenceUri, hl)
	foundReferenceIdentifier := ""
	var createdElement core.Element
	createdElementIdentifier := ""
	if foundReference == nil {
		t.Error("Reference not created")
	} else {
		foundReferenceIdentifier = foundReference.GetId(hl).String()
		createdElement = foundReference.GetReferencedElement(hl)
		if createdElement == nil {
			t.Error("Element not created")
		} else {
			createdElementIdentifier = createdElement.GetId(hl).String()
		}
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createElementInstanceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(foundReferenceIdentifier) != nil {
		t.Error("Element creation not undone")
	}
	if uOfD.GetElement(createdElementIdentifier) != nil {
		t.Error("Element creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createElementInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Element creation not redone")
	}
	redoneReference := uOfD.GetElement(foundReferenceIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, CreatedElementReferenceUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedElement := uOfD.GetElement(createdElementIdentifier)
		if redoneCreatedElement == nil {
			t.Error("Created element not redone")
		} else {
			if redoneReference.(core.ElementReference).GetReferencedElement(hl) != redoneCreatedElement {
				t.Error("Reference pointer to created element not restored")
			}
		}
	}
}
