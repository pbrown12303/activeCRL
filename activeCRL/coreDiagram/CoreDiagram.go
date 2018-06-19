// The coreDiagram package defines the CoreDiagram concept space. This is a pre-defined concept space (hence the term "core") that is, itself,
// represented as a CRLElement and identified with the CoreDiagramUri. This concept space contains the prototypes of all Elements used to construct CrlDiagrams.
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
// Instances of the prototpes can be conveniently instantiated using the core.CreateReplicateAsRefinementFromUri() function. This clones the prototype
// and, on an element-by-element basis, establishes a refinement relation between the instance elements and thier corresponding prototype elements.
// One essential side-effect of this is that changes that are made to the instnaces then trigger the execution of functions associated with the prototypes.
package coreDiagram

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
)

var CoreDiagramPrefix string = "http://activeCrl.com/coreDiagram/"
var CoreDiagramConceptSpace string = CoreDiagramPrefix + "CoreDiagram"

var CrlDiagramUri string = CoreDiagramConceptSpace + "/" + "CrlDiagram"
var CrlDiagramWidthUri string = CrlDiagramUri + "/" + "Width"
var CrlDiagramHeightUri string = CrlDiagramUri + "/" + "Height"

var CrlDiagramNodeUri string = CoreDiagramConceptSpace + "/" + "CrlDiagramNode"
var CrlDiagramNodeModelBaseElementReferenceUri string = CrlDiagramNodeUri + "/" + "ModelBaseElementReference"
var CrlDiagramNodeDisplayLabelUri string = CrlDiagramNodeUri + "/" + "DisplayLabel"
var CrlDiagramNodeXUri string = CrlDiagramNodeUri + "/" + "X"
var CrlDiagramNodeYUri string = CrlDiagramNodeUri + "/" + "Y"
var CrlDiagramNodeHeightUri string = CrlDiagramNodeUri + "/" + "Height"
var CrlDiagramNodeWidthUri string = CrlDiagramNodeUri + "/" + "Width"

var CrlDiagramLinkUri string = CoreDiagramConceptSpace + "/" + "CrlDiagramLink"

// AddCoreDiagramToUofD() constructs the core diagram concept space and adds it to the specified universe of discourse.
// It provides a means of bootstrapping, assembling the core diagram concepts space programmatically. It is expected that,
// after sufficient infrastructure has been developed, this concept space may be simply stored and retrieved rather than
// assembled programatically.
func AddCoreDiagramToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	coreDiagramConceptSpace := uOfD.GetElementWithUri(CoreDiagramConceptSpace)
	if coreDiagramConceptSpace == nil {
		coreDiagramConceptSpace = BuildCoreDiagramConceptSpace(uOfD, hl)
		if coreDiagramConceptSpace == nil {
			log.Printf("Build of CoreDiagram failed")
		}
	}
	return coreDiagramConceptSpace
}

func BuildCoreDiagramConceptSpace(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// Core
	coreDiagramConceptSpace := uOfD.NewElement(hl, CoreDiagramConceptSpace)
	core.SetLabel(coreDiagramConceptSpace, "CoreDiagramConceptSpace", hl)
	core.SetUri(coreDiagramConceptSpace, CoreDiagramConceptSpace, hl)

	// CrlDiagram
	crlDiagram := uOfD.NewElement(hl, CrlDiagramUri)
	core.SetLabel(crlDiagram, "CrlDiagram", hl)
	core.SetUri(crlDiagram, CrlDiagramUri, hl)
	core.SetOwningElement(crlDiagram, coreDiagramConceptSpace, hl)

	crlDiagramWidth := uOfD.NewLiteralReference(hl, CrlDiagramWidthUri)
	core.SetLabel(crlDiagramWidth, "Width", hl)
	core.SetUri(crlDiagramWidth, CrlDiagramWidthUri, hl)
	core.SetOwningElement(crlDiagramWidth, crlDiagram, hl)

	crlDiagramHeight := uOfD.NewLiteralReference(hl, CrlDiagramHeightUri)
	core.SetLabel(crlDiagramHeight, "Height", hl)
	core.SetUri(crlDiagramHeight, CrlDiagramHeightUri, hl)
	core.SetOwningElement(crlDiagramHeight, crlDiagram, hl)

	// CrlDiagramNode
	crlDiagramNode := uOfD.NewElement(hl, CrlDiagramNodeUri)
	core.SetLabel(crlDiagramNode, "CrlDiagramNode", hl)
	core.SetUri(crlDiagramNode, CrlDiagramNodeUri, hl)
	core.SetOwningElement(crlDiagramNode, coreDiagramConceptSpace, hl)

	crlDiagramNodeModelBaseElementReference := uOfD.NewBaseElementReference(hl, CrlDiagramNodeModelBaseElementReferenceUri)
	core.SetLabel(crlDiagramNodeModelBaseElementReference, "ModelBaseElementReference", hl)
	core.SetUri(crlDiagramNodeModelBaseElementReference, CrlDiagramNodeModelBaseElementReferenceUri, hl)
	core.SetOwningElement(crlDiagramNodeModelBaseElementReference, crlDiagramNode, hl)

	crlDiagramNodeDisplayLabel := uOfD.NewLiteralReference(hl, CrlDiagramNodeDisplayLabelUri)
	core.SetLabel(crlDiagramNodeDisplayLabel, "DisplayLabel", hl)
	core.SetUri(crlDiagramNodeDisplayLabel, CrlDiagramNodeDisplayLabelUri, hl)
	core.SetOwningElement(crlDiagramNodeDisplayLabel, crlDiagramNode, hl)

	crlDiagramNodeX := uOfD.NewLiteralReference(hl, CrlDiagramNodeXUri)
	core.SetLabel(crlDiagramNodeX, "X", hl)
	core.SetUri(crlDiagramNodeX, CrlDiagramNodeXUri, hl)
	core.SetOwningElement(crlDiagramNodeX, crlDiagramNode, hl)

	crlDiagramNodeY := uOfD.NewLiteralReference(hl, CrlDiagramNodeYUri)
	core.SetLabel(crlDiagramNodeY, "Y", hl)
	core.SetUri(crlDiagramNodeY, CrlDiagramNodeYUri, hl)
	core.SetOwningElement(crlDiagramNodeY, crlDiagramNode, hl)

	crlDiagramNodeHeight := uOfD.NewLiteralReference(hl, CrlDiagramNodeHeightUri)
	core.SetLabel(crlDiagramNodeHeight, "Height", hl)
	core.SetUri(crlDiagramNodeHeight, CrlDiagramNodeHeightUri, hl)
	core.SetOwningElement(crlDiagramNodeHeight, crlDiagramNode, hl)

	crlDiagramNodeWidth := uOfD.NewLiteralReference(hl, CrlDiagramNodeWidthUri)
	core.SetLabel(crlDiagramNodeWidth, "Width", hl)
	core.SetUri(crlDiagramNodeWidth, CrlDiagramNodeWidthUri, hl)
	core.SetOwningElement(crlDiagramNodeWidth, crlDiagramNode, hl)

	// CrlDiagramLink
	crlDiagramLink := uOfD.NewElement(hl, CrlDiagramLinkUri)
	core.SetLabel(crlDiagramLink, "CrlDiagramLink", hl)
	core.SetUri(crlDiagramLink, CrlDiagramLinkUri, hl)
	core.SetOwningElement(crlDiagramLink, coreDiagramConceptSpace, hl)

	//	BuildCoreBaseElementFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreBaseElementPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreBaseElementReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementPointerPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementPointerReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreElementReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralPointerPointerFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralPointerReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreLiteralReferenceFunctions(coreFunctionsElement, uOfD, hl)
	//	BuildCoreRefinementFunctions(coreFunctionsElement, uOfD, hl)

	return coreDiagramConceptSpace
}

func init() {
	//	baseElementFunctionsInit()
	//	baseElementPointerFunctionsInit()
	//	baseElementReferenceFunctionsInit()
	//	elementFunctionsInit()
	//	elementPointerFunctionsInit()
	//	elementPointerPointerFunctionsInit()
	//	elementPointerReferenceFunctionsInit()
	//	elementReferenceFunctionsInit()
	//	literalFunctionsInit()
	//	literalPointerFunctionsInit()
	//	literalPointerPointerFunctionsInit()
	//	literalPointerReferenceFunctionsInit()
	//	literalReferenceFunctionsInit()
	//	refinementFunctionsInit()
}
