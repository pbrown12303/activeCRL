// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
)

type literalPointerPointerList *[]LiteralPointerPointer

type StringLiteralPointerPointerListMap struct {
	sync.Mutex
	literalPointerPointerListMap map[string]literalPointerPointerList
}

func NewStringLiteralPointerPointerListMap() *StringLiteralPointerPointerListMap {
	var stringLiteralPointerPointerMap StringLiteralPointerPointerListMap
	stringLiteralPointerPointerMap.literalPointerPointerListMap = make(map[string]literalPointerPointerList)
	return &stringLiteralPointerPointerMap
}

func (sbeMap *StringLiteralPointerPointerListMap) AddEntry(key string, value LiteralPointerPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.literalPointerPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == value {
				// literal is already in list
				return
			}
		}
	}
	if currentList == nil {
		var newList [1]LiteralPointerPointer
		newList[0] = value
		newSlice := newList[:]
		sbeMap.literalPointerPointerListMap[key] = &newSlice
	} else {
		updatedList := append(*currentList, value)
		sbeMap.literalPointerPointerListMap[key] = &updatedList
	}
}

func (sbeMap *StringLiteralPointerPointerListMap) RemoveEntry(key string, entry LiteralPointerPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.literalPointerPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == entry {
				copy((*currentList)[i:], (*currentList)[i+1:])
				updatedList := (*currentList)[:len(*currentList)-1]
				sbeMap.literalPointerPointerListMap[key] = &updatedList
				return
			}
		}
	}
}

func (sbeMap *StringLiteralPointerPointerListMap) GetEntry(key string) literalPointerPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.literalPointerPointerListMap[key]
}

func (sbeMap *StringLiteralPointerPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock StringLiteralPointerPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *StringLiteralPointerPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock StringLiteralPointerPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
