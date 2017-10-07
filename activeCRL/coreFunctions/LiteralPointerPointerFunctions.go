package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"strconv"
	"sync"
)

var LiteralPointerPointerFunctionsUri string = CoreFunctionsPrefix + "LiteralPointerPointerFunctions"

var LiteralPointerPointerCreateLiteralPointerPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/CreateAbstractLiteralPointerPointer"
var LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri = CoreFunctionsPrefix + "LiteralPointerPointer/CreateAbstractLiteralPointerPointer/CreatedLiteralPointerPointerRef"

var LiteralPointerPointerGetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer"
var LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer/SourceLiteralPointerPointerRef"
var LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointer/IndicatedLiteralPointerRef"

var LiteralPointerPointerGetLiteralPointerIdUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId"
var LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId/SourceLiteralPointerPointerRef"
var LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerId/CreatedLiteralRef"

var LiteralPointerPointerGetLiteralPointerVersionUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion"
var LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion/SourceLiteralPointerPointerRef"
var LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/GetLiteralPointerVersion/CreatedLiteralRef"

var LiteralPointerPointerSetLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer"
var LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer/LiteralPointerRef"
var LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri string = CoreFunctionsPrefix + "LiteralPointerPointer/SetLiteralPointer/ModifiedLiteralPointerPointerRef"

func createLiteralPointerPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createLiteralPointerPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(element, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri, hl)
	if createdLiteralPointerPointerRef == nil {
		createdLiteralPointerPointerRef = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(createdLiteralPointerPointerRef, element, hl)
		core.SetName(createdLiteralPointerPointerRef, "CreatedLiteralPointerPointerRef", hl)
		rootCreatedElementReference := uOfD.GetBaseElementReferenceWithUri(LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerPointerRef, hl)
		refinement.SetRefinedElement(createdLiteralPointerPointerRef, hl)
		refinement.SetAbstractElement(rootCreatedElementReference, hl)
	}
	createdLiteralPointerPointer := createdLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	if createdLiteralPointerPointer == nil {
		createdLiteralPointerPointer = uOfD.NewLiteralPointerPointer(hl)
		createdLiteralPointerPointerRef.SetReferencedBaseElement(createdLiteralPointerPointer, hl)
	} else {
		switch createdLiteralPointerPointer.(type) {
		case core.LiteralPointerPointer:
		default:
			// It's the wrong type - create the correct type
			createdLiteralPointerPointer = uOfD.NewLiteralPointerPointer(hl)
			createdLiteralPointerPointerRef.SetReferencedBaseElement(createdLiteralPointerPointer, hl)
		}
	}
}

func getLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		log.Printf("In GetLiteralPointer, the SourceLiteralPointerPointerRef was not found in the replicate")
		return
	}

	indicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if indicatedLiteralPointerRef == nil {
		log.Printf("In GetLiteralPointer, the IndicatedLiteralPointerrRef was not found in the replicate")
		return
	}

	indicatedLiteralPointer := indicatedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	untypedSourceLiteralPointerPointer := sourceLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceLiteralPointerPointer core.LiteralPointerPointer
	var sourceLiteralPointer core.LiteralPointer
	if untypedSourceLiteralPointerPointer != nil {
		switch untypedSourceLiteralPointerPointer.(type) {
		case core.LiteralPointerPointer:
			sourceLiteralPointerPointer = untypedSourceLiteralPointerPointer.(core.LiteralPointerPointer)
			sourceLiteralPointer = sourceLiteralPointerPointer.GetLiteralPointer(hl)
		}
	}
	if sourceLiteralPointer != indicatedLiteralPointer {
		indicatedLiteralPointerRef.SetReferencedLiteralPointer(sourceLiteralPointer, hl)
	}
}

func getLiteralPointerId(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		log.Printf("In GetLiteralPointerId, the SourceLiteralPointerPointerRef was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetLiteralPointerId, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedSourceLiteralPointerPointer := sourceLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceLiteralPointerPointer core.LiteralPointerPointer
	if untypedSourceLiteralPointerPointer != nil {
		switch untypedSourceLiteralPointerPointer.(type) {
		case core.LiteralPointerPointer:
			sourceLiteralPointerPointer = untypedSourceLiteralPointerPointer.(core.LiteralPointerPointer)
		}
	}
	if sourceLiteralPointerPointer != nil {
		createdLiteral.SetLiteralValue(sourceLiteralPointerPointer.GetLiteralPointerId(hl).String(), hl)
	}
}

func getLiteralPointerVersion(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri, hl)
	if sourceLiteralPointerPointerRef == nil {
		log.Printf("In GetLiteralPointerVersion, the SourceLiteralPointerPointerRef was not found in the replicate")
		return
	}

	createdLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri, hl)
	if createdLiteralRef == nil {
		log.Printf("In GetLiteralPointerVersion, the CreatedLiteralRef was not found in the replicate")
		return
	}

	createdLiteral := createdLiteralRef.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, createdLiteralRef, hl)
		createdLiteralRef.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedSourceLiteralPointerPointer := sourceLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	var sourceLiteralPointerPointer core.LiteralPointerPointer
	if untypedSourceLiteralPointerPointer != nil {
		switch untypedSourceLiteralPointerPointer.(type) {
		case core.LiteralPointerPointer:
			sourceLiteralPointerPointer = untypedSourceLiteralPointerPointer.(core.LiteralPointerPointer)
		}
	}
	if sourceLiteralPointerPointer != nil {
		stringVersion := strconv.Itoa(sourceLiteralPointerPointer.GetLiteralPointerVersion(hl))
		createdLiteral.SetLiteralValue(stringVersion, hl)
	}
}

func setLiteralPointer(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerPointerSetLiteralPointerUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	literalPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri, hl)
	if literalPointerRef == nil {
		log.Printf("In SetLiteralPointer, the LiteralPointerRef was not found in the replicate")
		return
	}

	modifiedLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithAncestorUri(replicate, LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri, hl)
	if modifiedLiteralPointerPointerRef == nil {
		log.Printf("In SetLiteralPointer, the ModifiedLiteralPointerPointerRef was not found in the replicate")
		return
	}

	untypedBaseElement := modifiedLiteralPointerPointerRef.GetReferencedBaseElement(hl)
	literalPointer := literalPointerRef.GetReferencedLiteralPointer(hl)
	var modifiedLiteralPointerPointer core.LiteralPointerPointer
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.LiteralPointerPointer:
			modifiedLiteralPointerPointer = untypedBaseElement.(core.LiteralPointerPointer)
			modifiedLiteralPointerPointer.SetLiteralPointer(literalPointer, hl)
		default:
			log.Printf("In SetLiteralPointer, the ModifiedLiteralPointerPointerRef does not point to an LiteralPointer")
		}
	}
}

func UpdateRecoveredCoreLiteralPointerPointerFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralPointerPointerFunctions
	literalPointerPointerFunctions := uOfD.GetElementWithUri(LiteralPointerPointerFunctionsUri)
	if literalPointerPointerFunctions == nil {
		literalPointerPointerFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(literalPointerPointerFunctions, coreFunctionsElement, hl)
		core.SetName(literalPointerPointerFunctions, "LiteralPointerPointerFunctions", hl)
		core.SetUri(literalPointerPointerFunctions, LiteralPointerPointerFunctionsUri, hl)
	}

	// CreateAbstractLiteralPointerPointer
	createLiteralPointerPointer := uOfD.GetElementWithUri(LiteralPointerPointerCreateLiteralPointerPointerUri)
	if createLiteralPointerPointer == nil {
		createLiteralPointerPointer = uOfD.NewElement(hl)
		core.SetOwningElement(createLiteralPointerPointer, literalPointerPointerFunctions, hl)
		core.SetName(createLiteralPointerPointer, "CreateLiteralPointerPointer", hl)
		core.SetUri(createLiteralPointerPointer, LiteralPointerPointerCreateLiteralPointerPointerUri, hl)
	}
	// CreatedLiteralPointerPointerReference
	createdLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithUri(createLiteralPointerPointer, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri, hl)
	if createdLiteralPointerPointerRef == nil {
		createdLiteralPointerPointerRef = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(createdLiteralPointerPointerRef, createLiteralPointerPointer, hl)
		core.SetName(createdLiteralPointerPointerRef, "CreatedLiteralPointerdPointerRef", hl)
		core.SetUri(createdLiteralPointerPointerRef, LiteralPointerPointerCreateLiteralPointerPointerCreatedLiteralPointerPointerRefUri, hl)
	}

	// GetLiteralPointer
	getLiteralPointer := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerUri)
	if getLiteralPointer == nil {
		getLiteralPointer = uOfD.NewElement(hl)
		core.SetName(getLiteralPointer, "GetLiteralPointer", hl)
		core.SetOwningElement(getLiteralPointer, literalPointerPointerFunctions, hl)
		core.SetUri(getLiteralPointer, LiteralPointerPointerGetLiteralPointerUri, hl)
	}
	// GetLiteralPointer.SourceReference
	getLiteralPointerSourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithUri(getLiteralPointer, LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri, hl)
	if getLiteralPointerSourceLiteralPointerPointerRef == nil {
		getLiteralPointerSourceLiteralPointerPointerRef = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getLiteralPointerSourceLiteralPointerPointerRef, getLiteralPointer, hl)
		core.SetName(getLiteralPointerSourceLiteralPointerPointerRef, "SourceLiteralPointerPointerRef", hl)
		core.SetUri(getLiteralPointerSourceLiteralPointerPointerRef, LiteralPointerPointerGetLiteralPointerSourceLiteralPointerPointerRefUri, hl)
	}
	// GetLiteralPointerIndicatedLiteralPointerRef
	getLiteralPointerIndicatedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithUri(getLiteralPointer, LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	if getLiteralPointerIndicatedLiteralPointerRef == nil {
		getLiteralPointerIndicatedLiteralPointerRef = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(getLiteralPointerIndicatedLiteralPointerRef, getLiteralPointer, hl)
		core.SetName(getLiteralPointerIndicatedLiteralPointerRef, "IndicatedLiteralPointerRef", hl)
		core.SetUri(getLiteralPointerIndicatedLiteralPointerRef, LiteralPointerPointerGetLiteralPointerIndicatedLiteralPointerRefUri, hl)
	}

	// GetLiteralPointerId
	getLiteralPointerId := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerIdUri)
	if getLiteralPointerId == nil {
		getLiteralPointerId = uOfD.NewElement(hl)
		core.SetName(getLiteralPointerId, "GetLiteralPointerId", hl)
		core.SetOwningElement(getLiteralPointerId, literalPointerPointerFunctions, hl)
		core.SetUri(getLiteralPointerId, LiteralPointerPointerGetLiteralPointerIdUri, hl)
	}
	// GetLiteralPointerId.SourceLiteralPointerPointerRef
	getLiteralPointerIdSourceLiteralPointerPointerRef := core.GetChildBaseElementReferenceWithUri(getLiteralPointerId, LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri, hl)
	if getLiteralPointerIdSourceLiteralPointerPointerRef == nil {
		getLiteralPointerIdSourceLiteralPointerPointerRef = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getLiteralPointerIdSourceLiteralPointerPointerRef, getLiteralPointerId, hl)
		core.SetName(getLiteralPointerIdSourceLiteralPointerPointerRef, "SourceLiteralPointerPointerRef", hl)
		core.SetUri(getLiteralPointerIdSourceLiteralPointerPointerRef, LiteralPointerPointerGetLiteralPointerIdSourceLiteralPointerPointerRefUri, hl)
	}
	// GetLiteralPointerIdCreatedLiteralRef
	getLiteralPointerIdCreatedLiteralRef := core.GetChildLiteralReferenceWithUri(getLiteralPointerId, LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri, hl)
	if getLiteralPointerIdCreatedLiteralRef == nil {
		getLiteralPointerIdCreatedLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getLiteralPointerIdCreatedLiteralRef, getLiteralPointerId, hl)
		core.SetName(getLiteralPointerIdCreatedLiteralRef, "CreatedLiteralRef", hl)
		core.SetUri(getLiteralPointerIdCreatedLiteralRef, LiteralPointerPointerGetLiteralPointerIdCreatedLiteralUri, hl)
	}

	// GetLiteralPointerVersion
	getLiteralPointerVersion := uOfD.GetElementWithUri(LiteralPointerPointerGetLiteralPointerVersionUri)
	if getLiteralPointerVersion == nil {
		getLiteralPointerVersion = uOfD.NewElement(hl)
		core.SetName(getLiteralPointerVersion, "GetLiteralPointerVersion", hl)
		core.SetOwningElement(getLiteralPointerVersion, literalPointerPointerFunctions, hl)
		core.SetUri(getLiteralPointerVersion, LiteralPointerPointerGetLiteralPointerVersionUri, hl)
	}
	// GetLiteralPointerVersion.SourceReference
	getLiteralPointerVersionSourceReference := core.GetChildBaseElementReferenceWithUri(getLiteralPointerVersion, LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri, hl)
	if getLiteralPointerVersionSourceReference == nil {
		getLiteralPointerVersionSourceReference = uOfD.NewBaseElementReference(hl)
		core.SetOwningElement(getLiteralPointerVersionSourceReference, getLiteralPointerVersion, hl)
		core.SetName(getLiteralPointerVersionSourceReference, "SourceLiteralPointerRef", hl)
		core.SetUri(getLiteralPointerVersionSourceReference, LiteralPointerPointerGetLiteralPointerVersionSourceLiteralPointerPointerRefUri, hl)
	}
	// GetLiteralPointerVersionTargetLiteralReference
	getLiteralPointerVersionCreatedLiteralRef := core.GetChildLiteralReferenceWithUri(getLiteralPointerVersion, LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri, hl)
	if getLiteralPointerVersionCreatedLiteralRef == nil {
		getLiteralPointerVersionCreatedLiteralRef = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getLiteralPointerVersionCreatedLiteralRef, getLiteralPointerVersion, hl)
		core.SetName(getLiteralPointerVersionCreatedLiteralRef, "CreatedLiteralRef", hl)
		core.SetUri(getLiteralPointerVersionCreatedLiteralRef, LiteralPointerPointerGetLiteralPointerVersionCreatedLiteralRefUri, hl)
	}

	// SetLiteralPointer
	setLiteralPointer := uOfD.GetElementWithUri(LiteralPointerPointerSetLiteralPointerUri)
	if setLiteralPointer == nil {
		setLiteralPointer = uOfD.NewElement(hl)
		core.SetName(setLiteralPointer, "SetLiteralPointer", hl)
		core.SetOwningElement(setLiteralPointer, literalPointerPointerFunctions, hl)
		core.SetUri(setLiteralPointer, LiteralPointerPointerSetLiteralPointerUri, hl)
	}
	// SetLiteralPointer.LiteralPointerRef
	setLiteralPointerLiteralPointerRef := core.GetChildLiteralPointerReferenceWithUri(setLiteralPointer, LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri, hl)
	if setLiteralPointerLiteralPointerRef == nil {
		setLiteralPointerLiteralPointerRef = uOfD.NewLiteralPointerReference(hl)
		core.SetName(setLiteralPointerLiteralPointerRef, "LiteralPointerRef", hl)
		core.SetOwningElement(setLiteralPointerLiteralPointerRef, setLiteralPointer, hl)
		core.SetUri(setLiteralPointerLiteralPointerRef, LiteralPointerPointerSetLiteralPointerLiteralPointerRefUri, hl)
	}
	setLiteralPointerTargetLiteralPointerPointerReference := core.GetChildBaseElementReferenceWithUri(setLiteralPointer, LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri, hl)
	if setLiteralPointerTargetLiteralPointerPointerReference == nil {
		setLiteralPointerTargetLiteralPointerPointerReference = uOfD.NewBaseElementReference(hl)
		core.SetName(setLiteralPointerTargetLiteralPointerPointerReference, "ModifiedLiteralPointerPointerRef", hl)
		core.SetOwningElement(setLiteralPointerTargetLiteralPointerPointerReference, setLiteralPointer, hl)
		core.SetUri(setLiteralPointerTargetLiteralPointerPointerReference, LiteralPointerPointerSetLiteralPointerModifiedLiteralPointerPointerRefUri, hl)
	}
}

func literalPointerPointerFunctionsInit() {
	core.GetCore().AddFunction(LiteralPointerPointerCreateLiteralPointerPointerUri, createLiteralPointerPointer)
	core.GetCore().AddFunction(LiteralPointerPointerGetLiteralPointerUri, getLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerPointerGetLiteralPointerIdUri, getLiteralPointerId)
	core.GetCore().AddFunction(LiteralPointerPointerGetLiteralPointerVersionUri, getLiteralPointerVersion)
	core.GetCore().AddFunction(LiteralPointerPointerSetLiteralPointerUri, setLiteralPointer)
}
