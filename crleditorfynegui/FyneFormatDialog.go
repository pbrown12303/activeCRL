package crleditorfynegui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/diagramwidget"
)

// FyneFormatDialog is the screen used to edit DiagramElement properties
type FyneFormatDialog struct {
	properties diagramwidget.DiagramElementProperties
	callback   func(diagramwidget.DiagramElementProperties)
}

// ShowFyneFormatDialog displays the dialog for editing the DiagramElement properties
func ShowFyneFormatDialog(properties diagramwidget.DiagramElementProperties, callback func(diagramwidget.DiagramElementProperties)) {
	fd := &FyneFormatDialog{
		properties: properties,
		callback:   callback,
	}
	rect := canvas.NewRectangle(properties.ForegroundColor)
	rect.StrokeColor = properties.ForegroundColor
	rect.FillColor = properties.BackgroundColor
	rect.StrokeWidth = 2
	rect.SetMinSize(fyne.NewSize(200, 50))
	rectContainer := container.NewPadded(rect)
	fgButton := widget.NewButton("Choose Foreground Color", func() {
		picker := dialog.NewColorPicker("Color Picker", "Choose Foreground Color", func(c color.Color) {
			fd.properties.ForegroundColor = c
			rect.StrokeColor = fd.properties.ForegroundColor
			rect.Refresh()
		}, FyneGUISingleton.window)
		picker.Advanced = true
		picker.SetColor(properties.ForegroundColor)
		picker.Show()
	})
	bgButton := widget.NewButton("Choose Background Color", func() {
		picker := dialog.NewColorPicker("Color Picker", "Choose Background Color", func(c color.Color) {
			fd.properties.BackgroundColor = c
			rect.FillColor = fd.properties.BackgroundColor
			rect.Refresh()
		}, FyneGUISingleton.window)
		picker.Advanced = true
		picker.SetColor(properties.BackgroundColor)
		picker.Show()
	})
	dialogBox := container.NewVBox(rectContainer, fgButton, bgButton)
	dialog.ShowCustomConfirm("Diagram Element Format", "Apply", "Cancel", dialogBox, fd.apply, FyneGUISingleton.window)
}

func (fd *FyneFormatDialog) apply(apply bool) {
	if apply {
		fd.callback(fd.properties)
	}
}
