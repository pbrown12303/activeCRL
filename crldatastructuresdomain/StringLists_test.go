package crldatastructuresdomain

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("StringList test", func() {
	var uOfD *core.UniverseOfDiscourse
	var trans *core.Transaction

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
		BuildCrlDataStructuresDomain(uOfD, trans)
	})

	AfterEach(func() {
		trans.ReleaseLocks()
	})

	Describe("StringList should be created correctly", func() {
		Specify("Normal creation", func() {
			newStringList, err := NewStringList(uOfD, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(GetFirstMemberLiteral(newStringList, trans)).To(BeNil())
			Expect(GetLastMemberLiteral(newStringList, trans)).To(BeNil())
		})
	})

	Describe("NewStringListMemberLiteral should produce the member with prior and next references", func() {
		Specify("Replicate from URI without new URI", func() {
			newLiteralMember, err := NewStringListMemberLiteral(uOfD, trans)
			Expect(err).To(BeNil())
			priorLiteralReference, err2 := getReferenceToPriorMemberLiteral(newLiteralMember, trans)
			Expect(err2).To(BeNil())
			Expect(priorLiteralReference).ToNot(BeNil())
			nextLiteralReference, err3 := getReferenceToNextMemberLiteral(newLiteralMember, trans)
			Expect(err3).To(BeNil())
			Expect(nextLiteralReference).ToNot(BeNil())
		})
		Specify("Replicate from URI with new URI", func() {
			newLiteralMember, err := NewStringListMemberLiteral(uOfD, trans, CrlDataStructuresDomainURI+"/test")
			Expect(err).To(BeNil())
			priorLiteralReference, err2 := getReferenceToPriorMemberLiteral(newLiteralMember, trans)
			Expect(err2).To(BeNil())
			Expect(priorLiteralReference).ToNot(BeNil())
			nextLiteralReference, err3 := getReferenceToNextMemberLiteral(newLiteralMember, trans)
			Expect(err3).To(BeNil())
			Expect(nextLiteralReference).ToNot(BeNil())
		})
	})

	Describe("AddStringListMemberAfter should work correctly", func() {
		Specify("Add after existing solo member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			// Add valueA
			valueA := "A"
			memberLiteralA, err0 := AppendStringListMember(newStringList, valueA, trans)
			Expect(err0).ShouldNot(HaveOccurred())
			// Make sure this member literal has the requisite references as children
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			priorLiteralReference, err2 := getReferenceToPriorMemberLiteral(memberLiteralA, trans)
			Expect(err2).To(BeNil())
			Expect(priorLiteralReference).ToNot(BeNil())
			nextLiteralReference, err3 := getReferenceToNextMemberLiteral(memberLiteralA, trans)
			Expect(err3).To(BeNil())
			Expect(nextLiteralReference).ToNot(BeNil())
			// Add newValue
			newValue := "NewValue"
			newMemberLiteral, err := AddStringListMemberAfter(newStringList, memberLiteralA, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(trans)).To(Equal(newValue))
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			// Check priorMemberLiteral of newMemberLiteral
			priorMemberLiteral, err3 := GetPriorMemberLiteral(newMemberLiteral, trans)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			// Check nextMemberLiteral of memberLiteralA
			nextMemberLiteral, err4 := GetNextMemberLiteral(memberLiteralA, trans)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
		})
		Specify("Add after existing first member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AppendStringListMember(newStringList, valueB, trans)
			newValue := "NewValue"
			// Validate references between A and B
			nmr, _ := GetNextMemberLiteral(memberLiteralA, trans)
			pmr, _ := GetPriorMemberLiteral(memberReferenceB, trans)
			Expect(nmr.GetConceptID(trans)).To(Equal(memberReferenceB.GetConceptID(trans)))
			Expect(pmr.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			// Add new member
			newMemberLiteral, err := AddStringListMemberAfter(newStringList, memberLiteralA, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(trans)).To(Equal(newValue))
			// Check first member next reference
			nextMemberLiteralA, err6 := GetNextMemberLiteral(memberLiteralA, trans)
			Expect(err6).To(BeNil())
			Expect(nextMemberLiteralA.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			// Check new member prior reference
			priorMemberLiteral, err7 := GetPriorMemberLiteral(newMemberLiteral, trans)
			Expect(err7).To(BeNil())
			Expect(priorMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			// Check last member's prior refererence
			priorMemberLiteralB, err3 := GetPriorMemberLiteral(memberReferenceB, trans)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteralB).ToNot(BeNil())
			Expect(priorMemberLiteralB.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			// Check new member next reference
			nextMemberLiteral, err4 := GetNextMemberLiteral(newMemberLiteral, trans)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(trans)).To(Equal(memberReferenceB.GetConceptID(trans)))
			// Check first list member
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			// Check last list member
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(memberReferenceB.GetConceptID(trans)))
		})
	})
	Describe("AddStringListMemberBefore should work correctly", func() {
		Specify("Add before existing solo member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			newValue := "NewValue"
			newMemberLiteral, err := AddStringListMemberBefore(newStringList, memberLiteralA, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(trans)).To(Equal(newValue))
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberLiteralA, trans)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			nextMemberLiteral, err4 := GetNextMemberLiteral(newMemberLiteral, trans)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
		})
		Specify("Add before existing second member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberBefore(newStringList, memberLiteralA, valueB, trans)
			newValue := "NewValue"
			newMemberLiteral, err := AddStringListMemberBefore(newStringList, memberLiteralA, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(trans)).To(Equal(newValue))
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(memberReferenceB.GetConceptID(trans)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberLiteralA, trans)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			nextMemberLiteral, err4 := GetNextMemberLiteral(newMemberLiteral, trans)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			nextMemberLiteralB, err6 := GetNextMemberLiteral(memberReferenceB, trans)
			Expect(err6).To(BeNil())
			Expect(nextMemberLiteralB.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			priorMemberLiteralA, err7 := GetPriorMemberLiteral(memberLiteralA, trans)
			Expect(err7).To(BeNil())
			Expect(priorMemberLiteralA.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
		})
	})
	Describe("AppendStringListMember should work correctly", func() {
		Specify("Append with empty set should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			newValue := "NewValue"
			newMemberLiteral, err := AppendStringListMember(newStringList, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(trans)).To(Equal(newValue))
			Expect(newMemberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
		})
		Specify("Append with existing solo member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			// Add valueA
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			// Add newValue
			newValue := "NewValue"
			newMemberLiteral, err := AppendStringListMember(newStringList, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(trans)).To(Equal(newValue))
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			// Check priorMemberLiteral of newMemberLiteral
			priorMemberLiteral, err3 := GetPriorMemberLiteral(newMemberLiteral, trans)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			// Check nextMemberLiteral of memberLiteralA
			nextMemberLiteral, err4 := GetNextMemberLiteral(memberLiteralA, trans)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
		})
	})
	Describe("ClearStringList should work correctrly", func() {
		Specify("Clear with empty set should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			ClearStringList(newStringList, trans)
			Expect(GetFirstMemberLiteral(newStringList, trans)).To(BeNil())
			Expect(GetLastMemberLiteral(newStringList, trans)).To(BeNil())
		})
		Specify("Clear with single set member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			newValue := "NewValue"
			newMemberLiteral, _ := AppendStringListMember(newStringList, newValue, trans)
			ClearStringList(newStringList, trans)
			Expect(GetFirstMemberLiteral(newStringList, trans)).To(BeNil())
			Expect(GetLastMemberLiteral(newStringList, trans)).To(BeNil())
			Expect(newMemberLiteral.GetOwningConcept(trans)).To(BeNil())
		})
	})
	Describe("GetFirstLiteralForString should work correctly", func() {
		Specify("GetFirstLiteralForString should find each set member", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberBefore(newStringList, memberLiteralA, valueB, trans)
			newValue := "NewValue"
			newMemberLiteral, err := AddStringListMemberBefore(newStringList, memberLiteralA, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(GetFirstLiteralForString(newStringList, valueA, trans)).Should(Equal(memberLiteralA))
			Expect(GetFirstLiteralForString(newStringList, valueB, trans)).Should(Equal(memberReferenceB))
			Expect(GetFirstLiteralForString(newStringList, newValue, trans)).Should(Equal(newMemberLiteral))
		})
		Specify("GetFirstLiteralForString should return nil if element is not in list", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			Expect(GetFirstLiteralForString(newStringList, valueA, trans)).Should(BeNil())
		})
	})
	Describe("IsStringListMember should work correctly", func() {
		Specify("IsStringListMember on an empty set should return false", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			Expect(IsStringListMember(newStringList, valueA, trans)).To(BeFalse())
		})
		Specify("IsStringListMember should work on every member of the set", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			valueB := "B"
			AddStringListMemberBefore(newStringList, memberLiteralA, valueB, trans)
			newValue := "NewValue"
			AddStringListMemberBefore(newStringList, memberLiteralA, newValue, trans)
			Expect(IsStringListMember(newStringList, valueA, trans)).To(BeTrue())
			Expect(IsStringListMember(newStringList, valueB, trans)).To(BeTrue())
			Expect(IsStringListMember(newStringList, newValue, trans)).To(BeTrue())
		})
	})
	Describe("PrependStringListMember should work correctly", func() {
		Specify("Prepend with empty set should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			newValue := "NewValue"
			newMemberLiteral, err := PrependStringListMember(newStringList, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(trans)).To(Equal(newValue))
			Expect(newMemberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
		})
		Specify("Prepend with existing solo member should work correctly", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			// Add valueA
			valueA := "A"
			memberLiteralA, _ := PrependStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			// Add newValue
			newValue := "NewValue"
			newMemberLiteral, err := PrependStringListMember(newStringList, newValue, trans)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newMemberLiteral).ToNot(BeNil())
			Expect(newMemberLiteral.GetLiteralValue(trans)).To(Equal(newValue))
			// Check first member reference
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			// Check last member reference
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			// Check priorMemberLiteral of newMemberLiteral
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberLiteralA, trans)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).ToNot(BeNil())
			Expect(priorMemberLiteral.GetConceptID(trans)).To(Equal(newMemberLiteral.GetConceptID(trans)))
			// Check nextMemberLiteral of memberLiteralA
			nextMemberLiteral, err4 := GetNextMemberLiteral(newMemberLiteral, trans)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral).ToNot(BeNil())
			Expect(nextMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
		})
	})
	Describe("RemoveStringListMember should work correctly", func() {
		Specify("RemoveStringListMember on empty list should return an error", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			Expect(RemoveStringListMember(newStringList, valueA, trans)).ToNot(Succeed())
		})
		Specify("RemoveStringListMember on singleton list should result in the empty set", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			AppendStringListMember(newStringList, valueA, trans)
			Expect(RemoveStringListMember(newStringList, valueA, trans)).To(Succeed())
			Expect(GetFirstMemberLiteral(newStringList, trans)).To(BeNil())
			Expect(GetLastMemberLiteral(newStringList, trans)).To(BeNil())
			Expect(IsStringListMember(newStringList, valueA, trans)).To(BeFalse())
		})
		Specify("RemoveStringListMember on first element of list should work", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberAfter(newStringList, memberLiteralA, valueB, trans)
			valueC := "C"
			memberReferenceC, _ := AddStringListMemberAfter(newStringList, memberReferenceB, valueC, trans)
			Expect(RemoveStringListMember(newStringList, valueA, trans)).To(Succeed())
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(memberReferenceB.GetConceptID(trans)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(memberReferenceC.GetConceptID(trans)))
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberReferenceB, trans)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral).To(BeNil())
		})
		Specify("RemoveStringListMember on middle element of list should work", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberAfter(newStringList, memberLiteralA, valueB, trans)
			valueC := "C"
			memberReferenceC, _ := AddStringListMemberAfter(newStringList, memberReferenceB, valueC, trans)
			Expect(RemoveStringListMember(newStringList, valueB, trans)).To(Succeed())
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(memberReferenceC.GetConceptID(trans)))
			priorMemberLiteral, err3 := GetPriorMemberLiteral(memberReferenceC, trans)
			Expect(err3).To(BeNil())
			Expect(priorMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			nextMemberLiteral, err4 := GetNextMemberLiteral(memberLiteralA, trans)
			Expect(err4).To(BeNil())
			Expect(nextMemberLiteral.GetConceptID(trans)).To(Equal(memberReferenceC.GetConceptID(trans)))
		})
		Specify("RemoveStringListMember on last element of list should work", func() {
			newStringList, _ := NewStringList(uOfD, trans)
			valueA := "A"
			memberLiteralA, _ := AppendStringListMember(newStringList, valueA, trans)
			Expect(memberLiteralA.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)).To(BeTrue())
			valueB := "B"
			memberReferenceB, _ := AddStringListMemberAfter(newStringList, memberLiteralA, valueB, trans)
			valueC := "C"
			AddStringListMemberAfter(newStringList, memberReferenceB, valueC, trans)
			Expect(RemoveStringListMember(newStringList, valueC, trans)).To(Succeed())
			firstMemberLiteral, err1 := GetFirstMemberLiteral(newStringList, trans)
			Expect(err1).To(BeNil())
			Expect(firstMemberLiteral).ToNot(BeNil())
			Expect(firstMemberLiteral.GetConceptID(trans)).To(Equal(memberLiteralA.GetConceptID(trans)))
			lastMemberLiteral, err2 := GetLastMemberLiteral(newStringList, trans)
			Expect(err2).To(BeNil())
			Expect(lastMemberLiteral).ToNot(BeNil())
			Expect(lastMemberLiteral.GetConceptID(trans)).To(Equal(memberReferenceB.GetConceptID(trans)))
			nextMemberLiteral, err3 := GetNextMemberLiteral(memberReferenceB, trans)
			Expect(err3).To(BeNil())
			Expect(nextMemberLiteral).To(BeNil())
		})
	})
	Describe("Serialization tests", func() {
		Specify("Instantiated lists should serialize and de-serialze properly", func() {
			uOfD2 := core.NewUniverseOfDiscourse()
			hl2 := uOfD.NewTransaction()
			defer hl2.ReleaseLocks()
			BuildCrlDataStructuresDomain(uOfD2, trans)
			domain1, _ := uOfD.NewElement(trans)
			list1, err0 := NewStringList(uOfD, trans)
			list1.SetOwningConcept(domain1, trans)
			Expect(err0).To(BeNil())
			Expect(list1).ToNot(BeNil())
			serialized1, err := uOfD.MarshalDomain(domain1, trans)
			Expect(err).To(BeNil())
			domain2, err2 := uOfD2.RecoverDomain(serialized1, hl2)
			Expect(err2).To(BeNil())
			Expect(domain2).ToNot(BeNil())
			Expect(core.RecursivelyEquivalent(domain1, trans, domain2, hl2)).To(BeTrue())
			list2 := uOfD2.GetElement(list1.GetConceptID(trans))
			Expect(list2).ToNot(BeNil())
			list1FirstElementRefRef, err3 := getStringListReferenceToFirstMemberLiteral(list1, trans)
			Expect(err3).To(BeNil())
			Expect(list1FirstElementRefRef).ToNot(BeNil())
			list1FirstElementRef, err5 := GetFirstMemberLiteral(list1, trans)
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
