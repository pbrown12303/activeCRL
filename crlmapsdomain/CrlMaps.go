package crlmapsdomain

import (
	// "log"

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
// var CrlOneToOneMapRefinementURI = CrlOneToOneMapURI + "/Refinement"

// CrlOneToOneMapSourceReferenceURI is the URI for the source reference
// var CrlOneToOneMapSourceReferenceURI = CrlOneToOneMapURI + "/SourceReference"

// CrlOneToOneMapSourceReferenceRefinementURI is the URI of the refinement showing it to be a refinement of CrlMapSource
// var CrlOneToOneMapSourceReferenceRefinementURI = CrlOneToOneMapSourceReferenceURI + "/Refinement"

// CrlOneToOneMapTargetReferenceURI is the URI for the target reference
// var CrlOneToOneMapTargetReferenceURI = CrlOneToOneMapURI + "/TargetReference"

// CrlOneToOneMapTargetReferenceRefinementURI is the URI of the refinement showing it to be a refinement of CrlMapTarget
// var CrlOneToOneMapTargetReferenceRefinementURI = CrlOneToOneMapTargetReferenceURI + "/Refinement"

// Reference to Element Map

// CrlReferenceToElementMapURI is the URI for the Reference to Element Map
// var CrlReferenceToElementMapURI = CrlMapsDomainURI + "/ReferenceToElementMap"

// CrlReferenceToElementMapRefinementURI is the URI for the refinement from CrlMap
// var CrlReferenceToElementMapRefinementURI = CrlReferenceToElementMapURI + "/Refinement"

// CrlReferenceToElementMapSourceURI is the URI for the source
// var CrlReferenceToElementMapSourceURI = CrlReferenceToElementMapURI + "/Source"

// CrlReferenceToElementMapSourceRefinementURI is the URI for the refinement from CrlMapSource
// var CrlReferenceToElementMapSourceRefinementURI = CrlReferenceToElementMapSourceURI + "/Refinement"

// CrlReferenceToElementMapTargetURI is the URI for the target
// var CrlReferenceToElementMapTargetURI = CrlReferenceToElementMapURI + "/Target"

// CrlReferenceToElementMapTargetRefinementURI is the URI for the refinement from CrlMapTarget
// var CrlReferenceToElementMapTargetRefinementURI = CrlReferenceToElementMapTargetURI + "Refinement"

// ID to Reference Map

// CrlIDToReferenceMapURI is the URI for a map from an attribute that is an ID to a Reference
// var CrlIDToReferenceMapURI = CrlMapsDomainURI + "/IDToReferenceMap"

// CrlIDToReferenceMapRefinementURI is the URI for the refinement from CrlMap
// var CrlIDToReferenceMapRefinementURI = CrlIDToReferenceMapURI + "/Refinement"

// CrlIDToReferenceMapSourceURI is the URI for the source reference
// var CrlIDToReferenceMapSourceURI = CrlIDToReferenceMapURI + "/Source"

// CrlIDToReferenceMapSourceRefinementURI is the URI for the refinement from CrlMapSource
// var CrlIDToReferenceMapSourceRefinementURI = CrlIDToReferenceMapSourceURI + "/Refinement"

// CrlIDToReferenceMapTargetURI is the URI for the target reference
// var CrlIDToReferenceMapTargetURI = CrlIDToReferenceMapURI + "/Target"

// CrlIDToReferenceMapTargetRefinementURI is the URI for the refinement from CrlMapTarget
// var CrlIDToReferenceMapTargetRefinementURI = CrlIDToReferenceMapTargetURI + "/Refinement"

// NewOneToOneMap creates an instance of a one-to-one map with its source and target references
func NewOneToOneMap(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (core.Concept, error) {
	newMap, _ := uOfD.CreateRefinementOfConceptURI(CrlOneToOneMapURI, "OneToOneMap", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlMapSourceURI, newMap, "Source", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlMapTargetURI, newMap, "Target", trans)
	return newMap, nil
}

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
	uOfD.NewOwnedReference(crlMap, "MapSource", trans, CrlMapSourceURI)
	uOfD.NewOwnedReference(crlMap, "MapTarget", trans, CrlMapTargetURI)

	// One to One Map
	uOfD.CreateOwnedRefinementOfConcept(crlMap, crlMapsDomain, "CrlOneToOneMap", trans, CrlOneToOneMapURI)

	// Reference To Element Map
	// uOfD.NewOwnedElement(crlMapsDomain, "ReferenceToElementMap", trans, CrlReferenceToElementMapURI)

	// ID to Reference Map
	// uOfD.NewOwnedElement(crlMapsDomain, "IDToReferenceMap", trans, CrlIDToReferenceMapURI)

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
func executeOneToOneMap(mapInstance core.Concept, notification *core.ChangeNotification, trans *core.Transaction) error {
	uOfD := trans.GetUniverseOfDiscourse()
	var err error = nil
	trans.WriteLockElement(mapInstance)

	// As an initial assumption, it probably doesn't matter what kind of notification has been received.

	// The mapInstance must have an owner for meaningful execution of the map. If there is no owner, bail (this
	// occurs during the replicate as refinement creation of the map instance).
	if mapInstance.GetOwningConceptID(trans) == "" {
		return nil
	}

	// Validate that this instance is a refinement of a defining one-to-one map: a refinement of CrlOneToOneMap
	// This function can be called when creating a direct refinement of a one-to-one map. This is a defining map, and
	// we don't want to execute this function on defining maps, only on their instantiations. The following check eliminates
	// those calls
	var immediateAbstractions = map[string]core.Concept{}
	mapInstance.FindImmediateAbstractions(immediateAbstractions, trans)
	var definingMap core.Concept
	for _, abs := range immediateAbstractions {
		if abs.GetURI(trans) != CrlOneToOneMapURI && abs.IsRefinementOfURI(CrlOneToOneMapURI, trans) {
			definingMap = abs
			break
		}
	}
	if definingMap == nil {
		return nil
	}

	// Only report the maps that are actually being executed
	// log.Printf("Executing executeOneToOneMap for map labeled %s", mapInstance.GetLabel(trans))

	// Validate that the abstraction has a sourceRef and that the sourceRef is referencing an element
	definingSourceRef := definingMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
	if definingSourceRef == nil {
		return nil
	}
	definingSource := definingSourceRef.GetReferencedConcept(trans)
	if definingSource == nil {
		return nil
	}
	// Validate that the defining map has a targetRef and that the targetRef is referencing an element
	definingTargetRef := definingMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
	if definingTargetRef == nil {
		return nil
	}
	definingTarget := definingTargetRef.GetReferencedConcept(trans)
	if definingTarget == nil {
		return nil
	}
	// Check to see whether the source reference exists and references an element of the correct type
	sourceRef := mapInstance.GetFirstOwnedReferenceRefinedFrom(definingSourceRef, trans)
	if sourceRef == nil {
		sourceRef, err = uOfD.CreateOwnedRefinementOfConcept(definingSourceRef, mapInstance, "Source", trans)
		if err != nil {
			return errors.Wrap(err, "executeOneToOneMap failed")
		}
	}

	// Now explore the targetRef
	targetRef := mapInstance.GetFirstOwnedReferenceRefinedFrom(definingTargetRef, trans)
	// If the target ref does not exist, create it
	if targetRef == nil {
		targetRef, err = uOfD.CreateOwnedRefinementOfConcept(definingTargetRef, mapInstance, "Target", trans)
		if err != nil {
			return errors.Wrap(err, "executeOneToOneMap failed")
		}
	}

	// Now the source and target
	source := sourceRef.GetReferencedConcept(trans)
	if source == nil || !source.IsRefinementOf(definingSource, trans) {
		return nil
	}
	target := targetRef.GetReferencedConcept(trans)
	targetRefAttributeName := targetRef.GetReferencedAttributeName(trans)
	switch targetRefAttributeName {
	case core.NoAttribute:
		// We're pointing to an element. We need to create one if necessary
		if target == nil && mapInstance.GetOwningConcept(trans) != nil {
			if isRootMap(mapInstance, trans) {
				targetOwner := mapInstance.GetOwningConcept(trans)
				target, err = uOfD.NewOwnedElement(targetOwner, "Instance of Target Domain", trans)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
				_, err = uOfD.NewCompleteRefinement(definingTarget, target, "Refinement of "+definingTarget.GetLabel(trans), trans)
				if err != nil {
					return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
				}
				SetTarget(mapInstance, target, core.NoAttribute, trans)
			} else {
				target, err = uOfD.CreateRefinementOfConcept(definingTarget, definingTarget.GetLabel(trans), trans)
				if err != nil {
					return errors.Wrap(err, "executeOneToOneMap failed")
				}
				SetTarget(mapInstance, target, core.NoAttribute, trans)
				candidateTargetOwner := GetTarget(mapInstance.GetOwningConcept(trans), trans)
				if candidateTargetOwner != nil && target.GetOwningConceptID(trans) != candidateTargetOwner.GetConceptID(trans) {
					abstractTargetOwner := definingTarget.GetOwningConcept(trans)
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
		if target == nil {
			target, err = getAttributeTarget(mapInstance, trans)
			if err != nil {
				return errors.Wrap(err, "executeOneToOneMap failed")
			}
			SetTarget(mapInstance, target, targetRefAttributeName, trans)
		}
		// Make value assignments as required
		// If the sourceRef is an attribute value reference, get the source value
		sourceAttributeName := sourceRef.GetReferencedAttributeName(trans)
		if sourceAttributeName != core.NoAttribute {
			var sourceReferentValue string
			var sourceReferent core.Concept
			switch sourceAttributeName {
			case core.ReferencedConceptID:
				switch source.GetConceptType() {
				case core.Reference:
					sourceReferentValue = source.GetReferencedConceptID(trans)
					sourceReferent = uOfD.GetElement(sourceReferentValue)
				}
			case core.AbstractConceptID:
				switch source.GetConceptType() {
				case core.Refinement:
					sourceReferentValue = source.GetAbstractConceptID(trans)
					sourceReferent = uOfD.GetElement(sourceReferentValue)
				}
			case core.RefinedConceptID:
				switch source.GetConceptType() {
				case core.Refinement:
					sourceReferentValue = source.GetRefinedConceptID(trans)
					sourceReferent = uOfD.GetElement(sourceReferentValue)
				}
			case core.OwningConceptID:
				sourceReferentValue = source.GetOwningConceptID(trans)
				sourceReferent = uOfD.GetElement(sourceReferentValue)
			case core.LiteralValue, core.Definition, core.Label:
				sourceReferentValue = source.GetLiteralValue(trans)
			}

			// We now know the source referent value. If it is an Element and we are setting a target pointer,
			// we need to find the corresponding target referent Element

			var targetReferent core.Concept
			switch targetRefAttributeName {
			case core.OwningConceptID, core.ReferencedConceptID, core.RefinedConceptID, core.AbstractConceptID:
				if sourceReferent == nil {
					return nil
				}
				targetReferentMap := SearchForMapForSource(mapInstance, sourceReferent, trans)
				if targetReferentMap == nil {
					return nil
				}
				targetReferent = GetTarget(targetReferentMap, trans)
				if targetReferent == nil {
					return nil
				}
				if target == nil {
					return nil
				}
				// Now actually assign the values
				switch targetRefAttributeName {
				case core.OwningConceptID:
					target.SetOwningConcept(targetReferent, trans)
				case core.ReferencedConceptID:
					switch target.GetConceptType() {
					case core.Reference:
						target.SetReferencedConcept(targetReferent, core.NoAttribute, trans)
					}
				case core.RefinedConceptID:
					switch target.GetConceptType() {
					case core.Refinement:
						target.SetRefinedConcept(targetReferent, trans)
					}
				case core.AbstractConceptID:
					switch target.GetConceptType() {
					case core.Refinement:
						target.SetAbstractConcept(targetReferent, trans)
					}
				}
			case core.Label:
				target.SetLabel(sourceReferentValue, trans)
			case core.Definition:
				target.SetDefinition(sourceReferentValue, trans)
			case core.LiteralValue:
				switch target.GetConceptType() {
				case core.Literal:
					target.SetLiteralValue(sourceReferentValue, trans)
				}
			}
		}
	}

	// Now take care of map children.
	err = instantiateMapChildren(definingMap, mapInstance, source, target, uOfD, trans)
	if err != nil {
		return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
	}
	err = tickleMapChildren(mapInstance, trans)
	if err != nil {
		return errors.Wrap(err, "crlmaps.executeOneToOneMap failed")
	}
	return nil
}

func getAbstractMap(thisMap core.Concept, trans *core.Transaction) core.Concept {
	immediateAbstractions := map[string]core.Concept{}
	thisMap.FindImmediateAbstractions(immediateAbstractions, trans)
	for _, abstraction := range immediateAbstractions {
		if abstraction.IsRefinementOfURI(CrlOneToOneMapURI, trans) {
			return abstraction
		}
	}
	return nil
}

func getAttributeTarget(attributeMap core.Concept, trans *core.Transaction) (core.Concept, error) {
	// Assumes that the parent map's target is either the desired target or an ancestor of the desired target
	// The target is either going to be the parent map's target or one of its descendants
	parentTarget := getParentMapTarget(attributeMap, trans)
	if parentTarget == nil {
		return nil, nil
	}
	// Get the abstract attributeMap
	abstractMap := getAbstractMap(attributeMap, trans)
	// Get the abstract attributeMap's target
	abstractTarget := GetTarget(abstractMap, trans)
	// Find the descendent of the parent target that has the attributeMap's target as an ancestor
	if parentTarget.IsRefinementOf(abstractTarget, trans) {
		return parentTarget, nil
	}
	// See if the parent target is of the right type to be the parent of the abstract target.
	// If it is not, return nil (i.e. not found)
	// If it is, search for an existing unmapped child that is of the correct type and not already the target of a
	// map. If found, use it as the child target. If one is not found, create it and
	// make the parent target its parent
	candidateChildTargets := parentTarget.GetOwnedConceptsRefinedFrom(abstractTarget, trans)
	for _, candidateChildTarget := range candidateChildTargets {
		if !hasTargetListener(candidateChildTarget, trans) {
			return candidateChildTarget, nil
		}
	}
	childTarget, err := trans.GetUniverseOfDiscourse().CreateRefinementOfConcept(abstractTarget, abstractTarget.GetLabel(trans), trans)
	if err != nil {
		return nil, errors.Wrap(err, "getAttributeTarget failed")
	}
	childTarget.SetOwningConcept(parentTarget, trans)
	return childTarget, nil
}

// GetSource returns the source referenced by the given map
func GetSource(theMap core.Concept, trans *core.Transaction) core.Concept {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(trans)
}

// GetSourceReference returns the source reference for the given map
func GetSourceReference(theMap core.Concept, trans *core.Transaction) core.Concept {
	return theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
}

// GetTarget returns the target referenced by the given map
func GetTarget(theMap core.Concept, trans *core.Transaction) core.Concept {
	ref := theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
	if ref == nil {
		return nil
	}
	return ref.GetReferencedConcept(trans)
}

// GetTargetReference returns the target reference for the given map
func GetTargetReference(theMap core.Concept, trans *core.Transaction) core.Concept {
	return theMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapTargetURI, trans)
}

func getRootMap(theMap core.Concept, trans *core.Transaction) core.Concept {
	owner := theMap.GetOwningConcept(trans)
	if owner != nil && isMap(owner, trans) {
		return getRootMap(owner, trans)
	}
	return theMap
}

func getParentMapTarget(theMap core.Concept, trans *core.Transaction) core.Concept {
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
func FindAttributeMapForSource(currentMap core.Concept, source core.Concept, attributeName core.AttributeName, trans *core.Transaction) core.Concept {
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
func FindMapForSource(currentMap core.Concept, source core.Concept, trans *core.Transaction) core.Concept {
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
func SearchForMapForSource(currentMap core.Concept, source core.Concept, trans *core.Transaction) core.Concept {
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
func FindMapForSourceAttribute(currentMap core.Concept, source core.Concept, attributeName core.AttributeName, trans *core.Transaction) core.Concept {
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
func FindTargetForSource(currentMap core.Concept, source core.Concept, trans *core.Transaction) core.Concept {
	// get root map
	rootMap := getRootMap(currentMap, trans)
	// search the root map for a mapping whose source is the given source. If found, return target of the map
	foundMap := FindMapForSource(rootMap, source, trans)
	if foundMap != nil {
		return GetTarget(foundMap, trans)
	}
	return nil
}

func hasTargetListener(el core.Concept, trans *core.Transaction) bool {
	uOfD := trans.GetUniverseOfDiscourse()
	it := uOfD.GetListenerIDs(el.GetConceptID(trans)).Iterator()
	for listenerID := range it.C {
		listener := uOfD.GetElement(listenerID.(string))
		if listener.IsRefinementOfURI(CrlMapTargetURI, trans) {
			return true
		}
	}
	return false
}

func instantiateMapChildren(parentDefiningMap core.Concept, parentInstanceMap core.Concept, source core.Concept, target core.Concept, uOfD *core.UniverseOfDiscourse, trans *core.Transaction) error {
	// for each of the abstractMap's children that is a map
	for _, definingChildMap := range parentDefiningMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, trans) {
		definingChildMapSource := GetSource(definingChildMap, trans)
		if definingChildMapSource != nil {
			// There are two cases here, depending upon whether the source reference is to a pointer or an element.
			definingChildMapSourceReference := definingChildMap.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
			if definingChildMapSourceReference == nil {
				return errors.New("In CrlMaps.go instantiateMapChildren, the definingChildMapSource does not have a definingChildMapSourceReference")
			}
			definingChildMapSourceReferenceAttributeName := definingChildMapSourceReference.GetReferencedAttributeName(trans)
			if definingChildMapSourceReferenceAttributeName != core.NoAttribute {
				// If the abstractChildMap's source reference is to a pointer, then the actual source for the child is going to be
				// the parent's source. Error checking is required to ensure that the parent's source is of the appropriate type for the AttributeName
				// on the reference. In this case there will only be one instance of the abstractChildMap created.
				// Check to see whether there is already a map instance for this source
				parentInstanceMapSource := GetSource(parentInstanceMap, trans)
				if parentInstanceMapSource == nil {
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

				foundChildSource := parentInstanceMapSource // assume it's going to be the parent map source
				if !parentInstanceMapSource.IsRefinementOf(definingChildMapSource, trans) {
					if definingChildMapSourceReference.GetReferencedAttributeName(trans) == core.ReferencedConceptID {
						foundChildSource = parentInstanceMapSource.GetFirstOwnedReferenceRefinedFrom(definingChildMapSource, trans)
						if foundChildSource == nil {
							return nil
						}
					} else {
						return nil
					}
				}
				var newMapInstance core.Concept
				for _, mapInstance := range parentInstanceMap.GetOwnedConceptsRefinedFrom(definingChildMap, trans) {
					mapInstanceSource := GetSource(mapInstance, trans)
					if mapInstanceSource == nil || mapInstanceSource.GetConceptID(trans) == foundChildSource.GetConceptID(trans) {
						newMapInstance = mapInstance
						break
					}
				}
				if newMapInstance == nil {
					newMapInstance, err := uOfD.CreateRefinementOfConcept(definingChildMap, "Instance of "+definingChildMap.GetLabel(trans), trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateMapChildren failed")
					}
					err = newMapInstance.SetOwningConcept(parentInstanceMap, trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateMapChildren failed")
					}
					err = newMapInstance.SetLabel("Instance of "+definingChildMap.GetLabel(trans), trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateMapChildren failed")
					}
				}
				// During a deletion scenario newMapInstance can be null even though one was instantiated above. Not sure why
				if newMapInstance == nil {
					return errors.New("in crlmaps.instantiateMapChildred, newMapInstance is nil just after creating it")
				}
				newSourceRef := newMapInstance.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
				if newSourceRef == nil {
					return errors.New("In crlmaps.instantiateMapChildren, newSourceRef is nil")
				}
				if foundChildSource != nil && newSourceRef.GetReferencedConceptID(trans) != foundChildSource.GetConceptID(trans) {
					err := newSourceRef.SetReferencedConcept(foundChildSource, definingChildMapSourceReferenceAttributeName, trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateMapChildren failed")
					}
				}
			} else {
				// The abstractChildMapSourceReference is to an element. There is an assumption in this case that the parent's source
				// contains elements that are refinements of the abstractChildMapSource. For each element that is a refinement of the
				// abstractChildMapSource found in the parent's source, instantiate the abstractChildMap (replicate as refinement)
				// and wire up the element as the source
				for _, sourceEl := range source.GetOwnedDescendantsRefinedFrom(definingChildMapSource, trans) {
					// Check to see whether there is already a map instance for this source
					var newMapInstance core.Concept
					for _, mapInstance := range parentInstanceMap.GetOwnedConceptsRefinedFrom(definingChildMap, trans) {
						mapInstanceSource := GetSource(mapInstance, trans)
						if mapInstanceSource == nil || mapInstanceSource.GetConceptID(trans) == sourceEl.GetConceptID(trans) {
							newMapInstance = mapInstance
							break
						}
					}
					if newMapInstance == nil {
						var err error
						newMapInstance, err = uOfD.CreateRefinementOfConcept(definingChildMap, "Instance of "+definingChildMap.GetLabel(trans), trans)
						if err != nil {
							return errors.Wrap(err, "crlmaps.instantiateMapChildren failed")
						}
						err = newMapInstance.SetOwningConcept(parentInstanceMap, trans)
						if err != nil {
							return errors.Wrap(err, "crlmaps.instantiateMapChildren failed")
						}
					}
					newSourceRef := newMapInstance.GetFirstOwnedReferenceRefinedFromURI(CrlMapSourceURI, trans)
					if newSourceRef == nil {
						return errors.New("In crlmaps.instantiateMapChildren, newSourceRef is nil")
					}
					err := newSourceRef.SetReferencedConcept(sourceEl, core.NoAttribute, trans)
					if err != nil {
						return errors.Wrap(err, "crlmaps.instantiateMapChildren failed")
					}
				}
			}
		}
	}
	return nil
}

func isMap(candidate core.Concept, trans *core.Transaction) bool {
	return candidate != nil && candidate.IsRefinementOfURI(CrlMapURI, trans)
}

func isRootMap(candidate core.Concept, trans *core.Transaction) bool {
	rootMap := getRootMap(candidate, trans)
	return rootMap == candidate
}

// SetSource sets the source referenced by the given map
func SetSource(theMap core.Concept, newSource core.Concept, attributeName core.AttributeName, trans *core.Transaction) error {
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
func SetTarget(theMap core.Concept, newTarget core.Concept, attributeName core.AttributeName, trans *core.Transaction) error {
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

func tickleMapChildren(parentInstanceMap core.Concept, trans *core.Transaction) error {
	// for each of the abstractMap's children that is a map
	mapChildren := parentInstanceMap.GetOwnedConceptsRefinedFromURI(CrlMapURI, trans)
	for _, childMap := range mapChildren {
		err := trans.GetUniverseOfDiscourse().SendTickleNotification(parentInstanceMap, childMap, trans)
		if err != nil {
			return errors.Wrap(err, "tickleMapChildren failed")
		}
	}
	// We have to repeat this a second time since an attempt to set a pointer will have failed if the
	// element owning the pointer (or its parent) has not yet been created.
	for _, childMap := range mapChildren {
		err := trans.GetUniverseOfDiscourse().SendTickleNotification(parentInstanceMap, childMap, trans)
		if err != nil {
			return errors.Wrap(err, "tickleMapChildren failed")
		}
	}
	return nil
}
