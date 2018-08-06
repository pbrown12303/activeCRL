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

func TestElementFunctionsIds(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var ElementFunctionsUri string = CoreFunctionsPrefix + "ElementFunctions"
	validateElementId(t, uOfD, hl, ElementFunctionsUri)
	//
	//var ElementCreateUri string = CoreFunctionsPrefix + "Element/Create"
	validateElementId(t, uOfD, hl, ElementCreateUri)
	//var ElementCreateCreatedElementRefUri = CoreFunctionsPrefix + "Element/Create/CreatedElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementCreateCreatedElementRefUri)
	//
	//var ElementGetDefinitionUri string = CoreFunctionsPrefix + "Element/GetDefinition"
	validateElementId(t, uOfD, hl, ElementGetDefinitionUri)
	//var ElementGetDefinitionSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetDefinition/SourceElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementGetDefinitionSourceElementRefUri)
	//var ElementGetDefinitionCreatedLiteralRefUri string = CoreFunctionsPrefix + "Element/GetDefinition/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementGetDefinitionCreatedLiteralRefUri)
	//
	//var ElementGetDefinitionLiteralUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteral"
	validateElementId(t, uOfD, hl, ElementGetDefinitionLiteralUri)
	//var ElementGetDefinitionLiteralSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteral/SourceElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementGetDefinitionLiteralSourceElementRefUri)
	//var ElementGetDefinitionLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteral/IndicatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementGetDefinitionLiteralIndicatedLiteralRefUri)
	//
	//var ElementGetDefinitionLiteralPointerUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteralPointer"
	validateElementId(t, uOfD, hl, ElementGetDefinitionLiteralPointerUri)
	//var ElementGetDefinitionLiteralPointerSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteralPointer/SourceElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementGetDefinitionLiteralPointerSourceElementRefUri)
	//var ElementGetDefinitionLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteralPointer/IndicatedLiteralPointerRef"
	validateLiteralPointerReferenceId(t, uOfD, hl, ElementGetDefinitionLiteralPointerIndicatedLiteralPointerRefUri)
	//
	//var ElementGetLabelLiteralUri string = CoreFunctionsPrefix + "Element/GetLabelLiteral"
	validateElementId(t, uOfD, hl, ElementGetLabelLiteralUri)
	//var ElementGetLabelLiteralSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetLabelLiteral/SourceElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementGetLabelLiteralSourceElementRefUri)
	//var ElementGetLabelLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "Element/GetLabelLiteral/IndicatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementGetLabelLiteralIndicatedLiteralRefUri)
	//
	//var ElementGetLabelLiteralPointerUri string = CoreFunctionsPrefix + "Element/GetLabelLiteralPointer"
	validateElementId(t, uOfD, hl, ElementGetLabelLiteralPointerUri)
	//var ElementGetLabelLiteralPointerSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetLabelLiteralPointer/SourceElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementGetLabelLiteralPointerSourceElementRefUri)
	//var ElementGetLabelLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "Element/GetLabelLiteralPointer/IndicatedLiteralPointerRef"
	validateLiteralPointerReferenceId(t, uOfD, hl, ElementGetLabelLiteralPointerIndicatedLiteralPointerRefUri)
	//
	//var ElementGetUriLiteralUri string = CoreFunctionsPrefix + "Element/GetUriLiteral"
	validateElementId(t, uOfD, hl, ElementGetUriLiteralUri)
	//var ElementGetUriLiteralSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetUriLiteral/SourceElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementGetUriLiteralSourceElementRefUri)
	//var ElementGetUriLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "Element/GetUriLiteral/IndicatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementGetUriLiteralIndicatedLiteralRefUri)
	//
	//var ElementGetUriLiteralPointerUri string = CoreFunctionsPrefix + "Element/GetUriLiteralPointer"
	validateElementId(t, uOfD, hl, ElementGetUriLiteralPointerUri)
	//var ElementGetUriLiteralPointerSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetUriLiteralPointer/SourceElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementGetUriLiteralPointerSourceElementRefUri)
	//var ElementGetUriLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "Element/GetUriLiteralPointer/IndicatedLiteralPointerRef"
	validateLiteralPointerReferenceId(t, uOfD, hl, ElementGetUriLiteralPointerIndicatedLiteralPointerRefUri)
	//
	//var ElementSetDefinitionUri string = CoreFunctionsPrefix + "Element/SetDefinition"
	validateElementId(t, uOfD, hl, ElementSetDefinitionUri)
	//var ElementSetDefinitionSourceLiteralRefUri string = CoreFunctionsPrefix + "Element/SetDefinition/SourceLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementSetDefinitionSourceLiteralRefUri)
	//var ElementSetDefinitionModifiedElementRefUri string = CoreFunctionsPrefix + "Element/SetDefinition/ModifiedElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementSetDefinitionModifiedElementRefUri)
	//
	//var ElementSetLabelUri string = CoreFunctionsPrefix + "Element/SetLabel"
	validateElementId(t, uOfD, hl, ElementSetLabelUri)
	//var ElementSetLabelSourceLiteralRefUri string = CoreFunctionsPrefix + "Element/SetLabel/SourceLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, ElementSetLabelSourceLiteralRefUri)
	//var ElementSetLabelModifiedElementRefUri string = CoreFunctionsPrefix + "Element/SetLabel/ModifiedElementRef"
	validateElementReferenceId(t, uOfD, hl, ElementSetLabelModifiedElementRefUri)
}

func TestCreateElementFunction(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createElementFunction := uOfD.GetElementWithUri(ElementCreateUri)
	if createElementFunction == nil {
		t.Error("CreateElement Function not found")
	}
	createdElementReference := uOfD.GetElementReferenceWithUri(ElementCreateCreatedElementRefUri)
	if createdElementReference == nil {
		t.Error("CreatedElementReference not found")
	}

	// Now create the instance of the function
	createElementInstance := uOfD.NewElement(hl)
	createElementInstanceIdentifier := createElementInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createElementFunction, hl)

	refinementInstance.SetRefinedElement(createElementInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	//	log.Printf("Original instance:")
	//	core.Print(createElementInstance, "...", hl)

	foundReference := core.GetChildElementReferenceWithAncestorUri(createElementInstance, ElementCreateCreatedElementRefUri, hl)
	foundReferenceIdentifier := ""
	var createdElement core.Element
	createdElementIdentifier := ""
	if foundReference == nil {
		t.Error("Reference not created")
	} else {
		foundReferenceIdentifier = foundReference.GetId(hl)
		createdElement = foundReference.GetReferencedElement(hl)
		if createdElement == nil {
			t.Error("Element not created")
		} else {
			createdElementIdentifier = createdElement.GetId(hl)
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
		if core.GetChildElementReferenceWithAncestorUri(redoneInstance, ElementCreateCreatedElementRefUri, hl) != redoneReference {
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

func TestGetDefinition(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	getDefinitionFunction := uOfD.GetElementWithUri(ElementGetDefinitionUri)
	if getDefinitionFunction == nil {
		t.Error("GetDefinition Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetDefinitionSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	createdLiteralRef := uOfD.GetLiteralReferenceWithUri(ElementGetDefinitionCreatedLiteralRefUri)
	if createdLiteralRef == nil {
		t.Error("CreatedLiteralRef not found")
	}

	// Now create the instance of the function
	getDefinitionInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getDefinitionFunction, hl)
	refinementInstance.SetRefinedElement(getDefinitionInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getDefinitionInstance, ElementGetDefinitionSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not created")
	}
	foundCreatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(getDefinitionInstance, ElementGetDefinitionCreatedLiteralRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not created")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceDefinition := "SourceDefinition"
	core.SetDefinition(sourceElement, sourceDefinition, hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	createdLiteral := foundCreatedLiteralRef.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		t.Error("Literal not created")
	} else {
		if createdLiteral.GetLiteralValue(hl) != sourceDefinition {
			t.Error("Literal value not set properly")
			core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
			core.Print(createdLiteralRef, "foundCreatedLiteralRef: ", hl)
		}
	}
}

func TestGetDefinitionLiteral(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	getDefinitionFunction := uOfD.GetElementWithUri(ElementGetDefinitionLiteralUri)
	if getDefinitionFunction == nil {
		t.Error("GetDefinition Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetDefinitionLiteralSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	indicatedLiteralRef := uOfD.GetLiteralReferenceWithUri(ElementGetDefinitionLiteralIndicatedLiteralRefUri)
	if indicatedLiteralRef == nil {
		t.Error("IndicatedLiteralRef not found")
	}

	// Now create the instance of the function
	getDefinitionInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getDefinitionFunction, hl)
	refinementInstance.SetRefinedElement(getDefinitionInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getDefinitionInstance, ElementGetDefinitionLiteralSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	foundIndicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(getDefinitionInstance, ElementGetDefinitionLiteralIndicatedLiteralRefUri, hl)
	if foundIndicatedLiteralRef == nil {
		t.Error("IndicatedLiteralRef not found")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceDefinition := "SourceDefinition"
	core.SetDefinition(sourceElement, sourceDefinition, hl)
	sourceDefinitionLiteral := sourceElement.GetDefinitionLiteral(hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	indicatedLiteral := foundIndicatedLiteralRef.GetReferencedLiteral(hl)
	if indicatedLiteral != sourceDefinitionLiteral {
		t.Error("IndicatedLiteral not set properly")
		core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
		core.Print(indicatedLiteralRef, "foundIndicatedLiteralRef: ", hl)
	}
}

func TestGetDefinitionLiteralPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	getDefinitionFunction := uOfD.GetElementWithUri(ElementGetDefinitionLiteralPointerUri)
	if getDefinitionFunction == nil {
		t.Error("GetDefinition Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetDefinitionLiteralPointerSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	indicatedLiteralPointerRef := uOfD.GetLiteralPointerReferenceWithUri(ElementGetDefinitionLiteralPointerIndicatedLiteralPointerRefUri)
	if indicatedLiteralPointerRef == nil {
		t.Error("IndicatedLiteralPointerRef not found")
	}

	// Now create the instance of the function
	getDefinitionInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getDefinitionFunction, hl)
	refinementInstance.SetRefinedElement(getDefinitionInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getDefinitionInstance, ElementGetDefinitionLiteralPointerSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	foundIndicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(getDefinitionInstance, ElementGetDefinitionLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if foundIndicatedLiteralPointerRef == nil {
		t.Error("IndicatedLiteralPointerRef not found")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceDefinition := "SourceDefinition"
	core.SetDefinition(sourceElement, sourceDefinition, hl)
	sourceDefinitionLiteralPointer := sourceElement.GetDefinitionLiteralPointer(hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	indicatedLiteralPointer := foundIndicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	if indicatedLiteralPointer != sourceDefinitionLiteralPointer {
		t.Error("IndicatedLiteralPointer not set properly")
		core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
		core.Print(indicatedLiteralPointerRef, "foundIndicatedLiteralPointerRef: ", hl)
	}
}

func TestGetLabelLiteral(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	getLabelFunction := uOfD.GetElementWithUri(ElementGetLabelLiteralUri)
	if getLabelFunction == nil {
		t.Error("GetLabel Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetLabelLiteralSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	indicatedLiteralRef := uOfD.GetLiteralReferenceWithUri(ElementGetLabelLiteralIndicatedLiteralRefUri)
	if indicatedLiteralRef == nil {
		t.Error("IndicatedLiteralRef not found")
	}

	// Now create the instance of the function
	getLabelInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getLabelFunction, hl)
	refinementInstance.SetRefinedElement(getLabelInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getLabelInstance, ElementGetLabelLiteralSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	foundIndicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(getLabelInstance, ElementGetLabelLiteralIndicatedLiteralRefUri, hl)
	if foundIndicatedLiteralRef == nil {
		t.Error("IndicatedLiteralRef not found")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceLabel := "SourceLabel"
	core.SetLabel(sourceElement, sourceLabel, hl)
	sourceLabelLiteral := sourceElement.GetLabelLiteral(hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	indicatedLiteral := foundIndicatedLiteralRef.GetReferencedLiteral(hl)
	if indicatedLiteral != sourceLabelLiteral {
		t.Error("IndicatedLiteral not set properly")
		core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
		core.Print(indicatedLiteralRef, "foundIndicatedLiteralRef: ", hl)
	}
}

func TestGetLabelLiteralPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	getLabelFunction := uOfD.GetElementWithUri(ElementGetLabelLiteralPointerUri)
	if getLabelFunction == nil {
		t.Error("GetLabel Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetLabelLiteralPointerSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	indicatedLiteralPointerRef := uOfD.GetLiteralPointerReferenceWithUri(ElementGetLabelLiteralPointerIndicatedLiteralPointerRefUri)
	if indicatedLiteralPointerRef == nil {
		t.Error("IndicatedLiteralPointerRef not found")
	}

	// Now create the instance of the function
	getLabelInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getLabelFunction, hl)
	refinementInstance.SetRefinedElement(getLabelInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getLabelInstance, ElementGetLabelLiteralPointerSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	foundIndicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(getLabelInstance, ElementGetLabelLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if foundIndicatedLiteralPointerRef == nil {
		t.Error("IndicatedLiteralPointerRef not found")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceLabel := "SourceLabel"
	core.SetLabel(sourceElement, sourceLabel, hl)
	sourceLabelLiteralPointer := sourceElement.GetLabelLiteralPointer(hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	indicatedLiteralPointer := foundIndicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	if indicatedLiteralPointer != sourceLabelLiteralPointer {
		t.Error("IndicatedLiteralPointer not set properly")
		core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
		core.Print(indicatedLiteralPointerRef, "foundIndicatedLiteralPointerRef: ", hl)
	}
}

func TestGetUriLiteral(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	getUriFunction := uOfD.GetElementWithUri(ElementGetUriLiteralUri)
	if getUriFunction == nil {
		t.Error("GetUri Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetUriLiteralSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	indicatedLiteralRef := uOfD.GetLiteralReferenceWithUri(ElementGetUriLiteralIndicatedLiteralRefUri)
	if indicatedLiteralRef == nil {
		t.Error("IndicatedLiteralRef not found")
	}

	// Now create the instance of the function
	getUriInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getUriFunction, hl)
	refinementInstance.SetRefinedElement(getUriInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getUriInstance, ElementGetUriLiteralSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	foundIndicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(getUriInstance, ElementGetUriLiteralIndicatedLiteralRefUri, hl)
	if foundIndicatedLiteralRef == nil {
		t.Error("IndicatedLiteralRef not found")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceUri := "SourceUri"
	core.SetUri(sourceElement, sourceUri, hl)
	sourceUriLiteral := sourceElement.GetUriLiteral(hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	indicatedLiteral := foundIndicatedLiteralRef.GetReferencedLiteral(hl)
	if indicatedLiteral != sourceUriLiteral {
		t.Error("IndicatedLiteral not set properly")
		core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
		core.Print(indicatedLiteralRef, "foundIndicatedLiteralRef: ", hl)
	}
}

func TestGetUriLiteralPointer(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	getUriFunction := uOfD.GetElementWithUri(ElementGetUriLiteralPointerUri)
	if getUriFunction == nil {
		t.Error("GetUri Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetUriLiteralPointerSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	indicatedLiteralPointerRef := uOfD.GetLiteralPointerReferenceWithUri(ElementGetUriLiteralPointerIndicatedLiteralPointerRefUri)
	if indicatedLiteralPointerRef == nil {
		t.Error("IndicatedLiteralPointerRef not found")
	}

	// Now create the instance of the function
	getUriInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getUriFunction, hl)
	refinementInstance.SetRefinedElement(getUriInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getUriInstance, ElementGetUriLiteralPointerSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	foundIndicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(getUriInstance, ElementGetUriLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if foundIndicatedLiteralPointerRef == nil {
		t.Error("IndicatedLiteralPointerRef not found")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceUri := "SourceUri"
	core.SetUri(sourceElement, sourceUri, hl)
	sourceUriLiteralPointer := sourceElement.GetUriLiteralPointer(hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	indicatedLiteralPointer := foundIndicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	if indicatedLiteralPointer != sourceUriLiteralPointer {
		t.Error("IndicatedLiteralPointer not set properly")
		core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
		core.Print(indicatedLiteralPointerRef, "foundIndicatedLiteralPointerRef: ", hl)
	}
}

func TestSetDefinition(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setDefinition := uOfD.GetElementWithUri(ElementSetDefinitionUri)
	if setDefinition == nil {
		t.Errorf("SetDefinition function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setDefinition, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, setDefinition, hl) != true {
		t.Errorf("Replicate is not refinement of SetDefinition()")
	}
	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementSetDefinitionSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		t.Errorf("SourceLiteralRef child not found")
	}
	modifiedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementSetDefinitionModifiedElementRefUri, hl)
	if modifiedElementRef == nil {
		t.Errorf("ModifiedElementRef child not found")
		core.Print(replicate, "Replicate: ", hl)
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	uri := "TestDefinition"
	sourceLiteral.SetLiteralValue(uri, hl)
	sourceLiteralRef.SetReferencedLiteral(sourceLiteral, hl)
	modifiedElement := uOfD.NewElement(hl)
	modifiedElementRef.SetReferencedElement(modifiedElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if modifiedElement.GetDefinition(hl) != uri {
		t.Errorf("Definition not set properly")
	}
}

func TestSetLabel(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setLabel := uOfD.GetElementWithUri(ElementSetLabelUri)
	if setLabel == nil {
		t.Errorf("SetLabel function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setLabel, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, setLabel, hl) != true {
		t.Errorf("Replicate is not refinement of SetLabel()")
	}
	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementSetLabelSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		t.Errorf("SourceLiteralRef child not found")
	}
	modifiedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementSetLabelModifiedElementRefUri, hl)
	if modifiedElementRef == nil {
		t.Errorf("ModifiedElementRef child not found")
		core.Print(replicate, "Replicate: ", hl)
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	name := "TestLabel"
	sourceLiteral.SetLiteralValue(name, hl)
	sourceLiteralRef.SetReferencedLiteral(sourceLiteral, hl)
	modifiedElement := uOfD.NewElement(hl)
	modifiedElementRef.SetReferencedElement(modifiedElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if core.GetLabel(modifiedElement, hl) != name {
		t.Errorf("Label not set properly")
	}
}
