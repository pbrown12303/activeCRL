package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"testing"
	"time"
)

func TestCreateElementFunction(t *testing.T) {
	uOfD := core.NewUniverseOfDiscourse()
	GetCoreFunctionsConceptSpace(uOfD)
	createElement := uOfD.GetElementWithUri(CreateElememtUri)
	if createElement == nil {
		t.Error("CreateElement not found")
	}
	createdElementReference := uOfD.GetElementReferenceWithUri(CreatedElementReferenceUri)
	if createdElementReference == nil {
		t.Error("CreatedElementReference not found")
	}
	createElementInstance := uOfD.NewElement()
	refinementInstance := uOfD.NewRefinement()
	refinementInstance.SetAbstractElement(createElement)

	refinementInstance.SetRefinedElement(createElementInstance)
	time.Sleep(1 * time.Second)
	log.Printf("Created Element Instance:")
	core.Print(createElementInstance, "---")
	var foundReference core.Element
	children := createElementInstance.GetOwnedElementsNoLock()
	if len(children) != 1 {
		log.Printf("Length of children: %d", len(children))
		t.Error("Length of children != 1")
	}
	for _, child := range children {
		ancestors := child.GetAbstractElementsRecursivelyNoLock()
		for _, ancestor := range ancestors {
			if ancestor == createdElementReference {
				foundReference = child
			}
		}
	}
	if foundReference == nil {
		t.Error("Reference not created")
	}

}
