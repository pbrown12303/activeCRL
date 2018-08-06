package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
)

type PropertiesManager struct {
}

func NewPropertiesManager() *PropertiesManager {
	var newManager PropertiesManager
	return &newManager
}

func (pmPtr *PropertiesManager) BaseElementSelected(be core.BaseElement, hl *core.HeldLocks) {
	if be == nil {
		log.Printf("BaseElementSelected called with nil argument")
		return
	}
	properties := js.Global.Get("properties")

	// Type
	displayType(properties, be, 1, hl)

	// Id
	displayId(properties, be, 2, hl)

	// Version
	displayVersion(properties, be, 3, hl)
	// UniverseOfDiscourse
	displayUniverseOfDiscourse(properties, be, 4, hl)

	// Label
	displayLabel(properties, be, 5, hl)

	// URI
	displayUri(properties, be, 6, hl)

	switch be.(type) {
	case core.Element:
		// Definition
		displayDefinition(properties, be, 7, hl)
		clearRow(properties, 8)
	case core.Pointer:
		displayPointerProperties(properties, be, 7, hl)
	case core.Literal:
		clearRow(properties, 7)
		clearRow(properties, 8)
	}
}

func clearRow(properties *js.Object, row int) {
	propertyRow := properties.Get("rows").Index(row)
	if propertyRow != js.Undefined {
		properties.Call("deleteRow", row)
	}
}

func obtainPropertyRow(properties *js.Object, row int) *js.Object {
	propertyRow := properties.Get("rows").Index(row)
	if propertyRow == js.Undefined {
		propertyRow = properties.Call("insertRow", row)
		propertyRow.Call("insertCell", 0)
		propertyRow.Call("insertCell", 1)
	}
	return propertyRow
}

func displayDefinition(properties *js.Object, be core.BaseElement, row int, hl *core.HeldLocks) {
	definitionRow := obtainPropertyRow(properties, row)
	switch be.(type) {
	case core.Element:
		definitionRow.Get("cells").Index(0).Set("innerHTML", "Definition")
		definitionRow.Get("cells").Index(1).Set("innerHTML", be.(core.Element).GetDefinition(hl))
		definitionRow.Get("cells").Index(1).Set("id", "definition")
		uri := core.GetUri(be, hl)
		if uri != "" && core.BaseElement.GetId(be, hl) == uuid.NewV5(uuid.NamespaceURL, uri).String() {
			definitionRow.Get("cells").Index(1).Set("contentEditable", false)
		} else {
			definitionRow.Get("cells").Index(1).Set("contentEditable", true)
			definitionQuery := jquery.NewJQuery("#definition")
			definitionQuery.On(jquery.KEYUP, func(e jquery.Event) {
				definition := jquery.NewJQuery(e.Target).Text()
				CrlEditorSingleton.SetSelectionDefinition(definition)
			})
		}
	default:
		if definitionRow != js.Undefined {
			properties.Call("deleteRow", row)
		}
	}
}

func displayId(properties *js.Object, be core.BaseElement, row int, hl *core.HeldLocks) {
	if be == nil {
		log.Printf("PropertiesManager displayId called with nil base element")
		return
	}
	idRow := obtainPropertyRow(properties, row)
	idRow.Get("cells").Index(0).Set("innerHTML", "Id")
	idRow.Get("cells").Index(1).Set("innerHTML", core.BaseElement.GetId(be, hl))
}

func displayLabel(properties *js.Object, be core.BaseElement, row int, hl *core.HeldLocks) {
	labelRow := obtainPropertyRow(properties, row)
	labelRow.Get("cells").Index(0).Set("innerHTML", "Label")
	labelRow.Get("cells").Index(1).Set("innerHTML", core.GetLabel(be, hl))
	labelRow.Get("cells").Index(1).Set("id", "baseElementLabel")
	switch be.(type) {
	case core.Element:
		uri := core.GetUri(be, hl)
		if uri != "" && core.BaseElement.GetId(be, hl) == uuid.NewV5(uuid.NamespaceURL, uri).String() {
			labelRow.Get("cells").Index(1).Set("contentEditable", false)
		} else {
			labelRow.Get("cells").Index(1).Set("contentEditable", true)
			nameQuery := jquery.NewJQuery("#baseElementLabel")
			nameQuery.On(jquery.KEYUP, func(e jquery.Event) {
				name := jquery.NewJQuery(e.Target).Text()
				CrlEditorSingleton.SetSelectionLabel(name)
			})
		}
	default:
		labelRow.Get("cells").Index(1).Set("contentEditable", false)
	}
}

func displayPointerProperties(properties *js.Object, be core.BaseElement, row int, hl *core.HeldLocks) {
	indicatedBaseElementRow := obtainPropertyRow(properties, row)
	indicatedBaseElementVersionRow := obtainPropertyRow(properties, row+1)
	switch be.(type) {
	case core.BaseElementPointer:
		indicatedBaseElementRow.Get("cells").Index(0).Set("innerHTML", "Indicated BaseElement Id")
		indicatedBaseElementRow.Get("cells").Index(1).Set("innerHTML", be.(core.BaseElementPointer).GetBaseElementId(hl))
		indicatedBaseElementVersionRow.Get("cells").Index(0).Set("innerHTML", "Indicated BaseElement Version")
		indicatedBaseElementVersionRow.Get("cells").Index(1).Set("innerHTML", strconv.Itoa(be.(core.BaseElementPointer).GetBaseElementVersion(hl)))
	case core.ElementPointer:
		indicatedBaseElementRow.Get("cells").Index(0).Set("innerHTML", "Indicated Element Id")
		indicatedBaseElementRow.Get("cells").Index(1).Set("innerHTML", be.(core.ElementPointer).GetElementId(hl))
		indicatedBaseElementVersionRow.Get("cells").Index(0).Set("innerHTML", "Indicated Element Version")
		indicatedBaseElementVersionRow.Get("cells").Index(1).Set("innerHTML", strconv.Itoa(be.(core.ElementPointer).GetElementVersion(hl)))
	case core.ElementPointerPointer:
		indicatedBaseElementRow.Get("cells").Index(0).Set("innerHTML", "Indicated ElementPointer Id")
		indicatedBaseElementRow.Get("cells").Index(1).Set("innerHTML", be.(core.ElementPointerPointer).GetElementPointerId(hl))
		indicatedBaseElementVersionRow.Get("cells").Index(0).Set("innerHTML", "Indicated ElementPointer Version")
		indicatedBaseElementVersionRow.Get("cells").Index(1).Set("innerHTML", strconv.Itoa(be.(core.ElementPointerPointer).GetElementPointerVersion(hl)))
	case core.LiteralPointer:
		indicatedBaseElementRow.Get("cells").Index(0).Set("innerHTML", "Indicated Literal Id")
		indicatedBaseElementRow.Get("cells").Index(1).Set("innerHTML", be.(core.LiteralPointer).GetLiteralId(hl))
		indicatedBaseElementVersionRow.Get("cells").Index(0).Set("innerHTML", "Indicated Literal Version")
		indicatedBaseElementVersionRow.Get("cells").Index(1).Set("innerHTML", strconv.Itoa(be.(core.LiteralPointer).GetLiteralVersion(hl)))
	case core.LiteralPointerPointer:
		indicatedBaseElementRow.Get("cells").Index(0).Set("innerHTML", "Indicated LiteralPointer Id")
		indicatedBaseElementRow.Get("cells").Index(1).Set("innerHTML", be.(core.LiteralPointerPointer).GetLiteralPointerId(hl))
		indicatedBaseElementVersionRow.Get("cells").Index(0).Set("innerHTML", "Indicated LiteralPointer Version")
		indicatedBaseElementVersionRow.Get("cells").Index(1).Set("innerHTML", strconv.Itoa(be.(core.LiteralPointerPointer).GetLiteralPointerVersion(hl)))
	}
}

func displayType(properties *js.Object, be core.BaseElement, row int, hl *core.HeldLocks) {
	typeRow := obtainPropertyRow(properties, row)
	typeRow.Get("cells").Index(0).Set("innerHTML", "Type")
	typeRow.Get("cells").Index(1).Set("innerHTML", core.GetTypeName(be))
}

func displayUniverseOfDiscourse(properties *js.Object, be core.BaseElement, row int, hl *core.HeldLocks) {
	uOfDRow := obtainPropertyRow(properties, row)
	uOfDRow.Get("cells").Index(0).Set("innerHTML", "UofD Id")
	uOfDRow.Get("cells").Index(1).Set("innerHTML", be.GetUniverseOfDiscourse(hl).GetId(hl))
}

func displayUri(properties *js.Object, be core.BaseElement, row int, hl *core.HeldLocks) {
	uriRow := obtainPropertyRow(properties, row)
	uri := core.GetUri(be, hl)
	uriRow.Get("cells").Index(0).Set("innerHTML", "URI")
	uriRow.Get("cells").Index(1).Set("innerHTML", uri)
	uriRow.Get("cells").Index(1).Set("id", "uri")
	if uri != "" && core.BaseElement.GetId(be, hl) == uuid.NewV5(uuid.NamespaceURL, uri).String() {
		uriRow.Get("cells").Index(1).Set("contentEditable", false)
	} else {
		uriRow.Get("cells").Index(1).Set("contentEditable", true)
		uriQuery := jquery.NewJQuery("#uri")
		uriQuery.On(jquery.KEYUP, func(e jquery.Event) {
			uri := jquery.NewJQuery(e.Target).Text()
			CrlEditorSingleton.SetSelectionUri(uri)
		})
	}
}

func displayVersion(properties *js.Object, be core.BaseElement, row int, hl *core.HeldLocks) {
	versionRow := obtainPropertyRow(properties, row)
	versionRow.Get("cells").Index(0).Set("innerHTML", "Version")
	versionRow.Get("cells").Index(1).Set("innerHTML", strconv.Itoa(core.BaseElement.GetVersion(be, hl)))
}
