// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"sync"
	"testing"
)

func TestUpdateCoreElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(false)
	var emptyCore core.Element
	//	core.Print(emptyCore, "", hl)

	//Core
	recoveredCore := updateRecoveredCore(emptyCore, uOfD, hl)
	//	core.Print(recoveredCore, "", hl)
	if recoveredCore == nil {
		t.Error("updateRecoveredCore returned empty element")
	}
	if core.GetUri(recoveredCore, hl) != core.CoreConceptSpaceUri {
		t.Error("Core uri not set")
	}
	_, ok := recoveredCore.(core.Element)
	if !ok {
		t.Error("Core is of wrong type")
	}
	if uOfD.GetBaseElementWithUri(core.CoreConceptSpaceUri) == nil {
		t.Error("UofD uri index not updated")
	}

	// BaseElementPointer
	recoveredBaseElement := uOfD.GetBaseElementWithUri(core.BaseElementPointerUri)
	if recoveredBaseElement == nil {
		t.Error("BaseElementPointer not found")
	}
	_, ok = recoveredBaseElement.(core.BaseElementPointer)
	if !ok {
		t.Error("ElementPointer is of wrong type")
	}

	// BaseElementReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.BaseElementReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("BaseElementReference not found")
	}
	_, ok = recoveredBaseElement.(core.BaseElementReference)
	if !ok {
		t.Error("ElementPointer is of wrong type")
	}

	// Element
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.ElememtUri)
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
		core.Print(recoveredBaseElement, "", hl)
		//		core.Print(recoveredCore, "", hl)
		//		core.PrintUriIndex(uOfD, hl)
	}

	// ElementPointerRole and values
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.ElementPointerRoleUri)
	if recoveredBaseElement == nil {
		t.Error("ElementPointerRole not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("ElementPointerRole is of wrong type")
		core.Print(recoveredBaseElement, "", hl)
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.AbstractElementUri)
	if recoveredBaseElement == nil {
		t.Error("AbstractElement not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("AbstractElement is of wrong type")
		core.Print(recoveredBaseElement, "", hl)
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.RefinedElementUri)
	if recoveredBaseElement == nil {
		t.Error("RefinedElement not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("RefinedElement is of wrong type")
		core.Print(recoveredBaseElement, "", hl)
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.OwningElementUri)
	if recoveredBaseElement == nil {
		t.Error("OwningElement not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("OwningElement is of wrong type")
		core.Print(recoveredBaseElement, "", hl)
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.ReferencedElementUri)
	if recoveredBaseElement == nil {
		t.Error("ReferencedElement not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("ReferencedElement is of wrong type")
		core.Print(recoveredBaseElement, "", hl)
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

	// LiteralPointerRole and values
	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.LiteralPointerRoleUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralPointerRole not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("LiteralPointerRole is of wrong type")
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.NameUri)
	if recoveredBaseElement == nil {
		t.Error("Name not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("Name is of wrong type")
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.DefinitionUri)
	if recoveredBaseElement == nil {
		t.Error("Definition not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("Definition is of wrong type")
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.UriUri)
	if recoveredBaseElement == nil {
		t.Error("Uri not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("Uri is of wrong type")
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(core.ValueUri)
	if recoveredBaseElement == nil {
		t.Error("Value not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("Value is of wrong type")
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
