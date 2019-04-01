package editor_test

import (
	//	"fmt"

	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagram"
	"github.com/pbrown12303/activeCRL/crleditor/editor"
	"github.com/sclevine/agouti"

	//	"testing"
	. "github.com/sclevine/agouti/matchers"
)

var agoutiDriver *agouti.WebDriver

var _ = Describe("Test CrlEditor", func() {
	var uOfD core.UniverseOfDiscourse
	var hl *core.HeldLocks
	var diagramID string
	var diagram core.Element
	var newDiagramContainerID string

	BeforeEach(func() {
		// Start the editor server
		go editor.StartServer(false)
		// Start the browser
		// Choose a WebDriver:
		// agoutiDriver = agouti.PhantomJS()
		// agoutiDriver = agouti.Selenium()
		agoutiDriver = agouti.ChromeDriver()
		Expect(agoutiDriver.Start()).To(Succeed())

		var err error
		page, err = agoutiDriver.NewPage(agouti.Browser("chrome"))
		Expect(err).NotTo(HaveOccurred())

		Expect(page.Navigate("http://localhost:8082/index")).To(Succeed())
		Expect(page).To(HaveURL("http://localhost:8082/index/"))
		Eventually(func() bool {
			var initializationComplete interface{}
			page.RunScript("return crlInitializationComplete;", nil, &initializationComplete)
			return initializationComplete.(bool)
		}, 20).Should(BeTrue())
		uOfD = editor.CrlEditorSingleton.GetUofD()
		hl = uOfD.NewHeldLocks()
		var oldCurrentDiagramContainerID string
		Expect(page.RunScript("return crlCurrentDiagramContainerID;", nil, &oldCurrentDiagramContainerID)).To(Succeed())
		var result interface{}
		Expect(page.RunScript("crlSendNewDiagramRequest(null)", nil, &result)).To(Succeed())
		Eventually(func() string {
			var retrievedContainerID string
			page.RunScript("return crlCurrentDiagramContainerID;", nil, &retrievedContainerID)
			return retrievedContainerID
		}, 10).ShouldNot(Equal(oldCurrentDiagramContainerID))
		page.RunScript("return crlCurrentDiagramContainerID;", nil, &newDiagramContainerID)
		Expect(newDiagramContainerID).ToNot(Equal(""))
		page.RunScript("return crlGetConceptIDFromContainerID(containerID)", map[string]interface{}{"containerID": newDiagramContainerID}, &diagramID)
		Expect(diagramID).ToNot(Equal(""))
		// time.Sleep(10 * time.Second)
		Eventually(func() bool {
			diagram := editor.CrlEditorSingleton.GetUofD().GetElement(diagramID)
			return diagram != nil
		}, 10).Should(BeTrue())
		diagram = uOfD.GetElement(diagramID)
		hl.ReleaseLocksAndWait()
	})

	AfterEach(func() {
		editor.Exit()
		Expect(agoutiDriver.Stop()).To(Succeed())
	})

	Describe("Testing CrlEditor functionality", func() {
		Specify("The editor should be initialized", func() {
			Expect(editor.CrlEditorSingleton.IsInitialized()).To(BeTrue())
			var initializationComplete interface{}
			page.RunScript("return crlInitializationComplete;", nil, &initializationComplete)
			Expect(initializationComplete).To(BeTrue())
			coreConceptSpace := uOfD.GetElementWithURI(core.CoreConceptSpaceURI)
			var treeNodeID string
			page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
				map[string]interface{}{"conceptID": coreConceptSpace.GetConceptID(hl)},
				&treeNodeID)
			Expect(page.FindByID(treeNodeID)).To(BeFound())
		})
		Specify("NewDiagram should work", func() {
			var oldCurrentDiagramContainerID string
			Expect(page.RunScript("return crlCurrentDiagramContainerID;", nil, &oldCurrentDiagramContainerID)).To(Succeed())
			var result interface{}
			Expect(page.RunScript("crlSendNewDiagramRequest(null)", nil, &result)).To(Succeed())
			var newDiagramContainerID string
			Eventually(func() bool {
				page.RunScript("return crlCurrentDiagramContainerID;", nil, &newDiagramContainerID)
				return newDiagramContainerID != oldCurrentDiagramContainerID
			}).Should(BeTrue())
			Expect(newDiagramContainerID).ToNot(Equal(""))
			var diagramID string
			page.RunScript("return crlGetConceptIDFromContainerID(containerID)", map[string]interface{}{"containerID": newDiagramContainerID}, &diagramID)
			Expect(diagramID).ToNot(Equal(""))
			diagram := uOfD.GetElement(diagramID)
			Expect(diagram).ToNot(BeNil())
		})
		Specify("Tree selection should work", func() {
			coreConceptSpace := uOfD.GetElementWithURI(core.CoreConceptSpaceURI)
			var treeNodeID string
			page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
				map[string]interface{}{"conceptID": coreConceptSpace.GetConceptID(hl)},
				&treeNodeID)
			treeNode := page.FindByID(treeNodeID)
			treeNode.MouseToElement()
			page.Click(agouti.SingleClick, agouti.LeftButton)
			Eventually(func() bool {
				return editor.CrlEditorSingleton.GetCurrentSelection() == coreConceptSpace
			}).Should(BeTrue())
		})
		FSpecify("Drag TreeNode into Diagram should work", func() {
			coreConceptSpace := uOfD.GetElementWithURI(core.CoreConceptSpaceURI)
			var treeNodeID string
			page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
				map[string]interface{}{"conceptID": coreConceptSpace.GetConceptID(hl)},
				&treeNodeID)
			treeNode := page.FindByID(treeNodeID)
			treeNode.MouseToElement()
			pageError := page.Click(agouti.HoldClick, agouti.LeftButton)
			if pageError != nil {
				log.Printf(pageError.Error())
			}
			//				container := page.FindByID(newDiagramContainerID)
			ffError := treeNode.FlickFinger(-100, -300, 50)
			time.Sleep(10 * time.Second)
			if ffError != nil {
				log.Printf(ffError.Error())
			}
			// page.MoveMouseBy(-100, -200)
			Eventually(func() bool {
				return editor.CrlEditorSingleton.GetTreeDragSelection() == coreConceptSpace
			}).Should(BeTrue())
			//				container.MouseToElement()
			page.Click(agouti.ReleaseClick, agouti.LeftButton)
			hl.ReleaseLocksAndWait()
		})
		Specify("DiagramDrop should produce proper results", func() {
			coreConceptSpace := uOfD.GetElementWithURI(core.CoreConceptSpaceURI)
			Expect(page.RunScript("crlSendSetTreeDragSelection(ID)", map[string]interface{}{"ID": coreConceptSpace.GetConceptID(hl)}, nil)).To(Succeed())
			Expect(page.RunScript("crlSendDiagramDrop(ID, x, y)", map[string]interface{}{"ID": diagramID, "x": "100", "y": "100"}, nil)).To(Succeed())
			// Some form of sleep is required here as this thread blocks socket communications. Eventually accomplishes this as it will not
			// be true until after all of the expected client communication has completed.
			Eventually(func() bool {
				return editor.CrlEditorSingleton.GetTreeDragSelection() == nil
			}, 60).Should(BeTrue())
			Expect(len(diagram.GetOwnedConceptsRefinedFromURI(crldiagram.CrlDiagramNodeURI, hl))).To(Equal(1))
			newNode := diagram.GetFirstOwnedConceptRefinedFromURI(crldiagram.CrlDiagramNodeURI, hl)
			Expect(newNode).ToNot(BeNil())
			Expect(newNode.GetLabel(hl)).To(Equal(coreConceptSpace.GetLabel(hl)))
			Expect(crldiagram.GetDisplayLabel(newNode, hl)).To(Equal(coreConceptSpace.GetLabel(hl)))
			// Verify the tree structure
			var treeNodeID string
			Expect(page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
				map[string]interface{}{"conceptID": newNode.GetConceptID(hl)},
				&treeNodeID)).To(Succeed())
			var treeNodeParentID string
			Expect(page.RunScript("return $(\"#uOfD\").jstree(true).get_parent(treeNodeID);",
				map[string]interface{}{"treeNodeID": treeNodeID},
				&treeNodeParentID)).To(Succeed())
			var diagramTreeNodeID string
			Expect(page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
				map[string]interface{}{"conceptID": diagramID},
				&diagramTreeNodeID)).To(Succeed())
			Expect(treeNodeParentID).To(Equal(diagramTreeNodeID))
			// Now drop a second instance
			Expect(page.RunScript("crlSendSetTreeDragSelection(ID)", map[string]interface{}{"ID": coreConceptSpace.GetConceptID(hl)}, nil)).To(Succeed())
			Expect(page.RunScript("crlSendDiagramDrop(ID, x, y)", map[string]interface{}{"ID": diagramID, "x": "200", "y": "200"}, nil)).To(Succeed())
			// Some form of sleep is required here as this thread blocks socket communications. Eventually accomplishes this as it will not
			// be true until after all of the expected client communication has completed.
			Eventually(func() bool {
				return editor.CrlEditorSingleton.GetTreeDragSelection() == nil
			}, 60).Should(BeTrue())
			Expect(len(diagram.GetOwnedConceptsRefinedFromURI(crldiagram.CrlDiagramNodeURI, hl))).To(Equal(2))
			var newNode2 core.Element
			for _, el := range diagram.GetOwnedConceptsRefinedFromURI(crldiagram.CrlDiagramNodeURI, hl) {
				if el != newNode {
					newNode2 = el
				}
			}
			Expect(newNode2).ToNot(BeNil())
			Expect(newNode2.GetLabel(hl)).To(Equal(coreConceptSpace.GetLabel(hl)))
			Expect(crldiagram.GetDisplayLabel(newNode2, hl)).To(Equal(coreConceptSpace.GetLabel(hl)))
			hl.ReleaseLocksAndWait()
			// Verify the tree structure
			var treeNode2ID string
			Expect(page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
				map[string]interface{}{"conceptID": newNode2.GetConceptID(hl)},
				&treeNode2ID)).To(Succeed())
			var treeNode2ParentID string
			Expect(page.RunScript("return $(\"#uOfD\").jstree(true).get_parent(treeNodeID);",
				map[string]interface{}{"treeNodeID": treeNode2ID},
				&treeNode2ParentID)).To(Succeed())
			Expect(treeNode2ParentID).To(Equal(diagramTreeNodeID))
		})
	})
})
