package core

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestBaseElement(t *testing.T) {
	var be baseElement
	// Test id
	if be.GetId() != uuid.Nil {
		t.Error("baseElement identifier not nil before setting")
	}

	be.initializeBaseElement()
	if be.GetId() == uuid.Nil {
		t.Error("baseElement identifier nil after setting")
	}

	// Test version
	if be.GetVersion() != 0 {
		t.Error("baseElement version not 0 before increment")
	}
	be.internalIncrementVersion()
	if be.GetVersion() != 1 {
		t.Error("baseElement version not 1 after increment")
	}

}
