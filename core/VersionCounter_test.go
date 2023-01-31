package core

import (
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version counter testing", func() {
	var vc *versionCounter
	BeforeEach(func() {
		vc = newVersionCounter()
	})
	Specify("Initial count should be zero", func() {
		Expect(vc.getVersion()).To(Equal(0))
	})
	Specify("Increment should add one to version", func() {
		vc.incrementVersion()
		Expect(vc.getVersion()).To(Equal(1))
		vc.incrementVersion()
		Expect(vc.getVersion()).To(Equal(2))
	})
})
