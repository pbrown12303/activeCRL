// crlCurrentDiagramContainerIDGlobal is the identifier for the diagram container currently being displayed
var crlCurrentDiagramContainerIDGlobal
// crlDebugSettingsDialog is the initialized dialog used for editing debug settings
var crlDebugSettingsDialog
// crlDropReferenceAsLink when true causes references dragged from the tree into the diagram to be added as links
var crlDropReferenceAsLink = false
// crlDropRefinementAsLink when true causes refinements dragged from the tree into the diagram to be added as links
var crlDropRefinementAsLink = false
// crlEditorSettingsDialog is the initialized dialog used for editing editor settings
var crlEditorSettingsDialog
// crlEnableTracing is the client-side copy of the server-side value that turns on notification tracing
var crlEnableTracing = false
// crlGraphsGlobal is an array of existing graphs that is used to look up a graph given its identifier
var crlGraphsGlobal = {}
// crlInitializationCompleteGlobal indicates whether the server-side initialization has been completed
var crlInitializationCompleteGlobal = false
// crlMovedNodes is an array of nodes that have been moved. This is a temporary cache that is used to update the 
// server once a mouse up has occurred
var crlMovedNodes = {}
// crlOpenWorkspaceDialog is the initialized dialog used for opening a workspace
var crlOpenWorkspaceDialog
// crlPaperGlobal is an array of existing papers that is used to look up a paper given its identifier
var crlPapersGlobal = {}
// crlSelectedConceptIDGlobal contains the model identifier of the currently selected concept
var crlSelectedConceptIDGlobal
// crlTreeDragSelectionIDGlobal contains the model identifier of the concept currently being dragged from the tree
var crlTreeDragSelectionIDGlobal
// CrlWebSocketGlobal is the web socket being used for server-side communications
var crlWebsocketGlobal
// crlWorkspacePath is the path to the current workspace
var crlWorkspacePath

// <!-- Set css parameters -->
$(function () {
    $(".uofd-browser").resizable({
        handles: "e",
        resize: crlSizeAll()
    });
    $(".bottom").resizable({
        handles: "n",
        resize: crlSizeAll()
    });
});

// Initialize
$(function () {
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
                                        xhr.send(data);
                                    }
                                }
                            },
                            Diagram: {
                                "label": "Diagram",
                                "action": function (obj) {
                                    var xhr = crlCreateEmptyRequest();
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                    var data = JSON.stringify({ "Action": "AddDiagramChild", "RequestConceptID": conceptID });
                                    xhr.send(data);
                            }
                            },
                            Literal: {
                                "label": "Literal",
                                "action": function (obj) {
                                    var xhr = crlCreateEmptyRequest();
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                    var data = JSON.stringify({ "Action": "AddLiteralChild", "RequestConceptID": conceptID });
                                    xhr.send(data);
                           }
                            },
                            Reference: {
                                "label": "Reference",
                                "action": function (obj) {
                                    var xhr = crlCreateEmptyRequest();
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                    var data = JSON.stringify({ "Action": "AddReferenceChild", "RequestConceptID": conceptID });
                                    xhr.send(data);
                           }
                            },
                            Refinement: {
                                "label": "Refinement",
                                "action": function (obj) {
                                    var xhr = crlCreateEmptyRequest();
                                    var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                    var data = JSON.stringify({ "Action": "AddRefinementChild", "RequestConceptID": conceptID });
                                    xhr.send(data);
                            }
                            }
                        }
                    },
                    display: {
                        "label": "Display Diagram",
                        "action": function (obj) {
                            if ($node != undefined) {
                                var xhr = crlCreateEmptyRequest();
                                var conceptID = crlGetConceptIDFromTreeNodeID($node.id)
                                var data = JSON.stringify({ "Action": "DisplayDiagramSelected", "RequestConceptID": conceptID });
                                xhr.send(data);
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
    $("#uOfD").on("select_node.jstree", crlSendTreeNodeSelected);
    $("#uOfD").on("dragstart", crlOnTreeDragStart);
    $("#body").on("ondrop", crlOnEditorDrop);
    crlDebugSettingsDialog = $("#debugSettingsDialog").dialog({
        "resizable": false,
        "height": 200,
        "modal": true,
        "buttons": { "OK": crlDebugSettingsOK }
    });
    crlDebugSettingsDialog.dialog("close");
    crlEditorSettingsDialog = $("#editorSettingsDialog").dialog({
        "resizable": false,
        "height": 200,
        "modal": true,
        "buttons": { "OK": crlEditorSettingsOK }
    });
    crlEditorSettingsDialog.dialog("close");
    crlOpenWorkspaceDialog = $("#openWorkspaceDialog").dialog({
        "resizable": true,
        "width": 650,
        "height": 300,
        "modal": true,
        "buttons": { "OK": crlOpenWorkspaceOK }
    });
    crlOpenWorkspaceDialog.dialog("close");
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



function crlAddDiagramLink(data) {
    crlUpdateDiagramLink(data);
}

function crlAddDiagramNode(data) {
    crlUpdateDiagramNode(data);
}

function crlConstructDiagramLink(data, graphID, crlJointID) {
    var sourceJointID = crlGetJointCellIDFromConceptID(data.AdditionalParameters["LinkSourceID"])
    var targetJointID = crlGetJointCellIDFromConceptID(data.AdditionalParameters["LinkTargetID"])
    if (sourceJointID != "" && targetJointID != "") {
        var linkSource = crlFindElementInGraph(graphID, sourceJointID)
        var linkTarget = crlFindElementInGraph(graphID, targetJointID)
        if (linkSource != undefined && linkTarget != undefined) {
            var newLink;
            switch (data.AdditionalParameters["LinkType"]) {
                case "*core.reference":
                    newLink = new joint.shapes.crl.ReferenceLink({
                        source: { id: linkSource.id },
                        target: { id: linkTarget.id }
                    });
                    break;
                case "*core.refinement":
                    var newLink = new joint.shapes.crl.RefinementLink({
                        source: { id: linkSource.id },
                        target: { id: linkTarget.id }
                    });
                    break;
            }
            newLink.set("crlJointID", crlJointID);
            crlGraphsGlobal[graphID].addCell(newLink);
            return newLink;
        }
    }
    return undefined
}

function crlConstructDiagramNode(data, graphID, crlJointID) {
    var jointElement = new joint.shapes.crl.Element({
        attrs: {
            rect: {
                width: 300
            }
        }
    });
    jointElement.set("crlJointID", crlJointID);
    jointElement.set("name", data.AdditionalParameters["DisplayLabel"]);
    jointElement.set("position", { "x": Number(data.AdditionalParameters["NodeX"]), "y": Number(data.AdditionalParameters["NodeY"]) });
    jointElement.set("size", { "width": Number(data.AdditionalParameters["NodeWidth"]), "height": Number(data.AdditionalParameters["NodeHeight"]) });
    jointElement.set("icon", data.AdditionalParameters["Icon"]);
    jointElement.set("abstractions", data.AdditionalParameters["Abstractions"]);
    var graph = crlGraphsGlobal[graphID];
    graph.addCell(jointElement);
    return jointElement;
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

function crlFindElementInGraph(graphID, crlJointID) {
    var elements = crlGraphsGlobal[graphID].getElements();
    var elem = null;
    elements.forEach(function (item) {
        if (item.get("crlJointID") == crlJointID) {
            elem = item;
        }
    })
    return elem
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

// <!-- Set up the websockets connection and callbacks -->
function crlAddTreeNode(data) {
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
    crlSendNormalResponse();
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
    xhr.send(data);
}

function crlChangeTreeNode(data) {
    var concept = data.NotificationConcept;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var treeNodeOwnerID = ""
    if (owningConceptID == "") {
        treeNodeOwnerID = "#";
    } else {
        treeNodeOwnerID = crlGetTreeNodeIDFromConceptID(owningConceptID);
    }
    var nodeID = crlGetTreeNodeIDFromConceptID(concept.ConceptID);
    if ($('#uOfD').jstree().get_parent(nodeID) != treeNodeOwnerID) {
        $('#uOfD').jstree().move_node(nodeID, treeNodeOwnerID)
    }
    $('#uOfD').jstree().rename_node(nodeID, concept.Label)
    crlSendNormalResponse()
}

function crlClearRow(row) {
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
            console.log(xhr.responseText)
        };
    }
    return xhr
}

// crlDebugSettings creates and displays the Debug Settings dialog so that the debug settings can be updated from the UI
var crlDebugSettings = function () {
    $("#enableTracing").prop("checked", crlEnableTracing);
    crlDebugSettingsDialog.dialog("open");
}

// crlDebugSettingsOK is the callback function for the Debug Settings dialog OK button.
// It updaates the debug settings with the values from the dialog
var crlDebugSettingsOK = function () {
    crlSendDebugSettings();
    crlDebugSettingsDialog.dialog("close");
};


function crlDeleteTreeNode(data) {
    var concept = data.NotificationConcept;
    var params = data.AdditionalParameters;
    var nodeID = crlGetTreeNodeIDFromConceptID(concept.ConceptID);
    $('#uOfD').jstree().delete_node(nodeID);
    var data = {};
    data["Result"] = 0;
    data["ErrorMessage"] = "none"
    crlWebsocketGlobal.send(JSON.stringify(data))
}

function crlDisplayAbstractConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Abstract Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.AbstractConceptID;
}

function crlDisplayDefinition(data, row) {
    var definitionRow = crlObtainPropertyRow(row)
    definitionRow.cells[0].innerHTML = "Definition";
    definitionRow.cells[1].innerHTML = data.NotificationConcept.Definition;
    definitionRow.cells[1].id = "definition";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        definitionRow.cells[1].contentEditable = true;
        $("#definition").on("keyup", crlSendDefinitionChanged);
    } else {
        definitionRow.cells[1].contentEditable = false;

    };
}


function crlDisplayDiagram(data) {
    var diagramID = data.NotificationConceptID;
    var diagramLabel = data.NotificationConcept.Label;
    var diagramContainerID = crlGetDiagramContainerIDFromDiagramID(diagramID);
    var diagramContainer = document.getElementById(diagramContainerID);
    // Construct the container if it is not already present
    if (diagramContainer == undefined) {
        var topContent = document.getElementById("top-content");
        diagramContainer = document.createElement("DIV");
        diagramContainer.id = diagramContainerID;
        diagramContainer.className = "crlDiagramContainer";
        // It is not clear why, but the ondrop callback does not get called unless the ondragover callback is used,
        // even though the callback just calls preventDefault on the dragover event
        diagramContainer.ondragover = crlOnDragover;
        diagramContainer.ondrop = crlOnDiagramDrop;
        diagramContainer.style.display = "none";
        topContent.appendChild(diagramContainer);
        // Create the new tab
        var tabs = document.getElementById("tabs");
        var newTab = document.createElement("BUTTON");
        newTab.innerHTML = diagramLabel;
        newTab.className = "w3-bar-item w3-button";
        var newTabID = "DiagramTab" + diagramID;
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
        };

        var jointPaperID = crlGetJointPaperIDFromDiagramID(diagramID);
        var jointPaper = crlPapersGlobal[jointPaperID];
        if (jointPaper == undefined) {
            var diagramPaperDiv = document.createElement("DIV");
            diagramContainer.appendChild(diagramPaperDiv);
            jointPaper = new joint.dia.Paper({
                "el": diagramPaperDiv,
                "width": 1000,
                "height": 1000,
                "model": jointGraph,
                "gridSize": 1
            });
            jointPaper.on("cell:pointerdown", crlOnDiagramCellPointerDown);
            jointPaper.on("cell:pointerup", crlOnDiagramCellPointerUp);
            crlPapersGlobal[jointPaperID] = jointPaper;
        };
    }
    crlMakeDiagramVisible(diagramContainer.id);
    crlCurrentDiagramContainerIDGlobal = diagramContainerID
    crlSendNormalResponse()
}


function crlDisplayID(data, row) {
    var idRow = crlObtainPropertyRow(row)
    idRow.cells[0].innerHTML = "ID";
    idRow.cells[1].innerHTML = data.NotificationConceptID;
}

function crlDisplayLabel(data, row) {
    var labelRow = crlObtainPropertyRow(row);
    labelRow.cells[0].innerHTML = "Label";
    labelRow.cells[1].innerHTML = data.NotificationConcept.Label;
    labelRow.cells[1].id = "elementLabel";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        labelRow.cells[1].contentEditable = true;
        $("#elementLabel").on("keyup", crlSendLabelChanged);
    } else {
        labelRow.cells[1].contentEditable = false;
    };
}

function crlDisplayLiteralValue(data, row) {
    var labelRow = crlObtainPropertyRow(row);
    labelRow.cells[0].innerHTML = "Literal Value";
    labelRow.cells[1].innerHTML = data.NotificationConcept.LiteralValue;
    labelRow.cells[1].id = "literalValue";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        labelRow.cells[1].contentEditable = true;
        $("#literalValue").on("keyup", crlSendLiteralValueChanged);
    } else {
        labelRow.cells[1].contentEditable = false;
    };
}

function crlDisplayReferencedConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Referenced Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.ReferencedConceptID;
}

function crlDisplayRefinedConcept(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Refined Concept ID";
    typeRow.cells[1].innerHTML = data.NotificationConcept.RefinedConceptID;
}

function crlDisplayType(data, row) {
    var typeRow = crlObtainPropertyRow(row);
    typeRow.cells[0].innerHTML = "Type";
    typeRow.cells[1].innerHTML = data.NotificationConcept.Type;
}

function crlDisplayURI(data, row) {
    var uriRow = crlObtainPropertyRow(row);
    uriRow.cells[0].innerHTML = "URI";
    uriRow.cells[1].innerHTML = data.NotificationConcept.URI;
    uriRow.cells[1].id = "uri";
    if (data.NotificationConcept.IsCore == "false" && data.NotificationConcept.ReadOnly == "false") {
        uriRow.cells[1].contentEditable = true;
        $("#uri").on("keyup", crlSendURIChanged);
    } else {
        uriRow.cells[1].contentEditable = false;
    }
}


function crlDisplayVersion(data, row) {
    var versionRow = crlObtainPropertyRow(row)
    versionRow.cells[0].innerHTML = "Version";
    versionRow.cells[1].innerHTML = data.NotificationConcept.Version;
}

function crlDropdownMenu(dropdownId) {
    document.getElementById(dropdownId).classList.toggle("show");
}

// crlEditorSettingsOK is the callback function for the Editor Settings dialog OK button.
// It updaates the debug settings with the values from the dialog
var crlEditorSettingsOK = function () {
    crlSendEditorSettings();
    crlEditorSettingsDialog.dialog("close");
};

// crlEditorSettings creates and displays the Editor Settings dialog so that the debug settings can be updated from the UI
var crlEditorSettings = function () {
    $("#dropReferenceAsLink").prop("checked", crlDropReferenceAsLink);
    $("#dropRefinementAsLink").prop("checked", crlDropRefinementAsLink);
    crlEditorSettingsDialog.dialog("open");
}


function crlElementSelected(data) {
    if (data.NotificationConceptID != crlSelectedConceptIDGlobal) {
        selectedConceptId = data.NotificationConceptID

        // Update the properties
        crlDisplayType(data, 1);
        crlDisplayID(data, 2);
        crlDisplayVersion(data, 3);
        crlDisplayLabel(data, 4);
        crlDisplayDefinition(data, 5);
        crlDisplayURI(data, 6);
        switch (data.NotificationConcept.Type) {
            case "*core.element":
                crlClearRow(7);
                crlClearRow(8);
                break;
            case "*core.literal":
                crlDisplayLiteralValue(data, 7);
                crlClearRow(8);
                break
            case "*core.reference":
                crlDisplayReferencedConcept(data, 7);
                crlClearRow(8);
                break;
            case "*core.refinement":
                crlDisplayAbstractConcept(data, 7);
                crlDisplayRefinedConcept(data, 8);
                break;
        };
    }
    crlSendNormalResponse()
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
    xhr.send(data);
}

function crlInitializeWebSocket() {
    console.log("Initializing Web Socket")
    // ws initialization
    crlWebsocketGlobal = new WebSocket("ws://localhost:8080/index/ws");
    console.log("Web Socket Initialization complete")
    crlWebsocketGlobal.onmessage = function (e) {
        var data = JSON.parse(e.data)
        console.log("Notification:" + data.Notification)
        switch (data.Notification) {
            case 'AddDiagramLink':
                crlAddDiagramLink(data);
                break;
            case 'AddDiagramNode':
                crlAddDiagramNode(data);
                break;
            case 'AddTreeNode':
                crlAddTreeNode(data);
                break;
            case "ChangeTreeNode":
                crlChangeTreeNode(data);
                break;
            case "DebugSettings":
                crlSaveDebugSettings(data);
                break;
            case "DeleteTreeNode":
                crlDeleteTreeNode(data);
                break;
            case "DisplayDiagram":
                crlDisplayDiagram(data);
                break;
            case "EditorSettings":
                crlSaveEditorSettings(data);
                break;
            case "ElementSelected":
                crlElementSelected(data);
                break;
            case "InitializationComplete":
                crlInitializationCompleteGlobal = true;
                console.log("Initialization Complete")
                crlSendNormalResponse("Processed InitializationComplete")
                break;
            case "UpdateDiagramLink":
                crlUpdateDiagramLink(data);
                break;
            case "UpdateDiagramNode":
                crlUpdateDiagramNode(data);
                break;
            case "WorkspacePath":
                crlUpdateWorkspacePath(data);
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

function crlMakeDiagramVisible(diagramContainerID) {
    var x = document.getElementsByClassName("crlDiagramContainer");
    for (i = 0; i < x.length; i++) {
        if (x.item(i).id == diagramContainerID) {
            x.item(i).style.display = "block";
        } else {
            x.item(i).style.display = "none";
        }
    }
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

var crlOnChangePosition = function (modelElement, position) {
    var jointElementID = modelElement.get("crlJointID");
    var diagramNodeID = crlGetConceptIDFromJointElementID(jointElementID);
    crlMovedNodes[diagramNodeID] = position;
    //    crlSendDiagramNodeNewPosition(diagramNodeID, position)
}

var crlOnDiagramCellPointerDown = function (cellView, event, x, y) {
    var jointElementID = cellView.model.get("crlJointID");
    var diagramNodeID = crlGetConceptIDFromJointElementID(jointElementID);
    if (diagramNodeID == "") {
        console.log("In onDiagramManagerCellPointerDown diagramNodeID is empty")
    }
    crlSendDiagramNodeSelected(diagramNodeID)
}

var crlOnDiagramCellPointerUp = function (cellView, event, x, y) {
    $.each(crlMovedNodes, function (nodeID, position) {
        crlSendDiagramNodeNewPosition(nodeID, position)
    })
}

function crlOnDiagramDrop(event) {
    event.preventDefault();
    console.log("OnDiagramManagerDrop called");
    var conceptID = crlGetConceptIDFromContainerID(event.target.parentElement.parentElement.id);
    var x = event.layerX.toString();
    var y = event.layerY.toString();
    crlSendDiagramDrop(conceptID, x, y);
}

function crlOnDragover(event, data) {
    event.preventDefault()
}

function crlOnEditorDrop(e, data) {
    crlSendSetTreeDragSelection("")
}


function crlOnMakeDiagramVisible(e) {
    var diagramContainerId = e.target.getAttribute("diagramContainerID")
    crlMakeDiagramVisible(diagramContainerId)
}


function crlOnTreeDragStart(e, data) {
    var parentID = e.target.parentElement.id;
    console.log("On Tree Drag Start called, ParentId = " + parentID);
    var selectedElementID = crlGetConceptIDFromTreeNodeID(parentID);
    console.log("selectedElementID = " + selectedElementID)
    crlSendSetTreeDragSelection(selectedElementID);
}


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

var crlOpenWorkspace = function () {
    $("#selectedWorkspaceFolder").val(crlWorkspacePath);
    crlOpenWorkspaceDialog.dialog("open");
}

function crlOpenWorkspaceOK() {
    crlSendOpenWorkspace($("#selectedWorkspaceFolder").val());
    crlOpenWorkspaceDialog.dialog("close");
}

function crlSaveDebugSettings(data) {
    crlEnableTracing = JSON.parse(data.AdditionalParameters["EnableNotificationTracing"]);
    crlSendNormalResponse();
}

function crlSaveEditorSettings(data) {
    crlDropReferenceAsLink = JSON.parse(data.AdditionalParameters["DropReferenceAsLink"]);
    crlDropRefinementAsLink = JSON.parse(data.AdditionalParameters["DropRefinementAsLink"]);
    crlSendNormalResponse();
}

function crlSendDebugSettings() {
    var xhr = crlCreateEmptyRequest()
    var enableNotificationTracing = "false";
    if ($("#enableTracing").prop("checked") == true) {
        enableNotificationTracing = "true"
    };
    var maxTracingDepth = $("#maxTracingDepth").val()
    var data = JSON.stringify({
        "Action": "UpdateDebugSettings",
        "AdditionalParameters": {
            "EnableNotificationTracing": enableNotificationTracing,
            "MaxTracingDepth": maxTracingDepth
        }
    });
    xhr.send(data);
}

function crlSendDefinitionChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "DefinitionChanged",
        "RequestConceptID": selectedConceptId,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    xhr.send(data);
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
    xhr.send(data);
}

function crlSendLabelChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "LabelChanged",
        "RequestConceptID": selectedConceptId,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    xhr.send(data)
}

function crlSendLiteralValueChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "LiteralValueChanged",
        "RequestConceptID": selectedConceptId,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    xhr.send(data)
}

function crlSendOpenWorkspace(workspacePath) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "OpenWorkspace",
        "AdditionalParameters": {
            "WorkspacePath": workspacePath
        }
    });
    xhr.send(data);
}

function crlSendSaveWorkspace() {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "SaveWorkspace"
    });
    xhr.send(data);
}

function crlSendURIChanged(evt, obj) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({
        "Action": "URIChanged",
        "RequestConceptID": selectedConceptId,
        "AdditionalParameters":
            { "NewValue": evt.currentTarget.textContent }
    });
    xhr.send(data)
}

function crlSendNewConceptSpaceRequest(evt) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "NewConceptSpaceRequest" });
    xhr.send(data)
}

function crlSendNewDiagramRequest(evt) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "NewDiagramRequest" });
    xhr.send(data)
}

function crlSendDiagramDrop(diagramID, x, y) {
    console.log("Diagram Drop ID, x, y: " + diagramID + "  " + x.toString() + "  " + y.toString());
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
    console.log(data);
    xhr.send(data);
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
    xhr.send(data);
}

function crlSendDiagramNodeSelected(nodeID) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "DiagramNodeSelected", "RequestConceptID": nodeID });
    xhr.send(data);
}

function crlSendNormalResponse() {
    var data = {};
    data["Result"] = 0;
    data["ErrorMessage"] = "none"
    crlWebsocketGlobal.send(JSON.stringify(data))
    console.log('Sent normal response');
}

function crlSendSetTreeDragSelection(id) {
    var xhr = crlCreateEmptyRequest();
    var data = JSON.stringify({ "Action": "SetTreeDragSelection", "RequestConceptID": id });
    xhr.send(data);
}

function crlSendTreeNodeSelected(evt, obj) {
    if (obj != undefined) {
        var xhr = crlCreateEmptyRequest();
        var conceptID = crlGetConceptIDFromTreeNodeID(obj.node.id)
        var data = JSON.stringify({ "Action": "TreeNodeSelected", "RequestConceptID": conceptID });
        xhr.send(data);
    };
}

// <!-- Define sizeAll() to manage sizes of display components -->
function crlSizeAll() {
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

var crlUpdateDiagramLink = function (data) {
    var concept = data.NotificationConcept;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var linkID = crlGetJointCellIDFromConceptID(concept.ConceptID);
    var link = crlFindLinkInGraph(graphID, linkID)
    if (link == undefined) {
        link = crlConstructDiagramLink(data, graphID, linkID);
    }
    link.label(0, {
        attrs: {
            text: {
                text: data.AdditionalParameters["DisplayLabel"]
            }
        }
    });
    crlSendNormalResponse()
}

var crlUpdateDiagramNode = function (data) {
    var concept = data.NotificationConcept;
    var params = data.AdditionalParameters;
    var owningConceptID = concept.OwningConceptID;
    var graphID = crlGetJointGraphIDFromDiagramID(owningConceptID);
    var nodeID = crlGetJointCellIDFromConceptID(concept.ConceptID);
    var node = crlFindElementInGraph(graphID, nodeID);
    if (node == undefined) {
        node = crlConstructDiagramNode(data, graphID, nodeID);
    };
    node.set("displayLabelYOffset", Number(params["DisplayLabelYOffset"]));
    node.set('position', { "x": Number(params["NodeX"]), "y": Number(params["NodeY"]) });
    node.set('size', { "width": Number(params["NodeWidth"]), "height": Number(params["NodeHeight"]) });
    node.set('icon', params["Icon"]);
    node.set('name', params["DisplayLabel"]);
    node.set("abstractions", params["Abstractions"]);
    crlSendNormalResponse();
}

var crlUpdateWorkspacePath = function (data) {
    crlWorkspacePath = data.AdditionalParameters["WorkspacePath"];
    crlSendNormalResponse();
}