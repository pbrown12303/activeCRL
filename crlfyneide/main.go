package main

//go:generate fyne bundle --package images -o ../images/icons.go --prefix Resource ../crleditorbrowsergui/http/images/icons/AbstractPointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/CursorIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/DiagramIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/ElementIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/ElementPointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/LiteralIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/OwnerPointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/PointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/ReferenceIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/ReferenceLinkIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/RefinedPointerIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/RefinementIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/RefinementLinkIcon.png
//go:generate fyne bundle -a -o ../images/icons.go  --prefix Resource ../images/icons/OneToOneIcon.png

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

	// Common infrastructure
	crleditor.CrlEditorSingleton = crleditor.NewEditor(*userFolderArg)
	err := crleditor.CrlEditorSingleton.Initialize(*workspaceArg, true)
	if err != nil {
		log.Fatal(err)
	}

	// Fyne GUI
	fyneEditor := crleditorfynegui.NewFyneGUI(crleditor.CrlEditorSingleton)
	err = crleditor.CrlEditorSingleton.AddEditorGUI(fyneEditor)
	if err != nil {
		log.Fatal(err)
	}
	initialSize := fyne.NewSize(1600.0, 900.0)
	crleditorfynegui.FyneGUISingleton.GetWindow().Resize(initialSize)
	fyneEditor.GetWindow().ShowAndRun()
}
