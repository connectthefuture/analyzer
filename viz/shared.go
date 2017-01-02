package viz

import (
	"fmt"
	"image/color"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

var defaultFont vg.Font

func init() {
	var err error
	plot.DefaultFont = "Helvetica"
	defaultFont, err = vg.MakeFont("Helvetica", 6)
	if err != nil {
		fmt.Println(err.Error())
	}
}

var Red = color.RGBA{163, 36, 61, 255}
var Blue = color.RGBA{34, 64, 95, 255}
var Black = color.RGBA{0, 0, 0, 255}

func LineStyle(color color.RGBA, width float64, dashes ...[]vg.Length) draw.LineStyle {
	ls := draw.LineStyle{
		Color: color,
		Width: vg.Points(width),
	}
	if len(dashes) != 0 {
		ls.Dashes = dashes[0]
	}
	return ls
}

var OrderedColors = []color.RGBA{
	{255, 0, 0, 255},
	{0, 200, 0, 255},
	{0, 0, 255, 255},
	{125, 0, 0, 255},
	{0, 125, 0, 255},
	{0, 0, 125, 255},
	{125, 125, 0, 255},
	{125, 0, 125, 255},
	{0, 125, 125, 255},
	{125, 125, 125, 255},
	{200, 200, 200, 255},
	{255, 125, 0, 255},
	{0, 125, 255, 255},
	{0, 0, 0, 255},
	{255, 0, 0, 255},
	{0, 200, 0, 255},
	{0, 0, 255, 255},
	{125, 0, 0, 255},
	{0, 125, 0, 255},
	{0, 0, 125, 255},
	{125, 125, 0, 255},
	{125, 0, 125, 255},
	{0, 125, 125, 255},
	{125, 125, 125, 255},
	{200, 200, 200, 255},
	{255, 125, 0, 255},
	{0, 125, 255, 255},
}

func OrderedColor(i int) color.RGBA {
	return OrderedColors[i%len(OrderedColors)]
}

var Dot = []vg.Length{vg.Points(1), vg.Points(2)}
var Dash = []vg.Length{vg.Points(4), vg.Points(4)}

func pathRectangle(top vg.Length, right vg.Length, bottom vg.Length, left vg.Length) vg.Path {
	p := vg.Path{}
	p.Move(vg.Point{left, top})
	p.Line(vg.Point{right, top})
	p.Line(vg.Point{right, bottom})
	p.Line(vg.Point{left, bottom})
	p.Close()
	return p
}
