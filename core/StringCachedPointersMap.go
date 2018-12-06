package core

import (
	"errors"
	"sync"
)

// stringCachedPointersMap is used to keep track of unresolved pointers: Elements that are
// indicated by references or refinements using their ID but are not yet present in the Universe of Discourse.
// The key to the map is the ConceptID of the Element to which there are unresolved pointers
type stringCachedPointersMap struct {
	sync.Mutex
	scpMap map[string][]*cachedPointer
}

func newStringCachedPointersMap() *stringCachedPointersMap {
	var scpMap stringCachedPointersMap
	scpMap.scpMap = make(map[string][]*cachedPointer)
	return &scpMap
}

func (scpMapPtr *stringCachedPointersMap) addCachedPointer(cp *cachedPointer) error {
	if cp == nil {
		return errors.New("stringCachedPointersMap.addCachedPointer called with nil cachedPointer")
	}
	id := cp.getIndicatedConceptID()
	scpMapPtr.scpMap[id] = append(scpMapPtr.scpMap[id], cp)
	return nil
}

func (scpMapPtr *stringCachedPointersMap) resolveCachedPointers(el Element, hl *HeldLocks) error {
	if el == nil {
		return errors.New("stringCachedPointersMap.resolveCachedPointers called with nil Element")
	}
	id := el.getConceptIDNoLock()
	for _, cp := range scpMapPtr.scpMap[id] {
		cp.setIndicatedConcept(el)
		if cp.getIsOwnerPointer() == true {
			el.addRecoveredOwnedConcept(cp.parentConceptID, hl)
		} else {
			el.addListener(cp.parentConceptID, hl)
			// TODO: Add logic to handle differences between the referencedConcept version and the actual version
		}
	}
	delete(scpMapPtr.scpMap, id)
	return nil
}
