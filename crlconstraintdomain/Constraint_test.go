package crlconstraintdomain

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("Multiplicity evaluation function testing", func() {
	Specify("IsValidMultiplicity returns correct answers", func() {
		Expect(IsValidMultiplicity("")).To(BeTrue())
		Expect(IsValidMultiplicity("1")).To(BeTrue())
		Expect(IsValidMultiplicity("*")).To(BeTrue())
		Expect(IsValidMultiplicity("x")).To(BeFalse())
		Expect(IsValidMultiplicity("0..1")).To(BeTrue())
		Expect(IsValidMultiplicity("x..1")).To(BeFalse())
		Expect(IsValidMultiplicity("1..x")).To(BeFalse())
		Expect(IsValidMultiplicity("0..*")).To(BeTrue())
	})
	Specify("SatisfiesMultiplicity returns correct answers", func() {
		Expect(SatisfiesMultiplicity("x", 0)).To(BeFalse())
		Expect(SatisfiesMultiplicity("", 0)).To(BeTrue())
		Expect(SatisfiesMultiplicity("1", 0)).To(BeFalse())
		Expect(SatisfiesMultiplicity("0", 0)).To(BeTrue())
		Expect(SatisfiesMultiplicity("0", 1)).To(BeFalse())
		Expect(SatisfiesMultiplicity("1", 0)).To(BeFalse())
		Expect(SatisfiesMultiplicity("1", 1)).To(BeTrue())
		Expect(SatisfiesMultiplicity("1", 2)).To(BeFalse())
		Expect(SatisfiesMultiplicity("*", 0)).To(BeTrue())
		Expect(SatisfiesMultiplicity("0..1", 0)).To(BeTrue())
		Expect(SatisfiesMultiplicity("0..1", 1)).To(BeTrue())
		Expect(SatisfiesMultiplicity("0..1", 2)).To(BeFalse())
		Expect(SatisfiesMultiplicity("1..1", 1)).To(BeTrue())
		Expect(SatisfiesMultiplicity("1..1", 0)).To(BeFalse())
		Expect(SatisfiesMultiplicity("1..*", 0)).To(BeFalse())
		Expect(SatisfiesMultiplicity("1..*", 1)).To(BeTrue())
		Expect(SatisfiesMultiplicity("1..*", 5)).To(BeTrue())
	})
})

var _ = Describe("ConstraintSpecification functionality testing", func() {
	var uOfD *core.UniverseOfDiscourse
	var trans *core.Transaction
	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
		BuildCrlConstraintDomain(uOfD, trans)
	})
	Specify("NewMultiplicityConstraintSpecification is properly executed", func() {
		owner, _ := uOfD.NewElement(trans)
		reference, _ := uOfD.NewOwnedReference(owner, "Reference", trans)
		constraintSpecification, err := NewMultiplicityConstraintSpecification(owner, reference, "Reference Constraint", "*", trans)
		Expect(err).To(BeNil())
		Expect(constraintSpecification).ToNot(BeNil())
		Expect(owner.GetFirstOwnedConceptRefinedFromURI(CrlMultiplicityConstraintSpecificationURI, trans)).ToNot(BeNil())
		Expect(constraintSpecification.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlMultiplicityConstraintMultiplicityURI, trans)).ToNot(BeNil())
		constrainedConceptReference := constraintSpecification.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlMultiplicityConstraintConstrainedConceptURI, trans)
		Expect(constrainedConceptReference).ToNot(BeNil())
		Expect(constrainedConceptReference.GetReferencedConcept(trans)).To(Equal(reference))
		Expect(owner.IsRefinementOfURI(CrlMultiplicityConstrainedURI, trans)).To(BeTrue())
	})
	Specify("NewMultiplicityConstraintSpecification with nil owner should fail", func() {
		owner, _ := uOfD.NewElement(trans)
		reference, _ := uOfD.NewOwnedReference(owner, "Reference", trans)
		constraintSpecification, err := NewMultiplicityConstraintSpecification(nil, reference, "Reference Constraint", "*", trans)
		Expect(constraintSpecification).To(BeNil())
		Expect(err).ToNot(BeNil())
	})
	Specify("NewMultiplicityConstraintSpecification with nil reference should fail", func() {
		owner, _ := uOfD.NewElement(trans)
		constraintSpecification, err := NewMultiplicityConstraintSpecification(owner, nil, "Reference Constraint", "*", trans)
		Expect(constraintSpecification).To(BeNil())
		Expect(err).ToNot(BeNil())
	})
	Specify("NewMultiplicityConstraintSpecification with invalid multiplicity should fail", func() {
		owner, _ := uOfD.NewElement(trans)
		reference, _ := uOfD.NewOwnedReference(owner, "Reference", trans)
		constraintSpecification, err := NewMultiplicityConstraintSpecification(owner, reference, "Reference Constraint", "x", trans)
		Expect(constraintSpecification).To(BeNil())
		Expect(err).ToNot(BeNil())
	})
	Specify("NewMultiplicityConstraintSpecification with no Constraint domain in uOfD should fail", func() {
		uOfD := core.NewUniverseOfDiscourse()
		trans := uOfD.NewTransaction()
		owner, _ := uOfD.NewElement(trans)
		reference, _ := uOfD.NewOwnedReference(owner, "Reference", trans)
		constraintSpecification, err := NewMultiplicityConstraintSpecification(owner, reference, "Reference Constraint", "*", trans)
		Expect(constraintSpecification).To(BeNil())
		Expect(err).ToNot(BeNil())
	})
	Specify("Get and Set Multiplicity should perform properly", func() {
		owner, _ := uOfD.NewElement(trans)
		reference, _ := uOfD.NewOwnedReference(owner, "Reference", trans)
		constraintSpecification, _ := NewMultiplicityConstraintSpecification(owner, reference, "Reference Constraint", "*", trans)
		Expect(constraintSpecification.GetMultiplicity(trans)).To(Equal("*"))
		Expect(constraintSpecification.SetMultiplicity("1..*", trans)).To(Succeed())
	})
	Specify("Get and Set Multiplicity should fail with invalid target", func() {
		constraintSpecification, _ := uOfD.NewElement(trans)
		improperCast := (*CrlMultiplicityConstraintSpecification)(constraintSpecification)
		_, err := improperCast.GetMultiplicity(trans)
		Expect(err).ToNot(BeNil())
		Expect(improperCast.SetMultiplicity("1..*", trans)).ToNot(Succeed())
	})
	Specify("SetMultiplicity should fail with invalid multiplicity", func() {
		owner, _ := uOfD.NewElement(trans)
		reference, _ := uOfD.NewOwnedReference(owner, "Reference", trans)
		constraintSpecification, _ := NewMultiplicityConstraintSpecification(owner, reference, "Reference Constraint", "*", trans)
		Expect(constraintSpecification.SetMultiplicity("x", trans)).ToNot(Succeed())
	})
})

var _ = Describe("ConstraintCompliance functionality testing", func() {
	var uOfD *core.UniverseOfDiscourse
	var trans *core.Transaction
	var abstractConcept *core.Concept
	var reference *core.Concept
	var constraintSpecification *CrlMultiplicityConstraintSpecification
	var constrainedConcept *core.Concept
	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
		BuildCrlConstraintDomain(uOfD, trans)
		abstractConcept, _ = uOfD.NewElement(trans)
		reference, _ = uOfD.NewOwnedReference(abstractConcept, "Reference", trans)
		constraintSpecification, _ = NewMultiplicityConstraintSpecification(abstractConcept, reference, "Reference Constraint", "*", trans)
	})
	Specify("NewConstraintCompliance ", func() {
		constrainedConcept, _ = uOfD.NewElement(trans)
		constrainedConcept.SetLabel("Constrained Concept", trans)
		newCompliance := NewConstraintCompliance(constrainedConcept, constraintSpecification.AsCore(), trans)
		Expect(newCompliance.GetOwningConcept(trans)).To(Equal(constrainedConcept))
		Expect((*CrlMultiplicityConstraintSpecification)(GetConstraintSpecification(newCompliance, trans))).To(Equal(constraintSpecification))
	})
})

var _ = Describe("Multiplicity constrant compliance testing", func() {
	var uOfD *core.UniverseOfDiscourse
	var trans *core.Transaction
	var owner *core.Concept
	var reference *core.Concept
	var constraintSpecification *CrlMultiplicityConstraintSpecification
	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
		BuildCrlConstraintDomain(uOfD, trans)
		owner, _ = uOfD.NewElement(trans)
		owner.SetLabel("Owner", trans)
		reference, _ = uOfD.NewOwnedReference(owner, "Reference", trans)
		reference.SetLabel("Owned Reference", trans)
		constraintSpecification, _ = NewMultiplicityConstraintSpecification(owner, reference, "Reference Constraint", "*", trans)
	})
	Specify("Test evaluation of satisfied constraint", func() {
		child, err := uOfD.CreateRefinementOfConcept(owner, "Child", trans)
		Expect(err).To(BeNil())
		constraintCompliance := child.GetFirstOwnedConceptRefinedFromURI(CrlConstraintComplianceURI, trans)
		Expect(constraintCompliance).ToNot(BeNil())
		Expect((*CrlMultiplicityConstraintSpecification)(GetConstraintSpecification(constraintCompliance, trans))).To(Equal(constraintSpecification))
		Expect(IsSatisfied(constraintCompliance, trans)).To(BeTrue())
	})
	Specify("Test evaluation of unsatisfied constraint", func() {
		constraintSpecification.SetMultiplicity("1", trans)
		child, err := uOfD.CreateRefinementOfConcept(owner, "Child", trans)
		Expect(err).To(BeNil())
		constraintCompliance := child.GetFirstOwnedConceptRefinedFromURI(CrlConstraintComplianceURI, trans)
		Expect(constraintCompliance).ToNot(BeNil())
		Expect((*CrlMultiplicityConstraintSpecification)(GetConstraintSpecification(constraintCompliance, trans))).To(Equal(constraintSpecification))
		Expect(IsSatisfied(constraintCompliance, trans)).ToNot(BeTrue())
	})
})
