package core

import (
	"testing"
)

func TestEquivalence(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := uOfD.NewElement()
	child := uOfD.NewElement()
	child.SetOwningElement(parent)
	if Equivalent(parent, parent) != true {
		t.Errorf("Equivalence test failed")
	}
}
