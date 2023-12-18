package crleditorfynegui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// FynePropertyManager manages the property display in the interface
type FynePropertyManager struct {
	properties                          *fyne.Container
	propertyHeading                     *widget.Label
	valueHeading                        *widget.Label
	typeLabel                           *widget.Label
	typeValue                           *widget.Label
	idLabel                             *widget.Label
	idValue                             *copyableLabel
	owningConceptIDLabel                *widget.Label
	owningConceptIDValue                *copyableLabel
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
	referencedConceptValue              *copyableLabel
	referencedConceptAttributeNameLabel *widget.Label
	referencedConceptAttributeNameValue *widget.Label
	abstractConceptLabel                *widget.Label
	abstractConceptValue                *copyableLabel
	refinedConceptLabel                 *widget.Label
	refinedConceptValue                 *copyableLabel
	literalValueLabel                   *widget.Label
	literalValueValue                   *widget.Entry
}

// NewFynePropertyManager creates an initialized instance of the FynePropertyManager
func NewFynePropertyManager() *FynePropertyManager {
	var propertyManager FynePropertyManager
	propertyManager.propertyHeading = widget.NewLabel("Property")
	propertyManager.propertyHeading.TextStyle.Bold = true
	propertyManager.valueHeading = widget.NewLabel("Value")
	propertyManager.valueHeading.TextStyle.Bold = true
	propertyManager.typeLabel = widget.NewLabel("Type")
	propertyManager.typeValue = widget.NewLabel("")
	propertyManager.idLabel = widget.NewLabel("ID")
	propertyManager.idValue = newCopyableLabel()
	propertyManager.labelLabel = widget.NewLabel("Label")
	propertyManager.labelValue = widget.NewEntry()
	propertyManager.owningConceptIDLabel = widget.NewLabel("Owning Concept ID")
	propertyManager.owningConceptIDValue = newCopyableLabel()
	propertyManager.definitionLabel = widget.NewLabel("Definition")
	propertyManager.definitionValue = widget.NewEntry()
	propertyManager.uriLabel = widget.NewLabel("URI")
	propertyManager.uriValue = widget.NewEntry()
	propertyManager.literalValueLabel = widget.NewLabel("Literal Value")
	propertyManager.literalValueValue = widget.NewEntry()
	propertyManager.referencedConceptLabel = widget.NewLabel("Referenced Concept ID")
	propertyManager.referencedConceptValue = newCopyableLabel()
	propertyManager.referencedConceptAttributeNameLabel = widget.NewLabel("Referenced Attribute Name")
	propertyManager.referencedConceptAttributeNameValue = widget.NewLabel("")
	propertyManager.abstractConceptLabel = widget.NewLabel("Abstract Concept ID")
	propertyManager.abstractConceptValue = newCopyableLabel()
	propertyManager.refinedConceptLabel = widget.NewLabel("Refined Concept ID")
	propertyManager.refinedConceptValue = newCopyableLabel()
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
		propertyManager.abstractConceptLabel,
		propertyManager.abstractConceptValue,
		propertyManager.refinedConceptLabel,
		propertyManager.refinedConceptValue,
		propertyManager.isCoreLabel,
		propertyManager.isCoreValue,
		propertyManager.readOnlyLabel,
		propertyManager.readOnlyValue,
		propertyManager.versionLabel,
		propertyManager.versionValue)
	return &propertyManager
}

func (pMgr *FynePropertyManager) displayProperties(uid string) {
	csb := FyneGUISingleton.GetConceptStateBinding(uid)
	if uid == "" || csb == nil {
		pMgr.typeValue.SetText("")
		pMgr.idValue.SetText("")
		pMgr.owningConceptIDValue.Unbind()
		pMgr.owningConceptIDValue.SetText("")
		pMgr.versionValue.Unbind()
		pMgr.versionValue.SetText("")
		pMgr.labelValue.Unbind()
		pMgr.labelValue.SetText("")
		pMgr.definitionValue.Unbind()
		pMgr.definitionValue.SetText("")
		pMgr.uriValue.Unbind()
		pMgr.uriValue.SetText("")
		pMgr.isCoreValue.SetText("")
		pMgr.readOnlyValue.Unbind()
		pMgr.readOnlyValue.SetText("")
		pMgr.referencedConceptValue.Unbind()
		pMgr.referencedConceptValue.SetText("")
		pMgr.abstractConceptValue.Unbind()
		pMgr.abstractConceptValue.SetText("")
		pMgr.referencedConceptAttributeNameValue.Unbind()
		pMgr.referencedConceptAttributeNameValue.SetText("")
		pMgr.abstractConceptValue.Unbind()
		pMgr.abstractConceptValue.SetText("")
		pMgr.refinedConceptValue.Unbind()
		pMgr.refinedConceptValue.SetText("")
		pMgr.literalValueValue.Unbind()
		pMgr.literalValueValue.SetText("")
	} else {
		pMgr.typeValue.SetText(csb.conceptState.ConceptType)
		pMgr.idValue.SetText(csb.conceptState.ConceptID)
		pMgr.owningConceptIDValue.Bind(csb.owningConceptIDBinding)
		pMgr.versionValue.Bind(csb.versionBinding)
		pMgr.labelValue.Bind(csb.labelBinding)
		pMgr.definitionValue.Bind(csb.definitionBinding)
		pMgr.uriValue.Bind(csb.uriBinding)
		pMgr.isCoreValue.SetText(csb.conceptState.IsCore)
		pMgr.readOnlyValue.Bind(csb.readOnlyBinding)
		pMgr.referencedConceptValue.Bind(csb.referencedConceptIDBinding)
		pMgr.referencedConceptAttributeNameValue.Bind(csb.referencedAttributeNameBinding)
		pMgr.abstractConceptValue.Bind(csb.abstractConceptIDBinding)
		pMgr.refinedConceptValue.Bind(csb.refinedConceptIDBinding)
		pMgr.literalValueValue.Bind(csb.literalValueBinding)
	}
}

func (pMgr *FynePropertyManager) initialize() {
	pMgr.displayProperties("")
}

var _ fyne.Shortcutable = (*copyableLabel)(nil)
var _ fyne.Focusable = (*copyableLabel)(nil)
var _ fyne.Tappable = (*copyableLabel)(nil)

type copyableLabel struct {
	widget.Label
}

func newCopyableLabel() *copyableLabel {
	label := &copyableLabel{}
	label.ExtendBaseWidget(label)
	return label
}

func (cl *copyableLabel) FocusGained() {
}

func (cl *copyableLabel) FocusLost() {
}

func (cl *copyableLabel) Tapped(event *fyne.PointEvent) {
	FyneGUISingleton.GetWindow().RequestFocus()
}

// TappedSecondary is called when right or alternative tap is invoked.
// Implements: fyne.SecondaryTappable
func (cl *copyableLabel) TappedSecondary(pe *fyne.PointEvent) {
	clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	copyItem := fyne.NewMenuItem("Copy", func() {
		cl.TypedShortcut(&fyne.ShortcutCopy{Clipboard: clipboard})
	})

	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cl)
	popUpPos := entryPos.Add(fyne.NewPos(pe.Position.X, pe.Position.Y))
	c := fyne.CurrentApp().Driver().CanvasForObject(cl)

	menu := fyne.NewMenu("", copyItem)

	popUp := widget.NewPopUpMenu(menu, c)
	popUp.ShowAtPosition(popUpPos)
}

func (cl *copyableLabel) TypedKey(*fyne.KeyEvent) {

}

func (cl *copyableLabel) TypedRune(rune) {

}

func (cl *copyableLabel) TypedShortcut(shortcut fyne.Shortcut) {
	log.Print(shortcut.ShortcutName())
	switch typedShortcut := shortcut.(type) {
	case *fyne.ShortcutCopy:
		typedShortcut.Clipboard = FyneGUISingleton.window.Clipboard()
		typedShortcut.Clipboard.SetContent(cl.Text)
	}
}
