// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"sync"
	"testing"
)

func TestBuildCoreFunctions(t *testing.T) {
	//	log.Printf("Entering TestUpdateCoreElement")
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(false)

	//Core
	builtCoreFunctions := BuildCoreFunctionsConceptSpace(uOfD, hl)
	if builtCoreFunctions == nil {
		t.Error("builtCoreFunctions returned empty element")
	}
	if core.GetUri(builtCoreFunctions, hl) != CoreFunctionsUri {
		t.Error("CoreFunctions uri not set")
	}
	_, ok := builtCoreFunctions.(core.Element)
	if !ok {
		t.Error("Core is of wrong type")
	}

	// CreateElement
	recoveredBaseElement := uOfD.GetBaseElementWithUri(ElementCreateUri)
	if recoveredBaseElement == nil {
		t.Error("CreateElement not found")
	}
	_, ok = recoveredBaseElement.(core.Element)
	if !ok {
		t.Error("CreateElement is of wrong type")
	}

	// CreatedElementReference
	recoveredCreatedElementReference := core.GetChildElementReferenceWithUri(recoveredBaseElement.(core.Element), ElementCreateCreatedElementRefUri, hl)
	if recoveredCreatedElementReference == nil {
		t.Error("CreaedElementReference not found")
	}
}
