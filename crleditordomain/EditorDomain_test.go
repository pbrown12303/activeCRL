package crleditordomain

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructuresdomain"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
)

var _ = Describe("EdditorDomain tests", func() {
	Specify("Editor domain creation should be idempotent", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewTransaction()
		crldiagramdomain.BuildCrlDiagramDomain(uOfD1, hl1)
		crldatastructuresdomain.BuildCrlDataStructuresDomain(uOfD1, hl1)
		Expect(BuildEditorDomain(uOfD1, hl1)).ShouldNot(BeNil())
		cs1 := uOfD1.GetElementWithURI(EditorDomainURI)
		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewTransaction()
		crldiagramdomain.BuildCrlDiagramDomain(uOfD2, hl2)
		crldatastructuresdomain.BuildCrlDataStructuresDomain(uOfD2, hl2)
		Expect(BuildEditorDomain(uOfD2, hl2)).ShouldNot(BeNil())
		cs2 := uOfD2.GetElementWithURI(EditorDomainURI)
		Expect(core.RecursivelyEquivalent(cs1, hl1, cs2, hl2, true)).To(BeTrue())
	})
	Specify("Refinement of OpenDiagrams list with URIs should serialize and deserialzie properly", func() {
		uOfD1 := core.NewUniverseOfDiscourse()
		hl1 := uOfD1.NewTransaction()
		crldiagramdomain.BuildCrlDiagramDomain(uOfD1, hl1)
		crldatastructuresdomain.BuildCrlDataStructuresDomain(uOfD1, hl1)
		Expect(BuildEditorDomain(uOfD1, hl1)).ShouldNot(BeNil())

		uOfD2 := core.NewUniverseOfDiscourse()
		hl2 := uOfD2.NewTransaction()
		defer hl2.ReleaseLocks()
		crldiagramdomain.BuildCrlDiagramDomain(uOfD2, hl2)
		crldatastructuresdomain.BuildCrlDataStructuresDomain(uOfD2, hl2)
		Expect(BuildEditorDomain(uOfD2, hl2)).ShouldNot(BeNil())

		// Refine settings and persist
		coreSettings1 := uOfD1.GetElementWithURI(EditorSettingsURI)
		coreOpenDiagrams1 := uOfD1.GetElementWithURI(EditorOpenDiagramsURI)
		firstMemberReference1FromURI := coreOpenDiagrams1.GetFirstOwnedReferenceRefinedFromURI(crldatastructuresdomain.CrlStringListReferenceToFirstMemberLiteralURI, hl1)
		Expect(firstMemberReference1FromURI).ToNot(BeNil())
		settings1, err := uOfD1.CreateReplicateAsRefinementFromURI(EditorSettingsURI, hl1)
		Expect(err).To(BeNil())
		openDiagrams1 := settings1.GetFirstOwnedConceptRefinedFromURI(EditorOpenDiagramsURI, hl1)
		Expect(openDiagrams1).ToNot(BeNil())
		serialized1, err2 := uOfD1.MarshalDomain(settings1, hl1)
		Expect(err2).To(BeNil())

		// Recover settings
		coreSettings2 := uOfD1.GetElementWithURI(EditorSettingsURI)
		settings2, err3 := uOfD2.RecoverDomain(serialized1, hl2)
		Expect(err3).To(BeNil())
		Expect(settings2).ToNot(BeNil())
		openDiagrams2 := uOfD2.GetElement(openDiagrams1.GetConceptID(hl1))
		Expect(openDiagrams2).ToNot(BeNil())

		// Now the tests
		Expect(core.Equivalent(coreSettings1, hl1, coreSettings2, hl2))
		Expect(core.Equivalent(settings1, hl1, settings2, hl2)).To(BeTrue())
		Expect(core.Equivalent(openDiagrams1, hl1, openDiagrams2, hl2)).To(BeTrue())
		openDiagrams1FirstElementRef, err5 := crldatastructuresdomain.GetFirstMemberLiteral(openDiagrams1, hl1)
		Expect(err5).To(BeNil())
		Expect(openDiagrams1FirstElementRef).To(BeNil())
		openDiagrams2FirstElementRef, err7 := crldatastructuresdomain.GetFirstMemberLiteral(openDiagrams2, hl2)
		Expect(err7).To(BeNil())
		Expect(openDiagrams2FirstElementRef).To(BeNil())
	})

})
