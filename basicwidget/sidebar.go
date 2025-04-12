// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package basicwidget

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget/internal/draw"
)

type Sidebar struct {
	guigui.DefaultWidget

	scrollablePanel ScrollablePanel

	widthMinusDefault  int
	heightMinusDefault int
}

func (s *Sidebar) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	w, h := s.Size(context)
	s.scrollablePanel.SetSize(context, w, h)
	guigui.SetPosition(&s.scrollablePanel, guigui.Position(s))
	appender.AppendChildWidget(&s.scrollablePanel)

	return nil
}

func (s *Sidebar) SetContent(widget guigui.Widget) {
	s.scrollablePanel.SetContent(widget)
}

func (s *Sidebar) Draw(context *guigui.Context, dst *ebiten.Image) {
	dst.Fill(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.875))
	b := guigui.Bounds(s)
	b.Min.X = b.Max.X - int(1*context.Scale())
	dst.SubImage(b).(*ebiten.Image).Fill(draw.Color(context.ColorMode(), draw.ColorTypeBase, 0.85))
}

func defaultSidebarWidth(context *guigui.Context) (int, int) {
	return 6 * UnitSize(context), 6 * UnitSize(context)
}

func (s *Sidebar) Size(context *guigui.Context) (int, int) {
	dw, dh := defaultSidebarWidth(context)
	return s.widthMinusDefault + dw, s.heightMinusDefault + dh
}

func (s *Sidebar) SetSize(context *guigui.Context, width, height int) {
	dw, dh := defaultSidebarWidth(context)
	s.widthMinusDefault = width - dw
	s.heightMinusDefault = height - dh
}
