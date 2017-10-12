// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var CoreFunctionsPrefix string = "http://activeCrl.com/coreFunctions/"
var CoreFunctionsUri string = CoreFunctionsPrefix + "CoreFunctions"

func AddCoreFunctionsToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	coreFunctionsConceptSpace := uOfD.GetElementWithUri(CoreFunctionsUri)
	if coreFunctionsConceptSpace == nil {
		coreFunctionsConceptSpace = BuildCoreFunctionsConceptSpace(uOfD, hl)
		if coreFunctionsConceptSpace == nil {
			log.Printf("Build of CoreFunctions failed")
		}
	}
	return coreFunctionsConceptSpace
}

func BuildCoreFunctionsConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// Core
	coreFunctionsElement := uOfD.NewElement(hl, CoreFunctionsUri)
	core.SetName(coreFunctionsElement, "CoreFunctions", hl)
	core.SetUri(coreFunctionsElement, CoreFunctionsUri, hl)

	BuildCoreBaseElementFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreBaseElementPointerFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreBaseElementReferenceFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreElementFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreElementPointerFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreElementPointerPointerFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreElementPointerReferenceFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreElementReferenceFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreLiteralFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreLiteralPointerFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreLiteralPointerPointerFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreLiteralPointerReferenceFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreLiteralReferenceFunctions(coreFunctionsElement, uOfD, hl)
	BuildCoreRefinementFunctions(coreFunctionsElement, uOfD, hl)

	return coreFunctionsElement
}

func init() {
	baseElementFunctionsInit()
	baseElementPointerFunctionsInit()
	baseElementReferenceFunctionsInit()
	elementFunctionsInit()
	elementPointerFunctionsInit()
	elementPointerPointerFunctionsInit()
	elementPointerReferenceFunctionsInit()
	elementReferenceFunctionsInit()
	literalFunctionsInit()
	literalPointerFunctionsInit()
	literalPointerPointerFunctionsInit()
	literalPointerReferenceFunctionsInit()
	literalReferenceFunctionsInit()
	refinementFunctionsInit()
}
