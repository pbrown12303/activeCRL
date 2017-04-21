package core

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func Equivalent(be1 BaseElement, be2 BaseElement) bool {
	//	be1.Lock()
	//	defer be1.Unlock()
	//	be2.Lock()
	//	defer be2.Unlock()
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
	be.Lock()
	defer be.Unlock()
	printBe(be, prefix)
}

func printBe(be BaseElement, prefix string) {
	if be == nil {
		return
	}
	fmt.Printf("%s%s: \n", prefix, reflect.TypeOf(be).String())
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
			child.(*elementPointer).setOwningElement(el)
		case *elementPointerPointer:
			child.(*elementPointerPointer).setOwningElement(el)
		case *elementPointerReference:
			restoreValueOwningElementFieldsRecursively(child.(*elementPointerReference))
		case *elementReference:
			restoreValueOwningElementFieldsRecursively(child.(*elementReference))
		case *literal:
			child.(*literal).setOwningElement(el)
		case *literalPointer:
			child.(*literalPointer).setOwningElement(el)
		case *literalPointerPointer:
			child.(*literalPointerPointer).setOwningElement(el)
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
