// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Guigui Authors

package main

import (
	"image"
	"slices"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type TooltipAreas struct {
	guigui.DefaultWidget

	button       basicwidget.Button
	text         basicwidget.Text
	selectWidget basicwidget.Select[int]
	tooltipArea1 basicwidget.TooltipArea
	tooltipArea2 basicwidget.TooltipArea
	tooltipArea3 basicwidget.TooltipArea

	layoutItems     []guigui.LinearLayoutItem
	selectRowItems  []guigui.LinearLayoutItem
	selectRowLayout guigui.LinearLayout
	itemBoundsArr   []image.Rectangle
}

func (t *TooltipAreas) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&t.button)
	adder.AddWidget(&t.tooltipArea1)
	adder.AddWidget(&t.text)
	adder.AddWidget(&t.tooltipArea2)
	adder.AddWidget(&t.selectWidget)
	adder.AddWidget(&t.tooltipArea3)

	t.button.SetText("Hover me")
	t.tooltipArea1.SetText("This is a button tooltip")

	t.text.SetValue("Hover over this text to see a tooltip")
	t.tooltipArea2.SetText("This is a text tooltip")

	t.selectWidget.SetItemsByStrings([]string{
		"Apple",
		"Banana",
		"Cherry",
		"Date",
		"Elderberry",
		"Fig",
		"Grape",
		"Honeydew",
		"Kiwi",
		"Lemon",
	})
	if t.selectWidget.SelectedItemIndex() < 0 {
		t.selectWidget.SelectItemByIndex(0)
	}
	t.tooltipArea3.SetText("This is a select tooltip")

	return nil
}

func (t *TooltipAreas) layout(context *guigui.Context) guigui.LinearLayout {
	u := basicwidget.UnitSize(context)

	t.selectRowItems = slices.Delete(t.selectRowItems, 0, len(t.selectRowItems))
	t.selectRowItems = append(t.selectRowItems,
		guigui.LinearLayoutItem{
			Widget: &t.selectWidget,
		},
		guigui.LinearLayoutItem{
			Size: guigui.FlexibleSize(1),
		},
	)
	t.selectRowLayout = guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items:     t.selectRowItems,
	}

	t.layoutItems = slices.Delete(t.layoutItems, 0, len(t.layoutItems))
	t.layoutItems = append(t.layoutItems,
		guigui.LinearLayoutItem{
			Widget: &t.button,
		},
		guigui.LinearLayoutItem{
			Widget: &t.text,
		},
		guigui.LinearLayoutItem{
			Layout: &t.selectRowLayout,
		},
	)
	return guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     t.layoutItems,
		Gap:       u / 2,
		Padding: guigui.Padding{
			Start:  u / 2,
			Top:    u / 2,
			End:    u / 2,
			Bottom: u / 2,
		},
	}
}

func (t *TooltipAreas) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layout := t.layout(context)
	layout.LayoutWidgets(context, widgetBounds.Bounds(), layouter)

	t.itemBoundsArr = layout.AppendItemBounds(t.itemBoundsArr[:0], context, widgetBounds.Bounds())
	layouter.LayoutWidget(&t.tooltipArea1, t.itemBoundsArr[0])
	layouter.LayoutWidget(&t.tooltipArea2, t.itemBoundsArr[1])
	selectBounds := t.selectRowLayout.ItemBoundsAt(0, context, t.itemBoundsArr[2])
	layouter.LayoutWidget(&t.tooltipArea3, selectBounds)
}
