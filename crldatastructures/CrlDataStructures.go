package crldatastructures

import (
	"github.com/pbrown12303/activeCRL/core"
)

// CrlDataStructuresConceptSpaceURI is the uri for the concept space that defines the Crl Data Structures
var CrlDataStructuresConceptSpaceURI = "http://activeCRL.com/crldatastructures/CrlDataStructures"

// BuildCrlDataStructuresConceptSpace constructs the concept space for CRL data structures
func BuildCrlDataStructuresConceptSpace(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {
	crlDataStructures, _ := uOfD.NewElement(hl, CrlDataStructuresConceptSpaceURI)
	crlDataStructures.SetLabel("CrlDataStructures", hl)
	BuildCrlSetsConcepts(uOfD, crlDataStructures, hl)
	BuildCrlListsConcepts(uOfD, crlDataStructures, hl)
	crlDataStructures.SetIsCoreRecursively(hl)
}
