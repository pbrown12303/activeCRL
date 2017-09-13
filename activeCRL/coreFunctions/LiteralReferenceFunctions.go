package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var LiteralReferenceFunctionsUri string = CoreFunctionsPrefix + "LiteralReferenceFunctions"

var LiteralReferenceCreateUri string = CoreFunctionsPrefix + "LiteralReference/Create"
var LiteralReferenceCreateCreatedLiteralReferenceRefUri = CoreFunctionsPrefix + "LiteralReference/Create/CreatedLiteralReferenceRef"

var LiteralReferenceGetReferencedLiteralUri string = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral"
var LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral/SourceLiteralReferenceRef"
var LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralReference/GetReferencedLiteral/IndicatedLiteralRef"

var LiteralReferenceGetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer"
var LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer/SourceLiteralReferenceRef"
var LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralReference/GetLiteralPointer/IndicatedLiteralPointerRef"

var LiteralReferenceSetReferencedLiteralUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral"
var LiteralReferenceSetReferencedLiteralSourceLiteralRefUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral/SourceLiteralRef"
var LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri string = CoreFunctionsPrefix + "LiteralReference/SetReferencedLiteral/ModifiedLiteralReferenceRef"

func literalReferenceCreateLiteralReference(element core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)

	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(element, LiteralReferenceCreateCreatedLiteralReferenceRefUri, hl)
	if createdLiteralReferenceRef == nil {
		createdLiteralReferenceRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdLiteralReferenceRef, element, hl)
		core.SetName(createdLiteralReferenceRef, "CreatedLiteralReferenceReference", hl)
		rootCreatedLiteralReferenceReference := uOfD.GetElementReferenceWithUri(LiteralReferenceCreateCreatedLiteralReferenceRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralReferenceRef, hl)
		refinement.SetRefinedElement(createdLiteralReferenceRef, hl)
		refinement.SetAbstractElement(rootCreatedLiteralReferenceReference, hl)
	}
	createdLiteralReference := createdLiteralReferenceRef.GetReferencedElement(hl)
	if createdLiteralReference == nil {
		createdLiteralReference = uOfD.NewLiteralReference(hl)
		createdLiteralReferenceRef.SetReferencedElement(createdLiteralReference, hl)
	}
}

func literalReferenceGetReferencedLiteral(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralReferenceGetReferencedLiteralUri)
	if original == nil {
		log.Printf("In GetReferencedLiteral the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri, hl)
	if sourceLiteralReferenceRef == nil {
		log.Printf("In GetReferencedLiteral, the SourceReference was not found in the replicate")
		return
	}

	indicatedLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralReference == nil {
		log.Printf("In GetReferencedLiteral, the TargetLiteralReference was not found in the replicate")
		return
	}

	indicatedLiteral := indicatedLiteralReference.GetReferencedLiteral(hl)
	untypedLiteralReference := sourceLiteralReferenceRef.GetReferencedElement(hl)
	var sourceLiteralReference core.LiteralReference
	var sourceLiteral core.Literal
	if untypedLiteralReference != nil {
		switch untypedLiteralReference.(type) {
		case core.LiteralReference:
			sourceLiteralReference = untypedLiteralReference.(core.LiteralReference)
			sourceLiteral = sourceLiteralReference.GetReferencedLiteral(hl)
		}
	}
	if sourceLiteral != indicatedLiteral {
		indicatedLiteralReference.SetReferencedLiteral(sourceLiteral, hl)
	}
}

func literalReferenceGetLiteralPointer(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralReferenceGetLiteralPointerUri)
	if original == nil {
		log.Printf("In GetLiteralPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri, hl)
	if sourceLiteralReferenceRef == nil {
		log.Printf("In GetLiteralPointer, the SourceLiteralReferenceRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		log.Printf("In GetLiteralPointer, the IndicatedLiteralPointerRef was not found in the replicate")
		core.Print(replicate, "Replicate: ", hl)
		return
	}

	indicatedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	untypedSourceLiteralReference := sourceLiteralReferenceRef.GetReferencedElement(hl)
	var sourceLiteralPointer core.LiteralPointer
	if untypedSourceLiteralReference != nil {
		switch untypedSourceLiteralReference.(type) {
		case core.LiteralReference:
			sourceLiteralReference := untypedSourceLiteralReference.(core.LiteralReference)
			sourceLiteralPointer = sourceLiteralReference.GetLiteralPointer(hl)
		default:
			log.Printf("In GetLiteralPointer, the SourceElement is not a LiteralReference")
		}
	}
	if sourceLiteralPointer != indicatedLiteralPointer {
		indicatedLiteralPointerRef.SetReferencedLiteralPointer(sourceLiteralPointer, hl)
	}
}

func literalReferenceSetReferencedLiteral(replicate core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralReferenceSetReferencedLiteralUri)
	if original == nil {
		log.Printf("In SetReferencedLiteral the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralReferenceSetReferencedLiteralSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		log.Printf("In SetReferencedLiteral, the LiteralReference was not found in the replicate")
		return
	}

	modifiedLiteralReferenceRef := core.GetChildElementReferenceWithAncestorUri(replicate, LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri, hl)
	if modifiedLiteralReferenceRef == nil {
		log.Printf("In SetReferencedLiteral, the TargetElement was not found in the replicate")
		return
	}

	sourceLiteral := sourceLiteralRef.GetReferencedLiteral(hl)
	untypedLiteralReference := modifiedLiteralReferenceRef.GetReferencedElement(hl)
	var currentLiteral core.Literal
	var modifiedLiteralReference core.LiteralReference
	if untypedLiteralReference != nil {
		switch untypedLiteralReference.(type) {
		case core.LiteralReference:
			modifiedLiteralReference = untypedLiteralReference.(core.LiteralReference)
			currentLiteral = modifiedLiteralReference.GetReferencedLiteral(hl)
		default:
			log.Printf("In SetReferencedLiteral, the TargetedLiteralReference is not a LiteralReference")
		}
	}
	if sourceLiteral != currentLiteral {
		modifiedLiteralReference.SetReferencedLiteral(sourceLiteral, hl)
	}
}

func UpdateRecoveredCoreLiteralReferenceFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralReferenceFunctions
	elementReferenceFunctions := uOfD.GetElementWithUri(LiteralReferenceFunctionsUri)
	if elementReferenceFunctions == nil {
		elementReferenceFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(elementReferenceFunctions, coreFunctionsElement, hl)
		core.SetName(elementReferenceFunctions, "LiteralReferenceFunctions", hl)
		core.SetUri(elementReferenceFunctions, LiteralReferenceFunctionsUri, hl)
	}

	// CreateLiteralReference
	literalReferenceCreateLiteralReference := uOfD.GetElementWithUri(LiteralReferenceCreateUri)
	if literalReferenceCreateLiteralReference == nil {
		literalReferenceCreateLiteralReference = uOfD.NewElement(hl)
		core.SetOwningElement(literalReferenceCreateLiteralReference, elementReferenceFunctions, hl)
		core.SetName(literalReferenceCreateLiteralReference, "CreateLiteralReference", hl)
		core.SetUri(literalReferenceCreateLiteralReference, LiteralReferenceCreateUri, hl)
	}
	// CreatedLiteralReference
	createdLiteralReferenceRef := core.GetChildElementReferenceWithUri(literalReferenceCreateLiteralReference, LiteralReferenceCreateCreatedLiteralReferenceRefUri, hl)
	if createdLiteralReferenceRef == nil {
		createdLiteralReferenceRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdLiteralReferenceRef, literalReferenceCreateLiteralReference, hl)
		core.SetName(createdLiteralReferenceRef, "CreatedLiteralReferenceRef", hl)
		core.SetUri(createdLiteralReferenceRef, LiteralReferenceCreateCreatedLiteralReferenceRefUri, hl)
	}

	// GetReferencedLiteral
	literalReferenceGetReferencedLiteral := uOfD.GetElementWithUri(LiteralReferenceGetReferencedLiteralUri)
	if literalReferenceGetReferencedLiteral == nil {
		literalReferenceGetReferencedLiteral = uOfD.NewElement(hl)
		core.SetName(literalReferenceGetReferencedLiteral, "GetReferencedLiteral", hl)
		core.SetOwningElement(literalReferenceGetReferencedLiteral, elementReferenceFunctions, hl)
		core.SetUri(literalReferenceGetReferencedLiteral, LiteralReferenceGetReferencedLiteralUri, hl)
	}
	// GetReferencedLiteral.SourceReference
	getElementSourceReference := core.GetChildElementReferenceWithUri(literalReferenceGetReferencedLiteral, LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri, hl)
	if getElementSourceReference == nil {
		getElementSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementSourceReference, literalReferenceGetReferencedLiteral, hl)
		core.SetName(getElementSourceReference, "SourceLiteralReferenceRef", hl)
		core.SetUri(getElementSourceReference, LiteralReferenceGetReferencedLiteralSourceLiteralReferenceRefUri, hl)
	}
	// GetReferencedLiteralTargetLiteralReference
	getElementTargetReference := core.GetChildLiteralReferenceWithUri(literalReferenceGetReferencedLiteral, LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri, hl)
	if getElementTargetReference == nil {
		getElementTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getElementTargetReference, literalReferenceGetReferencedLiteral, hl)
		core.SetName(getElementTargetReference, "IndicatedLiteralRef", hl)
		core.SetUri(getElementTargetReference, LiteralReferenceGetReferencedLiteralIndicatedLiteralRefUri, hl)
	}

	// GetLiteralPointer
	literalReferenceGetLiteralPointer := uOfD.GetElementWithUri(LiteralReferenceGetLiteralPointerUri)
	if literalReferenceGetLiteralPointer == nil {
		literalReferenceGetLiteralPointer = uOfD.NewElement(hl)
		core.SetName(literalReferenceGetLiteralPointer, "GetLiteralPointer", hl)
		core.SetOwningElement(literalReferenceGetLiteralPointer, elementReferenceFunctions, hl)
		core.SetUri(literalReferenceGetLiteralPointer, LiteralReferenceGetLiteralPointerUri, hl)
	}
	// GetLiteralPointer.SourceReference
	getElementPointerSourceReference := core.GetChildElementReferenceWithUri(literalReferenceGetLiteralPointer, LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri, hl)
	if getElementPointerSourceReference == nil {
		getElementPointerSourceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(getElementPointerSourceReference, literalReferenceGetLiteralPointer, hl)
		core.SetName(getElementPointerSourceReference, "SourceLiteralReferenceRef", hl)
		core.SetUri(getElementPointerSourceReference, LiteralReferenceGetLiteralPointerSourceLiteralReferenceRefUri, hl)
	}
	// GetLiteralPointerIndicatedLiteralPointerRef
	getElementPointerIndicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithUri(literalReferenceGetLiteralPointer, LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if getElementPointerIndicatedLiteralPointerRef == nil {
		getElementPointerIndicatedLiteralPointerRef = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(getElementPointerIndicatedLiteralPointerRef, literalReferenceGetLiteralPointer, hl)
		core.SetName(getElementPointerIndicatedLiteralPointerRef, "IndicatedLiteralPointerRef", hl)
		core.SetUri(getElementPointerIndicatedLiteralPointerRef, LiteralReferenceGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	}

	// SetReferencedLiteral
	literalReferenceSetReferencedLiteral := uOfD.GetElementWithUri(LiteralReferenceSetReferencedLiteralUri)
	if literalReferenceSetReferencedLiteral == nil {
		literalReferenceSetReferencedLiteral = uOfD.NewElement(hl)
		core.SetName(literalReferenceSetReferencedLiteral, "SetReferencedLiteral", hl)
		core.SetOwningElement(literalReferenceSetReferencedLiteral, elementReferenceFunctions, hl)
		core.SetUri(literalReferenceSetReferencedLiteral, LiteralReferenceSetReferencedLiteralUri, hl)
	}
	// SetReferencedLiteral.LiteralReference
	setReferencedElementLiteralReference := core.GetChildLiteralReferenceWithUri(literalReferenceSetReferencedLiteral, LiteralReferenceSetReferencedLiteralSourceLiteralRefUri, hl)
	if setReferencedElementLiteralReference == nil {
		setReferencedElementLiteralReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(setReferencedElementLiteralReference, literalReferenceSetReferencedLiteral, hl)
		core.SetName(setReferencedElementLiteralReference, "SourceLiteralRef", hl)
		core.SetUri(setReferencedElementLiteralReference, LiteralReferenceSetReferencedLiteralSourceLiteralRefUri, hl)
	}
	// SetReferencedLiteralTargetLiteralReference
	setReferencedElementTargetLiteralReference := core.GetChildElementReferenceWithUri(literalReferenceSetReferencedLiteral, LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri, hl)
	if setReferencedElementTargetLiteralReference == nil {
		setReferencedElementTargetLiteralReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(setReferencedElementTargetLiteralReference, literalReferenceSetReferencedLiteral, hl)
		core.SetName(setReferencedElementTargetLiteralReference, "ModifiedLiteralReferenceRef", hl)
		core.SetUri(setReferencedElementTargetLiteralReference, LiteralReferenceSetReferencedLiteralModifiedLiteralReferenceRefUri, hl)
	}
}

func literalReferenceFunctionsInit() {
	core.GetCore().AddFunction(LiteralReferenceCreateUri, literalReferenceCreateLiteralReference)
	core.GetCore().AddFunction(LiteralReferenceGetReferencedLiteralUri, literalReferenceGetReferencedLiteral)
	core.GetCore().AddFunction(LiteralReferenceGetLiteralPointerUri, literalReferenceGetLiteralPointer)
	core.GetCore().AddFunction(LiteralReferenceSetReferencedLiteralUri, literalReferenceSetReferencedLiteral)
}
