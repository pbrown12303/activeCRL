// Current State Variables
// crlCurrentDiagramContainerIDGlobal is the identifier for the diagram container currently being displayed
var crlCurrentDiagramContainerID;
// crlCurrentToolbarButton is the id of the last toolbar button pressed
var crlCurrentToolbarButton;
// crlInitializationCompleteGlobal indicates whether the server-side initialization has been completed
var crlInitializationComplete = false;
// crlMovedNodes is an array of nodes that have been moved. This is a temporary cache that is used to update the 
// server once a mouse up has occurred
var crlMovedNodes = {};
// crlSelectedConceptIDGlobal contains the model identifier of the currently selected concept
var crlSelectedConceptID;
// crlTreeDragSelectionIDGlobal contains the model identifier of the concept currently being dragged from the tree
var crlTreeDragSelectionID;
// crlLineColor is the last copied color
var crlLineColor = ""
// crlBGColor is the last copied background color
var crlBGColor = ""
// crlSelectedDiagramElementIDForFormat is the diagram element whose format is being edited
var crlSelectedDiagramElementIDForFormat;

// Editor Settings
// crlDropReferenceAsLink when true causes references dragged from the tree into the diagram to be added as links
var crlDropReferenceAsLink = false;
// crlDropRefinementAsLink when true causes refinements dragged from the tree into the diagram to be added as links
var crlDropRefinementAsLink = false;
// crlEnableTracing is the client-side copy of the server-side value that turns on notification tracing
var crlEnableTracing = false;
// crlOmitHousekeepingCalls indicates whether housekeeping calls shouldl be included when tracing is enabled
var crlOmitHousekeepingCalls = false;
// crlOmitManageTreeNodesCalls indicates whether housekeeping calls shouldl be included when tracing is enabled
var crlOmitManageTreeNodesCalls = false;
// crlOmitDiagramRelatedCalls indicates whether housekeeping calls shouldl be included when tracing is enabled
var crlOmitDiagramRelatedCalls = false;
// crlAutomatedTestInProgress set to true during automated regression testing. Primary intent is to suppress alerts requiring user response
var crlAutomatedTestInProgress = false;

// Debug Settings
// crlDebugSettingsDialog is the initialized dialog used for editing debug settings
var crlDebugSettingsDialog;

// Dialogs
// crlUserPreferencesDialog is the initialized dialog used for editing user preferences
var crlUserPreferencesDialog;
// crlOpenWorkspaceDialog is the initialized dialog used for opening a workspace
var crlOpenWorkspaceDialog;
// crlDiagramElementFormatDialog is the initialized dialog used to set a diagram element format
var crlDiagramElementFormatDialog;
// crlSelectConceptByIDDialog is a dialog for entering the ConceptID to be selected
var crlSelectConceptByIDDialog;

// Lookup structures
// crlGraphsGlobal is an array of existing graphs that is used to look up a graph given its identifier
var crlGraphsGlobal = {};
// crlPaperGlobal is an array of existing papers that is used to look up a paper given its identifier
var crlPapersGlobal = {};
// CrlWebSocketGlobal is the web socket being used for server-side communications
var crlWebsocketGlobal;
// crlWorkspacePath is the path to the current workspace
var crlWorkspacePath;

var crlCursorToolbarButtonID = "cursorToolbarButton";
var crlElementToolbarButtonID = "elementToolbarButton";
var crlLiteralToolbarButtonID = "literalToolbarButton";
var crlReferenceToolbarButtonID = "referenceToolbarButton";
var crlReferenceLinkToolbarButtonID = "referenceLinkToolbarButton";
var crlRefinementToolbarButtonID = "refinementToolbarButton";
var crlRefinementLinkToolbarButtonID = "refinementLinkToolbarButton";
var crlDiagramToolbarButtonID = "diagramToolbarButton";
var crlOwnerPointerToolbarButtonID = "ownerPointerToolbarButton";
var crlElementPointerToolbarButtonID = "elementPointerToolbarButton";
var crlAbstractPointerToolbarButtonID = "abstractPointerToolbarButton";
var crlRefinedPointerToolbarButtonID = "refinedPointerToolbarButton";

var crlDiagramCellDropdownMenu = null;
var crlDiagramTabDropdownMenu = null;

var crlInCrlElementSelected = false;

var crlKeyPressed = { 16: false };
var crlMouseButtonPressed = { 0: false };
var crlMousePosition = { "x": 0, "y": 0 }

// Initialize
$(function () {
    $(".uofd-browser").resizable({
        resizeHeight: false
    });
    crlInitializeTree();
    $("#body").on("ondrop", crlOnEditorDrop);
    crlPopulateToolbar();
    crlDiagramCellDropdownMenu = document.getElementById("diagramCellDropdown");
    if (crlDiagramCellDropdownMenu) {
        crlDiagramCellDropdownMenu.addEventListener("mouseleave", function () {
            crlDiagramCellDropdownMenu.style.display = "none";
        });
        crlDiagramCellDropdownMenu.addEventListener("mouseup", function () {
            crlDiagramCellDropdownMenu.style.display = "none";
        });
    };
    crlDiagramTabDropdownMenu = document.getElementById("diagramTabDropdown");
    if (crlDiagramTabDropdownMenu) {
        crlDiagramTabDropdownMenu.addEventListener("mouseleave", function () {
            crlDiagramTabDropdownMenu.style.display = "none";
        });
        crlDiagramTabDropdownMenu.addEventListener("mouseup", function () {
            crlDiagramTabDropdownMenu.style.display = "none";
        });
    };
    crlDebugSettingsDialog = new jBox("Confirm", {
        title: "Notification Trace Settings",
        confirmButton: "OK",
        cancelButton: "Cancel",
        content: "" +
            "<form>" +
            "	<fieldset>" +
            "		<label for='enableTracing'>Enable Notification Tracing</label>" +
            "		<input type='checkbox' id='enableTracing'> <br> " +
            "       <label>Omit Houskeeping Calls</label>" +
            "       <input type='checkbox' id='omitHousekeepingCalls'> <br>" +
            "       <label>Omit ManageTreeNodes Calls</label>" +
            "       <input type='checkbox' id='omitManageTreeNodesCalls'> <br>" +
            "       <label>Omit Diagram-related Calls</label>" +
            "       <input type='checkbox' id='omitDiagramRelatedCalls'>" +
            "	</fieldset> " +
            "</form>",
        confirm: function () {
            var enableNotificationTracing = "false";
            if ($("#enableTracing").prop("checked") == true) {
                enableNotificationTracing = "true";
            };
            var omitHousekeepingCalls = "false";
            if ($("#omitHousekeepingCalls").prop("checked") == true) {
                omitHousekeepingCalls = "true";
            }
            var omitManageTreeNodesCalls = "false";
            if ($("#omitManageTreeNodesCalls").prop("checked") == true) {
                omitManageTreeNodesCalls = "true";
            }
            var omitDiagramRelatedCalls = "false";
            if ($("#omitDiagramRelatedCalls").prop("checked") == true) {
                omitDiagramRelatedCalls = "true";
            }
            crlSendDebugSettings(enableNotificationTracing, omitHousekeepingCalls, omitManageTreeNodesCalls, omitDiagramRelatedCalls, "0");
        },
        onOpen: function () {
            $("#enableTracing").prop("checked", crlEnableTracing);
            $("#omitHousekeepingCalls").prop("checked", crlOmitHousekeepingCalls);
            $("#omitManageTreeNodesCalls").prop("checked", crlOmitManageTreeNodesCalls);
            $("#omitDiagramRelatedCalls").prop("checked", crlOmitDiagramRelatedCalls);
        }
    });
    crlDiagramElementFormatDialog = new jBox("Confirm", {
        title: "Diagram Element Format",
        confirmButton: "OK",
        cancelButton: "Cancel",
        content: "" +
            "<form>" +
            "   <fieldset>" +
            "       <label>Line Color:</label><input id='lineColor' type='color'><br>" +
            "       <label>Background Color:</label><input id='bgColor' type='color'>" +
            "   </fieldset>" +
            "</form>",
        confirm: function () {
            crlLineColor = $("#lineColor").val()
            crlBGColor = $("#bgColor").val()
            crlSendDiagramElementFormatChanged(crlSelectedDiagramElementIDForFormat, $("#lineColor").val(), $("#bgColor").val())
        }
    })
    crlDisplayCallGraphsDialog = new jBox("Confirm", {
        title: "Display Call Graphs",
        confirmButton: "OK",
        cancelButton: "Cancel",
        content: "" +
            "<form>" +
            "    <fieldset>" +
            "        <label>Number of available call graphs: </label> <label id='numberOfAvailableGraphs'></label><br>" +
            "        <label>Selected Graph:</label><input id='selectedGraph' type='number'>" +
            "	</fieldset> " +
            "</form>",
        confirm: function () {
            var selectedNumber = $("#selectedGraph").val();
            crlSendDisplayCallGraph(selectedNumber);
        }
    })
    crlSelectConceptByIDDialog = new jBox("Confirm", {
        title: "Enter ID of concept to be selected",
        confirmButton: "OK",
        cancelButton: "Cancel",
        content: "" +
            "<form>" +
            "	<fieldset>" +
            "		<label for='enteredConceptID'>ConceptID: </label>" +
            "		<input type='text' id='enteredConceptID'><br>" +
            "	</fieldset>" +
            "</form>",
        confirm: function () {
            var enteredConceptIDElement = document.getElementById("enteredConceptID");
            var xhr = crlCreateEmptyRequest();
            var data = JSON.stringify({ "Action": "ShowConceptInNavigator", "RequestConceptID": enteredConceptIDElement.value });
            crlSendRequest(xhr, data);
        }
    });
    crlUserPreferencesDialog = new jBox("Confirm", {
        title: "User Preferences",
        confirmButton: "OK",
        cancelButton: "Cancel",
        content: "" +
            "<form>" +
            "	<fieldset>" +
            "		<label for='dropReferenceAsLink'>Drop Reference As link</label>" +
            "		<input type='checkbox' id='dropReferenceAsLink'><br>" +
            "		<label for='dropRefinementAsLink'>Drop Refinement As link</label>" +
            "		<input type='checkbox' id='dropRefinementAsLink'>" +
            "	</fieldset>" +
            "</form>",
        confirm: function () {
            crlSendUserPreferences();
        },
        onOpen: function () {
            $("#dropReferenceAsLink").prop("checked", crlDropReferenceAsLink);
            $("#dropRefinementAsLink").prop("checked", crlDropRefinementAsLink);
        }
    });
    crlOpenWorkspaceDialog = new jBox("Confirm", {
        title: "Select Workspace",
        confirmButton: "OK",
        cancelButton: "Cancel",
        content: "" +
            "<form>" +
            "	<fieldset>" +
            "		<p>Use the file selector to locate the folder you want to use for your workspace. Copy the path in the" +
            "			top of the browser and then paste it into the indicated box.</p>" +
            "		Identify Workspace Folder:<input type='file' webkitdirectory directory multiple><br>" +
            "		Paste Directory Path Here:<input type='text' id='selectedWorkspaceFolder'>" +
            "	</fieldset>" +
            "</form>",
        confirm: function () {
            crlSendOpenWorkspace($("#selectedWorkspaceFolder").val());
        },
        onOpen: function () {
            $("#selectedWorkspaceFolder").val(crlWorkspacePath);
        }
    });
    window.onkeydown = function (e) {
        e = e || window.event;
        crlKeyPressed[e.keyCode] = true;
    };
    window.onkeyup = function (e) {
        e = e || window.event;
        crlKeyPressed[e.keyCode] = false;
    };
    window.onmousedown = function (e) {
        e = e || window.event;
        crlMouseButtonPressed[e.button] = true;
    }
    window.onmouseup = function (e) {
        e = e || window.event;
        crlMouseButtonPressed[e.button] = false;
    }
    window.onmousemove = function (e) {
        e = e || window.event;
        crlMousePosition["x"] = e.pageX;
        crlMousePosition["y"] = e.pageY;
    }
    // Close the dropdown menu if the user clicks outside of it
    window.onclick = function (event) {
        if (!event.target.matches('.dropbtn')) {

            var dropdowns = document.getElementsByClassName("dropdown-content");
            var i;
            for (i = 0; i < dropdowns.length; i++) {
                var openDropdown = dropdowns[i];
                if (openDropdown.classList.contains('show')) {
                    openDropdown.classList.remove('show');
                }
            }
        }
    }


});


function crlAppendToolbarButton(toolbar, id, icon) {
    var btn = document.createElement("BUTTON");
    btn.setAttribute("class", "toolbar-button");
    btn.setAttribute("id", id);
    btn.onclick = crlOnToolbarButtonSelected;
    var image = document.createElement("IMG");
    image.setAttribute("class", "toolbar-button-icon");
    image.setAttribute("src", icon);
    image.style.backgroundColor = "grey";
    btn.appendChild(image);
    toolbar.appendChild(btn);
}

function crlBringToFront(evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    cellView.model.toFront();
}

function crlEditFormat(evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    crlSelectedDiagramElementIDForFormat = crlGetConceptIDFromJointElementID(cellView.model.attributes.crlJointID)
    $("#lineColor").val(cellView.model.attributes["lineColor"])
    $("#bgColor").val(cellView.model.attributes["bgColor"])
    crlDiagramElementFormatDialog.open()
}

function crlCopyFormat(evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    crlLineColor = cellView.model.attributes["lineColor"]
    crlBGColor = cellView.model.attributes["bgColor"]
}

function crlPasteFormat(evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    crlSelectedDiagramElementIDForFormat = crlGetConceptIDFromJointElementID(cellView.model.attributes.crlJointID)
    crlSendDiagramElementFormatChanged(crlSelectedDiagramElementIDForFormat, crlLineColor, crlBGColor)
}

function crlCloseDiagramView(diagramID) {
    crlCloseDiagramViewWithoutNotification(diagramID)
    // Notify the server
    crlSendDiagramViewHasBeenClosed(diagramID)
}

function crlCloseDiagramViewWithoutNotification(diagramID) {
    var diagramContainerID = crlGetDiagramContainerIDFromDiagramID(diagramID);
    var diagramContainer = document.getElementById(diagramContainerID);
    if (crlCurrentDiagramContainerID == diagramContainerID) {
        crlCurrentDiagramContainerID = "";
    }
    // Finalize any node moves
    crlFinalizeNodeMoves()
    // Remove the graph
    var jointGraphID = crlGetJointGraphIDFromDiagramID(diagramID);
    delete crlGraphsGlobal[jointGraphID];
    // Remove the paper
    var jointPaperID = crlGetJointPaperIDFromDiagramID(diagramID);
    delete crlPapersGlobal[jointPaperID];
    // Delete the diagram container
    var topContent = document.getElementById("top-content");
    if (diagramContainer) {
        topContent.removeChild(diagramContainer);
    }
    // Delete the tab
    var tabs = document.getElementById("tabs");
    var tabID = crlGetDiagramTabIDFromDiagramID(diagramID);
    var tab = document.getElementById(tabID);
    if (tab) {
        tabs.removeChild(tab);
    }
}


function crlConstructDiagramContainer(diagramContainer, diagramContainerID, diagramLabel, diagramID) {
    var topContent = document.getElementById("top-content");
    diagramContainer = document.createElement("DIV");
    diagramContainer.id = diagramContainerID;
    diagramContainer.className = "crlDiagramContainer";
    diagramContainer.onclick = crlOnDiagramClick;
    // It is not clear why, but the ondrop callback does not get called unless the ondragover callback is used,
    // even though the callback just calls preventDefault on the dragover event
    diagramContainer.ondragover = crlOnDragover;
    diagramContainer.onmouseover = crlOnDiagramMouseOver;
    diagramContainer.ondrop = crlOnDiagramDrop;
    diagramContainer.style.display = "none";
    topContent.appendChild(diagramContainer);
    // Create the new tab
    var tabs = document.getElementById("tabs");
    var newTab = document.createElement("BUTTON");
    newTab.innerHTML = diagramLabel;
    newTab.className = "w3-bar-item w3-button";
    var newTabID = crlGetDiagramTabIDFromDiagramID(diagramID);
    newTab.id = newTabID;
    newTab.setAttribute("diagramContainerID", diagramContainerID);
    newTab.onclick = crlOnMakeDiagramVisible;
    tabs.appendChild(newTab, -1);
    var jointGraphID = crlGetJointGraphIDFromDiagramID(diagramID);
    var jointGraph = crlGraphsGlobal[jointGraphID];
    //        jointGraph = document.getElementById(jointGraphID);
    if (jointGraph == undefined) {
        jointGraph = new joint.dia.Graph;
        jointGraph.id = jointGraphID;
        crlGraphsGlobal[jointGraphID] = jointGraph;
    }
    ;
    var jointPaperID = crlGetJointPaperIDFromDiagramID(diagramID);
    var jointPaper = crlPapersGlobal[jointPaperID];
    if (jointPaper == undefined) {
        jointPaper = crlConstructPaper(diagramContainer, jointGraph, jointPaperID);
    }
    ;
    return diagramContainer;
}

function crlConstructDiagramLink(data, graph, crlJointID) {
    var sourceJointID = crlGetJointCellIDFromConceptID(data.AdditionalParameters["LinkSourceID"])
    var targetJointID = crlGetJointCellIDFromConceptID(data.AdditionalParameters["LinkTargetID"])
    if (sourceJointID != "" && targetJointID != "") {
        var linkSource = crlFindCellInGraph(graph, sourceJointID)
        var linkTarget = crlFindCellInGraph(graph, targetJointID)
        if (linkSource != undefined && linkTarget != undefined) {
            var newLink;
            switch (data.AdditionalParameters["LinkType"]) {
                case "ReferenceLink":
                    newLink = new joint.shapes.crl.ReferenceLink({
                        source: linkSource,
                        target: linkTarget
                    });
                    break;
                case "RefinementLink":
                    newLink = new joint.shapes.crl.RefinementLink({
                        source: linkSource,
                        target: linkTarget
                    });
                    break;
                case "OwnerPointer":
                    newLink = new joint.shapes.crl.OwnerPointer({
                        source: linkSource,
                        target: linkTarget
                    });
                    break;
                case "ElementPointer":
                    newLink = new joint.shapes.crl.ElementPointer({
                        source: linkSource,
                        target: linkTarget
                    });
                    break;
                case "AbstractPointer":
                    newLink = new joint.shapes.crl.AbstractPointer({
                        source: linkSource,
                        target: linkTarget
                    });
                    break;
                case "RefinedPointer":
                    newLink = new joint.shapes.crl.RefinedPointer({
                        source: linkSource,
                        target: linkTarget
                    });
                    break;
            }
            newLink.set("crlJointID", crlJointID);
            newLink.set("represents", data.AdditionalParameters["Represents"]);
            newLink.set("dummyEndChangeToggle", false);
            graph.addCell(newLink);

            return newLink;
        }
    }
    return undefined
}

function crlConstructDiagramNode(data, graph, crlJointID) {
    var jointElement = new joint.shapes.crl.Element({});
    jointElement.set("crlJointID", crlJointID);
    jointElement.set("name", data.AdditionalParameters["DisplayLabel"]);
    jointElement.set("position", { "x": Number(data.AdditionalParameters["NodeX"]), "y": Number(data.AdditionalParameters["NodeY"]) });
    jointElement.set("size", { "width": Number(data.AdditionalParameters["NodeWidth"]), "height": Number(data.AdditionalParameters["NodeHeight"]) });
    jointElement.set("icon", data.AdditionalParameters["Icon"]);
    jointElement.set("abstractions", data.AdditionalParameters["Abstractions"]);
    jointElement.set("represents", data.AdditionalParameters["Represents"])
    graph.addCell(jointElement);
    return jointElement;
}

function crlConstructPaper(diagramContainer, jointGraph, jointPaperID) {
    var diagramPaperDiv = document.createElement("DIV");
    diagramContainer.appendChild(diagramPaperDiv);
    jointPaper = new joint.dia.Paper({
        "el": diagramPaperDiv,
        "width": 2000,
        "height": 2000,
        defaultLink: undefined,
        linkView: crlCustomLinkView,
        validateMagnet: crlValidateLinkStart,
        validateConnection: crlValidateConnection,
        "model": jointGraph,
        "gridSize": 1
    });
    jointPaper.options.connectionStrategy = joint.connectionStrategies.centerPort;
    // ConnectionPoint
    var linkConnectionPoint = function (linkView, view, magnet, reference) {
        var model = view.model;
        var spot;
        if (model.isElement()) {
            var bbox = model.getBBox();
            spot = bbox.intersectionWithLineFromCenterToPoint(reference);
        } else if (model.isLink()) {
            var label = model.labels()[0];
            spot = view.getLabelCoordinates(0.5);
        }
        return spot || model.getBBox();
    };
    jointPaper.options.linkConnectionPoint = linkConnectionPoint;
    // Event handlers
    jointPaper.on("cell:pointerdown", crlOnDiagramCellPointerDown);
    jointPaper.on("cell:pointerup", crlOnDiagramCellPointerUp);
    jointPaper.on("blank:pointerup", function (evt) {
        crlMouseButtonPressed[evt.button] = false;
    })
    jointPaper.on("cell:contextmenu", function (cellView, evt, x, y) {
        evt.preventDefault();
        var represents = cellView.model.attributes.represents;
        if (represents == "Reference") {
            document.getElementById("showReferencedConcept").style.display = "block";
            document.getElementById("nullifyReferencedConcept").style.display = "block";
        } else {
            document.getElementById("showReferencedConcept").style.display = "none";
            document.getElementById("nullifyReferencedConcept").style.display = "none";
        }
        if (represents == "OwnerPointer" || represents == "ElementPointer" || represents == "AbstractPointer" || represents == "RefinedPointer") {
            document.getElementById("showOwner").style.display = "none";
        } else {
            document.getElementById("showOwner").style.display = "block";
        }
        if (represents == "Refinement") {
            document.getElementById("showAbstractConcept").style.display = "block";
            document.getElementById("showRefinedConcept").style.display = "block";
        } else {
            document.getElementById("showAbstractConcept").style.display = "none";
            document.getElementById("showRefinedConcept").style.display = "none";
        }
        crlDiagramCellDropdownMenu.attributes.cellView = cellView;
        crlDiagramCellDropdownMenu.style.left = evt.pageX.toString() + "px";
        crlDiagramCellDropdownMenu.style.top = evt.pageY.toString() + "px";
        crlDiagramCellDropdownMenu.style.display = "block";
    });
    jointPaper.on("link:connect", crlLinkConnected);
    jointPaper.on("validateMagnet", crlValidateLinkStart);
    jointPaper.on('link:mouseenter', function (linkView) {
        var toolsView = linkView._toolsView;
        if (!toolsView) {
            var verticesTool = new joint.linkTools.Vertices();
            var segmentsTool = new joint.linkTools.Segments();
            var targetArrowheadTool = new joint.linkTools.TargetArrowhead();
            var sourceArrowheadTool = null;
            if (!crlLinkViewRepresentsPointer(linkView)) {
                sourceArrowheadTool = new joint.linkTools.SourceArrowhead();
            }
            toolsView = new joint.dia.ToolsView({
                tools: [verticesTool, segmentsTool, targetArrowheadTool, sourceArrowheadTool]
            });
            linkView.addTools(toolsView);
        }
        if (crlKeyPressed[16]) {
            linkView.showTools();
        } else {
            linkView.hideTools();
        }
    });
    jointPaper.on('link:mouseleave', function (linkView) {
        linkView.hideTools();
    });
    crlPapersGlobal[jointPaperID] = jointPaper;
    return jointPaper;
}

function crlCallExit() {
    console.log("Requesting Exit");
    var xhr = crlCreateEmptyRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            window.close();
        }
    };
    var data = JSON.stringify({ "Action": "Exit" });
    crlSendRequest(xhr, data);
}

function crlClearDiagrams() {
    var x = document.getElementsByClassName("crlDiagramContainer");
    // Note that the number of elements in x decreases as the diagrams are removed
    for (i = 0; x.length > 0;) {
        var container = x.item(i);
        var diagramID = crlGetDiagramIDFromDiagramContainerID(container.id);
        crlCloseDiagramViewWithoutNotification(diagramID);
    };
    crlSendNormalResponse();
}

var crlCloseWebsocket = function () {
    console.log("Closing websocket")
    crlWebsocketGlobal.close()
}

function crlCreateEmptyRequest() {
    var xhr = new XMLHttpRequest();
    var url = "request";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var response = JSON.parse(xhr.responseText)
            if (response.Result == 1) {
                // suppress alerts if automated regression testing is in progress
                if (crlAutomatedTestInProgress == false) {
                    alert(response.ResultDescription);
                }
            }
            console.log(response)
        };
    }
    return xhr
}

function crlCreateLink() {
    var link;
    switch (crlCurrentToolbarButton) {
        case crlCursorToolbarButtonID:
            link = undefined;
            break;
        case crlElementToolbarButtonID:
            link = undefined;
            break;
        case crlLiteralToolbarButtonID:
            link = undefined;
            break;
        case crlReferenceToolbarButtonID:
            link = undefined;
            break;
        case crlReferenceLinkToolbarButtonID:
            link = new joint.shapes.crl.ReferenceLink({});
            break;
        case crlRefinementToolbarButtonID:
            link = undefined;
            break;
        case crlRefinementLinkToolbarButtonID:
            link = new joint.shapes.crl.RefinementLink({});
            break;
        case crlDiagramToolbarButtonID:
            link = undefined;
            break;
        case crlOwnerPointerToolbarButtonID:
            link = new joint.shapes.crl.OwnerPointer({});
            break;
        case crlElementPointerToolbarButtonID:
            link = new joint.shapes.crl.ElementPointer({});
            break;
        case crlAbstractPointerToolbarButtonID:
            link = new joint.shapes.crl.AbstractPointer({});
            break;
        case crlRefinedPointerToolbarButtonID:
            link = new joint.shapes.crl.RefinedPointer({});
            break;
    }

    return link;
}

var crlCustomLinkView = joint.dia.LinkView.extend({
    onmagnet: function (evt, x, y) {
        this.dragMagnetStart(evt, x, y);
    },
    dragMagnetStart: function (evt, x, y) {

        if (!this.can('addLinkFromMagnet')) return;

        var magnet = evt.currentTarget;
        var paper = this.paper;
        this.eventData(evt, { targetMagnet: magnet });
        evt.stopPropagation();

        if (paper.options.validateMagnet(this, magnet)) {

            if (paper.options.magnetThreshold <= 0) {
                this.dragLinkStart(evt, magnet, x, y);
            }

            this.eventData(evt, { action: 'magnet' });
            this.stopPropagation(evt);

        } else {

            this.pointerdown(evt, x, y);
        }

        paper.delegateDragEvents(this, evt.data);
    },
    dragLinkStart: function (evt, magnet, x, y) {

        this.model.startBatch('add-link');

        var linkView = this.addLinkFromMagnet(magnet, x, y);

        // backwards compatiblity events
        joint.dia.CellView.prototype.pointerdown.apply(linkView, arguments);
        linkView.notify('link:pointerdown', evt, x, y);

        linkView.eventData(evt, linkView.startArrowheadMove('target', { whenNotAllowed: 'remove' }));
        this.eventData(evt, { linkView: linkView });
    },

    addLinkFromMagnet: function (magnet, x, y) {

        var paper = this.paper;
        var graph = paper.model;

        var link = paper.getDefaultLink(this, magnet);
        link.set({
            source: this.getLinkEnd(magnet, x, y, link, 'source'),
            target: { x: x, y: y }
        }).addTo(graph, {
            async: false,
            ui: true
        });

        return link.findView(paper);
    },
    getBBox: function (opt) {

        var bbox;
        if (opt && opt.useModelGeometry) {
            var model = this.model;
            bbox = model.getBBox().bbox(model.angle());
        } else {
            bbox = this.getNodeBBox(this.el);
        }

        return this.paper.localToPaperRect(bbox);
    },
    getNodeBBox: function (magnet) {

        var rect = this.getNodeBoundingRect(magnet);
        var magnetMatrix = this.getNodeMatrix(magnet);
        // var translateMatrix = this.getRootTranslateMatrix();
        // var rotateMatrix = this.getRootRotateMatrix();
        // return V.transformRect(rect, translateMatrix.multiply(rotateMatrix).multiply(magnetMatrix));
        return V.transformRect(rect, magnetMatrix);
    },
    getNodeBoundingRect: function (magnet) {

        var metrics = this.nodeCache(magnet);
        if (metrics.boundingRect === undefined) metrics.boundingRect = V(magnet).getBBox();
        return new g.Rect(metrics.boundingRect);
    },
    getBBox: function (opt) {

        var bbox;
        if (opt && opt.useModelGeometry) {
            var model = this.model;
            bbox = model.getBBox().bbox(model.angle());
        } else {
            bbox = this.getNodeBBox(this.el);
        }

        return this.paper.localToPaperRect(bbox);
    },
    nodeCache: function (magnet) {

        var metrics = this.metrics;
        if (!metrics) {
            // don't use cache
            // it most likely a custom view with overridden update
            return {};
        }

        var id = V.ensureId(magnet);

        var value = metrics[id];
        if (!value) value = metrics[id] = {};
        return value;
    },
    getNodeMatrix: function (magnet) {

        var metrics = this.nodeCache(magnet);
        if (metrics.magnetMatrix === undefined) {
            var target = this.rotatableNode || this.el;
            metrics.magnetMatrix = V(magnet).getTransformToElement(target);
        }
        return V.createSVGMatrix(metrics.magnetMatrix);
    },
    getRootTranslateMatrix: function () {

        var model = this.model;
        var position = model.position();
        var mt = V.createSVGMatrix().translate(position.x, position.y);
        return mt;
    },
    dragMagnet: function (evt, x, y) {

        var data = this.eventData(evt);
        var linkView = data.linkView;
        if (linkView) {
            linkView.pointermove(evt, x, y);
        } else {
            var paper = this.paper;
            var magnetThreshold = paper.options.magnetThreshold;
            var currentTarget = this.getEventTarget(evt);
            var targetMagnet = data.targetMagnet;
            if (magnetThreshold === 'onleave') {
                // magnetThreshold when the pointer leaves the magnet
                if (targetMagnet === currentTarget || V(targetMagnet).contains(currentTarget)) return;
            } else {
                // magnetThreshold defined as a number of movements
                if (paper.eventData(evt).mousemoved <= magnetThreshold) return;
            }
            this.dragLinkStart(evt, targetMagnet, x, y);
        }
    },
    pointermove: function (evt, x, y) {

        // Backwards compatibility
        var dragData = this._dragData;
        if (dragData) this.eventData(evt, dragData);

        var data = this.eventData(evt);
        switch (data.action) {
            case 'magnet':
                this.dragMagnet(evt, x, y);
                break;
            case 'vertex-move':
                this.dragVertex(evt, x, y);
                break;

            case 'label-move':
                this.dragLabel(evt, x, y);
                break;

            case 'arrowhead-move':
                this.dragArrowhead(evt, x, y);
                break;

            case 'move':
                this.drag(evt, x, y);
                break;
        }
        // Backwards compatibility
        if (dragData) joint.util.assign(dragData, this.eventData(evt));

        joint.dia.CellView.prototype.pointermove.apply(this, arguments);
        this.notify('link:pointermove', evt, x, y);
    },
    pointerup: function (evt, x, y) {

        // Backwards compatibility
        var dragData = this._dragData;
        if (dragData) {
            this.eventData(evt, dragData);
            this._dragData = null;
        }

        var data = this.eventData(evt);
        switch (data.action) {
            case 'magnet':
                this.dragMagnetEnd(evt, x, y);
                break;
            case 'vertex-move':
                this.dragVertexEnd(evt, x, y);
                break;

            case 'label-move':
                this.dragLabelEnd(evt, x, y);
                break;

            case 'arrowhead-move':
                this.dragArrowheadEnd(evt, x, y);
                break;

            case 'move':
                this.dragEnd(evt, x, y);
        }

        var magnet = data.targetMagnet;
        if (magnet) this.magnetpointerclick(evt, magnet, x, y);


        this.notify('link:pointerup', evt, x, y);
        joint.dia.CellView.prototype.pointerup.apply(this, arguments);
    },
    dragMagnetEnd: function (evt, x, y) {

        var data = this.eventData(evt);
        var linkView = data.linkView;
        if (!linkView) return;
        linkView.pointerup(evt, x, y);
        this.model.stopBatch('add-link');
    },
    magnetpointerclick: function (evt, magnet, x, y) {
        var paper = this.paper;
        if (paper.eventData(evt).mousemoved > paper.options.clickThreshold) return;
        this.notify('element:magnet:pointerclick', evt, magnet, x, y);
    },
    startListening: function () {
        // Code from joint.js version 2.2.1
        var model = this.model;

        this.listenTo(model, 'change:markup', this.render);
        this.listenTo(model, 'change:smooth change:manhattan change:router change:connector', this.update);
        this.listenTo(model, 'change:toolMarkup', this.onToolsChange);
        this.listenTo(model, 'change:labels change:labelMarkup', this.onLabelsChange);
        this.listenTo(model, 'change:vertices change:vertexMarkup', this.onVerticesChange);
        this.listenTo(model, 'change:source', this.onSourceChange);
        this.listenTo(model, 'change:target', this.onTargetChange);
    },
    // The custom update function adds a hack to ensure that changes to links that are endpoints of other links trigger updates of those other links
    // Default is to process the `attrs` object and set attributes on subelements based on the selectors.
    update: function (model, attributes, opt) {

        // Change the value of dummyEndChangeToggle. The change to this attribute is enough to notify links for which this is an endpoint that a change has 
        // occurred. The specific case being addressed is the movement of an Element (node) that is an endpoint of this link. Such movement does not actually
        // cause a change to the link model (without this hack): only the link view is changed.
        this.model.set("dummyEndChangeToggle", !this.model.get("dummyEndChangeToggle"));

        opt || (opt = {});


        // update the link path
        this.updateConnection(opt);

        // update SVG attributes defined by 'attrs/'.
        this.updateDOMSubtreeAttributes(this.el, this.model.attr(), { selectors: this.selectors });

        this.updateDefaultConnectionPath();

        // update the label position etc.
        this.updateLabelPositions();
        this.updateToolsPosition();
        this.updateArrowheadMarkers();

        this.updateTools(opt);
        // Local perpendicular flag (as opposed to one defined on paper).
        // Could be enabled inside a connector/router. It's valid only
        // during the update execution.
        this.options.perpendicular = null;
        // Mark that postponed update has been already executed.
        this.updatePostponed = false;

        return this;
    },

});

function crlDeleteDiagramElementView(evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    if (jointID) {
        var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
        var xhr = crlCreateEmptyRequest();
        var data = JSON.stringify({ "Action": "DeleteDiagramElementView", "RequestConceptID": diagramElementID });
        crlSendRequest(xhr, data);
    }
}

function crlPropertiesClearRow(row) {
    var properties = document.getElementById("properties");
    var propertyRow = properties.rows[row]
    if (propertyRow != undefined) {
        properties.deleteRow(row);
    }
}

function crlPropertiesDisplayAbstractConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Abstract Concept ID";
    var abstractConceptID = ""
    if (data.NotificationConceptState) {
        abstractConceptID = data.NotificationConceptState.AbstractConceptID
    }
    typeRow.cells[1].innerHTML = abstractConceptID;
}

function crlPropertiesDisplayDefinition(data, row) {
    var definitionRow = crlObtainPropertyRow(row)
    definitionRow.cells[0].innerHTML = "Definition";
    var definition = ""
    var isCore = false
    var isReadOnly = false
    if (data.NotificationConceptState) {
        definition = data.NotificationConceptState.Definition
        isCore = data.NotificationConceptState.IsCore
        isReadOnly = data.NotificationConceptState.ReadOnly
    }
    var input = definitionRow.cells[1].firstElementChild;
    var cursorPosition = input.selectionStart;
    input.value = definition;
    input.id = "definition";
    if (isCore == "false" && isReadOnly == "false") {
        input.read_only = false;
        if (!definitionRow.cells[1].callbackAssigned) {
            input.callbackAssigned = true;
            $("#definition").on("keyup", crlSendDefinitionChanged);
        }
    } else {
        input.read_only = true;
    };
    input.setSelectionRange(cursorPosition, cursorPosition);
}

function crlPropertiesDisplayID(data, row) {
    var idRow = crlObtainPropertyRow(row)
    idRow.cells[0].innerHTML = "ID";
    idRow.cells[1].innerHTML = data.NotificationConceptID;
}

function crlPropertiesDisplayLabel(data, row) {
    var labelRow = crlObtainPropertyRow(row);
    labelRow.cells[0].innerHTML = "Label";
    var label = "";
    var isCore = false;
    var isReadOnly = false;
    if (data.NotificationConceptState) {
        label = data.NotificationConceptState.Label;
        isCore = data.NotificationConceptState.IsCore;
        isReadOnly = data.NotificationConceptState.ReadOnly;
    }
    var input = labelRow.cells[1].firstElementChild;
    var cursorPosition = input.selectionStart;
    input.value = label;
    input.id = "elementLabel";
    if (isCore == "false" && isReadOnly == "false") {
        input.read_only = false;
        if (!input.callbackAssigned) {
            input.callbackAssigned = true;
            $("#elementLabel").on("keyup", crlSendLabelChanged);
        }
    } else {
        input.read_only = true;
    };
    input.setSelectionRange(cursorPosition, cursorPosition);
}

function crlPropertiesDisplayLiteralValue(data, row) {
    var literalValueRow = crlObtainPropertyRow(row);
    literalValueRow.cells[0].innerHTML = "Literal Value";
    var literalValue = ""
    var isCore = false
    var isReadOnly = false
    if (data.NotificationConceptState) {
        literalValue = data.NotificationConceptState.LiteralValue
        isCore = data.NotificationConceptState.IsCore
        isReadOnly = data.NotificationConceptState.ReadOnly
    }
    var input = literalValueRow.cells[1].firstElementChild;
    if (input == null) {
        // The existing row does not have an input: add one
        literalValueRow.cells[1].innerHTML = "";
        input = document.createElement("input");
        input.setAttribute("type", "text");
        literalValueRow.cells[1].appendChild(input);
    }
    var cursorPosition = input.selectionStart;
    input.value = literalValue;
    input.id = "literalValue";
    if (isCore == "false" && isReadOnly == "false") {
        input.read_only = false;
        if (!literalValueRow.cells[1].callbackAssigned) {
            input.callbackAssigned = true;
            $("#literalValue").on("keyup", crlSendLiteralValueChanged);
        }
    } else {
        input.read_only = true;
    };
    input.setSelectionRange(cursorPosition, cursorPosition);
}

function crlPropertiesDisplayOwningConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Owning Concept ID";
    var owningConceptID = ""
    if (data.NotificationConceptState) {
        owningConceptID = data.NotificationConceptState.OwningConceptID
    }
    typeRow.cells[1].innerHTML = owningConceptID;
}

function crlPropertiesDisplayReadOnly(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Read Only";
    var readOnly = ""
    if (data.NotificationConceptState) {
        readOnly = data.NotificationConceptState.ReadOnly
    }
    typeRow.cells[1].innerHTML = readOnly;
}

function crlPropertiesDisplayReferencedAttributeName(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Referenced AttributeName";
    var referencedAttributeName = ""
    if (data.NotificationConceptState) {
        referencedAttributeName = data.NotificationConceptState.ReferencedAttributeName
    }
    typeRow.cells[1].innerHTML = referencedAttributeName;
}

function crlPropertiesDisplayReferencedConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Referenced Concept ID";
    var referencedConceptID = ""
    if (data.NotificationConceptState) {
        referencedConceptID = data.NotificationConceptState.ReferencedConceptID
    }
    typeRow.cells[1].innerHTML = referencedConceptID;
}

function crlPropertiesDisplayRefinedConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Refined Concept ID";
    var refinedConceptID = ""
    if (data.NotificationConceptState) {
        refinedConceptID = data.NotificationConceptState.RefinedConceptID
    }
    typeRow.cells[1].innerHTML = refinedConceptID;
}

function crlPropertiesDisplayType(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Type";
    var type = ""
    if (data.NotificationConceptState) {
        type = data.NotificationConceptState.ConceptType
    }
    typeRow.cells[1].innerHTML = type;
}

function crlPropertiesDisplayURI(data, row) {
    var uriRow = crlObtainPropertyRow(row);
    uriRow.cells[0].innerHTML = "URI";
    var uri = ""
    var isCore = false
    var isReadOnly = false
    if (data.NotificationConceptState) {
        uri = data.NotificationConceptState.URI
        isCore = data.NotificationConceptState.IsCore
        isReadOnly = data.NotificationConceptState.ReadOnly
    }
    var input = uriRow.cells[1].firstElementChild;
    var cursorPosition = input.selectionStart;
    input.value = uri;
    input.id = "uri";
    if (isCore == "false" && isReadOnly == "false") {
        input.read_only = false;
        if (!uriRow.cells[1].callbackAssigned) {
            input.callbackAssigned = true;
            $("#uri").on("keyup", crlSendURIChanged);
        }
    } else {
        input.read_only = true;
    }
    input.setSelectionRange(cursorPosition, cursorPosition);
}


function crlPropertiesDisplayVersion(data, row) {
    var versionRow = crlObtainPropertyRow(row)
    versionRow.cells[0].innerHTML = "Version";
    var version = ""
    if (data.NotificationConceptState) {
        version = data.NotificationConceptState.Version
    }
    versionRow.cells[1].innerHTML = version;
}

function crlDropdownMenu(dropdownId) {
    document.getElementById(dropdownId).classList.toggle("show");
}

function crlDisplayDiagram(diagramContainerID) {
    var x = document.getElementsByClassName("crlDiagramContainer");
    for (i = 0; i < x.length; i++) {
        var container = x.item(i);
        var diagramID = crlGetDiagramIDFromDiagramContainerID(container.id);
        var tabID = crlGetDiagramTabIDFromDiagramID(diagramID);
        var tab = document.getElementById(tabID);
        if (container.id == diagramContainerID) {
            container.style.display = "block";
            tab.style.backgroundColor = "white";
            tab.oncontextmenu = function (evt) {
                evt.preventDefault();
                crlDiagramTabDropdownMenu.setAttribute("diagramID", diagramID);
                crlDiagramTabDropdownMenu.style.left = evt.pageX.toString() + "px";
                crlDiagramTabDropdownMenu.style.top = evt.pageY.toString() + "px";
                crlDiagramTabDropdownMenu.style.display = "block";
            };
            var graphID = crlGetJointGraphIDFromDiagramID(diagramID);
            var graph = crlGraphsGlobal[graphID]
            if (graph) {
                graph.resetCells(graph.getCells());
            }
        } else {
            container.style.display = "none";
            tab.style.backgroundColor = "grey";
        }
    }
}

function crlDisplayGraphDialog(numberOfAvailableGraphs) {
    crlDisplayCallGraphsDialog.open();
    $("#numberOfAvailableGraphs").text(numberOfAvailableGraphs);
}

function crlFindCellInGraphID(graphID, crlJointID) {
    var graph = crlGraphsGlobal[graphID];
    return crlFindCellInGraph(graph, crlJointID);
}

function crlFindCellInGraph(graph, crlJointID) {
    var cells = graph.getCells();
    var cell = null;
    cells.forEach(function (item) {
        if (item.get("crlJointID") == crlJointID) {
            cell = item;
        }
    })
    return cell;
}

function crlFindCellViewInPaperByDiagramID(diagramID, jointCellID) {
    var jointPaperID = crlGetJointPaperIDFromDiagramID(diagramID);
    var paper = crlPapersGlobal[jointPaperID];
    return crlFindCellViewInPaper(paper, jointCellID);
}

function crlFindCellViewInPaper(paper, jointCellID) {
    return paper.findViewByModel(jointCellID)
}

function crlFindLinkInGraph(graphID, crlJointID) {
    var links = crlGraphsGlobal[graphID].getLinks();
    var link = null;
    links.forEach(function (item) {
        if (item.get("crlJointID") == crlJointID) {
            link = item;
        }
    })
    return link
}

function crlGetConceptIDFromContainerID(containerID) {
    return containerID.replace("DiagramContainer", "")
}

function crlGetConceptIDFromTreeNodeID(treeNodeID) {
    return treeNodeID.replace("TreeNode", "");
}

function crlGetContainerIDFromConceptID(conceptID) {
    return "DiagramContainer" + conceptID;
}

function crlGetDiagramContainerIDFromDiagramID(diagramID) {
    return "DiagramContainer" + diagramID;
}

function crlGetDiagramIDFromDiagramContainerID(diagramContainerID) {
    return diagramContainerID.replace("DiagramContainer", "");
}

function crlGetDiagramIDFromJointGraphID(jointGraphID) {
    return jointGraphID.replace("JointGraph", "");
}

function crlGetDiagramIDFromJointPaperID(jointPaperID) {
    return jointPaperID.replace("JointPaper", "")
}

function crlGetDiagramTabIDFromDiagramID(diagramID) {
    return "DiagramTab" + diagramID;
}

function crlGetJointPaperIDFromDiagramID(diagramID) {
    return "JointPaper" + diagramID;
}

function crlGetJointGraphIDFromDiagramID(diagramID) {
    return "JointGraph" + diagramID;
}

function crlGetTreeNodeIDFromConceptID(conceptID) {
    return "TreeNode" + conceptID;
}


function crlGetConceptIDFromJointElementID(jointElementID) {
    return jointElementID.replace("JointElement", "")
}

function crlGetJointCellIDFromConceptID(conceptID) {
    return "JointElement" + conceptID
}

function crlIgnoreSpecialCharacters(inputString) {
    return inputString
        .replace(/[\\]/g, '')
        .replace(/[\/]/g, '')
        .replace(/[\b]/g, '')
        .replace(/[\f]/g, '')
        .replace(/[\n]/g, '')
        .replace(/[\r]/g, '')
        .replace(/[\t]/g, '')
        .replace(/[\"]/g, '');
}

function crlInitializeClient() {
    crlInitializeWebSocket();
    // console.log("Requesting InitializeClient");
    // var xhr = crlCreateEmptyRequest();
    // var data = JSON.stringify({ "Action": "InitializeClient" });
    // crlSendRequest(xhr, data);
}

function crlInitializeTree() {
    $("#uOfD").jstree({
        'core': {
            'check_callback': true,
            'multiple': false
        },
        'plugins': ['sort', 'contextmenu'],
        'sort': function (a, b) {
            aNode = this.get_node(a);
            bNode = this.get_node(b);
            var aNodeText = aNode.text
            var bNodeText = bNode.text
            if (aNodeText == bNodeText) {
                return aNode.id > bNode.id ? 1 : -1;
            }
            return aNodeText > bNodeText ? 1 : -1;
        },
        'contextmenu': {
            "items": function ($node) {
                var tree = $("uOfD").jstree(true);
                var items = {
                    addChild: {
                        "label": "Add Child",
                        "action": false,
                        "submenu": {
                            Element: {
                                "label": "Element",
                                "action": function (obj) {
                                    if ($node != undefined) {
                                        var conceptID = crlGetConceptIDFromTreeNodeID($node.id);
                                        crlSendAddElementChild(conceptID);
                                    }
                                }
                            },
                            Diagram: {
                                "label": "Diagram",
                                "action": function (obj) {
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id);
                                    crlSendAddDiagramChild(conceptID);
                                }
                            },
                            Literal: {
                                "label": "Literal",
                                "action": function (obj) {
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id);
                                    crlSendAddLiteralChild(conceptID);
                                }
                            },
                            Reference: {
                                "label": "Reference",
                                "action": function (obj) {
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id);
                                    crlSendAddReferenceChild(conceptID);
                                }
                            },
                            Refinement: {
                                "label": "Refinement",
                                "action": function (obj) {
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id);
                                    crlSendAddRefinementChild(conceptID);
                                }
                            }
                        }
                    },
                    display: {
                        "label": "Display Diagram",
                        "action": function (obj) {
                            if ($node != undefined) {
                                var conceptID = crlGetConceptIDFromTreeNodeID($node.id);
                                crlSendDisplayDiagramSelected(conceptID);
                            }
                        }
                    },
                    remove: {
                        "label": "Delete",
                        "action": function (obj) {
                            if ($node != undefined) {
                                var conceptID = crlGetConceptIDFromTreeNodeID($node.id);
                                crlSendTreeNodeDelete(conceptID);
                            };
                        }
                    }
                }
                if ($node.li_attr.is_diagram == "false") {
                    delete items.display
                }
                if ($node.li_attr.is_core == "true" || $node.li_attr.read_only == "true") {
                    delete items.remove
                }
                return items
            }
        }
    });
    $("#uOfD").on("select_node.jstree", crlSendTreeNodeSelected);
    $("#uOfD").on("dragstart", crlOnTreeDragStart);
}

function crlInitializeWebSocket() {
    // This function must be idempotent
    if (crlWebsocketGlobal == null) {
        console.log("Initializing Web Socket")
        // ws initialization
        crlWebsocketGlobal = new WebSocket("ws://localhost:8081/index/ws");
        console.log("Web Socket Initialization complete")
        crlWebsocketGlobal.onmessage = function (e) {
            var data = JSON.parse(e.data)
            console.log(data)
            switch (data.Notification) {
                case 'AddDiagramLink':
                    crlNotificationAddDiagramLink(data);
                    break;
                case 'AddDiagramNode':
                    crlNotificationAddDiagramNode(data);
                    break;
                case 'AddTreeNode':
                    crlNotificationAddTreeNode(data);
                    break;
                case "ChangeTreeNode":
                    crlNotificationUpdateTreeNode(data);
                    break;
                case "ClearDiagrams":
                    crlClearDiagrams();
                    break;
                case "ClearToolbarSelection":
                    crlNotificationClearToolbarSelection(data);
                    break;
                case "ClearTree":
                    crlNotificationClearTree();
                    break;
                case "CloseDiagramView":
                    crlNotificationCloseDiagramView(data);
                    break;
                case "DebugSettings":
                    crlNotificationSaveDebugSettings(data);
                    break;
                case "DeleteDiagramElement":
                    crlNotificationDeleteDiagramCell(data);
                    break;
                case "DeleteTreeNode":
                    crlNotificationDeleteTreeNode(data);
                    break;
                case "DiagramLabelChanged":
                    crlNotificationDiagramLabelChanged(data);
                    break;
                case "DisplayDiagram":
                    crlNotificationDisplayDiagram(data);
                    break;
                case "DisplayGraph":
                    crlNotificationDisplayGraph(data);
                    break;
                case "DoesLinkExist":
                    crlNotificationDoesLinkExist(data);
                    break;
                case "UserPreferences":
                    crlNotificationSaveUserPreferences(data);
                    break;
                case "ElementSelected":
                    crlNotificationElementSelected(data);
                    break;
                case "InitializationComplete":
                    crlInitializationComplete = true;
                    console.log("Initialization Complete")
                    crlSendNormalResponse("Processed InitializationComplete")
                    break;
                case "Refresh":
                    crlNotificationRefresh();
                    break;
                case "ShowTreeNode":
                    crlNotificationShowTreeNode(data);
                    break;
                case "UpdateDiagramLink":
                    crlNotificationUpdateDiagramLink(data);
                    break;
                case "UpdateDiagramNode":
                    crlNotificationUpdateDiagramNode(data);
                    break;
                case "UpdateProperties":
                    crlNotificationUpdateProperties(data);
                    break;
                case "WorkspacePath":
                    crlNotificationUpdateWorkspacePath(data);
                    break;
                default:
                    console.log('Unhandled notification: ' + e.data);
                    var data = {};
                    data["Result"] = 1;
                    data["ErrorMessage"] = "Unhandled notification";
                    crlWebsocketGlobal.send(JSON.stringify(data));
            }
        };
    } else {
        console.log("Web Socket already initialized")
    };
};

var crlInitiateGraphsDialogDisplay = function () {
    var xhr = new XMLHttpRequest();
    var url = "request";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var response = JSON.parse(xhr.responseText)
            console.log(response)
            if (response.Result == 1) {
                alert(response.ResultDescription);
            } else {
                crlDisplayGraphDialog(response.AdditionalParameters["NumberOfAvailableGraphs"]);
            }
        };
    }
    var data = JSON.stringify({
        "Action": "ReturnAvailableGraphCount"
    });
    crlSendRequest(xhr, data);
}

function crlLinkConnected(evt, cellView, magnet, arrowhead) {
    var linkType = evt.model.attributes.type;
    var linkJointID = evt.model.attributes.crlJointID;
    var linkID = "";
    if (linkJointID != "" && linkJointID != undefined) {
        linkID = crlGetConceptIDFromJointElementID(linkJointID);
    }
    switch (linkType) {
        case "crl.ReferenceLink":
            var sourceJointID = evt.sourceView.model.attributes.crlJointID;
            var targetJointID = evt.targetView.model.attributes.crlJointID;
            var sourceID = crlGetConceptIDFromJointElementID(sourceJointID);
            var targetID = crlGetConceptIDFromJointElementID(targetJointID);
            var targetAttributeName = "NoAttribute";
            switch (evt.targetView.model.attributes.represents) {
                case "OwnerPointer":
                    targetAttributeName = "OwningConceptID";
                    break;
                case "ElementPointer":
                    targetAttributeName = "ReferencedConceptID";
                    break;
                case "AbstractPointer":
                    targetAttributeName = "AbstractConceptID";
                    break;
                case "RefinedPointer":
                    targetAttributeName = "RefinedConceptID";
            }
            crlSendReferenceLinkChanged(evt.model, linkID, sourceID, targetID, targetAttributeName);
            break;
        case "crl.RefinementLink":
            var sourceJointID = evt.sourceView.model.attributes.crlJointID;
            var targetJointID = evt.targetView.model.attributes.crlJointID;
            var sourceID = crlGetConceptIDFromJointElementID(sourceJointID);
            var targetID = crlGetConceptIDFromJointElementID(targetJointID);
            crlSendRefinementLinkChanged(evt.model, linkID, sourceID, targetID);
            break;
        case "crl.OwnerPointer":
            var sourceJointID = evt.sourceView.model.attributes.crlJointID;
            var targetJointID = evt.targetView.model.attributes.crlJointID;
            var sourceID = crlGetConceptIDFromJointElementID(sourceJointID);
            var targetID = crlGetConceptIDFromJointElementID(targetJointID);
            crlSendOwnerPointerChanged(evt.model, linkID, sourceID, targetID);
            break;
        case "crl.ElementPointer":
            var sourceJointID = evt.sourceView.model.attributes.crlJointID;
            var targetJointID = evt.targetView.model.attributes.crlJointID;
            var sourceID = crlGetConceptIDFromJointElementID(sourceJointID);
            var targetID = crlGetConceptIDFromJointElementID(targetJointID);
            var targetAttributeName = "NoAttribute";
            switch (evt.targetView.model.attributes.represents) {
                case "OwnerPointer":
                    targetAttributeName = "OwningConceptID";
                    break;
                case "ElementPointer":
                    targetAttributeName = "ReferencedConceptID";
                    break;
                case "AbstractPointer":
                    targetAttributeName = "AbstractConceptID";
                    break;
                case "RefinedPointer":
                    targetAttributeName = "RefinedConceptID";
            }
            crlSendElementPointerChanged(evt.model, linkID, sourceID, targetID, targetAttributeName);
            break;
        case "crl.AbstractPointer":
            var sourceJointID = evt.sourceView.model.attributes.crlJointID;
            var targetJointID = evt.targetView.model.attributes.crlJointID;
            var sourceID = crlGetConceptIDFromJointElementID(sourceJointID);
            var targetID = crlGetConceptIDFromJointElementID(targetJointID);
            crlSendAbstractPointerChanged(evt.model, linkID, sourceID, targetID);
            break;
        case "crl.RefinedPointer":
            var sourceJointID = evt.sourceView.model.attributes.crlJointID;
            var targetJointID = evt.targetView.model.attributes.crlJointID;
            var sourceID = crlGetConceptIDFromJointElementID(sourceJointID);
            var targetID = crlGetConceptIDFromJointElementID(targetJointID);
            crlSendRefinedPointerChanged(evt.model, linkID, sourceID, targetID);
            break;
    }
    if (linkID == "") {
        evt.model.remove();
    }
}

function crlLinkViewRepresentsPointer(linkView) {
    var represents = linkView.model.attributes.represents
    if (represents == "OwnerPointer" || represents == "ElementPointer" || represents == "AbstractPointer" || represents == "RefinedPointer") {
        return true;
    }
    return false;
}

function crlNotificationRefresh() {
    crlSendNormalResponse();
    window.location.reload();
}

function crlNotificationSaveDebugSettings(data) {
    crlEnableTracing = JSON.parse(data.AdditionalParameters["EnableNotificationTracing"]);
    crlOmitHousekeepingCalls = JSON.parse(data.AdditionalParameters["OmitHousekeepingCalls"]);
    crlOmitManageTreeNodesCalls = JSON.parse(data.AdditionalParameters["OmitManageTreeNodesCalls"]);
    crlOmitDiagramRelatedCalls = JSON.parse(data.AdditionalParameters["OmitDiagramRelatedCalls"]);
    crlSendNormalResponse();
}

var crlPendingLinks = {};

function crlAddPendingLink(linkID, data) {
    crlPendingLinks[linkID] = data;
}

function crlProcessPendingLinks() {
    var beforeSize = _.size(crlPendingLinks);
    for (const linkID in crlPendingLinks) {
        data = crlPendingLinks[linkID];
        newlyCreatedLink = crlAddPendingDiagramLink(data);
        if (newlyCreatedLink) {
            delete crlPendingLinks[linkID];
        }
    }
    var afterSize = _.size(crlPendingLinks);
    if (afterSize < beforeSize) {
        crlProcessPendingLinks();
    }
}

function crlAddPendingDiagramLink(data) {
    var concept = data.NotificationConceptState;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var graph = crlGraphsGlobal[graphID];
    if (graph != null) {
        // The absence of a graph indicates that there is no view of the diagram at present
        var linkID = crlGetJointCellIDFromConceptID(concept.ConceptID);
        var link = crlFindLinkInGraph(graphID, linkID)
        var sourceJointID = crlGetJointCellIDFromConceptID(params["LinkSourceID"]);
        var targetJointID = crlGetJointCellIDFromConceptID(params["LinkTargetID"]);
        var linkSource = crlFindCellInGraph(graph, sourceJointID)
        var linkTarget = crlFindCellInGraph(graph, targetJointID)
        if (linkSource == null || linkTarget == null) {
            // the missing source or target still have not been created
            return null;
        }
        if (link == undefined || link == null) {
            link = crlConstructDiagramLink(data, graph, linkID);
            crlNotificationUpdateDiagramLink(data);
            return link;
        }
    }
}


function crlNotificationAddDiagramLink(data) {
    var concept = data.NotificationConceptState;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var graph = crlGraphsGlobal[graphID];
    if (graph != null) {
        // The absence of a graph indicates that there is no view of the diagram at present
        var linkID = crlGetJointCellIDFromConceptID(concept.ConceptID);
        var link = crlFindLinkInGraph(graphID, linkID)
        var sourceJointID = crlGetJointCellIDFromConceptID(params["LinkSourceID"]);
        var targetJointID = crlGetJointCellIDFromConceptID(params["LinkTargetID"]);
        var linkSource = crlFindCellInGraph(graph, sourceJointID)
        var linkTarget = crlFindCellInGraph(graph, targetJointID)
        if (linkSource == null || linkTarget == null) {
            // This case can arise when either the source or target is a link that has not yet been created
            crlPendingLinks[linkID] = data;
            crlSendNormalResponse();
            return;
        }
        if (link == undefined || link == null) {
            link = crlConstructDiagramLink(data, graph, linkID);
            newlyCreatedLink = link;
        }
        crlNotificationUpdateDiagramLink(data);
    }
    crlProcessPendingLinks();
    crlSendNormalResponse();
}

function crlNotificationAddDiagramNode(data) {
    var concept = data.NotificationConceptState;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var graph = crlGraphsGlobal[graphID];
    if (graph != null) {
        // The absence of a graph indicates the diagram is not being viewed
        var nodeID = crlGetJointCellIDFromConceptID(concept.ConceptID);
        var node = crlFindCellInGraph(graph, nodeID);
        if (node == undefined) {
            node = crlConstructDiagramNode(data, graph, nodeID);
        };
        crlNotificationUpdateDiagramNode(data);
        return; // crlNotificationUpdateDiagramNode will send the normal response
    }
    crlSendNormalResponse();
}

// <!-- Set up the websockets connection and callbacks -->
function crlNotificationAddTreeNode(data) {
    var concept = data.NotificationConceptState;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    if (owningConceptID == "") {
        owningConceptID = "#"
    } else {
        owningConceptID = "TreeNode" + owningConceptID
    }
    var nodeClass
    if (concept.ReadOnly == "true" || concept.IsCore == "true") {
        nodeClass = "node-read-only";
    } else {
        nodeClass = "node"
    }
    var tree = $('#uOfD').jstree();
    var nodeID = crlGetTreeNodeIDFromConceptID(concept.ConceptID);
    var node = tree.get_node(nodeID);
    if (node) {
        crlSendNormalResponse();
        return;
    }
    tree.create_node(owningConceptID,
        {
            'id': nodeID,
            'text': concept.Label,
            'icon': params.icon,
            'li_attr': {
                "read_only": concept.ReadOnly,
                "is_core": concept.IsCore,
                "is_diagram": params.isDiagram,
                "class": nodeClass
            }
        },
        'last');
    crlSendNormalResponse();
}

function crlNotificationClearToolbarSelection(data) {
    crlSelectToolbarButton("cursorToolbarButton");
    crlSendNormalResponse();
}

function crlNotificationClearTree() {
    $('#uOfD').jstree().destroy();
    crlInitializeTree();
    crlSendNormalResponse();
}

function crlNotificationCloseDiagramView(data) {
    var diagramID = data.NotificationConceptID;
    crlCloseDiagramView(diagramID);
    crlSendNormalResponse();
}

function crlNotificationDeleteDiagramCell(data) {
    var concept = data.NotificationConceptState;
    var elementID = crlGetJointCellIDFromConceptID(concept.ConceptID);
    var owningConceptID = data.AdditionalParameters["OwnerID"];
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var graph = crlGraphsGlobal[graphID];
    if (graph != null) {
        var element = crlFindCellInGraph(graph, elementID);
        if (element) {
            element.remove();
        }
    } else {
        console.log("************* In crlNotificationDeleteDiagramCell with null graph");
    }
    crlSendNormalResponse()
}

function crlNotificationDeleteTreeNode(data) {
    var conceptID = data.NotificationConceptID;
    var nodeID = crlGetTreeNodeIDFromConceptID(conceptID);
    $('#uOfD').jstree().delete_node(nodeID);
    crlSendNormalResponse();
}

function crlNotificationDiagramLabelChanged(data) {
    var tabID = crlGetDiagramTabIDFromDiagramID(data.NotificationConceptID);
    var tab = document.getElementById(tabID);
    tab.innerHTML = data.NotificationConceptState.Label;
    crlSendNormalResponse();
}

function crlNotificationDisplayDiagram(data) {
    var diagramID = data.NotificationConceptID;
    var diagramLabel = data.NotificationConceptState.Label;
    var diagramContainerID = crlGetDiagramContainerIDFromDiagramID(diagramID);
    var diagramContainer = document.getElementById(diagramContainerID);
    // Construct the container if it is not already present
    if (diagramContainer == undefined) {
        diagramContainer = crlConstructDiagramContainer(diagramContainer, diagramContainerID, diagramLabel, diagramID);
    }
    crlDisplayDiagram(diagramContainer.id);
    crlCurrentDiagramContainerID = diagramContainerID;
    crlSetDefaultLink();
    crlSendNormalResponse();
}

function crlNotificationDisplayGraph(data) {
    var graphString = data.AdditionalParameters["GraphString"];
    const workerURL = '/js/full.render.js';
    let viz = new Viz({ workerURL });
    var newTab = window.open("graph.html");
    viz.renderSVGElement(graphString)
        .then(function (element) {
            newTab.document.body.appendChild(element);
        })
        .catch(error => {
            // Create a new Viz instance (@see Caveats page for more info)
            viz = new Viz({ workerURL });

            // Possibly display the error
            console.error(error);
        });
    crlSendNormalResponse();
}

function crlNotificationDoesLinkExist(data) {
    var concept = data.NotificationConceptState;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var graph = crlGraphsGlobal[graphID];
    if (graph != null) {
        // The absence of a graph indicates that there is no view of the diagram at present
        var linkID = crlGetJointCellIDFromConceptID(concept.ConceptID);
        var link = crlFindLinkInGraph(graphID, linkID)
        var sourceJointID = crlGetJointCellIDFromConceptID(params["LinkSourceID"]);
        var targetJointID = crlGetJointCellIDFromConceptID(params["LinkTargetID"]);
        var linkSource = crlFindCellInGraph(graph, sourceJointID)
        var linkTarget = crlFindCellInGraph(graph, targetJointID)
        if (link != null && linkSource != null && linkTarget != null) {
            crlSendBooleanResponse(true);
            crlSendNormalResponse();
            return;
        }
    }
    crlSendBooleanResponse(false);
}

function crlNotificationElementSelected(data) {
    if (data.NotificationConceptID != crlSelectedConceptID) {
        crlSelectedConceptID = data.NotificationConceptID
        // Update the properties
        crlUpdateProperties(data);
        // Update the tree
        var treeNodeID = crlGetTreeNodeIDFromConceptID(crlSelectedConceptID);
        $("#uOfD").jstree(true).deselect_all(true);
        // a hack tp prevent infinite recursion
        crlInCrlElementSelected = true;
        $("#uOfD").jstree(true).select_node(treeNodeID, true);
        crlInCrlElementSelected = false;
    }
    crlSendNormalResponse()
}

function crlUpdateProperties(data) {
    crlPropertiesDisplayType(data, 1);
    crlPropertiesDisplayID(data, 2);
    crlPropertiesDisplayOwningConcept(data, 3);
    crlPropertiesDisplayVersion(data, 4);
    crlPropertiesDisplayLabel(data, 5);
    crlPropertiesDisplayDefinition(data, 6);
    crlPropertiesDisplayURI(data, 7);
    crlPropertiesDisplayReadOnly(data, 8);
    var type = "";
    if (data.NotificationConceptState) {
        type = data.NotificationConceptState.ConceptType;
    }
    switch (type) {
        case "*core.element":
            crlPropertiesClearRow(10);
            crlPropertiesClearRow(9);
            break;
        case "*core.literal":
            crlPropertiesDisplayLiteralValue(data, 9);
            crlPropertiesClearRow(10);
            break;
        case "*core.reference":
            crlPropertiesDisplayReferencedConcept(data, 9);
            crlPropertiesDisplayReferencedAttributeName(data, 10)
            break;
        case "*core.refinement":
            crlPropertiesDisplayAbstractConcept(data, 9);
            crlPropertiesDisplayRefinedConcept(data, 10);
            break;
        default:
            crlPropertiesClearRow(10);
            crlPropertiesClearRow(9);
    };
}

function crlNotificationSaveUserPreferences(data) {
    crlDropReferenceAsLink = JSON.parse(data.AdditionalParameters["DropReferenceAsLink"]);
    crlDropRefinementAsLink = JSON.parse(data.AdditionalParameters["DropRefinementAsLink"]);
    crlSendNormalResponse();
}

function crlNotificationShowTreeNode(data) {
    var concept = data.NotificationConceptState;
    var nodeID = crlGetTreeNodeIDFromConceptID(concept.ConceptID);
    var tree = $('#uOfD').jstree();
    tree.select_node(nodeID);
    crlSendNormalResponse();
}

var crlNotificationUpdateDiagramLink = function (data) {
    var concept = data.NotificationConceptState;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var graph = crlGraphsGlobal[graphID];
    if (graph != null) {
        // The absence of a graph indicates that there is no view of the diagram at present
        var linkID = crlGetJointCellIDFromConceptID(concept.ConceptID);
        var link = crlFindLinkInGraph(graphID, linkID)
        var sourceJointID = crlGetJointCellIDFromConceptID(data.AdditionalParameters["LinkSourceID"]);
        var targetJointID = crlGetJointCellIDFromConceptID(data.AdditionalParameters["LinkTargetID"]);
        var linkSource = crlFindCellInGraph(graph, sourceJointID)
        var linkTarget = crlFindCellInGraph(graph, targetJointID)
        if (link == undefined || link == null) {
            crlSendNormalResponse()
            return;
        }
        if ((linkSource == null || linkTarget == null) && (link != null && link != undefined)) {
            link.remove();
            crlSendNormalResponse()
            return;
        }
        link.label(0, {
            attrs: {
                text: {
                    text: data.AdditionalParameters["DisplayLabel"]
                }
            }
        });
        if (link.source().id != linkSource.id) {
            link.set("source", linkSource);
        }
        if (link.target().id != linkTarget.id) {
            link.set("target", linkTarget);
        }
    }
    crlSendNormalResponse()
}

var crlNotificationUpdateDiagramNode = function (data) {
    var concept = data.NotificationConceptState;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var graph = crlGraphsGlobal[graphID];
    if (graph != null) {
        // The absence of a graph indicates the diagram is not being viewed
        var nodeID = crlGetJointCellIDFromConceptID(concept.ConceptID);
        var node = crlFindCellInGraph(graph, nodeID);
        if (node == undefined) {
            crlSendNormalResponse();
            return;
        };
        node.set("displayLabelYOffset", Number(params["DisplayLabelYOffset"]));
        node.set('position', { "x": Number(params["NodeX"]), "y": Number(params["NodeY"]) });
        node.set('name', params["DisplayLabel"]);
        node.set('size', { "width": Number(params["NodeWidth"]), "height": Number(params["NodeHeight"]) });
        node.set('icon', params["Icon"]);
        node.set("abstractions", params["Abstractions"]);
        node.set("lineColor", params["LineColor"]);
        node.set("bgColor", params["BGColor"]);
    }
    crlSendNormalResponse();
}

function crlNotificationUpdateTreeNode(data) {
    var concept = data.NotificationConceptState;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var treeNodeOwnerID = ""
    if (owningConceptID == "") {
        treeNodeOwnerID = "#";
    } else {
        treeNodeOwnerID = crlGetTreeNodeIDFromConceptID(owningConceptID);
    };
    var nodeID = crlGetTreeNodeIDFromConceptID(concept.ConceptID);
    var tree = $('#uOfD').jstree();
    if (tree.get_parent(nodeID) != treeNodeOwnerID) {
        tree.move_node(nodeID, treeNodeOwnerID);
    }
    var nodeClass
    if (concept.ReadOnly == "true" || concept.IsCore == "true") {
        nodeClass = "node-read-only";
    } else {
        nodeClass = "node"
    }
    tree.rename_node(nodeID, concept.Label);
    tree.set_icon(nodeID, params.icon);
    var node = tree.get_node(nodeID);
    if (node) {
        node.li_attr.read_only = concept.ReadOnly;
        node.li_attr.is_core = concept.IsCore;
        node.li_attr.is_diagram = params.isDiagram;
        node.li_attr.class = nodeClass;
    }
    if (concept.ConceptID == crlSelectedConceptID) {
        crlUpdateProperties(data);
    }
    crlSendNormalResponse()
}

var crlNotificationUpdateProperties = function (data) {
    crlUpdateProperties(data);
    crlSendNormalResponse()
}

var crlNotificationUpdateWorkspacePath = function (data) {
    crlWorkspacePath = data.AdditionalParameters["WorkspacePath"];
    document.getElementById("CurrentWorkspace").innerHTML = crlWorkspacePath;
    crlSendNormalResponse();
}

function crlObtainPropertyRow(row) {
    var properties = document.getElementById("properties");
    var propertyRow = properties.rows[row];
    if (propertyRow == undefined) {
        propertyRow = properties.insertRow(row);
        propertyRow.insertCell(0);
        propertyRow.insertCell(1);
        var input = document.createElement("input");
        input.setAttribute("type", "text");
        propertyRow.cells[1].appendChild(input);
    }
    return propertyRow
}

var crlOnChangePosition = function (modelElement, position) {
    if (crlMouseButtonPressed[0] == true) {
        var jointElementID = modelElement.get("crlJointID");
        var diagramNodeID = crlGetConceptIDFromJointElementID(jointElementID);
        crlMovedNodes[diagramNodeID] = position;
    }
}

var crlOnCloseDiagramView = function () {
    var diagramID = crlGetDiagramIDFromDiagramContainerID(crlCurrentDiagramContainerID);
    crlCloseDiagramView(diagramID);
}

var crlOnDiagramCellPointerDown = function (cellView, event, x, y) {
    crlMouseButtonPressed[event.button] = true;
    var jointElementID = cellView.model.get("crlJointID");
    if (jointElementID && jointElementID != "") {
        var diagramNodeID = crlGetConceptIDFromJointElementID(jointElementID);
        if (diagramNodeID == "") {
            console.log("In onDiagramCellPointerDown diagramNodeID is empty")
        }
        crlSendDiagramElementSelected(diagramNodeID)
    }
}

var crlOnDiagramCellPointerUp = function (cellView, event, x, y) {
    crlMouseButtonPressed[event.button] = false;
    crlFinalizeNodeMoves();
}

function crlFinalizeNodeMoves() {
    $.each(crlMovedNodes, function (nodeID, position) {
        crlSendDiagramNodeNewPosition(nodeID, position);
    });
    crlMovedNodes = {};
}

function crlOnDiagramClick(event) {
    var nodeType = ""
    switch (crlCurrentToolbarButton) {
        case "cursorToolbarButton": {
            break;
        }
        case "elementToolbarButton": {
            nodeType = "Element";
            break;
        }
        case "literalToolbarButton": {
            nodeType = "Literal";
            break;
        }
        case "referenceToolbarButton": {
            nodeType = "Reference";
            break;
        }
        case "referenceLinkToolbarButton": {
            break;
        }
        case "refinementToolbarButton": {
            nodeType = "Refinement";
            break;
        }
        case "refinementLinkToolbarButton": {
            break;
        }
        case "diagramToolbarButton": {
            nodeType = "Diagram";
            break;
        }
        case "ownerPointerToolbarButton": {
            break;
        }
        case "elementPointerToolbarButton": {
            break;
        }
        case "abstractPointerToolbarButton": {
            break;
        }
        case "refinedPointerToolbarButton": {
            break;
        }
        default:
            console.log("In crlOnDiagramClick, unknown toolbar button type: " + crlCurrentToolbarButton)
    }
    if (nodeType != "") {
        var conceptID = crlGetConceptIDFromContainerID(event.target.parentElement.parentElement.id);
        var x = event.layerX.toString();
        var y = event.layerY.toString();
        crlSendDiagramClick(nodeType, conceptID, x, y);
    }
};

function crlOnDiagramDrop(event) {
    event.preventDefault();
    var conceptID = crlGetConceptIDFromContainerID(event.target.parentElement.parentElement.id);
    var x = event.layerX.toString();
    var y = event.layerY.toString();
    crlSendDiagramDrop(conceptID, x, y, event.shiftKey);
};

function crlOnDragover(event, data) {
    event.preventDefault();
};

function crlOnEditorDrop(e, data) {
    crlSendSetTreeDragSelection("");
};

function crlOnMagnet(evt, x, y) {
    this.dragMagnetStart(evt, x, y);
};

function crlOnMakeDiagramVisible(e) {
    var diagramContainerID = e.target.getAttribute("diagramContainerID")
    var diagramID = crlGetDiagramIDFromDiagramContainerID(diagramContainerID);
    crlSendDisplayDiagramSelected(diagramID);
};

var crlOnToolbarButtonSelected = function (e, data) {
    var img = e.target;
    var btn = img.parentElement;
    var id = btn.id;
    crlSelectToolbarButton(id);
};

function crlOnDiagramMouseOver(mouseEvent) {
    var diagram = mouseEvent.currentTarget;
    if (crlCurrentToolbarButton == crlElementToolbarButtonID ||
        crlCurrentToolbarButton == crlLiteralToolbarButtonID ||
        crlCurrentToolbarButton == crlReferenceToolbarButtonID ||
        crlCurrentToolbarButton == crlRefinementToolbarButtonID ||
        crlCurrentToolbarButton == crlDiagramToolbarButtonID) {
        diagram.style.cursor = "cell";
    } else {
        diagram.style.cursor = "default";
    }
};

function crlSelectToolbarButton(id) {
    crlCurrentToolbarButton = id;
    var toolbar = document.getElementById("toolbar");
    var buttons = toolbar.children;
    for (var i = 0; i < buttons.length; i++) {
        var button = buttons[i];
        if (button.id == id) {
            button.children[0].style.backgroundColor = "white";
        }
        else {
            button.children[0].style.backgroundColor = "grey";
        }
    }
    crlSetDefaultLink();
};

function crlOnTreeDragStart(e, data) {
    var parentID = e.target.parentElement.id;
    var selectedElementID = crlGetConceptIDFromTreeNodeID(parentID);
    crlSendSetTreeDragSelection(selectedElementID);
};


function crlOpenDiagramContainer(diagramContainerId) {
    var i;
    var x = document.getElementsByClassName("crlDiagramContainer");
    for (i = 0; i < x.length; i++) {
        if (x[i].id == diagramContainerId) {
            x[i].style.display = "block";
        } else {
            x[i].style.display = "none";
        }
    }
}



function crlPopulateToolbar() {
    var toolbar = document.getElementById("toolbar");
    crlAppendToolbarButton(toolbar, crlCursorToolbarButtonID, "/icons/CursorIcon.svg");
    crlAppendToolbarButton(toolbar, crlElementToolbarButtonID, "/icons/ElementIcon.svg");
    crlAppendToolbarButton(toolbar, crlLiteralToolbarButtonID, "/icons/LiteralIcon.svg");
    crlAppendToolbarButton(toolbar, crlReferenceToolbarButtonID, "/icons/ReferenceIcon.svg");
    crlAppendToolbarButton(toolbar, crlReferenceLinkToolbarButtonID, "/icons/ReferenceLinkIcon.svg");
    crlAppendToolbarButton(toolbar, crlRefinementToolbarButtonID, "/icons/RefinementIcon.svg");
    crlAppendToolbarButton(toolbar, crlRefinementLinkToolbarButtonID, "/icons/RefinementLinkIcon.svg");
    crlAppendToolbarButton(toolbar, crlDiagramToolbarButtonID, "/icons/DiagramIcon.svg");
    crlAppendToolbarButton(toolbar, crlOwnerPointerToolbarButtonID, "/icons/OwnerPointerIcon.svg");
    crlAppendToolbarButton(toolbar, crlElementPointerToolbarButtonID, "/icons/ElementPointerIcon.svg");
    crlAppendToolbarButton(toolbar, crlAbstractPointerToolbarButtonID, "/icons/AbstractPointerIcon.svg");
    crlAppendToolbarButton(toolbar, crlRefinedPointerToolbarButtonID, "/icons/RefinedPointerIcon.svg");
    crlSelectToolbarButton(crlCursorToolbarButtonID);
}

function crlSendAbstractPointerChanged(jointLink, linkID, sourceID, targetID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "AbstractPointerChanged",
        "RequestConceptID": linkID,
        "AdditionalParameters": {
            "SourceID": sourceID,
            "TargetID": targetID
        }
    })
    crlSendRequest(xhr, data);
}

function crlSendAddDiagramChild(conceptID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "AddDiagramChild", "RequestConceptID": conceptID });
    crlSendRequest(xhr, data);
}

function crlSendAddElementChild(conceptID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "AddElementChild", "RequestConceptID": conceptID });
    crlSendRequest(xhr, data);
}

function crlSendAddLiteralChild(conceptID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "AddLiteralChild", "RequestConceptID": conceptID });
    crlSendRequest(xhr, data);
}

function crlSendAddReferenceChild(conceptID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "AddReferenceChild", "RequestConceptID": conceptID });
    crlSendRequest(xhr, data);
}

function crlSendAddRefinementChild(conceptID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "AddRefinementChild", "RequestConceptID": conceptID });
    crlSendRequest(xhr, data);
}

function crlSendDiagramViewHasBeenClosed(diagramID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "DiagramViewHasBeenClosed",
        "RequestConceptID": diagramID
    });
    crlSendRequest(xhr, data);
}

function crlSendClearWorkspace() {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "ClearWorkspace"
    });
    crlSendRequest(xhr, data);
}

function crlSendCloseWorkspace() {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "CloseWorkspace"
    });
    crlSendRequest(xhr, data);
}

function crlSendDebugSettings(enableNotificationTracing, omitHousekeepingCalls, omitManageTreeNodesCalls, omitDiagramRelatedCalls, maxTracingDepth) {
    var xhr = crlCreateEmptyRequest()
    var data = JSON.stringify({
        "Action": "UpdateDebugSettings",
        "AdditionalParameters": {
            "EnableNotificationTracing": enableNotificationTracing,
            "OmitHousekeepingCalls": omitHousekeepingCalls,
            "OmitManageTreeNodesCalls": omitManageTreeNodesCalls,
            "OmitDiagramRelatedCalls": omitDiagramRelatedCalls,
            "MaxTracingDepth": maxTracingDepth
        }
    });
    crlSendRequest(xhr, data);
}

function crlSendDefinitionChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    // This game is required to eliminate special characters
    var newValue = crlIgnoreSpecialCharacters(evt.currentTarget.value);
    evt.currentTarget.value = newValue;
    var data = JSON.stringify({
        "Action": "DefinitionChanged",
        "RequestConceptID": crlSelectedConceptID,
        "AdditionalParameters":
            { "NewValue": newValue }
    });
    crlSendRequest(xhr, data);
}

function crlSendDiagramClick(nodeType, diagramID, x, y) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "DiagramClick",
        "AdditionalParameters":
        {
            "DiagramID": diagramID,
            "NodeType": nodeType,
            "NodeX": x,
            "NodeY": y
        }
    });
    crlSendRequest(xhr, data);
}

function crlSendDiagramDrop(diagramID, x, y, shiftKey) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "DiagramDrop",
        "AdditionalParameters":
        {
            "DiagramID": diagramID,
            "NodeX": x,
            "NodeY": y,
            "Shift": shiftKey.toString()
        }
    });
    crlSendRequest(xhr, data);
}

function crlSendDiagramNodeNewPosition(nodeID, position) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "DiagramNodeNewPosition",
        "RequestConceptID": nodeID,
        "AdditionalParameters": {
            "NodeX": position.x.toString(),
            "NodeY": position.y.toString()
        }
    })
    crlSendRequest(xhr, data);
}

function crlSendDiagramElementSelected(nodeID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "DiagramElementSelected", "RequestConceptID": nodeID });
    crlSendRequest(xhr, data);
}

function crlSendDisplayCallGraph(index) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "DisplayCallGraph",
        "AdditionalParameters": {
            "GraphIndex": index
        }
    });
    crlSendRequest(xhr, data);
}

function crlSendDisplayDiagramSelected(diagramID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "DisplayDiagramSelected", "RequestConceptID": diagramID });
    crlSendRequest(xhr, data);
}

function crlSendUserPreferences() {
    var xhr = crlCreateEmptyRequest()
    var dropReferenceAsLink = "false";
    var dropRefinementAsLink = "false";
    if ($("#dropReferenceAsLink").prop("checked") == true) {
        dropReferenceAsLink = "true";
    }
    if ($("#dropRefinementAsLink").prop("checked") == true) {
        dropRefinementAsLink = "true";
    }
    var data = JSON.stringify({
        "Action": "UpdateUserPreferences",
        "AdditionalParameters": {
            "DropReferenceAsLink": dropReferenceAsLink,
            "DropRefinementAsLink": dropRefinementAsLink
        }
    });
    crlSendRequest(xhr, data);
}

function crlSendElementPointerChanged(jointLink, linkID, sourceID, targetID, targetAttributeName) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "ElementPointerChanged",
        "RequestConceptID": linkID,
        "AdditionalParameters": {
            "SourceID": sourceID,
            "TargetID": targetID,
            "TargetAttributeName": targetAttributeName
        }
    })
    crlSendRequest(xhr, data);
}

function crlSendDiagramElementFormatChanged(diagramElementID, lineColor, bgColor) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "FormatChanged",
        "RequestConceptID": diagramElementID,
        "AdditionalParameters": {
            "LineColor": lineColor,
            "BGColor": bgColor
        }
    });
    crlSendRequest(xhr, data)
}

function crlSendLabelChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    // This game is required to eliminate special characters
    var newValue = crlIgnoreSpecialCharacters(evt.currentTarget.value);
    evt.currentTarget.value = newValue;
    var data = JSON.stringify({
        "Action": "LabelChanged",
        "RequestConceptID": crlSelectedConceptID,
        "AdditionalParameters":
            { "NewValue": newValue }
    });
    crlSendRequest(xhr, data)
}

function crlSendLiteralValueChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    // This game is required to eliminate special characters
    var newValue = crlIgnoreSpecialCharacters(evt.currentTarget.value);
    evt.currentTarget.value = newValue;
    var data = JSON.stringify({
        "Action": "LiteralValueChanged",
        "RequestConceptID": crlSelectedConceptID,
        "AdditionalParameters":
            { "NewValue": newValue }
    });
    crlSendRequest(xhr, data)
}

function crlSendNewDomainRequest(evt) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "NewDomainRequest" });
    crlSendRequest(xhr, data)
}

function crlSendOpenWorkspace() {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "OpenWorkspace",
    });
    crlSendRequest(xhr, data);
}

function crlSendOwnerPointerChanged(jointLink, linkID, sourceID, targetID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "OwnerPointerChanged",
        "RequestConceptID": linkID,
        "AdditionalParameters": {
            "SourceID": sourceID,
            "TargetID": targetID
        }
    })
    crlSendRequest(xhr, data);
}

function crlSendRedo() {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "Redo" });
    crlSendRequest(xhr, data);
}

function crlSendReferenceLinkChanged(jointLink, linkID, sourceID, targetID, targetAttributeName) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "ReferenceLinkChanged",
        "RequestConceptID": linkID,
        "AdditionalParameters": {
            "SourceID": sourceID,
            "TargetID": targetID,
            "TargetAttributeName": targetAttributeName
        }
    })
    crlSendRequest(xhr, data);
}

function crlSendRefinementLinkChanged(jointLink, linkID, sourceID, targetID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "RefinementLinkChanged",
        "RequestConceptID": linkID,
        "AdditionalParameters": {
            "SourceID": sourceID,
            "TargetID": targetID
        }
    })
    crlSendRequest(xhr, data);
}

function crlSendRefreshDiagram(diagramID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "RefreshDiagram",
        "RequestConceptID": diagramID
    });
    crlSendRequest(xhr, data);
}

function crlSendRefinedPointerChanged(jointLink, linkID, sourceID, targetID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "RefinedPointerChanged",
        "RequestConceptID": linkID,
        "AdditionalParameters": {
            "SourceID": sourceID,
            "TargetID": targetID
        }
    })
    crlSendRequest(xhr, data);
}

function crlSendRequest(xhr, data) {
    console.log(data);
    xhr.send(data);
}

function crlSendSaveWorkspace() {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "SaveWorkspace"
    });
    crlSendRequest(xhr, data);
}

function crlSendURIChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "URIChanged",
        "RequestConceptID": crlSelectedConceptID,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.value }
    });
    crlSendRequest(xhr, data);
}

function crlSendBooleanResponse(booleanValue) {
    var data = {};
    data["Result"] = 0;
    data["ErrorMessage"] = "none";
    data["ResultConceptID"] = "foo";
    if (booleanValue) {
        data["BooleanValue"] = "true";
    } else {
        data["BooleanValue"] = "false";
    };
    crlWebsocketGlobal.send(JSON.stringify(data));
    console.log(data);
}

function crlSendNormalResponse() {
    var data = {};
    data["Result"] = 0;
    data["ErrorMessage"] = "none";
    crlWebsocketGlobal.send(JSON.stringify(data));
    console.log(data);
}

function crlSendSetTreeDragSelection(id) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "SetTreeDragSelection", "RequestConceptID": id });
    crlSendRequest(xhr, data);
}

function crlSendTreeNodeDelete(id) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "TreeNodeDelete", "RequestConceptID": id });
    crlSendRequest(xhr, data);
}

var crlSendTreeNodeSelected = function (evt, obj) {
    if (obj != undefined) {
        var conceptID = crlGetConceptIDFromTreeNodeID(obj.node.id)
        if (conceptID != crlSelectedConceptID && crlInCrlElementSelected == false) {
            var xhr = crlCreateEmptyRequest();
            var data = JSON.stringify({ "Action": "TreeNodeSelected", "RequestConceptID": conceptID });
            crlSendRequest(xhr, data);
        }
    };
}

function crlSendUndo() {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "Undo" });
    crlSendRequest(xhr, data);
}

function crlSetDefaultLink() {
    var paper
    if (crlCurrentDiagramContainerID) {
        var diagramID = crlGetDiagramIDFromDiagramContainerID(crlCurrentDiagramContainerID);
        var paperID = crlGetJointPaperIDFromDiagramID(diagramID);
        paper = crlPapersGlobal[paperID]
        paper.options.defaultLink = crlCreateLink;
    }
}

var crlShowAbstractConcept = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "ShowAbstractConcept", "RequestConceptID": diagramElementID });
    crlSendRequest(xhr, data);
}

var crlShowModelConceptInNavigator = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "ShowModelConceptInNavigator", "RequestConceptID": diagramElementID });
    crlSendRequest(xhr, data);
}

var crlShowDiagramElementInNavigator = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "ShowDiagramElementInNavigator", "RequestConceptID": diagramElementID });
    crlSendRequest(xhr, data);
}

var crlShowOwnedConcepts = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "ShowOwnedConcepts", "RequestConceptID": diagramElementID });
    crlSendRequest(xhr, data);
}

var crlShowOwner = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "ShowOwner", "RequestConceptID": diagramElementID });
    crlSendRequest(xhr, data);
}

var crlNullifyReferencedConcept = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "NullifyReferencedConcept", "RequestConceptID": diagramElementID });
    crlSendRequest(xhr, data);
}

var crlShowReferencedConcept = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "ShowReferencedConcept", "RequestConceptID": diagramElementID });
    crlSendRequest(xhr, data);
}

var crlShowRefinedConcept = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "ShowRefinedConcept", "RequestConceptID": diagramElementID });
    crlSendRequest(xhr, data);
}

// The default validate connection allows connections to ElementViews only. CRL allows connections
// to links as well.
function crlValidateConnection(cellViewS, magnetS, cellViewT, magnetT, end, linkView) {
    var represents = linkView.model.attributes.represents;
    if (!cellViewT) {
        return false;
    }
    var targetRepresents = cellViewT.model.attributes.represents;
    switch (represents) {
        case "Reference":
            return true;
        case "ElementPointer":
            return true;
        case "Refinement":
            return targetRepresents == "Element" ||
                targetRepresents == "Literal" ||
                targetRepresents == "Refinement" ||
                targetRepresents == "Reference";
        default:
            // Must be an owner, abstract, or refined pointer
            return targetRepresents == "Element" ||
                targetRepresents == "Literal" ||
                targetRepresents == "Refinement" ||
                targetRepresents == "Reference";
    }
}

function crlValidateLinkStart(cellView, magnet) {
    var represents = cellView.model.attributes.represents;
    switch (crlCurrentToolbarButton) {
        case crlCursorToolbarButtonID:
            return false;
        case crlElementToolbarButtonID:
            return false;
        case crlLiteralToolbarButtonID:
            return false;
        case crlReferenceToolbarButtonID:
            return false;
        case crlReferenceLinkToolbarButtonID:
            if (represents == "Element" || represents == "Literal" || represents == "Refinement" || represents == "Reference") {
                return true;
            } else {
                return false;
            }
        case crlRefinementToolbarButtonID:
            return false;
        case crlRefinementLinkToolbarButtonID:
            if (represents == "Element" || represents == "Literal" || represents == "Refinement" || represents == "Reference") {
                return true;
            } else {
                return false;
            }
        case crlDiagramToolbarButtonID:
            return false;
        case crlOwnerPointerToolbarButtonID:
            if (represents == "Element" || represents == "Literal" || represents == "Refinement" || represents == "Reference") {
                return true;
            } else {
                return false;
            }
        case crlElementPointerToolbarButtonID:
            if (represents == "Reference") {
                return true;
            } else {
                return false;
            }
        case crlAbstractPointerToolbarButtonID:
            if (represents == "Refinement") {
                return true;
            } else {
                return false;
            }
        case crlRefinedPointerToolbarButtonID:
            if (represents == "Refinement") {
                return true;
            } else {
                return false;
            }

    }
    return false;
}

