package crldiagram

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("CrlDiagramtest", func() {
	var uOfD core.UniverseOfDiscourse
	var hl *core.HeldLocks

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("TestBuildCrlDiagramConceptSpace", func() {
		Specify("CrlDiagramConceptSpace should be created", func() {
			// CrlDiagramConceptSpace
			builtCrlDiagramConceptSpace := BuildCrlDiagramConceptSpace(uOfD, hl)
			Expect(builtCrlDiagramConceptSpace).ToNot(BeNil())
			Expect(builtCrlDiagramConceptSpace.GetURI(hl)).To(Equal(CrlDiagramConceptSpaceURI))
			Expect(uOfD.GetElementWithURI(CrlDiagramConceptSpaceURI)).To(Equal(builtCrlDiagramConceptSpace))
		})

		Specify("CrlDiagram should build correctly", func() {
			builtCrlDiagramConceptSpace := BuildCrlDiagramConceptSpace(uOfD, hl)
			crlDiagram := builtCrlDiagramConceptSpace.GetFirstChildWithURI(CrlDiagramURI, hl)
			Expect(crlDiagram).ToNot(BeNil())
			crlDiagramWidth := crlDiagram.GetFirstChildLiteralWithURI(CrlDiagramWidthURI, hl)
			Expect(crlDiagramWidth).ToNot(BeNil())
			crlDiagramHeight := crlDiagram.GetFirstChildLiteralWithURI(CrlDiagramHeightURI, hl)
			Expect(crlDiagramHeight).ToNot(BeNil())
		})

		Specify("CrlDiagramNode should build correctly", func() {
			builtCrlDiagramConceptSpace := BuildCrlDiagramConceptSpace(uOfD, hl)
			crlDiagramNode := builtCrlDiagramConceptSpace.GetFirstChildWithURI(CrlDiagramNodeURI, hl)
			Expect(crlDiagramNode).ToNot(BeNil())
			crlDiagramNodeModelReference := crlDiagramNode.GetFirstChildReferenceWithURI(CrlDiagramNodeModelReferenceURI, hl)
			Expect(crlDiagramNodeModelReference).ToNot(BeNil())
			crlDiagramNodeDisplayLabel := crlDiagramNode.GetFirstChildLiteralWithURI(CrlDiagramNodeDisplayLabelURI, hl)
			Expect(crlDiagramNodeDisplayLabel).ToNot(BeNil())
			crlDiagramNodeX := crlDiagramNode.GetFirstChildLiteralWithURI(CrlDiagramNodeXURI, hl)
			Expect(crlDiagramNodeX).ToNot(BeNil())
			crlDiagramNodeY := crlDiagramNode.GetFirstChildLiteralWithURI(CrlDiagramNodeYURI, hl)
			Expect(crlDiagramNodeY).ToNot(BeNil())
			crlDiagramNodeHeight := crlDiagramNode.GetFirstChildLiteralWithURI(CrlDiagramNodeHeightURI, hl)
			Expect(crlDiagramNodeHeight).ToNot(BeNil())
			crlDiagramNodeWidth := crlDiagramNode.GetFirstChildLiteralWithURI(CrlDiagramNodeWidthURI, hl)
			Expect(crlDiagramNodeWidth).ToNot(BeNil())
		})

		Specify("CrlDiagramLink should build correctly", func() {
			builtCrlDiagramConceptSpace := BuildCrlDiagramConceptSpace(uOfD, hl)
			crlDiagramLink := builtCrlDiagramConceptSpace.GetFirstChildWithURI(CrlDiagramLinkURI, hl)
			Expect(crlDiagramLink).ToNot(BeNil())
		})
	})
	Describe("Test convenience functions", func() {
		Specify("GetReferencedElement and SetReferencedElement should work", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			target, _ := uOfD.NewElement(hl)
			SetReferencedElement(node, target, hl)
			Expect(GetReferencedElement(node, hl)).To(Equal(target))
			Expect(node.GetFirstChildReferenceWithAbstractionURI(CrlDiagramNodeModelReferenceURI, hl)).ToNot(BeNil())
			Expect(node.GetFirstChildReferenceWithAbstractionURI(CrlDiagramNodeModelReferenceURI, hl).GetReferencedConcept(hl)).To(Equal(target))
		})
		Specify("SetReferencedElement should gracefully handle a nil argument", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			target, _ := uOfD.NewElement(hl)
			SetReferencedElement(nil, target, hl)
		})
		Specify("SetReferencedElement should gracefully handle a nil argument", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			SetReferencedElement(node, nil, hl)
		})
	})
})
