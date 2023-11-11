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

// TransientURI is the URI for transient information that needs to have undo/redo
// It is never saved and does not appear in the UniverseOfDiscourse.GetRootElements() set
var TransientURI = CorePrefix + "Transient"

// Transient is the instantiated Transient element
var Transient *Concept

// AdHocTrace is a global variable used in troubleshooting. Generally debugging logic is wrapped in a
// conditional expression contingent on the value of this variable
var AdHocTrace = false

// initCore is the core initialization, but it is made explicit so that it can be called for testing purposes
func init() {
	TraceChange = false
	notificationsLimit = 0
	// notificationsCount = 0
}

func buildCoreDomain(uOfD *UniverseOfDiscourse, trans *Transaction) *Concept {
	coreElement, _ := uOfD.NewElement(trans, CoreDomainURI)
	coreElementID := coreElement.getConceptIDNoLock()
	coreElement.SetLabel("CoreDomain", trans)

	// Element
	element, _ := uOfD.NewElement(trans, ElementURI)
	element.SetOwningConceptID(coreElementID, trans)
	element.SetLabel("Element", trans)

	// Literal
	literal, _ := uOfD.NewLiteral(trans, LiteralURI)
	literal.SetOwningConceptID(coreElementID, trans)
	literal.SetLabel("Literal", trans)

	// Reference
	reference, _ := uOfD.NewReference(trans, ReferenceURI)
	reference.SetOwningConceptID(coreElementID, trans)
	reference.SetLabel("Reference", trans)

	// Refinement
	refinement, _ := uOfD.NewRefinement(trans, RefinementURI)
	refinement.SetOwningConceptID(coreElementID, trans)
	refinement.SetLabel("Refinement", trans)

	coreElement.SetIsCoreRecursively(trans)

	Transient, _ = uOfD.NewElement(trans, TransientURI)

	return coreElement
}
