package core

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reference Tests", func() {
	var uOfD *UniverseOfDiscourse
	var trans *Transaction

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
	})

	AfterEach(func() {
		trans.ReleaseLocks()
	})

	Describe("Setting and getting the ReferencedConcept", func() {
		Specify("Referenced concept should initially be nil", func() {
			ref, _ := uOfD.NewReference(trans)
			Expect(ref.GetReferencedConceptID(trans)).To(Equal(""))
			Expect(ref.GetReferencedConcept(trans)).To(BeNil())
			Expect(ref.getReferencedConceptNoLock()).To(BeNil())
		})
		Specify("Referenced concept should set correctly", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			target.(*concept).Version.counter = 66
			initialVersion := ref.GetVersion(trans)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, trans)
			Expect(ref.GetReferencedConceptID(trans)).To(Equal(target.getConceptIDNoLock()))
			Expect(ref.GetReferencedConcept(trans)).To(Equal(target))
			Expect(ref.getReferencedConceptNoLock()).To(Equal(target))
			Expect(ref.GetVersion(trans)).To(Equal(initialVersion + 1))
		})
		Specify("Referenced concept should clear correctly", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			target.(*concept).Version.counter = 66
			initialVersion := ref.GetVersion(trans)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, trans)
			Expect(ref.GetVersion(trans)).To(Equal(initialVersion + 1))
			ref.SetReferencedConceptID("", NoAttribute, trans)
			Expect(ref.GetReferencedConceptID(trans)).To(Equal(""))
			Expect(ref.GetReferencedConcept(trans)).To(BeNil())
			Expect(ref.getReferencedConceptNoLock()).To(BeNil())
		})
		Specify("SetReferencedConcept should work correctly", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			initialVersion := ref.GetVersion(trans)
			ref.SetReferencedConcept(target, NoAttribute, trans)
			Expect(ref.GetReferencedConcept(trans)).To(Equal(target))
			Expect(ref.GetVersion(trans)).To(Equal(initialVersion + 1))
		})
		Specify("Referenced element should be retrieved from uOfD if cache does not contain pointer", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			ref.(*concept).ReferencedConceptID = target.getConceptIDNoLock()
			Expect(ref.GetReferencedConcept(trans)).To(Equal(target))
		})
	})

	Describe("Setting and getting the ReferencedConcept AttributeName", func() {
		Specify("ReferencedConcept AttributeName should initially be the empty string", func() {
			ref, _ := uOfD.NewReference(trans)
			Expect(ref.GetReferencedAttributeName(trans)).To(Equal(NoAttribute))
		})
		Specify("ReferencedConcept AttributeName should set correctly", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			initialVersion := ref.GetVersion(trans)
			ref.SetReferencedConcept(target, OwningConceptID, trans)
			Expect(ref.GetVersion(trans)).To(Equal(initialVersion + 1))
			Expect(ref.GetReferencedAttributeName(trans)).To(Equal(OwningConceptID))
			ref.SetReferencedConcept(target, NoAttribute, trans)
			Expect(ref.GetReferencedAttributeName(trans)).To(Equal(NoAttribute))
		})
	})

	Describe("Ensure that the read-only setting prevents setting the referenced concept", func() {
		Specify("SetReferencedConceptID should fail if read-only is set", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetReadOnly(true, trans)
			Expect(ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, trans)).ToNot(Succeed())
		})
	})

	Describe("Test clone and equivalence", func() {
		Specify("Newly initialized reference should be equivalent to its clone", func() {
			ref, _ := uOfD.NewReference(trans)
			clonedReference := clone(ref, trans)
			Expect(Equivalent(ref, trans, clonedReference, trans)).To(BeTrue())
		})
		Specify("After setting referenced element, reference should be equivalent to its clone", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, trans)
			clonedReference := clone(ref, trans)
			Expect(Equivalent(ref, trans, clonedReference, trans)).To(BeTrue())
		})
		Specify("Equivalent should fail if there is a difference in the ReferencedConceptID", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, trans)
			clonedReference := clone(ref, trans)
			ref.(*concept).ReferencedConceptID = ""
			Expect(Equivalent(ref, trans, clonedReference, trans)).To(BeFalse())
		})
		Specify("Equivalent should fail if there is a difference in the ReferencedConceptAttributeName", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetReferencedConcept(target, OwningConceptID, trans)
			clonedReference := clone(ref, trans)
			ref.SetReferencedConcept(target, ReferencedConceptID, trans)
			Expect(Equivalent(ref, trans, clonedReference, trans)).To(BeFalse())
		})
		Specify("Equivalence should also fail if there is any difference in the underlying element", func() {
			ref, _ := uOfD.NewReference(trans)
			clonedReference := clone(ref, trans)
			Expect(Equivalent(ref, trans, clonedReference, trans)).To(BeTrue())
			ref.(*concept).Version.counter = 123
			Expect(Equivalent(ref, trans, clonedReference, trans)).To(BeFalse())
		})
	})
	Describe("Marshal and Unmarshal Test", func() {
		Specify("Original and unmarshaled version should be equivalent", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), OwningConceptID, trans)
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
			Expect(Equivalent(ref, trans, rRef, hl2)).To(BeTrue())
			Expect(Equivalent(target, trans, rTarget, hl2)).To(BeTrue())
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
