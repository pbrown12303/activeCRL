package coreSimpleTypes

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	//	"log"
)

var BooleanUri string = CoreSimpleTypePrefix + "Boolean"
var BooleanTrueUri string = BooleanUri + "/" + "True"
var BooleanFalseUri string = BooleanUri + "/" + "False"

// BuildBooleanSimpleTypes constructs the Boolean concept and the true and false values programmatically, making them
// children of the coreSimpleTypesConceptSpace.
func BuildBooleanSimpleTypes(coreSimpleTypesConceptSpace core.Element, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {

	crlBoolean := uOfD.NewElement(hl, BooleanUri)
	core.SetLabel(crlBoolean, "Boolean", hl)
	core.SetUri(crlBoolean, BooleanUri, hl)
	core.SetOwningElement(crlBoolean, coreSimpleTypesConceptSpace, hl)

	crlTrueValue := "true"
	crlBooleanTrue := uOfD.NewLiteralReference(hl, BooleanTrueUri)
	core.SetLabel(crlBooleanTrue, crlTrueValue, hl)
	core.SetUri(crlBooleanTrue, BooleanTrueUri, hl)
	core.SetOwningElement(crlBooleanTrue, coreSimpleTypesConceptSpace, hl)

	crlBooleanTrueRefinement := uOfD.NewRefinement(hl)
	crlBooleanTrueRefinement.SetAbstractElement(crlBoolean, hl)
	crlBooleanTrueRefinement.SetRefinedElement(crlBooleanTrue, hl)
	core.SetOwningElement(crlBooleanTrueRefinement, crlBooleanTrue, hl)

	crlFalseValue := "false"
	crlBooleanFalse := uOfD.NewLiteralReference(hl, BooleanFalseUri)
	core.SetLabel(crlBooleanFalse, crlFalseValue, hl)
	core.SetUri(crlBooleanFalse, BooleanFalseUri, hl)
	core.SetOwningElement(crlBooleanFalse, coreSimpleTypesConceptSpace, hl)

	crlBooleanFalseRefinement := uOfD.NewRefinement(hl)
	crlBooleanFalseRefinement.SetAbstractElement(crlBoolean, hl)
	crlBooleanFalseRefinement.SetRefinedElement(crlBooleanFalse, hl)
	core.SetOwningElement(crlBooleanFalseRefinement, crlBooleanFalse, hl)
}
