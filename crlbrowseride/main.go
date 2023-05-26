// Copyright 2017, 2018 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"time"

	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crleditorbrowsergui/crleditorbrowsergui"
	// "github.com/pbrown12303/activeCRL/crleditorfynegui/fynegui"
)

func main() {
	workspaceArg := flag.String("workspace", "", "Path to workspace folder (optional)")
	userFolderArg := flag.String("userFolder", "", "Path to user folder (optional)")
	flag.Parse()
	log.Println("workspace: ", *workspaceArg)
	log.Println("user folder: ", *userFolderArg)
	// For debugging
	// browsergui.CrlLogClientRequests = true

	// Common infrastructure
	editor := crleditor.NewEditor(*userFolderArg)
	crleditorbrowsergui.InitializeBrowserGUISingleton(editor, true)
	err := editor.AddEditorGUI(crleditorbrowsergui.BrowserGUISingleton)
	if err != nil {
		log.Fatal(err)
	}
	// fyneGUI := fynegui.NewFyneGUI()
	// err = editor.AddEditorGUI(fyneGUI)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = editor.Initialize(*workspaceArg, true)
	if err != nil {
		log.Fatal(err)
	}

	// Fyne GUI
	// fyneGUI.GetWindow().ShowAndRun()

	// Alternate top-level when not using Fyne
	for editor.GetExitRequested() == false {
		time.Sleep(1 * time.Second)
	}
}
