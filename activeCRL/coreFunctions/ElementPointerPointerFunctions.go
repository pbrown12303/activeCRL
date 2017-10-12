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

var ElementPointerPointerFunctionsUri string = CoreFunctionsPrefix + "ElementPointerPointerFunctions"

var ElementPointerPointerCreateElementPointerPointerUri string = CoreFunctionsPrefix + "ElementPointerPointer/CreateAbstractElementPointerPointer"
var ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri = CoreFunctionsPrefix + "ElementPointerPointer/CreateAbstractElementPointerPointer/CreatedElementPointerPointerRef"

var ElementPointerPointerGetElementPointerUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointer"
var ElementPointerPointerGetElementPointerSourceElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointer/SourceElementPointerPointerRef"
var ElementPointerPointerGetElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointer/IndicatedElementPointerRef"

var ElementPointerPointerGetElementPointerIdUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerId"
var ElementPointerPointerGetElementPointerIdSourceElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerId/SourceElementPointerPointerRef"
var ElementPointerPointerGetElementPointerIdCreatedLiteralUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerId/CreatedLiteralRef"

var ElementPointerPointerGetElementPointerVersionUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerVersion"
var ElementPointerPointerGetElementPointerVersionSourceElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerVersion/SourceElementPointerPointerRef"
var ElementPointerPointerGetElementPointerVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/GetElementPointerVersion/CreatedLiteralRef"

var ElementPointerPointerSetElementPointerUri string = CoreFunctionsPrefix + "ElementPointerPointer/SetElementPointer"
var ElementPointerPointerSetElementPointerElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/SetElementPointer/ElementPointerRef"
var ElementPointerPointerSetElementPointerModifiedElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerPointer/SetElementPointer/ModifiedElementPointerPointerRef"

func createElementPointerPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createElementPointerPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(element, ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri, hl)
	if createdElementPointerPointerRef == nil {
		createdElementPointerPointerRef = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(createdElementPointerPointerRef, element, hl)
		core.SetName(createdElementPointerPointerRef, "CreatedElementPointerPointerRef", hl)
		rootCreatedElementReference := uOfD.GetBaseElementReferenceWithUri(ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementPointerPointerRef, hl)
		refinement.SetRefinedElement(createdElementPointerPointerRef, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdElementPointerPointer := createdElementPointerPointerRef.GetReferencedBaseElement(hl)
	if createdElementPointerPointer == nil {
		createdElementPointerPointer = uOfD.NewElementPointerPointer(hl)
		createdElementPointerPointerRef.SetReferencedBaseElement(createdElementPointerPointer, hl)
	} else {
		switch createdElementPointerPointer.(type) {
		case core.ElementPointerPointer:
		default:
			// It's the wrong type - create the correct type
			createdElementPointerPointer = uOfD.NewElementPointerPointer(hl)
			createdElementPointerPointerRef.SetReferencedBaseElement(createdElementPointerPointer, hl)
		}
	}
}

func getElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerPointerGetElementPointerUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerSourceElementPointerPointerRefUri, hl)
	if sourceElementPointerPointerRef == nil {
		log.Printf("In GetElementPointer, the SourceElementPointerPointerRef was not found in the replicate")
		return
	}

	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		log.Printf("In GetElementPointer, the IndicatedElementPointerrRef was not found in the replicate")
		return
	}

	indicatedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	untypedSourceElementPointerPointer := sourceElementPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceElementPointerPointer core.ElementPointerPointer
	var sourceElementPointer core.ElementPointer
	if untypedSourceElementPointerPointer != nil {
		switch untypedSourceElementPointerPointer.(type) {
		case core.ElementPointerPointer:
			sourceElementPointerPointer = untypedSourceElementPointerPointer.(core.ElementPointerPointer)
			sourceElementPointer = sourceElementPointerPointer.GetElementPointer(hl)
		}
	}
	if sourceElementPointer != indicatedElementPointer {
		indicatedElementPointerRef.SetReferencedElementPointer(sourceElementPointer, hl)
	}
}

func getElementPointerId(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerPointerGetElementPointerIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerIdSourceElementPointerPointerRefUri, hl)
	if sourceElementPointerPointerRef == nil {
		log.Printf("In GetElementPointerId, the SourceElementPointerPointerRef was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerIdCreatedLiteralUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetElementPointerId, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedSourceElementPointerPointer := sourceElementPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceElementPointerPointer core.ElementPointerPointer
	if untypedSourceElementPointerPointer != nil {
		switch untypedSourceElementPointerPointer.(type) {
		case core.ElementPointerPointer:
			sourceElementPointerPointer = untypedSourceElementPointerPointer.(core.ElementPointerPointer)
		}
	}
	if sourceElementPointerPointer != nil {
		createdLiteral.SetLiteralValue(sourceElementPointerPointer.GetElementPointerId(hl).String(), hl)
	}
}

func getElementPointerVersion(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerPointerGetElementPointerVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerVersionSourceElementPointerPointerRefUri, hl)
	if sourceElementPointerPointerRef == nil {
		log.Printf("In GetElementPointerVersion, the SourceElementPointerPointerRef was not found in the replicate")
		return
	}

	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerPointerGetElementPointerVersionCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		log.Printf("In GetElementPointerVersion, the CreatedLiteralRef was not found in the replicate")
		return
	}

	createdLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, createdLiteralRef, hl)
		createdLiteralRef.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedSourceElementPointerPointer := sourceElementPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceElementPointerPointer core.ElementPointerPointer
	if untypedSourceElementPointerPointer != nil {
		switch untypedSourceElementPointerPointer.(type) {
		case core.ElementPointerPointer:
			sourceElementPointerPointer = untypedSourceElementPointerPointer.(core.ElementPointerPointer)
		}
	}
	if sourceElementPointerPointer != nil {
		stringVersion := strconv.Itoa(sourceElementPointerPointer.GetElementPointerVersion(hl))
		createdLiteral.SetLiteralValue(stringVersion, hl)
	}
}

func setElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerPointerSetElementPointerUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	elementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerPointerSetElementPointerElementPointerRefUri, hl)
	if elementPointerRef == nil {
		log.Printf("In SetElementPointer, the ElementPointerRef was not found in the replicate")
		return
	}

	modifiedElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerPointerSetElementPointerModifiedElementPointerPointerRefUri, hl)
	if modifiedElementPointerPointerRef == nil {
		log.Printf("In SetElementPointer, the ModifiedElementPointerPointerRef was not found in the replicate")
		return
	}

	untypedBaseElement := modifiedElementPointerPointerRef.GetReferencedBaseElement(hl)
	elementPointer := elementPointerRef.GetReferencedElementPointer(hl)
	var modifiedElementPointerPointer core.ElementPointerPointer
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.ElementPointerPointer:
			modifiedElementPointerPointer = untypedBaseElement.(core.ElementPointerPointer)
			modifiedElementPointerPointer.SetElementPointer(elementPointer, hl)
		default:
			log.Printf("In SetElementPointer, the ModifiedElementPointerPointerRef does not point to an ElementPointer")
		}
	}
}

func BuildCoreElementPointerPointerFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// ElementPointerPointerFunctions
	elementPointerPointerFunctions := uOfD.NewElement(hl, ElementPointerPointerFunctionsUri)
	core.SetOwningElement(elementPointerPointerFunctions, coreFunctionsElement, hl)
	core.SetName(elementPointerPointerFunctions, "ElementPointerPointerFunctions", hl)
	core.SetUri(elementPointerPointerFunctions, ElementPointerPointerFunctionsUri, hl)

	// CreateAbstractElementPointerElement
	createElementPointerPointer := uOfD.NewElement(hl, ElementPointerPointerCreateElementPointerPointerUri)
	core.SetOwningElement(createElementPointerPointer, elementPointerPointerFunctions, hl)
	core.SetName(createElementPointerPointer, "CreateElementPointerPointer", hl)
	core.SetUri(createElementPointerPointer, ElementPointerPointerCreateElementPointerPointerUri, hl)
	// CreatedElementReference
	createdElementPointerPointerRef := uOfD.NewBaseElementReference(hl, ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri)
	core.SetOwningElement(createdElementPointerPointerRef, createElementPointerPointer, hl)
	core.SetName(createdElementPointerPointerRef, "CreatedElementPointerdPointerRef", hl)
	core.SetUri(createdElementPointerPointerRef, ElementPointerPointerCreateElementPointerPointerCreatedElementPointerPointerRefUri, hl)

	// GetElementPointer
	getElementPointer := uOfD.NewElement(hl, ElementPointerPointerGetElementPointerUri)
	core.SetName(getElementPointer, "GetElementPointer", hl)
	core.SetOwningElement(getElementPointer, elementPointerPointerFunctions, hl)
	core.SetUri(getElementPointer, ElementPointerPointerGetElementPointerUri, hl)
	// GetElementPointer.SourceReference
	getElementPointerSourceElementPointerPointerRef := uOfD.NewBaseElementReference(hl, ElementPointerPointerGetElementPointerSourceElementPointerPointerRefUri)
	core.SetOwningElement(getElementPointerSourceElementPointerPointerRef, getElementPointer, hl)
	core.SetName(getElementPointerSourceElementPointerPointerRef, "SourceElementPointerPointerRef", hl)
	core.SetUri(getElementPointerSourceElementPointerPointerRef, ElementPointerPointerGetElementPointerSourceElementPointerPointerRefUri, hl)
	// GetElementPointerIndicatedElementPointerRef
	getElementPointerIndicatedElementPointerRef := uOfD.NewElementPointerReference(hl, ElementPointerPointerGetElementPointerIndicatedElementPointerRefUri)
	core.SetOwningElement(getElementPointerIndicatedElementPointerRef, getElementPointer, hl)
	core.SetName(getElementPointerIndicatedElementPointerRef, "IndicatedElementPointerRef", hl)
	core.SetUri(getElementPointerIndicatedElementPointerRef, ElementPointerPointerGetElementPointerIndicatedElementPointerRefUri, hl)

	// GetElementPointerId
	getElementPointerId := uOfD.NewElement(hl, ElementPointerPointerGetElementPointerIdUri)
	core.SetName(getElementPointerId, "GetElementPointerId", hl)
	core.SetOwningElement(getElementPointerId, elementPointerPointerFunctions, hl)
	core.SetUri(getElementPointerId, ElementPointerPointerGetElementPointerIdUri, hl)
	// GetElementPointerId.SourceElementPointerPointerRef
	getElementPointerIdSourceElementPointerPointerRef := uOfD.NewBaseElementReference(hl, ElementPointerPointerGetElementPointerIdSourceElementPointerPointerRefUri)
	core.SetOwningElement(getElementPointerIdSourceElementPointerPointerRef, getElementPointerId, hl)
	core.SetName(getElementPointerIdSourceElementPointerPointerRef, "SourceElementPointerPointerRef", hl)
	core.SetUri(getElementPointerIdSourceElementPointerPointerRef, ElementPointerPointerGetElementPointerIdSourceElementPointerPointerRefUri, hl)
	// GetElementPointerIdCreatedLiteralRef
	getElementPointerIdCreatedLiteralRef := uOfD.NewLiteralReference(hl, ElementPointerPointerGetElementPointerIdCreatedLiteralUri)
	core.SetOwningElement(getElementPointerIdCreatedLiteralRef, getElementPointerId, hl)
	core.SetName(getElementPointerIdCreatedLiteralRef, "CreatedLiteralRef", hl)
	core.SetUri(getElementPointerIdCreatedLiteralRef, ElementPointerPointerGetElementPointerIdCreatedLiteralUri, hl)

	// GetElementPointerVersion
	getElementPointerVersion := uOfD.NewElement(hl, ElementPointerPointerGetElementPointerVersionUri)
	core.SetName(getElementPointerVersion, "GetElementPointerVersion", hl)
	core.SetOwningElement(getElementPointerVersion, elementPointerPointerFunctions, hl)
	core.SetUri(getElementPointerVersion, ElementPointerPointerGetElementPointerVersionUri, hl)
	// GetElementPointerVersion.SourceReference
	getElementPointerVersionSourceReference := uOfD.NewBaseElementReference(hl, ElementPointerPointerGetElementPointerVersionSourceElementPointerPointerRefUri)
	core.SetOwningElement(getElementPointerVersionSourceReference, getElementPointerVersion, hl)
	core.SetName(getElementPointerVersionSourceReference, "SourceElementPointerRef", hl)
	core.SetUri(getElementPointerVersionSourceReference, ElementPointerPointerGetElementPointerVersionSourceElementPointerPointerRefUri, hl)
	// GetElementPointerVersionTargetLiteralReference
	getElementPointerVersionCreatedLiteralRef := uOfD.NewLiteralReference(hl, ElementPointerPointerGetElementPointerVersionCreatedLiteralRefUri)
	core.SetOwningElement(getElementPointerVersionCreatedLiteralRef, getElementPointerVersion, hl)
	core.SetName(getElementPointerVersionCreatedLiteralRef, "CreatedLiteralRef", hl)
	core.SetUri(getElementPointerVersionCreatedLiteralRef, ElementPointerPointerGetElementPointerVersionCreatedLiteralRefUri, hl)

	// SetElementPointer
	setElementPointer := uOfD.NewElement(hl, ElementPointerPointerSetElementPointerUri)
	core.SetName(setElementPointer, "SetElementPointer", hl)
	core.SetOwningElement(setElementPointer, elementPointerPointerFunctions, hl)
	core.SetUri(setElementPointer, ElementPointerPointerSetElementPointerUri, hl)
	// SetElementPointer.ElementPointerRef
	setElementPointerElementPointerRef := uOfD.NewElementPointerReference(hl, ElementPointerPointerSetElementPointerElementPointerRefUri)
	core.SetName(setElementPointerElementPointerRef, "ElementPointerRef", hl)
	core.SetOwningElement(setElementPointerElementPointerRef, setElementPointer, hl)
	core.SetUri(setElementPointerElementPointerRef, ElementPointerPointerSetElementPointerElementPointerRefUri, hl)
	// SetElementPointer.TargetElementPointerPointerReference
	setElementPointerTargetElementPointerPointerReference := uOfD.NewBaseElementReference(hl, ElementPointerPointerSetElementPointerModifiedElementPointerPointerRefUri)
	core.SetName(setElementPointerTargetElementPointerPointerReference, "ModifiedElementPointerPointerRef", hl)
	core.SetOwningElement(setElementPointerTargetElementPointerPointerReference, setElementPointer, hl)
	core.SetUri(setElementPointerTargetElementPointerPointerReference, ElementPointerPointerSetElementPointerModifiedElementPointerPointerRefUri, hl)
}

func elementPointerPointerFunctionsInit() {
	core.GetCore().AddFunction(ElementPointerPointerCreateElementPointerPointerUri, createElementPointerPointer)
	core.GetCore().AddFunction(ElementPointerPointerGetElementPointerUri, getElementPointer)
	core.GetCore().AddFunction(ElementPointerPointerGetElementPointerIdUri, getElementPointerId)
	core.GetCore().AddFunction(ElementPointerPointerGetElementPointerVersionUri, getElementPointerVersion)
	core.GetCore().AddFunction(ElementPointerPointerSetElementPointerUri, setElementPointer)
}
