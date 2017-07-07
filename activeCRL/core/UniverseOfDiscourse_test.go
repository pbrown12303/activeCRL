package core

import (
	"testing"
)

func TestGetBaseElementWithUri(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()

	// Element
	element := uOfD.NewElement()
	element.SetName("Element")
	recoveredElement := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/Element")
	if recoveredElement != nil {
		t.Error("Wrong element returned for find Element by URI")
	}
	element.SetUri("http://activeCrl.com/core/Element")
	recoveredElement = uOfD.GetBaseElementWithUri("http://activeCrl.com/core/Element")
	if recoveredElement == nil {
		t.Error("Did not find Element by URI")
	}

	// ElementPointer
	elementPointer := uOfD.NewReferencedElementPointer()
	elementPointer.SetUri("http://activeCrl.com/core/ElementPointer")
	recoveredElementPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/ElementPointer")
	if recoveredElementPointer == nil {
		t.Error("Did not find ElementPointer by URI")
	}

	// ElementPointerPointer
	elementPointerPointer := uOfD.NewElementPointerPointer()
	elementPointerPointer.SetUri("http://activeCrl.com/core/ElementPointerPointer")
	recoveredElementPointerPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/ElementPointerPointer")
	if recoveredElementPointerPointer == nil {
		t.Error("Did not find ElementPointerPointer by URI")
	}

	// ElementPointerReference
	elementPointerReference := uOfD.NewElementPointerReference()
	elementPointerReference.SetName("ElementReference")
	elementPointerReference.SetUri("http://activeCrl.com/core/ElementPointerReference")
	recoveredElementPointerReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/ElementPointerReference")
	if recoveredElementPointerReference == nil {
		t.Error("Did not find ElementPointerReference by URI")
	}

	// ElementReference
	elementReference := uOfD.NewElementReference()
	elementReference.SetName("ElementReference")
	elementReference.SetUri("http://activeCrl.com/core/ElementReference")
	recoveredElementReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/ElementReference")
	if recoveredElementReference == nil {
		t.Error("Did not find ElementReference by URI")
	}

	// Literal
	literal := uOfD.NewLiteral()
	literal.SetUri("http://activeCrl.com/core/Literal")
	recoveredLiteral := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/Literal")
	if recoveredLiteral == nil {
		t.Error("Did not find Literal by URI")
	}

	// LiteralPointer
	literalPointer := uOfD.NewValueLiteralPointer()
	literalPointer.SetUri("http://activeCrl.com/core/LiteralPointer")
	recoveredLiteralPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/LiteralPointer")
	if recoveredLiteralPointer == nil {
		t.Error("Did not find LiteralPointer by URI")
	}

	// LiteralPointerPointer
	literalPointerPointer := uOfD.NewLiteralPointerPointer()
	literalPointerPointer.SetUri("http://activeCrl.com/core/LiteralPointerPointer")
	recoveredLiteralPointerPointer := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/LiteralPointerPointer")
	if recoveredLiteralPointerPointer == nil {
		t.Error("Did not find LiteralPointerPointer by URI")
	}

	// LiteralPointerReference
	literalPointerReference := uOfD.NewLiteralPointerReference()
	literalPointerReference.SetName("LiteralReference")
	literalPointerReference.SetUri("http://activeCrl.com/core/LiteralPointerReference")
	recoveredLiteralPointerReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/LiteralPointerReference")
	if recoveredLiteralPointerReference == nil {
		t.Error("Did not find LiteralPointerReference by URI")
	}

	// LiteralReference
	literalReference := uOfD.NewLiteralReference()
	literalReference.SetName("LiteralReference")
	literalReference.SetUri("http://activeCrl.com/core/LiteralReference")
	recoveredLiteralReference := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/LiteralReference")
	if recoveredLiteralReference == nil {
		t.Error("Did not find LiteralReference by URI")
	}

	// Refinement
	refinement := uOfD.NewRefinement()
	refinement.SetName("Refinement")
	refinement.SetUri("http://activeCrl.com/core/Refinement")
	recoveredRefinement := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/Refinement")
	if recoveredRefinement == nil {
		t.Error("Did not find Refinement by URI")
	}

	// Child of element
	child := uOfD.NewElement()
	child.SetName("Child")
	child.SetOwningElement(element)
	child.SetUri("http://activeCrl.com/core/Element/Child")
	recoveredChild := uOfD.GetBaseElementWithUri("http://activeCrl.com/core/Element/Child")
	if recoveredChild == nil {
		t.Error("Did not find Child by URI")
	}

}

func TestUndoRedoElementCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewElement()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*element) != e1.(*element) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*element) != e1.(*element) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*element) != e1.(*element) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*element) != e1.(*element) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*element) != e1.(*element) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoMarkUndoPoint(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	uOfD.NewElement()
	uOfD.markUndoPoint()
	if len(uOfD.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry := uOfD.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point change type not Marker")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.redoStack) != 2 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry = uOfD.redoStack.Peek()
	if undoPointEntry.changeType != Creation {
		t.Error("Redo point change type not Creation")
	}
	if undoPointEntry.changedElement == nil {
		t.Error("Undo point changed element is nil")
	}

	// Redo
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element, marking undo point, undo, and before redo")
	//	PrintUndoStack(uOfD.redoStack, "Redo stack after creating new element, marking undo point, undo, and before redo")
	uOfD.redo()
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element, marking undo point, undo, and redo")
	//	PrintUndoStack(uOfD.redoStack, "Redo stack after creating new element, marking undo point, undo, and redo")
	if len(uOfD.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry = uOfD.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}
}

func TestUndoRedoElementSetName(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewElement()
	uOfD.markUndoPoint()
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element and marking undo point, before SetName")

	if len(uOfD.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry := uOfD.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}

	// Verify initial state
	if e1.GetName() != "" {
		t.Error("Name not initially empty string")
	}
	if e1.getNameLiteralPointer() != nil {
		t.Error("Name literal pointer not initially nil")
	}
	if e1.getNameLiteral() != nil {
		t.Error("Name literal not initially nil")
	}

	// SetName
	testName := "Test name"
	//	uOfD.debugUndo = true
	e1.SetName(testName)
	//	uOfD.debugUndo = false
	undoStackSizeAfterSetName := len(uOfD.undoStack)
	nameLiteralPointer := e1.getNameLiteralPointer()
	nameLiteral := e1.getNameLiteral()
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] == nil {
		t.Error("Name literal pointer not in baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] == nil {
		t.Error("Name literal not in baseElementMap")
	}

	// Undo
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after SetName and before undo")
	//	uOfD.debugUndo = true
	uOfD.undo()
	//	uOfD.debugUndo = false
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after undo")
	if len(uOfD.undoStack) != 2 {
		t.Error("Undo stack size incorrect after undo of SetName")
	}
	if len(uOfD.redoStack) != (undoStackSizeAfterSetName - 2) {
		t.Error("Redo stack size incorrect after undo of SetName")
	}
	undoPointEntry = uOfD.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}
	if e1.GetName() != "" {
		t.Error("Undo did not remove name")
	}
	if e1.getNameLiteralPointer() != nil {
		t.Error("Undo did not remove name literal pointer")
	}
	if e1.getNameLiteral() != nil {
		t.Error("Undo did not remove name literal")
	}
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] != nil {
		t.Error("Name literal pointer not removed from baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] != nil {
		t.Error("Name literal not removed from baseElementMap")
	}

	// Redo
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element, marking undo point, settingName, undo, and before redo")
	//	PrintUndoStack(uOfD.redoStack, "Redo stack after creating new element, marking undo point, settingName, undo, and before redo")
	uOfD.redo()
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element, marking undo point, settingName, undo, and redo")
	//	PrintUndoStack(uOfD.redoStack, "Redo stack after creating new element, marking undo point, settingName, undo, and redo")

	if len(uOfD.redoStack) > 0 {
		redoPointEntry := uOfD.redoStack.Peek()
		if redoPointEntry.changeType != Marker {
			t.Error("redo point changeType not Marker")
		}
	}
	if e1.GetName() != testName {
		t.Error("Redo did not restore name")
	}
	if e1.getNameLiteralPointer() != nameLiteralPointer {
		t.Error("Redo did not restore name literal pointer")
	}
	if e1.getNameLiteral() != nameLiteral {
		t.Error("Redo did not restore name literal")
	}
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] == nil {
		t.Error("Name literal pointer not restored to baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] == nil {
		t.Error("Name literal not restored to baseElementMap")
	}

	// Now do two undos and two redos
	uOfD.undo()
	uOfD.undo()
	uOfD.redo()
	uOfD.redo()
	if len(uOfD.redoStack) > 0 {
		redoPointEntry := uOfD.redoStack.Peek()
		if redoPointEntry.changeType != Marker {
			t.Error("redo point changeType not Marker")
		}
	}
	if e1.GetName() != testName {
		t.Error("Double undo/redo did not restore name")
	}
	if e1.getNameLiteralPointer() != nameLiteralPointer {
		t.Error("Double undo/redo did not restore name literal pointer")
	}
	if e1.getNameLiteral() != nameLiteral {
		t.Error("Double undo/redo did not restore name literal")
	}
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] == nil {
		t.Error("Double undo/redo Name literal pointer not restored to baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] == nil {
		t.Error("Double undo/redoName literal not restored to baseElementMap")
	}
}

func TestEmptyStackUndoAndRedo(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	uOfD.undo()
	uOfD.redo()

}

func TestUndoRedoElementSetUri(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewElement()
	uOfD.markUndoPoint()
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element and marking undo point, before SetUri")

	if len(uOfD.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry := uOfD.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}

	// Verify initial state
	if e1.GetUri() != "" {
		t.Error("Uri not initially empty string")
	}
	if e1.getUriLiteralPointer() != nil {
		t.Error("Uri literal pointer not initially nil")
	}
	if e1.getUriLiteral() != nil {
		t.Error("Uri literal not initially nil")
	}

	// SetUri
	testUri := "foo.bar"
	//	uOfD.debugUndo = true
	e1.SetUri(testUri)
	//	uOfD.debugUndo = false
	undoStackSizeAfterSetUri := len(uOfD.undoStack)
	nameLiteralPointer := e1.getUriLiteralPointer()
	nameLiteral := e1.getUriLiteral()
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] == nil {
		t.Error("Uri literal pointer not in baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] == nil {
		t.Error("Uri literal not in baseElementMap")
	}

	// Undo
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after SetUri and before undo")
	//	uOfD.debugUndo = true
	uOfD.undo()
	//	uOfD.debugUndo = false
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after undo")
	if len(uOfD.undoStack) != 2 {
		t.Error("Undo stack size incorrect after undo of SetUri")
	}
	if len(uOfD.redoStack) != (undoStackSizeAfterSetUri - 2) {
		t.Error("Redo stack size incorrect after undo of SetUri")
	}
	undoPointEntry = uOfD.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}
	if e1.GetUri() != "" {
		t.Error("Undo did not remove name")
	}
	if e1.getUriLiteralPointer() != nil {
		t.Error("Undo did not remove name literal pointer")
	}
	if e1.getUriLiteral() != nil {
		t.Error("Undo did not remove name literal")
	}
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] != nil {
		t.Error("Uri literal pointer not removed from baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] != nil {
		t.Error("Uri literal not removed from baseElementMap")
	}

	// Redo
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element, marking undo point, settingUri, undo, and before redo")
	//	PrintUndoStack(uOfD.redoStack, "Redo stack after creating new element, marking undo point, settingUri, undo, and before redo")
	uOfD.redo()
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element, marking undo point, settingUri, undo, and redo")
	//	PrintUndoStack(uOfD.redoStack, "Redo stack after creating new element, marking undo point, settingUri, undo, and redo")

	if len(uOfD.redoStack) > 0 {
		redoPointEntry := uOfD.redoStack.Peek()
		if redoPointEntry.changeType != Marker {
			t.Error("redo point changeType not Marker")
		}
	}
	if e1.GetUri() != testUri {
		t.Error("Redo did not restore name")
	}
	if e1.getUriLiteralPointer() != nameLiteralPointer {
		t.Error("Redo did not restore name literal pointer")
	}
	if e1.getUriLiteral() != nameLiteral {
		t.Error("Redo did not restore name literal")
	}
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] == nil {
		t.Error("Uri literal pointer not restored to baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] == nil {
		t.Error("Uri literal not restored to baseElementMap")
	}
}

func TestUndoRedoElementSetDefinition(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewElement()
	uOfD.markUndoPoint()
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element and marking undo point, before SetDefinition")

	if len(uOfD.undoStack) != 2 {
		t.Error("Undo stack size incorrect after marking undo point")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after marking undo point")
	}
	undoPointEntry := uOfD.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}

	// Verify initial state
	if e1.GetDefinition() != "" {
		t.Error("Definition not initially empty string")
	}
	if e1.getDefinitionLiteralPointer() != nil {
		t.Error("Definition literal pointer not initially nil")
	}
	if e1.getDefinitionLiteral() != nil {
		t.Error("Definition literal not initially nil")
	}

	// SetDefinition
	testDefinition := "foo.bar"
	//	uOfD.debugUndo = true
	e1.SetDefinition(testDefinition)
	//	uOfD.debugUndo = false
	undoStackSizeAfterSetDefinition := len(uOfD.undoStack)
	nameLiteralPointer := e1.getDefinitionLiteralPointer()
	nameLiteral := e1.getDefinitionLiteral()
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] == nil {
		t.Error("Definition literal pointer not in baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] == nil {
		t.Error("Definition literal not in baseElementMap")
	}

	// Undo
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after SetDefinition and before undo")
	//	uOfD.debugUndo = true
	uOfD.undo()
	//	uOfD.debugUndo = false
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after undo")
	if len(uOfD.undoStack) != 2 {
		t.Error("Undo stack size incorrect after undo of SetDefinition")
	}
	if len(uOfD.redoStack) != (undoStackSizeAfterSetDefinition - 2) {
		t.Error("Redo stack size incorrect after undo of SetDefinition")
	}
	undoPointEntry = uOfD.undoStack.Peek()
	if undoPointEntry.changeType != Marker {
		t.Error("Undo point changeType not Marker")
	}
	if e1.GetDefinition() != "" {
		t.Error("Undo did not remove name")
	}
	if e1.getDefinitionLiteralPointer() != nil {
		t.Error("Undo did not remove name literal pointer")
	}
	if e1.getDefinitionLiteral() != nil {
		t.Error("Undo did not remove name literal")
	}
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] != nil {
		t.Error("Definition literal pointer not removed from baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] != nil {
		t.Error("Definition literal not removed from baseElementMap")
	}

	// Redo
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element, marking undo point, settingDefinition, undo, and before redo")
	//	PrintUndoStack(uOfD.redoStack, "Redo stack after creating new element, marking undo point, settingDefinition, undo, and before redo")
	uOfD.redo()
	//	PrintUndoStack(uOfD.undoStack, "Undo stack after creating new element, marking undo point, settingDefinition, undo, and redo")
	//	PrintUndoStack(uOfD.redoStack, "Redo stack after creating new element, marking undo point, settingDefinition, undo, and redo")

	if len(uOfD.redoStack) > 0 {
		redoPointEntry := uOfD.redoStack.Peek()
		if redoPointEntry.changeType != Marker {
			t.Error("redo point changeType not Marker")
		}
	}
	if e1.GetDefinition() != testDefinition {
		t.Error("Redo did not restore name")
	}
	if e1.getDefinitionLiteralPointer() != nameLiteralPointer {
		t.Error("Redo did not restore name literal pointer")
	}
	if e1.getDefinitionLiteral() != nameLiteral {
		t.Error("Redo did not restore name literal")
	}
	if uOfD.baseElementMap[nameLiteralPointer.getId().String()] == nil {
		t.Error("Definition literal pointer not restored to baseElementMap")
	}
	if uOfD.baseElementMap[nameLiteral.getId().String()] == nil {
		t.Error("Definition literal not restored to baseElementMap")
	}
}

func TestUndoRedoElementPointerCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	ep1 := uOfD.NewReferencedElementPointer()
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[ep1.GetId().String()].(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[ep1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[ep1.GetId().String()].(*elementPointer) != ep1.(*elementPointer) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoElementPointerSetOwningElement(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewElement()
	r1 := uOfD.NewElementReference()
	uOfD.markUndoPoint()
	r1.SetOwningElement(e1)
	oep := r1.getOwningElementPointer()
	if oep == nil {
		t.Error("Owning element pointer is nil")
	}
	if oep.GetOwningElement().GetId() != r1.GetId() {
		t.Error("Owning element not properly set")
	}
	if oep.GetElement() != e1 {
		t.Error("Elemenet not set properly")
	}
	if e1.getOwnedBaseElements()[r1.GetId().String()].GetId() != r1.GetId() {
		t.Error("E1.ownedBaseElements does not contain r1")
	}
	if r1.getOwnedBaseElements()[oep.GetId().String()].GetId() != oep.GetId() {
		t.Error("R1.ownedBaseElements does not contain oep")
	}

	// Undo
	uOfD.undo()
	if r1.getOwningElementPointer() != nil {
		t.Error("Owning element pointer is not nil")
	}
	if oep.GetOwningElement() != nil {
		t.Error("Owning element not properly cleared")
	}
	if oep.GetElement() != nil {
		t.Error("Elemenet not cleared properly")
	}
	if e1.getOwnedBaseElements()[r1.GetId().String()] != nil {
		t.Error("E1.ownedBaseElements still contains r1")
	}
	if r1.getOwnedBaseElements()[oep.GetId().String()] != nil {
		t.Error("R1.ownedBaseElements still contains oep")
	}

	// Redo
	uOfD.redo()
	if r1.getOwningElementPointer() == nil {
		t.Error("Owning element pointer is nil")
	}
	if oep.GetOwningElement().GetId() != r1.GetId() {
		t.Error("Owning element not properly set")
	}
	if oep.GetElement() != e1 {
		t.Error("Elemenet not set properly")
	}
	if e1.getOwnedBaseElements()[r1.GetId().String()].GetId() != r1.GetId() {
		t.Error("E1.ownedBaseElements does not contain r1")
	}
	if r1.getOwnedBaseElements()[oep.GetId().String()].GetId() != oep.GetId() {
		t.Error("R1.ownedBaseElements does not contain oep")
	}
}

func TestUndoRedoElementPointerPointerCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	epp1 := uOfD.NewElementPointerPointer()
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[epp1.GetId().String()].(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("ElementPointerPointer not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[epp1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[epp1.GetId().String()].(*elementPointerPointer) != epp1.(*elementPointerPointer) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoElementPointerPointerSetElementPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	ep1 := uOfD.NewReferencedElementPointer()
	r1 := uOfD.NewElementPointerReference()
	uOfD.markUndoPoint()
	r1.SetElementPointer(ep1)
	rep := r1.getElementPointerPointer()
	if rep == nil {
		t.Error("Referenced element pointer is nil")
	}
	if rep.GetOwningElement().GetId() != r1.GetId() {
		t.Error("Referenced element pointer owner not properly set")
	}
	if rep.GetElementPointer() != ep1 {
		t.Error("Element pointer not set properly")
	}
	if r1.getOwnedBaseElements()[rep.GetId().String()].GetId() != rep.GetId() {
		t.Error("R1.ownedBaseElements does not contain rep")
	}

	// Undo
	uOfD.undo()
	if r1.getElementPointerPointer() != nil {
		t.Error("Referenced element pointer is not nil")
	}
	if rep.GetOwningElement() != nil {
		t.Error("Owning element not properly cleared")
	}
	if rep.GetElementPointer() != nil {
		t.Error("ElementPointer not cleared properly")
	}
	if r1.getOwnedBaseElements()[rep.GetId().String()] != nil {
		t.Error("R1.ownedBaseElements still contains rep")
	}

	// Redo
	uOfD.redo()
	if r1.getElementPointerPointer() == nil {
		t.Error("Referenced element pointer is nil")
	}
	if rep.GetOwningElement().GetId() != r1.GetId() {
		t.Error("Owning element not properly set")
	}
	if rep.GetElementPointer() != ep1 {
		t.Error("ElementPointer not set properly")
	}
	if r1.getOwnedBaseElements()[rep.GetId().String()].GetId() != rep.GetId() {
		t.Error("R1.ownedBaseElements does not contain rep")
	}
}

func TestUndoRedoElementPointerReferenceCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewElementPointerReference()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*elementPointerReference) != e1.(*elementPointerReference) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoElementReferenceCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewElementReference()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*elementReference) != e1.(*elementReference) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*elementReference) != e1.(*elementReference) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*elementReference) != e1.(*elementReference) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*elementReference) != e1.(*elementReference) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*elementReference) != e1.(*elementReference) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoLiteralCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	l1 := uOfD.NewLiteral()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*literal) != l1.(*literal) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[l1.GetId().String()].(*literal) != l1.(*literal) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*literal) != l1.(*literal) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[l1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*literal) != l1.(*literal) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[l1.GetId().String()].(*literal) != l1.(*literal) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoLiteralPointerCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	lp1 := uOfD.NewNameLiteralPointer()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[lp1.GetId().String()].(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[lp1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[lp1.GetId().String()].(*literalPointer) != lp1.(*literalPointer) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoLiteralPointerPointerCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	lp1 := uOfD.NewLiteralPointerPointer()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[lp1.GetId().String()].(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[lp1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[lp1.GetId().String()].(*literalPointerPointer) != lp1.(*literalPointerPointer) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoLiteralPointerPointerSetLiteralPointer(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	lp1 := uOfD.NewNameLiteralPointer()
	r1 := uOfD.NewLiteralPointerReference()
	uOfD.markUndoPoint()
	r1.SetLiteralPointer(lp1)
	rlp := r1.getLiteralPointerPointer()
	if rlp == nil {
		t.Error("Referenced element pointer is nil")
	}
	if rlp.GetOwningElement().GetId() != r1.GetId() {
		t.Error("Referenced element pointer owner not properly set")
	}
	if rlp.GetLiteralPointer() != lp1 {
		t.Error("Element pointer not set properly")
	}
	if r1.getOwnedBaseElements()[rlp.GetId().String()].GetId() != rlp.GetId() {
		t.Error("R1.ownedBaseElements does not contain rlp")
	}

	// Undo
	uOfD.undo()
	if r1.getLiteralPointerPointer() != nil {
		t.Error("Referenced literal pointer is not nil")
	}
	if rlp.GetOwningElement() != nil {
		t.Error("Owning element not properly cleared")
	}
	if rlp.GetLiteralPointer() != nil {
		t.Error("LiteralPointer not cleared properly")
	}
	if r1.getOwnedBaseElements()[rlp.GetId().String()] != nil {
		t.Error("R1.ownedBaseElements still contains rlp")
	}

	// Redo
	uOfD.redo()
	if r1.getLiteralPointerPointer() == nil {
		t.Error("Literal pointer pointer is nil")
	}
	if rlp.GetOwningElement().GetId() != r1.GetId() {
		t.Error("Owning element not properly set")
	}
	if rlp.GetLiteralPointer() != lp1 {
		t.Error("LiteralPointer not set properly")
	}
	if r1.getOwnedBaseElements()[rlp.GetId().String()].GetId() != rlp.GetId() {
		t.Error("R1.ownedBaseElements does not contain rlp")
	}
}

func TestUndoRedoLiteralReferenceCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewLiteralReference()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*literalReference) != e1.(*literalReference) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*literalReference) != e1.(*literalReference) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*literalReference) != e1.(*literalReference) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*literalReference) != e1.(*literalReference) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*literalReference) != e1.(*literalReference) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestUndoRedoRefinementCreation(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	uOfD.setRecordingUndo(true)
	e1 := uOfD.NewRefinement()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after creating Element")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Reso stack size incorrect after creating Element")
	}
	creationEntry := uOfD.undoStack.Peek()
	if creationEntry.changeType != Creation {
		t.Error("Creation entry change type incorrect")
	}
	if creationEntry.changedElement.(*refinement) != e1.(*refinement) {
		t.Error("Creation entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*refinement) != e1.(*refinement) {
		t.Error("Element not added to uOfD.baseElementMap after creation")
	}

	// Undo
	uOfD.undo()
	if len(uOfD.undoStack) != 0 {
		t.Error("Undo stack size incorrect after undo")
	}
	if len(uOfD.redoStack) != 1 {
		t.Error("Redo stack size incorrect after undo")
	}
	redoEntry := uOfD.redoStack.Peek()
	if redoEntry != creationEntry {
		t.Error("Creation entry not moved to redo stack after undo")
	}
	if redoEntry.changeType != Creation {
		t.Error("Redo entry changeType incorrect")
	}
	if redoEntry.changedElement.(*refinement) != e1.(*refinement) {
		t.Error("Redo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()] != nil {
		t.Error("Element not removed from uOfD.baseElementMap after undo")
	}

	// Redo
	uOfD.redo()
	if len(uOfD.undoStack) != 1 {
		t.Error("Undo stack size incorrect after redo")
	}
	if len(uOfD.redoStack) != 0 {
		t.Error("Redo stack size incorrect after redo")
	}
	undoEntry := uOfD.undoStack.Peek()
	if undoEntry.changeType != Creation {
		t.Error("Undo entry change type not Creation")
	}
	if undoEntry.changedElement.(*refinement) != e1.(*refinement) {
		t.Error("Undo entry new entry not nil")
	}
	if uOfD.baseElementMap[e1.GetId().String()].(*refinement) != e1.(*refinement) {
		t.Error("Element not added to uOfD.baseElementMap after redo")
	}
}

func TestAddElementListener(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	e1 := uOfD.NewElement()
	ep := uOfD.NewReferencedElementPointer()
	ep.SetElement(e1)
	elm := uOfD.elementListenerMap[e1.GetId().String()]
	if elm == nil {
		t.Error("ElementListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("ElementListenerMap entry length != 1")
		} else {
			if (*elm)[0] != ep {
				t.Error("ElementListenerMap entry does not contain ElementPointer")
			}
		}
	}
	ep.SetElement(nil)
	elm = uOfD.elementListenerMap[e1.GetId().String()]
	if elm == nil {
		t.Error("ElementListenerMap entry is nil after SetElement(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("ElementListenerMap entry length != 0")
		}
	}

}

func TestAddElementPointerListener(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	ep := uOfD.NewReferencedElementPointer()
	epp := uOfD.NewElementPointerPointer()
	epp.SetElementPointer(ep)
	elm := uOfD.elementPointerListenerMap[ep.GetId().String()]
	if elm == nil {
		t.Error("ElementPointerListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("ElementPointerListenerMap entry length != 1")
		} else {
			if (*elm)[0] != epp {
				t.Error("ElementPointerListenerMap entry does not contain ElementPointer")
			}
		}
	}
	epp.SetElementPointer(nil)
	elm = uOfD.elementPointerListenerMap[ep.GetId().String()]
	if elm == nil {
		t.Error("ElementListenerMap entry is nil after SetElement(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("ElementListenerMap entry length != 0")
		}
	}

}

func TestAddLiteralListener(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	e1 := uOfD.NewLiteral()
	lp := uOfD.NewNameLiteralPointer()
	lp.SetLiteral(e1)
	elm := uOfD.literalListenerMap[e1.GetId().String()]
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("LiteralListenerMap entry length != 1")
		} else {
			if (*elm)[0] != lp {
				t.Error("LiteralListenerMap entry does not contain LiteralPointer")
			}
		}
	}
	lp.SetLiteral(nil)
	elm = uOfD.literalListenerMap[e1.GetId().String()]
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil after SetLiteral(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("LiteralListenerMap entry length != 0")
		}
	}

}

func TestAddLiteralPointerListener(t *testing.T) {
	uOfD := NewUniverseOfDiscourse()
	lp := uOfD.NewNameLiteralPointer()
	lpp := uOfD.NewLiteralPointerPointer()
	lpp.SetLiteralPointer(lp)
	elm := uOfD.literalPointerListenerMap[lp.GetId().String()]
	if elm == nil {
		t.Error("LiteralPointerListenerMap entry is nil")
	} else {
		if len(*elm) != 1 {
			t.Error("LiteralPointerListenerMap entry length != 1")
		} else {
			if (*elm)[0] != lpp {
				t.Error("LiteralPointerListenerMap entry does not contain LiteralPointer")
			}
		}
	}
	lpp.SetLiteralPointer(nil)
	elm = uOfD.literalPointerListenerMap[lp.GetId().String()]
	if elm == nil {
		t.Error("LiteralListenerMap entry is nil after SetLiteral(nil)")
	} else {
		if len(*elm) != 0 {
			t.Error("LiteralListenerMap entry length != 0")
		}
	}

}
