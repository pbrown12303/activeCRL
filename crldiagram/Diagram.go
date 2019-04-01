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
	"log"
	"math"
	"strconv"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
	"github.com/pbrown12303/activeCRL/core"

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

// CrlDiagramElementPointerURI identifies the element pointer of a Reference represented as a link
var CrlDiagramElementPointerURI = CrlDiagramConceptSpaceURI + "/" + "ElementPointer"

// CrlDiagramOwnerPointerURI identifies the owner of an Element represented as a link
var CrlDiagramOwnerPointerURI = CrlDiagramConceptSpaceURI + "/" + "OwnerPointer"

// CrlDiagramRefinedPointerURI identifies the refined element of a Refinement represented as a link
var CrlDiagramRefinedPointerURI = CrlDiagramConceptSpaceURI + "/" + "RefinedPointer"

// CrlDiagramReferenceLinkURI identifies the Reference represented as a link in the diagram
var CrlDiagramReferenceLinkURI = CrlDiagramConceptSpaceURI + "/" + "ReferenceLink"

// CrlDiagramRefinementLinkURI identifies the Refinement represented as a link in the diagram
var CrlDiagramRefinementLinkURI = CrlDiagramConceptSpaceURI + "/" + "RefinementLink"

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
		log.Printf("GetFirstElementRepresentingConcept called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramElementURI, hl) {
		if GetReferencedModelElement(el, hl).GetConceptID(hl) == conceptID && !el.IsRefinementOfURI(CrlDiagramPointerURI, hl) {
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

// GetOwnerPointer returns the ownerPoiner for the concept if one exists
func GetOwnerPointer(diagram core.Element, concept core.Element, hl *core.HeldLocks) core.Element {
	if diagram.IsRefinementOfURI(CrlDiagramURI, hl) == false {
		log.Printf("GetFirstElementRepresentingConcept called with diagram of incorrect type")
		return nil
	}
	for _, el := range diagram.GetOwnedConceptsRefinedFromURI(CrlDiagramOwnerPointerURI, hl) {
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

// IsDiagramElement returns true if the supplied element is a CrlDiagramElement
func IsDiagramElement(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramElementURI, hl)
}

// IsDiagramNode returns true if the supplied element is a CrlDiagramElement
func IsDiagramNode(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramNodeURI, hl)
}

// IsDiagramLink returns true if the supplied element is a CrlDiagramElement
func IsDiagramLink(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramLinkURI, hl)
}

// IsDiagramPointer returns true if the supplied element is a CrlDiagramElement
func IsDiagramPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramPointerURI, hl)
}

// IsDiagramOwnerPointer returns true if the supplied element is a CrlDiagramElement
func IsDiagramOwnerPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramOwnerPointerURI, hl)
}

// IsDiagramElementPointer returns true if the supplied element is a CrlDiagramElement
func IsDiagramElementPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramElementPointerURI, hl)
}

// IsDiagramAbstractPointer returns true if the supplied element is a CrlDiagramElement
func IsDiagramAbstractPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramAbstractPointerURI, hl)
}

// IsDiagramRefinedPointer returns true if the supplied element is a CrlDiagramElement
func IsDiagramRefinedPointer(el core.Element, hl *core.HeldLocks) bool {
	return el.IsRefinementOfURI(CrlDiagramRefinedPointerURI, hl)
}

// NewDiagram creates a new diagram
func NewDiagram(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramURI, hl)
}

// NewDiagramReferenceLink creates a new diagram link to represent a reference
func NewDiagramReferenceLink(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramReferenceLinkURI, hl)
}

// NewDiagramRefinementLink creates a new diagram link
func NewDiagramRefinementLink(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramRefinementLinkURI, hl)
}

// NewDiagramNode creates a new diagram node
func NewDiagramNode(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramNodeURI, hl)
}

// NewDiagramOwnerPointer creates a new DiagramOwnerPointer
func NewDiagramOwnerPointer(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramOwnerPointerURI, hl)
}

// NewDiagramElementPointer creates a new DiagramElementPointer
func NewDiagramElementPointer(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramElementPointerURI, hl)
}

// NewDiagramAbstractPointer creates a new DiagramAbstractPointer
func NewDiagramAbstractPointer(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	return uOfD.CreateReplicateAsRefinementFromURI(CrlDiagramAbstractPointerURI, hl)
}

// NewDiagramRefinedPointer creates a new DiagramRefinedPointer
func NewDiagramRefinedPointer(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
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

// SetDisplayLabel is a function on a CrlDiagramNode that sets the display label of the diagram node
func SetDisplayLabel(diagramElement core.Element, value string, hl *core.HeldLocks) {
	if diagramElement == nil {
		return
	}
	literal := diagramElement.GetFirstOwnedLiteralRefinementOfURI(CrlDiagramElementDisplayLabelURI, hl)
	if literal == nil {
		return
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
func BuildCrlDiagramConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// CrlDiagramConceptSpace
	crlDiagramConceptSpace, _ := uOfD.NewElement(hl, CrlDiagramConceptSpaceURI)
	crlDiagramConceptSpace.SetLabel("CrlDiagramConceptSpace", hl)
	crlDiagramConceptSpace.SetIsCore(hl)

	//
	// CrlDiagram
	//
	crlDiagram, _ := uOfD.NewElement(hl, CrlDiagramURI)
	crlDiagram.SetLabel("CrlDiagram", hl)
	crlDiagram.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagram.SetIsCore(hl)

	crlDiagramWidth, _ := uOfD.NewLiteral(hl, CrlDiagramWidthURI)
	crlDiagramWidth.SetLabel("Width", hl)
	crlDiagramWidth.SetOwningConcept(crlDiagram, hl)
	crlDiagramWidth.SetIsCore(hl)

	crlDiagramHeight, _ := uOfD.NewLiteral(hl, CrlDiagramHeightURI)
	crlDiagramHeight.SetLabel("Height", hl)
	crlDiagramHeight.SetOwningConcept(crlDiagram, hl)
	crlDiagramHeight.SetIsCore(hl)

	//
	// CrlDiagramElement
	//
	crlDiagramElement, _ := uOfD.NewElement(hl, CrlDiagramElementURI)
	crlDiagramElement.SetLabel("CrlDiagramElement", hl)
	crlDiagramElement.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramElement.SetIsCore(hl)
	hl.ReleaseLocksAndWait()

	crlDiagramElementModelReference, _ := uOfD.NewReference(hl, CrlDiagramElementModelReferenceURI)
	crlDiagramElementModelReference.SetLabel("ModelReference", hl)
	crlDiagramElementModelReference.SetOwningConcept(crlDiagramElement, hl)
	crlDiagramElementModelReference.SetIsCore(hl)

	crlDiagramElementDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementDisplayLabelURI)
	crlDiagramElementDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramElementDisplayLabel.SetOwningConcept(crlDiagramElement, hl)
	crlDiagramElementDisplayLabel.SetIsCore(hl)

	crlDiagramElementAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementAbstractionDisplayLabelURI)
	crlDiagramElementAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramElementAbstractionDisplayLabel.SetOwningConcept(crlDiagramElement, hl)
	crlDiagramElementAbstractionDisplayLabel.SetIsCore(hl)

	//
	// CrlDiagramNode
	//
	crlDiagramNode := uOfD.CreateReplicateAsRefinement(crlDiagramElement, hl, CrlDiagramNodeURI)
	crlDiagramNode.SetLabel("CrlDiagramNode", hl)
	crlDiagramNode.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramNode.SetIsCoreRecursively(hl)

	crlDiagramNodeX, _ := uOfD.NewLiteral(hl, CrlDiagramNodeXURI)
	crlDiagramNodeX.SetLabel("X", hl)
	crlDiagramNodeX.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeX.SetIsCore(hl)

	crlDiagramNodeY, _ := uOfD.NewLiteral(hl, CrlDiagramNodeYURI)
	crlDiagramNodeY.SetLabel("Y", hl)
	crlDiagramNodeY.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeY.SetIsCore(hl)

	crlDiagramNodeHeight, _ := uOfD.NewLiteral(hl, CrlDiagramNodeHeightURI)
	crlDiagramNodeHeight.SetLabel("Height", hl)
	crlDiagramNodeHeight.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeHeight.SetIsCore(hl)

	crlDiagramNodeWidth, _ := uOfD.NewLiteral(hl, CrlDiagramNodeWidthURI)
	crlDiagramNodeWidth.SetLabel("Width", hl)
	crlDiagramNodeWidth.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeWidth.SetIsCore(hl)

	crlDiagramNodeDisplayLabelYOffset, _ := uOfD.NewLiteral(hl, CrlDiagramNodeDisplayLabelYOffsetURI)
	crlDiagramNodeDisplayLabelYOffset.SetLabel("DisplayLabelYOffset", hl)
	crlDiagramNodeDisplayLabelYOffset.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeDisplayLabelYOffset.SetIsCore(hl)

	//
	// CrlDiagramLink
	//
	crlDiagramLink := uOfD.CreateReplicateAsRefinement(crlDiagramElement, hl, CrlDiagramLinkURI)
	crlDiagramLink.SetLabel("CrlDiagramLink", hl)
	crlDiagramLink.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramLink.SetIsCoreRecursively(hl)

	crlDiagramLinkSource, _ := uOfD.NewReference(hl, CrlDiagramLinkSourceURI)
	crlDiagramLinkSource.SetLabel("Source", hl)
	crlDiagramLinkSource.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkSource.SetIsCore(hl)

	crlDiagramLinkTarget, _ := uOfD.NewReference(hl, CrlDiagramLinkTargetURI)
	crlDiagramLinkTarget.SetLabel("Target", hl)
	crlDiagramLinkTarget.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkTarget.SetIsCore(hl)

	//
	// Pointer
	//
	crlDiagramPointer := uOfD.CreateReplicateAsRefinement(crlDiagramLink, hl, CrlDiagramPointerURI)
	crlDiagramPointer.SetLabel("Pointer", hl)
	crlDiagramPointer.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramPointer.SetIsCoreRecursively(hl)

	//
	// AbstractPointer
	//
	crlDiagramAbstractPointer := uOfD.CreateReplicateAsRefinement(crlDiagramPointer, hl, CrlDiagramAbstractPointerURI)
	crlDiagramAbstractPointer.SetLabel("AbstractPointer", hl)
	crlDiagramAbstractPointer.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramAbstractPointer.SetIsCore(hl)

	//
	// ElementPointer
	//
	crlDiagramElementPointer := uOfD.CreateReplicateAsRefinement(crlDiagramPointer, hl, CrlDiagramElementPointerURI)
	crlDiagramElementPointer.SetLabel("ElementPointer", hl)
	crlDiagramElementPointer.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramElementPointer.SetIsCore(hl)

	//
	// OwnerPointer
	//
	crlDiagramOwnerPointer := uOfD.CreateReplicateAsRefinement(crlDiagramPointer, hl, CrlDiagramOwnerPointerURI)
	crlDiagramOwnerPointer.SetLabel("OwnerPointer", hl)
	crlDiagramOwnerPointer.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramOwnerPointer.SetIsCore(hl)

	//
	// RefinedPointer
	//
	crlDiagramRefinedPointer := uOfD.CreateReplicateAsRefinement(crlDiagramPointer, hl, CrlDiagramRefinedPointerURI)
	crlDiagramRefinedPointer.SetLabel("RefinedPointer", hl)
	crlDiagramRefinedPointer.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramRefinedPointer.SetIsCore(hl)

	//
	// ReferenceLink
	//
	crlDiagramReferenceLink := uOfD.CreateReplicateAsRefinement(crlDiagramLink, hl, CrlDiagramReferenceLinkURI)
	crlDiagramReferenceLink.SetLabel("ReferenceLink", hl)
	crlDiagramReferenceLink.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramReferenceLink.SetIsCore(hl)

	//
	// RefinementLink
	//
	crlDiagramRefinementLink := uOfD.CreateReplicateAsRefinement(crlDiagramLink, hl, CrlDiagramRefinementLinkURI)
	crlDiagramRefinementLink.SetLabel("RefinementLink", hl)
	crlDiagramRefinementLink.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramRefinementLink.SetIsCore(hl)

	uOfD.AddFunction(CrlDiagramNodeURI, updateDiagramNode)
	uOfD.AddFunction(CrlDiagramOwnerPointerURI, updateDiagramOwnerPointer)
	// uOfD.AddFunction(CrlDiagramURI, updateDiagram)

	return crlDiagramConceptSpace
}

// updateDiagramNode updates the diagram node based on changes to the modelElement it represents
func updateDiagramNode(node core.Element, notification *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.WriteLockElement(node)
	// There are several notifications of interest here:
	//   - the deletion of the referenced model element
	//   - the label of the referenced model element
	//   - the list of immediate abstractions of the referenced model element.
	// First, determine whether it is the referenced model element that has changed
	diagramElementModelReference := node.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)
	modelElement := GetReferencedModelElement(node, hl)
	switch notification.GetNatureOfChange() {
	case core.IndicatedConceptChanged:
		if notification.GetReportingElement() == diagramElementModelReference {
			modelReferenceNotification := notification.GetUnderlyingChange()
			switch modelReferenceNotification.GetNatureOfChange() {
			case core.IndicatedConceptChanged:
				modelElementNotification := modelReferenceNotification.GetUnderlyingChange()
				switch modelElementNotification.GetNatureOfChange() {
				case core.ConceptChanged:
					currentModelElement := modelElementNotification.GetReportingElement()
					previousModelElement := modelElementNotification.GetPriorState()
					if currentModelElement != nil && previousModelElement != nil {
						updateNodeForModelElementChange(node, modelElement, hl)
					}
				case core.AbstractionChanged:

				}
			}
		}

	case core.ChildChanged:
		// We are looking for the model diagramElementModelReference reporting a ConceptChanged which would be the result of setting the referencedConcept
		if notification.GetReportingElement() != diagramElementModelReference {
			break
		}
		modelReferenceNotification := notification.GetUnderlyingChange()
		switch modelReferenceNotification.GetNatureOfChange() {
		case core.ConceptChanged:
			if modelReferenceNotification.GetReportingElement() != diagramElementModelReference {
				break
			}
			if diagramElementModelReference.(core.Reference).GetReferencedConceptID(hl) == "" {
				uOfD.DeleteElement(node, hl)
			} else {
				updateNodeForModelElementChange(node, modelElement, hl)
			}
		}
	}
}

// updateDiagramOwnerPointer updates the ownerPointer's target if the ownership of the represented modelElement changes
func updateDiagramOwnerPointer(diagramPointer core.Element, notification *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	// There is one change of interest here: the model element's owner has changed
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.WriteLockElement(diagramPointer)
	reference := diagramPointer.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)
	diagram := diagramPointer.GetOwningConcept(hl)
	modelElement := GetReferencedModelElement(diagramPointer, hl)
	switch notification.GetNatureOfChange() {
	case core.IndicatedConceptChanged:
		if notification.GetReportingElement() == reference {
			underlyingNotification := notification.GetUnderlyingChange()
			switch underlyingNotification.GetNatureOfChange() {
			case core.IndicatedConceptChanged:
				secondUnderlyingNotification := underlyingNotification.GetUnderlyingChange()
				switch secondUnderlyingNotification.GetNatureOfChange() {
				case core.ConceptChanged:
					if secondUnderlyingNotification.GetReportingElement() == modelElement {
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
								dEls := map[string]core.Element{diagramPointer.GetConceptID(hl): diagramPointer}
								uOfD.DeleteElements(dEls, hl)
							} else {
								SetLinkTarget(diagramPointer, newDiagramTarget, hl)
							}
						}
					}
				}
			}
		}
	}
}

func updateNodeForModelElementChange(node core.Element, modelElement core.Element, hl *core.HeldLocks) {
	modelElementLabel := ""
	if modelElement != nil {
		modelElementLabel = modelElement.GetLabel(hl)
		if modelElementLabel != node.GetLabel(hl) {
			node.SetLabel(modelElementLabel, hl)
			SetDisplayLabel(node, modelElementLabel, hl)
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
		if GetAbstractionDisplayLabel(node, hl) != abstractionsLabel {
			SetAbstractionDisplayLabel(node, abstractionsLabel, hl)
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
