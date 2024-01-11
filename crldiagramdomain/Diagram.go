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

// Multiplicity Map

// CrlDiagramMultiplicityMapURI identifies a mapping from a model multiplicity to an anchored text literal value
var CrlDiagramMultiplicityMapURI = CrlDiagramURI + "/" + "MultiplicityMap"

// CrlDiagramMultiplicityMap manages the mapping between a model multiplicity and a diagram anchored text literal value
type CrlDiagramMultiplicityMap core.Concept

// CrlDiagramMultiplicityMapMultiplicityReferenceURI identifies the model multiplicity associated with the anchored text
var CrlDiagramMultiplicityMapMultiplicityReferenceURI = CrlDiagramMultiplicityMapURI + "/" + "MultiplicityReference"

// CrlDiagramMultiplicityMapAnchoredTextReferenceURI identifies the diagram anchored text associated with the anchored text
var CrlDiagramMultiplicityMapAnchoredTextReferenceURI = CrlDiagramMultiplicityMapURI + "/" + "AnchoredText"

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

// CrlDiagramLinkSourceMultiplicityURI identifies the anchored text used to display target multiplicity
var CrlDiagramLinkSourceMultiplicityURI = CrlDiagramLinkURI + "/" + "SourceMultiplicity"

// CrlDiagramLinkTargetMultiplicityURI identifies the anchored text used to display target multiplicity
var CrlDiagramLinkTargetMultiplicityURI = CrlDiagramLinkURI + "/" + "TargetMultiplicity"

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
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextAnchorXURI, anchoredText.AsCore(), "AnchorX", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextAnchorYURI, anchoredText.AsCore(), "AnchorY", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextOffsetXURI, anchoredText.AsCore(), "OffsetX", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextOffsetYURI, anchoredText.AsCore(), "OffsetY", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramAnchoredTextVisibleURI, anchoredText.AsCore(), "Visible", trans)
}

func addDiagramElementConcepts(newElement *CrlDiagramElement, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementModelReferenceURI, newElement.AsCore(), "Model Reference", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementLineColorURI, newElement.AsCore(), "Line Color", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementBGColorURI, newElement.AsCore(), "BG Color", trans)
}

func addDiagramLinkConcepts(newLink *CrlDiagramLink, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkSourceURI, newLink.AsCore(), "Source", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkTargetURI, newLink.AsCore(), "Target", trans)
	displayLabel, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkDisplayLabelURI, newLink.AsCore(), "DisplayLabel", trans)
	uOfD.AddAbstractionURIToConcept(displayLabel, CrlDiagramAnchoredTextURI, trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(displayLabel), trans)
	abstractionDisplayLabel, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkAbstractionDisplayLabelURI, newLink.AsCore(), "AbstractionDisplayLabel", trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(abstractionDisplayLabel), trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(abstractionDisplayLabel), trans)
}

func addDiagramNodeConcepts(newNode *CrlDiagramNode, trans *core.Transaction) {
	uOfD := trans.GetUniverseOfDiscourse()
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeXURI, newNode.AsCore(), "X", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeYURI, newNode.AsCore(), "Y", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeHeightURI, newNode.AsCore(), "Height", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramNodeWidthURI, newNode.AsCore(), "Width", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementDisplayLabelURI, newNode.AsCore(), "DisplayLabel", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementAbstractionDisplayLabelURI, newNode.AsCore(), "AbstractionDisplayLabel", trans)
}

// AsCore casts a diagram element to core.Concept
func (diagramElement *CrlDiagramElement) AsCore() *core.Concept {
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
	abstractionDisplayLabelLiteral := diagramElement.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, trans)
	if abstractionDisplayLabelLiteral != nil {
		return abstractionDisplayLabelLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetAnchoredTextWithLabel returns the diagram element's first anchored text child with the given label
func (diagramElement *CrlDiagramElement) GetAnchoredTextWithLabel(label string, trans *core.Transaction) *CrlDiagramAnchoredText {
	if !diagramElement.AsCore().IsRefinementOfURI(CrlDiagramElementURI, trans) {
		log.Print("GetAnchoredTextWithLabel called for a concept that is not a CrlDiagramElement")
		return nil
	}
	anchoredTexts := diagramElement.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramAnchoredTextURI, trans)
	for _, anchoredText := range anchoredTexts {
		if anchoredText.GetLabel(trans) == label {
			return (*CrlDiagramAnchoredText)(anchoredText)
		}
	}
	return nil
}

// GetBGColor is a convenience function for getting the backgound color value of a DiagramElement
func (diagramElement *CrlDiagramElement) GetBGColor(trans *core.Transaction) string {
	BGColorLiteral := diagramElement.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementBGColorURI, trans)
	if BGColorLiteral != nil {
		return BGColorLiteral.GetLiteralValue(trans)
	}
	return ""
}

// GetDiagram returns the diagram within which the DiagramElement appears
func (diagramElement *CrlDiagramElement) GetDiagram(trans *core.Transaction) *CrlDiagram {
	return (*CrlDiagram)(diagramElement.AsCore().GetOwningConcept(trans))
}

// GetDisplayLabel is a convenience function for getting the DisplayLabel value of a DiagramElement
func (diagramElement *CrlDiagramElement) GetDisplayLabel(trans *core.Transaction) string {
	displayLabelLiteral := diagramElement.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
	if displayLabelLiteral != nil {
		return displayLabelLiteral.GetLiteralValue(trans)
	}
	return ""
}

// AsCore casts the diagram to core.Concept
func (diagram *CrlDiagram) AsCore() *core.Concept {
	return (*core.Concept)(diagram)
}

// GetFirstElementRepresentingConcept returns the first non-pointer diagram element that represents the indicated concept
func (diagram *CrlDiagram) GetFirstElementRepresentingConcept(concept *core.Concept, trans *core.Transaction) *CrlDiagramElement {
	if concept == nil {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept && !el.IsRefinementOfURI(CrlDiagramPointerURI, trans) {
			return (*CrlDiagramElement)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptID returns the first diagram element that represents the indicated concept
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptID(conceptID string, trans *core.Transaction) *CrlDiagramElement {
	if conceptID == "" {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID && !el.IsRefinementOfURI(CrlDiagramPointerURI, trans) {
			return (*CrlDiagramElement)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptOwnerPointer returns the first diagram element that represents the indicated concept's OwnerPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptOwnerPointer(concept *core.Concept, trans *core.Transaction) *CrlDiagramLink {
	if concept == nil {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDOwnerPointer returns the first diagram element that represents the indicated concept's OwnerPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptIDOwnerPointer(conceptID string, trans *core.Transaction) *CrlDiagramLink {
	if conceptID == "" {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptElementPointer returns the first diagram element that represents the indicated concept's ElementPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptElementPointer(concept *core.Concept, trans *core.Transaction) *CrlDiagramLink {
	if concept == nil {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDElementPointer returns the first diagram element that represents the indicated concept's ElementPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptIDElementPointer(conceptID string, trans *core.Transaction) *CrlDiagramLink {
	if conceptID == "" {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptAbstractPointer returns the first diagram element that represents the indicated concept's AbstractPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptAbstractPointer(concept *core.Concept, trans *core.Transaction) *CrlDiagramLink {
	if concept == nil {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramAbstractPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDAbstractPointer returns the first diagram element that represents the indicated concept's AbstractPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptIDAbstractPointer(conceptID string, trans *core.Transaction) *CrlDiagramLink {
	if conceptID == "" {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramAbstractPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans).GetConceptID(trans) == conceptID {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptRefinedPointer returns the first diagram element that represents the indicated concept's RefinedPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptRefinedPointer(concept *core.Concept, trans *core.Transaction) *CrlDiagramLink {
	if concept == nil {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramRefinedPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDRefinedPointer returns the first diagram element that represents the indicated concept's RefinedPointer
func (diagram *CrlDiagram) GetFirstElementRepresentingConceptIDRefinedPointer(conceptID string, trans *core.Transaction) *CrlDiagramLink {
	if conceptID == "" {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramRefinedPointerURI, trans) {
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
	lineColorLiteral := diagramElement.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementLineColorURI, trans)
	if lineColorLiteral != nil {
		return lineColorLiteral.GetLiteralValue(trans)
	}
	return ""
}

// AsCore casts the CrlDiagramLink to core.Concept
func (diagramLink *CrlDiagramLink) AsCore() *core.Concept {
	return (*core.Concept)(diagramLink)
}

// GetDisplayLabel returns the CrlAnchoredText for the display label
func (diagramLink *CrlDiagramLink) GetDisplayLabel(trans *core.Transaction) *CrlDiagramAnchoredText {
	at := diagramLink.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramLinkDisplayLabelURI, trans)
	return (*CrlDiagramAnchoredText)(at)
}

// GetLinkMultiplicityMap is a convenience function for getting the multiplicity map for a link
func (diagramLink *CrlDiagramLink) GetLinkMultiplicityMap(trans *core.Transaction) *CrlDiagramMultiplicityMap {
	multiplicityMap := diagramLink.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramMultiplicityMapURI, trans)
	return (*CrlDiagramMultiplicityMap)(multiplicityMap)
}

// GetLinkSource is a convenience function for getting the source concept of a link
func (diagramLink *CrlDiagramLink) GetLinkSource(trans *core.Transaction) *CrlDiagramElement {
	sourceReference := diagramLink.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
	if sourceReference != nil {
		referencedConcept := sourceReference.GetReferencedConcept(trans)
		if referencedConcept != nil && referencedConcept.IsRefinementOfURI(CrlDiagramElementURI, trans) {
			return (*CrlDiagramElement)(referencedConcept)
		}
	}
	return nil
}

// GetLinkSourceReference is a convenience function for getting the source reference of a link
func (diagramLink *CrlDiagramLink) GetLinkSourceReference(trans *core.Transaction) *core.Concept {
	return diagramLink.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
}

// GetLinkTarget is a convenience function for getting the target concept of a link
func (diagramLink *CrlDiagramLink) GetLinkTarget(trans *core.Transaction) *CrlDiagramElement {
	targetReference := diagramLink.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
	if targetReference != nil {
		referencedConcept := targetReference.GetReferencedConcept(trans)
		if referencedConcept != nil && referencedConcept.IsRefinementOfURI(CrlDiagramElementURI, trans) {
			return (*CrlDiagramElement)(referencedConcept)
		}
	}
	return nil
}

// GetLinkTargetReference is a convenience function for getting the target reference of a link
func (diagramLink *CrlDiagramLink) GetLinkTargetReference(trans *core.Transaction) *core.Concept {
	return diagramLink.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
}

// GetLinkSourceMultiplicity returns the AnchoredText representing the link's source multiplicity
func (diagramLink *CrlDiagramLink) GetLinkSourceMultiplicity(trans *core.Transaction) *CrlDiagramAnchoredText {
	return (*CrlDiagramAnchoredText)(diagramLink.AsCore().GetFirstOwnedLiteralRefinedFromURI(CrlDiagramLinkSourceMultiplicityURI, trans))
}

// GetLinkTargetMultiplicity returns the AnchoredText representing the link's target multiplicity
func (diagramLink *CrlDiagramLink) GetLinkTargetMultiplicity(trans *core.Transaction) *CrlDiagramAnchoredText {
	return (*CrlDiagramAnchoredText)(diagramLink.AsCore().GetFirstOwnedLiteralRefinedFromURI(CrlDiagramLinkTargetMultiplicityURI, trans))
}

// AsCore casts CrlDiagramMultiplicityMap to core.Concept
func (multiplicityMap *CrlDiagramMultiplicityMap) AsCore() *core.Concept {
	return (*core.Concept)(multiplicityMap)
}

// GetAnchoredText returns the diagram's AnchoredText (a Literal) for this map
func (multiplicityMap *CrlDiagramMultiplicityMap) GetAnchoredText(trans *core.Transaction) *CrlDiagramAnchoredText {
	anchoredTextReference := multiplicityMap.GetAnchoredTextReference(trans)
	if anchoredTextReference == nil {
		return nil
	}
	return (*CrlDiagramAnchoredText)(anchoredTextReference.GetReferencedConcept(trans))
}

// GetAnchoredTextReference returns the map's AnchoredText reference
func (multiplicityMap *CrlDiagramMultiplicityMap) GetAnchoredTextReference(trans *core.Transaction) *core.Concept {
	return multiplicityMap.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramMultiplicityMapAnchoredTextReferenceURI, trans)
}

// GetModelMultiplicity returns the model's Multiplicity (a Literal) for this map
func (multiplicityMap *CrlDiagramMultiplicityMap) GetModelMultiplicity(trans *core.Transaction) *core.Concept {
	multiplicityReference := multiplicityMap.GetModelMultiplicityReference(trans)
	if multiplicityReference == nil {
		return nil
	}
	return multiplicityReference.GetReferencedConcept(trans)
}

// GetModelMultiplicityReference returns the model multiplicity reference for this map
func (multiplicityMap *CrlDiagramMultiplicityMap) GetModelMultiplicityReference(trans *core.Transaction) *core.Concept {
	return multiplicityMap.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramMultiplicityMapMultiplicityReferenceURI, trans)
}

// SetAnchoredText sets the diagram's AnchoredText (a Literal) for this map
func (multiplicityMap *CrlDiagramMultiplicityMap) SetAnchoredText(anchoredText *CrlDiagramAnchoredText, trans *core.Transaction) {
	anchoredTextReference := multiplicityMap.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramMultiplicityMapAnchoredTextReferenceURI, trans)
	if anchoredTextReference == nil {
		return
	}
	anchoredTextReference.SetReferencedConcept(anchoredText.AsCore(), core.NoAttribute, trans)
}

// SetModelMultiplicity sets the model's Multiplicity (a Literal) for this map
func (multiplicityMap *CrlDiagramMultiplicityMap) SetModelMultiplicity(multiplicity *core.Concept, trans *core.Transaction) {
	multiplicityReference := multiplicityMap.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramMultiplicityMapMultiplicityReferenceURI, trans)
	if multiplicityReference == nil || !multiplicity.IsLiteral() {
		return
	}
	multiplicityReference.SetReferencedConcept(multiplicity, core.NoAttribute, trans)
}

// AsCore casts CrlDiagramNode to core.Concept
func (diagramNode *CrlDiagramNode) AsCore() *core.Concept {
	return (*core.Concept)(diagramNode)
}

// GetNodeHeight is a convenience function for getting the Height value of a node's position
func (diagramNode *CrlDiagramNode) GetNodeHeight(trans *core.Transaction) float64 {
	heightLiteral := diagramNode.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeHeightURI, trans)
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
	widthLiteral := diagramNode.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeWidthURI, trans)
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
	xLiteral := diagramNode.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeXURI, trans)
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
	yLiteral := diagramNode.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, trans)
	if yLiteral != nil {
		value := yLiteral.GetLiteralValue(trans)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// AsCore casts CrlDiagramAnchoredText to core.Concept
func (anchoredText *CrlDiagramAnchoredText) AsCore() *core.Concept {
	return (*core.Concept)(anchoredText)
}

// GetOffsetX returns the x offset value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) GetOffsetX(trans *core.Transaction) float64 {
	xOffsetLiteral := anchoredText.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextOffsetXURI, trans)
	if xOffsetLiteral == nil {
		return 0
	}
	xOffset, err := strconv.ParseFloat(xOffsetLiteral.GetLiteralValue(trans), 64)
	if err != nil {
		errors.Wrap(err, "GetOffsetX failed")
		return 0
	}
	return xOffset
}

// GetOffsetY returns the x offset value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) GetOffsetY(trans *core.Transaction) float64 {
	xOffsetLiteral := anchoredText.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextOffsetYURI, trans)
	if xOffsetLiteral == nil {
		return 0
	}
	xOffset, err := strconv.ParseFloat(xOffsetLiteral.GetLiteralValue(trans), 64)
	if err != nil {
		errors.Wrap(err, "GetOffsetY failed")
		return 0
	}
	return xOffset
}

// GetAnchorX returns the x anchor value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) GetAnchorX(trans *core.Transaction) float64 {
	xAnchorLiteral := anchoredText.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextAnchorXURI, trans)
	if xAnchorLiteral == nil {
		log.Printf("GetAnchorX called but no xAnchorLiteral was found")
		return 0
	}
	xAnchor, err := strconv.ParseFloat(xAnchorLiteral.GetLiteralValue(trans), 64)
	if err != nil {
		errors.Wrap(err, "GetAnchorX failed")
		return 0
	}
	return xAnchor
}

// GetAnchorY returns the x anchor value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) GetAnchorY(trans *core.Transaction) float64 {
	yAnchorLiteral := anchoredText.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextAnchorYURI, trans)
	if yAnchorLiteral == nil {
		log.Printf("GetAnchorY called but no xAnchorLiteral was found")
		return 0
	}
	yAnchor, err := strconv.ParseFloat(yAnchorLiteral.GetLiteralValue(trans), 64)
	if err != nil {
		errors.Wrap(err, "GetAnchorY failed")
		return 0
	}
	return yAnchor
}

// GetMultiplicityReference returns the multiplicity reference if one exists
func (anchoredText *CrlDiagramAnchoredText) GetMultiplicityReference(trans *core.Transaction) *core.Concept {
	return anchoredText.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramMultiplicityMapMultiplicityReferenceURI, trans)
}

// SetOffsetX sets the x offset value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) SetOffsetX(value float64, trans *core.Transaction) {
	xOffsetLiteral := anchoredText.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextOffsetXURI, trans)
	if xOffsetLiteral == nil {
		log.Printf("SetOffsetX called but no xOffsetLiteral was found")
		return
	}
	xOffset := fmt.Sprintf("%f", value)
	if xOffset != xOffsetLiteral.GetLiteralValue(trans) {
		xOffsetLiteral.SetLiteralValue(xOffset, trans)
	}
}

// SetOffsetY sets the y offset value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) SetOffsetY(value float64, trans *core.Transaction) {
	yOffsetLiteral := anchoredText.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextOffsetYURI, trans)
	if yOffsetLiteral == nil {
		log.Printf("SetOffsetY called but no xOffsetLiteral was found")
		return
	}
	yOffset := fmt.Sprintf("%f", value)
	if yOffset != yOffsetLiteral.GetLiteralValue(trans) {
		yOffsetLiteral.SetLiteralValue(yOffset, trans)
	}
}

// SetAnchorX sets the x anchor value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) SetAnchorX(value float64, trans *core.Transaction) {
	xAnchorLiteral := anchoredText.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextAnchorXURI, trans)
	if xAnchorLiteral == nil {
		log.Printf("SetAnchorX called but no xAnchorLiteral was found")
		return
	}
	xAnchor := fmt.Sprintf("%f", value)
	if xAnchor != xAnchorLiteral.GetLiteralValue(trans) {
		xAnchorLiteral.SetLiteralValue(xAnchor, trans)
	}
}

// SetAnchorY sets the y anchor value for an anchored text
func (anchoredText *CrlDiagramAnchoredText) SetAnchorY(value float64, trans *core.Transaction) {
	yAnchorLiteral := anchoredText.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlDiagramAnchoredTextAnchorYURI, trans)
	if yAnchorLiteral == nil {
		log.Printf("SetAnchorY called but no xAnchorLiteral was found")
		return
	}
	yAnchor := fmt.Sprintf("%f", value)
	if yAnchor != yAnchorLiteral.GetLiteralValue(trans) {
		yAnchorLiteral.SetLiteralValue(yAnchor, trans)
	}
}

// GetOwnerPointer returns the ownerPointer for the concept if one exists
func (diagram *CrlDiagram) GetOwnerPointer(concept *CrlDiagramElement, trans *core.Transaction) *CrlDiagramLink {
	if concept == nil {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept.AsCore() {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetElementPointer returns the elementPointer for the concept if one exists
func (diagram *CrlDiagram) GetElementPointer(concept *CrlDiagramElement, trans *core.Transaction) *CrlDiagramLink {
	if concept == nil {
		return nil
	}
	for _, el := range diagram.AsCore().GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, trans) {
		if (*CrlDiagramElement)(el).GetReferencedModelConcept(trans) == concept.AsCore() {
			return (*CrlDiagramLink)(el)
		}
	}
	return nil
}

// GetReferencedModelConcept is a function on a CrlDiagramElement that returns the model element represented by the
// diagram node
func (diagramElement *CrlDiagramElement) GetReferencedModelConcept(trans *core.Transaction) *core.Concept {
	reference := diagramElement.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
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
	if concept != nil && concept.IsRefinementOfURI(CrlDiagramURI, trans) {
		return (*CrlDiagram)(concept)
	}
	return nil
}

// GetCrlDiagramElement returns the CrlDiagramElement with the given ID
func GetCrlDiagramElement(id string, trans *core.Transaction) *CrlDiagramElement {
	concept := trans.GetUniverseOfDiscourse().GetElement(id)
	if concept != nil && concept.IsRefinementOfURI(CrlDiagramElementURI, trans) {
		return (*CrlDiagramElement)(concept)
	}
	return nil
}

// GetCrlDiagramLink returns the CrlDiagramLink with the given ID
func GetCrlDiagramLink(id string, trans *core.Transaction) *CrlDiagramLink {
	concept := trans.GetUniverseOfDiscourse().GetElement(id)
	if concept != nil && concept.IsRefinementOfURI(CrlDiagramLinkURI, trans) {
		return (*CrlDiagramLink)(concept)
	}
	return nil
}

// GetCrlDiagramNode returns the CrlDiagramNode with the given ID
func GetCrlDiagramNode(id string, trans *core.Transaction) *CrlDiagramNode {
	concept := trans.GetUniverseOfDiscourse().GetElement(id)
	if concept != nil && concept.IsRefinementOfURI(CrlDiagramNodeURI, trans) {
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
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramAbstractPointerURI, trans)
}

// IsDiagramElement returns true if the supplied element is a CrlDiagramElement
func IsDiagramElement(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramElementURI, trans)
}

// IsDiagramElementPointer returns true if the supplied element is a CrlDiagramElementPointer
func IsDiagramElementPointer(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramElementPointerURI, trans)
}

// IsDiagramLink returns true if the supplied element is a CrlDiagramLink
func IsDiagramLink(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramLinkURI, trans)
}

// IsDiagramNode returns true if the supplied element is a CrlDiagramNode
func IsDiagramNode(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramNodeURI, trans)
}

// IsDiagramOwnerPointer returns true if the supplied element is a CrlDiagramOwnerPointer
func IsDiagramOwnerPointer(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramOwnerPointerURI, trans)
}

// IsDiagramPointer returns true if the supplied element is a CrlDiagramPointer
func IsDiagramPointer(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramPointerURI, trans)
}

// IsDiagramRefinedPointer returns true if the supplied element is a CrlDiagramRefinedPointer
func IsDiagramRefinedPointer(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramRefinedPointerURI, trans)
}

// IsDiagramReferenceLink returns true if the supplied element is a CrlDiagramReferenceLink
func IsDiagramReferenceLink(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramReferenceLinkURI, trans)
}

// IsDiagramRefinementLink returns true if the supplied element is a CrlDiagramRefinementLink
func IsDiagramRefinementLink(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
	return el.IsRefinementOfURI(CrlDiagramRefinementLinkURI, trans)
}

// IsModelReference returns true if the supplied element is a ModelReference
func IsModelReference(el *core.Concept, trans *core.Transaction) bool {
	if el == nil {
		return false
	}
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
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramWidthURI, newDiagram.AsCore(), "Width", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramHeightURI, newDiagram.AsCore(), "Height", trans)
	return newDiagram, nil
}

// AsDiagramElement casts a CrlDiagramNode to a CrlDiagramElement
func (diagramNode *CrlDiagramNode) AsDiagramElement() *CrlDiagramElement {
	return (*CrlDiagramElement)(diagramNode)
}

// AsCrlDiagramElement casts a CrlDiagramLink to a CrlDiagramElement
func (diagramLink *CrlDiagramLink) AsCrlDiagramElement() *CrlDiagramElement {
	return (*CrlDiagramElement)(diagramLink)
}

// NewDiagramAbstractPointer creates a new DiagramAbstractPointer
func NewDiagramAbstractPointer(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramAbstractPointerURI, "AbstractPointer", trans)
	newPointer := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts(newPointer.AsCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// NewDiagramNode creates a new diagram node
func NewDiagramNode(trans *core.Transaction) (*CrlDiagramNode, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newElement, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramNodeURI, "New Node", trans)
	newNode := (*CrlDiagramNode)(newElement)
	addDiagramElementConcepts(newNode.AsDiagramElement(), trans)
	addDiagramNodeConcepts(newNode, trans)
	newNode.AsDiagramElement().SetLineColor("#00000000", trans)
	return newNode, nil
}

// NewDiagramOwnerPointer creates a new DiagramOwnerPointer
func NewDiagramOwnerPointer(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramOwnerPointerURI, "OwnerPointer", trans)
	newPointer := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts(newPointer.AsCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newPointer, trans)
	sourceMultiplicity, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkSourceMultiplicityURI, newPointer.AsCore(), "Multiplicity", trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(sourceMultiplicity), trans)
	NewDiagramMultiplicityMap(newObject, (*CrlDiagramAnchoredText)(sourceMultiplicity), trans)
	return newPointer, nil
}

// NewDiagramElementPointer creates a new DiagramElementPointer
func NewDiagramElementPointer(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramElementPointerURI, "ElementPointer", trans)
	newPointer := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts((*CrlDiagramElement)(newPointer.AsCore()), trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// NewDiagramMultiplicityMap creates a new multiplicity map
func NewDiagramMultiplicityMap(owner *core.Concept, anchoredText *CrlDiagramAnchoredText, trans *core.Transaction) *CrlDiagramMultiplicityMap {
	uOfD := trans.GetUniverseOfDiscourse()
	newMap, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramMultiplicityMapURI, owner, "MultiplicityMap", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramMultiplicityMapMultiplicityReferenceURI, newMap, "MultiplicityReference", trans)
	anchoredTextReference, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramMultiplicityMapAnchoredTextReferenceURI, newMap, "AnchoredTextReference", trans)
	anchoredTextReference.SetReferencedConcept(anchoredText.AsCore(), core.NoAttribute, trans)
	return (*CrlDiagramMultiplicityMap)(newMap)
}

// NewDiagramReferenceLink creates a new diagram link to represent a reference
func NewDiagramReferenceLink(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newElement, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramReferenceLinkURI, "ReferenceLink", trans)
	newLink := (*CrlDiagramLink)(newElement)
	addDiagramElementConcepts(newLink.AsCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newLink, trans)
	targetMultiplicity, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramLinkTargetMultiplicityURI, newLink.AsCore(), "TargetMultiplicity", trans)
	addAnchoredTextConcepts((*CrlDiagramAnchoredText)(targetMultiplicity), trans)
	NewDiagramMultiplicityMap(newElement, (*CrlDiagramAnchoredText)(targetMultiplicity), trans)
	return newLink, nil
}

// NewDiagramRefinementLink creates a new diagram link representing a refinement
func NewDiagramRefinementLink(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramRefinementLinkURI, "RefinementLink", trans)
	newLink := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts(newLink.AsCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newLink, trans)
	return newLink, nil
}

// NewDiagramRefinedPointer creates a new DiagramRefinedPointer
func NewDiagramRefinedPointer(trans *core.Transaction) (*CrlDiagramLink, error) {
	uOfD := trans.GetUniverseOfDiscourse()
	newObject, _ := uOfD.CreateRefinementOfConceptURI(CrlDiagramRefinedPointerURI, "RefinedPointer", trans)
	newPointer := (*CrlDiagramLink)(newObject)
	addDiagramElementConcepts(newPointer.AsCrlDiagramElement(), trans)
	addDiagramLinkConcepts(newPointer, trans)
	return newPointer, nil
}

// SetAbstractionDisplayLabel is a function on a CrlDiagramElement that sets the abstraction display label of the diagram node
func (diagramElement *CrlDiagramElement) SetAbstractionDisplayLabel(value string, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(value, trans)
	if IsDiagramNode(diagramElement.AsCore(), trans) {
		(*CrlDiagramNode)(diagramElement).updateNodeSize(trans)
	}
}

// SetDisplayLabel is a function on a CrlDiagramNode that sets the display label of the diagram element.
// If the diagram element is a pointer, the value is ignored and the label is set to the empty string
func (diagramElement *CrlDiagramElement) SetDisplayLabel(value string, trans *core.Transaction) {
	literal := diagramElement.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, trans)
	if literal == nil {
		return
	}
	if IsDiagramPointer(diagramElement.AsCore(), trans) {
		literal.SetLiteralValue("", trans)
	} else {
		literal.SetLiteralValue(value, trans)
	}
	if IsDiagramNode(diagramElement.AsCore(), trans) {
		(*CrlDiagramNode)(diagramElement).updateNodeSize(trans)
	}
}

// SetLineColor is a function on a CrlDiagramElement that sets the line color for the diagram element.
func (diagramElement *CrlDiagramElement) SetLineColor(value string, trans *core.Transaction) {
	literal := diagramElement.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementLineColorURI, trans)
	if literal == nil {
		// This is remedial code: the literal should already be there
		uOfD := trans.GetUniverseOfDiscourse()
		literal, _ = uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementLineColorURI, diagramElement.AsCore(), "Line Color", trans)
	}
	literal.SetLiteralValue(value, trans)
}

// SetBGColor is a function on a CrlDiagramNode that sets the background color for the diagram element.
// If the diagram element is a pointer, the value is ignored and the label is set to the empty string
func (diagramElement *CrlDiagramElement) SetBGColor(value string, trans *core.Transaction) {
	literal := diagramElement.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementBGColorURI, trans)
	if literal == nil {
		// This is remedial code: the literal should already be there
		uOfD := trans.GetUniverseOfDiscourse()
		literal, _ = uOfD.CreateOwnedRefinementOfConceptURI(CrlDiagramElementBGColorURI, diagramElement.AsCore(), "Background Color", trans)
	}
	if IsDiagramPointer(diagramElement.AsCore(), trans) {
		literal.SetLiteralValue("", trans)
	} else {
		literal.SetLiteralValue(value, trans)
	}
}

// SetDiagram sets the owner of the DiagramElement to the Diagram
func (diagramElement *CrlDiagramElement) SetDiagram(diagram *CrlDiagram, trans *core.Transaction) {
	diagramElement.AsCore().SetOwningConcept(diagram.AsCore(), trans)
}

// IsRefinementLink returns true if the link represents a refinement
func (diagramLink *CrlDiagramLink) IsRefinementLink(trans *core.Transaction) bool {
	return diagramLink.AsCore().IsRefinementOfURI(CrlDiagramRefinementLinkURI, trans)
}

// IsReferenceLink returns true if the link represents a reference
func (diagramLink *CrlDiagramLink) IsReferenceLink(trans *core.Transaction) bool {
	return diagramLink.AsCore().IsRefinementOfURI(CrlDiagramReferenceLinkURI, trans)
}

// IsOwnerPointer returns true if the link represents an owner pointer
func (diagramLink *CrlDiagramLink) IsOwnerPointer(trans *core.Transaction) bool {
	return diagramLink.AsCore().IsRefinementOfURI(CrlDiagramOwnerPointerURI, trans)
}

// IsElementPointer returns true if the link represents a referenced element pointer
func (diagramLink *CrlDiagramLink) IsElementPointer(trans *core.Transaction) bool {
	return diagramLink.AsCore().IsRefinementOfURI(CrlDiagramElementPointerURI, trans)
}

// IsAbstractPointer returns trkue if the link represents an abstract pointer
func (diagramLink *CrlDiagramLink) IsAbstractPointer(trans *core.Transaction) bool {
	return diagramLink.AsCore().IsRefinementOfURI(CrlDiagramAbstractPointerURI, trans)
}

// IsRefinedPointer returns true if the link represents a refined pointer
func (diagramLink *CrlDiagramLink) IsRefinedPointer(trans *core.Transaction) bool {
	return diagramLink.AsCore().IsRefinementOfURI(CrlDiagramRefinedPointerURI, trans)
}

// IsDiagramPointer returns true if the link represents a pointer
func (diagramLink *CrlDiagramLink) IsDiagramPointer(trans *core.Transaction) bool {
	return diagramLink.AsCore().IsRefinementOfURI(CrlDiagramPointerURI, trans)
}

// IsDiagramPointer returns true if the link represents a pointer
func (diagramElement *CrlDiagramElement) IsDiagramPointer(trans *core.Transaction) bool {
	return diagramElement.AsCore().IsRefinementOfURI(CrlDiagramPointerURI, trans)
}

// IsLink returns true if the diagram element is a link
func (diagramElement *CrlDiagramElement) IsLink(trans *core.Transaction) bool {
	return diagramElement.AsCore().IsRefinementOfURI(CrlDiagramLinkURI, trans)
}

// IsNode returns true if the diagram element is a node
func (diagramElement *CrlDiagramElement) IsNode(trans *core.Transaction) bool {
	return diagramElement.AsCore().IsRefinementOfURI(CrlDiagramNodeURI, trans)
}

// SetLinkSource is a convenience function for setting the source concept of a link
func (diagramLink *CrlDiagramLink) SetLinkSource(source *CrlDiagramElement, trans *core.Transaction) {
	sourceReference := diagramLink.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
	if sourceReference != nil {
		sourceReference.SetReferencedConcept(source.AsCore(), core.NoAttribute, trans)
	}
}

// SetLinkTarget is a convenience function for setting the target concept of a link
func (diagramLink *CrlDiagramLink) SetLinkTarget(target *CrlDiagramElement, trans *core.Transaction) {
	targetReference := diagramLink.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
	if targetReference != nil {
		targetReference.SetReferencedConcept(target.AsCore(), core.NoAttribute, trans)
	}
}

// SetNodeHeight is a function on a CrlDiagramNode that sets the height of the diagram node
func (diagramNode *CrlDiagramNode) SetNodeHeight(value float64, trans *core.Transaction) {
	literal := diagramNode.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeHeightURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeWidth is a function on a CrlDiagramNode that sets the width of the diagram node
func (diagramNode *CrlDiagramNode) SetNodeWidth(value float64, trans *core.Transaction) {
	literal := diagramNode.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeWidthURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetNodeX is a function on a CrlDiagramNode that sets the x of the diagram node
func (diagramNode *CrlDiagramNode) SetNodeX(value float64, trans *core.Transaction) {
	literal := diagramNode.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeXURI, trans)
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
	literal := diagramNode.AsCore().GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, trans)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), trans)
}

// SetReferencedModelConcept is a function on a CrlDiagramNode that sets the model element represented by the
// diagram node
func (diagramElement *CrlDiagramElement) SetReferencedModelConcept(el *core.Concept, trans *core.Transaction) {
	if diagramElement == nil {
		return
	}
	reference := diagramElement.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
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
	uOfD.CreateOwnedRefinementOfConcept(crlDiagramAnchoredText, crlDiagramDomain, "SourceMultiplicity", trans, CrlDiagramLinkSourceMultiplicityURI)
	uOfD.CreateOwnedRefinementOfConcept(crlDiagramAnchoredText, crlDiagramDomain, "TargetMultiplicity", trans, CrlDiagramLinkTargetMultiplicityURI)

	// MultiplicityMap
	crlMultiplicityMap, _ := uOfD.NewOwnedElement(crlDiagramDomain, "MultiplicityMap", trans, CrlDiagramMultiplicityMapURI)
	uOfD.NewOwnedReference(crlMultiplicityMap, "MultiplicityReference", trans, CrlDiagramMultiplicityMapMultiplicityReferenceURI)
	uOfD.NewOwnedReference(crlMultiplicityMap, "AnchoredTextReference", trans, CrlDiagramMultiplicityMapAnchoredTextReferenceURI)

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
	uOfD.AddFunction(CrlDiagramMultiplicityMapURI, updateMultiplicityMap)

	crlDiagramDomain.SetIsCoreRecursively(trans)
	return crlDiagramDomain
}

// updateMultiplicityMap updates the anchored text representing multiplicity
func updateMultiplicityMap(changedElement *core.Concept, notification *core.ChangeNotification, trans *core.Transaction) error {
	trans.WriteLockElement(changedElement)
	// core Elements should always be ignored
	if changedElement.GetIsCore(trans) {
		return nil
	}
	if !changedElement.IsRefinementOfURI(CrlDiagramMultiplicityMapURI, trans) {
		return nil
	}
	multiplicityMap := (*CrlDiagramMultiplicityMap)(changedElement)
	modelMultiplicityLiteral := multiplicityMap.GetModelMultiplicity(trans)
	if modelMultiplicityLiteral == nil {
		return nil
	}
	modelMultiplicityValue := modelMultiplicityLiteral.GetLiteralValue(trans)
	diagramMultiplicityValue := multiplicityMap.GetAnchoredText(trans).AsCore().GetLiteralValue(trans)
	if modelMultiplicityValue != diagramMultiplicityValue {
		// We need to determine which changed
		switch notification.GetNatureOfChange() {
		case core.OwnedConceptChanged:
			if notification.GetUnderlyingChange().GetReportingElementID() == multiplicityMap.GetModelMultiplicityReference(trans).ConceptID {
				// it's the model multiplicity
				multiplicityMap.GetAnchoredText(trans).AsCore().SetLiteralValue(modelMultiplicityValue, trans)
			} else if notification.GetUnderlyingChange().GetReportingElementID() == multiplicityMap.GetAnchoredTextReference(trans).ConceptID {
				// it's the diagram multiplicity
				multiplicityMap.GetModelMultiplicity(trans).SetLiteralValue(diagramMultiplicityValue, trans)
			}
		}
	}
	return nil
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
	if underlyingChange != nil && underlyingChange.IsReferenced(diagramElement.AsCore()) {
		return nil
	}

	// There are several notifications of interest here:
	//   - the deletion of the referenced model element
	//   - the label of the referenced model element
	//   - the list of immediate abstractions of the referenced model element.
	// First, determine whether it is the referenced model element that has changed
	diagramElementModelReference := diagramElement.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, trans)
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
			// If the reporting element of the underlying change is the display label, then
			// see if the label needs to change
			underlyingReportingElement := uOfD.GetElement(underlyingChange.GetReportingElementID())
			if IsDisplayLabel(underlyingReportingElement, trans) {
				displayLabelText := underlyingReportingElement.GetLiteralValue(trans)
				if modelElement.GetLabel(trans) != displayLabelText {
					modelElement.SetLabel(displayLabelText, trans)
				}
			}
		case core.ReferencedConceptChanged:
			underlyingReportingElementID := underlyingChange.GetReportingElementID()
			if underlyingReportingElementID == diagramElementModelReference.GetConceptID(trans) {
				// The underlying change is from the model reference
				if IsDiagramNode(diagramElement.AsCore(), trans) {
					currentModelElement := underlyingChange.GetAfterConceptState()
					previousModelElement := underlyingChange.GetBeforeConceptState()
					if currentModelElement != nil && previousModelElement != nil {
						if currentModelElement.ReferencedConceptID == "" && previousModelElement.ReferencedConceptID != "" {
							uOfD.DeleteElement(diagramElement.AsCore(), trans)
						} else {
							diagramElement.updateForModelElementChange(modelElement, trans)
						}
					}
				} else if IsDiagramLink(diagramElement.AsCore(), trans) {
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
							uOfD.DeleteElement(diagramElement.AsCore(), trans)
						} else {
							// Otherwise we update the diagram element
							previousReferencedModelElement := uOfD.GetElement(previousReferencedModelElementID)
							switch modelElement.GetConceptType() {
							case core.Reference:
								if IsDiagramElementPointer(diagramElement.AsCore(), trans) {
									currentReferencedModelElement := modelElement.GetReferencedConcept(trans)
									if previousReferencedModelElement != currentReferencedModelElement {
										if currentReferencedModelElement == nil {
											uOfD.DeleteElement(diagramElement.AsCore(), trans)
										} else {
											newTargetDiagramElement := diagram.GetFirstElementRepresentingConcept(currentReferencedModelElement, trans)
											(*CrlDiagramLink)(diagramElement).SetLinkTarget(newTargetDiagramElement, trans)
										}
									}
								} else if IsDiagramReferenceLink(diagramElement.AsCore(), trans) {
									diagramReferenceLink := (*CrlDiagramLink)(diagramElement)
									diagramElement.updateForModelElementChange(modelElement, trans)
									diagramElement.SetDisplayLabel(modelElement.GetLabel(trans), trans)
									newModelTarget := modelElement.GetReferencedConcept(trans)
									newModelSource := modelElement.GetOwningConcept(trans)
									if newModelSource == nil || newModelTarget == nil {
										uOfD.DeleteElement(diagramElement.AsCore(), trans)
										return nil
									}
									currentDiagramSource := diagramReferenceLink.GetLinkSource(trans)
									currentModelSource := currentDiagramSource.GetReferencedModelConcept(trans)
									currentDiagramTarget := diagramReferenceLink.GetLinkTarget(trans)
									currentModelTarget := currentDiagramTarget.GetReferencedModelConcept(trans)
									if currentModelSource != newModelSource {
										newDiagramSource := diagram.GetFirstElementRepresentingConcept(newModelSource, trans)
										if newDiagramSource == nil {
											uOfD.DeleteElement(diagramElement.AsCore(), trans)
											return nil
										}
										diagramReferenceLink.SetLinkSource(newDiagramSource, trans)
									}
									if currentModelTarget != newModelTarget {
										newDiagramTarget := diagram.GetFirstElementRepresentingConcept(newModelTarget, trans)
										if newDiagramTarget == nil {
											uOfD.DeleteElement(diagramElement.AsCore(), trans)
											return nil
										}
										diagramReferenceLink.SetLinkTarget(newDiagramTarget, trans)
									}
								}
							case core.Refinement:
								refinement := modelElement
								if IsDiagramPointer(diagramElement.AsCore(), trans) {
									diagramPointer := (*CrlDiagramLink)(diagramElement)
									var newTargetModelElement *core.Concept
									if IsDiagramAbstractPointer(diagramPointer.AsCore(), trans) {
										newTargetModelElement = refinement.GetAbstractConcept(trans)
									} else if IsDiagramRefinedPointer(diagramPointer.AsCore(), trans) {
										newTargetModelElement = refinement.GetRefinedConcept(trans)
									} else if IsDiagramOwnerPointer(diagramPointer.AsCore(), trans) {
										newTargetModelElement = refinement.GetOwningConcept(trans)
									}
									if previousReferencedModelElement != newTargetModelElement {
										if newTargetModelElement == nil {
											uOfD.DeleteElement(diagramElement.AsCore(), trans)
										} else {
											newTargetDiagramElement := diagram.GetFirstElementRepresentingConcept(newTargetModelElement, trans)
											diagramPointer.SetLinkTarget(newTargetDiagramElement, trans)
										}
									}
								} else if IsDiagramRefinementLink(diagramElement.AsCore(), trans) {
									refinementLink := (*CrlDiagramLink)(diagramElement)
									diagramElement.updateForModelElementChange(modelElement, trans)
									diagramElement.SetDisplayLabel(refinement.GetLabel(trans), trans)
									newModelTarget := refinement.GetAbstractConcept(trans)
									newModelSource := refinement.GetRefinedConcept(trans)
									if newModelTarget == nil || newModelSource == nil {
										uOfD.DeleteElement(diagramElement.AsCore(), trans)
										return nil
									}
									currentDiagramTarget := refinementLink.GetLinkTarget(trans)
									currentModelTarget := currentDiagramTarget.GetReferencedModelConcept(trans)
									currentDiagramSource := refinementLink.GetLinkSource(trans)
									currentModelSource := currentDiagramSource.GetReferencedModelConcept(trans)
									if currentModelTarget != newModelTarget {
										newDiagramTarget := diagram.GetFirstElementRepresentingConcept(newModelTarget, trans)
										if newDiagramTarget == nil {
											uOfD.DeleteElement(diagramElement.AsCore(), trans)
											return nil
										}
										refinementLink.SetLinkTarget(newDiagramTarget, trans)
									}
									if currentModelSource != newModelSource {
										newDiagramSource := diagram.GetFirstElementRepresentingConcept(newModelSource, trans)
										if newDiagramSource == nil {
											uOfD.DeleteElement(diagramElement.AsCore(), trans)
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
				if IsDiagramLink(diagramElement.AsCore(), trans) {
					diagramLink := (*CrlDiagramLink)(diagramElement)
					// If this is a diagram link and the underlying reporting element is either its source reference or its target reference
					// and the referenced element is now nil, we need to delete the link
					underlyingReportingElementIsSourceReference := diagramLink.GetLinkSourceReference(trans).GetConceptID(trans) == underlyingReportingElementID
					underlyingReportingElementIsTargetReference := diagramLink.GetLinkTargetReference(trans).GetConceptID(trans) == underlyingReportingElementID
					if (underlyingReportingElementIsSourceReference || underlyingReportingElementIsTargetReference) &&
						underlyingChange.GetAfterConceptState().ReferencedConceptID == "" {
						uOfD.DeleteElement(diagramElement.AsCore(), trans)
					} else {
						switch underlyingChange.GetNatureOfChange() {
						case core.ReferencedConceptChanged:
							if underlyingReportingElementIsTargetReference {
								// If the link's target has changed, we need to update the underlying model element to reflect the change.
								// Note that if the target is now null, the preceeding clause will have deleted the element
								targetDiagramElement := (*CrlDiagramElement)(uOfD.GetElement(underlyingChange.GetAfterConceptState().ReferencedConceptID))
								targetModelElement := targetDiagramElement.GetReferencedModelConcept(trans)
								if IsDiagramOwnerPointer(diagramElement.AsCore(), trans) {
									modelElement.SetOwningConcept(targetModelElement, trans)
								} else {
									switch modelElement.GetConceptType() {
									case core.Reference:
										// Setting the referenced concepts requires knowledge of what is being referenced
										targetAttribute := core.NoAttribute
										if IsDiagramElementPointer(targetDiagramElement.AsCore(), trans) {
											targetAttribute = core.ReferencedConceptID
										} else if IsDiagramOwnerPointer(targetDiagramElement.AsCore(), trans) {
											targetAttribute = core.OwningConceptID
										} else if IsDiagramAbstractPointer(targetDiagramElement.AsCore(), trans) {
											targetAttribute = core.AbstractConceptID
										} else if IsDiagramRefinedPointer(targetDiagramElement.AsCore(), trans) {
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
										if IsDiagramReferenceLink(diagramElement.AsCore(), trans) {
											modelElement.SetOwningConcept(sourceModelElement, trans)
										}
									case core.Refinement:
										if IsDiagramRefinementLink(diagramElement.AsCore(), trans) {
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
			uOfD.DeleteElement(diagramElement.AsCore(), trans)
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
	diagram := diagramPointer.AsCrlDiagramElement().GetDiagram(trans)
	modelElement := diagramPointer.AsCrlDiagramElement().GetReferencedModelConcept(trans)
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
							uOfD.DeleteElement(diagramPointer.AsCore(), trans)
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
		sourceReference := diagramPointer.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, trans)
		targetReference := diagramPointer.AsCore().GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, trans)
		if reportingElement == sourceReference || reportingElement == targetReference {
			underlyingNotification := notification.GetUnderlyingChange()
			switch underlyingNotification.GetNatureOfChange() {
			case core.ReferencedConceptChanged:
				switch reportingElement.GetConceptType() {
				case core.Reference:
					if reportingElement.GetReferencedConcept(trans) == nil {
						uOfD.DeleteElement(diagramPointer.AsCore(), trans)
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
		if modelElementLabel != diagramElement.AsCore().GetLabel(trans) {
			newLabel := modelElementLabel
			if IsDiagramPointer(diagramElement.AsCore(), trans) {
				if IsDiagramOwnerPointer(diagramElement.AsCore(), trans) {
					newLabel = newLabel + " Owner Pointer"
				} else if IsDiagramAbstractPointer(diagramElement.AsCore(), trans) {
					newLabel = newLabel + " Abstract Pointer"
				} else if IsDiagramRefinedPointer(diagramElement.AsCore(), trans) {
					newLabel = newLabel + " Refined Pointer"
				} else if IsDiagramElementPointer(diagramElement.AsCore(), trans) {
					newLabel = newLabel + " Referenced Concept Pointer"
				}
			}
			diagramElement.AsCore().SetLabel(newLabel, trans)
			if !IsDiagramPointer(diagramElement.AsCore(), trans) {
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
	displayLabel := diagramNode.AsDiagramElement().GetDisplayLabel(trans)
	displayLabelBounds, _ := font.BoundString(go12PtBoldFace, displayLabel)
	displayLabelMaxHeight := Int26_6ToFloat(displayLabelBounds.Max.Y)
	displayLabelMaxWidth := Int26_6ToFloat(displayLabelBounds.Max.X)
	displayLabelMinHeight := Int26_6ToFloat(displayLabelBounds.Min.Y)
	displayLabelMinWidth := Int26_6ToFloat(displayLabelBounds.Min.X)
	displayLabelHeight := displayLabelMaxHeight - displayLabelMinHeight
	displayLabelWidth := displayLabelMaxWidth - displayLabelMinWidth
	abstractionDisplayLabel := diagramNode.AsDiagramElement().GetAbstractionDisplayLabel(trans)
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
