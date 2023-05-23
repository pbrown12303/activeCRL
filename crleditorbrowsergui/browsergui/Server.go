package browsergui

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"sync"

	//	"log"
	"net/http"
	"time"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
)

var server *http.Server
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

var templates = template.Must(template.ParseFiles(root+"crleditorbrowsergui/http/index.html", root+"crleditorbrowsergui/http/graph.html"))

// Exit is used as a programmatic shutdown of the server. It is primarily intended to support testing scenarios.
func Exit() error {
	// Save the settings
	err := BrowserGUISingleton.editor.SaveUserPreferences()
	if err != nil {
		log.Print(err.Error())
	}
	BrowserGUISingleton.editor.SetExitRequested()
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

// func loadPage(title string) (*page, error) {
// 	filename := root + "crlEditor/data/" + title + ".txt"
// 	body, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &page{Title: title, Body: body}, nil
// }

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

type requestHandler struct {
	sync.Mutex
	ready bool
}

var rh requestHandler

// handler for client requests
func (rh *requestHandler) handleRequest(w http.ResponseWriter, r *http.Request) {
	rh.Lock()
	defer rh.Unlock()
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

	trans, isNew := crleditor.CrlEditorSingleton.GetTransaction()
	if isNew {
		defer crleditor.CrlEditorSingleton.EndTransaction()
	}

	if CrlLogClientRequests || core.TraceChange {
		log.Printf("Received request: %#v", request)
	}
	switch request.Action {
	case "AbstractPointerChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		linkID, err := BrowserGUISingleton.getDiagramManager().abstractPointerChanged(
			request.RequestConceptID, request.AdditionalParameters["SourceID"], request.AdditionalParameters["TargetID"], trans)
		if err != nil {
			sendReply(w, 1, "Error processing AbstractPointerChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed AbstractPointerChanged", linkID, nil)
		}
	case "AddElementChild":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		el, _ := BrowserGUISingleton.GetUofD().NewElement(trans)
		el.SetLabel(BrowserGUISingleton.editor.GetDefaultElementLabel(), trans)
		el.SetOwningConceptID(request.RequestConceptID, trans)
		BrowserGUISingleton.editor.SelectElement(el, trans)
		sendReply(w, 0, "Processed AddElementChild", el.GetConceptID(trans), el)
	case "AddDiagramChild":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		diagram, err := crleditor.CrlEditorSingleton.GetDiagramManager().AddDiagram(request.RequestConceptID, trans)
		if err != nil {
			sendReply(w, 1, "Error processing AddDiagramChild: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed AddDiagramChild", diagram.GetConceptID(trans), diagram)
		}
	case "AddLiteralChild":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		el, _ := BrowserGUISingleton.GetUofD().NewLiteral(trans)
		el.SetLabel(BrowserGUISingleton.editor.GetDefaultLiteralLabel(), trans)
		el.SetOwningConceptID(request.RequestConceptID, trans)
		BrowserGUISingleton.editor.SelectElement(el, trans)
		sendReply(w, 0, "Processed AddLiteralChild", el.GetConceptID(trans), el)
	case "AddReferenceChild":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		el, _ := BrowserGUISingleton.GetUofD().NewReference(trans)
		el.SetLabel(BrowserGUISingleton.editor.GetDefaultReferenceLabel(), trans)
		el.SetOwningConceptID(request.RequestConceptID, trans)
		BrowserGUISingleton.editor.SelectElement(el, trans)
		sendReply(w, 0, "Processed AddReferenceChild", el.GetConceptID(trans), el)
	case "AddRefinementChild":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		el, _ := BrowserGUISingleton.GetUofD().NewRefinement(trans)
		el.SetLabel(BrowserGUISingleton.editor.GetDefaultRefinementLabel(), trans)
		el.SetOwningConceptID(request.RequestConceptID, trans)
		BrowserGUISingleton.editor.SelectElement(el, trans)
		sendReply(w, 0, "Processed AddRefinementChild", el.GetConceptID(trans), el)
	case "ClearWorkspace":
		err := BrowserGUISingleton.editor.ClearWorkspace(trans)
		reply(w, "ClearWorkspace", err)
	case "CloseWorkspace":
		err := BrowserGUISingleton.editor.CloseWorkspace(trans)
		reply(w, "CloseWorkspace", err)
	case "DefinitionChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		el := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetDefinition(request.AdditionalParameters["NewValue"], trans)
		}
		sendReply(w, 0, "Processed DefinitionChanged", request.RequestConceptID, el)
	case "DeleteDiagramElementView":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().deleteDiagramElementView(request.RequestConceptID, trans)
		if err != nil {
			sendReply(w, 1, err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed DeleteDiagramElementView", request.RequestConceptID, nil)
		}
	case "DiagramClick":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().diagramClick(request, trans)
		if err != nil {
			sendReply(w, 1, err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed DiagramClick", request.RequestConceptID, nil)
		}
	case "DiagramDrop":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().diagramDrop(request, trans)
		if err != nil {
			sendReply(w, 1, err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed DiagramDrop", request.RequestConceptID, nil)
		}
	case "DiagramElementSelected":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		elementID := request.RequestConceptID
		element := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if element != nil {
			modelElement := crldiagramdomain.GetReferencedModelElement(element, trans)
			BrowserGUISingleton.editor.SelectElement(modelElement, trans)
		}
		sendReply(w, 0, "Processed DiagramElementSelected", elementID, BrowserGUISingleton.GetUofD().GetElement(elementID))
	case "DiagramNodeNewPosition":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		x, err := strconv.ParseFloat(request.AdditionalParameters["NodeX"], 64)
		if err != nil {
			sendReply(w, 1, err.Error(), "", nil)
		}
		y, err2 := strconv.ParseFloat(request.AdditionalParameters["NodeY"], 64)
		if err2 != nil {
			sendReply(w, 1, err2.Error(), "", nil)
		} else {
			BrowserGUISingleton.getDiagramManager().setDiagramNodePosition(request.RequestConceptID, x, y, trans)
			sendReply(w, 0, "Processed DiagramNodeNewPosition", "", nil)
		}
	case "DiagramViewHasBeenClosed":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().DiagramViewHasBeenClosed(request.RequestConceptID, trans)
		reply(w, "DiagramViewHasBeenClosed", err)
	case "DisplayCallGraph":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.DisplayCallGraph(request.AdditionalParameters["GraphIndex"], trans)
		reply(w, "DisplayCallGraph", err)
	case "DisplayDiagramSelected":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := crleditor.CrlEditorSingleton.GetDiagramManager().DisplayDiagram(request.RequestConceptID, trans)
		reply(w, "Processed DisplayDiagramSelected", err)
	case "ElementPointerChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		linkID, err := BrowserGUISingleton.getDiagramManager().elementPointerChanged(
			request.RequestConceptID,
			request.AdditionalParameters["SourceID"],
			request.AdditionalParameters["TargetID"],
			request.AdditionalParameters["TargetAttributeName"], trans)
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
	case "FormatChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		diagramElement := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		var err error
		if diagramElement != nil {
			err = BrowserGUISingleton.diagramManager.formatChanged(diagramElement, request.AdditionalParameters["LineColor"], request.AdditionalParameters["BGColor"], trans)
		}
		if err == nil {
			sendReply(w, 0, "Processed FormatChanged", request.RequestConceptID, diagramElement)
		} else {
			sendReply(w, 1, "Error processing FormatChanged: "+err.Error(), "", nil)
		}
	case "InitializeClient":
		log.Printf("InitializeClient requested")
		sendReply(w, 0, "Client will be initialized", "", nil)
		for !rh.ready {
			time.Sleep(100 * time.Millisecond)
		}
		err := BrowserGUISingleton.InitializeGUI(trans)
		if err != nil {
			SendNotification("Error initializing client: "+err.Error(), "", nil, nil)
		} else {
			SendNotification("InitializationComplete", "", nil, nil)
		}
		requestInProgress = false
		log.Printf("requestInProgress set to false in InitializeClient")
	case "LabelChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		el := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetLabel(request.AdditionalParameters["NewValue"], trans)
		}
		sendReply(w, 0, "Processed LabelChanged", request.RequestConceptID, el)
	case "LiteralValueChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		el := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			switch typedEl := el.(type) {
			case core.Literal:
				typedEl.SetLiteralValue(request.AdditionalParameters["NewValue"], trans)
			}
		}
		sendReply(w, 0, "Processed LiteralValueChanged", request.RequestConceptID, el)
	case "NewDomainRequest":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		cs, _ := BrowserGUISingleton.GetUofD().NewElement(trans)
		cs.SetLabel(BrowserGUISingleton.editor.GetDefaultDomainLabel(), trans)
		BrowserGUISingleton.editor.SelectElement(cs, trans)
		sendReply(w, 0, "Processed NewDomainRequest", cs.GetConceptID(trans), cs)
	case "NullifyReferencedConcept":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		diagramElement := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if diagramElement == nil {
			sendReply(w, 1, "Selected diagram element not found", "", nil)
			break
		}
		modelElement := crldiagramdomain.GetReferencedModelElement(diagramElement, trans)
		if modelElement == nil {
			sendReply(w, 1, "Model element corresponding to selected diagram element not found", "", nil)
			break
		}
		err := BrowserGUISingleton.nullifyReferencedConcept(modelElement.GetConceptID(trans), trans)
		reply(w, "NullifyReferencedConcept", err)
	case "OpenWorkspace":
		err := BrowserGUISingleton.editor.OpenWorkspace()
		if err != nil {
			sendReply(w, 1, "Error processing OpenWorkspace: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed OpenWorkspace", "", nil)
		}
	case "OwnerPointerChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		linkID, err := BrowserGUISingleton.getDiagramManager().ownerPointerChanged(
			request.RequestConceptID, request.AdditionalParameters["SourceID"], request.AdditionalParameters["TargetID"], trans)
		if err != nil {
			sendReply(w, 1, "Error processing OwnerPointerChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed OwnerPointerChanged", linkID, nil)
		}
	case "Redo":
		err := BrowserGUISingleton.editor.Redo(trans)
		reply(w, "Redo", err)
	case "RefinedPointerChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		linkID, err := BrowserGUISingleton.getDiagramManager().refinedPointerChanged(
			request.RequestConceptID, request.AdditionalParameters["SourceID"], request.AdditionalParameters["TargetID"], trans)
		if err != nil {
			sendReply(w, 1, "Error processing RefinedPointerChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed RefinedPointerChanged", linkID, nil)
		}
	case "ReferenceLinkChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		linkID, err := BrowserGUISingleton.getDiagramManager().ReferenceLinkChanged(
			request.RequestConceptID,
			request.AdditionalParameters["SourceID"],
			request.AdditionalParameters["TargetID"],
			request.AdditionalParameters["TargetAttributeName"], trans)
		if err != nil {
			sendReply(w, 1, "Error processing ReferenceLinkChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ReferenceLinkChanged", linkID, nil)
		}
	case "RefinementLinkChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		linkID, err := BrowserGUISingleton.getDiagramManager().RefinementLinkChanged(
			request.RequestConceptID, request.AdditionalParameters["SourceID"], request.AdditionalParameters["TargetID"], trans)
		if err != nil {
			sendReply(w, 1, "Error processing RefinementLinkChanged: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed RefinementLinkChanged", linkID, nil)
		}
	case "RefreshDiagram":
		err := BrowserGUISingleton.getDiagramManager().refreshDiagramUsingURI(request.RequestConceptID, trans)
		reply(w, "RefreshDiagram", err)
	case "ReturnAvailableGraphCount":
		count := BrowserGUISingleton.GetAvailableGraphCount()
		reply := newReply()
		reply.Result = 0
		reply.ResultDescription = "Processed ReturnAvailableGraphCount"
		reply.ResultConceptID = ""
		reply.ResultConcept = nil
		reply.AdditionalParameters = map[string]string{"NumberOfAvailableGraphs": strconv.Itoa(count)}
		err := json.NewEncoder(w).Encode(reply)
		if err != nil {
			log.Print(err.Error())
		}
	case "SaveWorkspace":
		err := BrowserGUISingleton.editor.SaveWorkspace(trans)
		if err != nil {
			sendReply(w, 1, "SaveWorkspace failed: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed SaveWorkspace", "", nil)
		}
	case "ShowConceptInNavigator":
		requestedConcept := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if requestedConcept == nil {
			sendReply(w, 1, "Selected concept not found", "", nil)
			break
		}
		err := BrowserGUISingleton.ShowConceptInTree(requestedConcept, trans)
		if err != nil {
			sendReply(w, 1, "ShowConceptInNavigator failed: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ShowConceptInNavigator", "", nil)
		}
	case "ShowModelConceptInNavigator":
		diagramElement := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if diagramElement == nil {
			sendReply(w, 1, "Selected diagram element not found", "", nil)
			break
		}
		modelElement := crldiagramdomain.GetReferencedModelElement(diagramElement, trans)
		if modelElement == nil {
			sendReply(w, 1, "Model element corresponding to selected diagram element not found", "", nil)
			break
		}
		err := BrowserGUISingleton.ShowConceptInTree(modelElement, trans)
		if err != nil {
			sendReply(w, 1, "ShowModelConceptInNavigator failed: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ShowModelConceptInNavigator", "", nil)
		}
	case "ShowDiagramElementInNavigator":
		diagramElement := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if diagramElement == nil {
			sendReply(w, 1, "Selected diagram element not found", "", nil)
			break
		}
		err := BrowserGUISingleton.ShowConceptInTree(diagramElement, trans)
		if err != nil {
			sendReply(w, 1, "ShowDiagramElementInNavigator failed: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ShowDiagramElementInNavigator", "", nil)
		}
	case "SetTreeDragSelection":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		BrowserGUISingleton.SetTreeDragSelection(request.RequestConceptID)
		sendReply(w, 0, "Processed SetTreeDragSelection", request.RequestConceptID, BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID))
	case "ShowAbstractConcept":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().showAbstractConcept(request.RequestConceptID, trans)
		reply(w, "ShowAbstractConcept", err)
	case "ShowOwnedConcepts":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().showOwnedConcepts(request.RequestConceptID, trans)
		if err != nil {
			sendReply(w, 1, "ShowOwner failed: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ShowOwner", "", nil)
		}
	case "ShowOwner":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().showOwner(request.RequestConceptID, trans)
		if err != nil {
			sendReply(w, 1, "ShowOwner failed: "+err.Error(), "", nil)
		} else {
			sendReply(w, 0, "Processed ShowOwner", "", nil)
		}
	case "ShowReferencedConcept":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().showReferencedConcept(request.RequestConceptID, trans)
		reply(w, "ShowReferencedConcept", err)
	case "ShowRefinedConcept":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		err := BrowserGUISingleton.getDiagramManager().showRefinedConcept(request.RequestConceptID, trans)
		reply(w, "ShowRefinedConcept", err)
	case "TreeNodeDelete":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		elementID := request.RequestConceptID
		// log.Printf("TreeNodeDelete called for node id: %s for elementID: %s", request.RequestConceptID, elementID)
		err := BrowserGUISingleton.editor.DeleteElement(elementID, trans)
		if err == nil {
			sendReply(w, 0, "Element has been deleted", elementID, nil)
		} else {
			sendReply(w, 1, "Delete failed", elementID, nil)
		}
	case "TreeNodeSelected":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		elementID := request.RequestConceptID
		if CrlLogClientNotifications {
			log.Printf("Selected node id: %s", request.RequestConceptID)
		}
		BrowserGUISingleton.editor.SelectElementUsingIDString(elementID, trans)
		sendReply(w, 0, "Processed TreeNodeSelected", elementID, BrowserGUISingleton.GetUofD().GetElement(elementID))
	case "Undo":
		err := BrowserGUISingleton.editor.Undo(trans)
		reply(w, "Undo", err)
	case "UpdateDebugSettings":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		BrowserGUISingleton.UpdateDebugSettings(request)
		sendReply(w, 0, "Processed UpdateDebugSettings", "", nil)
	case "UpdateUserPreferences":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		BrowserGUISingleton.UpdateUserPreferences(request, trans)
		sendReply(w, 0, "Processed UpdateUserPreferences", "", nil)
	case "URIChanged":
		BrowserGUISingleton.GetUofD().MarkUndoPoint()
		el := BrowserGUISingleton.GetUofD().GetElement(request.RequestConceptID)
		if el != nil {
			el.SetURI(request.AdditionalParameters["NewValue"], trans)
		}
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
		log.Print(err.Error())
	}
	if CrlLogClientRequests {
		log.Printf("Sent reply: %#v", reply)
	}
}

// StartServer starts the editor server. This will automatically launch a browser as an interface
func (bgPtr *BrowserGUI) StartServer() {

	// RequestServer
	mux := http.NewServeMux()
	mux.HandleFunc("/index/graph.html", graphHandler)
	mux.HandleFunc("/index/", indexHandler)
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(root+"crleditorbrowsergui/http/js"))))
	mux.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir(root+"crleditorbrowsergui/http/images/icons"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(root+"crleditorbrowsergui/http/css"))))
	mux.HandleFunc("/index/request", rh.handleRequest)
	server = &http.Server{Addr: "127.0.0.1:8080", Handler: mux}
	go server.ListenAndServe()
	// WebSocketts server
	bgPtr.clientNotificationManager.startWsServer()
	if bgPtr.startBrowser {
		go openBrowser("http://localhost:8080/index")
	}
	bgPtr.checkReady() // Make sure the websocket on the client is ready
}

func (bgPtr *BrowserGUI) checkReady() {
	rh.ready = <-bgPtr.clientNotificationManager.webSocketReady
	BrowserGUISingleton.SetInitialized()
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
