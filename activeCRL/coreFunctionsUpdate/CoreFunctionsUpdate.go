package main

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/pbrown12303/activeCRL/activeCRL/coreFunctions"
	"io/ioutil"
	"os"
)

var serializedCoreFunctionsPrefix string = `package coreFunctions

import ()

var serializedCoreFunctions string = `

func main() {
	uOfD := core.NewUniverseOfDiscourse()
	hl := core.NewHeldLocks()
	defer hl.ReleaseLocks()
	recoveredCoreFunctions := coreFunctions.GetCoreFunctionsConceptSpace(uOfD)
	core.Print(recoveredCoreFunctions, "---", hl)
	updatedCoreFunctions := updateRecoveredCoreFunctions(recoveredCoreFunctions, uOfD, hl)
	core.Print(updatedCoreFunctions, "+++", hl)
	marshaledCoreFunctions, err := updatedCoreFunctions.MarshalJSON()
	if err == nil {
		ioutil.WriteFile("CoreFunctions.acrl", marshaledCoreFunctions, os.ModePerm)
		serializedCoreFunctions := serializedCoreFunctionsPrefix + "`" + string(marshaledCoreFunctions) + "`"
		ioutil.WriteFile("./coreFunctions/SerializedCoreFunctions.go", []byte(serializedCoreFunctions), os.ModePerm)
	}
}

func updateRecoveredCoreFunctions(recoveredCoreFunctions core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) core.Element {
	// Core
	coreFunctionsElement := recoveredCoreFunctions
	if coreFunctionsElement == nil {
		coreFunctionsElement = uOfD.NewElement(hl)
		core.SetName(coreFunctionsElement, "CoreFunctions", hl)
		core.SetUri(coreFunctionsElement, coreFunctions.CoreFunctionsUri, hl)
	}

	// CreateElement
	createElement := uOfD.GetElementWithUri(coreFunctions.CreateElememtUri)
	if createElement == nil {
		createElement = uOfD.NewElement(hl)
		core.SetOwningElement(createElement, coreFunctionsElement, hl)
		core.SetName(createElement, "CreateElement", hl)
		core.SetUri(createElement, coreFunctions.CreateElememtUri, hl)
	}
	// CreatedElementReference
	createdElementReference := core.GetChildElementReferenceWithUri(createElement, coreFunctions.CreatedElementReferenceUri, hl)
	if createdElementReference == nil {
		createdElementReference = uOfD.NewElementReference(hl)
		core.SetOwningElement(createdElementReference, createElement, hl)
		core.SetName(createdElementReference, "CreatedElementReference", hl)
		core.SetUri(createdElementReference, coreFunctions.CreatedElementReferenceUri, hl)
	}
	return coreFunctionsElement
}
