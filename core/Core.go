package core

//	"time"

// Copyright 2017, 2018 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// CorePrefix is a string that prefixes all URIs for the core concepts
var CorePrefix = "http://activeCrl.com/core/"

// CoreDomainURI is the URI for the core concept space
var CoreDomainURI = CorePrefix + "CoreDomain"

// ElementURI is the URI for the Element core concept
var ElementURI = CorePrefix + "Element"

// LiteralURI is the URI for the Literal core concept
var LiteralURI = CorePrefix + "Literal"

// ReferenceURI is the URI for the Reference core concept
var ReferenceURI = CorePrefix + "Reference"

// RefinementURI is the URI for the Refinement core concept
var RefinementURI = CorePrefix + "Refinement"

// UniverseOfDiscourseURI is the URI for the UniverseOfDiscourse core concept
var UniverseOfDiscourseURI = CorePrefix + "UniverseOfDiscourse"

// AdHocTrace is a global variable used in troubleshooting. Generally debugging logic is wrapped in a
// conditional expression contingent on the value of this variable
var AdHocTrace = false

// initCore is the core initialization, but it is made explicit so that it can be called for testing purposes
func init() {
	TraceChange = false
	notificationsLimit = 0
	// notificationsCount = 0
}

func buildCoreDomain(uOfD *UniverseOfDiscourse, hl *Transaction) Element {
	coreElement, _ := uOfD.NewElement(hl, CoreDomainURI)
	coreElementID := coreElement.getConceptIDNoLock()
	coreElement.SetLabel("CoreDomain", hl)

	// Element
	element, _ := uOfD.NewElement(hl, ElementURI)
	element.SetOwningConceptID(coreElementID, hl)
	element.SetLabel("Element", hl)

	// Literal
	literal, _ := uOfD.NewLiteral(hl, LiteralURI)
	literal.SetOwningConceptID(coreElementID, hl)
	literal.SetLabel("Literal", hl)

	// Reference
	reference, _ := uOfD.NewReference(hl, ReferenceURI)
	reference.SetOwningConceptID(coreElementID, hl)
	reference.SetLabel("Reference", hl)

	// Refinement
	refinement, _ := uOfD.NewRefinement(hl, RefinementURI)
	refinement.SetOwningConceptID(coreElementID, hl)
	refinement.SetLabel("Refinement", hl)

	coreElement.SetIsCoreRecursively(hl)
	return coreElement
}
