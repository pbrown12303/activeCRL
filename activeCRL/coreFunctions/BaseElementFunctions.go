package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"strconv"
)

var BaseElementDeleteUri string = CoreFunctionsPrefix + "BaseElement/Delete"
var BaseElementDeleteTargetReferenceUri string = CoreFunctionsPrefix + "BaseElement/Delete/TargetReference"

var BaseElementGetIdUri string = CoreFunctionsPrefix + "BaseElement/GetId"
var BaseElementGetIdSourceReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetId/SourceReference"
var BaseElementGetIdTargetLiteralReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetId/TargetLiteralReference"

var BaseElementGetNameUri string = CoreFunctionsPrefix + "BaseElement/GetName"
var BaseElementGetNameSourceReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetName/SourceReference"
var BaseElementGetNameTargetLiteralReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetName/TargetLiteralReference"

var BaseElementGetOwningElementUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement"
var BaseElementGetOwningElementSourceReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement/SourceReference"
var BaseElementGetOwningElementTargetElementReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetOwningElement/TargetElementReference"

var BaseElementGetUriUri string = CoreFunctionsPrefix + "BaseElement/GetUri"
var BaseElementGetUriSourceReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetUri/SourceReference"
var BaseElementGetUriTargetLiteralReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetUri/TargetLiteralReference"

var BaseElementGetVersionUri string = CoreFunctionsPrefix + "BaseElement/GetVersion"
var BaseElementGetVersionSourceReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetVersion/SourceReference"
var BaseElementGetVersionTargetLiteralReferenceUri string = CoreFunctionsPrefix + "BaseElement/GetVersion/TargetLiteralReference"

var BaseElementSetOwningElementUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement"
var BaseElementSetOwningElementOwningElementReferenceUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement/OwningElementReference"
var BaseElementSetOwningElementTargetBaseElementReferenceUri string = CoreFunctionsPrefix + "BaseElement/SetOwningElement/TargetBaseElementReference"

var BaseElementSetUriUri string = CoreFunctionsPrefix + "BaseElement/SetUri"
var BaseElementSetUriUriReferenceUri string = CoreFunctionsPrefix + "BaseElement/SetUri/UriReference"
var BaseElementSetUriTargetBaseElementReferenceUri string = CoreFunctionsPrefix + "BaseElement/SetUri/TargetBaseElementReference"

func del(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementDeleteUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	targetReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementDeleteTargetReferenceUri, hl)
	if targetReference == nil {
		log.Printf("In Delete, the TargetReference was not found in the replicate")
		return
	}

	target := targetReference.GetReferencedBaseElement(hl)
	if target != nil {
		uOfD.DeleteBaseElement(target, hl)
	}

}

func getId(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetIdSourceReferenceUri, hl)
	if sourceReference == nil {
		log.Printf("In GetId, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetIdTargetLiteralReferenceUri, hl)
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

func getName(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetNameUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetNameSourceReferenceUri, hl)
	if sourceReference == nil {
		log.Printf("In GetName, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetNameTargetLiteralReferenceUri, hl)
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

func getOwningElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetOwningElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetOwningElementSourceReferenceUri, hl)
	if sourceReference == nil {
		log.Printf("In GetOwningElement, the SourceReference was not found in the replicate")
		return
	}

	targetElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementGetOwningElementTargetElementReferenceUri, hl)
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

func getUri(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetUriUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetUriSourceReferenceUri, hl)
	if sourceReference == nil {
		log.Printf("In GetUri, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetUriTargetLiteralReferenceUri, hl)
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

func getVersion(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementGetVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementGetVersionSourceReferenceUri, hl)
	if sourceReference == nil {
		log.Printf("In GetVersion, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementGetVersionTargetLiteralReferenceUri, hl)
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

func setOwningElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementSetOwningElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	owningElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, BaseElementSetOwningElementOwningElementReferenceUri, hl)
	if owningElementReference == nil {
		log.Printf("In SetOwningElement, the OwningElementReference was not found in the replicate")
		return
	}

	targetBaseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementSetOwningElementTargetBaseElementReferenceUri, hl)
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

func setUri(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementSetUriUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	uriReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementSetUriUriReferenceUri, hl)
	if uriReference == nil {
		log.Printf("In SetUri, the UriReference was not found in the replicate")
		return
	}

	targetBaseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementSetUriTargetBaseElementReferenceUri, hl)
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
func UpdateRecoveredCoreBaseElementFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// Delete
	del := uOfD.GetElementWithUri(BaseElementDeleteUri)
	if del == nil {
		del = uOfD.NewElement(hl)
		core.SetName(del, "Delete", hl)
		core.SetOwningElement(del, coreFunctionsElement, hl)
		core.SetUri(del, BaseElementDeleteUri, hl)
	}
	// Delete.TargetReference
	delTargetReference := core.GetChildBaseElementReferenceWithUri(del, BaseElementDeleteTargetReferenceUri, hl)
	if delTargetReference == nil {
		delTargetReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(delTargetReference, del, hl)
		core.SetName(delTargetReference, "TargetReference", hl)
		core.SetUri(delTargetReference, BaseElementDeleteTargetReferenceUri, hl)
	}

	// GetId
	getId := uOfD.GetElementWithUri(BaseElementGetIdUri)
	if getId == nil {
		getId = uOfD.NewElement(hl)
		core.SetName(getId, "GetId", hl)
		core.SetOwningElement(getId, coreFunctionsElement, hl)
		core.SetUri(getId, BaseElementGetIdUri, hl)
	}
	// GetId.SourceReference
	getIdSourceReference := core.GetChildBaseElementReferenceWithUri(getId, BaseElementGetIdSourceReferenceUri, hl)
	if getIdSourceReference == nil {
		getIdSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getIdSourceReference, getId, hl)
		core.SetName(getIdSourceReference, "SourceReference", hl)
		core.SetUri(getIdSourceReference, BaseElementGetIdSourceReferenceUri, hl)
	}
	// GetIdTargetLiteralReference
	getIdTargetReference := core.GetChildLiteralReferenceWithUri(getId, BaseElementGetIdTargetLiteralReferenceUri, hl)
	if getIdTargetReference == nil {
		getIdTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getIdTargetReference, getId, hl)
		core.SetName(getIdTargetReference, "TargetLiteralReference", hl)
		core.SetUri(getIdTargetReference, BaseElementGetIdTargetLiteralReferenceUri, hl)
	}

	// GetName
	getName := uOfD.GetElementWithUri(BaseElementGetNameUri)
	if getName == nil {
		getName = uOfD.NewElement(hl)
		core.SetName(getName, "GetName", hl)
		core.SetOwningElement(getName, coreFunctionsElement, hl)
		core.SetUri(getName, BaseElementGetNameUri, hl)
	}
	// GetName.SourceReference
	getNameSourceReference := core.GetChildBaseElementReferenceWithUri(getName, BaseElementGetNameSourceReferenceUri, hl)
	if getNameSourceReference == nil {
		getNameSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getNameSourceReference, getName, hl)
		core.SetName(getNameSourceReference, "SourceReference", hl)
		core.SetUri(getNameSourceReference, BaseElementGetNameSourceReferenceUri, hl)
	}
	// GetNameTargetLiteralReference
	getNameTargetReference := core.GetChildLiteralReferenceWithUri(getName, BaseElementGetNameTargetLiteralReferenceUri, hl)
	if getNameTargetReference == nil {
		getNameTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getNameTargetReference, getName, hl)
		core.SetName(getNameTargetReference, "TargetLiteralReference", hl)
		core.SetUri(getNameTargetReference, BaseElementGetNameTargetLiteralReferenceUri, hl)
	}

	// GetOwningElement
	getOwningElement := uOfD.GetElementWithUri(BaseElementGetOwningElementUri)
	if getOwningElement == nil {
		getOwningElement = uOfD.NewElement(hl)
		core.SetName(getOwningElement, "GetOwningElement", hl)
		core.SetOwningElement(getOwningElement, coreFunctionsElement, hl)
		core.SetUri(getOwningElement, BaseElementGetOwningElementUri, hl)
	}
	// GetOwningElement.SourceReference
	getOwningElementSourceReference := core.GetChildBaseElementReferenceWithUri(getOwningElement, BaseElementGetOwningElementSourceReferenceUri, hl)
	if getOwningElementSourceReference == nil {
		getOwningElementSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getOwningElementSourceReference, getOwningElement, hl)
		core.SetName(getOwningElementSourceReference, "SourceReference", hl)
		core.SetUri(getOwningElementSourceReference, BaseElementGetOwningElementSourceReferenceUri, hl)
	}
	// GetOwningElementTargetElementReference
	getOwningElementTargetReference := core.GetChildElementReferenceWithUri(getOwningElement, BaseElementGetOwningElementTargetElementReferenceUri, hl)
	if getOwningElementTargetReference == nil {
		getOwningElementTargetReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getOwningElementTargetReference, getOwningElement, hl)
		core.SetName(getOwningElementTargetReference, "TargetElementReference", hl)
		core.SetUri(getOwningElementTargetReference, BaseElementGetOwningElementTargetElementReferenceUri, hl)
	}

	// GetUri
	getUri := uOfD.GetElementWithUri(BaseElementGetUriUri)
	if getUri == nil {
		getUri = uOfD.NewElement(hl)
		core.SetName(getUri, "GetUri", hl)
		core.SetOwningElement(getUri, coreFunctionsElement, hl)
		core.SetUri(getUri, BaseElementGetUriUri, hl)
	}
	// GetUri.SourceReference
	getUriSourceReference := core.GetChildBaseElementReferenceWithUri(getUri, BaseElementGetUriSourceReferenceUri, hl)
	if getUriSourceReference == nil {
		getUriSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getUriSourceReference, getUri, hl)
		core.SetName(getUriSourceReference, "SourceReference", hl)
		core.SetUri(getUriSourceReference, BaseElementGetUriSourceReferenceUri, hl)
	}
	// GetUriTargetLiteralReference
	getUriTargetReference := core.GetChildLiteralReferenceWithUri(getUri, BaseElementGetUriTargetLiteralReferenceUri, hl)
	if getUriTargetReference == nil {
		getUriTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getUriTargetReference, getUri, hl)
		core.SetName(getUriTargetReference, "TargetLiteralReference", hl)
		core.SetUri(getUriTargetReference, BaseElementGetUriTargetLiteralReferenceUri, hl)
	}

	// GetVersion
	getVersion := uOfD.GetElementWithUri(BaseElementGetVersionUri)
	if getVersion == nil {
		getVersion = uOfD.NewElement(hl)
		core.SetName(getVersion, "GetVersion", hl)
		core.SetOwningElement(getVersion, coreFunctionsElement, hl)
		core.SetUri(getVersion, BaseElementGetVersionUri, hl)
	}
	// GetVersion.SourceReference
	getVersionSourceReference := core.GetChildBaseElementReferenceWithUri(getVersion, BaseElementGetVersionSourceReferenceUri, hl)
	if getVersionSourceReference == nil {
		getVersionSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getVersionSourceReference, getVersion, hl)
		core.SetName(getVersionSourceReference, "SourceReference", hl)
		core.SetUri(getVersionSourceReference, BaseElementGetVersionSourceReferenceUri, hl)
	}
	// GetVersionTargetLiteralReference
	getVersionTargetReference := core.GetChildLiteralReferenceWithUri(getVersion, BaseElementGetVersionTargetLiteralReferenceUri, hl)
	if getVersionTargetReference == nil {
		getVersionTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getVersionTargetReference, getVersion, hl)
		core.SetName(getVersionTargetReference, "TargetLiteralReference", hl)
		core.SetUri(getVersionTargetReference, BaseElementGetVersionTargetLiteralReferenceUri, hl)
	}

	// SetOwningElement
	setOwningElement := uOfD.GetElementWithUri(BaseElementSetOwningElementUri)
	if setOwningElement == nil {
		setOwningElement = uOfD.NewElement(hl)
		core.SetName(setOwningElement, "SetOwningElement", hl)
		core.SetOwningElement(setOwningElement, coreFunctionsElement, hl)
		core.SetUri(setOwningElement, BaseElementSetOwningElementUri, hl)
	}
	// SetOwningElement.SourceReference
	setOwningElementOwningElementReference := core.GetChildElementReferenceWithUri(setOwningElement, BaseElementSetOwningElementOwningElementReferenceUri, hl)
	if setOwningElementOwningElementReference == nil {
		setOwningElementOwningElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setOwningElementOwningElementReference, setOwningElement, hl)
		core.SetName(setOwningElementOwningElementReference, "OwningElementReference", hl)
		core.SetUri(setOwningElementOwningElementReference, BaseElementSetOwningElementOwningElementReferenceUri, hl)
	}
	// SetOwningElementTargetBaseElementReference
	setOwningElementTargetBaseElementReference := core.GetChildBaseElementReferenceWithUri(setOwningElement, BaseElementSetOwningElementTargetBaseElementReferenceUri, hl)
	if setOwningElementTargetBaseElementReference == nil {
		setOwningElementTargetBaseElementReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(setOwningElementTargetBaseElementReference, setOwningElement, hl)
		core.SetName(setOwningElementTargetBaseElementReference, "TargetBaseElementReference", hl)
		core.SetUri(setOwningElementTargetBaseElementReference, BaseElementSetOwningElementTargetBaseElementReferenceUri, hl)
	}

	// SetUri
	setUri := uOfD.GetElementWithUri(BaseElementSetUriUri)
	if setUri == nil {
		setUri = uOfD.NewElement(hl)
		core.SetName(setUri, "SetUri", hl)
		core.SetOwningElement(setUri, coreFunctionsElement, hl)
		core.SetUri(setUri, BaseElementSetUriUri, hl)
	}
	// SetUri.UriReference
	setUriUriReference := core.GetChildLiteralReferenceWithUri(setUri, BaseElementSetUriUriReferenceUri, hl)
	if setUriUriReference == nil {
		setUriUriReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(setUriUriReference, setUri, hl)
		core.SetName(setUriUriReference, "UriReference", hl)
		core.SetUri(setUriUriReference, BaseElementSetUriUriReferenceUri, hl)
	}
	// SetUriTargetBaseElementReference
	setUriTargetBaseElementReference := core.GetChildBaseElementReferenceWithUri(setUri, BaseElementSetUriTargetBaseElementReferenceUri, hl)
	if setUriTargetBaseElementReference == nil {
		setUriTargetBaseElementReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(setUriTargetBaseElementReference, setUri, hl)
		core.SetName(setUriTargetBaseElementReference, "TargetBaseElementReference", hl)
		core.SetUri(setUriTargetBaseElementReference, BaseElementSetUriTargetBaseElementReferenceUri, hl)
	}

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
