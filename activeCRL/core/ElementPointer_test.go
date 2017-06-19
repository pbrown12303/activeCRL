package core

import (
	"encoding/json"
	"testing"
)

func TestNewOwningElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := NewElement(uOfD)
	oep := NewOwningElementPointer(uOfD)
	oep.SetOwningElement(owner)
	if oep.GetOwningElement() != owner {
		t.Error("Owning element pointer's owner not set properly")
	}
	var found bool = false
	for key, _ := range owner.getOwnedBaseElements() {
		if key == oep.GetId().String() {
			found = true
		}
	}
	if found == false {
		t.Error("Owning Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if oep.getElementPointerRole() != OWNING_ELEMENT {
		t.Error("Owning Element Pointer role not OWNING_ELEMENT \n")
	}
	if owner.getOwningElementPointer() != oep {
		t.Error("Owner.getOwningElementPointer() did not return Owning Element Pointer")
	}
}

func TestReferencedElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := NewElementReference(uOfD)
	rep := NewReferencedElementPointer(uOfD)
	rep.SetOwningElement(owner)
	if rep.GetOwningElement() != owner {
		t.Error("Referenced element pointer's owner not set properly")
	}
	var found bool = false
	for key, _ := range owner.getOwnedBaseElements() {
		if key == rep.GetId().String() {
			found = true
		}
	}
	if found == false {
		t.Error("Referenced Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if rep.getElementPointerRole() != REFERENCED_ELEMENT {
		t.Error("Referenced Element Pointer role not REFERENCED_ELEMENT \n")
	}
	if owner.getReferencedElementPointer() != rep {
		t.Error("Owner.getReferencedElementPointer() did not return Referenced Element Pointer")
	}
}

func TestAbstractElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := NewRefinement(uOfD)
	aep := NewAbstractElementPointer(uOfD)
	aep.SetOwningElement(owner)
	if aep.GetOwningElement() != owner {
		t.Error("Abstract element pointer's owner not set properly")
	}
	var found bool = false
	for key, _ := range owner.getOwnedBaseElements() {
		if key == aep.GetId().String() {
			found = true
		}
	}
	if found == false {
		t.Error("Abstract Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if aep.getElementPointerRole() != ABSTRACT_ELEMENT {
		t.Error("Abstract Element Pointer role not ABSTRACT_ELEMENT \n")
	}
	if owner.getAbstractElementPointer() != aep {
		t.Error("Owner.getAbstractElementPointer() did not return Abstract Element Pointer")
	}
}

func TestRefinedElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := NewRefinement(uOfD)
	rep := NewRefinedElementPointer(uOfD)
	rep.SetOwningElement(owner)
	if rep.GetOwningElement() != owner {
		t.Error("Refined element pointer's owner not set properly")
	}
	var found bool = false
	for key, _ := range owner.getOwnedBaseElements() {
		if key == rep.GetId().String() {
			found = true
		}
	}
	if found == false {
		t.Error("Refined Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if rep.getElementPointerRole() != REFINED_ELEMENT {
		t.Error("Refined Element Pointer role not REFINED_ELEMENT \n")
	}
	if owner.getRefinedElementPointer() != rep {
		t.Error("Owner.getRefinedElementPointer() did not return Abstract Element Pointer")
	}
}

func TestSetElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := NewElementReference(uOfD)
	rep := NewReferencedElementPointer(uOfD)
	rep.SetOwningElement(owner)
	target := NewElement(uOfD)
	target.SetOwningElement(owner)
	rep.SetElement(target)
	if rep.GetElement() != target {
		t.Error("ReferencedElementPointer SetElement(target) did not work")
	}
	rep.SetElement(nil)
	if rep.GetElement() != nil {
		t.Error("ReferencedElementPointer SetElement(nil) did not work")
	}
}

func TestElementPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := NewElementReference(uOfD)
	rep := NewReferencedElementPointer(uOfD)
	rep.SetOwningElement(owner)
	target := NewElement(uOfD)
	target.SetOwningElement(owner)
	rep.SetElement(target)

	result, err := json.MarshalIndent(owner, "", "   ")
	if err != nil {
		t.Error(err)
	}

	uOfD2 := NewUniverseOfDiscourse()
	recoveredOwner := RecoverElement(result, uOfD2)
	if !Equivalent(owner, recoveredOwner) {
		t.Error("Recovered owner not equivalent to original owner")
	}
}

func TestElementPointerClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := NewElementReference(uOfD)
	rep := NewReferencedElementPointer(uOfD)
	rep.SetOwningElement(owner)
	target := NewElement(uOfD)
	target.SetOwningElement(owner)
	rep.SetElement(target)

	clone := rep.(*elementPointer).clone()
	if !Equivalent(rep, clone) {
		t.Error("ElementPointer clone failed")
		Print(rep, "   ")
		Print(clone, "   ")
	}

}
