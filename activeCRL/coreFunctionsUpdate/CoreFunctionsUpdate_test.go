package main

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/coreFunctions"

	"testing"
)

func TestUpdateCoreFunctions(t *testing.T) {
	//	log.Printf("Entering TestUpdateCoreElement")
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(false)
	var emptyCore core.Element

	//Core
	recoveredCoreFunctions := updateRecoveredCoreFunctions(emptyCore, uOfD, hl)
	if recoveredCoreFunctions == nil {
		t.Error("updateRecoveredCore returned empty element")
	}
	if core.GetUri(recoveredCoreFunctions, hl) != coreFunctions.CoreFunctionsUri {
		t.Error("CoreFunctions uri not set")
	}
	_, ok := recoveredCoreFunctions.(core.Element)
	if !ok {
		t.Error("Core is of wrong type")
	}

	// CreateElement
	recoveredBaseElement := uOfD.GetBaseElementWithUri(coreFunctions.CreateElememtUri)
	if recoveredBaseElement == nil {
		t.Error("CreateElement not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("CreateElement is of wrong type")
	}

	// CreatedElementReference
	recoveredCreatedElementReference := core.GetChildElementReferenceWithUri(recoveredBaseElement.(core.Element), coreFunctions.CreatedElementReferenceUri, hl)
	if recoveredCreatedElementReference == nil {
		t.Error("CreaedElementReference not found")
	}
}
