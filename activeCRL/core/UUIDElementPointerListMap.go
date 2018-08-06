// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
)

type elementPointerList *[]ElementPointer

type UUIDElementPointerListMap struct {
	sync.Mutex
	elementPointerListMap map[string]elementPointerList
}

func NewUUIDElementPointerListMap() *UUIDElementPointerListMap {
	var uuidElementPointerMap UUIDElementPointerListMap
	uuidElementPointerMap.elementPointerListMap = make(map[string]elementPointerList)
	return &uuidElementPointerMap
}

func (sbeMap *UUIDElementPointerListMap) AddEntry(key string, value ElementPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.elementPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == value {
				// element is already in list
				return
			}
		}
	}
	if currentList == nil {
		var newList [1]ElementPointer
		newList[0] = value
		newSlice := newList[:]
		sbeMap.elementPointerListMap[key] = &newSlice
	} else {
		updatedList := append(*currentList, value)
		sbeMap.elementPointerListMap[key] = &updatedList
	}
}

func (sbeMap *UUIDElementPointerListMap) RemoveEntry(key string, entry ElementPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.elementPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == entry {
				copy((*currentList)[i:], (*currentList)[i+1:])
				updatedList := (*currentList)[:len(*currentList)-1]
				sbeMap.elementPointerListMap[key] = &updatedList
				return
			}
		}
	}
}

func (sbeMap *UUIDElementPointerListMap) GetEntry(key string) elementPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.elementPointerListMap[key]
}

func (sbeMap *UUIDElementPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock UUIDElementPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *UUIDElementPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock UUIDElementPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
