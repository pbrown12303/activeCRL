package crlmapsdomain

import (
	"fmt"
	"log"
	"os"

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
	var hl *core.Transaction
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
	var tempDirPath string

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
		mapAbstractDomainOwnedConcepts := mapAbstractDomain.GetOwnedConcepts(hl)
		fmt.Print(mapAbstractDomainOwnedConcepts)
		mapAbstractSourceRef := mapAbstractDomain.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
		Expect(mapAbstractSourceRef.SetReferencedConcept(sourceAbstractDomain, core.NoAttribute, hl)).To(Succeed())
		mapAbstractTargetRef := mapAbstractDomain.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
		Expect(mapAbstractTargetRef.SetReferencedConcept(targetAbstractDomain, core.NoAttribute, hl)).To(Succeed())

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

		// Get the tempDir
		tempDirPath = os.TempDir()
		log.Printf("TempDirPath: " + tempDirPath)
		err = os.Mkdir(tempDirPath, os.ModeDir)
		if !(err == nil || os.IsExist(err)) {
			Expect(err).NotTo(HaveOccurred())
		}
		log.Printf("TempDir created")

	})
	Describe("Target domain creation", func() {
		Specify("The target domain should be created correctly", func() {
			log.Printf("About set set sourceInstanceDomain")
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
		})
	})
	Describe("Individual Concept Mapping - any to any", func() {
		var sourceAbstractElement core.Element
		var sourceAbstractReference core.Reference
		var sourceAbstractRefinement core.Refinement
		var sourceAbstractLiteral core.Literal
		var targetAbstractElement core.Element
		var targetAbstractReference core.Reference
		var targetAbstractRefinement core.Refinement
		var targetAbstractLiteral core.Literal
		BeforeEach(func() {
			var err error
			sourceAbstractElement, err = uOfD.NewOwnedElement(sourceAbstractDomain, "SourceAbstractElement", hl)
			Expect(err).To(BeNil())
			sourceAbstractReference, err = uOfD.NewOwnedReference(sourceAbstractDomain, "SourceAbstractReference", hl)
			Expect(err).To(BeNil())
			sourceAbstractRefinement, err = uOfD.NewOwnedRefinement(sourceAbstractDomain, "SourceAbstractRefinement", hl)
			Expect(err).To(BeNil())
			sourceAbstractLiteral, err = uOfD.NewOwnedLiteral(sourceAbstractDomain, "SourceAbstractLiteral", hl)
			Expect(err).To(BeNil())

			targetAbstractElement, err = uOfD.NewOwnedElement(targetAbstractDomain, "TargetAbstractElement", hl)
			Expect(err).To(BeNil())
			targetAbstractReference, err = uOfD.NewOwnedReference(targetAbstractDomain, "TargetAbstractReference", hl)
			Expect(err).To(BeNil())
			targetAbstractRefinement, err = uOfD.NewOwnedRefinement(targetAbstractDomain, "TargetAbstractRefinement", hl)
			Expect(err).To(BeNil())
			targetAbstractLiteral, err = uOfD.NewOwnedLiteral(targetAbstractDomain, "TargetAbstractLiteral", hl)
			Expect(err).To(BeNil())
			hl.ReleaseLocksAndWait()
		})
		Specify("Element to Element Map", func() {
			// Set up the abstract map
			elementToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(elementToElementMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(elementToElementMap.SetLabel("ElementToElementMap", hl)).To(Succeed())
			Expect(SetSource(elementToElementMap, sourceAbstractElement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(elementToElementMap, targetAbstractElement, core.NoAttribute, hl)).To(Succeed())
			// Add the element to the source instance
			sourceInstanceElement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractElement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceElement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement.SetLabel("SourceInstanceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("ElementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "ElementMapTest")).To(Succeed())

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
		Specify("Reference to Reference Map", func() {
			// Set up the abstract map
			reference2ReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(reference2ReferenceMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(reference2ReferenceMap.SetLabel("Reference2ReferenceMap", hl)).To(Succeed())
			Expect(SetSource(reference2ReferenceMap, sourceAbstractReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(reference2ReferenceMap, targetAbstractReference, core.NoAttribute, hl)).To(Succeed())
			// Add the reference to the source instance
			sourceInstanceReference, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractReference, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("ElementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Troubleshooting aid
			// Expect(graph.ExportDOT(tempDirPath, "ElementMapTest")).To(Succeed())

			// Check the result
			referenceMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceReference, hl)
			Expect(referenceMapInstance).ToNot(BeNil())
			targetInstanceReference := GetTarget(referenceMapInstance, hl)
			Expect(targetInstanceReference).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceReference2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractReference, hl)
			Expect(targetInstanceReference2).ToNot(BeNil())
			Expect(targetInstanceReference2).To(Equal(targetInstanceReference))
		})
		Specify("Literal to Literal Map", func() {
			// Set up the abstract map
			referemce2LiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(referemce2LiteralMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(referemce2LiteralMap.SetLabel("Literal2LiteralMap", hl)).To(Succeed())
			Expect(SetSource(referemce2LiteralMap, sourceAbstractLiteral, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(referemce2LiteralMap, targetAbstractLiteral, core.NoAttribute, hl)).To(Succeed())
			// Add the literal to the source instance
			sourceInstanceLiteral, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractLiteral, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceLiteral.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceLiteral.SetLabel("SourceInstanceLiteral", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("ElementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Troubleshooting aid
			// Expect(graph.ExportDOT(tempDirPath, "ElementMapTest")).To(Succeed())

			// Check the result
			literalMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceLiteral, hl)
			Expect(literalMapInstance).ToNot(BeNil())
			targetInstanceLiteral := GetTarget(literalMapInstance, hl)
			Expect(targetInstanceLiteral).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceLiteral2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractLiteral, hl)
			Expect(targetInstanceLiteral2).ToNot(BeNil())
			Expect(targetInstanceLiteral2).To(Equal(targetInstanceLiteral))
		})
		Specify("Refinement to Refinement Map", func() {
			// Set up the abstract map
			referemce2RefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(referemce2RefinementMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(referemce2RefinementMap.SetLabel("Refinement2RefinementMap", hl)).To(Succeed())
			Expect(SetSource(referemce2RefinementMap, sourceAbstractRefinement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(referemce2RefinementMap, targetAbstractRefinement, core.NoAttribute, hl)).To(Succeed())
			// Add the refinement to the source instance
			sourceInstanceRefinement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractRefinement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceRefinement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceRefinement.SetLabel("SourceInstanceRefinement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("ElementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Troubleshooting aid
			// Expect(graph.ExportDOT(tempDirPath, "ElementMapTest")).To(Succeed())

			// Check the result
			refinementMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceRefinement, hl)
			Expect(refinementMapInstance).ToNot(BeNil())
			targetInstanceRefinement := GetTarget(refinementMapInstance, hl)
			Expect(targetInstanceRefinement).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceRefinement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractRefinement, hl)
			Expect(targetInstanceRefinement2).ToNot(BeNil())
			Expect(targetInstanceRefinement2).To(Equal(targetInstanceRefinement))
		})
		Specify("Element2ReferenceMap", func() {
			// Set up the abstract map
			elementToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(elementToReferenceMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(elementToReferenceMap.SetLabel("ElementToReferenceMap", hl)).To(Succeed())
			Expect(SetSource(elementToReferenceMap, sourceAbstractElement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(elementToReferenceMap, targetAbstractReference, core.NoAttribute, hl)).To(Succeed())
			// Add the element to the source instance
			sourceInstanceElement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractElement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceElement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement.SetLabel("SourceInstanceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Element2ReferenceMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Element2ReferenceMapTest")).To(Succeed())

			// Check the result
			element2ReferenceMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceElement, hl)
			Expect(element2ReferenceMapInstance).ToNot(BeNil())
			targetInstanceReference := GetTarget(element2ReferenceMapInstance, hl)
			Expect(targetInstanceReference).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceReference2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractReference, hl)
			Expect(targetInstanceReference2).ToNot(BeNil())
			Expect(targetInstanceReference2).To(Equal(targetInstanceReference))
		})
		Specify("Element2LiteralMap", func() {
			// Set up the abstract map
			elementToLiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(elementToLiteralMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(elementToLiteralMap.SetLabel("ElementToLiteralMap", hl)).To(Succeed())
			Expect(SetSource(elementToLiteralMap, sourceAbstractElement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(elementToLiteralMap, targetAbstractLiteral, core.NoAttribute, hl)).To(Succeed())
			// Add the element to the source instance
			sourceInstanceElement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractElement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceElement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement.SetLabel("SourceInstanceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Element2LiteralMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Element2LiteralMapTest")).To(Succeed())

			// Check the result
			element2LiteralMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceElement, hl)
			Expect(element2LiteralMapInstance).ToNot(BeNil())
			targetInstanceLiteral := GetTarget(element2LiteralMapInstance, hl)
			Expect(targetInstanceLiteral).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceLiteral2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractLiteral, hl)
			Expect(targetInstanceLiteral2).ToNot(BeNil())
			Expect(targetInstanceLiteral2).To(Equal(targetInstanceLiteral))
		})
		Specify("Element2RefinementMap", func() {
			// Set up the abstract map
			elementToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(elementToRefinementMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(elementToRefinementMap.SetLabel("ElementToRefinementMap", hl)).To(Succeed())
			Expect(SetSource(elementToRefinementMap, sourceAbstractElement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(elementToRefinementMap, targetAbstractRefinement, core.NoAttribute, hl)).To(Succeed())
			// Add the element to the source instance
			sourceInstanceElement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractElement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceElement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement.SetLabel("SourceInstanceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Element2RefinementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Element2RefinementMapTest")).To(Succeed())

			// Check the result
			element2RefinementMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceElement, hl)
			Expect(element2RefinementMapInstance).ToNot(BeNil())
			targetInstanceRefinement := GetTarget(element2RefinementMapInstance, hl)
			Expect(targetInstanceRefinement).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceRefinement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractRefinement, hl)
			Expect(targetInstanceRefinement2).ToNot(BeNil())
			Expect(targetInstanceRefinement2).To(Equal(targetInstanceRefinement))
		})
		Specify("Reference2ElementMap", func() {
			// Set up the abstract map
			referenceToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(referenceToElementMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(referenceToElementMap.SetLabel("ReferenceToElementMap", hl)).To(Succeed())
			Expect(SetSource(referenceToElementMap, sourceAbstractReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(referenceToElementMap, targetAbstractElement, core.NoAttribute, hl)).To(Succeed())
			// Add the reference to the source instance
			sourceInstanceReference, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractReference, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Reference2ElementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Reference2ElementMapTest")).To(Succeed())

			// Check the result
			reference2ElementMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceReference, hl)
			Expect(reference2ElementMapInstance).ToNot(BeNil())
			targetInstanceElement := GetTarget(reference2ElementMapInstance, hl)
			Expect(targetInstanceElement).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceElement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractElement, hl)
			Expect(targetInstanceElement2).ToNot(BeNil())
			Expect(targetInstanceElement2).To(Equal(targetInstanceElement))
		})
		Specify("Reference2LiteralMap", func() {

		})
		Specify("Reference2RefinementMap", func() {
			// Set up the abstract map
			referenceToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(referenceToRefinementMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(referenceToRefinementMap.SetLabel("ReferenceToRefinementMap", hl)).To(Succeed())
			Expect(SetSource(referenceToRefinementMap, sourceAbstractReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(referenceToRefinementMap, targetAbstractRefinement, core.NoAttribute, hl)).To(Succeed())
			// Add the reference to the source instance
			sourceInstanceReference, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractReference, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Reference2RefinementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Reference2RefinementMapTest")).To(Succeed())

			// Check the result
			reference2RefinementMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceReference, hl)
			Expect(reference2RefinementMapInstance).ToNot(BeNil())
			targetInstanceRefinement := GetTarget(reference2RefinementMapInstance, hl)
			Expect(targetInstanceRefinement).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceRefinement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractRefinement, hl)
			Expect(targetInstanceRefinement2).ToNot(BeNil())
			Expect(targetInstanceRefinement2).To(Equal(targetInstanceRefinement))
		})
		Specify("Literal2ElementMap", func() {
			// Set up the abstract map
			literalToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(literalToElementMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(literalToElementMap.SetLabel("LiteralToElementMap", hl)).To(Succeed())
			Expect(SetSource(literalToElementMap, sourceAbstractLiteral, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(literalToElementMap, targetAbstractElement, core.NoAttribute, hl)).To(Succeed())
			// Add the literal to the source instance
			sourceInstanceLiteral, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractLiteral, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceLiteral.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceLiteral.SetLabel("SourceInstanceLiteral", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Literal2ElementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Literal2ElementMapTest")).To(Succeed())

			// Check the result
			literal2ElementMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceLiteral, hl)
			Expect(literal2ElementMapInstance).ToNot(BeNil())
			targetInstanceElement := GetTarget(literal2ElementMapInstance, hl)
			Expect(targetInstanceElement).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceElement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractElement, hl)
			Expect(targetInstanceElement2).ToNot(BeNil())
			Expect(targetInstanceElement2).To(Equal(targetInstanceElement))
		})
		Specify("Literal2ReferenceMap", func() {
			// Set up the abstract map
			literalToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(literalToReferenceMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(literalToReferenceMap.SetLabel("LiteralToReferenceMap", hl)).To(Succeed())
			Expect(SetSource(literalToReferenceMap, sourceAbstractLiteral, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(literalToReferenceMap, targetAbstractReference, core.NoAttribute, hl)).To(Succeed())
			// Add the literal to the source instance
			sourceInstanceLiteral, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractLiteral, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceLiteral.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceLiteral.SetLabel("SourceInstanceLiteral", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Literal2ReferenceMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Literal2ReferenceMapTest")).To(Succeed())

			// Check the result
			literal2ReferenceMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceLiteral, hl)
			Expect(literal2ReferenceMapInstance).ToNot(BeNil())
			targetInstanceReference := GetTarget(literal2ReferenceMapInstance, hl)
			Expect(targetInstanceReference).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceReference2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractReference, hl)
			Expect(targetInstanceReference2).ToNot(BeNil())
			Expect(targetInstanceReference2).To(Equal(targetInstanceReference))
		})
		Specify("Literal2RefinementMap", func() {
			// Set up the abstract map
			literalToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(literalToRefinementMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(literalToRefinementMap.SetLabel("LiteralToRefinementMap", hl)).To(Succeed())
			Expect(SetSource(literalToRefinementMap, sourceAbstractLiteral, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(literalToRefinementMap, targetAbstractRefinement, core.NoAttribute, hl)).To(Succeed())
			// Add the literal to the source instance
			sourceInstanceLiteral, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractLiteral, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceLiteral.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceLiteral.SetLabel("SourceInstanceLiteral", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Literal2RefinementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Literal2RefinementMapTest")).To(Succeed())

			// Check the result
			literal2RefinementMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceLiteral, hl)
			Expect(literal2RefinementMapInstance).ToNot(BeNil())
			targetInstanceRefinement := GetTarget(literal2RefinementMapInstance, hl)
			Expect(targetInstanceRefinement).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceRefinement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractRefinement, hl)
			Expect(targetInstanceRefinement2).ToNot(BeNil())
			Expect(targetInstanceRefinement2).To(Equal(targetInstanceRefinement))
		})
		Specify("Refinement2ElementMap", func() {
			// Set up the abstract map
			refinementToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(refinementToElementMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(refinementToElementMap.SetLabel("RefinementToElementMap", hl)).To(Succeed())
			Expect(SetSource(refinementToElementMap, sourceAbstractRefinement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(refinementToElementMap, targetAbstractElement, core.NoAttribute, hl)).To(Succeed())
			// Add the refinement to the source instance
			sourceInstanceRefinement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractRefinement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceRefinement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceRefinement.SetLabel("SourceInstanceRefinement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Refinement2ElementMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Refinement2ElementMapTest")).To(Succeed())

			// Check the result
			refinement2ElementMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceRefinement, hl)
			Expect(refinement2ElementMapInstance).ToNot(BeNil())
			targetInstanceElement := GetTarget(refinement2ElementMapInstance, hl)
			Expect(targetInstanceElement).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceElement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractElement, hl)
			Expect(targetInstanceElement2).ToNot(BeNil())
			Expect(targetInstanceElement2).To(Equal(targetInstanceElement))
		})
		Specify("Refinement2ReferenceMap", func() {
			// Set up the abstract map
			refinementToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(refinementToReferenceMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(refinementToReferenceMap.SetLabel("RefinementToReferenceMap", hl)).To(Succeed())
			Expect(SetSource(refinementToReferenceMap, sourceAbstractRefinement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(refinementToReferenceMap, targetAbstractReference, core.NoAttribute, hl)).To(Succeed())
			// Add the refinement to the source instance
			sourceInstanceRefinement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractRefinement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceRefinement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceRefinement.SetLabel("SourceInstanceRefinement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Refinement2ReferenceMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Refinement2ReferenceMapTest")).To(Succeed())

			// Check the result
			refinement2ReferenceMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceRefinement, hl)
			Expect(refinement2ReferenceMapInstance).ToNot(BeNil())
			targetInstanceReference := GetTarget(refinement2ReferenceMapInstance, hl)
			Expect(targetInstanceReference).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceReference2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractReference, hl)
			Expect(targetInstanceReference2).ToNot(BeNil())
			Expect(targetInstanceReference2).To(Equal(targetInstanceReference))
		})
		Specify("Refinement2LiteralMap", func() {
			// Set up the abstract map
			refinementToLiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(refinementToLiteralMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(refinementToLiteralMap.SetLabel("RefinementToLiteralMap", hl)).To(Succeed())
			Expect(SetSource(refinementToLiteralMap, sourceAbstractRefinement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(refinementToLiteralMap, targetAbstractLiteral, core.NoAttribute, hl)).To(Succeed())
			// Add the refinement to the source instance
			sourceInstanceRefinement, err2 := uOfD.CreateReplicateAsRefinement(sourceAbstractRefinement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceRefinement.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceRefinement.SetLabel("SourceInstanceRefinement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, core.NoAttribute, hl)).To(Succeed())
			hl.ReleaseLocksAndWait()

			// Diagnostics view
			// graph := core.NewCrlGraph("Refinement2LiteralMapTest")
			// Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
			// Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
			// hl.ReleaseLocksAndWait()
			// Expect(graph.ExportDOT(tempDirPath, "Refinement2LiteralMapTest")).To(Succeed())

			// Check the result
			refinement2LiteralMapInstance := FindMapForSource(mapInstanceDomain, sourceInstanceRefinement, hl)
			Expect(refinement2LiteralMapInstance).ToNot(BeNil())
			targetInstanceLiteral := GetTarget(refinement2LiteralMapInstance, hl)
			Expect(targetInstanceLiteral).ToNot(BeNil())
			targetInstanceDomain := mapInstanceFolder.GetFirstOwnedConceptRefinedFrom(targetAbstractDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceLiteral2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(targetAbstractLiteral, hl)
			Expect(targetInstanceLiteral2).ToNot(BeNil())
			Expect(targetInstanceLiteral2).To(Equal(targetInstanceLiteral))
		})
	})
	Describe("Individual Concept Mapping - any to any", func() {
		var sourceAbstractElement1 core.Element
		var sourceAbstractElement2 core.Element
		var sourceAbstractReference core.Reference
		// var sourceAbstractRefinement core.Refinement
		// var sourceAbstractLiteral core.Literal
		var targetAbstractElement1 core.Element
		var targetAbstractElement2 core.Element
		var targetAbstractReference core.Reference
		// var targetAbstractRefinement core.Refinement
		// var targetAbstractLiteral core.Literal

		var element1ToElement1Map core.Element
		var element2ToElement2Map core.Element
		var reference2ReferenceMap core.Element
		var sourceInstanceElement1 core.Element
		var sourceInstanceElement2 core.Element
		var sourceInstanceReference core.Reference

		BeforeEach(func() {
			var err error
			sourceAbstractElement1, err = uOfD.NewOwnedElement(sourceAbstractDomain, "SourceAbstractElement1", hl)
			Expect(err).To(BeNil())
			sourceAbstractElement2, err = uOfD.NewOwnedElement(sourceAbstractDomain, "SourceAbstractElement2", hl)
			Expect(err).To(BeNil())
			sourceAbstractReference, err = uOfD.NewOwnedReference(sourceAbstractElement1, "SourceAbstractReference", hl)
			Expect(err).To(BeNil())
			Expect(sourceAbstractReference.SetReferencedConcept(sourceAbstractElement2, core.NoAttribute, hl)).To(Succeed())
			// sourceAbstractRefinement, err = uOfD.NewOwnedRefinement(sourceAbstractDomain, "SourceAbstractRefinement", hl)
			// Expect(err).To(BeNil())
			// sourceAbstractLiteral, err = uOfD.NewOwnedLiteral(sourceAbstractDomain, "SourceAbstractLiteral", hl)
			// Expect(err).To(BeNil())

			targetAbstractElement1, err = uOfD.NewOwnedElement(targetAbstractDomain, "TargetAbstractElement1", hl)
			Expect(err).To(BeNil())
			targetAbstractElement2, err = uOfD.NewOwnedElement(targetAbstractDomain, "TargetAbstractElement2", hl)
			Expect(err).To(BeNil())
			targetAbstractReference, err = uOfD.NewOwnedReference(targetAbstractElement2, "TargetAbstractReference", hl)
			Expect(err).To(BeNil())
			Expect(targetAbstractReference.SetReferencedConcept(targetAbstractElement1, core.NoAttribute, hl)).To(Succeed())
			// targetAbstractRefinement, err = uOfD.NewOwnedRefinement(targetAbstractDomain, "TargetAbstractRefinement", hl)
			// Expect(err).To(BeNil())
			// targetAbstractLiteral, err = uOfD.NewOwnedLiteral(targetAbstractDomain, "TargetAbstractLiteral", hl)
			// Expect(err).To(BeNil())
			hl.ReleaseLocksAndWait()

			// Abstract Map Setup
			element1ToElement1Map, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(element1ToElement1Map.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(element1ToElement1Map.SetLabel("Element1ToElement1Map", hl)).To(Succeed())
			Expect(SetSource(element1ToElement1Map, sourceAbstractElement1, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(element1ToElement1Map, targetAbstractElement1, core.NoAttribute, hl)).To(Succeed())

			element2ToElement2Map, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(element2ToElement2Map.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(element2ToElement2Map.SetLabel("Element2ToElement2Map", hl)).To(Succeed())
			Expect(SetSource(element2ToElement2Map, sourceAbstractElement2, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(element2ToElement2Map, targetAbstractElement2, core.NoAttribute, hl)).To(Succeed())

			reference2ReferenceMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(reference2ReferenceMap.SetOwningConcept(mapAbstractDomain, hl)).To(Succeed())
			Expect(reference2ReferenceMap.SetLabel("Reference2ReferenceMap", hl)).To(Succeed())
			Expect(SetSource(reference2ReferenceMap, sourceAbstractReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(reference2ReferenceMap, targetAbstractReference, core.NoAttribute, hl)).To(Succeed())

			// Source Instance Setup
			sourceInstanceElement1, err = uOfD.CreateReplicateAsRefinement(sourceAbstractElement1, hl)
			Expect(err).To(BeNil())
			Expect(sourceInstanceElement1.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement1.SetLabel("SourceInstanceElement1", hl)).To(Succeed())

			sourceInstanceElement2, err = uOfD.CreateReplicateAsRefinement(sourceAbstractElement2, hl)
			Expect(err).To(BeNil())
			Expect(sourceInstanceElement2.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement2.SetLabel("SourceInstanceElement2", hl)).To(Succeed())

			sourceInstanceReference, err = uOfD.CreateReplicateReferenceAsRefinement(sourceAbstractReference, hl)
			Expect(err).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(sourceInstanceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())

			// Map Instance Setup
		})
		// FSpecify("Pointer to Reference Pointer", func() {
		// 	// Add the pointer map
		// 	pointer2ReferencePointerMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
		// 	Expect(err).To(BeNil())
		// 	Expect(pointer2ReferencePointerMap.SetOwningConcept(reference2ReferenceMap, hl)).To(Succeed())
		// 	Expect(pointer2ReferencePointerMap.SetLabel("Pointer2ReferencePointerMap", hl)).To(Succeed())
		// 	Expect(SetSource(pointer2ReferencePointerMap, sourceAbstractReference, hl)).To(Succeed())
		// 	Expect(SetSourceAttributeName(pointer2ReferencePointerMap, core.OwningConceptID, hl)).To(Succeed())
		// 	Expect(SetTarget(pointer2ReferencePointerMap, targetAbstractReference, hl)).To(Succeed())
		// 	Expect(SetTargetAttributeName(pointer2ReferencePointerMap, core.ReferencedConceptID, hl)).To(Succeed())

		// 	// Trigger the map
		// 	Expect(SetSource(mapInstanceDomain, sourceInstanceDomain, hl)).To(Succeed())
		// 	hl.ReleaseLocksAndWait()

		// 	// Diagnostics view
		// 	graph := core.NewCrlGraph("Pointer2RferencePointerMapTest")
		// 	Expect(graph.AddConceptRecursively(sourceAbstractFolder, hl)).To(Succeed())
		// 	Expect(graph.AddConceptRecursively(targetAbstractFolder, hl)).To(Succeed())
		// 	Expect(graph.AddConceptRecursively(mapAbstractFolder, hl)).To(Succeed())
		// 	Expect(graph.AddConceptRecursively(sourceInstanceFolder, hl)).To(Succeed())
		// 	Expect(graph.AddConceptRecursively(mapInstanceFolder, hl)).To(Succeed())
		// 	hl.ReleaseLocksAndWait()
		// 	Expect(graph.ExportDOT(tempDirPath, "Pointer2RferencePointerMapTest")).To(Succeed())

		// 	// Check the result
		// 	sourceInstanceReferenceMap := FindMapForSource(mapInstanceDomain, sourceInstanceReference, hl)
		// 	Expect(sourceInstanceReferenceMap).ToNot(BeNil())
		// 	sourceInstanceReferenceAttributeMap := FindMapForSourceAttribute(mapInstanceDomain, sourceInstanceReference, core.OwningConceptID, hl)
		// 	Expect(sourceInstanceReferenceAttributeMap).ToNot(BeNil())
		// 	Expect(sourceInstanceReferenceAttributeMap.GetOwningConcept(hl).GetConceptID(hl)).To(Equal(sourceInstanceReferenceMap.GetConceptID(hl)))
		// 	targetInstance := GetTarget(sourceInstanceReferenceMap, hl)
		// 	Expect(targetInstance).ToNot(BeNil())
		// 	targetInstanceAttribute := GetTarget(sourceInstanceReferenceAttributeMap, hl)
		// 	Expect(targetInstanceAttribute).To(Equal(targetInstance))

		// 	targetInstanceElement1Map := FindMapForSource(mapInstanceDomain, sourceInstanceElement1, hl)
		// 	Expect(targetInstanceElement1Map).ToNot(BeNil())
		// 	targetInstanceElement1 := GetTarget(targetInstanceElement1Map, hl)
		// 	Expect(targetInstanceElement1).ToNot(BeNil())

		// 	switch targetInstance.(type) {
		// 	case core.Reference:
		// 		Expect(targetInstance.(core.Reference).GetReferencedConcept(hl).GetConceptID(hl)).To(Equal(targetInstanceElement1.GetConceptID(hl)))
		// 	}
		// })
		Specify("Pointer to Reference Owner", func() {

		})
		Specify("Pointer to Refinement Abstract Pointer", func() {

		})
		Specify("Pointer to Refinement Refined Pointer", func() {

		})
		Specify("Pointer to Literal Value", func() {

		})
	})
})
