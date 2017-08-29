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
	//	core.Print(recoveredCoreFunctions, "---", hl)
	updatedCoreFunctions := updateRecoveredCoreFunctions(recoveredCoreFunctions, uOfD, hl)
	//	core.Print(updatedCoreFunctions, "+++", hl)
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

	coreFunctions.UpdateRecoveredCoreBaseElementFunctions(coreFunctionsElement, uOfD, hl)
	coreFunctions.UpdateRecoveredCoreBaseElementPointerFunctions(coreFunctionsElement, uOfD, hl)
	coreFunctions.UpdateRecoveredCoreBaseElementReferenceFunctions(coreFunctionsElement, uOfD, hl)
	coreFunctions.UpdateRecoveredCoreElementFunctions(coreFunctionsElement, uOfD, hl)

	return coreFunctionsElement
}
