// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
	"strconv"
	"sync"
	"testing"
	//	"time"
)

func TestBaseElementFunctionsIds(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// var BaseElementFunctionsUri string = CoreFunctionsPrefix + "BaseElement"
	validateElementId(t, uOfD, hl, BaseElementFunctionsUri)

	//var BaseElementDeleteUri string = CoreFunctionsPrefix + "BaseElement/Delete"
	validateElementId(t, uOfD, hl, BaseElementDeleteUri)
	//var BaseElementDeleteDeletedElementRefUri string = CoreFunctionsPrefix + "BaseElement/Delete/DeletedElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementDeleteDeletedElementRefUri)
	//
	//var BaseElementGetIdUri string = CoreFunctionsPrefix + "BaseElement/GetId"
	validateElementId(t, uOfD, hl, BaseElementGetIdUri)
	//var BaseElementGetIdSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetId/SourceBaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementGetIdSourceBaseElementRefUri)
	//var BaseElementGetIdCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElement/GetId/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, BaseElementGetIdCreatedLiteralRefUri)
	//
	//var BaseElementGetNameUri string = CoreFunctionsPrefix + "BaseElement/GetName"
	validateElementId(t, uOfD, hl, BaseElementGetNameUri)
	//var BaseElementGetNameSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetName/SourceBaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementGetNameSourceBaseElementRefUri)
	//var BaseElementGetNameCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElement/GetName/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, BaseElementGetNameCreatedLiteralRefUri)
	//
	//var BaseElementGetOwningElementUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement"
	validateElementId(t, uOfD, hl, BaseElementGetOwningElementUri)
	//var BaseElementGetOwningElementSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement/SourceBaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementGetOwningElementSourceBaseElementRefUri)
	//var BaseElementGetOwningElementOwningElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement/TargetElementReference"
	validateElementReferenceId(t, uOfD, hl, BaseElementGetOwningElementOwningElementRefUri)
	//
	//var BaseElementGetUriUri string = CoreFunctionsPrefix + "BaseElement/GetUri"
	validateElementId(t, uOfD, hl, BaseElementGetUriUri)
	//var BaseElementGetUriSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetUri/SourceBaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementGetUriSourceBaseElementRefUri)
	//var BaseElementGetUriCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElement/GetUri/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, BaseElementGetUriCreatedLiteralRefUri)
	//
	//var BaseElementGetVersionUri string = CoreFunctionsPrefix + "BaseElement/GetVersion"
	validateElementId(t, uOfD, hl, BaseElementGetVersionUri)
	//var BaseElementGetVersionSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetVersion/SourceBaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementGetVersionSourceBaseElementRefUri)
	//var BaseElementGetVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElement/GetVersion/CreatedLiteralRef"
	validateLiteralReferenceId(t, uOfD, hl, BaseElementGetVersionCreatedLiteralRefUri)
	//
	//var BaseElementSetOwningElementUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement"
	validateElementId(t, uOfD, hl, BaseElementSetOwningElementUri)
	//var BaseElementSetOwningElementOwningElementRefUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement/OwningElementRef"
	validateElementReferenceId(t, uOfD, hl, BaseElementSetOwningElementOwningElementRefUri)
	//var BaseElementSetOwningElementModifiedBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement/ModifiedBaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementSetOwningElementModifiedBaseElementRefUri)
	//
	//var BaseElementSetUriUri string = CoreFunctionsPrefix + "BaseElement/SetUri"
	validateElementId(t, uOfD, hl, BaseElementSetUriUri)
	//var BaseElementSetUriSourceUriRefUri string = CoreFunctionsPrefix + "BaseElement/SetUri/SourceUriRef"
	validateLiteralReferenceId(t, uOfD, hl, BaseElementSetUriSourceUriRefUri)
	//var BaseElementSetUriModifiedBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/SetUri/ModifiedBaseElementRef"
	validateBaseElementReferenceId(t, uOfD, hl, BaseElementSetUriModifiedBaseElementRefUri)

}

func TestDelete(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	deleteFunction := core.GetCore().GetFunction(BaseElementDeleteUri)
	if deleteFunction == nil {
		t.Errorf("Delete function not registered with core")
	}

	// Get Ancestor
	del := uOfD.GetElementWithUri(BaseElementDeleteUri)
	if del == nil {
		t.Errorf("Delete function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(del, hl)
	//	core.Print(replicate, "In TestDelete, replicate: ", hl)

	replicateFunctions := core.GetCore().FindFunctions(replicate, nil, hl)
	if len(replicateFunctions) != 1 {
		t.Errorf("Function not found associated with replicate")
	}

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	//	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replications
	if uOfD.IsRefinementOf(replicate, del, hl) != true {
		t.Errorf("Replicate is not refinement of Delete()")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementDeleteDeletedElementRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	targetElement := uOfD.NewElement(hl)
	targetElementId := targetElement.GetId(hl)
	if uOfD.GetBaseElement(targetElementId) != targetElement {
		t.Error("TargetElement not created successfully")
	}

	uOfD.MarkUndoPoint()

	// Set the targetReference, release the locks, and check for successful deletion
	// The setting of the target reference should trigger the deletion of the target reference
	targetReference.SetReferencedBaseElement(targetElement, hl)
	hl.ReleaseLocks()
	wg.Wait()
	//	time.Sleep(10000000 * time.Nanosecond)

	if uOfD.GetBaseElement(targetElementId) != nil {
		t.Error("TargetElement not deleted successfully")
	}
	if targetReference.GetReferencedBaseElement(hl) != nil {
		t.Error("TargetReference.ReferencedBaseElement not nil after deletion")
	}

	// Test Undo
	uOfD.Undo(hl)

	if uOfD.GetBaseElement(targetElementId) != targetElement {
		t.Error("TargetElement deletion not undone successfully")
	}
	if targetReference.GetReferencedBaseElement(hl) != nil {
		t.Error("TargetReference.ReferencedBaseElement not nil after undo")
	}

	uOfD.Redo(hl)
	if uOfD.GetBaseElement(targetElementId) != nil {
		t.Error("TargetElement not redone successfully")
	}
	if targetReference.GetReferencedBaseElement(hl) != nil {
		t.Error("TargetReference.ReferencedBaseElement not nil after redo")
	}

}

func TestGetId(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getId := uOfD.GetElementWithUri(BaseElementGetIdUri)
	if getId == nil {
		t.Errorf("GetId function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getId, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	//	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getId, hl) != true {
		t.Errorf("Replicate is not refinement of GetId()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetIdSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetIdCreatedLiteralRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test literal update functionality
	source := uOfD.NewElement(hl)
	sourceName := "SourceName"
	core.SetName(source, sourceName, hl)
	sourceReference.SetReferencedBaseElement(source, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	//	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetLiteral := targetReference.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != sourceName {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestGetName(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getName := uOfD.GetElementWithUri(BaseElementGetNameUri)
	if getName == nil {
		t.Errorf("GetName function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getName, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	//	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getName, hl) != true {
		t.Errorf("Replicate is not refinement of GetName()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetNameSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetNameCreatedLiteralRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test literal update functionality
	source := uOfD.NewElement(hl)
	sourceName := "SourceName"
	core.SetName(source, sourceName, hl)
	sourceReference.SetReferencedBaseElement(source, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	// time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetLiteral := targetReference.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != sourceName {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestGetOwningElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getOwningElement := uOfD.GetElementWithUri(BaseElementGetOwningElementUri)
	if getOwningElement == nil {
		t.Errorf("GetOwningElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getOwningElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	//	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getOwningElement, hl) != true {
		t.Errorf("Replicate is not refinement of GetOwningElement()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetOwningElementSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementGetOwningElementOwningElementRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	parent := uOfD.NewElement(hl)
	source := uOfD.NewElement(hl)
	core.SetOwningElement(source, parent, hl)
	sourceReference.SetReferencedBaseElement(source, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	//	time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetElement := targetReference.GetReferencedElement(hl)
	if targetElement == nil {
		t.Errorf("Target element not found")
	} else {
		if targetElement != parent {
			t.Errorf("Target element value incorrect")
		}
	}
}

func TestGetUri(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getUri := uOfD.GetElementWithUri(BaseElementGetUriUri)
	if getUri == nil {
		t.Errorf("GetUri function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getUri, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	//	time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getUri, hl) != true {
		t.Errorf("Replicate is not refinement of GetUri()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetUriSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetUriCreatedLiteralRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test literal update functionality
	source := uOfD.NewElement(hl)
	sourceName := "SourceName"
	core.SetName(source, sourceName, hl)
	sourceReference.SetReferencedBaseElement(source, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	// time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetLiteral := targetReference.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != sourceName {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestGetVersion(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	getVersion := uOfD.GetElementWithUri(BaseElementGetVersionUri)
	if getVersion == nil {
		t.Errorf("GetVersion function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(getVersion, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	// time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, getVersion, hl) != true {
		t.Errorf("Replicate is not refinement of GetVersion()")
	}
	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetVersionSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		t.Errorf("SourceReference child not found")
	}
	targetReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetVersionCreatedLiteralRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test literal update functionality
	source := uOfD.NewElement(hl)
	sourceName := "SourceName"
	core.SetName(source, sourceName, hl)
	sourceReference.SetReferencedBaseElement(source, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	// time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	targetLiteral := targetReference.GetReferencedLiteral(hl)
	if targetLiteral == nil {
		t.Errorf("Target literal not found")
	} else {
		if targetLiteral.GetLiteralValue(hl) != strconv.Itoa(source.GetVersion(hl)) {
			t.Errorf("Target literal value incorrect")
		}
	}
}

func TestSetOwningElement(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setOwningElement := uOfD.GetElementWithUri(BaseElementSetOwningElementUri)
	if setOwningElement == nil {
		t.Errorf("SetOwningElement function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setOwningElement, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	// time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, setOwningElement, hl) != true {
		t.Errorf("Replicate is not refinement of SetOwningElement()")
	}
	owningElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementSetOwningElementOwningElementRefUri, hl)
	if owningElementReference == nil {
		t.Errorf("OwningElementReference child not found")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementSetOwningElementModifiedBaseElementRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	parent := uOfD.NewElement(hl)
	target := uOfD.NewElement(hl)
	targetReference.SetReferencedBaseElement(target, hl)
	owningElementReference.SetReferencedElement(parent, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	// time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	if core.GetOwningElement(target, hl) != parent {
		t.Errorf("Target owner not set properly")
	}
}

func TestSetUri(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	AddCoreFunctionsToUofD(uOfD, hl)

	// Get Ancestor
	setUri := uOfD.GetElementWithUri(BaseElementSetUriUri)
	if setUri == nil {
		t.Errorf("SetUri function representation not found")
	}

	// Create the instance
	replicate := core.CreateReplicateAsRefinement(setUri, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	// time.Sleep(10000000 * time.Nanosecond)

	// Now check the replication
	if uOfD.IsRefinementOf(replicate, setUri, hl) != true {
		t.Errorf("Replicate is not refinement of SetUri()")
	}
	uriReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementSetUriSourceUriRefUri, hl)
	if uriReference == nil {
		t.Errorf("UriReference child not found")
	}
	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementSetUriModifiedBaseElementRefUri, hl)
	if targetReference == nil {
		t.Errorf("TargetReference child not found")
	}

	// Now test target reference update functionality
	uriLiteral := uOfD.NewLiteral(hl)
	uri := "TestUri"
	uriLiteral.SetLiteralValue(uri, hl)
	uriReference.SetReferencedLiteral(uriLiteral, hl)
	target := uOfD.NewElement(hl)
	targetReference.SetReferencedBaseElement(target, hl)

	// Locks must be released to allow function to execute
	hl.ReleaseLocks()
	wg.Wait()
	// time.Sleep(10000000 * time.Nanosecond)

	hl.LockBaseElement(replicate)
	if core.GetUri(target, hl) != uri {
		t.Errorf("Uri not set properly")
	}
}
