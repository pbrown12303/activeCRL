package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
)

func createElement(element core.Element, changeNotification *core.ChangeNotification) {
	core.PrintMutex.Lock()
	defer core.PrintMutex.Unlock()
	//	log.Printf("CreateElement called")
	element.TraceableLock()
	defer element.TraceableUnlock()
	//	core.Print(element, "+++")
	//	core.PrintNotification(changeNotification)
	uOfD := element.GetUniverseOfDiscourseNoLock()
	createdElementReference := core.GetChildElementReferenceWithAncestorUriNoLock(element, CreatedElementReferenceUri)
	if createdElementReference == nil {
		//		log.Printf("CreateElementReference nil, creating element reference")
		createdElementReference = uOfD.NewElementReference()
		createdElementReference.SetOwningElementNoLock(element)
		createdElementReference.SetNameNoLock("CreatedElementReference")
		rootCreatedElementReference := uOfD.GetElementReferenceWithUri(CreatedElementReferenceUri)
		refinement := uOfD.NewRefinement()
		refinement.SetOwningElementNoLock(createdElementReference)
		refinement.SetRefinedElement(createdElementReference)
		refinement.SetAbstractElement(rootCreatedElementReference)
	}
	//	log.Printf("CreateElement nil, creating element")
	createdElement := createdElementReference.GetReferencedElement()
	if createdElement == nil {
		createdElement = uOfD.NewElement()
		createdElementReference.SetReferencedElement(createdElement)
	}
	//	log.Printf("At end of createElement")
	//	core.Print(element, "   ")
}
