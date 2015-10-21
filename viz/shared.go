package viz

import (
	"fmt"
	"image/color"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/vg"
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

var Red = color.RGBA{255, 0, 0, 255}
var Blue = color.RGBA{0, 0, 255, 255}

func LineStyle(color color.RGBA, width float64) plot.LineStyle {
	return plot.LineStyle{
		Color: color,
		Width: vg.Points(width),
	}
}

var OrderedColors = []color.RGBA{
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

var Dot = []vg.Length{vg.Points(1), vg.Points(4)}
var Dash = []vg.Length{vg.Points(1), vg.Points(1)}

func pathRectangle(top vg.Length, right vg.Length, bottom vg.Length, left vg.Length) vg.Path {
	p := vg.Path{}
	p.Move(left, top)
	p.Line(right, top)
	p.Line(right, bottom)
	p.Line(left, bottom)
	p.Close()
	return p
}
