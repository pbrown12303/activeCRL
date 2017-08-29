package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"strconv"
)

var BaseElementPointerCreateUri string = CoreFunctionsPrefix + "BaseElementPointer/Create"
var BaseElementPointerCreateCreatedBaseElementPointerReferenceUri = CoreFunctionsPrefix + "BaseElementPointer/Create/CreatedBaseElementPointerReference"

var BaseElementPointerGetBaseElementUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement"
var BaseElementPointerGetBaseElementSourceReferenceUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement/SourceReference"
var BaseElementPointerGetBaseElementTargetBaseElementReferenceUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElement/TargetBaseElementReference"

var BaseElementPointerGetBaseElementIdUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId"
var BaseElementPointerGetBaseElementIdSourceReferenceUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId/SourceReference"
var BaseElementPointerGetBaseElementIdTargetLiteralReferenceUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementId/TargetLiteralReference"

var BaseElementPointerGetBaseElementVersionUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion"
var BaseElementPointerGetBaseElementVersionSourceReferenceUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion/SourceReference"
var BaseElementPointerGetBaseElementVersionTargetLiteralReferenceUri string = CoreFunctionsPrefix + "BaseElementPointer/GetBaseElementVersion/TargetLiteralReference"

var BaseElementPointerSetBaseElementUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement"
var BaseElementPointerSetBaseElementBaseElementReferenceUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement/BaseElementReference"
var BaseElementPointerSetBaseElementTargetBaseElementPointerReferenceUri string = CoreFunctionsPrefix + "BaseElementPointer/SetBaseElement/TargetBaseElementPointerReference"

func createBaseElementPointer(element core.Element, changeNotification *core.ChangeNotification) {
	//	log.Printf("In createBaseElementPointer")
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdBaseElementPointerReference := core.GetChildBaseElementReferenceWithAncestorUri(element, BaseElementPointerCreateCreatedBaseElementPointerReferenceUri, hl)
	if createdBaseElementPointerReference == nil {
		createdBaseElementPointerReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(createdBaseElementPointerReference, element, hl)
		core.SetName(createdBaseElementPointerReference, "CreatedBaseElementPointerReference", hl)
		rootCreatedElementReference := uOfD.GetBaseElementReferenceWithUri(BaseElementPointerCreateCreatedBaseElementPointerReferenceUri)
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

func getBaseElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementSourceReferenceUri, hl)
	if sourceReference == nil {
		log.Printf("In GetBaseElement, the SourceReference was not found in the replicate")
		return
	}

	targetElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementTargetBaseElementReferenceUri, hl)
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

func getBaseElementId(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIdSourceReferenceUri, hl)
	if sourceReference == nil {
		log.Printf("In GetBaseElementId, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementIdTargetLiteralReferenceUri, hl)
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

func getBaseElementVersion(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementVersionSourceReferenceUri, hl)
	if sourceReference == nil {
		log.Printf("In GetBaseElementVersion, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, BaseElementPointerGetBaseElementVersionTargetLiteralReferenceUri, hl)
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

func setBaseElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(BaseElementPointerSetBaseElementUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	baseElementReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerSetBaseElementBaseElementReferenceUri, hl)
	if baseElementReference == nil {
		log.Printf("In SetBaseElement, the BaseElementReference was not found in the replicate")
		return
	}

	targetBaseElementPointerReference := core.GetChildBaseElementReferenceWithAncestorUri(replicate, BaseElementPointerSetBaseElementTargetBaseElementPointerReferenceUri, hl)
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

func UpdateRecoveredCoreBaseElementPointerFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// CreateElement
	createBaseElementPointer := uOfD.GetElementWithUri(BaseElementPointerCreateUri)
	if createBaseElementPointer == nil {
		createBaseElementPointer = uOfD.NewElement(hl)
		core.SetOwningElement(createBaseElementPointer, coreFunctionsElement, hl)
		core.SetName(createBaseElementPointer, "CreateBaseElementPointer", hl)
		core.SetUri(createBaseElementPointer, BaseElementPointerCreateUri, hl)
	}
	// CreatedElementReference
	createdBaseElementPointerReference := core.GetChildBaseElementReferenceWithUri(createBaseElementPointer, BaseElementPointerCreateCreatedBaseElementPointerReferenceUri, hl)
	if createdBaseElementPointerReference == nil {
		createdBaseElementPointerReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(createdBaseElementPointerReference, createBaseElementPointer, hl)
		core.SetName(createdBaseElementPointerReference, "CreatedBaseElementPointerReference", hl)
		core.SetUri(createdBaseElementPointerReference, BaseElementPointerCreateCreatedBaseElementPointerReferenceUri, hl)
	}

	// GetBaseElement
	getBaseElement := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementUri)
	if getBaseElement == nil {
		getBaseElement = uOfD.NewElement(hl)
		core.SetName(getBaseElement, "GetBaseElement", hl)
		core.SetOwningElement(getBaseElement, coreFunctionsElement, hl)
		core.SetUri(getBaseElement, BaseElementPointerGetBaseElementUri, hl)
	}
	// GetBaseElement.SourceReference
	getBaseElementSourceReference := core.GetChildBaseElementReferenceWithUri(getBaseElement, BaseElementPointerGetBaseElementSourceReferenceUri, hl)
	if getBaseElementSourceReference == nil {
		getBaseElementSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementSourceReference, getBaseElement, hl)
		core.SetName(getBaseElementSourceReference, "SourceReference", hl)
		core.SetUri(getBaseElementSourceReference, BaseElementPointerGetBaseElementSourceReferenceUri, hl)
	}
	// GetBaseElementTargetBaseElementReference
	getBaseElementTargetReference := core.GetChildBaseElementReferenceWithUri(getBaseElement, BaseElementPointerGetBaseElementTargetBaseElementReferenceUri, hl)
	if getBaseElementTargetReference == nil {
		getBaseElementTargetReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementTargetReference, getBaseElement, hl)
		core.SetName(getBaseElementTargetReference, "TargetBaseElementReference", hl)
		core.SetUri(getBaseElementTargetReference, BaseElementPointerGetBaseElementTargetBaseElementReferenceUri, hl)
	}

	// GetBaseElementId
	getBaseElementId := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementIdUri)
	if getBaseElementId == nil {
		getBaseElementId = uOfD.NewElement(hl)
		core.SetName(getBaseElementId, "GetBaseElementId", hl)
		core.SetOwningElement(getBaseElementId, coreFunctionsElement, hl)
		core.SetUri(getBaseElementId, BaseElementPointerGetBaseElementIdUri, hl)
	}
	// GetBaseElementId.SourceReference
	getBaseElementIdSourceReference := core.GetChildBaseElementReferenceWithUri(getBaseElementId, BaseElementPointerGetBaseElementIdSourceReferenceUri, hl)
	if getBaseElementIdSourceReference == nil {
		getBaseElementIdSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementIdSourceReference, getBaseElementId, hl)
		core.SetName(getBaseElementIdSourceReference, "SourceReference", hl)
		core.SetUri(getBaseElementIdSourceReference, BaseElementPointerGetBaseElementIdSourceReferenceUri, hl)
	}
	// GetBaseElementIdTargetLiteralReference
	getBaseElementIdTargetReference := core.GetChildLiteralReferenceWithUri(getBaseElementId, BaseElementPointerGetBaseElementIdTargetLiteralReferenceUri, hl)
	if getBaseElementIdTargetReference == nil {
		getBaseElementIdTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getBaseElementIdTargetReference, getBaseElementId, hl)
		core.SetName(getBaseElementIdTargetReference, "TargetLiteralReference", hl)
		core.SetUri(getBaseElementIdTargetReference, BaseElementPointerGetBaseElementIdTargetLiteralReferenceUri, hl)
	}

	// GetBaseElementVersion
	getBaseElementVersion := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementVersionUri)
	if getBaseElementVersion == nil {
		getBaseElementVersion = uOfD.NewElement(hl)
		core.SetName(getBaseElementVersion, "GetBaseElementVersion", hl)
		core.SetOwningElement(getBaseElementVersion, coreFunctionsElement, hl)
		core.SetUri(getBaseElementVersion, BaseElementPointerGetBaseElementVersionUri, hl)
	}
	// GetBaseElementVersion.SourceReference
	getBaseElementVersionSourceReference := core.GetChildBaseElementReferenceWithUri(getBaseElementVersion, BaseElementPointerGetBaseElementVersionSourceReferenceUri, hl)
	if getBaseElementVersionSourceReference == nil {
		getBaseElementVersionSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementVersionSourceReference, getBaseElementVersion, hl)
		core.SetName(getBaseElementVersionSourceReference, "SourceReference", hl)
		core.SetUri(getBaseElementVersionSourceReference, BaseElementPointerGetBaseElementVersionSourceReferenceUri, hl)
	}
	// GetBaseElementVersionTargetLiteralReference
	getBaseElementVersionTargetReference := core.GetChildLiteralReferenceWithUri(getBaseElementVersion, BaseElementPointerGetBaseElementVersionTargetLiteralReferenceUri, hl)
	if getBaseElementVersionTargetReference == nil {
		getBaseElementVersionTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getBaseElementVersionTargetReference, getBaseElementVersion, hl)
		core.SetName(getBaseElementVersionTargetReference, "TargetLiteralReference", hl)
		core.SetUri(getBaseElementVersionTargetReference, BaseElementPointerGetBaseElementVersionTargetLiteralReferenceUri, hl)
	}

	// SetBaseElement
	setBaseElement := uOfD.GetElementWithUri(BaseElementPointerSetBaseElementUri)
	if setBaseElement == nil {
		setBaseElement = uOfD.NewElement(hl)
		core.SetName(setBaseElement, "SetBaseElement", hl)
		core.SetOwningElement(setBaseElement, coreFunctionsElement, hl)
		core.SetUri(setBaseElement, BaseElementPointerSetBaseElementUri, hl)
	}
	// SetBaseElement.BaseElementReference
	setBaseElementBaseElementReference := core.GetChildBaseElementReferenceWithUri(setBaseElement, BaseElementPointerSetBaseElementBaseElementReferenceUri, hl)
	if setBaseElementBaseElementReference == nil {
		setBaseElementBaseElementReference = uOfD.NewBaseElementReference(hl)
		core.SetName(setBaseElementBaseElementReference, "BaseElementReference", hl)
		core.SetOwningElement(setBaseElementBaseElementReference, setBaseElement, hl)
		core.SetUri(setBaseElementBaseElementReference, BaseElementPointerSetBaseElementBaseElementReferenceUri, hl)
	}
	setBaseElementTargetBaseElementPointerReference := core.GetChildBaseElementReferenceWithUri(setBaseElement, BaseElementPointerSetBaseElementTargetBaseElementPointerReferenceUri, hl)
	if setBaseElementTargetBaseElementPointerReference == nil {
		setBaseElementTargetBaseElementPointerReference = uOfD.NewBaseElementReference(hl)
		core.SetName(setBaseElementTargetBaseElementPointerReference, "TargetBaseElementPointerReference", hl)
		core.SetOwningElement(setBaseElementTargetBaseElementPointerReference, setBaseElement, hl)
		core.SetUri(setBaseElementTargetBaseElementPointerReference, BaseElementPointerSetBaseElementTargetBaseElementPointerReferenceUri, hl)
	}
}

func baseElementPointerFunctionsInit() {
	core.GetCore().AddFunction(BaseElementPointerCreateUri, createBaseElementPointer)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementUri, getBaseElement)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementIdUri, getBaseElementId)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementVersionUri, getBaseElementVersion)
	core.GetCore().AddFunction(BaseElementPointerSetBaseElementUri, setBaseElement)
}
