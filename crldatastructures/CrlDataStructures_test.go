package crldatastructures

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("CrlDataStructures domain test", func() {
	Specify("Domain generation should be idempotent", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewHeldLocks()
		BuildCrlDataStructuresConceptSpace(uOfD1, hl1)
		cs1 := uOfD1.GetElementWithURI(CrlDataStructuresConceptSpaceURI)
		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewHeldLocks()
		BuildCrlDataStructuresConceptSpace(uOfD2, hl2)
		cs2 := uOfD2.GetElementWithURI(CrlDataStructuresConceptSpaceURI)
		Expect(core.RecursivelyEquivalent(cs1, hl1, cs2, hl2)).To(BeTrue())
	})
})
