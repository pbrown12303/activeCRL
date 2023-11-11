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
func NewSet(uOfD *core.UniverseOfDiscourse, setType *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if setType == nil {
		return nil, errors.New("no type specified for set")
	}
	newSet, _ := uOfD.CreateRefinementOfConceptURI(CrlSetURI, "Set", trans)
	typeReference, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlSetTypeReferenceURI, newSet, "TypeReference", trans)
	typeReference.SetReferencedConcept(setType, core.NoAttribute, trans)
	return newSet, nil
}

// NewSetMemberReference creates a SetMemberReference with its child concepts
func NewSetMemberReference(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) (*core.Concept, error) {
	setMemberReference, _ := uOfD.CreateRefinementOfConceptURI(CrlSetMemberReferenceURI, "MemberReference", trans)
	return setMemberReference, nil
}

// AddSetMember adds a member to the set
func AddSetMember(set *core.Concept, newMember *core.Concept, trans *core.Transaction) error {
	uOfD := set.GetUniverseOfDiscourse(trans)
	if IsSetMember(set, newMember, trans) {
		return errors.New("newMember is already a member of the set")
	}
	setType, _ := GetSetType(set, trans)
	if !newMember.IsRefinementOf(setType, trans) {
		return errors.New("NewMember is of wrong type")
	}
	newMemberReference, _ := NewSetMemberReference(uOfD, trans)
	newMemberReference.SetOwningConcept(set, trans)
	newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, trans)
	return nil
}

// ClearSet removes all members from the set
func ClearSet(set *core.Concept, trans *core.Transaction) {
	uOfD := set.GetUniverseOfDiscourse(trans)
	it := set.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, trans) {
			uOfD.DeleteElement(memberReference, trans)
		}
	}
}

// GetSetType returns the element that should be an abstraction of every member. It returns an error if the argument is not a set
func GetSetType(set *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if !set.IsRefinementOfURI(CrlSetURI, trans) {
		return nil, errors.New("argument is not a set")
	}
	typeReference := set.GetFirstOwnedReferenceRefinedFromURI(CrlSetTypeReferenceURI, trans)
	return typeReference.GetReferencedConcept(trans), nil
}

// IsSetMember returns true if the element is a memeber of the given set
func IsSetMember(set *core.Concept, el *core.Concept, trans *core.Transaction) bool {
	uOfD := set.GetUniverseOfDiscourse(trans)
	it := set.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, trans) && memberReference.GetReferencedConcept(trans) == el {
			it.Stop()
			return true
		}
	}
	return false
}

// RemoveSetMember removes the element from the given set
func RemoveSetMember(set *core.Concept, el *core.Concept, trans *core.Transaction) error {
	uOfD := set.GetUniverseOfDiscourse(trans)
	it := set.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlSetMemberReferenceURI, trans) && memberReference.GetReferencedConcept(trans) == el {
			uOfD.DeleteElement(memberReference, trans)
			it.Stop()
			return nil
		}
	}
	return errors.New("element not member of set")
}

// BuildCrlSetsConcepts builds the CrlSets concept space and adds it as a child of the provided parent concept space
func BuildCrlSetsConcepts(uOfD *core.UniverseOfDiscourse, parentSpace *core.Concept, trans *core.Transaction) {
	crlSet, _ := uOfD.NewElement(trans, CrlSetURI)
	crlSet.SetLabel("CrlSet", trans)
	crlSet.SetOwningConcept(parentSpace, trans)

	crlSetMemberReference, _ := uOfD.NewReference(trans, CrlSetMemberReferenceURI)
	crlSetMemberReference.SetLabel("MemberReference", trans)
	crlSetMemberReference.SetOwningConcept(parentSpace, trans)

	CrlSetTypeReference, _ := uOfD.NewReference(trans, CrlSetTypeReferenceURI)
	CrlSetTypeReference.SetLabel("TypeReference", trans)
	CrlSetTypeReference.SetOwningConcept(crlSet, trans)

}
