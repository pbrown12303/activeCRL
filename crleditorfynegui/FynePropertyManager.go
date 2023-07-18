package crleditorfynegui

import (
	"log"

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
	owningConceptIDValue                *shortcutableLabel
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
	referencedConceptLabel              *widget.Label
	referencedConceptValue              *widget.Label
	referencedConceptAttributeNameLabel *widget.Label
	referencedConceptAttributeNameValue *widget.Label
	referencedConceptVersionLabel       *widget.Label
	referencedConceptVersionValue       *widget.Label
	abstractConceptLabel                *widget.Label
	abstractConceptValue                *widget.Label
	abstractConceptVersionLabel         *widget.Label
	abstractConceptVersionValue         *widget.Label
	refinedConceptLabel                 *widget.Label
	refinedConceptValue                 *widget.Label
	refinedConceptVersionLabel          *widget.Label
	refinedConceptVersionValue          *widget.Label
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
	propertyManager.owningConceptIDValue = newFocusableLabel()
	propertyManager.definitionLabel = widget.NewLabel("Definition")
	propertyManager.definitionValue = widget.NewEntry()
	propertyManager.uriLabel = widget.NewLabel("URI")
	propertyManager.uriValue = widget.NewEntry()
	propertyManager.literalValueLabel = widget.NewLabel("Literal Value")
	propertyManager.literalValueValue = widget.NewEntry()
	propertyManager.referencedConceptLabel = widget.NewLabel("Referenced Concept ID")
	propertyManager.referencedConceptValue = widget.NewLabel("")
	propertyManager.referencedConceptAttributeNameLabel = widget.NewLabel("Referenced Attribute Name")
	propertyManager.referencedConceptAttributeNameValue = widget.NewLabel("")
	propertyManager.referencedConceptVersionLabel = widget.NewLabel("Referenced Concept Version")
	propertyManager.referencedConceptVersionValue = widget.NewLabel("")
	propertyManager.abstractConceptLabel = widget.NewLabel("Abstract Concept ID")
	propertyManager.abstractConceptValue = widget.NewLabel("")
	propertyManager.abstractConceptVersionLabel = widget.NewLabel("Abstract Concept Version")
	propertyManager.abstractConceptVersionValue = widget.NewLabel("")
	propertyManager.refinedConceptLabel = widget.NewLabel("Refined Concept ID")
	propertyManager.refinedConceptValue = widget.NewLabel("")
	propertyManager.refinedConceptVersionLabel = widget.NewLabel("Refined Concept Version")
	propertyManager.refinedConceptVersionValue = widget.NewLabel("")
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
		propertyManager.referencedConceptLabel,
		propertyManager.referencedConceptValue,
		propertyManager.referencedConceptAttributeNameLabel,
		propertyManager.referencedConceptAttributeNameValue,
		propertyManager.referencedConceptVersionLabel,
		propertyManager.referencedConceptVersionValue,
		propertyManager.abstractConceptLabel,
		propertyManager.abstractConceptValue,
		propertyManager.abstractConceptVersionLabel,
		propertyManager.abstractConceptVersionValue,
		propertyManager.refinedConceptLabel,
		propertyManager.refinedConceptValue,
		propertyManager.refinedConceptVersionLabel,
		propertyManager.refinedConceptVersionValue,
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
		pMgr.referencedConceptValue.Unbind()
		pMgr.referencedConceptValue.Text = ""
		pMgr.referencedConceptValue.Refresh()
		pMgr.abstractConceptValue.Unbind()
		pMgr.referencedConceptAttributeNameValue.Unbind()
		pMgr.referencedConceptAttributeNameValue.Text = ""
		pMgr.referencedConceptAttributeNameValue.Refresh()
		pMgr.referencedConceptVersionValue.Unbind()
		pMgr.referencedConceptVersionValue.Text = ""
		pMgr.referencedConceptVersionValue.Refresh()
		pMgr.abstractConceptValue.Unbind()
		pMgr.abstractConceptValue.Text = ""
		pMgr.abstractConceptValue.Refresh()
		pMgr.abstractConceptVersionValue.Unbind()
		pMgr.abstractConceptVersionValue.Text = ""
		pMgr.abstractConceptVersionValue.Refresh()
		pMgr.refinedConceptValue.Unbind()
		pMgr.refinedConceptValue.Text = ""
		pMgr.refinedConceptValue.Refresh()
		pMgr.refinedConceptVersionValue.Unbind()
		pMgr.refinedConceptVersionValue.Text = ""
		pMgr.refinedConceptVersionValue.Refresh()
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
		pMgr.referencedConceptValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("ReferencedAttributeName")
		pMgr.referencedConceptAttributeNameValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("ReferencedConceptVersion")
		pMgr.referencedConceptVersionValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("AbstractConceptID")
		pMgr.abstractConceptValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("AbstractConceptVersion")
		pMgr.abstractConceptVersionValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("RefinedConceptID")
		pMgr.refinedConceptValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("RefinedConceptVersion")
		pMgr.refinedConceptVersionValue.Bind(itemBinding.(binding.String))
		itemBinding, _ = structBinding.GetItem("LiteralValue")
		pMgr.literalValueValue.Bind(itemBinding.(binding.String))
	}
}

var _ fyne.Shortcutable = (*shortcutableLabel)(nil)
var _ fyne.Focusable = (*shortcutableLabel)(nil)
var _ fyne.Tappable = (*shortcutableLabel)(nil)

type shortcutableLabel struct {
	widget.Entry
}

func (sl *shortcutableLabel) FocusGained() {
	log.Print("Focus Gained")
}

func (sl *shortcutableLabel) FocusLost() {
	log.Print("Focus Lost")
}

func (sl *shortcutableLabel) Tapped(event *fyne.PointEvent) {
	log.Print("Tapped")
	FyneGUISingleton.GetWindow().RequestFocus()
}

func (sl *shortcutableLabel) TypedKey(*fyne.KeyEvent) {

}

func (sl *shortcutableLabel) TypedRune(rune) {

}

func newFocusableLabel() *shortcutableLabel {
	label := &shortcutableLabel{}
	label.ExtendBaseWidget(label)
	label.Disable()
	return label
}

func (sl *shortcutableLabel) TypedShortcut(shortcut fyne.Shortcut) {
	log.Print(shortcut.ShortcutName())
	switch typedShortcut := shortcut.(type) {
	case *fyne.ShortcutCopy:
		typedShortcut.Clipboard = FyneGUISingleton.window.Clipboard()
		typedShortcut.Clipboard.SetContent(sl.Text)
	}
}
