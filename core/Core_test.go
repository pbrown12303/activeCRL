package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core should build properly", func() {
	Specify("All core concepts should be present and be marked as core", func() {
		uOfD := NewUniverseOfDiscourse()
		hl := uOfD.NewHeldLocks()
		core := uOfD.GetElementWithURI(CoreConceptSpaceURI)
		Expect(core).ToNot(BeNil())
		Expect(core.GetIsCore(hl)).To(BeTrue())
		el := uOfD.GetElementWithURI(ElementURI)
		Expect(el).ToNot(BeNil())
		Expect(el.GetIsCore(hl)).To(BeTrue())
		literal := uOfD.GetElementWithURI(LiteralURI)
		Expect(literal).ToNot(BeNil())
		Expect(literal.GetIsCore(hl)).To(BeTrue())
		reference := uOfD.GetElementWithURI(ReferenceURI)
		Expect(reference).ToNot(BeNil())
		Expect(reference.GetIsCore(hl)).To(BeTrue())
		refinement := uOfD.GetElementWithURI(RefinementURI)
		Expect(refinement).ToNot(BeNil())
		Expect(refinement.GetIsCore(hl)).To(BeTrue())
	})
})
