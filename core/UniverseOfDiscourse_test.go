package core

import (
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("UniverseOfDiscourse", func() {

	var uOfD *UniverseOfDiscourse
	var trans *Transaction

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
	})

	AfterEach(func() {
		trans.ReleaseLocks()
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
				lit, err := uOfD.NewLiteral(trans)
				Expect(lit).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				lit, err := uOfD.NewLiteral(trans, LiteralURI)
				Expect(lit).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, LiteralURI).String()
				Expect(lit.GetConceptID(trans)).To(Equal(expectedID))
				Expect(uOfD.GetLiteralWithURI(LiteralURI)).To(Equal(lit))
			})
		})
		Specify("UofD GetElement should return the correct type", func() {
			lit, err := uOfD.NewLiteral(trans)
			litID := lit.getConceptIDNoLock()
			Expect(lit).ShouldNot(BeNil())
			Expect(err).Should(BeNil())
			recoveredLit := uOfD.GetElement(litID)
			var correctType bool
			switch recoveredLit.GetConceptType() {
			case Literal:
				correctType = true
			default:
				correctType = false
			}
			Expect(correctType).To(BeTrue())
			Expect(uOfD.GetLiteral(litID)).To(Equal(lit))
		})
	})

	Describe("Creating an Element", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				el, err := uOfD.NewElement(trans)
				elID := el.getConceptIDNoLock()
				Expect(el).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				Expect(uOfD.GetElement(elID)).To(Equal(el))
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				el, err := uOfD.NewElement(trans, ElementURI)
				Expect(el).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, ElementURI).String()
				Expect(el.GetConceptID(trans)).To(Equal(expectedID))
				Expect(uOfD.GetElementWithURI(ElementURI)).To(Equal(el))
			})
		})
	})

	Describe("Creating a Reference", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				ref, err := uOfD.NewReference(trans)
				refID := ref.getConceptIDNoLock()
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				Expect(uOfD.GetReference(refID)).To(Equal(ref))
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				ref, err := uOfD.NewReference(trans, ReferenceURI)
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, ReferenceURI).String()
				Expect(ref.GetConceptID(trans)).To(Equal(expectedID))
				Expect(uOfD.GetReferenceWithURI(ReferenceURI)).To(Equal(ref))
			})
		})
	})

	Describe("Creating a Refinement", func() {
		Context("without URI specified", func() {
			It("should not be nil", func() {
				ref, err := uOfD.NewRefinement(trans)
				refID := ref.getConceptIDNoLock()
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				Expect(uOfD.GetRefinement(refID)).To(Equal(ref))
			})
		})
		Context("with URI specified", func() {
			It("should have the correct URI", func() {
				ref, err := uOfD.NewRefinement(trans, RefinementURI)
				Expect(ref).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				expectedID := uuid.NewV5(uuid.NamespaceURL, RefinementURI).String()
				Expect(ref.GetConceptID(trans)).To(Equal(expectedID))
				Expect(uOfD.GetRefinementWithURI(RefinementURI)).To(Equal(ref))
			})
		})
	})

	Describe("Changing the URI of an Element", func() {
		Specify("Setting the URI of an Element should update the uriElementMap of the uOfD", func() {
			el, _ := uOfD.NewElement(trans)
			uri := CorePrefix + "test"
			Expect(uOfD.GetElementWithURI(uri)).To(BeNil())
			el.SetURI(uri, trans)
			Expect(uOfD.GetElementWithURI(uri)).To(Equal(el))
		})
	})

	Describe("Test Replicate as Refinement for Element", func() {
		var original *Concept
		var oChild1 *Concept
		var oChild1Label string
		var oChild2 *Concept
		var oChild2Label string
		var oChild3 *Concept
		var oChild3Label string
		var replicateURI string

		BeforeEach(func() {
			replicateURI = "https://activeCRL.com/ReplicateURI"
			original, _ = uOfD.NewElement(trans)
			original.SetLabel("Root", trans)
			oChild1, _ = uOfD.NewElement(trans)
			oChild1.SetOwningConcept(original, trans)
			oChild1Label = "Element"
			oChild1.SetLabel(oChild1Label, trans)
			oChild2, _ = uOfD.NewReference(trans)
			oChild2.SetOwningConcept(original, trans)
			oChild2Label = "Reference"
			oChild2.SetLabel(oChild2Label, trans)
			oChild3, _ = uOfD.NewLiteral(trans)
			oChild3.SetOwningConcept(original, trans)
			oChild3Label = "Literal"
			oChild3.SetLabel(oChild3Label, trans)
		})
		Specify("Replicate should work properly", func() {
			replicate, err := uOfD.CreateReplicateAsRefinement(original, trans, replicateURI)
			Expect(err).To(BeNil())
			Expect(replicate.IsRefinementOf(original, trans)).To(BeTrue())
			Expect(replicate.GetURI(trans)).To(Equal(replicateURI))
			Expect(uOfD.GetElementWithURI(replicateURI)).To(Equal(replicate))
			foundChild1Replicate := false
			foundChild2Replicate := false
			foundChild3Replicate := false
			child1Suffix := ""
			child1PrefixFound := false
			child2Suffix := ""
			child2PrefixFound := false
			child3Suffix := ""
			child3PrefixFound := false
			it := uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Iterator()
			for id := range it.C {
				replicateChild := uOfD.GetElement(id.(string))
				if replicateChild.IsRefinementOf(oChild1, trans) {
					foundChild1Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					_, child1Suffix, child1PrefixFound = strings.Cut(replicateChildURI, replicateURI+".child")
					Expect(child1PrefixFound).To(BeTrue())
				}
				if replicateChild.IsRefinementOf(oChild2, trans) {
					foundChild2Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					_, child2Suffix, child2PrefixFound = strings.Cut(replicateChildURI, replicateURI+".child")
					Expect(child2PrefixFound).To(BeTrue())
				}
				if replicateChild.IsRefinementOf(oChild3, trans) {
					foundChild3Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					_, child3Suffix, child3PrefixFound = strings.Cut(replicateChildURI, replicateURI+".child")
					Expect(child3PrefixFound).To(BeTrue())
				}
			}
			Expect(foundChild1Replicate).To(BeTrue())
			Expect(foundChild2Replicate).To(BeTrue())
			Expect(foundChild3Replicate).To(BeTrue())
			Expect(child1Suffix != child2Suffix).To(BeTrue())
			Expect(child1Suffix != child3Suffix).To(BeTrue())
			Expect(child2Suffix != child3Suffix).To(BeTrue())
		})

		Specify("replicateAsRefinement should be idempotent", func() {
			replicate, err := uOfD.CreateReplicateAsRefinement(original, trans, replicateURI)
			Expect(err).To(BeNil())
			childCount := uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Cardinality()
			Expect(uOfD.replicateAsRefinement(original, replicate, trans)).To(Succeed())
			Expect(uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Cardinality()).To(Equal(childCount))
		})
	})

	Describe("Test Replicate as Refinement for Literal", func() {
		var original *Concept
		var oChild1 *Concept
		var oChild1Label string
		var oChild2 *Concept
		var oChild2Label string
		var oChild3 *Concept
		var oChild3Label string
		var replicateURI string

		BeforeEach(func() {
			replicateURI = "https://activeCRL.com/ReplicateURI"
			original, _ = uOfD.NewLiteral(trans)
			original.SetLabel("Root", trans)
			oChild1, _ = uOfD.NewElement(trans)
			oChild1.SetOwningConcept(original, trans)
			oChild1Label = "Element"
			oChild1.SetLabel(oChild1Label, trans)
			oChild2, _ = uOfD.NewReference(trans)
			oChild2.SetOwningConcept(original, trans)
			oChild2Label = "Reference"
			oChild2.SetLabel(oChild2Label, trans)
			oChild3, _ = uOfD.NewLiteral(trans)
			oChild3.SetOwningConcept(original, trans)
			oChild3Label = "Literal"
			oChild3.SetLabel(oChild3Label, trans)
		})
		Specify("Replicate should work properly", func() {
			replicate, err := uOfD.CreateReplicateAsRefinement(original, trans, replicateURI)
			Expect(err).To(BeNil())
			Expect(replicate.IsRefinementOf(original, trans)).To(BeTrue())
			Expect(replicate.GetURI(trans)).To(Equal(replicateURI))
			Expect(uOfD.GetElementWithURI(replicateURI)).To(Equal(replicate))
			var foundChild1Replicate = false
			var foundChild2Replicate = false
			var foundChild3Replicate = false
			it := uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Iterator()
			for id := range it.C {
				replicateChild := uOfD.GetElement(id.(string))
				if replicateChild.IsRefinementOf(oChild1, trans) {
					foundChild1Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					Expect(strings.HasPrefix(replicateChildURI, replicateURI+".child")).To(BeTrue())
				}
				if replicateChild.IsRefinementOf(oChild2, trans) {
					foundChild2Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					Expect(strings.HasPrefix(replicateChildURI, replicateURI+".child")).To(BeTrue())
				}
				if replicateChild.IsRefinementOf(oChild3, trans) {
					foundChild3Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					Expect(strings.HasPrefix(replicateChildURI, replicateURI+".child")).To(BeTrue())
				}
			}
			Expect(foundChild1Replicate).To(BeTrue())
			Expect(foundChild2Replicate).To(BeTrue())
			Expect(foundChild3Replicate).To(BeTrue())
		})
		Specify("replicateAsRefinement should be idempotent", func() {
			replicate, err := uOfD.CreateReplicateAsRefinement(original, trans, replicateURI)
			Expect(err).To(BeNil())
			childCount := uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Cardinality()
			Expect(uOfD.replicateAsRefinement(original, replicate, trans)).To(Succeed())
			Expect(uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Cardinality()).To(Equal(childCount))
		})
	})

	Describe("Test Replicate as Refinement for Reference", func() {
		var original *Concept
		var oChild1 *Concept
		var oChild1Label string
		var oChild2 *Concept
		var oChild2Label string
		var oChild3 *Concept
		var oChild3Label string
		var replicateURI string

		BeforeEach(func() {
			replicateURI = "https://activeCRL.com/ReplicateURI"
			original, _ = uOfD.NewReference(trans)
			original.SetLabel("Root", trans)
			oChild1, _ = uOfD.NewElement(trans)
			oChild1.SetOwningConcept(original, trans)
			oChild1Label = "Element"
			oChild1.SetLabel(oChild1Label, trans)
			oChild2, _ = uOfD.NewReference(trans)
			oChild2.SetOwningConcept(original, trans)
			oChild2Label = "Reference"
			oChild2.SetLabel(oChild2Label, trans)
			oChild3, _ = uOfD.NewLiteral(trans)
			oChild3.SetOwningConcept(original, trans)
			oChild3Label = "Literal"
			oChild3.SetLabel(oChild3Label, trans)
		})
		Specify("Replicate should work properly", func() {
			replicate, err := uOfD.CreateReplicateAsRefinement(original, trans, replicateURI)
			Expect(err).To(BeNil())
			Expect(replicate.IsRefinementOf(original, trans)).To(BeTrue())
			Expect(replicate.GetURI(trans)).To(Equal(replicateURI))
			Expect(uOfD.GetElementWithURI(replicateURI)).To(Equal(replicate))
			var foundChild1Replicate = false
			var foundChild2Replicate = false
			var foundChild3Replicate = false
			it := uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Iterator()
			for id := range it.C {
				replicateChild := uOfD.GetElement(id.(string))
				if replicateChild.IsRefinementOf(oChild1, trans) {
					foundChild1Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					Expect(strings.HasPrefix(replicateChildURI, replicateURI+".child")).To(BeTrue())
				}
				if replicateChild.IsRefinementOf(oChild2, trans) {
					foundChild2Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					Expect(strings.HasPrefix(replicateChildURI, replicateURI+".child")).To(BeTrue())
				}
				if replicateChild.IsRefinementOf(oChild3, trans) {
					foundChild3Replicate = true
					replicateChildURI := replicateChild.GetURI(trans)
					Expect(replicateChildURI).ToNot(BeNil())
					Expect(strings.HasPrefix(replicateChildURI, replicateURI+".child")).To(BeTrue())
				}
			}
			Expect(foundChild1Replicate).To(BeTrue())
			Expect(foundChild2Replicate).To(BeTrue())
			Expect(foundChild3Replicate).To(BeTrue())
		})
		Specify("replicateAsRefinement should be idempotent", func() {
			replicate, err := uOfD.CreateReplicateAsRefinement(original, trans, replicateURI)
			Expect(err).To(BeNil())
			childCount := uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Cardinality()
			Expect(uOfD.replicateAsRefinement(original, replicate, trans)).To(Succeed())
			Expect(uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Cardinality()).To(Equal(childCount))
		})
	})

	Describe("Test Replicate as Refinement for Refinement", func() {
		var original *Concept
		var oChild1 *Concept
		var oChild1Label string
		var oChild2 *Concept
		var oChild2Label string
		var oChild3 *Concept
		var oChild3Label string
		var replicateURI string

		BeforeEach(func() {
			replicateURI = "https://activeCRL.com/ReplicateURI"
			original, _ = uOfD.NewRefinement(trans)
			original.SetLabel("Root", trans)
			oChild1, _ = uOfD.NewElement(trans)
			oChild1.SetOwningConcept(original, trans)
			oChild1Label = "Element"
			oChild1.SetLabel(oChild1Label, trans)
			oChild2, _ = uOfD.NewReference(trans)
			oChild2.SetOwningConcept(original, trans)
			oChild2Label = "Reference"
			oChild2.SetLabel(oChild2Label, trans)
			oChild3, _ = uOfD.NewLiteral(trans)
			oChild3.SetOwningConcept(original, trans)
			oChild3Label = "Literal"
			oChild3.SetLabel(oChild3Label, trans)
		})
		Specify("Replicate should work properly", func() {
			replicate, err := uOfD.CreateReplicateAsRefinement(original, trans, replicateURI)
			Expect(err).To(BeNil())
			Expect(replicate.IsRefinementOf(original, trans)).To(BeTrue())
			Expect(replicate.GetURI(trans)).To(Equal(replicateURI))
			Expect(uOfD.GetElementWithURI(replicateURI)).To(Equal(replicate))
			var foundChild1Replicate = false
			var foundChild2Replicate = false
			var foundChild3Replicate = false
			it := uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Iterator()
			for id := range it.C {
				replicateChild := uOfD.GetElement(id.(string))
				if replicateChild.IsRefinementOf(oChild1, trans) {
					foundChild1Replicate = true
				}
				if replicateChild.IsRefinementOf(oChild2, trans) {
					foundChild2Replicate = true
				}
				if replicateChild.IsRefinementOf(oChild3, trans) {
					foundChild3Replicate = true
				}
			}
			Expect(foundChild1Replicate).To(BeTrue())
			Expect(foundChild2Replicate).To(BeTrue())
			Expect(foundChild3Replicate).To(BeTrue())
		})

		Specify("replicateAsRefinement should be idempotent", func() {
			replicate, err := uOfD.CreateReplicateAsRefinement(original, trans, replicateURI)
			Expect(err).To(BeNil())
			childCount := uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Cardinality()
			Expect(uOfD.replicateAsRefinement(original, replicate, trans)).To(Succeed())
			Expect(uOfD.GetConceptsOwnedConceptIDs(replicate.GetConceptID(trans)).Cardinality()).To(Equal(childCount))
		})
	})

	Describe("Test cloning of a Universe of Discourse", func() {
		Specify("Cloning of an empty uOfD should produce an equivalent uOfD", func() {
			uOfD1 := NewUniverseOfDiscourse()
			hl1 := uOfD1.NewTransaction()
			uOfD2 := uOfD1.Clone(hl1)
			hl2 := uOfD2.NewTransaction()
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2, true)).To(BeTrue())
		})
		Specify("Compute function entries should copy correctly", func() {
			uOfD1 := NewUniverseOfDiscourse()
			hl1 := uOfD1.NewTransaction()

			dummyURI := "dummyURI"
			// dummyChangeFunction declared in Housekeeping_test.go
			uOfD1.AddFunction(dummyURI, dummyChangeFunction)
			uOfD2 := uOfD1.Clone(hl1)
			hl2 := uOfD2.NewTransaction()
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2, true)).To(BeTrue())
			Expect(reflect.ValueOf(uOfD2.computeFunctions[dummyURI][0]).Pointer()).To(Equal(reflect.ValueOf(dummyChangeFunction).Pointer()))
		})
		Specify("uriUUIDs should copy correctly", func() {
			uOfD1 := NewUniverseOfDiscourse()
			hl1 := uOfD1.NewTransaction()

			uOfD1.uriUUIDMap.SetEntry("A", "X")
			uOfD2 := uOfD1.Clone(hl1)
			hl2 := uOfD2.NewTransaction()
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2, true)).To(BeTrue())
			Expect(uOfD2.uriUUIDMap.GetEntry("A")).To(Equal("X"))
		})
		Specify("uuidElementMap should copy correctly", func() {
			uOfD1 := NewUniverseOfDiscourse()
			hl1 := uOfD1.NewTransaction()
			el1, _ := uOfD1.NewElement(hl1)

			uOfD2 := uOfD1.Clone(hl1)
			hl2 := uOfD2.NewTransaction()
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2, true)).To(BeTrue())
			el2 := uOfD2.uuidElementMap.GetEntry(el1.GetConceptID(hl1))
			Expect(el2).ToNot(BeNil())
			Expect(el2.GetConceptID(hl2)).To(Equal(el1.GetConceptID(hl1)))
		})
		Specify("ownedIDs should copy correctly", func() {
			uOfD1 := NewUniverseOfDiscourse()
			hl1 := uOfD1.NewTransaction()

			uOfD1.ownedIDsMap.addMappedValue("A", "X")
			uOfD2 := uOfD1.Clone(hl1)
			hl2 := uOfD2.NewTransaction()
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2, true)).To(BeTrue())
			Expect(uOfD2.ownedIDsMap.GetMappedValues("A").Contains("X")).To(BeTrue())
		})
		Specify("listenersMap should copy correctly", func() {
			uOfD1 := NewUniverseOfDiscourse()
			hl1 := uOfD1.NewTransaction()

			uOfD1.listenersMap.addMappedValue("A", "X")
			uOfD2 := uOfD1.Clone(hl1)
			hl2 := uOfD2.NewTransaction()
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2, true)).To(BeTrue())
			Expect(uOfD2.listenersMap.GetMappedValues("A").Contains("X")).To(BeTrue())
		})
	})

	Describe("Test UofD Equivalence", func() {
		var uOfD1 *UniverseOfDiscourse
		var hl1 *Transaction
		var uOfD2 *UniverseOfDiscourse
		var hl2 *Transaction

		BeforeEach(func() {
			uOfD1 = NewUniverseOfDiscourse()
			hl1 = uOfD1.NewTransaction()
			uOfD2 = NewUniverseOfDiscourse()
			hl2 = uOfD2.NewTransaction()
		})

		Specify("Empty uOfDs should be equivalent", func() {
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2)).To(BeTrue())
		})
		Specify("UofDs with different elements should not be equivalent", func() {
			// Element in uOfD1 but not uOfD2
			el1a, _ := uOfD1.NewElement(hl1)
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2)).To(BeFalse())
			// HACK this is not the intended use of clone!
			el2a := clone(el1a, hl1)
			el2a.setUniverseOfDiscourse(nil, hl1)
			Expect(uOfD2.SetUniverseOfDiscourse(el2a, hl2)).To(Succeed())
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2)).To(BeTrue())
		})
		Specify("UofDs with different uuidElementMap should not be equivalent", func() {
			el1a, _ := uOfD1.NewElement(hl1)
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2)).To(BeFalse())
			// HACK this is not the intended use of clone!
			el2a := clone(el1a, hl1)
			el2a.setUniverseOfDiscourse(nil, hl1)
			Expect(uOfD2.SetUniverseOfDiscourse(el2a, hl2)).To(Succeed())
			// now remove the entry in uOfD2.uuidElementMap
			uOfD2.uuidElementMap.DeleteEntry(el2a.GetConceptID(hl2))
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2)).To(BeFalse())
		})
		Specify("UofDs with different uriUUIDMap should not be equivalent", func() {
			uOfD1.uriUUIDMap.SetEntry("A", "X")
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2)).To(BeFalse())
		})
		Specify("UofDs with different ownedIDMap should not be equivalent", func() {
			uOfD1.ownedIDsMap.addMappedValue("A", "X")
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2)).To(BeFalse())
		})
		Specify("UofDs with different listenersMap should not be equivalent", func() {
			uOfD1.listenersMap.addMappedValue("A", "X")
			Expect(uOfD1.IsEquivalent(hl1, uOfD2, hl2)).To(BeFalse())
		})
	})
})
