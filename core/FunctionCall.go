// Copyright 2017, 2018 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/pkg/errors"
)

// The crlExecutionFunction is the standard signature of a function that gets called when an element (including
// its children) experience a change. Its arguments are the element that changed, the array of ChangeNotifications, and
// a pointer to a WaitGroup that is used to determine (on a larger scale) when the execution of the triggered functions
// has completed.
type crlExecutionFunction func(Element, *ChangeNotification, *Transaction) error

type functionCallRecord struct {
	function     crlExecutionFunction
	functionID   string
	target       Element
	notification *ChangeNotification
}

func newFunctionCallRecord(functionID string, function crlExecutionFunction, target Element, notification *ChangeNotification) (*functionCallRecord, error) {
	if target == nil {
		return nil, errors.New("FunctionCallManager.go newPendingFunctionCall invoked with nil target")
	}
	var functionCall functionCallRecord
	functionCall.function = function
	functionCall.functionID = functionID
	functionCall.target = target
	functionCall.notification = notification
	return &functionCall, nil
}

// The functions type maps core Element identifiers to the array of crlExecutionFunctions associated with the identfier.
type functions map[string][]crlExecutionFunction

// isDiagramRelatedFunction returns true if the functionID matches one of the diagram related functions
func isDiagramRelatedFunction(functionID string) bool {
	if functionID == "http://activeCrl.com/corediagram/CoreDiagram/CrlDiagram" ||
		functionID == "http://activeCrl.com/corediagram/CoreDiagram/CrlDiagramElement" ||
		functionID == "http://activeCrl.com/corediagram/CoreDiagram/OwnerPointer" ||
		functionID == "http://activeCrl.com/crlEditor/EditorDomain/DiagramViewMonitor" ||
		functionID == "http://activeCrl.com/crlEditor/EditorDomain/TreeViews/TreeNodeManager" {
		return true
	}
	return false
}
