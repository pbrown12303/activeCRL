// Package crldiagramdomain defines the Diagram domain. This is a pre-defined domain that is, itself,
// represented as a CRLElement and identified with the CrlDiagramDomainURI. This concept space contains the prototypes of all Elements used to construct CrlDiagrams.
// Included are:
//
//		CrlDiagram: the diagram itself
//		CrlDiagramNode: a node in the diagram
//		CrlDiagramLink: a link in the diagram
//	    CrlDiagramPointer: a pointer shown as a link in the diagram
//
// These classes are intended to hold all of the information about the diagram that is not specific to the rendering engine.
//
// Intended Usage
// CRL Elements, in general, can have functions associated with them. When refinements of the elements are created, modified, or deleted, these functions are
// called. The strategy used for diagrams is to place all rendering-specific code in functions associated with the prototypes.
// This is accomplished using the FunctionCallManager.AddFunctionCall() method. Note that this registration is NOT done in the core diagram package, but
// rather in the package providing the rendering engine linkage. For example, the crlEditor package provides the linkages to JavaScript code that does
// the rendering in a browser.
//
// Instances of the prototpes can be conveniently instantiated using the core.CreateReplicateAsRefinementFromURI() function. This clones the prototype
// and, on an element-by-element basis, establishes a refinement relation between the instance elements and thier corresponding prototype elements.
// One essential side-effect of this is that changes that are made to the instnaces then trigger the execution of functions associated with the prototypes.
package crldiagramdomain

import (
	"log"
	"math"
	"strconv"

	"github.com/pkg/errors"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
	"github.com/pbrown12303/activeCRL/core"

	mapset "github.com/deckarep/golang-set"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"
)

// IconSize defines the height and width of icons
const IconSize = 16.0

// NodeLineWidth is the width of the line bordering the node image
const NodeLineWidth = 2.0

// NodePadWidth is the width of the padding surrounding the icon, displayLabel, and abstractionDisplayLabel
const NodePadWidth = 1.0

var goRegularFont *truetype.Font
var goBoldFont *truetype.Font

var go12PtRegularFace font.Face
var go12PtBoldFace font.Face

var go10PtRegularFace font.Face
var go10PtItalicFace font.Face

// CrlDiagramPrefix is the prefix for all URIs related to CrlDiagram
var CrlDiagramPrefix = "http://activeCrl.com"

// CrlDiagramDomainURI identifies concept space containing all concepts related to the CrlDiagram
var CrlDiagramDomainURI = CrlDiagramPrefix + "/CrlDiagramDomain"

// CrlDiagramURI identifies the CrlDiagram concept
var CrlDiagramURI = CrlDiagramDomainURI + "/" + "CrlDiagram"

// CrlDiagramWidthURI identifies the CrlDiagramWidth concept
var CrlDiagramWidthURI = CrlDiagramURI + "/" + "Width"

// CrlDiagramHeightURI identifies the CrlDiagramHeight concept
var CrlDiagramHeightURI = CrlDiagramURI + "/" + "Height"

// CrlDiagramElementURI identifies the CrlDiagramElement concept
var CrlDiagramElementURI = CrlDiagramDomainURI + "/" + "CrlDiagramElement"

// CrlDiagramElementModelReferenceURI identifies the reference to the model element represented by the element
var CrlDiagramElementModelReferenceURI = CrlDiagramElementURI + "/" + "ModelReference"

// CrlDiagramElementDisplayLabelURI identifies the display label concept to be used when displaying the element
var CrlDiagramElementDisplayLabelURI = CrlDiagramElementURI + "/" + "DisplayLabel"

// CrlDiagramElementAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the element
var CrlDiagramElementAbstractionDisplayLabelURI = CrlDiagramElementURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramElementLineColorURI identifies the line color to be used when displaying the element
var CrlDiagramElementLineColorURI = CrlDiagramElementURI + "/" + "LineColor"

// CrlDiagramElementBGColorURI identifies the background color to be used when displaying the element
var CrlDiagramElementBGColorURI = CrlDiagramElementURI + "/" + "BGColor"

// CrlDiagramLinkURI identifies the CrlDiagramLink concept
var CrlDiagramLinkURI = CrlDiagramDomainURI + "/" + "CrlDiagramLink"

// CrlDiagramLinkSourceURI identifies the concept that is the source of the link
var CrlDiagramLinkSourceURI = CrlDiagramLinkURI + "/" + "Source"

// CrlDiagramLinkTargetURI identifies the concept that is the target of the link
var CrlDiagramLinkTargetURI = CrlDiagramLinkURI + "/" + "Target"

// CrlDiagramNodeURI identifies the CrlDiagramNode concept
var CrlDiagramNodeURI = CrlDiagramDomainURI + "/" + "CrlDiagramNode"

// CrlDiagramNodeModelReferenceURI identifies the reference to the model element represented by the node
var CrlDiagramNodeModelReferenceURI = CrlDiagramNodeURI + "/" + "ModelReference"

// CrlDiagramNodeDisplayLabelURI identifies the display label concept to be used when displaying the node
var CrlDiagramNodeDisplayLabelURI = CrlDiagramNodeURI + "/" + "DisplayLabel"

// CrlDiagramNodeAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the node
var CrlDiagramNodeAbstractionDisplayLabelURI = CrlDiagramNodeURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramNodeLineColorURI identifies the line color to be used when displaying the element
var CrlDiagramNodeLineColorURI = CrlDiagramNodeURI + "/" + "LineColor"

// CrlDiagramNodeBGColorURI identifies the background color to be used when displaying the element
var CrlDiagramNodeBGColorURI = CrlDiagramNodeURI + "/" + "BGColor"

// CrlDiagramNodeXURI identifies the X coordinate of the node
var CrlDiagramNodeXURI = CrlDiagramNodeURI + "/" + "X"

// CrlDiagramNodeYURI identifies the Y coordinate of the node
var CrlDiagramNodeYURI = CrlDiagramNodeURI + "/" + "Y"

// CrlDiagramNodeHeightURI identifies the height of the node
var CrlDiagramNodeHeightURI = CrlDiagramNodeURI + "/" + "Height"

// CrlDiagramNodeWidthURI identifies the width of the node
var CrlDiagramNodeWidthURI = CrlDiagramNodeURI + "/" + "Width"

// CrlDiagramNodeDisplayLabelYOffsetURI identifies the Y offset for the display label within the node
var CrlDiagramNodeDisplayLabelYOffsetURI = CrlDiagramNodeURI + "/" + "DisplayLabelYOffset"

// CrlDiagramPointerURI identifies a pointer represented as a link
var CrlDiagramPointerURI = CrlDiagramDomainURI + "/" + "Pointer"

// CrlDiagramAbstractPointerURI identifies the Abstract of an Element represented as a link
var CrlDiagramAbstractPointerURI = CrlDiagramDomainURI + "/" + "AbstractPointer"

// CrlDiagramAbstractPointerModelReferenceURI identifies the reference to the model element represented by the link
var CrlDiagramAbstractPointerModelReferenceURI = CrlDiagramAbstractPointerURI + "/" + "ModelReference"

// CrlDiagramAbstractPointerDisplayLabelURI identifies the display label concept to be used when displaying the link
var CrlDiagramAbstractPointerDisplayLabelURI = CrlDiagramAbstractPointerURI + "/" + "DisplayLabel"

// CrlDiagramAbstractPointerAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the link
var CrlDiagramAbstractPointerAbstractionDisplayLabelURI = CrlDiagramAbstractPointerURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramAbstractPointerSourceURI identifies the concept that is the source of the link
var CrlDiagramAbstractPointerSourceURI = CrlDiagramAbstractPointerURI + "/" + "Source"

// CrlDiagramAbstractPointerTargetURI identifies the concept that is the target of the link
var CrlDiagramAbstractPointerTargetURI = CrlDiagramAbstractPointerURI + "/" + "Target"

// CrlDiagramElementPointerURI identifies the element pointer of a Reference represented as a link
var CrlDiagramElementPointerURI = CrlDiagramDomainURI + "/" + "ElementPointer"

// CrlDiagramElementPointerModelReferenceURI identifies the reference to the model element represented by the link
var CrlDiagramElementPointerModelReferenceURI = CrlDiagramElementPointerURI + "/" + "ModelReference"

// CrlDiagramElementPointerDisplayLabelURI identifies the display label concept to be used when displaying the link
var CrlDiagramElementPointerDisplayLabelURI = CrlDiagramElementPointerURI + "/" + "DisplayLabel"

// CrlDiagramElementPointerAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the link
var CrlDiagramElementPointerAbstractionDisplayLabelURI = CrlDiagramElementPointerURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramElementPointerSourceURI identifies the concept that is the source of the link
var CrlDiagramElementPointerSourceURI = CrlDiagramElementPointerURI + "/" + "Source"

// CrlDiagramElementPointerTargetURI identifies the concept that is the target of the link
var CrlDiagramElementPointerTargetURI = CrlDiagramElementPointerURI + "/" + "Target"

// CrlDiagramOwnerPointerURI identifies the owner of an Element represented as a link
var CrlDiagramOwnerPointerURI = CrlDiagramDomainURI + "/" + "OwnerPointer"

// CrlDiagramOwnerPointerModelReferenceURI identifies the reference to the model element represented by the link
var CrlDiagramOwnerPointerModelReferenceURI = CrlDiagramOwnerPointerURI + "/" + "ModelReference"

// CrlDiagramOwnerPointerDisplayLabelURI identifies the display label concept to be used when displaying the link
var CrlDiagramOwnerPointerDisplayLabelURI = CrlDiagramOwnerPointerURI + "/" + "DisplayLabel"

// CrlDiagramOwnerPointerAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the link
var CrlDiagramOwnerPointerAbstractionDisplayLabelURI = CrlDiagramOwnerPointerURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramOwnerPointerSourceURI identifies the concept that is the source of the link
var CrlDiagramOwnerPointerSourceURI = CrlDiagramOwnerPointerURI + "/" + "Source"

// CrlDiagramOwnerPointerTargetURI identifies the concept that is the target of the link
var CrlDiagramOwnerPointerTargetURI = CrlDiagramOwnerPointerURI + "/" + "Target"

// CrlDiagramRefinedPointerURI identifies the refined element of a Refinement represented as a link
var CrlDiagramRefinedPointerURI = CrlDiagramDomainURI + "/" + "RefinedPointer"

// CrlDiagramRefinedPointerModelReferenceURI identifies the reference to the model element represented by the link
var CrlDiagramRefinedPointerModelReferenceURI = CrlDiagramRefinedPointerURI + "/" + "ModelReference"

// CrlDiagramRefinedPointerDisplayLabelURI identifies the display label concept to be used when displaying the link
var CrlDiagramRefinedPointerDisplayLabelURI = CrlDiagramRefinedPointerURI + "/" + "DisplayLabel"

// CrlDiagramRefinedPointerAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the link
var CrlDiagramRefinedPointerAbstractionDisplayLabelURI = CrlDiagramRefinedPointerURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramRefinedPointerSourceURI identifies the concept that is the source of the link
var CrlDiagramRefinedPointerSourceURI = CrlDiagramRefinedPointerURI + "/" + "Source"

// CrlDiagramRefinedPointerTargetURI identifies the concept that is the target of the link
var CrlDiagramRefinedPointerTargetURI = CrlDiagramRefinedPointerURI + "/" + "Target"

// CrlDiagramReferenceLinkURI identifies the Reference represented as a link in the diagram
var CrlDiagramReferenceLinkURI = CrlDiagramDomainURI + "/" + "ReferenceLink"

// CrlDiagramReferenceLinkModelReferenceURI identifies the reference to the model element represented by the link
var CrlDiagramReferenceLinkModelReferenceURI = CrlDiagramReferenceLinkURI + "/" + "ModelReference"

// CrlDiagramReferenceLinkDisplayLabelURI identifies the display label concept to be used when displaying the link
var CrlDiagramReferenceLinkDisplayLabelURI = CrlDiagramReferenceLinkURI + "/" + "DisplayLabel"

// CrlDiagramReferenceLinkAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the link
var CrlDiagramReferenceLinkAbstractionDisplayLabelURI = CrlDiagramReferenceLinkURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramReferenceLinkSourceURI identifies the concept that is the source of the link
var CrlDiagramReferenceLinkSourceURI = CrlDiagramReferenceLinkURI + "/" + "Source"

// CrlDiagramReferenceLinkTargetURI identifies the concept that is the target of the link
var CrlDiagramReferenceLinkTargetURI = CrlDiagramReferenceLinkURI + "/" + "Target"

// CrlDiagramRefinementLinkURI identifies the Refinement represented as a link in the diagram
var CrlDiagramRefinementLinkURI = CrlDiagramDomainURI + "/" + "RefinementLink"

// CrlDiagramRefinementLinkModelReferenceURI identifies the reference to the model element represented by the link
var CrlDiagramRefinementLinkModelReferenceURI = CrlDiagramRefinementLinkURI + "/" + "ModelReference"

// CrlDiagramRefinementLinkDisplayLabelURI identifies the display label concept to be used when displaying the link
var CrlDiagramRefinementLinkDisplayLabelURI = CrlDiagramRefinementLinkURI + "/" + "DisplayLabel"

// CrlDiagramRefinementLinkAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the link
var CrlDiagramRefinementLinkAbstractionDisplayLabelURI = CrlDiagramRefinementLinkURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramRefinementLinkSourceURI identifies the concept that is the source of the link
var CrlDiagramRefinementLinkSourceURI = CrlDiagramRefinementLinkURI + "/" + "Source"

// CrlDiagramRefinementLinkTargetURI identifies the concept that is the target of the link
var CrlDiagramRefinementLinkTargetURI = CrlDiagramRefinementLinkURI + "/" + "Target"

// GetDisplayLabel is a convenience function for getting the DisplayLabel value of a DiagramElement
func GetDisplayLabel(diagramElement core.Element, trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	displayLabelLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
	if displayLabelLiteral != nil {
		return displayLabelLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetAbstractionDisplayLabel is a convenience function for getting the AbstractionDisplayLabel value for a node
func GetAbstractionDisplayLabel(diagramElement core.Element, trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	abstractionDisplayLabelLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, trans)
	if abstractionDisplayLabelLiteral != nil {
		return abstractionDisplayLabelLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetLineColor is a convenience function for getting the LineColor value of a DiagramElement
func GetLineColor(diagramElement core.Element, trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	lineColorLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementLineColorURI, trans)
	if lineColorLiteral != nil {
		return lineColorLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetBGColor is a convenience function for getting the backgound color value of a DiagramElement
func GetBGColor(diagramElement core.Element, trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	BGColorLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementBGColorURI, trans)
	if BGColorLiteral != nil {
		return BGColorLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetFirstElementRepresentingConcept returns the first diagram element that represents the indicated concept
func GetFirstElementRepresentingConcept(diagram core.Element, concept core.Element, trans *core.Transaction) core.Element {
	if diagram == nil || !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConcept called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, trans) {
		if GetReferencedModelConcept(el, trans) == concept && !el.IsRefinementOfURI(CrlDiagramPointerURI, trans) {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptID returns the first diagram element that represents the indicated concept
func GetFirstElementRepresentingConceptID(diagram core.Element, conceptID string, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptID called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, trans) {
		if GetReferencedModelConcept(el, trans).GetConceptID(trans) == conceptID && !el.IsRefinementOfURI(CrlDiagramPointerURI, trans) {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptOwnerPointer returns the first diagram element that represents the indicated concept's OwnerPointer
func GetFirstElementRepresentingConceptOwnerPointer(diagram core.Element, concept core.Element, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptOwnerPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if GetReferencedModelConcept(el, trans) == concept {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDOwnerPointer returns the first diagram element that represents the indicated concept's OwnerPointer
func GetFirstElementRepresentingConceptIDOwnerPointer(diagram core.Element, conceptID string, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptIDOwnerPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if GetReferencedModelConcept(el, trans).GetConceptID(trans) == conceptID {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptElementPointer returns the first diagram element that represents the indicated concept's ElementPointer
func GetFirstElementRepresentingConceptElementPointer(diagram core.Element, concept core.Reference, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptElementPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if GetReferencedModelConcept(el, trans) == concept {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDElementPointer returns the first diagram element that represents the indicated concept's ElementPointer
func GetFirstElementRepresentingConceptIDElementPointer(diagram core.Element, conceptID string, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptIDElementPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if GetReferencedModelConcept(el, trans).GetConceptID(trans) == conceptID {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptAbstractPointer returns the first diagram element that represents the indicated concept's AbstractPointer
func GetFirstElementRepresentingConceptAbstractPointer(diagram core.Element, concept core.Refinement, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptAbstractPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramAbstractPointerURI, trans) {
		if GetReferencedModelConcept(el, trans) == concept {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDAbstractPointer returns the first diagram element that represents the indicated concept's AbstractPointer
func GetFirstElementRepresentingConceptIDAbstractPointer(diagram core.Element, conceptID string, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptIDAbstractPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramAbstractPointerURI, trans) {
		if GetReferencedModelConcept(el, trans).GetConceptID(trans) == conceptID {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptRefinedPointer returns the first diagram element that represents the indicated concept's RefinedPointer
func GetFirstElementRepresentingConceptRefinedPointer(diagram core.Element, concept core.Refinement, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptRefinedPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramRefinedPointerURI, trans) {
		if GetReferencedModelConcept(el, trans) == concept {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDRefinedPointer returns the first diagram element that represents the indicated concept's RefinedPointer
func GetFirstElementRepresentingConceptIDRefinedPointer(diagram core.Element, conceptID string, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetFirstElementRepresentingConceptIDRefinedPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramRefinedPointerURI, trans) {
		if GetReferencedModelConcept(el, trans).GetConceptID(trans) == conceptID {
			return el
		}
	}
	return nil
}

// GetLinkSource is a convenience function for getting the source concept of a link
func GetLinkSource(diagramLink core.Element, trans *core.Transaction) core.Element {
	if diagramLink == nil {
		return nil
	}
	sourceReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
	if sourceReference != nil {
		return sourceReference.GetReferencedConcept(trans)
	}
	return nil
}

// GetLinkSourceReferemce is a convenience function for getting the source reference of a link
func GetLinkSourceReference(diagramLink core.Element, trans *core.Transaction) core.Reference {
	if diagramLink == nil {
		return nil
	}
	return diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
}

// GetLinkTarget is a convenience function for getting the target concept of a link
func GetLinkTarget(diagramLink core.Element, trans *core.Transaction) core.Element {
	if diagramLink == nil {
		return nil
	}
	targetReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
	if targetReference != nil {
		return targetReference.GetReferencedConcept(trans)
	}
	return nil
}

// GetLinkTargetReference is a convenience function for getting the target reference of a link
func GetLinkTargetReference(diagramLink core.Element, trans *core.Transaction) core.Reference {
	if diagramLink == nil {
		return nil
	}
	return diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
}

// GetNodeHeight is a convenience function for getting the Height value of a node's position
func GetNodeHeight(diagramNode core.Element, trans *core.Transaction) float64 {
	if diagramNode == nil {
		return 0.0
	}
	heightLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeHeightURI, trans)
	if heightLiteral != nil {
		value := heightLiteral.GetLiteralValue(trans)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetNodeWidth is a convenience function for getting the Width value of a node's position
func GetNodeWidth(diagramNode core.Element, trans *core.Transaction) float64 {
	if diagramNode == nil {
		return 0.0
	}
	widthLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeWidthURI, trans)
	if widthLiteral != nil {
		value := widthLiteral.GetLiteralValue(trans)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetNodeX is a convenience function for getting the X value of a node's position
func GetNodeX(diagramNode core.Element, trans *core.Transaction) float64 {
	if diagramNode == nil {
		return 0.0
	}
	xLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeXURI, trans)
	if xLiteral != nil {
		value := xLiteral.GetLiteralValue(trans)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetNodeY is a convenience function for getting the X value of a node's position
func GetNodeY(diagramNode core.Element, trans *core.Transaction) float64 {
	if diagramNode == nil {
		return 0.0
	}
	yLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, trans)
	if yLiteral != nil {
		value := yLiteral.GetLiteralValue(trans)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetDisplayLabelYOffset is a convenience function for getting the Display Label's Y offset within the node
func GetDisplayLabelYOffset(diagramNode core.Element, trans *core.Transaction) float64 {
	if diagramNode == nil {
		return 0.0
	}
	yOffsetLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeDisplayLabelYOffsetURI, trans)
	if yOffsetLiteral != nil {
		value := yOffsetLiteral.GetLiteralValue(trans)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetOwnerPointer returns the ownerPointer for the concept if one exists
func GetOwnerPointer(diagram core.Element, concept core.Element, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetOwnerPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if GetReferencedModelConcept(el, trans) == concept {
			return el
		}
	}
	return nil
}

// GetElementPointer returns the elementPointer for the concept if one exists
func GetElementPointer(diagram core.Element, concept core.Element, trans *core.Transaction) core.Element {
	if !diagram.IsRefinementOfURI(CrlDiagramURI, trans) {
		log.Printf("GetElementPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if GetReferencedModelConcept(el, trans) == concept {
			return el
		}
	}
	return nil
}

// GetReferencedModelConcept is a function on a CrlDiagramNode that returns the model element represented by the
// diagram node
func GetReferencedModelConcept(diagramElement core.Element, trans *core.Transaction) core.Element {
	if diagramElement == nil {
		return nil
	}
	reference := diagramElement.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
	if reference != nil {
		return reference.GetReferencedConcept(trans)
	}
	return nil
}

func init() {
	var err error

	// Set up fonts and faces
	goRegularFont, err = truetype.Parse(goregular.TTF)
	if err != nil {
		log.Print(err.Error())
	}
	goBoldFont, err = truetype.Parse(gobold.TTF)
	if err != nil {
		log.Print(err.Error())
	}

	goItalicFont, err := truetype.Parse(goitalic.TTF)
	if err != nil {
		log.Print(err.Error())
	}

	options12Pt := truetype.Options{Size: 12.0}
	go12PtRegularFace = truetype.NewFace(goRegularFont, &options12Pt)
	go12PtBoldFace = truetype.NewFace(goBoldFont, &options12Pt)

	options10Pt := truetype.Options{Size: 10.0}
	go10PtRegularFace = truetype.NewFace(goRegularFont, &options10Pt)
	go10PtItalicFace = truetype.NewFace(goItalicFont, &options10Pt)
}

// Int26_6ToFloat converts a fixed point 26_6 integer to a floating point number
func Int26_6ToFloat(val fixed.Int26_6) float64 {
	return float64(val) / 64.0
}

// IsDiagram returns true if the supplied element is a CrlDiagram
func IsDiagram(el core.Element, trans *core.Transaction) bool {
	switch el.(type) {
	case core.Element:
		return el.IsRefinementOfURI(CrlDiagramURI, trans)
	}
	return false
}

// IsDiagramAbstractPointer returns true if the supplied element is a CrlDiagramAbstractPointer
func IsDiagramAbstractPointer(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramAbstractPointerURI, trans)
}

// IsDiagramElement returns true if the supplied element is a CrlDiagramElement
func IsDiagramElement(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementURI, trans)
}

// IsDiagramElementPointer returns true if the supplied element is a CrlDiagramElementPointer
func IsDiagramElementPointer(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementPointerURI, trans)
}

// IsDiagramLink returns true if the supplied element is a CrlDiagramLink
func IsDiagramLink(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramLinkURI, trans)
}

// IsDiagramNode returns true if the supplied element is a CrlDiagramNode
func IsDiagramNode(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramNodeURI, trans)
}

// IsDiagramOwnerPointer returns true if the supplied element is a CrlDiagramOwnerPointer
func IsDiagramOwnerPointer(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramOwnerPointerURI, trans)
}

// IsDiagramPointer returns true if the supplied element is a CrlDiagramPointer
func IsDiagramPointer(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramPointerURI, trans)
}

// IsDiagramRefinedPointer returns true if the supplied element is a CrlDiagramRefinedPointer
func IsDiagramRefinedPointer(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinedPointerURI, trans)
}

// IsDiagramReferenceLink returns true if the supplied element is a CrlDiagramReferenceLink
func IsDiagramReferenceLink(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramReferenceLinkURI, trans)
}

// IsDiagramRefinementLink returns true if the supplied element is a CrlDiagramRefinementLink
func IsDiagramRefinementLink(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinementLinkURI, trans)
}

// IsModelReference returns true if the supplied element is a ModelReference
func IsModelReference(el core.Element, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementModelReferenceURI, trans)
}

// IsDisplayLabel returns true if the supplied Literal is the DisplayLabel
func IsDisplayLabel(el core.Element, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramNodeDisplayLabelURI, trans)
}

// NewDiagram creates a new diagram
func NewDiagram(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramURI, trans)
}

// NewDiagramReferenceLink creates a new diagram link to represent a reference
func NewDiagramReferenceLink(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramReferenceLinkURI, trans)
}

// NewDiagramRefinementLink creates a new diagram link
func NewDiagramRefinementLink(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramRefinementLinkURI, trans)
}

// NewDiagramNode creates a new diagram node
func NewDiagramNode(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Element, error) {
	newNode, err := uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, trans)
	if err != nil {
		return nil, errors.Wrap(err, "Diagram.go NewDiagramNode failed")
	}
	SetLineColor(newNode, "#000000", trans)
	return newNode, nil
}

// NewDiagramOwnerPointer creates a new DiagramOwnerPointer
func NewDiagramOwnerPointer(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramOwnerPointerURI, trans)
}

// NewDiagramElementPointer creates a new DiagramElementPointer
func NewDiagramElementPointer(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramElementPointerURI, trans)
}

// NewDiagramAbstractPointer creates a new DiagramAbstractPointer
func NewDiagramAbstractPointer(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramAbstractPointerURI, trans)
}

// NewDiagramRefinedPointer creates a new DiagramRefinedPointer
func NewDiagramRefinedPointer(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramRefinedPointerURI, trans)
}

// SetAbstractionDisplayLabel is a function on a CrlDiagramNode that sets the abstraction display label of the diagram node
func SetAbstractionDisplayLabel(diagramElement core.Element, value string, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(value, trans)
	updateNodeSize(diagramElement, trans)
}

// SetDisplayLabel is a function on a CrlDiagramNode that sets the display label of the diagram element.
// If the diagram element is a pointer, the value is ignored and the label is set to the empty string
func SetDisplayLabel(diagramElement core.Element, value string, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
	if literal == nil {
		return
	}
	if IsDiagramPointer(diagramElement, trans) {
		literal.SetLiteralValue("", trans)
	} else {
		literal.SetLiteralValue(value, trans)
	}
	updateNodeSize(diagramElement, trans)
}

// SetLineColor is a function on a CrlDiagramElement that sets the line color for the diagram element.
func SetLineColor(diagramElement core.Element, value string, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementLineColorURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(value, trans)
}

// SetBGColor is a function on a CrlDiagramNode that sets the background color for the diagram element.
// If the diagram element is a pointer, the value is ignored and the label is set to the empty string
func SetBGColor(diagramElement core.Element, value string, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementBGColorURI, trans)
	if literal == nil {
		return
	}
	if IsDiagramPointer(diagramElement, trans) {
		literal.SetLiteralValue("", trans)
	} else {
		literal.SetLiteralValue(value, trans)
	}
}

// SetLinkSource is a convenience function for setting the source concept of a link
func SetLinkSource(diagramLink core.Element, source core.Element, trans *core.Transaction) {
	if diagramLink == nil {
		return
	}
	sourceReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
	if sourceReference != nil {
		sourceReference.SetReferencedConcept(source, core.NoAttribute, trans)
	}
}

// SetLinkTarget is a convenience function for setting the target concept of a link
func SetLinkTarget(diagramLink core.Element, target core.Element, trans *core.Transaction) {
	if diagramLink == nil {
		return
	}
	// attributeName := core.NoAttribute
	// if target.IsRefinementOfURI(CrlDiagramAbstractPointerURI, trans) {
	// 	attributeName = core.AbstractConceptID
	// } else if target.IsRefinementOfURI(CrlDiagramOwnerPointerURI, trans) {
	// 	attributeName = core.OwningConceptID
	// } else if target.IsRefinementOfURI(CrlDiagramRefinedPointerURI, trans) {
	// 	attributeName = core.RefinedConceptID
	// } else if target.IsRefinementOfURI(CrlDiagramElementPointerURI, trans) {
	// 	attributeName = core.ReferencedConceptID
	// }
	targetReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
	if targetReference != nil {
		targetReference.SetReferencedConcept(target, core.NoAttribute, trans)
	}
}

// SetNodeHeight is a function on a CrlDiagramNode that sets the height of the diagram node
func SetNodeHeight(diagramNode core.Element, value float64, trans *core.Transaction) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeHeightURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeWidth is a function on a CrlDiagramNode that sets the width of the diagram node
func SetNodeWidth(diagramNode core.Element, value float64, trans *core.Transaction) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeWidthURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeX is a function on a CrlDiagramNode that sets the x of the diagram node
func SetNodeX(diagramNode core.Element, value float64, trans *core.Transaction) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeXURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeY is a function on a CrlDiagramNode that sets the y of the diagram node
func SetNodeY(diagramNode core.Element, value float64, trans *core.Transaction) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeDisplayLabelYOffset is a function on a CrlDiagramNode that sets the y offset of the display label within the node
func SetNodeDisplayLabelYOffset(diagramNode core.Element, value float64, trans *core.Transaction) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeDisplayLabelYOffsetURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetReferencedModelConcept is a function on a CrlDiagramNode that sets the model element represented by the
// diagram node
func SetReferencedModelConcept(diagramElement core.Element, el core.Element, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	reference := diagramElement.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
	if reference == nil {
		return
	}
	reference.SetReferencedConcept(el, core.NoAttribute, trans)
	updateDiagramElementForModelElementChange(diagramElement, el, trans)
}

// BuildCrlDiagramDomain builds the CrlDiagram concept space and adds it to the uOfD
func BuildCrlDiagramDomain(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) core.Element {
	// CrlDiagramDomain
	crlDiagramDomain, _ := uOfD.NewElement(trans, CrlDiagramDomainURI)
	crlDiagramDomain.SetLabel("CrlDiagramDomain", trans)

	//
	// CrlDiagram
	//
	crlDiagram, _ := uOfD.NewElement(trans, CrlDiagramURI)
	crlDiagram.SetLabel("CrlDiagram", trans)
	crlDiagram.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramWidth, _ := uOfD.NewLiteral(trans, CrlDiagramWidthURI)
	crlDiagramWidth.SetLabel("Width", trans)
	crlDiagramWidth.SetOwningConcept(crlDiagram, trans)

	crlDiagramHeight, _ := uOfD.NewLiteral(trans, CrlDiagramHeightURI)
	crlDiagramHeight.SetLabel("Height", trans)
	crlDiagramHeight.SetOwningConcept(crlDiagram, trans)

	//
	// CrlDiagramElement
	//
	crlDiagramElement, _ := uOfD.NewElement(trans, CrlDiagramElementURI)
	crlDiagramElement.SetLabel("CrlDiagramElement", trans)
	crlDiagramElement.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramElementModelReference, _ := uOfD.NewReference(trans, CrlDiagramElementModelReferenceURI)
	crlDiagramElementModelReference.SetLabel("ModelReference", trans)
	crlDiagramElementModelReference.SetOwningConcept(crlDiagramElement, trans)

	crlDiagramElementDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramElementDisplayLabelURI)
	crlDiagramElementDisplayLabel.SetLabel("DisplayLabel", trans)
	crlDiagramElementDisplayLabel.SetOwningConcept(crlDiagramElement, trans)

	crlDiagramElementAbstractionDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramElementAbstractionDisplayLabelURI)
	crlDiagramElementAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", trans)
	crlDiagramElementAbstractionDisplayLabel.SetOwningConcept(crlDiagramElement, trans)

	crlDiagramElementLineColor, _ := uOfD.NewOwnedLiteral(crlDiagramElement, "LineColor", trans, CrlDiagramElementLineColorURI)
	crlDiagramElementBGColor, _ := uOfD.NewOwnedLiteral(crlDiagramElement, "BGColor", trans, CrlDiagramElementBGColorURI)

	//
	// CrlDiagramNode
	//
	crlDiagramNode, _ := uOfD.NewElement(trans, CrlDiagramNodeURI)
	crlDiagramNode.SetLabel("CrlDiagramNode", trans)
	crlDiagramNode.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramNodeRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramNodeRefinement.SetOwningConcept(crlDiagramNode, trans)
	crlDiagramNodeRefinement.SetAbstractConcept(crlDiagramElement, trans)
	crlDiagramNodeRefinement.SetRefinedConcept(crlDiagramNode, trans)

	crlDiagramNodeModelReference, _ := uOfD.NewReference(trans, CrlDiagramNodeModelReferenceURI)
	crlDiagramNodeModelReference.SetLabel("ModelReference", trans)
	crlDiagramNodeModelReference.SetOwningConcept(crlDiagramNode, trans)

	crlDiagramNodeModelReferenceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramNodeModelReferenceRefinement.SetOwningConcept(crlDiagramNodeModelReference, trans)
	crlDiagramNodeModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, trans)
	crlDiagramNodeModelReferenceRefinement.SetRefinedConcept(crlDiagramNodeModelReference, trans)

	crlDiagramNodeDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramNodeDisplayLabelURI)
	crlDiagramNodeDisplayLabel.SetLabel("DisplayLabel", trans)
	crlDiagramNodeDisplayLabel.SetOwningConcept(crlDiagramNode, trans)

	crlDiagramNodeDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramNodeDisplayLabelRefinement.SetOwningConcept(crlDiagramNodeDisplayLabel, trans)
	crlDiagramNodeDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, trans)
	crlDiagramNodeDisplayLabelRefinement.SetRefinedConcept(crlDiagramNodeDisplayLabel, trans)

	crlDiagramNodeAbstractionDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramNodeAbstractionDisplayLabelURI)
	crlDiagramNodeAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", trans)
	crlDiagramNodeAbstractionDisplayLabel.SetOwningConcept(crlDiagramNode, trans)

	crlDiagramNodeAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramNodeAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramNodeAbstractionDisplayLabel, trans)
	crlDiagramNodeAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, trans)
	crlDiagramNodeAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramNodeAbstractionDisplayLabel, trans)

	crlDiagramNodeLineColor, _ := uOfD.NewOwnedLiteral(crlDiagramNode, "LineColor", trans, CrlDiagramNodeLineColorURI)
	crlDiagramNodeLineColor.SetLiteralValue("#000000", trans)
	uOfD.NewCompleteRefinement(crlDiagramElementLineColor, crlDiagramNodeLineColor, "LineColorRefinement", trans)

	crlDiagramNodeBGColor, _ := uOfD.NewOwnedLiteral(crlDiagramNode, "BGColor", trans, CrlDiagramNodeBGColorURI)
	uOfD.NewCompleteRefinement(crlDiagramElementBGColor, crlDiagramNodeBGColor, "BGColorRefinement", trans)

	crlDiagramNodeX, _ := uOfD.NewLiteral(trans, CrlDiagramNodeXURI)
	crlDiagramNodeX.SetLabel("X", trans)
	crlDiagramNodeX.SetOwningConcept(crlDiagramNode, trans)

	crlDiagramNodeY, _ := uOfD.NewLiteral(trans, CrlDiagramNodeYURI)
	crlDiagramNodeY.SetLabel("Y", trans)
	crlDiagramNodeY.SetOwningConcept(crlDiagramNode, trans)

	crlDiagramNodeHeight, _ := uOfD.NewLiteral(trans, CrlDiagramNodeHeightURI)
	crlDiagramNodeHeight.SetLabel("Height", trans)
	crlDiagramNodeHeight.SetOwningConcept(crlDiagramNode, trans)

	crlDiagramNodeWidth, _ := uOfD.NewLiteral(trans, CrlDiagramNodeWidthURI)
	crlDiagramNodeWidth.SetLabel("Width", trans)
	crlDiagramNodeWidth.SetOwningConcept(crlDiagramNode, trans)

	crlDiagramNodeDisplayLabelYOffset, _ := uOfD.NewLiteral(trans, CrlDiagramNodeDisplayLabelYOffsetURI)
	crlDiagramNodeDisplayLabelYOffset.SetLabel("DisplayLabelYOffset", trans)
	crlDiagramNodeDisplayLabelYOffset.SetOwningConcept(crlDiagramNode, trans)

	//
	// CrlDiagramLink
	//
	crlDiagramLink, _ := uOfD.NewElement(trans, CrlDiagramLinkURI)
	crlDiagramLink.SetLabel("CrlDiagramLink", trans)
	crlDiagramLink.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramLinkRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramLinkRefinement.SetOwningConcept(crlDiagramLink, trans)
	crlDiagramLinkRefinement.SetAbstractConcept(crlDiagramElement, trans)
	crlDiagramLinkRefinement.SetRefinedConcept(crlDiagramLink, trans)

	crlDiagramLinkSource, _ := uOfD.NewReference(trans, CrlDiagramLinkSourceURI)
	crlDiagramLinkSource.SetLabel("Source", trans)
	crlDiagramLinkSource.SetOwningConcept(crlDiagramLink, trans)

	crlDiagramLinkTarget, _ := uOfD.NewReference(trans, CrlDiagramLinkTargetURI)
	crlDiagramLinkTarget.SetLabel("Target", trans)
	crlDiagramLinkTarget.SetOwningConcept(crlDiagramLink, trans)

	//
	// Pointer
	//
	crlDiagramPointer, _ := uOfD.NewElement(trans, CrlDiagramPointerURI)
	crlDiagramPointer.SetLabel("Pointer", trans)
	crlDiagramPointer.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramPointerRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramPointerRefinement.SetOwningConcept(crlDiagramPointer, trans)
	crlDiagramPointerRefinement.SetAbstractConcept(crlDiagramLink, trans)
	crlDiagramPointerRefinement.SetRefinedConcept(crlDiagramPointer, trans)

	//
	// AbstractPointer
	//
	crlDiagramAbstractPointer, _ := uOfD.NewElement(trans, CrlDiagramAbstractPointerURI)
	crlDiagramAbstractPointer.SetLabel("AbstractPointer", trans)
	crlDiagramAbstractPointer.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramAbstractPointerRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramAbstractPointerRefinement.SetOwningConcept(crlDiagramAbstractPointer, trans)
	crlDiagramAbstractPointerRefinement.SetAbstractConcept(crlDiagramPointer, trans)
	crlDiagramAbstractPointerRefinement.SetRefinedConcept(crlDiagramAbstractPointer, trans)

	crlDiagramAbstractPointerModelReference, _ := uOfD.NewReference(trans, CrlDiagramAbstractPointerModelReferenceURI)
	crlDiagramAbstractPointerModelReference.SetLabel("ModelReference", trans)
	crlDiagramAbstractPointerModelReference.SetOwningConcept(crlDiagramAbstractPointer, trans)

	crlDiagramAbstractPointerModelReferenceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramAbstractPointerModelReferenceRefinement.SetOwningConcept(crlDiagramAbstractPointerModelReference, trans)
	crlDiagramAbstractPointerModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, trans)
	crlDiagramAbstractPointerModelReferenceRefinement.SetRefinedConcept(crlDiagramAbstractPointerModelReference, trans)

	crlDiagramAbstractPointerDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramAbstractPointerDisplayLabelURI)
	crlDiagramAbstractPointerDisplayLabel.SetLabel("DisplayLabel", trans)
	crlDiagramAbstractPointerDisplayLabel.SetOwningConcept(crlDiagramAbstractPointer, trans)

	crlDiagramAbstractPointerDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramAbstractPointerDisplayLabelRefinement.SetOwningConcept(crlDiagramAbstractPointerDisplayLabel, trans)
	crlDiagramAbstractPointerDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, trans)
	crlDiagramAbstractPointerDisplayLabelRefinement.SetRefinedConcept(crlDiagramAbstractPointerDisplayLabel, trans)

	crlDiagramAbstractPointerAbstractionDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramAbstractPointerAbstractionDisplayLabelURI)
	crlDiagramAbstractPointerAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", trans)
	crlDiagramAbstractPointerAbstractionDisplayLabel.SetOwningConcept(crlDiagramAbstractPointer, trans)

	crlDiagramAbstractPointerAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramAbstractPointerAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramAbstractPointerAbstractionDisplayLabel, trans)
	crlDiagramAbstractPointerAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, trans)
	crlDiagramAbstractPointerAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramAbstractPointerAbstractionDisplayLabel, trans)

	crlDiagramAbstractPointerSource, _ := uOfD.NewReference(trans, CrlDiagramAbstractPointerSourceURI)
	crlDiagramAbstractPointerSource.SetLabel("Source", trans)
	crlDiagramAbstractPointerSource.SetOwningConcept(crlDiagramAbstractPointer, trans)

	crlDiagramAbstractPointerSourceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramAbstractPointerSourceRefinement.SetOwningConcept(crlDiagramAbstractPointerSource, trans)
	crlDiagramAbstractPointerSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, trans)
	crlDiagramAbstractPointerSourceRefinement.SetRefinedConcept(crlDiagramAbstractPointerSource, trans)

	crlDiagramAbstractPointerTarget, _ := uOfD.NewReference(trans, CrlDiagramAbstractPointerTargetURI)
	crlDiagramAbstractPointerTarget.SetLabel("Target", trans)
	crlDiagramAbstractPointerTarget.SetOwningConcept(crlDiagramAbstractPointer, trans)

	crlDiagramAbstractPointerTargetRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramAbstractPointerTargetRefinement.SetOwningConcept(crlDiagramAbstractPointerTarget, trans)
	crlDiagramAbstractPointerTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, trans)
	crlDiagramAbstractPointerTargetRefinement.SetRefinedConcept(crlDiagramAbstractPointerTarget, trans)

	//
	// ElementPointer
	//
	crlDiagramElementPointer, _ := uOfD.NewElement(trans, CrlDiagramElementPointerURI)
	crlDiagramElementPointer.SetLabel("ElementPointer", trans)
	crlDiagramElementPointer.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramElementPointerRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramElementPointerRefinement.SetOwningConcept(crlDiagramElementPointer, trans)
	crlDiagramElementPointerRefinement.SetAbstractConcept(crlDiagramPointer, trans)
	crlDiagramElementPointerRefinement.SetRefinedConcept(crlDiagramElementPointer, trans)

	crlDiagramElementPointerModelReference, _ := uOfD.NewReference(trans, CrlDiagramElementPointerModelReferenceURI)
	crlDiagramElementPointerModelReference.SetLabel("ModelReference", trans)
	crlDiagramElementPointerModelReference.SetOwningConcept(crlDiagramElementPointer, trans)

	crlDiagramElementPointerModelReferenceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramElementPointerModelReferenceRefinement.SetOwningConcept(crlDiagramElementPointerModelReference, trans)
	crlDiagramElementPointerModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, trans)
	crlDiagramElementPointerModelReferenceRefinement.SetRefinedConcept(crlDiagramElementPointerModelReference, trans)

	crlDiagramElementPointerDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramElementPointerDisplayLabelURI)
	crlDiagramElementPointerDisplayLabel.SetLabel("DisplayLabel", trans)
	crlDiagramElementPointerDisplayLabel.SetOwningConcept(crlDiagramElementPointer, trans)

	crlDiagramElementPointerDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramElementPointerDisplayLabelRefinement.SetOwningConcept(crlDiagramElementPointerDisplayLabel, trans)
	crlDiagramElementPointerDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, trans)
	crlDiagramElementPointerDisplayLabelRefinement.SetRefinedConcept(crlDiagramElementPointerDisplayLabel, trans)

	crlDiagramElementPointerAbstractionDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramElementPointerAbstractionDisplayLabelURI)
	crlDiagramElementPointerAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", trans)
	crlDiagramElementPointerAbstractionDisplayLabel.SetOwningConcept(crlDiagramElementPointer, trans)

	crlDiagramElementPointerAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramElementPointerAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramElementPointerAbstractionDisplayLabel, trans)
	crlDiagramElementPointerAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, trans)
	crlDiagramElementPointerAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramElementPointerAbstractionDisplayLabel, trans)

	crlDiagramElementPointerSource, _ := uOfD.NewReference(trans, CrlDiagramElementPointerSourceURI)
	crlDiagramElementPointerSource.SetLabel("Source", trans)
	crlDiagramElementPointerSource.SetOwningConcept(crlDiagramElementPointer, trans)

	crlDiagramElementPointerSourceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramElementPointerSourceRefinement.SetOwningConcept(crlDiagramElementPointerSource, trans)
	crlDiagramElementPointerSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, trans)
	crlDiagramElementPointerSourceRefinement.SetRefinedConcept(crlDiagramElementPointerSource, trans)

	crlDiagramElementPointerTarget, _ := uOfD.NewReference(trans, CrlDiagramElementPointerTargetURI)
	crlDiagramElementPointerTarget.SetLabel("Target", trans)
	crlDiagramElementPointerTarget.SetOwningConcept(crlDiagramElementPointer, trans)

	crlDiagramElementPointerTargetRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramElementPointerTargetRefinement.SetOwningConcept(crlDiagramElementPointerTarget, trans)
	crlDiagramElementPointerTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, trans)
	crlDiagramElementPointerTargetRefinement.SetRefinedConcept(crlDiagramElementPointerTarget, trans)

	//
	// OwnerPointer
	//
	crlDiagramOwnerPointer, _ := uOfD.NewElement(trans, CrlDiagramOwnerPointerURI)
	crlDiagramOwnerPointer.SetLabel("OwnerPointer", trans)
	crlDiagramOwnerPointer.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramOwnerPointerRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramOwnerPointerRefinement.SetOwningConcept(crlDiagramOwnerPointer, trans)
	crlDiagramOwnerPointerRefinement.SetAbstractConcept(crlDiagramPointer, trans)
	crlDiagramOwnerPointerRefinement.SetRefinedConcept(crlDiagramOwnerPointer, trans)

	crlDiagramOwnerPointerModelReference, _ := uOfD.NewReference(trans, CrlDiagramOwnerPointerModelReferenceURI)
	crlDiagramOwnerPointerModelReference.SetLabel("ModelReference", trans)
	crlDiagramOwnerPointerModelReference.SetOwningConcept(crlDiagramOwnerPointer, trans)

	crlDiagramOwnerPointerModelReferenceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramOwnerPointerModelReferenceRefinement.SetOwningConcept(crlDiagramOwnerPointerModelReference, trans)
	crlDiagramOwnerPointerModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, trans)
	crlDiagramOwnerPointerModelReferenceRefinement.SetRefinedConcept(crlDiagramOwnerPointerModelReference, trans)

	crlDiagramOwnerPointerDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramOwnerPointerDisplayLabelURI)
	crlDiagramOwnerPointerDisplayLabel.SetLabel("DisplayLabel", trans)
	crlDiagramOwnerPointerDisplayLabel.SetOwningConcept(crlDiagramOwnerPointer, trans)

	crlDiagramOwnerPointerDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramOwnerPointerDisplayLabelRefinement.SetOwningConcept(crlDiagramOwnerPointerDisplayLabel, trans)
	crlDiagramOwnerPointerDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, trans)
	crlDiagramOwnerPointerDisplayLabelRefinement.SetRefinedConcept(crlDiagramOwnerPointerDisplayLabel, trans)

	crlDiagramOwnerPointerAbstractionDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramOwnerPointerAbstractionDisplayLabelURI)
	crlDiagramOwnerPointerAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", trans)
	crlDiagramOwnerPointerAbstractionDisplayLabel.SetOwningConcept(crlDiagramOwnerPointer, trans)

	crlDiagramOwnerPointerAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramOwnerPointerAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramOwnerPointerAbstractionDisplayLabel, trans)
	crlDiagramOwnerPointerAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, trans)
	crlDiagramOwnerPointerAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramOwnerPointerAbstractionDisplayLabel, trans)

	crlDiagramOwnerPointerSource, _ := uOfD.NewReference(trans, CrlDiagramOwnerPointerSourceURI)
	crlDiagramOwnerPointerSource.SetLabel("Source", trans)
	crlDiagramOwnerPointerSource.SetOwningConcept(crlDiagramOwnerPointer, trans)

	crlDiagramOwnerPointerSourceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramOwnerPointerSourceRefinement.SetOwningConcept(crlDiagramOwnerPointerSource, trans)
	crlDiagramOwnerPointerSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, trans)
	crlDiagramOwnerPointerSourceRefinement.SetRefinedConcept(crlDiagramOwnerPointerSource, trans)

	crlDiagramOwnerPointerTarget, _ := uOfD.NewReference(trans, CrlDiagramOwnerPointerTargetURI)
	crlDiagramOwnerPointerTarget.SetLabel("Target", trans)
	crlDiagramOwnerPointerTarget.SetOwningConcept(crlDiagramOwnerPointer, trans)

	crlDiagramOwnerPointerTargetRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramOwnerPointerTargetRefinement.SetOwningConcept(crlDiagramOwnerPointerTarget, trans)
	crlDiagramOwnerPointerTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, trans)
	crlDiagramOwnerPointerTargetRefinement.SetRefinedConcept(crlDiagramOwnerPointerTarget, trans)

	//
	// RefinedPointer
	//
	crlDiagramRefinedPointer, _ := uOfD.NewElement(trans, CrlDiagramRefinedPointerURI)
	crlDiagramRefinedPointer.SetLabel("RefinedPointer", trans)
	crlDiagramRefinedPointer.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramRefinedPointerRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinedPointerRefinement.SetOwningConcept(crlDiagramRefinedPointer, trans)
	crlDiagramRefinedPointerRefinement.SetAbstractConcept(crlDiagramPointer, trans)
	crlDiagramRefinedPointerRefinement.SetRefinedConcept(crlDiagramRefinedPointer, trans)

	crlDiagramRefinedPointerModelReference, _ := uOfD.NewReference(trans, CrlDiagramRefinedPointerModelReferenceURI)
	crlDiagramRefinedPointerModelReference.SetLabel("ModelReference", trans)
	crlDiagramRefinedPointerModelReference.SetOwningConcept(crlDiagramRefinedPointer, trans)

	crlDiagramRefinedPointerModelReferenceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinedPointerModelReferenceRefinement.SetOwningConcept(crlDiagramRefinedPointerModelReference, trans)
	crlDiagramRefinedPointerModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, trans)
	crlDiagramRefinedPointerModelReferenceRefinement.SetRefinedConcept(crlDiagramRefinedPointerModelReference, trans)

	crlDiagramRefinedPointerDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramRefinedPointerDisplayLabelURI)
	crlDiagramRefinedPointerDisplayLabel.SetLabel("DisplayLabel", trans)
	crlDiagramRefinedPointerDisplayLabel.SetOwningConcept(crlDiagramRefinedPointer, trans)

	crlDiagramRefinedPointerDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinedPointerDisplayLabelRefinement.SetOwningConcept(crlDiagramRefinedPointerDisplayLabel, trans)
	crlDiagramRefinedPointerDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, trans)
	crlDiagramRefinedPointerDisplayLabelRefinement.SetRefinedConcept(crlDiagramRefinedPointerDisplayLabel, trans)

	crlDiagramRefinedPointerAbstractionDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramRefinedPointerAbstractionDisplayLabelURI)
	crlDiagramRefinedPointerAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", trans)
	crlDiagramRefinedPointerAbstractionDisplayLabel.SetOwningConcept(crlDiagramRefinedPointer, trans)

	crlDiagramRefinedPointerAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinedPointerAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramRefinedPointerAbstractionDisplayLabel, trans)
	crlDiagramRefinedPointerAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, trans)
	crlDiagramRefinedPointerAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramRefinedPointerAbstractionDisplayLabel, trans)

	crlDiagramRefinedPointerSource, _ := uOfD.NewReference(trans, CrlDiagramRefinedPointerSourceURI)
	crlDiagramRefinedPointerSource.SetLabel("Source", trans)
	crlDiagramRefinedPointerSource.SetOwningConcept(crlDiagramRefinedPointer, trans)

	crlDiagramRefinedPointerSourceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinedPointerSourceRefinement.SetOwningConcept(crlDiagramRefinedPointerSource, trans)
	crlDiagramRefinedPointerSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, trans)
	crlDiagramRefinedPointerSourceRefinement.SetRefinedConcept(crlDiagramRefinedPointerSource, trans)

	crlDiagramRefinedPointerTarget, _ := uOfD.NewReference(trans, CrlDiagramRefinedPointerTargetURI)
	crlDiagramRefinedPointerTarget.SetLabel("Target", trans)
	crlDiagramRefinedPointerTarget.SetOwningConcept(crlDiagramRefinedPointer, trans)

	crlDiagramRefinedPointerTargetRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinedPointerTargetRefinement.SetOwningConcept(crlDiagramRefinedPointerTarget, trans)
	crlDiagramRefinedPointerTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, trans)
	crlDiagramRefinedPointerTargetRefinement.SetRefinedConcept(crlDiagramRefinedPointerTarget, trans)

	//
	// ReferenceLink
	//
	crlDiagramReferenceLink, _ := uOfD.NewElement(trans, CrlDiagramReferenceLinkURI)
	crlDiagramReferenceLink.SetLabel("ReferenceLink", trans)
	crlDiagramReferenceLink.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramReferenceLinkRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramReferenceLinkRefinement.SetOwningConcept(crlDiagramReferenceLink, trans)
	crlDiagramReferenceLinkRefinement.SetAbstractConcept(crlDiagramLink, trans)
	crlDiagramReferenceLinkRefinement.SetRefinedConcept(crlDiagramReferenceLink, trans)

	crlDiagramReferenceLinkModelReference, _ := uOfD.NewReference(trans, CrlDiagramReferenceLinkModelReferenceURI)
	crlDiagramReferenceLinkModelReference.SetLabel("ModelReference", trans)
	crlDiagramReferenceLinkModelReference.SetOwningConcept(crlDiagramReferenceLink, trans)

	crlDiagramReferenceLinkModelReferenceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramReferenceLinkModelReferenceRefinement.SetOwningConcept(crlDiagramReferenceLinkModelReference, trans)
	crlDiagramReferenceLinkModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, trans)
	crlDiagramReferenceLinkModelReferenceRefinement.SetRefinedConcept(crlDiagramReferenceLinkModelReference, trans)

	crlDiagramReferenceLinkDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramReferenceLinkDisplayLabelURI)
	crlDiagramReferenceLinkDisplayLabel.SetLabel("DisplayLabel", trans)
	crlDiagramReferenceLinkDisplayLabel.SetOwningConcept(crlDiagramReferenceLink, trans)

	crlDiagramReferenceLinkDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramReferenceLinkDisplayLabelRefinement.SetOwningConcept(crlDiagramReferenceLinkDisplayLabel, trans)
	crlDiagramReferenceLinkDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, trans)
	crlDiagramReferenceLinkDisplayLabelRefinement.SetRefinedConcept(crlDiagramReferenceLinkDisplayLabel, trans)

	crlDiagramReferenceLinkAbstractionDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramReferenceLinkAbstractionDisplayLabelURI)
	crlDiagramReferenceLinkAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", trans)
	crlDiagramReferenceLinkAbstractionDisplayLabel.SetOwningConcept(crlDiagramReferenceLink, trans)

	crlDiagramReferenceLinkAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramReferenceLinkAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramReferenceLinkAbstractionDisplayLabel, trans)
	crlDiagramReferenceLinkAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, trans)
	crlDiagramReferenceLinkAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramReferenceLinkAbstractionDisplayLabel, trans)

	crlDiagramReferenceLinkSource, _ := uOfD.NewReference(trans, CrlDiagramReferenceLinkSourceURI)
	crlDiagramReferenceLinkSource.SetLabel("Source", trans)
	crlDiagramReferenceLinkSource.SetOwningConcept(crlDiagramReferenceLink, trans)

	crlDiagramReferenceLinkSourceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramReferenceLinkSourceRefinement.SetOwningConcept(crlDiagramReferenceLinkSource, trans)
	crlDiagramReferenceLinkSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, trans)
	crlDiagramReferenceLinkSourceRefinement.SetRefinedConcept(crlDiagramReferenceLinkSource, trans)

	crlDiagramReferenceLinkTarget, _ := uOfD.NewReference(trans, CrlDiagramReferenceLinkTargetURI)
	crlDiagramReferenceLinkTarget.SetLabel("Target", trans)
	crlDiagramReferenceLinkTarget.SetOwningConcept(crlDiagramReferenceLink, trans)

	crlDiagramReferenceLinkTargetRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramReferenceLinkTargetRefinement.SetOwningConcept(crlDiagramReferenceLinkTarget, trans)
	crlDiagramReferenceLinkTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, trans)
	crlDiagramReferenceLinkTargetRefinement.SetRefinedConcept(crlDiagramReferenceLinkTarget, trans)

	//
	// RefinementLink
	//
	crlDiagramRefinementLink, _ := uOfD.NewElement(trans, CrlDiagramRefinementLinkURI)
	crlDiagramRefinementLink.SetLabel("RefinementLink", trans)
	crlDiagramRefinementLink.SetOwningConcept(crlDiagramDomain, trans)

	crlDiagramRefinementLinkRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinementLinkRefinement.SetOwningConcept(crlDiagramRefinementLink, trans)
	crlDiagramRefinementLinkRefinement.SetAbstractConcept(crlDiagramLink, trans)
	crlDiagramRefinementLinkRefinement.SetRefinedConcept(crlDiagramRefinementLink, trans)

	crlDiagramRefinementLinkModelReference, _ := uOfD.NewReference(trans, CrlDiagramRefinementLinkModelReferenceURI)
	crlDiagramRefinementLinkModelReference.SetLabel("ModelReference", trans)
	crlDiagramRefinementLinkModelReference.SetOwningConcept(crlDiagramRefinementLink, trans)

	crlDiagramRefinementLinkModelReferenceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinementLinkModelReferenceRefinement.SetOwningConcept(crlDiagramRefinementLinkModelReference, trans)
	crlDiagramRefinementLinkModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, trans)
	crlDiagramRefinementLinkModelReferenceRefinement.SetRefinedConcept(crlDiagramRefinementLinkModelReference, trans)

	crlDiagramRefinementLinkDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramRefinementLinkDisplayLabelURI)
	crlDiagramRefinementLinkDisplayLabel.SetLabel("DisplayLabel", trans)
	crlDiagramRefinementLinkDisplayLabel.SetOwningConcept(crlDiagramRefinementLink, trans)

	crlDiagramRefinementLinkDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinementLinkDisplayLabelRefinement.SetOwningConcept(crlDiagramRefinementLinkDisplayLabel, trans)
	crlDiagramRefinementLinkDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, trans)
	crlDiagramRefinementLinkDisplayLabelRefinement.SetRefinedConcept(crlDiagramRefinementLinkDisplayLabel, trans)

	crlDiagramRefinementLinkAbstractionDisplayLabel, _ := uOfD.NewLiteral(trans, CrlDiagramRefinementLinkAbstractionDisplayLabelURI)
	crlDiagramRefinementLinkAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", trans)
	crlDiagramRefinementLinkAbstractionDisplayLabel.SetOwningConcept(crlDiagramRefinementLink, trans)

	crlDiagramRefinementLinkAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinementLinkAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramRefinementLinkAbstractionDisplayLabel, trans)
	crlDiagramRefinementLinkAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, trans)
	crlDiagramRefinementLinkAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramRefinementLinkAbstractionDisplayLabel, trans)

	crlDiagramRefinementLinkSource, _ := uOfD.NewReference(trans, CrlDiagramRefinementLinkSourceURI)
	crlDiagramRefinementLinkSource.SetLabel("Source", trans)
	crlDiagramRefinementLinkSource.SetOwningConcept(crlDiagramRefinementLink, trans)

	crlDiagramRefinementLinkSourceRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinementLinkSourceRefinement.SetOwningConcept(crlDiagramRefinementLinkSource, trans)
	crlDiagramRefinementLinkSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, trans)
	crlDiagramRefinementLinkSourceRefinement.SetRefinedConcept(crlDiagramRefinementLinkSource, trans)

	crlDiagramRefinementLinkTarget, _ := uOfD.NewReference(trans, CrlDiagramRefinementLinkTargetURI)
	crlDiagramRefinementLinkTarget.SetLabel("Target", trans)
	crlDiagramRefinementLinkTarget.SetOwningConcept(crlDiagramRefinementLink, trans)

	crlDiagramRefinementLinkTargetRefinement, _ := uOfD.NewRefinement(trans)
	crlDiagramRefinementLinkTargetRefinement.SetOwningConcept(crlDiagramRefinementLinkTarget, trans)
	crlDiagramRefinementLinkTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, trans)
	crlDiagramRefinementLinkTargetRefinement.SetRefinedConcept(crlDiagramRefinementLinkTarget, trans)

	uOfD.AddFunction(CrlDiagramElementURI, updateDiagramElement)
	uOfD.AddFunction(CrlDiagramOwnerPointerURI, updateDiagramOwnerPointer)

	crlDiagramDomain.SetIsCoreRecursively(trans)
	return crlDiagramDomain
}

// updateDiagramElement updates the diagram element
func updateDiagramElement(diagramElement core.Element, notification *core.ChangeNotification, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	trans.WriteLockElement(diagramElement)
	// core Elements should always be ignored
	if diagramElement.GetIsCore(trans) {
		return nil
	}
	diagram := diagramElement.GetOwningConcept(trans)
	if diagram == nil {
		// There is nothing to do
		return nil
	}
	// Suppress circular notifications
	underlyingChange := notification.GetUnderlyingChange()
	if underlyingChange != nil && underlyingChange.IsReferenced(diagramElement) {
		return nil
	}

	// There are several notifications of interest here:
	//   - the deletion of the referenced model element
	//   - the label of the referenced model element
	//   - the list of immediate abstractions of the referenced model element.
	// First, determine whether it is the referenced model element that has changed
	diagramElementModelReference := diagramElement.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
	if diagramElementModelReference == nil {
		// Without a model reference, there is nothing to do. This scenario can occur during diagramElement deletion.
		return nil
	}
	modelElement := GetReferencedModelConcept(diagramElement, trans)
	switch notification.GetNatureOfChange() {
	case core.OwnedConceptChanged:
		switch underlyingChange.GetNatureOfChange() {
		case core.ConceptChanged:
			if underlyingChange.GetReportingElementID() == diagramElementModelReference.GetConceptID(trans) {
				// The underlying change is from the model reference
				updateDiagramElementForModelElementChange(diagramElement, modelElement, trans)
			}
		case core.ReferencedConceptChanged:
			underlyingReportingElementID := underlyingChange.GetReportingElementID()
			if underlyingReportingElementID == diagramElementModelReference.GetConceptID(trans) {
				// The underlying change is from the model reference
				if IsDiagramNode(diagramElement, trans) {
					currentModelElement := underlyingChange.GetAfterConceptState()
					previousModelElement := underlyingChange.GetBeforeConceptState()
					if currentModelElement != nil && previousModelElement != nil {
						if currentModelElement.ReferencedConceptID == "" && previousModelElement.ReferencedConceptID != "" {
							uOfD.DeleteElement(diagramElement, trans)
						} else {
							updateDiagramElementForModelElementChange(diagramElement, modelElement, trans)
						}
					}
				} else if IsDiagramLink(diagramElement, trans) {
					// First see if the underlying model element has been deleted
					currentModelElement := underlyingChange.GetAfterConceptState()
					previousModelElement := underlyingChange.GetBeforeConceptState()
					currentReferencedModelElementID := ""
					if currentModelElement != nil {
						currentReferencedModelElementID = currentModelElement.ReferencedConceptID
					}
					previousReferencedModelElementID := ""
					if previousModelElement != nil {
						previousReferencedModelElementID = previousModelElement.ReferencedConceptID
					}
					if currentModelElement != nil && previousModelElement != nil {
						if currentReferencedModelElementID == "" && previousReferencedModelElementID != "" {
							// If the underlying model element has been deleted then delete the diagram element
							uOfD.DeleteElement(diagramElement, trans)
						} else {
							// Otherwise we update the diagram element
							previousReferencedModelElement := uOfD.GetElement(previousReferencedModelElementID)
							switch typedModelElement := modelElement.(type) {
							case core.Reference:
								if IsDiagramElementPointer(diagramElement, trans) {
									currentReferencedModelElement := typedModelElement.GetReferencedConcept(trans)
									if previousReferencedModelElement != currentReferencedModelElement {
										if currentReferencedModelElement == nil {
											uOfD.DeleteElement(diagramElement, trans)
										} else {
											newTargetDiagramElement := GetFirstElementRepresentingConcept(diagram, currentReferencedModelElement, trans)
											SetLinkTarget(diagramElement, newTargetDiagramElement, trans)
										}
									}
								} else if IsDiagramReferenceLink(diagramElement, trans) {
									updateDiagramElementForModelElementChange(diagramElement, typedModelElement, trans)
									SetDisplayLabel(diagramElement, typedModelElement.GetLabel(trans), trans)
									newModelTarget := typedModelElement.GetReferencedConcept(trans)
									newModelSource := typedModelElement.GetOwningConcept(trans)
									if newModelSource == nil || newModelTarget == nil {
										uOfD.DeleteElement(diagramElement, trans)
										return nil
									}
									currentDiagramSource := GetLinkSource(diagramElement, trans)
									currentModelSource := GetReferencedModelConcept(currentDiagramSource, trans)
									currentDiagramTarget := GetLinkTarget(diagramElement, trans)
									currentModelTarget := GetReferencedModelConcept(currentDiagramTarget, trans)
									if currentModelSource != newModelSource {
										newDiagramSource := GetFirstElementRepresentingConcept(diagram, newModelSource, trans)
										if newDiagramSource == nil {
											uOfD.DeleteElement(diagramElement, trans)
											return nil
										}
										SetLinkSource(diagramElement, newDiagramSource, trans)
									}
									if currentModelTarget != newModelTarget {
										newDiagramTarget := GetFirstElementRepresentingConcept(diagram, newModelTarget, trans)
										if newDiagramTarget == nil {
											uOfD.DeleteElement(diagramElement, trans)
											return nil
										}
										SetLinkTarget(diagramElement, newDiagramTarget, trans)
									}
								}
							case core.Refinement:
								refinement := modelElement.(core.Refinement)
								if IsDiagramPointer(diagramElement, trans) {
									var newTargetModelElement core.Element
									if IsDiagramAbstractPointer(diagramElement, trans) {
										newTargetModelElement = refinement.GetAbstractConcept(trans)
									} else if IsDiagramRefinedPointer(diagramElement, trans) {
										newTargetModelElement = refinement.GetRefinedConcept(trans)
									} else if IsDiagramOwnerPointer(diagramElement, trans) {
										newTargetModelElement = refinement.GetOwningConcept(trans)
									}
									if previousReferencedModelElement != newTargetModelElement {
										if newTargetModelElement == nil {
											uOfD.DeleteElement(diagramElement, trans)
										} else {
											newTargetDiagramElement := GetFirstElementRepresentingConcept(diagram, newTargetModelElement, trans)
											SetLinkTarget(diagramElement, newTargetDiagramElement, trans)
										}
									}
								} else if IsDiagramRefinementLink(diagramElement, trans) {
									updateDiagramElementForModelElementChange(diagramElement, modelElement, trans)
									SetDisplayLabel(diagramElement, refinement.GetLabel(trans), trans)
									newModelTarget := refinement.GetAbstractConcept(trans)
									newModelSource := refinement.GetRefinedConcept(trans)
									if newModelTarget == nil || newModelSource == nil {
										uOfD.DeleteElement(diagramElement, trans)
										return nil
									}
									currentDiagramTarget := GetLinkTarget(diagramElement, trans)
									currentModelTarget := GetReferencedModelConcept(currentDiagramTarget, trans)
									currentDiagramSource := GetLinkSource(diagramElement, trans)
									currentModelSource := GetReferencedModelConcept(currentDiagramSource, trans)
									if currentModelTarget != newModelTarget {
										newDiagramTarget := GetFirstElementRepresentingConcept(diagram, newModelTarget, trans)
										if newDiagramTarget == nil {
											uOfD.DeleteElement(diagramElement, trans)
											return nil
										}
										SetLinkTarget(diagramElement, newDiagramTarget, trans)
									}
									if currentModelSource != newModelSource {
										newDiagramSource := GetFirstElementRepresentingConcept(diagram, newModelSource, trans)
										if newDiagramSource == nil {
											uOfD.DeleteElement(diagramElement, trans)
											return nil
										}
										SetLinkSource(diagramElement, newDiagramSource, trans)
									}
								}
							}
						}
					}
				}
			} else {
				// If this is a diagram link and the underlying reporting element is either its source reference or its target reference
				// and the referenced element is now nil, we need to delete the link
				if IsDiagramLink(diagramElement, trans) {
					// If this is a diagram link and the underlying reporting element is either its source reference or its target reference
					// and the referenced element is now nil, we need to delete the link
					underlyingReportingElementIsSourceReference := GetLinkSourceReference(diagramElement, trans).GetConceptID(trans) == underlyingReportingElementID
					underlyingReportingElementIsTargetReference := GetLinkTargetReference(diagramElement, trans).GetConceptID(trans) == underlyingReportingElementID
					if (underlyingReportingElementIsSourceReference || underlyingReportingElementIsTargetReference) &&
						underlyingChange.GetAfterConceptState().ReferencedConceptID == "" {
						uOfD.DeleteElement(diagramElement, trans)
					} else {
						switch underlyingChange.GetNatureOfChange() {
						case core.ReferencedConceptChanged:
							if underlyingReportingElementIsTargetReference {
								// If the link's target has changed, we need to update the underlying model element to reflect the change.
								// Note that if the target is now null, the preceeding clause will have deleted the element
								targetDiagramElement := uOfD.GetElement(underlyingChange.GetAfterConceptState().ReferencedConceptID)
								targetModelElement := GetReferencedModelConcept(targetDiagramElement, trans)
								if IsDiagramOwnerPointer(diagramElement, trans) {
									modelElement.SetOwningConcept(targetModelElement, trans)
								}
								switch typedModelElement := modelElement.(type) {
								case core.Reference:
									// Setting the referenced concepts requires knowledge of what is being referenced
									targetAttribute := core.NoAttribute
									if IsDiagramElementPointer(targetDiagramElement, trans) {
										targetAttribute = core.ReferencedConceptID
									} else if IsDiagramOwnerPointer(targetDiagramElement, trans) {
										targetAttribute = core.OwningConceptID
									} else if IsDiagramAbstractPointer(targetDiagramElement, trans) {
										targetAttribute = core.AbstractConceptID
									} else if IsDiagramRefinedPointer(targetDiagramElement, trans) {
										targetAttribute = core.RefinedConceptID
									}
									err := typedModelElement.SetReferencedConcept(targetModelElement, targetAttribute, trans)
									if err != nil {
										return errors.Wrap(err, "updateDiagramElement failed")
									}
								case core.Refinement:
									if underlyingReportingElementIsTargetReference {
										err := typedModelElement.SetAbstractConcept(targetModelElement, trans)
										if err != nil {
											return errors.Wrap(err, "updateDiagramElement failed")
										}
									} else if underlyingReportingElementIsSourceReference {
										err := typedModelElement.SetRefinedConcept(targetModelElement, trans)
										if err != nil {
											return errors.Wrap(err, "updateDiagramElement failed")
										}
									}
								}
							} else if underlyingReportingElementIsSourceReference {
								sourceDiagramElement := uOfD.GetElement(underlyingChange.GetAfterConceptState().ReferencedConceptID)
								sourceModelElement := GetReferencedModelConcept(sourceDiagramElement, trans)
								switch typedModelElement := modelElement.(type) {
								case core.Reference:
									if IsDiagramReferenceLink(diagramElement, trans) {
										typedModelElement.SetOwningConcept(sourceModelElement, trans)
									}
								case core.Refinement:
									if IsDiagramRefinementLink(diagramElement, trans) {
										typedModelElement.SetRefinedConcept(sourceModelElement, trans)
									}
								}
							}
						}
					}
				}
			}
		case core.IndicatedConceptChanged:
			// If the reporting element of the underlying change is the model element reference, then
			// see if the label needs to change
			if underlyingChange.GetReportingElementID() == diagramElementModelReference.GetConceptID(trans) {
				updateDiagramElementForModelElementChange(diagramElement, modelElement, trans)
			}
		}
	case core.ReferencedConceptChanged:
		// We are looking for the model diagramElementModelReference reporting a ConceptChanged which would be the result of setting the referencedConcept
		if notification.GetAfterConceptState().ConceptID != diagramElementModelReference.GetConceptID(trans) {
			break
		}
		if diagramElementModelReference.GetReferencedConceptID(trans) == "" {
			uOfD.DeleteElement(diagramElement, trans)
		} else {
			updateDiagramElementForModelElementChange(diagramElement, modelElement, trans)
		}
	}
	return nil
}

// updateDiagramOwnerPointer updates the ownerPointer's target if the ownership of the represented modelElement changes
func updateDiagramOwnerPointer(diagramPointer core.Element, notification *core.ChangeNotification, trans *core.Transaction) error {
	// There is one change of interest here: the model element's owner has changed
	uOfD := trans.GetUniverseOfDiscourse()
	trans.WriteLockElement(diagramPointer)
	reportingElement := uOfD.GetElement(notification.GetReportingElementID())
	diagram := diagramPointer.GetOwningConcept(trans)
	modelElement := GetReferencedModelConcept(diagramPointer, trans)
	switch notification.GetNatureOfChange() {
	case core.OwnedConceptChanged:
		if reportingElement == modelElement {
			underlyingNotification := notification.GetUnderlyingChange()
			switch underlyingNotification.GetNatureOfChange() {
			case core.OwningConceptChanged:
				if underlyingNotification.GetAfterConceptState().ConceptID == modelElement.GetConceptID(trans) {
					modelOwner := modelElement.GetOwningConcept(trans)
					var oldModelOwner core.Element
					diagramTarget := GetLinkTarget(diagramPointer, trans)
					if diagramTarget != nil {
						oldModelOwner = GetReferencedModelConcept(diagramTarget, trans)
					}
					if modelOwner != oldModelOwner {
						// Need to determine whether there is a view of the new owner in the diagram
						newDiagramTarget := GetFirstElementRepresentingConcept(diagram, modelOwner, trans)
						if newDiagramTarget == nil {
							// There is no view, delete the modelElement
							dEls := mapset.NewSet(diagramPointer.GetConceptID(trans))
							uOfD.DeleteElements(dEls, trans)
						} else {
							SetLinkTarget(diagramPointer, newDiagramTarget, trans)
						}
					}
				}
			}
			break
		}
		// We are looking for a notification from either the source or target reference in the diagram
		// If either source or target are nil, delete the pointer
		sourceReference := diagramPointer.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
		targetReference := diagramPointer.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
		if reportingElement == sourceReference || reportingElement == targetReference {
			underlyingNotification := notification.GetUnderlyingChange()
			switch underlyingNotification.GetNatureOfChange() {
			case core.ReferencedConceptChanged:
				switch typedElement := reportingElement.(type) {
				case core.Reference:
					if typedElement.GetReferencedConcept(trans) == nil {
						uOfD.DeleteElement(diagramPointer, trans)
					}
				}
			}
		}
	}
	return nil
}

func updateDiagramElementForModelElementChange(diagramElement core.Element, modelElement core.Element, trans *core.Transaction) {
	modelElementLabel := ""
	if modelElement != nil {
		modelElementLabel = modelElement.GetLabel(trans)
		if modelElementLabel != diagramElement.GetLabel(trans) {
			diagramElement.SetLabel(modelElementLabel, trans)
			if !IsDiagramPointer(diagramElement, trans) {
				SetDisplayLabel(diagramElement, modelElementLabel, trans)
			}
		}
		abstractions := make(map[string]core.Element)
		modelElement.FindImmediateAbstractions(abstractions, trans)
		abstractionsLabel := ""
		for _, abs := range abstractions {
			if len(abstractionsLabel) != 0 {
				abstractionsLabel += "\n"
			}
			abstractionsLabel += abs.GetLabel(trans)
		}
		if GetAbstractionDisplayLabel(diagramElement, trans) != abstractionsLabel {
			SetAbstractionDisplayLabel(diagramElement, abstractionsLabel, trans)
		}
	}
}

// updateNodeSize recalcualtes the size of the node based on the string sizes for the display label and
// abstractions listed
func updateNodeSize(node core.Element, trans *core.Transaction) {
	displayLabel := GetDisplayLabel(node, trans)
	displayLabelBounds, _ := font.BoundString(go12PtBoldFace, displayLabel)
	displayLabelMaxHeight := Int26_6ToFloat(displayLabelBounds.Max.Y)
	displayLabelMaxWidth := Int26_6ToFloat(displayLabelBounds.Max.X)
	displayLabelMinHeight := Int26_6ToFloat(displayLabelBounds.Min.Y)
	displayLabelMinWidth := Int26_6ToFloat(displayLabelBounds.Min.X)
	displayLabelHeight := displayLabelMaxHeight - displayLabelMinHeight
	displayLabelWidth := displayLabelMaxWidth - displayLabelMinWidth
	abstractionDisplayLabel := GetAbstractionDisplayLabel(node, trans)
	abstractionDisplayLabelBounds, _ := font.BoundString(go10PtItalicFace, abstractionDisplayLabel)
	abstractionDisplayLabelMaxWidth := Int26_6ToFloat(abstractionDisplayLabelBounds.Max.X)
	abstractionDisplayLabelMaxHeight := Int26_6ToFloat(abstractionDisplayLabelBounds.Max.Y)
	abstractionDisplayLabelMinWidth := Int26_6ToFloat(abstractionDisplayLabelBounds.Min.X)
	abstractionDisplayLabelMinHeight := Int26_6ToFloat(abstractionDisplayLabelBounds.Min.Y)
	abstractionDisplayLabelHeight := abstractionDisplayLabelMaxHeight - abstractionDisplayLabelMinHeight
	abstractionDisplayLabelWidth := abstractionDisplayLabelMaxWidth - abstractionDisplayLabelMinWidth
	topHeight := math.Max(IconSize, abstractionDisplayLabelHeight)
	height := topHeight + displayLabelHeight + 2*NodeLineWidth + 3*NodePadWidth
	topWidth := IconSize + 1*NodePadWidth + abstractionDisplayLabelWidth
	width := math.Max(topWidth, displayLabelWidth) + 2*NodeLineWidth + 2*NodePadWidth
	displayLabelYOffset := topHeight + NodeLineWidth + 2*NodePadWidth
	SetNodeHeight(node, height, trans)
	SetNodeWidth(node, width, trans)
	SetNodeDisplayLabelYOffset(node, displayLabelYOffset, trans)
}
