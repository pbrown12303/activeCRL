package crleditorbrowserguidomain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructuresdomain"
	"github.com/pbrown12303/activeCRL/crldiagram"
)

var _ = Describe("BrowserGUIDomain tests", func() {
	Specify("BrowserGUI domain creation should be idempotent", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewHeldLocks()
		crldiagram.BuildCrlDiagramDomain(uOfD1, hl1)
		crldatastructuresdomain.BuildCrlDataStructuresDomain(uOfD1, hl1)
		Expect(BuildBrowserGUIDomain(uOfD1, hl1)).ShouldNot(BeNil())
		cs1 := uOfD1.GetElementWithURI(BrowserGUIDomainURI)
		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewHeldLocks()
		crldiagram.BuildCrlDiagramDomain(uOfD2, hl2)
		crldatastructuresdomain.BuildCrlDataStructuresDomain(uOfD2, hl2)
		Expect(BuildBrowserGUIDomain(uOfD2, hl2)).ShouldNot(BeNil())
		cs2 := uOfD2.GetElementWithURI(BrowserGUIDomainURI)
		Expect(core.RecursivelyEquivalent(cs1, hl1, cs2, hl2, true)).To(BeTrue())
	})
})
