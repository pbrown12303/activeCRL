// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"sync"
)

var BaseElementReferenceFunctionsUri string = CoreFunctionsPrefix + "BaseElementReferenceFunctions"

var BaseElementReferenceCreateUri string = CoreFunctionsPrefix + "BaseElementReference/Create"
var BaseElementReferenceCreateCreatedBaseElementReferenceRefUri = CoreFunctionsPrefix + "BaseElementReference/Create/CreatedBaseElementReferenceRef"

var BaseElementReferenceGetBaseElementPointerUri string = CoreFunctionsPrefix + "BaseElementReference/GetBaseElementPointer"
var BaseElementReferenceGetBaseElementPointerSourceBaseElementReferenceRefUri string = CoreFunctionsPrefix + "BaseElementReference/GetBaseElementPointer/SourceBaseElementReferenceRef"
var BaseElementReferenceGetBaseElementPointerIndicatedBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementReference/GetBaseElementPointer/IndicatedBaseElementPointerRef"

var BaseElementReferenceGetReferencedBaseElementUri string = CoreFunctionsPrefix + "BaseElementReference/GetReferencedBaseElement"
var BaseElementReferenceGetReferencedBaseElementSourceBaseElementReferenceRefUri string = CoreFunctionsPrefix + "BaseElementReference/GetReferencedBaseElement/SourceBaseElementReferenceRef"
var BaseElementReferenceGetReferencedBaseElementIndicatedBaseElementRefUri string = CoreFunctionsPrefix + "BaseElementReference/GetReferencedBaseElement/IndicatedBaseElementRef"

var BaseElementReferenceSetReferencedBaseElementUri string = CoreFunctionsPrefix + "BaseElementReference/SetReferencedBaseElement"
var BaseElementReferenceSetReferencedBaseElementSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElementReference/SetReferencedBaseElement/SourceBaseElementRef"
var BaseElementReferenceSetReferencedBaseElementModifiedBaseElementReferenceRefUri string = CoreFunctionsPrefix + "BaseElementReference/SetReferencedBaseElement/ModifiedBaseElementReferenceRef"

func createBaseElementReference(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdBaseElementReferenceReference := core.GetChildElementReferenceWithAncestorUri(element, BaseElementReferenceCreateCreatedBaseElementReferenceRefUri, hl)
	if createdBaseElementReferenceReference == nil {
		createdBaseElementReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdBaseElementReferenceReference, element, hl)
		core.SetName(createdBaseElementReferenceReference, "CreatedBaseElementReferenceReference", hl)
		rootCreatedBaseElementReferenceReference := uOfD.GetElementReferenceWithUri(BaseElementReferenceCreateCreatedBaseElementReferenceRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdBaseElementReferenceReference, hl)
		refinement.SetRefinedElement(createdBaseElementReferenceReference, hl)
		refinement.SetAbstractElement(rootCreatedBaseElementReferenceReference, hl)
	}
	createdBaseElementReference := createdBaseElementReferenceReference.GetReferencedElement(hl)
	if createdBaseElementReference == nil {
		createdBaseElementReference = uOfD.NewBaseElementReference(hl)
		createdBaseElementReferenceReference.SetReferencedElement(createdBaseElementReference, hl)
	}
}

func getBaseElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementReferenceGetBaseElementPointerUri)
	if original == nil {
		log.Printf("In GetBaseElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementReferenceGetBaseElementPointerSourceBaseElementReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetBaseElementPointer, the SourceReference was not found in the replicate")
		return
	}

	targetBaseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementReferenceGetBaseElementPointerIndicatedBaseElementPointerRefUri, hl)
	if targetBaseElementReference == nil {
		log.Printf("In GetBaseElementPointer, the TargetBaseElementReference was not found in the replicate")
		return
	}

	targetBaseElementPointer := targetBaseElementReference.GetReferencedBaseElement(hl)
	untypedSourceElement := sourceReference.GetReferencedElement(hl)
	var sourceBaseElementPointer core.BaseElementPointer
	if untypedSourceElement != nil {
		switch untypedSourceElement.(type) {
		case core.BaseElementReference:
			sourceBaseElementReference := untypedSourceElement.(core.BaseElementReference)
			sourceBaseElementPointer = sourceBaseElementReference.GetBaseElementPointer(hl)
		default:
			log.Printf("In GetBaseElementPointer, the SourceElement is not a BaseElementReference")
		}
	}
	if sourceBaseElementPointer != targetBaseElementPointer {
		targetBaseElementReference.SetReferencedBaseElement(sourceBaseElementPointer, hl)
	}
}

func getReferencedBaseElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementReferenceGetReferencedBaseElementUri)
	if original == nil {
		log.Printf("In GetReferencedBaseElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementReferenceGetReferencedBaseElementSourceBaseElementReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetReferencedBaseElement, the SourceReference was not found in the replicate")
		return
	}

	targetBaseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementReferenceGetReferencedBaseElementIndicatedBaseElementRefUri, hl)
	if targetBaseElementReference == nil {
		log.Printf("In GetReferencedBaseElement, the TargetBaseElement was not found in the replicate")
		return
	}

	targetBaseElement := targetBaseElementReference.GetReferencedBaseElement(hl)
	untypedSourceElement := sourceReference.GetReferencedElement(hl)
	var sourceBaseElement core.BaseElement
	if untypedSourceElement != nil {
		switch untypedSourceElement.(type) {
		case core.BaseElementReference:
			sourceBaseElementReference := untypedSourceElement.(core.BaseElementReference)
			sourceBaseElement = sourceBaseElementReference.GetReferencedBaseElement(hl)
		default:
			log.Printf("In GetReferencedBaseElement, the SourceElement is not a BaseElementReference")
		}
	}
	if sourceBaseElement != targetBaseElement {
		targetBaseElementReference.SetReferencedBaseElement(sourceBaseElement, hl)
	}
}

func setReferencedBaseElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementReferenceSetReferencedBaseElementUri)
	if original == nil {
		log.Printf("In SetReferencedBaseElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	baseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementReferenceSetReferencedBaseElementSourceBaseElementRefUri, hl)
	if baseElementReference == nil {
		log.Printf("In SetReferencedBaseElement, the BaseElementReference was not found in the replicate")
		return
	}

	targetBaseElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementReferenceSetReferencedBaseElementModifiedBaseElementReferenceRefUri, hl)
	if targetBaseElementReference == nil {
		log.Printf("In SetReferencedBaseElement, the TargetBaseElement was not found in the replicate")
		return
	}

	sourcedBaseElement := baseElementReference.GetReferencedBaseElement(hl)
	untypedTargetedElement := targetBaseElementReference.GetReferencedElement(hl)
	var targetedBaseElement core.BaseElement
	var targetedBaseElementReference core.BaseElementReference
	if untypedTargetedElement != nil {
		switch untypedTargetedElement.(type) {
		case core.BaseElementReference:
			targetedBaseElementReference = untypedTargetedElement.(core.BaseElementReference)
			targetedBaseElement = targetedBaseElementReference.GetReferencedBaseElement(hl)
		default:
			log.Printf("In SetReferencedBaseElement, the TargetedBaseElementReference is not a BaseElementReference")
		}
	}
	if sourcedBaseElement != targetedBaseElement {
		targetedBaseElementReference.SetReferencedBaseElement(sourcedBaseElement, hl)
	}
}

func UpdateRecoveredCoreBaseElementReferenceFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// BaseElementReferenceFunctions
	baseElementReferenceFunctions := uOfD.GetElementWithUri(BaseElementReferenceFunctionsUri)
	if baseElementReferenceFunctions == nil {
		baseElementReferenceFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(baseElementReferenceFunctions, coreFunctionsElement, hl)
		core.SetName(baseElementReferenceFunctions, "BaseElementReferenceFunctions", hl)
		core.SetUri(baseElementReferenceFunctions, BaseElementReferenceFunctionsUri, hl)
	}

	// CreateBaseElementReference
	createBaseElementReference := uOfD.GetElementWithUri(BaseElementReferenceCreateUri)
	if createBaseElementReference == nil {
		createBaseElementReference = uOfD.NewElement(hl)
		core.SetOwningElement(createBaseElementReference, baseElementReferenceFunctions, hl)
		core.SetName(createBaseElementReference, "CreateBaseElementReference", hl)
		core.SetUri(createBaseElementReference, BaseElementReferenceCreateUri, hl)
	}
	// CreatedBaseElementReference
	createdBaseElementReferenceReference := core.GetChildElementReferenceWithUri(createBaseElementReference, BaseElementReferenceCreateCreatedBaseElementReferenceRefUri, hl)
	if createdBaseElementReferenceReference == nil {
		createdBaseElementReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdBaseElementReferenceReference, createBaseElementReference, hl)
		core.SetName(createdBaseElementReferenceReference, "CreatedBaseElementReferenceRef", hl)
		core.SetUri(createdBaseElementReferenceReference, BaseElementReferenceCreateCreatedBaseElementReferenceRefUri, hl)
	}

	// GetBaseElementPointer
	getBaseElementPointer := uOfD.GetElementWithUri(BaseElementReferenceGetBaseElementPointerUri)
	if getBaseElementPointer == nil {
		getBaseElementPointer = uOfD.NewElement(hl)
		core.SetName(getBaseElementPointer, "GetBaseElementPointer", hl)
		core.SetOwningElement(getBaseElementPointer, baseElementReferenceFunctions, hl)
		core.SetUri(getBaseElementPointer, BaseElementReferenceGetBaseElementPointerUri, hl)
	}
	// GetBaseElement.SourceReference
	getBaseElementPointerSourceReference := core.GetChildElementReferenceWithUri(getBaseElementPointer, BaseElementReferenceGetBaseElementPointerSourceBaseElementReferenceRefUri, hl)
	if getBaseElementPointerSourceReference == nil {
		getBaseElementPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getBaseElementPointerSourceReference, getBaseElementPointer, hl)
		core.SetName(getBaseElementPointerSourceReference, "SourceBaseElementReferenceRef", hl)
		core.SetUri(getBaseElementPointerSourceReference, BaseElementReferenceGetBaseElementPointerSourceBaseElementReferenceRefUri, hl)
	}
	// GetBaseElementTargetBaseElementReference
	getBaseElementPointerTargetReference := core.GetChildBaseElementReferenceWithUri(getBaseElementPointer, BaseElementReferenceGetBaseElementPointerIndicatedBaseElementPointerRefUri, hl)
	if getBaseElementPointerTargetReference == nil {
		getBaseElementPointerTargetReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementPointerTargetReference, getBaseElementPointer, hl)
		core.SetName(getBaseElementPointerTargetReference, "IndicatedBaseElementPointerRef", hl)
		core.SetUri(getBaseElementPointerTargetReference, BaseElementReferenceGetBaseElementPointerIndicatedBaseElementPointerRefUri, hl)
	}

	// GetReferencedBaseElement
	getReferencedBaseElement := uOfD.GetElementWithUri(BaseElementReferenceGetReferencedBaseElementUri)
	if getReferencedBaseElement == nil {
		getReferencedBaseElement = uOfD.NewElement(hl)
		core.SetName(getReferencedBaseElement, "GetReferencedBaseElement", hl)
		core.SetOwningElement(getReferencedBaseElement, baseElementReferenceFunctions, hl)
		core.SetUri(getReferencedBaseElement, BaseElementReferenceGetReferencedBaseElementUri, hl)
	}
	// GetBaseElement.SourceReference
	getReferencedBaseElementSourceReference := core.GetChildElementReferenceWithUri(getReferencedBaseElement, BaseElementReferenceGetReferencedBaseElementSourceBaseElementReferenceRefUri, hl)
	if getReferencedBaseElementSourceReference == nil {
		getReferencedBaseElementSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getReferencedBaseElementSourceReference, getReferencedBaseElement, hl)
		core.SetName(getReferencedBaseElementSourceReference, "SourceBaseElementReferenceRef", hl)
		core.SetUri(getReferencedBaseElementSourceReference, BaseElementReferenceGetReferencedBaseElementSourceBaseElementReferenceRefUri, hl)
	}
	// GetBaseElementTargetBaseElementReference
	getReferencedBaseElementTargetReference := core.GetChildBaseElementReferenceWithUri(getReferencedBaseElement, BaseElementReferenceGetReferencedBaseElementIndicatedBaseElementRefUri, hl)
	if getReferencedBaseElementTargetReference == nil {
		getReferencedBaseElementTargetReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getReferencedBaseElementTargetReference, getReferencedBaseElement, hl)
		core.SetName(getReferencedBaseElementTargetReference, "IndicatedBaseElementRef", hl)
		core.SetUri(getReferencedBaseElementTargetReference, BaseElementReferenceGetReferencedBaseElementIndicatedBaseElementRefUri, hl)
	}

	// SetReferencedBaseElement
	setReferencedBaseElement := uOfD.GetElementWithUri(BaseElementReferenceSetReferencedBaseElementUri)
	if setReferencedBaseElement == nil {
		setReferencedBaseElement = uOfD.NewElement(hl)
		core.SetName(setReferencedBaseElement, "SetReferencedBaseElement", hl)
		core.SetOwningElement(setReferencedBaseElement, baseElementReferenceFunctions, hl)
		core.SetUri(setReferencedBaseElement, BaseElementReferenceSetReferencedBaseElementUri, hl)
	}
	// GetBaseElement.BaseElementReference
	setReferencedBaseElementBaseElementReference := core.GetChildBaseElementReferenceWithUri(setReferencedBaseElement, BaseElementReferenceSetReferencedBaseElementSourceBaseElementRefUri, hl)
	if setReferencedBaseElementBaseElementReference == nil {
		setReferencedBaseElementBaseElementReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(setReferencedBaseElementBaseElementReference, setReferencedBaseElement, hl)
		core.SetName(setReferencedBaseElementBaseElementReference, "SourceBaseElementRef", hl)
		core.SetUri(setReferencedBaseElementBaseElementReference, BaseElementReferenceSetReferencedBaseElementSourceBaseElementRefUri, hl)
	}
	// GetBaseElementTargetBaseElementReference
	setReferencedBaseElementTargetBaseElementReference := core.GetChildElementReferenceWithUri(setReferencedBaseElement, BaseElementReferenceSetReferencedBaseElementModifiedBaseElementReferenceRefUri, hl)
	if setReferencedBaseElementTargetBaseElementReference == nil {
		setReferencedBaseElementTargetBaseElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedBaseElementTargetBaseElementReference, setReferencedBaseElement, hl)
		core.SetName(setReferencedBaseElementTargetBaseElementReference, "ModifiedBaseElementReferenceRef", hl)
		core.SetUri(setReferencedBaseElementTargetBaseElementReference, BaseElementReferenceSetReferencedBaseElementModifiedBaseElementReferenceRefUri, hl)
	}
}

func baseElementReferenceFunctionsInit() {
	core.GetCore().AddFunction(BaseElementReferenceCreateUri, createBaseElementReference)
	core.GetCore().AddFunction(BaseElementReferenceGetBaseElementPointerUri, getBaseElementPointer)
	core.GetCore().AddFunction(BaseElementReferenceGetReferencedBaseElementUri, getReferencedBaseElement)
	core.GetCore().AddFunction(BaseElementReferenceSetReferencedBaseElementUri, setReferencedBaseElement)
}
