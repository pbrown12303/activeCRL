package core

import ()

type ChangeNotification struct {
	changedObject  BaseElement
	natureOfChange NatureOfChange
}

type NatureOfChange int

const (
	Add NatureOfChange = iota
	Modify
	Remove
)
