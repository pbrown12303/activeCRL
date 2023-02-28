package main

//go:generate fyne bundle --package images -o ../images/icons.go --prefix Resource ../crleditorbrowsergui/http/images/icons/AbstractPointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/CursorIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/DiagramIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/ElementIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/ElementPointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/LiteralIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/OwnerPointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/PointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/ReferenceIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/ReferenceLinkIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/RefinedPointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/RefinementIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../crleditorbrowsergui/http/images/icons/RefinementLinkIcon.png

import (
	"flag"
	"log"

	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crleditorbrowsergui/browsergui"
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
	crleditor.CrlEditorSingleton = crleditor.NewEditor(*userFolderArg)

	browsergui.InitializeBrowserGUISingleton(crleditor.CrlEditorSingleton, true)
	err := crleditor.CrlEditorSingleton.AddEditorGUI(browsergui.BrowserGUISingleton)
	if err != nil {
		log.Fatal(err)
	}

	fyneEditor := fynegui.NewFyneGUI(crleditor.CrlEditorSingleton)
	err = crleditor.CrlEditorSingleton.AddEditorGUI(fyneEditor)
	if err != nil {
		log.Fatal(err)
	}
	err = crleditor.CrlEditorSingleton.Initialize(*workspaceArg, true)
	if err != nil {
		log.Fatal(err)
	}

	fyneEditor.GetWindow().ShowAndRun()
}
