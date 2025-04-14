// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package basicwidget

import (
	"image"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type FormItem struct {
	PrimaryWidget   guigui.Widget
	SecondaryWidget guigui.Widget
}

type Form struct {
	guigui.DefaultWidget

	items []*FormItem

	primaryBounds   []image.Rectangle
	secondaryBounds []image.Rectangle
}

func formItemPadding(context *guigui.Context) (int, int) {
	return UnitSize(context) / 2, UnitSize(context) / 4
}

func (f *Form) SetItems(items []*FormItem) {
	f.items = slices.Delete(f.items, 0, len(f.items))
	f.items = append(f.items, items...)
}

func (f *Form) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	f.calcItemBounds(context)

	for i, item := range f.items {
		if item.PrimaryWidget != nil {
			context.SetPosition(item.PrimaryWidget, f.primaryBounds[i].Min)
			appender.AppendChildWidget(item.PrimaryWidget)
		}
		if item.SecondaryWidget != nil {
			context.SetPosition(item.SecondaryWidget, f.secondaryBounds[i].Min)
			appender.AppendChildWidget(item.SecondaryWidget)
		}
	}

	return nil
}

func (f *Form) calcItemBounds(context *guigui.Context) {
	f.primaryBounds = slices.Delete(f.primaryBounds, 0, len(f.primaryBounds))
	f.secondaryBounds = slices.Delete(f.secondaryBounds, 0, len(f.secondaryBounds))

	paddingX, paddingY := formItemPadding(context)

	var y int
	for i, item := range f.items {
		f.primaryBounds = append(f.primaryBounds, image.Rectangle{})
		f.secondaryBounds = append(f.secondaryBounds, image.Rectangle{})

		if item.PrimaryWidget == nil && item.SecondaryWidget == nil {
			continue
		}
		if !context.IsVisible(item.SecondaryWidget) {
			continue
		}

		var primaryH int
		var secondaryH int
		if item.PrimaryWidget != nil {
			_, primaryH = context.Size(item.PrimaryWidget)
		}
		if item.SecondaryWidget != nil {
			_, secondaryH = context.Size(item.SecondaryWidget)
		}
		h := max(primaryH, secondaryH, minFormItemHeight(context))
		baseBounds := context.Bounds(f)
		baseBounds.Min.X += paddingX
		baseBounds.Max.X -= paddingX
		baseBounds.Min.Y += y
		baseBounds.Max.Y = baseBounds.Min.Y + h

		if item.PrimaryWidget != nil {
			bounds := baseBounds
			ww, wh := context.Size(item.PrimaryWidget)
			bounds.Max.X = bounds.Min.X + ww
			pY := (h + 2*paddingY - wh) / 2
			if wh < UnitSize(context)+2*paddingY {
				pY = min(pY, max(0, (UnitSize(context)+2*paddingY-wh)/2))
			}
			bounds.Min.Y += pY
			bounds.Max.Y += pY
			f.primaryBounds[i] = bounds
		}
		if item.SecondaryWidget != nil {
			bounds := baseBounds
			ww, wh := context.Size(item.SecondaryWidget)
			bounds.Min.X = bounds.Max.X - ww
			pY := (h + 2*paddingY - wh) / 2
			if wh < UnitSize(context)+2*paddingY {
				pY = min(pY, (UnitSize(context)+2*paddingY-wh)/2)
			}
			bounds.Min.Y += pY
			bounds.Max.Y += pY
			f.secondaryBounds[i] = bounds
		}

		y += h + 2*paddingY
	}
}

func (f *Form) Draw(context *guigui.Context, dst *ebiten.Image) {
	bounds := context.Bounds(f)
	bounds.Max.Y = bounds.Min.Y + f.height(context)
	draw.DrawRoundedRect(context, dst, bounds, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.925), RoundedCornerRadius(context))

	if len(f.items) > 0 {
		paddingX, paddingY := formItemPadding(context)
		y := paddingY
		for _, item := range f.items[:len(f.items)-1] {
			var primaryH int
			var secondaryH int
			if item.PrimaryWidget != nil {
				_, primaryH = context.Size(item.PrimaryWidget)
			}
			if item.SecondaryWidget != nil {
				_, secondaryH = context.Size(item.SecondaryWidget)
			}
			h := max(primaryH, secondaryH, minFormItemHeight(context))
			y += h + 2*paddingY

			x0 := float32(bounds.Min.X + paddingX)
			x1 := float32(bounds.Max.X - paddingX)
			y := float32(y) + float32(paddingY)
			width := 1 * float32(context.Scale())
			clr := draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.875)
			vector.StrokeLine(dst, x0, y, x1, y, width, clr, false)
		}
	}

	draw.DrawRoundedRectBorder(context, dst, bounds, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.875), RoundedCornerRadius(context), 1*float32(context.Scale()), draw.RoundedRectBorderTypeRegular)
}

func (f *Form) DefaultSize(context *guigui.Context) (int, int) {
	return 6 * UnitSize(context), f.height(context)
}

func (f *Form) height(context *guigui.Context) int {
	_, paddingY := formItemPadding(context)

	var y int
	for _, item := range f.items {
		if (item.PrimaryWidget == nil || !context.IsVisible(item.PrimaryWidget)) &&
			(item.SecondaryWidget == nil || !context.IsVisible(item.SecondaryWidget)) {
			continue
		}
		var primaryH int
		var secondaryH int
		if item.PrimaryWidget != nil {
			_, primaryH = context.Size(item.PrimaryWidget)
		}
		if item.SecondaryWidget != nil {
			_, secondaryH = context.Size(item.SecondaryWidget)
		}
		h := max(primaryH, secondaryH, minFormItemHeight(context))
		y += h + 2*paddingY
	}
	return y
}

func minFormItemHeight(context *guigui.Context) int {
	return UnitSize(context)
}
