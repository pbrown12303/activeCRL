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

	"fyne.io/fyne/v2"
	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crleditorbrowsergui/crleditorbrowsergui"
	fynegui "github.com/pbrown12303/activeCRL/crleditorfynegui"
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

	crleditorbrowsergui.InitializeBrowserGUISingleton(crleditor.CrlEditorSingleton, true)
	err := crleditor.CrlEditorSingleton.AddEditorGUI(crleditorbrowsergui.BrowserGUISingleton)
	if err != nil {
		log.Fatal(err)
	}

	fynegui.NewFyneGUI(crleditor.CrlEditorSingleton)
	err = crleditor.CrlEditorSingleton.AddEditorGUI(fynegui.FyneGUISingleton)
	if err != nil {
		log.Fatal(err)
	}
	err = crleditor.CrlEditorSingleton.Initialize(*workspaceArg, true)
	if err != nil {
		log.Fatal(err)
	}
	initialSize := fyne.NewSize(1600.0, 900.0)
	fynegui.FyneGUISingleton.GetWindow().Resize(initialSize)
	fynegui.FyneGUISingleton.GetWindow().ShowAndRun()
}