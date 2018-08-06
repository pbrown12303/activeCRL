// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"strconv"
	"sync"
)

var LiteralPointerPointerFunctionsUri string = CoreFunctionsPrefix + "LiteralPointerPointerFunctions"

var LiteralPointerPointerCreateLiteralPointerPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/CreateAbstractLiteralPointerPointer"
var LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri = CoreFunctionsPrefix + "LiteralPointerPointer/CreateAbstractLiteralPointerPointer/CreatedLiteralPointerPointerRef"

var LiteralPointerPointerGetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer"
var LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer/SourceLiteralPointerPointerRef"
var LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer/IndicatedLiteralPointerRef"

var LiteralPointerPointerGetLiteralPointerIdUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId"
var LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId/SourceLiteralPointerPointerRef"
var LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId/CreatedLiteralRef"

var LiteralPointerPointerGetLiteralPointerVersionUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion"
var LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion/SourceLiteralPointerPointerRef"
var LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion/CreatedLiteralRef"

var LiteralPointerPointerSetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer"
var LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer/LiteralPointerRef"
var LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer/ModifiedLiteralPointerPointerRef"

func createLiteralPointerPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createLiteralPointerPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(element, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri, hl)
	if createdLiteralPointerPointerRef == nil {
		createdLiteralPointerPointerRef = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(createdLiteralPointerPointerRef, element, hl)
		core.SetLabel(createdLiteralPointerPointerRef, "CreatedLiteralPointerPointerRef", hl)
		rootCreatedElementReference := uOfD.GetBaseElementReferenceWithUri(LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerPointerRef, hl)
		refinement.SetRefinedElement(createdLiteralPointerPointerRef, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdLiteralPointerPointer := createdLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	if createdLiteralPointerPointer == nil {
		createdLiteralPointerPointer = uOfD.NewLiteralPointerPointer(hl)
		createdLiteralPointerPointerRef.SetReferencedBaseElement(createdLiteralPointerPointer, hl)
	} else {
		switch createdLiteralPointerPointer.(type) {
		case core.LiteralPointerPointer:
		default:
			// It's the wrong type - create the correct type
			createdLiteralPointerPointer = uOfD.NewLiteralPointerPointer(hl)
			createdLiteralPointerPointerRef.SetReferencedBaseElement(createdLiteralPointerPointer, hl)
		}
	}
}

func getLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		log.Printf("In GetLiteralPointer, the SourceLiteralPointerPointerRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		log.Printf("In GetLiteralPointer, the IndicatedLiteralPointerrRef was not found in the replicate")
		return
	}

	indicatedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	untypedSourceLiteralPointerPointer := sourceLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceLiteralPointerPointer core.LiteralPointerPointer
	var sourceLiteralPointer core.LiteralPointer
	if untypedSourceLiteralPointerPointer != nil {
		switch untypedSourceLiteralPointerPointer.(type) {
		case core.LiteralPointerPointer:
			sourceLiteralPointerPointer = untypedSourceLiteralPointerPointer.(core.LiteralPointerPointer)
			sourceLiteralPointer = sourceLiteralPointerPointer.GetLiteralPointer(hl)
		}
	}
	if sourceLiteralPointer != indicatedLiteralPointer {
		indicatedLiteralPointerRef.SetReferencedLiteralPointer(sourceLiteralPointer, hl)
	}
}

func getLiteralPointerId(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		log.Printf("In GetLiteralPointerId, the SourceLiteralPointerPointerRef was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetLiteralPointerId, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedSourceLiteralPointerPointer := sourceLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceLiteralPointerPointer core.LiteralPointerPointer
	if untypedSourceLiteralPointerPointer != nil {
		switch untypedSourceLiteralPointerPointer.(type) {
		case core.LiteralPointerPointer:
			sourceLiteralPointerPointer = untypedSourceLiteralPointerPointer.(core.LiteralPointerPointer)
		}
	}
	if sourceLiteralPointerPointer != nil {
		createdLiteral.SetLiteralValue(sourceLiteralPointerPointer.GetLiteralPointerId(hl), hl)
	}
}

func getLiteralPointerVersion(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		log.Printf("In GetLiteralPointerVersion, the SourceLiteralPointerPointerRef was not found in the replicate")
		return
	}

	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		log.Printf("In GetLiteralPointerVersion, the CreatedLiteralRef was not found in the replicate")
		return
	}

	createdLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, createdLiteralRef, hl)
		createdLiteralRef.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedSourceLiteralPointerPointer := sourceLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceLiteralPointerPointer core.LiteralPointerPointer
	if untypedSourceLiteralPointerPointer != nil {
		switch untypedSourceLiteralPointerPointer.(type) {
		case core.LiteralPointerPointer:
			sourceLiteralPointerPointer = untypedSourceLiteralPointerPointer.(core.LiteralPointerPointer)
		}
	}
	if sourceLiteralPointerPointer != nil {
		stringVersion := strconv.Itoa(sourceLiteralPointerPointer.GetLiteralPointerVersion(hl))
		createdLiteral.SetLiteralValue(stringVersion, hl)
	}
}

func setLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerPointerSetLiteralPointerUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	literalPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri, hl)
	if literalPointerRef == nil {
		log.Printf("In SetLiteralPointer, the LiteralPointerRef was not found in the replicate")
		return
	}

	modifiedLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri, hl)
	if modifiedLiteralPointerPointerRef == nil {
		log.Printf("In SetLiteralPointer, the ModifiedLiteralPointerPointerRef was not found in the replicate")
		return
	}

	untypedBaseElement := modifiedLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	literalPointer := literalPointerRef.GetReferencedLiteralPointer(hl)
	var modifiedLiteralPointerPointer core.LiteralPointerPointer
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.LiteralPointerPointer:
			modifiedLiteralPointerPointer = untypedBaseElement.(core.LiteralPointerPointer)
			modifiedLiteralPointerPointer.SetLiteralPointer(literalPointer, hl)
		default:
			log.Printf("In SetLiteralPointer, the ModifiedLiteralPointerPointerRef does not point to an LiteralPointer")
		}
	}
}

func BuildCoreLiteralPointerPointerFunctions(coreFunctionsElement core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralPointerPointerFunctions
	literalPointerPointerFunctions := uOfD.NewElement(hl, LiteralPointerPointerFunctionsUri)
	core.SetOwningElement(literalPointerPointerFunctions, coreFunctionsElement, hl)
	core.SetLabel(literalPointerPointerFunctions, "LiteralPointerPointerFunctions", hl)
	core.SetUri(literalPointerPointerFunctions, LiteralPointerPointerFunctionsUri, hl)

	// CreateAbstractLiteralPointerPointer
	createLiteralPointerPointer := uOfD.NewElement(hl, LiteralPointerPointerCreateLiteralPointerPointerUri)
	core.SetOwningElement(createLiteralPointerPointer, literalPointerPointerFunctions, hl)
	core.SetLabel(createLiteralPointerPointer, "CreateLiteralPointerPointer", hl)
	core.SetUri(createLiteralPointerPointer, LiteralPointerPointerCreateLiteralPointerPointerUri, hl)
	// CreatedLiteralPointerPointerReference
	createdLiteralPointerPointerRef := uOfD.NewBaseElementReference(hl, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri)
	core.SetOwningElement(createdLiteralPointerPointerRef, createLiteralPointerPointer, hl)
	core.SetLabel(createdLiteralPointerPointerRef, "CreatedLiteralPointerdPointerRef", hl)
	core.SetUri(createdLiteralPointerPointerRef, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri, hl)

	// GetLiteralPointer
	getLiteralPointer := uOfD.NewElement(hl, LiteralPointerPointerGetLiteralPointerUri)
	core.SetLabel(getLiteralPointer, "GetLiteralPointer", hl)
	core.SetOwningElement(getLiteralPointer, literalPointerPointerFunctions, hl)
	core.SetUri(getLiteralPointer, LiteralPointerPointerGetLiteralPointerUri, hl)
	// GetLiteralPointer.SourceReference
	getLiteralPointerSourceLiteralPointerPointerRef := uOfD.NewBaseElementReference(hl, LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri)
	core.SetOwningElement(getLiteralPointerSourceLiteralPointerPointerRef, getLiteralPointer, hl)
	core.SetLabel(getLiteralPointerSourceLiteralPointerPointerRef, "SourceLiteralPointerPointerRef", hl)
	core.SetUri(getLiteralPointerSourceLiteralPointerPointerRef, LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri, hl)
	// GetLiteralPointerIndicatedLiteralPointerRef
	getLiteralPointerIndicatedLiteralPointerRef := uOfD.NewLiteralPointerReference(hl, LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri)
	core.SetOwningElement(getLiteralPointerIndicatedLiteralPointerRef, getLiteralPointer, hl)
	core.SetLabel(getLiteralPointerIndicatedLiteralPointerRef, "IndicatedLiteralPointerRef", hl)
	core.SetUri(getLiteralPointerIndicatedLiteralPointerRef, LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri, hl)

	// GetLiteralPointerId
	getLiteralPointerId := uOfD.NewElement(hl, LiteralPointerPointerGetLiteralPointerIdUri)
	core.SetLabel(getLiteralPointerId, "GetLiteralPointerId", hl)
	core.SetOwningElement(getLiteralPointerId, literalPointerPointerFunctions, hl)
	core.SetUri(getLiteralPointerId, LiteralPointerPointerGetLiteralPointerIdUri, hl)
	// GetLiteralPointerId.SourceLiteralPointerPointerRef
	getLiteralPointerIdSourceLiteralPointerPointerRef := uOfD.NewBaseElementReference(hl, LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri)
	core.SetOwningElement(getLiteralPointerIdSourceLiteralPointerPointerRef, getLiteralPointerId, hl)
	core.SetLabel(getLiteralPointerIdSourceLiteralPointerPointerRef, "SourceLiteralPointerPointerRef", hl)
	core.SetUri(getLiteralPointerIdSourceLiteralPointerPointerRef, LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri, hl)
	// GetLiteralPointerIdCreatedLiteralRef
	getLiteralPointerIdCreatedLiteralRef := uOfD.NewLiteralReference(hl, LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri)
	core.SetOwningElement(getLiteralPointerIdCreatedLiteralRef, getLiteralPointerId, hl)
	core.SetLabel(getLiteralPointerIdCreatedLiteralRef, "CreatedLiteralRef", hl)
	core.SetUri(getLiteralPointerIdCreatedLiteralRef, LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri, hl)

	// GetLiteralPointerVersion
	getLiteralPointerVersion := uOfD.NewElement(hl, LiteralPointerPointerGetLiteralPointerVersionUri)
	core.SetLabel(getLiteralPointerVersion, "GetLiteralPointerVersion", hl)
	core.SetOwningElement(getLiteralPointerVersion, literalPointerPointerFunctions, hl)
	core.SetUri(getLiteralPointerVersion, LiteralPointerPointerGetLiteralPointerVersionUri, hl)
	// GetLiteralPointerVersion.SourceReference
	getLiteralPointerVersionSourceReference := uOfD.NewBaseElementReference(hl, LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri)
	core.SetOwningElement(getLiteralPointerVersionSourceReference, getLiteralPointerVersion, hl)
	core.SetLabel(getLiteralPointerVersionSourceReference, "SourceLiteralPointerRef", hl)
	core.SetUri(getLiteralPointerVersionSourceReference, LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri, hl)
	// GetLiteralPointerVersionTargetLiteralReference
	getLiteralPointerVersionCreatedLiteralRef := uOfD.NewLiteralReference(hl, LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri)
	core.SetOwningElement(getLiteralPointerVersionCreatedLiteralRef, getLiteralPointerVersion, hl)
	core.SetLabel(getLiteralPointerVersionCreatedLiteralRef, "CreatedLiteralRef", hl)
	core.SetUri(getLiteralPointerVersionCreatedLiteralRef, LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri, hl)

	// SetLiteralPointer
	setLiteralPointer := uOfD.NewElement(hl, LiteralPointerPointerSetLiteralPointerUri)
	core.SetLabel(setLiteralPointer, "SetLiteralPointer", hl)
	core.SetOwningElement(setLiteralPointer, literalPointerPointerFunctions, hl)
	core.SetUri(setLiteralPointer, LiteralPointerPointerSetLiteralPointerUri, hl)
	// SetLiteralPointer.LiteralPointerRef
	setLiteralPointerLiteralPointerRef := uOfD.NewLiteralPointerReference(hl, LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri)
	core.SetLabel(setLiteralPointerLiteralPointerRef, "LiteralPointerRef", hl)
	core.SetOwningElement(setLiteralPointerLiteralPointerRef, setLiteralPointer, hl)
	core.SetUri(setLiteralPointerLiteralPointerRef, LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri, hl)
	// SetLiteralPointer.ModifiedLiteralPointerPointerRef
	setLiteralPointerTargetLiteralPointerPointerReference := uOfD.NewBaseElementReference(hl, LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri)
	core.SetLabel(setLiteralPointerTargetLiteralPointerPointerReference, "ModifiedLiteralPointerPointerRef", hl)
	core.SetOwningElement(setLiteralPointerTargetLiteralPointerPointerReference, setLiteralPointer, hl)
	core.SetUri(setLiteralPointerTargetLiteralPointerPointerReference, LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri, hl)
}

func literalPointerPointerFunctionsInit() {
	core.GetCore().AddFunction(LiteralPointerPointerCreateLiteralPointerPointerUri, createLiteralPointerPointer)
	core.GetCore().AddFunction(LiteralPointerPointerGetLiteralPointerUri, getLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerPointerGetLiteralPointerIdUri, getLiteralPointerId)
	core.GetCore().AddFunction(LiteralPointerPointerGetLiteralPointerVersionUri, getLiteralPointerVersion)
	core.GetCore().AddFunction(LiteralPointerPointerSetLiteralPointerUri, setLiteralPointer)
}
