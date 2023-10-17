package crldatatypesdomain

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("Boolean test", func() {
	var uOfD *core.UniverseOfDiscourse
	var trans *core.Transaction

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
		BuildCrlDataTypesDomain(uOfD, trans)
	})

	AfterEach(func() {
		trans.ReleaseLocks()
	})

	Specify("Boolean should be created correctly", func() {
		boolean := NewBoolean("", trans)
		Expect(boolean).ToNot(BeNil())
		Expect(boolean.IsRefinementOfURI(CrlBooleanURI, trans)).To(BeTrue())
		value, err := GetBooleanValue(boolean, trans)
		Expect(err).To(BeNil())
		Expect(value).To(Equal(false))
	})

	Specify("SetBooleanValue and GetBooleanValue should work correctly", func() {
		boolean := NewBoolean("", trans)
		Expect(GetBooleanValue(boolean, trans)).To(Equal(false))
		Expect(SetBooleanValue(boolean, true, trans)).To(BeNil())
		Expect(GetBooleanValue(boolean, trans)).To(Equal(true))
		Expect(SetBooleanValue(boolean, false, trans)).To(BeNil())
		Expect(GetBooleanValue(boolean, trans)).To(Equal(false))
	})

	Specify("GetBooleanValue and SetBooleanValue should produce errors if the argument is not a CrlBoolean", func() {
		argument, _ := uOfD.NewLiteral(trans)
		_, err := GetBooleanValue(argument, trans)
		Expect(err).ToNot(BeNil())
		err = SetBooleanValue(argument, true, trans)
		Expect(err).ToNot(BeNil())
	})

	Specify("GetBooleanValue should produce an error if the literal value is neither true or false", func() {
		boolean := NewBoolean("", trans)
		boolean.SetLiteralValue("foo", trans)
		_, err := GetBooleanValue(boolean, trans)
		Expect(err).ToNot(BeNil())
	})
})
