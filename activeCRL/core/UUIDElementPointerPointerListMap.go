// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

type elementPointerPointerList *[]ElementPointerPointer

type UUIDElementPointerPointerListMap struct {
	sync.Mutex
	elementPointerPointerListMap map[uuid.UUID]elementPointerPointerList
}

func NewUUIDElementPointerPointerListMap() *UUIDElementPointerPointerListMap {
	var uuidElementPointerPointerMap UUIDElementPointerPointerListMap
	uuidElementPointerPointerMap.elementPointerPointerListMap = make(map[uuid.UUID]elementPointerPointerList)
	return &uuidElementPointerPointerMap
}

func (sbeMap *UUIDElementPointerPointerListMap) AddEntry(key uuid.UUID, value ElementPointerPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.elementPointerPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == value {
				// element is already in list
				return
			}
		}
	}
	if currentList == nil {
		var newList [1]ElementPointerPointer
		newList[0] = value
		newSlice := newList[:]
		sbeMap.elementPointerPointerListMap[key] = &newSlice
	} else {
		updatedList := append(*currentList, value)
		sbeMap.elementPointerPointerListMap[key] = &updatedList
	}
}

func (sbeMap *UUIDElementPointerPointerListMap) RemoveEntry(key uuid.UUID, entry ElementPointerPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.elementPointerPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == entry {
				copy((*currentList)[i:], (*currentList)[i+1:])
				updatedList := (*currentList)[:len(*currentList)-1]
				sbeMap.elementPointerPointerListMap[key] = &updatedList
				return
			}
		}
	}
}

func (sbeMap *UUIDElementPointerPointerListMap) GetEntry(key uuid.UUID) elementPointerPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.elementPointerPointerListMap[key]
}

func (sbeMap *UUIDElementPointerPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock UUIDElementPointerPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *UUIDElementPointerPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock UUIDElementPointerPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
