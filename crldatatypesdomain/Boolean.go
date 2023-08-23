package crldatatypesdomain

import (
	"errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlBooleanURI is the URI that defines the prototype for Boolean
var CrlBooleanURI = CrlDataTypesDomainURI + "/Boolean"

// NewBoolean creates an instance of a Boolean
func NewBoolean(uOfD *core.UniverseOfDiscourse, hl *core.Transaction) core.Concept {
	newBoolean, _ := uOfD.CreateReplicateLiteralAsRefinementFromURI(CrlBooleanURI, hl)
	SetBooleanValue(newBoolean, false, hl)
	return newBoolean
}

// GetBooleanValue returns the Boolean value
func GetBooleanValue(literal core.Concept, hl *core.Transaction) (bool, error) {
	if !IsBoolean(literal, hl) {
		return false, errors.New("GetBooleanValue called with non-Boolean Literal")
	}
	literalValue := literal.GetLiteralValue(hl)
	if literalValue == "true" {
		return true, nil
	} else if literalValue == "false" {
		return false, nil
	}
	return false, errors.New("GetBooleanValue called with non-boolean value in Literal")
}

// IsBoolean returns true if the Literal is a refinement of Boolean
func IsBoolean(literal core.Concept, hl *core.Transaction) bool {
	return literal.IsRefinementOfURI(CrlBooleanURI, hl)
}

// SetBooleanValue sets the value of the Boolean Literal
func SetBooleanValue(literal core.Concept, value bool, hl *core.Transaction) error {
	if !IsBoolean(literal, hl) {
		return errors.New("GetBooleanValue called with non-Boolean Literal")
	}
	if value == true {
		literal.SetLiteralValue("true", hl)
	} else {
		literal.SetLiteralValue("false", hl)
	}
	return nil
}

// BuildCrlBooleanConcept builds the CrlBoolean concept and adds it to the parent space
func BuildCrlBooleanConcept(uOfD *core.UniverseOfDiscourse, parentSpace core.Concept, hl *core.Transaction) {
	crlBoolean, _ := uOfD.NewLiteral(hl, CrlBooleanURI)
	crlBoolean.SetLabel("CrlBoolean", hl)
	crlBoolean.SetOwningConcept(parentSpace, hl)
	crlBoolean.SetIsCore(hl)
}
