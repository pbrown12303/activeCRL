package crleditorfynegui

import (
	"fyne.io/x/fyne/widget/diagramwidget"
	"github.com/pbrown12303/activeCRL/core"
)

// AnchoredTextBinding is a bi-directional mapping transforming changes to and from the fyne-x AnchoredText
type AnchoredTextBinding interface {
	FyneAnchoredTextChanged(at *diagramwidget.AnchoredText)
	Update(notification *core.ChangeNotification, heldLocks *core.Transaction) error
}

// LinkLabelBinding maps the label on a link to the CRL DiagramLink label
type LinkLabelBinding struct {
	id               string
	fyneAnchoredText *diagramwidget.AnchoredText
	crlAnchoredText  *core.Concept
}

// NewLinkLabelBinding returns an initialized LinkLabelBinding
func NewLinkLabelBinding(id string, crlAnchoredText *core.Concept, fyneAnchoredText *diagramwidget.AnchoredText) (*LinkLabelBinding, error) {
	llb := &LinkLabelBinding{
		id:               id,
		fyneAnchoredText: fyneAnchoredText,
		crlAnchoredText:  crlAnchoredText,
	}
	return llb, nil
}

// FyneAnchoredTextChanged handles the notification that the fyne AnchoredText has changed
func (llb *LinkLabelBinding) FyneAnchoredTextChanged(at *diagramwidget.AnchoredText) {

}

// Update handles the notification that the crl AnchoredText has changed
func (llb *LinkLabelBinding) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
	return nil
}
