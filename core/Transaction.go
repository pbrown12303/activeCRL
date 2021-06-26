// Copyright 2017, 2018 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"errors"
	"log"
	"sync"
)

// Transaction maintains a record of which elements are currently read and write locked and provides facilities
// for locking additional elements.
type Transaction struct {
	sync.Mutex
	functionCallManager *functionCallManager
	readLocks           map[string]Element
	uOfD                *UniverseOfDiscourse
	writeLocks          map[string]Element
}

// GetUniverseOfDiscourse returns the UniverseOfDiscourse to which this HeldLocks belongs
func (hlPtr *Transaction) GetUniverseOfDiscourse() *UniverseOfDiscourse {
	return hlPtr.uOfD
}

// // IsLocked checks to see whether this HeldLocks structure already has a record of the Element being locked
// // and returns the result.
// func (hlPtr *HeldLocks) IsLocked(el Element) bool {
// 	hlPtr.Lock()
// 	defer hlPtr.Unlock()
// 	id := el.getConceptIDNoLock()
// 	return hlPtr.writeLocks[id] != nil
// }

// ReadLockElement checks to see whether this HeldLocks structure already has a record of the Element being
// locked, either read or write. If it does, it simply returns. If not, it attempts to acquire the read on the Element and makes
// a record of the fact that the read lock has been obtained.
func (hlPtr *Transaction) ReadLockElement(el Element) {
	hlPtr.Lock()
	defer hlPtr.Unlock()
	id := el.getConceptIDNoLock()
	_, writeLocked := hlPtr.writeLocks[id]
	if writeLocked {
		return
	}
	_, readLocked := hlPtr.readLocks[id]
	if !readLocked {
		el.TraceableReadLock(hlPtr)
		hlPtr.readLocks[id] = el
	}
}

// WriteLockElement checks to see whether this HeldLocks structure already has a record of the Element being
// write locked. If it does, it simply returns. If not, it attempts to acquire the write lock on the Element and makes
// a record of the fact that the lock has been obtained.
func (hlPtr *Transaction) WriteLockElement(el Element) error {
	hlPtr.Lock()
	defer hlPtr.Unlock()
	id := el.getConceptIDNoLock()
	_, readLocked := hlPtr.readLocks[id]
	if readLocked {
		return errors.New("Write lock attempted on Element with read lock: " + id)
	}
	_, writeLocked := hlPtr.writeLocks[id]
	if !writeLocked {
		el.TraceableWriteLock(hlPtr)
		hlPtr.writeLocks[id] = el
	}
	return nil
}

// ReleaseLocks releases all pending functions for execution (asynchronously) and releases all currently held locks
func (hlPtr *Transaction) ReleaseLocks() {
	hlPtr.Lock()
	if TraceLocks {
		log.Printf("HL %p about to ReleaseLocks", hlPtr)
	}
	defer hlPtr.Unlock()
	for _, el := range hlPtr.readLocks {
		el.TraceableReadUnlock(hlPtr)
		delete(hlPtr.readLocks, el.getConceptIDNoLock())
	}
	for _, el := range hlPtr.writeLocks {
		el.TraceableWriteUnlock(hlPtr)
		delete(hlPtr.writeLocks, el.getConceptIDNoLock())
	}
	err := hlPtr.functionCallManager.callQueuedFunctions(hlPtr)
	if err != nil {
		log.Print(err)
	}
}

// ReleaseLocksAndWait releases all pending functions for execution (asynchronously) and releases all currently held locks.
// It returns when all of the asynchronous function executions have completed. This is particularly useful to when you want
// to wait for all of the side effects for some change to complete before proceeding with additonal changes.
func (hlPtr *Transaction) ReleaseLocksAndWait() {
	if TraceLocks {
		log.Printf("HL %p about to ReleaseLocksAndWait 11111111111111111111111111111111111111111111", hlPtr)
	}
	hlPtr.ReleaseLocks()
	if TraceLocks {
		log.Printf("HL %p locks released, about to Wait 2222222222222222222222222222222222222222222", hlPtr)
	}
	if TraceLocks {
		log.Printf("HL %p finished waiting 33333333333333333333333333333333333333333333333333333333", hlPtr)
	}
}
