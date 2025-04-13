// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package main

import (
	"fmt"
	"image"
	"os"
	"slices"
	"strings"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	_ "github.com/hajimehoshi/guigui/basicwidget/cjkfont"
)

var theCurrentID int

func nextTaskID() int {
	theCurrentID++
	return theCurrentID
}

type Task struct {
	ID   int
	Text string
}

func NewTask(text string) Task {
	return Task{
		ID:   nextTaskID(),
		Text: text,
	}
}

type Root struct {
	guigui.DefaultWidget

	background        basicwidget.Background
	createButton      basicwidget.TextButton
	textField         basicwidget.TextField
	tasksPanel        basicwidget.ScrollablePanel
	tasksPanelContent tasksPanelContent

	tasks []Task
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	appender.AppendChildWidget(&r.background)

	u := float64(basicwidget.UnitSize(context))

	width, _ := context.Size(r)
	w := width - int(6.5*u)
	context.SetSize(&r.textField, w, int(u))
	r.textField.SetOnEnterPressed(func(text string) {
		r.tryCreateTask(context)
	})
	{
		context.SetPosition(&r.textField, context.Position(r).Add(image.Pt(int(0.5*u), int(0.5*u))))
		appender.AppendChildWidget(&r.textField)
	}

	r.createButton.SetText(context, "Create")
	context.SetSize(&r.createButton, int(5*u), guigui.AutoSize)
	r.createButton.SetOnUp(func() {
		r.tryCreateTask(context)
	})
	if r.canCreateTask() {
		context.Enable(&r.createButton)
	} else {
		context.Disable(&r.createButton)
	}
	{
		p := context.Position(r)
		w, _ := context.Size(r)
		p.X += w - int(0.5*u) - int(5*u)
		p.Y += int(0.5 * u)
		context.SetPosition(&r.createButton, p)
		appender.AppendChildWidget(&r.createButton)
	}

	w, h := context.Size(r)
	context.SetSize(&r.tasksPanel, w, h-int(2*u))
	r.tasksPanelContent.SetTasks(context, r.tasks)
	r.tasksPanelContent.SetOnDeleted(func(id int) {
		r.tasks = slices.DeleteFunc(r.tasks, func(t Task) bool {
			return t.ID == id
		})
	})
	r.tasksPanel.SetContent(&r.tasksPanelContent)
	context.SetPosition(&r.tasksPanel, context.Position(r).Add(image.Pt(0, int(2*u))))
	appender.AppendChildWidget(&r.tasksPanel)

	return nil
}

func (r *Root) canCreateTask() bool {
	str := r.textField.Text()
	str = strings.TrimSpace(str)
	return str != ""
}

func (r *Root) tryCreateTask(context *guigui.Context) {
	str := r.textField.Text()
	str = strings.TrimSpace(str)
	if str != "" {
		r.tasks = slices.Insert(r.tasks, 0, NewTask(str))
		r.textField.SetText(context, "")
	}
}

type taskWidget struct {
	guigui.DefaultWidget

	doneButton basicwidget.TextButton
	text       basicwidget.Text

	onDoneButtonPressed func()
}

func (t *taskWidget) SetOnDoneButtonPressed(f func()) {
	t.onDoneButtonPressed = f
}

func (t *taskWidget) SetText(context *guigui.Context, text string) {
	t.text.SetText(context, text)
}

func (t *taskWidget) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := float64(basicwidget.UnitSize(context))

	p := context.Position(t)
	t.doneButton.SetText(context, "Done")
	context.SetSize(&t.doneButton, int(3*u), guigui.AutoSize)
	t.doneButton.SetOnUp(func() {
		if t.onDoneButtonPressed != nil {
			t.onDoneButtonPressed()
		}
	})
	context.SetPosition(&t.doneButton, p)
	appender.AppendChildWidget(&t.doneButton)

	w, _ := context.Size(t)
	context.SetSize(&t.text, w-int(4.5*u), int(u))
	t.text.SetVerticalAlign(context, basicwidget.VerticalAlignMiddle)
	context.SetPosition(&t.text, image.Pt(p.X+int(3.5*u), p.Y))
	appender.AppendChildWidget(&t.text)
	return nil
}

func (t *taskWidget) DefaultSize(context *guigui.Context) (int, int) {
	w, _ := context.Size(guigui.Parent(t))
	return w, int(basicwidget.UnitSize(context))
}

type tasksPanelContent struct {
	guigui.DefaultWidget

	taskWidgets []taskWidget

	onDeleted func(id int)
}

func (t *tasksPanelContent) SetOnDeleted(f func(id int)) {
	t.onDeleted = f
}

func (t *tasksPanelContent) SetTasks(context *guigui.Context, tasks []Task) {
	if len(tasks) != len(t.taskWidgets) {
		if len(tasks) > len(t.taskWidgets) {
			t.taskWidgets = slices.Grow(t.taskWidgets, len(tasks)-len(t.taskWidgets))
			t.taskWidgets = t.taskWidgets[:len(tasks)]
		} else {
			t.taskWidgets = slices.Delete(t.taskWidgets, len(tasks), len(t.taskWidgets))
		}
	}
	for i, task := range tasks {
		t.taskWidgets[i].SetOnDoneButtonPressed(func() {
			if t.onDeleted != nil {
				t.onDeleted(task.ID)
			}
		})
		t.taskWidgets[i].SetText(context, task.Text)
	}
}

func (t *tasksPanelContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	u := float64(basicwidget.UnitSize(context))

	p := context.Position(t)
	x := p.X + int(0.5*u)
	y := p.Y
	for i := range t.taskWidgets {
		// Do not take a variable for the task widget in the for-range loop,
		// since the pointer value of the task widget matters.
		if i > 0 {
			y += int(u / 4)
		}
		context.SetPosition(&t.taskWidgets[i], image.Pt(x, y))
		appender.AppendChildWidget(&t.taskWidgets[i])
		y += int(u)
	}

	return nil
}

func (t *tasksPanelContent) DefaultSize(context *guigui.Context) (int, int) {
	u := basicwidget.UnitSize(context)

	w, _ := context.Size(guigui.Parent(t))
	c := len(t.taskWidgets)
	h := c * (u + u/4)
	return w, h
}

func main() {
	op := &guigui.RunOptions{
		Title:           "TODO",
		WindowMinWidth:  320,
		WindowMinHeight: 240,
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
