package crldatastructuresdomain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("List test", func() {
	var uOfD *core.UniverseOfDiscourse
	var hl *core.HeldLocks

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
		BuildCrlDataStructuresDomain(uOfD, hl)
		hl.ReleaseLocksAndWait()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("List should be created correctly", func() {
		Specify("Creation should fail with no specified type", func() {
			_, err := NewList(uOfD, nil, hl)
			Expect(err).Should(HaveOccurred())
		})
		Specify("Normal creation with Reference type", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, err := NewList(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			typeReference := newList.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
			Expect(typeReference).ToNot(BeNil())
			Expect(typeReference.GetReferencedConceptID(hl)).To(Equal(coreReference.GetConceptID(hl)))
			Expect(GetFirstMemberReference(newList, hl)).To(BeNil())
			Expect(GetLastMemberReference(newList, hl)).To(BeNil())
		})
	})
	Describe("AddListMemberAfter should work correctly", func() {
		Specify("Add after existing solo member should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			// Add referenceA
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, err0 := AppendListMember(newList, referenceA, hl)
			Expect(err0).ShouldNot(HaveOccurred())
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			// Add newReference
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, err := AddListMemberAfter(newList, memberReferenceA, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberReference).ToNot(BeNil())
			Expect(newMemberReference.GetReferencedConcept(hl)).To(Equal(newReference))
			// Check first member reference
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			// Check last member reference
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			// Check priorMemberReference of newMemberReference
			priorMemberReference, err3 := GetPriorMemberReference(newMemberReference, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberReference).ToNot(BeNil())
			Expect(priorMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			// Check nextMemberReference of memberReferenceA
			nextMemberReference, err4 := GetNextMemberReference(memberReferenceA, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberReference).ToNot(BeNil())
			Expect(nextMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
		})
		Specify("Add after existing first member should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			referenceB, _ := uOfD.NewReference(hl)
			memberReferenceB, _ := AppendListMember(newList, referenceB, hl)
			newReference, _ := uOfD.NewReference(hl)
			// Validate references between A and B
			nmr, _ := GetNextMemberReference(memberReferenceA, hl)
			pmr, _ := GetPriorMemberReference(memberReferenceB, hl)
			Expect(nmr.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			Expect(pmr.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			// Add new member
			newMemberReference, err := AddListMemberAfter(newList, memberReferenceA, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberReference).ToNot(BeNil())
			Expect(newMemberReference.GetReferencedConcept(hl)).To(Equal(newReference))
			// Check first member next reference
			nextMemberReferenceA, err6 := GetNextMemberReference(memberReferenceA, hl)
			Expect(err6).To(BeNil())
			Expect(nextMemberReferenceA.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			// Check new member prior reference
			priorMemberReference, err7 := GetPriorMemberReference(newMemberReference, hl)
			Expect(err7).To(BeNil())
			Expect(priorMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			// Check last member's prior refererence
			priorMemberReferenceB, err3 := GetPriorMemberReference(memberReferenceB, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberReferenceB).ToNot(BeNil())
			Expect(priorMemberReferenceB.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			// Check new member next reference
			nextMemberReference, err4 := GetNextMemberReference(newMemberReference, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberReference).ToNot(BeNil())
			Expect(nextMemberReference.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			// Check first list member
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			// Check last list member
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
		})
	})
	Describe("AddListMemberBefore should work correctly", func() {
		Specify("Add before existing solo member should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, err := AddListMemberBefore(newList, memberReferenceA, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberReference).ToNot(BeNil())
			Expect(newMemberReference.GetReferencedConcept(hl)).To(Equal(newReference))
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			priorMemberReference, err3 := GetPriorMemberReference(memberReferenceA, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberReference).ToNot(BeNil())
			Expect(priorMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			nextMemberReference, err4 := GetNextMemberReference(newMemberReference, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberReference).ToNot(BeNil())
			Expect(nextMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
		})
		Specify("Add before existing second member should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			referenceB, _ := uOfD.NewReference(hl)
			memberReferenceB, _ := AddListMemberBefore(newList, memberReferenceA, referenceB, hl)
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, err := AddListMemberBefore(newList, memberReferenceA, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberReference).ToNot(BeNil())
			Expect(newMemberReference.GetReferencedConcept(hl)).To(Equal(newReference))
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			priorMemberReference, err3 := GetPriorMemberReference(memberReferenceA, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberReference).ToNot(BeNil())
			Expect(priorMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			nextMemberReference, err4 := GetNextMemberReference(newMemberReference, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberReference).ToNot(BeNil())
			Expect(nextMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			nextMemberReferenceB, err6 := GetNextMemberReference(memberReferenceB, hl)
			Expect(err6).To(BeNil())
			Expect(nextMemberReferenceB.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			priorMemberReferenceA, err7 := GetPriorMemberReference(memberReferenceA, hl)
			Expect(err7).To(BeNil())
			Expect(priorMemberReferenceA.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
		})
	})
	Describe("AppendListMember should work correctly", func() {
		Specify("Append with empty set should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, err := AppendListMember(newList, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberReference).ToNot(BeNil())
			Expect(newMemberReference.GetReferencedConcept(hl)).To(Equal(newReference))
			Expect(newMemberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			// Check first member reference
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			// Check last member reference
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
		})
		Specify("Append with existing solo member should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			// Add referenceA
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			// Add newReference
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, err := AppendListMember(newList, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberReference).ToNot(BeNil())
			Expect(newMemberReference.GetReferencedConcept(hl)).To(Equal(newReference))
			// Check first member reference
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			// Check last member reference
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			// Check priorMemberReference of newMemberReference
			priorMemberReference, err3 := GetPriorMemberReference(newMemberReference, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberReference).ToNot(BeNil())
			Expect(priorMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			// Check nextMemberReference of memberReferenceA
			nextMemberReference, err4 := GetNextMemberReference(memberReferenceA, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberReference).ToNot(BeNil())
			Expect(nextMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
		})
	})
	Describe("ClearList should work correctrly", func() {
		Specify("Clear with empty set should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			ClearList(newList, hl)
			typeReference := newList.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
			Expect(typeReference).ToNot(BeNil())
			Expect(typeReference.GetReferencedConceptID(hl)).To(Equal(coreReference.GetConceptID(hl)))
			Expect(GetFirstMemberReference(newList, hl)).To(BeNil())
			Expect(GetLastMemberReference(newList, hl)).To(BeNil())
		})
		Specify("Clear with single set member should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, _ := AppendListMember(newList, newReference, hl)
			ClearList(newList, hl)
			typeReference := newList.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
			Expect(typeReference).ToNot(BeNil())
			Expect(typeReference.GetReferencedConceptID(hl)).To(Equal(coreReference.GetConceptID(hl)))
			Expect(GetFirstMemberReference(newList, hl)).To(BeNil())
			Expect(GetLastMemberReference(newList, hl)).To(BeNil())
			Expect(newMemberReference.GetOwningConcept(hl)).To(BeNil())
		})
	})
	Describe("GetFirstReferenceForMember should work correctly", func() {
		Specify("GetFirstReferenceForMember should find each set member", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			referenceB, _ := uOfD.NewReference(hl)
			memberReferenceB, _ := AddListMemberBefore(newList, memberReferenceA, referenceB, hl)
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, err := AddListMemberBefore(newList, memberReferenceA, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(GetFirstReferenceForMember(newList, referenceA, hl)).Should(Equal(memberReferenceA))
			Expect(GetFirstReferenceForMember(newList, referenceB, hl)).Should(Equal(memberReferenceB))
			Expect(GetFirstReferenceForMember(newList, newReference, hl)).Should(Equal(newMemberReference))
		})
		Specify("GetFirstReferenceForMember should return nil if element is not in list", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			Expect(GetFirstReferenceForMember(newList, referenceA, hl)).Should(BeNil())
		})
	})
	Describe("IsListMember should work correctly", func() {
		Specify("IsListMember on an empty set should return false", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			Expect(IsListMember(newList, referenceA, hl)).To(BeFalse())
		})
		Specify("IsListMember should work on every member of the set", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			referenceB, _ := uOfD.NewReference(hl)
			AddListMemberBefore(newList, memberReferenceA, referenceB, hl)
			newReference, _ := uOfD.NewReference(hl)
			AddListMemberBefore(newList, memberReferenceA, newReference, hl)
			Expect(IsListMember(newList, referenceA, hl)).To(BeTrue())
			Expect(IsListMember(newList, referenceB, hl)).To(BeTrue())
			Expect(IsListMember(newList, newReference, hl)).To(BeTrue())
		})
	})
	Describe("PrependListMember should work correctly", func() {
		Specify("Prepend with empty set should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, err := PrependListMember(newList, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberReference).ToNot(BeNil())
			Expect(newMemberReference.GetReferencedConcept(hl)).To(Equal(newReference))
			Expect(newMemberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			// Check first member reference
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			// Check last member reference
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
		})
		Specify("Prepend with existing solo member should work correctly", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			// Add referenceA
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := PrependListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			// Add newReference
			newReference, _ := uOfD.NewReference(hl)
			newMemberReference, err := PrependListMember(newList, newReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberReference).ToNot(BeNil())
			Expect(newMemberReference.GetReferencedConcept(hl)).To(Equal(newReference))
			// Check first member reference
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			// Check last member reference
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			// Check priorMemberReference of newMemberReference
			priorMemberReference, err3 := GetPriorMemberReference(memberReferenceA, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberReference).ToNot(BeNil())
			Expect(priorMemberReference.GetConceptID(hl)).To(Equal(newMemberReference.GetConceptID(hl)))
			// Check nextMemberReference of memberReferenceA
			nextMemberReference, err4 := GetNextMemberReference(newMemberReference, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberReference).ToNot(BeNil())
			Expect(nextMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
		})
	})
	Describe("RemoveListMember should work correctly", func() {
		Specify("RemoveListMember on empty list should return an error", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			Expect(RemoveListMember(newList, referenceA, hl)).ToNot(Succeed())
		})
		Specify("RemoveListMember on singleton list should result in the empty set", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			AppendListMember(newList, referenceA, hl)
			Expect(RemoveListMember(newList, referenceA, hl)).To(Succeed())
			typeReference := newList.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
			Expect(typeReference).ToNot(BeNil())
			Expect(typeReference.GetReferencedConceptID(hl)).To(Equal(coreReference.GetConceptID(hl)))
			Expect(GetFirstMemberReference(newList, hl)).To(BeNil())
			Expect(GetLastMemberReference(newList, hl)).To(BeNil())
			Expect(IsListMember(newList, referenceA, hl)).To(BeFalse())
		})
		Specify("RemoveListMember on first element of list should work", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			referenceB, _ := uOfD.NewReference(hl)
			memberReferenceB, _ := AddListMemberAfter(newList, memberReferenceA, referenceB, hl)
			referenceC, _ := uOfD.NewReference(hl)
			memberReferenceC, _ := AddListMemberAfter(newList, memberReferenceB, referenceC, hl)
			Expect(RemoveListMember(newList, referenceA, hl)).To(Succeed())
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(memberReferenceC.GetConceptID(hl)))
			priorMemberReference, err3 := GetPriorMemberReference(memberReferenceB, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberReference).To(BeNil())
		})
		Specify("RemoveListMember on middle element of list should work", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			referenceB, _ := uOfD.NewReference(hl)
			memberReferenceB, _ := AddListMemberAfter(newList, memberReferenceA, referenceB, hl)
			referenceC, _ := uOfD.NewReference(hl)
			memberReferenceC, _ := AddListMemberAfter(newList, memberReferenceB, referenceC, hl)
			Expect(RemoveListMember(newList, referenceB, hl)).To(Succeed())
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(memberReferenceC.GetConceptID(hl)))
			priorMemberReference, err3 := GetPriorMemberReference(memberReferenceC, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			nextMemberReference, err4 := GetNextMemberReference(memberReferenceA, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberReference.GetConceptID(hl)).To(Equal(memberReferenceC.GetConceptID(hl)))
		})
		Specify("RemoveListMember on last element of list should work", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newList, _ := NewList(uOfD, coreReference, hl)
			referenceA, _ := uOfD.NewReference(hl)
			memberReferenceA, _ := AppendListMember(newList, referenceA, hl)
			Expect(memberReferenceA.IsRefinementOfURI(CrlListMemberReferenceURI, hl)).To(BeTrue())
			referenceB, _ := uOfD.NewReference(hl)
			memberReferenceB, _ := AddListMemberAfter(newList, memberReferenceA, referenceB, hl)
			referenceC, _ := uOfD.NewReference(hl)
			AddListMemberAfter(newList, memberReferenceB, referenceC, hl)
			Expect(RemoveListMember(newList, referenceC, hl)).To(Succeed())
			firstMemberReference, err1 := GetFirstMemberReference(newList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberReference).ToNot(BeNil())
			Expect(firstMemberReference.GetConceptID(hl)).To(Equal(memberReferenceA.GetConceptID(hl)))
			lastMemberReference, err2 := GetLastMemberReference(newList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberReference).ToNot(BeNil())
			Expect(lastMemberReference.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			nextMemberReference, err3 := GetNextMemberReference(memberReferenceB, hl)
			Expect(err3).To(BeNil())
			Expect(nextMemberReference).To(BeNil())
		})
	})
	Describe("Serialization tests", func() {
		Specify("Instantiated lists should serialize and de-serialze properly", func() {
			uOfD2 := core.NewUniverseOfDiscourse()
			hl2 := uOfD.NewHeldLocks()
			BuildCrlDataStructuresDomain(uOfD2, hl)
			hl2.ReleaseLocksAndWait()
			type1 := uOfD.GetElementWithURI(core.ElementURI)
			domain1, _ := uOfD.NewElement(hl)
			list1, err0 := NewList(uOfD, type1, hl)
			list1.SetOwningConcept(domain1, hl)
			hl.ReleaseLocksAndWait()
			Expect(err0).To(BeNil())
			Expect(list1).ToNot(BeNil())
			serialized1, err := uOfD.MarshalDomain(domain1, hl)
			Expect(err).To(BeNil())
			domain2, err2 := uOfD2.RecoverDomain(serialized1, hl2)
			hl2.ReleaseLocksAndWait()
			Expect(err2).To(BeNil())
			Expect(domain2).ToNot(BeNil())
			Expect(core.RecursivelyEquivalent(domain1, hl, domain2, hl2)).To(BeTrue())
			list2 := uOfD2.GetElement(list1.GetConceptID(hl))
			Expect(list2).ToNot(BeNil())
			list1FirstElementRefRef, err3 := getListReferenceToFirstMemberReference(list1, hl)
			Expect(err3).To(BeNil())
			Expect(list1FirstElementRefRef).ToNot(BeNil())
			list1FirstElementRef, err5 := GetFirstMemberReference(list1, hl)
			Expect(err5).To(BeNil())
			Expect(list1FirstElementRef).To(BeNil())
			list2FirstElementRefRef, err4 := getListReferenceToFirstMemberReference(list2, hl2)
			Expect(err4).To(BeNil())
			Expect(list2FirstElementRefRef).ToNot(BeNil())
			list2FirstElementRef, err6 := GetFirstMemberReference(list2, hl2)
			Expect(err6).To(BeNil())
			Expect(list2FirstElementRef).To(BeNil())
		})
	})
})
