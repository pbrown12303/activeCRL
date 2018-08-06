package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"strconv"
)

func NewBeDefaultInstanceProperties() js.M {
	crlBeDefaultInstanceProps := js.M{
		"attrs": js.M{
			"rect": js.M{
				"width": 300},
			".image": js.M{
				"ref-x":  1.0,
				"ref-y":  1.0,
				"ref":    ".label-rect",
				"width":  16,
				"height": 16},
			".label-rect": js.M{
				"stroke":       "black",
				"stroke-width": 2,
				"fill":         "#ffffff",
				"height":       40,
				"transform":    "translate(0,0)"},
			".abstractions-text": js.M{
				"ref":         ".label-rect",
				"ref-y":       0.5,
				"ref-x":       0.5 + 18,
				"text-anchor": "right",
				"y-alignment": "middle",
				"font-weight": "normal",
				"font-style":  "italic",
				"fill":        "black",
				"font-size":   12,
				"font-family": "Go,  Helvetica, Ariel, sans-serif"},
			".label-text": js.M{
				"ref":         ".label-rect",
				"ref-y":       0.5,
				"ref-x":       0.5 + 18,
				"text-anchor": "left",
				"y-alignment": "middle",
				"font-weight": "bold",
				"fill":        "black",
				"font-size":   12,
				"font-family": "Go,  Helvetica, Ariel, sans-serif"}}}

	return crlBeDefaultInstanceProps
}

func NewBePrototypeProperties() js.M {
	// Create the prototype properties
	crlBePrototypeProps := js.M{
		"markup": "<g class=\"rotatable\">" +
			"<g class=\"scalable\">" +
			"<rect class=\"label-rect\"/>" +
			"</g>" +
			"<image class=\"image\"/><text class=\"abstractions-text\"/><text class=\"label-text\"/>" +
			"</g>",
		"initialize": js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			js.Global.Get("joint").Get("shapes").Get("basic").Get("Generic").Get("prototype").Get("initialize").Call("apply", this, arguments)
			this.Call("updateRectangles")
			return nil
		}),
		"updateRectangles": js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			offsetY := 0
			attributes := this.Get("attributes")
			attrs := attributes.Get("attrs")

			rectHeight := 1*12 + 6
			labelText := attributes.Get("name")
			labelTextAttr := attrs.Get(".label-text")
			labelTextAttr.Set("text", labelText)
			labelRectAttr := attrs.Get(".label-rect")
			labelRectAttr.Set("height", rectHeight)
			rectWidth := calculateTextWidth(labelText.String()) + 6 + 18
			labelRectAttr.Set("transform", "translate(0,"+strconv.Itoa(offsetY)+")")
			this.Call("resize", rectWidth, rectHeight)

			offsetY += rectHeight
			return nil
		})}
	return crlBePrototypeProps
}
