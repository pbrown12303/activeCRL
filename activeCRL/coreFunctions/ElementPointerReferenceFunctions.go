package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"sync"
)

var ElementPointerReferenceFunctionsUri string = CoreFunctionsPrefix + "ElementPointerReferenceFunctions"

var ElementPointerReferenceCreateUri string = CoreFunctionsPrefix + "ElementPointerReference/Create"
var ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri = CoreFunctionsPrefix + "ElementPointerReference/Create/CreatedElementPointerReferenceRef"

var ElementPointerReferenceGetReferencedElementPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer"
var ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer/SourceElementPointerReferenceRef"
var ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetReferencedElementPointer/IndicatedElementPointerRef"

var ElementPointerReferenceGetElementPointerPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer"
var ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer/SourceElementPointerReferenceRef"
var ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/GetElementPointerPointer/IndicatedElementPointerPointerRef"

var ElementPointerReferenceSetReferencedElementPointerUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer"
var ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer/SourceElementPointerRef"
var ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri string = CoreFunctionsPrefix + "ElementPointerReference/SetReferencedElementPointer/ModifiedElementPointerReferenceRef"

func createElementPointerReference(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementPointerReferenceReference := core.GetChildElementReferenceWithAncestorUri(element, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri, hl)
	if createdElementPointerReferenceReference == nil {
		createdElementPointerReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementPointerReferenceReference, element, hl)
		core.SetName(createdElementPointerReferenceReference, "CreatedElementPointerReferenceReference", hl)
		rootCreatedElementPointerReferenceReference := uOfD.GetElementReferenceWithUri(ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementPointerReferenceReference, hl)
		refinement.SetRefinedElement(createdElementPointerReferenceReference, hl)
		refinement.SetAbstractElement(rootCreatedElementPointerReferenceReference, hl)
	}
	createdElementPointerReference := createdElementPointerReferenceReference.GetReferencedElement(hl)
	if createdElementPointerReference == nil {
		createdElementPointerReference = uOfD.NewElementPointerReference(hl)
		createdElementPointerReferenceReference.SetReferencedElement(createdElementPointerReference, hl)
	}
}

func getReferencedElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerReferenceGetReferencedElementPointerUri)
	if original == nil {
		log.Printf("In GetReferencedElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetReferencedElementPointer, the SourceReference was not found in the replicate")
		return
	}

	targetElementPointerReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri, hl)
	if targetElementPointerReference == nil {
		log.Printf("In GetReferencedElementPointer, the TargetElementPointerReference was not found in the replicate")
		return
	}

	targetElementPointer := targetElementPointerReference.GetReferencedElementPointer(hl)
	untypedElementPointerReference := sourceReference.GetReferencedElement(hl)
	var sourceElementPointerReference core.ElementPointerReference
	var sourceElementPointer core.ElementPointer
	if untypedElementPointerReference != nil {
		switch untypedElementPointerReference.(type) {
		case core.ElementPointerReference:
			sourceElementPointerReference = untypedElementPointerReference.(core.ElementPointerReference)
			sourceElementPointer = sourceElementPointerReference.GetReferencedElementPointer(hl)
		}
	}
	if sourceElementPointer != targetElementPointer {
		targetElementPointerReference.SetReferencedElementPointer(sourceElementPointer, hl)
	}
}

func getElementPointerPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In getElementPointerPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerReferenceGetElementPointerPointerUri)
	if original == nil {
		log.Printf("In GetElementPointerPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetElementPointerPointer, the SourceElementPointerReferenceRef was not found in the replicate")
		return
	}

	indicatedElementPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri, hl)
	if indicatedElementPointerPointerRef == nil {
		log.Printf("In GetElementPointerPointer, the IndicatedElementPointerPointerRef was not found in the replicate")
		return
	}

	indicatedElementPointerPointer := indicatedElementPointerPointerRef.GetReferencedBaseElement(hl)
	untypedSourceElementPointerReference := sourceReference.GetReferencedElement(hl)
	var sourceElementPointerPointer core.ElementPointerPointer
	if untypedSourceElementPointerReference != nil {
		switch untypedSourceElementPointerReference.(type) {
		case core.ElementPointerReference:
			sourceElementPointerReference := untypedSourceElementPointerReference.(core.ElementPointerReference)
			sourceElementPointerPointer = sourceElementPointerReference.GetElementPointerPointer(hl)
		default:
			log.Printf("In GetElementPointerPointer, the SourceElement is not a ElementPointerReference")
		}
	}
	if sourceElementPointerPointer != indicatedElementPointerPointer {
		indicatedElementPointerPointerRef.SetReferencedBaseElement(sourceElementPointerPointer, hl)
	}
}

func setReferencedElementPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementPointerReferenceSetReferencedElementPointerUri)
	if original == nil {
		log.Printf("In SetReferencedElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	baseElementReference := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri, hl)
	if baseElementReference == nil {
		log.Printf("In SetReferencedElementPointer, the ElementPointerReference was not found in the replicate")
		return
	}

	targetElementPointerReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri, hl)
	if targetElementPointerReference == nil {
		log.Printf("In SetReferencedElementPointer, the TargetElementPointer was not found in the replicate")
		return
	}

	sourcedElementPointer := baseElementReference.GetReferencedElementPointer(hl)
	untypedTargetedElement := targetElementPointerReference.GetReferencedElement(hl)
	var targetedElementPointer core.ElementPointer
	var targetedElementPointerReference core.ElementPointerReference
	if untypedTargetedElement != nil {
		switch untypedTargetedElement.(type) {
		case core.ElementPointerReference:
			targetedElementPointerReference = untypedTargetedElement.(core.ElementPointerReference)
			targetedElementPointer = targetedElementPointerReference.GetReferencedElementPointer(hl)
		default:
			log.Printf("In SetReferencedElementPointer, the TargetedElementPointerReference is not a ElementPointerReference")
		}
	}
	if sourcedElementPointer != targetedElementPointer {
		targetedElementPointerReference.SetReferencedElementPointer(sourcedElementPointer, hl)
	}
}

func UpdateRecoveredCoreElementPointerReferenceFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// ElementPointerReferenceFunctions
	elementPointerReferenceFunctions := uOfD.GetElementWithUri(ElementPointerReferenceFunctionsUri)
	if elementPointerReferenceFunctions == nil {
		elementPointerReferenceFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(elementPointerReferenceFunctions, coreFunctionsElement, hl)
		core.SetName(elementPointerReferenceFunctions, "ElementPointerReferenceFunctions", hl)
		core.SetUri(elementPointerReferenceFunctions, ElementPointerReferenceFunctionsUri, hl)
	}

	// CreateElementPointerReference
	createElementPointerReference := uOfD.GetElementWithUri(ElementPointerReferenceCreateUri)
	if createElementPointerReference == nil {
		createElementPointerReference = uOfD.NewElement(hl)
		core.SetOwningElement(createElementPointerReference, elementPointerReferenceFunctions, hl)
		core.SetName(createElementPointerReference, "CreateElementPointerReference", hl)
		core.SetUri(createElementPointerReference, ElementPointerReferenceCreateUri, hl)
	}
	// CreatedElementPointerReference
	createdElementPointerReferenceReference := core.GetChildElementReferenceWithUri(createElementPointerReference, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri, hl)
	if createdElementPointerReferenceReference == nil {
		createdElementPointerReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementPointerReferenceReference, createElementPointerReference, hl)
		core.SetName(createdElementPointerReferenceReference, "CreatedElementPointerReferenceRef", hl)
		core.SetUri(createdElementPointerReferenceReference, ElementPointerReferenceCreateCreatedElementPointerReferenceRefUri, hl)
	}

	// GetReferencedElementPointer
	getReferencedElementPointer := uOfD.GetElementWithUri(ElementPointerReferenceGetReferencedElementPointerUri)
	if getReferencedElementPointer == nil {
		getReferencedElementPointer = uOfD.NewElement(hl)
		core.SetName(getReferencedElementPointer, "GetReferencedElementPointer", hl)
		core.SetOwningElement(getReferencedElementPointer, elementPointerReferenceFunctions, hl)
		core.SetUri(getReferencedElementPointer, ElementPointerReferenceGetReferencedElementPointerUri, hl)
	}
	// GetReferencedElementPointer.SourceReference
	getElementPointerSourceReference := core.GetChildElementReferenceWithUri(getReferencedElementPointer, ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri, hl)
	if getElementPointerSourceReference == nil {
		getElementPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementPointerSourceReference, getReferencedElementPointer, hl)
		core.SetName(getElementPointerSourceReference, "SourceElementPointerReferenceRef", hl)
		core.SetUri(getElementPointerSourceReference, ElementPointerReferenceGetReferencedElementPointerSourceElementPointerReferenceRefUri, hl)
	}
	// GetReferencedElementPointerTargetElementPointerReference
	getElementPointerTargetReference := core.GetChildElementPointerReferenceWithUri(getReferencedElementPointer, ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri, hl)
	if getElementPointerTargetReference == nil {
		getElementPointerTargetReference = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(getElementPointerTargetReference, getReferencedElementPointer, hl)
		core.SetName(getElementPointerTargetReference, "IndicatedElementPointerRef", hl)
		core.SetUri(getElementPointerTargetReference, ElementPointerReferenceGetReferencedElementPointerIndicatedElementPointerRefUri, hl)
	}

	// GetElementPointerPointer
	getElementPointerPointer := uOfD.GetElementWithUri(ElementPointerReferenceGetElementPointerPointerUri)
	if getElementPointerPointer == nil {
		getElementPointerPointer = uOfD.NewElement(hl)
		core.SetName(getElementPointerPointer, "GetElementPointerPointer", hl)
		core.SetOwningElement(getElementPointerPointer, elementPointerReferenceFunctions, hl)
		core.SetUri(getElementPointerPointer, ElementPointerReferenceGetElementPointerPointerUri, hl)
	}
	// GetElementPointerPointer.SourceReference
	getElementPointerPointerSourceReference := core.GetChildElementReferenceWithUri(getReferencedElementPointer, ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri, hl)
	if getElementPointerPointerSourceReference == nil {
		getElementPointerPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementPointerPointerSourceReference, getElementPointerPointer, hl)
		core.SetName(getElementPointerPointerSourceReference, "SourceElementPointerReferenceRef", hl)
		core.SetUri(getElementPointerPointerSourceReference, ElementPointerReferenceGetElementPointerPointerSourceElementPointerReferenceRefUri, hl)
	}
	// GetReferencedElementPointerTargetElementPointerReference
	getElementPointerPointerIndicatedElementPointerPointerRef := core.GetChildBaseElementReferenceWithUri(getReferencedElementPointer, ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri, hl)
	if getElementPointerPointerIndicatedElementPointerPointerRef == nil {
		getElementPointerPointerIndicatedElementPointerPointerRef = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getElementPointerPointerIndicatedElementPointerPointerRef, getElementPointerPointer, hl)
		core.SetName(getElementPointerPointerIndicatedElementPointerPointerRef, "IndicatedElementPointerPointerRef", hl)
		core.SetUri(getElementPointerPointerIndicatedElementPointerPointerRef, ElementPointerReferenceGetElementPointerPointerIndicatedElementPointerPointerRefUri, hl)
	}

	// SetReferencedElementPointer
	setReferencedElementPointer := uOfD.GetElementWithUri(ElementPointerReferenceSetReferencedElementPointerUri)
	if setReferencedElementPointer == nil {
		setReferencedElementPointer = uOfD.NewElement(hl)
		core.SetName(setReferencedElementPointer, "SetReferencedElementPointer", hl)
		core.SetOwningElement(setReferencedElementPointer, elementPointerReferenceFunctions, hl)
		core.SetUri(setReferencedElementPointer, ElementPointerReferenceSetReferencedElementPointerUri, hl)
	}
	// SetReferencedElementPointer.ElementPointerReference
	setReferencedElementPointerElementPointerReference := core.GetChildElementPointerReferenceWithUri(setReferencedElementPointer, ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri, hl)
	if setReferencedElementPointerElementPointerReference == nil {
		setReferencedElementPointerElementPointerReference = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(setReferencedElementPointerElementPointerReference, setReferencedElementPointer, hl)
		core.SetName(setReferencedElementPointerElementPointerReference, "SourceElementPointerRef", hl)
		core.SetUri(setReferencedElementPointerElementPointerReference, ElementPointerReferenceSetReferencedElementPointerSourceElementPointerRefUri, hl)
	}
	// SetReferencedElementPointerTargetElementPointerReference
	setReferencedElementPointerTargetElementPointerReference := core.GetChildElementReferenceWithUri(setReferencedElementPointer, ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri, hl)
	if setReferencedElementPointerTargetElementPointerReference == nil {
		setReferencedElementPointerTargetElementPointerReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedElementPointerTargetElementPointerReference, setReferencedElementPointer, hl)
		core.SetName(setReferencedElementPointerTargetElementPointerReference, "ModifiedElementPointerReferenceRef", hl)
		core.SetUri(setReferencedElementPointerTargetElementPointerReference, ElementPointerReferenceSetReferencedElementPointerModifiedElementPointerReferenceRefUri, hl)
	}
}

func elementPointerReferenceFunctionsInit() {
	core.GetCore().AddFunction(ElementPointerReferenceCreateUri, createElementPointerReference)
	core.GetCore().AddFunction(ElementPointerReferenceGetReferencedElementPointerUri, getReferencedElementPointer)
	core.GetCore().AddFunction(ElementPointerReferenceGetElementPointerPointerUri, getElementPointerPointer)
	core.GetCore().AddFunction(ElementPointerReferenceSetReferencedElementPointerUri, setReferencedElementPointer)
}
