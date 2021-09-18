package crldatastructuresdomain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("Set test", func() {
	var uOfD *core.UniverseOfDiscourse
	var hl *core.Transaction

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		hl = uOfD.NewTransaction()
		BuildCrlDataStructuresDomain(uOfD, hl)
		hl.ReleaseLocksAndWait()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("Set should be created correctly", func() {
		Specify("Creation should fail with no specified type", func() {
			_, err := NewSet(uOfD, nil, hl)
			Expect(err).Should(HaveOccurred())
		})
		Specify("Normal creation with Reference type", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			typeReference := newSet.GetFirstOwnedReferenceRefinedFromURI(CrlSetTypeReferenceURI, hl)
			Expect(typeReference).ToNot(BeNil())
			Expect(typeReference.GetReferencedConceptID(hl)).To(Equal(coreReference.GetConceptID(hl)))
		})
	})

	Describe("Adding and deleting members should work correctly", func() {
		Specify("Adding a member of the wrong type should fail", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewElement(hl)
			Expect(AddSetMember(newSet, el, hl)).ToNot(Succeed())
		})
		Specify("Adding a member of the correct type should succeed", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(hl)
			Expect(IsSetMember(newSet, el, hl)).ToNot(BeTrue())
			Expect(AddSetMember(newSet, el, hl)).To(Succeed())
			Expect(IsSetMember(newSet, el, hl)).To(BeTrue())
		})
		Specify("Adding a member for the second time should fail", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(hl)
			Expect(AddSetMember(newSet, el, hl)).To(Succeed())
			Expect(AddSetMember(newSet, el, hl)).ToNot(Succeed())
		})
		Specify("Removing a member that is not in the set should fail", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(hl)
			Expect(RemoveSetMember(newSet, el, hl)).ToNot(Succeed())
		})
		Specify("Removing a member that is in the set should succeed", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(hl)
			Expect(AddSetMember(newSet, el, hl)).To(Succeed())
			Expect(IsSetMember(newSet, el, hl)).To(BeTrue())
			Expect(RemoveSetMember(newSet, el, hl)).To(Succeed())
			Expect(IsSetMember(newSet, el, hl)).ToNot(BeTrue())
		})
		Specify("Removing a member for the second time should fail", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(hl)
			Expect(AddSetMember(newSet, el, hl)).To(Succeed())
			Expect(RemoveSetMember(newSet, el, hl)).To(Succeed())
			Expect(RemoveSetMember(newSet, el, hl)).ToNot(Succeed())
		})
		Specify("Clearing the set should remove members", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, hl)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(hl)
			Expect(AddSetMember(newSet, el, hl)).To(Succeed())
			Expect(IsSetMember(newSet, el, hl)).To(BeTrue())
			ClearSet(newSet, hl)
			Expect(IsSetMember(newSet, el, hl)).ToNot(BeTrue())
		})
	})
})
