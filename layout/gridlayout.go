// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 Hajime Hoshi

package layout

import (
	"image"
	"iter"
)

type Size struct {
	typ     sizeType
	value   int
	content func(index int) int
}

type sizeType int

const (
	sizeTypeFixed sizeType = iota
	sizeTypeFlexible
	sizeTypeMaxContent
)

func FixedSize(value int) Size {
	return Size{
		typ:   sizeTypeFixed,
		value: value,
	}
}

func FlexibleSize(value int) Size {
	return Size{
		typ:   sizeTypeFlexible,
		value: value,
	}
}

func MaxContentSize(f func(index int) int) Size {
	return Size{
		typ:     sizeTypeMaxContent,
		value:   0,
		content: f,
	}
}

var (
	defaultWidths  = []Size{FlexibleSize(1)}
	defaultHeights = []Size{FlexibleSize(1)}
)

type GridLayout struct {
	Bounds    image.Rectangle
	Widths    []Size
	Heights   []Size
	ColumnGap int
	RowGap    int
}

func (g GridLayout) CellBounds() iter.Seq2[int, image.Rectangle] {
	return g.cellBounds(max(len(g.Widths), 1) * max(len(g.Heights), 1))
}

func (g GridLayout) RepeatingCellBounds() iter.Seq2[int, image.Rectangle] {
	return g.cellBounds(-1)
}

func (g GridLayout) cellBounds(count int) iter.Seq2[int, image.Rectangle] {
	return func(yield func(index int, bounds image.Rectangle) bool) {
		widths := g.Widths
		if len(widths) == 0 {
			widths = defaultWidths
		}
		heights := g.Heights
		if len(heights) == 0 {
			heights = defaultHeights
		}

		widthsInPixels := make([]int, len(widths))

		// Calculate widths in pixels.
		restW := g.Bounds.Dx()
		restW -= (len(widths) - 1) * g.ColumnGap
		if restW < 0 {
			restW = 0
		}
		var denomW int

		for i, width := range widths {
			switch width.typ {
			case sizeTypeFixed:
				widthsInPixels[i] = width.value
			case sizeTypeFlexible:
				widthsInPixels[i] = 0
				denomW += width.value
			case sizeTypeMaxContent:
				if count < 0 {
					panic("layout: MaxContentSize is not supported with infinite count")
				}
				widthsInPixels[i] = 0
				for j := range (count-1)/len(widths) + 1 {
					if j*len(widths)+i >= count {
						break
					}
					if width.content == nil {
						continue
					}
					widthsInPixels[i] = max(widthsInPixels[i], width.content(j*len(widths)+i))
				}
			}
			restW -= widthsInPixels[i]
		}

		if denomW > 0 {
			origRestW := restW
			for i, width := range widths {
				if width.typ != sizeTypeFlexible {
					continue
				}
				w := int(float64(origRestW) * float64(width.value) / float64(denomW))
				widthsInPixels[i] = w
				restW -= w
			}
			// TODO: Use a better algorithm to distribute the rest.
			for restW > 0 {
				for i := len(widthsInPixels) - 1; i >= 0; i-- {
					if widths[i].typ != sizeTypeFlexible {
						continue
					}
					widthsInPixels[i]++
					restW--
					if restW <= 0 {
						break
					}
				}
				if restW <= 0 {
					break
				}
			}
		}

		y := g.Bounds.Min.Y
		heightsInPixels := make([]int, len(heights))
		var widgetIdx int
		for widgetBaseIdx := 0; count < 0 || widgetBaseIdx < count; widgetBaseIdx += len(widths) * len(heights) {
			// Calculate hights in pixels.
			// This is needed for each loop since the index starts with widgetBaseIdx for sizeTypeMaxContent.
			restH := g.Bounds.Dy()
			if restH < 0 {
				restH = 0
			}
			restH -= (len(heights) - 1) * g.RowGap
			var denomH int

			for j, height := range heights {
				switch height.typ {
				case sizeTypeFixed:
					heightsInPixels[j] = height.value
				case sizeTypeFlexible:
					heightsInPixels[j] = 0
					denomH += height.value
				case sizeTypeMaxContent:
					heightsInPixels[j] = 0
					for i := range widths {
						if count >= 0 && j*len(widths)+i >= count {
							break
						}
						if height.content == nil {
							continue
						}
						heightsInPixels[j] = max(heightsInPixels[j], height.content(widgetBaseIdx+j*len(widths)+i))
					}
				}
				restH -= heightsInPixels[j]
			}

			if denomH > 0 {
				origRestH := restH
				for j, height := range heights {
					if height.typ != sizeTypeFlexible {
						continue
					}
					h := int(float64(origRestH) * float64(height.value) / float64(denomH))
					heightsInPixels[j] = h
					restH -= h
				}
				for restH > 0 {
					for j := len(heightsInPixels) - 1; j >= 0; j-- {
						if heights[j].typ != sizeTypeFlexible {
							continue
						}
						heightsInPixels[j]++
						restH--
						if restH <= 0 {
							break
						}
					}
					if restH <= 0 {
						break
					}
				}
			}

			for j := 0; j < len(heights); j++ {
				x := g.Bounds.Min.X
				for i := 0; i < len(widths); i++ {
					bounds := image.Rect(x, y, x+widthsInPixels[i], y+heightsInPixels[j])
					if !yield(widgetIdx, bounds) {
						return
					}
					x += widthsInPixels[i]
					x += g.ColumnGap
					widgetIdx++
					if count >= 0 && widgetIdx >= count {
						return
					}
				}
				y += heightsInPixels[j]
				y += g.RowGap
			}
		}
	}
}
