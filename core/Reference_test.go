package core

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reference Tests", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewTransaction()
	})

	AfterEach(func() {
		hl.ReleaseLocks()
	})

	Describe("Setting and getting the ReferencedConcept", func() {
		Specify("Referenced concept should initially be nil", func() {
			ref, _ := uOfD.NewReference(hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(""))
			Expect(ref.GetReferencedConcept(hl)).To(BeNil())
			Expect(ref.getReferencedConceptNoLock()).To(BeNil())
			Expect(ref.GetReferencedConceptVersion(hl)).To(Equal(0))
		})
		Specify("Referenced concept should set correctly", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			target.(*element).Version.counter = 66
			initialVersion := ref.GetVersion(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(target.getConceptIDNoLock()))
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
			Expect(ref.getReferencedConceptNoLock()).To(Equal(target))
			Expect(ref.GetReferencedConceptVersion(hl)).To(Equal(target.GetVersion(hl)))
			Expect(ref.GetVersion(hl)).To(Equal(initialVersion + 1))
		})
		Specify("Referenced concept should clear correctly", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			target.(*element).Version.counter = 66
			initialVersion := ref.GetVersion(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, hl)
			Expect(ref.GetVersion(hl)).To(Equal(initialVersion + 1))
			ref.SetReferencedConceptID("", NoAttribute, hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(""))
			Expect(ref.GetReferencedConcept(hl)).To(BeNil())
			Expect(ref.getReferencedConceptNoLock()).To(BeNil())
			Expect(ref.GetReferencedConceptVersion(hl)).To(Equal(0))
		})
		Specify("SetReferencedConcept should work correctly", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			initialVersion := ref.GetVersion(hl)
			ref.SetReferencedConcept(target, NoAttribute, hl)
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
			Expect(ref.GetVersion(hl)).To(Equal(initialVersion + 1))
		})
		Specify("Referenced element should be retrieved from uOfD if cache does not contain pointer", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.(*reference).ReferencedConceptID = target.getConceptIDNoLock()
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
		})
	})

	Describe("Setting and getting the ReferencedConcept AttributeName", func() {
		Specify("ReferencedConcept AttributeName should initially be the empty string", func() {
			ref, _ := uOfD.NewReference(hl)
			Expect(ref.GetReferencedAttributeName(hl)).To(Equal(NoAttribute))
		})
		Specify("ReferencedConcept AttributeName should set correctly", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			initialVersion := ref.GetVersion(hl)
			ref.SetReferencedConcept(target, OwningConceptID, hl)
			Expect(ref.GetVersion(hl)).To(Equal(initialVersion + 1))
			Expect(ref.GetReferencedAttributeName(hl)).To(Equal(OwningConceptID))
			ref.SetReferencedConcept(target, NoAttribute, hl)
			Expect(ref.GetReferencedAttributeName(hl)).To(Equal(NoAttribute))
		})
	})

	Describe("Ensure that the read-only setting prevents setting the referenced concept", func() {
		Specify("SetReferencedConceptID should fail if read-only is set", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReadOnly(true, hl)
			Expect(ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, hl)).ToNot(Succeed())
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
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, hl)
			clonedReference := clone(ref, hl)
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeTrue())
		})
		Specify("Equivalent should fail if there is a difference in the ReferencedConceptID", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, hl)
			clonedReference := clone(ref, hl)
			ref.(*reference).ReferencedConceptID = ""
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the ReferencedConceptAttributeName", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConcept(target, OwningConceptID, hl)
			clonedReference := clone(ref, hl)
			ref.SetReferencedConcept(target, ReferencedConceptID, hl)
			Expect(Equivalent(ref, hl, clonedReference, hl)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the ReferencedConcept version", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, hl)
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
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), OwningConceptID, hl)
			mRef, err1 := ref.MarshalJSON()
			Expect(err1).To(BeNil())
			mTarget, err3 := target.MarshalJSON()
			Expect(err3).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewTransaction()
			rRef, err2 := uOfD2.RecoverElement(mRef, hl2)
			Expect(err2).To(BeNil())
			rTarget, err4 := uOfD2.RecoverElement(mTarget, hl2)
			Expect(err4).To(BeNil())
			Expect(Equivalent(ref, hl, rRef, hl2)).To(BeTrue())
			Expect(Equivalent(target, hl, rTarget, hl2)).To(BeTrue())
		})
	})
})

var _ = Describe("Test FindAttributeName", func() {
	Specify("The correct values should be returned", func() {
		foundAttribute, err := FindAttributeName("NoAttribute")
		Expect(err).To(BeNil())
		Expect(foundAttribute).To(Equal(NoAttribute))
		foundAttribute, err = FindAttributeName("OwningConceptID")
		Expect(err).To(BeNil())
		Expect(foundAttribute).To(Equal(OwningConceptID))
		foundAttribute, err = FindAttributeName("ReferencedConceptID")
		Expect(err).To(BeNil())
		Expect(foundAttribute).To(Equal(ReferencedConceptID))
		foundAttribute, err = FindAttributeName("AbstractConceptID")
		Expect(err).To(BeNil())
		Expect(foundAttribute).To(Equal(AbstractConceptID))
		foundAttribute, err = FindAttributeName("RefinedConceptID")
		Expect(err).To(BeNil())
		Expect(foundAttribute).To(Equal(RefinedConceptID))
		foundAttribute, err = FindAttributeName("garbage")
		Expect(err).ToNot(BeNil())
		Expect(foundAttribute).To(Equal(NoAttribute))
	})
})
