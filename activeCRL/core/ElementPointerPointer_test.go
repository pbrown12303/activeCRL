package core

import (
	"encoding/json"
	"testing"
)

func TestNewElementPointerPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := uOfD.NewElementPointerReference()
	epp := uOfD.NewElementPointerPointer()
	epp.SetOwningElement(owner)
	if epp.GetOwningElement() != owner {
		t.Error("Element pointer pointer's owner not set properly")
	}
	var found bool = false
	for key, _ := range owner.getOwnedBaseElements() {
		if key == epp.GetId().String() {
			found = true
		}
	}
	if found == false {
		t.Error("Element Pointer Pointer not found in parent's OwnedBaseElements \n")
	}
	if owner.getElementPointerPointer() != epp {
		t.Error("Owner.getElementPointerPointer() did not return Referenced Element Pointer")
	}
}

func TestSetElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := uOfD.NewElementPointerReference()
	epp := uOfD.NewElementPointerPointer()
	epp.SetOwningElement(owner)
	target := uOfD.NewReferencedElementPointer()
	target.SetOwningElement(owner)
	epp.SetElementPointer(target)
	if epp.GetElementPointer() != target {
		t.Error("ElementPointer not set to target properly")
	}
	epp.SetElementPointer(nil)
	if epp.GetElementPointer() != nil {
		t.Error("ElementPointer not set to nil properly")
	}
}

func TestElementPointerPointerMarshal(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := uOfD.NewElementPointerReference()
	epp := uOfD.NewElementPointerPointer()
	epp.SetOwningElement(owner)
	target := uOfD.NewReferencedElementPointer()
	target.SetOwningElement(owner)
	epp.SetElementPointer(target)

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

func TestElementPointerPointerClone(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	owner := uOfD.NewElementPointerReference()
	epp := uOfD.NewElementPointerPointer()
	epp.SetOwningElement(owner)
	target := uOfD.NewReferencedElementPointer()
	target.SetOwningElement(owner)
	epp.SetElementPointer(target)

	clone := epp.(*elementPointerPointer).clone()
	if !Equivalent(epp, clone) {
		t.Error("ElementPointerPointer clone failed")
		Print(epp, "   ")
		Print(clone, "   ")
	}

}
