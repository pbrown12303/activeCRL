package crldatastructures

import (
	"errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlListURI is the URI that identifies the prototype for list
var CrlListURI = CrlDataStructuresConceptSpaceURI + "/List"

// CrlListReferenceToFirstMemberReferenceURI is the URI that identifies the prototype for the first member reference
var CrlListReferenceToFirstMemberReferenceURI = CrlListURI + "/ListReferenceToFirstMemberReference"

// CrlListReferenceToLastMemberReferenceURI is the URI that identifies the prototype for a the last member reference
var CrlListReferenceToLastMemberReferenceURI = CrlListURI + "/ListReferenceToLastMemberReference"

// CrlListMemberReferenceURI is the URI that identifies the prototype for a list member reference
var CrlListMemberReferenceURI = CrlListURI + "/ListMemberReference"

// CrlListReferenceToNextMemberReferenceURI is the URI that identifies a member reference's next member reference
var CrlListReferenceToNextMemberReferenceURI = CrlListURI + "/ReferenceToNextMemberReference"

// CrlListReferenceToPriorMemberReferenceURI is the URI that identifies a member reference's prior member reference
var CrlListReferenceToPriorMemberReferenceURI = CrlListURI + "/ReferenceToPriorMemberReference"

// CrlListTypeReferenceURI is the URI that identifies the prototype for a list type reference
var CrlListTypeReferenceURI = CrlListURI + "/ListTypeReference"

// NewList creates an instance of a list
func NewList(uOfD *core.UniverseOfDiscourse, setType core.Element, hl *core.HeldLocks) (core.Element, error) {
	if setType == nil {
		return nil, errors.New("No type specified for list")
	}
	newList, err := uOfD.CreateReplicateAsRefinementFromURI(CrlListURI, hl)
	if err != nil {
		return nil, err
	}
	typeReference := newList.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
	typeReference.SetReferencedConcept(setType, hl)
	return newList, nil
}

// AddListMemberAfter adds a member to the list after the priorMemberReference and returns the newMemberReference.
// If the priorMemberReference is nil, the member is added to the beginning of the list. An error is returned if the
// supplied list is not a list, the newElement is nil, or the priorElementReference is not a CrlListMemberReference in this list.
func AddListMemberAfter(list core.Element, priorMemberReference core.Reference, newMember core.Element, hl *core.HeldLocks) (core.Reference, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsList(list, hl) {
		return nil, errors.New("Supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("Supplied Element is nil: nil members are not allowed in CRL Lists")
	}
	if priorMemberReference != nil {
		if !priorMemberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) {
			return nil, errors.New("Supplied priorElementReference is not a CrlListMemberReference")
		}
		if priorMemberReference.GetOwningConcept(hl) != list {
			return nil, errors.New("Supplied priorMemberReference does not belong to this list")
		}
	}
	listType, _ := GetListType(list, hl)
	if newMember.IsRefinementOf(listType, hl) == false {
		return nil, errors.New("NewMember is of wrong type")
	}
	var newPostMemberReference core.Reference
	if priorMemberReference != nil {
		referenceToPostMemberReference, err := getReferenceToNextMemberReference(priorMemberReference, hl)
		if err != nil {
			return nil, err
		}
		referencedPostMemberReference := referenceToPostMemberReference.GetReferencedConcept(hl)
		if referencedPostMemberReference != nil {
			newPostMemberReference = referencedPostMemberReference.(core.Reference)
		}
	}
	newMemberReference, _ := uOfD.CreateReplicateReferenceAsRefinementFromURI(CrlListMemberReferenceURI, hl)
	newMemberReference.SetOwningConcept(list, hl)
	newMemberReference.SetReferencedConcept(newMember, hl)
	// Wire up prior references
	if priorMemberReference == nil {
		// This is the new list beginning
		referenceToFirstMemberReference, _ := getListReferenceToFirstMemberReference(list, hl)
		if referenceToFirstMemberReference != nil {
			referenceToFirstMemberReference.SetReferencedConcept(newMemberReference, hl)
		}
	} else {
		setNextMemberReference(priorMemberReference, newMemberReference, hl)
		setPriorMemberReference(newMemberReference, priorMemberReference, hl)
	}
	// Wire up next references
	if newPostMemberReference == nil {
		referenceToLastMemberReference, _ := getListReferenceToLastMemberReference(list, hl)
		if referenceToLastMemberReference != nil {
			referenceToLastMemberReference.SetReferencedConcept(newMemberReference, hl)
		}
	} else {
		setPriorMemberReference(newPostMemberReference, newMemberReference, hl)
		setNextMemberReference(newMemberReference, newPostMemberReference, hl)
	}
	return newMemberReference, nil
}

// AddListMemberBefore adds a member to the list before the postMember.
// If the postMember is nil, the member is added at the end of the list.
func AddListMemberBefore(list core.Element, postMemberReference core.Reference, newMember core.Element, hl *core.HeldLocks) (core.Reference, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsList(list, hl) {
		return nil, errors.New("Supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("Supplied Element is nil: nil members are not allowed in CRL Lists")
	}
	if postMemberReference != nil {
		if !postMemberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) {
			return nil, errors.New("Supplied postMemberReference is not a CrlListMemberReference")
		}
		if postMemberReference.GetOwningConcept(hl) != list {
			return nil, errors.New("Supplied postMemberReference does not belong to this list")
		}
	}
	listType, _ := GetListType(list, hl)
	if newMember.IsRefinementOf(listType, hl) == false {
		return nil, errors.New("NewMember is of wrong type")
	}
	var newPriorMemberReference core.Reference
	if postMemberReference != nil {
		referenceToPriorMemberReference, err := getReferenceToPriorMemberReference(postMemberReference, hl)
		if err != nil {
			return nil, err
		}
		referencedPriorMemberReference := referenceToPriorMemberReference.GetReferencedConcept(hl)
		if referencedPriorMemberReference != nil {
			newPriorMemberReference = referencedPriorMemberReference.(core.Reference)
		}
	}
	newMemberReference, _ := uOfD.CreateReplicateReferenceAsRefinementFromURI(CrlListMemberReferenceURI, hl)
	newMemberReference.SetOwningConcept(list, hl)
	newMemberReference.SetReferencedConcept(newMember, hl)
	// Wire up post references
	if postMemberReference == nil {
		// This is the new list end
		referenceToLastMemberReference, _ := getListReferenceToLastMemberReference(list, hl)
		if referenceToLastMemberReference != nil {
			referenceToLastMemberReference.SetReferencedConcept(newMemberReference, hl)
		}
	} else {
		setPriorMemberReference(postMemberReference, newMemberReference, hl)
		setNextMemberReference(newMemberReference, postMemberReference, hl)
	}
	// Wire up prior references
	if newPriorMemberReference == nil {
		referenceToFirstMemberReference, _ := getListReferenceToFirstMemberReference(list, hl)
		if referenceToFirstMemberReference != nil {
			referenceToFirstMemberReference.SetReferencedConcept(newMemberReference, hl)
		}
	} else {
		setNextMemberReference(newPriorMemberReference, newMemberReference, hl)
		setPriorMemberReference(newMemberReference, newPriorMemberReference, hl)
	}
	return newMemberReference, nil
}

// ClearList removes all members from the list
func ClearList(list core.Element, hl *core.HeldLocks) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	for id := range list.GetOwnedConceptIDs(hl).Iterator().C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) {
			uOfD.DeleteElement(memberReference, hl)
		}
	}
}

// GetFirstMemberReference returns the reference to the first member of the list. It returns an error if the
// list is not a list. It returns nil if the list is empty
func GetFirstMemberReference(list core.Element, hl *core.HeldLocks) (core.Reference, error) {
	refRef, err := getListReferenceToFirstMemberReference(list, hl)
	if err != nil {
		return nil, err
	}
	if refRef == nil {
		return nil, errors.New("No reference to first member reference found")
	}
	firstMemberReference := refRef.GetReferencedConcept(hl)
	if firstMemberReference == nil {
		return nil, nil
	}
	return firstMemberReference.(core.Reference), nil
}

// GetFirstReferenceForMember returns the first reference to the given member. It returns an error if the list is not a list.
// It returns nil if the element it is not found in the list.
func GetFirstReferenceForMember(list core.Element, member core.Element, hl *core.HeldLocks) (core.Reference, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	for id := range list.GetOwnedConceptIDs(hl).Iterator().C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil &&
			memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) &&
			memberReference.GetReferencedConcept(hl) == member {
			return memberReference.(core.Reference), nil
		}
	}
	return nil, nil
}

// GetLastMemberReference returns the reference to the last member of the list. It returns an error if list is not a list.
// It returns nil if the list is empty
func GetLastMemberReference(list core.Element, hl *core.HeldLocks) (core.Reference, error) {
	refRef, err := getListReferenceToLastMemberReference(list, hl)
	if err != nil {
		return nil, err
	}
	if refRef == nil {
		return nil, errors.New("No reference to last member reference found")
	}
	lastMemberReference := refRef.GetReferencedConcept(hl)
	if lastMemberReference == nil {
		return nil, nil
	}
	return lastMemberReference.(core.Reference), nil
}

// getListReferenceToFirstMemberReference returns the reference to the first member reference. It returns an error if list is not a List
func getListReferenceToFirstMemberReference(list core.Element, hl *core.HeldLocks) (core.Reference, error) {
	if IsList(list, hl) == false {
		return nil, errors.New("Argument is not a CrlDataStructures.List")
	}
	return list.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToFirstMemberReferenceURI, hl), nil
}

// getListReferenceToLastMemberReference returns the reference to the last member reference. It returns an error if list is not a List
func getListReferenceToLastMemberReference(list core.Element, hl *core.HeldLocks) (core.Reference, error) {
	if IsList(list, hl) == false {
		return nil, errors.New("Argument is not a CrlDataStructures.List")
	}
	return list.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToLastMemberReferenceURI, hl), nil
}

// GetListType returns the element that should be an abstraction of every member. It returns an error if the argument is not a list
func GetListType(list core.Element, hl *core.HeldLocks) (core.Element, error) {
	if list.IsRefinementOfURI(CrlListURI, hl) == false {
		return nil, errors.New("Argument is not a list")
	}
	typeReference := list.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
	return typeReference.GetReferencedConcept(hl), nil
}

// GetNextMemberReference returns the successor member reference in the list
func GetNextMemberReference(memberReference core.Reference, hl *core.HeldLocks) (core.Reference, error) {
	if !IsListMemberReference(memberReference, hl) {
		return nil, errors.New("Supplied memberReference is not a refinement of CrlListMemberReference")
	}
	referenceToNextMemberReference, err := getReferenceToNextMemberReference(memberReference, hl)
	if err != nil {
		return nil, err
	}
	nextMemberReference := referenceToNextMemberReference.GetReferencedConcept(hl)
	if nextMemberReference == nil {
		return nil, nil
	}
	return nextMemberReference.(core.Reference), nil
}

// GetPriorMemberReference returns the predecessor member reference in the list
func GetPriorMemberReference(memberReference core.Reference, hl *core.HeldLocks) (core.Reference, error) {
	if !IsListMemberReference(memberReference, hl) {
		return nil, errors.New("Supplied memberReference is not a refinement of CrlListMemberReference")
	}
	referenceToPriorMemberReference, err := getReferenceToPriorMemberReference(memberReference, hl)
	if err != nil {
		return nil, err
	}
	priorMemberReference := referenceToPriorMemberReference.GetReferencedConcept(hl)
	if priorMemberReference == nil {
		return nil, nil
	}
	return priorMemberReference.(core.Reference), nil
}

// getReferenceToNextMemberReference returns the reference to the next member of the list.
// It returns nil if the reference is the last member of the list
func getReferenceToNextMemberReference(memberReference core.Reference, hl *core.HeldLocks) (core.Reference, error) {
	if memberReference == nil {
		return nil, errors.New("GetNextMemberReference called with nil memberReference")
	}
	if !IsListMemberReference(memberReference, hl) {
		return nil, errors.New("Supplied memberReference is not a refinement of CrlListMemberReference")
	}
	nextMemberReference := memberReference.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToNextMemberReferenceURI, hl)
	if nextMemberReference == nil {
		return nil, errors.New("In GetNextMemberReference, memberReference does not ave a NextMemberReferenceReference")
	}
	return nextMemberReference, nil
}

// getReferenceToPriorMemberReference returns the reference to the previous member of the list. It returns an error if the memberReference
// is either nil or is not a refinement of CrlListMemberReference
// It returns nil if the reference is the first member of the list
func getReferenceToPriorMemberReference(memberReference core.Reference, hl *core.HeldLocks) (core.Reference, error) {
	if memberReference == nil {
		return nil, errors.New("GetPriorMemberReference called with nil memberReference")
	}
	if !IsListMemberReference(memberReference, hl) {
		return nil, errors.New("Supplied memberReference is not a refinement of CrlListMemberReference")
	}
	priorMemberReference := memberReference.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToPriorMemberReferenceURI, hl)
	if priorMemberReference == nil {
		return nil, errors.New("In GetPriorMemberReference, memberReference does not have a PriorMemberReferenceReference")
	}
	return priorMemberReference, nil
}

// IsList returns true if the supplied Element is a refinement of List
func IsList(list core.Element, hl *core.HeldLocks) bool {
	return list.IsRefinementOfURI(CrlListURI, hl)
}

// IsListMember returns true if the element is a memeber of the given list
func IsListMember(list core.Element, el core.Element, hl *core.HeldLocks) bool {
	uOfD := list.GetUniverseOfDiscourse(hl)
	for id := range list.GetOwnedConceptIDs(hl).Iterator().C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) && memberReference.GetReferencedConcept(hl) == el {
			return true
		}
	}
	return false
}

// IsListMemberReference returns true if the supplied Reference is a refinement of ListMemberReference
func IsListMemberReference(memberReference core.Reference, hl *core.HeldLocks) bool {
	return memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl)
}

// RemoveListMember removes the first occurrance of an element from the given list
func RemoveListMember(list core.Element, el core.Element, hl *core.HeldLocks) error {
	uOfD := list.GetUniverseOfDiscourse(hl)
	for id := range list.GetOwnedConceptIDs(hl).Iterator().C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) && memberReference.GetReferencedConcept(hl) == el {
			// Fix previous and next pointers
			priorMemberReference, _ := GetPriorMemberReference(memberReference, hl)
			nextMemberReference, _ := GetNextMemberReference(memberReference, hl)
			if priorMemberReference != nil {
				setNextMemberReference(priorMemberReference, nextMemberReference, hl)
			} else {
				referenceToFirstMemberReference, _ := getListReferenceToFirstMemberReference(list, hl)
				referenceToFirstMemberReference.SetReferencedConcept(nextMemberReference, hl)
			}
			if nextMemberReference != nil {
				setPriorMemberReference(nextMemberReference, priorMemberReference, hl)
			} else {
				referenceToLastMemberReference, _ := getListReferenceToLastMemberReference(list, hl)
				referenceToLastMemberReference.SetReferencedConcept(priorMemberReference, hl)
			}
			// Now delete the member reference
			uOfD.DeleteElement(memberReference, hl)
			return nil
		}
	}
	return errors.New("element not member of list")
}

// setNextMemberReference takes a memberReference and sets its next reference
func setNextMemberReference(memberReference core.Reference, nextReference core.Reference, hl *core.HeldLocks) {
	// since this is an internal function we assume that the references are refinements of CrlListMemberReference
	nextReferenceReference, _ := getReferenceToNextMemberReference(memberReference, hl)
	nextReferenceReference.SetReferencedConcept(nextReference, hl)
}

// setPriorMemberReference takes a memberReference and sets its prior reference
func setPriorMemberReference(memberReference core.Reference, priorReference core.Reference, hl *core.HeldLocks) {
	// since this is an internal function we assume that the references are refinements of CrlListMemberReference
	priorReferenceReference, _ := getReferenceToPriorMemberReference(memberReference, hl)
	priorReferenceReference.SetReferencedConcept(priorReference, hl)
}

// BuildCrlListsConcepts builds the CrlList concept and adds it as a child of the provided parent concept space
func BuildCrlListsConcepts(uOfD *core.UniverseOfDiscourse, parentSpace core.Element, hl *core.HeldLocks) {
	crlList, _ := uOfD.NewElement(hl, CrlListURI)
	crlList.SetLabel("CrlList", hl)
	crlList.SetOwningConcept(parentSpace, hl)
	crlList.SetIsCore(hl)

	crlFirstMemberReference, _ := uOfD.NewReference(hl, CrlListReferenceToFirstMemberReferenceURI)
	crlFirstMemberReference.SetLabel("FirstMemberReference", hl)
	crlFirstMemberReference.SetOwningConcept(crlList, hl)
	crlFirstMemberReference.SetIsCore(hl)

	crlLastMemberReference, _ := uOfD.NewReference(hl, CrlListReferenceToLastMemberReferenceURI)
	crlLastMemberReference.SetLabel("LastMemberReference", hl)
	crlLastMemberReference.SetOwningConcept(crlList, hl)
	crlLastMemberReference.SetIsCore(hl)

	CrlListTypeReference, _ := uOfD.NewReference(hl, CrlListTypeReferenceURI)
	CrlListTypeReference.SetLabel("TypeReference", hl)
	CrlListTypeReference.SetOwningConcept(crlList, hl)
	CrlListTypeReference.SetIsCore(hl)

	crlListMemberReference, _ := uOfD.NewReference(hl, CrlListMemberReferenceURI)
	crlListMemberReference.SetLabel("MemberReference", hl)
	crlListMemberReference.SetOwningConcept(parentSpace, hl)
	crlListMemberReference.SetIsCore(hl)

	crlNextMemberReference, _ := uOfD.NewReference(hl, CrlListReferenceToNextMemberReferenceURI)
	crlNextMemberReference.SetLabel("NextMemberReference", hl)
	crlNextMemberReference.SetOwningConcept(crlListMemberReference, hl)
	crlNextMemberReference.SetIsCore(hl)

	crlPriorMemberReference, _ := uOfD.NewReference(hl, CrlListReferenceToPriorMemberReferenceURI)
	crlPriorMemberReference.SetLabel("PriorMemberReference", hl)
	crlPriorMemberReference.SetOwningConcept(crlListMemberReference, hl)
	crlPriorMemberReference.SetIsCore(hl)

}
