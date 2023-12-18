package crleditorfynegui

import (
	"errors"

	"fyne.io/fyne/v2/data/binding"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crleditor"
)

// ConceptStateBinding serves as an interface between fyne constructs and an activeCRL concept
type ConceptStateBinding struct {
	conceptID                      string
	conceptState                   core.ConceptState
	dcl                            *definitionChangeListener
	definitionBinding              binding.String
	lcl                            *labelChangeListener
	labelBinding                   binding.String
	owningConceptIDBinding         binding.String
	readOnlyBinding                binding.String
	ucl                            *uriChangeListener
	uriBinding                     binding.String
	versionBinding                 binding.String
	litcl                          *literalChangeListener
	literalValueBinding            binding.String
	referencedAttributeNameBinding binding.String
	referencedConceptIDBinding     binding.String
	abstractConceptIDBinding       binding.String
	refinedConceptIDBinding        binding.String
}

// NewConceptStateBinding creates a binding for the given id
func NewConceptStateBinding(id string) *ConceptStateBinding {
	view := ConceptStateBinding{}
	view.conceptID = id
	el := crleditor.CrlEditorSingleton.GetUofD().GetElement(id)
	if el != nil {
		conceptState, _ := core.NewConceptState(el)
		view.conceptState = *conceptState
		view.definitionBinding = binding.NewString()
		view.definitionBinding.Set(conceptState.Definition)
		view.dcl = newDefinitionChangeListener(&view)
		view.labelBinding = binding.NewString()
		view.labelBinding.Set(conceptState.Label)
		view.owningConceptIDBinding = binding.NewString()
		view.owningConceptIDBinding.Set(conceptState.OwningConceptID)
		view.readOnlyBinding = binding.NewString()
		view.readOnlyBinding.Set(conceptState.ReadOnly)
		view.uriBinding = binding.NewString()
		view.uriBinding.Set(conceptState.URI)
		view.versionBinding = binding.NewString()
		view.versionBinding.Set(conceptState.Version)
		view.ucl = newURIChangeListener(&view)
		view.lcl = newLabelChangeListener(&view)
		view.literalValueBinding = binding.NewString()
		view.literalValueBinding.Set(conceptState.LiteralValue)
		view.litcl = newLiteralChangeListener(&view)
		view.referencedAttributeNameBinding = binding.NewString()
		view.referencedAttributeNameBinding.Set(conceptState.ReferencedAttributeName)
		view.referencedConceptIDBinding = binding.NewString()
		view.referencedConceptIDBinding.Set(conceptState.ReferencedConceptID)
		view.abstractConceptIDBinding = binding.NewString()
		view.abstractConceptIDBinding.Set(conceptState.AbstractConceptID)
		view.refinedConceptIDBinding = binding.NewString()
		view.refinedConceptIDBinding.Set(conceptState.RefinedConceptID)
		el.Register(&view)
	}
	return &view
}

// Update updates the bound data based on the information in the notification
func (vPtr *ConceptStateBinding) Update(notification *core.ChangeNotification, trans *core.Transaction) error {
	// get the data from the notification
	afterState := notification.GetAfterConceptState()
	if afterState == nil {
		return nil
	}
	if afterState.ConceptID != vPtr.conceptID {
		return errors.New("elementTreeNodeView.Update called with invalid notification")
	}
	vPtr.conceptState = *afterState
	vPtr.labelBinding.Set(afterState.Label)
	vPtr.uriBinding.Set(afterState.URI)
	vPtr.definitionBinding.Set(afterState.Definition)
	vPtr.literalValueBinding.Set(afterState.LiteralValue)
	return nil
}

type labelChangeListener struct {
	csb *ConceptStateBinding
}

func newLabelChangeListener(csb *ConceptStateBinding) *labelChangeListener {
	var lcl labelChangeListener
	lcl.csb = csb
	csb.labelBinding.AddListener(&lcl)
	return &lcl
}

// DataChanged is the fyne callback for changes in the Label entry widget
func (lcl *labelChangeListener) DataChanged() {
	newValue, _ := lcl.csb.labelBinding.Get()
	if newValue != lcl.csb.conceptState.Label {
		lcl.csb.conceptState.Label = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
		el := uOfD.GetElement(lcl.csb.conceptID)
		if el != nil {
			el.SetLabel(newValue, trans)
		}
	}
}

type uriChangeListener struct {
	csb *ConceptStateBinding
}

func newURIChangeListener(csb *ConceptStateBinding) *uriChangeListener {
	var ucl uriChangeListener
	ucl.csb = csb
	ucl.csb.uriBinding.AddListener(&ucl)
	return &ucl
}

// DataChanged is the fyne callback callback for the URI Entry widget
func (ucl *uriChangeListener) DataChanged() {
	newValue, _ := ucl.csb.uriBinding.Get()
	if newValue != ucl.csb.conceptState.URI {
		ucl.csb.uriBinding.Set(newValue)
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
		el := uOfD.GetElement(ucl.csb.conceptID)
		if el != nil {
			el.SetURI(newValue, trans)
		}
	}
}

type definitionChangeListener struct {
	csb *ConceptStateBinding
}

func newDefinitionChangeListener(csb *ConceptStateBinding) *definitionChangeListener {
	var dcl definitionChangeListener
	dcl.csb = csb
	dcl.csb.definitionBinding.AddListener(&dcl)
	return &dcl
}

// DataChanged is the fyne callback for changes to the Definition Entry widget
func (dcl *definitionChangeListener) DataChanged() {
	newValue, _ := dcl.csb.definitionBinding.Get()
	if newValue != dcl.csb.conceptState.Definition {
		dcl.csb.conceptState.Definition = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
		el := uOfD.GetElement(dcl.csb.conceptID)
		if el != nil {
			el.SetDefinition(newValue, trans)
		}
	}
}

type literalChangeListener struct {
	csb *ConceptStateBinding
}

func newLiteralChangeListener(csb *ConceptStateBinding) *literalChangeListener {
	var litcl literalChangeListener
	litcl.csb = csb
	litcl.csb.literalValueBinding.AddListener(&litcl)
	return &litcl
}

// DataChanged is the fyne callback for changes to the Literal Value Entry widget
func (litcl *literalChangeListener) DataChanged() {
	newValue, _ := litcl.csb.literalValueBinding.Get()
	if newValue != litcl.csb.conceptState.LiteralValue {
		litcl.csb.conceptState.LiteralValue = newValue
		editor := crleditor.CrlEditorSingleton
		uOfD := editor.GetUofD()
		trans, isNew := FyneGUISingleton.editor.GetTransaction()
		if isNew {
			defer FyneGUISingleton.editor.EndTransaction()
		}
		el := uOfD.GetElement(litcl.csb.conceptID)
		if el != nil {
			switch el.GetConceptType() {
			case core.Literal:
				el.SetLiteralValue(newValue, trans)
			}
		}
	}
}
