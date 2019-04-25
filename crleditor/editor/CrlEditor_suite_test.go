package editor_test

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crleditor/editor"

	//	"github.com/pbrown12303/activeCRL/activeCRL/crlEditor/editor"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var startCrlEditorServerCmd *exec.Cmd
var page *agouti.Page
var agoutiDriver *agouti.WebDriver
var uOfD core.UniverseOfDiscourse

var _ = BeforeSuite(func() {
	// Start the editor server
	go editor.StartServer(false)
	// Start the browser
	// Choose a WebDriver:
	// agoutiDriver = agouti.PhantomJS()
	// agoutiDriver = agouti.Selenium()
	agoutiDriver = agouti.ChromeDriver()
	Expect(agoutiDriver.Start()).To(Succeed())

	var err error
	page, err = agoutiDriver.NewPage(agouti.Browser("chrome"))
	Expect(err).NotTo(HaveOccurred())

	Expect(page.Navigate("http://localhost:8082/index")).To(Succeed())
	Expect(page).To(HaveURL("http://localhost:8082/index/"))
	Eventually(func() bool {
		var initializationComplete interface{}
		page.RunScript("return crlInitializationComplete;", nil, &initializationComplete)
		return initializationComplete.(bool)
	}, 20).Should(BeTrue())
	uOfD = editor.CrlEditorSingleton.GetUofD()
})

var _ = AfterSuite(func() {
	editor.Exit()
	Expect(agoutiDriver.Stop()).To(Succeed())
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
