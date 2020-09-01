package browsergui_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"testing"
	// "time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"

	"github.com/pbrown12303/activeCRL/crleditor"
	"github.com/pbrown12303/activeCRL/crleditorbrowsergui/browsergui"
)

var startCrlEditorServerCmd *exec.Cmd
var page *agouti.Page
var agoutiDriver *agouti.WebDriver
var testRootDir string
var testUserDir string
var testWorkspaceDir string
var testEditor *crleditor.Editor

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

	testRootDir, err = ioutil.TempDir(tempDirPath, "crlEditorTestDir*")
	Expect(err).NotTo(HaveOccurred())
	testUserDir = testRootDir + "/testUserDir"
	err = os.Mkdir(testUserDir, os.ModeDir)
	Expect(err).NotTo(HaveOccurred())
	testWorkspaceDir = testRootDir + "/testWorkspace"
	err = os.Mkdir(testWorkspaceDir, os.ModeDir)
	Expect(err).NotTo(HaveOccurred())

	// Start the testEditor and GUI
	// Common infrastructure
	testEditor = crleditor.NewEditor(testUserDir)
	browsergui.InitializeBrowserGUISingleton(testEditor, false)
	err = testEditor.AddEditorGUI(browsergui.BrowserGUISingleton)
	if err != nil {
		log.Fatal(err)
	}
	// fyneGUI := fynegui.NewFyneGUI()
	// err = testEditor.AddEditorGUI(fyneGUI)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go testEditor.Initialize(testWorkspaceDir, false)

	// Start the browser
	// Choose a WebDriver:
	// agoutiDriver = agouti.PhantomJS()
	// agoutiDriver = agouti.Selenium()
	agoutiDriver = agouti.ChromeDriver()
	Expect(agoutiDriver.Start()).To(Succeed())

	page, err = agoutiDriver.NewPage(agouti.Browser("chrome"))
	Expect(err).NotTo(HaveOccurred())

	Expect(page.Navigate("http://localhost:8082/index")).To(Succeed())
	Expect(page).To(HaveURL("http://localhost:8082/index/"))
	Eventually(func() bool {
		var initializationComplete bool
		page.RunScript("return crlInitializationComplete;", nil, &initializationComplete)
		return initializationComplete
	}, 20).Should(BeTrue())
	// var fileMenuButton = page.FindByID("FileMenuButton")
	// Expect(fileMenuButton.Click()).To(Succeed())
	// var clearWorkspaceButton = page.FindByID("ClearWorkspaceButton")
	// Expect(clearWorkspaceButton.Click()).To(Succeed())
	// Eventually(func() bool {
	// 	return browsergui.GetRequestInProgress() == false
	// }, time.Second*10).Should(BeTrue())
	// Suppress alerts since they require user response
	Expect(page.RunScript("crlAutomatedTestInProgress = true;", nil, nil)).To(Succeed())
})

var _ = AfterSuite(func() {
	var exitButton = page.FindByID("Exit")
	Expect(exitButton.Click()).To(Succeed())
	Expect(agoutiDriver.Stop()).To(Succeed())
	Expect(os.RemoveAll(testRootDir)).To(Succeed())
})

func stop() error {
	var err error
	if runtime.GOOS == "windows" {
		err = startCrlEditorServerCmd.Process.Kill()
	} else {
		err = startCrlEditorServerCmd.Process.Signal(syscall.SIGTERM)
	}
	if err != nil {
		return fmt.Errorf("failed to stop command: %s", err)
	}

	startCrlEditorServerCmd.Wait()

	return nil
}

func TestCrlEditor(t *testing.T) {
	RegisterFailHandler(Fail)

	// startCrlEditorServerCmd = exec.Command("crlEditorServer")
	// startCrlEditorServerCmd.Dir = "C:/GoWorkspace/bin/"
	RunSpecs(t, "CrlEditor Suite")
}
