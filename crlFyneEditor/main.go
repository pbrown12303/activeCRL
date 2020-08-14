package main

import (
	"github.com/pbrown12303/activeCRL/crlfyneeditor/fyneeditor"
)

func main() {
	fyneEditor := fyneeditor.NewCrlFyneEditor()
	fyneEditor.GetWindow().ShowAndRun()
}
