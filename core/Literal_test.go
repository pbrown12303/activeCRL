package core

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Literal Tests", func() {
	var uOfD *UniverseOfDiscourse
	var trans *Transaction

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
	})

	AfterEach(func() {
		trans.ReleaseLocks()
	})

	Describe("Getting and setting the literal value", func() {
		Specify("Setting the literal value should succeed normally", func() {
			lit, _ := uOfD.NewLiteral(trans)
			testString := "Test"
			Expect(lit.SetLiteralValue(testString, trans)).To(Succeed())
			Expect(lit.GetLiteralValue(trans)).To(Equal(testString))
		})
		Specify("Setting the literal value should fail on a read-only Literal", func() {
			lit, _ := uOfD.NewLiteral(trans)
			testString := "Test"
			lit.SetReadOnly(true, trans)
			Expect(lit.SetLiteralValue(testString, trans)).ToNot(Succeed())
			Expect(lit.GetLiteralValue(trans)).To(Equal(""))
		})
	})

	Describe("Literal equivalence test", func() {
		Specify("Clone should be equivalent on newly initialized literals", func() {
			lit, _ := uOfD.NewLiteral(trans)
			cl := clone(lit, trans)
			Expect(Equivalent(lit, trans, cl, trans)).To(BeTrue())
		})
		Specify("Clone should be equivalent on literals with assigned literal value", func() {
			var testString = "Test"
			lit, _ := uOfD.NewLiteral(trans)
			lit.SetLiteralValue(testString, trans)
			cl := clone(lit, trans)
			Expect(Equivalent(lit, trans, cl, trans)).To(BeTrue())
		})
		Specify("Equivalence should fail if there is a difference in literal value", func() {
			lit, _ := uOfD.NewLiteral(trans)
			cl := clone(lit, trans)
			Expect(Equivalent(lit, trans, cl, trans)).To(BeTrue())
			lit.SetLiteralValue("Test", trans)
			Expect(Equivalent(lit, trans, cl, trans)).To(BeFalse())
		})
		Specify("Equivalence should also fail if there is any difference in the underlying element", func() {
			lit, _ := uOfD.NewLiteral(trans)
			cl := clone(lit, trans)
			Expect(Equivalent(lit, trans, cl, trans)).To(BeTrue())
			lit.Version.counter = 123
			Expect(Equivalent(lit, trans, cl, trans)).To(BeFalse())
		})
	})
	Describe("Marshal and unmarshal tests", func() {
		Specify("Marshal and unmarshal shoudl produce equivalent Literals", func() {
			lit, _ := uOfD.NewLiteral(trans)
			lit.SetLiteralValue("Test string", trans)
			mLit, err1 := lit.MarshalJSON()
			Expect(err1).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewTransaction()
			rLit, err2 := uOfD2.RecoverElement(mLit, hl2)
			Expect(err2).To(BeNil())
			Expect(Equivalent(lit, trans, rLit, hl2)).To(BeTrue())
		})
	})
})
