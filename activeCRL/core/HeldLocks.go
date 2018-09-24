// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
)

// HeldLocks maintains a record of which elements are currently locked and provides facilities
// for locking additional elements. It also contains a waitGroup that is used to verify the
// completion of asynchronous function executions.
type HeldLocks struct {
	sync.Mutex
	beLocks             map[string]BaseElement
	waitGroup           *sync.WaitGroup
	functionCallManager *FunctionCallManager
}

// NewHeldLocks creates and initializes a HeldLocks structure utilizing the supplied WaitGroup
func NewHeldLocks(wg *sync.WaitGroup) *HeldLocks {
	var hl HeldLocks
	hl.beLocks = make(map[string]BaseElement)
	hl.waitGroup = wg
	hl.functionCallManager = NewFunctionCallManager()
	return &hl
}

// GetWaitGroup returns the WaitGroup being used by HeldLocks
func (hlPtr *HeldLocks) GetWaitGroup() *sync.WaitGroup {
	return hlPtr.waitGroup
}

// IsLocked checks to see whether this HeldLocks structure already has a record of the BaseElement being locked
// and returns the result.
func (hlPtr *HeldLocks) IsLocked(be BaseElement) bool {
	hlPtr.Lock()
	defer hlPtr.Unlock()
	id := be.getIdNoLock()
	return hlPtr.beLocks[id] != nil
}

// LockBaseElement checks to see whether this HeldLocks structure already has a record of the BaseElement being
// locked. If it does, it simply returns. If not, it attempts to acquire the lock on the BaseElement and makes
// a record of the fact that the lock has been obtained.
func (hlPtr *HeldLocks) LockBaseElement(be BaseElement) {
	hlPtr.Lock()
	defer hlPtr.Unlock()
	id := be.getIdNoLock()
	if hlPtr.beLocks[id] == nil {
		hlPtr.beLocks[id] = be
		be.TraceableLock()
	}
}

// ReleaseLocks releases all pending functions for execution (asynchronously) and releases all currently held locks
func (hlPtr *HeldLocks) ReleaseLocks() {
	hlPtr.Lock()
	defer hlPtr.Unlock()
	hlPtr.functionCallManager.ExecuteFunctions(hlPtr)
	for _, be := range hlPtr.beLocks {
		be.TraceableUnlock()
	}
	hlPtr.beLocks = make(map[string]BaseElement)
}

// ReleaseLocksAndWait releases all pending functions for execution (asynchronously) and releases all currently held locks.
// It returns when all of the asynchronous function executions have completed. This is particularly useful to when you want
// to wait for all of the side effects for some change to complete before proceeding with additonal changes.
func (hlPtr *HeldLocks) ReleaseLocksAndWait() {
	hlPtr.ReleaseLocks()
	hlPtr.waitGroup.Wait()
}
