// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

type baseElementPointerList *[]BaseElementPointer

type UUIDBaseElementPointerListMap struct {
	sync.Mutex
	baseElementPointerListMap map[uuid.UUID]baseElementPointerList
}

func NewUUIDBaseElementPointerListMap() *UUIDBaseElementPointerListMap {
	var uuidBaseElementPointerListMap UUIDBaseElementPointerListMap
	uuidBaseElementPointerListMap.baseElementPointerListMap = make(map[uuid.UUID]baseElementPointerList)
	return &uuidBaseElementPointerListMap
}

func (sbeMap *UUIDBaseElementPointerListMap) AddEntry(key uuid.UUID, value BaseElementPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.baseElementPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == value {
				// element is already in list
				return
			}
		}
	}
	if currentList == nil {
		var newList [1]BaseElementPointer
		newList[0] = value
		newSlice := newList[:]
		sbeMap.baseElementPointerListMap[key] = &newSlice
	} else {
		updatedList := append(*currentList, value)
		sbeMap.baseElementPointerListMap[key] = &updatedList
	}
}

func (sbeMap *UUIDBaseElementPointerListMap) RemoveEntry(key uuid.UUID, entry BaseElementPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.baseElementPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == entry {
				copy((*currentList)[i:], (*currentList)[i+1:])
				updatedList := (*currentList)[:len(*currentList)-1]
				sbeMap.baseElementPointerListMap[key] = &updatedList
				return
			}
		}
	}
}

func (sbeMap *UUIDBaseElementPointerListMap) GetEntry(key uuid.UUID) baseElementPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.baseElementPointerListMap[key]
}

func (sbeMap *UUIDBaseElementPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock UUIDBaseElementPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *UUIDBaseElementPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock UUIDBaseElementPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
