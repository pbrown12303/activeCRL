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
		// Get current workspace path
		workspacePath := testWorkspaceDir
		// Open workspace (the same one - assumes nothing has been saved)
		crleditor.CrlEditorSingleton.Initialize(workspacePath, false)
		// log.Printf("Editor initialized with Workspace path: " + workspacePath)
		trans, _ = crleditor.CrlEditorSingleton.GetTransaction()
		Expect(trans).ToNot(BeNil())
		uOfD = trans.GetUniverseOfDiscourse()
		Expect(uOfD).ToNot(BeNil())
	})

	AfterEach(func() {
		// Clear existing workspace
		// log.Printf("**************************** About to hit ClearWorkspaceButton")
		clearWorkspaceItem := FyneGUISingleton.clearWorkspaceItem
		Expect(clearWorkspaceItem).ToNot(BeNil())
		clearWorkspaceItem.Action()
		crleditor.CrlEditorSingleton.EndTransaction()
	})

	Describe("Testing CrlEditor basic functionality", func() {
		Specify("The FyneGUISingleton should be populated", func() {
			Expect(FyneGUISingleton).ToNot(BeNil())
			coreDomain := uOfD.GetElementWithURI(core.CoreDomainURI)
			Expect(coreDomain).ToNot(BeNil())
			Expect(test.AssertRendersToImage(testT, "initialScreen.png", FyneGUISingleton.window.Canvas())).To(BeTrue())
		})

	})
})
