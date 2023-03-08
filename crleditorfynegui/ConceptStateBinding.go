package crleditorfynegui

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
	elementID       string
	rawData         core.ConceptState
	boundData       binding.Struct
	lcl             *labelChangeListener
	oldLabel        string
	ucl             *uriChangeListener
	oldUri          string
	dcl             *definitionChangeListener
	oldDefinition   string
	litcl           *literalChangeListener
	oldLiteralValue string
}

func NewConceptStateBinding(id string) ConceptStateBinding {
	view := conceptStateBinding{}
	view.elementID = id
	el := crleditor.CrlEditorSingleton.GetUofD().GetElement(id)
	if el != nil {
		conceptState, _ := core.NewConceptState(el)
		view.rawData = *conceptState
		view.oldLabel = conceptState.Label
		view.oldUri = conceptState.URI
		view.boundData = binding.BindStruct(&view.rawData)
		view.lcl = newLabelChangeListener(&view)
		view.ucl = newUriChangeListener(&view)
		view.dcl = newDefinitionChangeListener(&view)
		view.litcl = newLiteralChangeListener(&view)
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
		vPtr.rawData = *afterState
		vPtr.oldLabel = afterState.Label
		vPtr.oldUri = afterState.URI
		vPtr.oldDefinition = afterState.Definition
		vPtr.oldLiteralValue = afterState.LiteralValue
		vPtr.boundData.Reload()
	}
	return nil
}

func (vPtr *conceptStateBinding) DataChanged() {
	log.Print("DataChanged called")
}

type labelChangeListener struct {
	parentBinding *conceptStateBinding
}

func newLabelChangeListener(parentBinding *conceptStateBinding) *labelChangeListener {
	var lcl labelChangeListener
	lcl.parentBinding = parentBinding
	childField1Item, _ := parentBinding.boundData.GetItem("Label")
	childField1Item.AddListener(&lcl)
	return &lcl
}

func (lcl *labelChangeListener) DataChanged() {
	labelItem, _ := lcl.parentBinding.boundData.GetItem("Label")
	newValue, _ := labelItem.(binding.String).Get()
	if newValue != lcl.parentBinding.oldLabel {
		lcl.parentBinding.oldLabel = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans := GetTransaction()
		defer trans.ReleaseLocks()
		el := uOfD.GetElement(lcl.parentBinding.elementID)
		if el != nil {
			el.SetLabel(newValue, trans)
		}
	}
}

type uriChangeListener struct {
	parentBinding *conceptStateBinding
}

func newUriChangeListener(parentBinding *conceptStateBinding) *uriChangeListener {
	var ucl uriChangeListener
	ucl.parentBinding = parentBinding
	childField1Item, _ := parentBinding.boundData.GetItem("URI")
	childField1Item.AddListener(&ucl)
	return &ucl
}

func (ucl *uriChangeListener) DataChanged() {
	labelItem, _ := ucl.parentBinding.boundData.GetItem("URI")
	newValue, _ := labelItem.(binding.String).Get()
	if newValue != ucl.parentBinding.oldUri {
		ucl.parentBinding.oldUri = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans := GetTransaction()
		defer trans.ReleaseLocks()
		el := uOfD.GetElement(ucl.parentBinding.elementID)
		if el != nil {
			el.SetURI(newValue, trans)
		}
	}
}

type definitionChangeListener struct {
	parentBinding *conceptStateBinding
}

func newDefinitionChangeListener(parentBinding *conceptStateBinding) *definitionChangeListener {
	var dcl definitionChangeListener
	dcl.parentBinding = parentBinding
	childField1Item, _ := parentBinding.boundData.GetItem("Definition")
	childField1Item.AddListener(&dcl)
	return &dcl
}

func (dcl *definitionChangeListener) DataChanged() {
	labelItem, _ := dcl.parentBinding.boundData.GetItem("Definition")
	newValue, _ := labelItem.(binding.String).Get()
	if newValue != dcl.parentBinding.oldDefinition {
		dcl.parentBinding.oldDefinition = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans := GetTransaction()
		defer trans.ReleaseLocks()
		el := uOfD.GetElement(dcl.parentBinding.elementID)
		if el != nil {
			el.SetDefinition(newValue, trans)
		}
	}
}

type literalChangeListener struct {
	parentBinding *conceptStateBinding
}

func newLiteralChangeListener(parentBinding *conceptStateBinding) *literalChangeListener {
	var litcl literalChangeListener
	litcl.parentBinding = parentBinding
	childField1Item, _ := parentBinding.boundData.GetItem("LiteralValue")
	childField1Item.AddListener(&litcl)
	return &litcl
}

func (litcl *literalChangeListener) DataChanged() {
	labelItem, _ := litcl.parentBinding.boundData.GetItem("LiteralValue")
	newValue, _ := labelItem.(binding.String).Get()
	if newValue != litcl.parentBinding.oldLiteralValue {
		litcl.parentBinding.oldLiteralValue = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans := GetTransaction()
		defer trans.ReleaseLocks()
		el := uOfD.GetElement(litcl.parentBinding.elementID)
		if el != nil {
			switch typedElement := el.(type) {
			case core.Literal:
				typedElement.SetLiteralValue(newValue, trans)
			}
		}
	}
}
