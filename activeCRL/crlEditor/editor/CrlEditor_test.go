package editor_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//	"github.com/sclevine/agouti"
	//	"github.com/pbrown12303/activeCRL/activeCRL/crlEditor/editor"
	am "github.com/sclevine/agouti/matchers"
	//	"testing"
)

var _ = Describe("Test CrlEditor NewDiagram Function:", func() {

	BeforeEach(func() {

	})

	AfterEach(func() {
	})

	Describe("File Menu:", func() {
		Describe("New Diagram Menu Selection: ", func() {
			It("It should create a new diagram ", func() {
				Expect(page.First("#FileMenuButton").Click()).To(Succeed())
				Expect(page.First("#NewDiagramButton").Click()).To(Succeed())
				// Assumption: default diagram name and ID will  be "Diagram1"
				diagram1Tab := page.First("#tabs").FindByButton("Diagram1")
				Expect(diagram1Tab).To(am.HaveCount(1))
				diagram1ViewId, _ := diagram1Tab.Attribute("viewid")
				fmt.Printf("diagram1ViewId: %s \n", diagram1ViewId)
				diagramView := page.First("#" + diagram1ViewId)
				Expect(diagramView).To(am.HaveCount(1))
				var diagramId string
				Expect(page.RunScript("return DiagramManager.GetCrlDiagramIdForContainerId(\""+diagram1ViewId+"\");", nil, &diagramId)).To(Succeed())
				fmt.Printf("diagramId: %s \n", diagramId)
				Expect(diagramId).ToNot(Equal(""))
			})
		})
	})
})

//var _ = Describe("Test RunTest", func() {
//	var result interface{}
//	Expect(page.RunScript("return CrlEditor.RunTest();", nil, &result)).To(Succeed())
//})
