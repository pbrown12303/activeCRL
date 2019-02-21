// Package crldiagram defines the CoreDiagram concept space. This is a pre-defined concept space (hence the term "core") that is, itself,
// represented as a CRLElement and identified with the CoreDiagramURI. This concept space contains the prototypes of all Elements used to construct CrlDiagrams.
// Included are:
// 	CrlDiagram: the diagram itself
// 	CrlDiagramNode: a node in the diagram
// 	CrlDiagramLink: a link in the diagram
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
var CrlDiagramPrefix = "http://activeCrl.com/coreDiagram/"

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
		if GetReferencedModelElement(el, hl) == concept {
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

// IsDiagram returns true if the supplied element is a crldiagram
func IsDiagram(el core.Element, hl *core.HeldLocks) bool {
	switch el.(type) {
	case core.Element:
		return el.IsRefinementOfURI(CrlDiagramURI, hl)
	}
	return false
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
	crlDiagramConceptSpace.SetURI(CrlDiagramConceptSpaceURI, hl)
	crlDiagramConceptSpace.SetIsCore(hl)

	// CrlDiagram
	crlDiagram, _ := uOfD.NewElement(hl, CrlDiagramURI)
	crlDiagram.SetLabel("CrlDiagram", hl)
	crlDiagram.SetURI(CrlDiagramURI, hl)
	crlDiagram.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagram.SetIsCore(hl)

	crlDiagramWidth, _ := uOfD.NewLiteral(hl, CrlDiagramWidthURI)
	crlDiagramWidth.SetLabel("Width", hl)
	crlDiagramWidth.SetURI(CrlDiagramWidthURI, hl)
	crlDiagramWidth.SetOwningConcept(crlDiagram, hl)
	crlDiagramWidth.SetIsCore(hl)

	crlDiagramHeight, _ := uOfD.NewLiteral(hl, CrlDiagramHeightURI)
	crlDiagramHeight.SetLabel("Height", hl)
	crlDiagramHeight.SetURI(CrlDiagramHeightURI, hl)
	crlDiagramHeight.SetOwningConcept(crlDiagram, hl)
	crlDiagramHeight.SetIsCore(hl)

	// CrlDiagramElement
	crlDiagramElement, _ := uOfD.NewElement(hl, CrlDiagramElementURI)
	crlDiagramElement.SetLabel("CrlDiagramElement", hl)
	crlDiagramElement.SetURI(CrlDiagramElementURI, hl)
	crlDiagramElement.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramElement.SetIsCore(hl)

	// CrlDiagramNode
	crlDiagramNode, _ := uOfD.NewElement(hl, CrlDiagramNodeURI)
	crlDiagramNode.SetLabel("CrlDiagramNode", hl)
	crlDiagramNode.SetURI(CrlDiagramNodeURI, hl)
	crlDiagramNode.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramNode.SetIsCore(hl)

	crlDiagramNodeElementRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramNodeElementRefinement.SetLabel("CrlDiagramNode refines CrlDiagramElement", hl)
	crlDiagramNodeElementRefinement.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeElementRefinement.SetAbstractConcept(crlDiagramElement, hl)
	crlDiagramNodeElementRefinement.SetRefinedConcept(crlDiagramNode, hl)

	crlDiagramNodeModelReference, _ := uOfD.NewReference(hl, CrlDiagramElementModelReferenceURI)
	crlDiagramNodeModelReference.SetLabel("ModelReference", hl)
	crlDiagramNodeModelReference.SetURI(CrlDiagramElementModelReferenceURI, hl)
	crlDiagramNodeModelReference.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeModelReference.SetIsCore(hl)

	crlDiagramNodeDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementDisplayLabelURI)
	crlDiagramNodeDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramNodeDisplayLabel.SetURI(CrlDiagramElementDisplayLabelURI, hl)
	crlDiagramNodeDisplayLabel.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeDisplayLabel.SetIsCore(hl)

	crlDiagramNodeAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementAbstractionDisplayLabelURI)
	crlDiagramNodeAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramNodeAbstractionDisplayLabel.SetURI(CrlDiagramElementAbstractionDisplayLabelURI, hl)
	crlDiagramNodeAbstractionDisplayLabel.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeAbstractionDisplayLabel.SetIsCore(hl)

	crlDiagramNodeX, _ := uOfD.NewLiteral(hl, CrlDiagramNodeXURI)
	crlDiagramNodeX.SetLabel("X", hl)
	crlDiagramNodeX.SetURI(CrlDiagramNodeXURI, hl)
	crlDiagramNodeX.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeX.SetIsCore(hl)

	crlDiagramNodeY, _ := uOfD.NewLiteral(hl, CrlDiagramNodeYURI)
	crlDiagramNodeY.SetLabel("Y", hl)
	crlDiagramNodeY.SetURI(CrlDiagramNodeYURI, hl)
	crlDiagramNodeY.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeY.SetIsCore(hl)

	crlDiagramNodeHeight, _ := uOfD.NewLiteral(hl, CrlDiagramNodeHeightURI)
	crlDiagramNodeHeight.SetLabel("Height", hl)
	crlDiagramNodeHeight.SetURI(CrlDiagramNodeHeightURI, hl)
	crlDiagramNodeHeight.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeHeight.SetIsCore(hl)

	crlDiagramNodeWidth, _ := uOfD.NewLiteral(hl, CrlDiagramNodeWidthURI)
	crlDiagramNodeWidth.SetLabel("Width", hl)
	crlDiagramNodeWidth.SetURI(CrlDiagramNodeWidthURI, hl)
	crlDiagramNodeWidth.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeWidth.SetIsCore(hl)

	crlDiagramNodeDisplayLabelYOffset, _ := uOfD.NewLiteral(hl, CrlDiagramNodeDisplayLabelYOffsetURI)
	crlDiagramNodeDisplayLabelYOffset.SetLabel("DisplayLabelYOffset", hl)
	crlDiagramNodeDisplayLabelYOffset.SetURI(CrlDiagramNodeDisplayLabelYOffsetURI, hl)
	crlDiagramNodeDisplayLabelYOffset.SetOwningConcept(crlDiagramNode, hl)
	crlDiagramNodeDisplayLabelYOffset.SetIsCore(hl)

	// CrlDiagramLink
	crlDiagramLink, _ := uOfD.NewElement(hl, CrlDiagramLinkURI)
	crlDiagramLink.SetLabel("CrlDiagramLink", hl)
	crlDiagramLink.SetURI(CrlDiagramLinkURI, hl)
	crlDiagramLink.SetOwningConcept(crlDiagramConceptSpace, hl)
	crlDiagramLink.SetIsCore(hl)

	crlDiagramLinkElementRefinement, _ := uOfD.NewRefinement(hl)
	crlDiagramLinkElementRefinement.SetLabel("CrlDiagramLink refines CrlDiagramElement", hl)
	crlDiagramLinkElementRefinement.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkElementRefinement.SetAbstractConcept(crlDiagramElement, hl)
	crlDiagramLinkElementRefinement.SetRefinedConcept(crlDiagramLink, hl)

	crlDiagramLinkModelReference, _ := uOfD.NewReference(hl, CrlDiagramElementModelReferenceURI)
	crlDiagramLinkModelReference.SetLabel("ModelReference", hl)
	crlDiagramLinkModelReference.SetURI(CrlDiagramElementModelReferenceURI, hl)
	crlDiagramLinkModelReference.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkModelReference.SetIsCore(hl)

	crlDiagramLinkSource, _ := uOfD.NewReference(hl, CrlDiagramLinkSourceURI)
	crlDiagramLinkSource.SetLabel("Source", hl)
	crlDiagramLinkSource.SetURI(CrlDiagramLinkSourceURI, hl)
	crlDiagramLinkSource.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkSource.SetIsCore(hl)

	crlDiagramLinkTarget, _ := uOfD.NewReference(hl, CrlDiagramLinkTargetURI)
	crlDiagramLinkTarget.SetLabel("Target", hl)
	crlDiagramLinkTarget.SetURI(CrlDiagramLinkTargetURI, hl)
	crlDiagramLinkTarget.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkTarget.SetIsCore(hl)

	crlDiagramLinkDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementDisplayLabelURI)
	crlDiagramLinkDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramLinkDisplayLabel.SetURI(CrlDiagramElementDisplayLabelURI, hl)
	crlDiagramLinkDisplayLabel.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkDisplayLabel.SetIsCore(hl)

	crlDiagramLinkAbstractionDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramElementAbstractionDisplayLabelURI)
	crlDiagramLinkAbstractionDisplayLabel.SetLabel("AbstractionDisplayLabel", hl)
	crlDiagramLinkAbstractionDisplayLabel.SetURI(CrlDiagramElementAbstractionDisplayLabelURI, hl)
	crlDiagramLinkAbstractionDisplayLabel.SetOwningConcept(crlDiagramLink, hl)
	crlDiagramLinkAbstractionDisplayLabel.SetIsCore(hl)

	uOfD.AddFunction(CrlDiagramNodeURI, updateDiagramNode)

	return crlDiagramConceptSpace
}

func updateDiagramNode(node core.Element, notification *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.WriteLockElement(node)
	// There are two notifications of interest here: the label of the referenced model element
	// and the list of immediate abstractions of the referenced model element.
	// First, determine whether it is the referenced model element that has changed
	reference := node.GetFirstOwnedReferenceRefinedFromURI(CrlDiagramElementModelReferenceURI, hl)
	modelElement := GetReferencedModelElement(node, hl)
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
						updateNodeForModelElementChange(node, modelElement, hl)
					}
				case core.AbstractionChanged:

				}
			}
		}

	case core.ChildChanged:
		// We are looking for the model reference reporting a ConceptChanged which would be the result of setting the referencedConcept
		if notification.GetReportingElement() != reference {
			break
		}
		underlyingNotification := notification.GetUnderlyingChange()
		switch underlyingNotification.GetNatureOfChange() {
		case core.ConceptChanged:
			if underlyingNotification.GetReportingElement() != reference {
				break
			}
			updateNodeForModelElementChange(node, modelElement, hl)
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
