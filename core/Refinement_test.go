package core

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Refinement tests", func() {
	var uOfD *UniverseOfDiscourse
	var trans *Transaction
	var ref *Concept
	var abstractConcept *Concept
	var refinedConcept *Concept
	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
		ref, _ = uOfD.NewRefinement(trans)
		abstractConcept, _ = uOfD.NewElement(trans)
		refinedConcept, _ = uOfD.NewElement(trans)
	})

	AfterEach(func() {
		trans.ReleaseLocks()
	})

	Describe("Setting abstract and refined concepts should work properly", func() {

		Specify("Initially abstract and refined concepts should be nil", func() {
			Expect(ref.GetAbstractConcept(trans)).To(BeNil())
			Expect(ref.GetRefinedConcept(trans)).To(BeNil())
			Expect(ref.GetAbstractConceptID(trans)).To(Equal(""))
			Expect(ref.GetRefinedConceptID(trans)).To(Equal(""))
		})
		Specify("After assignment, abstract and refined concepts should be correctly set", func() {
			abstractConcept.Version.incrementVersion()
			refinedConcept.Version.incrementVersion()
			initialVersion := ref.GetVersion(trans)
			ref.SetAbstractConceptID(abstractConcept.getConceptIDNoLock(), trans)
			Expect(ref.GetVersion(trans)).To(Equal(initialVersion + 1))
			initialVersion = ref.GetVersion(trans)
			ref.SetRefinedConceptID(refinedConcept.getConceptIDNoLock(), trans)
			Expect(ref.GetVersion(trans)).To(Equal(initialVersion + 1))
			Expect(ref.GetAbstractConcept(trans)).To(Equal(abstractConcept))
			Expect(ref.GetRefinedConcept(trans)).To(Equal(refinedConcept))
			Expect(ref.GetAbstractConceptID(trans)).To(Equal(abstractConcept.getConceptIDNoLock()))
			Expect(ref.GetRefinedConceptID(trans)).To(Equal(refinedConcept.getConceptIDNoLock()))
			Expect(refinedConcept.IsRefinementOf(abstractConcept, trans)).To(BeTrue())
			// Now set to nil
			ref.SetAbstractConceptID("", trans)
			ref.SetRefinedConceptID("", trans)
			Expect(ref.GetAbstractConcept(trans)).To(BeNil())
			Expect(ref.GetRefinedConcept(trans)).To(BeNil())
			Expect(ref.GetAbstractConceptID(trans)).To(Equal(""))
			Expect(ref.GetRefinedConceptID(trans)).To(Equal(""))
		})
		Specify("Setting abstract and refined concepts using actual elements should work", func() {
			Expect(ref.SetAbstractConcept(abstractConcept, trans)).To(Succeed())
			Expect(ref.GetAbstractConcept(trans)).To(Equal(abstractConcept))
			Expect(ref.SetRefinedConcept(refinedConcept, trans)).To(Succeed())
			Expect(ref.GetRefinedConcept(trans)).To(Equal(refinedConcept))
		})
		Specify("If a referenced element becomes available after it's ID is set, GetElement should find it", func() {
			ref.AbstractConceptID = abstractConcept.getConceptIDNoLock()
			Expect(ref.GetAbstractConcept(trans)).To(Equal(abstractConcept))
			ref.RefinedConceptID = refinedConcept.getConceptIDNoLock()
			Expect(ref.GetRefinedConcept(trans)).To(Equal(refinedConcept))
		})
	})

	Describe("Cloning and equivalence should work properly", func() {
		Specify("Newly initialized refinement should be equivalent to its clone", func() {
			ref, _ := uOfD.NewRefinement(trans)
			clonedRefinement := clone(ref, trans)
			Expect(Equivalent(ref, trans, clonedRefinement, trans)).To(BeTrue())
		})
		Specify("After setting abstract element, refinement should be equivalent to its clone", func() {
			ref, _ := uOfD.NewRefinement(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetAbstractConceptID(target.getConceptIDNoLock(), trans)
			clonedRefinement := clone(ref, trans)
			Expect(Equivalent(ref, trans, clonedRefinement, trans)).To(BeTrue())
		})
		Specify("After setting refined element, refinement should be equivalent to its clone", func() {
			ref, _ := uOfD.NewRefinement(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetRefinedConceptID(target.getConceptIDNoLock(), trans)
			clonedRefinement := clone(ref, trans)
			Expect(Equivalent(ref, trans, clonedRefinement, trans)).To(BeTrue())
		})
		Specify("Equivalent should fail if there is a difference in the AbstractConceptID", func() {
			ref, _ := uOfD.NewRefinement(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetAbstractConceptID(target.getConceptIDNoLock(), trans)
			clonedRefinement := clone(ref, trans)
			ref.AbstractConceptID = ""
			Expect(Equivalent(ref, trans, clonedRefinement, trans)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the RefinementdConceptID", func() {
			ref, _ := uOfD.NewRefinement(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetRefinedConceptID(target.getConceptIDNoLock(), trans)
			clonedRefinement := clone(ref, trans)
			ref.RefinedConceptID = ""
			Expect(Equivalent(ref, trans, clonedRefinement, trans)).To(BeFalse())
		})
		Specify("Equivalence should also fail if there is any difference in the underlying element", func() {
			ref, _ := uOfD.NewRefinement(trans)
			clonedRefinement := clone(ref, trans)
			Expect(Equivalent(ref, trans, clonedRefinement, trans)).To(BeTrue())
			ref.Version.counter = 123
			Expect(Equivalent(ref, trans, clonedRefinement, trans)).To(BeFalse())
		})
	})
	Describe("Testing Marshal and Unmarshal", func() {
		Specify("After marshal and unmarshal the recovered refinement should be equivalent to the original", func() {
			ref, _ := uOfD.NewRefinement(trans)
			ac, _ := uOfD.NewElement(trans)
			rc, _ := uOfD.NewElement(trans)
			ref.SetAbstractConceptID(ac.getConceptIDNoLock(), trans)
			ref.SetRefinedConceptID(rc.getConceptIDNoLock(), trans)
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
			Expect(Equivalent(ref, trans, rRef, hl2)).To(BeTrue())
			Expect(Equivalent(ac, trans, rAc, hl2)).To(BeTrue())
			Expect(Equivalent(rc, trans, rRc, hl2)).To(BeTrue())
		})
	})
})
