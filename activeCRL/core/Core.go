package core

import (
	"sync"
)

var CoreUri string = "http://activeCrl.com/core/Core"
var ElememtUri string = "http://activeCrl.com/core/Element"
var ElementPointerUri string = "http://activeCrl.com/core/ElementPointer"
var ElementPointerPointerUri string = "http://activeCrl.com/core/ElementPointerPointer"
var ElementPointerReferenceUri string = "http://activeCrl.com/core/ElementPointerReference"
var ElementReferenceUri string = "http://activeCrl.com/core/ElementReference"
var LiteralUri string = "http://activeCrl.com/core/Literal"
var LiteralPointerUri string = "http://activeCrl.com/core/LiteralPointer"
var LiteralPointerPointerUri string = "http://activeCrl.com/core/LiteralPointerPointer"
var LiteralPointerReferenceUri string = "http://activeCrl.com/core/LiteralPointerReference"
var LiteralReferenceUri string = "http://activeCrl.com/core/LiteralReference"
var RefinementUri string = "http://activeCrl.com/core/Refinement"

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
}

func (c *core) findFunctions(element Element) []crlExecutionFunction {
	var functions []crlExecutionFunction
	if element == nil {
		return functions
	}
	abstractions := element.getAbstractElementsRecursively()
	for _, abstractElement := range abstractions {
		uri := abstractElement.getUri()
		if uri != "" {
			f := c.computeFunctions[uri]
			if f != nil {
				functions = append(functions, f)
			}
		}
	}
	return functions
}
