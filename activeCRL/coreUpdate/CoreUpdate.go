package main

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"io/ioutil"
	"log"
	"os"
)

var serializedCorePrefix string = `package core

import ()

var serializedCore string = `

func main() {
	uOfD := core.NewUniverseOfDiscourse()
	recoveredCore := uOfD.GetCoreConceptSpace()
	updatedCore := updateRecoveredCore(recoveredCore, uOfD)
	marshaledCore, err := updatedCore.MarshalJSON()
	if err == nil {
		ioutil.WriteFile("CoreConceptSpace.acrl", marshaledCore, os.ModePerm)
		serializedCore := serializedCorePrefix + "`" + string(marshaledCore) + "`"
		ioutil.WriteFile("./core/SerializedCore.go", []byte(serializedCore), os.ModePerm)
	}
}

func updateRecoveredCore(recoveredCore core.Element, uOfD *core.UniverseOfDiscourse) core.Element {
	// Core
	coreElement := recoveredCore
	if coreElement == nil {
		coreElement = uOfD.NewElement()
		coreElement.SetName("CoreConceptSpace")
		coreElement.SetUri(core.CoreConceptSpaceUri)
	}

	// Element
	re := uOfD.GetBaseElementWithUri(core.ElememtUri)
	var element core.Element
	if re != nil {
		var ok bool
		element, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url http://activeCrl.com/core/Element")
			return nil
		}
	}
	if element == nil {
		element = uOfD.NewElement()
		element.SetOwningElement(coreElement)
		element.SetName("Element")
		element.SetUri(core.ElememtUri)
	}

	// ElementPointer
	var elementPointer core.ElementPointer
	re = uOfD.GetBaseElementWithUri(core.ElementPointerUri)
	if re != nil {
		var ok bool
		elementPointer, ok = re.(core.ElementPointer)
		if !ok {
			log.Printf("Recovered object not of type ElementPointer with url http://activeCrl.com/core/ElementPointer")
			return nil
		}
	}
	if elementPointer == nil {
		elementPointer = uOfD.NewReferencedElementPointer()
		elementPointer.SetOwningElement(coreElement)
		elementPointer.SetUri(core.ElementPointerUri)
	}

	// ElementPointerPointer
	var elementPointerPointer core.ElementPointerPointer
	re = uOfD.GetBaseElementWithUri(core.ElementPointerPointerUri)
	if re != nil {
		var ok bool
		elementPointerPointer, ok = re.(core.ElementPointerPointer)
		if !ok {
			log.Printf("Recovered object not of type ElementPointerPointer with url http://activeCrl.com/core/ElementPointerPointer")
			return nil
		}
	}
	if elementPointerPointer == nil {
		elementPointerPointer = uOfD.NewElementPointerPointer()
		elementPointerPointer.SetOwningElement(coreElement)
		elementPointerPointer.SetUri(core.ElementPointerPointerUri)
	}

	// ElementPointerReference
	var elementPointerReference core.ElementPointerReference
	re = uOfD.GetBaseElementWithUri(core.ElementPointerReferenceUri)
	if re != nil {
		var ok bool
		elementPointerReference, ok = re.(core.ElementPointerReference)
		if !ok {
			log.Printf("Recovered object not of type ElementPointerReference with url http://activeCrl.com/core/ElementPointerRefernce")
			return nil
		}
	}
	if elementPointerReference == nil {
		elementPointerReference = uOfD.NewElementPointerReference()
		elementPointerReference.SetOwningElement(coreElement)
		elementPointerReference.SetName("ElementReference")
		elementPointerReference.SetUri(core.ElementPointerReferenceUri)
	}

	// ElementReference
	var elementReference core.ElementReference
	re = uOfD.GetBaseElementWithUri(core.ElementReferenceUri)
	if re != nil {
		var ok bool
		elementReference, ok = re.(core.ElementReference)
		if !ok {
			log.Printf("Recovered object not of type ElementReference with url http://activeCrl.com/core/ElementRefernce")
			return nil
		}
	}
	if elementReference == nil {
		elementReference = uOfD.NewElementReference()
		elementReference.SetOwningElement(coreElement)
		elementReference.SetName("ElementReference")
		elementReference.SetUri(core.ElementReferenceUri)
	}

	// Literal
	var literal core.Literal
	re = uOfD.GetBaseElementWithUri(core.LiteralUri)
	if re != nil {
		var ok bool
		literal, ok = re.(core.Literal)
		if !ok {
			log.Printf("Recovered object not of type Literal with url http://activeCrl.com/core/Literal")
			return nil
		}
	}
	if literal == nil {
		literal = uOfD.NewLiteral()
		literal.SetOwningElement(coreElement)
		literal.SetUri(core.LiteralUri)
	}

	// LiteralPointer
	var literalPointer core.LiteralPointer
	re = uOfD.GetBaseElementWithUri(core.LiteralPointerUri)
	if re != nil {
		var ok bool
		literalPointer, ok = re.(core.LiteralPointer)
		if !ok {
			log.Printf("Recovered object not of type LiteralPointer with url http://activeCrl.com/core/LiteralPointer")
			return nil
		}
	}
	if literalPointer == nil {
		literalPointer = uOfD.NewValueLiteralPointer()
		literalPointer.SetOwningElement(coreElement)
		literalPointer.SetUri(core.LiteralPointerUri)
	}

	// LiteralPointerPointer
	var literalPointerPointer core.LiteralPointerPointer
	re = uOfD.GetBaseElementWithUri(core.LiteralPointerPointerUri)
	if re != nil {
		var ok bool
		literalPointerPointer, ok = re.(core.LiteralPointerPointer)
		if !ok {
			log.Printf("Recovered object not of type LiteralPointerPointer with url http://activeCrl.com/core/LiteralPointerPointer")
			return nil
		}
	}
	if literalPointerPointer == nil {
		literalPointerPointer = uOfD.NewLiteralPointerPointer()
		literalPointerPointer.SetOwningElement(coreElement)
		literalPointerPointer.SetUri(core.LiteralPointerPointerUri)
	}

	// LiteralPointerReference
	var literalPointerReference core.LiteralPointerReference
	re = uOfD.GetBaseElementWithUri(core.LiteralPointerReferenceUri)
	if re != nil {
		var ok bool
		literalPointerReference, ok = re.(core.LiteralPointerReference)
		if !ok {
			log.Printf("Recovered object not of type LiteralPointerReference with url http://activeCrl.com/core/LiteralPointerReference")
			return nil
		}
	}
	if literalPointerReference == nil {
		literalPointerReference = uOfD.NewLiteralPointerReference()
		literalPointerReference.SetOwningElement(coreElement)
		literalPointerReference.SetName("LiteralReference")
		literalPointerReference.SetUri(core.LiteralPointerReferenceUri)
	}

	// LiteralReference
	var literalReference core.LiteralReference
	re = uOfD.GetBaseElementWithUri(core.LiteralReferenceUri)
	if re != nil {
		var ok bool
		literalReference, ok = re.(core.LiteralReference)
		if !ok {
			log.Printf("Recovered object not of type LiteralReference with url http://activeCrl.com/core/LiteralReference")
			return nil
		}
	}
	if literalReference == nil {
		literalReference = uOfD.NewLiteralReference()
		literalReference.SetOwningElement(coreElement)
		literalReference.SetName("LiteralReference")
		literalReference.SetUri(core.LiteralReferenceUri)
	}

	// Refinement
	var refinement core.Refinement
	re = uOfD.GetBaseElementWithUri(core.RefinementUri)
	if re != nil {
		var ok bool
		refinement, ok = re.(core.Refinement)
		if !ok {
			log.Printf("Recovered object not of type Refinement with url http://activeCrl.com/core/Refinement")
			return nil
		}
	}
	if refinement == nil {
		refinement = uOfD.NewRefinement()
		refinement.SetOwningElement(coreElement)
		refinement.SetName("Refinement")
		refinement.SetUri(core.RefinementUri)
	}
	return coreElement
}
