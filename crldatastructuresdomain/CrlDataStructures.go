package crldatastructuresdomain

import (
	"github.com/pbrown12303/activeCRL/core"
)

// CrlDataStructuresDomainURI is the uri for the concept space that defines the Crl Data Structures
var CrlDataStructuresDomainURI = "http://activeCRL.com/crldatastructuresdomain/CrlDataStructuresDomain"

// BuildCrlDataStructuresDomain constructs the concept space for CRL data structures
func BuildCrlDataStructuresDomain(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) {
	crlDataStructures, _ := uOfD.NewElement(trans, CrlDataStructuresDomainURI)
	crlDataStructures.SetLabel("CrlDataStructuresDomain", trans)
	BuildCrlSetsConcepts(uOfD, crlDataStructures, trans)
	BuildCrlListsConcepts(uOfD, crlDataStructures, trans)
	BuildCrlStringListsConcepts(uOfD, crlDataStructures, trans)
	crlDataStructures.SetIsCoreRecursively(trans)
}
