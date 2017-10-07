package core

import (
	"sync"
	"testing"

	"github.com/satori/go.uuid"
)

func TestBaseElement(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	var be baseElement
	// Test id
	if be.GetId(hl) != uuid.Nil {
		t.Error("baseElement identifier not nil before setting")
	}

	be.initializeBaseElement()
	if be.GetId(hl) == uuid.Nil {
		t.Error("baseElement identifier nil after setting")
	}

	// Test version
	if be.GetVersion(hl) != 0 {
		t.Error("baseElement version not 0 before increment")
	}
	be.internalIncrementVersion()
	if be.GetVersion(hl) != 1 {
		t.Error("baseElement version not 1 after increment")
	}

}

func TestGetUri(t *testing.T) {
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := NewUniverseOfDiscourse()
	e1 := uOfD.NewElement(hl)
	testUri := "testUri"
	SetUri(e1, testUri, hl)
	if GetUri(e1, hl) != testUri {
		t.Error("Uri not returned")
	}
}
