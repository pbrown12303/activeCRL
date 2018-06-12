package editor

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/coreDiagram"
	//	"log"
	"sync"
)

func addDiagramViewFunctionsToUofD(uOfD core.UniverseOfDiscourse, hl *core.HeldLocks) {
	core.GetCore().AddFunction(coreDiagram.CrlDiagramUri, updateDiagramView)
}

func updateDiagramView(el core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	go updateDiagramViewInternal(el, changeNotifications, wg)
}

func updateDiagramViewInternal(el core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	for _, changeNotification := range changeNotifications {
		changeNotification.Print("updateDiagramView ", hl)
		// Check whether the name has changed
		// Check to see whether a diagram node or edge has been added or removed from the diagram
	}
}

func init() {
	//	core.GetCore().AddFunction(coreDiagram.CrlDiagramUri, updateDiagramView)
}
