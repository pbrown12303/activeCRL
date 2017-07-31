package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var CoreFunctionsUri string = "http://activeCrl.com/coreFunctions/CoreFunctions"
var CreateElememtUri string = "http://activeCrl.com/coreFunctions/CreateElement"
var CreatedElementReferenceUri = "http://activeCrl.com/coreFunctions/CreatedElementReference"
var CreateElementPointerUri string = "http://activeCrl.com/coreFunctions/CreateElementPointer"
var CreateElementPointerPointerUri string = "http://activeCrl.com/coreFunctions/CreateElementPointerPointer"
var CreateElementPointerReferenceUri string = "http://activeCrl.com/coreFunctions/CreateElementPointerReference"
var CreateElementReferenceUri string = "http://activeCrl.com/coreFunctions/CreateElementReference"
var CreateLiteralUri string = "http://activeCrl.com/coreFunctions/CreateLiteral"
var CreateLiteralPointerUri string = "http://activeCrl.com/coreFunctions/CreateLiteralPointer"
var CreateLiteralPointerPointerUri string = "http://activeCrl.com/coreFunctions/CreateLiteralPointerPointer"
var CreateLiteralPointerReferenceUri string = "http://activeCrl.com/coreFunctions/CreateLiteralPointerReference"
var CreateLiteralReferenceUri string = "http://activeCrl.com/coreFunctions/CreateLiteralReference"
var CreateRefinementUri string = "http://activeCrl.com/coreFunctions/CreateRefinement"

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
	//	log.Printf("About to initialize CoreFunctions")
	core.GetCore().AddFunction(CreateElememtUri, createElement)
}
