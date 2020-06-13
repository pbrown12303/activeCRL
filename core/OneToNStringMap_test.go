package core

import (
	. "github.com/onsi/ginkgo"
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
			onsMap1.AddMappedValue("A", "X")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.AddMappedValue("A", "X")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.AddMappedValue("A", "Y")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.AddMappedValue("A", "Y")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.AddMappedValue("B", "Z")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.AddMappedValue("B", "Z")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.RemoveMappedValue("A", "X")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.RemoveMappedValue("A", "X")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.RemoveMappedValue("A", "Y")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.RemoveMappedValue("A", "Y")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
			onsMap1.RemoveMappedValue("B", "Z")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeFalse())
			onsMap2.RemoveMappedValue("B", "Z")
			Expect(onsMap1.IsEquivalent(onsMap2)).To(BeTrue())
		})
	})
})
