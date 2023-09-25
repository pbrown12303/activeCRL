package core

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("CrlGraph tests", func() {
	var uOfD *UniverseOfDiscourse
	var trans *Transaction

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		trans = uOfD.NewTransaction()
	})

	AfterEach(func() {
		trans.ReleaseLocks()
	})

	Specify("Creation of core diagram", func() {
		var err error
		// Get the tempDir
		tempDirPath := os.TempDir()
		log.Printf("TempDirPath: " + tempDirPath)
		err = os.Mkdir(tempDirPath, os.ModeDir)
		if !(err == nil || os.IsExist(err)) {
			Expect(err).NotTo(HaveOccurred())
		}
		log.Printf("TempDir created")

		graph := NewCrlGraph("CoreDomain")
		coreDomain := uOfD.GetElementWithURI(CoreDomainURI)
		Expect(graph.AddConceptRecursively(coreDomain, trans)).To(Succeed())
		Expect(graph.ExportDOT(tempDirPath, "CoreDomain")).To(Succeed())
	})
})
