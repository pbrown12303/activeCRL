package core

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core should build properly", func() {
	Specify("All core concepts should be present and be marked as core", func() {
		uOfD := NewUniverseOfDiscourse()
		trans := uOfD.NewTransaction()
		core := uOfD.GetElementWithURI(CoreDomainURI)
		Expect(core).ToNot(BeNil())
		Expect(core.GetIsCore(trans)).To(BeTrue())
		el := uOfD.GetElementWithURI(ElementURI)
		Expect(el).ToNot(BeNil())
		Expect(el.GetIsCore(trans)).To(BeTrue())
		literal := uOfD.GetElementWithURI(LiteralURI)
		Expect(literal).ToNot(BeNil())
		Expect(literal.GetIsCore(trans)).To(BeTrue())
		reference := uOfD.GetElementWithURI(ReferenceURI)
		Expect(reference).ToNot(BeNil())
		Expect(reference.GetIsCore(trans)).To(BeTrue())
		refinement := uOfD.GetElementWithURI(RefinementURI)
		Expect(refinement).ToNot(BeNil())
		Expect(refinement.GetIsCore(trans)).To(BeTrue())
	})
	Specify("The creation of the core domain should be idempotent", func() {
		uOfD1 := NewUniverseOfDiscourse()
		hl1 := uOfD1.NewTransaction()
		cs1 := uOfD1.GetElementWithURI(CoreDomainURI)
		uOfD2 := NewUniverseOfDiscourse()
		hl2 := uOfD2.NewTransaction()
		cs2 := uOfD2.GetElementWithURI(CoreDomainURI)
		Expect(RecursivelyEquivalent(cs1, hl1, cs2, hl2)).To(BeTrue())
	})
})
