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
			crlDiagram := builtCrlDiagramConceptSpace.GetFirstOwnedConceptWithURI(CrlDiagramURI, hl)
			Expect(crlDiagram).ToNot(BeNil())
			crlDiagramWidth := crlDiagram.GetFirstOwnedLiteralWithURI(CrlDiagramWidthURI, hl)
			Expect(crlDiagramWidth).ToNot(BeNil())
			crlDiagramHeight := crlDiagram.GetFirstOwnedLiteralWithURI(CrlDiagramHeightURI, hl)
			Expect(crlDiagramHeight).ToNot(BeNil())
		})

		Specify("CrlDiagramNode should build correctly", func() {
			builtCrlDiagramConceptSpace := BuildCrlDiagramConceptSpace(uOfD, hl)
			crlDiagramNode := builtCrlDiagramConceptSpace.GetFirstOwnedConceptWithURI(CrlDiagramNodeURI, hl)
			Expect(crlDiagramNode).ToNot(BeNil())
			crlDiagramNodeModelReference := crlDiagramNode.GetFirstOwnedReferenceWithURI(CrlDiagramNodeModelReferenceURI, hl)
			Expect(crlDiagramNodeModelReference).ToNot(BeNil())
			crlDiagramNodeDisplayLabel := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeDisplayLabelURI, hl)
			Expect(crlDiagramNodeDisplayLabel).ToNot(BeNil())
			crlDiagramNodeX := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeXURI, hl)
			Expect(crlDiagramNodeX).ToNot(BeNil())
			crlDiagramNodeY := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeYURI, hl)
			Expect(crlDiagramNodeY).ToNot(BeNil())
			crlDiagramNodeHeight := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeHeightURI, hl)
			Expect(crlDiagramNodeHeight).ToNot(BeNil())
			crlDiagramNodeWidth := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeWidthURI, hl)
			Expect(crlDiagramNodeWidth).ToNot(BeNil())
		})

		Specify("CrlDiagramLink should build correctly", func() {
			builtCrlDiagramConceptSpace := BuildCrlDiagramConceptSpace(uOfD, hl)
			crlDiagramLink := builtCrlDiagramConceptSpace.GetFirstOwnedConceptWithURI(CrlDiagramLinkURI, hl)
			Expect(crlDiagramLink).ToNot(BeNil())
		})
	})
	Describe("Test convenience functions", func() {
		Specify("GetReferencedElement and SetReferencedElement should work", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			target, _ := uOfD.NewElement(hl)
			SetReferencedModelElement(node, target, hl)
			Expect(GetReferencedModelElement(node, hl)).To(Equal(target))
			Expect(node.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramNodeModelReferenceURI, hl)).ToNot(BeNil())
			Expect(node.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramNodeModelReferenceURI, hl).GetReferencedConcept(hl)).To(Equal(target))
		})
		Specify("SetReferencedElement should gracefully handle a nil argument", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			target, _ := uOfD.NewElement(hl)
			SetReferencedModelElement(nil, target, hl)
		})
		Specify("SetReferencedElement should gracefully handle a nil argument", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			SetReferencedModelElement(node, nil, hl)
		})
		Specify("Test GetDisplayLabel and SetDisplayLabel", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			displayLabel := "displayLabel"
			SetDisplayLabel(node, displayLabel, hl)
			Expect(GetDisplayLabel(node, hl)).To(Equal(displayLabel))
		})
		Specify("Test GetNodeHeight and SetNodeHeight", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			value := 123.45
			SetNodeHeight(node, value, hl)
			Expect(GetNodeHeight(node, hl)).To(Equal(value))
		})
		Specify("Test GetNodeWidth and SetNodeWidth", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			value := 123.45
			SetNodeWidth(node, value, hl)
			Expect(GetNodeWidth(node, hl)).To(Equal(value))
		})
		Specify("Test GetNodeX and SetNodeX", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			value := 123.45
			SetNodeX(node, value, hl)
			Expect(GetNodeX(node, hl)).To(Equal(value))
		})
		Specify("Test GetNodeY and SetNodeY", func() {
			BuildCrlDiagramConceptSpace(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			value := 123.45
			SetNodeY(node, value, hl)
			Expect(GetNodeY(node, hl)).To(Equal(value))
		})
	})
})
