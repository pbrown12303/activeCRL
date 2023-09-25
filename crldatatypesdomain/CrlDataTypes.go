package crldatatypesdomain

import (
	"github.com/pbrown12303/activeCRL/core"
)

// CrlDataTypesDomainURI is the URI for the concpet space that defines the CRL Data Types
var CrlDataTypesDomainURI = "http://activeCRL.com/crldatastructuresdomain/CrlDataTypes"

// BuildCrlDataTypesDomain constructs the concept space for CRL data structures
func BuildCrlDataTypesDomain(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) {
	crlDataTypes, _ := uOfD.NewElement(trans, CrlDataTypesDomainURI)
	crlDataTypes.SetLabel("CrlDataTypesDomain", trans)
	BuildCrlBooleanConcept(uOfD, crlDataTypes, trans)
	crlDataTypes.SetReadOnlyRecursively(true, trans)
	crlDataTypes.SetIsCoreRecursively(trans)
}
