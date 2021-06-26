package core

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify function call graph generation", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction
	var df1 Element
	var df2 Element
	var df3 Element

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
		uOfD.AddFunction(df1URI, dummyChangeFunction)
		uOfD.AddFunction(df2URI, dummyChangeFunction)
		uOfD.AddFunction(df3URI, dummyChangeFunction)
		df1, _ = uOfD.NewElement(hl)
		df1.SetURI(df1URI, hl)
		df2, _ = uOfD.NewElement(hl)
		df2.SetURI(df2URI, hl)
		df3, _ = uOfD.NewElement(hl)
		df3.SetURI(df3URI, hl)
		hl.ReleaseLocksAndWait()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
	})

	Describe("Test FunctionCallGraph for Element ConceptChanged generation", func() {
		Specify("SetDefinition should generate a FunctionCallGraph for ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			// Initiate the graph capture
			TraceChange = true
			definition := "Definition"
			el.SetDefinition(definition, hl)
			hl.ReleaseLocksAndWait()
			Expect(len(functionCallGraphs) > 0).To(BeTrue())
			fcgZero := functionCallGraphs[0]
			Expect(fcgZero.executingElement).To(Equal(el))
			Expect(fcgZero.functionName).To(Equal("http://activeCrl.com/core/coreHousekeeping"))
			graphString := fcgZero.GetGraph().String()
			Expect(strings.Contains(graphString, "error")).To(BeFalse())
			TraceChange = false
		})
	})
})
