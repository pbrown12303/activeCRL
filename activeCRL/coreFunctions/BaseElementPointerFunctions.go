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

var BaseElementPointerFunctionsUri string = CoreFunctionsPrefix + "BaseElementPointerFunctions"

var BaseElementPointerCreateUri string = CoreFunctionsPrefix + "BaseElementPointer/Create"
var BaseElementPointerCreateCreatedBaseElementPointerRefUri = CoreFunctionsPrefix + "BaseElementPointer/Create/CreatedBaseElementPointerRef"

var BaseElementPointerGetBaseElementUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement"
var BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement/SourceBaseElementPointerRef"
var BaseElementPointerGetBaseElementIndicatedBaseElementRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement/IndicatedBaseElementRef"

var BaseElementPointerGetBaseElementIdUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId"
var BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId/SourceBaseElementPointerRef"
var BaseElementPointerGetBaseElementIdCreatedLiteralUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId/CreatedLiteralRef"

var BaseElementPointerGetBaseElementVersionUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion"
var BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion/SourceBaseElementPointerRef"
var BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion/CreatedLiteralRef"

var BaseElementPointerSetBaseElementUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement"
var BaseElementPointerSetBaseElementBaseElementRefUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement/BaseElementRef"
var BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement/ModifiedBaseElementPointerRef"

func createBaseElementPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createBaseElementPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdBaseElementPointerReference := core.GetChildBaseElementReferenceWithAncestorUri(element, BaseElementPointerCreateCreatedBaseElementPointerRefUri, hl)
	if createdBaseElementPointerReference == nil {
		createdBaseElementPointerReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(createdBaseElementPointerReference, element, hl)
		core.SetLabel(createdBaseElementPointerReference, "CreatedBaseElementPointerReference", hl)
		rootCreatedElementReference := uOfD.GetBaseElementReferenceWithUri(BaseElementPointerCreateCreatedBaseElementPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdBaseElementPointerReference, hl)
		refinement.SetRefinedElement(createdBaseElementPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdBaseElementPointer := createdBaseElementPointerReference.GetReferencedBaseElement(hl)
	if createdBaseElementPointer == nil {
		createdBaseElementPointer = uOfD.NewBaseElementPointer(hl)
		createdBaseElementPointerReference.SetReferencedBaseElement(createdBaseElementPointer, hl)
	}
}

func getBaseElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetBaseElement, the SourceReference was not found in the replicate")
		return
	}

	targetElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIndicatedBaseElementRefUri, hl)
	if targetElementReference == nil {
		log.Printf("In GetBaseElement, the TargetBaseElementReference was not found in the replicate")
		return
	}

	targetBaseElement := targetElementReference.GetReferencedBaseElement(hl)
	untypedBaseElement := sourceReference.GetReferencedBaseElement(hl)
	var sourceBaseElement core.BaseElement
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.BaseElementPointer:
			sourceBaseElementPointer := untypedBaseElement.(core.BaseElementPointer)
			sourceBaseElement = sourceBaseElementPointer.GetBaseElement(hl)
		default:
			log.Printf("In GetBaseElement, the SourceBaseElement is not a BaseElementPointer")
		}
	}
	if sourceBaseElement != targetBaseElement {
		targetElementReference.SetReferencedBaseElement(sourceBaseElement, hl)
	}
}

func getBaseElementId(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetBaseElementId, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIdCreatedLiteralUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetBaseElementId, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedBaseElement := sourceReference.GetReferencedBaseElement(hl)
	var sourceBaseElementPointer core.BaseElementPointer
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.BaseElementPointer:
			sourceBaseElementPointer = untypedBaseElement.(core.BaseElementPointer)
		default:
			log.Printf("In GetBaseElementId, the SourceBaseElement is not a BaseElementPointer")
		}
	}
	if sourceBaseElementPointer != nil {
		createdLiteral.SetLiteralValue(sourceBaseElementPointer.GetBaseElementId(hl).String(), hl)
	}
}

func getBaseElementVersion(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetBaseElementVersion, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetBaseElementVersion, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedBaseElement := sourceReference.GetReferencedBaseElement(hl)
	var sourceBaseElementPointer core.BaseElementPointer
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.BaseElementPointer:
			sourceBaseElementPointer = untypedBaseElement.(core.BaseElementPointer)
		default:
			log.Printf("In GetBaseElementVersion, the SourceBaseElement is not a BaseElementPointer")
		}
	}
	if sourceBaseElementPointer != nil {
		stringBaseElementVersion := strconv.Itoa(sourceBaseElementPointer.GetBaseElementVersion(hl))
		createdLiteral.SetLiteralValue(stringBaseElementVersion, hl)
	}
}

func setBaseElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementPointerSetBaseElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	baseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerSetBaseElementBaseElementRefUri, hl)
	if baseElementReference == nil {
		log.Printf("In SetBaseElement, the BaseElementReference was not found in the replicate")
		return
	}

	targetBaseElementPointerReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri, hl)
	if targetBaseElementPointerReference == nil {
		log.Printf("In SetBaseElement, the TargetBaseElementPointerReference was not found in the replicate")
		return
	}

	untypedBaseElement := targetBaseElementPointerReference.GetReferencedBaseElement(hl)
	baseElement := baseElementReference.GetReferencedBaseElement(hl)
	var targetBaseElementPointer core.BaseElementPointer
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.BaseElementPointer:
			targetBaseElementPointer = untypedBaseElement.(core.BaseElementPointer)
			targetBaseElementPointer.SetBaseElement(baseElement, hl)
		default:
			log.Printf("In SetBaseElement, the TargetBaseElementPointerReference does not point to a BaseElementPointer")
		}
	}
}

func BuildCoreBaseElementPointerFunctions(coreFunctionsElement core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// BaseElementPointerFunctions
	baseElementPointerFunctions := uOfD.NewElement(hl, BaseElementPointerFunctionsUri)
	core.SetOwningElement(baseElementPointerFunctions, coreFunctionsElement, hl)
	core.SetLabel(baseElementPointerFunctions, "BaseElementPointerFunctions", hl)
	core.SetUri(baseElementPointerFunctions, BaseElementPointerFunctionsUri, hl)

	// CreateElement
	createBaseElementPointer := uOfD.NewElement(hl, BaseElementPointerCreateUri)
	core.SetOwningElement(createBaseElementPointer, baseElementPointerFunctions, hl)
	core.SetLabel(createBaseElementPointer, "CreateBaseElementPointer", hl)
	core.SetUri(createBaseElementPointer, BaseElementPointerCreateUri, hl)
	// CreatedElementReference
	createdBaseElementPointerReference := uOfD.NewBaseElementReference(hl, BaseElementPointerCreateCreatedBaseElementPointerRefUri)
	core.SetOwningElement(createdBaseElementPointerReference, createBaseElementPointer, hl)
	core.SetLabel(createdBaseElementPointerReference, "CreatedBaseElementPointerRef", hl)
	core.SetUri(createdBaseElementPointerReference, BaseElementPointerCreateCreatedBaseElementPointerRefUri, hl)

	// GetBaseElement
	getBaseElement := uOfD.NewElement(hl, BaseElementPointerGetBaseElementUri)
	core.SetLabel(getBaseElement, "GetBaseElement", hl)
	core.SetOwningElement(getBaseElement, baseElementPointerFunctions, hl)
	core.SetUri(getBaseElement, BaseElementPointerGetBaseElementUri, hl)
	// GetBaseElement.SourceReference
	getBaseElementSourceReference := uOfD.NewBaseElementReference(hl, BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri)
	core.SetOwningElement(getBaseElementSourceReference, getBaseElement, hl)
	core.SetLabel(getBaseElementSourceReference, "SourceBaseElementPointerRef", hl)
	core.SetUri(getBaseElementSourceReference, BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri, hl)
	// GetBaseElementTargetBaseElementReference
	getBaseElementTargetReference := uOfD.NewBaseElementReference(hl, BaseElementPointerGetBaseElementIndicatedBaseElementRefUri)
	core.SetOwningElement(getBaseElementTargetReference, getBaseElement, hl)
	core.SetLabel(getBaseElementTargetReference, "IndicatedBaseElementRef", hl)
	core.SetUri(getBaseElementTargetReference, BaseElementPointerGetBaseElementIndicatedBaseElementRefUri, hl)

	// GetBaseElementId
	getBaseElementId := uOfD.NewElement(hl, BaseElementPointerGetBaseElementIdUri)
	core.SetLabel(getBaseElementId, "GetBaseElementId", hl)
	core.SetOwningElement(getBaseElementId, baseElementPointerFunctions, hl)
	core.SetUri(getBaseElementId, BaseElementPointerGetBaseElementIdUri, hl)
	// GetBaseElementId.SourceReference
	getBaseElementIdSourceReference := uOfD.NewBaseElementReference(hl, BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri)
	core.SetOwningElement(getBaseElementIdSourceReference, getBaseElementId, hl)
	core.SetLabel(getBaseElementIdSourceReference, "SourceBaseElementPointerRef", hl)
	core.SetUri(getBaseElementIdSourceReference, BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri, hl)
	// GetBaseElementIdTargetLiteralReference
	getBaseElementIdTargetReference := uOfD.NewLiteralReference(hl, BaseElementPointerGetBaseElementIdCreatedLiteralUri)
	core.SetOwningElement(getBaseElementIdTargetReference, getBaseElementId, hl)
	core.SetLabel(getBaseElementIdTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getBaseElementIdTargetReference, BaseElementPointerGetBaseElementIdCreatedLiteralUri, hl)

	// GetBaseElementVersion
	getBaseElementVersion := uOfD.NewElement(hl, BaseElementPointerGetBaseElementVersionUri)
	core.SetLabel(getBaseElementVersion, "GetBaseElementVersion", hl)
	core.SetOwningElement(getBaseElementVersion, baseElementPointerFunctions, hl)
	core.SetUri(getBaseElementVersion, BaseElementPointerGetBaseElementVersionUri, hl)
	// GetBaseElementVersion.SourceReference
	getBaseElementVersionSourceReference := uOfD.NewBaseElementReference(hl, BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri)
	core.SetOwningElement(getBaseElementVersionSourceReference, getBaseElementVersion, hl)
	core.SetLabel(getBaseElementVersionSourceReference, "SourceBaseElementPointerRef", hl)
	core.SetUri(getBaseElementVersionSourceReference, BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri, hl)
	// GetBaseElementVersionTargetLiteralReference
	getBaseElementVersionTargetReference := uOfD.NewLiteralReference(hl, BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri)
	core.SetOwningElement(getBaseElementVersionTargetReference, getBaseElementVersion, hl)
	core.SetLabel(getBaseElementVersionTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getBaseElementVersionTargetReference, BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri, hl)

	// SetBaseElement
	setBaseElement := uOfD.NewElement(hl, BaseElementPointerSetBaseElementUri)
	core.SetLabel(setBaseElement, "SetBaseElement", hl)
	core.SetOwningElement(setBaseElement, baseElementPointerFunctions, hl)
	core.SetUri(setBaseElement, BaseElementPointerSetBaseElementUri, hl)
	// SetBaseElement.BaseElementReference
	setBaseElementBaseElementReference := uOfD.NewBaseElementReference(hl, BaseElementPointerSetBaseElementBaseElementRefUri)
	core.SetLabel(setBaseElementBaseElementReference, "BaseElementRef", hl)
	core.SetOwningElement(setBaseElementBaseElementReference, setBaseElement, hl)
	core.SetUri(setBaseElementBaseElementReference, BaseElementPointerSetBaseElementBaseElementRefUri, hl)
	setBaseElementTargetBaseElementPointerReference := uOfD.NewBaseElementReference(hl, BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri)
	core.SetLabel(setBaseElementTargetBaseElementPointerReference, "ModifiedBaseElementPointerRef", hl)
	core.SetOwningElement(setBaseElementTargetBaseElementPointerReference, setBaseElement, hl)
	core.SetUri(setBaseElementTargetBaseElementPointerReference, BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri, hl)
}

func baseElementPointerFunctionsInit() {
	core.GetCore().AddFunction(BaseElementPointerCreateUri, createBaseElementPointer)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementUri, getBaseElement)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementIdUri, getBaseElementId)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementVersionUri, getBaseElementVersion)
	core.GetCore().AddFunction(BaseElementPointerSetBaseElementUri, setBaseElement)
}
