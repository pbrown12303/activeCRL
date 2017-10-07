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
		core.SetName(createdBaseElementPointerReference, "CreatedBaseElementPointerReference", hl)
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

func UpdateRecoveredCoreBaseElementPointerFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// BaseElementPointerFunctions
	baseElementPointerFunctions := uOfD.GetElementWithUri(BaseElementPointerFunctionsUri)
	if baseElementPointerFunctions == nil {
		baseElementPointerFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(baseElementPointerFunctions, coreFunctionsElement, hl)
		core.SetName(baseElementPointerFunctions, "BaseElementPointerFunctions", hl)
		core.SetUri(baseElementPointerFunctions, BaseElementPointerFunctionsUri, hl)
	}

	// CreateElement
	createBaseElementPointer := uOfD.GetElementWithUri(BaseElementPointerCreateUri)
	if createBaseElementPointer == nil {
		createBaseElementPointer = uOfD.NewElement(hl)
		core.SetOwningElement(createBaseElementPointer, baseElementPointerFunctions, hl)
		core.SetName(createBaseElementPointer, "CreateBaseElementPointer", hl)
		core.SetUri(createBaseElementPointer, BaseElementPointerCreateUri, hl)
	}
	// CreatedElementReference
	createdBaseElementPointerReference := core.GetChildBaseElementReferenceWithUri(createBaseElementPointer, BaseElementPointerCreateCreatedBaseElementPointerRefUri, hl)
	if createdBaseElementPointerReference == nil {
		createdBaseElementPointerReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(createdBaseElementPointerReference, createBaseElementPointer, hl)
		core.SetName(createdBaseElementPointerReference, "CreatedBaseElementPointerRef", hl)
		core.SetUri(createdBaseElementPointerReference, BaseElementPointerCreateCreatedBaseElementPointerRefUri, hl)
	}

	// GetBaseElement
	getBaseElement := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementUri)
	if getBaseElement == nil {
		getBaseElement = uOfD.NewElement(hl)
		core.SetName(getBaseElement, "GetBaseElement", hl)
		core.SetOwningElement(getBaseElement, baseElementPointerFunctions, hl)
		core.SetUri(getBaseElement, BaseElementPointerGetBaseElementUri, hl)
	}
	// GetBaseElement.SourceReference
	getBaseElementSourceReference := core.GetChildBaseElementReferenceWithUri(getBaseElement, BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri, hl)
	if getBaseElementSourceReference == nil {
		getBaseElementSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementSourceReference, getBaseElement, hl)
		core.SetName(getBaseElementSourceReference, "SourceBaseElementPointerRef", hl)
		core.SetUri(getBaseElementSourceReference, BaseElementPointerGetBaseElementSourceBaseElementPointerRefUri, hl)
	}
	// GetBaseElementTargetBaseElementReference
	getBaseElementTargetReference := core.GetChildBaseElementReferenceWithUri(getBaseElement, BaseElementPointerGetBaseElementIndicatedBaseElementRefUri, hl)
	if getBaseElementTargetReference == nil {
		getBaseElementTargetReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementTargetReference, getBaseElement, hl)
		core.SetName(getBaseElementTargetReference, "IndicatedBaseElementRef", hl)
		core.SetUri(getBaseElementTargetReference, BaseElementPointerGetBaseElementIndicatedBaseElementRefUri, hl)
	}

	// GetBaseElementId
	getBaseElementId := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementIdUri)
	if getBaseElementId == nil {
		getBaseElementId = uOfD.NewElement(hl)
		core.SetName(getBaseElementId, "GetBaseElementId", hl)
		core.SetOwningElement(getBaseElementId, baseElementPointerFunctions, hl)
		core.SetUri(getBaseElementId, BaseElementPointerGetBaseElementIdUri, hl)
	}
	// GetBaseElementId.SourceReference
	getBaseElementIdSourceReference := core.GetChildBaseElementReferenceWithUri(getBaseElementId, BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri, hl)
	if getBaseElementIdSourceReference == nil {
		getBaseElementIdSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementIdSourceReference, getBaseElementId, hl)
		core.SetName(getBaseElementIdSourceReference, "SourceBaseElementPointerRef", hl)
		core.SetUri(getBaseElementIdSourceReference, BaseElementPointerGetBaseElementIdSourceBaseElementPointerRefUri, hl)
	}
	// GetBaseElementIdTargetLiteralReference
	getBaseElementIdTargetReference := core.GetChildLiteralReferenceWithUri(getBaseElementId, BaseElementPointerGetBaseElementIdCreatedLiteralUri, hl)
	if getBaseElementIdTargetReference == nil {
		getBaseElementIdTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getBaseElementIdTargetReference, getBaseElementId, hl)
		core.SetName(getBaseElementIdTargetReference, "CreatedLiteralRef", hl)
		core.SetUri(getBaseElementIdTargetReference, BaseElementPointerGetBaseElementIdCreatedLiteralUri, hl)
	}

	// GetBaseElementVersion
	getBaseElementVersion := uOfD.GetElementWithUri(BaseElementPointerGetBaseElementVersionUri)
	if getBaseElementVersion == nil {
		getBaseElementVersion = uOfD.NewElement(hl)
		core.SetName(getBaseElementVersion, "GetBaseElementVersion", hl)
		core.SetOwningElement(getBaseElementVersion, baseElementPointerFunctions, hl)
		core.SetUri(getBaseElementVersion, BaseElementPointerGetBaseElementVersionUri, hl)
	}
	// GetBaseElementVersion.SourceReference
	getBaseElementVersionSourceReference := core.GetChildBaseElementReferenceWithUri(getBaseElementVersion, BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri, hl)
	if getBaseElementVersionSourceReference == nil {
		getBaseElementVersionSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getBaseElementVersionSourceReference, getBaseElementVersion, hl)
		core.SetName(getBaseElementVersionSourceReference, "SourceBaseElementPointerRef", hl)
		core.SetUri(getBaseElementVersionSourceReference, BaseElementPointerGetBaseElementVersionSourceBaseElementPointerRefUri, hl)
	}
	// GetBaseElementVersionTargetLiteralReference
	getBaseElementVersionTargetReference := core.GetChildLiteralReferenceWithUri(getBaseElementVersion, BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri, hl)
	if getBaseElementVersionTargetReference == nil {
		getBaseElementVersionTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getBaseElementVersionTargetReference, getBaseElementVersion, hl)
		core.SetName(getBaseElementVersionTargetReference, "CreatedLiteralRef", hl)
		core.SetUri(getBaseElementVersionTargetReference, BaseElementPointerGetBaseElementVersionCreatedLiteralRefUri, hl)
	}

	// SetBaseElement
	setBaseElement := uOfD.GetElementWithUri(BaseElementPointerSetBaseElementUri)
	if setBaseElement == nil {
		setBaseElement = uOfD.NewElement(hl)
		core.SetName(setBaseElement, "SetBaseElement", hl)
		core.SetOwningElement(setBaseElement, baseElementPointerFunctions, hl)
		core.SetUri(setBaseElement, BaseElementPointerSetBaseElementUri, hl)
	}
	// SetBaseElement.BaseElementReference
	setBaseElementBaseElementReference := core.GetChildBaseElementReferenceWithUri(setBaseElement, BaseElementPointerSetBaseElementBaseElementRefUri, hl)
	if setBaseElementBaseElementReference == nil {
		setBaseElementBaseElementReference = uOfD.NewBaseElementReference(hl)
		core.SetName(setBaseElementBaseElementReference, "BaseElementRef", hl)
		core.SetOwningElement(setBaseElementBaseElementReference, setBaseElement, hl)
		core.SetUri(setBaseElementBaseElementReference, BaseElementPointerSetBaseElementBaseElementRefUri, hl)
	}
	setBaseElementTargetBaseElementPointerReference := core.GetChildBaseElementReferenceWithUri(setBaseElement, BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri, hl)
	if setBaseElementTargetBaseElementPointerReference == nil {
		setBaseElementTargetBaseElementPointerReference = uOfD.NewBaseElementReference(hl)
		core.SetName(setBaseElementTargetBaseElementPointerReference, "ModifiedBaseElementPointerRef", hl)
		core.SetOwningElement(setBaseElementTargetBaseElementPointerReference, setBaseElement, hl)
		core.SetUri(setBaseElementTargetBaseElementPointerReference, BaseElementPointerSetBaseElementModifiedBaseElementPointerRefUri, hl)
	}
}

func baseElementPointerFunctionsInit() {
	core.GetCore().AddFunction(BaseElementPointerCreateUri, createBaseElementPointer)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementUri, getBaseElement)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementIdUri, getBaseElementId)
	core.GetCore().AddFunction(BaseElementPointerGetBaseElementVersionUri, getBaseElementVersion)
	core.GetCore().AddFunction(BaseElementPointerSetBaseElementUri, setBaseElement)
}
