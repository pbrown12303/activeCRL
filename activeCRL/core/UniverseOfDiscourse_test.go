package core

import (
	"testing"
)

func TestGetBaseElementWithUri(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()

	// Element
	element := uOfD.NewElement(hl)
	SetName(element, "Element", hl)
	recoveredElement := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Element")
	if recoveredElement != nil {
		t.Error("Wrong element returned for find Element by URI")
		Print(recoveredElement, "", hl)
	}
	SetUri(element, "http://activeCrl.com/test/Element", hl)
	recoveredElement = uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Element")
	if recoveredElement == nil {
		t.Error("Did not find Element by URI")
	}

	// ElementPointer
	elementPointer := uOfD.NewReferencedElementPointer(hl)
	elementPointer.SetUri("http://activeCrl.com/test/ElementPointer", hl)
	recoveredElementPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/ElementPointer")
	if recoveredElementPointer == nil {
		t.Error("Did not find ElementPointer by URI")
	}

	// ElementPointerPointer
	elementPointerPointer := uOfD.NewElementPointerPointer(hl)
	elementPointerPointer.SetUri("http://activeCrl.com/test/ElementPointerPointer", hl)
	recoveredElementPointerPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/ElementPointerPointer")
	if recoveredElementPointerPointer == nil {
		t.Error("Did not find ElementPointerPointer by URI")
	}

	// ElementPointerReference
	elementPointerReference := uOfD.NewElementPointerReference(hl)
	SetName(elementPointerReference, "ElementReference", hl)
	SetUri(elementPointerReference, "http://activeCrl.com/test/ElementPointerReference", hl)
	recoveredElementPointerReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/ElementPointerReference")
	if recoveredElementPointerReference == nil {
		t.Error("Did not find ElementPointerReference by URI")
	}

	// ElementReference
	elementReference := uOfD.NewElementReference(hl)
	SetName(elementReference, "ElementReference", hl)
	SetUri(elementReference, "http://activeCrl.com/test/ElementReference", hl)
	recoveredElementReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/ElementReference")
	if recoveredElementReference == nil {
		t.Error("Did not find ElementReference by URI")
	}

	// Literal
	literal := uOfD.NewLiteral(hl)
	literal.SetUri("http://activeCrl.com/test/Literal", hl)
	recoveredLiteral := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Literal")
	if recoveredLiteral == nil {
		t.Error("Did not find Literal by URI")
	}

	// LiteralPointer
	literalPointer := uOfD.NewValueLiteralPointer(hl)
	literalPointer.SetUri("http://activeCrl.com/test/LiteralPointer", hl)
	recoveredLiteralPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/LiteralPointer")
	if recoveredLiteralPointer == nil {
		t.Error("Did not find LiteralPointer by URI")
	}

	// LiteralPointerPointer
	literalPointerPointer := uOfD.NewLiteralPointerPointer(hl)
	literalPointerPointer.SetUri("http://activeCrl.com/test/LiteralPointerPointer", hl)
	recoveredLiteralPointerPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/LiteralPointerPointer")
	if recoveredLiteralPointerPointer == nil {
		t.Error("Did not find LiteralPointerPointer by URI")
	}

	// LiteralPointerReference
	literalPointerReference := uOfD.NewLiteralPointerReference(hl)
	SetName(literalPointerReference, "LiteralReference", hl)
	SetUri(literalPointerReference, "http://activeCrl.com/test/LiteralPointerReference", hl)
	recoveredLiteralPointerReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/LiteralPointerReference")
	if recoveredLiteralPointerReference == nil {
		t.Error("Did not find LiteralPointerReference by URI")
	}

	// LiteralReference
	literalReference := uOfD.NewLiteralReference(hl)
	SetName(literalReference, "LiteralReference", hl)
	SetUri(literalReference, "http://activeCrl.com/test/LiteralReference", hl)
	recoveredLiteralReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/LiteralReference")
	if recoveredLiteralReference == nil {
		t.Error("Did not find LiteralReference by URI")
	}

	// Refinement
	refinement := uOfD.NewRefinement(hl)
	SetName(refinement, "Refinement", hl)
	SetUri(refinement, "http://activeCrl.com/test/Refinement", hl)
	recoveredRefinement := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Refinement")
	if recoveredRefinement == nil {
		t.Error("Did not find Refinement by URI")
	}

	// Child of element
	child := uOfD.NewElement(hl)
	SetName(child, "Child", hl)
	SetOwningElement(child, element, hl)
	SetUri(child, "http://activeCrl.com/test/Element/Child", hl)
	recoveredChild := uOfD.GetBaseElementWithUri("http://activeCrl.com/test/Element/Child")
	if recoveredChild == nil {
		t.Error("Did not find Child by URI")
	}

}

func TestAddElementListener(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	e1 := uOfD.NewElement(hl)
	ep := uOfD.NewReferencedElementPointer(hl)
	ep.SetElement(e1, hl)
	elm := uOfD.elementListenerMap.GetEntry(e1.GetId(hl).String())
	if elm == nil {
		t.Error("ElementListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("ElementListenerMap entry length != 1")
		} else {
			if (*elm)[0] != ep {
				t.Error("ElementListenerMap entry does not contain ElementPointer")
			}
		}
	}
	ep.SetElement(nil, hl)
	elm = uOfD.elementListenerMap.GetEntry(e1.GetId(hl).String())
	if elm == nil {
		t.Error("ElementListenerMap entry is nil after SetElement(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("ElementListenerMap entry length != 0")
		}
	}

}

func TestAddElementPointerListener(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	ep := uOfD.NewReferencedElementPointer(hl)
	epp := uOfD.NewElementPointerPointer(hl)
	epp.SetElementPointer(ep, hl)
	elm := uOfD.elementPointerListenerMap.GetEntry(ep.GetId(hl).String())
	if elm == nil {
		t.Error("ElementPointerListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("ElementPointerListenerMap entry length != 1")
		} else {
			if (*elm)[0] != epp {
				t.Error("ElementPointerListenerMap entry does not contain ElementPointer")
			}
		}
	}
	epp.SetElementPointer(nil, hl)
	elm = uOfD.elementPointerListenerMap.GetEntry(ep.GetId(hl).String())
	if elm == nil {
		t.Error("ElementListenerMap entry is nil after SetElement(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("ElementListenerMap entry length != 0")
		}
	}

}

func TestAddLiteralListener(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	e1 := uOfD.NewLiteral(hl)
	lp := uOfD.NewNameLiteralPointer(hl)
	lp.SetLiteral(e1, hl)
	elm := uOfD.literalListenerMap.GetEntry(e1.GetId(hl).String())
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("LiteralListenerMap entry length != 1")
		} else {
			if (*elm)[0] != lp {
				t.Error("LiteralListenerMap entry does not contain LiteralPointer")
			}
		}
	}
	lp.SetLiteral(nil, hl)
	elm = uOfD.literalListenerMap.GetEntry(e1.GetId(hl).String())
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil after SetLiteral(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("LiteralListenerMap entry length != 0")
		}
	}

}

func TestAddLiteralPointerListener(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	lp := uOfD.NewNameLiteralPointer(hl)
	lpp := uOfD.NewLiteralPointerPointer(hl)
	lpp.SetLiteralPointer(lp, hl)
	elm := uOfD.literalPointerListenerMap.GetEntry(lp.GetId(hl).String())
	if elm == nil {
		t.Error("LiteralPointerListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("LiteralPointerListenerMap entry length != 1")
		} else {
			if (*elm)[0] != lpp {
				t.Error("LiteralPointerListenerMap entry does not contain LiteralPointer")
			}
		}
	}
	lpp.SetLiteralPointer(nil, hl)
	elm = uOfD.literalPointerListenerMap.GetEntry(lp.GetId(hl).String())
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil after SetLiteral(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("LiteralListenerMap entry length != 0")
		}
	}
}
