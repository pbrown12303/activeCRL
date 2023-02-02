package crlmapsdomain

import (
	"os"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("CrlMaps domain test", func() {
	Specify("Domain generation should be idempotent", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewTransaction()
		BuildCrlMapsDomain(uOfD1, hl1)
		cs1 := uOfD1.GetElementWithURI(CrlMapsDomainURI)
		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewTransaction()
		BuildCrlMapsDomain(uOfD2, hl2)
		cs2 := uOfD2.GetElementWithURI(CrlMapsDomainURI)
		Expect(core.RecursivelyEquivalent(cs1, hl1, cs2, hl2)).To(BeTrue())
	})
	Specify("Each URI should have an associated Element or Reference", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewTransaction()
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
	var definingSourceFolder core.Element
	var definingSourceDomain core.Element
	var definingTargetFolder core.Element
	var definingTargetDomain core.Element
	var definingMapFolder core.Element
	var definingDomainMap core.Element
	var instanceSourceFolder core.Element
	var instanceSourceDomain core.Element
	var instanceMapFolder core.Element
	var instanceDomainMap core.Element
	var tempDirPath string

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		hl = uOfD.NewTransaction()
		BuildCrlMapsDomain(uOfD, hl)
		var err error

		// Defining Source
		definingSourceFolder, err = uOfD.NewOwnedElement(nil, "DefiningSourceFolder", hl)
		Expect(err).To(BeNil())
		definingSourceDomain, err = uOfD.NewOwnedElement(definingSourceFolder, "DefiningSourceDomain", hl)
		Expect(err).To(BeNil())

		// Defining Target
		definingTargetFolder, err = uOfD.NewOwnedElement(nil, "DefiningTargetFolder", hl)
		Expect(err).To(BeNil())
		definingTargetDomain, err = uOfD.NewOwnedElement(definingTargetFolder, "DefiningTargetDomain", hl)
		Expect(err).To(BeNil())

		// Defining Map
		definingMapFolder, err = uOfD.NewOwnedElement(nil, "DefiningMapFolder", hl)
		Expect(err).To(BeNil())
		definingDomainMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
		Expect(err).To(BeNil())
		Expect(definingDomainMap.SetLabel("DefiningDomainMap", hl)).To(Succeed())
		Expect(definingDomainMap.SetOwningConcept(definingMapFolder, hl)).To(Succeed())
		// mapAbstractDomainOwnedConcepts := mapAbstractDomain.GetOwnedConcepts(hl)
		// fmt.Print(mapAbstractDomainOwnedConcepts)
		definingDomainMapSourceRef := definingDomainMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
		Expect(definingDomainMapSourceRef.SetReferencedConcept(definingSourceDomain, core.NoAttribute, hl)).To(Succeed())
		definingDomainMapTargetRef := definingDomainMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
		Expect(definingDomainMapTargetRef.SetReferencedConcept(definingTargetDomain, core.NoAttribute, hl)).To(Succeed())

		// Source Instance
		instanceSourceFolder, err = uOfD.NewOwnedElement(nil, "instanceSourceFolder", hl)
		Expect(err).To(BeNil())
		instanceSourceDomain, err = uOfD.CreateReplicateAsRefinement(definingSourceDomain, hl)
		Expect(err).To(BeNil())
		Expect(instanceSourceDomain.SetLabel("InstanceSourceDomain", hl)).To(Succeed())
		Expect(instanceSourceDomain.SetOwningConcept(instanceSourceFolder, hl)).To(Succeed())

		// Map Instance
		instanceMapFolder, err = uOfD.NewOwnedElement(nil, "instanceMapFolder", hl)
		Expect(err).To(BeNil())
		instanceDomainMap, err = uOfD.CreateReplicateAsRefinement(definingDomainMap, hl)
		Expect(err).To(BeNil())
		Expect(instanceDomainMap.SetLabel("InstanceDomainMap", hl)).To(Succeed())
		Expect(instanceDomainMap.SetOwningConcept(instanceMapFolder, hl)).To(Succeed())

		// Get the tempDir
		// This will be the location for the graphs generated for debugging purposes
		tempDirPath = os.TempDir()
		err = os.Mkdir(tempDirPath, os.ModeDir)
		if !(err == nil || os.IsExist(err)) {
			Expect(err).NotTo(HaveOccurred())
		}
	})
	Describe("Target domain creation", func() {
		Specify("The target domain should be created correctly", func() {
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
		})
	})
	Describe("Individual Concept Mapping - any to any", func() {
		var definingSourceElement core.Element
		var definingSourceReference core.Reference
		var definingSourceRefinement core.Refinement
		var definingSourceLiteral core.Literal
		var definingTargetElement core.Element
		var definingTargetReference core.Reference
		var definingTargetRefinement core.Refinement
		var definingTargetLiteral core.Literal
		BeforeEach(func() {
			var err error
			definingSourceElement, err = uOfD.NewOwnedElement(definingSourceDomain, "DefiningSourceElement", hl)
			Expect(err).To(BeNil())
			definingSourceReference, err = uOfD.NewOwnedReference(definingSourceDomain, "DefiningSourceReference", hl)
			Expect(err).To(BeNil())
			definingSourceRefinement, err = uOfD.NewOwnedRefinement(definingSourceDomain, "DefiningSourceRefinement", hl)
			Expect(err).To(BeNil())
			definingSourceLiteral, err = uOfD.NewOwnedLiteral(definingSourceDomain, "DefiningSourceLiteral", hl)
			Expect(err).To(BeNil())

			definingTargetElement, err = uOfD.NewOwnedElement(definingTargetDomain, "DefiningTargetElement", hl)
			Expect(err).To(BeNil())
			definingTargetReference, err = uOfD.NewOwnedReference(definingTargetDomain, "DefiningTargetReference", hl)
			Expect(err).To(BeNil())
			definingTargetRefinement, err = uOfD.NewOwnedRefinement(definingTargetDomain, "DefiningTargetRefinement", hl)
			Expect(err).To(BeNil())
			definingTargetLiteral, err = uOfD.NewOwnedLiteral(definingTargetDomain, "DefiningTargetLiteral", hl)
			Expect(err).To(BeNil())
		})
		Specify("Element to Element Map", func() {
			// Set up the abstract map
			definingElementToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(definingElementToElementMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(definingElementToElementMap.SetLabel("DefiningElementToElementMap", hl)).To(Succeed())
			Expect(SetSource(definingElementToElementMap, definingSourceElement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(definingElementToElementMap, definingTargetElement, core.NoAttribute, hl)).To(Succeed())
			// Add the element to the source instance
			instanceSourceElement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceElement, hl)
			Expect(err2).To(BeNil())
			Expect(instanceSourceElement.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(instanceSourceElement.SetLabel("InstanceSourceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			instanceElementToElementMap := FindMapForSource(instanceDomainMap, instanceSourceElement, hl)
			Expect(instanceElementToElementMap).ToNot(BeNil())
			instanceTargetElement := GetTarget(instanceElementToElementMap, hl)
			Expect(instanceTargetElement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetElement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetElement, hl)
			Expect(instanceTargetElement2).ToNot(BeNil())
			Expect(instanceTargetElement2).To(Equal(instanceTargetElement))
		})
		Specify("Reference to Reference Map", func() {
			// Set up the abstract map
			reference2ReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(reference2ReferenceMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(reference2ReferenceMap.SetLabel("Reference2ReferenceMap", hl)).To(Succeed())
			Expect(SetSource(reference2ReferenceMap, definingSourceReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(reference2ReferenceMap, definingTargetReference, core.NoAttribute, hl)).To(Succeed())
			// Add the reference to the source instance
			sourceInstanceReference, err2 := uOfD.CreateReplicateAsRefinement(definingSourceReference, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			referenceMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceReference, hl)
			Expect(referenceMapInstance).ToNot(BeNil())
			targetInstanceReference := GetTarget(referenceMapInstance, hl)
			Expect(targetInstanceReference).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceReference2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetReference, hl)
			Expect(targetInstanceReference2).ToNot(BeNil())
			Expect(targetInstanceReference2).To(Equal(targetInstanceReference))
		})
		Specify("Literal to Literal Map", func() {
			// Set up the abstract map
			referemce2LiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(referemce2LiteralMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(referemce2LiteralMap.SetLabel("Literal2LiteralMap", hl)).To(Succeed())
			Expect(SetSource(referemce2LiteralMap, definingSourceLiteral, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(referemce2LiteralMap, definingTargetLiteral, core.NoAttribute, hl)).To(Succeed())
			// Add the literal to the source instance
			sourceInstanceLiteral, err2 := uOfD.CreateReplicateAsRefinement(definingSourceLiteral, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceLiteral.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceLiteral.SetLabel("SourceInstanceLiteral", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			literalMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceLiteral, hl)
			Expect(literalMapInstance).ToNot(BeNil())
			targetInstanceLiteral := GetTarget(literalMapInstance, hl)
			Expect(targetInstanceLiteral).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceLiteral2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetLiteral, hl)
			Expect(targetInstanceLiteral2).ToNot(BeNil())
			Expect(targetInstanceLiteral2).To(Equal(targetInstanceLiteral))
		})
		Specify("Refinement to Refinement Map", func() {
			// Set up the abstract map
			referemce2RefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(referemce2RefinementMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(referemce2RefinementMap.SetLabel("Refinement2RefinementMap", hl)).To(Succeed())
			Expect(SetSource(referemce2RefinementMap, definingSourceRefinement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(referemce2RefinementMap, definingTargetRefinement, core.NoAttribute, hl)).To(Succeed())
			// Add the refinement to the source instance
			sourceInstanceRefinement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceRefinement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceRefinement.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceRefinement.SetLabel("SourceInstanceRefinement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			refinementMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceRefinement, hl)
			Expect(refinementMapInstance).ToNot(BeNil())
			targetInstanceRefinement := GetTarget(refinementMapInstance, hl)
			Expect(targetInstanceRefinement).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceRefinement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetRefinement, hl)
			Expect(targetInstanceRefinement2).ToNot(BeNil())
			Expect(targetInstanceRefinement2).To(Equal(targetInstanceRefinement))
		})
		Specify("Element2ReferenceMap", func() {
			// Set up the abstract map
			elementToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(elementToReferenceMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(elementToReferenceMap.SetLabel("ElementToReferenceMap", hl)).To(Succeed())
			Expect(SetSource(elementToReferenceMap, definingSourceElement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(elementToReferenceMap, definingTargetReference, core.NoAttribute, hl)).To(Succeed())
			// Add the element to the source instance
			sourceInstanceElement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceElement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceElement.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement.SetLabel("SourceInstanceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			element2ReferenceMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceElement, hl)
			Expect(element2ReferenceMapInstance).ToNot(BeNil())
			targetInstanceReference := GetTarget(element2ReferenceMapInstance, hl)
			Expect(targetInstanceReference).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceReference2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetReference, hl)
			Expect(targetInstanceReference2).ToNot(BeNil())
			Expect(targetInstanceReference2).To(Equal(targetInstanceReference))
		})
		Specify("Element2LiteralMap", func() {
			// Set up the abstract map
			elementToLiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(elementToLiteralMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(elementToLiteralMap.SetLabel("ElementToLiteralMap", hl)).To(Succeed())
			Expect(SetSource(elementToLiteralMap, definingSourceElement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(elementToLiteralMap, definingTargetLiteral, core.NoAttribute, hl)).To(Succeed())
			// Add the element to the source instance
			sourceInstanceElement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceElement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceElement.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement.SetLabel("SourceInstanceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			element2LiteralMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceElement, hl)
			Expect(element2LiteralMapInstance).ToNot(BeNil())
			targetInstanceLiteral := GetTarget(element2LiteralMapInstance, hl)
			Expect(targetInstanceLiteral).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceLiteral2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetLiteral, hl)
			Expect(targetInstanceLiteral2).ToNot(BeNil())
			Expect(targetInstanceLiteral2).To(Equal(targetInstanceLiteral))
		})
		Specify("Element2RefinementMap", func() {
			// Set up the abstract map
			elementToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(elementToRefinementMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(elementToRefinementMap.SetLabel("ElementToRefinementMap", hl)).To(Succeed())
			Expect(SetSource(elementToRefinementMap, definingSourceElement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(elementToRefinementMap, definingTargetRefinement, core.NoAttribute, hl)).To(Succeed())
			// Add the element to the source instance
			sourceInstanceElement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceElement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceElement.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceElement.SetLabel("SourceInstanceElement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			element2RefinementMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceElement, hl)
			Expect(element2RefinementMapInstance).ToNot(BeNil())
			targetInstanceRefinement := GetTarget(element2RefinementMapInstance, hl)
			Expect(targetInstanceRefinement).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceRefinement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetRefinement, hl)
			Expect(targetInstanceRefinement2).ToNot(BeNil())
			Expect(targetInstanceRefinement2).To(Equal(targetInstanceRefinement))
		})
		Specify("Reference2ElementMap", func() {
			// Set up the abstract map
			referenceToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(referenceToElementMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(referenceToElementMap.SetLabel("ReferenceToElementMap", hl)).To(Succeed())
			Expect(SetSource(referenceToElementMap, definingSourceReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(referenceToElementMap, definingTargetElement, core.NoAttribute, hl)).To(Succeed())
			// Add the reference to the source instance
			sourceInstanceReference, err2 := uOfD.CreateReplicateAsRefinement(definingSourceReference, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			reference2ElementMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceReference, hl)
			Expect(reference2ElementMapInstance).ToNot(BeNil())
			targetInstanceElement := GetTarget(reference2ElementMapInstance, hl)
			Expect(targetInstanceElement).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceElement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetElement, hl)
			Expect(targetInstanceElement2).ToNot(BeNil())
			Expect(targetInstanceElement2).To(Equal(targetInstanceElement))
		})
		Specify("Reference2LiteralMap", func() {
			// Set up the abstract map
			reference2LiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(reference2LiteralMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(reference2LiteralMap.SetLabel("Reference2LiteralMap", hl)).To(Succeed())
			Expect(SetSource(reference2LiteralMap, definingSourceReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(reference2LiteralMap, definingTargetLiteral, core.NoAttribute, hl)).To(Succeed())
			// Add the reference to the source instance
			sourceInstanceReference, err2 := uOfD.CreateReplicateAsRefinement(definingSourceReference, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			reference2LiteralMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceReference, hl)
			Expect(reference2LiteralMapInstance).ToNot(BeNil())
			targetInstanceElement := GetTarget(reference2LiteralMapInstance, hl)
			Expect(targetInstanceElement).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceElement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetLiteral, hl)
			Expect(targetInstanceElement2).ToNot(BeNil())
			Expect(targetInstanceElement2).To(Equal(targetInstanceElement))
		})
		Specify("Reference2RefinementMap", func() {
			// Set up the abstract map
			referenceToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(referenceToRefinementMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(referenceToRefinementMap.SetLabel("ReferenceToRefinementMap", hl)).To(Succeed())
			Expect(SetSource(referenceToRefinementMap, definingSourceReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(referenceToRefinementMap, definingTargetRefinement, core.NoAttribute, hl)).To(Succeed())
			// Add the reference to the source instance
			sourceInstanceReference, err2 := uOfD.CreateReplicateAsRefinement(definingSourceReference, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			reference2RefinementMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceReference, hl)
			Expect(reference2RefinementMapInstance).ToNot(BeNil())
			targetInstanceRefinement := GetTarget(reference2RefinementMapInstance, hl)
			Expect(targetInstanceRefinement).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceRefinement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetRefinement, hl)
			Expect(targetInstanceRefinement2).ToNot(BeNil())
			Expect(targetInstanceRefinement2).To(Equal(targetInstanceRefinement))
		})
		Specify("Literal2ElementMap", func() {
			// Set up the abstract map
			literalToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(literalToElementMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(literalToElementMap.SetLabel("LiteralToElementMap", hl)).To(Succeed())
			Expect(SetSource(literalToElementMap, definingSourceLiteral, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(literalToElementMap, definingTargetElement, core.NoAttribute, hl)).To(Succeed())
			// Add the literal to the source instance
			sourceInstanceLiteral, err2 := uOfD.CreateReplicateAsRefinement(definingSourceLiteral, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceLiteral.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceLiteral.SetLabel("SourceInstanceLiteral", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			literal2ElementMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceLiteral, hl)
			Expect(literal2ElementMapInstance).ToNot(BeNil())
			targetInstanceElement := GetTarget(literal2ElementMapInstance, hl)
			Expect(targetInstanceElement).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceElement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetElement, hl)
			Expect(targetInstanceElement2).ToNot(BeNil())
			Expect(targetInstanceElement2).To(Equal(targetInstanceElement))
		})
		Specify("Literal2ReferenceMap", func() {
			// Set up the abstract map
			literalToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(literalToReferenceMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(literalToReferenceMap.SetLabel("LiteralToReferenceMap", hl)).To(Succeed())
			Expect(SetSource(literalToReferenceMap, definingSourceLiteral, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(literalToReferenceMap, definingTargetReference, core.NoAttribute, hl)).To(Succeed())
			// Add the literal to the source instance
			sourceInstanceLiteral, err2 := uOfD.CreateReplicateAsRefinement(definingSourceLiteral, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceLiteral.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceLiteral.SetLabel("SourceInstanceLiteral", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			literal2ReferenceMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceLiteral, hl)
			Expect(literal2ReferenceMapInstance).ToNot(BeNil())
			targetInstanceReference := GetTarget(literal2ReferenceMapInstance, hl)
			Expect(targetInstanceReference).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceReference2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetReference, hl)
			Expect(targetInstanceReference2).ToNot(BeNil())
			Expect(targetInstanceReference2).To(Equal(targetInstanceReference))
		})
		Specify("Literal2RefinementMap", func() {
			// Set up the abstract map
			literalToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(literalToRefinementMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(literalToRefinementMap.SetLabel("LiteralToRefinementMap", hl)).To(Succeed())
			Expect(SetSource(literalToRefinementMap, definingSourceLiteral, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(literalToRefinementMap, definingTargetRefinement, core.NoAttribute, hl)).To(Succeed())
			// Add the literal to the source instance
			sourceInstanceLiteral, err2 := uOfD.CreateReplicateAsRefinement(definingSourceLiteral, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceLiteral.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceLiteral.SetLabel("SourceInstanceLiteral", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			literal2RefinementMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceLiteral, hl)
			Expect(literal2RefinementMapInstance).ToNot(BeNil())
			targetInstanceRefinement := GetTarget(literal2RefinementMapInstance, hl)
			Expect(targetInstanceRefinement).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceRefinement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetRefinement, hl)
			Expect(targetInstanceRefinement2).ToNot(BeNil())
			Expect(targetInstanceRefinement2).To(Equal(targetInstanceRefinement))
		})
		Specify("Refinement2ElementMap", func() {
			// Set up the abstract map
			refinementToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(refinementToElementMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(refinementToElementMap.SetLabel("RefinementToElementMap", hl)).To(Succeed())
			Expect(SetSource(refinementToElementMap, definingSourceRefinement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(refinementToElementMap, definingTargetElement, core.NoAttribute, hl)).To(Succeed())
			// Add the refinement to the source instance
			sourceInstanceRefinement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceRefinement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceRefinement.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceRefinement.SetLabel("SourceInstanceRefinement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			refinement2ElementMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceRefinement, hl)
			Expect(refinement2ElementMapInstance).ToNot(BeNil())
			targetInstanceElement := GetTarget(refinement2ElementMapInstance, hl)
			Expect(targetInstanceElement).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceElement2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetElement, hl)
			Expect(targetInstanceElement2).ToNot(BeNil())
			Expect(targetInstanceElement2).To(Equal(targetInstanceElement))
		})
		Specify("Refinement2ReferenceMap", func() {
			// Set up the abstract map
			refinementToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(refinementToReferenceMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(refinementToReferenceMap.SetLabel("RefinementToReferenceMap", hl)).To(Succeed())
			Expect(SetSource(refinementToReferenceMap, definingSourceRefinement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(refinementToReferenceMap, definingTargetReference, core.NoAttribute, hl)).To(Succeed())
			// Add the refinement to the source instance
			sourceInstanceRefinement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceRefinement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceRefinement.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceRefinement.SetLabel("SourceInstanceRefinement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			refinement2ReferenceMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceRefinement, hl)
			Expect(refinement2ReferenceMapInstance).ToNot(BeNil())
			targetInstanceReference := GetTarget(refinement2ReferenceMapInstance, hl)
			Expect(targetInstanceReference).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceReference2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetReference, hl)
			Expect(targetInstanceReference2).ToNot(BeNil())
			Expect(targetInstanceReference2).To(Equal(targetInstanceReference))
		})
		Specify("Refinement2LiteralMap", func() {
			// Set up the abstract map
			refinementToLiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(refinementToLiteralMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(refinementToLiteralMap.SetLabel("RefinementToLiteralMap", hl)).To(Succeed())
			Expect(SetSource(refinementToLiteralMap, definingSourceRefinement, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(refinementToLiteralMap, definingTargetLiteral, core.NoAttribute, hl)).To(Succeed())
			// Add the refinement to the source instance
			sourceInstanceRefinement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceRefinement, hl)
			Expect(err2).To(BeNil())
			Expect(sourceInstanceRefinement.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceRefinement.SetLabel("SourceInstanceRefinement", hl)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())
			// Check the result
			refinement2LiteralMapInstance := FindMapForSource(instanceDomainMap, sourceInstanceRefinement, hl)
			Expect(refinement2LiteralMapInstance).ToNot(BeNil())
			targetInstanceLiteral := GetTarget(refinement2LiteralMapInstance, hl)
			Expect(targetInstanceLiteral).ToNot(BeNil())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, hl)
			Expect(targetInstanceDomain).ToNot(BeNil())
			targetInstanceLiteral2 := targetInstanceDomain.GetFirstOwnedConceptRefinedFrom(definingTargetLiteral, hl)
			Expect(targetInstanceLiteral2).ToNot(BeNil())
			Expect(targetInstanceLiteral2).To(Equal(targetInstanceLiteral))
		})
	})
	Describe("Individual Pointer Mapping - any to any", func() {
		var sourceAbstractPointerTarget core.Element
		var sourceAbstractReference core.Reference
		// var sourceAbstractRefinement core.Refinement
		// var sourceAbstractLiteral core.Literal
		var targetAbstractPointerTarget core.Element
		var targetAbstractReference core.Reference
		// var targetAbstractRefinement core.Refinement
		// var targetAbstractLiteral core.Literal

		var abstractPointerTarget2PointerTargetMap core.Element
		var abstractReference2ReferenceMap core.Element

		var sourceInstancePointerTarget core.Element
		var sourceInstanceReference core.Reference

		BeforeEach(func() {
			var err error
			sourceAbstractPointerTarget, err = uOfD.NewOwnedElement(definingSourceDomain, "SourceAbstractPointerTarget", hl)
			Expect(err).To(BeNil())
			sourceAbstractReference, err = uOfD.NewOwnedReference(definingSourceDomain, "SourceAbstractReference", hl)
			Expect(err).To(BeNil())
			Expect(sourceAbstractReference.SetReferencedConcept(sourceAbstractPointerTarget, core.NoAttribute, hl)).To(Succeed())
			// sourceAbstractRefinement, err = uOfD.NewOwnedRefinement(sourceAbstractDomain, "SourceAbstractRefinement", hl)
			// Expect(err).To(BeNil())
			// sourceAbstractLiteral, err = uOfD.NewOwnedLiteral(sourceAbstractDomain, "SourceAbstractLiteral", hl)
			// Expect(err).To(BeNil())

			targetAbstractPointerTarget, err = uOfD.NewOwnedElement(definingTargetDomain, "TargetAbstractPointerTarget", hl)
			Expect(err).To(BeNil())
			targetAbstractReference, err = uOfD.NewOwnedReference(definingTargetDomain, "TargetAbstractReference", hl)
			Expect(err).To(BeNil())
			Expect(targetAbstractReference.SetReferencedConcept(targetAbstractPointerTarget, core.NoAttribute, hl)).To(Succeed())
			// targetAbstractRefinement, err = uOfD.NewOwnedRefinement(targetAbstractDomain, "TargetAbstractRefinement", hl)
			// Expect(err).To(BeNil())
			// targetAbstractLiteral, err = uOfD.NewOwnedLiteral(targetAbstractDomain, "TargetAbstractLiteral", hl)
			// Expect(err).To(BeNil())

			// Abstract Map Setup
			abstractPointerTarget2PointerTargetMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(abstractPointerTarget2PointerTargetMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(abstractPointerTarget2PointerTargetMap.SetLabel("Element2ToElement2Map", hl)).To(Succeed())
			Expect(SetSource(abstractPointerTarget2PointerTargetMap, sourceAbstractPointerTarget, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(abstractPointerTarget2PointerTargetMap, targetAbstractPointerTarget, core.NoAttribute, hl)).To(Succeed())

			abstractReference2ReferenceMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(abstractReference2ReferenceMap.SetOwningConcept(definingDomainMap, hl)).To(Succeed())
			Expect(abstractReference2ReferenceMap.SetLabel("Reference2ReferenceMap", hl)).To(Succeed())
			Expect(SetSource(abstractReference2ReferenceMap, sourceAbstractReference, core.NoAttribute, hl)).To(Succeed())
			Expect(SetTarget(abstractReference2ReferenceMap, targetAbstractReference, core.NoAttribute, hl)).To(Succeed())

			// Source Instance Setup

			sourceInstancePointerTarget, err = uOfD.CreateReplicateAsRefinement(sourceAbstractPointerTarget, hl)
			Expect(err).To(BeNil())
			Expect(sourceInstancePointerTarget.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstancePointerTarget.SetLabel("SourceInstancePointerTarget", hl)).To(Succeed())

			sourceInstanceReference, err = uOfD.CreateReplicateReferenceAsRefinement(sourceAbstractReference, hl)
			Expect(err).To(BeNil())
			Expect(sourceInstanceReference.SetOwningConcept(instanceSourceDomain, hl)).To(Succeed())
			Expect(sourceInstanceReference.SetLabel("SourceInstanceReference", hl)).To(Succeed())

		})
		Specify("Reference Pointer to Reference Pointer", func() {
			// Add the pointer map
			pointer2ReferencePointerMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, hl)
			Expect(err).To(BeNil())
			Expect(pointer2ReferencePointerMap.SetOwningConcept(abstractReference2ReferenceMap, hl)).To(Succeed())
			Expect(pointer2ReferencePointerMap.SetLabel("Pointer2ReferencePointerMap", hl)).To(Succeed())
			Expect(SetSource(pointer2ReferencePointerMap, sourceAbstractReference, core.ReferencedConceptID, hl)).To(Succeed())
			Expect(SetTarget(pointer2ReferencePointerMap, targetAbstractReference, core.ReferencedConceptID, hl)).To(Succeed())

			// Trigger the map
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, hl)).To(Succeed())

			// Diagnostics view
			graph := core.NewCrlGraph("Pointer2RferencePointerMapTest")
			Expect(graph.AddConceptRecursively(definingSourceFolder, hl)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingTargetFolder, hl)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingMapFolder, hl)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceSourceFolder, hl)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceMapFolder, hl)).To(Succeed())
			Expect(graph.ExportDOT(tempDirPath, "Pointer2RferencePointerMapTest")).To(Succeed())

			// Check the result
			sourceInstanceReferenceMap := FindMapForSource(instanceDomainMap, sourceInstanceReference, hl)
			Expect(sourceInstanceReferenceMap).ToNot(BeNil())
			sourceInstanceAttributeReferenceMap := FindAttributeMapForSource(instanceDomainMap, sourceInstanceReference, core.ReferencedConceptID, hl)
			Expect(sourceInstanceAttributeReferenceMap).ToNot(BeNil())
			Expect(sourceInstanceAttributeReferenceMap.GetOwningConcept(hl).GetConceptID(hl)).To(Equal(sourceInstanceReferenceMap.GetConceptID(hl)))
			targetInstance := GetTarget(sourceInstanceReferenceMap, hl)
			Expect(targetInstance).ToNot(BeNil())
			targetInstanceAttribute := GetTarget(sourceInstanceAttributeReferenceMap, hl)
			Expect(targetInstanceAttribute).To(Equal(targetInstance))

			targetInstanceElement1Map := FindMapForSource(instanceDomainMap, sourceInstancePointerTarget, hl)
			Expect(targetInstanceElement1Map).ToNot(BeNil())
			targetInstanceElement1 := GetTarget(targetInstanceElement1Map, hl)
			Expect(targetInstanceElement1).ToNot(BeNil())

			switch castInstance := targetInstance.(type) {
			case core.Reference:
				Expect(castInstance.GetReferencedConcept(hl).GetConceptID(hl)).To(Equal(targetInstanceElement1.GetConceptID(hl)))
			}
		})
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
