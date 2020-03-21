package crleditordomain

import (
	"github.com/pbrown12303/activeCRL/core"
	"github.com/pbrown12303/activeCRL/crldatastructures"
	"github.com/pbrown12303/activeCRL/crldiagram"
)

// EditorDomainURI is the URI for the domain of the editor
var EditorDomainURI = "http://activeCrl.com/crlEditor/EditorDomain"

// EditorURI is the URI for accessing the CrlEditor
var EditorURI = "http://activeCrl.com/crlEditor/Editor"

// EditorSettingsURI is the URI for accessing the CrlEditor settings
var EditorSettingsURI = EditorURI + "/Settings"

// EditorOpenDiagramsURI is the URI for accessing the CrlEditor list of open diagrams
var EditorOpenDiagramsURI = EditorSettingsURI + "/OpenDiagrams"

// editorDropReferenceAsLinkURI is the URI for the setting parameter DropReferenceAsLink
var editorDropReferenceAsLinkURI = EditorSettingsURI + "/DropReferenceAsLink"

// editorDropRefinmentAsLinkURI is the URI for the setting parameter DropRefinementAsLink
var editorDropRefinementAsLinkURI = EditorSettingsURI + "/DropRefinementAsLink"

// BuildEditorDomain programmatically constructs the EditorDomain
func BuildEditorDomain(uOfD *core.UniverseOfDiscourse, hl *core.HeldLocks) (core.Element, error) {
	// EditorDomain
	// Assumes that the DiagramConceptSpace has already been added to the uOfD
	domain, _ := uOfD.NewElement(hl, EditorDomainURI)
	domain.SetLabel("EditorDomain", hl)
	domain.SetURI(EditorDomainURI, hl)

	settings, _ := uOfD.NewElement(hl, EditorSettingsURI)
	settings.SetLabel("Settings", hl)
	settings.SetOwningConcept(domain, hl)

	diagram := uOfD.GetElementWithURI(crldiagram.CrlDiagramURI)
	editorOpenDiagramsList, err := crldatastructures.NewList(uOfD, diagram, hl, EditorOpenDiagramsURI)
	if err != nil {
		return nil, err
	}
	//	editorOpenDiagramsList.SetURI(EditorOpenDiagramsURI, hl)
	editorOpenDiagramsList.SetLabel("EditorOpenDiagrams", hl)
	editorOpenDiagramsList.SetOwningConcept(settings, hl)

	domain.SetIsCoreRecursively(hl)

	return domain, nil
}
