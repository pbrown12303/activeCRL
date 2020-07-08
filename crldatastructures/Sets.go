package crldatastructures

import (
	"errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlSetURI is the URI that identifies the prototype for sets
var CrlSetURI = CrlDataStructuresConceptSpaceURI + "/Set"

// CrlSetMemberReferenceURI is the URI that identifies the prototype for a set member reference
var CrlSetMemberReferenceURI = CrlSetURI + "/SetMemberReference"

// CrlSetTypeReferenceURI is the URI that identifies the prototype for a set type reference
var CrlSetTypeReferenceURI = CrlSetURI + "/SetTypeReference"

// NewSet creates an instance of a set
func NewSet(uOfD *core.UniverseOfDiscourse, setType core.Element, hl *core.HeldLocks) (core.Element, error) {
	if setType == nil {
		return nil, errors.New("No type specified for set")
	}
	newSet, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlSetURI, hl)
	typeReference := newSet.GetFirstOwnedReferenceRefinedFromURI(CrlSetTypeReferenceURI, hl)
	typeReference.SetReferencedConcept(setType, hl)
	return newSet, nil
}

// AddSetMember adds a member to the set
func AddSetMember(set core.Element, newMember core.Element, hl *core.HeldLocks) error {
	uOfD := set.GetUniverseOfDiscourse(hl)
	if IsSetMember(set, newMember, hl) {
		return errors.New("newMember is already a member of the set")
	}
	setType, _ := GetSetType(set, hl)
	if newMember.IsRefinementOf(setType, hl) == false {
		return errors.New("NewMember is of wrong type")
	}
	newMemberReference, _ := uOfD.CreateReplicateReferenceAsRefinementFromURI(CrlSetMemberReferenceURI, hl)
	newMemberReference.SetOwningConcept(set, hl)
	newMemberReference.SetReferencedConcept(newMember, hl)
	return nil
}

// ClearSet removes all members from the set
func ClearSet(set core.Element, hl *core.HeldLocks) {
	uOfD := set.GetUniverseOfDiscourse(hl)
	it := set.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, hl) {
			uOfD.DeleteElement(memberReference, hl)
		}
	}
}

// GetSetType returns the element that should be an abstraction of every member. It returns an error if the argument is not a set
func GetSetType(set core.Element, hl *core.HeldLocks) (core.Element, error) {
	if set.IsRefinementOfURI(CrlSetURI, hl) == false {
		return nil, errors.New("Argument is not a set")
	}
	typeReference := set.GetFirstOwnedReferenceRefinedFromURI(CrlSetTypeReferenceURI, hl)
	return typeReference.GetReferencedConcept(hl), nil
}

// IsSetMember returns true if the element is a memeber of the given set
func IsSetMember(set core.Element, el core.Element, hl *core.HeldLocks) bool {
	uOfD := set.GetUniverseOfDiscourse(hl)
	it := set.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, hl) && memberReference.GetReferencedConcept(hl) == el {
			return true
		}
	}
	return false
}

// RemoveSetMember removes the element from the given set
func RemoveSetMember(set core.Element, el core.Element, hl *core.HeldLocks) error {
	uOfD := set.GetUniverseOfDiscourse(hl)
	it := set.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, hl) && memberReference.GetReferencedConcept(hl) == el {
			uOfD.DeleteElement(memberReference, hl)
			return nil
		}
	}
	return errors.New("element not member of set")
}

// BuildCrlSetsConcepts builds the CrlSets concept space and adds it as a child of the provided parent concept space
func BuildCrlSetsConcepts(uOfD *core.UniverseOfDiscourse, parentSpace core.Element, hl *core.HeldLocks) {
	crlSet, _ := uOfD.NewElement(hl, CrlSetURI)
	crlSet.SetLabel("CrlSet", hl)
	crlSet.SetOwningConcept(parentSpace, hl)

	crlSetMemberReference, _ := uOfD.NewReference(hl, CrlSetMemberReferenceURI)
	crlSetMemberReference.SetLabel("MemberReference", hl)
	crlSetMemberReference.SetOwningConcept(parentSpace, hl)

	CrlSetTypeReference, _ := uOfD.NewReference(hl, CrlSetTypeReferenceURI)
	CrlSetTypeReference.SetLabel("TypeReference", hl)
	CrlSetTypeReference.SetOwningConcept(crlSet, hl)

}
