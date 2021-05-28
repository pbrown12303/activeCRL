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
func BuildCrlMapsDomain(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) error {
	crlMapsDomain, err1 := uOfD.NewOwnedElement(nil, "CrlMapsDomain", hl, CrlMapsDomainURI)
	if err1 != nil {
		return errors.Wrap(err1, "crlmaps.BuildCrlMapsDomain failed")
	}

	// CrlMap
	crlMap, err2 := uOfD.NewOwnedElement(crlMapsDomain, "CrlMap", hl, CrlMapURI)
	if err2 != nil {
		return errors.Wrap(err2, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlMapSource, err3 := uOfD.NewOwnedReference(crlMap, "MapSource", hl, CrlMapSourceURI)
	if err3 != nil {
		return errors.Wrap(err3, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlMapTarget, err4 := uOfD.NewOwnedReference(crlMap, "MapTarget", hl, CrlMapTargetURI)
	if err4 != nil {
		return errors.Wrap(err4, "crlmaps.BuildCrlMapsDomain failed")
	}

	// One to One Map
	crlOneToOneMap, err5 := uOfD.NewOwnedElement(crlMapsDomain, "CrlOneToOneMap", hl, CrlOneToOneMapURI)
	if err5 != nil {
		return errors.Wrap(err5, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err6 := uOfD.NewCompleteRefinement(crlMap, crlOneToOneMap, "Refines CrlMap", hl, CrlOneToOneMapRefinementURI)
	if err6 != nil {
		return errors.Wrap(err6, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlOneToOneMapSourceReference, err7 := uOfD.NewOwnedReference(crlOneToOneMap, "SourceReference", hl, CrlOneToOneMapSourceReferenceURI)
	if err7 != nil {
		return errors.Wrap(err7, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err8 := uOfD.NewCompleteRefinement(crlMapSource, crlOneToOneMapSourceReference, "Refines CrlMapSource", hl, CrlOneToOneMapSourceReferenceRefinementURI)
	if err8 != nil {
		return errors.Wrap(err8, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlOneToOneMapTargetReference, err9 := uOfD.NewOwnedReference(crlOneToOneMap, "TargetReference", hl, CrlOneToOneMapTargetReferenceURI)
	if err9 != nil {
		return errors.Wrap(err9, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err10 := uOfD.NewCompleteRefinement(crlMapTarget, crlOneToOneMapTargetReference, "Refines CrlMapTarget", hl, CrlOneToOneMapTargetReferenceRefinementURI)
	if err10 != nil {
		return errors.Wrap(err10, "crlmaps.BuildCrlMapsDomain failed")
	}

	// Reference To Element Map
	crlReferenceToElementMap, err11 := uOfD.NewOwnedElement(crlMapsDomain, "ReferenceToElementMap", hl, CrlReferenceToElementMapURI)
	if err11 != nil {
		return errors.Wrap(err11, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err12 := uOfD.NewCompleteRefinement(crlMap, crlReferenceToElementMap, "Refinement", hl, CrlReferenceToElementMapRefinementURI)
	if err12 != nil {
		return errors.Wrap(err12, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlReferenceToElementMapSource, err13 := uOfD.NewOwnedReference(crlReferenceToElementMap, "Source", hl, CrlReferenceToElementMapSourceURI)
	if err13 != nil {
		return errors.Wrap(err13, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err14 := uOfD.NewCompleteRefinement(crlMapSource, crlReferenceToElementMapSource, "Refinement", hl, CrlReferenceToElementMapSourceRefinementURI)
	if err14 != nil {
		return errors.Wrap(err14, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlReferenceToElementMapTarget, err15 := uOfD.NewOwnedReference(crlReferenceToElementMap, "Target", hl, CrlReferenceToElementMapTargetURI)
	if err15 != nil {
		return errors.Wrap(err15, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err16 := uOfD.NewCompleteRefinement(crlMapTarget, crlReferenceToElementMapTarget, "Refinement", hl, CrlReferenceToElementMapTargetRefinementURI)
	if err16 != nil {
		return errors.Wrap(err16, "crlmaps.BuildCrlMapsDomain failed")
	}

	// ID to Reference Map
	crlIDToReferenceMap, err17 := uOfD.NewOwnedElement(crlMapsDomain, "IDToReferenceMap", hl, CrlIDToReferenceMapURI)
	if err17 != nil {
		return errors.Wrap(err17, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err18 := uOfD.NewCompleteRefinement(crlMap, crlIDToReferenceMap, "Refinement", hl, CrlIDToReferenceMapRefinementURI)
	if err18 != nil {
		return errors.Wrap(err18, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlIDToReferenceMapSource, err19 := uOfD.NewOwnedReference(crlIDToReferenceMap, "Source", hl, CrlIDToReferenceMapSourceURI)
	if err19 != nil {
		return errors.Wrap(err19, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err20 := uOfD.NewCompleteRefinement(crlMapSource, crlIDToReferenceMapSource, "Source", hl, CrlIDToReferenceMapSourceRefinementURI)
	if err20 != nil {
		return errors.Wrap(err20, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlIDToReferenceMapTarget, err21 := uOfD.NewOwnedReference(crlIDToReferenceMap, "Target", hl, CrlIDToReferenceMapTargetURI)
	if err21 != nil {
		return errors.Wrap(err21, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err22 := uOfD.NewCompleteRefinement(crlMapTarget, crlIDToReferenceMapTarget, "Refinement", hl, CrlIDToReferenceMapTargetRefinementURI)
	if err22 != nil {
		return errors.Wrap(err22, "crlmaps.BuildCrlMapsDomain failed")
	}

	err := crlMapsDomain.SetReadOnlyRecursively(true, hl)
	if err != nil {
		return errors.Wrap(err, "CrlMaps BuildCrlMapsDomain failed")
	}
	err = crlMapsDomain.SetIsCoreRecursively(hl)
	if err != nil {
		return errors.Wrap(err, "CrlMaps BuildCrlMapsDomain failed")
	}

	uOfD.AddFunction(CrlOneToOneMapURI, executeOneToOneMap)
	return nil
}

// executeOneToOneMap performs the mapping function
func executeOneToOneMap(mapInstance core.Element, notification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.WriteLockElement(mapInstance)

	log.Printf("Executing executeOneToOneMap for map labeled %s", mapInstance.GetLabel(hl))

	// As an initial assumption, it probably doesn't matter what kind of notification has been received.
	// Validate that this instance is a refinement of an element that is, in turn, a refinement of CrlOneToOneMap
	var immediateAbstractions = map[string]core.Element{}
	mapInstance.FindImmediateAbstractions(immediateAbstractions, hl)
	var abstractMap core.Element
	for _, abs := range immediateAbstractions {
		if abs.IsRefinementOfURI(CrlOneToOneMapURI, hl) {
			abstractMap = abs
			break
		}
	}
	if abstractMap == nil {
		return nil
	}
	// Validate that the abstraction has a sourceRef and that the sourceRef is referencing an element
	absSourceRef := abstractMap.GetFirstOwnedReferenceRefinedFromURI(CrlOneToOneMapSourceReferenceURI, hl)
	if absSourceRef == nil {
		return nil
	}
	absSource := absSourceRef.GetReferencedConcept(hl)
	if absSource == nil {
		return nil
	}
	// Validate that the abstraction has a targetRef and that the targetRef is referencing an element
	absTargetRef := abstractMap.GetFirstOwnedReferenceRefinedFromURI(CrlOneToOneMapTargetReferenceURI, hl)
	if absTargetRef == nil {
		return nil
	}
	absTarget := absTargetRef.GetReferencedConcept(hl)
	if absTarget == nil {
		return nil
	}
	// Check to see whether the source reference exists and references an element of the correct type
	sourceRef := mapInstance.GetFirstOwnedReferenceRefinedFrom(absSourceRef, hl)
	if sourceRef == nil {
		return nil
	}
	source := sourceRef.GetReferencedConcept(hl)
	if source == nil || !source.IsRefinementOf(absSource, hl) {
		return nil
	}

	// Now explore the targetRef
	targetRef := mapInstance.GetFirstOwnedReferenceRefinedFrom(absTargetRef, hl)
	// If the target ref does not exist, create it
	if targetRef == nil {
		targetRef, _ = uOfD.NewReference(hl)
		targetRef.SetOwningConcept(mapInstance, hl)
		targetRefRefinement, _ := uOfD.NewRefinement(hl)
		targetRefRefinement.SetOwningConcept(targetRef, hl)
		targetRefRefinement.SetAbstractConcept(absTargetRef, hl)
		targetRefRefinement.SetRefinedConcept(targetRef, hl)
	}

	// Now the target
	target := targetRef.GetReferencedConcept(hl)
	switch targetRef.GetReferencedAttributeName(hl) {
	case core.NoAttribute:
		if target == nil {
			// create it
			switch absTarget.(type) {
			case core.Literal:
				target, _ = uOfD.NewLiteral(hl)
			case core.Reference:
				target, _ = uOfD.NewReference(hl)
			case core.Refinement:
				target, _ = uOfD.NewRefinement(hl)
			case core.Element:
				target, _ = uOfD.NewElement(hl)
			}
			targetRefinement, _ := uOfD.NewRefinement(hl)
			targetRefinement.SetOwningConcept(target, hl)
			targetRefinement.SetAbstractConcept(absTarget, hl)
			targetRefinement.SetRefinedConcept(target, hl)
			target.SetLabel("Refinement of "+absTarget.GetLabel(hl)+"From"+source.GetLabel(hl), hl)
			targetRef.SetReferencedConcept(target, core.NoAttribute, hl)
		}
		if mapInstance.GetOwningConcept(hl) != nil {
			if isRootMap(mapInstance, hl) {
				err := target.SetOwningConcept(mapInstance.GetOwningConcept(hl), hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			} else {
				candidateTargetOwner := GetTarget(mapInstance.GetOwningConcept(hl), hl)
				if candidateTargetOwner != nil && target.GetOwningConceptID(hl) != candidateTargetOwner.GetConceptID(hl) {
					abstractTargetOwner := absTarget.GetOwningConcept(hl)
					if candidateTargetOwner.IsRefinementOf(abstractTargetOwner, hl) {
						err := target.SetOwningConcept(candidateTargetOwner, hl)
						if err != nil {
							return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
						}
					}
				}
			}
		}
	case core.OwningConceptID, core.LiteralValue, core.ReferencedConceptID, core.AbstractConceptID, core.RefinedConceptID, core.Label, core.Definition:
		target = getParentMapTarget(mapInstance, hl)
	}
	if !target.IsRefinementOf(absTarget, hl) {
		return errors.New("In crlmaps.executeOneToOneMap, the found target is not refinement of abstraction target")
	}

	// Make value assignments as required
	// If the sourceRef is an attribute value reference, get the source value
	sourceAttributeName := sourceRef.GetReferencedAttributeName(hl)
	if sourceAttributeName != core.NoAttribute {
		var sourceAttributeValue string
		var sourceAttributeValueConcept core.Element
		switch sourceAttributeName {
		case core.ReferencedConceptID:
			switch sourceElement := source.(type) {
			case core.Reference:
				sourceAttributeValue = sourceElement.GetReferencedConceptID(hl)
				sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
			}
		case core.AbstractConceptID:
			switch sourceElement := source.(type) {
			case core.Refinement:
				sourceAttributeValue = sourceElement.GetAbstractConceptID(hl)
				sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
			}
		case core.RefinedConceptID:
			switch sourceElement := source.(type) {
			case core.Refinement:
				sourceAttributeValue = sourceElement.GetRefinedConceptID(hl)
				sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
			}
		case core.OwningConceptID:
			sourceAttributeValue = source.GetOwningConceptID(hl)
			sourceAttributeValueConcept = uOfD.GetElement(sourceAttributeValue)
		case core.LiteralValue:
			sourceAttributeValue = source.(core.Literal).GetLiteralValue(hl)
		}
		var foundTarget core.Element
		var foundTargetValue string
		switch source.(type) {
		case core.Literal:
			foundTargetValue = sourceAttributeValue
		default:
			foundTarget = FindTargetForSource(mapInstance, sourceAttributeValueConcept, hl)
			if foundTarget != nil {
				foundTargetValue = foundTarget.GetConceptID(hl)
			}
		}

		switch targetRef.GetReferencedAttributeName(hl) {
		case core.NoAttribute:
			targetRef.SetReferencedConcept(foundTarget, core.NoAttribute, hl)
		default:
			targetRef.SetReferencedConcept(target, core.NoAttribute, hl)
		}

		switch targetRef.GetReferencedAttributeName(hl) {
		case core.NoAttribute:
			// This case is valid only if the target is a reference, in which case we are setting the reference's referencedConceptID
			switch targetElement := target.(type) {
			case core.Reference:
				err := targetElement.SetReferencedConceptID(foundTargetValue, core.NoAttribute, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.ReferencedConceptID:
			// This case is valid only if the target is a reference
			switch targetElement := target.(type) {
			case core.Reference:
				err := targetElement.SetReferencedConceptID(foundTargetValue, core.NoAttribute, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.AbstractConceptID:
			switch targetElement := target.(type) {
			case core.Refinement:
				err := targetElement.SetAbstractConceptID(foundTargetValue, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.RefinedConceptID:
			// If the targetRef is an attribute value reference, set its value
			// If the target is a reference, set the referenced elementID
			switch targetElement := target.(type) {
			case core.Refinement:
				err := targetElement.SetRefinedConceptID(foundTargetValue, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.LiteralValue:
			switch targetElement := target.(type) {
			case core.Literal:
				err := targetElement.SetLiteralValue(foundTargetValue, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.OwningConceptID:
			err := target.SetOwningConceptID(foundTargetValue, hl)
			if err != nil {
				return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
			}
		}
	}

	// Now take care of map children.
	err := instantiateChildren(abstractMap, mapInstance, source, target, uOfD, hl)
	if err != nil {
		return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
	}
	hl.ReleaseLocksAndWait()
	return nil
}

// GetSource returns the source referenced by the given map
func GetSource(theMap core.Element, hl *core.HeldLocks) core.Element {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(hl)
}

// GetSourceReference returns the source reference for the given map
func GetSourceReference(theMap core.Element, hl *core.HeldLocks) core.Reference {
	return theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
}

// GetTarget returns the target referenced by the given map
func GetTarget(theMap core.Element, hl *core.HeldLocks) core.Element {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(hl)
}

// GetTargetReference returns the target reference for the given map
func GetTargetReference(theMap core.Element, hl *core.HeldLocks) core.Reference {
	return theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
}

func getRootMap(theMap core.Element, hl *core.HeldLocks) core.Element {
	owner := theMap.GetOwningConcept(hl)
	if owner != nil && isMap(owner, hl) {
		return getRootMap(owner, hl)
	}
	return theMap
}

func getParentMapTarget(theMap core.Element, hl *core.HeldLocks) core.Element {
	parentMap := theMap.GetOwningConcept(hl)
	ref := parentMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(hl)
}

func getRootMapTarget(theMap core.Element, hl *core.HeldLocks) core.Element {
	rootMap := getRootMap(theMap, hl)
	ref := rootMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(hl)
}

// FindMapForSource locates the map corresponding to the given source, if any.
func FindMapForSource(currentMap core.Element, source core.Element, hl *core.HeldLocks) core.Element {
	if GetSource(currentMap, hl) == source && GetSourceReference(currentMap, hl).GetReferencedAttributeName(hl) == core.NoAttribute {
		return currentMap
	}
	for _, childMap := range currentMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, hl) {
		foundMap := FindMapForSource(childMap, source, hl)
		if foundMap != nil {
			return foundMap
		}
	}
	return nil
}

// FindMapForSourceAttribute locates the map corresponding to the given source attribute, if any.
func FindMapForSourceAttribute(currentMap core.Element, source core.Element, attributeName core.AttributeName, hl *core.HeldLocks) core.Element {
	if GetSource(currentMap, hl) == source && GetSourceReference(currentMap, hl).GetReferencedAttributeName(hl) == attributeName {
		return currentMap
	}
	for _, childMap := range currentMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, hl) {
		foundMap := FindMapForSourceAttribute(childMap, source, attributeName, hl)
		if foundMap != nil {
			return foundMap
		}
	}
	return nil
}

// FindTargetForSource locates the map corresponding to the given source (if any) and then returns its target
func FindTargetForSource(currentMap core.Element, source core.Element, hl *core.HeldLocks) core.Element {
	// get root map
	rootMap := getRootMap(currentMap, hl)
	// search the root map for a mapping whose source is the given source. If found, return target of the map
	foundMap := FindMapForSource(rootMap, source, hl)
	if foundMap != nil {
		return GetTarget(foundMap, hl)
	}
	return nil
}

func instantiateChildren(abstractMap core.Element, parentMap core.Element, source core.Element, target core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) error {
	// for each of the abstractMap's children that is a map
	for _, abstractChildMap := range abstractMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, hl) {
		abstractChildMapSource := GetSource(abstractChildMap, hl)
		if abstractChildMapSource != nil {
			// There are two cases here, depending upon whether the source reference is to a pointer or an element.
			abstractChildMapSourceReference := abstractChildMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
			if abstractChildMapSourceReference == nil {
				return errors.New("In CrlMaps.go instantiateChildren, the abstractChildMapSource does not have a ChildMapSourceReference")
			}
			if abstractChildMapSourceReference.GetReferencedAttributeName(hl) != core.NoAttribute {
				// If the abstractChildMap's source reference is to a pointer, then the actual source for the child is going to be
				// the parent's source. Error checking is required to ensure that the parent's source is of the appropriate type for the AttributeName
				// on the reference. In this case there will only be one instance of the abstractChildMap created.
				// Check to see whether there is already a map instance for this source
				parentMapSource := GetSource(parentMap, hl)
				if parentMapSource == nil {
					// TODO This may not be an error - it may be a deletion that is being processed
					return errors.New("In CrlMaps.go instantiateChildren, the parentMap does not have a parentMapSource")
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
				if !parentMapSource.IsRefinementOf(abstractChildMapSource, hl) {
					if abstractChildMapSourceReference.GetReferencedAttributeName(hl) == core.ReferencedConceptID {
						foundChildSource = parentMapSource.GetFirstOwnedReferenceRefinedFrom(abstractChildMapSource, hl)
						if foundChildSource == nil {
							return nil
						}
					} else {
						return nil
					}
				}
				var newMapInstance core.Element
				for _, mapInstance := range parentMap.GetOwnedConceptsRefinedFrom(abstractChildMap, hl) {
					mapInstanceSource := GetSource(mapInstance, hl)
					if mapInstanceSource == nil || mapInstanceSource.GetConceptID(hl) == foundChildSource.GetConceptID(hl) {
						newMapInstance = mapInstance
						break
					}
				}
				if newMapInstance == nil {
					newMapInstance, err := uOfD.CreateReplicateAsRefinement(abstractChildMap, hl)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateChildren failed")
					}
					err = newMapInstance.SetOwningConcept(parentMap, hl)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateChildren failed")
					}
				}
				newSourceRef := newMapInstance.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
				if newSourceRef == nil {
					return errors.New("In crlmaps.instantiateChildren, newSourceRef is nil")
				}
				if foundChildSource != nil && newSourceRef.GetReferencedConceptID(hl) != foundChildSource.GetConceptID(hl) {
					err := newSourceRef.SetReferencedConcept(foundChildSource, core.NoAttribute, hl)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateChildren failed")
					}
				}
			} else {
				// The abstractChildMapSourceReference is to an element. There is an assumption in this case that the parent's source
				// contains elements that are refinements of the abstractChildMapSource. For each element that is a refinement of the
				// abstractChildMapSource found in the parent's source, instantiate the abstractChildMap (replicate as refinement)
				// and wire up the element as the source
				for _, sourceEl := range source.GetOwnedDescendantsRefinedFrom(abstractChildMapSource, hl) {
					// Check to see whether there is already a map instance for this source
					var newMapInstance core.Element
					for _, mapInstance := range parentMap.GetOwnedConceptsRefinedFrom(abstractChildMap, hl) {
						mapInstanceSource := GetSource(mapInstance, hl)
						if mapInstanceSource == nil || mapInstanceSource.GetConceptID(hl) == sourceEl.GetConceptID(hl) {
							newMapInstance = mapInstance
							break
						}
					}
					if newMapInstance == nil {
						var err error
						newMapInstance, err = uOfD.CreateReplicateAsRefinement(abstractChildMap, hl)
						if err != nil {
							return errors.Wrap(err, "crlmaps.instantiateChildren failed")
						}
						err = newMapInstance.SetOwningConcept(parentMap, hl)
						if err != nil {
							return errors.Wrap(err, "crlmaps.instantiateChildren failed")
						}
					}
					newSourceRef := newMapInstance.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
					if newSourceRef == nil {
						return errors.New("In crlmaps.instantiateChildren, newSourceRef is nil")
					}
					err := newSourceRef.SetReferencedConcept(sourceEl, core.NoAttribute, hl)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateChildren failed")
					}
				}
			}
		}
	}
	return nil
}

func isMap(candidate core.Element, hl *core.HeldLocks) bool {
	return candidate != nil && candidate.IsRefinementOfURI(CrlMapURI, hl)
}

func isRootMap(candidate core.Element, hl *core.HeldLocks) bool {
	rootMap := getRootMap(candidate, hl)
	return rootMap == candidate
}

// SetSource sets the source referenced by the given map
func SetSource(theMap core.Element, newSource core.Element, attributeName core.AttributeName, hl *core.HeldLocks) error {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
	if ref == nil {
		return errors.New("CrlMaps.SetSource called with map that does not have a source reference")
	}
	return ref.SetReferencedConcept(newSource, attributeName, hl)
}

// // SetSourceAttributeName sets the source attribute name referenced by the given map
// func SetSourceAttributeName(theMap core.Element, attributeName core.AttributeName, hl *core.HeldLocks) error {
// 	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
// 	if ref == nil {
// 		return errors.New("CrlMaps.SetSourceAttributeName called with map that does not have a source reference")
// 	}
// 	return ref.SetReferencedAttributeName(attributeName, hl)
// }

// SetTarget sets the target referenced by the given map
func SetTarget(theMap core.Element, newTarget core.Element, attributeName core.AttributeName, hl *core.HeldLocks) error {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
	if ref == nil {
		return errors.New("CrlMaps.SetTarget called with map that does not have a target reference")
	}
	return ref.SetReferencedConcept(newTarget, attributeName, hl)
}

// // SetTargetAttributeName sets the target attribute name referenced by the given map
// func SetTargetAttributeName(theMap core.Element, attributeName core.AttributeName, hl *core.HeldLocks) error {
// 	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
// 	if ref == nil {
// 		return errors.New("CrlMaps.SetTargetAttributeName called with map that does not have a target reference")
// 	}
// 	return ref.SetReferencedAttributeName(attributeName, hl)
// }
