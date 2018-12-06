package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("UniverseOfDiscourse", func() {

	var uOfD UniverseOfDiscourse
	var hl *HeldLocks

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocks()
	})

	Describe("Creating Initialized UniverseOfDiscourse", func() {
		It("should not be nil", func() {
			Expect(uOfD).ShouldNot(BeNil())
		})
	})

	Describe("Getting the maps", func() {
		It("shoulld return the correct and initialized maps", func() {
			Expect(uOfD.getURIElementMap()).ToNot(BeNil())
			Expect(uOfD.getURIElementMap() == uOfD.(*universeOfDiscourse).uriElementMap)
		})
	})

	Describe("Creating an Literal", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				lit, err := uOfD.NewLiteral(hl)
				Expect(lit).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				lit, err := uOfD.NewLiteral(hl, LiteralURI)
				Expect(lit).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, LiteralURI).String()
				Expect(lit.GetConceptID(hl)).To(Equal(expectedID))
			})
		})
		Specify("UofD GetElement should return the correct type", func() {
			lit, err := uOfD.NewLiteral(hl)
			Expect(lit).ShouldNot(BeNil())
			Expect(err).Should(BeNil())
			recoveredLit := uOfD.GetElement(lit.getConceptIDNoLock())
			var correctType bool
			switch recoveredLit.(type) {
			case *element:
				correctType = false
			case *literal:
				correctType = true
			}
			Expect(correctType).To(BeTrue())
		})
	})

	Describe("Creating an Element", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				el, err := uOfD.NewElement(hl)
				Expect(el).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				el, err := uOfD.NewElement(hl, ElementURI)
				Expect(el).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, ElementURI).String()
				Expect(el.GetConceptID(hl)).To(Equal(expectedID))
			})
		})
	})

	Describe("Creating a Reference", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				ref, err := uOfD.NewReference(hl)
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				ref, err := uOfD.NewReference(hl, ReferenceURI)
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, ReferenceURI).String()
				Expect(ref.GetConceptID(hl)).To(Equal(expectedID))
			})
		})
	})

	Describe("Creating a Refinement", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				ref, err := uOfD.NewRefinement(hl)
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				ref, err := uOfD.NewRefinement(hl, RefinementURI)
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, RefinementURI).String()
				Expect(ref.GetConceptID(hl)).To(Equal(expectedID))
			})
		})
	})

	Describe("Changing the URI of an Element", func() {
		Specify("Setting the URI of an Element should update the uriElementMap of the uOfD", func() {
			el, _ := uOfD.NewElement(hl)
			uri := CorePrefix + "test"
			Expect(uOfD.GetElementWithURI(uri)).To(BeNil())
			el.SetURI(uri, hl)
			hl.ReleaseLocksAndWait()
			Expect(uOfD.GetElementWithURI(uri)).To(Equal(el))
		})
	})
})
