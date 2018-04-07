// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	//	"runtime/debug"
	"sync"
)

// Due to a compiler issue leading to a "interface Element contains embedded non-interface BaseElement" error
// the key in this map is made a BaseElement rather than an Element.
type elementNotificationsMap map[BaseElement][]*ChangeNotification

type FunctionCallManager struct {
	functionTargetMap map[crlExecutionFunctionArrayIdentifier]elementNotificationsMap
}

func NewFunctionCallManager() *FunctionCallManager {
	var fcm FunctionCallManager
	fcm.functionTargetMap = make(map[crlExecutionFunctionArrayIdentifier]elementNotificationsMap)
	return &fcm
}

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

func (fcm *FunctionCallManager) ExecuteFunctions(wg *sync.WaitGroup) {
	for labeledFunction, enm := range fcm.functionTargetMap {
		for el, notifications := range enm {
			callLabeledFunction(labeledFunction, el.(Element), notifications, wg)
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

func makeGoCall(functionId crlExecutionFunctionArrayIdentifier, el Element, notifications []*ChangeNotification, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	functions := GetCore().computeFunctions[functionId]
	if functions != nil {
		for _, function := range functions {
			function(el, notifications, wg)
		}
	} else {
		log.Printf("In makeGoCall, function not found for identifier: %s ", functionId)
	}
}

func (fcm *FunctionCallManager) Print(prefix string) {
	log.Printf(prefix + "Pending Function Calls")
	for pendingFunctionIdentifier, enm := range fcm.functionTargetMap {
		log.Printf(prefix+"   Pending function: %s\n", pendingFunctionIdentifier)
		for el, notifications := range enm {
			log.Printf(prefix+"      Element Id: %s Notifications length: %d\n", el.getIdNoLock(), len(notifications))
		}
	}
}
