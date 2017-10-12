// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"sync"
)

var LiteralFunctionsUri string = CoreFunctionsPrefix + "LiteralFunctions"

var LiteralCreateUri string = CoreFunctionsPrefix + "Literal/Create"
var LiteralCreateCreatedLiteralRefUri string = CoreFunctionsPrefix + "Literal/Create/CreatedLiteralRef"

var LiteralGetLiteralValueUri string = CoreFunctionsPrefix + "Literal/GetLiteralValue"
var LiteralGetLiteralValueSourceLiteralRefUri string = CoreFunctionsPrefix + "Literal/GetLiteralValue/SourceLiteralRef"
var LiteralGetLiteralValueCreatedLiteralRefUri string = CoreFunctionsPrefix + "Literal/GetLiteralValue/CreatedLiteralRef"

var LiteralSetLiteralValueUri string = CoreFunctionsPrefix + "Literal/SetLiteralValue"
var LiteralSetLiteralValueSourceLiteralRefUri string = CoreFunctionsPrefix + "Literal/SetLiteralValue/SourceLiteralRef"
var LiteralSetLiteralValueModifiedLiteralRefUri string = CoreFunctionsPrefix + "Literal/SetLiteralValue/ModifiedLiteralRef"

func createLiteral(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(element, LiteralCreateCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		createdLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(createdLiteralRef, element, hl)
		core.SetName(createdLiteralRef, "CreatedLiteralRef", hl)
		rootCreatedLiteralRef := uOfD.GetLiteralReferenceWithUri(LiteralCreateCreatedLiteralRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralRef, hl)
		refinement.SetRefinedElement(createdLiteralRef, hl)
		refinement.SetAbstractElement(rootCreatedLiteralRef, hl)
	}
	createdLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		createdLiteralRef.SetReferencedLiteral(createdLiteral, hl)
	}
}

func getLiteralValue(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralGetLiteralValueUri)
	if original == nil {
		log.Printf("In GetLiteralValue the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralGetLiteralValueSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		log.Printf("In GetLiteralValue, the SourceLiteralRef was not found in the replicate")
		return
	}

	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralGetLiteralValueCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		log.Printf("In GetLiteralValue, the CreatedLiteralRef was not found in the replicate")
		return
	}

	currentLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if currentLiteral == nil {
		currentLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(currentLiteral, createdLiteralRef, hl)
		createdLiteralRef.SetReferencedLiteral(currentLiteral, hl)
	}

	var sourceLiteralValue string = ""
	sourceLiteral := sourceLiteralRef.GetReferencedLiteral(hl)
	if sourceLiteral != nil {
		sourceLiteralValue = sourceLiteral.GetLiteralValue(hl)
	}
	if sourceLiteralValue != currentLiteral.GetLiteralValue(hl) {
		currentLiteral.SetLiteralValue(sourceLiteralValue, hl)
	}
}

func setLiteralValue(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralSetLiteralValueUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralSetLiteralValueSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		log.Printf("In SetLiteralValue, the SourceLiteralRef was not found in the replicate")
		return
	}

	modifiedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralSetLiteralValueModifiedLiteralRefUri, hl)
	if modifiedLiteralRef == nil {
		log.Printf("In SetLiteralValue, the ModifiedLiteralRef was not found in the replicate")
		return
	}

	modifiedLiteral := modifiedLiteralRef.GetReferencedLiteral(hl)
	sourceLiteral := sourceLiteralRef.GetReferencedLiteral(hl)
	if modifiedLiteral != nil {
		modifiedLiteral.SetLiteralValue(sourceLiteral.GetLiteralValue(hl), hl)
	}
}

func BuildCoreLiteralFunctions(coreFunctionsElement core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralFunctions
	literalFunctions := uOfD.NewElement(hl, LiteralFunctionsUri)
	core.SetOwningElement(literalFunctions, coreFunctionsElement, hl)
	core.SetName(literalFunctions, "LiteralFunctions", hl)
	core.SetUri(literalFunctions, LiteralFunctionsUri, hl)

	// CreateLiteral
	createLiteral := uOfD.NewElement(hl, LiteralCreateUri)
	core.SetOwningElement(createLiteral, literalFunctions, hl)
	core.SetName(createLiteral, "CreateLiteral", hl)
	core.SetUri(createLiteral, LiteralCreateUri, hl)
	// CreatedLiteralReference
	createdLiteralRef0 := uOfD.NewLiteralReference(hl, LiteralCreateCreatedLiteralRefUri)
	core.SetOwningElement(createdLiteralRef0, createLiteral, hl)
	core.SetName(createdLiteralRef0, "CreatedLiteralReference", hl)
	core.SetUri(createdLiteralRef0, LiteralCreateCreatedLiteralRefUri, hl)

	// GetLiteralValue
	getLiteralValue := uOfD.NewElement(hl, LiteralGetLiteralValueUri)
	core.SetOwningElement(getLiteralValue, literalFunctions, hl)
	core.SetName(getLiteralValue, "GetLiteralValue", hl)
	core.SetUri(getLiteralValue, LiteralGetLiteralValueUri, hl)
	// SourceLiteralRef
	sourceLiteralRef := uOfD.NewLiteralReference(hl, LiteralGetLiteralValueSourceLiteralRefUri)
	core.SetOwningElement(sourceLiteralRef, getLiteralValue, hl)
	core.SetName(sourceLiteralRef, "SourceLiteralRef", hl)
	core.SetUri(sourceLiteralRef, LiteralGetLiteralValueSourceLiteralRefUri, hl)
	// CreatedLiteralRef
	createdLiteralRef1 := uOfD.NewLiteralReference(hl, LiteralGetLiteralValueCreatedLiteralRefUri)
	core.SetOwningElement(createdLiteralRef1, getLiteralValue, hl)
	core.SetName(createdLiteralRef1, "CreatedLiteralRef", hl)
	core.SetUri(createdLiteralRef1, LiteralGetLiteralValueCreatedLiteralRefUri, hl)

	// SetLiteralValue
	setLiteralValue := uOfD.NewElement(hl, LiteralSetLiteralValueUri)
	core.SetName(setLiteralValue, "SetLiteralValue", hl)
	core.SetOwningElement(setLiteralValue, literalFunctions, hl)
	core.SetUri(setLiteralValue, LiteralSetLiteralValueUri, hl)
	// SetLiteralValue.SourceLiteralRef
	setLiteralValueSourceLiteralRef := uOfD.NewLiteralReference(hl, LiteralSetLiteralValueSourceLiteralRefUri)
	core.SetOwningElement(setLiteralValueSourceLiteralRef, setLiteralValue, hl)
	core.SetName(setLiteralValueSourceLiteralRef, "SourceLiteralRefRef", hl)
	core.SetUri(setLiteralValueSourceLiteralRef, LiteralSetLiteralValueSourceLiteralRefUri, hl)
	// SetLiteralValueModifiedLiteralReference
	setLiteralValueModifiedLiteralRef := uOfD.NewLiteralReference(hl, LiteralSetLiteralValueModifiedLiteralRefUri)
	core.SetOwningElement(setLiteralValueModifiedLiteralRef, setLiteralValue, hl)
	core.SetName(setLiteralValueModifiedLiteralRef, "ModifiedLiteralRef", hl)
	core.SetUri(setLiteralValueModifiedLiteralRef, LiteralSetLiteralValueModifiedLiteralRefUri, hl)
}

func literalFunctionsInit() {
	core.GetCore().AddFunction(LiteralCreateUri, createLiteral)
	core.GetCore().AddFunction(LiteralGetLiteralValueUri, getLiteralValue)
	core.GetCore().AddFunction(LiteralSetLiteralValueUri, setLiteralValue)
}
