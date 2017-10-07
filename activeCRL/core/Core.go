// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"sync"
	//	"time"
)

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
	abstractions := element.GetAbstractElementsRecursively(hl)
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

func (c *coreConceptSpace) PrintFunctions() {
	for k, v := range c.computeFunctions {
		log.Printf("Key: %s Value: %p\n", k, v)
	}
}
