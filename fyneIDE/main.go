package main

import (
	"flag"
	"log"

	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crleditorfynegui/fynegui"
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
	crleditor := crleditor.NewEditor(*userFolderArg)
	fyneEditor := fynegui.NewFyneGUI(crleditor)
	err := crleditor.AddEditorGUI(fyneEditor)
	if err != nil {
		log.Fatal(err)
	}
	err = crleditor.Initialize(*workspaceArg, true)
	if err != nil {
		log.Fatal(err)
	}

	fyneEditor.GetWindow().ShowAndRun()
}
