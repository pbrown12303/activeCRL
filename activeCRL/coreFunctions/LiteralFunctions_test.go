// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/satori/go.uuid"
	//	"log"
	"sync"
	"testing"
	//	"time"
)

func TestLiteralFunctionsIds(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	//var LiteralFunctionsUri string = CoreFunctionsPrefix + "LiteralFunctions"
	validateElementId(t, uOfD, hl, LiteralFunctionsUri)
	//
	//var LiteralCreateUri string = CoreFunctionsPrefix + "Literal/Create"
	validateElementId(t, uOfD, hl, LiteralCreateUri)
	//var LiteralCreateCreatedLiteralRefUri string = CoreFunctionsPrefix + "Literal/Create/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralCreateCreatedLiteralRefUri)
	//
	//var LiteralGetLiteralValueUri string = CoreFunctionsPrefix + "Literal/GetLiteralValue"
	validateElementId(t, uOfD, hl, LiteralGetLiteralValueUri)
	//var LiteralGetLiteralValueSourceLiteralRefUri string = CoreFunctionsPrefix + "Literal/GetLiteralValue/SourceLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralGetLiteralValueSourceLiteralRefUri)
	//var LiteralGetLiteralValueCreatedLiteralRefUri string = CoreFunctionsPrefix + "Literal/GetLiteralValue/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralGetLiteralValueCreatedLiteralRefUri)
	//
	//var LiteralSetLiteralValueUri string = CoreFunctionsPrefix + "Literal/SetLiteralValue"
	validateElementId(t, uOfD, hl, LiteralSetLiteralValueUri)
	//var LiteralSetLiteralValueSourceLiteralRefUri string = CoreFunctionsPrefix + "Literal/SetLiteralValue/SourceLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralSetLiteralValueSourceLiteralRefUri)
	//var LiteralSetLiteralValueModifiedLiteralRefUri string = CoreFunctionsPrefix + "Literal/SetLiteralValue/ModifiedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, LiteralSetLiteralValueModifiedLiteralRefUri)
}

func TestCreateLiteralFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	createLiteralFunction := uOfD.GetElementWithUri(LiteralCreateUri)
	if createLiteralFunction == nil {
		t.Error("CreateLiteral Function not found")
	}
	createdLiteralReference := uOfD.GetLiteralReferenceWithUri(LiteralCreateCreatedLiteralRefUri)
	if createdLiteralReference == nil {
		t.Error("CreatedLiteralReference not found")
	}

	// Now create the instance of the function
	createLiteralInstance := uOfD.NewElement(hl)
	createLiteralInstanceIdentifier := createLiteralInstance.GetId(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(createLiteralFunction, hl)

	refinementInstance.SetRefinedElement(createLiteralInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	//	log.Printf("Original instance:")
	//	core.Print(createLiteralInstance, "...", hl)

	foundReference := core.GetChildLiteralReferenceWithAncestorUri(createLiteralInstance, LiteralCreateCreatedLiteralRefUri, hl)
	foundReferenceIdentifier := uuid.Nil
	var createdLiteral core.Literal
	createdLiteralIdentifier := uuid.Nil
	if foundReference == nil {
		t.Error("Reference not created")
	} else {
		foundReferenceIdentifier = foundReference.GetId(hl)
		createdLiteral = foundReference.GetReferencedLiteral(hl)
		if createdLiteral == nil {
			t.Error("Literal not created")
		} else {
			createdLiteralIdentifier = createdLiteral.GetId(hl)
		}
	}

	// Now undo
	uOfD.Undo(hl)
	if uOfD.GetElement(createLiteralInstanceIdentifier) != nil {
		t.Error("Literal creation not undone")
	}
	if uOfD.GetElement(foundReferenceIdentifier) != nil {
		t.Error("Literal creation not undone")
	}
	if uOfD.GetLiteral(createdLiteralIdentifier) != nil {
		t.Error("Literal creation not undone")
	}

	// Now Redo
	uOfD.Redo(hl)
	redoneInstance := uOfD.GetElement(createLiteralInstanceIdentifier)
	if redoneInstance == nil {
		t.Error("Literal creation not redone")
	}
	redoneReference := uOfD.GetElement(foundReferenceIdentifier)
	if redoneReference == nil {
		t.Error("Reference creation not redone")
	} else {
		if core.GetChildLiteralReferenceWithAncestorUri(redoneInstance, LiteralCreateCreatedLiteralRefUri, hl) != redoneReference {
			t.Error("Reference not restored as child of function instance")
		}
		redoneCreatedLiteral := uOfD.GetLiteral(createdLiteralIdentifier)
		if redoneCreatedLiteral == nil {
			t.Error("Created literal not redone")
		} else {
			if redoneReference.(core.LiteralReference).GetReferencedLiteral(hl) != redoneCreatedLiteral {
				t.Error("Reference pointer to created literal not restored")
			}
		}
	}
}

func TestGetLiteralValue(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get the reference elements
	getLiteralValueFunction := uOfD.GetElementWithUri(LiteralGetLiteralValueUri)
	if getLiteralValueFunction == nil {
		t.Error("GetLiteralValue Function not found")
	}
	sourceLiteralRef := uOfD.GetLiteralReferenceWithUri(LiteralGetLiteralValueSourceLiteralRefUri)
	if sourceLiteralRef == nil {
		t.Error("SourceLiteralRef not found")
	}
	createdLiteralRef := uOfD.GetLiteralReferenceWithUri(LiteralGetLiteralValueCreatedLiteralRefUri)
	if createdLiteralRef == nil {
		t.Error("CreatedLiteralRef not found")
	}

	// Now create the instance of the function
	getLiteralValueInstance := uOfD.NewElement(hl)
	refinementInstance := uOfD.NewRefinement(hl)
	refinementInstance.SetAbstractElement(getLiteralValueFunction, hl)
	refinementInstance.SetRefinedElement(getLiteralValueInstance, hl)
	hl.ReleaseLocks()
	wg.Wait()

	// Check the results
	foundSourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(getLiteralValueInstance, LiteralGetLiteralValueSourceLiteralRefUri, hl)
	if foundSourceLiteralRef == nil {
		t.Error("SourceLiteralRef not created")
	}
	foundCreatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(getLiteralValueInstance, LiteralGetLiteralValueCreatedLiteralRefUri, hl)
	if foundSourceLiteralRef == nil {
		t.Error("SourceLiteralRef not created")
	}

	// Now check function execution
	sourceLiteral := uOfD.NewLiteral(hl)
	sourceLiteralValue := "SourceLiteralValue"
	sourceLiteral.SetLiteralValue(sourceLiteralValue, hl)
	foundSourceLiteralRef.SetReferencedLiteral(sourceLiteral, hl)
	hl.ReleaseLocks()
	wg.Wait()

	createdLiteral := foundCreatedLiteralRef.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		t.Error("Literal not created")
	} else {
		if createdLiteral.GetLiteralValue(hl) != sourceLiteralValue {
			t.Error("Literal value not set properly")
			core.Print(sourceLiteralRef, "foundSourceLiteralRef: ", hl)
			core.Print(createdLiteralRef, "foundCreatedLiteralRef: ", hl)
		}
	}
}

func TestSetLiteralValue(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setLiteralValue := uOfD.GetElementWithUri(LiteralSetLiteralValueUri)
	if setLiteralValue == nil {
		t.Errorf("SetLiteralValue function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setLiteralValue, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	// Now check the replication
	if replicate.IsRefinementOf(setLiteralValue, hl) != true {
		t.Errorf("Replicate is not refinement of SetLiteralValue()")
	}
	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralSetLiteralValueSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		t.Errorf("SourceLiteralRef child not found")
	}
	modifiedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralSetLiteralValueModifiedLiteralRefUri, hl)
	if modifiedLiteralRef == nil {
		t.Errorf("ModifiedLiteralRef child not found")
		core.Print(replicate, "Replicate: ", hl)
	}

	// Now test target reference update functionality
	sourceLiteral := uOfD.NewLiteral(hl)
	literalValue := "TestLiteralValue"
	sourceLiteral.SetLiteralValue(literalValue, hl)
	sourceLiteralRef.SetReferencedLiteral(sourceLiteral, hl)
	modifiedLiteral := uOfD.NewLiteral(hl)
	modifiedLiteralRef.SetReferencedLiteral(modifiedLiteral, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()

	hl.LockBaseElement(replicate)
	if modifiedLiteral.GetLiteralValue(hl) != literalValue {
		t.Errorf("LiteralValue not set properly")
	}
}
