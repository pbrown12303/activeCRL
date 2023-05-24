package crleditorfynegui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FynePropertyManager struct {
	properties                          *fyne.Container
	propertyHeading                     *widget.Label
	valueHeading                        *widget.Label
	typeLabel                           *widget.Label
	typeValue                           *widget.Label
	idLabel                             *widget.Label
	idValue                             *widget.Label
	owningConceptIDLabel                *widget.Label
	owningConceptIDValue                *widget.Label
	versionLabel                        *widget.Label
	versionValue                        *widget.Label
	labelLabel                          *widget.Label
	labelValue                          *widget.Entry
	definitionLabel                     *widget.Label
	definitionValue                     *widget.Entry
	uriLabel                            *widget.Label
	uriValue                            *widget.Entry
	isCoreLabel                         *widget.Label
	isCoreValue                         *widget.Label
	readOnlyLabel                       *widget.Label
	readOnlyValue                       *widget.Label
	referencedElementLabel              *widget.Label
	referencedElementValue              *widget.Label
	referencedElementAttributeNameLabel *widget.Label
	referencedElementAttributeNameValue *widget.Label
	referencedElementVersionLabel       *widget.Label
	referencedElementVersionValue       *widget.Label
	abstractElementLabel                *widget.Label
	abstractElementValue                *widget.Label
	abstractElementVersionLabel         *widget.Label
	abstractElementVersionValue         *widget.Label
	refinedElementLabel                 *widget.Label
	refinedElementValue                 *widget.Label
	refinedElementVersionLabel          *widget.Label
	refinedElementVersionValue          *widget.Label
	literalValueLabel                   *widget.Label
	literalValueValue                   *widget.Entry
}

func NewFynePropertyManager() *FynePropertyManager {
	var propertyManager FynePropertyManager
	propertyManager.propertyHeading = widget.NewLabel("Property")
	propertyManager.propertyHeading.TextStyle.Bold = true
	propertyManager.valueHeading = widget.NewLabel("Value")
	propertyManager.valueHeading.TextStyle.Bold = true
	propertyManager.typeLabel = widget.NewLabel("Type")
	propertyManager.typeValue = widget.NewLabel("")
	propertyManager.idLabel = widget.NewLabel("ID")
	propertyManager.idValue = widget.NewLabel("")
	propertyManager.labelLabel = widget.NewLabel("Label")
	propertyManager.labelValue = widget.NewEntry()
	propertyManager.owningConceptIDLabel = widget.NewLabel("Owning Concept ID")
	propertyManager.owningConceptIDValue = widget.NewLabel("")
	propertyManager.definitionLabel = widget.NewLabel("Definition")
	propertyManager.definitionValue = widget.NewEntry()
	propertyManager.uriLabel = widget.NewLabel("URI")
	propertyManager.uriValue = widget.NewEntry()
	propertyManager.literalValueLabel = widget.NewLabel("Literal Value")
	propertyManager.literalValueValue = widget.NewEntry()
	propertyManager.referencedElementLabel = widget.NewLabel("Referenced Element ID")
	propertyManager.referencedElementValue = widget.NewLabel("")
	propertyManager.referencedElementAttributeNameLabel = widget.NewLabel("Referenced Attribute Name")
	propertyManager.referencedElementAttributeNameValue = widget.NewLabel("")
	propertyManager.referencedElementVersionLabel = widget.NewLabel("Referenced Element Version")
	propertyManager.referencedElementVersionValue = widget.NewLabel("")
	propertyManager.abstractElementLabel = widget.NewLabel("Abstract Element ID")
	propertyManager.abstractElementValue = widget.NewLabel("")
	propertyManager.abstractElementVersionLabel = widget.NewLabel("Abstract Element Version")
	propertyManager.abstractElementVersionValue = widget.NewLabel("")
	propertyManager.refinedElementLabel = widget.NewLabel("Refined Element ID")
	propertyManager.refinedElementValue = widget.NewLabel("")
	propertyManager.refinedElementVersionLabel = widget.NewLabel("Refined Element Version")
	propertyManager.refinedElementVersionValue = widget.NewLabel("")
	propertyManager.isCoreLabel = widget.NewLabel("Is Core")
	propertyManager.isCoreValue = widget.NewLabel("")
	propertyManager.readOnlyLabel = widget.NewLabel("Read Only")
	propertyManager.readOnlyValue = widget.NewLabel("")
	propertyManager.versionLabel = widget.NewLabel("Version")
	propertyManager.versionValue = widget.NewLabel("")

	// Properties
	propertyManager.properties = container.New(layout.NewGridLayout(2),
		propertyManager.propertyHeading,
		propertyManager.valueHeading,
		propertyManager.typeLabel,
		propertyManager.typeValue,
		propertyManager.idLabel,
		propertyManager.idValue,
		propertyManager.owningConceptIDLabel,
		propertyManager.owningConceptIDValue,
		propertyManager.labelLabel,
		propertyManager.labelValue,
		propertyManager.definitionLabel,
		propertyManager.definitionValue,
		propertyManager.literalValueLabel,
		propertyManager.literalValueValue,
		propertyManager.uriLabel,
		propertyManager.uriValue,
		propertyManager.referencedElementLabel,
		propertyManager.referencedElementValue,
		propertyManager.referencedElementAttributeNameLabel,
		propertyManager.referencedElementAttributeNameValue,
		propertyManager.referencedElementVersionLabel,
		propertyManager.referencedElementVersionValue,
		propertyManager.abstractElementLabel,
		propertyManager.abstractElementValue,
		propertyManager.abstractElementVersionLabel,
		propertyManager.abstractElementVersionValue,
		propertyManager.refinedElementLabel,
		propertyManager.refinedElementValue,
		propertyManager.refinedElementVersionLabel,
		propertyManager.refinedElementVersionValue,
		propertyManager.isCoreLabel,
		propertyManager.isCoreValue,
		propertyManager.readOnlyLabel,
		propertyManager.readOnlyValue,
		propertyManager.versionLabel,
		propertyManager.versionValue)
	return &propertyManager
}

func (pMgr *FynePropertyManager) displayProperties(uid string) {
	conceptBinding := GetConceptStateBinding(uid)
	if uid == "" || conceptBinding == nil {
		pMgr.typeValue.Unbind()
		pMgr.typeValue.Text = ""
		pMgr.typeValue.Refresh()
		pMgr.idValue.Unbind()
		pMgr.idValue.Text = ""
		pMgr.idValue.Refresh()
		pMgr.owningConceptIDValue.Unbind()
		pMgr.owningConceptIDValue.Text = ""
		pMgr.owningConceptIDValue.Refresh()
		pMgr.versionValue.Unbind()
		pMgr.versionValue.Text = ""
		pMgr.versionValue.Refresh()
		pMgr.labelValue.Unbind()
		pMgr.labelValue.Text = ""
		pMgr.labelValue.Refresh()
		pMgr.definitionValue.Unbind()
		pMgr.definitionValue.Text = ""
		pMgr.definitionValue.Refresh()
		pMgr.uriValue.Unbind()
		pMgr.uriValue.Text = ""
		pMgr.uriValue.Refresh()
		pMgr.isCoreValue.Unbind()
		pMgr.isCoreValue.Text = ""
		pMgr.isCoreValue.Refresh()
		pMgr.readOnlyValue.Unbind()
		pMgr.readOnlyValue.Text = ""
		pMgr.readOnlyValue.Refresh()
		pMgr.referencedElementValue.Unbind()
		pMgr.referencedElementValue.Text = ""
		pMgr.referencedElementValue.Refresh()
		pMgr.abstractElementValue.Unbind()
		pMgr.abstractElementValue.Text = ""
		pMgr.abstractElementValue.Refresh()
		pMgr.refinedElementValue.Unbind()
		pMgr.refinedElementValue.Text = ""
		pMgr.refinedElementValue.Refresh()
		pMgr.literalValueValue.Unbind()
		pMgr.literalValueValue.Text = ""
		pMgr.literalValueValue.Refresh()
	} else {
		structBinding := *conceptBinding.GetBoundData()
		itemBinding, _ := structBinding.GetItem("ConceptType")
		pMgr.typeValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("ConceptID")
		pMgr.idValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("OwningConceptID")
		pMgr.owningConceptIDValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("Version")
		pMgr.versionValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("Label")
		pMgr.labelValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("Definition")
		pMgr.definitionValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("URI")
		pMgr.uriValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("IsCore")
		pMgr.isCoreValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("ReadOnly")
		pMgr.readOnlyValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("ReferencedConceptID")
		pMgr.referencedElementValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("AbstractConceptID")
		pMgr.abstractElementValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("RefinedConceptID")
		pMgr.refinedElementValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("LiteralValue")
		pMgr.literalValueValue.Bind(itemBinding.(binding.String))
	}

}
