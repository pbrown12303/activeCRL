// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	//	"github.com/satori/go.uuid"
	"log"
	"sync"
)

// StringStringMap is a map from a URI (a string) to the UUID of an Element that represents it
type StringStringMap struct {
	sync.Mutex
	uriUUIDMap map[string]string
}

// NewStringStringMap creates and initializes a stringElementMap
func NewStringStringMap() *StringStringMap {
	var uriUUIDMap StringStringMap
	uriUUIDMap.uriUUIDMap = make(map[string]string)
	return &uriUUIDMap
}

// Clear clears the map
func (seMap *StringStringMap) Clear() {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	seMap.uriUUIDMap = make(map[string]string)
}

// CopyMap returns a copy of the map
func (seMap *StringStringMap) CopyMap() map[string]string {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	copy := make(map[string]string)
	for key, value := range seMap.uriUUIDMap {
		copy[key] = value
	}
	return copy
}

// DeleteEntry removes the map entry for the indicated UUID
func (seMap *StringStringMap) DeleteEntry(key string) {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	delete(seMap.uriUUIDMap, key)
}

// GetEntry returns the Element corresponding to the UUID
func (seMap *StringStringMap) GetEntry(key string) string {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	return seMap.uriUUIDMap[key]
}

// IsEquivalent returns true if the map contains the same number of entries
// and each key same value
func (seMap *StringStringMap) IsEquivalent(sem *StringStringMap, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) > 0 {
		print = printExceptions[0]
	}
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	sem.TraceableLock()
	defer sem.TraceableUnlock()
	if len(seMap.uriUUIDMap) != len(sem.uriUUIDMap) {
		if print {
			log.Printf("Map1 entries that differ from Map2:")
			for k := range seMap.uriUUIDMap {
				if seMap.uriUUIDMap[k] != sem.uriUUIDMap[k] {
					log.Printf("Map1 key: %s Map1 value %s Map2 value %s", k, seMap.uriUUIDMap[k], sem.uriUUIDMap[k])
				}
			}
			log.Printf("Map2 entries that differ frmo Map1")
			for k := range sem.uriUUIDMap {
				if sem.uriUUIDMap[k] != seMap.uriUUIDMap[k] {
					log.Printf("Map2 key: %s Map2 value %s Map1 value %s", k, sem.uriUUIDMap[k], seMap.uriUUIDMap[k])
				}
			}
		}
		return false
	}
	for k := range seMap.uriUUIDMap {
		if seMap.uriUUIDMap[k] != sem.uriUUIDMap[k] {
			if print {
				log.Printf("Map1 entries that differ from Map2:")
				for k := range seMap.uriUUIDMap {
					if seMap.uriUUIDMap[k] != sem.uriUUIDMap[k] {
						log.Printf("Map1 key: %s Map1 value %s Map2 value %s", k, seMap.uriUUIDMap[k], sem.uriUUIDMap[k])
					}
				}
			}
			return false
		}
	}
	return true
}

// Print prints the map. The function is intended for use in debugging
func (seMap *StringStringMap) Print(hl *HeldLocks) {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	for uri, uuid := range seMap.uriUUIDMap {
		log.Printf("Uri: %s  UUID: %s\n", uri, uuid)
	}
}

// SetEntry sets the UUID corresponding to the given URI
func (seMap *StringStringMap) SetEntry(key string, value string) {
	seMap.TraceableLock()
	defer seMap.TraceableUnlock()
	seMap.uriUUIDMap[key] = value
}

// TraceableLock locks the map. If TraceLocks is true in logs the acquisition of the lock
func (seMap *StringStringMap) TraceableLock() {
	// if TraceLocks {
	// 	log.Printf("About to lock stringElementMap %p\n", seMap)
	// }
	seMap.Lock()
}

// TraceableUnlock unlocks the map. If TraceLocks is true it logs the release of the lock
func (seMap *StringStringMap) TraceableUnlock() {
	// if TraceLocks {
	// 	log.Printf("About to unlock stringElementMap %p\n", seMap)
	// }
	seMap.Unlock()
}
