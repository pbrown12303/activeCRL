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
func NewList(uOfD *core.UniverseOfDiscourse, setType core.Element, hl *core.Transaction, newURI ...string) (core.Element, error) {
	if setType == nil {
		return nil, errors.New("No type specified for list")
	}
	newList, err := uOfD.CreateReplicateAsRefinementFromURI(CrlListURI, hl, newURI...)
	if err != nil {
		return nil, err
	}
	typeReference := newList.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
	if typeReference == nil {
		return nil, errors.New("In Lists.go, NewList failed to find a type reference")
	}
	typeReference.SetReferencedConcept(setType, core.NoAttribute, hl)
	return newList, nil
}

// AddListMemberAfter adds a member to the list after the priorMemberReference and returns the newMemberReference.
// If the priorMemberReference is nil, the member is added to the beginning of the list. An error is returned if the
// supplied list is not a list, the newElement is nil, or the priorElementReference is not a CrlListMemberReference in this list.
func AddListMemberAfter(list core.Element, priorMemberReference core.Reference, newMember core.Element, hl *core.Transaction) (core.Reference, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsList(list, hl) {
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
		if !priorMemberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) {
			return nil, errors.New("In AddListMemberAfter, supplied priorElementReference is not a CrlListMemberReference")
		}
		if priorMemberReference.GetOwningConcept(hl) != list {
			return nil, errors.New("In AddListMemberAfter, supplied priorMemberReference does not belong to this list")
		}
	}
	listType, _ := GetListType(list, hl)
	if newMember.IsRefinementOf(listType, hl) == false {
		return nil, errors.New("In AddListMemberAfter, newMember is of wrong type")
	}
	var newPostMemberReference core.Reference
	referenceToPostMemberReference, err := getReferenceToNextMemberReference(priorMemberReference, hl)
	if err != nil {
		return nil, err
	}
	referencedPostMemberReference := referenceToPostMemberReference.GetReferencedConcept(hl)
	if referencedPostMemberReference != nil {
		newPostMemberReference = referencedPostMemberReference.(core.Reference)
	}
	newMemberReference, _ := uOfD.CreateReplicateReferenceAsRefinementFromURI(CrlListMemberReferenceURI, hl)
	newMemberReference.SetOwningConcept(list, hl)
	newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, hl)
	// Wire up prior references
	setNextMemberReference(priorMemberReference, newMemberReference, hl)

	setPriorMemberReference(newMemberReference, priorMemberReference, hl)
	// Wire up next references
	if newPostMemberReference == nil {
		referenceToLastMemberReference, err2 := getListReferenceToLastMemberReference(list, hl)
		if err2 != nil {
			return nil, errors.Wrap(err2, "AddListMemberAfter failed")
		}
		if referenceToLastMemberReference != nil {
			referenceToLastMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, hl)
		}
	} else {
		setPriorMemberReference(newPostMemberReference, newMemberReference, hl)
		setNextMemberReference(newMemberReference, newPostMemberReference, hl)
	}

	return newMemberReference, nil
}

// AddListMemberBefore adds a member to the list before the postMember.
// If the postMember is nil, the member is added at the end of the list.
func AddListMemberBefore(list core.Element, postMemberReference core.Reference, newMember core.Element, hl *core.Transaction) (core.Reference, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsList(list, hl) {
		return nil, errors.New("In AddListMemberBefore, Supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("In AddListMemberBefore, Supplied Element is nil: nil members are not allowed in CRL Lists")
	}
	// Check to ensure that the postMemberReference is valid
	if postMemberReference == nil {
		return nil, errors.New("AddListMemberBefore called with nil postMemberReference")
	}
	if !postMemberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) {
		return nil, errors.New("In AddListMemberBefore, Supplied postMemberReference is not a CrlListMemberReference")
	}
	if postMemberReference.GetOwningConcept(hl) != list {
		return nil, errors.New("In AddListMemberBefore, Supplied postMemberReference does not belong to this list")
	}
	listType, _ := GetListType(list, hl)
	if newMember.IsRefinementOf(listType, hl) == false {
		return nil, errors.New("In AddListMemberBefore, NewMember is of wrong type")
	}
	var newPriorMemberReference core.Reference
	// If the postMemberReference exists, then its priorMemberReference should point to the newMemberReference
	referenceToPriorMemberReference, err := getReferenceToPriorMemberReference(postMemberReference, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AddListMemberBefore failed")
	}
	referencedPriorMemberReference := referenceToPriorMemberReference.GetReferencedConcept(hl)
	if referencedPriorMemberReference != nil {
		newPriorMemberReference = referencedPriorMemberReference.(core.Reference)
	}
	// Create the newMemberReference
	newMemberReference, _ := uOfD.CreateReplicateReferenceAsRefinementFromURI(CrlListMemberReferenceURI, hl)
	newMemberReference.SetOwningConcept(list, hl)
	newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, hl)
	// Wire up post references - be careful if inserting at the end
	setPriorMemberReference(postMemberReference, newMemberReference, hl)
	setNextMemberReference(newMemberReference, postMemberReference, hl)
	// Wire up prior references
	if newPriorMemberReference == nil {
		// The new member is the only member of the list
		referenceToFirstMemberReference, _ := getListReferenceToFirstMemberReference(list, hl)
		if referenceToFirstMemberReference != nil {
			referenceToFirstMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, hl)
		}
	} else {
		setNextMemberReference(newPriorMemberReference, newMemberReference, hl)
		setPriorMemberReference(newMemberReference, newPriorMemberReference, hl)
	}
	return newMemberReference, nil
}

// AppendListMember adds a member to the end of the list
func AppendListMember(list core.Element, newMember core.Element, hl *core.Transaction) (core.Reference, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsList(list, hl) {
		return nil, errors.New("In AddListMemberBefore, Supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("In AddListMemberBefore, Supplied Element is nil: nil members are not allowed in CRL Lists")
	}
	listType, _ := GetListType(list, hl)
	if newMember.IsRefinementOf(listType, hl) == false {
		return nil, errors.New("In AddListMemberBefore, NewMember is of wrong type")
	}
	oldLastMemberReference, err := GetLastMemberReference(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AppendListMember failed")
	}
	// Create the newMemberReference
	newMemberReference, err2 := uOfD.CreateReplicateReferenceAsRefinementFromURI(CrlListMemberReferenceURI, hl)
	if err2 != nil {
		return nil, errors.Wrap(err2, "AppendListMember failed")
	}
	err = newMemberReference.SetOwningConcept(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AppendListMember failed")
	}
	err = newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AppendListMember failed")
	}
	// Wire up references - be careful if inserting at the end
	referenceToLastMemberReference, err3 := getListReferenceToLastMemberReference(list, hl)
	if err3 != nil {
		return nil, errors.Wrap(err2, "AppendListMember failed")
	}
	if referenceToLastMemberReference != nil {
		err = referenceToLastMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, hl)
		if err != nil {
			return nil, errors.Wrap(err, "AppendListMember failed")
		}
	}
	if oldLastMemberReference == nil {
		referenceToFirstMemberReference, err4 := getListReferenceToFirstMemberReference(list, hl)
		if err4 != nil {
			return nil, errors.Wrap(err2, "AppendListMember failed")
		}
		err = referenceToFirstMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, hl)
		if err != nil {
			return nil, errors.Wrap(err, "AppendListMember failed")
		}
	} else {
		err = setNextMemberReference(oldLastMemberReference, newMemberReference, hl)
		if err != nil {
			return nil, errors.Wrap(err, "AppendListMember failed")
		}
	}
	err = setPriorMemberReference(newMemberReference, oldLastMemberReference, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AppendListMember failed")
	}
	return newMemberReference, nil
}

// ClearList removes all members from the list
func ClearList(list core.Element, hl *core.Transaction) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	it := list.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) {
			uOfD.DeleteElement(memberReference, hl)
		}
	}
}

// GetFirstMemberReference returns the reference to the first member of the list. It returns an error if the
// list is not a list. It returns nil if the list is empty
func GetFirstMemberReference(list core.Element, hl *core.Transaction) (core.Reference, error) {
	refRef, err := getListReferenceToFirstMemberReference(list, hl)
	if err != nil {
		return nil, err
	}
	if refRef == nil {
		return nil, errors.New("In List GetFirstMemberReference, No reference to first member reference found")
	}
	firstMemberReference := refRef.GetReferencedConcept(hl)
	if firstMemberReference == nil {
		return nil, nil
	}
	return firstMemberReference.(core.Reference), nil
}

// GetFirstReferenceForMember returns the first reference to the given member. It returns an error if the list is not a list.
// It returns nil if the element it is not found in the list.
func GetFirstReferenceForMember(list core.Element, member core.Element, hl *core.Transaction) (core.Reference, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	it := list.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	for id := range it.C {
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
func GetLastMemberReference(list core.Element, hl *core.Transaction) (core.Reference, error) {
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
func getListReferenceToFirstMemberReference(list core.Element, hl *core.Transaction) (core.Reference, error) {
	if IsList(list, hl) == false {
		return nil, errors.New("Argument is not a CrlDataStructures.List")
	}
	refToRef := list.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToFirstMemberReferenceURI, hl)
	if refToRef == nil {
		return nil, errors.New("In getListReferenceToFirstMemberReference, the reference was not found")
	}
	return refToRef, nil
}

// getListReferenceToLastMemberReference returns the reference to the last member reference. It returns an error if list is not a List
func getListReferenceToLastMemberReference(list core.Element, hl *core.Transaction) (core.Reference, error) {
	if IsList(list, hl) == false {
		return nil, errors.New("Argument is not a CrlDataStructures.List")
	}
	return list.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToLastMemberReferenceURI, hl), nil
}

// GetListType returns the element that should be an abstraction of every member. It returns an error if the argument is not a list
func GetListType(list core.Element, hl *core.Transaction) (core.Element, error) {
	if list.IsRefinementOfURI(CrlListURI, hl) == false {
		return nil, errors.New("Argument is not a list")
	}
	typeReference := list.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
	return typeReference.GetReferencedConcept(hl), nil
}

// GetNextMemberReference returns the successor member reference in the list
func GetNextMemberReference(memberReference core.Reference, hl *core.Transaction) (core.Reference, error) {
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
func GetPriorMemberReference(memberReference core.Reference, hl *core.Transaction) (core.Reference, error) {
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
func getReferenceToNextMemberReference(memberReference core.Reference, hl *core.Transaction) (core.Reference, error) {
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
func getReferenceToPriorMemberReference(memberReference core.Reference, hl *core.Transaction) (core.Reference, error) {
	if memberReference == nil {
		return nil, errors.New("getReferenceToPriorMemberReference called with nil memberReference")
	}
	if !IsListMemberReference(memberReference, hl) {
		return nil, errors.New("In getReferenceToPriorMemberReference, supplied memberReference is not a refinement of CrlListMemberReference")
	}
	priorMemberReference := memberReference.GetFirstOwnedReferenceRefinedFromURI(CrlListReferenceToPriorMemberReferenceURI, hl)
	if priorMemberReference == nil {
		return nil, errors.New("In getReferenceToPriorMemberReference, memberReference does not have a PriorMemberReferenceReference")
	}
	return priorMemberReference, nil
}

// IsList returns true if the supplied Element is a refinement of List
func IsList(list core.Element, hl *core.Transaction) bool {
	return list.IsRefinementOfURI(CrlListURI, hl)
}

// IsListMember returns true if the element is a memeber of the given list
func IsListMember(list core.Element, el core.Element, hl *core.Transaction) bool {
	uOfD := list.GetUniverseOfDiscourse(hl)
	it := list.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) && memberReference.GetReferencedConcept(hl) == el {
			return true
		}
	}
	return false
}

// IsListMemberReference returns true if the supplied Reference is a refinement of ListMemberReference
func IsListMemberReference(memberReference core.Reference, hl *core.Transaction) bool {
	return memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl)
}

// PrependListMember adds a member to the end of the list
func PrependListMember(list core.Element, newMember core.Element, hl *core.Transaction) (core.Reference, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsList(list, hl) {
		return nil, errors.New("In PrependListMember, Supplied Element is not a CRL List")
	}
	if newMember == nil {
		return nil, errors.New("In PrependListMember, Supplied Element is nil: nil members are not allowed in CRL Lists")
	}
	listType, _ := GetListType(list, hl)
	if newMember.IsRefinementOf(listType, hl) == false {
		return nil, errors.New("In PrependListMember, NewMember is of wrong type")
	}
	oldFirstMemberReference, err := GetFirstMemberReference(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "PrependListMember failed")
	}
	// Create the newMemberReference
	newMemberReference, err2 := uOfD.CreateReplicateReferenceAsRefinementFromURI(CrlListMemberReferenceURI, hl)
	if err2 != nil {
		return nil, errors.Wrap(err2, "PrependListMember failed")
	}
	err = newMemberReference.SetOwningConcept(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "PrependListMember failed")
	}
	err = newMemberReference.SetReferencedConcept(newMember, core.NoAttribute, hl)
	if err != nil {
		return nil, errors.Wrap(err, "PrependListMember failed")
	}
	// Wire up references - be careful if inserting at the end
	referenceToFirstMemberReference, err3 := getListReferenceToFirstMemberReference(list, hl)
	if err3 != nil {
		return nil, errors.Wrap(err2, "PrependListMember failed")
	}
	if referenceToFirstMemberReference != nil {
		err = referenceToFirstMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, hl)
		if err != nil {
			return nil, errors.Wrap(err, "PrependListMember failed")
		}
	}
	if oldFirstMemberReference == nil {
		referenceToLastMemberReference, err4 := getListReferenceToLastMemberReference(list, hl)
		if err4 != nil {
			return nil, errors.Wrap(err2, "PrependListMember failed")
		}
		err = referenceToLastMemberReference.SetReferencedConcept(newMemberReference, core.NoAttribute, hl)
		if err != nil {
			return nil, errors.Wrap(err, "PrependListMember failed")
		}
	} else {
		err = setPriorMemberReference(oldFirstMemberReference, newMemberReference, hl)
		if err != nil {
			return nil, errors.Wrap(err, "PrependListMember failed")
		}
	}
	err = setNextMemberReference(newMemberReference, oldFirstMemberReference, hl)
	if err != nil {
		return nil, errors.Wrap(err, "PrependListMember failed")
	}
	return newMemberReference, nil
}

// RemoveListMember removes the first occurrance of an element from the given list
func RemoveListMember(list core.Element, el core.Element, hl *core.Transaction) error {
	uOfD := list.GetUniverseOfDiscourse(hl)
	it := list.GetOwnedConceptIDs(hl).Iterator()
	defer it.Stop()
	for id := range it.C {
		memberReference := uOfD.GetReference(id.(string))
		if memberReference != nil && memberReference.IsRefinementOfURI(CrlListMemberReferenceURI, hl) && memberReference.GetReferencedConcept(hl) == el {
			// Modify previous and next pointers
			priorMemberReference, _ := GetPriorMemberReference(memberReference, hl)
			nextMemberReference, _ := GetNextMemberReference(memberReference, hl)
			if priorMemberReference != nil {
				setNextMemberReference(priorMemberReference, nextMemberReference, hl)
			} else {
				referenceToFirstMemberReference, _ := getListReferenceToFirstMemberReference(list, hl)
				referenceToFirstMemberReference.SetReferencedConcept(nextMemberReference, core.NoAttribute, hl)
			}
			if nextMemberReference != nil {
				setPriorMemberReference(nextMemberReference, priorMemberReference, hl)
			} else {
				referenceToLastMemberReference, _ := getListReferenceToLastMemberReference(list, hl)
				referenceToLastMemberReference.SetReferencedConcept(priorMemberReference, core.NoAttribute, hl)
			}
			// Now delete the member reference
			uOfD.DeleteElement(memberReference, hl)
			return nil
		}
	}
	return errors.New("element not member of list")
}

// SetListType sets the element that should be an abstraction of every member. It is only valid on a list that
// does not already have a list type assigned, i.e. you can't change the type of a list once it has been set.
// It returns an error if the argument is not a list or if the list already has a type assigned
func SetListType(list core.Element, listType core.Element, hl *core.Transaction) error {
	if list.IsRefinementOfURI(CrlListURI, hl) == false {
		return errors.New("Argument is not a list")
	}
	typeReference := list.GetFirstOwnedReferenceRefinedFromURI(CrlListTypeReferenceURI, hl)
	if typeReference == nil {
		return errors.New("ListTypeReference not found")
	}
	if typeReference.GetReferencedConcept(hl) != nil {
		return errors.New("List already has an assigned type")
	}
	return typeReference.SetReferencedConcept(listType, core.NoAttribute, hl)
}

// setNextMemberReference takes a memberReference and sets its next reference
func setNextMemberReference(memberReference core.Reference, nextReference core.Reference, hl *core.Transaction) error {
	// since this is an internal function we assume that the references are refinements of CrlListMemberReference
	nextReferenceReference, err := getReferenceToNextMemberReference(memberReference, hl)
	if err != nil {
		return errors.Wrap(err, "setNextMemberReference failed")
	}
	err = nextReferenceReference.SetReferencedConcept(nextReference, core.NoAttribute, hl)
	if err != nil {
		return errors.Wrap(err, "setNextMemberReference failed")
	}
	return nil
}

// setPriorMemberReference takes a memberReference and sets its prior reference
func setPriorMemberReference(memberReference core.Reference, priorReference core.Reference, hl *core.Transaction) error {
	// since this is an internal function we assume that the references are refinements of CrlListMemberReference
	priorReferenceReference, err := getReferenceToPriorMemberReference(memberReference, hl)
	if err != nil {
		return errors.Wrap(err, "setNextMemberReference failed")
	}
	err = priorReferenceReference.SetReferencedConcept(priorReference, core.NoAttribute, hl)
	if err != nil {
		return errors.Wrap(err, "setNextMemberReference failed")
	}
	return nil
}

// BuildCrlListsConcepts builds the CrlList concept and adds it as a child of the provided parent concept space
func BuildCrlListsConcepts(uOfD *core.UniverseOfDiscourse, parentSpace core.Element, hl *core.Transaction) {
	crlList, _ := uOfD.NewElement(hl, CrlListURI)
	crlList.SetLabel("CrlList", hl)
	crlList.SetOwningConcept(parentSpace, hl)

	crlFirstMemberReference, _ := uOfD.NewReference(hl, CrlListReferenceToFirstMemberReferenceURI)
	crlFirstMemberReference.SetLabel("FirstMemberReference", hl)
	crlFirstMemberReference.SetOwningConcept(crlList, hl)

	crlLastMemberReference, _ := uOfD.NewReference(hl, CrlListReferenceToLastMemberReferenceURI)
	crlLastMemberReference.SetLabel("LastMemberReference", hl)
	crlLastMemberReference.SetOwningConcept(crlList, hl)

	CrlListTypeReference, _ := uOfD.NewReference(hl, CrlListTypeReferenceURI)
	CrlListTypeReference.SetLabel("TypeReference", hl)
	CrlListTypeReference.SetOwningConcept(crlList, hl)

	crlListMemberReference, _ := uOfD.NewReference(hl, CrlListMemberReferenceURI)
	crlListMemberReference.SetLabel("MemberReference", hl)
	crlListMemberReference.SetOwningConcept(parentSpace, hl)

	crlNextMemberReference, _ := uOfD.NewReference(hl, CrlListReferenceToNextMemberReferenceURI)
	crlNextMemberReference.SetLabel("NextMemberReference", hl)
	crlNextMemberReference.SetOwningConcept(crlListMemberReference, hl)

	crlPriorMemberReference, _ := uOfD.NewReference(hl, CrlListReferenceToPriorMemberReferenceURI)
	crlPriorMemberReference.SetLabel("PriorMemberReference", hl)
	crlPriorMemberReference.SetOwningConcept(crlListMemberReference, hl)
}
