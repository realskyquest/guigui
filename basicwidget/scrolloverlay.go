// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package basicwidget

import (
	"image"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

func barMaxOpacity() int {
	return int(float64(ebiten.TPS()) / 6)
}

func barShowingTime() int {
	return ebiten.TPS()
}

type ScrollOverlay struct {
	guigui.DefaultWidget

	contentWidth  int
	contentHeight int
	offsetX       float64
	offsetY       float64

	lastWidth            int
	lastHeight           int
	lastCursorX          int
	lastCursorY          int
	lastWheelX           float64
	lastWheelY           float64
	lastOffsetX          float64
	lastOffsetY          float64
	draggingX            bool
	draggingY            bool
	draggingStartX       int
	draggingStartY       int
	draggingStartOffsetX float64
	draggingStartOffsetY float64
	onceUpdated          bool

	barOpacity     int
	barVisibleTime int

	contentSizeChanged bool

	onScroll func(offsetX, offsetY float64)
}

func (s *ScrollOverlay) SetOnScroll(f func(offsetX, offsetY float64)) {
	s.onScroll = f
}

func (s *ScrollOverlay) Reset() {
	s.offsetX = 0
	s.offsetY = 0
}

func (s *ScrollOverlay) SetContentSize(context *guigui.Context, contentWidth, contentHeight int) {
	if s.contentWidth == contentWidth && s.contentHeight == contentHeight {
		return
	}

	s.contentWidth = contentWidth
	s.contentHeight = contentHeight
	s.adjustOffset(context)
	if s.onceUpdated {
		s.contentSizeChanged = true
		guigui.RequestRedraw(s)
	}
}

func (s *ScrollOverlay) SetOffsetByDelta(context *guigui.Context, contentWidth, contentHeight int, dx, dy float64) {
	s.SetOffset(context, contentWidth, contentHeight, s.offsetX+dx, s.offsetY+dy)
}

func (s *ScrollOverlay) SetOffset(context *guigui.Context, contentWidth, contentHeight int, x, y float64) {
	s.SetContentSize(context, contentWidth, contentHeight)

	if s.offsetX == x && s.offsetY == y {
		return
	}
	s.offsetX = x
	s.offsetY = y
	s.adjustOffset(context)
	if s.onceUpdated {
		guigui.RequestRedraw(s)
	}
}

func (s *ScrollOverlay) isHovered(context *guigui.Context) bool {
	return context.IsWidgetHitAt(s, image.Pt(ebiten.CursorPosition()))
}

func (s *ScrollOverlay) setDragging(draggingX, draggingY bool) {
	if s.draggingX == draggingX && s.draggingY == draggingY {
		return
	}

	s.draggingX = draggingX
	s.draggingY = draggingY
}

func adjustedWheel() (float64, float64) {
	x, y := ebiten.Wheel()
	switch runtime.GOOS {
	case "darwin":
		x *= 2
		y *= 2
	}
	return x, y
}

func (s *ScrollOverlay) HandlePointingInput(context *guigui.Context) guigui.HandleInputResult {
	hovered := s.isHovered(context)
	if hovered {
		x, y := ebiten.CursorPosition()
		dx, dy := adjustedWheel()
		s.lastCursorX = x
		s.lastCursorY = y
		s.lastWheelX = dx
		s.lastWheelY = dy
	} else {
		s.lastCursorX = -1
		s.lastCursorY = -1
		s.lastWheelX = 0
		s.lastWheelY = 0
	}

	if !s.draggingX && !s.draggingY && hovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		hb, vb := s.barBounds(context)
		if image.Pt(x, y).In(hb) {
			s.setDragging(true, s.draggingY)
			s.draggingStartX = x
			s.draggingStartOffsetX = s.offsetX
		} else if image.Pt(x, y).In(vb) {
			s.setDragging(s.draggingX, true)
			s.draggingStartY = y
			s.draggingStartOffsetY = s.offsetY
		}
		if s.draggingX || s.draggingY {
			return guigui.HandleInputByWidget(s)
		}
	}

	if dx, dy := adjustedWheel(); dx != 0 || dy != 0 {
		s.setDragging(false, false)
	}

	if (s.draggingX || s.draggingY) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		var dx, dy float64
		if s.draggingX {
			dx = float64(x - s.draggingStartX)
		}
		if s.draggingY {
			dy = float64(y - s.draggingStartY)
		}
		if dx != 0 || dy != 0 {
			prevOffsetX := s.offsetX
			prevOffsetY := s.offsetY

			w, h := context.Size(s)
			barWidth, barHeight := s.barSize(context)
			if s.draggingX && barWidth > 0 && s.contentWidth-w > 0 {
				offsetPerPixel := float64(s.contentWidth-w) / (float64(w) - barWidth)
				s.offsetX = s.draggingStartOffsetX + float64(-dx)*offsetPerPixel
			}
			if s.draggingY && barHeight > 0 && s.contentHeight-h > 0 {
				offsetPerPixel := float64(s.contentHeight-h) / (float64(h) - barHeight)
				s.offsetY = s.draggingStartOffsetY + float64(-dy)*offsetPerPixel
			}
			s.adjustOffset(context)
			if prevOffsetX != s.offsetX || prevOffsetY != s.offsetY {
				if s.onScroll != nil {
					s.onScroll(s.offsetX, s.offsetY)
				}
				guigui.RequestRedraw(s)
			}
		}
		return guigui.HandleInputByWidget(s)
	}

	if (s.draggingX || s.draggingY) && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.setDragging(false, false)
	}

	if dx, dy := adjustedWheel(); dx != 0 || dy != 0 {
		if !hovered {
			return guigui.HandleInputResult{}
		}
		s.setDragging(false, false)

		prevOffsetX := s.offsetX
		prevOffsetY := s.offsetY
		s.offsetX += dx * 4 * context.Scale()
		s.offsetY += dy * 4 * context.Scale()
		s.adjustOffset(context)
		if prevOffsetX != s.offsetX || prevOffsetY != s.offsetY {
			if s.onScroll != nil {
				s.onScroll(s.offsetX, s.offsetY)
			}
			guigui.RequestRedraw(s)
			return guigui.HandleInputByWidget(s)
		}
		return guigui.HandleInputResult{}
	}

	return guigui.HandleInputResult{}
}

func (s *ScrollOverlay) CursorShape(context *guigui.Context) (ebiten.CursorShapeType, bool) {
	x, y := ebiten.CursorPosition()
	hb, vb := s.barBounds(context)
	if image.Pt(x, y).In(hb) || image.Pt(x, y).In(vb) {
		return ebiten.CursorShapeDefault, true
	}
	return 0, false
}

func (s *ScrollOverlay) Offset() (float64, float64) {
	return s.offsetX, s.offsetY
}

func (s *ScrollOverlay) adjustOffset(context *guigui.Context) {
	bounds := context.Bounds(s)

	// Adjust offsets.
	if s.offsetX > 0 {
		s.offsetX = 0
	}
	if s.offsetY > 0 {
		s.offsetY = 0
	}

	w := s.contentWidth - bounds.Dx()
	h := s.contentHeight - bounds.Dy()
	if w < 0 {
		s.offsetX = 0
	} else if s.offsetX < -float64(w) {
		s.offsetX = -float64(w)
	}
	if h < 0 {
		s.offsetY = 0
	} else if s.offsetY < -float64(h) {
		s.offsetY = -float64(h)
	}
}

func (s *ScrollOverlay) isBarVisible(context *guigui.Context) bool {
	if s.draggingX || s.draggingY {
		return true
	}
	if s.lastWheelX != 0 || s.lastWheelY != 0 {
		return true
	}
	if s.lastOffsetX != s.offsetX || s.lastOffsetY != s.offsetY {
		return true
	}

	bounds := context.Bounds(s)
	if s.contentWidth > bounds.Dx() && bounds.Max.Y-UnitSize(context) <= s.lastCursorY {
		return true
	}
	if s.contentHeight > bounds.Dy() && bounds.Max.X-UnitSize(context) <= s.lastCursorX {
		return true
	}
	return false
}

func (s *ScrollOverlay) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	w, h := context.Size(s)
	if s.lastWidth != w || s.lastHeight != h {
		s.adjustOffset(context)
		s.lastWidth = w
		s.lastHeight = h
	}

	context.SetOpacity(s, float64(s.barOpacity)/float64(barMaxOpacity())*3/4)
	return nil
}

func (s *ScrollOverlay) Update(context *guigui.Context) error {
	if s.contentSizeChanged {
		s.barVisibleTime = barShowingTime()
		s.contentSizeChanged = false
	}

	if s.isBarVisible(context) || (s.barVisibleTime == barShowingTime() && s.barOpacity < barMaxOpacity()) {
		if s.barOpacity < barMaxOpacity() {
			s.barOpacity++
			guigui.RequestRedraw(s)
		}
		s.barVisibleTime = barShowingTime()
	} else {
		if s.barVisibleTime > 0 {
			s.barVisibleTime--
		}
		if s.barVisibleTime == 0 && s.barOpacity > 0 {
			s.barOpacity--
			guigui.RequestRedraw(s)
		}
	}

	s.lastOffsetX = s.offsetX
	s.lastOffsetY = s.offsetY

	s.onceUpdated = true

	return nil
}

func (s *ScrollOverlay) Draw(context *guigui.Context, dst *ebiten.Image) {
	if s.barOpacity == 0 {
		return
	}

	barColor := draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.2)
	hb, vb := s.barBounds(context)

	// Show a horizontal bar.
	if !hb.Empty() {
		draw.DrawRoundedRect(context, dst, hb, barColor, RoundedCornerRadius(context))
	}

	// Show a vertical bar.
	if !vb.Empty() {
		draw.DrawRoundedRect(context, dst, vb, barColor, RoundedCornerRadius(context))
	}
}

func (s *ScrollOverlay) barWidth(scale float64) float64 {
	const scrollBarStrokeWidthInDIP = 8
	return scrollBarStrokeWidthInDIP * scale
}

func (s *ScrollOverlay) barSize(context *guigui.Context) (float64, float64) {
	bounds := context.Bounds(s)

	var w, h float64
	if s.contentWidth > bounds.Dx() {
		w = float64(bounds.Dx()) * float64(bounds.Dx()) / float64(s.contentWidth)
		if min := s.barWidth(context.Scale()); w < min {
			w = min
		}
	}
	if s.contentHeight > bounds.Dy() {
		h = float64(bounds.Dy()) * float64(bounds.Dy()) / float64(s.contentHeight)
		if min := s.barWidth(context.Scale()); h < min {
			h = min
		}
	}
	return w, h
}

func (s *ScrollOverlay) barBounds(context *guigui.Context) (image.Rectangle, image.Rectangle) {
	bounds := context.Bounds(s)

	offsetX, offsetY := s.Offset()
	barWidth, barHeight := s.barSize(context)

	padding := 2 * context.Scale()

	var horizontalBarBounds, verticalBarBounds image.Rectangle
	if s.contentWidth > bounds.Dx() {
		rate := -offsetX / float64(s.contentWidth-bounds.Dx())
		x0 := float64(bounds.Min.X) + rate*(float64(bounds.Dx())-barWidth)
		x1 := x0 + float64(barWidth)
		var y0, y1 float64
		if s.barWidth(context.Scale()) > float64(bounds.Dy())*0.3 {
			y0 = float64(bounds.Max.Y) - float64(bounds.Dy())*0.3
			y1 = float64(bounds.Max.Y)
		} else {
			y0 = float64(bounds.Max.Y) - padding - s.barWidth(context.Scale())
			y1 = float64(bounds.Max.Y) - padding
		}
		horizontalBarBounds = image.Rect(int(x0), int(y0), int(x1), int(y1))
	}
	if s.contentHeight > bounds.Dy() {
		rate := -offsetY / float64(s.contentHeight-bounds.Dy())
		y0 := float64(bounds.Min.Y) + rate*(float64(bounds.Dy())-barHeight)
		y1 := y0 + float64(barHeight)
		var x0, x1 float64
		if s.barWidth(context.Scale()) > float64(bounds.Dx())*0.3 {
			x0 = float64(bounds.Max.X) - float64(bounds.Dx())*0.3
			x1 = float64(bounds.Max.X)
		} else {
			x0 = float64(bounds.Max.X) - padding - s.barWidth(context.Scale())
			x1 = float64(bounds.Max.X) - padding
		}
		verticalBarBounds = image.Rect(int(x0), int(y0), int(x1), int(y1))
	}
	return horizontalBarBounds, verticalBarBounds
}
