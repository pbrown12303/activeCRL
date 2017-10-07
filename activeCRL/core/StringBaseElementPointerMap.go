// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
)

type baseElementPointerList *[]BaseElementPointer

type StringBaseElementPointerListMap struct {
	sync.Mutex
	baseElementPointerListMap map[string]baseElementPointerList
}

func NewStringBaseElementPointerListMap() *StringBaseElementPointerListMap {
	var stringBaseElementPointerMap StringBaseElementPointerListMap
	stringBaseElementPointerMap.baseElementPointerListMap = make(map[string]baseElementPointerList)
	return &stringBaseElementPointerMap
}

func (sbeMap *StringBaseElementPointerListMap) AddEntry(key string, value BaseElementPointer) {
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

func (sbeMap *StringBaseElementPointerListMap) RemoveEntry(key string, entry BaseElementPointer) {
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

func (sbeMap *StringBaseElementPointerListMap) GetEntry(key string) baseElementPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.baseElementPointerListMap[key]
}

func (sbeMap *StringBaseElementPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock StringBaseElementPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *StringBaseElementPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock StringBaseElementPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
