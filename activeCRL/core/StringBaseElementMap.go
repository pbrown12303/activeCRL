// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
)

type StringBaseElementMap struct {
	sync.Mutex
	baseElementMap map[string]BaseElement
}

func NewStringBaseElementMap() *StringBaseElementMap {
	var stringBaseElementMap StringBaseElementMap
	stringBaseElementMap.baseElementMap = make(map[string]BaseElement)
	return &stringBaseElementMap
}

func (sbeMap *StringBaseElementMap) GetRange() []BaseElement {
	var baseElements []BaseElement
	for _, be := range sbeMap.baseElementMap {
		baseElements = append(baseElements, be)
	}
	return baseElements
}

func (sbeMap *StringBaseElementMap) DeleteEntry(key string) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	delete(sbeMap.baseElementMap, key)
}

func (sbeMap *StringBaseElementMap) GetEntry(key string) BaseElement {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.baseElementMap[key]
}

func (sbeMap *StringBaseElementMap) Print(hl *HeldLocks) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	for uri, be := range sbeMap.baseElementMap {
		log.Printf("Uri: %s\n", uri)
		Print(be, "    ", hl)
	}
}

func (sbeMap *StringBaseElementMap) PrintJustIdentifiers(hl *HeldLocks) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	for uri, be := range sbeMap.baseElementMap {
		log.Printf("Uri: %s Id: %s\n", uri, be.GetId(hl).String())
	}
}

func (sbeMap *StringBaseElementMap) SetEntry(key string, value BaseElement) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	sbeMap.baseElementMap[key] = value
}

func (sbeMap *StringBaseElementMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock StringBaseElementMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *StringBaseElementMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock StringBaseElementMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
