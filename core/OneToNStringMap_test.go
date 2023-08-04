package core

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("OneToNStringMap Test", func() {
	Describe("IsEquivalent should work", func() {
		Specify("Empty sets should be equivalent", func() {
			onsMap1 := NewOneToNStringMap()
			onsMap2 := NewOneToNStringMap()
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
		})
		Specify("IsEquivalent should reflect actual differences", func() {
			onsMap1 := NewOneToNStringMap()
			onsMap2 := NewOneToNStringMap()
			onsMap1.addMappedValue("A", "X")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.addMappedValue("A", "X")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.addMappedValue("A", "Y")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.addMappedValue("A", "Y")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.addMappedValue("B", "Z")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.addMappedValue("B", "Z")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.removeMappedValue("A", "X")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.removeMappedValue("A", "X")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.removeMappedValue("A", "Y")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.removeMappedValue("A", "Y")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.removeMappedValue("B", "Z")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.removeMappedValue("B", "Z")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
		})
	})
})
