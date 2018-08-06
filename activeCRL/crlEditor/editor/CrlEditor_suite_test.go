package editor_test

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//	"github.com/pbrown12303/activeCRL/activeCRL/crlEditor/editor"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var startCrlEditorServerCmd *exec.Cmd
var page *agouti.Page

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

	startCrlEditorServerCmd = exec.Command("crlEditorServer")
	startCrlEditorServerCmd.Dir = "C:/GoWorkspace/bin/"
	RunSpecs(t, "CrlEditor Suite")
}

var agoutiDriver *agouti.WebDriver

var _ = BeforeSuite(func() {
	// Choose a WebDriver:

	// agoutiDriver = agouti.PhantomJS()
	// agoutiDriver = agouti.Selenium()
	agoutiDriver = agouti.ChromeDriver()

	Expect(agoutiDriver.Start()).To(Succeed())

	Expect(startCrlEditorServerCmd.Start()).To(Succeed())

	var err error
	page, err = agoutiDriver.NewPage(agouti.Browser("chrome"))
	Expect(err).NotTo(HaveOccurred())

	Expect(page.Navigate("http://localhost:8080/index")).To(Succeed())
	Eventually(func() bool {
		var isInitialized bool
		isInitialized = false
		page.RunScript("return CrlEditor.IsInitialized()", nil, &isInitialized)
		return isInitialized
	}, 10).Should(BeTrue())
	Expect(page).To(HaveURL("http://localhost:8080/index/"))
})

var _ = AfterSuite(func() {
	Expect(stop()).To(Succeed())
	Expect(page.Destroy()).To(Succeed())
	Expect(agoutiDriver.Stop()).To(Succeed())
})
