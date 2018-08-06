package coreFunctions

import (
	"github.com/pbrown12303/activeCRL/activeCRL/core"
	"github.com/satori/go.uuid"
	"testing"
)

func validateBaseElementReferenceId(t *testing.T, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks, uri string) {
	expectedId := uuid.NewV5(uuid.NamespaceURL, uri).String()
	representation := uOfD.GetBaseElementReferenceWithUri(uri)
	if representation == nil {
		t.Errorf("Representation not found for uri %s\n", uri)
		be := uOfD.GetBaseElementWithUri(uri)
		core.Print(be, "Found BE: ", hl)
	} else {
		if expectedId != representation.GetId(hl) {
			t.Errorf("ID incorrectly set for uri %s\n", uri)
		}
	}
}

func validateElementId(t *testing.T, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks, uri string) {
	expectedId := uuid.NewV5(uuid.NamespaceURL, uri).String()
	representation := uOfD.GetElementWithUri(uri)
	if representation == nil {
		t.Errorf("Representation not found for uri %s\n", uri)
		be := uOfD.GetBaseElementWithUri(uri)
		core.Print(be, "Found BE: ", hl)
	} else {
		if expectedId != representation.GetId(hl) {
			t.Errorf("ID incorrectly set for uri %s\n", uri)
		}
	}
}

func validateElementPointerReferenceId(t *testing.T, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks, uri string) {
	expectedId := uuid.NewV5(uuid.NamespaceURL, uri).String()
	representation := uOfD.GetElementPointerReferenceWithUri(uri)
	if representation == nil {
		t.Errorf("Representation not found for uri %s\n", uri)
		be := uOfD.GetBaseElementWithUri(uri)
		core.Print(be, "Found BE: ", hl)
	} else {
		if expectedId != representation.GetId(hl) {
			t.Errorf("ID incorrectly set for uri %s\n", uri)
		}
	}
}

func validateElementReferenceId(t *testing.T, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks, uri string) {
	expectedId := uuid.NewV5(uuid.NamespaceURL, uri).String()
	representation := uOfD.GetElementReferenceWithUri(uri)
	if representation == nil {
		t.Errorf("Representation not found for uri %s\n", uri)
		be := uOfD.GetBaseElementWithUri(uri)
		core.Print(be, "Found BE: ", hl)
	} else {
		if expectedId != representation.GetId(hl) {
			t.Errorf("ID incorrectly set for uri %s\n", uri)
		}
	}
}

func validateLiteralReferenceId(t *testing.T, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks, uri string) {
	expectedId := uuid.NewV5(uuid.NamespaceURL, uri).String()
	representation := uOfD.GetLiteralReferenceWithUri(uri)
	if representation == nil {
		t.Errorf("Representation not found for uri %s\n", uri)
		be := uOfD.GetBaseElementWithUri(uri)
		core.Print(be, "Found BE: ", hl)
	} else {
		if expectedId != representation.GetId(hl) {
			t.Errorf("ID incorrectly set for uri %s\n", uri)
		}
	}
}

func validateLiteralPointerReferenceId(t *testing.T, uOfD core.UniverseOfDiscourse, hl *core.HeldLocks, uri string) {
	expectedId := uuid.NewV5(uuid.NamespaceURL, uri).String()
	representation := uOfD.GetLiteralPointerReferenceWithUri(uri)
	if representation == nil {
		t.Errorf("Representation not found for uri %s\n", uri)
		be := uOfD.GetBaseElementWithUri(uri)
		core.Print(be, "Found BE: ", hl)
	} else {
		if expectedId != representation.GetId(hl) {
			t.Errorf("ID incorrectly set for uri %s\n", uri)
		}
	}
}
