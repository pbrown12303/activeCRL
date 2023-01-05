package crlmapsdomain

import (
	"log"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pkg/errors"
)

// CrlMapsDomainURI is the URI for the concpet space that defines the CRL Data Types
var CrlMapsDomainURI = "http://activeCRL.com/crlmaps/CrlMaps"

// Basic Map

// CrlMapURI is the uri of the ancestor of all maps
var CrlMapURI = CrlMapsDomainURI + "/Map"

// CrlMapSourceURI is the URI of the CrlMap source
var CrlMapSourceURI = CrlMapURI + "/Source"

// CrlMapTargetURI is the URI of the CrlMap target
var CrlMapTargetURI = CrlMapURI + "/Target"

// One to One Map

// CrlOneToOneMapURI is the URI for the Element-to-Element Map
var CrlOneToOneMapURI = CrlMapsDomainURI + "/OneToOneMap"

// CrlOneToOneMapRefinementURI is the URI of the refinement showing it to be a refinement of CrlMap
var CrlOneToOneMapRefinementURI = CrlOneToOneMapURI + "/Refinement"

// CrlOneToOneMapSourceReferenceURI is the URI for the source reference
var CrlOneToOneMapSourceReferenceURI = CrlOneToOneMapURI + "/SourceReference"

// CrlOneToOneMapSourceReferenceRefinementURI is the URI of the refinement showing it to be a refinement of CrlMapSource
var CrlOneToOneMapSourceReferenceRefinementURI = CrlOneToOneMapSourceReferenceURI + "/Refinement"

// CrlOneToOneMapTargetReferenceURI is the URI for the target reference
var CrlOneToOneMapTargetReferenceURI = CrlOneToOneMapURI + "/TargetReference"

// CrlOneToOneMapTargetReferenceRefinementURI is the URI of the refinement showing it to be a refinement of CrlMapTarget
var CrlOneToOneMapTargetReferenceRefinementURI = CrlOneToOneMapTargetReferenceURI + "/Refinement"

// Reference to Element Map

// CrlReferenceToElementMapURI is the URI for the Reference to Element Map
var CrlReferenceToElementMapURI = CrlMapsDomainURI + "/ReferenceToElementMap"

// CrlReferenceToElementMapRefinementURI is the URI for the refinement from CrlMap
var CrlReferenceToElementMapRefinementURI = CrlReferenceToElementMapURI + "/Refinement"

// CrlReferenceToElementMapSourceURI is the URI for the source
var CrlReferenceToElementMapSourceURI = CrlReferenceToElementMapURI + "/Source"

// CrlReferenceToElementMapSourceRefinementURI is the URI for the refinement from CrlMapSource
var CrlReferenceToElementMapSourceRefinementURI = CrlReferenceToElementMapSourceURI + "/Refinement"

// CrlReferenceToElementMapTargetURI is the URI for the target
var CrlReferenceToElementMapTargetURI = CrlReferenceToElementMapURI + "/Target"

// CrlReferenceToElementMapTargetRefinementURI is the URI for the refinement from CrlMapTarget
var CrlReferenceToElementMapTargetRefinementURI = CrlReferenceToElementMapTargetURI + "Refinement"

// ID to Reference Map

// CrlIDToReferenceMapURI is the URI for a map from an attribute that is an ID to a Reference
var CrlIDToReferenceMapURI = CrlMapsDomainURI + "/IDToReferenceMap"

// CrlIDToReferenceMapRefinementURI is the URI for the refinement from CrlMap
var CrlIDToReferenceMapRefinementURI = CrlIDToReferenceMapURI + "/Refinement"

// CrlIDToReferenceMapSourceURI is the URI for the source reference
var CrlIDToReferenceMapSourceURI = CrlIDToReferenceMapURI + "/Source"

// CrlIDToReferenceMapSourceRefinementURI is the URI for the refinement from CrlMapSource
var CrlIDToReferenceMapSourceRefinementURI = CrlIDToReferenceMapSourceURI + "/Refinement"

// CrlIDToReferenceMapTargetURI is the URI for the target reference
var CrlIDToReferenceMapTargetURI = CrlIDToReferenceMapURI + "/Target"

// CrlIDToReferenceMapTargetRefinementURI is the URI for the refinement from CrlMapTarget
var CrlIDToReferenceMapTargetRefinementURI = CrlIDToReferenceMapTargetURI + "/Refinement"

// BuildCrlMapsDomain constructs the domain for CRL maps
func BuildCrlMapsDomain(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) error {
	crlMapsDomain, err1 := uOfD.NewOwnedElement(nil, "CrlMapsDomain", trans, CrlMapsDomainURI)
	if err1 != nil {
		return errors.Wrap(err1, "crlmaps.BuildCrlMapsDomain failed")
	}

	// CrlMap
	crlMap, err2 := uOfD.NewOwnedElement(crlMapsDomain, "CrlMap", trans, CrlMapURI)
	if err2 != nil {
		return errors.Wrap(err2, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlMapSource, err3 := uOfD.NewOwnedReference(crlMap, "MapSource", trans, CrlMapSourceURI)
	if err3 != nil {
		return errors.Wrap(err3, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlMapTarget, err4 := uOfD.NewOwnedReference(crlMap, "MapTarget", trans, CrlMapTargetURI)
	if err4 != nil {
		return errors.Wrap(err4, "crlmaps.BuildCrlMapsDomain failed")
	}

	// One to One Map
	crlOneToOneMap, err5 := uOfD.NewOwnedElement(crlMapsDomain, "CrlOneToOneMap", trans, CrlOneToOneMapURI)
	if err5 != nil {
		return errors.Wrap(err5, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err6 := uOfD.NewCompleteRefinement(crlMap, crlOneToOneMap, "Refines CrlMap", trans, CrlOneToOneMapRefinementURI)
	if err6 != nil {
		return errors.Wrap(err6, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlOneToOneMapSourceReference, err7 := uOfD.NewOwnedReference(crlOneToOneMap, "SourceReference", trans, CrlOneToOneMapSourceReferenceURI)
	if err7 != nil {
		return errors.Wrap(err7, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err8 := uOfD.NewCompleteRefinement(crlMapSource, crlOneToOneMapSourceReference, "Refines CrlMapSource", trans, CrlOneToOneMapSourceReferenceRefinementURI)
	if err8 != nil {
		return errors.Wrap(err8, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlOneToOneMapTargetReference, err9 := uOfD.NewOwnedReference(crlOneToOneMap, "TargetReference", trans, CrlOneToOneMapTargetReferenceURI)
	if err9 != nil {
		return errors.Wrap(err9, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err10 := uOfD.NewCompleteRefinement(crlMapTarget, crlOneToOneMapTargetReference, "Refines CrlMapTarget", trans, CrlOneToOneMapTargetReferenceRefinementURI)
	if err10 != nil {
		return errors.Wrap(err10, "crlmaps.BuildCrlMapsDomain failed")
	}

	// Reference To Element Map
	crlReferenceToElementMap, err11 := uOfD.NewOwnedElement(crlMapsDomain, "ReferenceToElementMap", trans, CrlReferenceToElementMapURI)
	if err11 != nil {
		return errors.Wrap(err11, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err12 := uOfD.NewCompleteRefinement(crlMap, crlReferenceToElementMap, "Refinement", trans, CrlReferenceToElementMapRefinementURI)
	if err12 != nil {
		return errors.Wrap(err12, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlReferenceToElementMapSource, err13 := uOfD.NewOwnedReference(crlReferenceToElementMap, "Source", trans, CrlReferenceToElementMapSourceURI)
	if err13 != nil {
		return errors.Wrap(err13, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err14 := uOfD.NewCompleteRefinement(crlMapSource, crlReferenceToElementMapSource, "Refinement", trans, CrlReferenceToElementMapSourceRefinementURI)
	if err14 != nil {
		return errors.Wrap(err14, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlReferenceToElementMapTarget, err15 := uOfD.NewOwnedReference(crlReferenceToElementMap, "Target", trans, CrlReferenceToElementMapTargetURI)
	if err15 != nil {
		return errors.Wrap(err15, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err16 := uOfD.NewCompleteRefinement(crlMapTarget, crlReferenceToElementMapTarget, "Refinement", trans, CrlReferenceToElementMapTargetRefinementURI)
	if err16 != nil {
		return errors.Wrap(err16, "crlmaps.BuildCrlMapsDomain failed")
	}

	// ID to Reference Map
	crlIDToReferenceMap, err17 := uOfD.NewOwnedElement(crlMapsDomain, "IDToReferenceMap", trans, CrlIDToReferenceMapURI)
	if err17 != nil {
		return errors.Wrap(err17, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err18 := uOfD.NewCompleteRefinement(crlMap, crlIDToReferenceMap, "Refinement", trans, CrlIDToReferenceMapRefinementURI)
	if err18 != nil {
		return errors.Wrap(err18, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlIDToReferenceMapSource, err19 := uOfD.NewOwnedReference(crlIDToReferenceMap, "Source", trans, CrlIDToReferenceMapSourceURI)
	if err19 != nil {
		return errors.Wrap(err19, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err20 := uOfD.NewCompleteRefinement(crlMapSource, crlIDToReferenceMapSource, "Source", trans, CrlIDToReferenceMapSourceRefinementURI)
	if err20 != nil {
		return errors.Wrap(err20, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlIDToReferenceMapTarget, err21 := uOfD.NewOwnedReference(crlIDToReferenceMap, "Target", trans, CrlIDToReferenceMapTargetURI)
	if err21 != nil {
		return errors.Wrap(err21, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err22 := uOfD.NewCompleteRefinement(crlMapTarget, crlIDToReferenceMapTarget, "Refinement", trans, CrlIDToReferenceMapTargetRefinementURI)
	if err22 != nil {
		return errors.Wrap(err22, "crlmaps.BuildCrlMapsDomain failed")
	}

	err := crlMapsDomain.SetReadOnlyRecursively(true, trans)
	if err != nil {
		return errors.Wrap(err, "CrlMaps BuildCrlMapsDomain failed")
	}
	err = crlMapsDomain.SetIsCoreRecursively(trans)
	if err != nil {
		return errors.Wrap(err, "CrlMaps BuildCrlMapsDomain failed")
	}

	uOfD.AddFunction(CrlOneToOneMapURI, executeOneToOneMap)
	return nil
}

// executeOneToOneMap performs the mapping function
func executeOneToOneMap(mapInstance core.Element, notification *core.ChangeNotification, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	trans.WriteLockElement(mapInstance)

	log.Printf("Executing executeOneToOneMap for map labeled %s", mapInstance.GetLabel(trans))

	// As an initial assumption, it probably doesn't matter what kind of notification has been received.
	// Validate that this instance is a refinement of an element that is, in turn, a refinement of CrlOneToOneMap
	var immediateAbstractions = map[string]core.Element{}
	mapInstance.FindImmediateAbstractions(immediateAbstractions, trans)
	var abstractMap core.Element
	for _, abs := range immediateAbstractions {
		if abs.IsRefinementOfURI(CrlOneToOneMapURI, trans) {
			abstractMap = abs
			break
		}
	}
	if abstractMap == nil {
		return nil
	}
	// Validate that the abstraction has a sourceRef and that the sourceRef is referencing an element
	absSourceRef := abstractMap.GetFirstOwnedReferenceRefinedFromURI(CrlOneToOneMapSourceReferenceURI, trans)
	if absSourceRef == nil {
		return nil
	}
	absSource := absSourceRef.GetReferencedConcept(trans)
	if absSource == nil {
		return nil
	}
	// Validate that the abstraction has a targetRef and that the targetRef is referencing an element
	absTargetRef := abstractMap.GetFirstOwnedReferenceRefinedFromURI(CrlOneToOneMapTargetReferenceURI, trans)
	if absTargetRef == nil {
		return nil
	}
	absTarget := absTargetRef.GetReferencedConcept(trans)
	if absTarget == nil {
		return nil
	}
	// Check to see whether the source reference exists and references an element of the correct type
	sourceRef := mapInstance.GetFirstOwnedReferenceRefinedFrom(absSourceRef, trans)
	if sourceRef == nil {
		return nil
	}
	source := sourceRef.GetReferencedConcept(trans)
	if source == nil || !source.IsRefinementOf(absSource, trans) {
		return nil
	}

	// Now explore the targetRef
	targetRef := mapInstance.GetFirstOwnedReferenceRefinedFrom(absTargetRef, trans)
	// If the target ref does not exist, create it
	if targetRef == nil {
		targetRef, _ = uOfD.NewReference(trans)
		targetRef.SetOwningConcept(mapInstance, trans)
		targetRefRefinement, _ := uOfD.NewRefinement(trans)
		targetRefRefinement.SetOwningConcept(targetRef, trans)
		targetRefRefinement.SetAbstractConcept(absTargetRef, trans)
		targetRefRefinement.SetRefinedConcept(targetRef, trans)
	}

	// Now the target
	target := targetRef.GetReferencedConcept(trans)
	switch targetRef.GetReferencedAttributeName(trans) {
	case core.NoAttribute:
		if target == nil {
			// create it
			switch absTarget.(type) {
			case core.Literal:
				target, _ = uOfD.NewLiteral(trans)
			case core.Reference:
				target, _ = uOfD.NewReference(trans)
			case core.Refinement:
				target, _ = uOfD.NewRefinement(trans)
			case core.Element:
				target, _ = uOfD.NewElement(trans)
			}
			targetRefinement, _ := uOfD.NewRefinement(trans)
			targetRefinement.SetOwningConcept(target, trans)
			targetRefinement.SetAbstractConcept(absTarget, trans)
			targetRefinement.SetRefinedConcept(target, trans)
			target.SetLabel("Refinement of "+absTarget.GetLabel(trans)+"From"+source.GetLabel(trans), trans)
			targetRef.SetReferencedConcept(target, core.NoAttribute, trans)
		}
		if mapInstance.GetOwningConcept(trans) != nil {
			if isRootMap(mapInstance, trans) {
				err := target.SetOwningConcept(mapInstance.GetOwningConcept(trans), trans)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			} else {
				candidateTargetOwner := GetTarget(mapInstance.GetOwningConcept(trans), trans)
				if candidateTargetOwner != nil && target.GetOwningConceptID(trans) != candidateTargetOwner.GetConceptID(trans) {
					abstractTargetOwner := absTarget.GetOwningConcept(trans)
					if candidateTargetOwner.IsRefinementOf(abstractTargetOwner, trans) {
						err := target.SetOwningConcept(candidateTargetOwner, trans)
						if err != nil {
							return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
						}
					}
				}
			}
		}
	case core.OwningConceptID, core.LiteralValue, core.ReferencedConceptID, core.AbstractConceptID, core.RefinedConceptID, core.Label, core.Definition:
		// target = getParentMapTarget(mapInstance, trans)
		target = getAttributeTarget(mapInstance, trans)
	}
	if !target.IsRefinementOf(absTarget, trans) {
		return errors.New("In crlmaps.executeOneToOneMap, the found target is not refinement of abstraction target")
	}

	// Make value assignments as required
	// If the sourceRef is an attribute value reference, get the source value
	sourceAttributeName := sourceRef.GetReferencedAttributeName(trans)
	if sourceAttributeName != core.NoAttribute {
		var sourceAttributeValue string
		var sourceAttributeValueConcept core.Element
		switch sourceAttributeName {
		case core.ReferencedConceptID:
			switch sourceElement := source.(type) {
			case core.Reference:
				sourceAttributeValue = sourceElement.GetReferencedConceptID(trans)
				sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
			}
		case core.AbstractConceptID:
			switch sourceElement := source.(type) {
			case core.Refinement:
				sourceAttributeValue = sourceElement.GetAbstractConceptID(trans)
				sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
			}
		case core.RefinedConceptID:
			switch sourceElement := source.(type) {
			case core.Refinement:
				sourceAttributeValue = sourceElement.GetRefinedConceptID(trans)
				sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
			}
		case core.OwningConceptID:
			sourceAttributeValue = source.GetOwningConceptID(trans)
			sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
		case core.LiteralValue:
			sourceAttributeValue = source.(core.Literal).GetLiteralValue(trans)
		}

		// The found target is going to be the target of the parent map
		parentMap := mapInstance.GetOwningConcept(trans)
		if parentMap == nil {
			return errors.New("In crlmaps.executeOneToOneMap, parentMap was nil")
		}
		foundTarget := GetTarget(parentMap, trans)
		if foundTarget == nil {
			// mapping may not have been completed yet
			return nil
		}

		var foundTargetAttributeValue string
		switch source.(type) {
		case core.Literal:
			foundTargetAttributeValue = sourceAttributeValue
		default:
			sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
			if sourceAttributeValueConcept == nil {
				return errors.New("In crlmaps.executeOneToOneMap, the sourceAttributeValueConcept was not found")
			}
			sourceMap := SearchForMapForSource(mapInstance, sourceAttributeValueConcept, trans)
			if sourceMap != nil {
				pointerTargetConcept := GetTarget(sourceMap, trans)
				if pointerTargetConcept != nil {
					foundTargetAttributeValue = pointerTargetConcept.GetConceptID(trans)
				}
			}
		}

		referencedAttributeName := targetRef.GetReferencedAttributeName(trans)
		targetRef.SetReferencedConcept(foundTarget, referencedAttributeName, trans)

		switch referencedAttributeName {
		case core.NoAttribute:
			// Nothing to be done
		case core.ReferencedConceptID:
			// This case is valid only if the target is a reference
			switch targetElement := target.(type) {
			case core.Reference:
				err := targetElement.SetReferencedConceptID(foundTargetAttributeValue, core.NoAttribute, trans)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.AbstractConceptID:
			switch targetElement := target.(type) {
			case core.Refinement:
				err := targetElement.SetAbstractConceptID(foundTargetAttributeValue, trans)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.RefinedConceptID:
			// If the targetRef is an attribute value reference, set its value
			// If the target is a reference, set the referenced elementID
			switch targetElement := target.(type) {
			case core.Refinement:
				err := targetElement.SetRefinedConceptID(foundTargetAttributeValue, trans)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.LiteralValue:
			switch targetElement := target.(type) {
			case core.Literal:
				err := targetElement.SetLiteralValue(foundTargetAttributeValue, trans)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.OwningConceptID:
			err := target.SetOwningConceptID(foundTargetAttributeValue, trans)
			if err != nil {
				return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
			}
		}
	}

	// Now take care of map children.
	err := instantiateChildren(abstractMap, mapInstance, source, target, uOfD, trans)
	if err != nil {
		return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
	}
	return nil
}

func getAbstractMap(thisMap core.Element, trans *core.Transaction) core.Element {
	immediateAbstractions := map[string]core.Element{}
	thisMap.FindImmediateAbstractions(immediateAbstractions, trans)
	for _, abstraction := range immediateAbstractions {
		if abstraction.IsRefinementOfURI(CrlOneToOneMapURI, trans) {
			return abstraction
		}
	}
	return nil
}

func getAttributeTarget(attributeMap core.Element, trans *core.Transaction) core.Element {
	// Assumes that the parent map's target is either the desired target or an ancestor of the desired target
	// The target is either going to be the parent map's target or one of its descendants
	parentTarget := getParentMapTarget(attributeMap, trans)
	if parentTarget == nil {
		return nil
	}
	// Get the abstract attributeMap
	abstractMap := getAbstractMap(attributeMap, trans)
	// Get the abstract attributeMap's target
	abstractTarget := GetTarget(abstractMap, trans)
	// Find the descendent of the parent target that has the attributeMap's target as an ancestor
	if parentTarget.IsRefinementOf(abstractTarget, trans) {
		return parentTarget
	}
	childTarget := parentTarget.GetFirstOwnedConceptRefinedFrom(abstractTarget, trans)
	if childTarget != nil {
		return childTarget
	}
	return nil
}

// GetSource returns the source referenced by the given map
func GetSource(theMap core.Element, trans *core.Transaction) core.Element {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(trans)
}

// GetSourceReference returns the source reference for the given map
func GetSourceReference(theMap core.Element, trans *core.Transaction) core.Reference {
	return theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
}

// GetTarget returns the target referenced by the given map
func GetTarget(theMap core.Element, trans *core.Transaction) core.Element {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(trans)
}

// GetTargetReference returns the target reference for the given map
func GetTargetReference(theMap core.Element, trans *core.Transaction) core.Reference {
	return theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
}

func getRootMap(theMap core.Element, trans *core.Transaction) core.Element {
	owner := theMap.GetOwningConcept(trans)
	if owner != nil && isMap(owner, trans) {
		return getRootMap(owner, trans)
	}
	return theMap
}

func getParentMapTarget(theMap core.Element, trans *core.Transaction) core.Element {
	parentMap := theMap.GetOwningConcept(trans)
	ref := parentMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(trans)
}

// func getRootMapTarget(theMap core.Element, trans *core.Transaction) core.Element {
// 	rootMap := getRootMap(theMap, trans)
// 	ref := rootMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
// 	if ref == nil {
// 		return nil
// 	}
// 	return ref.GetReferencedConcept(trans)
// }

// FindAttributeMapForSource locates the attribute map referencing the given source, if any.
func FindAttributeMapForSource(currentMap core.Element, source core.Element, attributeName core.AttributeName, trans *core.Transaction) core.Element {
	if GetSource(currentMap, trans) == source && GetSourceReference(currentMap, trans).GetReferencedAttributeName(trans) == attributeName {
		return currentMap
	}
	for _, childMap := range currentMap.GetOwnedConceptsRefinedFromURI(CrlOneToOneMapURI, trans) {
		foundMap := FindAttributeMapForSource(childMap, source, attributeName, trans)
		if foundMap != nil {
			return foundMap
		}
	}
	return nil
}

// FindMapForSource locates the map corresponding to the given source, if any. It explores the current map and its descendants.
func FindMapForSource(currentMap core.Element, source core.Element, trans *core.Transaction) core.Element {
	if GetSource(currentMap, trans) == source && GetSourceReference(currentMap, trans).GetReferencedAttributeName(trans) == core.NoAttribute {
		return currentMap
	}
	for _, childMap := range currentMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, trans) {
		foundMap := FindMapForSource(childMap, source, trans)
		if foundMap != nil {
			return foundMap
		}
	}
	return nil
}

// SearchForMapForSource locates the map corresponding to the given source, if any. It starts with the current map.
// If not found, it then goes up one level to the parent map. It keeps going up until either a target is found or
// there is no parent. This method returns the first map found if there is more than one map.
func SearchForMapForSource(currentMap core.Element, source core.Element, trans *core.Transaction) core.Element {
	foundMap := FindMapForSource(currentMap, source, trans)
	if foundMap == nil {
		parentMap := currentMap.GetOwningConcept(trans)
		if parentMap != nil {
			foundMap = SearchForMapForSource(parentMap, source, trans)
		}
	}
	return foundMap
}

// FindMapForSourceAttribute locates the map corresponding to the given source attribute, if any.
func FindMapForSourceAttribute(currentMap core.Element, source core.Element, attributeName core.AttributeName, trans *core.Transaction) core.Element {
	if GetSource(currentMap, trans) == source && GetSourceReference(currentMap, trans).GetReferencedAttributeName(trans) == attributeName {
		return currentMap
	}
	for _, childMap := range currentMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, trans) {
		foundMap := FindMapForSourceAttribute(childMap, source, attributeName, trans)
		if foundMap != nil {
			return foundMap
		}
	}
	return nil
}

// FindTargetForSource locates the map corresponding to the given source (if any) and then returns its target
func FindTargetForSource(currentMap core.Element, source core.Element, trans *core.Transaction) core.Element {
	// get root map
	rootMap := getRootMap(currentMap, trans)
	// search the root map for a mapping whose source is the given source. If found, return target of the map
	foundMap := FindMapForSource(rootMap, source, trans)
	if foundMap != nil {
		return GetTarget(foundMap, trans)
	}
	return nil
}

func instantiateChildren(abstractMap core.Element, parentMap core.Element, source core.Element, target core.Element, uOfD *core.UniverseOfDiscourse, trans *core.Transaction) error {
	// for each of the abstractMap's children that is a map
	for _, abstractChildMap := range abstractMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, trans) {
		abstractChildMapSource := GetSource(abstractChildMap, trans)
		if abstractChildMapSource != nil {
			// There are two cases here, depending upon whether the source reference is to a pointer or an element.
			abstractChildMapSourceReference := abstractChildMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
			if abstractChildMapSourceReference == nil {
				return errors.New("In CrlMaps.go instantiateChildren, the abstractChildMapSource does not have a ChildMapSourceReference")
			}
			abstractChildMapSourceReferenceAttributeName := abstractChildMapSourceReference.GetReferencedAttributeName(trans)
			if abstractChildMapSourceReferenceAttributeName != core.NoAttribute {
				// If the abstractChildMap's source reference is to a pointer, then the actual source for the child is going to be
				// the parent's source. Error checking is required to ensure that the parent's source is of the appropriate type for the AttributeName
				// on the reference. In this case there will only be one instance of the abstractChildMap created.
				// Check to see whether there is already a map instance for this source
				parentMapSource := GetSource(parentMap, trans)
				if parentMapSource == nil {
					// This may not be an error - it may be a deletion that is being processed
					return nil
				}
				// We must find the Element whose attribute is being referenced. Two known cases are possible here (there may be others yet to be encountered).
				// Case 1: the parent's map source is that Element
				// Case 2: The sought-after Element is a Reference that is owned by the parent's map source.
				//         This latter case only occurs when the attribute name is ReferencedConceptID.
				// We first check to see whether the parent's map source is a refinement of the abstractChildMapSource. This condition may fail during editing scenarios
				// in which the owner of the childMap has not yet been assigned to the correct owner. This is not an error - it is an expected condition
				// If it is not, we then perform a secondary check to see whether the parent's map source has a child reference that is a refinement
				// of the abstractChildMapSource.

				// BUG There is a flaw in the following logic. The logic seems to assume that the reference to the owner pointer is a reference to the concept that
				// owns the pointer, while the reality is that the referenced conecpt is the concept to which the pointer refers. The assumption appears to be
				// the correct logic, so the reference to the pointer needs to be fixed.

				foundChildSource := parentMapSource // assume it's going to be the parent map source
				if !parentMapSource.IsRefinementOf(abstractChildMapSource, trans) {
					if abstractChildMapSourceReference.GetReferencedAttributeName(trans) == core.ReferencedConceptID {
						foundChildSource = parentMapSource.GetFirstOwnedReferenceRefinedFrom(abstractChildMapSource, trans)
						if foundChildSource == nil {
							return nil
						}
					} else {
						return nil
					}
				}
				var newMapInstance core.Element
				for _, mapInstance := range parentMap.GetOwnedConceptsRefinedFrom(abstractChildMap, trans) {
					mapInstanceSource := GetSource(mapInstance, trans)
					if mapInstanceSource == nil || mapInstanceSource.GetConceptID(trans) == foundChildSource.GetConceptID(trans) {
						newMapInstance = mapInstance
						break
					}
				}
				if newMapInstance == nil {
					newMapInstance, err := uOfD.CreateReplicateAsRefinement(abstractChildMap, trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateChildren failed")
					}
					err = newMapInstance.SetOwningConcept(parentMap, trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateChildren failed")
					}
				}
				newSourceRef := newMapInstance.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
				if newSourceRef == nil {
					return errors.New("In crlmaps.instantiateChildren, newSourceRef is nil")
				}
				if foundChildSource != nil && newSourceRef.GetReferencedConceptID(trans) != foundChildSource.GetConceptID(trans) {
					err := newSourceRef.SetReferencedConcept(foundChildSource, abstractChildMapSourceReferenceAttributeName, trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateChildren failed")
					}
				}
			} else {
				// The abstractChildMapSourceReference is to an element. There is an assumption in this case that the parent's source
				// contains elements that are refinements of the abstractChildMapSource. For each element that is a refinement of the
				// abstractChildMapSource found in the parent's source, instantiate the abstractChildMap (replicate as refinement)
				// and wire up the element as the source
				for _, sourceEl := range source.GetOwnedDescendantsRefinedFrom(abstractChildMapSource, trans) {
					// Check to see whether there is already a map instance for this source
					var newMapInstance core.Element
					for _, mapInstance := range parentMap.GetOwnedConceptsRefinedFrom(abstractChildMap, trans) {
						mapInstanceSource := GetSource(mapInstance, trans)
						if mapInstanceSource == nil || mapInstanceSource.GetConceptID(trans) == sourceEl.GetConceptID(trans) {
							newMapInstance = mapInstance
							break
						}
					}
					if newMapInstance == nil {
						var err error
						newMapInstance, err = uOfD.CreateReplicateAsRefinement(abstractChildMap, trans)
						if err != nil {
							return errors.Wrap(err, "crlmaps.instantiateChildren failed")
						}
						err = newMapInstance.SetOwningConcept(parentMap, trans)
						if err != nil {
							return errors.Wrap(err, "crlmaps.instantiateChildren failed")
						}
					}
					newSourceRef := newMapInstance.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
					if newSourceRef == nil {
						return errors.New("In crlmaps.instantiateChildren, newSourceRef is nil")
					}
					err := newSourceRef.SetReferencedConcept(sourceEl, core.NoAttribute, trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateChildren failed")
					}
				}
			}
		}
	}
	return nil
}

func isMap(candidate core.Element, trans *core.Transaction) bool {
	return candidate != nil && candidate.IsRefinementOfURI(CrlMapURI, trans)
}

func isRootMap(candidate core.Element, trans *core.Transaction) bool {
	rootMap := getRootMap(candidate, trans)
	return rootMap == candidate
}

// SetSource sets the source referenced by the given map
func SetSource(theMap core.Element, newSource core.Element, attributeName core.AttributeName, trans *core.Transaction) error {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
	if ref == nil {
		return errors.New("CrlMaps.SetSource called with map that does not have a source reference")
	}
	return ref.SetReferencedConcept(newSource, attributeName, trans)
}

// // SetSourceAttributeName sets the source attribute name referenced by the given map
// func SetSourceAttributeName(theMap core.Element, attributeName core.AttributeName, trans *core.HeldLocks) error {
// 	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
// 	if ref == nil {
// 		return errors.New("CrlMaps.SetSourceAttributeName called with map that does not have a source reference")
// 	}
// 	return ref.SetReferencedAttributeName(attributeName, trans)
// }

// SetTarget sets the target referenced by the given map
func SetTarget(theMap core.Element, newTarget core.Element, attributeName core.AttributeName, trans *core.Transaction) error {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
	if ref == nil {
		return errors.New("CrlMaps.SetTarget called with map that does not have a target reference")
	}
	return ref.SetReferencedConcept(newTarget, attributeName, trans)
}

// // SetTargetAttributeName sets the target attribute name referenced by the given map
// func SetTargetAttributeName(theMap core.Element, attributeName core.AttributeName, trans *core.HeldLocks) error {
// 	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
// 	if ref == nil {
// 		return errors.New("CrlMaps.SetTargetAttributeName called with map that does not have a target reference")
// 	}
// 	return ref.SetReferencedAttributeName(attributeName, trans)
// }
