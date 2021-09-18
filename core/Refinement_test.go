package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Refinement tests", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction
	var ref Refinement
	var abstractConcept Element
	var refinedConcept Element
	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewTransaction()
		ref, _ = uOfD.NewRefinement(hl)
		abstractConcept, _ = uOfD.NewElement(hl)
		refinedConcept, _ = uOfD.NewElement(hl)
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("Setting abstract and refined concepts should work properly", func() {

		Specify("Initially abstract and refined concepts should be nil", func() {
			Expect(ref.GetAbstractConcept(hl)).To(BeNil())
			Expect(ref.GetRefinedConcept(hl)).To(BeNil())
			Expect(ref.GetAbstractConceptID(hl)).To(Equal(""))
			Expect(ref.GetRefinedConceptID(hl)).To(Equal(""))
			Expect(ref.GetAbstractConceptVersion(hl)).To(Equal(0))
			Expect(ref.GetRefinedConceptVersion(hl)).To(Equal(0))
		})
		Specify("After assignment, abstract and refined concepts should be correctly set", func() {
			abstractConcept.incrementVersion(hl)
			refinedConcept.incrementVersion(hl)
			initialVersion := ref.GetVersion(hl)
			ref.SetAbstractConceptID(abstractConcept.getConceptIDNoLock(), hl)
			Expect(ref.GetVersion(hl)).To(Equal(initialVersion + 1))
			initialVersion = ref.GetVersion(hl)
			ref.SetRefinedConceptID(refinedConcept.getConceptIDNoLock(), hl)
			Expect(ref.GetVersion(hl)).To(Equal(initialVersion + 1))
			Expect(ref.GetAbstractConcept(hl)).To(Equal(abstractConcept))
			Expect(ref.GetRefinedConcept(hl)).To(Equal(refinedConcept))
			Expect(ref.GetAbstractConceptID(hl)).To(Equal(abstractConcept.getConceptIDNoLock()))
			Expect(ref.GetRefinedConceptID(hl)).To(Equal(refinedConcept.getConceptIDNoLock()))
			Expect(ref.GetAbstractConceptVersion(hl)).To(Equal(abstractConcept.GetVersion(hl)))
			Expect(ref.GetRefinedConceptVersion(hl)).To(Equal(refinedConcept.GetVersion(hl)))
			Expect(refinedConcept.IsRefinementOf(abstractConcept, hl)).To(BeTrue())
			// Now set to nil
			ref.SetAbstractConceptID("", hl)
			ref.SetRefinedConceptID("", hl)
			Expect(ref.GetAbstractConcept(hl)).To(BeNil())
			Expect(ref.GetRefinedConcept(hl)).To(BeNil())
			Expect(ref.GetAbstractConceptID(hl)).To(Equal(""))
			Expect(ref.GetRefinedConceptID(hl)).To(Equal(""))
			Expect(ref.GetAbstractConceptVersion(hl)).To(Equal(0))
			Expect(ref.GetRefinedConceptVersion(hl)).To(Equal(0))
		})
		Specify("Setting abstract and refined concepts using actual elements should work", func() {
			Expect(ref.SetAbstractConcept(abstractConcept, hl)).To(Succeed())
			Expect(ref.GetAbstractConcept(hl)).To(Equal(abstractConcept))
			Expect(ref.SetRefinedConcept(refinedConcept, hl)).To(Succeed())
			Expect(ref.GetRefinedConcept(hl)).To(Equal(refinedConcept))
		})
		Specify("If a referenced element becomes available after it's ID is set, GetElement should find it", func() {
			ref.(*refinement).AbstractConceptID = abstractConcept.getConceptIDNoLock()
			Expect(ref.GetAbstractConcept(hl)).To(Equal(abstractConcept))
			ref.(*refinement).RefinedConceptID = refinedConcept.getConceptIDNoLock()
			Expect(ref.GetRefinedConcept(hl)).To(Equal(refinedConcept))
		})
	})

	Describe("Cloning and equivalence should work properly", func() {
		Specify("Newly initialized refinement should be equivalent to its clone", func() {
			ref, _ := uOfD.NewRefinement(hl)
			clonedRefinement := clone(ref, hl)
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeTrue())
		})
		Specify("After setting abstract element, refinement should be equivalent to its clone", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetAbstractConceptID(target.getConceptIDNoLock(), hl)
			clonedRefinement := clone(ref, hl)
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeTrue())
		})
		Specify("After setting refined element, refinement should be equivalent to its clone", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetRefinedConceptID(target.getConceptIDNoLock(), hl)
			clonedRefinement := clone(ref, hl)
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeTrue())
		})
		Specify("Equivalent should fail if there is a difference in the AbstractConceptID", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetAbstractConceptID(target.getConceptIDNoLock(), hl)
			clonedRefinement := clone(ref, hl)
			ref.(*refinement).AbstractConceptID = ""
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the RefinementdConceptID", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetRefinedConceptID(target.getConceptIDNoLock(), hl)
			clonedRefinement := clone(ref, hl)
			ref.(*refinement).RefinedConceptID = ""
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the AbstractConcept version", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetAbstractConceptID(target.getConceptIDNoLock(), hl)
			clonedRefinement := clone(ref, hl)
			ref.(*refinement).AbstractConceptVersion = 100
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the RefinedConcept version", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetRefinedConceptID(target.getConceptIDNoLock(), hl)
			clonedRefinement := clone(ref, hl)
			ref.(*refinement).RefinedConceptVersion = 100
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeFalse())
		})
		Specify("Equivalence should also fail if there is any difference in the underlying element", func() {
			ref, _ := uOfD.NewRefinement(hl)
			clonedRefinement := clone(ref, hl)
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeTrue())
			ref.(*refinement).Version.counter = 123
			Expect(Equivalent(ref, hl, clonedRefinement, hl)).To(BeFalse())
		})
	})
	Describe("Testing Marshal and Unmarshal", func() {
		Specify("After marshal and unmarshal the recovered refinement should be equivalent to the original", func() {
			ref, _ := uOfD.NewRefinement(hl)
			ac, _ := uOfD.NewElement(hl)
			rc, _ := uOfD.NewElement(hl)
			ref.SetAbstractConceptID(ac.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(rc.getConceptIDNoLock(), hl)
			mRef, err1 := ref.MarshalJSON()
			Expect(err1).To(BeNil())
			mAc, err2 := ac.MarshalJSON()
			Expect(err2).To(BeNil())
			mRc, err3 := rc.MarshalJSON()
			Expect(err3).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewTransaction()
			rRef, err4 := uOfD2.RecoverElement(mRef, hl2)
			Expect(err4).To(BeNil())
			rAc, err5 := uOfD2.RecoverElement(mAc, hl2)
			Expect(err5).To(BeNil())
			rRc, err6 := uOfD2.RecoverElement(mRc, hl2)
			Expect(err6).To(BeNil())
			Expect(Equivalent(ref, hl, rRef, hl2)).To(BeTrue())
			Expect(Equivalent(ac, hl, rAc, hl2)).To(BeTrue())
			Expect(Equivalent(rc, hl, rRc, hl2)).To(BeTrue())
		})
	})
})
