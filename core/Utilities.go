// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can el found in the LICENSE file.

package core

import (
	"fmt"
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

// CreateReplicateAsRefinement replicates the indicated Element and all of its descendent Elements
// except that descendant Refinements are not replicated.
// For each replicated Element, a Refinement is created with the abstractElement being the original and the refinedElement
// being the replica. The root replicated element is returned.
func CreateReplicateAsRefinement(original Element, hl *HeldLocks) Element {
	uOfD := original.GetUniverseOfDiscourse(hl)
	var replicate Element
	switch original.(type) {
	case Literal:
		replicate, _ = uOfD.NewLiteral(hl)
	case Reference:
		replicate, _ = uOfD.NewReference(hl)
	case Refinement:
		replicate, _ = uOfD.NewRefinement(hl)
	case Element:
		replicate, _ = uOfD.NewElement(hl)
	}
	ReplicateAsRefinement(original, replicate, hl)
	return replicate
}

// CreateReplicateAsRefinementFromURI replicates the Element indicated by the URI
func CreateReplicateAsRefinementFromURI(uOfD UniverseOfDiscourse, originalURI string, hl *HeldLocks) (Element, error) {
	original := uOfD.GetElementWithURI(originalURI)
	if original == nil {
		return nil, fmt.Errorf("In CreateReplicateAsRefinementFromURI Element with uri %s not found", originalURI)
	}
	return CreateReplicateAsRefinement(original, hl), nil
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
func PrintURIIndex(uOfD UniverseOfDiscourse, hl *HeldLocks) {
	uOfD.(*universeOfDiscourse).uriElementMap.Print(hl)
}

// PrintURIIndexJustIdentifiers prints the URI indix with just identifiers
func PrintURIIndexJustIdentifiers(uOfD UniverseOfDiscourse, hl *HeldLocks) {
	uOfD.(*universeOfDiscourse).uriElementMap.PrintJustIdentifiers(hl)
}

// ReplicateAsRefinement replicates the structure of the original in the replicate, ignoring
// Refinements and Values. The name from each original element is copied into the name of the
// corresponding replicate element. This function is idempotent: if applied to an existing structure,
// Elements of that structure that have existing Refinement relationships with original Elements
// will not el re-created.
func ReplicateAsRefinement(original Element, replicate Element, hl *HeldLocks) {
	// if hl == nil {
	// 	hl = NewHeldLocks(nil)
	// 	defer hl.ReleaseLocks()
	// }
	// hl.LockElement(original)
	// hl.LockElement(replicate)
	// uOfD := replicate.GetUniverseOfDiscourse(hl)

	// SetLabel(replicate, GetLabel(original, hl), hl)
	// if uOfD.IsRefinementOf(replicate, original, hl) == false {
	// 	refinement := uOfD.NewRefinement(hl)
	// 	SetOwningElement(refinement, replicate, hl)
	// 	refinement.SetAbstractElement(original, hl)
	// 	refinement.SetRefinedElement(replicate, hl)
	// }

	// for _, originalChild := range original.GetOwnedElements(hl) {
	// 	var replicateChild Element
	// 	for _, currentChild := range replicate.GetOwnedElements(hl) {
	// 		for _, currentChildAncestor := range uOfD.GetAbstractElementsRecursively(currentChild, hl) {
	// 			if currentChildAncestor == originalChild {
	// 				replicateChild = currentChild
	// 			}
	// 		}
	// 	}
	// 	if replicateChild == nil {
	// 		switch originalChild.(type) {
	// 		case BaseElementReference:
	// 			replicateChild = uOfD.NewBaseElementReference(hl)
	// 		case ElementPointerReference:
	// 			replicateChild = uOfD.NewElementPointerReference(hl)
	// 		case ElementReference:
	// 			replicateChild = uOfD.NewElementReference(hl)
	// 		case LiteralPointerReference:
	// 			replicateChild = uOfD.NewLiteralPointerReference(hl)
	// 		case LiteralReference:
	// 			replicateChild = uOfD.NewLiteralReference(hl)
	// 		case Element:
	// 			replicateChild = uOfD.NewElement(hl)
	// 		}
	// 		SetOwningElement(replicateChild, replicate, hl)
	// 		refinement := uOfD.NewRefinement(hl)
	// 		SetOwningElement(refinement, replicateChild, hl)
	// 		refinement.SetAbstractElement(originalChild, hl)
	// 		refinement.SetRefinedElement(replicateChild, hl)
	// 		SetLabel(replicateChild, GetLabel(originalChild, hl), hl)
	// 	}
	// 	switch originalChild.(type) {
	// 	case Element:
	// 		ReplicateAsRefinement(originalChild.(Element), replicateChild.(Element), hl)
	// 	}
	// }

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
