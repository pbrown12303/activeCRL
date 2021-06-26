package crleditordomain

import (
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructuresdomain"
)

// EditorDomainURI is the URI for the domain of the editor
var EditorDomainURI = "http://activeCrl.com/crleditordomain/CrlEditorDomain"

// EditorSettingsURI is the URI for accessing the CrlEditor settings
var EditorSettingsURI = EditorDomainURI + "/Settings"

// EditorOpenDiagramsURI is the URI for accessing the CrlEditor list of open diagrams
var EditorOpenDiagramsURI = EditorSettingsURI + "/OpenDiagrams"

// BuildEditorDomain programmatically constructs the EditorDomain
func BuildEditorDomain(uOfD *core.UniverseOfDiscourse, hl *core.Transaction) (core.Element, error) {
	// EditorDomain
	// Assumes that the DiagramDomain has already been added to the uOfD
	domain, _ := uOfD.NewElement(hl, EditorDomainURI)
	domain.SetLabel("CrlEditorDomain", hl)

	settings, _ := uOfD.NewElement(hl, EditorSettingsURI)
	settings.SetLabel("Settings", hl)
	settings.SetOwningConcept(domain, hl)

	editorOpenDiagramsList, err := crldatastructuresdomain.NewStringList(uOfD, hl, EditorOpenDiagramsURI)
	if err != nil {
		return nil, err
	}
	//	editorOpenDiagramsList.SetURI(EditorOpenDiagramsURI, hl)
	editorOpenDiagramsList.SetLabel("EditorOpenDiagrams", hl)
	editorOpenDiagramsList.SetOwningConcept(settings, hl)

	domain.SetIsCoreRecursively(hl)

	return domain, nil
}
