package viz

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

type VerticalLine struct {
	X float64
	draw.LineStyle
}

func NewVerticalLine(x float64) *VerticalLine {
	return &VerticalLine{
		X:         x,
		LineStyle: plotter.DefaultLineStyle,
	}
}

func (pts *VerticalLine) Plot(da draw.Canvas, plt *plot.Plot) {
	trX, _ := plt.Transforms(&da)
	ps := make([]vg.Point, 2)
	ps[0].X = trX(pts.X)
	ps[1].X = ps[0].X

	ps[0].Y = da.Min.Y
	ps[1].Y = da.Max.Y

	da.StrokeLines(pts.LineStyle, da.ClipLinesXY(ps)...)
}
