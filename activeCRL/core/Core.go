package core

import (
	"log"
	"sync"
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

type crlExecutionFunction func(Element, *ChangeNotification)

type functions map[string]crlExecutionFunction

var coreSingleton *core

func GetCore() *core {
	if coreSingleton == nil {
		coreSingleton = newCore()
	}
	return coreSingleton
}

type core struct {
	sync.RWMutex
	computeFunctions functions
}

func newCore() *core {
	var newCore core
	newCore.computeFunctions = make(map[string]crlExecutionFunction)
	return &newCore
}

func init() {
	coreSingleton = newCore()
	TraceChange = false
}

func (c *core) AddFunction(uri string, function crlExecutionFunction) {
	c.computeFunctions[uri] = function
}

func (c *core) GetFunction(uri string) crlExecutionFunction {
	return c.computeFunctions[uri]
}

func (c *core) FindFunctions(element Element, hl *HeldLocks) []crlExecutionFunction {
	if AdHocTrace == true {
		Print(element, "Finding functions for: ", hl)
	}
	var functions []crlExecutionFunction
	if element == nil {
		return functions
	}
	abstractions := element.GetAbstractElementsRecursively(hl)
	for _, abstractElement := range abstractions {
		uri := GetUri(abstractElement, hl)
		if AdHocTrace == true {
			log.Printf("AbstractElement URI: %s\n", uri)
		}
		if uri != "" {
			f := c.computeFunctions[uri]
			if f != nil {
				functions = append(functions, f)
			}
		}
	}
	return functions
}

func (c *core) PrintFunctions() {
	for k, v := range c.computeFunctions {
		log.Printf("Key: %s Value: %p\n", k, v)
	}
}
