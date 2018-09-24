package editor

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/pbrown12303/activeCRL/activeCRL/crlDiagram"
)

func onDragover(event *js.Object, data *js.Object) {
	event.Call("preventDefault")
}

func onDiagramManagerDrop(event *js.Object) {
	hl := CrlEditorSingleton.getHeldLocks()
	defer hl.ReleaseLocksAndWait()
	event.Call("preventDefault")
	js.Global.Set("dropEvent", event)
	log.Printf("On Drop called")
	httpDiagramContainerID := event.Get("target").Get("parentElement").Get("parentElement").Get("id").String()
	be := CrlEditorSingleton.GetTreeDragSelection()
	js.Global.Set("diagramDroppedBaseElement", be.GetId(hl))
	js.Global.Get("console").Call("log", "In onDiagramManagerDrop")

	addNodeView(httpDiagramContainerID, be, event.Get("layerX").Float(), event.Get("layerY").Float(), hl)

	CrlEditorSingleton.SelectBaseElement(be)
	CrlEditorSingleton.SetTreeDragSelection(nil)
}

func onDiagramManagerCellPointerDown(cellView *js.Object, event *js.Object, x *js.Object, y *js.Object) {
	js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown")
	hl := CrlEditorSingleton.getHeldLocks()
	defer hl.ReleaseLocksAndWait()
	crlJointID := cellView.Get("model").Get("crlJointId").String()
	js.Global.Get("console").Call("log", "In onDiagramManagerCellPointerDown crlJointId = "+crlJointID)
	diagramManager := CrlEditorSingleton.GetDiagramManager()
	diagramNode := diagramManager.jointElementIDToCrlDiagramNode[crlJointID]
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
	httpDiagramContainerID := e.Get("target").Call("getAttribute", "viewId").String()
	makeDiagramVisible(httpDiagramContainerID)
}
