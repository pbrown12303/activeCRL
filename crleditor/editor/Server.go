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
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagram"
)

var server *http.Server
var webSocketReady = make(chan bool)

// Request is the data structure submitted by the client browser as part of an http request
type Request struct {
	Action               string
	AdditionalParameters map[string]string
	RequestConceptID     string
	RequestConcept       core.Element
}

func newRequest() *Request {
	var request Request
	request.AdditionalParameters = make(map[string]string)
	return &request
}

// Reply is the data structure returned by the editor server in response to an Request
type Reply struct {
	Result            int
	ResultDescription string
	ResultConceptID   string
	ResultConcept     core.Element
}

func newReply() *Reply {
	var reply Reply
	return &reply
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

// Exit is used as a programmatic shutdown of the server. It is primarily intended to support testing scenarios.
func Exit() {
	server.Shutdown(context.Background())
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
	request := newRequest()
	if r.Body == nil {
		log.Printf("Request received with no body")
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Request received but JSON decoding of body failed")
		http.Error(w, err.Error(), 400)
		return
	}
	hl := CrlEditorSingleton.GetUofD().NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	switch request.Action {
	case "DefinitionChanged":
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetDefinition(request.AdditionalParameters["NewValue"], hl)
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed DefinitionChanged", request.RequestConceptID, el)
	case "DiagramDrop":
		CrlEditorSingleton.GetDiagramManager().DiagramDrop(request, hl)
		sendReply(w, 0, "Processed DiagramDrop", request.RequestConceptID, nil)
	case "DiagramNodeSelected":
		// TODO: Finish DiagramNodeSelected
		sendReply(w, 0, "Processed DiagramNodeSelected", "", nil)
	case "DisplayDiagramSelected":
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil && el.IsRefinementOfURI(crldiagram.CrlDiagramURI, hl) {
			diagramManager := CrlEditorSingleton.GetDiagramManager()
			diagramManager.DisplayDiagram(el, hl)
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed DisplayDiagramSelected", request.RequestConceptID, el)
	case "Exit":
		log.Printf("Exit requested")
		sendReply(w, 0, "Server will close", "", nil)
		time.Sleep(5 * time.Second)
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
	case "InitializeClient":
		log.Printf("InitializeClient requested")
		sendReply(w, 0, "Client will be initialized", "", nil)
		InitializeClient()
		CrlEditorSingleton.GetClientNotificationManager().SendNotification("InitializationComplete", "", nil, nil)
	case "LabelChanged":
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetLabel(request.AdditionalParameters["NewValue"], hl)
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed LabelChanged", request.RequestConceptID, el)
	case "LiteralValueChanged":
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			switch el.(type) {
			case core.Literal:
				el.(core.Literal).SetLiteralValue(request.AdditionalParameters["NewValue"], hl)
			}
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed LiteralValueChanged", request.RequestConceptID, el)
	case "URIChanged":
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetURI(request.AdditionalParameters["NewValue"], hl)
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed URI changed", request.RequestConceptID, el)
	case "NewDiagramRequest":
		diagramManager := CrlEditorSingleton.GetDiagramManager()
		diagram := diagramManager.NewDiagram(hl)
		hl.ReleaseLocksAndWait()
		diagramManager.DisplayDiagram(diagram, hl)
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed NewDiagramRequest", diagram.GetConceptID(hl), diagram)
	case "SetTreeDragSelection":
		CrlEditorSingleton.SetTreeDragSelection(request.RequestConceptID)
		sendReply(w, 0, "Processed SetTreeDragSelection", request.RequestConceptID, CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID))
	case "TreeNodeDelete":
		elementID := request.RequestConceptID
		log.Printf("TreeNodeDelete called for node id: %s for elementID: %s", request.RequestConceptID, elementID)
		err := CrlEditorSingleton.Delete(elementID)
		if err == nil {
			sendReply(w, 0, "Element has been deleted", elementID, nil)
		} else {
			sendReply(w, 1, "Delete failed", elementID, nil)
		}
	case "TreeNodeSelected":
		elementID := request.RequestConceptID
		log.Printf("Selected node id: %s", request.RequestConceptID)
		CrlEditorSingleton.SelectElementUsingIDString(elementID, hl)
		sendReply(w, 0, "Element has been selected", elementID, CrlEditorSingleton.GetUofD().GetElement(elementID))
	default:
		log.Printf("Unhandled request: %s", request.Action)
		sendReply(w, 1, "Unhandled request: "+request.Action, "", nil)
	}
}

func sendReply(w http.ResponseWriter, code int, message string, resultConceptID string, resultConcept core.Element) {
	reply := newReply()
	reply.Result = code
	reply.ResultDescription = message
	reply.ResultConceptID = resultConceptID
	reply.ResultConcept = resultConcept
	json.NewEncoder(w).Encode(reply)
}

// StartServer starts the editor server. This will automatically launch a browser as an interface
func StartServer(startBrowser bool) {
	InitializeCrlEditorSingleton()
	mux := http.NewServeMux()
	mux.HandleFunc("/index/", indexHandler)
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(root+"crleditor/http/js"))))
	mux.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir(root+"crleditor/http/images/icons"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(root+"crleditor/http/css"))))
	mux.HandleFunc("/index/ws", wsHandler)
	mux.HandleFunc("/index/request", requestHandler)

	if startBrowser == true {
		openBrowser("http://localhost:8080/index")
	}

	server = &http.Server{Addr: "127.0.0.1:8080", Handler: mux}

	server.ListenAndServe()
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var wsConnection *websocket.Conn

// wsHandler is the handler for WebSocket Notifications
func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("wsHandler invoked")
	var err error
	wsConnection, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	CrlEditorSingleton.GetClientNotificationManager().setConnection(wsConnection)
	log.Printf("wsHandler complete")
	webSocketReady <- true
}
