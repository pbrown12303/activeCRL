package main

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/coreFunctions"

	"testing"
)

func TestUpdateCoreElement(t *testing.T) {
	//	log.Printf("Entering TestUpdateCoreElement")
	uOfD := core.NewUniverseOfDiscourse()
	uOfD.SetRecordingUndo(false)
	var emptyCore core.Element
	core.Print(emptyCore, "")

	//Core
	recoveredCore := updateRecoveredCoreFunctions(emptyCore, uOfD)
	if recoveredCore == nil {
		t.Error("updateRecoveredCore returned empty element")
	}
	if recoveredCore.GetUri() != core.CoreConceptSpaceUri {
		t.Error("Core uri not set")
	}
	_, ok := recoveredCore.(core.Element)
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
	recoveredCreatedElementReference := core.GetChildElementReferenceWithUri(recoveredBaseElement.(core.Element), coreFunctions.CreatedElementReferenceUri)
	if recoveredCreatedElementReference == nil {
		t.Error("CreaedElementReference not found")
	}
}
