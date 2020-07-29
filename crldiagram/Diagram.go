// Package crldiagram defines the CoreDiagram concept space. This is a pre-defined concept space (hence the term "core") that is, itself,
// represented as a CRLElement and identified with the CoreDiagramURI. This concept space contains the prototypes of all Elements used to construct CrlDiagrams.
// Included are:
// 	CrlDiagram: the diagram itself
// 	CrlDiagramNode: a node in the diagram
// 	CrlDiagramLink: a link in the diagram
//  CrlDiagramPointer: a pointer shown as a link in the diagram
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
package crldiagram

import (
	// "github.com/pkg/errors"
	"log"
	"math"
	"strconv"

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
var CrlDiagramPrefix = "http://activeCrl.com/corediagram/"

// CrlDiagramConceptSpaceURI identifies concept space containing all concepts related to the CrlDiagram
var CrlDiagramConceptSpaceURI = CrlDiagramPrefix + "CoreDiagram"

// CrlDiagramURI identifies the CrlDiagram concept
var CrlDiagramURI = CrlDiagramConceptSpaceURI + "/" + "CrlDiagram"

// CrlDiagramWidthURI identifies the CrlDiagramWidth concept
var CrlDiagramWidthURI = CrlDiagramURI + "/" + "Width"

// CrlDiagramHeightURI identifies the CrlDiagramHeight concept
var CrlDiagramHeightURI = CrlDiagramURI + "/" + "Height"

// CrlDiagramElementURI identifies the CrlDiagramElement concept
var CrlDiagramElementURI = CrlDiagramConceptSpaceURI + "/" + "CrlDiagramElement"

// CrlDiagramElementModelReferenceURI identifies the reference to the model element represented by the node
var CrlDiagramElementModelReferenceURI = CrlDiagramElementURI + "/" + "ModelReference"

// CrlDiagramElementDisplayLabelURI identifies the display label concept to be used when displaying the node
var CrlDiagramElementDisplayLabelURI = CrlDiagramElementURI + "/" + "DisplayLabel"

// CrlDiagramElementAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the node
var CrlDiagramElementAbstractionDisplayLabelURI = CrlDiagramElementURI + "/" + "AbstractionDisplayLabel"

// CrlDiagramLinkURI identifies the CrlDiagramLink concept
var CrlDiagramLinkURI = CrlDiagramConceptSpaceURI + "/" + "CrlDiagramLink"

// CrlDiagramLinkSourceURI identifies the concept that is the source of the link
var CrlDiagramLinkSourceURI = CrlDiagramLinkURI + "/" + "Source"

// CrlDiagramLinkTargetURI identifies the concept that is the target of the link
var CrlDiagramLinkTargetURI = CrlDiagramLinkURI + "/" + "Target"

// CrlDiagramNodeURI identifies the CrlDiagramNode concept
var CrlDiagramNodeURI = CrlDiagramConceptSpaceURI + "/" + "CrlDiagramNode"

// CrlDiagramNodeModelReferenceURI identifies the reference to the model element represented by the node
var CrlDiagramNodeModelReferenceURI = CrlDiagramNodeURI + "/" + "ModelReference"

// CrlDiagramNodeDisplayLabelURI identifies the display label concept to be used when displaying the node
var CrlDiagramNodeDisplayLabelURI = CrlDiagramNodeURI + "/" + "DisplayLabel"

// CrlDiagramNodeAbstractionDisplayLabelURI identifies the abstraction display label concept to be used when displaying the node
var CrlDiagramNodeAbstractionDisplayLabelURI = CrlDiagramNodeURI + "/" + "AbstractionDisplayLabel"

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
var CrlDiagramPointerURI = CrlDiagramConceptSpaceURI + "/" + "Pointer"

// CrlDiagramAbstractPointerURI identifies the Abstract of an Element represented as a link
var CrlDiagramAbstractPointerURI = CrlDiagramConceptSpaceURI + "/" + "AbstractPointer"

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
var CrlDiagramElementPointerURI = CrlDiagramConceptSpaceURI + "/" + "ElementPointer"

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
var CrlDiagramOwnerPointerURI = CrlDiagramConceptSpaceURI + "/" + "OwnerPointer"

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
var CrlDiagramRefinedPointerURI = CrlDiagramConceptSpaceURI + "/" + "RefinedPointer"

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
var CrlDiagramReferenceLinkURI = CrlDiagramConceptSpaceURI + "/" + "ReferenceLink"

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
var CrlDiagramRefinementLinkURI = CrlDiagramConceptSpaceURI + "/" + "RefinementLink"

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
func GetDisplayLabel(diagramElement core.Element, hl *core.HeldLocks) string {
	if diagramElement == nil {
		return ""
	}
	displayLabelLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, hl)
	if displayLabelLiteral != nil {
		return displayLabelLiteral.GetLiteralValue(hl)
	}
	return ""
}

// GetAbstractionDisplayLabel is a convenience function for getting the DisplayLabel value of a node's position
func GetAbstractionDisplayLabel(diagramElement core.Element, hl *core.HeldLocks) string {
	if diagramElement == nil {
		return ""
	}
	abstractionDisplayLabelLiteral := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, hl)
	if abstractionDisplayLabelLiteral != nil {
		return abstractionDisplayLabelLiteral.GetLiteralValue(hl)
	}
	return ""
}

// GetFirstElementRepresentingConcept returns the first diagram element that represents the indicated concept
func GetFirstElementRepresentingConcept(diagram core.Element, concept core.Element, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConcept called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, hl) {
		if GetReferencedModelElement(el, hl) == concept && !el.IsRefinementOfURI(CrlDiagramPointerURI, hl) {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptID returns the first diagram element that represents the indicated concept
func GetFirstElementRepresentingConceptID(diagram core.Element, conceptID string, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptID called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, hl) {
		if GetReferencedModelElement(el, hl).GetConceptID(hl) == conceptID && !el.IsRefinementOfURI(CrlDiagramPointerURI, hl) {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptOwnerPointer returns the first diagram element that represents the indicated concept's OwnerPointer
func GetFirstElementRepresentingConceptOwnerPointer(diagram core.Element, concept core.Element, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptOwnerPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, hl) {
		if GetReferencedModelElement(el, hl) == concept {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDOwnerPointer returns the first diagram element that represents the indicated concept's OwnerPointer
func GetFirstElementRepresentingConceptIDOwnerPointer(diagram core.Element, conceptID string, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptIDOwnerPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, hl) {
		if GetReferencedModelElement(el, hl).GetConceptID(hl) == conceptID {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptElementPointer returns the first diagram element that represents the indicated concept's ElementPointer
func GetFirstElementRepresentingConceptElementPointer(diagram core.Element, concept core.Reference, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptElementPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, hl) {
		if GetReferencedModelElement(el, hl) == concept {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDElementPointer returns the first diagram element that represents the indicated concept's ElementPointer
func GetFirstElementRepresentingConceptIDElementPointer(diagram core.Element, conceptID string, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptIDElementPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, hl) {
		if GetReferencedModelElement(el, hl).GetConceptID(hl) == conceptID {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptAbstractPointer returns the first diagram element that represents the indicated concept's AbstractPointer
func GetFirstElementRepresentingConceptAbstractPointer(diagram core.Element, concept core.Refinement, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptAbstractPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramAbstractPointerURI, hl) {
		if GetReferencedModelElement(el, hl) == concept {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDAbstractPointer returns the first diagram element that represents the indicated concept's AbstractPointer
func GetFirstElementRepresentingConceptIDAbstractPointer(diagram core.Element, conceptID string, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptIDAbstractPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramAbstractPointerURI, hl) {
		if GetReferencedModelElement(el, hl).GetConceptID(hl) == conceptID {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptRefinedPointer returns the first diagram element that represents the indicated concept's RefinedPointer
func GetFirstElementRepresentingConceptRefinedPointer(diagram core.Element, concept core.Refinement, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptRefinedPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramRefinedPointerURI, hl) {
		if GetReferencedModelElement(el, hl) == concept {
			return el
		}
	}
	return nil
}

// GetFirstElementRepresentingConceptIDRefinedPointer returns the first diagram element that represents the indicated concept's RefinedPointer
func GetFirstElementRepresentingConceptIDRefinedPointer(diagram core.Element, conceptID string, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConceptIDRefinedPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramRefinedPointerURI, hl) {
		if GetReferencedModelElement(el, hl).GetConceptID(hl) == conceptID {
			return el
		}
	}
	return nil
}

// GetLinkSource is a convenience function for getting the source concept of a link
func GetLinkSource(diagramLink core.Element, hl *core.HeldLocks) core.Element {
	if diagramLink == nil {
		return nil
	}
	sourceReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, hl)
	if sourceReference != nil {
		return sourceReference.GetReferencedConcept(hl)
	}
	return nil
}

// GetLinkTarget is a convenience function for getting the target concept of a link
func GetLinkTarget(diagramLink core.Element, hl *core.HeldLocks) core.Element {
	if diagramLink == nil {
		return nil
	}
	targetReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, hl)
	if targetReference != nil {
		return targetReference.GetReferencedConcept(hl)
	}
	return nil
}

// GetNodeHeight is a convenience function for getting the Height value of a node's position
func GetNodeHeight(diagramNode core.Element, hl *core.HeldLocks) float64 {
	if diagramNode == nil {
		return 0.0
	}
	heightLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeHeightURI, hl)
	if heightLiteral != nil {
		value := heightLiteral.GetLiteralValue(hl)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetNodeWidth is a convenience function for getting the Width value of a node's position
func GetNodeWidth(diagramNode core.Element, hl *core.HeldLocks) float64 {
	if diagramNode == nil {
		return 0.0
	}
	widthLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeWidthURI, hl)
	if widthLiteral != nil {
		value := widthLiteral.GetLiteralValue(hl)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetNodeX is a convenience function for getting the X value of a node's position
func GetNodeX(diagramNode core.Element, hl *core.HeldLocks) float64 {
	if diagramNode == nil {
		return 0.0
	}
	xLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeXURI, hl)
	if xLiteral != nil {
		value := xLiteral.GetLiteralValue(hl)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetNodeY is a convenience function for getting the X value of a node's position
func GetNodeY(diagramNode core.Element, hl *core.HeldLocks) float64 {
	if diagramNode == nil {
		return 0.0
	}
	yLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, hl)
	if yLiteral != nil {
		value := yLiteral.GetLiteralValue(hl)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetDisplayLabelYOffset is a convenience function for getting the Display Label's Y offset within the node
func GetDisplayLabelYOffset(diagramNode core.Element, hl *core.HeldLocks) float64 {
	if diagramNode == nil {
		return 0.0
	}
	yOffsetLiteral := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeDisplayLabelYOffsetURI, hl)
	if yOffsetLiteral != nil {
		value := yOffsetLiteral.GetLiteralValue(hl)
		numericValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return numericValue
		}
	}
	return 0.0
}

// GetOwnerPointer returns the ownerPointer for the concept if one exists
func GetOwnerPointer(diagram core.Element, concept core.Element, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetOwnerPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, hl) {
		if GetReferencedModelElement(el, hl) == concept {
			return el
		}
	}
	return nil
}

// GetElementPointer returns the elementPointer for the concept if one exists
func GetElementPointer(diagram core.Element, concept core.Element, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetElementPointer called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementPointerURI, hl) {
		if GetReferencedModelElement(el, hl) == concept {
			return el
		}
	}
	return nil
}

// GetReferencedModelElement is a function on a CrlDiagramNode that returns the model element represented by the
// diagram node
func GetReferencedModelElement(diagramElement core.Element, hl *core.HeldLocks) core.Element {
	if diagramElement == nil {
		return nil
	}
	reference := diagramElement.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)
	if reference != nil {
		return reference.GetReferencedConcept(hl)
	}
	return nil
}

func init() {
	var err error

	// Set up fonts and faces
	goRegularFont, err = truetype.Parse(goregular.TTF)
	if err != nil {
		log.Printf(err.Error())
	}
	goBoldFont, err = truetype.Parse(gobold.TTF)
	if err != nil {
		log.Printf(err.Error())
	}

	goItalicFont, err := truetype.Parse(goitalic.TTF)
	if err != nil {
		log.Printf(err.Error())
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
func IsDiagram(el core.Element, hl *core.HeldLocks) bool {
	switch el.(type) {
	case core.Element:
		return el.IsRefinementOfURI(CrlDiagramURI, hl)
	}
	return false
}

// IsDiagramAbstractPointer returns true if the supplied element is a CrlDiagramAbstractPointer
func IsDiagramAbstractPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramAbstractPointerURI, hl)
}

// IsDiagramElement returns true if the supplied element is a CrlDiagramElement
func IsDiagramElement(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramElementURI, hl)
}

// IsDiagramElementPointer returns true if the supplied element is a CrlDiagramElementPointer
func IsDiagramElementPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramElementPointerURI, hl)
}

// IsDiagramLink returns true if the supplied element is a CrlDiagramLink
func IsDiagramLink(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramLinkURI, hl)
}

// IsDiagramNode returns true if the supplied element is a CrlDiagramNode
func IsDiagramNode(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramNodeURI, hl)
}

// IsDiagramOwnerPointer returns true if the supplied element is a CrlDiagramOwnerPointer
func IsDiagramOwnerPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramOwnerPointerURI, hl)
}

// IsDiagramPointer returns true if the supplied element is a CrlDiagramPointer
func IsDiagramPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramPointerURI, hl)
}

// IsDiagramRefinedPointer returns true if the supplied element is a CrlDiagramRefinedPointer
func IsDiagramRefinedPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinedPointerURI, hl)
}

// IsDiagramReferenceLink returns true if the supplied element is a CrlDiagramReferenceLink
func IsDiagramReferenceLink(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramReferenceLinkURI, hl)
}

// IsDiagramRefinementLink returns true if the supplied element is a CrlDiagramRefinementLink
func IsDiagramRefinementLink(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinementLinkURI, hl)
}

// IsModelReference returns true if the supplied element is a ModelReference
func IsModelReference(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramElementModelReferenceURI, hl)
}

// NewDiagram creates a new diagram
func NewDiagram(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramURI, hl)
}

// NewDiagramReferenceLink creates a new diagram link to represent a reference
func NewDiagramReferenceLink(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramReferenceLinkURI, hl)
}

// NewDiagramRefinementLink creates a new diagram link
func NewDiagramRefinementLink(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramRefinementLinkURI, hl)
}

// NewDiagramNode creates a new diagram node
func NewDiagramNode(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
}

// NewDiagramOwnerPointer creates a new DiagramOwnerPointer
func NewDiagramOwnerPointer(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramOwnerPointerURI, hl)
}

// NewDiagramElementPointer creates a new DiagramElementPointer
func NewDiagramElementPointer(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramElementPointerURI, hl)
}

// NewDiagramAbstractPointer creates a new DiagramAbstractPointer
func NewDiagramAbstractPointer(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramAbstractPointerURI, hl)
}

// NewDiagramRefinedPointer creates a new DiagramRefinedPointer
func NewDiagramRefinedPointer(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramRefinedPointerURI, hl)
}

// SetAbstractionDisplayLabel is a function on a CrlDiagramNode that sets the display label of the diagram node
func SetAbstractionDisplayLabel(diagramElement core.Element, value string, hl *core.HeldLocks) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementAbstractionDisplayLabelURI, hl)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(value, hl)
	updateNodeSize(diagramElement, hl)
}

// SetDisplayLabel is a function on a CrlDiagramNode that sets the display label of the diagram element.
// If the diagram element is a pointer, the value is ignored and the label is set to the empty string
func SetDisplayLabel(diagramElement core.Element, value string, hl *core.HeldLocks) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, hl)
	if literal == nil {
		return
	}
	if IsDiagramPointer(diagramElement, hl) {
		literal.SetLiteralValue("", hl)
	} else {
		literal.SetLiteralValue(value, hl)
	}
	literal.SetLiteralValue(value, hl)
	updateNodeSize(diagramElement, hl)
}

// SetLinkSource is a convenience function for setting the source concept of a link
func SetLinkSource(diagramLink core.Element, source core.Element, hl *core.HeldLocks) {
	if diagramLink == nil {
		return
	}
	sourceReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, hl)
	if sourceReference != nil {
		sourceReference.SetReferencedConcept(source, hl)
	}
}

// SetLinkTarget is a convenience function for setting the target concept of a link
func SetLinkTarget(diagramLink core.Element, target core.Element, hl *core.HeldLocks) {
	if diagramLink == nil {
		return
	}
	targetReference := diagramLink.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, hl)
	if targetReference != nil {
		targetReference.SetReferencedConcept(target, hl)
	}
}

// SetNodeHeight is a function on a CrlDiagramNode that sets the height of the diagram node
func SetNodeHeight(diagramNode core.Element, value float64, hl *core.HeldLocks) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeHeightURI, hl)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), hl)
}

// SetNodeWidth is a function on a CrlDiagramNode that sets the width of the diagram node
func SetNodeWidth(diagramNode core.Element, value float64, hl *core.HeldLocks) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeWidthURI, hl)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), hl)
}

// SetNodeX is a function on a CrlDiagramNode that sets the x of the diagram node
func SetNodeX(diagramNode core.Element, value float64, hl *core.HeldLocks) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeXURI, hl)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), hl)
}

// SetNodeY is a function on a CrlDiagramNode that sets the y of the diagram node
func SetNodeY(diagramNode core.Element, value float64, hl *core.HeldLocks) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeYURI, hl)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), hl)
}

// SetNodeDisplayLabelYOffset is a function on a CrlDiagramNode that sets the y offset of the display label within the node
func SetNodeDisplayLabelYOffset(diagramNode core.Element, value float64, hl *core.HeldLocks) {
	if diagramNode == nil {
		return
	}
	literal := diagramNode.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramNodeDisplayLabelYOffsetURI, hl)
	if literal == nil {
		return
	}
	literal.SetLiteralValue(strconv.FormatFloat(value, 'f', -1, 64), hl)
}

// SetReferencedModelElement is a function on a CrlDiagramNode that sets the model element represented by the
// diagram node
func SetReferencedModelElement(diagramElement core.Element, el core.Element, hl *core.HeldLocks) {
	if diagramElement == nil {
		return
	}
	reference := diagramElement.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)
	if reference == nil {
		return
	}
	reference.SetReferencedConcept(el, hl)
}

// BuildCrlDiagramConceptSpace builds the CrlDiagram concept space and adds it to the uOfD
func BuildCrlDiagramConceptSpace(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// CrlDiagramConceptSpace
	crlDiagramConceptSpace, _ := uOfD.NewElement(hl, CrlDiagramConceptSpaceURI)
	crlDiagramConceptSpace.SetLabel("CrlDiagramConceptSpace", hl)

	//
	// CrlDiagram
	//
	crlDiagram, _ := uOfD.NewElement(hl, CrlDiagramURI)
	crlDiagram.SetLabel("CrlDiagram", hl)
	crlDiagram.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramWidth, _ := uOfD.NewLiteral(hl, CrlDiagramWidthURI)
	crlDiagramWidth.SetLabel("Width", hl)
	crlDiagramWidth.SetOwningConcept(crlDiagram, hl)

	crlDiagramHeight, _ := uOfD.NewLiteral(hl, CrlDiagramHeightURI)
	crlDiagramHeight.SetLabel("Height", hl)
	crlDiagramHeight.SetOwningConcept(crlDiagram, hl)

	//
	// CrlDiagramElement
	//
	crlDiagramElement, _ := uOfD.NewElement(hl, CrlDiagramElementURI)
	crlDiagramElement.SetLabel("CrlDiagramElement", hl)
	crlDiagramElement.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramElementModelReference, _ := uOfD.NewReference(hl, CrlDiagramElementModelReferenceURI)
	crlDiagramElementModelReference.SetLabel("ModelReference", hl)
	crlDiagramElementModelReference.SetOwningConcept(crlDiagramElement, hl)

	crlDiagramElementDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementDisplayLabelURI)
	crlDiagramElementDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramElementDisplayLabel.SetOwningConcept(crlDiagramElement, hl)

	crlDiagramElementAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementAbstractionDisplayLabelURI)
	crlDiagramElementAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramElementAbstractionDisplayLabel.SetOwningConcept(crlDiagramElement, hl)

	//
	// CrlDiagramNode
	//
	crlDiagramNode, _ := uOfD.NewElement(hl, CrlDiagramNodeURI)
	crlDiagramNode.SetLabel("CrlDiagramNode", hl)
	crlDiagramNode.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramNodeRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramNodeRefinement.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeRefinement.SetAbstractConcept(crlDiagramElement, hl)
	crlDiagramNodeRefinement.SetRefinedConcept(crlDiagramNode, hl)

	crlDiagramNodeModelReference, _ := uOfD.NewReference(hl, CrlDiagramNodeModelReferenceURI)
	crlDiagramNodeModelReference.SetLabel("ModelReference", hl)
	crlDiagramNodeModelReference.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeModelReferenceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramNodeModelReferenceRefinement.SetOwningConcept(crlDiagramNodeModelReference, hl)
	crlDiagramNodeModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, hl)
	crlDiagramNodeModelReferenceRefinement.SetRefinedConcept(crlDiagramNodeModelReference, hl)

	crlDiagramNodeDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramNodeDisplayLabelURI)
	crlDiagramNodeDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramNodeDisplayLabel.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramNodeDisplayLabelRefinement.SetOwningConcept(crlDiagramNodeDisplayLabel, hl)
	crlDiagramNodeDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, hl)
	crlDiagramNodeDisplayLabelRefinement.SetRefinedConcept(crlDiagramNodeDisplayLabel, hl)

	crlDiagramNodeAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramNodeAbstractionDisplayLabelURI)
	crlDiagramNodeAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramNodeAbstractionDisplayLabel.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramNodeAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramNodeAbstractionDisplayLabel, hl)
	crlDiagramNodeAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, hl)
	crlDiagramNodeAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramNodeAbstractionDisplayLabel, hl)

	crlDiagramNodeX, _ := uOfD.NewLiteral(hl, CrlDiagramNodeXURI)
	crlDiagramNodeX.SetLabel("X", hl)
	crlDiagramNodeX.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeY, _ := uOfD.NewLiteral(hl, CrlDiagramNodeYURI)
	crlDiagramNodeY.SetLabel("Y", hl)
	crlDiagramNodeY.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeHeight, _ := uOfD.NewLiteral(hl, CrlDiagramNodeHeightURI)
	crlDiagramNodeHeight.SetLabel("Height", hl)
	crlDiagramNodeHeight.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeWidth, _ := uOfD.NewLiteral(hl, CrlDiagramNodeWidthURI)
	crlDiagramNodeWidth.SetLabel("Width", hl)
	crlDiagramNodeWidth.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeDisplayLabelYOffset, _ := uOfD.NewLiteral(hl, CrlDiagramNodeDisplayLabelYOffsetURI)
	crlDiagramNodeDisplayLabelYOffset.SetLabel("DisplayLabelYOffset", hl)
	crlDiagramNodeDisplayLabelYOffset.SetOwningConcept(crlDiagramNode, hl)

	//
	// CrlDiagramLink
	//
	crlDiagramLink, _ := uOfD.NewElement(hl, CrlDiagramLinkURI)
	crlDiagramLink.SetLabel("CrlDiagramLink", hl)
	crlDiagramLink.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramLinkRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramLinkRefinement.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkRefinement.SetAbstractConcept(crlDiagramElement, hl)
	crlDiagramLinkRefinement.SetRefinedConcept(crlDiagramLink, hl)

	crlDiagramLinkSource, _ := uOfD.NewReference(hl, CrlDiagramLinkSourceURI)
	crlDiagramLinkSource.SetLabel("Source", hl)
	crlDiagramLinkSource.SetOwningConcept(crlDiagramLink, hl)

	crlDiagramLinkTarget, _ := uOfD.NewReference(hl, CrlDiagramLinkTargetURI)
	crlDiagramLinkTarget.SetLabel("Target", hl)
	crlDiagramLinkTarget.SetOwningConcept(crlDiagramLink, hl)

	//
	// Pointer
	//
	crlDiagramPointer, _ := uOfD.NewElement(hl, CrlDiagramPointerURI)
	crlDiagramPointer.SetLabel("Pointer", hl)
	crlDiagramPointer.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramPointerRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramPointerRefinement.SetOwningConcept(crlDiagramPointer, hl)
	crlDiagramPointerRefinement.SetAbstractConcept(crlDiagramLink, hl)
	crlDiagramPointerRefinement.SetRefinedConcept(crlDiagramPointer, hl)

	//
	// AbstractPointer
	//
	crlDiagramAbstractPointer, _ := uOfD.NewElement(hl, CrlDiagramAbstractPointerURI)
	crlDiagramAbstractPointer.SetLabel("AbstractPointer", hl)
	crlDiagramAbstractPointer.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramAbstractPointerRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramAbstractPointerRefinement.SetOwningConcept(crlDiagramAbstractPointer, hl)
	crlDiagramAbstractPointerRefinement.SetAbstractConcept(crlDiagramPointer, hl)
	crlDiagramAbstractPointerRefinement.SetRefinedConcept(crlDiagramAbstractPointer, hl)

	crlDiagramAbstractPointerModelReference, _ := uOfD.NewReference(hl, CrlDiagramAbstractPointerModelReferenceURI)
	crlDiagramAbstractPointerModelReference.SetLabel("ModelReference", hl)
	crlDiagramAbstractPointerModelReference.SetOwningConcept(crlDiagramAbstractPointer, hl)

	crlDiagramAbstractPointerModelReferenceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramAbstractPointerModelReferenceRefinement.SetOwningConcept(crlDiagramAbstractPointerModelReference, hl)
	crlDiagramAbstractPointerModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, hl)
	crlDiagramAbstractPointerModelReferenceRefinement.SetRefinedConcept(crlDiagramAbstractPointerModelReference, hl)

	crlDiagramAbstractPointerDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramAbstractPointerDisplayLabelURI)
	crlDiagramAbstractPointerDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramAbstractPointerDisplayLabel.SetOwningConcept(crlDiagramAbstractPointer, hl)

	crlDiagramAbstractPointerDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramAbstractPointerDisplayLabelRefinement.SetOwningConcept(crlDiagramAbstractPointerDisplayLabel, hl)
	crlDiagramAbstractPointerDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, hl)
	crlDiagramAbstractPointerDisplayLabelRefinement.SetRefinedConcept(crlDiagramAbstractPointerDisplayLabel, hl)

	crlDiagramAbstractPointerAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramAbstractPointerAbstractionDisplayLabelURI)
	crlDiagramAbstractPointerAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramAbstractPointerAbstractionDisplayLabel.SetOwningConcept(crlDiagramAbstractPointer, hl)

	crlDiagramAbstractPointerAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramAbstractPointerAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramAbstractPointerAbstractionDisplayLabel, hl)
	crlDiagramAbstractPointerAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, hl)
	crlDiagramAbstractPointerAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramAbstractPointerAbstractionDisplayLabel, hl)

	crlDiagramAbstractPointerSource, _ := uOfD.NewReference(hl, CrlDiagramAbstractPointerSourceURI)
	crlDiagramAbstractPointerSource.SetLabel("Source", hl)
	crlDiagramAbstractPointerSource.SetOwningConcept(crlDiagramAbstractPointer, hl)

	crlDiagramAbstractPointerSourceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramAbstractPointerSourceRefinement.SetOwningConcept(crlDiagramAbstractPointerSource, hl)
	crlDiagramAbstractPointerSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, hl)
	crlDiagramAbstractPointerSourceRefinement.SetRefinedConcept(crlDiagramAbstractPointerSource, hl)

	crlDiagramAbstractPointerTarget, _ := uOfD.NewReference(hl, CrlDiagramAbstractPointerTargetURI)
	crlDiagramAbstractPointerTarget.SetLabel("Target", hl)
	crlDiagramAbstractPointerTarget.SetOwningConcept(crlDiagramAbstractPointer, hl)

	crlDiagramAbstractPointerTargetRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramAbstractPointerTargetRefinement.SetOwningConcept(crlDiagramAbstractPointerTarget, hl)
	crlDiagramAbstractPointerTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, hl)
	crlDiagramAbstractPointerTargetRefinement.SetRefinedConcept(crlDiagramAbstractPointerTarget, hl)

	//
	// ElementPointer
	//
	crlDiagramElementPointer, _ := uOfD.NewElement(hl, CrlDiagramElementPointerURI)
	crlDiagramElementPointer.SetLabel("ElementPointer", hl)
	crlDiagramElementPointer.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramElementPointerRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramElementPointerRefinement.SetOwningConcept(crlDiagramElementPointer, hl)
	crlDiagramElementPointerRefinement.SetAbstractConcept(crlDiagramPointer, hl)
	crlDiagramElementPointerRefinement.SetRefinedConcept(crlDiagramElementPointer, hl)

	crlDiagramElementPointerModelReference, _ := uOfD.NewReference(hl, CrlDiagramElementPointerModelReferenceURI)
	crlDiagramElementPointerModelReference.SetLabel("ModelReference", hl)
	crlDiagramElementPointerModelReference.SetOwningConcept(crlDiagramElementPointer, hl)

	crlDiagramElementPointerModelReferenceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramElementPointerModelReferenceRefinement.SetOwningConcept(crlDiagramElementPointerModelReference, hl)
	crlDiagramElementPointerModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, hl)
	crlDiagramElementPointerModelReferenceRefinement.SetRefinedConcept(crlDiagramElementPointerModelReference, hl)

	crlDiagramElementPointerDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementPointerDisplayLabelURI)
	crlDiagramElementPointerDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramElementPointerDisplayLabel.SetOwningConcept(crlDiagramElementPointer, hl)

	crlDiagramElementPointerDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramElementPointerDisplayLabelRefinement.SetOwningConcept(crlDiagramElementPointerDisplayLabel, hl)
	crlDiagramElementPointerDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, hl)
	crlDiagramElementPointerDisplayLabelRefinement.SetRefinedConcept(crlDiagramElementPointerDisplayLabel, hl)

	crlDiagramElementPointerAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementPointerAbstractionDisplayLabelURI)
	crlDiagramElementPointerAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramElementPointerAbstractionDisplayLabel.SetOwningConcept(crlDiagramElementPointer, hl)

	crlDiagramElementPointerAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramElementPointerAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramElementPointerAbstractionDisplayLabel, hl)
	crlDiagramElementPointerAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, hl)
	crlDiagramElementPointerAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramElementPointerAbstractionDisplayLabel, hl)

	crlDiagramElementPointerSource, _ := uOfD.NewReference(hl, CrlDiagramElementPointerSourceURI)
	crlDiagramElementPointerSource.SetLabel("Source", hl)
	crlDiagramElementPointerSource.SetOwningConcept(crlDiagramElementPointer, hl)

	crlDiagramElementPointerSourceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramElementPointerSourceRefinement.SetOwningConcept(crlDiagramElementPointerSource, hl)
	crlDiagramElementPointerSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, hl)
	crlDiagramElementPointerSourceRefinement.SetRefinedConcept(crlDiagramElementPointerSource, hl)

	crlDiagramElementPointerTarget, _ := uOfD.NewReference(hl, CrlDiagramElementPointerTargetURI)
	crlDiagramElementPointerTarget.SetLabel("Target", hl)
	crlDiagramElementPointerTarget.SetOwningConcept(crlDiagramElementPointer, hl)

	crlDiagramElementPointerTargetRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramElementPointerTargetRefinement.SetOwningConcept(crlDiagramElementPointerTarget, hl)
	crlDiagramElementPointerTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, hl)
	crlDiagramElementPointerTargetRefinement.SetRefinedConcept(crlDiagramElementPointerTarget, hl)

	//
	// OwnerPointer
	//
	crlDiagramOwnerPointer, _ := uOfD.NewElement(hl, CrlDiagramOwnerPointerURI)
	crlDiagramOwnerPointer.SetLabel("OwnerPointer", hl)
	crlDiagramOwnerPointer.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramOwnerPointerRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramOwnerPointerRefinement.SetOwningConcept(crlDiagramOwnerPointer, hl)
	crlDiagramOwnerPointerRefinement.SetAbstractConcept(crlDiagramPointer, hl)
	crlDiagramOwnerPointerRefinement.SetRefinedConcept(crlDiagramOwnerPointer, hl)

	crlDiagramOwnerPointerModelReference, _ := uOfD.NewReference(hl, CrlDiagramOwnerPointerModelReferenceURI)
	crlDiagramOwnerPointerModelReference.SetLabel("ModelReference", hl)
	crlDiagramOwnerPointerModelReference.SetOwningConcept(crlDiagramOwnerPointer, hl)

	crlDiagramOwnerPointerModelReferenceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramOwnerPointerModelReferenceRefinement.SetOwningConcept(crlDiagramOwnerPointerModelReference, hl)
	crlDiagramOwnerPointerModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, hl)
	crlDiagramOwnerPointerModelReferenceRefinement.SetRefinedConcept(crlDiagramOwnerPointerModelReference, hl)

	crlDiagramOwnerPointerDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramOwnerPointerDisplayLabelURI)
	crlDiagramOwnerPointerDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramOwnerPointerDisplayLabel.SetOwningConcept(crlDiagramOwnerPointer, hl)

	crlDiagramOwnerPointerDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramOwnerPointerDisplayLabelRefinement.SetOwningConcept(crlDiagramOwnerPointerDisplayLabel, hl)
	crlDiagramOwnerPointerDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, hl)
	crlDiagramOwnerPointerDisplayLabelRefinement.SetRefinedConcept(crlDiagramOwnerPointerDisplayLabel, hl)

	crlDiagramOwnerPointerAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramOwnerPointerAbstractionDisplayLabelURI)
	crlDiagramOwnerPointerAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramOwnerPointerAbstractionDisplayLabel.SetOwningConcept(crlDiagramOwnerPointer, hl)

	crlDiagramOwnerPointerAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramOwnerPointerAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramOwnerPointerAbstractionDisplayLabel, hl)
	crlDiagramOwnerPointerAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, hl)
	crlDiagramOwnerPointerAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramOwnerPointerAbstractionDisplayLabel, hl)

	crlDiagramOwnerPointerSource, _ := uOfD.NewReference(hl, CrlDiagramOwnerPointerSourceURI)
	crlDiagramOwnerPointerSource.SetLabel("Source", hl)
	crlDiagramOwnerPointerSource.SetOwningConcept(crlDiagramOwnerPointer, hl)

	crlDiagramOwnerPointerSourceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramOwnerPointerSourceRefinement.SetOwningConcept(crlDiagramOwnerPointerSource, hl)
	crlDiagramOwnerPointerSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, hl)
	crlDiagramOwnerPointerSourceRefinement.SetRefinedConcept(crlDiagramOwnerPointerSource, hl)

	crlDiagramOwnerPointerTarget, _ := uOfD.NewReference(hl, CrlDiagramOwnerPointerTargetURI)
	crlDiagramOwnerPointerTarget.SetLabel("Target", hl)
	crlDiagramOwnerPointerTarget.SetOwningConcept(crlDiagramOwnerPointer, hl)

	crlDiagramOwnerPointerTargetRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramOwnerPointerTargetRefinement.SetOwningConcept(crlDiagramOwnerPointerTarget, hl)
	crlDiagramOwnerPointerTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, hl)
	crlDiagramOwnerPointerTargetRefinement.SetRefinedConcept(crlDiagramOwnerPointerTarget, hl)

	//
	// RefinedPointer
	//
	crlDiagramRefinedPointer, _ := uOfD.NewElement(hl, CrlDiagramRefinedPointerURI)
	crlDiagramRefinedPointer.SetLabel("RefinedPointer", hl)
	crlDiagramRefinedPointer.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramRefinedPointerRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinedPointerRefinement.SetOwningConcept(crlDiagramRefinedPointer, hl)
	crlDiagramRefinedPointerRefinement.SetAbstractConcept(crlDiagramPointer, hl)
	crlDiagramRefinedPointerRefinement.SetRefinedConcept(crlDiagramRefinedPointer, hl)

	crlDiagramRefinedPointerModelReference, _ := uOfD.NewReference(hl, CrlDiagramRefinedPointerModelReferenceURI)
	crlDiagramRefinedPointerModelReference.SetLabel("ModelReference", hl)
	crlDiagramRefinedPointerModelReference.SetOwningConcept(crlDiagramRefinedPointer, hl)

	crlDiagramRefinedPointerModelReferenceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinedPointerModelReferenceRefinement.SetOwningConcept(crlDiagramRefinedPointerModelReference, hl)
	crlDiagramRefinedPointerModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, hl)
	crlDiagramRefinedPointerModelReferenceRefinement.SetRefinedConcept(crlDiagramRefinedPointerModelReference, hl)

	crlDiagramRefinedPointerDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramRefinedPointerDisplayLabelURI)
	crlDiagramRefinedPointerDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramRefinedPointerDisplayLabel.SetOwningConcept(crlDiagramRefinedPointer, hl)

	crlDiagramRefinedPointerDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinedPointerDisplayLabelRefinement.SetOwningConcept(crlDiagramRefinedPointerDisplayLabel, hl)
	crlDiagramRefinedPointerDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, hl)
	crlDiagramRefinedPointerDisplayLabelRefinement.SetRefinedConcept(crlDiagramRefinedPointerDisplayLabel, hl)

	crlDiagramRefinedPointerAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramRefinedPointerAbstractionDisplayLabelURI)
	crlDiagramRefinedPointerAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramRefinedPointerAbstractionDisplayLabel.SetOwningConcept(crlDiagramRefinedPointer, hl)

	crlDiagramRefinedPointerAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinedPointerAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramRefinedPointerAbstractionDisplayLabel, hl)
	crlDiagramRefinedPointerAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, hl)
	crlDiagramRefinedPointerAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramRefinedPointerAbstractionDisplayLabel, hl)

	crlDiagramRefinedPointerSource, _ := uOfD.NewReference(hl, CrlDiagramRefinedPointerSourceURI)
	crlDiagramRefinedPointerSource.SetLabel("Source", hl)
	crlDiagramRefinedPointerSource.SetOwningConcept(crlDiagramRefinedPointer, hl)

	crlDiagramRefinedPointerSourceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinedPointerSourceRefinement.SetOwningConcept(crlDiagramRefinedPointerSource, hl)
	crlDiagramRefinedPointerSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, hl)
	crlDiagramRefinedPointerSourceRefinement.SetRefinedConcept(crlDiagramRefinedPointerSource, hl)

	crlDiagramRefinedPointerTarget, _ := uOfD.NewReference(hl, CrlDiagramRefinedPointerTargetURI)
	crlDiagramRefinedPointerTarget.SetLabel("Target", hl)
	crlDiagramRefinedPointerTarget.SetOwningConcept(crlDiagramRefinedPointer, hl)

	crlDiagramRefinedPointerTargetRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinedPointerTargetRefinement.SetOwningConcept(crlDiagramRefinedPointerTarget, hl)
	crlDiagramRefinedPointerTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, hl)
	crlDiagramRefinedPointerTargetRefinement.SetRefinedConcept(crlDiagramRefinedPointerTarget, hl)

	//
	// ReferenceLink
	//
	crlDiagramReferenceLink, _ := uOfD.NewElement(hl, CrlDiagramReferenceLinkURI)
	crlDiagramReferenceLink.SetLabel("ReferenceLink", hl)
	crlDiagramReferenceLink.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramReferenceLinkRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramReferenceLinkRefinement.SetOwningConcept(crlDiagramReferenceLink, hl)
	crlDiagramReferenceLinkRefinement.SetAbstractConcept(crlDiagramLink, hl)
	crlDiagramReferenceLinkRefinement.SetRefinedConcept(crlDiagramReferenceLink, hl)

	crlDiagramReferenceLinkModelReference, _ := uOfD.NewReference(hl, CrlDiagramReferenceLinkModelReferenceURI)
	crlDiagramReferenceLinkModelReference.SetLabel("ModelReference", hl)
	crlDiagramReferenceLinkModelReference.SetOwningConcept(crlDiagramReferenceLink, hl)

	crlDiagramReferenceLinkModelReferenceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramReferenceLinkModelReferenceRefinement.SetOwningConcept(crlDiagramReferenceLinkModelReference, hl)
	crlDiagramReferenceLinkModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, hl)
	crlDiagramReferenceLinkModelReferenceRefinement.SetRefinedConcept(crlDiagramReferenceLinkModelReference, hl)

	crlDiagramReferenceLinkDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramReferenceLinkDisplayLabelURI)
	crlDiagramReferenceLinkDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramReferenceLinkDisplayLabel.SetOwningConcept(crlDiagramReferenceLink, hl)

	crlDiagramReferenceLinkDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramReferenceLinkDisplayLabelRefinement.SetOwningConcept(crlDiagramReferenceLinkDisplayLabel, hl)
	crlDiagramReferenceLinkDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, hl)
	crlDiagramReferenceLinkDisplayLabelRefinement.SetRefinedConcept(crlDiagramReferenceLinkDisplayLabel, hl)

	crlDiagramReferenceLinkAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramReferenceLinkAbstractionDisplayLabelURI)
	crlDiagramReferenceLinkAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramReferenceLinkAbstractionDisplayLabel.SetOwningConcept(crlDiagramReferenceLink, hl)

	crlDiagramReferenceLinkAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramReferenceLinkAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramReferenceLinkAbstractionDisplayLabel, hl)
	crlDiagramReferenceLinkAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, hl)
	crlDiagramReferenceLinkAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramReferenceLinkAbstractionDisplayLabel, hl)

	crlDiagramReferenceLinkSource, _ := uOfD.NewReference(hl, CrlDiagramReferenceLinkSourceURI)
	crlDiagramReferenceLinkSource.SetLabel("Source", hl)
	crlDiagramReferenceLinkSource.SetOwningConcept(crlDiagramReferenceLink, hl)

	crlDiagramReferenceLinkSourceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramReferenceLinkSourceRefinement.SetOwningConcept(crlDiagramReferenceLinkSource, hl)
	crlDiagramReferenceLinkSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, hl)
	crlDiagramReferenceLinkSourceRefinement.SetRefinedConcept(crlDiagramReferenceLinkSource, hl)

	crlDiagramReferenceLinkTarget, _ := uOfD.NewReference(hl, CrlDiagramReferenceLinkTargetURI)
	crlDiagramReferenceLinkTarget.SetLabel("Target", hl)
	crlDiagramReferenceLinkTarget.SetOwningConcept(crlDiagramReferenceLink, hl)

	crlDiagramReferenceLinkTargetRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramReferenceLinkTargetRefinement.SetOwningConcept(crlDiagramReferenceLinkTarget, hl)
	crlDiagramReferenceLinkTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, hl)
	crlDiagramReferenceLinkTargetRefinement.SetRefinedConcept(crlDiagramReferenceLinkTarget, hl)

	//
	// RefinementLink
	//
	crlDiagramRefinementLink, _ := uOfD.NewElement(hl, CrlDiagramRefinementLinkURI)
	crlDiagramRefinementLink.SetLabel("RefinementLink", hl)
	crlDiagramRefinementLink.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramRefinementLinkRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinementLinkRefinement.SetOwningConcept(crlDiagramRefinementLink, hl)
	crlDiagramRefinementLinkRefinement.SetAbstractConcept(crlDiagramLink, hl)
	crlDiagramRefinementLinkRefinement.SetRefinedConcept(crlDiagramRefinementLink, hl)

	crlDiagramRefinementLinkModelReference, _ := uOfD.NewReference(hl, CrlDiagramRefinementLinkModelReferenceURI)
	crlDiagramRefinementLinkModelReference.SetLabel("ModelReference", hl)
	crlDiagramRefinementLinkModelReference.SetOwningConcept(crlDiagramRefinementLink, hl)

	crlDiagramRefinementLinkModelReferenceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinementLinkModelReferenceRefinement.SetOwningConcept(crlDiagramRefinementLinkModelReference, hl)
	crlDiagramRefinementLinkModelReferenceRefinement.SetAbstractConcept(crlDiagramElementModelReference, hl)
	crlDiagramRefinementLinkModelReferenceRefinement.SetRefinedConcept(crlDiagramRefinementLinkModelReference, hl)

	crlDiagramRefinementLinkDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramRefinementLinkDisplayLabelURI)
	crlDiagramRefinementLinkDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramRefinementLinkDisplayLabel.SetOwningConcept(crlDiagramRefinementLink, hl)

	crlDiagramRefinementLinkDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinementLinkDisplayLabelRefinement.SetOwningConcept(crlDiagramRefinementLinkDisplayLabel, hl)
	crlDiagramRefinementLinkDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementDisplayLabel, hl)
	crlDiagramRefinementLinkDisplayLabelRefinement.SetRefinedConcept(crlDiagramRefinementLinkDisplayLabel, hl)

	crlDiagramRefinementLinkAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramRefinementLinkAbstractionDisplayLabelURI)
	crlDiagramRefinementLinkAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramRefinementLinkAbstractionDisplayLabel.SetOwningConcept(crlDiagramRefinementLink, hl)

	crlDiagramRefinementLinkAbstractionDisplayLabelRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinementLinkAbstractionDisplayLabelRefinement.SetOwningConcept(crlDiagramRefinementLinkAbstractionDisplayLabel, hl)
	crlDiagramRefinementLinkAbstractionDisplayLabelRefinement.SetAbstractConcept(crlDiagramElementAbstractionDisplayLabel, hl)
	crlDiagramRefinementLinkAbstractionDisplayLabelRefinement.SetRefinedConcept(crlDiagramRefinementLinkAbstractionDisplayLabel, hl)

	crlDiagramRefinementLinkSource, _ := uOfD.NewReference(hl, CrlDiagramRefinementLinkSourceURI)
	crlDiagramRefinementLinkSource.SetLabel("Source", hl)
	crlDiagramRefinementLinkSource.SetOwningConcept(crlDiagramRefinementLink, hl)

	crlDiagramRefinementLinkSourceRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinementLinkSourceRefinement.SetOwningConcept(crlDiagramRefinementLinkSource, hl)
	crlDiagramRefinementLinkSourceRefinement.SetAbstractConcept(crlDiagramLinkSource, hl)
	crlDiagramRefinementLinkSourceRefinement.SetRefinedConcept(crlDiagramRefinementLinkSource, hl)

	crlDiagramRefinementLinkTarget, _ := uOfD.NewReference(hl, CrlDiagramRefinementLinkTargetURI)
	crlDiagramRefinementLinkTarget.SetLabel("Target", hl)
	crlDiagramRefinementLinkTarget.SetOwningConcept(crlDiagramRefinementLink, hl)

	crlDiagramRefinementLinkTargetRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramRefinementLinkTargetRefinement.SetOwningConcept(crlDiagramRefinementLinkTarget, hl)
	crlDiagramRefinementLinkTargetRefinement.SetAbstractConcept(crlDiagramLinkTarget, hl)
	crlDiagramRefinementLinkTargetRefinement.SetRefinedConcept(crlDiagramRefinementLinkTarget, hl)

	uOfD.AddFunction(CrlDiagramElementURI, updateDiagramElement)
	uOfD.AddFunction(CrlDiagramOwnerPointerURI, updateDiagramOwnerPointer)

	crlDiagramConceptSpace.SetIsCoreRecursively(hl)
	return crlDiagramConceptSpace
}

// updateDiagramElement updates the diagram node based on changes to the modelElement it represents
func updateDiagramElement(diagramElement core.Element, notification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.WriteLockElement(diagramElement)
	// core Elements should always be ignored
	if diagramElement.GetIsCore(hl) == true {
		return nil
	}
	// There are several notifications of interest here:
	//   - the deletion of the referenced model element
	//   - the label of the referenced model element
	//   - the list of immediate abstractions of the referenced model element.
	// First, determine whether it is the referenced model element that has changed
	diagramElementModelReference := diagramElement.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)
	if diagramElementModelReference == nil {
		// Without a model reference, there is nothing to do. This scenario can occur during diagramElement deletion.
		return nil
	}
	modelElement := GetReferencedModelElement(diagramElement, hl)
	switch notification.GetNatureOfChange() {
	case core.IndicatedConceptChanged:
		if notification.GetAfterState().ConceptID == diagramElementModelReference.GetConceptID(hl) {
			modelReferenceNotification := notification.GetUnderlyingChange()
			switch modelReferenceNotification.GetNatureOfChange() {
			case core.IndicatedConceptChanged:
				modelElementNotification := modelReferenceNotification.GetUnderlyingChange()
				switch modelElementNotification.GetNatureOfChange() {
				case core.ConceptChanged:
					if IsDiagramNode(diagramElement, hl) {
						currentModelElement := modelElementNotification.GetAfterState()
						previousModelElement := modelElementNotification.GetBeforeState()
						if currentModelElement != nil && previousModelElement != nil {
							updateDiagramElementForModelElementChange(diagramElement, modelElement, hl)
						}
					} else if IsDiagramLink(diagramElement, hl) {
						diagram := diagramElement.GetOwningConcept(hl)
						oldLinkTarget := GetLinkTarget(diagramElement, hl)
						oldTargetModelElement := GetReferencedModelElement(oldLinkTarget, hl)
						switch modelElement.(type) {
						case core.Reference:
							reference := modelElement.(core.Reference)
							if IsDiagramElementPointer(diagramElement, hl) {
								newTargetModelElement := reference.GetReferencedConcept(hl)
								if oldTargetModelElement != newTargetModelElement {
									if newTargetModelElement == nil {
										uOfD.DeleteElement(diagramElement, hl)
									} else {
										newTargetDiagramElement := GetFirstElementRepresentingConcept(diagram, newTargetModelElement, hl)
										SetLinkTarget(diagramElement, newTargetDiagramElement, hl)
									}
								}
							} else if IsDiagramReferenceLink(diagramElement, hl) {
								updateDiagramElementForModelElementChange(diagramElement, reference, hl)
								SetDisplayLabel(diagramElement, reference.GetLabel(hl), hl)
								newModelTarget := reference.GetReferencedConcept(hl)
								newModelSource := reference.GetOwningConcept(hl)
								if newModelSource == nil || newModelTarget == nil {
									uOfD.DeleteElement(diagramElement, hl)
									return nil
								}
								currentDiagramSource := GetLinkSource(diagramElement, hl)
								currentModelSource := GetReferencedModelElement(currentDiagramSource, hl)
								currentDiagramTarget := GetLinkTarget(diagramElement, hl)
								currentModelTarget := GetReferencedModelElement(currentDiagramTarget, hl)
								if currentModelSource != newModelSource {
									newDiagramSource := GetFirstElementRepresentingConcept(diagram, newModelSource, hl)
									if newDiagramSource == nil {
										uOfD.DeleteElement(diagramElement, hl)
										return nil
									}
									SetLinkSource(diagramElement, newDiagramSource, hl)
								}
								if currentModelTarget != newModelTarget {
									newDiagramTarget := GetFirstElementRepresentingConcept(diagram, newModelTarget, hl)
									if newDiagramTarget == nil {
										uOfD.DeleteElement(diagramElement, hl)
										return nil
									}
									SetLinkTarget(diagramElement, newDiagramTarget, hl)
								}
							}
						case core.Refinement:
							refinement := modelElement.(core.Refinement)
							if IsDiagramPointer(diagramElement, hl) {
								var newTargetModelElement core.Element
								if IsDiagramAbstractPointer(diagramElement, hl) {
									newTargetModelElement = refinement.GetAbstractConcept(hl)
								} else if IsDiagramRefinedPointer(diagramElement, hl) {
									newTargetModelElement = refinement.GetRefinedConcept(hl)
								} else if IsDiagramOwnerPointer(diagramElement, hl) {
									newTargetModelElement = refinement.GetOwningConcept(hl)
								}
								if oldTargetModelElement != newTargetModelElement {
									if newTargetModelElement == nil {
										uOfD.DeleteElement(diagramElement, hl)
									} else {
										newTargetDiagramElement := GetFirstElementRepresentingConcept(diagram, newTargetModelElement, hl)
										SetLinkTarget(diagramElement, newTargetDiagramElement, hl)
									}
								}
							} else if IsDiagramRefinementLink(diagramElement, hl) {
								updateDiagramElementForModelElementChange(diagramElement, modelElement, hl)
								SetDisplayLabel(diagramElement, refinement.GetLabel(hl), hl)
								newModelTarget := refinement.GetAbstractConcept(hl)
								newModelSource := refinement.GetRefinedConcept(hl)
								if newModelTarget == nil || newModelSource == nil {
									uOfD.DeleteElement(diagramElement, hl)
									return nil
								}
								currentDiagramTarget := GetLinkTarget(diagramElement, hl)
								currentModelTarget := GetReferencedModelElement(currentDiagramTarget, hl)
								currentDiagramSource := GetLinkSource(diagramElement, hl)
								currentModelSource := GetReferencedModelElement(currentDiagramSource, hl)
								if currentModelTarget != newModelTarget {
									newDiagramTarget := GetFirstElementRepresentingConcept(diagram, newModelTarget, hl)
									if newDiagramTarget == nil {
										uOfD.DeleteElement(diagramElement, hl)
										return nil
									}
									SetLinkTarget(diagramElement, newDiagramTarget, hl)
								}
								if currentModelSource != newModelSource {
									newDiagramSource := GetFirstElementRepresentingConcept(diagram, newModelSource, hl)
									if newDiagramSource == nil {
										uOfD.DeleteElement(diagramElement, hl)
										return nil
									}
									SetLinkSource(diagramElement, newDiagramSource, hl)
								}
							}

						}
					}
				}
			case core.AbstractionChanged:
				// TODO: Implement AbstractionChanged case
			}
		}

	case core.ChildChanged:
		// We are looking for the model diagramElementModelReference reporting a ConceptChanged which would be the result of setting the referencedConcept
		if diagramElementModelReference == nil {
			break
		}
		afterConceptID := notification.GetAfterState().ConceptID
		modelReferenceID := diagramElementModelReference.GetConceptID(hl)
		if afterConceptID != modelReferenceID {
			break
		}
		modelReferenceNotification := notification.GetUnderlyingChange()
		switch modelReferenceNotification.GetNatureOfChange() {
		case core.ConceptChanged:
			if modelReferenceNotification.GetAfterState().ConceptID != diagramElementModelReference.GetConceptID(hl) {
				break
			}
			if diagramElementModelReference.(core.Reference).GetReferencedConceptID(hl) == "" {
				uOfD.DeleteElement(diagramElement, hl)
			} else {
				updateDiagramElementForModelElementChange(diagramElement, modelElement, hl)
			}
		}
	}
	return nil
}

// updateDiagramOwnerPointer updates the ownerPointer's target if the ownership of the represented modelElement changes
func updateDiagramOwnerPointer(diagramPointer core.Element, notification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
	// There is one change of interest here: the model element's owner has changed
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.WriteLockElement(diagramPointer)
	changedElement := uOfD.GetElement(notification.GetAfterState().ConceptID)
	modelElementReference := diagramPointer.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)
	diagram := diagramPointer.GetOwningConcept(hl)
	modelElement := GetReferencedModelElement(diagramPointer, hl)
	switch notification.GetNatureOfChange() {
	case core.IndicatedConceptChanged:
		if changedElement == modelElementReference {
			underlyingNotification := notification.GetUnderlyingChange()
			switch underlyingNotification.GetNatureOfChange() {
			case core.IndicatedConceptChanged:
				secondUnderlyingNotification := underlyingNotification.GetUnderlyingChange()
				switch secondUnderlyingNotification.GetNatureOfChange() {
				case core.ConceptChanged:
					if secondUnderlyingNotification.GetAfterState().ConceptID == modelElement.GetConceptID(hl) {
						modelOwner := modelElement.GetOwningConcept(hl)
						var oldModelOwner core.Element
						diagramTarget := GetLinkTarget(diagramPointer, hl)
						if diagramTarget != nil {
							oldModelOwner = GetReferencedModelElement(diagramTarget, hl)
						}
						if modelOwner != oldModelOwner {
							// Need to determine whether there is a view of the new owner in the diagram
							newDiagramTarget := GetFirstElementRepresentingConcept(diagram, modelOwner, hl)
							if newDiagramTarget == nil {
								// There is no view, delete the modelElement
								dEls := mapset.NewSet(diagramPointer.GetConceptID(hl))
								uOfD.DeleteElements(dEls, hl)
							} else {
								SetLinkTarget(diagramPointer, newDiagramTarget, hl)
							}
						}
					}
				}
			}
		}
	case core.ChildChanged:
		// If either source or target are nil, delete the pointer
		sourceReference := diagramPointer.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkSourceURI, hl)
		targetReference := diagramPointer.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramLinkTargetURI, hl)
		if changedElement == sourceReference || changedElement == targetReference {
			underlyingNotification := notification.GetUnderlyingChange()
			switch underlyingNotification.GetNatureOfChange() {
			case core.ConceptChanged:
				switch changedElement.(type) {
				case core.Reference:
					if changedElement.(core.Reference).GetReferencedConcept(hl) == nil {
						uOfD.DeleteElement(diagramPointer, hl)
					}
				}
			}
		}
	}
	return nil
}

func updateDiagramElementForModelElementChange(diagramElement core.Element, modelElement core.Element, hl *core.HeldLocks) {
	modelElementLabel := ""
	if modelElement != nil {
		modelElementLabel = modelElement.GetLabel(hl)
		if modelElementLabel != diagramElement.GetLabel(hl) {
			diagramElement.SetLabel(modelElementLabel, hl)
			if !IsDiagramPointer(diagramElement, hl) {
				SetDisplayLabel(diagramElement, modelElementLabel, hl)
			}
		}
		abstractions := make(map[string]core.Element)
		modelElement.FindImmediateAbstractions(abstractions, hl)
		abstractionsLabel := ""
		for _, abs := range abstractions {
			if len(abstractionsLabel) != 0 {
				abstractionsLabel += "\n"
			}
			abstractionsLabel += abs.GetLabel(hl)
		}
		if GetAbstractionDisplayLabel(diagramElement, hl) != abstractionsLabel {
			SetAbstractionDisplayLabel(diagramElement, abstractionsLabel, hl)
		}
	}
}

// updateNodeSize recalcualtes the size of the node based on the string sizes for the display label and
// abstractions listed
func updateNodeSize(node core.Element, hl *core.HeldLocks) {
	displayLabel := GetDisplayLabel(node, hl)
	displayLabelBounds, _ := font.BoundString(go12PtBoldFace, displayLabel)
	displayLabelMaxHeight := Int26_6ToFloat(displayLabelBounds.Max.Y)
	displayLabelMaxWidth := Int26_6ToFloat(displayLabelBounds.Max.X)
	displayLabelMinHeight := Int26_6ToFloat(displayLabelBounds.Min.Y)
	displayLabelMinWidth := Int26_6ToFloat(displayLabelBounds.Min.X)
	displayLabelHeight := displayLabelMaxHeight - displayLabelMinHeight
	displayLabelWidth := displayLabelMaxWidth - displayLabelMinWidth
	abstractionDisplayLabel := GetAbstractionDisplayLabel(node, hl)
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
	SetNodeHeight(node, height, hl)
	SetNodeWidth(node, width, hl)
	SetNodeDisplayLabelYOffset(node, displayLabelYOffset, hl)
}
