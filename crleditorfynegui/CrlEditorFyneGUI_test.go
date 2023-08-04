package crleditorfynegui

import (
	//	"fmt"

	"time"

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
			time.Sleep(100 * time.Millisecond)
			Expect(test.AssertRendersToMarkup(testT, "initialScreen.xml", FyneGUISingleton.window.Canvas())).To(BeTrue())
		})
		Specify("Tree selection and crleditor selection should match", func() {
			coreDomain := uOfD.GetElementWithURI(core.CoreDomainURI)
			coreDomainID := coreDomain.GetConceptID(trans)
			Expect(FyneGUISingleton.editor.SelectElement(coreDomain, trans)).To(Succeed())
			time.Sleep(100 * time.Millisecond)
			Expect(test.AssertRendersToMarkup(testT, "/selectionTests/coreDomainSelectedViaEditor.xml", FyneGUISingleton.window.Canvas())).To(BeTrue())
			Expect(FyneGUISingleton.editor.GetCurrentSelectionID(trans)).To(Equal(coreDomainID))
			Expect(FyneGUISingleton.editor.SelectElementUsingIDString("", trans)).To(Succeed())
			Expect(FyneGUISingleton.editor.GetCurrentSelectionID(trans)).To(Equal(""))
			time.Sleep(100 * time.Millisecond)
			Expect(test.AssertRendersToMarkup(testT, "/selectionTests/noSelection.xml", FyneGUISingleton.window.Canvas())).To(BeTrue())
			FyneGUISingleton.treeManager.tree.Select(coreDomainID)
			time.Sleep(100 * time.Millisecond)
			Expect(test.AssertRendersToMarkup(testT, "/selectionTests/coreDomainSelectedViaTree.xml", FyneGUISingleton.window.Canvas())).To(BeTrue())
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
			time.Sleep(100 * time.Millisecond)
			Expect(test.AssertRendersToMarkup(testT, "/dragDropTests/PriorToFirstNodeDrop.xml", FyneGUISingleton.window.Canvas())).To(BeTrue())
			treeNode.Dragged(newDragEvent(diagramPosition.X-treeNodePosition.X+100, diagramPosition.Y-treeNodePosition.Y+100))
			Expect(FyneGUISingleton.dragDropTransaction).ToNot(BeNil())
			Expect(FyneGUISingleton.dragDropTransaction.id).ToNot(Equal(""))
			fyneDiagram.MouseMoved(newLeftMouseEventAt(fyne.NewPos(100, 100), diagramPosition.AddXY(100, 100)))
			Expect(FyneGUISingleton.dragDropTransaction.currentDiagramMousePosition).To(Equal(fyne.NewPos(100, 100)))
			treeNode.DragEnd()
			time.Sleep(100 * time.Millisecond)
			Expect(test.AssertRendersToMarkup(testT, "/dragDropTests/FirstNodeDrop.xml", FyneGUISingleton.window.Canvas())).To(BeTrue())
			fyneNode := fyneDiagram.GetPrimarySelection()
			Expect(fyneNode).ToNot(BeNil())
			newNodeID := fyneNode.GetDiagramElementID()
			crlNode := uOfD.GetElement(newNodeID)
			Expect(crlNode).ToNot(BeNil())
			crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlNode, trans)
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
				crlModelElement, _, fyneNode := createElementAt(selectedDiagram, 100, 100, trans)
				Expect(fyneNode).ToNot(BeNil())
				Expect(crlModelElement).ToNot(BeNil())
				Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(cs1ID))
				// Verify undo/redo
				PerformUndoRedoTest(1)
			})
			Specify("Literal node creation should work", func() {
				crlModelLiteral, _, fyneNode := createLiteralAt(selectedDiagram, 100, 100, trans)
				Expect(fyneNode).ToNot(BeNil())
				Expect(crlModelLiteral).ToNot(BeNil())
				// Verify undo/redo
				PerformUndoRedoTest(1)
			})
			Specify("Reference node creation should work", func() {
				crlModelReference, _, fyneNode := createReferenceAt(selectedDiagram, 100, 100, trans)
				Expect(fyneNode).ToNot(BeNil())
				Expect(crlModelReference).ToNot(BeNil())
				// Verify undo/redo
				PerformUndoRedoTest(1)
			})
			Specify("Refinement node creation should work", func() {
				crlModelRefinement, _, fyneNode := createRefinementAt(selectedDiagram, 100, 100, trans)
				Expect(fyneNode).ToNot(BeNil())
				Expect(crlModelRefinement).ToNot(BeNil())
				// Verify undo/redo
				PerformUndoRedoTest(1)
			})
			Describe("Reference link creation should work", func() {
				Specify("for a node source and target", func() {
					e1, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					e2, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					newReference, _, newReferenceView := createReferenceLink(selectedDiagram, e2View, e1View, trans)
					Expect(newReference).ToNot(BeNil())
					Expect(newReference.GetOwningConcept(trans)).To(Equal(e2))
					Expect(newReference.GetReferencedConcept(trans)).To(Equal(e1))
					Expect(newReferenceView).ToNot(BeNil())
					PerformUndoRedoTest(3)
				})
				Specify("for a link source and node target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					// create the node target
					e3, _, e3View := createElementAt(selectedDiagram, 200, 150, trans)
					// Create a reference link
					refLink1, _, refLink1View := createReferenceLink(selectedDiagram, e2View, e1View, trans)
					// Now the new reference
					refLink2, _, refLink2View := createReferenceLink(selectedDiagram, refLink1View, e3View, trans)
					// Now check the results
					Expect(refLink2).ToNot(BeNil())
					Expect(refLink2View).ToNot(BeNil())
					Expect(refLink2.GetOwningConceptID(trans)).To(Equal(refLink1.GetConceptID(trans)))
					Expect(refLink2.GetReferencedConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and link target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					// create the node target
					e3, _, e3View := createElementAt(selectedDiagram, 200, 150, trans)
					// Create a reference link
					refLink1, _, refLink1View := createReferenceLink(selectedDiagram, e2View, e1View, trans)
					// Now the new reference
					refLink2, _, refLink2View := createReferenceLink(selectedDiagram, e3View, refLink1View, trans)
					// Now check the results
					Expect(refLink2).ToNot(BeNil())
					Expect(refLink2View).ToNot(BeNil())
					Expect(refLink2.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
					Expect(refLink2.GetReferencedConceptID(trans)).To(Equal(refLink1.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a link source and link target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					_, _, e3View := createElementAt(selectedDiagram, 200, 100, trans)
					_, _, e4View := createElementAt(selectedDiagram, 200, 200, trans)
					// Create the source reference link
					refLink1, _, refLink1View := createReferenceLink(selectedDiagram, e2View, e1View, trans)
					// Create the target reference link
					refLink2, _, refLink2View := createReferenceLink(selectedDiagram, e4View, e3View, trans)
					// Now the new reference
					refLink3, _, refLink3View := createReferenceLink(selectedDiagram, refLink1View, refLink2View, trans)
					// Now check the results
					Expect(refLink3).ToNot(BeNil())
					Expect(refLink3View).ToNot(BeNil())
					Expect(refLink3.GetOwningConceptID(trans)).To(Equal(refLink1.GetConceptID(trans)))
					Expect(refLink3.GetReferencedConceptID(trans)).To(Equal(refLink2.GetConceptID(trans)))
					PerformUndoRedoTest(7)
				})
				Specify("for a node source and an OwnerPointer target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					e2, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					// create the node source
					e3, _, e3LinkView := createElementAt(selectedDiagram, 200, 150, trans)
					// create the owner pointer
					opModelElement, _, opView := createOwnerPointer(selectedDiagram, e2View, e1View, trans)
					// Create the Reference
					ref, _, refLinkView := createReferenceLink(selectedDiagram, e3LinkView, opView, trans)
					Expect(opModelElement.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(ref.GetReferencedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(ref.GetReferencedAttributeName(trans)).To(Equal(core.OwningConceptID))
					Expect(ref.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
					refView := uOfD.GetElement(refLinkView.GetDiagramElementID())
					Expect(refView).ToNot(BeNil())
					e3View := uOfD.GetElement(e3LinkView.GetDiagramElementID())
					Expect(e3View).ToNot(BeNil())
					Expect(crldiagramdomain.GetLinkSource(refView, trans).GetConceptID(trans)).To(Equal(e3View.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(refView, trans).GetConceptID(trans)).To(Equal(opView.GetDiagramElementID()))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and an ElementPointer target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					e2, _, e2View := createReferenceAt(selectedDiagram, 100, 200, trans)
					// create the node source
					e3, e3CrlView, e3View := createElementAt(selectedDiagram, 200, 150, trans)
					// create the referenced element pointer
					epModelElement, epCrlView, epView := createReferencedElementPointer(selectedDiagram, e2View, e1View, trans)
					// Create the Reference
					ref, crlRefView, _ := createReferenceLink(selectedDiagram, e3View, epView, trans)
					Expect(epModelElement.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(ref.GetReferencedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(ref.GetReferencedAttributeName(trans)).To(Equal(core.ReferencedConceptID))
					Expect(ref.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlRefView, trans).GetConceptID(trans)).To(Equal(e3CrlView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlRefView, trans).GetConceptID(trans)).To(Equal(epCrlView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and an AbstractPointer target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					e2, _, e2View := createRefinementAt(selectedDiagram, 100, 200, trans)
					// create the node source
					e3, crlE3View, e3View := createElementAt(selectedDiagram, 200, 150, trans)
					// create the owner pointer
					apModelElement, crlApView, apView := createAbstractPointer(selectedDiagram, e2View, e1View, trans)
					// Create the Reference
					ref, crlRefView, _ := createReferenceLink(selectedDiagram, e3View, apView, trans)
					Expect(apModelElement.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(ref.GetReferencedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(ref.GetReferencedAttributeName(trans)).To(Equal(core.AbstractConceptID))
					Expect(ref.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlRefView, trans).GetConceptID(trans)).To(Equal(crlE3View.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlRefView, trans).GetConceptID(trans)).To(Equal(crlApView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and an RefinedPointer target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					e2, _, e2View := createRefinementAt(selectedDiagram, 100, 200, trans)
					// create the node source
					e3, crlE3View, e3View := createElementAt(selectedDiagram, 200, 150, trans)
					// create the owner pointer
					apModelElement, crlAppView, apView := createRefinedPointer(selectedDiagram, e2View, e1View, trans)
					// Create the Reference
					ref, crlRefView, _ := createReferenceLink(selectedDiagram, e3View, apView, trans)
					Expect(apModelElement.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(ref.GetReferencedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(ref.GetReferencedAttributeName(trans)).To(Equal(core.RefinedConceptID))
					Expect(ref.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlRefView, trans).GetConceptID(trans)).To(Equal(crlE3View.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlRefView, trans).GetConceptID(trans)).To(Equal(crlAppView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
			})
			Describe("Refinement link creation should work", func() {
				Specify("for a node source and node target", func() {
					e1, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					e2, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					newRefinement, _, _ := createRefinementLink(selectedDiagram, e2View, e1View, trans)
					Expect(newRefinement.GetAbstractConceptID(trans)).To(Equal(e1.GetConceptID(trans)))
					Expect(newRefinement.GetRefinedConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					PerformUndoRedoTest(3)
				})
				Specify("for a link source and node target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					source, crlSourceView, sourceView := createRefinementLink(selectedDiagram, e2View, e1View, trans)
					target, crlTargetView, targetView := createRefinementAt(selectedDiagram, 200, 150, trans)
					newRefinement, crlNewRefinementView, _ := createRefinementLink(selectedDiagram, sourceView, targetView, trans)
					Expect(newRefinement.GetAbstractConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(newRefinement.GetRefinedConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlNewRefinementView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlNewRefinementView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and link target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					target, crlTargetView, targetView := createRefinementLink(selectedDiagram, e2View, e1View, trans)
					source, crlSourceView, sourceView := createRefinementAt(selectedDiagram, 200, 150, trans)
					newRefinement, crlNewRefinementView, _ := createRefinementLink(selectedDiagram, sourceView, targetView, trans)
					Expect(newRefinement.GetAbstractConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(newRefinement.GetRefinedConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlNewRefinementView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlNewRefinementView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
			})
			Describe("OwnerPointer creation should work", func() {
				Specify("For a node source and node target", func() {
					e1, crlE1View, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					e2, crlE2View, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					source, crlOwnerPointerView, ownerPointerView := createOwnerPointer(selectedDiagram, e2View, e1View, trans)
					// Now check the results
					Expect(ownerPointerView).ToNot(BeNil())
					Expect(source.GetConceptID(trans)).To(Equal(e2.GetConceptID(trans)))
					Expect(e2.GetOwningConceptID(trans)).To(Equal(e1.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlOwnerPointerView, trans).GetConceptID(trans)).To(Equal(crlE2View.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlOwnerPointerView, trans).GetConceptID(trans)).To(Equal(crlE1View.GetConceptID(trans)))
					PerformUndoRedoTest(3)
				})
				Specify("For a Refinement Link source and node target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					// Create the refinement link
					source, crlSourceView, sourceView := createRefinementLink(selectedDiagram, e1View, e2View, trans)
					// create the new owner
					e3, crlE3View, e3View := createElementAt(selectedDiagram, 200, 150, trans)
					// Now the ownerPointer
					ownerPointerConcept, crlOwnerPointerView, _ := createOwnerPointer(selectedDiagram, sourceView, e3View, trans)
					// Now check the results
					Expect(source.GetOwningConceptID(trans)).To(Equal(e3.GetConceptID(trans)))
					Expect(source.GetConceptID(trans)).To(Equal(ownerPointerConcept.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlOwnerPointerView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlOwnerPointerView, trans).GetConceptID(trans)).To(Equal(crlE3View.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("For a node source and ReferenceLink target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					// Create the Reference link
					target, crlTargetView, targetView := createReferenceLink(selectedDiagram, e1View, e2View, trans)
					// create the new owner
					source, crlSourceView, sourceView := createElementAt(selectedDiagram, 200, 150, trans)
					// Now the ownerPointer
					ownerPointerConcept, crlOwnerPointerView, _ := createOwnerPointer(selectedDiagram, sourceView, targetView, trans)
					// Now check the results
					Expect(source.GetOwningConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetConceptID(trans)).To(Equal(ownerPointerConcept.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlOwnerPointerView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlOwnerPointerView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("For a node source and RefinementLink target", func() {
					_, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 100, 200, trans)
					// Create the refinement link
					target, crlTargetView, targetView := createRefinementLink(selectedDiagram, e1View, e2View, trans)
					// create the new owner
					source, crlSourceView, sourceView := createElementAt(selectedDiagram, 200, 150, trans)
					// Now the ownerPointer
					ownerPointerConcept, crlOwnerPointerView, _ := createOwnerPointer(selectedDiagram, sourceView, targetView, trans)
					// Now check the results
					Expect(source.GetOwningConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetConceptID(trans)).To(Equal(ownerPointerConcept.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkSource(crlOwnerPointerView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlOwnerPointerView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
			})
			Describe("ElementPointer creation should work", func() {
				Specify("for a node source and node target", func() {
					target, crlTargetView, targetView := createElementAt(selectedDiagram, 100, 100, trans)
					source, crlSourceView, sourceView := createReferenceAt(selectedDiagram, 100, 200, trans)
					reference, crlEpView, _ := createReferencedElementPointer(selectedDiagram, sourceView, targetView, trans)
					// Now check the results
					Expect(source.GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(reference.GetConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(reference.(core.Reference).GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetReferencedAttributeName(trans)).To(Equal(core.NoAttribute))
					Expect(crldiagramdomain.GetLinkSource(crlEpView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlEpView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					modelReference := crldiagramdomain.GetReferencedModelConcept(crlEpView, trans)
					Expect(modelReference).ToNot(BeNil())
					Expect(modelReference.IsRefinementOfURI(core.ReferenceURI, trans)).To(BeTrue())
					Expect(modelReference.(core.Reference).GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					PerformUndoRedoTest(3)
				})
				Specify("for a node source and reference link target", func() {
					source, crlSourceView, sourceView := createReferenceAt(selectedDiagram, 100, 150, trans)
					_, _, e1View := createElementAt(selectedDiagram, 200, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 200, 200, trans)
					target, crlTargetView, targetView := createReferenceLink(selectedDiagram, e1View, e2View, trans)
					epModel, crlEpView, _ := createReferencedElementPointer(selectedDiagram, sourceView, targetView, trans)
					Expect(epModel.GetConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(source.GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetReferencedAttributeName(trans)).To(Equal(core.NoAttribute))
					Expect(crldiagramdomain.GetLinkSource(crlEpView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlEpView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and RefinementLink target", func() {
					source, crlSourceView, sourceView := createReferenceAt(selectedDiagram, 100, 150, trans)
					_, _, e1View := createElementAt(selectedDiagram, 200, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 200, 200, trans)
					target, crlTargetView, targetView := createRefinementLink(selectedDiagram, e1View, e2View, trans)
					epModel, crlEpView, _ := createReferencedElementPointer(selectedDiagram, sourceView, targetView, trans)
					Expect(epModel.GetConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(source.GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetReferencedAttributeName(trans)).To(Equal(core.NoAttribute))
					Expect(crldiagramdomain.GetLinkSource(crlEpView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlEpView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and an OwnerPointer target", func() {
					source, crlSourceView, sourceView := createReferenceAt(selectedDiagram, 100, 150, trans)
					e1, _, e1View := createElementAt(selectedDiagram, 200, 100, trans)
					e2, _, e2View := createElementAt(selectedDiagram, 200, 200, trans)
					Expect(e1).ToNot(BeNil())
					Expect(e2).ToNot(BeNil())
					target, crlTargetView, targetView := createOwnerPointer(selectedDiagram, e1View, e2View, trans)
					epModel, crlEpView, _ := createReferencedElementPointer(selectedDiagram, sourceView, targetView, trans)
					Expect(epModel.GetConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(source.GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetReferencedAttributeName(trans)).To(Equal(core.OwningConceptID))
					Expect(crldiagramdomain.GetLinkSource(crlEpView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlEpView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and an ElementPointer target", func() {
					source, crlSourceView, sourceView := createReferenceAt(selectedDiagram, 100, 150, trans)
					_, _, e1View := createReferenceAt(selectedDiagram, 200, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 200, 200, trans)
					target, crlTargetView, targetView := createReferencedElementPointer(selectedDiagram, e1View, e2View, trans)
					epModel, crlEpView, _ := createReferencedElementPointer(selectedDiagram, sourceView, targetView, trans)
					Expect(epModel.GetConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(source.GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetReferencedAttributeName(trans)).To(Equal(core.ReferencedConceptID))
					Expect(crldiagramdomain.GetLinkSource(crlEpView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlEpView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and an AbstractPointer target", func() {
					source, crlSourceView, sourceView := createReferenceAt(selectedDiagram, 100, 150, trans)
					_, _, e1View := createRefinementAt(selectedDiagram, 200, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 200, 200, trans)
					target, crlTargetView, targetView := createAbstractPointer(selectedDiagram, e1View, e2View, trans)
					epModel, crlEpView, _ := createReferencedElementPointer(selectedDiagram, sourceView, targetView, trans)
					Expect(epModel.GetConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(source.GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetReferencedAttributeName(trans)).To(Equal(core.AbstractConceptID))
					Expect(crldiagramdomain.GetLinkSource(crlEpView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlEpView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
				Specify("for a node source and an RefinedPointer target", func() {
					source, crlSourceView, sourceView := createReferenceAt(selectedDiagram, 100, 150, trans)
					_, _, e1View := createRefinementAt(selectedDiagram, 200, 100, trans)
					_, _, e2View := createElementAt(selectedDiagram, 200, 200, trans)
					target, crlTargetView, targetView := createRefinedPointer(selectedDiagram, e1View, e2View, trans)
					epModel, crlEpView, _ := createReferencedElementPointer(selectedDiagram, sourceView, targetView, trans)
					Expect(epModel.GetConceptID(trans)).To(Equal(source.GetConceptID(trans)))
					Expect(source.GetReferencedConceptID(trans)).To(Equal(target.GetConceptID(trans)))
					Expect(source.GetReferencedAttributeName(trans)).To(Equal(core.RefinedConceptID))
					Expect(crldiagramdomain.GetLinkSource(crlEpView, trans).GetConceptID(trans)).To(Equal(crlSourceView.GetConceptID(trans)))
					Expect(crldiagramdomain.GetLinkTarget(crlEpView, trans).GetConceptID(trans)).To(Equal(crlTargetView.GetConceptID(trans)))
					PerformUndoRedoTest(5)
				})
			})
			Specify("AbstractPointer creation should work", func() {
				e1, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
				r1, _, r1View := createRefinementAt(selectedDiagram, 100, 200, trans)
				createAbstractPointer(selectedDiagram, r1View, e1View, trans)
				Expect(r1.GetAbstractConcept(trans).GetConceptID((trans))).To(Equal(e1.GetConceptID(trans)))
				PerformUndoRedoTest(3)
			})
			Specify("RefinedPointer creation should work", func() {
				e1, _, e1View := createElementAt(selectedDiagram, 100, 100, trans)
				r1, _, r1View := createRefinementAt(selectedDiagram, 100, 200, trans)
				createRefinedPointer(selectedDiagram, r1View, e1View, trans)
				Expect(r1.GetRefinedConcept(trans)).To(Equal(e1))
				PerformUndoRedoTest(3)
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

func createElementAt(diagram *diagramwidget.DiagramWidget, x float32, y float32, trans *core.Transaction) (core.Element, core.Element, diagramwidget.DiagramElement) {
	// Create the element
	FyneGUISingleton.diagramManager.toolButtons[ElementSelected].Tapped(nil)
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
	crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetOwningConceptID(trans)))
	return crlModelElement, crlDiagramElement, fyneNode
}

func createLiteralAt(diagram *diagramwidget.DiagramWidget, x float32, y float32, trans *core.Transaction) (core.Literal, core.Element, diagramwidget.DiagramElement) {
	// Create the element
	FyneGUISingleton.diagramManager.toolButtons[LiteralSelected].Tapped(nil)
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
	crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetOwningConceptID(trans)))
	Expect(crlModelElement.IsRefinementOfURI(core.LiteralURI, trans)).To(BeTrue())
	return crlModelElement.(core.Literal), crlDiagramElement, fyneNode
}

func createReferenceAt(diagram *diagramwidget.DiagramWidget, x float32, y float32, trans *core.Transaction) (core.Reference, core.Element, diagramwidget.DiagramElement) {
	// Create the element
	FyneGUISingleton.diagramManager.toolButtons[ReferenceSelected].Tapped(nil)
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
	crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetOwningConceptID(trans)))
	Expect(crlModelElement.IsRefinementOfURI(core.ReferenceURI, trans)).To(BeTrue())
	return crlModelElement.(core.Reference), crlDiagramElement, fyneNode
}

func createRefinementAt(diagram *diagramwidget.DiagramWidget, x float32, y float32, trans *core.Transaction) (core.Refinement, core.Element, diagramwidget.DiagramElement) {
	// Create the element
	FyneGUISingleton.diagramManager.toolButtons[RefinementSelected].Tapped(nil)
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
	crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetOwningConceptID(trans)))
	Expect(crlModelElement.IsRefinementOfURI(core.RefinementURI, trans)).To(BeTrue())
	return crlModelElement.(core.Refinement), crlDiagramElement, fyneNode
}

func createReferenceLink(diagram *diagramwidget.DiagramWidget, sourceView diagramwidget.DiagramElement, targetView diagramwidget.DiagramElement, trans *core.Transaction) (core.Reference, core.Element, diagramwidget.DiagramLink) {
	FyneGUISingleton.diagramManager.toolButtons[ReferenceLinkSelected].Tapped(nil)
	newFyneLink := createFyneLink(sourceView, targetView, diagram)
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(newFyneLink.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	crlSource := trans.GetUniverseOfDiscourse().GetElement(sourceView.GetDiagramElementID())
	crlTarget := trans.GetUniverseOfDiscourse().GetElement(targetView.GetDiagramElementID())
	Expect(crldiagramdomain.GetLinkSource(crlDiagramElement, trans)).To(Equal(crlSource))
	Expect(crldiagramdomain.GetLinkTarget(crlDiagramElement, trans)).To(Equal(crlTarget))
	// Get the model element
	crlModelReference := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	crlSourceModelElement := crldiagramdomain.GetReferencedModelConcept(crlSource, trans)
	crlTargetModelElement := crldiagramdomain.GetReferencedModelConcept(crlTarget, trans)
	Expect(crlModelReference).ToNot(BeNil())
	Expect(crlModelReference.GetOwningConceptID(trans)).To(Equal(crlSourceModelElement.GetConceptID(trans)))
	Expect(crlModelReference.IsRefinementOfURI(core.ReferenceURI, trans)).To(BeTrue())
	Expect(crlModelReference.(core.Reference).GetReferencedConceptID(trans)).To(Equal(crlTargetModelElement.GetConceptID(trans)))

	return crlModelReference.(core.Reference), crlDiagramElement, newFyneLink
}

func createRefinementLink(diagram *diagramwidget.DiagramWidget, sourceView diagramwidget.DiagramElement, targetView diagramwidget.DiagramElement, trans *core.Transaction) (core.Refinement, core.Element, diagramwidget.DiagramLink) {
	FyneGUISingleton.diagramManager.toolButtons[RefinementLinkSelected].Tapped(nil)
	newFyneLink := createFyneLink(sourceView, targetView, diagram)
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(newFyneLink.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	crlSource := trans.GetUniverseOfDiscourse().GetElement(sourceView.GetDiagramElementID())
	crlTarget := trans.GetUniverseOfDiscourse().GetElement(targetView.GetDiagramElementID())
	Expect(crldiagramdomain.GetLinkSource(crlDiagramElement, trans)).To(Equal(crlSource))
	Expect(crldiagramdomain.GetLinkTarget(crlDiagramElement, trans)).To(Equal(crlTarget))
	// Get the model element
	crlModelRefinement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	crlSourceModelElement := crldiagramdomain.GetReferencedModelConcept(crlSource, trans)
	crlTargetModelElement := crldiagramdomain.GetReferencedModelConcept(crlTarget, trans)
	Expect(crlModelRefinement).ToNot(BeNil())
	Expect(crlModelRefinement.GetOwningConceptID(trans)).To(Equal(crlSourceModelElement.GetConceptID(trans)))
	Expect(crlModelRefinement.IsRefinementOfURI(core.RefinementURI, trans)).To(BeTrue())
	Expect(crlModelRefinement.(core.Refinement).GetAbstractConceptID(trans)).To(Equal(crlTargetModelElement.GetConceptID(trans)))
	Expect(crlModelRefinement.(core.Refinement).GetRefinedConceptID(trans)).To(Equal(crlSourceModelElement.GetConceptID(trans)))

	return crlModelRefinement.(core.Refinement), crlDiagramElement, newFyneLink
}

func createOwnerPointer(diagram *diagramwidget.DiagramWidget, sourceView diagramwidget.DiagramElement, targetView diagramwidget.DiagramElement, trans *core.Transaction) (core.Element, core.Element, diagramwidget.DiagramLink) {
	FyneGUISingleton.diagramManager.toolButtons[OwnerPointerSelected].Tapped(nil)
	newFyneLink := createFyneLink(sourceView, targetView, diagram)
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(newFyneLink.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	Expect(crlDiagramElement.IsRefinementOfURI(crldiagramdomain.CrlDiagramOwnerPointerURI, trans)).To(BeTrue())
	crlSource := trans.GetUniverseOfDiscourse().GetElement(sourceView.GetDiagramElementID())
	crlTarget := trans.GetUniverseOfDiscourse().GetElement(targetView.GetDiagramElementID())
	Expect(crldiagramdomain.GetLinkSource(crlDiagramElement, trans)).To(Equal(crlSource))
	Expect(crldiagramdomain.GetLinkTarget(crlDiagramElement, trans)).To(Equal(crlTarget))
	// Get the model element
	crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	crlSourceModelElement := crldiagramdomain.GetReferencedModelConcept(crlSource, trans)
	crlTargetModelElement := crldiagramdomain.GetReferencedModelConcept(crlTarget, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.GetConceptID(trans)).To(Equal(crlSourceModelElement.GetConceptID(trans)))
	Expect(crlModelElement.GetOwningConceptID(trans)).To(Equal(crlTargetModelElement.GetConceptID(trans)))

	return crlModelElement, crlDiagramElement, newFyneLink
}

func createReferencedElementPointer(diagram *diagramwidget.DiagramWidget, sourceView diagramwidget.DiagramElement, targetView diagramwidget.DiagramElement, trans *core.Transaction) (core.Element, core.Element, diagramwidget.DiagramLink) {
	FyneGUISingleton.diagramManager.toolButtons[ReferencedElementPointerSelected].Tapped(nil)
	newFyneLink := createFyneLink(sourceView, targetView, diagram)
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(newFyneLink.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	Expect(crlDiagramElement.IsRefinementOfURI(crldiagramdomain.CrlDiagramElementPointerURI, trans)).To(BeTrue())
	crlSource := trans.GetUniverseOfDiscourse().GetElement(sourceView.GetDiagramElementID())
	crlTarget := trans.GetUniverseOfDiscourse().GetElement(targetView.GetDiagramElementID())
	Expect(crldiagramdomain.GetLinkSource(crlDiagramElement, trans)).To(Equal(crlSource))
	Expect(crldiagramdomain.GetLinkTarget(crlDiagramElement, trans)).To(Equal(crlTarget))
	// Get the model element
	crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	crlSourceModelElement := crldiagramdomain.GetReferencedModelConcept(crlSource, trans)
	crlTargetModelElement := crldiagramdomain.GetReferencedModelConcept(crlTarget, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.IsRefinementOfURI(core.ReferenceURI, trans)).To(BeTrue())
	Expect(crlModelElement.GetConceptID(trans)).To(Equal(crlSourceModelElement.GetConceptID(trans)))
	Expect((crlModelElement.(core.Reference)).GetReferencedConceptID(trans)).To(Equal(crlTargetModelElement.GetConceptID(trans)))

	return crlModelElement, crlDiagramElement, newFyneLink
}

func createAbstractPointer(diagram *diagramwidget.DiagramWidget, sourceView diagramwidget.DiagramElement, targetView diagramwidget.DiagramElement, trans *core.Transaction) (core.Element, core.Element, diagramwidget.DiagramLink) {
	FyneGUISingleton.diagramManager.toolButtons[AbstractElementPointerSelected].Tapped(nil)
	newFyneLink := createFyneLink(sourceView, targetView, diagram)
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(newFyneLink.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	Expect(crlDiagramElement.IsRefinementOfURI(crldiagramdomain.CrlDiagramAbstractPointerURI, trans)).To(BeTrue())
	crlSource := trans.GetUniverseOfDiscourse().GetElement(sourceView.GetDiagramElementID())
	crlTarget := trans.GetUniverseOfDiscourse().GetElement(targetView.GetDiagramElementID())
	Expect(crldiagramdomain.GetLinkSource(crlDiagramElement, trans)).To(Equal(crlSource))
	Expect(crldiagramdomain.GetLinkTarget(crlDiagramElement, trans)).To(Equal(crlTarget))
	// Get the model element
	crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	crlSourceModelElement := crldiagramdomain.GetReferencedModelConcept(crlSource, trans)
	crlTargetModelElement := crldiagramdomain.GetReferencedModelConcept(crlTarget, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.IsRefinementOfURI(core.RefinementURI, trans)).To(BeTrue())
	Expect(crlModelElement.GetConceptID(trans)).To(Equal(crlSourceModelElement.GetConceptID(trans)))
	Expect((crlModelElement.(core.Refinement)).GetAbstractConceptID(trans)).To(Equal(crlTargetModelElement.GetConceptID(trans)))

	return crlModelElement, crlDiagramElement, newFyneLink
}

func createRefinedPointer(diagram *diagramwidget.DiagramWidget, sourceView diagramwidget.DiagramElement, targetView diagramwidget.DiagramElement, trans *core.Transaction) (core.Element, core.Element, diagramwidget.DiagramLink) {
	FyneGUISingleton.diagramManager.toolButtons[RefinedElementPointerSelected].Tapped(nil)
	newFyneLink := createFyneLink(sourceView, targetView, diagram)
	// Get the CRL diagram
	crlDiagram := trans.GetUniverseOfDiscourse().GetElement(diagram.ID)
	// Get the CRL diagram element
	crlDiagramElement := trans.GetUniverseOfDiscourse().GetElement(newFyneLink.GetDiagramElementID())
	Expect(crlDiagramElement).ToNot(BeNil())
	Expect(crlDiagramElement.GetOwningConceptID(trans)).To(Equal(crlDiagram.GetConceptID(trans)))
	Expect(crlDiagramElement.IsRefinementOfURI(crldiagramdomain.CrlDiagramRefinedPointerURI, trans)).To(BeTrue())
	crlSource := trans.GetUniverseOfDiscourse().GetElement(sourceView.GetDiagramElementID())
	crlTarget := trans.GetUniverseOfDiscourse().GetElement(targetView.GetDiagramElementID())
	Expect(crldiagramdomain.GetLinkSource(crlDiagramElement, trans)).To(Equal(crlSource))
	Expect(crldiagramdomain.GetLinkTarget(crlDiagramElement, trans)).To(Equal(crlTarget))
	// Get the model element
	crlModelElement := crldiagramdomain.GetReferencedModelConcept(crlDiagramElement, trans)
	crlSourceModelElement := crldiagramdomain.GetReferencedModelConcept(crlSource, trans)
	crlTargetModelElement := crldiagramdomain.GetReferencedModelConcept(crlTarget, trans)
	Expect(crlModelElement).ToNot(BeNil())
	Expect(crlModelElement.IsRefinementOfURI(core.RefinementURI, trans)).To(BeTrue())
	Expect(crlModelElement.GetConceptID(trans)).To(Equal(crlSourceModelElement.GetConceptID(trans)))
	Expect((crlModelElement.(core.Refinement)).GetRefinedConceptID(trans)).To(Equal(crlTargetModelElement.GetConceptID(trans)))

	return crlModelElement, crlDiagramElement, newFyneLink
}

// createFyneLink assumes that the caller has already selected the appropriate toolbar button
func createFyneLink(sourceView diagramwidget.DiagramElement, targetView diagramwidget.DiagramElement, diagram *diagramwidget.DiagramWidget) diagramwidget.DiagramLink {
	sourceViewPosition := FyneGUISingleton.app.Driver().AbsolutePositionForObject(sourceView)
	targetViewPosition := FyneGUISingleton.app.Driver().AbsolutePositionForObject(targetView)
	test.MoveMouse(FyneGUISingleton.window.Canvas(), sourceViewPosition)
	sourceView.GetDefaultConnectionPad().MouseDown(newLeftMouseEventAt(fyne.NewPos(0, 0), sourceViewPosition))
	Expect(diagram.ConnectionTransaction).ToNot(BeNil())
	newFyneLink := diagram.ConnectionTransaction.Link
	Expect(newFyneLink).ToNot(BeNil())
	targetView.GetDefaultConnectionPad().MouseIn(newLeftMouseEventAt(fyne.NewPos(0, 0), targetViewPosition))
	Expect(diagram.ConnectionTransaction.PendingPad).ToNot(BeNil())
	Expect(diagram.ConnectionTransaction.PendingPad).To(Equal(targetView.GetDefaultConnectionPad()))
	targetHandle := newFyneLink.GetTargetHandle()
	Expect(targetHandle).ToNot(BeNil())
	targetPadCurrentPosition := targetHandle.Position()
	dragEvent := newDragEvent(targetViewPosition.X-targetPadCurrentPosition.X, targetViewPosition.Y-targetPadCurrentPosition.Y)
	targetHandle.Dragged(dragEvent)
	targetHandle.DragEnd()
	Expect(newFyneLink.GetTargetPad()).To(Equal(targetView.GetDefaultConnectionPad()))
	return newFyneLink
}
