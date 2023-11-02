// Package crldiagramdomain defines the Diagram domain. This is a pre-defined domain that is, itself,
// represented as a CRLElement and identified with the CrlDiagramDomainURI. This concept space contains the prototypes of all Elements used to construct CrlDiagrams.
// Included are:
//
//		CrlDiagram: the diagram itself
//		CrlDiagramNode: a node in the diagram
//		CrlDiagramLink: a link in the diagram
//	    CrlDiagramPointer: a pointer shown as a link in the diagram
//		CrlDiagramAnchoredText: stand-alone text anchored to a reference point in the diagram
//
// These classes are intended to hold all of the information about the diagram that is not specific to the rendering engine.
//
// Intended Usage
// CRL Elements, in general, can have functions associated with them. When refinements of the elements are created, modified, or deleted, these functions are
// called. The strategy used for diagrams is to place all rendering-specific code in functions associated with the deining concepts.
// This is accomplished using the FunctionCallManager.AddFunctionCall() method. Note that this registration is NOT done in the core diagram package, but
// rather in the package providing the rendering engine linkage.
//
// Instances of the prototpes can be conveniently instantiated using the supplied New<type>() functions. This creates a refinement of the type and
// adds the appropriate children, with each child being a refinement of its defining type.
package crldiagramdomain

import (
	"log"
	"math"
	"strconv"

	"github.com/pkg/errors"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatatypesdomain"

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

// Diagram Anchored Text

// CrlDiagramAnchoredTextURI identifies the concept of an Anchored Text
var CrlDiagramAnchoredTextURI = CrlDiagramDomainURI + "/" + "AnchoredText"

// CrlDiagramAnchoredTextAnchorXURI identifies the x coordinate of the anchored text anchor point
var CrlDiagramAnchoredTextAnchorXURI = CrlDiagramAnchoredTextURI + "/" + "AnchorX"

// CrlDiagramAnchoredTextAnchorYURI identifies the y coordinate of the anchored text anchor point
var CrlDiagramAnchoredTextAnchorYURI = CrlDiagramAnchoredTextURI + "/" + "AnchorY"

// CrlDiagramAnchoredTextOffsetXURI identifies the x offset of the anchored text anchor point
var CrlDiagramAnchoredTextOffsetXURI = CrlDiagramAnchoredTextURI + "/" + "OffsetX"

// CrlDiagramAnchoredTextOffsetYURI identifies the y offset of the anchored text anchor point
var CrlDiagramAnchoredTextOffsetYURI = CrlDiagramAnchoredTextURI + "/" + "OffsetY"

// CrlDiagramAnchoredTextVisibleURI identifies the boolean indicating whether the anchored text is presently viewable
var CrlDiagramAnchoredTextVisibleURI = CrlDiagramAnchoredTextURI + "/" + "Visible"

// Diagram Element

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

// Diagram Node

// CrlDiagramNodeURI identifies the CrlDiagramNode concept
var CrlDiagramNodeURI = CrlDiagramDomainURI + "/" + "CrlDiagramNode"

// CrlDiagramNodeAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the element
var CrlDiagramNodeAbstractionDisplayLabelURI = CrlDiagramElementURI + "/" + "AbstractionDisplayLabel"

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

// Diagram Link

// CrlDiagramLinkURI identifies the CrlDiagramLink concept
var CrlDiagramLinkURI = CrlDiagramDomainURI + "/" + "CrlDiagramLink"

// CrlDiagramLinkSourceURI identifies the concept that is the source of the link
var CrlDiagramLinkSourceURI = CrlDiagramLinkURI + "/" + "Source"

// CrlDiagramLinkTargetURI identifies the concept that is the target of the link
var CrlDiagramLinkTargetURI = CrlDiagramLinkURI + "/" + "Target"

// CrlDiagramLinkDisplayLabelURI identifies the anchored text being used as the display label on a link
var CrlDiagramLinkDisplayLabelURI = CrlDiagramLinkURI + "/" + "DisplayLabel"

// CrlDiagramLinkAbstractionDisplayLabelURI identifies the anchored text used as the abstraction label on a link
var CrlDiagramLinkAbstractionDisplayLabelURI = CrlDiagramLinkURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramLinkMultiplicityURI identifies the anchored text used to display multiplicity
var CrlDiagramLinkMultiplicityURI = CrlDiagramLinkURI + "/" + "Multiplicity"

// Diagram Pointer

// CrlDiagramPointerURI identifies a pointer represented as a link
var CrlDiagramPointerURI = CrlDiagramDomainURI + "/" + "Pointer"

// CrlDiagramAbstractPointerURI identifies the Abstract of an Element represented as a link
var CrlDiagramAbstractPointerURI = CrlDiagramDomainURI + "/" + "AbstractPointer"

// CrlDiagramElementPointerURI identifies the element pointer of a Reference represented as a link
var CrlDiagramElementPointerURI = CrlDiagramDomainURI + "/" + "ElementPointer"

// CrlDiagramOwnerPointerURI identifies the owner of an Element represented as a link
var CrlDiagramOwnerPointerURI = CrlDiagramDomainURI + "/" + "OwnerPointer"

// CrlDiagramRefinedPointerURI identifies the refined element of a Refinement represented as a link
var CrlDiagramRefinedPointerURI = CrlDiagramDomainURI + "/" + "RefinedPointer"

// CrlDiagramReferenceLinkURI identifies the Reference represented as a link in the diagram
var CrlDiagramReferenceLinkURI = CrlDiagramDomainURI + "/" + "ReferenceLink"

// CrlDiagramRefinementLinkURI identifies the Refinement represented as a link in the diagram
var CrlDiagramRefinementLinkURI = CrlDiagramDomainURI + "/" + "RefinementLink"

func addAnchoredTextConcepts(anchoredText core.Concept, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextAnchorXURI, anchoredText, "AnchorX", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextAnchorYURI, anchoredText, "AnchorY", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextOffsetXURI, anchoredText, "OffsetX", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextOffsetYURI, anchoredText, "OffsetY", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextVisibleURI, anchoredText, "Visible", trans)
}

func addDiagramElementConcepts(newElement core.Concept, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementModelReferenceURI, newElement, "Model Reference", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementLineColorURI, newElement, "Line Color", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementBGColorURI, newElement, "BG Color", trans)
}

func addDiagramLinkConcepts(newLink core.Concept, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkSourceURI, newLink, "Source", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkTargetURI, newLink, "Target", trans)
	displayLabel, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkDisplayLabelURI, newLink, "DisplayLabel", trans)
	addAnchoredTextConcepts(displayLabel, trans)
	abstractionDisplayLabel, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkAbstractionDisplayLabelURI, newLink, "AbstractionDisplayLabel", trans)
	addAnchoredTextConcepts(abstractionDisplayLabel, trans)
}

func addDiagramNodeConcepts(newNode core.Concept, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeXURI, newNode, "X", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeYURI, newNode, "Y", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeHeightURI, newNode, "Height", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeWidthURI, newNode, "Width", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementDisplayLabelURI, newNode, "DisplayLabel", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementAbstractionDisplayLabelURI, newNode, "AbstractionDisplayLabel", trans)
}

// GetAbstractionDisplayLabel is a convenience function for getting the AbstractionDisplayLabel value for a DiagramElement
func GetAbstractionDisplayLabel(diagramElement core.Concept, trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	abstractionDisplayLabelLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, trans)
	if abstractionDisplayLabelLiteral != nil {
		return abstractionDisplayLabelLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetBGColor is a convenience function for getting the backgound color value of a DiagramElement
func GetBGColor(diagramElement core.Concept, trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	BGColorLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementBGColorURI, trans)
	if BGColorLiteral != nil {
		return BGColorLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetDisplayLabel is a convenience function for getting the DisplayLabel value of a DiagramElement
func GetDisplayLabel(diagramElement core.Concept, trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	displayLabelLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
	if displayLabelLiteral != nil {
		return displayLabelLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetFirstElementRepresentingConcept returns the first diagram element that represents the indicated concept
func GetFirstElementRepresentingConcept(diagram core.Concept, concept core.Concept, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptID(diagram core.Concept, conceptID string, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptOwnerPointer(diagram core.Concept, concept core.Concept, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptIDOwnerPointer(diagram core.Concept, conceptID string, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptElementPointer(diagram core.Concept, concept core.Concept, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptIDElementPointer(diagram core.Concept, conceptID string, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptAbstractPointer(diagram core.Concept, concept core.Concept, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptIDAbstractPointer(diagram core.Concept, conceptID string, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptRefinedPointer(diagram core.Concept, concept core.Concept, trans *core.Transaction) core.Concept {
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
func GetFirstElementRepresentingConceptIDRefinedPointer(diagram core.Concept, conceptID string, trans *core.Transaction) core.Concept {
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

// GetLineColor is a convenience function for getting the LineColor value of a DiagramElement
func GetLineColor(diagramElement core.Concept, trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	lineColorLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementLineColorURI, trans)
	if lineColorLiteral != nil {
		return lineColorLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetLinkSource is a convenience function for getting the source concept of a link
func GetLinkSource(diagramLink core.Concept, trans *core.Transaction) core.Concept {
	if diagramLink == nil {
		return nil
	}
	sourceReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
	if sourceReference != nil {
		return sourceReference.GetReferencedConcept(trans)
	}
	return nil
}

// GetLinkSourceReference is a convenience function for getting the source reference of a link
func GetLinkSourceReference(diagramLink core.Concept, trans *core.Transaction) core.Concept {
	if diagramLink == nil {
		return nil
	}
	return diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
}

// GetLinkTarget is a convenience function for getting the target concept of a link
func GetLinkTarget(diagramLink core.Concept, trans *core.Transaction) core.Concept {
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
func GetLinkTargetReference(diagramLink core.Concept, trans *core.Transaction) core.Concept {
	if diagramLink == nil {
		return nil
	}
	return diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
}

// GetNodeHeight is a convenience function for getting the Height value of a node's position
func GetNodeHeight(diagramNode core.Concept, trans *core.Transaction) float64 {
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
func GetNodeWidth(diagramNode core.Concept, trans *core.Transaction) float64 {
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
func GetNodeX(diagramNode core.Concept, trans *core.Transaction) float64 {
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
func GetNodeY(diagramNode core.Concept, trans *core.Transaction) float64 {
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

// GetOwnerPointer returns the ownerPointer for the concept if one exists
func GetOwnerPointer(diagram core.Concept, concept core.Concept, trans *core.Transaction) core.Concept {
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
func GetElementPointer(diagram core.Concept, concept core.Concept, trans *core.Transaction) core.Concept {
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
func GetReferencedModelConcept(diagramElement core.Concept, trans *core.Transaction) core.Concept {
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
func IsDiagram(el core.Concept, trans *core.Transaction) bool {
	switch el.(type) {
	case core.Concept:
		return el.IsRefinementOfURI(CrlDiagramURI, trans)
	}
	return false
}

// IsDiagramAbstractPointer returns true if the supplied element is a CrlDiagramAbstractPointer
func IsDiagramAbstractPointer(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramAbstractPointerURI, trans)
}

// IsDiagramElement returns true if the supplied element is a CrlDiagramElement
func IsDiagramElement(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementURI, trans)
}

// IsDiagramElementPointer returns true if the supplied element is a CrlDiagramElementPointer
func IsDiagramElementPointer(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementPointerURI, trans)
}

// IsDiagramLink returns true if the supplied element is a CrlDiagramLink
func IsDiagramLink(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramLinkURI, trans)
}

// IsDiagramNode returns true if the supplied element is a CrlDiagramNode
func IsDiagramNode(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramNodeURI, trans)
}

// IsDiagramOwnerPointer returns true if the supplied element is a CrlDiagramOwnerPointer
func IsDiagramOwnerPointer(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramOwnerPointerURI, trans)
}

// IsDiagramPointer returns true if the supplied element is a CrlDiagramPointer
func IsDiagramPointer(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramPointerURI, trans)
}

// IsDiagramRefinedPointer returns true if the supplied element is a CrlDiagramRefinedPointer
func IsDiagramRefinedPointer(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinedPointerURI, trans)
}

// IsDiagramReferenceLink returns true if the supplied element is a CrlDiagramReferenceLink
func IsDiagramReferenceLink(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramReferenceLinkURI, trans)
}

// IsDiagramRefinementLink returns true if the supplied element is a CrlDiagramRefinementLink
func IsDiagramRefinementLink(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinementLinkURI, trans)
}

// IsModelReference returns true if the supplied element is a ModelReference
func IsModelReference(el core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementModelReferenceURI, trans)
}

// IsDisplayLabel returns true if the supplied Literal is the DisplayLabel
func IsDisplayLabel(el core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
}

// NewDiagram creates a new diagram
func NewDiagram(trans *core.Transaction) (core.Concept, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newDiagram, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramURI, "New Diagram", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramWidthURI, newDiagram, "Width", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramHeightURI, newDiagram, "Height", trans)
	return newDiagram, nil
}

// NewDiagramNode creates a new diagram node
func NewDiagramNode(trans *core.Transaction) (core.Concept, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newNode, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramNodeURI, "New Node", trans)
	addDiagramElementConcepts(newNode, trans)
	addDiagramNodeConcepts(newNode, trans)
	SetLineColor(newNode, "#00000000", trans)
	return newNode, nil
}

// NewDiagramReferenceLink creates a new diagram link to represent a reference
func NewDiagramReferenceLink(trans *core.Transaction) (core.Concept, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newLink, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramReferenceLinkURI, "ReferenceLink", trans)
	addDiagramElementConcepts(newLink, trans)
	addDiagramLinkConcepts(newLink, trans)
	multiplicity, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkMultiplicityURI, newLink, "Multiplicity", trans)
	addAnchoredTextConcepts(multiplicity, trans)
	return newLink, nil
}

// NewDiagramRefinementLink creates a new diagram link
func NewDiagramRefinementLink(trans *core.Transaction) (core.Concept, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newLink, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramRefinementLinkURI, "RefinementLink", trans)
	addDiagramElementConcepts(newLink, trans)
	addDiagramLinkConcepts(newLink, trans)
	return newLink, nil
}

// NewDiagramOwnerPointer creates a new DiagramOwnerPointer
func NewDiagramOwnerPointer(trans *core.Transaction) (core.Concept, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newPointer, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramOwnerPointerURI, "OwnerPointer", trans)
	addDiagramElementConcepts(newPointer, trans)
	addDiagramLinkConcepts(newPointer, trans)
	multiplicity, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkMultiplicityURI, newPointer, "Multiplicity", trans)
	addAnchoredTextConcepts(multiplicity, trans)
	return newPointer, nil
}

// NewDiagramElementPointer creates a new DiagramElementPointer
func NewDiagramElementPointer(trans *core.Transaction) (core.Concept, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newPointer, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramElementPointerURI, "ElementPointer", trans)
	addDiagramElementConcepts(newPointer, trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// NewDiagramAbstractPointer creates a new DiagramAbstractPointer
func NewDiagramAbstractPointer(trans *core.Transaction) (core.Concept, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newPointer, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramAbstractPointerURI, "AbstractPointer", trans)
	addDiagramElementConcepts(newPointer, trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// NewDiagramRefinedPointer creates a new DiagramRefinedPointer
func NewDiagramRefinedPointer(trans *core.Transaction) (core.Concept, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newPointer, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramRefinedPointerURI, "RefinedPointer", trans)
	addDiagramElementConcepts(newPointer, trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// SetAbstractionDisplayLabel is a function on a CrlDiagramNode that sets the abstraction display label of the diagram node
func SetAbstractionDisplayLabel(diagramElement core.Concept, value string, trans *core.Transaction) {
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
func SetDisplayLabel(diagramElement core.Concept, value string, trans *core.Transaction) {
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
func SetLineColor(diagramElement core.Concept, value string, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementLineColorURI, trans)
	if literal == nil {
		// This is remedial code: the literal should already be there
		uOfD := trans.GetUniverseOfDiscourse()
		literal, _ = uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementLineColorURI, diagramElement, "Line Color", trans)
	}
	literal.SetLiteralValue(value, trans)
}

// SetBGColor is a function on a CrlDiagramNode that sets the background color for the diagram element.
// If the diagram element is a pointer, the value is ignored and the label is set to the empty string
func SetBGColor(diagramElement core.Concept, value string, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementBGColorURI, trans)
	if literal == nil {
		// This is remedial code: the literal should already be there
		uOfD := trans.GetUniverseOfDiscourse()
		literal, _ = uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementBGColorURI, diagramElement, "Background Color", trans)
	}
	if IsDiagramPointer(diagramElement, trans) {
		literal.SetLiteralValue("", trans)
	} else {
		literal.SetLiteralValue(value, trans)
	}
}

// SetLinkSource is a convenience function for setting the source concept of a link
func SetLinkSource(diagramLink core.Concept, source core.Concept, trans *core.Transaction) {
	if diagramLink == nil {
		return
	}
	sourceReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
	if sourceReference != nil {
		sourceReference.SetReferencedConcept(source, core.NoAttribute, trans)
	}
}

// SetLinkTarget is a convenience function for setting the target concept of a link
func SetLinkTarget(diagramLink core.Concept, target core.Concept, trans *core.Transaction) {
	if diagramLink == nil {
		return
	}
	targetReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
	if targetReference != nil {
		targetReference.SetReferencedConcept(target, core.NoAttribute, trans)
	}
}

// SetNodeHeight is a function on a CrlDiagramNode that sets the height of the diagram node
func SetNodeHeight(diagramNode core.Concept, value float64, trans *core.Transaction) {
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
func SetNodeWidth(diagramNode core.Concept, value float64, trans *core.Transaction) {
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
func SetNodeX(diagramNode core.Concept, value float64, trans *core.Transaction) {
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
func SetNodeY(diagramNode core.Concept, value float64, trans *core.Transaction) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetReferencedModelConcept is a function on a CrlDiagramNode that sets the model element represented by the
// diagram node
func SetReferencedModelConcept(diagramElement core.Concept, el core.Concept, trans *core.Transaction) {
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
func BuildCrlDiagramDomain(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) core.Concept {
	if uOfD.GetElementWithURI(crldatatypesdomain.CrlDataTypesDomainURI) == nil {
		crldatatypesdomain.BuildCrlDataTypesDomain(uOfD, trans)
	}
	// CrlDiagramDomain
	crlDiagramDomain, _ := uOfD.NewElement(trans, CrlDiagramDomainURI)
	crlDiagramDomain.SetLabel("CrlDiagramDomain", trans)

	//
	// CrlDiagram
	//
	crlDiagram, _ := uOfD.NewOwnedElement(crlDiagramDomain, "CrlDiagram", trans, CrlDiagramURI)
	uOfD.NewOwnedLiteral(crlDiagram, "Width", trans, CrlDiagramWidthURI)
	uOfD.NewOwnedLiteral(crlDiagram, "Height", trans, CrlDiagramHeightURI)

	// Diagram Anchored Text

	crlDiagramAnchoredText, _ := uOfD.NewOwnedLiteral(crlDiagramDomain, "AnchoredText", trans, CrlDiagramAnchoredTextURI)
	uOfD.NewOwnedLiteral(crlDiagramAnchoredText, "AnchorX", trans, CrlDiagramAnchoredTextAnchorXURI)
	uOfD.NewOwnedLiteral(crlDiagramAnchoredText, "AnchorY", trans, CrlDiagramAnchoredTextAnchorYURI)
	uOfD.NewOwnedLiteral(crlDiagramAnchoredText, "OffsetX", trans, CrlDiagramAnchoredTextOffsetXURI)
	uOfD.NewOwnedLiteral(crlDiagramAnchoredText, "OffsetY", trans, CrlDiagramAnchoredTextOffsetYURI)
	crldatatypesdomain.NewOwnedBoolean(crlDiagramAnchoredText, "Visible", trans, CrlDiagramAnchoredTextVisibleURI)

	// Multiplicity
	uOfD.CreateOwnedRefinementOfConcept(crlDiagramAnchoredText, crlDiagramDomain, "Multiplicity", trans, CrlDiagramLinkMultiplicityURI)

	//
	// CrlDiagramElement
	//
	crlDiagramElement, _ := uOfD.NewOwnedElement(crlDiagramDomain, "CrlDiagramElement", trans, CrlDiagramElementURI)
	uOfD.NewOwnedReference(crlDiagramElement, "ModelReference", trans, CrlDiagramElementModelReferenceURI)
	crlDisplayLabel, _ := uOfD.NewOwnedLiteral(crlDiagramElement, "DisplayLabel", trans, CrlDiagramElementDisplayLabelURI)
	crlAbstractionDisplayLabel, _ := uOfD.NewOwnedLiteral(crlDiagramElement, "AbstractionDisplayLabel", trans, CrlDiagramElementAbstractionDisplayLabelURI)
	uOfD.NewOwnedLiteral(crlDiagramElement, "LineColor", trans, CrlDiagramElementLineColorURI)
	uOfD.NewOwnedLiteral(crlDiagramElement, "BGColor", trans, CrlDiagramElementBGColorURI)

	//
	// CrlDiagramNode
	//
	crlDiagramNode, _ := uOfD.CreateOwnedRefinementOfConcept(crlDiagramElement, crlDiagramDomain, "CrlDiagramNode", trans, CrlDiagramNodeURI)
	uOfD.NewOwnedLiteral(crlDiagramNode, "X", trans, CrlDiagramNodeXURI)
	uOfD.NewOwnedLiteral(crlDiagramNode, "Y", trans, CrlDiagramNodeYURI)
	uOfD.NewOwnedLiteral(crlDiagramNode, "Height", trans, CrlDiagramNodeHeightURI)
	uOfD.NewOwnedLiteral(crlDiagramNode, "Width", trans, CrlDiagramNodeWidthURI)
	uOfD.NewOwnedLiteral(crlDiagramNode, "DisplayLabelYOffset", trans, CrlDiagramNodeDisplayLabelYOffsetURI)

	//
	// CrlDiagramLink
	//
	crlDiagramLink, _ := uOfD.CreateOwnedRefinementOfConcept(crlDiagramElement, crlDiagramDomain, "CrlDiagramLink", trans, CrlDiagramLinkURI)
	uOfD.NewOwnedReference(crlDiagramLink, "Source", trans, CrlDiagramLinkSourceURI)
	uOfD.NewOwnedReference(crlDiagramLink, "Target", trans, CrlDiagramLinkTargetURI)
	floatingDisplayLabel, _ := uOfD.CreateOwnedRefinementOfConcept(crlDiagramAnchoredText, crlDiagramLink, "DisplayLabel", trans, CrlDiagramLinkDisplayLabelURI)
	uOfD.AddAbstractionToConcept(floatingDisplayLabel, crlDisplayLabel, trans)
	floatingAbstractionDisplayLabel, _ := uOfD.CreateOwnedRefinementOfConcept(crlDiagramAnchoredText, crlDiagramLink, "DisplayLabel", trans, CrlDiagramLinkAbstractionDisplayLabelURI)
	uOfD.AddAbstractionToConcept(floatingAbstractionDisplayLabel, crlAbstractionDisplayLabel, trans)

	uOfD.CreateOwnedRefinementOfConcept(crlDiagramLink, crlDiagramDomain, "ReferenceLink", trans, CrlDiagramReferenceLinkURI)
	uOfD.CreateOwnedRefinementOfConcept(crlDiagramLink, crlDiagramDomain, "RefinementLink", trans, CrlDiagramRefinementLinkURI)

	//
	// Pointers
	//
	crlDiagramPointer, _ := uOfD.CreateOwnedRefinementOfConcept(crlDiagramLink, crlDiagramDomain, "Pointer", trans, CrlDiagramPointerURI)
	uOfD.CreateOwnedRefinementOfConcept(crlDiagramPointer, crlDiagramDomain, "AbstractPointer", trans, CrlDiagramAbstractPointerURI)
	uOfD.CreateOwnedRefinementOfConcept(crlDiagramPointer, crlDiagramDomain, "ElementPointer", trans, CrlDiagramElementPointerURI)
	uOfD.CreateOwnedRefinementOfConcept(crlDiagramPointer, crlDiagramDomain, "OwnerPointer", trans, CrlDiagramOwnerPointerURI)
	uOfD.CreateOwnedRefinementOfConcept(crlDiagramPointer, crlDiagramDomain, "RefinedPointer", trans, CrlDiagramRefinedPointerURI)

	uOfD.AddFunction(CrlDiagramElementURI, updateDiagramElement)
	uOfD.AddFunction(CrlDiagramOwnerPointerURI, updateDiagramOwnerPointer)

	crlDiagramDomain.SetIsCoreRecursively(trans)
	return crlDiagramDomain
}

// updateDiagramElement updates the diagram element
func updateDiagramElement(diagramElement core.Concept, notification *core.ChangeNotification, trans *core.Transaction) error {
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
							switch modelElement.GetConceptType() {
							case core.Reference:
								if IsDiagramElementPointer(diagramElement, trans) {
									currentReferencedModelElement := modelElement.GetReferencedConcept(trans)
									if previousReferencedModelElement != currentReferencedModelElement {
										if currentReferencedModelElement == nil {
											uOfD.DeleteElement(diagramElement, trans)
										} else {
											newTargetDiagramElement := GetFirstElementRepresentingConcept(diagram, currentReferencedModelElement, trans)
											SetLinkTarget(diagramElement, newTargetDiagramElement, trans)
										}
									}
								} else if IsDiagramReferenceLink(diagramElement, trans) {
									updateDiagramElementForModelElementChange(diagramElement, modelElement, trans)
									SetDisplayLabel(diagramElement, modelElement.GetLabel(trans), trans)
									newModelTarget := modelElement.GetReferencedConcept(trans)
									newModelSource := modelElement.GetOwningConcept(trans)
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
								refinement := modelElement
								if IsDiagramPointer(diagramElement, trans) {
									var newTargetModelElement core.Concept
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
								} else {
									switch modelElement.GetConceptType() {
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
										err := modelElement.SetReferencedConcept(targetModelElement, targetAttribute, trans)
										if err != nil {
											return errors.Wrap(err, "updateDiagramElement failed")
										}
									case core.Refinement:
										if underlyingReportingElementIsTargetReference {
											err := modelElement.SetAbstractConcept(targetModelElement, trans)
											if err != nil {
												return errors.Wrap(err, "updateDiagramElement failed")
											}
										} else if underlyingReportingElementIsSourceReference {
											err := modelElement.SetRefinedConcept(targetModelElement, trans)
											if err != nil {
												return errors.Wrap(err, "updateDiagramElement failed")
											}
										}
									}
								}
							} else if underlyingReportingElementIsSourceReference {
								sourceDiagramElement := uOfD.GetElement(underlyingChange.GetAfterConceptState().ReferencedConceptID)
								sourceModelElement := GetReferencedModelConcept(sourceDiagramElement, trans)
								if modelElement != nil {
									switch modelElement.GetConceptType() {
									case core.Reference:
										if IsDiagramReferenceLink(diagramElement, trans) {
											modelElement.SetOwningConcept(sourceModelElement, trans)
										}
									case core.Refinement:
										if IsDiagramRefinementLink(diagramElement, trans) {
											modelElement.SetRefinedConcept(sourceModelElement, trans)
										}
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
func updateDiagramOwnerPointer(diagramPointer core.Concept, notification *core.ChangeNotification, trans *core.Transaction) error {
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
					var oldModelOwner core.Concept
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
				switch reportingElement.GetConceptType() {
				case core.Reference:
					if reportingElement.GetReferencedConcept(trans) == nil {
						uOfD.DeleteElement(diagramPointer, trans)
					}
				}
			}
		}
	}
	return nil
}

func updateDiagramElementForModelElementChange(diagramElement core.Concept, modelElement core.Concept, trans *core.Transaction) {
	modelElementLabel := ""
	if modelElement != nil {
		modelElementLabel = modelElement.GetLabel(trans)
		if modelElementLabel != diagramElement.GetLabel(trans) {
			newLabel := modelElementLabel
			if IsDiagramPointer(diagramElement, trans) {
				if IsDiagramOwnerPointer(diagramElement, trans) {
					newLabel = newLabel + " Owner Pointer"
				} else if IsDiagramAbstractPointer(diagramElement, trans) {
					newLabel = newLabel + " Abstract Pointer"
				} else if IsDiagramRefinedPointer(diagramElement, trans) {
					newLabel = newLabel + " Refined Pointer"
				} else if IsDiagramElementPointer(diagramElement, trans) {
					newLabel = newLabel + " Referenced Concept Pointer"
				}
			}
			diagramElement.SetLabel(newLabel, trans)
			if !IsDiagramPointer(diagramElement, trans) {
				SetDisplayLabel(diagramElement, modelElementLabel, trans)
			}
		}
		abstractions := make(map[string]core.Concept)
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
func updateNodeSize(node core.Concept, trans *core.Transaction) {
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
	// displayLabelYOffset := topHeight + NodeLineWidth + 2*NodePadWidth
	SetNodeHeight(node, height, trans)
	SetNodeWidth(node, width, trans)
	// SetNodeDisplayLabelYOffset(node, displayLabelYOffset, trans)
}
