package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var CoreFunctionsPrefix string = "http://activeCrl.com/coreFunctions/"
var CoreFunctionsUri string = CoreFunctionsPrefix + "CoreFunctions"

func GetCoreFunctionsConceptSpace(uOfD *core.UniverseOfDiscourse) core.Element {
	coreFunctionsConceptSpace := uOfD.GetElementWithUri(CoreFunctionsUri)
	if coreFunctionsConceptSpace == nil {
		coreFunctionsConceptSpace = uOfD.RecoverElement([]byte(serializedCoreFunctions))
		if coreFunctionsConceptSpace == nil {
			log.Printf("Recovery of CoreFunctions failed")
		}
	}
	return coreFunctionsConceptSpace
}

func init() {
	baseElementFunctionsInit()
	baseElementPointerFunctionsInit()
	baseElementReferenceFunctionsInit()
	elementFunctionsInit()
}
