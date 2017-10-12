// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"sync"
)

var LiteralReferenceFunctionsUri string = CoreFunctionsPrefix + "LiteralReferenceFunctions"

var LiteralReferenceCreateUri string = CoreFunctionsPrefix + "LiteralReference/Create"
var LiteralReferenceCreateCreatedLiteralReferenceRefUri = CoreFunctionsPrefix + "LiteralReference/Create/CreatedLiteralReferenceRef"

var LiteralReferenceGetReferencedLiteralUri string = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral"
var LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral/SourceLiteralReferenceRef"
var LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral/IndicatedLiteralRef"

var LiteralReferenceGetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer"
var LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer/SourceLiteralReferenceRef"
var LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer/IndicatedLiteralPointerRef"

var LiteralReferenceSetReferencedLiteralUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral"
var LiteralReferenceSetReferencedLiteralSourceLiteralRefUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral/SourceLiteralRef"
var LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral/ModifiedLiteralReferenceRef"

func literalReferenceCreateLiteralReference(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(element, LiteralReferenceCreateCreatedLiteralReferenceRefUri, hl)
	if createdLiteralReferenceRef == nil {
		createdLiteralReferenceRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdLiteralReferenceRef, element, hl)
		core.SetName(createdLiteralReferenceRef, "CreatedLiteralReferenceReference", hl)
		rootCreatedLiteralReferenceReference := uOfD.GetElementReferenceWithUri(LiteralReferenceCreateCreatedLiteralReferenceRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralReferenceRef, hl)
		refinement.SetRefinedElement(createdLiteralReferenceRef, hl)
		refinement.SetAbstractElement(rootCreatedLiteralReferenceReference, hl)
	}
	createdLiteralReference := createdLiteralReferenceRef.GetReferencedElement(hl)
	if createdLiteralReference == nil {
		createdLiteralReference = uOfD.NewLiteralReference(hl)
		createdLiteralReferenceRef.SetReferencedElement(createdLiteralReference, hl)
	}
}

func literalReferenceGetReferencedLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralReferenceGetReferencedLiteralUri)
	if original == nil {
		log.Printf("In GetReferencedLiteral the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri, hl)
	if sourceLiteralReferenceRef == nil {
		log.Printf("In GetReferencedLiteral, the SourceReference was not found in the replicate")
		return
	}

	indicatedLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralReference == nil {
		log.Printf("In GetReferencedLiteral, the TargetLiteralReference was not found in the replicate")
		return
	}

	indicatedLiteral := indicatedLiteralReference.GetReferencedLiteral(hl)
	untypedLiteralReference := sourceLiteralReferenceRef.GetReferencedElement(hl)
	var sourceLiteralReference core.LiteralReference
	var sourceLiteral core.Literal
	if untypedLiteralReference != nil {
		switch untypedLiteralReference.(type) {
		case core.LiteralReference:
			sourceLiteralReference = untypedLiteralReference.(core.LiteralReference)
			sourceLiteral = sourceLiteralReference.GetReferencedLiteral(hl)
		}
	}
	if sourceLiteral != indicatedLiteral {
		indicatedLiteralReference.SetReferencedLiteral(sourceLiteral, hl)
	}
}

func literalReferenceGetLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralReferenceGetLiteralPointerUri)
	if original == nil {
		log.Printf("In GetLiteralPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri, hl)
	if sourceLiteralReferenceRef == nil {
		log.Printf("In GetLiteralPointer, the SourceLiteralReferenceRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		log.Printf("In GetLiteralPointer, the IndicatedLiteralPointerRef was not found in the replicate")
		core.Print(replicate, "Replicate: ", hl)
		return
	}

	indicatedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	untypedSourceLiteralReference := sourceLiteralReferenceRef.GetReferencedElement(hl)
	var sourceLiteralPointer core.LiteralPointer
	if untypedSourceLiteralReference != nil {
		switch untypedSourceLiteralReference.(type) {
		case core.LiteralReference:
			sourceLiteralReference := untypedSourceLiteralReference.(core.LiteralReference)
			sourceLiteralPointer = sourceLiteralReference.GetLiteralPointer(hl)
		default:
			log.Printf("In GetLiteralPointer, the SourceElement is not a LiteralReference")
		}
	}
	if sourceLiteralPointer != indicatedLiteralPointer {
		indicatedLiteralPointerRef.SetReferencedLiteralPointer(sourceLiteralPointer, hl)
	}
}

func literalReferenceSetReferencedLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralReferenceSetReferencedLiteralUri)
	if original == nil {
		log.Printf("In SetReferencedLiteral the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralReferenceSetReferencedLiteralSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		log.Printf("In SetReferencedLiteral, the LiteralReference was not found in the replicate")
		return
	}

	modifiedLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri, hl)
	if modifiedLiteralReferenceRef == nil {
		log.Printf("In SetReferencedLiteral, the TargetElement was not found in the replicate")
		return
	}

	sourceLiteral := sourceLiteralRef.GetReferencedLiteral(hl)
	untypedLiteralReference := modifiedLiteralReferenceRef.GetReferencedElement(hl)
	var currentLiteral core.Literal
	var modifiedLiteralReference core.LiteralReference
	if untypedLiteralReference != nil {
		switch untypedLiteralReference.(type) {
		case core.LiteralReference:
			modifiedLiteralReference = untypedLiteralReference.(core.LiteralReference)
			currentLiteral = modifiedLiteralReference.GetReferencedLiteral(hl)
		default:
			log.Printf("In SetReferencedLiteral, the TargetedLiteralReference is not a LiteralReference")
		}
	}
	if sourceLiteral != currentLiteral {
		modifiedLiteralReference.SetReferencedLiteral(sourceLiteral, hl)
	}
}

func BuildCoreLiteralReferenceFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralReferenceFunctions
	elementReferenceFunctions := uOfD.NewElement(hl, LiteralReferenceFunctionsUri)
	core.SetOwningElement(elementReferenceFunctions, coreFunctionsElement, hl)
	core.SetName(elementReferenceFunctions, "LiteralReferenceFunctions", hl)
	core.SetUri(elementReferenceFunctions, LiteralReferenceFunctionsUri, hl)

	// CreateLiteralReference
	literalReferenceCreateLiteralReference := uOfD.NewElement(hl, LiteralReferenceCreateUri)
	core.SetOwningElement(literalReferenceCreateLiteralReference, elementReferenceFunctions, hl)
	core.SetName(literalReferenceCreateLiteralReference, "CreateLiteralReference", hl)
	core.SetUri(literalReferenceCreateLiteralReference, LiteralReferenceCreateUri, hl)
	// CreatedLiteralReference
	createdLiteralReferenceRef := uOfD.NewElementReference(hl, LiteralReferenceCreateCreatedLiteralReferenceRefUri)
	core.SetOwningElement(createdLiteralReferenceRef, literalReferenceCreateLiteralReference, hl)
	core.SetName(createdLiteralReferenceRef, "CreatedLiteralReferenceRef", hl)
	core.SetUri(createdLiteralReferenceRef, LiteralReferenceCreateCreatedLiteralReferenceRefUri, hl)

	// GetReferencedLiteral
	literalReferenceGetReferencedLiteral := uOfD.NewElement(hl, LiteralReferenceGetReferencedLiteralUri)
	core.SetName(literalReferenceGetReferencedLiteral, "GetReferencedLiteral", hl)
	core.SetOwningElement(literalReferenceGetReferencedLiteral, elementReferenceFunctions, hl)
	core.SetUri(literalReferenceGetReferencedLiteral, LiteralReferenceGetReferencedLiteralUri, hl)
	// GetReferencedLiteral.SourceReference
	getElementSourceReference := uOfD.NewElementReference(hl, LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri)
	core.SetOwningElement(getElementSourceReference, literalReferenceGetReferencedLiteral, hl)
	core.SetName(getElementSourceReference, "SourceLiteralReferenceRef", hl)
	core.SetUri(getElementSourceReference, LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri, hl)
	// GetReferencedLiteralTargetLiteralReference
	getElementTargetReference := uOfD.NewLiteralReference(hl, LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri)
	core.SetOwningElement(getElementTargetReference, literalReferenceGetReferencedLiteral, hl)
	core.SetName(getElementTargetReference, "IndicatedLiteralRef", hl)
	core.SetUri(getElementTargetReference, LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri, hl)

	// GetLiteralPointer
	literalReferenceGetLiteralPointer := uOfD.NewElement(hl, LiteralReferenceGetLiteralPointerUri)
	core.SetName(literalReferenceGetLiteralPointer, "GetLiteralPointer", hl)
	core.SetOwningElement(literalReferenceGetLiteralPointer, elementReferenceFunctions, hl)
	core.SetUri(literalReferenceGetLiteralPointer, LiteralReferenceGetLiteralPointerUri, hl)
	// GetLiteralPointer.SourceReference
	getElementPointerSourceReference := uOfD.NewElementReference(hl, LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri)
	core.SetOwningElement(getElementPointerSourceReference, literalReferenceGetLiteralPointer, hl)
	core.SetName(getElementPointerSourceReference, "SourceLiteralReferenceRef", hl)
	core.SetUri(getElementPointerSourceReference, LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri, hl)
	// GetLiteralPointerIndicatedLiteralPointerRef
	getElementPointerIndicatedLiteralPointerRef := uOfD.NewLiteralPointerReference(hl, LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri)
	core.SetOwningElement(getElementPointerIndicatedLiteralPointerRef, literalReferenceGetLiteralPointer, hl)
	core.SetName(getElementPointerIndicatedLiteralPointerRef, "IndicatedLiteralPointerRef", hl)
	core.SetUri(getElementPointerIndicatedLiteralPointerRef, LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri, hl)

	// SetReferencedLiteral
	literalReferenceSetReferencedLiteral := uOfD.NewElement(hl, LiteralReferenceSetReferencedLiteralUri)
	core.SetName(literalReferenceSetReferencedLiteral, "SetReferencedLiteral", hl)
	core.SetOwningElement(literalReferenceSetReferencedLiteral, elementReferenceFunctions, hl)
	core.SetUri(literalReferenceSetReferencedLiteral, LiteralReferenceSetReferencedLiteralUri, hl)
	// SetReferencedLiteral.LiteralReference
	setReferencedElementLiteralReference := uOfD.NewLiteralReference(hl, LiteralReferenceSetReferencedLiteralSourceLiteralRefUri)
	core.SetOwningElement(setReferencedElementLiteralReference, literalReferenceSetReferencedLiteral, hl)
	core.SetName(setReferencedElementLiteralReference, "SourceLiteralRef", hl)
	core.SetUri(setReferencedElementLiteralReference, LiteralReferenceSetReferencedLiteralSourceLiteralRefUri, hl)
	// SetReferencedLiteralTargetLiteralReference
	setReferencedElementTargetLiteralReference := uOfD.NewElementReference(hl, LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri)
	core.SetOwningElement(setReferencedElementTargetLiteralReference, literalReferenceSetReferencedLiteral, hl)
	core.SetName(setReferencedElementTargetLiteralReference, "ModifiedLiteralReferenceRef", hl)
	core.SetUri(setReferencedElementTargetLiteralReference, LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri, hl)
}

func literalReferenceFunctionsInit() {
	core.GetCore().AddFunction(LiteralReferenceCreateUri, literalReferenceCreateLiteralReference)
	core.GetCore().AddFunction(LiteralReferenceGetReferencedLiteralUri, literalReferenceGetReferencedLiteral)
	core.GetCore().AddFunction(LiteralReferenceGetLiteralPointerUri, literalReferenceGetLiteralPointer)
	core.GetCore().AddFunction(LiteralReferenceSetReferencedLiteralUri, literalReferenceSetReferencedLiteral)
}
