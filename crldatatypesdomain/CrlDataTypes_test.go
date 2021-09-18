package crldatatypesdomain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("CrlDataTypes test", func() {
	Specify("Domain cration should be idempotent", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewTransaction()
		BuildCrlDataTypesDomain(uOfD1, hl1)
		cs1 := uOfD1.GetElementWithURI(CrlDataTypesDomainURI)
		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewTransaction()
		BuildCrlDataTypesDomain(uOfD2, hl2)
		cs2 := uOfD2.GetElementWithURI(CrlDataTypesDomainURI)
		Expect(core.RecursivelyEquivalent(cs1, hl1, cs2, hl2)).To(BeTrue())
	})
})
