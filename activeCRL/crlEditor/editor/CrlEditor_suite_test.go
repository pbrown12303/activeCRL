package editor_test

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
)

var startCrlEditorServerCmd *exec.Cmd

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
})

var _ = AfterSuite(func() {
	Expect(agoutiDriver.Stop()).To(Succeed())
	Expect(stop()).To(Succeed())
})
