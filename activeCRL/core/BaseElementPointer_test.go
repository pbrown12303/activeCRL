package core

import (
	"encoding/json"
	"testing"
)

func TestBaseElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := uOfD.NewBaseElementReference()
	rep := uOfD.NewBaseElementPointer()
	rep.SetOwningElement(owner)
	if rep.GetOwningElement() != owner {
		t.Error("Base element pointer's owner not set properly")
	}
	var found bool = false
	for key, _ := range owner.getOwnedBaseElements() {
		if key == rep.GetId().String() {
			found = true
		}
	}
	if found == false {
		t.Error("Base Element Pointer not found in parent's OwnedBaseElements \n")
	}
	if owner.getBaseElementPointer() != rep {
		t.Error("Owner.getBaseElementPointer() did not return Base Element Pointer")
	}
}

func TestSetBaseElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := uOfD.NewBaseElementReference()
	rep := uOfD.NewBaseElementPointer()
	rep.SetOwningElement(owner)
	target := uOfD.NewElement()
	target.SetOwningElement(owner)
	rep.SetBaseElement(target)
	if rep.GetBaseElement() != target {
		t.Error("BaseElementPointer SetBaseElement(target) did not work")
	}
	rep.SetBaseElement(nil)
	if rep.GetBaseElement() != nil {
		t.Error("BaseElementPointer SetBaseElement(nil) did not work")
	}
}

func TestBaseElementPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := uOfD.NewBaseElementReference()
	rep := uOfD.NewBaseElementPointer()
	rep.SetOwningElement(owner)
	target := uOfD.NewElement()
	target.SetOwningElement(owner)
	rep.SetBaseElement(target)

	result, err := json.MarshalIndent(owner, "", "   ")
	if err != nil {
		t.Error(err)
	}

	uOfD2 := NewUniverseOfDiscourse()
	recoveredOwner := uOfD2.RecoverElement(result)
	if !Equivalent(owner, recoveredOwner) {
		t.Error("Recovered owner not equivalent to original owner")
	}
}

func TestBaseElementPointerClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := uOfD.NewBaseElementReference()
	rep := uOfD.NewBaseElementPointer()
	rep.SetOwningElement(owner)
	target := uOfD.NewElement()
	target.SetOwningElement(owner)
	rep.SetBaseElement(target)

	clone := rep.(*baseElementPointer).clone()
	if !Equivalent(rep, clone) {
		t.Error("ElementPointer clone failed")
		Print(rep, "   ")
		Print(clone, "   ")
	}

}
