package crldatatypesdomain

import (
	"github.com/pbrown12303/activeCRL/core"
)

// CrlDataTypesDomainURI is the URI for the concpet space that defines the CRL Data Types
var CrlDataTypesDomainURI = "http://activeCRL.com/crldatastructuresdomain/CrlDataTypes"

// BuildCrlDataTypesDomain constructs the concept space for CRL data structures
func BuildCrlDataTypesDomain(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {
	crlDataTypes, _ := uOfD.NewElement(hl, CrlDataTypesDomainURI)
	crlDataTypes.SetLabel("CrlDataTypesDomain", hl)
	BuildCrlBooleanConcept(uOfD, crlDataTypes, hl)
	crlDataTypes.SetReadOnlyRecursively(true, hl)
	crlDataTypes.SetIsCoreRecursively(hl)
}
