package core

import ()

var CoreUri string = "http://activeCrl.com/core/Core"
var ElememtUri string = "http://activeCrl.com/core/Element"
var ElementPointerUri string = "http://activeCrl.com/core/ElementPointer"
var ElementPointerPointerUri string = "http://activeCrl.com/core/ElementPointerPointer"
var ElementPointerReferenceUri string = "http://activeCrl.com/core/ElementPointerReference"
var ElementReferenceUri string = "http://activeCrl.com/core/ElementReference"
var LiteralUri string = "http://activeCrl.com/core/Literal"
var LiteralPointerUri string = "http://activeCrl.com/core/LiteralPointer"
var LiteralPointerPointerUri string = "http://activeCrl.com/core/LiteralPointerPointer"
var LiteralPointerReferenceUri string = "http://activeCrl.com/core/LiteralPointerReference"
var LiteralReferenceUri string = "http://activeCrl.com/core/LiteralReference"
var RefinementUri string = "http://activeCrl.com/core/Refinement"

func GetCore(uOfD *UniverseOfDiscourse) Element {
	return RecoverElement([]byte(serializedCore), uOfD)
}
