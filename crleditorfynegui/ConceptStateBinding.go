package crleditorfynegui

import (
	"errors"

	"fyne.io/fyne/v2/data/binding"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crleditor"
)

// ConceptStateBinding serves as an interface between a fyne data binding and an activeCRL concept
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
	oldURI          string
	dcl             *definitionChangeListener
	oldDefinition   string
	litcl           *literalChangeListener
	oldLiteralValue string
}

// NewConceptStateBinding creates a binding for the given id
func NewConceptStateBinding(id string) ConceptStateBinding {
	view := conceptStateBinding{}
	view.elementID = id
	el := crleditor.CrlEditorSingleton.GetUofD().GetElement(id)
	if el != nil {
		conceptState, _ := core.NewConceptState(el)
		view.rawData = *conceptState
		view.oldLabel = conceptState.Label
		view.oldURI = conceptState.URI
		view.boundData = binding.BindStruct(&view.rawData)
		view.lcl = newLabelChangeListener(&view)
		view.ucl = newURIChangeListener(&view)
		view.dcl = newDefinitionChangeListener(&view)
		view.litcl = newLiteralChangeListener(&view)
		el.Register(&view)
		view.boundData.AddListener(&view)
	}
	return &view
}

// GetBoundData returns the bound data for the binding
func (vPtr *conceptStateBinding) GetBoundData() *binding.Struct {
	return &vPtr.boundData
}

// Update updates the bound data based on the information in the notification
func (vPtr *conceptStateBinding) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
	// get the data from the notification
	afterState := notification.GetAfterConceptState()
	if afterState == nil {
		return nil
	}
	if afterState.ConceptID != vPtr.elementID {
		return errors.New("elementTreeNodeView.Update called with invalid notification")
	}
	vPtr.rawData = *afterState
	vPtr.oldLabel = afterState.Label
	vPtr.oldURI = afterState.URI
	vPtr.oldDefinition = afterState.Definition
	vPtr.oldLiteralValue = afterState.LiteralValue
	vPtr.boundData.Reload()
	return nil
}

// DataChanged is the callback invoked when the fyne binding data changes
func (vPtr *conceptStateBinding) DataChanged() {
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

// DataChanged is the fyne callback for changes in the Label widget
func (lcl *labelChangeListener) DataChanged() {
	labelItem, _ := lcl.parentBinding.boundData.GetItem("Label")
	newValue, _ := labelItem.(binding.String).Get()
	if newValue != lcl.parentBinding.oldLabel {
		lcl.parentBinding.oldLabel = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
		el := uOfD.GetElement(lcl.parentBinding.elementID)
		if el != nil {
			el.SetLabel(newValue, trans)
		}
	}
}

type uriChangeListener struct {
	parentBinding *conceptStateBinding
}

func newURIChangeListener(parentBinding *conceptStateBinding) *uriChangeListener {
	var ucl uriChangeListener
	ucl.parentBinding = parentBinding
	childField1Item, _ := parentBinding.boundData.GetItem("URI")
	childField1Item.AddListener(&ucl)
	return &ucl
}

// DataChanged is the fyne callback callback for the URI Entry widget
func (ucl *uriChangeListener) DataChanged() {
	labelItem, _ := ucl.parentBinding.boundData.GetItem("URI")
	newValue, _ := labelItem.(binding.String).Get()
	if newValue != ucl.parentBinding.oldURI {
		ucl.parentBinding.oldURI = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
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

// DataChanged is the fyne callback for changes to the Definition Entry widget
func (dcl *definitionChangeListener) DataChanged() {
	labelItem, _ := dcl.parentBinding.boundData.GetItem("Definition")
	newValue, _ := labelItem.(binding.String).Get()
	if newValue != dcl.parentBinding.oldDefinition {
		dcl.parentBinding.oldDefinition = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
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

// DataChanged is the fyne callback for changes to the Literal Value Entry widget
func (litcl *literalChangeListener) DataChanged() {
	labelItem, _ := litcl.parentBinding.boundData.GetItem("LiteralValue")
	newValue, _ := labelItem.(binding.String).Get()
	if newValue != litcl.parentBinding.oldLiteralValue {
		litcl.parentBinding.oldLiteralValue = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
		el := uOfD.GetElement(litcl.parentBinding.elementID)
		if el != nil {
			switch typedElement := el.(type) {
			case core.Literal:
				typedElement.SetLiteralValue(newValue, trans)
			}
		}
	}
}
