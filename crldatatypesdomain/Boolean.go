package crldatatypesdomain

import (
	"errors"

	"github.com/pbrown12303/activeCRL/core"
)

// CrlBooleanURI is the URI that defines the prototype for Boolean
var CrlBooleanURI = CrlDataTypesDomainURI + "/Boolean"

// CrlBoolean is the crl representation of a boolean
type CrlBoolean core.Concept

// NewBoolean creates an instance of a Boolean
func NewBoolean(label string, trans *core.Transaction, newURI ...string) *CrlBoolean {
	uOfD := trans.GetUniverseOfDiscourse()
	newLiteral, _ := uOfD.CreateRefinementOfConceptURI(CrlBooleanURI, "CrlBoolean", trans, newURI...)
	newBoolean := (*CrlBoolean)(newLiteral)
	newBoolean.SetBooleanValue(false, trans)
	newBoolean.ToCore().SetLabel(label, trans)
	return newBoolean
}

// NewOwnedBoolean creates a refinement of the Boolean concept and sets both its label and owner
func NewOwnedBoolean(owner *core.Concept, label string, trans *core.Transaction, newURI ...string) *CrlBoolean {
	newBoolean := NewBoolean(label, trans, newURI...)
	newBoolean.ToCore().SetOwningConcept(owner, trans)
	return newBoolean
}

// ToCore casts CrlBoolean to core.Concept
func (crlb *CrlBoolean) ToCore() *core.Concept {
	return (*core.Concept)(crlb)
}

// GetBooleanValue returns the Boolean value
func (crlb *CrlBoolean) GetBooleanValue(trans *core.Transaction) (bool, error) {
	if !IsBoolean(crlb.ToCore(), trans) {
		return false, errors.New("GetBooleanValue called with non-Boolean Literal")
	}
	literalValue := crlb.ToCore().GetLiteralValue(trans)
	if literalValue == "true" {
		return true, nil
	} else if literalValue == "false" {
		return false, nil
	}
	return false, errors.New("GetBooleanValue called with non-boolean value in Literal")
}

// IsBoolean returns true if the Literal is a refinement of Boolean
func IsBoolean(literal *core.Concept, trans *core.Transaction) bool {
	return literal.IsRefinementOfURI(CrlBooleanURI, trans)
}

// SetBooleanValue sets the value of the Boolean Literal
func (crlb *CrlBoolean) SetBooleanValue(value bool, trans *core.Transaction) error {
	if !IsBoolean(crlb.ToCore(), trans) {
		return errors.New("SetBooleanValue called with non-Boolean Literal")
	}
	if value == true {
		crlb.ToCore().SetLiteralValue("true", trans)
	} else {
		crlb.ToCore().SetLiteralValue("false", trans)
	}
	return nil
}

// BuildCrlBooleanConcept builds the CrlBoolean concept and adds it to the parent space
func BuildCrlBooleanConcept(uOfD *core.UniverseOfDiscourse, parentSpace *core.Concept, trans *core.Transaction) {
	crlBoolean, _ := uOfD.NewLiteral(trans, CrlBooleanURI)
	crlBoolean.SetLabel("CrlBoolean", trans)
	crlBoolean.SetOwningConcept(parentSpace, trans)
	crlBoolean.SetIsCore(trans)
}
