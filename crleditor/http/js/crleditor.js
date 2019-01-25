var canvas
var currentDiagramContainerID
var graphs = {}
var initializationComplete = false
var papers = {}

// <!-- Set css parameters -->
$(function () {
    $(".uofd-browser").resizable({
        handles: "e",
        resize: sizeAll()
    });
    $(".bottom").resizable({
        handles: "n",
        resize: sizeAll()
    });
});

// Initialize
$(function () {
    $("#uOfD").jstree({
        'core': {
            'check_callback': true,
            'multiple': false
        },
        'plugins': ['sort', 'contextmenu', 'unique', 'wholerow'],
        'sort': function (a, b) {
            a1 = this.get_node(a);
            b1 = this.get_node(b);
            return (a1.text > b1.text) ? 1 : -1;
        },
        'contextmenu': {
            "items": function ($node) {
                var tree = $("uOfD").jstree(true);
                var items = {
                    display: {
                        "label": "Display Diagram",
                        "action": function (obj) {
                            if ($node != undefined) {
                                var xhr = createEmptyRequest();
                                var conceptID = getConceptIDFromTreeNodeID($node.id)
                                var data = JSON.stringify({ "Action": "DisplayDiagramSelected", "RequestConceptID": conceptID });
                                xhr.send(data);
                            }
                        }
                    },
                    remove: {
                        "label": "Delete",
                        "action": function (obj) {
                            if ($node != undefined) {
                                var xhr = createEmptyRequest();
                                var conceptID = getConceptIDFromTreeNodeID($node.id)
                                var data = JSON.stringify({ "Action": "TreeNodeDelete", "RequestConceptID": conceptID });
                                xhr.send(data);
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
    $("#uOfD").on("select_node.jstree", sendTreeNodeSelected);
    $("#uOfD").on("dragstart", onTreeDragStart);
    $("#body").on("ondrop", onEditorDrop)
    canvas = document.createElement("canvas");
});


// <!-- Display an example graph - this is throw-away code -->
$(function () {
    var graph = new joint.dia.Graph;

    var paper = new joint.dia.Paper({
        el: $('#diagram'),
        /*			width : 800,
         height : 800, */
        model: graph,
        gridSize: 1
    });

    var rect = new joint.shapes.basic.Rect({
        position: {
            x: 100,
            y: 30
        },
        size: {
            width: 100,
            height: 30
        },
        attrs: {
            rect: {
                fill: 'blue'
            },
            text: {
                text: 'my box',
                fill: 'white'
            }
        }
    });

    var rect3 = new joint.shapes.uml.Class({
        position: {
            x: 500,
            y: 30
        },
        size: {
            width: 100,
            height: 30
        },
        name: `Bar`
    });

    theClass = rect3

    graph.addCells([rect, rect3]);
});


var websocket
var selectedConceptID
var treeDragSelectionID

function addDiagramNode(data) {
    // Now construct the jointjs representation
    containerID = getContainerIDFromConceptID(data.AdditionalParameters.OwnerID);
    graphID = getJointGraphIDFromDiagramID(data.AdditionalParameters.OwnerID);
    jointElementID = getJointElementIDFromConceptID(data.NotificationConceptID);
    jointElement = new joint.shapes.crl.Element(
        {
            attrs: {
                rect: {
                    width: 300
                },
                '.image': {
                    'ref-x': 1.0,
                    'ref-y': 1.0,
                    ref: ".label-rect",
                    width: 16,
                    height: 16
                },
                '.label-rect': {
                    stroke: "black",
                    'stroke-width': 2,
                    fill: "#ffffff",
                    height: 40,
                    transform: "translate(0,0)"
                },
                '.abstractions-text': {
                    ref: ".label-rect",
                    'ref-y': 0.5,
                    'ref-x': 0.5 + 18,
                    'text-anchor': "right",
                    'y-alignment': "middle",
                    'font-weight': "normal",
                    'font-style': "italic",
                    fill: "black",
                    'font-size': 12,
                    'font-family': "Go,  Helvetica, Ariel, sans-serif"
                },
                '.label-text': {
                    ref: ".label-rect",
                    'ref-y': 0.5,
                    'ref-x': 0.5 + 18,
                    'text-anchor': "left",
                    'y-alignment': "middle",
                    'font-weight': "bold",
                    fill: "black",
                    'font-size': 12,
                    'font-family': "Go,  Helvetica, Ariel, sans-serif"
                }
            }
        },
        {
            markup: "<g class=\"rotatable\">" +
                "<g class=\"scalable\">" +
                "<rect class=\"label-rect\"/>" +
                "</g>" +
                "<image class=\"image\"/><text class=\"abstractions-text\"/><text class=\"label-text\"/>" +
                "</g>",
            initialize: function () {
                // js.Global.Get("joint").Get("shapes").Get("basic").Get("Generic").Get("prototype").Get("initialize").Call("apply", this, arguments)
                // this.Call("updateRectangles")
                return nil
            },
            updateRectangles: function () {
                offsetY = 0
                attributes = this.Get("attributes")
                attrs = attributes.Get("attrs")

                rectHeight = 1 * 12 + 6
                labelText = attributes.Get("name")
                labelTextAttr = attrs.Get(".label-text")
                labelTextAttr.Set("text", labelText)
                labelRectAttr = attrs.Get(".label-rect")
                labelRectAttr.Set("height", rectHeight)
                rectWidth = calculateTextWidth(labelText.String()) + 6 + 18
                labelRectAttr.Set("transform", "translate(0," + strconv.Itoa(offsetY) + ")")
                this.Call("resize", rectWidth, rectHeight)

                offsetY += rectHeight
                return nil
            }
        })

    jointElement.crlJointId = jointElementID;
    jointElement.attributes.name = data.AdditionalParameters["DisplayLabel"];
    jointElement.attributes.position = { "x": data.AdditionalParameters["NodeX"], "y": data.AdditionalParameters["NodeY"] };
    jointElement.attributes.attrs.image = { "xlink:href": data.AdditionalParameters["Icon"] };

    jointElement.updateRectangles();
    graphs[graphID].addCell(jointElement);
    sendNormalResponse();
}

// <!-- Set up the websockets connection and callbacks -->
function addTreeNode(data) {
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
            'id': "TreeNode" + concept.ConceptID,
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
    sendNormalResponse();
}

var calculateTextWidth = function (text) {
    return getTextWidth(text, "go12PtBoldFace")
}

function callExit() {
    console.log("Requesting Exit");
    var xhr = createEmptyRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            window.close();
        }
    };
    var data = JSON.stringify({ "Action": "Exit" });
    xhr.send(data);
}

function changeTreeNode(data) {
    var concept = data.NotificationConcept;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    treeNodeOwnerID = ""
    if (owningConceptID == "") {
        treeNodeOwnerID = "#";
    } else {
        treeNodeOwnerID = getTreeNodeIDFromConceptID(owningConceptID);
    }
    nodeID = getTreeNodeIDFromConceptID(concept.ConceptID);
    if ($('#uOfD').jstree().get_parent(nodeID) != treeNodeOwnerID) {
        $('#uOfD').jstree().move_node(nodeID, treeNodeOwnerID)
    }
    $('#uOfD').jstree().rename_node(nodeID, concept.Label)
    sendNormalResponse()
}

function clearRow(row) {
    properties = document.getElementById("properties");
    propertyRow = properties.rows[row]
    if (propertyRow != undefined) {
        properties.deleteRow(row);
    }
}


var closeWebsocket = function () {
    console.log("Closing websocket")
    websocket.close()
}

function createEmptyRequest() {
    var xhr = new XMLHttpRequest();
    var url = "request";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            console.log(xhr.responseText)
        };
    }
    return xhr
}

function deleteTreeNode(data) {
    var concept = data.NotificationConcept;
    var params = data.AdditionalParameters;
    nodeID = getTreeNodeIDFromConceptID(concept.ConceptID);
    $('#uOfD').jstree().delete_node(nodeID);
    var data = {};
    data["Result"] = 0;
    data["ErrorMessage"] = "none"
    websocket.send(JSON.stringify(data))
}

function displayAbstractConcept(data, row) {
    typeRow = obtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Abstract Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.AbstractConceptID;
}

function displayDefinition(data, row) {
    definitionRow = obtainPropertyRow(row)
    definitionRow.cells[0].innerHTML = "Definition";
    definitionRow.cells[1].innerHTML = data.NotificationConcept.Definition;
    definitionRow.cells[1].id = "definition";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        definitionRow.cells[1].contentEditable = true;
        $("#definition").on("keyup", sendDefinitionChanged);
    } else {
        definitionRow.cells[1].contentEditable = false;

    };
}


function displayDiagram(data) {
    diagramID = data.NotificationConceptID;
    diagramLabel = data.NotificationConcept.Label;
    diagramContainerID = getDiagramContainerIDFromDiagramID(diagramID);
    diagramContainer = document.getElementById(diagramContainerID);
    // Construct the container if it is not already present
    if (diagramContainer == undefined) {
        topContent = document.getElementById("top-content");
        diagramContainer = document.createElement("DIV");
        diagramContainer.id = diagramContainerID;
        diagramContainer.className = "crlDiagramContainer";
        // It is not clear why, but the ondrop callback does not get called unless the ondragover callback is used,
        // even though the callback just calls preventDefault on the dragover event
        diagramContainer.ondragover = onDragover;
        diagramContainer.ondrop = onDiagramDrop;
        diagramContainer.style.display = "none";
        topContent.appendChild(diagramContainer);
        // Create the new tab
        tabs = document.getElementById("tabs");
        newTab = document.createElement("BUTTON");
        newTab.innerHTML = diagramLabel;
        newTab.className = "w3-bar-item w3-button";
        newTabID = "DiagramTab" + diagramID;
        newTab.id = newTabID;
        newTab.setAttribute("diagramContainerID", diagramContainerID);
        newTab.onclick = onMakeDiagramVisible;
        tabs.appendChild(newTab, -1);

        jointGraphID = getJointGraphIDFromDiagramID(diagramID);
        jointGraph = graphs[jointGraphID];
        //        jointGraph = document.getElementById(jointGraphID);
        if (jointGraph == undefined) {
            jointGraph = new joint.dia.Graph;
            jointGraph.id = jointGraphID;
            graphs[jointGraphID] = jointGraph;
        };

        jointPaperID = getJointPaperIDFromDiagramID(diagramID);
        jointPaper = papers[jointPaperID];
        if (jointPaper == undefined) {
            diagramPaperDiv = document.createElement("DIV");
            diagramContainer.appendChild(diagramPaperDiv);
            jointPaper = new joint.dia.Paper({
                "el": diagramPaperDiv,
                "width": 1000,
                "height": 1000,
                "model": jointGraph,
                "gridSize": 1
            });
            jointPaper.on("cell:pointerdown", onDiagramCellPointerDown);
            papers[jointPaperID] = jointPaper;
        };
    }
    makeDiagramVisible(diagramContainer.id);
    currentDiagramContainerID = diagramContainerID
    sendNormalResponse()
}


function displayID(data, row) {
    idRow = obtainPropertyRow(row)
    idRow.cells[0].innerHTML = "ID";
    idRow.cells[1].innerHTML = data.NotificationConceptID;
}

function displayLabel(data, row) {
    labelRow = obtainPropertyRow(row);
    labelRow.cells[0].innerHTML = "Label";
    labelRow.cells[1].innerHTML = data.NotificationConcept.Label;
    labelRow.cells[1].id = "elementLabel";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        labelRow.cells[1].contentEditable = true;
        $("#elementLabel").on("keyup", sendLabelChanged);
    } else {
        labelRow.cells[1].contentEditable = false;
    };
}

function displayLiteralValue(data, row) {
    labelRow = obtainPropertyRow(row);
    labelRow.cells[0].innerHTML = "Literal Value";
    labelRow.cells[1].innerHTML = data.NotificationConcept.LiteralValue;
    labelRow.cells[1].id = "literalValue";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        labelRow.cells[1].contentEditable = true;
        $("#literalValue").on("keyup", sendLiteralValueChanged);
    } else {
        labelRow.cells[1].contentEditable = false;
    };
}

function displayReferencedConcept(data, row) {
    typeRow = obtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Referenced Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.ReferencedConceptID;
}

function displayRefinedConcept(data, row) {
    typeRow = obtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Refined Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.RefinedConceptID;
}

function displayType(data, row) {
    typeRow = obtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Type";
    typeRow.cells[1].innerHTML = data.NotificationConcept.Type;
}

function displayURI(data, row) {
    uriRow = obtainPropertyRow(row);
    uriRow.cells[0].innerHTML = "URI";
    uriRow.cells[1].innerHTML = data.NotificationConcept.URI;
    uriRow.cells[1].id = "uri";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        uriRow.cells[1].contentEditable = true;
        $("#uri").on("keyup", sendURIChanged);
    } else {
        uriRow.cells[1].contentEditable = false;
    }
}


function displayVersion(data, row) {
    versionRow = obtainPropertyRow(row)
    versionRow.cells[0].innerHTML = "Version";
    versionRow.cells[1].innerHTML = data.NotificationConcept.Version;
}

function dropdownMenu(dropdownId) {
    document.getElementById(dropdownId).classList.toggle("show");
}

function elementSelected(data) {
    if (data.NotificationConceptID != selectedConceptID) {
        selectedConceptId = data.NotificationConceptID

        // Update the properties
        displayType(data, 1);
        displayID(data, 2);
        displayVersion(data, 3);
        displayLabel(data, 4);
        displayDefinition(data, 5);
        displayURI(data, 6);
        switch (data.NotificationConcept.Type) {
            case "*core.element":
                clearRow(7);
                clearRow(8);
                break;
            case "*core.literal":
                displayLiteralValue(data, 7);
                clearRow(8);
                break
            case "*core.reference":
                displayReferencedConcept(data, 7);
                clearRow(8);
                break;
            case "*core.refinement":
                displayAbstractConcept(data, 7);
                displayRefinedConcept(data, 8);
                break;
        };
    }
    sendNormalResponse()
}

function getConceptIDFromContainerID(containerID) {
    return containerID.replace("DiagramContainer", "")
}

function getConceptIDFromTreeNodeID(treeNodeID) {
    return treeNodeID.replace("TreeNode", "");
}

function getContainerIDFromConceptID(conceptID) {
    return "DiagramContainer" + conceptID;
}

function getDiagramContainerIDFromDiagramID(diagramID) {
    return "DiagramContainer" + diagramID;
}

function getDiagramIDFromDiagramContainerID(diagramContainerID) {
    return diagramContainerID.replace("DiagramContainer", "");
}

function getDiagramIDFromJointGraphID(jointGraphID) {
    return jointGraphID.replace("JointGraph", "");
}

function getDiagramIDFromJointPaperID(jointPaperID) {
    return jointPaperID.replace("JointPaper", "")
}

function getJointPaperIDFromDiagramID(diagramID) {
    return "JointPaper" + diagramID;
}

function getJointGraphIDFromDiagramID(diagramID) {
    return "JointGraph" + diagramID;
}

function getTextWidth(text, font) {
    var context = canvas.getContext("2d");
    context.font = font;
    var metrics = context.measureText(text);
    return metrics.width;
}

function getTreeNodeIDFromConceptID(conceptID) {
    return "TreeNode" + conceptID;
}


function initializeClient() {
    initializeWebSocket();
    console.log("Requesting InitializeClient");
    var xhr = createEmptyRequest();
    var data = JSON.stringify({ "Action": "InitializeClient" });
    xhr.send(data);
}

function initializeWebSocket() {
    console.log("Initializing Web Socket")
    // ws initialization
    websocket = new WebSocket("ws://localhost:8080/index/ws");
    console.log("Web Socket Initialization complete")
    websocket.onmessage = function (e) {
        var data = JSON.parse(e.data)
        console.log("Notification:" + data.Notification)
        switch (data.Notification) {
            case 'AddDiagramNode':
                addDiagramNode(data);
                break;
            case 'AddTreeNode':
                addTreeNode(data);
                break;
            case "ChangeTreeNode":
                changeTreeNode(data);
                break;
            case "DeleteTreeNode":
                deleteTreeNode(data);
                break;
            case "DisplayDiagram":
                displayDiagram(data);
                break;
            case "ElementSelected":
                elementSelected(data);
                break;
            case "InitializationComplete":
                initializationComplete = true;
                console.log("Initialization Complete")
                sendNormalResponse("Processed InitializationComplete")
                break;
            default:
                console.log('Unhandled notification: ' + e.data);
                var data = {};
                data["Result"] = 1;
                data["ErrorMessage"] = "Unhandled notification";
                websocket.send(JSON.stringify(data));
        }
    };
};

function getConceptIDFromJointElementID(jointElementID) {
    return jointElementID.replace("JointElement", "")
}

function getJointElementIDFromConceptID(conceptID) {
    return "JointElement" + conceptID
}

function makeDiagramVisible(diagramContainerID) {
    x = document.getElementsByClassName("crlDiagramContainer");
    for (i = 0; i < x.length; i++) {
        if (x.item(i).id == diagramContainerID) {
            x.item(i).style.display = "block";
        } else {
            x.item(i).style.display = "none";
        }
    }
}


function obtainPropertyRow(row) {
    properties = document.getElementById("properties");
    propertyRow = properties.rows[row]
    if (propertyRow == undefined) {
        propertyRow = properties.insertRow(row)
        propertyRow.insertCell(0)
        propertyRow.insertCell(1)
    }
    return propertyRow
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

function onDiagramCellPointerDown(cellView, event, x, y) {
    jointElementID = cellView.model.crlJointId;
    diagramNodeID = getConceptIDFromJointElementID(jointElementID);
    if (diagramNodeID == "") {
        console.log("In onDiagramManagerCellPointerDown diagramNodeID is empty")
    }
    sendDiagramNodeSelected(diagramNodeID)
}

function onDiagramDrop(event) {
    event.preventDefault();
    console.log("OnDiagramManagerDrop called");
    conceptID = getConceptIDFromContainerID(event.target.parentElement.parentElement.id);
    x = event.layerX.toString();
    y = event.layerY.toString();
    sendDiagramDrop(conceptID, x, y);
}

function onDragover(event, data) {
    event.preventDefault()
}

function onEditorDrop(e, data) {
    sendSetTreeDragSelection("")
}


function onMakeDiagramVisible(e) {
    diagramContainerId = e.target.getAttribute("diagramContainerID")
    makeDiagramVisible(diagramContainerId)
}


function onTreeDragStart(e, data) {
    parentID = e.target.parentElement.id;
    console.log("On Tree Drag Start called, ParentId = " + parentID);
    selectedElementID = getConceptIDFromTreeNodeID(parentID);
    console.log("selectedElementID = " + selectedElementID)
    sendSetTreeDragSelection(selectedElementID);
}


function openDiagramContainer(diagramContainerId) {
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

function sendDefinitionChanged(evt, obj) {
    xhr = createEmptyRequest();
    data = JSON.stringify({
        "Action": "DefinitionChanged",
        "RequestConceptID": selectedConceptId,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    xhr.send(data)
}

function sendLabelChanged(evt, obj) {
    xhr = createEmptyRequest();
    data = JSON.stringify({
        "Action": "LabelChanged",
        "RequestConceptID": selectedConceptId,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    xhr.send(data)
}

function sendLiteralValueChanged(evt, obj) {
    xhr = createEmptyRequest();
    data = JSON.stringify({
        "Action": "LiteralValueChanged",
        "RequestConceptID": selectedConceptId,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    xhr.send(data)
}

function sendURIChanged(evt, obj) {
    xhr = createEmptyRequest();
    data = JSON.stringify({
        "Action": "URIChanged",
        "RequestConceptID": selectedConceptId,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    xhr.send(data)
}

function sendNewDiagramRequest(evt) {
    xhr = createEmptyRequest();
    data = JSON.stringify({ "Action": "NewDiagramRequest" });
    xhr.send(data)
}

function sendDiagramDrop(diagramID, x, y) {
    console.log("Diagram Drop ID, x, y: " + diagramID + "  " + x.toString() + "  " +  y.toString());
    var xhr = createEmptyRequest();
    var data = JSON.stringify({
        "Action": "DiagramDrop",
        "AdditionalParameters":
        {
            "DiagramID": diagramID,
            "NodeX": x,
            "NodeY": y
        }
    });
    console.log(data);
    xhr.send(data);
}

function sendDiagramNodeSelected(nodeID) {
    var xhr = createEmptyRequest();
    var data = JSON.stringify({ "Action": "DiagramNodeSelected", "RequestConceptID": nodeID });
    xhr.send(data);
}

function sendNormalResponse() {
    var data = {};
    data["Result"] = 0;
    data["ErrorMessage"] = "none"
    websocket.send(JSON.stringify(data))
    console.log('Sent normal response');
}

function sendSetTreeDragSelection(id) {
    var xhr = createEmptyRequest();
    var data = JSON.stringify({ "Action": "SetTreeDragSelection", "RequestConceptID": id });
    xhr.send(data);
}

function sendTreeNodeSelected(evt, obj) {
    if (obj != undefined) {
        var xhr = createEmptyRequest();
        var conceptID = getConceptIDFromTreeNodeID(obj.node.id)
        var data = JSON.stringify({ "Action": "TreeNodeSelected", "RequestConceptID": conceptID });
        xhr.send(data);
    };
}

// <!-- Define sizeAll() to manage sizes of display components -->
function sizeAll() {
    // Body
    var bodyHeight = jQuery("body").height();
    var bodyWidth = jQuery("body").width();
    // Wrapper
    var wrapperMargin = jQuery("#wrapper").outerWidth(true)
        - jQuery("#wrapper").width();
    var wrapperWidth = bodyWidth - wrapperMargin;
    var wrapperHeight = bodyHeight - jQuery("#navbar").outerHeight(true)
        - wrapperMargin;
    jQuery("#wrapper").width(wrapperWidth);
    jQuery("#wrapper").height(wrapperHeight);
    // UofDBrowser
    var uOfDBrowserOuterWidth = jQuery("#uofd-browser").outerWidth(true);
    var uOfDBrowserMargin = uOfDBrowserOuterWidth
        - jQuery("#uofd-browser").width();
    var uOfDBrowserHeight = wrapperHeight - uOfDBrowserMargin;
    jQuery("#uofd-browser").height(uOfDBrowserHeight);
    // Center Pane
    var centerPaneMargin = jQuery("center-pane").outerWidth(true)
        - jQuery("center-pane").width();
    var centerPaneWidth = wrapperWidth - uOfDBrowserOuterWidth
        - centerPaneMargin;
    jQuery("#center-pane").width(centerPaneWidth);
    var centerPaneHeight = wrapperHeight - centerPaneMargin;
    jQuery("#center-pane").height(centerPaneHeight);
    jQuery("#center-pane").position({
        my: "left top",
        at: "right top",
        of: "#uofd-browser"
    })

    // Top Pane
    var topPaneMargin = jQuery("#top-pane").outerWidth(true)
        - jQuery("#top-pane").width();
    jQuery("#top-pane").height(
        centerPaneHeight - jQuery("#bottom").outerHeight(true)
        - topPaneMargin);

    // Toolbar
    var toolbarOuterWidth = jQuery("#toolbar").outerWidth(true);
    var toolbarMargin = toolbarOuterWidth - jQuery("#toolbar").width();
    var toolbarHeight = jQuery("#top-pane").height() - toolbarMargin;
    jQuery("#toolbar").height(toolbarHeight);

    // Top Content
    var topContentMargin = jQuery("#top-content").outerWidth(true)
        - jQuery("#top-content").width();
    jQuery("#top-content").width(
        centerPaneWidth - toolbarOuterWidth - topContentMargin);
    topContentHeight = jQuery("#top-pane").height() - topContentMargin;
    jQuery("#top-caontent").height(topContentHeight);

    // crlDiagramContainer
    crlDiagramContainerHeight = topContentHeight - jQuery("#tabs").height();
    jQuery(".crlDiagramContainer").height(crlDiagramContainerHeight);

    // Bottom
    jQuery("#bottom").position({
        my: "left bottom",
        at: "right bottom",
        of: "#uofd-browser",
        collision: "fit"
    });
};


