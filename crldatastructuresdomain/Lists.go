package crldatastructuresdomain

import (
	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlListURI is the URI that identifies the prototype for list
var CrlListURI = CrlDataStructuresDomainURI + "/List"

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
func NewList(uOfD *core.UniverseOfDiscourse, setType *core.Concept, trans *core.Transaction, newURI ...string) (*core.Concept, error) {
	if setType == nil {
		return nil, errors.New("No type specified for list")
	}
	newList, err := uOfD.CreateRefinementOfConceptURI(CrlListURI, "List", trans, newURI...)
	if err != nil {
		return nil, err
	}
	addNewListConcepts(uOfD, newList, trans)
	typeReference := newList.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, trans)
	if typeReference == nil {
		return nil, errors.New("In Lists.go, NewList failed to find a type reference")
	}
	typeReference.SetReferencedConcept(setType, core.NoAttribute, trans)
	return newList, nil
}

// NewListMemberReference creates a list member reference with its child concepts
func NewListMemberReference(uOfD *core.UniverseOfDiscourse, trans *core.Transaction, newURI ...string) (*core.Concept, error) {
	newReference, _ := uOfD.CreateRefinementOfConceptURI(CrlListMemberReferenceURI, "MemberReference", trans, newURI...)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlListReferenceToNextMemberReferenceURI, newReference, "NextMemberReference", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlListReferenceToPriorMemberReferenceURI, newReference, "PriorMemberReference", trans)
	return newReference, nil
}

func addNewListConcepts(uOfD *core.UniverseOfDiscourse, newList *core.Concept, trans *core.Transaction) {
	uOfD.CreateOwnedRefinementOfConceptURI(CrlListReferenceToFirstMemberReferenceURI, newList, "FirstMemberReference", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlListReferenceToLastMemberReferenceURI, newList, "LastMemberReference", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlListTypeReferenceURI, newList, "TypeReference", trans)
}

// AddListMemberAfter adds a member to the list after the priorMemberReference and returns the newMemberReference.
// If the priorMemberReference is nil, the member is added to the beginning of the list. An error is returned if the
// supplied list is not a list, the newElement is nil, or the priorElementReference is not a CrlListMemberReference in this list.
func AddListMemberAfter(list *core.Concept, priorMemberReference *core.Concept, newMember *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	if !IsList(list, trans) {
		return nil, errors.New("In AddListMemberAfter, supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("In AddListMemberAfter, newMember is nil: nil members are not allowed in CRL Lists")
	}
	// validate prior member reference
	if priorMemberReference == nil {
		return nil, errors.New("In AddListMemberAfter, priorMemberReference is nil: this is a required value")
	}
	if priorMemberReference != nil {
		if !priorMemberReference.IsRefinementOfURI(CrlListMemberReferenceURI, trans) {
			return nil, errors.New("In AddListMemberAfter, supplied priorElementReference is not a CrlListMemberReference")
		}
		if priorMemberReference.GetOwningConcept(trans) != list {
			return nil, errors.New("In AddListMemberAfter, supplied priorMemberReference does not belong to this list")
		}
	}
	listType, _ := GetListType(list, trans)
	if !newMember.IsRefinementOf(listType, trans) {
		return nil, errors.New("In AddListMemberAfter, newMember is of wrong type")
	}
	var newPostMemberReference *core.Concept
	referenceToPostMemberReference, err := getReferenceToNextMemberReference(priorMemberReference, trans)
	if err != nil {
		return nil, err
	}
	referencedPostMemberReference := referenceToPostMemberReference.GetReferencedConcept(trans)
	if referencedPostMemberReference != nil {
		newPostMemberReference = referencedPostMemberReference
	}
	newMemberReference, _ := NewListMemberReference(uOfD, trans)
	newMemberReference.SetOwningConcept(list, trans)
	newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, trans)
	// Wire up prior references
	setNextMemberReference(priorMemberReference, newMemberReference, trans)

	setPriorMemberReference(newMemberReference, priorMemberReference, trans)
	// Wire up next references
	if newPostMemberReference == nil {
		referenceToLastMemberReference, err2 := getListReferenceToLastMemberReference(list, trans)
		if err2 != nil {
			return nil, errors.Wrap(err2, "AddListMemberAfter failed")
		}
		if referenceToLastMemberReference != nil {
			referenceToLastMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, trans)
		}
	} else {
		setPriorMemberReference(newPostMemberReference, newMemberReference, trans)
		setNextMemberReference(newMemberReference, newPostMemberReference, trans)
	}

	return newMemberReference, nil
}

// AddListMemberBefore adds a member to the list before the postMember.
// If the postMember is nil, the member is added at the end of the list.
func AddListMemberBefore(list *core.Concept, postMemberReference *core.Concept, newMember *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	if !IsList(list, trans) {
		return nil, errors.New("In AddListMemberBefore, Supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("In AddListMemberBefore, Supplied Element is nil: nil members are not allowed in CRL Lists")
	}
	// Check to ensure that the postMemberReference is valid
	if postMemberReference == nil {
		return nil, errors.New("AddListMemberBefore called with nil postMemberReference")
	}
	if !postMemberReference.IsRefinementOfURI(CrlListMemberReferenceURI, trans) {
		return nil, errors.New("In AddListMemberBefore, Supplied postMemberReference is not a CrlListMemberReference")
	}
	if postMemberReference.GetOwningConcept(trans) != list {
		return nil, errors.New("In AddListMemberBefore, Supplied postMemberReference does not belong to this list")
	}
	listType, _ := GetListType(list, trans)
	if !newMember.IsRefinementOf(listType, trans) {
		return nil, errors.New("In AddListMemberBefore, NewMember is of wrong type")
	}
	var newPriorMemberReference *core.Concept
	// If the postMemberReference exists, then its priorMemberReference should point to the newMemberReference
	referenceToPriorMemberReference, err := getReferenceToPriorMemberReference(postMemberReference, trans)
	if err != nil {
		return nil, errors.Wrap(err, "AddListMemberBefore failed")
	}
	referencedPriorMemberReference := referenceToPriorMemberReference.GetReferencedConcept(trans)
	if referencedPriorMemberReference != nil {
		newPriorMemberReference = referencedPriorMemberReference
	}
	// Create the newMemberReference
	newMemberReference, _ := NewListMemberReference(uOfD, trans)
	newMemberReference.SetOwningConcept(list, trans)
	newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, trans)
	// Wire up post references - be careful if inserting at the end
	setPriorMemberReference(postMemberReference, newMemberReference, trans)
	setNextMemberReference(newMemberReference, postMemberReference, trans)
	// Wire up prior references
	if newPriorMemberReference == nil {
		// The new member is the only member of the list
		referenceToFirstMemberReference, _ := getListReferenceToFirstMemberReference(list, trans)
		if referenceToFirstMemberReference != nil {
			referenceToFirstMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, trans)
		}
	} else {
		setNextMemberReference(newPriorMemberReference, newMemberReference, trans)
		setPriorMemberReference(newMemberReference, newPriorMemberReference, trans)
	}
	return newMemberReference, nil
}

// AppendListMember adds a member to the end of the list
func AppendListMember(list *core.Concept, newMember *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	if !IsList(list, trans) {
		return nil, errors.New("In AddListMemberBefore, Supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("In AddListMemberBefore, Supplied Element is nil: nil members are not allowed in CRL Lists")
	}
	listType, _ := GetListType(list, trans)
	if !newMember.IsRefinementOf(listType, trans) {
		return nil, errors.New("In AddListMemberBefore, NewMember is of wrong type")
	}
	oldLastMemberReference, err := GetLastMemberReference(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "AppendListMember failed")
	}
	// Create the newMemberReference
	newMemberReference, err2 := NewListMemberReference(uOfD, trans)
	if err2 != nil {
		return nil, errors.Wrap(err2, "AppendListMember failed")
	}
	err = newMemberReference.SetOwningConcept(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "AppendListMember failed")
	}
	err = newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, trans)
	if err != nil {
		return nil, errors.Wrap(err, "AppendListMember failed")
	}
	// Wire up references - be careful if inserting at the end
	referenceToLastMemberReference, err3 := getListReferenceToLastMemberReference(list, trans)
	if err3 != nil {
		return nil, errors.Wrap(err2, "AppendListMember failed")
	}
	if referenceToLastMemberReference != nil {
		err = referenceToLastMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, trans)
		if err != nil {
			return nil, errors.Wrap(err, "AppendListMember failed")
		}
	}
	if oldLastMemberReference == nil {
		referenceToFirstMemberReference, err4 := getListReferenceToFirstMemberReference(list, trans)
		if err4 != nil {
			return nil, errors.Wrap(err2, "AppendListMember failed")
		}
		err = referenceToFirstMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, trans)
		if err != nil {
			return nil, errors.Wrap(err, "AppendListMember failed")
		}
	} else {
		err = setNextMemberReference(oldLastMemberReference, newMemberReference, trans)
		if err != nil {
			return nil, errors.Wrap(err, "AppendListMember failed")
		}
		err = setPriorMemberReference(newMemberReference, oldLastMemberReference, trans)
		if err != nil {
			return nil, errors.Wrap(err, "AppendListMember failed")
		}
	}
	return newMemberReference, nil
}

// ClearList removes all members from the list
func ClearList(list *core.Concept, trans *core.Transaction) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	it := list.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, trans) {
			uOfD.DeleteElement(memberReference, trans)
		}
	}
}

// GetFirstMemberReference returns the reference to the first member of the list. It returns an error if the
// list is not a list. It returns nil if the list is empty
func GetFirstMemberReference(list *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	refRef, err := getListReferenceToFirstMemberReference(list, trans)
	if err != nil {
		return nil, err
	}
	if refRef == nil {
		return nil, errors.New("In List GetFirstMemberReference, No reference to first member reference found")
	}
	firstMemberReference := refRef.GetReferencedConcept(trans)
	if firstMemberReference == nil {
		return nil, nil
	}
	return firstMemberReference, nil
}

// GetFirstReferenceForMember returns the first reference to the given member. It returns an error if the list is not a list.
// It returns nil if the element it is not found in the list.
func GetFirstReferenceForMember(list *core.Concept, member *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	it := list.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil &&
			memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, trans) &&
			memberReference.GetReferencedConcept(trans) == member {
			it.Stop()
			return memberReference, nil
		}
	}
	return nil, nil
}

// GetLastMemberReference returns the reference to the last member of the list. It returns an error if list is not a list.
// It returns nil if the list is empty
func GetLastMemberReference(list *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	refRef, err := getListReferenceToLastMemberReference(list, trans)
	if err != nil {
		return nil, err
	}
	if refRef == nil {
		return nil, errors.New("No reference to last member reference found")
	}
	lastMemberReference := refRef.GetReferencedConcept(trans)
	if lastMemberReference == nil {
		return nil, nil
	}
	return lastMemberReference, nil
}

// getListReferenceToFirstMemberReference returns the reference to the first member reference. It returns an error if list is not a List
func getListReferenceToFirstMemberReference(list *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if !IsList(list, trans) {
		return nil, errors.New("Argument is not a CrlDataStructures.List")
	}
	refToRef := list.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToFirstMemberReferenceURI, trans)
	if refToRef == nil {
		return nil, errors.New("In getListReferenceToFirstMemberReference, the reference was not found")
	}
	return refToRef, nil
}

// getListReferenceToLastMemberReference returns the reference to the last member reference. It returns an error if list is not a List
func getListReferenceToLastMemberReference(list *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if !IsList(list, trans) {
		return nil, errors.New("Argument is not a CrlDataStructures.List")
	}
	return list.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToLastMemberReferenceURI, trans), nil
}

// GetListType returns the element that should be an abstraction of every member. It returns an error if the argument is not a list
func GetListType(list *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if !list.IsRefinementOfURI(CrlListURI, trans) {
		return nil, errors.New("Argument is not a list")
	}
	typeReference := list.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, trans)
	return typeReference.GetReferencedConcept(trans), nil
}

// GetNextMemberReference returns the successor member reference in the list
func GetNextMemberReference(memberReference *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if !IsListMemberReference(memberReference, trans) {
		return nil, errors.New("Supplied memberReference is not a refinement of CrlListMemberReference")
	}
	referenceToNextMemberReference, err := getReferenceToNextMemberReference(memberReference, trans)
	if err != nil {
		return nil, err
	}
	nextMemberReference := referenceToNextMemberReference.GetReferencedConcept(trans)
	if nextMemberReference == nil {
		return nil, nil
	}
	return nextMemberReference, nil
}

// GetPriorMemberReference returns the predecessor member reference in the list
func GetPriorMemberReference(memberReference *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if !IsListMemberReference(memberReference, trans) {
		return nil, errors.New("Supplied memberReference is not a refinement of CrlListMemberReference")
	}
	referenceToPriorMemberReference, err := getReferenceToPriorMemberReference(memberReference, trans)
	if err != nil {
		return nil, err
	}
	priorMemberReference := referenceToPriorMemberReference.GetReferencedConcept(trans)
	if priorMemberReference == nil {
		return nil, nil
	}
	return priorMemberReference, nil
}

// getReferenceToNextMemberReference returns the reference to the next member of the list.
// It returns nil if the reference is the last member of the list
func getReferenceToNextMemberReference(memberReference *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if memberReference == nil {
		return nil, errors.New("GetNextMemberReference called with nil memberReference")
	}
	if !IsListMemberReference(memberReference, trans) {
		return nil, errors.New("Supplied memberReference is not a refinement of CrlListMemberReference")
	}
	nextMemberReference := memberReference.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToNextMemberReferenceURI, trans)
	if nextMemberReference == nil {
		return nil, errors.New("In GetNextMemberReference, memberReference does not ave a NextMemberReferenceReference")
	}
	return nextMemberReference, nil
}

// getReferenceToPriorMemberReference returns the reference to the previous member of the list. It returns an error if the memberReference
// is either nil or is not a refinement of CrlListMemberReference
// It returns nil if the reference is the first member of the list
func getReferenceToPriorMemberReference(memberReference *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	if memberReference == nil {
		return nil, errors.New("getReferenceToPriorMemberReference called with nil memberReference")
	}
	if !IsListMemberReference(memberReference, trans) {
		return nil, errors.New("In getReferenceToPriorMemberReference, supplied memberReference is not a refinement of CrlListMemberReference")
	}
	priorMemberReference := memberReference.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToPriorMemberReferenceURI, trans)
	if priorMemberReference == nil {
		return nil, errors.New("In getReferenceToPriorMemberReference, memberReference does not have a PriorMemberReferenceReference")
	}
	return priorMemberReference, nil
}

// IsList returns true if the supplied Element is a refinement of List
func IsList(list *core.Concept, trans *core.Transaction) bool {
	return list.IsRefinementOfURI(CrlListURI, trans)
}

// IsListMember returns true if the element is a memeber of the given list
func IsListMember(list *core.Concept, el *core.Concept, trans *core.Transaction) bool {
	uOfD := list.GetUniverseOfDiscourse(trans)
	it := list.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, trans) && memberReference.GetReferencedConcept(trans) == el {
			it.Stop()
			return true
		}
	}
	return false
}

// IsListMemberReference returns true if the supplied Reference is a refinement of ListMemberReference
func IsListMemberReference(memberReference *core.Concept, trans *core.Transaction) bool {
	return memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, trans)
}

// PrependListMember adds a member to the end of the list
func PrependListMember(list *core.Concept, newMember *core.Concept, trans *core.Transaction) (*core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	if !IsList(list, trans) {
		return nil, errors.New("In PrependListMember, Supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("In PrependListMember, Supplied Element is nil: nil members are not allowed in CRL Lists")
	}
	listType, _ := GetListType(list, trans)
	if !newMember.IsRefinementOf(listType, trans) {
		return nil, errors.New("In PrependListMember, NewMember is of wrong type")
	}
	oldFirstMemberReference, err := GetFirstMemberReference(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "PrependListMember failed")
	}
	// Create the newMemberReference
	newMemberReference, err2 := NewListMemberReference(uOfD, trans)
	if err2 != nil {
		return nil, errors.Wrap(err2, "PrependListMember failed")
	}
	err = newMemberReference.SetOwningConcept(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "PrependListMember failed")
	}
	err = newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, trans)
	if err != nil {
		return nil, errors.Wrap(err, "PrependListMember failed")
	}
	// Wire up references - be careful if inserting at the end
	referenceToFirstMemberReference, err3 := getListReferenceToFirstMemberReference(list, trans)
	if err3 != nil {
		return nil, errors.Wrap(err2, "PrependListMember failed")
	}
	if referenceToFirstMemberReference != nil {
		err = referenceToFirstMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, trans)
		if err != nil {
			return nil, errors.Wrap(err, "PrependListMember failed")
		}
	}
	if oldFirstMemberReference == nil {
		referenceToLastMemberReference, err4 := getListReferenceToLastMemberReference(list, trans)
		if err4 != nil {
			return nil, errors.Wrap(err2, "PrependListMember failed")
		}
		err = referenceToLastMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, trans)
		if err != nil {
			return nil, errors.Wrap(err, "PrependListMember failed")
		}
	} else {
		err = setPriorMemberReference(oldFirstMemberReference, newMemberReference, trans)
		if err != nil {
			return nil, errors.Wrap(err, "PrependListMember failed")
		}
	}
	err = setNextMemberReference(newMemberReference, oldFirstMemberReference, trans)
	if err != nil {
		return nil, errors.Wrap(err, "PrependListMember failed")
	}
	return newMemberReference, nil
}

// RemoveListMember removes the first occurrance of an element from the given list
func RemoveListMember(list *core.Concept, el *core.Concept, trans *core.Transaction) error {
	uOfD := list.GetUniverseOfDiscourse(trans)
	it := list.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, trans) && memberReference.GetReferencedConcept(trans) == el {
			// Modify previous and next pointers
			priorMemberReference, _ := GetPriorMemberReference(memberReference, trans)
			nextMemberReference, _ := GetNextMemberReference(memberReference, trans)
			if priorMemberReference != nil {
				setNextMemberReference(priorMemberReference, nextMemberReference, trans)
			} else {
				referenceToFirstMemberReference, _ := getListReferenceToFirstMemberReference(list, trans)
				referenceToFirstMemberReference.SetReferencedConcept(nextMemberReference, core.NoAttribute, trans)
			}
			if nextMemberReference != nil {
				setPriorMemberReference(nextMemberReference, priorMemberReference, trans)
			} else {
				referenceToLastMemberReference, _ := getListReferenceToLastMemberReference(list, trans)
				referenceToLastMemberReference.SetReferencedConcept(priorMemberReference, core.NoAttribute, trans)
			}
			// Now delete the member reference
			uOfD.DeleteElement(memberReference, trans)
			it.Stop()
			return nil
		}
	}
	return errors.New("element not member of list")
}

// SetListType sets the element that should be an abstraction of every member. It is only valid on a list that
// does not already have a list type assigned, i.e. you can't change the type of a list once it has been set.
// It returns an error if the argument is not a list or if the list already has a type assigned
func SetListType(list *core.Concept, listType *core.Concept, trans *core.Transaction) error {
	if !list.IsRefinementOfURI(CrlListURI, trans) {
		return errors.New("Argument is not a list")
	}
	typeReference := list.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, trans)
	if typeReference == nil {
		return errors.New("ListTypeReference not found")
	}
	if typeReference.GetReferencedConcept(trans) != nil {
		return errors.New("List already has an assigned type")
	}
	return typeReference.SetReferencedConcept(listType, core.NoAttribute, trans)
}

// setNextMemberReference takes a memberReference and sets its next reference
func setNextMemberReference(memberReference *core.Concept, nextReference *core.Concept, trans *core.Transaction) error {
	// since this is an internal function we assume that the references are refinements of CrlListMemberReference
	nextReferenceReference, err := getReferenceToNextMemberReference(memberReference, trans)
	if err != nil {
		return errors.Wrap(err, "setNextMemberReference failed")
	}
	err = nextReferenceReference.SetReferencedConcept(nextReference, core.NoAttribute, trans)
	if err != nil {
		return errors.Wrap(err, "setNextMemberReference failed")
	}
	return nil
}

// setPriorMemberReference takes a memberReference and sets its prior reference
func setPriorMemberReference(memberReference *core.Concept, priorReference *core.Concept, trans *core.Transaction) error {
	// since this is an internal function we assume that the references are refinements of CrlListMemberReference
	priorReferenceReference, err := getReferenceToPriorMemberReference(memberReference, trans)
	if err != nil {
		return errors.Wrap(err, "setNextMemberReference failed")
	}
	err = priorReferenceReference.SetReferencedConcept(priorReference, core.NoAttribute, trans)
	if err != nil {
		return errors.Wrap(err, "setNextMemberReference failed")
	}
	return nil
}

// BuildCrlListsConcepts builds the CrlList concept and adds it as a child of the provided parent concept space
func BuildCrlListsConcepts(uOfD *core.UniverseOfDiscourse, parentSpace *core.Concept, trans *core.Transaction) {
	crlList, _ := uOfD.NewElement(trans, CrlListURI)
	crlList.SetLabel("CrlList", trans)
	crlList.SetOwningConcept(parentSpace, trans)

	crlFirstMemberReference, _ := uOfD.NewReference(trans, CrlListReferenceToFirstMemberReferenceURI)
	crlFirstMemberReference.SetLabel("FirstMemberReference", trans)
	crlFirstMemberReference.SetOwningConcept(crlList, trans)

	crlLastMemberReference, _ := uOfD.NewReference(trans, CrlListReferenceToLastMemberReferenceURI)
	crlLastMemberReference.SetLabel("LastMemberReference", trans)
	crlLastMemberReference.SetOwningConcept(crlList, trans)

	CrlListTypeReference, _ := uOfD.NewReference(trans, CrlListTypeReferenceURI)
	CrlListTypeReference.SetLabel("TypeReference", trans)
	CrlListTypeReference.SetOwningConcept(crlList, trans)

	crlListMemberReference, _ := uOfD.NewReference(trans, CrlListMemberReferenceURI)
	crlListMemberReference.SetLabel("MemberReference", trans)
	crlListMemberReference.SetOwningConcept(parentSpace, trans)

	crlNextMemberReference, _ := uOfD.NewReference(trans, CrlListReferenceToNextMemberReferenceURI)
	crlNextMemberReference.SetLabel("NextMemberReference", trans)
	crlNextMemberReference.SetOwningConcept(crlListMemberReference, trans)

	crlPriorMemberReference, _ := uOfD.NewReference(trans, CrlListReferenceToPriorMemberReferenceURI)
	crlPriorMemberReference.SetLabel("PriorMemberReference", trans)
	crlPriorMemberReference.SetOwningConcept(crlListMemberReference, trans)
}
