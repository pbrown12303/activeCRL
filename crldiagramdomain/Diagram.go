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
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/pkg/errors"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatatypesdomain"

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

// CrlDiagram is the CRL representation of a diagram
type CrlDiagram core.Concept

// CrlDiagramWidthURI identifies the CrlDiagramWidth concept
var CrlDiagramWidthURI = CrlDiagramURI + "/" + "Width"

// CrlDiagramHeightURI identifies the CrlDiagramHeight concept
var CrlDiagramHeightURI = CrlDiagramURI + "/" + "Height"

// Diagram Anchored Text

// CrlDiagramAnchoredTextURI identifies the concept of an Anchored Text
var CrlDiagramAnchoredTextURI = CrlDiagramDomainURI + "/" + "AnchoredText"

// CrlDiagramAnchoredText is the CRL representation of an anchored text
type CrlDiagramAnchoredText core.Concept

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

// CrlDiagramElement is the CRL representation of a diagram element
type CrlDiagramElement core.Concept

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

// CrlDiagramNode is the CRL representation of a diagram node
type CrlDiagramNode CrlDiagramElement

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

// CrlDiagramLink is the CRL representation of a diagram link
type CrlDiagramLink CrlDiagramElement

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

func addAnchoredTextConcepts(anchoredText *CrlDiagramAnchoredText, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextAnchorXURI, anchoredText.ToCore(), "AnchorX", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextAnchorYURI, anchoredText.ToCore(), "AnchorY", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextAnchorYURI, anchoredText.ToCore(), "AnchorY", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextOffsetXURI, anchoredText.ToCore(), "OffsetX", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextOffsetYURI, anchoredText.ToCore(), "OffsetY", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextVisibleURI, anchoredText.ToCore(), "Visible", trans)
}

func addDiagramElementConcepts(newElement *CrlDiagramElement, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementModelReferenceURI, newElement.ToCore(), "Model Reference", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementLineColorURI, newElement.ToCore(), "Line Color", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementBGColorURI, newElement.ToCore(), "BG Color", trans)
}

func addDiagramLinkConcepts(newLink *CrlDiagramLink, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkSourceURI, newLink.ToCore(), "Source", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkTargetURI, newLink.ToCore(), "Target", trans)
	displayLabel, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkDisplayLabelURI, newLink.ToCore(), "DisplayLabel", trans)
	uOfD.AddAbstractionURIToConcept(displayLabel, CrlDiagramAnchoredTextURI, trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(displayLabel), trans)
	abstractionDisplayLabel, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkAbstractionDisplayLabelURI, newLink.ToCore(), "AbstractionDisplayLabel", trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(abstractionDisplayLabel), trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(abstractionDisplayLabel), trans)
}

func addDiagramNodeConcepts(newNode *CrlDiagramNode, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeXURI, newNode.ToCore(), "X", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeYURI, newNode.ToCore(), "Y", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeHeightURI, newNode.ToCore(), "Height", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeWidthURI, newNode.ToCore(), "Width", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementDisplayLabelURI, newNode.ToCore(), "DisplayLabel", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementAbstractionDisplayLabelURI, newNode.ToCore(), "AbstractionDisplayLabel", trans)
}

// ToCore casts a diagram element to core.Concept
func (diagramElement *CrlDiagramElement) ToCore() *core.Concept {
	return (*core.Concept)(diagramElement)
}

// ToLink casts a diagram element to CrlDiagramLink after checking that it is the appopriate type
// It returns nil if it is the wrong type
func (diagramElement *CrlDiagramElement) ToLink(trans *core.Transaction) *CrlDiagramLink {
	if diagramElement.IsLink(trans) {
		return (*CrlDiagramLink)(diagramElement)
	}
	return nil
}

// ToNode casts a diagram element to CrlDiagramNode after checking that it is the appropriate type
// It returns nil if it is the wrong type
func (diagramElement *CrlDiagramElement) ToNode(trans *core.Transaction) *CrlDiagramNode {
	if diagramElement.IsNode(trans) {
		return (*CrlDiagramNode)(diagramElement)
	}
	return nil
}

// GetAbstractionDisplayLabel is a convenience function for getting the AbstractionDisplayLabel value for a DiagramElement
func (diagramElement *CrlDiagramElement) GetAbstractionDisplayLabel(trans *core.Transaction) string {
	abstractionDisplayLabelLiteral := diagramElement.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, trans)
	if abstractionDisplayLabelLiteral != nil {
		return abstractionDisplayLabelLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetAnchoredTextWithLabel returns the diagram element's first anchored text child with the given label
func (diagramElement *CrlDiagramElement) GetAnchoredTextWithLabel(label string, trans *core.Transaction) *CrlDiagramAnchoredText {
	if !diagramElement.ToCore().IsRefinementOfURI(CrlDiagramElementURI, trans) {
		log.Print("GetAnchoredTextWithLabel called for a concept that is not a CrlDiagramElement")
		return nil
	}
	anchoredTexts := diagramElement.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramAnchoredTextURI, trans)
	for _, anchoredText := range anchoredTexts {
		if anchoredText.GetLabel(trans) == label {
			return (*CrlDiagramAnchoredText)(anchoredText)
		}
	}
	return nil
}

// GetBGColor is a convenience function for getting the backgound color value of a DiagramElement
func (diagramElement *CrlDiagramElement) GetBGColor(trans *core.Transaction) string {
	BGColorLiteral := diagramElement.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementBGColorURI, trans)
	if BGColorLiteral != nil {
		return BGColorLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetDiagram returns the diagram within which the DiagramElement appears
func (diagramElement *CrlDiagramElement) GetDiagram(trans *core.Transaction) *CrlDiagram {
	return (*CrlDiagram)(diagramElement.ToCore().GetOwningConcept(trans))
}

// GetDisplayLabel is a convenience function for getting the DisplayLabel value of a DiagramElement
func (diagramElement *CrlDiagramElement) GetDisplayLabel(trans *core.Transaction) string {
	displayLabelLiteral := diagramElement.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
	if displayLabelLiteral != nil {
		return displayLabelLiteral.GetLiteralValue(trans)
	}
	return ""
}

// ToCore casts the diagram to core.Concept
func (diagram *CrlDiagram) ToCore() *core.Concept {
	return (*core.Concept)(diagram)
}

// GetFirstElementRepresentingConcept returns the first non-pointer diagram element that represents the indicated concept
func (diagram *CrlDiagram) GetFirstElementRepresentingConcept(concept *core.Concept, trans *core.Transaction) *CrlDiagramElement {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept && !el.IsRefinementOfURI(CrlDiagramPointerURI, trans) {
			return (*CrlDiagramElement)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptID returns the first diagram element that represents the indicated concept
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptID(conceptID string, trans *core.Transaction) *CrlDiagramElement {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID && !el.IsRefinementOfURI(CrlDiagramPointerURI, trans) {
			return (*CrlDiagramElement)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptOwnerPointer returns the first diagram element that represents the indicated concept's OwnerPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptOwnerPointer(concept *core.Concept, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDOwnerPointer returns the first diagram element that represents the indicated concept's OwnerPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptIDOwnerPointer(conceptID string, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptElementPointer returns the first diagram element that represents the indicated concept's ElementPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptElementPointer(concept *core.Concept, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDElementPointer returns the first diagram element that represents the indicated concept's ElementPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptIDElementPointer(conceptID string, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptAbstractPointer returns the first diagram element that represents the indicated concept's AbstractPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptAbstractPointer(concept *core.Concept, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramAbstractPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDAbstractPointer returns the first diagram element that represents the indicated concept's AbstractPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptIDAbstractPointer(conceptID string, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramAbstractPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptRefinedPointer returns the first diagram element that represents the indicated concept's RefinedPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptRefinedPointer(concept *core.Concept, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramRefinedPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDRefinedPointer returns the first diagram element that represents the indicated concept's RefinedPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptIDRefinedPointer(conceptID string, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramRefinedPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetLineColor is a convenience function for getting the LineColor value of a DiagramElement
func (diagramElement *CrlDiagramElement) GetLineColor(trans *core.Transaction) string {
	if diagramElement == nil {
		return ""
	}
	lineColorLiteral := diagramElement.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementLineColorURI, trans)
	if lineColorLiteral != nil {
		return lineColorLiteral.GetLiteralValue(trans)
	}
	return ""
}

// ToCore casts the CrlDiagramLink to core.Concept
func (diagramLink *CrlDiagramLink) ToCore() *core.Concept {
	return (*core.Concept)(diagramLink)
}

// GetLinkSource is a convenience function for getting the source concept of a link
func (diagramLink *CrlDiagramLink) GetLinkSource(trans *core.Transaction) *CrlDiagramElement {
	sourceReference := diagramLink.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
	if sourceReference != nil {
		referencedConcept := sourceReference.GetReferencedConcept(trans)
		if referencedConcept.IsRefinementOfURI(CrlDiagramElementURI, trans) {
			return (*CrlDiagramElement)(referencedConcept)
		}
	}
	return nil
}

// GetLinkSourceReference is a convenience function for getting the source reference of a link
func (diagramLink *CrlDiagramLink) GetLinkSourceReference(trans *core.Transaction) *core.Concept {
	return diagramLink.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
}

// GetLinkTarget is a convenience function for getting the target concept of a link
func (diagramLink *CrlDiagramLink) GetLinkTarget(trans *core.Transaction) *CrlDiagramElement {
	targetReference := diagramLink.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
	if targetReference != nil {
		referencedConcept := targetReference.GetReferencedConcept(trans)
		if referencedConcept.IsRefinementOfURI(CrlDiagramElementURI, trans) {
			return (*CrlDiagramElement)(referencedConcept)
		}
	}
	return nil
}

// GetLinkTargetReference is a convenience function for getting the target reference of a link
func (diagramLink *CrlDiagramLink) GetLinkTargetReference(trans *core.Transaction) *core.Concept {
	return diagramLink.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
}

// ToCore casts CrlDiagramNode to core.Concept
func (diagramNode *CrlDiagramNode) ToCore() *core.Concept {
	return (*core.Concept)(diagramNode)
}

// GetNodeHeight is a convenience function for getting the Height value of a node's position
func (diagramNode *CrlDiagramNode) GetNodeHeight(trans *core.Transaction) float64 {
	heightLiteral := diagramNode.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeHeightURI, trans)
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
func (diagramNode *CrlDiagramNode) GetNodeWidth(trans *core.Transaction) float64 {
	widthLiteral := diagramNode.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeWidthURI, trans)
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
func (diagramNode *CrlDiagramNode) GetNodeX(trans *core.Transaction) float64 {
	if diagramNode == nil {
		return 0.0
	}
	xLiteral := diagramNode.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeXURI, trans)
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
func (diagramNode *CrlDiagramNode) GetNodeY(trans *core.Transaction) float64 {
	yLiteral := diagramNode.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, trans)
	if yLiteral != nil {
		value := yLiteral.GetLiteralValue(trans)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// ToCore casts CrlDiagramAnchoredText to core.Concept
func (anchoredText *CrlDiagramAnchoredText) ToCore() *core.Concept {
	return (*core.Concept)(anchoredText)
}

// GetOffsetX returns the x offset value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) GetOffsetX(trans *core.Transaction) float64 {
	xOffsetLiteral := anchoredText.ToCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextOffsetXURI, trans)
	if xOffsetLiteral == nil {
		log.Printf("GetOffsetX called but no xOffsetLiteral was found")
		return 0
	}
	xOffset, err := strconv.ParseFloat(xOffsetLiteral.GetLiteralValue(trans), 64)
	if err != nil {
		errors.Wrap(err, "GetOffsetX failed")
		return 0
	}
	return xOffset
}

// GetOwnerPointer returns the ownerPointer for the concept if one exists
func (diagram *CrlDiagram) GetOwnerPointer(concept *CrlDiagramElement, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept.ToCore() {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetElementPointer returns the elementPointer for the concept if one exists
func (diagram *CrlDiagram) GetElementPointer(concept *CrlDiagramElement, trans *core.Transaction) *CrlDiagramLink {
	for _, el := range diagram.ToCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept.ToCore() {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetReferencedModelConcept is a function on a CrlDiagramElement that returns the model element represented by the
// diagram node
func (diagramElement *CrlDiagramElement) GetReferencedModelConcept(trans *core.Transaction) *core.Concept {
	reference := diagramElement.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
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

// GetCrlDiagram returns the CrlDiagram with the given ID
func GetCrlDiagram(id string, trans *core.Transaction) *CrlDiagram {
	concept := trans.GetUniverseOfDiscourse().GetElement(id)
	if concept.IsRefinementOfURI(CrlDiagramLinkURI, trans) {
		return (*CrlDiagram)(concept)
	}
	return nil
}

// GetCrlDiagramElement returns the CrlDiagramElement with the given ID
func GetCrlDiagramElement(id string, trans *core.Transaction) *CrlDiagramElement {
	concept := trans.GetUniverseOfDiscourse().GetElement(id)
	if concept.IsRefinementOfURI(CrlDiagramElementURI, trans) {
		return (*CrlDiagramElement)(concept)
	}
	return nil
}

// GetCrlDiagramLink returns the CrlDiagramLink with the given ID
func GetCrlDiagramLink(id string, trans *core.Transaction) *CrlDiagramLink {
	concept := trans.GetUniverseOfDiscourse().GetElement(id)
	if concept.IsRefinementOfURI(CrlDiagramLinkURI, trans) {
		return (*CrlDiagramLink)(concept)
	}
	return nil
}

// GetCrlDiagramNode returns the CrlDiagramNode with the given ID
func GetCrlDiagramNode(id string, trans *core.Transaction) *CrlDiagramNode {
	concept := trans.GetUniverseOfDiscourse().GetElement(id)
	if concept.IsRefinementOfURI(CrlDiagramNodeURI, trans) {
		return (*CrlDiagramNode)(concept)
	}
	return nil
}

// Int26_6ToFloat converts a fixed point 26_6 integer to a floating point number
func Int26_6ToFloat(val fixed.Int26_6) float64 {
	return float64(val) / 64.0
}

// IsDiagram returns true if the supplied element is a CrlDiagram
func IsDiagram(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramURI, trans)
}

// IsDiagramAbstractPointer returns true if the supplied element is a CrlDiagramAbstractPointer
func IsDiagramAbstractPointer(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramAbstractPointerURI, trans)
}

// IsDiagramElement returns true if the supplied element is a CrlDiagramElement
func IsDiagramElement(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementURI, trans)
}

// IsDiagramElementPointer returns true if the supplied element is a CrlDiagramElementPointer
func IsDiagramElementPointer(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementPointerURI, trans)
}

// IsDiagramLink returns true if the supplied element is a CrlDiagramLink
func IsDiagramLink(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramLinkURI, trans)
}

// IsDiagramNode returns true if the supplied element is a CrlDiagramNode
func IsDiagramNode(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramNodeURI, trans)
}

// IsDiagramOwnerPointer returns true if the supplied element is a CrlDiagramOwnerPointer
func IsDiagramOwnerPointer(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramOwnerPointerURI, trans)
}

// IsDiagramPointer returns true if the supplied element is a CrlDiagramPointer
func IsDiagramPointer(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramPointerURI, trans)
}

// IsDiagramRefinedPointer returns true if the supplied element is a CrlDiagramRefinedPointer
func IsDiagramRefinedPointer(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinedPointerURI, trans)
}

// IsDiagramReferenceLink returns true if the supplied element is a CrlDiagramReferenceLink
func IsDiagramReferenceLink(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramReferenceLinkURI, trans)
}

// IsDiagramRefinementLink returns true if the supplied element is a CrlDiagramRefinementLink
func IsDiagramRefinementLink(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinementLinkURI, trans)
}

// IsModelReference returns true if the supplied element is a ModelReference
func IsModelReference(el *core.Concept, trans *core.Transaction) bool {
	return el.IsRefinementOfURI(CrlDiagramElementModelReferenceURI, trans)
}

// IsDisplayLabel returns true if the supplied Literal is the DisplayLabel
func IsDisplayLabel(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
}

// NewDiagram creates a new diagram
func NewDiagram(trans *core.Transaction) (*CrlDiagram, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newElement, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramURI, "New Diagram", trans)
	newDiagram := (*CrlDiagram)(newElement)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramWidthURI, newDiagram.ToCore(), "Width", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramHeightURI, newDiagram.ToCore(), "Height", trans)
	return newDiagram, nil
}

// ToCrlDiagramElement casts a CrlDiagramNode to a CrlDiagramElement
func (diagramNode *CrlDiagramNode) ToCrlDiagramElement() *CrlDiagramElement {
	return (*CrlDiagramElement)(diagramNode)
}

// ToCrlDiagramElement casts a CrlDiagramLink to a CrlDiagramElement
func (diagramLink *CrlDiagramLink) ToCrlDiagramElement() *CrlDiagramElement {
	return (*CrlDiagramElement)(diagramLink)
}

// NewDiagramNode creates a new diagram node
func NewDiagramNode(trans *core.Transaction) (*CrlDiagramNode, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newElement, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramNodeURI, "New Node", trans)
	newNode := (*CrlDiagramNode)(newElement)
	addDiagramElementConcepts(newNode.ToCrlDiagramElement(), trans)
	addDiagramNodeConcepts(newNode, trans)
	newNode.ToCrlDiagramElement().SetLineColor("#00000000", trans)
	return newNode, nil
}

// NewDiagramReferenceLink creates a new diagram link to represent a reference
func NewDiagramReferenceLink(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newElement, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramReferenceLinkURI, "ReferenceLink", trans)
	newLink := (*CrlDiagramLink)(newElement)
	addDiagramElementConcepts(newLink.ToCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newLink, trans)
	multiplicity, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkMultiplicityURI, newLink.ToCore(), "Multiplicity", trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(multiplicity), trans)
	return newLink, nil
}

// NewDiagramRefinementLink creates a new diagram link representing a refinement
func NewDiagramRefinementLink(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramRefinementLinkURI, "RefinementLink", trans)
	newLink := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts(newLink.ToCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newLink, trans)
	return newLink, nil
}

// NewDiagramOwnerPointer creates a new DiagramOwnerPointer
func NewDiagramOwnerPointer(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramOwnerPointerURI, "OwnerPointer", trans)
	newPointer := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts(newPointer.ToCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newPointer, trans)
	multiplicity, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkMultiplicityURI, newPointer.ToCore(), "Multiplicity", trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(multiplicity), trans)
	return newPointer, nil
}

// NewDiagramElementPointer creates a new DiagramElementPointer
func NewDiagramElementPointer(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramElementPointerURI, "ElementPointer", trans)
	newPointer := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts((*CrlDiagramElement)(newPointer.ToCore()), trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// NewDiagramAbstractPointer creates a new DiagramAbstractPointer
func NewDiagramAbstractPointer(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramAbstractPointerURI, "AbstractPointer", trans)
	newPointer := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts(newPointer.ToCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// NewDiagramRefinedPointer creates a new DiagramRefinedPointer
func NewDiagramRefinedPointer(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramRefinedPointerURI, "RefinedPointer", trans)
	newPointer := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts(newPointer.ToCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// SetAbstractionDisplayLabel is a function on a CrlDiagramElement that sets the abstraction display label of the diagram node
func (diagramElement *CrlDiagramElement) SetAbstractionDisplayLabel(value string, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(value, trans)
	if IsDiagramNode(diagramElement.ToCore(), trans) {
		(*CrlDiagramNode)(diagramElement).updateNodeSize(trans)
	}
}

// SetDisplayLabel is a function on a CrlDiagramNode that sets the display label of the diagram element.
// If the diagram element is a pointer, the value is ignored and the label is set to the empty string
func (diagramElement *CrlDiagramElement) SetDisplayLabel(value string, trans *core.Transaction) {
	literal := diagramElement.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
	if literal == nil {
		return
	}
	if IsDiagramPointer(diagramElement.ToCore(), trans) {
		literal.SetLiteralValue("", trans)
	} else {
		literal.SetLiteralValue(value, trans)
	}
	if IsDiagramNode(diagramElement.ToCore(), trans) {
		(*CrlDiagramNode)(diagramElement).updateNodeSize(trans)
	}
}

// SetLineColor is a function on a CrlDiagramElement that sets the line color for the diagram element.
func (diagramElement *CrlDiagramElement) SetLineColor(value string, trans *core.Transaction) {
	literal := diagramElement.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementLineColorURI, trans)
	if literal == nil {
		// This is remedial code: the literal should already be there
		uOfD := trans.GetUniverseOfDiscourse()
		literal, _ = uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementLineColorURI, diagramElement.ToCore(), "Line Color", trans)
	}
	literal.SetLiteralValue(value, trans)
}

// SetBGColor is a function on a CrlDiagramNode that sets the background color for the diagram element.
// If the diagram element is a pointer, the value is ignored and the label is set to the empty string
func (diagramElement *CrlDiagramElement) SetBGColor(value string, trans *core.Transaction) {
	literal := diagramElement.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementBGColorURI, trans)
	if literal == nil {
		// This is remedial code: the literal should already be there
		uOfD := trans.GetUniverseOfDiscourse()
		literal, _ = uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementBGColorURI, diagramElement.ToCore(), "Background Color", trans)
	}
	if IsDiagramPointer(diagramElement.ToCore(), trans) {
		literal.SetLiteralValue("", trans)
	} else {
		literal.SetLiteralValue(value, trans)
	}
}

// SetDiagram sets the owner of the DiagramElement to the Diagram
func (diagramElement *CrlDiagramElement) SetDiagram(diagram *CrlDiagram, trans *core.Transaction) {
	diagramElement.ToCore().SetOwningConcept(diagram.ToCore(), trans)
}

// IsRefinementLink returns true if the link represents a refinement
func (diagramLink *CrlDiagramLink) IsRefinementLink(trans *core.Transaction) bool {
	return diagramLink.ToCore().IsRefinementOfURI(CrlDiagramRefinementLinkURI, trans)
}

// IsReferenceLink returns true if the link represents a reference
func (diagramLink *CrlDiagramLink) IsReferenceLink(trans *core.Transaction) bool {
	return diagramLink.ToCore().IsRefinementOfURI(CrlDiagramReferenceLinkURI, trans)
}

// IsOwnerPointer returns true if the link represents an owner pointer
func (diagramLink *CrlDiagramLink) IsOwnerPointer(trans *core.Transaction) bool {
	return diagramLink.ToCore().IsRefinementOfURI(CrlDiagramOwnerPointerURI, trans)
}

// IsElementPointer returns true if the link represents a referenced element pointer
func (diagramLink *CrlDiagramLink) IsElementPointer(trans *core.Transaction) bool {
	return diagramLink.ToCore().IsRefinementOfURI(CrlDiagramElementPointerURI, trans)
}

// IsAbstractPointer returns trkue if the link represents an abstract pointer
func (diagramLink *CrlDiagramLink) IsAbstractPointer(trans *core.Transaction) bool {
	return diagramLink.ToCore().IsRefinementOfURI(CrlDiagramAbstractPointerURI, trans)
}

// IsRefinedPointer returns true if the link represents a refined pointer
func (diagramLink *CrlDiagramLink) IsRefinedPointer(trans *core.Transaction) bool {
	return diagramLink.ToCore().IsRefinementOfURI(CrlDiagramRefinedPointerURI, trans)
}

// IsDiagramPointer returns true if the link represents a pointer
func (diagramLink *CrlDiagramLink) IsDiagramPointer(trans *core.Transaction) bool {
	return diagramLink.ToCore().IsRefinementOfURI(CrlDiagramPointerURI, trans)
}

// IsDiagramPointer returns true if the link represents a pointer
func (diagramElement *CrlDiagramElement) IsDiagramPointer(trans *core.Transaction) bool {
	return diagramElement.ToCore().IsRefinementOfURI(CrlDiagramPointerURI, trans)
}

// IsLink returns true if the diagram element is a link
func (diagramElement *CrlDiagramElement) IsLink(trans *core.Transaction) bool {
	return diagramElement.ToCore().IsRefinementOfURI(CrlDiagramLinkURI, trans)
}

// IsNode returns true if the diagram element is a node
func (diagramElement *CrlDiagramElement) IsNode(trans *core.Transaction) bool {
	return diagramElement.ToCore().IsRefinementOfURI(CrlDiagramNodeURI, trans)
}

// SetLinkSource is a convenience function for setting the source concept of a link
func (diagramLink *CrlDiagramLink) SetLinkSource(source *CrlDiagramElement, trans *core.Transaction) {
	sourceReference := diagramLink.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
	if sourceReference != nil {
		sourceReference.SetReferencedConcept(source.ToCore(), core.NoAttribute, trans)
	}
}

// SetLinkTarget is a convenience function for setting the target concept of a link
func (diagramLink *CrlDiagramLink) SetLinkTarget(target *CrlDiagramElement, trans *core.Transaction) {
	targetReference := diagramLink.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
	if targetReference != nil {
		targetReference.SetReferencedConcept(target.ToCore(), core.NoAttribute, trans)
	}
}

// SetNodeHeight is a function on a CrlDiagramNode that sets the height of the diagram node
func (diagramNode *CrlDiagramNode) SetNodeHeight(value float64, trans *core.Transaction) {
	literal := diagramNode.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeHeightURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeWidth is a function on a CrlDiagramNode that sets the width of the diagram node
func (diagramNode *CrlDiagramNode) SetNodeWidth(value float64, trans *core.Transaction) {
	literal := diagramNode.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeWidthURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeX is a function on a CrlDiagramNode that sets the x of the diagram node
func (diagramNode *CrlDiagramNode) SetNodeX(value float64, trans *core.Transaction) {
	literal := diagramNode.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeXURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeY is a function on a CrlDiagramNode that sets the y of the diagram node
func (diagramNode *CrlDiagramNode) SetNodeY(value float64, trans *core.Transaction) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.ToCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetOffsetX sets the x offset value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) SetOffsetX(value float64, trans *core.Transaction) {
	xOffsetLiteral := anchoredText.ToCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextOffsetXURI, trans)
	if xOffsetLiteral == nil {
		log.Printf("SetOffsetX called but no xOffsetLiteral was found")
		return
	}
	xOffset := fmt.Sprintf("%f", value)
	if xOffset != xOffsetLiteral.GetLiteralValue(trans) {
		xOffsetLiteral.SetLiteralValue(xOffset, trans)
	}
}

// SetReferencedModelConcept is a function on a CrlDiagramNode that sets the model element represented by the
// diagram node
func (diagramElement *CrlDiagramElement) SetReferencedModelConcept(el *core.Concept, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	reference := diagramElement.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
	if reference == nil {
		return
	}
	reference.SetReferencedConcept(el, core.NoAttribute, trans)
	diagramElement.updateForModelElementChange(el, trans)
}

// BuildCrlDiagramDomain builds the CrlDiagram concept space and adds it to the uOfD
func BuildCrlDiagramDomain(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) *core.Concept {
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
func updateDiagramElement(changedElement *core.Concept, notification *core.ChangeNotification, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	trans.WriteLockElement(changedElement)
	// core Elements should always be ignored
	if changedElement.GetIsCore(trans) {
		return nil
	}
	if !IsDiagramElement(changedElement, trans) {
		return nil
	}
	diagramElement := (*CrlDiagramElement)(changedElement)
	diagram := diagramElement.GetDiagram(trans)
	if diagram == nil {
		// There is nothing to do
		return nil
	}
	// Suppress circular notifications
	underlyingChange := notification.GetUnderlyingChange()
	if underlyingChange != nil && underlyingChange.IsReferenced(diagramElement.ToCore()) {
		return nil
	}

	// There are several notifications of interest here:
	//   - the deletion of the referenced model element
	//   - the label of the referenced model element
	//   - the list of immediate abstractions of the referenced model element.
	// First, determine whether it is the referenced model element that has changed
	diagramElementModelReference := diagramElement.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
	if diagramElementModelReference == nil {
		// Without a model reference, there is nothing to do. This scenario can occur during diagramElement deletion.
		return nil
	}
	modelElement := diagramElement.GetReferencedModelConcept(trans)
	switch notification.GetNatureOfChange() {
	case core.OwnedConceptChanged:
		switch underlyingChange.GetNatureOfChange() {
		case core.ConceptChanged:
			if underlyingChange.GetReportingElementID() == diagramElementModelReference.GetConceptID(trans) {
				// The underlying change is from the model reference
				diagramElement.updateForModelElementChange(modelElement, trans)
			}
		case core.ReferencedConceptChanged:
			underlyingReportingElementID := underlyingChange.GetReportingElementID()
			if underlyingReportingElementID == diagramElementModelReference.GetConceptID(trans) {
				// The underlying change is from the model reference
				if IsDiagramNode(diagramElement.ToCore(), trans) {
					currentModelElement := underlyingChange.GetAfterConceptState()
					previousModelElement := underlyingChange.GetBeforeConceptState()
					if currentModelElement != nil && previousModelElement != nil {
						if currentModelElement.ReferencedConceptID == "" && previousModelElement.ReferencedConceptID != "" {
							uOfD.DeleteElement(diagramElement.ToCore(), trans)
						} else {
							diagramElement.updateForModelElementChange(modelElement, trans)
						}
					}
				} else if IsDiagramLink(diagramElement.ToCore(), trans) {
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
							uOfD.DeleteElement(diagramElement.ToCore(), trans)
						} else {
							// Otherwise we update the diagram element
							previousReferencedModelElement := uOfD.GetElement(previousReferencedModelElementID)
							switch modelElement.GetConceptType() {
							case core.Reference:
								if IsDiagramElementPointer(diagramElement.ToCore(), trans) {
									currentReferencedModelElement := modelElement.GetReferencedConcept(trans)
									if previousReferencedModelElement != currentReferencedModelElement {
										if currentReferencedModelElement == nil {
											uOfD.DeleteElement(diagramElement.ToCore(), trans)
										} else {
											newTargetDiagramElement := diagram.GetFirstElementRepresentingConcept(currentReferencedModelElement, trans)
											(*CrlDiagramLink)(diagramElement).SetLinkTarget(newTargetDiagramElement, trans)
										}
									}
								} else if IsDiagramReferenceLink(diagramElement.ToCore(), trans) {
									diagramReferenceLink := (*CrlDiagramLink)(diagramElement)
									diagramElement.updateForModelElementChange(modelElement, trans)
									diagramElement.SetDisplayLabel(modelElement.GetLabel(trans), trans)
									newModelTarget := modelElement.GetReferencedConcept(trans)
									newModelSource := modelElement.GetOwningConcept(trans)
									if newModelSource == nil || newModelTarget == nil {
										uOfD.DeleteElement(diagramElement.ToCore(), trans)
										return nil
									}
									currentDiagramSource := diagramReferenceLink.GetLinkSource(trans)
									currentModelSource := currentDiagramSource.GetReferencedModelConcept(trans)
									currentDiagramTarget := diagramReferenceLink.GetLinkTarget(trans)
									currentModelTarget := currentDiagramTarget.GetReferencedModelConcept(trans)
									if currentModelSource != newModelSource {
										newDiagramSource := diagram.GetFirstElementRepresentingConcept(newModelSource, trans)
										if newDiagramSource == nil {
											uOfD.DeleteElement(diagramElement.ToCore(), trans)
											return nil
										}
										diagramReferenceLink.SetLinkSource(newDiagramSource, trans)
									}
									if currentModelTarget != newModelTarget {
										newDiagramTarget := diagram.GetFirstElementRepresentingConcept(newModelTarget, trans)
										if newDiagramTarget == nil {
											uOfD.DeleteElement(diagramElement.ToCore(), trans)
											return nil
										}
										diagramReferenceLink.SetLinkTarget(newDiagramTarget, trans)
									}
								}
							case core.Refinement:
								refinement := modelElement
								if IsDiagramPointer(diagramElement.ToCore(), trans) {
									diagramPointer := (*CrlDiagramLink)(diagramElement)
									var newTargetModelElement *core.Concept
									if IsDiagramAbstractPointer(diagramPointer.ToCore(), trans) {
										newTargetModelElement = refinement.GetAbstractConcept(trans)
									} else if IsDiagramRefinedPointer(diagramPointer.ToCore(), trans) {
										newTargetModelElement = refinement.GetRefinedConcept(trans)
									} else if IsDiagramOwnerPointer(diagramPointer.ToCore(), trans) {
										newTargetModelElement = refinement.GetOwningConcept(trans)
									}
									if previousReferencedModelElement != newTargetModelElement {
										if newTargetModelElement == nil {
											uOfD.DeleteElement(diagramElement.ToCore(), trans)
										} else {
											newTargetDiagramElement := diagram.GetFirstElementRepresentingConcept(newTargetModelElement, trans)
											diagramPointer.SetLinkTarget(newTargetDiagramElement, trans)
										}
									}
								} else if IsDiagramRefinementLink(diagramElement.ToCore(), trans) {
									refinementLink := (*CrlDiagramLink)(diagramElement)
									diagramElement.updateForModelElementChange(modelElement, trans)
									diagramElement.SetDisplayLabel(refinement.GetLabel(trans), trans)
									newModelTarget := refinement.GetAbstractConcept(trans)
									newModelSource := refinement.GetRefinedConcept(trans)
									if newModelTarget == nil || newModelSource == nil {
										uOfD.DeleteElement(diagramElement.ToCore(), trans)
										return nil
									}
									currentDiagramTarget := refinementLink.GetLinkTarget(trans)
									currentModelTarget := currentDiagramTarget.GetReferencedModelConcept(trans)
									currentDiagramSource := refinementLink.GetLinkSource(trans)
									currentModelSource := currentDiagramSource.GetReferencedModelConcept(trans)
									if currentModelTarget != newModelTarget {
										newDiagramTarget := diagram.GetFirstElementRepresentingConcept(newModelTarget, trans)
										if newDiagramTarget == nil {
											uOfD.DeleteElement(diagramElement.ToCore(), trans)
											return nil
										}
										refinementLink.SetLinkTarget(newDiagramTarget, trans)
									}
									if currentModelSource != newModelSource {
										newDiagramSource := diagram.GetFirstElementRepresentingConcept(newModelSource, trans)
										if newDiagramSource == nil {
											uOfD.DeleteElement(diagramElement.ToCore(), trans)
											return nil
										}
										refinementLink.SetLinkSource(newDiagramSource, trans)
									}
								}
							}
						}
					}
				}
			} else {
				// If this is a diagram link and the underlying reporting element is either its source reference or its target reference
				// and the referenced element is now nil, we need to delete the link
				if IsDiagramLink(diagramElement.ToCore(), trans) {
					diagramLink := (*CrlDiagramLink)(diagramElement)
					// If this is a diagram link and the underlying reporting element is either its source reference or its target reference
					// and the referenced element is now nil, we need to delete the link
					underlyingReportingElementIsSourceReference := diagramLink.GetLinkSourceReference(trans).GetConceptID(trans) == underlyingReportingElementID
					underlyingReportingElementIsTargetReference := diagramLink.GetLinkTargetReference(trans).GetConceptID(trans) == underlyingReportingElementID
					if (underlyingReportingElementIsSourceReference || underlyingReportingElementIsTargetReference) &&
						underlyingChange.GetAfterConceptState().ReferencedConceptID == "" {
						uOfD.DeleteElement(diagramElement.ToCore(), trans)
					} else {
						switch underlyingChange.GetNatureOfChange() {
						case core.ReferencedConceptChanged:
							if underlyingReportingElementIsTargetReference {
								// If the link's target has changed, we need to update the underlying model element to reflect the change.
								// Note that if the target is now null, the preceeding clause will have deleted the element
								targetDiagramElement := (*CrlDiagramElement)(uOfD.GetElement(underlyingChange.GetAfterConceptState().ReferencedConceptID))
								targetModelElement := targetDiagramElement.GetReferencedModelConcept(trans)
								if IsDiagramOwnerPointer(diagramElement.ToCore(), trans) {
									modelElement.SetOwningConcept(targetModelElement, trans)
								} else {
									switch modelElement.GetConceptType() {
									case core.Reference:
										// Setting the referenced concepts requires knowledge of what is being referenced
										targetAttribute := core.NoAttribute
										if IsDiagramElementPointer(targetDiagramElement.ToCore(), trans) {
											targetAttribute = core.ReferencedConceptID
										} else if IsDiagramOwnerPointer(targetDiagramElement.ToCore(), trans) {
											targetAttribute = core.OwningConceptID
										} else if IsDiagramAbstractPointer(targetDiagramElement.ToCore(), trans) {
											targetAttribute = core.AbstractConceptID
										} else if IsDiagramRefinedPointer(targetDiagramElement.ToCore(), trans) {
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
								sourceDiagramElement := (*CrlDiagramElement)(uOfD.GetElement(underlyingChange.GetAfterConceptState().ReferencedConceptID))
								sourceModelElement := sourceDiagramElement.GetReferencedModelConcept(trans)
								if modelElement != nil {
									switch modelElement.GetConceptType() {
									case core.Reference:
										if IsDiagramReferenceLink(diagramElement.ToCore(), trans) {
											modelElement.SetOwningConcept(sourceModelElement, trans)
										}
									case core.Refinement:
										if IsDiagramRefinementLink(diagramElement.ToCore(), trans) {
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
				diagramElement.updateForModelElementChange(modelElement, trans)
			}
		}
	case core.ReferencedConceptChanged:
		// We are looking for the model diagramElementModelReference reporting a ConceptChanged which would be the result of setting the referencedConcept
		if notification.GetAfterConceptState().ConceptID != diagramElementModelReference.GetConceptID(trans) {
			break
		}
		if diagramElementModelReference.GetReferencedConceptID(trans) == "" {
			uOfD.DeleteElement(diagramElement.ToCore(), trans)
		} else {
			diagramElement.updateForModelElementChange(modelElement, trans)
		}
	}
	return nil
}

// updateDiagramOwnerPointer updates the ownerPointer's target if the ownership of the represented modelElement changes
func updateDiagramOwnerPointer(concept *core.Concept, notification *core.ChangeNotification, trans *core.Transaction) error {
	// There is one change of interest here: the model element's owner has changed
	uOfD := trans.GetUniverseOfDiscourse()
	trans.WriteLockElement(concept)
	diagramPointer := (*CrlDiagramLink)(concept)
	reportingElement := uOfD.GetElement(notification.GetReportingElementID())
	diagram := diagramPointer.ToCrlDiagramElement().GetDiagram(trans)
	modelElement := diagramPointer.ToCrlDiagramElement().GetReferencedModelConcept(trans)
	switch notification.GetNatureOfChange() {
	case core.OwnedConceptChanged:
		if reportingElement == modelElement {
			underlyingNotification := notification.GetUnderlyingChange()
			switch underlyingNotification.GetNatureOfChange() {
			case core.OwningConceptChanged:
				if underlyingNotification.GetAfterConceptState().ConceptID == modelElement.GetConceptID(trans) {
					modelOwner := modelElement.GetOwningConcept(trans)
					var oldModelOwner *core.Concept
					diagramTarget := diagramPointer.GetLinkTarget(trans)
					if diagramTarget != nil {
						oldModelOwner = diagramTarget.GetReferencedModelConcept(trans)
					}
					if modelOwner != oldModelOwner {
						// Need to determine whether there is a view of the new owner in the diagram
						newDiagramTarget := (*CrlDiagramElement)(diagram.GetFirstElementRepresentingConcept(modelOwner, trans))
						if newDiagramTarget == nil {
							// There is no view, delete the modelElement
							uOfD.DeleteElement(diagramPointer.ToCore(), trans)
						} else {
							diagramPointer.SetLinkTarget(newDiagramTarget, trans)
						}
					}
				}
			}
			break
		}
		// We are looking for a notification from either the source or target reference in the diagram
		// If either source or target are nil, delete the pointer
		sourceReference := diagramPointer.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
		targetReference := diagramPointer.ToCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
		if reportingElement == sourceReference || reportingElement == targetReference {
			underlyingNotification := notification.GetUnderlyingChange()
			switch underlyingNotification.GetNatureOfChange() {
			case core.ReferencedConceptChanged:
				switch reportingElement.GetConceptType() {
				case core.Reference:
					if reportingElement.GetReferencedConcept(trans) == nil {
						uOfD.DeleteElement(diagramPointer.ToCore(), trans)
					}
				}
			}
		}
	}
	return nil
}

func (diagramElement *CrlDiagramElement) updateForModelElementChange(modelElement *core.Concept, trans *core.Transaction) {
	modelElementLabel := ""
	if modelElement != nil {
		modelElementLabel = modelElement.GetLabel(trans)
		if modelElementLabel != diagramElement.ToCore().GetLabel(trans) {
			newLabel := modelElementLabel
			if IsDiagramPointer(diagramElement.ToCore(), trans) {
				if IsDiagramOwnerPointer(diagramElement.ToCore(), trans) {
					newLabel = newLabel + " Owner Pointer"
				} else if IsDiagramAbstractPointer(diagramElement.ToCore(), trans) {
					newLabel = newLabel + " Abstract Pointer"
				} else if IsDiagramRefinedPointer(diagramElement.ToCore(), trans) {
					newLabel = newLabel + " Refined Pointer"
				} else if IsDiagramElementPointer(diagramElement.ToCore(), trans) {
					newLabel = newLabel + " Referenced Concept Pointer"
				}
			}
			diagramElement.ToCore().SetLabel(newLabel, trans)
			if !IsDiagramPointer(diagramElement.ToCore(), trans) {
				diagramElement.SetDisplayLabel(modelElementLabel, trans)
			}
		}
		abstractions := make(map[string]*core.Concept)
		modelElement.FindImmediateAbstractions(abstractions, trans)
		abstractionsLabel := ""
		for _, abs := range abstractions {
			if len(abstractionsLabel) != 0 {
				abstractionsLabel += "\n"
			}
			abstractionsLabel += abs.GetLabel(trans)
		}
		if diagramElement.GetAbstractionDisplayLabel(trans) != abstractionsLabel {
			diagramElement.SetAbstractionDisplayLabel(abstractionsLabel, trans)
		}
	}
}

// updateNodeSize recalcualtes the size of the node based on the string sizes for the display label and
// abstractions listed
func (diagramNode *CrlDiagramNode) updateNodeSize(trans *core.Transaction) {
	displayLabel := diagramNode.ToCrlDiagramElement().GetDisplayLabel(trans)
	displayLabelBounds, _ := font.BoundString(go12PtBoldFace, displayLabel)
	displayLabelMaxHeight := Int26_6ToFloat(displayLabelBounds.Max.Y)
	displayLabelMaxWidth := Int26_6ToFloat(displayLabelBounds.Max.X)
	displayLabelMinHeight := Int26_6ToFloat(displayLabelBounds.Min.Y)
	displayLabelMinWidth := Int26_6ToFloat(displayLabelBounds.Min.X)
	displayLabelHeight := displayLabelMaxHeight - displayLabelMinHeight
	displayLabelWidth := displayLabelMaxWidth - displayLabelMinWidth
	abstractionDisplayLabel := diagramNode.ToCrlDiagramElement().GetAbstractionDisplayLabel(trans)
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
	diagramNode.SetNodeHeight(height, trans)
	diagramNode.SetNodeWidth(width, trans)
	// SetNodeDisplayLabelYOffset(node, displayLabelYOffset, trans)
}
