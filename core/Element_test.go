package core

import (
	mapset "github.com/deckarep/golang-set"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Element internals test", func() {
	var uOfD *UniverseOfDiscourse
	var hl *HeldLocks

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("Creating an Element", func() {
		Context("without a URI specified", func() {
			It("should have a non-nil conceptID", func() {
				el, _ := uOfD.NewElement(hl)
				Expect(el.GetConceptID(hl)).ShouldNot(Equal(""))
				Expect(el.GetConceptID(hl) == el.(*element).ConceptID).To(BeTrue())
			})
			It("should have the UniverseOfDiscourse assigned", func() {
				el, _ := uOfD.NewElement(hl)
				Expect(el.GetUniverseOfDiscourse(hl) == uOfD).To(BeTrue())
				Expect(el.(*element).uOfD == uOfD).To(BeTrue())
			})
		})
	})

	Describe("Validating ID methods", func() {
		It("should return the correct ID", func() {
			el, _ := uOfD.NewElement(hl)
			Expect(el.GetConceptID(hl) == el.(*element).ConceptID).To(BeTrue())
			Expect(el.getConceptIDNoLock() == el.(*element).ConceptID).To(BeTrue())
		})
	})

	Describe("Cloning an Element", func() {
		It("should be equivalent for a newly creaed Element", func() {
			el, _ := uOfD.NewElement(hl)
			clone := clone(el, hl)
			Expect(Equivalent(el, hl, clone, hl)).To(BeTrue())
		})
		It("should be equivalent with an owning concept ID set", func() {
			el, _ := uOfD.NewElement(hl)
			x, _ := uOfD.NewElement(hl)
			el.(*element).OwningConceptID = x.GetConceptID(hl)
			clone := clone(el, hl)
			Expect(Equivalent(el, hl, clone, hl)).To(BeTrue())
		})
		It("should be equivalent with readOnly set", func() {
			el, _ := uOfD.NewElement(hl)
			el.(*element).ReadOnly = true
			clone := clone(el, hl)
			Expect(Equivalent(el, hl, clone, hl)).To(BeTrue())
		})
		It("should be equivalent with version set", func() {
			el, _ := uOfD.NewElement(hl)
			el.(*element).Version.counter = 3
			clone := clone(el, hl)
			Expect(Equivalent(el, hl, clone, hl)).To(BeTrue())
		})
		It("should be equivalent with a universeOfDiscourse set", func() {
			el, _ := uOfD.NewElement(hl)
			el.(*element).uOfD = uOfD
			clone := clone(el, hl)
			Expect(Equivalent(el, hl, clone, hl)).To(BeTrue())
		})
	})

	Describe("Test setting ownership", func() {
		Specify("Setting ownership via ID should work", func() {
			el, _ := uOfD.NewElement(hl)
			owner, _ := uOfD.NewElement(hl)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			Expect(el.GetOwningConcept(hl)).To(Equal(owner))
		})
		Specify("Setting ownership directly should work", func() {
			el, _ := uOfD.NewElement(hl)
			owner, _ := uOfD.NewElement(hl)
			el.SetOwningConcept(owner, hl)
			Expect(el.GetOwningConcept(hl)).To(Equal(owner))
		})
	})

	Describe("Managing ownedConcepts infrastructure", func() {
		Context("After creating an element", func() {
			Specify("ownedConcepts should be empty", func() {
				el, _ := uOfD.NewElement(hl)
				Expect(uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Cardinality() == 0).To(BeTrue())
			})
			Context("after adding an ownedConcept", func() {
				var el Element
				var ownedConcept Element
				BeforeEach(func() {
					el, _ = uOfD.NewElement(hl)
					ownedConcept, _ = uOfD.NewElement(hl)
					el.addOwnedConcept(ownedConcept.getConceptIDNoLock(), hl)
				})
				Specify("IsOwnedConcept should return false", func() {
					Expect(el.IsOwnedConcept(ownedConcept, hl)).To(BeTrue())
				})
				It("should be present in GetOwnedConcepts", func() {
					found := false
					it := uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Iterator()
					defer it.Stop()
					for id := range it.C {
						oc := uOfD.GetElement(id.(string))
						if oc.GetConceptID(hl) == ownedConcept.GetConceptID(hl) {
							found = true
						}
					}
					Expect(found).To(BeTrue())
				})
			})
			Context("after removing an owned concept", func() {
				var el Element
				var ownedConcept Element
				BeforeEach(func() {
					el, _ = uOfD.NewElement(hl)
					ownedConcept, _ = uOfD.NewElement(hl)
					el.addOwnedConcept(ownedConcept.getConceptIDNoLock(), hl)
					el.removeOwnedConcept(ownedConcept.getConceptIDNoLock(), hl)
				})
				Specify("IsOwnedConcept should return false", func() {
					Expect(el.IsOwnedConcept(ownedConcept, hl)).To(BeFalse())
				})
				It("should not be present in the OwnedConcepts", func() {
					found := false
					it := uOfD.GetConceptsOwnedConceptIDs(el.GetConceptID(hl)).Iterator()
					defer it.Stop()
					for id := range it.C {
						oc := uOfD.GetElement(id.(string))
						if oc.GetConceptID(hl) == ownedConcept.GetConceptID(hl) {
							found = true
						}
					}
					Expect(found).To(BeFalse())
				})
			})
		})
		Specify("Adding and removing an owned Literal should work", func() {
			el, _ := uOfD.NewElement(hl)
			elID := el.getConceptIDNoLock()
			lit, _ := uOfD.NewLiteral(hl)
			lit.SetOwningConceptID(elID, hl)
			Expect(el.IsOwnedConcept(lit, hl)).To(BeTrue())
			lit.SetOwningConceptID("", hl)
			Expect(el.IsOwnedConcept(lit, hl)).To(BeFalse())
		})
	})

	Describe("Managing listeningConcepts infrastructure", func() {
		Context("After creating an element", func() {
			Specify("listeningConcepts should be empty", func() {
				el, _ := uOfD.NewElement(hl)
				Expect(uOfD.getListenerIDs(el.GetConceptID(hl)).Cardinality()).To(Equal(0))
			})
			Context("after adding an referencingConcept", func() {
				var el Element
				var referencingConcept Element
				BeforeEach(func() {
					el, _ = uOfD.NewElement(hl)
					referencingConcept, _ = uOfD.NewElement(hl)
					el.addListener(referencingConcept.getConceptIDNoLock(), hl)
				})
				It("should be present in listeners", func() {
					found := false
					it := uOfD.getListenerIDs(el.GetConceptID(hl)).Iterator()
					defer it.Stop()
					for id := range it.C {
						oc := uOfD.GetElement(id.(string))
						if oc.GetConceptID(hl) == referencingConcept.GetConceptID(hl) {
							found = true
						}
					}
					Expect(found).To(BeTrue())
				})
			})
			Context("after removing an referencingConcept", func() {
				var el Element
				var referencingConcept Element
				BeforeEach(func() {
					el, _ = uOfD.NewElement(hl)
					referencingConcept, _ = uOfD.NewElement(hl)
					el.addListener(referencingConcept.getConceptIDNoLock(), hl)
					el.removeListener(referencingConcept.getConceptIDNoLock(), hl)
				})
				It("should not be present in the listeningConcepts", func() {
					found := false
					it := uOfD.getListenerIDs(el.GetConceptID(hl)).Iterator()
					defer it.Stop()
					for id := range it.C {
						oc := uOfD.GetElement(id.(string))
						if oc.GetConceptID(hl) == referencingConcept.GetConceptID(hl) {
							found = true
						}
					}
					Expect(found).To(BeFalse())
				})
			})
		})
	})

	Describe("Setting concept owner", func() {
		var el Element
		var owner Element
		var ownerID string
		BeforeEach(func() {
			el, _ = uOfD.NewElement(hl)
			owner, _ = uOfD.NewElement(hl)
			ownerID = owner.getConceptIDNoLock()
		})
		Context("After creating an Element", func() {
			Specify("conceptOwner should be nil", func() {
				Expect(el.GetOwningConceptID(hl) == "").To(BeTrue())
				Expect(el.GetOwningConcept(hl) == nil).To(BeTrue())
				Expect(el.(*element).OwningConceptID == "").To(BeTrue())
			})
		})
		Context("After setting the concept owner", func() {
			Specify("conceptOwner should indicate the owner", func() {
				initialVersion := el.GetVersion(hl)
				el.SetOwningConceptID(ownerID, hl)
				Expect(el.GetOwningConceptID(hl) == owner.GetConceptID(hl)).To(BeTrue())
				Expect(el.GetOwningConcept(hl) == owner).To(BeTrue())
				Expect(el.(*element).OwningConceptID == owner.GetConceptID(hl)).To(BeTrue())
				Expect(owner.IsOwnedConcept(el, hl)).To(BeTrue())
				Expect(el.GetVersion(hl)).To(Equal(initialVersion + 1))
			})
		})
		Context("After setting the concept owner and then setting it to nil", func() {
			Specify("conceptOwner should indicate nil", func() {
				el.SetOwningConceptID(ownerID, hl)
				initialVersion := el.GetVersion(hl)
				el.SetOwningConceptID("", hl)
				Expect(el.GetOwningConceptID(hl) == "").To(BeTrue())
				Expect(el.GetOwningConcept(hl) == nil).To(BeTrue())
				Expect(el.(*element).OwningConceptID == "").To(BeTrue())
				Expect(owner.IsOwnedConcept(el, hl)).To(BeFalse())
				Expect(el.GetVersion(hl)).To(Equal(initialVersion + 1))
			})
		})
		Context("if Element is read-only", func() {
			It("should fail", func() {
				el.SetReadOnly(true, hl)
				Expect(el.SetOwningConceptID(ownerID, hl)).ToNot(Succeed())
			})
		})
	})

	Describe("Setting read only", func() {
		var child Element
		var parent Element
		BeforeEach(func() {
			uOfD = NewUniverseOfDiscourse()
			hl = uOfD.NewHeldLocks()
		})
		AfterEach(func() {
			hl.ReleaseLocksAndWait()
		})
		Context("Owner is not readOnly", func() {
			It("should succed", func() {
				child, _ = uOfD.NewElement(hl)
				parent, _ = uOfD.NewElement(hl)
				child.SetOwningConceptID(parent.getConceptIDNoLock(), hl)
				Expect(child.IsReadOnly(hl)).ToNot(BeTrue())
				Expect(child.(*element).ReadOnly == false).To(BeTrue())
				Expect(child.SetReadOnly(true, hl)).To(Succeed())
				Expect(child.IsReadOnly(hl)).To(BeTrue())
				Expect(child.(*element).ReadOnly == true).To(BeTrue())
			})
		})
		Context("Owner is readOnly", func() {
			It("should fail", func() {
				child, _ = uOfD.NewElement(hl)
				parent, _ = uOfD.NewElement(hl)
				Expect(child.SetOwningConceptID(parent.getConceptIDNoLock(), hl)).To(Succeed())
				parent.SetReadOnly(true, hl)
				Expect(child.SetReadOnly(false, hl)).ToNot(Succeed())
			})
		})
		Context("Element is a core element", func() {
			It("should fail", func() {
				child = uOfD.GetElementWithURI(ElementURI)
				Expect(child.SetReadOnly(true, hl)).ToNot(Succeed())
			})
		})
	})

	Describe("Testing Version infrastructure", func() {
		var el Element
		BeforeEach(func() {
			el, _ = uOfD.NewElement(hl)
		})
		Specify("Version should increment when incrementVersion is called", func() {
			initialVersion := el.GetVersion(hl)
			Expect(initialVersion == el.(*element).Version.getVersion()).To(BeTrue())
			el.incrementVersion(hl)
			nextVersion := el.GetVersion(hl)
			Expect(nextVersion > initialVersion).To(BeTrue())
			Expect(nextVersion == el.(*element).Version.getVersion()).To(BeTrue())
		})
		Specify("Owner's version should increment when child's increment version is called", func() {
			owner, _ := uOfD.NewElement(hl)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			initialVersion := el.GetVersion(hl)
			initialOwnerVersion := owner.GetVersion(hl)
			el.incrementVersion(hl)
			nextVersion := el.GetVersion(hl)
			nextOwnerVersion := owner.GetVersion(hl)
			Expect(nextVersion > initialVersion).To(BeTrue())
			Expect(nextOwnerVersion > initialOwnerVersion).To(BeTrue())
		})
	})

	Describe("Setting universe of discourse", func() {
		It("should change the uOfD pointer correctly", func() {
			el, _ := uOfD.NewElement(hl)
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewHeldLocks()
			Expect(el.GetUniverseOfDiscourse(hl) == uOfD).To(BeTrue())
			Expect(el.(*element).uOfD == uOfD).To(BeTrue())
			// Can't set new uOfD without removing it from the old uOfD first
			Expect(uOfD2.SetUniverseOfDiscourse(el, hl)).ToNot(Succeed())
			deleteElements := mapset.NewSet(el.GetConceptID(hl))
			Expect(uOfD.DeleteElements(deleteElements, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()
			Expect(uOfD2.SetUniverseOfDiscourse(el, hl2)).To(Succeed())
			Expect(el.GetUniverseOfDiscourse(hl2) == uOfD2).To(BeTrue())
			Expect(el.(*element).uOfD == uOfD2).To(BeTrue())
		})
	})

	Describe("Setting URI", func() {
		var el Element
		BeforeEach(func() {
			el, _ = uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
		})
		Specify("URI should initially nil", func() {
			Expect(el.GetURI(hl)).To(Equal(""))
		})
		Specify("Setting to a valid URI should succeed", func() {
			uri := CorePrefix + "test"
			initialVersion := el.GetVersion(hl)
			Expect(el.SetURI(uri, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()
			Expect(el.GetURI(hl) == uri).To(BeTrue())
			Expect(uOfD.GetElementWithURI(uri)).To(Equal(el))
			Expect(el.GetVersion(hl)).To(Equal(initialVersion + 1))
			Expect(el.SetURI("", hl)).To(Succeed())
			Expect(uOfD.GetElementWithURI(uri)).To(BeNil())
		})
	})

	Describe("Setting Label", func() {
		var el Element
		BeforeEach(func() {
			el, _ = uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
		})
		Specify("Label should initially nil", func() {
			Expect(el.GetLabel(hl)).To(Equal(""))
		})
		Specify("Setting to a valid Label should succeed", func() {
			label := CorePrefix + "test"
			initialVersion := el.GetVersion(hl)
			Expect(el.SetLabel(label, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()
			Expect(el.GetLabel(hl) == label).To(BeTrue())
			Expect(el.GetVersion(hl)).To(Equal(initialVersion + 1))
		})
	})

	Describe("Setting Definition", func() {
		var el Element
		BeforeEach(func() {
			el, _ = uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
		})
		Specify("Definition should initially nil", func() {
			Expect(el.GetDefinition(hl)).To(Equal(""))
		})
		Specify("Setting to a valid Definition should succeed", func() {
			definition := CorePrefix + "test"
			initialVersion := el.GetVersion(hl)
			Expect(el.SetDefinition(definition, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()
			Expect(el.GetDefinition(hl) == definition).To(BeTrue())
			Expect(el.GetVersion(hl)).To(Equal(initialVersion + 1))
		})
	})

	Describe("Validating abstraction infrastructure", func() {
		var owner Element
		var child Element
		var firstAbstraction Element
		var secondAbstraction Element
		var firstAbstractionURI = "http://firstAbstraction"
		var secondAbstractionURI = "http://secondAbstraction"
		BeforeEach(func() {
			owner, _ = uOfD.NewElement(hl)
			child, _ = uOfD.NewElement(hl)
			child.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			firstAbstraction, _ = uOfD.NewElement(hl)
			firstAbstraction.SetURI(firstAbstractionURI, hl)
			secondAbstraction, _ = uOfD.NewElement(hl)
			secondAbstraction.SetURI(secondAbstractionURI, hl)
		})
		Specify("Initially HasAbstraction should return false", func() {
			Expect(child.IsRefinementOf(firstAbstraction, hl)).To(BeFalse())
			Expect(child.IsRefinementOfURI(firstAbstractionURI, hl)).To(BeFalse())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(firstAbstraction, hl)).To(BeNil())
		})
		Specify("After adding abstraction, child and owner abstraction-related methods should work", func() {
			ref, _ := uOfD.NewRefinement(hl)
			ref.SetAbstractConceptID(firstAbstraction.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(child.getConceptIDNoLock(), hl)
			Expect(child.IsRefinementOf(firstAbstraction, hl)).To(BeTrue())
			Expect(child.IsRefinementOfURI(firstAbstractionURI, hl)).To(BeTrue())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(firstAbstraction, hl)).To(Equal(child))
			foundAbstractions := make(map[string]Element)
			child.FindAbstractions(foundAbstractions, hl)
			Expect(foundAbstractions[firstAbstraction.getConceptIDNoLock()]).To(Equal(firstAbstraction))
		})
		Specify("After adding second-level abstraction, child and owner abstraction-related methods should work", func() {
			ref, _ := uOfD.NewRefinement(hl)
			ref.SetAbstractConceptID(firstAbstraction.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(child.getConceptIDNoLock(), hl)
			ref2, _ := uOfD.NewRefinement(hl)
			ref2.SetAbstractConceptID(secondAbstraction.getConceptIDNoLock(), hl)
			ref2.SetRefinedConceptID(firstAbstraction.getConceptIDNoLock(), hl)
			Expect(child.IsRefinementOf(secondAbstraction, hl)).To(BeTrue())
			Expect(child.IsRefinementOfURI(secondAbstractionURI, hl)).To(BeTrue())
			Expect(owner.GetFirstOwnedConceptRefinedFrom(secondAbstraction, hl)).To(Equal(child))
			foundAbstractions := make(map[string]Element)
			child.FindAbstractions(foundAbstractions, hl)
			Expect(foundAbstractions[firstAbstraction.getConceptIDNoLock()]).To(Equal(firstAbstraction))
			Expect(foundAbstractions[secondAbstraction.getConceptIDNoLock()]).To(Equal(secondAbstraction))
		})
		Specify("An Element should be a refinement of the core Element", func() {
			el, _ := uOfD.NewElement(hl)
			Expect(el.IsRefinementOfURI(ElementURI, hl)).Should(BeTrue())
		})
		Specify("A Literal should be a refinement of the core Element and core Literal", func() {
			el, _ := uOfD.NewLiteral(hl)
			Expect(el.IsRefinementOfURI(ElementURI, hl)).Should(BeTrue())
			Expect(el.IsRefinementOfURI(LiteralURI, hl)).Should(BeTrue())
		})
		Specify("A Reference should be a refinement of the core Element and core Reference", func() {
			el, _ := uOfD.NewReference(hl)
			Expect(el.IsRefinementOfURI(ElementURI, hl)).Should(BeTrue())
			Expect(el.IsRefinementOfURI(ReferenceURI, hl)).Should(BeTrue())
		})
		Specify("A Refinement should be a refinement of the core Element and core Refinement", func() {
			el, _ := uOfD.NewRefinement(hl)
			Expect(el.IsRefinementOfURI(ElementURI, hl)).Should(BeTrue())
			Expect(el.IsRefinementOfURI(RefinementURI, hl)).Should(BeTrue())
		})
	})

	Describe("Testing Element Equivalence", func() {
		var original Element
		var copy Element
		BeforeEach(func() {
			original, _ = uOfD.NewElement(hl)
			copy = clone(original, hl)
		})
		Specify("Differences in ConceptID should be detected", func() {
			// Have to release locks because HeldLocks keeps track by ConceptID
			hl.ReleaseLocksAndWait()
			original.(*element).ConceptID = "123"
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
		Specify("Differences in Definition should be detected", func() {
			original.SetDefinition("Definition", hl)
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
		Specify("Differences in IsCore should be detected", func() {
			original.(*element).IsCore = true
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
		Specify("Differences in Label should be detected", func() {
			original.SetLabel("Label", hl)
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
		Specify("Differences in owned concepts should be detected", func() {
			child, _ := uOfD.NewElement(hl)
			child.SetOwningConceptID(original.getConceptIDNoLock(), hl)
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
		Specify("Differences in owning concept should be detected", func() {
			owner, _ := uOfD.NewElement(hl)
			original.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
		Specify("Differences in readOnly should be detected", func() {
			original.SetReadOnly(true, hl)
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
		Specify("Differences in version should be detected", func() {
			original.incrementVersion(hl)
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
		Specify("Differences in URI should be detected", func() {
			original.SetURI("URI", hl)
			Expect(Equivalent(original, hl, copy, hl)).To(BeFalse())
		})
	})

	Describe("Marshal and Unmarshal Tests", func() {
		Specify("Marshal then unmarshal should produce equivalent Elements", func() {
			el, _ := uOfD.NewElement(hl)
			el.SetLabel("label value", hl)
			el.SetDefinition("definition value", hl)
			el.SetURI("URIValue", hl)
			el.SetReadOnly(true, hl)
			el.SetIsCore(hl)
			marshalledElement, err := el.MarshalJSON()
			Expect(err).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewHeldLocks()
			recoveredElement, err := uOfD2.RecoverElement(marshalledElement, hl2)
			Expect(err).To(BeNil())
			Expect(Equivalent(el, hl, recoveredElement, hl2))
		})
		Specify("Marshal and unmarshal of element and owner should re-establish owner relation", func() {
			el, _ := uOfD.NewElement(hl)
			owner, _ := uOfD.NewElement(hl)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			mEl, err1 := el.MarshalJSON()
			Expect(err1).To(BeNil())
			mOwner, err2 := owner.MarshalJSON()
			Expect(err2).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewHeldLocks()
			rEl, err3 := uOfD2.RecoverElement(mEl, hl2)
			Expect(err3).To(BeNil())
			rOwner, err4 := uOfD2.RecoverElement(mOwner, hl2)
			Expect(err4).To(BeNil())
			Expect(Equivalent(el, hl, rEl, hl2)).To(BeTrue())
			Expect(RecursivelyEquivalent(owner, hl, rOwner, hl2)).To(BeTrue())
		})
	})

	Describe("Getting owned concepts with abstractions", func() {
		Specify("Getting any concept with abstraction", func() {
			el, _ := uOfD.NewElement(hl)
			owner, _ := uOfD.NewElement(hl)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			ref, _ := uOfD.NewRefinement(hl)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(el.getConceptIDNoLock(), hl)
			Expect(owner.GetFirstOwnedConceptRefinedFrom(abs, hl)).To(Equal(el))
			Expect(len(owner.GetOwnedConceptsRefinedFrom(abs, hl))).To(Equal(1))
		})
		Specify("Getting any child with abstractionURI", func() {
			el, _ := uOfD.NewElement(hl)
			owner, _ := uOfD.NewElement(hl)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, hl)
			ref, _ := uOfD.NewRefinement(hl)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(el.getConceptIDNoLock(), hl)
			Expect(owner.GetFirstOwnedConceptRefinedFromURI(abstractionURI, hl)).To(Equal(el))
			Expect(len(owner.GetOwnedConceptsRefinedFromURI(abstractionURI, hl))).To(Equal(1))
		})
		Specify("Getting any descendant with abstractionURI", func() {
			el, _ := uOfD.NewElement(hl)
			el2, _ := uOfD.NewElement(hl)
			owner, _ := uOfD.NewElement(hl)
			el.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			el2.SetOwningConceptID(el.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, hl)
			ref, _ := uOfD.NewRefinement(hl)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(el.getConceptIDNoLock(), hl)
			ref2, _ := uOfD.NewRefinement(hl)
			ref2.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			ref2.SetRefinedConceptID(el2.getConceptIDNoLock(), hl)
			Expect(len(owner.GetOwnedDescendantsRefinedFromURI(abstractionURI, hl))).To(Equal(2))
		})
		Specify("Getting Literal child with abstraction", func() {
			lit, _ := uOfD.NewLiteral(hl)
			owner, _ := uOfD.NewElement(hl)
			lit.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			ref, _ := uOfD.NewRefinement(hl)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(lit.getConceptIDNoLock(), hl)
			Expect(owner.GetFirstOwnedLiteralRefinementOf(abs, hl)).To(Equal(lit))
			Expect(len(owner.GetOwnedLiteralsRefinedFrom(abs, hl))).To(Equal(1))
		})
		Specify("Getting Literal child with abstractionURI", func() {
			lit, _ := uOfD.NewLiteral(hl)
			owner, _ := uOfD.NewElement(hl)
			lit.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, hl)
			ref, _ := uOfD.NewRefinement(hl)
			ref.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(lit.getConceptIDNoLock(), hl)
			Expect(owner.GetFirstOwnedLiteralRefinementOfURI(abstractionURI, hl)).To(Equal(lit))
			Expect(len(owner.GetOwnedLiteralsRefinedFromURI(abstractionURI, hl))).To(Equal(1))
		})
		Specify("Getting Reference child with abstraction", func() {
			ref, _ := uOfD.NewReference(hl)
			owner, _ := uOfD.NewElement(hl)
			ref.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			refinement, _ := uOfD.NewRefinement(hl)
			refinement.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			refinement.SetRefinedConceptID(ref.getConceptIDNoLock(), hl)
			Expect(owner.GetFirstOwnedReferenceRefinedFrom(abs, hl)).To(Equal(ref))
			Expect(len(owner.GetOwnedReferencesRefinedFrom(abs, hl))).To(Equal(1))
		})
		Specify("Getting Reference child with abstractionURI", func() {
			ref, _ := uOfD.NewReference(hl)
			owner, _ := uOfD.NewElement(hl)
			ref.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, hl)
			refinement, _ := uOfD.NewRefinement(hl)
			refinement.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			refinement.SetRefinedConceptID(ref.getConceptIDNoLock(), hl)
			Expect(owner.GetFirstOwnedReferenceRefinedFromURI(abstractionURI, hl)).To(Equal(ref))
			Expect(len(owner.GetOwnedReferencesRefinedFromURI(abstractionURI, hl))).To(Equal(1))
		})
		Specify("Getting Refinement child with abstraction", func() {
			ref, _ := uOfD.NewRefinement(hl)
			owner, _ := uOfD.NewElement(hl)
			ref.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			refinement, _ := uOfD.NewRefinement(hl)
			refinement.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			refinement.SetRefinedConceptID(ref.getConceptIDNoLock(), hl)
			Expect(owner.GetFirstOwnedRefinementRefinedFrom(abs, hl)).To(Equal(ref))
			Expect(len(owner.GetOwnedRefinementsRefinedFrom(abs, hl))).To(Equal(1))
		})
		Specify("Getting Refinement child with abstractionURI", func() {
			ref, _ := uOfD.NewRefinement(hl)
			owner, _ := uOfD.NewElement(hl)
			ref.SetOwningConceptID(owner.getConceptIDNoLock(), hl)
			abs, _ := uOfD.NewElement(hl)
			abstractionURI := "http://test.uri"
			abs.SetURI(abstractionURI, hl)
			refinement, _ := uOfD.NewRefinement(hl)
			refinement.SetAbstractConceptID(abs.getConceptIDNoLock(), hl)
			refinement.SetRefinedConceptID(ref.getConceptIDNoLock(), hl)
			Expect(owner.GetFirstOwnedRefinementRefinedFromURI(abstractionURI, hl)).To(Equal(ref))
			Expect(len(owner.GetOwnedRefinementsRefinedFromURI(abstractionURI, hl))).To(Equal(1))
		})
	})

	Describe("Getting children with URI", func() {
		Specify("GetFirstChildWithURI should work", func() {
			owner, _ := uOfD.NewElement(hl)
			child, _ := uOfD.NewElement(hl)
			child.SetOwningConcept(owner, hl)
			uri := "http://test.uri"
			child.SetURI(uri, hl)
			Expect(owner.GetFirstOwnedConceptWithURI(uri, hl)).To(Equal(child))
		})
		Specify("GetFirstChildLiteralWithURI should work", func() {
			owner, _ := uOfD.NewElement(hl)
			child, _ := uOfD.NewLiteral(hl)
			child.SetOwningConcept(owner, hl)
			uri := "http://test.uri"
			child.SetURI(uri, hl)
			Expect(owner.GetFirstOwnedLiteralWithURI(uri, hl)).To(Equal(child))
		})
		Specify("GetFirstChildReferenceWithURI should work", func() {
			owner, _ := uOfD.NewElement(hl)
			child, _ := uOfD.NewReference(hl)
			child.SetOwningConcept(owner, hl)
			uri := "http://test.uri"
			child.SetURI(uri, hl)
			Expect(owner.GetFirstOwnedReferenceWithURI(uri, hl)).To(Equal(child))
		})
		Specify("GetFirstChildRefinementWithURI should work", func() {
			owner, _ := uOfD.NewElement(hl)
			child, _ := uOfD.NewRefinement(hl)
			child.SetOwningConcept(owner, hl)
			uri := "http://test.uri"
			child.SetURI(uri, hl)
			Expect(owner.GetFirstOwnedRefinementWithURI(uri, hl)).To(Equal(child))
		})
	})
})
