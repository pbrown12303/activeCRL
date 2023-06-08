package main

import (
	"flag"
	"log"

	"fyne.io/fyne/v2"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crleditorfynegui"
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
	crleditor.CrlEditorSingleton = crleditor.NewEditor(*userFolderArg)
	err := crleditor.CrlEditorSingleton.Initialize(*workspaceArg, true)
	if err != nil {
		log.Fatal(err)
	}
	fyneEditor := crleditorfynegui.NewFyneGUI(crleditor.CrlEditorSingleton)
	err = crleditor.CrlEditorSingleton.AddEditorGUI(fyneEditor)
	if err != nil {
		log.Fatal(err)
	}
	initialSize := fyne.NewSize(1600.0, 900.0)
	crleditorfynegui.FyneGUISingleton.GetWindow().Resize(initialSize)
	fyneEditor.GetWindow().ShowAndRun()
}
