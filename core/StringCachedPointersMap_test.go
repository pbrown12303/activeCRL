package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringCachedPointersMap test", func() {
	var uOfD UniverseOfDiscourse
	var hl *HeldLocks

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Specify("Empty map should handle resolve gracefully", func() {
		scpMap := newStringCachedPointersMap()
		Expect(scpMap.resolveCachedPointers(nil, hl)).ToNot(Succeed())
		el, _ := uOfD.NewElement(hl)
		Expect(scpMap.resolveCachedPointers(el, hl)).To(Succeed())
	})
	Specify("Map should add and resolve pointer properly", func() {
		scpMap := newStringCachedPointersMap()
		el, _ := uOfD.NewElement(hl)
		id := el.getConceptIDNoLock()
		cp := newCachedPointer("dummyOwnerID", false)
		cp.setIndicatedConceptID(id)
		scpMap.addCachedPointer(cp)
		Expect(len(scpMap.scpMap[id])).To(Equal(1))
		Expect(scpMap.scpMap[id][0]).To(Equal(cp))
		scpMap.resolveCachedPointers(el, hl)
		Expect(cp.getIndicatedConcept()).To(Equal(el))
		Expect(len(scpMap.scpMap[id])).To(Equal(0))
	})
	Specify("Map should update owner's ownedConcepts properly", func() {
		scpMap := newStringCachedPointersMap()
		child, _ := uOfD.NewElement(hl)
		owner, _ := uOfD.NewElement(hl)
		ownerID := owner.getConceptIDNoLock()
		cp := newCachedPointer(child.getConceptIDNoLock(), true)
		cp.setIndicatedConceptID(ownerID)
		scpMap.addCachedPointer(cp)
		scpMap.resolveCachedPointers(owner, hl)
		Expect(owner.IsOwnedConcept(child, hl)).To(BeTrue())
	})
	Specify("Map should update listeners properly", func() {
		scpMap := newStringCachedPointersMap()
		child, _ := uOfD.NewElement(hl)
		target, _ := uOfD.NewElement(hl)
		targetID := target.getConceptIDNoLock()
		cp := newCachedPointer(child.getConceptIDNoLock(), false)
		cp.setIndicatedConceptID(targetID)
		scpMap.addCachedPointer(cp)
		scpMap.resolveCachedPointers(target, hl)
		Expect(target.(*element).listeners.GetEntry(child.getConceptIDNoLock())).To(Equal(child))
	})
})
