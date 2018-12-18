package editor

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"

	//	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

var server *http.Server
var notificationChannel chan Notification
var notificationResponseChannel chan NotificationResponse

// Request is the data structure submitted by the client browser as part of an http request
type Request struct {
	Action           string
	RequestConceptID string
	//	RequestConcept core.Element
}

// Reply is the data structure returned by the editor server in response to an Request
type Reply struct {
	Result            int
	ResultDescription string
	ResultConceptID   string
	//	ResultConcept core.Element
}

// Notification is the data structure sent from the editor server to the browser client via websockets
type Notification struct {
	Notification          string
	NotificationConceptID string
	// PushConcept core.Element
}

// NotificationResponse is the data structure returned by the browser client in response to a PushRequest
type NotificationResponse struct {
	Result          int
	ResultConceptID string
	// ResultConcept core.Element
}

type page struct {
	Title string
	Body  []byte
}

var root = "C:/GoWorkspace/src/github.com/pbrown12303/activeCRL/"

var templates = template.Must(template.ParseFiles(root + "crleditor/http/index.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &page{Title: "CRL Editor"}
	renderTemplate(w, "index", p)
}

func loadPage(title string) (*page, error) {
	filename := root + "crlEditor/data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &page{Title: title, Body: body}, nil
}

func notificationResponder(conn *websocket.Conn) {
	for {
		var notificationResponse NotificationResponse
		err := conn.ReadJSON(&notificationResponse)
		if err != nil {
			log.Println("Error: ", err.Error())
			return
		}
		notificationResponseChannel <- notificationResponse
	}
}

func notificationSender(conn *websocket.Conn) {
	for {
		notification := <-notificationChannel
		err := conn.WriteJSON(&notification)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// openBrowser tries to open the URL in a browser,
// and returns whether it succeed in doing so.
func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

// handler for client requests
func requestHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	switch request.Action {
	case "Exit":
		var reply Reply
		reply.Result = 0
		reply.ResultDescription = "Server will close"
		json.NewEncoder(w).Encode(reply)
		time.Sleep(5 * time.Second)
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}

	}
}

// StartServer starts the editor server. This will automatically launch a browser as an interface
func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/index/", indexHandler)
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(root+"crleditor/http/js"))))
	mux.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir(root+"crleditor/http/images/icons"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(root+"crleditor/http/css"))))
	mux.HandleFunc("index/ws", wsHandler)
	mux.HandleFunc("/index/request", requestHandler)

	openBrowser("http://localhost:8080/index")

	server = &http.Server{Addr: "127.0.0.1:8080", Handler: mux}
	server.ListenAndServe()
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// wsHandler is the handler for WebSocket Notifications
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go notificationResponder(conn)
	go notificationSender(conn)
}
