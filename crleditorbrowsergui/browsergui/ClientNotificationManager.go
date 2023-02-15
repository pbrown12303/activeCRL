package browsergui

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pbrown12303/activeCRL/core"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Notification is the data structure sent from the editor server to the browser client via websockets
type Notification struct {
	Notification             string
	NotificationConceptID    string
	NotificationConceptState core.ConceptState
	AdditionalParameters     map[string]string
}

// NewNotification returns an initialized Notification
func NewNotification() *Notification {
	var notification Notification
	notification.AdditionalParameters = make(map[string]string)
	return &notification
}

// NotificationResponse is the data structure returned by the browser client in response to a PushRequest
type NotificationResponse struct {
	Result          int    `json:"Result"`
	ErrorMessage    string `json:"ErrorMessage"`
	ResultConceptID string `json:"ResultConceptID"`
	BooleanValue    string `json:"BooleanValue"`
}

// ClientNotificationManager manages notification communications from server to client
type ClientNotificationManager struct {
	sync.Mutex
	conn           *websocket.Conn
	context        *context.Context
	wsServer       *http.Server
	webSocketReady chan bool
}

func newClientNotificationManager() *ClientNotificationManager {
	var cnMgr ClientNotificationManager
	cnMgr.webSocketReady = make(chan bool)
	return &cnMgr
}

func (cnMgr *ClientNotificationManager) setConnection(conn *websocket.Conn, context *context.Context) {
	cnMgr.Lock()
	defer cnMgr.Unlock()
	cnMgr.conn = conn
}

// SendNotification sends the supplied Notification to the client and returns the Notification response.
// If there is no client connection or there is a problem in sending the Notification or receiving the NotificationResponse,
// an error is returned.
func (cnMgr *ClientNotificationManager) SendNotification(
	notificationDescription string, conceptID string, conceptState *core.ConceptState, params map[string]string) (*NotificationResponse, error) {
	cnMgr.Lock()
	defer cnMgr.Unlock()
	if cnMgr.conn == nil {
		return nil, nil
	}
	notification := NewNotification()
	notification.Notification = notificationDescription
	notification.NotificationConceptID = conceptID
	if conceptState != nil {
		notification.NotificationConceptState = *conceptState
	}
	notification.AdditionalParameters = params

	ctx := context.Background()

	err := wsjson.Write(ctx, cnMgr.conn, &notification)
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
	err = wsjson.Read(ctx, cnMgr.conn, &notificationResponse)
	if err != nil {
		switch err.(type) {
		case *websocket.CloseError:
			log.Printf("WebSocket closed for response to notification: %#v", notification)
			cnMgr.conn = nil
			return nil, fmt.Errorf("WebSocket closed for response to notification: %#v", notification)
		default:
			log.Printf("Error %s in parsing response to WebSocket notification: %#v", err.Error(), notification)
			return nil, fmt.Errorf("error %s in parsing response to WebSocket notification: %#v", err.Error(), notification)
		}
	}
	if CrlLogClientNotifications {
		log.Printf("Received notification response %#v", notificationResponse)
	}
	return &notificationResponse, nil
}

// SendNotification is a shortcut to the BrowserGUISingleton.GetClientNotificationManager().SendNotification() function
func SendNotification(notificationString string, conceptID string, conceptState *core.ConceptState, additionalParameters map[string]string) (*NotificationResponse, error) {
	return BrowserGUISingleton.GetClientNotificationManager().SendNotification(notificationString, conceptID, conceptState, additionalParameters)
}

func (cnMgr *ClientNotificationManager) startWsServer() {
	// This function must be idempotent
	if cnMgr.wsServer == nil {
		wsMux := http.NewServeMux()
		wsMux.HandleFunc("/index/ws", cnMgr.wsHandler)
		cnMgr.wsServer = &http.Server{Addr: "127.0.0.1:8081", Handler: wsMux}
		go cnMgr.wsServer.ListenAndServe()
	}
}

// wsHandler is the handler for WebSocket Notifications
func (cnMgr *ClientNotificationManager) wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("wsHandler invoked")

	options := websocket.AcceptOptions{
		Subprotocols:         []string{},
		InsecureSkipVerify:   false,
		OriginPatterns:       []string{"localhost*"},
		CompressionMode:      0,
		CompressionThreshold: 0,
	}
	wsConnection, err := websocket.Accept(w, r, &options)
	if err != nil {
		log.Println(err)
		return
	}
	// We keep the socket open so that notifications can be sent to the browser
	// defer wsConnection.Close(websocket.StatusInternalError, "wsHandler exited")

	wsContext, _ := context.WithTimeout(r.Context(), time.Second*10)
	// defer cancel()

	BrowserGUISingleton.GetClientNotificationManager().setConnection(wsConnection, &wsContext)
	cnMgr.webSocketReady <- true

	log.Printf("wsHandler complete")

}
