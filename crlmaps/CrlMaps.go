package crlmaps

import (
	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlMapsDomainURI is the URI for the concpet space that defines the CRL Data Types
var CrlMapsDomainURI = "http://activeCRL.com/crlmaps/CrlMaps"

// CrlMapURI is the uri of the ancestor of all maps
var CrlMapURI = CrlMapsDomainURI + "/Map"

// CrlMapSourceURI is the URI of the CrlMap source
var CrlMapSourceURI = CrlMapURI + "/Source"

// CrlMapTargetURI is the URI of the CrlMap target
var CrlMapTargetURI = CrlMapURI + "/Target"

// CrlElementToElementMapURI is the URI for the Element-to-Element Map
var CrlElementToElementMapURI = CrlMapsDomainURI + "/ElementToElementMap"

// CrlElementToElementMapRefinementURI is the URI of the refinement showing it to be a refinement of CrlMap
var CrlElementToElementMapRefinementURI = CrlElementToElementMapURI + "/Refinement"

// CrlElementToElementMapSourceReferenceURI is the URI for the source reference
var CrlElementToElementMapSourceReferenceURI = CrlElementToElementMapURI + "/SourceReference"

// CrlElementToElementMapSourceReferenceRefinementURI is the URI of the refinement showing it to be a refinement of CrlMapSource
var CrlElementToElementMapSourceReferenceRefinementURI = CrlElementToElementMapSourceReferenceURI + "/Refinement"

// CrlElementToElementMapTargetReferenceURI is the URI for the target reference
var CrlElementToElementMapTargetReferenceURI = CrlElementToElementMapURI + "/TargetReference"

// CrlElementToElementMapTargetReferenceRefinementURI is the URI of the refinement showing it to be a refinement of CrlMapTarget
var CrlElementToElementMapTargetReferenceRefinementURI = CrlElementToElementMapTargetReferenceURI + "/Refinement"

// BuildCrlMapsDomain constructs the domain for CRL maps
func BuildCrlMapsDomain(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) error {
	crlMaps, err1 := uOfD.NewOwnedElement(nil, "CrlMapsDomain", hl, CrlMapsDomainURI)
	if err1 != nil {
		return errors.Wrap(err1, "crlmaps.BuildCrlMapsDomain failed")
	}

	// CrlMap
	crlMap, err2 := uOfD.NewOwnedElement(crlMaps, "CrlMap", hl, CrlMapURI)
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

	// CrlElementToElementMap
	crlElementToElementMap, err5 := uOfD.NewOwnedElement(crlMaps, "CrlElementToElementMap", hl, CrlElementToElementMapURI)
	if err5 != nil {
		return errors.Wrap(err5, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err6 := uOfD.NewOwnedRefinement(crlMap, crlElementToElementMap, "Refines CrlMap", hl, CrlElementToElementMapRefinementURI)
	if err6 != nil {
		return errors.Wrap(err6, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlElementToElementMapSourceReference, err7 := uOfD.NewOwnedReference(crlElementToElementMap, "SourceReference", hl, CrlElementToElementMapSourceReferenceURI)
	if err7 != nil {
		return errors.Wrap(err7, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err8 := uOfD.NewOwnedRefinement(crlMapSource, crlElementToElementMapSourceReference, "Refines CrlMapSource", hl, CrlElementToElementMapSourceReferenceRefinementURI)
	if err8 != nil {
		return errors.Wrap(err8, "crlmaps.BuildCrlMapsDomain failed")
	}
	crlElementToElementMapTargetReference, err9 := uOfD.NewOwnedReference(crlElementToElementMap, "TargetReference", hl, CrlElementToElementMapTargetReferenceURI)
	if err9 != nil {
		return errors.Wrap(err9, "crlmaps.BuildCrlMapsDomain failed")
	}
	_, err10 := uOfD.NewOwnedRefinement(crlMapTarget, crlElementToElementMapTargetReference, "Refines CrlMapTarget", hl, CrlElementToElementMapTargetReferenceRefinementURI)
	if err10 != nil {
		return errors.Wrap(err10, "crlmaps.BuildCrlMapsDomain failed")
	}

	err := crlMaps.SetReadOnlyRecursively(true, hl)
	if err != nil {
		return errors.Wrap(err, "CrlMaps BuildCrlMapsDomain failed")
	}
	err = crlMaps.SetIsCoreRecursively(hl)
	if err != nil {
		return errors.Wrap(err, "CrlMaps BuildCrlMapsDomain failed")
	}

	uOfD.AddFunction(CrlElementToElementMapURI, executeElementToElementMap)
	return nil
}

// executeElementToElementMap performs the mapping function
func executeElementToElementMap(mapInstance core.Element, notification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) error {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.WriteLockElement(mapInstance)
	// As an initial assumption, it probably doesn't matter what kind of notification has been received.
	// Validate that this instance is a refinement of an element that is, in turn, a refinement of CrlElementToElementMap
	var immediateAbstractions = map[string]core.Element{}
	mapInstance.FindImmediateAbstractions(immediateAbstractions, hl)
	var abstractionMap core.Element
	for _, abs := range immediateAbstractions {
		if abs.IsRefinementOfURI(CrlElementToElementMapURI, hl) {
			abstractionMap = abs
			break
		}
	}
	if abstractionMap == nil {
		return nil
	}
	// Validate that the abstraction has a sourceRef and that the sourceRef is referencing an element
	absSourceRef := abstractionMap.GetFirstOwnedReferenceRefinedFromURI(CrlElementToElementMapSourceReferenceURI, hl)
	if absSourceRef == nil {
		return nil
	}
	absSource := absSourceRef.GetReferencedConcept(hl)
	if absSource == nil {
		return nil
	}
	// Validate that the abstraction has a targetRef and that the targetRef is referencing an element
	absTargetRef := abstractionMap.GetFirstOwnedReferenceRefinedFromURI(CrlElementToElementMapTargetReferenceURI, hl)
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
		return errors.New("In crlmaps.executeElementToElementMap, the found target is not refinement of abstraction target")
	}
	if getRootMap(mapInstance, hl) != mapInstance {
		err := target.SetOwningConcept(getRootMapTarget(mapInstance, hl), hl)
		if err != nil {
			return errors.Wrap(err, "crlmaps.executeElementToElementMap failed")
		}
	}

	// Now take care of map children.
	err := instantiateChildren(abstractionMap, mapInstance, source, target, uOfD, hl)
	if err != nil {
		return errors.Wrap(err, "crlmaps.executeElementToElementMap failed")
	}
	hl.ReleaseLocksAndWait()
	return nil
}

func instantiateChildren(abstractionMap core.Element, parentMap core.Element, source core.Element, target core.Element, uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) error {
	// for each of the abstractionMap's children that is a map
	for _, childMap := range abstractionMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, hl) {
		childMapSource := getSource(childMap, hl)
		if childMapSource != nil {
			// if the source contains an instance of the child's source (need to generalize this for multiple sources in the future)
			// then instantiate the map (replicate as refinement) and wire up the source
			for _, sourceEl := range source.GetOwnedConceptsRefinedFrom(childMapSource, hl) {
				// Check to see whether there is already a map instance for this source
				var newMapInstance core.Element
				for _, mapInstance := range parentMap.GetOwnedConceptsRefinedFrom(childMap, hl) {
					mapInstanceSource := getSource(mapInstance, hl)
					if mapInstanceSource.GetConceptID(hl) == sourceEl.GetConceptID(hl) {
						newMapInstance = mapInstance
						break
					}
				}
				if newMapInstance == nil {
					newMapInstance, err := uOfD.CreateReplicateAsRefinement(childMap, hl)
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
	return nil
}

func getSource(theMap core.Element, hl *core.HeldLocks) core.Element {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, hl)
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
