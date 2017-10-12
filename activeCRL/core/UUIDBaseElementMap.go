// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

type UUIDBaseElementMap struct {
	sync.Mutex
	baseElementMap map[uuid.UUID]BaseElement
}

func NewUUIDBaseElementMap() *UUIDBaseElementMap {
	var uuidBaseElementMap UUIDBaseElementMap
	uuidBaseElementMap.baseElementMap = make(map[uuid.UUID]BaseElement)
	return &uuidBaseElementMap
}

func (sbeMap *UUIDBaseElementMap) GetRange() []BaseElement {
	var baseElements []BaseElement
	for _, be := range sbeMap.baseElementMap {
		baseElements = append(baseElements, be)
	}
	return baseElements
}

func (sbeMap *UUIDBaseElementMap) DeleteEntry(key uuid.UUID) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	delete(sbeMap.baseElementMap, key)
}

func (sbeMap *UUIDBaseElementMap) GetEntry(key uuid.UUID) BaseElement {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.baseElementMap[key]
}

func (sbeMap *UUIDBaseElementMap) Print(hl *HeldLocks) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	for uuid, be := range sbeMap.baseElementMap {
		log.Printf("Uri: %s\n", uuid.String())
		Print(be, "    ", hl)
	}
}

func (sbeMap *UUIDBaseElementMap) PrintJustIdentifiers(hl *HeldLocks) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	for uuid, _ := range sbeMap.baseElementMap {
		log.Printf("UUID: %s \n", uuid.String())
	}
}

func (sbeMap *UUIDBaseElementMap) SetEntry(key uuid.UUID, value BaseElement) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	sbeMap.baseElementMap[key] = value
}

func (sbeMap *UUIDBaseElementMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock UUIDBaseElementMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *UUIDBaseElementMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock UUIDBaseElementMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
