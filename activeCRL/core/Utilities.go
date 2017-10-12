// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"log"
	"reflect"
	"sync"
)

type printMutexStruct struct {
	sync.Mutex
}

var PrintMutex printMutexStruct
var TraceLocks bool = false

func clone(be BaseElement) BaseElement {
	switch be.(type) {
	case *baseElementPointer:
		return be.(*baseElementPointer).clone()
	case *baseElementReference:
		return be.(*baseElementReference).clone()
	case *element:
		return be.(*element).clone()
	case *elementPointer:
		return be.(*elementPointer).clone()
	case *elementPointerPointer:
		return be.(*elementPointerPointer).clone()
	case *elementPointerReference:
		return be.(*elementPointerReference).clone()
	case *elementReference:
		return be.(*elementReference).clone()
	case *literal:
		return be.(*literal).clone()
	case *literalPointer:
		return be.(*literalPointer).clone()
	case *literalPointerPointer:
		return be.(*literalPointerPointer).clone()
	case *literalPointerReference:
		return be.(*literalPointerReference).clone()
	case *literalReference:
		return be.(*literalReference).clone()
	case *refinement:
		return be.(*refinement).clone()
	}
	log.Printf("clone called with unhandled type %T\n", be)
	return nil
}

// CreateReplicateAsRefinement() replicates the indicated Element and all of its descendent Elements
// except that descendant Refinements are not replicated.
// For each replicated Element, a Refinement is created with the abstractElement being the original and the refinedElement
// being the replica. The root replicated element is returned.
func CreateReplicateAsRefinement(original Element, hl *HeldLocks) Element {
	if hl == nil {
		hl := NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	uOfD := original.GetUniverseOfDiscourse(hl)
	var replicate Element
	switch original.(type) {
	case BaseElementReference:
		replicate = uOfD.NewBaseElementReference(hl)
	case ElementPointerReference:
		replicate = uOfD.NewElementPointerReference(hl)
	case ElementReference:
		replicate = uOfD.NewElementReference(hl)
	case LiteralPointerReference:
		replicate = uOfD.NewLiteralPointerReference(hl)
	case LiteralReference:
		replicate = uOfD.NewLiteralReference(hl)
	case Refinement:
		replicate = uOfD.NewRefinement(hl)
	case Element:
		replicate = uOfD.NewElement(hl)
	}
	ReplicateAsRefinement(original, replicate, hl)
	return replicate
}

func Equivalent(be1 BaseElement, be2 BaseElement, hl *HeldLocks) bool {
	if hl == nil {
		hl := NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be1)
	if be2 != be1 {
		hl.LockBaseElement(be2)
	}
	return equivalent(be1, be2, hl)
}

func equivalent(be1 BaseElement, be2 BaseElement, hl *HeldLocks) bool {
	if reflect.TypeOf(be1) != reflect.TypeOf(be2) {
		return false
	}
	switch be1.(type) {
	case *baseElementPointer:
		return be1.(*baseElementPointer).isEquivalent(be2.(*baseElementPointer), hl)
	case *baseElementReference:
		return be1.(*baseElementReference).isEquivalent(be2.(*baseElementReference), hl)
	case *element:
		return be1.(*element).isEquivalent(be2.(*element), hl)
	case *elementPointer:
		return be1.(*elementPointer).isEquivalent(be2.(*elementPointer), hl)
	case *elementPointerPointer:
		return be1.(*elementPointerPointer).isEquivalent(be2.(*elementPointerPointer), hl)
	case *elementPointerReference:
		return be1.(*elementPointerReference).isEquivalent(be2.(*elementPointerReference), hl)
	case *elementReference:
		return be1.(*elementReference).isEquivalent(be2.(*elementReference), hl)
	case *literal:
		return be1.(*literal).isEquivalent(be2.(*literal), hl)
	case *literalPointer:
		return be1.(*literalPointer).isEquivalent(be2.(*literalPointer), hl)
	case *literalPointerPointer:
		return be1.(*literalPointerPointer).isEquivalent(be2.(*literalPointerPointer), hl)
	case *literalPointerReference:
		return be1.(*literalPointerReference).isEquivalent(be2.(*literalPointerReference), hl)
	case *literalReference:
		return be1.(*literalReference).isEquivalent(be2.(*literalReference), hl)
	case *refinement:
		return be1.(*refinement).isEquivalent(be2.(*refinement), hl)
	default:
		log.Printf("Equivalent default case entered for object: \n")
		Print(be1, "   ", hl)
	}
	return false
}

func GetChildBaseElementReferenceWithAncestorUri(element Element, uri string, hl *HeldLocks) BaseElementReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildElementWithAncestorUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case BaseElementReference:
			return be.(BaseElementReference)
		}
	}
	return nil
}

// GetChildElementWithAncestorUri() is a locking function that returns the first child that has a refinement ancestor
// with the indicated uri
func GetChildElementWithAncestorUri(element Element, uri string, hl *HeldLocks) Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	if uOfD == nil {
		return nil
	}
	for _, child := range element.GetOwnedElements(hl) {
		for _, ancestor := range uOfD.GetAbstractElementsRecursively(child, hl) {
			if GetUri(ancestor, hl) == uri {
				return child
			}
		}
	}
	return nil
}

func GetChildBaseElementReferenceWithUri(element Element, uri string, hl *HeldLocks) BaseElementReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildWithUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case BaseElementReference:
			return be.(BaseElementReference)
		}
	}
	return nil
}

// GetChildElementWithUri() is a locking function that returns the first child with the indicated
// uri if that child is an element
func GetChildElementWithUri(element Element, uri string, hl *HeldLocks) Element {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildWithUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case Element:
			return be.(Element)
		}
	}
	return nil
}

func GetChildElementPointerReferenceWithAncestorUri(element Element, uri string, hl *HeldLocks) ElementPointerReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildElementWithAncestorUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case ElementPointerReference:
			return be.(ElementPointerReference)
		}
	}
	return nil
}

// GetChildElementPointerReferenceWithUri() is a locking function that returns the first child
// element pointer reference with the indicated uri
func GetChildElementPointerReferenceWithUri(element Element, uri string, hl *HeldLocks) ElementPointerReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildWithUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case ElementPointerReference:
			return be.(ElementPointerReference)
		}
	}
	return nil
}

// GetChildElementReferenceWithUri() is a locking function that returns the first child
// element reference with the indicated uri
func GetChildElementReferenceWithUri(element Element, uri string, hl *HeldLocks) ElementReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildWithUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case ElementReference:
			return be.(ElementReference)
		}
	}
	return nil
}

func GetChildElementReferenceWithAncestorUri(element Element, uri string, hl *HeldLocks) ElementReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildElementWithAncestorUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case ElementReference:
			return be.(ElementReference)
		}
	}
	return nil
}

func GetChildLiteralPointerReferenceWithAncestorUri(element Element, uri string, hl *HeldLocks) LiteralPointerReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildElementWithAncestorUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case LiteralPointerReference:
			return be.(LiteralPointerReference)
		}
	}
	return nil
}

func GetChildLiteralPointerReferenceWithUri(element Element, uri string, hl *HeldLocks) LiteralPointerReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildWithUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case LiteralPointerReference:
			return be.(LiteralPointerReference)
		}
	}
	return nil
}

func GetChildLiteralReferenceWithAncestorUri(element Element, uri string, hl *HeldLocks) LiteralReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildElementWithAncestorUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case LiteralReference:
			return be.(LiteralReference)
		}
	}
	return nil
}

func GetChildLiteralReferenceWithUri(element Element, uri string, hl *HeldLocks) LiteralReference {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	be := GetChildWithUri(element, uri, hl)
	if be != nil {
		switch be.(type) {
		case LiteralReference:
			return be.(LiteralReference)
		}
	}
	return nil
}

// GetChildWithUri() is a locking function that returns the first child with the indicated
// uri
func GetChildWithUri(element Element, uri string, hl *HeldLocks) BaseElement {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(element)
	for _, child := range element.GetOwnedBaseElements(hl) {
		if GetUri(child, hl) == uri {
			return child
		}
	}
	return nil
}

func Print(be BaseElement, prefix string, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	printBe(be, prefix, hl)
}

func printBe(be BaseElement, prefix string, hl *HeldLocks) {
	if be == nil {
		return
	}
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(be)
	log.Printf("%s%s: %p\n", prefix, reflect.TypeOf(be).String(), be)
	switch be.(type) {
	case *baseElementPointer:
		be.(*baseElementPointer).printBaseElementPointer(prefix, hl)
	case *baseElementReference:
		be.(*baseElementReference).printBaseElementReference(prefix, hl)
	case *element:
		be.(*element).printElement(prefix, hl)
	case *elementPointer:
		be.(*elementPointer).printElementPointer(prefix, hl)
	case *elementPointerPointer:
		be.(*elementPointerPointer).printElementPointerPointer(prefix, hl)
	case *elementPointerReference:
		be.(*elementPointerReference).printElementPointerReference(prefix, hl)
	case *elementReference:
		be.(*elementReference).printElementReference(prefix, hl)
	case *literal:
		be.(*literal).printLiteral(prefix, hl)
	case *literalPointer:
		be.(*literalPointer).printLiteralPointer(prefix, hl)
	case *literalPointerPointer:
		be.(*literalPointerPointer).printLiteralPointerPointer(prefix, hl)
	case *literalPointerReference:
		be.(*literalPointerReference).printLiteralPointerReference(prefix, hl)
	case *literalReference:
		be.(*literalReference).printLiteralReference(prefix, hl)
	case *refinement:
		be.(*refinement).printRefinement(prefix, hl)
	default:
		log.Printf("No case for %T in Print \n", be)
	}
}

func PrintUriIndex(uOfD UniverseOfDiscourse, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	uOfD.(*universeOfDiscourse).uriBaseElementMap.Print(hl)
}

func PrintUriIndexJustIdentifiers(uOfD UniverseOfDiscourse, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	uOfD.(*universeOfDiscourse).uriBaseElementMap.PrintJustIdentifiers(hl)
}

// ReplicateAsRefinement() replicates the structure of the original in the replicate, ignoring
// Refinements and Values. The name from each original element is copied into the name of the
// corresponding replicate element. This function is idempotent: if applied to an existing structure,
// Elements of that structure that have existing Refinement relationships with original Elements
// will not be re-created.
func ReplicateAsRefinement(original Element, replicate Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	hl.LockBaseElement(original)
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	SetName(replicate, GetName(original, hl), hl)
	if uOfD.IsRefinementOf(replicate, original, hl) == false {
		refinement := uOfD.NewRefinement(hl)
		SetOwningElement(refinement, replicate, hl)
		refinement.SetAbstractElement(original, hl)
		refinement.SetRefinedElement(replicate, hl)
	}

	for _, originalChild := range original.GetOwnedElements(hl) {
		var replicateChild Element
		for _, currentChild := range replicate.GetOwnedElements(hl) {
			for _, currentChildAncestor := range uOfD.GetAbstractElementsRecursively(currentChild, hl) {
				if currentChildAncestor == originalChild {
					replicateChild = currentChild
				}
			}
		}
		if replicateChild == nil {
			switch originalChild.(type) {
			case BaseElementReference:
				replicateChild = uOfD.NewBaseElementReference(hl)
			case ElementPointerReference:
				replicateChild = uOfD.NewElementPointerReference(hl)
			case ElementReference:
				replicateChild = uOfD.NewElementReference(hl)
			case LiteralPointerReference:
				replicateChild = uOfD.NewLiteralPointerReference(hl)
			case LiteralReference:
				replicateChild = uOfD.NewLiteralReference(hl)
			case Element:
				replicateChild = uOfD.NewElement(hl)
			}
			SetOwningElement(replicateChild, replicate, hl)
			refinement := uOfD.NewRefinement(hl)
			SetOwningElement(refinement, replicateChild, hl)
			refinement.SetAbstractElement(originalChild, hl)
			refinement.SetRefinedElement(replicateChild, hl)
			SetName(replicateChild, GetName(originalChild, hl), hl)
		}
		switch originalChild.(type) {
		case Element:
			ReplicateAsRefinement(originalChild.(Element), replicateChild.(Element), hl)
		}
	}

}

func restoreValueOwningElementFieldsRecursively(el Element, hl *HeldLocks) {
	if hl == nil {
		hl = NewHeldLocks(nil)
		defer hl.ReleaseLocks()
	}
	for _, child := range el.GetOwnedBaseElements(hl) {
		switch child.(type) {
		//@TODO add reference to case
		case *baseElementPointer:
			child.(*baseElementPointer).internalSetOwningElement(el, hl)
		case *baseElementReference:
			restoreValueOwningElementFieldsRecursively(child.(*baseElementReference), hl)
		case *element:
			restoreValueOwningElementFieldsRecursively(child.(*element), hl)
		case *elementPointer:
			child.(*elementPointer).internalSetOwningElement(el, hl)
		case *elementPointerPointer:
			child.(*elementPointerPointer).internalSetOwningElement(el, hl)
		case *elementPointerReference:
			restoreValueOwningElementFieldsRecursively(child.(*elementPointerReference), hl)
		case *elementReference:
			restoreValueOwningElementFieldsRecursively(child.(*elementReference), hl)
		case *literal:
			child.(*literal).internalSetOwningElement(el, hl)
		case *literalPointer:
			child.(*literalPointer).internalSetOwningElement(el, hl)
		case *literalPointerPointer:
			child.(*literalPointerPointer).internalSetOwningElement(el, hl)
		case *literalPointerReference:
			restoreValueOwningElementFieldsRecursively(child.(*literalPointerReference), hl)
		case *literalReference:
			restoreValueOwningElementFieldsRecursively(child.(*literalReference), hl)
		case *refinement:
			restoreValueOwningElementFieldsRecursively(child.(*refinement), hl)
		default:
			log.Printf("No case for %T in restoreValueOwningElementFieldsRecursively \n", child)
		}
	}
}

func unmarshalPolymorphicBaseElement(data []byte, result *BaseElement) error {
	var unmarshaledData map[string]json.RawMessage
	err := json.Unmarshal(data, &unmarshaledData)
	//	fmt.Printf("%s \n", unmarshaledData)
	var elementType string
	err = json.Unmarshal(unmarshaledData["Type"], &elementType)
	//	fmt.Printf("%s \n", elementType)

	switch elementType {
	case "*core.baseElementPointer":
		//		fmt.Printf("Switch choice *core.baseElementPointer \n")
		var recoveredBaseElementPointer baseElementPointer
		*result = &recoveredBaseElementPointer
		err = recoveredBaseElementPointer.recoverBaseElementPointerFields(&unmarshaledData)
	case "*core.baseElementReference":
		var recoveredBaseElementReference baseElementReference
		recoveredBaseElementReference.ownedBaseElements = make(map[uuid.UUID]BaseElement)
		*result = &recoveredBaseElementReference
		err = recoveredBaseElementReference.recoverBaseElementReferenceFields(&unmarshaledData)
		if err != nil {
			return err
		}
	case "*core.element":
		//		fmt.Printf("Switch choice *core.element \n")
		var recoveredElement element
		recoveredElement.ownedBaseElements = make(map[uuid.UUID]BaseElement)
		*result = &recoveredElement
		err = recoveredElement.recoverElementFields(&unmarshaledData)
	case "*core.elementPointer":
		//		fmt.Printf("Switch choice *core.elementPointer \n")
		var recoveredElementPointer elementPointer
		*result = &recoveredElementPointer
		err = recoveredElementPointer.recoverElementPointerFields(&unmarshaledData)
	case "*core.elementPointerPointer":
		//		fmt.Printf("Switch choice *core.elementPointerPointer \n")
		var recoveredElementPointerPointer elementPointerPointer
		*result = &recoveredElementPointerPointer
		err = recoveredElementPointerPointer.recoverElementPointerPointerFields(&unmarshaledData)
	case "*core.elementPointerReference":
		//		fmt.Printf("Switch choice *core.elementPointerReference \n")
		var recoveredElement elementPointerReference
		recoveredElement.ownedBaseElements = make(map[uuid.UUID]BaseElement)
		*result = &recoveredElement
		err = recoveredElement.recoverElementPointerReferenceFields(&unmarshaledData)
		if err != nil {
			return err
		}
	case "*core.elementReference":
		//		fmt.Printf("Switch choice *core.elementReference \n")
		var recoveredElement elementReference
		recoveredElement.ownedBaseElements = make(map[uuid.UUID]BaseElement)
		*result = &recoveredElement
		err = recoveredElement.recoverElementReferenceFields(&unmarshaledData)
		if err != nil {
			return err
		}
	case "*core.literal":
		//		fmt.Printf("Switch choice *core.literal \n")
		var recoveredLiteral literal
		*result = &recoveredLiteral
		err = recoveredLiteral.recoverLiteralFields(&unmarshaledData)
	case "*core.literalPointer":
		//		fmt.Printf("Switch choice *core.literalPointer \n")
		var recoveredLiteralPointer literalPointer
		*result = &recoveredLiteralPointer
		err = recoveredLiteralPointer.recoverLiteralPointerFields(&unmarshaledData)
	case "*core.literalPointerPointer":
		//		fmt.Printf("Switch choice *core.literalPointerPointer \n")
		var recoveredLiteralPointerPointer literalPointerPointer
		*result = &recoveredLiteralPointerPointer
		err = recoveredLiteralPointerPointer.recoverLiteralPointerPointerFields(&unmarshaledData)
	case "*core.literalPointerReference":
		//		fmt.Printf("Switch choice *core.literalPointerReference \n")
		var recoveredElement literalPointerReference
		recoveredElement.ownedBaseElements = make(map[uuid.UUID]BaseElement)
		*result = &recoveredElement
		err = recoveredElement.recoverLiteralPointerReferenceFields(&unmarshaledData)
		if err != nil {
			return err
		}
	case "*core.literalReference":
		//		fmt.Printf("Switch choice *core.literalPointer \n")
		var recoveredLiteralReference literalReference
		recoveredLiteralReference.ownedBaseElements = make(map[uuid.UUID]BaseElement)
		*result = &recoveredLiteralReference
		err = recoveredLiteralReference.recoverLiteralReferenceFields(&unmarshaledData)
	case "*core.refinement":
		var recoveredRefinement refinement
		recoveredRefinement.ownedBaseElements = make(map[uuid.UUID]BaseElement)
		*result = &recoveredRefinement
		err = recoveredRefinement.recoverRefinementFields(&unmarshaledData)
	default:
		log.Printf("No case for %s in unmarshalPolymorphicBaseElement \n", elementType)
	}
	return err
}
