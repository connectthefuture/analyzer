package viz

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
)

type VerticalLine struct {
	X float64
	plot.LineStyle
}

func NewVerticalLine(x float64) *VerticalLine {
	return &VerticalLine{
		X:         x,
		LineStyle: plotter.DefaultLineStyle,
	}
}

func (pts *VerticalLine) Plot(da plot.DrawArea, plt *plot.Plot) {
	trX, _ := plt.Transforms(&da)
	ps := make([]plot.Point, 2)
	ps[0].X = trX(pts.X)
	ps[1].X = ps[0].X

	ps[0].Y = da.Min.Y
	ps[1].Y = da.Max().Y

	da.StrokeLines(pts.LineStyle, da.ClipLinesXY(ps)...)
}
