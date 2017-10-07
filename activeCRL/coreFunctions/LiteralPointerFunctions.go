package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"log"
	"strconv"
	"sync"
)

var LiteralPointerFunctionsUri string = CoreFunctionsPrefix + "LiteralPointerFunctions"

var LiteralPointerCreateNameLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointer/CreateNameLiteralPointer"
var LiteralPointerCreateNameLiteralPointerCreatedLiteralPointerRefUri = CoreFunctionsPrefix + "LiteralPointer/CreateNameLiteralPointer/CreatedLiteralPointerRef"

var LiteralPointerCreateDefinitionLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointer/CreateDefinitionLiteralPointer"
var LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri = CoreFunctionsPrefix + "LiteralPointer/CreateDefinitionLiteralPointer/CreatedLiteralPointerRef"

var LiteralPointerCreateUriLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointer/CreateUriLiteralPointer"
var LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri = CoreFunctionsPrefix + "LiteralPointer/CreateUriLiteralPointer/CreatedLiteralPointerRef"

var LiteralPointerCreateValueLiteralPointerUri string = CoreFunctionsPrefix + "LiteralPointer/CreateValueLiteralPointer"
var LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri = CoreFunctionsPrefix + "LiteralPointer/CreateValueLiteralPointer/CreatedLiteralPointerRef"

var LiteralPointerGetLiteralUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteral"
var LiteralPointerGetLiteralSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteral/SourceLiteralPointerRef"
var LiteralPointerGetLiteralIndicatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteral/IndicatedLiteralRef"

var LiteralPointerGetLiteralIdUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralId"
var LiteralPointerGetLiteralIdSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralId/SourceLiteralPointerRef"
var LiteralPointerGetLiteralIdCreatedLiteralUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralId/CreatedLiteralRef"

var LiteralPointerGetLiteralPointerRoleUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralPointerRole"
var LiteralPointerGetLiteralPointerRoleSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralPointerRole/SourceLiteralPointerRef"
var LiteralPointerGetLiteralPointerRoleCreatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralPointerRole/CreatedLiteralRef"

var LiteralPointerGetLiteralVersionUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralVersion"
var LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralVersion/SourceLiteralPointerRef"
var LiteralPointerGetLiteralVersionCreatedLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointer/GetLiteralVersion/CreatedLiteralRef"

var LiteralPointerSetLiteralUri string = CoreFunctionsPrefix + "LiteralPointer/SetLiteral"
var LiteralPointerSetLiteralLiteralRefUri string = CoreFunctionsPrefix + "LiteralPointer/SetLiteral/LiteralRef"
var LiteralPointerSetLiteralModifiedLiteralPointerRefUri string = CoreFunctionsPrefix + "LiteralPointer/SetLiteral/ModifiedLiteralPointerRef"

func createNameLiteralPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createNameLiteralPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(element, LiteralPointerCreateNameLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, element, hl)
		core.SetName(createdLiteralPointerReference, "CreatedLiteralPointerReference", hl)
		rootCreatedLiteralReference := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateNameLiteralPointerCreatedLiteralPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralReference, hl)
	}
	createdLiteralPointer := createdLiteralPointerReference.GetReferencedLiteralPointer(hl)
	if createdLiteralPointer == nil {
		createdLiteralPointer = uOfD.NewNameLiteralPointer(hl)
		createdLiteralPointerReference.SetReferencedLiteralPointer(createdLiteralPointer, hl)
	}
}

func createDefinitionLiteralPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createDefinitionLiteralPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(element, LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, element, hl)
		core.SetName(createdLiteralPointerReference, "CreatedLiteralPointerReference", hl)
		rootCreatedLiteralPointerReference := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralPointerReference, hl)
	}
	createdLiteralPointer := createdLiteralPointerReference.GetReferencedLiteralPointer(hl)
	if createdLiteralPointer == nil {
		createdLiteralPointer = uOfD.NewDefinitionLiteralPointer(hl)
		createdLiteralPointerReference.SetReferencedLiteralPointer(createdLiteralPointer, hl)
	}
}

func createUriLiteralPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createUriLiteralPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(element, LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, element, hl)
		core.SetName(createdLiteralPointerReference, "CreatedLiteralPointerReference", hl)
		rootCreatedLiteralPointerReference := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralPointerReference, hl)
	}
	createdLiteralPointer := createdLiteralPointerReference.GetReferencedLiteralPointer(hl)
	if createdLiteralPointer == nil {
		createdLiteralPointer = uOfD.NewUriLiteralPointer(hl)
		createdLiteralPointerReference.SetReferencedLiteralPointer(createdLiteralPointer, hl)
	}
}

func createValueLiteralPointer(element core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	//	log.Printf("In createValueLiteralPointer")
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(element)
	uOfD := element.GetUniverseOfDiscourse(hl)
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithAncestorUri(element, LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, element, hl)
		core.SetName(createdLiteralPointerReference, "CreatedLiteralPointerReference", hl)
		rootCreatedLiteralPointerReference := uOfD.GetLiteralPointerReferenceWithUri(LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri)
		refinement := uOfD.NewRefinement(hl)
		core.SetOwningElement(refinement, createdLiteralPointerReference, hl)
		refinement.SetRefinedElement(createdLiteralPointerReference, hl)
		refinement.SetAbstractElement(rootCreatedLiteralPointerReference, hl)
	}
	createdLiteralPointer := createdLiteralPointerReference.GetReferencedLiteralPointer(hl)
	if createdLiteralPointer == nil {
		createdLiteralPointer = uOfD.NewValueLiteralPointer(hl)
		createdLiteralPointerReference.SetReferencedLiteralPointer(createdLiteralPointer, hl)
	}
}

func getLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerGetLiteralUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteral, the SourceReference was not found in the replicate")
		return
	}

	indicatedLiteralRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIndicatedLiteralRefUri, hl)
	if indicatedLiteralRef == nil {
		log.Printf("In GetLiteral, the TargetLiteralPointerReference was not found in the replicate")
		return
	}

	targetLiteral := indicatedLiteralRef.GetReferencedLiteral(hl)
	sourceLiteralPointer := sourceReference.GetReferencedLiteralPointer(hl)
	var sourceLiteral core.Literal
	if sourceLiteralPointer != nil {
		sourceLiteral = sourceLiteralPointer.GetLiteral(hl)
	}
	if sourceLiteral != targetLiteral {
		indicatedLiteralRef.SetReferencedLiteral(sourceLiteral, hl)
	}
}

func getLiteralId(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerGetLiteralIdUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIdSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteralId, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralIdCreatedLiteralUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetLiteralId, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	sourceLiteralPointer := sourceReference.GetReferencedLiteralPointer(hl)
	if sourceLiteralPointer != nil {
		createdLiteral.SetLiteralValue(sourceLiteralPointer.GetLiteralId(hl).String(), hl)
	}
}

func getLiteralPointerRole(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerGetLiteralPointerRoleUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralPointerRoleSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteralPointerRole, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralPointerRoleCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetLiteralPointerRole, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	untypedBaseElement := sourceReference.GetReferencedLiteralPointer(hl)
	var sourceLiteralPointer core.LiteralPointer
	if untypedBaseElement != nil {
		switch untypedBaseElement.(type) {
		case core.LiteralPointer:
			sourceLiteralPointer = untypedBaseElement.(core.LiteralPointer)
		default:
			log.Printf("In GetLiteralPointerRole, the SourceBaseElement is not a LiteralPointer")
		}
	}
	if sourceLiteralPointer != nil {
		stringBaseElementRole := sourceLiteralPointer.GetLiteralPointerRole(hl).RoleToString()
		createdLiteral.SetLiteralValue(stringBaseElementRole, hl)
	}
}

func getLiteralVersion(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerGetLiteralVersionUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	sourceReference := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri, hl)
	if sourceReference == nil {
		log.Printf("In GetLiteralVersion, the SourceReference was not found in the replicate")
		return
	}

	targetLiteralReference := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerGetLiteralVersionCreatedLiteralRefUri, hl)
	if targetLiteralReference == nil {
		log.Printf("In GetLiteralVersion, the TargetLiteralReference was not found in the replicate")
		return
	}

	createdLiteral := targetLiteralReference.GetReferencedLiteral(hl)
	if createdLiteral == nil {
		createdLiteral = uOfD.NewLiteral(hl)
		core.SetOwningElement(createdLiteral, targetLiteralReference, hl)
		targetLiteralReference.SetReferencedLiteral(createdLiteral, hl)
	}

	sourceLiteralPointer := sourceReference.GetReferencedLiteralPointer(hl)
	if sourceLiteralPointer != nil {
		stringBaseElementVersion := strconv.Itoa(sourceLiteralPointer.GetLiteralVersion(hl))
		createdLiteral.SetLiteralValue(stringBaseElementVersion, hl)
	}
}

func setLiteral(replicate core.Element, changeNotifications []*core.ChangeNotification, wg *sync.WaitGroup) {
	hl := core.NewHeldLocks(wg)
	defer hl.ReleaseLocks()
	hl.LockBaseElement(replicate)
	uOfD := replicate.GetUniverseOfDiscourse(hl)

	original := uOfD.GetElementWithUri(LiteralPointerSetLiteralUri)
	core.ReplicateAsRefinement(original, replicate, hl)

	literalRef := core.GetChildLiteralReferenceWithAncestorUri(replicate, LiteralPointerSetLiteralLiteralRefUri, hl)
	if literalRef == nil {
		log.Printf("In SetLiteral, the LiteralReference was not found in the replicate")
		return
	}

	modifiedLiteralPointerRef := core.GetChildLiteralPointerReferenceWithAncestorUri(replicate, LiteralPointerSetLiteralModifiedLiteralPointerRefUri, hl)
	if modifiedLiteralPointerRef == nil {
		log.Printf("In SetLiteral, the TargetLiteralPointerReference was not found in the replicate")
		return
	}

	modifiedLiteralPointer := modifiedLiteralPointerRef.GetReferencedLiteralPointer(hl)
	sourceLiteral := literalRef.GetReferencedLiteral(hl)
	if modifiedLiteralPointer != nil {
		modifiedLiteralPointer.SetLiteral(sourceLiteral, hl)
	}
}

func UpdateRecoveredCoreLiteralPointerFunctions(coreFunctionsElement core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) {

	// LiteralPointerFunctions
	literalPointerFunctions := uOfD.GetElementWithUri(LiteralPointerFunctionsUri)
	if literalPointerFunctions == nil {
		literalPointerFunctions = uOfD.NewElement(hl)
		core.SetOwningElement(literalPointerFunctions, coreFunctionsElement, hl)
		core.SetName(literalPointerFunctions, "LiteralPointerFunctions", hl)
		core.SetUri(literalPointerFunctions, LiteralPointerFunctionsUri, hl)
	}

	// CreateNameLiteralPointerElement
	createNameLiteralPointer := uOfD.GetElementWithUri(LiteralPointerCreateNameLiteralPointerUri)
	if createNameLiteralPointer == nil {
		createNameLiteralPointer = uOfD.NewElement(hl)
		core.SetOwningElement(createNameLiteralPointer, literalPointerFunctions, hl)
		core.SetName(createNameLiteralPointer, "CreateNameLiteralPointerLiteralPointer", hl)
		core.SetUri(createNameLiteralPointer, LiteralPointerCreateNameLiteralPointerUri, hl)
	}
	// CreatedLiteralReference
	createdLiteralPointerReference := core.GetChildLiteralPointerReferenceWithUri(createNameLiteralPointer, LiteralPointerCreateNameLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, createNameLiteralPointer, hl)
		core.SetName(createdLiteralPointerReference, "CreateNameLiteralPointerdLiteralPointerRef", hl)
		core.SetUri(createdLiteralPointerReference, LiteralPointerCreateNameLiteralPointerCreatedLiteralPointerRefUri, hl)
	}

	// CreateDefinitionLiteralPointerElement
	createDefinitionLiteralPointer := uOfD.GetElementWithUri(LiteralPointerCreateDefinitionLiteralPointerUri)
	if createDefinitionLiteralPointer == nil {
		createDefinitionLiteralPointer = uOfD.NewElement(hl)
		core.SetOwningElement(createDefinitionLiteralPointer, literalPointerFunctions, hl)
		core.SetName(createDefinitionLiteralPointer, "CreateDefinitionLiteralPointerLiteralPointer", hl)
		core.SetUri(createDefinitionLiteralPointer, LiteralPointerCreateDefinitionLiteralPointerUri, hl)
	}
	// CreatedLiteralReference
	createdLiteralPointerReference = core.GetChildLiteralPointerReferenceWithUri(createDefinitionLiteralPointer, LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, createDefinitionLiteralPointer, hl)
		core.SetName(createdLiteralPointerReference, "CreateDefinitionLiteralPointerdLiteralPointerRef", hl)
		core.SetUri(createdLiteralPointerReference, LiteralPointerCreateDefinitionLiteralPointerCreatedLiteralPointerRefUri, hl)
	}

	// CreateUriLiteralPointerElement
	createUriLiteralPointer := uOfD.GetElementWithUri(LiteralPointerCreateUriLiteralPointerUri)
	if createUriLiteralPointer == nil {
		createUriLiteralPointer = uOfD.NewElement(hl)
		core.SetOwningElement(createUriLiteralPointer, literalPointerFunctions, hl)
		core.SetName(createUriLiteralPointer, "CreateUriLiteralPointerLiteralPointer", hl)
		core.SetUri(createUriLiteralPointer, LiteralPointerCreateUriLiteralPointerUri, hl)
	}
	// CreatedLiteralReference
	createdLiteralPointerReference = core.GetChildLiteralPointerReferenceWithUri(createUriLiteralPointer, LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, createUriLiteralPointer, hl)
		core.SetName(createdLiteralPointerReference, "CreateUriLiteralPointerdLiteralPointerRef", hl)
		core.SetUri(createdLiteralPointerReference, LiteralPointerCreateUriLiteralPointerCreatedLiteralPointerRefUri, hl)
	}

	// CreateValueLiteralPointerElement
	createValueLiteralPointer := uOfD.GetElementWithUri(LiteralPointerCreateValueLiteralPointerUri)
	if createValueLiteralPointer == nil {
		createValueLiteralPointer = uOfD.NewElement(hl)
		core.SetOwningElement(createValueLiteralPointer, literalPointerFunctions, hl)
		core.SetName(createValueLiteralPointer, "CreateValueLiteralPointerLiteralPointer", hl)
		core.SetUri(createValueLiteralPointer, LiteralPointerCreateValueLiteralPointerUri, hl)
	}
	// CreatedLiteralReference
	createdLiteralPointerReference = core.GetChildLiteralPointerReferenceWithUri(createValueLiteralPointer, LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri, hl)
	if createdLiteralPointerReference == nil {
		createdLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(createdLiteralPointerReference, createValueLiteralPointer, hl)
		core.SetName(createdLiteralPointerReference, "CreateValueLiteralPointerdLiteralPointerRef", hl)
		core.SetUri(createdLiteralPointerReference, LiteralPointerCreateValueLiteralPointerCreatedLiteralPointerRefUri, hl)
	}

	// GetLiteral
	getLiteral := uOfD.GetElementWithUri(LiteralPointerGetLiteralUri)
	if getLiteral == nil {
		getLiteral = uOfD.NewElement(hl)
		core.SetName(getLiteral, "GetLiteral", hl)
		core.SetOwningElement(getLiteral, literalPointerFunctions, hl)
		core.SetUri(getLiteral, LiteralPointerGetLiteralUri, hl)
	}
	// GetLiteral.SourceReference
	getLiteralSourceReference := core.GetChildLiteralPointerReferenceWithUri(getLiteral, LiteralPointerGetLiteralSourceLiteralPointerRefUri, hl)
	if getLiteralSourceReference == nil {
		getLiteralSourceReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(getLiteralSourceReference, getLiteral, hl)
		core.SetName(getLiteralSourceReference, "SourceLiteralPointerRef", hl)
		core.SetUri(getLiteralSourceReference, LiteralPointerGetLiteralSourceLiteralPointerRefUri, hl)
	}
	// GetLiteralTargetLiteralPointerReference
	getLiteralTargetReference := core.GetChildLiteralReferenceWithUri(getLiteral, LiteralPointerGetLiteralIndicatedLiteralRefUri, hl)
	if getLiteralTargetReference == nil {
		getLiteralTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getLiteralTargetReference, getLiteral, hl)
		core.SetName(getLiteralTargetReference, "IndicatedBaseElementRef", hl)
		core.SetUri(getLiteralTargetReference, LiteralPointerGetLiteralIndicatedLiteralRefUri, hl)
	}

	// GetLiteralId
	getLiteralId := uOfD.GetElementWithUri(LiteralPointerGetLiteralIdUri)
	if getLiteralId == nil {
		getLiteralId = uOfD.NewElement(hl)
		core.SetName(getLiteralId, "GetLiteralId", hl)
		core.SetOwningElement(getLiteralId, literalPointerFunctions, hl)
		core.SetUri(getLiteralId, LiteralPointerGetLiteralIdUri, hl)
	}
	// GetLiteralId.SourceReference
	getLiteralIdSourceReference := core.GetChildLiteralPointerReferenceWithUri(getLiteralId, LiteralPointerGetLiteralIdSourceLiteralPointerRefUri, hl)
	if getLiteralIdSourceReference == nil {
		getLiteralIdSourceReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(getLiteralIdSourceReference, getLiteralId, hl)
		core.SetName(getLiteralIdSourceReference, "SourceLiteralPointerRef", hl)
		core.SetUri(getLiteralIdSourceReference, LiteralPointerGetLiteralIdSourceLiteralPointerRefUri, hl)
	}
	// GetLiteralIdTargetLiteralReference
	getLiteralIdTargetReference := core.GetChildLiteralReferenceWithUri(getLiteralId, LiteralPointerGetLiteralIdCreatedLiteralUri, hl)
	if getLiteralIdTargetReference == nil {
		getLiteralIdTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getLiteralIdTargetReference, getLiteralId, hl)
		core.SetName(getLiteralIdTargetReference, "CreatedLiteralRef", hl)
		core.SetUri(getLiteralIdTargetReference, LiteralPointerGetLiteralIdCreatedLiteralUri, hl)
	}

	// GetLiteralPointerRole
	getLiteralPointerRole := uOfD.GetElementWithUri(LiteralPointerGetLiteralPointerRoleUri)
	if getLiteralPointerRole == nil {
		getLiteralPointerRole = uOfD.NewElement(hl)
		core.SetName(getLiteralPointerRole, "GetLiteralPointerRole", hl)
		core.SetOwningElement(getLiteralPointerRole, literalPointerFunctions, hl)
		core.SetUri(getLiteralPointerRole, LiteralPointerGetLiteralPointerRoleUri, hl)
	}
	// GetLiteralPointerRole.SourceReference
	getLiteralPointerRoleSourceReference := core.GetChildLiteralPointerReferenceWithUri(getLiteralPointerRole, LiteralPointerGetLiteralPointerRoleSourceLiteralPointerRefUri, hl)
	if getLiteralPointerRoleSourceReference == nil {
		getLiteralPointerRoleSourceReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(getLiteralPointerRoleSourceReference, getLiteralPointerRole, hl)
		core.SetName(getLiteralPointerRoleSourceReference, "SourceLiteralPointerRef", hl)
		core.SetUri(getLiteralPointerRoleSourceReference, LiteralPointerGetLiteralPointerRoleSourceLiteralPointerRefUri, hl)
	}
	// GetLiteralPointerRoleTargetLiteralReference
	getLiteralPointerRoleTargetReference := core.GetChildLiteralReferenceWithUri(getLiteralPointerRole, LiteralPointerGetLiteralPointerRoleCreatedLiteralRefUri, hl)
	if getLiteralPointerRoleTargetReference == nil {
		getLiteralPointerRoleTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getLiteralPointerRoleTargetReference, getLiteralPointerRole, hl)
		core.SetName(getLiteralPointerRoleTargetReference, "CreatedLiteralRef", hl)
		core.SetUri(getLiteralPointerRoleTargetReference, LiteralPointerGetLiteralPointerRoleCreatedLiteralRefUri, hl)
	}

	// GetLiteralVersion
	getLiteralVersion := uOfD.GetElementWithUri(LiteralPointerGetLiteralVersionUri)
	if getLiteralVersion == nil {
		getLiteralVersion = uOfD.NewElement(hl)
		core.SetName(getLiteralVersion, "GetLiteralVersion", hl)
		core.SetOwningElement(getLiteralVersion, literalPointerFunctions, hl)
		core.SetUri(getLiteralVersion, LiteralPointerGetLiteralVersionUri, hl)
	}
	// GetLiteralVersion.SourceReference
	getLiteralVersionSourceReference := core.GetChildLiteralPointerReferenceWithUri(getLiteralVersion, LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri, hl)
	if getLiteralVersionSourceReference == nil {
		getLiteralVersionSourceReference = uOfD.NewLiteralPointerReference(hl)
		core.SetOwningElement(getLiteralVersionSourceReference, getLiteralVersion, hl)
		core.SetName(getLiteralVersionSourceReference, "SourceLiteralPointerRef", hl)
		core.SetUri(getLiteralVersionSourceReference, LiteralPointerGetLiteralVersionSourceLiteralPointerRefUri, hl)
	}
	// GetLiteralVersionTargetLiteralReference
	getLiteralVersionTargetReference := core.GetChildLiteralReferenceWithUri(getLiteralVersion, LiteralPointerGetLiteralVersionCreatedLiteralRefUri, hl)
	if getLiteralVersionTargetReference == nil {
		getLiteralVersionTargetReference = uOfD.NewLiteralReference(hl)
		core.SetOwningElement(getLiteralVersionTargetReference, getLiteralVersion, hl)
		core.SetName(getLiteralVersionTargetReference, "CreatedLiteralRef", hl)
		core.SetUri(getLiteralVersionTargetReference, LiteralPointerGetLiteralVersionCreatedLiteralRefUri, hl)
	}

	// SetLiteral
	setLiteral := uOfD.GetElementWithUri(LiteralPointerSetLiteralUri)
	if setLiteral == nil {
		setLiteral = uOfD.NewElement(hl)
		core.SetName(setLiteral, "SetLiteral", hl)
		core.SetOwningElement(setLiteral, literalPointerFunctions, hl)
		core.SetUri(setLiteral, LiteralPointerSetLiteralUri, hl)
	}
	// SetLiteral.LiteralReference
	setLiteralLiteralReference := core.GetChildLiteralReferenceWithUri(setLiteral, LiteralPointerSetLiteralLiteralRefUri, hl)
	if setLiteralLiteralReference == nil {
		setLiteralLiteralReference = uOfD.NewLiteralReference(hl)
		core.SetName(setLiteralLiteralReference, "BaseElementRef", hl)
		core.SetOwningElement(setLiteralLiteralReference, setLiteral, hl)
		core.SetUri(setLiteralLiteralReference, LiteralPointerSetLiteralLiteralRefUri, hl)
	}
	setLiteralTargetLiteralPointerReference := core.GetChildLiteralPointerReferenceWithUri(setLiteral, LiteralPointerSetLiteralModifiedLiteralPointerRefUri, hl)
	if setLiteralTargetLiteralPointerReference == nil {
		setLiteralTargetLiteralPointerReference = uOfD.NewLiteralPointerReference(hl)
		core.SetName(setLiteralTargetLiteralPointerReference, "ModifiedLiteralPointerRef", hl)
		core.SetOwningElement(setLiteralTargetLiteralPointerReference, setLiteral, hl)
		core.SetUri(setLiteralTargetLiteralPointerReference, LiteralPointerSetLiteralModifiedLiteralPointerRefUri, hl)
	}
}

func literalPointerFunctionsInit() {
	core.GetCore().AddFunction(LiteralPointerCreateNameLiteralPointerUri, createNameLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerCreateDefinitionLiteralPointerUri, createDefinitionLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerCreateUriLiteralPointerUri, createUriLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerCreateValueLiteralPointerUri, createValueLiteralPointer)
	core.GetCore().AddFunction(LiteralPointerGetLiteralUri, getLiteral)
	core.GetCore().AddFunction(LiteralPointerGetLiteralIdUri, getLiteralId)
	core.GetCore().AddFunction(LiteralPointerGetLiteralPointerRoleUri, getLiteralPointerRole)
	core.GetCore().AddFunction(LiteralPointerGetLiteralVersionUri, getLiteralVersion)
	core.GetCore().AddFunction(LiteralPointerSetLiteralUri, setLiteral)
}
