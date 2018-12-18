package editor

import (
	"log"
	"reflect"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/satori/go.uuid"
)

// PropertiesManager is the manager of the properties display in the client
type PropertiesManager struct {
}

// NewPropertiesManager creates an instance of the PropertiesManager
func NewPropertiesManager() *PropertiesManager {
	var newManager PropertiesManager
	return &newManager
}

// ElementSelected updates the client's property display after an element selection change
func (pmPtr *PropertiesManager) ElementSelected(el core.Element, hl *core.HeldLocks) {
	if el == nil {
		log.Printf("BaseElementSelected called with nil argument")
		return
	}
	properties := js.Global.Get("properties")

	// Type
	displayType(properties, el, 1, hl)

	// Id
	displayID(properties, el, 2, hl)

	// Version
	displayVersion(properties, el, 3, hl)
	// UniverseOfDiscourse
	displayUniverseOfDiscourse(properties, el, 4, hl)

	// Label
	displayLabel(properties, el, 5, hl)

	// URI
	displayURI(properties, el, 6, hl)
	// TODO: Fix URI display
	switch el.(type) {
	case core.Element:
		// Definition
		displayDefinition(properties, el, 7, hl)
		clearRow(properties, 8)
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

func displayDefinition(properties *js.Object, el core.Element, row int, hl *core.HeldLocks) {
	definitionRow := obtainPropertyRow(properties, row)
	switch el.(type) {
	case core.Element:
		definitionRow.Get("cells").Index(0).Set("innerHTML", "Definition")
		definitionRow.Get("cells").Index(1).Set("innerHTML", el.(core.Element).GetDefinition(hl))
		definitionRow.Get("cells").Index(1).Set("id", "definition")
		uri := el.GetURI(hl)
		if uri != "" && el.GetConceptID(hl) == uuid.NewV5(uuid.NamespaceURL, uri).String() {
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

func displayID(properties *js.Object, el core.Element, row int, hl *core.HeldLocks) {
	if el == nil {
		log.Printf("PropertiesManager displayID called with nil base element")
		return
	}
	idRow := obtainPropertyRow(properties, row)
	idRow.Get("cells").Index(0).Set("innerHTML", "Id")
	idRow.Get("cells").Index(1).Set("innerHTML", core.Element.GetConceptID(el, hl))
}

func displayLabel(properties *js.Object, el core.Element, row int, hl *core.HeldLocks) {
	labelRow := obtainPropertyRow(properties, row)
	labelRow.Get("cells").Index(0).Set("innerHTML", "Label")
	labelRow.Get("cells").Index(1).Set("innerHTML", el.GetLabel(hl))
	labelRow.Get("cells").Index(1).Set("id", "baseElementLabel")
	switch el.(type) {
	case core.Element:
		uri := el.GetURI(hl)
		if uri != "" && core.Element.GetConceptID(el, hl) == uuid.NewV5(uuid.NamespaceURL, uri).String() {
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

func displayType(properties *js.Object, el core.Element, row int, hl *core.HeldLocks) {
	typeRow := obtainPropertyRow(properties, row)
	typeRow.Get("cells").Index(0).Set("innerHTML", "Type")
	typeRow.Get("cells").Index(1).Set("innerHTML", reflect.TypeOf(el).String())
}

func displayUniverseOfDiscourse(properties *js.Object, el core.Element, row int, hl *core.HeldLocks) {
	uOfDRow := obtainPropertyRow(properties, row)
	uOfDRow.Get("cells").Index(0).Set("innerHTML", "UofD Id")
	uOfDRow.Get("cells").Index(1).Set("innerHTML", el.GetUniverseOfDiscourse(hl).GetConceptID(hl))
}

func displayURI(properties *js.Object, el core.Element, row int, hl *core.HeldLocks) {
	uriRow := obtainPropertyRow(properties, row)
	uri := el.GetURI(hl)
	uriRow.Get("cells").Index(0).Set("innerHTML", "URI")
	uriRow.Get("cells").Index(1).Set("innerHTML", uri)
	uriRow.Get("cells").Index(1).Set("id", "uri")
	if uri != "" && core.Element.GetConceptID(el, hl) == uuid.NewV5(uuid.NamespaceURL, uri).String() {
		uriRow.Get("cells").Index(1).Set("contentEditable", false)
	} else {
		uriRow.Get("cells").Index(1).Set("contentEditable", true)
		uriQuery := jquery.NewJQuery("#uri")
		uriQuery.On(jquery.KEYUP, func(e jquery.Event) {
			uri := jquery.NewJQuery(e.Target).Text()
			CrlEditorSingleton.SetSelectionURI(uri)
		})
	}
}

func displayVersion(properties *js.Object, el core.Element, row int, hl *core.HeldLocks) {
	versionRow := obtainPropertyRow(properties, row)
	versionRow.Get("cells").Index(0).Set("innerHTML", "Version")
	versionRow.Get("cells").Index(1).Set("innerHTML", strconv.Itoa(core.Element.GetVersion(el, hl)))
}
