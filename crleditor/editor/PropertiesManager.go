package editor

import (
	"log"

	// "github.com/gopherjs/gopherjs/js"
	// "github.com/gopherjs/jquery"

	"github.com/pbrown12303/activeCRL/core"
)

// PropertiesManager is the manager of the properties display in the client
type PropertiesManager struct {
	crlEditor *CrlEditor
}

// NewPropertiesManager creates an instance of the PropertiesManager
func NewPropertiesManager(crlEditor *CrlEditor) *PropertiesManager {
	var newManager PropertiesManager
	newManager.crlEditor = crlEditor
	return &newManager
}

// ElementSelected updates the client's property display after an element selection change
func (pmPtr *PropertiesManager) ElementSelected(el core.Element, hl *core.HeldLocks) {
	if el == nil {
		log.Printf("ElementSelected called with nil argument")
		return
	}

	//	log.Printf("About to send ElementSelected notification")
	_, err := pmPtr.crlEditor.GetClientNotificationManager().SendNotification("ElementSelected", el.GetConceptID(hl), el, nil)
	if err != nil {
		log.Printf(err.Error())
	}
	//	log.Printf("ElementSelected response received")

}
