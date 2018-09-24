// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
	//	"time"
)

var CorePrefix string = "http://activeCrl.com/core/"

var UniverseOfDiscourseUri string = CorePrefix + "UniverseOfDiscourse"

var CoreConceptSpaceUri string = "http://activeCrl.com/core/CoreConceptSpace"
var BaseElementPointerUri string = "http://activeCrl.com/core/BaseElementPointer"
var BaseElementReferenceUri string = "http://activeCrl.com/core/BaseElementReference"
var ElememtUri string = "http://activeCrl.com/core/Element"
var ElementPointerUri string = "http://activeCrl.com/core/ElementPointer"
var ElementPointerPointerUri string = "http://activeCrl.com/core/ElementPointerPointer"
var ElementPointerReferenceUri string = "http://activeCrl.com/core/ElementPointerReference"
var ElementPointerRoleUri string = "http://activeCrl.com/core/ElementPointerRole"
var AbstractElementUri string = "http://activeCrl.com/core/ElementPointerRole/AbstractElement"
var RefinedElementUri string = "http://activeCrl.com/core/ElementPointerRole/RefinedElement"
var OwningElementUri string = "http://activeCrl.com/core/ElementPointerRole/OwningElement"
var ReferencedElementUri string = "http://activeCrl.com/core/ElementPointerRole/ReferencedElement"
var ElementReferenceUri string = "http://activeCrl.com/core/ElementReference"
var LiteralUri string = "http://activeCrl.com/core/Literal"
var LiteralPointerUri string = "http://activeCrl.com/core/LiteralPointer"
var LiteralPointerPointerUri string = "http://activeCrl.com/core/LiteralPointerPointer"
var LiteralPointerReferenceUri string = "http://activeCrl.com/core/LiteralPointerReference"
var LiteralPointerRoleUri string = "http://activeCrl.com/core/LiteralPointerRole"
var LabelUri string = "http://activeCrl.com/core/LiteralPointerRole/Label"
var DefinitionUri string = "http://activeCrl.com/core/LiteralPointerRole/Definition"
var UriUri string = "http://activeCrl.com/core/LiteralPointerRole/Uri"
var ValueUri string = "http://activeCrl.com/core/LiteralPointerRole/Value"
var LiteralReferenceUri string = "http://activeCrl.com/core/LiteralReference"
var RefinementUri string = "http://activeCrl.com/core/Refinement"

var AdHocTrace bool = false

// coreSingleton is the singleton instance of the coreConceptSpace. It provides the singular instances of all of the core
// concepts and the mapping from Element identifiers to functions.
var coreSingleton *coreConceptSpace

// GetCore() is a static function that returns the singleton instance of the coreConceptSpace
func GetCore() *coreConceptSpace {
	if coreSingleton == nil {
		coreSingleton = newCore()
	}
	return coreSingleton
}

type coreConceptSpace struct {
	sync.RWMutex
	computeFunctions functions
}

func newCore() *coreConceptSpace {
	var newCore coreConceptSpace
	newCore.computeFunctions = make(map[crlExecutionFunctionArrayIdentifier][]crlExecutionFunction)
	return &newCore
}

func init() {
	coreSingleton = newCore()
	TraceChange = false
	notificationsLimit = 0
	notificationsCount = 0
}

func (c *coreConceptSpace) AddFunction(uri string, function crlExecutionFunction) {
	c.computeFunctions[crlExecutionFunctionArrayIdentifier(uri)] = append(c.computeFunctions[crlExecutionFunctionArrayIdentifier(uri)], function)
}

func (c *coreConceptSpace) GetFunctions(uri string) []crlExecutionFunction {
	return c.computeFunctions[crlExecutionFunctionArrayIdentifier(uri)]
}

func (c *coreConceptSpace) FindFunctions(element Element, notification *ChangeNotification, hl *HeldLocks) []crlExecutionFunctionArrayIdentifier {
	var functionIdentifiers []crlExecutionFunctionArrayIdentifier
	if element == nil {
		return functionIdentifiers
	}
	uOfD := element.GetUniverseOfDiscourse(hl)
	if uOfD == nil {
		return functionIdentifiers
	}
	abstractions := uOfD.GetAbstractElementsRecursively(element, hl)
	for _, abstractElement := range abstractions {
		uri := GetUri(abstractElement, hl)
		if uri != "" {
			f := c.computeFunctions[crlExecutionFunctionArrayIdentifier(uri)]
			if f != nil {
				//				var entry labeledFunction
				//				entry.function = f
				//				entry.label = uri
				functionIdentifiers = append(functionIdentifiers, crlExecutionFunctionArrayIdentifier(uri))
			}
		}
	}
	return functionIdentifiers
}

func buildCoreConceptSpace(uOfD UniverseOfDiscourse, hl *HeldLocks) Element {
	//	log.Printf("*** In buildCoreConceptSpace, held locks present is %v \n", hl != nil)
	// Core
	AdHocTrace = false
	//	log.Printf("*** In buildCoreConceptSpace, about to create Core Element \n")
	coreElement := uOfD.NewElement(hl, CoreConceptSpaceUri)
	//	log.Printf("*** In buildCoreConceptSpace, about to call SetLabel on Core Element \n")
	SetLabel(coreElement, "CoreConceptSpace", hl)
	//	log.Printf("*** In buildCoreConceptSpace, about to call SetURI on Core Element \n")
	SetUri(coreElement, CoreConceptSpaceUri, hl)
	//	log.Printf("*** In buildCoreConceptSpace, completed calling SetUri on core Element \n")
	AdHocTrace = false

	// BaseElementPointer
	baseElementPointer := uOfD.NewBaseElementPointer(hl, BaseElementPointerUri)
	SetOwningElement(baseElementPointer, coreElement, hl)
	SetUri(baseElementPointer, BaseElementPointerUri, hl)

	// BaseElementReference
	baseElementReference := uOfD.NewBaseElementReference(hl, BaseElementReferenceUri)
	SetOwningElement(baseElementReference, coreElement, hl)
	SetLabel(baseElementReference, "BaseElementReference", hl)
	SetUri(baseElementReference, BaseElementReferenceUri, hl)

	// Element
	element := uOfD.NewElement(hl, ElememtUri)
	SetOwningElement(element, coreElement, hl)
	SetLabel(element, "Element", hl)
	SetUri(element, ElememtUri, hl)

	// ElementPointer
	elementPointer := uOfD.NewReferencedElementPointer(hl, ElementPointerUri)
	SetOwningElement(elementPointer, coreElement, hl)
	SetUri(elementPointer, ElementPointerUri, hl)

	// ElementPointerPointer
	elementPointerPointer := uOfD.NewElementPointerPointer(hl, ElementPointerPointerUri)
	SetOwningElement(elementPointerPointer, coreElement, hl)
	SetUri(elementPointerPointer, ElementPointerPointerUri, hl)

	// ElementPointerReference
	elementPointerReference := uOfD.NewElementPointerReference(hl, ElementPointerReferenceUri)
	SetOwningElement(elementPointerReference, coreElement, hl)
	SetLabel(elementPointerReference, "ElementReference", hl)
	SetUri(elementPointerReference, ElementPointerReferenceUri, hl)

	// ElementPointerRole and values
	elementPointerRole := uOfD.NewElement(hl, ElementPointerRoleUri)
	SetOwningElement(elementPointerRole, coreElement, hl)
	SetLabel(elementPointerRole, "ElementPointerRole", hl)
	SetUri(elementPointerRole, ElementPointerRoleUri, hl)

	abstractElement := uOfD.NewElement(hl, AbstractElementUri)
	SetOwningElement(abstractElement, coreElement, hl)
	SetLabel(abstractElement, "AbstractElementRole", hl)
	SetUri(abstractElement, AbstractElementUri, hl)
	refinement0 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement0, abstractElement, hl)
	refinement0.SetAbstractElement(elementPointerRole, hl)
	refinement0.SetRefinedElement(abstractElement, hl)

	refinedElement := uOfD.NewElement(hl, RefinedElementUri)
	SetOwningElement(refinedElement, coreElement, hl)
	SetLabel(refinedElement, "RefinedElementRole", hl)
	SetUri(refinedElement, RefinedElementUri, hl)
	refinement1 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement1, refinedElement, hl)
	refinement1.SetAbstractElement(elementPointerRole, hl)
	refinement1.SetRefinedElement(refinedElement, hl)

	owningElement := uOfD.NewElement(hl, OwningElementUri)
	SetOwningElement(owningElement, coreElement, hl)
	SetLabel(owningElement, "OwningElementRole", hl)
	SetUri(owningElement, OwningElementUri, hl)
	refinement2 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement2, owningElement, hl)
	refinement2.SetAbstractElement(elementPointerRole, hl)
	refinement2.SetRefinedElement(owningElement, hl)

	referencedElement := uOfD.NewElement(hl, ReferencedElementUri)
	SetOwningElement(referencedElement, coreElement, hl)
	SetLabel(referencedElement, "ReferencedElementRole", hl)
	SetUri(referencedElement, ReferencedElementUri, hl)
	refinement3 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement3, referencedElement, hl)
	refinement3.SetAbstractElement(elementPointerRole, hl)
	refinement3.SetRefinedElement(referencedElement, hl)

	// ElementReference
	elementReference := uOfD.NewElementReference(hl, ElementReferenceUri)
	SetOwningElement(elementReference, coreElement, hl)
	SetLabel(elementReference, "ElementReference", hl)
	SetUri(elementReference, ElementReferenceUri, hl)

	// Literal
	literal := uOfD.NewLiteral(hl, LiteralUri)
	SetOwningElement(literal, coreElement, hl)
	SetUri(literal, LiteralUri, hl)

	// LiteralPointer
	literalPointer := uOfD.NewValueLiteralPointer(hl, LiteralPointerUri)
	SetOwningElement(literalPointer, coreElement, hl)
	SetUri(literalPointer, LiteralPointerUri, hl)

	// LiteralPointerPointer
	literalPointerPointer := uOfD.NewLiteralPointerPointer(hl, LiteralPointerPointerUri)
	SetOwningElement(literalPointerPointer, coreElement, hl)
	SetUri(literalPointerPointer, LiteralPointerPointerUri, hl)

	// LiteralPointerReference
	literalPointerReference := uOfD.NewLiteralPointerReference(hl, LiteralPointerReferenceUri)
	SetOwningElement(literalPointerReference, coreElement, hl)
	SetLabel(literalPointerReference, "LiteralReference", hl)
	SetUri(literalPointerReference, LiteralPointerReferenceUri, hl)

	// LiteralPointerRole and values
	literalPointerRole := uOfD.NewElement(hl, LiteralPointerRoleUri)
	SetOwningElement(literalPointerRole, coreElement, hl)
	SetLabel(literalPointerRole, "LiteralPointerRole", hl)
	SetUri(literalPointerRole, LiteralPointerRoleUri, hl)

	name := uOfD.NewElement(hl, LabelUri)
	SetOwningElement(name, literalPointerRole, hl)
	SetLabel(name, "Label", hl)
	SetUri(name, LabelUri, hl)
	refinement4 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement4, name, hl)
	refinement4.SetAbstractElement(literalPointerRole, hl)
	refinement4.SetRefinedElement(name, hl)

	definition := uOfD.NewElement(hl, DefinitionUri)
	SetOwningElement(definition, literalPointerRole, hl)
	SetLabel(definition, "Definition", hl)
	SetUri(definition, DefinitionUri, hl)
	refinement5 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement5, definition, hl)
	refinement5.SetAbstractElement(literalPointerRole, hl)
	refinement5.SetRefinedElement(definition, hl)

	uri := uOfD.NewElement(hl, UriUri)
	SetOwningElement(uri, literalPointerRole, hl)
	SetLabel(uri, "Uri", hl)
	SetUri(uri, UriUri, hl)
	refinement6 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement6, uri, hl)
	refinement6.SetAbstractElement(literalPointerRole, hl)
	refinement6.SetRefinedElement(uri, hl)

	value := uOfD.NewElement(hl, ValueUri)
	SetOwningElement(value, literalPointerRole, hl)
	SetLabel(value, "Value", hl)
	SetUri(value, ValueUri, hl)
	refinement7 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement7, value, hl)
	refinement7.SetAbstractElement(literalPointerRole, hl)
	refinement7.SetRefinedElement(value, hl)

	// LiteralReference
	literalReference := uOfD.NewLiteralReference(hl, LiteralReferenceUri)
	SetOwningElement(literalReference, coreElement, hl)
	SetLabel(literalReference, "LiteralReference", hl)
	SetUri(literalReference, LiteralReferenceUri, hl)

	// Refinement
	refinement := uOfD.NewRefinement(hl, RefinementUri)
	SetOwningElement(refinement, coreElement, hl)
	SetLabel(refinement, "Refinement", hl)
	SetUri(refinement, RefinementUri, hl)

	return coreElement
}

func (c *coreConceptSpace) PrintFunctions() {
	for k, v := range c.computeFunctions {
		log.Printf("Key: %s Value: %p\n", k, v)
	}
}
