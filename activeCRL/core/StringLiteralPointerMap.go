package core

import (
	"log"
	"sync"
)

type literalPointerList *[]LiteralPointer

type StringLiteralPointerListMap struct {
	sync.Mutex
	literalPointerListMap map[string]literalPointerList
}

func NewStringLiteralPointerListMap() *StringLiteralPointerListMap {
	var stringLiteralPointerMap StringLiteralPointerListMap
	stringLiteralPointerMap.literalPointerListMap = make(map[string]literalPointerList)
	return &stringLiteralPointerMap
}

func (sbeMap *StringLiteralPointerListMap) AddEntry(key string, value LiteralPointer) {
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

func (sbeMap *StringLiteralPointerListMap) RemoveEntry(key string, entry LiteralPointer) {
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

func (sbeMap *StringLiteralPointerListMap) GetEntry(key string) literalPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.literalPointerListMap[key]
}

func (sbeMap *StringLiteralPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock StringLiteralPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *StringLiteralPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock StringLiteralPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
