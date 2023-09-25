package core

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringElementMap Test", func() {
	var uOfD *UniverseOfDiscourse
	var trans *Transaction
	var seMap *StringElementMap
	var el Concept
	var elID string
	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
		seMap = NewStringElementMap()
		el, _ = uOfD.NewElement(trans)
		elID = el.getConceptIDNoLock()
	})
	Specify("Element should not initially show as a member", func() {
		Expect(seMap.GetEntry(elID)).To(BeNil())
	})
	Specify("Element should be found in map after adding", func() {
		seMap.SetEntry(elID, el)
		Expect(seMap.GetEntry(elID)).To(Equal(el))
	})
	Specify("Element should not be in map after removal", func() {
		seMap.SetEntry(elID, el)
		seMap.DeleteEntry(elID)
		Expect(seMap.GetEntry(elID)).To(BeNil())
	})
})
