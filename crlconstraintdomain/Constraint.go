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

// CrlConstraintComplianceURI is the URI for a constraint status
var CrlConstraintComplianceURI = CrlConstraintDomainURI + "/ConstraintCompliance"

// CrlConstraintSatisfiedURI is the URI for the boolean indicating whether the constraint is satisfied
var CrlConstraintSatisfiedURI = CrlConstraintComplianceURI + "/Satisfied"

// CrlConstraintSpecificationReferenceURI is the URI for reference to the constraint specification whose compliance is reported
var CrlConstraintSpecificationReferenceURI = CrlConstraintComplianceURI + "/SpecificationReference"

// Multiplicity Constraint

// CrlMultiplicityConstrainedURI is the URI for the concept of having multiplicity constraints
var CrlMultiplicityConstrainedURI = CrlConstraintDomainURI + "/MultiplicityConstrained"

// CrlMultiplicityConstraintSpecificationURI is the URI for a multiplicity constraint
var CrlMultiplicityConstraintSpecificationURI = CrlConstraintDomainURI + "/MultiplicityConstraintSpecification"

// CrlMultiplicityConstraintMultiplicityURI is the URI for the multiplicity specification
var CrlMultiplicityConstraintMultiplicityURI = CrlMultiplicityConstraintSpecificationURI + "/Multiplicity"

// CrlMultiplicityConstraintConstrainedConceptURI is the URI for the concecept whose multiplicity is being constrained
var CrlMultiplicityConstraintConstrainedConceptURI = CrlMultiplicityConstraintSpecificationURI + "/ConstrainedConcept"

// NewMultiplicityConstraintSpecification creates and initializes a multiplicity constraint specification
func NewMultiplicityConstraintSpecification(owner core.Concept, constrainedConcept core.Concept, label string, multiplicity string, trans *core.Transaction, newURI ...string) (core.Concept, error) {
	if owner == nil || constrainedConcept == nil {
		return nil, errors.New("NewMultiplicityConstraintSpecification called with nil owner or constrained concept")
	}
	uOfD := trans.GetUniverseOfDiscourse()
	newConstraint, err := uOfD.CreateOwnedRefinementOfConceptURI(CrlMultiplicityConstraintSpecificationURI, owner, label, trans, newURI...)
	if err != nil {
		return nil, errors.Wrap(err, "NewMultiplicityConstraintSpecification failed")
	}
	multiplicitySpecification, err3 := uOfD.CreateOwnedRefinementOfConceptURI(CrlMultiplicityConstraintMultiplicityURI, newConstraint, "IsSatisfied", trans)
	if err3 != nil {
		return nil, errors.Wrap(err, "NewMultiplicityConstraintSpecification failed")
	}
	if !IsValidMultiplicity(multiplicity) {
		return nil, errors.New("NewMultiplicityConstraintSpecification called with invalid multiplicity: " + multiplicity)
	}
	multiplicitySpecification.SetLiteralValue(multiplicity, trans)

	constrainedConceptReference, err2 := uOfD.CreateOwnedRefinementOfConceptURI(CrlMultiplicityConstraintConstrainedConceptURI, newConstraint, "Constrained Concept", trans)
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
	return newConstraint, nil
}

// NewConstraintCompliance creates and initializes a refinement of a ConstraintCompliance
func NewConstraintCompliance(owner core.Concept, constraintSpecification core.Concept, trans *core.Transaction) core.Concept {
	uOfD := trans.GetUniverseOfDiscourse()
	constraintCompliance, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlConstraintComplianceURI, owner, "ConstraintCompliance", trans)
	uOfD.CreateOwnedRefinementOfConceptURI(CrlConstraintSatisfiedURI, constraintCompliance, "Satisfied", trans)
	constraintSpecificationReference, _ := uOfD.CreateOwnedRefinementOfConceptURI(CrlConstraintSpecificationReferenceURI, constraintCompliance, "ConstraintSpecificationReference", trans)
	constraintSpecificationReference.SetReferencedConcept(constraintSpecification, core.NoAttribute, trans)
	return constraintCompliance
}

// GetConstrainedConceptType returns the concept whose multiplicity is being constrained
func GetConstrainedConceptType(target core.Concept, trans *core.Transaction) (core.Concept, error) {
	if !target.IsRefinementOfURI(CrlMultiplicityConstraintSpecificationURI, trans) {
		return nil, errors.New("GetMultiplicity called with invalid target")
	}
	conceptReference := target.GetFirstOwnedConceptRefinedFromURI(CrlMultiplicityConstraintConstrainedConceptURI, trans)
	if conceptReference == nil {
		return nil, errors.New("GetConstrainedConceptType failed: conceptReference not found")
	}
	return conceptReference.GetReferencedConcept(trans), nil
}

// GetConstraintSpecification returns the constraint specification for this compliance instance
func GetConstraintSpecification(constraintComplianceInstance core.Concept, trans *core.Transaction) core.Concept {
	if !constraintComplianceInstance.IsRefinementOfURI(CrlConstraintComplianceURI, trans) {
		return nil
	}
	constraintSpecificationReference := constraintComplianceInstance.GetFirstOwnedReferenceRefinedFromURI(CrlConstraintSpecificationReferenceURI, trans)
	return constraintSpecificationReference.GetReferencedConcept(trans)
}

// GetMultiplicity returns the literal value after checking that the target is valid
func GetMultiplicity(target core.Concept, trans *core.Transaction) (string, error) {
	if !target.IsRefinementOfURI(CrlMultiplicityConstraintSpecificationURI, trans) {
		return "", errors.New("GetMultiplicity called with invalid target")
	}
	spec := getMultiplicitySpecification(target, trans)
	return spec.GetLiteralValue(trans), nil
}

func getMultiplicitySpecification(target core.Concept, trans *core.Transaction) core.Concept {
	return target.GetFirstOwnedConceptRefinedFromURI(CrlMultiplicityConstraintMultiplicityURI, trans)
}

// IsSatisfied returns true if the ConstraintCompliance.ConstraintSatisfied is true
func IsSatisfied(constraintCompliance core.Concept, trans *core.Transaction) bool {
	if !constraintCompliance.IsRefinementOfURI(CrlConstraintComplianceURI, trans) {
		return false
	}
	constraintSatisfied := constraintCompliance.GetFirstOwnedConceptRefinedFromURI(CrlConstraintSatisfiedURI, trans)
	if constraintSatisfied == nil {
		return false
	}
	value, err := crldatatypesdomain.GetBooleanValue(constraintSatisfied, trans)
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
func SetMultiplicity(target core.Concept, multiplicity string, trans *core.Transaction) error {
	if !target.IsRefinementOfURI(CrlMultiplicityConstraintSpecificationURI, trans) {
		return errors.New("SetMultiplicity called with invalid target")
	}
	if !IsValidMultiplicity(multiplicity) {
		return errors.New("SetMultiplicity called with invalid multiplicity: " + multiplicity)
	}
	spec := getMultiplicitySpecification(target, trans)
	err := spec.SetLiteralValue(multiplicity, trans)
	if err != nil {
		return errors.Wrap(err, "SetMultiplicity failed")
	}
	return nil
}

func evaluateMultiplicityConstraints(constrainedConcept core.Concept, notification *core.ChangeNotification, trans *core.Transaction) error {
	// Determine which immediate ancestor owns the constraint specifications
	uOfD := trans.GetUniverseOfDiscourse()
	multiplicityConstrainedConcept := uOfD.GetElementWithURI(CrlMultiplicityConstrainedURI)
	var definingAbstraction core.Concept
	immediateAbstractions := make(map[string]core.Concept)
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
		constraintComplianceInstances := constrainedConcept.GetOwnedConceptsRefinedFromURI(CrlConstraintComplianceURI, trans)
		var foundComplianceInstance core.Concept
		for _, constraintComplianceInstance := range constraintComplianceInstances {
			if GetConstraintSpecification(constraintComplianceInstance, trans) == constraintSpecification {
				foundComplianceInstance = constraintComplianceInstance
				break
			}
		}
		if foundComplianceInstance == nil {
			foundComplianceInstance = NewConstraintCompliance(constrainedConcept, constraintSpecification, trans)
		}
		constrainedConceptType, err := GetConstrainedConceptType(constraintSpecification, trans)
		if err != nil {
			return errors.Wrap(err, "evaluateMultiplicityConstraints failed")
		}
		typedChildren := constrainedConcept.GetOwnedConceptsRefinedFrom(constrainedConceptType, trans)
		multiplicity, err2 := GetMultiplicity(constraintSpecification, trans)
		if err2 != nil {
			return errors.Wrap(err, "evaluateMultiplicityConstraints failed")
		}
		sat := SatisfiesMultiplicity(multiplicity, len(typedChildren))
		satisfied := foundComplianceInstance.GetFirstOwnedConceptRefinedFromURI(CrlConstraintSatisfiedURI, trans)
		crldatatypesdomain.SetBooleanValue(satisfied, sat, trans)
	}
	return nil
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
