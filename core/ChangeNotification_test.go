package core

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Element internals test", func() {
	var uOfD *UniverseOfDiscourse
	var trans *Transaction

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
	})

	AfterEach(func() {
		trans.ReleaseLocks()
	})

	Describe("Creating a ConceptState", func() {
		Specify("Creation from Element should work without error", func() {
			el, _ := uOfD.NewElement(trans)
			cs, err := NewConceptState(el)
			Expect(err).To(BeNil())
			Expect(el.GetConceptID(trans)).To(Equal(cs.ConceptID))
		})
		Specify("Element fields should be copied correctly", func() {
			var dummyURI = "http://DummyURI"
			el, err := uOfD.NewElement(trans, dummyURI)
			Expect(err).To(BeNil())
			parent, _ := uOfD.NewElement(trans)
			el.SetOwningConcept(parent, trans)
			el.SetLabel("ElementLabel", trans)
			el.SetDefinition("The definition", trans)
			el.SetIsCore(trans)
			el.SetReadOnly(true, trans)
			cs, err2 := NewConceptState(el)
			Expect(err2).To(BeNil())
			Expect(el.GetConceptID(trans)).To(Equal(cs.ConceptID))
			Expect(ConceptTypeToString(el.GetConceptType())).To(Equal(cs.ConceptType))
			Expect(el.GetOwningConceptID(trans)).To(Equal(cs.OwningConceptID))
			Expect(el.GetLabel(trans)).To(Equal(cs.Label))
			Expect(el.GetDefinition(trans)).To(Equal(cs.Definition))
			Expect(strconv.FormatBool(el.GetIsCore(trans))).To(Equal(cs.IsCore))
			Expect(strconv.FormatBool(el.IsReadOnly(trans))).To(Equal(cs.ReadOnly))
			Expect(strconv.Itoa(el.GetVersion(trans))).To(Equal(cs.Version))
			Expect(el.GetURI(trans)).To(Equal(dummyURI))
		})
		Specify("Creation from Literal should work without error", func() {
			lit, _ := uOfD.NewLiteral(trans)
			cs, err := NewConceptState(lit)
			Expect(err).To(BeNil())
			Expect(lit.GetConceptID(trans)).To(Equal(cs.ConceptID))
		})
		Specify("Literal fields should be copied correctly", func() {
			lit, _ := uOfD.NewLiteral(trans)
			lit.SetLiteralValue("The literal value", trans)
			cs, err := NewConceptState(lit)
			Expect(err).To(BeNil())
			Expect(lit.GetConceptID(trans)).To(Equal(cs.ConceptID))
			Expect(lit.GetLiteralValue(trans)).To(Equal(cs.LiteralValue))
		})
		Specify("Creation from Reference should work without error and fields copied correctly", func() {
			ref, _ := uOfD.NewReference(trans)
			target, _ := uOfD.NewElement(trans)
			target.SetLabel("TargetLabel", trans) // force the version to increment
			ref.SetReferencedConcept(target, OwningConceptID, trans)
			cs, err := NewConceptState(ref)
			Expect(err).To(BeNil())
			Expect(ref.GetConceptID(trans)).To(Equal(cs.ConceptID))
			Expect(target.GetConceptID(trans)).To(Equal(cs.ReferencedConceptID))
			Expect(ref.GetReferencedAttributeName(trans).String()).To(Equal(cs.ReferencedAttributeName))
		})
		Specify("Creation from Refinement should work without error and fields copied correctly", func() {
			ref, _ := uOfD.NewRefinement(trans)
			abs, _ := uOfD.NewElement(trans)
			abs.SetLabel("Abstract Concept", trans)
			rfne, _ := uOfD.NewElement(trans)
			rfne.SetLabel("Refined Concept", trans)
			ref.SetAbstractConcept(abs, trans)
			ref.SetRefinedConcept(rfne, trans)
			cs, err := NewConceptState(ref)
			Expect(err).To(BeNil())
			Expect(ref.GetConceptID(trans)).To(Equal(cs.ConceptID))
			Expect(ref.GetAbstractConceptID(trans)).To(Equal(cs.AbstractConceptID))
			Expect(ref.GetRefinedConceptID(trans)).To(Equal(cs.RefinedConceptID))
		})
	})
})
