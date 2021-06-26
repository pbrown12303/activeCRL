package crldatastructuresdomain

import (
	"github.com/pbrown12303/activeCRL/core"
)

// CrlDataStructuresDomainURI is the uri for the concept space that defines the Crl Data Structures
var CrlDataStructuresDomainURI = "http://activeCRL.com/crldatastructuresdomain/CrlDataStructuresDomain"

// BuildCrlDataStructuresDomain constructs the concept space for CRL data structures
func BuildCrlDataStructuresDomain(uOfD *core.UniverseOfDiscourse, hl *core.Transaction) {
	crlDataStructures, _ := uOfD.NewElement(hl, CrlDataStructuresDomainURI)
	crlDataStructures.SetLabel("CrlDataStructuresDomain", hl)
	BuildCrlSetsConcepts(uOfD, crlDataStructures, hl)
	BuildCrlListsConcepts(uOfD, crlDataStructures, hl)
	BuildCrlStringListsConcepts(uOfD, crlDataStructures, hl)
	crlDataStructures.SetIsCoreRecursively(hl)
}
