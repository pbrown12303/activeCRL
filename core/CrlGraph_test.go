package core

import (
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CrlGraph tests", func() {
	var uOfD *UniverseOfDiscourse
	var hl *HeldLocks

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
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
		Expect(graph.AddConceptRecursively(coreDomain, hl)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Expect(graph.ExportDOT(tempDirPath, "CoreDomain")).To(Succeed())
	})
})
