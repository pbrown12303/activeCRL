package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
)

var CreateElememtUri string = CoreFunctionsPrefix + "CreateElement"
var CreatedElementReferenceUri = CoreFunctionsPrefix + "CreatedElementReference"

func createElement(element core.Element, changeNotification *core.ChangeNotification) {
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementReference := core.GetChildElementReferenceWithAncestorUri(element, CreatedElementReferenceUri, hl)
	if createdElementReference == nil {
		createdElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementReference, element, hl)
		core.SetName(createdElementReference, "CreatedElementReference", hl)
		rootCreatedElementReference := uOfD.GetElementReferenceWithUri(CreatedElementReferenceUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdElementReference, hl)
		refinement.SetRefinedElement(createdElementReference, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdElement := createdElementReference.GetReferencedElement(hl)
	if createdElement == nil {
		createdElement = uOfD.NewElement(hl)
		createdElementReference.SetReferencedElement(createdElement, hl)
	}
}

func UpdateRecoveredCoreElementFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// CreateElement
	createElement := uOfD.GetElementWithUri(CreateElememtUri)
	if createElement == nil {
		createElement = uOfD.NewElement(hl)
		core.SetOwningElement(createElement, coreFunctionsElement, hl)
		core.SetName(createElement, "CreateElement", hl)
		core.SetUri(createElement, CreateElememtUri, hl)
	}
	// CreatedElementReference
	createdElementReference := core.GetChildElementReferenceWithUri(createElement, CreatedElementReferenceUri, hl)
	if createdElementReference == nil {
		createdElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementReference, createElement, hl)
		core.SetName(createdElementReference, "CreatedElementReference", hl)
		core.SetUri(createdElementReference, CreatedElementReferenceUri, hl)
	}

}

func elementFunctionsInit() {
	core.GetCore().AddFunction(CreateElememtUri, createElement)
}
