package crleditor

import (
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pkg/errors"
)

// DiagramManager manages the diagram display portion of the GUI
type DiagramManager struct {
	editor   *Editor
	diagrams map[string]*crldiagramdomain.CrlDiagram
}

// NewDiagramManager creates an instance of the DiagramManager
func NewDiagramManager(editor *Editor) *DiagramManager {
	var dMgr DiagramManager
	dMgr.editor = editor
	dMgr.diagrams = map[string]*crldiagramdomain.CrlDiagram{}
	return &dMgr
}

// AddDiagram adds a diagram tab, if needed, and displays the tab
func (dMgr *DiagramManager) AddDiagram(ownerID string, trans *core.Transaction) (*crldiagramdomain.CrlDiagram, error) {
	diagram, err := dMgr.NewDiagram(trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = diagram.ToCore().SetOwningConceptID(ownerID, trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = dMgr.editor.SelectElement(diagram.ToCore(), trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = dMgr.DisplayDiagram(diagram.ToCore().GetConceptID(trans), trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	return diagram, nil
}

// AddConceptView adds a view of the concept to the indicated diagram
func (dMgr *DiagramManager) AddConceptView(diagramID string, conceptID string, x float64, y float64, trans *core.Transaction) (*crldiagramdomain.CrlDiagramElement, error) {
	uOfD := dMgr.editor.GetUofD()
	diagram := (*crldiagramdomain.CrlDiagram)(uOfD.GetElement(diagramID))
	el := uOfD.GetElement(conceptID)
	if el == nil {
		return nil, errors.New("Indicated model element not found in addNodeView, ID: " + conceptID)
	}
	createAsLink := false
	switch el.GetConceptType() {
	case core.Reference:
		createAsLink = dMgr.editor.GetDropDiagramReferenceAsLink(trans)
	case core.Refinement:
		createAsLink = dMgr.editor.GetDropDiagramRefinementAsLink(trans)
	}

	var newElement *crldiagramdomain.CrlDiagramElement
	var err error
	if createAsLink {
		var newLink *crldiagramdomain.CrlDiagramLink
		var modelSourceConcept *core.Concept
		var modelTargetConcept *core.Concept
		switch el.GetConceptType() {
		case core.Reference:
			newLink, err = crldiagramdomain.NewDiagramReferenceLink(trans)
			if err != nil {
				return nil, err
			}
			reference := el
			modelSourceConcept = reference.GetOwningConcept(trans)
			modelTargetConcept = reference.GetReferencedConcept(trans)
		case core.Refinement:
			newLink, err = crldiagramdomain.NewDiagramRefinementLink(trans)
			if err != nil {
				return nil, err
			}
			refinement := el
			modelSourceConcept = refinement.GetRefinedConcept(trans)
			modelTargetConcept = refinement.GetAbstractConcept(trans)
		}
		if modelSourceConcept == nil {
			return nil, errors.New("In addConceptView for link, modelSourceConcept is nil")
		}
		if modelTargetConcept == nil {
			return nil, errors.New("In addConceptView for link, modelTargetConcept is nil")
		}
		diagramSourceElement := diagram.GetFirstElementRepresentingConcept(modelSourceConcept, trans)
		if diagramSourceElement == nil {
			return nil, errors.New("In addConceptView for reference link, diagramSourceElement is nil")
		}
		diagramTargetElement := diagram.GetFirstElementRepresentingConcept(modelTargetConcept, trans)
		if diagramTargetElement == nil {
			return nil, errors.New("In addConceptView for reference link, diagramTargetElement is nil")
		}
		newLink.SetLinkSource(diagramSourceElement, trans)
		newLink.SetLinkTarget(diagramTargetElement, trans)
		newElement = newLink.ToCrlDiagramElement()
	} else {
		newNode, err := crldiagramdomain.NewDiagramNode(trans)
		if err != nil {
			return nil, err
		}
		newNode.SetNodeX(x, trans)
		newNode.SetNodeY(y, trans)
		newElement = newNode.ToCrlDiagramElement()
		newElement.SetLineColor("#000000", trans)
	}

	err = newElement.ToCore().SetLabel(el.GetLabel(trans), trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addConceptView failed")
	}
	newElement.SetReferencedModelConcept(el, trans)
	newElement.SetDisplayLabel(el.GetLabel(trans), trans)

	newElement.SetDiagram(diagram, trans)
	// err = newElement.Register(dMgr.elementManager)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "DiagramManager.addConceptView failed")
	// }

	return newElement, nil
}

// DisplayDiagram tells the client to display the indicated diagram.
func (dMgr *DiagramManager) DisplayDiagram(diagramID string, trans *core.Transaction) error {
	diagram := crldiagramdomain.GetCrlDiagram(diagramID, trans)
	if diagram == nil {
		return errors.New("In DiagramManager.DisplayDiagram, the diagram does not exist")
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
	if !dMgr.editor.undoRedoInProgress {
		dMgr.editor.transientCurrentDiagram.SetLiteralValue(diagramID, trans)
	}
	return nil
}

// NewDiagram creates a new crldiagram
func (dMgr *DiagramManager) NewDiagram(trans *core.Transaction) (*crldiagramdomain.CrlDiagram, error) {
	name := dMgr.editor.GetDefaultDiagramLabel()
	diagram, err := crldiagramdomain.NewDiagram(trans)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.newDiagram failed")
	}
	diagram.ToCore().SetLabel(name, trans)
	dMgr.diagrams[diagram.ToCore().GetConceptID(trans)] = diagram
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.newDiagram failed")
	}
	dMgr.DisplayDiagram(diagram.ToCore().GetConceptID(trans), trans)
	return diagram, nil
}
