// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	_ "github.com/hajimehoshi/guigui/basicwidget/cjkfont"
)

type Root struct {
	guigui.DefaultWidget

	background basicwidget.Background
	sidebar    Sidebar
	settings   Settings
	basic      Basic
	buttons    Buttons
	lists      Lists
	popups     Popups
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	appender.AppendChildWidget(&r.background)
	appender.AppendChildWidget(&r.sidebar)

	guigui.SetPosition(&r.sidebar, guigui.Position(r))
	rw, _ := guigui.Size(r)
	sw, _ := guigui.Size(&r.sidebar)
	p := guigui.Position(r)
	p.X += sw
	pw := rw - sw
	guigui.SetPosition(&r.settings, p)
	guigui.SetSize(&r.settings, pw, guigui.AutoSize)
	guigui.SetPosition(&r.basic, p)
	guigui.SetSize(&r.basic, pw, guigui.AutoSize)
	guigui.SetPosition(&r.buttons, p)
	guigui.SetSize(&r.buttons, pw, guigui.AutoSize)
	guigui.SetPosition(&r.lists, p)
	guigui.SetSize(&r.lists, pw, guigui.AutoSize)
	guigui.SetPosition(&r.popups, p)
	guigui.SetSize(&r.popups, pw, guigui.AutoSize)

	switch r.sidebar.SelectedItemTag() {
	case "settings":
		appender.AppendChildWidget(&r.settings)
	case "basic":
		appender.AppendChildWidget(&r.basic)
	case "buttons":
		appender.AppendChildWidget(&r.buttons)
	case "lists":
		appender.AppendChildWidget(&r.lists)
	case "popups":
		appender.AppendChildWidget(&r.popups)
	}

	return nil
}

func main() {
	op := &guigui.RunOptions{
		Title: "Component Gallery",
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
