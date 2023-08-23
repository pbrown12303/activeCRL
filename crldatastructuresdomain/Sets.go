package crldatastructuresdomain

import (
	"errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlSetURI is the URI that identifies the prototype for sets
var CrlSetURI = CrlDataStructuresDomainURI + "/Set"

// CrlSetMemberReferenceURI is the URI that identifies the prototype for a set member reference
var CrlSetMemberReferenceURI = CrlSetURI + "/SetMemberReference"

// CrlSetTypeReferenceURI is the URI that identifies the prototype for a set type reference
var CrlSetTypeReferenceURI = CrlSetURI + "/SetTypeReference"

// NewSet creates an instance of a set
func NewSet(uOfD *core.UniverseOfDiscourse, setType core.Concept, hl *core.Transaction) (core.Concept, error) {
	if setType == nil {
		return nil, errors.New("no type specified for set")
	}
	newSet, _ := uOfD.CreateReplicateAsRefinementFromURI(CrlSetURI, hl)
	typeReference := newSet.GetFirstOwnedReferenceRefinedFromURI(CrlSetTypeReferenceURI, hl)
	typeReference.SetReferencedConcept(setType, core.NoAttribute, hl)
	return newSet, nil
}

// AddSetMember adds a member to the set
func AddSetMember(set core.Concept, newMember core.Concept, hl *core.Transaction) error {
	uOfD := set.GetUniverseOfDiscourse(hl)
	if IsSetMember(set, newMember, hl) {
		return errors.New("newMember is already a member of the set")
	}
	setType, _ := GetSetType(set, hl)
	if !newMember.IsRefinementOf(setType, hl) {
		return errors.New("NewMember is of wrong type")
	}
	newMemberReference, _ := uOfD.CreateReplicateReferenceAsRefinementFromURI(CrlSetMemberReferenceURI, hl)
	newMemberReference.SetOwningConcept(set, hl)
	newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, hl)
	return nil
}

// ClearSet removes all members from the set
func ClearSet(set core.Concept, hl *core.Transaction) {
	uOfD := set.GetUniverseOfDiscourse(hl)
	it := set.GetOwnedConceptIDs(hl).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, hl) {
			uOfD.DeleteElement(memberReference, hl)
		}
	}
}

// GetSetType returns the element that should be an abstraction of every member. It returns an error if the argument is not a set
func GetSetType(set core.Concept, hl *core.Transaction) (core.Concept, error) {
	if !set.IsRefinementOfURI(CrlSetURI, hl) {
		return nil, errors.New("argument is not a set")
	}
	typeReference := set.GetFirstOwnedReferenceRefinedFromURI(CrlSetTypeReferenceURI, hl)
	return typeReference.GetReferencedConcept(hl), nil
}

// IsSetMember returns true if the element is a memeber of the given set
func IsSetMember(set core.Concept, el core.Concept, hl *core.Transaction) bool {
	uOfD := set.GetUniverseOfDiscourse(hl)
	it := set.GetOwnedConceptIDs(hl).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, hl) && memberReference.GetReferencedConcept(hl) == el {
			it.Stop()
			return true
		}
	}
	return false
}

// RemoveSetMember removes the element from the given set
func RemoveSetMember(set core.Concept, el core.Concept, hl *core.Transaction) error {
	uOfD := set.GetUniverseOfDiscourse(hl)
	it := set.GetOwnedConceptIDs(hl).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, hl) && memberReference.GetReferencedConcept(hl) == el {
			uOfD.DeleteElement(memberReference, hl)
			it.Stop()
			return nil
		}
	}
	return errors.New("element not member of set")
}

// BuildCrlSetsConcepts builds the CrlSets concept space and adds it as a child of the provided parent concept space
func BuildCrlSetsConcepts(uOfD *core.UniverseOfDiscourse, parentSpace core.Concept, hl *core.Transaction) {
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
