package crleditorfynegui

import (
	"log"
	"os"
	"testing"
	"time"

	// "time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"

	"github.com/pbrown12303/activeCRL/crleditor"
)

var testRootDir string
var testUserDir string
var testWorkspaceDir string
var testT *testing.T

var _ = BeforeSuite(func() {
	var err error
	// Get the tempDir
	tempDirPath := os.TempDir()
	log.Printf("TempDirPath: " + tempDirPath)
	err = os.Mkdir(tempDirPath, os.ModeDir)
	if !(err == nil || os.IsExist(err)) {
		Expect(err).NotTo(HaveOccurred())
	}
	log.Printf("TempDir created")

	testRootDir, err = os.MkdirTemp(tempDirPath, "crlEditorTestDir*")
	Expect(err).NotTo(HaveOccurred())
	testUserDir = testRootDir + "/testUserDir"
	err = os.Mkdir(testUserDir, os.ModeDir)
	Expect(err).NotTo(HaveOccurred())
	testWorkspaceDir = testRootDir + "/testWorkspace"
	err = os.Mkdir(testWorkspaceDir, os.ModeDir)
	Expect(err).NotTo(HaveOccurred())

	// Common infrastructure
	crlEditor := crleditor.NewEditor(testWorkspaceDir)
	crleditor.CrlEditorSingleton = crlEditor
	err = crleditor.CrlEditorSingleton.Initialize(testWorkspaceDir, true)
	if err != nil {
		log.Fatal(err)
	}

	// create a test app instead of a normal Fyne app.
	FyneGUISingleton = NewFyneGUI(crlEditor, test.NewApp())
	initialSize := fyne.NewSize(1600.0, 900.0)
	FyneGUISingleton.GetWindow().Resize(initialSize)
	FyneGUISingleton.GetWindow().ShowAndRun()
	time.Sleep(1000 * time.Millisecond)
	Expect(test.AssertRendersToImage(testT, "afterSuiteInitializaqtion.png", FyneGUISingleton.window.Canvas())).To(BeTrue())
})

func TestCrlEditor(t *testing.T) {
	testT = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "CrlEditorFyneGUI Suite")
}
