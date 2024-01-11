package crlconstraintdomain

import (
	"strconv"
	"strings"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatatypesdomain"
	"github.com/pkg/errors"
)

// CrlConstraintDomainURI is the uri for the domain that defines the Crl Constraints
var CrlConstraintDomainURI = "http://activeCRL.com/crldatastructuresdomain/CrlConstraintDomain"

// CrlConstraintSpecificationURI is the URI for a constraint specification
var CrlConstraintSpecificationURI = CrlConstraintDomainURI + "/ConstraintSpecification"

// CrlConstraintSpecification is the CRL representation of a constraint specification
type CrlConstraintSpecification core.Concept

// CrlConstraintComplianceURI is the URI for a constraint status
var CrlConstraintComplianceURI = CrlConstraintDomainURI + "/ConstraintCompliance"

// CrlConstraintCompliance is the CRL representation of a constraint compliance

// CrlConstraintSatisfiedURI is the URI for the boolean indicating whether the constraint is satisfied
var CrlConstraintSatisfiedURI = CrlConstraintComplianceURI + "/Satisfied"

// CrlConstraintSpecificationReferenceURI is the URI for reference to the constraint specification whose compliance is reported
var CrlConstraintSpecificationReferenceURI = CrlConstraintComplianceURI + "/SpecificationReference"

// Multiplicity Constraint

// CrlMultiplicityConstrainedURI is the URI for the concept of having multiplicity constraints
var CrlMultiplicityConstrainedURI = CrlConstraintDomainURI + "/MultiplicityConstrained"

// CrlMultiplicityConstraintSpecificationURI is the URI for a multiplicity constraint
var CrlMultiplicityConstraintSpecificationURI = CrlConstraintDomainURI + "/MultiplicityConstraintSpecification"

// CrlMultiplicityConstraintSpecification is the CRL representation of a multiplicity constraint
type CrlMultiplicityConstraintSpecification CrlConstraintSpecification

// CrlMultiplicityConstraintMultiplicityURI is the URI for the multiplicity specification
var CrlMultiplicityConstraintMultiplicityURI = CrlMultiplicityConstraintSpecificationURI + "/Multiplicity"

// CrlMultiplicityConstraintConstrainedConceptURI is the URI for the concecept whose multiplicity is being constrained
var CrlMultiplicityConstraintConstrainedConceptURI = CrlMultiplicityConstraintSpecificationURI + "/ConstrainedConcept"

// NewMultiplicityConstraintSpecification creates and initializes a multiplicity constraint specification
func NewMultiplicityConstraintSpecification(owner *core.Concept, constrainedConcept *core.Concept, label string, multiplicity string, trans *core.Transaction, newURI ...string) (*CrlMultiplicityConstraintSpecification, error) {
	if owner == nil || constrainedConcept == nil {
		return nil, errors.New("NewMultiplicityConstraintSpecification called with nil owner or constrained concept")
	}
	uOfD := trans.GetUniverseOfDiscourse()
	newConcept, err := uOfD.CreateRefinementOfConceptURI(CrlMultiplicityConstraintSpecificationURI, label, trans, newURI...)
	if err != nil {
		return nil, errors.Wrap(err, "NewMultiplicityConstraintSpecification failed")
	}
	newMcs := (*CrlMultiplicityConstraintSpecification)(newConcept)
	multiplicitySpecification, err3 := uOfD.CreateOwnedRefinementOfConceptURI(CrlMultiplicityConstraintMultiplicityURI, newMcs.AsCore(), "Multiplicity", trans)
	if err3 != nil {
		return nil, errors.Wrap(err, "NewMultiplicityConstraintSpecification failed")
	}
	if !IsValidMultiplicity(multiplicity) {
		return nil, errors.New("NewMultiplicityConstraintSpecification called with invalid multiplicity: " + multiplicity)
	}
	multiplicitySpecification.SetLiteralValue(multiplicity, trans)

	constrainedConceptReference, err2 := uOfD.CreateOwnedRefinementOfConceptURI(CrlMultiplicityConstraintConstrainedConceptURI, newMcs.AsCore(), "Constrained Concept", trans)
	if err2 != nil {
		return nil, errors.Wrap(err, "NewMultiplicityConstraintSpecification failed")
	}
	constrainedConceptReference.SetReferencedConcept(constrainedConcept, core.NoAttribute, trans)
	if !owner.IsRefinementOfURI(CrlMultiplicityConstrainedURI, trans) {
		multiplicityConstrained := uOfD.GetElementWithURI(CrlMultiplicityConstrainedURI)
		_, err4 := uOfD.NewOwnedRefinement(owner, "", multiplicityConstrained, owner, trans)
		if err4 != nil {
			return nil, errors.Wrap(err, "NewMultiplicityConstraintSpecification failed")
		}
	}
	newConcept.SetOwningConcept(owner, trans)
	return newMcs, nil
}

// AsCore casts the CrlMultiplicityConstraintSpecification pointer to *core.Concept
func (mcs *CrlMultiplicityConstraintSpecification) AsCore() *core.Concept {
	return (*core.Concept)(mcs)
}

// NewConstraintCompliance creates and initializes a refinement of a ConstraintCompliance
func NewConstraintCompliance(owner *core.Concept, constraintSpecification *core.Concept, trans *core.Transaction) *core.Concept {
	uOfD := trans.GetUniverseOfDiscourse()
	constraintCompliance, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlConstraintComplianceURI, owner, "ConstraintCompliance", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlConstraintSatisfiedURI, constraintCompliance, "Satisfied", trans)
	constraintSpecificationReference, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlConstraintSpecificationReferenceURI, constraintCompliance, "ConstraintSpecificationReference", trans)
	constraintSpecificationReference.SetReferencedConcept(constraintSpecification, core.NoAttribute, trans)
	return constraintCompliance
}

// GetConstrainedConceptType returns the concept whose multiplicity is being constrained
func (mcs *CrlMultiplicityConstraintSpecification) GetConstrainedConceptType(trans *core.Transaction) (*core.Concept, error) {
	if !mcs.AsCore().IsRefinementOfURI(CrlMultiplicityConstraintSpecificationURI, trans) {
		return nil, errors.New("GetMultiplicity called with invalid target")
	}
	conceptReference := mcs.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlMultiplicityConstraintConstrainedConceptURI, trans)
	if conceptReference == nil {
		return nil, errors.New("GetConstrainedConceptType failed: conceptReference not found")
	}
	return conceptReference.GetReferencedConcept(trans), nil
}

// GetConstraintSpecification returns the constraint specification for this compliance instance
func GetConstraintSpecification(constraintComplianceInstance *core.Concept, trans *core.Transaction) *core.Concept {
	if !constraintComplianceInstance.IsRefinementOfURI(CrlConstraintComplianceURI, trans) {
		return nil
	}
	constraintSpecificationReference := constraintComplianceInstance.GetFirstOwnedReferenceRefinedFromURI(CrlConstraintSpecificationReferenceURI, trans)
	return constraintSpecificationReference.GetReferencedConcept(trans)
}

// GetMultiplicity returns the literal value of the multiplicity specification
func (mcs *CrlMultiplicityConstraintSpecification) GetMultiplicity(trans *core.Transaction) (string, error) {
	if !mcs.AsCore().IsRefinementOfURI(CrlMultiplicityConstraintSpecificationURI, trans) {
		return "", errors.New("GetMultiplicity called with invalid target")
	}
	spec := mcs.GetMultiplicityLiteral(trans)
	return spec.GetLiteralValue(trans), nil
}

func (mcs *CrlMultiplicityConstraintSpecification) GetMultiplicityLiteral(trans *core.Transaction) *core.Concept {
	return mcs.AsCore().GetFirstOwnedConceptRefinedFromURI(CrlMultiplicityConstraintMultiplicityURI, trans)
}

// IsMultiplicityConstraint returns true if the supplied concept is a multiplicity constraint
func IsMultiplicityConstraint(concept *core.Concept, trans *core.Transaction) bool {
	return concept.IsRefinementOfURI(CrlMultiplicityConstraintSpecificationURI, trans)
}

// IsSatisfied returns true if the ConstraintCompliance.ConstraintSatisfied is true
func IsSatisfied(constraintCompliance *core.Concept, trans *core.Transaction) bool {
	if !constraintCompliance.IsRefinementOfURI(CrlConstraintComplianceURI, trans) {
		return false
	}
	constraintSatisfied := constraintCompliance.GetFirstOwnedConceptRefinedFromURI(CrlConstraintSatisfiedURI, trans)
	if constraintSatisfied == nil {
		return false
	}
	value, err := (*crldatatypesdomain.CrlBoolean)(constraintSatisfied).GetBooleanValue(trans)
	if err != nil {
		return false
	}
	return value
}

// IsValidMultiplicity returns true if the supplied string is a valid multiplicity setting
func IsValidMultiplicity(multiplicity string) bool {
	if multiplicity == "" {
		return true
	}
	substrings := strings.Split(multiplicity, "..")
	switch len(substrings) {
	case 1:
		if substrings[0] == "*" {
			return true
		}
		_, err := strconv.Atoi(substrings[0])
		if err != nil {
			return false
		}
		return true
	case 2:
		_, err := strconv.Atoi(substrings[0])
		if err != nil {
			return false
		}
		if substrings[1] == "*" {
			return true
		}
		_, err = strconv.Atoi(substrings[1])
		if err == nil {
			return true
		}
	}
	return false
}

// SatisfiesMultiplicity returns true if the multiplicity is valid and the supplied value satisfies that multiplicity
func SatisfiesMultiplicity(multiplicity string, candidate int) bool {
	if !IsValidMultiplicity(multiplicity) {
		return false
	}
	if multiplicity == "" {
		return true
	}
	substrings := strings.Split(multiplicity, "..")
	switch len(substrings) {
	case 1:
		if substrings[0] == "*" {
			return true
		}
		i, _ := strconv.Atoi(substrings[0])
		return candidate == i
	case 2:
		i, _ := strconv.Atoi(substrings[0])
		if candidate < i {
			return false
		}
		if substrings[1] == "*" {
			return true
		}
		i, _ = strconv.Atoi(substrings[1])
		if candidate <= i {
			return true
		}
	}
	return false
}

// SetMultiplicity sets the multiplicity specification after checking that the target and the multiplicity are both valid
func (mcs *CrlMultiplicityConstraintSpecification) SetMultiplicity(multiplicity string, trans *core.Transaction) error {
	if !mcs.AsCore().IsRefinementOfURI(CrlMultiplicityConstraintSpecificationURI, trans) {
		return errors.New("SetMultiplicity called with invalid target")
	}
	if !IsValidMultiplicity(multiplicity) {
		return errors.New("SetMultiplicity called with invalid multiplicity: " + multiplicity)
	}
	spec := mcs.GetMultiplicityLiteral(trans)
	err := spec.SetLiteralValue(multiplicity, trans)
	if err != nil {
		return errors.Wrap(err, "SetMultiplicity failed")
	}
	return nil
}

func evaluateMultiplicityConstraints(constrainedConcept *core.Concept, notification *core.ChangeNotification, trans *core.Transaction) error {
	// Determine which immediate ancestor owns the constraint specifications
	uOfD := trans.GetUniverseOfDiscourse()
	multiplicityConstrainedConcept := uOfD.GetElementWithURI(CrlMultiplicityConstrainedURI)
	var definingAbstraction *core.Concept
	immediateAbstractions := make(map[string]*core.Concept)
	constrainedConcept.FindImmediateAbstractions(immediateAbstractions, trans)
	for _, abstraction := range immediateAbstractions {
		if abstraction != multiplicityConstrainedConcept && abstraction.IsRefinementOf(multiplicityConstrainedConcept, trans) {
			definingAbstraction = abstraction
			break
		}
	}
	if definingAbstraction == nil {
		return nil
	}
	// Now for each defined multiplicity constraint, ensure that the constrainedConcept has a corresponding compliance instance,
	// then evaluate the constraint and set the compliance value
	multiplicityConstraintSpecifications := definingAbstraction.GetOwnedConceptsRefinedFromURI(CrlMultiplicityConstraintSpecificationURI, trans)
	for _, constraintSpecification := range multiplicityConstraintSpecifications {
		mcs := (*CrlMultiplicityConstraintSpecification)(constraintSpecification)
		constraintComplianceInstances := constrainedConcept.GetOwnedConceptsRefinedFromURI(CrlConstraintComplianceURI, trans)
		var foundComplianceInstance *core.Concept
		for _, constraintComplianceInstance := range constraintComplianceInstances {
			if GetConstraintSpecification(constraintComplianceInstance, trans) == constraintSpecification {
				foundComplianceInstance = constraintComplianceInstance
				break
			}
		}
		if foundComplianceInstance == nil {
			foundComplianceInstance = NewConstraintCompliance(constrainedConcept, constraintSpecification, trans)
		}
		constrainedConceptType, err := mcs.GetConstrainedConceptType(trans)
		if err != nil {
			return errors.Wrap(err, "evaluateMultiplicityConstraints failed")
		}
		typedChildren := constrainedConcept.GetOwnedConceptsRefinedFrom(constrainedConceptType, trans)
		multiplicity, err2 := mcs.GetMultiplicity(trans)
		if err2 != nil {
			return errors.Wrap(err, "evaluateMultiplicityConstraints failed")
		}
		sat := SatisfiesMultiplicity(multiplicity, len(typedChildren))
		satisfied := foundComplianceInstance.GetFirstOwnedConceptRefinedFromURI(CrlConstraintSatisfiedURI, trans)
		(*crldatatypesdomain.CrlBoolean)(satisfied).SetBooleanValue(sat, trans)
	}
	return nil
}

// GetMultiplicityConstraint returns the multiplicity constraint if the indicated concept's owner has a multiplicity constraint on the concept
// Otherwise it returns nil
func GetMultiplicityConstraint(concept *core.Concept, trans *core.Transaction) *CrlMultiplicityConstraintSpecification {
	owner := concept.GetOwningConcept(trans)
	if owner == nil {
		return nil
	}
	multiplicityConstraints := owner.GetOwnedConceptsRefinedFromURI(CrlMultiplicityConstraintSpecificationURI, trans)
	for _, multiplicityConstraint := range multiplicityConstraints {
		typedConstraint := (*CrlMultiplicityConstraintSpecification)(multiplicityConstraint)
		constrainedConcept, err := typedConstraint.GetConstrainedConceptType(trans)
		if err != nil {
			return nil
		}
		if constrainedConcept.ConceptID == concept.ConceptID {
			return typedConstraint
		}

	}
	return nil
}

// HasMultiplicityConstraint returns true if the indicated concept's owner has a multiplicity constraint on the concept
func HasMultiplicityConstraint(concept *core.Concept, trans *core.Transaction) bool {
	owner := concept.GetOwningConcept(trans)
	if owner == nil {
		return false
	}
	multiplicityConstraints := owner.GetOwnedConceptsRefinedFromURI(CrlMultiplicityConstraintSpecificationURI, trans)
	for _, multiplicityConstraint := range multiplicityConstraints {
		typedConstraint := (*CrlMultiplicityConstraintSpecification)(multiplicityConstraint)
		constrainedConcept, err := typedConstraint.GetConstrainedConceptType(trans)
		if err != nil {
			return false
		}
		if constrainedConcept.ConceptID == concept.ConceptID {
			return true
		}

	}
	return false
}

// BuildCrlConstraintDomain constructs the concept space for CRL Constraints
func BuildCrlConstraintDomain(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) {
	if uOfD.GetElementWithURI(crldatatypesdomain.CrlDataTypesDomainURI) == nil {
		crldatatypesdomain.BuildCrlDataTypesDomain(uOfD, trans)
	}
	crlConstraintDomain, _ := uOfD.NewElement(trans, CrlConstraintDomainURI)
	crlConstraintDomain.SetLabel("CrlConstraintDomain", trans)

	crlConstraintCompliance, _ := uOfD.NewElement(trans, CrlConstraintComplianceURI)
	crlConstraintCompliance.SetLabel("ConstraintCompliance", trans)
	crldatatypesdomain.NewOwnedBoolean(crlConstraintCompliance, "Satisfied", trans, CrlConstraintSatisfiedURI)
	uOfD.NewOwnedReference(crlConstraintCompliance, "ConstraintSpecificationReference", trans, CrlConstraintSpecificationReferenceURI)

	crlConstraintSpecification, _ := uOfD.NewOwnedElement(crlConstraintDomain, "ConstraintSpecification", trans, CrlConstraintSpecificationURI)

	uOfD.NewOwnedElement(crlConstraintDomain, "Multiplicity Constrained", trans, CrlMultiplicityConstrainedURI)

	crlMultiplicityConstraintSpecification, _ := uOfD.CreateOwnedRefinementOfConcept(crlConstraintSpecification, crlConstraintDomain, "MultiplicityConstraintSpecification", trans, CrlMultiplicityConstraintSpecificationURI)
	uOfD.CreateOwnedRefinementOfConceptURI(core.LiteralURI, crlMultiplicityConstraintSpecification, "Multiplicity", trans, CrlMultiplicityConstraintMultiplicityURI)
	uOfD.CreateOwnedRefinementOfConceptURI(core.ReferenceURI, crlMultiplicityConstraintSpecification, "ConstrainedConceptReference", trans, CrlMultiplicityConstraintConstrainedConceptURI)
	uOfD.AddFunction(CrlMultiplicityConstrainedURI, evaluateMultiplicityConstraints)
}