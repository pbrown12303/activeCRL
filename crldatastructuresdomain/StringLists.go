package crldatastructuresdomain

import (
	"github.com/pkg/errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlStringListURI is the URI that identifies the prototype for list
var CrlStringListURI = CrlDataStructuresDomainURI + "/StringList"

// CrlStringListReferenceToFirstMemberLiteralURI is the URI that identifies the prototype for the first member literal
var CrlStringListReferenceToFirstMemberLiteralURI = CrlStringListURI + "/StringListReferenceToFirstMemberLiteral"

// CrlStringListReferenceToLastMemberLiteralURI is the URI that identifies the prototype for a the last member literal
var CrlStringListReferenceToLastMemberLiteralURI = CrlStringListURI + "/StringListReferenceToLastMemberLiteral"

// CrlStringListMemberLiteralURI is the URI that identifies the prototype for a list member literal
var CrlStringListMemberLiteralURI = CrlStringListURI + "/StringListMemberLiteral"

// CrlStringListReferenceToNextMemberLiteralURI is the URI that identifies a member literal's next member literal
var CrlStringListReferenceToNextMemberLiteralURI = CrlStringListURI + "/ReferenceToNextMemberLiteral"

// CrlStringListReferenceToPriorMemberLiteralURI is the URI that identifies a member literal's prior member literal
var CrlStringListReferenceToPriorMemberLiteralURI = CrlStringListURI + "/ReferenceToPriorMemberLiteral"

// NewStringList creates an instance of a list
func NewStringList(uOfD *core.UniverseOfDiscourse, trans *core.Transaction, newURI ...string) (core.Concept, error) {
	newStringList, err := uOfD.CreateRefinementOfConceptURI(CrlStringListURI, "StringList", trans, newURI...)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlStringListReferenceToFirstMemberLiteralURI, newStringList, "StringListFirstMemberLiteral", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlStringListReferenceToLastMemberLiteralURI, newStringList, "StringListLastMemberLiteral", trans)
	if err != nil {
		return nil, errors.Wrap(err, "StringLists.go NewStringList failed")
	}
	return newStringList, nil
}

// NewStringListMemberLiteral returns a new refinement of CrlStringListMemberLiteral
func NewStringListMemberLiteral(uOfD *core.UniverseOfDiscourse, trans *core.Transaction, newURI ...string) (core.Concept, error) {
	memberLiteral, _ := uOfD.CreateRefinementOfConceptURI(CrlStringListMemberLiteralURI, "StringListMemberLiteral", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlStringListReferenceToNextMemberLiteralURI, memberLiteral, "StringListNextMemberLiteral", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlStringListReferenceToPriorMemberLiteralURI, memberLiteral, "StringListPriorMemberLiteral", trans)
	return memberLiteral, nil
}

// AddStringListMemberAfter adds a string to the list after the priorMemberLiteral and returns the newMemberLiteral.
// If the priorMemberLiteral is nil, the string is added to the beginning of the list. An error is returned if the
// supplied list is not a list, the newMember is the empty string, or the priorMemberLiteral is not a CrlStringListMemberLiteral in this list.
func AddStringListMemberAfter(list core.Concept, priorMemberLiteral core.Concept, newMember string, trans *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	if !IsStringList(list, trans) {
		return nil, errors.New("In AddStringListMemberAfter, supplied Element is not a CRL StringList")
	}
	if newMember == "" {
		return nil, errors.New("In AddStringListMemberAfter, newMember is the emply string: empty strings are not allowed in CRL StringLists")
	}
	// validate prior member literal
	if priorMemberLiteral == nil {
		return nil, errors.New("In AddStringListMemberAfter, priorMemberLiteral is nil: this is a required value")
	}
	if priorMemberLiteral != nil {
		if !priorMemberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans) {
			return nil, errors.New("In AddStringListMemberAfter, supplied priorMemberLiteral is not a CrlStringListMemberLiteral")
		}
		if priorMemberLiteral.GetOwningConcept(trans) != list {
			return nil, errors.New("In AddStringListMemberAfter, supplied priorMemberLiteral does not belong to this list")
		}
	}
	var newPostMemberLiteral core.Concept
	referenceToPostMemberLiteral, err := getReferenceToNextMemberLiteral(priorMemberLiteral, trans)
	if err != nil {
		return nil, err
	}
	referencedPostMemberLiteral := referenceToPostMemberLiteral.GetReferencedConcept(trans)
	if referencedPostMemberLiteral != nil {
		newPostMemberLiteral = referencedPostMemberLiteral
	}
	newMemberLiteral, _ := NewStringListMemberLiteral(uOfD, trans)
	newMemberLiteral.SetOwningConcept(list, trans)
	newMemberLiteral.SetLiteralValue(newMember, trans)
	// Wire up prior references
	setNextMemberLiteral(priorMemberLiteral, newMemberLiteral, trans)

	setPriorMemberLiteral(newMemberLiteral, priorMemberLiteral, trans)
	// Wire up next references
	if newPostMemberLiteral == nil {
		referenceToLastMemberLiteral, err2 := getStringListReferenceToLastMemberLiteral(list, trans)
		if err2 != nil {
			return nil, errors.Wrap(err2, "AddStringListMemberAfter failed")
		}
		if referenceToLastMemberLiteral != nil {
			referenceToLastMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, trans)
		}
	} else {
		setPriorMemberLiteral(newPostMemberLiteral, newMemberLiteral, trans)
		setNextMemberLiteral(newMemberLiteral, newPostMemberLiteral, trans)
	}

	return newMemberLiteral, nil
}

// AddStringListMemberBefore adds a member to the list before the postMember.
// If the postMember is nil, the member is added at the end of the list.
func AddStringListMemberBefore(list core.Concept, postMemberLiteral core.Concept, newMember string, trans *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	if !IsStringList(list, trans) {
		return nil, errors.New("In AddStringListMemberBefore, Supplied Element is not a CRL StringList")
	}
	if newMember == "" {
		return nil, errors.New("In AddStringListMemberBefore, Supplied string is empty: empty strings are not allowed in CRL StringLists")
	}
	// Check to ensure that the postMemberLiteral is valid
	if postMemberLiteral == nil {
		return nil, errors.New("AddStringListMemberBefore called with nil postMemberLiteral")
	}
	if !postMemberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans) {
		return nil, errors.New("In AddStringListMemberBefore, Supplied postMemberLiteral is not a CrlStringListMemberLiteral")
	}
	if postMemberLiteral.GetOwningConcept(trans) != list {
		return nil, errors.New("In AddStringListMemberBefore, Supplied postMemberLiteral does not belong to this list")
	}
	var newPriorMemberLiteral core.Concept
	// If the postMemberLiteral exists, then its priorMemberLiteral should point to the newMemberLiteral
	referenceToPriorMemberLiteral, err := getReferenceToPriorMemberLiteral(postMemberLiteral, trans)
	if err != nil {
		return nil, errors.Wrap(err, "AddStringListMemberBefore failed")
	}
	referencedPriorMemberLiteral := referenceToPriorMemberLiteral.GetReferencedConcept(trans)
	if referencedPriorMemberLiteral != nil {
		newPriorMemberLiteral = referencedPriorMemberLiteral
	}
	// Create the newMemberLiteral
	newMemberLiteral, _ := NewStringListMemberLiteral(uOfD, trans)
	newMemberLiteral.SetOwningConcept(list, trans)
	newMemberLiteral.SetLiteralValue(newMember, trans)
	// Wire up post references - be careful if inserting at the end
	setPriorMemberLiteral(postMemberLiteral, newMemberLiteral, trans)
	setNextMemberLiteral(newMemberLiteral, postMemberLiteral, trans)
	// Wire up prior references
	if newPriorMemberLiteral == nil {
		// The new member is the only member of the list
		referenceToFirstMemberLiteral, _ := getStringListReferenceToFirstMemberLiteral(list, trans)
		if referenceToFirstMemberLiteral != nil {
			referenceToFirstMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, trans)
		}
	} else {
		setNextMemberLiteral(newPriorMemberLiteral, newMemberLiteral, trans)
		setPriorMemberLiteral(newMemberLiteral, newPriorMemberLiteral, trans)
	}
	return newMemberLiteral, nil
}

// AppendStringListMember adds a string to the end of the list
func AppendStringListMember(list core.Concept, value string, trans *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	if !IsStringList(list, trans) {
		return nil, errors.New("In AddStringListMemberBefore, Supplied Element is not a CRL StringList")
	}
	if value == "" {
		return nil, errors.New("In AddStringListMemberBefore, Supplied string is empty: empty strings are not allowed in CRL StringLists")
	}
	oldLastMemberLiteral, err := GetLastMemberLiteral(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "AppendStringListMember failed")
	}
	// Create the newMemberLiteral
	newMemberLiteral, err2 := NewStringListMemberLiteral(uOfD, trans)
	if err2 != nil {
		return nil, errors.Wrap(err2, "AppendStringListMember failed")
	}
	err = newMemberLiteral.SetOwningConcept(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "AppendStringListMember failed")
	}
	err = newMemberLiteral.SetLiteralValue(value, trans)
	if err != nil {
		return nil, errors.Wrap(err, "AppendStringListMember failed")
	}
	// Wire up references - be careful if inserting at the end
	referenceToLastMemberLiteral, err3 := getStringListReferenceToLastMemberLiteral(list, trans)
	if err3 != nil {
		return nil, errors.Wrap(err2, "AppendStringListMember failed")
	}
	if referenceToLastMemberLiteral != nil {
		err = referenceToLastMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, trans)
		if err != nil {
			return nil, errors.Wrap(err, "AppendStringListMember failed")
		}
	}
	if oldLastMemberLiteral == nil {
		referenceToFirstMemberLiteral, err4 := getStringListReferenceToFirstMemberLiteral(list, trans)
		if err4 != nil {
			return nil, errors.Wrap(err2, "AppendStringListMember failed")
		}
		err = referenceToFirstMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, trans)
		if err != nil {
			return nil, errors.Wrap(err, "AppendStringListMember failed")
		}
	} else {
		err = setNextMemberLiteral(oldLastMemberLiteral, newMemberLiteral, trans)
		if err != nil {
			return nil, errors.Wrap(err, "AppendStringListMember failed")
		}
		err = setPriorMemberLiteral(newMemberLiteral, oldLastMemberLiteral, trans)
		if err != nil {
			return nil, errors.Wrap(err, "AppendStringListMember failed")
		}
	}
	return newMemberLiteral, nil
}

// ClearStringList removes all members from the list
func ClearStringList(list core.Concept, trans *core.Transaction) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	it := list.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberLiteral := uOfD.GetLiteral(id.(string))
		if memberLiteral != nil && memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans) {
			uOfD.DeleteElement(memberLiteral, trans)
		}
	}
}

// GetFirstMemberLiteral returns the reference to the first member of the list. It returns an error if the
// list is not a list. It returns nil if the list is empty
func GetFirstMemberLiteral(list core.Concept, trans *core.Transaction) (core.Concept, error) {
	refRef, err := getStringListReferenceToFirstMemberLiteral(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "GetFirstMemberLiteral failed")
	}
	if refRef == nil {
		return nil, errors.New("In StringList GetFirstMemberLiteral, No reference to first member literal found")
	}
	firstMemberLiteral := refRef.GetReferencedConcept(trans)
	if firstMemberLiteral == nil {
		return nil, nil
	}
	return firstMemberLiteral, nil
}

// GetFirstLiteralForString returns the first Literal whose value is the given string. It returns an error if the list is not a list.
// It returns nil if the string it is not found in the list.
func GetFirstLiteralForString(list core.Concept, value string, trans *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	it := list.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberLiteral := uOfD.GetLiteral(id.(string))
		if memberLiteral != nil &&
			memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans) &&
			memberLiteral.GetLiteralValue(trans) == value {
			it.Stop()
			return memberLiteral, nil
		}
	}
	return nil, nil
}

// GetLastMemberLiteral returns the reference to the last member of the list. It returns an error if list is not a list.
// It returns nil if the list is empty
func GetLastMemberLiteral(list core.Concept, trans *core.Transaction) (core.Concept, error) {
	refRef, err := getStringListReferenceToLastMemberLiteral(list, trans)
	if err != nil {
		return nil, err
	}
	if refRef == nil {
		return nil, errors.New("No reference to last member literal found")
	}
	lastMemberLiteral := refRef.GetReferencedConcept(trans)
	if lastMemberLiteral == nil {
		return nil, nil
	}
	return lastMemberLiteral, nil
}

// getStringListReferenceToFirstMemberLiteral returns the reference to the first member literal. It returns an error if list is not a StringList
func getStringListReferenceToFirstMemberLiteral(list core.Concept, trans *core.Transaction) (core.Concept, error) {
	if !IsStringList(list, trans) {
		return nil, errors.New("Argument is not a CrlDataStructures.StringList")
	}
	refToLiteral := list.GetFirstOwnedReferenceRefinedFromURI(CrlStringListReferenceToFirstMemberLiteralURI, trans)
	if refToLiteral == nil {
		return nil, errors.New("In getStringListReferenceToFirstMemberLiteral, the reference was not found")
	}
	return refToLiteral, nil
}

// getStringListReferenceToLastMemberLiteral returns the reference to the last member literal. It returns an error if list is not a StringList
func getStringListReferenceToLastMemberLiteral(list core.Concept, trans *core.Transaction) (core.Concept, error) {
	if !IsStringList(list, trans) {
		return nil, errors.New("Argument is not a CrlDataStructures.StringList")
	}
	return list.GetFirstOwnedReferenceRefinedFromURI(CrlStringListReferenceToLastMemberLiteralURI, trans), nil
}

// GetNextMemberLiteral returns the successor member literal in the list
func GetNextMemberLiteral(memberLiteral core.Concept, trans *core.Transaction) (core.Concept, error) {
	if !IsStringListMemberLiteral(memberLiteral, trans) {
		return nil, errors.New("Supplied memberLiteral is not a refinement of CrlStringListMemberLiteral")
	}
	referenceToNextMemberLiteral, err := getReferenceToNextMemberLiteral(memberLiteral, trans)
	if err != nil {
		return nil, err
	}
	nextMemberLiteral := referenceToNextMemberLiteral.GetReferencedConcept(trans)
	if nextMemberLiteral == nil {
		return nil, nil
	}
	return nextMemberLiteral, nil
}

// GetPriorMemberLiteral returns the predecessor member literal in the list
func GetPriorMemberLiteral(memberLiteral core.Concept, trans *core.Transaction) (core.Concept, error) {
	if !IsStringListMemberLiteral(memberLiteral, trans) {
		return nil, errors.New("Supplied memberLiteral is not a refinement of CrlStringListMemberLiteral")
	}
	referenceToPriorMemberLiteral, err := getReferenceToPriorMemberLiteral(memberLiteral, trans)
	if err != nil {
		return nil, err
	}
	priorMemberLiteral := referenceToPriorMemberLiteral.GetReferencedConcept(trans)
	if priorMemberLiteral == nil {
		return nil, nil
	}
	return priorMemberLiteral.(core.Concept), nil
}

// getReferenceToNextMemberLiteral returns the reference to the next member of the list.
// It returns nil if the reference is the last member of the list
func getReferenceToNextMemberLiteral(memberLiteral core.Concept, trans *core.Transaction) (core.Concept, error) {
	if memberLiteral == nil {
		return nil, errors.New("GetNextMemberLiteral called with nil memberLiteral")
	}
	if !IsStringListMemberLiteral(memberLiteral, trans) {
		return nil, errors.New("Supplied memberLiteral is not a refinement of CrlStringListMemberLiteral")
	}
	nextMemberLiteral := memberLiteral.GetFirstOwnedReferenceRefinedFromURI(CrlStringListReferenceToNextMemberLiteralURI, trans)
	if nextMemberLiteral == nil {
		return nil, errors.New("In GetNextMemberLiteral, memberLiteral does not ave a NextMemberLiteralReference")
	}
	return nextMemberLiteral, nil
}

// getReferenceToPriorMemberLiteral returns the reference to the previous member of the list. It returns an error if the memberLiteral
// is either nil or is not a refinement of CrlStringListMemberLiteral
// It returns nil if the reference is the first member of the list
func getReferenceToPriorMemberLiteral(memberLiteral core.Concept, trans *core.Transaction) (core.Concept, error) {
	if memberLiteral == nil {
		return nil, errors.New("getReferenceToPriorMemberLiteral called with nil memberLiteral")
	}
	if !IsStringListMemberLiteral(memberLiteral, trans) {
		return nil, errors.New("In getReferenceToPriorMemberLiteral, supplied memberLiteral is not a refinement of CrlStringListMemberLiteral")
	}
	priorMemberLiteralReference := memberLiteral.GetFirstOwnedReferenceRefinedFromURI(CrlStringListReferenceToPriorMemberLiteralURI, trans)
	if priorMemberLiteralReference == nil {
		return nil, errors.New("In getReferenceToPriorMemberLiteral, memberLiteral does not have a PriorMemberLiteralReference")
	}
	return priorMemberLiteralReference, nil
}

// IsStringList returns true if the supplied Element is a refinement of StringList
func IsStringList(list core.Concept, trans *core.Transaction) bool {
	return list.IsRefinementOfURI(CrlStringListURI, trans)
}

// IsStringListMember returns true if the string is a memeber of the given list
func IsStringListMember(list core.Concept, value string, trans *core.Transaction) bool {
	uOfD := list.GetUniverseOfDiscourse(trans)
	it := list.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberLiteral := uOfD.GetLiteral(id.(string))
		if memberLiteral != nil && memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans) && memberLiteral.GetLiteralValue(trans) == value {
			it.Stop()
			return true
		}
	}
	return false
}

// IsStringListMemberLiteral returns true if the supplied Reference is a refinement of StringListMemberLiteral
func IsStringListMemberLiteral(memberLiteral core.Concept, trans *core.Transaction) bool {
	return memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans)
}

// PrependStringListMember adds a string to the beginning of the list
func PrependStringListMember(list core.Concept, value string, trans *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(trans)
	if !IsStringList(list, trans) {
		return nil, errors.New("In PrependStringListMember, Supplied Element is not a CRL StringList")
	}
	if value == "" {
		return nil, errors.New("In PrependStringListMember, Supplied string is empty: empty strings are not allowed in CRL StringLists")
	}
	oldFirstMemberLiteral, err := GetFirstMemberLiteral(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "PrependStringListMember failed")
	}
	// Create the newMemberLiteral
	newMemberLiteral, err2 := NewStringListMemberLiteral(uOfD, trans)
	if err2 != nil {
		return nil, errors.Wrap(err2, "PrependStringListMember failed")
	}
	err = newMemberLiteral.SetOwningConcept(list, trans)
	if err != nil {
		return nil, errors.Wrap(err, "PrependStringListMember failed")
	}
	err = newMemberLiteral.SetLiteralValue(value, trans)
	if err != nil {
		return nil, errors.Wrap(err, "PrependStringListMember failed")
	}
	// Wire up references - be careful if inserting at the end
	referenceToFirstMemberLiteral, err3 := getStringListReferenceToFirstMemberLiteral(list, trans)
	if err3 != nil {
		return nil, errors.Wrap(err2, "PrependStringListMember failed")
	}
	if referenceToFirstMemberLiteral != nil {
		err = referenceToFirstMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, trans)
		if err != nil {
			return nil, errors.Wrap(err, "PrependStringListMember failed")
		}
	}
	if oldFirstMemberLiteral == nil {
		referenceToLastMemberLiteral, err4 := getStringListReferenceToLastMemberLiteral(list, trans)
		if err4 != nil {
			return nil, errors.Wrap(err2, "PrependStringListMember failed")
		}
		err = referenceToLastMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, trans)
		if err != nil {
			return nil, errors.Wrap(err, "PrependStringListMember failed")
		}
	} else {
		err = setPriorMemberLiteral(oldFirstMemberLiteral, newMemberLiteral, trans)
		if err != nil {
			return nil, errors.Wrap(err, "PrependStringListMember failed")
		}
	}
	err = setNextMemberLiteral(newMemberLiteral, oldFirstMemberLiteral, trans)
	if err != nil {
		return nil, errors.Wrap(err, "PrependStringListMember failed")
	}
	return newMemberLiteral, nil
}

// RemoveStringListMember removes the first occurrance of an element from the given list
func RemoveStringListMember(list core.Concept, value string, trans *core.Transaction) error {
	uOfD := list.GetUniverseOfDiscourse(trans)
	it := list.GetOwnedConceptIDs(trans).Iterator()
	for id := range it.C {
		memberLiteral := uOfD.GetLiteral(id.(string))
		if memberLiteral != nil && memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, trans) && memberLiteral.GetLiteralValue(trans) == value {
			// Modify previous and next pointers
			priorMemberLiteral, _ := GetPriorMemberLiteral(memberLiteral, trans)
			nextMemberLiteral, _ := GetNextMemberLiteral(memberLiteral, trans)
			if priorMemberLiteral != nil {
				setNextMemberLiteral(priorMemberLiteral, nextMemberLiteral, trans)
			} else {
				referenceToFirstMemberLiteral, _ := getStringListReferenceToFirstMemberLiteral(list, trans)
				referenceToFirstMemberLiteral.SetReferencedConcept(nextMemberLiteral, core.NoAttribute, trans)
			}
			if nextMemberLiteral != nil {
				setPriorMemberLiteral(nextMemberLiteral, priorMemberLiteral, trans)
			} else {
				referenceToLastMemberLiteral, _ := getStringListReferenceToLastMemberLiteral(list, trans)
				referenceToLastMemberLiteral.SetReferencedConcept(priorMemberLiteral, core.NoAttribute, trans)
			}
			// Now delete the member literal
			uOfD.DeleteElement(memberLiteral, trans)
			it.Stop()
			return nil
		}
	}
	return errors.New("element not member of list")
}

// setNextMemberLiteral takes a memberLiteral and sets its next reference
func setNextMemberLiteral(memberLiteral core.Concept, nextLiteral core.Concept, trans *core.Transaction) error {
	// since this is an internal function we assume that the references are refinements of CrlStringListMemberLiteral
	nextLiteralReference, err := getReferenceToNextMemberLiteral(memberLiteral, trans)
	if err != nil {
		return errors.Wrap(err, "setNextMemberLiteral failed")
	}
	err = nextLiteralReference.SetReferencedConcept(nextLiteral, core.NoAttribute, trans)
	if err != nil {
		return errors.Wrap(err, "setNextMemberLiteral failed")
	}
	return nil
}

// setPriorMemberLiteral takes a memberLiteral and sets its prior reference
func setPriorMemberLiteral(memberLiteral core.Concept, priorLiteral core.Concept, trans *core.Transaction) error {
	// since this is an internal function we assume that the references are refinements of CrlStringListMemberLiteral
	priorLiteralReference, err := getReferenceToPriorMemberLiteral(memberLiteral, trans)
	if err != nil {
		return errors.Wrap(err, "setNextMemberLiteral failed")
	}
	err = priorLiteralReference.SetReferencedConcept(priorLiteral, core.NoAttribute, trans)
	if err != nil {
		return errors.Wrap(err, "setNextMemberLiteral failed")
	}
	return nil
}

// BuildCrlStringListsConcepts builds the CrlStringList concept and adds it as a child of the provided parent concept space
func BuildCrlStringListsConcepts(uOfD *core.UniverseOfDiscourse, parentSpace core.Concept, trans *core.Transaction) {
	crlStringList, _ := uOfD.NewElement(trans, CrlStringListURI)
	crlStringList.SetLabel("CrlStringList", trans)
	crlStringList.SetOwningConcept(parentSpace, trans)

	crlFirstMemberLiteral, _ := uOfD.NewReference(trans, CrlStringListReferenceToFirstMemberLiteralURI)
	crlFirstMemberLiteral.SetLabel("StringListFirstMemberLiteral", trans)
	crlFirstMemberLiteral.SetOwningConcept(crlStringList, trans)

	crlLastMemberLiteral, _ := uOfD.NewReference(trans, CrlStringListReferenceToLastMemberLiteralURI)
	crlLastMemberLiteral.SetLabel("StringListLastMemberLiteral", trans)
	crlLastMemberLiteral.SetOwningConcept(crlStringList, trans)

	crlStringListMemberLiteral, _ := uOfD.NewLiteral(trans, CrlStringListMemberLiteralURI)
	crlStringListMemberLiteral.SetLabel("StringListMemberLiteral", trans)
	crlStringListMemberLiteral.SetOwningConcept(parentSpace, trans)

	crlNextMemberLiteral, _ := uOfD.NewReference(trans, CrlStringListReferenceToNextMemberLiteralURI)
	crlNextMemberLiteral.SetLabel("StringListNextMemberLiteral", trans)
	crlNextMemberLiteral.SetOwningConcept(crlStringListMemberLiteral, trans)

	crlPriorMemberLiteral, _ := uOfD.NewReference(trans, CrlStringListReferenceToPriorMemberLiteralURI)
	crlPriorMemberLiteral.SetLabel("StringListPriorMemberLiteral", trans)
	crlPriorMemberLiteral.SetOwningConcept(crlStringListMemberLiteral, trans)
}
