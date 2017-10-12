// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

type literalPointerPointerList *[]LiteralPointerPointer

type UUIDLiteralPointerPointerListMap struct {
	sync.Mutex
	literalPointerPointerListMap map[uuid.UUID]literalPointerPointerList
}

func NewUUIDLiteralPointerPointerListMap() *UUIDLiteralPointerPointerListMap {
	var uuidLiteralPointerPointerMap UUIDLiteralPointerPointerListMap
	uuidLiteralPointerPointerMap.literalPointerPointerListMap = make(map[uuid.UUID]literalPointerPointerList)
	return &uuidLiteralPointerPointerMap
}

func (sbeMap *UUIDLiteralPointerPointerListMap) AddEntry(key uuid.UUID, value LiteralPointerPointer) {
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

func (sbeMap *UUIDLiteralPointerPointerListMap) RemoveEntry(key uuid.UUID, entry LiteralPointerPointer) {
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

func (sbeMap *UUIDLiteralPointerPointerListMap) GetEntry(key uuid.UUID) literalPointerPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.literalPointerPointerListMap[key]
}

func (sbeMap *UUIDLiteralPointerPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock UUIDLiteralPointerPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *UUIDLiteralPointerPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock UUIDLiteralPointerPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
