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
func NewStringList(uOfD *core.UniverseOfDiscourse, hl *core.Transaction, newURI ...string) (core.Concept, error) {
	newStringList, err := uOfD.CreateReplicateAsRefinementFromURI(CrlStringListURI, hl, newURI...)
	if err != nil {
		return nil, errors.Wrap(err, "StringLists.go NewStringList failed")
	}
	return newStringList, nil
}

// AddStringListMemberAfter adds a string to the list after the priorMemberLiteral and returns the newMemberLiteral.
// If the priorMemberLiteral is nil, the string is added to the beginning of the list. An error is returned if the
// supplied list is not a list, the newMember is the empty string, or the priorMemberLiteral is not a CrlStringListMemberLiteral in this list.
func AddStringListMemberAfter(list core.Concept, priorMemberLiteral core.Concept, newMember string, hl *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsStringList(list, hl) {
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
		if !priorMemberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl) {
			return nil, errors.New("In AddStringListMemberAfter, supplied priorMemberLiteral is not a CrlStringListMemberLiteral")
		}
		if priorMemberLiteral.GetOwningConcept(hl) != list {
			return nil, errors.New("In AddStringListMemberAfter, supplied priorMemberLiteral does not belong to this list")
		}
	}
	var newPostMemberLiteral core.Concept
	referenceToPostMemberLiteral, err := getReferenceToNextMemberLiteral(priorMemberLiteral, hl)
	if err != nil {
		return nil, err
	}
	referencedPostMemberLiteral := referenceToPostMemberLiteral.GetReferencedConcept(hl)
	if referencedPostMemberLiteral != nil {
		newPostMemberLiteral = referencedPostMemberLiteral
	}
	newMemberLiteral, _ := uOfD.CreateReplicateLiteralAsRefinementFromURI(CrlStringListMemberLiteralURI, hl)
	newMemberLiteral.SetOwningConcept(list, hl)
	newMemberLiteral.SetLiteralValue(newMember, hl)
	// Wire up prior references
	setNextMemberLiteral(priorMemberLiteral, newMemberLiteral, hl)

	setPriorMemberLiteral(newMemberLiteral, priorMemberLiteral, hl)
	// Wire up next references
	if newPostMemberLiteral == nil {
		referenceToLastMemberLiteral, err2 := getStringListReferenceToLastMemberLiteral(list, hl)
		if err2 != nil {
			return nil, errors.Wrap(err2, "AddStringListMemberAfter failed")
		}
		if referenceToLastMemberLiteral != nil {
			referenceToLastMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, hl)
		}
	} else {
		setPriorMemberLiteral(newPostMemberLiteral, newMemberLiteral, hl)
		setNextMemberLiteral(newMemberLiteral, newPostMemberLiteral, hl)
	}

	return newMemberLiteral, nil
}

// AddStringListMemberBefore adds a member to the list before the postMember.
// If the postMember is nil, the member is added at the end of the list.
func AddStringListMemberBefore(list core.Concept, postMemberLiteral core.Concept, newMember string, hl *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsStringList(list, hl) {
		return nil, errors.New("In AddStringListMemberBefore, Supplied Element is not a CRL StringList")
	}
	if newMember == "" {
		return nil, errors.New("In AddStringListMemberBefore, Supplied string is empty: empty strings are not allowed in CRL StringLists")
	}
	// Check to ensure that the postMemberLiteral is valid
	if postMemberLiteral == nil {
		return nil, errors.New("AddStringListMemberBefore called with nil postMemberLiteral")
	}
	if !postMemberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl) {
		return nil, errors.New("In AddStringListMemberBefore, Supplied postMemberLiteral is not a CrlStringListMemberLiteral")
	}
	if postMemberLiteral.GetOwningConcept(hl) != list {
		return nil, errors.New("In AddStringListMemberBefore, Supplied postMemberLiteral does not belong to this list")
	}
	var newPriorMemberLiteral core.Concept
	// If the postMemberLiteral exists, then its priorMemberLiteral should point to the newMemberLiteral
	referenceToPriorMemberLiteral, err := getReferenceToPriorMemberLiteral(postMemberLiteral, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AddStringListMemberBefore failed")
	}
	referencedPriorMemberLiteral := referenceToPriorMemberLiteral.GetReferencedConcept(hl)
	if referencedPriorMemberLiteral != nil {
		newPriorMemberLiteral = referencedPriorMemberLiteral
	}
	// Create the newMemberLiteral
	newMemberLiteral, _ := uOfD.CreateReplicateLiteralAsRefinementFromURI(CrlStringListMemberLiteralURI, hl)
	newMemberLiteral.SetOwningConcept(list, hl)
	newMemberLiteral.SetLiteralValue(newMember, hl)
	// Wire up post references - be careful if inserting at the end
	setPriorMemberLiteral(postMemberLiteral, newMemberLiteral, hl)
	setNextMemberLiteral(newMemberLiteral, postMemberLiteral, hl)
	// Wire up prior references
	if newPriorMemberLiteral == nil {
		// The new member is the only member of the list
		referenceToFirstMemberLiteral, _ := getStringListReferenceToFirstMemberLiteral(list, hl)
		if referenceToFirstMemberLiteral != nil {
			referenceToFirstMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, hl)
		}
	} else {
		setNextMemberLiteral(newPriorMemberLiteral, newMemberLiteral, hl)
		setPriorMemberLiteral(newMemberLiteral, newPriorMemberLiteral, hl)
	}
	return newMemberLiteral, nil
}

// AppendStringListMember adds a string to the end of the list
func AppendStringListMember(list core.Concept, value string, hl *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsStringList(list, hl) {
		return nil, errors.New("In AddStringListMemberBefore, Supplied Element is not a CRL StringList")
	}
	if value == "" {
		return nil, errors.New("In AddStringListMemberBefore, Supplied string is empty: empty strings are not allowed in CRL StringLists")
	}
	oldLastMemberLiteral, err := GetLastMemberLiteral(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AppendStringListMember failed")
	}
	// Create the newMemberLiteral
	newMemberLiteral, err2 := uOfD.CreateReplicateLiteralAsRefinementFromURI(CrlStringListMemberLiteralURI, hl)
	if err2 != nil {
		return nil, errors.Wrap(err2, "AppendStringListMember failed")
	}
	err = newMemberLiteral.SetOwningConcept(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AppendStringListMember failed")
	}
	err = newMemberLiteral.SetLiteralValue(value, hl)
	if err != nil {
		return nil, errors.Wrap(err, "AppendStringListMember failed")
	}
	// Wire up references - be careful if inserting at the end
	referenceToLastMemberLiteral, err3 := getStringListReferenceToLastMemberLiteral(list, hl)
	if err3 != nil {
		return nil, errors.Wrap(err2, "AppendStringListMember failed")
	}
	if referenceToLastMemberLiteral != nil {
		err = referenceToLastMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, hl)
		if err != nil {
			return nil, errors.Wrap(err, "AppendStringListMember failed")
		}
	}
	if oldLastMemberLiteral == nil {
		referenceToFirstMemberLiteral, err4 := getStringListReferenceToFirstMemberLiteral(list, hl)
		if err4 != nil {
			return nil, errors.Wrap(err2, "AppendStringListMember failed")
		}
		err = referenceToFirstMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, hl)
		if err != nil {
			return nil, errors.Wrap(err, "AppendStringListMember failed")
		}
	} else {
		err = setNextMemberLiteral(oldLastMemberLiteral, newMemberLiteral, hl)
		if err != nil {
			return nil, errors.Wrap(err, "AppendStringListMember failed")
		}
		err = setPriorMemberLiteral(newMemberLiteral, oldLastMemberLiteral, hl)
		if err != nil {
			return nil, errors.Wrap(err, "AppendStringListMember failed")
		}
	}
	return newMemberLiteral, nil
}

// ClearStringList removes all members from the list
func ClearStringList(list core.Concept, hl *core.Transaction) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	it := list.GetOwnedConceptIDs(hl).Iterator()
	for id := range it.C {
		memberLiteral := uOfD.GetLiteral(id.(string))
		if memberLiteral != nil && memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl) {
			uOfD.DeleteElement(memberLiteral, hl)
		}
	}
}

// GetFirstMemberLiteral returns the reference to the first member of the list. It returns an error if the
// list is not a list. It returns nil if the list is empty
func GetFirstMemberLiteral(list core.Concept, hl *core.Transaction) (core.Concept, error) {
	refRef, err := getStringListReferenceToFirstMemberLiteral(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "GetFirstMemberLiteral failed")
	}
	if refRef == nil {
		return nil, errors.New("In StringList GetFirstMemberLiteral, No reference to first member literal found")
	}
	firstMemberLiteral := refRef.GetReferencedConcept(hl)
	if firstMemberLiteral == nil {
		return nil, nil
	}
	return firstMemberLiteral, nil
}

// GetFirstLiteralForString returns the first Literal whose value is the given string. It returns an error if the list is not a list.
// It returns nil if the string it is not found in the list.
func GetFirstLiteralForString(list core.Concept, value string, hl *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	it := list.GetOwnedConceptIDs(hl).Iterator()
	for id := range it.C {
		memberLiteral := uOfD.GetLiteral(id.(string))
		if memberLiteral != nil &&
			memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl) &&
			memberLiteral.GetLiteralValue(hl) == value {
			it.Stop()
			return memberLiteral, nil
		}
	}
	return nil, nil
}

// GetLastMemberLiteral returns the reference to the last member of the list. It returns an error if list is not a list.
// It returns nil if the list is empty
func GetLastMemberLiteral(list core.Concept, hl *core.Transaction) (core.Concept, error) {
	refRef, err := getStringListReferenceToLastMemberLiteral(list, hl)
	if err != nil {
		return nil, err
	}
	if refRef == nil {
		return nil, errors.New("No reference to last member literal found")
	}
	lastMemberLiteral := refRef.GetReferencedConcept(hl)
	if lastMemberLiteral == nil {
		return nil, nil
	}
	return lastMemberLiteral, nil
}

// getStringListReferenceToFirstMemberLiteral returns the reference to the first member literal. It returns an error if list is not a StringList
func getStringListReferenceToFirstMemberLiteral(list core.Concept, hl *core.Transaction) (core.Concept, error) {
	if !IsStringList(list, hl) {
		return nil, errors.New("Argument is not a CrlDataStructures.StringList")
	}
	refToLiteral := list.GetFirstOwnedReferenceRefinedFromURI(CrlStringListReferenceToFirstMemberLiteralURI, hl)
	if refToLiteral == nil {
		return nil, errors.New("In getStringListReferenceToFirstMemberLiteral, the reference was not found")
	}
	return refToLiteral, nil
}

// getStringListReferenceToLastMemberLiteral returns the reference to the last member literal. It returns an error if list is not a StringList
func getStringListReferenceToLastMemberLiteral(list core.Concept, hl *core.Transaction) (core.Concept, error) {
	if !IsStringList(list, hl) {
		return nil, errors.New("Argument is not a CrlDataStructures.StringList")
	}
	return list.GetFirstOwnedReferenceRefinedFromURI(CrlStringListReferenceToLastMemberLiteralURI, hl), nil
}

// GetNextMemberLiteral returns the successor member literal in the list
func GetNextMemberLiteral(memberLiteral core.Concept, hl *core.Transaction) (core.Concept, error) {
	if !IsStringListMemberLiteral(memberLiteral, hl) {
		return nil, errors.New("Supplied memberLiteral is not a refinement of CrlStringListMemberLiteral")
	}
	referenceToNextMemberLiteral, err := getReferenceToNextMemberLiteral(memberLiteral, hl)
	if err != nil {
		return nil, err
	}
	nextMemberLiteral := referenceToNextMemberLiteral.GetReferencedConcept(hl)
	if nextMemberLiteral == nil {
		return nil, nil
	}
	return nextMemberLiteral, nil
}

// GetPriorMemberLiteral returns the predecessor member literal in the list
func GetPriorMemberLiteral(memberLiteral core.Concept, hl *core.Transaction) (core.Concept, error) {
	if !IsStringListMemberLiteral(memberLiteral, hl) {
		return nil, errors.New("Supplied memberLiteral is not a refinement of CrlStringListMemberLiteral")
	}
	referenceToPriorMemberLiteral, err := getReferenceToPriorMemberLiteral(memberLiteral, hl)
	if err != nil {
		return nil, err
	}
	priorMemberLiteral := referenceToPriorMemberLiteral.GetReferencedConcept(hl)
	if priorMemberLiteral == nil {
		return nil, nil
	}
	return priorMemberLiteral.(core.Concept), nil
}

// getReferenceToNextMemberLiteral returns the reference to the next member of the list.
// It returns nil if the reference is the last member of the list
func getReferenceToNextMemberLiteral(memberLiteral core.Concept, hl *core.Transaction) (core.Concept, error) {
	if memberLiteral == nil {
		return nil, errors.New("GetNextMemberLiteral called with nil memberLiteral")
	}
	if !IsStringListMemberLiteral(memberLiteral, hl) {
		return nil, errors.New("Supplied memberLiteral is not a refinement of CrlStringListMemberLiteral")
	}
	nextMemberLiteral := memberLiteral.GetFirstOwnedReferenceRefinedFromURI(CrlStringListReferenceToNextMemberLiteralURI, hl)
	if nextMemberLiteral == nil {
		return nil, errors.New("In GetNextMemberLiteral, memberLiteral does not ave a NextMemberLiteralReference")
	}
	return nextMemberLiteral, nil
}

// getReferenceToPriorMemberLiteral returns the reference to the previous member of the list. It returns an error if the memberLiteral
// is either nil or is not a refinement of CrlStringListMemberLiteral
// It returns nil if the reference is the first member of the list
func getReferenceToPriorMemberLiteral(memberLiteral core.Concept, hl *core.Transaction) (core.Concept, error) {
	if memberLiteral == nil {
		return nil, errors.New("getReferenceToPriorMemberLiteral called with nil memberLiteral")
	}
	if !IsStringListMemberLiteral(memberLiteral, hl) {
		return nil, errors.New("In getReferenceToPriorMemberLiteral, supplied memberLiteral is not a refinement of CrlStringListMemberLiteral")
	}
	priorMemberLiteralReference := memberLiteral.GetFirstOwnedReferenceRefinedFromURI(CrlStringListReferenceToPriorMemberLiteralURI, hl)
	if priorMemberLiteralReference == nil {
		return nil, errors.New("In getReferenceToPriorMemberLiteral, memberLiteral does not have a PriorMemberLiteralReference")
	}
	return priorMemberLiteralReference, nil
}

// IsStringList returns true if the supplied Element is a refinement of StringList
func IsStringList(list core.Concept, hl *core.Transaction) bool {
	return list.IsRefinementOfURI(CrlStringListURI, hl)
}

// IsStringListMember returns true if the string is a memeber of the given list
func IsStringListMember(list core.Concept, value string, hl *core.Transaction) bool {
	uOfD := list.GetUniverseOfDiscourse(hl)
	it := list.GetOwnedConceptIDs(hl).Iterator()
	for id := range it.C {
		memberLiteral := uOfD.GetLiteral(id.(string))
		if memberLiteral != nil && memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl) && memberLiteral.GetLiteralValue(hl) == value {
			it.Stop()
			return true
		}
	}
	return false
}

// IsStringListMemberLiteral returns true if the supplied Reference is a refinement of StringListMemberLiteral
func IsStringListMemberLiteral(memberLiteral core.Concept, hl *core.Transaction) bool {
	return memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl)
}

// PrependStringListMember adds a string to the beginning of the list
func PrependStringListMember(list core.Concept, value string, hl *core.Transaction) (core.Concept, error) {
	uOfD := list.GetUniverseOfDiscourse(hl)
	if !IsStringList(list, hl) {
		return nil, errors.New("In PrependStringListMember, Supplied Element is not a CRL StringList")
	}
	if value == "" {
		return nil, errors.New("In PrependStringListMember, Supplied string is empty: empty strings are not allowed in CRL StringLists")
	}
	oldFirstMemberLiteral, err := GetFirstMemberLiteral(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "PrependStringListMember failed")
	}
	// Create the newMemberLiteral
	newMemberLiteral, err2 := uOfD.CreateReplicateLiteralAsRefinementFromURI(CrlStringListMemberLiteralURI, hl)
	if err2 != nil {
		return nil, errors.Wrap(err2, "PrependStringListMember failed")
	}
	err = newMemberLiteral.SetOwningConcept(list, hl)
	if err != nil {
		return nil, errors.Wrap(err, "PrependStringListMember failed")
	}
	err = newMemberLiteral.SetLiteralValue(value, hl)
	if err != nil {
		return nil, errors.Wrap(err, "PrependStringListMember failed")
	}
	// Wire up references - be careful if inserting at the end
	referenceToFirstMemberLiteral, err3 := getStringListReferenceToFirstMemberLiteral(list, hl)
	if err3 != nil {
		return nil, errors.Wrap(err2, "PrependStringListMember failed")
	}
	if referenceToFirstMemberLiteral != nil {
		err = referenceToFirstMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, hl)
		if err != nil {
			return nil, errors.Wrap(err, "PrependStringListMember failed")
		}
	}
	if oldFirstMemberLiteral == nil {
		referenceToLastMemberLiteral, err4 := getStringListReferenceToLastMemberLiteral(list, hl)
		if err4 != nil {
			return nil, errors.Wrap(err2, "PrependStringListMember failed")
		}
		err = referenceToLastMemberLiteral.SetReferencedConcept(newMemberLiteral, core.NoAttribute, hl)
		if err != nil {
			return nil, errors.Wrap(err, "PrependStringListMember failed")
		}
	} else {
		err = setPriorMemberLiteral(oldFirstMemberLiteral, newMemberLiteral, hl)
		if err != nil {
			return nil, errors.Wrap(err, "PrependStringListMember failed")
		}
	}
	err = setNextMemberLiteral(newMemberLiteral, oldFirstMemberLiteral, hl)
	if err != nil {
		return nil, errors.Wrap(err, "PrependStringListMember failed")
	}
	return newMemberLiteral, nil
}

// RemoveStringListMember removes the first occurrance of an element from the given list
func RemoveStringListMember(list core.Concept, value string, hl *core.Transaction) error {
	uOfD := list.GetUniverseOfDiscourse(hl)
	it := list.GetOwnedConceptIDs(hl).Iterator()
	for id := range it.C {
		memberLiteral := uOfD.GetLiteral(id.(string))
		if memberLiteral != nil && memberLiteral.IsRefinementOfURI(CrlStringListMemberLiteralURI, hl) && memberLiteral.GetLiteralValue(hl) == value {
			// Modify previous and next pointers
			priorMemberLiteral, _ := GetPriorMemberLiteral(memberLiteral, hl)
			nextMemberLiteral, _ := GetNextMemberLiteral(memberLiteral, hl)
			if priorMemberLiteral != nil {
				setNextMemberLiteral(priorMemberLiteral, nextMemberLiteral, hl)
			} else {
				referenceToFirstMemberLiteral, _ := getStringListReferenceToFirstMemberLiteral(list, hl)
				referenceToFirstMemberLiteral.SetReferencedConcept(nextMemberLiteral, core.NoAttribute, hl)
			}
			if nextMemberLiteral != nil {
				setPriorMemberLiteral(nextMemberLiteral, priorMemberLiteral, hl)
			} else {
				referenceToLastMemberLiteral, _ := getStringListReferenceToLastMemberLiteral(list, hl)
				referenceToLastMemberLiteral.SetReferencedConcept(priorMemberLiteral, core.NoAttribute, hl)
			}
			// Now delete the member literal
			uOfD.DeleteElement(memberLiteral, hl)
			it.Stop()
			return nil
		}
	}
	return errors.New("element not member of list")
}

// setNextMemberLiteral takes a memberLiteral and sets its next reference
func setNextMemberLiteral(memberLiteral core.Concept, nextLiteral core.Concept, hl *core.Transaction) error {
	// since this is an internal function we assume that the references are refinements of CrlStringListMemberLiteral
	nextLiteralReference, err := getReferenceToNextMemberLiteral(memberLiteral, hl)
	if err != nil {
		return errors.Wrap(err, "setNextMemberLiteral failed")
	}
	err = nextLiteralReference.SetReferencedConcept(nextLiteral, core.NoAttribute, hl)
	if err != nil {
		return errors.Wrap(err, "setNextMemberLiteral failed")
	}
	return nil
}

// setPriorMemberLiteral takes a memberLiteral and sets its prior reference
func setPriorMemberLiteral(memberLiteral core.Concept, priorLiteral core.Concept, hl *core.Transaction) error {
	// since this is an internal function we assume that the references are refinements of CrlStringListMemberLiteral
	priorLiteralReference, err := getReferenceToPriorMemberLiteral(memberLiteral, hl)
	if err != nil {
		return errors.Wrap(err, "setNextMemberLiteral failed")
	}
	err = priorLiteralReference.SetReferencedConcept(priorLiteral, core.NoAttribute, hl)
	if err != nil {
		return errors.Wrap(err, "setNextMemberLiteral failed")
	}
	return nil
}

// BuildCrlStringListsConcepts builds the CrlStringList concept and adds it as a child of the provided parent concept space
func BuildCrlStringListsConcepts(uOfD *core.UniverseOfDiscourse, parentSpace core.Concept, hl *core.Transaction) {
	crlStringList, _ := uOfD.NewElement(hl, CrlStringListURI)
	crlStringList.SetLabel("CrlStringList", hl)
	crlStringList.SetOwningConcept(parentSpace, hl)

	crlFirstMemberLiteral, _ := uOfD.NewReference(hl, CrlStringListReferenceToFirstMemberLiteralURI)
	crlFirstMemberLiteral.SetLabel("StringListFirstMemberLiteral", hl)
	crlFirstMemberLiteral.SetOwningConcept(crlStringList, hl)

	crlLastMemberLiteral, _ := uOfD.NewReference(hl, CrlStringListReferenceToLastMemberLiteralURI)
	crlLastMemberLiteral.SetLabel("StringListLastMemberLiteral", hl)
	crlLastMemberLiteral.SetOwningConcept(crlStringList, hl)

	crlStringListMemberLiteral, _ := uOfD.NewLiteral(hl, CrlStringListMemberLiteralURI)
	crlStringListMemberLiteral.SetLabel("StringListMemberLiteral", hl)
	crlStringListMemberLiteral.SetOwningConcept(parentSpace, hl)

	crlNextMemberLiteral, _ := uOfD.NewReference(hl, CrlStringListReferenceToNextMemberLiteralURI)
	crlNextMemberLiteral.SetLabel("StringListNextMemberLiteral", hl)
	crlNextMemberLiteral.SetOwningConcept(crlStringListMemberLiteral, hl)

	crlPriorMemberLiteral, _ := uOfD.NewReference(hl, CrlStringListReferenceToPriorMemberLiteralURI)
	crlPriorMemberLiteral.SetLabel("StringListPriorMemberLiteral", hl)
	crlPriorMemberLiteral.SetOwningConcept(crlStringListMemberLiteral, hl)
}
