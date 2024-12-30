// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package main

import (
	"fmt"
	"image"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	_ "github.com/hajimehoshi/guigui/basicwidget/cjkfont"
)

type Task struct {
	ID        uuid.UUID
	Text      string
	CreatedAt time.Time
}

func NewTask(text string) Task {
	return Task{
		ID:        uuid.New(),
		Text:      text,
		CreatedAt: time.Now(),
	}
}

type TaskWidgets struct {
	doneButtonWidget *guigui.Widget
	textWidget       *guigui.Widget
}

type Root struct {
	guigui.DefaultWidgetBehavior

	createButtonWidget *guigui.Widget
	textFieldWidget    *guigui.Widget
	taskWidgets        map[uuid.UUID]*TaskWidgets
	tasksPanelWidget   *guigui.Widget

	tasks []Task
}

func (r *Root) AppendChildWidgets(context *guigui.Context, widget *guigui.Widget, appender *guigui.ChildWidgetAppender) {
	u := float64(basicwidget.UnitSize(context))

	if r.textFieldWidget == nil {
		var t basicwidget.TextField
		w := widget.Bounds().Dx() - int(6.5*u)
		t.SetSize(context, w, int(u))
		r.textFieldWidget = guigui.NewWidget(&t)
	}
	{
		b := widget.Bounds()
		x := b.Min.X + int(0.5*u)
		y := b.Min.Y + int(0.5*u)
		appender.AppendChildWidget(r.textFieldWidget, image.Pt(x, y))
	}

	if r.createButtonWidget == nil {
		var b basicwidget.TextButton
		b.SetText("Create")
		b.SetSize(context, int(5*u), int(u))
		r.createButtonWidget = guigui.NewWidget(&b)
	}
	{
		b := widget.Bounds()
		x := b.Max.X - int(0.5*u) - int(5*u)
		y := b.Min.Y + int(0.5*u)
		appender.AppendChildWidget(r.createButtonWidget, image.Pt(x, y))
	}

	if r.tasksPanelWidget == nil {
		var sp basicwidget.ScrollablePanel
		r.tasksPanelWidget = guigui.NewWidget(&sp)
	}
	tasksSP := r.tasksPanelWidget.Behavior().(*basicwidget.ScrollablePanel)
	tasksSP.SetContent(func(context *guigui.Context, widget *guigui.Widget, childAppender *basicwidget.ScrollablePanelChildWidgetAppender) {
		bounds := widget.Bounds()
		minX := bounds.Min.X + int(0.5*u)
		y := bounds.Min.Y
		for i, t := range r.tasks {
			if _, ok := r.taskWidgets[t.ID]; !ok {
				var b basicwidget.TextButton
				b.SetText("Done")
				b.SetSize(context, int(3*u), int(u))
				var text basicwidget.Text
				text.SetText(t.Text)
				text.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
				if r.taskWidgets == nil {
					r.taskWidgets = map[uuid.UUID]*TaskWidgets{}
				}
				r.taskWidgets[t.ID] = &TaskWidgets{
					doneButtonWidget: guigui.NewWidget(&b),
					textWidget:       guigui.NewWidget(&text),
				}
			}
			if i > 0 {
				y += int(u / 4)
			}
			childAppender.AppendChildWidget(r.taskWidgets[t.ID].doneButtonWidget, image.Pt(minX, y))
			_, textH := r.taskWidgets[t.ID].textWidget.Size(context)
			childAppender.AppendChildWidget(r.taskWidgets[t.ID].textWidget, image.Pt(minX+int(3.5*u), y+int((u-float64(textH))/2)))
			y += int(u)
		}
	})
	tasksSP.SetPadding(0, int(0.5*u))
	{
		b := widget.Bounds()
		b.Min.Y += int(2 * u)
		appender.AppendChildWidgetWithBounds(r.tasksPanelWidget, b)
	}

	// GC widgets
	for id := range r.taskWidgets {
		if slices.IndexFunc(r.tasks, func(t Task) bool {
			return t.ID == id
		}) >= 0 {
			continue
		}
		delete(r.taskWidgets, id)
	}
}

func (r *Root) Update(context *guigui.Context, widget *guigui.Widget) error {
	for event := range r.createButtonWidget.DequeueEvents() {
		switch e := event.(type) {
		case basicwidget.ButtonEvent:
			if e.Type == basicwidget.ButtonEventTypeUp {
				r.tryCreateTask()
			}
		}
	}
	for event := range r.textFieldWidget.DequeueEvents() {
		switch e := event.(type) {
		case basicwidget.TextEvent:
			if e.Type == basicwidget.TextEventTypeEnterPressed {
				r.tryCreateTask()
			}
		}
	}
	for id, t := range r.taskWidgets {
		for event := range t.doneButtonWidget.DequeueEvents() {
			switch e := event.(type) {
			case basicwidget.ButtonEvent:
				if e.Type == basicwidget.ButtonEventTypeUp {
					r.tasks = slices.DeleteFunc(r.tasks, func(task Task) bool {
						return task.ID == id
					})
				}
			}
		}
	}

	if r.canCreateTask() {
		r.createButtonWidget.Enable()
	} else {
		r.createButtonWidget.Disable()
	}

	return nil
}

func (r *Root) canCreateTask() bool {
	t := r.textFieldWidget.Behavior().(*basicwidget.TextField)
	str := t.Text()
	str = strings.TrimSpace(str)
	return str != ""
}

func (r *Root) tryCreateTask() {
	t := r.textFieldWidget.Behavior().(*basicwidget.TextField)
	str := t.Text()
	str = strings.TrimSpace(str)
	if str != "" {
		r.tasks = slices.Insert(r.tasks, 0, NewTask(str))
		t.SetText("")
	}
}

func (r *Root) Draw(context *guigui.Context, widget *guigui.Widget, dst *ebiten.Image) {
	basicwidget.FillBackground(dst, context)
}

func main() {
	op := &guigui.RunOptions{
		Title:           "TODO",
		WindowMinWidth:  320,
		WindowMinHeight: 240,
	}
	if err := guigui.Run(guigui.NewWidget(&Root{}), op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
