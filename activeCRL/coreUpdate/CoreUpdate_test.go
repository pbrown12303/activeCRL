package main

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"testing"
)

func TestUpdateCoreElement(t *testing.T) {
	log.Printf("Entering TestUpdateCoreElement")
	uOfD := core.NewUniverseOfDiscourse()
	uOfD.SetRecordingUndo(false)
	var emptyCore core.Element
	core.Print(emptyCore, "")

	//Core
	recoveredCore := updateRecoveredCore(emptyCore, uOfD)
	if recoveredCore == nil {
		t.Error("updateRecoveredCore returned empty element")
	}
	if recoveredCore.GetUri() != core.CoreUri {
		t.Error("Core uri not set")
	}
	_, ok := recoveredCore.(core.Element)
	if !ok {
		t.Error("Core is of wrong type")
	}

	recoveredBaseElement := uOfD.GetBaseElementWithUri(core.ElememtUri)
	if recoveredBaseElement == nil {
		t.Error("Element not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("Element is of wrong type")
	}

	// ElementPointer
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.ElementPointerUri)
	if recoveredBaseElement == nil {
		t.Error("ElementPointer not found")
	}
	_, ok = recoveredBaseElement.(core.ElementPointer)
	if !ok {
		t.Error("ElementPointer is of wrong type")
	}

	// ElementPointerPointer
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.ElementPointerPointerUri)
	if recoveredBaseElement == nil {
		t.Error("ElementPointerPointer not found")
	}
	_, ok = recoveredBaseElement.(core.ElementPointerPointer)
	if !ok {
		t.Error("ElementPointerPointer is of wrong type")
	}

	// ElementPointerReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.ElementPointerReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("ElementPointerReference not found")
	}
	_, ok = recoveredBaseElement.(core.ElementPointerReference)
	if !ok {
		t.Error("ElementPointerReference is of wrong type")
	}

	// ElementReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.ElementReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("ElementReference not found")
	}
	_, ok = recoveredBaseElement.(core.ElementReference)
	if !ok {
		t.Error("ElementReference is of wrong type")
	}

	// Literal
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.LiteralUri)
	if recoveredBaseElement == nil {
		t.Error("Literal not found")
	}
	_, ok = recoveredBaseElement.(core.Literal)
	if !ok {
		t.Error("Literal is of wrong type")
	}

	// LiteralPointer
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.LiteralPointerUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralPointer not found")
	}
	_, ok = recoveredBaseElement.(core.LiteralPointer)
	if !ok {
		t.Error("LiteralPointer is of wrong type")
	}

	// LiteralPointerPointer
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.LiteralPointerPointerUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralPointerPointer not found")
	}
	_, ok = recoveredBaseElement.(core.LiteralPointerPointer)
	if !ok {
		t.Error("LiteralPointerPointer is of wrong type")
	}

	// LiteralPointerReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.LiteralPointerReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralPointerReference not found")
	}
	_, ok = recoveredBaseElement.(core.LiteralPointerReference)
	if !ok {
		t.Error("LiteralPointerReference is of wrong type")
	}

	// LiteralReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.LiteralReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralReference not found")
	}
	_, ok = recoveredBaseElement.(core.LiteralReference)
	if !ok {
		t.Error("LiteralReference is of wrong type")
	}

	// Refinement
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.RefinementUri)
	if recoveredBaseElement == nil {
		t.Error("Refinement not found")
	}
	_, ok = recoveredBaseElement.(core.Refinement)
	if !ok {
		t.Error("Refinement is of wrong type")
	}
}
