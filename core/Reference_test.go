package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reference Tests", func() {
	var uOfD UniverseOfDiscourse
	var hl *HeldLocks

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("Setting and getting the ReferencedElement", func() {
		Specify("Referenced element should initially be nil", func() {
			ref, _ := uOfD.NewReference(hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(""))
			Expect(ref.GetReferencedConcept(hl)).To(BeNil())
			Expect(ref.getReferencedConceptNoLock()).To(BeNil())
			Expect(ref.GetReferencedConceptVersion(hl)).To(Equal(0))
		})
		Specify("Referenced element should set correctly", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			target.(*element).Version.counter = 66
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(target.getConceptIDNoLock()))
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
			Expect(ref.getReferencedConceptNoLock()).To(Equal(target))
			Expect(ref.GetReferencedConceptVersion(hl)).To(Equal(target.GetVersion(hl)))
		})
		Specify("Referenced element should clear correctly", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			target.(*element).Version.counter = 66
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)
			ref.SetReferencedConceptID("", hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(""))
			Expect(ref.GetReferencedConcept(hl)).To(BeNil())
			Expect(ref.getReferencedConceptNoLock()).To(BeNil())
			Expect(ref.GetReferencedConceptVersion(hl)).To(Equal(0))
		})
		Specify("SetReferencedConcept should work correctly", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConcept(target, hl)
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
		})
		Specify("Referenced element should be retrieved from uOfD if cache does not contain pointer", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.(*reference).ReferencedConceptID = target.getConceptIDNoLock()
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
		})
	})

	Describe("Ensure that the read-only setting prevents setting the referenced concept", func() {
		Specify("SetReferencedConceptID should fail if read-only is set", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReadOnly(true, hl)
			Expect(ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)).ToNot(Succeed())
		})
	})

	Describe("Test clone and equivalence", func() {
		Specify("Newly initialized reference should be equivalent to its clone", func() {
			ref, _ := uOfD.NewReference(hl)
			clonedReference := clone(ref, hl)
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeTrue())
		})
		Specify("After setting referenced element, reference should be equivalent to its clone", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)
			clonedReference := clone(ref, hl)
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeTrue())
		})
		Specify("Equivalent should fail if there is a difference in the ReferencedConceptID", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)
			clonedReference := clone(ref, hl)
			ref.(*reference).ReferencedConceptID = ""
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the ReferencedConcept", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)
			clonedReference := clone(ref, hl)
			ref.(*reference).referencedConcept.indicatedConcept = nil
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the ReferencedConcept version", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)
			clonedReference := clone(ref, hl)
			ref.(*reference).ReferencedConceptVersion = 100
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeFalse())
		})
		Specify("Equivalence should also fail if there is any difference in the underlying element", func() {
			ref, _ := uOfD.NewReference(hl)
			clonedReference := clone(ref, hl)
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeTrue())
			ref.(*reference).Version.counter = 123
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeFalse())
		})
	})
	Describe("Marshal and Unmarshal Test", func() {
		Specify("Original and unmarshaled version should be equivalent", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)
			mRef, err1 := ref.MarshalJSON()
			Expect(err1).To(BeNil())
			mTarget, err3 := target.MarshalJSON()
			Expect(err3).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewHeldLocks()
			rRef, err2 := uOfD2.RecoverElement(mRef, hl2)
			Expect(err2).To(BeNil())
			rTarget, err4 := uOfD2.RecoverElement(mTarget, hl2)
			Expect(err4).To(BeNil())
			Expect(Equivalent(ref, hl, rRef, hl2)).To(BeTrue())
			Expect(Equivalent(target, hl, rTarget, hl2)).To(BeTrue())
		})
	})
})
