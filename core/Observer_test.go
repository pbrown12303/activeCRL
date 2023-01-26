package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "strconv"
)

type testObserver struct {
	notifications []*ChangeNotification
}

func (toPtr *testObserver) Update(notification *ChangeNotification, heldLocks *Transaction) error {
	toPtr.notifications = append(toPtr.notifications, notification)
	return nil
}

var _ = Describe("Test Observer functionality", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction
	var obs1 *testObserver
	var obs2 *testObserver
	var obs3 *testObserver
	var obs4 *testObserver

	BeforeEach(func() {
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewTransaction()
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
	Describe("Any type of pointer change should be reported to uOfD observer of the pointer owner", func() {
		Specify("Owner change should be reported for both involved concepts", func() {
			Expect(uOfD.Register(obs1)).To(Succeed())
			el, _ := uOfD.NewElement(hl)
			Expect(el.Register(obs2)).To(Succeed())
			originalOwner, _ := uOfD.NewElement(hl)
			Expect(originalOwner.Register(obs3)).To(Succeed())
			newTarget, _ := uOfD.NewElement(hl)
			Expect(newTarget.Register(obs4)).To(Succeed())
			Expect(el.SetOwningConcept(originalOwner, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(4))
			Expect(obs1.notifications[3].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs1.notifications[3].GetAfterConceptState().OwningConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(1))
			Expect(obs2.notifications[0].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs2.notifications[0].GetAfterConceptState().OwningConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(1))
			Expect(obs3.notifications[0].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs3.notifications[0].GetUnderlyingChange().GetAfterConceptState().OwningConceptID).To(Equal(originalOwner.GetConceptID(hl)))
			// Now the new owner
			Expect(el.SetOwningConcept(newTarget, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(5))
			Expect(obs1.notifications[4].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs1.notifications[4].GetAfterConceptState().OwningConceptID).To(Equal(newTarget.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(2))
			Expect(obs2.notifications[1].GetNatureOfChange()).To(Equal(OwningConceptChanged))
			Expect(obs2.notifications[1].GetAfterConceptState().OwningConceptID).To(Equal(newTarget.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(2))
			Expect(obs3.notifications[1].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs3.notifications[1].GetUnderlyingChange().GetAfterConceptState().OwningConceptID).To(Equal(newTarget.GetConceptID(hl)))
			Expect(len(obs4.notifications)).To(Equal(1))
			Expect(obs4.notifications[0].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs4.notifications[0].GetUnderlyingChange().GetAfterConceptState().OwningConceptID).To(Equal(newTarget.GetConceptID(hl)))
		})
		Specify("Referenced element change should be reported for both reference and its owner", func() {
			Expect(uOfD.Register(obs1)).To(Succeed())
			ref, _ := uOfD.NewReference(hl)
			Expect(ref.Register(obs2)).To(Succeed())
			originalTarget, _ := uOfD.NewElement(hl)
			newTarget, _ := uOfD.NewElement(hl)
			referenceOwner, _ := uOfD.NewElement(hl)
			Expect(referenceOwner.Register(obs3)).To(Succeed())
			Expect(ref.SetOwningConcept(referenceOwner, hl)).To(Succeed())
			Expect(ref.SetReferencedConcept(originalTarget, NoAttribute, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(6))
			Expect(obs1.notifications[5].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs1.notifications[5].GetAfterConceptState().ReferencedConceptID).To(Equal(originalTarget.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(2))
			Expect(obs2.notifications[1].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs2.notifications[1].GetAfterConceptState().ReferencedConceptID).To(Equal(originalTarget.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(2))
			Expect(obs3.notifications[1].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs3.notifications[1].GetUnderlyingChange().GetAfterConceptState().ReferencedConceptID).To(Equal(originalTarget.GetConceptID(hl)))
			// Now the new target
			Expect(ref.SetReferencedConcept(newTarget, NoAttribute, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(7))
			Expect(obs1.notifications[6].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs1.notifications[6].GetAfterConceptState().ReferencedConceptID).To(Equal(newTarget.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(3))
			Expect(obs2.notifications[2].GetNatureOfChange()).To(Equal(ReferencedConceptChanged))
			Expect(obs2.notifications[2].GetAfterConceptState().ReferencedConceptID).To(Equal(newTarget.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(3))
			Expect(obs3.notifications[2].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs3.notifications[2].GetUnderlyingChange().GetAfterConceptState().ReferencedConceptID).To(Equal(newTarget.GetConceptID(hl)))
		})
		Specify("Abstract change should be reported for both refinement and its owner", func() {
			Expect(uOfD.Register(obs1)).To(Succeed())
			ref, _ := uOfD.NewRefinement(hl)
			Expect(ref.Register(obs2)).To(Succeed())
			originalAbstraction, _ := uOfD.NewElement(hl)
			newAbstraction, _ := uOfD.NewElement(hl)
			refinementOwner, _ := uOfD.NewElement(hl)
			Expect(refinementOwner.Register(obs3)).To(Succeed())
			Expect(ref.SetOwningConcept(refinementOwner, hl))
			Expect(ref.SetAbstractConcept(originalAbstraction, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(6))
			Expect(obs1.notifications[5].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs1.notifications[5].GetAfterConceptState().AbstractConceptID).To(Equal(originalAbstraction.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(2))
			Expect(obs2.notifications[1].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs2.notifications[1].GetAfterConceptState().AbstractConceptID).To(Equal(originalAbstraction.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(2))
			Expect(obs3.notifications[1].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs3.notifications[1].GetUnderlyingChange().GetAfterConceptState().AbstractConceptID).To(Equal(originalAbstraction.GetConceptID(hl)))
			// Now the new target
			Expect(ref.SetAbstractConcept(newAbstraction, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(7))
			Expect(obs1.notifications[6].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs1.notifications[6].GetAfterConceptState().AbstractConceptID).To(Equal(newAbstraction.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(3))
			Expect(obs2.notifications[2].GetNatureOfChange()).To(Equal(AbstractConceptChanged))
			Expect(obs2.notifications[2].GetAfterConceptState().AbstractConceptID).To(Equal(newAbstraction.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(3))
			Expect(obs3.notifications[2].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs3.notifications[2].GetUnderlyingChange().GetAfterConceptState().AbstractConceptID).To(Equal(newAbstraction.GetConceptID(hl)))
		})
		Specify("Refined change should be reported for both refinement and its owner", func() {
			Expect(uOfD.Register(obs1)).To(Succeed())
			ref, _ := uOfD.NewRefinement(hl)
			Expect(ref.Register(obs2)).To(Succeed())
			originalTarget, _ := uOfD.NewElement(hl)
			newTarget, _ := uOfD.NewElement(hl)
			refinementOwner, _ := uOfD.NewElement(hl)
			Expect(refinementOwner.Register(obs3)).To(Succeed())
			Expect(ref.SetOwningConcept(refinementOwner, hl)).To(Succeed())
			Expect(ref.SetRefinedConcept(originalTarget, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(6))
			Expect(obs1.notifications[5].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs1.notifications[5].GetAfterConceptState().RefinedConceptID).To(Equal(originalTarget.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(2))
			Expect(obs2.notifications[1].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs2.notifications[1].GetAfterConceptState().RefinedConceptID).To(Equal(originalTarget.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(2))
			Expect(obs3.notifications[1].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs3.notifications[1].GetUnderlyingChange().GetAfterConceptState().RefinedConceptID).To(Equal(originalTarget.GetConceptID(hl)))
			// Now the new owner
			Expect(ref.SetRefinedConcept(newTarget, hl)).To(Succeed())
			Expect(len(obs1.notifications)).To(Equal(7))
			Expect(obs1.notifications[6].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs1.notifications[6].GetAfterConceptState().RefinedConceptID).To(Equal(newTarget.GetConceptID(hl)))
			Expect(len(obs2.notifications)).To(Equal(3))
			Expect(obs2.notifications[2].GetNatureOfChange()).To(Equal(RefinedConceptChanged))
			Expect(obs2.notifications[2].GetAfterConceptState().RefinedConceptID).To(Equal(newTarget.GetConceptID(hl)))
			Expect(len(obs3.notifications)).To(Equal(3))
			Expect(obs3.notifications[2].GetNatureOfChange()).To(Equal(OwnedConceptChanged))
			Expect(obs3.notifications[2].GetUnderlyingChange().GetAfterConceptState().RefinedConceptID).To(Equal(newTarget.GetConceptID(hl)))
		})
	})
})
