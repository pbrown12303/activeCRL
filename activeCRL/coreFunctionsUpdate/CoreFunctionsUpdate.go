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
	recoveredCoreFunctions := coreFunctions.GetCoreFunctionsConceptSpace(uOfD)
	core.Print(recoveredCoreFunctions, "---")
	updatedCoreFunctions := updateRecoveredCoreFunctions(recoveredCoreFunctions, uOfD)
	core.Print(updatedCoreFunctions, "+++")
	marshaledCoreFunctions, err := updatedCoreFunctions.MarshalJSON()
	if err == nil {
		ioutil.WriteFile("CoreFunctions.acrl", marshaledCoreFunctions, os.ModePerm)
		serializedCoreFunctions := serializedCoreFunctionsPrefix + "`" + string(marshaledCoreFunctions) + "`"
		ioutil.WriteFile("./coreFunctions/SerializedCoreFunctions.go", []byte(serializedCoreFunctions), os.ModePerm)
	}
}

func updateRecoveredCoreFunctions(recoveredCoreFunctions core.Element, uOfD *core.UniverseOfDiscourse) core.Element {
	// Core
	coreFunctionsElement := recoveredCoreFunctions
	if coreFunctionsElement == nil {
		coreFunctionsElement = uOfD.NewElement()
		coreFunctionsElement.SetName("CoreConceptSpace")
		coreFunctionsElement.SetUri(core.CoreConceptSpaceUri)
	}

	// CreateElement
	createElement := uOfD.GetElementWithUri(coreFunctions.CreateElememtUri)
	if createElement == nil {
		createElement = uOfD.NewElement()
		createElement.SetOwningElement(coreFunctionsElement)
		createElement.SetName("CreateElement")
		createElement.SetUri(coreFunctions.CreateElememtUri)
	}
	// CreatedElementReference
	createdElementReference := core.GetChildElementReferenceWithUri(createElement, coreFunctions.CreatedElementReferenceUri)
	if createdElementReference == nil {
		createdElementReference = uOfD.NewElementReference()
		createdElementReference.SetOwningElement(createElement)
		createdElementReference.SetName("CreatedElementReference")
		createdElementReference.SetUri(coreFunctions.CreatedElementReferenceUri)
	}
	return coreFunctionsElement
}
