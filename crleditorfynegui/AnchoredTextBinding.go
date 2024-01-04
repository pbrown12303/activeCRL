package crleditorfynegui

import (
	"sync"

	"fyne.io/x/fyne/widget/diagramwidget"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
)

// AnchoredTextBinding is a bi-directional mapping transforming changes to and from the fyne-x AnchoredText
type AnchoredTextBinding struct {
	mu               sync.Mutex
	updateInProgress bool
	id               string
	fyneAnchoredText *diagramwidget.AnchoredText
	crlAnchoredText  *crldiagramdomain.CrlDiagramAnchoredText
}

// NewAnchoredTextBinding returns an initialized LinkLabelBinding
func NewAnchoredTextBinding(id string, crlAnchoredText *crldiagramdomain.CrlDiagramAnchoredText, fyneAnchoredText *diagramwidget.AnchoredText) (*AnchoredTextBinding, error) {
	llb := &AnchoredTextBinding{
		id:               id,
		fyneAnchoredText: fyneAnchoredText,
		crlAnchoredText:  crlAnchoredText,
		updateInProgress: true,
	}
	FyneGUISingleton.SetAnchoredTextBinding(id, llb)
	crlAnchoredText.AsCore().Register(llb)
	llb.updateInProgress = false
	return llb, nil
}

func (llb *AnchoredTextBinding) delete() {
	if llb.crlAnchoredText != nil {
		llb.crlAnchoredText.AsCore().Deregister(llb)
	}
	FyneGUISingleton.SetAnchoredTextBinding(llb.id, nil)
}

// FyneAnchoredTextChanged handles the notification that the fyne AnchoredText has changed
func (llb *AnchoredTextBinding) FyneAnchoredTextChanged(at *diagramwidget.AnchoredText) {
	llb.mu.Lock()
	if llb.updateInProgress {
		llb.mu.Unlock()
		return
	}
	llb.updateInProgress = true
	llb.mu.Unlock()
	trans, new := FyneGUISingleton.editor.GetTransaction()
	if new {
		defer FyneGUISingleton.editor.EndTransaction()
	}
	if llb.fyneAnchoredText.GetTextEntry().Text != llb.crlAnchoredText.LiteralValue {
		llb.crlAnchoredText.AsCore().SetLiteralValue(llb.fyneAnchoredText.GetTextEntry().Text, trans)
	}
	fyneOffsetX, fyneOffsetY := llb.fyneAnchoredText.GetOffset()
	crlOffsetX := llb.crlAnchoredText.GetOffsetX(trans)
	crlOffsetY := llb.crlAnchoredText.GetOffsetY(trans)
	if fyneOffsetX != crlOffsetX {
		llb.crlAnchoredText.SetOffsetX(fyneOffsetX, trans)
	}
	if fyneOffsetY != crlOffsetY {
		llb.crlAnchoredText.SetOffsetY(fyneOffsetY, trans)
	}
	fyneAnchorPosition := llb.fyneAnchoredText.Position()
	crlAnchorX := llb.crlAnchoredText.GetAnchorX(trans)
	crlAnchorY := llb.crlAnchoredText.GetAnchorY(trans)
	if fyneAnchorPosition.X != float32(crlAnchorX) {
		llb.crlAnchoredText.SetAnchorX(float64(fyneAnchorPosition.X), trans)
	}
	if fyneAnchorPosition.Y != float32(crlAnchorY) {
		llb.crlAnchoredText.SetAnchorY(float64(fyneAnchorPosition.Y), trans)
	}
	llb.mu.Lock()
	llb.updateInProgress = false
	llb.mu.Unlock()
}

// Update handles the notification that the crl AnchoredText has changed
func (llb *AnchoredTextBinding) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
	llb.mu.Lock()
	if llb.updateInProgress {
		llb.mu.Unlock()
		return nil
	}
	llb.updateInProgress = true
	llb.mu.Unlock()
	if llb.fyneAnchoredText.GetTextEntry().Text != llb.crlAnchoredText.LiteralValue {
		llb.fyneAnchoredText.SetText(llb.crlAnchoredText.LiteralValue)
	}
	fyneOffsetX, fyneOffsetY := llb.fyneAnchoredText.GetOffset()
	crlOffsetX := llb.crlAnchoredText.GetOffsetX(trans)
	crlOffsetY := llb.crlAnchoredText.GetOffsetY(trans)
	if fyneOffsetX != crlOffsetX || fyneOffsetY != crlOffsetY {
		llb.fyneAnchoredText.SetOffsetNoCallback(crlOffsetX, crlOffsetY)
	}
	llb.mu.Lock()
	llb.updateInProgress = false
	llb.mu.Unlock()
	return nil
}
