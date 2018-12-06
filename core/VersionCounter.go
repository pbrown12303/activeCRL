package core

import (
	"sync"
)

// versionCounter is simply a lockable way of counting versions
type versionCounter struct {
	sync.Mutex
	counter int
}

func newVersionCounter() *versionCounter {
	var vc versionCounter
	return &vc
}

func (vCtr *versionCounter) getVersion() int {
	vCtr.Lock()
	defer vCtr.Unlock()
	return vCtr.counter
}

func (vCtr *versionCounter) incrementVersion() int {
	vCtr.Lock()
	defer vCtr.Unlock()
	vCtr.counter++
	return vCtr.counter
}
