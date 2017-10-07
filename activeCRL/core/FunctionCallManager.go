// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
)

// Due to a compiler issue leading to a "interface Element contains embedded non-interface BaseElement" error
// the key in this map is made a BaseElement rather than an Element.
type elementNotificationsMap map[BaseElement][]*ChangeNotification

type FunctionCallManager struct {
	functionTargetMap map[*labeledFunction]elementNotificationsMap
}

func NewFunctionCallManager() *FunctionCallManager {
	var fcm FunctionCallManager
	fcm.functionTargetMap = make(map[*labeledFunction]elementNotificationsMap)
	return &fcm
}

func (fcm *FunctionCallManager) AddFunctionCall(lf *labeledFunction, el Element, notification *ChangeNotification) {
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
}

func callLabeledFunction(lf *labeledFunction, el Element, notifications []*ChangeNotification, wg *sync.WaitGroup) {
	// We have to call wg.Add() before the go call because there may be a delay between when the gorouting is invoked and
	// when it gets around to calling wg.Add()
	if wg != nil {
		wg.Add(1)
	}
	go makeGoCall(lf, el, notifications, wg)
}

func makeGoCall(lf *labeledFunction, el Element, notifications []*ChangeNotification, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	lf.function(el, notifications, wg)
}

func (fcm *FunctionCallManager) Print(prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	log.Printf(prefix + "Pending Function Calls")
	for pendingFunction, enm := range fcm.functionTargetMap {
		log.Printf(prefix+"   Pending function: %s\n", pendingFunction.label)
		for el, notifications := range enm {
			log.Printf(prefix+"      Element Id: %s\n", el.GetId(hl))
			for _, notification := range notifications {
				notification.Print(prefix+"         ", hl)
			}
		}
	}
}
