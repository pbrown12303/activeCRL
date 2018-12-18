package editor

import (
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagram"
	//	"log"
)

func addDiagramViewFunctionsToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {
	uOfD.AddFunction(crldiagram.CrlDiagramURI, updateDiagramView)
}

func updateDiagramView(el core.Element, changeNotifications *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	go updateDiagramViewInternal(el, changeNotifications, uOfD)
}

func updateDiagramViewInternal(el core.Element, changeNotifications *core.ChangeNotification, uOfD core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocks()
	// changeNotification.Print("updateDiagramView ", hl)
	// Check whether the name has changed
	// Check to see whether a diagram node or edge has been added or removed from the diagram
}

func init() {
	//	core.GetCore().AddFunction(crlDiagram.CrlDiagramUri, updateDiagramView)
}
