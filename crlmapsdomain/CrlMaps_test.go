package crlmapsdomain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("CrlMaps domain test", func() {
	Specify("Domain generation should be idempotent", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewHeldLocks()
		BuildCrlMapsDomain(uOfD1, hl1)
		cs1 := uOfD1.GetElementWithURI(CrlMapsDomainURI)
		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewHeldLocks()
		BuildCrlMapsDomain(uOfD2, hl2)
		cs2 := uOfD2.GetElementWithURI(CrlMapsDomainURI)
		Expect(core.RecursivelyEquivalent(cs1, hl1, cs2, hl2)).To(BeTrue())
	})
	Specify("Each URI should have an associated Element or Reference", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewHeldLocks()
		BuildCrlMapsDomain(uOfD1, hl1)
		domain := uOfD1.GetElementWithURI(CrlMapsDomainURI)
		Expect(domain).ShouldNot(BeNil())
		eToEMap := uOfD1.GetElementWithURI(CrlOneToOneMapURI)
		Expect(eToEMap).ShouldNot(BeNil())
		eToEMapSource := uOfD1.GetReferenceWithURI(CrlOneToOneMapSourceReferenceURI)
		Expect(eToEMapSource).ShouldNot(BeNil())
		eToEMapTarget := uOfD1.GetReferenceWithURI(CrlOneToOneMapTargetReferenceURI)
		Expect(eToEMapTarget).ShouldNot(BeNil())
	})
})

var _ = Describe("CrlMaps mapping tests", func() {
	var uOfD *core.UniverseOfDiscourse
	var hl *core.HeldLocks
	var sourceAbstractFolder core.Element
	var sourceAbstractDomain core.Element
	var targetAbstractFolder core.Element
	var targetAbstractDomain core.Element
	var mapAbstractFolder core.Element
	var mapAbstractDomain core.Element
	var sourceInstanceFolder core.Element
	var sourceInstanceDomain core.Element
	var mapInstanceFolder core.Element
	var mapInstanceDomain core.Element
	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
		BuildCrlMapsDomain(uOfD, hl)
		var err error

		// Abstract Source
		sourceAbstractFolder, err = uOfD.NewOwnedElement(nil, "SourceAbstractFolder", hl)
		Expect(err).To(BeNil())
		sourceAbstractDomain, err = uOfD.NewOwnedElement(sourceAbstractFolder, "SourceAbstractDomain", hl)
		Expect(err).To(BeNil())

		// Abstract Target
		targetAbstractFolder, err = uOfD.NewOwnedElement(nil, "TargetAbstractFolder", hl)
		Expect(err).To(BeNil())
		targetAbstractDomain, err = uOfD.NewOwnedElement(targetAbstractFolder, "TargetAbstractDomain", hl)
		Expect(err).To(BeNil())

		// Abstract Map
		mapAbstractFolder, err = uOfD.NewOwnedElement(nil, "MapAbstractFolder", hl)
		Expect(err).To(BeNil())
		mapAbstractDomain, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
		Expect(err).To(BeNil())
		Expect(mapAbstractDomain.SetLabel("MapAbstractDomain", hl)).To(Succeed())
		Expect(mapAbstractDomain.SetOwningConcept(mapAbstractFolder, hl)).To(Succeed())
		mapAbstractSourceRef := mapAbstractDomain.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
		Expect(mapAbstractSourceRef.SetReferencedConcept(sourceAbstractDomain, hl)).To(Succeed())
		mapAbstractTargetRef := mapAbstractDomain.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
		Expect(mapAbstractTargetRef.SetReferencedConcept(targetAbstractDomain, hl)).To(Succeed())

		// Source Instance
		sourceInstanceFolder, err = uOfD.NewOwnedElement(nil, "sourceInstanceFolder", hl)
		Expect(err).To(BeNil())
		sourceInstanceDomain, err = uOfD.CreateReplicateAsRefinement(sourceAbstractDomain, hl)
		Expect(err).To(BeNil())
		Expect(sourceInstanceDomain.SetLabel("SourceInstanceDomain", hl)).To(Succeed())
		Expect(sourceInstanceDomain.SetOwningConcept(sourceInstanceFolder, hl)).To(Succeed())

		// Map Instance
		mapInstanceFolder, err = uOfD.NewOwnedElement(nil, "mapInstanceFolder", hl)
		Expect(err).To(BeNil())
		mapInstanceDomain, err = uOfD.CreateReplicateAsRefinement(mapAbstractDomain, hl)
		Expect(err).To(BeNil())
		Expect(mapInstanceDomain.SetLabel("MapInstanceDomain", hl)).To(Succeed())
		Expect(mapInstanceDomain.SetOwningConcept(mapInstanceFolder, hl)).To(Succeed())
		hl.ReleaseLocksAndWait()
	})
	Describe("Target domain creation", func() {
		Specify("The target domain should be created correctly", func() {
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
		})
	})
	FDescribe("Individual Concept Mapping - any to any", func() {
		var sourceAbstractElement core.Element
		// var sourceAbstractReference core.Reference
		// var sourceAbstractRefinement core.Refinement
		// var sourceAbstractLiteral core.Literal
		var targetAbstractElement core.Element
		// var targetAbstractReference core.Reference
		// var targetAbstractRefinement core.Refinement
		// var targetAbstractLiteral core.Literal
		BeforeEach(func() {
			var err error
			sourceAbstractElement, err = uOfD.NewOwnedElement(sourceAbstractDomain, "SourceAbstractElement", hl)
			Expect(err).To(BeNil())
			// sourceAbstractReference, err = uOfD.NewOwnedReference(sourceAbstractDomain, "SourceAbstractReference", hl)
			// Expect(err).To(BeNil())
			// sourceAbstractRefinement, err = uOfD.NewOwnedRefinement(sourceAbstractDomain, "SourceAbstractRefinement", hl)
			// Expect(err).To(BeNil())
			// sourceAbstractLiteral, err = uOfD.NewOwnedLiteral(sourceAbstractDomain, "SourceAbstractLiteral", hl)
			// Expect(err).To(BeNil())

			targetAbstractElement, err = uOfD.NewOwnedElement(targetAbstractDomain, "TargetAbstractElement", hl)
			Expect(err).To(BeNil())
			// targetAbstractReference, err = uOfD.NewOwnedReference(targetAbstractDomain, "TargetAbstractReference", hl)
			// Expect(err).To(BeNil())
			// targetAbstractRefinement, err = uOfD.NewOwnedRefinement(targetAbstractDomain, "TargetAbstractRefinement", hl)
			// Expect(err).To(BeNil())
			// targetAbstractLiteral, err = uOfD.NewOwnedLiteral(targetAbstractDomain, "TargetAbstractLiteral", hl)
			// Expect(err).To(BeNil())
			hl.ReleaseLocksAndWait()
		})
		FSpecify("Element to Element Map", func() {
			// Set up the abstract map
			elementToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(elementToElementMap.SetOwningConcept(mapInstanceDomain, hl)).To(Succeed())
			Expect(elementToElementMap.SetLabel("ElementToElementMap", hl)).To(Succeed())
			Expect(SetSource(elementToElementMap, sourceAbstractElement, hl)).To(Succeed())
			Expect(SetTarget(elementToElementMap, targetAbstractElement, hl)).To(Succeed())
			// Add the element to the source instance
			sourceInstanceElement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractElement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceElement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement.SetLabel("SourceInstanceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()
			// Check the result
			elementMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceElement, hl)
			Expect(elementMapInstance).ToNot(BeNil())
			targetInstanceElement := GetTarget(elementMapInstance, hl)
			Expect(targetInstanceElement).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceElement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractElement, hl)
			Expect(targetInstanceElement2).ToNot(BeNil())
			Expect(targetInstanceElement2).To(Equal(targetInstanceElement))
		})
	})
})
