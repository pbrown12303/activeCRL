// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	//	"log"
	"sync"
)

type HeldLocks struct {
	sync.Mutex
	beLocks             map[string]BaseElement
	waitGroup           *sync.WaitGroup
	functionCallManager *FunctionCallManager
}

func NewHeldLocks(wg *sync.WaitGroup) *HeldLocks {
	var hl HeldLocks
	hl.beLocks = make(map[string]BaseElement)
	hl.waitGroup = wg
	hl.functionCallManager = NewFunctionCallManager()
	return &hl
}

func (hlPtr *HeldLocks) GetWaitGroup() *sync.WaitGroup {
	return hlPtr.waitGroup
}

func (hlPtr *HeldLocks) LockBaseElement(be BaseElement) {
	hlPtr.Lock()
	defer hlPtr.Unlock()
	id := be.getIdNoLock().String()
	if hlPtr.beLocks[id] == nil {
		//		log.Printf("Locking %s", id)
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
	hlPtr.beLocks = make(map[string]BaseElement)
	hlPtr.functionCallManager.ExecuteFunctions(hlPtr.waitGroup)
}
