package crleditor

import (
	"encoding/json"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pkg/errors"
)

type DiagramManager struct {
	editor   *Editor
	diagrams map[string]core.Element
}

func NewDiagramManager(editor *Editor) *DiagramManager {
	var dMgr DiagramManager
	dMgr.editor = editor
	dMgr.diagrams = map[string]core.Element{}
	return &dMgr
}

func (dMgr *DiagramManager) AddDiagram(ownerID string, trans *core.Transaction) (core.Element, error) {
	diagram, err := dMgr.NewDiagram(trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = diagram.SetOwningConceptID(ownerID, trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = dMgr.editor.SelectElement(diagram, trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = dMgr.DisplayDiagram(diagram.GetConceptID(trans), trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	return diagram, nil
}

func (dMgr *DiagramManager) AddConceptView(diagramID string, conceptID string, x float64, y float64, trans *core.Transaction) (core.Element, error) {
	uOfD := dMgr.editor.GetUofD()
	diagram := uOfD.GetElement(diagramID)
	el := uOfD.GetElement(conceptID)
	if el == nil {
		return nil, errors.New("Indicated model element not found in addNodeView, ID: " + conceptID)
	}
	createAsLink := false
	switch el.(type) {
	case core.Reference:
		createAsLink = dMgr.editor.GetDropDiagramReferenceAsLink(trans)
	case core.Refinement:
		createAsLink = dMgr.editor.GetDropDiagramRefinementAsLink(trans)
	}

	var newElement core.Element
	var err error
	if createAsLink {
		var modelSourceConcept core.Element
		var modelTargetConcept core.Element
		switch elTyped := el.(type) {
		case core.Reference:
			newElement, err = crldiagramdomain.NewDiagramReferenceLink(uOfD, trans)
			if err != nil {
				return nil, err
			}
			reference := elTyped
			modelSourceConcept = reference.GetOwningConcept(trans)
			modelTargetConcept = reference.GetReferencedConcept(trans)
		case core.Refinement:
			newElement, err = crldiagramdomain.NewDiagramRefinementLink(uOfD, trans)
			if err != nil {
				return nil, err
			}
			refinement := elTyped
			modelSourceConcept = refinement.GetRefinedConcept(trans)
			modelTargetConcept = refinement.GetAbstractConcept(trans)
		}
		if modelSourceConcept == nil {
			return nil, errors.New("In addConceptView for link, modelSourceConcept is nil")
		}
		if modelTargetConcept == nil {
			return nil, errors.New("In addConceptView for link, modelTargetConcept is nil")
		}
		diagramSourceElement := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelSourceConcept, trans)
		if diagramSourceElement == nil {
			return nil, errors.New("In addConceptView for reference link, diagramSourceElement is nil")
		}
		diagramTargetElement := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, modelTargetConcept, trans)
		if diagramTargetElement == nil {
			return nil, errors.New("In addConceptView for reference link, diagramTargetElement is nil")
		}
		crldiagramdomain.SetLinkSource(newElement, diagramSourceElement, trans)
		crldiagramdomain.SetLinkTarget(newElement, diagramTargetElement, trans)
	} else {
		newElement, err = crldiagramdomain.NewDiagramNode(uOfD, trans)
		if err != nil {
			return nil, err
		}
		crldiagramdomain.SetNodeX(newElement, x, trans)
		crldiagramdomain.SetNodeY(newElement, y, trans)
		crldiagramdomain.SetLineColor(newElement, "#000000", trans)
	}

	err = newElement.SetLabel(el.GetLabel(trans), trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addConceptView failed")
	}
	crldiagramdomain.SetReferencedModelElement(newElement, el, trans)
	crldiagramdomain.SetDisplayLabel(newElement, el.GetLabel(trans), trans)

	err = newElement.SetOwningConceptID(diagram.GetConceptID(trans), trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addConceptView failed")
	}
	// err = newElement.Register(dMgr.elementManager)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "DiagramManager.addConceptView failed")
	// }

	return newElement, nil
}

// displayDiagram tells the client to display the indicated diagram.
func (dMgr *DiagramManager) DisplayDiagram(diagramID string, trans *core.Transaction) error {
	diagram := dMgr.editor.GetUofD().GetElement(diagramID)
	if diagram == nil {
		return errors.New("In DiagramManager.DisplayDiagram, the diagram does not exist")
	}
	if !diagram.IsRefinementOfURI(crldiagramdomain.CrlDiagramURI, trans) {
		return errors.New("In DiagramManager.DisplayDiagram, the supplied diagram is not a refinement of CrlDiagramURI")
	}
	// Make sure the diagram is in the list of displayed diagrams
	if !dMgr.editor.IsDiagramDisplayed(diagramID, trans) {
		err3 := dMgr.editor.addDiagramToDisplayedList(diagramID, trans)
		if err3 != nil {
			return errors.Wrap(err3, "DiagramManager.displayDiagram failed")
		}
	}
	for _, gui := range dMgr.editor.editorGUIs {
		err := gui.DisplayDiagram(diagram, trans)
		if err != nil {
			return errors.Wrap(err, "DiagramManager.DisplayDiagram failed")
		}
	}
	dMgr.editor.transientCurrentDiagram.SetLiteralValue(diagramID, trans)
	jsonOpenDiagrams, _ := json.Marshal(dMgr.editor.settings.OpenDiagrams)
	dMgr.editor.transientDisplayedDiagrams.SetLiteralValue(string(jsonOpenDiagrams), trans)
	dMgr.editor.settings.CurrentDiagram = diagramID
	return nil
}

// NewDiagram creates a new crldiagram
func (dMgr *DiagramManager) NewDiagram(trans *core.Transaction) (core.Element, error) {
	name := dMgr.editor.GetDefaultDiagramLabel()
	uOfD := dMgr.editor.GetUofD()
	diagram, err := crldiagramdomain.NewDiagram(uOfD, trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.newDiagram failed")
	}
	diagram.SetLabel(name, trans)
	dMgr.diagrams[diagram.GetConceptID(trans)] = diagram
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.newDiagram failed")
	}
	dMgr.DisplayDiagram(diagram.GetConceptID(trans), trans)
	return diagram, nil
}
