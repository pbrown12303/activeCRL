// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/satori/go.uuid"
	"log"
	//	"runtime/debug"
	"sync"
)

var CountHeldLocksInvocations bool = false
var HeldLocksInvocations int = 0
var NewLockInvocations int = 0

type HeldLocks struct {
	sync.Mutex
	beLocks             map[uuid.UUID]BaseElement
	waitGroup           *sync.WaitGroup
	functionCallManager *FunctionCallManager
}

func NewHeldLocks(wg *sync.WaitGroup) *HeldLocks {
	var hl HeldLocks
	hl.beLocks = make(map[uuid.UUID]BaseElement)
	hl.waitGroup = wg
	hl.functionCallManager = NewFunctionCallManager()
	return &hl
}

func (hlPtr *HeldLocks) GetWaitGroup() *sync.WaitGroup {
	return hlPtr.waitGroup
}

func (hlPtr *HeldLocks) LockBaseElement(be BaseElement) {
	if CountHeldLocksInvocations == true {
		HeldLocksInvocations++
	}
	hlPtr.Lock()
	defer hlPtr.Unlock()
	id := be.getIdNoLock()
	if AdHocTrace == true {
		log.Printf("In LockBaseElement for be %s \n", id)
	}
	if hlPtr.beLocks[id] == nil {
		if AdHocTrace == true {
			log.Printf("In LockBaseElement, lock not found for be %s of type %T \n", id, be)
			//			debug.PrintStack()
		}
		//		log.Printf("Locking %s", id)
		if CountHeldLocksInvocations == true {
			NewLockInvocations++
		}
		hlPtr.beLocks[id] = be
		be.TraceableLock()
	}
}

func (hlPtr *HeldLocks) ReleaseLocks() {
	hlPtr.Lock()
	defer hlPtr.Unlock()
	for _, be := range hlPtr.beLocks {
		be.TraceableUnlock()
	}
	hlPtr.beLocks = make(map[uuid.UUID]BaseElement)
	hlPtr.functionCallManager.ExecuteFunctions(hlPtr.waitGroup)
}

func (hlPtr *HeldLocks) ReleaseLocksAndWait() {
	hlPtr.ReleaseLocks()
	hlPtr.waitGroup.Wait()
}
