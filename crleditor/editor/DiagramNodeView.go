package editor

import (
	"github.com/golang/freetype/truetype"
	"github.com/gopherjs/gopherjs/js"
	"github.com/pbrown12303/activeCRL/core"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

var goRegularFont *truetype.Font
var goBoldFont *truetype.Font
var go12PtRegularFace font.Face
var go12PtBoldFace font.Face

func calculateTextWidth(text string) int {
	size := font.MeasureString(go12PtBoldFace, text)
	return size.Ceil()
}

func defineNodeViews() {
	// Define the Element Graph Node and View
	crlBeDefaultInstanceProps := NewBeDefaultInstanceProperties()
	crlBePrototypeProps := NewBePrototypeProperties()
	js.Global.Get("joint").Get("dia").Get("Element").Call("define", "crl.Element", crlBeDefaultInstanceProps, crlBePrototypeProps)
	elementViewExtension := js.Global.Get("joint").Get("dia").Get("ElementView").Call("extend", js.M{}, js.M{
		"initialize": js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			js.Global.Get("joint").Get("dia").Get("ElementView").Get("prototype").Get("initialize").Call("apply", this, arguments)
			this.Call("listenTo", this.Get("model"), "crl-update", js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
				this.Call("update")
				this.Call("resize")
				return nil
			}))
			return nil
		})})
	js.Global.Get("joint").Get("shapes").Get("crl").Set("BaseElementView", elementViewExtension)
}

func updateDiagramNodeView(el core.Element, changeNotifications *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	//	for _, changeNotification :

}

func init() {
	var err error

	// Set up fonts and faces
	goRegularFont, err = truetype.Parse(goregular.TTF)
	if err != nil {
		js.Global.Get("console").Call("log", err)
	}
	goBoldFont, err = truetype.Parse(gobold.TTF)
	if err != nil {
		js.Global.Get("console").Call("log", err)
	}
	go12PtRegularFace = truetype.NewFace(goRegularFont, nil)
	go12PtBoldFace = truetype.NewFace(goBoldFont, nil)
}
