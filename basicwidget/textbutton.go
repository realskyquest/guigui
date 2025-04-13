// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package basicwidget

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type TextButton struct {
	guigui.DefaultWidget

	button Button
	text   Text
	image  Image

	textColor color.Color
}

func (t *TextButton) SetOnDown(f func()) {
	t.button.SetOnDown(f)
}

func (t *TextButton) SetOnUp(f func()) {
	t.button.SetOnUp(f)
}

func (t *TextButton) SetText(context *guigui.Context, text string) {
	t.text.SetText(context, text)
}

func (t *TextButton) SetImage(context *guigui.Context, image *ebiten.Image) {
	t.image.SetImage(context, image)
}

func (t *TextButton) SetTextColor(context *guigui.Context, clr color.Color) {
	if draw.EqualColor(t.textColor, clr) {
		return
	}
	t.textColor = clr
	context.RequestRedraw(t)
}

func (t *TextButton) SetForcePressed(forcePressed bool) {
	t.button.SetForcePressed(forcePressed)
}

func (t *TextButton) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	w, h := context.Size(t)
	context.SetSize(&t.button, w, h)
	context.SetPosition(&t.button, context.Position(t))
	appender.AppendChildWidget(&t.button)

	imgSize := textButtonImageSize(context)

	tw, _ := t.text.TextSize(context)
	context.SetSize(&t.text, tw, h)
	if !context.IsEnabled(&t.button) {
		t.text.SetColor(context, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.5))
	} else {
		t.text.SetColor(context, t.textColor)
	}
	t.text.SetHorizontalAlign(context, HorizontalAlignCenter)
	t.text.SetVerticalAlign(context, VerticalAlignMiddle)
	textP := context.Position(t)
	if t.image.HasImage() {
		textP.X += (w - tw + UnitSize(context)/4) / 2
		textP.X -= (textButtonTextAndImagePadding(context) + imgSize) / 2
	} else {
		textP.X += (w - tw) / 2
	}
	if t.button.isActive(context) {
		textP.Y += int(1 * context.Scale())
	}
	context.SetPosition(&t.text, textP)
	appender.AppendChildWidget(&t.text)

	context.SetSize(&t.image, imgSize, imgSize)
	imgP := context.Position(t)
	imgP.X = textP.X + tw + textButtonTextAndImagePadding(context)
	imgP.Y += (h - imgSize) / 2
	if t.button.isActive(context) {
		imgP.Y += int(1 * context.Scale())
	}
	context.SetPosition(&t.image, imgP)
	appender.AppendChildWidget(&t.image)

	return nil
}

func (t *TextButton) DefaultSize(context *guigui.Context) (int, int) {
	_, dh := defaultButtonSize(context)
	w, _ := t.text.TextSize(context)
	if t.image.HasImage() {
		imgSize := textButtonImageSize(context)
		return w + textButtonTextAndImagePadding(context) + imgSize + UnitSize(context)*3/4, dh
	}
	return w + UnitSize(context), dh
}

func textButtonImageSize(context *guigui.Context) int {
	return int(LineHeight(context))
}

func textButtonTextAndImagePadding(context *guigui.Context) int {
	return UnitSize(context) / 4
}
