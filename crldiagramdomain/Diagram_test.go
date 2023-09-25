package crldiagramdomain

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
	"golang.org/x/image/math/fixed"
)

var _ = Describe("CrlDiagramtest", func() {
	var uOfD *core.UniverseOfDiscourse
	var trans *core.Transaction

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
	})

	AfterEach(func() {
		trans.ReleaseLocks()
	})

	Describe("TestBuildCrlDiagramDomain", func() {
		Specify("CrlDiagramDomain should be created", func() {
			// CrlDiagramDomain
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, trans)
			Expect(builtCrlDiagramDomain).ToNot(BeNil())
			Expect(builtCrlDiagramDomain.GetURI(trans)).To(Equal(CrlDiagramDomainURI))
			Expect(uOfD.GetElementWithURI(CrlDiagramDomainURI)).To(Equal(builtCrlDiagramDomain))
		})

		Specify("CrlDiagram should build correctly", func() {
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, trans)
			crlDiagram := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramURI, trans)
			Expect(crlDiagram).ToNot(BeNil())
			crlDiagramWidth := crlDiagram.GetFirstOwnedLiteralWithURI(CrlDiagramWidthURI, trans)
			Expect(crlDiagramWidth).ToNot(BeNil())
			crlDiagramHeight := crlDiagram.GetFirstOwnedLiteralWithURI(CrlDiagramHeightURI, trans)
			Expect(crlDiagramHeight).ToNot(BeNil())
		})

		Specify("CrlDiagramNode should build correctly", func() {
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, trans)
			crlDiagramNode := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramNodeURI, trans)
			Expect(crlDiagramNode).ToNot(BeNil())
			crlDiagramNodeModelReference := crlDiagramNode.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
			Expect(crlDiagramNodeModelReference).ToNot(BeNil())
			crlDiagramNodeAbstractionDisplayLabel := crlDiagramNode.GetFirstOwnedLiteralRefinedFromURI(CrlDiagramElementAbstractionDisplayLabelURI, trans)
			Expect(crlDiagramNodeAbstractionDisplayLabel).ToNot(BeNil())
			crlDiagramNodeDisplayLabel := crlDiagramNode.GetFirstOwnedLiteralRefinedFromURI(CrlDiagramElementDisplayLabelURI, trans)
			Expect(crlDiagramNodeDisplayLabel).ToNot(BeNil())
			crlDiagramNodeX := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeXURI, trans)
			Expect(crlDiagramNodeX).ToNot(BeNil())
			crlDiagramNodeY := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeYURI, trans)
			Expect(crlDiagramNodeY).ToNot(BeNil())
			crlDiagramNodeHeight := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeHeightURI, trans)
			Expect(crlDiagramNodeHeight).ToNot(BeNil())
			crlDiagramNodeWidth := crlDiagramNode.GetFirstOwnedLiteralWithURI(CrlDiagramNodeWidthURI, trans)
			Expect(crlDiagramNodeWidth).ToNot(BeNil())
		})

		Specify("CrlDiagramLink should build correctly", func() {
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, trans)
			crlDiagramLink := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramLinkURI, trans)
			Expect(crlDiagramLink).ToNot(BeNil())
			crlDiagramLinkSourceReference := crlDiagramLink.GetFirstOwnedReferenceWithURI(CrlDiagramLinkSourceURI, trans)
			Expect(crlDiagramLinkSourceReference).ToNot(BeNil())
			crlDiagramLinkTargetReference := crlDiagramLink.GetFirstOwnedReferenceWithURI(CrlDiagramLinkSourceURI, trans)
			Expect(crlDiagramLinkTargetReference).ToNot(BeNil())
		})

		Specify("CrlDiagramPointer and refinements should build correctly", func() {
			builtCrlDiagramDomain := BuildCrlDiagramDomain(uOfD, trans)
			crlPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramPointerURI, trans)
			Expect(crlPointer).ToNot(BeNil())
			crlOwnerPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramOwnerPointerURI, trans)
			Expect(crlOwnerPointer).ToNot(BeNil())
			Expect(crlOwnerPointer.IsRefinementOf(crlPointer, trans)).To(BeTrue())
			crlAbstractPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramAbstractPointerURI, trans)
			Expect(crlAbstractPointer).ToNot(BeNil())
			Expect(crlAbstractPointer.IsRefinementOf(crlPointer, trans)).To(BeTrue())
			crlRefinedPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramRefinedPointerURI, trans)
			Expect(crlRefinedPointer).ToNot(BeNil())
			Expect(crlRefinedPointer.IsRefinementOf(crlPointer, trans)).To(BeTrue())
			crlElementPointer := builtCrlDiagramDomain.GetFirstOwnedConceptWithURI(CrlDiagramElementPointerURI, trans)
			Expect(crlElementPointer).ToNot(BeNil())
			Expect(crlElementPointer.IsRefinementOf(crlPointer, trans)).To(BeTrue())
		})
	})
	Describe("Test convenience functions", func() {
		Specify("GetReferencedElement and SetReferencedElement should work", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
			target, _ := uOfD.NewElement(trans)
			SetReferencedModelConcept(node, target, trans)
			Expect(GetReferencedModelConcept(node, trans)).To(Equal(target))
			Expect(node.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)).ToNot(BeNil())
			Expect(node.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans).GetReferencedConcept(trans)).To(Equal(target))
		})
		Specify("SetReferencedElement should gracefully handle a nil argument", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			target, _ := uOfD.NewElement(trans)
			SetReferencedModelConcept(nil, target, trans)
		})
		Specify("SetReferencedElement should gracefully handle a nil argument", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
			SetReferencedModelConcept(node, nil, trans)
		})
		Specify("Test GetAbstractionDisplayLabel and SetAbstractionDisplayLabel", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
			displayLabel := "displayLabel"
			SetAbstractionDisplayLabel(node, displayLabel, trans)
			Expect(GetAbstractionDisplayLabel(node, trans)).To(Equal(displayLabel))
		})
		Specify("Test GetDisplayLabel and SetDisplayLabel", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
			displayLabel := "displayLabel"
			SetDisplayLabel(node, displayLabel, trans)
			Expect(GetDisplayLabel(node, trans)).To(Equal(displayLabel))
		})
		Specify("Test GetNodeHeight and SetNodeHeight", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
			value := 123.45
			SetNodeHeight(node, value, trans)
			Expect(GetNodeHeight(node, trans)).To(Equal(value))
		})
		Specify("Test GetNodeWidth and SetNodeWidth", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
			value := 123.45
			SetNodeWidth(node, value, trans)
			Expect(GetNodeWidth(node, trans)).To(Equal(value))
		})
		Specify("Test GetNodeX and SetNodeX", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
			value := 123.45
			SetNodeX(node, value, trans)
			Expect(GetNodeX(node, trans)).To(Equal(value))
		})
		Specify("Test GetNodeY and SetNodeY", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			node, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
			value := 123.45
			SetNodeY(node, value, trans)
			Expect(GetNodeY(node, trans)).To(Equal(value))
		})
		Specify("GetLinkSource and SetLinkSource should work", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			link, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramLinkURI, trans)
			target, _ := uOfD.NewElement(trans)
			SetLinkSource(link, target, trans)
			Expect(GetLinkSource(link, trans)).To(Equal(target))
			Expect(link.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)).ToNot(BeNil())
			Expect(link.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans).GetReferencedConcept(trans)).To(Equal(target))
		})
		Specify("GetLinkTarget and SetLinkTarget should work", func() {
			BuildCrlDiagramDomain(uOfD, trans)
			link, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramLinkURI, trans)
			target, _ := uOfD.NewElement(trans)
			SetLinkTarget(link, target, trans)
			Expect(GetLinkTarget(link, trans)).To(Equal(target))
			Expect(link.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)).ToNot(BeNil())
			Expect(link.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans).GetReferencedConcept(trans)).To(Equal(target))
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
			BuildCrlDiagramDomain(uOfD, trans)
		})
		Specify("Creating a new node shoud work", func() {
			node, err := NewDiagramNode(uOfD, trans)
			Expect(err).To(BeNil())
			Expect(node).ToNot(BeNil())
			Expect(node.IsRefinementOfURI(CrlDiagramNodeURI, trans)).To(BeTrue())
		})
		Specify("Creating a new reference link shoud work", func() {
			link, err := NewDiagramReferenceLink(uOfD, trans)
			Expect(err).To(BeNil())
			Expect(link).ToNot(BeNil())
			Expect(link.IsRefinementOfURI(CrlDiagramReferenceLinkURI, trans)).To(BeTrue())
		})
		Specify("Creating a new refinement link shoud work", func() {
			link, err := NewDiagramRefinementLink(uOfD, trans)
			Expect(err).To(BeNil())
			Expect(link).ToNot(BeNil())
			Expect(link.IsRefinementOfURI(CrlDiagramRefinementLinkURI, trans)).To(BeTrue())
		})
		Specify("Creating a new ownerPointer shoud work", func() {
			ownerPointer, err := NewDiagramOwnerPointer(uOfD, trans)
			Expect(err).To(BeNil())
			Expect(ownerPointer).ToNot(BeNil())
			Expect(ownerPointer.IsRefinementOfURI(CrlDiagramOwnerPointerURI, trans)).To(BeTrue())
		})
		Specify("Creating a new elementPointer shoud work", func() {
			elementPointer, err := NewDiagramElementPointer(uOfD, trans)
			Expect(err).To(BeNil())
			Expect(elementPointer).ToNot(BeNil())
			Expect(elementPointer.IsRefinementOfURI(CrlDiagramElementPointerURI, trans)).To(BeTrue())
		})
		Specify("Creating a new abstractPointer shoud work", func() {
			abstractPointer, err := NewDiagramAbstractPointer(uOfD, trans)
			Expect(err).To(BeNil())
			Expect(abstractPointer).ToNot(BeNil())
			Expect(abstractPointer.IsRefinementOfURI(CrlDiagramAbstractPointerURI, trans)).To(BeTrue())
		})
		Specify("Creating a new refinedPointer shoud work", func() {
			refinedPointer, err := NewDiagramRefinedPointer(uOfD, trans)
			Expect(err).To(BeNil())
			Expect(refinedPointer).ToNot(BeNil())
			Expect(refinedPointer.IsRefinementOfURI(CrlDiagramRefinedPointerURI, trans)).To(BeTrue())
		})
	})
})
