// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Popups struct {
	guigui.DefaultWidget

	forms                        [2]basicwidget.Form
	blurBackgroundText           basicwidget.Text
	blurBackgroundToggle         basicwidget.Toggle
	closeByClickingOutsideText   basicwidget.Text
	closeByClickingOutsideToggle basicwidget.Toggle
	showButton                   basicwidget.TextButton

	contextMenuPopupText          basicwidget.Text
	contextMenuPopupClickHereText basicwidget.Text

	simplePopup        basicwidget.Popup
	simplePopupContent simplePopupContent

	contextMenuPopup basicwidget.PopupMenu
}

func (p *Popups) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	p.blurBackgroundText.SetText("Blur Background")
	p.closeByClickingOutsideText.SetText("Close by Clicking Outside")
	p.showButton.SetText("Show")
	p.showButton.SetOnUp(func() {
		p.simplePopup.Open(context)
	})

	p.forms[0].SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &p.blurBackgroundText,
			SecondaryWidget: &p.blurBackgroundToggle,
		},
		{
			PrimaryWidget:   &p.closeByClickingOutsideText,
			SecondaryWidget: &p.closeByClickingOutsideToggle,
		},
		{
			SecondaryWidget: &p.showButton,
		},
	})

	p.contextMenuPopupText.SetText("Context Menu")
	p.contextMenuPopupClickHereText.SetText("Click Here by the Right Button")

	p.forms[1].SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &p.contextMenuPopupText,
			SecondaryWidget: &p.contextMenuPopupClickHereText,
		},
	})

	u := basicwidget.UnitSize(context)
	for i, bounds := range (layout.GridLayout{
		Bounds: context.Bounds(p).Inset(u / 2),
		Heights: []layout.Size{
			layout.MaxContentSize(func(index int) int {
				if index >= len(p.forms) {
					return 0
				}
				return context.Size(&p.forms[index]).Y
			}),
		},
		RowGap: u / 2,
	}).RepeatingCellBounds() {
		if i >= len(p.forms) {
			break
		}
		appender.AppendChildWidgetWithBounds(&p.forms[i], bounds)
	}

	p.simplePopupContent.popup = &p.simplePopup
	p.simplePopup.SetContent(&p.simplePopupContent)
	p.simplePopup.SetBackgroundBlurred(p.blurBackgroundToggle.Value())
	p.simplePopup.SetCloseByClickingOutside(p.closeByClickingOutsideToggle.Value())
	p.simplePopup.SetAnimationDuringFade(true)

	appBounds := context.AppBounds()
	contentSize := image.Pt(int(12*u), int(6*u))
	simplePopupPosition := image.Point{
		X: appBounds.Min.X + (appBounds.Dx()-contentSize.X)/2,
		Y: appBounds.Min.Y + (appBounds.Dy()-contentSize.Y)/2,
	}
	simplePopupBounds := image.Rectangle{
		Min: simplePopupPosition,
		Max: simplePopupPosition.Add(contentSize),
	}
	context.SetSize(&p.simplePopupContent, simplePopupBounds.Size())
	appender.AppendChildWidgetWithBounds(&p.simplePopup, simplePopupBounds)

	p.contextMenuPopup.SetItemsByStrings([]string{"Item 1", "Item 2", "Item 3"})
	// A context menu's position is updated at HandlePointingInput.
	appender.AppendChildWidget(&p.contextMenuPopup)

	return nil
}

func (p *Popups) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		pt := image.Pt(ebiten.CursorPosition())
		if pt.In(context.VisibleBounds(&p.contextMenuPopupClickHereText)) {
			context.SetPosition(&p.contextMenuPopup, pt)
			p.contextMenuPopup.Open(context)
		}
	}
	return guigui.HandleInputResult{}
}

type simplePopupContent struct {
	guigui.DefaultWidget

	popup *basicwidget.Popup

	titleText   basicwidget.Text
	closeButton basicwidget.TextButton
}

func (s *simplePopupContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := basicwidget.UnitSize(context)

	s.titleText.SetText("Hello!")
	s.titleText.SetBold(true)

	s.closeButton.SetText("Close")
	s.closeButton.SetOnUp(func() {
		s.popup.Close()
	})

	for i, bounds := range (layout.GridLayout{
		Bounds: context.Bounds(s).Inset(u / 2),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.MaxContentSize(func(index int) int {
				if index != 1 {
					return 0
				}
				return s.closeButton.DefaultSize(context).Y
			}),
		},
	}).CellBounds() {
		switch i {
		case 0:
			appender.AppendChildWidgetWithBounds(&s.titleText, bounds)
		case 1:
			for i, bounds := range (layout.GridLayout{
				Bounds: bounds,
				Widths: []layout.Size{
					layout.FlexibleSize(1),
					layout.MaxContentSize(func(index int) int {
						if index != 1 {
							return 0
						}
						return s.closeButton.DefaultSize(context).X
					}),
				},
			}).CellBounds() {
				if i != 1 {
					continue
				}
				appender.AppendChildWidgetWithBounds(&s.closeButton, bounds)
			}
		}
	}

	return nil
}
