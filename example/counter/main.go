// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

type Root struct {
	guigui.RootWidgetBehavior

	resetButton basicwidget.TextButton
	incButton   basicwidget.TextButton
	decButton   basicwidget.TextButton
	counterText basicwidget.Text

	counter int
}

func (r *Root) AppendChildWidgets(context *guigui.Context, appender *guigui.ChildWidgetAppender) {
	{
		w, h := r.Size(context)
		w -= 2 * basicwidget.UnitSize(context)
		h -= 4 * basicwidget.UnitSize(context)
		r.counterText.SetSize(w, h)
		p := context.Widget(r).Position()
		p.X += basicwidget.UnitSize(context)
		p.Y += basicwidget.UnitSize(context)
		appender.AppendChildWidget(&r.counterText, p)
	}

	r.resetButton.SetText("Reset")
	r.resetButton.SetWidth(6 * basicwidget.UnitSize(context))
	{
		p := context.Widget(r).Position()
		_, h := r.Size(context)
		p.X += basicwidget.UnitSize(context)
		p.Y += h - 2*basicwidget.UnitSize(context)
		appender.AppendChildWidget(&r.resetButton, p)
	}

	r.incButton.SetText("Increment")
	r.incButton.SetWidth(6 * basicwidget.UnitSize(context))
	{
		p := context.Widget(r).Position()
		w, h := r.Size(context)
		p.X += w - 7*basicwidget.UnitSize(context)
		p.Y += h - 2*basicwidget.UnitSize(context)
		appender.AppendChildWidget(&r.incButton, p)
	}

	r.decButton.SetText("Decrement")
	r.decButton.SetWidth(6 * basicwidget.UnitSize(context))
	{
		p := context.Widget(r).Position()
		w, h := r.Size(context)
		p.X += w - int(13.5*float64(basicwidget.UnitSize(context)))
		p.Y += h - 2*basicwidget.UnitSize(context)
		appender.AppendChildWidget(&r.decButton, p)
	}
}

func (r *Root) Update(context *guigui.Context) error {
	for e := range context.Widget(&r.incButton).DequeueEvents() {
		args := e.(basicwidget.ButtonEvent)
		if args.Type == basicwidget.ButtonEventTypeUp {
			r.counter++
		}
	}
	for e := range context.Widget(&r.decButton).DequeueEvents() {
		args := e.(basicwidget.ButtonEvent)
		if args.Type == basicwidget.ButtonEventTypeUp {
			r.counter--
		}
	}
	for e := range context.Widget(&r.resetButton).DequeueEvents() {
		args := e.(basicwidget.ButtonEvent)
		if args.Type == basicwidget.ButtonEventTypeUp {
			r.counter = 0
		}
	}

	if r.counter == 0 {
		context.Widget(&r.resetButton).Disable()
	} else {
		context.Widget(&r.resetButton).Enable()
	}
	r.counterText.SetSelectable(true)
	r.counterText.SetBold(true)
	r.counterText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	r.counterText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	r.counterText.SetScale(4)
	r.counterText.SetText(fmt.Sprintf("%d", r.counter))

	return nil
}

func (r *Root) Draw(context *guigui.Context, dst *ebiten.Image) {
	basicwidget.FillBackground(dst, context)
}

func main() {
	op := &guigui.RunOptions{
		Title:           "Counter",
		WindowMinWidth:  640,
		WindowMinHeight: 480,
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
