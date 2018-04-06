// The coreSimpleTypes package defines the CoreSimpleTypes concept space. This is a pre-defined concept space (hence the term "core") that is, itself,
// represented as a CRLElement and identified with the CoreSimpleTypesUri. This concept space contains the elements representing all simple types and simple
// type constants used to construct Crl designs. This package not only defines these elements but gives them URIs that can be used consistently across all Crl
// designs.
package coreSimpleTypes

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var CoreSimpleTypePrefix string = "http://activeCrl.com/coreSimpleTypes/"
var CoreSimpleTypesConceptSpaceUri string = CoreSimpleTypePrefix + "CoreSimpleTypes"

// AddCoreSimpleTypesToUofD() constructs the core simple types concept space and adds it to the specified universe of discourse.
// It provides a means of bootstrapping, assembling the core simple types concepts space programmatically. It is expected that,
// after sufficient infrastructure has been developed, this concept space may be simply stored and retrieved rather than
// assembled programatically.
func AddCoreSimpleTypesToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	coreSimpleTypesConceptSpace := uOfD.GetElementWithUri(CoreSimpleTypesConceptSpaceUri)
	if coreSimpleTypesConceptSpace == nil {
		coreSimpleTypesConceptSpace = BuildCoreSimpleTypesConceptSpace(uOfD, hl)
		if coreSimpleTypesConceptSpace == nil {
			log.Printf("Build of CoreSimpleTypes failed")
		}
	}
	return coreSimpleTypesConceptSpace
}

// BuildCoreSimpleTypesConceptSpace() builds the concept space programatically. The actual functions that build the individual
// types and their related constants (if any) are defines separately. BuildSimpleTypesConceptSpace calls each of them and then
// returns the completed concept space.
func BuildCoreSimpleTypesConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// Core
	coreSimpleTypesConceptSpace := uOfD.NewElement(hl, CoreSimpleTypesConceptSpaceUri)
	core.SetName(coreSimpleTypesConceptSpace, "CoreSimpleTypes", hl)
	core.SetUri(coreSimpleTypesConceptSpace, CoreSimpleTypesConceptSpaceUri, hl)

	BuildBooleanSimpleTypes(coreSimpleTypesConceptSpace, uOfD, hl)

	//	BuildCoreBaseElementFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreBaseElementPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreBaseElementReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementPointerPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementPointerReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralPointerPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralPointerReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreRefinementFunctions(coreFunctionsElement, uOfD, hl)

	return coreSimpleTypesConceptSpace
}

func init() {
	//	baseElementFunctionsInit()
	//	baseElementPointerFunctionsInit()
	//	baseElementReferenceFunctionsInit()
	//	elementFunctionsInit()
	//	elementPointerFunctionsInit()
	//	elementPointerPointerFunctionsInit()
	//	elementPointerReferenceFunctionsInit()
	//	elementReferenceFunctionsInit()
	//	literalFunctionsInit()
	//	literalPointerFunctionsInit()
	//	literalPointerPointerFunctionsInit()
	//	literalPointerReferenceFunctionsInit()
	//	literalReferenceFunctionsInit()
	//	refinementFunctionsInit()
}
