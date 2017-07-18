package core

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
)

type printMutexStruct struct {
	sync.Mutex
}

var PrintMutex printMutexStruct

func Equivalent(be1 BaseElement, be2 BaseElement) bool {
	be1.TraceableLock()
	defer be1.TraceableUnlock()
	if be2 != be1 {
		be2.TraceableLock()
		defer be2.TraceableUnlock()
	}
	return equivalent(be1, be2)
}

func equivalent(be1 BaseElement, be2 BaseElement) bool {
	if reflect.TypeOf(be1) != reflect.TypeOf(be2) {
		return false
	}
	switch be1.(type) {
	case *baseElementPointer:
		return be1.(*baseElementPointer).isEquivalent(be2.(*baseElementPointer))
	case *baseElementReference:
		return be1.(*baseElementReference).isEquivalent(be2.(*baseElementReference))
	case *element:
		return be1.(*element).isEquivalent(be2.(*element))
	case *elementPointer:
		return be1.(*elementPointer).isEquivalent(be2.(*elementPointer))
	case *elementPointerPointer:
		return be1.(*elementPointerPointer).isEquivalent(be2.(*elementPointerPointer))
	case *elementPointerReference:
		return be1.(*elementPointerReference).isEquivalent(be2.(*elementPointerReference))
	case *elementReference:
		return be1.(*elementReference).isEquivalent(be2.(*elementReference))
	case *literal:
		return be1.(*literal).isEquivalent(be2.(*literal))
	case *literalPointer:
		return be1.(*literalPointer).isEquivalent(be2.(*literalPointer))
	case *literalPointerPointer:
		return be1.(*literalPointerPointer).isEquivalent(be2.(*literalPointerPointer))
	case *literalPointerReference:
		return be1.(*literalPointerReference).isEquivalent(be2.(*literalPointerReference))
	case *literalReference:
		return be1.(*literalReference).isEquivalent(be2.(*literalReference))
	case *refinement:
		return be1.(*refinement).isEquivalent(be2.(*refinement))
	default:
		log.Printf("Equivalent default case entered for object: \n")
		Print(be1, "   ")
	}
	return false
}

// GetChildWithUri() is a locking function that returns the first child with the indicated
// uri
func GetChildWithUri(element Element, uri string) BaseElement {
	element.TraceableLock()
	defer element.TraceableUnlock()
	return GetChildWithUriNoLock(element, uri)
}

// GetChildWithUriNoLock() is a non-locking function that returns the first child with the indicated
// uri
func GetChildWithUriNoLock(element Element, uri string) BaseElement {
	for _, child := range element.getOwnedBaseElements() {
		if child.GetUri() == uri {
			return child
		}
	}
	return nil
}

// GetChildElementWithAncestorUri() is a locking function that returns the first child that has a refinement ancestor
// with the indicated uri
func GetChildElementWithAncestorUri(element Element, uri string) Element {
	element.TraceableLock()
	defer element.TraceableUnlock()
	return GetChildElementWithAncestorUriNoLock(element, uri)
}

// GetChildElementWithAncestorUriNoLock() is a non-locking function that returns the first child that has a refinement
// ancestor with the indicated uri
func GetChildElementWithAncestorUriNoLock(element Element, uri string) Element {
	for _, child := range element.GetOwnedElementsNoLock() {
		for _, ancestor := range child.GetAbstractElementsRecursivelyNoLock() {
			if ancestor.GetUri() == uri {
				return child
			}
		}
	}
	return nil
}

// GetChildElementWithUri() is a locking function that returns the first child with the indicated
// uri if that child is an element
func GetChildElementWithUri(element Element, uri string) Element {
	element.TraceableLock()
	defer element.TraceableUnlock()
	return GetChildElementWithUriNoLock(element, uri)
}

// GetChildElementWithUriNoLock() is a non-locking function that returns the first child with the indicated
// uri if that child is an element
func GetChildElementWithUriNoLock(element Element, uri string) Element {
	be := GetChildWithUriNoLock(element, uri)
	if be != nil {
		switch be.(type) {
		case Element:
			return be.(Element)
		}
	}
	return nil
}

// GetChildElementReferenceWithUri() is a locking function that returns the first child
// element reference with the indicated uri
func GetChildElementReferenceWithUri(element Element, uri string) ElementReference {
	element.TraceableLock()
	defer element.TraceableUnlock()
	return GetChildElementReferenceWithUriNoLock(element, uri)
}

// GetChildElementReferenceWithUriNoLock() is a non-locking function that returns the first child
// element reference with the indicated uri
func GetChildElementReferenceWithUriNoLock(element Element, uri string) ElementReference {
	be := GetChildWithUriNoLock(element, uri)
	if be != nil {
		switch be.(type) {
		case ElementReference:
			return be.(ElementReference)
		}
	}
	return nil
}

// GetChildElementReferenceWithAncestorUri() is a locking function that returns the first child
// element reference with an ancestor having the indicated uri
func GetChildElementReferenceWithAncestorUri(element Element, uri string) ElementReference {
	element.TraceableLock()
	defer element.TraceableUnlock()
	return GetChildElementReferenceWithUriNoLock(element, uri)
}

// GetChildElementReferenceWithAncestorUriNoLock() is a non-locking function that returns the first child
// element reference with the indicated uri
func GetChildElementReferenceWithAncestorUriNoLock(element Element, uri string) ElementReference {
	be := GetChildElementWithAncestorUriNoLock(element, uri)
	if be != nil {
		switch be.(type) {
		case ElementReference:
			return be.(ElementReference)
		}
	}
	return nil
}

func Print(be BaseElement, prefix string) {
	printBe(be, prefix)
}

func PrintNotification(notification *ChangeNotification) {
	printNotification(notification, "")
}

func printNotification(notification *ChangeNotification, prefix string) {
	notificationType := ""
	switch notification.natureOfChange {
	case ADD:
		notificationType = "Add"
	case MODIFY:
		notificationType = "Modify"
	case REMOVE:
		notificationType = "Remove"
	}
	log.Printf("%s%s: \n", prefix, notificationType)
	Print(notification.changedObject, prefix+"   ")
	if notification.underlyingChange != nil {
		printNotification(notification.underlyingChange, prefix+"      ")
	}
}

func PrintUndoStack(s undoStack, stackName string) {
	log.Printf("%s:", stackName)
	for _, entry := range s {
		var changeType string
		switch entry.changeType {
		case Creation:
			{
				changeType = "Creation"
			}
		case Deletion:
			{
				changeType = "Deletion"
			}
		case Change:
			{
				changeType = "Change"
			}
		case Marker:
			{
				changeType = "Marker"
			}
		}
		log.Printf("   Change type: %s", changeType)
		log.Printf("   Prior state:")
		Print(entry.priorState, "      ")
		//		log.Printf("   Changed element:")
		//		Print(entry.changedElement, "      ")
	}
}

func printBe(be BaseElement, prefix string) {
	if be == nil {
		return
	}
	log.Printf("%s%s: \n", prefix, reflect.TypeOf(be).String())
	switch be.(type) {
	case *baseElementPointer:
		be.(*baseElementPointer).printBaseElementPointer(prefix)
	case *baseElementReference:
		be.(*baseElementReference).printBaseElementReference(prefix)
	case *element:
		be.(*element).printElement(prefix)
	case *elementPointer:
		be.(*elementPointer).printElementPointer(prefix)
	case *elementPointerPointer:
		be.(*elementPointerPointer).printElementPointerPointer(prefix)
	case *elementPointerReference:
		be.(*elementPointerReference).printElementPointerReference(prefix)
	case *elementReference:
		be.(*elementReference).printElementReference(prefix)
	case *literal:
		be.(*literal).printLiteral(prefix)
	case *literalPointer:
		be.(*literalPointer).printLiteralPointer(prefix)
	case *literalPointerPointer:
		be.(*literalPointerPointer).printLiteralPointerPointer(prefix)
	case *literalPointerReference:
		be.(*literalPointerReference).printLiteralPointerReference(prefix)
	case *literalReference:
		be.(*literalReference).printLiteralReference(prefix)
	case *refinement:
		be.(*refinement).printRefinement(prefix)
	default:
		log.Printf("No case for %T in Print \n", be)
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
		//		fmt.Printf("Switch choice *core.baseElementReference \n")
		var recoveredBaseElementReference baseElementReference
		recoveredBaseElementReference.ownedBaseElements = make(map[string]BaseElement)
		*result = &recoveredBaseElementReference
		err = recoveredBaseElementReference.recoverBaseElementReferenceFields(&unmarshaledData)
		if err != nil {
			return err
		}
	case "*core.element":
		//		fmt.Printf("Switch choice *core.element \n")
		var recoveredElement element
		recoveredElement.ownedBaseElements = make(map[string]BaseElement)
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
		recoveredElement.ownedBaseElements = make(map[string]BaseElement)
		*result = &recoveredElement
		err = recoveredElement.recoverElementPointerReferenceFields(&unmarshaledData)
		if err != nil {
			return err
		}
	case "*core.elementReference":
		//		fmt.Printf("Switch choice *core.elementReference \n")
		var recoveredElement elementReference
		recoveredElement.ownedBaseElements = make(map[string]BaseElement)
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
		recoveredElement.ownedBaseElements = make(map[string]BaseElement)
		*result = &recoveredElement
		err = recoveredElement.recoverLiteralPointerReferenceFields(&unmarshaledData)
		if err != nil {
			return err
		}
	case "*core.literalReference":
		//		fmt.Printf("Switch choice *core.literalPointer \n")
		var recoveredLiteralReference literalReference
		recoveredLiteralReference.ownedBaseElements = make(map[string]BaseElement)
		*result = &recoveredLiteralReference
		err = recoveredLiteralReference.recoverLiteralReferenceFields(&unmarshaledData)
	case "*core.refinement":
		var recoveredRefinement refinement
		recoveredRefinement.ownedBaseElements = make(map[string]BaseElement)
		*result = &recoveredRefinement
		err = recoveredRefinement.recoverRefinementFields(&unmarshaledData)
	default:
		log.Printf("No case for %s in unmarshalPolymorphicBaseElement \n", elementType)
	}
	return err
}

func restoreValueOwningElementFieldsRecursively(el Element) {
	for _, child := range el.getOwnedBaseElements() {
		switch child.(type) {
		//@TODO add reference to case
		case *baseElementPointer:
			child.(*baseElementPointer).internalSetOwningElement(el)
		case *baseElementReference:
			restoreValueOwningElementFieldsRecursively(child.(*baseElementReference))
		case *element:
			restoreValueOwningElementFieldsRecursively(child.(*element))
		case *elementPointer:
			child.(*elementPointer).internalSetOwningElement(el)
		case *elementPointerPointer:
			child.(*elementPointerPointer).internalSetOwningElement(el)
		case *elementPointerReference:
			restoreValueOwningElementFieldsRecursively(child.(*elementPointerReference))
		case *elementReference:
			restoreValueOwningElementFieldsRecursively(child.(*elementReference))
		case *literal:
			child.(*literal).internalSetOwningElement(el)
		case *literalPointer:
			child.(*literalPointer).internalSetOwningElement(el)
		case *literalPointerPointer:
			child.(*literalPointerPointer).internalSetOwningElement(el)
		case *literalPointerReference:
			restoreValueOwningElementFieldsRecursively(child.(*literalPointerReference))
		case *literalReference:
			restoreValueOwningElementFieldsRecursively(child.(*literalReference))
		case *refinement:
			restoreValueOwningElementFieldsRecursively(child.(*refinement))
		default:
			log.Printf("No case for %T in restoreValueOwningElementFieldsRecursively \n", child)
		}
	}
}

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

var TraceLocks bool = false
