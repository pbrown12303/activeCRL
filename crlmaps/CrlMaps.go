package crlmaps

import (
	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlMapsDomainURI is the URI for the concpet space that defines the CRL Data Types
var CrlMapsDomainURI = "http://activeCRL.com/crlmaps/CrlMaps"

// CrlElementToElementMapURI is the URI for the Element-to-Element Map
var CrlElementToElementMapURI = CrlMapsDomainURI + "/ElementToElementMap"

// CrlElementToElementMapSourceReferenceURI is the URI for the source reference
var CrlElementToElementMapSourceReferenceURI = CrlElementToElementMapURI + "/SourceReference"

// CrlElementToElementMapTargetReferenceURI is the URI for the target reference
var CrlElementToElementMapTargetReferenceURI = CrlElementToElementMapURI + "/TargetReference"

// BuildCrlMapsDomain constructs the domain for CRL maps
func BuildCrlMapsDomain(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) error {
	crlMaps, _ := uOfD.NewElement(hl, CrlMapsDomainURI)
	crlMaps.SetLabel("CrlMaps", hl)

	crlElementToElementMap, _ := uOfD.NewElement(hl, CrlElementToElementMapURI)
	crlElementToElementMap.SetLabel("CrlElementToElementMap", hl)
	crlElementToElementMap.SetOwningConcept(crlMaps, hl)

	crlElementToElementMapSourceReference, _ := uOfD.NewReference(hl, CrlElementToElementMapSourceReferenceURI)
	crlElementToElementMapSourceReference.SetLabel("SourceReference", hl)
	crlElementToElementMapSourceReference.SetOwningConcept(crlElementToElementMap, hl)

	crlElementToElementMapTargetReference, _ := uOfD.NewReference(hl, CrlElementToElementMapTargetReferenceURI)
	crlElementToElementMapTargetReference.SetLabel("TargetReference", hl)
	crlElementToElementMapTargetReference.SetOwningConcept(crlElementToElementMap, hl)

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
func executeElementToElementMap(mapInstance core.Element, notification *core.ChangeNotification, uOfD *core.UniverseOfDiscourse) {
	hl := uOfD.NewHeldLocks()
	defer hl.ReleaseLocksAndWait()
	hl.WriteLockElement(mapInstance)
	// As an initial assumption, it probably doesn't matter what kind of notification has been received.
	// Validate that this instance is a refinement of an element that is, in turn, a refinement of CrlElementToElementMap
	var immediateAbstractions = map[string]core.Element{}
	mapInstance.FindImmediateAbstractions(immediateAbstractions, hl)
	var foundAbstraction core.Element
	for _, abs := range immediateAbstractions {
		if abs.IsRefinementOfURI(CrlElementToElementMapURI, hl) {
			foundAbstraction = abs
			break
		}
	}
	if foundAbstraction == nil {
		return
	}
	// Validate that the abstraction has a sourceRef and that the sourceRef is referencing an element
	absSourceRef := foundAbstraction.GetFirstOwnedReferenceRefinedFromURI(CrlElementToElementMapSourceReferenceURI, hl)
	if absSourceRef == nil {
		return
	}
	absSource := absSourceRef.GetReferencedConcept(hl)
	if absSource == nil {
		return
	}
	// Validate that the abstraction has a targetRef and that the targetRef is referencing an element
	absTargetRef := foundAbstraction.GetFirstOwnedReferenceRefinedFromURI(CrlElementToElementMapTargetReferenceURI, hl)
	if absTargetRef == nil {
		return
	}
	absTarget := absTargetRef.GetReferencedConcept(hl)
	if absTarget == nil {
		return
	}
	// Check to see whether the source reference exists and references an element of the correct type
	sourceRef := mapInstance.GetFirstOwnedReferenceRefinedFrom(absSourceRef, hl)
	if sourceRef == nil {
		return
	}
	source := sourceRef.GetReferencedConcept(hl)
	if source == nil || !source.IsRefinementOf(absSource, hl) {
		return
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
		// something is wrong, but error handling for these functions is not yet implemented
		return
	}

	// Now take care of map children.
}
