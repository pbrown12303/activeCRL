package core

import (
	"log"
	"sync"
)

type StringStringMap struct {
	sync.Mutex
	stringMap map[string]string
}

func NewStringStringMap() *StringStringMap {
	var stringStringMap StringStringMap
	stringStringMap.stringMap = make(map[string]string)
	return &stringStringMap
}

func (ssMap *StringStringMap) GetRange() []string {
	var strings []string
	for _, be := range ssMap.stringMap {
		strings = append(strings, be)
	}
	return strings
}

func (ssMap *StringStringMap) DeleteEntry(key string) {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	delete(ssMap.stringMap, key)
}

func (ssMap *StringStringMap) GetEntry(key string) string {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	return ssMap.stringMap[key]
}

func (ssMap *StringStringMap) SetEntry(key string, value string) {
	ssMap.TraceableLock()
	defer ssMap.TraceableUnlock()
	ssMap.stringMap[key] = value
}

func (ssMap *StringStringMap) TraceableLock() {
	if TraceLocks {
		log.Printf("About to lock StringStringMap %p\n", ssMap)
	}
	ssMap.Lock()
}

func (ssMap *StringStringMap) TraceableUnlock() {
	if TraceLocks {
		log.Printf("About to unlock StringStringMap %p\n", ssMap)
	}
	ssMap.Unlock()
}
