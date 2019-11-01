package crldatatypes

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
)

var _ = Describe("Boolean test", func() {
	var uOfD *core.UniverseOfDiscourse
	var hl *core.HeldLocks

	BeforeEach(func() {
		uOfD = core.NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
		BuildCrlDataTypesConceptSpace(uOfD, hl)
		hl.ReleaseLocksAndWait()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Specify("Boolean should be created correctly", func() {
		boolean := NewBoolean(uOfD, hl)
		Expect(boolean).ToNot(BeNil())
		Expect(boolean.IsRefinementOfURI(CrlBooleanURI, hl)).To(BeTrue())
		value, err := GetBooleanValue(boolean, hl)
		Expect(err).To(BeNil())
		Expect(value).To(Equal(false))
	})

	Specify("SetBooleanValue and GetBooleanValue should work correctly", func() {
		boolean := NewBoolean(uOfD, hl)
		Expect(GetBooleanValue(boolean, hl)).To(Equal(false))
		Expect(SetBooleanValue(boolean, true, hl)).To(BeNil())
		Expect(GetBooleanValue(boolean, hl)).To(Equal(true))
		Expect(SetBooleanValue(boolean, false, hl)).To(BeNil())
		Expect(GetBooleanValue(boolean, hl)).To(Equal(false))
	})

	Specify("GetBooleanValue and SetBooleanValue should produce errors if the argument is not a CrlBoolean", func() {
		argument, _ := uOfD.NewLiteral(hl)
		_, err := GetBooleanValue(argument, hl)
		Expect(err).ToNot(BeNil())
		err = SetBooleanValue(argument, true, hl)
		Expect(err).ToNot(BeNil())
	})

	Specify("GetBooleanValue should produce an error if the literal value is neither true or false", func() {
		boolean := NewBoolean(uOfD, hl)
		boolean.SetLiteralValue("foo", hl)
		_, err := GetBooleanValue(boolean, hl)
		Expect(err).ToNot(BeNil())
	})
})
