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
	"github.com/pbrown12303/activeCRL/core"
)

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

// CrlDiagramNodeURI identifies teh CrlDiagramNode conceot
var CrlDiagramNodeURI = CrlDiagramConceptSpaceURI + "/" + "CrlDiagramNode"

// CrlDiagramNodeModelReferenceURI identifies the reference to the model element represented by the node
var CrlDiagramNodeModelReferenceURI = CrlDiagramNodeURI + "/" + "ModelReference"

// CrlDiagramNodeDisplayLabelURI identifies the display label concept to be used when displaying the diagram
var CrlDiagramNodeDisplayLabelURI = CrlDiagramNodeURI + "/" + "DisplayLabel"

// CrlDiagramNodeXURI identifies the X coordinate of the node
var CrlDiagramNodeXURI = CrlDiagramNodeURI + "/" + "X"

// CrlDiagramNodeYURI identifies the Y coordinate of the node
var CrlDiagramNodeYURI = CrlDiagramNodeURI + "/" + "Y"

// CrlDiagramNodeHeightURI identifies the height of the node
var CrlDiagramNodeHeightURI = CrlDiagramNodeURI + "/" + "Height"

// CrlDiagramNodeWidthURI identifies the width of the node
var CrlDiagramNodeWidthURI = CrlDiagramNodeURI + "/" + "Width"

// CrlDiagramLinkURI identifies the concept of a link
var CrlDiagramLinkURI = CrlDiagramConceptSpaceURI + "/" + "CrlDiagramLink"

// GetReferencedElement is a function on a CrlDiagramNode that returns the model element represented by the
// diagram node
func GetReferencedElement(diagramNode core.Element, hl *core.HeldLocks) core.Element {
	if diagramNode == nil {
		return nil
	}
	reference := diagramNode.GetFirstChildReferenceWithAbstractionURI(CrlDiagramNodeModelReferenceURI, hl)
	if reference != nil {
		return reference.GetReferencedConcept(hl)
	}
	return nil
}

// SetReferencedElement is a function on a CrlDiagramNode that sets the model element represented by the
// diagram node
func SetReferencedElement(diagramNode core.Element, el core.Element, hl *core.HeldLocks) {
	if diagramNode == nil {
		return
	}
	reference := diagramNode.GetFirstChildReferenceWithAbstractionURI(CrlDiagramNodeModelReferenceURI, hl)
	if reference == nil {
		return
	}
	reference.SetReferencedConcept(el, hl)
}

// BuildCrlDiagramConceptSpace builds the CrlDiagram concept space and adds it to the uOfD
func BuildCrlDiagramConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// CrlDiagramConceptSpace
	crlDiagramConceptSpace, _ := uOfD.NewElement(hl, CrlDiagramConceptSpaceURI)
	crlDiagramConceptSpace.SetLabel("CrlDiagramConceptSpaceURI", hl)
	crlDiagramConceptSpace.SetURI(CrlDiagramConceptSpaceURI, hl)

	// CrlDiagram
	crlDiagram, _ := uOfD.NewElement(hl, CrlDiagramURI)
	crlDiagram.SetLabel("CrlDiagram", hl)
	crlDiagram.SetURI(CrlDiagramURI, hl)
	crlDiagram.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramWidth, _ := uOfD.NewLiteral(hl, CrlDiagramWidthURI)
	crlDiagramWidth.SetLabel("Width", hl)
	crlDiagramWidth.SetURI(CrlDiagramWidthURI, hl)
	crlDiagramWidth.SetOwningConcept(crlDiagram, hl)

	crlDiagramHeight, _ := uOfD.NewLiteral(hl, CrlDiagramHeightURI)
	crlDiagramHeight.SetLabel("Height", hl)
	crlDiagramHeight.SetURI(CrlDiagramHeightURI, hl)
	crlDiagramHeight.SetOwningConcept(crlDiagram, hl)

	// CrlDiagramNode
	crlDiagramNode, _ := uOfD.NewElement(hl, CrlDiagramNodeURI)
	crlDiagramNode.SetLabel("CrlDiagramNode", hl)
	crlDiagramNode.SetURI(CrlDiagramNodeURI, hl)
	crlDiagramNode.SetOwningConcept(crlDiagramConceptSpace, hl)

	crlDiagramNodeModelReference, _ := uOfD.NewReference(hl, CrlDiagramNodeModelReferenceURI)
	crlDiagramNodeModelReference.SetLabel("ModelReference", hl)
	crlDiagramNodeModelReference.SetURI(CrlDiagramNodeModelReferenceURI, hl)
	crlDiagramNodeModelReference.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeDisplayLabel, _ := uOfD.NewLiteral(hl, CrlDiagramNodeDisplayLabelURI)
	crlDiagramNodeDisplayLabel.SetLabel("DisplayLabel", hl)
	crlDiagramNodeDisplayLabel.SetURI(CrlDiagramNodeDisplayLabelURI, hl)
	crlDiagramNodeDisplayLabel.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeX, _ := uOfD.NewLiteral(hl, CrlDiagramNodeXURI)
	crlDiagramNodeX.SetLabel("X", hl)
	crlDiagramNodeX.SetURI(CrlDiagramNodeXURI, hl)
	crlDiagramNodeX.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeY, _ := uOfD.NewLiteral(hl, CrlDiagramNodeYURI)
	crlDiagramNodeY.SetLabel("Y", hl)
	crlDiagramNodeY.SetURI(CrlDiagramNodeYURI, hl)
	crlDiagramNodeY.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeHeight, _ := uOfD.NewLiteral(hl, CrlDiagramNodeHeightURI)
	crlDiagramNodeHeight.SetLabel("Height", hl)
	crlDiagramNodeHeight.SetURI(CrlDiagramNodeHeightURI, hl)
	crlDiagramNodeHeight.SetOwningConcept(crlDiagramNode, hl)

	crlDiagramNodeWidth, _ := uOfD.NewLiteral(hl, CrlDiagramNodeWidthURI)
	crlDiagramNodeWidth.SetLabel("Width", hl)
	crlDiagramNodeWidth.SetURI(CrlDiagramNodeWidthURI, hl)
	crlDiagramNodeWidth.SetOwningConcept(crlDiagramNode, hl)

	// CrlDiagramLink
	crlDiagramLink, _ := uOfD.NewElement(hl, CrlDiagramLinkURI)
	crlDiagramLink.SetLabel("CrlDiagramLink", hl)
	crlDiagramLink.SetURI(CrlDiagramLinkURI, hl)
	crlDiagramLink.SetOwningConcept(crlDiagramConceptSpace, hl)

	return crlDiagramConceptSpace
}
