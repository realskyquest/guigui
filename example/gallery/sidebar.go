// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package main

import (
	"sync"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Sidebar struct {
	guigui.DefaultWidget

	sidebar        basicwidget.Sidebar
	sidebarContent sidebarContent
}

func (s *Sidebar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	w, h := context.Size(s)
	context.SetSize(&s.sidebar, w, h)
	context.SetSize(&s.sidebarContent, w, h)
	s.sidebar.SetContent(&s.sidebarContent)
	context.SetPosition(&s.sidebar, context.Position(s))
	appender.AppendChildWidget(&s.sidebar)

	return nil
}

func (s *Sidebar) SelectedItemTag() string {
	return s.sidebarContent.SelectedItemTag()
}

func (s *Sidebar) SetSelectedItemIndex(index int) {
	s.sidebarContent.SetSelectedItemIndex(index)
}

type sidebarContent struct {
	guigui.DefaultWidget

	list            basicwidget.List
	listItemWidgets []basicwidget.Text

	initOnce sync.Once
}

func (s *sidebarContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.list.SetStyle(basicwidget.ListStyleSidebar)

	type item struct {
		text string
		tag  string
	}
	items := []item{
		{
			text: "Settings",
			tag:  "settings",
		},
		{
			text: "Basic",
			tag:  "basic",
		},
		{
			text: "Buttons",
			tag:  "buttons",
		},
		{
			text: "Lists",
			tag:  "lists",
		},
		{
			text: "Popups",
			tag:  "popups",
		},
	}

	if len(s.listItemWidgets) == 0 {
		s.listItemWidgets = make([]basicwidget.Text, len(items))
	}
	listItems := make([]basicwidget.ListItem, len(items))
	for i, item := range items {
		t := &s.listItemWidgets[i]
		t.SetText(item.text)
		listItems[i] = basicwidget.ListItem{
			Content:    t,
			Selectable: true,
			Tag:        item.tag,
		}
	}
	for i := range s.listItemWidgets {
		t := &s.listItemWidgets[i]
		t.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
		context.SetSize(t, guigui.DefaultSize, basicwidget.UnitSize(context))
		if s.list.SelectedItemIndex() == i {
			t.SetColor(basicwidget.DefaultActiveListItemTextColor(context))
		} else {
			t.SetColor(basicwidget.DefaultTextColor(context))
		}
	}
	s.list.SetItems(listItems)

	s.initOnce.Do(func() {
		s.list.SetSelectedItemIndex(0)
	})

	w, h := context.Size(s)
	context.SetSize(&s.list, w, h)
	context.SetPosition(&s.list, context.Position(s))
	appender.AppendChildWidget(&s.list)

	return nil
}

func (s *sidebarContent) SelectedItemTag() string {
	item, ok := s.list.SelectedItem()
	if !ok {
		return ""
	}
	return item.Tag.(string)
}

func (s *sidebarContent) SetSelectedItemIndex(index int) {
	s.list.SetSelectedItemIndex(index)
}
