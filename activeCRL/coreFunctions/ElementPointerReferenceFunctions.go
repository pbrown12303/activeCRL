// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"sync"
)

var ElementPointerReferenceFunctionsUri string = CoreFunctionsPrefix + "ElementPointerReferenceFunctions"

var ElementPointerReferenceCreateUri string = CoreFunctionsPrefix + "ElementPointerReference/Create"
var ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri = CoreFunctionsPrefix + "ElementPointerReference/Create/CreatedElementPointerReferenceRef"

var ElementPointerReferenceGetReferencedElementPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer"
var ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer/SourceElementPointerReferenceRef"
var ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer/IndicatedElementPointerRef"

var ElementPointerReferenceGetElementPointerPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer"
var ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer/SourceElementPointerReferenceRef"
var ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer/IndicatedElementPointerPointerRef"

var ElementPointerReferenceSetReferencedElementPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer"
var ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer/SourceElementPointerRef"
var ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer/ModifiedElementPointerReferenceRef"

func createElementPointerReference(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementPointerReferenceReference := core.GetChildElementReferenceWithAncestorUri(element, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri, hl)
	if createdElementPointerReferenceReference == nil {
		createdElementPointerReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementPointerReferenceReference, element, hl)
		core.SetLabel(createdElementPointerReferenceReference, "CreatedElementPointerReferenceReference", hl)
		rootCreatedElementPointerReferenceReference := uOfD.GetElementReferenceWithUri(ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementPointerReferenceReference, hl)
		refinement.SetRefinedElement(createdElementPointerReferenceReference, hl)
		refinement.SetAbstractElement(rootCreatedElementPointerReferenceReference, hl)
	}
	createdElementPointerReference := createdElementPointerReferenceReference.GetReferencedElement(hl)
	if createdElementPointerReference == nil {
		createdElementPointerReference = uOfD.NewElementPointerReference(hl)
		createdElementPointerReferenceReference.SetReferencedElement(createdElementPointerReference, hl)
	}
}

func getReferencedElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerReferenceGetReferencedElementPointerUri)
	if original == nil {
		log.Printf("In GetReferencedElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetReferencedElementPointer, the SourceReference was not found in the replicate")
		return
	}

	targetElementPointerReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri, hl)
	if targetElementPointerReference == nil {
		log.Printf("In GetReferencedElementPointer, the TargetElementPointerReference was not found in the replicate")
		return
	}

	targetElementPointer := targetElementPointerReference.GetReferencedElementPointer(hl)
	untypedElementPointerReference := sourceReference.GetReferencedElement(hl)
	var sourceElementPointerReference core.ElementPointerReference
	var sourceElementPointer core.ElementPointer
	if untypedElementPointerReference != nil {
		switch untypedElementPointerReference.(type) {
		case core.ElementPointerReference:
			sourceElementPointerReference = untypedElementPointerReference.(core.ElementPointerReference)
			sourceElementPointer = sourceElementPointerReference.GetReferencedElementPointer(hl)
		}
	}
	if sourceElementPointer != targetElementPointer {
		targetElementPointerReference.SetReferencedElementPointer(sourceElementPointer, hl)
	}
}

func getElementPointerPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In getElementPointerPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerReferenceGetElementPointerPointerUri)
	if original == nil {
		log.Printf("In GetElementPointerPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetElementPointerPointer, the SourceElementPointerReferenceRef was not found in the replicate")
		return
	}

	indicatedElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri, hl)
	if indicatedElementPointerPointerRef == nil {
		log.Printf("In GetElementPointerPointer, the IndicatedElementPointerPointerRef was not found in the replicate")
		return
	}

	indicatedElementPointerPointer := indicatedElementPointerPointerRef.GetReferencedBaseElement(hl)
	untypedSourceElementPointerReference := sourceReference.GetReferencedElement(hl)
	var sourceElementPointerPointer core.ElementPointerPointer
	if untypedSourceElementPointerReference != nil {
		switch untypedSourceElementPointerReference.(type) {
		case core.ElementPointerReference:
			sourceElementPointerReference := untypedSourceElementPointerReference.(core.ElementPointerReference)
			sourceElementPointerPointer = sourceElementPointerReference.GetElementPointerPointer(hl)
		default:
			log.Printf("In GetElementPointerPointer, the SourceElement is not a ElementPointerReference")
		}
	}
	if sourceElementPointerPointer != indicatedElementPointerPointer {
		indicatedElementPointerPointerRef.SetReferencedBaseElement(sourceElementPointerPointer, hl)
	}
}

func setReferencedElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerReferenceSetReferencedElementPointerUri)
	if original == nil {
		log.Printf("In SetReferencedElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	baseElementReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri, hl)
	if baseElementReference == nil {
		log.Printf("In SetReferencedElementPointer, the ElementPointerReference was not found in the replicate")
		return
	}

	targetElementPointerReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri, hl)
	if targetElementPointerReference == nil {
		log.Printf("In SetReferencedElementPointer, the TargetElementPointer was not found in the replicate")
		return
	}

	sourcedElementPointer := baseElementReference.GetReferencedElementPointer(hl)
	untypedTargetedElement := targetElementPointerReference.GetReferencedElement(hl)
	var targetedElementPointer core.ElementPointer
	var targetedElementPointerReference core.ElementPointerReference
	if untypedTargetedElement != nil {
		switch untypedTargetedElement.(type) {
		case core.ElementPointerReference:
			targetedElementPointerReference = untypedTargetedElement.(core.ElementPointerReference)
			targetedElementPointer = targetedElementPointerReference.GetReferencedElementPointer(hl)
		default:
			log.Printf("In SetReferencedElementPointer, the TargetedElementPointerReference is not a ElementPointerReference")
		}
	}
	if sourcedElementPointer != targetedElementPointer {
		targetedElementPointerReference.SetReferencedElementPointer(sourcedElementPointer, hl)
	}
}

func BuildCoreElementPointerReferenceFunctions(coreFunctionsElement core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// ElementPointerReferenceFunctions
	elementPointerReferenceFunctions := uOfD.NewElement(hl, ElementPointerReferenceFunctionsUri)
	core.SetOwningElement(elementPointerReferenceFunctions, coreFunctionsElement, hl)
	core.SetLabel(elementPointerReferenceFunctions, "ElementPointerReferenceFunctions", hl)
	core.SetUri(elementPointerReferenceFunctions, ElementPointerReferenceFunctionsUri, hl)

	// CreateElementPointerReference
	createElementPointerReference := uOfD.NewElement(hl, ElementPointerReferenceCreateUri)
	core.SetOwningElement(createElementPointerReference, elementPointerReferenceFunctions, hl)
	core.SetLabel(createElementPointerReference, "CreateElementPointerReference", hl)
	core.SetUri(createElementPointerReference, ElementPointerReferenceCreateUri, hl)
	// CreatedElementPointerReference
	createdElementPointerReferenceReference := uOfD.NewElementReference(hl, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri)
	core.SetOwningElement(createdElementPointerReferenceReference, createElementPointerReference, hl)
	core.SetLabel(createdElementPointerReferenceReference, "CreatedElementPointerReferenceRef", hl)
	core.SetUri(createdElementPointerReferenceReference, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri, hl)

	// GetReferencedElementPointer
	getReferencedElementPointer := uOfD.NewElement(hl, ElementPointerReferenceGetReferencedElementPointerUri)
	core.SetLabel(getReferencedElementPointer, "GetReferencedElementPointer", hl)
	core.SetOwningElement(getReferencedElementPointer, elementPointerReferenceFunctions, hl)
	core.SetUri(getReferencedElementPointer, ElementPointerReferenceGetReferencedElementPointerUri, hl)
	// GetReferencedElementPointer.SourceReference
	getElementPointerSourceReference := uOfD.NewElementReference(hl, ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri)
	core.SetOwningElement(getElementPointerSourceReference, getReferencedElementPointer, hl)
	core.SetLabel(getElementPointerSourceReference, "SourceElementPointerReferenceRef", hl)
	core.SetUri(getElementPointerSourceReference, ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri, hl)
	// GetReferencedElementPointerTargetElementPointerReference
	getElementPointerTargetReference := uOfD.NewElementPointerReference(hl, ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri)
	core.SetOwningElement(getElementPointerTargetReference, getReferencedElementPointer, hl)
	core.SetLabel(getElementPointerTargetReference, "IndicatedElementPointerRef", hl)
	core.SetUri(getElementPointerTargetReference, ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri, hl)

	// GetElementPointerPointer
	getElementPointerPointer := uOfD.NewElement(hl, ElementPointerReferenceGetElementPointerPointerUri)
	core.SetLabel(getElementPointerPointer, "GetElementPointerPointer", hl)
	core.SetOwningElement(getElementPointerPointer, elementPointerReferenceFunctions, hl)
	core.SetUri(getElementPointerPointer, ElementPointerReferenceGetElementPointerPointerUri, hl)
	// GetElementPointerPointer.SourceReference
	getElementPointerPointerSourceReference := uOfD.NewElementReference(hl, ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri)
	core.SetOwningElement(getElementPointerPointerSourceReference, getElementPointerPointer, hl)
	core.SetLabel(getElementPointerPointerSourceReference, "SourceElementPointerReferenceRef", hl)
	core.SetUri(getElementPointerPointerSourceReference, ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri, hl)
	// GetReferencedElementPointerTargetElementPointerReference
	getElementPointerPointerIndicatedElementPointerPointerRef := uOfD.NewBaseElementReference(hl, ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri)
	core.SetOwningElement(getElementPointerPointerIndicatedElementPointerPointerRef, getElementPointerPointer, hl)
	core.SetLabel(getElementPointerPointerIndicatedElementPointerPointerRef, "IndicatedElementPointerPointerRef", hl)
	core.SetUri(getElementPointerPointerIndicatedElementPointerPointerRef, ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri, hl)

	// SetReferencedElementPointer
	setReferencedElementPointer := uOfD.NewElement(hl, ElementPointerReferenceSetReferencedElementPointerUri)
	core.SetLabel(setReferencedElementPointer, "SetReferencedElementPointer", hl)
	core.SetOwningElement(setReferencedElementPointer, elementPointerReferenceFunctions, hl)
	core.SetUri(setReferencedElementPointer, ElementPointerReferenceSetReferencedElementPointerUri, hl)
	// SetReferencedElementPointer.ElementPointerReference
	setReferencedElementPointerElementPointerReference := uOfD.NewElementPointerReference(hl, ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri)
	core.SetOwningElement(setReferencedElementPointerElementPointerReference, setReferencedElementPointer, hl)
	core.SetLabel(setReferencedElementPointerElementPointerReference, "SourceElementPointerRef", hl)
	core.SetUri(setReferencedElementPointerElementPointerReference, ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri, hl)
	// SetReferencedElementPointerTargetElementPointerReference
	setReferencedElementPointerTargetElementPointerReference := uOfD.NewElementReference(hl, ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri)
	core.SetOwningElement(setReferencedElementPointerTargetElementPointerReference, setReferencedElementPointer, hl)
	core.SetLabel(setReferencedElementPointerTargetElementPointerReference, "ModifiedElementPointerReferenceRef", hl)
	core.SetUri(setReferencedElementPointerTargetElementPointerReference, ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri, hl)
}

func elementPointerReferenceFunctionsInit() {
	core.GetCore().AddFunction(ElementPointerReferenceCreateUri, createElementPointerReference)
	core.GetCore().AddFunction(ElementPointerReferenceGetReferencedElementPointerUri, getReferencedElementPointer)
	core.GetCore().AddFunction(ElementPointerReferenceGetElementPointerPointerUri, getElementPointerPointer)
	core.GetCore().AddFunction(ElementPointerReferenceSetReferencedElementPointerUri, setReferencedElementPointer)
}
