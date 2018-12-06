// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
)

// StringStringMap is used to map UUIDs and URIs to the identifiers of the Elements with which they are associated
type StringStringMap struct {
	sync.Mutex
	stringMap map[string]string
}

// NewStringStringMap creates and initializes a stringStringMap that is lockable
func NewStringStringMap() *StringStringMap {
	var stringStringMap StringStringMap
	stringStringMap.stringMap = make(map[string]string)
	return &stringStringMap
}

// GetRange returns the range of the StringStringMap
func (ssMap *StringStringMap) GetRange() []string {
	var strings []string
	for _, be := range ssMap.stringMap {
		strings = append(strings, be)
	}
	return strings
}

// DeleteEntry removes an entry fron the StringStringMap
func (ssMap *StringStringMap) DeleteEntry(key string) {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	delete(ssMap.stringMap, key)
}

// GetEntry gets an entry fron the StringStringMap
func (ssMap *StringStringMap) GetEntry(key string) string {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	return ssMap.stringMap[key]
}

// SetEntry sets an entry in the StringStringMap
func (ssMap *StringStringMap) SetEntry(key string, value string) {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	ssMap.stringMap[key] = value
}

// TraceableLock locks ithe StringStringMap in a traceable manner
func (ssMap *StringStringMap) TraceableLock() {
	// if TraceLocks {
	// 	log.Printf("About to lock stringStringMap %p\n", ssMap)
	// }
	ssMap.Lock()
}

// TraceableUnlock unlocks ithe StringStringMap in a traceable manner
func (ssMap *StringStringMap) TraceableUnlock() {
	// if TraceLocks {
	// 	log.Printf("About to unlock stringStringMap %p\n", ssMap)
	// }
	ssMap.Unlock()
}
