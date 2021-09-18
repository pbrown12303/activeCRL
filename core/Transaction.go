// Copyright 2017, 2018 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

// Transaction maintains a record of which elements are currently read and write locked and provides facilities
// for locking additional elements.
type Transaction struct {
	sync.Mutex
	// functionCallManager *functionCallManager
	readLocks         map[string]Element
	uOfD              *UniverseOfDiscourse
	writeLocks        map[string]Element
	functionCallQueue *pendingFunctionCallQueue
}

// addFunctionCall adds a pending function call to the Transaction for each function associated with the functionID.
// The Element is the element that will eventually "execute" the function, and the ChangeNotification is the trigger
// that caused the function to be queued for execution.
func (transPtr *Transaction) addFunctionCall(functionID string, targetElement Element, notification *ChangeNotification) error {
	for _, function := range transPtr.uOfD.getFunctions(functionID) {
		pendingCall, err := newPendingFunctionCall(functionID, function, targetElement, notification)
		if err != nil {
			return errors.Wrap(err, "functionCallManager.addFunctionCall failed")
		}
		newCount := atomic.AddInt32(&pendingFunctionCount, 1)
		if CrlLogPendingFunctionCount {
			log.Printf("Pending function count: %d", newCount)
		}
		transPtr.functionCallQueue.enqueue(pendingCall)
	}
	return nil
}

// callQueuedFunctions calls each function on the pending function queue
func (transPtr *Transaction) callQueuedFunctions(hl *Transaction) error {
	for transPtr.functionCallQueue.queueHead != nil {
		pendingCall := transPtr.functionCallQueue.dequeue()
		if transPtr.uOfD.getExecutedCalls() != nil {
			transPtr.uOfD.getExecutedCalls() <- pendingCall
		}
		if TraceLocks || TraceChange {
			omitCall := (OmitHousekeepingCalls && pendingCall.functionID == "http://activeCrl.com/core/coreHousekeeping") ||
				(OmitManageTreeNodesCalls && pendingCall.functionID == "http://activeCrl.com/crlEditor/EditorDomain/TreeViews/TreeNodeManager") ||
				(OmitDiagramRelatedCalls && isDiagramRelatedFunction(pendingCall.functionID))
			if !omitCall {
				log.Printf("About to execute %s with notification %s target %p", pendingCall.functionID, pendingCall.notification.GetNatureOfChange().String(), pendingCall.target)
				log.Printf("   Function target: %T %s %s %p", pendingCall.target, pendingCall.target.getConceptIDNoLock(), pendingCall.target.getLabelNoLock(), pendingCall.target)
				functionCallGraphs = append(functionCallGraphs, NewFunctionCallGraph(pendingCall.functionID, pendingCall.target, pendingCall.notification, hl))
			}
		}
		err := pendingCall.function(pendingCall.target, pendingCall.notification, transPtr)
		if err != nil {
			return errors.Wrap(err, "functionCallManager.callQueuedFunctions failed")
		}
		newCount := atomic.AddInt32(&pendingFunctionCount, -1)
		if CrlLogPendingFunctionCount {
			log.Printf("Pending function count: %d", newCount)
			log.Printf("Dequeued call: %+v", pendingCall)
		}
	}
	return nil
}

// GetUniverseOfDiscourse returns the UniverseOfDiscourse to which this HeldLocks belongs
func (transPtr *Transaction) GetUniverseOfDiscourse() *UniverseOfDiscourse {
	return transPtr.uOfD
}

// // IsLocked checks to see whether this HeldLocks structure already has a record of the Element being locked
// // and returns the result.
// func (transPtr *HeldLocks) IsLocked(el Element) bool {
// 	transPtr.Lock()
// 	defer transPtr.Unlock()
// 	id := el.getConceptIDNoLock()
// 	return transPtr.writeLocks[id] != nil
// }

// ReadLockElement checks to see whether this HeldLocks structure already has a record of the Element being
// locked, either read or write. If it does, it simply returns. If not, it attempts to acquire the read on the Element and makes
// a record of the fact that the read lock has been obtained.
func (transPtr *Transaction) ReadLockElement(el Element) {
	transPtr.Lock()
	defer transPtr.Unlock()
	id := el.getConceptIDNoLock()
	_, writeLocked := transPtr.writeLocks[id]
	if writeLocked {
		return
	}
	_, readLocked := transPtr.readLocks[id]
	if !readLocked {
		el.TraceableReadLock(transPtr)
		transPtr.readLocks[id] = el
	}
}

// WriteLockElement checks to see whether this HeldLocks structure already has a record of the Element being
// write locked. If it does, it simply returns. If not, it attempts to acquire the write lock on the Element and makes
// a record of the fact that the lock has been obtained.
func (transPtr *Transaction) WriteLockElement(el Element) error {
	transPtr.Lock()
	defer transPtr.Unlock()
	id := el.getConceptIDNoLock()
	_, readLocked := transPtr.readLocks[id]
	if readLocked {
		return errors.New("Write lock attempted on Element with read lock: " + id)
	}
	_, writeLocked := transPtr.writeLocks[id]
	if !writeLocked {
		el.TraceableWriteLock(transPtr)
		transPtr.writeLocks[id] = el
	}
	return nil
}

// ReleaseLocks releases all pending functions for execution (asynchronously) and releases all currently held locks
func (transPtr *Transaction) ReleaseLocks() {
	// Execute all the queued functions before releasing locks
	err := transPtr.callQueuedFunctions(transPtr)
	if err != nil {
		log.Print(err)
	}
	transPtr.Lock()
	defer transPtr.Unlock()
	if TraceLocks {
		log.Printf("HL %p about to ReleaseLocks", transPtr)
	}
	for _, el := range transPtr.readLocks {
		el.TraceableReadUnlock(transPtr)
		delete(transPtr.readLocks, el.getConceptIDNoLock())
	}
	for _, el := range transPtr.writeLocks {
		el.TraceableWriteUnlock(transPtr)
		delete(transPtr.writeLocks, el.getConceptIDNoLock())
	}
}

// ReleaseLocksAndWait releases all pending functions for execution (asynchronously) and releases all currently held locks.
// It returns when all of the asynchronous function executions have completed. This is particularly useful to when you want
// to wait for all of the side effects for some change to complete before proceeding with additonal changes.
func (transPtr *Transaction) ReleaseLocksAndWait() {
	if TraceLocks {
		log.Printf("HL %p about to ReleaseLocksAndWait 11111111111111111111111111111111111111111111", transPtr)
	}
	transPtr.ReleaseLocks()
	if TraceLocks {
		log.Printf("HL %p locks released, about to Wait 2222222222222222222222222222222222222222222", transPtr)
	}
	if TraceLocks {
		log.Printf("HL %p finished waiting 33333333333333333333333333333333333333333333333333333333", transPtr)
	}
}
