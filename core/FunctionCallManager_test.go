package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FunctionCall Queue Tests", func() {
	var uOfD *UniverseOfDiscourse
	var hl *Transaction
	var queue *pendingFunctionCallQueue
	var element Element
	BeforeEach(func() {
		queue = newPendingFunctionCallQueue()
		uOfD = NewUniverseOfDiscourse()
		hl = uOfD.NewHeldLocks()
		element, _ = uOfD.NewElement(hl)
	})
	Specify("An empty queue should handle calls gracefully", func() {
		Expect(queue.dequeue()).To(BeNil())
		Expect(queue.enqueue(nil)).ToNot(Succeed())
	})
	Specify("Enque on an empty queue should succeed", func() {
		entry, err := newPendingFunctionCall("", nil, element, nil)
		Expect(err).To(BeNil())
		Expect(queue.enqueue(entry)).To(Succeed())
		Expect(queue.queueHead.pendingCall).To(Equal(entry))
		Expect(queue.queueHead).To(Equal(queue.queueTail))
	})
	Specify("Dequeue on a single-entry queue should leave an empty queue", func() {
		entry, err := newPendingFunctionCall("", nil, element, nil)
		Expect(err).To(BeNil())
		Expect(queue.enqueue(entry)).To(Succeed())
		Expect(queue.dequeue()).To(Equal(entry))
		Expect(queue.queueHead).To(BeNil())
		Expect(queue.queueTail).To(BeNil())
	})
	Specify("Enqueue on a single-entry queue should leave a two-entry queue", func() {
		entry1, err1 := newPendingFunctionCall("", nil, element, nil)
		Expect(err1).To(BeNil())
		entry2, err2 := newPendingFunctionCall("", nil, element, nil)
		Expect(err2).To(BeNil())
		Expect(queue.enqueue(entry1)).To(Succeed())
		Expect(queue.enqueue(entry2)).To(Succeed())
		Expect(queue.queueHead.pendingCall).To(Equal(entry1))
		Expect(queue.queueTail.pendingCall).To(Equal(entry2))
		Expect(queue.queueHead.next).To(Equal(queue.queueTail))
	})
	Specify("Enqueue on a two-entry queue should leave a single-entry queue", func() {
		entry1, err1 := newPendingFunctionCall("", nil, element, nil)
		Expect(err1).To(BeNil())
		entry2, err2 := newPendingFunctionCall("", nil, element, nil)
		Expect(err2).To(BeNil())
		Expect(queue.enqueue(entry1)).To(Succeed())
		Expect(queue.enqueue(entry2)).To(Succeed())
		Expect(queue.dequeue()).To(Equal(entry1))
		Expect(queue.queueHead.pendingCall).To(Equal(entry2))
		Expect(queue.queueTail).To(Equal(queue.queueHead))
	})
	Specify("findFirstPendingCall on empty queue should return nil", func() {
		Expect(queue.findFirstPendingCall("", nil)).To(BeNil())
	})
	Specify("findFirstPendincCall should find matching entry", func() {
		uOfD := NewUniverseOfDiscourse()
		hl := uOfD.NewHeldLocks()
		el, _ := uOfD.NewElement(hl)
		entry, err := newPendingFunctionCall("ABC", nil, el, nil)
		Expect(err).To(BeNil())
		Expect(queue.enqueue(entry)).To(Succeed())
		foundEntry := queue.findFirstPendingCall("ABC", el)
		Expect(foundEntry).ToNot(BeNil())
		Expect(foundEntry).To(Equal(entry))
	})
})
