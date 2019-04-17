// crlCurrentDiagramContainerIDGlobal is the identifier for the diagram container currently being displayed
var crlCurrentDiagramContainerID;
// crlCurrentToolbarButton is the id of the last toolbar button pressed
var crlCurrentToolbarButton;
// crlDebugSettingsDialog is the initialized dialog used for editing debug settings
var crlDebugSettingsDialog;
// crlDropReferenceAsLink when true causes references dragged from the tree into the diagram to be added as links
var crlDropReferenceAsLink = false;
// crlDropRefinementAsLink when true causes refinements dragged from the tree into the diagram to be added as links
var crlDropRefinementAsLink = false;
// crlEditorSettingsDialog is the initialized dialog used for editing editor settings
var crlEditorSettingsDialog;
// crlEnableTracing is the client-side copy of the server-side value that turns on notification tracing
var crlEnableTracing = false;
// crlOmitHousekeepingCalls indicates whether housekeeping calls shouldl be included when tracing is enabled
var crlOmitHousekeepingCalls = false;
// crlGraphsGlobal is an array of existing graphs that is used to look up a graph given its identifier
var crlGraphsGlobal = {};
// crlInitializationCompleteGlobal indicates whether the server-side initialization has been completed
var crlInitializationComplete = false;
// crlMovedNodes is an array of nodes that have been moved. This is a temporary cache that is used to update the 
// server once a mouse up has occurred
var crlMovedNodes = {};
// crlOpenWorkspaceDialog is the initialized dialog used for opening a workspace
var crlOpenWorkspaceDialog;
// crlPaperGlobal is an array of existing papers that is used to look up a paper given its identifier
var crlPapersGlobal = {};
// crlSelectedConceptIDGlobal contains the model identifier of the currently selected concept
var crlSelectedConceptIDGlobal;
// crlTreeDragSelectionIDGlobal contains the model identifier of the concept currently being dragged from the tree
var crlTreeDragSelectionIDGlobal;
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

var crlInCrlElementSelected = false;

var crlKeyPressed = { 16: false };

// Initialize
$(function () {
    $(".uofd-browser").resizable({
        resizeHeight: false
    });
    $("#uOfD").jstree({
        'core': {
            'check_callback': true,
            'multiple': false
        },
        'plugins': ['sort', 'contextmenu', 'wholerow'],
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
                                        var xhr = crlCreateEmptyRequest();
                                        var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                        var data = JSON.stringify({ "Action": "AddElementChild", "RequestConceptID": conceptID });
                                        crlSendRequest(xhr, data);
                                    }
                                }
                            },
                            Diagram: {
                                "label": "Diagram",
                                "action": function (obj) {
                                    var xhr = crlCreateEmptyRequest();
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                    var data = JSON.stringify({ "Action": "AddDiagramChild", "RequestConceptID": conceptID });
                                    crlSendRequest(xhr, data);
                                }
                            },
                            Literal: {
                                "label": "Literal",
                                "action": function (obj) {
                                    var xhr = crlCreateEmptyRequest();
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                    var data = JSON.stringify({ "Action": "AddLiteralChild", "RequestConceptID": conceptID });
                                    crlSendRequest(xhr, data);
                                }
                            },
                            Reference: {
                                "label": "Reference",
                                "action": function (obj) {
                                    var xhr = crlCreateEmptyRequest();
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                    var data = JSON.stringify({ "Action": "AddReferenceChild", "RequestConceptID": conceptID });
                                    crlSendRequest(xhr, data);
                                }
                            },
                            Refinement: {
                                "label": "Refinement",
                                "action": function (obj) {
                                    var xhr = crlCreateEmptyRequest();
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                    var data = JSON.stringify({ "Action": "AddRefinementChild", "RequestConceptID": conceptID });
                                    crlSendRequest(xhr, data);
                                }
                            }
                        }
                    },
                    display: {
                        "label": "Display Diagram",
                        "action": function (obj) {
                            if ($node != undefined) {
                                //                                var xhr = crlCreateEmptyRequest();
                                var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                crlSendDisplayDiagramSelected(conceptID);
                                //     var data = JSON.stringify({ "Action": "DisplayDiagramSelected", "RequestConceptID": conceptID });
                                //     crlSendRequest(xhr, data);
                            }
                        }
                    },
                    remove: {
                        "label": "Delete",
                        "action": function (obj) {
                            if ($node != undefined) {
                                var xhr = crlCreateEmptyRequest();
                                var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                var data = JSON.stringify({ "Action": "TreeNodeDelete", "RequestConceptID": conceptID });
                                crlSendRequest(xhr, data);
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
    $("#body").on("ondrop", crlOnEditorDrop);
    crlPopulateToolbar();
    crlDiagramCellDropdownMenu = document.getElementById("diagramCellDropdown");
    if (crlDiagramCellDropdownMenu) {
        crlDiagramCellDropdownMenu.addEventListener("mouseleave", function () {
            crlDiagramCellDropdownMenu.style.display = "none";
        })
        crlDiagramCellDropdownMenu.addEventListener("mouseup", function () {
            crlDiagramCellDropdownMenu.style.display = "none";
        })
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
            "       <input type='checkbox' id='omitManageTreeNodesCalls'>" +
            "	</fieldset> " +
            "</form>",
        confirm: function () {
            var enableNotificationTracing = "false";
            if ($("#enableTracing").prop("checked") == true) {
                enableNotificationTracing = "true";
            };
            var omitHousekeepingCalls = "false"
            if ($("#omitHousekeepingCalls").prop("checked") == true) {
                omitHousekeepingCalls = "true";
            }
            var omitManageTreeNodesCalls = "false"
            if ($("#omitManageTreeNodesCalls").prop("checked") == true) {
                omitManageTreeNodesCalls = "true";
            }
            crlSendDebugSettings(enableNotificationTracing, omitHousekeepingCalls, omitManageTreeNodesCalls, "0");
        },
        onOpen: function () {
            $("#enableTracing").prop("checked", crlEnableTracing);
            $("#omitHousekeepingCalls").prop("checked", crlOmitHousekeepingCalls);
        }
    });
    crlDisplayCallGraphsDialog = new jBox("Confirm", {
        title: "Display Call Graphs",
        confirmButton: "OK",
        cancelButton: "Cancel",
        content: "" +
            "<form>" +
            "    <fieldseet>" +
            "        <label>Number of available call graphs: </label> <label id='numberOfAvailableGraphs'></label><br>" +
            "        <label>Selected Graph:</label><input id='selectedGraph' type='number'>" +
            "	</fieldset> " +
            "</form>",
        confirm: function () {
            var selectedNumber = $("#selectedGraph").val();
            crlSendDisplayCallGraph(selectedNumber);
        }
    })
    crlEditorSettingsDialog = new jBox("Confirm", {
        title: "Editor Settings",
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
            crlSendEditorSettings();
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
            "		Identify Workspace Folder:<input type='file'><br>" +
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
    // Close the dropdown menu if the user clicks outside of it
    window.onkeydown = function (e) {
        e = e || window.event;
        crlKeyPressed[e.keyCode] = true;
    };
    window.onkeyup = function (e) {
        e = e || window.event;
        crlKeyPressed[e.keyCode] = false;
    };
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
        "width": 1000,
        "height": 1000,
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
    jointPaper.on("cell:contextmenu", function (cellView, evt, x, y) {
        evt.preventDefault();
        var represents = cellView.model.attributes.represents;
        if (represents == "Reference") {
            document.getElementById("showReferencedConcept").style.display = "block";
        } else {
            document.getElementById("showReferencedConcept").style.display = "none";
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

function crlPropertiesClearRow(row) {
    var properties = document.getElementById("properties");
    var propertyRow = properties.rows[row]
    if (propertyRow != undefined) {
        properties.deleteRow(row);
    }
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
                alert(response.ResultDescription);
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

var  crlCustomLinkView = joint.dia.LinkView.extend({
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



function crlDeleteView(evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    if (jointID) {
        var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
        var xhr = crlCreateEmptyRequest();
        var data = JSON.stringify({ "Action": "DeleteView", "RequestConceptID": diagramElementID });
        crlSendRequest(xhr, data);
    }
}

function crlPropertiesDisplayAbstractConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Abstract Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.AbstractConceptID;
}

function crlPropertiesDisplayDefinition(data, row) {
    var definitionRow = crlObtainPropertyRow(row)
    definitionRow.cells[0].innerHTML = "Definition";
    definitionRow.cells[1].innerHTML = data.NotificationConcept.Definition;
    definitionRow.cells[1].id = "definition";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        definitionRow.cells[1].contentEditable = true;
        if (!definitionRow.cells[1].callbackAssigned) {
            definitionRow.cells[1].callbackAssigned = true;
            $("#definition").on("keyup", crlSendDefinitionChanged);
        }
    } else {
        definitionRow.cells[1].contentEditable = false;
    };
}


function crlPropertiesDisplayID(data, row) {
    var idRow = crlObtainPropertyRow(row)
    idRow.cells[0].innerHTML = "ID";
    idRow.cells[1].innerHTML = data.NotificationConceptID;
}

function crlPropertiesDisplayLabel(data, row) {
    var labelRow = crlObtainPropertyRow(row);
    labelRow.cells[0].innerHTML = "Label";
    labelRow.cells[1].innerHTML = data.NotificationConcept.Label;
    labelRow.cells[1].id = "elementLabel";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        labelRow.cells[1].contentEditable = true;
        if (!labelRow.cells[1].callbackAssigned) {
            labelRow.cells[1].callbackAssigned = true;
            $("#elementLabel").on("keyup", crlSendLabelChanged);
        }
    } else {
        labelRow.cells[1].contentEditable = false;
    };
}

function crlPropertiesDisplayLiteralValue(data, row) {
    var labelRow = crlObtainPropertyRow(row);
    labelRow.cells[0].innerHTML = "Literal Value";
    labelRow.cells[1].innerHTML = data.NotificationConcept.LiteralValue;
    labelRow.cells[1].id = "literalValue";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        labelRow.cells[1].contentEditable = true;
        if (!labelRow.cells[1].callbackAssigned) {
            labelRow.cells[1].callbackAssigned = true;
            $("#literalValue").on("keyup", crlSendLiteralValueChanged);
        }
    } else {
        labelRow.cells[1].contentEditable = false;
    };
}

function crlPropertiesDisplayReferencedConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Referenced Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.ReferencedConceptID;
}

function crlPropertiesDisplayRefinedConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Refined Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.RefinedConceptID;
}

function crlPropertiesDisplayType(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Type";
    typeRow.cells[1].innerHTML = data.NotificationConcept.Type;
}

function crlPropertiesDisplayURI(data, row) {
    var uriRow = crlObtainPropertyRow(row);
    uriRow.cells[0].innerHTML = "URI";
    uriRow.cells[1].innerHTML = data.NotificationConcept.URI;
    uriRow.cells[1].id = "uri";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        uriRow.cells[1].contentEditable = true;
        if (!uriRow.cells[1].callbackAssigned) {
            uriRow.cells[1].callbackAssigned = true;
            $("#uri").on("keyup", crlSendURIChanged);
        }
    } else {
        uriRow.cells[1].contentEditable = false;
    }
}


function crlPropertiesDisplayVersion(data, row) {
    var versionRow = crlObtainPropertyRow(row)
    versionRow.cells[0].innerHTML = "Version";
    versionRow.cells[1].innerHTML = data.NotificationConcept.Version;
}

function crlDropdownMenu(dropdownId) {
    document.getElementById(dropdownId).classList.toggle("show");
}

function crlDisplayGraphDialog(numberOfAvailableGraphs) {
    crlDisplayCallGraphsDialog.open();
    $("#numberOfAvailableGraphs").text(numberOfAvailableGraphs);
}

function crlFindCellInGraph(graphID, crlJointID) {
    var cells = crlGraphsGlobal[graphID].getCells();
    var cell = null;
    cells.forEach(function (item) {
        if (item.get("crlJointID") == crlJointID) {
            cell = item;
        }
    })
    return cell
}

function crlFindCellInGraph(graph, crlJointID) {
    var cells = graph.getCells();
    var cell = null;
    cells.forEach(function (item) {
        if (item.get("crlJointID") == crlJointID) {
            cell = item;
        }
    })
    return cell
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

function crlInitializeClient() {
    crlInitializeWebSocket();
    console.log("Requesting InitializeClient");
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "InitializeClient" });
    crlSendRequest(xhr, data);
}

function crlInitializeWebSocket() {
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
                crlNotificationChangeTreeNode(data);
                break;
            case "ClearToolbarSelection":
                crlNotificationClearToolbarSelection(data);
                break;
            case "DebugSettings":
                crlNotificationSaveDebugSettings(data);
                break;
            case "DeleteDiagram":
                crlNotificationDeleteDiagram(data);
                break;
            case "DeleteDiagramElement":
                crlNotificationDeleteDiagramCell(data);
                break;
            case "DeleteTreeNode":
                crlNotificationDeleteTreeNode(data);
                break;
            case "DisplayDiagram":
                crlNotificationDisplayDiagram(data);
                break;
            case "DisplayGraph":
                crlNotificationDisplayGraph(data);
                break;
            case "EditorSettings":
                crlNotificationSaveEditorSettings(data);
                break;
            case "ElementSelected":
                crlNotificationElementSelected(data);
                break;
            case "InitializationComplete":
                crlInitializationComplete = true;
                console.log("Initialization Complete")
                crlSendNormalResponse("Processed InitializationComplete")
                break;
            case "UpdateDiagramLink":
                crlNotificationUpdateDiagramLink(data);
                break;
            case "UpdateDiagramNode":
                crlNotificationUpdateDiagramNode(data);
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
    crlSelectToolbarButton(crlCursorToolbarButtonID);
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
            crlSendReferenceLinkChanged(evt.model, linkID, sourceID, targetID);
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
            crlSendElementPointerChanged(evt.model, linkID, sourceID, targetID);
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

function crlNotificationSaveDebugSettings(data) {
    crlEnableTracing = JSON.parse(data.AdditionalParameters["EnableNotificationTracing"]);
    crlOmitHousekeepingCalls = JSON.parse(data.AdditionalParameters["OmitHousekeepingCalls"]);
    crlOmitManageTreeNodesCalls = JSON.parse(data.AdditionalParameters["OmitManageTreeNodesCalls"]);
    crlSendNormalResponse();
}

function crlMakeDiagramVisible(diagramContainerID) {
    var x = document.getElementsByClassName("crlDiagramContainer");
    for (i = 0; i < x.length; i++) {
        var container = x.item(i);
        var diagramID = crlGetDiagramIDFromDiagramContainerID(container.id);
        var tabID = crlGetDiagramTabIDFromDiagramID(diagramID);
        var tab = document.getElementById(tabID);
        if (container.id == diagramContainerID) {
            container.style.display = "block";
            tab.style.backgroundColor = "white";
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
    var concept = data.NotificationConcept;
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
    var concept = data.NotificationConcept;
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
}

function crlNotificationAddDiagramNode(data) {
    var concept = data.NotificationConcept;
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
    var concept = data.NotificationConcept;
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
    var nodeID = $('#uOfD').jstree().create_node(owningConceptID,
        {
            'id': crlGetTreeNodeIDFromConceptID(concept.ConceptID),
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

function crlNotificationChangeTreeNode(data) {
    var concept = data.NotificationConcept;
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
    crlSendNormalResponse()
}

function crlNotificationClearToolbarSelection(data) {
    crlSelectToolbarButton("cursorToolbarButton");
    crlSendNormalResponse();
}

function crlNotificationDeleteDiagram(data) {
    var diagramID = data.NotificationConceptID;
    var diagramContainerID = crlGetDiagramContainerIDFromDiagramID(diagramID);
    var diagramContainer = document.getElementById(diagramContainerID);
    if (crlCurrentDiagramContainerID == diagramContainerID) {
        crlCurrentDiagramContainerID = "";
    }
    // Remove the graph
    var jointGraphID = crlGetJointGraphIDFromDiagramID(diagramID);
    delete crlGraphsGlobal[jointGraphID]
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
    var tab = document.getElementById(tabID)
    if (tab) {
        tabs.removeChild(tab);
    }
    crlSendNormalResponse();
}

function crlNotificationDeleteDiagramCell(data) {
    var concept = data.NotificationConcept;
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
    var concept = data.NotificationConcept;
    var params = data.AdditionalParameters;
    var nodeID = crlGetTreeNodeIDFromConceptID(concept.ConceptID);
    $('#uOfD').jstree().delete_node(nodeID);
    crlSendNormalResponse();
}

function crlNotificationDisplayDiagram(data) {
    var diagramID = data.NotificationConceptID;
    var diagramLabel = data.NotificationConcept.Label;
    var diagramContainerID = crlGetDiagramContainerIDFromDiagramID(diagramID);
    var diagramContainer = document.getElementById(diagramContainerID);
    // Construct the container if it is not already present
    if (diagramContainer == undefined) {
        diagramContainer = crlConstructDiagramContainer(diagramContainer, diagramContainerID, diagramLabel, diagramID);
    }
    crlMakeDiagramVisible(diagramContainer.id);
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

function crlNotificationElementSelected(data) {
    if (data.NotificationConceptID != crlSelectedConceptIDGlobal) {
        crlSelectedConceptIDGlobal = data.NotificationConceptID
        // Update the properties
        crlPropertiesDisplayType(data, 1);
        crlPropertiesDisplayID(data, 2);
        crlPropertiesDisplayVersion(data, 3);
        crlPropertiesDisplayLabel(data, 4);
        crlPropertiesDisplayDefinition(data, 5);
        crlPropertiesDisplayURI(data, 6);
        switch (data.NotificationConcept.Type) {
            case "*core.element":
                crlPropertiesClearRow(7);
                crlPropertiesClearRow(8);
                break;
            case "*core.literal":
                crlPropertiesDisplayLiteralValue(data, 7);
                crlPropertiesClearRow(8);
                break
            case "*core.reference":
                crlPropertiesDisplayReferencedConcept(data, 7);
                crlPropertiesClearRow(8);
                break;
            case "*core.refinement":
                crlPropertiesDisplayAbstractConcept(data, 7);
                crlPropertiesDisplayRefinedConcept(data, 8);
                break;
        };

        // Update the tree
        var treeNodeID = crlGetTreeNodeIDFromConceptID(crlSelectedConceptIDGlobal);
        $("#uOfD").jstree(true).deselect_all(true);
        // a hack tp prevent infinite recursion
        crlInCrlElementSelected = true;
        $("#uOfD").jstree(true).select_node(treeNodeID, true);
        crlInCrlElementSelected = false;
    }
    crlSendNormalResponse()
}

function crlNotificationSaveEditorSettings(data) {
    crlDropReferenceAsLink = JSON.parse(data.AdditionalParameters["DropReferenceAsLink"]);
    crlDropRefinementAsLink = JSON.parse(data.AdditionalParameters["DropRefinementAsLink"]);
    crlSendNormalResponse();
}

var crlNotificationUpdateDiagramLink = function (data) {
    var concept = data.NotificationConcept;
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
    var concept = data.NotificationConcept;
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
        node.set('size', { "width": Number(params["NodeWidth"]), "height": Number(params["NodeHeight"]) });
        node.set('icon', params["Icon"]);
        node.set('name', params["DisplayLabel"]);
        node.set("abstractions", params["Abstractions"]);
    }
    crlSendNormalResponse();
}

var crlNotificationUpdateWorkspacePath = function (data) {
    crlWorkspacePath = data.AdditionalParameters["WorkspacePath"];
    crlSendNormalResponse();
}

function crlObtainPropertyRow(row) {
    var properties = document.getElementById("properties");
    var propertyRow = properties.rows[row]
    if (propertyRow == undefined) {
        propertyRow = properties.insertRow(row)
        propertyRow.insertCell(0)
        propertyRow.insertCell(1)
    }
    return propertyRow
}

var crlOnChangePosition = function (modelElement, position) {
    var jointElementID = modelElement.get("crlJointID");
    var diagramNodeID = crlGetConceptIDFromJointElementID(jointElementID);
    crlMovedNodes[diagramNodeID] = position;
    //    crlSendDiagramNodeNewPosition(diagramNodeID, position)
}

var crlOnDiagramCellPointerDown = function (cellView, event, x, y) {
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
    $.each(crlMovedNodes, function (nodeID, position) {
        crlSendDiagramNodeNewPosition(nodeID, position)
    })
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
        case "ownerLinkToolbarButton": {
            break;
        }
        case "pointerToolbarButton": {
            break;
        }
        case "abstractionLinkToolbarButton": {
            break;
        }
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
    console.log("OnDiagramManagerDrop called");
    var conceptID = crlGetConceptIDFromContainerID(event.target.parentElement.parentElement.id);
    var x = event.layerX.toString();
    var y = event.layerY.toString();
    crlSendDiagramDrop(conceptID, x, y);
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
    console.log("On Tree Drag Start called, ParentId = " + parentID);
    var selectedElementID = crlGetConceptIDFromTreeNodeID(parentID);
    console.log("selectedElementID = " + selectedElementID)
    crlSendSetTreeDragSelection(selectedElementID);
};


function crlOpenDiagramContainer(diagramContainerId) {
    var i;
    var x = document.getElementsByClassName("crlDiagramContainer");
    for (i = 0; i < x.length; i++) {
        if (x[i].id == diagramContainerId) {
            x[i].style.display = "block";
            console.log("Showing: " + diagramContainerId.toString())
        } else {
            x[i].style.display = "none";
            console.log("Hiding: " + diagramContainerId.toString())
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

function crlSendDebugSettings(enableNotificationTracing, omitHousekeepingCalls, omitManageTreeNodesCalls, maxTracingDepth) {
    var xhr = crlCreateEmptyRequest()
    var data = JSON.stringify({
        "Action": "UpdateDebugSettings",
        "AdditionalParameters": {
            "EnableNotificationTracing": enableNotificationTracing,
            "OmitHousekeepingCalls": omitHousekeepingCalls,
            "OmitManageTreeNodesCalls": omitManageTreeNodesCalls,
            "MaxTracingDepth": maxTracingDepth
        }
    });
    crlSendRequest(xhr, data);
}

function crlSendDefinitionChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "DefinitionChanged",
        "RequestConceptID": crlSelectedConceptIDGlobal,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
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

function crlSendDiagramDrop(diagramID, x, y) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "DiagramDrop",
        "AdditionalParameters":
        {
            "DiagramID": diagramID,
            "NodeX": x,
            "NodeY": y
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

function crlSendEditorSettings() {
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
        "Action": "UpdateEditorSettings",
        "AdditionalParameters": {
            "DropReferenceAsLink": dropReferenceAsLink,
            "DropRefinementAsLink": dropRefinementAsLink
        }
    });
    crlSendRequest(xhr, data);
}

function crlSendElementPointerChanged(jointLink, linkID, sourceID, targetID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "ElementPointerChanged",
        "RequestConceptID": linkID,
        "AdditionalParameters": {
            "SourceID": sourceID,
            "TargetID": targetID
        }
    })
    crlSendRequest(xhr, data);
}

function crlSendLabelChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "LabelChanged",
        "RequestConceptID": crlSelectedConceptIDGlobal,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    crlSendRequest(xhr, data)
}

function crlSendLiteralValueChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "LiteralValueChanged",
        "RequestConceptID": crlSelectedConceptIDGlobal,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    crlSendRequest(xhr, data)
}

function crlSendNewConceptSpaceRequest(evt) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "NewConceptSpaceRequest" });
    crlSendRequest(xhr, data)
}

function crlSendNewDiagramRequest(evt) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "NewDiagramRequest" });
    crlSendRequest(xhr, data)
}

function crlSendOpenWorkspace(workspacePath) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "OpenWorkspace",
        "AdditionalParameters": {
            "WorkspacePath": workspacePath
        }
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

function crlSendReferenceLinkChanged(jointLink, linkID, sourceID, targetID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "ReferenceLinkChanged",
        "RequestConceptID": linkID,
        "AdditionalParameters": {
            "SourceID": sourceID,
            "TargetID": targetID
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
        "RequestConceptID": crlSelectedConceptIDGlobal,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    crlSendRequest(xhr, data)
}

function crlSendNormalResponse() {
    var data = {};
    data["Result"] = 0;
    data["ErrorMessage"] = "none"
    crlWebsocketGlobal.send(JSON.stringify(data))
    console.log(data);
}

function crlSendSetTreeDragSelection(id) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "SetTreeDragSelection", "RequestConceptID": id });
    crlSendRequest(xhr, data);
}

function crlSendTreeNodeSelected(evt, obj) {
    if (obj != undefined) {
        var conceptID = crlGetConceptIDFromTreeNodeID(obj.node.id)
        if (conceptID != crlSelectedConceptIDGlobal && crlInCrlElementSelected == false) {
            var xhr = crlCreateEmptyRequest();
            var data = JSON.stringify({ "Action": "TreeNodeSelected", "RequestConceptID": conceptID });
            crlSendRequest(xhr, data);
        }
    };
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

var crlShowOwner = function (evt) {
    var cellView = crlDiagramCellDropdownMenu.attributes.cellView;
    var jointID = cellView.model.attributes.crlJointID;
    var diagramElementID = crlGetConceptIDFromJointElementID(jointID)
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "ShowOwner", "RequestConceptID": diagramElementID });
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
        case "Refinement":
            return targetRepresents == "Element" || targetRepresents == "Literal" || targetRepresents == "Refinement" || targetRepresents == "Reference";
        default:
            // Must be a pointer
            return targetRepresents == "Element" || targetRepresents == "Literal" || targetRepresents == "Refinement" || targetRepresents == "Reference";
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

