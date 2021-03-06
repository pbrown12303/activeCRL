package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "strconv"
)

type testObserver struct {
	notifications []*ChangeNotification
}

func (toPtr *testObserver) Update(notification *ChangeNotification, heldLocks *HeldLocks) error {
	toPtr.notifications = append(toPtr.notifications, notification)
	return nil
}

var _ = Describe("Test Observer functionality", func() {
	var uOfD *UniverseOfDiscourse
	var hl *HeldLocks
	var obs1 *testObserver
	var obs2 *testObserver
	var obs3 *testObserver
	var obs4 *testObserver

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
		obs1 = &testObserver{}
		obs2 = &testObserver{}
		obs3 = &testObserver{}
		obs4 = &testObserver{}
	})
	Specify("Element added to uOfD should be reported to uOfD observer", func() {
		Expect(uOfD.Register(obs1)).To(Succeed())
		el, err := uOfD.NewElement(hl)
		Expect(err).To(BeNil())
		Expect(len(obs1.notifications)).To(Equal(1))
		Expect(obs1.notifications[0].GetNatureOfChange()).To(Equal(ConceptAdded))
		Expect(obs1.notifications[0].GetAfterConceptState().ConceptID).To(Equal(el.GetConceptID(hl)))
	})
	Specify("Element removed from uOfD should be reported to uOfD observer", func() {
		Expect(uOfD.Register(obs1)).To(Succeed())
		el, err := uOfD.NewElement(hl)
		uOfD.DeleteElement(el, hl)
		Expect(err).To(BeNil())
		Expect(len(obs1.notifications)).To(Equal(2))
		Expect(obs1.notifications[1].GetNatureOfChange()).To(Equal(ConceptRemoved))
		Expect(obs1.notifications[1].GetBeforeConceptState().ConceptID).To(Equal(el.GetConceptID(hl)))
	})
	Specify("Element changed should be reported to uOfD observer and concept observer", func() {
		Expect(uOfD.Register(obs1)).To(Succeed())
		el, err := uOfD.NewElement(hl)
		Expect(el.Register(obs2)).To(Succeed())
		el.SetLabel("TestLabel", hl)
		Expect(err).To(BeNil())
		Expect(len(obs1.notifications)).To(Equal(2))
		Expect(obs1.notifications[1].GetNatureOfChange()).To(Equal(ConceptChanged))
		Expect(obs1.notifications[1].GetAfterConceptState().Label).To(Equal(el.GetLabel(hl)))
		Expect(len(obs2.notifications)).To(Equal(1))
		Expect(obs2.notifications[0].GetNatureOfChange()).To(Equal(ConceptChanged))
		Expect(obs2.notifications[0].GetAfterConceptState().Label).To(Equal(el.GetLabel(hl)))
	})
	Describe("Any type of pointer change should be reported to uOfD observer", func() {
		Specify("Owner change should be reported for both involved concepts", func() {
			Expect(uOfD.Register(obs1)).To(Succeed())
			el, _ := uOfD.NewElement(hl)
			Expect(el.Register(obs2)).To(Succeed())
			originalOwner, _ := uOfD.NewElement(hl)
			Expect(originalOwner.Register(obs3)).To(Succeed())
			newOwner, _ := uOfD.NewElement(hl)
			Expect(newOwner.Register(obs4)).To(Succeed())
			Expect(el.SetOwningConcept(originalOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(4))
			Expect(obs1.notifications[3].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs1.notifications[3].GetAfterConceptState().OwningConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(1))
			Expect(obs2.notifications[0].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs2.notifications[0].GetAfterConceptState().OwningConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(1))
			Expect(obs3.notifications[0].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs3.notifications[0].GetAfterConceptState().OwningConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(obs3.notifications[0].GetAfterReferencedState().ConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			// Now the new owner
			Expect(el.SetOwningConcept(newOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(5))
			Expect(obs1.notifications[4].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs1.notifications[4].GetAfterConceptState().OwningConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(2))
			Expect(obs2.notifications[1].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs2.notifications[1].GetAfterConceptState().OwningConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(2))
			Expect(obs3.notifications[1].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs3.notifications[1].GetAfterConceptState().OwningConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(obs3.notifications[1].GetAfterReferencedState().ConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs4.notifications)).To(Equal(1))
			Expect(obs4.notifications[0].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs4.notifications[0].GetAfterConceptState().OwningConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(obs4.notifications[0].GetAfterReferencedState().ConceptID).To(Equal(newOwner.GetConceptID(hl)))
		})
		Specify("Referenced element change should be reported for both involved concepts", func() {
			Expect(uOfD.Register(obs1)).To(Succeed())
			ref, _ := uOfD.NewReference(hl)
			Expect(ref.Register(obs2)).To(Succeed())
			originalOwner, _ := uOfD.NewElement(hl)
			Expect(originalOwner.Register(obs3)).To(Succeed())
			newOwner, _ := uOfD.NewElement(hl)
			Expect(newOwner.Register(obs4)).To(Succeed())
			Expect(ref.SetReferencedConcept(originalOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(4))
			Expect(obs1.notifications[3].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs1.notifications[3].GetAfterConceptState().ReferencedConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(1))
			Expect(obs2.notifications[0].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs2.notifications[0].GetAfterConceptState().ReferencedConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(1))
			Expect(obs3.notifications[0].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs3.notifications[0].GetAfterConceptState().ReferencedConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(obs3.notifications[0].GetAfterReferencedState().ConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			// Now the new owner
			Expect(ref.SetReferencedConcept(newOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(5))
			Expect(obs1.notifications[4].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs1.notifications[4].GetAfterConceptState().ReferencedConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(2))
			Expect(obs2.notifications[1].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs2.notifications[1].GetAfterConceptState().ReferencedConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(2))
			Expect(obs3.notifications[1].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs3.notifications[1].GetAfterConceptState().ReferencedConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(obs3.notifications[1].GetAfterReferencedState().ConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs4.notifications)).To(Equal(1))
			Expect(obs4.notifications[0].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs4.notifications[0].GetAfterConceptState().ReferencedConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(obs4.notifications[0].GetAfterReferencedState().ConceptID).To(Equal(newOwner.GetConceptID(hl)))
		})
		Specify("Abstract change should be reported for both involved concepts", func() {
			Expect(uOfD.Register(obs1)).To(Succeed())
			ref, _ := uOfD.NewRefinement(hl)
			Expect(ref.Register(obs2)).To(Succeed())
			originalOwner, _ := uOfD.NewElement(hl)
			Expect(originalOwner.Register(obs3)).To(Succeed())
			newOwner, _ := uOfD.NewElement(hl)
			Expect(newOwner.Register(obs4)).To(Succeed())
			Expect(ref.SetAbstractConcept(originalOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(4))
			Expect(obs1.notifications[3].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs1.notifications[3].GetAfterConceptState().AbstractConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(1))
			Expect(obs2.notifications[0].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs2.notifications[0].GetAfterConceptState().AbstractConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(1))
			Expect(obs3.notifications[0].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs3.notifications[0].GetAfterConceptState().AbstractConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(obs3.notifications[0].GetAfterReferencedState().ConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			// Now the new owner
			Expect(ref.SetAbstractConcept(newOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(5))
			Expect(obs1.notifications[4].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs1.notifications[4].GetAfterConceptState().AbstractConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(2))
			Expect(obs2.notifications[1].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs2.notifications[1].GetAfterConceptState().AbstractConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(2))
			Expect(obs3.notifications[1].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs3.notifications[1].GetAfterConceptState().AbstractConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(obs3.notifications[1].GetAfterReferencedState().ConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs4.notifications)).To(Equal(1))
			Expect(obs4.notifications[0].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs4.notifications[0].GetAfterConceptState().AbstractConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(obs4.notifications[0].GetAfterReferencedState().ConceptID).To(Equal(newOwner.GetConceptID(hl)))
		})
		Specify("Refined change should be reported for both involved concepts", func() {
			Expect(uOfD.Register(obs1)).To(Succeed())
			ref, _ := uOfD.NewRefinement(hl)
			Expect(ref.Register(obs2)).To(Succeed())
			originalOwner, _ := uOfD.NewElement(hl)
			Expect(originalOwner.Register(obs3)).To(Succeed())
			newOwner, _ := uOfD.NewElement(hl)
			Expect(newOwner.Register(obs4)).To(Succeed())
			Expect(ref.SetRefinedConcept(originalOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(4))
			Expect(obs1.notifications[3].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs1.notifications[3].GetAfterConceptState().RefinedConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(1))
			Expect(obs2.notifications[0].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs2.notifications[0].GetAfterConceptState().RefinedConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(1))
			Expect(obs3.notifications[0].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs3.notifications[0].GetAfterConceptState().RefinedConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(obs3.notifications[0].GetAfterReferencedState().ConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			// Now the new owner
			Expect(ref.SetRefinedConcept(newOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(5))
			Expect(obs1.notifications[4].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs1.notifications[4].GetAfterConceptState().RefinedConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(2))
			Expect(obs2.notifications[1].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs2.notifications[1].GetAfterConceptState().RefinedConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(2))
			Expect(obs3.notifications[1].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs3.notifications[1].GetAfterConceptState().RefinedConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(obs3.notifications[1].GetAfterReferencedState().ConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(len(obs4.notifications)).To(Equal(1))
			Expect(obs4.notifications[0].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs4.notifications[0].GetAfterConceptState().RefinedConceptID).To(Equal(newOwner.GetConceptID(hl)))
			Expect(obs4.notifications[0].GetAfterReferencedState().ConceptID).To(Equal(newOwner.GetConceptID(hl)))
		})
	})
})
