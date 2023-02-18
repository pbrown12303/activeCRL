// Copyright 2017, 2018 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"log"
	"sync"

	"github.com/pkg/errors"
)

// Transaction maintains a record of which elements are currently read and write locked and provides facilities
// for locking additional elements. It also manages function calls. Part of the management is the
// suppression of circular function calls.
type Transaction struct {
	sync.Mutex
	readLocks  map[string]Element
	uOfD       *UniverseOfDiscourse
	writeLocks map[string]Element
	// The key to inProgressCalls is the catenation of the functionID and the target element ID
	inProgressCalls map[string]bool
}

// callFunction calls the referenced function on the target element
func (transPtr *Transaction) callFunctions(functionID string, targetElement Element, notification *ChangeNotification) error {
	// First, check to see whether the targetElement is in the process of being deleted. If it is, simply return: we don't
	// execute functions on deleted elements
	targetID := targetElement.GetConceptID(transPtr)
	if notification.uOfD.inProgressDeletions.Contains(targetID) {
		return nil
	}
	for _, function := range transPtr.uOfD.getFunctions(functionID) {
		if transPtr.uOfD.getExecutedCalls() != nil {
			functionCallRecord, err := newFunctionCallRecord(functionID, function, targetElement, notification)
			transPtr.uOfD.getExecutedCalls() <- functionCallRecord
			if err != nil {
				return errors.Wrap(err, "Transaction.callFunctions failed to queue functionCallRecord")
			}
		}
		if TraceLocks || TraceChange {
			omitCall := (OmitManageTreeNodesCalls && functionID == "http://activeCrl.com/crlEditor/EditorDomain/TreeViews/TreeNodeManager") ||
				(OmitDiagramRelatedCalls && isDiagramRelatedFunction(functionID))
			if !omitCall {
				functionCallGraphs = append(functionCallGraphs, NewFunctionCallGraph(functionID, targetElement, notification, transPtr))
			}
		}
		inProgressKey := functionID + targetID
		if !transPtr.inProgressCalls[inProgressKey] {
			transPtr.inProgressCalls[inProgressKey] = true
			err := function(targetElement, notification, transPtr)
			if err != nil {
				delete(transPtr.inProgressCalls, inProgressKey)
				return errors.Wrap(err, "Transaction.callFunctions failed")
			}
			delete(transPtr.inProgressCalls, inProgressKey)
		}
	}
	return nil
}

// GetUniverseOfDiscourse returns the UniverseOfDiscourse to which this HeldLocks belongs
func (transPtr *Transaction) GetUniverseOfDiscourse() *UniverseOfDiscourse {
	return transPtr.uOfD
}

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
