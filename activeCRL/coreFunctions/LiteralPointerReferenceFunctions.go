// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"sync"
)

var LiteralPointerReferenceFunctionsUri string = CoreFunctionsPrefix + "LiteralPointerReferenceFunctions"

var LiteralPointerReferenceCreateUri string = CoreFunctionsPrefix + "LiteralPointerReference/Create"
var LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri = CoreFunctionsPrefix + "LiteralPointerReference/Create/CreatedLiteralPointerReferenceRef"

var LiteralPointerReferenceGetReferencedLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointerReference/GetReferencedLiteralPointer"
var LiteralPointerReferenceGetReferencedLiteralPointerSourceLiteralPointerReferenceRefUri = CoreFunctionsPrefix + "LiteralPointerReference/GetReferencedLiteralPointer/SourceLiteralPointerReferenceRef"
var LiteralPointerReferenceGetReferencedLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerReference/GetReferencedLiteralPointer/IndicatedLiteralPointerRef"

var LiteralPointerReferenceGetLiteralPointerPointerUri string = CoreFunctionsPrefix + "LiteralPointerReference/GetLiteralPointerPointer"
var LiteralPointerReferenceGetLiteralPointerPointerSourceLiteralPointerReferenceRefUri string = CoreFunctionsPrefix + "LiteralPointerReference/GetLiteralPointerPointer/SourceLiteralPointerReferenceRef"
var LiteralPointerReferenceGetLiteralPointerPointerIndicatedLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerReference/GetLiteralPointerPointer/IndicatedLiteralPointerPointerRef"

var LiteralPointerReferenceSetReferencedLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointerReference/SetReferencedLiteralPointer"
var LiteralPointerReferenceSetReferencedLiteralPointerSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerReference/SetReferencedLiteralPointer/SourceLiteralPointerRef"
var LiteralPointerReferenceSetReferencedLiteralPointerModifiedLiteralPointerReferenceRefUri string = CoreFunctionsPrefix + "LiteralPointerReference/SetReferencedLiteralPointer/ModifiedLiteralPointerReferenceRef"

func createLiteralPointerReference(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReferenceReference := core.GetChildElementReferenceWithAncestorUri(element, LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri, hl)
	if createdLiteralPointerReferenceReference == nil {
		createdLiteralPointerReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdLiteralPointerReferenceReference, element, hl)
		core.SetName(createdLiteralPointerReferenceReference, "CreatedLiteralPointerReferenceReference", hl)
		rootCreatedLiteralPointerReferenceReference := uOfD.GetElementReferenceWithUri(LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReferenceReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReferenceReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralPointerReferenceReference, hl)
	}
	createdLiteralPointerReference := createdLiteralPointerReferenceReference.GetReferencedElement(hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		createdLiteralPointerReferenceReference.SetReferencedElement(createdLiteralPointerReference, hl)
	}
}

func getReferencedLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerReferenceGetReferencedLiteralPointerUri)
	if original == nil {
		log.Printf("In GetReferencedLiteralPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralPointerReferenceGetReferencedLiteralPointerSourceLiteralPointerReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetReferencedLiteralPointer, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerReferenceGetReferencedLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if targetLiteralPointerReference == nil {
		log.Printf("In GetReferencedLiteralPointer, the TargetLiteralPointerReference was not found in the replicate")
		return
	}

	targetLiteralPointer := targetLiteralPointerReference.GetReferencedLiteralPointer(hl)
	untypedLiteralPointerReference := sourceReference.GetReferencedElement(hl)
	var sourceLiteralPointerReference core.LiteralPointerReference
	var sourceLiteralPointer core.LiteralPointer
	if untypedLiteralPointerReference != nil {
		switch untypedLiteralPointerReference.(type) {
		case core.LiteralPointerReference:
			sourceLiteralPointerReference = untypedLiteralPointerReference.(core.LiteralPointerReference)
			sourceLiteralPointer = sourceLiteralPointerReference.GetReferencedLiteralPointer(hl)
		}
	}
	if sourceLiteralPointer != targetLiteralPointer {
		targetLiteralPointerReference.SetReferencedLiteralPointer(sourceLiteralPointer, hl)
	}
}

func getLiteralPointerPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerReferenceGetLiteralPointerPointerUri)
	if original == nil {
		log.Printf("In GetLiteralPointerPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralPointerReferenceGetLiteralPointerPointerSourceLiteralPointerReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteralPointerPointer, the SourceLiteralPointerReferenceRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerReferenceGetLiteralPointerPointerIndicatedLiteralPointerPointerRefUri, hl)
	if indicatedLiteralPointerPointerRef == nil {
		log.Printf("In GetLiteralPointerPointer, the IndicatedLiteralPointerPointerRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerPointer := indicatedLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	untypedSourceLiteralPointerReference := sourceReference.GetReferencedElement(hl)
	var sourceLiteralPointerPointer core.LiteralPointerPointer
	if untypedSourceLiteralPointerReference != nil {
		switch untypedSourceLiteralPointerReference.(type) {
		case core.LiteralPointerReference:
			sourceLiteralPointerReference := untypedSourceLiteralPointerReference.(core.LiteralPointerReference)
			sourceLiteralPointerPointer = sourceLiteralPointerReference.GetLiteralPointerPointer(hl)
		default:
			log.Printf("In GetLiteralPointerPointer, the SourceElement is not a LiteralPointerReference")
		}
	}
	if sourceLiteralPointerPointer != indicatedLiteralPointerPointer {
		indicatedLiteralPointerPointerRef.SetReferencedBaseElement(sourceLiteralPointerPointer, hl)
	}
}

func setReferencedLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerReferenceSetReferencedLiteralPointerUri)
	if original == nil {
		log.Printf("In SetReferencedLiteralPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	baseElementReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerReferenceSetReferencedLiteralPointerSourceLiteralPointerRefUri, hl)
	if baseElementReference == nil {
		log.Printf("In SetReferencedLiteralPointer, the LiteralPointerReference was not found in the replicate")
		return
	}

	targetLiteralPointerReference := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralPointerReferenceSetReferencedLiteralPointerModifiedLiteralPointerReferenceRefUri, hl)
	if targetLiteralPointerReference == nil {
		log.Printf("In SetReferencedLiteralPointer, the TargetLiteralPointer was not found in the replicate")
		return
	}

	sourcedLiteralPointer := baseElementReference.GetReferencedLiteralPointer(hl)
	untypedTargetedElement := targetLiteralPointerReference.GetReferencedElement(hl)
	var targetedLiteralPointer core.LiteralPointer
	var targetedLiteralPointerReference core.LiteralPointerReference
	if untypedTargetedElement != nil {
		switch untypedTargetedElement.(type) {
		case core.LiteralPointerReference:
			targetedLiteralPointerReference = untypedTargetedElement.(core.LiteralPointerReference)
			targetedLiteralPointer = targetedLiteralPointerReference.GetReferencedLiteralPointer(hl)
		default:
			log.Printf("In SetReferencedLiteralPointer, the TargetedLiteralPointerReference is not a LiteralPointerReference")
		}
	}
	if sourcedLiteralPointer != targetedLiteralPointer {
		targetedLiteralPointerReference.SetReferencedLiteralPointer(sourcedLiteralPointer, hl)
	}
}

func UpdateRecoveredCoreLiteralPointerReferenceFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralPointerReferenceFunctions
	literalPointerReferenceFunctions := uOfD.GetElementWithUri(LiteralPointerReferenceFunctionsUri)
	if literalPointerReferenceFunctions == nil {
		literalPointerReferenceFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(literalPointerReferenceFunctions, coreFunctionsElement, hl)
		core.SetName(literalPointerReferenceFunctions, "LiteralPointerReferenceFunctions", hl)
		core.SetUri(literalPointerReferenceFunctions, LiteralPointerReferenceFunctionsUri, hl)
	}

	// CreateLiteralPointerReference
	createLiteralPointerReference := uOfD.GetElementWithUri(LiteralPointerReferenceCreateUri)
	if createLiteralPointerReference == nil {
		createLiteralPointerReference = uOfD.NewElement(hl)
		core.SetOwningElement(createLiteralPointerReference, literalPointerReferenceFunctions, hl)
		core.SetName(createLiteralPointerReference, "CreateLiteralPointerReference", hl)
		core.SetUri(createLiteralPointerReference, LiteralPointerReferenceCreateUri, hl)
	}
	// CreatedLiteralPointerReference
	createdLiteralPointerReferenceReference := core.GetChildElementReferenceWithUri(createLiteralPointerReference, LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri, hl)
	if createdLiteralPointerReferenceReference == nil {
		createdLiteralPointerReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdLiteralPointerReferenceReference, createLiteralPointerReference, hl)
		core.SetName(createdLiteralPointerReferenceReference, "CreatedLiteralPointerReferenceRef", hl)
		core.SetUri(createdLiteralPointerReferenceReference, LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri, hl)
	}

	// GetReferencedLiteralPointer
	getReferencedLiteralPointer := uOfD.GetElementWithUri(LiteralPointerReferenceGetReferencedLiteralPointerUri)
	if getReferencedLiteralPointer == nil {
		getReferencedLiteralPointer = uOfD.NewElement(hl)
		core.SetName(getReferencedLiteralPointer, "GetReferencedLiteralPointer", hl)
		core.SetOwningElement(getReferencedLiteralPointer, literalPointerReferenceFunctions, hl)
		core.SetUri(getReferencedLiteralPointer, LiteralPointerReferenceGetReferencedLiteralPointerUri, hl)
	}
	// GetReferencedLiteralPointer.SourceReference
	getLiteralPointerSourceReference := core.GetChildElementReferenceWithUri(getReferencedLiteralPointer, LiteralPointerReferenceGetReferencedLiteralPointerSourceLiteralPointerReferenceRefUri, hl)
	if getLiteralPointerSourceReference == nil {
		getLiteralPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getLiteralPointerSourceReference, getReferencedLiteralPointer, hl)
		core.SetName(getLiteralPointerSourceReference, "SourceLiteralPointerReferenceRef", hl)
		core.SetUri(getLiteralPointerSourceReference, LiteralPointerReferenceGetReferencedLiteralPointerSourceLiteralPointerReferenceRefUri, hl)
	}
	// GetReferencedLiteralPointerTargetLiteralPointerReference
	getLiteralPointerTargetReference := core.GetChildLiteralPointerReferenceWithUri(getReferencedLiteralPointer, LiteralPointerReferenceGetReferencedLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if getLiteralPointerTargetReference == nil {
		getLiteralPointerTargetReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(getLiteralPointerTargetReference, getReferencedLiteralPointer, hl)
		core.SetName(getLiteralPointerTargetReference, "IndicatedLiteralPointerRef", hl)
		core.SetUri(getLiteralPointerTargetReference, LiteralPointerReferenceGetReferencedLiteralPointerIndicatedLiteralPointerRefUri, hl)
	}

	// GetLiteralPointerPointer
	getLiteralPointerPointer := uOfD.GetElementWithUri(LiteralPointerReferenceGetLiteralPointerPointerUri)
	if getLiteralPointerPointer == nil {
		getLiteralPointerPointer = uOfD.NewElement(hl)
		core.SetName(getLiteralPointerPointer, "GetLiteralPointerPointer", hl)
		core.SetOwningElement(getLiteralPointerPointer, literalPointerReferenceFunctions, hl)
		core.SetUri(getLiteralPointerPointer, LiteralPointerReferenceGetLiteralPointerPointerUri, hl)
	}
	// GetLiteralPointerPointer.SourceReference
	getLiteralPointerPointerSourceReference := core.GetChildElementReferenceWithUri(getReferencedLiteralPointer, LiteralPointerReferenceGetLiteralPointerPointerSourceLiteralPointerReferenceRefUri, hl)
	if getLiteralPointerPointerSourceReference == nil {
		getLiteralPointerPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getLiteralPointerPointerSourceReference, getLiteralPointerPointer, hl)
		core.SetName(getLiteralPointerPointerSourceReference, "SourceLiteralPointerReferenceRef", hl)
		core.SetUri(getLiteralPointerPointerSourceReference, LiteralPointerReferenceGetLiteralPointerPointerSourceLiteralPointerReferenceRefUri, hl)
	}
	// GetReferencedLiteralPointerTargetLiteralPointerReference
	getLiteralPointerPointerIndicatedLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithUri(getReferencedLiteralPointer, LiteralPointerReferenceGetLiteralPointerPointerIndicatedLiteralPointerPointerRefUri, hl)
	if getLiteralPointerPointerIndicatedLiteralPointerPointerRef == nil {
		getLiteralPointerPointerIndicatedLiteralPointerPointerRef = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getLiteralPointerPointerIndicatedLiteralPointerPointerRef, getLiteralPointerPointer, hl)
		core.SetName(getLiteralPointerPointerIndicatedLiteralPointerPointerRef, "IndicatedLiteralPointerPointerRef", hl)
		core.SetUri(getLiteralPointerPointerIndicatedLiteralPointerPointerRef, LiteralPointerReferenceGetLiteralPointerPointerIndicatedLiteralPointerPointerRefUri, hl)
	}

	// SetReferencedLiteralPointer
	setReferencedLiteralPointer := uOfD.GetElementWithUri(LiteralPointerReferenceSetReferencedLiteralPointerUri)
	if setReferencedLiteralPointer == nil {
		setReferencedLiteralPointer = uOfD.NewElement(hl)
		core.SetName(setReferencedLiteralPointer, "SetReferencedLiteralPointer", hl)
		core.SetOwningElement(setReferencedLiteralPointer, literalPointerReferenceFunctions, hl)
		core.SetUri(setReferencedLiteralPointer, LiteralPointerReferenceSetReferencedLiteralPointerUri, hl)
	}
	// SetReferencedLiteralPointer.LiteralPointerReference
	setReferencedLiteralPointerLiteralPointerReference := core.GetChildLiteralPointerReferenceWithUri(setReferencedLiteralPointer, LiteralPointerReferenceSetReferencedLiteralPointerSourceLiteralPointerRefUri, hl)
	if setReferencedLiteralPointerLiteralPointerReference == nil {
		setReferencedLiteralPointerLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(setReferencedLiteralPointerLiteralPointerReference, setReferencedLiteralPointer, hl)
		core.SetName(setReferencedLiteralPointerLiteralPointerReference, "SourceLiteralPointerRef", hl)
		core.SetUri(setReferencedLiteralPointerLiteralPointerReference, LiteralPointerReferenceSetReferencedLiteralPointerSourceLiteralPointerRefUri, hl)
	}
	// SetReferencedLiteralPointerTargetLiteralPointerReference
	setReferencedLiteralPointerTargetLiteralPointerReference := core.GetChildElementReferenceWithUri(setReferencedLiteralPointer, LiteralPointerReferenceSetReferencedLiteralPointerModifiedLiteralPointerReferenceRefUri, hl)
	if setReferencedLiteralPointerTargetLiteralPointerReference == nil {
		setReferencedLiteralPointerTargetLiteralPointerReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedLiteralPointerTargetLiteralPointerReference, setReferencedLiteralPointer, hl)
		core.SetName(setReferencedLiteralPointerTargetLiteralPointerReference, "ModifiedLiteralPointerReferenceRef", hl)
		core.SetUri(setReferencedLiteralPointerTargetLiteralPointerReference, LiteralPointerReferenceSetReferencedLiteralPointerModifiedLiteralPointerReferenceRefUri, hl)
	}
}

func literalPointerReferenceFunctionsInit() {
	core.GetCore().AddFunction(LiteralPointerReferenceCreateUri, createLiteralPointerReference)
	core.GetCore().AddFunction(LiteralPointerReferenceGetReferencedLiteralPointerUri, getReferencedLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerReferenceGetLiteralPointerPointerUri, getLiteralPointerPointer)
	core.GetCore().AddFunction(LiteralPointerReferenceSetReferencedLiteralPointerUri, setReferencedLiteralPointer)
}
