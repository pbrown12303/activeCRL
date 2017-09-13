package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var ElementReferenceFunctionsUri string = CoreFunctionsPrefix + "ElementReferenceFunctions"

var ElementReferenceCreateUri string = CoreFunctionsPrefix + "ElementReference/Create"
var ElementReferenceCreateCreatedElementReferenceRefUri = CoreFunctionsPrefix + "ElementReference/Create/CreatedElementReferenceRef"

var ElementReferenceGetReferencedElementUri string = CoreFunctionsPrefix + "ElementReference/GetReferencedElement"
var ElementReferenceGetReferencedElementSourceElementReferenceRefUri = CoreFunctionsPrefix + "ElementReference/GetReferencedElement/SourceElementReferenceRef"
var ElementReferenceGetReferencedElementIndicatedElementRefUri string = CoreFunctionsPrefix + "ElementReference/GetReferencedElement/IndicatedElementRef"

var ElementReferenceGetElementPointerUri string = CoreFunctionsPrefix + "ElementReference/GetElementPointer"
var ElementReferenceGetElementPointerSourceElementReferenceRefUri string = CoreFunctionsPrefix + "ElementReference/GetElementPointer/SourceElementReferenceRef"
var ElementReferenceGetElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "ElementReference/GetElementPointer/IndicatedElementPointerRef"

var ElementReferenceSetReferencedElementUri string = CoreFunctionsPrefix + "ElementReference/SetReferencedElement"
var ElementReferenceSetReferencedElementSourceElementRefUri string = CoreFunctionsPrefix + "ElementReference/SetReferencedElement/SourceElementRef"
var ElementReferenceSetReferencedElementModifiedElementReferenceRefUri string = CoreFunctionsPrefix + "ElementReference/SetReferencedElement/ModifiedElementReferenceRef"

func elementReferenceCreateElementReference(element core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementReferenceReference := core.GetChildElementReferenceWithAncestorUri(element, ElementReferenceCreateCreatedElementReferenceRefUri, hl)
	if createdElementReferenceReference == nil {
		createdElementReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementReferenceReference, element, hl)
		core.SetName(createdElementReferenceReference, "CreatedElementReferenceReference", hl)
		rootCreatedElementReferenceReference := uOfD.GetElementReferenceWithUri(ElementReferenceCreateCreatedElementReferenceRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementReferenceReference, hl)
		refinement.SetRefinedElement(createdElementReferenceReference, hl)
		refinement.SetAbstractElement(rootCreatedElementReferenceReference, hl)
	}
	createdElementReference := createdElementReferenceReference.GetReferencedElement(hl)
	if createdElementReference == nil {
		createdElementReference = uOfD.NewElementReference(hl)
		createdElementReferenceReference.SetReferencedElement(createdElementReference, hl)
	}
}

func elementReferenceGetReferencedElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementReferenceGetReferencedElementUri)
	if original == nil {
		log.Printf("In GetReferencedElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceGetReferencedElementSourceElementReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetReferencedElement, the SourceReference was not found in the replicate")
		return
	}

	targetElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceGetReferencedElementIndicatedElementRefUri, hl)
	if targetElementReference == nil {
		log.Printf("In GetReferencedElement, the TargetElementReference was not found in the replicate")
		return
	}

	targetElement := targetElementReference.GetReferencedElement(hl)
	untypedElementReference := sourceReference.GetReferencedElement(hl)
	var sourceElementReference core.ElementReference
	var sourceElement core.Element
	if untypedElementReference != nil {
		switch untypedElementReference.(type) {
		case core.ElementReference:
			sourceElementReference = untypedElementReference.(core.ElementReference)
			sourceElement = sourceElementReference.GetReferencedElement(hl)
		}
	}
	if sourceElement != targetElement {
		targetElementReference.SetReferencedElement(sourceElement, hl)
	}
}

func elementReferenceGetElementPointer(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementReferenceGetElementPointerUri)
	if original == nil {
		log.Printf("In GetElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceGetElementPointerSourceElementReferenceRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetElementPointer, the SourceElementReferenceRef was not found in the replicate")
		return
	}

	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, ElementReferenceGetElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		log.Printf("In GetElementPointer, the IndicatedElementPointerRef was not found in the replicate")
		return
	}

	indicatedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	untypedSourceElementReference := sourceReference.GetReferencedElement(hl)
	var sourceElementPointer core.ElementPointer
	if untypedSourceElementReference != nil {
		switch untypedSourceElementReference.(type) {
		case core.ElementReference:
			sourceElementReference := untypedSourceElementReference.(core.ElementReference)
			sourceElementPointer = sourceElementReference.GetElementPointer(hl)
		default:
			log.Printf("In GetElementPointer, the SourceElement is not a ElementReference")
		}
	}
	if sourceElementPointer != indicatedElementPointer {
		indicatedElementPointerRef.SetReferencedElementPointer(sourceElementPointer, hl)
	}
}

func elementReferenceSetReferencedElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementReferenceSetReferencedElementUri)
	if original == nil {
		log.Printf("In SetReferencedElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	baseElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceSetReferencedElementSourceElementRefUri, hl)
	if baseElementReference == nil {
		log.Printf("In SetReferencedElement, the ElementReference was not found in the replicate")
		return
	}

	targetElementReference := core.GetChildElementReferenceWithAncestorUri(replicate, ElementReferenceSetReferencedElementModifiedElementReferenceRefUri, hl)
	if targetElementReference == nil {
		log.Printf("In SetReferencedElement, the TargetElement was not found in the replicate")
		return
	}

	sourcedElement := baseElementReference.GetReferencedElement(hl)
	untypedTargetedElement := targetElementReference.GetReferencedElement(hl)
	var targetedElement core.Element
	var targetedElementReference core.ElementReference
	if untypedTargetedElement != nil {
		switch untypedTargetedElement.(type) {
		case core.ElementReference:
			targetedElementReference = untypedTargetedElement.(core.ElementReference)
			targetedElement = targetedElementReference.GetReferencedElement(hl)
		default:
			log.Printf("In SetReferencedElement, the TargetedElementReference is not a ElementReference")
		}
	}
	if sourcedElement != targetedElement {
		targetedElementReference.SetReferencedElement(sourcedElement, hl)
	}
}

func UpdateRecoveredCoreElementReferenceFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// ElementReferenceFunctions
	elementReferenceFunctions := uOfD.GetElementWithUri(ElementReferenceFunctionsUri)
	if elementReferenceFunctions == nil {
		elementReferenceFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(elementReferenceFunctions, coreFunctionsElement, hl)
		core.SetName(elementReferenceFunctions, "ElementReferenceFunctions", hl)
		core.SetUri(elementReferenceFunctions, ElementReferenceFunctionsUri, hl)
	}

	// CreateElementReference
	elementReferenceCreateElementReference := uOfD.GetElementWithUri(ElementReferenceCreateUri)
	if elementReferenceCreateElementReference == nil {
		elementReferenceCreateElementReference = uOfD.NewElement(hl)
		core.SetOwningElement(elementReferenceCreateElementReference, elementReferenceFunctions, hl)
		core.SetName(elementReferenceCreateElementReference, "CreateElementReference", hl)
		core.SetUri(elementReferenceCreateElementReference, ElementReferenceCreateUri, hl)
	}
	// CreatedElementReference
	createdElementReferenceReference := core.GetChildElementReferenceWithUri(elementReferenceCreateElementReference, ElementReferenceCreateCreatedElementReferenceRefUri, hl)
	if createdElementReferenceReference == nil {
		createdElementReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementReferenceReference, elementReferenceCreateElementReference, hl)
		core.SetName(createdElementReferenceReference, "CreatedElementReferenceRef", hl)
		core.SetUri(createdElementReferenceReference, ElementReferenceCreateCreatedElementReferenceRefUri, hl)
	}

	// GetReferencedElement
	elementReferenceGetReferencedElement := uOfD.GetElementWithUri(ElementReferenceGetReferencedElementUri)
	if elementReferenceGetReferencedElement == nil {
		elementReferenceGetReferencedElement = uOfD.NewElement(hl)
		core.SetName(elementReferenceGetReferencedElement, "GetReferencedElement", hl)
		core.SetOwningElement(elementReferenceGetReferencedElement, elementReferenceFunctions, hl)
		core.SetUri(elementReferenceGetReferencedElement, ElementReferenceGetReferencedElementUri, hl)
	}
	// GetReferencedElement.SourceReference
	getElementSourceReference := core.GetChildElementReferenceWithUri(elementReferenceGetReferencedElement, ElementReferenceGetReferencedElementSourceElementReferenceRefUri, hl)
	if getElementSourceReference == nil {
		getElementSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementSourceReference, elementReferenceGetReferencedElement, hl)
		core.SetName(getElementSourceReference, "SourceElementReferenceRef", hl)
		core.SetUri(getElementSourceReference, ElementReferenceGetReferencedElementSourceElementReferenceRefUri, hl)
	}
	// GetReferencedElementTargetElementReference
	getElementTargetReference := core.GetChildElementReferenceWithUri(elementReferenceGetReferencedElement, ElementReferenceGetReferencedElementIndicatedElementRefUri, hl)
	if getElementTargetReference == nil {
		getElementTargetReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementTargetReference, elementReferenceGetReferencedElement, hl)
		core.SetName(getElementTargetReference, "IndicatedElementRef", hl)
		core.SetUri(getElementTargetReference, ElementReferenceGetReferencedElementIndicatedElementRefUri, hl)
	}

	// GetElementPointer
	elementReferenceGetElementPointer := uOfD.GetElementWithUri(ElementReferenceGetElementPointerUri)
	if elementReferenceGetElementPointer == nil {
		elementReferenceGetElementPointer = uOfD.NewElement(hl)
		core.SetName(elementReferenceGetElementPointer, "GetElementPointer", hl)
		core.SetOwningElement(elementReferenceGetElementPointer, elementReferenceFunctions, hl)
		core.SetUri(elementReferenceGetElementPointer, ElementReferenceGetElementPointerUri, hl)
	}
	// GetElementPointer.SourceReference
	getElementPointerSourceReference := core.GetChildElementReferenceWithUri(elementReferenceGetElementPointer, ElementReferenceGetElementPointerSourceElementReferenceRefUri, hl)
	if getElementPointerSourceReference == nil {
		getElementPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementPointerSourceReference, elementReferenceGetElementPointer, hl)
		core.SetName(getElementPointerSourceReference, "SourceElementReferenceRef", hl)
		core.SetUri(getElementPointerSourceReference, ElementReferenceGetElementPointerSourceElementReferenceRefUri, hl)
	}
	// GetElementPointerIndicatedElementPointerRef
	getElementPointerIndicatedElementPointerRef := core.GetChildElementPointerReferenceWithUri(elementReferenceGetElementPointer, ElementReferenceGetElementPointerIndicatedElementPointerRefUri, hl)
	if getElementPointerIndicatedElementPointerRef == nil {
		getElementPointerIndicatedElementPointerRef = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(getElementPointerIndicatedElementPointerRef, elementReferenceGetElementPointer, hl)
		core.SetName(getElementPointerIndicatedElementPointerRef, "IndicatedElementPointerRef", hl)
		core.SetUri(getElementPointerIndicatedElementPointerRef, ElementReferenceGetElementPointerIndicatedElementPointerRefUri, hl)
	}

	// SetReferencedElement
	elementReferenceSetReferencedElement := uOfD.GetElementWithUri(ElementReferenceSetReferencedElementUri)
	if elementReferenceSetReferencedElement == nil {
		elementReferenceSetReferencedElement = uOfD.NewElement(hl)
		core.SetName(elementReferenceSetReferencedElement, "SetReferencedElement", hl)
		core.SetOwningElement(elementReferenceSetReferencedElement, elementReferenceFunctions, hl)
		core.SetUri(elementReferenceSetReferencedElement, ElementReferenceSetReferencedElementUri, hl)
	}
	// SetReferencedElement.ElementReference
	setReferencedElementElementReference := core.GetChildElementReferenceWithUri(elementReferenceSetReferencedElement, ElementReferenceSetReferencedElementSourceElementRefUri, hl)
	if setReferencedElementElementReference == nil {
		setReferencedElementElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedElementElementReference, elementReferenceSetReferencedElement, hl)
		core.SetName(setReferencedElementElementReference, "SourceElementRef", hl)
		core.SetUri(setReferencedElementElementReference, ElementReferenceSetReferencedElementSourceElementRefUri, hl)
	}
	// SetReferencedElementTargetElementReference
	setReferencedElementTargetElementReference := core.GetChildElementReferenceWithUri(elementReferenceSetReferencedElement, ElementReferenceSetReferencedElementModifiedElementReferenceRefUri, hl)
	if setReferencedElementTargetElementReference == nil {
		setReferencedElementTargetElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedElementTargetElementReference, elementReferenceSetReferencedElement, hl)
		core.SetName(setReferencedElementTargetElementReference, "ModifiedElementReferenceRef", hl)
		core.SetUri(setReferencedElementTargetElementReference, ElementReferenceSetReferencedElementModifiedElementReferenceRefUri, hl)
	}
}

func elementReferenceFunctionsInit() {
	core.GetCore().AddFunction(ElementReferenceCreateUri, elementReferenceCreateElementReference)
	core.GetCore().AddFunction(ElementReferenceGetReferencedElementUri, elementReferenceGetReferencedElement)
	core.GetCore().AddFunction(ElementReferenceGetElementPointerUri, elementReferenceGetElementPointer)
	core.GetCore().AddFunction(ElementReferenceSetReferencedElementUri, elementReferenceSetReferencedElement)
}
