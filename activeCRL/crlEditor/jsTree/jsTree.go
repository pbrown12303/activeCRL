// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
