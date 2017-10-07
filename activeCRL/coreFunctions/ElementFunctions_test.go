package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"sync"
	"testing"
	//	"time"
)

func TestCreateElementFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

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
	createElementInstanceIdentifier := createElementInstance.GetId(hl).String()
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
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

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
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

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
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

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

func TestGetNameLiteral(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	getNameFunction := uOfD.GetElementWithUri(ElementGetNameLiteralUri)
	if getNameFunction == nil {
		t.Error("GetName Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetNameLiteralSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	indicatedLiteralRef := uOfD.GetLiteralReferenceWithUri(ElementGetNameLiteralIndicatedLiteralRefUri)
	if indicatedLiteralRef == nil {
		t.Error("IndicatedLiteralRef not found")
	}

	// Now create the instance of the function
	getNameInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getNameFunction, hl)
	refinementInstance.SetRefinedElement(getNameInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getNameInstance, ElementGetNameLiteralSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	foundIndicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(getNameInstance, ElementGetNameLiteralIndicatedLiteralRefUri, hl)
	if foundIndicatedLiteralRef == nil {
		t.Error("IndicatedLiteralRef not found")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceName := "SourceName"
	core.SetName(sourceElement, sourceName, hl)
	sourceNameLiteral := sourceElement.GetNameLiteral(hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	indicatedLiteral := foundIndicatedLiteralRef.GetReferencedLiteral(hl)
	if indicatedLiteral != sourceNameLiteral {
		t.Error("IndicatedLiteral not set properly")
		core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
		core.Print(indicatedLiteralRef, "foundIndicatedLiteralRef: ", hl)
	}
}

func TestGetNameLiteralPointer(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get the reference elements
	getNameFunction := uOfD.GetElementWithUri(ElementGetNameLiteralPointerUri)
	if getNameFunction == nil {
		t.Error("GetName Function not found")
	}
	sourceElementRef := uOfD.GetElementReferenceWithUri(ElementGetNameLiteralPointerSourceElementRefUri)
	if sourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	indicatedLiteralPointerRef := uOfD.GetLiteralPointerReferenceWithUri(ElementGetNameLiteralPointerIndicatedLiteralPointerRefUri)
	if indicatedLiteralPointerRef == nil {
		t.Error("IndicatedLiteralPointerRef not found")
	}

	// Now create the instance of the function
	getNameInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getNameFunction, hl)
	refinementInstance.SetRefinedElement(getNameInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceElementRef := core.GetChildElementReferenceWithAncestorUri(getNameInstance, ElementGetNameLiteralPointerSourceElementRefUri, hl)
	if foundSourceElementRef == nil {
		t.Error("SourceElementRef not found")
	}
	foundIndicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(getNameInstance, ElementGetNameLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if foundIndicatedLiteralPointerRef == nil {
		t.Error("IndicatedLiteralPointerRef not found")
	}

	// Now check function execution
	sourceElement := uOfD.NewElement(hl)
	sourceName := "SourceName"
	core.SetName(sourceElement, sourceName, hl)
	sourceNameLiteralPointer := sourceElement.GetNameLiteralPointer(hl)
	foundSourceElementRef.SetReferencedElement(sourceElement, hl)
	hl.ReleaseLocks()
	wg.Wait()

	indicatedLiteralPointer := foundIndicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	if indicatedLiteralPointer != sourceNameLiteralPointer {
		t.Error("IndicatedLiteralPointer not set properly")
		core.Print(sourceElementRef, "foundSourceElementRef: ", hl)
		core.Print(indicatedLiteralPointerRef, "foundIndicatedLiteralPointerRef: ", hl)
	}
}

func TestGetUriLiteral(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

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
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

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
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

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
	if replicate.IsRefinementOf(setDefinition, hl) != true {
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

func TestSetName(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	GetCoreFunctionsConceptSpace(uOfD)

	// Get Ancestor
	setName := uOfD.GetElementWithUri(ElementSetNameUri)
	if setName == nil {
		t.Errorf("SetName function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setName, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setName, hl) != true {
		t.Errorf("Replicate is not refinement of SetName()")
	}
	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementSetNameSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		t.Errorf("SourceLiteralRef child not found")
	}
	modifiedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementSetNameModifiedElementRefUri, hl)
	if modifiedElementRef == nil {
		t.Errorf("ModifiedElementRef child not found")
		core.Print(replicate, "Replicate: ", hl)
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	name := "TestName"
	sourceLiteral.SetLiteralValue(name, hl)
	sourceLiteralRef.SetReferencedLiteral(sourceLiteral, hl)
	modifiedElement := uOfD.NewElement(hl)
	modifiedElementRef.SetReferencedElement(modifiedElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if core.GetName(modifiedElement, hl) != name {
		t.Errorf("Name not set properly")
	}
}
