package editor_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pbrown12303/activeCRL/activeCRL/crlEditor/editor"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

func editorNotNil(page *agouti.Page) bool {
	var editorSingleton *editor.CrlEditor
	page.RunScript("return CrlEditor", nil, &editorSingleton)
	return editorSingleton != nil
}

func editorInitialized(page *agouti.Page) bool {
	//	var theResult interface{}
	//	page.RunScript("return CrlEditor.IsInitialized()", nil, &theResult)
	//	fmt.Printf("The result: %+v", theResult)

	var editorSingleton *editor.CrlEditor
	page.RunScript("return CrlEditor", nil, &editorSingleton)
	fmt.Printf("Editor singleton: %+v", editorSingleton)

	var isInitialized bool
	isInitialized = false
	page.RunScript("return CrlEditor.IsInitialized()", nil, &isInitialized)
	return isInitialized
}

var _ = Describe("CrlEditor Initialization", func() {
	var page *agouti.Page

	BeforeEach(func() {

		var err error
		page, err = agoutiDriver.NewPage(agouti.Browser("chrome"))
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})

	It("it should load the default Universe of Discourse", func() {
		By("Loading the Index page", func() {
			Expect(page.Navigate("http://localhost:8080/index")).To(Succeed())
			Eventually(func() bool {
				var editorSingleton *editor.CrlEditor
				page.RunScript("return CrlEditor", nil, &editorSingleton)
				return editorSingleton != nil
			}, 10).Should(BeTrue())
			Eventually(func() bool {
				var isInitialized bool
				isInitialized = false
				page.RunScript("return CrlEditor.IsInitialized()", nil, &isInitialized)
				return isInitialized
			}, 10).Should(BeTrue())
			Expect(page).To(HaveURL("http://localhost:8080/index/"))
		})
	})
})

var _ = Describe("CrlEditor Selection", func() {
		
	}
