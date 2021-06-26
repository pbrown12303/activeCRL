// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"log"
	"reflect"
	"runtime/debug"
	// "sync"
)

// type printMutexStruct struct {
// 	sync.Mutex
// }

// GetConceptTypeString returns the string representing the reflected type
func GetConceptTypeString(el Element) string {
	if el == nil {
		return ""
	}
	return reflect.TypeOf(el).String()
}

// // PrintMutex provides a mututal exclusion for print routhines shared across threads
// var PrintMutex printMutexStruct

func clone(el Element, hl *Transaction) Element {
	switch typedEl := el.(type) {
	case *literal:
		return typedEl.clone(hl)
	case *reference:
		return typedEl.clone(hl)
	case *refinement:
		return typedEl.clone(hl)
	case *element:
		return typedEl.clone(hl)
	}
	log.Printf("clone called with unhandled type %T\n", el)
	debug.PrintStack()
	return nil
}

// Equivalent returns true if the two elements are equivalent
func Equivalent(be1 Element, hl1 *Transaction, be2 Element, hl2 *Transaction, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) > 0 {
		print = printExceptions[0]
	}
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
	return equivalent(be1, hl1, be2, hl2, print)
}

// RecursivelyEquivalent returns true if two elements and all of their children are equivalent
func RecursivelyEquivalent(e1 Element, hl1 *Transaction, e2 Element, hl2 *Transaction, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) == 1 {
		print = printExceptions[0]
	}
	if !Equivalent(e1, hl1, e2, hl2, print) {
		if print {
			log.Print("Equivalence failed")
			Print(e1, "e1: ", hl1)
			Print(e2, "e2: ", hl2)
		}
		return false
	}
	children1 := e1.GetOwnedConceptIDs(hl1)
	children2 := e2.GetOwnedConceptIDs(hl2)
	if children1.Cardinality() != children2.Cardinality() {
		if print {
			log.Print("Children Cardinality failed")
			Print(e1, "Element1:", hl1)
			log.Printf("Children1: %s", children1.String())
			log.Printf("Children2: %s", children2.String())
		}
		return false
	}
	if !children1.Equal(children2) {
		if print {
			log.Print("Children Equal failed")
			Print(e1, "Element1:", hl1)
			log.Printf("Children1: %s", children1.String())
			log.Printf("Children2: %s", children2.String())
		}
		return false
	}
	it := children1.Iterator()
	for childEntry := range it.C {
		childID := childEntry.(string)
		child1 := e1.GetUniverseOfDiscourse(hl1).GetElement(childID)
		child2 := e2.GetUniverseOfDiscourse(hl2).GetElement(childID)
		if child1 == nil || child2 == nil || !RecursivelyEquivalent(child1, hl1, child2, hl2, print) {
			return false
		}
	}
	return true
}

func equivalent(be1 Element, hl1 *Transaction, be2 Element, hl2 *Transaction, printExceptions ...bool) bool {
	var print bool
	if len(printExceptions) > 0 {
		print = printExceptions[0]
	}
	if reflect.TypeOf(be1) != reflect.TypeOf(be2) {
		if print {
			log.Printf("In equivalent, element types do not match")
		}
		return false
	}
	switch be1.(type) {
	case *element:
		return be1.(*element).isEquivalent(hl1, be2.(*element), hl2, print)
	case *reference:
		return be1.(*reference).isEquivalent(hl1, be2.(*reference), hl2, print)
	case *literal:
		return be1.(*literal).isEquivalent(hl1, be2.(*literal), hl2, print)
	case *refinement:
		return be1.(*refinement).isEquivalent(hl1, be2.(*refinement), hl2, print)
	// case *UniverseOfDiscourse:
	// 	return be1.(*UniverseOfDiscourse).element.isEquivalent(hl1, &be2.(*UniverseOfDiscourse).element, hl2, print)
	default:
		log.Printf("Equivalent default case entered for object: \n")
		Print(be1, "   ", hl1)
	}
	return false
}

// Print prints the indicated element and its ownedConcepts, recursively
func Print(el Element, prefix string, hl *Transaction) {
	printElement(el, prefix, hl)
}

func printElement(el Element, prefix string, hl *Transaction) {
	if el == nil {
		return
	}
	hl.ReadLockElement(el)
	serializedElement, _ := el.MarshalJSON()
	log.Printf("%s%s", prefix, string(serializedElement))
	ownedIDs := el.GetUniverseOfDiscourse(hl).ownedIDsMap.GetMappedValues(el.GetConceptID(hl))
	it := ownedIDs.Iterator()
	defer it.Stop()
	for id := range it.C {
		ownedElement := el.GetUniverseOfDiscourse(hl).GetElement(id.(string))
		printElement(ownedElement, prefix+"  ", hl)
	}
}

// PrintURIIndex prints the URI index of the uOfD with full Element information
func PrintURIIndex(uOfD *UniverseOfDiscourse, hl *Transaction) {
	uOfD.uriUUIDMap.Print(hl)
}

// func restoreValueOwningElementFieldsRecursively(el Element, hl *HeldLocks) {
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
// }
