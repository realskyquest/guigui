// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package main

import (
	"sync"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Texts struct {
	guigui.DefaultWidget

	horizontalAlign basicwidget.HorizontalAlign
	verticalAlign   basicwidget.VerticalAlign
	unwrap          bool
	bold            bool
	selectable      bool
	editable        bool

	form                        basicwidget.Form
	horizontalAlignText         basicwidget.Text
	horizontalAlignDropdownList basicwidget.DropdownList
	verticalAlignText           basicwidget.Text
	verticalAlignDropdownList   basicwidget.DropdownList
	autoWrapText                basicwidget.Text
	autoWrapToggleSwitch        basicwidget.ToggleSwitch
	boldText                    basicwidget.Text
	boldToggleSwitch            basicwidget.ToggleSwitch
	selectableText              basicwidget.Text
	selectableToggleSwitch      basicwidget.ToggleSwitch
	editableText                basicwidget.Text
	editableToggleSwitch        basicwidget.ToggleSwitch
	sampleText                  basicwidget.Text

	initOnce sync.Once
}

const sampleText = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
隴西の李徴は博学才穎、天宝の末年、若くして名を虎榜に連ね、ついで江南尉に補せられたが、性、狷介、自ら恃むところ頗る厚く、賤吏に甘んずるを潔しとしなかった。
(Note: Selecting and editing auto-wrap texts are broken now.)`

func (t *Texts) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	t.horizontalAlignText.SetText("Horizontal Align")
	t.horizontalAlignDropdownList.SetItemsByStrings([]string{
		"Start",
		"Center",
		"End",
	})
	t.horizontalAlignDropdownList.SetOnValueChanged(func(index int) {
		switch index {
		case 0:
			t.horizontalAlign = basicwidget.HorizontalAlignStart
		case 1:
			t.horizontalAlign = basicwidget.HorizontalAlignCenter
		case 2:
			t.horizontalAlign = basicwidget.HorizontalAlignEnd
		}
	})

	t.verticalAlignText.SetText("Vertical Align")
	t.verticalAlignDropdownList.SetItemsByStrings([]string{
		"Top",
		"Middle",
		"Bottom",
	})
	t.verticalAlignDropdownList.SetOnValueChanged(func(index int) {
		switch index {
		case 0:
			t.verticalAlign = basicwidget.VerticalAlignTop
		case 1:
			t.verticalAlign = basicwidget.VerticalAlignMiddle
		case 2:
			t.verticalAlign = basicwidget.VerticalAlignBottom
		}
	})

	t.autoWrapText.SetText("Auto Wrap")
	t.autoWrapToggleSwitch.SetValue(!t.unwrap)
	t.autoWrapToggleSwitch.SetOnValueChanged(func(checked bool) {
		t.unwrap = !checked
	})

	t.boldText.SetText("Bold")
	t.boldToggleSwitch.SetValue(t.bold)
	t.boldToggleSwitch.SetOnValueChanged(func(checked bool) {
		t.bold = checked
	})

	t.selectableText.SetText("Selectable")
	t.selectableToggleSwitch.SetValue(t.selectable)
	t.selectableToggleSwitch.SetOnValueChanged(func(checked bool) {
		t.selectable = checked
		if !t.selectable {
			t.editable = false
		}
	})

	t.editableText.SetText("Editable")
	t.editableToggleSwitch.SetValue(t.editable)
	t.editableToggleSwitch.SetOnValueChanged(func(checked bool) {
		t.editable = checked
		if t.editable {
			t.selectable = true
		}
	})

	t.form.SetItems([]*basicwidget.FormItem{
		{
			PrimaryWidget:   &t.horizontalAlignText,
			SecondaryWidget: &t.horizontalAlignDropdownList,
		},
		{
			PrimaryWidget:   &t.verticalAlignText,
			SecondaryWidget: &t.verticalAlignDropdownList,
		},
		{
			PrimaryWidget:   &t.autoWrapText,
			SecondaryWidget: &t.autoWrapToggleSwitch,
		},
		{
			PrimaryWidget:   &t.boldText,
			SecondaryWidget: &t.boldToggleSwitch,
		},
		{
			PrimaryWidget:   &t.selectableText,
			SecondaryWidget: &t.selectableToggleSwitch,
		},
		{
			PrimaryWidget:   &t.editableText,
			SecondaryWidget: &t.editableToggleSwitch,
		},
	})

	t.sampleText.SetMultiline(true)
	t.sampleText.SetHorizontalAlign(t.horizontalAlign)
	t.sampleText.SetVerticalAlign(t.verticalAlign)
	t.sampleText.SetAutoWrap(!t.unwrap)
	t.sampleText.SetBold(t.bold)
	t.sampleText.SetSelectable(t.selectable)
	t.sampleText.SetEditable(t.editable)

	t.initOnce.Do(func() {
		t.sampleText.SetText(sampleText)
		t.horizontalAlignDropdownList.SetSelectedItemIndex(0)
		t.verticalAlignDropdownList.SetSelectedItemIndex(0)
	})

	u := basicwidget.UnitSize(context)
	for i, bounds := range (layout.GridLayout{
		Bounds: context.Bounds(t).Inset(u / 2),
		Heights: []layout.Size{
			layout.FixedSize(t.form.DefaultSize(context).Y),
			layout.FlexibleSize(1),
		},
		RowGap: u / 2,
	}).CellBounds() {
		switch i {
		case 0:
			appender.AppendChildWidgetWithBounds(&t.form, bounds)
		case 1:
			appender.AppendChildWidgetWithBounds(&t.sampleText, bounds)
		}
	}

	return nil
}
