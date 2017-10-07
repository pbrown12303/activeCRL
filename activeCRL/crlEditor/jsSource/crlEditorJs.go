// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//go:generate gopherjs build crlEditorJs.go -o ../js/crlEditor.js -m
// +build ignore

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	//	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/crlEditor/jsTree"
	"log"
)

//convenience:
var jQuery = jquery.NewJQuery

const (
	INPUT  = "input#name"
	OUTPUT = "span#output"
	TREE   = "div#uOfD"
)

func main() {

	//	uOfD := core.NewUniverseOfDiscourse()

	// *** EXPERIMENTAL CODE BELOW THIS LINE **************************************************************

	//show jQuery Version on console:
	print("Your current jQuery version is: " + jQuery().Jquery)

	//catch keyup events on input#name element:
	jQuery(INPUT).On(jquery.KEYUP, func(e jquery.Event) {

		name := jQuery(e.Target).Val()
		name = jquery.Trim(name)

		//show welcome message:
		if len(name) > 0 {
			jQuery(OUTPUT).SetText("Welcome to GopherJS, " + name + " !")
		} else {
			jQuery(OUTPUT).Empty()
		}
	})

	js.Global.Set("newJQuery", newJQuery)

	js.Global.Set("insertNode", insertNode)

	js.Global.Set("logEntry", logEntry)

	js.Global.Set("NewJsTree", jsTree.NewJsTree)

	js.Global.Set("selectChildNode1", selectChildNode1)

	js.Global.Set("addCallback", addCallback)

	js.Global.Set("pet", map[string]interface{}{
		"NewPet": NewPet,
	})

}

// *** EXPERIMENTAL CODE BELOW THIS LINE **************************************************************

func newJQuery(args ...interface{}) jquery.JQuery {
	return jquery.NewJQuery(args...)
}

func logEntry(x string) {
	log.Printf("Log Entry %s\n", x)
}

func selectChildNode1() {
	jquery.NewJQuery("#uOfD").Call("jstree", "select_node", "child_node_1")
}

func getJsTree() {

}

func insertNode() jquery.JQuery {
	return jquery.NewJQuery("#uOfD").Call("jstree", "create_node", "child_node_1", "New Node", "last")
}

func addCallback() {
	//catch keyup events on input#name element:
	jQuery(INPUT).On(jquery.KEYUP, func(e jquery.Event) {

		name := jQuery(e.Target).Val()
		name = jquery.Trim(name)

		//show welcome message:
		if len(name) > 0 {
			jQuery(OUTPUT).SetText("Welcome to GopherJS, " + name + " !")
		} else {
			jQuery(OUTPUT).Empty()
		}
	})
}

type Pet struct {
	name string
}

func NewPet(name string) *js.Object {
	log.Printf("Pet name: %s\n", name)
	return js.MakeWrapper(&Pet{name})
}

func (p *Pet) Name() string {
	return p.name
}

func (p *Pet) SetName(name string) {
	p.name = name
}
