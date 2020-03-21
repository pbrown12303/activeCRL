package crldatatypes

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("CrlDataTypes test", func() {
	Specify("Domain cration should be idempotent", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewHeldLocks()
		BuildCrlDataTypesConceptSpace(uOfD1, hl1)
		cs1 := uOfD1.GetElementWithURI(CrlDataTypesConceptSpaceURI)
		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewHeldLocks()
		BuildCrlDataTypesConceptSpace(uOfD2, hl2)
		cs2 := uOfD2.GetElementWithURI(CrlDataTypesConceptSpaceURI)
		Expect(core.RecursivelyEquivalent(cs1, hl1, cs2, hl2)).To(BeTrue())
	})
})
