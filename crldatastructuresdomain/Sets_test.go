package crldatastructuresdomain

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("Set test", func() {
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

	Describe("Set should be created correctly", func() {
		Specify("Creation should fail with no specified type", func() {
			_, err := NewSet(uOfD, nil, trans)
			Expect(err).Should(HaveOccurred())
		})
		Specify("Normal creation with Reference type", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, trans)
			Expect(err).ShouldNot(HaveOccurred())
			typeReference := newSet.GetFirstOwnedReferenceRefinedFromURI(CrlSetTypeReferenceURI, trans)
			Expect(typeReference).ToNot(BeNil())
			Expect(typeReference.GetReferencedConceptID(trans)).To(Equal(coreReference.GetConceptID(trans)))
		})
	})

	Describe("Adding and deleting members should work correctly", func() {
		Specify("Adding a member of the wrong type should fail", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, trans)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewElement(trans)
			Expect(AddSetMember(newSet, el, trans)).ToNot(Succeed())
		})
		Specify("Adding a member of the correct type should succeed", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, trans)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(trans)
			Expect(IsSetMember(newSet, el, trans)).ToNot(BeTrue())
			Expect(AddSetMember(newSet, el, trans)).To(Succeed())
			Expect(IsSetMember(newSet, el, trans)).To(BeTrue())
		})
		Specify("Adding a member for the second time should fail", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, trans)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(trans)
			Expect(AddSetMember(newSet, el, trans)).To(Succeed())
			Expect(AddSetMember(newSet, el, trans)).ToNot(Succeed())
		})
		Specify("Removing a member that is not in the set should fail", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, trans)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(trans)
			Expect(RemoveSetMember(newSet, el, trans)).ToNot(Succeed())
		})
		Specify("Removing a member that is in the set should succeed", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, trans)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(trans)
			Expect(AddSetMember(newSet, el, trans)).To(Succeed())
			Expect(IsSetMember(newSet, el, trans)).To(BeTrue())
			Expect(RemoveSetMember(newSet, el, trans)).To(Succeed())
			Expect(IsSetMember(newSet, el, trans)).ToNot(BeTrue())
		})
		Specify("Removing a member for the second time should fail", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, trans)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(trans)
			Expect(AddSetMember(newSet, el, trans)).To(Succeed())
			Expect(RemoveSetMember(newSet, el, trans)).To(Succeed())
			Expect(RemoveSetMember(newSet, el, trans)).ToNot(Succeed())
		})
		Specify("Clearing the set should remove members", func() {
			coreReference := uOfD.GetReferenceWithURI(core.ReferenceURI)
			newSet, err := NewSet(uOfD, coreReference, trans)
			Expect(err).ShouldNot(HaveOccurred())
			el, _ := uOfD.NewReference(trans)
			Expect(AddSetMember(newSet, el, trans)).To(Succeed())
			Expect(IsSetMember(newSet, el, trans)).To(BeTrue())
			ClearSet(newSet, trans)
			Expect(IsSetMember(newSet, el, trans)).ToNot(BeTrue())
		})
	})
})
