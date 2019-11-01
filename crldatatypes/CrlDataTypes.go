package crldatatypes

import (
	"github.com/pbrown12303/activeCRL/core"
)

// CrlDataTypesConceptSpaceURI is the URI for the concpet space that defines the CRL Data Types
var CrlDataTypesConceptSpaceURI = "http://activeCRL.com/crldatastructures/CrlDataTypes"

// BuildCrlDataTypesConceptSpace constructs the concept space for CRL data structures
func BuildCrlDataTypesConceptSpace(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {
	crlDataTypes, _ := uOfD.NewElement(hl, CrlDataTypesConceptSpaceURI)
	crlDataTypes.SetLabel("CrlDataTypes", hl)
	BuildCrlBooleanConcept(uOfD, crlDataTypes, hl)
}
