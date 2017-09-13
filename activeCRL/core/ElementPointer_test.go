package core

import (
	"encoding/json"
	"testing"
)

func TestNewOwningElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	owner := uOfD.NewElement(hl)
	oep := uOfD.NewOwningElementPointer(hl)
	SetOwningElement(oep, owner, hl)
	if GetOwningElement(oep, hl) != owner {
		t.Error("Owning element pointer's owner not set properly")
	}
	var found bool = false
	for _, be := range owner.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == oep.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("Owning Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if oep.GetElementPointerRole(hl) != OWNING_ELEMENT {
		t.Error("Owning Element Pointer role not OWNING_ELEMENT \n")
	}
	if owner.GetOwningElementPointer(hl) != oep {
		t.Error("Owner.getOwningElementPointer() did not return Owning Element Pointer")
	}
}

func TestReferencedElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	owner := uOfD.NewElementReference(hl)
	rep := uOfD.NewReferencedElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	if GetOwningElement(rep, hl) != owner {
		t.Error("Referenced element pointer's owner not set properly")
	}
	var found bool = false
	for _, be := range owner.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == rep.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("Referenced Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if rep.GetElementPointerRole(hl) != REFERENCED_ELEMENT {
		t.Error("Referenced Element Pointer role not REFERENCED_ELEMENT \n")
	}
	if owner.GetElementPointer(hl) != rep {
		t.Error("Owner.getElementPointer() did not return Element Pointer")
	}
}

func TestAbstractElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	owner := uOfD.NewRefinement(hl)
	aep := uOfD.NewAbstractElementPointer(hl)
	SetOwningElement(aep, owner, hl)
	if GetOwningElement(aep, hl) != owner {
		t.Error("Abstract element pointer's owner not set properly")
	}
	var found bool = false
	for _, be := range owner.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == aep.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("Abstract Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if aep.GetElementPointerRole(hl) != ABSTRACT_ELEMENT {
		t.Error("Abstract Element Pointer role not ABSTRACT_ELEMENT \n")
	}
	if owner.GetAbstractElementPointer(hl) != aep {
		t.Error("Owner.getAbstractElementPointer() did not return Abstract Element Pointer")
	}
}

func TestRefinedElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	owner := uOfD.NewRefinement(hl)
	rep := uOfD.NewRefinedElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	if GetOwningElement(rep, hl) != owner {
		t.Error("Refined element pointer's owner not set properly")
	}
	var found bool = false
	for _, be := range owner.GetOwnedBaseElements(hl) {
		if be.GetId(hl) == rep.GetId(hl) {
			found = true
		}
	}
	if found == false {
		t.Error("Refined Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if rep.GetElementPointerRole(hl) != REFINED_ELEMENT {
		t.Error("Refined Element Pointer role not REFINED_ELEMENT \n")
	}
	if owner.GetRefinedElementPointer(hl) != rep {
		t.Error("Owner.getRefinedElementPointer() did not return Abstract Element Pointer")
	}
}

func TestSetElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	owner := uOfD.NewElementReference(hl)
	rep := uOfD.NewReferencedElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	target := uOfD.NewElement(hl)
	SetOwningElement(target, owner, hl)
	rep.SetElement(target, hl)
	if rep.GetElement(hl) != target {
		t.Error("ReferencedElementPointer SetElement(target) did not work")
	}
	if rep.GetElementId(hl) != target.GetId(hl) {
		t.Error("ReferencedElementPointer GetElementId() returns wrong value after SetElement()")
	}
	if rep.GetElementVersion(hl) != target.GetVersion(hl) {
		t.Error("ReferencedElementPointer GetElementVersion() returns wrong value after SetElement()")
	}
	rep.SetElement(nil, hl)
	if rep.GetElement(hl) != nil {
		t.Error("ReferencedElementPointer SetElement(nil) did not work")
	}
}

func TestElementPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	owner := uOfD.NewElementReference(hl)
	rep := uOfD.NewReferencedElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	target := uOfD.NewElement(hl)
	SetOwningElement(target, owner, hl)
	rep.SetElement(target, hl)

	result, err := json.MarshalIndent(owner, "", "   ")
	if err != nil {
		t.Error(err)
	}

	uOfD2 := NewUniverseOfDiscourse()
	recoveredOwner := uOfD2.RecoverElement(result)
	if !Equivalent(owner, recoveredOwner, hl) {
		t.Error("Recovered owner not equivalent to original owner")
	}
}

func TestElementPointerClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	hl := NewHeldLocks()
	defer hl.ReleaseLocks()
	owner := uOfD.NewElementReference(hl)
	rep := uOfD.NewReferencedElementPointer(hl)
	SetOwningElement(rep, owner, hl)
	target := uOfD.NewElement(hl)
	SetOwningElement(target, owner, hl)
	rep.SetElement(target, hl)

	clone := rep.(*elementPointer).clone()
	if !Equivalent(rep, clone, hl) {
		t.Error("ElementPointer clone failed")
		Print(rep, "   ", hl)
		Print(clone, "   ", hl)
	}

}
