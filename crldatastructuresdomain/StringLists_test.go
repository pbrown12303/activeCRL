package crldatastructuresdomain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("StringList test", func() {
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

	Describe("StringList should be created correctly", func() {
		Specify("Normal creation", func() {
			newStringList, err := NewStringList(uOfD, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(GetFirstMemberLiteral(newStringList, hl)).To(BeNil())
			Expect(GetLastMemberLiteral(newStringList, hl)).To(BeNil())
		})
	})
	Describe("AddStringListMemberAfter should work correctly", func() {
		Specify("Add after existing solo member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			// Add valueA
			valueA := "A"
			memberLiteralA, err0 := AppendStringListMember(newStringList, valueA, hl)
			Expect(err0).ShouldNot(HaveOccurred())
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			// Add newValue
			newValue := "NewValue"
			newMemberLiteral, err := AddStringListMemberAfter(newStringList, memberLiteralA, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(hl)).To(Equal(newValue))
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			// Check priorMemberLiteral of newMemberLiteral
			priorMemberLiteral, err3 := GetPriorMemberLiteral(newMemberLiteral, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			// Check nextMemberLiteral of memberLiteralA
			nextMemberLiteral, err4 := GetNextMemberLiteral(memberLiteralA, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
		})
		Specify("Add after existing first member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AppendStringListMember(newStringList, valueB, hl)
			newValue := "NewValue"
			// Validate references between A and B
			nmr, _ := GetNextMemberLiteral(memberLiteralA, hl)
			pmr, _ := GetPriorMemberLiteral(memberReferenceB, hl)
			Expect(nmr.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			Expect(pmr.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			// Add new member
			newMemberLiteral, err := AddStringListMemberAfter(newStringList, memberLiteralA, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(hl)).To(Equal(newValue))
			// Check first member next reference
			nextMemberLiteralA, err6 := GetNextMemberLiteral(memberLiteralA, hl)
			Expect(err6).To(BeNil())
			Expect(nextMemberLiteralA.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			// Check new member prior reference
			priorMemberLiteral, err7 := GetPriorMemberLiteral(newMemberLiteral, hl)
			Expect(err7).To(BeNil())
			Expect(priorMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			// Check last member's prior refererence
			priorMemberLiteralB, err3 := GetPriorMemberLiteral(memberReferenceB, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteralB).ToNot(BeNil())
			Expect(priorMemberLiteralB.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			// Check new member next reference
			nextMemberLiteral, err4 := GetNextMemberLiteral(newMemberLiteral, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			// Check first list member
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			// Check last list member
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
		})
	})
	Describe("AddStringListMemberBefore should work correctly", func() {
		Specify("Add before existing solo member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			newValue := "NewValue"
			newMemberLiteral, err := AddStringListMemberBefore(newStringList, memberLiteralA, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(hl)).To(Equal(newValue))
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberLiteralA, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			nextMemberLiteral, err4 := GetNextMemberLiteral(newMemberLiteral, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
		})
		Specify("Add before existing second member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberBefore(newStringList, memberLiteralA, valueB, hl)
			newValue := "NewValue"
			newMemberLiteral, err := AddStringListMemberBefore(newStringList, memberLiteralA, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(hl)).To(Equal(newValue))
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberLiteralA, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			nextMemberLiteral, err4 := GetNextMemberLiteral(newMemberLiteral, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			nextMemberLiteralB, err6 := GetNextMemberLiteral(memberReferenceB, hl)
			Expect(err6).To(BeNil())
			Expect(nextMemberLiteralB.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			priorMemberLiteralA, err7 := GetPriorMemberLiteral(memberLiteralA, hl)
			Expect(err7).To(BeNil())
			Expect(priorMemberLiteralA.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
		})
	})
	Describe("AppendStringListMember should work correctly", func() {
		Specify("Append with empty set should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			newValue := "NewValue"
			newMemberLiteral, err := AppendStringListMember(newStringList, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(hl)).To(Equal(newValue))
			Expect(newMemberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
		})
		Specify("Append with existing solo member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			// Add valueA
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			// Add newValue
			newValue := "NewValue"
			newMemberLiteral, err := AppendStringListMember(newStringList, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(hl)).To(Equal(newValue))
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			// Check priorMemberLiteral of newMemberLiteral
			priorMemberLiteral, err3 := GetPriorMemberLiteral(newMemberLiteral, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			// Check nextMemberLiteral of memberLiteralA
			nextMemberLiteral, err4 := GetNextMemberLiteral(memberLiteralA, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
		})
	})
	Describe("ClearStringList should work correctrly", func() {
		Specify("Clear with empty set should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			ClearStringList(newStringList, hl)
			Expect(GetFirstMemberLiteral(newStringList, hl)).To(BeNil())
			Expect(GetLastMemberLiteral(newStringList, hl)).To(BeNil())
		})
		Specify("Clear with single set member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			newValue := "NewValue"
			newMemberLiteral, _ := AppendStringListMember(newStringList, newValue, hl)
			ClearStringList(newStringList, hl)
			Expect(GetFirstMemberLiteral(newStringList, hl)).To(BeNil())
			Expect(GetLastMemberLiteral(newStringList, hl)).To(BeNil())
			Expect(newMemberLiteral.GetOwningConcept(hl)).To(BeNil())
		})
	})
	Describe("GetFirstLiteralForString should work correctly", func() {
		Specify("GetFirstLiteralForString should find each set member", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberBefore(newStringList, memberLiteralA, valueB, hl)
			newValue := "NewValue"
			newMemberLiteral, err := AddStringListMemberBefore(newStringList, memberLiteralA, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(GetFirstLiteralForString(newStringList, valueA, hl)).Should(Equal(memberLiteralA))
			Expect(GetFirstLiteralForString(newStringList, valueB, hl)).Should(Equal(memberReferenceB))
			Expect(GetFirstLiteralForString(newStringList, newValue, hl)).Should(Equal(newMemberLiteral))
		})
		Specify("GetFirstLiteralForString should return nil if element is not in list", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			Expect(GetFirstLiteralForString(newStringList, valueA, hl)).Should(BeNil())
		})
	})
	Describe("IsStringListMember should work correctly", func() {
		Specify("IsStringListMember on an empty set should return false", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			Expect(IsStringListMember(newStringList, valueA, hl)).To(BeFalse())
		})
		Specify("IsStringListMember should work on every member of the set", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			valueB := "B"
			AddStringListMemberBefore(newStringList, memberLiteralA, valueB, hl)
			newValue := "NewValue"
			AddStringListMemberBefore(newStringList, memberLiteralA, newValue, hl)
			Expect(IsStringListMember(newStringList, valueA, hl)).To(BeTrue())
			Expect(IsStringListMember(newStringList, valueB, hl)).To(BeTrue())
			Expect(IsStringListMember(newStringList, newValue, hl)).To(BeTrue())
		})
	})
	Describe("PrependStringListMember should work correctly", func() {
		Specify("Prepend with empty set should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			newValue := "NewValue"
			newMemberLiteral, err := PrependStringListMember(newStringList, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(hl)).To(Equal(newValue))
			Expect(newMemberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
		})
		Specify("Prepend with existing solo member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			// Add valueA
			valueA := "A"
			memberLiteralA, _ := PrependStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			// Add newValue
			newValue := "NewValue"
			newMemberLiteral, err := PrependStringListMember(newStringList, newValue, hl)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(hl)).To(Equal(newValue))
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			// Check priorMemberLiteral of newMemberLiteral
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberLiteralA, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(hl)).To(Equal(newMemberLiteral.GetConceptID(hl)))
			// Check nextMemberLiteral of memberLiteralA
			nextMemberLiteral, err4 := GetNextMemberLiteral(newMemberLiteral, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
		})
	})
	Describe("RemoveStringListMember should work correctly", func() {
		Specify("RemoveStringListMember on empty list should return an error", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			Expect(RemoveStringListMember(newStringList, valueA, hl)).ToNot(Succeed())
		})
		Specify("RemoveStringListMember on singleton list should result in the empty set", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			AppendStringListMember(newStringList, valueA, hl)
			Expect(RemoveStringListMember(newStringList, valueA, hl)).To(Succeed())
			Expect(GetFirstMemberLiteral(newStringList, hl)).To(BeNil())
			Expect(GetLastMemberLiteral(newStringList, hl)).To(BeNil())
			Expect(IsStringListMember(newStringList, valueA, hl)).To(BeFalse())
		})
		Specify("RemoveStringListMember on first element of list should work", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberAfter(newStringList, memberLiteralA, valueB, hl)
			valueC := "C"
			memberReferenceC, _ := AddStringListMemberAfter(newStringList, memberReferenceB, valueC, hl)
			Expect(RemoveStringListMember(newStringList, valueA, hl)).To(Succeed())
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(memberReferenceC.GetConceptID(hl)))
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberReferenceB, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).To(BeNil())
		})
		Specify("RemoveStringListMember on middle element of list should work", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberAfter(newStringList, memberLiteralA, valueB, hl)
			valueC := "C"
			memberReferenceC, _ := AddStringListMemberAfter(newStringList, memberReferenceB, valueC, hl)
			Expect(RemoveStringListMember(newStringList, valueB, hl)).To(Succeed())
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(memberReferenceC.GetConceptID(hl)))
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberReferenceC, hl)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			nextMemberLiteral, err4 := GetNextMemberLiteral(memberLiteralA, hl)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral.GetConceptID(hl)).To(Equal(memberReferenceC.GetConceptID(hl)))
		})
		Specify("RemoveStringListMember on last element of list should work", func() {
			newStringList, _ := NewStringList(uOfD, hl)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, hl)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberAfter(newStringList, memberLiteralA, valueB, hl)
			valueC := "C"
			AddStringListMemberAfter(newStringList, memberReferenceB, valueC, hl)
			Expect(RemoveStringListMember(newStringList, valueC, hl)).To(Succeed())
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, hl)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(hl)).To(Equal(memberLiteralA.GetConceptID(hl)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, hl)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(hl)).To(Equal(memberReferenceB.GetConceptID(hl)))
			nextMemberLiteral, err3 := GetNextMemberLiteral(memberReferenceB, hl)
			Expect(err3).To(BeNil())
			Expect(nextMemberLiteral).To(BeNil())
		})
	})
	Describe("Serialization tests", func() {
		Specify("Instantiated lists should serialize and de-serialze properly", func() {
			uOfD2 := core.NewUniverseOfDiscourse()
			hl2 := uOfD.NewHeldLocks()
			BuildCrlDataStructuresDomain(uOfD2, hl)
			hl2.ReleaseLocksAndWait()
			domain1, _ := uOfD.NewElement(hl)
			list1, err0 := NewStringList(uOfD, hl)
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
			list1FirstElementRefRef, err3 := getStringListReferenceToFirstMemberLiteral(list1, hl)
			Expect(err3).To(BeNil())
			Expect(list1FirstElementRefRef).ToNot(BeNil())
			list1FirstElementRef, err5 := GetFirstMemberLiteral(list1, hl)
			Expect(err5).To(BeNil())
			Expect(list1FirstElementRef).To(BeNil())
			list2FirstElementRefRef, err4 := getStringListReferenceToFirstMemberLiteral(list2, hl2)
			Expect(err4).To(BeNil())
			Expect(list2FirstElementRefRef).ToNot(BeNil())
			list2FirstElementRef, err6 := GetFirstMemberLiteral(list2, hl2)
			Expect(err6).To(BeNil())
			Expect(list2FirstElementRef).To(BeNil())
		})
	})
})
