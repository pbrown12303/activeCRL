package fynegui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FynePropertyManager struct {
	properties             *fyne.Container
	propertyHeading        *widget.Label
	valueHeading           *widget.Label
	typeLabel              *widget.Label
	typeValue              *widget.Label
	idLabel                *widget.Label
	idValue                *widget.Label
	owningConceptIDLabel   *widget.Label
	owningConceptIDValue   *widget.Label
	versionLabel           *widget.Label
	versionValue           *widget.Label
	labelLabel             *widget.Label
	labelValue             *widget.Entry
	definitionLabel        *widget.Label
	definitionValue        *widget.Entry
	uriLabel               *widget.Label
	uriValue               *widget.Entry
	readOnlyLabel          *widget.Label
	readOnlyValue          *widget.Label
	referencedElementLabel *widget.Label
	referencedElementValue *widget.Label
	abstractElementLabel   *widget.Label
	abstractElementValue   *widget.Label
	refinedElementLabel    *widget.Label
	refinedElementValue    *widget.Label
	literalValueLabel      *widget.Label
	literalValueValue      *widget.Entry
}

func NewFynePropertyManager() *FynePropertyManager {
	var propertyManager FynePropertyManager
	propertyManager.propertyHeading = widget.NewLabel("Property")
	propertyManager.valueHeading = widget.NewLabel("Value")
	propertyManager.typeLabel = widget.NewLabel("Type")
	propertyManager.typeValue = widget.NewLabel("")
	propertyManager.idLabel = widget.NewLabel("ID")
	propertyManager.idValue = widget.NewLabel("")
	propertyManager.owningConceptIDLabel = widget.NewLabel("Owning Concept ID")
	propertyManager.owningConceptIDValue = widget.NewLabel("")
	propertyManager.versionLabel = widget.NewLabel("Version")
	propertyManager.versionValue = widget.NewLabel("")
	propertyManager.labelLabel = widget.NewLabel("Label")
	propertyManager.labelValue = widget.NewEntry()
	propertyManager.definitionLabel = widget.NewLabel("Definition")
	propertyManager.definitionValue = widget.NewEntry()
	propertyManager.uriLabel = widget.NewLabel("URI")
	propertyManager.uriValue = widget.NewEntry()
	propertyManager.readOnlyLabel = widget.NewLabel("Read Only")
	propertyManager.readOnlyValue = widget.NewLabel("")
	propertyManager.referencedElementLabel = widget.NewLabel("Referenced Element ID")
	propertyManager.referencedElementValue = widget.NewLabel("")
	propertyManager.abstractElementLabel = widget.NewLabel("Abstract Element ID")
	propertyManager.abstractElementValue = widget.NewLabel("")
	propertyManager.refinedElementLabel = widget.NewLabel("Refined Element ID")
	propertyManager.refinedElementValue = widget.NewLabel("")
	propertyManager.literalValueLabel = widget.NewLabel("Literal Value")
	propertyManager.literalValueValue = widget.NewEntry()

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
		propertyManager.versionLabel,
		propertyManager.versionValue,
		propertyManager.labelLabel,
		propertyManager.labelValue,
		propertyManager.definitionLabel,
		propertyManager.definitionValue,
		propertyManager.uriLabel,
		propertyManager.uriValue,
		propertyManager.readOnlyLabel,
		propertyManager.readOnlyValue,
		propertyManager.referencedElementLabel,
		propertyManager.referencedElementValue,
		propertyManager.abstractElementLabel,
		propertyManager.abstractElementValue,
		propertyManager.refinedElementLabel,
		propertyManager.refinedElementValue,
		propertyManager.literalValueLabel,
		propertyManager.literalValueValue)

	return &propertyManager
}
