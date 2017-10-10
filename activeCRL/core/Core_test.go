package core

import (
	"sync"
	"testing"
)

func TestBuildCoreConceptSpace(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(false)

	//Core
	builtCore := buildCoreConceptSpace(uOfD, hl)
	//	Print(builtCore, "", hl)
	if builtCore == nil {
		t.Error("buildCoreConceptSpace returned empty element")
	}
	if GetUri(builtCore, hl) != CoreConceptSpaceUri {
		t.Error("Core uri not set")
	}
	_, ok := builtCore.(Element)
	if !ok {
		t.Error("Core is of wrong type")
	}
	if uOfD.GetBaseElementWithUri(CoreConceptSpaceUri) == nil {
		t.Error("UofD uri index not updated")
	}

	// BaseElementPointer
	recoveredBaseElement := uOfD.GetBaseElementWithUri(BaseElementPointerUri)
	if recoveredBaseElement == nil {
		t.Error("BaseElementPointer not found")
	}
	_, ok = recoveredBaseElement.(BaseElementPointer)
	if !ok {
		t.Error("ElementPointer is of wrong type")
	}

	// BaseElementReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(BaseElementReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("BaseElementReference not found")
	}
	_, ok = recoveredBaseElement.(BaseElementReference)
	if !ok {
		t.Error("ElementPointer is of wrong type")
	}

	// Element
	recoveredBaseElement = uOfD.GetBaseElementWithUri(ElememtUri)
	if recoveredBaseElement == nil {
		t.Error("Element not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("Element is of wrong type")
	}

	// ElementPointer
	recoveredBaseElement = uOfD.GetBaseElementWithUri(ElementPointerUri)
	if recoveredBaseElement == nil {
		t.Error("ElementPointer not found")
	}
	_, ok = recoveredBaseElement.(ElementPointer)
	if !ok {
		t.Error("ElementPointer is of wrong type")
	}

	// ElementPointerPointer
	recoveredBaseElement = uOfD.GetBaseElementWithUri(ElementPointerPointerUri)
	if recoveredBaseElement == nil {
		t.Error("ElementPointerPointer not found")
	}
	_, ok = recoveredBaseElement.(ElementPointerPointer)
	if !ok {
		t.Error("ElementPointerPointer is of wrong type")
	}

	// ElementPointerReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(ElementPointerReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("ElementPointerReference not found")
	}
	_, ok = recoveredBaseElement.(ElementPointerReference)
	if !ok {
		t.Error("ElementPointerReference is of wrong type")
		Print(recoveredBaseElement, "", hl)
		//		Print(builtCore, "", hl)
		//		PrintUriIndex(uOfD, hl)
	}

	// ElementPointerRole and values
	recoveredBaseElement = uOfD.GetBaseElementWithUri(ElementPointerRoleUri)
	if recoveredBaseElement == nil {
		t.Error("ElementPointerRole not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("ElementPointerRole is of wrong type")
		Print(recoveredBaseElement, "", hl)
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(AbstractElementUri)
	if recoveredBaseElement == nil {
		t.Error("AbstractElement not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("AbstractElement is of wrong type")
		Print(recoveredBaseElement, "", hl)
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(RefinedElementUri)
	if recoveredBaseElement == nil {
		t.Error("RefinedElement not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("RefinedElement is of wrong type")
		Print(recoveredBaseElement, "", hl)
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(OwningElementUri)
	if recoveredBaseElement == nil {
		t.Error("OwningElement not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("OwningElement is of wrong type")
		Print(recoveredBaseElement, "", hl)
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(ReferencedElementUri)
	if recoveredBaseElement == nil {
		t.Error("ReferencedElement not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("ReferencedElement is of wrong type")
		Print(recoveredBaseElement, "", hl)
	}

	// ElementReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(ElementReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("ElementReference not found")
	}
	_, ok = recoveredBaseElement.(ElementReference)
	if !ok {
		t.Error("ElementReference is of wrong type")
	}

	// Literal
	recoveredBaseElement = uOfD.GetBaseElementWithUri(LiteralUri)
	if recoveredBaseElement == nil {
		t.Error("Literal not found")
	}
	_, ok = recoveredBaseElement.(Literal)
	if !ok {
		t.Error("Literal is of wrong type")
	}

	// LiteralPointer
	recoveredBaseElement = uOfD.GetBaseElementWithUri(LiteralPointerUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralPointer not found")
	}
	_, ok = recoveredBaseElement.(LiteralPointer)
	if !ok {
		t.Error("LiteralPointer is of wrong type")
	}

	// LiteralPointerPointer
	recoveredBaseElement = uOfD.GetBaseElementWithUri(LiteralPointerPointerUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralPointerPointer not found")
	}
	_, ok = recoveredBaseElement.(LiteralPointerPointer)
	if !ok {
		t.Error("LiteralPointerPointer is of wrong type")
	}

	// LiteralPointerReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(LiteralPointerReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralPointerReference not found")
	}
	_, ok = recoveredBaseElement.(LiteralPointerReference)
	if !ok {
		t.Error("LiteralPointerReference is of wrong type")
	}

	// LiteralPointerRole and values
	recoveredBaseElement = uOfD.GetBaseElementWithUri(LiteralPointerRoleUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralPointerRole not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("LiteralPointerRole is of wrong type")
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(NameUri)
	if recoveredBaseElement == nil {
		t.Error("Name not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("Name is of wrong type")
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(DefinitionUri)
	if recoveredBaseElement == nil {
		t.Error("Definition not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("Definition is of wrong type")
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(UriUri)
	if recoveredBaseElement == nil {
		t.Error("Uri not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("Uri is of wrong type")
	}

	recoveredBaseElement = uOfD.GetBaseElementWithUri(ValueUri)
	if recoveredBaseElement == nil {
		t.Error("Value not found")
	}
	_, ok = recoveredBaseElement.(Element)
	if !ok {
		t.Error("Value is of wrong type")
	}

	// LiteralReference
	recoveredBaseElement = uOfD.GetBaseElementWithUri(LiteralReferenceUri)
	if recoveredBaseElement == nil {
		t.Error("LiteralReference not found")
	}
	_, ok = recoveredBaseElement.(LiteralReference)
	if !ok {
		t.Error("LiteralReference is of wrong type")
	}

	// Refinement
	recoveredBaseElement = uOfD.GetBaseElementWithUri(RefinementUri)
	if recoveredBaseElement == nil {
		t.Error("Refinement not found")
	}
	_, ok = recoveredBaseElement.(Refinement)
	if !ok {
		t.Error("Refinement is of wrong type")
	}
}
