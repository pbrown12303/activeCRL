package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var BaseElementReferenceCreateUri string = CoreFunctionsPrefix + "BaseElementReference/Create"
var BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri = CoreFunctionsPrefix + "BaseElementReference/Create/CreatedBaseElementReferenceReference"

func createBaseElementReference(element core.Element, changeNotification *core.ChangeNotification) {
	log.Printf("In createBaseElementReference")
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdBaseElementReferenceReference := core.GetChildElementReferenceWithAncestorUri(element, BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri, hl)
	if createdBaseElementReferenceReference == nil {
		createdBaseElementReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdBaseElementReferenceReference, element, hl)
		core.SetName(createdBaseElementReferenceReference, "CreatedBaseElementReferenceReference", hl)
		rootCreatedBaseElementReferenceReference := uOfD.GetBaseElementReferenceWithUri(BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdBaseElementReferenceReference, hl)
		refinement.SetRefinedElement(createdBaseElementReferenceReference, hl)
		refinement.SetAbstractElement(rootCreatedBaseElementReferenceReference, hl)
	}
	createdBaseElementReference := createdBaseElementReferenceReference.GetReferencedElement(hl)
	if createdBaseElementReference == nil {
		createdBaseElementReference = uOfD.NewBaseElementReference(hl)
		createdBaseElementReferenceReference.SetReferencedElement(createdBaseElementReference, hl)
	}
}

func UpdateRecoveredCoreBaseElementReferenceFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// CreateBaseElementReference
	createBaseElementReference := uOfD.GetElementWithUri(BaseElementReferenceCreateUri)
	if createBaseElementReference == nil {
		createBaseElementReference = uOfD.NewElement(hl)
		core.SetOwningElement(createBaseElementReference, coreFunctionsElement, hl)
		core.SetName(createBaseElementReference, "CreateBaseElementReference", hl)
		core.SetUri(createBaseElementReference, BaseElementReferenceCreateUri, hl)
	}
	// CreatedBaseElementReference
	createdBaseElementReferenceReference := core.GetChildElementReferenceWithUri(createBaseElementReference, BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri, hl)
	if createdBaseElementReferenceReference == nil {
		createdBaseElementReferenceReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdBaseElementReferenceReference, createBaseElementReference, hl)
		core.SetName(createdBaseElementReferenceReference, "CreatedBaseElementReferenceReference", hl)
		core.SetUri(createdBaseElementReferenceReference, BaseElementReferenceCreateCreatedBaseElementReferenceReferenceUri, hl)
	}
}

func baseElementReferenceFunctionsInit() {
	core.GetCore().AddFunction(BaseElementReferenceCreateUri, createBaseElementReference)
	//	core.GetCore().AddFunction(BaseElementPointerGetBaseElementUri, getBaseElement)
	//	core.GetCore().AddFunction(BaseElementPointerGetBaseElementIdUri, getBaseElementId)
	//	core.GetCore().AddFunction(BaseElementPointerGetBaseElementVersionUri, getBaseElementVersion)
	//	core.GetCore().AddFunction(BaseElementPointerSetBaseElementUri, setBaseElement)
}
