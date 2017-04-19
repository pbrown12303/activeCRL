package core

import (
	"testing"
)

func TestEquivaalence(t *testing.T) {
	var uOfD UniverseOfDiscourse
	parent := NewElement(&uOfD)
	child := NewElement(&uOfD)
	child.setOwningElement(parent)
	if Equivalent(parent, parent) != true {
		t.Errorf("Equivalence test failed")
	}
}
