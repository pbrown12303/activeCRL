package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"strconv"
)

type paperProperties struct {
	*js.Object
	el       []*js.Object `js:"el"`
	width    float64      `js:"width"`
	height   float64      `js:"height"`
	model    *js.Object   `js:"model"`
	gridSize float64      `js:"gridSize"`
}

type positionProperties struct {
	*js.Object
	x float64 `js:"x"`
	y float64 `js:"y"`
}

type sizeProperties struct {
	*js.Object
	width  float64 `js:"width"`
	height float64 `js:"height"`
}

type shapeProperties struct {
	*js.Object
	fill        string `js:"fill"`
	stroke      string `js:"stroke"`
	strokeWidth string `js:"stroke-width"`
}

type textProperties struct {
	*js.Object
	text string `js:"text"`
	fill string `js:"fill"`
}

type rectProperties struct {
	*js.Object
	position *positionProperties `js:"position"`
	size     *sizeProperties     `js:"size"`
}

type beProperties struct {
	*js.Object
	position *positionProperties `js:"position"`
	size     *sizeProperties     `js:"size"`
	attrs    js.M                `js:"attrs"`
	name     string              `js:"name"`
}

type beLabelProperty struct {
	*js.Object
	name string `js:"label"`
}

type beDotLabelRectProperties struct {
	*js.Object
	stroke      string  `js:"stroke"`
	strokeWidth float64 `js:"stroke-width"`
	fill        string  `js:"fill"`
}

type beDotLabelTextProperties struct {
	*js.Object
	ref        string  `js:"ref"`
	refY       float64 `js:"ref-y"`
	refX       float64 `js:"ref-x"`
	textAnchor string  `js:"text-anchor"`
	yAlignment string  `js:"y-alignment"`
	fontWeight string  `js:"font-weight"`
	fill       string  `js:"fill"`
	fontSize   float64 `js:"font-size"`
	fontFamily string  `js:"font-family"`
}

type beDefaultInstanceProperties struct {
	*js.Object
	attrs        *beDefaultInstanceAttrsProperties `js:"attrs"`
	crlJointId   string                            `js:"crlJointId"`
	label        string                            `js: "label"`
	abstractions []string                          `js: "abstractions"`
}

type beDefaultInstanceAttrsProperties struct {
	*js.Object
	rect         *rectProperties           `js:"rect"`
	dotLabelRect *beDotLabelRectProperties `js:"label-rect"`
	dotLabelText *beDotLabelTextProperties `js:".label-text"`
}

type beAbstractionsProperties struct {
	*js.Object
}

type bePrototypeProperties struct {
	*js.Object
	markup           string     `js:"markup"`
	initialize       *js.Object `js:"initialize"`
	updateRectangles *js.Object `js:"updateRectangles"`
}

type beInitializeProperties struct {
	*js.Object
}

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
