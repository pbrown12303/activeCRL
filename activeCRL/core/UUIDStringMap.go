// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
)

type UUIDStringMap struct {
	sync.Mutex
	stringMap map[string]string
}

func NewUUIDStringMap() *UUIDStringMap {
	var uuidStringMap UUIDStringMap
	uuidStringMap.stringMap = make(map[string]string)
	return &uuidStringMap
}

func (ssMap *UUIDStringMap) GetRange() []string {
	var strings []string
	for _, be := range ssMap.stringMap {
		strings = append(strings, be)
	}
	return strings
}

func (ssMap *UUIDStringMap) DeleteEntry(key string) {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	delete(ssMap.stringMap, key)
}

func (ssMap *UUIDStringMap) GetEntry(key string) string {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	return ssMap.stringMap[key]
}

func (ssMap *UUIDStringMap) SetEntry(key string, value string) {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	ssMap.stringMap[key] = value
}

func (ssMap *UUIDStringMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock UUIDStringMap %p\n", ssMap)
	}
	ssMap.Lock()
}

func (ssMap *UUIDStringMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock UUIDStringMap %p\n", ssMap)
	}
	ssMap.Unlock()
}
