// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"log"
	"reflect"
	"runtime/debug"
	"sync"
)

type printMutexStruct struct {
	sync.Mutex
}

// PrintMutex provides a mututal exclusion for print routhines shared across threads
var PrintMutex printMutexStruct

func clone(el Element, hl *HeldLocks) Element {
	switch el.(type) {
	case *literal:
		return el.(*literal).clone(hl)
	case *reference:
		return el.(*reference).clone(hl)
	case *refinement:
		return el.(*refinement).clone(hl)
	case *element:
		return el.(*element).clone(hl)
	}
	log.Printf("clone called with unhandled type %T\n", el)
	debug.PrintStack()
	return nil
}

// Equivalent returns true if the two elements are equivalent
func Equivalent(be1 Element, hl1 *HeldLocks, be2 Element, hl2 *HeldLocks) bool {
	if be1 == nil && be2 == nil {
		return true
	}
	if (be1 == nil && be2 != nil) || (be1 != nil && be2 == nil) {
		return false
	}
	hl1.ReadLockElement(be1)
	if be2 != be1 {
		hl2.ReadLockElement(be2)
	}
	return equivalent(be1, hl1, be2, hl2)
}

func equivalent(be1 Element, hl1 *HeldLocks, be2 Element, hl2 *HeldLocks) bool {
	if reflect.TypeOf(be1) != reflect.TypeOf(be2) {
		return false
	}
	switch be1.(type) {
	case *element:
		return be1.(*element).isEquivalent(hl1, be2.(*element), hl2)
	case *reference:
		return be1.(*reference).isEquivalent(hl1, be2.(*reference), hl2)
	case *literal:
		return be1.(*literal).isEquivalent(hl1, be2.(*literal), hl2)
	case *refinement:
		return be1.(*refinement).isEquivalent(hl1, be2.(*refinement), hl2)
	default:
		log.Printf("Equivalent default case entered for object: \n")
		Print(be1, "   ", hl1)
	}
	return false
}

// Print prints the indicated element and its ownedConcepts, recursively
func Print(el Element, prefix string, hl *HeldLocks) {
	printBe(el, prefix, hl)
}

func printBe(el Element, prefix string, hl *HeldLocks) {
	// if el == nil {
	// 	return
	// }
	// if hl == nil {
	// 	hl = NewHeldLocks(nil)
	// 	defer hl.ReleaseLocks()
	// }
	// hl.LockElement(el)
	// log.Printf("%s%s: %p\n", prefix, reflect.TypeOf(el).String(), el)
	// switch el.(type) {
	// case *element:
	// 	el.(*element).printElement(prefix, hl)
	// case *reference:
	// 	el.(*reference).printReference(prefix, hl)
	// case *literal:
	// 	el.(*literal).printLiteral(prefix, hl)
	// case *refinement:
	// 	el.(*refinement).printRefinement(prefix, hl)
	// case *universeOfDiscourse:
	// 	el.(*universeOfDiscourse).printElement(prefix, hl)
	// default:
	// 	log.Printf("No case for %T in Print \n", el)
	// }
}

// PrintURIIndex prints the URI index of the uOfD with full Element information
func PrintURIIndex(uOfD *UniverseOfDiscourse, hl *HeldLocks) {
	uOfD.uriUUIDMap.Print(hl)
}

func restoreValueOwningElementFieldsRecursively(el Element, hl *HeldLocks) {
	// if hl == nil {
	// 	hl = NewHeldLocks(nil)
	// 	defer hl.ReleaseLocks()
	// }
	// for _, child := range el.GetOwnedConcepts(hl) {
	// 	switch child.(type) {
	// 	//@TODO add reference to case
	// 	case *element:
	// 		restoreValueOwningElementFieldsRecursively(child.(*element), hl)
	// 	case *reference:
	// 		restoreValueOwningElementFieldsRecursively(child.(*reference), hl)
	// 	case *literal:
	// 		child.(*literal).internalSetOwningElement(el, hl)
	// 	case *refinement:
	// 		restoreValueOwningElementFieldsRecursively(child.(*refinement), hl)
	// 	default:
	// 		log.Printf("No case for %T in restoreValueOwningElementFieldsRecursively \n", child)
	// 	}
	// }
}
