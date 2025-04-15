// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package basicwidget

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type ScrollablePanel struct {
	guigui.DefaultWidget

	content      guigui.Widget
	scollOverlay ScrollOverlay
	border       scrollablePanelBorder
}

func (s *ScrollablePanel) SetContent(widget guigui.Widget) {
	s.content = widget
}

func (s *ScrollablePanel) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if s.content == nil {
		return nil
	}

	offsetX, offsetY := s.scollOverlay.Offset()
	appender.AppendChildWidgetWithPosition(s.content, context.Position(s).Add(image.Pt(int(offsetX), int(offsetY))))

	cw, ch := context.Size(s.content)
	s.scollOverlay.SetContentSize(context, cw, ch)
	appender.AppendChildWidgetWithBounds(&s.scollOverlay, context.Bounds(s))

	s.border.scrollOverlay = &s.scollOverlay
	appender.AppendChildWidgetWithPosition(&s.border, context.Position(s))

	return nil
}

type scrollablePanelBorder struct {
	guigui.DefaultWidget

	scrollOverlay *ScrollOverlay
}

func (s *scrollablePanelBorder) Draw(context *guigui.Context, dst *ebiten.Image) {
	// Render borders.
	strokeWidth := float32(1 * context.Scale())
	bounds := context.Bounds(s)
	x0 := float32(bounds.Min.X)
	x1 := float32(bounds.Max.X)
	y0 := float32(bounds.Min.Y)
	y1 := float32(bounds.Max.Y)
	offsetX, offsetY := s.scrollOverlay.Offset()
	if offsetX < 0 {
		vector.StrokeLine(dst, x0+strokeWidth/2, y0, x0+strokeWidth/2, y1, strokeWidth, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.85), false)
	}
	if offsetY < 0 {
		vector.StrokeLine(dst, x0, y0+strokeWidth/2, x1, y0+strokeWidth/2, strokeWidth, draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.85), false)
	}
}
