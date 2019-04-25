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

var _ = Describe("Test CrlEditor", func() {
	var hl *core.HeldLocks
	var cs1ID string
	var cs1 core.Element
	var diagramID string
	var diagram core.Element
	var diagramContainerID string

	BeforeEach(func() {
		hl = uOfD.NewHeldLocks()
		var oldSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID;", nil, &oldSelectionID)).To(Succeed())
		Expect(page.FindByID("FileMenuButton").Click()).To(Succeed())
		Expect(page.FindByID("NewConceptSpaceButton").Click()).To(Succeed())
		Eventually(func() string {
			var retrievedSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
			return retrievedSelectionID
		}, 10).ShouldNot(Equal(oldSelectionID))
		Expect(page.RunScript("return crlSelectedConceptID;", nil, &cs1ID)).To(Succeed())
		Expect(cs1ID).ToNot(Equal(""))
		Eventually(func() bool {
			conceptSpace := editor.CrlEditorSingleton.GetUofD().GetElement(cs1ID)
			return conceptSpace != nil
		}, 10).Should(BeTrue())
		cs1 = uOfD.GetElement(cs1ID)
		Expect(cs1).ToNot(BeNil())
		// At this point the newly created concept space is the selected concept

		// Now add a diagram
		page.RunScript("crlSendAddDiagramChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)
		Eventually(func() string {
			var retrievedSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
			return retrievedSelectionID
		}, 10).ShouldNot(Equal(cs1ID))
		Expect(page.RunScript("return crlSelectedConceptID", nil, &diagramID)).To(Succeed())
		diagram = uOfD.GetElement(diagramID)
		Expect(diagram).ToNot(BeNil())
		Expect(page.RunScript("return crlGetContainerIDFromConceptID(conceptID)", map[string]interface{}{"conceptID": diagramID}, &diagramContainerID)).To(Succeed())

		hl.ReleaseLocksAndWait()
	})

	AfterEach(func() {
		uOfD.DeleteElement(cs1, hl)
		hl.ReleaseLocksAndWait()
	})

	Describe("Testing CrlEditor basic functionality", func() {
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
		PSpecify("Drag TreeNode into Diagram should work", func() {
			// There is a bug in Agouti with respect to both FlickFinger and MoveMouseBy
			// This test will not work until that bug is fixed
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
		Specify("DiagramDrop should produce view of treeDragSelection", func() {
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

	Describe("Test AddChild functionality", func() {
		Specify("AddChild Diagram should work", func() {
			var initialSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
			Expect(page.RunScript("crlSendAddDiagramChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 10).ShouldNot(Equal(initialSelectionID))
			var newDiagramID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &newDiagramID)).To(Succeed())
			newDiagram := uOfD.GetElement(newDiagramID)
			Expect(newDiagram).ToNot(BeNil())
			Expect(newDiagram.IsRefinementOfURI(crldiagram.CrlDiagramURI, hl)).To(BeTrue())
			Expect(newDiagram.GetOwningConcept(hl)).To(Equal(cs1))
		})
		Specify("AddChild Element should work", func() {
			var initialSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
			Expect(page.RunScript("crlSendAddElementChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 10).ShouldNot(Equal(initialSelectionID))
			var newID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &newID)).To(Succeed())
			el := uOfD.GetElement(newID)
			Expect(el).ToNot(BeNil())
			Expect(el.GetOwningConcept(hl)).To(Equal(cs1))
		})
		Specify("AddChild Literal should work", func() {
			var initialSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
			Expect(page.RunScript("crlSendAddLiteralChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 10).ShouldNot(Equal(initialSelectionID))
			var newID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &newID)).To(Succeed())
			el := uOfD.GetElement(newID)
			Expect(el).ToNot(BeNil())
			Expect(el.GetOwningConcept(hl)).To(Equal(cs1))
			isLiteral := false
			switch el.(type) {
			case core.Literal:
				isLiteral = true
			}
			Expect(isLiteral).To(BeTrue())
		})
		Specify("AddChild Reference should work", func() {
			var initialSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
			Expect(page.RunScript("crlSendAddReferenceChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 10).ShouldNot(Equal(initialSelectionID))
			var newID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &newID)).To(Succeed())
			el := uOfD.GetElement(newID)
			Expect(el).ToNot(BeNil())
			Expect(el.GetOwningConcept(hl)).To(Equal(cs1))
			isReference := false
			switch el.(type) {
			case core.Reference:
				isReference = true
			}
			Expect(isReference).To(BeTrue())
		})
		Specify("AddChild Refinement should work", func() {
			var initialSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
			Expect(page.RunScript("crlSendAddRefinementChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 10).ShouldNot(Equal(initialSelectionID))
			var newID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &newID)).To(Succeed())
			el := uOfD.GetElement(newID)
			Expect(el).ToNot(BeNil())
			Expect(el.GetOwningConcept(hl)).To(Equal(cs1))
			isRefinement := false
			switch el.(type) {
			case core.Refinement:
				isRefinement = true
			}
			Expect(isRefinement).To(BeTrue())
		})
	})

	Describe("Test Toolbar Functionality", func() {
		Specify("Element node creation should work", func() {
			var toolbarID string
			Expect(page.RunScript("return crlElementToolbarButtonID", nil, &toolbarID)).To(Succeed())
			Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			var correctToolbarSelection bool
			Expect(page.RunScript("return crlCurrentToolbarButton == crlElementToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
			Expect(correctToolbarSelection).To(BeTrue())
			// Now move mouse to correct position
			// Expect(page.FindByID(diagramContainerID).MouseToElement()).To(Succeed())
			Expect(page.MoveMouseBy(100, 100)).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			Eventually(func() bool {
				var correctToolbarSelection bool
				page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
				return correctToolbarSelection
			}, 60).Should(BeTrue())
			var currentSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
			newNode := uOfD.GetElement(currentSelectionID)
			Expect(newNode).ToNot(BeNil())
			Expect(newNode.GetOwningConcept(hl)).To(Equal(cs1))
			// Check to see that the diagram view of the element has been created correctly
			diagramView := crldiagram.GetFirstElementRepresentingConcept(diagram, newNode, hl)
			Expect(diagramView).ToNot(BeNil())
			Expect(crldiagram.GetReferencedModelElement(diagramView, hl)).To(Equal(newNode))
		})
		Specify("Literal node creation should work", func() {
			var toolbarID string
			Expect(page.RunScript("return crlLiteralToolbarButtonID", nil, &toolbarID)).To(Succeed())
			Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			var correctToolbarSelection bool
			Expect(page.RunScript("return crlCurrentToolbarButton == crlLiteralToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
			Expect(correctToolbarSelection).To(BeTrue())
			// Now move mouse to correct position
			// Expect(page.FindByID(diagramContainerID).MouseToElement()).To(Succeed())
			Expect(page.MoveMouseBy(100, 100)).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			Eventually(func() bool {
				var correctToolbarSelection bool
				page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
				return correctToolbarSelection
			}, 60).Should(BeTrue())
			var currentSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
			newNode := uOfD.GetElement(currentSelectionID)
			Expect(newNode).ToNot(BeNil())
			Expect(newNode.GetOwningConcept(hl)).To(Equal(cs1))
			correctType := false
			switch newNode.(type) {
			case core.Literal:
				correctType = true
			}
			Expect(correctType).To(BeTrue())
			// Check to see that the diagram view of the element has been created correctly
			diagramView := crldiagram.GetFirstElementRepresentingConcept(diagram, newNode, hl)
			Expect(diagramView).ToNot(BeNil())
			Expect(crldiagram.GetReferencedModelElement(diagramView, hl)).To(Equal(newNode))
		})
		Specify("Reference node creation should work", func() {
			var toolbarID string
			Expect(page.RunScript("return crlReferenceToolbarButtonID", nil, &toolbarID)).To(Succeed())
			Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			var correctToolbarSelection bool
			Expect(page.RunScript("return crlCurrentToolbarButton == crlReferenceToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
			Expect(correctToolbarSelection).To(BeTrue())
			// Now move mouse to correct position
			// Expect(page.FindByID(diagramContainerID).MouseToElement()).To(Succeed())
			Expect(page.MoveMouseBy(100, 100)).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			Eventually(func() bool {
				var correctToolbarSelection bool
				page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
				return correctToolbarSelection
			}, 60).Should(BeTrue())
			var currentSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
			newNode := uOfD.GetElement(currentSelectionID)
			Expect(newNode).ToNot(BeNil())
			Expect(newNode.GetOwningConcept(hl)).To(Equal(cs1))
			correctType := false
			switch newNode.(type) {
			case core.Reference:
				correctType = true
			}
			Expect(correctType).To(BeTrue())
			// Check to see that the diagram view of the element has been created correctly
			diagramView := crldiagram.GetFirstElementRepresentingConcept(diagram, newNode, hl)
			Expect(diagramView).ToNot(BeNil())
			Expect(crldiagram.GetReferencedModelElement(diagramView, hl)).To(Equal(newNode))
		})
		Specify("Refinement node creation should work", func() {
			var toolbarID string
			Expect(page.RunScript("return crlRefinementToolbarButtonID", nil, &toolbarID)).To(Succeed())
			Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			var correctToolbarSelection bool
			Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinementToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
			Expect(correctToolbarSelection).To(BeTrue())
			// Now move mouse to correct position
			// Expect(page.FindByID(diagramContainerID).MouseToElement()).To(Succeed())
			Expect(page.MoveMouseBy(100, 100)).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			Eventually(func() bool {
				var correctToolbarSelection bool
				page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
				return correctToolbarSelection
			}, 60).Should(BeTrue())
			var currentSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
			newNode := uOfD.GetElement(currentSelectionID)
			Expect(newNode).ToNot(BeNil())
			Expect(newNode.GetOwningConcept(hl)).To(Equal(cs1))
			correctType := false
			switch newNode.(type) {
			case core.Refinement:
				correctType = true
			}
			Expect(correctType).To(BeTrue())
			// Check to see that the diagram view of the element has been created correctly
			diagramView := crldiagram.GetFirstElementRepresentingConcept(diagram, newNode, hl)
			Expect(diagramView).ToNot(BeNil())
			Expect(crldiagram.GetReferencedModelElement(diagramView, hl)).To(Equal(newNode))
		})
	})
})
