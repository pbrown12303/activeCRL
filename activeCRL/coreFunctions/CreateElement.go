package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
)

func createElement(element core.Element, changeNotification *core.ChangeNotification) {
	core.PrintMutex.Lock()
	defer core.PrintMutex.Unlock()
	//	log.Printf("CreateElement called")
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	//	core.Print(element, "+++", hl)
	//	core.PrintNotification(changeNotification, hl)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdElementReference := core.GetChildElementReferenceWithAncestorUri(element, CreatedElementReferenceUri, hl)
	if createdElementReference == nil {
		//		log.Printf("CreateElementReference nil, creating element reference")
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
		//		log.Printf("CreatedElement nil, creating element")
		createdElement = uOfD.NewElement(hl)
		createdElementReference.SetReferencedElement(createdElement, hl)
	}
	//	log.Printf("At end of createElement")
	//	core.Print(element, "   ", hl)
}
