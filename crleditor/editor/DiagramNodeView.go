package editor

import (
	"log"

	"github.com/pbrown12303/activeCRL/crldiagram"

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

func addNodeView(httpDiagramContainerID string, be core.Element, x float64, y float64, hl *core.HeldLocks) (core.Element, error) {
	//	create crlDiagramNode
	uOfD := CrlEditorSingleton.GetUofD()
	diagramManager := CrlEditorSingleton.GetDiagramManager()
	crlDiag := diagramManager.diagramContainerIDToCrlDiagram[httpDiagramContainerID]

	// Tracing
	if core.AdHocTrace == true {
		log.Printf("In addNodeView CrlDiagramId is " + crlDiag.GetConceptID(hl))
	}

	newCrlDiagramNode, err := uOfD.CreateReplicateAsRefinementFromURI(crldiagram.CrlDiagramNodeURI, hl)
	if err != nil {
		js.Global.Get("console").Call("log", "Failed to create CrlDiagramNode"+err.Error())
		return nil, err
	}
	crldiagram.SetReferencedElement(newCrlDiagramNode, be, hl)

	// Tracing
	if core.AdHocTrace == true {
		log.Printf("In addNodeView about to call SetOwningElement on new diagram node")
	}
	newCrlDiagramNode.SetOwningConcept(crlDiag, hl)

	// Tracing
	if core.AdHocTrace == true {
		log.Printf("In addNodeView CrlDiagramNodeId is " + newCrlDiagramNode.GetConceptID(hl))
	}

	// Now construct the jointjs representation
	graph := diagramManager.crlDiagramIDToJointGraph[httpDiagramContainerID]
	jointBaseElementID := createJointBaseElementNodePrefix() + newCrlDiagramNode.GetConceptID(hl)
	jointBaseElement := js.Global.Get("joint").Get("shapes").Get("crl").Get("Element").New(NewBeDefaultInstanceProperties(), NewBePrototypeProperties())
	jointBaseElement.Set("crlJointId", jointBaseElementID)
	js.Global.Set("jointBaseElement", jointBaseElement)
	// name
	name := be.GetLabel(hl)
	jointBaseElement.Get("attributes").Set("name", name)
	// position
	jointBaseElement.Get("attributes").Set("position", js.M{"x": x, "y": y})
	// image
	jointBaseElement.Get("attributes").Get("attrs").Set("image", js.M{"xlink:href": "/icons/ElementIcon.svg"})

	diagramManager.jointElementIDToCrlDiagramNode[jointBaseElementID] = newCrlDiagramNode

	jointBaseElement.Call("updateRectangles")
	js.Global.Set("graph", graph)
	graph.Call("addCell", jointBaseElement)

	return newCrlDiagramNode, nil
}

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
