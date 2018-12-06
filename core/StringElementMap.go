// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	//	"github.com/satori/go.uuid"
	"log"
	"sync"
)

// StringElementMap is a map from a UUID (a string) to the Element that represents it
type StringElementMap struct {
	sync.Mutex
	elementMap map[string]Element
}

// NewStringElementMap creates and initializes a stringElementMap
func NewStringElementMap() *StringElementMap {
	var uuidElementMap StringElementMap
	uuidElementMap.elementMap = make(map[string]Element)
	return &uuidElementMap
}

// CopyMap returns a copy of the map
func (seMap *StringElementMap) CopyMap() *map[string]Element {
	copy := make(map[string]Element)
	for key, value := range seMap.elementMap {
		copy[key] = value
	}
	return &copy
}

// DeleteEntry removes the map entry for the indicated UUID
func (seMap *StringElementMap) DeleteEntry(key string) {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	delete(seMap.elementMap, key)
}

// GetEntry returns the Element corresponding to the UUID
func (seMap *StringElementMap) GetEntry(key string) Element {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	return seMap.elementMap[key]
}

// IsEquivalent returns true if the map contains the same number of elements
// and each has the same set of keys. No comparison is performed on the elements
// themselves
func (seMap *StringElementMap) IsEquivalent(sem *StringElementMap) bool {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	sem.TraceableLock()
	defer sem.TraceableUnlock()
	if len(seMap.elementMap) != len(sem.elementMap) {
		return false
	}
	for k := range seMap.elementMap {
		if sem.elementMap[k] == nil {
			return false
		}
	}
	return true
}

// Print prints the map. The function is intended for use in debugging
func (seMap *StringElementMap) Print(hl *HeldLocks) {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	for uuid, be := range seMap.elementMap {
		log.Printf("Uri: %s\n", uuid)
		Print(be, "    ", hl)
	}
}

// PrintJustIdentifiers prints just the UUIDs (keys) of the map. It is intended for use in debugging
func (seMap *StringElementMap) PrintJustIdentifiers(hl *HeldLocks) {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	for uuid := range seMap.elementMap {
		log.Printf("UUID: %s \n", uuid)
	}
}

// SetEntry sets the Element corresponding to the given UUID
func (seMap *StringElementMap) SetEntry(key string, value Element) {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	seMap.elementMap[key] = value
}

// TraceableLock locks the map. If TraceLocks is true in logs the acquisition of the lock
func (seMap *StringElementMap) TraceableLock() {
	// if TraceLocks {
	// 	log.Printf("About to lock stringElementMap %p\n", seMap)
	// }
	seMap.Lock()
}

// TraceableUnlock unlocks the map. If TraceLocks is true it logs the release of the lock
func (seMap *StringElementMap) TraceableUnlock() {
	// if TraceLocks {
	// 	log.Printf("About to unlock stringElementMap %p\n", seMap)
	// }
	seMap.Unlock()
}
