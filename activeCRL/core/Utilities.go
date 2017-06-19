package core

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func Equivalent(be1 BaseElement, be2 BaseElement) bool {
	be1.traceableLock()
	defer be1.traceableUnlock()
	if be2 != be1 {
		be2.traceableLock()
		defer be2.traceableUnlock()
	}
	return equivalent(be1, be2)
}

func equivalent(be1 BaseElement, be2 BaseElement) bool {
	if reflect.TypeOf(be1) != reflect.TypeOf(be2) {
		return false
	}
	switch be1.(type) {
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

func Print(be BaseElement, prefix string) {
	//	if be != nil {
	//		be.traceableLock()
	//		defer be.traceableUnlock()
	//	}
	printBe(be, prefix)
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

func RecoverElement(data []byte, uOfD *UniverseOfDiscourse) Element {
	if len(data) == 0 {
		return nil
	}
	var recoveredElement BaseElement
	err := unmarshalPolymorphicBaseElement(data, &recoveredElement)
	//	fmt.Printf("Recovered Element: \n")
	//	Print(recoveredElement, "   ")
	if err != nil {
		fmt.Printf("Error recovering Element: %s \n", err)
		return nil
	}
	uOfD.SetUniverseOfDiscourseRecursively(recoveredElement)
	restoreValueOwningElementFieldsRecursively(recoveredElement.(Element))
	return recoveredElement.(Element)
}

func unmarshalPolymorphicBaseElement(data []byte, result *BaseElement) error {
	var unmarshaledData map[string]json.RawMessage
	err := json.Unmarshal(data, &unmarshaledData)
	//	fmt.Printf("%s \n", unmarshaledData)
	var elementType string
	err = json.Unmarshal(unmarshaledData["Type"], &elementType)
	//	fmt.Printf("%s \n", elementType)

	switch elementType {
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

func preChange(be BaseElement) {
	if be != nil && be.getUniverseOfDiscourse().recordingUndo == true {
		be.getUniverseOfDiscourse().markChangedBaseElement(be)
	}
}

func postChange(be BaseElement) {
	be.internalIncrementVersion()
	parent := be.getOwningElement()
	if parent != nil {
		parent.childChanged()
	}
}

var TraceLocks bool = false
