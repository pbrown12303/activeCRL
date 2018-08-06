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

var ElementPointerFunctionsUri string = CoreFunctionsPrefix + "ElementPointerFunctions"

var ElementPointerCreateAbstractElementPointerUri string = CoreFunctionsPrefix + "ElementPointer/CreateAbstractElementPointer"
var ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri = CoreFunctionsPrefix + "ElementPointer/CreateAbstractElementPointer/CreatedElementPointerRef"

var ElementPointerCreateRefinedElementPointerUri string = CoreFunctionsPrefix + "ElementPointer/CreateRefinedElementPointer"
var ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri = CoreFunctionsPrefix + "ElementPointer/CreateRefinedElementPointer/CreatedElementPointerRef"

var ElementPointerCreateOwningElementPointerUri string = CoreFunctionsPrefix + "ElementPointer/CreateOwningElementPointer"
var ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri = CoreFunctionsPrefix + "ElementPointer/CreateOwningElementPointer/CreatedElementPointerRef"

var ElementPointerCreateReferencedElementPointerUri string = CoreFunctionsPrefix + "ElementPointer/CreateReferencedElementPointer"
var ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri = CoreFunctionsPrefix + "ElementPointer/CreateReferencedElementPointer/CreatedElementPointerRef"

var ElementPointerGetElementUri string = CoreFunctionsPrefix + "ElementPointer/GetElement"
var ElementPointerGetElementSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElement/SourceElementPointerRef"
var ElementPointerGetElementIndicatedElementRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElement/IndicatedElementRef"

var ElementPointerGetElementIdUri string = CoreFunctionsPrefix + "ElementPointer/GetElementId"
var ElementPointerGetElementIdSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementId/SourceElementPointerRef"
var ElementPointerGetElementIdCreatedLiteralUri string = CoreFunctionsPrefix + "ElementPointer/GetElementId/CreatedLiteralRef"

var ElementPointerGetElementPointerRoleUri string = CoreFunctionsPrefix + "ElementPointer/GetElementPointerRole"
var ElementPointerGetElementPointerRoleSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementPointerRole/SourceElementPointerRef"
var ElementPointerGetElementPointerRoleCreatedLiteralRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementPointerRole/CreatedLiteralRef"

var ElementPointerGetElementVersionUri string = CoreFunctionsPrefix + "ElementPointer/GetElementVersion"
var ElementPointerGetElementVersionSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementVersion/SourceElementPointerRef"
var ElementPointerGetElementVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "ElementPointer/GetElementVersion/CreatedLiteralRef"

var ElementPointerSetElementUri string = CoreFunctionsPrefix + "ElementPointer/SetElement"
var ElementPointerSetElementElementRefUri string = CoreFunctionsPrefix + "ElementPointer/SetElement/ElementRef"
var ElementPointerSetElementModifiedElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointer/SetElement/ModifiedElementPointerRef"

func createAbstractElementPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createAbstractElementPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementPointerReference := core.GetChildElementPointerReferenceWithAncestorUri(element, ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri, hl)
	if createdElementPointerReference == nil {
		createdElementPointerReference = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(createdElementPointerReference, element, hl)
		core.SetLabel(createdElementPointerReference, "CreatedElementPointerReference", hl)
		rootCreatedElementReference := uOfD.GetElementPointerReferenceWithUri(ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementPointerReference, hl)
		refinement.SetRefinedElement(createdElementPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdElementPointer := createdElementPointerReference.GetReferencedElementPointer(hl)
	if createdElementPointer == nil {
		createdElementPointer = uOfD.NewAbstractElementPointer(hl)
		createdElementPointerReference.SetReferencedElementPointer(createdElementPointer, hl)
	}
}

func createRefinedElementPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createRefinedElementPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementPointerReference := core.GetChildElementPointerReferenceWithAncestorUri(element, ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri, hl)
	if createdElementPointerReference == nil {
		createdElementPointerReference = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(createdElementPointerReference, element, hl)
		core.SetLabel(createdElementPointerReference, "CreatedElementPointerReference", hl)
		rootCreatedElementReference := uOfD.GetElementPointerReferenceWithUri(ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementPointerReference, hl)
		refinement.SetRefinedElement(createdElementPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdElementPointer := createdElementPointerReference.GetReferencedElementPointer(hl)
	if createdElementPointer == nil {
		createdElementPointer = uOfD.NewRefinedElementPointer(hl)
		createdElementPointerReference.SetReferencedElementPointer(createdElementPointer, hl)
	}
}

func createOwningElementPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createOwningElementPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementPointerReference := core.GetChildElementPointerReferenceWithAncestorUri(element, ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri, hl)
	if createdElementPointerReference == nil {
		createdElementPointerReference = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(createdElementPointerReference, element, hl)
		core.SetLabel(createdElementPointerReference, "CreatedElementPointerReference", hl)
		rootCreatedElementReference := uOfD.GetElementPointerReferenceWithUri(ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementPointerReference, hl)
		refinement.SetRefinedElement(createdElementPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdElementPointer := createdElementPointerReference.GetReferencedElementPointer(hl)
	if createdElementPointer == nil {
		createdElementPointer = uOfD.NewOwningElementPointer(hl)
		createdElementPointerReference.SetReferencedElementPointer(createdElementPointer, hl)
	}
}

func createReferencedElementPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createReferencedElementPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementPointerReference := core.GetChildElementPointerReferenceWithAncestorUri(element, ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri, hl)
	if createdElementPointerReference == nil {
		createdElementPointerReference = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(createdElementPointerReference, element, hl)
		core.SetLabel(createdElementPointerReference, "CreatedElementPointerReference", hl)
		rootCreatedElementReference := uOfD.GetElementPointerReferenceWithUri(ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementPointerReference, hl)
		refinement.SetRefinedElement(createdElementPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdElementPointer := createdElementPointerReference.GetReferencedElementPointer(hl)
	if createdElementPointer == nil {
		createdElementPointer = uOfD.NewReferencedElementPointer(hl)
		createdElementPointerReference.SetReferencedElementPointer(createdElementPointer, hl)
	}
}

func getElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerGetElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerGetElementSourceElementPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetElement, the SourceReference was not found in the replicate")
		return
	}

	targetElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerGetElementIndicatedElementRefUri, hl)
	if targetElementReference == nil {
		log.Printf("In GetElement, the TargetElementPointerReference was not found in the replicate")
		return
	}

	targetElement := targetElementReference.GetReferencedElement(hl)
	sourceElementPointer := sourceReference.GetReferencedElementPointer(hl)
	var sourceElement core.Element
	if sourceElementPointer != nil {
		sourceElement = sourceElementPointer.GetElement(hl)
	}
	if sourceElement != targetElement {
		targetElementReference.SetReferencedElement(sourceElement, hl)
	}
}

func getElementId(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerGetElementIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerGetElementIdSourceElementPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetElementId, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerGetElementIdCreatedLiteralUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetElementId, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	sourceElementPointer := sourceReference.GetReferencedElementPointer(hl)
	if sourceElementPointer != nil {
		createdLiteral.SetLiteralValue(sourceElementPointer.GetElementId(hl), hl)
	}
}

func getElementPointerRole(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerGetElementPointerRoleUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerGetElementPointerRoleSourceElementPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetElementPointerRole, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerGetElementPointerRoleCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetElementPointerRole, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	sourceElementPointer := sourceReference.GetReferencedElementPointer(hl)
	if sourceElementPointer != nil {
		stringBaseElementRole := sourceElementPointer.GetElementPointerRole(hl).RoleToString()
		createdLiteral.SetLiteralValue(stringBaseElementRole, hl)
	}
}

func getElementVersion(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerGetElementVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerGetElementVersionSourceElementPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetElementVersion, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementPointerGetElementVersionCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetElementVersion, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	sourceElementPointer := sourceReference.GetReferencedElementPointer(hl)
	if sourceElementPointer != nil {
		stringBaseElementVersion := strconv.Itoa(sourceElementPointer.GetElementVersion(hl))
		createdLiteral.SetLiteralValue(stringBaseElementVersion, hl)
	}
}

func setElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerSetElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	elementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerSetElementElementRefUri, hl)
	if elementRef == nil {
		log.Printf("In SetElement, the ElementReference was not found in the replicate")
		return
	}

	targetElementPointerReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerSetElementModifiedElementPointerRefUri, hl)
	if targetElementPointerReference == nil {
		log.Printf("In SetElement, the TargetElementPointerReference was not found in the replicate")
		return
	}

	targetElementPointer := targetElementPointerReference.GetReferencedElementPointer(hl)
	baseElement := elementRef.GetReferencedElement(hl)
	if targetElementPointer != nil {
		targetElementPointer.SetElement(baseElement, hl)
	}
}

func BuildCoreElementPointerFunctions(coreFunctionsElement core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// ElementPointerFunctions
	elementPointerFunctions := uOfD.NewElement(hl, ElementPointerFunctionsUri)
	core.SetOwningElement(elementPointerFunctions, coreFunctionsElement, hl)
	core.SetLabel(elementPointerFunctions, "ElementPointerFunctions", hl)
	core.SetUri(elementPointerFunctions, ElementPointerFunctionsUri, hl)

	// CreateAbstractElementPointerElement
	createAbstractElementPointer := uOfD.NewElement(hl, ElementPointerCreateAbstractElementPointerUri)
	core.SetOwningElement(createAbstractElementPointer, elementPointerFunctions, hl)
	core.SetLabel(createAbstractElementPointer, "CreateAbstractElementPointerElementPointer", hl)
	core.SetUri(createAbstractElementPointer, ElementPointerCreateAbstractElementPointerUri, hl)
	// CreatedElementReference
	createdElementPointerReference0 := uOfD.NewElementPointerReference(hl, ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri)
	core.SetOwningElement(createdElementPointerReference0, createAbstractElementPointer, hl)
	core.SetLabel(createdElementPointerReference0, "CreateAbstractElementPointerdElementPointerRef", hl)
	core.SetUri(createdElementPointerReference0, ElementPointerCreateAbstractElementPointerCreatedElementPointerRefUri, hl)

	// CreateRefinedElementPointerElement
	createRefinedElementPointer := uOfD.NewElement(hl, ElementPointerCreateRefinedElementPointerUri)
	core.SetOwningElement(createRefinedElementPointer, elementPointerFunctions, hl)
	core.SetLabel(createRefinedElementPointer, "CreateRefinedElementPointerElementPointer", hl)
	core.SetUri(createRefinedElementPointer, ElementPointerCreateRefinedElementPointerUri, hl)
	// CreatedElementReference
	createdElementPointerReference1 := uOfD.NewElementPointerReference(hl, ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri)
	core.SetOwningElement(createdElementPointerReference1, createRefinedElementPointer, hl)
	core.SetLabel(createdElementPointerReference1, "CreateRefinedElementPointerdElementPointerRef", hl)
	core.SetUri(createdElementPointerReference1, ElementPointerCreateRefinedElementPointerCreatedElementPointerRefUri, hl)

	// CreateOwningElementPointerElement
	createOwningElementPointer := uOfD.NewElement(hl, ElementPointerCreateOwningElementPointerUri)
	core.SetOwningElement(createOwningElementPointer, elementPointerFunctions, hl)
	core.SetLabel(createOwningElementPointer, "CreateOwningElementPointerElementPointer", hl)
	core.SetUri(createOwningElementPointer, ElementPointerCreateOwningElementPointerUri, hl)
	// CreatedElementReference
	createdElementPointerReference2 := uOfD.NewElementPointerReference(hl, ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri)
	core.SetOwningElement(createdElementPointerReference2, createOwningElementPointer, hl)
	core.SetLabel(createdElementPointerReference2, "CreateOwningElementPointerdElementPointerRef", hl)
	core.SetUri(createdElementPointerReference2, ElementPointerCreateOwningElementPointerCreatedElementPointerRefUri, hl)

	// CreateReferencedElementPointerElement
	createReferencedElementPointer := uOfD.NewElement(hl, ElementPointerCreateReferencedElementPointerUri)
	core.SetOwningElement(createReferencedElementPointer, elementPointerFunctions, hl)
	core.SetLabel(createReferencedElementPointer, "CreateReferencedElementPointerElementPointer", hl)
	core.SetUri(createReferencedElementPointer, ElementPointerCreateReferencedElementPointerUri, hl)
	// CreatedElementReference
	createdElementPointerReference3 := uOfD.NewElementPointerReference(hl, ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri)
	core.SetOwningElement(createdElementPointerReference3, createReferencedElementPointer, hl)
	core.SetLabel(createdElementPointerReference3, "CreateReferencedElementPointerdElementPointerRef", hl)
	core.SetUri(createdElementPointerReference3, ElementPointerCreateReferencedElementPointerCreatedElementPointerRefUri, hl)

	// GetElement
	getElement := uOfD.NewElement(hl, ElementPointerGetElementUri)
	core.SetLabel(getElement, "GetElement", hl)
	core.SetOwningElement(getElement, elementPointerFunctions, hl)
	core.SetUri(getElement, ElementPointerGetElementUri, hl)
	// GetElement.SourceReference
	getElementSourceReference := uOfD.NewElementPointerReference(hl, ElementPointerGetElementSourceElementPointerRefUri)
	core.SetOwningElement(getElementSourceReference, getElement, hl)
	core.SetLabel(getElementSourceReference, "SourceElementPointerRef", hl)
	core.SetUri(getElementSourceReference, ElementPointerGetElementSourceElementPointerRefUri, hl)
	// GetElementTargetElementPointerReference
	getElementTargetReference := uOfD.NewElementReference(hl, ElementPointerGetElementIndicatedElementRefUri)
	core.SetOwningElement(getElementTargetReference, getElement, hl)
	core.SetLabel(getElementTargetReference, "IndicatedBaseElementRef", hl)
	core.SetUri(getElementTargetReference, ElementPointerGetElementIndicatedElementRefUri, hl)

	// GetElementId
	getElementId := uOfD.NewElement(hl, ElementPointerGetElementIdUri)
	core.SetLabel(getElementId, "GetElementId", hl)
	core.SetOwningElement(getElementId, elementPointerFunctions, hl)
	core.SetUri(getElementId, ElementPointerGetElementIdUri, hl)
	// GetElementId.SourceReference
	getElementIdSourceReference := uOfD.NewElementPointerReference(hl, ElementPointerGetElementIdSourceElementPointerRefUri)
	core.SetOwningElement(getElementIdSourceReference, getElementId, hl)
	core.SetLabel(getElementIdSourceReference, "SourceElementPointerRef", hl)
	core.SetUri(getElementIdSourceReference, ElementPointerGetElementIdSourceElementPointerRefUri, hl)
	// GetElementIdTargetLiteralReference
	getElementIdTargetReference := uOfD.NewLiteralReference(hl, ElementPointerGetElementIdCreatedLiteralUri)
	core.SetOwningElement(getElementIdTargetReference, getElementId, hl)
	core.SetLabel(getElementIdTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getElementIdTargetReference, ElementPointerGetElementIdCreatedLiteralUri, hl)

	// GetElementPointerRole
	getElementPointerRole := uOfD.NewElement(hl, ElementPointerGetElementPointerRoleUri)
	core.SetLabel(getElementPointerRole, "GetElementPointerRole", hl)
	core.SetOwningElement(getElementPointerRole, elementPointerFunctions, hl)
	core.SetUri(getElementPointerRole, ElementPointerGetElementPointerRoleUri, hl)
	// GetElementPointerRole.SourceReference
	getElementPointerRoleSourceReference := uOfD.NewElementPointerReference(hl, ElementPointerGetElementPointerRoleSourceElementPointerRefUri)
	core.SetOwningElement(getElementPointerRoleSourceReference, getElementPointerRole, hl)
	core.SetLabel(getElementPointerRoleSourceReference, "SourceElementPointerRef", hl)
	core.SetUri(getElementPointerRoleSourceReference, ElementPointerGetElementPointerRoleSourceElementPointerRefUri, hl)
	// GetElementPointerRoleTargetLiteralReference
	getElementPointerRoleTargetReference := uOfD.NewLiteralReference(hl, ElementPointerGetElementPointerRoleCreatedLiteralRefUri)
	core.SetOwningElement(getElementPointerRoleTargetReference, getElementPointerRole, hl)
	core.SetLabel(getElementPointerRoleTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getElementPointerRoleTargetReference, ElementPointerGetElementPointerRoleCreatedLiteralRefUri, hl)

	// GetElementVersion
	getElementVersion := uOfD.NewElement(hl, ElementPointerGetElementVersionUri)
	core.SetLabel(getElementVersion, "GetElementVersion", hl)
	core.SetOwningElement(getElementVersion, elementPointerFunctions, hl)
	core.SetUri(getElementVersion, ElementPointerGetElementVersionUri, hl)
	// GetElementVersion.SourceReference
	getElementVersionSourceReference := uOfD.NewElementPointerReference(hl, ElementPointerGetElementVersionSourceElementPointerRefUri)
	core.SetOwningElement(getElementVersionSourceReference, getElementVersion, hl)
	core.SetLabel(getElementVersionSourceReference, "SourceElementPointerRef", hl)
	core.SetUri(getElementVersionSourceReference, ElementPointerGetElementVersionSourceElementPointerRefUri, hl)
	// GetElementVersionTargetLiteralReference
	getElementVersionTargetReference := uOfD.NewLiteralReference(hl, ElementPointerGetElementVersionCreatedLiteralRefUri)
	core.SetOwningElement(getElementVersionTargetReference, getElementVersion, hl)
	core.SetLabel(getElementVersionTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getElementVersionTargetReference, ElementPointerGetElementVersionCreatedLiteralRefUri, hl)

	// SetElement
	setElement := uOfD.NewElement(hl, ElementPointerSetElementUri)
	core.SetLabel(setElement, "SetElement", hl)
	core.SetOwningElement(setElement, elementPointerFunctions, hl)
	core.SetUri(setElement, ElementPointerSetElementUri, hl)
	// SetElement.ElementReference
	setElementElementReference := uOfD.NewElementReference(hl, ElementPointerSetElementElementRefUri)
	core.SetLabel(setElementElementReference, "BaseElementRef", hl)
	core.SetOwningElement(setElementElementReference, setElement, hl)
	core.SetUri(setElementElementReference, ElementPointerSetElementElementRefUri, hl)
	// SetElement.TargetElementPointerReference
	setElementTargetElementPointerReference := uOfD.NewElementPointerReference(hl, ElementPointerSetElementModifiedElementPointerRefUri)
	core.SetLabel(setElementTargetElementPointerReference, "ModifiedElementPointerRef", hl)
	core.SetOwningElement(setElementTargetElementPointerReference, setElement, hl)
	core.SetUri(setElementTargetElementPointerReference, ElementPointerSetElementModifiedElementPointerRefUri, hl)
}

func elementPointerFunctionsInit() {
	core.GetCore().AddFunction(ElementPointerCreateAbstractElementPointerUri, createAbstractElementPointer)
	core.GetCore().AddFunction(ElementPointerCreateRefinedElementPointerUri, createRefinedElementPointer)
	core.GetCore().AddFunction(ElementPointerCreateOwningElementPointerUri, createOwningElementPointer)
	core.GetCore().AddFunction(ElementPointerCreateReferencedElementPointerUri, createReferencedElementPointer)
	core.GetCore().AddFunction(ElementPointerGetElementUri, getElement)
	core.GetCore().AddFunction(ElementPointerGetElementIdUri, getElementId)
	core.GetCore().AddFunction(ElementPointerGetElementPointerRoleUri, getElementPointerRole)
	core.GetCore().AddFunction(ElementPointerGetElementVersionUri, getElementVersion)
	core.GetCore().AddFunction(ElementPointerSetElementUri, setElement)
}
