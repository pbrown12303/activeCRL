package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"sync"
)

var ElementFunctionsUri string = CoreFunctionsPrefix + "ElementFunctions"

var ElementCreateUri string = CoreFunctionsPrefix + "Element/Create"
var ElementCreateCreatedElementRefUri = CoreFunctionsPrefix + "Element/Create/CreatedElementRef"

var ElementGetDefinitionUri string = CoreFunctionsPrefix + "Element/GetDefinition"
var ElementGetDefinitionSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetDefinition/SourceElementRef"
var ElementGetDefinitionCreatedLiteralRefUri string = CoreFunctionsPrefix + "Element/GetDefinition/CreatedLiteralRef"

var ElementGetDefinitionLiteralUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteral"
var ElementGetDefinitionLiteralSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteral/SourceElementRef"
var ElementGetDefinitionLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteral/IndicatedLiteralRef"

var ElementGetDefinitionLiteralPointerUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteralPointer"
var ElementGetDefinitionLiteralPointerSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteralPointer/SourceElementRef"
var ElementGetDefinitionLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "Element/GetDefinitionLiteralPointer/IndicatedLiteralPointerRef"

var ElementGetNameLiteralUri string = CoreFunctionsPrefix + "Element/GetNameLiteral"
var ElementGetNameLiteralSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetNameLiteral/SourceElementRef"
var ElementGetNameLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "Element/GetNameLiteral/IndicatedLiteralRef"

var ElementGetNameLiteralPointerUri string = CoreFunctionsPrefix + "Element/GetNameLiteralPointer"
var ElementGetNameLiteralPointerSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetNameLiteralPointer/SourceElementRef"
var ElementGetNameLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "Element/GetNameLiteralPointer/IndicatedLiteralPointerRef"

var ElementGetUriLiteralUri string = CoreFunctionsPrefix + "Element/GetUriLiteral"
var ElementGetUriLiteralSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetUriLiteral/SourceElementRef"
var ElementGetUriLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "Element/GetUriLiteral/IndicatedLiteralRef"

var ElementGetUriLiteralPointerUri string = CoreFunctionsPrefix + "Element/GetUriLiteralPointer"
var ElementGetUriLiteralPointerSourceElementRefUri string = CoreFunctionsPrefix + "Element/GetUriLiteralPointer/SourceElementRef"
var ElementGetUriLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "Element/GetUriLiteralPointer/IndicatedLiteralPointerRef"

var ElementSetDefinitionUri string = CoreFunctionsPrefix + "Element/SetDefinition"
var ElementSetDefinitionSourceLiteralRefUri string = CoreFunctionsPrefix + "Element/SetDefinition/SourceLiteralRef"
var ElementSetDefinitionModifiedElementRefUri string = CoreFunctionsPrefix + "Element/SetDefinition/ModifiedElementRef"

var ElementSetNameUri string = CoreFunctionsPrefix + "Element/SetName"
var ElementSetNameSourceLiteralRefUri string = CoreFunctionsPrefix + "Element/SetName/SourceLiteralRef"
var ElementSetNameModifiedElementRefUri string = CoreFunctionsPrefix + "Element/SetName/ModifiedElementRef"

func createElement(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementRef := core.GetChildElementReferenceWithAncestorUri(element, ElementCreateCreatedElementRefUri, hl)
	if createdElementRef == nil {
		createdElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementRef, element, hl)
		core.SetName(createdElementRef, "CreatedElementRef", hl)
		rootCreatedElementReference := uOfD.GetElementReferenceWithUri(ElementCreateCreatedElementRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementRef, hl)
		refinement.SetRefinedElement(createdElementRef, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdElement := createdElementRef.GetReferencedElement(hl)
	if createdElement == nil {
		createdElement = uOfD.NewElement(hl)
		createdElementRef.SetReferencedElement(createdElement, hl)
	}
}

func getDefinition(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementGetDefinitionUri)
	if original == nil {
		log.Printf("In GetDefinition the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementGetDefinitionSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In GetDefinition, the SourceElementRef was not found in the replicate")
		return
	}

	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementGetDefinitionCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		log.Printf("In GetDefinition, the CreatedLiteralRef was not found in the replicate")
		return
	}

	currentLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if currentLiteral == nil {
		currentLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(currentLiteral, createdLiteralRef, hl)
		createdLiteralRef.SetReferencedLiteral(currentLiteral, hl)
	}

	var sourceDefinition string = ""
	sourceElement := sourceElementRef.GetReferencedElement(hl)
	if sourceElement != nil {
		sourceDefinition = sourceElement.GetDefinition(hl)
	}
	if sourceDefinition != currentLiteral.GetLiteralValue(hl) {
		currentLiteral.SetLiteralValue(sourceDefinition, hl)
	}
}

func getDefinitionLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementGetDefinitionLiteralUri)
	if original == nil {
		log.Printf("In GetDefinitionLiteral the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementGetDefinitionLiteralSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In GetDefinitionLiteral, the SourceElementRef was not found in the replicate")
		return
	}

	indicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementGetDefinitionLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		log.Printf("In GetGetDefinitionLiteral, the IndicatedLiteralRef was not found in the replicate")
		return
	}

	indicatedLiteral := indicatedLiteralRef.GetReferencedLiteral(hl)
	sourceElement := sourceElementRef.GetReferencedElement(hl)
	var sourceDefinitionLiteral core.Literal
	if sourceElement != nil {
		sourceDefinitionLiteral = sourceElement.GetDefinitionLiteral(hl)
	}
	if sourceDefinitionLiteral != indicatedLiteral {
		indicatedLiteralRef.SetReferencedLiteral(sourceDefinitionLiteral, hl)
	}
}

func getDefinitionLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementGetDefinitionLiteralPointerUri)
	if original == nil {
		log.Printf("In GetDefinitionLiteralPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementGetDefinitionLiteralPointerSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In GetDefinitionLiteralPointer, the SourceElementRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, ElementGetDefinitionLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		log.Printf("In GetGetDefinitionLiteralPointer, the IndicatedLiteralPointerRef was not found in the replicate")
		return
	}

	indicatedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	sourceElement := sourceElementRef.GetReferencedElement(hl)
	var sourceDefinitionLiteralPointer core.LiteralPointer
	if sourceElement != nil {
		sourceDefinitionLiteralPointer = sourceElement.GetDefinitionLiteralPointer(hl)
	}
	if sourceDefinitionLiteralPointer != indicatedLiteralPointer {
		indicatedLiteralPointerRef.SetReferencedLiteralPointer(sourceDefinitionLiteralPointer, hl)
	}
}

func getNameLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementGetNameLiteralUri)
	if original == nil {
		log.Printf("In GetNameLiteral the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementGetNameLiteralSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In GetNameLiteral, the SourceElementRef was not found in the replicate")
		return
	}

	indicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementGetNameLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		log.Printf("In GetGetNameLiteral, the IndicatedLiteralRef was not found in the replicate")
		return
	}

	indicatedLiteral := indicatedLiteralRef.GetReferencedLiteral(hl)
	sourceElement := sourceElementRef.GetReferencedElement(hl)
	var sourceNameLiteral core.Literal
	if sourceElement != nil {
		sourceNameLiteral = sourceElement.GetNameLiteral(hl)
	}
	if sourceNameLiteral != indicatedLiteral {
		indicatedLiteralRef.SetReferencedLiteral(sourceNameLiteral, hl)
	}
}

func getNameLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementGetNameLiteralPointerUri)
	if original == nil {
		log.Printf("In GetNameLiteralPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementGetNameLiteralPointerSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In GetNameLiteralPointer, the SourceElementRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, ElementGetNameLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		log.Printf("In GetGetNameLiteralPointer, the IndicatedLiteralPointerRef was not found in the replicate")
		return
	}

	indicatedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	sourceElement := sourceElementRef.GetReferencedElement(hl)
	var sourceNameLiteralPointer core.LiteralPointer
	if sourceElement != nil {
		sourceNameLiteralPointer = sourceElement.GetNameLiteralPointer(hl)
	}
	if sourceNameLiteralPointer != indicatedLiteralPointer {
		indicatedLiteralPointerRef.SetReferencedLiteralPointer(sourceNameLiteralPointer, hl)
	}
}

func getUriLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementGetUriLiteralUri)
	if original == nil {
		log.Printf("In GetUriLiteral the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementGetUriLiteralSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In GetUriLiteral, the SourceElementRef was not found in the replicate")
		return
	}

	indicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementGetUriLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		log.Printf("In GetGetUriLiteral, the IndicatedLiteralRef was not found in the replicate")
		return
	}

	indicatedLiteral := indicatedLiteralRef.GetReferencedLiteral(hl)
	sourceElement := sourceElementRef.GetReferencedElement(hl)
	var sourceUriLiteral core.Literal
	if sourceElement != nil {
		sourceUriLiteral = sourceElement.GetUriLiteral(hl)
	}
	if sourceUriLiteral != indicatedLiteral {
		indicatedLiteralRef.SetReferencedLiteral(sourceUriLiteral, hl)
	}
}

func getUriLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementGetUriLiteralPointerUri)
	if original == nil {
		log.Printf("In GetUriLiteralPointer the original operation was not found")
		return
	}
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementGetUriLiteralPointerSourceElementRefUri, hl)
	if sourceElementRef == nil {
		log.Printf("In GetUriLiteralPointer, the SourceElementRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, ElementGetUriLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		log.Printf("In GetGetUriLiteralPointer, the IndicatedLiteralPointerRef was not found in the replicate")
		return
	}

	indicatedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	sourceElement := sourceElementRef.GetReferencedElement(hl)
	var sourceUriLiteralPointer core.LiteralPointer
	if sourceElement != nil {
		sourceUriLiteralPointer = sourceElement.GetUriLiteralPointer(hl)
	}
	if sourceUriLiteralPointer != indicatedLiteralPointer {
		indicatedLiteralPointerRef.SetReferencedLiteralPointer(sourceUriLiteralPointer, hl)
	}

}

func setDefinition(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementSetDefinitionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementSetDefinitionSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		log.Printf("In SetDefinition, the SourceLiteralRef was not found in the replicate")
		return
	}

	modifiedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementSetDefinitionModifiedElementRefUri, hl)
	if modifiedElementRef == nil {
		log.Printf("In SetDefinition, the ModifiedElementRef was not found in the replicate")
		return
	}

	modifiedElement := modifiedElementRef.GetReferencedElement(hl)
	sourceLiteral := sourceLiteralRef.GetReferencedLiteral(hl)
	if modifiedElement != nil {
		core.SetDefinition(modifiedElement, sourceLiteral.GetLiteralValue(hl), hl)
	}
}

func setName(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(ElementSetNameUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, ElementSetNameSourceLiteralRefUri, hl)
	if sourceLiteralRef == nil {
		log.Printf("In SetName, the SourceLiteralRef was not found in the replicate")
		return
	}

	modifiedElementRef := core.GetChildElementReferenceWithAncestorUri(replicate, ElementSetNameModifiedElementRefUri, hl)
	if modifiedElementRef == nil {
		log.Printf("In SetName, the ModifiedElementRef was not found in the replicate")
		return
	}

	modifiedElement := modifiedElementRef.GetReferencedElement(hl)
	sourceLiteral := sourceLiteralRef.GetReferencedLiteral(hl)
	if modifiedElement != nil {
		core.SetName(modifiedElement, sourceLiteral.GetLiteralValue(hl), hl)
	}
}

func UpdateRecoveredCoreElementFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// ElementFunctions
	elementFunctions := uOfD.GetElementWithUri(ElementFunctionsUri)
	if elementFunctions == nil {
		elementFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(elementFunctions, coreFunctionsElement, hl)
		core.SetName(elementFunctions, "ElementFunctions", hl)
		core.SetUri(elementFunctions, ElementFunctionsUri, hl)
	}

	// CreateElement
	createElement := uOfD.GetElementWithUri(ElementCreateUri)
	if createElement == nil {
		createElement = uOfD.NewElement(hl)
		core.SetOwningElement(createElement, elementFunctions, hl)
		core.SetName(createElement, "CreateElement", hl)
		core.SetUri(createElement, ElementCreateUri, hl)
	}
	// CreatedElementReference
	createdElementRef := core.GetChildElementReferenceWithUri(createElement, ElementCreateCreatedElementRefUri, hl)
	if createdElementRef == nil {
		createdElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementRef, createElement, hl)
		core.SetName(createdElementRef, "CreatedElementReference", hl)
		core.SetUri(createdElementRef, ElementCreateCreatedElementRefUri, hl)
	}

	// GetDefinition
	getDefinition := uOfD.GetElementWithUri(ElementGetDefinitionUri)
	if getDefinition == nil {
		getDefinition = uOfD.NewElement(hl)
		core.SetOwningElement(getDefinition, elementFunctions, hl)
		core.SetName(getDefinition, "GetDefinition", hl)
		core.SetUri(getDefinition, ElementGetDefinitionUri, hl)
	}
	// SourceElementRef
	sourceElementRef := core.GetChildElementReferenceWithUri(getDefinition, ElementGetDefinitionSourceElementRefUri, hl)
	if sourceElementRef == nil {
		sourceElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(sourceElementRef, getDefinition, hl)
		core.SetName(sourceElementRef, "SourceElementRef", hl)
		core.SetUri(sourceElementRef, ElementGetDefinitionSourceElementRefUri, hl)
	}
	// CreatedLiteralRef
	createdLiteralRef := core.GetChildLiteralReferenceWithUri(getDefinition, ElementGetDefinitionCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		createdLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(createdLiteralRef, getDefinition, hl)
		core.SetName(createdLiteralRef, "CreatedLiteralRef", hl)
		core.SetUri(createdLiteralRef, ElementGetDefinitionCreatedLiteralRefUri, hl)
	}

	// GetDefinitionLiteral
	getDefinitionLiteral := uOfD.GetElementWithUri(ElementGetDefinitionLiteralUri)
	if getDefinitionLiteral == nil {
		getDefinitionLiteral = uOfD.NewElement(hl)
		core.SetOwningElement(getDefinitionLiteral, elementFunctions, hl)
		core.SetName(getDefinitionLiteral, "GetDefinition", hl)
		core.SetUri(getDefinitionLiteral, ElementGetDefinitionLiteralUri, hl)
	}
	// SourceElementRef
	sourceElementRef = core.GetChildElementReferenceWithUri(getDefinitionLiteral, ElementGetDefinitionLiteralSourceElementRefUri, hl)
	if sourceElementRef == nil {
		sourceElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(sourceElementRef, getDefinitionLiteral, hl)
		core.SetName(sourceElementRef, "SourceElementRef", hl)
		core.SetUri(sourceElementRef, ElementGetDefinitionLiteralSourceElementRefUri, hl)
	}
	// IndicatedLiteralRef
	indicatedLiteralRef := core.GetChildLiteralReferenceWithUri(getDefinitionLiteral, ElementGetDefinitionLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		indicatedLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(indicatedLiteralRef, getDefinitionLiteral, hl)
		core.SetName(indicatedLiteralRef, "CreatedLiteralRef", hl)
		core.SetUri(indicatedLiteralRef, ElementGetDefinitionLiteralIndicatedLiteralRefUri, hl)
	}

	// GetDefinitionLiteralPointer
	getDefinitionLiteralPointer := uOfD.GetElementWithUri(ElementGetDefinitionLiteralPointerUri)
	if getDefinitionLiteralPointer == nil {
		getDefinitionLiteralPointer = uOfD.NewElement(hl)
		core.SetOwningElement(getDefinitionLiteralPointer, elementFunctions, hl)
		core.SetName(getDefinitionLiteralPointer, "GetDefinition", hl)
		core.SetUri(getDefinitionLiteralPointer, ElementGetDefinitionLiteralPointerUri, hl)
	}
	// SourceElementRef
	sourceElementRef = core.GetChildElementReferenceWithUri(getDefinitionLiteralPointer, ElementGetDefinitionLiteralPointerSourceElementRefUri, hl)
	if sourceElementRef == nil {
		sourceElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(sourceElementRef, getDefinitionLiteralPointer, hl)
		core.SetName(sourceElementRef, "SourceElementRef", hl)
		core.SetUri(sourceElementRef, ElementGetDefinitionLiteralPointerSourceElementRefUri, hl)
	}
	// IndicatedLiteralPointerRef
	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithUri(getDefinitionLiteralPointer, ElementGetDefinitionLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		indicatedLiteralPointerRef = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(indicatedLiteralPointerRef, getDefinitionLiteralPointer, hl)
		core.SetName(indicatedLiteralPointerRef, "CreatedLiteralPointerRef", hl)
		core.SetUri(indicatedLiteralPointerRef, ElementGetDefinitionLiteralPointerIndicatedLiteralPointerRefUri, hl)
	}

	// GetNameLiteral
	getNameLiteral := uOfD.GetElementWithUri(ElementGetNameLiteralUri)
	if getNameLiteral == nil {
		getNameLiteral = uOfD.NewElement(hl)
		core.SetOwningElement(getNameLiteral, elementFunctions, hl)
		core.SetName(getNameLiteral, "GetName", hl)
		core.SetUri(getNameLiteral, ElementGetNameLiteralUri, hl)
	}
	// SourceElementRef
	sourceElementRef = core.GetChildElementReferenceWithUri(getNameLiteral, ElementGetNameLiteralSourceElementRefUri, hl)
	if sourceElementRef == nil {
		sourceElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(sourceElementRef, getNameLiteral, hl)
		core.SetName(sourceElementRef, "SourceElementRef", hl)
		core.SetUri(sourceElementRef, ElementGetNameLiteralSourceElementRefUri, hl)
	}
	// IndicatedLiteralRef
	indicatedLiteralRef = core.GetChildLiteralReferenceWithUri(getNameLiteral, ElementGetNameLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		indicatedLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(indicatedLiteralRef, getNameLiteral, hl)
		core.SetName(indicatedLiteralRef, "CreatedLiteralRef", hl)
		core.SetUri(indicatedLiteralRef, ElementGetNameLiteralIndicatedLiteralRefUri, hl)
	}

	// GetNameLiteralPointer
	getNameLiteralPointer := uOfD.GetElementWithUri(ElementGetNameLiteralPointerUri)
	if getNameLiteralPointer == nil {
		getNameLiteralPointer = uOfD.NewElement(hl)
		core.SetOwningElement(getNameLiteralPointer, elementFunctions, hl)
		core.SetName(getNameLiteralPointer, "GetName", hl)
		core.SetUri(getNameLiteralPointer, ElementGetNameLiteralPointerUri, hl)
	}
	// SourceElementRef
	sourceElementRef = core.GetChildElementReferenceWithUri(getNameLiteralPointer, ElementGetNameLiteralPointerSourceElementRefUri, hl)
	if sourceElementRef == nil {
		sourceElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(sourceElementRef, getNameLiteralPointer, hl)
		core.SetName(sourceElementRef, "SourceElementRef", hl)
		core.SetUri(sourceElementRef, ElementGetNameLiteralPointerSourceElementRefUri, hl)
	}
	// IndicatedLiteralPointerRef
	indicatedLiteralPointerRef = core.GetChildLiteralPointerReferenceWithUri(getNameLiteralPointer, ElementGetNameLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		indicatedLiteralPointerRef = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(indicatedLiteralPointerRef, getNameLiteralPointer, hl)
		core.SetName(indicatedLiteralPointerRef, "CreatedLiteralPointerRef", hl)
		core.SetUri(indicatedLiteralPointerRef, ElementGetNameLiteralPointerIndicatedLiteralPointerRefUri, hl)
	}

	// GetUriLiteral
	getUriLiteral := uOfD.GetElementWithUri(ElementGetUriLiteralUri)
	if getUriLiteral == nil {
		getUriLiteral = uOfD.NewElement(hl)
		core.SetOwningElement(getUriLiteral, elementFunctions, hl)
		core.SetName(getUriLiteral, "GetUri", hl)
		core.SetUri(getUriLiteral, ElementGetUriLiteralUri, hl)
	}
	// SourceElementRef
	sourceElementRef = core.GetChildElementReferenceWithUri(getUriLiteral, ElementGetUriLiteralSourceElementRefUri, hl)
	if sourceElementRef == nil {
		sourceElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(sourceElementRef, getUriLiteral, hl)
		core.SetName(sourceElementRef, "SourceElementRef", hl)
		core.SetUri(sourceElementRef, ElementGetUriLiteralSourceElementRefUri, hl)
	}
	// IndicatedLiteralRef
	indicatedLiteralRef = core.GetChildLiteralReferenceWithUri(getUriLiteral, ElementGetUriLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		indicatedLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(indicatedLiteralRef, getUriLiteral, hl)
		core.SetName(indicatedLiteralRef, "CreatedLiteralRef", hl)
		core.SetUri(indicatedLiteralRef, ElementGetUriLiteralIndicatedLiteralRefUri, hl)
	}

	// GetUriLiteralPointer
	getUriLiteralPointer := uOfD.GetElementWithUri(ElementGetUriLiteralPointerUri)
	if getUriLiteralPointer == nil {
		getUriLiteralPointer = uOfD.NewElement(hl)
		core.SetOwningElement(getUriLiteralPointer, elementFunctions, hl)
		core.SetName(getUriLiteralPointer, "GetUri", hl)
		core.SetUri(getUriLiteralPointer, ElementGetUriLiteralPointerUri, hl)
	}
	// SourceElementRef
	sourceElementRef = core.GetChildElementReferenceWithUri(getUriLiteralPointer, ElementGetUriLiteralPointerSourceElementRefUri, hl)
	if sourceElementRef == nil {
		sourceElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(sourceElementRef, getUriLiteralPointer, hl)
		core.SetName(sourceElementRef, "SourceElementRef", hl)
		core.SetUri(sourceElementRef, ElementGetUriLiteralPointerSourceElementRefUri, hl)
	}
	// IndicatedLiteralPointerRef
	indicatedLiteralPointerRef = core.GetChildLiteralPointerReferenceWithUri(getUriLiteralPointer, ElementGetUriLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		indicatedLiteralPointerRef = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(indicatedLiteralPointerRef, getUriLiteralPointer, hl)
		core.SetName(indicatedLiteralPointerRef, "CreatedLiteralPointerRef", hl)
		core.SetUri(indicatedLiteralPointerRef, ElementGetUriLiteralPointerIndicatedLiteralPointerRefUri, hl)
	}

	// SetDefinition
	setDefinition := uOfD.GetElementWithUri(ElementSetDefinitionUri)
	if setDefinition == nil {
		setDefinition = uOfD.NewElement(hl)
		core.SetName(setDefinition, "SetDefinition", hl)
		core.SetOwningElement(setDefinition, elementFunctions, hl)
		core.SetUri(setDefinition, ElementSetDefinitionUri, hl)
	}
	// SetDefinition.SourceLiteralRef
	setDefinitionSourceLiteralRef := core.GetChildLiteralReferenceWithUri(setDefinition, ElementSetDefinitionSourceLiteralRefUri, hl)
	if setDefinitionSourceLiteralRef == nil {
		setDefinitionSourceLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(setDefinitionSourceLiteralRef, setDefinition, hl)
		core.SetName(setDefinitionSourceLiteralRef, "SourceLiteralRefRef", hl)
		core.SetUri(setDefinitionSourceLiteralRef, ElementSetDefinitionSourceLiteralRefUri, hl)
	}
	// SetDefinitionModifiedElementReference
	setDefinitionModifiedElementRef := core.GetChildElementReferenceWithUri(setDefinition, ElementSetDefinitionModifiedElementRefUri, hl)
	if setDefinitionModifiedElementRef == nil {
		setDefinitionModifiedElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(setDefinitionModifiedElementRef, setDefinition, hl)
		core.SetName(setDefinitionModifiedElementRef, "ModifiedElementRef", hl)
		core.SetUri(setDefinitionModifiedElementRef, ElementSetDefinitionModifiedElementRefUri, hl)
	}

	// SetName
	setName := uOfD.GetElementWithUri(ElementSetNameUri)
	if setName == nil {
		setName = uOfD.NewElement(hl)
		core.SetName(setName, "SetName", hl)
		core.SetOwningElement(setName, elementFunctions, hl)
		core.SetUri(setName, ElementSetNameUri, hl)
	}
	// SetName.SourceLiteralRef
	setNameSourceLiteralRef := core.GetChildLiteralReferenceWithUri(setName, ElementSetNameSourceLiteralRefUri, hl)
	if setNameSourceLiteralRef == nil {
		setNameSourceLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(setNameSourceLiteralRef, setName, hl)
		core.SetName(setNameSourceLiteralRef, "SourceLiteralRefRef", hl)
		core.SetUri(setNameSourceLiteralRef, ElementSetNameSourceLiteralRefUri, hl)
	}
	// SetNameModifiedElementReference
	setNameModifiedElementRef := core.GetChildElementReferenceWithUri(setName, ElementSetNameModifiedElementRefUri, hl)
	if setNameModifiedElementRef == nil {
		setNameModifiedElementRef = uOfD.NewElementReference(hl)
		core.SetOwningElement(setNameModifiedElementRef, setName, hl)
		core.SetName(setNameModifiedElementRef, "ModifiedElementRef", hl)
		core.SetUri(setNameModifiedElementRef, ElementSetNameModifiedElementRefUri, hl)
	}

}

func elementFunctionsInit() {
	core.GetCore().AddFunction(ElementCreateUri, createElement)
	core.GetCore().AddFunction(ElementGetDefinitionUri, getDefinition)
	core.GetCore().AddFunction(ElementGetDefinitionLiteralUri, getDefinitionLiteral)
	core.GetCore().AddFunction(ElementGetDefinitionLiteralPointerUri, getDefinitionLiteralPointer)
	core.GetCore().AddFunction(ElementGetNameLiteralUri, getNameLiteral)
	core.GetCore().AddFunction(ElementGetNameLiteralPointerUri, getNameLiteralPointer)
	core.GetCore().AddFunction(ElementGetUriLiteralUri, getUriLiteral)
	core.GetCore().AddFunction(ElementGetUriLiteralPointerUri, getUriLiteralPointer)
	core.GetCore().AddFunction(ElementSetDefinitionUri, setDefinition)
	core.GetCore().AddFunction(ElementSetNameUri, setName)
}
