// Copyright 2017, 2018 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/pbrown12303/activeCRL/crleditor/editor"
	"github.com/pbrown12303/activeCRL/crlfyneeditor/fyneeditor"
)

func main() {
	workspaceArg := flag.String("workspace", "", "Path to workspace folder (optional)")
	userFolderArg := flag.String("userFolder", "", "Path to user folder (optional)")
	flag.Parse()
	fmt.Println("workspace: ", *workspaceArg)
	fmt.Println("user folder: ", *userFolderArg)
	editor.CrlLogClientRequests = true
	go editor.StartServer(true, *workspaceArg, *userFolderArg)

	// Fyne
	fyneEditor := fyneeditor.NewCrlFyneEditor()
	fyneEditor.GetWindow().ShowAndRun()
}
