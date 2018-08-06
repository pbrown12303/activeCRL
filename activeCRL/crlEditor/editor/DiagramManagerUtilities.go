package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"strconv"
)

var defaultLabelCount int
var diagramTabCount int
var diagramContainerCount int
var jointGraphCount int
var jointBaseElementNodeCount int

func createDiagramTabPrefix() string {
	diagramTabCount++
	countString := strconv.Itoa(diagramTabCount)
	return "DiagramTab" + countString
}

func createDiagramContainerPrefix() string {
	diagramContainerCount++
	countString := strconv.Itoa(diagramContainerCount)
	return "DiagramView" + countString
}

func createJointBaseElementNodePrefix() string {
	jointBaseElementNodeCount++
	countString := strconv.Itoa(jointGraphCount)
	return "JointBaseElementNode" + countString
}

func createJointGraphPrefix() string {
	jointGraphCount++
	countString := strconv.Itoa(jointGraphCount)
	return "DiagramGraph" + countString
}

func getDefaultDiagramLabel() string {
	defaultLabelCount++
	countString := strconv.Itoa(defaultLabelCount)
	return "Diagram" + countString
}

func makeDiagramVisible(httpDiagramContainerId string) {
	x := js.Global.Get("document").Call("getElementsByClassName", "crlDiagramContainer")
	lengthString := strconv.Itoa(x.Length())
	js.Global.Get("console").Call("log", "List length: "+lengthString)
	for i := 0; i < x.Length(); i++ {
		js.Global.Get("console").Call("log", "Container id: ", x.Index(i).Get("id").String())
		if x.Index(i).Get("id").String() == httpDiagramContainerId {
			x.Index(i).Get("style").Set("display", "block")
			js.Global.Get("console").Call("log", "Showing: "+httpDiagramContainerId)
		} else {
			x.Index(i).Get("style").Set("display", "none")
			js.Global.Get("console").Call("log", "Hiding: "+httpDiagramContainerId)
		}
	}

}
