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
	var diagramGraphID string

	BeforeEach(func() {
		hl = uOfD.NewHeldLocks()
		var oldSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID;", nil, &oldSelectionID)).To(Succeed())
		Expect(page.FindByID("FileMenuButton").Click()).To(Succeed())
		Expect(page.FindByID("NewConceptSpaceButton").Click()).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() string {
			var retrievedSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
			return retrievedSelectionID
		}, 3).ShouldNot(Equal(oldSelectionID))
		Expect(page.RunScript("return crlSelectedConceptID;", nil, &cs1ID)).To(Succeed())
		Expect(cs1ID).ToNot(Equal(""))
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			conceptSpace := editor.CrlEditorSingleton.GetUofD().GetElement(cs1ID)
			return conceptSpace != nil
		}, 3).Should(BeTrue())
		cs1 = uOfD.GetElement(cs1ID)
		Expect(cs1).ToNot(BeNil())
		// At this point the newly created concept space is the selected concept

		// Now add a diagram
		page.RunScript("crlSendAddDiagramChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)
		hl.ReleaseLocksAndWait()
		Eventually(func() string {
			var retrievedSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
			return retrievedSelectionID
		}, 3).ShouldNot(Equal(cs1ID))
		Expect(page.RunScript("return crlSelectedConceptID", nil, &diagramID)).To(Succeed())
		diagram = uOfD.GetElement(diagramID)
		Expect(diagram).ToNot(BeNil())
		Expect(page.RunScript("return crlGetContainerIDFromConceptID(conceptID)", map[string]interface{}{"conceptID": diagramID}, &diagramContainerID)).To(Succeed())
		Expect(page.RunScript("return crlGetJointGraphIDFromDiagramID(diagramID)", map[string]interface{}{"diagramID": diagramID}, &diagramGraphID)).To(Succeed())
		hl.ReleaseLocksAndWait()
	})

	AfterEach(func() {
		hl.ReleaseLocksAndWait()
		uOfD.DeleteElement(cs1, hl)
		hl.ReleaseLocksAndWait()
	})

	GetCellViewIDFromViewElementID := func(viewElementID string) string {
		// First get the expected crl ID for the cell in the graph
		var crlCellID string
		Expect(page.RunScript("return crlGetJointCellIDFromConceptID(conceptID);",
			map[string]interface{}{"conceptID": viewElementID},
			&crlCellID)).To(Succeed())
		// Now find the DOM id for the cell in the graph
		var cellID string
		Expect(page.RunScript("return crlFindCellInGraphID(graphID, crlCellID).id",
			map[string]interface{}{"graphID": diagramGraphID, "crlCellID": crlCellID},
			&cellID)).To(Succeed())
		// Finally, find the cell view id on the paper
		var cellViewID string
		Expect(page.RunScript("return crlFindCellViewInPaperByDiagramID(diagramID, cellID).id",
			map[string]interface{}{"diagramID": diagramID, "cellID": cellID},
			&cellViewID)).To(Succeed())
		return cellViewID
	}
	GetCurrentSelection := func() core.Element {
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		return uOfD.GetElement(currentSelectionID)
	}
	MouseToDiagramPosition := func(x int, y int) {
		var mousePosition map[string]interface{}
		Expect(page.RunScript("return crlMousePosition;", nil, &mousePosition))
		currentMouseX := int(mousePosition["x"].(float64))
		currentMouseY := int(mousePosition["y"].(float64))
		var pageOffset map[string]interface{}
		Expect(page.RunScript("var jointPaperID = crlGetJointPaperIDFromDiagramID(diagramID); return crlPapersGlobal[jointPaperID].pageOffset();",
			map[string]interface{}{"diagramID": diagramID}, &pageOffset)).To(Succeed())
		pageX := int(pageOffset["x"].(float64))
		pageY := int(pageOffset["y"].(float64))
		xMove := pageX + x - currentMouseX
		yMove := pageY + y - currentMouseY
		Expect(page.MoveMouseBy(xMove, yMove)).To(Succeed())
	}
	CreateElement := func(x int, y int) (core.Element, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlElementToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlElementToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move mouse to correct position
		MouseToDiagramPosition(x, y)
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		newConcept := uOfD.GetElement(currentSelectionID)
		Expect(newConcept).ToNot(BeNil())
		Expect(newConcept.GetOwningConcept(hl)).To(Equal(cs1))
		// Check to see that the diagram view of the element has been created correctly
		conceptView := crldiagram.GetFirstElementRepresentingConcept(diagram, newConcept, hl)
		Expect(conceptView).ToNot(BeNil())
		Expect(crldiagram.GetReferencedModelElement(conceptView, hl)).To(Equal(newConcept))
		hl.ReleaseLocksAndWait()
		return newConcept, conceptView
	}
	CreateLiteral := func(x int, y int) (core.Literal, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlLiteralToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlLiteralToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move mouse to correct position
		MouseToDiagramPosition(x, y)
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		newConcept := uOfD.GetElement(currentSelectionID)
		Expect(newConcept).ToNot(BeNil())
		Expect(newConcept.GetOwningConcept(hl)).To(Equal(cs1))
		correctType := false
		switch newConcept.(type) {
		case core.Literal:
			correctType = true
		}
		Expect(correctType).To(BeTrue())
		// Check to see that the diagram view of the element has been created correctly
		conceptView := crldiagram.GetFirstElementRepresentingConcept(diagram, newConcept, hl)
		Expect(conceptView).ToNot(BeNil())
		Expect(crldiagram.GetReferencedModelElement(conceptView, hl)).To(Equal(newConcept))
		hl.ReleaseLocksAndWait()
		return newConcept.(core.Literal), conceptView
	}
	CreateReferenceNode := func(x int, y int) (core.Reference, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlReferenceToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlReferenceToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move mouse to correct position
		MouseToDiagramPosition(x, y)
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			Expect(page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		newConcept := uOfD.GetElement(currentSelectionID)
		Expect(newConcept).ToNot(BeNil())
		Expect(newConcept.GetOwningConcept(hl)).To(Equal(cs1))
		correctType := false
		switch newConcept.(type) {
		case core.Reference:
			correctType = true
		}
		Expect(correctType).To(BeTrue())
		// Check to see that the diagram view of the element has been created correctly
		conceptView := crldiagram.GetFirstElementRepresentingConcept(diagram, newConcept, hl)
		Expect(conceptView).ToNot(BeNil())
		Expect(crldiagram.GetReferencedModelElement(conceptView, hl)).To(Equal(newConcept))
		hl.ReleaseLocksAndWait()
		return newConcept.(core.Reference), conceptView
	}
	CreateReferenceLink := func(sourceView core.Element, targetView core.Element) (core.Reference, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlReferenceLinkToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlReferenceLinkToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to source, click, drag to target, and release
		targetCellID := GetCellViewIDFromViewElementID(targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		newElement := GetCurrentSelection()
		Expect(newElement).ToNot(BeNil())
		correctType := false
		switch newElement.(type) {
		case core.Reference:
			correctType = true
		}
		Expect(correctType).To(BeTrue())
		newReference := newElement.(core.Reference)
		source := crldiagram.GetReferencedModelElement(sourceView, hl)
		Expect(newReference.GetOwningConceptID(hl)).To(Equal(source.GetConceptID(hl)))
		newReferenceView := crldiagram.GetFirstElementRepresentingConcept(diagram, newReference, hl)
		hl.ReleaseLocksAndWait()
		return newReference, newReferenceView
	}
	CreateRefinementNode := func(x int, y int) (core.Refinement, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlRefinementToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinementToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move mouse to correct position
		MouseToDiagramPosition(x, y)
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		newConcept := uOfD.GetElement(currentSelectionID)
		Expect(newConcept).ToNot(BeNil())
		Expect(newConcept.GetOwningConcept(hl)).To(Equal(cs1))
		correctType := false
		switch newConcept.(type) {
		case core.Refinement:
			correctType = true
		}
		Expect(correctType).To(BeTrue())
		// Check to see that the diagram view of the element has been created correctly
		conceptView := crldiagram.GetFirstElementRepresentingConcept(diagram, newConcept, hl)
		Expect(conceptView).ToNot(BeNil())
		Expect(crldiagram.GetReferencedModelElement(conceptView, hl)).To(Equal(newConcept))
		Expect(newConcept.GetOwningConceptID(hl)).To(Equal(cs1.GetConceptID(hl)))
		hl.ReleaseLocksAndWait()
		return newConcept.(core.Refinement), conceptView
	}
	CreateRefinementLink := func(sourceView core.Element, targetView core.Element) (core.Refinement, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlRefinementLinkToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinementLinkToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to e2, click, drag to e1, and release
		targetCellID := GetCellViewIDFromViewElementID(targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		// Now check the results
		newElement := GetCurrentSelection()
		Expect(newElement).ToNot(BeNil())
		correctType := false
		switch newElement.(type) {
		case core.Refinement:
			correctType = true
		}
		Expect(correctType).To(BeTrue())
		newRefinement := newElement.(core.Refinement)
		newRefinementView := crldiagram.GetFirstElementRepresentingConcept(diagram, newRefinement, hl)
		Expect(newRefinement.GetOwningConceptID(hl)).To(Equal(cs1.GetConceptID(hl)))
		hl.ReleaseLocksAndWait()
		return newRefinement, newRefinementView
	}
	CreateOwnerPointer := func(sourceView core.Element, targetView core.Element) (core.Element, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlOwnerPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlOwnerPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to source, click, drag to target, and release
		targetCellID := GetCellViewIDFromViewElementID(targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		source := crldiagram.GetReferencedModelElement(sourceView, hl)
		ownerPointerView := crldiagram.GetFirstElementRepresentingConceptOwnerPointer(diagram, source, hl)
		hl.ReleaseLocksAndWait()
		return source, ownerPointerView
	}
	CreateElementPointer := func(sourceView core.Element, targetView core.Element) (core.Reference, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlElementPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlElementPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to r1, click, drag to e1, and release
		targetCellID := GetCellViewIDFromViewElementID(targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		referenceID := crldiagram.GetReferencedModelElement(sourceView, hl).GetConceptID(hl)
		reference := uOfD.GetReference(referenceID)
		elementPointerView := crldiagram.GetFirstElementRepresentingConceptElementPointer(diagram, reference, hl)
		hl.ReleaseLocksAndWait()
		return reference, elementPointerView
	}

	CreateAbstractPointer := func(sourceView core.Element, targetView core.Element) (core.Refinement, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlAbstractPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlAbstractPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to r1, click, drag to e1, and release
		targetCellID := GetCellViewIDFromViewElementID(targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		refinementID := crldiagram.GetReferencedModelElement(sourceView, hl).GetConceptID(hl)
		refinement := uOfD.GetRefinement(refinementID)
		elementPointerView := crldiagram.GetFirstElementRepresentingConceptAbstractPointer(diagram, refinement, hl)
		hl.ReleaseLocksAndWait()
		return refinement, elementPointerView
	}

	CreateRefinedPointer := func(sourceView core.Element, targetView core.Element) (core.Refinement, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlRefinedPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinedPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to r1, click, drag to e1, and release
		targetCellID := GetCellViewIDFromViewElementID(targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		hl.ReleaseLocksAndWait()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		refinementID := crldiagram.GetReferencedModelElement(sourceView, hl).GetConceptID(hl)
		refinement := uOfD.GetRefinement(refinementID)
		elementPointerView := crldiagram.GetFirstElementRepresentingConceptRefinedPointer(diagram, refinement, hl)
		hl.ReleaseLocksAndWait()
		return refinement, elementPointerView
	}

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
			hl.ReleaseLocksAndWait()
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
			hl.ReleaseLocksAndWait()
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
			hl.ReleaseLocksAndWait()
			Eventually(func() bool {
				return editor.CrlEditorSingleton.GetTreeDragSelection() == nil
			}, 3).Should(BeTrue())
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
			hl.ReleaseLocksAndWait()
			Eventually(func() bool {
				return editor.CrlEditorSingleton.GetTreeDragSelection() == nil
			}, 3).Should(BeTrue())
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
			hl.ReleaseLocksAndWait()
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 3).ShouldNot(Equal(initialSelectionID))
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
			hl.ReleaseLocksAndWait()
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 3).ShouldNot(Equal(initialSelectionID))
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
			hl.ReleaseLocksAndWait()
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 3).ShouldNot(Equal(initialSelectionID))
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
			hl.ReleaseLocksAndWait()
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 3).ShouldNot(Equal(initialSelectionID))
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
			hl.ReleaseLocksAndWait()
			Eventually(func() string {
				var retrievedSelectionID string
				Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
				return retrievedSelectionID
			}, 3).ShouldNot(Equal(initialSelectionID))
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
			e1, e1View := CreateElement(100, 100)
			Expect(e1).ToNot(BeNil())
			Expect(e1View).ToNot(BeNil())
		})
		Specify("Literal node creation should work", func() {
			l1, l1View := CreateLiteral(100, 100)
			Expect(l1).ToNot(BeNil())
			Expect(l1View).ToNot(BeNil())
		})
		Specify("Reference node creation should work", func() {
			r1, r1View := CreateReferenceNode(100, 100)
			Expect(r1).ToNot(BeNil())
			Expect(r1View).ToNot(BeNil())
		})
		Specify("Refinement node creation should work", func() {
			// editor.CrlLogClientDialog = true
			r1, r1View := CreateRefinementNode(100, 100)
			Expect(r1).ToNot(BeNil())
			Expect(r1View).ToNot(BeNil())
		})
		Describe("Reference link creation should work", func() {
			Specify("for a node source and target", func() {
				e1, e1View := CreateElement(100, 100)
				e2, e2View := CreateElement(100, 200)
				newRefinement, _ := CreateReferenceLink(e2View, e1View)
				Expect(newRefinement.GetOwningConcept(hl)).To(Equal(e2))
				Expect(newRefinement.GetReferencedConcept(hl)).To(Equal(e1))
			})
			Specify("for a link source and node target", func() {
				_, e1View := CreateElement(100, 100)
				_, e2View := CreateElement(100, 200)
				// create the node target
				e3, e3View := CreateElement(200, 150)
				// Create a reference link
				refLink1, refLink1View := CreateReferenceLink(e2View, e1View)
				// Now the new reference
				refLink2, _ := CreateReferenceLink(refLink1View, e3View)
				// Now check the results
				Expect(refLink2.GetOwningConceptID(hl)).To(Equal(refLink1.GetConceptID(hl)))
				Expect(refLink2.GetReferencedConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
			})
			Specify("for a node source and link target", func() {
				_, e1View := CreateElement(100, 100)
				_, e2View := CreateElement(100, 200)
				// create the node source
				e3, e3View := CreateElement(200, 150)
				// Create a reference link
				refLink1, refLink1View := CreateReferenceLink(e2View, e1View)
				// Now the new reference
				refLink2, _ := CreateReferenceLink(e3View, refLink1View)
				// Now check the results
				Expect(refLink2.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
				Expect(refLink2.GetReferencedConceptID(hl)).To(Equal(refLink1.GetConceptID(hl)))
			})
			Specify("for a link source and link target", func() {
				_, e1View := CreateElement(100, 100)
				_, e2View := CreateElement(100, 200)
				_, e3View := CreateElement(200, 100)
				_, e4View := CreateElement(200, 200)
				// Create the source reference link
				refLink1, refLink1View := CreateReferenceLink(e2View, e1View)
				// Create the target reference link
				refLink2, refLink2View := CreateReferenceLink(e4View, e3View)
				// Now the new reference
				refLink3, _ := CreateReferenceLink(refLink1View, refLink2View)
				// Now check the results
				Expect(refLink3.GetOwningConceptID(hl)).To(Equal(refLink1.GetConceptID(hl)))
				Expect(refLink3.GetReferencedConceptID(hl)).To(Equal(refLink2.GetConceptID(hl)))
			})
			Specify("for a node source and an OwnerPointer target", func() {
				_, e1View := CreateElement(100, 100)
				e2, e2View := CreateElement(100, 200)
				// create the node source
				e3, e3View := CreateElement(200, 150)
				// create the owner pointer
				opModelElement, opView := CreateOwnerPointer(e2View, e1View)
				// Create the Reference
				ref, refView := CreateReferenceLink(e3View, opView)
				Expect(opModelElement.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(ref.GetReferencedConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(ref.GetReferencedAttributeName(hl)).To(Equal(core.OwningConceptID))
				Expect(ref.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(refView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(refView, hl).GetConceptID(hl)).To(Equal(opView.GetConceptID(hl)))
			})
			Specify("for a node source and an ElementPointer target", func() {
				_, e1View := CreateElement(100, 100)
				e2, e2View := CreateReferenceNode(100, 200)
				// create the node source
				e3, e3View := CreateElement(200, 150)
				// create the owner pointer
				epModelElement, epView := CreateElementPointer(e2View, e1View)
				// Create the Reference
				ref, refView := CreateReferenceLink(e3View, epView)
				Expect(epModelElement.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(ref.GetReferencedConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(ref.GetReferencedAttributeName(hl)).To(Equal(core.ReferencedConceptID))
				Expect(ref.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(refView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(refView, hl).GetConceptID(hl)).To(Equal(epView.GetConceptID(hl)))
			})
			Specify("for a node source and an AbstractPointer target", func() {
				_, e1View := CreateElement(100, 100)
				e2, e2View := CreateRefinementNode(100, 200)
				// create the node source
				e3, e3View := CreateElement(200, 150)
				// create the owner pointer
				apModelElement, apView := CreateAbstractPointer(e2View, e1View)
				// Create the Reference
				ref, refView := CreateReferenceLink(e3View, apView)
				Expect(apModelElement.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(ref.GetReferencedConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(ref.GetReferencedAttributeName(hl)).To(Equal(core.AbstractConceptID))
				Expect(ref.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(refView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(refView, hl).GetConceptID(hl)).To(Equal(apView.GetConceptID(hl)))
			})
			Specify("for a node source and an RefinedPointer target", func() {
				_, e1View := CreateElement(100, 100)
				e2, e2View := CreateRefinementNode(100, 200)
				// create the node source
				e3, e3View := CreateElement(200, 150)
				// create the owner pointer
				apModelElement, apView := CreateRefinedPointer(e2View, e1View)
				// Create the Reference
				ref, refView := CreateReferenceLink(e3View, apView)
				Expect(apModelElement.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(ref.GetReferencedConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(ref.GetReferencedAttributeName(hl)).To(Equal(core.RefinedConceptID))
				Expect(ref.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(refView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(refView, hl).GetConceptID(hl)).To(Equal(apView.GetConceptID(hl)))
			})
		})
		Describe("Refinement link creation should work", func() {
			Specify("for a node source and node target", func() {
				e1, e1View := CreateElement(100, 100)
				e2, e2View := CreateElement(100, 200)
				newRefinement, _ := CreateRefinementLink(e2View, e1View)
				Expect(newRefinement.GetAbstractConcept(hl)).To(Equal(e1))
				Expect(newRefinement.GetRefinedConcept(hl)).To(Equal(e2))
			})
			Specify("for a link source and node target", func() {
				_, e1View := CreateElement(100, 100)
				_, e2View := CreateElement(100, 200)
				source, sourceView := CreateRefinementLink(e2View, e1View)
				target, targetView := CreateRefinementNode(200, 150)
				newRefinement, newRefinementView := CreateRefinementLink(sourceView, targetView)
				Expect(newRefinement.GetAbstractConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(newRefinement.GetRefinedConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(newRefinementView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(newRefinementView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
			Specify("for a node source and link target", func() {
				_, e1View := CreateElement(100, 100)
				_, e2View := CreateElement(100, 200)
				target, targetView := CreateRefinementLink(e2View, e1View)
				source, sourceView := CreateRefinementNode(200, 150)
				newRefinement, newRefinementView := CreateRefinementLink(sourceView, targetView)
				Expect(newRefinement.GetAbstractConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(newRefinement.GetRefinedConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(newRefinementView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(newRefinementView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
		})
		Describe("OwnerPointer creation should work", func() {
			Specify("For a node source and node target", func() {
				e1, e1View := CreateElement(100, 100)
				e2, e2View := CreateElement(100, 200)
				source, ownerPointerView := CreateOwnerPointer(e2View, e1View)
				// Now check the results
				Expect(source.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
				Expect(e2.GetOwningConceptID(hl)).To(Equal(e1.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(ownerPointerView, hl).GetConceptID(hl)).To(Equal(e2View.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(ownerPointerView, hl).GetConceptID(hl)).To(Equal(e1View.GetConceptID(hl)))
			})
			Specify("For a Refinement Link source and node target", func() {
				_, e1View := CreateElement(100, 100)
				_, e2View := CreateElement(100, 200)
				// Create the refinement link
				source, sourceView := CreateRefinementLink(e1View, e2View)
				// create the new owner
				e3, e3View := CreateElement(200, 150)
				// Now the ownerPointer
				ownerPointerConcept, ownerPointerView := CreateOwnerPointer(sourceView, e3View)
				// Now check the results
				Expect(source.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
				Expect(source.GetConceptID(hl)).To(Equal(ownerPointerConcept.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(ownerPointerView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(ownerPointerView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
			})
			Specify("For a node source and ReferenceLink target", func() {
				_, e1View := CreateElement(100, 100)
				_, e2View := CreateElement(100, 200)
				// Create the Reference link
				target, targetView := CreateReferenceLink(e1View, e2View)
				// create the new owner
				source, sourceView := CreateElement(200, 150)
				// Now the ownerPointer
				ownerPointerConcept, ownerPointerView := CreateOwnerPointer(sourceView, targetView)
				// Now check the results
				Expect(source.GetOwningConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetConceptID(hl)).To(Equal(ownerPointerConcept.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(ownerPointerView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(ownerPointerView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
			Specify("For a node source and RefinementLink target", func() {
				_, e1View := CreateElement(100, 100)
				_, e2View := CreateElement(100, 200)
				// Create the refinement link
				target, targetView := CreateRefinementLink(e1View, e2View)
				// create the new owner
				source, sourceView := CreateElement(200, 150)
				// Now the ownerPointer
				ownerPointerConcept, ownerPointerView := CreateOwnerPointer(sourceView, targetView)
				// Now check the results
				Expect(source.GetOwningConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetConceptID(hl)).To(Equal(ownerPointerConcept.GetConceptID(hl)))
				Expect(crldiagram.GetLinkSource(ownerPointerView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(ownerPointerView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
		})
		Describe("ElementPointer creation should work", func() {
			Specify("for a node source and node target", func() {
				target, targetView := CreateElement(100, 100)
				source, sourceView := CreateReferenceNode(100, 200)
				reference, epView := CreateElementPointer(sourceView, targetView)
				// Now check the results
				Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(reference.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(reference.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.NoAttribute))
				Expect(crldiagram.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
			Specify("for a node source and reference link target", func() {
				source, sourceView := CreateReferenceNode(100, 150)
				_, e1View := CreateElement(200, 100)
				_, e2View := CreateElement(200, 200)
				target, targetView := CreateReferenceLink(e1View, e2View)
				epModel, epView := CreateElementPointer(sourceView, targetView)
				Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.NoAttribute))
				Expect(crldiagram.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
			Specify("for a node source and RefinementLink target", func() {
				source, sourceView := CreateReferenceNode(100, 150)
				_, e1View := CreateElement(200, 100)
				_, e2View := CreateElement(200, 200)
				target, targetView := CreateRefinementLink(e1View, e2View)
				epModel, epView := CreateElementPointer(sourceView, targetView)
				Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.NoAttribute))
				Expect(crldiagram.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
			Specify("for a node source and an OwnerPointer target", func() {
				source, sourceView := CreateReferenceNode(100, 150)
				_, e1View := CreateElement(200, 100)
				_, e2View := CreateElement(200, 200)
				target, targetView := CreateOwnerPointer(e1View, e2View)
				epModel, epView := CreateElementPointer(sourceView, targetView)
				Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.OwningConceptID))
				Expect(crldiagram.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
			Specify("for a node source and an ElementPointer target", func() {
				source, sourceView := CreateReferenceNode(100, 150)
				_, e1View := CreateReferenceNode(200, 100)
				_, e2View := CreateElement(200, 200)
				target, targetView := CreateElementPointer(e1View, e2View)
				epModel, epView := CreateElementPointer(sourceView, targetView)
				Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.ReferencedConceptID))
				Expect(crldiagram.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
			Specify("for a node source and an AbstractPointer target", func() {
				source, sourceView := CreateReferenceNode(100, 150)
				_, e1View := CreateRefinementNode(200, 100)
				_, e2View := CreateElement(200, 200)
				target, targetView := CreateAbstractPointer(e1View, e2View)
				epModel, epView := CreateElementPointer(sourceView, targetView)
				Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.AbstractConceptID))
				Expect(crldiagram.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
			Specify("for a node source and an RefinedPointer target", func() {
				source, sourceView := CreateReferenceNode(100, 150)
				_, e1View := CreateRefinementNode(200, 100)
				_, e2View := CreateElement(200, 200)
				target, targetView := CreateRefinedPointer(e1View, e2View)
				epModel, epView := CreateElementPointer(sourceView, targetView)
				Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
				Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
				Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.RefinedConceptID))
				Expect(crldiagram.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
				Expect(crldiagram.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
			})
		})
		Specify("AbstractPointer creation should work", func() {
			e1, e1View := CreateElement(100, 100)
			r1, r1View := CreateRefinementNode(100, 200)
			var toolbarID string
			Expect(page.RunScript("return crlAbstractPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
			Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			var correctToolbarSelection bool
			Expect(page.RunScript("return crlCurrentToolbarButton == crlAbstractPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
			Expect(correctToolbarSelection).To(BeTrue())
			// Now move the mouse to r1, click, drag to e1, and release
			e1CellID := GetCellViewIDFromViewElementID(e1View.GetConceptID(hl))
			r1CellID := GetCellViewIDFromViewElementID(r1View.GetConceptID(hl))
			Expect(page.FindByID(r1CellID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
			Expect(page.FindByID(e1CellID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
			hl.ReleaseLocksAndWait()
			Eventually(func() bool {
				var correctToolbarSelection bool
				page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
				return correctToolbarSelection
			}, 3).Should(BeTrue())
			hl.ReleaseLocksAndWait()
			// Now check the results
			Expect(r1.GetAbstractConcept(hl)).To(Equal(e1))
		})
		Specify("RefinedPointer creation should work", func() {
			e1, e1View := CreateElement(100, 100)
			r1, r1View := CreateRefinementNode(100, 200)
			var toolbarID string
			Expect(page.RunScript("return crlRefinedPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
			Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
			var correctToolbarSelection bool
			Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinedPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
			Expect(correctToolbarSelection).To(BeTrue())
			// Now move the mouse to r1, click, drag to e1, and release
			e1CellID := GetCellViewIDFromViewElementID(e1View.GetConceptID(hl))
			r1CellID := GetCellViewIDFromViewElementID(r1View.GetConceptID(hl))
			Expect(page.FindByID(r1CellID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
			Expect(page.FindByID(e1CellID).MouseToElement()).To(Succeed())
			Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
			hl.ReleaseLocksAndWait()
			Eventually(func() bool {
				var correctToolbarSelection bool
				page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
				return correctToolbarSelection
			}, 3).Should(BeTrue())
			hl.ReleaseLocksAndWait()
			// Now check the results
			Expect(r1.GetRefinedConcept(hl)).To(Equal(e1))
		})
	})
})
