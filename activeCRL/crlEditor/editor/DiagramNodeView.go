package editor

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/coreDiagram"
	"sync"
)

func updateDiagramNodeView(el core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	for _, changeNotification :

}

func init() {
	core.GetCore().AddFunction(coreDiagram.CrlDiagramNodeUri, updateDiagramNodeView)
}
