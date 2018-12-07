package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var df1URI = "http://dummy.function.uri,df1"
var df2URI = "http://dummy.function.uri.df2"
var df3URI = "http://dummy.function.uri.df3"

func dummyChangeFunction(Element, *ChangeNotification, UniverseOfDiscourse) {
	// noop
}

var _ = Describe("Verify housekeeping function execution", func() {
	var uOfD UniverseOfDiscourse
	var hl *HeldLocks
	var df1 Element
	var df2 Element
	var df3 Element

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
		uOfD.AddFunction(df1URI, dummyChangeFunction)
		uOfD.AddFunction(df2URI, dummyChangeFunction)
		uOfD.AddFunction(df3URI, dummyChangeFunction)
		df1, _ = uOfD.NewElement(hl)
		df1.SetURI(df1URI, hl)
		df2, _ = uOfD.NewElement(hl)
		df2.SetURI(df2URI, hl)
		df3, _ = uOfD.NewElement(hl)
		df3.SetURI(df3URI, hl)
		hl.ReleaseLocksAndWait()
	})

	AfterEach(func() {
		hl.ReleaseLocks()
	})

	Describe("Test Element ConceptChanged generation", func() {
		Specify("SetDefinition should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetDefinition("Definition", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == el {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("SetLabel should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetLabel("Label", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == el {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("SetOwningConcept should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == el {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("SetReadOnly should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetReadOnly(true, hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == el {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("SetURI should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetURI("URI", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == el {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test Literal ConceptChanged generation", func() {
		Specify("SetLiteralValue should generate ConceptChanged", func() {
			lit, _ := uOfD.NewLiteral(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			lit.SetLiteralValue("LiteralValue", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == lit {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test Reference ConceptChanged generation", func() {
		Specify("SetReferencedConcept should generate ConceptChanged", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == ref {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test Refinement ConceptChanged generation", func() {
		Specify("SetAbstractConcept should generate ConceptChanged", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			ref.SetAbstractConceptID(target.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == ref {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("SetRefinedConcept should generate ConceptChanged", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			ref.SetRefinedConceptID(target.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == ref {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test Refinement Abstraction Changed generation", func() {
		Specify("Abstraction changed should be generated when an IndicatedConceptChanged is received from the AbstactConcept", func() {
			ref, _ := uOfD.NewRefinement(hl)
			abstractConcept, _ := uOfD.NewElement(hl)
			refinedConcept, _ := uOfD.NewElement(hl)
			ref.SetAbstractConceptID(abstractConcept.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(refinedConcept.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			abstractConcept.SetLabel("Label", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == refinedConcept {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(AbstractionChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(ConceptChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test ChildAbstractionChange generation", func() {
		Specify("ChildAbstractionChange should be generated when an AbstractionChanged is received by the RefinedConcept", func() {
			ref, _ := uOfD.NewRefinement(hl)
			abstractConcept, _ := uOfD.NewElement(hl)
			refinedConcept, _ := uOfD.NewElement(hl)
			refinedConceptOwner, _ := uOfD.NewElement(hl)
			refinedConcept.SetOwningConceptID(refinedConceptOwner.getConceptIDNoLock(), hl)
			ref.SetAbstractConceptID(abstractConcept.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(refinedConcept.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			abstractConcept.SetLabel("Label", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == refinedConceptOwner {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ChildAbstractionChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(AbstractionChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})

	})

	Describe("Test ConceptChanged propagation", func() {
		Specify("After SetOwningConcept, UofDChanged should be sent to uOfD", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == uOfD {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(UofDConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(ConceptChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("After SetOwningConcept, ChildChanged should be sent to owner", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == newOwner {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ChildChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(ConceptChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("After SetOwningConcept, function associated with element should be invoked", func() {
			el, _ := uOfD.NewElement(hl)
			df1Ref, _ := uOfD.NewRefinement(hl)
			df1Ref.SetAbstractConceptID(df1.getConceptIDNoLock(), hl)
			df1Ref.SetRefinedConceptID(el.getConceptIDNoLock(), hl)
			newOwner, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == df1URI && pc.target == el {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ConceptChanged))
					Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("After SetOwningConcept, IndicatedConceptChanged should be sent to listeners", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			ref, _ := uOfD.NewReference(hl)
			ref.SetReferencedConceptID(el.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == ref {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(ConceptChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test ChildChange propagation", func() {
		Specify("After ChildChanged, another ChildChanged should be sent to owner", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			grandparent, _ := uOfD.NewElement(hl)
			newOwner.SetOwningConceptID(grandparent.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == grandparent {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ChildChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(ChildChanged))
					Expect(pc.notification.GetDepth()).To(Equal(3))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("After ChildChanged, IndicatedConceptChanged should be sent to listeners", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			grandparent, _ := uOfD.NewElement(hl)
			newOwner.SetOwningConceptID(grandparent.getConceptIDNoLock(), hl)
			ref, _ := uOfD.NewReference(hl)
			ref.SetReferencedConceptID(grandparent.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == ref {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(ChildChanged))
					Expect(pc.notification.GetDepth()).To(Equal(4))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test IndicatedConceptChanged propagation", func() {
		Specify("After IndicatedConceptChanged, IndicatedConceptChanged should be sent to listener's owner", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			ref, _ := uOfD.NewReference(hl)
			ref.SetReferencedConceptID(el.getConceptIDNoLock(), hl)
			refOwner, _ := uOfD.NewElement(hl)
			ref.SetOwningConceptID(refOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == refOwner {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetDepth()).To(Equal(3))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("After IndicatedConceptChanged, IndicatedConceptChanged should be sent to listener's grandparent", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			ref, _ := uOfD.NewReference(hl)
			ref.SetReferencedConceptID(el.getConceptIDNoLock(), hl)
			refOwner, _ := uOfD.NewElement(hl)
			ref.SetOwningConceptID(refOwner.getConceptIDNoLock(), hl)
			refGrandparent, _ := uOfD.NewElement(hl)
			refOwner.SetOwningConceptID(refGrandparent.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == refGrandparent {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetDepth()).To(Equal(4))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test AbstractionChanged propagation", func() {
		Specify("When a refinedConcept is also the abstract concept of another refinement, AbstractionChanged is propagated to the other refinement's refined concept", func() {
			ref, _ := uOfD.NewRefinement(hl)
			abstractConcept, _ := uOfD.NewElement(hl)
			refinedConcept, _ := uOfD.NewElement(hl)
			ref.SetAbstractConceptID(abstractConcept.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(refinedConcept.getConceptIDNoLock(), hl)
			refinedConcept2, _ := uOfD.NewElement(hl)
			ref2, _ := uOfD.NewRefinement(hl)
			ref2.SetAbstractConceptID(refinedConcept.getConceptIDNoLock(), hl)
			ref2.SetRefinedConceptID(refinedConcept2.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			abstractConcept.SetLabel("Label", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == refinedConcept2 {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(AbstractionChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(AbstractionChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("AbstractionChanged propagates as an IndicatedElementChanged to other listeners", func() {
			ref, _ := uOfD.NewRefinement(hl)
			abstractConcept, _ := uOfD.NewElement(hl)
			refinedConcept, _ := uOfD.NewElement(hl)
			ref.SetAbstractConceptID(abstractConcept.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(refinedConcept.getConceptIDNoLock(), hl)
			listener, _ := uOfD.NewReference(hl)
			listener.SetReferencedConceptID(refinedConcept.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			abstractConcept.SetLabel("Label", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == listener {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(AbstractionChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("Test ChildAbstractionChanged propagation", func() {
		Specify("ChildAbstractionChanged should propagate as ChildAbstractionChanged to owning concept", func() {
			ref, _ := uOfD.NewRefinement(hl)
			abstractConcept, _ := uOfD.NewElement(hl)
			refinedConcept, _ := uOfD.NewElement(hl)
			refinedConceptOwner, _ := uOfD.NewElement(hl)
			refinedConcept.SetOwningConceptID(refinedConceptOwner.getConceptIDNoLock(), hl)
			ref.SetAbstractConceptID(abstractConcept.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(refinedConcept.getConceptIDNoLock(), hl)
			refinedConceptGrandparent, _ := uOfD.NewElement(hl)
			refinedConceptOwner.SetOwningConceptID(refinedConceptGrandparent.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			abstractConcept.SetLabel("Label", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == refinedConceptGrandparent {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(ChildAbstractionChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(ChildAbstractionChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
		Specify("ChildAbstractionChanged should propagate as AbstractionChanged to refinements of the recipient", func() {
			ref, _ := uOfD.NewRefinement(hl)
			abstractConcept, _ := uOfD.NewElement(hl)
			refinedConcept, _ := uOfD.NewElement(hl)
			refinedConceptOwner, _ := uOfD.NewElement(hl)
			refinedConcept.SetOwningConceptID(refinedConceptOwner.getConceptIDNoLock(), hl)
			ref.SetAbstractConceptID(abstractConcept.getConceptIDNoLock(), hl)
			ref.SetRefinedConceptID(refinedConcept.getConceptIDNoLock(), hl)
			refinedConcept2, _ := uOfD.NewElement(hl)
			ref2, _ := uOfD.NewRefinement(hl)
			ref2.SetAbstractConceptID(refinedConceptOwner.getConceptIDNoLock(), hl)
			ref2.SetRefinedConceptID(refinedConcept2.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			abstractConcept.SetLabel("Label", hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == refinedConcept2 {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(AbstractionChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(ChildAbstractionChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})

	Describe("UofDConceptChanged propagation", func() {
		Specify("UofDConceptChange should propagate as IndicatedConceptChanged to uOfD listeners", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			listener, _ := uOfD.NewReference(hl)
			listener.SetReferencedConceptID(uOfD.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el.SetOwningConceptID(newOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == listener {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(UofDConceptChanged))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})
	Describe("UofDConceptAdded propagation", func() {
		Specify("UofDConceptAdded should propagate as IndicatedConceptChanged to uOfD listeners", func() {
			listener, _ := uOfD.NewReference(hl)
			listener.SetReferencedConceptID(uOfD.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == listener {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(UofDConceptAdded))
					Expect(pc.notification.GetUnderlyingChange().GetPriorConceptState().getConceptIDNoLock()).To(Equal(el.getConceptIDNoLock()))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})
	Describe("UofDConceptRemoved propagation", func() {
		Specify("UofDConceptRemoved should propagate as IndicatedConceptChanged to uOfD listeners", func() {
			listener, _ := uOfD.NewReference(hl)
			listener.SetReferencedConceptID(uOfD.getConceptIDNoLock(), hl)
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.(*universeOfDiscourse).executedCalls = make(chan *pendingFunctionCall, 100)
			uOfD.ClearUniverseOfDiscourse(el, hl)
			hl.ReleaseLocksAndWait()
			var calls []*pendingFunctionCall
			done := false
			for done == false {
				select {
				case pc := <-uOfD.getExecutedCalls():
					calls = append(calls, pc)
				default:
					done = true
				}
			}
			found := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" && pc.target == listener {
					Expect(pc.notification.GetNatureOfChange()).To(Equal(IndicatedConceptChanged))
					Expect(pc.notification.GetUnderlyingChange().GetNatureOfChange()).To(Equal(UofDConceptRemoved))
					Expect(pc.notification.GetUnderlyingChange().GetPriorConceptState().getConceptIDNoLock()).To(Equal(el.getConceptIDNoLock()))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.(*universeOfDiscourse).executedCalls = nil
		})
	})
})
