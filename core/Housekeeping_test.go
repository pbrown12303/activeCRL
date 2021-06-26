package core

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var df1URI = "http://dummy.function.uri,df1"
var df2URI = "http://dummy.function.uri.df2"
var df3URI = "http://dummy.function.uri.df3"

func dummyChangeFunction(Element, *ChangeNotification, *UniverseOfDiscourse) error {
	return nil
}

var _ = Describe("Verify housekeeping function execution", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction
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
		hl.ReleaseLocksAndWait()
	})

	Describe("Test Element ConceptChanged generation", func() {
		Specify("SetDefinition should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
			definition := "Definition"
			el.SetDefinition(definition, hl)
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
					Expect(pc.notification.GetAfterConceptState().Definition).To(Equal(definition))
					Expect(pc.notification.GetBeforeConceptState().Definition).To(Equal(""))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetLabel should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
			label := "Label"
			el.SetLabel(label, hl)
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
					Expect(pc.notification.GetAfterConceptState().Label).To(Equal(label))
					Expect(pc.notification.GetBeforeConceptState().Label).To(Equal(""))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetOwningConcept should generate an OwningConceptChanged for both the child and the owner", func() {
			el, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
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
			childFound := false
			ownerFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {
					if pc.target == el {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(OwningConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().OwningConceptID).To(Equal(newOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().OwningConceptID).To(Equal(""))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(newOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState()).To(BeNil())
						childFound = true
					} else if pc.target == newOwner {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(OwningConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().OwningConceptID).To(Equal(newOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().OwningConceptID).To(Equal(""))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(newOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState()).To(BeNil())
						ownerFound = true
					}
				}
			}
			Expect(childFound).To(BeTrue())
			Expect(ownerFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetOwningConcept should generate a OwningConceptChanged for the old owner as well", func() {
			el, _ := uOfD.NewElement(hl)
			oldOwner, _ := uOfD.NewElement(hl)
			newOwner, _ := uOfD.NewElement(hl)
			el.SetOwningConceptID(oldOwner.getConceptIDNoLock(), hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
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
			childFound := false
			oldOwnerFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {
					if pc.target == el {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(OwningConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().OwningConceptID).To(Equal(newOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().OwningConceptID).To(Equal(oldOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(newOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState().ConceptID).To(Equal(oldOwner.getConceptIDNoLock()))
						childFound = true
					} else if pc.target == oldOwner {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(OwningConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().OwningConceptID).To(Equal(newOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().OwningConceptID).To(Equal(oldOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(newOwner.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState().ConceptID).To(Equal(oldOwner.getConceptIDNoLock()))
						oldOwnerFound = true
					}
				}
			}
			Expect(childFound).To(BeTrue())
			Expect(oldOwnerFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetReadOnly should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
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
					Expect(strconv.ParseBool(pc.notification.GetAfterConceptState().ReadOnly)).To(BeTrue())
					Expect(strconv.ParseBool(pc.notification.GetBeforeConceptState().ReadOnly)).To(BeFalse())
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetURI should generate a ConceptChanged", func() {
			el, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
			uri := "URI"
			el.SetURI(uri, hl)
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
					Expect(pc.notification.GetAfterConceptState().URI).To(Equal(uri))
					Expect(pc.notification.GetBeforeConceptState().URI).To(Equal(""))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.executedCalls = nil
		})
	})

	Describe("Test Literal ConceptChanged generation", func() {
		Specify("SetLiteralValue should generate ConceptChanged", func() {
			lit, _ := uOfD.NewLiteral(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
			literalValue := "LiteralValue"
			lit.SetLiteralValue(literalValue, hl)
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
					Expect(pc.notification.GetAfterConceptState().LiteralValue).To(Equal(literalValue))
					Expect(pc.notification.GetBeforeConceptState().LiteralValue).To(Equal(""))
					found = true
				}
			}
			Expect(found).To(BeTrue())
			uOfD.executedCalls = nil
		})
	})

	Describe("Test ReferencedConceptChanged generation", func() {
		Specify("SetReferencedConcept should generate ReferencedConceptChanged", func() {
			ref, _ := uOfD.NewReference(hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, hl)
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
			refFound := false
			targetFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {

					if pc.target == ref {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().ReferencedConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().ReferencedConceptID).To(Equal(""))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState()).To(BeNil())
						refFound = true
					} else if pc.target == target {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().ReferencedConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().ReferencedConceptID).To(Equal(""))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState()).To(BeNil())
						targetFound = true
					}
				}
			}
			Expect(refFound).To(BeTrue())
			Expect(targetFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetReferencedConcept should generate ReferencedConceptChanged for old target", func() {
			ref, _ := uOfD.NewReference(hl)
			oldTarget, _ := uOfD.NewElement(hl)
			ref.SetReferencedConcept(oldTarget, NoAttribute, hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
			ref.SetReferencedConceptID(target.getConceptIDNoLock(), NoAttribute, hl)
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
			refFound := false
			oldTargetFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {

					if pc.target == ref {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().ReferencedConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().ReferencedConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState().ConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						refFound = true
					} else if pc.target == oldTarget {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().ReferencedConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().ReferencedConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState().ConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						oldTargetFound = true
					}
				}
			}
			Expect(refFound).To(BeTrue())
			Expect(oldTargetFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
	})

	Describe("Test Refinement AbstractConceptChanged and RefinedConceptChanged generation", func() {
		Specify("SetAbstractConcept should generate AbstractConceptChanged", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
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
			refFound := false
			targetFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {
					if pc.target == ref {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(AbstractConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().AbstractConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().AbstractConceptID).To(Equal(""))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState()).To(BeNil())
						refFound = true
					} else if pc.target == target {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(AbstractConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().AbstractConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().AbstractConceptID).To(Equal(""))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState()).To(BeNil())
						targetFound = true
					}
				}
			}
			Expect(refFound).To(BeTrue())
			Expect(targetFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetAbstractConcept should generate AbstractConceptChanged for old target", func() {
			ref, _ := uOfD.NewRefinement(hl)
			oldTarget, _ := uOfD.NewElement(hl)
			ref.SetAbstractConcept(oldTarget, hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
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
			refFound := false
			oldTargetFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {
					if pc.target == ref {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(AbstractConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().AbstractConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().AbstractConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState().ConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						refFound = true
					} else if pc.target == oldTarget {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(AbstractConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().AbstractConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().AbstractConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState().ConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						oldTargetFound = true
					}
				}
			}
			Expect(refFound).To(BeTrue())
			Expect(oldTargetFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetRefinedConcept should generate RefinedConceptChanged", func() {
			ref, _ := uOfD.NewRefinement(hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
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
			refFound := false
			targetFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {
					if pc.target == ref {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(RefinedConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().RefinedConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().RefinedConceptID).To(Equal(""))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState()).To(BeNil())
						refFound = true
					} else if pc.target == target {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(RefinedConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().RefinedConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().RefinedConceptID).To(Equal(""))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState()).To(BeNil())
						targetFound = true
					}
				}
			}
			Expect(refFound).To(BeTrue())
			Expect(targetFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
		Specify("SetRefinedConcept should generate RefinedConceptChanged for old refined concept", func() {
			ref, _ := uOfD.NewRefinement(hl)
			oldTarget, _ := uOfD.NewElement(hl)
			ref.SetRefinedConcept(oldTarget, hl)
			target, _ := uOfD.NewElement(hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
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
			refFound := false
			oldTargetFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {
					if pc.target == ref {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(RefinedConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().RefinedConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().RefinedConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState().ConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						refFound = true
					} else if pc.target == oldTarget {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(RefinedConceptChanged))
						Expect(pc.notification.GetUnderlyingChange()).To(BeNil())
						Expect(pc.notification.GetAfterConceptState().RefinedConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeConceptState().RefinedConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						Expect(pc.notification.GetAfterReferencedState().ConceptID).To(Equal(target.getConceptIDNoLock()))
						Expect(pc.notification.GetBeforeReferencedState().ConceptID).To(Equal(oldTarget.getConceptIDNoLock()))
						oldTargetFound = true
					}
				}
			}
			Expect(refFound).To(BeTrue())
			Expect(oldTargetFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
	})

	Describe("Test OwningConceptChanged propagation", func() {
		Specify("After SetOwningConcept, OwningConceptChanged should be sent to listeners for both owner and child", func() {
			el, _ := uOfD.NewElement(hl)
			oldOwner, _ := uOfD.NewElement(hl)
			el.SetOwningConcept(oldOwner, hl)
			newOwner, _ := uOfD.NewElement(hl)
			childRef, _ := uOfD.NewReference(hl)
			childRef.SetReferencedConceptID(el.getConceptIDNoLock(), NoAttribute, hl)
			newOwnerRef, _ := uOfD.NewReference(hl)
			newOwnerRef.SetReferencedConceptID(newOwner.getConceptIDNoLock(), NoAttribute, hl)
			oldOwnerRef, _ := uOfD.NewReference(hl)
			oldOwnerRef.SetReferencedConceptID(oldOwner.getConceptIDNoLock(), NoAttribute, hl)
			hl.ReleaseLocksAndWait()
			uOfD.executedCalls = make(chan *pendingFunctionCall, 100)
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
			childRefFound := false
			oldOwnerRefFound := false
			newOwnerRefFound := false
			for _, pc := range calls {
				if pc.functionID == "http://activeCrl.com/core/coreHousekeeping" {
					if pc.target == childRef {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(OwningConceptChanged))
						childRefFound = true
					} else if pc.target == oldOwnerRef {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(OwningConceptChanged))
						oldOwnerRefFound = true
					} else if pc.target == newOwnerRef {
						Expect(pc.notification.GetNatureOfChange()).To(Equal(OwningConceptChanged))
						newOwnerRefFound = true
					}
				}
			}
			Expect(childRefFound).To(BeTrue())
			Expect(oldOwnerRefFound).To(BeTrue())
			Expect(newOwnerRefFound).To(BeTrue())
			uOfD.executedCalls = nil
		})
	})

})
