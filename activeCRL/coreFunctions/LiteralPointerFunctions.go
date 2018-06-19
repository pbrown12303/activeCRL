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

var LiteralPointerFunctionsUri string = CoreFunctionsPrefix + "LiteralPointerFunctions"

var LiteralPointerCreateLabelLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointer/CreateLabelLiteralPointer"
var LiteralPointerCreateLabelLiteralPointerCreatedLiteralPointerRefUri = CoreFunctionsPrefix + "LiteralPointer/CreateLabelLiteralPointer/CreatedLiteralPointerRef"

var LiteralPointerCreateDefinitionLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointer/CreateDefinitionLiteralPointer"
var LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri = CoreFunctionsPrefix + "LiteralPointer/CreateDefinitionLiteralPointer/CreatedLiteralPointerRef"

var LiteralPointerCreateUriLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointer/CreateUriLiteralPointer"
var LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri = CoreFunctionsPrefix + "LiteralPointer/CreateUriLiteralPointer/CreatedLiteralPointerRef"

var LiteralPointerCreateValueLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointer/CreateValueLiteralPointer"
var LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri = CoreFunctionsPrefix + "LiteralPointer/CreateValueLiteralPointer/CreatedLiteralPointerRef"

var LiteralPointerGetLiteralUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteral"
var LiteralPointerGetLiteralSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteral/SourceLiteralPointerRef"
var LiteralPointerGetLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteral/IndicatedLiteralRef"

var LiteralPointerGetLiteralIdUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralId"
var LiteralPointerGetLiteralIdSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralId/SourceLiteralPointerRef"
var LiteralPointerGetLiteralIdCreatedLiteralUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralId/CreatedLiteralRef"

var LiteralPointerGetLiteralPointerRoleUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralPointerRole"
var LiteralPointerGetLiteralPointerRoleSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralPointerRole/SourceLiteralPointerRef"
var LiteralPointerGetLiteralPointerRoleCreatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralPointerRole/CreatedLiteralRef"

var LiteralPointerGetLiteralVersionUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralVersion"
var LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralVersion/SourceLiteralPointerRef"
var LiteralPointerGetLiteralVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralVersion/CreatedLiteralRef"

var LiteralPointerSetLiteralUri string = CoreFunctionsPrefix + "LiteralPointer/SetLiteral"
var LiteralPointerSetLiteralLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointer/SetLiteral/LiteralRef"
var LiteralPointerSetLiteralModifiedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/SetLiteral/ModifiedLiteralPointerRef"

func createLabelLiteralPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createLabelLiteralPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(element, LiteralPointerCreateLabelLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, element, hl)
		core.SetLabel(createdLiteralPointerReference, "CreatedLiteralPointerReference", hl)
		rootCreatedLiteralReference := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateLabelLiteralPointerCreatedLiteralPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralReference, hl)
	}
	createdLiteralPointer := createdLiteralPointerReference.GetReferencedLiteralPointer(hl)
	if createdLiteralPointer == nil {
		createdLiteralPointer = uOfD.NewLabelLiteralPointer(hl)
		createdLiteralPointerReference.SetReferencedLiteralPointer(createdLiteralPointer, hl)
	}
}

func createDefinitionLiteralPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createDefinitionLiteralPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(element, LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, element, hl)
		core.SetLabel(createdLiteralPointerReference, "CreatedLiteralPointerReference", hl)
		rootCreatedLiteralPointerReference := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralPointerReference, hl)
	}
	createdLiteralPointer := createdLiteralPointerReference.GetReferencedLiteralPointer(hl)
	if createdLiteralPointer == nil {
		createdLiteralPointer = uOfD.NewDefinitionLiteralPointer(hl)
		createdLiteralPointerReference.SetReferencedLiteralPointer(createdLiteralPointer, hl)
	}
}

func createUriLiteralPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createUriLiteralPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(element, LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, element, hl)
		core.SetLabel(createdLiteralPointerReference, "CreatedLiteralPointerReference", hl)
		rootCreatedLiteralPointerReference := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralPointerReference, hl)
	}
	createdLiteralPointer := createdLiteralPointerReference.GetReferencedLiteralPointer(hl)
	if createdLiteralPointer == nil {
		createdLiteralPointer = uOfD.NewUriLiteralPointer(hl)
		createdLiteralPointerReference.SetReferencedLiteralPointer(createdLiteralPointer, hl)
	}
}

func createValueLiteralPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createValueLiteralPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(element, LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, element, hl)
		core.SetLabel(createdLiteralPointerReference, "CreatedLiteralPointerReference", hl)
		rootCreatedLiteralPointerReference := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralPointerReference, hl)
	}
	createdLiteralPointer := createdLiteralPointerReference.GetReferencedLiteralPointer(hl)
	if createdLiteralPointer == nil {
		createdLiteralPointer = uOfD.NewValueLiteralPointer(hl)
		createdLiteralPointerReference.SetReferencedLiteralPointer(createdLiteralPointer, hl)
	}
}

func getLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerGetLiteralUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteral, the SourceReference was not found in the replicate")
		return
	}

	indicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		log.Printf("In GetLiteral, the TargetLiteralPointerReference was not found in the replicate")
		return
	}

	targetLiteral := indicatedLiteralRef.GetReferencedLiteral(hl)
	sourceLiteralPointer := sourceReference.GetReferencedLiteralPointer(hl)
	var sourceLiteral core.Literal
	if sourceLiteralPointer != nil {
		sourceLiteral = sourceLiteralPointer.GetLiteral(hl)
	}
	if sourceLiteral != targetLiteral {
		indicatedLiteralRef.SetReferencedLiteral(sourceLiteral, hl)
	}
}

func getLiteralId(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerGetLiteralIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIdSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteralId, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIdCreatedLiteralUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetLiteralId, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	sourceLiteralPointer := sourceReference.GetReferencedLiteralPointer(hl)
	if sourceLiteralPointer != nil {
		createdLiteral.SetLiteralValue(sourceLiteralPointer.GetLiteralId(hl).String(), hl)
	}
}

func getLiteralPointerRole(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerGetLiteralPointerRoleUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralPointerRoleSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteralPointerRole, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralPointerRoleCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetLiteralPointerRole, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedBaseElement := sourceReference.GetReferencedLiteralPointer(hl)
	var sourceLiteralPointer core.LiteralPointer
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.LiteralPointer:
			sourceLiteralPointer = untypedBaseElement.(core.LiteralPointer)
		default:
			log.Printf("In GetLiteralPointerRole, the SourceBaseElement is not a LiteralPointer")
		}
	}
	if sourceLiteralPointer != nil {
		stringBaseElementRole := sourceLiteralPointer.GetLiteralPointerRole(hl).RoleToString()
		createdLiteral.SetLiteralValue(stringBaseElementRole, hl)
	}
}

func getLiteralVersion(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerGetLiteralVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteralVersion, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralVersionCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetLiteralVersion, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	sourceLiteralPointer := sourceReference.GetReferencedLiteralPointer(hl)
	if sourceLiteralPointer != nil {
		stringBaseElementVersion := strconv.Itoa(sourceLiteralPointer.GetLiteralVersion(hl))
		createdLiteral.SetLiteralValue(stringBaseElementVersion, hl)
	}
}

func setLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerSetLiteralUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	literalRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerSetLiteralLiteralRefUri, hl)
	if literalRef == nil {
		log.Printf("In SetLiteral, the LiteralReference was not found in the replicate")
		return
	}

	modifiedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerSetLiteralModifiedLiteralPointerRefUri, hl)
	if modifiedLiteralPointerRef == nil {
		log.Printf("In SetLiteral, the TargetLiteralPointerReference was not found in the replicate")
		return
	}

	modifiedLiteralPointer := modifiedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	sourceLiteral := literalRef.GetReferencedLiteral(hl)
	if modifiedLiteralPointer != nil {
		modifiedLiteralPointer.SetLiteral(sourceLiteral, hl)
	}
}

func BuildCoreLiteralPointerFunctions(coreFunctionsElement core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralPointerFunctions
	literalPointerFunctions := uOfD.NewElement(hl, LiteralPointerFunctionsUri)
	core.SetOwningElement(literalPointerFunctions, coreFunctionsElement, hl)
	core.SetLabel(literalPointerFunctions, "LiteralPointerFunctions", hl)
	core.SetUri(literalPointerFunctions, LiteralPointerFunctionsUri, hl)

	// CreateLabelLiteralPointerElement
	createLabelLiteralPointer := uOfD.NewElement(hl, LiteralPointerCreateLabelLiteralPointerUri)
	core.SetOwningElement(createLabelLiteralPointer, literalPointerFunctions, hl)
	core.SetLabel(createLabelLiteralPointer, "CreateLabelLiteralPointerLiteralPointer", hl)
	core.SetUri(createLabelLiteralPointer, LiteralPointerCreateLabelLiteralPointerUri, hl)
	// CreatedLiteralReference
	createdLiteralPointerReference0 := uOfD.NewLiteralPointerReference(hl, LiteralPointerCreateLabelLiteralPointerCreatedLiteralPointerRefUri)
	core.SetOwningElement(createdLiteralPointerReference0, createLabelLiteralPointer, hl)
	core.SetLabel(createdLiteralPointerReference0, "CreateLabelLiteralPointerdLiteralPointerRef", hl)
	core.SetUri(createdLiteralPointerReference0, LiteralPointerCreateLabelLiteralPointerCreatedLiteralPointerRefUri, hl)

	// CreateDefinitionLiteralPointerElement
	createDefinitionLiteralPointer := uOfD.NewElement(hl, LiteralPointerCreateDefinitionLiteralPointerUri)
	core.SetOwningElement(createDefinitionLiteralPointer, literalPointerFunctions, hl)
	core.SetLabel(createDefinitionLiteralPointer, "CreateDefinitionLiteralPointerLiteralPointer", hl)
	core.SetUri(createDefinitionLiteralPointer, LiteralPointerCreateDefinitionLiteralPointerUri, hl)
	// CreatedLiteralReference
	createdLiteralPointerReference1 := uOfD.NewLiteralPointerReference(hl, LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri)
	core.SetOwningElement(createdLiteralPointerReference1, createDefinitionLiteralPointer, hl)
	core.SetLabel(createdLiteralPointerReference1, "CreateDefinitionLiteralPointerdLiteralPointerRef", hl)
	core.SetUri(createdLiteralPointerReference1, LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri, hl)

	// CreateUriLiteralPointerElement
	createUriLiteralPointer := uOfD.NewElement(hl, LiteralPointerCreateUriLiteralPointerUri)
	core.SetOwningElement(createUriLiteralPointer, literalPointerFunctions, hl)
	core.SetLabel(createUriLiteralPointer, "CreateUriLiteralPointerLiteralPointer", hl)
	core.SetUri(createUriLiteralPointer, LiteralPointerCreateUriLiteralPointerUri, hl)
	// CreatedLiteralReference
	createdLiteralPointerReference2 := uOfD.NewLiteralPointerReference(hl, LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri)
	core.SetOwningElement(createdLiteralPointerReference2, createUriLiteralPointer, hl)
	core.SetLabel(createdLiteralPointerReference2, "CreateUriLiteralPointerdLiteralPointerRef", hl)
	core.SetUri(createdLiteralPointerReference2, LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri, hl)

	// CreateValueLiteralPointerElement
	createValueLiteralPointer := uOfD.NewElement(hl, LiteralPointerCreateValueLiteralPointerUri)
	core.SetOwningElement(createValueLiteralPointer, literalPointerFunctions, hl)
	core.SetLabel(createValueLiteralPointer, "CreateValueLiteralPointerLiteralPointer", hl)
	core.SetUri(createValueLiteralPointer, LiteralPointerCreateValueLiteralPointerUri, hl)
	// CreatedLiteralReference
	createdLiteralPointerReference3 := uOfD.NewLiteralPointerReference(hl, LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri)
	core.SetOwningElement(createdLiteralPointerReference3, createValueLiteralPointer, hl)
	core.SetLabel(createdLiteralPointerReference3, "CreateValueLiteralPointerdLiteralPointerRef", hl)
	core.SetUri(createdLiteralPointerReference3, LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri, hl)

	// GetLiteral
	getLiteral := uOfD.NewElement(hl, LiteralPointerGetLiteralUri)
	core.SetLabel(getLiteral, "GetLiteral", hl)
	core.SetOwningElement(getLiteral, literalPointerFunctions, hl)
	core.SetUri(getLiteral, LiteralPointerGetLiteralUri, hl)
	// GetLiteral.SourceReference
	getLiteralSourceReference := uOfD.NewLiteralPointerReference(hl, LiteralPointerGetLiteralSourceLiteralPointerRefUri)
	core.SetOwningElement(getLiteralSourceReference, getLiteral, hl)
	core.SetLabel(getLiteralSourceReference, "SourceLiteralPointerRef", hl)
	core.SetUri(getLiteralSourceReference, LiteralPointerGetLiteralSourceLiteralPointerRefUri, hl)
	// GetLiteralTargetLiteralPointerReference
	getLiteralTargetReference := uOfD.NewLiteralReference(hl, LiteralPointerGetLiteralIndicatedLiteralRefUri)
	core.SetOwningElement(getLiteralTargetReference, getLiteral, hl)
	core.SetLabel(getLiteralTargetReference, "IndicatedBaseElementRef", hl)
	core.SetUri(getLiteralTargetReference, LiteralPointerGetLiteralIndicatedLiteralRefUri, hl)

	// GetLiteralId
	getLiteralId := uOfD.NewElement(hl, LiteralPointerGetLiteralIdUri)
	core.SetLabel(getLiteralId, "GetLiteralId", hl)
	core.SetOwningElement(getLiteralId, literalPointerFunctions, hl)
	core.SetUri(getLiteralId, LiteralPointerGetLiteralIdUri, hl)
	// GetLiteralId.SourceReference
	getLiteralIdSourceReference := uOfD.NewLiteralPointerReference(hl, LiteralPointerGetLiteralIdSourceLiteralPointerRefUri)
	core.SetOwningElement(getLiteralIdSourceReference, getLiteralId, hl)
	core.SetLabel(getLiteralIdSourceReference, "SourceLiteralPointerRef", hl)
	core.SetUri(getLiteralIdSourceReference, LiteralPointerGetLiteralIdSourceLiteralPointerRefUri, hl)
	// GetLiteralIdTargetLiteralReference
	getLiteralIdTargetReference := uOfD.NewLiteralReference(hl, LiteralPointerGetLiteralIdCreatedLiteralUri)
	core.SetOwningElement(getLiteralIdTargetReference, getLiteralId, hl)
	core.SetLabel(getLiteralIdTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getLiteralIdTargetReference, LiteralPointerGetLiteralIdCreatedLiteralUri, hl)

	// GetLiteralPointerRole
	getLiteralPointerRole := uOfD.NewElement(hl, LiteralPointerGetLiteralPointerRoleUri)
	core.SetLabel(getLiteralPointerRole, "GetLiteralPointerRole", hl)
	core.SetOwningElement(getLiteralPointerRole, literalPointerFunctions, hl)
	core.SetUri(getLiteralPointerRole, LiteralPointerGetLiteralPointerRoleUri, hl)
	// GetLiteralPointerRole.SourceReference
	getLiteralPointerRoleSourceReference := uOfD.NewLiteralPointerReference(hl, LiteralPointerGetLiteralPointerRoleSourceLiteralPointerRefUri)
	core.SetOwningElement(getLiteralPointerRoleSourceReference, getLiteralPointerRole, hl)
	core.SetLabel(getLiteralPointerRoleSourceReference, "SourceLiteralPointerRef", hl)
	core.SetUri(getLiteralPointerRoleSourceReference, LiteralPointerGetLiteralPointerRoleSourceLiteralPointerRefUri, hl)
	// GetLiteralPointerRoleTargetLiteralReference
	getLiteralPointerRoleTargetReference := uOfD.NewLiteralReference(hl, LiteralPointerGetLiteralPointerRoleCreatedLiteralRefUri)
	core.SetOwningElement(getLiteralPointerRoleTargetReference, getLiteralPointerRole, hl)
	core.SetLabel(getLiteralPointerRoleTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getLiteralPointerRoleTargetReference, LiteralPointerGetLiteralPointerRoleCreatedLiteralRefUri, hl)

	// GetLiteralVersion
	getLiteralVersion := uOfD.NewElement(hl, LiteralPointerGetLiteralVersionUri)
	core.SetLabel(getLiteralVersion, "GetLiteralVersion", hl)
	core.SetOwningElement(getLiteralVersion, literalPointerFunctions, hl)
	core.SetUri(getLiteralVersion, LiteralPointerGetLiteralVersionUri, hl)
	// GetLiteralVersion.SourceReference
	getLiteralVersionSourceReference := uOfD.NewLiteralPointerReference(hl, LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri)
	core.SetOwningElement(getLiteralVersionSourceReference, getLiteralVersion, hl)
	core.SetLabel(getLiteralVersionSourceReference, "SourceLiteralPointerRef", hl)
	core.SetUri(getLiteralVersionSourceReference, LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri, hl)
	// GetLiteralVersionTargetLiteralReference
	getLiteralVersionTargetReference := uOfD.NewLiteralReference(hl, LiteralPointerGetLiteralVersionCreatedLiteralRefUri)
	core.SetOwningElement(getLiteralVersionTargetReference, getLiteralVersion, hl)
	core.SetLabel(getLiteralVersionTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getLiteralVersionTargetReference, LiteralPointerGetLiteralVersionCreatedLiteralRefUri, hl)

	// SetLiteral
	setLiteral := uOfD.NewElement(hl, LiteralPointerSetLiteralUri)
	core.SetLabel(setLiteral, "SetLiteral", hl)
	core.SetOwningElement(setLiteral, literalPointerFunctions, hl)
	core.SetUri(setLiteral, LiteralPointerSetLiteralUri, hl)
	// SetLiteral.LiteralReference
	setLiteralLiteralReference := uOfD.NewLiteralReference(hl, LiteralPointerSetLiteralLiteralRefUri)
	core.SetLabel(setLiteralLiteralReference, "BaseElementRef", hl)
	core.SetOwningElement(setLiteralLiteralReference, setLiteral, hl)
	core.SetUri(setLiteralLiteralReference, LiteralPointerSetLiteralLiteralRefUri, hl)
	// SetLiteral.ModifiedLiteralPointerReference
	setLiteralTargetLiteralPointerReference := uOfD.NewLiteralPointerReference(hl, LiteralPointerSetLiteralModifiedLiteralPointerRefUri)
	core.SetLabel(setLiteralTargetLiteralPointerReference, "ModifiedLiteralPointerRef", hl)
	core.SetOwningElement(setLiteralTargetLiteralPointerReference, setLiteral, hl)
	core.SetUri(setLiteralTargetLiteralPointerReference, LiteralPointerSetLiteralModifiedLiteralPointerRefUri, hl)
}

func literalPointerFunctionsInit() {
	core.GetCore().AddFunction(LiteralPointerCreateLabelLiteralPointerUri, createLabelLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerCreateDefinitionLiteralPointerUri, createDefinitionLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerCreateUriLiteralPointerUri, createUriLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerCreateValueLiteralPointerUri, createValueLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerGetLiteralUri, getLiteral)
	core.GetCore().AddFunction(LiteralPointerGetLiteralIdUri, getLiteralId)
	core.GetCore().AddFunction(LiteralPointerGetLiteralPointerRoleUri, getLiteralPointerRole)
	core.GetCore().AddFunction(LiteralPointerGetLiteralVersionUri, getLiteralVersion)
	core.GetCore().AddFunction(LiteralPointerSetLiteralUri, setLiteral)
}
