// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

type literalPointerList *[]LiteralPointer

type UUIDLiteralPointerListMap struct {
	sync.Mutex
	literalPointerListMap map[uuid.UUID]literalPointerList
}

func NewUUIDLiteralPointerListMap() *UUIDLiteralPointerListMap {
	var uuidLiteralPointerMap UUIDLiteralPointerListMap
	uuidLiteralPointerMap.literalPointerListMap = make(map[uuid.UUID]literalPointerList)
	return &uuidLiteralPointerMap
}

func (sbeMap *UUIDLiteralPointerListMap) AddEntry(key uuid.UUID, value LiteralPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.literalPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == value {
				// literal is already in list
				return
			}
		}
	}
	if currentList == nil {
		var newList [1]LiteralPointer
		newList[0] = value
		newSlice := newList[:]
		sbeMap.literalPointerListMap[key] = &newSlice
	} else {
		updatedList := append(*currentList, value)
		sbeMap.literalPointerListMap[key] = &updatedList
	}
}

func (sbeMap *UUIDLiteralPointerListMap) RemoveEntry(key uuid.UUID, entry LiteralPointer) {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	currentList := sbeMap.literalPointerListMap[key]
	if currentList != nil && len(*currentList) > 0 {
		for i := range *currentList {
			if (*currentList)[i] == entry {
				copy((*currentList)[i:], (*currentList)[i+1:])
				updatedList := (*currentList)[:len(*currentList)-1]
				sbeMap.literalPointerListMap[key] = &updatedList
				return
			}
		}
	}
}

func (sbeMap *UUIDLiteralPointerListMap) GetEntry(key uuid.UUID) literalPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.literalPointerListMap[key]
}

func (sbeMap *UUIDLiteralPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock UUIDLiteralPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *UUIDLiteralPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock UUIDLiteralPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
