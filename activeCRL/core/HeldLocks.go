package core

import (
	"sync"
)

type HeldLocks struct {
	sync.Mutex
	beLocks             []BaseElement
	waitGroup           *sync.WaitGroup
	functionCallManager *FunctionCallManager
}

func NewHeldLocks(wg *sync.WaitGroup) *HeldLocks {
	var hl HeldLocks
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
	found := false
	for _, lbe := range hlPtr.beLocks {
		if lbe.getIdNoLock() == be.getIdNoLock() {
			found = true
		}
	}
	if found == false {
		be.TraceableLock()
		hlPtr.beLocks = append(hlPtr.beLocks, be)
	}
}

func (hlPtr *HeldLocks) ReleaseLocks() {
	hlPtr.Lock()
	defer hlPtr.Unlock()
	for _, be := range hlPtr.beLocks {
		be.TraceableUnlock()
	}
	hlPtr.beLocks = nil
	hlPtr.functionCallManager.ExecuteFunctions(hlPtr.waitGroup)
}
