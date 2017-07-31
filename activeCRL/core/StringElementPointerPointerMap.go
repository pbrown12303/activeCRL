package core

import (
	"log"
	"sync"
)

type elementPointerPointerList *[]ElementPointerPointer

type StringElementPointerPointerListMap struct {
	sync.Mutex
	elementPointerPointerListMap map[string]elementPointerPointerList
}

func NewStringElementPointerPointerListMap() *StringElementPointerPointerListMap {
	var stringElementPointerPointerMap StringElementPointerPointerListMap
	stringElementPointerPointerMap.elementPointerPointerListMap = make(map[string]elementPointerPointerList)
	return &stringElementPointerPointerMap
}

func (sbeMap *StringElementPointerPointerListMap) AddEntry(key string, value ElementPointerPointer) {
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

func (sbeMap *StringElementPointerPointerListMap) RemoveEntry(key string, entry ElementPointerPointer) {
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

func (sbeMap *StringElementPointerPointerListMap) GetEntry(key string) elementPointerPointerList {
	sbeMap.TraceableLock()
	defer sbeMap.TraceableUnlock()
	return sbeMap.elementPointerPointerListMap[key]
}

func (sbeMap *StringElementPointerPointerListMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock StringElementPointerPointerListMap %p\n", sbeMap)
	}
	sbeMap.Lock()
}

func (sbeMap *StringElementPointerPointerListMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock StringElementPointerPointerListMap %p\n", sbeMap)
	}
	sbeMap.Unlock()
}
