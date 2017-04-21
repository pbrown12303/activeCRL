package core

import (
	"testing"
)

func TestEquivalence(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	parent := NewElement(uOfD)
	child := NewElement(uOfD)
	child.SetOwningElement(parent)
	if Equivalent(parent, parent) != true {
		t.Errorf("Equivalence test failed")
	}
}
