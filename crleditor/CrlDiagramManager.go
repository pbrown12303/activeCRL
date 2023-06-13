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

func (dMgr *DiagramManager) AddDiagram(ownerID string, hl *core.Transaction) (core.Element, error) {
	diagram, err := dMgr.NewDiagram(hl)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = diagram.SetOwningConceptID(ownerID, hl)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = dMgr.editor.SelectElement(diagram, hl)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	err = dMgr.DisplayDiagram(diagram.GetConceptID(hl), hl)
	if err != nil {
		return nil, errors.Wrap(err, "DiagramManager.addDiagram failed")
	}
	return diagram, nil
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
			return errors.Wrap(err3, "diagramManager.displayDiagram failed")
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
		return nil, errors.Wrap(err, "diagramManager.newDiagram failed")
	}
	diagram.SetLabel(name, trans)
	dMgr.diagrams[diagram.GetConceptID(trans)] = diagram
	if err != nil {
		return nil, errors.Wrap(err, "diagramManager.newDiagram failed")
	}
	dMgr.DisplayDiagram(diagram.GetConceptID(trans), trans)
	return diagram, nil
}
