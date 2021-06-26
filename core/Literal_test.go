package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Literal Tests", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("Getting and setting the literal value", func() {
		Specify("Setting the literal value should succeed normally", func() {
			lit, _ := uOfD.NewLiteral(hl)
			testString := "Test"
			Expect(lit.SetLiteralValue(testString, hl)).To(Succeed())
			Expect(lit.GetLiteralValue(hl)).To(Equal(testString))
		})
		Specify("Setting the literal value should fail on a read-only Literal", func() {
			lit, _ := uOfD.NewLiteral(hl)
			testString := "Test"
			lit.SetReadOnly(true, hl)
			Expect(lit.SetLiteralValue(testString, hl)).ToNot(Succeed())
			Expect(lit.GetLiteralValue(hl)).To(Equal(""))
		})
	})

	Describe("Literal equivalence test", func() {
		Specify("Clone should be equivalent on newly initialized literals", func() {
			lit, _ := uOfD.NewLiteral(hl)
			cl := clone(lit, hl)
			Expect(Equivalent(lit, hl, cl, hl)).To(BeTrue())
		})
		Specify("Clone should be equivalent on literals with assigned literal value", func() {
			var testString = "Test"
			lit, _ := uOfD.NewLiteral(hl)
			lit.SetLiteralValue(testString, hl)
			cl := clone(lit, hl)
			Expect(Equivalent(lit, hl, cl, hl)).To(BeTrue())
		})
		Specify("Equivalence should fail if there is a difference in literal value", func() {
			lit, _ := uOfD.NewLiteral(hl)
			cl := clone(lit, hl)
			Expect(Equivalent(lit, hl, cl, hl)).To(BeTrue())
			lit.SetLiteralValue("Test", hl)
			Expect(Equivalent(lit, hl, cl, hl)).To(BeFalse())
		})
		Specify("Equivalence should also fail if there is any difference in the underlying element", func() {
			lit, _ := uOfD.NewLiteral(hl)
			cl := clone(lit, hl)
			Expect(Equivalent(lit, hl, cl, hl)).To(BeTrue())
			lit.(*literal).Version.counter = 123
			Expect(Equivalent(lit, hl, cl, hl)).To(BeFalse())
		})
	})
	Describe("Marshal and unmarshal tests", func() {
		Specify("Marshal and unmarshal shoudl produce equivalent Literals", func() {
			lit, _ := uOfD.NewLiteral(hl)
			lit.SetLiteralValue("Test string", hl)
			mLit, err1 := lit.MarshalJSON()
			Expect(err1).To(BeNil())
			uOfD2 := NewUniverseOfDiscourse()
			hl2 := uOfD2.NewHeldLocks()
			rLit, err2 := uOfD2.RecoverElement(mLit, hl2)
			Expect(err2).To(BeNil())
			Expect(Equivalent(lit, hl, rLit, hl2)).To(BeTrue())
		})
	})
})
