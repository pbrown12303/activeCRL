package editor

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	//	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagram"
)

var home string
var server *http.Server
var wsServer *http.Server
var webSocketReady = make(chan bool)
var requestInProgress bool

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

// Reply is the data structure returned by the editor server in response to a Request
type Reply struct {
	Result               int
	ResultDescription    string
	ResultConceptID      string
	ResultConcept        core.Element
	AdditionalParameters map[string]string
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

var templates = template.Must(template.ParseFiles(root+"crleditor/http/index.html", root+"crleditor/http/graph.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Exit is used as a programmatic shutdown of the server. It is primarily intended to support testing scenarios.
func Exit() error {
	// Save the settings
	err := CrlEditorSingleton.saveUserPreferences()
	if err != nil {
		log.Printf(err.Error())
	}
	err = server.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// GetRequestInProgress returns true if the server is actively processing a request
func GetRequestInProgress() bool {
	return requestInProgress
}

func graphHandler(w http.ResponseWriter, r *http.Request) {
	p := &page{Title: "Function Call Notification Graph"}
	renderTemplate(w, "graph", p)
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
	requestInProgress = true
	// log.Printf("requestInProgress set to true")
	defer func() {
		requestInProgress = false
		// log.Printf("requestInProgress set to false")
	}()
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
	if CrlLogClientRequests {
		log.Printf("Received request: %#v", request)
	}
	switch request.Action {
	case "AbstractPointerChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		linkID, err := CrlEditorSingleton.getDiagramManager().abstractPointerChanged(
			request.RequestConceptID, request.AdditionalParameters["SourceID"], request.AdditionalParameters["TargetID"], hl)
		if err != nil {
			sendReply(w, 1, "Error processing AbstractPointerChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed AbstractPointerChanged", linkID, nil)
		}
	case "AddElementChild":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el, _ := CrlEditorSingleton.GetUofD().NewElement(hl)
		el.SetLabel(CrlEditorSingleton.getDefaultElementLabel(), hl)
		el.SetOwningConceptID(request.RequestConceptID, hl)
		CrlEditorSingleton.SelectElement(el, hl)
		sendReply(w, 0, "Processed AddElementChild", el.GetConceptID(hl), el)
	case "AddDiagramChild":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		diagramManager := CrlEditorSingleton.getDiagramManager()
		diagram := diagramManager.newDiagram(hl)
		diagram.SetOwningConceptID(request.RequestConceptID, hl)
		CrlEditorSingleton.SelectElement(diagram, hl)
		hl.ReleaseLocksAndWait()
		err := diagramManager.displayDiagram(diagram, hl)
		hl.ReleaseLocksAndWait()
		if err != nil {
			sendReply(w, 1, "Error processing AddDiagramChild: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed AddDiagramChild", diagram.GetConceptID(hl), diagram)
		}
	case "AddLiteralChild":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el, _ := CrlEditorSingleton.GetUofD().NewLiteral(hl)
		el.SetLabel(CrlEditorSingleton.getDefaultLiteralLabel(), hl)
		el.SetOwningConceptID(request.RequestConceptID, hl)
		CrlEditorSingleton.SelectElement(el, hl)
		sendReply(w, 0, "Processed AddLiteralChild", el.GetConceptID(hl), el)
	case "AddReferenceChild":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el, _ := CrlEditorSingleton.GetUofD().NewReference(hl)
		el.SetLabel(CrlEditorSingleton.getDefaultReferenceLabel(), hl)
		el.SetOwningConceptID(request.RequestConceptID, hl)
		CrlEditorSingleton.SelectElement(el, hl)
		sendReply(w, 0, "Processed AddReferenceChild", el.GetConceptID(hl), el)
	case "AddRefinementChild":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el, _ := CrlEditorSingleton.GetUofD().NewRefinement(hl)
		el.SetLabel(CrlEditorSingleton.getDefaultRefinementLabel(), hl)
		el.SetOwningConceptID(request.RequestConceptID, hl)
		CrlEditorSingleton.SelectElement(el, hl)
		sendReply(w, 0, "Processed AddRefinementChild", el.GetConceptID(hl), el)
	case "ClearWorkspace":
		err := CrlEditorSingleton.ClearWorkspace(hl)
		reply(w, "ClearWorkspace", err)
	case "CloseWorkspace":
		err := CrlEditorSingleton.CloseWorkspace(hl)
		reply(w, "CloseWorkspace", err)
	case "DefinitionChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetDefinition(request.AdditionalParameters["NewValue"], hl)
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed DefinitionChanged", request.RequestConceptID, el)
	case "DeleteDiagramElementView":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		err := CrlEditorSingleton.getDiagramManager().deleteDiagramElementView(request.RequestConceptID, hl)
		if err != nil {
			sendReply(w, 1, err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed DeleteDiagramElementView", request.RequestConceptID, nil)
		}
	case "DiagramClick":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		err := CrlEditorSingleton.getDiagramManager().diagramClick(request, hl)
		if err != nil {
			sendReply(w, 1, err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed DiagramClick", request.RequestConceptID, nil)
		}
	case "DiagramDrop":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		err := CrlEditorSingleton.getDiagramManager().diagramDrop(request, hl)
		if err != nil {
			sendReply(w, 1, err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed DiagramDrop", request.RequestConceptID, nil)
		}
	case "DiagramNodeNewPosition":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		x, err := strconv.ParseFloat(request.AdditionalParameters["NodeX"], 64)
		if err != nil {
			sendReply(w, 1, err.Error(), "", nil)
		}
		y, err2 := strconv.ParseFloat(request.AdditionalParameters["NodeY"], 64)
		if err2 != nil {
			sendReply(w, 1, err2.Error(), "", nil)
		} else {
			CrlEditorSingleton.getDiagramManager().setDiagramNodePosition(request.RequestConceptID, x, y, hl)
			sendReply(w, 0, "Processed DiagramNodeNewPosition", "", nil)
		}
	case "DiagramElementSelected":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		elementID := request.RequestConceptID
		element := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if element != nil {
			modelElement := crldiagram.GetReferencedModelElement(element, hl)
			CrlEditorSingleton.SelectElement(modelElement, hl)
		}
		sendReply(w, 0, "Processed DiagramElementSelected", elementID, CrlEditorSingleton.GetUofD().GetElement(elementID))
	case "DiagramViewHasBeenClosed":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		CrlEditorSingleton.getDiagramManager().DiagramViewHasBeenClosed(request.RequestConceptID, hl)
	case "DisplayCallGraph":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		err := CrlEditorSingleton.DisplayCallGraph(request.AdditionalParameters["GraphIndex"], hl)
		reply(w, "DisplayCallGraph", err)
	case "DisplayDiagramSelected":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil && el.IsRefinementOfURI(crldiagram.CrlDiagramURI, hl) {
			diagramManager := CrlEditorSingleton.getDiagramManager()
			diagramManager.displayDiagram(el, hl)
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed DisplayDiagramSelected", request.RequestConceptID, el)
	case "ElementPointerChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		linkID, err := CrlEditorSingleton.getDiagramManager().elementPointerChanged(
			request.RequestConceptID,
			request.AdditionalParameters["SourceID"],
			request.AdditionalParameters["TargetID"],
			request.AdditionalParameters["TargetAttributeName"], hl)
		if err != nil {
			sendReply(w, 1, "Error processing ElementPointerChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ElementPointerChanged", linkID, nil)
		}
	case "Exit":
		log.Printf("Exit requested")
		sendReply(w, 0, "Server will close", "", nil)
		time.Sleep(5 * time.Second)
		if err := Exit(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown error: %s", err.Error())
		}
	case "InitializeClient":
		log.Printf("InitializeClient requested")
		sendReply(w, 0, "Client will be initialized", "", nil)
		err := CrlEditorSingleton.InitializeClient()
		if err != nil {
			SendNotification("Error initializing client: "+err.Error(), "", nil, nil)
		} else {
			SendNotification("InitializationComplete", "", nil, nil)
		}
		requestInProgress = false
		log.Printf("requestInProgress set to false in InitializeClient")
	case "LabelChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetLabel(request.AdditionalParameters["NewValue"], hl)
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed LabelChanged", request.RequestConceptID, el)
	case "LiteralValueChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			switch el.(type) {
			case core.Literal:
				el.(core.Literal).SetLiteralValue(request.AdditionalParameters["NewValue"], hl)
			}
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed LiteralValueChanged", request.RequestConceptID, el)
	case "NewConceptSpaceRequest":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		cs, _ := CrlEditorSingleton.GetUofD().NewElement(hl)
		cs.SetLabel(CrlEditorSingleton.getDefaultConceptSpaceLabel(), hl)
		CrlEditorSingleton.SelectElement(cs, hl)
		sendReply(w, 0, "Processed NewConceptSpaceRequest", cs.GetConceptID(hl), cs)
	case "OpenWorkspace":
		err := CrlEditorSingleton.openWorkspace(hl)
		if err != nil {
			sendReply(w, 1, "Error processing OpenWorkspace: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed OpenWorkspace", "", nil)
		}
	case "OwnerPointerChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		linkID, err := CrlEditorSingleton.getDiagramManager().ownerPointerChanged(
			request.RequestConceptID, request.AdditionalParameters["SourceID"], request.AdditionalParameters["TargetID"], hl)
		if err != nil {
			sendReply(w, 1, "Error processing OwnerPointerChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed OwnerPointerChanged", linkID, nil)
		}
	case "Redo":
		err := CrlEditorSingleton.Redo(hl)
		reply(w, "Redo", err)
	case "RefinedPointerChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		linkID, err := CrlEditorSingleton.getDiagramManager().refinedPointerChanged(
			request.RequestConceptID, request.AdditionalParameters["SourceID"], request.AdditionalParameters["TargetID"], hl)
		if err != nil {
			sendReply(w, 1, "Error processing RefinedPointerChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed RefinedPointerChanged", linkID, nil)
		}
	case "ReferenceLinkChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		linkID, err := CrlEditorSingleton.getDiagramManager().ReferenceLinkChanged(
			request.RequestConceptID,
			request.AdditionalParameters["SourceID"],
			request.AdditionalParameters["TargetID"],
			request.AdditionalParameters["TargetAttributeName"], hl)
		if err != nil {
			sendReply(w, 1, "Error processing ReferenceLinkChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ReferenceLinkChanged", linkID, nil)
		}
	case "RefinementLinkChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		linkID, err := CrlEditorSingleton.getDiagramManager().RefinementLinkChanged(
			request.RequestConceptID, request.AdditionalParameters["SourceID"], request.AdditionalParameters["TargetID"], hl)
		if err != nil {
			sendReply(w, 1, "Error processing RefinementLinkChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed RefinementLinkChanged", linkID, nil)
		}
	case "RefreshDiagram":
		err := CrlEditorSingleton.getDiagramManager().refreshDiagramUsingURI(request.RequestConceptID, hl)
		reply(w, "RefreshDiagram", err)
	case "ReturnAvailableGraphCount":
		count := CrlEditorSingleton.GetAvailableGraphCount()
		reply := newReply()
		reply.Result = 0
		reply.ResultDescription = "Processed ReturnAvailableGraphCount"
		reply.ResultConceptID = ""
		reply.ResultConcept = nil
		reply.AdditionalParameters = map[string]string{"NumberOfAvailableGraphs": strconv.Itoa(count)}
		err := json.NewEncoder(w).Encode(reply)
		if err != nil {
			log.Printf(err.Error())
		}
	case "SaveWorkspace":
		err := CrlEditorSingleton.SaveWorkspace(hl)
		if err != nil {
			sendReply(w, 1, "SaveWorkspace failed: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed SaveWorkspace", "", nil)
		}
	case "SetTreeDragSelection":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		CrlEditorSingleton.SetTreeDragSelection(request.RequestConceptID)
		sendReply(w, 0, "Processed SetTreeDragSelection", request.RequestConceptID, CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID))
	case "ShowAbstractConcept":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		err := CrlEditorSingleton.getDiagramManager().showAbstractConcept(request.RequestConceptID, hl)
		reply(w, "ShowAbstractConcept", err)
	case "ShowOwner":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		err := CrlEditorSingleton.getDiagramManager().showOwner(request.RequestConceptID, hl)
		if err != nil {
			sendReply(w, 1, "ShowOwner failed: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ShowOwner", "", nil)
		}
	case "ShowReferencedConcept":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		err := CrlEditorSingleton.getDiagramManager().showReferencedConcept(request.RequestConceptID, hl)
		reply(w, "ShowReferencedConcept", err)
	case "ShowRefinedConcept":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		err := CrlEditorSingleton.getDiagramManager().showRefinedConcept(request.RequestConceptID, hl)
		reply(w, "ShowRefinedConcept", err)
	case "TreeNodeDelete":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		elementID := request.RequestConceptID
		// log.Printf("TreeNodeDelete called for node id: %s for elementID: %s", request.RequestConceptID, elementID)
		err := CrlEditorSingleton.Delete(elementID)
		if err == nil {
			sendReply(w, 0, "Element has been deleted", elementID, nil)
		} else {
			sendReply(w, 1, "Delete failed", elementID, nil)
		}
	case "TreeNodeSelected":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		elementID := request.RequestConceptID
		if CrlLogClientNotifications {
			log.Printf("Selected node id: %s", request.RequestConceptID)
		}
		CrlEditorSingleton.SelectElementUsingIDString(elementID, hl)
		sendReply(w, 0, "Processed TreeNodeSelected", elementID, CrlEditorSingleton.GetUofD().GetElement(elementID))
	case "Undo":
		err := CrlEditorSingleton.Undo(hl)
		reply(w, "Undo", err)
	case "UpdateDebugSettings":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		CrlEditorSingleton.UpdateDebugSettings(request)
		sendReply(w, 0, "Processed UpdateDebugSettings", "", nil)
	case "UpdateUserPreferences":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		CrlEditorSingleton.UpdateUserPreferences(request, hl)
		sendReply(w, 0, "Processed UpdateUserPreferences", "", nil)
	case "URIChanged":
		CrlEditorSingleton.uOfD.MarkUndoPoint()
		el := CrlEditorSingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetURI(request.AdditionalParameters["NewValue"], hl)
		}
		hl.ReleaseLocksAndWait()
		sendReply(w, 0, "Processed URI changed", request.RequestConceptID, el)
	default:
		log.Printf("Unhandled request: %s", request.Action)
		sendReply(w, 1, "Unhandled request: "+request.Action, "", nil)
	}
}

func reply(w http.ResponseWriter, requestName string, err error) {
	if err != nil {
		sendReply(w, 1, requestName+": "+err.Error(), "", nil)
	} else {
		sendReply(w, 0, " Processed "+requestName, "", nil)
	}
}

func sendReply(w http.ResponseWriter, code int, message string, resultConceptID string, resultConcept core.Element) {
	reply := newReply()
	reply.Result = code
	reply.ResultDescription = message
	reply.ResultConceptID = resultConceptID
	reply.ResultConcept = resultConcept
	err := json.NewEncoder(w).Encode(reply)
	if err != nil {
		log.Printf(err.Error())
	}
}

// StartServer starts the editor server. This will automatically launch a browser as an interface
func StartServer(startBrowser bool, workspaceArg string, userFolderArg string) {
	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		log.Fatalf("User home directory not found")
	}
	InitializeCrlEditorSingleton(userFolderArg)
	err = CrlEditorSingleton.LoadUserPreferences(workspaceArg)
	if err != nil {
		log.Fatal(err)
	}
	workspacePath := workspaceArg
	if workspacePath == "" {
		workspacePath = CrlEditorSingleton.GetUserPreferences().WorkspacePath
	}
	if workspacePath == "" {
		workspacePath, err2 := CrlEditorSingleton.SelectWorkspace()
		if err2 != nil {
			log.Fatal(err2)
		}
		err = CrlEditorSingleton.SetWorkspacePath(workspacePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	hl := CrlEditorSingleton.GetUofD().NewHeldLocks()
	err = CrlEditorSingleton.LoadWorkspace(hl)
	if err != nil {
		log.Fatal(err)
	}
	hl.ReleaseLocksAndWait()
	CrlEditorSingleton.uOfD.SetRecordingUndo(true)

	go CrlEditorSingleton.InitializeClient()
	// WebSocketts server
	go startWsServer()
	// RequestServer
	mux := http.NewServeMux()
	mux.HandleFunc("/index/graph.html", graphHandler)
	mux.HandleFunc("/index/", indexHandler)
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(root+"crleditor/http/js"))))
	mux.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir(root+"crleditor/http/images/icons"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(root+"crleditor/http/css"))))
	mux.HandleFunc("/index/request", requestHandler)

	if startBrowser == true {
		openBrowser("http://localhost:8082/index")
	}

	server = &http.Server{Addr: "127.0.0.1:8082", Handler: mux}

	server.ListenAndServe()
}

func startWsServer() {
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/index/ws", wsHandler)
	wsServer = &http.Server{Addr: "127.0.0.1:8081", Handler: wsMux}
	wsServer.ListenAndServe()
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
	// TODO: Fix the upgrader.CheckOrigin() to do something intelligent
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	wsConnection, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	CrlEditorSingleton.GetClientNotificationManager().setConnection(wsConnection)
	log.Printf("wsHandler complete")
	webSocketReady <- true
}
