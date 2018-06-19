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
		core.SetLabel(createdLiteralPointerReferenceReference, "CreatedLiteralPointerReferenceReference", hl)
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

func BuildCoreLiteralPointerReferenceFunctions(coreFunctionsElement core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralPointerReferenceFunctions
	literalPointerReferenceFunctions := uOfD.NewElement(hl, LiteralPointerReferenceFunctionsUri)
	core.SetOwningElement(literalPointerReferenceFunctions, coreFunctionsElement, hl)
	core.SetLabel(literalPointerReferenceFunctions, "LiteralPointerReferenceFunctions", hl)
	core.SetUri(literalPointerReferenceFunctions, LiteralPointerReferenceFunctionsUri, hl)

	// CreateLiteralPointerReference
	createLiteralPointerReference := uOfD.NewElement(hl, LiteralPointerReferenceCreateUri)
	core.SetOwningElement(createLiteralPointerReference, literalPointerReferenceFunctions, hl)
	core.SetLabel(createLiteralPointerReference, "CreateLiteralPointerReference", hl)
	core.SetUri(createLiteralPointerReference, LiteralPointerReferenceCreateUri, hl)
	// CreatedLiteralPointerReference
	createdLiteralPointerReferenceReference := uOfD.NewElementReference(hl, LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri)
	core.SetOwningElement(createdLiteralPointerReferenceReference, createLiteralPointerReference, hl)
	core.SetLabel(createdLiteralPointerReferenceReference, "CreatedLiteralPointerReferenceRef", hl)
	core.SetUri(createdLiteralPointerReferenceReference, LiteralPointerReferenceCreateCreatedLiteralPointerReferenceRefUri, hl)

	// GetReferencedLiteralPointer
	getReferencedLiteralPointer := uOfD.NewElement(hl, LiteralPointerReferenceGetReferencedLiteralPointerUri)
	core.SetLabel(getReferencedLiteralPointer, "GetReferencedLiteralPointer", hl)
	core.SetOwningElement(getReferencedLiteralPointer, literalPointerReferenceFunctions, hl)
	core.SetUri(getReferencedLiteralPointer, LiteralPointerReferenceGetReferencedLiteralPointerUri, hl)
	// GetReferencedLiteralPointer.SourceReference
	getLiteralPointerSourceReference := uOfD.NewElementReference(hl, LiteralPointerReferenceGetReferencedLiteralPointerSourceLiteralPointerReferenceRefUri)
	core.SetOwningElement(getLiteralPointerSourceReference, getReferencedLiteralPointer, hl)
	core.SetLabel(getLiteralPointerSourceReference, "SourceLiteralPointerReferenceRef", hl)
	core.SetUri(getLiteralPointerSourceReference, LiteralPointerReferenceGetReferencedLiteralPointerSourceLiteralPointerReferenceRefUri, hl)
	// GetReferencedLiteralPointerTargetLiteralPointerReference
	getLiteralPointerTargetReference := uOfD.NewLiteralPointerReference(hl, LiteralPointerReferenceGetReferencedLiteralPointerIndicatedLiteralPointerRefUri)
	core.SetOwningElement(getLiteralPointerTargetReference, getReferencedLiteralPointer, hl)
	core.SetLabel(getLiteralPointerTargetReference, "IndicatedLiteralPointerRef", hl)
	core.SetUri(getLiteralPointerTargetReference, LiteralPointerReferenceGetReferencedLiteralPointerIndicatedLiteralPointerRefUri, hl)

	// GetLiteralPointerPointer
	getLiteralPointerPointer := uOfD.NewElement(hl, LiteralPointerReferenceGetLiteralPointerPointerUri)
	core.SetLabel(getLiteralPointerPointer, "GetLiteralPointerPointer", hl)
	core.SetOwningElement(getLiteralPointerPointer, literalPointerReferenceFunctions, hl)
	core.SetUri(getLiteralPointerPointer, LiteralPointerReferenceGetLiteralPointerPointerUri, hl)
	// GetLiteralPointerPointer.SourceReference
	getLiteralPointerPointerSourceReference := uOfD.NewElementReference(hl, LiteralPointerReferenceGetLiteralPointerPointerSourceLiteralPointerReferenceRefUri)
	core.SetOwningElement(getLiteralPointerPointerSourceReference, getLiteralPointerPointer, hl)
	core.SetLabel(getLiteralPointerPointerSourceReference, "SourceLiteralPointerReferenceRef", hl)
	core.SetUri(getLiteralPointerPointerSourceReference, LiteralPointerReferenceGetLiteralPointerPointerSourceLiteralPointerReferenceRefUri, hl)
	// GetReferencedLiteralPointerTargetLiteralPointerReference
	getLiteralPointerPointerIndicatedLiteralPointerPointerRef := uOfD.NewBaseElementReference(hl, LiteralPointerReferenceGetLiteralPointerPointerIndicatedLiteralPointerPointerRefUri)
	core.SetOwningElement(getLiteralPointerPointerIndicatedLiteralPointerPointerRef, getLiteralPointerPointer, hl)
	core.SetLabel(getLiteralPointerPointerIndicatedLiteralPointerPointerRef, "IndicatedLiteralPointerPointerRef", hl)
	core.SetUri(getLiteralPointerPointerIndicatedLiteralPointerPointerRef, LiteralPointerReferenceGetLiteralPointerPointerIndicatedLiteralPointerPointerRefUri, hl)

	// SetReferencedLiteralPointer
	setReferencedLiteralPointer := uOfD.NewElement(hl, LiteralPointerReferenceSetReferencedLiteralPointerUri)
	core.SetLabel(setReferencedLiteralPointer, "SetReferencedLiteralPointer", hl)
	core.SetOwningElement(setReferencedLiteralPointer, literalPointerReferenceFunctions, hl)
	core.SetUri(setReferencedLiteralPointer, LiteralPointerReferenceSetReferencedLiteralPointerUri, hl)
	// SetReferencedLiteralPointer.LiteralPointerReference
	setReferencedLiteralPointerLiteralPointerReference := uOfD.NewLiteralPointerReference(hl, LiteralPointerReferenceSetReferencedLiteralPointerSourceLiteralPointerRefUri)
	core.SetOwningElement(setReferencedLiteralPointerLiteralPointerReference, setReferencedLiteralPointer, hl)
	core.SetLabel(setReferencedLiteralPointerLiteralPointerReference, "SourceLiteralPointerRef", hl)
	core.SetUri(setReferencedLiteralPointerLiteralPointerReference, LiteralPointerReferenceSetReferencedLiteralPointerSourceLiteralPointerRefUri, hl)
	// SetReferencedLiteralPointerTargetLiteralPointerReference
	setReferencedLiteralPointerTargetLiteralPointerReference := uOfD.NewElementReference(hl, LiteralPointerReferenceSetReferencedLiteralPointerModifiedLiteralPointerReferenceRefUri)
	core.SetOwningElement(setReferencedLiteralPointerTargetLiteralPointerReference, setReferencedLiteralPointer, hl)
	core.SetLabel(setReferencedLiteralPointerTargetLiteralPointerReference, "ModifiedLiteralPointerReferenceRef", hl)
	core.SetUri(setReferencedLiteralPointerTargetLiteralPointerReference, LiteralPointerReferenceSetReferencedLiteralPointerModifiedLiteralPointerReferenceRefUri, hl)
}

func literalPointerReferenceFunctionsInit() {
	core.GetCore().AddFunction(LiteralPointerReferenceCreateUri, createLiteralPointerReference)
	core.GetCore().AddFunction(LiteralPointerReferenceGetReferencedLiteralPointerUri, getReferencedLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerReferenceGetLiteralPointerPointerUri, getLiteralPointerPointer)
	core.GetCore().AddFunction(LiteralPointerReferenceSetReferencedLiteralPointerUri, setReferencedLiteralPointer)
}
