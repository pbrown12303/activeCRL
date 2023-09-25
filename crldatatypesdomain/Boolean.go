package crldatatypesdomain

import (
	"errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlBooleanURI is the URI that defines the prototype for Boolean
var CrlBooleanURI = CrlDataTypesDomainURI + "/Boolean"

// NewBoolean creates an instance of a Boolean
func NewBoolean(uOfD *core.UniverseOfDiscourse, trans *core.Transaction) core.Concept {
	newBoolean, _ := uOfD.CreateReplicateLiteralAsRefinementFromURI(CrlBooleanURI, trans)
	SetBooleanValue(newBoolean, false, trans)
	return newBoolean
}

// GetBooleanValue returns the Boolean value
func GetBooleanValue(literal core.Concept, trans *core.Transaction) (bool, error) {
	if !IsBoolean(literal, trans) {
		return false, errors.New("GetBooleanValue called with non-Boolean Literal")
	}
	literalValue := literal.GetLiteralValue(trans)
	if literalValue == "true" {
		return true, nil
	} else if literalValue == "false" {
		return false, nil
	}
	return false, errors.New("GetBooleanValue called with non-boolean value in Literal")
}

// IsBoolean returns true if the Literal is a refinement of Boolean
func IsBoolean(literal core.Concept, trans *core.Transaction) bool {
	return literal.IsRefinementOfURI(CrlBooleanURI, trans)
}

// SetBooleanValue sets the value of the Boolean Literal
func SetBooleanValue(literal core.Concept, value bool, trans *core.Transaction) error {
	if !IsBoolean(literal, trans) {
		return errors.New("GetBooleanValue called with non-Boolean Literal")
	}
	if value == true {
		literal.SetLiteralValue("true", trans)
	} else {
		literal.SetLiteralValue("false", trans)
	}
	return nil
}

// BuildCrlBooleanConcept builds the CrlBoolean concept and adds it to the parent space
func BuildCrlBooleanConcept(uOfD *core.UniverseOfDiscourse, parentSpace core.Concept, trans *core.Transaction) {
	crlBoolean, _ := uOfD.NewLiteral(trans, CrlBooleanURI)
	crlBoolean.SetLabel("CrlBoolean", trans)
	crlBoolean.SetOwningConcept(parentSpace, trans)
	crlBoolean.SetIsCore(trans)
}
