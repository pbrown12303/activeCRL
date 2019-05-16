package editor

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pbrown12303/activeCRL/core"
)

// Notification is the data structure sent from the editor server to the browser client via websockets
type Notification struct {
	Notification          string
	NotificationConceptID string
	NotificationConcept   core.Element
	AdditionalParameters  map[string]string
}

// NewNotification returns an initialized Notification
func NewNotification() *Notification {
	var notification Notification
	notification.AdditionalParameters = make(map[string]string)
	return &notification
}

// NotificationResponse is the data structure returned by the browser client in response to a PushRequest
type NotificationResponse struct {
	Result          int
	ErrorMessage    string
	ResultConceptID string
	// ResultConcept core.Element
}

// ClientNotificationManager manages notification communications from server to client
type ClientNotificationManager struct {
	sync.Mutex
	conn *websocket.Conn
}

func newClientNotificationManager() *ClientNotificationManager {
	var cnMgr ClientNotificationManager
	return &cnMgr
}

func (cnMgr *ClientNotificationManager) setConnection(conn *websocket.Conn) {
	cnMgr.Lock()
	defer cnMgr.Unlock()
	cnMgr.conn = conn
}

// SendNotification sends the supplied Notification to the client and returns the Notification response.
// If there is no client connection or there is a problem in sending the Notification or receiving the NotificationResponse,
// an error is returned.
func (cnMgr *ClientNotificationManager) SendNotification(
	notificationDescription string, conceptID string, concept core.Element, params map[string]string) (*NotificationResponse, error) {
	cnMgr.Lock()
	defer cnMgr.Unlock()
	if cnMgr.conn == nil {
		return nil, nil
	}
	notification := NewNotification()
	notification.Notification = notificationDescription
	notification.NotificationConceptID = conceptID
	notification.NotificationConcept = concept
	notification.AdditionalParameters = params

	err := cnMgr.conn.WriteJSON(&notification)
	if err != nil {
		switch err.(type) {
		case *websocket.CloseError:
			log.Printf("WebSocket closed sending notification: %#v", notification)
			cnMgr.conn = nil
			return nil, fmt.Errorf("WebSocket closed sending notification: %#v", notification)
		default:
			log.Printf("Error: %s", err.Error())
			return nil, err
		}
	}
	if CrlLogClientNotifications {
		log.Printf("Sent notification: %#v", notification)
	}
	var notificationResponse NotificationResponse
	cnMgr.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	err = cnMgr.conn.ReadJSON(&notificationResponse)
	if err != nil {
		switch err.(type) {
		case *websocket.CloseError:
			log.Printf("WebSocket closed for response to notification: %#v", notification)
			cnMgr.conn = nil
			return nil, fmt.Errorf("WebSocket closed for response to notification: %#v", notification)
		default:
			log.Printf("Error %s in parsing response to WebSocket notification: %#v", err.Error(), notification)
			return nil, fmt.Errorf("Error %s in parsing response to WebSocket notification: %#v", err.Error(), notification)
		}
	}
	if CrlLogClientNotifications {
		log.Printf("Received notification response %#v", notificationResponse)
	}
	return &notificationResponse, nil
}

// SendNotification is a shortcut to the CrlEditorSingleton.GetClientNotificationManager().SendNotification() function
func SendNotification(notificationString string, conceptID string, concept core.Element, additionalParameters map[string]string) (*NotificationResponse, error) {
	return CrlEditorSingleton.GetClientNotificationManager().SendNotification(notificationString, conceptID, concept, additionalParameters)
}
