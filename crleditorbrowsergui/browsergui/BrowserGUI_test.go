package browsergui_test

import (
	//	"fmt"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditorbrowsergui/browsergui"
	"github.com/sclevine/agouti"

	//	"testing"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Test CrlEditor", func() {

	var uOfD *core.UniverseOfDiscourse
	var hl *core.Transaction

	AssertServerRequestProcessingComplete := func() {
		EventuallyWithOffset(1, func() bool {
			// log.Printf("GetRequestInProgress: %t", browsergui.GetRequestInProgress())
			return browsergui.GetRequestInProgress() == false
		}, time.Second*10).Should(BeTrue())
	}

	SelectElementInTree := func(el core.Element) {
		var treeNodeID string
		page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
			map[string]interface{}{"conceptID": el.GetConceptID(hl)},
			&treeNodeID)
		// log.Printf("TreeNodeID: " + treeNodeID)
		// browsergui.CrlLogClientRequests = true
		// browsergui.CrlLogClientNotifications = true
		treeNode := page.FindByID(treeNodeID)
		treeNode.MouseToElement()
		page.Click(agouti.SingleClick, agouti.LeftButton)
		Eventually(func() bool {
			return testEditor.GetCurrentSelection() != nil && testEditor.GetCurrentSelection().GetConceptID(hl) == el.GetConceptID(hl)
		}).Should(BeTrue())
	}

	CreateDomain := func() (string, core.Element) {
		var oldSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID;", nil, &oldSelectionID)).To(Succeed())
		Expect(page.FindByID("FileMenuButton").Click()).To(Succeed())
		Expect(page.FindByID("NewDomainButton").Click()).To(Succeed())
		AssertServerRequestProcessingComplete()
		Eventually(func() string {
			var retrievedSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
			return retrievedSelectionID
		}, 3).ShouldNot(Equal(oldSelectionID))
		var newID string
		Expect(page.RunScript("return crlSelectedConceptID;", nil, &newID)).To(Succeed())
		Expect(newID).ToNot(Equal(""))
		Eventually(func() bool {
			conceptSpace := browsergui.BrowserGUISingleton.GetUofD().GetElement(newID)
			return conceptSpace != nil
		}, 3).Should(BeTrue())
		newCS := uOfD.GetElement(newID)
		Expect(newCS).ToNot(BeNil())
		// log.Printf("**************************** Domain creation complete")
		// At this point the newly created concept space is the selected concept
		return newID, newCS
	}

	CreateDiagram := func(parent core.Element) (string, core.Element) {
		SelectElementInTree(parent)
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID;", nil, &currentSelectionID)).To(Succeed())
		page.RunScript("crlSendAddDiagramChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": currentSelectionID}, nil)
		AssertServerRequestProcessingComplete()
		Eventually(func() string {
			var retrievedSelectionID string
			Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
			return retrievedSelectionID
		}, 3).ShouldNot(Equal(currentSelectionID))
		var newDiagramID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &newDiagramID)).To(Succeed())
		newDiagram := uOfD.GetElement(newDiagramID)
		Expect(newDiagram).ToNot(BeNil())
		// browsergui.CrlLogClientRequests = true
		Eventually(func() bool {
			return browsergui.GetRequestInProgress() == false
		}).Should(BeTrue())
		return newDiagramID, newDiagram
	}

	BeforeEach(func() {
		// Get current workspace path
		workspacePath := testWorkspaceDir
		// Open workspace (the same one - assumes nothing has been saved)
		// Because the present implementation of the editor uses a server-side dialog to select the new workspace
		// the Agouti driver can't access this. Instead, we use a convenience function on the server side to open
		// the new workspace.
		testEditor.Initialize(workspacePath, false)
		// log.Printf("Editor initialized with Workspace path: " + workspacePath)
		AssertServerRequestProcessingComplete()
		uOfD = testEditor.GetUofD()
		hl = uOfD.NewTransaction()
	})

	AfterEach(func() {
		// Clear existing workspace
		// log.Printf("**************************** About to hit ClearWorkspaceButton")
		var fileMenuButton = page.FindByID("FileMenuButton")
		Expect(fileMenuButton.Click()).To(Succeed())
		var clearWorkspaceButton = page.FindByID("ClearWorkspaceButton")
		Expect(clearWorkspaceButton.Click()).To(Succeed())
		Eventually(func() bool {
			return browsergui.GetRequestInProgress() == false
		}, time.Second*5).Should(BeTrue())
		// log.Printf("**************************** ClearWorkspace Request Complete")
		hl.ReleaseLocks()
	})

	GetCellViewIDFromViewElementID := func(diagram core.Element, viewElementID string) string {
		// First get the expected crl ID for the cell in the graph
		var crlCellID string
		Expect(page.RunScript("return crlGetJointCellIDFromConceptID(conceptID);",
			map[string]interface{}{"conceptID": viewElementID},
			&crlCellID)).To(Succeed())
		// Now find the DOM id for the cell in the graph
		var cellID string
		var diagramGraphID string
		Expect(page.RunScript("return crlGetJointGraphIDFromDiagramID(diagramID)", map[string]interface{}{"diagramID": diagram.GetConceptID(hl)}, &diagramGraphID)).To(Succeed())
		Expect(page.RunScript("return crlFindCellInGraphID(graphID, crlCellID).id",
			map[string]interface{}{"graphID": diagramGraphID, "crlCellID": crlCellID},
			&cellID)).To(Succeed())
		// Finally, find the cell view id on the paper
		var cellViewID string
		Expect(page.RunScript("return crlFindCellViewInPaperByDiagramID(diagramID, cellID).id",
			map[string]interface{}{"diagramID": diagram.GetConceptID(hl), "cellID": cellID},
			&cellViewID)).To(Succeed())
		return cellViewID
	}
	GetCurrentSelection := func() core.Element {
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		return uOfD.GetElement(currentSelectionID)
	}
	MouseToDiagramPosition := func(diagram core.Element, x int, y int) {
		var mousePosition map[string]interface{}
		Expect(page.RunScript("return crlMousePosition;", nil, &mousePosition))
		currentMouseX := int(mousePosition["x"].(float64))
		currentMouseY := int(mousePosition["y"].(float64))
		var pageOffset map[string]interface{}
		Expect(page.RunScript("var jointPaperID = crlGetJointPaperIDFromDiagramID(diagramID); return crlPapersGlobal[jointPaperID].pageOffset();",
			map[string]interface{}{"diagramID": diagram.GetConceptID(hl)}, &pageOffset)).To(Succeed())
		pageX := int(pageOffset["x"].(float64))
		pageY := int(pageOffset["y"].(float64))
		xMove := pageX + x - currentMouseX
		yMove := pageY + y - currentMouseY
		Expect(page.MoveMouseBy(xMove, yMove)).To(Succeed())
	}
	CreateElement := func(diagram core.Element, x int, y int) (core.Element, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlElementToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlElementToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move mouse to correct position
		MouseToDiagramPosition(diagram, x, y)
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		newConcept := uOfD.GetElement(currentSelectionID)
		Expect(newConcept).ToNot(BeNil())
		Expect(newConcept.GetOwningConcept(hl)).To(Equal(diagram.GetOwningConcept(hl)))
		// Check to see that the diagram view of the element has been created correctly
		conceptView := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, newConcept, hl)
		Expect(conceptView).ToNot(BeNil())
		Expect(crldiagramdomain.GetReferencedModelElement(conceptView, hl)).To(Equal(newConcept))
		return newConcept, conceptView
	}
	CreateLiteral := func(diagram core.Element, x int, y int) (core.Literal, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlLiteralToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlLiteralToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move mouse to correct position
		MouseToDiagramPosition(diagram, x, y)
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		newConcept := uOfD.GetElement(currentSelectionID)
		Expect(newConcept).ToNot(BeNil())
		Expect(newConcept.GetOwningConcept(hl)).To(Equal(diagram.GetOwningConcept(hl)))
		correctType := false
		switch newConcept.(type) {
		case core.Literal:
			correctType = true
		}
		Expect(correctType).To(BeTrue())
		// Check to see that the diagram view of the element has been created correctly
		conceptView := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, newConcept, hl)
		Expect(conceptView).ToNot(BeNil())
		Expect(crldiagramdomain.GetReferencedModelElement(conceptView, hl)).To(Equal(newConcept))
		return newConcept.(core.Literal), conceptView
	}
	CreateReferenceNode := func(diagram core.Element, x int, y int) (core.Reference, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlReferenceToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlReferenceToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move mouse to correct position
		MouseToDiagramPosition(diagram, x, y)
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
		Eventually(func() bool {
			var correctToolbarSelection bool
			Expect(page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		newConcept := uOfD.GetElement(currentSelectionID)
		Expect(newConcept).ToNot(BeNil())
		Expect(newConcept.GetOwningConcept(hl)).To(Equal(diagram.GetOwningConcept(hl)))
		correctType := false
		switch newConcept.(type) {
		case core.Reference:
			correctType = true
		}
		Expect(correctType).To(BeTrue())
		// Check to see that the diagram view of the element has been created correctly
		conceptView := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, newConcept, hl)
		Expect(conceptView).ToNot(BeNil())
		Expect(crldiagramdomain.GetReferencedModelElement(conceptView, hl)).To(Equal(newConcept))
		return newConcept.(core.Reference), conceptView
	}
	CreateReferenceLink := func(diagram core.Element, sourceView core.Element, targetView core.Element) (core.Reference, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlReferenceLinkToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlReferenceLinkToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to source, click, drag to target, and release
		targetCellID := GetCellViewIDFromViewElementID(diagram, targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(diagram, sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
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
		source := crldiagramdomain.GetReferencedModelElement(sourceView, hl)
		Expect(newReference.GetOwningConceptID(hl)).To(Equal(source.GetConceptID(hl)))
		newReferenceView := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, newReference, hl)
		return newReference, newReferenceView
	}
	CreateRefinementNode := func(diagram core.Element, x int, y int) (core.Refinement, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlRefinementToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinementToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move mouse to correct position
		MouseToDiagramPosition(diagram, x, y)
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		var currentSelectionID string
		Expect(page.RunScript("return crlSelectedConceptID", nil, &currentSelectionID)).To(Succeed())
		newConcept := uOfD.GetElement(currentSelectionID)
		Expect(newConcept).ToNot(BeNil())
		Expect(newConcept.GetOwningConcept(hl)).To(Equal(diagram.GetOwningConcept(hl)))
		correctType := false
		switch newConcept.(type) {
		case core.Refinement:
			correctType = true
		}
		Expect(correctType).To(BeTrue())
		// Check to see that the diagram view of the element has been created correctly
		conceptView := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, newConcept, hl)
		Expect(conceptView).ToNot(BeNil())
		Expect(crldiagramdomain.GetReferencedModelElement(conceptView, hl)).To(Equal(newConcept))
		Expect(newConcept.GetOwningConceptID(hl)).To(Equal(diagram.GetOwningConceptID(hl)))
		return newConcept.(core.Refinement), conceptView
	}
	CreateRefinementLink := func(diagram core.Element, sourceView core.Element, targetView core.Element) (core.Refinement, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlRefinementLinkToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinementLinkToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to e2, click, drag to e1, and release
		targetCellID := GetCellViewIDFromViewElementID(diagram, targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(diagram, sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
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
		newRefinementView := crldiagramdomain.GetFirstElementRepresentingConcept(diagram, newRefinement, hl)
		Expect(newRefinementView).ToNot(BeNil())
		Expect(newRefinement.GetOwningConceptID(hl)).To(Equal(diagram.GetOwningConceptID(hl)))
		return newRefinement, newRefinementView
	}
	CreateOwnerPointer := func(diagram core.Element, sourceView core.Element, targetView core.Element) (core.Element, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlOwnerPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlOwnerPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to source, click, drag to target, and release
		targetCellID := GetCellViewIDFromViewElementID(diagram, targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(diagram, sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		source := crldiagramdomain.GetReferencedModelElement(sourceView, hl)
		ownerPointerView := crldiagramdomain.GetFirstElementRepresentingConceptOwnerPointer(diagram, source, hl)
		return source, ownerPointerView
	}
	CreateElementPointer := func(diagram core.Element, sourceView core.Element, targetView core.Element) (core.Reference, core.Element) {
		var toolbarID string
		// core.TraceLocks = true
		Expect(page.RunScript("return crlElementPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlElementPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to r1, click, drag to e1, and release
		// browsergui.CrlLogClientRequests = true
		// browsergui.CrlLogClientNotifications = true
		targetCellID := GetCellViewIDFromViewElementID(diagram, targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(diagram, sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		// core.TraceChange = true
		// core.EnableNotificationPrint = true
		// AssertServerRequestProcessingComplete()
		time.Sleep(time.Second * 3)
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		// core.TraceChange = false
		// core.EnableNotificationPrint = false
		// browsergui.CrlLogClientRequests = false
		// browsergui.CrlLogClientNotifications = false
		referenceID := crldiagramdomain.GetReferencedModelElement(sourceView, hl).GetConceptID(hl)
		reference := uOfD.GetReference(referenceID)
		elementPointerView := crldiagramdomain.GetFirstElementRepresentingConceptElementPointer(diagram, reference, hl)
		// core.TraceLocks = false
		return reference, elementPointerView
	}

	CreateAbstractPointer := func(diagram core.Element, sourceView core.Element, targetView core.Element) (core.Refinement, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlAbstractPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlAbstractPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to r1, click, drag to e1, and release
		targetCellID := GetCellViewIDFromViewElementID(diagram, targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(diagram, sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		refinementID := crldiagramdomain.GetReferencedModelElement(sourceView, hl).GetConceptID(hl)
		refinement := uOfD.GetRefinement(refinementID)
		elementPointerView := crldiagramdomain.GetFirstElementRepresentingConceptAbstractPointer(diagram, refinement, hl)
		return refinement, elementPointerView
	}

	CreateRefinedPointer := func(diagram core.Element, sourceView core.Element, targetView core.Element) (core.Refinement, core.Element) {
		var toolbarID string
		Expect(page.RunScript("return crlRefinedPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
		Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		var correctToolbarSelection bool
		Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinedPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
		Expect(correctToolbarSelection).To(BeTrue())
		// Now move the mouse to r1, click, drag to e1, and release
		targetCellID := GetCellViewIDFromViewElementID(diagram, targetView.GetConceptID(hl))
		sourceCellID := GetCellViewIDFromViewElementID(diagram, sourceView.GetConceptID(hl))
		Expect(page.FindByID(sourceCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID(targetCellID).MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
		Eventually(func() bool {
			var correctToolbarSelection bool
			page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
			return correctToolbarSelection
		}, 3).Should(BeTrue())
		refinementID := crldiagramdomain.GetReferencedModelElement(sourceView, hl).GetConceptID(hl)
		refinement := uOfD.GetRefinement(refinementID)
		elementPointerView := crldiagramdomain.GetFirstElementRepresentingConceptRefinedPointer(diagram, refinement, hl)
		return refinement, elementPointerView
	}

	Undo := func() {
		Expect(page.FindByID("EditMenuButton").MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID("UndoButton").MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
	}

	Redo := func() {
		Expect(page.FindByID("EditMenuButton").MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		Expect(page.FindByID("RedoButton").MouseToElement()).To(Succeed())
		Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
		AssertServerRequestProcessingComplete()
	}

	Describe("Testing CrlEditor basic functionality", func() {
		Specify("The editor should be initialized", func() {
			Expect(browsergui.BrowserGUISingleton.IsInitialized()).To(BeTrue())
			var initializationComplete interface{}
			page.RunScript("return crlInitializationComplete;", nil, &initializationComplete)
			Expect(initializationComplete).To(BeTrue())
			coreDomain := uOfD.GetElementWithURI(core.CoreDomainURI)
			var treeNodeID string
			page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
				map[string]interface{}{"conceptID": coreDomain.GetConceptID(hl)},
				&treeNodeID)
			Expect(page.FindByID(treeNodeID)).To(BeFound())
		})
		Specify("Tree selection should work", func() {
			coreDomain := uOfD.GetElementWithURI(core.CoreDomainURI)
			Expect(coreDomain).ToNot(BeNil())
			var treeNodeID string
			page.RunScript("return crlGetTreeNodeIDFromConceptID(conceptID);",
				map[string]interface{}{"conceptID": coreDomain.GetConceptID(hl)},
				&treeNodeID)
			// log.Printf("TreeNodeID: " + treeNodeID)
			// browsergui.CrlLogClientRequests = true
			// browsergui.CrlLogClientNotifications = true
			treeNode := page.FindByID(treeNodeID)
			treeNode.MouseToElement()
			page.Click(agouti.SingleClick, agouti.LeftButton)
			AssertServerRequestProcessingComplete()
			Eventually(func() bool {
				return testEditor.GetCurrentSelection() != nil && testEditor.GetCurrentSelection().GetConceptID(hl) == coreDomain.GetConceptID(hl)
			}).Should(BeTrue())
			// browsergui.CrlLogClientRequests = false
			// browsergui.CrlLogClientNotifications = false
		})
		Specify("UndoRedo of a concept space should work", func() {
			uOfD.MarkUndoPoint()
			beforeUofD := uOfD.Clone(hl)
			beforeHL := beforeUofD.NewTransaction()
			_, cs1 := CreateDomain()
			Expect(cs1).ToNot(BeNil())
			afterUofD := uOfD.Clone(hl)
			afterHL := afterUofD.NewTransaction()
			Undo()
			Expect(uOfD.IsEquivalent(hl, beforeUofD, beforeHL, true)).To(BeTrue())
			Redo()
			Expect(uOfD.IsEquivalent(hl, afterUofD, afterHL, true)).To(BeTrue())
		})

		Specify("UndoRedo of a diagram creation should work", func() {
			_, cs1 := CreateDomain()
			Expect(cs1).ToNot(BeNil())
			uOfD.MarkUndoPoint()
			beforeUofD := uOfD.Clone(hl)
			beforeHL := beforeUofD.NewTransaction()
			_, diag := CreateDiagram(cs1)
			Expect(diag).ToNot(BeNil())
			afterUofD := uOfD.Clone(hl)
			afterHL := afterUofD.NewTransaction()
			Undo()
			Expect(uOfD.IsEquivalent(hl, beforeUofD, beforeHL, true)).To(BeTrue())
			Redo()
			Expect(uOfD.IsEquivalent(hl, afterUofD, afterHL, true)).To(BeTrue())
		})

		Describe("Single Diagram Tests", func() {
			var cs1ID string
			var cs1 core.Element
			var diagramID string
			var diagram core.Element
			var beforeUofD *core.UniverseOfDiscourse
			var beforeHL *core.Transaction
			var afterUofD *core.UniverseOfDiscourse
			var afterHL *core.Transaction

			BeforeEach(func() {
				cs1ID, cs1 = CreateDomain()

				// Now add a diagram
				diagramID, diagram = CreateDiagram(cs1)
				// Expect(page.RunScript("return crlGetContainerIDFromConceptID(conceptID)", map[string]interface{}{"conceptID": diagramID}, &diagramContainerID)).To(Succeed())
				// Expect(page.RunScript("return crlGetJointGraphIDFromDiagramID(diagramID)", map[string]interface{}{"diagramID": diagramID}, &diagramGraphID)).To(Succeed())
				uOfD.MarkUndoPoint()
				beforeUofD = uOfD.Clone(hl)
				beforeHL = beforeUofD.NewTransaction()
			})

			PerformUndoRedoTest := func(count int) {
				afterUofD = uOfD.Clone(hl)
				afterHL = afterUofD.NewTransaction()
				for i := 0; i < count; i++ {
					Undo()
				}
				Expect(uOfD.IsEquivalent(hl, beforeUofD, beforeHL, true)).To(BeTrue())
				for i := 0; i < count; i++ {
					Redo()
				}
				Expect(uOfD.IsEquivalent(hl, afterUofD, afterHL, true)).To(BeTrue())
			}
			Specify("DiagramDrop should produce view of treeDragSelection", func() {
				coreDomain := uOfD.GetElementWithURI(core.CoreDomainURI)
				Expect(page.RunScript("crlSendSetTreeDragSelection(ID)", map[string]interface{}{"ID": coreDomain.GetConceptID(hl)}, nil)).To(Succeed())
				AssertServerRequestProcessingComplete()
				Expect(page.RunScript("crlSendDiagramDrop(ID, x, y, shiftKey)", map[string]interface{}{"ID": diagramID, "x": "100", "y": "100", "shiftKey": "false"}, nil)).To(Succeed())
				// Some form of sleep is required here as this thread blocks socket communications. Eventually accomplishes this as it will not
				// be true until after all of the expected client communication has completed.
				AssertServerRequestProcessingComplete()
				Eventually(func() bool {
					return browsergui.BrowserGUISingleton.GetTreeDragSelection() == nil
				}, 3).Should(BeTrue())
				Expect(len(diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, hl))).To(Equal(1))
				newNode := diagram.GetFirstOwnedConceptRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, hl)
				Expect(newNode).ToNot(BeNil())
				Expect(newNode.GetLabel(hl)).To(Equal(coreDomain.GetLabel(hl)))
				Expect(crldiagramdomain.GetDisplayLabel(newNode, hl)).To(Equal(coreDomain.GetLabel(hl)))
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
				Expect(page.RunScript("crlSendSetTreeDragSelection(ID)", map[string]interface{}{"ID": coreDomain.GetConceptID(hl)}, nil)).To(Succeed())
				AssertServerRequestProcessingComplete()
				Expect(page.RunScript("crlSendDiagramDrop(ID, x, y, shiftKey)", map[string]interface{}{"ID": diagramID, "x": "200", "y": "200", "shiftKey": "false"}, nil)).To(Succeed())
				// Some form of sleep is required here as this thread blocks socket communications. Eventually accomplishes this as it will not
				// be true until after all of the expected client communication has completed.
				AssertServerRequestProcessingComplete()
				Eventually(func() bool {
					return browsergui.BrowserGUISingleton.GetTreeDragSelection() == nil
				}, 3).Should(BeTrue())
				Expect(len(diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, hl))).To(Equal(2))
				var newNode2 core.Element
				for _, el := range diagram.GetOwnedConceptsRefinedFromURI(crldiagramdomain.CrlDiagramNodeURI, hl) {
					if el != newNode {
						newNode2 = el
					}
				}
				Expect(newNode2).ToNot(BeNil())
				Expect(newNode2.GetLabel(hl)).To(Equal(coreDomain.GetLabel(hl)))
				Expect(crldiagramdomain.GetDisplayLabel(newNode2, hl)).To(Equal(coreDomain.GetLabel(hl)))
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
				PerformUndoRedoTest(4)
			})
			Describe("Test AddChild functionality", func() {
				Specify("AddChild Diagram should work", func() {
					var initialSelectionID string
					Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
					Expect(page.RunScript("crlSendAddDiagramChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
					AssertServerRequestProcessingComplete()
					Eventually(func() string {
						var retrievedSelectionID string
						Expect(page.RunScript("return crlSelectedConceptID;", nil, &retrievedSelectionID)).To(Succeed())
						return retrievedSelectionID
					}, 3).ShouldNot(Equal(initialSelectionID))
					var newDiagramID string
					Expect(page.RunScript("return crlSelectedConceptID", nil, &newDiagramID)).To(Succeed())
					newDiagram := uOfD.GetElement(newDiagramID)
					Expect(newDiagram).ToNot(BeNil())
					Expect(newDiagram.IsRefinementOfURI(crldiagramdomain.CrlDiagramURI, hl)).To(BeTrue())
					Expect(newDiagram.GetOwningConcept(hl)).To(Equal(cs1))
					PerformUndoRedoTest(1)
				})
				Specify("AddChild Element should work", func() {
					var initialSelectionID string
					Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
					Expect(page.RunScript("crlSendAddElementChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
					AssertServerRequestProcessingComplete()
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
					PerformUndoRedoTest(1)
				})
				Specify("AddChild Literal should work", func() {
					var initialSelectionID string
					Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
					Expect(page.RunScript("crlSendAddLiteralChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
					AssertServerRequestProcessingComplete()
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
					PerformUndoRedoTest(1)
				})
				Specify("AddChild Reference should work", func() {
					var initialSelectionID string
					Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
					Expect(page.RunScript("crlSendAddReferenceChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
					AssertServerRequestProcessingComplete()
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
					PerformUndoRedoTest(1)
				})
				Specify("AddChild Refinement should work", func() {
					var initialSelectionID string
					Expect(page.RunScript("return crlSelectedConceptID", nil, &initialSelectionID)).To(Succeed())
					Expect(page.RunScript("crlSendAddRefinementChild(conceptSpaceID);", map[string]interface{}{"conceptSpaceID": cs1ID}, nil)).To(Succeed())
					AssertServerRequestProcessingComplete()
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
					PerformUndoRedoTest(1)
				})
			})

			Describe("Test Toolbar Functionality", func() {
				Specify("Element node creation should work", func() {
					e1, e1View := CreateElement(diagram, 100, 100)
					Expect(e1).ToNot(BeNil())
					Expect(e1View).ToNot(BeNil())
					PerformUndoRedoTest(1)
				})
				Specify("Literal node creation should work", func() {
					l1, l1View := CreateLiteral(diagram, 100, 100)
					Expect(l1).ToNot(BeNil())
					Expect(l1View).ToNot(BeNil())
					PerformUndoRedoTest(1)
				})
				Specify("Reference node creation should work", func() {
					r1, r1View := CreateReferenceNode(diagram, 100, 100)
					Expect(r1).ToNot(BeNil())
					Expect(r1View).ToNot(BeNil())
					PerformUndoRedoTest(1)
				})
				Specify("Refinement node creation should work", func() {
					// browsergui.CrlLogClientDialog = true
					r1, r1View := CreateRefinementNode(diagram, 100, 100)
					Expect(r1).ToNot(BeNil())
					Expect(r1View).ToNot(BeNil())
					PerformUndoRedoTest(1)
				})
				Describe("Reference link creation should work", func() {
					Specify("for a node source and target", func() {
						e1, e1View := CreateElement(diagram, 100, 100)
						e2, e2View := CreateElement(diagram, 100, 200)
						newRefinement, _ := CreateReferenceLink(diagram, e2View, e1View)
						Expect(newRefinement.GetOwningConcept(hl)).To(Equal(e2))
						Expect(newRefinement.GetReferencedConcept(hl)).To(Equal(e1))
						PerformUndoRedoTest(3)
					})
					Specify("for a link source and node target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						_, e2View := CreateElement(diagram, 100, 200)
						// create the node target
						e3, e3View := CreateElement(diagram, 200, 150)
						// Create a reference link
						refLink1, refLink1View := CreateReferenceLink(diagram, e2View, e1View)
						// Now the new reference
						refLink2, _ := CreateReferenceLink(diagram, refLink1View, e3View)
						// Now check the results
						Expect(refLink2.GetOwningConceptID(hl)).To(Equal(refLink1.GetConceptID(hl)))
						Expect(refLink2.GetReferencedConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and link target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						_, e2View := CreateElement(diagram, 100, 200)
						// create the node source
						e3, e3View := CreateElement(diagram, 200, 150)
						// Create a reference link
						refLink1, refLink1View := CreateReferenceLink(diagram, e2View, e1View)
						// Now the new reference
						refLink2, _ := CreateReferenceLink(diagram, e3View, refLink1View)
						// Now check the results
						Expect(refLink2.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
						Expect(refLink2.GetReferencedConceptID(hl)).To(Equal(refLink1.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a link source and link target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						_, e2View := CreateElement(diagram, 100, 200)
						_, e3View := CreateElement(diagram, 200, 100)
						_, e4View := CreateElement(diagram, 200, 200)
						// Create the source reference link
						refLink1, refLink1View := CreateReferenceLink(diagram, e2View, e1View)
						// Create the target reference link
						refLink2, refLink2View := CreateReferenceLink(diagram, e4View, e3View)
						// Now the new reference
						refLink3, _ := CreateReferenceLink(diagram, refLink1View, refLink2View)
						// Now check the results
						Expect(refLink3.GetOwningConceptID(hl)).To(Equal(refLink1.GetConceptID(hl)))
						Expect(refLink3.GetReferencedConceptID(hl)).To(Equal(refLink2.GetConceptID(hl)))
						PerformUndoRedoTest(7)
					})
					Specify("for a node source and an OwnerPointer target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						e2, e2View := CreateElement(diagram, 100, 200)
						// create the node source
						e3, e3View := CreateElement(diagram, 200, 150)
						// create the owner pointer
						opModelElement, opView := CreateOwnerPointer(diagram, e2View, e1View)
						// Create the Reference
						ref, refView := CreateReferenceLink(diagram, e3View, opView)
						Expect(opModelElement.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(ref.GetReferencedConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(ref.GetReferencedAttributeName(hl)).To(Equal(core.OwningConceptID))
						Expect(ref.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(refView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(refView, hl).GetConceptID(hl)).To(Equal(opView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and an ElementPointer target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						e2, e2View := CreateReferenceNode(diagram, 100, 200)
						// create the node source
						e3, e3View := CreateElement(diagram, 200, 150)
						// create the owner pointer
						epModelElement, epView := CreateElementPointer(diagram, e2View, e1View)
						// Create the Reference
						ref, refView := CreateReferenceLink(diagram, e3View, epView)
						Expect(epModelElement.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(ref.GetReferencedConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(ref.GetReferencedAttributeName(hl)).To(Equal(core.ReferencedConceptID))
						Expect(ref.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(refView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(refView, hl).GetConceptID(hl)).To(Equal(epView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and an AbstractPointer target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						e2, e2View := CreateRefinementNode(diagram, 100, 200)
						// create the node source
						e3, e3View := CreateElement(diagram, 200, 150)
						// create the owner pointer
						apModelElement, apView := CreateAbstractPointer(diagram, e2View, e1View)
						// Create the Reference
						ref, refView := CreateReferenceLink(diagram, e3View, apView)
						Expect(apModelElement.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(ref.GetReferencedConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(ref.GetReferencedAttributeName(hl)).To(Equal(core.AbstractConceptID))
						Expect(ref.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(refView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(refView, hl).GetConceptID(hl)).To(Equal(apView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and an RefinedPointer target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						e2, e2View := CreateRefinementNode(diagram, 100, 200)
						// create the node source
						e3, e3View := CreateElement(diagram, 200, 150)
						// create the owner pointer
						apModelElement, apView := CreateRefinedPointer(diagram, e2View, e1View)
						// Create the Reference
						ref, refView := CreateReferenceLink(diagram, e3View, apView)
						Expect(apModelElement.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(ref.GetReferencedConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(ref.GetReferencedAttributeName(hl)).To(Equal(core.RefinedConceptID))
						Expect(ref.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(refView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(refView, hl).GetConceptID(hl)).To(Equal(apView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
				})
				Describe("Refinement link creation should work", func() {
					Specify("for a node source and node target", func() {
						e1, e1View := CreateElement(diagram, 100, 100)
						e2, e2View := CreateElement(diagram, 100, 200)
						newRefinement, _ := CreateRefinementLink(diagram, e2View, e1View)
						Expect(newRefinement.GetAbstractConcept(hl)).To(Equal(e1))
						Expect(newRefinement.GetRefinedConcept(hl)).To(Equal(e2))
						PerformUndoRedoTest(3)
					})
					Specify("for a link source and node target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						_, e2View := CreateElement(diagram, 100, 200)
						source, sourceView := CreateRefinementLink(diagram, e2View, e1View)
						target, targetView := CreateRefinementNode(diagram, 200, 150)
						newRefinement, newRefinementView := CreateRefinementLink(diagram, sourceView, targetView)
						Expect(newRefinement.GetAbstractConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(newRefinement.GetRefinedConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(newRefinementView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(newRefinementView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and link target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						_, e2View := CreateElement(diagram, 100, 200)
						target, targetView := CreateRefinementLink(diagram, e2View, e1View)
						source, sourceView := CreateRefinementNode(diagram, 200, 150)
						newRefinement, newRefinementView := CreateRefinementLink(diagram, sourceView, targetView)
						Expect(newRefinement.GetAbstractConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(newRefinement.GetRefinedConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(newRefinementView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(newRefinementView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
				})
				Describe("OwnerPointer creation should work", func() {
					Specify("For a node source and node target", func() {
						e1, e1View := CreateElement(diagram, 100, 100)
						e2, e2View := CreateElement(diagram, 100, 200)
						source, ownerPointerView := CreateOwnerPointer(diagram, e2View, e1View)
						// Now check the results
						Expect(source.GetConceptID(hl)).To(Equal(e2.GetConceptID(hl)))
						Expect(e2.GetOwningConceptID(hl)).To(Equal(e1.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(ownerPointerView, hl).GetConceptID(hl)).To(Equal(e2View.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(ownerPointerView, hl).GetConceptID(hl)).To(Equal(e1View.GetConceptID(hl)))
						PerformUndoRedoTest(3)
					})
					Specify("For a Refinement Link source and node target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						_, e2View := CreateElement(diagram, 100, 200)
						// Create the refinement link
						source, sourceView := CreateRefinementLink(diagram, e1View, e2View)
						// create the new owner
						e3, e3View := CreateElement(diagram, 200, 150)
						// Now the ownerPointer
						ownerPointerConcept, ownerPointerView := CreateOwnerPointer(diagram, sourceView, e3View)
						// Now check the results
						Expect(source.GetOwningConceptID(hl)).To(Equal(e3.GetConceptID(hl)))
						Expect(source.GetConceptID(hl)).To(Equal(ownerPointerConcept.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(ownerPointerView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(ownerPointerView, hl).GetConceptID(hl)).To(Equal(e3View.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("For a node source and ReferenceLink target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						_, e2View := CreateElement(diagram, 100, 200)
						// Create the Reference link
						target, targetView := CreateReferenceLink(diagram, e1View, e2View)
						// create the new owner
						source, sourceView := CreateElement(diagram, 200, 150)
						// Now the ownerPointer
						ownerPointerConcept, ownerPointerView := CreateOwnerPointer(diagram, sourceView, targetView)
						// Now check the results
						Expect(source.GetOwningConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetConceptID(hl)).To(Equal(ownerPointerConcept.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(ownerPointerView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(ownerPointerView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("For a node source and RefinementLink target", func() {
						_, e1View := CreateElement(diagram, 100, 100)
						_, e2View := CreateElement(diagram, 100, 200)
						// Create the refinement link
						target, targetView := CreateRefinementLink(diagram, e1View, e2View)
						// create the new owner
						source, sourceView := CreateElement(diagram, 200, 150)
						// Now the ownerPointer
						ownerPointerConcept, ownerPointerView := CreateOwnerPointer(diagram, sourceView, targetView)
						// Now check the results
						Expect(source.GetOwningConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetConceptID(hl)).To(Equal(ownerPointerConcept.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkSource(ownerPointerView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(ownerPointerView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
				})
				Describe("ElementPointer creation should work", func() {
					Specify("for a node source and node target", func() {
						target, targetView := CreateElement(diagram, 100, 100)
						source, sourceView := CreateReferenceNode(diagram, 100, 200)
						reference, epView := CreateElementPointer(diagram, sourceView, targetView)
						// Now check the results
						Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(reference.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(reference.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.NoAttribute))
						Expect(crldiagramdomain.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(3)
					})
					Specify("for a node source and reference link target", func() {
						source, sourceView := CreateReferenceNode(diagram, 100, 150)
						_, e1View := CreateElement(diagram, 200, 100)
						_, e2View := CreateElement(diagram, 200, 200)
						target, targetView := CreateReferenceLink(diagram, e1View, e2View)
						epModel, epView := CreateElementPointer(diagram, sourceView, targetView)
						Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.NoAttribute))
						Expect(crldiagramdomain.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and RefinementLink target", func() {
						source, sourceView := CreateReferenceNode(diagram, 100, 150)
						_, e1View := CreateElement(diagram, 200, 100)
						_, e2View := CreateElement(diagram, 200, 200)
						target, targetView := CreateRefinementLink(diagram, e1View, e2View)
						epModel, epView := CreateElementPointer(diagram, sourceView, targetView)
						Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.NoAttribute))
						Expect(crldiagramdomain.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and an OwnerPointer target", func() {
						source, sourceView := CreateReferenceNode(diagram, 100, 150)
						e1, e1View := CreateElement(diagram, 200, 100)
						e2, e2View := CreateElement(diagram, 200, 200)
						Expect(e1).ToNot(BeNil())
						Expect(e2).ToNot(BeNil())
						target, targetView := CreateOwnerPointer(diagram, e1View, e2View)
						epModel, epView := CreateElementPointer(diagram, sourceView, targetView)
						Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.OwningConceptID))
						Expect(crldiagramdomain.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and an ElementPointer target", func() {
						source, sourceView := CreateReferenceNode(diagram, 100, 150)
						_, e1View := CreateReferenceNode(diagram, 200, 100)
						_, e2View := CreateElement(diagram, 200, 200)
						target, targetView := CreateElementPointer(diagram, e1View, e2View)
						epModel, epView := CreateElementPointer(diagram, sourceView, targetView)
						Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.ReferencedConceptID))
						Expect(crldiagramdomain.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and an AbstractPointer target", func() {
						source, sourceView := CreateReferenceNode(diagram, 100, 150)
						_, e1View := CreateRefinementNode(diagram, 200, 100)
						_, e2View := CreateElement(diagram, 200, 200)
						target, targetView := CreateAbstractPointer(diagram, e1View, e2View)
						epModel, epView := CreateElementPointer(diagram, sourceView, targetView)
						Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.AbstractConceptID))
						Expect(crldiagramdomain.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
					Specify("for a node source and an RefinedPointer target", func() {
						source, sourceView := CreateReferenceNode(diagram, 100, 150)
						_, e1View := CreateRefinementNode(diagram, 200, 100)
						_, e2View := CreateElement(diagram, 200, 200)
						target, targetView := CreateRefinedPointer(diagram, e1View, e2View)
						epModel, epView := CreateElementPointer(diagram, sourceView, targetView)
						Expect(epModel.GetConceptID(hl)).To(Equal(source.GetConceptID(hl)))
						Expect(source.GetReferencedConceptID(hl)).To(Equal(target.GetConceptID(hl)))
						Expect(source.GetReferencedAttributeName(hl)).To(Equal(core.RefinedConceptID))
						Expect(crldiagramdomain.GetLinkSource(epView, hl).GetConceptID(hl)).To(Equal(sourceView.GetConceptID(hl)))
						Expect(crldiagramdomain.GetLinkTarget(epView, hl).GetConceptID(hl)).To(Equal(targetView.GetConceptID(hl)))
						PerformUndoRedoTest(5)
					})
				})
				Specify("AbstractPointer creation should work", func() {
					e1, e1View := CreateElement(diagram, 100, 100)
					r1, r1View := CreateRefinementNode(diagram, 100, 200)
					var toolbarID string
					Expect(page.RunScript("return crlAbstractPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
					Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
					Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
					var correctToolbarSelection bool
					Expect(page.RunScript("return crlCurrentToolbarButton == crlAbstractPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
					Expect(correctToolbarSelection).To(BeTrue())
					// Now move the mouse to r1, click, drag to e1, and release
					e1CellID := GetCellViewIDFromViewElementID(diagram, e1View.GetConceptID(hl))
					r1CellID := GetCellViewIDFromViewElementID(diagram, r1View.GetConceptID(hl))
					Expect(page.FindByID(r1CellID).MouseToElement()).To(Succeed())
					Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
					Expect(page.FindByID(e1CellID).MouseToElement()).To(Succeed())
					Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
					AssertServerRequestProcessingComplete()
					Eventually(func() bool {
						var correctToolbarSelection bool
						page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
						return correctToolbarSelection
					}, 3).Should(BeTrue())
					// Now check the results
					Expect(r1.GetAbstractConcept(hl)).To(Equal(e1))
					PerformUndoRedoTest(3)
				})
				Specify("RefinedPointer creation should work", func() {
					e1, e1View := CreateElement(diagram, 100, 100)
					r1, r1View := CreateRefinementNode(diagram, 100, 200)
					var toolbarID string
					Expect(page.RunScript("return crlRefinedPointerToolbarButtonID", nil, &toolbarID)).To(Succeed())
					Expect(page.FindByID(toolbarID).MouseToElement()).To(Succeed())
					Expect(page.Click(agouti.SingleClick, agouti.LeftButton)).To(Succeed())
					var correctToolbarSelection bool
					Expect(page.RunScript("return crlCurrentToolbarButton == crlRefinedPointerToolbarButtonID;", nil, &correctToolbarSelection)).To(Succeed())
					Expect(correctToolbarSelection).To(BeTrue())
					// Now move the mouse to r1, click, drag to e1, and release
					e1CellID := GetCellViewIDFromViewElementID(diagram, e1View.GetConceptID(hl))
					r1CellID := GetCellViewIDFromViewElementID(diagram, r1View.GetConceptID(hl))
					Expect(page.FindByID(r1CellID).MouseToElement()).To(Succeed())
					Expect(page.Click(agouti.HoldClick, agouti.LeftButton)).To(Succeed())
					Expect(page.FindByID(e1CellID).MouseToElement()).To(Succeed())
					Expect(page.Click(agouti.ReleaseClick, agouti.LeftButton)).To(Succeed())
					AssertServerRequestProcessingComplete()
					Eventually(func() bool {
						var correctToolbarSelection bool
						page.RunScript("return crlCurrentToolbarButton == crlCursorToolbarButtonID;", nil, &correctToolbarSelection)
						return correctToolbarSelection
					}, 3).Should(BeTrue())
					// Now check the results
					Expect(r1.GetRefinedConcept(hl)).To(Equal(e1))
					PerformUndoRedoTest(3)
				})
			})

		})
	})

})
