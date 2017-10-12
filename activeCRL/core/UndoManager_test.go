// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
	"testing"
)

func TestUndoRedoElementCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewElement(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*element) != e1.(*element) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*element) != e1.(*element) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.undoMgr.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*element) != e1.(*element) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*element) != e1.(*element) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*element) != e1.(*element) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoMarkUndoPoint(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	uOfD.NewElement(hl)
	uOfD.MarkUndoPoint()
	if len(uOfD.undoMgr.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry := uOfD.undoMgr.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point change type not Marker")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.undoMgr.redoStack) != 2 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry = uOfD.undoMgr.redoStack.Peek()
	if undoPointEntry.changeType != Creation {
		t.Error("Redo point change type not Creation")
	}
	if undoPointEntry.changedElement == nil {
		t.Error("Undo point changed element is nil")
	}

	// Redo
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element, marking undo point, undo, and before redo")
	//	PrintUndoStack(uOfD.undoMgr.redoStack, "Redo stack after creating new element, marking undo point, undo, and before redo")
	uOfD.Redo(hl)
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element, marking undo point, undo, and redo")
	//	PrintUndoStack(uOfD.undoMgr.redoStack, "Redo stack after creating new element, marking undo point, undo, and redo")
	if len(uOfD.undoMgr.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry = uOfD.undoMgr.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}
}

func TestUndoRedoElementSetName(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewElement(hl)
	uOfD.MarkUndoPoint()
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element and marking undo point, before SetName")

	if len(uOfD.undoMgr.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry := uOfD.undoMgr.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}

	// Verify initial state
	if GetName(e1, hl) != "" {
		t.Error("Name not initially empty string")
	}
	if e1.GetNameLiteralPointer(hl) != nil {
		t.Error("Name literal pointer not initially nil")
	}
	if e1.GetNameLiteral(hl) != nil {
		t.Error("Name literal not initially nil")
	}

	// SetName
	testName := "Test name"
	//	uOfD.debugUndo = true
	SetName(e1, testName, hl)
	//	uOfD.debugUndo = false
	undoStackSizeAfterSetName := len(uOfD.undoMgr.undoStack)
	nameLiteralPointer := e1.GetNameLiteralPointer(hl)
	nameLiteral := e1.GetNameLiteral(hl)
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) == nil {
		t.Error("Name literal pointer not in baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) == nil {
		t.Error("Name literal not in baseElementMap")
	}

	// Undo
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after SetName and before undo")
	//	uOfD.debugUndo = true
	uOfD.Undo(hl)
	//	uOfD.debugUndo = false
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after undo")
	if len(uOfD.undoMgr.undoStack) != 2 {
		t.Error("Undo stack size incorrect after undo of SetName")
	}
	if len(uOfD.undoMgr.redoStack) != (undoStackSizeAfterSetName - 2) {
		t.Error("Redo stack size incorrect after undo of SetName")
	}
	undoPointEntry = uOfD.undoMgr.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}
	if GetName(e1, hl) != "" {
		t.Error("Undo did not remove name")
	}
	if e1.GetNameLiteralPointer(hl) != nil {
		t.Error("Undo did not remove name literal pointer")
	}
	if e1.GetNameLiteral(hl) != nil {
		t.Error("Undo did not remove name literal")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) != nil {
		t.Error("Name literal pointer not removed from baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) != nil {
		t.Error("Name literal not removed from baseElementMap")
	}

	// Redo
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element, marking undo point, settingName, undo, and before redo")
	//	PrintUndoStack(uOfD.undoMgr.redoStack, "Redo stack after creating new element, marking undo point, settingName, undo, and before redo")
	uOfD.Redo(hl)
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element, marking undo point, settingName, undo, and redo")
	//	PrintUndoStack(uOfD.undoMgr.redoStack, "Redo stack after creating new element, marking undo point, settingName, undo, and redo")

	if len(uOfD.undoMgr.redoStack) > 0 {
		redoPointEntry := uOfD.undoMgr.redoStack.Peek()
		if redoPointEntry.changeType != Marker {
			t.Error("redo point changeType not Marker")
		}
	}
	if GetName(e1, hl) != testName {
		t.Error("Redo did not restore name")
	}
	if e1.GetNameLiteralPointer(hl) != nameLiteralPointer {
		t.Error("Redo did not restore name literal pointer")
	}
	if e1.GetNameLiteral(hl) != nameLiteral {
		t.Error("Redo did not restore name literal")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) == nil {
		t.Error("Name literal pointer not restored to baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) == nil {
		t.Error("Name literal not restored to baseElementMap")
	}

	// Now do two undos and two redos
	uOfD.Undo(hl)
	uOfD.Undo(hl)
	uOfD.Redo(hl)
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.redoStack) > 0 {
		redoPointEntry := uOfD.undoMgr.redoStack.Peek()
		if redoPointEntry.changeType != Marker {
			t.Error("redo point changeType not Marker")
		}
	}
	if GetName(e1, hl) != testName {
		t.Error("Double undo/redo did not restore name")
	}
	if e1.GetNameLiteralPointer(hl) != nameLiteralPointer {
		t.Error("Double undo/redo did not restore name literal pointer")
	}
	if e1.GetNameLiteral(hl) != nameLiteral {
		t.Error("Double undo/redo did not restore name literal")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) == nil {
		t.Error("Double undo/redo Name literal pointer not restored to baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) == nil {
		t.Error("Double undo/redoName literal not restored to baseElementMap")
	}
}

func TestUndoRedoElementSetOwner(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	parent := uOfD.NewElement(hl)
	uOfD.MarkUndoPoint()
	child := uOfD.NewElement(hl)
	childId := child.GetId(hl)
	SetOwningElement(child, parent, hl)

	// Undo
	uOfD.Undo(hl)
	if uOfD.GetElement(childId) != nil {
		t.Errorf("Child not removed from uOfD")
	}
	if len(parent.GetOwnedBaseElements(hl)) != 0 {
		t.Errorf("Child not removed from parent")
	}

	// Redo
	uOfD.Redo(hl)
	if uOfD.GetElement(childId) == nil {
		t.Errorf("Child not restored to uOfD")
	}
	if parent.GetOwnedBaseElements(hl)[0].GetId(hl) != childId {
		t.Errorf("Child not restored to parent")
	}

}

func TestUndoRedoReferenceAndReferencedElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	parent := uOfD.NewElement(hl)
	uOfD.MarkUndoPoint()
	child := uOfD.NewElementReference(hl)
	childId := child.GetId(hl)
	SetOwningElement(child, parent, hl)
	target := uOfD.NewElement(hl)
	targetId := target.GetId(hl)
	child.SetReferencedElement(target, hl)

	// Undo
	uOfD.Undo(hl)
	if uOfD.GetElement(childId) != nil {
		t.Errorf("Child not removed from uOfD")
	}
	if len(parent.GetOwnedBaseElements(hl)) != 0 {
		t.Errorf("Child not removed from parent")
	}
	if uOfD.GetElement(targetId) != nil {
		t.Errorf("Target not removed from uOfD")
	}

	// Redo
	uOfD.Redo(hl)
	recoveredChild := uOfD.GetElement(childId)
	if recoveredChild == nil {
		t.Errorf("Child not restored to uOfD, id: %s\n", childId)
	}
	if parent.GetOwnedBaseElements(hl)[0].GetId(hl) != childId {
		t.Errorf("Child not restored to parent")
	}
	recoveredTarget := uOfD.GetElement(targetId)
	if recoveredTarget == nil {
		t.Errorf("Target not restored to uOfD, id: %s\n", targetId)
	} else {
		if recoveredChild.(ElementReference).GetReferencedElement(hl).GetId(hl) != targetId {
			t.Errorf("Child's referencedElement not restored")
		}
	}
}

func TestUndoRedoDeleteReferenceAndReferencedElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	parent := uOfD.NewElement(hl)
	uOfD.MarkUndoPoint()
	child := uOfD.NewElementReference(hl)
	childId := child.GetId(hl)
	SetOwningElement(child, parent, hl)
	target := uOfD.NewElement(hl)
	targetId := target.GetId(hl)
	child.SetReferencedElement(target, hl)
	uOfD.MarkUndoPoint()

	// Delete the child
	uOfD.DeleteBaseElement(child, hl)
	if uOfD.GetElement(childId) != nil {
		t.Errorf("Child not removed from uOfD, id: %s\n", childId)
	}
	if len(parent.GetOwnedBaseElements(hl)) != 0 {
		t.Errorf("Child not removed from parent")
	}

	// Undo the deletion
	uOfD.Undo(hl)
	recoveredChild := uOfD.GetElement(childId)
	if recoveredChild == nil {
		t.Errorf("Child not restored to uOfD, id: %s\n", childId)
	}
	if parent.GetOwnedBaseElements(hl)[0].GetId(hl) != childId {
		t.Errorf("Child not restored to parent")
	}
	recoveredTarget := uOfD.GetElement(targetId)
	if recoveredTarget == nil {
		t.Errorf("Target not restored to uOfD, id: %s\n", targetId)
	} else {
		if recoveredChild.(ElementReference).GetReferencedElement(hl).GetId(hl) != targetId {
			t.Errorf("Child's referencedElement not restored")
		}
	}
}

func TestEmptyStackUndoAndRedo(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.SetRecordingUndo(true)
	uOfD.Undo(nil)
	uOfD.Redo(nil)
}

func TestUndoRedoElementSetUri(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewElement(hl)
	uOfD.MarkUndoPoint()
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element and marking undo point, before SetUri")

	if len(uOfD.undoMgr.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry := uOfD.undoMgr.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}

	// Verify initial state
	if GetUri(e1, hl) != "" {
		t.Error("Uri not initially empty string")
	}
	if e1.GetUriLiteralPointer(hl) != nil {
		t.Error("Uri literal pointer not initially nil")
	}
	if e1.GetUriLiteral(hl) != nil {
		t.Error("Uri literal not initially nil")
	}

	// SetUri
	testUri := "foo.bar"
	//	uOfD.debugUndo = true
	SetUri(e1, testUri, hl)
	//	uOfD.debugUndo = false
	undoStackSizeAfterSetUri := len(uOfD.undoMgr.undoStack)
	nameLiteralPointer := e1.GetUriLiteralPointer(hl)
	nameLiteral := e1.GetUriLiteral(hl)
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) == nil {
		t.Error("Uri literal pointer not in baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) == nil {
		t.Error("Uri literal not in baseElementMap")
	}

	// Undo
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after SetUri and before undo")
	//	uOfD.debugUndo = true
	uOfD.Undo(hl)
	//	uOfD.debugUndo = false
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after undo")
	if len(uOfD.undoMgr.undoStack) != 2 {
		t.Error("Undo stack size incorrect after undo of SetUri")
	}
	if len(uOfD.undoMgr.redoStack) != (undoStackSizeAfterSetUri - 2) {
		t.Error("Redo stack size incorrect after undo of SetUri")
	}
	undoPointEntry = uOfD.undoMgr.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}
	if GetUri(e1, hl) != "" {
		t.Error("Undo did not remove name")
	}
	if e1.GetUriLiteralPointer(hl) != nil {
		t.Error("Undo did not remove name literal pointer")
	}
	if e1.GetUriLiteral(hl) != nil {
		t.Error("Undo did not remove name literal")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) != nil {
		t.Error("Uri literal pointer not removed from baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) != nil {
		t.Error("Uri literal not removed from baseElementMap")
	}

	// Redo
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element, marking undo point, settingUri, undo, and before redo")
	//	PrintUndoStack(uOfD.undoMgr.redoStack, "Redo stack after creating new element, marking undo point, settingUri, undo, and before redo")
	uOfD.Redo(hl)
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element, marking undo point, settingUri, undo, and redo")
	//	PrintUndoStack(uOfD.undoMgr.redoStack, "Redo stack after creating new element, marking undo point, settingUri, undo, and redo")

	if len(uOfD.undoMgr.redoStack) > 0 {
		redoPointEntry := uOfD.undoMgr.redoStack.Peek()
		if redoPointEntry.changeType != Marker {
			t.Error("redo point changeType not Marker")
		}
	}
	if GetUri(e1, hl) != testUri {
		t.Error("Redo did not restore name")
	}
	if e1.GetUriLiteralPointer(hl) != nameLiteralPointer {
		t.Error("Redo did not restore name literal pointer")
	}
	if e1.GetUriLiteral(hl) != nameLiteral {
		t.Error("Redo did not restore name literal")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) == nil {
		t.Error("Uri literal pointer not restored to baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) == nil {
		t.Error("Uri literal not restored to baseElementMap")
	}
}

func TestUndoRedoElementSetDefinition(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewElement(hl)
	uOfD.MarkUndoPoint()
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element and marking undo point, before SetDefinition")

	if len(uOfD.undoMgr.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry := uOfD.undoMgr.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}

	// Verify initial state
	if e1.GetDefinition(hl) != "" {
		t.Error("Definition not initially empty string")
	}
	if e1.GetDefinitionLiteralPointer(hl) != nil {
		t.Error("Definition literal pointer not initially nil")
	}
	if e1.GetDefinitionLiteral(hl) != nil {
		t.Error("Definition literal not initially nil")
	}

	// SetDefinition
	testDefinition := "foo.bar"
	//	uOfD.debugUndo = true
	SetDefinition(e1, testDefinition, hl)
	//	uOfD.debugUndo = false
	undoStackSizeAfterSetDefinition := len(uOfD.undoMgr.undoStack)
	nameLiteralPointer := e1.GetDefinitionLiteralPointer(hl)
	nameLiteral := e1.GetDefinitionLiteral(hl)
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) == nil {
		t.Error("Definition literal pointer not in baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) == nil {
		t.Error("Definition literal not in baseElementMap")
	}

	// Undo
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after SetDefinition and before undo")
	//	uOfD.debugUndo = true
	uOfD.Undo(hl)
	//	uOfD.debugUndo = false
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after undo")
	if len(uOfD.undoMgr.undoStack) != 2 {
		t.Error("Undo stack size incorrect after undo of SetDefinition")
	}
	if len(uOfD.undoMgr.redoStack) != (undoStackSizeAfterSetDefinition - 2) {
		t.Error("Redo stack size incorrect after undo of SetDefinition")
	}
	undoPointEntry = uOfD.undoMgr.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}
	if e1.GetDefinition(hl) != "" {
		t.Error("Undo did not remove name")
	}
	if e1.GetDefinitionLiteralPointer(hl) != nil {
		t.Error("Undo did not remove name literal pointer")
	}
	if e1.GetDefinitionLiteral(hl) != nil {
		t.Error("Undo did not remove name literal")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) != nil {
		t.Error("Definition literal pointer not removed from baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) != nil {
		t.Error("Definition literal not removed from baseElementMap")
	}

	// Redo
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element, marking undo point, settingDefinition, undo, and before redo")
	//	PrintUndoStack(uOfD.undoMgr.redoStack, "Redo stack after creating new element, marking undo point, settingDefinition, undo, and before redo")
	uOfD.Redo(hl)
	//	PrintUndoStack(uOfD.undoMgr.undoStack, "Undo stack after creating new element, marking undo point, settingDefinition, undo, and redo")
	//	PrintUndoStack(uOfD.undoMgr.redoStack, "Redo stack after creating new element, marking undo point, settingDefinition, undo, and redo")

	if len(uOfD.undoMgr.redoStack) > 0 {
		redoPointEntry := uOfD.undoMgr.redoStack.Peek()
		if redoPointEntry.changeType != Marker {
			t.Error("redo point changeType not Marker")
		}
	}
	if e1.GetDefinition(hl) != testDefinition {
		t.Error("Redo did not restore name")
	}
	if e1.GetDefinitionLiteralPointer(hl) != nameLiteralPointer {
		t.Error("Redo did not restore name literal pointer")
	}
	if e1.GetDefinitionLiteral(hl) != nameLiteral {
		t.Error("Redo did not restore name literal")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteralPointer.GetId(hl)) == nil {
		t.Error("Definition literal pointer not restored to baseElementMap")
	}
	if uOfD.baseElementMap.GetEntry(nameLiteral.GetId(hl)) == nil {
		t.Error("Definition literal not restored to baseElementMap")
	}
}

func TestUndoRedoElementPointerCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	ep1 := uOfD.NewReferencedElementPointer(hl)
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(ep1.GetId(hl)).(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(ep1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(ep1.GetId(hl)).(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoElementPointerSetOwningElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewElement(hl)
	r1 := uOfD.NewElementReference(hl)
	uOfD.MarkUndoPoint()
	SetOwningElement(r1, e1, hl)
	oep := r1.GetOwningElementPointer(hl)
	if oep == nil {
		t.Error("Owning element pointer is nil")
	}
	if GetOwningElement(oep, hl).GetId(hl) != r1.GetId(hl) {
		t.Error("Owning element not properly set")
	}
	if oep.GetElement(hl) != e1 {
		t.Error("Elemenet not set properly")
	}
	if e1.IsOwnedBaseElement(r1, hl) == false {
		t.Error("E1.ownedBaseElements does not contain r1")
	}
	if r1.IsOwnedBaseElement(oep, hl) == false {
		t.Error("R1.ownedBaseElements does not contain oep")
	}

	// Undo
	uOfD.Undo(hl)
	if r1.GetOwningElementPointer(hl) != nil {
		t.Error("Owning element pointer is not nil")
	}
	if GetOwningElement(oep, hl) != nil {
		t.Error("Owning element not properly cleared")
	}
	if oep.GetElement(hl) != nil {
		t.Error("Elemenet not cleared properly")
	}
	if e1.IsOwnedBaseElement(r1, hl) {
		t.Error("E1.ownedBaseElements still contains r1")
	}
	if r1.IsOwnedBaseElement(oep, hl) {
		t.Error("R1.ownedBaseElements still contains oep")
	}

	// Redo
	uOfD.Redo(hl)
	if r1.GetOwningElementPointer(hl) == nil {
		t.Error("Owning element pointer is nil")
	}
	if GetOwningElement(oep, hl).GetId(hl) != r1.GetId(hl) {
		t.Error("Owning element not properly set")
	}
	if oep.GetElement(hl) != e1 {
		t.Error("Elemenet not set properly")
	}
	if e1.IsOwnedBaseElement(r1, hl) == false {
		t.Error("E1.ownedBaseElements does not contain r1")
	}
	if r1.IsOwnedBaseElement(oep, hl) == false {
		t.Error("R1.ownedBaseElements does not contain oep")
	}
}

func TestUndoRedoElementPointerPointerCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	epp1 := uOfD.NewElementPointerPointer(hl)
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(epp1.GetId(hl)).(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("ElementPointerPointer not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(epp1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(epp1.GetId(hl)).(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoElementPointerPointerSetElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	ep1 := uOfD.NewReferencedElementPointer(hl)
	r1 := uOfD.NewElementPointerReference(hl)
	uOfD.MarkUndoPoint()
	r1.SetReferencedElementPointer(ep1, hl)
	rep := r1.GetElementPointerPointer(hl)
	if rep == nil {
		t.Error("Referenced element pointer is nil")
	}
	if GetOwningElement(rep, hl).GetId(hl) != r1.GetId(hl) {
		t.Error("Referenced element pointer owner not properly set")
	}
	if rep.GetElementPointer(hl) != ep1 {
		t.Error("Element pointer not set properly")
	}
	if r1.IsOwnedBaseElement(rep, hl) == false {
		t.Error("R1.ownedBaseElements does not contain rep")
	}

	// Undo
	uOfD.Undo(hl)
	if r1.GetElementPointerPointer(hl) != nil {
		t.Error("Referenced element pointer is not nil")
	}
	if GetOwningElement(rep, hl) != nil {
		t.Error("Owning element not properly cleared")
	}
	if rep.GetElementPointer(hl) != nil {
		t.Error("ElementPointer not cleared properly")
	}
	if r1.IsOwnedBaseElement(rep, hl) {
		t.Error("R1.ownedBaseElements still contains rep")
	}

	// Redo
	uOfD.Redo(hl)
	if r1.GetElementPointerPointer(hl) == nil {
		t.Error("Referenced element pointer is nil")
	}
	if GetOwningElement(rep, hl).GetId(hl) != r1.GetId(hl) {
		t.Error("Owning element not properly set")
	}
	if rep.GetElementPointer(hl) != ep1 {
		t.Error("ElementPointer not set properly")
	}
	if r1.IsOwnedBaseElement(rep, hl) == false {
		t.Error("R1.ownedBaseElements does not contain rep")
	}
}

func TestUndoRedoElementPointerReferenceCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewElementPointerReference(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.undoMgr.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoElementReferenceCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewElementReference(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*elementReference) != e1.(*elementReference) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*elementReference) != e1.(*elementReference) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.undoMgr.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*elementReference) != e1.(*elementReference) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*elementReference) != e1.(*elementReference) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*elementReference) != e1.(*elementReference) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoLiteralCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	l1 := uOfD.NewLiteral(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*literal) != l1.(*literal) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(l1.GetId(hl)).(*literal) != l1.(*literal) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.undoMgr.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*literal) != l1.(*literal) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(l1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*literal) != l1.(*literal) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(l1.GetId(hl)).(*literal) != l1.(*literal) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoLiteralPointerCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	lp1 := uOfD.NewNameLiteralPointer(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(lp1.GetId(hl)).(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.undoMgr.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(lp1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(lp1.GetId(hl)).(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoLiteralPointerPointerCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	lp1 := uOfD.NewLiteralPointerPointer(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(lp1.GetId(hl)).(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.undoMgr.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(lp1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(lp1.GetId(hl)).(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoLiteralPointerPointerSetLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	lp1 := uOfD.NewNameLiteralPointer(hl)
	r1 := uOfD.NewLiteralPointerReference(hl)
	uOfD.MarkUndoPoint()
	r1.SetReferencedLiteralPointer(lp1, hl)
	rlp := r1.GetLiteralPointerPointer(hl)
	if rlp == nil {
		t.Error("Referenced element pointer is nil")
	}
	if GetOwningElement(rlp, hl).GetId(hl) != r1.GetId(hl) {
		t.Error("Referenced element pointer owner not properly set")
	}
	if rlp.GetLiteralPointer(hl) != lp1 {
		t.Error("Element pointer not set properly")
	}
	if r1.IsOwnedBaseElement(rlp, hl) == false {
		t.Error("R1.ownedBaseElements does not contain rlp")
	}

	// Undo
	uOfD.Undo(hl)
	if r1.GetLiteralPointerPointer(hl) != nil {
		t.Error("Referenced literal pointer is not nil")
	}
	if GetOwningElement(rlp, hl) != nil {
		t.Error("Owning element not properly cleared")
	}
	if rlp.GetLiteralPointer(hl) != nil {
		t.Error("LiteralPointer not cleared properly")
	}
	if r1.IsOwnedBaseElement(rlp, hl) {
		t.Error("R1.ownedBaseElements still contains rlp")
	}

	// Redo
	uOfD.Redo(hl)
	if r1.GetLiteralPointerPointer(hl) == nil {
		t.Error("Literal pointer pointer is nil")
	}
	if GetOwningElement(rlp, hl).GetId(hl) != r1.GetId(hl) {
		t.Error("Owning element not properly set")
	}
	if rlp.GetLiteralPointer(hl) != lp1 {
		t.Error("LiteralPointer not set properly")
	}
	if r1.IsOwnedBaseElement(rlp, hl) == false {
		t.Error("R1.ownedBaseElements does not contain rlp")
	}
}

func TestUndoRedoLiteralReferenceCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewLiteralReference(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*literalReference) != e1.(*literalReference) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*literalReference) != e1.(*literalReference) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.undoMgr.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*literalReference) != e1.(*literalReference) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*literalReference) != e1.(*literalReference) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*literalReference) != e1.(*literalReference) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoRefinementCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	var wg sync.WaitGroup
	hl := NewHeldLocks(&wg)
	defer hl.ReleaseLocks()
	uOfD.SetRecordingUndo(true)
	e1 := uOfD.NewRefinement(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoMgr.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*refinement) != e1.(*refinement) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*refinement) != e1.(*refinement) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.Undo(hl)
	if len(uOfD.undoMgr.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.undoMgr.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.undoMgr.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*refinement) != e1.(*refinement) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)) != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.Redo(hl)
	if len(uOfD.undoMgr.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.undoMgr.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoMgr.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*refinement) != e1.(*refinement) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap.GetEntry(e1.GetId(hl)).(*refinement) != e1.(*refinement) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}
