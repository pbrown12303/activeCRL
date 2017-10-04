package jsTree

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	JST = "jstree"
)

type JsTree struct {
	o *js.Object
}

//JQuery constructor
func NewJsTree(args ...interface{}) JsTree {
	return JsTree{o: js.Global.Get(JST).New(args...)}
}
