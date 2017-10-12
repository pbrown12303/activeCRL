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
var NameUri string = "http://activeCrl.com/core/LiteralPointerRole/Name"
var DefinitionUri string = "http://activeCrl.com/core/LiteralPointerRole/Definition"
var UriUri string = "http://activeCrl.com/core/LiteralPointerRole/Uri"
var ValueUri string = "http://activeCrl.com/core/LiteralPointerRole/Value"
var LiteralReferenceUri string = "http://activeCrl.com/core/LiteralReference"
var RefinementUri string = "http://activeCrl.com/core/Refinement"

var AdHocTrace bool = false

type crlExecutionFunction func(Element, []*ChangeNotification, *sync.WaitGroup)

type functions map[string]crlExecutionFunction

type labeledFunction struct {
	function crlExecutionFunction
	label    string
}

var coreSingleton *coreConceptSpace

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
	newCore.computeFunctions = make(map[string]crlExecutionFunction)
	return &newCore
}

func init() {
	coreSingleton = newCore()
	TraceChange = false
	notificationsLimit = 0
	notificationsCount = 0
}

func (c *coreConceptSpace) AddFunction(uri string, function crlExecutionFunction) {
	c.computeFunctions[uri] = function
}

func (c *coreConceptSpace) GetFunction(uri string) crlExecutionFunction {
	return c.computeFunctions[uri]
}

func (c *coreConceptSpace) FindFunctions(element Element, notification *ChangeNotification, hl *HeldLocks) []*labeledFunction {
	var labeledFunctions []*labeledFunction
	if element == nil {
		return labeledFunctions
	}
	uOfD := element.GetUniverseOfDiscourse(hl)
	if uOfD == nil {
		return labeledFunctions
	}
	abstractions := uOfD.GetAbstractElementsRecursively(element, hl)
	for _, abstractElement := range abstractions {
		uri := GetUri(abstractElement, hl)
		if uri != "" {
			f := c.computeFunctions[uri]
			if f != nil {
				var entry labeledFunction
				entry.function = f
				entry.label = uri
				labeledFunctions = append(labeledFunctions, &entry)
			}
		}
	}
	return labeledFunctions
}

func buildCoreConceptSpace(uOfD UniverseOfDiscourse, hl *HeldLocks) Element {
	// Core
	coreElement := uOfD.NewElement(hl, CoreConceptSpaceUri)
	SetName(coreElement, "CoreConceptSpace", hl)
	SetUri(coreElement, CoreConceptSpaceUri, hl)

	// BaseElementPointer
	baseElementPointer := uOfD.NewBaseElementPointer(hl, BaseElementPointerUri)
	SetOwningElement(baseElementPointer, coreElement, hl)
	SetUri(baseElementPointer, BaseElementPointerUri, hl)

	// BaseElementReference
	baseElementReference := uOfD.NewBaseElementReference(hl, BaseElementReferenceUri)
	SetOwningElement(baseElementReference, coreElement, hl)
	SetName(baseElementReference, "BaseElementReference", hl)
	SetUri(baseElementReference, BaseElementReferenceUri, hl)

	// Element
	element := uOfD.NewElement(hl, ElememtUri)
	SetOwningElement(element, coreElement, hl)
	SetName(element, "Element", hl)
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
	SetName(elementPointerReference, "ElementReference", hl)
	SetUri(elementPointerReference, ElementPointerReferenceUri, hl)

	// ElementPointerRole and values
	elementPointerRole := uOfD.NewElement(hl, ElementPointerRoleUri)
	SetOwningElement(elementPointerRole, coreElement, hl)
	SetName(elementPointerRole, "ElementPointerRole", hl)
	SetUri(elementPointerRole, ElementPointerRoleUri, hl)

	abstractElement := uOfD.NewElement(hl, AbstractElementUri)
	SetOwningElement(abstractElement, coreElement, hl)
	SetName(abstractElement, "AbstractElementRole", hl)
	SetUri(abstractElement, AbstractElementUri, hl)
	refinement0 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement0, abstractElement, hl)
	refinement0.SetAbstractElement(elementPointerRole, hl)
	refinement0.SetRefinedElement(abstractElement, hl)

	refinedElement := uOfD.NewElement(hl, RefinedElementUri)
	SetOwningElement(refinedElement, coreElement, hl)
	SetName(refinedElement, "RefinedElementRole", hl)
	SetUri(refinedElement, RefinedElementUri, hl)
	refinement1 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement1, refinedElement, hl)
	refinement1.SetAbstractElement(elementPointerRole, hl)
	refinement1.SetRefinedElement(refinedElement, hl)

	owningElement := uOfD.NewElement(hl, OwningElementUri)
	SetOwningElement(owningElement, coreElement, hl)
	SetName(owningElement, "OwningElementRole", hl)
	SetUri(owningElement, OwningElementUri, hl)
	refinement2 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement2, owningElement, hl)
	refinement2.SetAbstractElement(elementPointerRole, hl)
	refinement2.SetRefinedElement(owningElement, hl)

	referencedElement := uOfD.NewElement(hl, ReferencedElementUri)
	SetOwningElement(referencedElement, coreElement, hl)
	SetName(referencedElement, "ReferencedElementRole", hl)
	SetUri(referencedElement, ReferencedElementUri, hl)
	refinement3 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement3, referencedElement, hl)
	refinement3.SetAbstractElement(elementPointerRole, hl)
	refinement3.SetRefinedElement(referencedElement, hl)

	// ElementReference
	elementReference := uOfD.NewElementReference(hl, ElementReferenceUri)
	SetOwningElement(elementReference, coreElement, hl)
	SetName(elementReference, "ElementReference", hl)
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
	SetName(literalPointerReference, "LiteralReference", hl)
	SetUri(literalPointerReference, LiteralPointerReferenceUri, hl)

	// LiteralPointerRole and values
	literalPointerRole := uOfD.NewElement(hl, LiteralPointerRoleUri)
	SetOwningElement(literalPointerRole, coreElement, hl)
	SetName(literalPointerRole, "LiteralPointerRole", hl)
	SetUri(literalPointerRole, LiteralPointerRoleUri, hl)

	name := uOfD.NewElement(hl, NameUri)
	SetOwningElement(name, literalPointerRole, hl)
	SetName(name, "Name", hl)
	SetUri(name, NameUri, hl)
	refinement4 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement4, name, hl)
	refinement4.SetAbstractElement(literalPointerRole, hl)
	refinement4.SetRefinedElement(name, hl)

	definition := uOfD.NewElement(hl, DefinitionUri)
	SetOwningElement(definition, literalPointerRole, hl)
	SetName(definition, "Definition", hl)
	SetUri(definition, DefinitionUri, hl)
	refinement5 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement5, definition, hl)
	refinement5.SetAbstractElement(literalPointerRole, hl)
	refinement5.SetRefinedElement(definition, hl)

	uri := uOfD.NewElement(hl, UriUri)
	SetOwningElement(uri, literalPointerRole, hl)
	SetName(uri, "Uri", hl)
	SetUri(uri, UriUri, hl)
	refinement6 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement6, uri, hl)
	refinement6.SetAbstractElement(literalPointerRole, hl)
	refinement6.SetRefinedElement(uri, hl)

	value := uOfD.NewElement(hl, ValueUri)
	SetOwningElement(value, literalPointerRole, hl)
	SetName(value, "Value", hl)
	SetUri(value, ValueUri, hl)
	refinement7 := uOfD.NewRefinement(hl)
	SetOwningElement(refinement7, value, hl)
	refinement7.SetAbstractElement(literalPointerRole, hl)
	refinement7.SetRefinedElement(value, hl)

	// LiteralReference
	literalReference := uOfD.NewLiteralReference(hl, LiteralReferenceUri)
	SetOwningElement(literalReference, coreElement, hl)
	SetName(literalReference, "LiteralReference", hl)
	SetUri(literalReference, LiteralReferenceUri, hl)

	// Refinement
	refinement := uOfD.NewRefinement(hl, RefinementUri)
	SetOwningElement(refinement, coreElement, hl)
	SetName(refinement, "Refinement", hl)
	SetUri(refinement, RefinementUri, hl)

	return coreElement
}

func (c *coreConceptSpace) PrintFunctions() {
	for k, v := range c.computeFunctions {
		log.Printf("Key: %s Value: %p\n", k, v)
	}
}
