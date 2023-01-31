// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	//	"github.com/satori/go.uuid"

	"sync"

	mapset "github.com/deckarep/golang-set"
)

// OneToNStringMap is a map from a string (a string) to a set of strings. In most usages the string is a UUID
type OneToNStringMap struct {
	sync.Mutex
	oneToNStringMap map[string]mapset.Set
}

// NewOneToNStringMap creates and initializes a stringElementMap
func NewOneToNStringMap() *OneToNStringMap {
	var newMap OneToNStringMap
	newMap.oneToNStringMap = make(map[string]mapset.Set)
	return &newMap
}

// Clear clears the map
func (onMap *OneToNStringMap) Clear() {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	onMap.oneToNStringMap = make(map[string]mapset.Set)
}

// CopyMap returns a copy of the map
func (onMap *OneToNStringMap) CopyMap() map[string]mapset.Set {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	copy := make(map[string]mapset.Set)
	for key, childSet := range onMap.oneToNStringMap {
		copy[key] = childSet.Clone()
	}
	return copy
}

// DeleteKey removes the map entry for the indicated UUID
func (onMap *OneToNStringMap) DeleteKey(key string) {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	delete(onMap.oneToNStringMap, key)
}

// GetMappedValues returns the set of strings corresponding to the key
func (onMap *OneToNStringMap) GetMappedValues(key string) mapset.Set {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	return onMap.getMappedValuesNoLock(key)
}

// getMappedValuesNoLock returns the Element corresponding to the UUID
func (onMap *OneToNStringMap) getMappedValuesNoLock(key string) mapset.Set {
	if onMap.oneToNStringMap[key] == nil {
		onMap.oneToNStringMap[key] = mapset.NewSet()
	}
	return onMap.oneToNStringMap[key].Clone()
}

// IsEquivalent returns true if the map contains the same number of elements
// and each has the same set of keys. No comparison is performed on the elements
// themselves
func (onMap *OneToNStringMap) IsEquivalent(sem *OneToNStringMap) bool {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	sem.TraceableLock()
	defer sem.TraceableUnlock()
	for k := range onMap.oneToNStringMap {
		if !sem.getMappedValuesNoLock(k).Equal(onMap.oneToNStringMap[k]) {
			return false
		}
	}
	return len(onMap.oneToNStringMap) == len(sem.oneToNStringMap)
}

// AddMappedValue adds the Element UUID as a child of the owner UUID
func (onMap *OneToNStringMap) AddMappedValue(key string, value string) {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	if onMap.oneToNStringMap[key] == nil {
		onMap.oneToNStringMap[key] = mapset.NewSet()
	}
	onMap.oneToNStringMap[key].Add(value)
}

// ContainsMappedValue adds the Element UUID as a child of the owner UUID
func (onMap *OneToNStringMap) ContainsMappedValue(key string, value string) bool {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	if onMap.oneToNStringMap[key] == nil {
		onMap.oneToNStringMap[key] = mapset.NewSet()
	}
	return onMap.oneToNStringMap[key].Contains(value)
}

// RemoveMappedValue removes the Element UUID as a child of the owner UUID
func (onMap *OneToNStringMap) RemoveMappedValue(key string, value string) {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	if onMap.oneToNStringMap[key] == nil {
		return
	}
	onMap.oneToNStringMap[key].Remove(value)
}

// SetMappedValues sets the mapped values for the given key
func (onMap *OneToNStringMap) SetMappedValues(key string, newValues mapset.Set) {
	onMap.TraceableLock()
	defer onMap.TraceableUnlock()
	if onMap.oneToNStringMap[key] == nil {
		onMap.oneToNStringMap[key] = mapset.NewSet()
	}
	selectedMap := onMap.oneToNStringMap[key]
	selectedMap.Clear()
	newValuesIterator := newValues.Iterator()
	for value := range newValuesIterator.C {
		selectedMap.Add(value)
	}
}

// TraceableLock locks the map. If TraceLocks is true in logs the acquisition of the lock
func (onMap *OneToNStringMap) TraceableLock() {
	// if TraceLocks {
	// 	log.Printf("About to lock stringElementMap %p\n", onMap)
	// }
	onMap.Lock()
}

// TraceableUnlock unlocks the map. If TraceLocks is true it logs the release of the lock
func (onMap *OneToNStringMap) TraceableUnlock() {
	// if TraceLocks {
	// 	log.Printf("About to unlock stringElementMap %p\n", onMap)
	// }
	onMap.Unlock()
}
