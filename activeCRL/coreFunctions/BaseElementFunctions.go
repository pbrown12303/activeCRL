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

var BaseElementFunctionsUri string = CoreFunctionsPrefix + "BaseElement"

var BaseElementDeleteUri string = CoreFunctionsPrefix + "BaseElement/Delete"
var BaseElementDeleteDeletedElementRefUri string = CoreFunctionsPrefix + "BaseElement/Delete/DeletedElementRef"

var BaseElementGetIdUri string = CoreFunctionsPrefix + "BaseElement/GetId"
var BaseElementGetIdSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetId/SourceBaseElementRef"
var BaseElementGetIdCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElement/GetId/CreatedLiteralRef"

var BaseElementGetNameUri string = CoreFunctionsPrefix + "BaseElement/GetName"
var BaseElementGetNameSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetName/SourceBaseElementRef"
var BaseElementGetNameCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElement/GetName/CreatedLiteralRef"

var BaseElementGetOwningElementUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement"
var BaseElementGetOwningElementSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement/SourceBaseElementRef"
var BaseElementGetOwningElementOwningElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement/TargetElementReference"

var BaseElementGetUriUri string = CoreFunctionsPrefix + "BaseElement/GetUri"
var BaseElementGetUriSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetUri/SourceBaseElementRef"
var BaseElementGetUriCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElement/GetUri/CreatedLiteralRef"

var BaseElementGetVersionUri string = CoreFunctionsPrefix + "BaseElement/GetVersion"
var BaseElementGetVersionSourceBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/GetVersion/SourceBaseElementRef"
var BaseElementGetVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "BaseElement/GetVersion/CreatedLiteralRef"

var BaseElementSetOwningElementUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement"
var BaseElementSetOwningElementOwningElementRefUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement/OwningElementRef"
var BaseElementSetOwningElementModifiedBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement/ModifiedBaseElementRef"

var BaseElementSetUriUri string = CoreFunctionsPrefix + "BaseElement/SetUri"
var BaseElementSetUriSourceUriRefUri string = CoreFunctionsPrefix + "BaseElement/SetUri/SourceUriRef"
var BaseElementSetUriModifiedBaseElementRefUri string = CoreFunctionsPrefix + "BaseElement/SetUri/ModifiedBaseElementRef"

func del(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementDeleteUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementDeleteDeletedElementRefUri, hl)
	if targetReference == nil {
		log.Printf("In Delete, the TargetReference was not found in the replicate")
		return
	}

	target := targetReference.GetReferencedBaseElement(hl)
	if target != nil {
		uOfD.DeleteBaseElement(target, hl)
	}

}

func getId(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetIdSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetId, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetIdCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetId, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	source := sourceReference.GetReferencedBaseElement(hl)
	if source != nil {
		createdLiteral.SetLiteralValue(core.GetName(source, hl), hl)
	}
}

func getName(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetNameUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetNameSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetName, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetNameCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetName, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	source := sourceReference.GetReferencedBaseElement(hl)
	if source != nil {
		createdLiteral.SetLiteralValue(core.GetName(source, hl), hl)
	}
}

func getOwningElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetOwningElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetOwningElementSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetOwningElement, the SourceReference was not found in the replicate")
		return
	}

	targetElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementGetOwningElementOwningElementRefUri, hl)
	if targetElementReference == nil {
		log.Printf("In GetOwningElement, the TargetElementReference was not found in the replicate")
		return
	}

	referencedElement := targetElementReference.GetReferencedElement(hl)
	source := sourceReference.GetReferencedBaseElement(hl)
	sourceOwner := core.GetOwningElement(source, hl)
	if sourceOwner != referencedElement {
		targetElementReference.SetReferencedElement(sourceOwner, hl)
	}
}

func getUri(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetUriUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetUriSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetUri, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetUriCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetUri, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	source := sourceReference.GetReferencedBaseElement(hl)
	if source != nil {
		createdLiteral.SetLiteralValue(core.GetName(source, hl), hl)
	}
}

func getVersion(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetVersionSourceBaseElementRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetVersion, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetVersionCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetVersion, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	source := sourceReference.GetReferencedBaseElement(hl)
	if source != nil {
		stringVersion := strconv.Itoa(source.GetVersion(hl))
		createdLiteral.SetLiteralValue(stringVersion, hl)
	}
}

func setOwningElement(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementSetOwningElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	owningElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementSetOwningElementOwningElementRefUri, hl)
	if owningElementReference == nil {
		log.Printf("In SetOwningElement, the OwningElementReference was not found in the replicate")
		return
	}

	targetBaseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementSetOwningElementModifiedBaseElementRefUri, hl)
	if targetBaseElementReference == nil {
		log.Printf("In SetOwningElement, the TargetBaseElementReference was not found in the replicate")
		return
	}

	targetBaseElement := targetBaseElementReference.GetReferencedBaseElement(hl)
	owner := owningElementReference.GetReferencedElement(hl)
	if targetBaseElement != nil {
		core.SetOwningElement(targetBaseElement, owner, hl)
	}
}

func setUri(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementSetUriUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	uriReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementSetUriSourceUriRefUri, hl)
	if uriReference == nil {
		log.Printf("In SetUri, the UriReference was not found in the replicate")
		return
	}

	targetBaseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementSetUriModifiedBaseElementRefUri, hl)
	if targetBaseElementReference == nil {
		log.Printf("In SetUri, the TargetBaseElementReference was not found in the replicate")
		return
	}

	targetBaseElement := targetBaseElementReference.GetReferencedBaseElement(hl)
	uriLiteral := uriReference.GetReferencedLiteral(hl)
	if targetBaseElement != nil {
		core.SetUri(targetBaseElement, uriLiteral.GetLiteralValue(hl), hl)
	}
}

// UpdateRecoveredCoreBaseElementFunctions() updates the representations of all BaseElementFunctions. The function is idempotent.
func BuildCoreBaseElementFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// BaseElementFunctions
	baseElementFunctions := uOfD.NewElement(hl, BaseElementFunctionsUri)
	core.SetName(baseElementFunctions, "BaseElementFunctions", hl)
	core.SetOwningElement(baseElementFunctions, coreFunctionsElement, hl)
	core.SetUri(baseElementFunctions, BaseElementFunctionsUri, hl)

	// Delete
	del := uOfD.NewElement(hl, BaseElementDeleteUri)
	core.SetName(del, "Delete", hl)
	core.SetOwningElement(del, baseElementFunctions, hl)
	core.SetUri(del, BaseElementDeleteUri, hl)

	// Delete.TargetReference
	delTargetReference := uOfD.NewBaseElementReference(hl, BaseElementDeleteDeletedElementRefUri)
	core.SetOwningElement(delTargetReference, del, hl)
	core.SetName(delTargetReference, "DeletedElementRef", hl)
	core.SetUri(delTargetReference, BaseElementDeleteDeletedElementRefUri, hl)

	// GetId
	getId := uOfD.NewElement(hl, BaseElementGetIdUri)
	core.SetName(getId, "GetId", hl)
	core.SetOwningElement(getId, baseElementFunctions, hl)
	core.SetUri(getId, BaseElementGetIdUri, hl)
	// GetId.SourceReference
	getIdSourceReference := uOfD.NewBaseElementReference(hl, BaseElementGetIdSourceBaseElementRefUri)
	core.SetOwningElement(getIdSourceReference, getId, hl)
	core.SetName(getIdSourceReference, "SourceBaseElementRef", hl)
	core.SetUri(getIdSourceReference, BaseElementGetIdSourceBaseElementRefUri, hl)
	// GetIdTargetLiteralReference
	getIdTargetReference := uOfD.NewLiteralReference(hl, BaseElementGetIdCreatedLiteralRefUri)
	core.SetOwningElement(getIdTargetReference, getId, hl)
	core.SetName(getIdTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getIdTargetReference, BaseElementGetIdCreatedLiteralRefUri, hl)

	// GetName
	getName := uOfD.NewElement(hl, BaseElementGetNameUri)
	core.SetName(getName, "GetName", hl)
	core.SetOwningElement(getName, baseElementFunctions, hl)
	core.SetUri(getName, BaseElementGetNameUri, hl)
	// GetName.SourceReference
	getNameSourceReference := uOfD.NewBaseElementReference(hl, BaseElementGetNameSourceBaseElementRefUri)
	core.SetOwningElement(getNameSourceReference, getName, hl)
	core.SetName(getNameSourceReference, "SourceBaseElementRef", hl)
	core.SetUri(getNameSourceReference, BaseElementGetNameSourceBaseElementRefUri, hl)
	// GetNameTargetLiteralReference
	getNameTargetReference := uOfD.NewLiteralReference(hl, BaseElementGetNameCreatedLiteralRefUri)
	core.SetOwningElement(getNameTargetReference, getName, hl)
	core.SetName(getNameTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getNameTargetReference, BaseElementGetNameCreatedLiteralRefUri, hl)

	// GetOwningElement
	getOwningElement := uOfD.NewElement(hl, BaseElementGetOwningElementUri)
	core.SetName(getOwningElement, "GetOwningElement", hl)
	core.SetOwningElement(getOwningElement, baseElementFunctions, hl)
	core.SetUri(getOwningElement, BaseElementGetOwningElementUri, hl)
	// GetOwningElement.SourceReference
	getOwningElementSourceReference := uOfD.NewBaseElementReference(hl, BaseElementGetOwningElementSourceBaseElementRefUri)
	core.SetOwningElement(getOwningElementSourceReference, getOwningElement, hl)
	core.SetName(getOwningElementSourceReference, "SourceBaseElementRef", hl)
	core.SetUri(getOwningElementSourceReference, BaseElementGetOwningElementSourceBaseElementRefUri, hl)
	// GetOwningElementTargetElementReference
	getOwningElementTargetReference := uOfD.NewElementReference(hl, BaseElementGetOwningElementOwningElementRefUri)
	core.SetOwningElement(getOwningElementTargetReference, getOwningElement, hl)
	core.SetName(getOwningElementTargetReference, "OwningElementRef", hl)
	core.SetUri(getOwningElementTargetReference, BaseElementGetOwningElementOwningElementRefUri, hl)

	// GetUri
	getUri := uOfD.NewElement(hl, BaseElementGetUriUri)
	core.SetName(getUri, "GetUri", hl)
	core.SetOwningElement(getUri, baseElementFunctions, hl)
	core.SetUri(getUri, BaseElementGetUriUri, hl)
	// GetUri.SourceReference
	getUriSourceReference := uOfD.NewBaseElementReference(hl, BaseElementGetUriSourceBaseElementRefUri)
	core.SetOwningElement(getUriSourceReference, getUri, hl)
	core.SetName(getUriSourceReference, "SourceBaseElementRef", hl)
	core.SetUri(getUriSourceReference, BaseElementGetUriSourceBaseElementRefUri, hl)
	// GetUriTargetLiteralReference
	getUriTargetReference := uOfD.NewLiteralReference(hl, BaseElementGetUriCreatedLiteralRefUri)
	core.SetOwningElement(getUriTargetReference, getUri, hl)
	core.SetName(getUriTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getUriTargetReference, BaseElementGetUriCreatedLiteralRefUri, hl)

	// GetVersion
	getVersion := uOfD.NewElement(hl, BaseElementGetVersionUri)
	core.SetName(getVersion, "GetVersion", hl)
	core.SetOwningElement(getVersion, baseElementFunctions, hl)
	core.SetUri(getVersion, BaseElementGetVersionUri, hl)
	// GetVersion.SourceReference
	getVersionSourceReference := uOfD.NewBaseElementReference(hl, BaseElementGetVersionSourceBaseElementRefUri)
	core.SetOwningElement(getVersionSourceReference, getVersion, hl)
	core.SetName(getVersionSourceReference, "SourceBaseElementRef", hl)
	core.SetUri(getVersionSourceReference, BaseElementGetVersionSourceBaseElementRefUri, hl)
	// GetVersionTargetLiteralReference
	getVersionTargetReference := uOfD.NewLiteralReference(hl, BaseElementGetVersionCreatedLiteralRefUri)
	core.SetOwningElement(getVersionTargetReference, getVersion, hl)
	core.SetName(getVersionTargetReference, "CreatedLiteralRef", hl)
	core.SetUri(getVersionTargetReference, BaseElementGetVersionCreatedLiteralRefUri, hl)

	// SetOwningElement
	setOwningElement := uOfD.NewElement(hl, BaseElementSetOwningElementUri)
	core.SetName(setOwningElement, "SetOwningElement", hl)
	core.SetOwningElement(setOwningElement, baseElementFunctions, hl)
	core.SetUri(setOwningElement, BaseElementSetOwningElementUri, hl)
	// SetOwningElement.SourceReference
	setOwningElementOwningElementReference := uOfD.NewElementReference(hl, BaseElementSetOwningElementOwningElementRefUri)
	core.SetOwningElement(setOwningElementOwningElementReference, setOwningElement, hl)
	core.SetName(setOwningElementOwningElementReference, "OwningElementRef", hl)
	core.SetUri(setOwningElementOwningElementReference, BaseElementSetOwningElementOwningElementRefUri, hl)
	// SetOwningElementTargetBaseElementReference
	setOwningElementTargetBaseElementReference := uOfD.NewBaseElementReference(hl, BaseElementSetOwningElementModifiedBaseElementRefUri)
	core.SetOwningElement(setOwningElementTargetBaseElementReference, setOwningElement, hl)
	core.SetName(setOwningElementTargetBaseElementReference, "ModifiedBaseElementRef", hl)
	core.SetUri(setOwningElementTargetBaseElementReference, BaseElementSetOwningElementModifiedBaseElementRefUri, hl)

	// SetUri
	setUri := uOfD.NewElement(hl, BaseElementSetUriUri)
	core.SetName(setUri, "SetUri", hl)
	core.SetOwningElement(setUri, baseElementFunctions, hl)
	core.SetUri(setUri, BaseElementSetUriUri, hl)
	// SetUri.UriReference
	setUriUriReference := uOfD.NewLiteralReference(hl, BaseElementSetUriSourceUriRefUri)
	core.SetOwningElement(setUriUriReference, setUri, hl)
	core.SetName(setUriUriReference, "SourceUriRef", hl)
	core.SetUri(setUriUriReference, BaseElementSetUriSourceUriRefUri, hl)
	// SetUriTargetBaseElementReference
	setUriTargetBaseElementReference := uOfD.NewBaseElementReference(hl, BaseElementSetUriModifiedBaseElementRefUri)
	core.SetOwningElement(setUriTargetBaseElementReference, setUri, hl)
	core.SetName(setUriTargetBaseElementReference, "ModifiedBaseElementRef", hl)
	core.SetUri(setUriTargetBaseElementReference, BaseElementSetUriModifiedBaseElementRefUri, hl)

}

func baseElementFunctionsInit() {
	core.GetCore().AddFunction(BaseElementDeleteUri, del)
	core.GetCore().AddFunction(BaseElementGetIdUri, getId)
	core.GetCore().AddFunction(BaseElementGetNameUri, getName)
	core.GetCore().AddFunction(BaseElementGetOwningElementUri, getOwningElement)
	core.GetCore().AddFunction(BaseElementGetUriUri, getUri)
	core.GetCore().AddFunction(BaseElementGetVersionUri, getVersion)
	core.GetCore().AddFunction(BaseElementSetOwningElementUri, setOwningElement)
	core.GetCore().AddFunction(BaseElementSetUriUri, setUri)
}
