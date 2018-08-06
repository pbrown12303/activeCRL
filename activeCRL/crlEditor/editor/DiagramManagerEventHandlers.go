package editor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/crlDiagram"
	"log"
)

func onDragover(event *js.Object, data *js.Object) {
	event.Call("preventDefault")
}

func onDiagramManagerDrop(event *js.Object) {
	hl := CrlEditorSingleton.getHeldLocks()
	event.Call("preventDefault")
	js.Global.Set("dropEvent", event)
	log.Printf("On Drop called")
	httpDiagramContainerId := event.Get("target").Get("parentElement").Get("parentElement").Get("id").String()
	be := CrlEditorSingleton.GetTreeDragSelection()
	js.Global.Set("diagramDroppedBaseElement", be.GetId(hl))
	js.Global.Get("console").Call("log", "In onDiagramManagerDrop")

	addNodeView(httpDiagramContainerId, be, event.Get("layerX").Float(), event.Get("layerY").Float(), hl)

	CrlEditorSingleton.SelectBaseElement(be)
	CrlEditorSingleton.SetTreeDragSelection(nil)
}

func onDiagramManagerCellPointerDown(cellView *js.Object, event *js.Object, x *js.Object, y *js.Object) {
	js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown")
	hl := CrlEditorSingleton.getHeldLocks()
	crlJointId := cellView.Get("model").Get("crlJointId").String()
	js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown crlJointId = "+crlJointId)
	diagramManager := CrlEditorSingleton.GetDiagramManager()
	diagramNode := diagramManager.jointElementIdToCrlDiagramNode[crlJointId]
	if diagramNode == nil {
		js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown diagramNode is nil")
	} else {
		js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown diagramNode id = "+diagramNode.GetId(hl))
	}

	be, err := crlDiagram.GetReferencedBaseElement(diagramNode, hl)
	if err == nil {
		CrlEditorSingleton.SelectBaseElement(be)
	}
}

func onMakeDiagramVisible(e jquery.Event) {
	httpDiagramContainerId := e.Get("target").Call("getAttribute", "viewId").String()
	//	js.Global.Get("console").Call("log", "In : onMakeDiagramVisible with: "+httpDiagramContainerId)
	//	js.Global.Set("clickEvent", e)
	//	js.Global.Set("clickEventTarget", e.Get("target"))
	//	js.Global.Set("clickEventViewId", e.Get("target").Call("getAttribute", "viewId"))
	makeDiagramVisible(httpDiagramContainerId)
}
