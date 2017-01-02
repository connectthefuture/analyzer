package viz

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

type HorizontalLine struct {
	Y float64
	draw.LineStyle
}

func NewHorizontalLine(y float64) *HorizontalLine {
	return &HorizontalLine{
		Y:         y,
		LineStyle: plotter.DefaultLineStyle,
	}
}

func (pts *HorizontalLine) Plot(da draw.Canvas, plt *plot.Plot) {
	_, trY := plt.Transforms(&da)
	ps := make([]vg.Point, 2)

	ps[0].X = da.Min.X
	ps[1].X = da.Max.X

	ps[0].Y = trY(pts.Y)
	ps[1].Y = ps[0].Y

	da.StrokeLines(pts.LineStyle, da.ClipLinesXY(ps)...)
}
