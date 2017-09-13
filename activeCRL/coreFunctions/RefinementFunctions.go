package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var RefinementFunctionsUri string = CoreFunctionsPrefix + "RefinementFunctions"

var RefinementCreateUri string = CoreFunctionsPrefix + "Refinement/Create"
var RefinementCreateCreatedRefinementRefUri = CoreFunctionsPrefix + "Refinement/Create/CreatedRefinementRef"

var RefinementGetAbstractElementUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElement"
var RefinementGetAbstractElementSourceRefinementRefUri = CoreFunctionsPrefix + "Refinement/GetAbstractElement/SourceRefinementRef"
var RefinementGetAbstractElementIndicatedElementRefUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElement/IndicatedElementRef"

var RefinementGetAbstractElementPointerUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElementPointer"
var RefinementGetAbstractElementPointerSourceRefinementRefUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElementPointer/SourceRefinementRef"
var RefinementGetAbstractElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "Refinement/GetAbstractElementPointer/IndicatedElementPointerRef"

var RefinementGetRefinedElementUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElement"
var RefinementGetRefinedElementSourceRefinementRefUri = CoreFunctionsPrefix + "Refinement/GetRefinedElement/SourceRefinementRef"
var RefinementGetRefinedElementIndicatedElementRefUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElement/IndicatedElementRef"

var RefinementGetRefinedElementPointerUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElementPointer"
var RefinementGetRefinedElementPointerSourceRefinementRefUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElementPointer/SourceRefinementRef"
var RefinementGetRefinedElementPointerIndicatedElementPointerRefUri string = CoreFunctionsPrefix + "Refinement/GetRefinedElementPointer/IndicatedElementPointerRef"

var RefinementSetAbstractElementUri string = CoreFunctionsPrefix + "Refinement/SetAbstractElement"
var RefinementSetAbstractElementSourceElementRefUri string = CoreFunctionsPrefix + "Refinement/SetAbstractElement/SourceElementRef"
var RefinementSetAbstractElementModifiedRefinementRefUri string = CoreFunctionsPrefix + "Refinement/SetAbstractElement/ModifiedRefinementRef"

var RefinementSetRefinedElementUri string = CoreFunctionsPrefix + "Refinement/SetRefinedElement"
var RefinementSetRefinedElementSourceElementRefUri string = CoreFunctionsPrefix + "Refinement/SetRefinedElement/SourceElementRef"
var RefinementSetRefinedElementModifiedRefinementRefUri string = CoreFunctionsPrefix + "Refinement/SetRefinedElement/ModifiedRefinementRef"

func refinementCreateRefinement(element core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdRefinementRef := core.GetChildElementReferenceWithAncestorUri(element, RefinementCreateCreatedRefinementRefUri, hl)
	if createdRefinementRef == nil {
		createdRefinementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdRefinementRef, element, hl)
		core.SetName(createdRefinementRef, "CreatedRefinementReference", hl)
		rootCreatedRefinementReference := uOfD.GetElementReferenceWithUri(RefinementCreateCreatedRefinementRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdRefinementRef, hl)
		refinement.SetRefinedElement(createdRefinementRef, hl)
		refinement.SetAbstractElement(rootCreatedRefinementReference, hl)
	}
	createdRefinement := createdRefinementRef.GetReferencedElement(hl)
	if createdRefinement == nil {
		createdRefinement = uOfD.NewRefinement(hl)
		createdRefinementRef.SetReferencedElement(createdRefinement, hl)
	}
}

func refinementGetAbstractElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementGetAbstractElementUri)
	if original == nil {
		log.Printf("In GetAbstractElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		log.Printf("In GetAbstractElement, the SourceRefinementRef was not found in the replicate")
		return
	}

	indicatedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementIndicatedElementRefUri, hl)
	if indicatedElementRef == nil {
		log.Printf("In GetAbstractElement, the IndicatedElementReference was not found in the replicate")
		return
	}

	targetElement := indicatedElementRef.GetReferencedElement(hl)
	untypedRefinement := sourceRefinementRef.GetReferencedElement(hl)
	var sourceRefinement core.Refinement
	var sourceElement core.Element
	if untypedRefinement != nil {
		switch untypedRefinement.(type) {
		case core.Refinement:
			sourceRefinement = untypedRefinement.(core.Refinement)
			sourceElement = sourceRefinement.GetAbstractElement(hl)
		}
	}
	if sourceElement != targetElement {
		indicatedElementRef.SetReferencedElement(sourceElement, hl)
	}
}

func refinementGetAbstractElementPointer(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementGetAbstractElementPointerUri)
	if original == nil {
		log.Printf("In GetAbstractElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetAbstractElementPointerSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		log.Printf("In GetAbstractElementPointer, the SourceRefinementRef was not found in the replicate")
		return
	}

	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, RefinementGetAbstractElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		log.Printf("In GetAbstractElementPointer, the IndicatedElementPointerRef was not found in the replicate")
		return
	}

	indicatedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	untypedSourceRefinement := sourceRefinementRef.GetReferencedElement(hl)
	var sourceElementPointer core.ElementPointer
	if untypedSourceRefinement != nil {
		switch untypedSourceRefinement.(type) {
		case core.Refinement:
			sourceRefinement := untypedSourceRefinement.(core.Refinement)
			sourceElementPointer = sourceRefinement.GetAbstractElementPointer(hl)
		default:
			log.Printf("In GetAbstractElementPointer, the SourceElement is not a Refinement")
		}
	}
	if sourceElementPointer != indicatedElementPointer {
		indicatedElementPointerRef.SetReferencedElementPointer(sourceElementPointer, hl)
	}
}

func refinementGetRefinedElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementGetRefinedElementUri)
	if original == nil {
		log.Printf("In GetRefinedElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		log.Printf("In GetRefinedElement, the SourceRefinementRef was not found in the replicate")
		return
	}

	indicatedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementIndicatedElementRefUri, hl)
	if indicatedElementRef == nil {
		log.Printf("In GetRefinedElement, the IndicatedElementReference was not found in the replicate")
		return
	}

	targetElement := indicatedElementRef.GetReferencedElement(hl)
	untypedRefinement := sourceRefinementRef.GetReferencedElement(hl)
	var sourceRefinement core.Refinement
	var sourceElement core.Element
	if untypedRefinement != nil {
		switch untypedRefinement.(type) {
		case core.Refinement:
			sourceRefinement = untypedRefinement.(core.Refinement)
			sourceElement = sourceRefinement.GetRefinedElement(hl)
		}
	}
	if sourceElement != targetElement {
		indicatedElementRef.SetReferencedElement(sourceElement, hl)
	}
}

func refinementGetRefinedElementPointer(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementGetRefinedElementPointerUri)
	if original == nil {
		log.Printf("In GetRefinedElementPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementGetRefinedElementPointerSourceRefinementRefUri, hl)
	if sourceRefinementRef == nil {
		log.Printf("In GetRefinedElementPointer, the SourceRefinementRef was not found in the replicate")
		return
	}

	indicatedElementPointerRef := core.GetChildElementPointerReferenceWithAncestorUri(replicate, RefinementGetRefinedElementPointerIndicatedElementPointerRefUri, hl)
	if indicatedElementPointerRef == nil {
		log.Printf("In GetRefinedElementPointer, the IndicatedElementPointerRef was not found in the replicate")
		return
	}

	indicatedElementPointer := indicatedElementPointerRef.GetReferencedElementPointer(hl)
	untypedSourceRefinement := sourceRefinementRef.GetReferencedElement(hl)
	var sourceElementPointer core.ElementPointer
	if untypedSourceRefinement != nil {
		switch untypedSourceRefinement.(type) {
		case core.Refinement:
			sourceRefinement := untypedSourceRefinement.(core.Refinement)
			sourceElementPointer = sourceRefinement.GetRefinedElementPointer(hl)
		default:
			log.Printf("In GetRefinedElementPointer, the SourceElement is not a Refinement")
		}
	}
	if sourceElementPointer != indicatedElementPointer {
		indicatedElementPointerRef.SetReferencedElementPointer(sourceElementPointer, hl)
	}
}

func refinementSetAbstractElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementSetAbstractElementUri)
	if original == nil {
		log.Printf("In SetAbstractElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetAbstractElementSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In SetAbstractElement, the ElementReference was not found in the replicate")
		return
	}

	targetRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetAbstractElementModifiedRefinementRefUri, hl)
	if targetRefinementRef == nil {
		log.Printf("In SetAbstractElement, the TargetRefinementRef was not found in the replicate")
		return
	}

	sourcedElement := sourceElementRef.GetReferencedElement(hl)
	untypedModifiedRefinement := targetRefinementRef.GetReferencedElement(hl)
	var targetElement core.Element
	var modifiedRefinement core.Refinement
	if untypedModifiedRefinement != nil {
		switch untypedModifiedRefinement.(type) {
		case core.Refinement:
			modifiedRefinement = untypedModifiedRefinement.(core.Refinement)
			targetElement = modifiedRefinement.GetAbstractElement(hl)
		default:
			log.Printf("In SetAbstractElement, the TargetedRefinement is not a Refinement")
		}
	}
	if sourcedElement != targetElement {
		modifiedRefinement.SetAbstractElement(sourcedElement, hl)
	}
}

func refinementSetRefinedElement(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(RefinementSetRefinedElementUri)
	if original == nil {
		log.Printf("In SetRefinedElement the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetRefinedElementSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In SetRefinedElement, the ElementReference was not found in the replicate")
		return
	}

	targetRefinementRef := core.GetChildElementReferenceWithAncestorUri(replicate, RefinementSetRefinedElementModifiedRefinementRefUri, hl)
	if targetRefinementRef == nil {
		log.Printf("In SetRefinedElement, the TargetRefinementRef was not found in the replicate")
		return
	}

	sourcedElement := sourceElementRef.GetReferencedElement(hl)
	untypedModifiedRefinement := targetRefinementRef.GetReferencedElement(hl)
	var targetElement core.Element
	var modifiedRefinement core.Refinement
	if untypedModifiedRefinement != nil {
		switch untypedModifiedRefinement.(type) {
		case core.Refinement:
			modifiedRefinement = untypedModifiedRefinement.(core.Refinement)
			targetElement = modifiedRefinement.GetRefinedElement(hl)
		default:
			log.Printf("In SetRefinedElement, the TargetedRefinement is not a Refinement")
		}
	}
	if sourcedElement != targetElement {
		modifiedRefinement.SetRefinedElement(sourcedElement, hl)
	}
}

func UpdateRecoveredCoreRefinementFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// RefinementFunctions
	refinementFunctions := uOfD.GetElementWithUri(RefinementFunctionsUri)
	if refinementFunctions == nil {
		refinementFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(refinementFunctions, coreFunctionsElement, hl)
		core.SetName(refinementFunctions, "RefinementFunctions", hl)
		core.SetUri(refinementFunctions, RefinementFunctionsUri, hl)
	}

	// CreateRefinement
	refinementCreateRefinement := uOfD.GetElementWithUri(RefinementCreateUri)
	if refinementCreateRefinement == nil {
		refinementCreateRefinement = uOfD.NewElement(hl)
		core.SetOwningElement(refinementCreateRefinement, refinementFunctions, hl)
		core.SetName(refinementCreateRefinement, "CreateRefinement", hl)
		core.SetUri(refinementCreateRefinement, RefinementCreateUri, hl)
	}
	// CreatedRefinement
	createdRefinementRef := core.GetChildElementReferenceWithUri(refinementCreateRefinement, RefinementCreateCreatedRefinementRefUri, hl)
	if createdRefinementRef == nil {
		createdRefinementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdRefinementRef, refinementCreateRefinement, hl)
		core.SetName(createdRefinementRef, "CreatedRefinementRef", hl)
		core.SetUri(createdRefinementRef, RefinementCreateCreatedRefinementRefUri, hl)
	}

	// GetAbstractElement
	refinementGetAbstractElement := uOfD.GetElementWithUri(RefinementGetAbstractElementUri)
	if refinementGetAbstractElement == nil {
		refinementGetAbstractElement = uOfD.NewElement(hl)
		core.SetName(refinementGetAbstractElement, "GetAbstractElement", hl)
		core.SetOwningElement(refinementGetAbstractElement, refinementFunctions, hl)
		core.SetUri(refinementGetAbstractElement, RefinementGetAbstractElementUri, hl)
	}
	// GetAbstractElement.SourceReference
	getElementSourceReference := core.GetChildElementReferenceWithUri(refinementGetAbstractElement, RefinementGetAbstractElementSourceRefinementRefUri, hl)
	if getElementSourceReference == nil {
		getElementSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementSourceReference, refinementGetAbstractElement, hl)
		core.SetName(getElementSourceReference, "SourceRefinementRef", hl)
		core.SetUri(getElementSourceReference, RefinementGetAbstractElementSourceRefinementRefUri, hl)
	}
	// GetAbstractElementTargetElementReference
	getElementTargetReference := core.GetChildElementReferenceWithUri(refinementGetAbstractElement, RefinementGetAbstractElementIndicatedElementRefUri, hl)
	if getElementTargetReference == nil {
		getElementTargetReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementTargetReference, refinementGetAbstractElement, hl)
		core.SetName(getElementTargetReference, "IndicatedElementRef", hl)
		core.SetUri(getElementTargetReference, RefinementGetAbstractElementIndicatedElementRefUri, hl)
	}

	// GetAbstractElementPointer
	refinementGetAbstractElementPointer := uOfD.GetElementWithUri(RefinementGetAbstractElementPointerUri)
	if refinementGetAbstractElementPointer == nil {
		refinementGetAbstractElementPointer = uOfD.NewElement(hl)
		core.SetName(refinementGetAbstractElementPointer, "GetAbstractElementPointer", hl)
		core.SetOwningElement(refinementGetAbstractElementPointer, refinementFunctions, hl)
		core.SetUri(refinementGetAbstractElementPointer, RefinementGetAbstractElementPointerUri, hl)
	}
	// GetAbstractElementPointer.SourceReference
	getElementPointerSourceReference := core.GetChildElementReferenceWithUri(refinementGetAbstractElementPointer, RefinementGetAbstractElementPointerSourceRefinementRefUri, hl)
	if getElementPointerSourceReference == nil {
		getElementPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementPointerSourceReference, refinementGetAbstractElementPointer, hl)
		core.SetName(getElementPointerSourceReference, "SourceRefinementRef", hl)
		core.SetUri(getElementPointerSourceReference, RefinementGetAbstractElementPointerSourceRefinementRefUri, hl)
	}
	// GetAbstractElementPointerIndicatedElementPointerRef
	getElementPointerIndicatedElementPointerRef := core.GetChildElementPointerReferenceWithUri(refinementGetAbstractElementPointer, RefinementGetAbstractElementPointerIndicatedElementPointerRefUri, hl)
	if getElementPointerIndicatedElementPointerRef == nil {
		getElementPointerIndicatedElementPointerRef = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(getElementPointerIndicatedElementPointerRef, refinementGetAbstractElementPointer, hl)
		core.SetName(getElementPointerIndicatedElementPointerRef, "IndicatedElementPointerRef", hl)
		core.SetUri(getElementPointerIndicatedElementPointerRef, RefinementGetAbstractElementPointerIndicatedElementPointerRefUri, hl)
	}

	// GetRefinedElement
	refinementGetRefinedElement := uOfD.GetElementWithUri(RefinementGetRefinedElementUri)
	if refinementGetRefinedElement == nil {
		refinementGetRefinedElement = uOfD.NewElement(hl)
		core.SetName(refinementGetRefinedElement, "GetRefinedElement", hl)
		core.SetOwningElement(refinementGetRefinedElement, refinementFunctions, hl)
		core.SetUri(refinementGetRefinedElement, RefinementGetRefinedElementUri, hl)
	}
	// GetRefinedElement.SourceReference
	getElementSourceReference = core.GetChildElementReferenceWithUri(refinementGetRefinedElement, RefinementGetRefinedElementSourceRefinementRefUri, hl)
	if getElementSourceReference == nil {
		getElementSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementSourceReference, refinementGetRefinedElement, hl)
		core.SetName(getElementSourceReference, "SourceRefinementRef", hl)
		core.SetUri(getElementSourceReference, RefinementGetRefinedElementSourceRefinementRefUri, hl)
	}
	// GetRefinedElementTargetElementReference
	getElementTargetReference = core.GetChildElementReferenceWithUri(refinementGetRefinedElement, RefinementGetRefinedElementIndicatedElementRefUri, hl)
	if getElementTargetReference == nil {
		getElementTargetReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementTargetReference, refinementGetRefinedElement, hl)
		core.SetName(getElementTargetReference, "IndicatedElementRef", hl)
		core.SetUri(getElementTargetReference, RefinementGetRefinedElementIndicatedElementRefUri, hl)
	}

	// GetRefinedElementPointer
	refinementGetRefinedElementPointer := uOfD.GetElementWithUri(RefinementGetRefinedElementPointerUri)
	if refinementGetRefinedElementPointer == nil {
		refinementGetRefinedElementPointer = uOfD.NewElement(hl)
		core.SetName(refinementGetRefinedElementPointer, "GetRefinedElementPointer", hl)
		core.SetOwningElement(refinementGetRefinedElementPointer, refinementFunctions, hl)
		core.SetUri(refinementGetRefinedElementPointer, RefinementGetRefinedElementPointerUri, hl)
	}
	// GetRefinedElementPointer.SourceReference
	getElementPointerSourceReference = core.GetChildElementReferenceWithUri(refinementGetRefinedElementPointer, RefinementGetRefinedElementPointerSourceRefinementRefUri, hl)
	if getElementPointerSourceReference == nil {
		getElementPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementPointerSourceReference, refinementGetRefinedElementPointer, hl)
		core.SetName(getElementPointerSourceReference, "SourceRefinementRef", hl)
		core.SetUri(getElementPointerSourceReference, RefinementGetRefinedElementPointerSourceRefinementRefUri, hl)
	}
	// GetRefinedElementPointerIndicatedElementPointerRef
	getElementPointerIndicatedElementPointerRef = core.GetChildElementPointerReferenceWithUri(refinementGetRefinedElementPointer, RefinementGetRefinedElementPointerIndicatedElementPointerRefUri, hl)
	if getElementPointerIndicatedElementPointerRef == nil {
		getElementPointerIndicatedElementPointerRef = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(getElementPointerIndicatedElementPointerRef, refinementGetRefinedElementPointer, hl)
		core.SetName(getElementPointerIndicatedElementPointerRef, "IndicatedElementPointerRef", hl)
		core.SetUri(getElementPointerIndicatedElementPointerRef, RefinementGetRefinedElementPointerIndicatedElementPointerRefUri, hl)
	}

	// SetAbstractElement
	refinementSetAbstractElement := uOfD.GetElementWithUri(RefinementSetAbstractElementUri)
	if refinementSetAbstractElement == nil {
		refinementSetAbstractElement = uOfD.NewElement(hl)
		core.SetName(refinementSetAbstractElement, "SetAbstractElement", hl)
		core.SetOwningElement(refinementSetAbstractElement, refinementFunctions, hl)
		core.SetUri(refinementSetAbstractElement, RefinementSetAbstractElementUri, hl)
	}
	// SetAbstractElement.ElementReference
	setReferencedElementElementReference := core.GetChildElementReferenceWithUri(refinementSetAbstractElement, RefinementSetAbstractElementSourceElementRefUri, hl)
	if setReferencedElementElementReference == nil {
		setReferencedElementElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedElementElementReference, refinementSetAbstractElement, hl)
		core.SetName(setReferencedElementElementReference, "SourceElementRef", hl)
		core.SetUri(setReferencedElementElementReference, RefinementSetAbstractElementSourceElementRefUri, hl)
	}
	// SetAbstractElementTargetRefinement
	setReferencedElementTargetRefinement := core.GetChildElementReferenceWithUri(refinementSetAbstractElement, RefinementSetAbstractElementModifiedRefinementRefUri, hl)
	if setReferencedElementTargetRefinement == nil {
		setReferencedElementTargetRefinement = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedElementTargetRefinement, refinementSetAbstractElement, hl)
		core.SetName(setReferencedElementTargetRefinement, "ModifiedRefinementRef", hl)
		core.SetUri(setReferencedElementTargetRefinement, RefinementSetAbstractElementModifiedRefinementRefUri, hl)
	}

	// SetRefinedElement
	refinementSetRefinedElement := uOfD.GetElementWithUri(RefinementSetRefinedElementUri)
	if refinementSetRefinedElement == nil {
		refinementSetRefinedElement = uOfD.NewElement(hl)
		core.SetName(refinementSetRefinedElement, "SetRefinedElement", hl)
		core.SetOwningElement(refinementSetRefinedElement, refinementFunctions, hl)
		core.SetUri(refinementSetRefinedElement, RefinementSetRefinedElementUri, hl)
	}
	// SetRefinedElement.ElementReference
	setReferencedElementElementReference = core.GetChildElementReferenceWithUri(refinementSetRefinedElement, RefinementSetRefinedElementSourceElementRefUri, hl)
	if setReferencedElementElementReference == nil {
		setReferencedElementElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedElementElementReference, refinementSetRefinedElement, hl)
		core.SetName(setReferencedElementElementReference, "SourceElementRef", hl)
		core.SetUri(setReferencedElementElementReference, RefinementSetRefinedElementSourceElementRefUri, hl)
	}
	// SetRefinedElementTargetRefinement
	setReferencedElementTargetRefinement = core.GetChildElementReferenceWithUri(refinementSetRefinedElement, RefinementSetRefinedElementModifiedRefinementRefUri, hl)
	if setReferencedElementTargetRefinement == nil {
		setReferencedElementTargetRefinement = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedElementTargetRefinement, refinementSetRefinedElement, hl)
		core.SetName(setReferencedElementTargetRefinement, "ModifiedRefinementRef", hl)
		core.SetUri(setReferencedElementTargetRefinement, RefinementSetRefinedElementModifiedRefinementRefUri, hl)
	}

}

func refinementFunctionsInit() {
	core.GetCore().AddFunction(RefinementCreateUri, refinementCreateRefinement)
	core.GetCore().AddFunction(RefinementGetAbstractElementUri, refinementGetAbstractElement)
	core.GetCore().AddFunction(RefinementGetAbstractElementPointerUri, refinementGetAbstractElementPointer)
	core.GetCore().AddFunction(RefinementGetRefinedElementUri, refinementGetRefinedElement)
	core.GetCore().AddFunction(RefinementGetRefinedElementPointerUri, refinementGetRefinedElementPointer)
	core.GetCore().AddFunction(RefinementSetAbstractElementUri, refinementSetAbstractElement)
	core.GetCore().AddFunction(RefinementSetRefinedElementUri, refinementSetRefinedElement)
}
