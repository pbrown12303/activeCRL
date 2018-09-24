// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	//	"runtime/debug"
	"sync"
)

// The crlExecutionFunction is the standard signature of a function that gets called when an element (including
// its children) experience a change. Its arguments are the element that changed, the array of ChangeNotifications, and
// a pointer to a WaitGroup that is used to determine (on a larger scale) when the execution of the triggered functions
// has completed.
type crlExecutionFunction func(Element, []*ChangeNotification, *sync.WaitGroup)

// The crlExecutionFunctionArrayIdentifier is expected to be the string version of the URI associated with the element.
// It serves as the index into the array of functions associated with the core Element.
type crlExecutionFunctionArrayIdentifier string

// The functions type maps core Element identifiers to the array of crlExecutionFunctions associated with the identfier.
type functions map[crlExecutionFunctionArrayIdentifier][]crlExecutionFunction

// Due to a compiler issue leading to a "interface Element contains embedded non-interface BaseElement" error
// the key in this map is made a BaseElement rather than an Element.
type elementNotificationsMap map[BaseElement][]*ChangeNotification

// FunctionCallManager manages the set of pending function calls
type FunctionCallManager struct {
	functionTargetMap map[crlExecutionFunctionArrayIdentifier]elementNotificationsMap
}

// NewFunctionCallManager creates and initializes a FunctionCallManager
func NewFunctionCallManager() *FunctionCallManager {
	var fcm FunctionCallManager
	fcm.functionTargetMap = make(map[crlExecutionFunctionArrayIdentifier]elementNotificationsMap)
	return &fcm
}

// AddFunctionCall adds a pending function call to the manager. The function is identified by a crlExecutionFunctionArrayIdentifier, which is
// a string. The Element is the element that will eventually "execute" the function, and the ChangeNotification is the trigger
// that caused the function to be queued for execution.
func (fcm *FunctionCallManager) AddFunctionCall(lf crlExecutionFunctionArrayIdentifier, el Element, notification *ChangeNotification) {
	enm := fcm.functionTargetMap[lf]
	if enm == nil {
		enm = make(map[BaseElement][]*ChangeNotification)
	}
	changeNotificationArray := enm[el]
	modifiedChangeNotificationArray := append(changeNotificationArray, notification)
	enm[el] = modifiedChangeNotificationArray
	fcm.functionTargetMap[lf] = enm
}

// ExecuteFunctions does not actually execute the functions: instead, for each function, it creates a goroutine
// that executes the function asynchronously. It uses the WaitGroup from the HeldLocks to keep track of the
// function execution completions.
func (fcm *FunctionCallManager) ExecuteFunctions(hl *HeldLocks) {
	for functionID, enm := range fcm.functionTargetMap {
		for el, notifications := range enm {
			if TraceChange == true {
				log.Printf("Adding a function call graph \n")
				functionCallGraphs = append(functionCallGraphs, NewFunctionCallGraph(functionID, el.(Element), notifications))
			}
			callLabeledFunction(functionID, el.(Element), notifications, hl.waitGroup)
		}
	}
	fcm.clearFunctionCalls()
}

func callLabeledFunction(lf crlExecutionFunctionArrayIdentifier, el Element, notifications []*ChangeNotification, wg *sync.WaitGroup) {
	// We have to call wg.Add() before the go call because there may be a delay between when the gorouting is invoked and
	// when it gets around to calling wg.Add()
	if wg != nil {
		wg.Add(1)
	}
	go makeGoCall(lf, el, notifications, wg)
}

func (fcm *FunctionCallManager) clearFunctionCalls() {
	fcm.functionTargetMap = make(map[crlExecutionFunctionArrayIdentifier]elementNotificationsMap)
}

func makeGoCall(functionID crlExecutionFunctionArrayIdentifier, el Element, notifications []*ChangeNotification, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	functions := GetCore().computeFunctions[functionID]
	if functions != nil {
		for _, function := range functions {
			function(el, notifications, wg)
		}
	} else {
		log.Printf("In makeGoCall, function not found for identifier: %s ", functionID)
	}
}

// Print prints out the pending function calls
func (fcm *FunctionCallManager) Print(prefix string) {
	log.Printf(prefix + "Pending Function Calls")
	for pendingFunctionIdentifier, enm := range fcm.functionTargetMap {
		log.Printf(prefix+"   Pending function: %s\n", pendingFunctionIdentifier)
		for el, notifications := range enm {
			log.Printf(prefix+"      Element Id: %s Notifications length: %d\n", el.getIdNoLock(), len(notifications))
		}
	}
}
