package crleditorfynegui

import (
	//	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/x/fyne/widget/diagramwidget"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"

	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldiagramdomain"
	"github.com/pbrown12303/activeCRL/crleditor"
)

var _ = Describe("Basic CRLEditorFyneGUI testing", func() {
	var uOfD *core.UniverseOfDiscourse
	var trans *core.Transaction

	BeforeEach(func() {
		uOfD, trans = beforeEachTest()
	})

	AfterEach(func() {
		afterEachTest()
	})

	Describe("Testing CrlEditor basic functionality", func() {
		Specify("The FyneGUISingleton should be populated", func() {
			Expect(trans).ToNot(BeNil())
			Expect(FyneGUISingleton).ToNot(BeNil())
			coreDomain := uOfD.GetElementWithURI(core.CoreDomainURI)
			Expect(coreDomain).ToNot(BeNil())
			Expect(test.AssertRendersToImage(testT, "initialScreen.png", FyneGUISingleton.window.Canvas())).To(BeTrue())
		})
		Specify("Tree selection and crleditor selection should match", func() {
			coreDomain := uOfD.GetElementWithURI(core.CoreDomainURI)
			coreDomainID := coreDomain.GetConceptID(trans)
			Expect(FyneGUISingleton.editor.SelectElement(coreDomain, trans)).To(Succeed())
			Expect(test.AssertRendersToImage(testT, "/selectionTests/coreDomainSelectedViaEditor.png", FyneGUISingleton.window.Canvas())).To(BeTrue())
			Expect(FyneGUISingleton.editor.GetCurrentSelectionID(trans)).To(Equal(coreDomainID))
			Expect(FyneGUISingleton.editor.SelectElementUsingIDString("", trans)).To(Succeed())
			Expect(FyneGUISingleton.editor.GetCurrentSelectionID(trans)).To(Equal(""))
			Expect(test.AssertRendersToImage(testT, "/selectionTests/noSelection.png", FyneGUISingleton.window.Canvas())).To(BeTrue())
			FyneGUISingleton.treeManager.tree.Select(coreDomainID)
			Expect(test.AssertRendersToImage(testT, "/selectionTests/coreDomainSelectedViaTree.png", FyneGUISingleton.window.Canvas())).To(BeTrue())
			Expect(FyneGUISingleton.editor.GetCurrentSelectionID(trans)).To(Equal(coreDomainID))
		})
		Specify("Domain creation should Undo and Redo successfully", func() {
			uOfD.MarkUndoPoint()
			beforeUofD := uOfD.Clone(trans)
			beforeTrans := beforeUofD.NewTransaction()
			cs1 := FyneGUISingleton.addElement("", FyneGUISingleton.editor.GetDefaultDomainLabel())
			Expect(cs1).ToNot(BeNil())
			Expect(crleditor.CrlEditorSingleton.GetCurrentSelection()).To(Equal(cs1))
			afterUofD := uOfD.Clone(trans)
			afterTrans := afterUofD.NewTransaction()
			FyneGUISingleton.undo()
			Expect(uOfD.IsEquivalent(trans, beforeUofD, beforeTrans, true)).To(BeTrue())
			Expect(crleditor.CrlEditorSingleton.GetCurrentSelection()).To(BeNil())
			FyneGUISingleton.redo()
			Expect(uOfD.IsEquivalent(trans, afterUofD, afterTrans, true)).To(BeTrue())
			Expect(crleditor.CrlEditorSingleton.GetCurrentSelection()).To(Equal(cs1))
		})
		Specify("UndoRedo of a diagram creation should work", func() {
			cs1 := FyneGUISingleton.addElement("", FyneGUISingleton.editor.GetDefaultDomainLabel())
			uOfD.MarkUndoPoint()
			Expect(cs1).ToNot(BeNil())
			beforeUofD := uOfD.Clone(trans)
			beforeTrans := beforeUofD.NewTransaction()
			diag := FyneGUISingleton.addDiagram(cs1.GetConceptID(trans))
			Expect(diag).ToNot(BeNil())
			afterUofD := uOfD.Clone(trans)
			afterTrans := afterUofD.NewTransaction()
			FyneGUISingleton.undo()
			Expect(uOfD.IsEquivalent(trans, beforeUofD, beforeTrans, true)).To(BeTrue())
			FyneGUISingleton.redo()
			Expect(uOfD.IsEquivalent(trans, afterUofD, afterTrans, true)).To(BeTrue())
		})
	})

	Describe("Single Diagram Tests", func() {
		var cs1ID string
		var cs1 core.Element
		var diagramID string
		var diagram core.Element
		var beforeUofD *core.UniverseOfDiscourse
		var beforeTrans *core.Transaction
		var afterUofD *core.UniverseOfDiscourse
		var afterTrans *core.Transaction

		BeforeEach(func() {
			cs1 = FyneGUISingleton.addElement("", FyneGUISingleton.editor.GetDefaultDomainLabel())
			cs1ID = cs1.GetConceptID(trans)
			Expect(cs1ID).ToNot(Equal(""))
			diagram = FyneGUISingleton.addDiagram(cs1.GetConceptID(trans))
			diagramID = diagram.GetConceptID(trans)
			Expect(diagramID).ToNot(Equal(""))
			uOfD.MarkUndoPoint()
			beforeUofD = uOfD.Clone(trans)
			beforeTrans = beforeUofD.NewTransaction()
		})

		PerformUndoRedoTest := func(count int) {
			afterUofD = uOfD.Clone(trans)
			afterTrans = afterUofD.NewTransaction()
			for i := 0; i < count; i++ {
				FyneGUISingleton.undo()
			}
			Expect(uOfD.IsEquivalent(trans, beforeUofD, beforeTrans, true)).To(BeTrue())
			for i := 0; i < count; i++ {
				FyneGUISingleton.redo()
			}
			Expect(uOfD.IsEquivalent(trans, afterUofD, afterTrans, true)).To(BeTrue())
		}
		Specify("Drag and drop of a tree node should produce a view of the element represented by the tree node", func() {
			fyneDiagram := FyneGUISingleton.diagramManager.getDiagramWidget(diagramID)
			coreDomain := uOfD.GetElementWithURI(core.CoreDomainURI)
			coreDomainID := coreDomain.GetConceptID(trans)
			FyneGUISingleton.editor.SelectElementUsingIDString(coreDomainID, trans)
			treeNode := FyneGUISingleton.treeManager.treeNodes[coreDomainID]
			treeNodePosition := FyneGUISingleton.app.Driver().AbsolutePositionForObject(treeNode)
			diagramPosition := FyneGUISingleton.app.Driver().AbsolutePositionForObject(fyneDiagram)
			Expect(test.AssertRendersToImage(testT, "/dragDropTests/PriorToFirstNodeDrop.png", FyneGUISingleton.window.Canvas())).To(BeTrue())
			treeNode.Dragged(newDragEvent(diagramPosition.X-treeNodePosition.X+100, diagramPosition.Y-treeNodePosition.Y+100))
			Expect(FyneGUISingleton.dragDropTransaction).ToNot(BeNil())
			Expect(FyneGUISingleton.dragDropTransaction.id).ToNot(Equal(""))
			fyneDiagram.MouseMoved(newLeftMouseEventAt(fyne.NewPos(100, 100), diagramPosition.AddXY(100, 100)))
			Expect(FyneGUISingleton.dragDropTransaction.currentDiagramMousePosition).To(Equal(fyne.NewPos(100, 100)))
			treeNode.DragEnd()
			Expect(test.AssertRendersToImage(testT, "/dragDropTests/FirstNodeDrop.png", FyneGUISingleton.window.Canvas())).To(BeTrue())
			fyneNode := fyneDiagram.GetPrimarySelection()
			Expect(fyneNode).ToNot(BeNil())
			newNodeID := fyneNode.GetDiagramElementID()
			crlNode := uOfD.GetElement(newNodeID)
			Expect(crlNode).ToNot(BeNil())
			crlModelElement := crldiagramdomain.GetReferencedModelElement(crlNode, trans)
			Expect(crlModelElement.GetConceptID(trans)).To(Equal(coreDomainID))
			PerformUndoRedoTest(1)
		})

		Describe("Test AddChild functionality", func() {
			Specify("AddChild Diagram should work", func() {
				newDiagram := FyneGUISingleton.addDiagram(cs1ID)
				Expect(newDiagram).ToNot(BeNil())
				Expect(newDiagram.IsRefinementOfURI(crldiagramdomain.CrlDiagramURI, trans)).To(BeTrue())
				Expect(newDiagram.GetOwningConcept(trans)).To(Equal(cs1))
				PerformUndoRedoTest(1)
			})
			Specify("AddChild Element should work", func() {
				el := FyneGUISingleton.addElement(cs1ID, "")
				Expect(el).ToNot(BeNil())
				Expect(el.GetOwningConcept(trans)).To(Equal(cs1))
				PerformUndoRedoTest(1)
			})
			Specify("AddChild Literal should work", func() {
				el := FyneGUISingleton.addLiteral(cs1ID, "")
				Expect(el).ToNot(BeNil())
				Expect(el.GetOwningConcept(trans)).To(Equal(cs1))
				isLiteral := false
				switch el.(type) {
				case core.Literal:
					isLiteral = true
				}
				Expect(isLiteral).To(BeTrue())
				PerformUndoRedoTest(1)
			})
			Specify("AddChild Reference should work", func() {
				el := FyneGUISingleton.addReference(cs1ID, "")
				Expect(el).ToNot(BeNil())
				Expect(el.GetOwningConcept(trans)).To(Equal(cs1))
				isReference := false
				switch el.(type) {
				case core.Reference:
					isReference = true
				}
				Expect(isReference).To(BeTrue())
				PerformUndoRedoTest(1)
			})
			Specify("AddChild Refinement should work", func() {
				el := FyneGUISingleton.addRefinement(cs1ID, "")
				Expect(el).ToNot(BeNil())
				Expect(el.GetOwningConcept(trans)).To(Equal(cs1))
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
			var selectedDiagram *diagramwidget.DiagramWidget
			BeforeEach(func() {
				selectedDiagram = FyneGUISingleton.diagramManager.GetSelectedDiagram()
				Expect(len(selectedDiagram.Nodes)).To(Equal(0))
				Expect(len(selectedDiagram.Links)).To(Equal(0))
			})
			Specify("Element node creation should work", func() {
				crlModelElement, fyneNode := createElementAt(selectedDiagram, 100, 100, trans)
				Expect(fyneNode).ToNot(BeNil())
				Expect(crlModelElement).ToNot(BeNil())
				Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(cs1ID))
				// Verify undo/redo
				PerformUndoRedoTest(1)
			})
			Specify("Literal node creation should work", func() {
				crlModelLiteral, fyneNode := createLiteralAt(selectedDiagram, 100, 100, trans)
				Expect(fyneNode).ToNot(BeNil())
				Expect(crlModelLiteral).ToNot(BeNil())
				// Verify undo/redo
				PerformUndoRedoTest(1)
			})
			Specify("Reference node creation should work", func() {
				crlModelReference, fyneNode := createReferenceAt(selectedDiagram, 100, 100, trans)
				Expect(fyneNode).ToNot(BeNil())
				Expect(crlModelReference).ToNot(BeNil())
				// Verify undo/redo
				PerformUndoRedoTest(1)
			})
			Specify("Refinement node creation should work", func() {
				crlModelRefinement, fyneNode := createRefinementAt(selectedDiagram, 100, 100, trans)
				Expect(fyneNode).ToNot(BeNil())
				Expect(crlModelRefinement).ToNot(BeNil())
				// Verify undo/redo
				PerformUndoRedoTest(1)
			})
			Describe("Reference link creation should work", func() {
				Specify("for a node source and target", func() {
					e1, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					e2, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					newReference, newReferenceView := createReferenceLink(selectedDiagram, e2View, e1View, trans)
					Expect(newReference).ToNot(BeNil())
					Expect(newReference.GetOwningConcept(trans)).To(Equal(e2))
					Expect(newReference.GetReferencedConcept(trans)).To(Equal(e1))
					Expect(newReferenceView).ToNot(BeNil())
					PerformUndoRedoTest(3)
				})
				// 	Specify("for a link source and node target", func() {
				// 		_, e1View := CreateElement(diagram, 100, 100)
				// 		_, e2View := CreateElement(diagram, 100, 200)
				// 		// create the node target
				// 		e3, e3View := CreateElement(diagram, 200, 150)
				// 		// Create a reference link
				// 		refLink1, refLink1View := CreateReferenceLink(diagram, e2View, e1View)
				// 		// Now the new reference
				// 		refLink2, _ := CreateReferenceLink(diagram, refLink1View, e3View)
				// 		// Now check the results
				// 		Expect(refLink2.GetOwningConceptID(trans)).To(Equal(refLink1.GetConceptID(trans)))
				// 		Expect(refLink2.GetReferencedConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
				// 		PerformUndoRedoTest(5)
				// 	})
				// 	Specify("for a node source and link target", func() {
				// 		_, e1View := CreateElement(diagram, 100, 100)
				// 		_, e2View := CreateElement(diagram, 100, 200)
				// 		// create the node source
				// 		e3, e3View := CreateElement(diagram, 200, 150)
				// 		// Create a reference link
				// 		refLink1, refLink1View := CreateReferenceLink(diagram, e2View, e1View)
				// 		// Now the new reference
				// 		refLink2, _ := CreateReferenceLink(diagram, e3View, refLink1View)
				// 		// Now check the results
				// 		Expect(refLink2.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
				// 		Expect(refLink2.GetReferencedConceptID(trans)).To(Equal(refLink1.GetConceptID(trans)))
				// 		PerformUndoRedoTest(5)
				// 	})
				// 	Specify("for a link source and link target", func() {
				// 		_, e1View := CreateElement(diagram, 100, 100)
				// 		_, e2View := CreateElement(diagram, 100, 200)
				// 		_, e3View := CreateElement(diagram, 200, 100)
				// 		_, e4View := CreateElement(diagram, 200, 200)
				// 		// Create the source reference link
				// 		refLink1, refLink1View := CreateReferenceLink(diagram, e2View, e1View)
				// 		// Create the target reference link
				// 		refLink2, refLink2View := CreateReferenceLink(diagram, e4View, e3View)
				// 		// Now the new reference
				// 		refLink3, _ := CreateReferenceLink(diagram, refLink1View, refLink2View)
				// 		// Now check the results
				// 		Expect(refLink3.GetOwningConceptID(trans)).To(Equal(refLink1.GetConceptID(trans)))
				// 		Expect(refLink3.GetReferencedConceptID(trans)).To(Equal(refLink2.GetConceptID(trans)))
				// 		PerformUndoRedoTest(7)
				// 	})
				// 	Specify("for a node source and an OwnerPointer target", func() {
				// 		_, e1View := CreateElement(diagram, 100, 100)
				// 		e2, e2View := CreateElement(diagram, 100, 200)
				// 		// create the node source
				// 		e3, e3View := CreateElement(diagram, 200, 150)
				// 		// create the owner pointer
				// 		opModelElement, opView := CreateOwnerPointer(diagram, e2View, e1View)
				// 		// Create the Reference
				// 		ref, refView := CreateReferenceLink(diagram, e3View, opView)
				// 		Expect(opModelElement.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
				// 		Expect(ref.GetReferencedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
				// 		Expect(ref.GetReferencedAttributeName(trans)).To(Equal(core.OwningConceptID))
				// 		Expect(ref.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
				// 		Expect(crldiagramdomain.GetLinkSource(refView, trans).GetConceptID(trans)).To(Equal(e3View.GetConceptID(trans)))
				// 		Expect(crldiagramdomain.GetLinkTarget(refView, trans).GetConceptID(trans)).To(Equal(opView.GetConceptID(trans)))
				// 		PerformUndoRedoTest(5)
				// 	})
				// 	Specify("for a node source and an ElementPointer target", func() {
				// 		_, e1View := CreateElement(diagram, 100, 100)
				// 		e2, e2View := CreateReferenceNode(diagram, 100, 200)
				// 		// create the node source
				// 		e3, e3View := CreateElement(diagram, 200, 150)
				// 		// create the owner pointer
				// 		epModelElement, epView := CreateElementPointer(diagram, e2View, e1View)
				// 		// Create the Reference
				// 		ref, refView := CreateReferenceLink(diagram, e3View, epView)
				// 		Expect(epModelElement.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
				// 		Expect(ref.GetReferencedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
				// 		Expect(ref.GetReferencedAttributeName(trans)).To(Equal(core.ReferencedConceptID))
				// 		Expect(ref.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
				// 		Expect(crldiagramdomain.GetLinkSource(refView, trans).GetConceptID(trans)).To(Equal(e3View.GetConceptID(trans)))
				// 		Expect(crldiagramdomain.GetLinkTarget(refView, trans).GetConceptID(trans)).To(Equal(epView.GetConceptID(trans)))
				// 		PerformUndoRedoTest(5)
				// 	})
				// 	Specify("for a node source and an AbstractPointer target", func() {
				// 		_, e1View := CreateElement(diagram, 100, 100)
				// 		e2, e2View := CreateRefinementNode(diagram, 100, 200)
				// 		// create the node source
				// 		e3, e3View := CreateElement(diagram, 200, 150)
				// 		// create the owner pointer
				// 		apModelElement, apView := CreateAbstractPointer(diagram, e2View, e1View)
				// 		// Create the Reference
				// 		ref, refView := CreateReferenceLink(diagram, e3View, apView)
				// 		Expect(apModelElement.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
				// 		Expect(ref.GetReferencedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
				// 		Expect(ref.GetReferencedAttributeName(trans)).To(Equal(core.AbstractConceptID))
				// 		Expect(ref.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
				// 		Expect(crldiagramdomain.GetLinkSource(refView, trans).GetConceptID(trans)).To(Equal(e3View.GetConceptID(trans)))
				// 		Expect(crldiagramdomain.GetLinkTarget(refView, trans).GetConceptID(trans)).To(Equal(apView.GetConceptID(trans)))
				// 		PerformUndoRedoTest(5)
				// 	})
				// 	Specify("for a node source and an RefinedPointer target", func() {
				// 		_, e1View := CreateElement(diagram, 100, 100)
				// 		e2, e2View := CreateRefinementNode(diagram, 100, 200)
				// 		// create the node source
				// 		e3, e3View := CreateElement(diagram, 200, 150)
				// 		// create the owner pointer
				// 		apModelElement, apView := CreateRefinedPointer(diagram, e2View, e1View)
				// 		// Create the Reference
				// 		ref, refView := CreateReferenceLink(diagram, e3View, apView)
				// 		Expect(apModelElement.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
				// 		Expect(ref.GetReferencedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
				// 		Expect(ref.GetReferencedAttributeName(trans)).To(Equal(core.RefinedConceptID))
				// 		Expect(ref.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
				// 		Expect(crldiagramdomain.GetLinkSource(refView, trans).GetConceptID(trans)).To(Equal(e3View.GetConceptID(trans)))
				// 		Expect(crldiagramdomain.GetLinkTarget(refView, trans).GetConceptID(trans)).To(Equal(apView.GetConceptID(trans)))
				// 		PerformUndoRedoTest(5)
				// 	})
			})
		})

	})

})

func beforeEachTest() (*core.UniverseOfDiscourse, *core.Transaction) {
	// Get current workspace path
	workspacePath := testWorkspaceDir
	// Open workspace (the same one - assumes nothing has been saved)
	crleditor.CrlEditorSingleton.Initialize(workspacePath, false)
	// log.Printf("Editor initialized with Workspace path: " + workspacePath)
	trans, _ := crleditor.CrlEditorSingleton.GetTransaction()
	Expect(trans).ToNot(BeNil())
	uOfD := trans.GetUniverseOfDiscourse()
	Expect(uOfD).ToNot(BeNil())
	return uOfD, trans
}

func afterEachTest() {
	// Clear existing workspace
	// log.Printf("**************************** About to hit ClearWorkspaceButton")
	clearWorkspaceItem := FyneGUISingleton.clearWorkspaceItem
	Expect(clearWorkspaceItem).ToNot(BeNil())
	clearWorkspaceItem.Action()
	crleditor.CrlEditorSingleton.EndTransaction()
}

func newDragEvent(dx float32, dy float32) *fyne.DragEvent {
	dragEvent := &fyne.DragEvent{}
	dragEvent.Dragged.DX = dx
	dragEvent.Dragged.DY = dy
	return dragEvent
}

func newPointEventAt(x float32, y float32) *fyne.PointEvent {
	inDiagramPosition := fyne.NewPos(x, y)
	absolutePosition := fyne.NewPos(0, 0)
	return &fyne.PointEvent{AbsolutePosition: absolutePosition, Position: inDiagramPosition}
}

func newLeftMouseEventAt(position fyne.Position, absolutePosition fyne.Position) *desktop.MouseEvent {
	newMouseEvent := &desktop.MouseEvent{}
	newMouseEvent.Position = position
	newMouseEvent.PointEvent.AbsolutePosition = absolutePosition
	newMouseEvent.Button = desktop.MouseButtonPrimary
	return newMouseEvent
}

func createElementAt(diagram *diagramwidget.DiagramWidget, x float32, y float32, trans *core.Transaction) (modelElement core.Element, fyneDiagramElement diagramwidget.DiagramElement) {
	// Create the element
	FyneGUISingleton.diagramManager.toolButtons[ELEMENT].Tapped(nil)
	diagram.Tapped(newPointEventAt(100, 100))
	// Get the Fyne node
	fyneNode := diagram.GetPrimarySelection()
	Expect(fyneNode).ToNot(BeNil())
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(fyneNode.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	// Get the model element
	crlModelElement := crldiagramdomain.GetReferencedModelElement(crlDiagramElement, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetOwningConceptID(trans)))
	return crlModelElement, fyneNode
}

func createLiteralAt(diagram *diagramwidget.DiagramWidget, x float32, y float32, trans *core.Transaction) (modelElement core.Literal, fyneDiagramElement diagramwidget.DiagramElement) {
	// Create the element
	FyneGUISingleton.diagramManager.toolButtons[LITERAL].Tapped(nil)
	diagram.Tapped(newPointEventAt(100, 100))
	// Get the Fyne node
	fyneNode := diagram.GetPrimarySelection()
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(fyneNode.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	// Get the model element
	crlModelElement := crldiagramdomain.GetReferencedModelElement(crlDiagramElement, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetOwningConceptID(trans)))
	Expect(crlModelElement.IsRefinementOfURI(core.LiteralURI, trans)).To(BeTrue())
	return crlModelElement.(core.Literal), fyneNode
}

func createReferenceAt(diagram *diagramwidget.DiagramWidget, x float32, y float32, trans *core.Transaction) (modelElement core.Reference, fyneDiagramElement diagramwidget.DiagramElement) {
	// Create the element
	FyneGUISingleton.diagramManager.toolButtons[REFERENCE].Tapped(nil)
	diagram.Tapped(newPointEventAt(100, 100))
	// Get the Fyne node
	fyneNode := diagram.GetPrimarySelection()
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(fyneNode.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	// Get the model element
	crlModelElement := crldiagramdomain.GetReferencedModelElement(crlDiagramElement, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetOwningConceptID(trans)))
	Expect(crlModelElement.IsRefinementOfURI(core.ReferenceURI, trans)).To(BeTrue())
	return crlModelElement.(core.Reference), fyneNode
}

func createRefinementAt(diagram *diagramwidget.DiagramWidget, x float32, y float32, trans *core.Transaction) (modelElement core.Refinement, fyneDiagramElement diagramwidget.DiagramElement) {
	// Create the element
	FyneGUISingleton.diagramManager.toolButtons[REFINEMENT].Tapped(nil)
	diagram.Tapped(newPointEventAt(100, 100))
	// Get the Fyne node
	fyneNode := diagram.GetPrimarySelection()
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(fyneNode.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	// Get the model element
	crlModelElement := crldiagramdomain.GetReferencedModelElement(crlDiagramElement, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetOwningConceptID(trans)))
	Expect(crlModelElement.IsRefinementOfURI(core.RefinementURI, trans)).To(BeTrue())
	return crlModelElement.(core.Refinement), fyneNode
}

func createReferenceLink(diagram *diagramwidget.DiagramWidget, sourceView diagramwidget.DiagramElement, targetView diagramwidget.DiagramElement, trans *core.Transaction) (newReference core.Reference, newReferenceView diagramwidget.DiagramLink) {
	FyneGUISingleton.diagramManager.toolButtons[REFERENCE_LINK].Tapped(nil)
	sourceViewPosition := FyneGUISingleton.app.Driver().AbsolutePositionForObject(sourceView)
	targetViewPosition := FyneGUISingleton.app.Driver().AbsolutePositionForObject(targetView)
	test.MoveMouse(FyneGUISingleton.window.Canvas(), sourceViewPosition)
	sourceView.GetDefaultConnectionPad().MouseDown(newLeftMouseEventAt(fyne.NewPos(0, 0), sourceViewPosition))
	Expect(diagram.ConnectionTransaction).ToNot(BeNil())
	newFyneLink := diagram.ConnectionTransaction.Link
	Expect(newFyneLink).ToNot(BeNil())
	targetView.GetDefaultConnectionPad().MouseIn(newLeftMouseEventAt(fyne.NewPos(0, 0), targetViewPosition))
	Expect(diagram.ConnectionTransaction.PendingPad).To(Equal(targetView.GetDefaultConnectionPad()))
	targetHandle := newFyneLink.GetTargetHandle()
	Expect(targetHandle).ToNot(BeNil())
	targetPadCurrentPosition := targetHandle.Position()
	dragEvent := newDragEvent(targetViewPosition.X-targetPadCurrentPosition.X, targetViewPosition.Y-targetPadCurrentPosition.Y)
	targetHandle.Dragged(dragEvent)
	targetHandle.DragEnd()
	Expect(newFyneLink.GetTargetPad()).To(Equal(targetView.GetDefaultConnectionPad()))
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlLink := trans.GetUniverseOfDiscourse().GetElement(newFyneLink.GetDiagramElementID())
	Expect(crlLink).ToNot(BeNil())
	Expect(crlLink.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	crlSource := trans.GetUniverseOfDiscourse().GetElement(sourceView.GetDiagramElementID())
	crlTarget := trans.GetUniverseOfDiscourse().GetElement(targetView.GetDiagramElementID())
	Expect(crldiagramdomain.GetLinkSource(crlLink, trans)).To(Equal(crlSource))
	Expect(crldiagramdomain.GetLinkTarget(crlLink, trans)).To(Equal(crlTarget))
	// Get the model element
	crlModelReference := crldiagramdomain.GetReferencedModelElement(crlLink, trans)
	crlSourceModelElement := crldiagramdomain.GetReferencedModelElement(crlSource, trans)
	crlTargetModelElement := crldiagramdomain.GetReferencedModelElement(crlTarget, trans)
	Expect(crlModelReference).ToNot(BeNil())
	Expect(crlModelReference.GetOwningConceptID(trans)).To(Equal(crlSourceModelElement.GetConceptID(trans)))
	Expect(crlModelReference.IsRefinementOfURI(core.ReferenceURI, trans)).To(BeTrue())
	Expect(crlModelReference.(core.Reference).GetReferencedConceptID(trans)).To(Equal(crlTargetModelElement.GetConceptID(trans)))

	return crlModelReference.(core.Reference), newFyneLink
}
