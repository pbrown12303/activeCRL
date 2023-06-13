package crleditorfynegui

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

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
var RequestFocus func()

var _ = BeforeSuite(func() {
	RequestFocus = func() {
		FyneGUISingleton.GetWindow().RequestFocus()
	}

	var err error
	// Get the tempDir
	tempDirPath := os.TempDir()
	log.Printf("TempDirPath: " + tempDirPath)
	err = os.Mkdir(tempDirPath, os.ModeDir)
	if !(err == nil || os.IsExist(err)) {
		Expect(err).NotTo(HaveOccurred())
	}
	log.Printf("TempDir created")

	testRootDir, err = ioutil.TempDir(tempDirPath, "crlEditorTestDir*")
	Expect(err).NotTo(HaveOccurred())
	testUserDir = testRootDir + "/testUserDir"
	err = os.Mkdir(testUserDir, os.ModeDir)
	Expect(err).NotTo(HaveOccurred())
	testWorkspaceDir = testRootDir + "/testWorkspace"
	err = os.Mkdir(testWorkspaceDir, os.ModeDir)
	Expect(err).NotTo(HaveOccurred())

	// Common infrastructure
	crleditor.CrlEditorSingleton = crleditor.NewEditor(testWorkspaceDir)
	err = crleditor.CrlEditorSingleton.Initialize(testWorkspaceDir, true)
	if err != nil {
		log.Fatal(err)
	}

	fyneGUI := &CrlEditorFyneGUI{}
	fyneGUI.app = test.NewApp()
	initializeFyneGUI(fyneGUI, crleditor.CrlEditorSingleton)

	err = crleditor.CrlEditorSingleton.AddEditorGUI(FyneGUISingleton)
	if err != nil {
		log.Fatal(err)
	}
	initialSize := fyne.NewSize(1600.0, 900.0)
	FyneGUISingleton.GetWindow().Resize(initialSize)
	FyneGUISingleton.GetWindow().ShowAndRun()

})

func TestCrlEditor(t *testing.T) {
	testT = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "CrlEditorFyneGUI Suite")
}
