package core

import (
	"sync"
)

type cachedPointer struct {
	sync.Mutex
	indicatedConceptID string
	isOwnerPointer     bool
	indicatedConcept   Element
	parentConceptID    string
}

func newCachedPointer(parentConceptID string, isOwnerPointer bool) *cachedPointer {
	var ptr cachedPointer
	ptr.isOwnerPointer = isOwnerPointer
	ptr.parentConceptID = parentConceptID
	return &ptr
}

func (ptr *cachedPointer) getIndicatedConcept() Element {
	ptr.Lock()
	defer ptr.Unlock()
	return ptr.indicatedConcept
}

func (ptr *cachedPointer) getIndicatedConceptID() string {
	ptr.Lock()
	defer ptr.Unlock()
	return ptr.indicatedConceptID
}

func (ptr *cachedPointer) getIsOwnerPointer() bool {
	ptr.Lock()
	defer ptr.Unlock()
	return ptr.isOwnerPointer
}

func (ptr *cachedPointer) getParentConceptID() string {
	ptr.Lock()
	defer ptr.Unlock()
	return ptr.parentConceptID
}

func (ptr *cachedPointer) isEquivalent(cp *cachedPointer) bool {
	ptr.Lock()
	defer ptr.Unlock()
	cp.Lock()
	defer cp.Unlock()
	if ptr.indicatedConceptID != cp.indicatedConceptID {
		return false
	}
	if ptr.parentConceptID != cp.parentConceptID {
		return false
	}
	if (ptr.indicatedConcept == nil && cp.indicatedConcept != nil) || (ptr.indicatedConcept != nil && cp.indicatedConcept == nil) {
		return false
	}
	if ptr.indicatedConcept != nil && cp.indicatedConcept != nil {
		if ptr.indicatedConcept.getConceptIDNoLock() != cp.indicatedConcept.getConceptIDNoLock() {
			return false
		}
	}
	return true
}

func (ptr *cachedPointer) setIndicatedConceptID(id string) {
	ptr.Lock()
	defer ptr.Unlock()
	ptr.indicatedConceptID = id
}

func (ptr *cachedPointer) setIndicatedConcept(el Element) {
	ptr.Lock()
	defer ptr.Unlock()
	ptr.indicatedConcept = el
}
