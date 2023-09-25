package core

import (
	mapset "github.com/deckarep/golang-set"
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

	Describe("Creating an Element", func() {
		Context("without a URI specified", func() {
			It("should have a non-nil conceptID", func() {
				el, _ := uOfD.NewElement(trans)
				Expect(el.GetConceptID(trans)).ShouldNot(Equal(""))
				Expect(el.GetConceptID(trans) == el.(*concept).ConceptID).To(BeTrue())
			})
			It("should have the UniverseOfDiscourse assigned", func() {
				el, _ := uOfD.NewElement(trans)
				Expect(el.GetUniverseOfDiscourse(trans) == uOfD).To(BeTrue())
				Expect(el.(*concept).uOfD == uOfD).To(BeTrue())
			})
		})
	})

	Describe("Validating ID methods", func() {
		It("should return the correct ID", func() {
			el, _ := uOfD.NewElement(trans)
			Expect(el.GetConceptID(trans) == el.(*concept).ConceptID).To(BeTrue())
			Expect(el.getConceptIDNoLock() == el.(*concept).ConceptID).To(BeTrue())
		})
	})

	Describe("Cloning an Element", func() {
		It("should be equivalent for a newly creaed Element", func() {
			el, _ := uOfD.NewElement(trans)
			clone := clone(el, trans)
			Expect(Equivalent(el, trans, clone, trans)).To(BeTrue())
		})
		It("should be equivalent with an owning concept ID set", func() {
			el, _ := uOfD.NewElement(trans)
			x, _ := uOfD.NewElement(trans)
			el.(*concept).OwningConceptID = x.GetConceptID(trans)
			clone := clone(el, trans)
			Expect(Equivalent(el, trans, clone, trans)).To(BeTrue())
		})
		It("should be equivalent with readOnly set", func() {
			el, _ := uOfD.NewElement(trans)
			el.(*concept).ReadOnly = true
			clone := clone(el, trans)
			Expect(Equivalent(el, trans, clone, trans)).To(BeTrue())
		})
		It("should be equivalent with version set", func() {
			el, _ := uOfD.NewElement(trans)
			el.(*concept).Version.counter = 3
			clone := clone(el, trans)
			Expect(Equivalent(el, trans, clone, trans)).To(BeTrue())
		})
		It("should be equivalent with a universeOfDiscourse set", func() {
			el, _ := uOfD.NewElement(trans)
			el.(*concept).uOfD = uOfD
			clone := clone(el, trans)
			Expect(Equivalent(el, trans, clone, trans)).To(BeTrue())
		})
	})

	Describe("Test setting ownership", func() {
		Specify("Setting ownership via ID should work", func() {
			el, _ := uOfD.NewElement(trans)
			owner, _ := uOfD.NewElement(trans)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			Expect(el.GetOwningConcept(trans)).To(Equal(owner))
		})
		Specify("Setting ownership directly should work", func() {
			el, _ := uOfD.NewElement(trans)
			owner, _ := uOfD.NewElement(trans)
			el.SetOwningConcept(owner, trans)
			Expect(el.GetOwningConcept(trans)).To(Equal(owner))
		})
	})

	Describe("Managing ownedConcepts infrastructure", func() {
		Context("After creating an element", func() {
			Specify("ownedConcepts should be empty", func() {
				el, _ := uOfD.NewElement(trans)
				Expect(uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(trans)).Cardinality() == 0).To(BeTrue())
			})
			Context("after adding an ownedConcept", func() {
				var el Concept
				var ownedConcept Concept
				BeforeEach(func() {
					el, _ = uOfD.NewElement(trans)
					ownedConcept, _ = uOfD.NewElement(trans)
					el.addOwnedConcept(ownedConcept.getConceptIDNoLock(), trans)
				})
				Specify("IsOwnedConcept should return false", func() {
					Expect(el.IsOwnedConcept(ownedConcept, trans)).To(BeTrue())
				})
				It("should be present in GetOwnedConcepts", func() {
					found := false
					it := uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(trans)).Iterator()
					for id := range it.C {
						oc := uOfD.GetElement(id.(string))
						if oc.GetConceptID(trans) == ownedConcept.GetConceptID(trans) {
							found = true
							it.Stop()
						}
					}
					Expect(found).To(BeTrue())
				})
			})
			Context("after removing an owned concept", func() {
				var el Concept
				var ownedConcept Concept
				BeforeEach(func() {
					el, _ = uOfD.NewElement(trans)
					ownedConcept, _ = uOfD.NewElement(trans)
					el.addOwnedConcept(ownedConcept.getConceptIDNoLock(), trans)
					el.removeOwnedConcept(ownedConcept.getConceptIDNoLock(), trans)
				})
				Specify("IsOwnedConcept should return false", func() {
					Expect(el.IsOwnedConcept(ownedConcept, trans)).To(BeFalse())
				})
				It("should not be present in the OwnedConcepts", func() {
					found := false
					it := uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(trans)).Iterator()
					for id := range it.C {
						oc := uOfD.GetElement(id.(string))
						if oc.GetConceptID(trans) == ownedConcept.GetConceptID(trans) {
							found = true
							it.Stop()
						}
					}
					Expect(found).To(BeFalse())
				})
			})
		})
		Specify("Adding and removing an owned Literal should work", func() {
			el, _ := uOfD.NewElement(trans)
			elID := el.getConceptIDNoLock()
			lit, _ := uOfD.NewLiteral(trans)
			lit.SetOwningConceptID(elID, trans)
			Expect(el.IsOwnedConcept(lit, trans)).To(BeTrue())
			lit.SetOwningConceptID("", trans)
			Expect(el.IsOwnedConcept(lit, trans)).To(BeFalse())
		})
	})

	Describe("Managing listeningConcepts infrastructure", func() {
		Context("After creating an element", func() {
			Specify("listeningConcepts should be empty", func() {
				el, _ := uOfD.NewElement(trans)
				Expect(uOfD.getListenerIDs(el.GetConceptID(trans)).Cardinality()).To(Equal(0))
			})
			Context("after adding an referencingConcept", func() {
				var el Concept
				var referencingConcept Concept
				BeforeEach(func() {
					el, _ = uOfD.NewElement(trans)
					referencingConcept, _ = uOfD.NewElement(trans)
					el.addListener(referencingConcept.getConceptIDNoLock(), trans)
				})
				It("should be present in listeners", func() {
					found := false
					it := uOfD.getListenerIDs(el.GetConceptID(trans)).Iterator()
					for id := range it.C {
						oc := uOfD.GetElement(id.(string))
						if oc.GetConceptID(trans) == referencingConcept.GetConceptID(trans) {
							found = true
							it.Stop()
						}
					}
					Expect(found).To(BeTrue())
				})
			})
			Context("after removing an referencingConcept", func() {
				var el Concept
				var referencingConcept Concept
				BeforeEach(func() {
					el, _ = uOfD.NewElement(trans)
					referencingConcept, _ = uOfD.NewElement(trans)
					el.addListener(referencingConcept.getConceptIDNoLock(), trans)
					el.removeListener(referencingConcept.getConceptIDNoLock(), trans)
				})
				It("should not be present in the listeningConcepts", func() {
					found := false
					it := uOfD.getListenerIDs(el.GetConceptID(trans)).Iterator()
					for id := range it.C {
						oc := uOfD.GetElement(id.(string))
						if oc.GetConceptID(trans) == referencingConcept.GetConceptID(trans) {
							found = true
							it.Stop()
						}
					}
					Expect(found).To(BeFalse())
				})
			})
		})
	})

	Describe("Setting concept owner", func() {
		var el Concept
		var owner Concept
		var ownerID string
		BeforeEach(func() {
			el, _ = uOfD.NewElement(trans)
			owner, _ = uOfD.NewElement(trans)
			ownerID = owner.getConceptIDNoLock()
		})
		Context("After creating an Element", func() {
			Specify("conceptOwner should be nil", func() {
				Expect(el.GetOwningConceptID(trans) == "").To(BeTrue())
				Expect(el.GetOwningConcept(trans) == nil).To(BeTrue())
				Expect(el.(*concept).OwningConceptID == "").To(BeTrue())
			})
		})
		Context("After setting the concept owner", func() {
			Specify("conceptOwner should indicate the owner", func() {
				initialVersion := el.GetVersion(trans)
				el.SetOwningConceptID(ownerID, trans)
				Expect(el.GetOwningConceptID(trans) == owner.GetConceptID(trans)).To(BeTrue())
				Expect(el.GetOwningConcept(trans) == owner).To(BeTrue())
				Expect(el.(*concept).OwningConceptID == owner.GetConceptID(trans)).To(BeTrue())
				Expect(owner.IsOwnedConcept(el, trans)).To(BeTrue())
				Expect(el.GetVersion(trans)).To(Equal(initialVersion + 1))
			})
		})
		Context("After setting the concept owner and then setting it to nil", func() {
			Specify("conceptOwner should indicate nil", func() {
				el.SetOwningConceptID(ownerID, trans)
				initialVersion := el.GetVersion(trans)
				el.SetOwningConceptID("", trans)
				Expect(el.GetOwningConceptID(trans) == "").To(BeTrue())
				Expect(el.GetOwningConcept(trans) == nil).To(BeTrue())
				Expect(el.(*concept).OwningConceptID == "").To(BeTrue())
				Expect(owner.IsOwnedConcept(el, trans)).To(BeFalse())
				Expect(el.GetVersion(trans)).To(Equal(initialVersion + 1))
			})
		})
		Context("if Element is read-only", func() {
			It("should fail", func() {
				el.SetReadOnly(true, trans)
				Expect(el.SetOwningConceptID(ownerID, trans)).ToNot(Succeed())
			})
		})
	})

	Describe("Setting read only", func() {
		var child Concept
		var parent Concept
		BeforeEach(func() {
			uOfD = NewUniverseOfDiscourse()
			trans = uOfD.NewTransaction()
		})
		AfterEach(func() {
			trans.ReleaseLocks()
		})
		Context("Owner is not readOnly", func() {
			It("should succed", func() {
				child, _ = uOfD.NewElement(trans)
				parent, _ = uOfD.NewElement(trans)
				child.SetOwningConceptID(parent.getConceptIDNoLock(), trans)
				Expect(child.IsReadOnly(trans)).ToNot(BeTrue())
				Expect(child.(*concept).ReadOnly == false).To(BeTrue())
				Expect(child.SetReadOnly(true, trans)).To(Succeed())
				Expect(child.IsReadOnly(trans)).To(BeTrue())
				Expect(child.(*concept).ReadOnly == true).To(BeTrue())
			})
		})
		Context("Owner is readOnly", func() {
			It("should fail", func() {
				child, _ = uOfD.NewElement(trans)
				parent, _ = uOfD.NewElement(trans)
				Expect(child.SetOwningConceptID(parent.getConceptIDNoLock(), trans)).To(Succeed())
				parent.SetReadOnly(true, trans)
				Expect(child.SetReadOnly(false, trans)).ToNot(Succeed())
			})
		})
		Context("Element is a core element", func() {
			It("should fail", func() {
				child = uOfD.GetElementWithURI(ElementURI)
				Expect(child.SetReadOnly(true, trans)).ToNot(Succeed())
			})
		})
	})

	Describe("Setting universe of discourse", func() {
		It("should change the uOfD pointer correctly", func() {
			el, _ := uOfD.NewElement(trans)
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewTransaction()
			defer hl2.ReleaseLocks()
			Expect(el.GetUniverseOfDiscourse(trans) == uOfD).To(BeTrue())
			Expect(el.(*concept).uOfD == uOfD).To(BeTrue())
			// Can't set new uOfD without removing it from the old uOfD first
			Expect(uOfD2.SetUniverseOfDiscourse(el, trans)).ToNot(Succeed())
			deleteElements := mapset.NewSet(el.GetConceptID(trans))
			Expect(uOfD.DeleteElements(deleteElements, trans)).To(Succeed())
			trans.ReleaseLocks()
			Expect(uOfD2.SetUniverseOfDiscourse(el, hl2)).To(Succeed())
			Expect(el.GetUniverseOfDiscourse(hl2) == uOfD2).To(BeTrue())
			Expect(el.(*concept).uOfD == uOfD2).To(BeTrue())
		})
	})

	Describe("Setting URI", func() {
		var el Concept
		BeforeEach(func() {
			el, _ = uOfD.NewElement(trans)
		})
		Specify("URI should initially nil", func() {
			Expect(el.GetURI(trans)).To(Equal(""))
		})
		Specify("Setting to a valid URI should succeed", func() {
			uri := CorePrefix + "test"
			initialVersion := el.GetVersion(trans)
			Expect(el.SetURI(uri, trans)).To(Succeed())
			Expect(el.GetURI(trans) == uri).To(BeTrue())
			Expect(uOfD.GetElementWithURI(uri)).To(Equal(el))
			Expect(el.GetVersion(trans)).To(Equal(initialVersion + 1))
			Expect(el.SetURI("", trans)).To(Succeed())
			Expect(uOfD.GetElementWithURI(uri)).To(BeNil())
		})
	})

	Describe("Setting Label", func() {
		var el Concept
		BeforeEach(func() {
			el, _ = uOfD.NewElement(trans)
		})
		Specify("Label should initially nil", func() {
			Expect(el.GetLabel(trans)).To(Equal(""))
		})
		Specify("Setting to a valid Label should succeed", func() {
			label := CorePrefix + "test"
			initialVersion := el.GetVersion(trans)
			Expect(el.SetLabel(label, trans)).To(Succeed())
			Expect(el.GetLabel(trans) == label).To(BeTrue())
			Expect(el.GetVersion(trans)).To(Equal(initialVersion + 1))
		})
	})

	Describe("Setting Definition", func() {
		var el Concept
		BeforeEach(func() {
			el, _ = uOfD.NewElement(trans)
		})
		Specify("Definition should initially nil", func() {
			Expect(el.GetDefinition(trans)).To(Equal(""))
		})
		Specify("Setting to a valid Definition should succeed", func() {
			definition := CorePrefix + "test"
			initialVersion := el.GetVersion(trans)
			Expect(el.SetDefinition(definition, trans)).To(Succeed())
			Expect(el.GetDefinition(trans) == definition).To(BeTrue())
			Expect(el.GetVersion(trans)).To(Equal(initialVersion + 1))
		})
	})

	Describe("Validating abstraction infrastructure", func() {
		var owner Concept
		var child Concept
		var firstAbstraction Concept
		var secondAbstraction Concept
		var thirdAbstraction Concept
		var firstAbstractionURI = "http://firstAbstraction"
		var secondAbstractionURI = "http://secondAbstraction"
		var thirdAbstractionURI = "http://thirdAbstraction"
		BeforeEach(func() {
			owner, _ = uOfD.NewElement(trans)
			child, _ = uOfD.NewElement(trans)
			child.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			firstAbstraction, _ = uOfD.NewElement(trans, firstAbstractionURI)
			secondAbstraction, _ = uOfD.NewElement(trans, secondAbstractionURI)
			thirdAbstraction, _ = uOfD.NewElement(trans, thirdAbstractionURI)
		})
		Specify("Initially HasAbstraction should return false", func() {
			Expect(child.IsRefinementOf(firstAbstraction, trans)).To(BeFalse())
			Expect(child.IsRefinementOfURI(firstAbstractionURI, trans)).To(BeFalse())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(firstAbstraction, trans)).To(BeNil())
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(firstAbstractionURI, trans)).To(BeNil())
			Expect(child.IsRefinementOf(secondAbstraction, trans)).To(BeFalse())
			Expect(child.IsRefinementOfURI(secondAbstractionURI, trans)).To(BeFalse())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(secondAbstraction, trans)).To(BeNil())
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(secondAbstractionURI, trans)).To(BeNil())
			Expect(child.IsRefinementOf(thirdAbstraction, trans)).To(BeFalse())
			Expect(child.IsRefinementOfURI(thirdAbstractionURI, trans)).To(BeFalse())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(thirdAbstraction, trans)).To(BeNil())
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(thirdAbstractionURI, trans)).To(BeNil())
		})
		Specify("After adding abstraction, child and owner abstraction-related methods should work", func() {
			ref, _ := uOfD.NewRefinement(trans)
			ref.SetAbstractConceptID(firstAbstraction.getConceptIDNoLock(), trans)
			ref.SetRefinedConceptID(child.getConceptIDNoLock(), trans)
			Expect(child.IsRefinementOf(firstAbstraction, trans)).To(BeTrue())
			Expect(child.IsRefinementOfURI(firstAbstractionURI, trans)).To(BeTrue())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(firstAbstraction, trans)).To(Equal(child))
			foundAbstractions := make(map[string]Concept)
			child.FindAbstractions(foundAbstractions, trans)
			Expect(foundAbstractions[firstAbstraction.getConceptIDNoLock()]).To(Equal(firstAbstraction))
			Expect(foundAbstractions[secondAbstraction.getConceptIDNoLock()]).To(BeNil())
			Expect(foundAbstractions[thirdAbstraction.getConceptIDNoLock()]).To(BeNil())
			Expect(child.IsRefinementOf(secondAbstraction, trans)).To(BeFalse())
			Expect(child.IsRefinementOfURI(secondAbstractionURI, trans)).To(BeFalse())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(secondAbstraction, trans)).To(BeNil())
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(secondAbstractionURI, trans)).To(BeNil())
			Expect(child.IsRefinementOf(thirdAbstraction, trans)).To(BeFalse())
			Expect(child.IsRefinementOfURI(thirdAbstractionURI, trans)).To(BeFalse())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(thirdAbstraction, trans)).To(BeNil())
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(thirdAbstractionURI, trans)).To(BeNil())
		})
		Specify("After adding second-level abstraction, child and owner abstraction-related methods should work", func() {
			ref, _ := uOfD.NewRefinement(trans)
			ref.SetAbstractConceptID(firstAbstraction.getConceptIDNoLock(), trans)
			ref.SetRefinedConceptID(child.getConceptIDNoLock(), trans)
			ref2, _ := uOfD.NewRefinement(trans)
			ref2.SetAbstractConceptID(secondAbstraction.getConceptIDNoLock(), trans)
			ref2.SetRefinedConceptID(firstAbstraction.getConceptIDNoLock(), trans)
			Expect(child.IsRefinementOf(secondAbstraction, trans)).To(BeTrue())
			Expect(child.IsRefinementOfURI(secondAbstractionURI, trans)).To(BeTrue())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(secondAbstraction, trans)).To(Equal(child))
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(secondAbstractionURI, trans)).To(Equal(child))
			Expect(child.IsRefinementOf(thirdAbstraction, trans)).To(BeFalse())
			Expect(child.IsRefinementOfURI(thirdAbstractionURI, trans)).To(BeFalse())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(thirdAbstraction, trans)).To(BeNil())
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(thirdAbstractionURI, trans)).To(BeNil())
			foundAbstractions := make(map[string]Concept)
			child.FindAbstractions(foundAbstractions, trans)
			Expect(foundAbstractions[firstAbstraction.getConceptIDNoLock()]).To(Equal(firstAbstraction))
			Expect(foundAbstractions[secondAbstraction.getConceptIDNoLock()]).To(Equal(secondAbstraction))
			Expect(foundAbstractions[thirdAbstraction.getConceptIDNoLock()]).To(BeNil())
		})
		Specify("An Element should be a refinement of the core Element", func() {
			el, _ := uOfD.NewElement(trans)
			Expect(el.IsRefinementOfURI(ElementURI, trans)).Should(BeTrue())
		})
		Specify("A Literal should be a refinement of the core Element and core Literal", func() {
			el, _ := uOfD.NewLiteral(trans)
			Expect(el.IsRefinementOfURI(ElementURI, trans)).Should(BeTrue())
			Expect(el.IsRefinementOfURI(LiteralURI, trans)).Should(BeTrue())
		})
		Specify("A Reference should be a refinement of the core Element and core Reference", func() {
			el, _ := uOfD.NewReference(trans)
			Expect(el.IsRefinementOfURI(ElementURI, trans)).Should(BeTrue())
			Expect(el.IsRefinementOfURI(ReferenceURI, trans)).Should(BeTrue())
		})
		Specify("A Refinement should be a refinement of the core Element and core Refinement", func() {
			el, _ := uOfD.NewRefinement(trans)
			Expect(el.IsRefinementOfURI(ElementURI, trans)).Should(BeTrue())
			Expect(el.IsRefinementOfURI(RefinementURI, trans)).Should(BeTrue())
		})
	})

	Describe("Testing Element Equivalence", func() {
		var original Concept
		var copy Concept
		BeforeEach(func() {
			original, _ = uOfD.NewElement(trans)
			copy = clone(original, trans)
		})
		Specify("Differences in ConceptID should be detected", func() {
			copy.(*concept).ConceptID = "123"
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
		Specify("Differences in Definition should be detected", func() {
			original.SetDefinition("Definition", trans)
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
		Specify("Differences in IsCore should be detected", func() {
			original.(*concept).IsCore = true
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
		Specify("Differences in Label should be detected", func() {
			original.SetLabel("Label", trans)
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
		Specify("Differences in owned concepts should be detected", func() {
			child, _ := uOfD.NewElement(trans)
			child.SetOwningConceptID(original.getConceptIDNoLock(), trans)
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
		Specify("Differences in owning concept should be detected", func() {
			owner, _ := uOfD.NewElement(trans)
			original.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
		Specify("Differences in readOnly should be detected", func() {
			original.SetReadOnly(true, trans)
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
		Specify("Differences in version should be detected", func() {
			original.(*concept).Version.incrementVersion()
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
		Specify("Differences in URI should be detected", func() {
			original.SetURI("URI", trans)
			Expect(Equivalent(original, trans, copy, trans)).To(BeFalse())
		})
	})

	Describe("Marshal and Unmarshal Tests", func() {
		Specify("Marshal then unmarshal should produce equivalent Elements", func() {
			el, _ := uOfD.NewElement(trans)
			el.SetLabel("label value", trans)
			el.SetDefinition("definition value", trans)
			el.SetURI("URIValue", trans)
			el.SetReadOnly(true, trans)
			el.SetIsCore(trans)
			marshalledElement, err := el.MarshalJSON()
			Expect(err).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewTransaction()
			recoveredElement, err := uOfD2.RecoverElement(marshalledElement, hl2)
			Expect(err).To(BeNil())
			Expect(Equivalent(el, trans, recoveredElement, hl2))
		})
		Specify("Marshal and unmarshal of element and owner should re-establish owner relation", func() {
			el, _ := uOfD.NewElement(trans)
			owner, _ := uOfD.NewElement(trans)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			mEl, err1 := el.MarshalJSON()
			Expect(err1).To(BeNil())
			mOwner, err2 := owner.MarshalJSON()
			Expect(err2).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewTransaction()
			rEl, err3 := uOfD2.RecoverElement(mEl, hl2)
			Expect(err3).To(BeNil())
			rOwner, err4 := uOfD2.RecoverElement(mOwner, hl2)
			Expect(err4).To(BeNil())
			Expect(Equivalent(el, trans, rEl, hl2)).To(BeTrue())
			Expect(RecursivelyEquivalent(owner, trans, rOwner, hl2)).To(BeTrue())
		})
	})

	Describe("Getting owned concepts with abstractions", func() {
		Specify("Getting any concept with abstraction", func() {
			el, _ := uOfD.NewElement(trans)
			owner, _ := uOfD.NewElement(trans)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			ref, _ := uOfD.NewRefinement(trans)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			ref.SetRefinedConceptID(el.getConceptIDNoLock(), trans)
			Expect(owner.GetFirstOwnedConceptRefinedFrom(abs, trans)).To(Equal(el))
			Expect(len(owner.GetOwnedConceptsRefinedFrom(abs, trans))).To(Equal(1))
		})
		Specify("Getting any child with abstractionURI", func() {
			el, _ := uOfD.NewElement(trans)
			owner, _ := uOfD.NewElement(trans)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, trans)
			ref, _ := uOfD.NewRefinement(trans)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			ref.SetRefinedConceptID(el.getConceptIDNoLock(), trans)
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(abstractionURI, trans)).To(Equal(el))
			Expect(len(owner.GetOwnedConceptsRefinedFromURI(abstractionURI, trans))).To(Equal(1))
		})
		Specify("Getting any descendant with abstractionURI", func() {
			el, _ := uOfD.NewElement(trans)
			el2, _ := uOfD.NewElement(trans)
			owner, _ := uOfD.NewElement(trans)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			el2.SetOwningConceptID(el.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, trans)
			ref, _ := uOfD.NewRefinement(trans)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			ref.SetRefinedConceptID(el.getConceptIDNoLock(), trans)
			ref2, _ := uOfD.NewRefinement(trans)
			ref2.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			ref2.SetRefinedConceptID(el2.getConceptIDNoLock(), trans)
			Expect(len(owner.GetOwnedDescendantsRefinedFromURI(abstractionURI, trans))).To(Equal(2))
		})
		Specify("Getting Literal child with abstraction", func() {
			lit, _ := uOfD.NewLiteral(trans)
			owner, _ := uOfD.NewElement(trans)
			lit.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			ref, _ := uOfD.NewRefinement(trans)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			ref.SetRefinedConceptID(lit.getConceptIDNoLock(), trans)
			Expect(owner.GetFirstOwnedLiteralRefinementOf(abs, trans)).To(Equal(lit))
			Expect(len(owner.GetOwnedLiteralsRefinedFrom(abs, trans))).To(Equal(1))
		})
		Specify("Getting Literal child with abstractionURI", func() {
			lit, _ := uOfD.NewLiteral(trans)
			owner, _ := uOfD.NewElement(trans)
			lit.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, trans)
			ref, _ := uOfD.NewRefinement(trans)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			ref.SetRefinedConceptID(lit.getConceptIDNoLock(), trans)
			Expect(owner.GetFirstOwnedLiteralRefinementOfURI(abstractionURI, trans)).To(Equal(lit))
			Expect(len(owner.GetOwnedLiteralsRefinedFromURI(abstractionURI, trans))).To(Equal(1))
		})
		Specify("Getting Reference child with abstraction", func() {
			ref, _ := uOfD.NewReference(trans)
			owner, _ := uOfD.NewElement(trans)
			ref.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			refinement, _ := uOfD.NewRefinement(trans)
			refinement.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			refinement.SetRefinedConceptID(ref.getConceptIDNoLock(), trans)
			Expect(owner.GetFirstOwnedReferenceRefinedFrom(abs, trans)).To(Equal(ref))
			Expect(len(owner.GetOwnedReferencesRefinedFrom(abs, trans))).To(Equal(1))
		})
		Specify("Getting Reference child with abstractionURI", func() {
			ref, _ := uOfD.NewReference(trans)
			owner, _ := uOfD.NewElement(trans)
			ref.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, trans)
			refinement, _ := uOfD.NewRefinement(trans)
			refinement.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			refinement.SetRefinedConceptID(ref.getConceptIDNoLock(), trans)
			Expect(owner.GetFirstOwnedReferenceRefinedFromURI(abstractionURI, trans)).To(Equal(ref))
			Expect(len(owner.GetOwnedReferencesRefinedFromURI(abstractionURI, trans))).To(Equal(1))
		})
		Specify("Getting Refinement child with abstraction", func() {
			ref, _ := uOfD.NewRefinement(trans)
			owner, _ := uOfD.NewElement(trans)
			ref.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			refinement, _ := uOfD.NewRefinement(trans)
			refinement.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			refinement.SetRefinedConceptID(ref.getConceptIDNoLock(), trans)
			Expect(owner.GetFirstOwnedRefinementRefinedFrom(abs, trans)).To(Equal(ref))
			Expect(len(owner.GetOwnedRefinementsRefinedFrom(abs, trans))).To(Equal(1))
		})
		Specify("Getting Refinement child with abstractionURI", func() {
			ref, _ := uOfD.NewRefinement(trans)
			owner, _ := uOfD.NewElement(trans)
			ref.SetOwningConceptID(owner.getConceptIDNoLock(), trans)
			abs, _ := uOfD.NewElement(trans)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, trans)
			refinement, _ := uOfD.NewRefinement(trans)
			refinement.SetAbstractConceptID(abs.getConceptIDNoLock(), trans)
			refinement.SetRefinedConceptID(ref.getConceptIDNoLock(), trans)
			Expect(owner.GetFirstOwnedRefinementRefinedFromURI(abstractionURI, trans)).To(Equal(ref))
			Expect(len(owner.GetOwnedRefinementsRefinedFromURI(abstractionURI, trans))).To(Equal(1))
		})
	})

	Describe("Getting children with URI", func() {
		Specify("GetFirstChildWithURI should work", func() {
			owner, _ := uOfD.NewElement(trans)
			child, _ := uOfD.NewElement(trans)
			child.SetOwningConcept(owner, trans)
			uri := "http://test.uri"
			child.SetURI(uri, trans)
			Expect(owner.GetFirstOwnedConceptWithURI(uri, trans)).To(Equal(child))
		})
		Specify("GetFirstChildLiteralWithURI should work", func() {
			owner, _ := uOfD.NewElement(trans)
			child, _ := uOfD.NewLiteral(trans)
			child.SetOwningConcept(owner, trans)
			uri := "http://test.uri"
			child.SetURI(uri, trans)
			Expect(owner.GetFirstOwnedLiteralWithURI(uri, trans)).To(Equal(child))
		})
		Specify("GetFirstChildReferenceWithURI should work", func() {
			owner, _ := uOfD.NewElement(trans)
			child, _ := uOfD.NewReference(trans)
			child.SetOwningConcept(owner, trans)
			uri := "http://test.uri"
			child.SetURI(uri, trans)
			Expect(owner.GetFirstOwnedReferenceWithURI(uri, trans)).To(Equal(child))
		})
		Specify("GetFirstChildRefinementWithURI should work", func() {
			owner, _ := uOfD.NewElement(trans)
			child, _ := uOfD.NewRefinement(trans)
			child.SetOwningConcept(owner, trans)
			uri := "http://test.uri"
			child.SetURI(uri, trans)
			Expect(owner.GetFirstOwnedRefinementWithURI(uri, trans)).To(Equal(child))
		})
	})
})
