package crleditorfynegui

import (
	//	"fmt"

	"fyne.io/fyne/v2/test"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"

	"github.com/pbrown12303/activeCRL/core"
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
			// _, cs1 := CreateDomain()
			// Expect(cs1).ToNot(BeNil())
			// uOfD.MarkUndoPoint()
			// beforeUofD := uOfD.Clone(hl)
			// beforeHL := beforeUofD.NewTransaction()
			// crleditorbrowsergui.CrlLogClientRequests = true
			// _, diag := CreateDiagram(cs1)
			// Expect(diag).ToNot(BeNil())
			// crleditorbrowsergui.CrlLogClientRequests = false
			// afterUofD := uOfD.Clone(hl)
			// afterHL := afterUofD.NewTransaction()
			// Undo()
			// Expect(uOfD.IsEquivalent(hl, beforeUofD, beforeHL, true)).To(BeTrue())
			// Redo()
			// Expect(uOfD.IsEquivalent(hl, afterUofD, afterHL, true)).To(BeTrue())
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
