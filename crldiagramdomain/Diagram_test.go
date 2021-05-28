package crldiagramdomain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
	"golang.org/x/image/math/fixed"
)

var _ = Describe("CrlDiagramtest", func() {
	var uOfD *core.UniverseOfDiscourse
	var hl *core.HeldLocks

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("TestBuildCrlDiagramDomain", func() {
		Specify("CrlDiagramDomain should be created", func() {
			// CrlDiagramDomain
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, hl)
			Expect(builtCrlDiagramDomain).ToNot(BeNil())
			Expect(builtCrlDiagramDomain.GetURI(hl)).To(Equal(CrlDiagramDomainURI))
			Expect(uOfD.GetElementWithURI(CrlDiagramDomainURI)).To(Equal(builtCrlDiagramDomain))
		})

		Specify("CrlDiagram should build correctly", func() {
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, hl)
			crlDiagram := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramURI, hl)
			Expect(crlDiagram).ToNot(BeNil())
			crlDiagramWidth := crlDiagram.GetFirstOwnedLiteralWithURI(CrlDiagramWidthURI, hl)
			Expect(crlDiagramWidth).ToNot(BeNil())
			crlDiagramHeight := crlDiagram.GetFirstOwnedLiteralWithURI(CrlDiagramHeightURI, hl)
			Expect(crlDiagramHeight).ToNot(BeNil())
		})

		Specify("CrlDiagramNode should build correctly", func() {
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, hl)
			crlDiagramNode := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramNodeURI, hl)
			Expect(crlDiagramNode).ToNot(BeNil())
			crlDiagramNodeModelReference := crlDiagramNode.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)
			Expect(crlDiagramNodeModelReference).ToNot(BeNil())
			crlDiagramNodeAbstractionDisplayLabel := crlDiagramNode.GetFirstOwnedLiteralRefinedFromURI(CrlDiagramElementAbstractionDisplayLabelURI, hl)
			Expect(crlDiagramNodeAbstractionDisplayLabel).ToNot(BeNil())
			crlDiagramNodeDisplayLabel := crlDiagramNode.GetFirstOwnedLiteralRefinedFromURI(CrlDiagramElementDisplayLabelURI, hl)
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
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, hl)
			crlDiagramLink := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramLinkURI, hl)
			Expect(crlDiagramLink).ToNot(BeNil())
			crlDiagramLinkSourceReference := crlDiagramLink.GetFirstOwnedReferenceWithURI(CrlDiagramLinkSourceURI, hl)
			Expect(crlDiagramLinkSourceReference).ToNot(BeNil())
			crlDiagramLinkTargetReference := crlDiagramLink.GetFirstOwnedReferenceWithURI(CrlDiagramLinkSourceURI, hl)
			Expect(crlDiagramLinkTargetReference).ToNot(BeNil())
		})

		Specify("CrlDiagramPointer and refinements should build correctly", func() {
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, hl)
			crlPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramPointerURI, hl)
			Expect(crlPointer).ToNot(BeNil())
			crlOwnerPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramOwnerPointerURI, hl)
			Expect(crlOwnerPointer).ToNot(BeNil())
			Expect(crlOwnerPointer.IsRefinementOf(crlPointer, hl)).To(BeTrue())
			crlAbstractPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramAbstractPointerURI, hl)
			Expect(crlAbstractPointer).ToNot(BeNil())
			Expect(crlAbstractPointer.IsRefinementOf(crlPointer, hl)).To(BeTrue())
			crlRefinedPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramRefinedPointerURI, hl)
			Expect(crlRefinedPointer).ToNot(BeNil())
			Expect(crlRefinedPointer.IsRefinementOf(crlPointer, hl)).To(BeTrue())
			crlElementPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramElementPointerURI, hl)
			Expect(crlElementPointer).ToNot(BeNil())
			Expect(crlElementPointer.IsRefinementOf(crlPointer, hl)).To(BeTrue())
		})
	})
	Describe("Test convenience functions", func() {
		Specify("GetReferencedElement and SetReferencedElement should work", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			target, _ := uOfD.NewElement(hl)
			SetReferencedModelElement(node, target, hl)
			Expect(GetReferencedModelElement(node, hl)).To(Equal(target))
			Expect(node.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)).ToNot(BeNil())
			Expect(node.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl).GetReferencedConcept(hl)).To(Equal(target))
		})
		Specify("SetReferencedElement should gracefully handle a nil argument", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			target, _ := uOfD.NewElement(hl)
			SetReferencedModelElement(nil, target, hl)
		})
		Specify("SetReferencedElement should gracefully handle a nil argument", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			SetReferencedModelElement(node, nil, hl)
		})
		Specify("Test GetAbstractionDisplayLabel and SetAbstractionDisplayLabel", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			displayLabel := "displayLabel"
			SetAbstractionDisplayLabel(node, displayLabel, hl)
			Expect(GetAbstractionDisplayLabel(node, hl)).To(Equal(displayLabel))
		})
		Specify("Test GetDisplayLabel and SetDisplayLabel", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			displayLabel := "displayLabel"
			SetDisplayLabel(node, displayLabel, hl)
			Expect(GetDisplayLabel(node, hl)).To(Equal(displayLabel))
		})
		Specify("Test GetNodeHeight and SetNodeHeight", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			value := 123.45
			SetNodeHeight(node, value, hl)
			Expect(GetNodeHeight(node, hl)).To(Equal(value))
		})
		Specify("Test GetNodeWidth and SetNodeWidth", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			value := 123.45
			SetNodeWidth(node, value, hl)
			Expect(GetNodeWidth(node, hl)).To(Equal(value))
		})
		Specify("Test GetNodeX and SetNodeX", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			value := 123.45
			SetNodeX(node, value, hl)
			Expect(GetNodeX(node, hl)).To(Equal(value))
		})
		Specify("Test GetNodeY and SetNodeY", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
			value := 123.45
			SetNodeY(node, value, hl)
			Expect(GetNodeY(node, hl)).To(Equal(value))
		})
		Specify("GetLinkSource and SetLinkSource should work", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			link, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramLinkURI, hl)
			target, _ := uOfD.NewElement(hl)
			SetLinkSource(link, target, hl)
			Expect(GetLinkSource(link, hl)).To(Equal(target))
			Expect(link.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, hl)).ToNot(BeNil())
			Expect(link.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, hl).GetReferencedConcept(hl)).To(Equal(target))
		})
		Specify("GetLinkTarget and SetLinkTarget should work", func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
			link, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramLinkURI, hl)
			target, _ := uOfD.NewElement(hl)
			SetLinkTarget(link, target, hl)
			Expect(GetLinkTarget(link, hl)).To(Equal(target))
			Expect(link.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, hl)).ToNot(BeNil())
			Expect(link.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, hl).GetReferencedConcept(hl)).To(Equal(target))
		})
	})
	Describe("Test Fixed-to-Float conversions", func() {
		Specify("Zero should equal zero", func() {
			fixedX := fixed.I(0)
			floatX := Int26_6ToFloat(fixedX)
			Expect(floatX).To(Equal(0.0))
		})
		Specify("One should equal one", func() {
			fixedX := fixed.I(1)
			floatX := Int26_6ToFloat(fixedX)
			Expect(floatX).To(Equal(1.0))
		})
		Specify("Minus One should equal minus one", func() {
			fixedX := fixed.I(-1)
			floatX := Int26_6ToFloat(fixedX)
			Expect(floatX).To(Equal(-1.0))
		})
		Specify("0.5 should equal 0.5", func() {
			var fixedX fixed.Int26_6 = 1 << 5
			floatX := Int26_6ToFloat(fixedX)
			Expect(floatX).To(Equal(0.5))
		})
	})

	Describe("Test New... functions", func() {
		BeforeEach(func() {
			BuildCrlDiagramDomain(uOfD, hl)
			hl.ReleaseLocksAndWait()
		})
		Specify("Creating a new node shoud work", func() {
			node, err := NewDiagramNode(uOfD, hl)
			Expect(err).To(BeNil())
			Expect(node).ToNot(BeNil())
			Expect(node.IsRefinementOfURI(CrlDiagramNodeURI, hl)).To(BeTrue())
		})
		Specify("Creating a new reference link shoud work", func() {
			link, err := NewDiagramReferenceLink(uOfD, hl)
			Expect(err).To(BeNil())
			Expect(link).ToNot(BeNil())
			Expect(link.IsRefinementOfURI(CrlDiagramReferenceLinkURI, hl)).To(BeTrue())
		})
		Specify("Creating a new refinement link shoud work", func() {
			link, err := NewDiagramRefinementLink(uOfD, hl)
			Expect(err).To(BeNil())
			Expect(link).ToNot(BeNil())
			Expect(link.IsRefinementOfURI(CrlDiagramRefinementLinkURI, hl)).To(BeTrue())
		})
		Specify("Creating a new ownerPointer shoud work", func() {
			ownerPointer, err := NewDiagramOwnerPointer(uOfD, hl)
			Expect(err).To(BeNil())
			Expect(ownerPointer).ToNot(BeNil())
			Expect(ownerPointer.IsRefinementOfURI(CrlDiagramOwnerPointerURI, hl)).To(BeTrue())
		})
		Specify("Creating a new elementPointer shoud work", func() {
			elementPointer, err := NewDiagramElementPointer(uOfD, hl)
			Expect(err).To(BeNil())
			Expect(elementPointer).ToNot(BeNil())
			Expect(elementPointer.IsRefinementOfURI(CrlDiagramElementPointerURI, hl)).To(BeTrue())
		})
		Specify("Creating a new abstractPointer shoud work", func() {
			abstractPointer, err := NewDiagramAbstractPointer(uOfD, hl)
			Expect(err).To(BeNil())
			Expect(abstractPointer).ToNot(BeNil())
			Expect(abstractPointer.IsRefinementOfURI(CrlDiagramAbstractPointerURI, hl)).To(BeTrue())
		})
		Specify("Creating a new refinedPointer shoud work", func() {
			refinedPointer, err := NewDiagramRefinedPointer(uOfD, hl)
			Expect(err).To(BeNil())
			Expect(refinedPointer).ToNot(BeNil())
			Expect(refinedPointer.IsRefinementOfURI(CrlDiagramRefinedPointerURI, hl)).To(BeTrue())
		})
	})
})
