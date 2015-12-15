package viz

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
)

type HorizontalLine struct {
	Y float64
	plot.LineStyle
}

func NewHorizontalLine(y float64) *HorizontalLine {
	return &HorizontalLine{
		Y:         y,
		LineStyle: plotter.DefaultLineStyle,
	}
}

func (pts *HorizontalLine) Plot(da plot.DrawArea, plt *plot.Plot) {
	_, trY := plt.Transforms(&da)
	ps := make([]plot.Point, 2)

	ps[0].X = da.Min.X
	ps[1].X = da.Max().X

	ps[0].Y = trY(pts.Y)
	ps[1].Y = ps[0].Y

	da.StrokeLines(pts.LineStyle, da.ClipLinesXY(ps)...)
}
