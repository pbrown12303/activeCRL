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
		Expect(boolean.AsCore().IsRefinementOfURI(CrlBooleanURI, trans)).To(BeTrue())
		value, err := boolean.GetBooleanValue(trans)
		Expect(err).To(BeNil())
		Expect(value).To(Equal(false))
	})

	Specify("SetBooleanValue and GetBooleanValue should work correctly", func() {
		boolean := NewBoolean("", trans)
		Expect(boolean.GetBooleanValue(trans)).To(Equal(false))
		Expect(boolean.SetBooleanValue(true, trans)).To(BeNil())
		Expect(boolean.GetBooleanValue(trans)).To(Equal(true))
		Expect(boolean.SetBooleanValue(false, trans)).To(BeNil())
		Expect(boolean.GetBooleanValue(trans)).To(Equal(false))
	})
})
