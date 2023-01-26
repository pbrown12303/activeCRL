package core

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var df1URI = "http://dummy.function.uri,df1"

var df1Called = false
var df1CalledElement Element = nil

func dummyChangeFunction(el Element, cn *ChangeNotification, tran *Transaction) error {
	df1Called = true
	df1CalledElement = el
	return nil
}

var _ = Describe("Verify function call graph generation", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction
	var df1 Element

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewTransaction()
		uOfD.AddFunction(df1URI, dummyChangeFunction)
		df1, _ = uOfD.NewElement(hl)
		df1.SetURI(df1URI, hl)
	})

	AfterEach(func() {
		hl.ReleaseLocks()
	})

	Describe("Test FunctionCallGraph for Element ConceptChanged generation", func() {
		Specify("SetDefinition should generate a FunctionCallGraph for ConceptChanged", func() {
			// Initiate the graph capture
			TraceChange = true
			definition := "Definition"
			df1.SetDefinition(definition, hl)
			Expect(df1Called).To(BeTrue())
			Expect(df1CalledElement == df1).To(BeTrue())
			fcgZero := functionCallGraphs[0]
			Expect(fcgZero.executingElement).To(Equal(df1))
			Expect(fcgZero.functionName).To(Equal(df1URI))
			graphString := fcgZero.GetGraph().String()
			Expect(strings.Contains(graphString, "error")).To(BeFalse())
			TraceChange = false
		})
	})
})
