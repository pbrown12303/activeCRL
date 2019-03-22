// Copyright 2017, 2018 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"log"
	"sync"
)

// The crlExecutionFunction is the standard signature of a function that gets called when an element (including
// its children) experience a change. Its arguments are the element that changed, the array of ChangeNotifications, and
// a pointer to a WaitGroup that is used to determine (on a larger scale) when the execution of the triggered functions
// has completed.
type crlExecutionFunction func(Element, *ChangeNotification, UniverseOfDiscourse)

type pendingFunctionCall struct {
	function     crlExecutionFunction
	functionID   string
	target       Element
	notification *ChangeNotification
}

func newPendingFunctionCall(functionID string, function crlExecutionFunction, target Element, notification *ChangeNotification) *pendingFunctionCall {
	var pendingCall pendingFunctionCall
	pendingCall.function = function
	pendingCall.functionID = functionID
	pendingCall.target = target
	pendingCall.notification = notification
	return &pendingCall
}

type pendingFunctionCallEntry struct {
	pendingCall *pendingFunctionCall
	next        *pendingFunctionCallEntry
}

func newPendingFunctionCallEntry(pendingCall *pendingFunctionCall) *pendingFunctionCallEntry {
	var entry pendingFunctionCallEntry
	entry.pendingCall = pendingCall
	return &entry
}

// pendingFunctionCallQueue maintains a linked list of pending function calls
type pendingFunctionCallQueue struct {
	sync.Mutex
	queueHead *pendingFunctionCallEntry
	queueTail *pendingFunctionCallEntry
}

func newPendingFunctionCallQueue() *pendingFunctionCallQueue {
	var queue pendingFunctionCallQueue
	return &queue
}

func (queue *pendingFunctionCallQueue) enqueue(pendingCall *pendingFunctionCall) error {
	queue.Mutex.Lock()
	defer queue.Mutex.Unlock()
	if pendingCall == nil {
		return errors.New("pendingFunctionCallQueue.enqueue called with nil pendingCall")
	}
	currentTail := queue.queueTail
	newTail := newPendingFunctionCallEntry(pendingCall)
	if currentTail == nil {
		queue.queueHead = newTail
	} else {
		currentTail.next = newTail
	}
	queue.queueTail = newTail
	return nil
}

func (queue *pendingFunctionCallQueue) dequeue() *pendingFunctionCall {
	queue.Mutex.Lock()
	defer queue.Mutex.Unlock()
	currentHead := queue.queueHead
	if currentHead != nil {
		queue.queueHead = currentHead.next
		if currentHead.next == nil {
			queue.queueTail = nil
		} else {
			currentHead.next = nil
		}
		return currentHead.pendingCall
	}
	return nil
}

func (queue *pendingFunctionCallQueue) findFirstPendingCall(functionID string, target Element) *pendingFunctionCall {
	currentCandidate := queue.queueHead
	for currentCandidate != nil {
		currentCall := currentCandidate.pendingCall
		if currentCall.functionID == functionID && currentCall.target == target {
			return currentCall
		}
		currentCandidate = currentCandidate.next
	}
	return nil
}

func (queue *pendingFunctionCallQueue) isEmpty() bool {
	return queue.queueHead == nil
}

// The functions type maps core Element identifiers to the array of crlExecutionFunctions associated with the identfier.
type functions map[string][]crlExecutionFunction

// functionCallManager manages the set of pending function calls
type functionCallManager struct {
	functionCallQueue *pendingFunctionCallQueue
	uOfD              UniverseOfDiscourse
}

// newFunctionCallManager creates and initializes a FunctionCallManager
func newFunctionCallManager(uOfD UniverseOfDiscourse) *functionCallManager {
	var fcm functionCallManager
	fcm.uOfD = uOfD
	fcm.functionCallQueue = newPendingFunctionCallQueue()
	return &fcm
}

// addFunctionCall adds a pending function call to the manager for each function associated with the functionIK.
// The Element is the element that will eventually "execute" the function, and the ChangeNotification is the trigger
// that caused the function to be queued for execution.
func (fcm *functionCallManager) addFunctionCall(functionID string, targetElement Element, notification *ChangeNotification) {
	for _, function := range fcm.uOfD.GetFunctions(functionID) {
		pendingCall := newPendingFunctionCall(functionID, function, targetElement, notification)
		fcm.functionCallQueue.enqueue(pendingCall)
	}
}

// callQueuedFunctions calls each function on the pending function queue
func (fcm *functionCallManager) callQueuedFunctions(hl *HeldLocks) {
	for fcm.functionCallQueue.queueHead != nil {
		pendingCall := fcm.functionCallQueue.dequeue()
		if fcm.uOfD.getExecutedCalls() != nil {
			fcm.uOfD.getExecutedCalls() <- pendingCall
		}
		if TraceLocks == true || TraceChange == true {
			log.Printf("About to execute %s with notification %s target %p", pendingCall.functionID, pendingCall.notification.GetNatureOfChange().String(), pendingCall.target)
		}
		pendingCall.function(pendingCall.target, pendingCall.notification, fcm.uOfD)
	}
}
