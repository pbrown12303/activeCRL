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

	Describe("Setting and clearing the uOfD should update relevant cached pointers", func() {
		Specify("The owner relationship should be maintained", func() {
			el, _ := uOfD.NewElement(hl)
			elID := el.getConceptIDNoLock()
			owner, _ := uOfD.NewElement(hl)
			ownerID := owner.getConceptIDNoLock()
			el.SetOwningConceptID(ownerID, hl)
			Expect(el.GetOwningConceptID(hl)).To(Equal(ownerID))
			Expect(el.GetOwningConcept(hl)).To(Equal(owner))
			Expect((*owner.GetOwnedConcepts(hl))[elID]).To(Equal(el))
			uOfD.ClearUniverseOfDiscourse(el, hl)
			Expect((*owner.GetOwnedConcepts(hl))[elID]).To(BeNil())
			uOfD.SetUniverseOfDiscourse(el, hl)
			Expect(el.GetOwningConceptID(hl)).To(Equal(ownerID))
			Expect(el.GetOwningConcept(hl)).To(Equal(owner))
			Expect((*owner.GetOwnedConcepts(hl))[elID]).To(Equal(el))
			uOfD.ClearUniverseOfDiscourse(owner, hl)
			Expect(el.GetOwningConceptID(hl)).To(Equal(ownerID))
			Expect(el.GetOwningConcept(hl)).To(BeNil())
			uOfD.SetUniverseOfDiscourse(owner, hl)
			Expect(el.GetOwningConceptID(hl)).To(Equal(ownerID))
			Expect(el.GetOwningConcept(hl)).To(Equal(owner))
			Expect((*owner.GetOwnedConcepts(hl))[elID]).To(Equal(el))
		})
		Specify("The ReferencedConcept relationship should be maintained", func() {
			ref, _ := uOfD.NewReference(hl)
			refID := ref.getConceptIDNoLock()
			target, _ := uOfD.NewElement(hl)
			targetID := target.getConceptIDNoLock()
			ref.SetReferencedConceptID(targetID, hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
			uOfD.ClearUniverseOfDiscourse(target, hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetReferencedConcept(hl)).To(BeNil())
			uOfD.SetUniverseOfDiscourse(target, hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
			uOfD.ClearUniverseOfDiscourse(ref, hl)
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(BeNil())
			uOfD.SetUniverseOfDiscourse(ref, hl)
			Expect(ref.GetReferencedConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetReferencedConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
		})
		Specify("The AbstractConcept relationship should be maintained", func() {
			ref, _ := uOfD.NewRefinement(hl)
			refID := ref.getConceptIDNoLock()
			target, _ := uOfD.NewElement(hl)
			targetID := target.getConceptIDNoLock()
			ref.SetAbstractConceptID(targetID, hl)
			Expect(ref.GetAbstractConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetAbstractConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
			uOfD.ClearUniverseOfDiscourse(target, hl)
			Expect(ref.GetAbstractConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetAbstractConcept(hl)).To(BeNil())
			uOfD.SetUniverseOfDiscourse(target, hl)
			Expect(ref.GetAbstractConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetAbstractConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
			uOfD.ClearUniverseOfDiscourse(ref, hl)
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(BeNil())
			uOfD.SetUniverseOfDiscourse(ref, hl)
			Expect(ref.GetAbstractConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetAbstractConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
		})
		Specify("The RefinedConcept relationship should be maintained", func() {
			ref, _ := uOfD.NewRefinement(hl)
			refID := ref.getConceptIDNoLock()
			target, _ := uOfD.NewElement(hl)
			targetID := target.getConceptIDNoLock()
			ref.SetRefinedConceptID(targetID, hl)
			Expect(ref.GetRefinedConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetRefinedConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
			uOfD.ClearUniverseOfDiscourse(target, hl)
			Expect(ref.GetRefinedConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetRefinedConcept(hl)).To(BeNil())
			uOfD.SetUniverseOfDiscourse(target, hl)
			Expect(ref.GetRefinedConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetRefinedConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
			uOfD.ClearUniverseOfDiscourse(ref, hl)
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(BeNil())
			uOfD.SetUniverseOfDiscourse(ref, hl)
			Expect(ref.GetRefinedConceptID(hl)).To(Equal(targetID))
			Expect(ref.GetRefinedConcept(hl)).To(Equal(target))
			Expect((*target.(*element).listeners.CopyMap())[refID]).To(Equal(ref))
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
			Expect(replicate.HasAbstraction(original, hl)).To(BeTrue())
			var foundChild1Replicate = false
			var foundChild2Replicate = false
			var foundChild3Replicate = false
			for _, replicateChild := range *replicate.GetOwnedConcepts(hl) {
				if replicateChild.HasAbstraction(oChild1, hl) {
					foundChild1Replicate = true
				}
				if replicateChild.HasAbstraction(oChild2, hl) {
					foundChild2Replicate = true
				}
				if replicateChild.HasAbstraction(oChild3, hl) {
					foundChild3Replicate = true
				}
			}
			Expect(foundChild1Replicate).To(BeTrue())
			Expect(foundChild2Replicate).To(BeTrue())
			Expect(foundChild3Replicate).To(BeTrue())
		})

		Specify("replicateAsRefinement should be idempotent", func() {
			replicate := uOfD.CreateReplicateAsRefinement(original, hl)
			childCount := len(*replicate.GetOwnedConcepts(hl))
			uOfD.(*universeOfDiscourse).replicateAsRefinement(original, replicate, hl)
			Expect(len(*replicate.GetOwnedConcepts(hl))).To(Equal(childCount))
		})
	})
})
