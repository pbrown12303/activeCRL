package core

import (
	"reflect"
	"strconv"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Element internals test", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewTransaction()
	})

	AfterEach(func() {
		hl.ReleaseLocks()
	})

	Describe("Creating a ConceptState", func() {
		Specify("Creation from Element should work without error", func() {
			el, _ := uOfD.NewElement(hl)
			cs, err := NewConceptState(el)
			Expect(err).To(BeNil())
			Expect(el.GetConceptID(hl)).To(Equal(cs.ConceptID))
		})
		Specify("Element fields should be copied correctly", func() {
			var dummyURI = "http://DummyURI"
			el, err := uOfD.NewElement(hl, dummyURI)
			Expect(err).To(BeNil())
			parent, _ := uOfD.NewElement(hl)
			el.SetOwningConcept(parent, hl)
			el.SetLabel("ElementLabel", hl)
			el.SetDefinition("The definition", hl)
			el.SetIsCore(hl)
			el.SetReadOnly(true, hl)
			cs, err2 := NewConceptState(el)
			Expect(err2).To(BeNil())
			Expect(el.GetConceptID(hl)).To(Equal(cs.ConceptID))
			Expect(reflect.TypeOf(el).String()).To(Equal(cs.ConceptType))
			Expect(el.GetOwningConceptID(hl)).To(Equal(cs.OwningConceptID))
			Expect(el.GetLabel(hl)).To(Equal(cs.Label))
			Expect(el.GetDefinition(hl)).To(Equal(cs.Definition))
			Expect(strconv.FormatBool(el.GetIsCore(hl))).To(Equal(cs.IsCore))
			Expect(strconv.FormatBool(el.IsReadOnly(hl))).To(Equal(cs.ReadOnly))
			Expect(strconv.Itoa(el.GetVersion(hl))).To(Equal(cs.Version))
			Expect(el.GetURI(hl)).To(Equal(dummyURI))
		})
		Specify("Creation from Literal should work without error", func() {
			lit, _ := uOfD.NewLiteral(hl)
			cs, err := NewConceptState(lit)
			Expect(err).To(BeNil())
			Expect(lit.GetConceptID(hl)).To(Equal(cs.ConceptID))
		})
		Specify("Literal fields should be copied correctly", func() {
			lit, _ := uOfD.NewLiteral(hl)
			lit.SetLiteralValue("The literal value", hl)
			cs, err := NewConceptState(lit)
			Expect(err).To(BeNil())
			Expect(lit.GetConceptID(hl)).To(Equal(cs.ConceptID))
			Expect(lit.GetLiteralValue(hl)).To(Equal(cs.LiteralValue))
		})
		Specify("Creation from Reference should work without error and fields copied correctly", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			target.SetLabel("TargetLabel", hl) // force the version to increment
			ref.SetReferencedConcept(target, OwningConceptID, hl)
			cs, err := NewConceptState(ref)
			Expect(err).To(BeNil())
			Expect(ref.GetConceptID(hl)).To(Equal(cs.ConceptID))
			Expect(target.GetConceptID(hl)).To(Equal(cs.ReferencedConceptID))
			Expect(strconv.Itoa(target.GetVersion(hl))).To(Equal(cs.ReferencedConceptVersion))
			Expect(ref.GetReferencedAttributeName(hl).String()).To(Equal(cs.ReferencedAttributeName))
		})
		Specify("Creation from Refinement should work without error and fields copied correctly", func() {
			ref, _ := uOfD.NewRefinement(hl)
			abs, _ := uOfD.NewElement(hl)
			abs.SetLabel("Abstract Concept", hl)
			rfne, _ := uOfD.NewElement(hl)
			rfne.SetLabel("Refined Concept", hl)
			ref.SetAbstractConcept(abs, hl)
			ref.SetRefinedConcept(rfne, hl)
			cs, err := NewConceptState(ref)
			Expect(err).To(BeNil())
			Expect(ref.GetConceptID(hl)).To(Equal(cs.ConceptID))
			Expect(ref.GetAbstractConceptID(hl)).To(Equal(cs.AbstractConceptID))
			Expect(strconv.Itoa(ref.GetAbstractConceptVersion(hl))).To(Equal(cs.AbstractConceptVersion))
			Expect(ref.GetRefinedConceptID(hl)).To(Equal(cs.RefinedConceptID))
			Expect(strconv.Itoa(ref.GetRefinedConceptVersion(hl))).To(Equal(cs.RefinedConceptVersion))
		})
	})
})
