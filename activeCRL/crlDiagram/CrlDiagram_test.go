package crlDiagram

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"sync"
	"testing"
)

func TestBuildCrlDiagramConceptSpace(t *testing.T) {
	var wg sync.WaitGroup
	hl := core.NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD := core.NewUniverseOfDiscourse(hl)
	uOfD.SetRecordingUndo(false)

	// CrlDiagramConceptSpace
	builtCrlDiagramConceptSpace := BuildCrlDiagramConceptSpace(uOfD, hl)
	if builtCrlDiagramConceptSpace == nil {
		t.Error("buildCoreDiagramConceptSpace returned empty element")
	}
	if core.GetUri(builtCrlDiagramConceptSpace, hl) != CrlDiagramConceptSpaceUri {
		t.Error("CrleDiagramConceptSpace uri not set")
	}
	_, ok := builtCrlDiagramConceptSpace.(core.Element)
	if !ok {
		t.Error("CrlDiagramConceptSpace is of wrong type")
	}
	if uOfD.GetElementWithUri(CrlDiagramConceptSpaceUri) == nil {
		t.Error("UofD uri index not updated with CrlDiagramConceptSpaceUri")
	}

	// CrlDiagram
	crlDiagram := core.GetChildElementWithUri(builtCrlDiagramConceptSpace, CrlDiagramUri, hl)
	if crlDiagram == nil {
		t.Error("CrlDiagram not found")
	}
	crlDiagramWidth := core.GetChildLiteralReferenceWithUri(crlDiagram, CrlDiagramWidthUri, hl)
	if crlDiagramWidth == nil {
		t.Error("CrlDiagramWidth not found")
	}
	crlDiagramHeight := core.GetChildLiteralReferenceWithUri(crlDiagram, CrlDiagramHeightUri, hl)
	if crlDiagramHeight == nil {
		t.Error("CrlDiagramHeight not found")
	}

	// CrlDiagramNode
	crlDiagramNode := core.GetChildElementWithUri(builtCrlDiagramConceptSpace, CrlDiagramNodeUri, hl)
	if crlDiagramNode == nil {
		t.Error("CrlDiagramNode not found")
	}
	crlDiagramNodeModelBaseElementReference := core.GetChildBaseElementReferenceWithUri(crlDiagramNode, CrlDiagramNodeModelBaseElementReferenceUri, hl)
	if crlDiagramNodeModelBaseElementReference == nil {
		t.Error("CrlDiagramNodeModelBaseElementReference not found")
	}
	crlDiagramNodeDisplayLabel := core.GetChildLiteralReferenceWithUri(crlDiagramNode, CrlDiagramNodeDisplayLabelUri, hl)
	if crlDiagramNodeDisplayLabel == nil {
		t.Error("CrlDiagramNodeDisplayLabel not found")
	}
	crlDiagramNodeX := core.GetChildLiteralReferenceWithUri(crlDiagramNode, CrlDiagramNodeXUri, hl)
	if crlDiagramNodeX == nil {
		t.Error("CrlDiagramNodeX not found")
	}
	crlDiagramNodeY := core.GetChildLiteralReferenceWithUri(crlDiagramNode, CrlDiagramNodeYUri, hl)
	if crlDiagramNodeY == nil {
		t.Error("CrlDiagramNodeY not found")
	}
	crlDiagramNodeHeight := core.GetChildLiteralReferenceWithUri(crlDiagramNode, CrlDiagramNodeHeightUri, hl)
	if crlDiagramNodeHeight == nil {
		t.Error("CrlDiagramNodeHeight not found")
	}
	crlDiagramNodeWidth := core.GetChildLiteralReferenceWithUri(crlDiagramNode, CrlDiagramNodeWidthUri, hl)
	if crlDiagramNodeWidth == nil {
		t.Error("crlDiagramNodeHeight not found")
	}

	// CrlDiagramLink
	crlDiagramLink := core.GetChildElementWithUri(builtCrlDiagramConceptSpace, CrlDiagramLinkUri, hl)
	if crlDiagramLink == nil {
		t.Error("CrlDiagramLink not found")
	}

}
