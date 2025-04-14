// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package guigui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type bounds3D struct {
	bounds image.Rectangle
	z      int
}

type widgetsAndBounds struct {
	bounds3Ds map[Widget]bounds3D
}

func (w *widgetsAndBounds) reset() {
	clear(w.bounds3Ds)
}

func (w *widgetsAndBounds) append(widget Widget, bounds image.Rectangle) {
	if w.bounds3Ds == nil {
		w.bounds3Ds = map[Widget]bounds3D{}
	}
	w.bounds3Ds[widget] = bounds3D{
		bounds: bounds,
		z:      z(widget),
	}
}

func (w *widgetsAndBounds) equals(context *Context, currentWidgets []Widget) bool {
	if len(w.bounds3Ds) != len(currentWidgets) {
		return false
	}
	for _, widget := range currentWidgets {
		b, ok := w.bounds3Ds[widget]
		if !ok {
			return false
		}
		if b.bounds != context.Bounds(widget) {
			return false
		}
		if b.z != z(widget) {
			return false
		}
	}
	return true
}

func (w *widgetsAndBounds) redrawIfDifferentParentZ(app *app) {
	for widget, bounds3D := range w.bounds3Ds {
		if widget.ZDelta() != 0 {
			app.requestRedraw(bounds3D.bounds)
			RequestRedraw(widget)
		}
	}
}

type widgetState struct {
	root bool

	position    image.Point
	widthPlus1  int
	heightPlus1 int

	parent   Widget
	children []Widget
	prev     widgetsAndBounds

	hidden       bool
	disabled     bool
	transparency float64

	offscreen *ebiten.Image

	dirty bool
}

func (w *widgetState) isInTree() bool {
	p := w
	for ; p.parent != nil; p = p.parent.widgetState() {
	}
	return p.root
}

func (w *widgetState) isVisible() bool {
	if w.parent != nil {
		if w.hidden {
			return false
		}
		return w.parent.widgetState().isVisible()
	}
	return !w.hidden
}

func (w *widgetState) isEnabled() bool {
	if w.parent != nil {
		if w.disabled {
			return false
		}
		return w.parent.widgetState().isEnabled()
	}
	return !w.disabled
}

func (w *widgetState) opacity() float64 {
	return 1 - w.transparency
}

func (w *widgetState) ensureOffscreen(bounds image.Rectangle) *ebiten.Image {
	if w.offscreen != nil {
		if !bounds.In(w.offscreen.Bounds()) {
			w.offscreen.Deallocate()
			w.offscreen = nil
		}
	}
	if w.offscreen == nil {
		w.offscreen = ebiten.NewImageWithOptions(bounds, nil)
	}
	return w.offscreen.SubImage(bounds).(*ebiten.Image)
}

func traverseWidget(widget Widget, f func(widget Widget) error) error {
	if err := f(widget); err != nil {
		return err
	}
	for _, child := range widget.widgetState().children {
		if err := traverseWidget(child, f); err != nil {
			return err
		}
	}
	return nil
}

func RequestRedraw(widget Widget) {
	widget.widgetState().dirty = true
}

func z(widget Widget) int {
	var r int
	if parent := widget.widgetState().parent; parent != nil {
		r = z(parent)
	}
	r += widget.ZDelta()
	return r
}
