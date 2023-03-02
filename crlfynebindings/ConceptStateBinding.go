package crlfynebindings

import (
	"errors"
	"log"

	"fyne.io/fyne/v2/data/binding"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crleditor"
)

var conceptStateBindingMap map[string]ConceptStateBinding

func InitBindings() {
	conceptStateBindingMap = make(map[string]ConceptStateBinding)
}

func GetConceptStateBinding(uid string) ConceptStateBinding {
	binding := conceptStateBindingMap[uid]
	if binding == nil {
		binding = NewConceptStateBinding(uid)
		conceptStateBindingMap[uid] = binding
	}
	return binding
}

type ConceptStateBinding interface {
	GetBoundData() *binding.Struct
	Update(*core.ChangeNotification, *core.Transaction) error
}

type conceptStateBinding struct {
	elementID string
	rawData   core.ConceptState
	boundData binding.Struct
}

func NewConceptStateBinding(id string) ConceptStateBinding {
	view := conceptStateBinding{}
	view.elementID = id
	el := crleditor.CrlEditorSingleton.GetUofD().GetElement(id)
	if el != nil {
		conceptState, _ := core.NewConceptState(el)
		view.rawData = *conceptState
		view.boundData = binding.BindStruct(&view.rawData)
		el.Register(&view)
		view.boundData.AddListener(&view)
	}
	return &view
}

func (vPtr *conceptStateBinding) GetBoundData() *binding.Struct {
	return &vPtr.boundData
}

func (vPtr *conceptStateBinding) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
	// get the data from the notification
	switch notification.GetNatureOfChange() {
	case core.ConceptChanged:
		afterState := notification.GetAfterConceptState()
		if afterState == nil || afterState.ConceptID != vPtr.elementID {
			return errors.New("elementTreeNodeView.Update called with invalid notification")
		}

		// Debugging
		item, _ := vPtr.boundData.GetItem("Label")
		v, _ := item.(binding.String).Get()
		log.Print("Label before assignment: " + v)

		vPtr.rawData = *afterState
		vPtr.boundData.Reload()

		// Debugging
		item, _ = vPtr.boundData.GetItem("Label")
		v, _ = item.(binding.String).Get()
		log.Print("Label after assignment: " + v)

	}
	return nil
}

func (vPtr *conceptStateBinding) DataChanged() {
	log.Print("DataChanged called")
}
