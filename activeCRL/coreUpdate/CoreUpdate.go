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
	hl := core.NewHeldLocks(nil)
	defer hl.ReleaseLocks()
	recoveredCore := uOfD.GetCoreConceptSpace()
	updatedCore := updateRecoveredCore(recoveredCore, uOfD, hl)
	marshaledCore, err := updatedCore.MarshalJSON()
	if err == nil {
		ioutil.WriteFile("CoreConceptSpace.acrl", marshaledCore, os.ModePerm)
		serializedCore := serializedCorePrefix + "`" + string(marshaledCore) + "`"
		ioutil.WriteFile("./core/SerializedCore.go", []byte(serializedCore), os.ModePerm)
	}
}

func updateRecoveredCore(recoveredCore core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// Core
	coreElement := recoveredCore
	if coreElement == nil {
		coreElement = uOfD.NewElement(hl)
		core.SetName(coreElement, "CoreConceptSpace", hl)
		core.SetUri(coreElement, core.CoreConceptSpaceUri, hl)
	}

	// BaseElementPointer
	var baseElementPointer core.BaseElementPointer
	re := uOfD.GetBaseElementWithUri(core.BaseElementPointerUri)
	if re != nil {
		var ok bool
		baseElementPointer, ok = re.(core.BaseElementPointer)
		if !ok {
			log.Printf("Recovered object not of type BaseElementPointer with url %s\n", core.BaseElementPointerUri)
			core.Print(re, "", hl)
			return nil
		}
	}
	if baseElementPointer == nil {
		baseElementPointer = uOfD.NewBaseElementPointer(hl)
		core.SetOwningElement(baseElementPointer, coreElement, hl)
		core.SetUri(baseElementPointer, core.BaseElementPointerUri, hl)
	}

	// BaseElementReference
	var baseElementReference core.BaseElementReference
	re = uOfD.GetBaseElementWithUri(core.BaseElementReferenceUri)
	if re != nil {
		var ok bool
		baseElementReference, ok = re.(core.BaseElementReference)
		if !ok {
			log.Printf("Recovered object not of type ElementReference with url %s\n", core.BaseElementReferenceUri)
			return nil
		}
	}
	if baseElementReference == nil {
		baseElementReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(baseElementReference, coreElement, hl)
		core.SetName(baseElementReference, "BaseElementReference", hl)
		core.SetUri(baseElementReference, core.BaseElementReferenceUri, hl)
	}

	// Element
	re = uOfD.GetBaseElementWithUri(core.ElememtUri)
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
		element = uOfD.NewElement(hl)
		core.SetOwningElement(element, coreElement, hl)
		core.SetName(element, "Element", hl)
		core.SetUri(element, core.ElememtUri, hl)
	}

	// ElementPointer
	var elementPointer core.ElementPointer
	re = uOfD.GetBaseElementWithUri(core.ElementPointerUri)
	if re != nil {
		var ok bool
		elementPointer, ok = re.(core.ElementPointer)
		if !ok {
			log.Printf("Recovered object not of type ElementPointer with url http://activeCrl.com/core/ElementPointer")
			core.Print(re, "", hl)
			uOfD.DeleteBaseElement(re, hl)
			re = nil
		}
	}
	if elementPointer == nil {
		elementPointer = uOfD.NewReferencedElementPointer(hl)
		core.SetOwningElement(elementPointer, coreElement, hl)
		core.SetUri(elementPointer, core.ElementPointerUri, hl)
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
		elementPointerPointer = uOfD.NewElementPointerPointer(hl)
		core.SetOwningElement(elementPointerPointer, coreElement, hl)
		core.SetUri(elementPointerPointer, core.ElementPointerPointerUri, hl)
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
		elementPointerReference = uOfD.NewElementPointerReference(hl)
		core.SetOwningElement(elementPointerReference, coreElement, hl)
		core.SetName(elementPointerReference, "ElementReference", hl)
		core.SetUri(elementPointerReference, core.ElementPointerReferenceUri, hl)
	}

	// ElementPointerRole and values
	var elementPointerRole core.Element
	re = uOfD.GetBaseElementWithUri(core.ElementPointerRoleUri)
	if re != nil {
		var ok bool
		elementPointerRole, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url http://activeCrl.com/core/ElementPointerRole")
			return nil
		}
	}
	if elementPointerRole == nil {
		elementPointerRole = uOfD.NewElement(hl)
		core.SetOwningElement(elementPointerRole, coreElement, hl)
		core.SetName(elementPointerRole, "ElementPointerRole", hl)
		core.SetUri(elementPointerRole, core.ElementPointerRoleUri, hl)
	}

	var abstractElement core.Element
	re = uOfD.GetBaseElementWithUri(core.ReferencedElementUri)
	if re != nil {
		var ok bool
		abstractElement, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.AbstractElementUri)
			return nil
		}
	}
	if abstractElement == nil {
		abstractElement = uOfD.NewElement(hl)
		core.SetOwningElement(abstractElement, coreElement, hl)
		core.SetName(abstractElement, "AbstractElementRole", hl)
		core.SetUri(abstractElement, core.AbstractElementUri, hl)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, abstractElement, hl)
		refinement.SetAbstractElement(elementPointerRole, hl)
		refinement.SetRefinedElement(abstractElement, hl)
	}

	var refinedElement core.Element
	re = uOfD.GetBaseElementWithUri(core.RefinedElementUri)
	if re != nil {
		var ok bool
		refinedElement, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.RefinedElementUri)
			return nil
		}
	}
	if refinedElement == nil {
		refinedElement = uOfD.NewElement(hl)
		core.SetOwningElement(refinedElement, coreElement, hl)
		core.SetName(refinedElement, "RefinedElementRole", hl)
		core.SetUri(refinedElement, core.RefinedElementUri, hl)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, refinedElement, hl)
		refinement.SetAbstractElement(elementPointerRole, hl)
		refinement.SetRefinedElement(refinedElement, hl)
	}

	var owningElement core.Element
	re = uOfD.GetBaseElementWithUri(core.OwningElementUri)
	if re != nil {
		var ok bool
		owningElement, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.OwningElementUri)
			return nil
		}
	}
	if owningElement == nil {
		owningElement = uOfD.NewElement(hl)
		core.SetOwningElement(owningElement, coreElement, hl)
		core.SetName(owningElement, "OwningElementRole", hl)
		core.SetUri(owningElement, core.OwningElementUri, hl)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, owningElement, hl)
		refinement.SetAbstractElement(elementPointerRole, hl)
		refinement.SetRefinedElement(owningElement, hl)
	}

	var referencedElement core.Element
	re = uOfD.GetBaseElementWithUri(core.ReferencedElementUri)
	if re != nil {
		var ok bool
		referencedElement, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.RefinedElementUri)
			return nil
		}
	}
	if referencedElement == nil {
		referencedElement = uOfD.NewElement(hl)
		core.SetOwningElement(referencedElement, coreElement, hl)
		core.SetName(referencedElement, "ReferencedElementRole", hl)
		core.SetUri(referencedElement, core.ReferencedElementUri, hl)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, referencedElement, hl)
		refinement.SetAbstractElement(elementPointerRole, hl)
		refinement.SetRefinedElement(referencedElement, hl)
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
		elementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(elementReference, coreElement, hl)
		core.SetName(elementReference, "ElementReference", hl)
		core.SetUri(elementReference, core.ElementReferenceUri, hl)
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
		literal = uOfD.NewLiteral(hl)
		core.SetOwningElement(literal, coreElement, hl)
		core.SetUri(literal, core.LiteralUri, hl)
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
		literalPointer = uOfD.NewValueLiteralPointer(hl)
		core.SetOwningElement(literalPointer, coreElement, hl)
		core.SetUri(literalPointer, core.LiteralPointerUri, hl)
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
		literalPointerPointer = uOfD.NewLiteralPointerPointer(hl)
		core.SetOwningElement(literalPointerPointer, coreElement, hl)
		core.SetUri(literalPointerPointer, core.LiteralPointerPointerUri, hl)
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
		literalPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(literalPointerReference, coreElement, hl)
		core.SetName(literalPointerReference, "LiteralReference", hl)
		core.SetUri(literalPointerReference, core.LiteralPointerReferenceUri, hl)
	}

	// LiteralPointerRole and values
	var literalPointerRole core.Element
	re = uOfD.GetBaseElementWithUri(core.LiteralPointerRoleUri)
	if re != nil {
		var ok bool
		literalPointerRole, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.LiteralPointerRoleUri)
			return nil
		}
	}
	if literalPointerRole == nil {
		literalPointerRole = uOfD.NewElement(hl)
		core.SetOwningElement(literalPointerRole, coreElement, hl)
		core.SetName(literalPointerRole, "LiteralPointerRole", hl)
		core.SetUri(literalPointerRole, core.LiteralPointerRoleUri, hl)
	}

	var name core.Element
	re = uOfD.GetBaseElementWithUri(core.NameUri)
	if re != nil {
		var ok bool
		name, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.NameUri)
			return nil
		}
	}
	if name == nil {
		name = uOfD.NewElement(hl)
		core.SetOwningElement(name, literalPointerRole, hl)
		core.SetName(name, "Name", hl)
		core.SetUri(name, core.NameUri, hl)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, name, hl)
		refinement.SetAbstractElement(literalPointerRole, hl)
		refinement.SetRefinedElement(name, hl)
	}

	var definition core.Element
	re = uOfD.GetBaseElementWithUri(core.DefinitionUri)
	if re != nil {
		var ok bool
		definition, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.DefinitionUri)
			return nil
		}
	}
	if definition == nil {
		definition = uOfD.NewElement(hl)
		core.SetOwningElement(definition, literalPointerRole, hl)
		core.SetName(definition, "Definition", hl)
		core.SetUri(definition, core.DefinitionUri, hl)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, definition, hl)
		refinement.SetAbstractElement(literalPointerRole, hl)
		refinement.SetRefinedElement(definition, hl)
	}

	var uri core.Element
	re = uOfD.GetBaseElementWithUri(core.UriUri)
	if re != nil {
		var ok bool
		uri, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.UriUri)
			return nil
		}
	}
	if uri == nil {
		uri = uOfD.NewElement(hl)
		core.SetOwningElement(uri, literalPointerRole, hl)
		core.SetName(uri, "Uri", hl)
		core.SetUri(uri, core.UriUri, hl)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, uri, hl)
		refinement.SetAbstractElement(literalPointerRole, hl)
		refinement.SetRefinedElement(uri, hl)
	}

	var value core.Element
	re = uOfD.GetBaseElementWithUri(core.ValueUri)
	if re != nil {
		var ok bool
		value, ok = re.(core.Element)
		if !ok {
			log.Printf("Recovered object not of type Element with url %s\n", core.ValueUri)
			return nil
		}
	}
	if value == nil {
		value = uOfD.NewElement(hl)
		core.SetOwningElement(value, literalPointerRole, hl)
		core.SetName(value, "Value", hl)
		core.SetUri(value, core.ValueUri, hl)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, value, hl)
		refinement.SetAbstractElement(literalPointerRole, hl)
		refinement.SetRefinedElement(value, hl)
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
		literalReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(literalReference, coreElement, hl)
		core.SetName(literalReference, "LiteralReference", hl)
		core.SetUri(literalReference, core.LiteralReferenceUri, hl)
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
		refinement = uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, coreElement, hl)
		core.SetName(refinement, "Refinement", hl)
		core.SetUri(refinement, core.RefinementUri, hl)
	}
	return coreElement
}
