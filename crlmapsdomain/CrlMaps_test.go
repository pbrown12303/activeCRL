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
		Expect(BuildCrlMapsDomain(uOfD1, hl1)).To(Succeed())
		md1 := uOfD1.GetElementWithURI(CrlMapsDomainURI)
		Expect(md1).ToNot(BeNil())
		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewTransaction()
		Expect(BuildCrlMapsDomain(uOfD2, hl2)).To(Succeed())
		md2 := uOfD2.GetElementWithURI(CrlMapsDomainURI)
		Expect(md2).ToNot(BeNil())
		Expect(core.RecursivelyEquivalent(md1, hl1, md2, hl2, true)).To(BeTrue())
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
	var trans *core.Transaction
	var definingSourceFolder core.Concept
	var definingSourceDomain core.Concept
	var definingTargetFolder core.Concept
	var definingTargetDomain core.Concept
	var definingMapFolder core.Concept
	var definingDomainMap core.Concept
	var instanceSourceFolder core.Concept
	var instanceSourceDomain core.Concept
	var instanceMapFolder core.Concept
	var instanceDomainMap core.Concept
	var tempDirPath string

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
		BuildCrlMapsDomain(uOfD, trans)
		var err error

		// Defining Source
		definingSourceFolder, err = uOfD.NewOwnedElement(nil, "DefiningSourceFolder", trans)
		Expect(err).To(BeNil())
		definingSourceDomain, err = uOfD.NewOwnedElement(definingSourceFolder, "DefiningSourceDomain", trans)
		Expect(err).To(BeNil())

		// Defining Target
		definingTargetFolder, err = uOfD.NewOwnedElement(nil, "DefiningTargetFolder", trans)
		Expect(err).To(BeNil())
		definingTargetDomain, err = uOfD.NewOwnedElement(definingTargetFolder, "DefiningTargetDomain", trans)
		Expect(err).To(BeNil())

		// Defining Map
		definingMapFolder, err = uOfD.NewOwnedElement(nil, "DefiningMapFolder", trans)
		Expect(err).To(BeNil())
		definingDomainMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
		Expect(err).To(BeNil())
		Expect(definingDomainMap.SetLabel("DefiningDomainMap", trans)).To(Succeed())
		Expect(definingDomainMap.SetOwningConcept(definingMapFolder, trans)).To(Succeed())
		// mapAbstractDomainOwnedConcepts := mapAbstractDomain.GetOwnedConcepts(trans)
		// fmt.Print(mapAbstractDomainOwnedConcepts)
		definingDomainMapSourceRef := definingDomainMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
		Expect(definingDomainMapSourceRef.SetReferencedConcept(definingSourceDomain, core.NoAttribute, trans)).To(Succeed())
		definingDomainMapTargetRef := definingDomainMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
		Expect(definingDomainMapTargetRef.SetReferencedConcept(definingTargetDomain, core.NoAttribute, trans)).To(Succeed())

		// Source Instance
		instanceSourceFolder, err = uOfD.NewOwnedElement(nil, "instanceSourceFolder", trans)
		Expect(err).To(BeNil())
		instanceSourceDomain, err = uOfD.CreateReplicateAsRefinement(definingSourceDomain, trans)
		Expect(err).To(BeNil())
		Expect(instanceSourceDomain.SetLabel("InstanceSourceDomain", trans)).To(Succeed())
		Expect(instanceSourceDomain.SetOwningConcept(instanceSourceFolder, trans)).To(Succeed())

		// Map Instance
		instanceMapFolder, err = uOfD.NewOwnedElement(nil, "instanceMapFolder", trans)
		Expect(err).To(BeNil())
		instanceDomainMap, err = uOfD.CreateReplicateAsRefinement(definingDomainMap, trans)
		Expect(err).To(BeNil())
		Expect(instanceDomainMap.SetLabel("InstanceDomainMap", trans)).To(Succeed())
		Expect(instanceDomainMap.SetOwningConcept(instanceMapFolder, trans)).To(Succeed())

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
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			targetInstanceDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(targetInstanceDomain).ToNot(BeNil())
		})
	})
	Describe("Individual Concept Mapping - any to any", func() {
		var definingSourceElement core.Concept
		var definingSourceReference core.Concept
		var definingSourceRefinement core.Concept
		var definingSourceLiteral core.Concept
		var definingTargetElement core.Concept
		var definingTargetReference core.Concept
		var definingTargetRefinement core.Concept
		var definingTargetLiteral core.Concept
		BeforeEach(func() {
			var err error
			definingSourceElement, err = uOfD.NewOwnedElement(definingSourceDomain, "DefiningSourceElement", trans)
			Expect(err).To(BeNil())
			definingSourceReference, err = uOfD.NewOwnedReference(definingSourceDomain, "DefiningSourceReference", trans)
			Expect(err).To(BeNil())
			definingSourceRefinement, err = uOfD.NewOwnedRefinement(definingSourceDomain, "DefiningSourceRefinement", trans)
			Expect(err).To(BeNil())
			definingSourceLiteral, err = uOfD.NewOwnedLiteral(definingSourceDomain, "DefiningSourceLiteral", trans)
			Expect(err).To(BeNil())

			definingTargetElement, err = uOfD.NewOwnedElement(definingTargetDomain, "DefiningTargetElement", trans)
			Expect(err).To(BeNil())
			definingTargetReference, err = uOfD.NewOwnedReference(definingTargetDomain, "DefiningTargetReference", trans)
			Expect(err).To(BeNil())
			definingTargetRefinement, err = uOfD.NewOwnedRefinement(definingTargetDomain, "DefiningTargetRefinement", trans)
			Expect(err).To(BeNil())
			definingTargetLiteral, err = uOfD.NewOwnedLiteral(definingTargetDomain, "DefiningTargetLiteral", trans)
			Expect(err).To(BeNil())
		})
		Specify("Element to Element Map", func() {
			// Set up the abstract map
			definingElementToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingElementToElementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingElementToElementMap.SetLabel("DefiningElementToElementMap", trans)).To(Succeed())
			Expect(SetSource(definingElementToElementMap, definingSourceElement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingElementToElementMap, definingTargetElement, core.NoAttribute, trans)).To(Succeed())
			// Add the element to the source instance
			instanceSourceElement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceElement, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceElement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceElement.SetLabel("InstanceSourceElement", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceElementToElementMap := FindMapForSource(instanceDomainMap, instanceSourceElement, trans)
			Expect(instanceElementToElementMap).ToNot(BeNil())
			instanceTargetElement := GetTarget(instanceElementToElementMap, trans)
			Expect(instanceTargetElement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetElement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetElement, trans)
			Expect(instanceTargetElement2).ToNot(BeNil())
			Expect(instanceTargetElement2).To(Equal(instanceTargetElement))
		})
		Specify("Reference to Reference Map", func() {
			// Set up the abstract map
			definingReference2ReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingReference2ReferenceMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingReference2ReferenceMap.SetLabel("DefiningReference2ReferenceMap", trans)).To(Succeed())
			Expect(SetSource(definingReference2ReferenceMap, definingSourceReference, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingReference2ReferenceMap, definingTargetReference, core.NoAttribute, trans)).To(Succeed())
			// Add the reference to the source instance
			instanceSourceReference, err2 := uOfD.CreateReplicateAsRefinement(definingSourceReference, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceReference.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceReference.SetLabel("InstanceSourceReference", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())

			// Diagnostics view
			graph := core.NewCrlGraph("ReferencetoReferenceMapTest")
			// Expect(graph.AddConceptRecursively(definingSourceFolder, trans)).To(Succeed())
			// Expect(graph.AddConceptRecursively(definingTargetFolder, trans)).To(Succeed())
			// Expect(graph.AddConceptRecursively(definingMapFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceMapFolder, trans)).To(Succeed())
			Expect(graph.ExportDOT(tempDirPath, "ReferencetoReferenceMapTest")).To(Succeed())

			// Check the result
			instanceReferenceMap := FindMapForSource(instanceDomainMap, instanceSourceReference, trans)
			Expect(instanceReferenceMap).ToNot(BeNil())
			instanceTargetReference := GetTarget(instanceReferenceMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetReference2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetReference, trans)
			Expect(instanceTargetReference2).ToNot(BeNil())
			Expect(instanceTargetReference2).To(Equal(instanceTargetReference))
		})
		Specify("Literal to Literal Map", func() {
			// Set up the abstract map
			definingLiteral2LiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingLiteral2LiteralMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingLiteral2LiteralMap.SetLabel("DefiningLiteral2LiteralMap", trans)).To(Succeed())
			Expect(SetSource(definingLiteral2LiteralMap, definingSourceLiteral, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingLiteral2LiteralMap, definingTargetLiteral, core.NoAttribute, trans)).To(Succeed())
			// Add the literal to the source instance
			instanceSourceLiteral, err2 := uOfD.CreateReplicateAsRefinement(definingSourceLiteral, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceLiteral.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceLiteral.SetLabel("InstanceSourceLiteral", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceLiteralMap := FindMapForSource(instanceDomainMap, instanceSourceLiteral, trans)
			Expect(instanceLiteralMap).ToNot(BeNil())
			instanceTargetLiteral := GetTarget(instanceLiteralMap, trans)
			Expect(instanceTargetLiteral).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetLiteral2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetLiteral, trans)
			Expect(instanceTargetLiteral2).ToNot(BeNil())
			Expect(instanceTargetLiteral2).To(Equal(instanceTargetLiteral))
		})
		Specify("Refinement to Refinement Map", func() {
			// Set up the abstract map
			definingRefinement2RefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingRefinement2RefinementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingRefinement2RefinementMap.SetLabel("DefiningRefinement2RefinementMap", trans)).To(Succeed())
			Expect(SetSource(definingRefinement2RefinementMap, definingSourceRefinement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingRefinement2RefinementMap, definingTargetRefinement, core.NoAttribute, trans)).To(Succeed())
			// Add the refinement to the source instance
			instanceSourceRefinement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceRefinement, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceRefinement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceRefinement.SetLabel("InstanceSourceRefinement", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceRefinementMap := FindMapForSource(instanceDomainMap, instanceSourceRefinement, trans)
			Expect(instanceRefinementMap).ToNot(BeNil())
			instanceTargetRefinement := GetTarget(instanceRefinementMap, trans)
			Expect(instanceTargetRefinement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetRefinement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetRefinement, trans)
			Expect(instanceTargetRefinement2).ToNot(BeNil())
			Expect(instanceTargetRefinement2).To(Equal(instanceTargetRefinement))
		})
		Specify("Element2ReferenceMap", func() {
			// Set up the abstract map
			definingElementToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingElementToReferenceMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingElementToReferenceMap.SetLabel("DefiningElementToReferenceMap", trans)).To(Succeed())
			Expect(SetSource(definingElementToReferenceMap, definingSourceElement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingElementToReferenceMap, definingTargetReference, core.NoAttribute, trans)).To(Succeed())
			// Add the element to the source instance
			instanceSourceElement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceElement, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceElement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceElement.SetLabel("InstanceSourceElement", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceElement2ReferenceMap := FindMapForSource(instanceDomainMap, instanceSourceElement, trans)
			Expect(instanceElement2ReferenceMap).ToNot(BeNil())
			instanceTargetReference := GetTarget(instanceElement2ReferenceMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetReference2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetReference, trans)
			Expect(instanceTargetReference2).ToNot(BeNil())
			Expect(instanceTargetReference2).To(Equal(instanceTargetReference))
		})
		Specify("Element2LiteralMap", func() {
			// Set up the abstract map
			definingElementToLiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingElementToLiteralMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingElementToLiteralMap.SetLabel("DefiningElementToLiteralMap", trans)).To(Succeed())
			Expect(SetSource(definingElementToLiteralMap, definingSourceElement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingElementToLiteralMap, definingTargetLiteral, core.NoAttribute, trans)).To(Succeed())
			// Add the element to the source instance
			instanceSourceElement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceElement, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceElement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceElement.SetLabel("InstanceSourceElement", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceElement2LiteralMap := FindMapForSource(instanceDomainMap, instanceSourceElement, trans)
			Expect(instanceElement2LiteralMap).ToNot(BeNil())
			instanceTargetLiteral := GetTarget(instanceElement2LiteralMap, trans)
			Expect(instanceTargetLiteral).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetLiteral2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetLiteral, trans)
			Expect(instanceTargetLiteral2).ToNot(BeNil())
			Expect(instanceTargetLiteral2).To(Equal(instanceTargetLiteral))
		})
		Specify("Element2RefinementMap", func() {
			// Set up the abstract map
			definingElementToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingElementToRefinementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingElementToRefinementMap.SetLabel("DefiningElementToRefinementMap", trans)).To(Succeed())
			Expect(SetSource(definingElementToRefinementMap, definingSourceElement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingElementToRefinementMap, definingTargetRefinement, core.NoAttribute, trans)).To(Succeed())
			// Add the element to the source instance
			instanceSourceElement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceElement, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceElement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceElement.SetLabel("InstanceSourceElement", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceElement2RefinementMap := FindMapForSource(instanceDomainMap, instanceSourceElement, trans)
			Expect(instanceElement2RefinementMap).ToNot(BeNil())
			instanceTargetRefinement := GetTarget(instanceElement2RefinementMap, trans)
			Expect(instanceTargetRefinement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetRefinement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetRefinement, trans)
			Expect(instanceTargetRefinement2).ToNot(BeNil())
			Expect(instanceTargetRefinement2).To(Equal(instanceTargetRefinement))
		})
		Specify("Reference2ElementMap", func() {
			// Set up the abstract map
			definingReferenceToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingReferenceToElementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingReferenceToElementMap.SetLabel("DefiningReferenceToElementMap", trans)).To(Succeed())
			Expect(SetSource(definingReferenceToElementMap, definingSourceReference, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingReferenceToElementMap, definingTargetElement, core.NoAttribute, trans)).To(Succeed())
			// Add the reference to the source instance
			instanceSourceReference, err2 := uOfD.CreateReplicateAsRefinement(definingSourceReference, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceReference.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceReference.SetLabel("InstanceSourceReference", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceReference2ElementMap := FindMapForSource(instanceDomainMap, instanceSourceReference, trans)
			Expect(instanceReference2ElementMap).ToNot(BeNil())
			instanceTargetElement := GetTarget(instanceReference2ElementMap, trans)
			Expect(instanceTargetElement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetElement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetElement, trans)
			Expect(instanceTargetElement2).ToNot(BeNil())
			Expect(instanceTargetElement2).To(Equal(instanceTargetElement))
		})
		Specify("Reference2LiteralMap", func() {
			// Set up the abstract map
			definingReference2LiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingReference2LiteralMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingReference2LiteralMap.SetLabel("DefiningReference2LiteralMap", trans)).To(Succeed())
			Expect(SetSource(definingReference2LiteralMap, definingSourceReference, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingReference2LiteralMap, definingTargetLiteral, core.NoAttribute, trans)).To(Succeed())
			// Add the reference to the source instance
			instanceSourceReference, err2 := uOfD.CreateReplicateAsRefinement(definingSourceReference, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceReference.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceReference.SetLabel("InstanceSourceReference", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceReference2LiteralMap := FindMapForSource(instanceDomainMap, instanceSourceReference, trans)
			Expect(instanceReference2LiteralMap).ToNot(BeNil())
			instanceTargetElement := GetTarget(instanceReference2LiteralMap, trans)
			Expect(instanceTargetElement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetElement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetLiteral, trans)
			Expect(instanceTargetElement2).ToNot(BeNil())
			Expect(instanceTargetElement2).To(Equal(instanceTargetElement))
		})
		Specify("Reference2RefinementMap", func() {
			// Set up the abstract map
			definingReferenceToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingReferenceToRefinementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingReferenceToRefinementMap.SetLabel("DefiningReferenceToRefinementMap", trans)).To(Succeed())
			Expect(SetSource(definingReferenceToRefinementMap, definingSourceReference, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingReferenceToRefinementMap, definingTargetRefinement, core.NoAttribute, trans)).To(Succeed())
			// Add the reference to the source instance
			instanceSourceReference, err2 := uOfD.CreateReplicateAsRefinement(definingSourceReference, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceReference.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceReference.SetLabel("InstanceSourceReference", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceReference2RefinementMap := FindMapForSource(instanceDomainMap, instanceSourceReference, trans)
			Expect(instanceReference2RefinementMap).ToNot(BeNil())
			instanceTargetRefinement := GetTarget(instanceReference2RefinementMap, trans)
			Expect(instanceTargetRefinement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetRefinement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetRefinement, trans)
			Expect(instanceTargetRefinement2).ToNot(BeNil())
			Expect(instanceTargetRefinement2).To(Equal(instanceTargetRefinement))
		})
		Specify("Literal2ElementMap", func() {
			// Set up the abstract map
			definingLiteralToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingLiteralToElementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingLiteralToElementMap.SetLabel("DefiningLiteralToElementMap", trans)).To(Succeed())
			Expect(SetSource(definingLiteralToElementMap, definingSourceLiteral, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingLiteralToElementMap, definingTargetElement, core.NoAttribute, trans)).To(Succeed())
			// Add the literal to the source instance
			instanceSourceLiteral, err2 := uOfD.CreateReplicateAsRefinement(definingSourceLiteral, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceLiteral.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceLiteral.SetLabel("InstanceSourceLiteral", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceLiteral2ElementMap := FindMapForSource(instanceDomainMap, instanceSourceLiteral, trans)
			Expect(instanceLiteral2ElementMap).ToNot(BeNil())
			instanceTargetElement := GetTarget(instanceLiteral2ElementMap, trans)
			Expect(instanceTargetElement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetElement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetElement, trans)
			Expect(instanceTargetElement2).ToNot(BeNil())
			Expect(instanceTargetElement2).To(Equal(instanceTargetElement))
		})
		Specify("Literal2ReferenceMap", func() {
			// Set up the abstract map
			definingLiteralToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingLiteralToReferenceMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingLiteralToReferenceMap.SetLabel("DefiningLiteralToReferenceMap", trans)).To(Succeed())
			Expect(SetSource(definingLiteralToReferenceMap, definingSourceLiteral, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingLiteralToReferenceMap, definingTargetReference, core.NoAttribute, trans)).To(Succeed())
			// Add the literal to the source instance
			instanceSourceLiteral, err2 := uOfD.CreateReplicateAsRefinement(definingSourceLiteral, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceLiteral.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceLiteral.SetLabel("InstanceSourceLiteral", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceLiteral2ReferenceMap := FindMapForSource(instanceDomainMap, instanceSourceLiteral, trans)
			Expect(instanceLiteral2ReferenceMap).ToNot(BeNil())
			instanceTargetReference := GetTarget(instanceLiteral2ReferenceMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetReference2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetReference, trans)
			Expect(instanceTargetReference2).ToNot(BeNil())
			Expect(instanceTargetReference2).To(Equal(instanceTargetReference))
		})
		Specify("Literal2RefinementMap", func() {
			// Set up the abstract map
			definingLiteralToRefinementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingLiteralToRefinementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingLiteralToRefinementMap.SetLabel("DefiningLiteralToRefinementMap", trans)).To(Succeed())
			Expect(SetSource(definingLiteralToRefinementMap, definingSourceLiteral, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingLiteralToRefinementMap, definingTargetRefinement, core.NoAttribute, trans)).To(Succeed())
			// Add the literal to the source instance
			instanceSourceLiteral, err2 := uOfD.CreateReplicateAsRefinement(definingSourceLiteral, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceLiteral.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceLiteral.SetLabel("InstanceSourceLiteral", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceLiteral2RefinementMap := FindMapForSource(instanceDomainMap, instanceSourceLiteral, trans)
			Expect(instanceLiteral2RefinementMap).ToNot(BeNil())
			instanceTargetRefinement := GetTarget(instanceLiteral2RefinementMap, trans)
			Expect(instanceTargetRefinement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetRefinement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetRefinement, trans)
			Expect(instanceTargetRefinement2).ToNot(BeNil())
			Expect(instanceTargetRefinement2).To(Equal(instanceTargetRefinement))
		})
		Specify("Refinement2ElementMap", func() {
			// Set up the abstract map
			definingRefinementToElementMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingRefinementToElementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingRefinementToElementMap.SetLabel("DefiningRefinementToElementMap", trans)).To(Succeed())
			Expect(SetSource(definingRefinementToElementMap, definingSourceRefinement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingRefinementToElementMap, definingTargetElement, core.NoAttribute, trans)).To(Succeed())
			// Add the refinement to the source instance
			instanceSourceRefinement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceRefinement, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceRefinement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceRefinement.SetLabel("InstanceSourceRefinement", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceRefinement2ElementMap := FindMapForSource(instanceDomainMap, instanceSourceRefinement, trans)
			Expect(instanceRefinement2ElementMap).ToNot(BeNil())
			instanceTargetElement := GetTarget(instanceRefinement2ElementMap, trans)
			Expect(instanceTargetElement).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetElement2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetElement, trans)
			Expect(instanceTargetElement2).ToNot(BeNil())
			Expect(instanceTargetElement2).To(Equal(instanceTargetElement))
		})
		Specify("Refinement2ReferenceMap", func() {
			// Set up the abstract map
			definingRefinementToReferenceMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingRefinementToReferenceMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingRefinementToReferenceMap.SetLabel("DefiningRefinementToReferenceMap", trans)).To(Succeed())
			Expect(SetSource(definingRefinementToReferenceMap, definingSourceRefinement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingRefinementToReferenceMap, definingTargetReference, core.NoAttribute, trans)).To(Succeed())
			// Add the refinement to the source instance
			instaneSourceRefinement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceRefinement, trans)
			Expect(err2).To(BeNil())
			Expect(instaneSourceRefinement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instaneSourceRefinement.SetLabel("InstanceSourceRefinement", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceRefinement2ReferenceMap := FindMapForSource(instanceDomainMap, instaneSourceRefinement, trans)
			Expect(instanceRefinement2ReferenceMap).ToNot(BeNil())
			instanceTargetReference := GetTarget(instanceRefinement2ReferenceMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetReference2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetReference, trans)
			Expect(instanceTargetReference2).ToNot(BeNil())
			Expect(instanceTargetReference2).To(Equal(instanceTargetReference))
		})
		Specify("Refinement2LiteralMap", func() {
			// Set up the abstract map
			definingRefinementToLiteralMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingRefinementToLiteralMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingRefinementToLiteralMap.SetLabel("DefiningRefinementToLiteralMap", trans)).To(Succeed())
			Expect(SetSource(definingRefinementToLiteralMap, definingSourceRefinement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingRefinementToLiteralMap, definingTargetLiteral, core.NoAttribute, trans)).To(Succeed())
			// Add the refinement to the source instance
			instanceSourceRefinement, err2 := uOfD.CreateReplicateAsRefinement(definingSourceRefinement, trans)
			Expect(err2).To(BeNil())
			Expect(instanceSourceRefinement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceRefinement.SetLabel("InstanceSourceRefinement", trans)).To(Succeed())
			// Trigger the mapping
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())
			// Check the result
			instanceRefinement2LiteralMap := FindMapForSource(instanceDomainMap, instanceSourceRefinement, trans)
			Expect(instanceRefinement2LiteralMap).ToNot(BeNil())
			instanceTargetLiteral := GetTarget(instanceRefinement2LiteralMap, trans)
			Expect(instanceTargetLiteral).ToNot(BeNil())
			instanceTargetDomain := instanceMapFolder.GetFirstOwnedConceptRefinedFrom(definingTargetDomain, trans)
			Expect(instanceTargetDomain).ToNot(BeNil())
			instanceTargetLiteral2 := instanceTargetDomain.GetFirstOwnedConceptRefinedFrom(definingTargetLiteral, trans)
			Expect(instanceTargetLiteral2).ToNot(BeNil())
			Expect(instanceTargetLiteral2).To(Equal(instanceTargetLiteral))
		})
	})
	Describe("Individual Pointer Mapping - any to any", func() {
		var definingSourceReferent core.Concept  // The source element being referenced
		var definingSourceReference core.Concept // The reference to the source element
		var definingSourceRefinement core.Concept
		var definingSourceLiteral core.Concept
		var definingTargetReferent core.Concept
		var definingTargetReference core.Concept
		var definingTargetRefinement core.Concept
		var definingTargetLiteral core.Concept

		var definingReferent2ReferentMap core.Concept
		var definingReference2ReferenceMap core.Concept
		var definingRefinement2RefinementMap core.Concept
		var definingLiteral2LiteralMap core.Concept

		var instanceSourceReferent core.Concept
		var instanceSourceReference core.Concept
		var instanceSourceRefinement core.Concept
		var instanceSourceLiteral core.Concept

		BeforeEach(func() {
			var err error
			definingSourceReferent, err = uOfD.NewOwnedElement(definingSourceDomain, "DefiningSourceReferent", trans)
			Expect(err).To(BeNil())
			definingSourceReference, err = uOfD.NewOwnedReference(definingSourceDomain, "DefiningSourceReference", trans)
			Expect(err).To(BeNil())
			definingSourceRefinement, err = uOfD.NewOwnedRefinement(definingSourceDomain, "DefiningSourceRefinement", trans)
			Expect(err).To(BeNil())
			definingSourceLiteral, err = uOfD.NewOwnedLiteral(definingSourceDomain, "DefiningSourceLiteral", trans)
			Expect(err).To(BeNil())

			definingTargetReferent, err = uOfD.NewOwnedElement(definingTargetDomain, "DefiningTargetReferent", trans)
			Expect(err).To(BeNil())
			definingTargetReference, err = uOfD.NewOwnedReference(definingTargetDomain, "DefiningTargetReference", trans)
			Expect(err).To(BeNil())
			definingTargetRefinement, err = uOfD.NewOwnedRefinement(definingTargetDomain, "DefiningTargetRefinement", trans)
			Expect(err).To(BeNil())
			definingTargetLiteral, err = uOfD.NewOwnedLiteral(definingTargetDomain, "DefiningTargetLiteral", trans)
			Expect(err).To(BeNil())

			// Defining Map Setup
			definingReferent2ReferentMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingReferent2ReferentMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingReferent2ReferentMap.SetLabel("DefiningReferent2ReferentMap", trans)).To(Succeed())
			Expect(SetSource(definingReferent2ReferentMap, definingSourceReferent, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingReferent2ReferentMap, definingTargetReferent, core.NoAttribute, trans)).To(Succeed())

			definingReference2ReferenceMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingReference2ReferenceMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingReference2ReferenceMap.SetLabel("DefiningReference2ReferenceMap", trans)).To(Succeed())
			Expect(SetSource(definingReference2ReferenceMap, definingSourceReference, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingReference2ReferenceMap, definingTargetReference, core.NoAttribute, trans)).To(Succeed())

			definingRefinement2RefinementMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingRefinement2RefinementMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingRefinement2RefinementMap.SetLabel("DefiningRefinement2RefinementMap", trans)).To(Succeed())
			Expect(SetSource(definingRefinement2RefinementMap, definingSourceRefinement, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingRefinement2RefinementMap, definingTargetRefinement, core.NoAttribute, trans)).To(Succeed())

			definingLiteral2LiteralMap, err = uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingLiteral2LiteralMap.SetOwningConcept(definingDomainMap, trans)).To(Succeed())
			Expect(definingLiteral2LiteralMap.SetLabel("DefiningLiteral2LiteralMap", trans)).To(Succeed())
			Expect(SetSource(definingLiteral2LiteralMap, definingSourceLiteral, core.NoAttribute, trans)).To(Succeed())
			Expect(SetTarget(definingLiteral2LiteralMap, definingTargetLiteral, core.NoAttribute, trans)).To(Succeed())

			// Source Instance Setup

			instanceSourceReferent, err = uOfD.CreateReplicateAsRefinement(definingSourceReferent, trans)
			Expect(err).To(BeNil())
			Expect(instanceSourceReferent.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceReferent.SetLabel("InstanceSourceReferent", trans)).To(Succeed())

			instanceSourceReference, err = uOfD.CreateReplicateReferenceAsRefinement(definingSourceReference, trans)
			Expect(err).To(BeNil())
			Expect(instanceSourceReference.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceReference.SetLabel("InstanceSourceReference", trans)).To(Succeed())

			instanceSourceRefinement, err = uOfD.CreateReplicateRefinementAsRefinement(definingSourceRefinement, trans)
			Expect(err).To(BeNil())
			Expect(instanceSourceRefinement.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceRefinement.SetLabel("InstanceSourceRefinement", trans)).To(Succeed())

			instanceSourceLiteral, err = uOfD.CreateReplicateLiteralAsRefinement(definingSourceLiteral, trans)
			Expect(err).To(BeNil())
			Expect(instanceSourceLiteral.SetOwningConcept(instanceSourceDomain, trans)).To(Succeed())
			Expect(instanceSourceLiteral.SetLabel("InstanceSourceLiteral", trans)).To(Succeed())

		})
		Specify("Referenced Element Pointer to Referenced Element Pointer", func() {
			// Establish the value that is going to be mapped
			Expect(instanceSourceReference.SetReferencedConcept(instanceSourceReferent, core.NoAttribute, trans)).To(Succeed())
			// Set up the pointers to be mapped
			Expect(definingSourceReference.SetReferencedConcept(definingSourceReferent, core.NoAttribute, trans)).To(Succeed())
			Expect(definingTargetReference.SetReferencedConcept(definingTargetReferent, core.NoAttribute, trans)).To(Succeed())
			// Add the pointer map
			definingReferencedElementPointer2ReferencedElementPointerMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingReferencedElementPointer2ReferencedElementPointerMap.SetOwningConcept(definingReference2ReferenceMap, trans)).To(Succeed())
			Expect(definingReferencedElementPointer2ReferencedElementPointerMap.SetLabel("DefiningReferencedElementPointer2ReferencedElementPointerMap", trans)).To(Succeed())
			Expect(SetSource(definingReferencedElementPointer2ReferencedElementPointerMap, definingSourceReference, core.ReferencedConceptID, trans)).To(Succeed())
			Expect(SetTarget(definingReferencedElementPointer2ReferencedElementPointerMap, definingTargetReference, core.ReferencedConceptID, trans)).To(Succeed())

			// Trigger the map
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())

			// Diagnostics view
			graph := core.NewCrlGraph("ReferencedElementPointer2ReferencePointerMapTest")
			Expect(graph.AddConceptRecursively(definingSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingTargetFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingMapFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceMapFolder, trans)).To(Succeed())
			Expect(graph.ExportDOT(tempDirPath, "ReferencedElementPointer2ReferencedElementPointerMapTest")).To(Succeed())

			// Check the result
			instanceReference2ReferenceMap := FindMapForSource(instanceDomainMap, instanceSourceReference, trans)
			Expect(instanceReference2ReferenceMap).ToNot(BeNil())
			instanceReferencedElementPointer2ReferencedElementPointerMap := FindAttributeMapForSource(instanceDomainMap, instanceSourceReference, core.ReferencedConceptID, trans)
			Expect(instanceReferencedElementPointer2ReferencedElementPointerMap).ToNot(BeNil())
			Expect(instanceReferencedElementPointer2ReferencedElementPointerMap.GetOwningConcept(trans).GetConceptID(trans)).To(Equal(instanceReference2ReferenceMap.GetConceptID(trans)))
			instanceTargetReference := GetTarget(instanceReference2ReferenceMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceAttributeTarget := GetTarget(instanceReferencedElementPointer2ReferencedElementPointerMap, trans)
			Expect(instanceAttributeTarget).To(Equal(instanceTargetReference))

			instanceReferent2ReferentMap := FindMapForSource(instanceDomainMap, instanceSourceReferent, trans)
			Expect(instanceReferent2ReferentMap).ToNot(BeNil())
			instanceTargetReferent2 := GetTarget(instanceReferent2ReferentMap, trans)
			Expect(instanceTargetReferent2).ToNot(BeNil())

			switch instanceTargetReference.GetConceptType() {
			case core.Reference:
				Expect(instanceTargetReference.GetReferencedConcept(trans).GetConceptID(trans)).To(Equal(instanceTargetReferent2.GetConceptID(trans)))
			}
		})
		Specify("Reference Owner Pointer to Reference Owner Pointer", func() {
			// Establish the value that is going to be mapped
			Expect(instanceSourceReference.SetOwningConcept(instanceSourceReferent, trans)).To(Succeed())
			// Set up the pointers to be mapped
			Expect(definingSourceReference.SetOwningConcept(definingSourceReferent, trans)).To(Succeed())
			Expect(definingTargetReference.SetOwningConcept(definingTargetReferent, trans)).To(Succeed())
			// Add the pointer map
			definingOwnerPointer2OwnerPointerMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingOwnerPointer2OwnerPointerMap.SetOwningConcept(definingReference2ReferenceMap, trans)).To(Succeed())
			Expect(definingOwnerPointer2OwnerPointerMap.SetLabel("definingOwnerPointer2OwnerPointerMap", trans)).To(Succeed())
			Expect(SetSource(definingOwnerPointer2OwnerPointerMap, definingSourceReference, core.OwningConceptID, trans)).To(Succeed())
			Expect(SetTarget(definingOwnerPointer2OwnerPointerMap, definingTargetReference, core.OwningConceptID, trans)).To(Succeed())

			// Trigger the map
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())

			// Diagnostics view
			graph := core.NewCrlGraph("OwnerPointer2OwnerPointerMapTest")
			Expect(graph.AddConceptRecursively(definingSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingTargetFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingMapFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceMapFolder, trans)).To(Succeed())
			Expect(graph.ExportDOT(tempDirPath, "OwnerPointer2OwnerPointerMapTest")).To(Succeed())

			// Check the result
			instanceReference2ReferenceMap := FindMapForSource(instanceDomainMap, instanceSourceReference, trans)
			Expect(instanceReference2ReferenceMap).ToNot(BeNil())
			instanceOwnerPointer2OwnerPointerMap := FindAttributeMapForSource(instanceDomainMap, instanceSourceReference, core.OwningConceptID, trans)
			Expect(instanceOwnerPointer2OwnerPointerMap).ToNot(BeNil())
			Expect(instanceOwnerPointer2OwnerPointerMap.GetOwningConcept(trans).GetConceptID(trans)).To(Equal(instanceReference2ReferenceMap.GetConceptID(trans)))
			instanceTargetReference := GetTarget(instanceReference2ReferenceMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceAttributeTarget := GetTarget(instanceOwnerPointer2OwnerPointerMap, trans)
			Expect(instanceAttributeTarget).To(Equal(instanceTargetReference))

			instanceReferent2ReferentMap := FindMapForSource(instanceDomainMap, instanceSourceReferent, trans)
			Expect(instanceReferent2ReferentMap).ToNot(BeNil())
			instanceTargetReferent2 := GetTarget(instanceReferent2ReferentMap, trans)
			Expect(instanceTargetReferent2).ToNot(BeNil())

			Expect(instanceTargetReference.GetOwningConceptID(trans)).To(Equal(instanceTargetReferent2.GetConceptID(trans)))
		})

		Specify("Refinement Abstract Pointer to Refinement Abstract Pointer", func() {
			// Establish the value that is going to be mapped
			Expect(instanceSourceRefinement.SetAbstractConcept(instanceSourceReferent, trans)).To(Succeed())
			// Set up the pointers to be mapped
			Expect(definingSourceRefinement.SetAbstractConcept(definingSourceReferent, trans)).To(Succeed())
			Expect(definingTargetRefinement.SetAbstractConcept(definingTargetReferent, trans)).To(Succeed())
			// Add the pointer map
			definingAbstractPointer2AbstractPointerMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingAbstractPointer2AbstractPointerMap.SetOwningConcept(definingRefinement2RefinementMap, trans)).To(Succeed())
			Expect(definingAbstractPointer2AbstractPointerMap.SetLabel("definingAbstractPointer2AbstractPointerMap", trans)).To(Succeed())
			Expect(SetSource(definingAbstractPointer2AbstractPointerMap, definingSourceRefinement, core.AbstractConceptID, trans)).To(Succeed())
			Expect(SetTarget(definingAbstractPointer2AbstractPointerMap, definingTargetRefinement, core.AbstractConceptID, trans)).To(Succeed())

			// Trigger the map
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())

			// Diagnostics view
			graph := core.NewCrlGraph("AbstractPointer2AbstractPointerMapTest")
			Expect(graph.AddConceptRecursively(definingSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingTargetFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingMapFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceMapFolder, trans)).To(Succeed())
			Expect(graph.ExportDOT(tempDirPath, "AbstractPointer2AbstractPointerMapTest")).To(Succeed())

			// Check the result
			instanceRefinement2RefinementMap := FindMapForSource(instanceDomainMap, instanceSourceRefinement, trans)
			Expect(instanceRefinement2RefinementMap).ToNot(BeNil())
			instanceAbstractPointer2AbstractPointerMap := FindAttributeMapForSource(instanceDomainMap, instanceSourceRefinement, core.AbstractConceptID, trans)
			Expect(instanceAbstractPointer2AbstractPointerMap).ToNot(BeNil())
			Expect(instanceAbstractPointer2AbstractPointerMap.GetOwningConcept(trans).GetConceptID(trans)).To(Equal(instanceRefinement2RefinementMap.GetConceptID(trans)))
			instanceTargetReference := GetTarget(instanceRefinement2RefinementMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceAttributeTarget := GetTarget(instanceAbstractPointer2AbstractPointerMap, trans)
			Expect(instanceAttributeTarget).To(Equal(instanceTargetReference))

			instanceReferent2ReferentMap := FindMapForSource(instanceDomainMap, instanceSourceReferent, trans)
			Expect(instanceReferent2ReferentMap).ToNot(BeNil())
			instanceTargetReferent2 := GetTarget(instanceReferent2ReferentMap, trans)
			Expect(instanceTargetReferent2).ToNot(BeNil())

			switch instanceTargetReference.GetConceptType() {
			case core.Refinement:
				Expect(instanceTargetReference.GetAbstractConceptID(trans)).To(Equal(instanceTargetReferent2.GetConceptID(trans)))
			}
		})
		Specify("Refinement Refined Pointer to Refinement Refined Pointer", func() {
			// Establish the value that is going to be mapped
			Expect(instanceSourceRefinement.SetRefinedConcept(instanceSourceReferent, trans)).To(Succeed())
			// Set up the pointers to be mapped
			Expect(definingSourceRefinement.SetRefinedConcept(definingSourceReferent, trans)).To(Succeed())
			Expect(definingTargetRefinement.SetRefinedConcept(definingTargetReferent, trans)).To(Succeed())
			// Add the pointer map
			definingRefinedPointer2RefinedPointerMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingRefinedPointer2RefinedPointerMap.SetOwningConcept(definingRefinement2RefinementMap, trans)).To(Succeed())
			Expect(definingRefinedPointer2RefinedPointerMap.SetLabel("definingRefinedPointer2RefinedPointerMap", trans)).To(Succeed())
			Expect(SetSource(definingRefinedPointer2RefinedPointerMap, definingSourceRefinement, core.RefinedConceptID, trans)).To(Succeed())
			Expect(SetTarget(definingRefinedPointer2RefinedPointerMap, definingTargetRefinement, core.RefinedConceptID, trans)).To(Succeed())

			// Trigger the map
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())

			// Diagnostics view
			graph := core.NewCrlGraph("AbstractPointer2AbstractPointerMapTest")
			Expect(graph.AddConceptRecursively(definingSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingTargetFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingMapFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceMapFolder, trans)).To(Succeed())
			Expect(graph.ExportDOT(tempDirPath, "AbstractPointer2AbstractPointerMapTest")).To(Succeed())

			// Check the result
			instanceRefinement2RefinementMap := FindMapForSource(instanceDomainMap, instanceSourceRefinement, trans)
			Expect(instanceRefinement2RefinementMap).ToNot(BeNil())
			instanceRefinedPointer2RefinedPointerMap := FindAttributeMapForSource(instanceDomainMap, instanceSourceRefinement, core.RefinedConceptID, trans)
			Expect(instanceRefinedPointer2RefinedPointerMap).ToNot(BeNil())
			Expect(instanceRefinedPointer2RefinedPointerMap.GetOwningConcept(trans).GetConceptID(trans)).To(Equal(instanceRefinement2RefinementMap.GetConceptID(trans)))
			instanceTargetReference := GetTarget(instanceRefinement2RefinementMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceAttributeTarget := GetTarget(instanceRefinedPointer2RefinedPointerMap, trans)
			Expect(instanceAttributeTarget).To(Equal(instanceTargetReference))

			instanceReferent2ReferentMap := FindMapForSource(instanceDomainMap, instanceSourceReferent, trans)
			Expect(instanceReferent2ReferentMap).ToNot(BeNil())
			instanceTargetReferent2 := GetTarget(instanceReferent2ReferentMap, trans)
			Expect(instanceTargetReferent2).ToNot(BeNil())

			switch instanceTargetReference.GetConceptType() {
			case core.Refinement:
				Expect(instanceTargetReference.GetRefinedConceptID(trans)).To(Equal(instanceTargetReferent2.GetConceptID(trans)))
			}
		})
		Specify("Literal Value to Literal Value", func() {
			testString := "TestString"
			// Establish the value that is going to be mapped
			Expect(instanceSourceLiteral.SetLiteralValue(testString, trans)).To(Succeed())
			// // Set up the pointers to be mapped
			// Expect(definingSourceRefinement.SetRefinedConcept(definingSourceReferent, trans)).To(Succeed())
			// Expect(definingTargetRefinement.SetRefinedConcept(definingTargetReferent, trans)).To(Succeed())
			// Add the pointer map
			definingLiteralValue2LiteralValueMap, err := uOfD.CreateReplicateAsRefinementFromURI(CrlOneToOneMapURI, trans)
			Expect(err).To(BeNil())
			Expect(definingLiteralValue2LiteralValueMap.SetOwningConcept(definingLiteral2LiteralMap, trans)).To(Succeed())
			Expect(definingLiteralValue2LiteralValueMap.SetLabel("definingLiteralValue2LiteralValueMap", trans)).To(Succeed())
			Expect(SetSource(definingLiteralValue2LiteralValueMap, definingSourceLiteral, core.LiteralValue, trans)).To(Succeed())
			Expect(SetTarget(definingLiteralValue2LiteralValueMap, definingTargetLiteral, core.LiteralValue, trans)).To(Succeed())

			// Trigger the map
			Expect(SetSource(instanceDomainMap, instanceSourceDomain, core.NoAttribute, trans)).To(Succeed())

			// Diagnostics view
			graph := core.NewCrlGraph("AbstractPointer2AbstractPointerMapTest")
			Expect(graph.AddConceptRecursively(definingSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingTargetFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(definingMapFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceSourceFolder, trans)).To(Succeed())
			Expect(graph.AddConceptRecursively(instanceMapFolder, trans)).To(Succeed())
			Expect(graph.ExportDOT(tempDirPath, "AbstractPointer2AbstractPointerMapTest")).To(Succeed())

			// Check the result
			instanceLiteral2LiteralMap := FindMapForSource(instanceDomainMap, instanceSourceLiteral, trans)
			Expect(instanceLiteral2LiteralMap).ToNot(BeNil())
			instanceLiteralValue2LiteralValueMap := FindAttributeMapForSource(instanceDomainMap, instanceSourceLiteral, core.LiteralValue, trans)
			Expect(instanceLiteralValue2LiteralValueMap).ToNot(BeNil())
			Expect(instanceLiteralValue2LiteralValueMap.GetOwningConcept(trans).GetConceptID(trans)).To(Equal(instanceLiteral2LiteralMap.GetConceptID(trans)))
			instanceTargetReference := GetTarget(instanceLiteral2LiteralMap, trans)
			Expect(instanceTargetReference).ToNot(BeNil())
			instanceAttributeTarget := GetTarget(instanceLiteralValue2LiteralValueMap, trans)
			Expect(instanceAttributeTarget).To(Equal(instanceTargetReference))

			switch instanceTargetReference.GetConceptType() {
			case core.Literal:
				Expect(instanceTargetReference.GetLiteralValue(trans)).To(Equal(testString))
			}
		})
	})
})
