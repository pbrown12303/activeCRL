package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("UniverseOfDiscourse", func() {

	var uOfD *UniverseOfDiscourse
	var hl *HeldLocks

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("Creating Initialized UniverseOfDiscourse", func() {
		It("should not be nil", func() {
			Expect(uOfD).ShouldNot(BeNil())
		})
	})

	Describe("Getting the maps", func() {
		It("shoulld return the correct and initialized maps", func() {
			Expect(uOfD.getURIUUIDMap()).ToNot(BeNil())
			Expect(uOfD.getURIUUIDMap() == uOfD.uriUUIDMap)
		})
	})

	Describe("Deriving UUIDs from URIs", func() {
		Specify("The same UUID should be generated each time", func() {
			testURI := "http://activeCrl.com/foo"
			uuid1, err1 := uOfD.generateConceptID(testURI)
			uuid2, err2 := uOfD.generateConceptID(testURI)
			Expect(err1).To(BeNil())
			Expect(err2).To(BeNil())
			Expect(uuid1).To(Equal(uuid2))
		})
	})

	Describe("Creating a Literal", func() {
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
				Expect(uOfD.GetLiteralWithURI(LiteralURI)).To(Equal(lit))
			})
		})
		Specify("UofD GetElement should return the correct type", func() {
			lit, err := uOfD.NewLiteral(hl)
			litID := lit.getConceptIDNoLock()
			Expect(lit).ShouldNot(BeNil())
			Expect(err).Should(BeNil())
			recoveredLit := uOfD.GetElement(litID)
			var correctType bool
			switch recoveredLit.(type) {
			case *element:
				correctType = false
			case *literal:
				correctType = true
			}
			Expect(correctType).To(BeTrue())
			Expect(uOfD.GetLiteral(litID)).To(Equal(lit))
		})
	})

	Describe("Creating an Element", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				el, err := uOfD.NewElement(hl)
				elID := el.getConceptIDNoLock()
				Expect(el).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				Expect(uOfD.GetElement(elID)).To(Equal(el))
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				el, err := uOfD.NewElement(hl, ElementURI)
				Expect(el).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, ElementURI).String()
				Expect(el.GetConceptID(hl)).To(Equal(expectedID))
				Expect(uOfD.GetElementWithURI(ElementURI)).To(Equal(el))
			})
		})
	})

	Describe("Creating a Reference", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				ref, err := uOfD.NewReference(hl)
				refID := ref.getConceptIDNoLock()
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				Expect(uOfD.GetReference(refID)).To(Equal(ref))
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				ref, err := uOfD.NewReference(hl, ReferenceURI)
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, ReferenceURI).String()
				Expect(ref.GetConceptID(hl)).To(Equal(expectedID))
				Expect(uOfD.GetReferenceWithURI(ReferenceURI)).To(Equal(ref))
			})
		})
	})

	Describe("Creating a Refinement", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				ref, err := uOfD.NewRefinement(hl)
				refID := ref.getConceptIDNoLock()
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				Expect(uOfD.GetRefinement(refID)).To(Equal(ref))
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				ref, err := uOfD.NewRefinement(hl, RefinementURI)
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, RefinementURI).String()
				Expect(ref.GetConceptID(hl)).To(Equal(expectedID))
				Expect(uOfD.GetRefinementWithURI(RefinementURI)).To(Equal(ref))
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

	Describe("Test Replicate as Refinement", func() {
		var original Element
		var oChild1 Element
		var oChild1Label string
		var oChild2 Reference
		var oChild2Label string
		var oChild3 Literal
		var oChild3Label string

		BeforeEach(func() {
			original, _ = uOfD.NewElement(hl)
			original.SetLabel("Root", hl)
			oChild1, _ = uOfD.NewElement(hl)
			oChild1.SetOwningConcept(original, hl)
			oChild1Label = "Element"
			oChild1.SetLabel(oChild1Label, hl)
			oChild2, _ = uOfD.NewReference(hl)
			oChild2.SetOwningConcept(original, hl)
			oChild2Label = "Reference"
			oChild2.SetLabel(oChild2Label, hl)
			oChild3, _ = uOfD.NewLiteral(hl)
			oChild3.SetOwningConcept(original, hl)
			oChild3Label = "Literal"
			oChild3.SetLabel(oChild3Label, hl)
		})
		Specify("Replicate should work properly", func() {
			replicate := uOfD.CreateReplicateAsRefinement(original, hl)
			hl.ReleaseLocksAndWait()
			Expect(replicate.IsRefinementOf(original, hl)).To(BeTrue())
			var foundChild1Replicate = false
			var foundChild2Replicate = false
			var foundChild3Replicate = false
			for id := range uOfD.GetOwnedConceptIDs(replicate.GetConceptID(hl)).Iterator().C {
				replicateChild := uOfD.GetElement(id.(string))
				if replicateChild.IsRefinementOf(oChild1, hl) {
					foundChild1Replicate = true
				}
				if replicateChild.IsRefinementOf(oChild2, hl) {
					foundChild2Replicate = true
				}
				if replicateChild.IsRefinementOf(oChild3, hl) {
					foundChild3Replicate = true
				}
			}
			Expect(foundChild1Replicate).To(BeTrue())
			Expect(foundChild2Replicate).To(BeTrue())
			Expect(foundChild3Replicate).To(BeTrue())
		})

		Specify("replicateAsRefinement should be idempotent", func() {
			replicate := uOfD.CreateReplicateAsRefinement(original, hl)
			childCount := uOfD.GetOwnedConceptIDs(replicate.GetConceptID(hl)).Cardinality()
			uOfD.replicateAsRefinement(original, replicate, hl)
			Expect(uOfD.GetOwnedConceptIDs(replicate.GetConceptID(hl)).Cardinality()).To(Equal(childCount))
		})
	})
})
