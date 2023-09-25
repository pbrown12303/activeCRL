// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"log"
	"reflect"
	// "sync"
)

// type printMutexStruct struct {
// 	sync.Mutex
// }

// GetConceptTypeString returns the string representing the reflected type
func GetConceptTypeString(el Concept) string {
	if el == nil {
		return ""
	}
	return reflect.TypeOf(el).String()
}

// // PrintMutex provides a mututal exclusion for print routhines shared across threads
// var PrintMutex printMutexStruct

func clone(el Concept, trans *Transaction) Concept {
	return el.(*concept).clone(trans)
}

// Equivalent returns true if the two elements are equivalent
func Equivalent(be1 Concept, hl1 *Transaction, be2 Concept, hl2 *Transaction, printExceptions ...bool) bool {
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
func RecursivelyEquivalent(e1 Concept, hl1 *Transaction, e2 Concept, hl2 *Transaction, printExceptions ...bool) bool {
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
			it.Stop()
			return false
		}
	}
	return true
}

func equivalent(be1 Concept, hl1 *Transaction, be2 Concept, hl2 *Transaction, printExceptions ...bool) bool {
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
	case *concept:
		return be1.(*concept).isEquivalent(hl1, be2.(*concept), hl2, print)
	default:
		log.Printf("Equivalent default case entered for object: \n")
		Print(be1, "   ", hl1)
	}
	return false
}

// Print prints the indicated element and its ownedConcepts, recursively
func Print(el Concept, prefix string, trans *Transaction) {
	printElement(el, prefix, trans)
}

func printElement(el Concept, prefix string, trans *Transaction) {
	if el == nil {
		return
	}
	trans.ReadLockElement(el)
	serializedElement, _ := el.MarshalJSON()
	log.Printf("%s%s", prefix, string(serializedElement))
	uOfD := el.GetUniverseOfDiscourse(trans)
	if uOfD == nil {
		return
	}
	ownedIDsMap := uOfD.ownedIDsMap
	if ownedIDsMap == nil {
		return
	}
	ownedIDs := ownedIDsMap.GetMappedValues(el.GetConceptID(trans))
	it := ownedIDs.Iterator()
	for id := range it.C {
		ownedElement := el.GetUniverseOfDiscourse(trans).GetElement(id.(string))
		printElement(ownedElement, prefix+"  ", trans)
	}
}

// PrintURIIndex prints the URI index of the uOfD with full Element information
func PrintURIIndex(uOfD *UniverseOfDiscourse, trans *Transaction) {
	uOfD.uriUUIDMap.Print(trans)
}
