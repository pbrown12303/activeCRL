// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"sync"
)

var RefinementFunctionsUri string = CoreFunctionsPrefix + "RefinementFunctions"

var RefinementCreateUri string = CoreFunctionsPrefix + "Refinement/Create"
var RefinementCreateCreatedRefinementRefUri = CoreFunctionsPrefix + "Refinement/Create/CreatedRefinementRef"

var RefinementGetAbstractElementUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElement"
var RefinementGetAbstractElementSourceRefinementRefUri = CoreFunctionsPrefix + "Refinement/GetAbstractElement/SourceRefinementRef"
var RefinementGetAbstractElementIndicatedElementRefUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElement/IndicatedElementRef"

var RefinementGetAbstractElementPointerUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElementPointer"
var RefinementGetAbstractElementPointerSourceRefinementRefUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElementPointer/SourceRefinementRef"
var RefinementGetAbstractElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElementPointer/IndicatedElementPointerRef"

var RefinementGetRefinedElementUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElement"
var RefinementGetRefinedElementSourceRefinementRefUri = CoreFunctionsPrefix + "Refinement/GetRefinedElement/SourceRefinementRef"
var RefinementGetRefinedElementIndicatedElementRefUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElement/IndicatedElementRef"

var RefinementGetRefinedElementPointerUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElementPointer"
var RefinementGetRefinedElementPointerSourceRefinementRefUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElementPointer/SourceRefinementRef"
var RefinementGetRefinedElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElementPointer/IndicatedElementPointerRef"

var RefinementSetAbstractElementUri string = CoreFunctionsPrefix + "Refinement/SetAbstractElement"
var RefinementSetAbstractElementSourceElementRefUri string = CoreFunctionsPrefix + "Refinement/SetAbstractElement/SourceElementRef"
var RefinementSetAbstractElementModifiedRefinementRefUri string = CoreFunctionsPrefix + "Refinement/SetAbstractElement/ModifiedRefinementRef"

var RefinementSetRefinedElementUri string = CoreFunctionsPrefix + "Refinement/SetRefinedElement"
var RefinementSetRefinedElementSourceElementRefUri string = CoreFunctionsPrefix + "Refinement/SetRefinedElement/SourceElementRef"
var RefinementSetRefinedElementModifiedRefinementRefUri string = CoreFunctionsPrefix + "Refinement/SetRefinedElement/ModifiedRefinementRef"

func refinementCreateRefinement(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdRefinementRef := core.GetChildElementReferenceWithAncestorUri(element, RefinementCreateCreatedRefinementRefUri, hl)
	if createdRefinementRef == nil {
		createdRefinementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdRefinementRef, element, hl)
		core.SetName(createdRefinementRef, "CreatedRefinementReference", hl)
		rootCreatedRefinementReference := uOfD.GetElementReferenceWithUri(RefinementCreateCreatedRefinementRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdRefinementRef, hl)
		refinement.SetRefinedElement(createdRefinementRef, hl)
		refinement.SetAbstractElement(rootCreatedRefinementReference, hl)
	}
	createdRefinement := createdRefinementRef.GetReferencedElement(hl)
	if createdRefinement == nil {
		createdRefinement = uOfD.NewRefinement(hl)
		createdRefinementRef.SetReferencedElement(createdRefinement, hl)
	}
}

func refinementGetAbstractElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementGetAbstractElementUri)
	if original == nil {
		log.Printf("In GetAbstractElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		log.Printf("In GetAbstractElement, the SourceRefinementRef was not found in the replicate")
		return
	}

	indicatedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementIndicatedElementRefUri, hl)
	if indicatedElementRef == nil {
		log.Printf("In GetAbstractElement, the IndicatedElementReference was not found in the replicate")
		return
	}

	targetElement := indicatedElementRef.GetReferencedElement(hl)
	untypedRefinement := sourceRefinementRef.GetReferencedElement(hl)
	var sourceRefinement core.Refinement
	var sourceElement core.Element
	if untypedRefinement != nil {
		switch untypedRefinement.(type) {
		case core.Refinement:
			sourceRefinement = untypedRefinement.(core.Refinement)
			sourceElement = sourceRefinement.GetAbstractElement(hl)
		}
	}
	if sourceElement != targetElement {
		indicatedElementRef.SetReferencedElement(sourceElement, hl)
	}
}

func refinementGetAbstractElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementGetAbstractElementPointerUri)
	if original == nil {
		log.Printf("In GetAbstractElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementPointerSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		log.Printf("In GetAbstractElementPointer, the SourceRefinementRef was not found in the replicate")
		return
	}

	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, RefinementGetAbstractElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		log.Printf("In GetAbstractElementPointer, the IndicatedElementPointerRef was not found in the replicate")
		return
	}

	indicatedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	untypedSourceRefinement := sourceRefinementRef.GetReferencedElement(hl)
	var sourceElementPointer core.ElementPointer
	if untypedSourceRefinement != nil {
		switch untypedSourceRefinement.(type) {
		case core.Refinement:
			sourceRefinement := untypedSourceRefinement.(core.Refinement)
			sourceElementPointer = sourceRefinement.GetAbstractElementPointer(hl)
		default:
			log.Printf("In GetAbstractElementPointer, the SourceElement is not a Refinement")
		}
	}
	if sourceElementPointer != indicatedElementPointer {
		indicatedElementPointerRef.SetReferencedElementPointer(sourceElementPointer, hl)
	}
}

func refinementGetRefinedElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementGetRefinedElementUri)
	if original == nil {
		log.Printf("In GetRefinedElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		log.Printf("In GetRefinedElement, the SourceRefinementRef was not found in the replicate")
		return
	}

	indicatedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementIndicatedElementRefUri, hl)
	if indicatedElementRef == nil {
		log.Printf("In GetRefinedElement, the IndicatedElementReference was not found in the replicate")
		return
	}

	targetElement := indicatedElementRef.GetReferencedElement(hl)
	untypedRefinement := sourceRefinementRef.GetReferencedElement(hl)
	var sourceRefinement core.Refinement
	var sourceElement core.Element
	if untypedRefinement != nil {
		switch untypedRefinement.(type) {
		case core.Refinement:
			sourceRefinement = untypedRefinement.(core.Refinement)
			sourceElement = sourceRefinement.GetRefinedElement(hl)
		}
	}
	if sourceElement != targetElement {
		indicatedElementRef.SetReferencedElement(sourceElement, hl)
	}
}

func refinementGetRefinedElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementGetRefinedElementPointerUri)
	if original == nil {
		log.Printf("In GetRefinedElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementPointerSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		log.Printf("In GetRefinedElementPointer, the SourceRefinementRef was not found in the replicate")
		return
	}

	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, RefinementGetRefinedElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		log.Printf("In GetRefinedElementPointer, the IndicatedElementPointerRef was not found in the replicate")
		return
	}

	indicatedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	untypedSourceRefinement := sourceRefinementRef.GetReferencedElement(hl)
	var sourceElementPointer core.ElementPointer
	if untypedSourceRefinement != nil {
		switch untypedSourceRefinement.(type) {
		case core.Refinement:
			sourceRefinement := untypedSourceRefinement.(core.Refinement)
			sourceElementPointer = sourceRefinement.GetRefinedElementPointer(hl)
		default:
			log.Printf("In GetRefinedElementPointer, the SourceElement is not a Refinement")
		}
	}
	if sourceElementPointer != indicatedElementPointer {
		indicatedElementPointerRef.SetReferencedElementPointer(sourceElementPointer, hl)
	}
}

func refinementSetAbstractElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementSetAbstractElementUri)
	if original == nil {
		log.Printf("In SetAbstractElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetAbstractElementSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In SetAbstractElement, the ElementReference was not found in the replicate")
		return
	}

	targetRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetAbstractElementModifiedRefinementRefUri, hl)
	if targetRefinementRef == nil {
		log.Printf("In SetAbstractElement, the TargetRefinementRef was not found in the replicate")
		return
	}

	sourcedElement := sourceElementRef.GetReferencedElement(hl)
	untypedModifiedRefinement := targetRefinementRef.GetReferencedElement(hl)
	var targetElement core.Element
	var modifiedRefinement core.Refinement
	if untypedModifiedRefinement != nil {
		switch untypedModifiedRefinement.(type) {
		case core.Refinement:
			modifiedRefinement = untypedModifiedRefinement.(core.Refinement)
			targetElement = modifiedRefinement.GetAbstractElement(hl)
		default:
			log.Printf("In SetAbstractElement, the TargetedRefinement is not a Refinement")
		}
	}
	if sourcedElement != targetElement {
		modifiedRefinement.SetAbstractElement(sourcedElement, hl)
	}
}

func refinementSetRefinedElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementSetRefinedElementUri)
	if original == nil {
		log.Printf("In SetRefinedElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetRefinedElementSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In SetRefinedElement, the ElementReference was not found in the replicate")
		return
	}

	targetRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetRefinedElementModifiedRefinementRefUri, hl)
	if targetRefinementRef == nil {
		log.Printf("In SetRefinedElement, the TargetRefinementRef was not found in the replicate")
		return
	}

	sourcedElement := sourceElementRef.GetReferencedElement(hl)
	untypedModifiedRefinement := targetRefinementRef.GetReferencedElement(hl)
	var targetElement core.Element
	var modifiedRefinement core.Refinement
	if untypedModifiedRefinement != nil {
		switch untypedModifiedRefinement.(type) {
		case core.Refinement:
			modifiedRefinement = untypedModifiedRefinement.(core.Refinement)
			targetElement = modifiedRefinement.GetRefinedElement(hl)
		default:
			log.Printf("In SetRefinedElement, the TargetedRefinement is not a Refinement")
		}
	}
	if sourcedElement != targetElement {
		modifiedRefinement.SetRefinedElement(sourcedElement, hl)
	}
}

func BuildCoreRefinementFunctions(coreFunctionsElement core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// RefinementFunctions
	refinementFunctions := uOfD.NewElement(hl, RefinementFunctionsUri)
	core.SetOwningElement(refinementFunctions, coreFunctionsElement, hl)
	core.SetName(refinementFunctions, "RefinementFunctions", hl)
	core.SetUri(refinementFunctions, RefinementFunctionsUri, hl)

	// CreateRefinement
	refinementCreateRefinement := uOfD.NewElement(hl, RefinementCreateUri)
	core.SetOwningElement(refinementCreateRefinement, refinementFunctions, hl)
	core.SetName(refinementCreateRefinement, "CreateRefinement", hl)
	core.SetUri(refinementCreateRefinement, RefinementCreateUri, hl)
	// CreatedRefinement
	createdRefinementRef := uOfD.NewElementReference(hl, RefinementCreateCreatedRefinementRefUri)
	core.SetOwningElement(createdRefinementRef, refinementCreateRefinement, hl)
	core.SetName(createdRefinementRef, "CreatedRefinementRef", hl)
	core.SetUri(createdRefinementRef, RefinementCreateCreatedRefinementRefUri, hl)

	// GetAbstractElement
	refinementGetAbstractElement := uOfD.NewElement(hl, RefinementGetAbstractElementUri)
	core.SetName(refinementGetAbstractElement, "GetAbstractElement", hl)
	core.SetOwningElement(refinementGetAbstractElement, refinementFunctions, hl)
	core.SetUri(refinementGetAbstractElement, RefinementGetAbstractElementUri, hl)
	// GetAbstractElement.SourceReference
	getElementSourceReference0 := uOfD.NewElementReference(hl, RefinementGetAbstractElementSourceRefinementRefUri)
	core.SetOwningElement(getElementSourceReference0, refinementGetAbstractElement, hl)
	core.SetName(getElementSourceReference0, "SourceRefinementRef", hl)
	core.SetUri(getElementSourceReference0, RefinementGetAbstractElementSourceRefinementRefUri, hl)
	// GetAbstractElementTargetElementReference
	getElementTargetReference0 := uOfD.NewElementReference(hl, RefinementGetAbstractElementIndicatedElementRefUri)
	core.SetOwningElement(getElementTargetReference0, refinementGetAbstractElement, hl)
	core.SetName(getElementTargetReference0, "IndicatedElementRef", hl)
	core.SetUri(getElementTargetReference0, RefinementGetAbstractElementIndicatedElementRefUri, hl)

	// GetAbstractElementPointer
	refinementGetAbstractElementPointer := uOfD.NewElement(hl, RefinementGetAbstractElementPointerUri)
	core.SetName(refinementGetAbstractElementPointer, "GetAbstractElementPointer", hl)
	core.SetOwningElement(refinementGetAbstractElementPointer, refinementFunctions, hl)
	core.SetUri(refinementGetAbstractElementPointer, RefinementGetAbstractElementPointerUri, hl)
	// GetAbstractElementPointer.SourceReference
	getElementPointerSourceReference0 := uOfD.NewElementReference(hl, RefinementGetAbstractElementPointerSourceRefinementRefUri)
	core.SetOwningElement(getElementPointerSourceReference0, refinementGetAbstractElementPointer, hl)
	core.SetName(getElementPointerSourceReference0, "SourceRefinementRef", hl)
	core.SetUri(getElementPointerSourceReference0, RefinementGetAbstractElementPointerSourceRefinementRefUri, hl)
	// GetAbstractElementPointerIndicatedElementPointerRef
	getElementPointerIndicatedElementPointerRef0 := uOfD.NewElementPointerReference(hl, RefinementGetAbstractElementPointerIndicatedElementPointerRefUri)
	core.SetOwningElement(getElementPointerIndicatedElementPointerRef0, refinementGetAbstractElementPointer, hl)
	core.SetName(getElementPointerIndicatedElementPointerRef0, "IndicatedElementPointerRef", hl)
	core.SetUri(getElementPointerIndicatedElementPointerRef0, RefinementGetAbstractElementPointerIndicatedElementPointerRefUri, hl)

	// GetRefinedElement
	refinementGetRefinedElement := uOfD.NewElement(hl, RefinementGetRefinedElementUri)
	core.SetName(refinementGetRefinedElement, "GetRefinedElement", hl)
	core.SetOwningElement(refinementGetRefinedElement, refinementFunctions, hl)
	core.SetUri(refinementGetRefinedElement, RefinementGetRefinedElementUri, hl)
	// GetRefinedElement.SourceReference
	getElementSourceReference1 := uOfD.NewElementReference(hl, RefinementGetRefinedElementSourceRefinementRefUri)
	core.SetOwningElement(getElementSourceReference1, refinementGetRefinedElement, hl)
	core.SetName(getElementSourceReference1, "SourceRefinementRef", hl)
	core.SetUri(getElementSourceReference1, RefinementGetRefinedElementSourceRefinementRefUri, hl)
	// GetRefinedElementTargetElementReference
	getElementTargetReference1 := uOfD.NewElementReference(hl, RefinementGetRefinedElementIndicatedElementRefUri)
	core.SetOwningElement(getElementTargetReference1, refinementGetRefinedElement, hl)
	core.SetName(getElementTargetReference1, "IndicatedElementRef", hl)
	core.SetUri(getElementTargetReference1, RefinementGetRefinedElementIndicatedElementRefUri, hl)

	// GetRefinedElementPointer
	refinementGetRefinedElementPointer := uOfD.NewElement(hl, RefinementGetRefinedElementPointerUri)
	core.SetName(refinementGetRefinedElementPointer, "GetRefinedElementPointer", hl)
	core.SetOwningElement(refinementGetRefinedElementPointer, refinementFunctions, hl)
	core.SetUri(refinementGetRefinedElementPointer, RefinementGetRefinedElementPointerUri, hl)
	// GetRefinedElementPointer.SourceReference
	getElementPointerSourceReference1 := uOfD.NewElementReference(hl, RefinementGetRefinedElementPointerSourceRefinementRefUri)
	core.SetOwningElement(getElementPointerSourceReference1, refinementGetRefinedElementPointer, hl)
	core.SetName(getElementPointerSourceReference1, "SourceRefinementRef", hl)
	core.SetUri(getElementPointerSourceReference1, RefinementGetRefinedElementPointerSourceRefinementRefUri, hl)
	// GetRefinedElementPointerIndicatedElementPointerRef
	getElementPointerIndicatedElementPointerRef1 := uOfD.NewElementPointerReference(hl, RefinementGetRefinedElementPointerIndicatedElementPointerRefUri)
	core.SetOwningElement(getElementPointerIndicatedElementPointerRef1, refinementGetRefinedElementPointer, hl)
	core.SetName(getElementPointerIndicatedElementPointerRef1, "IndicatedElementPointerRef", hl)
	core.SetUri(getElementPointerIndicatedElementPointerRef1, RefinementGetRefinedElementPointerIndicatedElementPointerRefUri, hl)

	// SetAbstractElement
	refinementSetAbstractElement := uOfD.NewElement(hl, RefinementSetAbstractElementUri)
	core.SetName(refinementSetAbstractElement, "SetAbstractElement", hl)
	core.SetOwningElement(refinementSetAbstractElement, refinementFunctions, hl)
	core.SetUri(refinementSetAbstractElement, RefinementSetAbstractElementUri, hl)
	// SetAbstractElement.ElementReference
	setReferencedElementElementReference0 := uOfD.NewElementReference(hl, RefinementSetAbstractElementSourceElementRefUri)
	core.SetOwningElement(setReferencedElementElementReference0, refinementSetAbstractElement, hl)
	core.SetName(setReferencedElementElementReference0, "SourceElementRef", hl)
	core.SetUri(setReferencedElementElementReference0, RefinementSetAbstractElementSourceElementRefUri, hl)
	// SetAbstractElementTargetRefinement
	setReferencedElementTargetRefinement0 := uOfD.NewElementReference(hl, RefinementSetAbstractElementModifiedRefinementRefUri)
	core.SetOwningElement(setReferencedElementTargetRefinement0, refinementSetAbstractElement, hl)
	core.SetName(setReferencedElementTargetRefinement0, "ModifiedRefinementRef", hl)
	core.SetUri(setReferencedElementTargetRefinement0, RefinementSetAbstractElementModifiedRefinementRefUri, hl)

	// SetRefinedElement
	refinementSetRefinedElement := uOfD.NewElement(hl, RefinementSetRefinedElementUri)
	core.SetName(refinementSetRefinedElement, "SetRefinedElement", hl)
	core.SetOwningElement(refinementSetRefinedElement, refinementFunctions, hl)
	core.SetUri(refinementSetRefinedElement, RefinementSetRefinedElementUri, hl)
	// SetRefinedElement.ElementReference
	setReferencedElementElementReference1 := uOfD.NewElementReference(hl, RefinementSetRefinedElementSourceElementRefUri)
	core.SetOwningElement(setReferencedElementElementReference1, refinementSetRefinedElement, hl)
	core.SetName(setReferencedElementElementReference1, "SourceElementRef", hl)
	core.SetUri(setReferencedElementElementReference1, RefinementSetRefinedElementSourceElementRefUri, hl)
	// SetRefinedElementTargetRefinement
	setReferencedElementTargetRefinement1 := uOfD.NewElementReference(hl, RefinementSetRefinedElementModifiedRefinementRefUri)
	core.SetOwningElement(setReferencedElementTargetRefinement1, refinementSetRefinedElement, hl)
	core.SetName(setReferencedElementTargetRefinement1, "ModifiedRefinementRef", hl)
	core.SetUri(setReferencedElementTargetRefinement1, RefinementSetRefinedElementModifiedRefinementRefUri, hl)
}

func refinementFunctionsInit() {
	core.GetCore().AddFunction(RefinementCreateUri, refinementCreateRefinement)
	core.GetCore().AddFunction(RefinementGetAbstractElementUri, refinementGetAbstractElement)
	core.GetCore().AddFunction(RefinementGetAbstractElementPointerUri, refinementGetAbstractElementPointer)
	core.GetCore().AddFunction(RefinementGetRefinedElementUri, refinementGetRefinedElement)
	core.GetCore().AddFunction(RefinementGetRefinedElementPointerUri, refinementGetRefinedElementPointer)
	core.GetCore().AddFunction(RefinementSetAbstractElementUri, refinementSetAbstractElement)
	core.GetCore().AddFunction(RefinementSetRefinedElementUri, refinementSetRefinedElement)
}
