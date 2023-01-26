package core

// Observer receives ChangeNotifications when the observed subject changes
type Observer interface {
	Update(notification *ChangeNotification, heldLocks *Transaction) error
}

// Subject notifies Observers when changes occur
type Subject interface {
	Register(observer Observer) error
	Deregister(observer Observer) error
	notifyAll(notification *ChangeNotification, heldLocks *Transaction) error
}
