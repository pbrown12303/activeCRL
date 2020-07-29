package crlmaps

import (
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

	// Element To Element Map
	crlOneToOneMap, err5 := uOfD.NewOwnedElement(crlMapsDomain, "CrlOneToOneMap", hl, CrlOneToOneMapURI)
	if err5 != nil {
		return errors.Wrap(err5, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err6 := uOfD.NewOwnedRefinement(crlMap, crlOneToOneMap, "Refines CrlMap", hl, CrlOneToOneMapRefinementURI)
	if err6 != nil {
		return errors.Wrap(err6, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlOneToOneMapSourceReference, err7 := uOfD.NewOwnedReference(crlOneToOneMap, "SourceReference", hl, CrlOneToOneMapSourceReferenceURI)
	if err7 != nil {
		return errors.Wrap(err7, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err8 := uOfD.NewOwnedRefinement(crlMapSource, crlOneToOneMapSourceReference, "Refines CrlMapSource", hl, CrlOneToOneMapSourceReferenceRefinementURI)
	if err8 != nil {
		return errors.Wrap(err8, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlOneToOneMapTargetReference, err9 := uOfD.NewOwnedReference(crlOneToOneMap, "TargetReference", hl, CrlOneToOneMapTargetReferenceURI)
	if err9 != nil {
		return errors.Wrap(err9, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err10 := uOfD.NewOwnedRefinement(crlMapTarget, crlOneToOneMapTargetReference, "Refines CrlMapTarget", hl, CrlOneToOneMapTargetReferenceRefinementURI)
	if err10 != nil {
		return errors.Wrap(err10, "crlmaps.BuildCrlMapsDomain failed")
	}

	// Reference To Element Map
	crlReferenceToElementMap, err11 := uOfD.NewOwnedElement(crlMapsDomain, "ReferenceToElementMap", hl, CrlReferenceToElementMapURI)
	if err11 != nil {
		return errors.Wrap(err11, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err12 := uOfD.NewOwnedRefinement(crlMap, crlReferenceToElementMap, "Refinement", hl, CrlReferenceToElementMapRefinementURI)
	if err12 != nil {
		return errors.Wrap(err12, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlReferenceToElementMapSource, err13 := uOfD.NewOwnedReference(crlReferenceToElementMap, "Source", hl, CrlReferenceToElementMapSourceURI)
	if err13 != nil {
		return errors.Wrap(err13, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err14 := uOfD.NewOwnedRefinement(crlMapSource, crlReferenceToElementMapSource, "Refinement", hl, CrlReferenceToElementMapSourceRefinementURI)
	if err14 != nil {
		return errors.Wrap(err14, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlReferenceToElementMapTarget, err15 := uOfD.NewOwnedReference(crlReferenceToElementMap, "Target", hl, CrlReferenceToElementMapTargetURI)
	if err15 != nil {
		return errors.Wrap(err15, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err16 := uOfD.NewOwnedRefinement(crlMapTarget, crlReferenceToElementMapTarget, "Refinement", hl, CrlReferenceToElementMapTargetRefinementURI)
	if err16 != nil {
		return errors.Wrap(err16, "crlmaps.BuildCrlMapsDomain failed")
	}

	// ID to Reference Map
	crlIDToReferenceMap, err17 := uOfD.NewOwnedElement(crlMapsDomain, "IDToReferenceMap", hl, CrlIDToReferenceMapURI)
	if err17 != nil {
		return errors.Wrap(err17, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err18 := uOfD.NewOwnedRefinement(crlMap, crlIDToReferenceMap, "Refinement", hl, CrlIDToReferenceMapRefinementURI)
	if err18 != nil {
		return errors.Wrap(err18, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlIDToReferenceMapSource, err19 := uOfD.NewOwnedReference(crlIDToReferenceMap, "Source", hl, CrlIDToReferenceMapSourceURI)
	if err19 != nil {
		return errors.Wrap(err19, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err20 := uOfD.NewOwnedRefinement(crlMapSource, crlIDToReferenceMapSource, "Source", hl, CrlIDToReferenceMapSourceRefinementURI)
	if err20 != nil {
		return errors.Wrap(err20, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlIDToReferenceMapTarget, err21 := uOfD.NewOwnedReference(crlIDToReferenceMap, "Target", hl, CrlIDToReferenceMapTargetURI)
	if err21 != nil {
		return errors.Wrap(err21, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err22 := uOfD.NewOwnedRefinement(crlMapTarget, crlIDToReferenceMapTarget, "Refinement", hl, CrlIDToReferenceMapTargetRefinementURI)
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
		target.SetLabel(absTarget.GetLabel(hl)+"From"+source.GetLabel(hl), hl)
		targetRef.SetReferencedConcept(target, hl)
	}
	if !target.IsRefinementOf(absTarget, hl) {
		return errors.New("In crlmaps.executeOneToOneMap, the found target is not refinement of abstraction target")
	}
	if mapInstance.GetOwningConcept(hl) != nil {
		targetOwner := getTarget(mapInstance.GetOwningConcept(hl), hl)
		if targetOwner != nil {
			err := target.SetOwningConcept(targetOwner, hl)
			if err != nil {
				return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
			}
		}
	}

	// Make value assignments as required
	// If the sourceRef is an attribute value reference, get the source value
	sourceAttributeName := sourceRef.GetReferencedAttributeName(hl)
	if sourceAttributeName != core.NoAttribute {
		sourceValue := sourceRef.GetReferencedAttributeValue(hl)

		switch targetRef.GetReferencedAttributeName(hl) {
		case core.NoAttribute:
			// This case is valid only if the target is a reference, in which case we are setting the reference's referencedConceptID
			switch target.(type) {
			case core.Reference:
				err := target.(core.Reference).SetReferencedConceptID(sourceValue, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.ReferencedConceptID:
			// This case is valid only if the target is a reference
			switch target.(type) {
			case core.Reference:
				err := target.(core.Reference).SetReferencedConceptID(sourceValue, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.AbstractConceptID:
			switch target.(type) {
			case core.Refinement:
				err := target.(core.Refinement).SetAbstractConceptID(sourceValue, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.RefinedConceptID:
			// If the targetRef is an attribute value reference, set its value
			// If the target is a reference, set the referenced elementID
			switch target.(type) {
			case core.Refinement:
				err := target.(core.Refinement).SetRefinedConceptID(sourceValue, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.LiteralValue:
			switch target.(type) {
			case core.Literal:
				err := target.(core.Literal).SetLiteralValue(sourceValue, hl)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
			}
		case core.OwningConceptID:
			err := target.SetOwningConceptID(sourceValue, hl)
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

func instantiateChildren(abstractMap core.Element, parentMap core.Element, source core.Element, target core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) error {
	// for each of the abstractMap's children that is a map
	for _, abstractChildMap := range abstractMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, hl) {
		abstractChildMapSource := getSource(abstractChildMap, hl)
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
				parentMapSource := getSource(parentMap, hl)
				if parentMapSource == nil {
					return errors.New("In CrlMaps.go instantiateChildren, the parentMap does not have a parentMapSource")
				}
				var newMapInstance core.Element
				for _, mapInstance := range parentMap.GetOwnedConceptsRefinedFrom(abstractChildMap, hl) {
					mapInstanceSource := getSource(mapInstance, hl)
					if mapInstanceSource.GetConceptID(hl) == parentMapSource.GetConceptID(hl) {
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
					newSourceRef := newMapInstance.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
					if newSourceRef == nil {
						return errors.New("In crlmaps.instantiateChildren, newSourceRef is nil")
					}
					err = newSourceRef.SetReferencedConcept(parentMapSource, hl)
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
						mapInstanceSource := getSource(mapInstance, hl)
						if mapInstanceSource.GetConceptID(hl) == sourceEl.GetConceptID(hl) {
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
						newSourceRef := newMapInstance.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
						if newSourceRef == nil {
							return errors.New("In crlmaps.instantiateChildren, newSourceRef is nil")
						}
						err = newSourceRef.SetReferencedConcept(sourceEl, hl)
						if err != nil {
							return errors.Wrap(err, "crlmaps.instantiateChildren failed")
						}
					}
				}
			}
		}
	}
	return nil
}

func getSource(theMap core.Element, hl *core.HeldLocks) core.Element {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(hl)
}

func getTarget(theMap core.Element, hl *core.HeldLocks) core.Element {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(hl)
}

func getRootMap(theMap core.Element, hl *core.HeldLocks) core.Element {
	owner := theMap.GetOwningConcept(hl)
	if owner != nil && isMap(owner, hl) {
		return getRootMap(owner, hl)
	}
	return theMap
}

func getRootMapTarget(theMap core.Element, hl *core.HeldLocks) core.Element {
	rootMap := getRootMap(theMap, hl)
	ref := rootMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, hl)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(hl)
}

func isMap(candidate core.Element, hl *core.HeldLocks) bool {
	return candidate != nil && candidate.IsRefinementOfURI(CrlMapURI, hl)
}
