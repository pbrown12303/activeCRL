package core

import (
	"log"
	"sync"
)

type elementPointerList *[]ElementPointer

type StringElementPointerListMap struct {
	sync.Mutex
	elementPointerListMap map[string]elementPointerList
}

func NewStringElementPointerListMap() *StringElementPointerListMap {
	var stringElementPointerMap StringElementPointerListMap
	stringElementPointerMap.elementPointerListMap = make(map[string]elementPointerList)
	return &stringElementPointerMap
}

func (sbeMap *StringElementPointerListMap) AddEntry(key string, value ElementPointer) {
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

func (sbeMap *StringElementPointerListMap) RemoveEntry(key string, entry ElementPointer) {
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

func (sbeMap *StringElementPointerListMap) GetEntry(key string) elementPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.elementPointerListMap[key]
}

func (sbeMap *StringElementPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock StringElementPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *StringElementPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock StringElementPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
