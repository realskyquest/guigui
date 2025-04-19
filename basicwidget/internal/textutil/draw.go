// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package textutil

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DrawOptions struct {
	Options

	TextColor color.Color

	DrawSelection  bool
	SelectionStart int
	SelectionEnd   int
	SelectionColor color.Color

	DrawComposition          bool
	CompositionStart         int
	CompositionEnd           int
	CompositionActiveStart   int
	CompositionActiveEnd     int
	InactiveCompositionColor color.Color
	ActiveCompositionColor   color.Color
	CompositionBorderWidth   float32
}

func Draw(bounds image.Rectangle, dst *ebiten.Image, str string, options *DrawOptions) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(bounds.Min.X), float64(bounds.Min.Y))
	op.ColorScale.ScaleWithColor(options.TextColor)
	if dst.Bounds() != bounds {
		dst = dst.SubImage(bounds).(*ebiten.Image)
	}

	op.LineSpacing = options.LineHeight

	switch options.HorizontalAlign {
	case HorizontalAlignStart:
		op.PrimaryAlign = text.AlignStart
	case HorizontalAlignCenter:
		op.GeoM.Translate(float64(bounds.Dx())/2, 0)
		op.PrimaryAlign = text.AlignCenter
	case HorizontalAlignEnd:
		op.GeoM.Translate(float64(bounds.Dx()), 0)
		op.PrimaryAlign = text.AlignEnd
	}

	c := lineCount(bounds.Dx(), str, options.AutoWrap, options.Face)
	if c == 0 {
		return
	}
	height := options.LineHeight * float64(c)

	m := options.Face.Metrics()
	padding := (options.LineHeight - (m.HAscent + m.HDescent)) / 2
	op.GeoM.Translate(0, padding)
	switch options.VerticalAlign {
	case VerticalAlignTop:
	case VerticalAlignMiddle:
		op.GeoM.Translate(0, (float64(bounds.Dy())-height)/2)
	case VerticalAlignBottom:
		op.GeoM.Translate(0, float64(bounds.Dy())-height)
	}

	for pos, line := range lines(bounds.Dx(), str, options.AutoWrap, options.Face) {
		start := pos
		end := pos + len(line) - tailingLineBreakLen(line)

		if options.DrawSelection {
			if start <= options.SelectionEnd && end >= options.SelectionStart {
				start := max(start, options.SelectionStart)
				end := min(end, options.SelectionEnd)
				posStart0, posStart1, okStart0, okStart1 := TextPositionFromIndex(bounds.Dx(), str, start, &options.Options)
				posEnd0, posEnd1, okEnd0, okEnd1 := TextPositionFromIndex(bounds.Dx(), str, end, &options.Options)
				if (okStart0 || okStart1) && (okEnd0 || okEnd1) {
					posStart := posStart1
					if !okStart1 {
						posStart = posStart0
					}
					posEnd := posEnd0
					if !okEnd0 {
						posEnd = posEnd1
					}
					x := float32(posStart.X) + float32(bounds.Min.X)
					y := float32(posStart.Top) + float32(bounds.Min.Y)
					width := float32(posEnd.X - posStart.X)
					height := float32(posStart.Bottom - posStart.Top)
					vector.DrawFilledRect(dst, x, y, width, height, options.SelectionColor, false)
				}
			}
		}

		if options.DrawComposition {
			if start <= options.CompositionEnd && end >= options.CompositionStart {
				start := max(start, options.CompositionStart)
				end := min(end, options.CompositionEnd)
				posStart0, posStart1, okStart0, okStart1 := TextPositionFromIndex(bounds.Dx(), str, start, &options.Options)
				posEnd0, posEnd1, okEnd0, okEnd1 := TextPositionFromIndex(bounds.Dx(), str, end, &options.Options)
				if (okStart0 || okStart1) && (okEnd0 || okEnd1) {
					posStart := posStart1
					if !okStart1 {
						posStart = posStart0
					}
					posEnd := posEnd0
					if !okEnd0 {
						posEnd = posEnd1
					}
					x := float32(posStart.X) + float32(bounds.Min.X)
					y := float32(posStart.Bottom) + float32(bounds.Min.Y) - options.CompositionBorderWidth
					w := float32(posEnd.X - posStart.X)
					h := options.CompositionBorderWidth
					vector.DrawFilledRect(dst, x, y, w, h, options.InactiveCompositionColor, false)
				}
			}
			if start <= options.CompositionActiveEnd && end >= options.CompositionActiveStart {
				start := max(start, options.CompositionActiveStart)
				end := min(end, options.CompositionActiveEnd)
				posStart0, posStart1, okStart0, okStart1 := TextPositionFromIndex(bounds.Dx(), str, start, &options.Options)
				posEnd0, posEnd1, okEnd0, okEnd1 := TextPositionFromIndex(bounds.Dx(), str, end, &options.Options)
				if (okStart0 || okStart1) && (okEnd0 || okEnd1) {
					posStart := posStart1
					if !okStart1 {
						posStart = posStart0
					}
					posEnd := posEnd0
					if !okEnd0 {
						posEnd = posEnd1
					}
					x := float32(posStart.X) + float32(bounds.Min.X)
					y := float32(posStart.Bottom) + float32(bounds.Min.Y) - options.CompositionBorderWidth
					w := float32(posEnd.X - posStart.X)
					h := options.CompositionBorderWidth
					vector.DrawFilledRect(dst, x, y, w, h, options.ActiveCompositionColor, false)
				}
			}
		}

		// Draw the text.
		text.Draw(dst, line, options.Face, op)
		op.GeoM.Translate(0, options.LineHeight)
	}
}
