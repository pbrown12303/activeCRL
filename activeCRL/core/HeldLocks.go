package core

import (
	"sync"
)

type HeldLocks struct {
	sync.Mutex
	beLocks []BaseElement
}

func NewHeldLocks() *HeldLocks {
	var hl HeldLocks
	return &hl
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
}
